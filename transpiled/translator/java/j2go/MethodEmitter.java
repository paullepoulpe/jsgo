package j2go;

import com.sun.source.tree.BlockTree;
import com.sun.source.tree.ExpressionStatementTree;
import com.sun.source.tree.ExpressionTree;
import com.sun.source.tree.IdentifierTree;
import com.sun.source.tree.MemberSelectTree;
import com.sun.source.tree.MethodInvocationTree;
import com.sun.source.tree.MethodTree;
import com.sun.source.tree.ReturnTree;
import com.sun.source.tree.StatementTree;
import com.sun.source.tree.ThrowTree;
import com.sun.source.tree.VariableTree;
import com.sun.source.util.TreePath;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Set;
import javax.lang.model.element.Element;
import javax.lang.model.element.ElementKind;
import javax.lang.model.element.ExecutableElement;
import javax.lang.model.element.Modifier;
import javax.lang.model.element.RecordComponentElement;
import javax.lang.model.element.TypeElement;
import javax.lang.model.element.TypeParameterElement;
import javax.lang.model.element.VariableElement;
import javax.lang.model.type.ExecutableType;
import javax.lang.model.type.TypeKind;
import javax.lang.model.type.TypeMirror;

/**
 * Emits Go constructors, instance methods, static methods, interface default helpers,
 * and enum constant-specific method overrides.
 */
public final class MethodEmitter {

  private final TypeResolver resolver;
  private final ExprEmitter exprEmitter;
  private final StmtEmitter stmtEmitter;

  public MethodEmitter(
      TypeResolver resolver, ExprEmitter exprEmitter, StmtEmitter stmtEmitter) {
    this.resolver = resolver;
    this.exprEmitter = exprEmitter;
    this.stmtEmitter = stmtEmitter;
  }

  public void emitConstructor(
      MethodTree mt,
      TreePath mtPath,
      TypeElement typeEl,
      boolean hasClinit,
      Map<VariableElement, String> capturedFields,
      StringBuilder out,
      EmitContext fileCtx) {
    Element el = resolver.trees.getElement(mtPath);
    if (!(el instanceof ExecutableElement ctor)) {
      return;
    }

    // Special case: Enum constructors -> emit `initEnum_<GoName>(this_ *<Enum>, name_ string, ord_ int32, params...)`
    if (typeEl.getKind() == ElementKind.ENUM) {
      emitEnumConstructorInit(mt, mtPath, typeEl, ctor, out, fileCtx);
      return;
    }

    // Special case: Node and its 4 private inner classes
    if (resolver.isNodeClass(typeEl)) {
      emitNodeConstructor(mt, mtPath, typeEl, ctor, hasClinit, out, fileCtx);
      return;
    }

    String newFnName = resolver.getConstructorGoName(ctor);
    String initFnName = "Init" + newFnName.substring(3);
    String goTypeName = resolver.planner.getGoTypeName(typeEl);
    boolean isSubclassed = resolver.isSubclassed(typeEl);
    String structName = isSubclassed ? (goTypeName + "Base") : goTypeName;
    List<TypeParameterElement> effTypeParams = resolver.getEffectiveTypeParameters(typeEl);
    String declParams = resolver.formatGoTypeParamDecl(effTypeParams);
    String useParams = resolver.formatGoTypeParamUse(effTypeParams);
    String structWithParams = structName + useParams;
    String retGoType = resolver.goTypeForElement(typeEl, fileCtx.imports);

    EmitContext ctx = fileCtx.copyForMember(typeEl, false);
    ctx.currentExec = ctor;
    ctx.capturedVarFields.putAll(capturedFields);

    // Extra leading parameters: [outer_] + [capturedVars...]
    List<String> extraSigParams = new ArrayList<>();
    List<String> extraCallArgs = new ArrayList<>();
    StringBuilder allocAssigns = new StringBuilder();

    if (resolver.hasOuterInstance(typeEl)) {
      TypeElement outerTe = resolver.getEnclosingInstanceType(typeEl);
      String outerGoType = resolver.goTypeForElement(outerTe, ctx.imports);
      extraSigParams.add("outer_ " + outerGoType);
      extraCallArgs.add("outer_");
      allocAssigns.append("\tthis_.outer_ = outer_\n");
    }

    List<VariableElement> caps = resolver.getCapturedVars(typeEl);
    for (VariableElement cap : caps) {
      String capParam = "cap_p_" + GoNames.cleanChars(cap.getSimpleName().toString());
      String capField = capturedFields.get(cap);
      String capGoType = resolver.goType(cap.asType(), ctx.imports);
      extraSigParams.add(capParam + " " + capGoType);
      extraCallArgs.add(capParam);
      allocAssigns.append("\tthis_.").append(capField).append(" = ").append(capParam).append('\n');
    }

    List<String> javaSigParams = new ArrayList<>();
    List<String> javaCallArgs = new ArrayList<>();
    StringBuilder unusedParamGuards = new StringBuilder();

    for (int i = 0; i < mt.getParameters().size(); i++) {
      VariableTree pvt = mt.getParameters().get(i);
      VariableElement pve =
          (VariableElement) resolver.trees.getElement(new TreePath(mtPath, pvt));
      String pName = ctx.registerLocalVar(pve);
      String pGoType = resolver.goType(pve.asType(), ctx.imports);
      javaSigParams.add(pName + " " + pGoType);
      javaCallArgs.add(pName);
      unusedParamGuards.append("\t_ = ").append(pName).append('\n');
    }

    List<String> fullSigParams = new ArrayList<>(extraSigParams);
    fullSigParams.addAll(javaSigParams);

    // 1. Emit `New<GoType>[_...](...) <RetGoType>` (only for concrete non-abstract classes!)
    if (!typeEl.getModifiers().contains(Modifier.ABSTRACT)) {
      out.append("// Java constructor: ")
          .append(typeEl.getSimpleName())
          .append(formatJavaParamsComment(ctor))
          .append('\n');
      out.append("func ")
          .append(newFnName)
          .append(declParams)
          .append('(')
          .append(String.join(", ", fullSigParams))
          .append(") ")
          .append(retGoType)
          .append(" {\n");
      if (hasClinit) {
        out.append("\tclinit_").append(goTypeName).append("()\n");
      }
      if (resolver.isValueStructClass(typeEl)) {
        out.append("\tthis_ := ").append(structWithParams).append("{}\n");
        out.append(allocAssigns);
        List<String> initArgs = new ArrayList<>();
        initArgs.add("&this_");
        initArgs.addAll(javaCallArgs);
        out.append('\t')
            .append(initFnName)
            .append(useParams)
            .append('(')
            .append(String.join(", ", initArgs))
            .append(")\n");
        out.append("\treturn this_\n}\n\n");
      } else {
        out.append("\tthis_ := &").append(structWithParams).append("{}\n");
        if (resolver.needsBindSelf(typeEl)) {
          out.append("\tthis_.BindSelf_(this_)\n");
        }
        out.append(allocAssigns);
        List<String> initArgs = new ArrayList<>();
        initArgs.add("this_");
        initArgs.addAll(javaCallArgs);
        out.append('\t')
            .append(initFnName)
            .append(useParams)
            .append('(')
            .append(String.join(", ", initArgs))
            .append(")\n");
        out.append("\treturn this_\n}\n\n");
      }
    }

    // 2. Emit `Init_<GoType>[_...](this_ *<StructName>, javaParams...)`
    List<String> initSigParams = new ArrayList<>();
    initSigParams.add("this_ *" + structWithParams);
    initSigParams.addAll(javaSigParams);

    out.append("func ")
        .append(initFnName)
        .append(declParams)
        .append('(')
        .append(String.join(", ", initSigParams))
        .append(") {\n");
    out.append(unusedParamGuards);

    BlockTree body = mt.getBody();
    TreePath bodyPath = new TreePath(mtPath, body);
    List<? extends StatementTree> stmts = body != null ? body.getStatements() : List.of();

    int startIdx = 0;
    boolean delegatedToThis = false;

    if (!stmts.isEmpty() && stmts.get(0) instanceof ExpressionStatementTree est) {
      ExpressionTree expr0 = est.getExpression();
      if (expr0 instanceof MethodInvocationTree call0) {
        String kind = classifyCtorChainingCall(call0);
        if (kind != null) {
          startIdx = 1;
          TreePath call0Path = new TreePath(new TreePath(bodyPath, est), call0);
          if (kind.equals("this")) {
            delegatedToThis = true;
            emitThisCtorDelegation(call0, call0Path, out, ctx);
          } else {
            emitSuperCtorDelegation(call0, call0Path, typeEl, out, ctx);
          }
        }
      }
    }

    if (!delegatedToThis) {
      if (startIdx == 0) {
        // Implicit super() if not present in AST
        emitImplicitSuperInit(typeEl, out, ctx);
      }
      out.append("\tthis_.initFields_").append(goTypeName).append("()\n");
    }

    for (int i = startIdx; i < stmts.size(); i++) {
      StatementTree s = stmts.get(i);
      stmtEmitter.emitStmt(s, new TreePath(bodyPath, s), 1, null, out, ctx);
    }

    // Ensure Java record compact constructors assign parameters to record component fields
    if (typeEl.getKind() == ElementKind.RECORD
        && mt.getParameters().size() == typeEl.getRecordComponents().size()) {
      List<? extends RecordComponentElement> comps = typeEl.getRecordComponents();
      for (int i = 0; i < comps.size(); i++) {
        String compName = comps.get(i).getSimpleName().toString();
        VariableElement pve =
            (VariableElement)
                resolver.trees.getElement(new TreePath(mtPath, mt.getParameters().get(i)));
        String pGo = ctx.lookupLocalVar(pve);
        VariableElement fieldVe = findFieldByName(typeEl, compName);
        if (fieldVe != null && pGo != null) {
          out.append("\tthis_.")
              .append(resolver.getFieldGoName(fieldVe))
              .append(" = ")
              .append(pGo)
              .append('\n');
        }
      }
    }

    if (!delegatedToThis
        && typeEl
            .getQualifiedName()
            .contentEquals("com.google.javascript.rhino.RhinoStringPool.LazyInternedStringList")) {
      String atomicPkg = fileCtx.imports.std("sync/atomic");
      String ns = resolver.runtime.nodestoreAlias(fileCtx.imports);
      String nsPref = ns.isEmpty() ? "" : (ns + ".");
      out.append("\tif this_.size > 0 {\n\t\tthis_.ids_ = make([]uint32, int(this_.size))\n\t\t")
          .append(atomicPkg)
          .append(".StoreUint64(&this_.epoch_, ")
          .append(nsPref)
          .append("Epoch())\n\t}\n");
    }

    out.append("}\n\n");
    fileCtx.syncTempSeqFrom(ctx);
    fileCtx.deferredLocalTypes.addAll(ctx.deferredLocalTypes);
  }

