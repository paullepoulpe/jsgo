// Port of com/google/javascript/rhino/IR.java @ closure-compiler v20260415

package rhino

// IR is an AST construction helper. All methods are package-level functions
// because Java statics become free functions in Go. Mirrors
// com.google.javascript.rhino.IR.
//
// Naming convention: every function carries the "IR" prefix so that callers
// can write IRBlock(...) instead of block(...), avoiding clashes with the
// many same-named methods on *Node. The Java class uses plain static names
// (e.g. IR.block); the prefix is a language-adapter choice.

// Go-only: language-adapter — newNodes is a variadic helper used by every
// builder below; Java has no equivalent — each builder inlines new Node(tok, children).
func newNodes(tok Token, children ...*Node) *Node {
	n := New(tok)
	for _, c := range children {
		n.AddChildToBack(c)
	}
	return n
}

// Go-only: language-adapter — irBinaryOp ports IR.binaryOp(Token, Node, Node).
// Private in Java; exposed here because Go lacks private package helpers.
func irBinaryOp(tok Token, expr1, expr2 *Node) *Node {
	return newNodes(tok, expr1, expr2)
}

// Go-only: language-adapter — irUnaryOp ports IR.unaryOp(Token, Node).
// Private in Java; exposed here because Go lacks private package helpers.
func irUnaryOp(tok Token, expr *Node) *Node {
	return newNodes(tok, expr)
}

// Go-only: language-adapter — IREmpty ports IR.empty(). The IR prefix
// prevents collision with the Token constant EMPTY in the same package.
func IREmpty() *Node { return New(EMPTY) }

// Go-only: language-adapter — IRExport ports IR.export(Node).
func IRExport(declaration *Node) *Node { return newNodes(EXPORT, declaration) }

// Go-only: language-adapter — IRImportNode ports IR.importNode(Node, Node, Node).
func IRImportNode(name, importSpecs, moduleIdentifier *Node) *Node {
	return newNodes(IMPORT, name, importSpecs, moduleIdentifier)
}

// Go-only: language-adapter — IRImportStar ports IR.importStar(String).
func IRImportStar(name string) *Node { return NewString(IMPORT_STAR, name) }

// Go-only: language-adapter — IRFunction ports IR.function(Node, Node, Node).
func IRFunction(name, params, body *Node) *Node {
	return newNodes(FUNCTION, name, params, body)
}

// Go-only: language-adapter — IRArrowFunction ports IR.arrowFunction(Node, Node, Node).
func IRArrowFunction(name, params, body *Node) *Node {
	fn := newNodes(FUNCTION, name, params, body)
	fn.SetIsArrowFunction(true)
	return fn
}

// Go-only: language-adapter — IRParamList ports IR.paramList(Node...).
func IRParamList(params ...*Node) *Node {
	return newNodes(PARAM_LIST, params...)
}

// Go-only: language-adapter — IRRoot ports IR.root(Node...).
func IRRoot(rootChildren ...*Node) *Node {
	return newNodes(ROOT, rootChildren...)
}

// Go-only: language-adapter — IRBlock ports IR.block() / IR.block(Node...).
// Java has four overloads; Go's variadic unifies them.
func IRBlock(stmts ...*Node) *Node {
	return newNodes(BLOCK, stmts...)
}

// Go-only: language-adapter — IRScript ports IR.script() / IR.script(Node...).
func IRScript(stmts ...*Node) *Node {
	return newNodes(SCRIPT, stmts...)
}

// Go-only: language-adapter — IRVar ports IR.var(Node, Node).
func IRVar(lhs, value *Node) *Node { return IRDeclaration(lhs, value, VAR) }

// Go-only: language-adapter — IRVarLhs ports IR.var(Node) (lhs-only overload).
func IRVarLhs(lhs *Node) *Node { return IRDeclarationLhs(lhs, VAR) }

// Go-only: language-adapter — IRLet ports IR.let(Node, Node).
func IRLet(lhs, value *Node) *Node { return IRDeclaration(lhs, value, LET) }

