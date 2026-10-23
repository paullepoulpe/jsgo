// Port of com/google/javascript/jscomp/parsing/IRFactory.java @ closure-compiler v20261015

/*
 * Copyright 2013 The Closure Compiler Authors.
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

package parsing

import (
	"math"
	"math/big"
	"reflect"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/dtoa"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// IRFactory transforms the external AST to the internal AST.
//
// Java: class IRFactory
type IRFactory struct {
	sourceFile rhino.StaticSourceFile
	// Go-only: language-adapter — rhino.ErrorReporter takes a Go string, so
	// Java's null source name (no source file, in tests) is "".
	sourceName string

	// Source file reference that also contains the file content.
	fileWithContent *parser.SourceFile

	config              *Config
	errorReporter       rhino.ErrorReporter
	transformDispatcher *IRFactory_TransformDispatcher

	// If non-null, use this set of keywords instead of TokenStream.isKeyword().
	reservedKeywords *jutil.LinkedHashSet[string]

	parsedCommentBits    []uint64
	parsedCommentBitsBuf [4]uint64

	licenseBuilder    *jutil.LinkedHashSet[string]
	firstFileoverview *rhino.JSDocInfo

	// Use a template node for properties set on all nodes to minimize the
	// memory footprint associated with these.
	// templateNode acts as the current template node for the IRFactory. The
	// template node is updated from defaultTemplateNode to
	// closureUnawareTemplateNode whenever a new closure-unaware code range is
	// entered in the transform method and swapped back to the
	// defaultTemplateNode when the closure-unaware code range is exited.
	// This is safe because while on a closure unaware code range, only nodes
	// that should have the IS_IN_CLOSURE_UNAWARE_SUBTREE property are created.
	// An alternative was to directly update the templateNode's propListHead,
	// but we don't want to expose the propListHead outside of the Node class.
	templateNode               *rhino.Node
	defaultTemplateNode        *rhino.Node
	closureUnawareTemplateNode *rhino.Node
	defaultTemplateNodeVal     rhino.Node

	jsdocTracker       *IRFactory_CommentTracker
	nonJsdocTracker    *IRFactory_CommentTracker
	jsdocTrackerVal    IRFactory_CommentTracker
	nonJsdocTrackerVal IRFactory_CommentTracker
	jsDocSourceKind    JsDocInfoParser_JsDocSourceKind

	currentFileIsExterns bool

	features   *parser.FeatureSet
	resultNode *rhino.Node

	closureUnawareCodeRanges *jutil.TreeRangeSet[IRFactory_LineAndColumn]

	jsdocSlabsVal            rhino.JSDocBuilderSlabs
	jsdocSlabs               *rhino.JSDocBuilderSlabs
	jsdocStream              JsDocTokenStream
	jsdocParser              JsDocInfoParser
	transformDispatcherVal   IRFactory_TransformDispatcher
	classScopeVal            IRFactory_TransformDispatcher_ScopeTracker
	privateIdLhsOfInScopeVal IRFactory_TransformDispatcher_ScopeTracker
}

const IRFactory_GETTER_ERROR_MESSAGE = "getters are not supported in older versions of JavaScript. " +
	"If you are targeting newer versions of JavaScript, " +
	"set the appropriate language_in option."

const IRFactory_SETTER_ERROR_MESSAGE = "setters are not supported in older versions of JavaScript. " +
	"If you are targeting newer versions of JavaScript, " +
	"set the appropriate language_in option."

const IRFactory_INVALID_ES3_PROP_NAME = "Keywords and reserved words are not allowed as unquoted property " +
	"names in older versions of JavaScript. " +
	"If you are targeting newer versions of JavaScript, " +
	"set the appropriate language_in option."

const IRFactory_INVALID_ES5_STRICT_OCTAL = "Octal integer literals are not supported in strict mode."

const IRFactory_INVALID_OCTAL_DIGIT = "Invalid octal digit in octal literal."

const IRFactory_STRING_CONTINUATION_WARNING = "String continuations are not recommended. See" +
	" https://google.github.io/styleguide/jsguide.html#features-strings-no-line-continuations"

const IRFactory_OCTAL_STRING_LITERAL_WARNING = "Octal literals in strings are not supported in this language mode."

const IRFactory_DUPLICATE_PARAMETER = "Duplicate parameter name \"%s\""

const IRFactory_DUPLICATE_LABEL = "Duplicate label \"%s\""

const IRFactory_UNLABELED_BREAK = "unlabelled break must be inside loop or switch"

const IRFactory_UNEXPECTED_CONTINUE = "continue must be inside loop"

const IRFactory_UNEXPECTED_LABELLED_CONTINUE = "continue can only use labels of iteration statements"

const IRFactory_UNEXPECTED_RETURN = "return must be inside function"
const IRFactory_UNEXPECTED_YIELD = "yield must be inside generator function"
const IRFactory_UNEXPECTED_AWAIT = "await must be inside asynchronous function"
const IRFactory_UNEXPECTED_FOR_AWAIT_OF = "'for-await-of' used in a non-async function context"
const IRFactory_UNEXPECTED_NEW_DOT_TARGET = "new.target must be inside a function"
const IRFactory_UNDEFINED_LABEL = "undefined label \"%s\""

var IRFactory_ES5_RESERVED_KEYWORDS = jutil.NewLinkedHashSet(
	// From Section 7.6.1.2
	"class", "const", "enum", "export", "extends", "import", "super")

var IRFactory_ES5_STRICT_RESERVED_KEYWORDS = jutil.NewLinkedHashSet(
	// From Section 7.6.1.2
	"class",
	"const",
	"enum",
	"export",
	"extends",
	"import",
	"super",
	"implements",
	"interface",
	"let",
	"package",
	"private",
	"protected",
	"public",
	"static",
	"yield")

// IRFactory_LineAndColumn represents a line and column number in a file.
//
// Java: record IRFactory.LineAndColumn(int lineNo, int column) implements Comparable<LineAndColumn>
type IRFactory_LineAndColumn struct {
	lineNo int
	column int
}

// Java: IRFactory.LineAndColumn(int, int)
func NewIRFactory_LineAndColumn(lineNo int, column int) IRFactory_LineAndColumn {
	return IRFactory_LineAndColumn{lineNo: lineNo, column: column}
}

// Java: IRFactory.LineAndColumn.lineNo()
func (l IRFactory_LineAndColumn) LineNo() int { return l.lineNo }

// Java: IRFactory.LineAndColumn.column()
func (l IRFactory_LineAndColumn) Column() int { return l.column }

// Java: IRFactory.LineAndColumn.fromSourcePosition(SourcePosition)
func IRFactory_LineAndColumn_FromSourcePosition(pos parser.SourcePosition) IRFactory_LineAndColumn {
	return NewIRFactory_LineAndColumn(int(pos.Line), int(pos.Column))
}

// Java: IRFactory.LineAndColumn.compareTo(LineAndColumn)
func (l IRFactory_LineAndColumn) CompareTo(o IRFactory_LineAndColumn) int {
	if c := irFactory_compareInt(l.lineNo, o.lineNo); c != 0 {
		return c
	}
	return irFactory_compareInt(l.column, o.column)
}

// irFactory_compareInt is Integer.compare(int, int), the comparison of
// Comparator.comparingInt.
//
// Go-only: language-adapter — JDK method.
func irFactory_compareInt(x, y int) int {
	if x < y {
		return -1
	}
	if x == y {
		return 0
	}
	return 1
}

// Equals is the record's implicit equals(Object).
func (l IRFactory_LineAndColumn) Equals(o any) bool {
	other, ok := o.(IRFactory_LineAndColumn)
	return ok && other == l
}

// HashCode is the record's implicit hashCode().
func (l IRFactory_LineAndColumn) HashCode() int32 {
	return 31*int32(l.lineNo) + int32(l.column)
}

// String is the record's implicit toString().
func (l IRFactory_LineAndColumn) String() string {
	return "LineAndColumn[lineNo=" + jutil.ToJavaString(l.lineNo) + ", column=" + jutil.ToJavaString(l.column) + "]"
}

// irFactory_LineAndColumn_compare is the natural ordering of LineAndColumn
// (Comparable.compareTo) as a comparator.
//
// Go-only: language-adapter — TreeRangeSet<LineAndColumn> orders by
// Comparable; the jutil TreeRangeSet takes the comparator.
func irFactory_LineAndColumn_compare(a, b IRFactory_LineAndColumn) int { return a.CompareTo(b) }

// Java: IRFactory(StaticSourceFile, Config, ErrorReporter, ImmutableList<Comment>, SourceFile, JsDocInfoParser.JsDocSourceKind)
func NewIRFactory(
	sourceFile rhino.StaticSourceFile,
	config *Config,
	errorReporter rhino.ErrorReporter,
	comments []*parser.Comment,
	fileWithContent *parser.SourceFile,
	jsDocSourceKind JsDocInfoParser_JsDocSourceKind) *IRFactory {
	return NewIRFactoryWithSlabs(sourceFile, config, errorReporter, comments, fileWithContent, jsDocSourceKind, nil)
}

// NewIRFactoryWithSlabs constructs an IRFactory using optional external JSDoc builder slabs.
func NewIRFactoryWithSlabs(
	sourceFile rhino.StaticSourceFile,
	config *Config,
	errorReporter rhino.ErrorReporter,
	comments []*parser.Comment,
	fileWithContent *parser.SourceFile,
	jsDocSourceKind JsDocInfoParser_JsDocSourceKind,
	jsdocSlabs *rhino.JSDocBuilderSlabs) *IRFactory {
	f := &IRFactory{}
	if jsdocSlabs != nil {
		f.jsdocSlabs = jsdocSlabs
	} else {
		f.jsdocSlabs = &f.jsdocSlabsVal
	}
	if nWords := (len(comments) + 63) / 64; nWords <= len(f.parsedCommentBitsBuf) {
		f.parsedCommentBits = f.parsedCommentBitsBuf[:nWords]
	} else {
		f.parsedCommentBits = make([]uint64, nWords)
	}
	f.firstFileoverview = nil
	f.currentFileIsExterns = false
	f.features = parser.FeatureSet_BARE_MINIMUM

	f.jsdocTrackerVal = IRFactory_CommentTracker{index: -1, previousIndex: -1, source: comments, filter: irFactory_isJsdocComment}
	f.jsdocTrackerVal.Advance()
	f.jsdocTracker = &f.jsdocTrackerVal

	f.nonJsdocTrackerVal = IRFactory_CommentTracker{index: -1, previousIndex: -1, source: comments, filter: irFactory_isNonJsdocComment}
	f.nonJsdocTrackerVal.Advance()
	f.nonJsdocTracker = &f.nonJsdocTrackerVal

	f.sourceFile = sourceFile
	f.fileWithContent = fileWithContent
	// The template node properties are applied to all nodes in this transform.
	f.defaultTemplateNode = f.CreateTemplateNode()
	f.templateNode = f.defaultTemplateNode

	// Sometimes this will be null in tests.
	if sourceFile == nil {
		f.sourceName = ""
	} else {
		f.sourceName = sourceFile.GetName()
	}

	f.config = config
	f.errorReporter = errorReporter
	f.transformDispatcherVal = IRFactory_TransformDispatcher{
		IRFactory:             f,
		classScope:            &f.classScopeVal,
		privateIdLhsOfInScope: &f.privateIdLhsOfInScopeVal,
	}
	f.transformDispatcher = &f.transformDispatcherVal

	if config.StrictMode().IsStrict() {
		f.reservedKeywords = IRFactory_ES5_STRICT_RESERVED_KEYWORDS
	} else if config.LanguageMode() == Config_LanguageMode_ECMASCRIPT3 {
		f.reservedKeywords = nil // use TokenStream.isKeyword instead
	} else {
		f.reservedKeywords = IRFactory_ES5_RESERVED_KEYWORDS
	}
	f.jsDocSourceKind = jsDocSourceKind
	return f
}

func (f *IRFactory) markCommentParsed(idx int) {
	if uint(idx) < uint(len(f.parsedCommentBits)*64) {
		f.parsedCommentBits[idx>>6] |= 1 << (uint(idx) & 63)
	}
}

func (f *IRFactory) markCommentPtrParsed(comment *parser.Comment) {
	if idx := f.jsdocTracker.previousIndex; idx >= 0 && idx < len(f.jsdocTracker.source) && f.jsdocTracker.source[idx] == comment {
		f.markCommentParsed(idx)
		return
	}
	for i, c := range f.jsdocTracker.source {
		if c == comment {
			f.markCommentParsed(i)
			return
		}
	}
}

func (f *IRFactory) isCommentParsed(idx int) bool {
	return uint(idx) < uint(len(f.parsedCommentBits)*64) && (f.parsedCommentBits[idx>>6]&(1<<(uint(idx)&63))) != 0
}

// The CommentTracker filters of the IRFactory constructor (the lambdas
// (c) -> c.type == Comment.Type.JSDOC and (c) -> c.type != Comment.Type.JSDOC).
func irFactory_isJsdocComment(c *parser.Comment) bool    { return c.Type == parser.Comment_Type_JSDOC }
func irFactory_isNonJsdocComment(c *parser.Comment) bool { return c.Type != parser.Comment_Type_JSDOC }

// Java: final class IRFactory.CommentTracker
type IRFactory_CommentTracker struct {
	source        []*parser.Comment
	filter        func(*parser.Comment) bool
	index         int
	previousIndex int
}

// Java: IRFactory.CommentTracker(ImmutableList<Comment>, Predicate<Comment>)
func NewIRFactory_CommentTracker(source []*parser.Comment, filter func(*parser.Comment) bool) *IRFactory_CommentTracker {
	t := &IRFactory_CommentTracker{}
	t.index = -1
	t.previousIndex = -1
	t.source = source
	t.filter = filter

	t.Advance()
	return t
}

// Java: IRFactory.CommentTracker.current()
func (t *IRFactory_CommentTracker) Current() *parser.Comment {
	if t.index >= len(t.source) {
		return nil
	}
	return t.source[t.index]
}

// Java: IRFactory.CommentTracker.advance()
func (t *IRFactory_CommentTracker) Advance() {
	t.previousIndex = t.index
	for {
		t.index++ // Always advance at least one element.

		c := t.Current()
		if c == nil || t.filter(c) {
			break
		}
	}
}

// Java: IRFactory.CommentTracker.backtrack()
func (t *IRFactory_CommentTracker) Backtrack() {
	t.index = t.previousIndex
}

// Java: IRFactory.CommentTracker.hasPendingCommentBefore(SourcePosition)
func (t *IRFactory_CommentTracker) HasPendingCommentBefore(pos parser.SourcePosition) bool {
	c := t.Current()
	// The line number matters for trailing comments
	return c != nil && c.Location.End.Line <= pos.Line && c.Location.End.Offset <= pos.Offset
}

// CreateTemplateNode creates a template node to use as a source of common
// attributes, this allows the prop structure to be shared among all the node
// from this source file. This reduces the cost of these properties to
// O(nodes) to O(files).
//
// Java: IRFactory.createTemplateNode()
func (f *IRFactory) CreateTemplateNode() *rhino.Node {
	// The Node type choice is arbitrary.
	templateNode := f.defaultTemplateNodeVal.Init(rhino.Token_SCRIPT)
	templateNode.SetStaticSourceFile(f.sourceFile)
	return templateNode
}

// Java: IRFactory.createClosureUnawareTemplateNode()
func (f *IRFactory) CreateClosureUnawareTemplateNode() *rhino.Node {
	// We clone the default template node, so that the source file property,
	// and other potential properties are shared
	closureUnawareTemplateNode := f.defaultTemplateNode.CloneNode()

	// We add the IS_IN_CLOSURE_UNAWARE_SUBTREE property to the head of the
	// property list.
	closureUnawareTemplateNode.SetIsInClosureUnawareSubtree(true)
	return closureUnawareTemplateNode
}

// Java: IRFactory.transformTree(ProgramTree, StaticSourceFile, Config, ErrorReporter, SourceFile)
func IRFactory_TransformTree(
	tree *parser.ProgramTree,
	sourceFile rhino.StaticSourceFile,
	config *Config,
	errorReporter rhino.ErrorReporter,
	file *parser.SourceFile) *IRFactory {
	return IRFactory_TransformTreeWithSlabs(tree, sourceFile, config, errorReporter, file, nil)
}

// IRFactory_TransformTreeWithSlabs transforms tree into an IR AST using optional external JSDoc builder slabs.
func IRFactory_TransformTreeWithSlabs(
	tree *parser.ProgramTree,
	sourceFile rhino.StaticSourceFile,
	config *Config,
	errorReporter rhino.ErrorReporter,
	file *parser.SourceFile,
	jsdocSlabs *rhino.JSDocBuilderSlabs) *IRFactory {
	jsDocSourceKind := JsDocInfoParser_JsDocSourceKind_NORMAL
	if sourceFile.IsTypeScriptSource() {
		jsDocSourceKind = JsDocInfoParser_JsDocSourceKind_TSICKLE
	}
	irFactory := NewIRFactoryWithSlabs(
		sourceFile, config, errorReporter, tree.SourceComments, file, jsDocSourceKind, jsdocSlabs)

	// don't call transform as we don't want standard jsdoc handling.
	n := irFactory.transformDispatcher.Process(tree)
	irFactory.SetSourceInfo(n, tree)

	if tree.SourceComments != nil {
		for i, comment := range tree.SourceComments {
			if (comment.Type == parser.Comment_Type_JSDOC || comment.Type == parser.Comment_Type_IMPORTANT) &&
				!irFactory.isCommentParsed(i) {
				useLicensesOnlyConfig := irFactory.WithinClosureUnawareCodeRange(
					int(comment.Location.Start.Line), int(comment.Location.Start.Column))
				irFactory.markCommentParsed(i)
				irFactory.HandlePossibleFileOverviewJsDocComment(comment, useLicensesOnlyConfig)
			}
		}

		if n.IsScript() {
			endOfFilePos := tree.Location.End
			// Handle end of file comments that are still pending
			nonJSDocComment := irFactory.ParseNonJSDocCommentAt(endOfFilePos /* isInline= */, false)
			if nonJSDocComment != nil {
				// pending end-of-file comment exists && parsing mode is set to
				// INCLUDE_ALL_COMMENTS
				n.SetTrailingNonJSDocComment(nonJSDocComment)
			}
		}
	}

	irFactory.SetFileOverviewJsDoc(n)

	irFactory.ValidateAll(n)
	irFactory.resultNode = n

	return irFactory
}

// Java: IRFactory.getResultNode()
func (f *IRFactory) GetResultNode() *rhino.Node {
	return f.resultNode
}

// Java: IRFactory.getFeatures()
func (f *IRFactory) GetFeatures() *parser.FeatureSet {
	return f.features
}

// Java: IRFactory.validateAll(Node)
func (f *IRFactory) ValidateAll(n *rhino.Node) {
	// Go-only: language-adapter — the ArrayDeque used as a stack (push/poll at
	// the head) is a slice used as a stack.
	var workBuf [32]*rhino.Node
	work := workBuf[:0]
	for n != nil {
		f.Validate(n)
		nextSibling := n.GetNext()
		firstChild := n.GetFirstChild()
		if firstChild != nil {
			if nextSibling != nil {
				// handle the siblings later
				work = append(work, nextSibling)
			}
			n = firstChild
		} else if nextSibling != nil {
			// no children, handle the next sibling
			n = nextSibling
		} else {
			// no siblings, continue with work we have saved on the work queue
			if len(work) == 0 {
				n = nil
			} else {
				n = work[len(work)-1]
				work = work[:len(work)-1]
			}
		}
	}
	jutil.CheckState(len(work) == 0)
}

// Java: IRFactory.validate(Node)
func (f *IRFactory) Validate(n *rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_PARAM_LIST:
		f.ValidateParameters(n)
	case rhino.Token_BREAK, rhino.Token_CONTINUE:
		f.ValidateBreakContinue(n)
	case rhino.Token_RETURN:
		f.ValidateReturn(n)
	case rhino.Token_YIELD:
		f.ValidateYield(n)
	case rhino.Token_AWAIT, rhino.Token_FOR_AWAIT_OF:
		f.ValidateAwait(n)
	case rhino.Token_NEW_TARGET:
		f.ValidateNewDotTarget(n)
	case rhino.Token_LABEL:
		f.ValidateLabel(n)
	case rhino.Token_FUNCTION:
		f.ValidateBlockScopedFunctions(n)
	case rhino.Token_NAME:
		f.ValidateShadowHost(n)
	}
}

// Java: IRFactory.validateShadowHost(Node)
func (f *IRFactory) ValidateShadowHost(n *rhino.Node) {
	if n.GetClosureUnawareShadow() == nil {
		return
	}

	if f.firstFileoverview == nil || !f.firstFileoverview.IsClosureUnawareCode() {
		f.Error(n, "invalid @closureUnaware node: must be in a file annotated with @closureUnaware")
		return
	}

	parent := n.GetParent()
	// We only support:

	//  - FUNCTION nodes that are direct CALL callees
	if parent.IsCall() {
		return
	}
	//  - FUNCTION nodes that are children of a GETPROP node that is a direct
	//    CALL callee (e.g. fn.call or fn.apply can be the CALL callee)
	grandparent := parent.GetParent()
	if parent.IsGetProp() && grandparent != nil && grandparent.IsCall() {
		calledFn := parent.GetString()
		if calledFn == "call" || calledFn == "apply" {
			return
		}
	}
	f.Error(n, "invalid @closureUnaware node: must be directly called")
}

// Java: IRFactory.validateAwait(Node)
func (f *IRFactory) ValidateAwait(n *rhino.Node) {
	var errorMsg string
	if n.IsAwait() {
		errorMsg = IRFactory_UNEXPECTED_AWAIT
	} else if n.IsForAwaitOf() {
		errorMsg = IRFactory_UNEXPECTED_FOR_AWAIT_OF
	} else {
		return
	}

	parent := n
	for {
		parent = parent.GetParent()
		if parent == nil {
			break
		}
		// The await is in a class static block or a field initializer.
		// e.g. `class C { static { await y; } }` or `class C { x = await y; }`
		if (parent.IsBlock() || parent.IsMemberFieldDef()) && parent.GetParent().IsClassMembers() {
			f.errorReporter.Error(errorMsg, f.sourceName, n.GetLineno(), n.GetCharno())
			return
		}
		if parent.IsAsyncFunction() {
			return
		}
		if parent.IsFunction() {
			// The await is in a non-async function.
			// e.g. `function f() { return await 5; }`
			// nested function e.g. `async function f() { function f2() { return await 5; } }`
			f.errorReporter.Error(errorMsg, f.sourceName, n.GetLineno(), n.GetCharno())
			return
		}
		if parent.IsScript() && parent.GetBooleanProp(rhino.Node_ES6_MODULE) {
			// This is a valid top-level await in an ES module.
			f.MaybeWarnForFeatureNode(n, parser.FeatureSet_Feature_TOP_LEVEL_AWAIT)
			return
		}
	}
	f.errorReporter.Error(errorMsg, f.sourceName, n.GetLineno(), n.GetCharno())
}

// Java: IRFactory.validateYield(Node)
func (f *IRFactory) ValidateYield(n *rhino.Node) {
	if !n.IsYield() {
		return
	}

	parent := n
	for {
		parent = parent.GetParent()
		if parent == nil {
			break
		}
		// The yield is in a class static block.
		// e.g. `class C { static { yield; } }`
		if parent.IsClassMembers() {
			f.errorReporter.Error(IRFactory_UNEXPECTED_YIELD, f.sourceName, n.GetLineno(), n.GetCharno())
			return
		}
		if parent.IsGeneratorFunction() {
			return
		}
	}
}

// Java: IRFactory.validateReturn(Node)
func (f *IRFactory) ValidateReturn(n *rhino.Node) {
	if !n.IsReturn() {
		return
	}

	parent := n
	for {
		parent = parent.GetParent()
		if parent == nil {
			break
		}
		// The return is in a class static block.
		// e.g. `class C { static { return; } }`
		if parent.IsClassMembers() {
			f.errorReporter.Error(IRFactory_UNEXPECTED_RETURN, f.sourceName, n.GetLineno(), n.GetCharno())
			return
		}
		if parent.IsFunction() {
			return
		}
	}
	f.errorReporter.Error(IRFactory_UNEXPECTED_RETURN, f.sourceName, n.GetLineno(), n.GetCharno())
}

// Java: IRFactory.validateNewDotTarget(Node)
func (f *IRFactory) ValidateNewDotTarget(n *rhino.Node) {
	if n.GetToken() != rhino.Token_NEW_TARGET {
		return
	}

	parent := n
	for {
		parent = parent.GetParent()
		if parent == nil {
			break
		}
		if parent.IsFunction() {
			return
		}
	}
	f.errorReporter.Error(IRFactory_UNEXPECTED_NEW_DOT_TARGET, f.sourceName, n.GetLineno(), n.GetCharno())
}

// Java: IRFactory.validateBreakContinue(Node)
func (f *IRFactory) ValidateBreakContinue(n *rhino.Node) {
	if !n.IsBreak() && !n.IsContinue() {
		return
	}

	labelName := n.GetFirstChild()
	if labelName != nil {
		parent := n.GetParent()
		for !parent.IsLabel() || !IRFactory_LabelsMatch(parent, labelName) {
			if parent.IsFunction() || parent.IsScript() || parent.IsClassMembers() {
				// report missing label
				f.errorReporter.Error(
					jutil.StringFormat(IRFactory_UNDEFINED_LABEL, labelName.GetString()),
					f.sourceName,
					n.GetLineno(),
					n.GetCharno())
				break
			}
			parent = parent.GetParent()
		}

		if parent.IsLabel() && IRFactory_LabelsMatch(parent, labelName) {
			if n.IsContinue() && !IRFactory_IsContinueTarget(parent.GetLastChild()) {
				// report invalid continue target
				f.errorReporter.Error(
					IRFactory_UNEXPECTED_LABELLED_CONTINUE, f.sourceName, n.GetLineno(), n.GetCharno())
			}
		}
	} else if n.IsContinue() {
		parent := n.GetParent()
		for !IRFactory_IsContinueTarget(parent) {
			if parent.IsFunction() || parent.IsScript() || parent.IsClassMembers() {
				// report invalid continue
				f.errorReporter.Error(IRFactory_UNEXPECTED_CONTINUE, f.sourceName, n.GetLineno(), n.GetCharno())
				break
			}
			parent = parent.GetParent()
		}
	} else {
		parent := n.GetParent()
		for !IRFactory_IsBreakTarget(parent) {
			if parent.IsFunction() || parent.IsScript() || parent.IsClassMembers() {
				// report invalid break
				f.errorReporter.Error(IRFactory_UNLABELED_BREAK, f.sourceName, n.GetLineno(), n.GetCharno())
				break
			}
			parent = parent.GetParent()
		}
	}
}

