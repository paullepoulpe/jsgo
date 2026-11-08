// Port of com/google/javascript/jscomp/CodeGenerator.java @ closure-compiler v20261104

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
	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// addStmtOrDecl handles statement, loop, function, class, module, and
// declaration nodes, returning true when node's token was handled.
func (g *CodeGenerator) addStmtOrDecl(node rhino.Node, context CodeGenerator_Context) bool {
	cc := &g.cc
	type_ := node.GetToken()
	switch type_ {
	case rhino.Token_TRY:
		g.addTry(ast.MustTry(node))
	case rhino.Token_CATCH:
		g.addCatch(ast.MustCatch(node))
	case rhino.Token_THROW:
		g.Add("throw")
		cc.MaybeInsertSpace()
		g.AddNode(ast.MustThrow(node).Value().Node())
		// Must have a ';' after a throw statement, otherwise safari can't
		// parse this.
		cc.EndStatementWithSemiColon(true /* needSemiColon */, g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_RETURN:
		g.addReturn(ast.MustReturn(node))
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET:
		g.addNameDecl(ast.MustNameDeclaration(node), context)
	case rhino.Token_PARAM_LIST:
		g.addParamList(ast.MustParamList(node))
	case rhino.Token_FUNCTION:
		g.addFunctionNode(ast.MustFunction(node), context)
	case rhino.Token_EXPORT:
		g.addExport(ast.MustExport(node), context)
	case rhino.Token_IMPORT:
		g.addImport(ast.MustImport(node))
	case rhino.Token_EXPORT_SPECS:
		firstSpec, _ := ast.MustExportSpecs(node).FirstSpec()
		g.addModuleSpecs(firstSpec.Node())
	case rhino.Token_IMPORT_SPECS:
		firstSpec, _ := ast.MustImportSpecs(node).FirstSpec()
		g.addModuleSpecs(firstSpec.Node())
	case rhino.Token_EXPORT_SPEC:
		g.addExportSpec(ast.MustExportSpec(node))
	case rhino.Token_IMPORT_SPEC:
		g.addImportSpec(ast.MustImportSpec(node))
	case rhino.Token_IMPORT_STAR:
		jutil.CheckStateMsg1(!node.HasChildren(), "%s", node)
		star := ast.MustImportStar(node)
		g.Add("*")
		g.Add("as")
		g.Add(star.Ident())
	case rhino.Token_CLASS:
		g.addClass(ast.MustClass(node), context)
	case rhino.Token_CLASS_MEMBERS:
		firstMember, _ := ast.MustClassBody(node).FirstMember()
		g.addMemberBlock(firstMember.Node(), context)
	case rhino.Token_INTERFACE_MEMBERS, rhino.Token_NAMESPACE_ELEMENTS:
		g.addMemberBlock(node.GetFirstChild(), context)
	case rhino.Token_ENUM_MEMBERS:
		g.addEnumMembers(node.GetFirstChild())
	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_VARIABLE_DEF:
		g.addMemberDef(node, type_)
	case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_FIELD_DEF:
		g.addFieldDef(ast.MustFieldDef(node), type_)
	case rhino.Token_SCRIPT, rhino.Token_MODULE_BODY, rhino.Token_BLOCK, rhino.Token_ROOT:
		g.addBlockOrRoot(node, type_, context)
	case rhino.Token_FOR:
		g.addFor(ast.MustFor(node), context)
	case rhino.Token_FOR_IN, rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		g.addForInOf(node, type_, context)
	case rhino.Token_DO:
		g.addDo(ast.MustDo(node))
	case rhino.Token_WHILE:
		g.addWhile(ast.MustWhile(node), context)
	case rhino.Token_EMPTY:
		jutil.CheckStateMsg1(!node.HasChildren(), "%s", node)
	case rhino.Token_WITH:
		withStmt := ast.MustWith(node)
		g.Add("with(")
		g.AddNode(withStmt.Object().Node())
		g.Add(")")
		g.AddNonEmptyStatement(withStmt.Body().Node(), CodeGenerator_GetContextForNonEmptyExpression(context), false)
	case rhino.Token_IF:
		g.addIf(ast.MustIf(node), context)
	case rhino.Token_CONTINUE, rhino.Token_BREAK:
		g.addJump(ast.MustJump(node), type_)
	case rhino.Token_DEBUGGER:
		jutil.CheckStateMsg1(!node.HasChildren(), "%s", node)
		g.Add("debugger")
		cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_EXPR_RESULT:
		g.AddNodeWithContext(ast.MustExprStmt(node).Expr().Node(), CodeGenerator_Context_START_OF_EXPR)
		cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_SWITCH:
		g.addSwitch(ast.MustSwitch(node), context)
	case rhino.Token_SWITCH_BODY:
		g.addSwitchBody(ast.MustSwitchBody(node), context)
	case rhino.Token_CASE:
		caseClause := ast.MustCase(node)
		g.Add("case ")
		g.AddNode(caseClause.Value().Node())
		g.AddCaseBody(caseClause.Body().Node())
	case rhino.Token_DEFAULT_CASE:
		g.Add("default")
		g.AddCaseBody(ast.MustDefaultCase(node).Body().Node())
	case rhino.Token_LABEL:
		g.addLabel(ast.MustLabel(node), context)
	case rhino.Token_INTERFACE:
		g.addInterface(ast.MustInterface(node))
	case rhino.Token_ENUM:
		en := ast.MustEnum(node)
		jutil.CheckStateMsg1(en.NumChildren() == 2, "%s", node)
		last, _ := en.LastChild()
		g.Add("enum")
		g.AddNode(en.MustFirstChild())
		g.AddNode(last)
	case rhino.Token_NAMESPACE:
		ns := ast.MustNamespace(node)
		jutil.CheckStateMsg1(ns.NumChildren() == 2, "%s", node)
		last, _ := ns.LastChild()
		g.Add("namespace")
		g.AddNode(ns.MustFirstChild())
		g.AddNode(last)
	case rhino.Token_TYPE_ALIAS:
		ta := ast.MustTypeAlias(node)
		last, _ := ta.LastChild()
		g.Add("type")
		g.Add(node.GetString())
		cc.AddOp("=", true)
		g.AddNode(last)
		cc.EndStatementWithSemiColon(true /* needSemiColon */, g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_DECLARE:
		g.Add("declare")
		g.AddNode(ast.MustDeclare(node).MustFirstChild())
		g.ProcessEnd(node, context)
	case rhino.Token_INDEX_SIGNATURE:
		g.Add("[")
		g.AddNode(ast.MustIndexSignature(node).MustFirstChild())
		g.Add("]")
		g.MaybeAddTypeDecl(node)
		cc.EndStatementWithSemiColon(true /* needSemiColon */, g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_CALL_SIGNATURE:
		if node.GetBooleanProp(rhino.Node_CONSTRUCT_SIGNATURE) {
			g.Add("new ")
		}
		g.MaybeAddGenericTypes(node)
		g.AddNode(ast.MustCallSignature(node).MustFirstChild())
		g.MaybeAddTypeDecl(node)
		cc.EndStatementWithSemiColon(true /* needSemiColon */, g.HasTrailingCommentOnSameLine(node))
	default:
		return false
	}
	return true
}

func (g *CodeGenerator) addTry(tryStmt ast.Try) {
	catchSlot := tryStmt.Catches().Node()
	jutil.CheckState(catchSlot.IsBlock() && !catchSlot.HasMoreThanOneChild())

	g.Add("try")
	g.AddNode(tryStmt.Body().Node())

	if catchClause, ok := tryStmt.Catch(); ok {
		g.AddNode(catchClause.Node())
	}

	if finallyBlock, ok := tryStmt.Finally(); ok {
		g.cc.MaybeInsertSpace()
		g.Add("finally")
		g.AddNode(finallyBlock.Node())
	}
}

func (g *CodeGenerator) addCatch(catchClause ast.Catch) {
	g.cc.MaybeInsertSpace()
	g.Add("catch")
	g.cc.MaybeInsertSpace()

	if param, ok := catchClause.Param(); ok {
		// optional catch binding
		g.Add("(")
		g.AddNode(param.Node())
		g.Add(")")
	}

	g.AddNode(catchClause.Body().Node())
}

func (g *CodeGenerator) addReturn(retStmt ast.Return) {
	g.Add("return")
	if val, ok := retStmt.Value(); ok {
		valNode := val.Node()
		g.cc.MaybeInsertSpace()
		if g.preserveTypeAnnotations && valNode.GetJSDocInfo() != nil {
			g.Add("(")
			g.AddNode(valNode)
			g.Add(")")
		} else {
			g.AddNode(valNode)
		}
	}
	g.cc.EndStatement(g.HasTrailingCommentOnSameLine(retStmt.Node()))
}

func (g *CodeGenerator) addNameDecl(decl ast.NameDeclaration, context CodeGenerator_Context) {
	node := decl.Node()
	switch decl.Token() {
	case rhino.Token_VAR:
		g.Add("var ")
	case rhino.Token_CONST:
		g.Add("const ")
	default:
		g.Add("let ")
	}
	if firstDecl, ok := decl.FirstDeclarator(); ok {
		g.AddListWithContextAndSeparator(firstDecl.Node(), false, CodeGenerator_GetContextForNoInOperator(context), ",")
	}
	if node.GetParent().IsNil() || NodeUtil_IsStatement(node) {
		g.cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
	}
}

func (g *CodeGenerator) addParamList(params ast.ParamList) {
	// If this is the list for a non-TypeScript arrow function with one
	// simple name param.
	parentFn, isFn := ast.AsFunction(params.Node().GetParent())
	if param, ok := params.OnlyParam(); ok &&
		isFn && parentFn.IsArrow() &&
		!g.gentsMode {
		if name, isName := param.AsName(); isName {
			g.AddNode(name.Node())
			return
		}
	}
	g.Add("(")
	if firstParam, ok := params.FirstParam(); ok {
		g.AddList(firstParam.Node())
	}
	g.Add(")")
}

func (g *CodeGenerator) addFunctionNode(fn ast.Function, context CodeGenerator_Context) {
	if !fn.Node().IsPlainNode() {
		panic(jutil.NewJavaLangError("Unexpected Node subclass."))
	}
	if fn.IsArrow() {
		g.addArrowFunction(fn, context)
	} else {
		g.addFunction(fn, context)
	}
}

func (g *CodeGenerator) addExport(exp ast.Export, context CodeGenerator_Context) {
	g.Add("export")
	if exp.IsDefault() {
		g.Add("default")
	}
	bodyNode := exp.Body().Node()
	if exp.IsAllFrom() {
		g.Add("*")
		jutil.CheckStateMsg1(!bodyNode.IsNil() && bodyNode.IsEmpty(), "%s", exp.Node())
	} else {
		g.AddNode(bodyNode)
	}
	if src, ok := exp.Source(); ok {
		g.Add("from")
		g.AddNode(src.Node())
	}
	g.ProcessEnd(bodyNode, context)
}

func (g *CodeGenerator) addImport(imp ast.Import) {
	g.Add("import")

	def, hasDef := imp.Default()
	bind, hasBind := imp.Bindings()
	if hasDef {
		g.AddNode(def.Node())
		if hasBind {
			g.cc.ListSeparator()
		}
	}
	if hasBind {
		g.AddNode(bind.Node())
	}
	if hasDef || hasBind {
		g.Add("from")
	}
	g.AddNode(imp.Source().Node())
	g.cc.EndStatement(g.HasTrailingCommentOnSameLine(imp.Node()))
}

func (g *CodeGenerator) addModuleSpecs(first rhino.Node) {
	g.Add("{")
	for c := first; !c.IsNil(); c = c.GetNext() {
		if c != first {
			g.cc.ListSeparator()
		}
		g.AddNode(c)
	}
	g.Add("}")
}

func (g *CodeGenerator) addExportSpec(spec ast.ExportSpec) {
	local := spec.Local()
	exported := spec.Exported()
	g.AddNode(local.Node())
	if !spec.Node().IsShorthandProperty() || local.Ident() != exported.Ident() {
		g.Add("as")
		g.AddNode(exported.Node())
	}
}

func (g *CodeGenerator) addImportSpec(spec ast.ImportSpec) {
	imported := spec.Imported()
	local := spec.Local()
	g.AddNode(imported.Node())
	if !spec.Node().IsShorthandProperty() || imported.Ident() != local.Ident() {
		g.Add("as")
		g.AddNode(local.Node())
	}
}

func (g *CodeGenerator) addClass(cls ast.Class, context CodeGenerator_Context) {
	classNeedsParens := context == CodeGenerator_Context_START_OF_EXPR
	if classNeedsParens {
		g.Add("(")
	}

	g.Add("class")
	if name, ok := cls.Name(); ok {
		g.AddNode(name.Node())
	}

	g.MaybeAddGenericTypes(cls.NameNode())

	if superClass, ok := cls.Extends(); ok {
		g.Add("extends")
		// ClassHeritage[Yield, Await] : extends LeftHandSideExpression[?Yield, ?Await]
		superNode := superClass.Node()
		minPrec := NodeUtil_Precedence(rhino.Token_CALL)
		if g.skipTransparentCasts(superNode).IsNew() {
			minPrec = NodeUtil_Precedence(rhino.Token_NEW)
		}
		g.AddExpr(superNode, minPrec, CodeGenerator_Context_OTHER)
	}

	interfaces := codeGenerator_nodeProp(cls.Node().GetProp(rhino.Node_IMPLEMENTS))
	if !interfaces.IsNil() {
		g.Add("implements")
		child := interfaces.GetFirstChild()
		g.AddNode(child)
		for child = child.GetNext(); !child.IsNil(); child = child.GetNext() {
			g.Add(",")
			g.cc.MaybeInsertSpace()
			g.AddNode(child)
		}
	}
	g.AddNode(cls.Body().Node())
	g.cc.EndClass(context == CodeGenerator_Context_STATEMENT)

	if classNeedsParens {
		g.Add(")")
	}
}

func (g *CodeGenerator) addMemberBlock(first rhino.Node, context CodeGenerator_Context) {
	cc := &g.cc
	cc.BeginBlock()
	for c := first; !c.IsNil(); c = c.GetNext() {
		g.AddNode(c)
		g.ProcessEnd(c, context)
		cc.EndLine()
	}
	cc.EndBlockWithShouldEndLine(false)
}

func (g *CodeGenerator) addEnumMembers(first rhino.Node) {
	cc := &g.cc
	cc.BeginBlock()
	for c := first; !c.IsNil(); c = c.GetNext() {
		g.AddNode(c)
		if !c.GetNext().IsNil() {
			g.Add(",")
		}
		cc.EndLine()
	}
	cc.EndBlockWithShouldEndLine(false)
}

func (g *CodeGenerator) addMemberDef(node rhino.Node, type_ rhino.Token) {
	parent := node.GetParent()
	jutil.CheckState(
		parent.IsObjectLit() ||
			parent.IsClassMembers() ||
			parent.IsInterfaceMembers() ||
			parent.IsRecordType() ||
			parent.IsIndexSignature())

	g.MaybeAddAccessibilityModifier(node)
	if node.IsStaticMember() {
		g.Add("static ")
	}

	if type_ == rhino.Token_MEMBER_VARIABLE_DEF {
		mv := ast.MustMemberVariableDef(node)
		g.Add(mv.Key())
		g.MaybeAddOptional(node)
		g.MaybeAddTypeDecl(node)
		return
	}

	method := ast.MustMethodDef(node)
	fn := method.Function()
	if type_ == rhino.Token_MEMBER_FUNCTION_DEF && fn.IsAsync() {
		g.Add("async ")
	}
	if fn.IsGenerator() {
		jutil.CheckStateMsg1(type_ == rhino.Token_MEMBER_FUNCTION_DEF, "%s", node)
		g.Add("*")
	}

	switch type_ {
	case rhino.Token_GETTER_DEF:
		jutil.CheckStateMsg1(!fn.Params().HasParams(), "%s", node)
		g.Add("get ")
	case rhino.Token_SETTER_DEF:
		_, hasOneParam := fn.Params().OnlyParam()
		jutil.CheckStateMsg1(hasOneParam, "%s", node)
		g.Add("set ")
	}

	// The function referenced by the definition should always be unnamed.
	fnNameNode := fn.NameNode()
	jutil.CheckStateMsg1(fnNameNode.Ident() == "", "%s", fn.Node())

	name := method.Key()
	if node.IsPrivateIdentifier() {
		// Never quote a private name (falling through to `addJsString` when `!isLatin(name)`
		// would turn `#π() {}` into a quoted public method `"#\u03c0"() {}`). Non-Latin
		// characters are Unicode-escaped via `addIdentifier`.
		g.AddIdentifier(name)
		g.MaybeAddGenericTypes(fnNameNode.Node())
	} else if !method.IsQuoted() &&
		rhino.TokenStream_IsJSIdentifier(name) &&
		NodeUtil_IsLatin(name) {
		g.Add(name)
		g.MaybeAddGenericTypes(fnNameNode.Node())
	} else {
		g.addNumericOrStringKey(name, node)
	}
	fnNode := fn.Node()
	g.MaybeAddOptional(fnNode)
	g.AddNode(fn.Params().Node())
	g.MaybeAddTypeDecl(fnNode)
	g.AddNode(fn.Body().Node())
}

func (g *CodeGenerator) addFieldDef(field ast.FieldDef, type_ rhino.Token) {
	node := field.Node()
	jutil.CheckState(node.GetParent().IsClassMembers())
	if field.IsStatic() {
		g.Add("static ")
	}
	if type_ == rhino.Token_MEMBER_FIELD_DEF {
		// Escape non-Latin characters as in other identifiers.
		g.AddIdentifier(ast.MustMemberFieldDef(node).Key())
	} else {
		g.Add("[")
		// Must use addExpr() with a priority of 1, because comma
		// expressions aren't allowed.
		g.AddExpr(ast.MustComputedFieldDef(node).Key().Node(), 1, CodeGenerator_Context_OTHER)
		g.Add("]")
	}
	if init, ok := field.Value(); ok {
		g.cc.AddOp("=", true)
		g.AddExpr(init.Node(), 1, CodeGenerator_Context_OTHER)
	}
	g.Add(";")
}

func firstBlockOrRootChild(node rhino.Node, type_ rhino.Token) rhino.Node {
	switch type_ {
	case rhino.Token_SCRIPT:
		first, _ := ast.MustScript(node).FirstItem()
		return first.Node()
	case rhino.Token_MODULE_BODY:
		first, _ := ast.MustModuleBody(node).FirstStatement()
		return first.Node()
	case rhino.Token_BLOCK:
		first, _ := ast.MustBlock(node).FirstStatement()
		return first.Node()
	default:
		first, _ := ast.MustRoot(node).FirstItem()
		return first.Node()
	}
}

func (g *CodeGenerator) addBlockOrRoot(node rhino.Node, type_ rhino.Token, context CodeGenerator_Context) {
	if !node.IsPlainNode() {
		panic(jutil.NewJavaLangError("Unexpected Node subclass."))
	}
	cc := &g.cc
	preserveBlock := false
	if blk, ok := ast.AsBlock(node); ok {
		if _, inClass := ast.AsClassBody(node.GetParent()); inClass {
			g.Add("static")
		}
		// A BLOCK marked as synthetic is not a real JS block with {}
		// around it in the JS code. It just represents a span of
		// statements that need to be kept together.
		preserveBlock = !blk.IsSynthetic()
	}
	if preserveBlock {
		cc.BeginBlock()
	}

	preferLineBreaks := type_ == rhino.Token_SCRIPT ||
		(type_ == rhino.Token_BLOCK && !preserveBlock && node.GetParent().IsScript())
	for c := firstBlockOrRootChild(node, type_); !c.IsNil(); c = c.GetNext() {
		g.AddNodeWithContext(c, CodeGenerator_Context_STATEMENT)

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
}

func (g *CodeGenerator) addFor(forStmt ast.For, context CodeGenerator_Context) {
	cc := &g.cc
	g.Add("for")
	cc.MaybeInsertSpace()
	g.Add("(")
	// EMPTY placeholder children must still pass through AddExpr/AddNode so
	// comments and source mappings attached to empty `for (;;)` slots are emitted.
	initNode := forStmt.Node().GetFirstChild()
	if _, isDecl := ast.AsNameDeclaration(initNode); isDecl {
		g.AddNodeWithContext(initNode, CodeGenerator_Context_IN_FOR_INIT_CLAUSE)
	} else {
		g.AddExpr(initNode, 0, CodeGenerator_Context_IN_FOR_INIT_CLAUSE)
	}
	g.Add(";")
	condNode := initNode.GetNext()
	if !condNode.IsEmpty() {
		cc.MaybeInsertSpace()
	}
	g.AddNode(condNode)
	g.Add(";")
	updateNode := condNode.GetNext()
	if !updateNode.IsEmpty() {
		cc.MaybeInsertSpace()
	}
	g.AddNode(updateNode)
	g.Add(")")
	g.AddNonEmptyStatement(forStmt.Body().Node(), CodeGenerator_GetContextForNonEmptyExpression(context), false)
}

func (g *CodeGenerator) addForInOf(node rhino.Node, type_ rhino.Token, context CodeGenerator_Context) {
	cc := &g.cc
	if type_ == rhino.Token_FOR_IN {
		forIn := ast.MustForIn(node)
		g.Add("for")
		cc.MaybeInsertSpace()
		g.Add("(")
		g.AddNode(forIn.Target().Node())
		g.Add("in")
		g.AddNode(forIn.Object().Node())
		g.Add(")")
		g.AddNonEmptyStatement(forIn.Body().Node(), CodeGenerator_GetContextForNonEmptyExpression(context), false)
		return
	}
	forOf := ast.MustForOfLoop(node)
	if type_ == rhino.Token_FOR_AWAIT_OF {
		g.Add("for await")
	} else {
		g.Add("for")
	}
	cc.MaybeInsertSpace()
	g.Add("(")
	g.AddNode(forOf.Target().Node())
	cc.MaybeInsertSpace()
	g.Add("of")
	cc.MaybeInsertSpace()
	// the iterable must be an AssignmentExpression
	g.AddExpr(forOf.Iterable().Node(), NodeUtil_Precedence(rhino.Token_ASSIGN), CodeGenerator_Context_OTHER)
	g.Add(")")
	g.AddNonEmptyStatement(forOf.AsLoop().Body().Node(), CodeGenerator_GetContextForNonEmptyExpression(context), false)
}

func (g *CodeGenerator) addDo(doStmt ast.Do) {
	cc := &g.cc
	g.Add("do")
	g.AddNonEmptyStatement(doStmt.Body().Node(), CodeGenerator_Context_OTHER, false)
	cc.MaybeInsertSpace()
	g.Add("while")
	cc.MaybeInsertSpace()
	g.Add("(")
	g.AddNode(doStmt.Cond().Node())
	g.Add(")")
	cc.EndStatement(g.HasTrailingCommentOnSameLine(doStmt.Node()))
}

func (g *CodeGenerator) addWhile(whileStmt ast.While, context CodeGenerator_Context) {
	g.Add("while")
	g.cc.MaybeInsertSpace()
	g.Add("(")
	g.AddNode(whileStmt.Cond().Node())
	g.Add(")")
	g.AddNonEmptyStatement(whileStmt.Body().Node(), CodeGenerator_GetContextForNonEmptyExpression(context), false)
}

func (g *CodeGenerator) addIf(ifStmt ast.If, context CodeGenerator_Context) {
	cc := &g.cc
	elseBlock, hasElse := ifStmt.Else()
	ambiguousElseClause := context == CodeGenerator_Context_BEFORE_DANGLING_ELSE && !hasElse
	if ambiguousElseClause {
		cc.BeginBlock()
	}

	g.Add("if")
	cc.MaybeInsertSpace()
	g.Add("(")
	g.AddNode(ifStmt.Cond().Node())
	g.Add(")")

	if hasElse {
		g.AddNonEmptyStatement(ifStmt.Then().Node(), CodeGenerator_Context_BEFORE_DANGLING_ELSE, false)
		cc.MaybeInsertSpace()
		g.Add("else")
		g.AddNonEmptyStatement(elseBlock.Node(), CodeGenerator_GetContextForNonEmptyExpression(context), false)
	} else {
		g.AddNonEmptyStatement(ifStmt.Then().Node(), CodeGenerator_Context_OTHER, false)
	}

	if ambiguousElseClause {
		cc.EndBlock()
	}
}

func (g *CodeGenerator) addJump(jmp ast.Jump, type_ rhino.Token) {
	if type_ == rhino.Token_CONTINUE {
		g.Add("continue")
	} else {
		g.Add("break")
	}
	if lbl, ok := jmp.Label(); ok {
		g.Add(" ")
		g.AddNode(lbl.Node())
	}
	g.cc.EndStatement(g.HasTrailingCommentOnSameLine(jmp.Node()))
}

func (g *CodeGenerator) addSwitch(sw ast.Switch, context CodeGenerator_Context) {
	g.Add("switch(")
	g.AddNode(sw.Discriminant().Node())
	g.Add(")")
	g.AddNodeWithContext(sw.Body().Node(), context)
}

func (g *CodeGenerator) addSwitchBody(swBody ast.SwitchBody, context CodeGenerator_Context) {
	g.cc.BeginBlock()
	if firstCase, ok := swBody.FirstCase(); ok {
		g.AddAllSiblings(firstCase.Node())
	}
	g.cc.EndBlockWithShouldEndLine(context == CodeGenerator_Context_STATEMENT)
}

func (g *CodeGenerator) addLabel(lbl ast.Label, context CodeGenerator_Context) {
	nameNode := lbl.Name().Node()
	if !nameNode.IsLabelName() {
		panic(jutil.NewJavaLangError("Unexpected token type. Should be LABEL_NAME."))
	}
	g.AddNode(nameNode)
	g.Add(":")
	bodyNode := lbl.Body().Node()
	if !bodyNode.IsBlock() {
		g.cc.MaybeInsertSpace()
	}
	g.AddNonEmptyStatement(bodyNode, CodeGenerator_GetContextForNonEmptyExpression(context), true)
}

func (g *CodeGenerator) addInterface(iface ast.Interface) {
	jutil.CheckStateMsg1(iface.NumChildren() == 3, "%s", iface.Node())
	name := iface.MustFirstChild()
	superTypes := name.GetNext()
	members := superTypes.GetNext()

	g.Add("interface")
	g.AddNode(name)
	g.MaybeAddGenericTypes(name)
	if !superTypes.IsEmpty() {
		g.Add("extends")
		superType := superTypes.GetFirstChild()
		g.AddNode(superType)
		for superType = superType.GetNext(); !superType.IsNil(); superType = superType.GetNext() {
			g.Add(",")
			g.cc.MaybeInsertSpace()
			g.AddNode(superType)
		}
	}
	g.AddNode(members)
}

// Java: CodeGenerator.arrowFunctionNeedsParens(Node)
func (g *CodeGenerator) ArrowFunctionNeedsParens(n rhino.Node) bool {
	parent := n.GetParent()
	expressionOrEnclosingCast := n
	for {
		if _, isCast := ast.AsCast(parent); !isCast {
			break
		}
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
		return true
	} else if getElem, ok := ast.AsNormalOrOptChainGetElem(parent); ok {
		return getElem.Object().Is(expressionOrEnclosingCast)
	} else if callOrNew, ok := ast.AsCallOrNew(parent); ok {
		return callOrNew.Callee().Is(expressionOrEnclosingCast)
	} else if hook, ok := ast.AsHook(parent); ok {
		return hook.Cond().Is(expressionOrEnclosingCast)
	}
	return false
}

func (g *CodeGenerator) addArrowFunction(fn ast.Function, context CodeGenerator_Context) {
	n := fn.Node()
	nameNode := fn.NameNode()
	jutil.CheckStateMsg1(nameNode.Ident() == "", "%s", nameNode.Node())
	funcNeedsParens := g.ArrowFunctionNeedsParens(n)
	if funcNeedsParens {
		g.Add("(")
	}

	g.MaybeAddGenericTypes(nameNode.Node())

	if fn.IsAsync() {
		g.Add("async")
	}
	g.AddNode(fn.Params().Node())
	g.MaybeAddTypeDecl(n)

	g.cc.AddOp("=>", true)

	bodyNode := fn.Body().Node()
	if bodyNode.IsBlock() {
		g.AddNode(bodyNode)
	} else {
		// Arrow functions have no token type, but blockless arrow function
		// bodies have lower precedence than anything other than commas.
		g.AddExpr(bodyNode, NodeUtil_Precedence(rhino.Token_COMMA)+1, CodeGenerator_GetContextForArrowFunctionBody(context))
	}
	g.cc.EndFunction(context == CodeGenerator_Context_STATEMENT)

	if funcNeedsParens {
		g.Add(")")
	}
}

func (g *CodeGenerator) addFunction(fn ast.Function, context CodeGenerator_Context) {
	n := fn.Node()
	funcNeedsParens := context == CodeGenerator_Context_START_OF_EXPR
	if funcNeedsParens {
		g.Add("(")
	}

	if fn.IsAsync() {
		g.Add("async function")
	} else {
		g.Add("function")
	}
	nameNode := fn.NameNode()
	if fn.IsGenerator() {
		g.Add("*")
		if nameNode.Ident() != "" {
			g.cc.MaybeInsertSpace()
		}
	}

	g.AddNode(nameNode.Node())
	g.MaybeAddGenericTypes(nameNode.Node())

	g.AddNode(fn.Params().Node())
	g.MaybeAddTypeDecl(n)

	g.AddNode(fn.Body().Node())
	g.cc.EndFunction(context == CodeGenerator_Context_STATEMENT)

	if funcNeedsParens {
		g.Add(")")
	}
}

// AddNonEmptyStatement adds a block or expression, substituting a VOID with
// an empty statement. This is used for "for (...);" and "if (...);" type
// statements.
//
// Java: CodeGenerator.addNonEmptyStatement(Node, Context, boolean)
func (g *CodeGenerator) AddNonEmptyStatement(n rhino.Node, context CodeGenerator_Context, allowNonBlockChild bool) {
	nodeToProcess := n

	_, isBlock := ast.AsBlock(n)
	if !allowNonBlockChild && !isBlock {
		panic(jutil.NewJavaLangError("Missing BLOCK child."))
	}

	// Strip unneeded blocks, that is blocks with <2 children unless
	// the CodePrinter specifically wants to keep them.
	if isBlock {
		count := CodeGenerator_GetNonEmptyChildCount(n, 2)
		if count == 0 {
			if g.cc.ShouldPreserveExtras(n) {
				g.cc.BeginBlock()
				g.PrintTrailingComment(n)
				g.cc.EndBlockWithShouldEndLine(g.cc.BreakAfterBlockFor(n, context == CodeGenerator_Context_STATEMENT))
			} else {
				g.PrintTrailingComment(n)
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
				g.AddNodeWithContext(firstAndOnlyChild, CodeGenerator_Context_STATEMENT)
				g.PrintTrailingComment(n)
				g.cc.MaybeLineBreak()
				g.cc.EndBlockWithShouldEndLine(g.cc.BreakAfterBlockFor(n, context == CodeGenerator_Context_STATEMENT))
				return
			}
			// Continue with the only child.
			nodeToProcess = firstAndOnlyChild
		}
	}

	if nodeToProcess.IsEmpty() {
		g.PrintTrailingComment(n)
		g.cc.EndStatementWithSemiColon(true /* needSemiColon */, false /* hasTrailingCommentOnSameLine */)
	} else {
		g.AddNodeWithContext(nodeToProcess, context)
		g.PrintTrailingComment(n)
	}
}

// CodeGenerator_IsBlockDeclOrDo returns whether the Node is a DO or a
// declaration that is only allowed in restricted contexts.
//
// Java: CodeGenerator.isBlockDeclOrDo(Node)
func CodeGenerator_IsBlockDeclOrDo(n rhino.Node) bool {
	if lbl, ok := ast.AsLabel(n); ok {
		labeledStatement := lbl.Body().Node()
		if !labeledStatement.IsBlock() {
			return CodeGenerator_IsBlockDeclOrDo(labeledStatement)
		}
		// For labels with block children, we need to ensure that a labeled
		// FUNCTION or DO isn't generated when extraneous BLOCKs are skipped.
		if CodeGenerator_GetNonEmptyChildCount(n, 2) == 1 {
			return CodeGenerator_IsBlockDeclOrDo(CodeGenerator_GetFirstNonEmptyChild(n))
		}
		return false
	}
	switch n.GetToken() {
	case rhino.Token_LET, rhino.Token_CONST, rhino.Token_FUNCTION, rhino.Token_CLASS, rhino.Token_DO:
		return true
	default:
		return false
	}
}

// CodeGenerator_GetNonEmptyChildCount returns the number of children of
// this node that are non empty up to maxCount.
//
// Java: CodeGenerator.getNonEmptyChildCount(Node, int)
func CodeGenerator_GetNonEmptyChildCount(n rhino.Node, maxCount int) int {
	i := 0
	for c := n.GetFirstChild(); !c.IsNil() && i < maxCount; c = c.GetNext() {
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

// Java: CodeGenerator.addCaseBody(Node)
func (g *CodeGenerator) AddCaseBody(caseBody rhino.Node) {
	blk := ast.MustBlock(caseBody)
	g.cc.BeginCaseBody()
	if firstStmt, ok := blk.FirstStatement(); ok {
		g.AddAllSiblings(firstStmt.Node())
	}
	g.cc.EndCaseBody()
}