// Go-only: language-adapter — IRLetLhs ports IR.let(Node) (lhs-only overload).
func IRLetLhs(lhs *Node) *Node { return IRDeclarationLhs(lhs, LET) }

// Go-only: language-adapter — IRConst ports IR.constNode(Node, Node).
// Renamed from Java's constNode to avoid collision with the CONST Token.
func IRConst(lhs, value *Node) *Node { return IRDeclaration(lhs, value, CONST) }

// Go-only: language-adapter — IRDeclarationLhs ports IR.declaration(Node, Token).
func IRDeclarationLhs(lhs *Node, typ Token) *Node {
	return newNodes(typ, lhs)
}

// Go-only: language-adapter — IRDeclaration ports IR.declaration(Node, Node, Token).
func IRDeclaration(lhs, value *Node, typ Token) *Node {
	lhs.AddChildToBack(value)
	return newNodes(typ, lhs)
}

// Go-only: language-adapter — IRReturn ports IR.returnNode() / IR.returnNode(Node).
// Java has two overloads; Go's variadic unifies them.
func IRReturn(exprs ...*Node) *Node {
	if len(exprs) == 0 {
		return New(RETURN)
	}
	return newNodes(RETURN, exprs[0])
}

// Go-only: language-adapter — IRYield ports IR.yieldNode(Node).
// Renamed: Java uses yieldNode; Go drops the Node suffix for brevity.
func IRYield(expr *Node) *Node { return newNodes(YIELD, expr) }

// Go-only: language-adapter — IRAwait ports IR.await(Node).
func IRAwait(expr *Node) *Node { return newNodes(AWAIT, expr) }

// Go-only: language-adapter — IRThrow ports IR.throwNode(Node).
func IRThrow(expr *Node) *Node { return newNodes(THROW, expr) }

// Go-only: language-adapter — IRExprResult ports IR.exprResult(Node).
func IRExprResult(expr *Node) *Node { return newNodes(EXPR_RESULT, expr) }

// Go-only: language-adapter — IRIf ports IR.ifNode(Node, Node) / IR.ifNode(Node, Node, Node).
// Java has two overloads; Go's variadic else unifies them.
func IRIf(cond, then *Node, elseNodes ...*Node) *Node {
	n := newNodes(IF, cond, then)
	for _, e := range elseNodes {
		n.AddChildToBack(e)
	}
	return n
}

// Go-only: language-adapter — IRDo ports IR.doNode(Node, Node).
func IRDo(body, cond *Node) *Node { return newNodes(DO, body, cond) }

// Go-only: language-adapter — IRWhile ports IR.whileNode(Node, Node).
func IRWhile(cond, body *Node) *Node { return newNodes(WHILE, cond, body) }

// Go-only: language-adapter — IRForIn ports IR.forIn(Node, Node, Node).
func IRForIn(target, cond, body *Node) *Node { return newNodes(FOR_IN, target, cond, body) }

// Go-only: language-adapter — IRFor ports IR.forNode(Node, Node, Node, Node).
func IRFor(init, cond, incr, body *Node) *Node {
	n := newNodes(FOR, init, cond, incr)
	n.AddChildToBack(body)
	return n
}

// Go-only: language-adapter — IRSwitch ports IR.switchNode(Node, Node...).
func IRSwitch(cond *Node, cases ...*Node) *Node {
	n := newNodes(SWITCH, cond)
	body := New(SWITCH_BODY)
	n.AddChildToBack(body)
	for _, c := range cases {
		body.AddChildToBack(c)
	}
	return n
}

// Go-only: language-adapter — IRCase ports IR.caseNode(Node, Node).
func IRCase(expr, body *Node) *Node { return newNodes(CASE, expr, body) }

// Go-only: language-adapter — IRDefaultCase ports IR.defaultCase(Node).
func IRDefaultCase(body *Node) *Node { return newNodes(DEFAULT_CASE, body) }

