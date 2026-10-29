// Port of com/google/javascript/jscomp/ReplaceCssNames.java @ closure-compiler v20261021
// Port of com/google/javascript/jscomp/CssRenamingMap.java @ closure-compiler v20261021

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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CssRenamingMap is the interface used by ReplaceCssNames to substitute
// CSS class names.
type CssRenamingMap interface {
	RenamingMap

	// Java: CssRenamingMap.getStyle()
	GetStyle() CssRenamingMap_Style
}

// CssRenamingMap_Style is the kind of renaming map.
type CssRenamingMap_Style int32

const (
	CssRenamingMap_Style_BY_WHOLE CssRenamingMap_Style = iota
	CssRenamingMap_Style_BY_PART
)

var cssRenamingMap_Style_names = [...]string{"BY_WHOLE", "BY_PART"}

// Name is Enum.name().
func (s CssRenamingMap_Style) Name() string { return cssRenamingMap_Style_names[s] }

// String is Enum.toString().
func (s CssRenamingMap_Style) String() string { return s.Name() }

// Ordinal is Enum.ordinal().
func (s CssRenamingMap_Style) Ordinal() int { return int(s) }

// CssRenamingMap_Style_ValueOf is Style.valueOf(String).
func CssRenamingMap_Style_ValueOf(name string) CssRenamingMap_Style {
	for i, n := range cssRenamingMap_Style_names {
		if n == name {
			return CssRenamingMap_Style(i)
		}
	}
	panic(jutil.NewIllegalArgumentExceptionMsg(
		"No enum constant com.google.javascript.jscomp.CssRenamingMap.Style." + name))
}

// CssRenamingMap_ByPartBase provides GetStyle() returning BY_PART for ByPart maps.
type CssRenamingMap_ByPartBase struct{}

// Java: CssRenamingMap.ByPart.getStyle()
func (CssRenamingMap_ByPartBase) GetStyle() CssRenamingMap_Style {
	return CssRenamingMap_Style_BY_PART
}

// CssRenamingMap_ByWholeBase provides GetStyle() returning BY_WHOLE for ByWhole maps.
type CssRenamingMap_ByWholeBase struct{}

// Java: CssRenamingMap.ByWhole.getStyle()
func (CssRenamingMap_ByWholeBase) GetStyle() CssRenamingMap_Style {
	return CssRenamingMap_Style_BY_WHOLE
}

// ReplaceCssNames replaces occurrences of goog.getCssName('foo') with a
// shorter version from the passed in renaming map. There are two styles of
// operation: for 'BY_WHOLE' we look up the whole string in the renaming map.
// For 'BY_PART', all the class name's components, separated by '-', are
// renamed individually and then recombined.
//
// In addition, the CSS names before replacement can optionally be gathered.
//
// Java: class ReplaceCssNames implements CompilerPass
type ReplaceCssNames struct {
	compiler   AbstractCompiler
	astFactory *AstFactory

	cssNameCollector             ReplaceCssNames_CssNameCollector
	cssNamesBySymbol             map[string]string
	classesObjectsQualifiedNames map[string]struct{}

	symbolMap CssRenamingMap    // may be nil
	skiplist  jutil.Set[string] // may be nil
}

var (
	replaceCssNamesGetCssNameFn      = rhino.QualifiedName_Of("goog.getCssName")
	replaceCssNamesSetCssNameMapping = rhino.QualifiedName_Of("goog.setCssNameMapping")
)

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
	"goog.getCssName called with unrecognized symbol \"{0}\" in class \"{1}\".")

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
type ReplaceCssNames_CssNameCollector interface {
	Add(className string)
}

// ReplaceCssNames_CssNameCollectorFunc lets a Go closure stand for a Java
// lambda implementing CssNameCollector.
type ReplaceCssNames_CssNameCollectorFunc func(className string)

// Add calls f.
func (f ReplaceCssNames_CssNameCollectorFunc) Add(className string) { f(className) }

// Java: ReplaceCssNames(AbstractCompiler, CssRenamingMap, ReplaceCssNames.CssNameCollector, Set<String>)
func NewReplaceCssNames(
	compiler AbstractCompiler,
	symbolMap CssRenamingMap,
	cssNameCollector ReplaceCssNames_CssNameCollector,
	skiplist jutil.Set[string]) *ReplaceCssNames {
	return &ReplaceCssNames{
		compiler:         compiler,
		astFactory:       compiler.CreateAstFactory(),
		symbolMap:        symbolMap,
		cssNameCollector: cssNameCollector,
		skiplist:         skiplist,
	}
}

