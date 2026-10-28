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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// addStmtOrDecl handles statement, loop, function, class, module, and
// declaration nodes, returning true when type_ was handled.
func (g *CodeGenerator) addStmtOrDecl(
	node rhino.Node,
	type_ rhino.Token,
	first rhino.Node,
	last rhino.Node,
	childCount int,
	context CodeGenerator_Context,
) bool {
	cc := &g.cc
	switch type_ {
	case rhino.Token_TRY:
		jutil.CheckState(first.GetNext().IsBlock() && !first.GetNext().HasMoreThanOneChild())
		jutil.CheckState(childCount >= 2 && childCount <= 3)

		tryStmt := ast.MustTry(node)
		g.Add("try")
		g.AddNode(tryStmt.Body().Node())

		if catchClause, ok := tryStmt.Catch(); ok {
			g.AddNode(catchClause.Node())
		}

		if finallyBlock, ok := tryStmt.Finally(); ok {
			cc.MaybeInsertSpace()
			g.Add("finally")
			g.AddNode(finallyBlock.Node())
		}
	case rhino.Token_CATCH:
		jutil.CheckStateMsg1(childCount == 2, "%s", node)
		catchClause := ast.MustCatch(node)
		cc.MaybeInsertSpace()
		g.Add("catch")
		cc.MaybeInsertSpace()

		if param, ok := catchClause.Param(); ok {
			// optional catch binding
			g.Add("(")
			g.AddNode(param.Node())
			g.Add(")")
		}

		g.AddNode(catchClause.Body().Node())
	case rhino.Token_THROW:
		jutil.CheckStateMsg1(childCount == 1, "%s", node)
		throwStmt := ast.MustThrow(node)
		g.Add("throw")
		cc.MaybeInsertSpace()
		g.AddNode(throwStmt.Value().Node())

		// Must have a ';' after a throw statement, otherwise safari can't
		// parse this.
		cc.EndStatementWithSemiColon(true /* needSemiColon */, g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_RETURN:
		jutil.CheckStateMsg1(childCount <= 1, "%s", node)
		retStmt := ast.MustReturn(node)
		g.Add("return")
		if val, ok := retStmt.Value(); ok {
			valNode := val.Node()
			cc.MaybeInsertSpace()
			if g.preserveTypeAnnotations && valNode.GetJSDocInfo() != nil {
				g.Add("(")
				g.AddNode(valNode)
				g.Add(")")
			} else {
				g.AddNode(valNode)
			}
		}
		cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_VAR, rhino.Token_CONST, rhino.Token_LET:
		decl := ast.MustNameDeclaration(node)
		switch type_ {
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
			cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
		}
	case rhino.Token_PARAM_LIST:
		params := ast.MustParamList(node)
		// If this is the list for a non-TypeScript arrow function with one
		// simple name param.
		if param, ok := params.OnlyParam(); ok &&
			node.GetParent().IsArrowFunction() &&
			param.Node().IsName() &&
			!g.gentsMode {
			g.AddNode(param.Node())
		} else {
			g.Add("(")
			if firstParam, ok := params.FirstParam(); ok {
				g.AddList(firstParam.Node())
			}
			g.Add(")")
		}
	case rhino.Token_FUNCTION:
		if !node.IsPlainNode() {
			panic(jutil.NewJavaLangError("Unexpected Node subclass."))
		}
		jutil.CheckStateMsg1(childCount == 3, "%s", node)
		fn := ast.MustFunction(node)
		if fn.IsArrow() {
			g.AddArrowFunction(node, first, last, context)
		} else {
			g.AddFunction(node, first, last, context)
		}
	case rhino.Token_EXPORT:
		exp := ast.MustExport(node)
		g.Add("export")
		if exp.IsDefault() {
			g.Add("default")
		}
		bodyNode := exp.Body().Node()
		if exp.IsAllFrom() {
			g.Add("*")
			jutil.CheckStateMsg1(!bodyNode.IsNil() && bodyNode.IsEmpty(), "%s", node)
		} else {
			g.AddNode(bodyNode)
		}
		if src, ok := exp.Source(); ok {
			g.Add("from")
			g.AddNode(src.Node())
		}
		g.ProcessEnd(bodyNode, context)
	case rhino.Token_IMPORT:
		imp := ast.MustImport(node)
		g.Add("import")

		def, hasDef := imp.Default()
		bind, hasBind := imp.Bindings()
		if hasDef {
			g.AddNode(def.Node())
			if hasBind {
				cc.ListSeparator()
			}
		}
		if hasBind {
			g.AddNode(bind.Node())
		}
		if hasDef || hasBind {
			g.Add("from")
		}
		g.AddNode(imp.Source().Node())
		cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_EXPORT_SPECS, rhino.Token_IMPORT_SPECS:
		g.Add("{")
		for c := first; !c.IsNil(); c = c.GetNext() {
			if c != first {
				cc.ListSeparator()
			}
			g.AddNode(c)
		}
		g.Add("}")
	case rhino.Token_EXPORT_SPEC:
		spec := ast.MustExportSpec(node)
		local := spec.Local()
		exported := spec.Exported()
		g.AddNode(local.Node())
		if !node.IsShorthandProperty() || local.Ident() != exported.Ident() {
			g.Add("as")
			g.AddNode(exported.Node())
		}
	case rhino.Token_IMPORT_SPEC:
		spec := ast.MustImportSpec(node)
		imported := spec.Imported()
		local := spec.Local()
		g.AddNode(imported.Node())
		if !node.IsShorthandProperty() || imported.Ident() != local.Ident() {
			g.Add("as")
			g.AddNode(local.Node())
		}
	case rhino.Token_IMPORT_STAR:
		star := ast.MustImportStar(node)
		g.Add("*")
		g.Add("as")
		g.Add(star.Ident())
	case rhino.Token_CLASS:
		jutil.CheckStateMsg1(childCount == 3, "%s", node)
		cls := ast.MustClass(node)
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
			// Parentheses are required for a comma expression or an
			// assignment expression.
			g.AddExpr(superClass.Node(), 1, CodeGenerator_Context_OTHER)
		}

		interfaces := codeGenerator_nodeProp(node.GetProp(rhino.Node_IMPLEMENTS))
		if !interfaces.IsNil() {
			g.Add("implements")
			child := interfaces.GetFirstChild()
			g.AddNode(child)
			for child = child.GetNext(); !child.IsNil(); child = child.GetNext() {
				g.Add(",")
				cc.MaybeInsertSpace()
				g.AddNode(child)
			}
		}
		g.AddNode(cls.Body().Node())
		cc.EndClass(context == CodeGenerator_Context_STATEMENT)

		if classNeedsParens {
			g.Add(")")
		}
	case rhino.Token_CLASS_MEMBERS, rhino.Token_INTERFACE_MEMBERS, rhino.Token_NAMESPACE_ELEMENTS:
		cc.BeginBlock()
		for c := first; !c.IsNil(); c = c.GetNext() {
			g.AddNode(c)
			g.ProcessEnd(c, context)
			cc.EndLine()
		}
		cc.EndBlockWithShouldEndLine(false)
	case rhino.Token_ENUM_MEMBERS:
		cc.BeginBlock()
		for c := first; !c.IsNil(); c = c.GetNext() {
			g.AddNode(c)
			if !c.GetNext().IsNil() {
				g.Add(",")
			}
			cc.EndLine()
		}
		cc.EndBlockWithShouldEndLine(false)
	case rhino.Token_GETTER_DEF, rhino.Token_SETTER_DEF, rhino.Token_MEMBER_FUNCTION_DEF, rhino.Token_MEMBER_VARIABLE_DEF:
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
			break
		}

		jutil.CheckStateMsg1(childCount == 1, "%s", node)
		jutil.CheckStateMsg1(first.IsFunction(), "%s", first)

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
			jutil.CheckStateMsg1(fn.Params().Node().HasOneChild(), "%s", node)
			g.Add("set ")
		}

		// The function referenced by the definition should always be unnamed.
		fnNameNode := fn.NameNode()
		jutil.CheckStateMsg1(fnNameNode.Ident() == "", "%s", first)

		name := method.Key()
		if !method.IsQuoted() &&
			(rhino.TokenStream_IsJSIdentifier(name) || node.IsPrivateIdentifier()) &&
			NodeUtil_IsLatin(name) {
			g.Add(name)
			g.MaybeAddGenericTypes(fnNameNode.Node())
		} else {
			d := CodeGenerator_GetSimpleNumber(name)
			if !math.IsNaN(d) {
				cc.AddNumber(d, node)
			} else {
				g.AddJsString(node)
			}
		}
		fnNode := fn.Node()
		g.MaybeAddOptional(fnNode)
		g.AddNode(fn.Params().Node())
		g.MaybeAddTypeDecl(fnNode)
		g.AddNode(fn.Body().Node())
	case rhino.Token_MEMBER_FIELD_DEF, rhino.Token_COMPUTED_FIELD_DEF:
		jutil.CheckState(node.GetParent().IsClassMembers())
		field := ast.MustFieldDef(node)
		if field.IsStatic() {
			g.Add("static ")
		}
		if type_ == rhino.Token_MEMBER_FIELD_DEF {
			g.Add(ast.MustMemberFieldDef(node).Key())
		} else {
			g.Add("[")
			// Must use addExpr() with a priority of 1, because comma
			// expressions aren't allowed.
			g.AddExpr(ast.MustComputedFieldDef(node).Key().Node(), 1, CodeGenerator_Context_OTHER)
			g.Add("]")
		}
		if init, ok := field.Value(); ok {
			cc.AddOp("=", true)
			g.AddExpr(init.Node(), 1, CodeGenerator_Context_OTHER)
		}
		g.Add(";")
	case rhino.Token_SCRIPT, rhino.Token_MODULE_BODY, rhino.Token_BLOCK, rhino.Token_ROOT:
		if !node.IsPlainNode() {
			panic(jutil.NewJavaLangError("Unexpected Node subclass."))
		}
		preserveBlock := false
		if blk, ok := ast.AsBlock(node); ok {
			if node.HasParent() && node.GetParent().IsClassMembers() {
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
		for c := first; !c.IsNil(); c = c.GetNext() {
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
	case rhino.Token_FOR:
		jutil.CheckStateMsg1(childCount == 4, "%s", node)
		forStmt := ast.MustFor(node)
		g.Add("for")
		cc.MaybeInsertSpace()
		g.Add("(")
		if NodeUtil_IsNameDeclaration(first) {
			g.AddNodeWithContext(first, CodeGenerator_Context_IN_FOR_INIT_CLAUSE)
		} else {
			g.AddExpr(first, 0, CodeGenerator_Context_IN_FOR_INIT_CLAUSE)
		}
		g.Add(";")
		condNode := first.GetNext()
		if _, hasCond := forStmt.Cond(); hasCond {
			cc.MaybeInsertSpace()
		}
		g.AddNode(condNode)
		g.Add(";")
		updateNode := condNode.GetNext()
		if _, hasUpdate := forStmt.Update(); hasUpdate {
			cc.MaybeInsertSpace()
		}
		g.AddNode(updateNode)
		g.Add(")")
		g.AddNonEmptyStatement(forStmt.Body().Node(), CodeGenerator_GetContextForNonEmptyExpression(context), false)
	case rhino.Token_FOR_IN:
		jutil.CheckStateMsg1(childCount == 3, "%s", node)
		forIn := ast.MustForIn(node)
		g.Add("for")
		cc.MaybeInsertSpace()
		g.Add("(")
		g.AddNode(forIn.Target().Node())
		g.Add("in")
		g.AddNode(forIn.Object().Node())
		g.Add(")")
		g.AddNonEmptyStatement(forIn.Body().Node(), CodeGenerator_GetContextForNonEmptyExpression(context), false)
	case rhino.Token_FOR_OF, rhino.Token_FOR_AWAIT_OF:
		jutil.CheckStateMsg1(childCount == 3, "%s", node)
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
	case rhino.Token_DO:
		jutil.CheckStateMsg1(childCount == 2, "%s", node)
		doStmt := ast.MustDo(node)
		g.Add("do")
		g.AddNonEmptyStatement(doStmt.Body().Node(), CodeGenerator_Context_OTHER, false)
		cc.MaybeInsertSpace()
		g.Add("while")
		cc.MaybeInsertSpace()
		g.Add("(")
		g.AddNode(doStmt.Cond().Node())
		g.Add(")")
		cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_WHILE:
		jutil.CheckStateMsg1(childCount == 2, "%s", node)
		whileStmt := ast.MustWhile(node)
		g.Add("while")
		cc.MaybeInsertSpace()
		g.Add("(")
		g.AddNode(whileStmt.Cond().Node())
		g.Add(")")
		g.AddNonEmptyStatement(whileStmt.Body().Node(), CodeGenerator_GetContextForNonEmptyExpression(context), false)
	case rhino.Token_EMPTY:
		jutil.CheckStateMsg1(childCount == 0, "%s", node)
	case rhino.Token_WITH:
		jutil.CheckStateMsg1(childCount == 2, "%s", node)
		withStmt := ast.MustWith(node)
		g.Add("with(")
		g.AddNode(withStmt.Object().Node())
		g.Add(")")
		g.AddNonEmptyStatement(withStmt.Body().Node(), CodeGenerator_GetContextForNonEmptyExpression(context), false)
	case rhino.Token_IF:
		jutil.CheckStateMsg1(childCount == 2 || childCount == 3, "%s", node)
		ifStmt := ast.MustIf(node)
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
	case rhino.Token_CONTINUE, rhino.Token_BREAK:
		jutil.CheckStateMsg1(childCount <= 1, "%s", node)
		if type_ == rhino.Token_CONTINUE {
			g.Add("continue")
		} else {
			g.Add("break")
		}
		if childCount == 1 {
			if !first.IsLabelName() {
				panic(jutil.NewJavaLangError("Unexpected token type. Should be LABEL_NAME."))
			}
			jmp := ast.MustJump(node)
			if lbl, ok := jmp.Label(); ok {
				g.Add(" ")
				g.AddNode(lbl.Node())
			}
		}
		cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_DEBUGGER:
		jutil.CheckStateMsg1(childCount == 0, "%s", node)
		g.Add("debugger")
		cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_EXPR_RESULT:
		jutil.CheckStateMsg1(childCount == 1, "%s", node)
		exprStmt := ast.MustExprStmt(node)
		g.AddNodeWithContext(exprStmt.Expr().Node(), CodeGenerator_Context_START_OF_EXPR)
		cc.EndStatement(g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_SWITCH:
		sw := ast.MustSwitch(node)
		g.Add("switch(")
		g.AddNode(sw.Discriminant().Node())
		g.Add(")")
		g.AddNodeWithContext(sw.Body().Node(), context)
	case rhino.Token_SWITCH_BODY:
		swBody := ast.MustSwitchBody(node)
		cc.BeginBlock()
		if firstCase, ok := swBody.FirstCase(); ok {
			g.AddAllSiblings(firstCase.Node())
		}
		cc.EndBlockWithShouldEndLine(context == CodeGenerator_Context_STATEMENT)
	case rhino.Token_CASE:
		jutil.CheckStateMsg1(childCount == 2, "%s", node)
		caseClause := ast.MustCase(node)
		g.Add("case ")
		g.AddNode(caseClause.Value().Node())
		g.AddCaseBody(caseClause.Body().Node())
	case rhino.Token_DEFAULT_CASE:
		jutil.CheckStateMsg1(childCount == 1, "%s", node)
		defClause := ast.MustDefaultCase(node)
		g.Add("default")
		g.AddCaseBody(defClause.Body().Node())
	case rhino.Token_LABEL:
		jutil.CheckStateMsg1(childCount == 2, "%s", node)
		if !first.IsLabelName() {
			panic(jutil.NewJavaLangError("Unexpected token type. Should be LABEL_NAME."))
		}
		lbl := ast.MustLabel(node)
		g.AddNode(lbl.Name().Node())
		g.Add(":")
		bodyNode := lbl.Body().Node()
		if !bodyNode.IsBlock() {
			cc.MaybeInsertSpace()
		}
		g.AddNonEmptyStatement(bodyNode, CodeGenerator_GetContextForNonEmptyExpression(context), true)
	case rhino.Token_INTERFACE:
		jutil.CheckStateMsg1(childCount == 3, "%s", node)
		name := first
		superTypes := first.GetNext()
		members := last

		g.Add("interface")
		g.AddNode(name)
		g.MaybeAddGenericTypes(name)
		if !superTypes.IsEmpty() {
			g.Add("extends")
			superType := superTypes.GetFirstChild()
			g.AddNode(superType)
			for superType = superType.GetNext(); !superType.IsNil(); superType = superType.GetNext() {
				g.Add(",")
				cc.MaybeInsertSpace()
				g.AddNode(superType)
			}
		}
		g.AddNode(members)
	case rhino.Token_ENUM:
		jutil.CheckStateMsg1(childCount == 2, "%s", node)
		g.Add("enum")
		g.AddNode(first)
		g.AddNode(last)
	case rhino.Token_NAMESPACE:
		jutil.CheckStateMsg1(childCount == 2, "%s", node)
		g.Add("namespace")
		g.AddNode(first)
		g.AddNode(last)
	case rhino.Token_TYPE_ALIAS:
		g.Add("type")
		g.Add(node.GetString())
		cc.AddOp("=", true)
		g.AddNode(last)
		cc.EndStatementWithSemiColon(true /* needSemiColon */, g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_DECLARE:
		g.Add("declare")
		g.AddNode(first)
		g.ProcessEnd(node, context)
	case rhino.Token_INDEX_SIGNATURE:
		g.Add("[")
		g.AddNode(first)
		g.Add("]")
		g.MaybeAddTypeDecl(node)
		cc.EndStatementWithSemiColon(true /* needSemiColon */, g.HasTrailingCommentOnSameLine(node))
	case rhino.Token_CALL_SIGNATURE:
		if node.GetBooleanProp(rhino.Node_CONSTRUCT_SIGNATURE) {
			g.Add("new ")
		}
		g.MaybeAddGenericTypes(node)
		g.AddNode(first)
		g.MaybeAddTypeDecl(node)
		cc.EndStatementWithSemiColon(true /* needSemiColon */, g.HasTrailingCommentOnSameLine(node))
	default:
		return false
	}
	return true
}

// Java: CodeGenerator.arrowFunctionNeedsParens(Node)
func (g *CodeGenerator) ArrowFunctionNeedsParens(n rhino.Node) bool {
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
	} else if parent.IsGetElem() ||
		parent.IsCall() ||
		parent.IsHook() ||
		parent.IsOptChainGetElem() ||
		parent.IsOptChainCall() ||
		parent.IsNew() {
		return expressionOrEnclosingCast.IsFirstChildOf(parent)
	}
	return false
}

// Java: CodeGenerator.addArrowFunction(Node, Node, Node, Context)
func (g *CodeGenerator) AddArrowFunction(n rhino.Node, first rhino.Node, last rhino.Node, context CodeGenerator_Context) {
	fn := ast.MustFunction(n)
	nameNode := fn.NameNode()
	jutil.CheckStateMsg1(nameNode.Ident() == "", "%s", first)
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

// Java: CodeGenerator.addFunction(Node, Node, Node, Context)
func (g *CodeGenerator) AddFunction(n rhino.Node, first rhino.Node, last rhino.Node, context CodeGenerator_Context) {
	fn := ast.MustFunction(n)
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
