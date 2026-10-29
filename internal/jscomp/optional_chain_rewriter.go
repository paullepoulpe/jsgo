// Port of com/google/javascript/jscomp/OptionalChainRewriter.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/RewriteOptionalChainingOperator.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// OptionalChainRewriter rewrites a single optional chain as one or more
// nested hook expressions.
//
// Optional chains contained in OPTCHAIN_GETELEM indices or OPTCHAIN_CALL
// arguments are not rewritten.
//
// Example:
//
//	a?.b[obj?.index]?.c?.(obj?.arg)
//	// becomes
//	(tmp0 = a) == null
//	    ? void 0
//	    : (tmp1 = tmp0.b[obj?.index]) == null
//	        ? void 0
//	        : (tmp2 = tmp1.c) == null
//	            ? void 0
//	            : tmp2.call(tmp1, obj?.arg);
//
// The unit tests for this class are in RewriteOptionalChainingOperatorTest,
// because it's most convenient to test this class as part of the
// transpilation pass that uses it.
type OptionalChainRewriter struct {
	compiler          AbstractCompiler
	astFactory        *AstFactory
	tmpVarNameCreator OptionalChainRewriter_TmpVarNameCreator
	// If a scope is provided, newly created variables will be declared in that scope
	scope           *Scope
	chainParent     rhino.Node
	wholeChain      rhino.Node
	insertionPoint  rhino.Node
	deletesToDelete []rhino.Node // ArrayDeque used as a FIFO queue
}

// OptionalChainRewriter_TmpVarNameCreator creates unique names to be used
// for temporary variables.
//
// Java: OptionalChainRewriter.TmpVarNameCreator
type OptionalChainRewriter_TmpVarNameCreator func() string

// OptionalChainRewriter_Builder is OptionalChainRewriter.Builder.
type OptionalChainRewriter_Builder struct {
	compiler          AbstractCompiler
	astFactory        *AstFactory
	tmpVarNameCreator OptionalChainRewriter_TmpVarNameCreator
	scope             *Scope
}

// Java: OptionalChainRewriter.Builder(AbstractCompiler)
func NewOptionalChainRewriter_Builder(compiler AbstractCompiler) *OptionalChainRewriter_Builder {
	b := &OptionalChainRewriter_Builder{}
	b.compiler = jutil.CheckNotNull(compiler)
	b.astFactory = compiler.CreateAstFactory()
	return b
}

// Java: OptionalChainRewriter.Builder.setTmpVarNameCreator(OptionalChainRewriter.TmpVarNameCreator)
func (b *OptionalChainRewriter_Builder) SetTmpVarNameCreator(
	tmpVarNameCreator OptionalChainRewriter_TmpVarNameCreator) *OptionalChainRewriter_Builder {
	if tmpVarNameCreator == nil {
		panic(jutil.NewNullPointerException())
	}
	b.tmpVarNameCreator = tmpVarNameCreator
	return b
}

// SetScope optionally sets a scope in which the rewriter will declare all
// temporary variables
//
// Java: OptionalChainRewriter.Builder.setScope(Scope)
func (b *OptionalChainRewriter_Builder) SetScope(scope *Scope) *OptionalChainRewriter_Builder {
	b.scope = jutil.CheckNotNull(scope)
	return b
}

// Build: wholeChain is the last Node in the optional chain. Parent of all
// the rest.
//
// Java: OptionalChainRewriter.Builder.build(Node)
func (b *OptionalChainRewriter_Builder) Build(wholeChain rhino.Node) *OptionalChainRewriter {
	return NewOptionalChainRewriter(b, wholeChain)
}

// Java: OptionalChainRewriter.builder(AbstractCompiler)
func OptionalChainRewriter_Builder_(compiler AbstractCompiler) *OptionalChainRewriter_Builder {
	return NewOptionalChainRewriter_Builder(compiler)
}

