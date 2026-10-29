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
	"slices"

	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
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

	parsedCommentBits    []uint64
	parsedCommentBitsBuf [4]uint64

	licenses          []string
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

	closureUnawareRanges []irFactorySpan

	jsdocStream JsDocTokenStream
	jsdocParser JsDocInfoParser

	// Go-only: performance — where the factory takes its node ids (nil: the
	// store's counter); a parse worker passes its own (rhino.NodeAllocator).
	alloc                    *rhino.NodeAllocator
	transformDispatcherVal   IRFactory_TransformDispatcher
	classScopeVal            IRFactory_TransformDispatcher_ScopeTracker
	privateIdLhsOfInScopeVal IRFactory_TransformDispatcher_ScopeTracker

	// Go-only: performance — single-threaded line table + cursor, fast-path
	// flags, L1 string intern cache, and reusable literal scratch buffers.
	lineStarts         []int32
	lastLine           int
	includeAllComments bool
	reservedMode       uint8 // 0 = ES3 (nil), 1 = ES5, 2 = ES5_STRICT
	strCache           [128]strpool.ID
	litUnitsBuf        []uint16
	litUnitsScratch    [128]uint16
	litUTF16Buf        []uint16
	litUTF16Scratch    [128]uint16
	litByteBuf         []byte
	litByteScratch     [128]byte
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

// IRFactory_LineAndColumn represents a line and column number in a file.
//
// Java: record IRFactory.LineAndColumn(int lineNo, int column) implements Comparable<LineAndColumn>
type IRFactory_LineAndColumn struct {
	lineNo int
	column int
}

