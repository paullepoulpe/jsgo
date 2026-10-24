// Port of com/google/javascript/jscomp/Es6RewriteSpreadExpressions.java @ closure-compiler v20261015

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
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RewriteSpreadExpressions converts SPREAD expressions in array literals,
// new calls, and function calls.
type Es6RewriteSpreadExpressions struct {
	NodeTraversal_AbstractPostOrderCallbackBase

	compiler   AbstractCompiler
	astFactory *AstFactory
	namespace  rhino.StaticScope
}

// Java: Es6RewriteSpreadExpressions.TRANSPILED_FEATURES
var Es6RewriteSpreadExpressions_TRANSPILED_FEATURES = parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_SPREAD_EXPRESSIONS)

// Java: Es6RewriteSpreadExpressions.FRESH_SPREAD_VAR
const Es6RewriteSpreadExpressions_FRESH_SPREAD_VAR = "$jscomp$spread$args"

// Java: Es6RewriteSpreadExpressions.arrayType
var Es6RewriteSpreadExpressions_arrayType = AstFactory_Type_OfColorId(colors.StandardColors_ARRAY_ID)

// Java: Es6RewriteSpreadExpressions.concatFnType
var Es6RewriteSpreadExpressions_concatFnType = AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT)

// Java: Es6RewriteSpreadExpressions(AbstractCompiler)
func NewEs6RewriteSpreadExpressions(compiler AbstractCompiler) *Es6RewriteSpreadExpressions {
	e := &Es6RewriteSpreadExpressions{}
	e.compiler = compiler
	e.astFactory = compiler.CreateAstFactory()
	e.namespace = compiler.GetTranspilationNamespace()
	return e
}

// Java: Es6RewriteSpreadExpressions.process(Node, Node)
func (e *Es6RewriteSpreadExpressions) Process(externs rhino.Node, root rhino.Node) {
	TranspilationPasses_ProcessTranspile(e.compiler, root, Es6RewriteSpreadExpressions_TRANSPILED_FEATURES, e)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, Es6RewriteSpreadExpressions_TRANSPILED_FEATURES)
}

// Java: Es6RewriteSpreadExpressions.visit(NodeTraversal, Node, Node)
func (e *Es6RewriteSpreadExpressions) Visit(traversal *NodeTraversal, current rhino.Node, parent rhino.Node) {
	token := current.GetToken()
	if token == rhino.Token_ARRAYLIT || token == rhino.Token_NEW || token == rhino.Token_CALL {
		for child := current.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
			if child.IsSpread() {
				e.VisitArrayLitOrCallWithSpread(traversal, current)
				break
			}
		}
	}
}

// VisitArrayLitOrCallWithSpread processes array literals or calls to
// eliminate spreads.
//
// Examples:
//
//   - [1, 2, ...x, 4, 5] => [1, 2].concat($jscomp.arrayFromIterable(x), [4, 5])
//   - f(1, ...arr) => f.apply(null, [1].concat($jscomp.arrayFromIterable(arr)))
//   - new F(...args) => new Function.prototype.bind.apply(F,
//     [null].concat($jscomp.arrayFromIterable(args)))
//
// Java: Es6RewriteSpreadExpressions.visitArrayLitOrCallWithSpread(NodeTraversal, Node)
func (e *Es6RewriteSpreadExpressions) VisitArrayLitOrCallWithSpread(t *NodeTraversal, spreadParent rhino.Node) {
	if spreadParent.IsArrayLit() {
		e.VisitArrayLitContainingSpread(spreadParent)
	} else if spreadParent.IsCall() {
		e.VisitCallContainingSpread(t, spreadParent)
	} else {
		jutil.CheckArgumentMsg1(spreadParent.IsNew(), "%s", spreadParent)
		e.VisitNewWithSpread(spreadParent)
	}
}

