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
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// generatorLabelMap is a lazy native Go map from label name to its break/continue
// state cases. It allocates nothing until a generator contains labeled statements
// wrapping yields.
type generatorLabelMap map[string]generatorLabelCases

func (m *generatorLabelMap) Put(k string, v generatorLabelCases) (generatorLabelCases, bool) {
	if *m == nil {
		*m = make(generatorLabelMap, 4)
	}
	old, existed := (*m)[k]
	(*m)[k] = v
	return old, existed
}

func (m generatorLabelMap) Get(k string) generatorLabelCases {
	return m[k]
}

func (m generatorLabelMap) GetOk(k string) (generatorLabelCases, bool) {
	v, ok := m[k]
	return v, ok
}

func (m generatorLabelMap) ContainsKey(k string) bool {
	_, ok := m[k]
	return ok
}

func (m generatorLabelMap) Remove(k string) (generatorLabelCases, bool) {
	v, ok := m[k]
	if ok {
		delete(m, k)
	}
	return v, ok
}

func (m generatorLabelMap) IsEmpty() bool {
	return len(m) == 0
}

func (m generatorLabelMap) Size() int {
	return len(m)
}

// generatorStringSet is a lazy native Go string set used to deduplicate hoisted
// catch variable declarations in source encounter order.
type generatorStringSet map[string]struct{}

func (s *generatorStringSet) Add(k string) bool {
	if _, exists := (*s)[k]; exists {
		return false
	}
	if *s == nil {
		*s = make(generatorStringSet, 4)
	}
	(*s)[k] = struct{}{}
	return true
}

func (s generatorStringSet) Contains(k string) bool {
	_, ok := s[k]
	return ok
}

func (s generatorStringSet) Remove(k string) bool {
	if _, ok := s[k]; ok {
		delete(s, k)
		return true
	}
	return false
}

func (s generatorStringSet) IsEmpty() bool {
	return len(s) == 0
}

func (s generatorStringSet) Size() int {
	return len(s)
}

// generatorTranspilationContext holds the state machine context used during
// generator function transpilation.
type generatorTranspilationContext struct {
	pass *Es6RewriteGenerators
	fn   *generatorFunctionTranspiler

	namedLabels   generatorLabelMap
	breakCases    []*generatorStateCase
	continueCases []*generatorStateCase

	catchCases   []generatorCatchCase
	finallyCases []*generatorStateCase
	catchNames   generatorStringSet

	// All "case" sections that will be added to generator program.
	allCases []*generatorStateCase

	// All "break" nodes that exit from the generator primary switch statement.
	switchBreaks []rhino.Node

	// A virtual case that indicates end of program.
	programEndCase *generatorStateCase

	// Next case ID to assign.
	caseIdCounter int

	// Switch case currently being populated with transpiled instructions.
	currentCase *generatorStateCase

	// Counter for the number of finally blocks currently enclosing the cursor.
	nestedFinallyBlockCount int
	thisReferenceFound      bool
	argumentsReferenceFound bool
}

func newGeneratorTranspilationContext(pass *Es6RewriteGenerators, fn *generatorFunctionTranspiler) *generatorTranspilationContext {
	c := &generatorTranspilationContext{
		pass: pass,
		fn:   fn,
	}
	c.programEndCase = newGeneratorStateCase(pass, c)
	jutil.CheckState(c.programEndCase.id == 0)
	c.currentCase = newGeneratorStateCase(pass, c)
	jutil.CheckState(c.currentCase.id == 1)
	c.allCases = append(c.allCases, c.currentCase)
	return c
}