// Java: ReplaceCssNames.process(Node, Node)
func (r *ReplaceCssNames) Process(externs rhino.Node, root rhino.Node) {
	NodeTraversal_Traverse(r.compiler, root, &findSetCssNameTraversal{r: r})

	gather := gatherCssNamesTraversal{r: r}
	NodeTraversal_Traverse(r.compiler, root, &gather)
	for i := range gather.instances {
		inst := &gather.instances[i]
		cssClassName := inst.cssClassName()
		if inst.classesQName != "" {
			if r.cssNamesBySymbol == nil {
				r.cssNamesBySymbol = make(map[string]string)
				r.classesObjectsQualifiedNames = make(map[string]struct{})
			}
			r.cssNamesBySymbol[inst.classesMemberName()] = cssClassName
			r.classesObjectsQualifiedNames[inst.classesQName] = struct{}{}
		} else {
			r.cssNameCollector.Add(cssClassName)
		}
		r.replaceGetCssNameCall(inst.callNode)
	}
	if len(r.classesObjectsQualifiedNames) > 0 {
		NodeTraversal_Traverse(r.compiler, root, &countCssNamesBySymbol{r: r})
	}
}

var (
	replaceCssNamesCallVisitTokens          = makeVisitTokenSet(rhino.Token_CALL)
	replaceCssNamesQualifiedNameVisitTokens = makeVisitTokenSet(
		rhino.Token_NAME, rhino.Token_GETPROP, rhino.Token_THIS, rhino.Token_SUPER)
)

type findSetCssNameTraversal struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	r *ReplaceCssNames
}

var (
	_ nodeTraversal_noVisitLeaves = (*findSetCssNameTraversal)(nil)
	_ nodeTraversal_visitTokens   = (*findSetCssNameTraversal)(nil)
)

func (*findSetCssNameTraversal) noVisitLeaves()             {}
func (*findSetCssNameTraversal) visitTokens() *ast.TokenSet { return &replaceCssNamesCallVisitTokens }

func (f *findSetCssNameTraversal) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	call, ok := ast.AsCall(n)
	if !ok || parent.IsNil() || !parent.IsExprResult() || !replaceCssNamesSetCssNameMapping.Matches(call.Callee().Node()) {
		return
	}
	f.r.symbolMap = ProcessClosurePrimitives_ProcessSetCssNameMapping(f.r.compiler, n, parent)
	f.r.compiler.ReportChangeToEnclosingScope(parent)
	parent.Detach()
}

type gatherCssNamesTraversal struct {
	r                          *ReplaceCssNames
	instances                  []cssNameInstance
	inSassGeneratedCssTsScript bool
	classesQualifiedName       string
}

var (
	_ nodeTraversal_noVisitLeaves = (*gatherCssNamesTraversal)(nil)
	_ nodeTraversal_visitTokens   = (*gatherCssNamesTraversal)(nil)
)

func (*gatherCssNamesTraversal) noVisitLeaves()             {}
func (*gatherCssNamesTraversal) visitTokens() *ast.TokenSet { return &replaceCssNamesCallVisitTokens }

func (g *gatherCssNamesTraversal) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	jsdoc := n.GetJSDocInfo()
	hasSassGeneratedCssTsJsDoc := jsdoc != nil && jsdoc.IsSassGeneratedCssTs()
	if hasSassGeneratedCssTsJsDoc && (!n.IsScript() || isNotInCssClosureFile(n)) {
		g.r.compiler.Report(JSError_MakeWithNode(n, ReplaceCssNames_UNEXPECTED_SASS_GENERATED_CSS_TS_ERROR))
		// Skip all nodes that are descendants of this one, since we found an invalid
		// @sassGeneratedCssTs JSDoc annotation, and are in an invalid state.
		return false
	}
	if n.IsScript() {
		g.inSassGeneratedCssTsScript = hasSassGeneratedCssTsJsDoc
		g.classesQualifiedName = ""
	} else if g.inSassGeneratedCssTsScript {
		if assign, ok := ast.AsAssign(n); ok {
			if gp, ok := assign.Target().AsGetProp(); ok && gp.Property() == "classes" {
				if qname, ok := gp.Node().GetQualifiedNameString(); ok {
					g.classesQualifiedName = qname
				} else {
					g.classesQualifiedName = ""
				}
			}
		}
	}
	return true
}

func (g *gatherCssNamesTraversal) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	if !isGetCssNameCall(n) {
		return
	}
	if err := validateGetCssNameCall(n); err != nil {
		g.r.compiler.Report(err)
		return
	}
	// Collect valid goog.getCssName() calls in depth-first left-to-right post-order so
	// nested calls are processed before containing calls.
	g.instances = append(g.instances, cssNameInstance{
		callNode:     n,
		classesQName: g.classesQualifiedName,
	})
}

