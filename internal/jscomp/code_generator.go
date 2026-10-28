// Port of com/google/javascript/jscomp/CodeGenerator.java @ closure-compiler v20261021

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CodeGenerator generates JS code from a parse tree, sending it to the
// specified concrete CodeConsumer.
//
// Go-only: language-adapter — TypedCodeGenerator's single override
// (add(Node, Context) calling maybeAddTypeAnnotation(Node)) is folded into
// CodeGenerator via emitTypeAnnotations + registry so all internal calls
// dispatch directly without interface indirection.
type CodeGenerator struct {
	cc CodeConsumer

	outputCharsetEncoder *OutputCharsetEncoder // nil is null

	preferSingleQuotes      bool
	preserveTypeAnnotations bool
	printNonJSDocComments   bool

	// To distinguish between gents and non-gents mode so that we can turn
	// off checking the sanity of the source location of comments, and also
	// provide a different mode for comment printing between those two.
	gentsMode bool

	trustedStrings         bool
	quoteKeywordProperties bool
	useOriginalName        bool
	emitTypeAnnotations    bool
	outputFeatureSet       *parser.FeatureSet
	jsDocInfoPrinter       *JSDocInfoPrinter
	registry               *rhino.JSTypeRegistry

	// Go-only: language-adapter — reused UTF-16 buffer for escaping strings
	// larger than the stack buffer.
	scratchUnits []uint16
}

// NewCodeGenerator is the constructor used for cost estimation.
//
// Java: CodeGenerator(CodeConsumer)
func NewCodeGenerator(consumer *CodeConsumer) *CodeGenerator {
	g := &CodeGenerator{}
	g.InitCodeGenerator(consumer)
	return g
}

var (
	codeGenerator_jsDocInfoPrinterFalse = &JSDocInfoPrinter{useOriginalName: false, printDesc: false}
	codeGenerator_jsDocInfoPrinterTrue  = &JSDocInfoPrinter{useOriginalName: true, printDesc: false}
)

func (g *CodeGenerator) initDefaults() {
	g.outputCharsetEncoder = nil
	g.preferSingleQuotes = false
	g.trustedStrings = true
	g.preserveTypeAnnotations = false
	g.printNonJSDocComments = false
	g.gentsMode = false
	g.quoteKeywordProperties = false
	g.useOriginalName = false
	g.emitTypeAnnotations = false
	g.outputFeatureSet = parser.FeatureSet_BARE_MINIMUM
	g.jsDocInfoPrinter = codeGenerator_jsDocInfoPrinterFalse
	g.registry = nil
}

func (g *CodeGenerator) initForCostEstimation(costThreshold int) {
	g.cc = CodeConsumer{
		mode:               sinkCost,
		continueProcessing: true,
		maxCost:            costThreshold,
	}
	g.initDefaults()
}

func (g *CodeGenerator) initWithOptions(options *CompilerOptions) {
	if outputCharset := options.GetOutputCharset(); outputCharset != nil && outputCharset != jutil.StandardCharsets_US_ASCII {
		g.outputCharsetEncoder = NewOutputCharsetEncoder(outputCharset)
	} else {
		g.outputCharsetEncoder = nil
	}
	g.preferSingleQuotes = options.ShouldPreferSingleQuotes()
	g.trustedStrings = options.AssumeTrustedStrings()
	g.preserveTypeAnnotations = options.ShouldPreserveTypeAnnotations()
	g.printNonJSDocComments = options.GetPreserveNonJSDocComments()
	g.gentsMode = options.GetGentsMode()
	g.quoteKeywordProperties = options.ShouldQuoteKeywordProperties()
	g.useOriginalName = options.GetUseOriginalNamesInOutput()
	g.emitTypeAnnotations = false
	g.outputFeatureSet = options.GetOutputFeatureSet()
	if g.useOriginalName {
		g.jsDocInfoPrinter = codeGenerator_jsDocInfoPrinterTrue
	} else {
		g.jsDocInfoPrinter = codeGenerator_jsDocInfoPrinterFalse
	}
	g.registry = nil
}

