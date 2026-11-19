package j2go;

import com.sun.source.tree.AssertTree;
import com.sun.source.tree.AssignmentTree;
import com.sun.source.tree.BinaryTree;
import com.sun.source.tree.BindingPatternTree;
import com.sun.source.tree.BlockTree;
import com.sun.source.tree.BreakTree;
import com.sun.source.tree.CaseLabelTree;
import com.sun.source.tree.CaseTree;
import com.sun.source.tree.CatchTree;
import com.sun.source.tree.ClassTree;
import com.sun.source.tree.CompoundAssignmentTree;
import com.sun.source.tree.ConditionalExpressionTree;
import com.sun.source.tree.ConstantCaseLabelTree;
import com.sun.source.tree.ContinueTree;
import com.sun.source.tree.DefaultCaseLabelTree;
import com.sun.source.tree.DoWhileLoopTree;
import com.sun.source.tree.EnhancedForLoopTree;
import com.sun.source.tree.ExpressionStatementTree;
import com.sun.source.tree.ExpressionTree;
import com.sun.source.tree.ForLoopTree;
import com.sun.source.tree.IfTree;
import com.sun.source.tree.IdentifierTree;
import com.sun.source.tree.InstanceOfTree;
import com.sun.source.tree.LabeledStatementTree;
import com.sun.source.tree.LambdaExpressionTree;
import com.sun.source.tree.LiteralTree;
import com.sun.source.tree.MemberSelectTree;
import com.sun.source.tree.MethodInvocationTree;
import com.sun.source.tree.MethodTree;
import com.sun.source.tree.NewClassTree;
import com.sun.source.tree.ParenthesizedTree;
import com.sun.source.tree.PatternCaseLabelTree;
import com.sun.source.tree.PatternTree;
import com.sun.source.tree.ReturnTree;
import com.sun.source.tree.StatementTree;
import com.sun.source.tree.SwitchExpressionTree;
import com.sun.source.tree.SwitchTree;
import com.sun.source.tree.SynchronizedTree;
import com.sun.source.tree.ThrowTree;
import com.sun.source.tree.Tree;
import com.sun.source.tree.TryTree;
import com.sun.source.tree.UnaryTree;
import com.sun.source.tree.UnionTypeTree;
import com.sun.source.tree.VariableTree;
import com.sun.source.tree.WhileLoopTree;
import com.sun.source.tree.YieldTree;
import com.sun.source.util.TreePath;
import com.sun.source.util.TreePathScanner;
import com.sun.source.util.TreeScanner;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import javax.lang.model.element.Element;
import javax.lang.model.element.ElementKind;
import javax.lang.model.element.ExecutableElement;
import javax.lang.model.element.Modifier;
import javax.lang.model.element.Name;
import javax.lang.model.element.NestingKind;
import javax.lang.model.element.TypeElement;
import javax.lang.model.element.VariableElement;
import javax.lang.model.type.ArrayType;
import javax.lang.model.type.DeclaredType;
import javax.lang.model.type.ExecutableType;
import javax.lang.model.type.TypeKind;
import javax.lang.model.type.TypeMirror;

/**
 * Translates Java AST {@link StatementTree} and {@link SwitchExpressionTree} nodes into idiomatic Go statements.
 */
public final class StmtEmitter {

  private final TypeResolver resolver;
  private final ExprEmitter exprEmitter;

  public StmtEmitter(TypeResolver resolver, ExprEmitter exprEmitter) {
    this.resolver = resolver;
    this.exprEmitter = exprEmitter;
  }

  public String emitBlockStatements(
      BlockTree block, TreePath blockPath, int indent, EmitContext ctx) {
    StringBuilder sb = new StringBuilder();
    List<? extends StatementTree> stmts = block.getStatements();
    for (int i = 0; i < stmts.size(); i++) {
      StatementTree stmt = stmts.get(i);
      TreePath stmtPath = new TreePath(blockPath, stmt);
      emitPatternVarPredecls(stmtPath, indent, ctx.predeclaredPatternVars, sb, ctx);
      if (i + 1 < stmts.size()
          && stmt instanceof VariableTree vt
          && stmts.get(i + 1) instanceof SwitchTree swt) {
        TreePath nextPath = new TreePath(blockPath, swt);
        if (tryEmitAdjacentNodeTokenVarSwitch(
            vt, stmtPath, swt, nextPath, stmts, i + 2, blockPath, indent, sb, ctx)) {
          i++;
          continue;
        }
      }
      emitStmt(stmt, stmtPath, indent, null, sb, ctx);
    }
    return sb.toString();
  }

  private boolean tryEmitAdjacentNodeTokenVarSwitch(
      VariableTree vt,
      TreePath vtPath,
      SwitchTree swt,
      TreePath swtPath,
      List<? extends StatementTree> blockStmts,
      int afterSwitchIdx,
      TreePath blockPath,
      int indent,
      StringBuilder out,
      EmitContext ctx) {
    if (vt.getInitializer() == null || swt.getCases().isEmpty()) {
      return false;
    }
    if (!(resolver.trees.getElement(vtPath) instanceof VariableElement ve)
        || ve.getKind() != ElementKind.LOCAL_VARIABLE) {
      return false;
    }
    if (!(ve.asType() instanceof DeclaredType vdt)
        || !(vdt.asElement() instanceof TypeElement vte)
        || !vte.getQualifiedName().contentEquals("com.google.javascript.rhino.Token")) {
      return false;
    }
    TreePath selPath = unwrapParensPath(new TreePath(swtPath, swt.getExpression()));
    if (!(selPath.getLeaf() instanceof IdentifierTree)
        || !ve.equals(resolver.trees.getElement(selPath))) {
      return false;
    }

    Set<VariableElement> tokenConsts = new HashSet<>();
    for (Element enc : vte.getEnclosedElements()) {
      if (enc instanceof VariableElement eve && eve.getKind() == ElementKind.ENUM_CONSTANT) {
        tokenConsts.add(eve);
      }
    }
    for (CaseTree ct : swt.getCases()) {
      if (ct.getGuard() != null) {
        return false;
      }
      for (CaseLabelTree lbl : ct.getLabels()) {
        if (lbl instanceof DefaultCaseLabelTree) {
          continue;
        }
        if (lbl instanceof ConstantCaseLabelTree ccl) {
          ExpressionTree ce = ccl.getConstantExpression();
          if (ce.getKind() != Tree.Kind.NULL_LITERAL) {
            TreePath cePath = new TreePath(new TreePath(new TreePath(swtPath, ct), lbl), ce);
            if (resolver.trees.getElement(cePath) instanceof VariableElement cve
                && tokenConsts.contains(cve)) {
              continue;
            }
          }
        }
        return false;
      }
    }

    TreePath initPath = new TreePath(vtPath, vt.getInitializer());
    if (referencesVariable(initPath, ve)) {
      return false;
    }
    for (CaseTree ct : swt.getCases()) {
      if (referencesVariable(new TreePath(swtPath, ct), ve)) {
        return false;
      }
    }
    for (int j = afterSwitchIdx; j < blockStmts.size(); j++) {
      if (referencesVariable(new TreePath(blockPath, blockStmts.get(j)), ve)) {
        return false;
      }
    }

    String nodeTokenRecv = exprEmitter.tryEmitNodeForTokenRead(initPath, ctx);
    if (nodeTokenRecv == null) {
      return false;
    }

    emitPatternVarPredecls(swtPath, indent, ctx.predeclaredPatternVars, out, ctx);
    String tGo = ctx.registerLocalVar(ve);
    tabs(out, indent)
        .append("var ")
        .append(tGo)
        .append(" int32 = ")
        .append(nodeTokenRecv)
        .append(".TokenSwitchOrd_()\n");
    emitSwitchCommon(
        swt.getExpression(), swt.getCases(), swtPath, indent, false, null, tGo, out, ctx);
    return true;
  }

  public void emitStmt(
      StatementTree stmt,
      TreePath path,
      int indent,
      String pendingLoopLabel,
      StringBuilder out,
      EmitContext ctx) {
    if (stmt == null) {
      return;
    }
    if (!(stmt instanceof BlockTree)) {
      emitPatternVarPredecls(path, indent, ctx.predeclaredPatternVars, out, ctx);
    }

    switch (stmt.getKind()) {
      case EMPTY_STATEMENT -> {}
      case BLOCK -> {
        tabs(out, indent).append("{\n");
        out.append(emitBlockStatements((BlockTree) stmt, path, indent + 1, ctx));
        tabs(out, indent).append("}\n");
      }
      case VARIABLE -> emitVariableDecl((VariableTree) stmt, path, indent, out, ctx);
      case CLASS, RECORD, ENUM, INTERFACE -> {
        ClassTree ct = (ClassTree) stmt;
        TypeElement te = resolver.getClassTreeElement(ct);
        if (te == null && resolver.trees.getElement(path) instanceof TypeElement t) {
          te = t;
        }
        ctx.deferLocalOrAnonClass(ct, path, te);
      }
      case EXPRESSION_STATEMENT -> {
        ExpressionTree expr = ((ExpressionStatementTree) stmt).getExpression();
        if (tryEmitStackNewClassImmediateCallStmt(expr, path, indent, out, ctx)) {
          break;
        }
        if (tryEmitUnwrappedConditionalAssignStmt(expr, path, indent, out, ctx)) {
          break;
        }
        tabs(out, indent).append(emitExprAsStmt(expr, path, ctx)).append('\n');
      }
      case IF -> emitIf((IfTree) stmt, path, indent, false, out, ctx);
      case WHILE_LOOP -> emitWhile((WhileLoopTree) stmt, path, indent, pendingLoopLabel, out, ctx);
      case DO_WHILE_LOOP ->
          emitDoWhile((DoWhileLoopTree) stmt, path, indent, pendingLoopLabel, out, ctx);
      case FOR_LOOP -> emitFor((ForLoopTree) stmt, path, indent, pendingLoopLabel, out, ctx);
      case ENHANCED_FOR_LOOP ->
          emitEnhancedFor((EnhancedForLoopTree) stmt, path, indent, pendingLoopLabel, out, ctx);
      case LABELED_STATEMENT -> emitLabeled((LabeledStatementTree) stmt, path, indent, out, ctx);
      case BREAK -> emitBreak((BreakTree) stmt, indent, out, ctx);
      case CONTINUE -> emitContinue((ContinueTree) stmt, indent, out, ctx);
      case RETURN -> emitReturn((ReturnTree) stmt, path, indent, out, ctx);
      case YIELD -> emitYield((YieldTree) stmt, path, indent, out, ctx);
      case THROW -> emitThrow((ThrowTree) stmt, path, indent, out, ctx);
      case ASSERT -> emitAssert((AssertTree) stmt, path, indent, out, ctx);
      case SYNCHRONIZED -> emitSynchronized((SynchronizedTree) stmt, path, indent, out, ctx);
      case SWITCH ->
          emitSwitchCommon(
              ((SwitchTree) stmt).getExpression(),
              ((SwitchTree) stmt).getCases(),
              path,
              indent,
              false,
              null,
              null,
              out,
              ctx);
      case TRY -> emitTry((TryTree) stmt, path, indent, out, ctx);
      default -> tabs(out, indent).append("// unsupported stmt: ").append(stmt.getKind()).append('\n');
    }
  }

  public String emitExprAsStmt(ExpressionTree expr, TreePath parentPath, EmitContext ctx) {
    while (expr instanceof ParenthesizedTree pt) {
      parentPath = new TreePath(parentPath, pt);
      expr = pt.getExpression();
    }
    TreePath exprPath = new TreePath(parentPath, expr);
    if (expr instanceof AssignmentTree at) {
      return exprEmitter.emitAssignmentStatement(at, exprPath, ctx);
    }
    if (expr instanceof CompoundAssignmentTree cat) {
      return exprEmitter.emitCompoundAssignmentStatement(cat, exprPath, ctx);
    }
    if (expr instanceof UnaryTree ut) {
      return exprEmitter.emitUnaryStatement(ut, exprPath, ctx);
    }
    if (expr.getKind() == Tree.Kind.METHOD_INVOCATION) {
      String precondStmt = exprEmitter.tryEmitPreconditionsStmt(exprPath, ctx);
      if (precondStmt != null) {
        return precondStmt;
      }
    }
    String go = exprEmitter.emitPath(exprPath, ctx);
    TypeMirror tm = exprEmitter.typeOf(exprPath);
    if (tm != null && tm.getKind() == TypeKind.VOID) {
      return go;
    }
    if (expr.getKind() == Tree.Kind.METHOD_INVOCATION && isBareGoCallStmt(go)) {
      return go;
    }
    return "_ = " + go;
  }

  private static boolean isBareGoCallStmt(String go) {
    int lp = go.indexOf('(');
    if (lp <= 0 || !go.endsWith(")") || go.startsWith("(")) {
      return false;
    }
    String head = go.substring(0, lp);
    return !head.equals("int8")
        && !head.equals("int16")
        && !head.equals("int32")
        && !head.equals("int64")
        && !head.equals("uint16")
        && !head.equals("float32")
        && !head.equals("float64")
        && !head.equals("bool")
        && !head.equals("Node")
        && !head.endsWith(".Node");
  }

