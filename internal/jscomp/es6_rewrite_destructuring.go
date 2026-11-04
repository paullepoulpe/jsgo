// Port of com/google/javascript/jscomp/Es6RewriteDestructuring.java @ closure-compiler v20261028

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

package jscomp

import (
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast) where it reads
// pattern shapes: the properties of an object pattern (key.Value(),
// computed.Key(), rest.Target()), the elements of an array pattern, default
// values (dv.Target(), dv.Value()), the parameters of a function, and the
// declaration, assignment, loop or catch that holds a pattern. The rewrites
// take the kind they work on (ast.ObjectPattern, ast.ArrayPattern,
// ast.DestructuringLhs, ast.EnhancedFor...) instead of a node they check
// again, so the port no longer follows the Java line by line. The per-node
// dispatch of shouldTraverse and visit stays a token switch, and the new code
// is built with the rhino and AstFactory builders, in the Java's order.

// Es6RewriteDestructuring rewrites destructuring patterns and default
// parameters to valid ES3 code or to a different form of destructuring.
type Es6RewriteDestructuring struct {
	compiler    AbstractCompiler
	rewriteMode Es6RewriteDestructuring_ObjectDestructuringRewriteMode
	astFactory  *AstFactory
	namespace   rhino.StaticScope

	featuresToTriggerRunningPass parser.FeatureSet
	featuresToMarkAsRemoved      parser.FeatureSet

	patternNestingStack []patternNestingLevel
	stackBuf            [4]patternNestingLevel

	destructuringVarCounter int
}

// Es6RewriteDestructuring_ObjectDestructuringRewriteMode is
// Es6RewriteDestructuring.ObjectDestructuringRewriteMode.
type Es6RewriteDestructuring_ObjectDestructuringRewriteMode int32

const (
	// Rewrite all object destructuring patterns. This is the default mode used
	// if no ObjectDestructuringRewriteMode is provided to the Builder.
	//
	// Used to transpile ES2018 -> ES5
	Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_ALL_OBJECT_PATTERNS Es6RewriteDestructuring_ObjectDestructuringRewriteMode = iota

	// Rewrite only destructuring patterns that contain object pattern rest
	// properties (whether the rest is on the top level or nested within a
	// property).
	//
	// Used to transpile ES2018 -> ES2017
	Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST
)

var es6RewriteDestructuring_ObjectDestructuringRewriteMode_names = [...]string{
	"REWRITE_ALL_OBJECT_PATTERNS", "REWRITE_OBJECT_REST"}

// String is Enum.name().
//
// Java: ObjectDestructuringRewriteMode.name()
func (e Es6RewriteDestructuring_ObjectDestructuringRewriteMode) String() string {
	return es6RewriteDestructuring_ObjectDestructuringRewriteMode_names[e]
}

// Java: Es6RewriteDestructuring.DESTRUCTURING_TEMP_VAR
const Es6RewriteDestructuring_DESTRUCTURING_TEMP_VAR = "$jscomp$destructuring$var"

// NewEs6RewriteDestructuring constructs an Es6RewriteDestructuring pass.
func NewEs6RewriteDestructuring(
	compiler AbstractCompiler,
	rewriteMode Es6RewriteDestructuring_ObjectDestructuringRewriteMode,
) *Es6RewriteDestructuring {
	e := &Es6RewriteDestructuring{}
	e.patternNestingStack = e.stackBuf[:0]
	e.destructuringVarCounter = 0
	e.compiler = compiler
	e.rewriteMode = rewriteMode
	e.astFactory = e.compiler.CreateAstFactory()
	e.namespace = e.compiler.GetTranspilationNamespace()

	switch e.rewriteMode {
	case Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_ALL_OBJECT_PATTERNS:
		e.featuresToTriggerRunningPass =
			parser.FeatureSet_BARE_MINIMUM.WithFeatures(
				parser.FeatureSet_Feature_DEFAULT_PARAMETERS,
				parser.FeatureSet_Feature_ARRAY_DESTRUCTURING,
				parser.FeatureSet_Feature_ARRAY_PATTERN_REST,
				parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)

		// If OBJECT_PATTERN_REST were to be present in featuresToTriggerRunningPass and not the
		// input language featureSet (such as ES6=>ES5) the pass would be skipped.
		e.featuresToMarkAsRemoved =
			e.featuresToTriggerRunningPass.With(parser.FeatureSet_Feature_OBJECT_PATTERN_REST)
	case Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST:
		// TODO(bradfordcsmith): We shouldn't really need to remove default parameters for this
		// case.
		e.featuresToTriggerRunningPass =
			parser.FeatureSet_BARE_MINIMUM.With(parser.FeatureSet_Feature_OBJECT_PATTERN_REST)
		e.featuresToMarkAsRemoved = e.featuresToTriggerRunningPass
	default:
		panic(jutil.NewAssertionError(
			"Es6RewriteDestructuring cannot handle ObjectDestructuringRewriteMode " +
				e.rewriteMode.String()))
	}
	return e
}

type patternNestingLevel struct {
	pattern             ast.Pattern
	hasNestedObjectRest bool
}

func (e *Es6RewriteDestructuring) pushPatternNestingLevel(pattern ast.Pattern, hasNestedRest bool) {
	e.patternNestingStack = append(e.patternNestingStack, patternNestingLevel{
		pattern:             pattern,
		hasNestedObjectRest: hasNestedRest,
	})
}

func (e *Es6RewriteDestructuring) peekPatternNestingLevel() *patternNestingLevel {
	if len(e.patternNestingStack) == 0 {
		return nil
	}
	return &e.patternNestingStack[len(e.patternNestingStack)-1]
}

func (e *Es6RewriteDestructuring) popPatternNestingLevel() {
	last := len(e.patternNestingStack) - 1
	e.patternNestingStack[last] = patternNestingLevel{}
	e.patternNestingStack = e.patternNestingStack[:last]
}

// Java: Es6RewriteDestructuring.process(Node, Node)
func (e *Es6RewriteDestructuring) Process(externs rhino.Node, root rhino.Node) {
	if len(e.patternNestingStack) != 0 {
		panic(jutil.NewIllegalStateException())
	}
	NodeTraversal_Traverse(e.compiler, root, e)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, e.featuresToMarkAsRemoved)
	if len(e.patternNestingStack) != 0 {
		panic(jutil.NewIllegalStateException())
	}
}