// optimizeCaseIds shortcuts jump chains and merges unreferenced/embeddable cases.
func (c *generatorTranspilationContext) optimizeCaseIds() {
	jutil.CheckStateMsg(len(c.allCases) != 0, "[]")

	for _, currentCase := range c.allCases {
		if currentCase.jumpTo != nil {
			for currentCase.jumpTo.jumpTo != nil {
				currentCase.jumpTo = currentCase.jumpTo.jumpTo
			}

			if !currentCase.embedInto.IsNil() && len(currentCase.references) == 1 {
				currentCase.jumpTo.embedInto = currentCase.embedInto
			}
			currentCase.embedInto = rhino.NilNode

			for _, reference := range currentCase.references {
				reference.SetDouble(float64(currentCase.jumpTo.id))
			}
			currentCase.jumpTo.references = append(currentCase.jumpTo.references, currentCase.references...)
			currentCase.references = currentCase.references[:0]
		}
	}

	kept := c.allCases[:1]
	prevCase := c.allCases[0]
	jutil.CheckState(prevCase.id == 1)
	for _, currentCase := range c.allCases[1:] {
		if len(currentCase.references) == 0 {
			jutil.CheckState(currentCase.embedInto.IsNil())
			if prevCase.mayFallThrough {
				prevCase.caseBlock.AddChildrenToBack(currentCase.caseBlock.RemoveChildren())
				prevCase.mayFallThrough = currentCase.mayFallThrough
			}
			continue
		}
		if !currentCase.embedInto.IsNil() {
			jutil.CheckState(currentCase.jumpTo == nil)
			if len(currentCase.references) == 1 && !currentCase.mayFallThrough {
				currentCase.embedInto.ReplaceWith(currentCase.caseBlock)
				continue
			}
		}
		if prevCase.jumpTo == currentCase {
			jutil.CheckState(prevCase.mayFallThrough)
			jutil.CheckState(!prevCase.caseBlock.HasChildren())
			jutil.CheckState(currentCase.jumpTo == nil)

			prevCase.caseBlock.AddChildrenToBack(currentCase.caseBlock.RemoveChildren())
			prevCase.mayFallThrough = currentCase.mayFallThrough
			for _, reference := range currentCase.references {
				reference.SetDouble(float64(prevCase.id))
			}
			prevCase.jumpTo = currentCase.jumpTo
			prevCase.references = append(prevCase.references, currentCase.references...)
			continue
		}
		kept = append(kept, currentCase)
		prevCase = currentCase
	}
	clear(c.allCases[len(kept):])
	c.allCases = kept
}

// eliminateSwitchBreaks replaces `...; break;` with `return ...;` when rendering without a switch.
func (c *generatorTranspilationContext) eliminateSwitchBreaks() {
	for _, breakNode := range c.switchBreaks {
		prevStatement := breakNode.GetPrevious()
		jutil.CheckState(!prevStatement.IsNil())
		exprStmt := ast.MustExprStmt(prevStatement)
		prevStatement.ReplaceWith(rhino.IR_ReturnNodeWithExpr(exprStmt.Expr().Node().Detach()))
		breakNode.Detach()
	}
	c.switchBreaks = nil
}

func (c *generatorTranspilationContext) allCasesRemoveFirst() *generatorStateCase {
	first := c.allCases[0]
	c.allCases = c.allCases[1:]
	return first
}

func (c *generatorTranspilationContext) finalizeTransformation(generatorBody rhino.Node) {
	astFactory := c.pass.astFactory
	c.optimizeCaseIds()

	if len(c.allCases) == 2 || len(c.allCases) == 3 {
		generatorBody.AddChildToBack(
			astFactory.
				CreateIf(
					astFactory.CreateEq(
						c.callContextMethod(
							generatorBody,
							"getNextAddressJsc",
							AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_NUMBER_TYPE, colors.StandardColors_NUMBER)),
						astFactory.CreateNumber(1)),
					c.allCasesRemoveFirst().caseBlock).
				SrcrefTreeIfMissing(generatorBody))
	}

	if len(c.allCases) == 2 {
		generatorBody.AddChildToBack(
			astFactory.
				CreateIf(
					astFactory.CreateNe(
						c.callContextMethod(
							generatorBody,
							"getNextAddressJsc",
							AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_NUMBER_TYPE, colors.StandardColors_NUMBER)),
						astFactory.CreateNumber(float64(c.allCases[1].id))),
					c.allCasesRemoveFirst().caseBlock).
				SrcrefTreeIfMissing(generatorBody))
	}

	if len(c.allCases) == 1 {
		generatorBody.AddChildrenToBack(c.allCasesRemoveFirst().caseBlock.RemoveChildren())
		c.eliminateSwitchBreaks()
		return
	}

	switchNode :=
		rhino.IR_SwitchNode(
			c.callContextMethod(
				generatorBody,
				"getNextAddressJsc",
				AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_NUMBER_TYPE, colors.StandardColors_NUMBER))).
			Srcref(generatorBody)
	generatorBody.AddChildToBack(switchNode)
	switchBody := switchNode.GetSecondChild().Srcref(generatorBody)

	for _, currentCase := range c.allCases {
		switchBody.AddChildToBack(currentCase.createCaseNode())
	}
	c.allCases = nil
}

