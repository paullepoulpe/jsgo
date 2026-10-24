// Port of com/google/javascript/jscomp/Es6RewriteDestructuring.java @ closure-compiler v20261021

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

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Es6RewriteDestructuring rewrites destructuring patterns and default
// parameters to valid ES3 code or to a different form of destructuring.
type Es6RewriteDestructuring struct {
	compiler    AbstractCompiler
	rewriteMode Es6RewriteDestructuring_ObjectDestructuringRewriteMode
	astFactory  *AstFactory
	namespace   rhino.StaticScope

	featuresToTriggerRunningPass *parser.FeatureSet
	featuresToMarkAsRemoved      *parser.FeatureSet

	// Go-only: language-adapter — the ArrayDeque is a slice: addLast/removeLast
	// append to and drop from its end, and iteration runs first to last.
	patternNestingStack []*Es6RewriteDestructuring_PatternNestingLevel

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

// Es6RewriteDestructuring_ObjectDestructuringRewriteMode_values is
// ObjectDestructuringRewriteMode.values().
//
// Java: ObjectDestructuringRewriteMode.values()
func Es6RewriteDestructuring_ObjectDestructuringRewriteMode_values() []Es6RewriteDestructuring_ObjectDestructuringRewriteMode {
	return []Es6RewriteDestructuring_ObjectDestructuringRewriteMode{
		Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_ALL_OBJECT_PATTERNS,
		Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST,
	}
}

// String is Enum.name().
//
// Java: ObjectDestructuringRewriteMode.name()
func (e Es6RewriteDestructuring_ObjectDestructuringRewriteMode) String() string {
	return es6RewriteDestructuring_ObjectDestructuringRewriteMode_names[e]
}

// Ordinal is Enum.ordinal().
//
// Java: ObjectDestructuringRewriteMode.ordinal()
func (e Es6RewriteDestructuring_ObjectDestructuringRewriteMode) Ordinal() int { return int(e) }

// Java: Es6RewriteDestructuring.DESTRUCTURING_TEMP_VAR
const Es6RewriteDestructuring_DESTRUCTURING_TEMP_VAR = "$jscomp$destructuring$var"

// Java: Es6RewriteDestructuring(Builder)
func newEs6RewriteDestructuring(builder *Es6RewriteDestructuring_Builder) *Es6RewriteDestructuring {
	e := &Es6RewriteDestructuring{}
	e.patternNestingStack = nil
	e.destructuringVarCounter = 0
	e.compiler = builder.compiler
	e.rewriteMode = builder.rewriteMode
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

// Es6RewriteDestructuring_Builder is Es6RewriteDestructuring.Builder.
type Es6RewriteDestructuring_Builder struct {
	compiler    AbstractCompiler
	rewriteMode Es6RewriteDestructuring_ObjectDestructuringRewriteMode
}

// Java: Es6RewriteDestructuring.Builder(AbstractCompiler)
func NewEs6RewriteDestructuring_Builder(compiler AbstractCompiler) *Es6RewriteDestructuring_Builder {
	b := &Es6RewriteDestructuring_Builder{}
	b.rewriteMode = Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_ALL_OBJECT_PATTERNS
	b.compiler = compiler
	return b
}

// Java: Es6RewriteDestructuring.Builder.setDestructuringRewriteMode(ObjectDestructuringRewriteMode)
func (b *Es6RewriteDestructuring_Builder) SetDestructuringRewriteMode(rewriteMode Es6RewriteDestructuring_ObjectDestructuringRewriteMode) *Es6RewriteDestructuring_Builder {
	b.rewriteMode = rewriteMode
	return b
}

// Java: Es6RewriteDestructuring.Builder.build()
func (b *Es6RewriteDestructuring_Builder) Build() *Es6RewriteDestructuring {
	return newEs6RewriteDestructuring(b)
}

// Es6RewriteDestructuring_PatternNestingLevel is
// Es6RewriteDestructuring.PatternNestingLevel.
type Es6RewriteDestructuring_PatternNestingLevel struct {
	pattern             rhino.Node
	hasNestedObjectRest bool
}

// Java: PatternNestingLevel(Node, boolean)
func NewEs6RewriteDestructuring_PatternNestingLevel(pattern rhino.Node, hasNestedRest bool) *Es6RewriteDestructuring_PatternNestingLevel {
	l := &Es6RewriteDestructuring_PatternNestingLevel{}
	l.pattern = pattern
	l.hasNestedObjectRest = hasNestedRest
	return l
}

// patternNestingStackPeekLast is Deque.peekLast(): the last element, or null.
//
// Go-only: language-adapter — the ArrayDeque operation on the slice stack.
func (e *Es6RewriteDestructuring) patternNestingStackPeekLast() *Es6RewriteDestructuring_PatternNestingLevel {
	if len(e.patternNestingStack) == 0 {
		return nil
	}
	return e.patternNestingStack[len(e.patternNestingStack)-1]
}

// patternNestingStackGetLast is Deque.getLast(): the last element, or
// NoSuchElementException.
//
// Go-only: language-adapter — the ArrayDeque operation on the slice stack.
func (e *Es6RewriteDestructuring) patternNestingStackGetLast() *Es6RewriteDestructuring_PatternNestingLevel {
	if len(e.patternNestingStack) == 0 {
		panic(jutil.NewNoSuchElementException())
	}
	return e.patternNestingStack[len(e.patternNestingStack)-1]
}

// Java: Es6RewriteDestructuring.process(Node, Node)
func (e *Es6RewriteDestructuring) Process(externs rhino.Node, root rhino.Node) {
	jutil.CheckState(len(e.patternNestingStack) == 0)
	NodeTraversal_Traverse(e.compiler, root, e)
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(e.compiler, root, e.featuresToMarkAsRemoved)
	jutil.CheckState(len(e.patternNestingStack) == 0)
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
		e.PullDestructuringOutOfParams(n, parent)
	case rhino.Token_ARRAY_PATTERN, rhino.Token_OBJECT_PATTERN:
		hasRest := n.IsObjectPattern() && n.HasChildren() && n.GetLastChild().IsRest()
		if len(e.patternNestingStack) != 0 && hasRest {
			for _, level := range e.patternNestingStack {
				if level.hasNestedObjectRest {
					break
				}
				level.hasNestedObjectRest = true
			}
			e.patternNestingStackPeekLast().hasNestedObjectRest = true
		}
		e.patternNestingStack = append(e.patternNestingStack, NewEs6RewriteDestructuring_PatternNestingLevel(n, hasRest))
	default:
	}
	return true
}

// Java: Es6RewriteDestructuring.visit(NodeTraversal, Node, Node)
func (e *Es6RewriteDestructuring) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_ARRAY_PATTERN, rhino.Token_OBJECT_PATTERN:
		e.VisitPattern(t, n)
		if n == e.patternNestingStackGetLast().pattern {
			e.patternNestingStack = e.patternNestingStack[:len(e.patternNestingStack)-1]
		}
	default:
	}
}