// InitCodeGenerator is the CodeGenerator(CodeConsumer) constructor.
//
// Java: CodeGenerator(CodeConsumer)
func (g *CodeGenerator) InitCodeGenerator(consumer *CodeConsumer) {
	if consumer != nil {
		g.cc = *consumer
	}
	g.initDefaults()
}

// Java: CodeGenerator(CodeConsumer, CompilerOptions)
func NewCodeGeneratorWithOptions(consumer *CodeConsumer, options *CompilerOptions) *CodeGenerator {
	g := &CodeGenerator{}
	g.InitCodeGeneratorWithOptions(consumer, options)
	return g
}

// InitCodeGeneratorWithOptions is the CodeGenerator(CodeConsumer,
// CompilerOptions) constructor.
//
// Java: CodeGenerator(CodeConsumer, CompilerOptions)
func (g *CodeGenerator) InitCodeGeneratorWithOptions(consumer *CodeConsumer, options *CompilerOptions) {
	if consumer != nil {
		g.cc = *consumer
	}
	g.initWithOptions(options)
}

// Java: CodeGenerator.forCostEstimation(CodeConsumer)
func CodeGenerator_ForCostEstimation(consumer *CodeConsumer) *CodeGenerator {
	return NewCodeGenerator(consumer)
}

// TagAsTypeSummary inserts a top-level identifying file as .i.js generated
// typing file.
//
// Java: CodeGenerator.tagAsTypeSummary()
func (g *CodeGenerator) TagAsTypeSummary() {
	g.Add("/** @fileoverview @typeSummary */\n")
}

// TagAsStrict inserts a ECMASCRIPT 5 strict annotation.
//
// Java: CodeGenerator.tagAsStrict()
func (g *CodeGenerator) TagAsStrict() {
	g.Add("'use strict';")
	g.cc.EndLine()
}

// Java: CodeGenerator.printJSDocComment(Node, JSDocInfo)
func (g *CodeGenerator) PrintJSDocComment(node rhino.Node, jsDocInfo *rhino.JSDocInfo) {
	var jsdocAsString *string
	// In gents mode, we print NonJSDocInfo without special handling, as we
	// already have pre-filtered the comments properly.
	if g.gentsMode {
		jsdocAsString = jsDocInfo.GetOriginalCommentString()
	} else {
		s := g.jsDocInfoPrinter.Print(jsDocInfo)
		jsdocAsString = &s
	}
	// Don't print an empty jsdoc
	if jsdocAsString != nil && *jsdocAsString != "/** */ " {
		g.Add(*jsdocAsString)
		if !node.IsCast() {
			g.cc.EndLine()
		}
	}
}

// Java: CodeGenerator.printNonJSDocComment(Node, NonJSDocComment)
func (g *CodeGenerator) PrintNonJSDocComment(node rhino.Node, nonJSDocComment *rhino.NonJSDocComment) {
	nonJSDocCommentString := nonJSDocComment.GetCommentString()
	if nonJSDocCommentString != "" {
		g.AddNonJsDoc_nonTrailing(node, nonJSDocComment)
	}
}

