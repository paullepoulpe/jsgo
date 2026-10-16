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

// Go-only: language-adapter — private variadic builder that unifies Java's many IR constructor overloads into one Go helper.
// newNodes builds a Node with the given token and appends each child.
func newNodes(tok Token, children ...*Node) *Node {
	n := New(tok)
	for _, c := range children {
		n.AddChildToBack(c)
	}
	return n
}

// Go-only: language-adapter — Java calls IR.empty() as a static; Go uses a free function with IR prefix to avoid name collision in the rhino package.
// IREmpty returns a new EMPTY node. Mirrors IR.empty().
func IREmpty() *Node { return New(EMPTY) }

// Go-only: language-adapter — Java calls IR.export() as a static; Go uses IR prefix to avoid name collision.
// IRExport wraps declaration in an EXPORT node. Mirrors IR.export(Node).
func IRExport(declaration *Node) *Node { return newNodes(EXPORT, declaration) }

// Go-only: language-adapter — Java calls IR.importNode() as a static; Go uses IR prefix.
// IRImportNode builds an IMPORT node. Mirrors IR.importNode(Node, Node, Node).
func IRImportNode(name, importSpecs, moduleIdentifier *Node) *Node {
	return newNodes(IMPORT, name, importSpecs, moduleIdentifier)
}

// Go-only: language-adapter — Java calls IR.importStar() as a static; Go uses IR prefix.
// IRImportStar creates an IMPORT_STAR node. Mirrors IR.importStar(String).
func IRImportStar(name string) *Node { return NewString(IMPORT_STAR, name) }

// Go-only: language-adapter — Java calls IR.function() as a static; Go uses IR prefix.
// IRFunction creates a FUNCTION node. Mirrors IR.function(Node, Node, Node).
func IRFunction(name, params, body *Node) *Node {
	return newNodes(FUNCTION, name, params, body)
}

// Go-only: language-adapter — Java calls IR.arrowFunction() as a static; Go uses IR prefix.
// IRArrowFunction creates an arrow FUNCTION node. Mirrors IR.arrowFunction(Node, Node, Node).
func IRArrowFunction(name, params, body *Node) *Node {
	fn := newNodes(FUNCTION, name, params, body)
	fn.SetIsArrowFunction(true)
	return fn
}

// Go-only: language-adapter — Java calls IR.paramList() as a static; Go uses IR prefix and variadic for the Node... overload.
// IRParamList creates a PARAM_LIST node. Mirrors IR.paramList(Node...).
func IRParamList(params ...*Node) *Node {
	return newNodes(PARAM_LIST, params...)
}

// Go-only: language-adapter — Java calls IR.root() as a static; Go uses IR prefix and variadic.
// IRRoot creates a ROOT node. Mirrors IR.root(Node...).
func IRRoot(rootChildren ...*Node) *Node {
	return newNodes(ROOT, rootChildren...)
}

// Go-only: language-adapter — Java calls IR.block() as a static with multiple overloads; Go merges them via variadic.
// IRBlock creates an empty BLOCK node. Mirrors IR.block().
func IRBlock(stmts ...*Node) *Node {
	return newNodes(BLOCK, stmts...)
}

// Go-only: language-adapter — Java calls IR.script() as a static with multiple overloads; Go merges them via variadic.
// IRScript creates an empty SCRIPT node. Mirrors IR.script().
func IRScript(stmts ...*Node) *Node {
	return newNodes(SCRIPT, stmts...)
}

// Go-only: language-adapter — Java calls IR.var() as a static; Go uses IR prefix since "var" is a keyword.
// IRVar creates a VAR declaration node. Mirrors IR.var(Node, Node).
func IRVar(lhs, value *Node) *Node { return IRDeclaration(lhs, value, VAR) }

// Go-only: language-adapter — Java calls IR.var(Node) overload; Go splits overloads into separate named functions.
// IRVarLhs creates a VAR declaration without initialiser. Mirrors IR.var(Node).
func IRVarLhs(lhs *Node) *Node { return IRDeclarationLhs(lhs, VAR) }