// ExtractSpreadGroups extracts child nodes from an ARRAYLIT, CALL or NEW node
// that may contain spread operators into a list of nodes that may be
// concatenated with Array.concat() to get an array.
//
// Example: [a, b, ...x, c, ...arguments] returns a list containing [ [a, b],
// $jscomp.arrayFromIterable(x), [c], $jscomp.arrayFromIterable(arguments) ]
//
// IMPORTANT: CALL and NEW nodes must have the first, callee, child removed
// already.
//
// Note that all elements of the returned list will be one of:
//
//   - array literal
//   - $jscomp.arrayFromIterable(spreadExpression)
//
// TODO(bradfordcsmith): When this pass moves after type checking, we can use
// type information to avoid unnecessary calls to $jscomp.arrayFromIterable().
//
// TODO(nickreid): Stop mutating `spreadParent`.
//
// Java: Es6RewriteSpreadExpressions.extractSpreadGroups(Node)
func (e *Es6RewriteSpreadExpressions) ExtractSpreadGroups(spreadParent rhino.Node) []rhino.Node {
	jutil.CheckArgument(spreadParent.IsCall() || spreadParent.IsArrayLit() || spreadParent.IsNew())

	groups := []rhino.Node{}
	var currGroup rhino.Node = rhino.Node{}

	for currElement := spreadParent.RemoveFirstChild(); !currElement.IsNil(); currElement = spreadParent.RemoveFirstChild() {
		if currElement.IsSpread() {
			spreadExpression := currElement.RemoveFirstChild()
			if spreadExpression.IsArrayLit() {
				// We can expand an array literal spread in place.
				if currGroup.IsNil() {
					// [...[spread, contents], a, b]
					// we can use this array lit itself as a group and append following elements to it
					currGroup = spreadExpression
				} else {
					// [ a, b, ...[spread, contents], c]
					// Just add contents of this array lit to the group we were already collecting.
					currGroup.AddChildrenToBack(spreadExpression.RemoveChildren())
				}
			} else {
				// We need to treat the spread expression as a separate group
				if !currGroup.IsNil() {
					// finish off and add the group we were collecting before
					groups = append(groups, currGroup)
					currGroup = rhino.Node{}
				}

				groups = append(groups,
					e.astFactory.CreateJscompArrayFromIterableCall(spreadExpression, e.namespace))
			}
		} else {
			if currGroup.IsNil() {
				currGroup = e.astFactory.CreateArraylit()
			}
			currGroup.AddChildToBack(currElement)
		}
	}

	if !currGroup.IsNil() {
		groups = append(groups, currGroup)
	}
	return groups
}

// VisitArrayLitContainingSpread processes array literals containing spreads.
//
// Example:
//
//	[1, 2, ...x, 4, 5] => [1, 2].concat($jscomp.arrayFromIterable(x), [4, 5])
//
// Java: Es6RewriteSpreadExpressions.visitArrayLitContainingSpread(Node)
func (e *Es6RewriteSpreadExpressions) VisitArrayLitContainingSpread(spreadParent rhino.Node) {
	jutil.CheckArgument(spreadParent.IsArrayLit())

	groups := e.ExtractSpreadGroups(spreadParent)

	var baseArrayLit rhino.Node
	if groups[0].IsArrayLit() {
		// g0.concat(g1, g2, ..., gn)
		baseArrayLit = groups[0]
		groups = groups[1:]
	} else {
		// [].concat(g0, g1, g2, ..., gn)
		baseArrayLit = e.astFactory.CreateArraylit()
	}

	var joinedGroups rhino.Node
	if len(groups) == 0 {
		joinedGroups = baseArrayLit
	} else {
		concat := e.astFactory.CreateGetProp(baseArrayLit, "concat", Es6RewriteSpreadExpressions_concatFnType)
		joinedGroups = e.astFactory.CreateCall(concat, AstFactory_Type_(spreadParent), groups...)
	}

	joinedGroups.SrcrefTreeIfMissing(spreadParent)

	spreadParent.ReplaceWith(joinedGroups)
	e.compiler.ReportChangeToEnclosingScope(joinedGroups)
}