// Java: IRFactory.isBreakTarget(Node)
func IRFactory_IsBreakTarget(n *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF,
		rhino.Token_WHILE, rhino.Token_DO, rhino.Token_SWITCH:
		return true
	default:
		return false
	}
}

// Java: IRFactory.isContinueTarget(Node)
func IRFactory_IsContinueTarget(n *rhino.Node) bool {
	switch n.GetToken() {
	case rhino.Token_FOR, rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF,
		rhino.Token_WHILE, rhino.Token_DO:
		return true
	default:
		return false
	}
}

// Java: IRFactory.labelsMatch(Node, Node)
func IRFactory_LabelsMatch(label *rhino.Node, labelName *rhino.Node) bool {
	return label.GetFirstChild().GetString() == labelName.GetString()
}

// Java: IRFactory.validateLabel(Node)
func (f *IRFactory) ValidateLabel(n *rhino.Node) {
	if !n.IsLabel() {
		return
	}

	labelName := n.GetFirstChild()
	for parent := n.GetParent(); parent != nil && !parent.IsFunction(); parent = parent.GetParent() {
		if parent.IsLabel() && IRFactory_LabelsMatch(parent, labelName) {
			f.errorReporter.Error(
				jutil.StringFormat(IRFactory_DUPLICATE_LABEL, labelName.GetString()),
				f.sourceName,
				n.GetLineno(),
				n.GetCharno())
			break
		}
	}
}

// Java: IRFactory.validateParameters(Node)
func (f *IRFactory) ValidateParameters(n *rhino.Node) {
	if !n.IsParamList() {
		return
	}
	first := n.GetFirstChild()
	if first == nil {
		return
	}
	if first.IsName() {
		second := first.GetNext()
		if second == nil {
			return
		}
		var seen [8]string
		seen[0] = first.GetString()
		count := 1
		allSimple := true
		for c := second; c != nil; c = c.GetNext() {
			if !c.IsName() || count >= len(seen) {
				allSimple = false
				break
			}
			seen[count] = c.GetString()
			count++
		}
		if allSimple {
			hasDup := false
			for i := 1; i < count; i++ {
				for j := 0; j < i; j++ {
					if seen[i] == seen[j] {
						hasDup = true
						break
					}
				}
				if hasDup {
					break
				}
			}
			if !hasDup {
				return
			}
		}
	}
	f.validateParametersSlow(n)
}

func (f *IRFactory) validateParametersSlow(n *rhino.Node) {
	seenNames := jutil.NewLinkedHashSet[string]()
	for c := n.GetFirstChild(); c != nil; c = c.GetNext() {
		ParsingUtil_GetParamOrPatternNames(
			c,
			func(param *rhino.Node) {
				paramName := param.GetString()
				if !seenNames.Add(paramName) {
					f.errorReporter.Warning(
						jutil.StringFormat(IRFactory_DUPLICATE_PARAMETER, paramName),
						f.sourceName,
						param.GetLineno(),
						param.GetCharno())
				}
			})
	}
}

// Java: IRFactory.validateBlockScopedFunctions(Node)
func (f *IRFactory) ValidateBlockScopedFunctions(n *rhino.Node) {
	if n.IsFunction() && n.GetParent().IsBlock() && !n.GetGrandparent().IsFunction() {
		f.MaybeWarnForFeatureNode(n, parser.FeatureSet_Feature_BLOCK_SCOPED_FUNCTION_DECLARATION)
	}
}

// Java: IRFactory.setFileOverviewJsDoc(Node)
func (f *IRFactory) SetFileOverviewJsDoc(irNode *rhino.Node) {
	var fileoverview *rhino.JSDocInfo_Builder
	if f.firstFileoverview == nil {
		fileoverview = rhino.JSDocInfo_Builder_()
	} else {
		fileoverview = f.firstFileoverview.ToBuilder()
	}

	if f.licenseBuilder != nil && !f.licenseBuilder.IsEmpty() {
		var sb strings.Builder
		for license := range f.licenseBuilder.All() {
			sb.WriteString(license)
		}
		fileoverview.RecordLicense(sb.String())
	}

	irNode.SetJSDocInfo(fileoverview.BuildWithAlways(false))
}

// Java: IRFactory.transformBlock(ParseTree)
func (f *IRFactory) TransformBlock(node parser.ParseTree) *rhino.Node {
	irNode := f.Transform(node)
	if irNode.IsBlock() {
		return irNode
	}

	var newBlock *rhino.Node
	if irNode.IsEmpty() {
		newBlock = f.NewNode(rhino.Token_BLOCK)
	} else {
		newBlock = f.NewNodeWithChild(rhino.Token_BLOCK, irNode)
	}
	f.SetSourceInfoFromNode(newBlock, irNode)
	newBlock.SetIsAddedBlock(true)

	return newBlock
}

// HandlePossibleFileOverviewJsDoc returns true if the jsDocParser represents
// a fileoverview.
//
// Java: IRFactory.handlePossibleFileOverviewJsDoc(JsDocInfoParser)
func (f *IRFactory) HandlePossibleFileOverviewJsDoc(jsDocParser *JsDocInfoParser) bool {
	if jsDocParser.GetLicenseText() != nil {
		if f.licenseBuilder == nil {
			f.licenseBuilder = jutil.NewLinkedHashSet[string]()
		}
		f.licenseBuilder.Add(*jsDocParser.GetLicenseText())
	}

	newFileoverview := jsDocParser.GetFileOverviewJSDocInfo()
	if base.JSCompObjects_Identical(newFileoverview, f.firstFileoverview) {
		return false
	}

	if f.firstFileoverview == nil {
		f.firstFileoverview = newFileoverview
		f.currentFileIsExterns = newFileoverview.IsExterns()
	} else {
		merged := f.firstFileoverview.ToBuilder()
		if !newFileoverview.GetSuppressions().IsEmpty() {
			merged.RecordSuppressions(newFileoverview.GetSuppressions())
		}
		if newFileoverview.IsExterns() {
			f.currentFileIsExterns = true
			merged.RecordExterns()
		}
		if newFileoverview.IsNoCoverage() {
			merged.RecordNoCoverage()
		}
		if newFileoverview.IsTypeSummary() {
			merged.RecordTypeSummary()
		}
		f.firstFileoverview = merged.Build()
	}

	return true
}

// Java: IRFactory.handlePossibleFileOverviewJsDoc(Comment, boolean)
func (f *IRFactory) HandlePossibleFileOverviewJsDocComment(comment *parser.Comment, useLicensesOnlyConfig bool) {
	jsDocParser := f.CreateJsDocInfoParserWithLicensesOnlyConfig(comment, useLicensesOnlyConfig)
	f.markCommentPtrParsed(comment)
	f.HandlePossibleFileOverviewJsDoc(jsDocParser)
}

// Java: IRFactory.getJSDocCommentAt(SourcePosition)
func (f *IRFactory) GetJSDocCommentAt(pos parser.SourcePosition) *parser.Comment {
	var closestPreviousComment *parser.Comment
	for f.jsdocTracker.HasPendingCommentBefore(pos) {
		closestPreviousComment = f.jsdocTracker.Current()
		f.jsdocTracker.Advance()
	}

	return closestPreviousComment
}

// Java: IRFactory.parseJSDocInfoFrom(Comment)
func (f *IRFactory) ParseJSDocInfoFrom(comment *parser.Comment) *rhino.JSDocInfo {
	if comment == nil {
		return nil
	}

	if f.WithinClosureUnawareCodeRange(int(comment.Location.Start.Line), int(comment.Location.Start.Column)) {
		// Within closure-unaware code, we don't parse jsdoc as if it is
		// load-bearing. This early-return prevents the comments here from
		// being recorded as being parsed, and this later allows all these
		// jsdoc comments to be treated as "top-level" comments that might
		// contain license info / "important" comments (the only form of JSDoc
		// that should happen for closure-unaware code).
		f.jsdocTracker.Backtrack()
		return nil
	}

	jsDocParser := f.CreateJsDocInfoParser(comment)
	f.markCommentPtrParsed(comment)
	if f.HandlePossibleFileOverviewJsDoc(jsDocParser) {
		// Returning null here means that the comment was treated as a
		// fileoverview comment, and JSDocInfo should not then be attached to a
		// specific node, so we return null;
		return nil
	}
	return jsDocParser.RetrieveAndResetParsedJSDocInfo()
}

// Java: IRFactory.parseJSDocInfoOnTree(ParseTree)
func (f *IRFactory) ParseJSDocInfoOnTree(tree parser.ParseTree) *rhino.JSDocInfo {
	switch tree.GetParseTreeBase().Type {
	case parser.ParseTreeType_EXPRESSION_STATEMENT, parser.ParseTreeType_LABELLED_STATEMENT,
		parser.ParseTreeType_EXPORT_DECLARATION, parser.ParseTreeType_TEMPLATE_SUBSTITUTION:
		return nil
	case parser.ParseTreeType_CALL_EXPRESSION,
		parser.ParseTreeType_CONDITIONAL_EXPRESSION,
		parser.ParseTreeType_COMMA_EXPRESSION,
		parser.ParseTreeType_BINARY_OPERATOR,
		parser.ParseTreeType_MEMBER_EXPRESSION,
		parser.ParseTreeType_MEMBER_LOOKUP_EXPRESSION,
		parser.ParseTreeType_UPDATE_EXPRESSION:
		{
			nearest := IRFactory_FindNearestNode(tree)
			if nearest.GetParseTreeBase().Type == parser.ParseTreeType_PAREN_EXPRESSION {
				return nil
			}
		}
	default:
	}

	return f.ParseJSDocInfoFrom(f.GetJSDocCommentAt(tree.GetStart()))
}

// Java: IRFactory.parseJSDocInfoOnToken(parser.Token)
func (f *IRFactory) ParseJSDocInfoOnToken(token parser.Token) *rhino.JSDocInfo {
	return f.ParseJSDocInfoFrom(f.GetJSDocCommentAt(token.GetStart()))
}

// Java: IRFactory.parseInlineJSDocAt(SourcePosition)
func (f *IRFactory) ParseInlineJSDocAt(pos parser.SourcePosition) *rhino.JSDocInfo {
	if f.WithinClosureUnawareCodeRange(int(pos.Line), int(pos.Column)) {
		// Within closure-unaware code, we don't parse jsdoc as if it is
		// load-bearing. This early-return prevents the comments here from
		// being recorded as being parsed, and this later allows all these
		// jsdoc comments to be treated as "top-level" comments that might
		// contain license info / "important" comments (the only form of JSDoc
		// that should happen for closure-unaware code).
		// Unlike other points in IRFactory where we explicitly backtrack the
		// jsdocTracker, we don't do that here because we've actually never
		// advanced it yet - that happens in getJSDocCommentAt.
		return nil
	}
	comment := f.GetJSDocCommentAt(pos)
	if comment != nil && !strings.Contains(comment.Value, "@") {
		return f.ParseInlineTypeDoc(comment)
	}
	return f.ParseJSDocInfoFrom(comment)
}

// IRFactory_COMMENT_START_POSITION_COMPARATOR compares comments by start line,
// then start column.
//
// Java: IRFactory.COMMENT_START_POSITION_COMPARATOR
var IRFactory_COMMENT_START_POSITION_COMPARATOR = func(a, b *parser.Comment) int {
	if c := irFactory_compareInt(int(a.Location.Start.Line), int(b.Location.Start.Line)); c != 0 {
		return c
	}
	return irFactory_compareInt(int(a.Location.Start.Column), int(b.Location.Start.Column))
}

// Java: IRFactory.hasAnyPendingCommentBefore(SourcePosition, boolean)
func (f *IRFactory) HasAnyPendingCommentBefore(pos parser.SourcePosition, withinClosureUnawareCodeRange bool) bool {
	return f.nonJsdocTracker.HasPendingCommentBefore(pos) ||
		(withinClosureUnawareCodeRange && f.jsdocTracker.HasPendingCommentBefore(pos))
}

// ParseNonJSDocCommentAt creates a single NonJSDocComment from every comment
// associated with this node; or null if there are no such comments.
//
// It would be legal to replace all comments associated with this node with
// that one string.
//
// Java: IRFactory.parseNonJSDocCommentAt(SourcePosition, boolean)
func (f *IRFactory) ParseNonJSDocCommentAt(pos parser.SourcePosition, isInline bool) *rhino.NonJSDocComment {
	withinClosureUnawareCodeRange := f.WithinClosureUnawareCodeRange(int(pos.Line), int(pos.Column))
	if !f.HasAnyPendingCommentBefore(pos, withinClosureUnawareCodeRange) {
		return nil
	}
	includeAll := f.config.JsDocParsingMode() == Config_JsDocParsing_INCLUDE_ALL_COMMENTS

	var result strings.Builder
	var firstComment *parser.Comment
	var lastComment *parser.Comment

	for f.HasAnyPendingCommentBefore(pos, withinClosureUnawareCodeRange) {
		currentComment := f.nonJsdocTracker.Current()
		commentFromJsdocTracker := false

		if withinClosureUnawareCodeRange {
			currentJsDocComment := f.jsdocTracker.Current()
			if currentComment == nil ||
				(currentJsDocComment != nil &&
					IRFactory_COMMENT_START_POSITION_COMPARATOR(currentComment, currentJsDocComment) > 0) {

				if strings.Contains(currentJsDocComment.Value, "@license") {
					// Skip this so it can be parsed for the license contents
					// (which happens when a comment isn't inserted into the
					// parsedComments set).
					f.jsdocTracker.Advance()
					continue
				}

				commentFromJsdocTracker = true
				currentComment = currentJsDocComment
			}
		}

		if firstComment == nil {
			firstComment = currentComment
		}

		if includeAll {
			if lastComment != nil {
				for blankCount := currentComment.Location.Start.Line - lastComment.Location.End.Line; blankCount > 0; blankCount-- {
					result.WriteString("\n")
				}
			}
			result.WriteString(currentComment.Value)
		}

		lastComment = currentComment
		if commentFromJsdocTracker {
			f.markCommentParsed(f.jsdocTracker.index)
			f.jsdocTracker.Advance()
		} else {
			f.nonJsdocTracker.Advance()
		}
	}

	if firstComment == nil || !includeAll {
		return nil // We didn't find any relevant comments.
	}

	contents := result.String()
	nonJSDocComment := rhino.NewNonJSDocComment(
		firstComment.Location.Start, lastComment.Location.End, &contents)
	nonJSDocComment.SetEndsAsLineComment(lastComment.Type == parser.Comment_Type_LINE)
	nonJSDocComment.SetIsInline(isInline)
	return nonJSDocComment
}

// ParseTrailingNonJSDocCommentAt creates a single NonJSDocComment from the
// comment immediately following this node; or null if there is no such
// comment.
//
// tokenEnd: The end position after which we're looking for a trailing
// comment. possibleNextTokenStart: last source position that we're
// interested in when looking for a trailing comment.
//
// Java: IRFactory.parseTrailingNonJSDocCommentAt(SourcePosition, SourcePosition)
func (f *IRFactory) ParseTrailingNonJSDocCommentAt(
	tokenEnd parser.SourcePosition, possibleNextTokenStart parser.SourcePosition) *rhino.NonJSDocComment {
	if f.config.JsDocParsingMode() != Config_JsDocParsing_INCLUDE_ALL_COMMENTS {
		return nil
	}

	if !f.nonJsdocTracker.HasPendingCommentBefore(possibleNextTokenStart) {
		return nil
	}

	var result strings.Builder
	comment := f.nonJsdocTracker.Current()

	// Disregard comments that don't start on the same line
	if tokenEnd.Line != comment.Location.Start.Line {
		return nil
	}

	// Disregard comments that may be within the current statement
	if comment.Location.Start.Offset <= tokenEnd.Offset {
		return nil
	}

	// Check if there are only whitespace characters between the current token
	// and the start of the comment.
	preCommentText := f.fileWithContent.Substring(int(tokenEnd.Offset)+1, int(comment.Location.Start.Offset))

	if jutil.StringTrim(preCommentText) != "" {
		return nil
	}
	// TODO(user): handle multiple trailing comments like /* c1 */ /* c2 */ // c3
	result.WriteString(comment.Value)
	f.nonJsdocTracker.Advance()

	contents := result.String()
	nonJSDocComment := rhino.NewNonJSDocComment(comment.Location.Start, comment.Location.End, &contents)
	nonJSDocComment.SetEndsAsLineComment(comment.Type == parser.Comment_Type_LINE)
	return nonJSDocComment
}

// Java: IRFactory.findNearestNode(ParseTree)
func IRFactory_FindNearestNode(tree parser.ParseTree) parser.ParseTree {
	for {
		switch tree.GetParseTreeBase().Type {
		case parser.ParseTreeType_EXPRESSION_STATEMENT:
			tree = tree.(*parser.ExpressionStatementTree).Expression
			continue
		case parser.ParseTreeType_CALL_EXPRESSION:
			tree = tree.(*parser.CallExpressionTree).Operand
			continue
		case parser.ParseTreeType_BINARY_OPERATOR:
			tree = tree.(*parser.BinaryOperatorTree).Left
			continue
		case parser.ParseTreeType_CONDITIONAL_EXPRESSION:
			tree = tree.(*parser.ConditionalExpressionTree).Condition
			continue
		case parser.ParseTreeType_MEMBER_EXPRESSION:
			tree = tree.(*parser.MemberExpressionTree).Operand
			continue
		case parser.ParseTreeType_MEMBER_LOOKUP_EXPRESSION:
			tree = tree.(*parser.MemberLookupExpressionTree).Operand
			continue
		case parser.ParseTreeType_UPDATE_EXPRESSION:
			tree = tree.(*parser.UpdateExpressionTree).Operand
			continue
		case parser.ParseTreeType_COMMA_EXPRESSION:
			tree = tree.(*parser.CommaExpressionTree).Expressions[0]
			continue
		default:
			return tree
		}
	}
}

// Java: IRFactory.transform(ParseTree)
func (f *IRFactory) Transform(tree parser.ParseTree) *rhino.Node {
	info := f.ParseJSDocInfoOnTree(tree)
	comment := f.ParseNonJSDocCommentAt(tree.GetStart(), false)
	isRootOfClosureUnawareSubTree := info != nil && info.IsClosureUnawareCode()

	if isRootOfClosureUnawareSubTree {
		if f.closureUnawareTemplateNode == nil {
			f.closureUnawareTemplateNode = f.CreateClosureUnawareTemplateNode()
		}
		if f.closureUnawareCodeRanges == nil {
			f.closureUnawareCodeRanges = jutil.NewTreeRangeSet(irFactory_LineAndColumn_compare)
		}
		f.templateNode = f.closureUnawareTemplateNode
		nextClosureUnawareRange := tree.GetParseTreeBase().Location
		nextClosureUnawareRangeStart := IRFactory_LineAndColumn_FromSourcePosition(nextClosureUnawareRange.Start)
		nextClosureUnawareRangeEnd := IRFactory_LineAndColumn_FromSourcePosition(nextClosureUnawareRange.End)
		f.closureUnawareCodeRanges.Add(
			jutil.Range_Closed(nextClosureUnawareRangeStart, nextClosureUnawareRangeEnd, irFactory_LineAndColumn_compare))
	}

	node := f.transformDispatcher.Process(tree)

	if info != nil {
		node = f.MaybeInjectCastNode(tree, info, node)
		node.SetJSDocInfo(info)
	}
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	if isRootOfClosureUnawareSubTree {
		f.SetSourceInfo(node, tree)
		f.templateNode = f.defaultTemplateNode
		node = f.HideClosureUnawareCodeRoot(node)
	}
	f.SetSourceInfo(node, tree)
	return node
}

const IRFactory_JSCOMP_CLOSURE_UNAWARE_CODE_SHADOW_HOST_NAME = "$jscomp_wrap_closure_unaware_code"
const IRFactory_JSCOMP_CLOSURE_UNAWARE_SINK_VALUE = "$jscomp_sink_closure_unaware_impl"

// Java: IRFactory.hideClosureUnawareCodeRoot(Node)
func (f *IRFactory) HideClosureUnawareCodeRoot(fn *rhino.Node) *rhino.Node {
	if !fn.IsFunction() {
		f.Error(fn, "@closureUnaware root must be function")
		return fn
	}
	info := fn.GetJSDocInfo()
	if mode := info.GetPerFileClosureUnawareMode(); mode == nil || *mode != rhino.JSDocInfo_PerFileClosureUnawareMode_UNSPECIFIED {
		f.Error(fn, "@closureUnaware mode can only be specified at the fileoverview level")
	}
	fn.SetJSDocInfo(nil)

	shadowNameNode := rhino.IR_Name(IRFactory_JSCOMP_CLOSURE_UNAWARE_CODE_SHADOW_HOST_NAME)

	sinkCall := rhino.IR_Call(rhino.IR_Name(IRFactory_JSCOMP_CLOSURE_UNAWARE_SINK_VALUE).Srcref(f.templateNode), fn).
		Srcref(f.templateNode)
	sinkCall.PutBooleanProp(rhino.Node_FREE_CALL, true)
	script := rhino.IR_ScriptOfNodes(rhino.IR_ExprResult(sinkCall).Srcref(f.templateNode))
	shadowJsRoot := rhino.IR_Root(script)
	script.SetStaticSourceFileFrom(f.templateNode)
	script.SetIsInClosureUnawareSubtree(true)

	shadowNameNode.ClonePropsFrom(f.templateNode)
	shadowNameNode.SetClosureUnawareShadow(shadowJsRoot)
	return shadowNameNode
}

// Java: IRFactory.withinClosureUnawareCodeRange(int, int)
func (f *IRFactory) WithinClosureUnawareCodeRange(line int, lineColumnNo int) bool {
	if f.closureUnawareCodeRanges == nil || f.closureUnawareCodeRanges.IsEmpty() {
		return false
	}

	return f.closureUnawareCodeRanges.Contains(NewIRFactory_LineAndColumn(line, lineColumnNo))
}

// Java: IRFactory.maybeInjectCastNode(ParseTree, JSDocInfo, Node)
func (f *IRFactory) MaybeInjectCastNode(node parser.ParseTree, info *rhino.JSDocInfo, irNode *rhino.Node) *rhino.Node {
	if node.GetParseTreeBase().Type == parser.ParseTreeType_PAREN_EXPRESSION && info.HasType() {
		irNode = f.NewNodeWithChild(rhino.Token_CAST, irNode)
	}
	return irNode
}

// TransformNodeWithInlineComments: names and destructuring patterns, in
// parameters or variable declarations are special, because they can have
// inline type docs attached.
//
//	function f(/** string */ x) {}
//
// annotates 'x' as a string.
//
// See: Using Inline Doc Comments
// (http://code.google.com/p/jsdoc-toolkit/wiki/InlineDocs)
//
// Java: IRFactory.transformNodeWithInlineComments(ParseTree)
func (f *IRFactory) TransformNodeWithInlineComments(tree parser.ParseTree) *rhino.Node {
	info := f.ParseInlineJSDocAt(tree.GetStart())
	comment := f.ParseNonJSDocCommentAt(tree.GetStart(), true)

	node := f.transformDispatcher.Process(tree)

	if info != nil {
		node.SetJSDocInfo(info)
	}
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	f.SetSourceInfo(node, tree)
	return node
}

// Java: IRFactory.lineno(ParseTree)
func IRFactory_Lineno(node parser.ParseTree) int {
	return IRFactory_LinenoOfPosition(node.GetParseTreeBase().Location.Start)
}

// Java: IRFactory.lineno(parser.Token)
func IRFactory_LinenoOfToken(token parser.Token) int {
	return IRFactory_LinenoOfPosition(token.GetTokenBase().Location.Start)
}

// Java: IRFactory.lineno(SourcePosition)
func IRFactory_LinenoOfPosition(location parser.SourcePosition) int {
	// location lines start at zero, our AST starts at 1.
	return int(location.Line + 1)
}

// Java: IRFactory.charno(ParseTree)
func IRFactory_Charno(node parser.ParseTree) int {
	return IRFactory_CharnoOfPosition(node.GetParseTreeBase().Location.Start)
}

// Java: IRFactory.charno(parser.Token)
func IRFactory_CharnoOfToken(token parser.Token) int {
	return IRFactory_CharnoOfPosition(token.GetTokenBase().Location.Start)
}

// Java: IRFactory.charno(SourcePosition)
func IRFactory_CharnoOfPosition(location parser.SourcePosition) int {
	return int(location.Column)
}

// Java: IRFactory.languageFeatureWarningMessage(Feature)
func IRFactory_LanguageFeatureWarningMessage(feature parser.FeatureSet_Feature) string {
	forFeature := Config_LanguageMode_MinimumRequiredFor(feature)

	if forFeature == Config_LanguageMode_UNSUPPORTED {
		return "This language feature is not currently supported by the compiler: " + feature.String()
	} else {
		return "This language feature is only supported for " +
			Config_LanguageMode_MinimumRequiredFor(feature).String() +
			" mode or better: " +
			feature.String()
	}
}

// Java: IRFactory.maybeWarnForFeature(ParseTree, Feature)
func (f *IRFactory) MaybeWarnForFeature(node parser.ParseTree, feature parser.FeatureSet_Feature) {
	f.features = f.features.With(feature)
	location := node.GetParseTreeBase().Location
	if f.WithinClosureUnawareCodeRange(int(location.Start.Line), int(location.Start.Column)) {
		return
	}
	// Externs should not produce warnings about the input language mode.
	if f.currentFileIsExterns {
		return
	}
	if !f.IsSupportedForInputLanguageMode(feature) {
		f.errorReporter.Warning(
			IRFactory_LanguageFeatureWarningMessage(feature), f.sourceName, IRFactory_Lineno(node), IRFactory_Charno(node))
	}
}

// Java: IRFactory.maybeWarnForFeature(parser.Token, Feature)
func (f *IRFactory) MaybeWarnForFeatureToken(token parser.Token, feature parser.FeatureSet_Feature) {
	f.features = f.features.With(feature)
	location := token.GetTokenBase().Location
	if f.WithinClosureUnawareCodeRange(int(location.Start.Line), int(location.Start.Column)) {
		return
	}
	// Externs should not produce warnings about the input language mode.
	if f.currentFileIsExterns {
		return
	}
	if !f.IsSupportedForInputLanguageMode(feature) {
		f.errorReporter.Warning(
			IRFactory_LanguageFeatureWarningMessage(feature), f.sourceName, IRFactory_LinenoOfToken(token), IRFactory_CharnoOfToken(token))
	}
}