// Java: Es6RewriteDestructuring.shouldTraverse(NodeTraversal, Node, Node)
func (e *Es6RewriteDestructuring) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		scriptFeatures := NodeUtil_GetFeatureSetOfScript(n)
		// This will ensure that we run this only when features exist in the script
		return scriptFeatures.ContainsAtLeastOneOf(e.featuresToTriggerRunningPass)
	}
	switch n.GetToken() {
	case rhino.Token_PARAM_LIST:
		e.pullDestructuringOutOfParams(ast.MustFunction(parent))
	case rhino.Token_ARRAY_PATTERN, rhino.Token_OBJECT_PATTERN:
		pattern := ast.MustPattern(n)
		hasRest := false
		if objectPattern, ok := pattern.AsObjectPattern(); ok {
			_, hasRest = objectPatternRest(objectPattern)
		}
		if hasRest {
			for i := len(e.patternNestingStack) - 1; i >= 0; i-- {
				if e.patternNestingStack[i].hasNestedObjectRest {
					break
				}
				e.patternNestingStack[i].hasNestedObjectRest = true
			}
		}
		e.pushPatternNestingLevel(pattern, hasRest)
	default:
	}
	return true
}

// objectPatternRest returns the ...rest that ends objectPattern; ok is false
// if it has none.
//
// Go-only: Java tests objectPattern.getLastChild().isRest() where it needs
// this.
func objectPatternRest(objectPattern ast.ObjectPattern) (ast.ObjectRest, bool) {
	last, ok := objectPattern.LastProp()
	if !ok {
		return ast.ObjectRest{}, false
	}
	return last.AsObjectRest()
}

var es6RewriteDestructuringVisitTokens = func() ast.TokenSet {
	var s ast.TokenSet
	s.Add(rhino.Token_ARRAY_PATTERN)
	s.Add(rhino.Token_OBJECT_PATTERN)
	return s
}()

func (*Es6RewriteDestructuring) visitTokens() *ast.TokenSet {
	return &es6RewriteDestructuringVisitTokens
}

var _ nodeTraversal_visitTokens = (*Es6RewriteDestructuring)(nil)

// Java: Es6RewriteDestructuring.visit(NodeTraversal, Node, Node)
func (e *Es6RewriteDestructuring) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_ARRAY_PATTERN, rhino.Token_OBJECT_PATTERN:
		pattern := ast.MustPattern(n)
		e.visitPattern(t, pattern)
		if top := e.peekPatternNestingLevel(); top != nil && top.pattern.Is(n) {
			e.popPatternNestingLevel()
		}
	default:
	}
}

// pullDestructuringOutOfParams pulls all default and destructuring
// parameters out of function parameters.
//
// TODO(bradfordcsmith): Ideally if we're only removing OBJECT_REST, we should
// only do this when the parameter list contains a usage of OBJECT_REST.
//
// Go-only: it takes the function alone; Java also passes its parameter list.
//
// Java: Es6RewriteDestructuring.pullDestructuringOutOfParams(Node, Node)
func (e *Es6RewriteDestructuring) pullDestructuringOutOfParams(function ast.Function) {
	var insertSpot rhino.Node = rhino.NilNode
	body := function.Body().Node()
	var nextParam ast.Param
	for param, hasParam := function.Params().FirstParam(); hasParam; param = nextParam {
		nextParam, hasParam = param.Next()
		if defaultParam, ok := param.AsDefaultValue(); ok {
			nameOrPattern, defaultValue := defaultParam.Target(), defaultParam.Value() // read both before detaching either
			nameOrPattern.Detach()
			defaultValue.Detach()
			// We'll be cloning nameOrPattern below, and we don't want to clone the JSDoc info with it
			jsDoc := nameOrPattern.Node().GetJSDocInfo()
			nameOrPattern.Node().SetJSDocInfo(nil)
			var newParam rhino.Node

			// Treat name=undefined (and equivalent) as if it was just name.  There
			// is no need to generate a (name===void 0?void 0:name) statement for
			// such arguments.
			_, isName := nameOrPattern.AsName()
			isNoop := false
			if !isName {
				// Do not try to optimize unless nameOrPattern is a simple name.
			} else if name, ok := defaultValue.AsName(); ok {
				isNoop = "undefined" == name.Ident()
			} else if void, ok := defaultValue.AsUnaryOf(rhino.Token_VOID); ok {
				// Any kind of 'void literal' is fine, but 'void fun()' or anything
				// else with side effects isn't.  We're not trying to be particularly
				// smart here and treat 'void {}' for example as if it could cause side effects.
				isNoop = NodeUtil_IsImmutableValue(void.Operand().Node())
			}

			if isNoop {
				newParam = nameOrPattern.Node().CloneTree()
			} else {
				if isName {
					newParam = nameOrPattern.Node()
				} else {
					newParam = e.createTempVarNameNode(e.getTempVariableName(), AstFactory_Type_(nameOrPattern.Node()))
				}
				lhs := nameOrPattern.Node().CloneTree()
				rhs := e.defaultValueHook(newParam.CloneTree(), defaultValue.Node())
				var newStatement rhino.Node
				if isName {
					newStatement = rhino.IR_ExprResult(e.astFactory.CreateAssign(lhs, rhs))
				} else {
					newStatement = rhino.IR_VarWithValue(lhs, rhs)
				}
				newStatement.SrcrefTreeIfMissing(param.Node())
				if insertSpot.IsNil() {
					// insert new declarations only after all inner function declarations in that function
					// body to preserve normalization
					insertSpot = NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(body)
					if !insertSpot.IsNil() {
						newStatement.InsertBefore(insertSpot)
					} else {
						// functionBody only contains hoisted function declarations
						body.AddChildToBack(newStatement)
					}
				} else {
					newStatement.InsertAfter(insertSpot)
				}
				insertSpot = newStatement
				if varDecl, ok := ast.AsNameDeclaration(newStatement); ok {
					if destructuringLhs, ok := varDecl.MustOnlyDeclarator().AsDestructuringLhs(); ok {
						insertSpot = e.rewriteVarDestructuringDeclaration(destructuringLhs)
					}
				}
			}

			param.Node().ReplaceWith(newParam)
			newParam.SetJSDocInfo(jsDoc)

			e.compiler.ReportChangeToChangeScope(function.Node())
		} else if pattern, ok := param.AsPattern(); ok {
			insertSpot =
				e.replacePatternParamWithTempVar(function, insertSpot, pattern, e.getTempVariableName())
			e.compiler.ReportChangeToChangeScope(function.Node())
		} else if rest, ok := param.AsIterRest(); ok {
			if pattern, ok := rest.Target().AsPattern(); ok {
				insertSpot =
					e.replacePatternParamWithTempVar(function, insertSpot, pattern, e.getTempVariableName())
				e.compiler.ReportChangeToChangeScope(function.Node())
			}
		}
	}
}

