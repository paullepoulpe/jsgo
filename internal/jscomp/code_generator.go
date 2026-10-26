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
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const codeGenerator_LT_ESCAPED = "\\x3c"
const codeGenerator_GT_ESCAPED = "\\x3e"

// CodeGenerator generates codes from a parse tree, sending it to the
// specified CodeConsumer.
//
// Go-only: language-adapter — a class with a subclass (TypedCodeGenerator,
// which overrides add(Node, Context)): this interface declares the
// non-private methods, and CodeGeneratorBase holds the fields and the base
// implementations. Every call that Java dispatches virtually goes through
// CodeGeneratorBase.self.
type CodeGenerator interface {
	TagAsTypeSummary()
	TagAsStrict()
	PrintLeadingCommentsInOrder(node rhino.Node)
	PrintTrailingComment(node rhino.Node)
	Add(str string)
	AddNode(n rhino.Node)
	AddNodeWithContext(node rhino.Node, context CodeGenerator_Context)
	AddNodeWithContextAndComments(node rhino.Node, context CodeGenerator_Context, printComments bool)
	AddList(firstInList rhino.Node)
	AddListWithSeparator(firstInList rhino.Node, separator string)
	AddListWithContextAndSeparator(firstInList rhino.Node, isArrayOrFunctionArgument bool, lhsContext CodeGenerator_Context, separator string)
	AddStringKey(n rhino.Node)
	AddObjectPattern(n rhino.Node)
	AddArrayList(firstInList rhino.Node)
	AddCaseBody(caseBody rhino.Node)
	AddAllSiblings(n rhino.Node)
	RegexpEscape(s string) string

	GetCodeGeneratorBase() *CodeGeneratorBase
}

// CodeGeneratorBase holds CodeGenerator's fields and base implementations.
type CodeGeneratorBase struct {
	self CodeGenerator

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
	outputFeatureSet       *parser.FeatureSet
	jsDocInfoPrinter       *JSDocInfoPrinter

	// Go-only: language-adapter — the UTF-16 view of the string being
	// escaped (Java indexes its chars), reused across calls.
	scratchUnits []uint16
}

// utf16Units is the UTF-16 view of s, in the reused scratch buffer when s
// is ASCII.
//
// Go-only: language-adapter — see scratchUnits.
func (g *CodeGeneratorBase) utf16Units(s string) []uint16 {
	units := g.scratchUnits[:0]
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			return jutil.StringToUTF16(s)
		}
		units = append(units, uint16(s[i]))
	}
	g.scratchUnits = units
	return units
}

// GetCodeGeneratorBase gives access to the base fields through the
// interface.
func (g *CodeGeneratorBase) GetCodeGeneratorBase() *CodeGeneratorBase { return g }

// NewCodeGenerator is the private constructor used for cost estimation.
//
// Java: CodeGenerator(CodeConsumer)
func NewCodeGenerator(consumer CodeConsumer) *CodeGeneratorBase {
	g := &CodeGeneratorBase{}
	g.InitCodeGeneratorBase(g, consumer)
	return g
}

var (
	codeGenerator_jsDocInfoPrinterFalse = &JSDocInfoPrinter{useOriginalName: false, printDesc: false}
	codeGenerator_jsDocInfoPrinterTrue  = &JSDocInfoPrinter{useOriginalName: true, printDesc: false}
)

// InitCodeGeneratorBase is the CodeGenerator(CodeConsumer) constructor;
// self is the object being constructed.
//
// Java: CodeGenerator(CodeConsumer)
func (g *CodeGeneratorBase) InitCodeGeneratorBase(self CodeGenerator, consumer CodeConsumer) {
	g.self = self
	g.cc = consumer
	g.outputCharsetEncoder = nil
	g.preferSingleQuotes = false
	g.trustedStrings = true
	g.preserveTypeAnnotations = false
	g.printNonJSDocComments = false
	g.gentsMode = false
	g.quoteKeywordProperties = false
	g.useOriginalName = false
	g.outputFeatureSet = parser.FeatureSet_BARE_MINIMUM
	g.jsDocInfoPrinter = codeGenerator_jsDocInfoPrinterFalse
}

// Java: CodeGenerator(CodeConsumer, CompilerOptions)
func NewCodeGeneratorWithOptions(consumer CodeConsumer, options *CompilerOptions) *CodeGeneratorBase {
	g := &CodeGeneratorBase{}
	g.InitCodeGeneratorBaseWithOptions(g, consumer, options)
	return g
}

// InitCodeGeneratorBaseWithOptions is the CodeGenerator(CodeConsumer,
// CompilerOptions) constructor; self is the object being constructed.
//
// Java: CodeGenerator(CodeConsumer, CompilerOptions)
func (g *CodeGeneratorBase) InitCodeGeneratorBaseWithOptions(self CodeGenerator, consumer CodeConsumer, options *CompilerOptions) {
	g.self = self
	g.cc = consumer

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
	g.outputFeatureSet = options.GetOutputFeatureSet()
	if g.useOriginalName {
		g.jsDocInfoPrinter = codeGenerator_jsDocInfoPrinterTrue
	} else {
		g.jsDocInfoPrinter = codeGenerator_jsDocInfoPrinterFalse
	}
}

// Java: CodeGenerator.forCostEstimation(CodeConsumer)
func CodeGenerator_ForCostEstimation(consumer CodeConsumer) CodeGenerator {
	return NewCodeGenerator(consumer)
}

// TagAsTypeSummary inserts a top-level identifying file as .i.js generated
// typing file.
//
// Java: CodeGenerator.tagAsTypeSummary()
func (g *CodeGeneratorBase) TagAsTypeSummary() {
	g.self.Add("/** @fileoverview @typeSummary */\n")
}

// TagAsStrict inserts a ECMASCRIPT 5 strict annotation.
//
// Java: CodeGenerator.tagAsStrict()
func (g *CodeGeneratorBase) TagAsStrict() {
	g.self.Add("'use strict';")
	g.cc.EndLine()
}

// Java: CodeGenerator.printJSDocComment(Node, JSDocInfo)
func (g *CodeGeneratorBase) PrintJSDocComment(node rhino.Node, jsDocInfo *rhino.JSDocInfo) {
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
		g.self.Add(*jsdocAsString)
		if !node.IsCast() {
			g.cc.EndLine()
		}
	}
}