// Java: IRFactory.maybeWarnForFeature(Node, Feature)
func (f *IRFactory) MaybeWarnForFeatureNode(node *rhino.Node, feature parser.FeatureSet_Feature) {
	f.features = f.features.With(feature)
	if f.WithinClosureUnawareCodeRange(node.GetLineno(), node.GetCharno()) {
		return
	}
	// Externs should not produce warnings about the input language mode.
	if f.currentFileIsExterns {
		return
	}
	if !f.IsSupportedForInputLanguageMode(feature) {
		f.errorReporter.Warning(
			IRFactory_LanguageFeatureWarningMessage(feature), f.sourceName, node.GetLineno(), node.GetCharno())
	}
}

// Java: IRFactory.setSourceInfo(Node, Node)
func (f *IRFactory) SetSourceInfoFromNode(node *rhino.Node, ref *rhino.Node) {
	node.SetLinenoCharno(ref.GetLineno(), ref.GetCharno())
	f.SetLengthFrom(node, ref)
}

// Java: IRFactory.setSourceInfo(Node, ParseTree)
func (f *IRFactory) SetSourceInfo(irNode *rhino.Node, node parser.ParseTree) {
	if irNode.GetLineno() == -1 {
		location := node.GetParseTreeBase().Location
		f.SetSourceInfoFromPositions(irNode, location.Start, location.End)
	}
}

// Java: IRFactory.setSourceInfo(Node, parser.Token)
func (f *IRFactory) SetSourceInfoFromToken(irNode *rhino.Node, token parser.Token) {
	location := token.GetTokenBase().Location
	f.SetSourceInfoFromPositions(irNode, location.Start, location.End)
}

// Java: IRFactory.setSourceInfo(Node, SourcePosition, SourcePosition)
func (f *IRFactory) SetSourceInfoFromPositions(node *rhino.Node, start parser.SourcePosition, end parser.SourcePosition) {
	if node.GetLineno() == -1 {
		// If we didn't already set the line, then set it now. This avoids
		// cases like ParenthesizedExpression where we just return a previous
		// node, but don't want the new node to get its parent's line number.
		node.SetLinenoCharno(IRFactory_LinenoOfPosition(start), IRFactory_CharnoOfPosition(start))
		f.SetLength(node, start, end)
	}
}

// Java: IRFactory.createJsDocInfoParser(Comment)
func (f *IRFactory) CreateJsDocInfoParser(node *parser.Comment) *JsDocInfoParser {
	return f.CreateJsDocInfoParserWithLicensesOnlyConfig(node, false)
}

// CreateJsDocInfoParserWithLicensesOnlyConfig creates a JsDocInfoParser and
// parses the JsDoc string.
//
// Used both for handling individual JSDoc comments and for handling
// file-level JSDoc comments (@fileoverview and @license).
//
// node: The JsDoc Comment node to parse. Returns a JsDocInfoParser. Will
// contain either fileoverview JsDoc, or normal JsDoc, or no JsDoc (if the
// method parses to the wrong level).
//
// Java: IRFactory.createJsDocInfoParser(Comment, boolean)
func (f *IRFactory) CreateJsDocInfoParserWithLicensesOnlyConfig(node *parser.Comment, useLicensesOnlyConfig bool) *JsDocInfoParser {
	comment := node.Value
	lineno := IRFactory_LinenoOfPosition(node.Location.Start)
	charno := IRFactory_CharnoOfPosition(node.Location.Start)
	position := node.Location.Start.Offset

	config := f.config
	if useLicensesOnlyConfig {
		config = config.ToBuilder().SetJsDocParsingMode(Config_JsDocParsing_LICENSE_COMMENTS_ONLY).Build()
	}

	// The JsDocInfoParser expects the comment without the initial '/**'.
	const numOpeningChars = 3
	var commentBody string
	if len(comment) >= numOpeningChars && comment[0] < 0x80 && comment[1] < 0x80 && comment[2] < 0x80 {
		commentBody = comment[numOpeningChars:]
	} else {
		commentBody = jutil.StringSubstringFrom(comment, numOpeningChars)
	}
	f.jsdocStream.Init(commentBody, lineno, charno+numOpeningChars)
	jsdocParser := &f.jsdocParser
	jsdocParser.Init(
		&f.jsdocStream,
		comment,
		int(position),
		f.templateNode,
		config,
		f.jsDocSourceKind,
		f.errorReporter,
		f.jsdocSlabs)
	jsdocParser.SetFileOverviewJSDocInfo(f.firstFileoverview)
	if node.Type == parser.Comment_Type_IMPORTANT && len(node.Value) > 0 {
		jsdocParser.ParseImportantComment()
	} else {
		jsdocParser.Parse()
	}

	return jsdocParser
}

// ParseInlineTypeDoc parses inline type info.
//
// Java: IRFactory.parseInlineTypeDoc(Comment)
func (f *IRFactory) ParseInlineTypeDoc(node *parser.Comment) *rhino.JSDocInfo {
	comment := node.Value
	lineno := IRFactory_LinenoOfPosition(node.Location.Start)
	charno := IRFactory_CharnoOfPosition(node.Location.Start)

	// The JsDocInfoParser expects the comment without the initial '/**'.
	const numOpeningChars = 3
	var commentBody string
	if len(comment) >= numOpeningChars && comment[0] < 0x80 && comment[1] < 0x80 && comment[2] < 0x80 {
		commentBody = comment[numOpeningChars:]
	} else {
		commentBody = jutil.StringSubstringFrom(comment, numOpeningChars)
	}
	f.jsdocStream.Init(commentBody, lineno, charno+numOpeningChars)
	parser_ := &f.jsdocParser
	parser_.Init(
		&f.jsdocStream,
		comment,
		int(node.Location.Start.Offset),
		f.templateNode,
		f.config,
		f.jsDocSourceKind,
		f.errorReporter,
		f.jsdocSlabs)
	return parser_.ParseInlineTypeDoc()
}

// SetLength sets the length on the node if we're in IDE mode.
//
// Java: IRFactory.setLength(Node, SourcePosition, SourcePosition)
func (f *IRFactory) SetLength(node *rhino.Node, start parser.SourcePosition, end parser.SourcePosition) {
	node.SetLength(int(end.Offset - start.Offset))
}

// Java: IRFactory.setLengthFrom(Node, Node)
func (f *IRFactory) SetLengthFrom(node *rhino.Node, ref *rhino.Node) {
	node.SetLength(ref.GetLength())
}

var IRFactory_GOOG_MODULE = rhino.QualifiedName_Of("goog.module")

// IRFactory_TransformDispatcher dispatches each ParseTree to its process
// method.
//
// Java: private class IRFactory.TransformDispatcher (an inner class)
type IRFactory_TransformDispatcher struct {
	// Go-only: language-adapter — the enclosing IRFactory instance of the
	// inner class, embedded so that its members resolve as they do in Java.
	*IRFactory

	// Tracks whether a class is currently in scope. Presently used to detect
	// valid usage of private fields.
	classScope *IRFactory_TransformDispatcher_ScopeTracker

	// Used to detect an IdentifierExpression in the following context:
	// Valid: `class A { #x; static isA(o) { return #x in o; } }`
	// Invalid: `class A { #x; } function isA(o) { return #x in o; } }`
	privateIdLhsOfInScope *IRFactory_TransformDispatcher_ScopeTracker
}

// Java: IRFactory.TransformDispatcher()
func NewIRFactory_TransformDispatcher(outer *IRFactory) *IRFactory_TransformDispatcher {
	d := &IRFactory_TransformDispatcher{IRFactory: outer}
	d.classScope = NewIRFactory_TransformDispatcher_ScopeTracker()
	d.privateIdLhsOfInScope = NewIRFactory_TransformDispatcher_ScopeTracker()
	return d
}

// IRFactory_TransformDispatcher_ScopeTracker tracks whether a certain
// language feature is currently in scope. Use in conjunction with
// try-with-resources to auto-handle decrementing.
//
// Java: private static class IRFactory.TransformDispatcher.ScopeTracker
type IRFactory_TransformDispatcher_ScopeTracker struct {
	usageCount int
}

// Java: IRFactory.TransformDispatcher.ScopeTracker()
func NewIRFactory_TransformDispatcher_ScopeTracker() *IRFactory_TransformDispatcher_ScopeTracker {
	return &IRFactory_TransformDispatcher_ScopeTracker{usageCount: 0}
}

// Java: IRFactory.TransformDispatcher.ScopeTracker.inScope()
func (s *IRFactory_TransformDispatcher_ScopeTracker) InScope() bool {
	return s.usageCount > 0
}

// Java: IRFactory.TransformDispatcher.ScopeTracker.increment()
func (s *IRFactory_TransformDispatcher_ScopeTracker) Increment() irFactory_TransformDispatcher_ScopeTracker_autoDecrement {
	return s.MaybeIncrement(true)
}

// MaybeIncrement increments the usage count if condition holds and returns
// the AutoDecrement that undoes it.
//
// Go-only: language-adapter — returns the anonymous AutoDecrement class by
// value rather than as the AutoDecrement interface, so that it stays on the
// stack (the JVM's escape analysis removes the allocation; Go's would not
// once it is boxed in an interface).
//
// Java: IRFactory.TransformDispatcher.ScopeTracker.maybeIncrement(boolean)
func (s *IRFactory_TransformDispatcher_ScopeTracker) MaybeIncrement(condition bool) irFactory_TransformDispatcher_ScopeTracker_autoDecrement {
	if condition {
		jutil.CheckState(s.usageCount >= 0)
		s.usageCount++
	}
	return irFactory_TransformDispatcher_ScopeTracker_autoDecrement{tracker: s, condition: condition, closed: false}
}

// irFactory_TransformDispatcher_ScopeTracker_autoDecrement is the anonymous
// AutoDecrement class of ScopeTracker.maybeIncrement.
type irFactory_TransformDispatcher_ScopeTracker_autoDecrement struct {
	tracker   *IRFactory_TransformDispatcher_ScopeTracker
	condition bool
	closed    bool
}

// Java: IRFactory.TransformDispatcher.ScopeTracker.maybeIncrement(boolean).new AutoDecrement() {...}.close()
func (a *irFactory_TransformDispatcher_ScopeTracker_autoDecrement) Close() {
	jutil.CheckState(!a.closed)
	if a.condition {
		a.tracker.usageCount--
		jutil.CheckState(a.tracker.usageCount >= 0)
	}
	a.closed = true
}

// IRFactory_TransformDispatcher_AutoDecrement is a variant of AutoCloseable
// that does not throw.
//
// Java: private interface IRFactory.TransformDispatcher.AutoDecrement extends AutoCloseable
type IRFactory_TransformDispatcher_AutoDecrement interface {
	Close()
}

var _ IRFactory_TransformDispatcher_AutoDecrement = (*irFactory_TransformDispatcher_ScopeTracker_autoDecrement)(nil)

// IRFactory_TransformDispatcher_IdentifierType indicates the valid values
// usages for an identifier.
//
// Java: private enum IRFactory.TransformDispatcher.IdentifierType
type IRFactory_TransformDispatcher_IdentifierType int32

const (
	// The identifier can never be private.
	IRFactory_TransformDispatcher_IdentifierType_STANDARD IRFactory_TransformDispatcher_IdentifierType = iota
	// The identifier can be private if it is in a class scope.
	IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE
)

var irFactory_TransformDispatcher_IdentifierType_names = [...]string{"STANDARD", "CAN_BE_PRIVATE"}

// Java: IRFactory.TransformDispatcher.IdentifierType.values()
func IRFactory_TransformDispatcher_IdentifierType_values() []IRFactory_TransformDispatcher_IdentifierType {
	return []IRFactory_TransformDispatcher_IdentifierType{
		IRFactory_TransformDispatcher_IdentifierType_STANDARD,
		IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE,
	}
}

// String is Enum.name().
func (t IRFactory_TransformDispatcher_IdentifierType) String() string {
	return irFactory_TransformDispatcher_IdentifierType_names[t]
}

// NewStringNodeFromIdentifier creates a StringNode using the value of an
// IdentifierToken. Reports an error if the identifier value is private and
// the identifier type is STANDARD or the identifier is not within the scope
// of a class.
//
// Java: IRFactory.TransformDispatcher.newStringNodeFromIdentifier(Token, IdentifierType, IdentifierToken)
func (d *IRFactory_TransformDispatcher) NewStringNodeFromIdentifier(
	type_ rhino.Token, identifierType IRFactory_TransformDispatcher_IdentifierType, identifierToken *parser.IdentifierToken) *rhino.Node {
	// Private properties can only ever be referenced from within the scope of
	// a class. If we're not in a class, we always use STANDARD rules.
	if !d.classScope.InScope() {
		identifierType = IRFactory_TransformDispatcher_IdentifierType_STANDARD
	}

	var value string
	if identifierToken.IsPrivateIdentifier() {
		if identifierType == IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE {
			d.MaybeWarnForFeatureToken(identifierToken, parser.FeatureSet_Feature_PRIVATE_ELEMENTS)
		} else {
			d.errorReporter.Error(
				"Private identifiers may not be used in this context",
				d.sourceName,
				IRFactory_LinenoOfToken(identifierToken),
				IRFactory_CharnoOfToken(identifierToken))
		}
		value = identifierToken.GetMaybePrivateValue()
	} else {
		value = identifierToken.GetValue()
	}

	node := d.NewStringNodeWithToken(type_, value)
	if identifierType == IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE &&
		identifierToken.IsPrivateIdentifier() {
		node.SetPrivateIdentifier()
	}
	return node
}

// ProcessObjectLitKey transforms the given object key `input` into a Node
// with token `output`.
//
// Depending on `input`, this may add quoting, such as for numerical values.
// For example, in `{2: null}`, `2` will be transformed into a quoted string.
// This adjustment loses some accuracy about the source code, but simplifies
// the AST.
//
// Get and set accessors process their name using this method. For classes,
// this can be a private property in which case identifierType will be set to
// CAN_BE_PRIVATE. E.g. `class C { get #f() { return value; } }`.
//
// Java: IRFactory.TransformDispatcher.processObjectLitKey(parser.Token, Token, IdentifierType)
func (d *IRFactory_TransformDispatcher) ProcessObjectLitKey(
	input parser.Token,
	output rhino.Token,
	identifierType IRFactory_TransformDispatcher_IdentifierType) *rhino.Node {
	if input == nil {
		return d.CreateMissingExpressionNode()
	}

	if input.GetTokenBase().Type == parser.TokenType_IDENTIFIER {
		return d.ProcessNameOfToken(input.(*parser.IdentifierToken), output, identifierType)
	}

	literal := input.(parser.LiteralToken)
	jsDocInfo := d.ParseJSDocInfoOnToken(literal)
	comment := d.ParseNonJSDocCommentAt(literal.GetStart(), true)

	var node *rhino.Node
	switch input.GetTokenBase().Type {
	case parser.TokenType_NUMBER:
		node = d.NewStringNodeWithToken(output, dtoa.DToA_NumberToString(d.NormalizeNumber(literal)))
	case parser.TokenType_BIGINT:
		node = d.NewStringNodeWithToken(output, d.NormalizeBigInt(literal).String())
	default:
		node = d.NewStringNodeWithToken(output, d.NormalizeString(literal, false))
	}

	if jsDocInfo != nil {
		node.SetJSDocInfo(jsDocInfo)
	}
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}

	d.SetSourceInfoFromToken(node, literal)
	node.PutBooleanProp(rhino.Node_QUOTED_PROP, true)
	return node
}

// Java: IRFactory.TransformDispatcher.processComprehension(ComprehensionTree)
func (d *IRFactory_TransformDispatcher) ProcessComprehension(tree *parser.ComprehensionTree) *rhino.Node {
	return d.UnsupportedLanguageFeature(tree, "array/generator comprehensions")
}

// Java: IRFactory.TransformDispatcher.processComprehensionFor(ComprehensionForTree)
func (d *IRFactory_TransformDispatcher) ProcessComprehensionFor(tree *parser.ComprehensionForTree) *rhino.Node {
	return d.UnsupportedLanguageFeature(tree, "array/generator comprehensions")
}

// Java: IRFactory.TransformDispatcher.processComprehensionIf(ComprehensionIfTree)
func (d *IRFactory_TransformDispatcher) ProcessComprehensionIf(tree *parser.ComprehensionIfTree) *rhino.Node {
	return d.UnsupportedLanguageFeature(tree, "array/generator comprehensions")
}

// Java: IRFactory.TransformDispatcher.processArrayLiteral(ArrayLiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessArrayLiteral(tree *parser.ArrayLiteralExpressionTree) *rhino.Node {
	node := d.NewNode(rhino.Token_ARRAYLIT)
	node.SetTrailingComma(tree.HasTrailingComma)
	for _, child := range tree.Elements {
		c := d.Transform(child)
		node.AddChildToBack(c)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processArrayPattern(ArrayPatternTree)
func (d *IRFactory_TransformDispatcher) ProcessArrayPattern(tree *parser.ArrayPatternTree) *rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_ARRAY_DESTRUCTURING)

	node := d.NewNode(rhino.Token_ARRAY_PATTERN)
	for _, child := range tree.Elements {
		var elementNode *rhino.Node
		switch child.GetParseTreeBase().Type {
		case parser.ParseTreeType_DEFAULT_PARAMETER:
			// processDefaultParameter() knows how to find and apply inline
			// JSDoc to the right node
			elementNode = d.ProcessDefaultParameter(child.(*parser.DefaultParameterTree))
		case parser.ParseTreeType_ITER_REST:
			d.MaybeWarnForFeature(child, parser.FeatureSet_Feature_ARRAY_PATTERN_REST)
			elementNode = d.TransformNodeWithInlineComments(child)
		default:
			elementNode = d.TransformNodeWithInlineComments(child)
		}
		node.AddChildToBack(elementNode)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processObjectPattern(ObjectPatternTree)
func (d *IRFactory_TransformDispatcher) ProcessObjectPattern(tree *parser.ObjectPatternTree) *rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_OBJECT_DESTRUCTURING)

	node := d.NewNode(rhino.Token_OBJECT_PATTERN)
	for _, child := range tree.Fields {
		childNode := d.ProcessObjectPatternElement(child)
		node.AddChildToBack(childNode)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processObjectPatternElement(ParseTree)
func (d *IRFactory_TransformDispatcher) ProcessObjectPatternElement(child parser.ParseTree) *rhino.Node {
	switch child.GetParseTreeBase().Type {
	case parser.ParseTreeType_DEFAULT_PARAMETER:
		// shorthand with a default value
		// let { /** inlineType */ name = default } = something;
		return d.ProcessObjectPatternShorthandWithDefault(child.(*parser.DefaultParameterTree))
	case parser.ParseTreeType_PROPERTY_NAME_ASSIGNMENT:
		return d.ProcessObjectPatternPropertyNameAssignment(child.(*parser.PropertyNameAssignmentTree))
	case parser.ParseTreeType_COMPUTED_PROPERTY_DEFINITION:
		// let {[expression]: /** inlineType */ name} = something;
		computedPropertyDefinition := child.(*parser.ComputedPropertyDefinitionTree)
		return d.ProcessObjectPatternComputedPropertyDefinition(computedPropertyDefinition)
	case parser.ParseTreeType_OBJECT_REST:
		// let {...restObject} = someObject;
		d.MaybeWarnForFeature(child, parser.FeatureSet_Feature_OBJECT_PATTERN_REST)
		target := d.TransformNodeWithInlineComments(child.(*parser.ObjectRestTree).AssignmentTarget)
		rest := d.NewNodeWithChild(rhino.Token_OBJECT_REST, target)
		d.SetSourceInfo(rest, child)
		return rest
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unexpected object pattern element: " + child.String()))
	}
}

// ProcessObjectPatternShorthandWithDefault processes an object pattern
// element using shorthand and a default value.
//
// e.g. the `/** inlineType * / name = default` part of this code
//
//	let { /** inlineType * / name = default } = something;
//
// Java: IRFactory.TransformDispatcher.processObjectPatternShorthandWithDefault(DefaultParameterTree)
func (d *IRFactory_TransformDispatcher) ProcessObjectPatternShorthandWithDefault(defaultParameter *parser.DefaultParameterTree) *rhino.Node {
	defaultValueNode := d.ProcessDefaultParameter(defaultParameter)
	// Store in AST as non-shorthand form & just note it was originally
	// shorthand {name: /**inlineType */ name = default }
	nameNode := defaultValueNode.GetFirstChild()
	stringKeyNode := d.MaybeAddNonJsDocComment(
		d.NewStringNodeWithToken(rhino.Token_STRING_KEY, nameNode.GetString()), defaultParameter.GetStart())
	d.SetSourceInfoFromNode(stringKeyNode, nameNode)
	stringKeyNode.SetShorthandProperty(true)
	stringKeyNode.AddChildToBack(defaultValueNode)
	return stringKeyNode
}

// ProcessObjectPatternPropertyNameAssignment processes property name
// assignments in an object pattern.
//
// Covers these cases.
//
//	let {name} = someObject;
//	let {key: name} = someObject;
//	let {key: name = something} = someObject;
//
// Java: IRFactory.TransformDispatcher.processObjectPatternPropertyNameAssignment(PropertyNameAssignmentTree)
func (d *IRFactory_TransformDispatcher) ProcessObjectPatternPropertyNameAssignment(
	propertyNameAssignment *parser.PropertyNameAssignmentTree) *rhino.Node {
	key := d.ProcessObjectLitKey(
		propertyNameAssignment.Name, rhino.Token_STRING_KEY, IRFactory_TransformDispatcher_IdentifierType_STANDARD)
	targetTree := propertyNameAssignment.Value
	var valueNode *rhino.Node
	if targetTree == nil {
		// `let { /** inlineType */ key } = something;`
		// The key is also the target name.
		valueNode = d.ProcessNameWithInlineCommentsOfToken(
			propertyNameAssignment.Name.(*parser.IdentifierToken), IRFactory_TransformDispatcher_IdentifierType_STANDARD)
		key.SetShorthandProperty(true)
	} else {
		valueNode = d.ProcessDestructuringElementTarget(targetTree)
	}
	key.AddChildToFront(valueNode)
	return key
}

// Java: IRFactory.TransformDispatcher.processDestructuringElementTarget(ParseTree)
func (d *IRFactory_TransformDispatcher) ProcessDestructuringElementTarget(targetTree parser.ParseTree) *rhino.Node {
	var valueNode *rhino.Node
	if targetTree.GetParseTreeBase().Type == parser.ParseTreeType_DEFAULT_PARAMETER {
		// let {key: /** inlineType */ name = default} = something;
		// let [/** inlineType */ name = default] = something;
		// processDefaultParameter() knows how to apply the inline JSDoc, if
		// any, to the right node
		valueNode = d.ProcessDefaultParameter(targetTree.(*parser.DefaultParameterTree))
	} else if targetTree.GetParseTreeBase().Type == parser.ParseTreeType_IDENTIFIER_EXPRESSION {
		// let {key: /** inlineType */ name} = something
		// let [/** inlineType */ name] = something
		// Allow inline JSDoc on the name, since we may well be declaring it
		// here.
		valueNode = d.ProcessNameWithInlineComments(
			targetTree.(*parser.IdentifierExpressionTree), IRFactory_TransformDispatcher_IdentifierType_STANDARD)
	} else {
		// ({prop: /** string */ ns.a.b} = someObject);
		// NOTE: CheckJSDoc will report an error for this case, since we want
		// qualified names to be declared with individual statements, like
		// `/** @type {string} */ ns.a.b;`
		valueNode = d.TransformNodeWithInlineComments(targetTree)
	}
	return valueNode
}

// ProcessObjectPatternComputedPropertyDefinition processes a computed
// property in an object pattern.
//
// Covers these cases:
//
//	let {[expression]: name} = someObject;
//	let {[expression]: name = defaultValue} = someObject;
//
// Java: IRFactory.TransformDispatcher.processObjectPatternComputedPropertyDefinition(ComputedPropertyDefinitionTree)
func (d *IRFactory_TransformDispatcher) ProcessObjectPatternComputedPropertyDefinition(
	computedPropertyDefinition *parser.ComputedPropertyDefinitionTree) *rhino.Node {
	d.MaybeWarnForFeature(computedPropertyDefinition, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)

	expressionNode := d.Transform(computedPropertyDefinition.Property)

	valueTree := computedPropertyDefinition.Value
	valueNode := d.ProcessDestructuringElementTarget(valueTree)
	computedPropertyNode := d.NewNodeWith2Children(rhino.Token_COMPUTED_PROP, expressionNode, valueNode)
	d.SetSourceInfo(computedPropertyNode, computedPropertyDefinition)
	return computedPropertyNode
}

// Java: IRFactory.TransformDispatcher.processAstRoot(ProgramTree)
func (d *IRFactory_TransformDispatcher) ProcessAstRoot(rootNode *parser.ProgramTree) *rhino.Node {
	scriptNode := d.NewNode(rhino.Token_SCRIPT)
	for _, child := range rootNode.SourceElements {
		scriptNode.AddChildToBack(d.Transform(child))
	}
	d.ParseDirectives(scriptNode)
	isGoogModule := d.IsGoogModuleFile(scriptNode)
	if isGoogModule || d.features.Has(parser.FeatureSet_Feature_MODULES) {
		moduleNode := d.NewNode(rhino.Token_MODULE_BODY)
		d.SetSourceInfo(moduleNode, rootNode)
		moduleNode.AddChildrenToBack(scriptNode.RemoveChildren())
		scriptNode.AddChildToBack(moduleNode)
		if isGoogModule {
			scriptNode.PutBooleanProp(rhino.Node_GOOG_MODULE, true)
		} else {
			scriptNode.PutBooleanProp(rhino.Node_ES6_MODULE, true)
		}
	}
	return scriptNode
}

// Java: IRFactory.TransformDispatcher.isGoogModuleFile(Node)
func (d *IRFactory_TransformDispatcher) IsGoogModuleFile(scriptNode *rhino.Node) bool {
	jutil.CheckArgument(scriptNode.IsScript())
	if !scriptNode.HasChildren() {
		return false
	}
	exprResult := scriptNode.GetFirstChild()
	if !exprResult.IsExprResult() {
		return false
	}
	call := exprResult.GetFirstChild()
	if !call.IsCall() {
		return false
	}
	return IRFactory_GOOG_MODULE.Matches(call.GetFirstChild()) &&
		call.HasTwoChildren() &&
		call.GetSecondChild().IsStringLit()
}

// ParseDirectives parses the directives, encodes them in the AST, and
// removes their nodes.
//
// The only suported directive is "use strict".
//
// For information on ES5 directives, see section 14.1 of ECMA-262, Edition 5.
//
// Java: IRFactory.TransformDispatcher.parseDirectives(Node)
func (d *IRFactory_TransformDispatcher) ParseDirectives(node *rhino.Node) {
	useStrict := false
	for {
		statement := node.GetFirstChild()
		if statement == nil || !statement.IsExprResult() {
			break
		}

		directive := statement.GetFirstChild()
		if !directive.IsStringLit() || directive.GetString() != "use strict" {
			break
		}

		useStrict = true
		// TODO(nickreid): Should we also remove other directives that aren't
		// "use strict"?
		statement.Detach()
	}

	if useStrict {
		node.SetUseStrict(true)
	}
}

// Java: IRFactory.TransformDispatcher.processBlock(BlockTree)
func (d *IRFactory_TransformDispatcher) ProcessBlock(blockNode *parser.BlockTree) *rhino.Node {
	node := d.NewNode(rhino.Token_BLOCK)
	for i, child := range blockNode.Statements {
		childNode := d.Transform(child)
		node.AddChildToBack(childNode)
		d.AttachPossibleTrailingComment(childNode, child.GetEnd())
		// Go-only: performance — nothing reads a statement's parse tree once it
		// is transformed: dropping it lets the GC free it while the rest of the
		// file is transformed, instead of holding a large file's whole parse
		// tree next to its whole AST.
		blockNode.Statements[i] = nil
	}

	lastComment := d.ParseNonJSDocCommentAt(blockNode.GetEnd(), false)
	d.AddExtraTrailingComment(node, lastComment)
	return node
}

// Java: IRFactory.TransformDispatcher.processBreakStatement(BreakStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessBreakStatement(statementNode *parser.BreakStatementTree) *rhino.Node {
	node := d.NewNode(rhino.Token_BREAK)
	if statementNode.HasLabel() {
		labelName := d.TransformLabelName(statementNode.Name)
		node.AddChildToBack(labelName)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.transformLabelName(IdentifierToken)
func (d *IRFactory_TransformDispatcher) TransformLabelName(token *parser.IdentifierToken) *rhino.Node {
	label := d.MaybeAddNonJsDocComment(
		d.NewStringNodeFromIdentifier(rhino.Token_LABEL_NAME, IRFactory_TransformDispatcher_IdentifierType_STANDARD, token),
		token.GetStart())
	d.SetSourceInfoFromToken(label, token)
	return label
}

// Java: IRFactory.TransformDispatcher.processConditionalExpression(ConditionalExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessConditionalExpression(exprNode *parser.ConditionalExpressionTree) *rhino.Node {
	return d.NewNodeWith3Children(
		rhino.Token_HOOK,
		d.Transform(exprNode.Condition),
		d.Transform(exprNode.Left),
		d.Transform(exprNode.Right))
}

// Java: IRFactory.TransformDispatcher.processContinueStatement(ContinueStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessContinueStatement(statementNode *parser.ContinueStatementTree) *rhino.Node {
	node := d.NewNode(rhino.Token_CONTINUE)
	if statementNode.HasLabel() {
		labelName := d.TransformLabelName(statementNode.Name)
		node.AddChildToBack(labelName)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processDoLoop(DoWhileStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessDoLoop(loopNode *parser.DoWhileStatementTree) *rhino.Node {
	return d.NewNodeWith2Children(rhino.Token_DO, d.TransformBlock(loopNode.Body), d.Transform(loopNode.Condition))
}

// Java: IRFactory.TransformDispatcher.processElementGet(MemberLookupExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessElementGet(getNode *parser.MemberLookupExpressionTree) *rhino.Node {
	return d.NewNodeWith2Children(
		rhino.Token_GETELEM, d.Transform(getNode.Operand), d.Transform(getNode.MemberExpression))
}

// Java: IRFactory.TransformDispatcher.processOptChainElementGet(OptionalMemberLookupExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessOptChainElementGet(getNode *parser.OptionalMemberLookupExpressionTree) *rhino.Node {
	d.MaybeWarnForFeature(getNode, parser.FeatureSet_Feature_OPTIONAL_CHAINING)
	getElem := d.NewNodeWith2Children(
		rhino.Token_OPTCHAIN_GETELEM,
		d.Transform(getNode.Operand),
		d.Transform(getNode.MemberExpression))
	getElem.SetIsOptionalChainStart(getNode.IsStartOfOptionalChain)
	return getElem
}

// Java: IRFactory.TransformDispatcher.processEmptyStatement(EmptyStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessEmptyStatement(unused *parser.EmptyStatementTree) *rhino.Node {
	return d.NewNode(rhino.Token_EMPTY)
}