// Go-only: language-adapter — IRLabel ports IR.label(Node, Node).
func IRLabel(name, stmt *Node) *Node { return newNodes(LABEL, name, stmt) }

// Go-only: language-adapter — IRLabelName ports IR.labelName(String).
func IRLabelName(name string) *Node { return NewString(LABEL_NAME, name) }

// Go-only: language-adapter — IRTryFinally ports IR.tryFinally(Node, Node).
func IRTryFinally(tryBody, finallyBody *Node) *Node {
	catchBody := New(BLOCK)
	return newNodes(TRY, tryBody, catchBody, finallyBody)
}

// Go-only: language-adapter — IRTryCatch ports IR.tryCatch(Node, Node).
func IRTryCatch(tryBody, catchNode *Node) *Node {
	catchBody := newNodes(BLOCK, catchNode)
	return newNodes(TRY, tryBody, catchBody)
}

// Go-only: language-adapter — IRTryCatchFinally ports IR.tryCatchFinally(Node, Node, Node).
func IRTryCatchFinally(tryBody, catchNode, finallyBody *Node) *Node {
	n := IRTryCatch(tryBody, catchNode)
	n.AddChildToBack(finallyBody)
	return n
}

// Go-only: language-adapter — IRCatch ports IR.catchNode(Node, Node).
func IRCatch(expr, body *Node) *Node { return newNodes(CATCH, expr, body) }

// Go-only: language-adapter — IRBreak ports IR.breakNode() / IR.breakNode(Node).
func IRBreak(labels ...*Node) *Node {
	if len(labels) == 0 {
		return New(BREAK)
	}
	return newNodes(BREAK, labels[0])
}

// Go-only: language-adapter — IRContinue ports IR.continueNode() / IR.continueNode(Node).
func IRContinue(labels ...*Node) *Node {
	if len(labels) == 0 {
		return New(CONTINUE)
	}
	return newNodes(CONTINUE, labels[0])
}

// Go-only: language-adapter — IRCall ports IR.call(Node, Node...).
func IRCall(target *Node, args ...*Node) *Node {
	n := newNodes(CALL, target)
	for _, a := range args {
		n.AddChildToBack(a)
	}
	return n
}

// Go-only: language-adapter — IRStartOptChainCall ports IR.startOptChainCall(Node, Node...).
func IRStartOptChainCall(target *Node, args ...*Node) *Node {
	n := newNodes(OPTCHAIN_CALL, target)
	for _, a := range args {
		n.AddChildToBack(a)
	}
	n.SetIsOptionalChainStart(true)
	return n
}

// Go-only: language-adapter — IRContinueOptChainCall ports IR.continueOptChainCall(Node, Node...).
func IRContinueOptChainCall(target *Node, args ...*Node) *Node {
	n := newNodes(OPTCHAIN_CALL, target)
	for _, a := range args {
		n.AddChildToBack(a)
	}
	n.SetIsOptionalChainStart(false)
	return n
}

// Go-only: language-adapter — IRNew ports IR.newNode(Node, Node...).
// Renamed from Java's newNode to avoid collision with Go's built-in new.
func IRNew(target *Node, args ...*Node) *Node {
	n := newNodes(NEW, target)
	for _, a := range args {
		n.AddChildToBack(a)
	}
	return n
}

// Go-only: language-adapter — IRName ports IR.name(String).
func IRName(name string) *Node { return NewString(NAME, name) }

// Go-only: language-adapter — IRStartOptChainGetprop ports IR.startOptChainGetprop(Node, String).
func IRStartOptChainGetprop(target *Node, prop string) *Node {
	n := NewString(OPTCHAIN_GETPROP, prop)
	n.AddChildToBack(target)
	n.SetIsOptionalChainStart(true)
	return n
}

// Go-only: language-adapter — IRContinueOptChainGetprop ports IR.continueOptChainGetprop(Node, String).
func IRContinueOptChainGetprop(target *Node, prop string) *Node {
	n := NewString(OPTCHAIN_GETPROP, prop)
	n.AddChildToBack(target)
	n.SetIsOptionalChainStart(false)
	return n
}