// PrintLeadingCommentsInOrder prints Leading JSDocComments or
// NonJSDocComments for the given node in order, depending on their source
// location.
//
// Java: CodeGenerator.printLeadingCommentsInOrder(Node)
func (g *CodeGenerator) PrintLeadingCommentsInOrder(node rhino.Node) {
	if !g.preserveTypeAnnotations && !g.printNonJSDocComments {
		return
	}
	var jsDocInfo *rhino.JSDocInfo
	if g.preserveTypeAnnotations {
		jsDocInfo = node.GetJSDocInfo()
	}
	var nonJSDocComment *rhino.NonJSDocComment
	if g.printNonJSDocComments {
		nonJSDocComment = node.GetNonJSDocComment()
	}

	printJSDoc := jsDocInfo != nil
	printNonJSDoc := nonJSDocComment != nil
	if printJSDoc && printNonJSDoc {
		if jsDocInfo.GetOriginalCommentPosition() < nonJSDocComment.GetStartPosition().GetOffset() {
			g.PrintJSDocComment(node, jsDocInfo)
			g.PrintNonJSDocComment(node, nonJSDocComment)
		} else {
			g.PrintNonJSDocComment(node, nonJSDocComment)
			g.PrintJSDocComment(node, jsDocInfo)
		}
		return
	}

	if printJSDoc {
		g.PrintJSDocComment(node, jsDocInfo)
		return
	}

	if printNonJSDoc {
		g.PrintNonJSDocComment(node, nonJSDocComment)
		return
	}
}

// HasTrailingCommentOnSameLine returns true when a node has a trailing
// comment.
//
// Java: CodeGenerator.hasTrailingCommentOnSameLine(Node)
func (g *CodeGenerator) HasTrailingCommentOnSameLine(node rhino.Node) bool {
	if !g.printNonJSDocComments {
		return false
	}
	return node.GetTrailingNonJSDocCommentString() != ""
}

// Java: CodeGenerator.printTrailingComment(Node)
func (g *CodeGenerator) PrintTrailingComment(node rhino.Node) {
	if !g.printNonJSDocComments {
		return
	}
	nonJSDocComment := node.GetTrailingNonJSDocComment()
	if nonJSDocComment != nil {
		nonJSDocCommentString := node.GetTrailingNonJSDocCommentString()
		if nonJSDocCommentString != "" {
			g.AddNonJsDoctrailing(nonJSDocComment, g.HasTrailingCommentOnSameLine(node))
		}
	}
}

// Java: CodeGenerator.addNonJsDoc_nonTrailing(Node, NonJSDocComment)
func (g *CodeGenerator) AddNonJsDoc_nonTrailing(node rhino.Node, nonJSDocComment *rhino.NonJSDocComment) {
	content := nonJSDocComment.GetCommentString()
	commentEndPosition := nonJSDocComment.GetEndPosition()

	nodeLineNumber := node.GetLineno() - 1 // source lines are 1-indexed

	if nonJSDocComment.IsEndingAsLineComment() {
		// Non trailing line comments can not be on the same line as the node.
		jutil.CheckStateMsg(
			g.gentsMode || int(commentEndPosition.Line) < nodeLineNumber,
			"Non trailing line comments can not be on the same line as the node.")
		g.Add(content + "\n")
	} else {
		if nodeLineNumber == int(commentEndPosition.Line) {
			// e.g. ``` /* comment */ let x; ```
			g.Add(content + " ")
		} else {
			// e.g.
			// ```
			// /* comment */
			// let x;
			// ```
			g.Add(content + "\n")
		}
	}
}

// Java: CodeGenerator.addNonJsDoctrailing(NonJSDocComment, boolean)
func (g *CodeGenerator) AddNonJsDoctrailing(nonJSDocComment *rhino.NonJSDocComment, sameLine bool) {
	content := nonJSDocComment.GetCommentString()
	if nonJSDocComment.IsEndingAsLineComment() {
		// Trailing line comments *must* end with a `\n`. E.g.. `let x; //comment\n`
		g.Add(" " + content)
		if sameLine {
			g.cc.StartNewLine()
		}
	} else {
		if nonJSDocComment.IsInline() {
			// e.g. `foo(x /*comment*/);` is inline
			g.Add(" " + content)
		} else {
			// e.g. `let x; /*comment*/` is non-inline
			g.Add(" " + content)
			if sameLine {
				g.cc.StartNewLine()
			}
		}
	}
}

// Java: CodeGenerator.add(String)
func (g *CodeGenerator) Add(str string) {
	g.cc.Add(str)
}