// Go-only: language-adapter — Java calls IR.let() as a static; Go uses IR prefix since "let" could conflict.
// IRLet creates a LET declaration node. Mirrors IR.let(Node, Node).
func IRLet(lhs, value *Node) *Node { return IRDeclaration(lhs, value, LET) }

// Go-only: language-adapter — Java calls IR.let(Node) overload; Go splits overloads into separate named functions.
// IRLetLhs creates a LET declaration without initialiser. Mirrors IR.let(Node).
func IRLetLhs(lhs *Node) *Node { return IRDeclarationLhs(lhs, LET) }

// Go-only: language-adapter — Java calls IR.constNode() as a static; Go uses IR prefix.
// IRConst creates a CONST declaration node. Mirrors IR.constNode(Node, Node).
func IRConst(lhs, value *Node) *Node { return IRDeclaration(lhs, value, CONST) }

// Go-only: language-adapter — Java calls IR.declaration() as a static; Go uses IR prefix and splits overloads.
// IRDeclarationLhs creates a declaration without initialiser.
// Mirrors IR.declaration(Node, Token).
func IRDeclarationLhs(lhs *Node, typ Token) *Node {
	return newNodes(typ, lhs)
}

// Go-only: language-adapter — Java calls IR.declaration() as a static; Go uses IR prefix.
// IRDeclaration creates a declaration with initialiser.
// Mirrors IR.declaration(Node, Node, Token).
func IRDeclaration(lhs, value *Node, typ Token) *Node {
	lhs.AddChildToBack(value)
	return newNodes(typ, lhs)
}

// Go-only: language-adapter — Java calls IR.returnNode() with two overloads; Go merges them via variadic.
// IRReturn creates a RETURN node. Mirrors IR.returnNode() / IR.returnNode(Node).
func IRReturn(exprs ...*Node) *Node {
	if len(exprs) == 0 {
		return New(RETURN)
	}
	return newNodes(RETURN, exprs[0])
}

// Go-only: language-adapter — Java calls IR.yieldNode() as a static; Go uses IR prefix.
// IRYield creates a YIELD node. Mirrors IR.yieldNode(Node).
func IRYield(expr *Node) *Node { return newNodes(YIELD, expr) }

// Go-only: language-adapter — Java calls IR.await() as a static; Go uses IR prefix.
// IRAwait creates an AWAIT node. Mirrors IR.await(Node).
func IRAwait(expr *Node) *Node { return newNodes(AWAIT, expr) }

// Go-only: language-adapter — Java calls IR.throwNode() as a static; Go uses IR prefix.
// IRThrow creates a THROW node. Mirrors IR.throwNode(Node).
func IRThrow(expr *Node) *Node { return newNodes(THROW, expr) }

// Go-only: language-adapter — Java calls IR.exprResult() as a static; Go uses IR prefix.
// IRExprResult creates an EXPR_RESULT node. Mirrors IR.exprResult(Node).
func IRExprResult(expr *Node) *Node { return newNodes(EXPR_RESULT, expr) }

// Go-only: language-adapter — Java calls IR.ifNode() with two overloads; Go merges them via variadic.
// IRIf creates an IF node. Mirrors IR.ifNode(Node, Node) / IR.ifNode(Node, Node, Node).
func IRIf(cond, then *Node, elseNodes ...*Node) *Node {
	n := newNodes(IF, cond, then)
	for _, e := range elseNodes {
		n.AddChildToBack(e)
	}
	return n
}

// Go-only: language-adapter — Java calls IR.doNode() as a static; Go uses IR prefix.
// IRDo creates a DO node. Mirrors IR.doNode(Node, Node).
func IRDo(body, cond *Node) *Node { return newNodes(DO, body, cond) }

// Go-only: language-adapter — Java calls IR.whileNode() as a static; Go uses IR prefix.
// IRWhile creates a WHILE node. Mirrors IR.whileNode(Node, Node).
func IRWhile(cond, body *Node) *Node { return newNodes(WHILE, cond, body) }