// Go-only: language-adapter — IRGetprop ports IR.getprop(Node, String).
// The variadic moreProps extension is a Go convenience; Java callers chain manually.
func IRGetprop(target *Node, prop string, moreProps ...string) *Node {
	n := NewString(GETPROP, prop)
	n.AddChildToBack(target)
	result := n
	for _, mp := range moreProps {
		result = IRGetprop(result, mp)
	}
	return result
}

// Go-only: language-adapter — IRStartOptChainGetelem ports IR.startOptChainGetelem(Node, Node).
func IRStartOptChainGetelem(target, elem *Node) *Node {
	n := newNodes(OPTCHAIN_GETELEM, target, elem)
	n.SetIsOptionalChainStart(true)
	return n
}

// Go-only: language-adapter — IRContinueOptChainGetelem ports IR.continueOptChainGetelem(Node, Node).
func IRContinueOptChainGetelem(target, elem *Node) *Node {
	n := newNodes(OPTCHAIN_GETELEM, target, elem)
	n.SetIsOptionalChainStart(false)
	return n
}

// Go-only: language-adapter — IRGetelem ports IR.getelem(Node, Node).
func IRGetelem(target, elem *Node) *Node { return newNodes(GETELEM, target, elem) }

// Go-only: language-adapter — IRDelprop ports IR.delprop(Node).
func IRDelprop(target *Node) *Node { return newNodes(DELPROP, target) }

// Go-only: language-adapter — IRAssign ports IR.assign(Node, Node).
func IRAssign(target, expr *Node) *Node { return newNodes(ASSIGN, target, expr) }

// Go-only: language-adapter — IRHook ports IR.hook(Node, Node, Node).
func IRHook(cond, trueval, falseval *Node) *Node {
	return newNodes(HOOK, cond, trueval, falseval)
}

// Go-only: language-adapter — IRIn ports IR.in(Node, Node).
func IRIn(expr1, expr2 *Node) *Node { return irBinaryOp(IN, expr1, expr2) }

// Go-only: language-adapter — IRComma ports IR.comma(Node, Node).
func IRComma(expr1, expr2 *Node) *Node { return irBinaryOp(COMMA, expr1, expr2) }

// Go-only: language-adapter — IRCoalesce ports IR.coalesce(Node, Node).
func IRCoalesce(expr1, expr2 *Node) *Node { return irBinaryOp(COALESCE, expr1, expr2) }

// Go-only: language-adapter — IRLT ports IR.lt(Node, Node) ("<").
func IRLT(expr1, expr2 *Node) *Node { return irBinaryOp(LT, expr1, expr2) }

// Go-only: language-adapter — IRGE ports IR.ge(Node, Node) (">=").
func IRGE(expr1, expr2 *Node) *Node { return irBinaryOp(GE, expr1, expr2) }

// Go-only: language-adapter — IRAssignOr ports IR.assignOr(Node, Node).
func IRAssignOr(expr1, expr2 *Node) *Node { return newNodes(ASSIGN_OR, expr1, expr2) }

// Go-only: language-adapter — IRAssignAnd ports IR.assignAnd(Node, Node).
func IRAssignAnd(expr1, expr2 *Node) *Node { return newNodes(ASSIGN_AND, expr1, expr2) }

// Go-only: language-adapter — IRAssignCoalesce ports IR.assignCoalesce(Node, Node).
func IRAssignCoalesce(expr1, expr2 *Node) *Node { return newNodes(ASSIGN_COALESCE, expr1, expr2) }

// Go-only: language-adapter — IRAnd ports IR.and(Node, Node).
func IRAnd(expr1, expr2 *Node) *Node { return newNodes(AND, expr1, expr2) }

// Go-only: language-adapter — IROr ports IR.or(Node, Node).
func IROr(expr1, expr2 *Node) *Node { return newNodes(OR, expr1, expr2) }

// Go-only: language-adapter — IRNot ports IR.not(Node).
func IRNot(expr *Node) *Node { return newNodes(NOT, expr) }