type irFactorySpan struct {
	start IRFactory_LineAndColumn
	end   IRFactory_LineAndColumn
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

// Java: IRFactory(StaticSourceFile, Config, ErrorReporter, ImmutableList<Comment>, SourceFile, JsDocInfoParser.JsDocSourceKind)
func NewIRFactory(
	sourceFile rhino.StaticSourceFile,
	config *Config,
	errorReporter rhino.ErrorReporter,
	comments []parser.Comment,
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
	if fileWithContent != nil {
		f.lineStarts = fileWithContent.LineStarts()
	}
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
	f.includeAllComments = config.JsDocParsingMode() == Config_JsDocParsing_INCLUDE_ALL_COMMENTS
	f.errorReporter = errorReporter
	f.transformDispatcherVal = IRFactory_TransformDispatcher{
		IRFactory:             f,
		classScope:            &f.classScopeVal,
		privateIdLhsOfInScope: &f.privateIdLhsOfInScopeVal,
	}
	f.transformDispatcher = &f.transformDispatcherVal

	if config.StrictMode().IsStrict() {
		f.reservedMode = 2
	} else if config.LanguageMode() == Config_LanguageMode_ECMASCRIPT3 {
		f.reservedMode = 0
	} else {
		f.reservedMode = 1
	}
	f.jsDocSourceKind = jsDocSourceKind
	return f
}

func (f *IRFactory) isReservedKeyword(identifier string) bool {
	switch f.reservedMode {
	case 0:
		return false
	case 1:
		switch identifier {
		case "class", "const", "enum", "export", "extends", "import", "super":
			return true
		}
		return false
	default:
		switch identifier {
		case "class", "const", "enum", "export", "extends", "import", "super",
			"implements", "interface", "let", "package", "private", "protected", "public", "static", "yield":
			return true
		}
		return false
	}
}

func (f *IRFactory) internString(s string) strpool.ID {
	if len(s) == 0 {
		return 1
	}
	if len(s) <= 32 {
		h := uint32(2166136261)
		for i := 0; i < len(s); i++ {
			h ^= uint32(s[i])
			h *= 16777619
		}
		slot := h & uint32(len(f.strCache)-1)
		if id := f.strCache[slot]; id != 0 && strpool.String(id) == s {
			return id
		}
		id := strpool.Intern(s)
		f.strCache[slot] = id
		return id
	}
	return strpool.Intern(s)
}

func (f *IRFactory) internStringBytes(b []byte) strpool.ID {
	if len(b) == 0 {
		return 1
	}
	if len(b) <= 32 {
		h := uint32(2166136261)
		for _, c := range b {
			h ^= uint32(c)
			h *= 16777619
		}
		slot := h & uint32(len(f.strCache)-1)
		if id := f.strCache[slot]; id != 0 && strpool.String(id) == string(b) {
			return id
		}
		id := strpool.InternBytes(b)
		f.strCache[slot] = id
		return id
	}
	return strpool.InternBytes(b)
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
		for i := range tree.SourceComments {
			comment := &tree.SourceComments[i]
			if (comment.Type == parser.Comment_Type_JSDOC || comment.Type == parser.Comment_Type_IMPORTANT) &&
				!irFactory.isCommentParsed(i) {
				useLicensesOnlyConfig := irFactory.WithinClosureUnawareCodeRange(
					int(comment.Location.Start.Line), int(comment.Location.Start.Column))
				irFactory.markCommentParsed(i)
				irFactory.HandlePossibleFileOverviewJsDocComment(comment, useLicensesOnlyConfig)
			}
		}

		if n.IsScript() && irFactory.hasAnyPendingCommentBeforeOffset(tree.Location.End.Offset) {
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
	loc := tree.GetParseTreeBase().Location
	startPos := f.linePos(loc.Start)
	info := f.parseJSDocInfoOnTreeAt(tree, startPos)
	var comment *rhino.NonJSDocComment
	if f.hasAnyPendingCommentBeforeOffset(startPos.Offset) {
		comment = f.ParseNonJSDocCommentAt(startPos, false)
	}
	isRootOfClosureUnawareSubTree := info != nil && info.IsClosureUnawareCode()

	if isRootOfClosureUnawareSubTree {
		if f.closureUnawareTemplateNode.IsNil() {
			f.closureUnawareTemplateNode = f.CreateClosureUnawareTemplateNode()
		}
		f.templateNode = f.closureUnawareTemplateNode
		f.addClosureUnawareRange(
			IRFactory_LineAndColumn_FromSourcePosition(startPos),
			IRFactory_LineAndColumn_FromSourcePosition(f.linePos(loc.End)),
		)
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
		f.setSourceInfoAt(node, startPos, loc.End)
		f.templateNode = f.defaultTemplateNode
		node = f.HideClosureUnawareCodeRoot(node)
	}
	f.setSourceInfoAt(node, startPos, loc.End)
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

func (f *IRFactory) addClosureUnawareRange(start, end IRFactory_LineAndColumn) {
	spans := f.closureUnawareRanges
	i, _ := slices.BinarySearchFunc(spans, start, func(s irFactorySpan, p IRFactory_LineAndColumn) int {
		return s.end.CompareTo(p)
	})
	j := i
	for j < len(spans) && spans[j].start.CompareTo(end) <= 0 {
		if spans[j].start.CompareTo(start) < 0 {
			start = spans[j].start
		}
		if spans[j].end.CompareTo(end) > 0 {
			end = spans[j].end
		}
		j++
	}
	f.closureUnawareRanges = slices.Replace(spans, i, j, irFactorySpan{start: start, end: end})
}

// Java: IRFactory.withinClosureUnawareCodeRange(int, int)
func (f *IRFactory) WithinClosureUnawareCodeRange(line int, lineColumnNo int) bool {
	spans := f.closureUnawareRanges
	if len(spans) == 0 {
		return false
	}
	pt := IRFactory_LineAndColumn{lineNo: line, column: lineColumnNo}
	i, _ := slices.BinarySearchFunc(spans, pt, func(s irFactorySpan, p IRFactory_LineAndColumn) int {
		return s.end.CompareTo(p)
	})
	return i < len(spans) && spans[i].start.CompareTo(pt) <= 0
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
	loc := tree.GetParseTreeBase().Location
	startPos := f.linePos(loc.Start)
	info := f.ParseInlineJSDocAt(startPos)
	var comment *rhino.NonJSDocComment
	if f.hasAnyPendingCommentBeforeOffset(startPos.Offset) {
		comment = f.ParseNonJSDocCommentAt(startPos, true)
	}

	node := f.transformDispatcher.Process(tree)

	if info != nil {
		node.SetJSDocInfo(info)
	}
	if comment != nil {
		node.SetNonJSDocComment(comment)
	}
	f.setSourceInfoAt(node, startPos, loc.End)
	return node
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
	if f.currentFileIsExterns || f.IsSupportedForInputLanguageMode(feature) {
		return
	}
	start := f.linePos(node.GetParseTreeBase().Location.Start)
	if f.WithinClosureUnawareCodeRange(int(start.Line), int(start.Column)) {
		return
	}
	f.errorReporter.Warning(
		IRFactory_LanguageFeatureWarningMessage(feature), f.sourceName, int(start.Line)+1, int(start.Column))
}

// Java: IRFactory.maybeWarnForFeature(parser.Token, Feature)
func (f *IRFactory) MaybeWarnForFeatureToken(token parser.Token, feature parser.FeatureSet_Feature) {
	f.features = f.features.With(feature)
	if f.currentFileIsExterns || f.IsSupportedForInputLanguageMode(feature) {
		return
	}
	checkTokenNotNull(token)
	start := f.linePos(token.Location.Start)
	if f.WithinClosureUnawareCodeRange(int(start.Line), int(start.Column)) {
		return
	}
	f.errorReporter.Warning(
		IRFactory_LanguageFeatureWarningMessage(feature), f.sourceName, int(start.Line)+1, int(start.Column))
}

// Java: IRFactory.maybeWarnForFeature(Node, Feature)
func (f *IRFactory) MaybeWarnForFeatureNode(node rhino.Node, feature parser.FeatureSet_Feature) {
	f.features = f.features.With(feature)
	if f.currentFileIsExterns || f.IsSupportedForInputLanguageMode(feature) {
		return
	}
	if f.WithinClosureUnawareCodeRange(node.GetLineno(), node.GetCharno()) {
		return
	}
	f.errorReporter.Warning(
		IRFactory_LanguageFeatureWarningMessage(feature), f.sourceName, node.GetLineno(), node.GetCharno())
}

func (f *IRFactory) setSourceInfoAt(node rhino.Node, startPos parser.LinePosition, end parser.SourcePosition) {
	if node.GetLineno() == -1 {
		node.SetLinenoCharno(int(startPos.Line)+1, int(startPos.Column))
		node.SetLength(int(end.Offset - startPos.Offset))
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
	checkTokenNotNull(token)
	location := token.Location
	f.SetSourceInfoFromPositions(irNode, location.Start, location.End)
}

// Java: IRFactory.setSourceInfo(Node, SourcePosition, SourcePosition)
func (f *IRFactory) SetSourceInfoFromPositions(node rhino.Node, start parser.SourcePosition, end parser.SourcePosition) {
	if node.GetLineno() == -1 {
		// If we didn't already set the line, then set it now. This avoids
		// cases like ParenthesizedExpression where we just return a previous
		// node, but don't want the new node to get its parent's line number.
		line, column := f.lineColumn(int(start.Offset))
		// location lines start at zero, our AST starts at 1.
		node.SetLinenoCharno(line+1, column)
		f.SetLength(node, start, end)
	}
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
	return f.alloc.Alloc().InitStringIDFromTemplate(rhino.Token_STRINGLIT, f.internString(value), f.templateNode)
}

// Java: IRFactory.newStringNode(Token, String)
func (f *IRFactory) NewStringNodeWithToken(type_ rhino.Token, value string) rhino.Node {
	return f.alloc.Alloc().InitStringIDFromTemplate(type_, f.internString(value), f.templateNode)
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