// PullDestructuringOutOfParams pulls all default and destructuring
// parameters out of function parameters.
//
// TODO(bradfordcsmith): Ideally if we're only removing OBJECT_REST, we should
// only do this when the parameter list contains a usage of OBJECT_REST.
//
// Java: Es6RewriteDestructuring.pullDestructuringOutOfParams(Node, Node)
func (e *Es6RewriteDestructuring) PullDestructuringOutOfParams(paramList rhino.Node, function rhino.Node) {
	var insertSpot rhino.Node = rhino.NilNode
	body := function.GetLastChild()
	var next rhino.Node = rhino.NilNode
	for param := paramList.GetFirstChild(); !param.IsNil(); param = next {
		next = param.GetNext()
		if param.IsDefaultValue() {
			nameOrPattern := param.RemoveFirstChild()
			// We'll be cloning nameOrPattern below, and we don't want to clone the JSDoc info with it
			jsDoc := nameOrPattern.GetJSDocInfo()
			nameOrPattern.SetJSDocInfo(nil)
			defaultValue := param.RemoveFirstChild()
			var newParam rhino.Node

			// Treat name=undefined (and equivalent) as if it was just name.  There
			// is no need to generate a (name===void 0?void 0:name) statement for
			// such arguments.
			isNoop := false
			if !nameOrPattern.IsName() {
				// Do not try to optimize unless nameOrPattern is a simple name.
			} else if defaultValue.IsName() {
				isNoop = "undefined" == defaultValue.GetString()
			} else if defaultValue.IsVoid() {
				// Any kind of 'void literal' is fine, but 'void fun()' or anything
				// else with side effects isn't.  We're not trying to be particularly
				// smart here and treat 'void {}' for example as if it could cause side effects.
				isNoop = NodeUtil_IsImmutableValue(defaultValue.GetFirstChild())
			}

			if isNoop {
				newParam = nameOrPattern.CloneTree()
			} else {
				if nameOrPattern.IsName() {
					newParam = nameOrPattern
				} else {
					newParam = e.CreateTempVarNameNode(e.GetTempVariableName(), AstFactory_Type_(nameOrPattern))
				}
				lhs := nameOrPattern.CloneTree()
				rhs := e.DefaultValueHook(newParam.CloneTree(), defaultValue)
				var newStatement rhino.Node
				if nameOrPattern.IsName() {
					newStatement = rhino.IR_ExprResult(e.astFactory.CreateAssign(lhs, rhs))
				} else {
					newStatement = rhino.IR_VarWithValue(lhs, rhs)
				}
				newStatement.SrcrefTreeIfMissing(param)
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
				if newStatement.GetFirstChild().IsDestructuringLhs() {
					insertSpot = e.RewriteVarDestructuringDeclaration(newStatement.GetFirstChild())
				}
			}

			param.ReplaceWith(newParam)
			newParam.SetJSDocInfo(jsDoc)

			e.compiler.ReportChangeToChangeScope(function)
		} else if param.IsDestructuringPattern() {
			insertSpot =
				e.ReplacePatternParamWithTempVar(function, insertSpot, param, e.GetTempVariableName())
			e.compiler.ReportChangeToChangeScope(function)
		} else if param.IsRest() && param.GetFirstChild().IsDestructuringPattern() {
			insertSpot =
				e.ReplacePatternParamWithTempVar(
					function, insertSpot, param.GetFirstChild(), e.GetTempVariableName())
			e.compiler.ReportChangeToChangeScope(function)
		}
	}
}