// Java: IRFactory.TransformDispatcher.processExpressionStatement(ExpressionStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessExpressionStatement(statementNode *parser.ExpressionStatementTree) *rhino.Node {
	node := d.NewNode(rhino.Token_EXPR_RESULT)
	node.AddChildToBack(d.Transform(statementNode.Expression))
	return node
}

// Java: IRFactory.TransformDispatcher.processForInLoop(ForInStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessForInLoop(loopNode *parser.ForInStatementTree) *rhino.Node {
	// TODO(bradfordcsmith): Rename initializer to something more intuitive
	// like "lhs"
	initializer := d.Transform(loopNode.Initializer)
	return d.NewNodeWith3Children(
		rhino.Token_FOR_IN, initializer, d.Transform(loopNode.Collection), d.TransformBlock(loopNode.Body))
}

// Java: IRFactory.TransformDispatcher.processForOf(ForOfStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessForOf(loopNode *parser.ForOfStatementTree) *rhino.Node {
	d.MaybeWarnForFeature(loopNode, parser.FeatureSet_Feature_FOR_OF)
	initializer := d.Transform(loopNode.Initializer)
	return d.NewNodeWith3Children(
		rhino.Token_FOR_OF, initializer, d.Transform(loopNode.Collection), d.TransformBlock(loopNode.Body))
}

// Java: IRFactory.TransformDispatcher.processForAwaitOf(ForAwaitOfStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessForAwaitOf(loopNode *parser.ForAwaitOfStatementTree) *rhino.Node {
	d.MaybeWarnForFeature(loopNode, parser.FeatureSet_Feature_FOR_AWAIT_OF)
	initializer := d.Transform(loopNode.Initializer)
	return d.NewNodeWith3Children(
		rhino.Token_FOR_AWAIT_OF,
		initializer,
		d.Transform(loopNode.Collection),
		d.TransformBlock(loopNode.Body))
}

// Java: IRFactory.TransformDispatcher.processForLoop(ForStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessForLoop(loopNode *parser.ForStatementTree) *rhino.Node {
	node := d.NewNodeWith3Children(
		rhino.Token_FOR,
		d.TransformOrEmpty(loopNode.Initializer, loopNode),
		d.TransformOrEmpty(loopNode.Condition, loopNode),
		d.TransformOrEmpty(loopNode.Increment, loopNode))
	node.AddChildToBack(d.TransformBlock(loopNode.Body))
	return node
}

// Java: IRFactory.TransformDispatcher.transformOrEmpty(ParseTree, ParseTree)
func (d *IRFactory_TransformDispatcher) TransformOrEmpty(tree parser.ParseTree, parent parser.ParseTree) *rhino.Node {
	if tree == nil {
		n := d.NewNode(rhino.Token_EMPTY)
		d.SetSourceInfo(n, parent)
		return n
	}
	return d.Transform(tree)
}

// Java: IRFactory.TransformDispatcher.transformOrEmpty(IdentifierToken, ParseTree)
func (d *IRFactory_TransformDispatcher) TransformOrEmptyIdentifier(token *parser.IdentifierToken, parent parser.ParseTree) *rhino.Node {
	if token == nil {
		n := d.NewNode(rhino.Token_EMPTY)
		d.SetSourceInfo(n, parent)
		return n
	}
	return d.ProcessNameOfToken(token, rhino.Token_NAME, IRFactory_TransformDispatcher_IdentifierType_STANDARD)
}

// Java: IRFactory.TransformDispatcher.processFunctionCall(CallExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessFunctionCall(callNode *parser.CallExpressionTree) *rhino.Node {
	node := d.NewNodeWithChild(rhino.Token_CALL, d.Transform(callNode.Operand))
	node.SetTrailingComma(callNode.Arguments.HasTrailingComma)
	argumentsTree := callNode.Arguments
	argCount := 0
	for _, child := range argumentsTree.Arguments {
		childNode := d.Transform(child)
		node.AddChildToBack(childNode)
		// The non-trailing comments are already attached to `childNode` in
		// `transform(child)` call. Now we must attach possible trailing
		// comments to `childNode`.

		d.attachPossibleTrailingCommentsForArgAt(
			childNode, argumentsTree.Arguments, argCount, argumentsTree.CommaPositions, argumentsTree.Location.End)
	}
	d.AnnotateCalls(node)
	return node
}

// AnnotateCalls: there are two types of calls we are interested in calls
// without explicit "this" values (what we are call "free" calls) and direct
// call to eval.
//
// Java: IRFactory.TransformDispatcher.annotateCalls(Node)
func (d *IRFactory_TransformDispatcher) AnnotateCalls(n *rhino.Node) {
	jutil.CheckStateMsg1(n.IsCall() || n.IsOptChainCall() || n.IsTaggedTemplateLit(), "%s", n)

	// Keep track of of the "this" context of a call.  A call without an
	// explicit "this" is a free call.
	callee := n.GetFirstChild()

	// ignore cast nodes.
	for callee.IsCast() {
		callee = callee.GetFirstChild()
	}

	if !d.IsNormalOrOptChainGet(callee) {
		// This call originally was not passed a `this` value.
		// Inlining could change the callee into a property reference of some
		// kind. The code printer will recognize the `FREE_CALL` property and
		// wrap the real callee with `(0, real.callee)(args)` when necessary to
		// avoid changing the calling behavior.
		n.PutBooleanProp(rhino.Node_FREE_CALL, true)

		if callee.IsName() && callee.GetString() == "eval" {
			// Keep track of the context in which eval is called. It is
			// important to distinguish between "(0, eval)()" and "eval()".
			callee.PutBooleanProp(rhino.Node_DIRECT_EVAL, true)
		} else if callee.IsComma() && callee.GetFirstChild().IsNumber() {
			// The input code may actually already contain calls of the form
			// `(0, real.callee)(arg1, arg2)`. For example, the TypeScript
			// compiler outputs code like this in some cases.
			// This makes it hard for us to connect calls with the called
			// functions, though, so we will simplify these cases. The
			// `FREE_CALL` property we've just applied will tell the code
			// printer to put this wrapping back later.
			realCallee := callee.GetSecondChild()
			callee.ReplaceWith(realCallee.Detach())
			// TODO(bradfordcsmith): Why do I get an NPE if I try to report
			// this code change?
		}
	}
}

// IsNormalOrOptChainGet: is this a GETPROP, OPTCHAIN_GETPROP, GETELEM, or
// OPTCHAIN_GETELEM?
//
// Java: IRFactory.TransformDispatcher.isNormalOrOptChainGet(Node)
func (d *IRFactory_TransformDispatcher) IsNormalOrOptChainGet(n *rhino.Node) bool {
	return n.IsGetProp() || n.IsGetElem() || n.IsOptChainGetProp() || n.IsOptChainGetElem()
}

// GetEndOfArgCommentZones calculates, for each arg, a location (end
// SourcePosition) such that all trailing comments before this location
// correspond to that arg. Can be used while processing both function calls
// (ArgsList) as well as declarations (ParamList).
//
// args: list of arguments or formal parameters (ParseTree nodes).
// commaPositions: list of SourcePositions corresponding to commas in the
// argsList or formal parameter list. argListEndPosition: SourcePosition of
// the end of argsList or paramList.
//
// Java: IRFactory.TransformDispatcher.getEndOfArgCommentZones(ImmutableList<ParseTree>, ImmutableList<SourcePosition>, SourcePosition)
func (d *IRFactory_TransformDispatcher) GetEndOfArgCommentZones(
	args []parser.ParseTree,
	commaPositions []parser.SourcePosition,
	argListEndPosition parser.SourcePosition) []parser.SourcePosition {
	zones := make([]parser.SourcePosition, 0, len(args))
	commaCount := 0
	for _, arg := range args {
		if len(args) > commaCount+1 {
			// there is a next arg after this arg
			nextParam := args[commaCount+1]
			if nextParam.GetParseTreeBase().Location.Start.Line > arg.GetParseTreeBase().Location.End.Line {
				// Next arg is on a new line; all trailing comments on this line
				// belong to this arg create a source position to represent the
				// end of current line
				tempSourcePos := parser.SourcePosition{
					Offset: math.MaxInt32,
					Line:   arg.GetParseTreeBase().Location.End.Line,
					Column: math.MaxInt32,
				}
				zones = append(zones, tempSourcePos)
			} else {
				// Next arg is on the same line; trailing comments before the
				// comma belong to this arg
				commaPosition := commaPositions[commaCount]
				zones = append(zones, commaPosition)
			}
		} else {
			// last arg; trailing comments till the end of argList belong to
			// this arg
			zones = append(zones, argListEndPosition)
		}
		commaCount++
	}
	return zones
}

// AttachPossibleTrailingCommentsForArg attaches trailing comments
// associated with this arg or formal param to it.
//
// paramNode: The node to which we're attaching trailing comment. endZone:
// The end location until which we fetch pending comments for attachment.
//
// Java: IRFactory.TransformDispatcher.attachPossibleTrailingCommentsForArg(Node, SourcePosition)
func (d *IRFactory_TransformDispatcher) AttachPossibleTrailingCommentsForArg(paramNode *rhino.Node, endZone parser.SourcePosition) {
	trailingComment := d.ParseNonJSDocCommentAt(endZone, true)
	if trailingComment == nil {
		return
	}
	paramNode.SetTrailingNonJSDocComment(trailingComment)
}

func (d *IRFactory_TransformDispatcher) attachPossibleTrailingCommentsForArgAt(
	paramNode *rhino.Node,
	args []parser.ParseTree,
	commaCount int,
	commaPositions []parser.SourcePosition,
	argListEndPosition parser.SourcePosition,
) {
	if len(args) > commaCount+1 {
		arg := args[commaCount]
		nextParam := args[commaCount+1]
		if nextParam.GetParseTreeBase().Location.Start.Line > arg.GetParseTreeBase().Location.End.Line {
			tempSourcePos := parser.SourcePosition{
				Offset: math.MaxInt32,
				Line:   arg.GetParseTreeBase().Location.End.Line,
				Column: math.MaxInt32,
			}
			d.AttachPossibleTrailingCommentsForArg(paramNode, tempSourcePos)
			return
		}
		d.AttachPossibleTrailingCommentsForArg(paramNode, commaPositions[commaCount])
		return
	}
	d.AttachPossibleTrailingCommentsForArg(paramNode, argListEndPosition)
}

// AttachPossibleTrailingComment attaches trailing comments associated with
// this node.
//
// node: The node to which we're attaching trailing comment. tokenEnd: The
// end location after which we're looking for a trailing comment.
//
// Java: IRFactory.TransformDispatcher.attachPossibleTrailingComment(Node, SourcePosition)
func (d *IRFactory_TransformDispatcher) AttachPossibleTrailingComment(node *rhino.Node, tokenEnd parser.SourcePosition) {
	// Check for pending comments on the rest of the current line.
	nextLine := parser.SourcePosition{Offset: math.MaxInt32, Line: tokenEnd.Line + 1, Column: 0}
	trailingComment := d.ParseTrailingNonJSDocCommentAt(tokenEnd, nextLine)
	if trailingComment == nil {
		return
	}
	node.SetTrailingNonJSDocComment(trailingComment)
}

// Java: IRFactory.TransformDispatcher.addExtraTrailingComment(Node, NonJSDocComment)
func (d *IRFactory_TransformDispatcher) AddExtraTrailingComment(node *rhino.Node, lastComment *rhino.NonJSDocComment) {
	if lastComment == nil {
		return
	}
	if !node.HasChildren() {
		node.SetTrailingNonJSDocComment(lastComment)
		return
	}

	lastChild := node.GetLastChild()
	currentComment := lastChild.GetTrailingNonJSDocComment()

	if currentComment == nil {
		// We add a line break here as this is cannot be a trailing comment on
		// the same line as the last child.
		// TODO(user): pass in proper child tree end position to get the number
		// of line breaks right.
		newStart := parser.SourcePosition{
			Offset: lastComment.GetEndPosition().Offset - 1,
			Line:   lastComment.GetEndPosition().Line - 1,
			Column: 0,
		}

		contents := "\n" + lastComment.GetCommentString()
		newlineComment := rhino.NewNonJSDocComment(newStart, lastComment.GetEndPosition(), &contents)
		node.GetLastChild().SetTrailingNonJSDocComment(newlineComment)
		return
	}

	blankLines := int(lastComment.GetStartPosition().Line - currentComment.GetEndPosition().Line)
	numWhiteSpace := 0
	if blankLines == 0 {
		numWhiteSpace = int(lastComment.GetStartPosition().Column - currentComment.GetEndPosition().Column - 1)
	}

	var comment strings.Builder
	comment.WriteString(currentComment.GetCommentString())
	for i := 0; i < numWhiteSpace; i++ {
		comment.WriteString(" ")
	}
	for i := 0; i < blankLines; i++ {
		comment.WriteString("\n")
	}
	comment.WriteString(lastComment.GetCommentString())

	contents := comment.String()
	allComments := rhino.NewNonJSDocComment(
		currentComment.GetStartPosition(), lastComment.GetEndPosition(), &contents)
	allComments.SetEndsAsLineComment(lastComment.IsEndingAsLineComment())
	allComments.SetIsInline(true)
	node.GetLastChild().SetTrailingNonJSDocComment(allComments)
}

// Java: IRFactory.TransformDispatcher.processOptChainFunctionCall(OptChainCallExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessOptChainFunctionCall(callNode *parser.OptChainCallExpressionTree) *rhino.Node {
	d.MaybeWarnForFeature(callNode, parser.FeatureSet_Feature_OPTIONAL_CHAINING)
	node := d.NewNodeWithChild(rhino.Token_OPTCHAIN_CALL, d.Transform(callNode.Operand))
	node.SetTrailingComma(callNode.HasTrailingComma)
	for _, child := range callNode.Arguments.Arguments {
		node.AddChildToBack(d.Transform(child))
	}
	node.SetIsOptionalChainStart(callNode.IsStartOfOptionalChain)
	d.AnnotateCalls(node)
	return node
}

// Java: IRFactory.TransformDispatcher.processFunction(FunctionDeclarationTree)
func (d *IRFactory_TransformDispatcher) ProcessFunction(functionTree *parser.FunctionDeclarationTree) *rhino.Node {
	isDeclaration := functionTree.Kind == parser.FunctionDeclarationTree_Kind_DECLARATION
	isMember := functionTree.Kind == parser.FunctionDeclarationTree_Kind_MEMBER
	isArrow := functionTree.Kind == parser.FunctionDeclarationTree_Kind_ARROW
	isAsync := functionTree.IsAsync
	isGenerator := functionTree.IsGenerator
	isSignature := functionTree.FunctionBody.GetParseTreeBase().Type == parser.ParseTreeType_EMPTY_STATEMENT

	identifierType := IRFactory_TransformDispatcher_IdentifierType_STANDARD
	if functionTree.IsClassMember {
		identifierType = IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE
	}

	if isGenerator {
		d.MaybeWarnForFeature(functionTree, parser.FeatureSet_Feature_GENERATORS)
	}

	if isMember {
		d.MaybeWarnForFeature(functionTree, parser.FeatureSet_Feature_MEMBER_DECLARATIONS)
	}

	if isArrow {
		d.MaybeWarnForFeature(functionTree, parser.FeatureSet_Feature_ARROW_FUNCTIONS)
	}

	if isAsync {
		d.MaybeWarnForFeature(functionTree, parser.FeatureSet_Feature_ASYNC_FUNCTIONS)
	}

	if isGenerator && isAsync {
		d.MaybeWarnForFeature(functionTree, parser.FeatureSet_Feature_ASYNC_GENERATORS)
	}

	name := functionTree.Name
	var newName *rhino.Node
	if name != nil {
		newName = d.ProcessNameWithInlineCommentsOfToken(name, identifierType)
	} else {
		if isDeclaration || isMember {
			d.errorReporter.Error(
				"unnamed function statement", d.sourceName, IRFactory_Lineno(functionTree), IRFactory_Charno(functionTree))

			// Return the bare minimum to put the AST in a valid state.
			newName = d.CreateMissingNameNode()
		} else {
			newName = d.NewStringNodeWithToken(rhino.Token_NAME, "")
		}

		// Old Rhino tagged the empty name node with the line number of the
		// declaration.
		d.SetSourceInfo(newName, functionTree)
	}

	node := d.NewNode(rhino.Token_FUNCTION)
	if isMember {
		newName.SetString("")
	}
	node.AddChildToBack(newName)

	node.AddChildToBack(d.Transform(functionTree.FormalParameterList))

	bodyNode := d.Transform(functionTree.FunctionBody)
	if !isArrow && !isSignature && !bodyNode.IsBlock() {
		// When in "keep going" mode the parser tries to parse some constructs
		// the compiler doesn't support, repair it here.
		jutil.CheckState(d.config.RunMode() == Config_RunMode_KEEP_GOING)
		bodyNode = rhino.IR_Block()
	}
	d.ParseDirectives(bodyNode)
	node.AddChildToBack(bodyNode)

	node.SetIsGeneratorFunction(isGenerator)
	node.SetIsArrowFunction(isArrow)
	node.SetIsAsyncFunction(isAsync)
	node.PutBooleanProp(rhino.Node_OPT_ES6_TYPED, functionTree.IsOptional)

	d.AttachPossibleTrailingComment(node, functionTree.GetEnd())

	var result *rhino.Node

	if isMember {
		d.SetSourceInfo(node, functionTree)
		member := d.NewStringNodeFromIdentifier(rhino.Token_MEMBER_FUNCTION_DEF, identifierType, name)
		member.AddChildToBack(node)
		member.SetStaticMember(functionTree.IsStatic)
		// The source info should only include the identifier, not the entire
		// function expression
		d.SetSourceInfoFromToken(member, name)
		result = member
	} else {
		result = node
	}

	return result
}

// ProcessField processes class fields: general case.
//
// Java: IRFactory.TransformDispatcher.processField(FieldDeclarationTree)
func (d *IRFactory_TransformDispatcher) ProcessField(tree *parser.FieldDeclarationTree) *rhino.Node {
	if tree.Name.IsPrivateIdentifier() {
		d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_PRIVATE_ELEMENTS)
	} else {
		d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)
	}

	node := d.MaybeAddNonJsDocComment(
		d.NewStringNodeFromIdentifier(
			rhino.Token_MEMBER_FIELD_DEF, IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE, tree.Name),
		tree.GetStart())
	if tree.Initializer != nil {
		initializer := d.Transform(tree.Initializer)
		node.AddChildToBack(initializer)
		d.SetLength(node, tree.Location.Start, tree.Location.End)
	}

	node.PutBooleanProp(rhino.Node_STATIC_MEMBER, tree.IsStatic)

	return node
}

// ProcessComputedPropertyField processes class fields: computed property
// fields.
//
// Java: IRFactory.TransformDispatcher.processComputedPropertyField(ComputedPropertyFieldTree)
func (d *IRFactory_TransformDispatcher) ProcessComputedPropertyField(tree *parser.ComputedPropertyFieldTree) *rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_PUBLIC_CLASS_FIELDS)

	key := d.Transform(tree.Property)

	var node *rhino.Node
	if tree.Initializer != nil {
		node = d.NewNodeWith2Children(rhino.Token_COMPUTED_FIELD_DEF, key, d.Transform(tree.Initializer))
	} else {
		node = d.NewNodeWithChild(rhino.Token_COMPUTED_FIELD_DEF, key)
	}

	node.PutBooleanProp(rhino.Node_STATIC_MEMBER, tree.IsStatic)

	return node
}

// Java: IRFactory.TransformDispatcher.processFormalParameterList(FormalParameterListTree)
func (d *IRFactory_TransformDispatcher) ProcessFormalParameterList(tree *parser.FormalParameterListTree) *rhino.Node {
	params := d.NewNode(rhino.Token_PARAM_LIST)
	params.SetTrailingComma(tree.HasTrailingComma)
	if !d.CheckParameters(tree) {
		return params
	}

	parameters := tree.Parameters
	argCount := 0

	for _, param := range parameters {
		var paramNode *rhino.Node
		switch param.GetParseTreeBase().Type {
		case parser.ParseTreeType_DEFAULT_PARAMETER:
			// processDefaultParameter() knows how to find and apply inline
			// JSDoc to the right node
			paramNode = d.ProcessDefaultParameter(param.(*parser.DefaultParameterTree))
		case parser.ParseTreeType_ITER_REST:
			d.MaybeWarnForFeature(param, parser.FeatureSet_Feature_REST_PARAMETERS)
			paramNode = d.TransformNodeWithInlineComments(param)
		default:
			paramNode = d.TransformNodeWithInlineComments(param)
			// Reusing the logic to attach trailing comments used from
			// call-site argsList
			d.attachPossibleTrailingCommentsForArgAt(
				paramNode, parameters, argCount, tree.CommaPositions, tree.Location.End)
		}

		// Children must be simple names, default parameters, rest
		// parameters, or destructuring patterns.
		jutil.CheckState(
			paramNode.IsName() ||
				paramNode.IsRest() ||
				paramNode.IsArrayPattern() ||
				paramNode.IsObjectPattern() ||
				paramNode.IsDefaultValue())
		params.AddChildToBack(paramNode)
		argCount++
	}

	return params
}