// Go-only: language-adapter — Java calls IR.forIn() as a static; Go uses IR prefix.
// IRForIn creates a FOR_IN node. Mirrors IR.forIn(Node, Node, Node).
func IRForIn(target, cond, body *Node) *Node { return newNodes(FOR_IN, target, cond, body) }

// Go-only: language-adapter — Java calls IR.forNode() as a static; Go uses IR prefix.
// IRFor creates a FOR node. Mirrors IR.forNode(Node, Node, Node, Node).
func IRFor(init, cond, incr, body *Node) *Node {
	n := newNodes(FOR, init, cond, incr)
	n.AddChildToBack(body)
	return n
}

// Go-only: language-adapter — Java calls IR.switchNode() as a static with varargs; Go uses IR prefix and variadic.
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

// Go-only: language-adapter — Java calls IR.caseNode() as a static; Go uses IR prefix.
// IRCase creates a CASE node. Mirrors IR.caseNode(Node, Node).
func IRCase(expr, body *Node) *Node { return newNodes(CASE, expr, body) }

// Go-only: language-adapter — Java calls IR.defaultCase() as a static; Go uses IR prefix.
// IRDefaultCase creates a DEFAULT_CASE node. Mirrors IR.defaultCase(Node).
func IRDefaultCase(body *Node) *Node { return newNodes(DEFAULT_CASE, body) }

// Go-only: language-adapter — Java calls IR.label() as a static; Go uses IR prefix.
// IRLabel creates a LABEL node. Mirrors IR.label(Node, Node).
func IRLabel(name, stmt *Node) *Node { return newNodes(LABEL, name, stmt) }

// Go-only: language-adapter — Java calls IR.labelName() as a static; Go uses IR prefix.
// IRLabelName creates a LABEL_NAME node. Mirrors IR.labelName(String).
func IRLabelName(name string) *Node { return NewString(LABEL_NAME, name) }

// Go-only: language-adapter — Java calls IR.tryFinally() as a static; Go uses IR prefix.
// IRTryFinally creates a TRY node with a finally clause. Mirrors IR.tryFinally(Node, Node).
func IRTryFinally(tryBody, finallyBody *Node) *Node {
	catchBody := New(BLOCK)
	return newNodes(TRY, tryBody, catchBody, finallyBody)
}

// Go-only: language-adapter — Java calls IR.tryCatch() as a static; Go uses IR prefix.
// IRTryCatch creates a TRY node with a catch clause. Mirrors IR.tryCatch(Node, Node).
func IRTryCatch(tryBody, catchNode *Node) *Node {
	catchBody := newNodes(BLOCK, catchNode)
	return newNodes(TRY, tryBody, catchBody)
}

// Go-only: language-adapter — Java calls IR.tryCatchFinally() as a static; Go uses IR prefix.
// IRTryCatchFinally creates a TRY node with both catch and finally.
// Mirrors IR.tryCatchFinally(Node, Node, Node).
func IRTryCatchFinally(tryBody, catchNode, finallyBody *Node) *Node {
	n := IRTryCatch(tryBody, catchNode)
	n.AddChildToBack(finallyBody)
	return n
}

// Go-only: language-adapter — Java calls IR.catchNode() as a static; Go uses IR prefix.
// IRCatch creates a CATCH node. Mirrors IR.catchNode(Node, Node).
func IRCatch(expr, body *Node) *Node { return newNodes(CATCH, expr, body) }

// Go-only: language-adapter — Java calls IR.breakNode() with two overloads; Go merges them via variadic.
// IRBreak creates a BREAK node. Mirrors IR.breakNode() / IR.breakNode(Node).
func IRBreak(labels ...*Node) *Node {
	if len(labels) == 0 {
		return New(BREAK)
	}
	return newNodes(BREAK, labels[0])
}

// Go-only: language-adapter — Java calls IR.continueNode() with two overloads; Go merges them via variadic.
// IRContinue creates a CONTINUE node. Mirrors IR.continueNode() / IR.continueNode(Node).
func IRContinue(labels ...*Node) *Node {
	if len(labels) == 0 {
		return New(CONTINUE)
	}
	return newNodes(CONTINUE, labels[0])
}

