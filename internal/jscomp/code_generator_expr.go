// Port of com/google/javascript/jscomp/CodeGenerator.java @ closure-compiler v20261028

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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// addExprOrTypeNode handles expression, literal, destructuring pattern,
// spread/rest, and TypeScript type-expression nodes.
func (g *CodeGenerator) addExprOrTypeNode(node rhino.Node, context CodeGenerator_Context) {
	cc := &g.cc
	type_ := node.GetToken()
	switch type_ {
	case rhino.Token_LABEL_NAME:
		jutil.CheckStateMsg1(!node.HasChildren(), "%s", node)
		lbl := ast.MustLabelName(node)
		jutil.CheckStateMsg1(lbl.Ident() != "", "%s", node)
		g.AddIdentifier(lbl.Ident())
	case rhino.Token_DESTRUCTURING_LHS:
		g.addDestructuringLhs(ast.MustDestructuringLhs(node), context)
	case rhino.Token_NAME:
		g.addName(ast.MustName(node), context)
	case rhino.Token_ARRAYLIT:
		g.Add("[")
		g.AddArrayList(ast.MustArrayLit(node).Node().GetFirstChild())
		g.Add("]")
	case rhino.Token_ARRAY_PATTERN:
		g.Add("[")
		g.AddArrayList(ast.MustArrayPattern(node).Node().GetFirstChild())
		g.Add("]")
		g.MaybeAddTypeDecl(node)
	case rhino.Token_DEFAULT_VALUE:
		dv := ast.MustDefaultValue(node)
		g.AddNode(dv.Target().Node())
		g.MaybeAddTypeDecl(node)
		cc.AddOp("=", true)
		g.AddExpr(dv.Value().Node(), 1, CodeGenerator_Context_OTHER)
	case rhino.Token_COMMA:
		jutil.CheckStateMsg1(node.HasTwoChildren(), "%s", node)
		g.UnrollBinaryOperator(
			node, rhino.Token_COMMA, ",", context, CodeGenerator_GetContextForNoInOperator(context), 0, 0)
	case rhino.Token_NUMBER:
		cc.AddNumber(ast.MustNumber(node).Value(), node)
	case rhino.Token_BIGINT:
		jutil.CheckStateMsg1(!node.HasChildren(), "%s", node)
		cc.AddBigInt(node.GetBigInt())
	case rhino.Token_TYPEOF, rhino.Token_VOID, rhino.Token_NOT, rhino.Token_BITNOT, rhino.Token_POS, rhino.Token_NEG:
		// All of these unary operators are right-associative
		cc.AddOp(NodeUtil_OpToStrNoFail(type_), false)
		g.AddExpr(ast.MustUnary(node).Operand().Node(), NodeUtil_Precedence(type_), CodeGenerator_Context_OTHER)
	case rhino.Token_DELPROP:
		g.Add("delete ")
		g.AddNode(ast.MustUnary(node).Operand().Node())
	case rhino.Token_INC, rhino.Token_DEC:
		g.addIncDec(ast.MustIncDec(node), type_, context)
	case rhino.Token_HOOK:
		g.addHook(ast.MustHook(node), context)
	case rhino.Token_REGEXP:
		g.addRegExpNode(ast.MustRegExp(node))
	case rhino.Token_ITER_REST, rhino.Token_OBJECT_REST:
		g.Add("...")
		g.AddNode(ast.MustRest(node).Target().Node())
		g.MaybeAddTypeDecl(node)
	case rhino.Token_ITER_SPREAD, rhino.Token_OBJECT_SPREAD:
		g.Add("...")
		g.AddExpr(ast.MustSpread(node).Expr().Node(), NodeUtil_Precedence(type_), CodeGenerator_Context_OTHER)
	case rhino.Token_DYNAMIC_IMPORT:
		g.Add("import(")
		g.AddExpr(ast.MustDynamicImport(node).Specifier().Node(), NodeUtil_Precedence(type_), context)
		g.Add(")")
	case rhino.Token_IMPORT_META:
		jutil.CheckStateMsg1(!node.HasChildren(), "%s", node)
		g.Add("import.meta")
	case rhino.Token_OPTCHAIN_GETPROP:
		g.addOptChainGetProp(ast.MustOptChainGetProp(node), context)
	case rhino.Token_GETPROP:
		g.addGetProp(ast.MustGetProp(node), context)
	case rhino.Token_OPTCHAIN_GETELEM:
		g.addOptChainGetElem(ast.MustOptChainGetElem(node), context)
	case rhino.Token_GETELEM:
		g.addGetElem(ast.MustGetElem(node), context)
	case rhino.Token_OPTCHAIN_CALL:
		g.addOptChainCall(ast.MustOptChainCall(node), context)
	case rhino.Token_CALL:
		g.addCall(ast.MustCall(node), context)
	case rhino.Token_NEW:
		g.addNew(ast.MustNew(node))
	case rhino.Token_NULL:
		jutil.CheckStateMsg1(!node.HasChildren(), "%s", node)
		cc.AddConstant("null")
	case rhino.Token_THIS:
		jutil.CheckStateMsg1(!node.HasChildren(), "%s", node)
		g.Add("this")
	case rhino.Token_SUPER:
		jutil.CheckStateMsg1(!node.HasChildren(), "%s", node)
		g.Add("super")
	case rhino.Token_NEW_TARGET:
		jutil.CheckStateMsg1(!node.HasChildren(), "%s", node)
		g.Add("new.target")
	case rhino.Token_YIELD:
		g.addYield(ast.MustYield(node))
	case rhino.Token_AWAIT:
		g.Add("await ")
		g.AddExpr(ast.MustAwait(node).Value().Node(), NodeUtil_Precedence(type_), CodeGenerator_Context_OTHER)
	case rhino.Token_FALSE:
		jutil.CheckStateMsg1(!node.HasChildren(), "%s", node)
		cc.AddConstant("false")
	case rhino.Token_TRUE:
		jutil.CheckStateMsg1(!node.HasChildren(), "%s", node)
		cc.AddConstant("true")
	case rhino.Token_STRING_KEY:
		g.AddStringKey(node)
	case rhino.Token_STRINGLIT:
		jutil.CheckStateMsg1(!node.HasChildren(), "String node %s may not have children", node)
		g.AddJsString(node)
	case rhino.Token_OBJECTLIT:
		g.addObjectLit(ast.MustObjectLit(node), context)
	case rhino.Token_COMPUTED_PROP:
		g.addComputedProp(ast.MustComputedProp(node))
	case rhino.Token_OBJECT_PATTERN:
		g.AddObjectPattern(node)
		g.MaybeAddTypeDecl(node)
	case rhino.Token_CAST:
		g.addCast(ast.MustCast(node), context)
	case rhino.Token_TAGGED_TEMPLATELIT:
		g.AddInvocationTarget(node, context)
		g.AddNode(ast.MustTaggedTemplateLit(node).Template().Node())
	case rhino.Token_TEMPLATELIT:
		g.addTemplateLit(ast.MustTemplateLit(node))
	// Type Declaration ASTs.
	case rhino.Token_STRING_TYPE:
		g.Add("string")
	case rhino.Token_BOOLEAN_TYPE:
		g.Add("boolean")
	case rhino.Token_NUMBER_TYPE:
		g.Add("number")
	case rhino.Token_ANY_TYPE:
		g.Add("any")
	case rhino.Token_VOID_TYPE:
		g.Add("void")
	case rhino.Token_NAMED_TYPE:
		g.AddNode(node.GetFirstChild())
	case rhino.Token_ARRAY_TYPE:
		g.AddExpr(node.GetFirstChild(), NodeUtil_Precedence(rhino.Token_ARRAY_TYPE), context)
		g.Add("[]")
	case rhino.Token_FUNCTION_TYPE:
		first := node.GetFirstChild()
		g.Add("(")
		g.AddList(first.GetNext())
		g.Add(")")
		cc.AddOp("=>", true)
		g.AddNode(first)
	case rhino.Token_UNION_TYPE:
		g.AddListWithSeparator(node.GetFirstChild(), "|")
	case rhino.Token_RECORD_TYPE:
		g.Add("{")
		g.AddListWithContextAndSeparator(node.GetFirstChild(), false, CodeGenerator_Context_OTHER, ",")
		g.Add("}")
	case rhino.Token_PARAMETERIZED_TYPE:
		first := node.GetFirstChild()
		g.AddNode(first)
		g.Add("<")
		g.AddList(first.GetNext())
		g.Add(">")
	case rhino.Token_GENERIC_TYPE_LIST:
		g.Add("<")
		g.AddListWithContextAndSeparator(node.GetFirstChild(), false, CodeGenerator_Context_STATEMENT, ",")
		g.Add(">")
	case rhino.Token_GENERIC_TYPE:
		g.AddIdentifier(node.GetString())
		if node.HasChildren() {
			g.Add("extends")
			cc.MaybeInsertSpace()
			g.AddNode(node.GetFirstChild())
		}
	default:
		panic(jutil.NewIllegalStateExceptionMsg("Unknown token " + type_.String() + "\n" + node.ToStringTree()))
	}
}