  private void emitPatternVarPredecls(
      TreePath stmtPath,
      int indent,
      Set<VariableElement> seen,
      StringBuilder out,
      EmitContext ctx) {
    new TreePathScanner<Void, Void>() {
      @Override
      public Void visitBlock(BlockTree node, Void p) {
        if (getCurrentPath() == stmtPath) {
          return super.visitBlock(node, p);
        }
        return null;
      }

      @Override
      public Void visitLambdaExpression(LambdaExpressionTree node, Void p) {
        return null;
      }

      @Override
      public Void visitClass(ClassTree node, Void p) {
        return null;
      }

      @Override
      public Void visitSwitch(SwitchTree node, Void p) {
        if (getCurrentPath() == stmtPath) {
          scan(node.getExpression(), p);
        }
        return null;
      }

      @Override
      public Void visitSwitchExpression(SwitchExpressionTree node, Void p) {
        return null;
      }

      @Override
      public Void visitInstanceOf(InstanceOfTree node, Void p) {
        if (node.getPattern() instanceof BindingPatternTree bpt) {
          TreePath varPath =
              new TreePath(new TreePath(getCurrentPath(), node.getPattern()), bpt.getVariable());
          if (resolver.trees.getElement(varPath) instanceof VariableElement ve && seen.add(ve)) {
            String vName = ctx.registerLocalVar(ve);
            String vGoType = resolver.goType(ve.asType(), ctx.imports);
            tabs(out, indent).append("var ").append(vName).append(' ').append(vGoType).append('\n');
            tabs(out, indent).append("_ = ").append(vName).append('\n');
          }
        }
        return super.visitInstanceOf(node, p);
      }
    }.scan(stmtPath, null);
  }

  private void emitVariableDecl(
      VariableTree vt, TreePath path, int indent, StringBuilder out, EmitContext ctx) {
    Element el = resolver.trees.getElement(path);
    if (!(el instanceof VariableElement ve)) {
      return;
    }
    String varName = ctx.registerLocalVar(ve);
    TypeMirror tm = ve.asType();
    String goType = resolver.goType(tm, ctx.imports);
    if (vt.getInitializer() != null) {
      TreePath initPath = unwrapParensPath(new TreePath(path, vt.getInitializer()));
      if (initPath.getLeaf() instanceof SwitchExpressionTree set
          && isArrowExprOrThrowSwitch(set)
          && !referencesVariable(initPath, ve)) {
        tabs(out, indent)
            .append("var ")
            .append(varName)
            .append(' ')
            .append(goType)
            .append(" = ")
            .append(resolver.zeroValue(tm, ctx.imports))
            .append('\n');
        tabs(out, indent).append("_ = ").append(varName).append('\n');
        ctx.switchExprTargetStack.push(new EmitContext.SwitchExprTarget(varName, tm, false));
        emitSwitchCommon(
            set.getExpression(), set.getCases(), initPath, indent, true, tm, null, out, ctx);
        ctx.switchExprTargetStack.pop();
        return;
      }
      if (ve.getKind() == ElementKind.LOCAL_VARIABLE
          && initPath.getLeaf() instanceof NewClassTree nct
          && tryEmitStackNewClassVarDecl(
              nct, initPath, ve, varName, goType, path, indent, out, ctx)) {
        return;
      }
      if (initPath.getLeaf() instanceof ConditionalExpressionTree ce) {
        tabs(out, indent).append("var ").append(varName).append(' ').append(goType).append('\n');
        tabs(out, indent).append("_ = ").append(varName).append('\n');
        emitUnwrappedConditional(
            ce,
            initPath,
            indent,
            false,
            tm,
            false,
            (val, branchIndent) ->
                tabs(out, branchIndent).append(varName).append(" = ").append(val).append('\n'),
            out,
            ctx);
        return;
      }
    }
    if (vt.getInitializer() != null
        && referencesVariable(new TreePath(path, vt.getInitializer()), ve)) {
      String zero = resolver.zeroValue(tm, ctx.imports);
      tabs(out, indent)
          .append("var ")
          .append(varName)
          .append(' ')
          .append(goType)
          .append(" = ")
          .append(zero)
          .append('\n');
      tabs(out, indent).append("_ = ").append(varName).append('\n');
      String init = exprEmitter.emitAs(vt.getInitializer(), path, tm, false, ctx);
      tabs(out, indent).append(varName).append(" = ").append(init).append('\n');
      return;
    }
    String init =
        vt.getInitializer() != null
            ? exprEmitter.emitAs(vt.getInitializer(), path, tm, false, ctx)
            : resolver.zeroValue(tm, ctx.imports);
    tabs(out, indent)
        .append("var ")
        .append(varName)
        .append(' ')
        .append(goType)
        .append(" = ")
        .append(init)
        .append('\n');
    tabs(out, indent).append("_ = ").append(varName).append('\n');
  }

  private record StackNewClassTarget(
      ExecutableElement ctorEe,
      TypeElement ctorOwner,
      String goTypeName,
      boolean hasClinit,
      String outerVal) {}

  private StackNewClassTarget resolveStackNewClassTarget(
      NewClassTree nct, TreePath nctPath, TreePath stmtPath, EmitContext ctx) {
    if (nct.getClassBody() != null || nct.getEnclosingExpression() != null) {
      return null;
    }
    Element ctorEl = resolver.trees.getElement(nctPath);
    if (!(ctorEl instanceof ExecutableElement ctorEe)
        || ctorEe.getKind() != ElementKind.CONSTRUCTOR
        || !(ctorEe.getEnclosingElement() instanceof TypeElement ctorOwner)) {
      return null;
    }
    if (!resolver.isSourceType(ctorOwner)
        || ctorOwner.getKind() != ElementKind.CLASS
        || ctorOwner.getModifiers().contains(Modifier.ABSTRACT)
        || ctorOwner.getNestingKind() == NestingKind.ANONYMOUS
        || resolver.isNodeClass(ctorOwner)
        || resolver.isValueStructClass(ctorOwner)
        || resolver.isValueEnum(ctorOwner)
        || resolver.hasEffectiveTypeParameters(ctorOwner)
        || resolver.hasFastInstanceSyncMutex(ctorOwner)
        || resolver.isEnumConstantAnonType(ctorOwner)
        || resolver.isSubclassed(ctorOwner)
        || !resolver.getCapturedVars(ctorOwner).isEmpty()) {
      return null;
    }
    TypeElement superTe = resolver.getSuperclassElement(ctorOwner);
    if (superTe != null && !resolver.runtime.mapsToAny(superTe)) {
      return null;
    }
    if (!resolver.planner.getGoPkg(ctorOwner).equals(ctx.unitPlan.goPkg)) {
      return null;
    }
    if (isInsideEnclosingLoop(stmtPath)) {
      return null;
    }
    TreePath ownerCtPath = resolver.trees.getPath(ctorOwner);
    if (ownerCtPath == null || !(ownerCtPath.getLeaf() instanceof ClassTree ownerCt)) {
      return null;
    }
    String outerVal = null;
    if (resolver.hasOuterInstance(ctorOwner)) {
      TypeElement outerTe = resolver.getEnclosingInstanceType(ctorOwner);
      if (outerTe == null) {
        return null;
      }
      outerVal = exprEmitter.resolveImplicitThisValue(outerTe, ctx);
    }
    String goTypeName = resolver.planner.getGoTypeName(ctorOwner);
    boolean hasClinit =
        ClassEmitter.hasEagerMethodEntryClinit(resolver, ownerCt, ownerCtPath, ctorOwner);
    return new StackNewClassTarget(ctorEe, ctorOwner, goTypeName, hasClinit, outerVal);
  }

  private boolean tryEmitStackNewClassVarDecl(
      NewClassTree nct,
      TreePath initPath,
      VariableElement ve,
      String varName,
      String goType,
      TreePath stmtPath,
      int indent,
      StringBuilder out,
      EmitContext ctx) {
    if (referencesVariable(initPath, ve)) {
      return false;
    }
    StackNewClassTarget target = resolveStackNewClassTarget(nct, initPath, stmtPath, ctx);
    if (target == null || !goType.equals("*" + target.goTypeName())) {
      return false;
    }

    String valTmp = ctx.nextTemp("val");
    String initFnName = "Init" + resolver.getConstructorGoName(target.ctorEe()).substring(3);
    List<String> initArgs = new ArrayList<>();
    initArgs.add("&" + valTmp);
    exprEmitter.emitNewClassInitArgs(nct, initPath, target.ctorEe(), initArgs, ctx);

    if (target.hasClinit()) {
      tabs(out, indent).append("clinit_").append(target.goTypeName()).append("()\n");
    }
    tabs(out, indent).append("var ").append(valTmp).append(' ').append(target.goTypeName()).append('\n');
    if (target.outerVal() != null) {
      tabs(out, indent).append(valTmp).append(".outer_ = ").append(target.outerVal()).append('\n');
    }
    tabs(out, indent)
        .append(initFnName)
        .append('(')
        .append(String.join(", ", initArgs))
        .append(")\n");
    tabs(out, indent)
        .append("var ")
        .append(varName)
        .append(' ')
        .append(goType)
        .append(" = &")
        .append(valTmp)
        .append('\n');
    tabs(out, indent).append("_ = ").append(varName).append('\n');
    return true;
  }

  private boolean tryEmitStackNewClassImmediateCallStmt(
      ExpressionTree expr, TreePath stmtPath, int indent, StringBuilder out, EmitContext ctx) {
    TreePath callPath = unwrapParensPath(new TreePath(stmtPath, expr));
    if (!(callPath.getLeaf() instanceof MethodInvocationTree mi)
        || !(mi.getMethodSelect() instanceof MemberSelectTree ms)) {
      return false;
    }
    TreePath msPath = new TreePath(callPath, ms);
    TreePath recvPath = unwrapParensPath(new TreePath(msPath, ms.getExpression()));
    if (!(recvPath.getLeaf() instanceof NewClassTree nct)) {
      return false;
    }
    if (!(resolver.trees.getElement(callPath) instanceof ExecutableElement methodEe)
        || methodEe.getModifiers().contains(Modifier.STATIC)
        || !(methodEe.getEnclosingElement() instanceof TypeElement methodOwner)
        || resolver.runtime.mapsToAny(methodOwner)
        || resolver.runtime.isShimException(methodOwner)) {
      return false;
    }
    TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(methodEe);
    if (sig == null) {
      return false;
    }
    StackNewClassTarget target = resolveStackNewClassTarget(nct, recvPath, stmtPath, ctx);
    if (target == null) {
      return false;
    }
    if (!methodOwner.equals(target.ctorOwner())
        && !(methodOwner.getKind() == ElementKind.INTERFACE
            && (resolver.isSourceType(methodOwner) || resolver.isNeededExternalSig(sig)))) {
      return false;
    }

    String valTmp = ctx.nextTemp("val");
    String initFnName = "Init" + resolver.getConstructorGoName(target.ctorEe()).substring(3);
    List<String> initArgs = new ArrayList<>();
    initArgs.add("&" + valTmp);
    exprEmitter.emitNewClassInitArgs(nct, recvPath, target.ctorEe(), initArgs, ctx);

    ExecutableType execType =
        resolver.resolveExecutableType(target.ctorOwner().asType(), methodEe);
    List<String> mArgs = new ArrayList<>();
    exprEmitter.appendFormattedExecArgs(
        mi.getArguments(), callPath, methodEe, execType, sig, mArgs, ctx);

    tabs(out, indent).append("{\n");
    if (target.hasClinit()) {
      tabs(out, indent + 1).append("clinit_").append(target.goTypeName()).append("()\n");
    }
    tabs(out, indent + 1)
        .append("var ")
        .append(valTmp)
        .append(' ')
        .append(target.goTypeName())
        .append('\n');
    if (target.outerVal() != null) {
      tabs(out, indent + 1)
          .append(valTmp)
          .append(".outer_ = ")
          .append(target.outerVal())
          .append('\n');
    }
    tabs(out, indent + 1)
        .append(initFnName)
        .append('(')
        .append(String.join(", ", initArgs))
        .append(")\n");
    tabs(out, indent + 1)
        .append("(&")
        .append(valTmp)
        .append(").")
        .append(sig.goName)
        .append('(')
        .append(String.join(", ", mArgs))
        .append(")\n");
    tabs(out, indent).append("}\n");
    return true;
  }

  private boolean referencesVariable(TreePath exprPath, VariableElement targetVe) {
    boolean[] found = {false};
    new TreePathScanner<Void, Void>() {
      @Override
      public Void scan(Tree tree, Void p) {
        if (found[0] || tree == null) {
          return null;
        }
        TreePath tp = tree == exprPath.getLeaf() ? exprPath : new TreePath(getCurrentPath(), tree);
        if (targetVe.equals(resolver.trees.getElement(tp))) {
          found[0] = true;
          return null;
        }
        return super.scan(tree, p);
      }
    }.scan(exprPath, null);
    return found[0];
  }

  private boolean tryEmitUnwrappedConditionalAssignStmt(
      ExpressionTree expr, TreePath stmtPath, int indent, StringBuilder out, EmitContext ctx) {
    TreePath exprPath = unwrapParensPath(new TreePath(stmtPath, expr));
    if (!(exprPath.getLeaf() instanceof AssignmentTree at)) {
      return false;
    }
    TreePath rhsPath = unwrapParensPath(new TreePath(exprPath, at.getExpression()));
    if (!(rhsPath.getLeaf() instanceof ConditionalExpressionTree ce)) {
      return false;
    }
    TreePath lhsPath = new TreePath(exprPath, at.getVariable());
    if (!exprEmitter.isPureSideEffectFreeAssignLhs(lhsPath, ctx)) {
      return false;
    }
    TypeMirror lhsType = exprEmitter.typeOf(lhsPath);
    emitUnwrappedConditional(
        ce,
        rhsPath,
        indent,
        false,
        lhsType,
        false,
        (val, branchIndent) ->
            tabs(out, branchIndent)
                .append(exprEmitter.emitPureLhsWrite(lhsPath, val, ctx))
                .append('\n'),
        out,
        ctx);
    return true;
  }