func (c *generatorTranspilationContext) checkStateIsEmpty() {
	jutil.CheckState(c.namedLabels.IsEmpty())
	jutil.CheckState(len(c.breakCases) == 0)
	jutil.CheckState(len(c.continueCases) == 0)
	jutil.CheckState(len(c.catchCases) == 0)
	jutil.CheckState(len(c.finallyCases) == 0)
	jutil.CheckState(c.nestedFinallyBlockCount == 0)
	jutil.CheckState(len(c.allCases) == 0)
}

func (c *generatorTranspilationContext) transpileUnmarkedBlock(block rhino.Node) {
	if block.HasChildren() {
		NodeTraversal_Traverse(c.pass.compiler, block, newGeneratorUnmarkedNodeTranspiler(c.pass, c.fn, c))
		for block.HasChildren() {
			c.writeGeneratedNode(block.RemoveFirstChild())
		}
	}
}

func (c *generatorTranspilationContext) writeGeneratedNode(n rhino.Node) {
	c.currentCase.addNode(n)
}

func (c *generatorTranspilationContext) writeGeneratedNodeAndBreak(n rhino.Node) {
	c.writeGeneratedNode(n)
	c.writeGeneratedNode(c.createBreakNodeFor(n))
	c.currentCase.mayFallThrough = false
}

func (c *generatorTranspilationContext) createCase() *generatorStateCase {
	return newGeneratorStateCase(c.pass, c)
}

func (c *generatorTranspilationContext) maybeCreateCase(other *generatorStateCase) *generatorStateCase {
	if other != nil {
		return other
	}
	return c.createCase()
}

func (c *generatorTranspilationContext) getJsContextNameNode(sourceNode rhino.Node) rhino.Node {
	return c.pass.astFactory.
		CreateName(c.getScopedName(Es6RewriteGenerators_GENERATOR_CONTEXT), c.pass.generatorContextType()).
		Srcref(sourceNode)
}

func (c *generatorTranspilationContext) getScopedName(name string) string {
	return name + "$" + c.fn.uniqueId
}

func (c *generatorTranspilationContext) getContextField(sourceNode rhino.Node, fieldName string) rhino.Node {
	return c.pass.astFactory.
		CreateGetPropWithUnknownType(c.getJsContextNameNode(sourceNode), fieldName).
		Srcref(sourceNode)
}

func (c *generatorTranspilationContext) callContextMethod(
	sourceNode rhino.Node, methodName string, type_ AstFactory_Type, args ...rhino.Node) rhino.Node {
	contextField := c.getContextField(sourceNode, methodName)
	return c.pass.astFactory.CreateCall(contextField, type_, args...).Srcref(sourceNode)
}

func (c *generatorTranspilationContext) callContextMethodResult(
	sourceNode rhino.Node, methodName string, type_ AstFactory_Type, args ...rhino.Node) rhino.Node {
	return c.pass.astFactory.
		ExprResult(c.callContextMethod(sourceNode, methodName, type_, args...)).
		Srcref(sourceNode)
}

func (c *generatorTranspilationContext) returnContextMethod(
	sourceNode rhino.Node, methodName string, type_ AstFactory_Type, args ...rhino.Node) rhino.Node {
	return c.pass.astFactory.
		CreateReturn(c.callContextMethod(sourceNode, methodName, type_, args...)).
		Srcref(sourceNode)
}

