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
	"slices"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func (f *IRFactory) markCommentParsed(idx int) {
	if uint(idx) < uint(len(f.parsedCommentBits)*64) {
		f.parsedCommentBits[idx>>6] |= 1 << (uint(idx) & 63)
	}
}

func (f *IRFactory) markCommentPtrParsed(comment *parser.Comment) {
	src := f.jsdocTracker.source
	if idx := f.jsdocTracker.previousIndex; idx >= 0 && idx < len(src) && &src[idx] == comment {
		f.markCommentParsed(idx)
		return
	}
	for i := range src {
		if &src[i] == comment || src[i] == *comment {
			f.markCommentParsed(i)
			return
		}
	}
}

func (f *IRFactory) isCommentParsed(idx int) bool {
	return uint(idx) < uint(len(f.parsedCommentBits)*64) && (f.parsedCommentBits[idx>>6]&(1<<(uint(idx)&63))) != 0
}

// Java: final class IRFactory.CommentTracker
type IRFactory_CommentTracker struct {
	source        []parser.Comment
	index         int
	previousIndex int
	jsdocOnly     bool
}

// Java: IRFactory.CommentTracker.current()
func (t *IRFactory_CommentTracker) Current() *parser.Comment {
	if t.index >= len(t.source) {
		return nil
	}
	return &t.source[t.index]
}

func (t *IRFactory_CommentTracker) hasPendingBeforeOffset(offset int32) bool {
	return t.index < len(t.source) && t.source[t.index].Location.End.Offset <= offset
}

