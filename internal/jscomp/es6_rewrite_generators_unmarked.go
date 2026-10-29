// Port of com/google/javascript/jscomp/Es6RewriteGenerators.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var generatorArgumentsID = rhino.RhinoStringPool_InternID("arguments")

var generatorUnmarkedVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_FOR)
	s.Add(rhino.Token_FOR_IN)
	s.Add(rhino.Token_FOR_OF)
	s.Add(rhino.Token_FOR_AWAIT_OF)
	s.Add(rhino.Token_WHILE)
	s.Add(rhino.Token_DO)
	s.Add(rhino.Token_SWITCH)
	s.Add(rhino.Token_THIS)
	s.Add(rhino.Token_RETURN)
	s.Add(rhino.Token_NAME)
	s.Add(rhino.Token_VAR)
	return s
}()

var _ nodeTraversal_visitTokens = (*generatorUnmarkedNodeTranspiler)(nil)

func (*generatorUnmarkedNodeTranspiler) visitTokens() *ast.TokenSet {
	return &generatorUnmarkedVisitTokens
}

// generatorUnmarkedNodeTranspiler adjusts YIELD-free nodes to run correctly
// inside a generator state-machine program:
//   - moving `var` declarations into the outer hoist scope;
//   - transpiling `return` statements;
//   - transpiling `break` and `continue` statements;
//   - transpiling references to `this` and `arguments`.
type generatorUnmarkedNodeTranspiler struct {
	pass *Es6RewriteGenerators
	fn   *generatorFunctionTranspiler
	ctx  *generatorTranspilationContext

	// Count of enclosing statements within this unmarked subtree that a bare
	// `break` / `continue` addresses locally (so it needs no state-machine jump).
	breakSuppressors    int
	continueSuppressors int
}

func newGeneratorUnmarkedNodeTranspiler(
	pass *Es6RewriteGenerators,
	fn *generatorFunctionTranspiler,
	ctx *generatorTranspilationContext,
) *generatorUnmarkedNodeTranspiler {
	return &generatorUnmarkedNodeTranspiler{pass: pass, fn: fn, ctx: ctx}
}

func (u *generatorUnmarkedNodeTranspiler) ShouldTraverse(nodeTraversal *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsGeneratorSafe() {
		n.SetGeneratorSafe(false)
		return false
	}

	jutil.CheckState(!n.IsGeneratorMarker())
	jutil.CheckStateMsg(!n.IsSuper(), "Reference to SUPER is not supported")

	if NodeUtil_IsLoopStructure(n) {
		u.continueSuppressors++
		u.breakSuppressors++
	} else if _, ok := ast.AsSwitch(n); ok {
		u.breakSuppressors++
	}

	if jmp, ok := ast.AsJump(n); ok {
		if lbl, hasLabel := jmp.Label(); hasLabel {
			u.visitNamedBreakContinue(n, lbl)
		} else {
			u.visitBreakContinue(n)
		}
		return false
	}

	_, isFn := ast.AsFunction(n)
	return !isFn
}

func (u *generatorUnmarkedNodeTranspiler) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if NodeUtil_IsLoopStructure(n) {
		u.continueSuppressors--
		u.breakSuppressors--
	} else if _, ok := ast.AsSwitch(n); ok {
		u.breakSuppressors--
	} else if n.IsThis() {
		u.visitThis(n)
	} else if n.IsReturn() {
		u.visitReturn(n)
	} else if n.IsName() && n.GetStringID() == generatorArgumentsID {
		u.visitArguments(n)
	} else if n.IsVar() {
		if parent.IsVanillaFor() {
			u.visitVanillaForLoopVar(n)
		} else if parent.IsForIn() {
			u.visitForInLoopVar(n)
		} else {
			// NOTE: for-of loops are transpiled away before this pass
			u.visitVarStatement(n)
		}
	}
}

func (u *generatorUnmarkedNodeTranspiler) visitReturn(n rhino.Node) {
	n.AddChildToFront(u.ctx.returnExpression(n, n.RemoveFirstChild()))
}

func (u *generatorUnmarkedNodeTranspiler) visitNamedBreakContinue(n rhino.Node, lbl ast.LabelName) {
	labelIdent := lbl.Ident()
	if u.ctx.namedLabels.ContainsKey(labelIdent) {
		cases := u.ctx.namedLabels.Get(labelIdent)
		var caseSection *generatorStateCase
		if n.IsBreak() {
			caseSection = cases.breakCase
		} else {
			caseSection = cases.continueCase
		}
		u.ctx.replaceBreakContinueWithJump(n, caseSection, u.breakSuppressors)
	}
}