func isNotInCssClosureFile(node rhino.Node) bool {
	sourceFileName := node.GetSourceFileName()
	return sourceFileName == nil || !strings.HasSuffix(*sourceFileName, ".css.closure.js")
}

func isGetCssNameCall(node rhino.Node) bool {
	call, ok := ast.AsCall(node)
	return ok && replaceCssNamesGetCssNameFn.Matches(call.Callee().Node())
}

// cssNameInstance represents a validated `goog.getCssName()` call.
type cssNameInstance struct {
	callNode     rhino.Node
	classesQName string
}

func (inst *cssNameInstance) classesMemberName() string {
	// We support 2 structures of getCssName() calls in Sass-generated files:
	// 1. { Foo: goog.getCssName('foo') }
	// 2. { Foo: expression ? … : goog.getCssName('foo') }
	// Traversing our ancestry allows us to find the `Foo` string key in both cases.
	for ancestor := inst.callNode.GetParent(); !ancestor.IsNil(); ancestor = ancestor.GetParent() {
		if ancestor.IsStringKey() {
			return inst.classesQName + "." + ancestor.GetString()
		}
	}
	panic(jutil.NewIllegalStateExceptionMsg(
		"No css closure classes member name found: " + inst.callNode.String()))
}

func (inst *cssNameInstance) cssClassName() string {
	rhino.CheckNotNullMsg1(inst.callNode.GetParent(), "already replaced: %s", inst.callNode)
	childCount := inst.callNode.GetChildCount()
	firstArg := rhino.CheckNotNullMsg1(inst.callNode.GetSecondChild(), "%s", inst.callNode)
	switch childCount {
	case 2:
		jutil.CheckStateMsg1(firstArg.IsStringLit(), "not a string literal: %s", firstArg)
		return firstArg.GetString()
	case 3:
		secondArg := rhino.CheckNotNullMsg1(firstArg.GetNext(), "%s", firstArg)
		jutil.CheckStateMsg1(secondArg.IsStringLit(), "not a string literal: %s", secondArg)
		return secondArg.GetString()
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"invalid number of children: " + strconv.Itoa(childCount) + " for: " + inst.callNode.String()))
	}
}

func (r *ReplaceCssNames) replaceGetCssNameCall(callNode rhino.Node) {
	rhino.CheckNotNullMsg1(callNode.GetParent(), "already replaced: %s", callNode)
	childCount := callNode.GetChildCount()
	switch childCount {
	case 2:
		// `goog.getCssName('some-literal')` -> `'mapped-literal'`
		stringLitArg := rhino.CheckNotNullMsg1(callNode.GetSecondChild(), "%s", callNode)
		jutil.CheckStateMsg1(stringLitArg.IsStringLit(), "not a string literal: %s", stringLitArg)
		stringLitArg.Detach()
		r.processStringNode(stringLitArg)
		callNode.ReplaceWith(stringLitArg)
		r.compiler.ReportChangeToEnclosingScope(stringLitArg)
	case 3:
		// `goog.getCssName(someExpr, 'some-literal')` -> `someExpr + '-mapped-literal'`
		firstArg := rhino.CheckNotNullMsg1(callNode.GetSecondChild(), "%s", callNode)
		secondArg := rhino.CheckNotNullMsg1(firstArg.GetNext(), "%s", firstArg)
		firstArg.Detach()
		secondArg.Detach()
		r.processStringNode(secondArg)
		secondArg.SetString("-" + secondArg.GetString())
		replacement := r.astFactory.CreateAdd(firstArg, secondArg).Srcref(callNode)
		callNode.ReplaceWith(replacement)
		r.compiler.ReportChangeToEnclosingScope(replacement)
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"invalid number of children: " + strconv.Itoa(childCount) + " for: " + callNode.String()))
	}
}

func validateGetCssNameCall(callNode rhino.Node) *JSError {
	childCount := callNode.GetChildCount()
	switch childCount {
	case 2:
		// `goog.getCssName('css-name')`
		stringLiteralArg := rhino.CheckNotNull(callNode.GetLastChild())
		if !stringLiteralArg.IsStringLit() {
			return JSError_MakeWithNode(
				callNode, ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR, stringLiteralArg.GetToken().String())
		}
		return nil
	case 3:
		// `goog.getCssName(BASE_CSS_NAME, 'css-suffix')`
		firstArg := rhino.CheckNotNull(callNode.GetSecondChild())
		secondArg := rhino.CheckNotNull(firstArg.GetNext())
		if !secondArg.IsStringLit() {
			return JSError_MakeWithNode(callNode, ReplaceCssNames_STRING_LITERAL_EXPECTED_ERROR, secondArg.GetToken().String())
		}
		if firstArg.IsStringLit() {
			return JSError_MakeWithNode(
				callNode, ReplaceCssNames_UNEXPECTED_STRING_LITERAL_ERROR, firstArg.GetString(), secondArg.GetString())
		}
		if isGetCssNameCall(firstArg) {
			return JSError_MakeWithNode(firstArg, ReplaceCssNames_NESTED_CALL_ERROR)
		}
		return nil
	default:
		return JSError_MakeWithNode(callNode, ReplaceCssNames_INVALID_NUM_ARGUMENTS_ERROR, strconv.Itoa(childCount))
	}
}