// Java: IRFactory.TransformDispatcher.processDefaultParameter(DefaultParameterTree)
func (d *IRFactory_TransformDispatcher) ProcessDefaultParameter(tree *parser.DefaultParameterTree) *rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_DEFAULT_PARAMETERS)
	targetTree := tree.Lhs
	var targetNode *rhino.Node
	if targetTree.GetParseTreeBase().Type == parser.ParseTreeType_IDENTIFIER_EXPRESSION {
		// allow inline JSDoc on an identifier
		// let { /** inlineType */ x = defaultValue } = someObject;
		// TODO(bradfordcsmith): Do we need to allow inline JSDoc for
		// qualified names, too?
		targetNode = d.ProcessNameWithInlineComments(
			targetTree.(*parser.IdentifierExpressionTree), IRFactory_TransformDispatcher_IdentifierType_STANDARD)
	} else {
		// ({prop: /** string */ ns.a.b = 'foo'} = someObject);
		// NOTE: CheckJSDoc will report an error for this case, since we want
		// qualified names to be declared with individual statements, like
		// `/** @type {string} */ ns.a.b;`
		targetNode = d.TransformNodeWithInlineComments(targetTree)
	}
	defaultValueExpression := d.Transform(tree.DefaultValue)
	defaultValueNode := d.NewNodeWith2Children(rhino.Token_DEFAULT_VALUE, targetNode, defaultValueExpression)
	d.ReportErrorIfYieldOrAwaitInDefaultValue(defaultValueNode)
	d.SetSourceInfo(defaultValueNode, tree)
	return defaultValueNode
}

// Java: IRFactory.TransformDispatcher.processIterRest(IterRestTree)
func (d *IRFactory_TransformDispatcher) ProcessIterRest(tree *parser.IterRestTree) *rhino.Node {
	target := d.TransformNodeWithInlineComments(tree.AssignmentTarget)
	return d.NewNodeWithChild(rhino.Token_ITER_REST, target)
}

// Java: IRFactory.TransformDispatcher.processIterSpread(IterSpreadTree)
func (d *IRFactory_TransformDispatcher) ProcessIterSpread(tree *parser.IterSpreadTree) *rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_SPREAD_EXPRESSIONS)

	return d.NewNodeWithChild(rhino.Token_ITER_SPREAD, d.Transform(tree.Expression))
}

// Java: IRFactory.TransformDispatcher.processObjectSpread(ObjectSpreadTree)
func (d *IRFactory_TransformDispatcher) ProcessObjectSpread(tree *parser.ObjectSpreadTree) *rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_OBJECT_LITERALS_WITH_SPREAD)

	return d.NewNodeWithChild(rhino.Token_OBJECT_SPREAD, d.Transform(tree.Expression))
}

// Java: IRFactory.TransformDispatcher.processIfStatement(IfStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessIfStatement(statementNode *parser.IfStatementTree) *rhino.Node {
	node := d.NewNode(rhino.Token_IF)
	node.AddChildToBack(d.Transform(statementNode.Condition))
	node.AddChildToBack(d.TransformBlock(statementNode.IfClause))
	if statementNode.ElseClause != nil {
		node.AddChildToBack(d.TransformBlock(statementNode.ElseClause))
	}
	return node
}

// Java: IRFactory.TransformDispatcher.markBinaryExpressionFeatures(BinaryOperatorTree)
func (d *IRFactory_TransformDispatcher) MarkBinaryExpressionFeatures(exprNode *parser.BinaryOperatorTree) {
	switch exprNode.Operator.GetTokenBase().Type {
	case parser.TokenType_STAR_STAR, parser.TokenType_STAR_STAR_EQUAL:
		d.MaybeWarnForFeature(exprNode, parser.FeatureSet_Feature_EXPONENT_OP)
	case parser.TokenType_QUESTION_QUESTION:
		d.MaybeWarnForFeature(exprNode, parser.FeatureSet_Feature_NULL_COALESCE_OP)
	case parser.TokenType_QUESTION_QUESTION_EQUAL:
		d.MaybeWarnForFeature(exprNode, parser.FeatureSet_Feature_NULL_COALESCE_OP)
		fallthrough
	case parser.TokenType_OR_EQUAL, parser.TokenType_AND_EQUAL:
		d.MaybeWarnForFeature(exprNode, parser.FeatureSet_Feature_LOGICAL_ASSIGNMENT)
	default:
	}
}

// Java: IRFactory.TransformDispatcher.processBinaryExpression(BinaryOperatorTree)
func (d *IRFactory_TransformDispatcher) ProcessBinaryExpression(exprNode *parser.BinaryOperatorTree) *rhino.Node {
	if d.jsdocTracker.HasPendingCommentBefore(exprNode.Right.GetParseTreeBase().Location.Start) ||
		d.nonJsdocTracker.HasPendingCommentBefore(exprNode.Right.GetParseTreeBase().Location.Start) {
		d.MarkBinaryExpressionFeatures(exprNode)
		return d.NewNodeWith2Children(
			IRFactory_TransformBinaryTokenType(exprNode.Operator.GetTokenBase().Type),
			d.TransformLhsOfBinaryTree(exprNode),
			d.Transform(exprNode.Right))
	} else {
		// No pending comments, we can traverse out of order.
		return d.ProcessBinaryExpressionHelper(exprNode)
	}
}

// ProcessBinaryExpressionHelper: deep binary ops (typical string
// concatentations) can cause stack overflows, avoid recursing in this case
// and loop instead.
//
// Java: IRFactory.TransformDispatcher.processBinaryExpressionHelper(BinaryOperatorTree)
func (d *IRFactory_TransformDispatcher) ProcessBinaryExpressionHelper(exprTree *parser.BinaryOperatorTree) *rhino.Node {
	var root *rhino.Node
	var current *rhino.Node
	var previous *rhino.Node
	for exprTree != nil {
		d.MarkBinaryExpressionFeatures(exprTree)
		previous = current
		// Skip the first child but recurse normally into the right operand as
		// typically this isn't deep and because we have already checked that
		// there isn't any JSDoc we can traverse out of order.
		current = d.NewNodeWithChild(
			IRFactory_TransformBinaryTokenType(exprTree.Operator.GetTokenBase().Type), d.Transform(exprTree.Right))
		// We have inlined "transform" here. Normally, we would need to handle
		// the JSDoc here but we know there is no JSDoc to attach, which
		// simplifies things.
		d.SetSourceInfo(current, exprTree)

		// As the iteration continues add the left operand.
		if previous != nil {
			previous.AddChildToFront(current)
		}

		if binaryOperatorTree, ok := exprTree.Left.(*parser.BinaryOperatorTree); ok {
			// continue with the left hand child
			exprTree = binaryOperatorTree
		} else {
			// Finish things off, add the left operand to the current node.
			leftNode := d.TransformLhsOfBinaryTree(exprTree)
			current.AddChildToFront(leftNode)
			// Nothing left to do.
			exprTree = nil
		}

		// Save the top binary op, this is the result to return.
		if root == nil {
			root = current
		}
	}
	return root
}

// Java: IRFactory.TransformDispatcher.transformLhsOfBinaryTree(BinaryOperatorTree)
func (d *IRFactory_TransformDispatcher) TransformLhsOfBinaryTree(exprTree *parser.BinaryOperatorTree) *rhino.Node {
	_, leftIsIdentifierExpression := exprTree.Left.(*parser.IdentifierExpressionTree)
	lhsOfInIsPrivateId := exprTree.Operator.GetTokenBase().Type == parser.TokenType_IN &&
		leftIsIdentifierExpression &&
		exprTree.Left.(*parser.IdentifierExpressionTree).IdentifierToken.IsPrivateIdentifier()

	ignored := d.privateIdLhsOfInScope.MaybeIncrement(lhsOfInIsPrivateId)
	defer ignored.Close()
	return d.Transform(exprTree.Left)
}

// Java: IRFactory.TransformDispatcher.processDebuggerStatement(DebuggerStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessDebuggerStatement(unused *parser.DebuggerStatementTree) *rhino.Node {
	return d.NewNode(rhino.Token_DEBUGGER)
}

// Java: IRFactory.TransformDispatcher.processThisExpression(ThisExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessThisExpression(unused *parser.ThisExpressionTree) *rhino.Node {
	return d.NewNode(rhino.Token_THIS)
}

// Java: IRFactory.TransformDispatcher.processLabeledStatement(LabelledStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessLabeledStatement(labelTree *parser.LabelledStatementTree) *rhino.Node {
	statement := d.Transform(labelTree.Statement)
	if statement.IsFunction() ||
		statement.IsClass() ||
		statement.IsLet() ||
		statement.IsConst() {
		d.errorReporter.Error(
			"Lexical declarations are only allowed at top level or inside a block.",
			d.sourceName,
			IRFactory_Lineno(labelTree),
			IRFactory_Charno(labelTree))
		return statement // drop the LABEL node so that the resulting AST is valid
	}
	return d.NewNodeWith2Children(rhino.Token_LABEL, d.TransformLabelName(labelTree.Name), statement)
}

// Java: IRFactory.TransformDispatcher.processName(IdentifierExpressionTree, IdentifierType)
func (d *IRFactory_TransformDispatcher) ProcessName(nameNode *parser.IdentifierExpressionTree, identifierType IRFactory_TransformDispatcher_IdentifierType) *rhino.Node {
	return d.ProcessNameOfToken(nameNode.IdentifierToken, rhino.Token_NAME, identifierType)
}

// Java: IRFactory.TransformDispatcher.processName(IdentifierToken, Token, IdentifierType)
func (d *IRFactory_TransformDispatcher) ProcessNameOfToken(
	identifierToken *parser.IdentifierToken, output rhino.Token, identifierType IRFactory_TransformDispatcher_IdentifierType) *rhino.Node {
	comment := d.ParseNonJSDocCommentAt(identifierToken.GetStart(), true)

	node := d.NewStringNodeFromIdentifier(output, identifierType, identifierToken)

	if output == rhino.Token_NAME {
		d.MaybeWarnReservedKeyword(identifierToken)

		info := d.ParseJSDocInfoOnToken(identifierToken)
		if info != nil {
			node.SetJSDocInfo(info)
		}
	}

	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	d.SetSourceInfoFromToken(node, identifierToken)
	return node
}

// Java: IRFactory.TransformDispatcher.processString(LiteralToken)
func (d *IRFactory_TransformDispatcher) ProcessString(token parser.LiteralToken) *rhino.Node {
	jutil.CheckArgument(token.GetTokenBase().Type == parser.TokenType_STRING)
	comment := d.ParseNonJSDocCommentAt(token.GetStart(), true)

	node := d.NewStringNodeWithToken(rhino.Token_STRINGLIT, d.NormalizeString(token, false))

	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	d.SetSourceInfoFromToken(node, token)
	return node
}

// Java: IRFactory.TransformDispatcher.processTemplateLiteralToken(TemplateLiteralToken)
func (d *IRFactory_TransformDispatcher) ProcessTemplateLiteralToken(token *parser.TemplateLiteralToken) *rhino.Node {
	jutil.CheckArgument(
		token.Type == parser.TokenType_NO_SUBSTITUTION_TEMPLATE ||
			token.Type == parser.TokenType_TEMPLATE_HEAD ||
			token.Type == parser.TokenType_TEMPLATE_MIDDLE ||
			token.Type == parser.TokenType_TEMPLATE_TAIL)
	var node *rhino.Node
	if token.HasError() {
		node = d.NewTemplateLitStringNode(nil, token.Value)
	} else {
		cooked := d.NormalizeString(token, true)
		node = d.NewTemplateLitStringNode(&cooked, token.Value)
	}
	d.SetSourceInfoFromToken(node, token)
	return node
}

// Java: IRFactory.TransformDispatcher.processNameWithInlineComments(IdentifierExpressionTree, IdentifierType)
func (d *IRFactory_TransformDispatcher) ProcessNameWithInlineComments(
	identifierExpression *parser.IdentifierExpressionTree, identifierType IRFactory_TransformDispatcher_IdentifierType) *rhino.Node {
	return d.ProcessNameWithInlineCommentsOfToken(identifierExpression.IdentifierToken, identifierType)
}

// Java: IRFactory.TransformDispatcher.processNameWithInlineComments(IdentifierToken, IdentifierType)
func (d *IRFactory_TransformDispatcher) ProcessNameWithInlineCommentsOfToken(
	identifierToken *parser.IdentifierToken, identifierType IRFactory_TransformDispatcher_IdentifierType) *rhino.Node {
	info := d.ParseInlineJSDocAt(identifierToken.GetStart())
	comment := d.ParseNonJSDocCommentAt(identifierToken.GetStart(), false)

	d.MaybeWarnReservedKeyword(identifierToken)
	node := d.NewStringNodeFromIdentifier(rhino.Token_NAME, identifierType, identifierToken)

	if info != nil {
		node.SetJSDocInfo(info)
	}
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	d.SetSourceInfoFromToken(node, identifierToken)
	return node
}

// Java: IRFactory.TransformDispatcher.maybeWarnKeywordProperty(Node)
func (d *IRFactory_TransformDispatcher) MaybeWarnKeywordProperty(node *rhino.Node) {
	if d.currentFileIsExterns {
		return
	}
	if !rhino.TokenStream_IsKeyword(node.GetString()) {
		return
	}

	d.features = d.features.With(parser.FeatureSet_Feature_KEYWORDS_AS_PROPERTIES)
	if d.config.LanguageMode() == Config_LanguageMode_ECMASCRIPT3 {
		d.errorReporter.Warning(
			IRFactory_INVALID_ES3_PROP_NAME, d.sourceName, node.GetLineno(), node.GetCharno())
	}
}

// Java: IRFactory.TransformDispatcher.maybeWarnReservedKeyword(IdentifierToken)
func (d *IRFactory_TransformDispatcher) MaybeWarnReservedKeyword(token *parser.IdentifierToken) {
	// We're just checking for keywords here, not actually using the
	// identifier so we don't care about potentially running into a misplaced
	// private identifier.
	identifier := token.GetMaybePrivateValue()
	isIdentifier := false
	if rhino.TokenStream_IsKeyword(identifier) {
		d.features = d.features.With(parser.FeatureSet_Feature_ES3_KEYWORDS_AS_IDENTIFIERS)
		isIdentifier = d.config.LanguageMode() == Config_LanguageMode_ECMASCRIPT3
	}
	if d.reservedKeywords != nil && d.reservedKeywords.Contains(identifier) {
		d.features = d.features.With(parser.FeatureSet_Feature_KEYWORDS_AS_PROPERTIES)
		isIdentifier = d.config.LanguageMode() == Config_LanguageMode_ECMASCRIPT3
	}
	if isIdentifier {
		d.errorReporter.Error(
			"identifier is a reserved word",
			d.sourceName,
			IRFactory_LinenoOfPosition(token.Location.Start),
			IRFactory_CharnoOfPosition(token.Location.Start))
	}
}

// Java: IRFactory.TransformDispatcher.processNewExpression(NewExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessNewExpression(exprNode *parser.NewExpressionTree) *rhino.Node {
	node := d.NewNodeWithChild(rhino.Token_NEW, d.Transform(exprNode.Operand))
	node.SetTrailingComma(exprNode.HasTrailingComma)
	if exprNode.Arguments != nil {
		for _, arg := range exprNode.Arguments.Arguments {
			node.AddChildToBack(d.Transform(arg))
		}
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processNumberLiteral(LiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessNumberLiteral(literalNode *parser.LiteralExpressionTree) *rhino.Node {
	value := d.NormalizeNumber(literalNode.LiteralToken.(parser.LiteralToken))
	number := d.NewNumberNode(value)
	d.SetSourceInfo(number, literalNode)
	return number
}

// Java: IRFactory.TransformDispatcher.processBigIntLiteral(LiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessBigIntLiteral(literalNode *parser.LiteralExpressionTree) *rhino.Node {
	d.MaybeWarnForFeature(literalNode, parser.FeatureSet_Feature_BIGINT)
	value := d.NormalizeBigInt(literalNode.LiteralToken.(parser.LiteralToken))
	bigint := d.NewBigIntNode(value)
	d.SetSourceInfo(bigint, literalNode)
	return bigint
}

// Java: IRFactory.TransformDispatcher.processObjectLiteral(ObjectLiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessObjectLiteral(objTree *parser.ObjectLiteralExpressionTree) *rhino.Node {
	node := d.NewNode(rhino.Token_OBJECTLIT)
	node.SetTrailingComma(objTree.HasTrailingComma)
	maybeWarn := false
	for _, el := range objTree.PropertyNameAndValues {
		if el.GetParseTreeBase().Type == parser.ParseTreeType_DEFAULT_PARAMETER {
			// (e.g. var o = { x=4 };) This is only parsed for compatibility
			// with object patterns.
			d.errorReporter.Error(
				"Default value cannot appear at top level of an object literal.",
				d.sourceName,
				IRFactory_Lineno(el),
				0)
			continue
		} else if el.GetParseTreeBase().Type == parser.ParseTreeType_GET_ACCESSOR && d.MaybeReportGetter(el) {
			continue
		} else if el.GetParseTreeBase().Type == parser.ParseTreeType_SET_ACCESSOR && d.MaybeReportSetter(el) {
			continue
		}

		key := d.Transform(el)
		if !key.IsComputedProp() &&
			!key.IsQuotedStringKey() &&
			!key.IsSpread() &&
			!d.currentFileIsExterns {
			d.MaybeWarnKeywordProperty(key)
		}
		if key.IsShorthandProperty() {
			maybeWarn = true
		}

		node.AddChildToBack(key)
	}
	if maybeWarn {
		d.MaybeWarnForFeature(objTree, parser.FeatureSet_Feature_SHORTHAND_OBJECT_PROPERTIES)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processComputedPropertyDefinition(ComputedPropertyDefinitionTree)
func (d *IRFactory_TransformDispatcher) ProcessComputedPropertyDefinition(tree *parser.ComputedPropertyDefinitionTree) *rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)

	return d.NewNodeWith2Children(rhino.Token_COMPUTED_PROP, d.Transform(tree.Property), d.Transform(tree.Value))
}

// Java: IRFactory.TransformDispatcher.processComputedPropertyMethod(ComputedPropertyMethodTree)
func (d *IRFactory_TransformDispatcher) ProcessComputedPropertyMethod(tree *parser.ComputedPropertyMethodTree) *rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)

	n := d.NewNodeWith2Children(rhino.Token_COMPUTED_PROP, d.Transform(tree.Property), d.Transform(tree.Method))
	n.PutBooleanProp(rhino.Node_COMPUTED_PROP_METHOD, true)
	if tree.Method.(*parser.FunctionDeclarationTree).IsStatic {
		n.SetStaticMember(true)
	}
	return n
}

// Java: IRFactory.TransformDispatcher.processComputedPropertyGetter(ComputedPropertyGetterTree)
func (d *IRFactory_TransformDispatcher) ProcessComputedPropertyGetter(tree *parser.ComputedPropertyGetterTree) *rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)

	key := d.Transform(tree.Property)
	body := d.Transform(tree.Body)
	function := rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(), body)
	function.SrcrefTreeIfMissing(body)
	n := d.NewNodeWith2Children(rhino.Token_COMPUTED_PROP, key, function)
	n.PutBooleanProp(rhino.Node_COMPUTED_PROP_GETTER, true)
	n.PutBooleanProp(rhino.Node_STATIC_MEMBER, tree.IsStatic)
	return n
}

// Java: IRFactory.TransformDispatcher.processComputedPropertySetter(ComputedPropertySetterTree)
func (d *IRFactory_TransformDispatcher) ProcessComputedPropertySetter(tree *parser.ComputedPropertySetterTree) *rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_COMPUTED_PROPERTIES)

	key := d.Transform(tree.Property)

	paramList := d.ProcessFormalParameterList(tree.Parameter)
	d.SetSourceInfo(paramList, tree.Parameter)

	body := d.Transform(tree.Body)

	function := rhino.IR_Function(rhino.IR_Name(""), paramList, body)
	function.SrcrefTreeIfMissing(body)

	n := d.NewNodeWith2Children(rhino.Token_COMPUTED_PROP, key, function)
	n.PutBooleanProp(rhino.Node_COMPUTED_PROP_SETTER, true)
	n.PutBooleanProp(rhino.Node_STATIC_MEMBER, tree.IsStatic)
	return n
}

// Java: IRFactory.TransformDispatcher.processGetAccessor(GetAccessorTree)
func (d *IRFactory_TransformDispatcher) ProcessGetAccessor(tree *parser.GetAccessorTree) *rhino.Node {
	identifierType := IRFactory_TransformDispatcher_IdentifierType_STANDARD
	if tree.IsClassMember {
		identifierType = IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE
	}
	key := d.ProcessObjectLitKey(tree.PropertyName, rhino.Token_GETTER_DEF, identifierType)
	body := d.Transform(tree.Body)
	dummyName := d.NewStringNodeWithToken(rhino.Token_NAME, "")
	d.SetSourceInfo(dummyName, tree.Body)
	paramList := d.NewNode(rhino.Token_PARAM_LIST)
	d.SetSourceInfo(paramList, tree.Body)
	value := d.NewNodeWith3Children(rhino.Token_FUNCTION, dummyName, paramList, body)
	d.SetSourceInfo(value, tree.Body)
	key.AddChildToFront(value)
	key.SetStaticMember(tree.IsStatic)
	return key
}

// Java: IRFactory.TransformDispatcher.processSetAccessor(SetAccessorTree)
func (d *IRFactory_TransformDispatcher) ProcessSetAccessor(tree *parser.SetAccessorTree) *rhino.Node {
	identifierType := IRFactory_TransformDispatcher_IdentifierType_STANDARD
	if tree.IsClassMember {
		identifierType = IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE
	}
	key := d.ProcessObjectLitKey(tree.PropertyName, rhino.Token_SETTER_DEF, identifierType)

	paramList := d.ProcessFormalParameterList(tree.Parameter)
	d.SetSourceInfo(paramList, tree.Parameter)

	body := d.Transform(tree.Body)

	dummyName := d.NewStringNodeWithToken(rhino.Token_NAME, "")
	d.SetSourceInfoFromToken(dummyName, tree.PropertyName)

	value := d.NewNodeWith3Children(rhino.Token_FUNCTION, dummyName, paramList, body)
	d.SetSourceInfo(value, tree.Body)

	key.AddChildToFront(value)
	key.SetStaticMember(tree.IsStatic)
	return key
}

// Java: IRFactory.TransformDispatcher.processPropertyNameAssignment(PropertyNameAssignmentTree)
func (d *IRFactory_TransformDispatcher) ProcessPropertyNameAssignment(tree *parser.PropertyNameAssignmentTree) *rhino.Node {
	key := d.ProcessObjectLitKey(tree.Name, rhino.Token_STRING_KEY, IRFactory_TransformDispatcher_IdentifierType_STANDARD)
	if tree.Value != nil {
		key.AddChildToFront(d.Transform(tree.Value))
	} else {
		value := d.MaybeAddNonJsDocComment(
			d.NewStringNodeWithToken(rhino.Token_NAME, key.GetString()), tree.Name.GetStart()).
			Srcref(key)
		key.SetShorthandProperty(true)
		key.AddChildToFront(value)
	}
	return key
}

// Java: IRFactory.TransformDispatcher.checkParenthesizedExpression(ParenExpressionTree)
func (d *IRFactory_TransformDispatcher) CheckParenthesizedExpression(exprNode *parser.ParenExpressionTree) {
	if exprNode.Expression.GetParseTreeBase().Type == parser.ParseTreeType_COMMA_EXPRESSION {
		commaNodes := exprNode.Expression.(*parser.CommaExpressionTree).Expressions
		lastChild := commaNodes[len(commaNodes)-1]
		if lastChild.GetParseTreeBase().Type == parser.ParseTreeType_ITER_REST {
			d.errorReporter.Error(
				"A rest parameter must be in a parameter list.",
				d.sourceName,
				IRFactory_Lineno(lastChild),
				IRFactory_Charno(lastChild))
		}
	}
}

// Java: IRFactory.TransformDispatcher.processParenthesizedExpression(ParenExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessParenthesizedExpression(exprNode *parser.ParenExpressionTree) *rhino.Node {
	d.CheckParenthesizedExpression(exprNode)
	expr := d.Transform(exprNode.Expression)
	expr.SetIsParenthesized(true)
	return expr
}

// Java: IRFactory.TransformDispatcher.processPropertyGet(MemberExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessPropertyGet(getNode *parser.MemberExpressionTree) *rhino.Node {
	leftChild := d.Transform(getNode.Operand)
	propName := getNode.MemberName
	if propName == nil {
		return leftChild
	}

	// Private element get prop not valid on super receiver.
	identifierType := IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE
	if leftChild.IsSuper() {
		identifierType = IRFactory_TransformDispatcher_IdentifierType_STANDARD
	}
	getProp := d.MaybeAddNonJsDocComment(
		d.NewStringNodeFromIdentifier(rhino.Token_GETPROP, identifierType, propName),
		getNode.MemberName.GetStart())
	getProp.AddChildToBack(leftChild)
	d.SetSourceInfoFromToken(getProp, propName)
	d.MaybeWarnKeywordProperty(getProp)
	return getProp
}

// Java: IRFactory.TransformDispatcher.processOptChainPropertyGet(OptionalMemberExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessOptChainPropertyGet(getNode *parser.OptionalMemberExpressionTree) *rhino.Node {
	d.MaybeWarnForFeature(getNode, parser.FeatureSet_Feature_OPTIONAL_CHAINING)

	leftChild := d.Transform(getNode.Operand)
	propName := getNode.MemberName
	if propName == nil {
		return leftChild
	}

	// Private element get prop not valid on super receiver.
	identifierType := IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE
	if leftChild.IsSuper() {
		identifierType = IRFactory_TransformDispatcher_IdentifierType_STANDARD
	}
	getProp := d.MaybeAddNonJsDocComment(
		d.NewStringNodeFromIdentifier(rhino.Token_OPTCHAIN_GETPROP, identifierType, propName),
		getNode.MemberName.GetStart())
	getProp.AddChildToBack(leftChild)
	getProp.SetIsOptionalChainStart(getNode.IsStartOfOptionalChain)
	d.SetSourceInfoFromToken(getProp, propName)
	d.MaybeWarnKeywordProperty(getProp)
	return getProp
}