// Java: CodeGenerator.add(Node)
func (g *CodeGenerator) AddNode(n rhino.Node) {
	g.AddNodeWithContext(n, CodeGenerator_Context_OTHER)
}

var codeGenerator_JSCOMP_SCOPE = rhino.QualifiedName_Of("$jscomp.scope")

// Java: CodeGenerator.add(Node, Context)
func (g *CodeGenerator) AddNodeWithContext(node rhino.Node, context CodeGenerator_Context) {
	if g.emitTypeAnnotations {
		g.MaybeAddTypeAnnotation(node)
	}
	g.AddNodeWithContextAndComments(node, context, true)
}

// AddNodeWithContextAndComments generates the current node.
//
// Java: CodeGenerator.add(Node, Context, boolean)
func (g *CodeGenerator) AddNodeWithContextAndComments(node rhino.Node, context CodeGenerator_Context, printComments bool) {
	cc := &g.cc
	if !cc.ContinueProcessing() {
		return
	}
	// If this node is actually a shadow host node, then print the shadow
	// content instead.
	shadow := node.GetClosureUnawareShadow()
	if !shadow.IsNil() {
		script := shadow.GetOnlyChild()
		sinkCall := script.GetOnlyChild().GetOnlyChild()
		jutil.CheckStateMsg1(sinkCall.IsCall(), "%s", sinkCall)
		g.AddNodeWithContextAndComments(sinkCall.GetLastChild(), context, printComments)
		return
	}
	if printComments && (g.preserveTypeAnnotations || g.printNonJSDocComments) {
		g.PrintLeadingCommentsInOrder(node)
	}
	cc.TrackLicenses(node)

	type_ := node.GetToken()
	opstr := NodeUtil_OpToStr(type_)
	first := node.GetFirstChild()
	last := node.GetLastChild()
	childCount := 0
	if !first.IsNil() {
		if first == last {
			childCount = 1
		} else if second := first.GetNext(); second == last {
			childCount = 2
		} else if third := second.GetNext(); third == last {
			childCount = 3
		} else if third.GetNext() == last {
			childCount = 4
		} else {
			childCount = 5
		}
	}

	// Handle all binary operators
	if opstr != nil && first != last {
		if childCount != 2 {
			jutil.CheckStateMsg2(
				false,
				"Bad binary operator \"%s\": expected 2 arguments but got %s",
				*opstr,
				node.GetChildCount())
		}
		p := g.Precedence(node)

		// For right-hand-side of operations, only pass context if it's
		// the IN_FOR_INIT_CLAUSE one.
		rhsContext := CodeGenerator_GetContextForNoInOperator(context)

		needsParens := (context == CodeGenerator_Context_START_OF_EXPR || context.AtArrowFunctionBody()) &&
			first.IsObjectPattern()
		if node.IsAssign() && needsParens {
			g.Add("(")
		}

		if NodeUtil_IsAssignmentOp(node) || type_ == rhino.Token_EXPONENT {
			// Assignment operators and '**' are the only right-associative
			// binary operators
			g.AddExpr(first, p+1, context)
			cc.AddOp(*opstr, true)
			g.AddExpr(last, p, rhsContext)
		} else {
			g.UnrollBinaryOperator(node, type_, *opstr, context, rhsContext, p, p+1)
		}

		if node.IsAssign() && needsParens {
			g.Add(")")
		}
		return
	}

	cc.StartSourceMapping(node)

	if !g.addStmtOrDecl(node, type_, first, last, childCount, context) {
		g.addExprOrTypeNode(node, type_, first, last, childCount, context)
	}

	if g.printNonJSDocComments {
		g.PrintTrailingComment(node)
	}

	cc.EndSourceMapping(node)
}

// Java: CodeGenerator.addAllSiblings(Node)
func (g *CodeGenerator) AddAllSiblings(n rhino.Node) {
	for c := n; !c.IsNil(); c = c.GetNext() {
		g.AddNode(c)
	}
}