func (c *generatorTranspilationContext) createBreakNodeFor(preBreak rhino.Node) rhino.Node {
	breakNode := rhino.IR_BreakNode().Srcref(preBreak)
	c.switchBreaks = append(c.switchBreaks, breakNode)
	return breakNode
}

func (c *generatorTranspilationContext) createJumpToNode(section *generatorStateCase, sourceNode rhino.Node) rhino.Node {
	return c.returnContextMethod(
		sourceNode, "jumpTo", AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID), section.getNumber(sourceNode))
}

func (c *generatorTranspilationContext) writeJumpTo(section *generatorStateCase, sourceNode rhino.Node) {
	c.currentCase.jumpToCase(
		section, c.createJumpToBlock(section /* allowEmbedding= */, false, sourceNode))
}

func (c *generatorTranspilationContext) createJumpToBlock(section *generatorStateCase, allowEmbedding bool, sourceNode rhino.Node) rhino.Node {
	jutil.CheckState(section.embedInto.IsNil())
	jumpBlock :=
		c.pass.astFactory.
			CreateBlock(
				c.callContextMethodResult(
					sourceNode,
					"jumpTo",
					AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID),
					section.getNumber(sourceNode)),
				c.createBreakNodeFor(sourceNode)).
			Srcref(sourceNode)
	if allowEmbedding {
		section.embedInto = jumpBlock
	}
	return jumpBlock
}

func (c *generatorTranspilationContext) replaceBreakContinueWithJump(sourceNode rhino.Node, section *generatorStateCase, breakSuppressors int) {
	var jumpMethod string
	if c.nestedFinallyBlockCount > 0 {
		jumpMethod = "jumpThroughFinallyBlocks"
	} else if len(c.finallyCases) == 0 || c.finallyCasesGetFirst().id < section.id {
		jumpMethod = "jumpTo"
	} else {
		jutil.CheckState(c.finallyCasesGetFirst().id != section.id)
		jumpMethod = "jumpThroughFinallyBlocks"
	}
	if breakSuppressors == 0 {
		c.callContextMethodResult(
			sourceNode,
			jumpMethod,
			AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID),
			section.getNumber(sourceNode)).
			InsertBefore(sourceNode)
		sourceNode.ReplaceWith(c.createBreakNodeFor(sourceNode))
	} else {
		sourceNode.ReplaceWith(
			c.returnContextMethod(
				sourceNode,
				jumpMethod,
				AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID),
				section.getNumber(sourceNode)))
	}
}

func (c *generatorTranspilationContext) yieldValue(
	expression rhino.Node, jumpToSection *generatorStateCase, sourceNode rhino.Node) {
	var exprArg rhino.Node
	if expression.IsNil() {
		exprArg = c.pass.astFactory.CreateUndefinedValue().SrcrefTree(sourceNode)
	} else {
		exprArg = expression
	}
	c.writeGeneratedNode(
		c.returnContextMethod(
			sourceNode, "yield", AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN), exprArg, jumpToSection.getNumber(sourceNode)))
	c.currentCase.mayFallThrough = false
}

func (c *generatorTranspilationContext) yieldAll(expression rhino.Node, jumpToSection *generatorStateCase, sourceNode rhino.Node) {
	c.writeGeneratedNode(
		c.returnContextMethod(
			sourceNode,
			"yieldAll",
			AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN),
			expression,
			jumpToSection.getNumber(sourceNode)))
	c.currentCase.mayFallThrough = false
}

func (c *generatorTranspilationContext) returnExpression(sourceNode rhino.Node, expression rhino.Node) rhino.Node {
	if expression.IsNil() {
		return c.callContextMethod(sourceNode, "return", AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID))
	}
	return c.callContextMethod(
		sourceNode, "return", AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID), expression)
}