  private void emitUnwrappedConditional(
      ConditionalExpressionTree ce,
      TreePath cePath,
      int indent,
      boolean isElseIf,
      TypeMirror targetTm,
      boolean forceAny,
      java.util.function.ObjIntConsumer<String> leafSink,
      StringBuilder out,
      EmitContext ctx) {
    TypeMirror ceTm = exprEmitter.typeOf(cePath);
    String cond = null;
    if (!isElseIf) {
      cond = tryEmitPatternCondition(new TreePath(cePath, ce.getCondition()), indent, false, out, ctx);
    }
    if (cond == null) {
      cond =
          exprEmitter.emitAs(
              ce.getCondition(),
              cePath,
              resolver.types.getPrimitiveType(TypeKind.BOOLEAN),
              false,
              ctx);
    }
    if (!isElseIf) {
      tabs(out, indent);
    }
    out.append("if ").append(cond).append(" {\n");
    emitConditionalBranchLeaf(
        new TreePath(cePath, ce.getTrueExpression()),
        ceTm,
        targetTm,
        forceAny,
        indent + 1,
        false,
        leafSink,
        out,
        ctx);
    TreePath falsePath = unwrapParensPath(new TreePath(cePath, ce.getFalseExpression()));
    if (falsePath.getLeaf() instanceof ConditionalExpressionTree nestedFalseCe
        && (ceTm == null || !ceTm.getKind().isPrimitive())
        && !containsImmediateBindingPattern(new TreePath(falsePath, nestedFalseCe.getCondition()))) {
      tabs(out, indent).append("} else ");
      emitUnwrappedConditional(
          nestedFalseCe, falsePath, indent, true, targetTm, forceAny, leafSink, out, ctx);
    } else {
      tabs(out, indent).append("} else {\n");
      emitConditionalBranchLeaf(
          new TreePath(cePath, ce.getFalseExpression()),
          ceTm,
          targetTm,
          forceAny,
          indent + 1,
          true,
          leafSink,
          out,
          ctx);
      tabs(out, indent).append("}\n");
    }
  }

  private void emitConditionalBranchLeaf(
      TreePath rawBranchPath,
      TypeMirror ceTm,
      TypeMirror targetTm,
      boolean forceAny,
      int indent,
      boolean allowNestedUnwrap,
      java.util.function.ObjIntConsumer<String> leafSink,
      StringBuilder out,
      EmitContext ctx) {
    TreePath unwrapped = unwrapParensPath(rawBranchPath);
    if (unwrapped.getLeaf() instanceof ConditionalExpressionTree nestedCe
        && (ceTm == null || !ceTm.getKind().isPrimitive())) {
      emitUnwrappedConditional(
          nestedCe, unwrapped, indent, false, targetTm, forceAny, leafSink, out, ctx);
      return;
    }
    if ((ceTm == null || !ceTm.getKind().isPrimitive())
        && exprEmitter.typeOf(unwrapped) != null
        && exprEmitter.typeOf(unwrapped).getKind() == TypeKind.NULL) {
      String zero = forceAny ? "nil" : resolver.zeroValue(targetTm, ctx.imports);
      leafSink.accept(zero, indent);
      return;
    }
    String branchVal =
        exprEmitter.emitAs(
            (ExpressionTree) rawBranchPath.getLeaf(), rawBranchPath.getParentPath(), ceTm, false, ctx);
    String converted = exprEmitter.convert(branchVal, ceTm, targetTm, forceAny, ctx);
    leafSink.accept(converted, indent);
  }

  private String emitBoolTreePath(TreePath p, EmitContext ctx) {
    return exprEmitter.emitAs(
        (ExpressionTree) p.getLeaf(),
        p.getParentPath(),
        resolver.types.getPrimitiveType(TypeKind.BOOLEAN),
        false,
        ctx);
  }

  private void collectBinaryChain(TreePath p, Tree.Kind opKind, List<TreePath> out) {
    TreePath u = unwrapParensPath(p);
    if (u.getLeaf() instanceof BinaryTree bt && bt.getKind() == opKind) {
      collectBinaryChain(new TreePath(u, bt.getLeftOperand()), opKind, out);
      collectBinaryChain(new TreePath(u, bt.getRightOperand()), opKind, out);
    } else {
      out.add(u);
    }
  }

  private boolean containsImmediateBindingPattern(TreePath exprPath) {
    boolean[] found = {false};
    new TreePathScanner<Void, Void>() {
      @Override
      public Void visitLambdaExpression(LambdaExpressionTree n, Void p) {
        return null;
      }

      @Override
      public Void visitClass(ClassTree n, Void p) {
        return null;
      }

      @Override
      public Void visitSwitchExpression(SwitchExpressionTree n, Void p) {
        return null;
      }

      @Override
      public Void visitInstanceOf(InstanceOfTree n, Void p) {
        if (n.getPattern() instanceof BindingPatternTree) {
          found[0] = true;
          return null;
        }
        return super.visitInstanceOf(n, p);
      }
    }.scan(exprPath, null);
    return found[0];
  }

  private boolean isSimpleBindingInstanceOf(TreePath p) {
    return p.getLeaf() instanceof InstanceOfTree io
        && io.getPattern() instanceof BindingPatternTree
        && !containsImmediateBindingPattern(new TreePath(p, io.getExpression()));
  }

  private boolean canEmitConjunctivePatternCond(TreePath condPath) {
    List<TreePath> conjuncts = new ArrayList<>();
    collectBinaryChain(condPath, Tree.Kind.CONDITIONAL_AND, conjuncts);
    boolean hasBind = false;
    for (TreePath c : conjuncts) {
      if (isSimpleBindingInstanceOf(c)) {
        hasBind = true;
      } else if (containsImmediateBindingPattern(c)) {
        return false;
      }
    }
    return hasBind;
  }

  private String tryEmitConjunctivePatternCond(
      TreePath condPath,
      int indent,
      boolean foldTrailingIntoOkVar,
      StringBuilder out,
      EmitContext ctx) {
    List<TreePath> conjuncts = new ArrayList<>();
    collectBinaryChain(condPath, Tree.Kind.CONDITIONAL_AND, conjuncts);
    int lastBindIdx = -1;
    for (int i = 0; i < conjuncts.size(); i++) {
      TreePath c = conjuncts.get(i);
      if (isSimpleBindingInstanceOf(c)) {
        lastBindIdx = i;
      } else if (containsImmediateBindingPattern(c)) {
        return null;
      }
    }
    if (lastBindIdx < 0) {
      return null;
    }
    String okVar = ctx.nextTemp("ok");
    int evalThrough = foldTrailingIntoOkVar ? (conjuncts.size() - 1) : lastBindIdx;
    TreePath c0 = conjuncts.get(0);
    if (isSimpleBindingInstanceOf(c0)) {
      exprEmitter.emitPatternTryCastAssign(
          (InstanceOfTree) c0.getLeaf(), c0, okVar, true, "\t".repeat(indent), out, ctx);
    } else {
      tabs(out, indent)
          .append(okVar)
          .append(" := ")
          .append(emitBoolTreePath(c0, ctx))
          .append('\n');
    }
    for (int i = 1; i <= evalThrough; i++) {
      TreePath ci = conjuncts.get(i);
      tabs(out, indent).append("if ").append(okVar).append(" {\n");
      if (isSimpleBindingInstanceOf(ci)) {
        exprEmitter.emitPatternTryCastAssign(
            (InstanceOfTree) ci.getLeaf(), ci, okVar, false, "\t".repeat(indent + 1), out, ctx);
      } else {
        tabs(out, indent + 1)
            .append(okVar)
            .append(" = ")
            .append(emitBoolTreePath(ci, ctx))
            .append('\n');
      }
      tabs(out, indent).append("}\n");
    }
    if (evalThrough == conjuncts.size() - 1) {
      return okVar;
    }
    StringBuilder rem = new StringBuilder(okVar);
    for (int i = evalThrough + 1; i < conjuncts.size(); i++) {
      rem.append(" && (").append(emitBoolTreePath(conjuncts.get(i), ctx)).append(')');
    }
    return rem.toString();
  }

  private String tryEmitPatternCondition(
      TreePath rawCondPath,
      int indent,
      boolean foldTrailingIntoOkVar,
      StringBuilder out,
      EmitContext ctx) {
    TreePath condPath = unwrapParensPath(rawCondPath);
    String pos = tryEmitConjunctivePatternCond(condPath, indent, foldTrailingIntoOkVar, out, ctx);
    if (pos != null) {
      return pos;
    }
    if (condPath.getLeaf() instanceof UnaryTree ut
        && ut.getKind() == Tree.Kind.LOGICAL_COMPLEMENT) {
      TreePath inner = unwrapParensPath(new TreePath(condPath, ut.getExpression()));
      String negOk = tryEmitConjunctivePatternCond(inner, indent, true, out, ctx);
      if (negOk != null) {
        return "!" + negOk;
      }
    }
    List<TreePath> disjuncts = new ArrayList<>();
    collectBinaryChain(condPath, Tree.Kind.CONDITIONAL_OR, disjuncts);
    if (disjuncts.size() >= 2) {
      int lastBindDisjIdx = -1;
      for (int j = 0; j < disjuncts.size(); j++) {
        TreePath d = disjuncts.get(j);
        if (d.getLeaf() instanceof UnaryTree ut && ut.getKind() == Tree.Kind.LOGICAL_COMPLEMENT) {
          TreePath inner = unwrapParensPath(new TreePath(d, ut.getExpression()));
          if (canEmitConjunctivePatternCond(inner)) {
            lastBindDisjIdx = j;
            continue;
          }
        }
        if (containsImmediateBindingPattern(d)) {
          return null;
        }
      }
      if (lastBindDisjIdx >= 0) {
        String orVar = ctx.nextTemp("ok");
        int evalThrough = foldTrailingIntoOkVar ? (disjuncts.size() - 1) : lastBindDisjIdx;
        TreePath d0 = disjuncts.get(0);
        if (d0.getLeaf() instanceof UnaryTree ut0
            && ut0.getKind() == Tree.Kind.LOGICAL_COMPLEMENT
            && canEmitConjunctivePatternCond(unwrapParensPath(new TreePath(d0, ut0.getExpression())))) {
          String innerOk =
              tryEmitConjunctivePatternCond(
                  unwrapParensPath(new TreePath(d0, ut0.getExpression())), indent, true, out, ctx);
          tabs(out, indent).append(orVar).append(" := !").append(innerOk).append('\n');
        } else {
          tabs(out, indent)
              .append(orVar)
              .append(" := ")
              .append(emitBoolTreePath(d0, ctx))
              .append('\n');
        }
        for (int j = 1; j <= evalThrough; j++) {
          TreePath dj = disjuncts.get(j);
          tabs(out, indent).append("if !").append(orVar).append(" {\n");
          if (dj.getLeaf() instanceof UnaryTree utj
              && utj.getKind() == Tree.Kind.LOGICAL_COMPLEMENT
              && canEmitConjunctivePatternCond(
                  unwrapParensPath(new TreePath(dj, utj.getExpression())))) {
            String innerOk =
                tryEmitConjunctivePatternCond(
                    unwrapParensPath(new TreePath(dj, utj.getExpression())),
                    indent + 1,
                    true,
                    out,
                    ctx);
            tabs(out, indent + 1).append(orVar).append(" = !").append(innerOk).append('\n');
          } else {
            tabs(out, indent + 1)
                .append(orVar)
                .append(" = ")
                .append(emitBoolTreePath(dj, ctx))
                .append('\n');
          }
          tabs(out, indent).append("}\n");
        }
        if (evalThrough == disjuncts.size() - 1) {
          return orVar;
        }
        StringBuilder rem = new StringBuilder(orVar);
        for (int j = evalThrough + 1; j < disjuncts.size(); j++) {
          rem.append(" || (").append(emitBoolTreePath(disjuncts.get(j), ctx)).append(')');
        }
        return rem.toString();
      }
    }
    return null;
  }

  private void emitIf(
      IfTree it, TreePath path, int indent, boolean isElseIf, StringBuilder out, EmitContext ctx) {
    String cond = null;
    if (!isElseIf) {
      cond = tryEmitPatternCondition(new TreePath(path, it.getCondition()), indent, false, out, ctx);
    }
    if (cond == null) {
      cond =
          exprEmitter.emitAs(
              it.getCondition(),
              path,
              resolver.types.getPrimitiveType(TypeKind.BOOLEAN),
              false,
              ctx);
    }
    if (!isElseIf) {
      tabs(out, indent);
    }
    out.append("if ").append(cond).append(" {\n");
    emitAsBody(it.getThenStatement(), new TreePath(path, it.getThenStatement()), indent + 1, out, ctx);
    if (it.getElseStatement() != null) {
      StatementTree elseStmt = it.getElseStatement();
      TreePath elsePath = new TreePath(path, elseStmt);
      // Check if else-if has pattern variables that need predeclaration; if so, emit nested block instead of `else if`!
      if (elseStmt instanceof IfTree nestedIf && !hasImmediatePatternVars(elsePath)) {
        tabs(out, indent).append("} else ");
        emitIf(nestedIf, elsePath, indent, true, out, ctx);
      } else {
        tabs(out, indent).append("} else {\n");
        emitAsBody(elseStmt, elsePath, indent + 1, out, ctx);
        tabs(out, indent).append("}\n");
      }
    } else {
      tabs(out, indent).append("}\n");
    }
  }

  private boolean hasImmediatePatternVars(TreePath stmtPath) {
    boolean[] found = {false};
    new TreePathScanner<Void, Void>() {
      @Override
      public Void visitBlock(BlockTree n, Void p) {
        return null;
      }

      @Override
      public Void visitLambdaExpression(LambdaExpressionTree n, Void p) {
        return null;
      }

      @Override
      public Void visitClass(ClassTree n, Void p) {
        return null;
      }

      @Override
      public Void visitInstanceOf(InstanceOfTree n, Void p) {
        if (n.getPattern() != null) found[0] = true;
        return super.visitInstanceOf(n, p);
      }
    }.scan(stmtPath, null);
    return found[0];
  }

  private void emitAsBody(
      StatementTree stmt, TreePath stmtPath, int indent, StringBuilder out, EmitContext ctx) {
    if (stmt instanceof BlockTree bt) {
      out.append(emitBlockStatements(bt, stmtPath, indent, ctx));
    } else {
      emitStmt(stmt, stmtPath, indent, null, out, ctx);
    }
  }

  // --- Loops & Labeled Statements ---