// codeGenerator_nodeProp is the cast `(Node) n.getProp(...)`.
//
// Go-only: language-adapter — Java's reference cast of a property value.
func codeGenerator_nodeProp(value any) rhino.Node {
	if value == nil {
		return rhino.NilNode
	}
	return value.(rhino.Node)
}

// Java: CodeGenerator.precedence(Node)
func (g *CodeGenerator) Precedence(n rhino.Node) int {
	if n.IsCast() {
		return g.Precedence(n.GetFirstChild())
	}
	return NodeUtil_Precedence(n.GetToken())
}

// Java: CodeGenerator.addExpr(Node, int, Context)
func (g *CodeGenerator) AddExpr(n rhino.Node, minPrecedence int, context CodeGenerator_Context) {
	if g.OpRequiresParentheses(n, minPrecedence, context) {
		g.Add("(")
		g.AddNodeWithContext(n, CodeGenerator_Context_OTHER)
		g.Add(")")
	} else {
		g.AddNodeWithContext(n, context)
	}
}

// Java: CodeGenerator.opRequiresParentheses(Node, int, Context)
func (g *CodeGenerator) OpRequiresParentheses(n rhino.Node, minPrecedence int, context CodeGenerator_Context) bool {
	if context.InForInInitClause() && n.IsIn() {
		// make sure this operator 'in' isn't confused with the for-loop 'in'
		return true
	} else if NodeUtil_IsUnaryOperator(n) && g.IsFirstOperandOfExponentiationExpression(n) {
		// Unary operators are higher precedence than '**', but
		// ExponentiationExpression cannot expand to
		//     UnaryExpression ** ExponentiationExpression
		return true
	} else if CodeGenerator_IsLogicalANDorLogicalORChildOfNullishCoalesce(n) ||
		CodeGenerator_IsNullishCoalesceChildOfLogicalANDorLogicalOR(n) {
		// precedence is not enough here since using && or || with ?? without
		// parentheses is a syntax error as ?? expands directly to |
		return true
	} else if n.IsAssign() && n.GetParent().IsClass() {
		// Class declarations with assignments should be wrapped in parentheses.
		return true
	} else {
		return g.Precedence(n) < minPrecedence
	}
}

// Java: CodeGenerator.isLogicalANDorLogicalORChildOfNullishCoalesce(Node)
func CodeGenerator_IsLogicalANDorLogicalORChildOfNullishCoalesce(n rhino.Node) bool {
	parent := n.GetParent()
	logicalANDorLogicalOR := n.IsAnd() || n.IsOr()
	childOfNullishCoalesce := !parent.IsNil() && parent.IsNullishCoalesce()
	return logicalANDorLogicalOR && childOfNullishCoalesce
}

// Java: CodeGenerator.isNullishCoalesceChildOfLogicalANDorLogicalOR(Node)
func CodeGenerator_IsNullishCoalesceChildOfLogicalANDorLogicalOR(n rhino.Node) bool {
	parent := n.GetParent()
	childOfLogicalANDorLogicalOR := !parent.IsNil() && (parent.IsAnd() || parent.IsOr())
	return n.IsNullishCoalesce() && childOfLogicalANDorLogicalOR
}

// Java: CodeGenerator.isFirstOperandOfExponentiationExpression(Node)
func (g *CodeGenerator) IsFirstOperandOfExponentiationExpression(n rhino.Node) bool {
	parent := n.GetParent()
	return !parent.IsNil() && parent.IsExponent() && parent.GetFirstChild() == n
}