func (g *CodeGenerator) addDestructuringLhs(lhs ast.DestructuringLhs, context CodeGenerator_Context) {
	g.AddNode(lhs.Pattern().Node())
	if init, ok := lhs.Init(); ok {
		g.cc.AddOp("=", true)
		g.AddExpr(init.Node(), NodeUtil_Precedence(rhino.Token_ASSIGN), CodeGenerator_GetContextForNoInOperator(context))
	}
}

func (g *CodeGenerator) addName(nm ast.Name, context CodeGenerator_Context) {
	node := nm.Node()
	if orig, ok := node.OriginalName(); g.useOriginalName && ok {
		g.AddIdentifier(orig)
	} else {
		g.AddIdentifier(nm.Ident())
	}
	g.MaybeAddOptional(node)
	g.MaybeAddTypeDecl(node)

	if init, ok := nm.Init(); ok && !init.Node().IsEmpty() {
		g.cc.AddOp("=", true)
		g.AddExpr(init.Node(), NodeUtil_Precedence(rhino.Token_ASSIGN), CodeGenerator_GetContextForNoInOperator(context))
	}
}

func (g *CodeGenerator) addIncDec(incDec ast.IncDec, type_ rhino.Token, context CodeGenerator_Context) {
	o := "--"
	if type_ == rhino.Token_INC {
		o = "++"
	}
	target := incDec.Target().Node()
	if incDec.IsPostfix() {
		g.AddExpr(target, NodeUtil_Precedence(type_), context)
		g.cc.AddOp(o, false)
	} else {
		g.cc.AddOp(o, false)
		g.AddNode(target)
	}
}