// replacePatternParamWithTempVar replaces a destructuring pattern parameter
// with a a temporary parameter name. To do this, this function adds a new
// local assignment in the function body assigning the temporary parameter to
// the pattern and creates stub declarations for the individual names in the
// destructuring pattern.
//
// For example, given the function `function f([a, b]) {}`, this function will
// convert it to `function f(tempVar) { var a; var b; [a, b] = tempVar; }`
//
// Note: Rewrites of the destructuring pattern will happen later to rewrite
// this destructuring pattern as non-destructured code.
//
// insertSpot: the local variable declaration will be inserted after this
// statement. tempVarName: the name to use for the temporary variable. Returns
// the last stub declaration that was generated for the local variable.
//
// Java: Es6RewriteDestructuring.replacePatternParamWithTempVar(Node, Node, Node, String)
func (e *Es6RewriteDestructuring) replacePatternParamWithTempVar(
	function ast.Function, insertSpot rhino.Node, patternParam ast.Pattern, tempVarName string) rhino.Node {
	// Convert `function f([a, b]) {}` to `function f(tempVar) { var [a, b] = tempVar; }`
	pattern := patternParam.Node()
	paramType := AstFactory_Type_(pattern)
	newParam := e.createTempVarNameNode(tempVarName, paramType)
	newParam.SetJSDocInfo(pattern.GetJSDocInfo())
	pattern.ReplaceWith(newParam)
	newDecl := rhino.IR_VarWithValue(pattern, e.createTempVarNameNode(tempVarName, paramType))
	newDecl.SrcrefTreeIfMissing(pattern)

	if insertSpot.IsNil() {
		// insert new declarations only after all inner function declarations in that function body to
		// preserve normalization
		body := function.Body().Node()
		insertSpot = NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(body)
		if !insertSpot.IsNil() {
			newDecl.InsertBefore(insertSpot)
		} else {
			// functionBody only contains hoisted function declarations
			body.AddChildToBack(newDecl)
		}
	} else {
		newDecl.InsertAfter(insertSpot)
	}

	destructuringLhs := ast.MustNameDeclaration(newDecl).MustOnlyDeclarator().MustDestructuringLhs()
	return e.rewriteVarDestructuringDeclaration(destructuringLhs)
}

// rewriteVarDestructuringDeclaration rewrites a destructuring declaration to
// a destructuring assignment.
//
// For example, given the code `var [a, b] = ...;`, this function will convert
// it to `var a; var b; [a, b] = ...;`
//
// destructuringLhs is the destructuring declaration to rewrite. Returns the
// expression result containing the assignment.
//
// Java: Es6RewriteDestructuring.rewriteVarDestructuringDeclaration(Node)
func (e *Es6RewriteDestructuring) rewriteVarDestructuringDeclaration(destructuringLhs ast.DestructuringLhs) rhino.Node {
	if _, isExport := ast.AsExport(destructuringLhs.Node().GetGrandparent()); isExport {
		panic(jutil.NewIllegalStateExceptionMsg("Export destructuring declarations not expected inside a function body."))
	}
	return rewriteVarDestructuringDeclaration(e.astFactory, destructuringLhs)
}

// createTempVarNameNode creates a new name node and adds the constant name
// property because a constant variable is used (the compiler only assigns
// $jscomp$destructuring$var[num] once)
//
// Java: Es6RewriteDestructuring.createTempVarNameNode(String, AstFactory.Type)
func (e *Es6RewriteDestructuring) createTempVarNameNode(name string, type_ AstFactory_Type) rhino.Node {
	// NOTE: This does not really create a constant node as this pass runs before normalization. See
	// b/322009741.
	return e.astFactory.CreateConstantName(name, type_)
}

// getTempVariableName creates a new unique name to use for a pattern we need
// to rewrite.
//
// Java: Es6RewriteDestructuring.getTempVariableName()
func (e *Es6RewriteDestructuring) getTempVariableName() string {
	counter := e.destructuringVarCounter
	e.destructuringVarCounter++
	return Es6RewriteDestructuring_DESTRUCTURING_TEMP_VAR + strconv.Itoa(counter)
}

// Java: Es6RewriteDestructuring.visitPattern(NodeTraversal, Node)
func (e *Es6RewriteDestructuring) visitPattern(t *NodeTraversal, pattern ast.Pattern) {
	parent := pattern.Node().GetParent()
	if parent.IsNil() {
		panic(jutil.NewIllegalStateExceptionMsg("unexpected parent"))
	}

	switch parent.GetToken() {
	case rhino.Token_DESTRUCTURING_LHS:
		destructuringLhs := ast.MustDestructuringLhs(parent)
		declaration := parent.GetParent()
		declarationParent := declaration.GetParent()
		if _, ok := ast.AsFor(declarationParent); ok {
			e.visitDestructuringPatternInVanillaForInnerVars(t, destructuringLhs)
		} else if forNode, ok := ast.AsEnhancedFor(declarationParent); ok {
			e.visitDestructuringPatternInEnhancedForInnerVars(forNode, destructuringLhs)
		} else {
			rhs, ok := destructuringLhs.Init()
			if !ok {
				// only the declaration of a for-in or for-of loop has no value
				panic(jutil.NewIllegalStateException())
			}
			e.replacePattern(t, pattern, rhs.Node(), declaration, declaration)
		}
	case rhino.Token_ASSIGN:
		assign := ast.MustAssign(parent)
		if exprStmt, ok := ast.AsExprStmt(parent.GetParent()); ok {
			e.replacePattern(t, pattern, assign.Value().Node(), parent, exprStmt.Node())
		} else {
			e.wrapAssignOrDestructuringInCallToArrow(t, parent)
		}
	case rhino.Token_OBJECT_REST, rhino.Token_ITER_REST, rhino.Token_STRING_KEY, rhino.Token_ARRAY_PATTERN, rhino.Token_DEFAULT_VALUE, rhino.Token_COMPUTED_PROP:
		// Nested pattern; do nothing. We will visit it after rewriting the parent.
	case rhino.Token_FOR_OF, rhino.Token_FOR_IN, rhino.Token_FOR_AWAIT_OF:
		e.visitDestructuringPatternInEnhancedForWithOuterVars(ast.MustEnhancedFor(parent), pattern)
	case rhino.Token_CATCH:
		e.visitDestructuringPatternInCatch(t, ast.MustCatch(parent), pattern)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("unexpected parent"))
	}
}