func (c *generatorTranspilationContext) yieldResult(sourceNode rhino.Node) rhino.Node {
	return c.callContextMethod(sourceNode, "getYieldResultJsc", AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN))
}

func (c *generatorTranspilationContext) catchCasesGetFirst() *generatorCatchCase {
	if len(c.catchCases) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	return &c.catchCases[len(c.catchCases)-1]
}

func (c *generatorTranspilationContext) catchCasesRemoveFirst() generatorCatchCase {
	first := *c.catchCasesGetFirst()
	c.catchCases = c.catchCases[:len(c.catchCases)-1]
	return first
}

func (c *generatorTranspilationContext) finallyCasesGetFirst() *generatorStateCase {
	if len(c.finallyCases) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	return c.finallyCases[len(c.finallyCases)-1]
}

func (c *generatorTranspilationContext) finallyCasesRemoveFirst() *generatorStateCase {
	first := c.finallyCasesGetFirst()
	c.finallyCases = c.finallyCases[:len(c.finallyCases)-1]
	return first
}

func (c *generatorTranspilationContext) addCatchFinallyCases(catchCase *generatorStateCase, finallyCase *generatorStateCase) {
	if finallyCase != nil {
		if len(c.catchCases) != 0 {
			c.catchCasesGetFirst().finallyBlocks++
		}
		c.finallyCases = append(c.finallyCases, finallyCase)
	}
	if catchCase != nil {
		c.catchCases = append(c.catchCases, generatorCatchCase{catchCase: catchCase})
	}
}

func (c *generatorTranspilationContext) getNextCatchCase() *generatorStateCase {
	if len(c.catchCases) > 0 {
		top := c.catchCases[len(c.catchCases)-1]
		if top.finallyBlocks == 0 {
			return top.catchCase
		}
	}
	return nil
}

func (c *generatorTranspilationContext) getNextFinallyCase() *generatorStateCase {
	if len(c.finallyCases) == 0 {
		return nil
	}
	return c.finallyCasesGetFirst()
}

func (c *generatorTranspilationContext) removeCatchFinallyCases(catchCase *generatorStateCase, finallyCase *generatorStateCase) {
	if catchCase != nil {
		lastCatch := c.catchCasesRemoveFirst()
		jutil.CheckState(lastCatch.finallyBlocks == 0)
		jutil.CheckState(lastCatch.catchCase == catchCase)
	}
	if finallyCase != nil {
		if len(c.catchCases) != 0 {
			top := c.catchCasesGetFirst()
			top.finallyBlocks--
			jutil.CheckState(top.finallyBlocks >= 0)
		}
		lastFinally := c.finallyCasesRemoveFirst()
		jutil.CheckState(lastFinally == finallyCase)
	}
}

func (c *generatorTranspilationContext) enterTryBlock(catchCase *generatorStateCase, finallyCase *generatorStateCase, sourceNode rhino.Node) {
	c.addCatchFinallyCases(catchCase, finallyCase)

	var methodName string
	var args []rhino.Node
	if catchCase == nil {
		methodName = "setFinallyBlock"
		args = append(args, finallyCase.getNumber(sourceNode))
	} else {
		methodName = "setCatchFinallyBlocks"
		args = append(args, catchCase.getNumber(sourceNode))
		if finallyCase != nil {
			args = append(args, finallyCase.getNumber(sourceNode))
		}
	}
	c.writeGeneratedNode(
		c.callContextMethodResult(
			sourceNode,
			methodName,
			AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID),
			args...))
}

func (c *generatorTranspilationContext) leaveTryBlock(catchCase *generatorStateCase, endCase *generatorStateCase, sourceNode rhino.Node) {
	c.removeCatchFinallyCases(catchCase, nil)
	args := []rhino.Node{endCase.getNumber(sourceNode)}
	if nextCatchCase := c.getNextCatchCase(); nextCatchCase != nil {
		args = append(args, nextCatchCase.getNumber(sourceNode))
	}
	c.writeGeneratedNodeAndBreak(
		c.callContextMethodResult(
			sourceNode,
			"leaveTryBlock",
			AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID),
			args...))
}

