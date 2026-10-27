// Port of com/google/javascript/jscomp/parsing/IRFactory.java @ closure-compiler v20261021

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
	"cmp"
	"math/big"
	"strings"

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
	reservedKeywords map[string]bool

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
	templateNode               rhino.Node
	defaultTemplateNode        rhino.Node
	closureUnawareTemplateNode rhino.Node

	jsdocTracker       *IRFactory_CommentTracker
	nonJsdocTracker    *IRFactory_CommentTracker
	jsdocTrackerVal    IRFactory_CommentTracker
	nonJsdocTrackerVal IRFactory_CommentTracker
	jsDocSourceKind    JsDocInfoParser_JsDocSourceKind

	currentFileIsExterns bool

	features   *parser.FeatureSet
	resultNode rhino.Node

	closureUnawareCodeRanges *jutil.TreeRangeSet[IRFactory_LineAndColumn]

	jsdocStream JsDocTokenStream
	jsdocParser JsDocInfoParser

	// Go-only: performance — where the factory takes its node ids (nil: the
	// store's counter); a parse worker passes its own (rhino.NodeAllocator).
	alloc                    *rhino.NodeAllocator
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
const IRFactory_NEW_DOT_TARGET_IN_CLASS_FIELD_OR_STATIC_BLOCK = "new.target is not supported in class field initializers or static blocks"
const IRFactory_UNDEFINED_LABEL = "undefined label \"%s\""

var IRFactory_ES5_RESERVED_KEYWORDS = map[string]bool{
	// From Section 7.6.1.2
	"class": true, "const": true, "enum": true, "export": true, "extends": true, "import": true, "super": true,
}

var IRFactory_ES5_STRICT_RESERVED_KEYWORDS = map[string]bool{
	// From Section 7.6.1.2
	"class":      true,
	"const":      true,
	"enum":       true,
	"export":     true,
	"extends":    true,
	"import":     true,
	"super":      true,
	"implements": true,
	"interface":  true,
	"let":        true,
	"package":    true,
	"private":    true,
	"protected":  true,
	"public":     true,
	"static":     true,
	"yield":      true,
}

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

// Java: IRFactory.LineAndColumn.fromSourcePosition(SourcePosition)
//
// Go-only: takes a LinePosition (positions are offsets, see
// parser.SourcePosition).
func IRFactory_LineAndColumn_FromSourcePosition(pos parser.LinePosition) IRFactory_LineAndColumn {
	return NewIRFactory_LineAndColumn(int(pos.Line), int(pos.Column))
}