// Java: CodeGenerator.printNonJSDocComment(Node, NonJSDocComment)
func (g *CodeGeneratorBase) PrintNonJSDocComment(node rhino.Node, nonJSDocComment *rhino.NonJSDocComment) {
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
func (g *CodeGeneratorBase) PrintLeadingCommentsInOrder(node rhino.Node) {
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
func (g *CodeGeneratorBase) HasTrailingCommentOnSameLine(node rhino.Node) bool {
	if !g.printNonJSDocComments {
		return false
	}
	return node.GetTrailingNonJSDocCommentString() != ""
}

// Java: CodeGenerator.printTrailingComment(Node)
func (g *CodeGeneratorBase) PrintTrailingComment(node rhino.Node) {
	// print any trailing nonJSDoc comment attached to this node
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

// Java: CodeGenerator.add(String)
func (g *CodeGeneratorBase) Add(str string) {
	g.cc.Add(str)
}

// Java: CodeGenerator.add(Node)
func (g *CodeGeneratorBase) AddNode(n rhino.Node) {
	g.self.AddNodeWithContext(n, CodeGenerator_Context_OTHER)
}

var codeGenerator_JSCOMP_SCOPE = rhino.QualifiedName_Of("$jscomp.scope")

// Java: CodeGenerator.add(Node, Context)
func (g *CodeGeneratorBase) AddNodeWithContext(node rhino.Node, context CodeGenerator_Context) {
	g.self.AddNodeWithContextAndComments(node, context, true)
}

// codeGenerator_NODE_CLASS is Node.class.getName().
const codeGenerator_NODE_CLASS = "com.google.javascript.rhino.Node"

// AddNodeWithContextAndComments generates the current node.
//
// Java: CodeGenerator.add(Node, Context, boolean)
func (g *CodeGeneratorBase) AddNodeWithContextAndComments(node rhino.Node, context CodeGenerator_Context, printComments bool) {
	cc := g.cc
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
		g.self.AddNodeWithContextAndComments(sinkCall.GetLastChild(), context, printComments)
		return
	}
	if printComments && (g.preserveTypeAnnotations || g.printNonJSDocComments) {
		g.self.PrintLeadingCommentsInOrder(node)
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
			g.self.Add("(")
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
			g.self.Add(")")
		}
		return
	}

	cc.StartSourceMapping(node)

	switch type_ {
	case rhino.Token_TRY:
		jutil.CheckState(first.GetNext().IsBlock() && !first.GetNext().HasMoreThanOneChild())
		jutil.CheckState(childCount >= 2 && childCount <= 3)

		g.self.Add("try")
		g.self.AddNode(first)

		// second child contains the catch block, or nothing if there
		// isn't a catch block
		catchblock := first.GetNext().GetFirstChild()
		if !catchblock.IsNil() {
			g.self.AddNode(catchblock)
		}

		if childCount == 3 {
			cc.MaybeInsertSpace()
			g.self.Add("finally")
			g.self.AddNode(last)
		}
	case rhino.Token_CATCH:
		jutil.CheckStateMsg1(childCount == 2, "%s", node)
		cc.MaybeInsertSpace()
		g.self.Add("catch")
		cc.MaybeInsertSpace()

		if !first.IsEmpty() {
			// optional catch binding
			g.self.Add("(")
			g.self.AddNode(first)
			g.self.Add(")")
		}

		g.self.AddNode(last)
	case rhino.Token_THROW:
		jutil.CheckStateMsg1(childCount == 1, "%s", node)
		g.self.Add("throw")
		cc.MaybeInsertSpace()
		g.self.AddNode(first)

		// Must have a ';' after a throw statement, otherwise safari can't
		// parse this.
		cc.EndStatementWithSemiColon(true /* needSemiColon */, g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_RETURN:
		g.self.Add("return")
		if childCount == 1 {
			cc.MaybeInsertSpace()
			if g.preserveTypeAnnotations && first.GetJSDocInfo() != nil {
				g.self.Add("(")
				g.self.AddNode(first)
				g.self.Add(")")
			} else {
				g.self.AddNode(first)
			}
		} else {
			jutil.CheckStateMsg1(childCount == 0, "%s", node)
		}
		cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_VAR:
		g.self.Add("var ")
		g.self.AddListWithContextAndSeparator(first, false, CodeGenerator_GetContextForNoInOperator(context), ",")
		if node.GetParent().IsNil() || NodeUtil_IsStatement(node) {
			cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
		}
	case rhino.Token_CONST:
		g.self.Add("const ")
		g.self.AddListWithContextAndSeparator(first, false, CodeGenerator_GetContextForNoInOperator(context), ",")
		if node.GetParent().IsNil() || NodeUtil_IsStatement(node) {
			cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
		}
	case rhino.Token_LET:
		g.self.Add("let ")
		g.self.AddListWithContextAndSeparator(first, false, CodeGenerator_GetContextForNoInOperator(context), ",")
		if node.GetParent().IsNil() || NodeUtil_IsStatement(node) {
			cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
		}
	case rhino.Token_LABEL_NAME:
		jutil.CheckStateMsg1(node.GetString() != "", "%s", node)
		g.AddIdentifier(node.GetString())
	case rhino.Token_DESTRUCTURING_LHS:
		g.self.AddNode(first)
		if first != last {
			jutil.CheckStateMsg1(childCount == 2, "%s", node)
			cc.AddOp("=", true)
			g.AddExpr(last, NodeUtil_Precedence(rhino.Token_ASSIGN), CodeGenerator_GetContextForNoInOperator(context))
		}
	case rhino.Token_NAME:
		if g.useOriginalName && node.GetOriginalName() != nil {
			g.AddIdentifier(*node.GetOriginalName())
		} else {
			g.AddIdentifier(node.GetString())
		}
		g.MaybeAddOptional(node)
		g.MaybeAddTypeDecl(node)

		if !first.IsNil() && !first.IsEmpty() {
			jutil.CheckStateMsg1(childCount == 1, "%s", node)
			cc.AddOp("=", true)
			g.AddExpr(first, NodeUtil_Precedence(rhino.Token_ASSIGN), CodeGenerator_GetContextForNoInOperator(context))
		}
	case rhino.Token_ARRAYLIT:
		g.self.Add("[")
		g.self.AddArrayList(first)
		g.self.Add("]")
	case rhino.Token_ARRAY_PATTERN:
		g.self.Add("[")
		g.self.AddArrayList(first)
		g.self.Add("]")
		g.MaybeAddTypeDecl(node)
	case rhino.Token_PARAM_LIST:
		// If this is the list for a non-TypeScript arrow function with one
		// simple name param.
		if node.GetParent().IsArrowFunction() &&
			node.HasOneChild() &&
			first.IsName() &&
			!g.gentsMode {
			g.self.AddNode(first)
		} else {
			g.self.Add("(")
			g.self.AddList(first)
			g.self.Add(")")
		}
	case rhino.Token_DEFAULT_VALUE:
		g.self.AddNode(first)
		g.MaybeAddTypeDecl(node)
		cc.AddOp("=", true)
		g.AddExpr(first.GetNext(), 1, CodeGenerator_Context_OTHER)
	case rhino.Token_COMMA:
		jutil.CheckStateMsg1(childCount == 2, "%s", node)
		g.UnrollBinaryOperator(
			node, rhino.Token_COMMA, ",", context, CodeGenerator_GetContextForNoInOperator(context), 0, 0)
	case rhino.Token_NUMBER:
		jutil.CheckStateMsg1(childCount == 0, "%s", node)
		cc.AddNumber(node.GetDouble(), node)
	case rhino.Token_BIGINT:
		jutil.CheckStateMsg1(childCount == 0, "%s", node)
		cc.AddBigInt(node.GetBigInt())
	case rhino.Token_TYPEOF, rhino.Token_VOID, rhino.Token_NOT, rhino.Token_BITNOT, rhino.Token_POS, rhino.Token_NEG:
		// All of these unary operators are right-associative
		jutil.CheckStateMsg1(childCount == 1, "%s", node)
		cc.AddOp(NodeUtil_OpToStrNoFail(type_), false)
		g.AddExpr(first, NodeUtil_Precedence(type_), CodeGenerator_Context_OTHER)
	case rhino.Token_HOOK:
		if childCount != 3 {
			jutil.CheckStateMsg2(false, "%s wrong number of children: %s", node, node.GetChildCount())
		}
		p := NodeUtil_Precedence(type_)
		rhsContext := CodeGenerator_GetContextForNoInOperator(context)
		g.AddExpr(first, p+1, context)
		cc.AddOp("?", true)
		g.AddExpr(first.GetNext(), 1, rhsContext)
		cc.AddOp(":", true)
		g.AddExpr(last, 1, rhsContext)
	case rhino.Token_REGEXP:
		if !first.IsStringLit() || !last.IsStringLit() {
			panic(jutil.NewJavaLangError("Expected children to be strings"))
		}

		regexp := g.self.RegexpEscape(first.GetString())

		// I only use one .add because whitespace matters
		if childCount == 2 {
			g.self.Add(regexp + last.GetString())
		} else {
			jutil.CheckStateMsg1(childCount == 1, "%s", node)
			g.self.Add(regexp)
		}
	case rhino.Token_FUNCTION:
		if node.GetClassName() != codeGenerator_NODE_CLASS {
			panic(jutil.NewJavaLangError("Unexpected Node subclass."))
		}
		jutil.CheckStateMsg1(childCount == 3, "%s", node)
		if node.IsArrowFunction() {
			g.AddArrowFunction(node, first, last, context)
		} else {
			g.AddFunction(node, first, last, context)
		}
	case rhino.Token_ITER_REST, rhino.Token_OBJECT_REST:
		g.self.Add("...")
		g.self.AddNode(first)
		g.MaybeAddTypeDecl(node)
	case rhino.Token_ITER_SPREAD, rhino.Token_OBJECT_SPREAD:
		g.self.Add("...")
		g.AddExpr(first, NodeUtil_Precedence(type_), CodeGenerator_Context_OTHER)
	case rhino.Token_EXPORT:
		g.self.Add("export")
		if node.GetBooleanProp(rhino.Node_EXPORT_DEFAULT) {
			g.self.Add("default")
		}
		if node.GetBooleanProp(rhino.Node_EXPORT_ALL_FROM) {
			g.self.Add("*")
			jutil.CheckStateMsg1(!first.IsNil() && first.IsEmpty(), "%s", node)
		} else {
			g.self.AddNode(first)
		}
		if childCount == 2 {
			g.self.Add("from")
			g.self.AddNode(last)
		}
		g.ProcessEnd(first, context)
	case rhino.Token_IMPORT:
		g.self.Add("import")

		second := first.GetNext()
		if !first.IsEmpty() {
			g.self.AddNode(first)
			if !second.IsEmpty() {
				cc.ListSeparator()
			}
		}
		if !second.IsEmpty() {
			g.self.AddNode(second)
		}
		if !first.IsEmpty() || !second.IsEmpty() {
			g.self.Add("from")
		}
		g.self.AddNode(last)
		cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_EXPORT_SPECS, rhino.Token_IMPORT_SPECS:
		g.self.Add("{")
		for c := first; !c.IsNil(); c = c.GetNext() {
			if c != first {
				cc.ListSeparator()
			}
			g.self.AddNode(c)
		}
		g.self.Add("}")
	case rhino.Token_EXPORT_SPEC, rhino.Token_IMPORT_SPEC:
		g.self.AddNode(first)
		if node.IsShorthandProperty() && first.GetString() == last.GetString() {
			break
		}
		g.self.Add("as")
		g.self.AddNode(last)
	case rhino.Token_IMPORT_STAR:
		g.self.Add("*")
		g.self.Add("as")
		g.self.Add(node.GetString())
	case rhino.Token_DYNAMIC_IMPORT:
		g.self.Add("import(")
		g.AddExpr(first, NodeUtil_Precedence(type_), context)
		g.self.Add(")")
	case rhino.Token_IMPORT_META:
		g.self.Add("import.meta")
	// CLASS -> NAME,EXPR|EMPTY,BLOCK
	case rhino.Token_CLASS:
		{
			jutil.CheckStateMsg1(childCount == 3, "%s", node)
			classNeedsParens := (context == CodeGenerator_Context_START_OF_EXPR)
			if classNeedsParens {
				g.self.Add("(")
			}

			name := first
			superClass := first.GetNext()
			members := last

			g.self.Add("class")
			if !name.IsEmpty() {
				g.self.AddNode(name)
			}

			g.MaybeAddGenericTypes(first)

			if !superClass.IsEmpty() {
				g.self.Add("extends")

				// Parentheses are required for a comma expression or an
				// assignment expression.
				g.AddExpr(superClass, 1, CodeGenerator_Context_OTHER)
			}

			interfaces := codeGenerator_nodeProp(node.GetProp(rhino.Node_IMPLEMENTS))
			if !interfaces.IsNil() {
				g.self.Add("implements")
				child := interfaces.GetFirstChild()
				g.self.AddNode(child)
				for {
					child = child.GetNext()
					if child.IsNil() {
						break
					}
					g.self.Add(",")
					cc.MaybeInsertSpace()
					g.self.AddNode(child)
				}
			}
			g.self.AddNode(members)
			cc.EndClass(context == CodeGenerator_Context_STATEMENT)

			if classNeedsParens {
				g.self.Add(")")
			}
		}
	case rhino.Token_CLASS_MEMBERS, rhino.Token_INTERFACE_MEMBERS, rhino.Token_NAMESPACE_ELEMENTS:
		cc.BeginBlock()
		for c := first; !c.IsNil(); c = c.GetNext() {
			g.self.AddNode(c)
			g.ProcessEnd(c, context)
			cc.EndLine()
		}
		cc.EndBlockWithShouldEndLine(false)
	case rhino.Token_ENUM_MEMBERS:
		cc.BeginBlock()
		for c := first; !c.IsNil(); c = c.GetNext() {
			g.self.AddNode(c)
			if !c.GetNext().IsNil() {
				g.self.Add(",")
			}
			cc.EndLine()
		}
		cc.EndBlockWithShouldEndLine(false)
	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_VARIABLE_DEF:
		jutil.CheckState(
			node.GetParent().IsObjectLit() ||
				node.GetParent().IsClassMembers() ||
				node.GetParent().IsInterfaceMembers() ||
				node.GetParent().IsRecordType() ||
				node.GetParent().IsIndexSignature())

		g.MaybeAddAccessibilityModifier(node)
		if node.IsStaticMember() {
			g.self.Add("static ")
		}

		if node.IsMemberFunctionDef() && node.GetFirstChild().IsAsyncFunction() {
			g.self.Add("async ")
		}

		if !node.IsMemberVariableDef() && node.GetFirstChild().IsGeneratorFunction() {
			jutil.CheckStateMsg1(type_ == rhino.Token_MEMBER_FUNCTION_DEF, "%s", node)
			g.self.Add("*")
		}

		switch type_ {
		case rhino.Token_GETTER_DEF:
			// Get methods have no parameters.
			jutil.CheckStateMsg1(!first.GetSecondChild().HasChildren(), "%s", node)
			g.self.Add("get ")
		case rhino.Token_SETTER_DEF:
			// Set methods have one parameter.
			jutil.CheckStateMsg1(first.GetSecondChild().HasOneChild(), "%s", node)
			g.self.Add("set ")
		case rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_VARIABLE_DEF:
			// nothing to do.
		default:
		}

		// The name is on the GET or SET node.
		name := node.GetString()
		if node.IsMemberVariableDef() {
			g.self.Add(node.GetString())
			g.MaybeAddOptional(node)
			g.MaybeAddTypeDecl(node)
		} else {
			jutil.CheckStateMsg1(childCount == 1, "%s", node)
			jutil.CheckStateMsg1(first.IsFunction(), "%s", first)

			// The function referenced by the definition should always be
			// unnamed.
			jutil.CheckStateMsg1(first.GetFirstChild().GetString() == "", "%s", first)

			fn := first
			parameters := fn.GetSecondChild()
			body := fn.GetLastChild()

			// Add the property name.
			if !node.IsQuotedStringKey() &&
				(rhino.TokenStream_IsJSIdentifier(name) || node.IsPrivateIdentifier()) &&
				// do not encode literally any non-literal characters that
				// were Unicode escaped.
				NodeUtil_IsLatin(name) {
				g.self.Add(name)
				g.MaybeAddGenericTypes(fn.GetFirstChild())
			} else {
				// Determine if the string is a simple number.
				d := CodeGenerator_GetSimpleNumber(name)
				if !math.IsNaN(d) {
					cc.AddNumber(d, node)
				} else {
					g.AddJsString(node)
				}
			}
			g.MaybeAddOptional(fn)
			g.self.AddNode(parameters)
			g.MaybeAddTypeDecl(fn)
			g.self.AddNode(body)
		}
	case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_FIELD_DEF:
		jutil.CheckState(node.GetParent().IsClassMembers())
		if node.GetBooleanProp(rhino.Node_STATIC_MEMBER) {
			g.self.Add("static ")
		}
		var init rhino.Node
		switch type_ {
		case rhino.Token_MEMBER_FIELD_DEF:
			propertyName := node.GetString()
			g.self.Add(propertyName)
			init = first
		case rhino.Token_COMPUTED_FIELD_DEF:
			g.self.Add("[")
			// Must use addExpr() with a priority of 1, because comma
			// expressions aren't allowed.
			// https://www.ecma-international.org/ecma-262/9.0/index.html#prod-ComputedPropertyName
			g.AddExpr(first, 1, CodeGenerator_Context_OTHER)
			g.self.Add("]")
			init = node.GetSecondChild()
		default:
		}
		if !init.IsNil() {
			cc.AddOp("=", true)
			g.AddExpr(init, 1, CodeGenerator_Context_OTHER)
		}
		g.self.Add(";")
	case rhino.Token_SCRIPT, rhino.Token_MODULE_BODY, rhino.Token_BLOCK, rhino.Token_ROOT:
		if node.GetClassName() != codeGenerator_NODE_CLASS {
			panic(jutil.NewJavaLangError("Unexpected Node subclass."))
		}
		if node.HasParent() {
			staticBlock := node.IsBlock() && node.GetParent().IsClassMembers()
			if staticBlock {
				g.self.Add("static")
			}
		}
		// A BLOCK marked as synthetic is not a real JS block with {} around
		// it in the JS code. It just represents a span of statements that
		// need to be kept together.
		preserveBlock := node.IsBlock() && !node.IsSyntheticBlock()
		if preserveBlock {
			cc.BeginBlock()
		}

		preferLineBreaks := type_ == rhino.Token_SCRIPT ||
			(type_ == rhino.Token_BLOCK && !preserveBlock && node.GetParent().IsScript())
		for c := first; !c.IsNil(); c = c.GetNext() {
			g.self.AddNodeWithContext(c, CodeGenerator_Context_STATEMENT)

			if c.IsFunction() || c.IsClass() {
				cc.MaybeLineBreak()
			}

			// Prefer to break lines in between top-level statements
			// because top-level statements are more homogeneous.
			if preferLineBreaks {
				cc.NotePreferredLineBreak()
			}
		}
		if preserveBlock {
			cc.EndBlockWithShouldEndLine(cc.BreakAfterBlockFor(node, context == CodeGenerator_Context_STATEMENT))
		}
	case rhino.Token_FOR:
		jutil.CheckStateMsg1(childCount == 4, "%s", node)
		g.self.Add("for")
		cc.MaybeInsertSpace()
		g.self.Add("(")
		if NodeUtil_IsNameDeclaration(first) {
			g.self.AddNodeWithContext(first, CodeGenerator_Context_IN_FOR_INIT_CLAUSE)
		} else {
			g.AddExpr(first, 0, CodeGenerator_Context_IN_FOR_INIT_CLAUSE)
		}
		g.self.Add(";")
		if !first.GetNext().IsEmpty() {
			cc.MaybeInsertSpace()
		}
		g.self.AddNode(first.GetNext())
		g.self.Add(";")
		if !first.GetNext().GetNext().IsEmpty() {
			cc.MaybeInsertSpace()
		}
		g.self.AddNode(first.GetNext().GetNext())
		g.self.Add(")")
		g.AddNonEmptyStatement(last, CodeGenerator_GetContextForNonEmptyExpression(context), false)
	case rhino.Token_FOR_IN:
		jutil.CheckStateMsg1(childCount == 3, "%s", node)
		g.self.Add("for")
		cc.MaybeInsertSpace()
		g.self.Add("(")
		g.self.AddNode(first)
		g.self.Add("in")
		g.self.AddNode(first.GetNext())
		g.self.Add(")")
		g.AddNonEmptyStatement(last, CodeGenerator_GetContextForNonEmptyExpression(context), false)
	case rhino.Token_FOR_OF:
		jutil.CheckStateMsg1(childCount == 3, "%s", node)
		g.self.Add("for")
		cc.MaybeInsertSpace()
		g.self.Add("(")
		g.self.AddNode(first)
		cc.MaybeInsertSpace()
		g.self.Add("of")
		cc.MaybeInsertSpace()
		// the iterable must be an AssignmentExpression
		g.AddExpr(first.GetNext(), NodeUtil_Precedence(rhino.Token_ASSIGN), CodeGenerator_Context_OTHER)
		g.self.Add(")")
		g.AddNonEmptyStatement(last, CodeGenerator_GetContextForNonEmptyExpression(context), false)
	case rhino.Token_FOR_AWAIT_OF:
		jutil.CheckStateMsg1(childCount == 3, "%s", node)
		g.self.Add("for await")
		cc.MaybeInsertSpace()
		g.self.Add("(")
		g.self.AddNode(first)
		cc.MaybeInsertSpace()
		g.self.Add("of")
		cc.MaybeInsertSpace()
		// the iterable must be an AssignmentExpression
		g.AddExpr(first.GetNext(), NodeUtil_Precedence(rhino.Token_ASSIGN), CodeGenerator_Context_OTHER)
		g.self.Add(")")
		g.AddNonEmptyStatement(last, CodeGenerator_GetContextForNonEmptyExpression(context), false)
	case rhino.Token_DO:
		jutil.CheckStateMsg1(childCount == 2, "%s", node)
		g.self.Add("do")
		g.AddNonEmptyStatement(first, CodeGenerator_Context_OTHER, false)
		cc.MaybeInsertSpace()
		g.self.Add("while")
		cc.MaybeInsertSpace()
		g.self.Add("(")
		g.self.AddNode(last)
		g.self.Add(")")
		cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_WHILE:
		jutil.CheckStateMsg1(childCount == 2, "%s", node)
		g.self.Add("while")
		cc.MaybeInsertSpace()
		g.self.Add("(")
		g.self.AddNode(first)
		g.self.Add(")")
		g.AddNonEmptyStatement(last, CodeGenerator_GetContextForNonEmptyExpression(context), false)
	case rhino.Token_EMPTY:
		jutil.CheckStateMsg1(childCount == 0, "%s", node)
	case rhino.Token_OPTCHAIN_GETPROP:
		g.AddExpr(first, NodeUtil_Precedence(type_), context)
		if node.IsOptionalChainStart() {
			g.self.Add("?.")
		} else {
			g.self.Add(".")
		}
		g.AddGetpropIdentifier(node)
	case rhino.Token_GETPROP:
		// This attempts to convert rewritten aliased code back to the
		// original code, such as when using goog.scope(). See
		// ScopedAliases.java for the original code.
		// NOTE: OPTCHAIN_GETPROP case doesn't need this logic, because it
		// only applies to qualified names.
		if g.useOriginalName && node.GetOriginalName() != nil {
			// The ScopedAliases pass will convert variable assignments and
			// function declarations to assignments to GETPROP nodes, like
			// $jscomp.scope.SOME_VAR = 3;. This attempts to rewrite it back
			// to the original code.
			if codeGenerator_JSCOMP_SCOPE.Matches(node.GetFirstChild()) && node.GetParent().IsAssign() {
				g.self.Add("var ")
			}
			g.AddGetpropIdentifier(node)
			break
		}
		// We need parentheses to distinguish
		// `a?.b.c` from `(a?.b).c`
		breakOutOfOptionalChain := NodeUtil_IsOptChainNode(first)
		// `2.toString()` is invalid - it must be `(2).toString()`
		needsParens := first.IsNumber() || breakOutOfOptionalChain
		if needsParens {
			g.self.Add("(")
		}
		g.AddExpr(first, NodeUtil_Precedence(type_), context)
		if needsParens {
			g.self.Add(")")
		}
		if g.quoteKeywordProperties && rhino.TokenStream_IsKeyword(node.GetString()) {
			// NOTE: We don't have to worry about quoting keyword properties
			// in the OPTCHAIN_GETPROP case above, because we only need to
			// quote keywords for ES3-compatible output.
			//
			// Must be a single call to `add` otherwise the generator will
			// add a trailing space.
			g.self.Add("[\"" + node.GetString() + "\"]")
		} else {
			g.self.Add(".")
			g.AddGetpropIdentifier(node)
		}
	case rhino.Token_OPTCHAIN_GETELEM:
		if childCount != 2 {
			jutil.CheckStateMsg2(
				false,
				"Bad GETELEM node: Expected 2 children but got %s. For node: %s",
				node.GetChildCount(),
				node)
		}
		g.AddExpr(first, NodeUtil_Precedence(type_), context)
		if node.IsOptionalChainStart() {
			g.self.Add("?.")
		}
		g.self.Add("[")
		g.self.AddNode(first.GetNext())
		g.self.Add("]")
	case rhino.Token_GETELEM:
		if childCount != 2 {
			jutil.CheckStateMsg2(
				false,
				"Bad GETELEM node: Expected 2 children but got %s. For node: %s",
				node.GetChildCount(),
				node)
		}
		needsParens := NodeUtil_IsOptChainNode(first)
		if needsParens {
			g.self.Add("(")
		}
		g.AddExpr(first, NodeUtil_Precedence(type_), context)
		if needsParens {
			g.self.Add(")")
		}
		g.self.Add("[")
		g.self.AddNode(first.GetNext())
		g.self.Add("]")
	case rhino.Token_WITH:
		jutil.CheckStateMsg1(childCount == 2, "%s", node)
		g.self.Add("with(")
		g.self.AddNode(first)
		g.self.Add(")")
		g.AddNonEmptyStatement(last, CodeGenerator_GetContextForNonEmptyExpression(context), false)
	case rhino.Token_INC, rhino.Token_DEC:
		jutil.CheckStateMsg1(childCount == 1, "%s", node)
		o := "--"
		if type_ == rhino.Token_INC {
			o = "++"
		}
		postProp := node.GetBooleanProp(rhino.Node_INCRDECR_PROP)
		if postProp {
			g.AddExpr(first, NodeUtil_Precedence(type_), context)
			cc.AddOp(o, false)
		} else {
			cc.AddOp(o, false)
			g.self.AddNode(first)
		}
	case rhino.Token_OPTCHAIN_CALL:
		// We have two special cases here:
		// 1) If the left hand side of the call is a direct reference to eval,
		// then it must have a DIRECT_EVAL annotation. If it does not, then
		// that means it was originally an indirect call to eval, and that
		// indirectness must be preserved.
		// 2) If the left hand side of the call is a property reference,
		// then the call must not a FREE_CALL annotation. If it does, then
		// that means it was originally an call without an explicit this and
		// that must be preserved.
		if CodeGenerator_IsIndirectEval(first) ||
			(node.GetBooleanProp(rhino.Node_FREE_CALL) && NodeUtil_IsNormalOrOptChainGet(first)) {
			g.self.Add("(0,")
			g.AddExpr(first, NodeUtil_Precedence(rhino.Token_COMMA), CodeGenerator_Context_OTHER)
			g.self.Add(")")
		} else {
			g.AddExpr(first, NodeUtil_Precedence(type_), context)
		}
		args := first.GetNext()
		if node.IsOptionalChainStart() {
			g.self.Add("?.")
		}
		g.self.Add("(")
		g.self.AddList(args)
		g.self.Add(")")
	case rhino.Token_CALL:
		g.AddInvocationTarget(node, context)

		g.self.Add("(")
		g.self.AddList(first.GetNext())
		g.self.Add(")")
	case rhino.Token_IF:
		jutil.CheckStateMsg1(childCount == 2 || childCount == 3, "%s", node)
		hasElse := childCount == 3
		ambiguousElseClause := context == CodeGenerator_Context_BEFORE_DANGLING_ELSE && !hasElse
		if ambiguousElseClause {
			cc.BeginBlock()
		}

		g.self.Add("if")
		cc.MaybeInsertSpace()
		g.self.Add("(")
		g.self.AddNode(first)
		g.self.Add(")")

		if hasElse {
			g.AddNonEmptyStatement(first.GetNext(), CodeGenerator_Context_BEFORE_DANGLING_ELSE, false)
			cc.MaybeInsertSpace()
			g.self.Add("else")
			g.AddNonEmptyStatement(last, CodeGenerator_GetContextForNonEmptyExpression(context), false)
		} else {
			g.AddNonEmptyStatement(first.GetNext(), CodeGenerator_Context_OTHER, false)
		}

		if ambiguousElseClause {
			cc.EndBlock()
		}
	case rhino.Token_NULL:
		jutil.CheckStateMsg1(childCount == 0, "%s", node)
		cc.AddConstant("null")
	case rhino.Token_THIS:
		jutil.CheckStateMsg1(childCount == 0, "%s", node)
		g.self.Add("this")
	case rhino.Token_SUPER:
		jutil.CheckStateMsg1(childCount == 0, "%s", node)
		g.self.Add("super")
	case rhino.Token_NEW_TARGET:
		jutil.CheckStateMsg1(childCount == 0, "%s", node)
		g.self.Add("new.target")
	case rhino.Token_YIELD:
		g.self.Add("yield")
		if node.IsYieldAll() {
			rhino.CheckNotNull(first)
			g.self.Add("*")
		}
		if !first.IsNil() {
			cc.MaybeInsertSpace()
			g.AddExpr(first, NodeUtil_Precedence(type_), CodeGenerator_Context_OTHER)
		}
	case rhino.Token_AWAIT:
		g.self.Add("await ")
		g.AddExpr(first, NodeUtil_Precedence(type_), CodeGenerator_Context_OTHER)
	case rhino.Token_FALSE:
		jutil.CheckStateMsg1(childCount == 0, "%s", node)
		cc.AddConstant("false")
	case rhino.Token_TRUE:
		jutil.CheckStateMsg1(childCount == 0, "%s", node)
		cc.AddConstant("true")
	case rhino.Token_CONTINUE:
		jutil.CheckStateMsg1(childCount <= 1, "%s", node)
		g.self.Add("continue")
		if childCount == 1 {
			if !first.IsLabelName() {
				panic(jutil.NewJavaLangError("Unexpected token type. Should be LABEL_NAME."))
			}
			g.self.Add(" ")
			g.self.AddNode(first)
		}
		cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_DEBUGGER:
		jutil.CheckStateMsg1(childCount == 0, "%s", node)
		g.self.Add("debugger")
		cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_BREAK:
		jutil.CheckStateMsg1(childCount <= 1, "%s", node)
		g.self.Add("break")
		if childCount == 1 {
			if !first.IsLabelName() {
				panic(jutil.NewJavaLangError("Unexpected token type. Should be LABEL_NAME."))
			}
			g.self.Add(" ")
			g.self.AddNode(first)
		}
		cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_EXPR_RESULT:
		jutil.CheckStateMsg1(childCount == 1, "%s", node)
		g.self.AddNodeWithContext(first, CodeGenerator_Context_START_OF_EXPR)
		cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_NEW:
		g.self.Add("new ")
		precedence := NodeUtil_Precedence(type_)

		// `new void 0` is a syntax error add parenthese in this case. This
		// is only particularly interesting for code in dead branches.
		precedenceOfFirst := NodeUtil_Precedence(first.GetToken())
		if precedenceOfFirst == precedence {
			precedence = precedence + 1
		}

		// If the first child contains a CALL, then claim higher precedence
		// to force parentheses. Otherwise, when parsed, NEW will bind to the
		// first viable parentheses (don't traverse into functions).
		// Also, NEW requires parentheses around an optional chain callee.
		// If the first child is an arrow function, then parentheses is
		// needed
		if NodeUtil_Has(first, (rhino.Node).IsCall, NodeUtil_MATCH_NOT_FUNCTION) ||
			NodeUtil_IsOptChainNode(first) {
			precedence = NodeUtil_Precedence(first.GetToken()) + 1
		}
		g.AddExpr(first, precedence, CodeGenerator_Context_OTHER)

		// '()' is optional when no arguments are present
		next := first.GetNext()
		if !next.IsNil() {
			g.self.Add("(")
			g.self.AddList(next)
			g.self.Add(")")
		} else {
			if cc.ShouldPreserveExtras(node) {
				g.self.Add("(")
				g.self.Add(")")
			}
		}
	case rhino.Token_STRING_KEY:
		g.self.AddStringKey(node)
	case rhino.Token_STRINGLIT:
		jutil.CheckStateMsg1(childCount == 0, "String node %s may not have children", node)
		g.AddJsString(node)
	case rhino.Token_DELPROP:
		jutil.CheckStateMsg1(childCount == 1, "%s", node)
		g.self.Add("delete ")
		g.self.AddNode(first)
	case rhino.Token_OBJECTLIT:
		needsParens := context == CodeGenerator_Context_START_OF_EXPR || context.AtArrowFunctionBody()
		if needsParens {
			g.self.Add("(")
		}
		g.self.Add("{")
		for c := first; !c.IsNil(); c = c.GetNext() {
			if c != first {
				cc.ListSeparator()
			}

			jutil.CheckStateMsg1(NodeUtil_IsObjLitProperty(c) || c.IsSpread(), "%s", c)
			g.self.AddNode(c)
		}
		if !first.IsNil() && node.HasTrailingComma() {
			cc.OptionalListSeparator()
		}
		g.self.Add("}")
		if needsParens {
			g.self.Add(")")
		}
	case rhino.Token_COMPUTED_PROP:
		g.MaybeAddAccessibilityModifier(node)
		if node.GetBooleanProp(rhino.Node_STATIC_MEMBER) {
			g.self.Add("static ")
		}

		if node.GetBooleanProp(rhino.Node_COMPUTED_PROP_GETTER) {
			g.self.Add("get ")
		} else if node.GetBooleanProp(rhino.Node_COMPUTED_PROP_SETTER) {
			g.self.Add("set ")
		} else if node.GetBooleanProp(rhino.Node_COMPUTED_PROP_METHOD) {
			if last.IsAsyncFunction() {
				g.self.Add("async")
			}
			if last.GetBooleanProp(rhino.Node_GENERATOR_FN) {
				g.self.Add("*")
			}
		}
		g.self.Add("[")
		// Must use addExpr() with a priority of 1, because comma expressions
		// aren't allowed.
		// https://www.ecma-international.org/ecma-262/9.0/index.html#prod-ComputedPropertyName
		g.AddExpr(first, 1, CodeGenerator_Context_OTHER)
		g.self.Add("]")
		// TODO(martinprobst): There's currently no syntax for properties in
		// object literals that have type declarations on them (a la `{foo:
		// number: 12}`). This comes up for, e.g., function parameters with
		// default values. Support when figured out.
		g.MaybeAddTypeDecl(node)
		if node.GetBooleanProp(rhino.Node_COMPUTED_PROP_METHOD) ||
			node.GetBooleanProp(rhino.Node_COMPUTED_PROP_GETTER) ||
			node.GetBooleanProp(rhino.Node_COMPUTED_PROP_SETTER) {
			function := first.GetNext()
			params := function.GetSecondChild()
			body := function.GetLastChild()

			g.self.AddNode(params)
			g.self.AddNode(body)
		} else {
			// This is a field or object literal property.
			isInClass := node.GetParent().IsClassMembers()
			initializer := first.GetNext()
			if !initializer.IsNil() {
				// Object literal value.
				jutil.CheckStateMsg(
					!isInClass, "initializers should only exist in object literals, not classes")
				cc.Add(":")
				// Must use addExpr() with a priority of 1, because a comma
				// expression here would cause a syntax error within the
				// object literal.
				g.AddExpr(initializer, 1, CodeGenerator_Context_OTHER)
			} else {
				// Computed properties must either have an initializer or be
				// computed member-variable properties that exist for their
				// type declaration.
				jutil.CheckStateMsg1(node.GetBooleanProp(rhino.Node_COMPUTED_PROP_VARIABLE), "%s", node)
			}
		}
	case rhino.Token_OBJECT_PATTERN:
		g.self.AddObjectPattern(node)
		g.MaybeAddTypeDecl(node)
	case rhino.Token_SWITCH:
		g.self.Add("switch(")
		g.self.AddNode(first)
		g.self.Add(")")
		g.self.AddNodeWithContext(last, context)
	case rhino.Token_SWITCH_BODY:
		cc.BeginBlock()
		g.self.AddAllSiblings(first)
		cc.EndBlockWithShouldEndLine(context == CodeGenerator_Context_STATEMENT)
	case rhino.Token_CASE:
		jutil.CheckStateMsg1(childCount == 2, "%s", node)
		g.self.Add("case ")
		g.self.AddNode(first)
		g.self.AddCaseBody(last)
	case rhino.Token_DEFAULT_CASE:
		jutil.CheckStateMsg1(childCount == 1, "%s", node)
		g.self.Add("default")
		g.self.AddCaseBody(first)
	case rhino.Token_LABEL:
		jutil.CheckStateMsg1(childCount == 2, "%s", node)
		if !first.IsLabelName() {
			panic(jutil.NewJavaLangError("Unexpected token type. Should be LABEL_NAME."))
		}
		g.self.AddNode(first)
		g.self.Add(":")
		if !last.IsBlock() {
			cc.MaybeInsertSpace()
		}
		g.AddNonEmptyStatement(last, CodeGenerator_GetContextForNonEmptyExpression(context), true)
	case rhino.Token_CAST:
		if g.preserveTypeAnnotations {
			g.self.Add("(")
			g.self.AddNode(first) // drop context because of added parentheses
			g.self.Add(")")
		} else {
			g.self.AddNodeWithContext(first, context) // preserve context
		}
	case rhino.Token_TAGGED_TEMPLATELIT:
		g.AddInvocationTarget(node, context)
		g.self.AddNode(first.GetNext())
	case rhino.Token_TEMPLATELIT:
		cc.BeginTemplateLit()
		for c := first; !c.IsNil(); c = c.GetNext() {
			if c.IsTemplateLitString() {
				g.self.Add(g.EscapeUnrecognizedCharacters(c.GetRawString()))
			} else {
				cc.BeginTemplateLitSub()
				g.self.AddNodeWithContext(c.GetFirstChild(), CodeGenerator_Context_START_OF_EXPR)
				cc.EndTemplateLitSub()
			}
		}
		cc.EndTemplateLit()
	// Type Declaration ASTs.
	case rhino.Token_STRING_TYPE:
		g.self.Add("string")
	case rhino.Token_BOOLEAN_TYPE:
		g.self.Add("boolean")
	case rhino.Token_NUMBER_TYPE:
		g.self.Add("number")
	case rhino.Token_ANY_TYPE:
		g.self.Add("any")
	case rhino.Token_VOID_TYPE:
		g.self.Add("void")
	case rhino.Token_NAMED_TYPE:
		// Children are a chain of getprop nodes.
		g.self.AddNode(first)
	case rhino.Token_ARRAY_TYPE:
		g.AddExpr(first, NodeUtil_Precedence(rhino.Token_ARRAY_TYPE), context)
		g.self.Add("[]")
	case rhino.Token_FUNCTION_TYPE:
		returnType := first
		g.self.Add("(")
		g.self.AddList(first.GetNext())
		g.self.Add(")")
		cc.AddOp("=>", true)
		g.self.AddNode(returnType)
	case rhino.Token_UNION_TYPE:
		g.self.AddListWithSeparator(first, "|")
	case rhino.Token_RECORD_TYPE:
		g.self.Add("{")
		g.self.AddListWithContextAndSeparator(first, false, CodeGenerator_Context_OTHER, ",")
		g.self.Add("}")
	case rhino.Token_PARAMETERIZED_TYPE:
		// First child is the type that's parameterized, later children are
		// the arguments.
		g.self.AddNode(first)
		g.self.Add("<")
		g.self.AddList(first.GetNext())
		g.self.Add(">")
		// CLASS -> NAME,EXPR|EMPTY,BLOCK
	case rhino.Token_GENERIC_TYPE_LIST:
		g.self.Add("<")
		g.self.AddListWithContextAndSeparator(first, false, CodeGenerator_Context_STATEMENT, ",")
		g.self.Add(">")
	case rhino.Token_GENERIC_TYPE:
		g.AddIdentifier(node.GetString())
		if node.HasChildren() {
			g.self.Add("extends")
			cc.MaybeInsertSpace()
			g.self.AddNode(node.GetFirstChild())
		}
	case rhino.Token_INTERFACE:
		{
			jutil.CheckStateMsg1(childCount == 3, "%s", node)
			name := first
			superTypes := first.GetNext()
			members := last

			g.self.Add("interface")
			g.self.AddNode(name)
			g.MaybeAddGenericTypes(name)
			if !superTypes.IsEmpty() {
				g.self.Add("extends")
				superType := superTypes.GetFirstChild()
				g.self.AddNode(superType)
				for {
					superType = superType.GetNext()
					if superType.IsNil() {
						break
					}
					g.self.Add(",")
					cc.MaybeInsertSpace()
					g.self.AddNode(superType)
				}
			}
			g.self.AddNode(members)
		}
	case rhino.Token_ENUM:
		jutil.CheckStateMsg1(childCount == 2, "%s", node)
		name := first
		members := last
		g.self.Add("enum")
		g.self.AddNode(name)
		g.self.AddNode(members)
	case rhino.Token_NAMESPACE:
		jutil.CheckStateMsg1(childCount == 2, "%s", node)
		name := first
		elements := last
		g.self.Add("namespace")
		g.self.AddNode(name)
		g.self.AddNode(elements)
	case rhino.Token_TYPE_ALIAS:
		g.self.Add("type")
		g.self.Add(node.GetString())
		cc.AddOp("=", true)
		g.self.AddNode(last)
		cc.EndStatementWithSemiColon(true /* needSemiColon */, g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_DECLARE:
		g.self.Add("declare")
		g.self.AddNode(first)
		g.ProcessEnd(node, context)
	case rhino.Token_INDEX_SIGNATURE:
		g.self.Add("[")
		g.self.AddNode(first)
		g.self.Add("]")
		g.MaybeAddTypeDecl(node)
		cc.EndStatementWithSemiColon(true /* needSemiColon */, g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_CALL_SIGNATURE:
		if node.GetBooleanProp(rhino.Node_CONSTRUCT_SIGNATURE) {
			g.self.Add("new ")
		}
		g.MaybeAddGenericTypes(node)
		g.self.AddNode(first)
		g.MaybeAddTypeDecl(node)
		cc.EndStatementWithSemiColon(true /* needSemiColon */, g.HasTrailingCommentOnSameLine(node))
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unknown token " + type_.String() + "\n" + node.ToStringTree()))
	}
	if g.printNonJSDocComments {
		g.self.PrintTrailingComment(node)
	}

	cc.EndSourceMapping(node)
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

// Java: CodeGenerator.addIdentifier(String)
func (g *CodeGeneratorBase) AddIdentifier(identifier string) {
	g.cc.AddIdentifier(CodeGenerator_IdentifierEscape(identifier))
}

// Java: CodeGenerator.addGetpropIdentifier(Node)
func (g *CodeGeneratorBase) AddGetpropIdentifier(getprop rhino.Node) {
	g.cc.StartSourceMapping(getprop)
	g.AddIdentifier(getprop.GetString())
	g.cc.EndSourceMapping(getprop)
}

// Java: CodeGenerator.precedence(Node)
func (g *CodeGeneratorBase) Precedence(n rhino.Node) int {
	if n.IsCast() {
		return g.Precedence(n.GetFirstChild())
	}
	return NodeUtil_Precedence(n.GetToken())
}

// AddInvocationTarget handles two special cases:
//
//   - If the left hand side of the call is a direct reference to eval, then
//     it must have a DIRECT_EVAL annotation. If it does not, then that means
//     it was originally an indirect call to eval, and that indirectness must
//     be preserved.
//   - If the left hand side of the call is a property reference, then the
//     call must not a FREE_CALL annotation. If it does, then that means it
//     was originally an call without an explicit this and that must be
//     preserved.
//
// Java: CodeGenerator.addInvocationTarget(Node, Context)
func (g *CodeGeneratorBase) AddInvocationTarget(node rhino.Node, context CodeGenerator_Context) {
	first := node.GetFirstChild()

	needsParens := NodeUtil_IsOptChainNode(first)
	if CodeGenerator_IsIndirectEval(first) ||
		(node.GetBooleanProp(rhino.Node_FREE_CALL) && NodeUtil_IsNormalOrOptChainGet(first)) {
		g.self.Add("(0,")
		g.AddExpr(first, NodeUtil_Precedence(rhino.Token_COMMA), CodeGenerator_Context_OTHER)
		g.self.Add(")")
	} else {
		if needsParens {
			g.self.Add("(")
		}
		g.AddExpr(first, NodeUtil_Precedence(node.GetToken()), context)
		if needsParens {
			g.self.Add(")")
		}
	}
}

// Java: CodeGenerator.arrowFunctionNeedsParens(Node)
func (g *CodeGeneratorBase) ArrowFunctionNeedsParens(n rhino.Node) bool {
	parent := n.GetParent()
	expressionOrEnclosingCast := n
	for !parent.IsNil() && parent.IsCast() {
		if g.preserveTypeAnnotations {
			// If printing type annotations, any expression in a CAST
			// automatically is wrapped in parentheses when printing the
			// CAST. Returning true here would add a second, unnecessary pair
			// of parentheses.
			return false
		}
		// If not printing type annotations, then pretend the CAST node is
		// not there and check the parent of the CAST.
		expressionOrEnclosingCast = parent
		parent = parent.GetParent()
	}

	// Once you cut through the layers of non-terminals used to define
	// operator precedence, you can see the following are true.
	// (Read => as "may expand to" and "!=>" as "may not expand to")
	//
	// 1. You can substitute an ArrowFunction into rules where an Expression
	//    or AssignmentExpression is required, because
	//      Expression => AssignmentExpression => ArrowFunction
	//
	// 2. However, most operators act on LeftHandSideExpression,
	//    CallExpression, or MemberExpression. None of these expand to
	//    ArrowFunction.
	//
	// 3. CallExpression cannot expand to an ArrowFunction at all, because
	//    all of its expansions produce multiple symbols and none can be
	//    logically equivalent to ArrowFunction.
	//
	// 4. LeftHandSideExpression and MemberExpression may be replaced with an
	//    ArrowFunction in parentheses, because:
	//      LeftHandSideExpression => MemberExpression => PrimaryExpression
	//      PrimaryExpression => '(' Expression ')' => '(' ArrowFunction ')'
	if parent.IsNil() {
		return false
	} else if NodeUtil_IsBinaryOperator(parent) ||
		NodeUtil_IsUnaryOperator(parent) ||
		NodeUtil_IsUpdateOperator(parent) ||
		parent.IsTaggedTemplateLit() ||
		parent.IsGetProp() ||
		parent.IsOptChainGetProp() ||
		parent.IsAwait() ||
		parent.IsYield() {
		// LeftHandSideExpression OP LeftHandSideExpression
		// OP LeftHandSideExpression | LeftHandSideExpression OP
		// MemberExpression TemplateLiteral
		// MemberExpression '.' IdentifierName
		return true
	} else if parent.IsGetElem() ||
		parent.IsCall() ||
		parent.IsHook() ||
		parent.IsOptChainGetElem() ||
		parent.IsOptChainCall() ||
		parent.IsNew() {
		// MemberExpression '[' Expression ']'
		// MemberFunction '(' AssignmentExpressionList ')'
		// LeftHandSideExpression ? AssignmentExpression : AssignmentExpression
		return expressionOrEnclosingCast.IsFirstChildOf(parent)
	} else {
		// All other cases are either illegal (e.g. because you cannot assign
		// a value to an ArrowFunction) or do not require parens.
		return false
	}
}

// Java: CodeGenerator.addArrowFunction(Node, Node, Node, Context)
func (g *CodeGeneratorBase) AddArrowFunction(n rhino.Node, first rhino.Node, last rhino.Node, context CodeGenerator_Context) {
	jutil.CheckStateMsg1(first.GetString() == "", "%s", first)
	funcNeedsParens := g.ArrowFunctionNeedsParens(n)
	if funcNeedsParens {
		g.self.Add("(")
	}

	g.MaybeAddGenericTypes(first)

	if n.IsAsyncFunction() {
		g.self.Add("async")
	}
	g.self.AddNode(first.GetNext()) // param list
	g.MaybeAddTypeDecl(n)

	g.cc.AddOp("=>", true)

	if last.IsBlock() {
		g.self.AddNode(last)
	} else {
		// This is a hack. Arrow functions have no token type, but
		// blockless arrow function bodies have lower precedence than
		// anything other than commas.
		g.AddExpr(last, NodeUtil_Precedence(rhino.Token_COMMA)+1, CodeGenerator_GetContextForArrowFunctionBody(context))
	}
	g.cc.EndFunction(context == CodeGenerator_Context_STATEMENT)

	if funcNeedsParens {
		g.self.Add(")")
	}
}

// Java: CodeGenerator.addFunction(Node, Node, Node, Context)
func (g *CodeGeneratorBase) AddFunction(n rhino.Node, first rhino.Node, last rhino.Node, context CodeGenerator_Context) {
	funcNeedsParens := (context == CodeGenerator_Context_START_OF_EXPR)
	if funcNeedsParens {
		g.self.Add("(")
	}

	if n.IsAsyncFunction() {
		g.self.Add("async function")
	} else {
		g.self.Add("function")
	}
	if n.IsGeneratorFunction() {
		g.self.Add("*")
		if first.GetString() != "" {
			g.cc.MaybeInsertSpace()
		}
	}

	g.self.AddNode(first)
	g.MaybeAddGenericTypes(first)

	g.self.AddNode(first.GetNext()) // param list
	g.MaybeAddTypeDecl(n)

	g.self.AddNode(last)
	g.cc.EndFunction(context == CodeGenerator_Context_STATEMENT)

	if funcNeedsParens {
		g.self.Add(")")
	}
}

// Java: CodeGenerator.maybeAddAccessibilityModifier(Node)
func (g *CodeGeneratorBase) MaybeAddAccessibilityModifier(n rhino.Node) {
	access := n.GetProp(rhino.Node_ACCESS_MODIFIER)
	if access != nil {
		g.self.Add(jutil.StringToLowerCaseASCII(access.(rhino.JSDocInfo_Visibility).String()) + " ")
	}
}

// Java: CodeGenerator.maybeAddTypeDecl(Node)
func (g *CodeGeneratorBase) MaybeAddTypeDecl(n rhino.Node) {
	if !n.GetDeclaredTypeExpression().IsNil() {
		g.self.Add(":")
		g.cc.MaybeInsertSpace()
		g.self.AddNode(n.GetDeclaredTypeExpression())
	}
}

// Java: CodeGenerator.maybeAddGenericTypes(Node)
func (g *CodeGeneratorBase) MaybeAddGenericTypes(n rhino.Node) {
	generics := codeGenerator_nodeProp(n.GetProp(rhino.Node_GENERIC_TYPE_LIST))
	if !generics.IsNil() {
		g.self.AddNode(generics)
	}
}

// Java: CodeGenerator.maybeAddOptional(Node)
func (g *CodeGeneratorBase) MaybeAddOptional(n rhino.Node) {
	if n.GetBooleanProp(rhino.Node_OPT_ES6_TYPED) {
		g.self.Add("?")
	}
}

// UnrollBinaryOperator: we could use addList recursively here, but
// sometimes we produce very deeply nested operators and run out of stack
// space, so we just unroll the recursion when possible.
//
// We assume nodes are left-recursive.
//
// Java: CodeGenerator.unrollBinaryOperator(Node, Token, String, Context, Context, int, int)
func (g *CodeGeneratorBase) UnrollBinaryOperator(
	n rhino.Node,
	op rhino.Token,
	opStr string,
	context CodeGenerator_Context,
	rhsContext CodeGenerator_Context,
	leftPrecedence int,
	rightPrecedence int) {
	firstNonOperator := n.GetFirstChild()
	for firstNonOperator.GetToken() == op {
		firstNonOperator = firstNonOperator.GetFirstChild()
	}

	g.AddExpr(firstNonOperator, leftPrecedence, context)

	current := firstNonOperator
	for {
		current = current.GetParent()
		g.cc.AddOp(opStr, true)
		g.AddExpr(current.GetSecondChild(), rightPrecedence, rhsContext)
		if current == n {
			break
		}
	}
}

// Java: CodeGenerator.isSimpleNumber(String)
func CodeGenerator_IsSimpleNumber(s string) bool {
	// Go-only: language-adapter — every accepted char is an ASCII digit, so
	// the byte length equals the UTF-16 length whenever the loop completes.
	len_ := len(s)
	if len_ == 0 {
		return false
	}
	for index := 0; index < len_; index++ {
		c := s[index]
		if c < '0' || c > '9' {
			return false
		}
	}
	return len_ == 1 || s[0] != '0'
}

// Java: CodeGenerator.getSimpleNumber(String)
func CodeGenerator_GetSimpleNumber(s string) float64 {
	if CodeGenerator_IsSimpleNumber(s) {
		// Long.parseLong throws NumberFormatException when the number is
		// too long to parse; that falls through to NaN.
		l, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			if l <= NodeUtil_MAX_POSITIVE_INTEGER_NUMBER {
				return float64(l)
			}
		}
	}
	return math.NaN()
}

// CodeGenerator_IsIndirectEval returns whether the name is an indirect
// eval.
//
// Java: CodeGenerator.isIndirectEval(Node)
func CodeGenerator_IsIndirectEval(n rhino.Node) bool {
	return n.IsName() && n.GetString() == "eval" && !n.GetBooleanProp(rhino.Node_DIRECT_EVAL)
}

// AddNonEmptyStatement adds a block or expression, substituting a VOID with
// an empty statement. This is used for "for (...);" and "if (...);" type
// statements.
//
// n is the node to print; context is the context to determine how the node
// should be printed.
//
// Java: CodeGenerator.addNonEmptyStatement(Node, Context, boolean)
func (g *CodeGeneratorBase) AddNonEmptyStatement(n rhino.Node, context CodeGenerator_Context, allowNonBlockChild bool) {
	nodeToProcess := n

	if !allowNonBlockChild && !n.IsBlock() {
		panic(jutil.NewJavaLangError("Missing BLOCK child."))
	}

	// Strip unneeded blocks, that is blocks with <2 children unless
	// the CodePrinter specifically wants to keep them.
	if n.IsBlock() {
		count := CodeGenerator_GetNonEmptyChildCount(n, 2)
		if count == 0 {
			if g.cc.ShouldPreserveExtras(n) {
				g.cc.BeginBlock()
				g.self.PrintTrailingComment(n)
				g.cc.EndBlockWithShouldEndLine(g.cc.BreakAfterBlockFor(n, context == CodeGenerator_Context_STATEMENT))
			} else {
				g.self.PrintTrailingComment(n)
				g.cc.EndStatementWithSemiColon(true /* needSemiColon */, false /* hasTrailingCommentOnSameLine */)
			}
			return
		}

		if count == 1 {
			// Preserve the block only if needed or requested.
			// 'let', 'const', etc are not allowed by themselves in "if" and
			// other structures. Also, hack around a IE6/7 browser bug that
			// needs a block around DOs.
			firstAndOnlyChild := CodeGenerator_GetFirstNonEmptyChild(n)
			alwaysWrapInBlock := g.cc.ShouldPreserveExtras(n)
			if alwaysWrapInBlock || CodeGenerator_IsBlockDeclOrDo(firstAndOnlyChild) {
				g.cc.BeginBlock()
				g.self.AddNodeWithContext(firstAndOnlyChild, CodeGenerator_Context_STATEMENT)
				g.self.PrintTrailingComment(n)
				g.cc.MaybeLineBreak()
				g.cc.EndBlockWithShouldEndLine(g.cc.BreakAfterBlockFor(n, context == CodeGenerator_Context_STATEMENT))
				return
			} else {
				// Continue with the only child.
				nodeToProcess = firstAndOnlyChild
			}
		}
	}

	if nodeToProcess.IsEmpty() {
		g.self.PrintTrailingComment(n)
		g.cc.EndStatementWithSemiColon(true /* needSemiColon */, false /* hasTrailingCommentOnSameLine */)
	} else {
		g.self.AddNodeWithContext(nodeToProcess, context)
		g.self.PrintTrailingComment(n)
	}
}

// CodeGenerator_IsBlockDeclOrDo returns whether the Node is a DO or a
// declaration that is only allowed in restricted contexts.
//
// Java: CodeGenerator.isBlockDeclOrDo(Node)
func CodeGenerator_IsBlockDeclOrDo(n rhino.Node) bool {
	if n.IsLabel() {
		labeledStatement := n.GetLastChild()
		if !labeledStatement.IsBlock() {
			return CodeGenerator_IsBlockDeclOrDo(labeledStatement)
		} else {
			// For labels with block children, we need to ensure that a
			// labeled FUNCTION or DO isn't generated when extraneous BLOCKs
			// are skipped.
			if CodeGenerator_GetNonEmptyChildCount(n, 2) == 1 {
				return CodeGenerator_IsBlockDeclOrDo(CodeGenerator_GetFirstNonEmptyChild(n))
			} else {
				// Either a empty statement or an block with more than one
				// child, way it isn't a FUNCTION or DO.
				return false
			}
		}
	} else {
		switch n.GetToken() {
		case rhino.Token_LET, rhino.Token_CONST, rhino.Token_FUNCTION, rhino.Token_CLASS, rhino.Token_DO:
			return true
		default:
			return false
		}
	}
}

// Java: CodeGenerator.addExpr(Node, int, Context)
func (g *CodeGeneratorBase) AddExpr(n rhino.Node, minPrecedence int, context CodeGenerator_Context) {
	if g.OpRequiresParentheses(n, minPrecedence, context) {
		g.self.Add("(")
		g.self.AddNodeWithContext(n, CodeGenerator_Context_OTHER)
		g.self.Add(")")
	} else {
		g.self.AddNodeWithContext(n, context)
	}
}

// Java: CodeGenerator.opRequiresParentheses(Node, int, Context)
func (g *CodeGeneratorBase) OpRequiresParentheses(n rhino.Node, minPrecedence int, context CodeGenerator_Context) bool {
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
		// Class declarations with assignments should be wrapped in
		// parentheses.
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
func (g *CodeGeneratorBase) IsFirstOperandOfExponentiationExpression(n rhino.Node) bool {
	parent := n.GetParent()
	return !parent.IsNil() && parent.IsExponent() && parent.GetFirstChild() == n
}

// Java: CodeGenerator.addList(Node)
func (g *CodeGeneratorBase) AddList(firstInList rhino.Node) {
	g.self.AddListWithContextAndSeparator(firstInList, true, CodeGenerator_Context_OTHER, ",")
}

// Java: CodeGenerator.addList(Node, String)
func (g *CodeGeneratorBase) AddListWithSeparator(firstInList rhino.Node, separator string) {
	g.self.AddListWithContextAndSeparator(firstInList, true, CodeGenerator_Context_OTHER, separator)
}

// Java: CodeGenerator.addList(Node, boolean, Context, String)
func (g *CodeGeneratorBase) AddListWithContextAndSeparator(
	firstInList rhino.Node, isArrayOrFunctionArgument bool, lhsContext CodeGenerator_Context, separator string) {
	if firstInList.IsNil() {
		return
	}
	for n := firstInList; !n.IsNil(); n = n.GetNext() {
		isFirst := n == firstInList
		minPrecedence := 0
		if isArrayOrFunctionArgument {
			minPrecedence = 1
		}
		if isFirst {
			g.AddExpr(n, minPrecedence, lhsContext)
		} else {
			g.cc.AddOp(separator, true)
			g.AddExpr(n, minPrecedence, CodeGenerator_GetContextForNoInOperator(lhsContext))
		}
	}
	if isArrayOrFunctionArgument && rhino.CheckNotNull(firstInList.GetParent()).HasTrailingComma() {
		g.cc.OptionalListSeparator()
	}
}

// Java: CodeGenerator.addStringKey(Node)
func (g *CodeGeneratorBase) AddStringKey(n rhino.Node) {
	key := n.GetString()
	// Object literal property names don't have to be quoted if they are not
	// JavaScript keywords.
	mustBeQuoted := n.IsQuotedStringKey() ||
		(g.quoteKeywordProperties && rhino.TokenStream_IsKeyword(key)) ||
		!rhino.TokenStream_IsJSIdentifier(key) ||
		// do not encode literally any non-literal characters that were
		// Unicode escaped.
		!NodeUtil_IsLatin(key)
	if !mustBeQuoted {
		// Check if the property is eligible to be printed as shorthand.
		if n.IsShorthandProperty() {
			child := n.GetFirstChild()
			if child.MatchesQualifiedName(key) ||
				(child.IsDefaultValue() && child.GetFirstChild().MatchesQualifiedName(key)) {
				g.self.AddNode(child)
				return
			}
		}
		g.self.Add(key)
	} else {
		// Determine if the string is a simple number.
		d := CodeGenerator_GetSimpleNumber(key)
		if !math.IsNaN(d) {
			g.cc.AddNumber(d, n)
		} else {
			g.AddJsString(n)
		}
	}
	if n.HasChildren() {
		// NOTE: the only time a STRING_KEY node does *not* have children is
		// when it's inside a TypeScript enum. We should change these to
		// their own ENUM_KEY token so that the bifurcating logic can be
		// removed from STRING_KEY.
		g.self.Add(":")
		g.AddExpr(n.GetFirstChild(), 1, CodeGenerator_Context_OTHER)
	}
}

// Java: CodeGenerator.addObjectPattern(Node)
func (g *CodeGeneratorBase) AddObjectPattern(n rhino.Node) {
	g.self.Add("{")
	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		if child != n.GetFirstChild() {
			g.cc.ListSeparator()
		}

		g.self.AddNode(child)
	}
	g.self.Add("}")
}

// AddArrayList adds a comma-separated list as is specified by an ARRAYLIT
// node.
//
// firstInList is the first in the node list (chained through the next
// property).
//
// Java: CodeGenerator.addArrayList(Node)
func (g *CodeGeneratorBase) AddArrayList(firstInList rhino.Node) {
	if firstInList.IsNil() {
		return
	}
	lastWasEmpty := false
	for n := firstInList; !n.IsNil(); n = n.GetNext() {
		if n != firstInList {
			g.cc.ListSeparator()
		}
		g.AddExpr(n, 1, CodeGenerator_Context_OTHER)
		lastWasEmpty = n.IsEmpty()
	}

	if lastWasEmpty {
		g.cc.ListSeparator()
	} else if firstInList.GetParent().HasTrailingComma() {
		g.cc.OptionalListSeparator()
	}
}

// Java: CodeGenerator.addCaseBody(Node)
func (g *CodeGeneratorBase) AddCaseBody(caseBody rhino.Node) {
	jutil.CheckStateMsg1(caseBody.IsBlock(), "%s", caseBody)
	g.cc.BeginCaseBody()
	g.self.AddAllSiblings(caseBody.GetFirstChild())
	g.cc.EndCaseBody()
}

// Java: CodeGenerator.addAllSiblings(Node)
func (g *CodeGeneratorBase) AddAllSiblings(n rhino.Node) {
	for c := n; !c.IsNil(); c = c.GetNext() {
		g.self.AddNode(c)
	}
}

// Java: CodeGenerator.addNonJsDoc_nonTrailing(Node, NonJSDocComment)
func (g *CodeGeneratorBase) AddNonJsDoc_nonTrailing(node rhino.Node, nonJSDocComment *rhino.NonJSDocComment) {
	content := nonJSDocComment.GetCommentString()
	commentEndPosition := nonJSDocComment.GetEndPosition()

	nodeLineNumber := node.GetLineno() - 1 // source lines are 1-indexed

	if nonJSDocComment.IsEndingAsLineComment() {
		// Non trailing line comments can not be on the same line as the node.
		jutil.CheckStateMsg(
			g.gentsMode || int(commentEndPosition.Line) < nodeLineNumber,
			"Non trailing line comments can not be on the same line as the node.")
		g.self.Add(content + "\n")
	} else {
		if nodeLineNumber == int(commentEndPosition.Line) {
			// e.g. ``` /* comment */ let x; ```
			g.self.Add(content + " ")
		} else {
			// e.g.
			// ```
			// /* comment */
			// let x;
			// ```
			g.self.Add(content + "\n")
		}
	}
}

// Java: CodeGenerator.addNonJsDoctrailing(NonJSDocComment, boolean)
func (g *CodeGeneratorBase) AddNonJsDoctrailing(nonJSDocComment *rhino.NonJSDocComment, sameLine bool) {
	content := nonJSDocComment.GetCommentString()
	if nonJSDocComment.IsEndingAsLineComment() {
		// Trailing line comments *must* end with a `\n`. E.g.. `let x;
		// //comment\n`
		g.self.Add(" " + content)
		if sameLine {
			g.cc.StartNewLine()
		}
	} else {
		if nonJSDocComment.IsInline() {
			// e.g. `foo(x /*comment*/);` is inline
			g.self.Add(" " + content)
		} else {
			// e.g. `let x; /*comment*/` is non-inline
			g.self.Add(" " + content)
			if sameLine {
				g.cc.StartNewLine()
			}
		}
	}
}

// AddJsString outputs a JS string, using the optimal (single/double) quote
// character.
//
// Java: CodeGenerator.addJsString(Node)
func (g *CodeGeneratorBase) AddJsString(n rhino.Node) {
	g.self.Add(g.JsString(n.GetString()))
}

var codeGenerator_plainASCII = func() [256]bool {
	var t [256]bool
	for c := 0x20; c <= 0x7e; c++ {
		switch byte(c) {
		case '\\', '"', '\'', '$', '`', '=', '&', '>', '<':
		default:
			t[c] = true
		}
	}
	return t
}()

const codeGenerator_HEX_CHARS = "0123456789abcdef"

func codeGenerator_appendHexJS(sb *strings.Builder, c uint16) {
	var buf [6]byte
	buf[0] = '\\'
	buf[1] = 'u'
	buf[2] = codeGenerator_HEX_CHARS[(c>>12)&0xf]
	buf[3] = codeGenerator_HEX_CHARS[(c>>8)&0xf]
	buf[4] = codeGenerator_HEX_CHARS[(c>>4)&0xf]
	buf[5] = codeGenerator_HEX_CHARS[c&0xf]
	sb.Write(buf[:])
}

// Java: CodeGenerator.jsString(String)
func (g *CodeGeneratorBase) JsString(s string) string {
	singleq := 0
	doubleq := 0
	plain := true

	// could count the quotes and pick the optimal quote character
	//
	// Go-only: language-adapter — both quotes are ASCII, and no byte of a
	// multi-byte WTF-8 sequence is ASCII, so counting bytes counts the same
	// chars as Java's UTF-16 loop.
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !codeGenerator_plainASCII[c] {
			plain = false
			switch c {
			case '"':
				doubleq++
			case '\'':
				singleq++
			}
		}
	}

	if plain {
		quote := byte('"')
		if g.preferSingleQuotes {
			quote = '\''
		}
		var sb strings.Builder
		sb.Grow(len(s) + 2)
		sb.WriteByte(quote)
		sb.WriteString(s)
		sb.WriteByte(quote)
		return sb.String()
	}

	var doublequote string
	var singlequote string
	var quote byte
	useSingle := singleq < doubleq
	if g.preferSingleQuotes {
		useSingle = singleq <= doubleq
	}
	if useSingle {
		// more double quotes so enclose in single quotes.
		quote = '\''
		doublequote = "\""
		singlequote = "\\'"
	} else {
		// more single quotes so escape the doubles
		quote = '"'
		doublequote = "\\\""
		singlequote = "'"
	}

	var sb strings.Builder
	sb.Grow(len(s) + 4)
	sb.WriteByte(quote)
	g.strEscapeInto(&sb, s, doublequote, singlequote, "`", "\\\\", "$", false)
	sb.WriteByte(quote)
	return sb.String()
}

// RegexpEscape escapes regular expression.
//
// Java: CodeGenerator.regexpEscape(String)
func (g *CodeGeneratorBase) RegexpEscape(s string) string {
	var sb strings.Builder
	sb.Grow(len(s) + 2)
	sb.WriteByte('/')
	g.strEscapeInto(&sb, s, "\"", "'", "`", "\\", "$", true)
	sb.WriteByte('/')
	return sb.String()
}

// codeGenerator_appendChar is sb.append(char) for a char that is not
// escaped.
//
// Go-only: language-adapter — appends one UTF-16 char to WTF-8 output.
func codeGenerator_appendChar(sb *strings.Builder, c uint16) {
	if c < utf8.RuneSelf {
		sb.WriteByte(byte(c))
	} else if c >= 0xD800 && c <= 0xDFFF {
		// A surrogate half on its own: kept as WTF-8.
		sb.WriteString(jutil.StringFromChar(c))
	} else {
		sb.WriteRune(rune(c))
	}
}

// codeGenerator_regionMatches is String.regionMatches(boolean, int, String,
// int, int) on the UTF-16 units of the string being escaped; other is
// ASCII.
//
// Go-only: language-adapter — Java's regionMatches on the UTF-16 view.
func codeGenerator_regionMatches(units []uint16, ignoreCase bool, toffset int, other string) bool {
	return jutil.UTF16RegionMatches(ignoreCase, units, toffset, jutil.StringToUTF16(other), 0, len(other))
}

// StrEscape is a helper to escape JavaScript string as well as regular
// expression.
//
// Java: CodeGenerator.strEscape(String, String, String, String, String, String, boolean)
func (g *CodeGeneratorBase) StrEscape(
	s string,
	doublequoteEscape string,
	singlequoteEscape string,
	backtickEscape string,
	backslashEscape string,
	dollarEscape string,
	isRegexp bool) string {
	plain := true
	for i := 0; i < len(s); i++ {
		if !codeGenerator_plainASCII[s[i]] {
			plain = false
			break
		}
	}
	if plain {
		return s
	}
	var sb strings.Builder
	sb.Grow(len(s) + 2)
	g.strEscapeInto(&sb, s, doublequoteEscape, singlequoteEscape, backtickEscape, backslashEscape, dollarEscape, isRegexp)
	return sb.String()
}

func (g *CodeGeneratorBase) strEscapeInto(
	sb *strings.Builder,
	s string,
	doublequoteEscape string,
	singlequoteEscape string,
	backtickEscape string,
	backslashEscape string,
	dollarEscape string,
	isRegexp bool) {
	// Go-only: language-adapter — Java indexes the UTF-16 chars of s.
	units := g.utf16Units(s)
	for i := 0; i < len(units); i++ {
		c := units[i]
		switch c {
		case 0:
			sb.WriteString("\\x00")
		case '\u000B':
			if !isRegexp {
				sb.WriteString("\\v")
			} else {
				sb.WriteString("\\x0B")
			}
			// From the SingleEscapeCharacter grammar production.
		case '\b':
			sb.WriteString("\\b")
		case '\f':
			sb.WriteString("\\f")
		case '\n':
			sb.WriteString("\\n")
		case '\r':
			sb.WriteString("\\r")
		case '\t':
			sb.WriteString("\\t")
		case '\\':
			sb.WriteString(backslashEscape)
		case '"':
			sb.WriteString(doublequoteEscape)
		case '\'':
			sb.WriteString(singlequoteEscape)
		case '$':
			sb.WriteString(dollarEscape)
		case '`':
			sb.WriteString(backtickEscape)
		case '=':
			g.AppendEscapedEq(sb, isRegexp)
		case '&':
			g.AppendEscapedAmp(sb, isRegexp)
		case '>':
			g.AppendEscapedGt(sb, units, i, isRegexp)
		case '<':
			g.AppendEscapedLt(sb, units, i, isRegexp)
		default:
			if isRegexp ||
				!g.outputFeatureSet.ContainsFeature(parser.FeatureSet_Feature_UNESCAPED_UNICODE_LINE_OR_PARAGRAPH_SEP) {
				// In 2019 these characters (line and paragraph separators)
				// are valid in strings but not regular expressions.
				// https://github.com/tc39/proposal-json-superset
				if c == ' ' {
					sb.WriteString("\\u2028")
					break
				}
				if c == ' ' {
					sb.WriteString("\\u2029")
					break
				}
			}

			if (g.outputCharsetEncoder != nil && g.outputCharsetEncoder.CanEncode(c)) ||
				(c > 0x1f && c < 0x7f) {
				// If we're given an outputCharsetEncoder, then check if the
				// character can be represented in this character set. If no
				// charsetEncoder provided - pass straight Latin characters
				// through, and escape the rest. Doing the explicit character
				// check is measurably faster than using the CharsetEncoder.
				codeGenerator_appendChar(sb, c)
			} else {
				// Other characters can be misinterpreted by some JS parsers,
				// or perhaps mangled by proxies along the way,
				// so we play it safe and Unicode escape them.
				codeGenerator_appendHexJS(sb, c)
			}
		}
	}
}

// Java: CodeGenerator.appendEscapedEq(StringBuilder, boolean)
func (g *CodeGeneratorBase) AppendEscapedEq(sb *strings.Builder, isRegexp bool) {
	// '=' is a syntactically significant regexp character.
	if g.trustedStrings || isRegexp {
		sb.WriteByte('=')
	} else {
		sb.WriteString("\\x3d")
	}
}

// Java: CodeGenerator.appendEscapedAmp(StringBuilder, boolean)
func (g *CodeGeneratorBase) AppendEscapedAmp(sb *strings.Builder, isRegexp bool) {
	if g.trustedStrings || isRegexp {
		sb.WriteByte('&')
	} else {
		sb.WriteString("\\x26")
	}
}

// AppendEscapedGt: s is the string being escaped, as its UTF-16 units.
//
// Java: CodeGenerator.appendEscapedGt(StringBuilder, String, int, boolean)
func (g *CodeGeneratorBase) AppendEscapedGt(sb *strings.Builder, s []uint16, index int, isRegexp bool) {
	if !g.trustedStrings && !isRegexp {
		sb.WriteString(codeGenerator_GT_ESCAPED)
		return
	}

	// Break --> into --\> or ]]> into ]]\>
	//
	// This is just to prevent developers from shooting themselves in the
	// foot, and does not provide the level of security that you get with
	// trustedString == false.
	if index >= 2 &&
		((s[index-1] == '-' && s[index-2] == '-') ||
			(s[index-1] == ']' && s[index-2] == ']')) {
		sb.WriteString(codeGenerator_GT_ESCAPED)
	} else {
		sb.WriteByte('>')
	}
}

// AppendEscapedLt: s is the string being escaped, as its UTF-16 units.
//
// Java: CodeGenerator.appendEscapedLt(StringBuilder, String, int, boolean)
func (g *CodeGeneratorBase) AppendEscapedLt(sb *strings.Builder, s []uint16, index int, isRegexp bool) {
	if !g.trustedStrings && !isRegexp {
		sb.WriteString(codeGenerator_LT_ESCAPED)
		return
	}

	// Break </script into <\/script and </style into <\/style
	// As above, this is just to prevent developers from doing this accidentally.
	const endScript = "/script"
	const endStyle = "/style"

	// Break <!-- into <\!--
	const startComment = "!--"

	if codeGenerator_regionMatches(s, true, index+1, endScript) ||
		codeGenerator_regionMatches(s, true, index+1, endStyle) ||
		codeGenerator_regionMatches(s, false, index+1, startComment) {
		sb.WriteString(codeGenerator_LT_ESCAPED)
	} else {
		sb.WriteByte('<')
	}
}

// EscapeUnrecognizedCharacters is a helper to escape the characters that
// might be misinterpreted.
//
// s is the string to modify; returns the string with unrecognizable
// characters escaped.
//
// Java: CodeGenerator.escapeUnrecognizedCharacters(String)
func (g *CodeGeneratorBase) EscapeUnrecognizedCharacters(s string) string {
	// TODO(yitingwang) Move this method to a suitable place
	plain := true
	for i := 0; i < len(s); i++ {
		if !codeGenerator_plainASCII[s[i]] {
			plain = false
			break
		}
	}
	if plain {
		return s
	}
	// Go-only: language-adapter — Java indexes the UTF-16 chars of s.
	units := g.utf16Units(s)
	var sb strings.Builder
	sb.Grow(len(s))
	for i := 0; i < len(units); i++ {
		c := units[i]
		switch c {
		case '\b', '\f', '\n', '\r', '\t', '\\', '"', '\'', '$', '`', ' ', ' ':
			// From the SingleEscapeCharacter grammar production.
			codeGenerator_appendChar(&sb, c)
		case '=':
			g.AppendEscapedEq(&sb /* isRegexp= */, false)
		case '&':
			g.AppendEscapedAmp(&sb /* isRegexp= */, false)
		case '>':
			g.AppendEscapedGt(&sb, units, i /* isRegexp= */, false)
		case '<':
			g.AppendEscapedLt(&sb, units, i /* isRegexp= */, false)
		default:
			if (g.outputCharsetEncoder != nil && g.outputCharsetEncoder.CanEncode(c)) ||
				(c > 0x1f && c < 0x7f) {
				// If we're given an outputCharsetEncoder, then check if the
				// character can be represented in this character set. If no
				// charsetEncoder provided - pass straight Latin characters
				// through, and escape the rest. Doing the explicit character
				// check is measurably faster than using the CharsetEncoder.
				codeGenerator_appendChar(&sb, c)
			} else {
				// Other characters can be misinterpreted by some JS parsers,
				// or perhaps mangled by proxies along the way,
				// so we play it safe and Unicode escape them.
				codeGenerator_appendHexJS(&sb, c)
			}
		}
	}
	return sb.String()
}

// Java: CodeGenerator.identifierEscape(String)
func CodeGenerator_IdentifierEscape(s string) string {
	// First check if escaping is needed at all -- in most cases it isn't.
	if NodeUtil_IsLatin(s) {
		return s
	}

	// Now going through the string to escape non-Latin characters if needed.
	//
	// Go-only: language-adapter — Java indexes the UTF-16 chars of s.
	units := jutil.StringToUTF16(s)
	var sb strings.Builder
	sb.Grow(len(s) + 6)
	for i := 0; i < len(units); i++ {
		c := units[i]
		// Identifiers should always go to Latin1/ ASCII characters because
		// different browser's rules for valid identifier characters are
		// crazy.
		if c > 0x1F && c < 0x7F {
			sb.WriteByte(byte(c))
		} else {
			codeGenerator_appendHexJS(&sb, c)
		}
	}
	return sb.String()
}

// CodeGenerator_GetNonEmptyChildCount returns the number of children of
// this node that are non empty up to maxCount (the maximum number of
// children to look for).
//
// Java: CodeGenerator.getNonEmptyChildCount(Node, int)
func CodeGenerator_GetNonEmptyChildCount(n rhino.Node, maxCount int) int {
	i := 0
	c := n.GetFirstChild()
	for ; !c.IsNil() && i < maxCount; c = c.GetNext() {
		if c.IsBlock() {
			i += CodeGenerator_GetNonEmptyChildCount(c, maxCount-i)
		} else if !c.IsEmpty() {
			i++
		}
	}
	return i
}

// CodeGenerator_GetFirstNonEmptyChild gets the first non-empty child of the
// given node.
//
// Java: CodeGenerator.getFirstNonEmptyChild(Node)
func CodeGenerator_GetFirstNonEmptyChild(n rhino.Node) rhino.Node {
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		if c.IsBlock() {
			result := CodeGenerator_GetFirstNonEmptyChild(c)
			if !result.IsNil() {
				return result
			}
		} else if !c.IsEmpty() {
			return c
		}
	}
	return rhino.NilNode
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
	// body. This is only important when the first token after the "=>" is
	// "{".
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

// The following two cases are independent, unlike the other enum states, so
// we have separate booleans for them.
//
// Context.inForInitClause, by ordinal.
var codeGenerator_Context_inForInitClause = [...]bool{
	false, // STATEMENT
	false, // BEFORE_DANGLING_ELSE
	false, // START_OF_EXPR
	/* inForInitClause= */ true, // IN_FOR_INIT_CLAUSE
	/* inForInitClause= */ false, // START_OF_ARROW_FN_BODY
	/* inForInitClause= */ true, // START_OF_ARROW_FN_IN_FOR_INIT
	false,                       // OTHER
}

// Context.atArrowFnBody, by ordinal.
var codeGenerator_Context_atArrowFnBody = [...]bool{
	false, // STATEMENT
	false, // BEFORE_DANGLING_ELSE
	false, // START_OF_EXPR
	/* at start of arrow fn */ false, // IN_FOR_INIT_CLAUSE
	/* at start of arrow fn */ true, // START_OF_ARROW_FN_BODY
	/* atArrowFunctionBody */ true, // START_OF_ARROW_FN_IN_FOR_INIT
	false,                          // OTHER
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

// CodeGenerator_GetContextForNoInOperator: if we're in a
// IN_FOR_INIT_CLAUSE, we can't permit in operators in the expression. Pass
// on the IN_FOR_INIT_CLAUSE flag through subexpressions.
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
// object patterns. We also must also pass the IN_FOR_INIT_CLAUSE flag into
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
func (g *CodeGeneratorBase) ProcessEnd(n rhino.Node, context CodeGenerator_Context) {
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