func (r *ReplaceCssNames) processStringNode(n rhino.Node) {
	name := n.GetString()
	if r.skiplist != nil && r.skiplist.Contains(name) {
		// We apply the skiplist before splitting on dashes, and not after.
		// External substitution maps should do the same.
		return
	}
	if r.symbolMap == nil {
		return
	}

	var replacement string
	if strings.HasPrefix(name, "--") || r.symbolMap.GetStyle() == CssRenamingMap_Style_BY_WHOLE {
		// Force BY_WHOLE style for CSS variables.
		mapped, ok := r.symbolMap.Get(name)
		if !ok {
			r.compiler.Report(JSError_MakeWithNode(n, ReplaceCssNames_UNKNOWN_SYMBOL_WARNING, name, name))
			return
		}
		replacement = mapped
	} else if r.symbolMap.GetStyle() == CssRenamingMap_Style_BY_PART {
		if strings.IndexByte(name, '-') < 0 {
			part, ok := r.symbolMap.Get(name)
			if !ok {
				r.compiler.Report(JSError_MakeWithNode(n, ReplaceCssNames_UNKNOWN_SYMBOL_WARNING, name, name))
				return
			}
			replacement = part
		} else {
			trimmed := strings.TrimRight(name, "-")
			if trimmed != "" {
				var sb strings.Builder
				sb.Grow(len(trimmed))
				for rest := trimmed; ; {
					dash := strings.IndexByte(rest, '-')
					var p string
					if dash < 0 {
						p = rest
					} else {
						p = rest[:dash]
					}
					part, ok := r.symbolMap.Get(p)
					if !ok {
						// If we can't encode all parts, don't encode any of it.
						r.compiler.Report(JSError_MakeWithNode(n, ReplaceCssNames_UNKNOWN_SYMBOL_WARNING, p, name))
						return
					}
					sb.WriteString(part)
					if dash < 0 {
						break
					}
					sb.WriteByte('-')
					rest = rest[dash+1:]
				}
				replacement = sb.String()
			}
		}
	}
	n.SetString(replacement)
}

type countCssNamesBySymbol struct {
	r *ReplaceCssNames
}

var (
	_ nodeTraversal_visitTokens                         = (*countCssNamesBySymbol)(nil)
	_ nodeTraversal_onlyScriptAndFunctionShouldTraverse = (*countCssNamesBySymbol)(nil)
)

func (*countCssNamesBySymbol) visitTokens() *ast.TokenSet {
	return &replaceCssNamesQualifiedNameVisitTokens
}

func (*countCssNamesBySymbol) shouldTraverseOnlyScriptAndFunction() bool { return true }

func (c *countCssNamesBySymbol) ShouldTraverse(t *NodeTraversal, n rhino.Node, parent rhino.Node) bool {
	if n.IsScript() {
		// Only descend into source files NOT named `*.css.closure.js`
		return isNotInCssClosureFile(n)
	}
	return true
}

func (c *countCssNamesBySymbol) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {
	r := c.r
	if len(r.classesObjectsQualifiedNames) == 0 || !n.IsQualifiedName() {
		return
	}
	classesObjQName, ok := n.GetQualifiedNameString()
	if !ok {
		return
	}
	if _, found := r.classesObjectsQualifiedNames[classesObjQName]; !found {
		return
	}
	classNameNode := n.GetParent()
	if !n.IsGetProp() || classNameNode.IsNil() || !classNameNode.IsGetProp() {
		r.compiler.Report(JSError_MakeWithNode(n, ReplaceCssNames_INVALID_USE_OF_CLASSES_OBJECT_ERROR))
		return
	}
	classQName, ok := classNameNode.GetQualifiedNameString()
	if !ok {
		r.compiler.Report(JSError_MakeWithNode(n, ReplaceCssNames_UNKNOWN_SYMBOL_ERROR, classNameNode.GetString()))
		return
	}
	className, defined := r.cssNamesBySymbol[classQName]
	if !defined {
		r.compiler.Report(JSError_MakeWithNode(n, ReplaceCssNames_UNKNOWN_SYMBOL_ERROR, classNameNode.GetString()))
		return
	}
	r.cssNameCollector.Add(className)
}