// replacePattern transpiles a destructuring pattern in a declaration or
// assignment to ES5
//
// nodeToDetach is a statement node containing the pattern. This method will
// replace the node with one or more other statements.
//
// Go-only: rhs is the value assigned to the pattern, and parent the var, let
// or const that declares it or the assignment to it (kinds that share no
// view, so it stays a rhino.Node).
//
// Java: Es6RewriteDestructuring.replacePattern(NodeTraversal, Node, Node, Node, Node)
func (e *Es6RewriteDestructuring) replacePattern(
	t *NodeTraversal, pattern ast.Pattern, rhs rhino.Node, parent rhino.Node, nodeToDetach rhino.Node) {
	if arrayPattern, ok := pattern.AsArrayPattern(); ok {
		e.replaceArrayPattern(t, arrayPattern, rhs, parent, nodeToDetach)
	} else {
		e.replaceObjectPattern(t, pattern.MustObjectPattern(), rhs, parent, nodeToDetach)
	}
}

// replaceObjectPattern converts 'var {a: b, c: d} = rhs' to:
//
// @const var temp = rhs; var b = temp.a; var d = temp.c;
//
// Java: Es6RewriteDestructuring.replaceObjectPattern(NodeTraversal, Node, Node, Node, Node)
func (e *Es6RewriteDestructuring) replaceObjectPattern(
	t *NodeTraversal, objectPattern ast.ObjectPattern, rhs rhino.Node, parent rhino.Node, nodeToDetach rhino.Node) {
	tempVarName := e.getTempVariableName()
	tempVarType := AstFactory_Type_(objectPattern.Node())

	var restTempVarName string
	// If the last child is a rest node we will want a list of the stated properties so we can
	// exclude them from being written to the rest variable.
	var propsToDeleteForRest []rhino.Node
	_, hasRest := objectPatternRest(objectPattern)
	if hasRest {
		restTempVarName = e.getTempVariableName()
	} else if e.rewriteMode == Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST {
		// We are configured to only break object pattern rest, but this destructure has none.
		if level := e.peekPatternNestingLevel(); level == nil || !level.hasNestedObjectRest {
			// Replacement is performed after the post-order visit has reached the root pattern node, so
			// peeking last represents if there is a rest property anywhere in the entire pattern. All
			// nesting levels of lower levels have already been popped.
			e.destructuringVarCounter--
			return
		}
	}

	if _, isLabel := ast.AsLabel(nodeToDetach.GetParent()); isLabel {
		block := rhino.IR_Block().Srcref(nodeToDetach)
		nodeToDetach.ReplaceWith(block)
		block.AddChildToBack(nodeToDetach)
	}

	// create the declaration `var temp = rhs;`
	tempDecl :=
		rhino.IR_VarWithValue(e.createTempVarNameNode(tempVarName, tempVarType), rhs.Detach()).
			SrcrefTreeIfMissing(objectPattern.Node())
	// TODO(tbreisacher): Remove the "if" and add this JSDoc unconditionally.
	if parent.IsConst() {
		jsDoc := rhino.JSDocInfo_Builder_()
		jsDoc.RecordConstancy()
		tempDecl.SetJSDocInfo(jsDoc.Build())
	}
	tempDecl.InsertBefore(nodeToDetach)

	// Go-only: the loop empties each property (it moves out the key, target and
	// default value) but leaves it in the pattern, so the iterator, which reads
	// the next property after the body, sees the same one Java saved before.
	for child, ok := objectPattern.FirstProp(); ok; child, ok = child.Next() {
		var newLHS rhino.Node
		var newRHS rhino.Node
		if key, ok := child.AsStringKey(); ok {
			// const {a: b} = obj;
			tempVarNameNode := e.createTempVarNameNode(tempVarName, tempVarType)
			var getprop rhino.Node
			if key.IsQuoted() {
				getprop = e.astFactory.CreateGetElem(
					tempVarNameNode, e.astFactory.CreateString(key.Key()))
			} else {
				getprop = e.astFactory.CreateGetProp(tempVarNameNode, key.Key(), tempVarType)
			}

			value := key.Value().Detach()
			if defaultNode, ok := value.AsDefaultValue(); ok {
				target, defaultValue := defaultNode.Target(), defaultNode.Value() // read both before detaching either
				newLHS = target.Detach().Node()
				defaultValue.Detach()
				intermediateTempVarName := e.getTempVariableName()
				intermediateDecl :=
					rhino.IR_VarWithValue(e.createTempVarNameNode(intermediateTempVarName, AstFactory_Type_(getprop)), getprop)
				intermediateDecl.SrcrefTreeIfMissing(key.Node())
				intermediateDecl.InsertBefore(nodeToDetach)
				newRHS =
					e.defaultValueHook(
						e.createTempVarNameNode(intermediateTempVarName, AstFactory_Type_(getprop)), defaultValue.Node())
			} else {
				newLHS = value.Node()
				newRHS = getprop
			}
			if hasRest {
				propsToDeleteForRest = append(propsToDeleteForRest, key.Node())
			}
		} else if computed, ok := child.AsComputedProp(); ok {
			// const {[propExpr]: newLHS = defaultValue} = newRHS;
			// Read the key and the target before detaching either. The target is
			// always there in a pattern.
			key := computed.Key()
			value, ok := computed.Value()
			if !ok {
				panic(jutil.NewIllegalStateExceptionMsg(computed.Node().String()))
			}
			defaultNode, hasDefault := value.AsDefaultValue()
			propExpr := key.Detach().Node()
			var defaultValue rhino.Node
			if hasDefault {
				target, defaultExpr := defaultNode.Target(), defaultNode.Value() // read both before detaching either
				newLHS = target.Detach().Node()
				defaultValue = defaultExpr.Detach().Node()
			} else {
				newLHS = value.Detach().Node()
				defaultValue = rhino.NilNode
			}
			if hasRest {
				// A "...rest" variable is present and result of computation must be cached
				exprEvalTempVarName := e.getTempVariableName()
				exprEvalTempVarModel :=
					e.createTempVarNameNode(exprEvalTempVarName, AstFactory_Type_(propExpr)) // clone this node
				exprEvalDecl := rhino.IR_VarWithValue(exprEvalTempVarModel.CloneNode(), propExpr)
				exprEvalDecl.SrcrefTreeIfMissing(computed.Node())
				exprEvalDecl.InsertBefore(nodeToDetach)
				propExpr = exprEvalTempVarModel.CloneNode()
				propsToDeleteForRest = append(propsToDeleteForRest, exprEvalTempVarModel.CloneNode())
			}
			if hasDefault {
				// tempVarName[propExpr]
				getelem :=
					e.astFactory.CreateGetElem(e.createTempVarNameNode(tempVarName, tempVarType), propExpr)

				// var tempVarName2 = tempVarName1[propExpr]
				intermediateTempVarName := e.getTempVariableName()
				intermediateDecl :=
					rhino.IR_VarWithValue(e.createTempVarNameNode(intermediateTempVarName, AstFactory_Type_(getelem)), getelem)
				intermediateDecl.SrcrefTreeIfMissing(computed.Node())
				intermediateDecl.InsertBefore(nodeToDetach)

				// tempVarName2 === undefined ? defaultValue : tempVarName2
				newRHS =
					e.defaultValueHook(
						e.createTempVarNameNode(intermediateTempVarName, AstFactory_Type_(getelem)), defaultValue)
			} else {
				newRHS =
					e.astFactory.CreateGetElem(e.createTempVarNameNode(tempVarName, AstFactory_Type_(newLHS)), propExpr)
			}
		} else if rest, ok := child.AsObjectRest(); ok {
			if _, hasNext := child.Next(); hasNext {
				panic(jutil.NewIllegalStateExceptionMsg("object rest may not be followed by any properties"))
			}
			// Emit Object.assign here after preceding property bindings to preserve evaluation
			// side-effect ordering (e.g. default value evaluations, property extractions).
			// TODO(b/538156291): Properties already bound are still copied during Object.assign and
			// then deleted, so getters for those properties may still be invoked during Object.assign.
			// TODO(b/116532470): see if casting this to a more specific type fixes disambiguation
			assignCall :=
				e.astFactory.CreateCall(
					e.astFactory.CreateQName(e.namespace, "Object.assign"),
					AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT))
			assignCall.AddChildToBack(e.astFactory.CreateObjectLit())
			assignCall.AddChildToBack(e.createTempVarNameNode(tempVarName, tempVarType))

			restTempDecl := rhino.IR_VarWithValue(e.createTempVarNameNode(restTempVarName, tempVarType), assignCall)
			restTempDecl.SrcrefTreeIfMissing(objectPattern.Node())
			restTempDecl.InsertBefore(nodeToDetach)

			restTarget := rest.Target() // e.g. get `rest` from `const {...rest} = {};`
			// Go-only: Java reads the string of the target whatever its kind, which
			// fails for an element access (`({...a[0]} = obj)` crashed the compiler)
			// and, for a property access whose name has the prefix, trips the check
			// that it is a name. Only a temporary name is cloned; any other target
			// (a name, a property or element access) is moved.
			if name, ok := restTarget.AsName(); ok && strings.HasPrefix(name.Ident(), Es6RewriteDestructuring_DESTRUCTURING_TEMP_VAR) {
				newLHS = e.createTempVarNameNode(name.Ident(), AstFactory_Type_(name.Node()))
			} else {
				newLHS = restTarget.Detach().Node()
			}
			newRHS = e.objectPatternRestRHS(objectPattern, rest, restTempVarName, propsToDeleteForRest)
		} else {
			panic(jutil.NewIllegalStateExceptionMsg("unexpected child"))
		}

		var newNode rhino.Node
		if NodeUtil_IsNameDeclaration(parent) {
			if parent.IsConst() {
				newLHS.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
			}
			newNode = rhino.IR_DeclarationWithValue(newLHS, newRHS, parent.GetToken())
		} else if parent.IsAssign() {
			newNode = rhino.IR_ExprResult(e.astFactory.CreateAssign(newLHS, newRHS))
		} else {
			panic(jutil.NewIllegalStateExceptionMsg("not reached"))
		}
		newNode.SrcrefTreeIfMissing(child.Node())

		newNode.InsertBefore(nodeToDetach)

		// Explicitly visit the LHS of the new node since it may be a nested
		// destructuring pattern.
		e.Visit(t, newLHS, newLHS.GetParent())
	}

	nodeToDetach.Detach()
	t.ReportCodeChange()
}