// Go-only: language-adapter — Java calls IR.call() as a static with varargs; Go uses IR prefix and variadic.
// IRCall creates a CALL node. Mirrors IR.call(Node, Node...).
func IRCall(target *Node, args ...*Node) *Node {
	n := newNodes(CALL, target)
	for _, a := range args {
		n.AddChildToBack(a)
	}
	return n
}

// Go-only: language-adapter — Java calls IR.startOptChainCall() as a static with varargs; Go uses IR prefix and variadic.
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

// Go-only: language-adapter — Java calls IR.continueOptChainCall() as a static with varargs; Go uses IR prefix and variadic.
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

// Go-only: language-adapter — Java calls IR.newNode() as a static with varargs; Go uses IR prefix and variadic.
// IRNew creates a NEW node. Mirrors IR.newNode(Node, Node...).
func IRNew(target *Node, args ...*Node) *Node {
	n := newNodes(NEW, target)
	for _, a := range args {
		n.AddChildToBack(a)
	}
	return n
}

// Go-only: language-adapter — Java calls IR.name() as a static; Go uses IR prefix.
// IRName creates a NAME node. Mirrors IR.name(String).
func IRName(name string) *Node { return NewString(NAME, name) }

// Go-only: language-adapter — Java calls IR.startOptChainGetprop() as a static; Go uses IR prefix.
// IRStartOptChainGetprop creates an optional-chaining GETPROP (start).
// Mirrors IR.startOptChainGetprop(Node, String).
func IRStartOptChainGetprop(target *Node, prop string) *Node {
	n := NewString(OPTCHAIN_GETPROP, prop)
	n.AddChildToBack(target)
	n.SetIsOptionalChainStart(true)
	return n
}

// Go-only: language-adapter — Java calls IR.continueOptChainGetprop() as a static; Go uses IR prefix.
// IRContinueOptChainGetprop creates an optional-chaining GETPROP (cont).
// Mirrors IR.continueOptChainGetprop(Node, String).
func IRContinueOptChainGetprop(target *Node, prop string) *Node {
	n := NewString(OPTCHAIN_GETPROP, prop)
	n.AddChildToBack(target)
	n.SetIsOptionalChainStart(false)
	return n
}

// Go-only: language-adapter — Java calls IR.getprop() as a static; Go uses IR prefix and adds variadic for chaining multiple props.
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

// Go-only: language-adapter — Java calls IR.startOptChainGetelem() as a static; Go uses IR prefix.
// IRStartOptChainGetelem creates an OPTCHAIN_GETELEM (start).
// Mirrors IR.startOptChainGetelem(Node, Node).
func IRStartOptChainGetelem(target, elem *Node) *Node {
	n := newNodes(OPTCHAIN_GETELEM, target, elem)
	n.SetIsOptionalChainStart(true)
	return n
}

// Go-only: language-adapter — Java calls IR.continueOptChainGetelem() as a static; Go uses IR prefix.
// IRContinueOptChainGetelem creates an OPTCHAIN_GETELEM (cont).
// Mirrors IR.continueOptChainGetelem(Node, Node).
func IRContinueOptChainGetelem(target, elem *Node) *Node {
	n := newNodes(OPTCHAIN_GETELEM, target, elem)
	n.SetIsOptionalChainStart(false)
	return n
}

// Go-only: language-adapter — Java calls IR.getelem() as a static; Go uses IR prefix.
// IRGetelem creates a GETELEM node. Mirrors IR.getelem(Node, Node).
func IRGetelem(target, elem *Node) *Node { return newNodes(GETELEM, target, elem) }

// Go-only: language-adapter — Java calls IR.delprop() as a static; Go uses IR prefix.
// IRDelprop creates a DELPROP node. Mirrors IR.delprop(Node).
func IRDelprop(target *Node) *Node { return newNodes(DELPROP, target) }