func (u *generatorUnmarkedNodeTranspiler) visitBreakContinue(n rhino.Node) {
	var caseSection *generatorStateCase
	if n.IsBreak() && u.breakSuppressors == 0 {
		if len(u.ctx.breakCases) == 0 {
			panic(jutil.NewNoSuchElementException())
		}
		caseSection = u.ctx.breakCases[len(u.ctx.breakCases)-1]
	}
	if n.IsContinue() && u.continueSuppressors == 0 {
		if len(u.ctx.continueCases) == 0 {
			panic(jutil.NewNoSuchElementException())
		}
		caseSection = u.ctx.continueCases[len(u.ctx.continueCases)-1]
	}
	if caseSection != nil {
		u.ctx.replaceBreakContinueWithJump(n, caseSection, u.breakSuppressors)
	}
}

func (u *generatorUnmarkedNodeTranspiler) visitThis(n rhino.Node) {
	newThis := u.pass.astFactory.CreateName(u.ctx.getScopedName(Es6RewriteGenerators_GENERATOR_THIS), AstFactory_Type_(n))
	n.ReplaceWith(newThis)
	if !u.ctx.thisReferenceFound {
		varNode := rhino.IR_VarWithValue(newThis.CloneNode().Srcref(n), n).Srcref(u.fn.newGeneratorHoistBlock)
		u.fn.hoistVarNode(varNode)
		u.ctx.thisReferenceFound = true
	}
}

func (u *generatorUnmarkedNodeTranspiler) visitArguments(n rhino.Node) {
	newArguments :=
		u.pass.astFactory.CreateName(u.ctx.getScopedName(Es6RewriteGenerators_GENERATOR_ARGUMENTS), AstFactory_Type_(n)).Srcref(n)
	n.ReplaceWith(newArguments)
	if !u.ctx.argumentsReferenceFound {
		varNode := rhino.IR_VarWithValue(newArguments.CloneNode(), n).Srcref(u.fn.newGeneratorHoistBlock)
		u.fn.hoistVarNode(varNode)
		u.ctx.argumentsReferenceFound = true
	}
}

func (u *generatorUnmarkedNodeTranspiler) visitVarStatement(varStatement rhino.Node) {
	commaExpression := u.extractAssignmentsToCommaExpression(varStatement)
	if commaExpression.IsNil() {
		varStatement.Detach()
	} else {
		varStatement.ReplaceWith(u.pass.astFactory.ExprResult(commaExpression))
	}
	u.fn.hoistVarNode(varStatement)
}

func (u *generatorUnmarkedNodeTranspiler) visitVanillaForLoopVar(varDeclaration rhino.Node) {
	commaExpression := u.extractAssignmentsToCommaExpression(varDeclaration)
	if commaExpression.IsNil() {
		varDeclaration.ReplaceWith(rhino.IR_Empty())
	} else {
		varDeclaration.ReplaceWith(commaExpression)
	}
	u.fn.hoistVarNode(varDeclaration)
}

func (u *generatorUnmarkedNodeTranspiler) visitForInLoopVar(varDeclaration rhino.Node) {
	varName := varDeclaration.GetOnlyChild()
	jutil.CheckStateMsg1(!varName.HasChildren(), "%s", varName)
	clonedVarName := varName.CloneNode().SetJSDocInfo(nil)
	varDeclaration.ReplaceWith(clonedVarName)
	u.fn.hoistVarNode(varDeclaration)
}

func (u *generatorUnmarkedNodeTranspiler) extractAssignmentsToCommaExpression(varDeclaration rhino.Node) rhino.Node {
	astFactory := u.pass.astFactory
	var commaExpression rhino.Node
	for varName := varDeclaration.GetFirstChild(); !varName.IsNil(); varName = varName.GetNext() {
		if varName.HasChildren() {
			copiedVarName := varName.CloneNode().SetJSDocInfo(nil)
			assign :=
				astFactory.
					CreateAssign(copiedVarName, varName.RemoveFirstChild()).
					Srcref(varName)
			if commaExpression.IsNil() {
				commaExpression = assign
			} else {
				commaExpression = astFactory.CreateComma(commaExpression, assign).Srcref(assign)
			}
		}
	}
	return commaExpression
}
