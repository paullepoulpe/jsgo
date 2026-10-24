// Port of com/google/javascript/jscomp/ReplaceCssNames.java @ closure-compiler v20261021

/*
 * Copyright 2009 The Closure Compiler Authors.
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
	"fmt"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ReplaceCssNames replaces occurrences of goog.getCssName('foo') with a
// shorter version from the passed in renaming map. There are two styles of
// operation: for 'BY_WHOLE' we look up the whole string in the renaming map.
// For 'BY_PART', all the class name's components, separated by '-', are
// renamed individually and then recombined.
//
// Given the renaming map:
//
//	{
//	  once:  'a',
//	  upon:  'b',
//	  atime: 'c',
//	  long:  'd',
//	  time:  'e',
//	  ago:   'f'
//	}
//
// The following outputs are expected with the 'BY_PART' renaming style:
//
//	goog.getCssName('once') -> 'a'
//	goog.getCssName('once-upon-atime') -> 'a-b-c'
//
//	var baseClass = goog.getCssName('long-time');
//	el.className = goog.getCssName(baseClass, 'ago');
//	->
//	var baseClass = 'd-e';
//	el.className = baseClass + '-f';
//
// However if we have the following renaming map with the 'BY_WHOLE' renaming
// style:
//
//	{
//	  once: 'a',
//	  upon-atime: 'b',
//	  long-time: 'c',
//	  ago: 'd'
//	}
//
// Then we would expect:
//
//	goog.getCssName('once') -> 'a'
//
//	var baseClass = goog.getCssName('long-time');
//	el.className = goog.getCssName(baseClass, 'ago');
//	->
//	var baseClass = 'c';
//	el.className = baseClass + '-d';
//
// In addition, the CSS names before replacement can optionally be gathered.
//
// Java: class ReplaceCssNames implements CompilerPass
type ReplaceCssNames struct {
	compiler   AbstractCompiler
	astFactory *AstFactory

	cssNameCollector ReplaceCssNames_CssNameCollector

	cssNamesBySymbol *jutil.LinkedHashMap[string, string]

	classesObjectsQualifiedNames *jutil.LinkedHashSet[string]

	symbolMap CssRenamingMap // may be nil

	skiplist jutil.Set[string] // may be nil

	partsScratch    []string
	replacedScratch []string
	instanceSlab    jutil.Slab[ReplaceCssNames_GetCssNameInstance]
}

// This is used only for qualified name comparison, so it's OK to build it with IR instead of
// AstFactory. Type information isn't important here.
//
// Java: ReplaceCssNames.GET_CSS_NAME_FUNCTION
var ReplaceCssNames_GET_CSS_NAME_FUNCTION = rhino.IR_Getprop(rhino.IR_Name("goog"), "getCssName")

// Java: ReplaceCssNames.INVALID_NUM_ARGUMENTS_ERROR
var ReplaceCssNames_INVALID_NUM_ARGUMENTS_ERROR = DiagnosticType_Error(
	"JSC_GETCSSNAME_NUM_ARGS",
	"goog.getCssName called with \"{0}\" arguments, expected 1 or 2.")

// Java: ReplaceCssNames.STRING_LITERAL_EXPECTED_ERROR
var ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR = DiagnosticType_Error(
	"JSC_GETCSSNAME_STRING_LITERAL_EXPECTED",
	"goog.getCssName called with invalid argument, string literal "+
		"expected.  Was \"{0}\".")

// Java: ReplaceCssNames.UNEXPECTED_STRING_LITERAL_ERROR
var ReplaceCssNames_UNEXPECTED_STRING_LITERAL_ERROR = DiagnosticType_Error(
	"JSC_GETCSSNAME_UNEXPECTED_STRING_LITERAL",
	"goog.getCssName called with invalid arguments, string literal "+
		"passed as first of two arguments.  Did you mean "+
		"goog.getCssName(\"{0}-{1}\")?")

// Java: ReplaceCssNames.NESTED_CALL_ERROR
var ReplaceCssNames_NESTED_CALL_ERROR = DiagnosticType_Error(
	"JSC_GETCSSNAME_NESTED_CALL", "goog.getCssName: nested call is not allowed.")

// Java: ReplaceCssNames.UNKNOWN_SYMBOL_WARNING
var ReplaceCssNames_UNKNOWN_SYMBOL_WARNING = DiagnosticType_Warning(
	"JSC_GETCSSNAME_UNKNOWN_CSS_SYMBOL",
	"goog.getCssName called with unrecognized symbol \"{0}\" in class "+"\"{1}\".")

// Java: ReplaceCssNames.UNEXPECTED_SASS_GENERATED_CSS_TS_ERROR
var ReplaceCssNames_UNEXPECTED_SASS_GENERATED_CSS_TS_ERROR = DiagnosticType_Error(
	"JSC_UNEXPECTED_SASS_GENERATED_CSS_TS",
	"@sass_generated_css_ts JSDoc annotation is only allowed on .css.closure.js files.")

// Java: ReplaceCssNames.UNKNOWN_SYMBOL_ERROR
var ReplaceCssNames_UNKNOWN_SYMBOL_ERROR = DiagnosticType_Error(
	"JSC_UNKNOWN_CSS_SYMBOL_IN_CLASSES_OBJECT",
	"Symbol was not defined \"{0}\" in classes object.")

// Java: ReplaceCssNames.INVALID_USE_OF_CLASSES_OBJECT_ERROR
var ReplaceCssNames_INVALID_USE_OF_CLASSES_OBJECT_ERROR = DiagnosticType_Error(
	"JSC_INVALID_USE_OF_CLASSES_OBJECT",
	"invalid use of generated classes object. Only accessing its members is allowed.")

// ReplaceCssNames_CssNameCollector is called for each class name seen when
// replacing.
//
// Java: public interface ReplaceCssNames.CssNameCollector (a @FunctionalInterface)
type ReplaceCssNames_CssNameCollector interface {
	// Java: ReplaceCssNames.CssNameCollector.add(String)
	Add(className string)
}

// ReplaceCssNames_CssNameCollectorFunc lets a Go closure stand for a Java
// lambda implementing CssNameCollector.
//
// Go-only: language-adapter
type ReplaceCssNames_CssNameCollectorFunc func(className string)

// Add calls f.
func (f ReplaceCssNames_CssNameCollectorFunc) Add(className string) { f(className) }

// Java: ReplaceCssNames(AbstractCompiler, CssRenamingMap, ReplaceCssNames.CssNameCollector, Set<String>)
func NewReplaceCssNames(
	compiler AbstractCompiler,
	symbolMap CssRenamingMap,
	cssNameCollector ReplaceCssNames_CssNameCollector,
	skiplist jutil.Set[string]) *ReplaceCssNames {
	r := &ReplaceCssNames{}
	r.compiler = compiler
	r.astFactory = compiler.CreateAstFactory()
	r.symbolMap = symbolMap
	r.cssNameCollector = cssNameCollector
	r.skiplist = skiplist
	r.cssNamesBySymbol = jutil.NewLinkedHashMap[string, string]()
	r.classesObjectsQualifiedNames = jutil.NewLinkedHashSet[string]()
	return r
}

// Java: ReplaceCssNames.process(Node, Node)
func (r *ReplaceCssNames) Process(externs *rhino.Node, root *rhino.Node) {
	NodeTraversal_Traverse(r.compiler, root, &ReplaceCssNames_FindSetCssNameTraversal{outer: r})
	getCssNameInstances := r.GatherGetCssNameInstances(root)
	for getCssNameInstance := range getCssNameInstances.All() {
		cssClassName := getCssNameInstance.GetCssClassName()
		if getCssNameInstance.IsCssClosureFileClassesMember() {
			r.cssNamesBySymbol.Put(*getCssNameInstance.GetCssClosureClassesMemberName(), *cssClassName)
			r.classesObjectsQualifiedNames.Add(*getCssNameInstance.GetCssClosureClassesQualifiedName())
		} else {
			r.cssNameCollector.Add(*cssClassName)
		}
		getCssNameInstance.ReplaceWithExpression()
	}
	if !r.classesObjectsQualifiedNames.IsEmpty() {
		NodeTraversal_Traverse(r.compiler, root, &ReplaceCssNames_CountCssNamesBySymbol{outer: r})
	}
}

// GatherGetCssNameInstances creates a list of objects representing all the
// `goog.getCssName()` calls in the AST.
//
// The elements in this list will be in depth-first left-to-right order. So,
// if you perform actions on them in this order you'll affect nested calls
// before the containing calls. e.g.
//
//	goog.getCssName(goog.getCssName('me-first'), 'me-second')
//
// Java: ReplaceCssNames.gatherGetCssNameInstances(Node)
func (r *ReplaceCssNames) GatherGetCssNameInstances(root *rhino.Node) *jutil.ArrayList[*ReplaceCssNames_GetCssNameInstance] {
	gatherCssNamesTraversal := NewReplaceCssNames_GatherCssNamesTraversal(r)
	NodeTraversal_Traverse(r.compiler, root, gatherCssNamesTraversal)
	return gatherCssNamesTraversal.listOfCssNameInstances
}

// This is used only for qualified name comparison, so it's OK to build it with IR instead of
// AstFactory. Type information isn't important here.
var replaceCssNames_GOOG_SET_CSS_NAME_MAPPING = rhino.IR_Getprop(rhino.IR_Name("goog"), "setCssNameMapping")

// Java: private class ReplaceCssNames.FindSetCssNameTraversal extends AbstractPostOrderCallback
type ReplaceCssNames_FindSetCssNameTraversal struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer *ReplaceCssNames
}

// Java: ReplaceCssNames.FindSetCssNameTraversal.visit(NodeTraversal, Node, Node)
func (f *ReplaceCssNames_FindSetCssNameTraversal) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	r := f.outer
	if !n.IsCall() || !parent.IsExprResult() {
		return
	}
	callee := n.GetFirstChild()
	if !callee.MatchesQualifiedNameNode(replaceCssNames_GOOG_SET_CSS_NAME_MAPPING) {
		return
	}
	cssRenamingMap := ProcessClosurePrimitives_ProcessSetCssNameMapping(r.compiler, n, parent)
	r.symbolMap = cssRenamingMap
	r.compiler.ReportChangeToEnclosingScope(parent)
	parent.Detach()
}

// Java: private class ReplaceCssNames.GatherCssNamesTraversal implements NodeTraversal.Callback
type ReplaceCssNames_GatherCssNamesTraversal struct {
	outer                  *ReplaceCssNames
	listOfCssNameInstances *jutil.ArrayList[*ReplaceCssNames_GetCssNameInstance]
	traversalState         ReplaceCssNames_GatherCssNamesTraversal_TraversalState
}

// Java: ReplaceCssNames.GatherCssNamesTraversal()
func NewReplaceCssNames_GatherCssNamesTraversal(outer *ReplaceCssNames) *ReplaceCssNames_GatherCssNamesTraversal {
	return &ReplaceCssNames_GatherCssNamesTraversal{
		outer:                  outer,
		listOfCssNameInstances: jutil.NewArrayList[*ReplaceCssNames_GetCssNameInstance](),
	}
}

// Java: ReplaceCssNames.GatherCssNamesTraversal.shouldTraverse(NodeTraversal, Node, Node)
func (g *ReplaceCssNames_GatherCssNamesTraversal) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	r := g.outer
	jsdoc := n.GetJSDocInfo()
	hasSassGeneratedCssTsJsDoc := jsdoc != nil && jsdoc.IsSassGeneratedCssTs()
	if hasSassGeneratedCssTsJsDoc && (!n.IsScript() || r.IsNotInCssClosureFile(n)) {
		r.compiler.Report(JSError_MakeWithNode(n, ReplaceCssNames_UNEXPECTED_SASS_GENERATED_CSS_TS_ERROR))
		// Skip all nodes that are descendants of this one, since we found an invalid
		// @sassGeneratedCssTs JSDoc annotation, and are in an invalid state.
		return false
	}
	if n.IsScript() {
		g.traversalState.inSassGeneratedCssTsScript = hasSassGeneratedCssTsJsDoc
		g.traversalState.cssClosureClassesQualifiedName = nil
	} else if g.traversalState.inSassGeneratedCssTsScript && n.IsAssign() {
		assignmentTarget := n.GetFirstChild()
		if assignmentTarget.IsGetProp() && assignmentTarget.GetString() == "classes" {
			g.traversalState.cssClosureClassesQualifiedName = assignmentTarget.GetQualifiedName()
		}
	}
	return true
}

// Java: ReplaceCssNames.GatherCssNamesTraversal.visit(NodeTraversal, Node, Node)
func (g *ReplaceCssNames_GatherCssNamesTraversal) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	r := g.outer
	if r.IsGetCssNameCall(n) {
		getCssNameInstance := r.CreateGetCssNameInstance(n, g.traversalState.cssClosureClassesQualifiedName)
		if getCssNameInstance.IsValid() {
			g.listOfCssNameInstances.Add(getCssNameInstance)
		} else {
			r.compiler.Report(getCssNameInstance.GetValidationError())
		}
	}
}

// Java: ReplaceCssNames.GatherCssNamesTraversal.createSassGeneratedCssTsExpert(Node)
func (g *ReplaceCssNames_GatherCssNamesTraversal) CreateSassGeneratedCssTsExpert(n *rhino.Node) *ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert {
	r := g.outer
	jsdoc := n.GetJSDocInfo()
	hasSassGeneratedCssTsJsDoc := jsdoc != nil && jsdoc.IsSassGeneratedCssTs()

	var sassGeneratedCssTsValidationError *JSError
	if hasSassGeneratedCssTsJsDoc && (!n.IsScript() || r.IsNotInCssClosureFile(n)) {
		sassGeneratedCssTsValidationError = JSError_MakeWithNode(n, ReplaceCssNames_UNEXPECTED_SASS_GENERATED_CSS_TS_ERROR)
	}

	assignmentTarget := n.GetFirstChild()
	isCssClosureClassesAssignment := n.IsAssign() &&
		assignmentTarget.IsGetProp() &&
		assignmentTarget.GetString() == "classes"

	var cssClosureClassesQualifiedName *string
	if isCssClosureClassesAssignment {
		cssClosureClassesQualifiedName = assignmentTarget.GetQualifiedName()
	}

	return NewReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert(
		hasSassGeneratedCssTsJsDoc,
		sassGeneratedCssTsValidationError,
		isCssClosureClassesAssignment,
		cssClosureClassesQualifiedName)
}

// Java: private static class ReplaceCssNames.GatherCssNamesTraversal.SassGeneratedCssTsExpert
type ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert struct {
	hasSassGeneratedCssTsJsDoc        bool
	sassGeneratedCssTsValidationError *JSError
	isCssClosureClassesAssignment     bool
	cssClosureClassesQualifiedName    *string
}

// Java: ReplaceCssNames.GatherCssNamesTraversal.SassGeneratedCssTsExpert(boolean, JSError, boolean, String)
func NewReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert(
	hasSassGeneratedCssTsJsDoc bool,
	sassGeneratedCssTsValidationError *JSError,
	isCssClosureClassesAssignment bool,
	cssClosureClassesQualifiedName *string) *ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert {
	return &ReplaceCssNames_GatherCssNamesTraversal_SassGeneratedCssTsExpert{
		hasSassGeneratedCssTsJsDoc:        hasSassGeneratedCssTsJsDoc,
		sassGeneratedCssTsValidationError: sassGeneratedCssTsValidationError,
		isCssClosureClassesAssignment:     isCssClosureClassesAssignment,
		cssClosureClassesQualifiedName:    cssClosureClassesQualifiedName,
	}
}

// Java: private static class ReplaceCssNames.GatherCssNamesTraversal.TraversalState
type ReplaceCssNames_GatherCssNamesTraversal_TraversalState struct {
	inSassGeneratedCssTsScript     bool
	cssClosureClassesQualifiedName *string
}

// Java: ReplaceCssNames.isNotInCssClosureFile(Node)
func (r *ReplaceCssNames) IsNotInCssClosureFile(node *rhino.Node) bool {
	sourceFileName := node.GetSourceFileName()
	return sourceFileName == nil || !strings.HasSuffix(*sourceFileName, ".css.closure.js")
}

// Java: ReplaceCssNames.isGetCssNameCall(Node)
func (r *ReplaceCssNames) IsGetCssNameCall(node *rhino.Node) bool {
	return node.IsCall() && node.GetFirstChild().MatchesQualifiedNameNode(ReplaceCssNames_GET_CSS_NAME_FUNCTION)
}

// Java: ReplaceCssNames.createGetCssNameInstance(Node, String)
func (r *ReplaceCssNames) CreateGetCssNameInstance(
	n *rhino.Node, cssClosureClassesQualifiedName *string) *ReplaceCssNames_GetCssNameInstance {
	validationError := r.ValidateGetCssNameCall(n)
	inst := r.instanceSlab.New()
	inst.outer = r
	inst.callNode = n
	inst.validationError = validationError
	inst.cssClosureClassesQualifiedName = cssClosureClassesQualifiedName
	return inst
}

// ReplaceCssNames_GetCssNameInstance represents a `goog.getCssName()` call.
//
// Java: private class ReplaceCssNames.GetCssNameInstance
type ReplaceCssNames_GetCssNameInstance struct {
	outer *ReplaceCssNames

	// The CALL node representing `goog.getCssName()`.
	//
	// We shouldn't store final pointers to any of the children of this node, because they might
	// be changed between creation of this object and performance of other methods.
	callNode *rhino.Node

	// Non-null if the shape of the function call was invalid, when this object was created.
	validationError *JSError

	// The qualified name for the classes object if we're in a Sass-generated .css.ts file, or null
	// otherwise. For example "module$exports$foo.classes".
	cssClosureClassesQualifiedName *string
}

// Java: ReplaceCssNames.GetCssNameInstance(Node, JSError, String)
func NewReplaceCssNames_GetCssNameInstance(
	outer *ReplaceCssNames, callNode *rhino.Node, validationError *JSError, cssClosureClassesQualifiedName *string) *ReplaceCssNames_GetCssNameInstance {
	return &ReplaceCssNames_GetCssNameInstance{
		outer:                          outer,
		callNode:                       callNode,
		validationError:                validationError,
		cssClosureClassesQualifiedName: cssClosureClassesQualifiedName,
	}
}

// Java: ReplaceCssNames.GetCssNameInstance.isValid()
func (g *ReplaceCssNames_GetCssNameInstance) IsValid() bool {
	return g.validationError == nil
}

// Java: ReplaceCssNames.GetCssNameInstance.getValidationError()
func (g *ReplaceCssNames_GetCssNameInstance) GetValidationError() *JSError {
	jutil.CheckNotNullMsg1(g.validationError, "No validation error found: %s", g.callNode)
	return g.validationError
}

// ReplaceWithExpression replaces this `goog.getCssName()` call in the AST
// with a string typed expression.
//
// Java: ReplaceCssNames.GetCssNameInstance.replaceWithExpression()
func (g *ReplaceCssNames_GetCssNameInstance) ReplaceWithExpression() {
	jutil.CheckStateMsg1(g.IsValid(), "not a valid goog.getCssName() call: %s", g.callNode)
	jutil.CheckNotNullMsg1(g.callNode.GetParent(), "already replaced: %s", g.callNode)
	childCount := g.callNode.GetChildCount()
	switch childCount {
	case 2:
		g.ReplaceWithConvertedSingleArg()
	case 3:
		g.ReplaceWithConcatenatedArgs()
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			fmt.Sprintf("invalid number of children: %s for: %s", strconv.Itoa(childCount), g.callNode.String())))
	}
}

// Java: ReplaceCssNames.GetCssNameInstance.replaceWithConvertedSingleArg()
func (g *ReplaceCssNames_GetCssNameInstance) ReplaceWithConvertedSingleArg() {
	r := g.outer
	// `goog.getCssName('some-literal')` -> `'mapped-literal'`
	stringLitArg := jutil.CheckNotNullMsg1(g.callNode.GetSecondChild(), "%s", g.callNode)
	jutil.CheckStateMsg1(stringLitArg.IsStringLit(), "not a string literal: %s", stringLitArg)
	stringLitArg.Detach()
	r.ProcessStringNode(stringLitArg)
	g.callNode.ReplaceWith(stringLitArg)
	r.compiler.ReportChangeToEnclosingScope(stringLitArg)
}

// Java: ReplaceCssNames.GetCssNameInstance.replaceWithConcatenatedArgs()
func (g *ReplaceCssNames_GetCssNameInstance) ReplaceWithConcatenatedArgs() {
	r := g.outer
	// `goog.getCssName(someExpr, 'some-literal')` -> `someExpr + '-mapped-literal'`
	firstArg := jutil.CheckNotNullMsg1(g.callNode.GetSecondChild(), "%s", g.callNode)
	secondArg := jutil.CheckNotNullMsg1(firstArg.GetNext(), "%s", firstArg)
	firstArg.Detach()
	secondArg.Detach()
	r.ProcessStringNode(secondArg)
	secondArg.SetString("-" + secondArg.GetString())
	replacement := r.astFactory.CreateAdd(firstArg, secondArg).Srcref(g.callNode)
	g.callNode.ReplaceWith(replacement)
	r.compiler.ReportChangeToEnclosingScope(replacement)
}

// Java: ReplaceCssNames.GetCssNameInstance.isCssClosureFileClassesMember()
func (g *ReplaceCssNames_GetCssNameInstance) IsCssClosureFileClassesMember() bool {
	return g.cssClosureClassesQualifiedName != nil
}

// Java: ReplaceCssNames.GetCssNameInstance.getCssClosureClassesQualifiedName()
func (g *ReplaceCssNames_GetCssNameInstance) GetCssClosureClassesQualifiedName() *string {
	jutil.CheckNotNullMsg1(
		g.cssClosureClassesQualifiedName, "Not a css closure file classes object: %s", g.callNode)
	return g.cssClosureClassesQualifiedName
}

// Java: ReplaceCssNames.GetCssNameInstance.getCssClosureClassesMemberName()
func (g *ReplaceCssNames_GetCssNameInstance) GetCssClosureClassesMemberName() *string {
	jutil.CheckNotNullMsg1(
		g.cssClosureClassesQualifiedName, "Not a css closure file classes object: %s", g.callNode)
	// We support 2 structures getCssName() calls in Sass-generated files:
	// 1. { Foo: goog.getCssName('foo') }
	// 2. { Foo: expression ? … : goog.getCssName('foo') }
	// Traversing our ancestry allows us to find the `Foo` string key in both cases.
	ancestors := g.callNode.GetAncestors()
	for ancestor := range ancestors.All() {
		if ancestor.IsStringKey() {
			name := *g.cssClosureClassesQualifiedName + "." + ancestor.GetString()
			return &name
		}
	}
	panic(jutil.NewIllegalStateExceptionMsg(
		fmt.Sprintf("No css closure classes member name found: %s", g.callNode.String())))
}

// Java: ReplaceCssNames.GetCssNameInstance.getCssClassName()
func (g *ReplaceCssNames_GetCssNameInstance) GetCssClassName() *string {
	jutil.CheckStateMsg1(g.IsValid(), "not a valid goog.getCssName() call: %s", g.callNode)
	jutil.CheckNotNullMsg1(g.callNode.GetParent(), "already replaced: %s", g.callNode)
	childCount := g.callNode.GetChildCount()
	firstArg := jutil.CheckNotNullMsg1(g.callNode.GetSecondChild(), "%s", g.callNode)
	switch childCount {
	case 2:
		jutil.CheckStateMsg1(firstArg.IsStringLit(), "not a string literal: %s", firstArg)
		name := firstArg.GetString()
		return &name
	case 3:
		secondArg := jutil.CheckNotNullMsg1(firstArg.GetNext(), "%s", firstArg)
		jutil.CheckStateMsg1(secondArg.IsStringLit(), "not a string literal: %s", secondArg)
		name := secondArg.GetString()
		return &name
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			fmt.Sprintf("invalid number of children: %s for: %s", strconv.Itoa(childCount), g.callNode.String())))
	}
}

// Java: ReplaceCssNames.validateGetCssNameCall(Node)
func (r *ReplaceCssNames) ValidateGetCssNameCall(callNode *rhino.Node) *JSError {
	childCount := callNode.GetChildCount()
	switch childCount {
	case 2:
		return r.ValidateSingleArgGetCssNameCall(callNode)
	case 3:
		return r.ValidateTwoArgGetCssNameCall(callNode)
	default:
		return JSError_MakeWithNode(callNode, ReplaceCssNames_INVALID_NUM_ARGUMENTS_ERROR, strconv.Itoa(childCount))
	}
}

// Java: ReplaceCssNames.validateSingleArgGetCssNameCall(Node)
func (r *ReplaceCssNames) ValidateSingleArgGetCssNameCall(callNode *rhino.Node) *JSError {
	// `goog.getCssName('css-name')`
	stringLiteralArg := jutil.CheckNotNull(callNode.GetLastChild())
	if !stringLiteralArg.IsStringLit() {
		return JSError_MakeWithNode(
			callNode, ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR, stringLiteralArg.GetToken().String())
	}
	return nil
}

// Java: ReplaceCssNames.validateTwoArgGetCssNameCall(Node)
func (r *ReplaceCssNames) ValidateTwoArgGetCssNameCall(callNode *rhino.Node) *JSError {
	// `goog.getCssName(BASE_CSS_NAME, 'css-suffix')`
	// The first child is the callee, so the first argument is the second child
	firstArg := jutil.CheckNotNull(callNode.GetSecondChild())
	secondArg := jutil.CheckNotNull(firstArg.GetNext())

	if !secondArg.IsStringLit() {
		return JSError_MakeWithNode(callNode, ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR, secondArg.GetToken().String())
	}
	if firstArg.IsStringLit() {
		return JSError_MakeWithNode(
			callNode, ReplaceCssNames_UNEXPECTED_STRING_LITERAL_ERROR, firstArg.GetString(), secondArg.GetString())
	}
	if r.IsGetCssNameCall(firstArg) {
		// Disallow `goog.getCssName(goog.getCssName('n1'), 'n2')`.
		// Disallowing this is a bit arbitrary and is done for historical reasons.
		// The 2-argument form results in more JS code shipped to the browser and
		// forces users to have camelCase CSS names, which contravenes normal CSS style.
		// Since we'd like to stop supporting the 2-argument form, we don't want to
		// start allowing any form of it that was previously an error.
		return JSError_MakeWithNode(firstArg, ReplaceCssNames_NESTED_CALL_ERROR)
	}

	return nil
}

// ProcessStringNode processes a string argument to goog.getCssName(). The
// string will be renamed based off the symbol map. If there is no map or any
// part of the name can't be renamed, a warning is reported to the compiler
// and the node is left unchanged.
//
// If the type is unexpected then an error is reported to the compiler.
//
// n is the string node to process.
//
// Java: ReplaceCssNames.processStringNode(Node)
func (r *ReplaceCssNames) ProcessStringNode(n *rhino.Node) {
	name := n.GetString()
	if r.skiplist != nil && r.skiplist.Contains(name) {
		// We apply the skiplist before splitting on dashes, and not after.
		// External substitution maps should do the same.
		return
	}
	if r.symbolMap != nil {
		var replacement *string

		if strings.HasPrefix(name, "--") {
			// Force BY_WHOLE style for CSS variables.
			replacement = r.symbolMap.Get(name)
			if replacement == nil {
				r.compiler.Report(JSError_MakeWithNode(n, ReplaceCssNames_UNKNOWN_SYMBOL_WARNING, name, name))
				return
			}
		} else {
			switch r.symbolMap.GetStyle() {
			case CssRenamingMap_Style_BY_WHOLE:
				replacement = r.symbolMap.Get(name)
				if replacement == nil {
					r.compiler.Report(JSError_MakeWithNode(n, ReplaceCssNames_UNKNOWN_SYMBOL_WARNING, name, name))
					return
				}
			case CssRenamingMap_Style_BY_PART:
				parts := r.partsScratch[:0]
				if !strings.Contains(name, "-") {
					parts = append(parts, name)
				} else {
					for rem := name; ; {
						idx := strings.IndexByte(rem, '-')
						if idx < 0 {
							parts = append(parts, rem)
							break
						}
						parts = append(parts, rem[:idx])
						rem = rem[idx+1:]
					}
					for len(parts) > 0 && parts[len(parts)-1] == "" {
						parts = parts[:len(parts)-1]
					}
				}
				r.partsScratch = parts
				replaced := r.replacedScratch[:0]
				for i := 0; i < len(parts); i++ {
					part := r.symbolMap.Get(parts[i])
					if part == nil {
						// If we can't encode all parts, don't encode any of it.
						r.compiler.Report(JSError_MakeWithNode(n, ReplaceCssNames_UNKNOWN_SYMBOL_WARNING, parts[i], name))
						return
					}
					replaced = append(replaced, *part)
				}
				r.replacedScratch = replaced
				joined := strings.Join(replaced, "-")
				replacement = &joined
			}
		}
		// n.setString(replacement): a null replacement throws.
		n.SetString(*replacement)
	}
}

// Java: private class ReplaceCssNames.CountCssNamesBySymbol implements NodeTraversal.Callback
type ReplaceCssNames_CountCssNamesBySymbol struct {
	outer *ReplaceCssNames
}

// Java: ReplaceCssNames.CountCssNamesBySymbol.shouldTraverse(NodeTraversal, Node, Node)
func (c *ReplaceCssNames_CountCssNamesBySymbol) ShouldTraverse(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) bool {
	if n.IsScript() {
		// Only descend into source files NOT named `*.css.closure.js`
		return c.outer.IsNotInCssClosureFile(n)
	} else {
		return true // descend into every other node
	}
}

// Java: ReplaceCssNames.CountCssNamesBySymbol.visit(NodeTraversal, Node, Node)
func (c *ReplaceCssNames_CountCssNamesBySymbol) Visit(t *NodeTraversal, n *rhino.Node, parent *rhino.Node) {
	r := c.outer
	if r.classesObjectsQualifiedNames.IsEmpty() || !n.IsQualifiedName() {
		return
	}
	classesObjectQualifiedName := n.GetQualifiedName()
	if classesObjectQualifiedName == nil ||
		!r.classesObjectsQualifiedNames.Contains(*classesObjectQualifiedName) {
		return
	}
	classNameNode := n.GetParent()
	if !n.IsGetProp() || !classNameNode.IsGetProp() {
		r.compiler.Report(JSError_MakeWithNode(n, ReplaceCssNames_INVALID_USE_OF_CLASSES_OBJECT_ERROR))
		return
	}
	classQualifiedName := classNameNode.GetQualifiedName()
	if classQualifiedName == nil || !r.cssNamesBySymbol.ContainsKey(*classQualifiedName) {
		r.compiler.Report(JSError_MakeWithNode(n, ReplaceCssNames_UNKNOWN_SYMBOL_ERROR, classNameNode.GetString()))
		return
	}
	className := r.cssNamesBySymbol.Get(*classQualifiedName)
	r.cssNameCollector.Add(className)
}
