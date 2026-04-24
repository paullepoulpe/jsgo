// Port of com/google/javascript/rhino/IR.java @ closure-compiler v20260415

package rhino

// IR is an AST construction helper. All methods are package-level functions
// because Java statics become free functions in Go. Mirrors
// com.google.javascript.rhino.IR.

// newNodes builds a Node with the given token and appends each child.
func newNodes(tok Token, children ...*Node) *Node {
	n := New(tok)
	for _, c := range children {
		n.AddChildToBack(c)
	}
	return n
}

// IREmpty returns a new EMPTY node. Mirrors IR.empty().
func IREmpty() *Node { return New(EMPTY) }

// IRExport wraps declaration in an EXPORT node. Mirrors IR.export(Node).
func IRExport(declaration *Node) *Node { return newNodes(EXPORT, declaration) }

// IRImportNode builds an IMPORT node. Mirrors IR.importNode(Node, Node, Node).
func IRImportNode(name, importSpecs, moduleIdentifier *Node) *Node {
	return newNodes(IMPORT, name, importSpecs, moduleIdentifier)
}

// IRImportStar creates an IMPORT_STAR node. Mirrors IR.importStar(String).
func IRImportStar(name string) *Node { return NewString(IMPORT_STAR, name) }

// IRFunction creates a FUNCTION node. Mirrors IR.function(Node, Node, Node).
func IRFunction(name, params, body *Node) *Node {
	return newNodes(FUNCTION, name, params, body)
}

// IRArrowFunction creates an arrow FUNCTION node. Mirrors IR.arrowFunction(Node, Node, Node).
func IRArrowFunction(name, params, body *Node) *Node {
	fn := newNodes(FUNCTION, name, params, body)
	fn.SetIsArrowFunction(true)
	return fn
}

// IRParamList creates a PARAM_LIST node. Mirrors IR.paramList(Node...).
func IRParamList(params ...*Node) *Node {
	return newNodes(PARAM_LIST, params...)
}

// IRRoot creates a ROOT node. Mirrors IR.root(Node...).
func IRRoot(rootChildren ...*Node) *Node {
	return newNodes(ROOT, rootChildren...)
}

// IRBlock creates an empty BLOCK node. Mirrors IR.block().
func IRBlock(stmts ...*Node) *Node {
	return newNodes(BLOCK, stmts...)
}

// IRScript creates an empty SCRIPT node. Mirrors IR.script().
func IRScript(stmts ...*Node) *Node {
	return newNodes(SCRIPT, stmts...)
}

// IRVar creates a VAR declaration node. Mirrors IR.var(Node, Node).
func IRVar(lhs, value *Node) *Node { return IRDeclaration(lhs, value, VAR) }

// IRVarLhs creates a VAR declaration without initialiser. Mirrors IR.var(Node).
func IRVarLhs(lhs *Node) *Node { return IRDeclarationLhs(lhs, VAR) }

// IRLet creates a LET declaration node. Mirrors IR.let(Node, Node).
func IRLet(lhs, value *Node) *Node { return IRDeclaration(lhs, value, LET) }

// IRLetLhs creates a LET declaration without initialiser. Mirrors IR.let(Node).
func IRLetLhs(lhs *Node) *Node { return IRDeclarationLhs(lhs, LET) }

// IRConst creates a CONST declaration node. Mirrors IR.constNode(Node, Node).
func IRConst(lhs, value *Node) *Node { return IRDeclaration(lhs, value, CONST) }

// IRDeclarationLhs creates a declaration without initialiser.
// Mirrors IR.declaration(Node, Token).
func IRDeclarationLhs(lhs *Node, typ Token) *Node {
	return newNodes(typ, lhs)
}

// IRDeclaration creates a declaration with initialiser.
// Mirrors IR.declaration(Node, Node, Token).
func IRDeclaration(lhs, value *Node, typ Token) *Node {
	lhs.AddChildToBack(value)
	return newNodes(typ, lhs)
}