// VisitCallContainingSpread processes calls containing spreads.
//
// Examples:
//
//	f(...arr) => f.apply(null, $jscomp.arrayFromIterable(arr))
//	f(a, ...arr) => f.apply(null, [a].concat($jscomp.arrayFromIterable(arr)))
//	f(...arr, b) => f.apply(null, [].concat($jscomp.arrayFromIterable(arr), [b]))
//
// Java: Es6RewriteSpreadExpressions.visitCallContainingSpread(NodeTraversal, Node)
func (e *Es6RewriteSpreadExpressions) VisitCallContainingSpread(t *NodeTraversal, spreadParent rhino.Node) {
	jutil.CheckArgument(spreadParent.IsCall())

	callee := spreadParent.GetFirstChild()
	// ES6 classes must all be transpiled away before this pass runs.
	jutil.CheckStateMsg(!callee.IsSuper(), "Cannot spread into super calls")
	// Check if the callee has side effects before removing it from the AST (since some NodeUtil
	// methods assume the node they are passed has a non-null parent).
	calleeMayHaveSideEffects := e.compiler.GetAstAnalyzer().MayHaveSideEffects(callee)
	// Must remove callee before extracting argument groups.
	callee.Detach()

	for callee.IsCast() {
		// Drop any CAST nodes. They're not needed anymore since this pass runs at the end of
		// the checks phase, and they complicate detecting GETPROP/GETELEM callees.
		callee = callee.RemoveFirstChild()
	}
	var joinedGroups rhino.Node
	if spreadParent.HasOneChild() && e.IsSpreadOfArguments(spreadParent.GetOnlyChild()) {
		// Check for special case of `foo(...arguments)` and pass `arguments` directly to
		// `foo.apply(null, arguments)`. We want to avoid calling $jscomp.arrayFromIterable(arguments)
		// for this case, because it can have side effects, which prevents code removal.
		joinedGroups = spreadParent.RemoveFirstChild().RemoveFirstChild()
	} else {
		groups := e.ExtractSpreadGroups(spreadParent)
		jutil.CheckState(len(groups) != 0)

		if len(groups) == 1 {
			// A single group can just be passed to `apply()` as-is
			// It could be `arguments`, an array literal, or $jscomp.arrayFromIterable(someExpression).
			joinedGroups = groups[0]
			groups = groups[1:]
		} else {
			// If the first group is an array literal, we can just use that for concatenation,
			// otherwise use an empty array literal.
			//
			// TODO(bradfordcsmith): Now that this pass runs after type checking, it would be nice
			// to skip creating an array literal when when the type of the first element says it is
			// an Array.
			var baseArrayLit rhino.Node
			if groups[0].IsArrayLit() {
				baseArrayLit = groups[0]
				groups = groups[1:]
			} else {
				baseArrayLit = e.astFactory.CreateArraylit()
			}
			concat := e.astFactory.CreateGetProp(baseArrayLit, "concat", Es6RewriteSpreadExpressions_concatFnType)
			joinedGroups = e.astFactory.CreateCall(concat, Es6RewriteSpreadExpressions_arrayType, groups...)
		}
	}
	isFreeCall := spreadParent.GetBooleanProp(rhino.Node_FREE_CALL)

	var callToApply rhino.Node
	if calleeMayHaveSideEffects && (callee.IsGetProp() || callee.IsGetElem()) && !isFreeCall {
		// foo().method(...[a, b, c]) or foo()['method'](...[a, b, c])
		//   must convert to
		// var freshVar;
		// (freshVar = foo()).method.apply(freshVar, [a, b, c])
		freshVar :=
			e.astFactory.CreateName(
				Es6RewriteSpreadExpressions_FRESH_SPREAD_VAR+e.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput()),
				// Type of `foo()`.
				AstFactory_Type_(callee.GetFirstChild()))
		freshVarDeclaration := rhino.IR_Var(freshVar.CloneTree())

		statementContainingSpread := NodeUtil_GetEnclosingStatement(spreadParent)
		freshVarDeclaration.SrcrefTreeIfMissing(statementContainingSpread)

		freshVarDeclaration.InsertBefore(statementContainingSpread)
		callee.AddChildToFront(
			e.astFactory.CreateAssign(freshVar.CloneTree(), callee.RemoveFirstChild()))

		callToApply =
			e.astFactory.CreateCallWithUnknownType(
				e.astFactory.CreateGetPropWithUnknownType(callee, "apply"),
				freshVar.CloneTree(),
				joinedGroups)
	} else {
		// foo.method(...[a, b, c]) -> foo.method.apply(foo, [a, b, c])
		// foo['method'](...[a, b, c]) -> foo['method'].apply(foo, [a, b, c])
		// or
		// foo(...[a, b, c]) -> foo.apply(null, [a, b, c])
		var context rhino.Node
		if (callee.IsGetProp() || callee.IsGetElem()) && !isFreeCall {
			context = callee.GetFirstChild().CloneTree()
		} else {
			context = e.astFactory.CreateNull()
		}
		callToApply =
			e.astFactory.CreateCall(
				e.astFactory.CreateGetPropWithUnknownType(callee, "apply"),
				AstFactory_Type_(spreadParent),
				context,
				joinedGroups)
	}

	callToApply.SetColor(spreadParent.GetColor())
	callToApply.SrcrefTreeIfMissing(spreadParent)
	spreadParent.ReplaceWith(callToApply)
	e.compiler.ReportChangeToEnclosingScope(callToApply)
}