// Java: CodeGenerator.maybeAddAccessibilityModifier(Node)
func (g *CodeGenerator) MaybeAddAccessibilityModifier(n rhino.Node) {
	access := n.GetProp(rhino.Node_ACCESS_MODIFIER)
	if access != nil {
		switch access.(rhino.JSDocInfo_Visibility) {
		case rhino.JSDocInfo_Visibility_PRIVATE:
			g.Add("private ")
		case rhino.JSDocInfo_Visibility_PACKAGE:
			g.Add("package ")
		case rhino.JSDocInfo_Visibility_PROTECTED:
			g.Add("protected ")
		case rhino.JSDocInfo_Visibility_PUBLIC:
			g.Add("public ")
		case rhino.JSDocInfo_Visibility_INHERITED:
			g.Add("inherited ")
		default:
			g.Add(jutil.StringToLowerCaseASCII(access.(rhino.JSDocInfo_Visibility).String()) + " ")
		}
	}
}

// Java: CodeGenerator.maybeAddTypeDecl(Node)
func (g *CodeGenerator) MaybeAddTypeDecl(n rhino.Node) {
	if !n.GetDeclaredTypeExpression().IsNil() {
		g.Add(":")
		g.cc.MaybeInsertSpace()
		g.AddNode(n.GetDeclaredTypeExpression())
	}
}

// Java: CodeGenerator.maybeAddGenericTypes(Node)
func (g *CodeGenerator) MaybeAddGenericTypes(n rhino.Node) {
	generics := codeGenerator_nodeProp(n.GetProp(rhino.Node_GENERIC_TYPE_LIST))
	if !generics.IsNil() {
		g.AddNode(generics)
	}
}

// Java: CodeGenerator.maybeAddOptional(Node)
func (g *CodeGenerator) MaybeAddOptional(n rhino.Node) {
	if n.GetBooleanProp(rhino.Node_OPT_ES6_TYPED) {
		g.Add("?")
	}
}

// CodeGenerator_Context is information on the current context. Used for
// disambiguating special cases. For example, a "{" could indicate the start
// of an object literal or a block, depending on the current context.
//
// Java: CodeGenerator.Context
type CodeGenerator_Context int32

const (
	CodeGenerator_Context_STATEMENT CodeGenerator_Context = iota
	// a hack to resolve the else-clause ambiguity
	CodeGenerator_Context_BEFORE_DANGLING_ELSE
	CodeGenerator_Context_START_OF_EXPR
	// Are we inside the init clause of a for loop? If so, the containing
	// expression can't contain an in operator. Pass this context flag down
	// until we reach expressions which no longer have the limitation.
	CodeGenerator_Context_IN_FOR_INIT_CLAUSE
	// Handle object literals at the start of a non-block arrow function
	// body. This is only important when the first token after the "=>" is "{".
	CodeGenerator_Context_START_OF_ARROW_FN_BODY
	CodeGenerator_Context_START_OF_ARROW_FN_IN_FOR_INIT
	// nothing special to watch out for.
	CodeGenerator_Context_OTHER
)

var codeGenerator_Context_names = [...]string{
	"STATEMENT",
	"BEFORE_DANGLING_ELSE",
	"START_OF_EXPR",
	"IN_FOR_INIT_CLAUSE",
	"START_OF_ARROW_FN_BODY",
	"START_OF_ARROW_FN_IN_FOR_INIT",
	"OTHER",
}

var codeGenerator_Context_inForInitClause = [...]bool{
	false, // STATEMENT
	false, // BEFORE_DANGLING_ELSE
	false, // START_OF_EXPR
	true,  // IN_FOR_INIT_CLAUSE
	false, // START_OF_ARROW_FN_BODY
	true,  // START_OF_ARROW_FN_IN_FOR_INIT
	false, // OTHER
}

var codeGenerator_Context_atArrowFnBody = [...]bool{
	false, // STATEMENT
	false, // BEFORE_DANGLING_ELSE
	false, // START_OF_EXPR
	false, // IN_FOR_INIT_CLAUSE
	true,  // START_OF_ARROW_FN_BODY
	true,  // START_OF_ARROW_FN_IN_FOR_INIT
	false, // OTHER
}

// String is Context.name().
func (c CodeGenerator_Context) String() string { return codeGenerator_Context_names[c] }

// Java: CodeGenerator.Context.inForInInitClause()
func (c CodeGenerator_Context) InForInInitClause() bool {
	return codeGenerator_Context_inForInitClause[c]
}

