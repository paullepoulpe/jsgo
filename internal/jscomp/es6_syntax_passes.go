// Port of com/google/javascript/jscomp/Es6RewriteBlockScopedFunctionDeclaration.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/Es6RewriteRestParameters.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/SubstituteEs6Syntax.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/RewriteObjectSpread.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/RewriteGlobalDeclarationsForTryCatchWrapping.java @ closure-compiler v20261021

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

/*
 * Copyright 2015 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RewriteBlockScopedFunctionDeclaration rewrites block-scoped function
// declarations as "let"s. This pass must happen before
// Es6RewriteBlockScopedDeclaration, which rewrites "let" to "var".
type Es6RewriteBlockScopedFunctionDeclaration struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler
}

// Java: Es6RewriteBlockScopedFunctionDeclaration.transpiledFeatures
var es6RewriteBlockScopedFunctionDeclaration_transpiledFeatures = parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION)

var es6RewriteBlockScopedFnVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_FUNCTION)
	return s
}()

func (*Es6RewriteBlockScopedFunctionDeclaration) noVisitLeaves() {}
func (*Es6RewriteBlockScopedFunctionDeclaration) visitTokens() *ast.TokenSet {
	return &es6RewriteBlockScopedFnVisitTokens
}

// Java: Es6RewriteBlockScopedFunctionDeclaration(AbstractCompiler)
func NewEs6RewriteBlockScopedFunctionDeclaration(compiler AbstractCompiler) *Es6RewriteBlockScopedFunctionDeclaration {
	e := &Es6RewriteBlockScopedFunctionDeclaration{}
	e.compiler = compiler
	return e
}

// Java: Es6RewriteBlockScopedFunctionDeclaration.process(Node, Node)
func (e *Es6RewriteBlockScopedFunctionDeclaration) Process(externs rhino.Node, root rhino.Node) {
	TranspilationPasses_ProcessTranspile(e.compiler, root, es6RewriteBlockScopedFunctionDeclaration_transpiledFeatures, e)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, es6RewriteBlockScopedFunctionDeclaration_transpiledFeatures)
}

// Java: Es6RewriteBlockScopedFunctionDeclaration.visit(NodeTraversal, Node, Node)
func (e *Es6RewriteBlockScopedFunctionDeclaration) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if n.IsFunction() &&
		!parent.IsNil() &&
		parent.IsBlock() &&
		!parent.GetParent().IsFunction() {
		// Only consider declarations (all expressions have non-block parents) that are not directly
		// within a function or top-level.
		e.VisitBlockScopedFunctionDeclaration(t, n, parent)
	}
}

// VisitBlockScopedFunctionDeclaration rewrites the function declaration
// from:
//
//	function f() {}
//	FUNCTION
//	  NAME x
//	  PARAM_LIST
//	  BLOCK
//
// to
//
//	let f = function() {};
//	LET
//	  NAME f
//	    FUNCTION
//	      NAME (w/ empty string)
//	      PARAM_LIST
//	      BLOCK
//
// This is similar to Normalize.NormalizeStatements#rewriteFunctionDeclaration
// but rewrites to "let" instead of "var".
//
// Java: Es6RewriteBlockScopedFunctionDeclaration.visitBlockScopedFunctionDeclaration(NodeTraversal, Node, Node)
func (e *Es6RewriteBlockScopedFunctionDeclaration) VisitBlockScopedFunctionDeclaration(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	// Prepare a spot for the function.
	oldNameNode := n.GetFirstChild()
	fnNameNode := oldNameNode.CloneNode()
	let := rhino.IR_Declaration(fnNameNode, rhino.Token_LET).Srcref(n)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_LET_DECLARATIONS, e.compiler)

	// Prepare the function.
	oldNameNode.SetString("")
	e.compiler.ReportChangeToEnclosingScope(oldNameNode)

	// Move the function to the front of the parent.
	n.Detach()
	parent.AddChildToFront(let)
	e.compiler.ReportChangeToEnclosingScope(let)
	fnNameNode.AddChildToFront(n)
}

// Es6RewriteRestParameters converts REST parameters.
type Es6RewriteRestParameters struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler         AbstractCompiler
	astFactory       *AstFactory
	namespace        rhino.StaticScope
	getRestArguments RuntimeJsLibManager_JsLibField
}

// Java: Es6RewriteRestParameters.TRANSPILED_FEATURES
var Es6RewriteRestParameters_TRANSPILED_FEATURES = parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_REST_PARAMETERS)

var es6RewriteRestParametersVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_ITER_REST)
	return s
}()

func (*Es6RewriteRestParameters) noVisitLeaves() {}
func (*Es6RewriteRestParameters) visitTokens() *ast.TokenSet {
	return &es6RewriteRestParametersVisitTokens
}

// Java: Es6RewriteRestParameters(AbstractCompiler)
func NewEs6RewriteRestParameters(compiler AbstractCompiler) *Es6RewriteRestParameters {
	e := &Es6RewriteRestParameters{}
	e.compiler = compiler
	e.astFactory = compiler.CreateAstFactory()
	e.namespace = compiler.GetTranspilationNamespace()
	e.getRestArguments =
		compiler.GetRuntimeJsLibManager().GetJsLibField("$jscomp.getRestArguments")
	return e
}

// Java: Es6RewriteRestParameters.process(Node, Node)
func (e *Es6RewriteRestParameters) Process(externs rhino.Node, root rhino.Node) {
	TranspilationPasses_ProcessTranspile(e.compiler, root, Es6RewriteRestParameters_TRANSPILED_FEATURES, e)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, Es6RewriteRestParameters_TRANSPILED_FEATURES)
}

// Java: Es6RewriteRestParameters.visit(NodeTraversal, Node, Node)
func (e *Es6RewriteRestParameters) Visit(traversal *NodeTraversal, current rhino.Node, parent rhino.Node) {
	if current.GetToken() == rhino.Token_ITER_REST {
		e.VisitRestParam(traversal, current, parent)
	}
}

// VisitRestParam processes a rest parameter
//
// Java: Es6RewriteRestParameters.visitRestParam(NodeTraversal, Node, Node)
func (e *Es6RewriteRestParameters) VisitRestParam(t *NodeTraversal, restParam rhino.Node, paramList rhino.Node) {
	functionBody := paramList.GetNext()
	restIndex := paramList.GetIndexOfChild(restParam)
	nameNode := restParam.GetOnlyChild()
	paramName := nameNode.GetString()

	// Remove the existing param from the list, as it will be replaced with a declaration with the
	// same name.
	restParam.Detach()

	if !functionBody.HasChildren() {
		// If function has no body, we are done!
		t.ReportCodeChange()
		return
	}

	// Now that the restParam is deleted, create a let declaration by making a new NAME node of the
	// same name `paramName`
	let :=
		e.astFactory.
			CreateSingleLetNameDeclarationWithValue(
				paramName, // creates a new NAME node with name `paramName`
				e.astFactory.CreateCall(
					e.astFactory.CreateGetPropWithUnknownType(
						e.astFactory.CreateQNameOfField(e.namespace, e.getRestArguments), "apply"),
					AstFactory_Type_(nameNode),
					e.astFactory.CreateNumber(float64(restIndex)),
					e.astFactory.CreateArgumentsReference())).
			SrcrefTreeIfMissing(functionBody)
	insertBeforePoint :=
		NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(functionBody)
	if !insertBeforePoint.IsNil() {
		let.InsertBefore(insertBeforePoint)
	} else {
		// functionBody only contains hoisted function declarations
		functionBody.AddChildToBack(let)
	}
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_LET_DECLARATIONS, e.compiler)
	t.ReportCodeChange()
}

// SubstituteEs6Syntax is an optimization that does peephole optimizations
// of ES6 code.
type SubstituteEs6Syntax struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler                       AbstractCompiler
	objectLiteralShorthandWasAdded bool
}

var substituteEs6SyntaxVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_FUNCTION)
	s.Add(rhino.Token_STRING_KEY)
	s.Add(rhino.Token_SCRIPT)
	return s
}()

func (*SubstituteEs6Syntax) noVisitLeaves()             {}
func (*SubstituteEs6Syntax) visitTokens() *ast.TokenSet { return &substituteEs6SyntaxVisitTokens }

// Java: SubstituteEs6Syntax(AbstractCompiler)
func NewSubstituteEs6Syntax(compiler AbstractCompiler) *SubstituteEs6Syntax {
	s := &SubstituteEs6Syntax{}
	s.compiler = compiler
	return s
}

// Java: SubstituteEs6Syntax.process(Node, Node)
func (s *SubstituteEs6Syntax) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(s.compiler, root, s)
}

// Java: SubstituteEs6Syntax.visit(NodeTraversal, Node, Node)
func (s *SubstituteEs6Syntax) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_FUNCTION:
		if n.IsArrowFunction() {
			s.MaybeSimplifyArrowFunctionBody(n, n.GetLastChild())
		}
	case rhino.Token_STRING_KEY:
		if n.GetFirstChild().IsName() && n.GetFirstChild().GetStringID() == n.GetStringID() {
			n.SetShorthandProperty(true)
			s.objectLiteralShorthandWasAdded = true
		}
	case rhino.Token_SCRIPT:
		if s.objectLiteralShorthandWasAdded {
			NodeUtil_AddFeatureToScript(n, parser.FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES, s.compiler)
			s.objectLiteralShorthandWasAdded = false // false for the next script
		}
	default:
	}
}

// MaybeSimplifyArrowFunctionBody: if possible, replace functions of the
// form ()=>{ return x; } with ()=>x
//
// Java: SubstituteEs6Syntax.maybeSimplifyArrowFunctionBody(Node, Node)
func (s *SubstituteEs6Syntax) MaybeSimplifyArrowFunctionBody(arrowFunction rhino.Node, body rhino.Node) {
	jutil.CheckArgument(arrowFunction.IsArrowFunction())
	if !body.IsBlock() || !body.HasOneChild() || !body.GetFirstChild().IsReturn() {
		return
	}
	returnValue := body.GetFirstChild().RemoveFirstChild()
	var replacement rhino.Node
	if !returnValue.IsNil() {
		replacement = returnValue
	} else {
		replacement = rhino.IR_Name("undefined")
	}
	body.ReplaceWith(replacement)
	s.compiler.ReportChangeToEnclosingScope(replacement)
}

// Go-only: this pass uses the typed AST view (internal/ast) to read the
// object literal it rewrites: its members (ast.ObjectMember), which of them
// are spreads, and the expression of each spread (spread.Expr()). The
// Object.assign call and the object literals it receives come from the
// AstFactory (which sets their colors), so they stay rhino.Nodes.

// RewriteObjectSpread converts object spread to valid ES2017 code.
//
// Currently this class converts Object spread properties as documented in
// tc39. https://github.com/tc39/proposal-object-rest-spread. For example:
//
// `const bar = {a: 1, ...foo};`
//
// Note that object rest is handled by Es6RewriteDestructuring.
type RewriteObjectSpread struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler AbstractCompiler

	astFactory *AstFactory
	namespace  rhino.StaticScope
}

// Java: RewriteObjectSpread.transpiledFeatures
var RewriteObjectSpread_transpiledFeatures = parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD)

var rewriteObjectSpreadVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_OBJECTLIT)
	return s
}()

// Java: RewriteObjectSpread(AbstractCompiler)
func NewRewriteObjectSpread(compiler AbstractCompiler) *RewriteObjectSpread {
	return &RewriteObjectSpread{
		compiler:   compiler,
		astFactory: compiler.CreateAstFactory(),
		namespace:  compiler.GetTranspilationNamespace(),
	}
}

func (*RewriteObjectSpread) visitTokens() *ast.TokenSet {
	return &rewriteObjectSpreadVisitTokens
}

// Java: RewriteObjectSpread.process(Node, Node)
func (r *RewriteObjectSpread) Process(externs rhino.Node, root rhino.Node) {
	TranspilationPasses_ProcessTranspile(r.compiler, root, RewriteObjectSpread_transpiledFeatures, r)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(r.compiler, root, RewriteObjectSpread_transpiledFeatures)
}

// Java: RewriteObjectSpread.visit(NodeTraversal, Node, Node)
func (r *RewriteObjectSpread) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if obj, ok := ast.AsObjectLit(n); ok {
		r.VisitObject(obj)
	}
}

// Java: RewriteObjectSpread.visitObject(Node)
func (r *RewriteObjectSpread) VisitObject(obj ast.ObjectLit) {
	for member := range obj.Members() {
		if _, ok := member.AsObjectSpread(); ok {
			r.VisitObjectWithSpread(obj)
			return
		}
	}
}

// Convert '{first: b, c, ...spread, d: e, last}' to:
//
// Object.assign({}, {first:b, c}, spread, {d:e, last});
//
// Java: RewriteObjectSpread.visitObjectWithSpread(Node)
func (r *RewriteObjectSpread) VisitObjectWithSpread(obj ast.ObjectLit) {
	// Add an empty target object literal so changes made by Object.assign will not affect any other
	// variables.
	result :=
		r.astFactory.CreateObjectDotAssignCall(r.namespace, AstFactory_Type_(obj.Node()), r.astFactory.CreateObjectLit())

	// An indicator whether the current last thing in the param list is an object literal to which
	// properties may be added.  Initialized to null since nothing should be added to the empty
	// object literal in first position of the param list.
	var trailingObjectLiteral rhino.Node = rhino.NilNode

	// ReadAhead: the current member may move.
	for member := range ast.ReadAhead(obj.FirstMember()) {
		if spread, ok := member.AsObjectSpread(); ok {
			// Add the object directly to the param list.
			spreaded := spread.Expr().Detach()
			result.AddChildToBack(spreaded.Node())

			// Properties should not be added to the trailing object.
			trailingObjectLiteral = rhino.NilNode
		} else {
			if trailingObjectLiteral.IsNil() {
				// Add a new object to which properties may be added.
				trailingObjectLiteral = r.astFactory.CreateObjectLit()
				result.AddChildToBack(trailingObjectLiteral)
			}
			// Add the property to the object literal.
			trailingObjectLiteral.AddChildToBack(member.Detach().Node())
		}
	}

	result.SrcrefTreeIfMissing(obj.Node())
	obj.Node().ReplaceWith(result)
	r.compiler.ReportChangeToEnclosingScope(result)
}

type tryCatchChunkFunctions struct {
	chunk   *JSChunk
	fnNodes []rhino.Node
}

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
	compiler  AbstractCompiler
	functions []tryCatchChunkFunctions
	classes   []rhino.Node
}

var rewriteGlobalDeclarationsVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_FUNCTION)
	s.Add(rhino.Token_CLASS)
	s.Add(rhino.Token_CONST)
	s.Add(rhino.Token_LET)
	return s
}()

func (*RewriteGlobalDeclarationsForTryCatchWrapping) noVisitLeaves() {}
func (*RewriteGlobalDeclarationsForTryCatchWrapping) visitTokens() *ast.TokenSet {
	return &rewriteGlobalDeclarationsVisitTokens
}

// Java: RewriteGlobalDeclarationsForTryCatchWrapping(AbstractCompiler)
func NewRewriteGlobalDeclarationsForTryCatchWrapping(compiler AbstractCompiler) *RewriteGlobalDeclarationsForTryCatchWrapping {
	return &RewriteGlobalDeclarationsForTryCatchWrapping{
		compiler: compiler,
	}
}

func (r *RewriteGlobalDeclarationsForTryCatchWrapping) addFunction(chunk *JSChunk, n rhino.Node) {
	if m := len(r.functions); m > 0 && r.functions[m-1].chunk == chunk {
		r.functions[m-1].fnNodes = append(r.functions[m-1].fnNodes, n)
		return
	}
	for i := range r.functions {
		if r.functions[i].chunk == chunk {
			r.functions[i].fnNodes = append(r.functions[i].fnNodes, n)
			return
		}
	}
	r.functions = append(r.functions, tryCatchChunkFunctions{
		chunk:   chunk,
		fnNodes: []rhino.Node{n},
	})
}

// Java: RewriteGlobalDeclarationsForTryCatchWrapping.process(Node, Node)
func (r *RewriteGlobalDeclarationsForTryCatchWrapping) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(r.compiler, root, r)
	for i := range r.functions {
		chunk := r.functions[i].chunk
		fnNodes := r.functions[i].fnNodes
		addingRoot := r.compiler.GetNodeForCodeInsertion(chunk)
		if len(fnNodes) != 0 {
			// Lists.reverse(entry.getValue())
			for j := len(fnNodes) - 1; j >= 0; j-- {
				n := fnNodes[j]
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

	for _, n := range r.classes {
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
		r.addFunction(t.GetChunk(), n)
	}

	if NodeUtil_IsClassDeclaration(n) {
		r.classes = append(r.classes, n)
	}

	if n.IsConst() || n.IsLet() {
		n.SetToken(rhino.Token_VAR)
		r.compiler.ReportChangeToEnclosingScope(n)
	}
}