// Java: IRFactory.TransformDispatcher.processRegExpLiteral(LiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessRegExpLiteral(literalTree *parser.LiteralExpressionTree) *rhino.Node {
	token := literalTree.LiteralToken.(parser.LiteralToken)
	literalStringNode := d.NewStringNode(d.NormalizeRegex(token))
	// TODO(johnlenz): fix the source location.
	d.SetSourceInfoFromToken(literalStringNode, token)
	node := d.NewNodeWithChild(rhino.Token_REGEXP, literalStringNode)

	rawRegex := token.GetLiteralTokenBase().Value
	// Go-only: language-adapter — '/' is ASCII, so the byte index of the last
	// '/' delimits the same substring as Java's UTF-16 index.
	lastSlash := strings.LastIndexByte(rawRegex, '/')
	flags := ""
	if lastSlash < len(rawRegex) {
		flags = rawRegex[lastSlash+1:]
	}
	d.ValidateRegExpFlags(literalTree, flags)

	if flags != "" {
		flagsNode := d.NewStringNode(flags)
		// TODO(johnlenz): fix the source location.
		d.SetSourceInfoFromToken(flagsNode, token)
		node.AddChildToBack(flagsNode)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.validateRegExpFlags(LiteralExpressionTree, String)
func (d *IRFactory_TransformDispatcher) ValidateRegExpFlags(tree *parser.LiteralExpressionTree, flags string) {
	for _, flag := range jutil.StringToUTF16(flags) {
		switch flag {
		case 'g', 'i', 'm':
		case 'u', 'y':
			feature := parser.FeatureSet_Feature_REGEXP_FLAG_Y
			if flag == 'u' {
				feature = parser.FeatureSet_Feature_REGEXP_FLAG_U
			}
			d.MaybeWarnForFeature(tree, feature)
		case 's':
			d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_REGEXP_FLAG_S)
		case 'd':
			d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_REGEXP_FLAG_D)
		default:
			d.errorReporter.Error(
				"Invalid RegExp flag '"+jutil.StringFromChar(flag)+"'", d.sourceName, IRFactory_Lineno(tree), IRFactory_Charno(tree))
		}
	}
}

// Java: IRFactory.TransformDispatcher.processReturnStatement(ReturnStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessReturnStatement(statementNode *parser.ReturnStatementTree) *rhino.Node {
	node := d.NewNode(rhino.Token_RETURN)
	if statementNode.Expression != nil {
		node.AddChildToBack(d.Transform(statementNode.Expression))
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processStringLiteral(LiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessStringLiteral(literalTree *parser.LiteralExpressionTree) *rhino.Node {
	token := literalTree.LiteralToken.(parser.LiteralToken)
	return d.ProcessString(token)
}

// Java: IRFactory.TransformDispatcher.processTemplateLiteral(TemplateLiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessTemplateLiteral(tree *parser.TemplateLiteralExpressionTree) *rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_TEMPLATE_LITERALS)
	templateLitNode := d.NewNode(rhino.Token_TEMPLATELIT)
	d.SetSourceInfo(templateLitNode, tree)
	var node *rhino.Node
	if tree.Operand == nil {
		node = templateLitNode
	} else {
		node = d.NewNodeWith2Children(rhino.Token_TAGGED_TEMPLATELIT, d.Transform(tree.Operand), templateLitNode)
	}
	for _, child := range tree.Elements {
		templateLitNode.AddChildToBack(d.Transform(child))
	}
	if node.IsTaggedTemplateLit() {
		d.AnnotateCalls(node)
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processTemplateLiteralPortion(TemplateLiteralPortionTree)
func (d *IRFactory_TransformDispatcher) ProcessTemplateLiteralPortion(tree *parser.TemplateLiteralPortionTree) *rhino.Node {
	return d.ProcessTemplateLiteralToken(tree.Value.(*parser.TemplateLiteralToken))
}

// Java: IRFactory.TransformDispatcher.processTemplateSubstitution(TemplateSubstitutionTree)
func (d *IRFactory_TransformDispatcher) ProcessTemplateSubstitution(tree *parser.TemplateSubstitutionTree) *rhino.Node {
	return d.NewNodeWithChild(rhino.Token_TEMPLATELIT_SUB, d.Transform(tree.Expression))
}

// Java: IRFactory.TransformDispatcher.processSwitchCase(CaseClauseTree)
func (d *IRFactory_TransformDispatcher) ProcessSwitchCase(caseNode *parser.CaseClauseTree) *rhino.Node {
	expr := caseNode.Expression
	node := d.NewNodeWithChild(rhino.Token_CASE, d.Transform(expr))
	block := d.NewNode(rhino.Token_BLOCK)
	block.SetIsAddedBlock(true)
	d.SetSourceInfo(block, caseNode)
	if caseNode.Statements != nil {
		for _, child := range caseNode.Statements {
			block.AddChildToBack(d.Transform(child))
		}
	}
	node.AddChildToBack(block)
	return node
}

// Java: IRFactory.TransformDispatcher.processSwitchDefault(DefaultClauseTree)
func (d *IRFactory_TransformDispatcher) ProcessSwitchDefault(caseNode *parser.DefaultClauseTree) *rhino.Node {
	node := d.NewNode(rhino.Token_DEFAULT_CASE)
	block := d.NewNode(rhino.Token_BLOCK)
	block.SetIsAddedBlock(true)
	d.SetSourceInfo(block, caseNode)
	if caseNode.Statements != nil {
		for _, child := range caseNode.Statements {
			block.AddChildToBack(d.Transform(child))
		}
	}
	node.AddChildToBack(block)
	return node
}

// Java: IRFactory.TransformDispatcher.processSwitchStatement(SwitchStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessSwitchStatement(statementNode *parser.SwitchStatementTree) *rhino.Node {
	node := d.NewNodeWithChild(rhino.Token_SWITCH, d.Transform(statementNode.Expression))
	switchBody := d.NewNode(rhino.Token_SWITCH_BODY)
	d.SetSourceInfo(switchBody, statementNode)
	node.AddChildToBack(switchBody)
	for _, child := range statementNode.CaseClauses {
		switchBody.AddChildToBack(d.Transform(child))
	}
	return node
}

// Java: IRFactory.TransformDispatcher.processThrowStatement(ThrowStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessThrowStatement(statementNode *parser.ThrowStatementTree) *rhino.Node {
	return d.NewNodeWithChild(rhino.Token_THROW, d.Transform(statementNode.Value))
}

// Java: IRFactory.TransformDispatcher.processTryStatement(TryStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessTryStatement(statementNode *parser.TryStatementTree) *rhino.Node {
	node := d.NewNodeWithChild(rhino.Token_TRY, d.TransformBlock(statementNode.Body))
	block := d.NewNode(rhino.Token_BLOCK)
	lineSet := false

	cc := statementNode.CatchBlock
	if cc != nil {
		// Mark the enclosing block at the same line as the first catch
		// clause.
		d.SetSourceInfo(block, cc)
		lineSet = true
		block.AddChildToBack(d.Transform(cc))
	}

	node.AddChildToBack(block)

	finallyBlock := statementNode.FinallyBlock
	if finallyBlock != nil {
		node.AddChildToBack(d.TransformBlock(finallyBlock))
	}

	// If we didn't set the line on the catch clause, then
	// we've got an empty catch clause.  Set its line to be the same
	// as the finally block (to match Old Rhino's behavior.)
	if !lineSet && (finallyBlock != nil) {
		d.SetSourceInfo(block, finallyBlock)
	}

	return node
}

// Java: IRFactory.TransformDispatcher.processCatchClause(CatchTree)
func (d *IRFactory_TransformDispatcher) ProcessCatchClause(clauseNode *parser.CatchTree) *rhino.Node {
	if clauseNode.Exception.GetParseTreeBase().Type == parser.ParseTreeType_EMPTY_STATEMENT {
		d.MaybeWarnForFeature(clauseNode, parser.FeatureSet_Feature_OPTIONAL_CATCH_BINDING)
	}

	return d.NewNodeWith2Children(
		rhino.Token_CATCH, d.Transform(clauseNode.Exception), d.TransformBlock(clauseNode.CatchBody))
}

// Java: IRFactory.TransformDispatcher.processFinally(FinallyTree)
func (d *IRFactory_TransformDispatcher) ProcessFinally(finallyNode *parser.FinallyTree) *rhino.Node {
	return d.TransformBlock(finallyNode.Block)
}

// Java: IRFactory.TransformDispatcher.processUnaryExpression(UnaryExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessUnaryExpression(exprNode *parser.UnaryExpressionTree) *rhino.Node {
	type_ := IRFactory_TransformUnaryTokenType(exprNode.Operator.GetTokenBase().Type)
	operand := d.Transform(exprNode.Operand)

	switch type_ {
	case rhino.Token_DELPROP:
		if !(operand.IsGetProp() ||
			operand.IsGetElem() ||
			operand.IsName() ||
			operand.IsOptChainGetProp() ||
			operand.IsOptChainGetElem()) {
			msg := "Invalid delete operand. Only properties can be deleted."
			d.errorReporter.Error(msg, d.sourceName, operand.GetLineno(), 0)
		}
	case rhino.Token_POS:
		if operand.IsBigInt() {
			d.errorReporter.Error(
				"Cannot convert a BigInt value to a number", d.sourceName, operand.GetLineno(), 0)
		}
	default:
	}

	return d.NewNodeWithChild(type_, operand)
}

// Java: IRFactory.TransformDispatcher.processUpdateExpression(UpdateExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessUpdateExpression(updateExpr *parser.UpdateExpressionTree) *rhino.Node {
	type_ := IRFactory_TransformUpdateTokenType(updateExpr.Operator.GetTokenBase().Type)
	operand := d.Transform(updateExpr.Operand)
	return d.CreateUpdateNode(
		type_, updateExpr.OperatorPosition == parser.UpdateExpressionTree_OperatorPosition_POSTFIX, operand)
}

// Java: IRFactory.TransformDispatcher.createUpdateNode(Token, boolean, Node)
func (d *IRFactory_TransformDispatcher) CreateUpdateNode(type_ rhino.Token, postfix bool, operand *rhino.Node) *rhino.Node {
	assignTarget := operand
	if operand.IsCast() {
		assignTarget = operand.GetFirstChild()
	}
	if !assignTarget.IsValidAssignmentTarget() {
		fix := "prefix"
		if postfix {
			fix = "postfix"
		}
		op := "decrement"
		if type_ == rhino.Token_INC {
			op = "increment"
		}
		d.errorReporter.Error(
			jutil.StringFormat("Invalid %s %s operand.", fix, op),
			d.sourceName,
			operand.GetLineno(),
			operand.GetCharno())
	}
	node := d.NewNodeWithChild(type_, operand)
	node.PutBooleanProp(rhino.Node_INCRDECR_PROP, postfix)
	return node
}

// Java: IRFactory.TransformDispatcher.processVariableStatement(VariableStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessVariableStatement(stmt *parser.VariableStatementTree) *rhino.Node {
	// TODO(moz): Figure out why we still need the special handling
	return d.transformDispatcher.Process(stmt.Declarations)
}

// Java: IRFactory.TransformDispatcher.processVariableDeclarationList(VariableDeclarationListTree)
func (d *IRFactory_TransformDispatcher) ProcessVariableDeclarationList(decl *parser.VariableDeclarationListTree) *rhino.Node {
	var declType rhino.Token
	switch decl.DeclarationType {
	case parser.TokenType_CONST:
		d.MaybeWarnForFeature(decl, parser.FeatureSet_Feature_CONST_DECLARATIONS)
		declType = rhino.Token_CONST
	case parser.TokenType_LET:
		d.MaybeWarnForFeature(decl, parser.FeatureSet_Feature_LET_DECLARATIONS)
		declType = rhino.Token_LET
	case parser.TokenType_VAR:
		declType = rhino.Token_VAR
	default:
		panic(jutil.NewIllegalStateException())
	}

	node := d.NewNode(declType)
	for _, child := range decl.Declarations {
		node.AddChildToBack(d.TransformNodeWithInlineComments(child))
	}

	d.AttachPossibleTrailingComment(node, decl.GetEnd())
	return node
}

// Java: IRFactory.TransformDispatcher.processVariableDeclaration(VariableDeclarationTree)
func (d *IRFactory_TransformDispatcher) ProcessVariableDeclaration(decl *parser.VariableDeclarationTree) *rhino.Node {
	node := d.TransformNodeWithInlineComments(decl.Lvalue)
	lhs := node
	if node.IsDestructuringPattern() {
		lhs = d.NewNodeWithChild(rhino.Token_DESTRUCTURING_LHS, node)
	}
	if decl.Initializer != nil {
		initializer := d.Transform(decl.Initializer)
		lhs.AddChildToBack(initializer)
		d.SetLength(lhs, decl.Location.Start, decl.Location.End)
	}
	return lhs
}

// Java: IRFactory.TransformDispatcher.processWhileLoop(WhileStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessWhileLoop(stmt *parser.WhileStatementTree) *rhino.Node {
	return d.NewNodeWith2Children(rhino.Token_WHILE, d.Transform(stmt.Condition), d.TransformBlock(stmt.Body))
}

// Java: IRFactory.TransformDispatcher.processWithStatement(WithStatementTree)
func (d *IRFactory_TransformDispatcher) ProcessWithStatement(stmt *parser.WithStatementTree) *rhino.Node {
	return d.NewNodeWith2Children(rhino.Token_WITH, d.Transform(stmt.Expression), d.TransformBlock(stmt.Body))
}

// Java: IRFactory.TransformDispatcher.processMissingExpression(MissingPrimaryExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessMissingExpression(unused *parser.MissingPrimaryExpressionTree) *rhino.Node {
	// This will already have been reported as an error by the parser.
	// Try to create something valid that ide mode might be able to
	// continue with.
	return d.CreateMissingExpressionNode()
}

// Java: IRFactory.TransformDispatcher.createMissingNameNode()
func (d *IRFactory_TransformDispatcher) CreateMissingNameNode() *rhino.Node {
	return d.NewStringNodeWithToken(rhino.Token_NAME, "__missing_name__")
}

// Java: IRFactory.TransformDispatcher.createMissingExpressionNode()
func (d *IRFactory_TransformDispatcher) CreateMissingExpressionNode() *rhino.Node {
	return d.NewStringNodeWithToken(rhino.Token_NAME, "__missing_expression__")
}

// Java: IRFactory.TransformDispatcher.processIllegalToken(ParseTree)
func (d *IRFactory_TransformDispatcher) ProcessIllegalToken(node parser.ParseTree) *rhino.Node {
	d.errorReporter.Error("Unsupported syntax: "+node.GetParseTreeBase().Type.String(), d.sourceName, IRFactory_Lineno(node), 0)
	return d.NewNode(rhino.Token_EMPTY)
}

// MaybeReportGetter reports an illegal getter and returns true if the
// language mode is too low.
//
// Java: IRFactory.TransformDispatcher.maybeReportGetter(ParseTree)
func (d *IRFactory_TransformDispatcher) MaybeReportGetter(node parser.ParseTree) bool {
	d.features = d.features.With(parser.FeatureSet_Feature_GETTER)
	if d.config.LanguageMode() == Config_LanguageMode_ECMASCRIPT3 {
		d.errorReporter.Error(IRFactory_GETTER_ERROR_MESSAGE, d.sourceName, IRFactory_Lineno(node), 0)
		return true
	}
	return false
}

// MaybeReportSetter reports an illegal setter and returns true if the
// language mode is too low.
//
// Java: IRFactory.TransformDispatcher.maybeReportSetter(ParseTree)
func (d *IRFactory_TransformDispatcher) MaybeReportSetter(node parser.ParseTree) bool {
	d.features = d.features.With(parser.FeatureSet_Feature_SETTER)
	if d.config.LanguageMode() == Config_LanguageMode_ECMASCRIPT3 {
		d.errorReporter.Error(IRFactory_SETTER_ERROR_MESSAGE, d.sourceName, IRFactory_Lineno(node), 0)
		return true
	}
	return false
}

// Java: IRFactory.TransformDispatcher.processBooleanLiteral(LiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessBooleanLiteral(literal *parser.LiteralExpressionTree) *rhino.Node {
	return d.NewNode(IRFactory_TransformBooleanTokenType(literal.LiteralToken.GetTokenBase().Type))
}

// Java: IRFactory.TransformDispatcher.processNullLiteral(LiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessNullLiteral(unused *parser.LiteralExpressionTree) *rhino.Node {
	return d.NewNode(rhino.Token_NULL)
}

// Java: IRFactory.TransformDispatcher.processNull(NullTree)
func (d *IRFactory_TransformDispatcher) ProcessNull(unused *parser.NullTree) *rhino.Node {
	// NOTE: This is not a NULL literal but a placeholder node such as in
	// an array with "holes".
	return d.NewNode(rhino.Token_EMPTY)
}

// Java: IRFactory.TransformDispatcher.processCommaExpression(CommaExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessCommaExpression(tree *parser.CommaExpressionTree) *rhino.Node {
	root := d.NewNode(rhino.Token_COMMA)
	start := tree.Expressions[0].GetParseTreeBase().Location.Start
	end := tree.Expressions[1].GetParseTreeBase().Location.End
	d.SetSourceInfoFromPositions(root, start, end)
	for _, expr := range tree.Expressions {
		count := root.GetChildCount()
		if count < 2 {
			root.AddChildToBack(d.Transform(expr))
		} else {
			end = expr.GetParseTreeBase().Location.End
			root = d.NewNodeWith2Children(rhino.Token_COMMA, root, d.Transform(expr))
			d.SetSourceInfoFromPositions(root, start, end)
		}
	}
	return root
}

// Java: IRFactory.TransformDispatcher.processClassDeclaration(ClassDeclarationTree)
func (d *IRFactory_TransformDispatcher) ProcessClassDeclaration(tree *parser.ClassDeclarationTree) *rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_CLASSES)

	name := d.TransformOrEmptyIdentifier(tree.Name, tree)
	superClass := d.TransformOrEmpty(tree.SuperClass, tree)

	classMembers := d.NewNode(rhino.Token_CLASS_MEMBERS)
	d.SetSourceInfo(classMembers, tree)

	func() {
		ignored := d.classScope.Increment()
		defer ignored.Close()
		hasConstructor := false
		for _, child := range tree.Elements {
			switch child.GetParseTreeBase().Type {
			case parser.ParseTreeType_COMPUTED_PROPERTY_GETTER, parser.ParseTreeType_COMPUTED_PROPERTY_SETTER,
				parser.ParseTreeType_GET_ACCESSOR, parser.ParseTreeType_SET_ACCESSOR:
				d.features = d.features.With(parser.FeatureSet_Feature_CLASS_GETTER_SETTER)
			case parser.ParseTreeType_BLOCK:
				d.features = d.features.With(parser.FeatureSet_Feature_CLASS_STATIC_BLOCK)
			default:
			}

			childIsCtor := d.ValidateClassConstructorMember(child) // Has side-effects.
			if childIsCtor {
				if hasConstructor {
					d.errorReporter.Error(
						"Class may have only one constructor.", //
						d.sourceName,
						IRFactory_Lineno(child),
						IRFactory_Charno(child))
				}
				hasConstructor = true
			}

			classMembers.AddChildToBack(d.Transform(child))
		}
	}()

	// When nested inside of multiple classes, only run the validation once
	// the top-most class exits scope.
	if !d.classScope.InScope() {
		// Go-only: language-adapter — Collections.unmodifiableSet is a
		// read-only view; the set is passed as it is.
		d.ValidatePrivatePropertyUsage(d.GetPrivatePropsAndValidateUniqueness(classMembers), classMembers)
	}

	classNode := d.NewNodeWith3Children(rhino.Token_CLASS, name, superClass, classMembers)
	d.AttachPossibleTrailingComment(classNode, tree.GetEnd())

	return classNode
}

// ValidatePrivatePropertyUsage validates that all referenced private
// properties within the scope of classMembers are found in privatePropNames
// and private fields are not deleted (both of these cases are syntax
// errors).
//
// Recurses into inner classes and successively adds to privatePropNames such
// that outer class private properties are available to inner classes.
//
// Java: IRFactory.TransformDispatcher.validatePrivatePropertyUsage(Set<String>, Node)
func (d *IRFactory_TransformDispatcher) ValidatePrivatePropertyUsage(privatePropNames *jutil.TreeSet[string], classMembers *rhino.Node) {
	jutil.CheckState(classMembers.IsClassMembers())
	if !classMembers.HasChildren() {
		return
	}

	// Uses a stack object instead of recursion to avoid a stack overflow.
	var visitStack []*rhino.Node
	visitStack = append(visitStack, classMembers.GetFirstChild())
	for len(visitStack) != 0 {
		node := visitStack[len(visitStack)-1]
		visitStack = visitStack[:len(visitStack)-1]

		// Check that these usages reference a private property declared in
		// the class:
		//   this.#privateProp
		//   x.#privateProp
		//   x?.#privateProp
		//   #privateProp in x
		if node.IsPrivateIdentifier() &&
			(node.IsGetProp() ||
				node.IsOptChainGetProp() ||
				(node.IsName() && node.HasParent() && node.GetParent().IsIn())) &&
			!privatePropNames.Contains(node.GetString()) {
			propType := "fields"
			if node.GetParent().IsCall() {
				propType = "methods"
			}
			d.errorReporter.Error(
				"Private "+propType+" must be declared in an enclosing class",
				d.sourceName,
				node.GetLineno(),
				node.GetCharno())
		}

		// Checks that private fields are not deleted:
		//   delete this.#privateProp
		//   delete x.#privateProp
		//   delete x?.#privateProp
		if node.IsDelProp() {
			firstChild := jutil.CheckNotNull(node.GetFirstChild())
			if firstChild.IsPrivateIdentifier() &&
				(firstChild.IsGetProp() || firstChild.IsOptChainGetProp()) {
				d.errorReporter.Error(
					"Private fields cannot be deleted", d.sourceName, node.GetLineno(), node.GetCharno())
			}
		}

		if node.IsClassMembers() {
			// Recurse into the inner class with the outer and inner class
			// private properties.
			innerClassPrivatePropNames := d.GetPrivatePropsAndValidateUniqueness(node)
			innerClassPrivatePropNames.AddAll(privatePropNames.All())
			d.ValidatePrivatePropertyUsage(innerClassPrivatePropNames, node)
		} else {
			if node.GetNext() != nil {
				visitStack = append(visitStack, node.GetNext())
			}
			if node.HasChildren() {
				visitStack = append(visitStack, node.GetFirstChild())
			}
		}
	}
}

// GetPrivatePropsAndValidateUniqueness returns the set of private properties
// declared in the class corresponding to classMembers.
//
// While gathering private properties, validates that there are no duplicate
// property names. Note that, for private getters and setters, one getter and
// one setter can have the same name provided they have the same static or
// non-static modifier.
//
// Java: IRFactory.TransformDispatcher.getPrivatePropsAndValidateUniqueness(Node)
func (d *IRFactory_TransformDispatcher) GetPrivatePropsAndValidateUniqueness(classMembers *rhino.Node) *jutil.TreeSet[string] {
	// Track all this to ensure we don't have duplicate private getters and
	// setters.
	privateGetterNames := jutil.NewTreeSet(jutil.StringCompareTo)
	privateSetterNames := jutil.NewTreeSet(jutil.StringCompareTo)
	privateStaticGetterNames := jutil.NewTreeSet(jutil.StringCompareTo)
	privateStaticSetterNames := jutil.NewTreeSet(jutil.StringCompareTo)
	privateFieldAndMethodNames := jutil.NewTreeSet(jutil.StringCompareTo)

	privatePropNames := jutil.NewTreeSet(jutil.StringCompareTo)
	for curNode := classMembers.GetFirstChild(); curNode != nil; curNode = curNode.GetNext() {
		if !curNode.IsPrivateIdentifier() {
			continue
		}

		propName := curNode.GetString()
		jutil.CheckStateMsg2(
			curNode.IsGetterDef() ||
				curNode.IsSetterDef() ||
				curNode.IsMemberFieldDef() ||
				curNode.IsMemberFunctionDef(),
			"Private property '%s' has an unsupported token type: %s",
			propName,
			curNode.GetToken())

		isStatic := curNode.IsStaticMember()
		alreadyDeclared := privatePropNames.Contains(propName)
		switch curNode.GetToken() {
		case rhino.Token_GETTER_DEF:
			if alreadyDeclared {
				var otherSetter bool
				if isStatic {
					otherSetter = privateSetterNames.Contains(propName)
				} else {
					otherSetter = privateStaticSetterNames.Contains(propName)
				}
				alreadyDeclaredIsOtherThanItsSetter := privateFieldAndMethodNames.Contains(propName) ||
					privateGetterNames.Contains(propName) ||
					privateStaticGetterNames.Contains(propName) ||
					otherSetter
				var itsSetter bool
				if isStatic {
					itsSetter = privateStaticSetterNames.Contains(propName)
				} else {
					itsSetter = privateSetterNames.Contains(propName)
				}
				jutil.CheckState(alreadyDeclaredIsOtherThanItsSetter || itsSetter)
				alreadyDeclared = alreadyDeclaredIsOtherThanItsSetter
			}

			if isStatic {
				privateStaticGetterNames.Add(propName)
			} else {
				privateGetterNames.Add(propName)
			}
		case rhino.Token_SETTER_DEF:
			if alreadyDeclared {
				var otherGetter bool
				if isStatic {
					otherGetter = privateGetterNames.Contains(propName)
				} else {
					otherGetter = privateStaticGetterNames.Contains(propName)
				}
				alreadyDeclaredIsOtherThanItsGetter := privateFieldAndMethodNames.Contains(propName) ||
					privateSetterNames.Contains(propName) ||
					privateStaticSetterNames.Contains(propName) ||
					otherGetter
				var itsGetter bool
				if isStatic {
					itsGetter = privateStaticGetterNames.Contains(propName)
				} else {
					itsGetter = privateGetterNames.Contains(propName)
				}
				jutil.CheckState(alreadyDeclaredIsOtherThanItsGetter || itsGetter)
				alreadyDeclared = alreadyDeclaredIsOtherThanItsGetter
			}

			if curNode.IsStaticMember() {
				privateStaticSetterNames.Add(propName)
			} else {
				privateSetterNames.Add(propName)
			}
		case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_MEMBER_FUNCTION_DEF:
			privateFieldAndMethodNames.Add(propName)
		default:
		}
		if alreadyDeclared {
			d.errorReporter.Error(
				"Identifier '"+propName+"' has already been declared",
				d.sourceName,
				curNode.GetLineno(),
				curNode.GetCharno())
		}

		privatePropNames.Add(propName)
	}
	return privatePropNames
}