// Go-only: language-adapter — IREQ ports IR.eq(Node, Node).
func IREQ(expr1, expr2 *Node) *Node { return newNodes(EQ, expr1, expr2) }

// Go-only: language-adapter — IRNE ports IR.ne(Node, Node).
func IRNE(expr1, expr2 *Node) *Node { return newNodes(NE, expr1, expr2) }

// Go-only: language-adapter — IRSHEQ ports IR.sheq(Node, Node).
func IRSHEQ(expr1, expr2 *Node) *Node { return newNodes(SHEQ, expr1, expr2) }

// Go-only: language-adapter — IRSHNE ports IR.shne(Node, Node).
func IRSHNE(expr1, expr2 *Node) *Node { return newNodes(SHNE, expr1, expr2) }

// Go-only: language-adapter — IRVoid ports IR.voidNode(Node).
func IRVoid(expr *Node) *Node { return newNodes(VOID, expr) }

// Go-only: language-adapter — IRNeg ports IR.neg(Node).
func IRNeg(expr *Node) *Node { return newNodes(NEG, expr) }

// Go-only: language-adapter — IRPos ports IR.pos(Node).
func IRPos(expr *Node) *Node { return newNodes(POS, expr) }

// Go-only: language-adapter — IRCast ports IR.cast(Node, JSDocInfo).
func IRCast(expr *Node, jsdoc *JSDocInfo) *Node {
	op := irUnaryOp(CAST, expr)
	op.SetJSDocInfo(jsdoc)
	return op
}

// Go-only: language-adapter — IRInc ports IR.inc(Node, boolean).
func IRInc(exp *Node, isPost bool) *Node {
	op := irUnaryOp(INC, exp)
	op.PutBooleanProp(PropIncrDecr, isPost)
	return op
}

// Go-only: language-adapter — IRDec ports IR.dec(Node, boolean).
func IRDec(exp *Node, isPost bool) *Node {
	op := irUnaryOp(DEC, exp)
	op.PutBooleanProp(PropIncrDecr, isPost)
	return op
}

// Go-only: language-adapter — IRAdd ports IR.add(Node, Node).
func IRAdd(expr1, expr2 *Node) *Node { return newNodes(ADD, expr1, expr2) }

// Go-only: language-adapter — IRSub ports IR.sub(Node, Node).
func IRSub(expr1, expr2 *Node) *Node { return newNodes(SUB, expr1, expr2) }

// Go-only: language-adapter — IRBitwiseAnd ports IR.bitwiseAnd(Node, Node).
func IRBitwiseAnd(expr1, expr2 *Node) *Node { return newNodes(BITAND, expr1, expr2) }

// Go-only: language-adapter — IRRightShift ports IR.rightShift(Node, Node).
func IRRightShift(expr1, expr2 *Node) *Node { return newNodes(RSH, expr1, expr2) }

// Go-only: language-adapter — IRObjectlit ports IR.objectlit(Node...).
func IRObjectlit(propdefs ...*Node) *Node {
	n := New(OBJECTLIT)
	for _, p := range propdefs {
		n.AddChildToBack(p)
	}
	return n
}

// Go-only: language-adapter — IRObjectPattern ports IR.objectPattern(Node...).
func IRObjectPattern(keys ...*Node) *Node {
	return newNodes(OBJECT_PATTERN, keys...)
}

// Go-only: language-adapter — IRArrayPattern ports IR.arrayPattern(Node...).
func IRArrayPattern(keys ...*Node) *Node {
	return newNodes(ARRAY_PATTERN, keys...)
}

// Go-only: language-adapter — IRComputedProp ports IR.computedProp(Node, Node).
func IRComputedProp(key, value *Node) *Node { return newNodes(COMPUTED_PROP, key, value) }

// Go-only: language-adapter — IRPropdef ports IR.propdef(Node, Node).
func IRPropdef(str, value *Node) *Node {
	str.AddChildToFront(value)
	return str
}