// ReplacePatternParamWithTempVar replaces a destructuring pattern parameter
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
func (e *Es6RewriteDestructuring) ReplacePatternParamWithTempVar(
	function rhino.Node, insertSpot rhino.Node, patternParam rhino.Node, tempVarName string) rhino.Node {
	// Convert `function f([a, b]) {}` to `function f(tempVar) { var [a, b] = tempVar; }`
	paramType := AstFactory_Type_(patternParam)
	newParam := e.CreateTempVarNameNode(tempVarName, paramType)
	newParam.SetJSDocInfo(patternParam.GetJSDocInfo())
	patternParam.ReplaceWith(newParam)
	newDecl := rhino.IR_VarWithValue(patternParam, e.CreateTempVarNameNode(tempVarName, paramType))
	newDecl.SrcrefTreeIfMissing(patternParam)

	if insertSpot.IsNil() {
		// insert new declarations only after all inner function declarations in that function body to
		// preserve normalization
		insertSpot =
			NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(function.GetLastChild())
		if !insertSpot.IsNil() {
			newDecl.InsertBefore(insertSpot)
		} else {
			// functionBody only contains hoisted function declarations
			function.GetLastChild().AddChildToBack(newDecl)
		}
	} else {
		newDecl.InsertAfter(insertSpot)
	}

	return e.RewriteVarDestructuringDeclaration(newDecl.GetFirstChild())
}

// RewriteVarDestructuringDeclaration rewrites a destructuring declaration to
// a destructuring assignment.
//
// For example, given the code `var [a, b] = ...;`, this function will convert
// it to `var a; var b; [a, b] = ...;`
//
// destructuringLhs is the destructuring declaration to rewrite. Returns the
// expression result containing the assignment.
//
// Java: Es6RewriteDestructuring.rewriteVarDestructuringDeclaration(Node)
func (e *Es6RewriteDestructuring) RewriteVarDestructuringDeclaration(destructuringLhs rhino.Node) rhino.Node {
	jutil.CheckStateMsg1(destructuringLhs.IsDestructuringLhs(), "%s", destructuringLhs)
	jutil.CheckStateMsg(
		!destructuringLhs.GetGrandparent().IsExport(),
		"Export destructuring declarations not expected inside a function body.")
	var_ := destructuringLhs.GetParent()

	// create a stub declaration for each name in the destructuring pattern
	// e.g. generate `var a; var b;` from `var [a, b] = ...;`
	NodeUtil_VisitLhsNodesInNode(
		destructuringLhs,
		func(name rhino.Node) {
			// Add a declaration outside the destructuring pattern for the given name.
			jutil.CheckStateMsg1(
				name.IsName(),
				"lhs in destructuring declaration should be a simple name. (%s)",
				name)
			newName := rhino.IR_Name(name.GetString()).Srcref(name)
			if name.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME) {
				// if old name was a const, new name should be too
				// e.g. when rewriting `{VALUE} = ...` the `VALUE` is const by coding convention
				newName.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, true)
			}
			newVar := rhino.IR_Var(newName).Srcref(name)
			newVar.InsertBefore(var_)
		})

	// Transform destructuring var declaration to assignment. That is, `var [a, b] = ...` to `[a,
	// b] = ...` and `var {a, b} = ...` to `({a, b} = ...);`
	destructuringPattern := destructuringLhs.RemoveFirstChild()
	if !destructuringPattern.IsDestructuringPattern() {
		jutil.CheckStateMsg1(
			false,
			"Expected destructuring pattern but got %s",
			destructuringPattern.ToStringTree())
	}

	rhs := destructuringLhs.RemoveFirstChild()
	assign := e.astFactory.CreateAssign(destructuringPattern, rhs)
	assign.Srcref(var_)
	expr := e.astFactory.ExprResult(assign)
	expr.Srcref(var_)

	var_.ReplaceWith(expr)
	return expr
}

// CreateTempVarNameNode creates a new name node and adds the constant name
// property because a constant variable is used (the compiler only assigns
// $jscomp$destructuring$var[num] once)
//
// Java: Es6RewriteDestructuring.createTempVarNameNode(String, AstFactory.Type)
func (e *Es6RewriteDestructuring) CreateTempVarNameNode(name string, type_ AstFactory_Type) rhino.Node {
	// NOTE: This does not really create a constant node as this pass runs before normalization. See
	// b/322009741.
	return e.astFactory.CreateConstantName(name, type_)
}

// GetTempVariableName creates a new unique name to use for a pattern we need
// to rewrite.
//
// Java: Es6RewriteDestructuring.getTempVariableName()
func (e *Es6RewriteDestructuring) GetTempVariableName() string {
	counter := e.destructuringVarCounter
	e.destructuringVarCounter++
	return Es6RewriteDestructuring_DESTRUCTURING_TEMP_VAR + strconv.Itoa(counter)
}