// Java: IRFactory.CommentTracker.advance()
func (t *IRFactory_CommentTracker) Advance() {
	t.previousIndex = t.index
	src := t.source
	jsdocOnly := t.jsdocOnly
	for {
		t.index++ // Always advance at least one element.
		if t.index >= len(src) || (src[t.index].Type == parser.Comment_Type_JSDOC) == jsdocOnly {
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

// Java: IRFactory.setFileOverviewJsDoc(Node)
func (f *IRFactory) SetFileOverviewJsDoc(irNode rhino.Node) {
	if len(f.licenses) == 0 {
		if f.firstFileoverview != nil {
			irNode.SetJSDocInfo(f.firstFileoverview)
		}
		return
	}

	var fileoverview *rhino.JSDocInfo_Builder
	if f.firstFileoverview == nil {
		fileoverview = rhino.JSDocInfo_Builder_()
	} else {
		fileoverview = f.firstFileoverview.ToBuilder()
	}

	if len(f.licenses) == 1 {
		fileoverview.RecordLicense(f.licenses[0])
	} else {
		fileoverview.RecordLicense(strings.Join(f.licenses, ""))
	}
	irNode.SetJSDocInfo(fileoverview.BuildWithAlways(false))
}

// HandlePossibleFileOverviewJsDoc returns true if the jsDocParser represents
// a fileoverview.
//
// Java: IRFactory.handlePossibleFileOverviewJsDoc(JsDocInfoParser)
func (f *IRFactory) HandlePossibleFileOverviewJsDoc(jsDocParser *JsDocInfoParser) bool {
	if lic, ok := jsDocParser.GetLicenseText(); ok && !slices.Contains(f.licenses, lic) {
		f.licenses = append(f.licenses, lic)
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
		if newFileoverview.HasAnySuppressions() {
			var supBuf [8]string
			suppressions := supBuf[:0]
			newFileoverview.ForEachSuppression(func(name string) bool {
				if !slices.Contains(suppressions, name) {
					suppressions = append(suppressions, name)
				}
				return true
			})
			merged.RecordSuppressionsSlice(suppressions)
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
	if !f.jsdocTracker.hasPendingBeforeOffset(tree.GetStart().Offset) {
		return nil
	}
	return f.parseJSDocInfoOnTreeAt(tree, f.linePos(tree.GetStart()))
}

func (f *IRFactory) parseJSDocInfoOnTreeAt(tree parser.ParseTree, startPos parser.LinePosition) *rhino.JSDocInfo {
	if !f.jsdocTracker.hasPendingBeforeOffset(startPos.Offset) {
		return nil
	}
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

	return f.ParseJSDocInfoFrom(f.GetJSDocCommentAt(startPos))
}

// Java: IRFactory.parseJSDocInfoOnToken(parser.Token)
func (f *IRFactory) ParseJSDocInfoOnToken(token parser.Token) *rhino.JSDocInfo {
	if !f.jsdocTracker.hasPendingBeforeOffset(token.Location.Start.Offset) {
		return nil
	}
	return f.ParseJSDocInfoFrom(f.GetJSDocCommentAt(f.linePos(token.GetStart())))
}

// Java: IRFactory.parseInlineJSDocAt(SourcePosition)
func (f *IRFactory) ParseInlineJSDocAt(pos parser.LinePosition) *rhino.JSDocInfo {
	if !f.jsdocTracker.hasPendingBeforeOffset(pos.Offset) {
		return nil
	}
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

func (f *IRFactory) hasAnyPendingCommentBeforeOffset(offset int32) bool {
	if f.nonJsdocTracker.hasPendingBeforeOffset(offset) {
		return true
	}
	return len(f.closureUnawareRanges) > 0 && f.jsdocTracker.hasPendingBeforeOffset(offset)
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
	if !f.hasAnyPendingCommentBeforeOffset(pos.Offset) {
		return nil
	}
	withinClosureUnawareCodeRange := f.WithinClosureUnawareCodeRange(int(pos.Line), int(pos.Column))
	if !f.HasAnyPendingCommentBefore(pos, withinClosureUnawareCodeRange) {
		return nil
	}
	includeAll := f.includeAllComments

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
	if !f.includeAllComments {
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

func (f *IRFactory) lineColumn(offset int) (int, int) {
	starts := f.lineStarts
	if starts == nil {
		return f.fileWithContent.LineColumn(offset)
	}
	o := int32(offset)
	i := f.lastLine
	if i >= len(starts) || o < starts[i] || (i+1 < len(starts) && starts[i+1] <= o) {
		if i+1 < len(starts) && starts[i+1] <= o && (i+2 >= len(starts) || o < starts[i+2]) {
			i++
		} else if i > 0 && i < len(starts) && o < starts[i] && starts[i-1] <= o {
			i--
		} else {
			i = parser.FindLine(starts, o)
		}
		f.lastLine = i
	}
	return i, offset - int(starts[i])
}

// linePos returns pos with its line and column.
//
// Go-only: parse-tree positions are offsets; the lines and columns come from
// the file's line table (see parser.SourcePosition).
func (f *IRFactory) linePos(pos parser.SourcePosition) parser.LinePosition {
	line, col := f.lineColumn(int(pos.Offset))
	return parser.LinePosition{Offset: pos.Offset, Line: int32(line), Column: int32(col)}
}

// Java: IRFactory.lineno(ParseTree)
//
// Go-only: lineno and charno are methods (Java: static), as positions are
// offsets and the line table is the file's.
func (f *IRFactory) Lineno(node parser.ParseTree) int {
	return f.LinenoOfPosition(node.GetParseTreeBase().Location.Start)
}

func checkTokenNotNull(token parser.Token) {
	if token.IsNull() {
		panic(jutil.NewNullPointerException())
	}
}

// Java: IRFactory.lineno(parser.Token)
func (f *IRFactory) LinenoOfToken(token parser.Token) int {
	checkTokenNotNull(token)
	return f.LinenoOfPosition(token.Location.Start)
}

// Java: IRFactory.lineno(SourcePosition)
func (f *IRFactory) LinenoOfPosition(location parser.SourcePosition) int {
	// location lines start at zero, our AST starts at 1.
	line, _ := f.lineColumn(int(location.Offset))
	return line + 1
}

// Java: IRFactory.charno(ParseTree)
func (f *IRFactory) Charno(node parser.ParseTree) int {
	return f.CharnoOfPosition(node.GetParseTreeBase().Location.Start)
}

// Java: IRFactory.charno(parser.Token)
func (f *IRFactory) CharnoOfToken(token parser.Token) int {
	checkTokenNotNull(token)
	return f.CharnoOfPosition(token.Location.Start)
}

// Java: IRFactory.charno(SourcePosition)
func (f *IRFactory) CharnoOfPosition(location parser.SourcePosition) int {
	_, col := f.lineColumn(int(location.Offset))
	return col
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

// MaybeAddNonJsDocComment attaches any pending JSDoc comments to the given
// node.
//
// Java: IRFactory.maybeAddNonJsDocComment(Node, SourcePosition)
func (f *IRFactory) MaybeAddNonJsDocComment(node rhino.Node, start parser.SourcePosition) rhino.Node {
	if !f.hasAnyPendingCommentBeforeOffset(start.Offset) {
		return node
	}
	comment := f.ParseNonJSDocCommentAt(f.linePos(start), false)
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	return node
}