func (c *generatorTranspilationContext) enterCatchBlock(finallyCase *generatorStateCase, exceptionName rhino.Node) {
	astFactory := c.pass.astFactory
	excName := ast.MustName(exceptionName)
	c.addCatchFinallyCases(nil, finallyCase)

	nextCatchCase := c.getNextCatchCase()

	if c.catchNames.Add(excName.Ident()) {
		c.fn.hoistVarNode(rhino.IR_Var(exceptionName.CloneNode()).Srcref(exceptionName))
	}

	var args []rhino.Node
	if nextCatchCase != nil {
		args = append(args, nextCatchCase.getNumber(exceptionName))
	}

	enterCatchBlockCall :=
		c.callContextMethod(
			exceptionName,
			"enterCatchBlock",
			AstFactory_Type_OfColor(colors.StandardColors_UNKNOWN),
			args...)
	exceptionName.SetColor(enterCatchBlockCall.GetColor())
	c.writeGeneratedNode(
		astFactory.
			ExprResult(
				astFactory.
					CreateAssign(exceptionName, enterCatchBlockCall).
					Srcref(exceptionName)).
			Srcref(exceptionName))
}

func (c *generatorTranspilationContext) leaveCatchBlock(finallyCase *generatorStateCase, sourceNode rhino.Node) {
	if finallyCase != nil {
		c.removeCatchFinallyCases(nil, finallyCase)
		c.writeJumpTo(finallyCase, sourceNode)
	}
}

func (c *generatorTranspilationContext) enterFinallyBlock(
	catchCase *generatorStateCase, finallyCase *generatorStateCase, sourceNode rhino.Node) {
	astFactory := c.pass.astFactory
	c.removeCatchFinallyCases(catchCase, finallyCase)

	nextCatchCase := c.getNextCatchCase()
	nextFinallyCase := c.getNextFinallyCase()

	var args []rhino.Node
	if c.nestedFinallyBlockCount == 0 {
		if nextCatchCase != nil || nextFinallyCase != nil {
			if nextCatchCase == nil {
				args = append(args, astFactory.CreateNumber(0).Srcref(sourceNode))
			} else {
				args = append(args, nextCatchCase.getNumber(sourceNode))
			}
			if nextFinallyCase != nil {
				args = append(args, nextFinallyCase.getNumber(sourceNode))
			}
		}
	} else {
		if nextCatchCase == nil {
			args = append(args, astFactory.CreateNumber(0).Srcref(sourceNode))
		} else {
			args = append(args, nextCatchCase.getNumber(sourceNode))
		}
		if nextFinallyCase == nil {
			args = append(args, astFactory.CreateNumber(0).Srcref(sourceNode))
		} else {
			args = append(args, nextFinallyCase.getNumber(sourceNode))
		}
		args = append(args, astFactory.CreateNumber(float64(c.nestedFinallyBlockCount)).Srcref(sourceNode))
	}

	c.writeGeneratedNode(
		c.callContextMethodResult(
			sourceNode,
			"enterFinallyBlock",
			AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID),
			args...))

	c.nestedFinallyBlockCount++
}

func (c *generatorTranspilationContext) leaveFinallyBlock(endCase *generatorStateCase, sourceNode rhino.Node) {
	args := []rhino.Node{endCase.getNumber(sourceNode)}
	c.nestedFinallyBlockCount--
	if c.nestedFinallyBlockCount != 0 {
		args = append(args, c.pass.astFactory.CreateNumber(float64(c.nestedFinallyBlockCount)).Srcref(sourceNode))
	}

	c.writeGeneratedNodeAndBreak(
		c.callContextMethodResult(
			sourceNode,
			"leaveFinallyBlock",
			AstFactory_Type_OfColor(colors.StandardColors_NULL_OR_VOID),
			args...))
}

func (c *generatorTranspilationContext) switchCaseTo(caseSection *generatorStateCase) {
	c.currentCase.willFollowBy(caseSection)
	c.allCases = append(c.allCases, caseSection)
	c.currentCase = caseSection
}