// IRReturn creates a RETURN node. Mirrors IR.returnNode() / IR.returnNode(Node).
func IRReturn(exprs ...*Node) *Node {
	if len(exprs) == 0 {
		return New(RETURN)
	}
	return newNodes(RETURN, exprs[0])
}

// IRYield creates a YIELD node. Mirrors IR.yieldNode(Node).
func IRYield(expr *Node) *Node { return newNodes(YIELD, expr) }

// IRAwait creates an AWAIT node. Mirrors IR.await(Node).
func IRAwait(expr *Node) *Node { return newNodes(AWAIT, expr) }

// IRThrow creates a THROW node. Mirrors IR.throwNode(Node).
func IRThrow(expr *Node) *Node { return newNodes(THROW, expr) }

// IRExprResult creates an EXPR_RESULT node. Mirrors IR.exprResult(Node).
func IRExprResult(expr *Node) *Node { return newNodes(EXPR_RESULT, expr) }

// IRIf creates an IF node. Mirrors IR.ifNode(Node, Node) / IR.ifNode(Node, Node, Node).
func IRIf(cond, then *Node, elseNodes ...*Node) *Node {
	n := newNodes(IF, cond, then)
	for _, e := range elseNodes {
		n.AddChildToBack(e)
	}
	return n
}

// IRDo creates a DO node. Mirrors IR.doNode(Node, Node).
func IRDo(body, cond *Node) *Node { return newNodes(DO, body, cond) }

// IRWhile creates a WHILE node. Mirrors IR.whileNode(Node, Node).
func IRWhile(cond, body *Node) *Node { return newNodes(WHILE, cond, body) }

// IRForIn creates a FOR_IN node. Mirrors IR.forIn(Node, Node, Node).
func IRForIn(target, cond, body *Node) *Node { return newNodes(FOR_IN, target, cond, body) }

// IRFor creates a FOR node. Mirrors IR.forNode(Node, Node, Node, Node).
func IRFor(init, cond, incr, body *Node) *Node {
	n := newNodes(FOR, init, cond, incr)
	n.AddChildToBack(body)
	return n
}

// IRSwitch creates a SWITCH node. Mirrors IR.switchNode(Node, Node...).
func IRSwitch(cond *Node, cases ...*Node) *Node {
	n := newNodes(SWITCH, cond)
	body := New(SWITCH_BODY)
	n.AddChildToBack(body)
	for _, c := range cases {
		body.AddChildToBack(c)
	}
	return n
}

// IRCase creates a CASE node. Mirrors IR.caseNode(Node, Node).
func IRCase(expr, body *Node) *Node { return newNodes(CASE, expr, body) }

// IRDefaultCase creates a DEFAULT_CASE node. Mirrors IR.defaultCase(Node).
func IRDefaultCase(body *Node) *Node { return newNodes(DEFAULT_CASE, body) }

// IRLabel creates a LABEL node. Mirrors IR.label(Node, Node).
func IRLabel(name, stmt *Node) *Node { return newNodes(LABEL, name, stmt) }

// IRLabelName creates a LABEL_NAME node. Mirrors IR.labelName(String).
func IRLabelName(name string) *Node { return NewString(LABEL_NAME, name) }

// IRTryFinally creates a TRY node with a finally clause. Mirrors IR.tryFinally(Node, Node).
func IRTryFinally(tryBody, finallyBody *Node) *Node {
	catchBody := New(BLOCK)
	return newNodes(TRY, tryBody, catchBody, finallyBody)
}

// IRTryCatch creates a TRY node with a catch clause. Mirrors IR.tryCatch(Node, Node).
func IRTryCatch(tryBody, catchNode *Node) *Node {
	catchBody := newNodes(BLOCK, catchNode)
	return newNodes(TRY, tryBody, catchBody)
}