  private VariableElement findFieldByName(TypeElement te, String simpleName) {
    for (Element m : te.getEnclosedElements()) {
      if (m instanceof VariableElement ve
          && ve.getKind() == ElementKind.FIELD
          && !ve.getModifiers().contains(Modifier.STATIC)
          && ve.getSimpleName().contentEquals(simpleName)) {
        return ve;
      }
    }
    return null;
  }

  private String classifyCtorChainingCall(MethodInvocationTree call) {
    ExpressionTree sel = call.getMethodSelect();
    if (sel instanceof IdentifierTree id) {
      String n = id.getName().toString();
      if (n.equals("this") || n.equals("super")) {
        return n;
      }
    } else if (sel instanceof MemberSelectTree ms) {
      if (ms.getIdentifier().contentEquals("super")) {
        return "super";
      }
    }
    return null;
  }

  private void emitThisCtorDelegation(
      MethodInvocationTree call, TreePath callPath, StringBuilder out, EmitContext ctx) {
    Element el = resolver.trees.getElement(callPath);
    if (!(el instanceof ExecutableElement targetCtor)) {
      return;
    }
    String targetNew = resolver.getConstructorGoName(targetCtor);
    String targetInit = "Init" + targetNew.substring(3);
    TypeElement tOwner =
        targetCtor.getEnclosingElement() instanceof TypeElement t ? t : null;
    String useParams =
        tOwner != null
            ? resolver.formatGoTypeParamUse(resolver.getEffectiveTypeParameters(tOwner))
            : "";
    List<String> args = new ArrayList<>();
    args.add("this_");
    appendCtorCallArgs(
        call.getArguments(),
        callPath,
        targetCtor,
        tOwner != null ? tOwner.asType() : null,
        args,
        ctx);
    out.append('\t')
        .append(targetInit)
        .append(useParams)
        .append('(')
        .append(String.join(", ", args))
        .append(")\n");
  }

  private void emitSuperCtorDelegation(
      MethodInvocationTree call,
      TreePath callPath,
      TypeElement subType,
      StringBuilder out,
      EmitContext ctx) {
    TypeElement superTe = resolver.getSuperclassElement(subType);
    if (superTe == null) {
      return;
    }

    // Shim exception superclass (`Exception`, `RuntimeException`, `IllegalArgumentException`, `IOException`, etc.)
    if (resolver.runtime.isShimException(superTe)) {
      List<String> exArgs = new ArrayList<>();
      exArgs.add("this_");
      exArgs.add(GoNames.goStringLiteral(subType.getQualifiedName().toString()));
      for (ExpressionTree a : call.getArguments()) {
        exArgs.add(
            exprEmitter.emitAs(a, callPath, resolver.getObjectElement().asType(), true, ctx));
      }
      String rt = rtPrefix(ctx);
      out.append('\t')
          .append(rt)
          .append("InitUserThrowable(")
          .append(String.join(", ", exArgs))
          .append(")\n");
      return;
    }

    Element el = resolver.trees.getElement(callPath);
    if (!(el instanceof ExecutableElement superCtor)) {
      return;
    }

    String superStruct = resolver.goStructTypeForElement(superTe, ctx.imports);
    int dot = superStruct.lastIndexOf('.');
    String embeddedField = dot >= 0 ? superStruct.substring(dot + 1) : superStruct;

    // If superTe is a non-static inner class, assign `&this_.<SuperBase>.outer_`!
    if (resolver.isSourceType(superTe) && resolver.hasOuterInstance(superTe)) {
      TypeElement superOuterTe = resolver.getEnclosingInstanceType(superTe);
      String outerVal;
      if (call.getMethodSelect() instanceof MemberSelectTree ms) {
        outerVal =
            exprEmitter.emitAs(
                ms.getExpression(),
                new TreePath(callPath, ms),
                superOuterTe.asType(),
                false,
                ctx);
      } else {
        outerVal = exprEmitter.resolveImplicitThisValue(superOuterTe, ctx);
      }
      String superPkg = resolver.planner.getGoPkg(superTe);
      if (superPkg.equals(ctx.unitPlan.goPkg)) {
        out.append("\tthis_.").append(embeddedField).append(".outer_ = ").append(outerVal).append('\n');
      } else {
        out.append("\tthis_.").append(embeddedField).append(".SetOuter_(").append(outerVal).append(")\n");
      }
    }

    String superNew = resolver.getConstructorGoName(superCtor);
    String superInit =
        "Init"
            + superNew.substring(3)
            + resolver.formatSuperTypeArgs(subType.getSuperclass(), superTe, ctx.imports);
    List<String> args = new ArrayList<>();
    args.add("&this_." + embeddedField);
    appendCtorCallArgs(
        call.getArguments(), callPath, superCtor, subType.getSuperclass(), args, ctx);

    if (resolver.isSourceType(superTe)) {
      String superPkg = resolver.planner.getGoPkg(superTe);
      String alias = ctx.imports.ext(resolver.planner.getGoImportPath(superPkg), superPkg);
      String prefix = alias.isEmpty() ? "" : (alias + ".");
      out.append('\t').append(prefix).append(superInit).append('(').append(String.join(", ", args)).append(")\n");
    } else {
      out.append('\t').append(rtPrefix(ctx)).append(superInit).append('(').append(String.join(", ", args)).append(")\n");
    }
  }