// objectPatternRestRHS converts "rest" of object destructuring lhs by making a
// clone and deleting any properties that were stated in the original object
// pattern.
//
// Nodes in statedProperties that are a stringKey will be used in a getprop
// when deleting. All other types will be used in a getelem such as what is
// done for computed properties.
//
//	  {a, [foo()]:b, ...x} = rhs;
//	becomes
//	  var temp = rhs;
//	  a = temp.a;
//	  var temp2 = foo();
//	  b = temp[temp2];
//	  var temp1 = Object.assign({}, temp);
//	  x = (delete temp1.a, delete temp1[temp2], temp1);
//
// rest is the node representing the "...rest" of objectPattern;
// restTempVarName is the name of the var containing a clone of the result of
// the rhs evaluation; statedProperties is the list of properties to delete
// from the clone.
//
// Java: Es6RewriteDestructuring.objectPatternRestRHS(Node, Node, String, ArrayList<Node>)
func (e *Es6RewriteDestructuring) objectPatternRestRHS(
	objectPattern ast.ObjectPattern, rest ast.ObjectRest, restTempVarName string, statedProperties []rhino.Node) rhino.Node {
	restTempVarModel := e.createTempVarNameNode(restTempVarName, AstFactory_Type_(objectPattern.Node()))
	result := restTempVarModel.CloneNode()
	if len(statedProperties) > 0 {
		comma := e.deletionNodeForRestProperty(restTempVarModel.CloneNode(), statedProperties[0])
		for _, prop := range statedProperties[1:] {
			comma =
				e.astFactory.CreateComma(
					comma, e.deletionNodeForRestProperty(restTempVarModel.CloneNode(), prop))
		}
		result = e.astFactory.CreateComma(comma, result)
	}
	result.SrcrefTreeIfMissing(rest.Node())
	return result
}