// Java: Es6RewriteDestructuring.visitPattern(NodeTraversal, Node)
func (e *Es6RewriteDestructuring) VisitPattern(t *NodeTraversal, pattern rhino.Node) {
	parent := pattern.GetParent()

	switch parent.GetToken() {
	case rhino.Token_DESTRUCTURING_LHS:
		{
			declaration := parent.GetParent()
			declarationParent := declaration.GetParent()
			if declarationParent.IsVanillaFor() {
				e.VisitDestructuringPatternInVanillaForInnerVars(t, pattern)
			} else if NodeUtil_IsEnhancedFor(declarationParent) {
				e.VisitDestructuringPatternInEnhancedForInnerVars(pattern)
			} else {
				e.ReplacePattern(t, pattern, pattern.GetNext(), declaration, declaration)
			}
		}
	case rhino.Token_ASSIGN:
		if parent.GetParent().IsExprResult() {
			e.ReplacePattern(t, pattern, pattern.GetNext(), parent, parent.GetParent())
		} else {
			e.WrapAssignOrDestructuringInCallToArrow(t, parent)
		}
	case rhino.Token_OBJECT_REST, rhino.Token_ITER_REST, rhino.Token_STRING_KEY, rhino.Token_ARRAY_PATTERN, rhino.Token_DEFAULT_VALUE, rhino.Token_COMPUTED_PROP:
		// Nested pattern; do nothing. We will visit it after rewriting the parent.
	case rhino.Token_FOR_OF, rhino.Token_FOR_IN, rhino.Token_FOR_AWAIT_OF:
		e.VisitDestructuringPatternInEnhancedForWithOuterVars(pattern)
	case rhino.Token_CATCH:
		e.VisitDestructuringPatternInCatch(t, pattern)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("unexpected parent"))
	}
}

// ReplacePattern transpiles a destructuring pattern in a declaration or
// assignment to ES5
//
// nodeToDetach is a statement node containing the pattern. This method will
// replace the node with one or more other statements.
//
// Java: Es6RewriteDestructuring.replacePattern(NodeTraversal, Node, Node, Node, Node)
func (e *Es6RewriteDestructuring) ReplacePattern(
	t *NodeTraversal, pattern rhino.Node, rhs rhino.Node, parent rhino.Node, nodeToDetach rhino.Node) {
	jutil.CheckArgumentMsg1(NodeUtil_IsStatement(nodeToDetach), "%s", nodeToDetach)
	switch pattern.GetToken() {
	case rhino.Token_ARRAY_PATTERN:
		e.ReplaceArrayPattern(t, pattern, rhs, parent, nodeToDetach)
	case rhino.Token_OBJECT_PATTERN:
		e.ReplaceObjectPattern(t, pattern, rhs, parent, nodeToDetach)
	default:
		panic(jutil.NewIllegalStateExceptionMsg("unexpected"))
	}
}