// Go-only: language-adapter — Java calls IR.assign() as a static; Go uses IR prefix.
// IRAssign creates an ASSIGN node. Mirrors IR.assign(Node, Node).
func IRAssign(target, expr *Node) *Node { return newNodes(ASSIGN, target, expr) }

// Go-only: language-adapter — Java calls IR.assignOr() as a static; Go uses IR prefix.
// IRAssignOr creates an ASSIGN_OR node. Mirrors IR.assignOr(Node, Node).
func IRAssignOr(expr1, expr2 *Node) *Node { return newNodes(ASSIGN_OR, expr1, expr2) }

// Go-only: language-adapter — Java calls IR.assignAnd() as a static; Go uses IR prefix.
// IRAssignAnd creates an ASSIGN_AND node. Mirrors IR.assignAnd(Node, Node).
func IRAssignAnd(expr1, expr2 *Node) *Node { return newNodes(ASSIGN_AND, expr1, expr2) }

// Go-only: language-adapter — Java calls IR.assignCoalesce() as a static; Go uses IR prefix.
// IRAssignCoalesce creates an ASSIGN_COALESCE node.
// Mirrors IR.assignCoalesce(Node, Node).
func IRAssignCoalesce(expr1, expr2 *Node) *Node { return newNodes(ASSIGN_COALESCE, expr1, expr2) }

// Go-only: language-adapter — Java calls IR.and() as a static; Go uses IR prefix.
// IRAnd creates an AND node. Mirrors IR.and(Node, Node).
func IRAnd(expr1, expr2 *Node) *Node { return newNodes(AND, expr1, expr2) }

// Go-only: language-adapter — Java calls IR.or() as a static; Go uses IR prefix.
// IROr creates an OR node. Mirrors IR.or(Node, Node).
func IROr(expr1, expr2 *Node) *Node { return newNodes(OR, expr1, expr2) }

// Go-only: language-adapter — Java calls IR.not() as a static; Go uses IR prefix.
// IRNot creates a NOT node. Mirrors IR.not(Node).
func IRNot(expr *Node) *Node { return newNodes(NOT, expr) }

// Go-only: language-adapter — Java calls IR.eq() as a static; Go uses IREQ prefix (eq conflicts with built-in semantics and EQ token).
// IREQ creates an EQ node. Mirrors IR.eq(Node, Node).
func IREQ(expr1, expr2 *Node) *Node { return newNodes(EQ, expr1, expr2) }

// Go-only: language-adapter — Java calls IR.ne() as a static; Go uses IRNE prefix.
// IRNE creates a NE node. Mirrors IR.ne(Node, Node).
func IRNE(expr1, expr2 *Node) *Node { return newNodes(NE, expr1, expr2) }

// Go-only: language-adapter — Java calls IR.sheq() as a static; Go uses IRSHEQ prefix.
// IRSHEQ creates a SHEQ node. Mirrors IR.sheq(Node, Node).
func IRSHEQ(expr1, expr2 *Node) *Node { return newNodes(SHEQ, expr1, expr2) }

// Go-only: language-adapter — Java calls IR.shne() as a static; Go uses IRSHNE prefix.
// IRSHNE creates a SHNE node. Mirrors IR.shne(Node, Node).
func IRSHNE(expr1, expr2 *Node) *Node { return newNodes(SHNE, expr1, expr2) }

// Go-only: language-adapter — Java calls IR.voidNode() as a static; Go uses IR prefix.
// IRVoid creates a VOID node. Mirrors IR.voidNode(Node).
func IRVoid(expr *Node) *Node { return newNodes(VOID, expr) }

// Go-only: language-adapter — Java calls IR.neg() as a static; Go uses IR prefix.
// IRNeg creates a NEG node. Mirrors IR.neg(Node).
func IRNeg(expr *Node) *Node { return newNodes(NEG, expr) }

// Go-only: language-adapter — Java calls IR.pos() as a static; Go uses IR prefix.
// IRPos creates a POS node. Mirrors IR.pos(Node).
func IRPos(expr *Node) *Node { return newNodes(POS, expr) }