  private void emitImplicitSuperInit(TypeElement subType, StringBuilder out, EmitContext ctx) {
    TypeElement superTe = resolver.getSuperclassElement(subType);
    if (superTe != null && resolver.runtime.isShimException(superTe)) {
      out.append('\t')
          .append(rtPrefix(ctx))
          .append("InitUserThrowable(this_, ")
          .append(GoNames.goStringLiteral(subType.getQualifiedName().toString()))
          .append(")\n");
    }
  }

  private void appendCtorCallArgs(
      List<? extends ExpressionTree> argTrees,
      TreePath callPath,
      ExecutableElement ctor,
      List<String> outArgs,
      EmitContext ctx) {
    appendCtorCallArgs(argTrees, callPath, ctor, null, outArgs, ctx);
  }

  private void appendCtorCallArgs(
      List<? extends ExpressionTree> argTrees,
      TreePath callPath,
      ExecutableElement ctor,
      TypeMirror ownerSiteTm,
      List<String> outArgs,
      EmitContext ctx) {
    int n = ctor.getParameters().size();
    boolean varArgs = ctor.isVarArgs() && n > 0;
    int fixed = varArgs ? (n - 1) : n;
    TypeElement ctorOwner =
        ctor.getEnclosingElement() instanceof TypeElement t ? t : null;
    for (int i = 0; i < fixed && i < argTrees.size(); i++) {
      TypeMirror paramTm = ctor.getParameters().get(i).asType();
      if (ownerSiteTm != null && ctorOwner != null) {
        paramTm = resolver.substituteOwnerTypeArgs(paramTm, ctorOwner, ownerSiteTm);
      }
      outArgs.add(
          exprEmitter.emitAs(
              argTrees.get(i),
              callPath,
              resolver.goEffectiveType(paramTm),
              false,
              ctx));
    }
    if (varArgs) {
      TypeMirror lastParamTm = ctor.getParameters().get(n - 1).asType();
      if (ownerSiteTm != null && ctorOwner != null) {
        lastParamTm = resolver.substituteOwnerTypeArgs(lastParamTm, ctorOwner, ownerSiteTm);
      }
      javax.lang.model.type.ArrayType at =
          (javax.lang.model.type.ArrayType) resolver.goEffectiveType(lastParamTm);
      if (argTrees.size() == n
          && exprEmitter.typeOf(new TreePath(callPath, argTrees.get(n - 1))).getKind()
              == TypeKind.ARRAY) {
        outArgs.add(exprEmitter.emitAs(argTrees.get(n - 1), callPath, at, false, ctx));
      } else {
        List<String> packed = new ArrayList<>();
        for (int i = fixed; i < argTrees.size(); i++) {
          packed.add(
              exprEmitter.emitAs(argTrees.get(i), callPath, at.getComponentType(), false, ctx));
        }
        outArgs.add(resolver.goType(at, ctx.imports) + "{" + String.join(", ", packed) + "}");
      }
    }
  }

  private void emitEnumConstructorInit(
      MethodTree mt,
      TreePath mtPath,
      TypeElement enumTe,
      ExecutableElement ctor,
      StringBuilder out,
      EmitContext fileCtx) {
    String newFn = resolver.getConstructorGoName(ctor);
    String initFn = "initEnum" + newFn.substring(3);
    String enumGoName = resolver.planner.getGoTypeName(enumTe);

    EmitContext ctx = fileCtx.copyForMember(enumTe, false);
    ctx.currentExec = ctor;

    List<String> sigParams = new ArrayList<>();
    sigParams.add("this_ *" + enumGoName);
    sigParams.add("enumName_ string");
    sigParams.add("enumOrd_ int32");
    StringBuilder unusedGuards = new StringBuilder();

    for (VariableTree pvt : mt.getParameters()) {
      VariableElement pve =
          (VariableElement) resolver.trees.getElement(new TreePath(mtPath, pvt));
      String pName = ctx.registerLocalVar(pve);
      sigParams.add(pName + " " + resolver.goType(pve.asType(), ctx.imports));
      unusedGuards.append("\t_ = ").append(pName).append('\n');
    }

    out.append("func ")
        .append(initFn)
        .append('(')
        .append(String.join(", ", sigParams))
        .append(") *" + enumGoName + " {\n");
    out.append("\tthis_.EnumBase.Name_ = enumName_\n");
    out.append("\tthis_.EnumBase.Ordinal_ = enumOrd_\n");
    out.append(unusedGuards);

    BlockTree body = mt.getBody();
    TreePath bodyPath = new TreePath(mtPath, body);
    List<? extends StatementTree> stmts = body != null ? body.getStatements() : List.of();
    int startIdx = 0;
    boolean delegatedToThis = false;

    if (!stmts.isEmpty()
        && stmts.get(0) instanceof ExpressionStatementTree est
        && est.getExpression() instanceof MethodInvocationTree call0) {
      String kind = classifyCtorChainingCall(call0);
      if (kind != null) {
        startIdx = 1;
        if (kind.equals("this")) {
          delegatedToThis = true;
          TreePath call0Path = new TreePath(new TreePath(bodyPath, est), call0);
          Element tEl = resolver.trees.getElement(call0Path);
          if (tEl instanceof ExecutableElement targetCtor) {
            String tInit = "initEnum" + resolver.getConstructorGoName(targetCtor).substring(3);
            List<String> args = new ArrayList<>();
            args.add("this_");
            args.add("enumName_");
            args.add("enumOrd_");
            appendCtorCallArgs(call0.getArguments(), call0Path, targetCtor, args, ctx);
            out.append('\t')
                .append(tInit)
                .append('(')
                .append(String.join(", ", args))
                .append(")\n");
          }
        }
      }
    }

    if (!delegatedToThis) {
      out.append("\tthis_.initFields_").append(enumGoName).append("()\n");
    }

    for (int i = startIdx; i < stmts.size(); i++) {
      StatementTree s = stmts.get(i);
      stmtEmitter.emitStmt(s, new TreePath(bodyPath, s), 1, null, out, ctx);
    }

    out.append("\treturn this_\n}\n\n");
    fileCtx.syncTempSeqFrom(ctx);
    fileCtx.deferredLocalTypes.addAll(ctx.deferredLocalTypes);
  }