func (g *CodeGenerator) addHook(hook ast.Hook, context CodeGenerator_Context) {
	p := NodeUtil_Precedence(rhino.Token_HOOK)
	rhsContext := CodeGenerator_GetContextForNoInOperator(context)
	g.AddExpr(hook.Cond().Node(), p+1, context)
	g.cc.AddOp("?", true)
	g.AddExpr(hook.Then().Node(), 1, rhsContext)
	g.cc.AddOp(":", true)
	g.AddExpr(hook.Else().Node(), 1, rhsContext)
}

func (g *CodeGenerator) addRegExpNode(re ast.RegExp) {
	pat := re.Pattern()
	if !pat.Node().IsStringLit() {
		panic(jutil.NewJavaLangError("Expected children to be strings"))
	}
	flags := ""
	if f, ok := re.Flags(); ok {
		if !f.Node().IsStringLit() {
			panic(jutil.NewJavaLangError("Expected children to be strings"))
		}
		flags = f.Value()
	}
	g.addRegexp(pat.Value(), flags)
}

func (g *CodeGenerator) addOptChainGetProp(gp ast.OptChainGetProp, context CodeGenerator_Context) {
	g.AddExpr(gp.Object().Node(), NodeUtil_Precedence(rhino.Token_OPTCHAIN_GETPROP), context)
	if gp.IsStart() {
		g.Add("?.")
	} else {
		g.Add(".")
	}
	g.AddGetpropIdentifier(gp.Node())
}