// Java: Es6RewriteSpreadExpressions.isSpreadOfArguments(Node)
func (e *Es6RewriteSpreadExpressions) IsSpreadOfArguments(n rhino.Node) bool {
	return n.IsSpread() && n.GetOnlyChild().MatchesName("arguments")
}

// VisitNewWithSpread processes new calls containing spreads.
//
// Example:
//
//	new F(...args) =>
//	    new Function.prototype.bind.apply(F, [].concat($jscomp.arrayFromIterable(args)))
//
// Java: Es6RewriteSpreadExpressions.visitNewWithSpread(Node)
func (e *Es6RewriteSpreadExpressions) VisitNewWithSpread(spreadParent rhino.Node) {
	jutil.CheckArgument(spreadParent.IsNew())

	// Must remove callee before extracting argument groups.
	callee := spreadParent.RemoveFirstChild()
	groups := e.ExtractSpreadGroups(spreadParent)

	// We need to generate
	// `new (Function.prototype.bind.apply(callee, [null].concat(other, args))();`.
	// `null` stands in for the 'this' arg to the contructor.
	var baseArrayLit rhino.Node
	if groups[0].IsArrayLit() {
		baseArrayLit = groups[0]
		groups = groups[1:]
	} else {
		baseArrayLit = e.astFactory.CreateArraylit()
	}
	baseArrayLit.AddChildToFront(e.astFactory.CreateNull())
	var joinedGroups rhino.Node
	if len(groups) == 0 {
		joinedGroups = baseArrayLit
	} else {
		joinedGroups = e.astFactory.CreateCall(
			e.astFactory.CreateGetProp(baseArrayLit, "concat", Es6RewriteSpreadExpressions_concatFnType),
			Es6RewriteSpreadExpressions_arrayType,
			groups...)
	}

	if parser.FeatureSet_ES3.Contains(e.compiler.GetOptions().GetOutputFeatureSet()) {
		// TODO(tbreisacher): Support this in ES3 too by not relying on Function.bind.
		TranspilationUtil_CannotConvert(
			e.compiler,
			spreadParent,
			"\"...\" passed to a constructor (consider using --language_out=ES5)")
	}

	// Function.prototype.bind =>
	//      function(this:function(new:[spreadParent], ...?), ...?):function(new:[spreadParent])
	// Function.prototype.bind.apply =>
	//      function(function(new:[spreadParent], ...?), !Array<?>):function(new:[spreadParent])
	bindApply :=
		e.astFactory.CreateGetPropWithUnknownType(
			e.astFactory.CreateGetPropWithUnknownType(
				e.astFactory.CreatePrototypeAccess(e.astFactory.CreateNameInScope(e.namespace, "Function")),
				"bind"),
			"apply")
	result :=
		rhino.IR_NewNode(
			e.astFactory.CreateCallWithUnknownType(
				bindApply, callee, joinedGroups /* function(new:[spreadParent]) */)).
			SetColor(spreadParent.GetColor())

	result.SrcrefTreeIfMissing(spreadParent)
	spreadParent.ReplaceWith(result)
	e.compiler.ReportChangeToEnclosingScope(result)
}