// Java: OptionalChainRewriter(OptionalChainRewriter.Builder, Node)
func NewOptionalChainRewriter(builder *OptionalChainRewriter_Builder, wholeChain rhino.Node) *OptionalChainRewriter {
	o := &OptionalChainRewriter{}
	initOptionalChainRewriter(o, builder, wholeChain)
	return o
}

func initOptionalChainRewriter(o *OptionalChainRewriter, builder *OptionalChainRewriter_Builder, wholeChain rhino.Node) {
	// This class will only operate on an entire chain.
	jutil.CheckArgumentMsg1(NodeUtil_IsEndOfFullOptChain(wholeChain), "%s", wholeChain)
	if builder.tmpVarNameCreator == nil {
		panic(jutil.NewNullPointerException())
	}
	o.compiler = builder.compiler
	o.astFactory = builder.astFactory
	o.tmpVarNameCreator = builder.tmpVarNameCreator
	o.scope = builder.scope
	o.wholeChain = wholeChain
	o.chainParent = rhino.CheckNotNullMsg1(wholeChain.GetParent(), "%s", wholeChain)
	o.deletesToDelete = o.deletesToDelete[:0]
	// Insert temporaries just before the enclosing statement.
	enclosingStatement := NodeUtil_GetEnclosingStatement(wholeChain)
	if _, isLabel := ast.AsLabel(enclosingStatement.GetParent()); isLabel {
		// Do not insert temporaries between a label name and its statement.
		// e.g. `{ label: for (const a of b?.c) {...}`
		// AST shape is
		// BLOCK
		//      LABEL // parent of enclosingStatement - insert above this
		//          LABEL_NAME
		//          enclosingStatement
		o.insertionPoint = enclosingStatement.GetParent()
	} else {
		o.insertionPoint = enclosingStatement
	}
}