// Go-only: language-adapter — IRArraylit ports IR.arraylit(Node...).
func IRArraylit(exprs ...*Node) *Node {
	n := New(ARRAYLIT)
	for _, e := range exprs {
		n.AddChildToBack(e)
	}
	return n
}

// Go-only: language-adapter — IRRegexp ports IR.regexp(Node) / IR.regexp(Node, Node).
func IRRegexp(expr *Node, flags ...*Node) *Node {
	n := newNodes(REGEXP, expr)
	for _, f := range flags {
		n.AddChildToBack(f)
	}
	return n
}

// Go-only: language-adapter — IRString ports IR.string(String).
func IRString(s string) *Node { return NewString(STRINGLIT, s) }

// Go-only: language-adapter — IRStringKey ports IR.stringKey(String).
func IRStringKey(s string) *Node { return NewString(STRING_KEY, s) }

// Go-only: language-adapter — IRStringKeyValue ports IR.stringKey(String, Node).
// The two-argument overload is split into IRStringKey / IRStringKeyValue in Go.
func IRStringKeyValue(s string, value *Node) *Node {
	k := IRStringKey(s)
	k.AddChildToFront(value)
	return k
}

// Go-only: language-adapter — IRQuotedStringKey ports IR.quotedStringKey(String, Node).
func IRQuotedStringKey(s string, value *Node) *Node {
	k := IRStringKeyValue(s, value)
	k.PutBooleanProp(PropQuotedStringKey, true)
	return k
}

// Go-only: language-adapter — IRTemplateLiteral ports IR.templateLiteral().
func IRTemplateLiteral() *Node { return New(TEMPLATELIT) }

// Go-only: language-adapter — IRTemplateLiteralString ports IR.templateLiteralString(String, String).
func IRTemplateLiteralString(cooked, raw string) *Node {
	return NewTemplateLitString(cooked, raw)
}

// Go-only: language-adapter — IRTemplateLiteralSubstitution ports IR.templateLiteralSubstitution(Node).
func IRTemplateLiteralSubstitution(child *Node) *Node {
	n := New(TEMPLATELIT_SUB)
	n.AddChildToBack(child)
	return n
}

// Go-only: language-adapter — IRIterRest ports IR.iterRest(Node).
func IRIterRest(target *Node) *Node { return newNodes(ITER_REST, target) }

// Go-only: language-adapter — IRObjectRest ports IR.objectRest(Node).
func IRObjectRest(target *Node) *Node { return newNodes(OBJECT_REST, target) }

// Go-only: language-adapter — IRIterSpread ports IR.iterSpread(Node).
func IRIterSpread(expr *Node) *Node { return newNodes(ITER_SPREAD, expr) }

// Go-only: language-adapter — IRObjectSpread ports IR.objectSpread(Node).
func IRObjectSpread(expr *Node) *Node { return newNodes(OBJECT_SPREAD, expr) }

// Go-only: language-adapter — IRSuper ports IR.superNode().
func IRSuper() *Node { return New(SUPER) }

// Go-only: language-adapter — IRGetterDef ports IR.getterDef(String, Node).
func IRGetterDef(name string, value *Node) *Node {
	n := NewString(GETTER_DEF, name)
	n.AddChildToFront(value)
	return n
}

// Go-only: language-adapter — IRSetterDef ports IR.setterDef(String, Node).
func IRSetterDef(name string, value *Node) *Node {
	n := NewString(SETTER_DEF, name)
	n.AddChildToFront(value)
	return n
}

// Go-only: language-adapter — IRMemberFieldDef ports IR.memberFieldDef(String, Node).
func IRMemberFieldDef(name string, value *Node) *Node {
	n := NewString(MEMBER_FIELD_DEF, name)
	n.AddChildToFront(value)
	return n
}

// Go-only: language-adapter — IRMemberFunctionDef ports IR.memberFunctionDef(String, Node).
func IRMemberFunctionDef(name string, function *Node) *Node {
	n := NewString(MEMBER_FUNCTION_DEF, name)
	n.AddChildToBack(function)
	return n
}