  private void emitLabeled(
      LabeledStatementTree lst, TreePath path, int indent, StringBuilder out, EmitContext ctx) {
    String goLabel = GoNames.localIdent(lst.getLabel().toString());
    StatementTree inner = lst.getStatement();
    TreePath innerPath = new TreePath(path, inner);
    if (inner instanceof WhileLoopTree
        || inner instanceof DoWhileLoopTree
        || inner instanceof ForLoopTree
        || inner instanceof EnhancedForLoopTree) {
      emitStmt(inner, innerPath, indent, goLabel, out, ctx);
      return;
    }
    // Non-loop labeled block (`complete: { ... break complete; }`): lower to labeled Go switch!
    tabs(out, indent).append(goLabel).append(":\n");
    tabs(out, indent).append("switch {\n");
    tabs(out, indent).append("default:\n");
    emitAsBody(inner, innerPath, indent + 1, out, ctx);
    tabs(out, indent).append("}\n");
  }

  private void emitWhile(
      WhileLoopTree wt,
      TreePath path,
      int indent,
      String loopLabel,
      StringBuilder out,
      EmitContext ctx) {
    if (loopLabel != null) {
      tabs(out, indent).append(loopLabel).append(":\n");
    }
    String cond =
        exprEmitter.emitAs(
            wt.getCondition(),
            path,
            resolver.types.getPrimitiveType(TypeKind.BOOLEAN),
            false,
            ctx);
    if (cond.equals("true") || cond.equals("(true)")) {
      tabs(out, indent).append("for {\n");
    } else {
      tabs(out, indent).append("for ").append(cond).append(" {\n");
    }
    if (ctx.inTryWrapper) ctx.loopDepthInTryIife++;
    emitAsBody(wt.getStatement(), new TreePath(path, wt.getStatement()), indent + 1, out, ctx);
    if (ctx.inTryWrapper) ctx.loopDepthInTryIife--;
    tabs(out, indent).append("}\n");
  }

  private void emitDoWhile(
      DoWhileLoopTree dwt,
      TreePath path,
      int indent,
      String loopLabel,
      StringBuilder out,
      EmitContext ctx) {
    if (loopLabel != null) {
      tabs(out, indent).append(loopLabel).append(":\n");
    }
    String firstVar = ctx.nextTemp("doFirst");
    String cond =
        exprEmitter.emitAs(
            dwt.getCondition(),
            path,
            resolver.types.getPrimitiveType(TypeKind.BOOLEAN),
            false,
            ctx);
    tabs(out, indent)
        .append("for ")
        .append(firstVar)
        .append(" := true; ")
        .append(firstVar)
        .append(" || (")
        .append(cond)
        .append("); ")
        .append(firstVar)
        .append(" = false {\n");
    if (ctx.inTryWrapper) ctx.loopDepthInTryIife++;
    emitAsBody(dwt.getStatement(), new TreePath(path, dwt.getStatement()), indent + 1, out, ctx);
    if (ctx.inTryWrapper) ctx.loopDepthInTryIife--;
    tabs(out, indent).append("}\n");
  }

  private void emitFor(
      ForLoopTree ft,
      TreePath path,
      int indent,
      String loopLabel,
      StringBuilder out,
      EmitContext ctx) {
    boolean hasInit = !ft.getInitializer().isEmpty();
    int bodyIndent = indent;
    if (hasInit) {
      tabs(out, indent).append("{\n");
      bodyIndent = indent + 1;
      for (StatementTree initStmt : ft.getInitializer()) {
        emitStmt(initStmt, new TreePath(path, initStmt), bodyIndent, null, out, ctx);
      }
    }
    if (loopLabel != null) {
      tabs(out, bodyIndent).append(loopLabel).append(":\n");
    }

    String cond =
        ft.getCondition() != null
            ? exprEmitter.emitAs(
                ft.getCondition(),
                path,
                resolver.types.getPrimitiveType(TypeKind.BOOLEAN),
                false,
                ctx)
            : "";

    String post = "";
    if (!ft.getUpdate().isEmpty()) {
      if (ft.getUpdate().size() == 1) {
        ExpressionStatementTree u0 = ft.getUpdate().get(0);
        TreePath u0Path = new TreePath(path, u0);
        post = emitExprAsStmt(u0.getExpression(), u0Path, ctx);
        // If `post` wrapped in a `{ ... }` block or `if ...` stmt, turn into `func() { ... }()` so it is a valid Go SimpleStmt
        if (post.startsWith("{")) {
          post = "func() " + post + "()";
        } else if (post.startsWith("if ")) {
          post = "func() { " + post + " }()";
        }
      } else {
        List<String> upds = new ArrayList<>();
        for (ExpressionStatementTree u : ft.getUpdate()) {
          upds.add(emitExprAsStmt(u.getExpression(), new TreePath(path, u), ctx));
        }
        post = "func() { " + String.join("; ", upds) + " }()";
      }
    }

    tabs(out, bodyIndent).append("for ");
    if (!cond.isEmpty() || !post.isEmpty()) {
      out.append("; ").append(cond).append("; ").append(post).append(' ');
    }
    out.append("{\n");
    if (ctx.inTryWrapper) ctx.loopDepthInTryIife++;
    emitAsBody(ft.getStatement(), new TreePath(path, ft.getStatement()), bodyIndent + 1, out, ctx);
    if (ctx.inTryWrapper) ctx.loopDepthInTryIife--;
    tabs(out, bodyIndent).append("}\n");

    if (hasInit) {
      tabs(out, indent).append("}\n");
    }
  }

  private void emitEnhancedFor(
      EnhancedForLoopTree eft,
      TreePath path,
      int indent,
      String loopLabel,
      StringBuilder out,
      EmitContext ctx) {
    TreePath varPath = new TreePath(path, eft.getVariable());
    VariableElement ve = (VariableElement) resolver.trees.getElement(varPath);
    String varName = ctx.registerLocalVar(ve);
    TypeMirror varTm = ve.asType();

    TreePath exprPath = new TreePath(path, eft.getExpression());
    if (tryEmitMapCursorEnhancedFor(
        eft, path, ve, varName, varTm, exprPath, indent, loopLabel, out, ctx)) {
      return;
    }
    TypeMirror exprTm = exprEmitter.typeOf(exprPath);
    String exprGo = exprEmitter.emitPath(exprPath, ctx);

    if (exprTm != null && exprTm.getKind() == TypeKind.ARRAY) {
      if (loopLabel != null) {
        tabs(out, indent).append(loopLabel).append(":\n");
      }
      ArrayType at = (ArrayType) resolver.goEffectiveType(exprTm);
      String elemTmp = ctx.nextTemp("el");
      String rt = rtPrefix(ctx);
      tabs(out, indent)
          .append("for _, ")
          .append(elemTmp)
          .append(" := range ")
          .append(rt)
          .append("CheckNotNullSlice(")
          .append(exprGo)
          .append(") {\n");
      String conv = exprEmitter.convert(elemTmp, at.getComponentType(), varTm, false, ctx);
      tabs(out, indent + 1).append(varName).append(" := ").append(conv).append('\n');
      tabs(out, indent + 1).append("_ = ").append(varName).append('\n');

      if (ctx.inTryWrapper) ctx.loopDepthInTryIife++;
      emitAsBody(eft.getStatement(), new TreePath(path, eft.getStatement()), indent + 1, out, ctx);
      if (ctx.inTryWrapper) ctx.loopDepthInTryIife--;
      tabs(out, indent).append("}\n");
      return;
    }

    String rt = rtPrefix(ctx);
    String checkedColl = exprEmitter.emitCheckNotNull(eft.getExpression(), exprGo, exprTm, ctx);

    String curTmp = ctx.nextTemp("cur");
    tabs(out, indent).append("{\n");
    if ((resolver.isValueStructType(varTm) || resolver.isValueEnumType(varTm))
        && !resolver.isActiveTypeVariable(varTm)
        && exprEmitter.isNonLoweredCollectionExpr(exprPath)) {
      String elemGo = resolver.goType(resolver.safeErasure(varTm), ctx.imports);
      tabs(out, indent + 1)
          .append(curTmp)
          .append(" := ")
          .append(rt)
          .append("NewTypedIterCursor[")
          .append(elemGo)
          .append("](")
          .append(checkedColl)
          .append(")\n");
      if (loopLabel != null) {
        tabs(out, indent + 1).append(loopLabel).append(":\n");
      }
      tabs(out, indent + 1).append("for ").append(curTmp).append(".HasNext() {\n");
      tabs(out, indent + 2)
          .append(varName)
          .append(" := ")
          .append(curTmp)
          .append(".Next()\n");
      tabs(out, indent + 2).append("_ = ").append(varName).append('\n');

      if (ctx.inTryWrapper) ctx.loopDepthInTryIife++;
      emitAsBody(eft.getStatement(), new TreePath(path, eft.getStatement()), indent + 2, out, ctx);
      if (ctx.inTryWrapper) ctx.loopDepthInTryIife--;
      tabs(out, indent + 1).append("}\n");
      tabs(out, indent).append("}\n");
      return;
    }

    tabs(out, indent + 1)
        .append(curTmp)
        .append(" := ")
        .append(rt)
        .append("NewIterCursor(")
        .append(checkedColl)
        .append(")\n");
    if (loopLabel != null) {
      tabs(out, indent + 1).append(loopLabel).append(":\n");
    }
    tabs(out, indent + 1).append("for ").append(curTmp).append(".HasNext() {\n");
    if (resolver.isStringType(varTm)) {
      tabs(out, indent + 2)
          .append(varName)
          .append(" := ")
          .append(curTmp)
          .append(".NextStr()\n");
    } else if (resolver.isNodeType(varTm)) {
      String nodeGo = resolver.goType(resolver.safeErasure(varTm), ctx.imports);
      tabs(out, indent + 2)
          .append(varName)
          .append(" := ")
          .append(nodeGo)
          .append('(')
          .append(curTmp)
          .append(".NextNode())\n");
    } else {
      String rawElemTmp = ctx.nextTemp("rawElem");
      tabs(out, indent + 2)
          .append(rawElemTmp)
          .append(" := ")
          .append(curTmp)
          .append(".Next()\n");
      String conv =
          exprEmitter.convert(rawElemTmp, resolver.getObjectElement().asType(), varTm, false, ctx);
      tabs(out, indent + 2).append(varName).append(" := ").append(conv).append('\n');
    }
    tabs(out, indent + 2).append("_ = ").append(varName).append('\n');

    if (ctx.inTryWrapper) ctx.loopDepthInTryIife++;
    emitAsBody(eft.getStatement(), new TreePath(path, eft.getStatement()), indent + 2, out, ctx);
    if (ctx.inTryWrapper) ctx.loopDepthInTryIife--;
    tabs(out, indent + 1).append("}\n");
    tabs(out, indent).append("}\n");
  }

  private boolean tryEmitMapCursorEnhancedFor(
      EnhancedForLoopTree eft,
      TreePath path,
      VariableElement ve,
      String varName,
      TypeMirror varTm,
      TreePath exprPath,
      int indent,
      String loopLabel,
      StringBuilder out,
      EmitContext ctx) {
    TreePath strippedPath = unwrapParensPath(exprPath);
    if (!(strippedPath.getLeaf() instanceof MethodInvocationTree mi)
        || !mi.getArguments().isEmpty()) {
      return false;
    }
    if (!(resolver.trees.getElement(strippedPath) instanceof ExecutableElement ee)
        || ee.getModifiers().contains(Modifier.STATIC)
        || ee.getModifiers().contains(Modifier.PRIVATE)
        || !ee.getParameters().isEmpty()) {
      return false;
    }
    String mName = ee.getSimpleName().toString();
    if (!mName.equals("values") && !mName.equals("keySet") && !mName.equals("entrySet")) {
      return false;
    }
    if (!(ee.getEnclosingElement() instanceof TypeElement ownerTe)
        || !exprEmitter.isSubtypeOfQualified(ownerTe, "java.util.Map")) {
      return false;
    }

    TreePath bodyPath = new TreePath(path, eft.getStatement());
    if (mName.equals("entrySet")) {
      if (!(resolver.safeErasure(varTm) instanceof DeclaredType vdt)
          || !(vdt.asElement() instanceof TypeElement vte)
          || !vte.getQualifiedName().contentEquals("java.util.Map.Entry")
          || !canUseStackEntryCursor(bodyPath, ve)) {
        return false;
      }
    }

    String recvGo;
    ExpressionTree sel = mi.getMethodSelect();
    if (sel instanceof MemberSelectTree ms) {
      ExpressionTree recvTree = ms.getExpression();
      if ((recvTree instanceof IdentifierTree id && id.getName().contentEquals("super"))
          || (recvTree instanceof MemberSelectTree rms
              && rms.getIdentifier().contentEquals("super"))) {
        return false;
      }
      recvGo = exprEmitter.emitPath(new TreePath(new TreePath(strippedPath, ms), recvTree), ctx);
    } else if (sel instanceof IdentifierTree id) {
      if (id.getName().contentEquals("super")) {
        return false;
      }
      recvGo = exprEmitter.resolveImplicitThisValue(ownerTe, ctx);
    } else {
      return false;
    }

    String rt = rtPrefix(ctx);
    String curTmp = ctx.nextTemp("cur");

    if (mName.equals("entrySet")) {
      tabs(out, indent).append("{\n");
      tabs(out, indent + 1)
          .append(curTmp)
          .append(" := ")
          .append(rt)
          .append("NewMapEntryCursor(")
          .append(recvGo)
          .append(")\n");
      if (loopLabel != null) {
        tabs(out, indent + 1).append(loopLabel).append(":\n");
      }
      tabs(out, indent + 1).append("for ").append(curTmp).append(".NextEntry() {\n");
      exprEmitter.pushMapEntryCursor(ve, curTmp);
      try {
        if (ctx.inTryWrapper) ctx.loopDepthInTryIife++;
        emitAsBody(eft.getStatement(), bodyPath, indent + 2, out, ctx);
        if (ctx.inTryWrapper) ctx.loopDepthInTryIife--;
      } finally {
        exprEmitter.popMapEntryCursor(ve);
      }
      tabs(out, indent + 1).append("}\n");
      tabs(out, indent).append("}\n");
      return true;
    }

    boolean strKey = mName.equals("keySet") && resolver.isStringType(varTm);
    boolean nodeKey = mName.equals("keySet") && !strKey && resolver.isNodeType(varTm);
    String cursorFn =
        mName.equals("values")
            ? "NewMapValCursor"
            : (strKey ? "NewMapKeyStrCursor" : "NewMapKeyCursor");

    tabs(out, indent).append("{\n");
    tabs(out, indent + 1)
        .append(curTmp)
        .append(" := ")
        .append(rt)
        .append(cursorFn)
        .append('(')
        .append(recvGo)
        .append(")\n");
    if (loopLabel != null) {
      tabs(out, indent + 1).append(loopLabel).append(":\n");
    }
    tabs(out, indent + 1).append("for ").append(curTmp).append(".HasNext() {\n");
    if (strKey) {
      tabs(out, indent + 2)
          .append(varName)
          .append(" := ")
          .append(curTmp)
          .append(".NextStr()\n");
    } else if (nodeKey) {
      String nodeGo = resolver.goType(resolver.safeErasure(varTm), ctx.imports);
      tabs(out, indent + 2)
          .append(varName)
          .append(" := ")
          .append(nodeGo)
          .append('(')
          .append(curTmp)
          .append(".NextNode())\n");
    } else {
      String rawElemTmp = ctx.nextTemp("rawElem");
      tabs(out, indent + 2)
          .append(rawElemTmp)
          .append(" := ")
          .append(curTmp)
          .append(".Next()\n");
      String conv =
          exprEmitter.convert(rawElemTmp, resolver.getObjectElement().asType(), varTm, false, ctx);
      tabs(out, indent + 2).append(varName).append(" := ").append(conv).append('\n');
    }
    tabs(out, indent + 2).append("_ = ").append(varName).append('\n');

    if (ctx.inTryWrapper) ctx.loopDepthInTryIife++;
    emitAsBody(eft.getStatement(), bodyPath, indent + 2, out, ctx);
    if (ctx.inTryWrapper) ctx.loopDepthInTryIife--;
    tabs(out, indent + 1).append("}\n");
    tabs(out, indent).append("}\n");
    return true;
  }