// Java: CodeGenerator.Context.atArrowFunctionBody()
func (c CodeGenerator_Context) AtArrowFunctionBody() bool {
	return codeGenerator_Context_atArrowFnBody[c]
}

// Java: CodeGenerator.getContextForNonEmptyExpression(Context)
func CodeGenerator_GetContextForNonEmptyExpression(currentContext CodeGenerator_Context) CodeGenerator_Context {
	if currentContext == CodeGenerator_Context_BEFORE_DANGLING_ELSE {
		return CodeGenerator_Context_BEFORE_DANGLING_ELSE
	}
	return CodeGenerator_Context_OTHER
}

// CodeGenerator_GetContextForNoInOperator: if we're in a IN_FOR_INIT_CLAUSE,
// we can't permit in operators in the expression. Pass on the
// IN_FOR_INIT_CLAUSE flag through subexpressions.
//
// Java: CodeGenerator.getContextForNoInOperator(Context)
func CodeGenerator_GetContextForNoInOperator(context CodeGenerator_Context) CodeGenerator_Context {
	if context.InForInInitClause() {
		return context
	}
	return CodeGenerator_Context_OTHER
}

// CodeGenerator_GetContextForArrowFunctionBody: if we're at the start of an
// arrow function body, we need parentheses around object literals and
// object patterns. We also must pass the IN_FOR_INIT_CLAUSE flag into
// subexpressions.
//
// Java: CodeGenerator.getContextForArrowFunctionBody(Context)
func CodeGenerator_GetContextForArrowFunctionBody(context CodeGenerator_Context) CodeGenerator_Context {
	if context.InForInInitClause() {
		return CodeGenerator_Context_START_OF_ARROW_FN_IN_FOR_INIT
	}
	return CodeGenerator_Context_START_OF_ARROW_FN_BODY
}

// Java: CodeGenerator.processEnd(Node, Context)
func (g *CodeGenerator) ProcessEnd(n rhino.Node, context CodeGenerator_Context) {
	switch n.GetToken() {
	case rhino.Token_CLASS, rhino.Token_INTERFACE, rhino.Token_ENUM, rhino.Token_NAMESPACE:
		g.cc.EndClass(context == CodeGenerator_Context_STATEMENT)
	case rhino.Token_FUNCTION:
		if n.GetLastChild().IsEmpty() {
			g.cc.EndStatementWithSemiColon(true /* needSemiColon */, false /* hasTrailingCommentOnSameLine */)
		} else {
			g.cc.EndFunction(context == CodeGenerator_Context_STATEMENT)
		}
	case rhino.Token_DECLARE:
		if n.GetParent().GetToken() != rhino.Token_NAMESPACE_ELEMENTS {
			g.ProcessEnd(n.GetFirstChild(), context)
		}
	case rhino.Token_EXPORT:
		if n.GetParent().GetToken() != rhino.Token_NAMESPACE_ELEMENTS &&
			!n.GetFirstChild().IsDeclare() {
			g.ProcessEnd(n.GetFirstChild(), context)
		}
	case rhino.Token_COMPUTED_PROP:
		if n.HasOneChild() {
			g.cc.EndStatementWithSemiColon(true /* needSemiColon */, false /* hasTrailingCommentOnSameLine */)
		}
	case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF:
		if n.GetFirstChild().GetLastChild().IsEmpty() {
			g.cc.EndStatementWithSemiColon(true /* needSemiColon */, false /* hasTrailingCommentOnSameLine */)
		}
	case rhino.Token_MEMBER_VARIABLE_DEF:
		g.cc.EndStatementWithSemiColon(true /* needSemiColon */, false /* hasTrailingCommentOnSameLine */)
	default:
		if context == CodeGenerator_Context_STATEMENT {
			g.cc.EndStatement(false /* hasTrailingCommentOnSameLine */)
		}
	}
}