// deletionNodeForRestProperty returns `delete rest.key` for a property of
// the pattern (a STRING_KEY), or `delete rest[temp]` for the temporary
// variable that holds the value of a computed key (a NAME).
//
// Java: Es6RewriteDestructuring.deletionNodeForRestProperty(Node, Node)
func (e *Es6RewriteDestructuring) deletionNodeForRestProperty(restTempVarNameNode rhino.Node, property rhino.Node) rhino.Node {
	var get rhino.Node
	if key, ok := ast.AsStringKey(property); ok {
		if key.IsQuoted() {
			get = e.astFactory.CreateGetElem(
				restTempVarNameNode, e.astFactory.CreateString(key.Key()))
		} else {
			get = e.astFactory.CreateGetPropWithUnknownType(
				restTempVarNameNode, key.Key())
		}
	} else if _, ok := ast.AsName(property); ok {
		get = e.astFactory.CreateGetElem(restTempVarNameNode, property)
	} else {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected property to delete node: " + property.ToStringTree()))
	}

	return e.astFactory.CreateDelProp(get)
}

// replaceArrayPattern converts
//
//	var [x, y] = rhs
//
// to:
//
//	var temp = $jscomp.makeIterator(rhs);
//	var x = temp.next().value;
//	var y = temp.next().value;
//
// Java: Es6RewriteDestructuring.replaceArrayPattern(NodeTraversal, Node, Node, Node, Node)
func (e *Es6RewriteDestructuring) replaceArrayPattern(
	t *NodeTraversal, arrayPattern ast.ArrayPattern, rhs rhino.Node, parent rhino.Node, nodeToDetach rhino.Node) {

	if e.rewriteMode == Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST {
		if level := e.peekPatternNestingLevel(); level == nil || !level.hasNestedObjectRest {
			return
		}
	}

	if _, isLabel := ast.AsLabel(nodeToDetach.GetParent()); isLabel {
		block := rhino.IR_Block().Srcref(nodeToDetach)
		nodeToDetach.ReplaceWith(block)
		block.AddChildToBack(nodeToDetach)
	}

	tempVarName := e.getTempVariableName()
	makeIteratorCall := e.astFactory.CreateJSCompMakeIteratorCall(rhs.Detach(), e.namespace)
	tempDecl := e.astFactory.CreateSingleVarNameDeclarationWithValue(tempVarName, makeIteratorCall)
	tempVarModel := ast.MustNameDeclaration(tempDecl).MustOnlyDeclarator().Node()
	tempDecl.SrcrefTreeIfMissing(arrayPattern.Node())
	tempDecl.InsertBefore(nodeToDetach)

	var nextChild ast.ArrayPatternElem
	for child, hasChild := arrayPattern.FirstElement(); hasChild; child = nextChild {
		nextChild, hasChild = child.Next()
		if _, ok := child.AsEmpty(); ok {
			// Just call the next() method to advance the iterator, but throw away the value.
			nextCall :=
				rhino.IR_ExprResult(
					e.astFactory.CreateCallWithUnknownType(
						e.astFactory.CreateGetProp(
							tempVarModel.CloneNode(), "next", AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT))))
			nextCall.SrcrefTreeIfMissing(child.Node())
			nextCall.InsertBefore(nodeToDetach)
			continue
		}

		var newLHS rhino.Node
		var newRHS rhino.Node
		if defaultNode, ok := child.AsDefaultValue(); ok {
			//   [x = defaultValue] = rhs;
			// becomes
			//   var temp0 = $jscomp.makeIterator(rhs);
			//   var temp1 = temp.next().value
			//   x = (temp1 === undefined) ? defaultValue : temp1;
			nextVarName := e.getTempVariableName()
			// `temp.next().value`
			nextCallDotValue :=
				e.astFactory.CreateGetPropWithUnknownType(
					e.astFactory.CreateCallWithUnknownType(
						e.astFactory.CreateGetPropWithUnknownType(tempVarModel.CloneNode(), "next")),
					"value")
			nextVarType := AstFactory_Type_(nextCallDotValue)
			// `var temp1 = temp.next().value`
			var_ := rhino.IR_VarWithValue(e.createTempVarNameNode(nextVarName, nextVarType), nextCallDotValue)
			var_.SrcrefTreeIfMissing(defaultNode.Node())
			var_.InsertBefore(nodeToDetach)

			target, defaultValue := defaultNode.Target(), defaultNode.Value() // read both before detaching either
			// `x`
			newLHS = target.Detach().Node()
			// `(temp1 === undefined) ? defaultValue : temp1;
			newRHS =
				e.defaultValueHook(
					e.createTempVarNameNode(nextVarName, nextVarType), defaultValue.Detach().Node())
		} else if rest, ok := child.AsIterRest(); ok {
			//   [...x] = rhs;
			// becomes
			//   var temp = $jscomp.makeIterator(rhs);
			//   x = $jscomp.arrayFromIterator(temp);
			newLHS = rest.Target().Detach().Node()
			newRHS = e.astFactory.CreateJscompArrayFromIteratorCall(tempVarModel.CloneNode(), e.namespace)
		} else {
			// LHS is just a name (or a nested pattern).
			//   var [x] = rhs;
			// becomes
			//   var temp = $jscomp.makeIterator(rhs);
			//   var x = temp.next().value;
			newLHS = child.Detach().Node()
			newRHS =
				e.astFactory.CreateGetProp(
					e.astFactory.CreateCallWithUnknownType(
						e.astFactory.CreateGetPropWithUnknownType(tempVarModel.CloneNode(), "next")),
					"value",
					AstFactory_Type_(newLHS))
		}
		var newNode rhino.Node
		if parent.IsAssign() {
			assignment := e.astFactory.CreateAssign(newLHS, newRHS)
			newNode = rhino.IR_ExprResult(assignment)
		} else {
			newNode = rhino.IR_DeclarationWithValue(newLHS, newRHS, parent.GetToken())
		}
		newNode.SrcrefTreeIfMissing(arrayPattern.Node())

		newNode.InsertBefore(nodeToDetach)
		// Explicitly visit the LHS of the new node since it may be a nested
		// destructuring pattern.
		e.Visit(t, newLHS, newLHS.GetParent())
	}

	nodeToDetach.Detach()
	t.ReportCodeChange()
}