func (g *CodeGenerator) addGetProp(gp ast.GetProp, context CodeGenerator_Context) {
	node := gp.Node()
	objNode := gp.Object().Node()
	// Convert rewritten aliased code back to original code when using goog.scope().
	if _, ok := node.OriginalName(); g.useOriginalName && ok {
		if codeGenerator_JSCOMP_SCOPE.Matches(objNode) && node.GetParent().IsAssign() {
			g.Add("var ")
		}
		g.AddGetpropIdentifier(node)
		return
	}
	// We need parentheses to distinguish `a?.b.c` from `(a?.b).c`, and
	// `2.toString()` must be `(2).toString()`.
	needsParens := objNode.IsNumber() || NodeUtil_IsOptChainNode(objNode)
	if needsParens {
		g.Add("(")
	}
	g.AddExpr(objNode, NodeUtil_Precedence(rhino.Token_GETPROP), context)
	if needsParens {
		g.Add(")")
	}
	prop := gp.Property()
	if g.quoteKeywordProperties && rhino.TokenStream_IsKeyword(prop) {
		g.Add("[")
		g.cc.addPlainQuotedASCII('"', prop)
		g.Add("]")
	} else {
		g.Add(".")
		g.AddGetpropIdentifier(node)
	}
}

func (g *CodeGenerator) addOptChainGetElem(ge ast.OptChainGetElem, context CodeGenerator_Context) {
	g.AddExpr(ge.Object().Node(), NodeUtil_Precedence(rhino.Token_OPTCHAIN_GETELEM), context)
	if ge.IsStart() {
		g.Add("?.")
	}
	g.Add("[")
	g.AddNode(ge.Index().Node())
	g.Add("]")
}

func (g *CodeGenerator) addGetElem(ge ast.GetElem, context CodeGenerator_Context) {
	objNode := ge.Object().Node()
	needsParens := NodeUtil_IsOptChainNode(objNode)
	if needsParens {
		g.Add("(")
	}
	g.AddExpr(objNode, NodeUtil_Precedence(rhino.Token_GETELEM), context)
	if needsParens {
		g.Add(")")
	}
	g.Add("[")
	g.AddNode(ge.Index().Node())
	g.Add("]")
}

func (g *CodeGenerator) addOptChainCall(call ast.OptChainCall, context CodeGenerator_Context) {
	calleeNode := call.Callee().Node()
	if CodeGenerator_IsIndirectEval(calleeNode) ||
		(call.IsFreeCall() && NodeUtil_IsNormalOrOptChainGet(calleeNode)) {
		g.Add("(0,")
		g.AddExpr(calleeNode, NodeUtil_Precedence(rhino.Token_COMMA), CodeGenerator_Context_OTHER)
		g.Add(")")
	} else {
		g.AddExpr(calleeNode, NodeUtil_Precedence(rhino.Token_OPTCHAIN_CALL), context)
	}
	if call.IsStart() {
		g.Add("?.")
	}
	g.Add("(")
	if firstArg, ok := call.FirstArg(); ok {
		g.AddList(firstArg.Node())
	}
	g.Add(")")
}

func (g *CodeGenerator) addCall(call ast.Call, context CodeGenerator_Context) {
	g.AddInvocationTarget(call.Node(), context)
	g.Add("(")
	if firstArg, ok := call.FirstArg(); ok {
		g.AddList(firstArg.Node())
	}
	g.Add(")")
}