// IRTryCatchFinally creates a TRY node with both catch and finally.
// Mirrors IR.tryCatchFinally(Node, Node, Node).
func IRTryCatchFinally(tryBody, catchNode, finallyBody *Node) *Node {
	n := IRTryCatch(tryBody, catchNode)
	n.AddChildToBack(finallyBody)
	return n
}

// IRCatch creates a CATCH node. Mirrors IR.catchNode(Node, Node).
func IRCatch(expr, body *Node) *Node { return newNodes(CATCH, expr, body) }

// IRBreak creates a BREAK node. Mirrors IR.breakNode() / IR.breakNode(Node).
func IRBreak(labels ...*Node) *Node {
	if len(labels) == 0 {
		return New(BREAK)
	}
	return newNodes(BREAK, labels[0])
}

// IRContinue creates a CONTINUE node. Mirrors IR.continueNode() / IR.continueNode(Node).
func IRContinue(labels ...*Node) *Node {
	if len(labels) == 0 {
		return New(CONTINUE)
	}
	return newNodes(CONTINUE, labels[0])
}

// IRCall creates a CALL node. Mirrors IR.call(Node, Node...).
func IRCall(target *Node, args ...*Node) *Node {
	n := newNodes(CALL, target)
	for _, a := range args {
		n.AddChildToBack(a)
	}
	return n
}

// IRStartOptChainCall creates an optional-chaining CALL node (chain start).
// Mirrors IR.startOptChainCall(Node, Node...).
func IRStartOptChainCall(target *Node, args ...*Node) *Node {
	n := newNodes(OPTCHAIN_CALL, target)
	for _, a := range args {
		n.AddChildToBack(a)
	}
	n.SetIsOptionalChainStart(true)
	return n
}

// IRContinueOptChainCall creates an optional-chaining CALL node (chain cont).
// Mirrors IR.continueOptChainCall(Node, Node...).
func IRContinueOptChainCall(target *Node, args ...*Node) *Node {
	n := newNodes(OPTCHAIN_CALL, target)
	for _, a := range args {
		n.AddChildToBack(a)
	}
	n.SetIsOptionalChainStart(false)
	return n
}

// IRNew creates a NEW node. Mirrors IR.newNode(Node, Node...).
func IRNew(target *Node, args ...*Node) *Node {
	n := newNodes(NEW, target)
	for _, a := range args {
		n.AddChildToBack(a)
	}
	return n
}

// IRName creates a NAME node. Mirrors IR.name(String).
func IRName(name string) *Node { return NewString(NAME, name) }

// IRStartOptChainGetprop creates an optional-chaining GETPROP (start).
// Mirrors IR.startOptChainGetprop(Node, String).
func IRStartOptChainGetprop(target *Node, prop string) *Node {
	n := NewString(OPTCHAIN_GETPROP, prop)
	n.AddChildToBack(target)
	n.SetIsOptionalChainStart(true)
	return n
}

// IRContinueOptChainGetprop creates an optional-chaining GETPROP (cont).
// Mirrors IR.continueOptChainGetprop(Node, String).
func IRContinueOptChainGetprop(target *Node, prop string) *Node {
	n := NewString(OPTCHAIN_GETPROP, prop)
	n.AddChildToBack(target)
	n.SetIsOptionalChainStart(false)
	return n
}

// IRGetprop creates a GETPROP node. Mirrors IR.getprop(Node, String).
func IRGetprop(target *Node, prop string, moreProps ...string) *Node {
	n := NewString(GETPROP, prop)
	n.AddChildToBack(target)
	result := n
	for _, mp := range moreProps {
		result = IRGetprop(result, mp)
	}
	return result
}

// IRStartOptChainGetelem creates an OPTCHAIN_GETELEM (start).
// Mirrors IR.startOptChainGetelem(Node, Node).
func IRStartOptChainGetelem(target, elem *Node) *Node {
	n := newNodes(OPTCHAIN_GETELEM, target, elem)
	n.SetIsOptionalChainStart(true)
	return n
}