  private boolean canUseStackEntryCursor(TreePath bodyPath, VariableElement entryVar) {
    final boolean[] safe = {true};
    new TreePathScanner<Void, Void>() {
      int closureDepth = 0;

      @Override
      public Void visitLambdaExpression(LambdaExpressionTree node, Void p) {
        closureDepth++;
        try {
          return super.visitLambdaExpression(node, p);
        } finally {
          closureDepth--;
        }
      }

      @Override
      public Void visitClass(ClassTree node, Void p) {
        closureDepth++;
        try {
          return super.visitClass(node, p);
        } finally {
          closureDepth--;
        }
      }

      @Override
      public Void visitNewClass(NewClassTree node, Void p) {
        if (node.getClassBody() != null) {
          closureDepth++;
          try {
            return super.visitNewClass(node, p);
          } finally {
            closureDepth--;
          }
        }
        return super.visitNewClass(node, p);
      }

      @Override
      public Void visitIdentifier(IdentifierTree node, Void p) {
        if (!safe[0]) return null;
        Element el = resolver.trees.getElement(getCurrentPath());
        if (entryVar.equals(el)) {
          if (closureDepth > 0 || !isAllowedEntryMethodReceiver(getCurrentPath())) {
            safe[0] = false;
            return null;
          }
        }
        return null;
      }
    }.scan(bodyPath, null);
    return safe[0];
  }

  private static boolean isAllowedEntryMethodReceiver(TreePath idPath) {
    Tree child = idPath.getLeaf();
    TreePath parentPath = idPath.getParentPath();
    while (parentPath != null && parentPath.getLeaf() instanceof ParenthesizedTree) {
      child = parentPath.getLeaf();
      parentPath = parentPath.getParentPath();
    }
    if (parentPath == null
        || !(parentPath.getLeaf() instanceof MemberSelectTree ms)
        || ms.getExpression() != child) {
      return false;
    }
    TreePath callPath = parentPath.getParentPath();
    if (callPath == null
        || !(callPath.getLeaf() instanceof MethodInvocationTree call)
        || call.getMethodSelect() != ms) {
      return false;
    }
    String name = ms.getIdentifier().toString();
    int arity = call.getArguments().size();
    return ((name.equals("getKey") || name.equals("getValue")) && arity == 0)
        || (name.equals("setValue") && arity == 1);
  }

  // --- Jump & Control Statements ---

  private void emitBreak(BreakTree bt, int indent, StringBuilder out, EmitContext ctx) {
    if (bt.getLabel() == null
        && ctx.inTryWrapper
        && ctx.loopDepthInTryIife == 0
        && ctx.switchDepthInTryIife == 0
        && ctx.tryFlowVar != null) {
      tabs(out, indent).append(ctx.tryFlowVar).append(" = 1\n");
      tabs(out, indent).append("return\n");
      return;
    }
    tabs(out, indent).append("break");
    if (bt.getLabel() != null) {
      out.append(' ').append(GoNames.localIdent(bt.getLabel().toString()));
    }
    out.append('\n');
  }

  private void emitContinue(ContinueTree ct, int indent, StringBuilder out, EmitContext ctx) {
    if (ct.getLabel() == null && ctx.inTryWrapper && ctx.loopDepthInTryIife == 0 && ctx.tryFlowVar != null) {
      tabs(out, indent).append(ctx.tryFlowVar).append(" = 2\n");
      tabs(out, indent).append("return\n");
      return;
    }
    tabs(out, indent).append("continue");
    if (ct.getLabel() != null) {
      out.append(' ').append(GoNames.localIdent(ct.getLabel().toString()));
    }
    out.append('\n');
  }

  private void emitReturn(
      ReturnTree rt, TreePath path, int indent, StringBuilder out, EmitContext ctx) {
    if (ctx.inTryWrapper && ctx.switchExprTargetStack.isEmpty()) {
      if (rt.getExpression() != null && ctx.expectedReturnType != null) {
        TreePath tryRetExprPath = unwrapParensPath(new TreePath(path, rt.getExpression()));
        if (tryRetExprPath.getLeaf() instanceof ConditionalExpressionTree ce) {
          emitUnwrappedConditional(
              ce,
              tryRetExprPath,
              indent,
              false,
              ctx.expectedReturnType,
              ctx.returnWidenedToAny,
              (val, branchIndent) ->
                  tabs(out, branchIndent)
                      .append(ctx.tryRetValVar)
                      .append(" = ")
                      .append(val)
                      .append('\n'),
              out,
              ctx);
          tabs(out, indent).append(ctx.tryDidReturnVar).append(" = true\n");
          tabs(out, indent).append("return\n");
          return;
        }
        String val =
            exprEmitter.emitAs(
                rt.getExpression(), path, ctx.expectedReturnType, ctx.returnWidenedToAny, ctx);
        tabs(out, indent).append(ctx.tryRetValVar).append(" = ").append(val).append('\n');
      }
      tabs(out, indent).append(ctx.tryDidReturnVar).append(" = true\n");
      tabs(out, indent).append("return\n");
      return;
    }

    if (rt.getExpression() == null || ctx.expectedReturnType == null) {
      tabs(out, indent).append("return\n");
    } else {
      TreePath unwrappedRetPath = unwrapParensPath(new TreePath(path, rt.getExpression()));
      if (unwrappedRetPath.getLeaf() instanceof ConditionalExpressionTree ce) {
        emitUnwrappedConditional(
            ce,
            unwrappedRetPath,
            indent,
            false,
            ctx.expectedReturnType,
            ctx.returnWidenedToAny,
            (val, branchIndent) -> tabs(out, branchIndent).append("return ").append(val).append('\n'),
            out,
            ctx);
        return;
      }
      if (!ctx.inTryWrapper && ctx.switchExprTargetStack.isEmpty()) {
        if (unwrappedRetPath.getLeaf() instanceof SwitchExpressionTree set
            && !containsDirectTryTree(set)) {
          ctx.switchExprTargetStack.push(
              new EmitContext.SwitchExprTarget(null, ctx.expectedReturnType, ctx.returnWidenedToAny));
          emitSwitchCommon(
              set.getExpression(),
              set.getCases(),
              unwrappedRetPath,
              indent,
              true,
              ctx.expectedReturnType,
              null,
              out,
              ctx);
          ctx.switchExprTargetStack.pop();
          tabs(out, indent).append("panic(\"unreachable switch expression\")\n");
          return;
        }
      }
      if (ctx.expectedReturnType.getKind() == TypeKind.BOOLEAN && !ctx.returnWidenedToAny) {
        String patBool =
            tryEmitConjunctivePatternCond(
                unwrappedRetPath, indent, false, out, ctx);
        if (patBool != null) {
          tabs(out, indent).append("return ").append(patBool).append('\n');
          return;
        }
      }
      String val =
          exprEmitter.emitAs(
              rt.getExpression(), path, ctx.expectedReturnType, ctx.returnWidenedToAny, ctx);
      tabs(out, indent).append("return ").append(val).append('\n');
    }
  }

  private void emitYield(
      YieldTree yt, TreePath path, int indent, StringBuilder out, EmitContext ctx) {
    EmitContext.SwitchExprTarget sink = ctx.switchExprTargetStack.peek();
    String val = exprEmitter.emitAs(yt.getValue(), path, sink.targetType(), sink.forceAny(), ctx);
    if (sink.assignVar() != null) {
      tabs(out, indent).append(sink.assignVar()).append(" = ").append(val).append('\n');
      if (ctx.inTryWrapper) {
        if (ctx.loopDepthInTryIife == 0 && ctx.switchDepthInTryIife == 0 && ctx.tryFlowVar != null) {
          tabs(out, indent).append(ctx.tryFlowVar).append(" = 1\n");
          tabs(out, indent).append("return\n");
        } else {
          tabs(out, indent).append("break\n");
        }
      }
    } else if (ctx.inTryWrapper) {
      if (ctx.tryRetValVar != null) {
        tabs(out, indent).append(ctx.tryRetValVar).append(" = ").append(val).append('\n');
      }
      tabs(out, indent).append(ctx.tryDidReturnVar).append(" = true\n");
      tabs(out, indent).append("return\n");
    } else {
      tabs(out, indent).append("return ").append(val).append('\n');
    }
  }

  private void emitThrow(
      ThrowTree tt, TreePath path, int indent, StringBuilder out, EmitContext ctx) {
    ExpressionTree exTree = tt.getExpression();
    String ex = exprEmitter.emit(exTree, path, ctx);
    TypeMirror exTm = exprEmitter.typeOf(new TreePath(path, exTree));
    tabs(out, indent)
        .append("panic(")
        .append(exprEmitter.emitCheckNotNull(exTree, ex, exTm, ctx))
        .append(")\n");
  }

  private void emitAssert(
      AssertTree at, TreePath path, int indent, StringBuilder out, EmitContext ctx) {
    String cond =
        exprEmitter.emitAs(
            at.getCondition(),
            path,
            resolver.types.getPrimitiveType(TypeKind.BOOLEAN),
            false,
            ctx);
    String detail =
        at.getDetail() != null
            ? exprEmitter.emitAs(
                at.getDetail(), path, resolver.getObjectElement().asType(), true, ctx)
            : "nil";
    String jl = jlPrefix(ctx);
    tabs(out, indent)
        .append("if !(")
        .append(cond)
        .append(") { panic(")
        .append(jl)
        .append("NewAssertionError(")
        .append(detail)
        .append(")) }\n");
  }

  private void emitSynchronized(
      SynchronizedTree st, TreePath path, int indent, StringBuilder out, EmitContext ctx) {
    ExpressionTree unwrappedLock = st.getExpression();
    while (unwrappedLock instanceof ParenthesizedTree pt) {
      unwrappedLock = pt.getExpression();
    }
    if (ctx.currentType != null
        && !ctx.inStaticContext
        && resolver.hasFastInstanceSyncMutex(ctx.currentType)
        && unwrappedLock instanceof IdentifierTree id
        && id.getName().contentEquals("this")) {
      boolean wrapIife =
          !ctx.inTryWrapper
              && ctx.switchExprTargetStack.isEmpty()
              && !hasDisallowedJumpForDirectDefer(st.getBlock(), false);
      tabs(out, indent).append(wrapIife ? "func() {\n" : "{\n");
      tabs(out, indent + 1)
          .append("if this_ == nil { ")
          .append(rtPrefix(ctx))
          .append("ThrowNPE() }\n");
      tabs(out, indent + 1).append("this_._syncMu.Lock()\n");
      tabs(out, indent + 1).append("defer this_._syncMu.Unlock()\n");
      out.append(emitBlockStatements(st.getBlock(), new TreePath(path, st.getBlock()), indent + 1, ctx));
      tabs(out, indent).append(wrapIife ? "}()\n" : "}\n");
      return;
    }
    String lockExpr =
        exprEmitter.emitAs(
            st.getExpression(), path, resolver.getObjectElement().asType(), true, ctx);
    String unlockTmp = ctx.nextTemp("unlock");
    tabs(out, indent).append("{\n");
    tabs(out, indent + 1)
        .append(unlockTmp)
        .append(" := ")
        .append(rtPrefix(ctx))
        .append("LockMonitor(")
        .append(lockExpr)
        .append(")\n");
    tabs(out, indent + 1).append("defer ").append(unlockTmp).append("()\n");
    out.append(emitBlockStatements(st.getBlock(), new TreePath(path, st.getBlock()), indent + 1, ctx));
    tabs(out, indent).append("}\n");
  }

  // --- Switch Statements & Switch Expressions ---