  private void emitNodeConstructor(
      MethodTree mt,
      TreePath mtPath,
      TypeElement nodeTe,
      ExecutableElement ctor,
      boolean hasClinit,
      StringBuilder out,
      EmitContext fileCtx) {
    String newFn = resolver.getConstructorGoName(ctor);
    String initFn = "initNode" + newFn.substring(3);
    int kindFlag = resolver.getNodeKindFlag(nodeTe);

    EmitContext ctx = fileCtx.copyForMember(nodeTe, false);
    ctx.currentExec = ctor;

    List<String> sigParams = new ArrayList<>();
    List<String> callArgs = new ArrayList<>();
    StringBuilder unusedGuards = new StringBuilder();

    for (VariableTree pvt : mt.getParameters()) {
      VariableElement pve =
          (VariableElement) resolver.trees.getElement(new TreePath(mtPath, pvt));
      String pName = ctx.registerLocalVar(pve);
      sigParams.add(pName + " " + resolver.goType(pve.asType(), ctx.imports));
      callArgs.add(pName);
      unusedGuards.append("\t_ = ").append(pName).append('\n');
    }

    out.append("func ")
        .append(newFn)
        .append('(')
        .append(String.join(", ", sigParams))
        .append(") Node {\n");
    if (hasClinit && kindFlag == 0) {
      out.append("\tclinit_Node()\n");
    }
    String ns = resolver.runtime.nodestoreAlias(ctx.imports);
    String nsPref = ns.isEmpty() ? "" : (ns + ".");
    out.append("\traw_, rec_ := ").append(nsPref).append("AllocNodeFastRec()\n");
    out.append("\trec_.LinenoCharno = -1\n");
    out.append("\trec_.Token = 0xFF\n");
    if (kindFlag != 0) {
      out.append("\trec_.KindFlags = ").append(kindFlag).append("\n");
    }
    out.append("\tthis_ := Node(raw_)\n");
    List<String> initCall = new ArrayList<>();
    initCall.add("this_");
    initCall.addAll(callArgs);
    out.append('\t').append(initFn).append('(').append(String.join(", ", initCall)).append(")\n");
    out.append("\treturn this_\n}\n\n");

    List<String> initSig = new ArrayList<>();
    initSig.add("this_ Node");
    initSig.addAll(sigParams);
    out.append("func ").append(initFn).append('(').append(String.join(", ", initSig)).append(") {\n");
    out.append(unusedGuards);

    BlockTree body = mt.getBody();
    TreePath bodyPath = new TreePath(mtPath, body);
    List<? extends StatementTree> stmts = body != null ? body.getStatements() : List.of();
    int startIdx = 0;

    if (!stmts.isEmpty()
        && stmts.get(0) instanceof ExpressionStatementTree est
        && est.getExpression() instanceof MethodInvocationTree call0) {
      String kind = classifyCtorChainingCall(call0);
      if (kind != null) {
        startIdx = 1;
        TreePath call0Path = new TreePath(new TreePath(bodyPath, est), call0);
        Element targetEl = resolver.trees.getElement(call0Path);
        if (targetEl instanceof ExecutableElement targetCtor
            && targetCtor.getEnclosingElement() instanceof TypeElement targetOwner
            && resolver.isNodeClass(targetOwner)) {
          String targetInit = "initNode" + resolver.getConstructorGoName(targetCtor).substring(3);
          List<String> args = new ArrayList<>();
          args.add("this_");
          appendCtorCallArgs(call0.getArguments(), call0Path, targetCtor, args, ctx);
          out.append('\t')
              .append(targetInit)
              .append('(')
              .append(String.join(", ", args))
              .append(")\n");
        }
      }
    }

    for (int i = startIdx; i < stmts.size(); i++) {
      StatementTree s = stmts.get(i);
      stmtEmitter.emitStmt(s, new TreePath(bodyPath, s), 1, null, out, ctx);
    }
    out.append("}\n\n");
    fileCtx.syncTempSeqFrom(ctx);
    fileCtx.deferredLocalTypes.addAll(ctx.deferredLocalTypes);
  }

  // --- Methods ---

  public void emitMethod(
      MethodTree mt,
      TreePath mtPath,
      TypeElement typeEl,
      boolean hasClinit,
      Set<String> enumOverriddenMethodGoNames,
      Map<VariableElement, String> capturedFields,
      StringBuilder out,
      EmitContext fileCtx) {
    Element el = resolver.trees.getElement(mtPath);
    if (!(el instanceof ExecutableElement ee)) {
      return;
    }

    boolean isStatic = ee.getModifiers().contains(Modifier.STATIC);
    boolean isDefault = ee.getModifiers().contains(Modifier.DEFAULT);
    boolean isAbstract = ee.getModifiers().contains(Modifier.ABSTRACT);
    boolean isNative = ee.getModifiers().contains(Modifier.NATIVE);

    // Abstract methods on interfaces have no Go method body; abstract methods on classes/enums
    // emit panic stubs on `*<Class>Base` (or dispatch via enum `fn_*`) so `*<Class>Base` satisfies `<Class>`.
    if (isAbstract && typeEl.getKind() == ElementKind.INTERFACE) {
      return;
    }

    EmitContext ctx = fileCtx.copyForMember(typeEl, isStatic);
    ctx.currentExec = ee;
    ctx.capturedVarFields.putAll(capturedFields);

    if (isStatic) {
      emitStaticMethod(mt, mtPath, typeEl, ee, hasClinit, out, ctx);
      fileCtx.syncTempSeqFrom(ctx);
      fileCtx.deferredLocalTypes.addAll(ctx.deferredLocalTypes);
      return;
    }

    if ((isDefault || ee.getModifiers().contains(Modifier.PRIVATE))
        && typeEl.getKind() == ElementKind.INTERFACE) {
      emitInterfaceDefaultHelper(mt, mtPath, typeEl, ee, out, ctx);
      fileCtx.syncTempSeqFrom(ctx);
      fileCtx.deferredLocalTypes.addAll(ctx.deferredLocalTypes);
      return;
    }

    TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
    String goTypeName = resolver.planner.getGoTypeName(typeEl);

    // Special case: Node's 4 private inner classes (`NumberNode`, `BigIntNode`, `StringNode`, `TemplateLiteralSubstringNode`)
    if (resolver.isNodeClass(typeEl) && resolver.getNodeKindFlag(typeEl) > 0) {
      int kf = resolver.getNodeKindFlag(typeEl);
      String helperName = "nodeSub" + kf + "_" + sig.goName;
      emitInstanceMethodCore(
          mt, mtPath, typeEl, ee, sig, "(this_ Node)", helperName, false, false, false, out, ctx);
      fileCtx.syncTempSeqFrom(ctx);
      fileCtx.deferredLocalTypes.addAll(ctx.deferredLocalTypes);
      return;
    }

    // Special case: Node's own methods overridden by any of its 4 inner classes
    if (typeEl.getQualifiedName().contentEquals(TypeResolver.NODE_QUAL_NAME)) {
      List<Integer> overridingKfs = findNodeOverridingKindFlags(typeEl, sig);
      if (!overridingKfs.isEmpty()) {
        emitNodePolymorphicMethod(mt, mtPath, typeEl, ee, sig, overridingKfs, out, ctx);
        fileCtx.syncTempSeqFrom(ctx);
        fileCtx.deferredLocalTypes.addAll(ctx.deferredLocalTypes);
        return;
      }
    }

    List<TypeParameterElement> effTypeParams = resolver.getEffectiveTypeParameters(typeEl);
    String useParams = resolver.formatGoTypeParamUse(effTypeParams);
    String recvSpec;
    if (resolver.isNodeClass(typeEl)) {
      recvSpec = "(this_ Node)";
    } else if (resolver.isValueStructClass(typeEl) || resolver.isValueEnum(typeEl)) {
      recvSpec = "(this_ " + goTypeName + useParams + ")";
    } else if (resolver.isSubclassed(typeEl)) {
      recvSpec = "(this_ *" + goTypeName + "Base" + useParams + ")";
    } else {
      recvSpec = "(this_ *" + goTypeName + useParams + ")";
    }

    boolean hasEnumOverrideHook =
        typeEl.getKind() == ElementKind.ENUM
            && enumOverriddenMethodGoNames != null
            && enumOverriddenMethodGoNames.contains(sig.goName);

    emitInstanceMethodCore(
        mt,
        mtPath,
        typeEl,
        ee,
        sig,
        recvSpec,
        sig.goName,
        isAbstract,
        isNative,
        hasEnumOverrideHook,
        out,
        ctx);
    if (typeEl.getQualifiedName().contentEquals("com.google.javascript.jscomp.SimpleMultimap")
        && !resolver.isSubclassed(typeEl)
        && ee.getSimpleName().contentEquals("clear")
        && ee.getParameters().isEmpty()) {
      emitSimpleMultimapFastMethods(goTypeName + useParams, out, ctx);
    }
    fileCtx.syncTempSeqFrom(ctx);
    fileCtx.deferredLocalTypes.addAll(ctx.deferredLocalTypes);
  }