// Rewrite rewrites the optional chain as a hook with temporary variables
// introduced as needed.
//
// Java: OptionalChainRewriter.rewrite()
func (o *OptionalChainRewriter) Rewrite() {
	jutil.CheckStateMsg1(NodeUtil_IsOptChainNode(o.wholeChain), "already rewritten: %s", o.wholeChain)

	// `first?.start.second?.start`
	// We search from the end of the chain and push the start nodes onto the stack, so the first
	// one ends up on top.
	var startNodeStackBuf [8]rhino.Node
	startNodeStack := startNodeStackBuf[:0] // ArrayDeque used as a stack; the top is the last element
	subchainEnd := o.wholeChain
	for NodeUtil_IsOptChainNode(subchainEnd) {
		subchainStart := NodeUtil_GetStartOfOptChainSegment(subchainEnd)
		startNodeStack = append(startNodeStack, subchainStart)
		subchainEnd = subchainStart.GetFirstChild()
	}

	jutil.CheckState(len(startNodeStack) != 0)
	// Each time we rewrite the initial segment of the chain, the remaining chain gets wrapped
	// in a hook statement like `(tmp0 = a.b) == null ? void 0 : tmp0.rest?.of.chain?.()`,
	// So wholeChain ends up more deeply nested on each rewrite.
	// We only care about the top-most replacement here.
	optChainReplacement := o.RewriteInitialSegment(optionalChainRewriterPop(&startNodeStack), o.wholeChain)
	for len(startNodeStack) != 0 {
		o.RewriteInitialSegment(optionalChainRewriterPop(&startNodeStack), o.wholeChain)
	}

	// Handle a non-optional call to an optional chain that ends in an element or property
	// access.
	// `(a?.optional.chain)(arg1)`
	// Writing JavaScript code like this is a bad idea, but it might get automatically
	// generated, so we must handle it.
	// The optional chain could evaluate to `undefined`, which we then try to call as a
	// function. However, if it isn't undefined, we have to preserve the correct `this` value
	// for the call.
	if parentCall, ok := ast.AsCall(o.chainParent); ok &&
		// If the call is explicitly a free call (e.g. `(0, a?.b)()`), we should not preserve
		// `this`.
		!o.chainParent.GetBooleanProp(rhino.Node_FREE_CALL) &&
		// The chain will have been replaced by optChainReplacement during the rewriting above.
		parentCall.Callee().Is(optChainReplacement) {
		// The wholeChain variable will still point to the rewritten final Node of the
		// chain. It will no longer be optional.
		if finalGet, ok := ast.AsNormalGet(o.wholeChain); ok {
			thisValue := finalGet.Object().Node()
			tmpThisNode := o.GetSubExprNameNode(thisValue)
			optChainReplacement.Detach()
			o.chainParent.AddChildToFront(tmpThisNode)
			dotCallNode :=
				o.astFactory.
					CreateGetPropWithUnknownType(optChainReplacement, "call").
					SrcrefTreeIfMissing(optChainReplacement)
			o.chainParent.AddChildToFront(dotCallNode)
		}
	}

	// Report changes here; chainParent can get deleted below this.
	// E.g. code `(p = a?.b)=>{ return p}` changes to `let tmp0; (p = (tmp0=a)==null?void
	// 0:tmp0.b)=>{return p}`
	// This requires recording scope changes at two places:
	// 1. the function scope changed from rewriting to HOOK (recorded here),
	// 2. SCRIPT scope changed due to the `let tmp0` inserted (recorded in `declareTempVarName`
	// where declarations are created).
	o.compiler.ReportChangeToEnclosingScope(o.chainParent)

	if o.chainParent.IsDelProp() {
		// With rewriting above,
		// `delete a?.b?.c`
		//  synthesizes additional deletes
		// `delete (tmp0 = a) == null ? true : delete (tmp1 = tmp0.b) == null ? true : delete tmp1.c;`
		//  ^^^^^^                             ^^^^^^
		//
		// But we must generate only:
		// `(tmp0 = a) == null ? true : (tmp1 = tmp0.b) == null ? true : delete tmp1.c;`
		// That is, the preceding deletes for every hook must be removed.

		for len(o.deletesToDelete) != 0 {
			delete_ := o.deletesToDelete[0]
			o.deletesToDelete = o.deletesToDelete[1:]
			delUnary := ast.MustUnary(delete_)
			hook, ok := delUnary.Operand().AsHook()
			jutil.CheckStateMsg1(ok, "%s", delete_)
			hookNode := hook.Detach().Node()
			delete_.ReplaceWith(hookNode)
			o.compiler.ReportChangeToEnclosingScope(hookNode)
		}
	}

	// Transpilation of the optional chain adds `let` declarations for temporary variables.
	// NOTE: If this class is being used before transpilation, it's OK to use `let`, since it will
	// be transpiled away, if necessary. If it is being used after transpilation, then using `let`
	// must be OK, because optional chains weren't transpiled away and `let` existed before they
	// did.
	enclosingScript := NodeUtil_GetEnclosingScript(o.insertionPoint)
	NodeUtil_AddFeatureToScript(enclosingScript, parser.FeatureSet_Feature_LET_DECLARATIONS, o.compiler)
}

// optionalChainRewriterPop is ArrayDeque.pop() on a stack whose top is the
// last element.
//
// Go-only: language-adapter — ArrayDeque.
func optionalChainRewriterPop(stack *[]rhino.Node) rhino.Node {
	s := *stack
	if len(s) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	top := s[len(s)-1]
	*stack = s[:len(s)-1]
	return top
}