// Go-only: language-adapter — Java calls IR.add() as a static; Go uses IR prefix.
// IRAdd creates an ADD node. Mirrors IR.add(Node, Node).
func IRAdd(expr1, expr2 *Node) *Node { return newNodes(ADD, expr1, expr2) }

// Go-only: language-adapter — Java calls IR.sub() as a static; Go uses IR prefix.
// IRSub creates a SUB node. Mirrors IR.sub(Node, Node).
func IRSub(expr1, expr2 *Node) *Node { return newNodes(SUB, expr1, expr2) }

// Go-only: language-adapter — Java calls IR.bitwiseAnd() as a static; Go uses IR prefix.
// IRBitwiseAnd creates a BITAND node. Mirrors IR.bitwiseAnd(Node, Node).
func IRBitwiseAnd(expr1, expr2 *Node) *Node { return newNodes(BITAND, expr1, expr2) }

// Go-only: language-adapter — Java calls IR.rightShift() as a static; Go uses IR prefix.
// IRRightShift creates a RSH node. Mirrors IR.rightShift(Node, Node).
func IRRightShift(expr1, expr2 *Node) *Node { return newNodes(RSH, expr1, expr2) }

// Go-only: language-adapter — Java calls IR.objectlit() as a static with varargs; Go uses IR prefix and variadic.
// IRObjectlit creates an OBJECTLIT node. Mirrors IR.objectlit(Node...).
func IRObjectlit(propdefs ...*Node) *Node {
	n := New(OBJECTLIT)
	for _, p := range propdefs {
		n.AddChildToBack(p)
	}
	return n
}

// Go-only: language-adapter — Java calls IR.objectPattern() as a static with varargs; Go uses IR prefix and variadic.
// IRObjectPattern creates an OBJECT_PATTERN node. Mirrors IR.objectPattern(Node...).
func IRObjectPattern(keys ...*Node) *Node {
	return newNodes(OBJECT_PATTERN, keys...)
}

// Go-only: language-adapter — Java calls IR.arrayPattern() as a static with varargs; Go uses IR prefix and variadic.
// IRArrayPattern creates an ARRAY_PATTERN node. Mirrors IR.arrayPattern(Node...).
func IRArrayPattern(keys ...*Node) *Node {
	return newNodes(ARRAY_PATTERN, keys...)
}

// Go-only: language-adapter — Java calls IR.computedProp() as a static; Go uses IR prefix.
// IRComputedProp creates a COMPUTED_PROP node. Mirrors IR.computedProp(Node, Node).
func IRComputedProp(key, value *Node) *Node { return newNodes(COMPUTED_PROP, key, value) }

// Go-only: language-adapter — Java calls IR.propdef() as a static; Go uses IR prefix.
// IRPropdef defines a string-key property by adding value to string as a child.
// Mirrors IR.propdef(Node, Node).
func IRPropdef(str, value *Node) *Node {
	str.AddChildToFront(value)
	return str
}

// Go-only: language-adapter — Java calls IR.arraylit() as a static with varargs; Go uses IR prefix and variadic.
// IRArraylit creates an ARRAYLIT node. Mirrors IR.arraylit(Node...).
func IRArraylit(exprs ...*Node) *Node {
	n := New(ARRAYLIT)
	for _, e := range exprs {
		n.AddChildToBack(e)
	}
	return n
}

// Go-only: language-adapter — Java calls IR.regexp() with two overloads; Go merges them via variadic.
// IRRegexp creates a REGEXP node. Mirrors IR.regexp(Node) / IR.regexp(Node, Node).
func IRRegexp(expr *Node, flags ...*Node) *Node {
	n := newNodes(REGEXP, expr)
	for _, f := range flags {
		n.AddChildToBack(f)
	}
	return n
}

// Go-only: language-adapter — Java calls IR.string() as a static; Go uses IR prefix since "string" is a built-in type name.
// IRString creates a STRINGLIT node. Mirrors IR.string(String).
func IRString(s string) *Node { return NewString(STRINGLIT, s) }