func (g *CodeGenerator) addNew(nw ast.New) {
	calleeNode := nw.Callee().Node()
	g.Add("new ")
	precedence := NodeUtil_Precedence(rhino.Token_NEW)

	// `new void 0` is a syntax error; add parentheses in this case.
	if NodeUtil_Precedence(calleeNode.GetToken()) == precedence {
		precedence = precedence + 1
	}

	// If the callee contains a CALL (outside functions) or is an optional
	// chain, force parentheses.
	if NodeUtil_Has(calleeNode, (rhino.Node).IsCall, NodeUtil_MATCH_NOT_FUNCTION) ||
		NodeUtil_IsOptChainNode(calleeNode) {
		precedence = NodeUtil_Precedence(calleeNode.GetToken()) + 1
	}
	g.AddExpr(calleeNode, precedence, CodeGenerator_Context_OTHER)

	// '()' is optional when no arguments are present
	if firstArg, ok := nw.FirstArg(); ok {
		g.Add("(")
		g.AddList(firstArg.Node())
		g.Add(")")
	} else if g.cc.ShouldPreserveExtras(nw.Node()) {
		g.Add("(")
		g.Add(")")
	}
}

func (g *CodeGenerator) addYield(yld ast.Yield) {
	g.Add("yield")
	val, hasVal := yld.Value()
	if yld.IsYieldAll() {
		jutil.CheckState(hasVal)
		g.Add("*")
	}
	if hasVal {
		g.cc.MaybeInsertSpace()
		g.AddExpr(val.Node(), NodeUtil_Precedence(rhino.Token_YIELD), CodeGenerator_Context_OTHER)
	}
}

func (g *CodeGenerator) addObjectLit(lit ast.ObjectLit, context CodeGenerator_Context) {
	needsParens := context == CodeGenerator_Context_START_OF_EXPR || context.AtArrowFunctionBody()
	if needsParens {
		g.Add("(")
	}
	g.Add("{")
	first := lit.Node().GetFirstChild()
	for c := first; !c.IsNil(); c = c.GetNext() {
		if c != first {
			g.cc.ListSeparator()
		}
		jutil.CheckStateMsg1(NodeUtil_IsObjLitProperty(c) || c.IsSpread(), "%s", c)
		g.AddNode(c)
	}
	if !first.IsNil() && lit.Node().HasTrailingComma() {
		g.cc.OptionalListSeparator()
	}
	g.Add("}")
	if needsParens {
		g.Add(")")
	}
}

func (g *CodeGenerator) addComputedProp(cp ast.ComputedProp) {
	node := cp.Node()
	g.MaybeAddAccessibilityModifier(node)
	if cp.IsStatic() {
		g.Add("static ")
	}

	val, hasVal := cp.Value()
	isGetter := cp.IsGetter()
	isSetter := cp.IsSetter()
	isMethod := cp.IsMethod()
	if isGetter {
		g.Add("get ")
	} else if isSetter {
		g.Add("set ")
	} else if isMethod {
		if val.Node().IsAsyncFunction() {
			g.Add("async")
		}
		if val.Node().GetBooleanProp(rhino.Node_GENERATOR_FN) {
			g.Add("*")
		}
	}
	g.Add("[")
	// Must use addExpr() with a priority of 1, because comma expressions
	// aren't allowed.
	g.AddExpr(cp.Key().Node(), 1, CodeGenerator_Context_OTHER)
	g.Add("]")
	g.MaybeAddTypeDecl(node)
	if isMethod || isGetter || isSetter {
		fn := ast.MustFunction(val.Node())
		g.AddNode(fn.Params().Node())
		g.AddNode(fn.Body().Node())
	} else {
		isInClass := node.GetParent().IsClassMembers()
		if hasVal {
			jutil.CheckStateMsg(
				!isInClass, "initializers should only exist in object literals, not classes")
			g.cc.Add(":")
			g.AddExpr(val.Node(), 1, CodeGenerator_Context_OTHER)
		} else {
			jutil.CheckStateMsg1(node.GetBooleanProp(rhino.Node_COMPUTED_PROP_VARIABLE), "%s", node)
		}
	}
}