  private void emitSimpleMultimapFastMethods(
      String goTypeName, StringBuilder out, EmitContext ctx) {
    String rt = rtPrefix(ctx);
    String colContains = resolver.runtime.findExtMethodGoName("java.util.Collection", "contains", 1);
    String colAdd = resolver.runtime.findExtMethodGoName("java.util.Collection", "add", 1);
    String colRemove = resolver.runtime.findExtMethodGoName("java.util.Collection", "remove", 1);
    String colIsEmpty = resolver.runtime.findExtMethodGoName("java.util.Collection", "isEmpty", 0);
    String colSize = resolver.runtime.findExtMethodGoName("java.util.Collection", "size", 0);
    out.append("func (this_ *")
        .append(goTypeName)
        .append(") ContainsKeyStr_(key string) bool {\n\treturn ")
        .append(rt)
        .append("MapContainsKeyStr(this_.map_, key)\n}\n\n")
        .append("func (this_ *")
        .append(goTypeName)
        .append(") ContainsKeyNode_(key uint32) bool {\n\treturn ")
        .append(rt)
        .append("MapContainsKeyNode(this_.map_, key)\n}\n\n")
        .append("func (this_ *")
        .append(goTypeName)
        .append(") ContainsEntryStr_(key string, value any) bool {\n\tlist := ")
        .append(rt)
        .append("AsIface[")
        .append(rt)
        .append("Ext_java_util_List](")
        .append(rt)
        .append("MapGetStr(this_.map_, key))\n\treturn list != nil && list.")
        .append(colContains)
        .append("(")
        .append(rt)
        .append("ToAny(value))\n}\n\n")
        .append("func (this_ *")
        .append(goTypeName)
        .append(") ContainsEntryNode_(key uint32, value any) bool {\n\tlist := ")
        .append(rt)
        .append("AsIface[")
        .append(rt)
        .append("Ext_java_util_List](")
        .append(rt)
        .append("MapGetNode(this_.map_, key))\n\treturn list != nil && list.")
        .append(colContains)
        .append("(")
        .append(rt)
        .append("ToAny(value))\n}\n\n")
        .append("func (this_ *")
        .append(goTypeName)
        .append(") PutStr_(key string, value any) bool {\n\tvalue = ")
        .append(rt)
        .append("ToAny(value)\n\tlist := ")
        .append(rt)
        .append("AsIface[")
        .append(rt)
        .append("Ext_java_util_List](")
        .append(rt)
        .append("MapGetStr(this_.map_, key))\n\tif list == nil {\n\t\tal := ")
        .append(rt)
        .append("NewExt_java_util_ArrayList()\n\t\t")
        .append(rt)
        .append("MapPutStr(this_.map_, key, al)\n\t\tlist = al\n\t}\n\tif list.")
        .append(colAdd)
        .append("(value) {\n\t\tthis_.totalSize++\n\t\treturn true\n\t}\n\treturn false\n}\n\n")
        .append("func (this_ *")
        .append(goTypeName)
        .append(") PutNode_(key uint32, value any) bool {\n\tvalue = ")
        .append(rt)
        .append("ToAny(value)\n\tlist := ")
        .append(rt)
        .append("AsIface[")
        .append(rt)
        .append("Ext_java_util_List](")
        .append(rt)
        .append("MapGetNode(this_.map_, key))\n\tif list == nil {\n\t\tal := ")
        .append(rt)
        .append("NewExt_java_util_ArrayList()\n\t\t")
        .append(rt)
        .append("MapPutNode(this_.map_, key, al)\n\t\tlist = al\n\t}\n\tif list.")
        .append(colAdd)
        .append("(value) {\n\t\tthis_.totalSize++\n\t\treturn true\n\t}\n\treturn false\n}\n\n")
        .append("func (this_ *")
        .append(goTypeName)
        .append(") GetStr_(key string) ")
        .append(rt)
        .append("Ext_java_util_List {\n\tlist := ")
        .append(rt)
        .append("AsIface[")
        .append(rt)
        .append("Ext_java_util_List](")
        .append(rt)
        .append("MapGetStr(this_.map_, key))\n\tif list == nil {\n\t\treturn ")
        .append(rt)
        .append("Ext_java_util_Collections_EmptyList()\n\t}\n\treturn list\n}\n\n")
        .append("func (this_ *")
        .append(goTypeName)
        .append(") GetNode_(key uint32) ")
        .append(rt)
        .append("Ext_java_util_List {\n\tlist := ")
        .append(rt)
        .append("AsIface[")
        .append(rt)
        .append("Ext_java_util_List](")
        .append(rt)
        .append("MapGetNode(this_.map_, key))\n\tif list == nil {\n\t\treturn ")
        .append(rt)
        .append("Ext_java_util_Collections_EmptyList()\n\t}\n\treturn list\n}\n\n")
        .append("func (this_ *")
        .append(goTypeName)
        .append(") RemoveStr_(key string, value any) bool {\n\tlist := ")
        .append(rt)
        .append("AsIface[")
        .append(rt)
        .append("Ext_java_util_List](")
        .append(rt)
        .append("MapGetStr(this_.map_, key))\n\tif list == nil {\n\t\treturn false\n\t}\n\tif list.")
        .append(colRemove)
        .append("(")
        .append(rt)
        .append("ToAny(value)) {\n\t\tthis_.totalSize--\n\t\tif list.")
        .append(colIsEmpty)
        .append("() {\n\t\t\t")
        .append(rt)
        .append("MapRemoveStr(this_.map_, key)\n\t\t}\n\t\treturn true\n\t}\n\treturn false\n}\n\n")
        .append("func (this_ *")
        .append(goTypeName)
        .append(") RemoveNode_(key uint32, value any) bool {\n\tlist := ")
        .append(rt)
        .append("AsIface[")
        .append(rt)
        .append("Ext_java_util_List](")
        .append(rt)
        .append("MapGetNode(this_.map_, key))\n\tif list == nil {\n\t\treturn false\n\t}\n\tif list.")
        .append(colRemove)
        .append("(")
        .append(rt)
        .append("ToAny(value)) {\n\t\tthis_.totalSize--\n\t\tif list.")
        .append(colIsEmpty)
        .append("() {\n\t\t\t")
        .append(rt)
        .append("MapRemoveNode(this_.map_, key)\n\t\t}\n\t\treturn true\n\t}\n\treturn false\n}\n\n")
        .append("func (this_ *")
        .append(goTypeName)
        .append(") RemoveAllStr_(key string) ")
        .append(rt)
        .append("Ext_java_util_List {\n\tlist := ")
        .append(rt)
        .append("AsIface[")
        .append(rt)
        .append("Ext_java_util_List](")
        .append(rt)
        .append("MapRemoveStr(this_.map_, key))\n\tif list == nil {\n\t\treturn ")
        .append(rt)
        .append("Ext_java_util_Collections_EmptyList()\n\t}\n\tthis_.totalSize -= list.")
        .append(colSize)
        .append("()\n\treturn list\n}\n\n")
        .append("func (this_ *")
        .append(goTypeName)
        .append(") RemoveAllNode_(key uint32) ")
        .append(rt)
        .append("Ext_java_util_List {\n\tlist := ")
        .append(rt)
        .append("AsIface[")
        .append(rt)
        .append("Ext_java_util_List](")
        .append(rt)
        .append("MapRemoveNode(this_.map_, key))\n\tif list == nil {\n\t\treturn ")
        .append(rt)
        .append("Ext_java_util_Collections_EmptyList()\n\t}\n\tthis_.totalSize -= list.")
        .append(colSize)
        .append("()\n\treturn list\n}\n\n")
        .append("func (this_ *")
        .append(goTypeName)
        .append(") PutAllStr_(key string, values ")
        .append(rt)
        .append("Ext_java_lang_Iterable) bool {\n\t")
        .append(rt)
        .append("CheckNotNull(values)\n\tmodified := false\n\tcur := ")
        .append(rt)
        .append("NewIterCursor(values)\n\tfor cur.HasNext() {\n\t\tif this_.PutStr_(key, cur.Next()) {\n\t\t\tmodified = true\n\t\t}\n\t}\n\treturn modified\n}\n\n")
        .append("func (this_ *")
        .append(goTypeName)
        .append(") PutAllNode_(key uint32, values ")
        .append(rt)
        .append("Ext_java_lang_Iterable) bool {\n\t")
        .append(rt)
        .append("CheckNotNull(values)\n\tmodified := false\n\tcur := ")
        .append(rt)
        .append("NewIterCursor(values)\n\tfor cur.HasNext() {\n\t\tif this_.PutNode_(key, cur.Next()) {\n\t\t\tmodified = true\n\t\t}\n\t}\n\treturn modified\n}\n\n");
  }