// RewriteInitialSegment rewrites the first part of a possibly-multi-part
// optional chain.
//
// e.g.
//
//	a()?.b.c?.d;
//	// becomes
//	let tmp0;
//	(tmp0 = a()) == null
//	    ? void 0
//	    : tmp0.b.c?d;
//
// If this optional chain is under a delete, the l-r rewriting must
// synthesize another `delete` to ensure the next chain(if present), knows
// that it must delete the `fullChainEnd`.
//
// e.g.
//
//	delete a?.b.c?.d;
//	// becomes
//	let tmp0;
//	(tmp0 = a()) == null
//	    ? true
//	    : delete tmp0.b.c?d;
//
// fullChainStart is the very first `?.` node; fullChainEnd the very last
// optional chain node. Returns the hook expression that replaced the chain.
//
// Java: OptionalChainRewriter.rewriteInitialSegment(Node, Node)
func (o *OptionalChainRewriter) RewriteInitialSegment(fullChainStart rhino.Node, fullChainEnd rhino.Node) rhino.Node {
	// `receiverNode?.restOfChain`
	receiverNode := fullChainStart.GetFirstChild()
	// for `a?.b.c?.d`, this will be `a?.b.c`, because the NodeUtil method finds the end
	// of the sub-chain, not the full chain.
	initialChainEnd := NodeUtil_GetEndOfOptChainSegment(fullChainStart)

	// Is this optional chain under delete
	isBeingDeleted := fullChainEnd.GetParent().IsDelProp()
	if isBeingDeleted {
		o.deletesToDelete = append(o.deletesToDelete, fullChainEnd.GetParent())
	}

	// If the receiver is an optional chain, we weren't really given the start of a full
	// chain.
	jutil.CheckArgumentMsg1(!NodeUtil_IsOptChainNode(receiverNode), "%s", receiverNode)

	// change the initial chain's nodes to be non-optional
	NodeUtil_ConvertToNonOptionalChainSegment(initialChainEnd)

	placeholder := rhino.IR_Empty()
	fullChainEnd.ReplaceWith(placeholder)
	// NOTE: convertToNonOptionalChain() above will have made the chain start
	// and all the other nodes in the first segment of the chain non-optional,
	// so ast.AsCall(fullChainStart) is the right test here.
	if recvGet, ok := ast.AsNormalGet(receiverNode); ok &&
		fullChainStart.IsCall() &&
		// Check !FREE_CALL to ensure we don't treat free calls on property
		// access (e.g. `(0, expr.prop)?.()`) as method calls.
		!fullChainStart.GetBooleanProp(rhino.Node_FREE_CALL) {
		// `expr.prop?.(x).y`
		// Needs to become
		// `(t1 = (t0 = expr).prop) == null ? void 0 : t1.call(t0, x).y`
		thisValue := recvGet.Object().Node()
		tmpThisNode := o.GetSubExprNameNode(thisValue)
		tmpReceiverNode := o.GetSubExprNameNode(receiverNode)
		receiverNode = fullChainStart.RemoveFirstChild()
		fullChainStart.AddChildToFront(tmpThisNode)
		fullChainStart.AddChildToFront(
			o.astFactory.
				CreateGetPropWithUnknownType(tmpReceiverNode, "call").
				SrcrefTreeIfMissing(receiverNode))
	} else {
		// `expr?.x.y`
		// needs to become
		// `((t0 = expr) == null) ? void 0 : t0.x.y`
		tmpReceiverNode := o.GetSubExprNameNode(receiverNode)
		receiverNode = fullChainStart.GetFirstChild()
		receiverNode.ReplaceWith(tmpReceiverNode)
	}
	// Arguments of createHook, in Java's evaluation order.
	cond := o.astFactory.CreateEq(receiverNode, o.astFactory.CreateNull())
	var trueValue, falseValue rhino.Node
	if isBeingDeleted {
		trueValue = o.astFactory.CreateBoolean(true)
		falseValue = o.astFactory.CreateDelProp(fullChainEnd)
	} else {
		trueValue = o.astFactory.CreateUndefinedValue()
		falseValue = fullChainEnd
	}
	optChainReplacement :=
		o.astFactory.
			CreateHook(
				cond,
				trueValue,
				falseValue).
			SrcrefTreeIfMissing(fullChainEnd)

	placeholder.ReplaceWith(optChainReplacement)

	return optChainReplacement
}