// wrapAssignOrDestructuringInCallToArrow replaces ASSIGN or
// DESTRUCTURING_LHS with a IIFE that contains the transpiled destructuring.
//
//	Transform
//	  [x, y] = rhs
//	into
//	  ((temp0) => {
//	    var temp1 = $jscomp.makeIterator(temp0);
//	    var x = temp0.next().value;
//	    var y = temp0.next().value;
//	    return temp0;
//	  })(rhs)
//
//	Transform
//	  {x: a, y: b} = rhs
//	into
//	  ((temp0) => {
//	    var temp1 = temp0;
//	    var a = temp0.x;
//	    var b = temp0.y;
//	    return temp0;
//	  })(rhs)
//
// Java: Es6RewriteDestructuring.wrapAssignOrDestructuringInCallToArrow(NodeTraversal, Node)
func (e *Es6RewriteDestructuring) wrapAssignOrDestructuringInCallToArrow(t *NodeTraversal, assignment rhino.Node) {
	var lhs ast.Pattern
	var rhs rhino.Node
	if assign, ok := ast.AsAssign(assignment); ok {
		lhs = ast.MustPattern(assign.Target().Node())
		rhs = assign.Value().Detach().Node()
	} else {
		dlhs := ast.MustDestructuringLhs(assignment)
		lhs = dlhs.Pattern()
		initExpr, hasInit := dlhs.Init()
		if !hasInit {
			panic(jutil.NewIllegalStateExceptionMsg("missing destructuring initializer"))
		}
		rhs = initExpr.Detach().Node()
	}
	// NOTE: we do not presently support await/yield in the LHS of nested destructuring assignments.
	// See b/475296868.
	e.validateNoAwaitOrYieldInNestedAssignmentPattern(lhs)

	tempVarName := e.getTempVariableName()

	tempVarModel := e.createTempVarNameNode(tempVarName, AstFactory_Type_(rhs))
	//  ((temp0) => {...})
	paramList := rhino.IR_ParamList(tempVarModel.CloneNode())
	// [x, y] = temp0;
	patternValue := tempVarModel.CloneNode()
	replacementExpr := e.astFactory.CreateAssign(lhs.Detach().Node(), patternValue)
	exprResult := rhino.IR_ExprResult(replacementExpr)
	// return temp0;
	returnNode := rhino.IR_ReturnNodeWithExpr(tempVarModel.CloneNode())

	// Create a function to hold these assignments:
	block := rhino.IR_BlockOfNodes(exprResult, returnNode)
	arrowFn :=
		e.astFactory.CreateFunction(
			/* name= */ "",
			paramList,
			block,
			AstFactory_Type_OfJSTypeNativeAndColor(rhino.JSTypeNative_UNKNOWN_TYPE, colors.StandardColors_UNKNOWN))
	arrowFn.SetIsArrowFunction(true)

	// Create a call to the function, and replace the pattern with the call.
	call := e.astFactory.CreateCall(arrowFn, AstFactory_Type_(rhs), rhs)
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_ARROW_FUNCTIONS, e.compiler)
	call.SrcrefTreeIfMissing(assignment)
	ast.MustCall(call).SetFreeCall(true)
	assignment.ReplaceWith(call)
	NodeUtil_MarkNewScopesChanged(call, e.compiler)

	e.replacePattern(t, lhs, patternValue, replacementExpr, exprResult)
}

func isAwaitOrYieldNode(n rhino.Node) bool {
	return n.IsAwait() || n.IsYield()
}

// Java: Es6RewriteDestructuring.validateNoAwaitOrYieldInNestedAssignmentPattern(Node)
func (e *Es6RewriteDestructuring) validateNoAwaitOrYieldInNestedAssignmentPattern(pattern ast.Pattern) {
	if !NodeUtil_FindPreorder(pattern.Node(), isAwaitOrYieldNode, NodeUtil_MATCH_NOT_FUNCTION).IsNil() {
		panic(jutil.NewIllegalStateExceptionMsg(
			"Cannot transpile yet: destructuring assignment referencing await or yield in lhs, in" +
				" nested sub-expression: " + pattern.Node().String()))
	}
}