  public String emitSwitchExpression(
      SwitchExpressionTree set, TreePath path, EmitContext ctx) {
    TypeMirror resTm = exprEmitter.typeOf(path);
    String resGo = resolver.goType(resTm, ctx.imports);
    StringBuilder sb = new StringBuilder();
    sb.append("func() ").append(resGo).append(" {\n");
    boolean savedInTry = ctx.inTryWrapper;
    TypeMirror savedExpRet = ctx.expectedReturnType;
    boolean savedWidened = ctx.returnWidenedToAny;
    ctx.inTryWrapper = false;
    ctx.expectedReturnType = resTm;
    ctx.returnWidenedToAny = false;
    ctx.switchExprTargetStack.push(new EmitContext.SwitchExprTarget(null, resTm, false));
    emitSwitchCommon(set.getExpression(), set.getCases(), path, 1, true, resTm, null, sb, ctx);
    ctx.switchExprTargetStack.pop();
    ctx.inTryWrapper = savedInTry;
    ctx.expectedReturnType = savedExpRet;
    ctx.returnWidenedToAny = savedWidened;
    sb.append("\tpanic(\"unreachable switch expression\")\n}()");
    return sb.toString();
  }

  private void emitSwitchCommon(
      ExpressionTree selectorTree,
      List<? extends CaseTree> cases,
      TreePath switchPath,
      int indent,
      boolean isExpr,
      TypeMirror exprTargetTm,
      String precomputedEnumOrdGo,
      StringBuilder out,
      EmitContext ctx) {
    if (cases.isEmpty()) {
      tabs(out, indent).append("_ = ").append(exprEmitter.emit(selectorTree, switchPath, ctx)).append('\n');
      return;
    }

    boolean isRuleSwitch =
        cases.stream().anyMatch(c -> c.getCaseKind() == CaseTree.CaseKind.RULE);
    boolean isPatternSwitch = hasPatternOrNullCase(cases);

    // Strip outer parens on selector
    ExpressionTree rawSel = selectorTree;
    TreePath selParent = switchPath;
    while (rawSel instanceof ParenthesizedTree pt) {
      selParent = new TreePath(selParent, rawSel);
      rawSel = pt.getExpression();
    }
    TreePath selPath = new TreePath(selParent, rawSel);
    TypeMirror selTm = exprEmitter.typeOf(selPath);
    TypeMirror switchSelTm = computeSwitchSelectorTargetType(selTm);

    // Fast path: compile-time integer `.Ordinal()` switch for Java enums (with or without `case null`)!
    TypeElement enumSwitchTe = null;
    Map<VariableElement, Integer> enumOrdinals = null;
    boolean enumHasNullCase = false;
    if (switchSelTm instanceof DeclaredType dt
        && dt.asElement() instanceof TypeElement ete
        && ete.getKind() == ElementKind.ENUM) {
      Map<VariableElement, Integer> ordMap = new HashMap<>();
      int ord = 0;
      for (Element enc : ete.getEnclosedElements()) {
        if (enc instanceof VariableElement eve && eve.getKind() == ElementKind.ENUM_CONSTANT) {
          ordMap.put(eve, ord++);
        }
      }
      boolean validEnumSwitch = true;
      for (CaseTree ct : cases) {
        if (ct.getGuard() != null) {
          validEnumSwitch = false;
          break;
        }
        for (CaseLabelTree lbl : ct.getLabels()) {
          if (lbl instanceof DefaultCaseLabelTree) {
            continue;
          }
          if (lbl instanceof ConstantCaseLabelTree ccl) {
            ExpressionTree ce = ccl.getConstantExpression();
            if (ce.getKind() == Tree.Kind.NULL_LITERAL) {
              enumHasNullCase = true;
              continue;
            }
            TreePath cePath = new TreePath(new TreePath(new TreePath(switchPath, ct), lbl), ce);
            Element cel = resolver.trees.getElement(cePath);
            if (cel instanceof VariableElement ve && ordMap.containsKey(ve)) {
              continue;
            }
          }
          validEnumSwitch = false;
          break;
        }
        if (!validEnumSwitch) {
          break;
        }
      }
      if (validEnumSwitch) {
        enumSwitchTe = ete;
        enumOrdinals = ordMap;
      }
    }

    if (isPatternSwitch && enumOrdinals == null) {
      emitPatternSwitch(selPath, selTm, cases, switchPath, indent, isExpr, exprTargetTm, isRuleSwitch, out, ctx);
      return;
    }

    // Standard value/enum/string/int switch
    // For classic `:` switch, hoist top-level local variable declarations across cases!
    List<VariableElement> hoistedVars = new ArrayList<>();
    if (!isRuleSwitch) {
      for (CaseTree ct : cases) {
        if (ct.getStatements() != null) {
          for (StatementTree s : ct.getStatements()) {
            if (s instanceof VariableTree vt) {
              TreePath vp = new TreePath(new TreePath(switchPath, ct), vt);
              if (resolver.trees.getElement(vp) instanceof VariableElement ve) {
                hoistedVars.add(ve);
              }
            }
          }
        }
      }
    }

    int swIndent = indent;
    if (!hoistedVars.isEmpty()) {
      tabs(out, indent).append("{\n");
      swIndent = indent + 1;
      for (VariableElement ve : hoistedVars) {
        String vn = ctx.registerLocalVar(ve);
        String vtGo = resolver.goType(ve.asType(), ctx.imports);
        tabs(out, swIndent)
            .append("var ")
            .append(vn)
            .append(' ')
            .append(vtGo)
            .append(" = ")
            .append(resolver.zeroValue(ve.asType(), ctx.imports))
            .append('\n');
        tabs(out, swIndent).append("_ = ").append(vn).append('\n');
      }
    }

    // Evaluate selector (unboxing boxed primitives/enums/strings with NPE check per JLS §14.11.1!)
    String selGo;
    if (precomputedEnumOrdGo != null && enumOrdinals != null) {
      selGo = precomputedEnumOrdGo;
    } else {
      String nodeTokenRecv =
          (enumOrdinals != null
                  && enumSwitchTe != null
                  && enumSwitchTe
                      .getQualifiedName()
                      .contentEquals("com.google.javascript.rhino.Token"))
              ? exprEmitter.tryEmitNodeForTokenRead(selPath, ctx)
              : null;
      if (nodeTokenRecv != null) {
        selGo =
            enumHasNullCase
                ? (nodeTokenRecv + ".TokenOrd_()")
                : (nodeTokenRecv + ".TokenSwitchOrd_()");
      } else {
        selGo =
            exprEmitter.convert(exprEmitter.emitPath(selPath, ctx), selTm, switchSelTm, false, ctx);
        if (enumOrdinals != null) {
          if (resolver.isValueEnum(enumSwitchTe)) {
            if (enumHasNullCase) {
              selGo = "(int32(" + selGo + ") - 1)";
            } else {
              selGo = "(" + selGo + ").Ordinal()";
            }
          } else if (enumHasNullCase) {
            selGo = rtPrefix(ctx) + "EnumSwitchOrd(" + selGo + ")";
          } else {
            selGo = rtPrefix(ctx) + "CheckNotNullPtr(" + selGo + ").Ordinal()";
          }
        } else if (!switchSelTm.getKind().isPrimitive()) {
          selGo = exprEmitter.emitCheckNotNull(rawSel, selGo, switchSelTm, ctx);
        }
      }
    }

    if (ctx.inTryWrapper) ctx.switchDepthInTryIife++; // `break` inside switch targets switch, not outer loop!

    tabs(out, swIndent).append("switch ").append(selGo).append(" {\n");

    if (isRuleSwitch) {
      for (CaseTree ct : cases) {
        TreePath casePath = new TreePath(switchPath, ct);
        emitCaseHeader(List.of(ct), switchPath, switchSelTm, enumSwitchTe, enumOrdinals, swIndent, out, ctx);
        emitRuleCaseBody(ct, casePath, swIndent + 1, isExpr, exprTargetTm, out, ctx);
      }
    } else {
      // Group consecutive empty classic `:` cases (`case A: case B: stmts`)
      int i = 0;
      while (i < cases.size()) {
        List<CaseTree> grouped = new ArrayList<>();
        while (i < cases.size() - 1
            && (cases.get(i).getStatements() == null || cases.get(i).getStatements().isEmpty())
            && !isDefaultCase(cases.get(i))
            && !isDefaultCase(cases.get(i + 1))) {
          grouped.add(cases.get(i));
          i++;
        }
        CaseTree active = cases.get(i);
        grouped.add(active);
        boolean isLastGroup = (i == cases.size() - 1);
        i++;

        emitCaseHeader(grouped, switchPath, switchSelTm, enumSwitchTe, enumOrdinals, swIndent, out, ctx);
        List<? extends StatementTree> stmts = active.getStatements();
        TreePath activePath = new TreePath(switchPath, active);
        if (stmts != null) {
          emitClassicCaseStatements(stmts, activePath, swIndent + 1, hoistedVars, out, ctx);
        }
        if (!isLastGroup && fallsThrough(stmts)) {
          tabs(out, swIndent + 1).append("fallthrough\n");
        }
      }
    }

    tabs(out, swIndent).append("}\n");
    if (ctx.inTryWrapper) ctx.switchDepthInTryIife--;

    if (!hoistedVars.isEmpty()) {
      tabs(out, indent).append("}\n");
    }
  }

  private void emitPatternSwitch(
      TreePath selPath,
      TypeMirror selTm,
      List<? extends CaseTree> cases,
      TreePath switchPath,
      int indent,
      boolean isExpr,
      TypeMirror exprTargetTm,
      boolean isRuleSwitch,
      StringBuilder out,
      EmitContext ctx) {
    String selVar = ctx.nextTemp("sw");
    String selAny =
        exprEmitter.convert(
            exprEmitter.emitPath(selPath, ctx),
            selTm,
            resolver.getObjectElement().asType(),
            true,
            ctx);

    boolean hasExplicitNull = false;
    for (CaseTree ct : cases) {
      for (CaseLabelTree lbl : ct.getLabels()) {
        if (lbl instanceof ConstantCaseLabelTree ccl
            && ccl.getConstantExpression().getKind() == Tree.Kind.NULL_LITERAL) {
          hasExplicitNull = true;
        }
      }
    }

    tabs(out, indent).append("{\n");
    tabs(out, indent + 1).append(selVar).append(" := ").append(selAny).append('\n');
    if (!hasExplicitNull) {
      tabs(out, indent + 1)
          .append("if ")
          .append(selVar)
          .append(" == nil { panic(")
          .append(jlPrefix(ctx))
          .append("NewNullPointerException()) }\n");
    }

    // Predeclare any pattern binding variables from case labels
    for (CaseTree ct : cases) {
      TreePath cp = new TreePath(switchPath, ct);
      for (CaseLabelTree lbl : ct.getLabels()) {
        if (lbl instanceof PatternCaseLabelTree pcl && pcl.getPattern() instanceof BindingPatternTree bpt) {
          TreePath vp =
              new TreePath(new TreePath(new TreePath(cp, lbl), pcl.getPattern()), bpt.getVariable());
          if (resolver.trees.getElement(vp) instanceof VariableElement ve) {
            String vn = ctx.registerLocalVar(ve);
            tabs(out, indent + 1)
                .append("var ")
                .append(vn)
                .append(' ')
                .append(resolver.goType(ve.asType(), ctx.imports))
                .append('\n');
            tabs(out, indent + 1).append("_ = ").append(vn).append('\n');
          }
        }
      }
    }

    if (ctx.inTryWrapper) ctx.switchDepthInTryIife++;
    tabs(out, indent + 1).append("switch {\n");
    for (CaseTree ct : cases) {
      TreePath cp = new TreePath(switchPath, ct);
      if (isDefaultCase(ct) && ct.getLabels().size() == 1) {
        tabs(out, indent + 1).append("default:\n");
      } else {
        List<String> conds = new ArrayList<>();
        boolean hasDefaultAmongLabels = false;
        for (CaseLabelTree lbl : ct.getLabels()) {
          if (lbl instanceof DefaultCaseLabelTree) {
            hasDefaultAmongLabels = true;
          } else if (lbl instanceof ConstantCaseLabelTree ccl) {
            ExpressionTree ce = ccl.getConstantExpression();
            if (ce.getKind() == Tree.Kind.NULL_LITERAL) {
              conds.add("(" + selVar + " == nil)");
            } else {
              String cv =
                  exprEmitter.emitAs(
                      ce, new TreePath(cp, lbl), resolver.getObjectElement().asType(), true, ctx);
              conds.add(jlPrefix(ctx) + "ObjectsEquals(" + selVar + ", " + cv + ")");
            }
          } else if (lbl instanceof PatternCaseLabelTree pcl
              && pcl.getPattern() instanceof BindingPatternTree bpt) {
            TreePath vp =
                new TreePath(new TreePath(new TreePath(cp, lbl), pcl.getPattern()), bpt.getVariable());
            VariableElement ve = (VariableElement) resolver.trees.getElement(vp);
            String vn = ctx.registerLocalVar(ve);
            TypeMirror pTm = ve.asType();
            String rt = rtPrefix(ctx);
            String bindCond =
                resolver.isStringType(pTm)
                    ? (rt + "InstanceOfStringBind(" + selVar + ", &" + vn + ")")
                    : (rt + "InstanceOfBind[" + resolver.goType(pTm, ctx.imports) + "](" + selVar + ", &" + vn + ")");
            conds.add(bindCond);
          }
        }
        if (hasDefaultAmongLabels) {
          tabs(out, indent + 1).append("default:\n");
        } else {
          String fullCond = String.join(" || ", conds);
          if (ct.getGuard() != null) {
            String guard =
                exprEmitter.emitAs(
                    ct.getGuard(),
                    cp,
                    resolver.types.getPrimitiveType(TypeKind.BOOLEAN),
                    false,
                    ctx);
            fullCond = "(" + fullCond + ") && (" + guard + ")";
          }
          tabs(out, indent + 1).append("case ").append(fullCond).append(":\n");
        }
      }

      if (isRuleSwitch) {
        emitRuleCaseBody(ct, cp, indent + 2, isExpr, exprTargetTm, out, ctx);
      } else if (ct.getStatements() != null) {
        emitClassicCaseStatements(ct.getStatements(), cp, indent + 2, List.of(), out, ctx);
      }
    }
    tabs(out, indent + 1).append("}\n");
    if (ctx.inTryWrapper) ctx.switchDepthInTryIife--;
    tabs(out, indent).append("}\n");
  }