// GetSubExprNameNode: given an expression node, declare a temporary
// variable to hold that expression and replace the expression with
// `(tmp = expr)`.
//
// e.g. `subExpr.moreExpr` becomes `(tmp = subExpr).moreExpr`, and `let
// tmp;` gets inserted before the enclosing statement of this optional
// chain.
//
// subExpr is the sub expression Node. Returns a detached NAME node for the
// temporary variable name and with source info and type matching
// `subExpr`, that may be inserted where needed.
//
// Java: OptionalChainRewriter.getSubExprNameNode(Node)
func (o *OptionalChainRewriter) GetSubExprNameNode(subExpr rhino.Node) rhino.Node {
	tempVarName := o.DeclareTempVarName(subExpr)
	placeholder := rhino.IR_Empty()
	subExpr.ReplaceWith(placeholder)
	subExpr.RemoveProp(rhino.Node_DIRECT_EVAL)
	replacement := o.astFactory.CreateAssignToName(tempVarName, subExpr).SrcrefTreeIfMissing(subExpr)
	placeholder.ReplaceWith(replacement)
	return ast.MustAssign(replacement).Target().Node().CloneNode()
}

// DeclareTempVarName declares a temporary variable name that will be used
// to hold the given value.
//
// The generated declaration has no assignment, it's just `let tmp;`.
//
// valueNode is a node from which to copy the source info and type to be
// used for the new variable. Returns the name used for the new temporary
// variable.
//
// Java: OptionalChainRewriter.declareTempVarName(Node)
func (o *OptionalChainRewriter) DeclareTempVarName(valueNode rhino.Node) string {
	tempVarName := o.tmpVarNameCreator()
	declarationStatement :=
		o.astFactory.CreateSingleLetNameDeclaration(tempVarName).SrcrefTree(valueNode)
	nameNode, _ := ast.MustNameDeclaration(declarationStatement).OnlyDeclarator()
	nameNode.Node().SetInferredConstantVar(true)
	declarationStatement.InsertBefore(o.insertionPoint)
	o.compiler.ReportChangeToEnclosingScope(declarationStatement)
	if o.scope != nil {
		o.scope.Declare(tempVarName, nameNode.Node() /* input= */, nil)
	}
	return tempVarName
}

var rewriteOptChainVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_SCRIPT)
	s.Add(rhino.Token_OPTCHAIN_GETPROP)
	s.Add(rhino.Token_OPTCHAIN_GETELEM)
	s.Add(rhino.Token_OPTCHAIN_CALL)
	return s
}()

// RewriteOptionalChainingOperator replaces the ES2020 `?.` operator with
// conditional (? :).
type RewriteOptionalChainingOperator struct {
	compiler                AbstractCompiler
	uniqueIdSupplier        *UniqueIdSupplier
	customTmpVarNameCreator OptionalChainRewriter_TmpVarNameCreator
	callback                RewriteOptionalChainingOperator_TranspilationCallback
}

// NewRewriteOptionalChainingOperator is the constructor to be used for
// actual transpilation.
//
// Java: RewriteOptionalChainingOperator(AbstractCompiler)
func NewRewriteOptionalChainingOperator(compiler AbstractCompiler) *RewriteOptionalChainingOperator {
	return &RewriteOptionalChainingOperator{
		compiler:         compiler,
		uniqueIdSupplier: compiler.GetUniqueIdSupplier(),
	}
}

// NewRewriteOptionalChainingOperatorWithTmpVarNameCreator is the
// constructor for testing.
//
// Java: RewriteOptionalChainingOperator(AbstractCompiler, TmpVarNameCreator)
func NewRewriteOptionalChainingOperatorWithTmpVarNameCreator(
	compiler AbstractCompiler, tmpVarNameCreator OptionalChainRewriter_TmpVarNameCreator) *RewriteOptionalChainingOperator {
	return &RewriteOptionalChainingOperator{
		compiler:                compiler,
		customTmpVarNameCreator: jutil.CheckNotNull(tmpVarNameCreator),
	}
}

// Java: RewriteOptionalChainingOperator.process(Node, Node)
func (r *RewriteOptionalChainingOperator) Process(externs rhino.Node, root rhino.Node) {
	r.callback.initForProcess(r)
	NodeTraversal_TraverseRoots(r.compiler, &r.callback, externs, root)
	TranspilationPasses_MaybeMarkFeatureAsTranspiledAway(r.compiler, root, parser.FeatureSet_Feature_OPTIONAL_CHAINING)
}