func (c *generatorTranspilationContext) pushLabels(labelNames []rhino.Node, breakCase *generatorStateCase, continueCase *generatorStateCase) {
	for _, labelName := range labelNames {
		lbl := ast.MustLabelName(labelName)
		c.namedLabels.Put(lbl.Ident(), generatorLabelCases{breakCase: breakCase, continueCase: continueCase})
	}
}

func (c *generatorTranspilationContext) popLabels(labelNames []rhino.Node) {
	for _, labelName := range labelNames {
		lbl := ast.MustLabelName(labelName)
		c.namedLabels.Remove(lbl.Ident())
	}
}

func (c *generatorTranspilationContext) pushBreakContext(breakCase *generatorStateCase) {
	c.breakCases = append(c.breakCases, breakCase)
}

func (c *generatorTranspilationContext) pushBreakContinueContext(breakCase *generatorStateCase, continueCase *generatorStateCase) {
	c.pushBreakContext(breakCase)
	c.continueCases = append(c.continueCases, continueCase)
}

func (c *generatorTranspilationContext) popBreakContext() {
	if len(c.breakCases) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	c.breakCases = c.breakCases[:len(c.breakCases)-1]
}

func (c *generatorTranspilationContext) popBreakContinueContext() {
	c.popBreakContext()
	if len(c.continueCases) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	c.continueCases = c.continueCases[:len(c.continueCases)-1]
}

// generatorStateCase is a single state-machine case section in a generator program.
type generatorStateCase struct {
	pass *Es6RewriteGenerators
	ctx  *generatorTranspilationContext

	id             int
	caseBlock      rhino.Node
	references     []rhino.Node
	jumpTo         *generatorStateCase
	embedInto      rhino.Node
	mayFallThrough bool
}

func newGeneratorStateCase(pass *Es6RewriteGenerators, ctx *generatorTranspilationContext) *generatorStateCase {
	id := ctx.caseIdCounter
	ctx.caseIdCounter++
	return &generatorStateCase{
		pass:           pass,
		ctx:            ctx,
		id:             id,
		caseBlock:      pass.astFactory.CreateBlock().Srcref(ctx.fn.originalGeneratorBody),
		mayFallThrough: true,
	}
}

func (cs *generatorStateCase) createCaseNode() rhino.Node {
	return rhino.IR_CaseNode(cs.pass.astFactory.CreateNumber(float64(cs.id)).Srcref(cs.caseBlock), cs.caseBlock).
		Srcref(cs.caseBlock)
}

func (cs *generatorStateCase) getNumber(sourceNode rhino.Node) rhino.Node {
	if cs.jumpTo != nil {
		return cs.jumpTo.getNumber(sourceNode)
	}
	node := cs.pass.astFactory.CreateNumber(float64(cs.id)).Srcref(sourceNode)
	cs.references = append(cs.references, node)
	return node
}

func (cs *generatorStateCase) jumpToCase(other *generatorStateCase, jumpBlock rhino.Node) {
	jutil.CheckState(jumpBlock.IsBlock())
	jutil.CheckState(cs.jumpTo == nil)
	cs.willFollowBy(other)
	cs.caseBlock.AddChildrenToBack(jumpBlock.RemoveChildren())
	cs.mayFallThrough = false
}

func (cs *generatorStateCase) willFollowBy(other *generatorStateCase) {
	if cs.jumpTo == nil && !cs.caseBlock.HasChildren() {
		jutil.CheckState(other.jumpTo == nil)
		cs.jumpTo = other
	}
}

func (cs *generatorStateCase) addNode(n rhino.Node) {
	jutil.CheckState(cs.jumpTo == nil)
	jutil.CheckState(rhino.IR_MayBeStatement(n))
	cs.caseBlock.AddChildToBack(n)
}

type generatorCatchCase struct {
	catchCase     *generatorStateCase
	finallyBlocks int
}

type generatorLabelCases struct {
	breakCase    *generatorStateCase
	continueCase *generatorStateCase
}