// Go-only: language-adapter — IRNumber ports IR.number(double).
func IRNumber(d float64) *Node { return NewNumber(d) }

// Go-only: language-adapter — IRBigInt ports IR.bigint(BigInteger).
// The value is stored as a string in the node for portability.
func IRBigInt(val string) *Node { return NewString(BIGINT, val) }

// Go-only: language-adapter — IRThis ports IR.thisNode().
func IRThis() *Node { return New(THIS) }

// Go-only: language-adapter — IRTrue ports IR.trueNode().
func IRTrue() *Node { return New(TRUE) }

// Go-only: language-adapter — IRFalse ports IR.falseNode().
func IRFalse() *Node { return New(FALSE) }

// Go-only: language-adapter — IRNull ports IR.nullNode().
func IRNull() *Node { return New(NULL) }

// Go-only: language-adapter — IRTypeof ports IR.typeof(Node).
func IRTypeof(expr *Node) *Node { return newNodes(TYPEOF, expr) }

// Go-only: language-adapter — IRImportMeta ports IR.importMeta().
func IRImportMeta() *Node { return New(IMPORT_META) }

// Go-only: language-adapter — IRMayBeExpressionOrEmpty ports IR.mayBeExpressionOrEmpty(Node).
func IRMayBeExpressionOrEmpty(n *Node) bool {
	return n.IsEmpty() || IRMayBeExpression(n)
}

// Go-only: language-adapter — IRMayBeStatementNoReturn ports IR.mayBeStatementNoReturn(Node).
func IRMayBeStatementNoReturn(n *Node) bool {
	switch n.Token() {
	case EMPTY, FUNCTION,
		BLOCK, BREAK, CLASS, CONST, CONTINUE, DEBUGGER, DO,
		EXPR_RESULT, FOR, FOR_IN, FOR_OF, FOR_AWAIT_OF,
		IF, LABEL, LET, SWITCH, THROW, TRY, VAR, WHILE:
		return true
	}
	return false
}

// Go-only: language-adapter — IRMayBeStatement ports IR.mayBeStatement(Node).
func IRMayBeStatement(n *Node) bool {
	if !IRMayBeStatementNoReturn(n) {
		return n.IsReturn()
	}
	return true
}

// Go-only: language-adapter — IRMayBeExpression ports IR.mayBeExpression(Node).
func IRMayBeExpression(n *Node) bool {
	switch n.Token() {
	case FUNCTION, CLASS,
		ADD, AND, ARRAYLIT, ASSIGN,
		ASSIGN_BITOR, ASSIGN_BITXOR, ASSIGN_BITAND,
		ASSIGN_LSH, ASSIGN_RSH, ASSIGN_URSH,
		ASSIGN_ADD, ASSIGN_SUB, ASSIGN_MUL, ASSIGN_EXPONENT,
		ASSIGN_DIV, ASSIGN_MOD, ASSIGN_OR, ASSIGN_AND, ASSIGN_COALESCE,
		AWAIT, BIGINT, BITAND, BITOR, BITNOT, BITXOR,
		CALL, CAST, COALESCE, COMMA, DEC, DELPROP, DIV,
		DYNAMIC_IMPORT, EQ, EXPONENT, FALSE, GE,
		GETPROP, GETELEM, GT, HOOK, IMPORT_META,
		IN, INC, INSTANCEOF, LE, LSH, LT, MOD, MUL,
		NAME, NE, NEG, NEW, NEW_TARGET, NOT,
		NUMBER, NULL, OBJECTLIT,
		OPTCHAIN_CALL, OPTCHAIN_GETELEM, OPTCHAIN_GETPROP,
		OR, POS, REGEXP, RSH, SHEQ, SHNE, STRINGLIT,
		SUB, SUPER, TEMPLATELIT, TAGGED_TEMPLATELIT,
		THIS, TYPEOF, TRUE, URSH, VOID, YIELD:
		return true
	}
	return false
}