// Go-only: language-adapter — Java calls IR.stringKey() as a static; Go uses IR prefix.
// IRStringKey creates a STRING_KEY node. Mirrors IR.stringKey(String).
func IRStringKey(s string) *Node { return NewString(STRING_KEY, s) }

// Go-only: language-adapter — Java calls IR.stringKey(String, Node) overload; Go splits into separate named function.
// IRStringKeyValue creates a STRING_KEY with a child value.
// Mirrors IR.stringKey(String, Node).
func IRStringKeyValue(s string, value *Node) *Node {
	k := IRStringKey(s)
	k.AddChildToFront(value)
	return k
}

// Go-only: language-adapter — Java calls IR.quotedStringKey() as a static; Go uses IR prefix.
// IRQuotedStringKey creates a quoted STRING_KEY. Mirrors IR.quotedStringKey(String, Node).
func IRQuotedStringKey(s string, value *Node) *Node {
	k := IRStringKeyValue(s, value)
	k.PutBooleanProp(PropQuotedStringKey, true)
	return k
}

// Go-only: language-adapter — Java calls IR.templateLiteral() as a static; Go uses IR prefix.
// IRTemplateLiteral creates a TEMPLATELIT node. Mirrors IR.templateLiteral().
func IRTemplateLiteral() *Node { return New(TEMPLATELIT) }

// Go-only: language-adapter — Java calls IR.templateLiteralString() as a static; Go uses IR prefix.
// IRTemplateLiteralString creates a TEMPLATELIT_STRING node.
// Mirrors IR.templateLiteralString(String, String).
func IRTemplateLiteralString(cooked, raw string) *Node {
	return NewTemplateLitString(cooked, raw)
}

// Go-only: language-adapter — Java calls IR.templateLiteralSubstitution() as a static; Go uses IR prefix.
// IRTemplateLiteralSubstitution creates a TEMPLATELIT_SUB node.
// Mirrors IR.templateLiteralSubstitution(Node).
func IRTemplateLiteralSubstitution(child *Node) *Node {
	n := New(TEMPLATELIT_SUB)
	n.AddChildToBack(child)
	return n
}

// Go-only: language-adapter — Java calls IR.iterRest() as a static; Go uses IR prefix.
// IRIterRest creates an ITER_REST node. Mirrors IR.iterRest(Node).
func IRIterRest(target *Node) *Node { return newNodes(ITER_REST, target) }

// Go-only: language-adapter — Java calls IR.objectRest() as a static; Go uses IR prefix.
// IRObjectRest creates an OBJECT_REST node. Mirrors IR.objectRest(Node).
func IRObjectRest(target *Node) *Node { return newNodes(OBJECT_REST, target) }

// Go-only: language-adapter — Java calls IR.iterSpread() as a static; Go uses IR prefix.
// IRIterSpread creates an ITER_SPREAD node. Mirrors IR.iterSpread(Node).
func IRIterSpread(expr *Node) *Node { return newNodes(ITER_SPREAD, expr) }

// Go-only: language-adapter — Java calls IR.objectSpread() as a static; Go uses IR prefix.
// IRObjectSpread creates an OBJECT_SPREAD node. Mirrors IR.objectSpread(Node).
func IRObjectSpread(expr *Node) *Node { return newNodes(OBJECT_SPREAD, expr) }

// Go-only: language-adapter — Java calls IR.superNode() as a static; Go uses IR prefix.
// IRSuper creates a SUPER node. Mirrors IR.superNode().
func IRSuper() *Node { return New(SUPER) }

// Go-only: language-adapter — Java calls IR.getterDef() as a static; Go uses IR prefix.
// IRGetterDef creates a GETTER_DEF node. Mirrors IR.getterDef(String, Node).
func IRGetterDef(name string, value *Node) *Node {
	n := NewString(GETTER_DEF, name)
	n.AddChildToFront(value)
	return n
}