// visitDestructuringPatternInVanillaForInnerVars handles
// for (let [a, b, c] = arr; a < b; a++)
//
// Go-only: it takes the declarator of the pattern; Java takes the pattern.
//
// Java: Es6RewriteDestructuring.visitDestructuringPatternInVanillaForInnerVars(NodeTraversal, Node)
func (e *Es6RewriteDestructuring) visitDestructuringPatternInVanillaForInnerVars(t *NodeTraversal, destructuringLhs ast.DestructuringLhs) {
	pattern := destructuringLhs.Pattern()
	lhs := destructuringLhs.Node()
	declaration := ast.MustNameDeclaration(lhs.GetParent())

	insertionPoint := declaration.Node().GetParent()
	for _, ok := ast.AsLabel(insertionPoint.GetParent()); ok; _, ok = ast.AsLabel(insertionPoint.GetParent()) {
		insertionPoint = insertionPoint.GetParent()
	}

	switch declaration.Token() {
	case rhino.Token_CONST:
		{
			block := rhino.IR_Block().Srcref(insertionPoint)
			insertionPoint.ReplaceWith(block)
			block.AddChildToBack(insertionPoint)
		}
		// Fall through
		fallthrough

	case rhino.Token_VAR:
		{
			// Move any earlier variables out of the loop initializer
			for c, ok := declaration.FirstDeclarator(); ok && !c.Is(lhs); c, ok = declaration.FirstDeclarator() {
				newDeclaration := declaration.Node().CloneNode()
				newDeclaration.AddChildToBack(c.Detach().Node())
				newDeclaration.InsertBefore(insertionPoint)
			}

			// Move the pattern out of the initializer and transpile it
			rhs, ok := destructuringLhs.Init()
			if !ok {
				// a pattern declared in a for initializer has a value
				panic(jutil.NewIllegalStateException())
			}
			newDeclaration := declaration.Node().CloneNode()
			newDeclaration.AddChildToBack(destructuringLhs.Detach().Node())
			newDeclaration.InsertBefore(insertionPoint)
			e.replacePattern(t, pattern, rhs.Node(), newDeclaration, newDeclaration)

			if !declaration.HasDeclarators() {
				declaration.Node().ReplaceWith(rhino.IR_Empty())
			}
		}

	case rhino.Token_LET:
		// See https://tc39.es/ecma262/#sec-createperiterationenvironment
		// for (let a, b, c, unusedTmp = (() => [a, b, c] = arr); a < b; a++)
		{
			parsing.ParsingUtil_GetParamOrPatternNames(
				pattern.Node(), func(name rhino.Node) { name.CloneNode().InsertBefore(lhs) })

			unusedVar :=
				e.astFactory.
					CreateNameWithUnknownType(e.getTempVariableName() + "$unused").
					Srcref(lhs)
			lhs.ReplaceWith(unusedVar)
			unusedVar.AddChildToBack(lhs)
			e.wrapAssignOrDestructuringInCallToArrow(t, lhs)
		}
	}

	e.compiler.ReportChangeToEnclosingScope(insertionPoint)
}

// visitDestructuringPatternInEnhancedForInnerVars handles
// for (const [a, b, c] of arr)
//
// Go-only: it takes the loop and the declarator of the pattern; Java takes
// the pattern.
//
// Java: Es6RewriteDestructuring.visitDestructuringPatternInEnhancedForInnerVars(Node)
func (e *Es6RewriteDestructuring) visitDestructuringPatternInEnhancedForInnerVars(forNode ast.EnhancedFor, destructuringLhs ast.DestructuringLhs) {
	tempVarName := e.getTempVariableName()

	pattern := destructuringLhs.Pattern().Node()
	declarationType := ast.MustNameDeclaration(destructuringLhs.Node().GetParent()).Token()
	block := forNode.Body().Node()

	destructuringLhs.Node().ReplaceWith(e.createTempVarNameNode(tempVarName, AstFactory_Type_(pattern)).Srcref(pattern))
	decl :=
		rhino.IR_DeclarationWithValue(
			pattern.Detach(), e.createTempVarNameNode(tempVarName, AstFactory_Type_(pattern)), declarationType)
	decl.SrcrefTreeIfMissing(pattern)
	// Move the body into an inner block to handle cases where declared variables in the for
	// loop initializer are shadowed by variables in the for loop body. e.g.
	//   for (const [value] of []) { const value = 1; }
	newBlock := rhino.IR_BlockOfNode(decl)
	block.ReplaceWith(newBlock)
	newBlock.AddChildToBack(block)
}

// visitDestructuringPatternInEnhancedForWithOuterVars handles
// for ([a, b, c] of arr)
//
//	Transform
//	  for ({x} of y) {}
//	into
//	  var TEMP_VAR0;
//	  for (TEMP_VAR0 of y) {
//	    var TEMP_VAR1 = TEMP_VAR0;
//	    x = TEMP_VAR1.x;
//	  }
//
// Go-only: it takes the loop too; Java takes the pattern alone.
//
// Java: Es6RewriteDestructuring.visitDestructuringPatternInEnhancedForWithOuterVars(Node)
func (e *Es6RewriteDestructuring) visitDestructuringPatternInEnhancedForWithOuterVars(forNode ast.EnhancedFor, patternTarget ast.Pattern) {
	tempVarName := e.getTempVariableName()

	pattern := patternTarget.Node()
	block := forNode.Body().Node()

	name := e.createTempVarNameNode(tempVarName, AstFactory_Type_(pattern))
	decl := rhino.IR_Var(name)
	decl.SrcrefTreeIfMissing(pattern)
	insertionPoint := forNode.Node()
	for !insertionPoint.IsNil() && !insertionPoint.GetParent().IsNil() && insertionPoint.GetParent().IsLabel() {
		insertionPoint = insertionPoint.GetParent()
	}
	decl.InsertBefore(insertionPoint)
	clonedName := name.CloneNode()
	clonedName.SrcrefTreeIfMissing(pattern)
	pattern.ReplaceWith(clonedName)
	exprResult :=
		rhino.IR_ExprResult(
			e.astFactory.CreateAssign(pattern, e.createTempVarNameNode(tempVarName, AstFactory_Type_(pattern))))
	exprResult.SrcrefTreeIfMissing(pattern)
	block.AddChildToFront(exprResult)
}

// Go-only: it takes the catch clause too; Java takes the pattern alone.
//
// Java: Es6RewriteDestructuring.visitDestructuringPatternInCatch(NodeTraversal, Node)
func (e *Es6RewriteDestructuring) visitDestructuringPatternInCatch(t *NodeTraversal, catch ast.Catch, patternParam ast.Pattern) {
	tempVarName := e.getTempVariableName()
	catchBlock := catch.Body().Node()
	pattern := patternParam.Node()
	patternType := AstFactory_Type_(pattern)

	pattern.ReplaceWith(e.createTempVarNameNode(tempVarName, patternType))
	catchBlock.AddChildToFront(
		rhino.IR_DeclarationWithValue(pattern, e.createTempVarNameNode(tempVarName, patternType), rhino.Token_LET))
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_LET_DECLARATIONS, e.compiler)
}

// defaultValueHook is a helper for transpiling DEFAULT_VALUE trees.
//
// Java: Es6RewriteDestructuring.defaultValueHook(Node, Node)
func (e *Es6RewriteDestructuring) defaultValueHook(getprop rhino.Node, defaultValue rhino.Node) rhino.Node {
	undefined := e.astFactory.CreateUndefinedValue()
	undefined.MakeNonIndexable()
	getpropClone := getprop.CloneTree().SetColor(getprop.GetColor())
	return e.astFactory.CreateHook(
		e.astFactory.CreateSheq(getprop, undefined), defaultValue, getpropClone)
}