// Java: IRFactory.LineAndColumn.compareTo(LineAndColumn)
func (l IRFactory_LineAndColumn) CompareTo(o IRFactory_LineAndColumn) int {
	if c := cmp.Compare(l.lineNo, o.lineNo); c != 0 {
		return c
	}
	return cmp.Compare(l.column, o.column)
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
	f := &IRFactory{}
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
func (t *IRFactory_CommentTracker) HasPendingCommentBefore(pos parser.LinePosition) bool {
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
func (f *IRFactory) CreateTemplateNode() rhino.Node {
	// The Node type choice is arbitrary.
	templateNode := rhino.NewNode(rhino.Token_SCRIPT)
	templateNode.SetStaticSourceFile(f.sourceFile)
	return templateNode
}

// Java: IRFactory.createClosureUnawareTemplateNode()
func (f *IRFactory) CreateClosureUnawareTemplateNode() rhino.Node {
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
	return irFactoryTransformTree(tree, sourceFile, config, errorReporter, file, nil)
}

// irFactoryTransformTree is IRFactory_TransformTree with the factory's
// node allocator.
func irFactoryTransformTree(
	tree *parser.ProgramTree,
	sourceFile rhino.StaticSourceFile,
	config *Config,
	errorReporter rhino.ErrorReporter,
	file *parser.SourceFile,
	alloc *rhino.NodeAllocator) *IRFactory {
	jsDocSourceKind := JsDocInfoParser_JsDocSourceKind_NORMAL
	if sourceFile.IsTypeScriptSource() {
		jsDocSourceKind = JsDocInfoParser_JsDocSourceKind_TSICKLE
	}
	irFactory := NewIRFactory(
		sourceFile, config, errorReporter, tree.SourceComments, file, jsDocSourceKind)
	irFactory.alloc = alloc
	irFactory.jsdocParser.alloc = alloc

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
			endOfFilePos := irFactory.linePos(tree.Location.End)
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
func (f *IRFactory) GetResultNode() rhino.Node {
	return f.resultNode
}

// Java: IRFactory.getFeatures()
func (f *IRFactory) GetFeatures() *parser.FeatureSet {
	return f.features
}

// Java: IRFactory.setFileOverviewJsDoc(Node)
func (f *IRFactory) SetFileOverviewJsDoc(irNode rhino.Node) {
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
func (f *IRFactory) TransformBlock(node parser.ParseTree) rhino.Node {
	irNode := f.Transform(node)
	if irNode.IsBlock() {
		return irNode
	}

	var newBlock rhino.Node
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
	if newFileoverview == f.firstFileoverview {
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
func (f *IRFactory) GetJSDocCommentAt(pos parser.LinePosition) *parser.Comment {
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

	return f.ParseJSDocInfoFrom(f.GetJSDocCommentAt(f.linePos(tree.GetStart())))
}

// Java: IRFactory.parseJSDocInfoOnToken(parser.Token)
func (f *IRFactory) ParseJSDocInfoOnToken(token parser.Token) *rhino.JSDocInfo {
	return f.ParseJSDocInfoFrom(f.GetJSDocCommentAt(f.linePos(token.GetStart())))
}

// Java: IRFactory.parseInlineJSDocAt(SourcePosition)
func (f *IRFactory) ParseInlineJSDocAt(pos parser.LinePosition) *rhino.JSDocInfo {
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
	if c := cmp.Compare(a.Location.Start.Line, b.Location.Start.Line); c != 0 {
		return c
	}
	return cmp.Compare(a.Location.Start.Column, b.Location.Start.Column)
}

// Java: IRFactory.hasAnyPendingCommentBefore(SourcePosition, boolean)
func (f *IRFactory) HasAnyPendingCommentBefore(pos parser.LinePosition, withinClosureUnawareCodeRange bool) bool {
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
func (f *IRFactory) ParseNonJSDocCommentAt(pos parser.LinePosition, isInline bool) *rhino.NonJSDocComment {
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
	tokenEnd parser.LinePosition, possibleNextTokenStart parser.LinePosition) *rhino.NonJSDocComment {
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
func (f *IRFactory) Transform(tree parser.ParseTree) rhino.Node {
	info := f.ParseJSDocInfoOnTree(tree)
	comment := f.ParseNonJSDocCommentAt(f.linePos(tree.GetStart()), false)
	isRootOfClosureUnawareSubTree := info != nil && info.IsClosureUnawareCode()

	if isRootOfClosureUnawareSubTree {
		if f.closureUnawareTemplateNode.IsNil() {
			f.closureUnawareTemplateNode = f.CreateClosureUnawareTemplateNode()
		}
		if f.closureUnawareCodeRanges == nil {
			f.closureUnawareCodeRanges = jutil.NewTreeRangeSet(irFactory_LineAndColumn_compare)
		}
		f.templateNode = f.closureUnawareTemplateNode
		nextClosureUnawareRange := tree.GetParseTreeBase().Location
		nextClosureUnawareRangeStart := IRFactory_LineAndColumn_FromSourcePosition(f.linePos(nextClosureUnawareRange.Start))
		nextClosureUnawareRangeEnd := IRFactory_LineAndColumn_FromSourcePosition(f.linePos(nextClosureUnawareRange.End))
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
func (f *IRFactory) HideClosureUnawareCodeRoot(fn rhino.Node) rhino.Node {
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
func (f *IRFactory) MaybeInjectCastNode(node parser.ParseTree, info *rhino.JSDocInfo, irNode rhino.Node) rhino.Node {
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
func (f *IRFactory) TransformNodeWithInlineComments(tree parser.ParseTree) rhino.Node {
	info := f.ParseInlineJSDocAt(f.linePos(tree.GetStart()))
	comment := f.ParseNonJSDocCommentAt(f.linePos(tree.GetStart()), true)

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

// linePos returns pos with its line and column.
//
// Go-only: parse-tree positions are offsets; the lines and columns come from
// the file's line table (see parser.SourcePosition).
func (f *IRFactory) linePos(pos parser.SourcePosition) parser.LinePosition {
	return f.fileWithContent.LinePosition(pos)
}

// Java: IRFactory.lineno(ParseTree)
//
// Go-only: lineno and charno are methods (Java: static), as positions are
// offsets and the line table is the file's.
func (f *IRFactory) Lineno(node parser.ParseTree) int {
	return f.LinenoOfPosition(node.GetParseTreeBase().Location.Start)
}

// Java: IRFactory.lineno(parser.Token)
func (f *IRFactory) LinenoOfToken(token parser.Token) int {
	return f.LinenoOfPosition(token.GetTokenBase().Location.Start)
}

// Java: IRFactory.lineno(SourcePosition)
func (f *IRFactory) LinenoOfPosition(location parser.SourcePosition) int {
	// location lines start at zero, our AST starts at 1.
	return f.fileWithContent.Line(location) + 1
}

// Java: IRFactory.charno(ParseTree)
func (f *IRFactory) Charno(node parser.ParseTree) int {
	return f.CharnoOfPosition(node.GetParseTreeBase().Location.Start)
}

// Java: IRFactory.charno(parser.Token)
func (f *IRFactory) CharnoOfToken(token parser.Token) int {
	return f.CharnoOfPosition(token.GetTokenBase().Location.Start)
}

// Java: IRFactory.charno(SourcePosition)
func (f *IRFactory) CharnoOfPosition(location parser.SourcePosition) int {
	return int(f.linePos(location).Column)
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
	start := f.linePos(node.GetParseTreeBase().Location.Start)
	if f.WithinClosureUnawareCodeRange(int(start.Line), int(start.Column)) {
		return
	}
	// Externs should not produce warnings about the input language mode.
	if f.currentFileIsExterns {
		return
	}
	if !f.IsSupportedForInputLanguageMode(feature) {
		f.errorReporter.Warning(
			IRFactory_LanguageFeatureWarningMessage(feature), f.sourceName, f.Lineno(node), f.Charno(node))
	}
}

// Java: IRFactory.maybeWarnForFeature(parser.Token, Feature)
func (f *IRFactory) MaybeWarnForFeatureToken(token parser.Token, feature parser.FeatureSet_Feature) {
	f.features = f.features.With(feature)
	start := f.linePos(token.GetTokenBase().Location.Start)
	if f.WithinClosureUnawareCodeRange(int(start.Line), int(start.Column)) {
		return
	}
	// Externs should not produce warnings about the input language mode.
	if f.currentFileIsExterns {
		return
	}
	if !f.IsSupportedForInputLanguageMode(feature) {
		f.errorReporter.Warning(
			IRFactory_LanguageFeatureWarningMessage(feature), f.sourceName, f.LinenoOfToken(token), f.CharnoOfToken(token))
	}
}

// Java: IRFactory.maybeWarnForFeature(Node, Feature)
func (f *IRFactory) MaybeWarnForFeatureNode(node rhino.Node, feature parser.FeatureSet_Feature) {
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
func (f *IRFactory) SetSourceInfoFromNode(node rhino.Node, ref rhino.Node) {
	node.SetLinenoCharno(ref.GetLineno(), ref.GetCharno())
	f.SetLengthFrom(node, ref)
}

// Java: IRFactory.setSourceInfo(Node, ParseTree)
func (f *IRFactory) SetSourceInfo(irNode rhino.Node, node parser.ParseTree) {
	if irNode.GetLineno() == -1 {
		location := node.GetParseTreeBase().Location
		f.SetSourceInfoFromPositions(irNode, location.Start, location.End)
	}
}

// Java: IRFactory.setSourceInfo(Node, parser.Token)
func (f *IRFactory) SetSourceInfoFromToken(irNode rhino.Node, token parser.Token) {
	location := token.GetTokenBase().Location
	f.SetSourceInfoFromPositions(irNode, location.Start, location.End)
}

// Java: IRFactory.setSourceInfo(Node, SourcePosition, SourcePosition)
func (f *IRFactory) SetSourceInfoFromPositions(node rhino.Node, start parser.SourcePosition, end parser.SourcePosition) {
	if node.GetLineno() == -1 {
		// If we didn't already set the line, then set it now. This avoids
		// cases like ParenthesizedExpression where we just return a previous
		// node, but don't want the new node to get its parent's line number.
		line, column := f.fileWithContent.LineColumn(int(start.Offset))
		// location lines start at zero, our AST starts at 1.
		node.SetLinenoCharno(line+1, column)
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
	lineno := int(node.Location.Start.Line) + 1
	charno := int(node.Location.Start.Column)
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
		f.errorReporter)
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
	lineno := int(node.Location.Start.Line) + 1
	charno := int(node.Location.Start.Column)

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
		f.errorReporter)
	return parser_.ParseInlineTypeDoc()
}

// SetLength sets the length on the node if we're in IDE mode.
//
// Java: IRFactory.setLength(Node, SourcePosition, SourcePosition)
func (f *IRFactory) SetLength(node rhino.Node, start parser.SourcePosition, end parser.SourcePosition) {
	node.SetLength(int(end.Offset - start.Offset))
}

// Java: IRFactory.setLengthFrom(Node, Node)
func (f *IRFactory) SetLengthFrom(node rhino.Node, ref rhino.Node) {
	node.SetLength(ref.GetLength())
}

// NewNode is a simple helper to create nodes and set the initial node
// properties.
//
// Java: IRFactory.newNode(Token)
func (f *IRFactory) NewNode(type_ rhino.Token) rhino.Node {
	return f.alloc.Alloc().InitFromTemplate(type_, f.templateNode)
}

// Java: IRFactory.newNode(Token, Node)
func (f *IRFactory) NewNodeWithChild(type_ rhino.Token, child1 rhino.Node) rhino.Node {
	return f.alloc.Alloc().InitWithChildFromTemplate(type_, child1, f.templateNode)
}

// Java: IRFactory.newNode(Token, Node, Node)
func (f *IRFactory) NewNodeWith2Children(type_ rhino.Token, child1 rhino.Node, child2 rhino.Node) rhino.Node {
	return f.alloc.Alloc().InitWith2ChildrenFromTemplate(type_, child1, child2, f.templateNode)
}

// Java: IRFactory.newNode(Token, Node, Node, Node)
func (f *IRFactory) NewNodeWith3Children(type_ rhino.Token, child1 rhino.Node, child2 rhino.Node, child3 rhino.Node) rhino.Node {
	return f.alloc.Alloc().InitWith3ChildrenFromTemplate(type_, child1, child2, child3, f.templateNode)
}

// Java: IRFactory.newStringNode(String)
func (f *IRFactory) NewStringNode(value string) rhino.Node {
	return f.alloc.Alloc().InitStringFromTemplate(rhino.Token_STRINGLIT, value, f.templateNode)
}

// Java: IRFactory.newStringNode(Token, String)
func (f *IRFactory) NewStringNodeWithToken(type_ rhino.Token, value string) rhino.Node {
	return f.alloc.Alloc().InitStringFromTemplate(type_, value, f.templateNode)
}

// MaybeAddNonJsDocComment attaches any pending JSDoc comments to the given
// node.
//
// Java: IRFactory.maybeAddNonJsDocComment(Node, SourcePosition)
func (f *IRFactory) MaybeAddNonJsDocComment(node rhino.Node, start parser.SourcePosition) rhino.Node {
	comment := f.ParseNonJSDocCommentAt(f.linePos(start), false)
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	return node
}

// Java: IRFactory.newTemplateLitStringNode(String, String)
func (f *IRFactory) NewTemplateLitStringNode(cooked *string, raw string) rhino.Node {
	return f.alloc.Alloc().InitTemplateLitStringFromTemplate(cooked, raw, f.templateNode)
}

// Java: IRFactory.newNumberNode(Double)
func (f *IRFactory) NewNumberNode(value float64) rhino.Node {
	return f.alloc.Alloc().InitFromTemplate(rhino.Token_NUMBER, f.templateNode).SetNumberRaw(value)
}

// Java: IRFactory.newBigIntNode(BigInteger)
func (f *IRFactory) NewBigIntNode(value *big.Int) rhino.Node {
	return f.alloc.Alloc().InitBigIntFromTemplate(value, f.templateNode)
}

// Java: IRFactory.error(Node, String)
func (f *IRFactory) Error(rhinoNode rhino.Node, message string) {
	f.errorReporter.Error(
		message, irFactory_nullableString(rhinoNode.GetSourceFileName()), rhinoNode.GetLineno(), rhinoNode.GetCharno())
}