// IRContinueOptChainGetelem creates an OPTCHAIN_GETELEM (cont).
// Mirrors IR.continueOptChainGetelem(Node, Node).
func IRContinueOptChainGetelem(target, elem *Node) *Node {
	n := newNodes(OPTCHAIN_GETELEM, target, elem)
	n.SetIsOptionalChainStart(false)
	return n
}

// IRGetelem creates a GETELEM node. Mirrors IR.getelem(Node, Node).
func IRGetelem(target, elem *Node) *Node { return newNodes(GETELEM, target, elem) }

// IRDelprop creates a DELPROP node. Mirrors IR.delprop(Node).
func IRDelprop(target *Node) *Node { return newNodes(DELPROP, target) }

// IRAssign creates an ASSIGN node. Mirrors IR.assign(Node, Node).
func IRAssign(target, expr *Node) *Node { return newNodes(ASSIGN, target, expr) }

// IRAssignOr creates an ASSIGN_OR node. Mirrors IR.assignOr(Node, Node).
func IRAssignOr(expr1, expr2 *Node) *Node { return newNodes(ASSIGN_OR, expr1, expr2) }

// IRAssignAnd creates an ASSIGN_AND node. Mirrors IR.assignAnd(Node, Node).
func IRAssignAnd(expr1, expr2 *Node) *Node { return newNodes(ASSIGN_AND, expr1, expr2) }

// IRAssignCoalesce creates an ASSIGN_COALESCE node.
// Mirrors IR.assignCoalesce(Node, Node).
func IRAssignCoalesce(expr1, expr2 *Node) *Node { return newNodes(ASSIGN_COALESCE, expr1, expr2) }

// IRAnd creates an AND node. Mirrors IR.and(Node, Node).
func IRAnd(expr1, expr2 *Node) *Node { return newNodes(AND, expr1, expr2) }

// IROr creates an OR node. Mirrors IR.or(Node, Node).
func IROr(expr1, expr2 *Node) *Node { return newNodes(OR, expr1, expr2) }

// IRNot creates a NOT node. Mirrors IR.not(Node).
func IRNot(expr *Node) *Node { return newNodes(NOT, expr) }

// IREQ creates an EQ node. Mirrors IR.eq(Node, Node).
func IREQ(expr1, expr2 *Node) *Node { return newNodes(EQ, expr1, expr2) }

// IRNE creates a NE node. Mirrors IR.ne(Node, Node).
func IRNE(expr1, expr2 *Node) *Node { return newNodes(NE, expr1, expr2) }

// IRSHEQ creates a SHEQ node. Mirrors IR.sheq(Node, Node).
func IRSHEQ(expr1, expr2 *Node) *Node { return newNodes(SHEQ, expr1, expr2) }

// IRSHNE creates a SHNE node. Mirrors IR.shne(Node, Node).
func IRSHNE(expr1, expr2 *Node) *Node { return newNodes(SHNE, expr1, expr2) }

// IRVoid creates a VOID node. Mirrors IR.voidNode(Node).
func IRVoid(expr *Node) *Node { return newNodes(VOID, expr) }

// IRNeg creates a NEG node. Mirrors IR.neg(Node).
func IRNeg(expr *Node) *Node { return newNodes(NEG, expr) }

// IRPos creates a POS node. Mirrors IR.pos(Node).
func IRPos(expr *Node) *Node { return newNodes(POS, expr) }

// IRAdd creates an ADD node. Mirrors IR.add(Node, Node).
func IRAdd(expr1, expr2 *Node) *Node { return newNodes(ADD, expr1, expr2) }

// IRSub creates a SUB node. Mirrors IR.sub(Node, Node).
func IRSub(expr1, expr2 *Node) *Node { return newNodes(SUB, expr1, expr2) }