  private TypeMirror computeSwitchSelectorTargetType(TypeMirror selTm) {
    TypeMirror er = resolver.safeErasure(selTm);
    if (er instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      return switch (te.getQualifiedName().toString()) {
        case "java.lang.Character" -> resolver.types.getPrimitiveType(TypeKind.CHAR);
        case "java.lang.Byte" -> resolver.types.getPrimitiveType(TypeKind.BYTE);
        case "java.lang.Short" -> resolver.types.getPrimitiveType(TypeKind.SHORT);
        case "java.lang.Integer" -> resolver.types.getPrimitiveType(TypeKind.INT);
        default -> er;
      };
    }
    return er;
  }

  private boolean hasPatternOrNullCase(List<? extends CaseTree> cases) {
    for (CaseTree ct : cases) {
      if (ct.getGuard() != null) return true;
      for (CaseLabelTree lbl : ct.getLabels()) {
        if (lbl instanceof PatternCaseLabelTree) return true;
        if (lbl instanceof ConstantCaseLabelTree ccl
            && ccl.getConstantExpression().getKind() == Tree.Kind.NULL_LITERAL) {
          return true;
        }
      }
    }
    return false;
  }

  private boolean isDefaultCase(CaseTree ct) {
    if (ct.getLabels().isEmpty()) {
      return true;
    }
    for (CaseLabelTree lbl : ct.getLabels()) {
      if (lbl instanceof DefaultCaseLabelTree) {
        return true;
      }
    }
    return false;
  }

  private void emitCaseHeader(
      List<CaseTree> grouped,
      TreePath switchPath,
      TypeMirror switchSelTm,
      TypeElement enumSwitchTe,
      Map<VariableElement, Integer> enumOrdinals,
      int indent,
      StringBuilder out,
      EmitContext ctx) {
    boolean hasDefault = false;
    List<String> exprs = new ArrayList<>();
    for (CaseTree ct : grouped) {
      TreePath cp = new TreePath(switchPath, ct);
      if (ct.getLabels().isEmpty()) {
        hasDefault = true;
      }
      for (CaseLabelTree lbl : ct.getLabels()) {
        if (lbl instanceof DefaultCaseLabelTree) {
          hasDefault = true;
        } else if (lbl instanceof ConstantCaseLabelTree ccl) {
          ExpressionTree ce = ccl.getConstantExpression();
          if (enumOrdinals != null) {
            if (ce.getKind() == Tree.Kind.NULL_LITERAL) {
              exprs.add("int32(-1)");
            } else {
              TreePath cePath = new TreePath(new TreePath(cp, lbl), ce);
              VariableElement ve = (VariableElement) resolver.trees.getElement(cePath);
              int ord = enumOrdinals.get(ve);
              exprs.add("int32(" + ord + ") /* " + enumSwitchTe.getSimpleName() + "." + ve.getSimpleName() + " */");
            }
          } else {
            exprs.add(
                exprEmitter.emitAs(
                    ce, new TreePath(cp, lbl), switchSelTm, false, ctx));
          }
        }
      }
    }
    if (hasDefault) {
      tabs(out, indent).append("default:\n");
    } else {
      tabs(out, indent).append("case ").append(String.join(", ", exprs)).append(":\n");
    }
  }

  private void emitRuleCaseBody(
      CaseTree ct,
      TreePath casePath,
      int indent,
      boolean isExpr,
      TypeMirror exprTargetTm,
      StringBuilder out,
      EmitContext ctx) {
    Tree body = ct.getBody();
    if (body == null) {
      return;
    }
    TreePath bodyPath = new TreePath(casePath, body);
    if (body instanceof ExpressionTree et) {
      if (isExpr) {
        EmitContext.SwitchExprTarget sink = ctx.switchExprTargetStack.peek();
        TypeMirror targetTm = sink != null ? sink.targetType() : exprTargetTm;
        boolean forceAny = sink != null && sink.forceAny();
        String val = exprEmitter.emitAs(et, casePath, targetTm, forceAny, ctx);
        if (sink != null && sink.assignVar() != null) {
          tabs(out, indent).append(sink.assignVar()).append(" = ").append(val).append('\n');
        } else {
          tabs(out, indent).append("return ").append(val).append('\n');
        }
      } else {
        tabs(out, indent).append(emitExprAsStmt(et, casePath, ctx)).append('\n');
      }
    } else if (body instanceof StatementTree st) {
      emitAsBody(st, bodyPath, indent, out, ctx);
    }
  }

  private void emitClassicCaseStatements(
      List<? extends StatementTree> stmts,
      TreePath casePath,
      int indent,
      List<VariableElement> hoistedVars,
      StringBuilder out,
      EmitContext ctx) {
    Set<VariableElement> hoistedSet = new HashSet<>(hoistedVars);
    Set<VariableElement> preds = new HashSet<>();
    for (StatementTree s : stmts) {
      TreePath sp = new TreePath(casePath, s);
      emitPatternVarPredecls(sp, indent, preds, out, ctx);
      if (s instanceof VariableTree vt
          && resolver.trees.getElement(sp) instanceof VariableElement ve
          && hoistedSet.contains(ve)) {
        if (vt.getInitializer() != null) {
          String vn = ctx.lookupLocalVar(ve);
          String rhs = exprEmitter.emitAs(vt.getInitializer(), sp, ve.asType(), false, ctx);
          tabs(out, indent).append(vn).append(" = ").append(rhs).append('\n');
        }
        continue;
      }
      // Strip trailing unlabeled `break;` because Go switch does not fall through by default!
      if (s == stmts.get(stmts.size() - 1)
          && s instanceof BreakTree bt
          && bt.getLabel() == null) {
        continue;
      }
      emitStmt(s, sp, indent, null, out, ctx);
    }
  }

  private boolean fallsThrough(List<? extends StatementTree> stmts) {
    if (stmts == null || stmts.isEmpty()) {
      return true;
    }
    StatementTree last = stmts.get(stmts.size() - 1);
    return !stmtAlwaysExitsSwitchOrMethod(last);
  }

  private boolean stmtAlwaysExitsSwitchOrMethod(StatementTree s) {
    return switch (s.getKind()) {
      case BREAK, RETURN, THROW, CONTINUE, YIELD -> true;
      case BLOCK -> {
        List<? extends StatementTree> list = ((BlockTree) s).getStatements();
        yield !list.isEmpty() && stmtAlwaysExitsSwitchOrMethod(list.get(list.size() - 1));
      }
      case IF -> {
        IfTree it = (IfTree) s;
        yield it.getElseStatement() != null
            && stmtAlwaysExitsSwitchOrMethod(it.getThenStatement())
            && stmtAlwaysExitsSwitchOrMethod(it.getElseStatement());
      }
      default -> false;
    };
  }

  // --- Try / Catch / Finally / Try-With-Resources ---

  private void emitTry(TryTree tt, TreePath path, int indent, StringBuilder out, EmitContext ctx) {
    if (canEmitDirectDeferFinally(tt, path, ctx)) {
      tabs(out, indent).append("defer func() {\n");
      out.append(
          emitBlockStatements(
              tt.getFinallyBlock(), new TreePath(path, tt.getFinallyBlock()), indent + 1, ctx));
      tabs(out, indent).append("}()\n");
      tabs(out, indent).append("{\n");
      out.append(
          emitBlockStatements(tt.getBlock(), new TreePath(path, tt.getBlock()), indent + 1, ctx));
      tabs(out, indent).append("}\n");
      return;
    }

    // Analyze if `try` or `catch` contains `return` or outer `break`/`continue`
    boolean[] hasReturn = {false};
    boolean[] hasOuterBreak = {false};
    boolean[] hasOuterContinue = {false};

    // d[0] = loopDepth, d[1] = switchDepth inside the `try`/`catch`
    TreePathScanner<Void, int[]> jumpScanner =
        new TreePathScanner<>() {
          @Override
          public Void visitLambdaExpression(LambdaExpressionTree n, int[] d) {
            return null;
          }

          @Override
          public Void visitClass(ClassTree n, int[] d) {
            return null;
          }

          @Override
          public Void visitSwitchExpression(SwitchExpressionTree n, int[] d) {
            return null;
          }

          @Override
          public Void visitWhileLoop(WhileLoopTree n, int[] d) {
            return super.visitWhileLoop(n, new int[] {d[0] + 1, d[1]});
          }

          @Override
          public Void visitDoWhileLoop(DoWhileLoopTree n, int[] d) {
            return super.visitDoWhileLoop(n, new int[] {d[0] + 1, d[1]});
          }

          @Override
          public Void visitForLoop(ForLoopTree n, int[] d) {
            return super.visitForLoop(n, new int[] {d[0] + 1, d[1]});
          }

          @Override
          public Void visitEnhancedForLoop(EnhancedForLoopTree n, int[] d) {
            return super.visitEnhancedForLoop(n, new int[] {d[0] + 1, d[1]});
          }

          @Override
          public Void visitSwitch(SwitchTree n, int[] d) {
            return super.visitSwitch(n, new int[] {d[0], d[1] + 1});
          }

          @Override
          public Void visitReturn(ReturnTree n, int[] d) {
            hasReturn[0] = true;
            return null;
          }

          @Override
          public Void visitYield(YieldTree n, int[] d) {
            if (!ctx.switchExprTargetStack.isEmpty()) {
              if (ctx.switchExprTargetStack.peek().assignVar() == null) {
                hasReturn[0] = true;
              } else if (d[0] == 0 && d[1] == 0) {
                hasOuterBreak[0] = true;
              }
            }
            return null;
          }

          @Override
          public Void visitBreak(BreakTree n, int[] d) {
            if (n.getLabel() == null && d[0] == 0 && d[1] == 0) {
              hasOuterBreak[0] = true;
            }
            return null;
          }

          @Override
          public Void visitContinue(ContinueTree n, int[] d) {
            if (n.getLabel() == null && d[0] == 0) {
              hasOuterContinue[0] = true;
            }
            return null;
          }
        };

    int[] zeroDepth = {0, 0};
    jumpScanner.scan(new TreePath(path, tt.getBlock()), zeroDepth);
    for (CatchTree ct : tt.getCatches()) {
      jumpScanner.scan(new TreePath(path, ct), zeroDepth);
    }
    boolean hasOuterJump = hasOuterBreak[0] || hasOuterContinue[0];

    // Save outer try-wrapper state in case of nested `try` blocks!
    boolean prevInTry = ctx.inTryWrapper;
    String prevRetVal = ctx.tryRetValVar;
    String prevDidRet = ctx.tryDidReturnVar;
    String prevFlow = ctx.tryFlowVar;
    int prevLoopDepth = ctx.loopDepthInTryIife;
    int prevSwitchDepth = ctx.switchDepthInTryIife;

    String retValVar = null;
    String didRetVar = null;
    String flowVar = null;

    if (hasReturn[0]) {
      didRetVar = ctx.nextTemp("didRet");
      tabs(out, indent).append("var ").append(didRetVar).append(" bool\n");
      if (ctx.expectedReturnType != null) {
        retValVar = ctx.nextTemp("retVal");
        tabs(out, indent)
            .append("var ")
            .append(retValVar)
            .append(' ')
            .append(resolver.goType(ctx.expectedReturnType, ctx.imports))
            .append('\n');
      }
    }
    if (hasOuterJump) {
      flowVar = ctx.nextTemp("tryFlow");
      tabs(out, indent).append("var ").append(flowVar).append(" int\n");
    }

    ctx.inTryWrapper = true;
    ctx.tryRetValVar = retValVar;
    ctx.tryDidReturnVar = didRetVar;
    ctx.tryFlowVar = flowVar;
    ctx.loopDepthInTryIife = 0;
    ctx.switchDepthInTryIife = 0;

    tabs(out, indent).append("func() {\n");

    // 1. `finally` block runs unconditionally on exit of outer func() via `defer`
    if (tt.getFinallyBlock() != null) {
      tabs(out, indent + 1).append("defer func() {\n");
      // Temporarily restore outer try context inside `finally` if `finally` itself returns
      boolean savedInTry = ctx.inTryWrapper;
      ctx.inTryWrapper = hasReturn[0];
      out.append(
          emitBlockStatements(
              tt.getFinallyBlock(), new TreePath(path, tt.getFinallyBlock()), indent + 2, ctx));
      ctx.inTryWrapper = savedInTry;
      tabs(out, indent + 1).append("}()\n");
    }

    // 2. `catch` recovery handler via `defer` (runs BEFORE `finally` because LIFO defer order!)
    if (!tt.getCatches().isEmpty()) {
      String recVar = ctx.nextTemp("rec");
      String exVar = ctx.nextTemp("jex");
      String jl = jlPrefix(ctx);
      String rt = rtPrefix(ctx);
      tabs(out, indent + 1).append("defer func() {\n");
      tabs(out, indent + 2).append("if ").append(recVar).append(" := recover(); ").append(recVar).append(" != nil {\n");
      tabs(out, indent + 3)
          .append(exVar)
          .append(" := ")
          .append(jl)
          .append("RecoverException(")
          .append(recVar)
          .append(")\n");
      tabs(out, indent + 3).append("_ = ").append(exVar).append('\n');

      for (int i = 0; i < tt.getCatches().size(); i++) {
        CatchTree ct = tt.getCatches().get(i);
        TreePath cp = new TreePath(path, ct);
        TreePath paramPath = new TreePath(cp, ct.getParameter());
        VariableElement catchVe = (VariableElement) resolver.trees.getElement(paramPath);
        String catchVarName = ctx.registerLocalVar(catchVe);
        Tree paramTypeTree = ct.getParameter().getType();

        tabs(out, indent + 3).append(i == 0 ? "if " : "} else if ");
        if (paramTypeTree instanceof UnionTypeTree utt) {
          List<String> checks = new ArrayList<>();
          for (Tree alt : utt.getTypeAlternatives()) {
            TypeMirror altTm = exprEmitter.typeOf(new TreePath(paramPath, alt));
            String altGo = resolver.goType(altTm, ctx.imports);
            checks.add(rt + "InstanceOf[" + altGo + "](" + exVar + ")");
          }
          out.append(String.join(" || ", checks)).append(" {\n");
          tabs(out, indent + 4).append(catchVarName).append(" := ").append(exVar).append('\n');
          tabs(out, indent + 4).append("_ = ").append(catchVarName).append('\n');
        } else {
          TypeMirror catchTm = catchVe.asType();
          String catchGo = resolver.goType(catchTm, ctx.imports);
          out.append(rt).append("InstanceOf[").append(catchGo).append("](").append(exVar).append(") {\n");
          tabs(out, indent + 4)
              .append(catchVarName)
              .append(" := ")
              .append(rt)
              .append("AsIface[")
              .append(catchGo)
              .append("](")
              .append(exVar)
              .append(")\n");
          tabs(out, indent + 4).append("_ = ").append(catchVarName).append('\n');
        }
        out.append(emitBlockStatements(ct.getBlock(), new TreePath(cp, ct.getBlock()), indent + 4, ctx));
      }
      tabs(out, indent + 3).append("} else {\n");
      tabs(out, indent + 4).append("panic(").append(recVar).append(")\n");
      tabs(out, indent + 3).append("}\n");
      tabs(out, indent + 2).append("}\n");
      tabs(out, indent + 1).append("}()\n");
    }

    // 3. `try-with-resources` declarations + LIFO `defer CloseResource(...)`
    for (Tree res : tt.getResources()) {
      TreePath resPath = new TreePath(path, res);
      if (res instanceof VariableTree vt) {
        emitVariableDecl(vt, resPath, indent + 1, out, ctx);
        VariableElement ve = (VariableElement) resolver.trees.getElement(resPath);
        String vn = ctx.lookupLocalVar(ve);
        tabs(out, indent + 1)
            .append("defer ")
            .append(rtPrefix(ctx))
            .append("CloseResource(")
            .append(vn)
            .append(")\n");
      } else if (res instanceof ExpressionTree et) {
        String rExpr = exprEmitter.emit(et, path, ctx);
        tabs(out, indent + 1)
            .append("defer ")
            .append(rtPrefix(ctx))
            .append("CloseResource(")
            .append(rExpr)
            .append(")\n");
      }
    }

    // 4. `try` block body
    out.append(emitBlockStatements(tt.getBlock(), new TreePath(path, tt.getBlock()), indent + 1, ctx));
    tabs(out, indent).append("}()\n");

    // Restore outer try context
    ctx.inTryWrapper = prevInTry;
    ctx.tryRetValVar = prevRetVal;
    ctx.tryDidReturnVar = prevDidRet;
    ctx.tryFlowVar = prevFlow;
    ctx.loopDepthInTryIife = prevLoopDepth;
    ctx.switchDepthInTryIife = prevSwitchDepth;

    // Propagate return or break/continue out of try wrapper!
    if (hasReturn[0]) {
      if (ctx.inTryWrapper) {
        tabs(out, indent).append("if ").append(didRetVar).append(" {\n");
        if (retValVar != null && ctx.tryRetValVar != null) {
          tabs(out, indent + 1).append(ctx.tryRetValVar).append(" = ").append(retValVar).append('\n');
        }
        tabs(out, indent + 1).append(ctx.tryDidReturnVar).append(" = true\n");
        tabs(out, indent + 1).append("return\n");
        tabs(out, indent).append("}\n");
      } else {
        tabs(out, indent)
            .append("if ")
            .append(didRetVar)
            .append(" { return")
            .append(retValVar != null ? (" " + retValVar) : "")
            .append(" }\n");
      }
    }
    if (hasOuterBreak[0]) {
      if (ctx.inTryWrapper
          && ctx.loopDepthInTryIife == 0
          && ctx.switchDepthInTryIife == 0
          && ctx.tryFlowVar != null) {
        tabs(out, indent)
            .append("if ")
            .append(flowVar)
            .append(" == 1 { ")
            .append(ctx.tryFlowVar)
            .append(" = 1; return }\n");
      } else {
        tabs(out, indent).append("if ").append(flowVar).append(" == 1 { break }\n");
      }
    }
    if (hasOuterContinue[0]) {
      if (ctx.inTryWrapper && ctx.loopDepthInTryIife == 0 && ctx.tryFlowVar != null) {
        tabs(out, indent)
            .append("if ")
            .append(flowVar)
            .append(" == 2 { ")
            .append(ctx.tryFlowVar)
            .append(" = 2; return }\n");
      } else {
        tabs(out, indent).append("if ").append(flowVar).append(" == 2 { continue }\n");
      }
    }
  }