// ValidateClassConstructorMember returns true iff this member is a legal
// class constructor.
//
// Java: IRFactory.TransformDispatcher.validateClassConstructorMember(ParseTree)
func (d *IRFactory_TransformDispatcher) ValidateClassConstructorMember(member parser.ParseTree) bool {
	var memberName parser.Token
	var isStatic bool
	var hasIllegalModifier bool
	switch member.GetParseTreeBase().Type {
	case parser.ParseTreeType_GET_ACCESSOR:
		getter := member.(*parser.GetAccessorTree)
		memberName = getter.PropertyName
		isStatic = getter.IsStatic
		hasIllegalModifier = true
	case parser.ParseTreeType_SET_ACCESSOR:
		setter := member.(*parser.SetAccessorTree)
		memberName = setter.PropertyName
		isStatic = setter.IsStatic
		hasIllegalModifier = true
	case parser.ParseTreeType_FUNCTION_DECLARATION:
		method := member.(*parser.FunctionDeclarationTree)
		if method.Name != nil {
			// Go-only: language-adapter — no typed nil in the interface.
			memberName = method.Name
		}
		isStatic = method.IsStatic
		hasIllegalModifier = method.IsGenerator || method.IsAsync
	default:
		// Computed properties aren't an issue here because they aren't used
		// as the class constructor, regardless of their name.
		return false
	}

	if isStatic {
		// Statics are fine because they're never the class constructor.
		return false
	}

	if memberName.GetTokenBase().Type != parser.TokenType_IDENTIFIER ||
		!memberName.(*parser.IdentifierToken).ValueEquals("constructor") {
		// There's only a potential issue if the member is named
		// "constructor".
		// TODO(b/123769080): Also check for quoted string literals with the
		// value "constructor".
		return false
	}

	if hasIllegalModifier {
		d.errorReporter.Error(
			"Class constructor may not be getter, setter, async, or generator.",
			d.sourceName,
			IRFactory_Lineno(member),
			IRFactory_Charno(member))
		return false
	}

	return true
}

// Java: IRFactory.TransformDispatcher.processSuper(SuperExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessSuper(tree *parser.SuperExpressionTree) *rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_SUPER)
	return d.NewNode(rhino.Token_SUPER)
}

// Java: IRFactory.TransformDispatcher.processNewTarget(NewTargetExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessNewTarget(tree *parser.NewTargetExpressionTree) *rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_NEW_TARGET)
	return d.NewNode(rhino.Token_NEW_TARGET)
}

// Java: IRFactory.TransformDispatcher.processYield(YieldExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessYield(tree *parser.YieldExpressionTree) *rhino.Node {
	yield := d.NewNode(rhino.Token_YIELD)
	if tree.Expression != nil {
		yield.AddChildToBack(d.Transform(tree.Expression))
	}
	yield.SetYieldAll(tree.IsYieldAll)
	return yield
}

// Java: IRFactory.TransformDispatcher.processAwait(AwaitExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessAwait(tree *parser.AwaitExpressionTree) *rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_ASYNC_FUNCTIONS)
	await := d.NewNode(rhino.Token_AWAIT)
	await.AddChildToBack(d.Transform(tree.Expression))
	return await
}

// Java: IRFactory.TransformDispatcher.processExportDecl(ExportDeclarationTree)
func (d *IRFactory_TransformDispatcher) ProcessExportDecl(tree *parser.ExportDeclarationTree) *rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_MODULES)
	var decls *rhino.Node
	if tree.IsExportAll {
		jutil.CheckState(tree.Declaration == nil && tree.ExportSpecifierList == nil)
	} else if tree.Declaration != nil {
		jutil.CheckState(tree.ExportSpecifierList == nil)
		decls = d.Transform(tree.Declaration)
	} else {
		decls = d.TransformList(rhino.Token_EXPORT_SPECS, tree.ExportSpecifierList)
	}
	if decls == nil {
		decls = d.NewNode(rhino.Token_EMPTY)
	}
	d.SetSourceInfo(decls, tree)
	export := d.NewNodeWithChild(rhino.Token_EXPORT, decls)
	if tree.From != nil {
		from := d.ProcessString(tree.From)
		export.AddChildToBack(from)
	}

	export.PutBooleanProp(rhino.Node_EXPORT_ALL_FROM, tree.IsExportAll)
	export.PutBooleanProp(rhino.Node_EXPORT_DEFAULT, tree.IsDefault)
	return export
}

// Java: IRFactory.TransformDispatcher.processExportSpec(ExportSpecifierTree)
func (d *IRFactory_TransformDispatcher) ProcessExportSpec(tree *parser.ExportSpecifierTree) *rhino.Node {
	importedName := d.ProcessNameOfToken(tree.ImportedName, rhino.Token_NAME, IRFactory_TransformDispatcher_IdentifierType_STANDARD)
	exportSpec := d.NewNodeWithChild(rhino.Token_EXPORT_SPEC, importedName)
	if tree.DestinationName == nil {
		exportSpec.SetShorthandProperty(true)
		exportSpec.AddChildToBack(importedName.CloneTree())
	} else {
		destinationName := d.ProcessNameOfToken(tree.DestinationName, rhino.Token_NAME, IRFactory_TransformDispatcher_IdentifierType_STANDARD)
		exportSpec.AddChildToBack(destinationName)
	}
	return exportSpec
}

// Java: IRFactory.TransformDispatcher.processImportDecl(ImportDeclarationTree)
func (d *IRFactory_TransformDispatcher) ProcessImportDecl(tree *parser.ImportDeclarationTree) *rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_MODULES)

	firstChild := d.TransformOrEmptyIdentifier(tree.DefaultBindingIdentifier, tree)
	var secondChild *rhino.Node
	if tree.NameSpaceImportIdentifier == nil {
		secondChild = d.TransformListOrEmpty(rhino.Token_IMPORT_SPECS, tree.ImportSpecifierList)
		// Currently source info is "import {foo} from '...';" expression. If
		// needed this should be changed to use only "{foo}" part.
		d.SetSourceInfo(secondChild, tree)
	} else {
		secondChild = d.NewStringNodeFromIdentifier(
			rhino.Token_IMPORT_STAR, IRFactory_TransformDispatcher_IdentifierType_STANDARD, tree.NameSpaceImportIdentifier)
		d.SetSourceInfoFromToken(secondChild, tree.NameSpaceImportIdentifier)
	}
	thirdChild := d.ProcessString(tree.ModuleSpecifier)

	return d.NewNodeWith3Children(rhino.Token_IMPORT, firstChild, secondChild, thirdChild)
}

// Java: IRFactory.TransformDispatcher.processImportSpec(ImportSpecifierTree)
func (d *IRFactory_TransformDispatcher) ProcessImportSpec(tree *parser.ImportSpecifierTree) *rhino.Node {
	importedName := d.ProcessNameOfToken(tree.ImportedName, rhino.Token_NAME, IRFactory_TransformDispatcher_IdentifierType_STANDARD)
	importSpec := d.NewNodeWithChild(rhino.Token_IMPORT_SPEC, importedName)
	if tree.DestinationName == nil {
		importSpec.SetShorthandProperty(true)
		importSpec.AddChildToBack(importedName.CloneTree())
	} else {
		importSpec.AddChildToBack(
			d.ProcessNameOfToken(tree.DestinationName, rhino.Token_NAME, IRFactory_TransformDispatcher_IdentifierType_STANDARD))
	}
	return importSpec
}

// Java: IRFactory.TransformDispatcher.processDynamicImport(DynamicImportTree)
func (d *IRFactory_TransformDispatcher) ProcessDynamicImport(dynamicImportNode *parser.DynamicImportTree) *rhino.Node {
	d.MaybeWarnForFeature(dynamicImportNode, parser.FeatureSet_Feature_DYNAMIC_IMPORT)
	argument := d.Transform(dynamicImportNode.Argument)
	return d.NewNodeWithChild(rhino.Token_DYNAMIC_IMPORT, argument)
}

// Java: IRFactory.TransformDispatcher.processImportMeta(ImportMetaExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessImportMeta(tree *parser.ImportMetaExpressionTree) *rhino.Node {
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_MODULES)
	d.MaybeWarnForFeature(tree, parser.FeatureSet_Feature_IMPORT_META)
	return d.NewNode(rhino.Token_IMPORT_META)
}

// Java: IRFactory.TransformDispatcher.checkParameters(FormalParameterListTree)
func (d *IRFactory_TransformDispatcher) CheckParameters(tree *parser.FormalParameterListTree) bool {
	params := tree.Parameters
	good := true
	for i := 0; i < len(params); i++ {
		param := params[i]
		if param.GetParseTreeBase().Type == parser.ParseTreeType_ITER_REST {
			if i != len(params)-1 {
				d.errorReporter.Error(
					"A rest parameter must be last in a parameter list.",
					d.sourceName,
					IRFactory_Lineno(param),
					IRFactory_Charno(param))
				good = false
			} else if tree.HasTrailingComma {
				d.errorReporter.Error(
					"A trailing comma must not follow a rest parameter.",
					d.sourceName,
					IRFactory_Lineno(param),
					IRFactory_Charno(param))
				good = false
			}
		}
	}
	return good
}

// Java: IRFactory.TransformDispatcher.transformList(Token, ImmutableList<ParseTree>)
func (d *IRFactory_TransformDispatcher) TransformList(type_ rhino.Token, list []parser.ParseTree) *rhino.Node {
	n := d.NewNode(type_)
	// Go-only: language-adapter — a nil slice is Java null (the parser keeps
	// empty lists non-nil), and an enhanced for statement over a null list
	// throws NullPointerException where a Go range does nothing. Of the lists
	// this file iterates, only ExportDeclarationTree.exportSpecifierList
	// reaches here as null (an export the parser recovered from in
	// KEEP_GOING mode).
	if list == nil {
		panic(jutil.NewNullPointerException())
	}
	for _, tree := range list {
		n.AddChildToBack(d.Transform(tree))
	}
	return n
}

// Java: IRFactory.TransformDispatcher.transformListOrEmpty(Token, ImmutableList<ParseTree>)
func (d *IRFactory_TransformDispatcher) TransformListOrEmpty(type_ rhino.Token, list []parser.ParseTree) *rhino.Node {
	if list == nil || len(list) == 0 {
		return d.NewNode(rhino.Token_EMPTY)
	} else {
		return d.TransformList(type_, list)
	}
}

// Java: IRFactory.TransformDispatcher.unsupportedLanguageFeature(ParseTree, String)
func (d *IRFactory_TransformDispatcher) UnsupportedLanguageFeature(node parser.ParseTree, feature string) *rhino.Node {
	d.errorReporter.Error(
		"unsupported language feature: "+feature, d.sourceName, IRFactory_Lineno(node), IRFactory_Charno(node))
	return d.CreateMissingExpressionNode()
}

// Java: IRFactory.TransformDispatcher.processLiteralExpression(LiteralExpressionTree)
func (d *IRFactory_TransformDispatcher) ProcessLiteralExpression(expr *parser.LiteralExpressionTree) *rhino.Node {
	switch expr.LiteralToken.GetTokenBase().Type {
	case parser.TokenType_NUMBER:
		return d.ProcessNumberLiteral(expr)
	case parser.TokenType_STRING:
		return d.ProcessStringLiteral(expr)
	case parser.TokenType_BIGINT:
		return d.ProcessBigIntLiteral(expr)
	case parser.TokenType_FALSE, parser.TokenType_TRUE:
		return d.ProcessBooleanLiteral(expr)
	case parser.TokenType_NULL:
		return d.ProcessNullLiteral(expr)
	case parser.TokenType_REGULAR_EXPRESSION:
		return d.ProcessRegExpLiteral(expr)
	default:
		panic(jutil.NewIllegalStateExceptionMsg(
			"Unexpected literal type: " +
				irFactory_javaClassOf(expr.LiteralToken) +
				" type: " +
				expr.LiteralToken.GetTokenBase().Type.String()))
	}
}

// irFactory_javaClassOf is Object.getClass().toString() for a parser token:
// "class com.google.javascript.jscomp.parsing.parser.<simple name>".
//
// Go-only: language-adapter — the Go token types are named after the Java
// classes (a *Base struct is the non-abstract Java base class itself).
func irFactory_javaClassOf(token parser.Token) string {
	t := reflect.TypeOf(token)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return "class com.google.javascript.jscomp.parsing.parser." + strings.TrimSuffix(t.Name(), "Base")
}

// Java: IRFactory.TransformDispatcher.process(ParseTree)
func (d *IRFactory_TransformDispatcher) Process(node parser.ParseTree) *rhino.Node {
	switch node.GetParseTreeBase().Type {
	case parser.ParseTreeType_BINARY_OPERATOR:
		return d.ProcessBinaryExpression(node.(*parser.BinaryOperatorTree))
	case parser.ParseTreeType_ARRAY_LITERAL_EXPRESSION:
		return d.ProcessArrayLiteral(node.(*parser.ArrayLiteralExpressionTree))
	case parser.ParseTreeType_TEMPLATE_LITERAL_EXPRESSION:
		return d.ProcessTemplateLiteral(node.(*parser.TemplateLiteralExpressionTree))
	case parser.ParseTreeType_TEMPLATE_LITERAL_PORTION:
		return d.ProcessTemplateLiteralPortion(node.(*parser.TemplateLiteralPortionTree))
	case parser.ParseTreeType_TEMPLATE_SUBSTITUTION:
		return d.ProcessTemplateSubstitution(node.(*parser.TemplateSubstitutionTree))
	case parser.ParseTreeType_UNARY_EXPRESSION:
		return d.ProcessUnaryExpression(node.(*parser.UnaryExpressionTree))
	case parser.ParseTreeType_BLOCK:
		return d.ProcessBlock(node.(*parser.BlockTree))
	case parser.ParseTreeType_BREAK_STATEMENT:
		return d.ProcessBreakStatement(node.(*parser.BreakStatementTree))
	case parser.ParseTreeType_CALL_EXPRESSION:
		return d.ProcessFunctionCall(node.(*parser.CallExpressionTree))
	case parser.ParseTreeType_OPT_CHAIN_CALL_EXPRESSION:
		return d.ProcessOptChainFunctionCall(node.(*parser.OptChainCallExpressionTree))
	case parser.ParseTreeType_CASE_CLAUSE:
		return d.ProcessSwitchCase(node.(*parser.CaseClauseTree))
	case parser.ParseTreeType_DEFAULT_CLAUSE:
		return d.ProcessSwitchDefault(node.(*parser.DefaultClauseTree))
	case parser.ParseTreeType_CATCH:
		return d.ProcessCatchClause(node.(*parser.CatchTree))
	case parser.ParseTreeType_CONTINUE_STATEMENT:
		return d.ProcessContinueStatement(node.(*parser.ContinueStatementTree))
	case parser.ParseTreeType_DO_WHILE_STATEMENT:
		return d.ProcessDoLoop(node.(*parser.DoWhileStatementTree))
	case parser.ParseTreeType_EMPTY_STATEMENT:
		return d.ProcessEmptyStatement(node.(*parser.EmptyStatementTree))
	case parser.ParseTreeType_EXPRESSION_STATEMENT:
		return d.ProcessExpressionStatement(node.(*parser.ExpressionStatementTree))
	case parser.ParseTreeType_DEBUGGER_STATEMENT:
		return d.ProcessDebuggerStatement(node.(*parser.DebuggerStatementTree))
	case parser.ParseTreeType_THIS_EXPRESSION:
		return d.ProcessThisExpression(node.(*parser.ThisExpressionTree))
	case parser.ParseTreeType_FOR_STATEMENT:
		return d.ProcessForLoop(node.(*parser.ForStatementTree))
	case parser.ParseTreeType_FOR_IN_STATEMENT:
		return d.ProcessForInLoop(node.(*parser.ForInStatementTree))
	case parser.ParseTreeType_FUNCTION_DECLARATION:
		return d.ProcessFunction(node.(*parser.FunctionDeclarationTree))
	case parser.ParseTreeType_MEMBER_LOOKUP_EXPRESSION:
		return d.ProcessElementGet(node.(*parser.MemberLookupExpressionTree))
	case parser.ParseTreeType_OPT_CHAIN_MEMBER_LOOKUP_EXPRESSION:
		return d.ProcessOptChainElementGet(node.(*parser.OptionalMemberLookupExpressionTree))
	case parser.ParseTreeType_MEMBER_EXPRESSION:
		return d.ProcessPropertyGet(node.(*parser.MemberExpressionTree))
	case parser.ParseTreeType_OPT_CHAIN_MEMBER_EXPRESSION:
		return d.ProcessOptChainPropertyGet(node.(*parser.OptionalMemberExpressionTree))
	case parser.ParseTreeType_CONDITIONAL_EXPRESSION:
		return d.ProcessConditionalExpression(node.(*parser.ConditionalExpressionTree))
	case parser.ParseTreeType_IF_STATEMENT:
		return d.ProcessIfStatement(node.(*parser.IfStatementTree))
	case parser.ParseTreeType_LABELLED_STATEMENT:
		return d.ProcessLabeledStatement(node.(*parser.LabelledStatementTree))
	case parser.ParseTreeType_PAREN_EXPRESSION:
		return d.ProcessParenthesizedExpression(node.(*parser.ParenExpressionTree))
	case parser.ParseTreeType_IDENTIFIER_EXPRESSION:
		identifierType := IRFactory_TransformDispatcher_IdentifierType_STANDARD
		if d.privateIdLhsOfInScope.InScope() {
			identifierType = IRFactory_TransformDispatcher_IdentifierType_CAN_BE_PRIVATE
		}
		return d.ProcessName(node.(*parser.IdentifierExpressionTree), identifierType)
	case parser.ParseTreeType_NEW_EXPRESSION:
		return d.ProcessNewExpression(node.(*parser.NewExpressionTree))
	case parser.ParseTreeType_OBJECT_LITERAL_EXPRESSION:
		return d.ProcessObjectLiteral(node.(*parser.ObjectLiteralExpressionTree))
	case parser.ParseTreeType_COMPUTED_PROPERTY_DEFINITION:
		return d.ProcessComputedPropertyDefinition(node.(*parser.ComputedPropertyDefinitionTree))
	case parser.ParseTreeType_COMPUTED_PROPERTY_GETTER:
		return d.ProcessComputedPropertyGetter(node.(*parser.ComputedPropertyGetterTree))
	case parser.ParseTreeType_COMPUTED_PROPERTY_METHOD:
		return d.ProcessComputedPropertyMethod(node.(*parser.ComputedPropertyMethodTree))
	case parser.ParseTreeType_COMPUTED_PROPERTY_SETTER:
		return d.ProcessComputedPropertySetter(node.(*parser.ComputedPropertySetterTree))
	case parser.ParseTreeType_RETURN_STATEMENT:
		return d.ProcessReturnStatement(node.(*parser.ReturnStatementTree))
	case parser.ParseTreeType_UPDATE_EXPRESSION:
		return d.ProcessUpdateExpression(node.(*parser.UpdateExpressionTree))
	case parser.ParseTreeType_PROGRAM:
		return d.ProcessAstRoot(node.(*parser.ProgramTree))
	case parser.ParseTreeType_LITERAL_EXPRESSION:
		// STRING, NUMBER, TRUE, FALSE, NULL, REGEXP
		return d.ProcessLiteralExpression(node.(*parser.LiteralExpressionTree))
	case parser.ParseTreeType_SWITCH_STATEMENT:
		return d.ProcessSwitchStatement(node.(*parser.SwitchStatementTree))
	case parser.ParseTreeType_THROW_STATEMENT:
		return d.ProcessThrowStatement(node.(*parser.ThrowStatementTree))
	case parser.ParseTreeType_TRY_STATEMENT:
		return d.ProcessTryStatement(node.(*parser.TryStatementTree))
	case parser.ParseTreeType_VARIABLE_STATEMENT:
		// var const let
		return d.ProcessVariableStatement(node.(*parser.VariableStatementTree))
	case parser.ParseTreeType_VARIABLE_DECLARATION_LIST:
		return d.ProcessVariableDeclarationList(node.(*parser.VariableDeclarationListTree))
	case parser.ParseTreeType_VARIABLE_DECLARATION:
		return d.ProcessVariableDeclaration(node.(*parser.VariableDeclarationTree))
	case parser.ParseTreeType_WHILE_STATEMENT:
		return d.ProcessWhileLoop(node.(*parser.WhileStatementTree))
	case parser.ParseTreeType_WITH_STATEMENT:
		return d.ProcessWithStatement(node.(*parser.WithStatementTree))
	case parser.ParseTreeType_COMMA_EXPRESSION:
		return d.ProcessCommaExpression(node.(*parser.CommaExpressionTree))
	case parser.ParseTreeType_NULL:
		// this is not the null literal
		return d.ProcessNull(node.(*parser.NullTree))
	case parser.ParseTreeType_FINALLY:
		return d.ProcessFinally(node.(*parser.FinallyTree))
	case parser.ParseTreeType_MISSING_PRIMARY_EXPRESSION:
		return d.ProcessMissingExpression(node.(*parser.MissingPrimaryExpressionTree))
	case parser.ParseTreeType_PROPERTY_NAME_ASSIGNMENT:
		return d.ProcessPropertyNameAssignment(node.(*parser.PropertyNameAssignmentTree))
	case parser.ParseTreeType_GET_ACCESSOR:
		return d.ProcessGetAccessor(node.(*parser.GetAccessorTree))
	case parser.ParseTreeType_SET_ACCESSOR:
		return d.ProcessSetAccessor(node.(*parser.SetAccessorTree))
	case parser.ParseTreeType_FORMAL_PARAMETER_LIST:
		return d.ProcessFormalParameterList(node.(*parser.FormalParameterListTree))
	case parser.ParseTreeType_CLASS_DECLARATION:
		return d.ProcessClassDeclaration(node.(*parser.ClassDeclarationTree))
	case parser.ParseTreeType_SUPER_EXPRESSION:
		return d.ProcessSuper(node.(*parser.SuperExpressionTree))
	case parser.ParseTreeType_NEW_TARGET_EXPRESSION:
		return d.ProcessNewTarget(node.(*parser.NewTargetExpressionTree))
	case parser.ParseTreeType_YIELD_EXPRESSION:
		return d.ProcessYield(node.(*parser.YieldExpressionTree))
	case parser.ParseTreeType_AWAIT_EXPRESSION:
		return d.ProcessAwait(node.(*parser.AwaitExpressionTree))
	case parser.ParseTreeType_FOR_OF_STATEMENT:
		return d.ProcessForOf(node.(*parser.ForOfStatementTree))
	case parser.ParseTreeType_FOR_AWAIT_OF_STATEMENT:
		return d.ProcessForAwaitOf(node.(*parser.ForAwaitOfStatementTree))
	case parser.ParseTreeType_EXPORT_DECLARATION:
		return d.ProcessExportDecl(node.(*parser.ExportDeclarationTree))
	case parser.ParseTreeType_EXPORT_SPECIFIER:
		return d.ProcessExportSpec(node.(*parser.ExportSpecifierTree))
	case parser.ParseTreeType_IMPORT_DECLARATION:
		return d.ProcessImportDecl(node.(*parser.ImportDeclarationTree))
	case parser.ParseTreeType_IMPORT_SPECIFIER:
		return d.ProcessImportSpec(node.(*parser.ImportSpecifierTree))
	case parser.ParseTreeType_DYNAMIC_IMPORT_EXPRESSION:
		return d.ProcessDynamicImport(node.(*parser.DynamicImportTree))
	case parser.ParseTreeType_IMPORT_META_EXPRESSION:
		return d.ProcessImportMeta(node.(*parser.ImportMetaExpressionTree))
	case parser.ParseTreeType_ARRAY_PATTERN:
		return d.ProcessArrayPattern(node.(*parser.ArrayPatternTree))
	case parser.ParseTreeType_OBJECT_PATTERN:
		return d.ProcessObjectPattern(node.(*parser.ObjectPatternTree))
	case parser.ParseTreeType_COMPREHENSION:
		return d.ProcessComprehension(node.(*parser.ComprehensionTree))
	case parser.ParseTreeType_COMPREHENSION_FOR:
		return d.ProcessComprehensionFor(node.(*parser.ComprehensionForTree))
	case parser.ParseTreeType_COMPREHENSION_IF:
		return d.ProcessComprehensionIf(node.(*parser.ComprehensionIfTree))
	case parser.ParseTreeType_DEFAULT_PARAMETER:
		return d.ProcessDefaultParameter(node.(*parser.DefaultParameterTree))
	case parser.ParseTreeType_ITER_REST:
		return d.ProcessIterRest(node.(*parser.IterRestTree))
	case parser.ParseTreeType_ITER_SPREAD:
		return d.ProcessIterSpread(node.(*parser.IterSpreadTree))
		// ES2019
	case parser.ParseTreeType_OBJECT_REST:
		return d.ProcessObjectPatternElement(node.(*parser.ObjectRestTree))
	case parser.ParseTreeType_OBJECT_SPREAD:
		return d.ProcessObjectSpread(node.(*parser.ObjectSpreadTree))
		// ES2022
	case parser.ParseTreeType_FIELD_DECLARATION:
		return d.ProcessField(node.(*parser.FieldDeclarationTree))
	case parser.ParseTreeType_COMPUTED_PROPERTY_FIELD:
		return d.ProcessComputedPropertyField(node.(*parser.ComputedPropertyFieldTree))
	case parser.ParseTreeType_ARGUMENT_LIST:
		// TODO(johnlenz): handle these or remove parser support
	default:
	}
	return d.ProcessIllegalToken(node)
}

// Java: IRFactory.reportErrorIfYieldOrAwaitInDefaultValue(Node)
func (f *IRFactory) ReportErrorIfYieldOrAwaitInDefaultValue(defaultValueNode *rhino.Node) {
	yieldNode := IRFactory_FindNodeTypeInExpression(defaultValueNode, rhino.Token_YIELD)
	if yieldNode != nil {
		f.errorReporter.Error(
			"`yield` is illegal in parameter default value.",
			irFactory_nullableString(yieldNode.GetSourceFileName()),
			yieldNode.GetLineno(),
			yieldNode.GetCharno())
	}
	awaitNode := IRFactory_FindNodeTypeInExpression(defaultValueNode, rhino.Token_AWAIT)
	if awaitNode != nil {
		f.errorReporter.Error(
			"`await` is illegal in parameter default value.",
			irFactory_nullableString(awaitNode.GetSourceFileName()),
			awaitNode.GetLineno(),
			awaitNode.GetCharno())
	}
}

// irFactory_nullableString passes a nullable Java String to a rhino
// ErrorReporter.
//
// Go-only: language-adapter — rhino.ErrorReporter takes a Go string, so
// Java's null is "".
func irFactory_nullableString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// IRFactory_FindNodeTypeInExpression tries to find a node with the given
// token in the given expression. Returns the first one found in pre-order
// traversal or `null` if none found. Will not traverse into function or
// class expressions.
//
// Java: IRFactory.findNodeTypeInExpression(Node, Token)
func IRFactory_FindNodeTypeInExpression(expressionNode *rhino.Node, token rhino.Token) *rhino.Node {
	// Go-only: language-adapter — the ArrayDeque used as a FIFO queue
	// (add at the tail, remove at the head) is a slice with a head index.
	worklist := []*rhino.Node{expressionNode}
	for head := 0; head < len(worklist); {
		node := worklist[head]
		head++
		if node.GetToken() == token {
			return node
		} else if !node.IsFunction() && !node.IsClass() {
			for child := node.GetFirstChild(); child != nil; child = child.GetNext() {
				worklist = append(worklist, child)
			}
		}
	}
	return nil
}