  /** Emits an enum constant's overridden method body as a function literal assigned in `clinit`. */
  public String emitEnumConstantMethodOverrideFunc(
      MethodTree mt, TreePath mtPath, TypeElement enumTe, EmitContext fileCtx) {
    ExecutableElement ee = (ExecutableElement) resolver.trees.getElement(mtPath);
    TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
    String enumGoName = resolver.planner.getGoTypeName(enumTe);

    EmitContext ctx = fileCtx.copyForMember(enumTe, false);
    ctx.currentExec = ee;

    StringBuilder sb = new StringBuilder();
    List<String> sigParams = new ArrayList<>();
    sigParams.add("this_ *" + enumGoName);
    StringBuilder preamble = new StringBuilder();
    TypeMirror enumSiteTm = enumTe.asType();
    buildVirtualParamHeaderAndPreamble(mt, mtPath, sig, enumSiteTm, sigParams, preamble, ctx);

    sb.append("func(").append(String.join(", ", sigParams)).append(')');
    TypeMirror canonRetTm = resolver.resolveSigReturnType(sig, enumSiteTm);
    String retGo = canonRetTm != null ? resolver.goType(canonRetTm, ctx.imports) : "";
    if (!retGo.isEmpty()) {
      sb.append(' ').append(retGo);
    }
    sb.append(" {\n");
    sb.append(preamble);
    ctx.expectedReturnType = canonRetTm;
    ctx.returnWidenedToAny = sig.returnWidenedToAny;

    if (mt.getBody() != null) {
      sb.append(
          stmtEmitter.emitBlockStatements(
              mt.getBody(), new TreePath(mtPath, mt.getBody()), 1, ctx));
      if (canonRetTm != null && !endsWithGoTerminatingReturnOrThrow(mt.getBody())) {
        sb.append("\tpanic(\"unreachable\")\n");
      }
    } else if (canonRetTm != null) {
      sb.append("\treturn ").append(resolver.zeroValue(canonRetTm, ctx.imports)).append('\n');
    }
    sb.append('}');

    fileCtx.syncTempSeqFrom(ctx);
    fileCtx.deferredLocalTypes.addAll(ctx.deferredLocalTypes);
    return sb.toString();
  }

  private void emitStaticMethod(
      MethodTree mt,
      TreePath mtPath,
      TypeElement typeEl,
      ExecutableElement ee,
      boolean hasClinit,
      StringBuilder out,
      EmitContext ctx) {
    boolean isGenStatic = resolver.isGenericStaticMethod(ee);
    List<? extends TypeParameterElement> methodTpes =
        isGenStatic ? ee.getTypeParameters() : List.of();
    resolver.pushTypeParamScope(null, methodTpes);
    try {
      String fnName = resolver.getStaticMethodGoName(ee);
      String declParams = resolver.formatGoTypeParamDecl(methodTpes);
      String goTypeName = resolver.planner.getGoTypeName(typeEl);

      List<String> sigParams = new ArrayList<>();
      StringBuilder guards = new StringBuilder();
      for (VariableTree pvt : mt.getParameters()) {
        VariableElement pve =
            (VariableElement) resolver.trees.getElement(new TreePath(mtPath, pvt));
        String pName = ctx.registerLocalVar(pve);
        sigParams.add(pName + " " + resolver.goType(pve.asType(), ctx.imports));
        guards.append("\t_ = ").append(pName).append('\n');
      }

      TypeMirror retTm =
          ee.getReturnType().getKind() == TypeKind.VOID
              ? null
              : resolver.goEffectiveType(ee.getReturnType());
      String retGo = retTm != null ? resolver.goType(retTm, ctx.imports) : "";

      out.append("// Java static method: ")
          .append(typeEl.getSimpleName())
          .append('.')
          .append(ee.getSimpleName())
          .append(formatJavaParamsComment(ee))
          .append('\n');
      out.append("func ")
          .append(fnName)
          .append(declParams)
          .append('(')
          .append(String.join(", ", sigParams))
          .append(')');
      if (!retGo.isEmpty()) {
        out.append(' ').append(retGo);
      }
      out.append(" {\n");
      if (hasClinit) {
        out.append("\tclinit_").append(goTypeName).append("()\n");
      }
      if (ee.getModifiers().contains(Modifier.SYNCHRONIZED)) {
        if (resolver.hasFastStaticSyncMutex(typeEl)) {
          out.append("\tstaticSyncMu_").append(goTypeName).append(".Lock()\n");
          out.append("\tdefer staticSyncMu_").append(goTypeName).append(".Unlock()\n");
        } else {
          String lockCls =
              rtPrefix(ctx)
                  + "ClassLiteral("
                  + GoNames.goStringLiteral(typeEl.getQualifiedName().toString())
                  + ", "
                  + GoNames.goStringLiteral(typeEl.getSimpleName().toString())
                  + ")";
          out.append("\tunlock_ := ")
              .append(rtPrefix(ctx))
              .append("LockMonitor(")
              .append(lockCls)
              .append(")\n");
          out.append("\tdefer unlock_()\n");
        }
      }
      out.append(guards);

      ctx.expectedReturnType = retTm;
      ctx.returnWidenedToAny = false;

      if (mt.getBody() != null) {
        out.append(
            stmtEmitter.emitBlockStatements(
                mt.getBody(), new TreePath(mtPath, mt.getBody()), 1, ctx));
        if (retTm != null && !endsWithGoTerminatingReturnOrThrow(mt.getBody())) {
          out.append("\tpanic(\"unreachable\")\n");
        }
      } else if (retTm != null) {
        out.append("\treturn ").append(resolver.zeroValue(retTm, ctx.imports)).append('\n');
      }
      out.append("}\n\n");
    } finally {
      resolver.popTypeParamScope();
    }
  }