// RewriteOptionalChainingOperator_TranspilationCallback locates and
// transpiles all optional chains.
type RewriteOptionalChainingOperator_TranspilationCallback struct {
	outer           *RewriteOptionalChainingOperator
	rewriterBuilder *OptionalChainRewriter_Builder
	builderVal      OptionalChainRewriter_Builder
	currentInput    *CompilerInput
	boundTmpCreator OptionalChainRewriter_TmpVarNameCreator
	optChains       []OptionalChainRewriter
	optChainsBuf    [4]OptionalChainRewriter
}

// Java: RewriteOptionalChainingOperator.TranspilationCallback()
func NewRewriteOptionalChainingOperator_TranspilationCallback(outer *RewriteOptionalChainingOperator) *RewriteOptionalChainingOperator_TranspilationCallback {
	c := &RewriteOptionalChainingOperator_TranspilationCallback{}
	c.initForProcess(outer)
	return c
}

func (c *RewriteOptionalChainingOperator_TranspilationCallback) initForProcess(outer *RewriteOptionalChainingOperator) {
	c.outer = outer
	c.builderVal = OptionalChainRewriter_Builder{
		compiler:   jutil.CheckNotNull(outer.compiler),
		astFactory: outer.compiler.CreateAstFactory(),
	}
	c.rewriterBuilder = &c.builderVal
	if outer.customTmpVarNameCreator != nil {
		c.rewriterBuilder.SetTmpVarNameCreator(outer.customTmpVarNameCreator)
	} else {
		if c.boundTmpCreator == nil {
			c.boundTmpCreator = c.nextTmpVarName
		}
		c.rewriterBuilder.SetTmpVarNameCreator(c.boundTmpCreator)
	}
	if c.optChains == nil {
		c.optChains = c.optChainsBuf[:0]
	} else {
		clear(c.optChains)
		c.optChains = c.optChains[:0]
	}
}

func (c *RewriteOptionalChainingOperator_TranspilationCallback) nextTmpVarName() string {
	return "$jscomp$optchain$tmp" + c.outer.uniqueIdSupplier.GetUniqueId(c.currentInput)
}

func (*RewriteOptionalChainingOperator_TranspilationCallback) noVisitLeaves() {}
func (*RewriteOptionalChainingOperator_TranspilationCallback) visitTokens() *ast.TokenSet {
	return &rewriteOptChainVisitTokens
}

// Java: RewriteOptionalChainingOperator.TranspilationCallback.shouldTraverse(NodeTraversal, Node, Node)
func (c *RewriteOptionalChainingOperator_TranspilationCallback) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		scriptFeatures := NodeUtil_GetFeatureSetOfScript(n)
		// ensures that the pass early exits if script does not contain the feature
		if scriptFeatures != nil && !scriptFeatures.ContainsFeature(parser.FeatureSet_Feature_OPTIONAL_CHAINING) {
			return false
		}
		c.currentInput = t.GetInput()
	}
	return true
}

// Java: RewriteOptionalChainingOperator.TranspilationCallback.visit(NodeTraversal, Node, Node)
func (c *RewriteOptionalChainingOperator_TranspilationCallback) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsScript() {
		// We transpile all of the optional chains in a single script as a batch because,
		// rewriting changes the AST in ways that could interfere with traversal
		// if we change the chains as we visit them.
		if len(c.optChains) > 0 {
			for i := range c.optChains {
				c.optChains[i].Rewrite()
			}
			clear(c.optChains)
			c.optChains = c.optChains[:0]
		}
	} else if NodeUtil_IsEndOfFullOptChain(n) {
		c.optChains = append(c.optChains, OptionalChainRewriter{})
		initOptionalChainRewriter(&c.optChains[len(c.optChains)-1], c.rewriterBuilder, n)
	}
}