  public static StringBuilder tabs(StringBuilder sb, int count) {
    for (int i = 0; i < count; i++) {
      sb.append('\t');
    }
    return sb;
  }

  private static TreePath unwrapParensPath(TreePath tp) {
    while (tp != null && tp.getLeaf() instanceof ParenthesizedTree pt) {
      tp = new TreePath(tp, pt.getExpression());
    }
    return tp;
  }

  private static boolean isArrowExprOrThrowSwitch(SwitchExpressionTree set) {
    if (set.getCases().isEmpty() || containsDirectTryTree(set)) {
      return false;
    }
    for (CaseTree ct : set.getCases()) {
      if (ct.getCaseKind() != CaseTree.CaseKind.RULE) {
        return false;
      }
      Tree body = ct.getBody();
      if (!(body instanceof ExpressionTree) && !(body instanceof ThrowTree)) {
        return false;
      }
    }
    return true;
  }

  private static boolean containsDirectTryTree(Tree root) {
    boolean[] found = {false};
    new TreeScanner<Void, Void>() {
      @Override
      public Void scan(Tree tree, Void p) {
        if (found[0] || tree == null) {
          return null;
        }
        return super.scan(tree, p);
      }

      @Override
      public Void visitTry(TryTree node, Void p) {
        found[0] = true;
        return null;
      }

      @Override
      public Void visitLambdaExpression(LambdaExpressionTree node, Void p) {
        return null;
      }

      @Override
      public Void visitClass(ClassTree node, Void p) {
        return null;
      }
    }.scan(root, null);
    return found[0];
  }

  private boolean canEmitDirectDeferFinally(TryTree tt, TreePath path, EmitContext ctx) {
    if (!tt.getCatches().isEmpty()
        || !tt.getResources().isEmpty()
        || tt.getFinallyBlock() == null
        || ctx.inTryWrapper
        || !ctx.switchExprTargetStack.isEmpty()) {
      return false;
    }
    if (isInsideEnclosingLoop(path)) {
      return false;
    }
    if (hasDisallowedJumpForDirectDefer(tt.getFinallyBlock(), false)) {
      return false;
    }
    if (hasDisallowedJumpForDirectDefer(tt.getBlock(), true)) {
      return false;
    }
    return stmtAlwaysReturnsOrThrows(tt.getBlock()) || isAtExecutableTailPosition(path, ctx);
  }

  private static boolean isInsideEnclosingLoop(TreePath path) {
    for (TreePath cur = path.getParentPath(); cur != null; cur = cur.getParentPath()) {
      Tree leaf = cur.getLeaf();
      if (leaf instanceof MethodTree
          || leaf instanceof LambdaExpressionTree
          || leaf instanceof ClassTree) {
        return false;
      }
      if (leaf instanceof WhileLoopTree
          || leaf instanceof DoWhileLoopTree
          || leaf instanceof ForLoopTree
          || leaf instanceof EnhancedForLoopTree) {
        return true;
      }
    }
    return false;
  }

  private static boolean isAtExecutableTailPosition(TreePath path, EmitContext ctx) {
    Tree child = path.getLeaf();
    for (TreePath cur = path.getParentPath(); cur != null; cur = cur.getParentPath()) {
      Tree parent = cur.getLeaf();
      if (parent instanceof MethodTree mt) {
        if (mt.getReturnType() == null
            && ctx.currentType != null
            && ctx.currentType.getKind() == ElementKind.RECORD) {
          return false;
        }
        return mt.getBody() == child;
      }
      if (parent instanceof LambdaExpressionTree let) {
        return let.getBody() == child;
      }
      if (parent instanceof BlockTree bt) {
        List<? extends StatementTree> stmts = bt.getStatements();
        if (stmts.isEmpty() || stmts.get(stmts.size() - 1) != child) {
          return false;
        }
      } else if (parent instanceof LabeledStatementTree lst) {
        if (lst.getStatement() != child) {
          return false;
        }
      } else if (parent instanceof IfTree it) {
        if (it.getThenStatement() != child && it.getElseStatement() != child) {
          return false;
        }
      } else if (parent instanceof TryTree tt) {
        if (tt.getBlock() != child || !tt.getCatches().isEmpty() || !tt.getResources().isEmpty()) {
          return false;
        }
      } else {
        return false;
      }
      child = parent;
    }
    return false;
  }

  private static boolean hasDisallowedJumpForDirectDefer(BlockTree block, boolean allowReturn) {
    boolean[] bad = {false};
    Set<Name> localLabels = new HashSet<>();
    // d[0] = loopDepth, d[1] = switchDepth, d[2] = switchExprDepth
    new TreeScanner<Void, int[]>() {
      @Override
      public Void scan(Tree tree, int[] d) {
        if (bad[0] || tree == null) {
          return null;
        }
        return super.scan(tree, d);
      }

      @Override
      public Void visitLambdaExpression(LambdaExpressionTree node, int[] d) {
        return null;
      }

      @Override
      public Void visitClass(ClassTree node, int[] d) {
        return null;
      }

      @Override
      public Void visitLabeledStatement(LabeledStatementTree node, int[] d) {
        boolean added = localLabels.add(node.getLabel());
        super.visitLabeledStatement(node, d);
        if (added) {
          localLabels.remove(node.getLabel());
        }
        return null;
      }

      @Override
      public Void visitWhileLoop(WhileLoopTree node, int[] d) {
        return super.visitWhileLoop(node, new int[] {d[0] + 1, d[1], d[2]});
      }

      @Override
      public Void visitDoWhileLoop(DoWhileLoopTree node, int[] d) {
        return super.visitDoWhileLoop(node, new int[] {d[0] + 1, d[1], d[2]});
      }

      @Override
      public Void visitForLoop(ForLoopTree node, int[] d) {
        return super.visitForLoop(node, new int[] {d[0] + 1, d[1], d[2]});
      }

      @Override
      public Void visitEnhancedForLoop(EnhancedForLoopTree node, int[] d) {
        return super.visitEnhancedForLoop(node, new int[] {d[0] + 1, d[1], d[2]});
      }

      @Override
      public Void visitSwitch(SwitchTree node, int[] d) {
        return super.visitSwitch(node, new int[] {d[0], d[1] + 1, d[2]});
      }

      @Override
      public Void visitSwitchExpression(SwitchExpressionTree node, int[] d) {
        return super.visitSwitchExpression(node, new int[] {d[0], d[1], d[2] + 1});
      }

      @Override
      public Void visitReturn(ReturnTree node, int[] d) {
        if (!allowReturn) {
          bad[0] = true;
          return null;
        }
        return super.visitReturn(node, d);
      }

      @Override
      public Void visitYield(YieldTree node, int[] d) {
        if (d[2] == 0) {
          bad[0] = true;
          return null;
        }
        return super.visitYield(node, d);
      }

      @Override
      public Void visitBreak(BreakTree node, int[] d) {
        if (node.getLabel() == null) {
          if (d[0] == 0 && d[1] == 0) {
            bad[0] = true;
          }
        } else if (!localLabels.contains(node.getLabel())) {
          bad[0] = true;
        }
        return null;
      }

      @Override
      public Void visitContinue(ContinueTree node, int[] d) {
        if (node.getLabel() == null) {
          if (d[0] == 0) {
            bad[0] = true;
          }
        } else if (!localLabels.contains(node.getLabel())) {
          bad[0] = true;
        }
        return null;
      }
    }.scan(block, new int[] {0, 0, 0});
    return bad[0];
  }

  private static boolean stmtAlwaysReturnsOrThrows(StatementTree s) {
    if (s == null) {
      return false;
    }
    return switch (s.getKind()) {
      case RETURN, THROW -> true;
      case BLOCK -> {
        for (StatementTree st : ((BlockTree) s).getStatements()) {
          if (stmtAlwaysReturnsOrThrows(st)) {
            yield true;
          }
        }
        yield false;
      }
      case IF -> {
        IfTree it = (IfTree) s;
        yield it.getElseStatement() != null
            && stmtAlwaysReturnsOrThrows(it.getThenStatement())
            && stmtAlwaysReturnsOrThrows(it.getElseStatement());
      }
      case TRY -> {
        TryTree tt = (TryTree) s;
        if (!stmtAlwaysReturnsOrThrows(tt.getBlock())) {
          yield false;
        }
        for (CatchTree ct : tt.getCatches()) {
          if (!stmtAlwaysReturnsOrThrows(ct.getBlock())) {
            yield false;
          }
        }
        yield true;
      }
      default -> false;
    };
  }

  private String rtPrefix(EmitContext ctx) {
    String a = resolver.runtime.j2gortAlias(ctx.imports);
    return a.isEmpty() ? "" : (a + ".");
  }

  private String jlPrefix(EmitContext ctx) {
    String a = resolver.runtime.javalangAlias(ctx.imports);
    return a.isEmpty() ? "" : (a + ".");
  }
}