  private void emitInterfaceDefaultHelper(
      MethodTree mt,
      TreePath mtPath,
      TypeElement ifaceTe,
      ExecutableElement ee,
      StringBuilder out,
      EmitContext ctx) {
    TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
    String ifaceGoName = resolver.planner.getGoTypeName(ifaceTe);
    String helperName = "IfaceDefault_" + ifaceGoName + "_" + sig.goName;
    List<TypeParameterElement> effTypeParams = resolver.getEffectiveTypeParameters(ifaceTe);
    String declParams = resolver.formatGoTypeParamDecl(effTypeParams);
    String useParams = resolver.formatGoTypeParamUse(effTypeParams);
    TypeMirror ifaceSiteTm = ifaceTe.asType();

    List<String> sigParams = new ArrayList<>();
    sigParams.add("this_ " + ifaceGoName + useParams);
    StringBuilder preamble = new StringBuilder();
    buildVirtualParamHeaderAndPreamble(mt, mtPath, sig, ifaceSiteTm, sigParams, preamble, ctx);

    TypeMirror canonRetTm = resolver.resolveSigReturnType(sig, ifaceSiteTm);
    String retGo = canonRetTm != null ? resolver.goType(canonRetTm, ctx.imports) : "";

    out.append("// Java default interface method: ")
        .append(ifaceTe.getSimpleName())
        .append('.')
        .append(ee.getSimpleName())
        .append(formatJavaParamsComment(ee))
        .append('\n');
    out.append("func ")
        .append(helperName)
        .append(declParams)
        .append('(')
        .append(String.join(", ", sigParams))
        .append(')');
    if (!retGo.isEmpty()) {
      out.append(' ').append(retGo);
    }
    out.append(" {\n");
    out.append("\t_ = this_\n");
    out.append(preamble);

    ctx.expectedReturnType = canonRetTm;
    ctx.returnWidenedToAny = sig.returnWidenedToAny;

    if (mt.getBody() != null) {
      out.append(
          stmtEmitter.emitBlockStatements(
              mt.getBody(), new TreePath(mtPath, mt.getBody()), 1, ctx));
      if (canonRetTm != null && !endsWithGoTerminatingReturnOrThrow(mt.getBody())) {
        out.append("\tpanic(\"unreachable\")\n");
      }
    }
    out.append("}\n\n");
  }

  private List<Integer> findNodeOverridingKindFlags(
      TypeElement nodeTe, TypeResolver.MethodSig sig) {
    List<Integer> kfs = new ArrayList<>();
    for (Element enc : nodeTe.getEnclosedElements()) {
      if (enc instanceof TypeElement subTe && resolver.getNodeKindFlag(subTe) > 0) {
        int kf = resolver.getNodeKindFlag(subTe);
        for (Element m : subTe.getEnclosedElements()) {
          if (m instanceof ExecutableElement subEe
              && subEe.getKind() == ElementKind.METHOD
              && !subEe.getModifiers().contains(Modifier.STATIC)) {
            if (resolver.getInstanceMethodSig(subEe) == sig) {
              kfs.add(kf);
              break;
            }
          }
        }
      }
    }
    kfs.sort(Integer::compareTo);
    return kfs;
  }

  private void emitNodePolymorphicMethod(
      MethodTree mt,
      TreePath mtPath,
      TypeElement nodeTe,
      ExecutableElement ee,
      TypeResolver.MethodSig sig,
      List<Integer> overridingKfs,
      StringBuilder out,
      EmitContext ctx) {
    // 1. Emit `base_<goName>` containing Node's own base implementation
    emitInstanceMethodCore(
        mt,
        mtPath,
        nodeTe,
        ee,
        sig,
        "(this_ Node)",
        "base_" + sig.goName,
        false,
        false,
        false,
        out,
        ctx);

    // 2. Emit `<goName>` dispatching on `nodestore.Node(this_).KindFlags()` -> `nodeSub<k>_<goName>`
    List<String> dispParams = new ArrayList<>();
    List<String> dispArgs = new ArrayList<>();
    TypeMirror nodeSiteTm = nodeTe.asType();
    for (int i = 0; i < sig.canonParamTypes.size(); i++) {
      String pn = "p" + i + "_";
      dispParams.add(
          pn + " " + resolver.goType(resolver.resolveSigParamType(sig, i, nodeSiteTm), ctx.imports));
      dispArgs.add(pn);
    }
    TypeMirror dispRetTm = resolver.resolveSigReturnType(sig, nodeSiteTm);
    String retGo = dispRetTm != null ? resolver.goType(dispRetTm, ctx.imports) : "";
    String ns = resolver.runtime.nodestoreAlias(ctx.imports);
    String nsPref = ns.isEmpty() ? "" : (ns + ".");
    String argStr = String.join(", ", dispArgs);
    String retPrefix = retGo.isEmpty() ? "" : "return ";

    out.append("func (this_ Node) ")
        .append(sig.goName)
        .append('(')
        .append(String.join(", ", dispParams))
        .append(')');
    if (!retGo.isEmpty()) {
      out.append(' ').append(retGo);
    }
    out.append(" {\n");
    out.append("\tswitch ").append(nsPref).append("Node(this_).KindFlags() {\n");
    for (int kf : overridingKfs) {
      out.append("\tcase ")
          .append(kf)
          .append(":\n\t\t")
          .append(retPrefix)
          .append("this_.nodeSub")
          .append(kf)
          .append('_')
          .append(sig.goName)
          .append('(')
          .append(argStr)
          .append(")\n");
    }
    out.append("\tdefault:\n\t\t")
        .append(retPrefix)
        .append("this_.base_")
        .append(sig.goName)
        .append('(')
        .append(argStr)
        .append(")\n");
    out.append("\t}\n}\n\n");
  }