func (g *CodeGenerator) addCast(cast ast.Cast, context CodeGenerator_Context) {
	exprNode := cast.Expr().Node()
	if g.preserveTypeAnnotations {
		g.Add("(")
		g.AddNode(exprNode) // drop context because of added parentheses
		g.Add(")")
	} else {
		g.AddNodeWithContext(exprNode, context) // preserve context
	}
}

func (g *CodeGenerator) addTemplateLit(tl ast.TemplateLit) {
	cc := &g.cc
	cc.BeginTemplateLit()
	for part, ok := tl.FirstPart(); ok; part, ok = part.Next() {
		if tls, isStr := part.AsTemplateLitString(); isStr {
			g.Add(g.EscapeUnrecognizedCharacters(tls.Raw()))
		} else {
			jutil.CheckStateMsg1(part.Node().HasOneChild(), "%s", part.Node())
			sub := part.MustTemplateLitSub()
			cc.BeginTemplateLitSub()
			g.AddNodeWithContext(sub.Expr().Node(), CodeGenerator_Context_START_OF_EXPR)
			cc.EndTemplateLitSub()
		}
	}
	cc.EndTemplateLit()
}

// Java: CodeGenerator.addIdentifier(String)
func (g *CodeGenerator) AddIdentifier(identifier string) {
	if g.cc.mode == sinkCost {
		g.cc.Add(InlineCostEstimator_ESTIMATED_IDENTIFIER)
		return
	}
	if NodeUtil_IsLatin(identifier) {
		g.cc.Add(identifier)
		return
	}
	g.addNonLatinIdentifier(identifier)
}

// Java: CodeGenerator.addGetpropIdentifier(Node)
func (g *CodeGenerator) AddGetpropIdentifier(getprop rhino.Node) {
	g.cc.StartSourceMapping(getprop)
	g.AddIdentifier(getprop.GetString())
	g.cc.EndSourceMapping(getprop)
}

// AddInvocationTarget handles indirect eval and free calls on property
// accesses before invoking the target.
//
// Java: CodeGenerator.addInvocationTarget(Node, Context)
func (g *CodeGenerator) AddInvocationTarget(node rhino.Node, context CodeGenerator_Context) {
	first := ast.MustInvocation(node).Callee().Node()

	needsParens := NodeUtil_IsOptChainNode(first)
	if CodeGenerator_IsIndirectEval(first) ||
		(node.GetBooleanProp(rhino.Node_FREE_CALL) && NodeUtil_IsNormalOrOptChainGet(first)) {
		g.Add("(0,")
		g.AddExpr(first, NodeUtil_Precedence(rhino.Token_COMMA), CodeGenerator_Context_OTHER)
		g.Add(")")
	} else {
		if needsParens {
			g.Add("(")
		}
		g.AddExpr(first, NodeUtil_Precedence(node.GetToken()), context)
		if needsParens {
			g.Add(")")
		}
	}
}

// UnrollBinaryOperator iteratively emits left-recursive binary operator chains
// without exhausting the call stack on deeply nested expressions.
//
// Java: CodeGenerator.unrollBinaryOperator(Node, Token, String, Context, Context, int, int)
func (g *CodeGenerator) UnrollBinaryOperator(
	n rhino.Node,
	op rhino.Token,
	opStr string,
	context CodeGenerator_Context,
	rhsContext CodeGenerator_Context,
	leftPrecedence int,
	rightPrecedence int,
) {
	firstNonOperator := ast.MustBinary(n).Left().Node()
	for firstNonOperator.GetToken() == op {
		firstNonOperator = ast.MustBinary(firstNonOperator).Left().Node()
	}

	g.AddExpr(firstNonOperator, leftPrecedence, context)

	current := firstNonOperator
	for {
		current = current.GetParent()
		g.cc.AddOp(opStr, true)
		g.AddExpr(ast.MustBinary(current).Right().Node(), rightPrecedence, rhsContext)
		if current == n {
			break
		}
	}
}