// ReplaceObjectPattern converts 'var {a: b, c: d} = rhs' to:
//
// @const var temp = rhs; var b = temp.a; var d = temp.c;
//
// Java: Es6RewriteDestructuring.replaceObjectPattern(NodeTraversal, Node, Node, Node, Node)
func (e *Es6RewriteDestructuring) ReplaceObjectPattern(
	t *NodeTraversal, objectPattern rhino.Node, rhs rhino.Node, parent rhino.Node, nodeToDetach rhino.Node) {
	tempVarName := e.GetTempVariableName()
	tempVarType := AstFactory_Type_(objectPattern)

	var restTempVarName *string = nil
	// If the last child is a rest node we will want a list of the stated properties so we can
	// exclude them from being written to the rest variable.
	var propsToDeleteForRest *jutil.ArrayList[rhino.Node] = nil
	if objectPattern.HasChildren() && objectPattern.GetLastChild().IsRest() {
		propsToDeleteForRest = jutil.NewArrayList[rhino.Node]()
		name := e.GetTempVariableName()
		restTempVarName = &name
	} else if e.rewriteMode == Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST {
		// We are configured to only break object pattern rest, but this destructure has none.
		if !e.patternNestingStackPeekLast().hasNestedObjectRest {
			// Replacement is performed after the post-order visit has reached the root pattern node, so
			// peeking last represents if there is a rest property anywhere in the entire pattern. All
			// nesting levels of lower levels have already been popped.
			e.destructuringVarCounter--
			return
		}
	}

	// create the declaration `var temp = rhs;`
	tempDecl :=
		rhino.IR_VarWithValue(e.CreateTempVarNameNode(tempVarName, tempVarType), rhs.Detach()).
			SrcrefTreeIfMissing(objectPattern)
	// TODO(tbreisacher): Remove the "if" and add this JSDoc unconditionally.
	if parent.IsConst() {
		jsDoc := rhino.JSDocInfo_Builder_()
		jsDoc.RecordConstancy()
		tempDecl.SetJSDocInfo(jsDoc.Build())
	}
	tempDecl.InsertBefore(nodeToDetach)

	var next rhino.Node
	for child := objectPattern.GetFirstChild(); !child.IsNil(); child = next {
		next = child.GetNext()

		var newLHS rhino.Node
		var newRHS rhino.Node
		if child.IsStringKey() {
			// const {a: b} = obj;
			tempVarNameNode := e.CreateTempVarNameNode(tempVarName, tempVarType)
			var getprop rhino.Node
			if child.IsQuotedStringKey() {
				getprop = e.astFactory.CreateGetElem(
					tempVarNameNode, e.astFactory.CreateString(child.GetString()))
			} else {
				getprop = e.astFactory.CreateGetProp(tempVarNameNode, child.GetString(), tempVarType)
			}

			value := child.RemoveFirstChild()
			if !value.IsDefaultValue() {
				newLHS = value
				newRHS = getprop
			} else {
				newLHS = value.RemoveFirstChild()
				defaultValue := value.RemoveFirstChild()
				intermediateTempVarName := e.GetTempVariableName()
				intermediateDecl :=
					rhino.IR_VarWithValue(e.CreateTempVarNameNode(intermediateTempVarName, AstFactory_Type_(getprop)), getprop)
				intermediateDecl.SrcrefTreeIfMissing(child)
				intermediateDecl.InsertBefore(nodeToDetach)
				newRHS =
					e.DefaultValueHook(
						e.CreateTempVarNameNode(intermediateTempVarName, AstFactory_Type_(getprop)), defaultValue)
			}
			if propsToDeleteForRest != nil {
				propsToDeleteForRest.Add(child)
			}
		} else if child.IsComputedProp() {
			// const {[propExpr]: newLHS = defaultValue} = newRHS;
			hasDefault := child.GetLastChild().IsDefaultValue()
			var defaultValue rhino.Node
			propExpr := child.RemoveFirstChild()
			if hasDefault {
				defaultNode := child.GetLastChild()
				newLHS = defaultNode.RemoveFirstChild()
				defaultValue = defaultNode.RemoveFirstChild()
			} else {
				newLHS = child.RemoveFirstChild()
				defaultValue = rhino.NilNode
			}
			if propsToDeleteForRest != nil {
				// A "...rest" variable is present and result of computation must be cached
				exprEvalTempVarName := e.GetTempVariableName()
				exprEvalTempVarModel :=
					e.CreateTempVarNameNode(exprEvalTempVarName, AstFactory_Type_(propExpr)) // clone this node
				exprEvalDecl := rhino.IR_VarWithValue(exprEvalTempVarModel.CloneNode(), propExpr)
				exprEvalDecl.SrcrefTreeIfMissing(child)
				exprEvalDecl.InsertBefore(nodeToDetach)
				propExpr = exprEvalTempVarModel.CloneNode()
				propsToDeleteForRest.Add(exprEvalTempVarModel.CloneNode())
			}
			if hasDefault {
				// tempVarName[propExpr]
				getelem :=
					e.astFactory.CreateGetElem(e.CreateTempVarNameNode(tempVarName, tempVarType), propExpr)

				// var tempVarName2 = tempVarName1[propExpr]
				intermediateTempVarName := e.GetTempVariableName()
				intermediateDecl :=
					rhino.IR_VarWithValue(e.CreateTempVarNameNode(intermediateTempVarName, AstFactory_Type_(getelem)), getelem)
				intermediateDecl.SrcrefTreeIfMissing(child)
				intermediateDecl.InsertBefore(nodeToDetach)

				// tempVarName2 === undefined ? defaultValue : tempVarName2
				newRHS =
					e.DefaultValueHook(
						e.CreateTempVarNameNode(intermediateTempVarName, AstFactory_Type_(getelem)), defaultValue)
			} else {
				newRHS =
					e.astFactory.CreateGetElem(e.CreateTempVarNameNode(tempVarName, AstFactory_Type_(newLHS)), propExpr)
			}
		} else if child.IsRest() {
			if !next.IsNil() {
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
			assignCall.AddChildToBack(e.CreateTempVarNameNode(tempVarName, tempVarType))

			restTempDecl := rhino.IR_VarWithValue(e.CreateTempVarNameNode(*restTempVarName, tempVarType), assignCall)
			restTempDecl.SrcrefTreeIfMissing(objectPattern)
			restTempDecl.InsertBefore(nodeToDetach)

			restName := child.GetOnlyChild() // e.g. get `rest` from `const {...rest} = {};`
			if strings.HasPrefix(restName.GetString(), Es6RewriteDestructuring_DESTRUCTURING_TEMP_VAR) {
				jutil.CheckStateMsg1(restName.IsName(), "%s", restName)
				newLHS = e.CreateTempVarNameNode(restName.GetString(), AstFactory_Type_(restName))
			} else {
				newLHS = restName.Detach()
			}
			newRHS = e.ObjectPatternRestRHS(objectPattern, child, *restTempVarName, propsToDeleteForRest)
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
		newNode.SrcrefTreeIfMissing(child)

		newNode.InsertBefore(nodeToDetach)

		// Explicitly visit the LHS of the new node since it may be a nested
		// destructuring pattern.
		e.Visit(t, newLHS, newLHS.GetParent())
	}

	nodeToDetach.Detach()
	t.ReportCodeChange()
}

// ObjectPatternRestRHS converts "rest" of object destructuring lhs by making a
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
func (e *Es6RewriteDestructuring) ObjectPatternRestRHS(
	objectPattern rhino.Node, rest rhino.Node, restTempVarName string, statedProperties *jutil.ArrayList[rhino.Node]) rhino.Node {
	jutil.CheckArgument(objectPattern.GetLastChild() == rest)
	restTempVarModel := e.CreateTempVarNameNode(restTempVarName, AstFactory_Type_(objectPattern))
	result := restTempVarModel.CloneNode()
	if !statedProperties.IsEmpty() {
		props := statedProperties.Slice()
		comma := e.DeletionNodeForRestProperty(restTempVarModel.CloneNode(), props[0])
		for _, prop := range props[1:] {
			comma =
				e.astFactory.CreateComma(
					comma, e.DeletionNodeForRestProperty(restTempVarModel.CloneNode(), prop))
		}
		result = e.astFactory.CreateComma(comma, result)
	}
	result.SrcrefTreeIfMissing(rest)
	return result
}

// Java: Es6RewriteDestructuring.deletionNodeForRestProperty(Node, Node)
func (e *Es6RewriteDestructuring) DeletionNodeForRestProperty(restTempVarNameNode rhino.Node, property rhino.Node) rhino.Node {
	var get rhino.Node
	switch property.GetToken() {
	case rhino.Token_STRING_KEY:
		if property.IsQuotedStringKey() {
			get = e.astFactory.CreateGetElem(
				restTempVarNameNode, e.astFactory.CreateString(property.GetString()))
		} else {
			get = e.astFactory.CreateGetPropWithUnknownType(
				restTempVarNameNode, property.GetString())
		}
	case rhino.Token_NAME:
		get = e.astFactory.CreateGetElem(restTempVarNameNode, property)
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected property to delete node: " + property.ToStringTree()))
	}

	return e.astFactory.CreateDelProp(get)
}

// ReplaceArrayPattern converts
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
func (e *Es6RewriteDestructuring) ReplaceArrayPattern(
	t *NodeTraversal, arrayPattern rhino.Node, rhs rhino.Node, parent rhino.Node, nodeToDetach rhino.Node) {

	if e.rewriteMode == Es6RewriteDestructuring_ObjectDestructuringRewriteMode_REWRITE_OBJECT_REST {
		if len(e.patternNestingStack) == 0 || !e.patternNestingStackPeekLast().hasNestedObjectRest {
			return
		}
	}

	tempVarName := e.GetTempVariableName()
	makeIteratorCall := e.astFactory.CreateJSCompMakeIteratorCall(rhs.Detach(), e.namespace)
	tempDecl := e.astFactory.CreateSingleVarNameDeclarationWithValue(tempVarName, makeIteratorCall)
	tempVarModel := tempDecl.GetFirstChild()
	tempDecl.SrcrefTreeIfMissing(arrayPattern)
	tempDecl.InsertBefore(nodeToDetach)

	var next rhino.Node
	for child := arrayPattern.GetFirstChild(); !child.IsNil(); child = next {
		next = child.GetNext()
		if child.IsEmpty() {
			// Just call the next() method to advance the iterator, but throw away the value.
			nextCall :=
				rhino.IR_ExprResult(
					e.astFactory.CreateCallWithUnknownType(
						e.astFactory.CreateGetProp(
							tempVarModel.CloneNode(), "next", AstFactory_Type_OfColor(colors.StandardColors_TOP_OBJECT))))
			nextCall.SrcrefTreeIfMissing(child)
			nextCall.InsertBefore(nodeToDetach)
			continue
		}

		var newLHS rhino.Node
		var newRHS rhino.Node
		if child.IsDefaultValue() {
			//   [x = defaultValue] = rhs;
			// becomes
			//   var temp0 = $jscomp.makeIterator(rhs);
			//   var temp1 = temp.next().value
			//   x = (temp1 === undefined) ? defaultValue : temp1;
			nextVarName := e.GetTempVariableName()
			// `temp.next().value`
			nextCallDotValue :=
				e.astFactory.CreateGetPropWithUnknownType(
					e.astFactory.CreateCallWithUnknownType(
						e.astFactory.CreateGetPropWithUnknownType(tempVarModel.CloneNode(), "next")),
					"value")
			nextVarType := AstFactory_Type_(nextCallDotValue)
			// `var temp1 = temp.next().value`
			var_ := rhino.IR_VarWithValue(e.CreateTempVarNameNode(nextVarName, nextVarType), nextCallDotValue)
			var_.SrcrefTreeIfMissing(child)
			var_.InsertBefore(nodeToDetach)

			// `x`
			newLHS = child.RemoveFirstChild()
			// `(temp1 === undefined) ? defaultValue : temp1;
			newRHS =
				e.DefaultValueHook(
					e.CreateTempVarNameNode(nextVarName, nextVarType), child.GetLastChild().Detach())
		} else if child.IsRest() {
			//   [...x] = rhs;
			// becomes
			//   var temp = $jscomp.makeIterator(rhs);
			//   x = $jscomp.arrayFromIterator(temp);
			newLHS = child.RemoveFirstChild()
			newRHS = e.astFactory.CreateJscompArrayFromIteratorCall(tempVarModel.CloneNode(), e.namespace)
		} else {
			// LHS is just a name (or a nested pattern).
			//   var [x] = rhs;
			// becomes
			//   var temp = $jscomp.makeIterator(rhs);
			//   var x = temp.next().value;
			newLHS = child.Detach()
			newRHS =
				e.astFactory.CreateGetProp(
					e.astFactory.CreateCallWithUnknownType(
						e.astFactory.CreateGetPropWithUnknownType(tempVarModel.CloneNode(), "next")),
					"value",
					AstFactory_Type_(child))
		}
		var newNode rhino.Node
		if parent.IsAssign() {
			assignment := e.astFactory.CreateAssign(newLHS, newRHS)
			newNode = rhino.IR_ExprResult(assignment)
		} else {
			newNode = rhino.IR_DeclarationWithValue(newLHS, newRHS, parent.GetToken())
		}
		newNode.SrcrefTreeIfMissing(arrayPattern)

		newNode.InsertBefore(nodeToDetach)
		// Explicitly visit the LHS of the new node since it may be a nested
		// destructuring pattern.
		e.Visit(t, newLHS, newLHS.GetParent())
	}

	nodeToDetach.Detach()
	t.ReportCodeChange()
}

// WrapAssignOrDestructuringInCallToArrow replaces ASSIGN or
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
func (e *Es6RewriteDestructuring) WrapAssignOrDestructuringInCallToArrow(t *NodeTraversal, assignment rhino.Node) {
	lhs := assignment.GetFirstChild()
	rhs := assignment.GetLastChild().Detach()
	// NOTE: we do not presently support await/yield in the LHS of nested destructuring assignments.
	// See b/475296868.
	e.ValidateNoAwaitOrYieldInNestedAssignmentPattern(lhs)

	tempVarName := e.GetTempVariableName()

	tempVarModel := e.CreateTempVarNameNode(tempVarName, AstFactory_Type_(rhs))
	//  ((temp0) => {...})
	paramList := rhino.IR_ParamList(tempVarModel.CloneNode())
	// [x, y] = temp0;
	replacementExpr :=
		e.astFactory.CreateAssign(assignment.RemoveFirstChild(), tempVarModel.CloneNode())
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
	call.PutBooleanProp(rhino.Node_FREE_CALL, true)
	assignment.ReplaceWith(call)
	NodeUtil_MarkNewScopesChanged(call, e.compiler)

	e.ReplacePattern(
		t,
		replacementExpr.GetFirstChild(),
		replacementExpr.GetLastChild(),
		replacementExpr,
		exprResult)
}

// Java: Es6RewriteDestructuring.validateNoAwaitOrYieldInNestedAssignmentPattern(Node)
func (e *Es6RewriteDestructuring) ValidateNoAwaitOrYieldInNestedAssignmentPattern(pattern rhino.Node) {
	hasAwaitOrYieldInLhs :=
		!NodeUtil_FindPreorder(
			pattern, func(n rhino.Node) bool { return n.IsAwait() || n.IsYield() }, NodeUtil_MATCH_NOT_FUNCTION).IsNil()
	jutil.CheckStateMsg1(
		!hasAwaitOrYieldInLhs,
		"Cannot transpile yet: destructuring assignment referencing await or yield in lhs, in"+
			" nested sub-expression: %s",
		pattern)
}

// VisitDestructuringPatternInVanillaForInnerVars handles
// for (let [a, b, c] = arr; a < b; a++)
//
// Java: Es6RewriteDestructuring.visitDestructuringPatternInVanillaForInnerVars(NodeTraversal, Node)
func (e *Es6RewriteDestructuring) VisitDestructuringPatternInVanillaForInnerVars(t *NodeTraversal, pattern rhino.Node) {
	jutil.CheckArgument(pattern.IsDestructuringPattern())
	destructuringLhs := pattern.GetParent()
	declaration := destructuringLhs.GetParent()

	insertionPoint := declaration.GetParent()
	for insertionPoint.GetParent().IsLabel() {
		insertionPoint = insertionPoint.GetParent()
	}

	switch declaration.GetToken() {
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
			for c := declaration.GetFirstChild(); c != destructuringLhs; c = declaration.GetFirstChild() {
				newDeclaration := declaration.CloneNode()
				newDeclaration.AddChildToBack(c.Detach())
				newDeclaration.InsertBefore(insertionPoint)
			}

			// Move the pattern out of the initializer and transpile it
			newDeclaration := declaration.CloneNode()
			newDeclaration.AddChildToBack(destructuringLhs.Detach())
			newDeclaration.InsertBefore(insertionPoint)
			e.ReplacePattern(t, pattern, pattern.GetNext(), newDeclaration, newDeclaration)

			if !declaration.HasChildren() {
				declaration.ReplaceWith(rhino.IR_Empty())
			}
		}

	case rhino.Token_LET:
		// See https://tc39.es/ecma262/#sec-createperiterationenvironment
		// for (let a, b, c, unusedTmp = (() => [a, b, c] = arr); a < b; a++)
		{
			parsing.ParsingUtil_GetParamOrPatternNames(
				pattern, func(name rhino.Node) { name.CloneNode().InsertBefore(destructuringLhs) })

			unusedVar :=
				e.astFactory.
					CreateNameWithUnknownType(e.GetTempVariableName() + "$unused").
					Srcref(destructuringLhs)
			destructuringLhs.ReplaceWith(unusedVar)
			unusedVar.AddChildToBack(destructuringLhs)
			e.WrapAssignOrDestructuringInCallToArrow(t, destructuringLhs)
		}

	default:
		panic(jutil.NewIllegalStateExceptionMsg(declaration.String()))
	}

	e.compiler.ReportChangeToEnclosingScope(insertionPoint)
}

// VisitDestructuringPatternInEnhancedForInnerVars handles
// for (const [a, b, c] of arr)
//
// Java: Es6RewriteDestructuring.visitDestructuringPatternInEnhancedForInnerVars(Node)
func (e *Es6RewriteDestructuring) VisitDestructuringPatternInEnhancedForInnerVars(pattern rhino.Node) {
	jutil.CheckArgument(pattern.IsDestructuringPattern())
	tempVarName := e.GetTempVariableName()

	destructuringLhs := pattern.GetParent()
	jutil.CheckState(destructuringLhs.IsDestructuringLhs())
	declarationNode := destructuringLhs.GetParent()
	forNode := declarationNode.GetParent()
	jutil.CheckState(NodeUtil_IsEnhancedFor(forNode))
	block := forNode.GetLastChild()

	destructuringLhs.ReplaceWith(e.CreateTempVarNameNode(tempVarName, AstFactory_Type_(pattern)).Srcref(pattern))
	declarationType := declarationNode.GetToken()
	decl :=
		rhino.IR_DeclarationWithValue(
			pattern.Detach(), e.CreateTempVarNameNode(tempVarName, AstFactory_Type_(pattern)), declarationType)
	decl.SrcrefTreeIfMissing(pattern)
	// Move the body into an inner block to handle cases where declared variables in the for
	// loop initializer are shadowed by variables in the for loop body. e.g.
	//   for (const [value] of []) { const value = 1; }
	newBlock := rhino.IR_BlockOfNode(decl)
	block.ReplaceWith(newBlock)
	newBlock.AddChildToBack(block)
}

// VisitDestructuringPatternInEnhancedForWithOuterVars handles
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
// Java: Es6RewriteDestructuring.visitDestructuringPatternInEnhancedForWithOuterVars(Node)
func (e *Es6RewriteDestructuring) VisitDestructuringPatternInEnhancedForWithOuterVars(pattern rhino.Node) {
	jutil.CheckArgument(pattern.IsDestructuringPattern())
	tempVarName := e.GetTempVariableName()

	forNode := pattern.GetParent()
	block := forNode.GetLastChild()

	name := e.CreateTempVarNameNode(tempVarName, AstFactory_Type_(pattern))
	decl := rhino.IR_Var(name)
	decl.SrcrefTreeIfMissing(pattern)
	decl.InsertBefore(forNode)
	clonedName := name.CloneNode()
	clonedName.SrcrefTreeIfMissing(pattern)
	pattern.ReplaceWith(clonedName)
	exprResult :=
		rhino.IR_ExprResult(
			e.astFactory.CreateAssign(pattern, e.CreateTempVarNameNode(tempVarName, AstFactory_Type_(pattern))))
	exprResult.SrcrefTreeIfMissing(pattern)
	block.AddChildToFront(exprResult)
}

// Java: Es6RewriteDestructuring.visitDestructuringPatternInCatch(NodeTraversal, Node)
func (e *Es6RewriteDestructuring) VisitDestructuringPatternInCatch(t *NodeTraversal, pattern rhino.Node) {
	tempVarName := e.GetTempVariableName()
	catchBlock := pattern.GetNext()
	patternType := AstFactory_Type_(pattern)

	pattern.ReplaceWith(e.CreateTempVarNameNode(tempVarName, patternType))
	catchBlock.AddChildToFront(
		rhino.IR_DeclarationWithValue(pattern, e.CreateTempVarNameNode(tempVarName, patternType), rhino.Token_LET))
	NodeUtil_AddFeatureToScript(t.GetCurrentScript(), parser.FeatureSet_Feature_LET_DECLARATIONS, e.compiler)
}

// DefaultValueHook is a helper for transpiling DEFAULT_VALUE trees.
//
// Java: Es6RewriteDestructuring.defaultValueHook(Node, Node)
func (e *Es6RewriteDestructuring) DefaultValueHook(getprop rhino.Node, defaultValue rhino.Node) rhino.Node {
	undefined := e.astFactory.CreateUndefinedValue()
	undefined.MakeNonIndexable()
	getpropClone := getprop.CloneTree().SetColor(getprop.GetColor())
	return e.astFactory.CreateHook(
		e.astFactory.CreateSheq(getprop, undefined), defaultValue, getpropClone)
}