  private void emitInstanceMethodCore(
      MethodTree mt,
      TreePath mtPath,
      TypeElement typeEl,
      ExecutableElement ee,
      TypeResolver.MethodSig sig,
      String recvSpec,
      String methodGoName,
      boolean isAbstractEnumMethod,
      boolean isNative,
      boolean hasEnumOverrideHook,
      StringBuilder out,
      EmitContext ctx) {
    List<String> sigParams = new ArrayList<>();
    List<String> rawArgPassThrough = new ArrayList<>();
    StringBuilder preamble = new StringBuilder();
    TypeMirror ownerSiteTm = typeEl.asType();

    for (int i = 0; i < mt.getParameters().size(); i++) {
      VariableTree pvt = mt.getParameters().get(i);
      VariableElement pve =
          (VariableElement) resolver.trees.getElement(new TreePath(mtPath, pvt));
      TypeMirror canonTm = resolver.resolveSigParamType(sig, i, ownerSiteTm);
      String canonGo = resolver.goType(canonTm, ctx.imports);
      TypeMirror actualTm = resolver.goEffectiveType(pve.asType());
      String actualGo = resolver.goType(actualTm, ctx.imports);

      String userVar = ctx.registerLocalVar(pve);
      if (canonGo.equals(actualGo) && !sig.paramWidenedToAny[i]) {
        sigParams.add(userVar + " " + canonGo);
        rawArgPassThrough.add(userVar);
        preamble.append("\t_ = ").append(userVar).append('\n');
      } else {
        String rawArg = ctx.nextTemp("arg");
        sigParams.add(rawArg + " " + canonGo);
        rawArgPassThrough.add(rawArg);
        String conv = exprEmitter.convert(rawArg, canonTm, actualTm, false, ctx);
        preamble.append('\t').append(userVar).append(" := ").append(conv).append('\n');
        preamble.append("\t_ = ").append(userVar).append('\n');
      }
    }

    TypeMirror canonRetTm = resolver.resolveSigReturnType(sig, ownerSiteTm);
    String retGo = canonRetTm != null ? resolver.goType(canonRetTm, ctx.imports) : "";

    out.append("// Java method: ")
        .append(typeEl.getSimpleName())
        .append('.')
        .append(ee.getSimpleName())
        .append(formatJavaParamsComment(ee))
        .append('\n');
    out.append("func ")
        .append(recvSpec)
        .append(' ')
        .append(methodGoName)
        .append('(')
        .append(String.join(", ", sigParams))
        .append(')');
    if (!retGo.isEmpty()) {
      out.append(' ').append(retGo);
    }
    out.append(" {\n");

    if (resolver.isValueEnum(typeEl)) {
      out.append("\tif this_ == 0 {\n\t\t").append(rtPrefix(ctx)).append("ThrowNPE()\n\t}\n");
    } else if (resolver.hasValueStructValidFlag(typeEl)) {
      out.append("\tif !this_.valid_ {\n\t\t").append(rtPrefix(ctx)).append("ThrowNPE()\n\t}\n");
    } else if (resolver.isNullableValueStructClass(typeEl)) {
      String goTypeName = resolver.planner.getGoTypeName(typeEl);
      String useParams =
          resolver.formatGoTypeParamUse(resolver.getEffectiveTypeParameters(typeEl));
      out.append("\tif this_ == (")
          .append(goTypeName)
          .append(useParams)
          .append("{}) {\n\t\t")
          .append(rtPrefix(ctx))
          .append("ThrowNPE()\n\t}\n");
    }

    if (ee.getModifiers().contains(Modifier.SYNCHRONIZED)) {
      if (resolver.hasFastInstanceSyncMutex(typeEl)) {
        out.append("\tif this_ == nil {\n\t\t").append(rtPrefix(ctx)).append("ThrowNPE()\n\t}\n");
        out.append("\tthis_._syncMu.Lock()\n");
        out.append("\tdefer this_._syncMu.Unlock()\n");
      } else {
        out.append("\tunlock_ := ").append(rtPrefix(ctx)).append("LockMonitor(this_)\n");
        out.append("\tdefer unlock_()\n");
      }
    }

    if (hasEnumOverrideHook) {
      List<String> hookArgs = new ArrayList<>();
      hookArgs.add("this_");
      hookArgs.addAll(rawArgPassThrough);
      out.append("\tif this_.fn_").append(methodGoName).append(" != nil {\n\t\t");
      if (!retGo.isEmpty()) {
        out.append("return ");
      }
      out.append("this_.fn_")
          .append(methodGoName)
          .append('(')
          .append(String.join(", ", hookArgs))
          .append(")\n");
      if (retGo.isEmpty()) {
        out.append("\t\treturn\n");
      }
      out.append("\t}\n");
      out.append('\t')
          .append(retGo.isEmpty() ? "" : "return ")
          .append("this_.baseEnum_")
          .append(methodGoName)
          .append('(')
          .append(String.join(", ", rawArgPassThrough))
          .append(")\n}\n\n");

      out.append("func ")
          .append(recvSpec)
          .append(" baseEnum_")
          .append(methodGoName)
          .append('(')
          .append(String.join(", ", sigParams))
          .append(')');
      if (!retGo.isEmpty()) {
        out.append(' ').append(retGo);
      }
      out.append(" {\n");
    }

    out.append(preamble);

    ctx.expectedReturnType = canonRetTm;
    ctx.returnWidenedToAny = sig.returnWidenedToAny;

    if (isAbstractEnumMethod || isNative || mt.getBody() == null) {
      out.append("\tpanic(\"abstract or native method invoked\")\n");
    } else {
      out.append(
          stmtEmitter.emitBlockStatements(
              mt.getBody(), new TreePath(mtPath, mt.getBody()), 1, ctx));
      if (canonRetTm != null && !endsWithGoTerminatingReturnOrThrow(mt.getBody())) {
        out.append("\tpanic(\"unreachable\")\n");
      }
    }
    out.append("}\n\n");
  }

  private void buildVirtualParamHeaderAndPreamble(
      MethodTree mt,
      TreePath mtPath,
      TypeResolver.MethodSig sig,
      TypeMirror ownerSiteTm,
      List<String> outSigParams,
      StringBuilder outPreamble,
      EmitContext ctx) {
    for (int i = 0; i < mt.getParameters().size(); i++) {
      VariableTree pvt = mt.getParameters().get(i);
      VariableElement pve =
          (VariableElement) resolver.trees.getElement(new TreePath(mtPath, pvt));
      TypeMirror canonTm = resolver.resolveSigParamType(sig, i, ownerSiteTm);
      String canonGo = resolver.goType(canonTm, ctx.imports);
      TypeMirror actualTm = resolver.goEffectiveType(pve.asType());
      String actualGo = resolver.goType(actualTm, ctx.imports);

      String userVar = ctx.registerLocalVar(pve);
      if (canonGo.equals(actualGo) && !sig.paramWidenedToAny[i]) {
        outSigParams.add(userVar + " " + canonGo);
        outPreamble.append("\t_ = ").append(userVar).append('\n');
      } else {
        String rawArg = ctx.nextTemp("arg");
        outSigParams.add(rawArg + " " + canonGo);
        String conv = exprEmitter.convert(rawArg, canonTm, actualTm, false, ctx);
        outPreamble.append('\t').append(userVar).append(" := ").append(conv).append('\n');
        outPreamble.append("\t_ = ").append(userVar).append('\n');
      }
    }
  }

  public static boolean endsWithGoTerminatingReturnOrThrow(BlockTree block) {
    List<? extends StatementTree> stmts = block.getStatements();
    if (stmts.isEmpty()) {
      return false;
    }
    return isGoTerminating(stmts.get(stmts.size() - 1));
  }

  private static boolean isGoTerminating(StatementTree stmt) {
    if (stmt instanceof ReturnTree || stmt instanceof ThrowTree) {
      return true;
    }
    if (stmt instanceof BlockTree bt) {
      return endsWithGoTerminatingReturnOrThrow(bt);
    }
    if (stmt instanceof com.sun.source.tree.IfTree it && it.getElseStatement() != null) {
      return isGoTerminating(it.getThenStatement()) && isGoTerminating(it.getElseStatement());
    }
    return false;
  }

  private String formatJavaParamsComment(ExecutableElement ee) {
    StringBuilder sb = new StringBuilder("(");
    for (int i = 0; i < ee.getParameters().size(); i++) {
      if (i > 0) sb.append(", ");
      VariableElement p = ee.getParameters().get(i);
      sb.append(p.asType().toString()).append(' ').append(p.getSimpleName());
    }
    sb.append(')');
    return sb.toString();
  }

  private String rtPrefix(EmitContext ctx) {
    String a = resolver.runtime.j2gortAlias(ctx.imports);
    return a.isEmpty() ? "" : (a + ".");
  }
}