// IRBitwiseAnd creates a BITAND node. Mirrors IR.bitwiseAnd(Node, Node).
func IRBitwiseAnd(expr1, expr2 *Node) *Node { return newNodes(BITAND, expr1, expr2) }

// IRRightShift creates a RSH node. Mirrors IR.rightShift(Node, Node).
func IRRightShift(expr1, expr2 *Node) *Node { return newNodes(RSH, expr1, expr2) }

// IRObjectlit creates an OBJECTLIT node. Mirrors IR.objectlit(Node...).
func IRObjectlit(propdefs ...*Node) *Node {
	n := New(OBJECTLIT)
	for _, p := range propdefs {
		n.AddChildToBack(p)
	}
	return n
}

// IRObjectPattern creates an OBJECT_PATTERN node. Mirrors IR.objectPattern(Node...).
func IRObjectPattern(keys ...*Node) *Node {
	return newNodes(OBJECT_PATTERN, keys...)
}

// IRArrayPattern creates an ARRAY_PATTERN node. Mirrors IR.arrayPattern(Node...).
func IRArrayPattern(keys ...*Node) *Node {
	return newNodes(ARRAY_PATTERN, keys...)
}

// IRComputedProp creates a COMPUTED_PROP node. Mirrors IR.computedProp(Node, Node).
func IRComputedProp(key, value *Node) *Node { return newNodes(COMPUTED_PROP, key, value) }

// IRPropdef defines a string-key property by adding value to string as a child.
// Mirrors IR.propdef(Node, Node).
func IRPropdef(str, value *Node) *Node {
	str.AddChildToFront(value)
	return str
}

// IRArraylit creates an ARRAYLIT node. Mirrors IR.arraylit(Node...).
func IRArraylit(exprs ...*Node) *Node {
	n := New(ARRAYLIT)
	for _, e := range exprs {
		n.AddChildToBack(e)
	}
	return n
}

// IRRegexp creates a REGEXP node. Mirrors IR.regexp(Node) / IR.regexp(Node, Node).
func IRRegexp(expr *Node, flags ...*Node) *Node {
	n := newNodes(REGEXP, expr)
	for _, f := range flags {
		n.AddChildToBack(f)
	}
	return n
}

// IRString creates a STRINGLIT node. Mirrors IR.string(String).
func IRString(s string) *Node { return NewString(STRINGLIT, s) }

// IRStringKey creates a STRING_KEY node. Mirrors IR.stringKey(String).
func IRStringKey(s string) *Node { return NewString(STRING_KEY, s) }

// IRStringKeyValue creates a STRING_KEY with a child value.
// Mirrors IR.stringKey(String, Node).
func IRStringKeyValue(s string, value *Node) *Node {
	k := IRStringKey(s)
	k.AddChildToFront(value)
	return k
}

// IRQuotedStringKey creates a quoted STRING_KEY. Mirrors IR.quotedStringKey(String, Node).
func IRQuotedStringKey(s string, value *Node) *Node {
	k := IRStringKeyValue(s, value)
	k.PutBooleanProp(PropQuotedStringKey, true)
	return k
}

// IRTemplateLiteral creates a TEMPLATELIT node. Mirrors IR.templateLiteral().
func IRTemplateLiteral() *Node { return New(TEMPLATELIT) }

// IRTemplateLiteralString creates a TEMPLATELIT_STRING node.
// Mirrors IR.templateLiteralString(String, String).
func IRTemplateLiteralString(cooked, raw string) *Node {
	return NewTemplateLitString(cooked, raw)
}

// IRTemplateLiteralSubstitution creates a TEMPLATELIT_SUB node.
// Mirrors IR.templateLiteralSubstitution(Node).
func IRTemplateLiteralSubstitution(child *Node) *Node {
	n := New(TEMPLATELIT_SUB)
	n.AddChildToBack(child)
	return n
}