// Java: CodeGenerator.isSimpleNumber(String)
func CodeGenerator_IsSimpleNumber(s string) bool {
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
func codeGenerator_getSimpleNumber(s string) (float64, bool) {
	if CodeGenerator_IsSimpleNumber(s) {
		l, err := strconv.ParseInt(s, 10, 64)
		if err == nil && l <= NodeUtil_MAX_POSITIVE_INTEGER_NUMBER {
			return float64(l), true
		}
	}
	return 0, false
}

func (g *CodeGenerator) addNumericOrStringKey(key string, n rhino.Node) {
	if d, ok := codeGenerator_getSimpleNumber(key); ok {
		g.cc.AddNumber(d, n)
	} else {
		g.AddJsString(n)
	}
}

var codeGenIDEval = strpool.Intern("eval")

// CodeGenerator_IsIndirectEval returns whether the name is an indirect eval.
//
// Java: CodeGenerator.isIndirectEval(Node)
func CodeGenerator_IsIndirectEval(n rhino.Node) bool {
	return n.IsName() && nodeStringID(n) == codeGenIDEval && !n.GetBooleanProp(rhino.Node_DIRECT_EVAL)
}

// Java: CodeGenerator.addList(Node)
func (g *CodeGenerator) AddList(firstInList rhino.Node) {
	g.AddListWithContextAndSeparator(firstInList, true, CodeGenerator_Context_OTHER, ",")
}

// Java: CodeGenerator.addList(Node, String)
func (g *CodeGenerator) AddListWithSeparator(firstInList rhino.Node, separator string) {
	g.AddListWithContextAndSeparator(firstInList, true, CodeGenerator_Context_OTHER, separator)
}

// Java: CodeGenerator.addList(Node, boolean, Context, String)
func (g *CodeGenerator) AddListWithContextAndSeparator(
	firstInList rhino.Node, isArrayOrFunctionArgument bool, lhsContext CodeGenerator_Context, separator string,
) {
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
func (g *CodeGenerator) AddStringKey(n rhino.Node) {
	sk := ast.MustStringKey(n)
	key := sk.Key()
	// Object literal property names don't have to be quoted if they are not
	// JavaScript keywords.
	mustBeQuoted := sk.IsQuoted() ||
		(g.quoteKeywordProperties && rhino.TokenStream_IsKeyword(key)) ||
		!rhino.TokenStream_IsJSIdentifier(key) ||
		!NodeUtil_IsLatin(key)
	if !mustBeQuoted {
		// Check if the property is eligible to be printed as shorthand.
		if n.IsShorthandProperty() {
			val := sk.Value()
			child := val.Node()
			dv, isDv := val.AsDefaultValue()
			if child.MatchesQualifiedName(key) ||
				(isDv && dv.Target().Node().MatchesQualifiedName(key)) {
				g.AddNode(child)
				return
			}
		}
		g.Add(key)
	} else {
		g.addNumericOrStringKey(key, n)
	}
	if n.HasChildren() {
		// NOTE: the only time a STRING_KEY node does *not* have children is
		// when it's inside a TypeScript enum.
		g.Add(":")
		g.AddExpr(sk.Value().Node(), 1, CodeGenerator_Context_OTHER)
	}
}

// Java: CodeGenerator.addObjectPattern(Node)
func (g *CodeGenerator) AddObjectPattern(n rhino.Node) {
	pat := ast.MustObjectPattern(n)
	g.Add("{")
	firstProp, _ := pat.FirstProp()
	for child := firstProp.Node(); !child.IsNil(); child = child.GetNext() {
		if child != firstProp.Node() {
			g.cc.ListSeparator()
		}
		g.AddNode(child)
	}
	g.Add("}")
}

// AddArrayList adds a comma-separated list as is specified by an ARRAYLIT
// or ARRAY_PATTERN node.
//
// Java: CodeGenerator.addArrayList(Node)
func (g *CodeGenerator) AddArrayList(firstInList rhino.Node) {
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