// irFactory_substring is String.substring(begin, end) on a string given as
// its UTF-16 code units.
//
// Go-only: language-adapter — Java String API on UTF-16 code units.
func irFactory_substring(units []uint16, begin, end int) string {
	if begin < 0 || end > len(units) || begin > end {
		panic(jutil.NewStringIndexOutOfBoundsException(begin, end))
	}
	return jutil.StringFromUTF16(units[begin:end])
}

// irFactory_charAt is String.charAt(int) on a string given as its UTF-16
// code units.
//
// Go-only: language-adapter — Java String API on UTF-16 code units.
func irFactory_charAt(units []uint16, index int) uint16 {
	if index < 0 || index >= len(units) {
		panic(jutil.NewStringIndexOutOfBoundsException(index, len(units)))
	}
	return units[index]
}

// irFactory_indexOf is String.indexOf(int ch, int fromIndex) on a string
// given as its UTF-16 code units.
//
// Go-only: language-adapter — Java String API on UTF-16 code units.
func irFactory_indexOf(units []uint16, ch uint16, fromIndex int) int {
	if fromIndex < 0 {
		fromIndex = 0
	}
	for i := fromIndex; i < len(units); i++ {
		if units[i] == ch {
			return i
		}
	}
	return -1
}

// Java: IRFactory.normalizeRegex(LiteralToken)
func (f *IRFactory) NormalizeRegex(token parser.LiteralToken) string {
	value := token.GetLiteralTokenBase().Value
	if strings.IndexByte(value, '\\') == -1 {
		// Go-only: language-adapter — the delimiting '/' are ASCII, so their
		// byte indices delimit the same substring as Java's UTF-16 indices.
		lastSlash := strings.LastIndexByte(value, '/')
		if lastSlash < 1 {
			panic(jutil.NewStringIndexOutOfBoundsException(1, lastSlash))
		}
		return value[1:lastSlash]
	}
	// Go-only: language-adapter — the String and StringBuilder operations of
	// the escape path work on the UTF-16 code units of the value.
	units := jutil.StringToUTF16(value)
	lastSlash := -1
	for i := len(units) - 1; i >= 0; i-- {
		if units[i] == '/' {
			lastSlash = i
			break
		}
	}
	cur := irFactory_indexOf(units, '\\', 0)

	result := make([]uint16, 0, len(units))
	start := 1
	for cur != -1 {
		result = append(result, units[start:cur]...)

		cur++ // skip the escape char.
		c := irFactory_charAt(units, cur)
		switch c {
		// Characters for which the backslash is semantically important.
		case '^', '$', '\\', '/', '.', '*', '+', '?', '(', ')', '[', ']', '{', '}', '|', '-',
			'b', 'B', 'c', 'd', 'D', 'f', 'n',
			'k', // 2018 named backreferences
			'p', // 2018 unicode property escapes
			'P', // 2018 unicode property escapes
			'r', 's', 'S', 't', 'u', 'v', 'w', 'W', 'x',
			'0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			result = append(result, '\\')
			fallthrough
		default:
			// For all other characters, the backslash has no effect, so just
			// append the next char.
			result = append(result, c)
		}

		start = cur + 1
		cur = irFactory_indexOf(units, '\\', start)
	}
	if start > lastSlash {
		panic(jutil.NewStringIndexOutOfBoundsException(start, lastSlash))
	}
	result = append(result, units[start:lastSlash]...)
	return jutil.StringFromUTF16(result)
}

// Java: IRFactory.normalizeString(LiteralToken, boolean)
func (f *IRFactory) NormalizeString(token parser.LiteralToken, templateLiteral bool) string {
	value := token.GetLiteralTokenBase().Value
	// <CR><LF> and <CR> are normalized as <LF>. For raw template literal
	// string values: this is the spec behaviour. For regular string literals:
	// they can only be part of a line continuation, which we want to scrub.
	if strings.IndexByte(value, '\r') >= 0 {
		// Go-only: language-adapter — value.replaceAll("\r\n?", "\n"); the
		// pattern has no regex features beyond the optional '\n'.
		value = strings.ReplaceAll(strings.ReplaceAll(value, "\r\n", "\n"), "\r", "\n")
	}

	start := 1 // skip the leading quote
	if templateLiteral {
		start = 0
	}
	if strings.IndexByte(value, '\\') == -1 {
		// short circuit no escapes.
		if templateLiteral {
			return value
		}
		return irFactory_substringWithoutQuotes(value)
	}
	// Go-only: language-adapter — the String and StringBuilder operations of
	// the escape path work on the UTF-16 code units of the value.
	units := jutil.StringToUTF16(value)
	cur := irFactory_indexOf(units, '\\', 0)
	result := make([]uint16, 0, len(units))
	tokenStart := token.GetTokenBase().Location.Start
	for cur != -1 {
		result = append(result, units[start:cur]...)

		cur += 1 // skip the escape char.
		c := irFactory_charAt(units, cur)
		switch c {
		case 'b':
			result = append(result, '\b')
		case 'f':
			result = append(result, '\f')
		case 'n':
			result = append(result, '\n')
		case 'r':
			result = append(result, '\r')
		case 't':
			result = append(result, '\t')
		case 'v':
			result = append(result, '\u000B')
		case '\n':
			// line continuation, skip the line break
			f.MaybeWarnForFeatureToken(token, parser.FeatureSet_Feature_STRING_CONTINUATION)
			f.errorReporter.Warning(
				IRFactory_STRING_CONTINUATION_WARNING,
				f.sourceName,
				IRFactory_LinenoOfPosition(tokenStart),
				IRFactory_CharnoOfPosition(tokenStart))
		case '0', '1', '2', '3', '4', '5', '6', '7':
			var numDigits int

			if cur+1 < len(units) && IRFactory_IsOctalDigit(units[cur+1]) {
				if cur+2 < len(units) && IRFactory_IsOctalDigit(units[cur+2]) {
					numDigits = 3
				} else {
					numDigits = 2
				}
			} else {
				numDigits = 1
			}

			if f.InStrictContext() || templateLiteral {
				if c == '0' && numDigits == 1 {
					// No warning: "\0" followed by a character which is not an
					// octal digit is allowed in strict mode.
				} else {
					f.errorReporter.Warning(
						IRFactory_OCTAL_STRING_LITERAL_WARNING,
						f.sourceName,
						IRFactory_LinenoOfPosition(tokenStart),
						IRFactory_CharnoOfPosition(tokenStart))
				}
			}

			result = append(result, uint16(jutil.IntegerParseInt(irFactory_substring(units, cur, cur+numDigits), 8)))
			cur += numDigits - 1
		case 'x':
			result = append(result,
				uint16(IRFactory_Hexdigit(irFactory_charAt(units, cur+1))*0x10+IRFactory_Hexdigit(irFactory_charAt(units, cur+2))))
			cur += 2
		case 'u':
			var escapeEnd int
			var hexDigits string
			if irFactory_charAt(units, cur+1) != '{' {
				// Simple escape with exactly four hex digits: \\uXXXX
				escapeEnd = cur + 5
				hexDigits = irFactory_substring(units, cur+1, escapeEnd)
			} else {
				// Escape with braces can have any number of hex digits:
				// \\u{XXXXXXX}
				escapeEnd = cur + 2
				for jutil.CharacterDigit(irFactory_charAt(units, escapeEnd), 0x10) >= 0 {
					escapeEnd++
				}
				hexDigits = irFactory_substring(units, cur+2, escapeEnd)
				escapeEnd++
			}
			codePointValue := jutil.IntegerParseInt(hexDigits, 0x10)
			if codePointValue > 0x10ffff {
				f.errorReporter.Error(
					"Undefined Unicode code-point",
					f.sourceName,
					IRFactory_LinenoOfPosition(tokenStart),
					IRFactory_CharnoOfPosition(tokenStart))

				// Compilation should stop, but we should finish the string and
				// find more errors. These appends are just to have a
				// placeholder for the errored normalization.
				result = append(result, '\\', 'u', '{')
				result = append(result, jutil.StringToUTF16(hexDigits)...)
				result = append(result, '}')
			} else {
				result = irFactory_appendCodePoint(result, codePointValue)
			}
			cur = escapeEnd - 1
		case '\'', '"', '\\':
			result = append(result, c)
		default:
			result = append(result, c)
		}
		start = cur + 1
		cur = irFactory_indexOf(units, '\\', start)
	}
	// skip the trailing quote.
	end := len(units) - 1
	if templateLiteral {
		end = len(units)
	}
	if start > end {
		panic(jutil.NewStringIndexOutOfBoundsException(start, end))
	}
	result = append(result, units[start:end]...)

	return jutil.StringFromUTF16(result)
}

// irFactory_substringWithoutQuotes is value.substring(1, value.length() - 1).
//
// Go-only: language-adapter — the first char of a string token is an ASCII
// quote; when the last byte is ASCII too, the UTF-16 range [1, length - 1)
// is the byte range [1, len - 1).
func irFactory_substringWithoutQuotes(value string) string {
	if n := len(value); n >= 2 && value[0] < 0x80 && value[n-1] < 0x80 {
		return value[1 : n-1]
	}
	return jutil.StringSubstring(value, 1, jutil.StringLength(value)-1)
}

// irFactory_appendCodePoint is StringBuilder.append(Character.toChars(int)).
//
// Go-only: language-adapter — JDK methods on UTF-16 code units.
func irFactory_appendCodePoint(units []uint16, codePoint int32) []uint16 {
	if codePoint < 0 || codePoint > 0x10ffff {
		panic(jutil.NewIllegalArgumentExceptionMsg(
			"Not a valid Unicode code point: 0x" + strings.ToUpper(jutil.IntegerToHexString(codePoint))))
	}
	if codePoint < 0x10000 {
		return append(units, uint16(codePoint))
	}
	cp := codePoint - 0x10000
	return append(units, uint16(0xD800+(cp>>10)), uint16(0xDC00+(cp&0x3FF)))
}

// Java: IRFactory.isSupportedForInputLanguageMode(Feature)
func (f *IRFactory) IsSupportedForInputLanguageMode(feature parser.FeatureSet_Feature) bool {
	return f.config.LanguageMode().FeatureSet().Has(feature)
}

// Java: IRFactory.inStrictContext()
func (f *IRFactory) InStrictContext() bool {
	// TODO(johnlenz): in ECMASCRIPT5/6 is a "mixed" mode and we should track
	// the context that we are in, if we want to support it.
	return f.config.StrictMode().IsStrict()
}

// NormalizeNumber is the value of a NUMBER token.
//
// Go-only: language-adapter — a NUMBER token's value only holds ASCII
// characters (digits, '_', '.', exponent and radix letters, signs; see
// Scanner), so byte indices are Java char indices.
//
// Java: IRFactory.normalizeNumber(LiteralToken)
func (f *IRFactory) NormalizeNumber(token parser.LiteralToken) float64 {
	value := token.GetLiteralTokenBase().Value
	if strings.Contains(value, "_") {
		value = f.RemoveNumericSeparators(value, token)
	}
	location := token.GetTokenBase().Location
	length := len(value)
	jutil.CheckState(length > 0)
	jutil.CheckState(value[0] != '-' && value[0] != '+')
	if value[0] == '.' {
		return jutil.DoubleParseDouble("0" + value)
	} else if value[0] == '0' && length > 1 {
		switch value[1] {
		case '.', 'e', 'E':
			return jutil.DoubleParseDouble(value)
		case 'b', 'B':
			f.MaybeWarnForFeatureToken(token, parser.FeatureSet_Feature_BINARY_LITERALS)
			v := 0.0
			c := 1
			for c++; c < length; c++ {
				v = (v * 2) + float64(IRFactory_Binarydigit(uint16(value[c])))
			}
			return v
		case 'o', 'O':
			f.MaybeWarnForFeatureToken(token, parser.FeatureSet_Feature_OCTAL_LITERALS)
			v := 0.0
			c := 1
			for c++; c < length; c++ {
				v = (v * 8) + float64(IRFactory_Octaldigit(uint16(value[c])))
			}
			return v
		case 'x', 'X':
			v := 0.0
			c := 1
			for c++; c < length; c++ {
				v = (v * 0x10) + float64(IRFactory_Hexdigit(uint16(value[c])))
			}
			return v
		case '0', '1', '2', '3', '4', '5', '6', '7':
			v := 0.0
			c := 0
			for c++; c < length; c++ {
				digit := uint16(value[c])
				if IRFactory_IsOctalDigit(digit) {
					v = (v * 8) + float64(IRFactory_Octaldigit(digit))
				} else {
					f.errorReporter.Error(
						IRFactory_INVALID_OCTAL_DIGIT, f.sourceName, IRFactory_LinenoOfPosition(location.Start), IRFactory_CharnoOfPosition(location.Start))
					return 0
				}
			}
			if f.InStrictContext() {
				f.errorReporter.Error(
					IRFactory_INVALID_ES5_STRICT_OCTAL,
					f.sourceName,
					IRFactory_LinenoOfPosition(location.Start),
					IRFactory_CharnoOfPosition(location.Start))
			} else {
				f.errorReporter.Warning(
					IRFactory_INVALID_ES5_STRICT_OCTAL,
					f.sourceName,
					IRFactory_LinenoOfPosition(location.Start),
					IRFactory_CharnoOfPosition(location.Start))
			}
			return v
		case '8', '9':
			f.errorReporter.Error(
				IRFactory_INVALID_OCTAL_DIGIT, f.sourceName, IRFactory_LinenoOfPosition(location.Start), IRFactory_CharnoOfPosition(location.Start))
			return 0
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected character in number literal: " + value[1:2]))
		}
	} else if value == "0" {
		return 0
	} else {
		if length <= 15 && value[0] >= '1' && value[0] <= '9' {
			v := uint64(value[0] - '0')
			allDigits := true
			for i := 1; i < length; i++ {
				d := value[i] - '0'
				if d > 9 {
					allDigits = false
					break
				}
				v = v*10 + uint64(d)
			}
			if allDigits {
				return float64(v)
			}
		}
		return jutil.DoubleParseDouble(value)
	}
}

// NormalizeBigInt is the value of a BIGINT token.
//
// Go-only: language-adapter — a BIGINT token's value only holds ASCII
// characters (see NormalizeNumber).
//
// Java: IRFactory.normalizeBigInt(LiteralToken)
func (f *IRFactory) NormalizeBigInt(token parser.LiteralToken) *big.Int {
	value := token.GetLiteralTokenBase().Value
	n := strings.IndexByte(value, 'n')
	if n < 0 {
		panic(jutil.NewStringIndexOutOfBoundsException(0, n))
	}
	value = value[0:n]
	if strings.Contains(value, "_") {
		value = f.RemoveNumericSeparators(value, token)
	}
	length := len(value)
	jutil.CheckState(length > 0)
	jutil.CheckState(value[0] != '-' && value[0] != '+')
	if value[0] == '0' && length > 1 {
		switch value[1] {
		case 'b', 'B':
			f.MaybeWarnForFeatureToken(token, parser.FeatureSet_Feature_BINARY_LITERALS)
			return jutil.NewBigIntegerFromStringRadix(value[2:], 2)
		case 'o', 'O':
			f.MaybeWarnForFeatureToken(token, parser.FeatureSet_Feature_OCTAL_LITERALS)
			return jutil.NewBigIntegerFromStringRadix(value[2:], 8)
		case 'x', 'X':
			return jutil.NewBigIntegerFromStringRadix(value[2:], 16)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			panic(jutil.NewIllegalStateExceptionMsg("Nonzero BigInts can't have a leading zero"))
		default:
			panic(jutil.NewIllegalStateExceptionMsg(
				"Unexpected character in bigint literal: " + value[1:2]))
		}
	} else {
		return jutil.NewBigIntegerFromString(value)
	}
}

// Java: IRFactory.removeNumericSeparators(String, LiteralToken)
func (f *IRFactory) RemoveNumericSeparators(value string, token parser.LiteralToken) string {
	f.MaybeWarnForFeatureToken(token, parser.FeatureSet_Feature_NUMERIC_SEPARATOR)
	return jutil.StringReplace(value, "_", "")
}

// Java: IRFactory.binarydigit(char)
func IRFactory_Binarydigit(c uint16) int {
	if c >= '0' && c <= '1' {
		return int(c - '0')
	}
	panic(jutil.NewIllegalStateExceptionMsg("unexpected: " + jutil.StringFromChar(c)))
}

// Java: IRFactory.isOctalDigit(char)
func IRFactory_IsOctalDigit(c uint16) bool {
	return c >= '0' && c <= '7'
}

// Java: IRFactory.octaldigit(char)
func IRFactory_Octaldigit(c uint16) int {
	if IRFactory_IsOctalDigit(c) {
		return int(c - '0')
	}
	panic(jutil.NewIllegalStateExceptionMsg("unexpected: " + jutil.StringFromChar(c)))
}

// Java: IRFactory.hexdigit(char)
func IRFactory_Hexdigit(c uint16) int {
	switch c {
	case '0':
		return 0
	case '1':
		return 1
	case '2':
		return 2
	case '3':
		return 3
	case '4':
		return 4
	case '5':
		return 5
	case '6':
		return 6
	case '7':
		return 7
	case '8':
		return 8
	case '9':
		return 9
	case 'a', 'A':
		return 10
	case 'b', 'B':
		return 11
	case 'c', 'C':
		return 12
	case 'd', 'D':
		return 13
	case 'e', 'E':
		return 14
	case 'f', 'F':
		return 15
	default:
		panic(jutil.NewIllegalStateExceptionMsg("unexpected: " + jutil.StringFromChar(c)))
	}
}

// Java: IRFactory.transformBooleanTokenType(TokenType)
func IRFactory_TransformBooleanTokenType(token parser.TokenType) rhino.Token {
	switch token {
	case parser.TokenType_TRUE:
		return rhino.Token_TRUE
	case parser.TokenType_FALSE:
		return rhino.Token_FALSE
	default:
		panic(jutil.NewIllegalStateExceptionMsg(token.String()))
	}
}

// Java: IRFactory.transformUpdateTokenType(TokenType)
func IRFactory_TransformUpdateTokenType(token parser.TokenType) rhino.Token {
	switch token {
	case parser.TokenType_PLUS_PLUS:
		return rhino.Token_INC
	case parser.TokenType_MINUS_MINUS:
		return rhino.Token_DEC
	default:
		panic(jutil.NewIllegalStateExceptionMsg(token.String()))
	}
}

// Java: IRFactory.transformUnaryTokenType(TokenType)
func IRFactory_TransformUnaryTokenType(token parser.TokenType) rhino.Token {
	switch token {
	case parser.TokenType_BANG:
		return rhino.Token_NOT
	case parser.TokenType_TILDE:
		return rhino.Token_BITNOT
	case parser.TokenType_PLUS:
		return rhino.Token_POS
	case parser.TokenType_MINUS:
		return rhino.Token_NEG
	case parser.TokenType_DELETE:
		return rhino.Token_DELPROP
	case parser.TokenType_TYPEOF:
		return rhino.Token_TYPEOF
	case parser.TokenType_VOID:
		return rhino.Token_VOID
	default:
		panic(jutil.NewIllegalStateExceptionMsg(token.String()))
	}
}

// Java: IRFactory.transformBinaryTokenType(TokenType)
func IRFactory_TransformBinaryTokenType(token parser.TokenType) rhino.Token {
	switch token {
	case parser.TokenType_BAR:
		return rhino.Token_BITOR
	case parser.TokenType_CARET:
		return rhino.Token_BITXOR
	case parser.TokenType_AMPERSAND:
		return rhino.Token_BITAND
	case parser.TokenType_EQUAL_EQUAL:
		return rhino.Token_EQ
	case parser.TokenType_NOT_EQUAL:
		return rhino.Token_NE
	case parser.TokenType_OPEN_ANGLE:
		return rhino.Token_LT
	case parser.TokenType_LESS_EQUAL:
		return rhino.Token_LE
	case parser.TokenType_CLOSE_ANGLE:
		return rhino.Token_GT
	case parser.TokenType_GREATER_EQUAL:
		return rhino.Token_GE
	case parser.TokenType_LEFT_SHIFT:
		return rhino.Token_LSH
	case parser.TokenType_RIGHT_SHIFT:
		return rhino.Token_RSH
	case parser.TokenType_UNSIGNED_RIGHT_SHIFT:
		return rhino.Token_URSH
	case parser.TokenType_PLUS:
		return rhino.Token_ADD
	case parser.TokenType_MINUS:
		return rhino.Token_SUB
	case parser.TokenType_STAR:
		return rhino.Token_MUL
	case parser.TokenType_SLASH:
		return rhino.Token_DIV
	case parser.TokenType_PERCENT:
		return rhino.Token_MOD
	case parser.TokenType_STAR_STAR:
		return rhino.Token_EXPONENT
	case parser.TokenType_EQUAL_EQUAL_EQUAL:
		return rhino.Token_SHEQ
	case parser.TokenType_NOT_EQUAL_EQUAL:
		return rhino.Token_SHNE
	case parser.TokenType_IN:
		return rhino.Token_IN
	case parser.TokenType_INSTANCEOF:
		return rhino.Token_INSTANCEOF
	case parser.TokenType_COMMA:
		return rhino.Token_COMMA
	case parser.TokenType_EQUAL:
		return rhino.Token_ASSIGN
	case parser.TokenType_BAR_EQUAL:
		return rhino.Token_ASSIGN_BITOR
	case parser.TokenType_CARET_EQUAL:
		return rhino.Token_ASSIGN_BITXOR
	case parser.TokenType_AMPERSAND_EQUAL:
		return rhino.Token_ASSIGN_BITAND
	case parser.TokenType_LEFT_SHIFT_EQUAL:
		return rhino.Token_ASSIGN_LSH
	case parser.TokenType_RIGHT_SHIFT_EQUAL:
		return rhino.Token_ASSIGN_RSH
	case parser.TokenType_UNSIGNED_RIGHT_SHIFT_EQUAL:
		return rhino.Token_ASSIGN_URSH
	case parser.TokenType_PLUS_EQUAL:
		return rhino.Token_ASSIGN_ADD
	case parser.TokenType_MINUS_EQUAL:
		return rhino.Token_ASSIGN_SUB
	case parser.TokenType_STAR_EQUAL:
		return rhino.Token_ASSIGN_MUL
	case parser.TokenType_STAR_STAR_EQUAL:
		return rhino.Token_ASSIGN_EXPONENT
	case parser.TokenType_SLASH_EQUAL:
		return rhino.Token_ASSIGN_DIV
	case parser.TokenType_PERCENT_EQUAL:
		return rhino.Token_ASSIGN_MOD
	case parser.TokenType_OR:
		return rhino.Token_OR
	case parser.TokenType_AND:
		return rhino.Token_AND
	case parser.TokenType_QUESTION_QUESTION:
		return rhino.Token_COALESCE
	case parser.TokenType_OR_EQUAL:
		return rhino.Token_ASSIGN_OR
	case parser.TokenType_AND_EQUAL:
		return rhino.Token_ASSIGN_AND
	case parser.TokenType_QUESTION_QUESTION_EQUAL:
		return rhino.Token_ASSIGN_COALESCE
	default:
		panic(jutil.NewIllegalStateExceptionMsg(token.String()))
	}
}

// NewNode is a simple helper to create nodes and set the initial node
// properties.
//
// Java: IRFactory.newNode(Token)
func (f *IRFactory) NewNode(type_ rhino.Token) *rhino.Node {
	return new(rhino.Node).InitFromTemplate(type_, f.templateNode)
}

// Java: IRFactory.newNode(Token, Node)
func (f *IRFactory) NewNodeWithChild(type_ rhino.Token, child1 *rhino.Node) *rhino.Node {
	return new(rhino.Node).InitWithChildFromTemplate(type_, child1, f.templateNode)
}

// Java: IRFactory.newNode(Token, Node, Node)
func (f *IRFactory) NewNodeWith2Children(type_ rhino.Token, child1 *rhino.Node, child2 *rhino.Node) *rhino.Node {
	return new(rhino.Node).InitWith2ChildrenFromTemplate(type_, child1, child2, f.templateNode)
}

// Java: IRFactory.newNode(Token, Node, Node, Node)
func (f *IRFactory) NewNodeWith3Children(type_ rhino.Token, child1 *rhino.Node, child2 *rhino.Node, child3 *rhino.Node) *rhino.Node {
	return new(rhino.Node).InitWith3ChildrenFromTemplate(type_, child1, child2, child3, f.templateNode)
}

// Java: IRFactory.newStringNode(String)
func (f *IRFactory) NewStringNode(value string) *rhino.Node {
	return new(rhino.Node).InitStringFromTemplate(rhino.Token_STRINGLIT, value, f.templateNode)
}

// Java: IRFactory.newStringNode(Token, String)
func (f *IRFactory) NewStringNodeWithToken(type_ rhino.Token, value string) *rhino.Node {
	return new(rhino.Node).InitStringFromTemplate(type_, value, f.templateNode)
}

// MaybeAddNonJsDocComment attaches any pending JSDoc comments to the given
// node.
//
// Java: IRFactory.maybeAddNonJsDocComment(Node, SourcePosition)
func (f *IRFactory) MaybeAddNonJsDocComment(node *rhino.Node, start parser.SourcePosition) *rhino.Node {
	comment := f.ParseNonJSDocCommentAt(start, false)
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	return node
}

// Java: IRFactory.newTemplateLitStringNode(String, String)
func (f *IRFactory) NewTemplateLitStringNode(cooked *string, raw string) *rhino.Node {
	return new(rhino.Node).InitTemplateLitStringFromTemplate(cooked, raw, f.templateNode)
}

// Java: IRFactory.newNumberNode(Double)
func (f *IRFactory) NewNumberNode(value float64) *rhino.Node {
	return new(rhino.Node).InitFromTemplate(rhino.Token_NUMBER, f.templateNode).SetNumberRaw(value)
}

// Java: IRFactory.newBigIntNode(BigInteger)
func (f *IRFactory) NewBigIntNode(value *big.Int) *rhino.Node {
	return new(rhino.Node).InitBigIntFromTemplate(value, f.templateNode)
}

// Java: IRFactory.error(Node, String)
func (f *IRFactory) Error(rhinoNode *rhino.Node, message string) {
	f.errorReporter.Error(
		message, irFactory_nullableString(rhinoNode.GetSourceFileName()), rhinoNode.GetLineno(), rhinoNode.GetCharno())
}