// Go-only: language-adapter — Java calls IR.setterDef() as a static; Go uses IR prefix.
// IRSetterDef creates a SETTER_DEF node. Mirrors IR.setterDef(String, Node).
func IRSetterDef(name string, value *Node) *Node {
	n := NewString(SETTER_DEF, name)
	n.AddChildToFront(value)
	return n
}

// Go-only: language-adapter — Java calls IR.memberFieldDef() as a static; Go uses IR prefix.
// IRMemberFieldDef creates a MEMBER_FIELD_DEF node.
// Mirrors IR.memberFieldDef(String, Node).
func IRMemberFieldDef(name string, value *Node) *Node {
	n := NewString(MEMBER_FIELD_DEF, name)
	n.AddChildToFront(value)
	return n
}

// Go-only: language-adapter — Java calls IR.memberFunctionDef() as a static; Go uses IR prefix.
// IRMemberFunctionDef creates a MEMBER_FUNCTION_DEF node.
// Mirrors IR.memberFunctionDef(String, Node).
func IRMemberFunctionDef(name string, function *Node) *Node {
	n := NewString(MEMBER_FUNCTION_DEF, name)
	n.AddChildToBack(function)
	return n
}

// Go-only: language-adapter — Java calls IR.number() as a static; Go uses IR prefix since "number" conflicts with type semantics.
// IRNumber creates a NUMBER node. Mirrors IR.number(double).
func IRNumber(d float64) *Node { return NewNumber(d) }

// Go-only: language-adapter — Java's IR.bigint() takes BigInteger; Go uses string since there's no BigInteger type.
// IRBigInt creates a BIGINT node. Mirrors IR.bigint(BigInteger).
// The value is stored as a string in the node for portability.
func IRBigInt(val string) *Node { return NewString(BIGINT, val) }

// Go-only: language-adapter — Java calls IR.thisNode() as a static; Go uses IR prefix.
// IRThis creates a THIS node. Mirrors IR.thisNode().
func IRThis() *Node { return New(THIS) }

// Go-only: language-adapter — Java calls IR.trueNode() as a static; Go uses IR prefix since "true" is a keyword.
// IRTrue creates a TRUE node. Mirrors IR.trueNode().
func IRTrue() *Node { return New(TRUE) }

// Go-only: language-adapter — Java calls IR.falseNode() as a static; Go uses IR prefix since "false" is a keyword.
// IRFalse creates a FALSE node. Mirrors IR.falseNode().
func IRFalse() *Node { return New(FALSE) }

// Go-only: language-adapter — Java calls IR.nullNode() as a static; Go uses IR prefix since "null" is not a Go keyword.
// IRNull creates a NULL node. Mirrors IR.nullNode().
func IRNull() *Node { return New(NULL) }

// Go-only: language-adapter — Java calls IR.typeof() as a static; Go uses IR prefix since "typeof" is a JS keyword not a Go name.
// IRTypeof creates a TYPEOF node. Mirrors IR.typeof(Node).
func IRTypeof(expr *Node) *Node { return newNodes(TYPEOF, expr) }

// Go-only: language-adapter — Java calls IR.importMeta() as a static; Go uses IR prefix.
// IRImportMeta creates an IMPORT_META node. Mirrors IR.importMeta().
func IRImportMeta() *Node { return New(IMPORT_META) }

// Go-only: language-adapter — Java calls IR.mayBeStatement() as a static; Go uses IR prefix.
// IRMayBeStatement reports whether n can appear as a statement.
// Mirrors IR.mayBeStatement(Node).
func IRMayBeStatement(n *Node) bool {
	switch n.Token() {
	case EMPTY, FUNCTION,
		BLOCK, BREAK, CLASS, CONST, CONTINUE, DEBUGGER, DO,
		EXPR_RESULT, FOR, FOR_IN, FOR_OF, FOR_AWAIT_OF,
		IF, LABEL, LET, SWITCH, THROW, TRY, VAR, WHILE:
		return true
	}
	return false
}

// Go-only: language-adapter — Java calls IR.mayBeExpression() as a static; Go uses IR prefix.
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