// IRIterRest creates an ITER_REST node. Mirrors IR.iterRest(Node).
func IRIterRest(target *Node) *Node { return newNodes(ITER_REST, target) }

// IRObjectRest creates an OBJECT_REST node. Mirrors IR.objectRest(Node).
func IRObjectRest(target *Node) *Node { return newNodes(OBJECT_REST, target) }

// IRIterSpread creates an ITER_SPREAD node. Mirrors IR.iterSpread(Node).
func IRIterSpread(expr *Node) *Node { return newNodes(ITER_SPREAD, expr) }

// IRObjectSpread creates an OBJECT_SPREAD node. Mirrors IR.objectSpread(Node).
func IRObjectSpread(expr *Node) *Node { return newNodes(OBJECT_SPREAD, expr) }

// IRSuper creates a SUPER node. Mirrors IR.superNode().
func IRSuper() *Node { return New(SUPER) }

// IRGetterDef creates a GETTER_DEF node. Mirrors IR.getterDef(String, Node).
func IRGetterDef(name string, value *Node) *Node {
	n := NewString(GETTER_DEF, name)
	n.AddChildToFront(value)
	return n
}

// IRSetterDef creates a SETTER_DEF node. Mirrors IR.setterDef(String, Node).
func IRSetterDef(name string, value *Node) *Node {
	n := NewString(SETTER_DEF, name)
	n.AddChildToFront(value)
	return n
}

// IRMemberFieldDef creates a MEMBER_FIELD_DEF node.
// Mirrors IR.memberFieldDef(String, Node).
func IRMemberFieldDef(name string, value *Node) *Node {
	n := NewString(MEMBER_FIELD_DEF, name)
	n.AddChildToFront(value)
	return n
}

// IRMemberFunctionDef creates a MEMBER_FUNCTION_DEF node.
// Mirrors IR.memberFunctionDef(String, Node).
func IRMemberFunctionDef(name string, function *Node) *Node {
	n := NewString(MEMBER_FUNCTION_DEF, name)
	n.AddChildToBack(function)
	return n
}

// IRNumber creates a NUMBER node. Mirrors IR.number(double).
func IRNumber(d float64) *Node { return NewNumber(d) }

// IRBigInt creates a BIGINT node. Mirrors IR.bigint(BigInteger).
// The value is stored as a string in the node for portability.
func IRBigInt(val string) *Node { return NewString(BIGINT, val) }

// IRThis creates a THIS node. Mirrors IR.thisNode().
func IRThis() *Node { return New(THIS) }

// IRTrue creates a TRUE node. Mirrors IR.trueNode().
func IRTrue() *Node { return New(TRUE) }

// IRFalse creates a FALSE node. Mirrors IR.falseNode().
func IRFalse() *Node { return New(FALSE) }

// IRNull creates a NULL node. Mirrors IR.nullNode().
func IRNull() *Node { return New(NULL) }

// IRTypeof creates a TYPEOF node. Mirrors IR.typeof(Node).
func IRTypeof(expr *Node) *Node { return newNodes(TYPEOF, expr) }

// IRImportMeta creates an IMPORT_META node. Mirrors IR.importMeta().
func IRImportMeta() *Node { return New(IMPORT_META) }

// IRMayBeStatement reports whether n can appear as a statement.
// Mirrors IR.mayBeStatement(Node).
func IRMayBeStatement(n *Node) bool {
	switch n.Token() {
	case EMPTY, FUNCTION,
		BLOCK, BREAK, CLASS, CONST, CONTINUE, DEBUGGER, DO,
		EXPR_RESULT, FOR, FOR_IN, FOR_OF, FOR_AWAIT_OF,
		IF, LABEL, LET, RETURN, SWITCH, THROW, TRY, VAR, WHILE:
		return true
	}
	return false
}

// IRMayBeExpression reports whether n can appear as an expression.
// Mirrors IR.mayBeExpression(Node).
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
