package j2go;

import com.sun.source.tree.AssignmentTree;
import com.sun.source.tree.BlockTree;
import com.sun.source.tree.ClassTree;
import com.sun.source.tree.ExpressionStatementTree;
import com.sun.source.tree.ExpressionTree;
import com.sun.source.tree.IdentifierTree;
import com.sun.source.tree.LambdaExpressionTree;
import com.sun.source.tree.LiteralTree;
import com.sun.source.tree.MemberReferenceTree;
import com.sun.source.tree.MemberSelectTree;
import com.sun.source.tree.MethodInvocationTree;
import com.sun.source.tree.MethodTree;
import com.sun.source.tree.NewClassTree;
import com.sun.source.tree.ParenthesizedTree;
import com.sun.source.tree.StatementTree;
import com.sun.source.tree.Tree;
import com.sun.source.tree.UnaryTree;
import com.sun.source.tree.VariableTree;
import com.sun.source.util.TreePath;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import javax.lang.model.element.AnnotationMirror;
import javax.lang.model.element.AnnotationValue;
import javax.lang.model.element.Element;
import javax.lang.model.element.ElementKind;
import javax.lang.model.element.ExecutableElement;
import javax.lang.model.element.Modifier;
import javax.lang.model.element.NestingKind;
import javax.lang.model.element.RecordComponentElement;
import javax.lang.model.element.TypeElement;
import javax.lang.model.element.TypeParameterElement;
import javax.lang.model.element.VariableElement;
import javax.lang.model.type.DeclaredType;
import javax.lang.model.type.ExecutableType;
import javax.lang.model.type.TypeKind;
import javax.lang.model.type.TypeMirror;

/**
 * Emits 1:1 Go source files for Java {@code CompilationUnitTree}s, handling classes, abstract
 * base/interface pairs, interfaces, default methods, lambda adapters, enums (with fields,
 * constructors, abstract methods, and constant-specific bodies), Java 17/21 records, inner
 * classes, local/anonymous classes, static {@code clinit} blocks, and {@code rhino.Node}
 * arena lowering onto {@code transpiled/shims/nodestore}.
 */
public final class ClassEmitter {

  private final TypeResolver resolver;
  private final ExprEmitter exprEmitter;
  private final StmtEmitter stmtEmitter;
  private final MethodEmitter methodEmitter;

  public ClassEmitter(TypeResolver resolver) {
    this.resolver = resolver;
    this.exprEmitter = new ExprEmitter(resolver);
    this.stmtEmitter = new StmtEmitter(resolver, exprEmitter);
    this.exprEmitter.setStmtEmitter(stmtEmitter);
    this.methodEmitter = new MethodEmitter(resolver, exprEmitter, stmtEmitter);
  }

  public String emitCompilationUnit(PackagePlanner.UnitPlan plan, Path srcRoot) throws IOException {
    String rawJava = Files.exists(plan.sourcePath) ? Files.readString(plan.sourcePath) : "";
    String licenseHeader = GoNames.extractSanitizedLicenseHeader(rawJava);
    String relPath =
        srcRoot != null && plan.sourcePath.startsWith(srcRoot)
            ? srcRoot.relativize(plan.sourcePath).toString().replace('\\', '/')
            : plan.sourcePath.getFileName().toString();

    PackagePlanner.ImportTracker imports = new PackagePlanner.ImportTracker(plan.goPkg);
    EmitContext fileCtx = new EmitContext(resolver, plan, imports);
    Set<TypeElement> emittedTypes = new HashSet<>();
    StringBuilder body = new StringBuilder(16384);

    TreePath unitPath = new TreePath(plan.unit);
    for (Tree decl : plan.unit.getTypeDecls()) {
      if (decl instanceof ClassTree ct) {
        emitTypeDecl(ct, new TreePath(unitPath, ct), emittedTypes, body, fileCtx);
      }
    }

    // Drain any method-local and anonymous classes discovered while emitting methods/initializers
    while (!fileCtx.deferredLocalTypes.isEmpty()) {
      List<EmitContext.DeferredLocalType> batch = new ArrayList<>(fileCtx.deferredLocalTypes);
      fileCtx.deferredLocalTypes.clear();
      for (EmitContext.DeferredLocalType def : batch) {
        if (def.element() != null && !emittedTypes.contains(def.element())) {
          emitTypeDecl(def.tree(), def.path(), emittedTypes, body, fileCtx);
        }
      }
    }

    StringBuilder hoisted = new StringBuilder();
    fileCtx.appendHoistedDecls(hoisted);

    String bodyStr = hoisted.toString() + body.toString();
    StringBuilder full = new StringBuilder(bodyStr.length() + 512);
    if (!licenseHeader.isEmpty()) {
      full.append(licenseHeader).append('\n');
    }
    full.append("// Port of ").append(relPath).append(" (generated mechanically by j2go).\n\n");
    full.append("package ").append(plan.goPkg).append("\n\n");
    full.append(imports.emitImportBlock(bodyStr));
    full.append(bodyStr);
    return full.toString();
  }

  private void emitTypeDecl(
      ClassTree ct,
      TreePath ctPath,
      Set<TypeElement> emittedTypes,
      StringBuilder out,
      EmitContext fileCtx) {
    TypeElement te = resolver.getClassTreeElement(ct);
    if (te == null) {
      Element el = resolver.trees.getElement(ctPath);
      if (el instanceof TypeElement t) {
        te = t;
      }
    }
    if (te == null || resolver.isEnumConstantAnonType(te) || !emittedTypes.add(te)) {
      return;
    }

    switch (te.getKind()) {
      case ANNOTATION_TYPE -> {
        out.append("// Java @interface ").append(te.getQualifiedName()).append(" (erased at runtime)\n\n");
        emitStaticFieldsAndClinit(ct, ctPath, te, Map.of(), out, fileCtx);
        emitNestedMemberTypes(ct, ctPath, emittedTypes, out, fileCtx);
      }
      case INTERFACE -> emitInterface(ct, ctPath, te, emittedTypes, out, fileCtx);
      case ENUM -> emitEnum(ct, ctPath, te, emittedTypes, out, fileCtx);
      case CLASS, RECORD -> {
        if (resolver.isNodeClass(te)) {
          emitNodeType(ct, ctPath, te, emittedTypes, out, fileCtx);
        } else {
          emitClassOrRecord(ct, ctPath, te, emittedTypes, out, fileCtx);
        }
      }
      default -> {}
    }
  }

  private void emitNestedMemberTypes(
      ClassTree ct,
      TreePath ctPath,
      Set<TypeElement> emittedTypes,
      StringBuilder out,
      EmitContext fileCtx) {
    for (Tree member : ct.getMembers()) {
      if (member instanceof ClassTree nestedCt) {
        emitTypeDecl(nestedCt, new TreePath(ctPath, nestedCt), emittedTypes, out, fileCtx);
      }
    }
  }

  // =========================================================================
  // 1. Interfaces (Default Methods, Static Members, Lambda Adapters)
  // =========================================================================

  private void emitInterface(
      ClassTree ct,
      TreePath ctPath,
      TypeElement te,
      Set<TypeElement> emittedTypes,
      StringBuilder out,
      EmitContext fileCtx) {
    String goName = resolver.planner.getGoTypeName(te);
    StringBuilder staticOut = new StringBuilder();
    boolean hasClinit = emitStaticFieldsAndClinit(ct, ctPath, te, Map.of(), staticOut, fileCtx);
    boolean hasMethodEntryClinit =
        hasClinit && !isPureStaticFieldOnlyClinit(resolver, ct, ctPath, te);

    List<TypeParameterElement> effTypeParams = resolver.getEffectiveTypeParameters(te);
    boolean hasEffTypeParams = !effTypeParams.isEmpty();
    if (hasEffTypeParams) {
      resolver.pushTypeParamScope(te, effTypeParams);
    }
    try {
      String declParams = resolver.formatGoTypeParamDecl(effTypeParams);
      String useParams = resolver.formatGoTypeParamUse(effTypeParams);

      out.append("// Java interface: ").append(te.getQualifiedName()).append('\n');
      out.append("type ").append(goName).append(declParams).append(" interface {\n");
      out.append('\t').append(resolver.getInterfaceMarkerMethod(te)).append("()\n");
      out.append("\tIsJavaNil_() bool\n");

      for (TypeMirror im : te.getInterfaces()) {
        if (im instanceof DeclaredType dt
            && dt.asElement() instanceof TypeElement ste
            && !resolver.runtime.mapsToAny(ste)) {
          out.append('\t').append(resolver.goType(im, fileCtx.imports)).append('\n');
        }
      }

      Set<String> inheritedNames = collectEmbeddedAncestorMethodGoNames(te);
      Set<String> emittedSigNames = new HashSet<>(inheritedNames);
      emittedSigNames.add("IsJavaNil_");
      for (Element m : te.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)
            && !ee.getModifiers().contains(Modifier.PRIVATE)) {
          TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
          if (sig != null && emittedSigNames.add(sig.goName)) {
            out.append('\t').append(formatInterfaceMethodSpec(te, sig, fileCtx.imports)).append('\n');
          }
        }
      }
      out.append("}\n");
      emitToIfaceHelper(goName, declParams, useParams, out);

      if (hasEffTypeParams) {
        emitMarkerTypePredicateRegistration(te, out, fileCtx);
      } else {
        emitTypePredicateRegistration(te, goName, out, fileCtx);
      }

      out.append(staticOut);

      for (Tree member : ct.getMembers()) {
        if (member instanceof MethodTree mt) {
          methodEmitter.emitMethod(
              mt,
              new TreePath(ctPath, mt),
              te,
              hasMethodEntryClinit,
              null,
              Map.of(),
              out,
              fileCtx);
        }
      }

      if (resolver.getLambdaTargetInterfaces().contains(te)) {
        emitSourceLambdaAdapter(te, declParams, useParams, out, fileCtx);
      }
    } finally {
      if (hasEffTypeParams) {
        resolver.popTypeParamScope();
      }
    }

    emitNestedMemberTypes(ct, ctPath, emittedTypes, out, fileCtx);
  }

  private void emitToIfaceHelper(
      String goName, String declParams, String useParams, StringBuilder out) {
    String ifaceWithParams = goName + useParams;
    if (declParams.isEmpty()) {
      out.append("func ToIface_")
          .append(goName)
          .append("[T interface {\n\t")
          .append(ifaceWithParams)
          .append("\n\tcomparable\n}](v T) ")
          .append(ifaceWithParams)
          .append(" {\n\tvar z T\n\tif v == z {\n\t\treturn nil\n\t}\n\treturn v\n}\n\n");
      out.append("func Upcast_")
          .append(goName)
          .append("[T ")
          .append(goName)
          .append("](v T) ")
          .append(goName)
          .append(" {\n\tvar vi ")
          .append(goName)
          .append(" = v\n\tif vi == nil || vi.IsJavaNil_() {\n\t\treturn nil\n\t}\n\treturn vi\n}\n\n");
    } else {
      String declInner = declParams.substring(1, declParams.length() - 1);
      out.append("func ToIface_")
          .append(goName)
          .append('[')
          .append(declInner)
          .append(", T_iface_ interface {\n\t")
          .append(ifaceWithParams)
          .append("\n\tcomparable\n}](v T_iface_) ")
          .append(ifaceWithParams)
          .append(" {\n\tvar z T_iface_\n\tif v == z {\n\t\treturn nil\n\t}\n\treturn v\n}\n\n");
    }
  }

  private void emitSourceLambdaAdapter(
      TypeElement fi,
      String declParams,
      String useParams,
      StringBuilder out,
      EmitContext fileCtx) {
    ExecutableElement sam = resolver.findSamMethod(fi);
    if (sam == null) {
      return;
    }
    TypeResolver.MethodSig samSig = resolver.getInstanceMethodSig(sam);
    String goName = resolver.planner.getGoTypeName(fi);
    String ifaceWithParams = goName + useParams;
    String structName = "lambda_" + goName;
    String structWithParams = structName + useParams;
    TypeMirror fiSiteTm = fi.asType();

    List<String> fnParamTypes = new ArrayList<>();
    List<String> namedParams = new ArrayList<>();
    List<String> callArgs = new ArrayList<>();
    for (int i = 0; i < samSig.canonParamTypes.size(); i++) {
      TypeMirror ptTm = resolver.resolveSigParamType(samSig, i, fiSiteTm);
      String pt = resolver.goType(ptTm, fileCtx.imports);
      fnParamTypes.add(pt);
      String pn = "p" + i + "_";
      namedParams.add(pn + " " + pt);
      callArgs.add(pn);
    }
    TypeMirror retTm = resolver.resolveSigReturnType(samSig, fiSiteTm);
    String retGo = retTm != null ? resolver.goType(retTm, fileCtx.imports) : "";
    String fnType =
        "func(" + String.join(", ", fnParamTypes) + ")" + (retGo.isEmpty() ? "" : (" " + retGo));

    out.append("type ")
        .append(structName)
        .append(declParams)
        .append(" struct {\n\tfn ")
        .append(fnType)
        .append("\n}\n\n");
    out.append("func (l *")
        .append(structWithParams)
        .append(") IsJavaNil_() bool { return l == nil }\n");
    out.append("func NewLambda_")
        .append(goName)
        .append(declParams)
        .append("(fn ")
        .append(fnType)
        .append(") ")
        .append(ifaceWithParams)
        .append(" {\n\tif fn == nil {\n\t\treturn nil\n\t}\n\treturn &")
        .append(structWithParams)
        .append("{fn: fn}\n}\n\n");

    emitMarkerMethodsForType(fi, structWithParams, true, out);

    // SAM method implementation
    out.append("func (l *")
        .append(structWithParams)
        .append(") ")
        .append(samSig.goName)
        .append('(')
        .append(String.join(", ", namedParams))
        .append(')')
        .append(retGo.isEmpty() ? "" : (" " + retGo))
        .append(" {\n\t")
        .append(retGo.isEmpty() ? "" : "return ")
        .append("l.fn(")
        .append(String.join(", ", callArgs))
        .append(")\n}\n\n");

    // Forward any default methods (or non-SAM interface methods) so `*lambda_<goName>` satisfies `goName`
    Set<String> handled = new HashSet<>();
    handled.add(samSig.goName);

    List<TypeElement> allIfaces = new ArrayList<>();
    allIfaces.add(fi);
    for (TypeElement anc : resolver.getAllAncestorTypes(fi)) {
      if (anc.getKind() == ElementKind.INTERFACE && !resolver.runtime.mapsToAny(anc)) {
        allIfaces.add(anc);
      }
    }

    for (TypeElement iface : allIfaces) {
      for (Element m : iface.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)
            && !ee.getModifiers().contains(Modifier.PRIVATE)) {
          TypeResolver.MethodSig dSig = resolver.getInstanceMethodSig(ee);
          if (dSig != null
              && (resolver.isSourceType(iface) || resolver.isNeededExternalSig(dSig))
              && handled.add(dSig.goName)) {
            emitDefaultOrFallbackMethodOnStruct(
                fi, structWithParams, false, iface, ee, dSig, out, fileCtx);
          }
        }
      }
    }
  }

  // =========================================================================
  // 2. Classes & Java 17/21 Records
  // =========================================================================

  private void emitClassOrRecord(
      ClassTree ct,
      TreePath ctPath,
      TypeElement te,
      Set<TypeElement> emittedTypes,
      StringBuilder out,
      EmitContext fileCtx) {
    String goName = resolver.planner.getGoTypeName(te);
    boolean isSubclassed = resolver.isSubclassed(te);
    String structName = isSubclassed ? (goName + "Base") : goName;
    TypeElement superTe = resolver.getSuperclassElement(te);

    Map<VariableElement, String> capturedFields = buildCapturedFieldMap(te);

    String qn = te.getQualifiedName().toString();
    boolean hasTcHandle =
        qn.equals("com.google.javascript.rhino.jstype.JSType")
            || qn.equals("com.google.javascript.jscomp.colors.Color");
    boolean hasSfHandle =
        implementsInterface(te, "com.google.javascript.rhino.StaticSourceFile")
            && (superTe == null
                || !implementsInterface(superTe, "com.google.javascript.rhino.StaticSourceFile"));
    boolean hasJsDocHandle = qn.equals("com.google.javascript.rhino.JSDocInfo");
    boolean isLazyInternedStringList =
        qn.equals("com.google.javascript.rhino.RhinoStringPool.LazyInternedStringList");
    String handleAtomicPkg =
        (hasTcHandle || hasSfHandle || hasJsDocHandle || isLazyInternedStringList)
            ? fileCtx.imports.std("sync/atomic")
            : null;

    List<VariableElement> instanceFields = new ArrayList<>();
    for (Element m : te.getEnclosedElements()) {
      if (m instanceof VariableElement ve
          && ve.getKind() == ElementKind.FIELD
          && !ve.getModifiers().contains(Modifier.STATIC)) {
        instanceFields.add(ve);
      }
    }

    // 4. Emit static fields & `clinit_<GoName>` into buffer so `hasClinit` is known for constructors/methods
    StringBuilder staticOut = new StringBuilder();
    boolean hasClinit = emitStaticFieldsAndClinit(ct, ctPath, te, capturedFields, staticOut, fileCtx);
    boolean hasMethodEntryClinit =
        hasClinit && !isPureStaticFieldOnlyClinit(resolver, ct, ctPath, te);

    List<TypeParameterElement> effTypeParams = resolver.getEffectiveTypeParameters(te);
    boolean hasEffTypeParams = !effTypeParams.isEmpty();
    if (hasEffTypeParams) {
      resolver.pushTypeParamScope(te, effTypeParams);
    }
    try {
      String declParams = resolver.formatGoTypeParamDecl(effTypeParams);
      String useParams = resolver.formatGoTypeParamUse(effTypeParams);
      String structWithParams = structName + useParams;
      String ifaceWithParams = goName + useParams;
      NodePredInfo nodePred = inspectNodePredicateImpl(te);

      String superStructDecl = null;
      String superEmbeddedField = null;
      if (superTe != null && !resolver.runtime.mapsToAny(superTe)) {
        String rawSuperStruct = resolver.goStructTypeForElement(superTe, fileCtx.imports);
        int dot = rawSuperStruct.lastIndexOf('.');
        superEmbeddedField = dot >= 0 ? rawSuperStruct.substring(dot + 1) : rawSuperStruct;
        superStructDecl =
            resolver.goEmbeddedSuperStructType(te.getSuperclass(), superTe, fileCtx.imports);
      }

      // 1. If subclassed or abstract, emit polymorphic Go interface `<GoName>`
      if (isSubclassed) {
        out.append("// Java polymorphic class interface: ")
            .append(te.getQualifiedName())
            .append('\n');
        out.append("type ").append(goName).append(declParams).append(" interface {\n");
        out.append('\t').append(resolver.getInterfaceMarkerMethod(te)).append("()\n");
        out.append("\tIsJavaNil_() bool\n");
        out.append("\tGet").append(goName).append("Base() *").append(structWithParams).append('\n');

        if (superTe != null && !resolver.runtime.mapsToAny(superTe)) {
          out.append('\t').append(resolver.goType(te.getSuperclass(), fileCtx.imports)).append('\n');
        }
        for (TypeMirror im : te.getInterfaces()) {
          if (im instanceof DeclaredType dt
              && dt.asElement() instanceof TypeElement ite
              && !resolver.runtime.mapsToAny(ite)) {
            out.append('\t').append(resolver.goType(im, fileCtx.imports)).append('\n');
          }
        }

        Set<String> inheritedNames = collectEmbeddedAncestorMethodGoNames(te);
        Set<String> emittedIfaceSigs = new HashSet<>(inheritedNames);
        emittedIfaceSigs.add("IsJavaNil_");
        // Include cross-package field getters/setters on the polymorphic interface so cross-pkg subclasses/callers work!
        for (Element m : te.getEnclosedElements()) {
          if (m instanceof VariableElement ve
              && ve.getKind() == ElementKind.FIELD
              && !ve.getModifiers().contains(Modifier.STATIC)) {
            String fGo = resolver.getFieldGoName(ve);
            String fType = resolver.goType(ve.asType(), fileCtx.imports);
            if (resolver.isFieldCrossPkgRead(ve) && emittedIfaceSigs.add("FieldGet_" + fGo)) {
              out.append("\tFieldGet_").append(fGo).append("() ").append(fType).append('\n');
            }
            if (resolver.isFieldCrossPkgWrite(ve) && emittedIfaceSigs.add("FieldSet_" + fGo)) {
              out.append("\tFieldSet_").append(fGo).append('(').append(fType).append(")\n");
            }
          }
        }
        for (Element m : te.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && !ee.getModifiers().contains(Modifier.STATIC)
              && !ee.getModifiers().contains(Modifier.PRIVATE)) {
            TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
            if (sig != null && emittedIfaceSigs.add(sig.goName)) {
              out.append('\t')
                  .append(formatInterfaceMethodSpec(te, sig, fileCtx.imports))
                  .append('\n');
            }
          }
        }
        if (nodePred != null && emittedIfaceSigs.add("TestNode_")) {
          out.append("\tTestNode_(uint32) bool\n");
        }
        out.append("}\n");
        emitToIfaceHelper(goName, declParams, useParams, out);
      }

      // Emit steps 5-9b into `methodsOut` BEFORE step 10 so we can detect whether `this_.self_` is ever read!
      StringBuilder methodsOut = new StringBuilder();

      // 5. Emit `initFields_<GoName>` (instance field initializers + non-static `{ ... }` blocks in source order)
      emitInstanceFieldsInit(
          ct, ctPath, te, goName, structWithParams, capturedFields, methodsOut, fileCtx);

      // 6. Emit constructors
      if (te.getNestingKind() == NestingKind.ANONYMOUS) {
        emitAnonymousClassConstructor(
            ctPath,
            te,
            goName,
            declParams,
            structWithParams,
            superTe,
            superEmbeddedField,
            capturedFields,
            methodsOut,
            fileCtx);
      } else {
        boolean hasCtorTree = false;
        for (Tree member : ct.getMembers()) {
          if (member instanceof MethodTree mt && isConstructorTree(mt)) {
            hasCtorTree = true;
            methodEmitter.emitConstructor(
                mt,
                new TreePath(ctPath, mt),
                te,
                hasMethodEntryClinit,
                capturedFields,
                methodsOut,
                fileCtx);
          }
        }
        if (!hasCtorTree) {
          emitSynthesizedDefaultConstructor(
              te,
              goName,
              declParams,
              useParams,
              structWithParams,
              superTe,
              superEmbeddedField,
              hasMethodEntryClinit,
              capturedFields,
              methodsOut,
              fileCtx);
        }
      }

      // 7. Emit methods declared on `ct`
      Set<String> emittedMethodGoNames = new HashSet<>();
      boolean emittedTestNode = false;
      for (Tree member : ct.getMembers()) {
        if (member instanceof MethodTree mt && !isConstructorTree(mt)) {
          TreePath mtPath = new TreePath(ctPath, mt);
          Element el = resolver.trees.getElement(mtPath);
          if (el instanceof ExecutableElement ee
              && !ee.getModifiers().contains(Modifier.STATIC)) {
            TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
            if (!ee.getModifiers().contains(Modifier.ABSTRACT) && sig != null) {
              emittedMethodGoNames.add(sig.goName);
            }
            if (nodePred != null
                && !emittedTestNode
                && !ee.getModifiers().contains(Modifier.PRIVATE)
                && mt.getParameters().size() == 1
                && sig != null
                && sig.goName.equals(nodePred.primarySamGoName())) {
              emitNodePredicateMethodAndBridge(
                  mt, mtPath, te, ee, sig, nodePred, false, capturedFields, methodsOut, fileCtx);
              emittedTestNode = true;
              continue;
            }
          }
          methodEmitter.emitMethod(
              mt,
              mtPath,
              te,
              hasMethodEntryClinit,
              null,
              capturedFields,
              methodsOut,
              fileCtx);
        }
      }

      // 8. For Java `record`s, emit any synthesized component accessors + value Equals/HashCode/String not explicitly declared
      if (te.getKind() == ElementKind.RECORD) {
        emitSynthesizedRecordMethods(
            te, structWithParams, emittedMethodGoNames, methodsOut, fileCtx);
      }

      // 9. Emit forwarding methods for inherited interface `default` methods not already implemented on `te` or superclasses
      emitInheritedInterfaceDefaultForwarders(
          te, structWithParams, isSubclassed, emittedMethodGoNames, methodsOut, fileCtx);

      if (nodePred != null && !emittedTestNode) {
        emitFallbackOrAbstractNodePredicateMethod(
            te,
            goName,
            useParams,
            structWithParams,
            superTe,
            nodePred,
            emittedMethodGoNames,
            methodsOut,
            fileCtx);
      }

      // 9b. Emit mechanical org.kohsuke.args4j @Option / @Argument registration if present
      if (!hasEffTypeParams) {
        Args4jEmitter.emitArgs4jRegistrationIfNeeded(
            te, structName, resolver, methodsOut, fileCtx);
      }

      TypeElement reusableSelfRoot =
          isSubclassed ? resolver.getSamePkgReusableSelfRoot(te) : null;
      String selfHelperName = "self_" + goName + "_";
      boolean needsSelfHelper =
          reusableSelfRoot != null && methodsOut.indexOf("this_." + selfHelperName + "()") >= 0;
      boolean needsSelfField =
          isSubclassed
              && reusableSelfRoot == null
              && (methodsOut.indexOf("this_.self_") >= 0
                  || resolver.isReusableSelfAncestor(te));
      if (needsSelfField && !resolver.mayHaveSelfField(te)) {
        throw new IllegalStateException(
            "mayHaveSelfField under-approximated for " + te.getQualifiedName());
      }
      boolean hasOuterField = resolver.hasOuterInstance(te);
      boolean hasFastInstSync = resolver.hasFastInstanceSyncMutex(te);
      TypeElement rootSrcTe = rootSourceClass(te);
      boolean rootHasIdHash = classHasIdentityHashField(rootSrcTe);
      boolean hasIdHash = rootHasIdHash && te.equals(rootSrcTe);
      boolean inheritsIdHash =
          rootHasIdHash && !te.equals(rootSrcTe) && superEmbeddedField != null;
      if (hasIdHash && handleAtomicPkg == null) {
        handleAtomicPkg = fileCtx.imports.std("sync/atomic");
      }
      int handleCount = (hasTcHandle ? 1 : 0) + (hasSfHandle ? 1 : 0) + (hasJsDocHandle ? 1 : 0);
      boolean packIdHashWithHandles = hasIdHash && (handleCount % 2 == 1);
      boolean hasAnyStructField =
          hasTcHandle
              || hasSfHandle
              || hasJsDocHandle
              || isLazyInternedStringList
              || hasIdHash
              || hasFastInstSync
              || needsSelfField
              || hasOuterField
              || !capturedFields.isEmpty()
              || superEmbeddedField != null
              || !instanceFields.isEmpty();

      // 2. Emit concrete / base struct `<StructName>`
      out.append("// Java ")
          .append(te.getKind() == ElementKind.RECORD ? "record" : "class")
          .append(": ")
          .append(te.getQualifiedName())
          .append('\n');
      out.append("type ").append(structName).append(declParams).append(" struct {\n");
      if (!hasAnyStructField) {
        out.append("\t_id byte\n");
      } else if (shouldSkipStructFieldPacking(te, instanceFields)) {
        if (hasFastInstSync) {
          String syncPkg = fileCtx.imports.std("sync");
          out.append("\t_syncMu ").append(syncPkg).append(".Mutex `gson:\"-\"`\n");
        }
        if (hasTcHandle) {
          out.append("\t_tcHandle ").append(handleAtomicPkg).append(".Uint32\n");
        }
        if (hasSfHandle) {
          out.append("\t_sfHandle ").append(handleAtomicPkg).append(".Uint32\n");
        }
        if (hasJsDocHandle) {
          out.append("\t_jsDocHandle ").append(handleAtomicPkg).append(".Uint32\n");
        }
        if (packIdHashWithHandles) {
          out.append("\tidHash_ ").append(handleAtomicPkg).append(".Int32\n");
        }
        if (needsSelfField) {
          out.append("\tself_ ").append(ifaceWithParams).append('\n');
        }
        if (hasOuterField) {
          TypeElement outerTe = resolver.getEnclosingInstanceType(te);
          out.append("\touter_ ")
              .append(resolver.goTypeForElement(outerTe, fileCtx.imports))
              .append('\n');
        }
        for (Map.Entry<VariableElement, String> cap : capturedFields.entrySet()) {
          out.append('\t')
              .append(cap.getValue())
              .append(' ')
              .append(resolver.goType(cap.getKey().asType(), fileCtx.imports))
              .append('\n');
        }
        if (superStructDecl != null) {
          out.append('\t').append(superStructDecl).append('\n');
        }
        for (VariableElement ve : instanceFields) {
          out.append('\t')
              .append(resolver.getFieldGoName(ve))
              .append(' ')
              .append(resolver.goType(ve.asType(), fileCtx.imports))
              .append(' ')
              .append(gsonStructTag(ve))
              .append('\n');
        }
        if (isLazyInternedStringList) {
          out.append("\tids_ []uint32 `gson:\"-\"`\n");
          out.append("\tepoch_ uint64 `gson:\"-\"`\n");
        }
        if (hasIdHash && !packIdHashWithHandles) {
          out.append("\tidHash_ ").append(handleAtomicPkg).append(".Int32\n");
        }
        if (resolver.hasValueStructValidFlag(te)) {
          out.append("\tvalid_ bool `gson:\"-\"`\n");
        }
      } else {
        // Keep embedded base struct first, then sort declared fields stably by pointer-containment
        // and alignment descending:
        //   0: align 8 containing pointers (interfaces, pointers, slices, string, maps, funcs, typevars)
        //   1: align 8 pointer-free (int64, uint64, float64)
        //   2: align 4 pointer-free (sync.Mutex, atomic.{Int32,Uint32}, int32, uint32, float32, Node)
        //   3: align 2 pointer-free (int16, uint16)
        //   4: align 1 pointer-free (bool, int8, uint8, byte-sized value enums)
        if (superStructDecl != null) {
          out.append('\t').append(superStructDecl).append('\n');
        }
        @SuppressWarnings({"unchecked", "rawtypes"})
        List<String>[] buckets = new List[5];
        for (int b = 0; b < 5; b++) {
          buckets[b] = new ArrayList<>();
        }
        if (needsSelfField) {
          buckets[0].add("\tself_ " + ifaceWithParams + "\n");
        }
        if (hasOuterField) {
          TypeElement outerTe = resolver.getEnclosingInstanceType(te);
          buckets[0].add(
              "\touter_ " + resolver.goTypeForElement(outerTe, fileCtx.imports) + "\n");
        }
        if (hasFastInstSync) {
          String syncPkg = fileCtx.imports.std("sync");
          buckets[2].add("\t_syncMu " + syncPkg + ".Mutex `gson:\"-\"`\n");
        }
        if (hasTcHandle) {
          buckets[2].add("\t_tcHandle " + handleAtomicPkg + ".Uint32\n");
        }
        if (hasSfHandle) {
          buckets[2].add("\t_sfHandle " + handleAtomicPkg + ".Uint32\n");
        }
        if (hasJsDocHandle) {
          buckets[2].add("\t_jsDocHandle " + handleAtomicPkg + ".Uint32\n");
        }
        for (Map.Entry<VariableElement, String> cap : capturedFields.entrySet()) {
          String capGoType = resolver.goType(cap.getKey().asType(), fileCtx.imports);
          int b = classifyFieldPackingBucket(cap.getKey().asType(), null);
          buckets[b].add("\t" + cap.getValue() + " " + capGoType + "\n");
        }
        for (VariableElement ve : instanceFields) {
          String fGoName = resolver.getFieldGoName(ve);
          String fGoType = resolver.goType(ve.asType(), fileCtx.imports);
          String fTag = gsonStructTag(ve);
          int b = classifyFieldPackingBucket(ve.asType(), null);
          buckets[b].add("\t" + fGoName + " " + fGoType + " " + fTag + "\n");
        }
        if (isLazyInternedStringList) {
          buckets[0].add("\tids_ []uint32 `gson:\"-\"`\n");
          buckets[1].add("\tepoch_ uint64 `gson:\"-\"`\n");
        }
        if (hasIdHash) {
          buckets[2].add("\tidHash_ " + handleAtomicPkg + ".Int32\n");
        }
        if (resolver.hasValueStructValidFlag(te)) {
          buckets[4].add("\tvalid_ bool `gson:\"-\"`\n");
        }
        for (int b = 0; b < 5; b++) {
          for (String line : buckets[b]) {
            out.append(line);
          }
        }
      }
      out.append("}\n\n");

      // 3. Emit `BindSelf_`, `Get<GoName>Base`, `SetOuter_`, marker methods, and cross-pkg field accessors
      boolean isValueStruct = resolver.isValueStructClass(te);
      emitBindSelfMethod(
          te,
          ifaceWithParams,
          structWithParams,
          needsSelfField,
          superTe,
          superEmbeddedField,
          out,
          fileCtx);
      emitJavaClassInfoMethod(te, structWithParams, out);
      if (hasTcHandle) {
        out.append("func (this_ *")
            .append(structWithParams)
            .append(") TCHandlePtr_() *")
            .append(handleAtomicPkg)
            .append(".Uint32 { return &this_._tcHandle }\n");
        if (!isSubclassed && qn.equals("com.google.javascript.jscomp.colors.Color")) {
          out.append("func (*").append(structWithParams).append(") Implements_colors_Color() {}\n");
        }
      }
      if (hasSfHandle) {
        out.append("func (this_ *")
            .append(structWithParams)
            .append(") SourceFileHandlePtr_() *")
            .append(handleAtomicPkg)
            .append(".Uint32 { return &this_._sfHandle }\n");
      }
      if (hasJsDocHandle) {
        out.append("func (this_ *")
            .append(structWithParams)
            .append(") JSDocHandlePtr_() *")
            .append(handleAtomicPkg)
            .append(".Uint32 { return &this_._jsDocHandle }\n");
      }
      if (hasIdHash) {
        String rt = rtPrefix(fileCtx);
        out.append("func (this_ *")
            .append(structWithParams)
            .append(") IdentityHash_() int32 {\n\tif this_ == nil {\n\t\treturn 0\n\t}\n\tif h := this_.idHash_.Load(); h != 0 {\n\t\treturn h\n\t}\n\treturn ")
            .append(rt)
            .append("AssignIdentHash(&this_.idHash_)\n}\n");
      } else if (inheritsIdHash) {
        out.append("func (this_ *")
            .append(structWithParams)
            .append(") IdentityHash_() int32 {\n\tif this_ == nil {\n\t\treturn 0\n\t}\n\treturn this_.")
            .append(superEmbeddedField)
            .append(".IdentityHash_()\n}\n");
      }
      if (hasIdentityEquals(te, isSubclassed)) {
        out.append("func (*").append(structWithParams).append(") HasIdentityEquals_() {}\n");
      } else if (isValueStruct && !hasEffTypeParams && exprEmitter.valueStructHasIdentityEquals(te)) {
        out.append("func (").append(structWithParams).append(") HasIdentityEquals_() {}\n");
        emitValueStructIdentityHashMethod(
            te, structWithParams, hasOuterField, instanceFields, out, fileCtx);
        if (resolver.isNullableValueStructClass(te)) {
          out.append("func (this_ ")
              .append(structWithParams)
              .append(") FastValHash_() (int32, bool, bool) {\n\tif this_.IsJavaNil_() {\n\t\treturn 0, false, true\n\t}\n\treturn this_.IdentityHash_(), true, true\n}\n");
        } else {
          out.append("func (this_ ")
              .append(structWithParams)
              .append(") FastValHash_() (int32, bool, bool) { return this_.IdentityHash_(), true, true }\n");
        }
      } else if (isValueStruct
          && !hasEffTypeParams
          && qn.equals("com.google.javascript.jscomp.colors.ColorId")) {
        String rt = rtPrefix(fileCtx);
        out.append("func (this_ ")
            .append(structWithParams)
            .append(") FastValHash_() (int32, bool, bool) {\n\tif !this_.valid_ {\n\t\treturn 0, false, false\n\t}\n\treturn int32(this_.rightAligned), true, false\n}\n");
        out.append("func (this_ ")
            .append(structWithParams)
            .append(") RawColorId_() (int64, bool) { return this_.rightAligned, this_.valid_ }\n");
        out.append("func init() {\n\t")
            .append(rt)
            .append("RegisterColorIdBoxer(func(raw_ int64) any { return ")
            .append(structWithParams)
            .append("{rightAligned: raw_, valid_: true} })\n}\n");
      }
      if (hasEffTypeParams) {
        String genMarker = resolver.getGenericClassMarkerMethod(te);
        out.append("func (*")
            .append(structWithParams)
            .append(") ")
            .append(genMarker)
            .append("() {}\n");
        emitGenericTypePredicateRegistration(te, genMarker, out, fileCtx);
      } else {
        emitTypePredicateRegistration(
            te, (isSubclassed || isValueStruct) ? goName : ("*" + structName), out, fileCtx);
      }
      if (isSubclassed) {
        out.append("func (b *")
            .append(structWithParams)
            .append(") Get")
            .append(goName)
            .append("Base() *")
            .append(structWithParams)
            .append(" { return b }\n");
      }
      if (needsSelfHelper) {
        String rootBaseName = resolver.planner.getGoTypeName(reusableSelfRoot) + "Base";
        out.append("func (this_ *")
            .append(structWithParams)
            .append(") ")
            .append(selfHelperName)
            .append("() ")
            .append(ifaceWithParams)
            .append(" {\n\tif s := this_.")
            .append(rootBaseName)
            .append(".self_; s != nil {\n\t\treturn s.(")
            .append(ifaceWithParams)
            .append(")\n\t}\n\treturn nil\n}\n");
      }
      if (hasOuterField) {
        TypeElement outerTe = resolver.getEnclosingInstanceType(te);
        String outerGo = resolver.goTypeForElement(outerTe, fileCtx.imports);
        out.append("func (this_ *")
            .append(structWithParams)
            .append(") SetOuter_(o ")
            .append(outerGo)
            .append(") { this_.outer_ = o }\n");
        out.append("func (this_ *")
            .append(structWithParams)
            .append(") GetOuter_() ")
            .append(outerGo)
            .append(" { return this_.outer_ }\n");
      }
      emitMarkerMethodsForType(te, structWithParams, isSubclassed, out);
      emitCrossPkgFieldAccessors(te, structWithParams, isValueStruct, instanceFields, out, fileCtx);
      if (isLazyInternedStringList) {
        String syncPkg = fileCtx.imports.std("sync");
        String ns = resolver.runtime.nodestoreAlias(fileCtx.imports);
        String nsPref = ns.isEmpty() ? "" : (ns + ".");
        String rt = rtPrefix(fileCtx);
        out.append("var lazyStringPoolEpochMu_ ").append(syncPkg).append(".Mutex\n\n");
        out.append("func (this_ *")
            .append(structWithParams)
            .append(") InternID_(offset int32) ")
            .append(nsPref)
            .append("StringID {\n")
            .append("\treturn internIDFast_(this_, offset, (*")
            .append(structWithParams)
            .append(").internIDSlow_)\n}\n\n");
        out.append("func internIDFast_(this_ *")
            .append(structWithParams)
            .append(", offset int32, slow func(*")
            .append(structWithParams)
            .append(", int32) ")
            .append(nsPref)
            .append("StringID) ")
            .append(nsPref)
            .append("StringID {\n")
            .append("\tif this_ != nil {\n")
            .append("\t\tif ids := this_.ids_; uint(offset) < uint(len(ids)) && ")
            .append(handleAtomicPkg)
            .append(".LoadUint64(&this_.epoch_) == ")
            .append(nsPref)
            .append("Epoch() {\n")
            .append("\t\t\tif id := ")
            .append(handleAtomicPkg)
            .append(".LoadUint32(&ids[offset]); id != 0 {\n")
            .append("\t\t\t\treturn ")
            .append(nsPref)
            .append("StringID(id)\n")
            .append("\t\t\t}\n")
            .append("\t\t}\n")
            .append("\t}\n")
            .append("\treturn slow(this_, offset)\n}\n\n");
        out.append("//go:noinline\n")
            .append("func (this_ *")
            .append(structWithParams)
            .append(") internIDSlow_(offset int32) ")
            .append(nsPref)
            .append("StringID {\n")
            .append("\tif this_ == nil {\n\t\t")
            .append(rt)
            .append("ThrowNPE()\n\t}\n")
            .append("\tif offset < 0 || offset >= this_.size {\n\t\t_ = this_.Get(offset)\n\t}\n")
            .append("\tep := ")
            .append(nsPref)
            .append("Epoch()\n")
            .append("\tif ")
            .append(handleAtomicPkg)
            .append(".LoadUint64(&this_.epoch_) != ep {\n")
            .append("\t\tlazyStringPoolEpochMu_.Lock()\n")
            .append("\t\tif ep = ")
            .append(nsPref)
            .append("Epoch(); ")
            .append(handleAtomicPkg)
            .append(".LoadUint64(&this_.epoch_) != ep {\n")
            .append("\t\t\tfor i := range this_.ids_ {\n")
            .append("\t\t\t\t")
            .append(handleAtomicPkg)
            .append(".StoreUint32(&this_.ids_[i], 0)\n")
            .append("\t\t\t}\n")
            .append("\t\t\t")
            .append(handleAtomicPkg)
            .append(".StoreUint64(&this_.epoch_, ep)\n")
            .append("\t\t}\n")
            .append("\t\tlazyStringPoolEpochMu_.Unlock()\n")
            .append("\t}\n")
            .append("\ts := this_.pool[offset]\n")
            .append("\tif s == ")
            .append(rt)
            .append("NullString {\n\t\treturn 0\n\t}\n")
            .append("\tid := ")
            .append(nsPref)
            .append("InternOwned(s)\n")
            .append("\tif int(offset) < len(this_.ids_) {\n")
            .append("\t\t")
            .append(handleAtomicPkg)
            .append(".StoreUint32(&this_.ids_[offset], uint32(id))\n")
            .append("\t}\n")
            .append("\treturn id\n}\n");
      }
      out.append('\n');

      out.append(staticOut);
      out.append(methodsOut);
    } finally {
      if (hasEffTypeParams) {
        resolver.popTypeParamScope();
      }
    }

    // 10. Nested member classes/interfaces/enums/records
    emitNestedMemberTypes(ct, ctPath, emittedTypes, out, fileCtx);
  }

  private void emitBindSelfMethod(
      TypeElement te,
      String goName,
      String structName,
      boolean needsSelfField,
      TypeElement superTe,
      String superEmbeddedField,
      StringBuilder out,
      EmitContext fileCtx) {
    if (resolver.isValueStructClass(te)) {
      out.append("func (this_ ").append(structName).append(") BindSelf_(s any) {\n");
      out.append("\t_ = s\n}\n");
      return;
    }
    out.append("func (this_ *").append(structName).append(") BindSelf_(s any) {\n");
    out.append("\t_ = s\n");
    if (needsSelfField) {
      out.append("\tthis_.self_ = s.(").append(goName).append(")\n");
    }
    if (superTe != null && superEmbeddedField != null && !resolver.runtime.mapsToAny(superTe)) {
      if (resolver.isSourceType(superTe)) {
        if (resolver.needsBindSelf(superTe)) {
          out.append("\tthis_.").append(superEmbeddedField).append(".BindSelf_(s)\n");
        }
      } else if (!resolver.runtime.isShimException(superTe)
          && resolver.runtime.isSubclassedExternal(superTe)) {
        out.append("\tthis_.").append(superEmbeddedField).append(".BindExtSelf_(s)\n");
        if (superTe.getQualifiedName().contentEquals("com.google.gson.reflect.TypeToken")
            && te.getSuperclass() instanceof DeclaredType supDt
            && !supDt.getTypeArguments().isEmpty()) {
          out.append("\tthis_.")
              .append(superEmbeddedField)
              .append(".Core_ = ")
              .append(GoNames.goStringLiteral(supDt.getTypeArguments().get(0).toString()))
              .append('\n');
        }
      }
    }
    out.append("}\n");
  }

  private boolean shouldSkipStructFieldPacking(
      TypeElement te, List<VariableElement> instanceFields) {
    String qn = te.getQualifiedName().toString();
    if (qn.startsWith("com.google.javascript.jscomp.AbstractCommandLineRunner")
        || qn.startsWith("com.google.javascript.jscomp.CommandLineRunner")
        || qn.startsWith("com.google.javascript.jscomp.disambiguate.DisambiguateProperties.")
        || qn.startsWith("com.google.javascript.jscomp.serialization.SerializeTypesToPointers.")
        || qn.startsWith("com.google.javascript.jscomp.diagnostic.")
        || qn.startsWith("com.google.javascript.jscomp.instrumentation.")
        || qn.startsWith("com.google.javascript.jscomp.deps.")
        || te.getSimpleName().toString().endsWith("Json")) {
      return true;
    }
    for (VariableElement ve : instanceFields) {
      for (AnnotationMirror am : ve.getAnnotationMirrors()) {
        Element annEl = am.getAnnotationType().asElement();
        if (annEl instanceof TypeElement annTe) {
          String annQn = annTe.getQualifiedName().toString();
          if (annQn.equals("com.google.gson.annotations.SerializedName")
              || annQn.endsWith(".Keep")) {
            return true;
          }
        }
      }
    }
    return false;
  }

  private int classifyFieldPackingBucket(TypeMirror tm, Set<TypeElement> visited) {
    if (tm == null) {
      return 0;
    }
    TypeMirror er = resolver.goEffectiveType(tm);
    return switch (er.getKind()) {
      case BOOLEAN, BYTE -> 4;
      case SHORT, CHAR -> 3;
      case INT, FLOAT -> 2;
      case LONG, DOUBLE -> 1;
      case ARRAY, TYPEVAR -> 0;
      case DECLARED -> {
        TypeElement fte = (TypeElement) ((DeclaredType) er).asElement();
        if (resolver.isNodeClass(fte)) {
          yield 2;
        }
        if (resolver.isValueEnum(fte)) {
          yield 4;
        }
        if (resolver.isValueStructClass(fte)) {
          Set<TypeElement> seen = visited != null ? visited : new HashSet<>();
          if (!seen.add(fte)) {
            yield 0;
          }
          int minBucket = 4;
          TypeElement cur = fte;
          while (cur != null && resolver.isSourceType(cur)) {
            for (Element m : cur.getEnclosedElements()) {
              if (m instanceof VariableElement fve
                  && fve.getKind() == ElementKind.FIELD
                  && !fve.getModifiers().contains(Modifier.STATIC)) {
                int cb = classifyFieldPackingBucket(fve.asType(), seen);
                if (cb == 0) {
                  seen.remove(fte);
                  yield 0;
                }
                if (cb < minBucket) {
                  minBucket = cb;
                }
              }
            }
            cur = resolver.getSuperclassElement(cur);
          }
          seen.remove(fte);
          yield minBucket;
        }
        yield 0;
      }
      default -> 0;
    };
  }

  private String gsonStructTag(VariableElement ve) {
    if (ve.getModifiers().contains(Modifier.TRANSIENT)) {
      return "`gson:\"-\"`";
    }
    String jsonName = ve.getSimpleName().toString();
    for (AnnotationMirror am : ve.getAnnotationMirrors()) {
      Element annEl = am.getAnnotationType().asElement();
      if (annEl instanceof TypeElement annTe
          && annTe.getQualifiedName().contentEquals("com.google.gson.annotations.SerializedName")) {
        for (Map.Entry<? extends ExecutableElement, ? extends AnnotationValue> e :
            am.getElementValues().entrySet()) {
          if (e.getKey().getSimpleName().contentEquals("value") && e.getValue().getValue() != null) {
            jsonName = e.getValue().getValue().toString();
          }
        }
      }
    }
    return "`gson:\"" + jsonName.replace("\\", "\\\\").replace("\"", "\\\"") + "\"`";
  }

  private void emitJavaClassInfoMethod(TypeElement te, String structName, StringBuilder out) {
    String binName = resolver.elements.getBinaryName(te).toString();
    String simpleName = te.getSimpleName().toString();
    if (resolver.isValueStructClass(te)) {
      String sentinelField = resolver.getValueStructEnumSentinelFieldGoName(te);
      String nilBody =
          resolver.hasValueStructValidFlag(te)
              ? "!this_.valid_"
              : (sentinelField != null
                  ? ("this_." + sentinelField + " == 0")
                  : (resolver.isNullableValueStructClass(te)
                      ? ("this_ == (" + structName + "{})")
                      : "false"));
      out.append("func (this_ ")
          .append(structName)
          .append(") IsJavaNil_() bool { return ")
          .append(nilBody)
          .append(" }\n");
      out.append("func (")
          .append(structName)
          .append(") JavaClassInfo_() (string, string) { return ")
          .append(GoNames.goStringLiteral(binName))
          .append(", ")
          .append(GoNames.goStringLiteral(simpleName))
          .append(" }\n");
      return;
    }
    out.append("func (this_ *")
        .append(structName)
        .append(") IsJavaNil_() bool { return this_ == nil }\n");
    out.append("func (*")
        .append(structName)
        .append(") JavaClassInfo_() (string, string) { return ")
        .append(GoNames.goStringLiteral(binName))
        .append(", ")
        .append(GoNames.goStringLiteral(simpleName))
        .append(" }\n");
  }

  private void emitValueStructIdentityHashMethod(
      TypeElement te,
      String structWithParams,
      boolean hasOuterField,
      List<VariableElement> instanceFields,
      StringBuilder out,
      EmitContext fileCtx) {
    out.append("func (this_ ").append(structWithParams).append(") IdentityHash_() int32 {\n");
    String sentinelField = resolver.getValueStructEnumSentinelFieldGoName(te);
    if (resolver.hasValueStructValidFlag(te)) {
      out.append("\tif !this_.valid_ {\n\t\treturn 0\n\t}\n");
    } else if (sentinelField != null) {
      out.append("\tif this_.").append(sentinelField).append(" == 0 {\n\t\treturn 0\n\t}\n");
    } else if (resolver.isNullableValueStructClass(te)) {
      out.append("\tif this_ == (").append(structWithParams).append("{}) {\n\t\treturn 0\n\t}\n");
    }
    out.append("\tvar h_ int32 = 1\n");
    EmitContext ctx = fileCtx.copyForMember(te, false);
    if (hasOuterField) {
      TypeElement outerTe = resolver.getEnclosingInstanceType(te);
      if (outerTe != null) {
        out.append('\t');
        exprEmitter.appendTypedHashAccum(out, "h_", "this_.outer_", outerTe.asType(), false, ctx);
        out.append('\n');
      }
    }
    for (VariableElement f : instanceFields) {
      String fGo = resolver.getFieldGoName(f);
      out.append('\t');
      exprEmitter.appendTypedHashAccum(out, "h_", "this_." + fGo, f.asType(), false, ctx);
      out.append('\n');
    }
    out.append("\treturn h_\n}\n");
    fileCtx.syncTempSeqFrom(ctx);
  }

  /**
   * Returns true iff {@code te} is a concrete leaf Java class ({@code !isSubclassed}, so never a
   * {@code *Base} struct whose methods could be promoted onto overriding subclasses) whose
   * {@code equals(Object)} and {@code hashCode()} semantics are strictly inherited from
   * {@code java.lang.Object}.
   */
  private boolean hasIdentityEquals(TypeElement te, boolean isSubclassed) {
    if (isSubclassed
        || te.getKind() != ElementKind.CLASS
        || resolver.isNodeClass(te)
        || resolver.isValueStructClass(te)) {
      return false;
    }
    TypeElement cur = te;
    while (cur != null) {
      String cqn = cur.getQualifiedName().toString();
      if (cqn.equals("java.lang.Object")) {
        break;
      }
      if (!resolver.isSourceType(cur)) {
        return false;
      }
      for (Element m : cur.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)) {
          String mn = ee.getSimpleName().toString();
          if ((mn.equals("equals") && ee.getParameters().size() == 1)
              || (mn.equals("hashCode") && ee.getParameters().isEmpty())) {
            return false;
          }
          TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
          if (sig != null && (sig.goName.equals("Equals") || sig.goName.equals("HashCode"))) {
            return false;
          }
        }
      }
      TypeMirror supTm = cur.getSuperclass();
      if (!(supTm instanceof DeclaredType supDt) || !(supDt.asElement() instanceof TypeElement ste)) {
        return false;
      }
      cur = ste;
    }
    for (TypeElement anc : resolver.getAllAncestorTypes(te)) {
      if (anc.getKind() == ElementKind.INTERFACE) {
        for (Element m : anc.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && ee.getModifiers().contains(Modifier.DEFAULT)) {
            String mn = ee.getSimpleName().toString();
            if ((mn.equals("equals") && ee.getParameters().size() == 1)
                || (mn.equals("hashCode") && ee.getParameters().isEmpty())) {
              return false;
            }
          }
        }
      }
    }
    return true;
  }

  private TypeElement rootSourceClass(TypeElement te) {
    TypeElement cur = te;
    while (cur != null) {
      TypeElement sup = resolver.getSuperclassElement(cur);
      if (sup == null || resolver.runtime.mapsToAny(sup) || !resolver.isSourceType(sup)) {
        return cur;
      }
      cur = sup;
    }
    return te;
  }

  private boolean classHasIdentityHashField(TypeElement rootTe) {
    if (rootTe == null
        || rootTe.getKind() != ElementKind.CLASS
        || resolver.isNodeClass(rootTe)
        || resolver.isValueStructClass(rootTe)) {
      return false;
    }
    String qn = rootTe.getQualifiedName().toString();
    if (qn.equals("com.google.javascript.jscomp.parsing.parser.trees.ParseTree")) {
      return false;
    }
    TypeElement cur = rootTe;
    while (cur != null) {
      String cqn = cur.getQualifiedName().toString();
      if (cqn.equals("java.lang.Object")) {
        break;
      }
      for (Element m : cur.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)) {
          String mn = ee.getSimpleName().toString();
          if (mn.equals("hashCode") && ee.getParameters().isEmpty()) {
            return false;
          }
          TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
          if (sig != null && sig.goName.equals("HashCode")) {
            return false;
          }
        }
      }
      cur = resolver.getSuperclassElement(cur);
    }
    for (TypeElement anc : resolver.getAllAncestorTypes(rootTe)) {
      if (anc.getKind() == ElementKind.INTERFACE) {
        for (Element m : anc.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && ee.getModifiers().contains(Modifier.DEFAULT)) {
            String mn = ee.getSimpleName().toString();
            if (mn.equals("hashCode") && ee.getParameters().isEmpty()) {
              return false;
            }
          }
        }
      }
    }
    return true;
  }

  private void emitTypePredicateRegistration(
      TypeElement te, String goCheckType, StringBuilder out, EmitContext fileCtx) {
    if (te.getNestingKind() == NestingKind.ANONYMOUS || te.getNestingKind() == NestingKind.LOCAL) {
      return;
    }
    String qn = te.getQualifiedName().toString();
    String binName = resolver.elements.getBinaryName(te).toString();
    String rt = resolver.runtime.j2gortAlias(fileCtx.imports);
    String pfx = rt.isEmpty() ? "" : (rt + ".");
    out.append("func init() {\n\t")
        .append(pfx)
        .append("RegisterTypePredicate(")
        .append(GoNames.goStringLiteral(qn))
        .append(", ")
        .append(GoNames.goStringLiteral(binName))
        .append(", func(v_ any) bool {\n\t\treturn ")
        .append(pfx)
        .append("InstanceOf[")
        .append(goCheckType)
        .append("](v_)\n\t})\n}\n");
  }

  private void emitGenericTypePredicateRegistration(
      TypeElement te, String genMarker, StringBuilder out, EmitContext fileCtx) {
    if (te.getNestingKind() == NestingKind.ANONYMOUS || te.getNestingKind() == NestingKind.LOCAL) {
      return;
    }
    String qn = te.getQualifiedName().toString();
    String binName = resolver.elements.getBinaryName(te).toString();
    String rt = resolver.runtime.j2gortAlias(fileCtx.imports);
    String pfx = rt.isEmpty() ? "" : (rt + ".");
    out.append("func init() {\n\t")
        .append(pfx)
        .append("RegisterTypePredicate(")
        .append(GoNames.goStringLiteral(qn))
        .append(", ")
        .append(GoNames.goStringLiteral(binName))
        .append(", func(v_ any) bool {\n\t\tm_, ok_ := v_.(interface {\n\t\t\t")
        .append(genMarker)
        .append("()\n\t\t\tIsJavaNil_() bool\n\t\t})\n\t\treturn ok_ && !m_.IsJavaNil_()\n\t})\n}\n");
  }

  private void emitMarkerTypePredicateRegistration(
      TypeElement te, StringBuilder out, EmitContext fileCtx) {
    emitGenericTypePredicateRegistration(
        te, resolver.getInterfaceMarkerMethod(te), out, fileCtx);
  }

  private void emitCrossPkgFieldAccessors(
      TypeElement te,
      String structName,
      boolean isValueStruct,
      List<VariableElement> instanceFields,
      StringBuilder out,
      EmitContext fileCtx) {
    String getRecv = isValueStruct ? structName : ("*" + structName);
    boolean checkValid = isValueStruct && resolver.hasValueStructValidFlag(te);
    for (VariableElement ve : instanceFields) {
      String fGo = resolver.getFieldGoName(ve);
      String fType = resolver.goType(ve.asType(), fileCtx.imports);
      if (resolver.isFieldCrossPkgRead(ve)) {
        out.append("func (this_ ")
            .append(getRecv)
            .append(") FieldGet_")
            .append(fGo)
            .append("() ")
            .append(fType)
            .append(" { ");
        if (checkValid) {
          out.append("if !this_.valid_ { ").append(rtPrefix(fileCtx)).append("ThrowNPE() }; ");
        }
        out.append("return this_.").append(fGo).append(" }\n");
      }
      if (resolver.isFieldCrossPkgWrite(ve)) {
        out.append("func (this_ *")
            .append(structName)
            .append(") FieldSet_")
            .append(fGo)
            .append("(v ")
            .append(fType)
            .append(") { this_.")
            .append(fGo)
            .append(" = v }\n");
      }
    }
  }

  private void emitInstanceFieldsInit(
      ClassTree ct,
      TreePath ctPath,
      TypeElement te,
      String goName,
      String structName,
      Map<VariableElement, String> capturedFields,
      StringBuilder out,
      EmitContext fileCtx) {
    EmitContext ctx = fileCtx.copyForMember(te, false);
    ctx.capturedVarFields.putAll(capturedFields);

    StringBuilder body = new StringBuilder();
    for (Tree member : ct.getMembers()) {
      if (member instanceof VariableTree vt) {
        TreePath vtPath = new TreePath(ctPath, vt);
        Element el = resolver.trees.getElement(vtPath);
        if (el instanceof VariableElement ve
            && ve.getKind() == ElementKind.FIELD
            && !ve.getModifiers().contains(Modifier.STATIC)) {
          String fGo = resolver.getFieldGoName(ve);
          if (vt.getInitializer() != null) {
            String val = exprEmitter.emitAs(vt.getInitializer(), vtPath, ve.asType(), false, ctx);
            body.append("\tthis_.").append(fGo).append(" = ").append(val).append('\n');
          } else if (resolver.isStringType(ve.asType()) && te.getKind() != ElementKind.RECORD) {
            String rt = rtPrefix(ctx);
            body.append("\tthis_.").append(fGo).append(" = ").append(rt).append("NullString\n");
          }
        }
      } else if (member instanceof BlockTree bt && !bt.isStatic()) {
        body.append(stmtEmitter.emitBlockStatements(bt, new TreePath(ctPath, bt), 1, ctx));
      }
    }

    out.append("func (this_ *").append(structName).append(") initFields_").append(goName).append("() {\n");
    out.append("\t_ = this_\n");
    if (resolver.hasValueStructValidFlag(te)) {
      out.append("\tthis_.valid_ = true\n");
    }
    out.append(body);
    out.append("}\n\n");

    fileCtx.syncTempSeqFrom(ctx);
    fileCtx.deferredLocalTypes.addAll(ctx.deferredLocalTypes);
  }

  private void emitAnonymousClassConstructor(
      TreePath ctPath,
      TypeElement anonTe,
      String goName,
      String declParams,
      String structName,
      TypeElement superTe,
      String superEmbeddedField,
      Map<VariableElement, String> capturedFields,
      StringBuilder out,
      EmitContext fileCtx) {
    Tree parentLeaf = ctPath.getParentPath() != null ? ctPath.getParentPath().getLeaf() : null;
    ExecutableElement superCtor = null;
    if (parentLeaf instanceof NewClassTree) {
      Element el = resolver.trees.getElement(ctPath.getParentPath());
      if (el instanceof ExecutableElement ee && ee.getKind() == ElementKind.CONSTRUCTOR) {
        // Note: in JavacTask, trees.getElement(NewClassTree) for an anonymous class may return
        // the anonymous class's synthesized constructor; if so, find the matching super constructor!
        if (ee.getEnclosingElement().equals(anonTe)) {
          superCtor = findSuperCtorCalledFromSynthesizedAnonInit(ctPath, superTe, ee);
        } else {
          superCtor = ee;
        }
      }
    }

    List<String> sigParams = new ArrayList<>();
    StringBuilder body = new StringBuilder();

    if (resolver.hasOuterInstance(anonTe)) {
      TypeElement outerTe = resolver.getEnclosingInstanceType(anonTe);
      sigParams.add("outer_ " + resolver.goTypeForElement(outerTe, fileCtx.imports));
      body.append("\tthis_.outer_ = outer_\n");
    }

    boolean needsSuperOuter =
        superTe != null && resolver.isSourceType(superTe) && resolver.hasOuterInstance(superTe);
    if (needsSuperOuter) {
      TypeElement superOuterTe = resolver.getEnclosingInstanceType(superTe);
      sigParams.add("superOuter_ " + resolver.goTypeForElement(superOuterTe, fileCtx.imports));
    }

    for (VariableElement cap : resolver.getCapturedVars(anonTe)) {
      String pName = "cap_p_" + GoNames.cleanChars(cap.getSimpleName().toString());
      sigParams.add(pName + " " + resolver.goType(cap.asType(), fileCtx.imports));
      body.append("\tthis_.").append(capturedFields.get(cap)).append(" = ").append(pName).append('\n');
    }

    List<String> superArgs = new ArrayList<>();
    if (superCtor != null) {
      for (int i = 0; i < superCtor.getParameters().size(); i++) {
        VariableElement p = superCtor.getParameters().get(i);
        TypeMirror paramTm =
            superTe != null
                ? resolver.substituteOwnerTypeArgs(p.asType(), superTe, anonTe.getSuperclass())
                : p.asType();
        String pName = "sp" + i + "_";
        sigParams.add(pName + " " + resolver.goType(paramTm, fileCtx.imports));
        superArgs.add(pName);
      }
    } else if (parentLeaf instanceof NewClassTree nc && !nc.getArguments().isEmpty()) {
      for (int i = 0; i < nc.getArguments().size(); i++) {
        TreePath argPath = new TreePath(ctPath.getParentPath(), nc.getArguments().get(i));
        TypeMirror tm = resolver.trees.getTypeMirror(argPath);
        String pName = "sp" + i + "_";
        sigParams.add(pName + " " + resolver.goType(tm, fileCtx.imports));
        superArgs.add(pName);
      }
    }

    String retGo = resolver.goTypeForElement(anonTe, fileCtx.imports);
    out.append("func New")
        .append(goName)
        .append(declParams)
        .append('(')
        .append(String.join(", ", sigParams))
        .append(") ")
        .append(retGo)
        .append(" {\n");
    out.append("\tthis_ := &").append(structName).append("{}\n");
    if (resolver.needsBindSelf(anonTe)) {
      out.append("\tthis_.BindSelf_(this_)\n");
    }
    out.append(body);

    if (superTe != null && !resolver.runtime.mapsToAny(superTe)) {
      if (resolver.runtime.isShimException(superTe)) {
        List<String> exArgs = new ArrayList<>();
        exArgs.add("this_");
        exArgs.add(GoNames.goStringLiteral(superTe.getQualifiedName().toString() + "$Anon"));
        for (String sa : superArgs) {
          exArgs.add(rtPrefix(fileCtx) + "ToAny(" + sa + ")");
        }
        out.append('\t')
            .append(rtPrefix(fileCtx))
            .append("InitUserThrowable(")
            .append(String.join(", ", exArgs))
            .append(")\n");
      } else if (superCtor != null && superEmbeddedField != null) {
        if (needsSuperOuter) {
          String superPkg = resolver.planner.getGoPkg(superTe);
          if (superPkg.equals(fileCtx.unitPlan.goPkg)) {
            out.append("\tthis_.").append(superEmbeddedField).append(".outer_ = superOuter_\n");
          } else {
            out.append("\tthis_.").append(superEmbeddedField).append(".SetOuter_(superOuter_)\n");
          }
        }
        String superNew = resolver.getConstructorGoName(superCtor);
        String superInit =
            "Init"
                + superNew.substring(3)
                + resolver.formatSuperTypeArgs(anonTe.getSuperclass(), superTe, fileCtx.imports);
        List<String> initCallArgs = new ArrayList<>();
        initCallArgs.add("&this_." + superEmbeddedField);
        initCallArgs.addAll(superArgs);
        if (resolver.isSourceType(superTe)) {
          String superPkg = resolver.planner.getGoPkg(superTe);
          String alias = fileCtx.imports.ext(resolver.planner.getGoImportPath(superPkg), superPkg);
          String pref = alias.isEmpty() ? "" : (alias + ".");
          out.append('\t')
              .append(pref)
              .append(superInit)
              .append('(')
              .append(String.join(", ", initCallArgs))
              .append(")\n");
        } else {
          out.append('\t')
              .append(rtPrefix(fileCtx))
              .append(superInit)
              .append('(')
              .append(String.join(", ", initCallArgs))
              .append(")\n");
        }
      }
    }

    out.append("\tthis_.initFields_").append(goName).append("()\n");
    out.append("\treturn this_\n}\n\n");
  }

  private ExecutableElement findSuperCtorCalledFromSynthesizedAnonInit(
      TreePath anonCtPath, TypeElement superTe, ExecutableElement anonInit) {
    if (superTe == null) {
      return null;
    }
    ClassTree ct = (ClassTree) anonCtPath.getLeaf();
    for (Tree m : ct.getMembers()) {
      if (m instanceof MethodTree mt && isConstructorTree(mt) && mt.getBody() != null) {
        for (StatementTree st : mt.getBody().getStatements()) {
          if (st instanceof com.sun.source.tree.ExpressionStatementTree est
              && est.getExpression() instanceof com.sun.source.tree.MethodInvocationTree mit) {
            TreePath callPath =
                new TreePath(new TreePath(new TreePath(anonCtPath, mt), mt.getBody()), st);
            callPath = new TreePath(callPath, mit);
            Element t = resolver.trees.getElement(callPath);
            if (t instanceof ExecutableElement targetCtor
                && targetCtor.getKind() == ElementKind.CONSTRUCTOR
                && targetCtor.getEnclosingElement().equals(superTe)) {
              return targetCtor;
            }
          }
        }
      }
    }
    // Fallback: match by parameter count/erasure on superTe
    for (Element m : superTe.getEnclosedElements()) {
      if (m instanceof ExecutableElement sc
          && sc.getKind() == ElementKind.CONSTRUCTOR
          && sc.getParameters().size() == anonInit.getParameters().size()) {
        return sc;
      }
    }
    return null;
  }

  private void emitSynthesizedDefaultConstructor(
      TypeElement te,
      String goName,
      String declParams,
      String useParams,
      String structName,
      TypeElement superTe,
      String superEmbeddedField,
      boolean hasClinit,
      Map<VariableElement, String> capturedFields,
      StringBuilder out,
      EmitContext fileCtx) {
    String retGo = resolver.goTypeForElement(te, fileCtx.imports);
    List<String> extraParams = new ArrayList<>();
    StringBuilder allocAssigns = new StringBuilder();

    if (resolver.hasOuterInstance(te)) {
      TypeElement outerTe = resolver.getEnclosingInstanceType(te);
      extraParams.add("outer_ " + resolver.goTypeForElement(outerTe, fileCtx.imports));
      allocAssigns.append("\tthis_.outer_ = outer_\n");
    }
    for (VariableElement cap : resolver.getCapturedVars(te)) {
      String pName = "cap_p_" + GoNames.cleanChars(cap.getSimpleName().toString());
      extraParams.add(pName + " " + resolver.goType(cap.asType(), fileCtx.imports));
      allocAssigns.append("\tthis_.").append(capturedFields.get(cap)).append(" = ").append(pName).append('\n');
    }

    out.append("func New")
        .append(goName)
        .append(declParams)
        .append('(')
        .append(String.join(", ", extraParams))
        .append(") ")
        .append(retGo)
        .append(" {\n");
    if (hasClinit) {
      out.append("\tclinit_").append(goName).append("()\n");
    }
    if (resolver.isValueStructClass(te)) {
      out.append("\tthis_ := ").append(structName).append("{}\n");
      out.append(allocAssigns);
      out.append("\tInit").append(goName).append(useParams).append("(&this_)\n");
      out.append("\treturn this_\n}\n\n");
    } else {
      out.append("\tthis_ := &").append(structName).append("{}\n");
      if (resolver.needsBindSelf(te)) {
        out.append("\tthis_.BindSelf_(this_)\n");
      }
      out.append(allocAssigns);
      out.append("\tInit").append(goName).append(useParams).append("(this_)\n");
      out.append("\treturn this_\n}\n\n");
    }

    out.append("func Init")
        .append(goName)
        .append(declParams)
        .append("(this_ *")
        .append(structName)
        .append(") {\n");
    if (superTe != null && !resolver.runtime.mapsToAny(superTe)) {
      if (resolver.runtime.isShimException(superTe)) {
        out.append('\t')
            .append(rtPrefix(fileCtx))
            .append("InitUserThrowable(this_, ")
            .append(GoNames.goStringLiteral(te.getQualifiedName().toString()))
            .append(")\n");
      } else if (superEmbeddedField != null) {
        ExecutableElement noArgSuper = findNoArgConstructor(superTe);
        if (noArgSuper != null) {
          String superInit =
              "Init"
                  + resolver.getConstructorGoName(noArgSuper).substring(3)
                  + resolver.formatSuperTypeArgs(te.getSuperclass(), superTe, fileCtx.imports);
          if (resolver.isSourceType(superTe)) {
            String superPkg = resolver.planner.getGoPkg(superTe);
            String alias = fileCtx.imports.ext(resolver.planner.getGoImportPath(superPkg), superPkg);
            String pref = alias.isEmpty() ? "" : (alias + ".");
            out.append('\t').append(pref).append(superInit).append("(&this_.").append(superEmbeddedField).append(")\n");
          } else {
            out.append('\t')
                .append(rtPrefix(fileCtx))
                .append(superInit)
                .append("(&this_.")
                .append(superEmbeddedField)
                .append(")\n");
          }
        }
      }
    }
    out.append("\tthis_.initFields_").append(goName).append("()\n");
    out.append("}\n\n");
  }

  private ExecutableElement findNoArgConstructor(TypeElement te) {
    for (Element m : te.getEnclosedElements()) {
      if (m instanceof ExecutableElement ee
          && ee.getKind() == ElementKind.CONSTRUCTOR
          && ee.getParameters().isEmpty()) {
        return ee;
      }
    }
    return null;
  }

  private void emitSynthesizedRecordMethods(
      TypeElement te,
      String structName,
      Set<String> emittedMethodGoNames,
      StringBuilder out,
      EmitContext fileCtx) {
    List<? extends RecordComponentElement> comps = te.getRecordComponents();
    List<VariableElement> compFields = new ArrayList<>();
    for (RecordComponentElement rc : comps) {
      ExecutableElement accessor = rc.getAccessor();
      VariableElement fieldVe = findInstanceFieldByName(te, rc.getSimpleName().toString());
      if (fieldVe != null) {
        compFields.add(fieldVe);
      }
      if (accessor != null && fieldVe != null) {
        TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(accessor);
        if (sig != null && emittedMethodGoNames.add(sig.goName)) {
          TypeMirror canonRetTm = resolver.resolveSigReturnType(sig, te.asType());
          String retGo = resolver.goType(canonRetTm, fileCtx.imports);
          EmitContext ctx = fileCtx.copyForMember(te, false);
          String val =
              exprEmitter.convert(
                  "this_." + resolver.getFieldGoName(fieldVe),
                  fieldVe.asType(),
                  canonRetTm,
                  sig.returnWidenedToAny,
                  ctx);
          out.append("func (this_ *")
              .append(structName)
              .append(") ")
              .append(sig.goName)
              .append("() ")
              .append(retGo)
              .append(" { return ")
              .append(val)
              .append(" }\n");
        }
      }
    }

    String jl = resolver.runtime.javalangAlias(fileCtx.imports);
    String jlPref = jl.isEmpty() ? "" : (jl + ".");
    String rt = rtPrefix(fileCtx);

    if (emittedMethodGoNames.add("Equals")) {
      out.append("func (this_ *").append(structName).append(") Equals(other_ any) bool {\n");
      out.append("\to_, ok_ := other_.(*").append(structName).append(")\n");
      out.append("\tif !ok_ || o_ == nil {\n\t\treturn false\n\t}\n");
      out.append("\tif this_ == o_ {\n\t\treturn true\n\t}\n");
      EmitContext ctx = fileCtx.copyForMember(te, false);
      for (VariableElement f : compFields) {
        String fGo = resolver.getFieldGoName(f);
        TypeMirror er = resolver.safeErasure(f.asType());
        if (er.getKind().isPrimitive() && er.getKind() != TypeKind.FLOAT && er.getKind() != TypeKind.DOUBLE) {
          out.append("\tif this_.").append(fGo).append(" != o_.").append(fGo).append(" { return false }\n");
        } else if (resolver.isStringType(er)
            || resolver.isNodeType(er)
            || (resolver.isValueStructType(er)
                && er instanceof DeclaredType vsdt
                && vsdt.asElement() instanceof TypeElement vste
                && exprEmitter.valueStructHasIdentityEquals(vste))
            || (er instanceof DeclaredType fdt
                && fdt.asElement() instanceof TypeElement fte
                && (fte.getKind() == ElementKind.ENUM || exprEmitter.hasIdentityEquals(fte)))) {
          out.append("\tif this_.").append(fGo).append(" != o_.").append(fGo).append(" { return false }\n");
        } else {
          String a = exprEmitter.convert("this_." + fGo, er, resolver.getObjectElement().asType(), true, ctx);
          String b = exprEmitter.convert("o_." + fGo, er, resolver.getObjectElement().asType(), true, ctx);
          out.append("\tif !").append(jlPref).append("ObjectsEquals(").append(a).append(", ").append(b).append(") { return false }\n");
        }
      }
      out.append("\treturn true\n}\n");
    }

    if (emittedMethodGoNames.add("HashCode")) {
      out.append("func (this_ *").append(structName).append(") HashCode() int32 {\n");
      out.append("\tvar h_ int32 = 0\n");
      EmitContext ctx = fileCtx.copyForMember(te, false);
      for (VariableElement f : compFields) {
        String fGo = resolver.getFieldGoName(f);
        out.append('\t');
        exprEmitter.appendTypedHashAccum(out, "h_", "this_." + fGo, f.asType(), false, ctx);
        out.append('\n');
      }
      out.append("\treturn h_\n}\n");
    }

    if (emittedMethodGoNames.add("String")) {
      String fmtPkg = fileCtx.imports.std("fmt");
      out.append("func (this_ *").append(structName).append(") String() string {\n");
      out.append("\t_ = ").append(fmtPkg).append(".Sprint\n");
      List<String> parts = new ArrayList<>();
      EmitContext ctx = fileCtx.copyForMember(te, false);
      for (VariableElement f : compFields) {
        String fGo = resolver.getFieldGoName(f);
        String boxed =
            exprEmitter.convert(
                "this_." + fGo, f.asType(), resolver.getObjectElement().asType(), true, ctx);
        parts.add(
            GoNames.goStringLiteral(f.getSimpleName() + "=")
                + " + "
                + rt
                + "StrVal("
                + boxed
                + ")");
      }
      if (parts.isEmpty()) {
        out.append("\treturn ").append(GoNames.goStringLiteral(te.getSimpleName() + "[]")).append('\n');
      } else {
        out.append("\treturn ")
            .append(GoNames.goStringLiteral(te.getSimpleName() + "["))
            .append(" + ")
            .append(String.join(" + \", \" + ", parts))
            .append(" + \"]\"\n");
      }
      out.append("}\n\n");
    }
  }

  private void emitInheritedInterfaceDefaultForwarders(
      TypeElement te,
      String structName,
      boolean hasSelfField,
      Set<String> alreadyEmittedOnThisClass,
      StringBuilder out,
      EmitContext fileCtx) {
    Set<String> covered = new HashSet<>(alreadyEmittedOnThisClass);

    // Collect all non-abstract methods already implemented on any embedded superclass of `te`
    TypeElement curSup = resolver.getSuperclassElement(te);
    while (curSup != null && !resolver.runtime.mapsToAny(curSup)) {
      if (resolver.runtime.isShimException(curSup)) {
        covered.add("GetMessage");
        covered.add("GetCause");
        covered.add("Error");
      } else if (resolver.isSourceType(curSup)) {
        for (Element m : curSup.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && !ee.getModifiers().contains(Modifier.STATIC)
              && !ee.getModifiers().contains(Modifier.PRIVATE)
              && !ee.getModifiers().contains(Modifier.ABSTRACT)) {
            TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
            if (sig != null) {
              covered.add(sig.goName);
            }
          }
        }
        for (TypeElement supAnc : resolver.getAllAncestorTypes(curSup)) {
          if (supAnc.getKind() == ElementKind.INTERFACE && !resolver.runtime.mapsToAny(supAnc)) {
            for (Element m : supAnc.getEnclosedElements()) {
              if (m instanceof ExecutableElement ee
                  && ee.getKind() == ElementKind.METHOD
                  && !ee.getModifiers().contains(Modifier.STATIC)
                  && !ee.getModifiers().contains(Modifier.PRIVATE)) {
                TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
                if (sig != null
                    && (ee.getModifiers().contains(Modifier.DEFAULT)
                        || isObjectMethodOnInterface(sig))
                    && (resolver.isSourceType(supAnc) || resolver.isNeededExternalSig(sig))) {
                  covered.add(sig.goName);
                }
              }
            }
          }
        }
      } else {
        List<TypeElement> extChain = new ArrayList<>();
        extChain.add(curSup);
        extChain.addAll(resolver.getAllAncestorTypes(curSup));
        for (TypeElement extAnc : extChain) {
          if (resolver.runtime.mapsToAny(extAnc) || resolver.runtime.isShimException(extAnc)) {
            continue;
          }
          for (Element m : extAnc.getEnclosedElements()) {
            if (m instanceof ExecutableElement ee
                && ee.getKind() == ElementKind.METHOD
                && !ee.getModifiers().contains(Modifier.STATIC)
                && !ee.getModifiers().contains(Modifier.PRIVATE)) {
              TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
              if (sig != null && resolver.isNeededExternalSig(sig)) {
                covered.add(sig.goName);
              }
            }
          }
        }
      }
      curSup = resolver.getSuperclassElement(curSup);
    }

    // Check all ancestor interfaces in topological/BFS order (most specific interfaces first!)
    // Also synthesize java.lang.Object methods (`equals`, `hashCode`, `toString`) redeclared on interfaces.
    for (TypeElement anc : resolver.getAllAncestorTypes(te)) {
      if (anc.getKind() == ElementKind.INTERFACE && !resolver.runtime.mapsToAny(anc)) {
        for (Element m : anc.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && !ee.getModifiers().contains(Modifier.STATIC)
              && !ee.getModifiers().contains(Modifier.PRIVATE)) {
            TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
            if (sig != null
                && (ee.getModifiers().contains(Modifier.DEFAULT) || isObjectMethodOnInterface(sig))
                && (resolver.isSourceType(anc) || resolver.isNeededExternalSig(sig))
                && covered.add(sig.goName)) {
              emitDefaultOrFallbackMethodOnStruct(
                  te, structName, hasSelfField, anc, ee, sig, out, fileCtx);
            }
          }
        }
      }
    }
  }

  private boolean isObjectMethodOnInterface(TypeResolver.MethodSig sig) {
    return (sig.simpleName.equals("equals")
            && sig.canonParamTypes.size() == 1
            && sig.goName.equals("Equals"))
        || (sig.simpleName.equals("hashCode")
            && sig.canonParamTypes.isEmpty()
            && sig.goName.equals("HashCode"))
        || (sig.simpleName.equals("toString")
            && sig.canonParamTypes.isEmpty()
            && sig.goName.equals("String"));
  }

  private void emitDefaultOrFallbackMethodOnStruct(
      TypeElement hostTe,
      String structName,
      boolean hasSelfField,
      TypeElement iface,
      ExecutableElement ee,
      TypeResolver.MethodSig sig,
      StringBuilder out,
      EmitContext fileCtx) {
    List<String> sigParams = new ArrayList<>();
    List<String> callArgs = new ArrayList<>();
    TypeMirror hostSiteTm = hostTe != null ? hostTe.asType() : null;
    // Note: if the class is subclassed, pass dynamic `self_` (reusing an ancestor's `self_` when
    // compatible so intermediate `*Base` structs don't emit a redundant `self_` field).
    String selfExpr = "this_";
    if (hasSelfField) {
      TypeElement samePkgSelfRoot = resolver.getSamePkgReusableSelfRoot(hostTe);
      selfExpr =
          samePkgSelfRoot != null
              ? ("this_.self_" + resolver.planner.getGoTypeName(hostTe) + "_()")
              : "this_.self_";
      if (!ee.getModifiers().contains(Modifier.DEFAULT) && isObjectMethodOnInterface(sig)) {
        TypeElement anc =
            resolver.findReusableSelfAncestor(hostTe, resolver.getObjectElement().asType(), true);
        if (anc != null) {
          selfExpr = "this_." + resolver.planner.getGoTypeName(anc) + "Base.self_";
        }
      } else if (ee.getModifiers().contains(Modifier.DEFAULT)
          && resolver.isSourceType(iface)
          && hostSiteTm instanceof DeclaredType hostDt) {
        DeclaredType ifaceSite = resolver.findSuperType(hostDt, iface);
        if (ifaceSite != null) {
          TypeElement anc = resolver.findReusableSelfAncestor(hostTe, ifaceSite, false);
          if (anc != null) {
            selfExpr = "this_." + resolver.planner.getGoTypeName(anc) + "Base.self_";
          }
        }
      }
    }
    callArgs.add(selfExpr);

    for (int i = 0; i < sig.canonParamTypes.size(); i++) {
      String pn = "p" + i + "_";
      TypeMirror pTm = resolver.resolveSigParamType(sig, i, hostSiteTm);
      sigParams.add(pn + " " + resolver.goType(pTm, fileCtx.imports));
      callArgs.add(pn);
    }
    TypeMirror retTm = resolver.resolveSigReturnType(sig, hostSiteTm);
    String retGo = retTm != null ? resolver.goType(retTm, fileCtx.imports) : "";

    out.append("func (this_ *")
        .append(structName)
        .append(") ")
        .append(sig.goName)
        .append('(')
        .append(String.join(", ", sigParams))
        .append(')')
        .append(retGo.isEmpty() ? "" : (" " + retGo))
        .append(" {\n\t");

    if (resolver.isSourceType(iface) && ee.getModifiers().contains(Modifier.DEFAULT)) {
      String ifacePkg = resolver.planner.getGoPkg(iface);
      String alias = fileCtx.imports.ext(resolver.planner.getGoImportPath(ifacePkg), ifacePkg);
      String pref = alias.isEmpty() ? "" : (alias + ".");
      String ifaceTypeArgs =
          hostSiteTm != null ? resolver.formatSuperTypeArgs(hostSiteTm, iface, fileCtx.imports) : "";
      String helper =
          pref
              + "IfaceDefault_"
              + resolver.planner.getGoTypeName(iface)
              + "_"
              + sig.goName
              + ifaceTypeArgs;
      if (!retGo.isEmpty()) {
        out.append("return ");
      }
      out.append(helper).append('(').append(String.join(", ", callArgs)).append(")\n}\n\n");
    } else if (!ee.getModifiers().contains(Modifier.DEFAULT) && isObjectMethodOnInterface(sig)) {
      String rt = rtPrefix(fileCtx);
      String selfArg = callArgs.get(0);
      if (sig.goName.equals("Equals")) {
        out.append("return ").append(rt).append("RefEquals(").append(selfArg).append(", p0_)\n}\n\n");
      } else if (sig.goName.equals("HashCode")) {
        out.append("return ").append(rt).append("IdentityHash(").append(selfArg).append(")\n}\n\n");
      } else {
        String jl = resolver.runtime.javalangAlias(fileCtx.imports);
        String jlPref = jl.isEmpty() ? "" : (jl + ".");
        out.append("return ")
            .append(jlPref)
            .append("ObjectToString(")
            .append(GoNames.goStringLiteral(structName))
            .append(", ")
            .append(selfArg)
            .append(")\n}\n\n");
      }
    } else {
      // External interface default method
      String rt = rtPrefix(fileCtx);
      String selfArg = callArgs.get(0);
      String ifaceQn = iface.getQualifiedName().toString();
      if (ifaceQn.equals("com.google.common.base.Predicate")
          && sig.simpleName.equals("test")
          && sig.canonParamTypes.size() == 1) {
        out.append("return ").append(selfArg).append(".Apply(p0_)\n}\n\n");
      } else if (ifaceQn.equals("java.lang.Iterable")
          && sig.simpleName.equals("forEach")
          && sig.canonParamTypes.size() == 1) {
        out.append(rt).append("ExtDefault_Iterable_ForEach(").append(selfArg).append(", p0_)\n}\n\n");
      } else if (ifaceQn.equals("java.util.Iterator")
          && sig.simpleName.equals("remove")
          && sig.canonParamTypes.isEmpty()) {
        String jl = resolver.runtime.javalangAlias(fileCtx.imports);
        String jlPref = jl.isEmpty() ? "" : (jl + ".");
        out.append("panic(").append(jlPref).append("NewUnsupportedOperationException(\"remove\"))\n}\n\n");
      } else if (ifaceQn.equals("java.util.Iterator")
          && sig.simpleName.equals("forEachRemaining")
          && sig.canonParamTypes.size() == 1) {
        out.append(rt).append("ExtDefault_Iterator_ForEachRemaining(").append(selfArg).append(", p0_)\n}\n\n");
      } else if (ifaceQn.equals("java.util.Comparator")) {
        if (sig.simpleName.equals("reversed") && sig.canonParamTypes.isEmpty()) {
          out.append("return ").append(rt).append("ExtDefault_Comparator_Reversed(").append(selfArg).append(")\n}\n\n");
        } else if (sig.simpleName.equals("thenComparing") && sig.canonParamTypes.size() == 1) {
          String p0Type = resolver.goType(sig.canonParamTypes.get(0), fileCtx.imports);
          if (p0Type.contains("Comparator")) {
            out.append("return ").append(rt).append("ExtDefault_Comparator_ThenComparing(").append(selfArg).append(", p0_)\n}\n\n");
          } else {
            out.append("return ").append(rt).append("ExtDefault_Comparator_ThenComparingKey(").append(selfArg).append(", p0_)\n}\n\n");
          }
        } else if (sig.simpleName.equals("thenComparing") && sig.canonParamTypes.size() == 2) {
          out.append("return ").append(rt).append("ExtDefault_Comparator_ThenComparingKeyCmp(").append(selfArg).append(", p0_, p1_)\n}\n\n");
        } else if (sig.simpleName.equals("thenComparingInt") && sig.canonParamTypes.size() == 1) {
          out.append("return ").append(rt).append("ExtDefault_Comparator_ThenComparingInt(").append(selfArg).append(", p0_)\n}\n\n");
        } else if (sig.simpleName.equals("thenComparingLong") && sig.canonParamTypes.size() == 1) {
          out.append("return ").append(rt).append("ExtDefault_Comparator_ThenComparingLong(").append(selfArg).append(", p0_)\n}\n\n");
        } else if (sig.simpleName.equals("thenComparingDouble") && sig.canonParamTypes.size() == 1) {
          out.append("return ").append(rt).append("ExtDefault_Comparator_ThenComparingDouble(").append(selfArg).append(", p0_)\n}\n\n");
        } else {
          for (int i = 1; i < callArgs.size(); i++) {
            out.append("_ = ").append(callArgs.get(i)).append("; ");
          }
          out.append("panic(").append(GoNames.goStringLiteral("UNIMPLEMENTED_EXT_DEFAULT: " + ifaceQn + "." + sig.goName)).append(")\n}\n\n");
        }
      } else if (ifaceQn.equals("java.util.function.Predicate") || ifaceQn.equals("com.google.common.base.Predicate")) {
        if (sig.simpleName.equals("negate") && sig.canonParamTypes.isEmpty()) {
          out.append("return ").append(rt).append("ExtDefault_Predicate_Negate(").append(selfArg).append(")\n}\n\n");
        } else if (sig.simpleName.equals("and") && sig.canonParamTypes.size() == 1) {
          out.append("return ").append(rt).append("ExtDefault_Predicate_And(").append(selfArg).append(", p0_)\n}\n\n");
        } else if (sig.simpleName.equals("or") && sig.canonParamTypes.size() == 1) {
          out.append("return ").append(rt).append("ExtDefault_Predicate_Or(").append(selfArg).append(", p0_)\n}\n\n");
        } else {
          for (int i = 1; i < callArgs.size(); i++) {
            out.append("_ = ").append(callArgs.get(i)).append("; ");
          }
          out.append("panic(").append(GoNames.goStringLiteral("UNIMPLEMENTED_EXT_DEFAULT: " + ifaceQn + "." + sig.goName)).append(")\n}\n\n");
        }
      } else if (ifaceQn.equals("java.util.function.Function") || ifaceQn.equals("com.google.common.base.Function")) {
        if (sig.simpleName.equals("andThen") && sig.canonParamTypes.size() == 1) {
          out.append("return ").append(rt).append("ExtDefault_Function_AndThen(").append(selfArg).append(", p0_)\n}\n\n");
        } else if (sig.simpleName.equals("compose") && sig.canonParamTypes.size() == 1) {
          out.append("return ").append(rt).append("ExtDefault_Function_Compose(").append(selfArg).append(", p0_)\n}\n\n");
        } else {
          for (int i = 1; i < callArgs.size(); i++) {
            out.append("_ = ").append(callArgs.get(i)).append("; ");
          }
          out.append("panic(").append(GoNames.goStringLiteral("UNIMPLEMENTED_EXT_DEFAULT: " + ifaceQn + "." + sig.goName)).append(")\n}\n\n");
        }
      } else if (ifaceQn.equals("java.util.function.Consumer") && sig.simpleName.equals("andThen") && sig.canonParamTypes.size() == 1) {
        out.append("return ").append(rt).append("ExtDefault_Consumer_AndThen(").append(selfArg).append(", p0_)\n}\n\n");
      } else {
        for (int i = 1; i < callArgs.size(); i++) {
          out.append("_ = ").append(callArgs.get(i)).append("; ");
        }
        out.append("panic(").append(GoNames.goStringLiteral("UNIMPLEMENTED_EXT_DEFAULT: " + ifaceQn + "." + sig.goName)).append(")\n}\n\n");
      }
    }
  }

  // =========================================================================
  // 3. Enums (Fields, Constructors, Abstract Methods, Constant-Specific Overrides)
  // =========================================================================

  private void emitEnum(
      ClassTree ct,
      TreePath ctPath,
      TypeElement te,
      Set<TypeElement> emittedTypes,
      StringBuilder out,
      EmitContext fileCtx) {
    if (resolver.isValueEnum(te)) {
      emitValueEnum(ct, ctPath, te, emittedTypes, out, fileCtx);
      return;
    }
    String goName = resolver.planner.getGoTypeName(te);
    String rt = rtPrefix(fileCtx);

    // Discover enum constants and any constant-specific overridden methods (`A { @Override ... }`)
    List<VariableTree> enumConstTrees = new ArrayList<>();
    Set<String> overriddenMethodGoNames = new LinkedHashSet<>();
    Map<String, TypeResolver.MethodSig> overriddenMethodSigs = new LinkedHashMap<>();

    for (Tree member : ct.getMembers()) {
      if (member instanceof VariableTree vt) {
        TreePath vtPath = new TreePath(ctPath, vt);
        Element el = resolver.trees.getElement(vtPath);
        if (el instanceof VariableElement ve && ve.getKind() == ElementKind.ENUM_CONSTANT) {
          enumConstTrees.add(vt);
          if (vt.getInitializer() instanceof NewClassTree nct && nct.getClassBody() != null) {
            TreePath bodyPath = new TreePath(new TreePath(vtPath, nct), nct.getClassBody());
            for (Tree bm : nct.getClassBody().getMembers()) {
              if (bm instanceof MethodTree bmt && !isConstructorTree(bmt)) {
                Element bel = resolver.trees.getElement(new TreePath(bodyPath, bmt));
                if (bel instanceof ExecutableElement bee
                    && !bee.getModifiers().contains(Modifier.STATIC)) {
                  TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(bee);
                  if (sig != null) {
                    overriddenMethodGoNames.add(sig.goName);
                    overriddenMethodSigs.put(sig.goName, sig);
                  }
                }
              }
            }
          }
        }
      }
    }

    // Also include any `abstract` methods declared on the enum itself
    for (Element m : te.getEnclosedElements()) {
      if (m instanceof ExecutableElement ee
          && ee.getKind() == ElementKind.METHOD
          && ee.getModifiers().contains(Modifier.ABSTRACT)) {
        TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
        if (sig != null) {
          overriddenMethodGoNames.add(sig.goName);
          overriddenMethodSigs.put(sig.goName, sig);
        }
      }
    }

    // 1. Emit `type <GoName> struct`
    out.append("// Java enum: ").append(te.getQualifiedName()).append('\n');
    out.append("type ").append(goName).append(" struct {\n");
    out.append('\t').append(rt).append("EnumBase\n");

    List<VariableElement> instanceFields = new ArrayList<>();
    for (Element m : te.getEnclosedElements()) {
      if (m instanceof VariableElement ve
          && ve.getKind() == ElementKind.FIELD
          && !ve.getModifiers().contains(Modifier.STATIC)) {
        instanceFields.add(ve);
        out.append('\t')
            .append(resolver.getFieldGoName(ve))
            .append(' ')
            .append(resolver.goType(ve.asType(), fileCtx.imports))
            .append('\n');
      }
    }

    NodePredInfo enumNodePred = inspectNodePredicateImpl(te);
    for (Map.Entry<String, TypeResolver.MethodSig> e : overriddenMethodSigs.entrySet()) {
      TypeResolver.MethodSig sig = e.getValue();
      List<String> fnParams = new ArrayList<>();
      fnParams.add("*" + goName);
      for (int i = 0; i < sig.canonParamTypes.size(); i++) {
        fnParams.add(
            resolver.goType(resolver.resolveSigParamType(sig, i, te.asType()), fileCtx.imports));
      }
      TypeMirror enumRetTm = resolver.resolveSigReturnType(sig, te.asType());
      String retGo = enumRetTm != null ? resolver.goType(enumRetTm, fileCtx.imports) : "";
      out.append("\tfn_")
          .append(e.getKey())
          .append(" func(")
          .append(String.join(", ", fnParams))
          .append(')')
          .append(retGo.isEmpty() ? "" : (" " + retGo))
          .append('\n');
    }
    if (enumNodePred != null && overriddenMethodGoNames.contains(enumNodePred.primarySamGoName())) {
      out.append("\tfn_TestNode_ func(*").append(goName).append(", uint32) bool\n");
    }
    out.append("}\n\n");

    out.append("func (this_ *")
        .append(goName)
        .append(") GetEnumBase() *")
        .append(rt)
        .append("EnumBase {\n\tif this_ == nil {\n\t\treturn nil\n\t}\n\treturn &this_.EnumBase\n}\n");
    out.append("func (this_ *")
        .append(goName)
        .append(") EnumOrdinal_() int32 {\n\tif this_ == nil {\n\t\treturn -1\n\t}\n\treturn this_.EnumBase.Ordinal_\n}\n");
    emitMarkerMethodsForType(te, goName, false, out);
    emitCrossPkgFieldAccessors(te, goName, false, instanceFields, out, fileCtx);

    List<String> staticConstLiterals =
        tryBuildStaticEnumConstantLiterals(
            ct, ctPath, te, enumConstTrees, instanceFields, overriddenMethodSigs, fileCtx);
    boolean hasStaticInitBlock = hasNonEmptyStaticInitBlock(ct);
    boolean hasEnumClinit;

    if (staticConstLiterals != null) {
      hasEnumClinit =
          emitStaticEnumConstantsAndClinit(
              ct,
              ctPath,
              te,
              goName,
              enumConstTrees,
              staticConstLiterals,
              hasStaticInitBlock,
              out,
              fileCtx);
    } else {
      // 2. Emit `initFields_<GoName>`
      emitInstanceFieldsInit(ct, ctPath, te, goName, goName, Map.of(), out, fileCtx);

      // 3. Emit enum constructors (`initEnum_*`)
      boolean hasEnumCtor = false;
      for (Tree member : ct.getMembers()) {
        if (member instanceof MethodTree mt && isConstructorTree(mt)) {
          hasEnumCtor = true;
          methodEmitter.emitConstructor(
              mt, new TreePath(ctPath, mt), te, true, Map.of(), out, fileCtx);
        }
      }
      if (!hasEnumCtor) {
        out.append("func initEnum_")
            .append(goName)
            .append("(this_ *")
            .append(goName)
            .append(", enumName_ string, enumOrd_ int32) *")
            .append(goName)
            .append(
                " {\n\tthis_.EnumBase.Name_ = enumName_\n\tthis_.EnumBase.Ordinal_ = enumOrd_\n\tthis_.initFields_")
            .append(goName)
            .append("()\n\treturn this_\n}\n\n");
      }

      // 4. Emit enum constants + static fields + unified `clinit_<GoName>()`
      emitEnumClinitAndStaticMembers(ct, ctPath, te, goName, enumConstTrees, out, fileCtx);
      out.append("var enumValues_").append(goName).append(" []*").append(goName).append('\n');
      hasEnumClinit = true;
    }

    // 5. Emit `<GoName>_Values()` and `<GoName>_ValueOf(name string)`
    boolean valuesNeedsClinit = staticConstLiterals == null || hasStaticInitBlock;
    out.append("func ").append(goName).append("_Values() []*").append(goName).append(" {\n");
    if (valuesNeedsClinit) {
      out.append("\tclinit_").append(goName).append("()\n");
    }
    out.append("\treturn ").append(rt).append("CloneSlice(enumValues_").append(goName).append(")\n}\n\n");

    out.append("func ").append(goName).append("_ValueOf(name string) *").append(goName).append(" {\n");
    if (valuesNeedsClinit) {
      out.append("\tclinit_").append(goName).append("()\n");
    }
    out.append("\tfor _, c := range enumValues_").append(goName).append(" {\n");
    out.append("\t\tif c.Name() == name {\n\t\t\treturn c\n\t\t}\n\t}\n");
    String jl = resolver.runtime.javalangAlias(fileCtx.imports);
    String jlPref = jl.isEmpty() ? "" : (jl + ".");
    out.append("\tpanic(")
        .append(jlPref)
        .append("NewIllegalArgumentException(\"No enum constant ")
        .append(te.getQualifiedName())
        .append(".\" + name))\n}\n\n");

    // 6. Emit enum methods + inherited interface default forwarders + nested member types
    Set<String> emittedMethodGoNames = new HashSet<>();
    emittedMethodGoNames.add("Name");
    emittedMethodGoNames.add("Ordinal");
    emittedMethodGoNames.add("HashCode");
    boolean emittedEnumTestNode = false;
    for (Tree member : ct.getMembers()) {
      if (member instanceof MethodTree mt && !isConstructorTree(mt)) {
        TreePath mtPath = new TreePath(ctPath, mt);
        Element el = resolver.trees.getElement(mtPath);
        if (el instanceof ExecutableElement ee && !ee.getModifiers().contains(Modifier.STATIC)) {
          TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
          if (sig != null) {
            emittedMethodGoNames.add(sig.goName);
          }
          if (enumNodePred != null
              && !emittedEnumTestNode
              && !ee.getModifiers().contains(Modifier.PRIVATE)
              && mt.getParameters().size() == 1
              && sig != null
              && sig.goName.equals(enumNodePred.primarySamGoName())) {
            emitNodePredicateMethodAndBridge(
                mt,
                mtPath,
                te,
                ee,
                sig,
                enumNodePred,
                overriddenMethodGoNames.contains(sig.goName),
                Map.of(),
                out,
                fileCtx);
            emittedEnumTestNode = true;
            continue;
          }
        }
        methodEmitter.emitMethod(
            mt, mtPath, te, hasEnumClinit, overriddenMethodGoNames, Map.of(), out, fileCtx);
      }
    }
    // Emit dispatchers for inherited interface/Object methods overridden only inside enum constant bodies
    for (TypeResolver.MethodSig sig : overriddenMethodSigs.values()) {
      if (emittedMethodGoNames.add(sig.goName)) {
        if (enumNodePred != null
            && !emittedEnumTestNode
            && sig.goName.equals(enumNodePred.primarySamGoName())) {
          emitEnumInheritedNodePredicateDispatcher(te, goName, sig, out, fileCtx);
          emittedEnumTestNode = true;
        } else {
          emitEnumInheritedOverrideDispatcher(te, goName, sig, out, fileCtx);
        }
      }
    }
    emittedMethodGoNames.add("String");
    emittedMethodGoNames.add("Equals");
    emittedMethodGoNames.add("CompareTo");
    emitInheritedInterfaceDefaultForwarders(te, goName, false, emittedMethodGoNames, out, fileCtx);
    emitJavaClassInfoMethod(te, goName, out);
    emitTypePredicateRegistration(te, "*" + goName, out, fileCtx);
    emitNestedMemberTypes(ct, ctPath, emittedTypes, out, fileCtx);
  }

  private void emitEnumInheritedOverrideDispatcher(
      TypeElement enumTe,
      String goName,
      TypeResolver.MethodSig sig,
      StringBuilder out,
      EmitContext fileCtx) {
    List<String> sigParams = new ArrayList<>();
    List<String> fnArgs = new ArrayList<>();
    List<String> baseArgs = new ArrayList<>();
    fnArgs.add("this_");
    for (int i = 0; i < sig.canonParamTypes.size(); i++) {
      String pn = "p" + i + "_";
      sigParams.add(
          pn
              + " "
              + resolver.goType(
                  resolver.resolveSigParamType(sig, i, enumTe.asType()), fileCtx.imports));
      fnArgs.add(pn);
      baseArgs.add(pn);
    }
    TypeMirror enumDispRetTm = resolver.resolveSigReturnType(sig, enumTe.asType());
    String retGo = enumDispRetTm != null ? resolver.goType(enumDispRetTm, fileCtx.imports) : "";
    String retPrefix = retGo.isEmpty() ? "" : "return ";

    out.append("func (this_ *")
        .append(goName)
        .append(") ")
        .append(sig.goName)
        .append('(')
        .append(String.join(", ", sigParams))
        .append(')')
        .append(retGo.isEmpty() ? "" : (" " + retGo))
        .append(" {\n\tif this_.fn_")
        .append(sig.goName)
        .append(" != nil {\n\t\t")
        .append(retPrefix)
        .append("this_.fn_")
        .append(sig.goName)
        .append('(')
        .append(String.join(", ", fnArgs))
        .append(")\n\t}\n\t")
        .append(retPrefix)
        .append("this_.baseEnum_")
        .append(sig.goName)
        .append('(')
        .append(String.join(", ", baseArgs))
        .append(")\n}\n\n");

    out.append("func (this_ *")
        .append(goName)
        .append(") baseEnum_")
        .append(sig.goName)
        .append('(')
        .append(String.join(", ", sigParams))
        .append(')')
        .append(retGo.isEmpty() ? "" : (" " + retGo))
        .append(" {\n\t");

    if (sig.goName.equals("String") && sig.canonParamTypes.isEmpty()) {
      out.append("return this_.EnumBase.String()\n}\n\n");
      return;
    }
    if (sig.goName.equals("Equals") && sig.canonParamTypes.size() == 1) {
      out.append("return any(this_) == p0_\n}\n\n");
      return;
    }

    // Check if an ancestor source interface provides a `default` implementation for `sig`
    for (TypeElement anc : resolver.getAllAncestorTypes(enumTe)) {
      if (anc.getKind() == ElementKind.INTERFACE && resolver.isSourceType(anc)) {
        for (Element m : anc.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && ee.getModifiers().contains(Modifier.DEFAULT)
              && resolver.getInstanceMethodSig(ee) == sig) {
            String ifacePkg = resolver.planner.getGoPkg(anc);
            String alias = fileCtx.imports.ext(resolver.planner.getGoImportPath(ifacePkg), ifacePkg);
            String pref = alias.isEmpty() ? "" : (alias + ".");
            String helper =
                pref
                    + "IfaceDefault_"
                    + resolver.planner.getGoTypeName(anc)
                    + "_"
                    + sig.goName
                    + resolver.formatSuperTypeArgs(enumTe.asType(), anc, fileCtx.imports);
            out.append(retPrefix)
                .append(helper)
                .append('(')
                .append(String.join(", ", fnArgs))
                .append(")\n}\n\n");
            return;
          }
        }
      }
    }

    for (String ba : baseArgs) {
      out.append("_ = ").append(ba).append("; ");
    }
    out.append("panic(\"abstract or native method invoked\")\n}\n\n");
  }

  private void emitEnumClinitAndStaticMembers(
      ClassTree ct,
      TreePath ctPath,
      TypeElement enumTe,
      String goName,
      List<VariableTree> enumConstTrees,
      StringBuilder out,
      EmitContext fileCtx) {
    String rt = rtPrefix(fileCtx);
    EmitContext clinitCtx = fileCtx.copyForMember(enumTe, true);
    StringBuilder clinitBody = new StringBuilder();
    List<String> constVarNames = new ArrayList<>();

    // Pre-allocate enum constant pointers FIRST (so circular refs in constructors/static fields see non-nil pointers!)
    for (VariableTree vt : enumConstTrees) {
      String cName = vt.getName().toString();
      String varName = "enum_" + goName + "_" + GoNames.cleanChars(cName);
      constVarNames.add(varName);
      out.append("var ").append(varName).append(" *").append(goName).append('\n');
      out.append("func StaticGet_")
          .append(goName)
          .append('_')
          .append(GoNames.exported(cName))
          .append("() *")
          .append(goName)
          .append(" {\n\tclinit_")
          .append(goName)
          .append("()\n\treturn ")
          .append(varName)
          .append("\n}\n");
    }

    // Initialize enum constants in ordinal order
    for (int ord = 0; ord < enumConstTrees.size(); ord++) {
      VariableTree vt = enumConstTrees.get(ord);
      TreePath vtPath = new TreePath(ctPath, vt);
      String cName = vt.getName().toString();
      String varName = constVarNames.get(ord);
      clinitBody.append('\t').append(varName).append(" = &").append(goName).append("{}\n");

      if (vt.getInitializer() instanceof NewClassTree nct) {
        TreePath nctPath = new TreePath(vtPath, nct);
        // Attach constant-specific method overrides BEFORE running constructor (in case ctor calls virtual method!)
        if (nct.getClassBody() != null) {
          TreePath bodyPath = new TreePath(nctPath, nct.getClassBody());
          for (Tree bm : nct.getClassBody().getMembers()) {
            if (bm instanceof MethodTree bmt && !isConstructorTree(bmt)) {
              TreePath bmtPath = new TreePath(bodyPath, bmt);
              Element bel = resolver.trees.getElement(bmtPath);
              if (bel instanceof ExecutableElement bee
                  && !bee.getModifiers().contains(Modifier.STATIC)) {
                TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(bee);
                NodePredInfo enumNodePred = inspectNodePredicateImpl(enumTe);
                if (enumNodePred != null
                    && sig != null
                    && sig.goName.equals(enumNodePred.primarySamGoName())
                    && bmt.getParameters().size() == 1) {
                  String typedLit =
                      emitEnumConstantNodePredicateFunc(
                          bmt, bmtPath, enumTe, enumNodePred, fileCtx);
                  TypeMirror enumRetTm = resolver.resolveSigReturnType(sig, enumTe.asType());
                  String retGo =
                      enumRetTm != null ? resolver.goType(enumRetTm, fileCtx.imports) : "bool";
                  clinitBody
                      .append('\t')
                      .append(varName)
                      .append(".fn_TestNode_ = ")
                      .append(typedLit)
                      .append('\n')
                      .append('\t')
                      .append(varName)
                      .append(".fn_")
                      .append(sig.goName)
                      .append(" = func(this_ *")
                      .append(goName)
                      .append(", arg_ any) ")
                      .append(retGo)
                      .append(" { return this_.fn_TestNode_(this_, ")
                      .append(rt)
                      .append("CastNode[uint32](arg_)) }\n");
                } else {
                  String fnLit =
                      methodEmitter.emitEnumConstantMethodOverrideFunc(bmt, bmtPath, enumTe, fileCtx);
                  clinitBody
                      .append('\t')
                      .append(varName)
                      .append(".fn_")
                      .append(sig.goName)
                      .append(" = ")
                      .append(fnLit)
                      .append('\n');
                }
              }
            }
          }
        }

        Element ctorEl = resolver.trees.getElement(nctPath);
        ExecutableElement targetCtor = null;
        if (ctorEl instanceof ExecutableElement ee) {
          if (ee.getEnclosingElement().equals(enumTe)) {
            targetCtor = ee;
          } else if (nct.getClassBody() != null) {
            targetCtor =
                findSuperCtorCalledFromSynthesizedAnonInit(
                    new TreePath(nctPath, nct.getClassBody()), enumTe, ee);
          }
        }

        List<String> initArgs = new ArrayList<>();
        initArgs.add(varName);
        initArgs.add(GoNames.goStringLiteral(cName));
        initArgs.add("int32(" + ord + ")");

        String initFn = "initEnum_" + goName;
        if (targetCtor != null) {
          initFn = "initEnum" + resolver.getConstructorGoName(targetCtor).substring(3);
          appendStaticInitCtorArgs(nct.getArguments(), nctPath, targetCtor, initArgs, clinitCtx);
        }
        clinitBody
            .append('\t')
            .append(initFn)
            .append('(')
            .append(String.join(", ", initArgs))
            .append(")\n");
      } else {
        clinitBody
            .append("\tinitEnum_")
            .append(goName)
            .append('(')
            .append(varName)
            .append(", ")
            .append(GoNames.goStringLiteral(cName))
            .append(", int32(")
            .append(ord)
            .append("))\n");
      }
    }

    clinitBody
        .append("\tenumValues_")
        .append(goName)
        .append(" = []*")
        .append(goName)
        .append('{')
        .append(String.join(", ", constVarNames))
        .append("}\n");

    List<String> anyConsts = new ArrayList<>();
    for (String v : constVarNames) {
      anyConsts.add(v);
    }
    clinitBody
        .append('\t')
        .append(rt)
        .append("RegisterEnumConstants(")
        .append(GoNames.goStringLiteral(enumTe.getQualifiedName().toString()))
        .append(", []any{")
        .append(String.join(", ", anyConsts))
        .append("})\n");
    String binEnumName = resolver.elements.getBinaryName(enumTe).toString();
    if (!binEnumName.equals(enumTe.getQualifiedName().toString())) {
      clinitBody
          .append('\t')
          .append(rt)
          .append("RegisterEnumConstants(")
          .append(GoNames.goStringLiteral(binEnumName))
          .append(", []any{")
          .append(String.join(", ", anyConsts))
          .append("})\n");
    }

    // Static fields & `static { ... }` blocks AFTER enum constants (JLS §8.9.2 initialization order!)
    appendStaticMemberDeclsAndClinitStatements(
        ct, ctPath, enumTe, goName, true, clinitCtx, out, clinitBody);

    String atomicPkg = fileCtx.imports.std("sync/atomic");
    out.append("var clinitDone_").append(goName).append(" uint32\n");
    out.append("var clinitState_").append(goName).append(' ').append(rt).append("ClinitState\n");
    out.append("func clinit_").append(goName).append("() {\n");
    out.append("\tif ").append(atomicPkg).append(".LoadUint32(&clinitDone_").append(goName).append(") != 0 {\n\t\treturn\n\t}\n");
    out.append("\tclinitSlow_").append(goName).append("()\n}\n\n");
    out.append("func clinitSlow_").append(goName).append("() {\n");
    out.append("\tif !").append(rt).append("EnterClinit(&clinitState_").append(goName).append(", &clinitDone_").append(goName).append(") {\n\t\treturn\n\t}\n");
    out.append("\tdefer ").append(rt).append("ExitClinit(&clinitState_").append(goName).append(", &clinitDone_").append(goName).append(")\n");
    out.append(clinitBody);
    out.append("}\n\n");
    String enumQnLit = GoNames.goStringLiteral(enumTe.getQualifiedName().toString());
    String enumBinLit = GoNames.goStringLiteral(resolver.elements.getBinaryName(enumTe).toString());
    out.append("func init() {\n\t")
        .append(rt)
        .append("RegisterEnumConstantsLoader(")
        .append(enumQnLit)
        .append(", ")
        .append(enumBinLit)
        .append(", func() {\n\t\tclinit_")
        .append(goName)
        .append("()\n\t})\n}\n\n");

    fileCtx.syncTempSeqFrom(clinitCtx);
    fileCtx.deferredLocalTypes.addAll(clinitCtx.deferredLocalTypes);
  }

  private List<String> tryBuildStaticEnumConstantLiterals(
      ClassTree ct,
      TreePath ctPath,
      TypeElement enumTe,
      List<VariableTree> enumConstTrees,
      List<VariableElement> instanceFields,
      Map<String, TypeResolver.MethodSig> overriddenMethodSigs,
      EmitContext fileCtx) {
    if (!overriddenMethodSigs.isEmpty()
        || resolver.hasOuterInstance(enumTe)
        || !resolver.getCapturedVars(enumTe).isEmpty()) {
      return null;
    }
    List<MethodTree> ctors = new ArrayList<>();
    for (Tree member : ct.getMembers()) {
      if (member instanceof BlockTree bt && !bt.isStatic() && !bt.getStatements().isEmpty()) {
        return null;
      }
      if (member instanceof MethodTree mt && isConstructorTree(mt)) {
        ctors.add(mt);
      }
      if (member instanceof VariableTree vt && !instanceFields.isEmpty()) {
        Element el = resolver.trees.getElement(new TreePath(ctPath, vt));
        if (el instanceof VariableElement ve
            && ve.getKind() == ElementKind.FIELD
            && !ve.getModifiers().contains(Modifier.STATIC)
            && vt.getInitializer() != null) {
          return null;
        }
      }
    }
    for (VariableTree vt : enumConstTrees) {
      if (vt.getInitializer() instanceof NewClassTree nct && nct.getClassBody() != null) {
        return null;
      }
    }

    String rt = rtPrefix(fileCtx);
    if (instanceFields.isEmpty()) {
      for (MethodTree mt : ctors) {
        if (!mt.getParameters().isEmpty() || !isTrivialEmptyEnumCtorBody(mt.getBody())) {
          return null;
        }
      }
      for (VariableTree vt : enumConstTrees) {
        if (vt.getInitializer() instanceof NewClassTree nct && !nct.getArguments().isEmpty()) {
          return null;
        }
      }
      List<String> lits = new ArrayList<>(enumConstTrees.size());
      for (int ord = 0; ord < enumConstTrees.size(); ord++) {
        String cName = enumConstTrees.get(ord).getName().toString();
        lits.add(
            "{EnumBase: "
                + rt
                + "EnumBase{Name_: "
                + GoNames.goStringLiteral(cName)
                + ", Ordinal_: "
                + ord
                + "}}");
      }
      return lits;
    }

    // Simple literal-field enum (`Tri`, `MatchStatus`, etc.): single non-varargs ctor assigning fields directly from params
    for (VariableElement fve : instanceFields) {
      if (!fve.asType().getKind().isPrimitive() && !resolver.isStringType(fve.asType())) {
        return null;
      }
    }
    if (ctors.size() != 1) {
      return null;
    }
    MethodTree ctorMt = ctors.get(0);
    TreePath ctorPath = new TreePath(ctPath, ctorMt);
    if (!(resolver.trees.getElement(ctorPath) instanceof ExecutableElement ctorEe)
        || ctorEe.isVarArgs()
        || ctorEe.getParameters().size() != instanceFields.size()
        || ctorMt.getBody() == null) {
      return null;
    }
    List<? extends VariableElement> ctorParams = ctorEe.getParameters();
    for (VariableElement pve : ctorParams) {
      if (!pve.asType().getKind().isPrimitive() && !resolver.isStringType(pve.asType())) {
        return null;
      }
    }

    Map<VariableElement, Integer> fieldToParamIdx = new LinkedHashMap<>();
    TreePath bodyPath = new TreePath(ctorPath, ctorMt.getBody());
    List<? extends StatementTree> stmts = ctorMt.getBody().getStatements();
    for (int i = 0; i < stmts.size(); i++) {
      StatementTree st = stmts.get(i);
      if (st.getKind() == Tree.Kind.EMPTY_STATEMENT) {
        continue;
      }
      if (i == 0 && isSuperNoArgCallStmt(st)) {
        continue;
      }
      if (!(st instanceof ExpressionStatementTree est)
          || !(est.getExpression() instanceof AssignmentTree at)) {
        return null;
      }
      TreePath atPath = new TreePath(new TreePath(bodyPath, est), at);
      TreePath lhsPath = unwrapParensTreePath(new TreePath(atPath, at.getVariable()));
      TreePath rhsPath = unwrapParensTreePath(new TreePath(atPath, at.getExpression()));
      Tree lhsLeaf = lhsPath.getLeaf();
      if (lhsLeaf instanceof MemberSelectTree ms) {
        ExpressionTree q = unwrapParensExpr(ms.getExpression());
        if (!(q instanceof IdentifierTree qid) || !qid.getName().contentEquals("this")) {
          return null;
        }
      } else if (!(lhsLeaf instanceof IdentifierTree)) {
        return null;
      }
      if (!(rhsPath.getLeaf() instanceof IdentifierTree)) {
        return null;
      }
      Element lhsEl = resolver.trees.getElement(lhsPath);
      Element rhsEl = resolver.trees.getElement(rhsPath);
      if (!(lhsEl instanceof VariableElement lhsVe)
          || !instanceFields.contains(lhsVe)
          || fieldToParamIdx.containsKey(lhsVe)
          || !(rhsEl instanceof VariableElement rhsVe)) {
        return null;
      }
      int pIdx = ctorParams.indexOf(rhsVe);
      if (pIdx < 0) {
        return null;
      }
      fieldToParamIdx.put(lhsVe, pIdx);
    }
    if (fieldToParamIdx.size() != instanceFields.size()) {
      return null;
    }

    for (VariableTree vt : enumConstTrees) {
      if (!(vt.getInitializer() instanceof NewClassTree nct)
          || nct.getArguments().size() != ctorParams.size()) {
        return null;
      }
      for (ExpressionTree arg : nct.getArguments()) {
        if (!isPureGoCompileTimeConstExpr(arg)) {
          return null;
        }
      }
    }

    List<String> lits = new ArrayList<>(enumConstTrees.size());
    for (int ord = 0; ord < enumConstTrees.size(); ord++) {
      VariableTree vt = enumConstTrees.get(ord);
      TreePath vtPath = new TreePath(ctPath, vt);
      NewClassTree nct = (NewClassTree) vt.getInitializer();
      TreePath nctPath = new TreePath(vtPath, nct);
      String cName = vt.getName().toString();
      StringBuilder sb = new StringBuilder();
      sb.append("{EnumBase: ")
          .append(rt)
          .append("EnumBase{Name_: ")
          .append(GoNames.goStringLiteral(cName))
          .append(", Ordinal_: ")
          .append(ord)
          .append('}');
      for (VariableElement fve : instanceFields) {
        int pIdx = fieldToParamIdx.get(fve);
        String valGo =
            exprEmitter.emitAs(nct.getArguments().get(pIdx), nctPath, fve.asType(), false, fileCtx);
        sb.append(", ").append(resolver.getFieldGoName(fve)).append(": ").append(valGo);
      }
      sb.append('}');
      lits.add(sb.toString());
    }
    return lits;
  }

  private static boolean isTrivialEmptyEnumCtorBody(BlockTree body) {
    if (body == null) {
      return true;
    }
    List<? extends StatementTree> stmts = body.getStatements();
    for (int i = 0; i < stmts.size(); i++) {
      StatementTree st = stmts.get(i);
      if (st.getKind() == Tree.Kind.EMPTY_STATEMENT) {
        continue;
      }
      if (i == 0 && isSuperNoArgCallStmt(st)) {
        continue;
      }
      return false;
    }
    return true;
  }

  private static boolean isSuperNoArgCallStmt(StatementTree st) {
    return st instanceof ExpressionStatementTree est
        && est.getExpression() instanceof MethodInvocationTree call0
        && call0.getArguments().isEmpty()
        && call0.getMethodSelect() instanceof IdentifierTree id
        && id.getName().contentEquals("super");
  }

  private static ExpressionTree unwrapParensExpr(ExpressionTree e) {
    while (e instanceof ParenthesizedTree pt) {
      e = pt.getExpression();
    }
    return e;
  }

  private static TreePath unwrapParensTreePath(TreePath tp) {
    while (tp != null && tp.getLeaf() instanceof ParenthesizedTree pt) {
      tp = new TreePath(tp, pt.getExpression());
    }
    return tp;
  }

  private static boolean isPureGoCompileTimeConstExpr(ExpressionTree expr) {
    expr = unwrapParensExpr(expr);
    if (expr instanceof UnaryTree ut) {
      return switch (ut.getKind()) {
        case UNARY_MINUS, UNARY_PLUS, BITWISE_COMPLEMENT ->
            isPureGoCompileTimeConstExpr(ut.getExpression());
        default -> false;
      };
    }
    if (expr instanceof LiteralTree lt) {
      return switch (lt.getKind()) {
        case INT_LITERAL, LONG_LITERAL, BOOLEAN_LITERAL, CHAR_LITERAL, STRING_LITERAL -> true;
        case FLOAT_LITERAL ->
            lt.getValue() instanceof Float f && !f.isNaN() && !f.isInfinite();
        case DOUBLE_LITERAL ->
            lt.getValue() instanceof Double d && !d.isNaN() && !d.isInfinite();
        default -> false;
      };
    }
    return false;
  }

  private static boolean hasNonEmptyStaticInitBlock(ClassTree ct) {
    for (Tree member : ct.getMembers()) {
      if (member instanceof BlockTree bt && bt.isStatic() && !bt.getStatements().isEmpty()) {
        return true;
      }
    }
    return false;
  }

  private boolean emitStaticEnumConstantsAndClinit(
      ClassTree ct,
      TreePath ctPath,
      TypeElement enumTe,
      String goName,
      List<VariableTree> enumConstTrees,
      List<String> staticConstLiterals,
      boolean hasStaticInitBlock,
      StringBuilder out,
      EmitContext fileCtx) {
    String rt = rtPrefix(fileCtx);
    int n = enumConstTrees.size();
    String backingName = "enumBacking_" + goName;

    if (n > 0) {
      out.append("var ").append(backingName).append(" = [").append(n).append(']').append(goName).append("{\n");
      for (String lit : staticConstLiterals) {
        out.append('\t').append(lit).append(",\n");
      }
      out.append("}\n");
    }

    List<String> sliceElems = new ArrayList<>(n);
    for (int ord = 0; ord < n; ord++) {
      String cName = enumConstTrees.get(ord).getName().toString();
      String varName = "enum_" + goName + "_" + GoNames.cleanChars(cName);
      String addrExpr = "&" + backingName + "[" + ord + "]";
      sliceElems.add(addrExpr);
      out.append("var ").append(varName).append(" *").append(goName).append(" = ").append(addrExpr).append('\n');
      out.append("func StaticGet_")
          .append(goName)
          .append('_')
          .append(GoNames.exported(cName))
          .append("() *")
          .append(goName)
          .append(" {\n");
      if (hasStaticInitBlock) {
        out.append("\tclinit_").append(goName).append("()\n");
      }
      out.append("\treturn ").append(addrExpr).append("\n}\n");
    }

    out.append("var enumValues_")
        .append(goName)
        .append(" = []*")
        .append(goName)
        .append('{')
        .append(String.join(", ", sliceElems))
        .append("}\n");

    EmitContext clinitCtx = fileCtx.copyForMember(enumTe, true);
    StringBuilder clinitBody = new StringBuilder();
    boolean needsStaticClinit = hasStaticMemberClinit(resolver, ct, ctPath);
    appendStaticMemberDeclsAndClinitStatements(
        ct, ctPath, enumTe, goName, needsStaticClinit, clinitCtx, out, clinitBody);

    fileCtx.syncTempSeqFrom(clinitCtx);
    fileCtx.deferredLocalTypes.addAll(clinitCtx.deferredLocalTypes);

    boolean hasClinit = needsStaticClinit || clinitBody.length() > 0;
    if (hasClinit) {
      String atomicPkg = fileCtx.imports.std("sync/atomic");
      out.append("var clinitDone_").append(goName).append(" uint32\n");
      out.append("var clinitState_").append(goName).append(' ').append(rt).append("ClinitState\n");
      out.append("func clinit_").append(goName).append("() {\n");
      out.append("\tif ").append(atomicPkg).append(".LoadUint32(&clinitDone_").append(goName).append(") != 0 {\n\t\treturn\n\t}\n");
      out.append("\tclinitSlow_").append(goName).append("()\n}\n\n");
      out.append("func clinitSlow_").append(goName).append("() {\n");
      out.append("\tif !").append(rt).append("EnterClinit(&clinitState_").append(goName).append(", &clinitDone_").append(goName).append(") {\n\t\treturn\n\t}\n");
      out.append("\tdefer ").append(rt).append("ExitClinit(&clinitState_").append(goName).append(", &clinitDone_").append(goName).append(")\n");
      out.append(clinitBody);
      out.append("}\n\n");
    }

    String enumQn = enumTe.getQualifiedName().toString();
    String binEnumName = resolver.elements.getBinaryName(enumTe).toString();
    String enumQnLit = GoNames.goStringLiteral(enumQn);
    String enumBinLit = GoNames.goStringLiteral(binEnumName);
    out.append("func init() {\n\t")
        .append(rt)
        .append("RegisterEnumConstantsLoader(")
        .append(enumQnLit)
        .append(", ")
        .append(enumBinLit)
        .append(", func() {\n");
    if (hasStaticInitBlock && hasClinit) {
      out.append("\t\tclinit_").append(goName).append("()\n");
    }
    out.append("\t\tvals_ := make([]any, len(enumValues_").append(goName).append("))\n");
    out.append("\t\tfor i_, v_ := range enumValues_").append(goName).append(" {\n");
    out.append("\t\t\tvals_[i_] = v_\n\t\t}\n");
    out.append("\t\t").append(rt).append("RegisterEnumConstants(").append(enumQnLit).append(", vals_)\n");
    if (!binEnumName.equals(enumQn)) {
      out.append("\t\t").append(rt).append("RegisterEnumConstants(").append(enumBinLit).append(", vals_)\n");
    }
    out.append("\t})\n}\n\n");

    return hasClinit;
  }

  private void appendStaticInitCtorArgs(
      List<? extends ExpressionTree> argTrees,
      TreePath callPath,
      ExecutableElement ctor,
      List<String> outArgs,
      EmitContext ctx) {
    int n = ctor.getParameters().size();
    boolean varArgs = ctor.isVarArgs() && n > 0;
    int fixed = varArgs ? (n - 1) : n;
    for (int i = 0; i < fixed && i < argTrees.size(); i++) {
      outArgs.add(
          exprEmitter.emitAs(
              argTrees.get(i),
              callPath,
              resolver.safeErasure(ctor.getParameters().get(i).asType()),
              false,
              ctx));
    }
    if (varArgs) {
      javax.lang.model.type.ArrayType at =
          (javax.lang.model.type.ArrayType)
              resolver.safeErasure(ctor.getParameters().get(n - 1).asType());
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

  // =========================================================================
  // 4. Static Fields & `clinit_<GoName>()`
  // =========================================================================

  private boolean emitStaticFieldsAndClinit(
      ClassTree ct,
      TreePath ctPath,
      TypeElement te,
      Map<VariableElement, String> capturedFields,
      StringBuilder out,
      EmitContext fileCtx) {
    String goName = resolver.planner.getGoTypeName(te);
    EmitContext clinitCtx = fileCtx.copyForMember(te, true);
    clinitCtx.capturedVarFields.putAll(capturedFields);
    StringBuilder clinitBody = new StringBuilder();

    boolean needsRuntimeClinit = hasRuntimeClinit(resolver, ct, ctPath, te);
    appendStaticMemberDeclsAndClinitStatements(
        ct, ctPath, te, goName, needsRuntimeClinit, clinitCtx, out, clinitBody);

    fileCtx.syncTempSeqFrom(clinitCtx);
    fileCtx.deferredLocalTypes.addAll(clinitCtx.deferredLocalTypes);

    if (!needsRuntimeClinit && clinitBody.length() == 0) {
      return false;
    }

    String atomicPkg = fileCtx.imports.std("sync/atomic");
    String rt = rtPrefix(fileCtx);
    out.append("var clinitDone_").append(goName).append(" uint32\n");
    out.append("var clinitState_").append(goName).append(' ').append(rt).append("ClinitState\n");
    out.append("func clinit_").append(goName).append("() {\n");
    out.append("\tif ").append(atomicPkg).append(".LoadUint32(&clinitDone_").append(goName).append(") != 0 {\n\t\treturn\n\t}\n");
    out.append("\tclinitSlow_").append(goName).append("()\n}\n\n");
    out.append("func clinitSlow_").append(goName).append("() {\n");
    out.append("\tif !").append(rt).append("EnterClinit(&clinitState_").append(goName).append(", &clinitDone_").append(goName).append(") {\n\t\treturn\n\t}\n");
    out.append("\tdefer ").append(rt).append("ExitClinit(&clinitState_").append(goName).append(", &clinitDone_").append(goName).append(")\n");
    out.append(clinitBody);
    out.append("}\n\n");
    return true;
  }

  private void emitValueEnum(
      ClassTree ct,
      TreePath ctPath,
      TypeElement te,
      Set<TypeElement> emittedTypes,
      StringBuilder out,
      EmitContext fileCtx) {
    String goName = resolver.planner.getGoTypeName(te);
    String rt = rtPrefix(fileCtx);
    String jl = resolver.runtime.javalangAlias(fileCtx.imports);
    String jlPref = jl.isEmpty() ? "" : (jl + ".");
    String gsonAlias = fileCtx.imports.ext(RuntimeBindings.SHIM_BASE + "/gson", "gson");
    String gsonPref = gsonAlias.isEmpty() ? "" : (gsonAlias + ".");

    List<VariableTree> enumConstTrees = new ArrayList<>();
    for (Tree member : ct.getMembers()) {
      if (member instanceof VariableTree vt) {
        Element el = resolver.trees.getElement(new TreePath(ctPath, vt));
        if (el instanceof VariableElement ve && ve.getKind() == ElementKind.ENUM_CONSTANT) {
          enumConstTrees.add(vt);
        }
      }
    }
    int n = enumConstTrees.size();

    out.append("// Java enum: ").append(te.getQualifiedName()).append('\n');
    out.append("type ").append(goName).append(" uint8\n\n");

    out.append("const (\n");
    for (int ord = 0; ord < n; ord++) {
      String cName = enumConstTrees.get(ord).getName().toString();
      out.append("\tenum_")
          .append(goName)
          .append('_')
          .append(GoNames.cleanChars(cName))
          .append(' ')
          .append(goName)
          .append(" = ")
          .append(ord + 1)
          .append('\n');
    }
    out.append(")\n\n");

    out.append("var enumBaseTable_").append(goName).append(" = [...]").append(rt).append("EnumBase{\n");
    List<String> valSliceElems = new ArrayList<>(n);
    for (int ord = 0; ord < n; ord++) {
      String cName = enumConstTrees.get(ord).getName().toString();
      out.append("\t{Name_: ")
          .append(GoNames.goStringLiteral(cName))
          .append(", Ordinal_: ")
          .append(ord)
          .append("},\n");
      valSliceElems.add(Integer.toString(ord + 1));
    }
    out.append("}\n");
    out.append("var enumValues_")
        .append(goName)
        .append(" = []")
        .append(goName)
        .append('{')
        .append(String.join(", ", valSliceElems))
        .append("}\n\n");

    Map<VariableElement, List<Object>> fieldTable = resolver.getValueEnumFieldTable(te);
    for (Map.Entry<VariableElement, List<Object>> fe : fieldTable.entrySet()) {
      VariableElement fve = fe.getKey();
      List<Object> vals = fe.getValue();
      String fGo = resolver.getFieldGoName(fve);
      String fTypeGo = resolver.goType(fve.asType(), fileCtx.imports);
      String tableIdent = "enumField_" + goName + "_" + fGo;
      out.append("var ").append(tableIdent).append(" = [...]").append(fTypeGo).append("{\n");
      for (Object v : vals) {
        String elemGo =
            (v == TypeResolver.ENUM_NULL_STRING)
                ? (rt + "NullString")
                : exprEmitter.emitConstantValue(v, fve.asType(), fileCtx);
        out.append('\t').append(elemGo).append(",\n");
      }
      out.append("}\n");
      out.append("func (this_ ")
          .append(goName)
          .append(") FieldGet_")
          .append(fGo)
          .append("() ")
          .append(fTypeGo)
          .append(" {\n\tif this_ == 0 {\n\t\t")
          .append(rt)
          .append("ThrowNPE()\n\t}\n\treturn ")
          .append(tableIdent)
          .append("[this_-1]\n}\n\n");
    }

    boolean hasStaticInitBlock = hasNonEmptyStaticInitBlock(ct);
    boolean needsStaticClinit = hasStaticMemberClinit(resolver, ct, ctPath);
    EmitContext clinitCtx = fileCtx.copyForMember(te, true);
    StringBuilder staticDeclsOut = new StringBuilder();
    StringBuilder clinitBody = new StringBuilder();
    appendStaticMemberDeclsAndClinitStatements(
        ct, ctPath, te, goName, needsStaticClinit, clinitCtx, staticDeclsOut, clinitBody);
    fileCtx.syncTempSeqFrom(clinitCtx);
    fileCtx.deferredLocalTypes.addAll(clinitCtx.deferredLocalTypes);
    boolean hasClinit = needsStaticClinit || clinitBody.length() > 0;

    for (VariableTree vt : enumConstTrees) {
      String cName = vt.getName().toString();
      String constIdent = "enum_" + goName + "_" + GoNames.cleanChars(cName);
      out.append("func StaticGet_")
          .append(goName)
          .append('_')
          .append(GoNames.exported(cName))
          .append("() ")
          .append(goName)
          .append(" {\n");
      if (hasStaticInitBlock && hasClinit) {
        out.append("\tclinit_").append(goName).append("()\n");
      }
      out.append("\treturn ").append(constIdent).append("\n}\n");
    }
    out.append('\n');
    out.append(staticDeclsOut);
    String enumQn = te.getQualifiedName().toString();
    String binEnumName = resolver.elements.getBinaryName(te).toString();
    String enumQnLit = GoNames.goStringLiteral(enumQn);
    String enumBinLit = GoNames.goStringLiteral(binEnumName);
    List<String> anyVals = new ArrayList<>(n);
    for (int ord = 0; ord < n; ord++) {
      anyVals.add(goName + "(" + (ord + 1) + ")");
    }
    if (hasClinit) {
      String atomicPkg = fileCtx.imports.std("sync/atomic");
      out.append("var clinitDone_").append(goName).append(" uint32\n");
      out.append("var clinitState_").append(goName).append(' ').append(rt).append("ClinitState\n");
      out.append("func clinit_").append(goName).append("() {\n");
      out.append("\tif ")
          .append(atomicPkg)
          .append(".LoadUint32(&clinitDone_")
          .append(goName)
          .append(") != 0 {\n\t\treturn\n\t}\n");
      out.append("\tclinitSlow_").append(goName).append("()\n}\n\n");
      out.append("func clinitSlow_").append(goName).append("() {\n");
      out.append("\tif !")
          .append(rt)
          .append("EnterClinit(&clinitState_")
          .append(goName)
          .append(", &clinitDone_")
          .append(goName)
          .append(") {\n\t\treturn\n\t}\n");
      out.append("\tdefer ")
          .append(rt)
          .append("ExitClinit(&clinitState_")
          .append(goName)
          .append(", &clinitDone_")
          .append(goName)
          .append(")\n");
      if (hasStaticInitBlock) {
        out.append("\tvals_ := []any{").append(String.join(", ", anyVals)).append("}\n");
        out.append('\t')
            .append(rt)
            .append("RegisterEnumConstants(")
            .append(enumQnLit)
            .append(", vals_)\n");
        if (!binEnumName.equals(enumQn)) {
          out.append('\t')
              .append(rt)
              .append("RegisterEnumConstants(")
              .append(enumBinLit)
              .append(", vals_)\n");
        }
      }
      out.append(clinitBody);
      out.append("}\n\n");
    }

    out.append("func ").append(goName).append("_Values() []").append(goName).append(" {\n");
    if (hasClinit) {
      out.append("\tclinit_").append(goName).append("()\n");
    }
    out.append("\treturn ").append(rt).append("CloneSlice(enumValues_").append(goName).append(")\n}\n\n");

    out.append("func ").append(goName).append("_ValueOf(name string) ").append(goName).append(" {\n");
    if (hasClinit) {
      out.append("\tclinit_").append(goName).append("()\n");
    }
    out.append("\tif name == ")
        .append(rt)
        .append("NullString {\n\t\t")
        .append(rt)
        .append("ThrowNPE()\n\t}\n");
    out.append("\tfor i_ := range enumBaseTable_").append(goName).append(" {\n");
    out.append("\t\tif enumBaseTable_")
        .append(goName)
        .append("[i_].Name_ == name {\n\t\t\treturn ")
        .append(goName)
        .append("(i_ + 1)\n\t\t}\n\t}\n");
    out.append("\tpanic(")
        .append(jlPref)
        .append("NewIllegalArgumentException(\"No enum constant ")
        .append(te.getQualifiedName())
        .append(".\" + name))\n}\n\n");

    Set<String> emittedMethodGoNames = new HashSet<>();
    emittedMethodGoNames.add("Name");
    emittedMethodGoNames.add("Ordinal");
    emittedMethodGoNames.add("HashCode");
    emittedMethodGoNames.add("Equals");
    emittedMethodGoNames.add("CompareTo");
    for (Tree member : ct.getMembers()) {
      if (member instanceof MethodTree mt && !isConstructorTree(mt)) {
        TreePath mtPath = new TreePath(ctPath, mt);
        Element el = resolver.trees.getElement(mtPath);
        if (el instanceof ExecutableElement ee && !ee.getModifiers().contains(Modifier.STATIC)) {
          TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
          if (sig != null) {
            emittedMethodGoNames.add(sig.goName);
          }
        }
        methodEmitter.emitMethod(mt, mtPath, te, hasClinit, null, Map.of(), out, fileCtx);
      }
    }

    String simpleName = te.getSimpleName().toString();
    out.append("func (this_ ").append(goName).append(") IsJavaNil_() bool { return this_ == 0 }\n");
    out.append("func (").append(goName).append(") HasIdentityEquals_() {}\n");
    out.append("func (").append(goName).append(") Implements_j2gort_Ext_java_lang_Enum() {}\n");
    out.append("func (")
        .append(goName)
        .append(") Implements_j2gort_Ext_java_lang_constant_Constable() {}\n");
    out.append("func (this_ ")
        .append(goName)
        .append(") GetExt_java_lang_EnumBase() *")
        .append(rt)
        .append("Ext_java_lang_EnumBase { return &")
        .append(rt)
        .append("Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()} }\n");
    out.append("func (")
        .append(goName)
        .append(") JavaClassInfo_() (string, string) { return ")
        .append(GoNames.goStringLiteral(binEnumName))
        .append(", ")
        .append(GoNames.goStringLiteral(simpleName))
        .append(" }\n");
    out.append("func (this_ ")
        .append(goName)
        .append(") GetEnumBase() *")
        .append(rt)
        .append("EnumBase {\n\tif this_ == 0 {\n\t\treturn nil\n\t}\n\treturn &enumBaseTable_")
        .append(goName)
        .append("[this_-1]\n}\n");
    out.append("func (this_ ")
        .append(goName)
        .append(") Ordinal() int32 {\n\tif this_ == 0 {\n\t\t")
        .append(rt)
        .append("ThrowNPE()\n\t}\n\treturn int32(this_ - 1)\n}\n");
    out.append("func (this_ ")
        .append(goName)
        .append(") EnumOrdinal_() int32 { return int32(this_) - 1 }\n");
    out.append("func (this_ ")
        .append(goName)
        .append(") Name() string {\n\tif this_ == 0 {\n\t\t")
        .append(rt)
        .append("ThrowNPE()\n\t}\n\treturn enumBaseTable_")
        .append(goName)
        .append("[this_-1].Name_\n}\n");
    if (!emittedMethodGoNames.contains("String")) {
      out.append("func (this_ ")
          .append(goName)
          .append(") String() string {\n\tif this_ == 0 {\n\t\treturn \"null\"\n\t}\n\treturn enumBaseTable_")
          .append(goName)
          .append("[this_-1].Name_\n}\n");
    }
    out.append("func (this_ ")
        .append(goName)
        .append(") Equals(other_ any) bool {\n\to_, ok_ := other_.(")
        .append(goName)
        .append(")\n\treturn this_ != 0 && ok_ && this_ == o_\n}\n");
    out.append("func (this_ ")
        .append(goName)
        .append(") HashCode() int32 {\n\tif this_ == 0 {\n\t\treturn 0\n\t}\n\treturn enumBaseTable_")
        .append(goName)
        .append("[this_-1].HashCode()\n}\n");
    out.append("func (this_ ")
        .append(goName)
        .append(") CompareTo(other_ any) int32 {\n\tif this_ == 0 {\n\t\t")
        .append(rt)
        .append("ThrowNPE()\n\t}\n\treturn enumBaseTable_")
        .append(goName)
        .append("[this_-1].CompareTo(other_)\n}\n");
    out.append("func (this_ ")
        .append(goName)
        .append(") WriteGson(out_ *")
        .append(gsonPref)
        .append("JsonWriter) {\n\tif this_ == 0 {\n\t\tout_.NullValue()\n\t\treturn\n\t}\n\tout_.ValueString(enumBaseTable_")
        .append(goName)
        .append("[this_-1].Name_)\n}\n");

    if (hasStaticInitBlock && hasClinit) {
      out.append("func init() {\n\t")
          .append(rt)
          .append("RegisterEnumConstantsLoader(")
          .append(enumQnLit)
          .append(", ")
          .append(enumBinLit)
          .append(", func() {\n\t\tvals_ := []any{")
          .append(String.join(", ", anyVals))
          .append("}\n\t\t")
          .append(rt)
          .append("RegisterEnumConstants(")
          .append(enumQnLit)
          .append(", vals_)\n");
      if (!binEnumName.equals(enumQn)) {
        out.append("\t\t")
            .append(rt)
            .append("RegisterEnumConstants(")
            .append(enumBinLit)
            .append(", vals_)\n");
      }
      out.append("\t\tclinit_").append(goName).append("()\n\t})\n}\n");
    } else {
      out.append("func init() {\n\tvals_ := []any{")
          .append(String.join(", ", anyVals))
          .append("}\n\t")
          .append(rt)
          .append("RegisterEnumConstants(")
          .append(enumQnLit)
          .append(", vals_)\n");
      if (!binEnumName.equals(enumQn)) {
        out.append('\t')
            .append(rt)
            .append("RegisterEnumConstants(")
            .append(enumBinLit)
            .append(", vals_)\n");
      }
      out.append("}\n");
    }
    emitTypePredicateRegistration(te, goName, out, fileCtx);
    out.append('\n');
    emitNestedMemberTypes(ct, ctPath, emittedTypes, out, fileCtx);
  }

  static boolean hasRuntimeClinit(
      TypeResolver resolver, ClassTree ct, TreePath ctPath, TypeElement te) {
    if (te.getKind() == ElementKind.ENUM) {
      return !resolver.isValueEnum(te) || hasStaticMemberClinit(resolver, ct, ctPath);
    }
    return hasStaticMemberClinit(resolver, ct, ctPath);
  }

  static boolean hasEagerMethodEntryClinit(
      TypeResolver resolver, ClassTree ct, TreePath ctPath, TypeElement te) {
    return hasRuntimeClinit(resolver, ct, ctPath, te)
        && !isPureStaticFieldOnlyClinit(resolver, ct, ctPath, te);
  }

  private static boolean isPureStaticFieldOnlyClinit(
      TypeResolver resolver, ClassTree ct, TreePath ctPath, TypeElement te) {
    if (te.getKind() != ElementKind.CLASS || hasNonEmptyStaticInitBlock(ct)) {
      return false;
    }
    for (Tree member : ct.getMembers()) {
      if (member instanceof VariableTree vt) {
        TreePath vtPath = new TreePath(ctPath, vt);
        Element el = resolver.trees.getElement(vtPath);
        if (el instanceof VariableElement ve
            && ve.getKind() == ElementKind.FIELD
            && ve.getModifiers().contains(Modifier.STATIC)) {
          boolean isCompileConst =
              ve.getConstantValue() != null && ve.getModifiers().contains(Modifier.FINAL);
          if (isCompileConst) {
            continue;
          }
          if (!ve.getModifiers().contains(Modifier.FINAL) || vt.getInitializer() == null) {
            return false;
          }
          if (!isPureImmutableInitExpr(
              resolver, te, new TreePath(vtPath, vt.getInitializer()))) {
            return false;
          }
        }
      }
    }
    return true;
  }

  private static boolean isPureImmutableInitExpr(
      TypeResolver resolver, TypeElement ownerTe, TreePath exprPath) {
    exprPath = unwrapParensTreePath(exprPath);
    if (exprPath == null) {
      return false;
    }
    Tree leaf = exprPath.getLeaf();
    if (leaf instanceof LiteralTree || leaf instanceof LambdaExpressionTree) {
      return true;
    }
    if (leaf instanceof UnaryTree ut) {
      return switch (ut.getKind()) {
        case UNARY_MINUS, UNARY_PLUS, BITWISE_COMPLEMENT, LOGICAL_COMPLEMENT ->
            isPureImmutableInitExpr(resolver, ownerTe, new TreePath(exprPath, ut.getExpression()));
        default -> false;
      };
    }
    if (leaf instanceof MemberReferenceTree mrt) {
      Element qualEl =
          resolver.trees.getElement(new TreePath(exprPath, mrt.getQualifierExpression()));
      return qualEl instanceof TypeElement;
    }
    if (leaf instanceof IdentifierTree || leaf instanceof MemberSelectTree) {
      if (leaf instanceof MemberSelectTree ms) {
        Element qualEl = resolver.trees.getElement(new TreePath(exprPath, ms.getExpression()));
        if (!(qualEl instanceof TypeElement)) {
          return false;
        }
      }
      Element el = resolver.trees.getElement(exprPath);
      if (el instanceof VariableElement ve) {
        if (ve.getKind() == ElementKind.ENUM_CONSTANT) {
          Element enc = ve.getEnclosingElement();
          return enc instanceof TypeElement ete
              && (resolver.isValueEnum(ete)
                  || !resolver.isSourceType(ete)
                  || !resolver.planner.getGoPkg(ete).equals(resolver.planner.getGoPkg(ownerTe)));
        }
        return ve.getModifiers().contains(Modifier.STATIC)
            && ve.getModifiers().contains(Modifier.FINAL)
            && ve.getConstantValue() != null;
      }
      return false;
    }
    if (leaf instanceof MethodInvocationTree mit) {
      Element el = resolver.trees.getElement(exprPath);
      if (!(el instanceof ExecutableElement ee)
          || !ee.getModifiers().contains(Modifier.STATIC)
          || !(ee.getEnclosingElement() instanceof TypeElement owner)) {
        return false;
      }
      String qn = owner.getQualifiedName().toString();
      String mn = ee.getSimpleName().toString();
      List<? extends ExpressionTree> args = mit.getArguments();
      for (ExpressionTree arg : args) {
        if (!isPureImmutableInitExpr(resolver, ownerTe, new TreePath(exprPath, arg))) {
          return false;
        }
      }
      if (qn.equals("com.google.common.collect.ImmutableSet")
          && (mn.equals("of") || mn.equals("copyOf"))) {
        for (ExpressionTree arg : args) {
          if (unwrapParensExpr(arg).getKind() == Tree.Kind.NULL_LITERAL) {
            return false;
          }
        }
        return true;
      }
      if (qn.equals("com.google.common.collect.ImmutableList")
          && (mn.equals("of") || mn.equals("copyOf"))) {
        for (ExpressionTree arg : args) {
          if (unwrapParensExpr(arg).getKind() == Tree.Kind.NULL_LITERAL) {
            return false;
          }
        }
        return true;
      }
      if (qn.equals("com.google.common.collect.Sets") && mn.equals("immutableEnumSet")) {
        for (ExpressionTree arg : args) {
          if (unwrapParensExpr(arg).getKind() == Tree.Kind.NULL_LITERAL) {
            return false;
          }
        }
        return true;
      }
      if (qn.equals("com.google.javascript.rhino.QualifiedName")
          && mn.equals("of")
          && args.size() == 1
          && unwrapParensExpr(args.get(0)).getKind() == Tree.Kind.STRING_LITERAL) {
        return true;
      }
      if (qn.equals("com.google.javascript.rhino.IR")
          && ((mn.equals("name")
                  && args.size() == 1
                  && unwrapParensExpr(args.get(0)).getKind() == Tree.Kind.STRING_LITERAL)
              || (mn.equals("string")
                  && args.size() == 1
                  && unwrapParensExpr(args.get(0)).getKind() == Tree.Kind.STRING_LITERAL)
              || (mn.equals("getprop")
                  && args.size() == 2
                  && unwrapParensExpr(args.get(1)).getKind() == Tree.Kind.STRING_LITERAL))) {
        return true;
      }
    }
    return false;
  }

  private static boolean hasStaticMemberClinit(
      TypeResolver resolver, ClassTree ct, TreePath ctPath) {
    for (Tree member : ct.getMembers()) {
      if (member instanceof BlockTree bt && bt.isStatic() && !bt.getStatements().isEmpty()) {
        return true;
      }
      if (member instanceof VariableTree vt) {
        Element el = resolver.trees.getElement(new TreePath(ctPath, vt));
        if (el instanceof VariableElement ve
            && ve.getKind() == ElementKind.FIELD
            && ve.getModifiers().contains(Modifier.STATIC)) {
          boolean isCompileConst =
              ve.getConstantValue() != null && ve.getModifiers().contains(Modifier.FINAL);
          if (!isCompileConst
              && (vt.getInitializer() != null || resolver.isStringType(ve.asType()))) {
            return true;
          }
        }
      }
    }
    return false;
  }

  private void appendStaticMemberDeclsAndClinitStatements(
      ClassTree ct,
      TreePath ctPath,
      TypeElement te,
      String goName,
      boolean needsRuntimeClinit,
      EmitContext clinitCtx,
      StringBuilder declsOut,
      StringBuilder clinitBodyOut) {
    if (resolver.hasFastStaticSyncMutex(te)) {
      String syncPkg = clinitCtx.imports.std("sync");
      declsOut.append("var staticSyncMu_").append(goName).append(' ').append(syncPkg).append(".Mutex\n");
    }
    // Emit `var static_...`, `StaticGet_...`, `StaticSet_...` and append ordered `clinit` statements
    for (Tree member : ct.getMembers()) {
      if (member instanceof VariableTree vt) {
        TreePath vtPath = new TreePath(ctPath, vt);
        Element el = resolver.trees.getElement(vtPath);
        if (el instanceof VariableElement ve
            && ve.getKind() == ElementKind.FIELD
            && ve.getModifiers().contains(Modifier.STATIC)) {
          String expField = GoNames.exported(ve.getSimpleName().toString());
          String fType = resolver.goType(ve.asType(), clinitCtx.imports);
          Object constVal = ve.getConstantValue();

          if (constVal != null && ve.getModifiers().contains(Modifier.FINAL)) {
            String lit = exprEmitter.emitConstantValue(constVal, ve.asType(), clinitCtx);
            declsOut
                .append("func StaticGet_")
                .append(goName)
                .append('_')
                .append(expField)
                .append("() ")
                .append(fType)
                .append(" { return ")
                .append(lit)
                .append(" }\n");
            continue;
          }

          String varName = "static_" + goName + "_" + GoNames.cleanChars(ve.getSimpleName().toString());
          declsOut.append("var ").append(varName).append(' ').append(fType).append('\n');

          declsOut
              .append("func StaticGet_")
              .append(goName)
              .append('_')
              .append(expField)
              .append("() ")
              .append(fType)
              .append(" {\n");
          if (needsRuntimeClinit) {
            declsOut.append("\tclinit_").append(goName).append("()\n");
          }
          declsOut.append("\treturn ").append(varName).append("\n}\n");

          declsOut
              .append("func StaticSet_")
              .append(goName)
              .append('_')
              .append(expField)
              .append("(v ")
              .append(fType)
              .append(") {\n");
          if (needsRuntimeClinit) {
            declsOut.append("\tclinit_").append(goName).append("()\n");
          }
          declsOut.append('\t').append(varName).append(" = v\n}\n");

          if (vt.getInitializer() != null) {
            String initExpr =
                exprEmitter.emitAs(vt.getInitializer(), vtPath, ve.asType(), false, clinitCtx);
            clinitBodyOut.append('\t').append(varName).append(" = ").append(initExpr).append('\n');
          } else if (resolver.isStringType(ve.asType())) {
            clinitBodyOut
                .append('\t')
                .append(varName)
                .append(" = ")
                .append(rtPrefix(clinitCtx))
                .append("NullString\n");
          }
        }
      } else if (member instanceof BlockTree bt && bt.isStatic()) {
        clinitBodyOut.append(
            stmtEmitter.emitBlockStatements(bt, new TreePath(ctPath, bt), 1, clinitCtx));
      }
    }
  }

  // =========================================================================
  // 5. Special `com.google.javascript.rhino.Node` Lowering onto `nodestore`
  // =========================================================================

  private void emitNodeType(
      ClassTree ct,
      TreePath ctPath,
      TypeElement te,
      Set<TypeElement> emittedTypes,
      StringBuilder out,
      EmitContext fileCtx) {
    int kindFlag = resolver.getNodeKindFlag(te);

    if (kindFlag == 0) {
      // Top-level `com.google.javascript.rhino.Node`
      String ns = resolver.runtime.nodestoreAlias(fileCtx.imports);
      String nsPref = ns.isEmpty() ? "" : (ns + ".");
      String jl = resolver.runtime.javalangAlias(fileCtx.imports);
      String jlPref = jl.isEmpty() ? "" : (jl + ".");
      String rt = rtPrefix(fileCtx);

      out.append("// Node is backed directly by transpiled/shims/nodestore's flat 32-byte NodeData arena.\n");
      out.append("type Node ").append(nsPref).append("Node\n");
      out.append("type Node_NumberNode = Node\n");
      out.append("type Node_BigIntNode = Node\n");
      out.append("type Node_StringNode = Node\n");
      out.append("type Node_TemplateLiteralSubstringNode = Node\n\n");
      out.append("var nodeNullPtrEx any = ").append(jlPref).append("NewNullPointerException()\n\n");

      // Collect dynamic Go types of auxiliary fields from `Node` and its 4 subclasses
      String jstypeOrColorType = findNodeFieldGoType(te, "jstypeOrColor", "any", fileCtx);
      String propListHeadType = findNodeFieldGoType(te, "propListHead", "*Node_PropListItem", fileCtx);
      String sourceFileType = findNodeFieldGoType(te, "sourceFile", "StaticSourceFile", fileCtx);
      String jsDocInfoType = findNodeFieldGoType(te, "jsDocInfo", "*JSDocInfo", fileCtx);
      TypeElement jsTypeTeForNode =
          resolver.elements.getTypeElement("com.google.javascript.rhino.jstype.JSType");
      String defaultTypeBeforeCastGo =
          jsTypeTeForNode != null
              ? resolver.goType(jsTypeTeForNode.asType(), fileCtx.imports)
              : "any";
      String typeBeforeCastType =
          findNodeFieldGoType(te, "typeBeforeCast", defaultTypeBeforeCastGo, fileCtx);
      String bigIntType = findNodeSubFieldGoType(ct, ctPath, "BigIntNode", "bigint", "any", fileCtx);

      out.append("type nodeExtraFields struct {\n");
      out.append("\tbigint ").append(bigIntType).append('\n');
      out.append("\tcooked string\n");
      out.append("\traw    string\n");
      out.append("}\n\n");

      out.append("func allocNodeHandle(kind uint8) Node {\n");
      out.append("\treturn Node(").append(nsPref).append("AllocNodeInit(kind))\n}\n\n");

      out.append("func allocNodeHandleToken(kind uint8, token uint8) Node {\n");
      out.append("\treturn Node(").append(nsPref).append("AllocNodeToken(kind, token))\n}\n\n");

      TypeElement tokenTeForNode =
          resolver.elements.getTypeElement("com.google.javascript.rhino.Token");
      boolean valueToken = tokenTeForNode != null && resolver.isValueEnum(tokenTeForNode);
      if (valueToken) {
        out.append("func tokenToByte(v Token) uint8 {\n");
        out.append("\treturn uint8(v) - 1\n}\n\n");
      } else {
        out.append("func tokenToByte(v *Token) uint8 {\n");
        out.append("\tif v == nil {\n\t\treturn 0xFF\n\t}\n");
        out.append("\treturn uint8(v.Ordinal())\n}\n\n");
      }

      out.append("//go:noinline\n");
      out.append("func (this_ Node) checkDetachedSlow_() {\n");
      out.append("\trec := ").append(nsPref).append("Node(this_).Rec()\n");
      out.append("\tif rec.Parent != 0 {\n\t\t")
          .append(rt)
          .append("CheckStateFmt1(false, \"Has parent: %s\", this_)\n\t}\n");
      out.append("\tif rec.Next != 0 {\n\t\t")
          .append(rt)
          .append("CheckStateFmt1(false, \"Has next: %s\", this_)\n\t}\n");
      out.append("\tif rec.Previous != 0 {\n\t\t")
          .append(rt)
          .append("CheckStateFmt1(false, \"Has previous: %s\", this_)\n\t}\n");
      out.append("}\n\n");

      out.append("func checkDetachedRec_(n Node) *").append(nsPref).append("NodeData {\n");
      out.append("\tif n == 0 {\n\t\tpanic(nodeNullPtrEx)\n\t}\n");
      out.append("\trec := ").append(nsPref).append("Node(n).Rec()\n");
      out.append("\tif (rec.Parent | rec.Next | rec.Previous) != 0 {\n\t\tn.checkDetachedSlow_()\n\t}\n");
      out.append("\treturn rec\n}\n\n");

      out.append("func getNodeExtra(n Node) *nodeExtraFields {\n");
      out.append("\taux, ok := ").append(nsPref).append("Node(n).GetAux()\n");
      out.append("\tif !ok || aux.Extra == nil {\n\t\treturn nil\n\t}\n");
      out.append("\treturn aux.Extra.(*nodeExtraFields)\n}\n\n");

      out.append("func ensureNodeExtra(n Node) *nodeExtraFields {\n");
      out.append("\taux, ok := ").append(nsPref).append("Node(n).GetAux()\n");
      out.append("\tif ok && aux.Extra != nil {\n\t\treturn aux.Extra.(*nodeExtraFields)\n\t}\n");
      out.append("\tex := &nodeExtraFields{cooked: ").append(rt).append("NullString, raw: ").append(rt).append("NullString}\n");
      out.append("\taux.Extra = ex\n");
      out.append('\t').append(nsPref).append("Node(n).SetAux(aux)\n");
      out.append("\treturn ex\n}\n\n");

      // Direct arena field accessors (`token`, `first`, `next`, `previous`, `parent`, `linenoCharno`, `number`, `str`, `originalName`)
      // Note: `nodeBlocks[0]` is permanently nil, so `nodestore.Node(0).Rec()` faults immediately on nil dereference.
      if (valueToken) {
        out.append("func (this_ Node) nodeGet_token() Token {\n");
        out.append("\treturn Token(").append(nsPref).append("Node(this_).Rec().Token + 1)\n}\n");
        out.append("func (this_ Node) nodeSet_token(v Token) {\n");
        out.append('\t').append(nsPref).append("Node(this_).Rec().Token = tokenToByte(v)\n}\n");
      } else {
        out.append("func (this_ Node) nodeGet_token() *Token {\n");
        out.append("\tif t := ").append(nsPref).append("Node(this_).Rec().Token; int(t) < len(enumValues_Token) {\n\t\treturn enumValues_Token[t]\n\t}\n");
        out.append("\treturn nil\n}\n");
        out.append("func (this_ Node) nodeSet_token(v *Token) {\n");
        out.append('\t').append(nsPref).append("Node(this_).Rec().Token = tokenToByte(v)\n}\n");
      }
      out.append("func (this_ Node) TokenOrd_() int32 {\n");
      out.append("\treturn int32(").append(nsPref).append("Node(this_).Rec().Token+1) - 1\n}\n");
      out.append("func (this_ Node) RawToken_() uint8 {\n");
      out.append("\treturn ").append(nsPref).append("Node(this_).Rec().Token\n}\n");
      out.append("func (this_ Node) TokenSwitchOrd_() int32 {\n");
      out.append("\tt := ").append(nsPref).append("Node(this_).Rec().Token\n");
      out.append("\tif t == 0xFF { panic(nodeNullPtrEx) }\n\treturn int32(t)\n}\n");
      out.append("func (this_ Node) FirstNextTok_() (Node, Node, uint8) {\n");
      out.append("\tf, nx, tok := ").append(nsPref).append("Node(this_).FirstNextTok()\n");
      out.append("\treturn Node(f), Node(nx), tok\n}\n");
      out.append("func (this_ Node) HasColdProps_() bool {\n");
      out.append("\treturn ").append(nsPref).append("Node(this_).HasColdProps()\n}\n");

      emitNodePtrAccessor(out, nsPref, "first", "First");
      emitNodePtrAccessor(out, nsPref, "next", "Next");
      emitNodePtrAccessor(out, nsPref, "previous", "Previous");
      emitNodePtrAccessor(out, nsPref, "parent", "Parent");

      out.append("func (this_ Node) nodeGet_linenoCharno() int32 {\n");
      out.append("\treturn ").append(nsPref).append("Node(this_).Rec().LinenoCharno\n}\n");
      out.append("func (this_ Node) nodeSet_linenoCharno(v int32) {\n");
      out.append('\t').append(nsPref).append("Node(this_).Rec().LinenoCharno = v\n}\n");

      // `NumberNode.number` stored zero-alloc in `nodestore`'s noscan float64 table
      out.append("func (this_ Node) nodeGet_number() float64 {\n");
      out.append("\treturn ").append(nsPref).append("Node(this_).NumberVal()\n}\n");
      out.append("func (this_ Node) nodeSet_number(v float64) {\n");
      out.append('\t').append(nsPref).append("Node(this_).SetNumberVal(v)\n}\n");

      // `sourceFile` & `propFlags` stored inline on `NodeData` (`SourceFile16` handle + overflow column + `PropFlags` int32)
      out.append("func (this_ Node) nodeGet_sourceFile() ").append(sourceFileType).append(" {\n");
      out.append("\tif sf := ").append(nsPref).append("SourceFileByHandle(").append(nsPref).append("Node(this_).SourceFile()); sf != nil {\n\t\treturn sf.(").append(sourceFileType).append(")\n\t}\n\treturn nil\n}\n");
      out.append("func (this_ Node) nodeSet_sourceFile(v ").append(sourceFileType).append(") {\n");
      out.append("\tif this_ == 0 { panic(nodeNullPtrEx) }\n");
      out.append("\tif v == nil || v.IsJavaNil_() {\n\t\t").append(nsPref).append("Node(this_).SetSourceFileObj(nil)\n\t} else {\n\t\t")
          .append(nsPref).append("Node(this_).SetSourceFileObj(v)\n\t}\n}\n");

      out.append("func (this_ Node) nodeGet_propFlags() int32 {\n");
      out.append("\treturn int32(").append(nsPref).append("Node(this_).Rec().PropFlags)\n}\n");
      out.append("func (this_ Node) nodeSet_propFlags(v int32) {\n");
      out.append('\t').append(nsPref).append("Node(this_).Rec().PropFlags = uint32(v)\n}\n");

      // `StringNode.str` interned zero-alloc in `nodestore.Intern` via `NodeData.Payload` (0 = NullString)
      out.append("func (this_ Node) nodeGet_str() string {\n");
      out.append("\tp := ").append(nsPref).append("StringID(").append(nsPref).append("Node(this_).Rec().Payload)\n");
      out.append("\tif p == 0 {\n\t\treturn ").append(rt).append("NullString\n\t}\n");
      out.append("\treturn ").append(nsPref).append("String(p)\n}\n");
      out.append("func (this_ Node) nodeSet_str(v string) {\n");
      out.append("\tif v == ").append(rt).append("NullString {\n\t\t")
          .append(nsPref).append("Node(this_).SetStringPayload(0)\n\t} else {\n\t\t")
          .append(nsPref).append("Node(this_).SetStringPayload(").append(nsPref).append("Intern(v))\n\t}\n}\n");

      // `originalName` interned in sparse `nodestore.OrigNameID` paged column (0 = NullString)
      out.append("func (this_ Node) nodeGet_originalName() string {\n");
      out.append("\tif this_ == 0 { panic(nodeNullPtrEx) }\n");
      out.append("\tid := ").append(nsPref).append("Node(this_).OrigNameID()\n");
      out.append("\tif id == 0 {\n\t\treturn ").append(rt).append("NullString\n\t}\n");
      out.append("\treturn ").append(nsPref).append("String(id)\n}\n");
      out.append("func (this_ Node) nodeSet_originalName(v string) {\n");
      out.append("\tif this_ == 0 { panic(nodeNullPtrEx) }\n");
      out.append("\tif v == ").append(rt).append("NullString {\n\t\t").append(nsPref).append("Node(this_).SetOrigNameID(0)\n\t} else {\n\t\t")
          .append(nsPref).append("Node(this_).SetOrigNameID(").append(nsPref).append("Intern(v))\n\t}\n}\n");
      out.append("func (this_ Node) nodeCopy_originalName(other Node) {\n");
      out.append("\tif this_ == 0 || other == 0 {\n\t\tpanic(nodeNullPtrEx)\n\t}\n");
      out.append('\t').append(nsPref).append("Node(this_).SetOrigNameID(").append(nsPref).append("Node(other).OrigNameID())\n}\n");

      // `length` & `changeTime` stored in lock-free `nodestore` columns
      out.append("func (this_ Node) nodeGet_length() int32 {\n");
      out.append("\tif this_ == 0 { panic(nodeNullPtrEx) }\n");
      out.append("\treturn ").append(nsPref).append("Node(this_).Length()\n}\n");
      out.append("func (this_ Node) nodeSet_length(v int32) {\n");
      out.append("\tif this_ == 0 { panic(nodeNullPtrEx) }\n");
      out.append('\t').append(nsPref).append("Node(this_).SetLength(v)\n}\n");
      out.append("func (this_ Node) nodeGet_changeTime() int32 {\n");
      out.append("\tif this_ == 0 { panic(nodeNullPtrEx) }\n");
      out.append("\treturn ").append(nsPref).append("Node(this_).ChangeTime()\n}\n");
      out.append("func (this_ Node) nodeSet_changeTime(v int32) {\n");
      out.append("\tif this_ == 0 { panic(nodeNullPtrEx) }\n");
      out.append('\t').append(nsPref).append("Node(this_).SetChangeTime(v)\n}\n");

      // `jstypeOrColor` & `propListHead` stored in lock-free `nodestore` columns
      out.append("func (this_ Node) nodeGet_jstypeOrColor() ").append(jstypeOrColorType).append(" {\n");
      out.append("\tif this_ == 0 { panic(nodeNullPtrEx) }\n");
      out.append("\treturn ").append(nsPref).append("Node(this_).TypeOrColor()\n}\n");
      out.append("func (this_ Node) nodeSet_jstypeOrColor(v ").append(jstypeOrColorType).append(") {\n");
      out.append("\tif this_ == 0 { panic(nodeNullPtrEx) }\n");
      out.append('\t').append(nsPref).append("Node(this_).SetTypeOrColor(v)\n}\n");

      out.append("func (this_ Node) nodeGet_propListHead() ").append(propListHeadType).append(" {\n");
      out.append("\tif this_ == 0 { panic(nodeNullPtrEx) }\n");
      out.append("\tif p := ").append(nsPref).append("Node(this_).PropList(); p != nil {\n\t\treturn p.(").append(propListHeadType).append(")\n\t}\n\treturn nil\n}\n");
      out.append("func (this_ Node) nodeSet_propListHead(v ").append(propListHeadType).append(") {\n");
      out.append("\tif this_ == 0 { panic(nodeNullPtrEx) }\n");
      out.append("\tif v == nil {\n\t\tif !").append(nsPref).append("Node(this_).HasColdProps() {\n\t\t\treturn\n\t\t}\n\t\t").append(nsPref).append("Node(this_).SetPropList(nil)\n\t} else {\n\t\t")
          .append(nsPref).append("Node(this_).SetPropList(v)\n\t}\n}\n");

      out.append("func (this_ Node) nodeGet_jsDocInfo() ").append(jsDocInfoType).append(" {\n");
      out.append("\tif this_ == 0 { panic(nodeNullPtrEx) }\n");
      out.append("\tif d := ").append(nsPref).append("Node(this_).JSDocObj(); d != nil {\n\t\treturn d.(").append(jsDocInfoType).append(")\n\t}\n\treturn nil\n}\n");
      out.append("func (this_ Node) nodeSet_jsDocInfo(v ").append(jsDocInfoType).append(") {\n");
      out.append("\tif v == nil {\n\t\t").append(nsPref).append("Node(this_).SetJSDocObj(nil)\n\t} else {\n\t\t")
          .append(nsPref).append("Node(this_).SetJSDocObj(v)\n\t}\n}\n");

      out.append("func (this_ Node) nodeGet_typeBeforeCast() ").append(typeBeforeCastType).append(" {\n");
      out.append("\tif this_ == 0 { panic(nodeNullPtrEx) }\n");
      if (typeBeforeCastType.equals("any")) {
        out.append("\treturn ").append(nsPref).append("Node(this_).JSTypeBeforeCastObj()\n}\n");
      } else {
        out.append("\tif t := ").append(nsPref).append("Node(this_).JSTypeBeforeCastObj(); t != nil {\n\t\treturn t.(").append(typeBeforeCastType).append(")\n\t}\n\treturn nil\n}\n");
      }
      out.append("func (this_ Node) nodeSet_typeBeforeCast(v ").append(typeBeforeCastType).append(") {\n");
      out.append("\tif this_ == 0 { panic(nodeNullPtrEx) }\n");
      if (typeBeforeCastType.equals("any")) {
        out.append('\t').append(nsPref).append("Node(this_).SetJSTypeBeforeCastObj(v)\n}\n");
      } else {
        out.append("\tif v == nil || v.IsJavaNil_() {\n\t\t").append(nsPref).append("Node(this_).SetJSTypeBeforeCastObj(nil)\n\t} else {\n\t\t")
            .append(nsPref).append("Node(this_).SetJSTypeBeforeCastObj(v)\n\t}\n}\n");
      }

      emitNodeExtraAccessor(out, "bigint", bigIntType, "nil");
      emitNodeExtraAccessor(out, "cooked", "string", rt + "NullString");
      emitNodeExtraAccessor(out, "raw", "string", rt + "NullString");
      out.append('\n');
    }

    boolean hasClinit = emitStaticFieldsAndClinit(ct, ctPath, te, Map.of(), out, fileCtx);

    Set<String> emittedMethodGoNames = new HashSet<>();
    for (Tree member : ct.getMembers()) {
      if (member instanceof MethodTree mt) {
        TreePath mtPath = new TreePath(ctPath, mt);
        if (isConstructorTree(mt)) {
          emitNodeConstructor(mt, mtPath, te, kindFlag, hasClinit, out, fileCtx);
        } else {
          Element el = resolver.trees.getElement(mtPath);
          if (el instanceof ExecutableElement ee && !ee.getModifiers().contains(Modifier.STATIC)) {
            TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
            if (sig != null) {
              emittedMethodGoNames.add(sig.goName);
              if (kindFlag == 0) {
                if (sig.goName.equals("ClonePropsFrom")
                    && ee.getParameters().size() == 1
                    && resolver.isNodeType(ee.getParameters().get(0).asType())) {
                  String ns = resolver.runtime.nodestoreAlias(fileCtx.imports);
                  String nsPref = ns.isEmpty() ? "" : (ns + ".");
                  String rt = rtPrefix(fileCtx);
                  out.append("// Java method: Node.clonePropsFrom(com.google.javascript.rhino.Node other)\n");
                  out.append("func (this_ Node) ClonePropsFrom(other Node) Node {\n");
                  out.append("\tif !")
                      .append(nsPref)
                      .append("Node(this_).ClonePropsFrom(")
                      .append(nsPref)
                      .append("Node(other)) {\n");
                  out.append("\t\t")
                      .append(rt)
                      .append("CheckStateStr(false, \"Node has existing properties.\")\n");
                  out.append("\t}\n");
                  out.append("\treturn this_\n}\n\n");
                  continue;
                }
                if (sig.goName.equals("GetJSType") && ee.getParameters().isEmpty()) {
                  String ns = resolver.runtime.nodestoreAlias(fileCtx.imports);
                  String nsPref = ns.isEmpty() ? "" : (ns + ".");
                  String retGo = resolver.goType(ee.getReturnType(), fileCtx.imports);
                  out.append("// Java method: Node.getJSType()\n");
                  out.append("func (this_ Node) GetJSType() ").append(retGo).append(" {\n");
                  out.append("\tif this_ == 0 {\n\t\tpanic(nodeNullPtrEx)\n\t}\n");
                  out.append("\tif v := ").append(nsPref).append("Node(this_).JSTypeObj(); v != nil {\n");
                  out.append("\t\treturn v.(").append(retGo).append(")\n\t}\n");
                  out.append("\treturn nil\n}\n\n");
                  continue;
                }
                if (sig.goName.equals("SetJSType")
                    && ee.getParameters().size() == 1
                    && resolver.isNodeType(ee.getReturnType())) {
                  String ns = resolver.runtime.nodestoreAlias(fileCtx.imports);
                  String nsPref = ns.isEmpty() ? "" : (ns + ".");
                  String rt = rtPrefix(fileCtx);
                  String paramGo = resolver.goType(ee.getParameters().get(0).asType(), fileCtx.imports);
                  out.append("// Java method: Node.setJSType(JSType)\n");
                  out.append("func (this_ Node) SetJSType(x ").append(paramGo).append(") Node {\n");
                  out.append("\tif this_ == 0 {\n\t\tpanic(nodeNullPtrEx)\n\t}\n");
                  out.append("\tvar v any\n\tif x != nil && !x.IsJavaNil_() {\n\t\tv = x\n\t}\n");
                  out.append("\tif !").append(nsPref).append("Node(this_).SetJSTypeObj(v) {\n");
                  out.append("\t\t")
                      .append(rt)
                      .append("CheckStateMsg(false, ")
                      .append(rt)
                      .append("BoxNode(this_))\n\t}\n");
                  out.append("\treturn this_\n}\n\n");
                  continue;
                }
                if (sig.goName.equals("GetColor") && ee.getParameters().isEmpty()) {
                  String ns = resolver.runtime.nodestoreAlias(fileCtx.imports);
                  String nsPref = ns.isEmpty() ? "" : (ns + ".");
                  String retGo = resolver.goType(ee.getReturnType(), fileCtx.imports);
                  out.append("// Java method: Node.getColor()\n");
                  out.append("func (this_ Node) GetColor() ").append(retGo).append(" {\n");
                  out.append("\tif this_ == 0 {\n\t\tpanic(nodeNullPtrEx)\n\t}\n");
                  out.append("\tif v := ").append(nsPref).append("Node(this_).ColorObj(); v != nil {\n");
                  out.append("\t\treturn v.(").append(retGo).append(")\n\t}\n");
                  out.append("\treturn nil\n}\n\n");
                  continue;
                }
                if (sig.goName.equals("SetColor")
                    && ee.getParameters().size() == 1
                    && resolver.isNodeType(ee.getReturnType())) {
                  String ns = resolver.runtime.nodestoreAlias(fileCtx.imports);
                  String nsPref = ns.isEmpty() ? "" : (ns + ".");
                  String rt = rtPrefix(fileCtx);
                  String paramGo = resolver.goType(ee.getParameters().get(0).asType(), fileCtx.imports);
                  out.append("// Java method: Node.setColor(Color)\n");
                  out.append("func (this_ Node) SetColor(x ").append(paramGo).append(") Node {\n");
                  out.append("\tif this_ == 0 {\n\t\tpanic(nodeNullPtrEx)\n\t}\n");
                  out.append("\tvar v any\n\tif x != nil && !x.IsJavaNil_() {\n\t\tv = x\n\t}\n");
                  out.append("\tif !").append(nsPref).append("Node(this_).SetColorObj(v) {\n");
                  out.append("\t\t")
                      .append(rt)
                      .append("CheckStateMsg(false, ")
                      .append(rt)
                      .append("BoxNode(this_))\n\t}\n");
                  out.append("\treturn this_\n}\n\n");
                  continue;
                }
                if (sig.goName.equals("SetStringFromStringPool")
                    && ee.getParameters().size() == 2) {
                  String ns = resolver.runtime.nodestoreAlias(fileCtx.imports);
                  String nsPref = ns.isEmpty() ? "" : (ns + ".");
                  String p0Go =
                      resolver.goType(ee.getParameters().get(0).asType(), fileCtx.imports);
                  out.append("// Java method: Node.setStringFromStringPool(LazyInternedStringList, int)\n");
                  out.append("func (this_ Node) SetStringFromStringPool(stringPool ")
                      .append(p0Go)
                      .append(", offset int32) {\n\t")
                      .append(nsPref)
                      .append("Node(this_).SetStringPayload(stringPool.InternID_(offset))\n}\n\n");
                  continue;
                }
                if (sig.goName.equals("SetOriginalNameFromStringPool")
                    && ee.getParameters().size() == 2) {
                  String ns = resolver.runtime.nodestoreAlias(fileCtx.imports);
                  String nsPref = ns.isEmpty() ? "" : (ns + ".");
                  String p0Go =
                      resolver.goType(ee.getParameters().get(0).asType(), fileCtx.imports);
                  out.append("// Java method: Node.setOriginalNameFromStringPool(LazyInternedStringList, int)\n");
                  out.append("func (this_ Node) SetOriginalNameFromStringPool(stringPool ")
                      .append(p0Go)
                      .append(", offset int32) {\n\tid := stringPool.InternID_(offset)\n\tif this_ == 0 {\n\t\tpanic(nodeNullPtrEx)\n\t}\n\t")
                      .append(nsPref)
                      .append("Node(this_).SetOrigNameID(id)\n}\n\n");
                  continue;
                }
              }
            }
          }
          methodEmitter.emitMethod(mt, mtPath, te, hasClinit, null, Map.of(), out, fileCtx);
        }
      }
    }

    if (kindFlag == 0) {
      String ns = resolver.runtime.nodestoreAlias(fileCtx.imports);
      String nsPref = ns.isEmpty() ? "" : (ns + ".");
      String rt = rtPrefix(fileCtx);
      if (!emittedMethodGoNames.contains("Equals")) {
        out.append("func (this_ Node) Equals(other any) bool {\n");
        out.append("\to, ok := other.(Node)\n\treturn ok && this_ == o\n}\n");
      }
      if (!emittedMethodGoNames.contains("HashCode")) {
        out.append("func (this_ Node) HashCode() int32 { return int32(this_) }\n");
      }
      out.append("func (this_ Node) IsJavaNil_() bool { return this_ == 0 }\n");
      out.append("func (this_ Node) RawNodeId_() uint32 { return uint32(this_) }\n");
      out.append("func (this_ Node) JavaClassInfo_() (string, string) {\n");
      out.append("\tif this_ == 0 { ").append(rt).append("ThrowNPE() }\n");
      out.append("\tswitch ").append(nsPref).append("Node(this_).KindFlags() {\n");
      out.append("\tcase 1:\n\t\treturn \"com.google.javascript.rhino.Node$NumberNode\", \"NumberNode\"\n");
      out.append("\tcase 2:\n\t\treturn \"com.google.javascript.rhino.Node$BigIntNode\", \"BigIntNode\"\n");
      out.append("\tcase 3:\n\t\treturn \"com.google.javascript.rhino.Node$StringNode\", \"StringNode\"\n");
      out.append("\tcase 4:\n\t\treturn \"com.google.javascript.rhino.Node$TemplateLiteralSubstringNode\", \"TemplateLiteralSubstringNode\"\n");
      out.append("\tdefault:\n\t\treturn \"com.google.javascript.rhino.Node\", \"Node\"\n\t}\n}\n");
      out.append("func init() {\n");
      TypeElement tokenTe = resolver.elements.getTypeElement("com.google.javascript.rhino.Token");
      if (tokenTe != null) {
        TreePath tokenPath = resolver.trees.getPath(tokenTe);
        if (tokenPath != null
            && tokenPath.getLeaf() instanceof ClassTree tokenCt
            && hasStaticMemberClinit(resolver, tokenCt, tokenPath)) {
          out.append("\tclinit_Token()\n");
        }
      }
      out.append('\t').append(rt).append("RegisterNodeBoxer(func(id_ uint32) any { return Node(id_) })\n");
      out.append('\t').append(rt).append("RegisterTypePredicate(\"com.google.javascript.rhino.Node\", \"com.google.javascript.rhino.Node\", func(v_ any) bool { n, ok := v_.(Node); return ok && n != 0 })\n");
      out.append('\t').append(rt).append("RegisterTypePredicate(\"com.google.javascript.rhino.Node.NumberNode\", \"com.google.javascript.rhino.Node$NumberNode\", func(v_ any) bool { n, ok := v_.(Node); return ok && n != 0 && ").append(nsPref).append("Node(n).KindFlags() == 1 })\n");
      out.append('\t').append(rt).append("RegisterTypePredicate(\"com.google.javascript.rhino.Node.BigIntNode\", \"com.google.javascript.rhino.Node$BigIntNode\", func(v_ any) bool { n, ok := v_.(Node); return ok && n != 0 && ").append(nsPref).append("Node(n).KindFlags() == 2 })\n");
      out.append('\t').append(rt).append("RegisterTypePredicate(\"com.google.javascript.rhino.Node.StringNode\", \"com.google.javascript.rhino.Node$StringNode\", func(v_ any) bool { n, ok := v_.(Node); return ok && n != 0 && ").append(nsPref).append("Node(n).KindFlags() == 3 })\n");
      out.append('\t').append(rt).append("RegisterTypePredicate(\"com.google.javascript.rhino.Node.TemplateLiteralSubstringNode\", \"com.google.javascript.rhino.Node$TemplateLiteralSubstringNode\", func(v_ any) bool { n, ok := v_.(Node); return ok && n != 0 && ").append(nsPref).append("Node(n).KindFlags() == 4 })\n");
      out.append("}\n\n");
    }

    emitNestedMemberTypes(ct, ctPath, emittedTypes, out, fileCtx);
  }

  private void emitNodePtrAccessor(
      StringBuilder out, String nsPref, String javaField, String goField) {
    out.append("func (this_ Node) nodeGet_").append(javaField).append("() Node {\n");
    out.append("\treturn Node(").append(nsPref).append("Node(this_).Rec().").append(goField).append(")\n}\n");
    out.append("func (this_ Node) nodeSet_").append(javaField).append("(v Node) {\n");
    out.append('\t').append(nsPref).append("Node(this_).Rec().").append(goField).append(" = ").append(nsPref).append("Node(v)\n}\n");
  }

  private void emitNodeConstructor(
      MethodTree mt,
      TreePath mtPath,
      TypeElement nodeTe,
      int kindFlag,
      boolean hasClinit,
      StringBuilder out,
      EmitContext fileCtx) {
    Element el = resolver.trees.getElement(mtPath);
    if (!(el instanceof ExecutableElement ctor)) {
      return;
    }
    BlockTree body = mt.getBody();
    List<? extends StatementTree> stmts = body != null ? body.getStatements() : List.of();

    if (isBaseTokenOnlyNodeCtor(nodeTe, ctor)) {
      String ns = resolver.runtime.nodestoreAlias(fileCtx.imports);
      String nsPref = ns.isEmpty() ? "" : (ns + ".");
      String newFn = resolver.getConstructorGoName(ctor);
      String initFn = "initNode" + newFn.substring(3);
      EmitContext ctx = fileCtx.copyForMember(nodeTe, false);
      ctx.currentExec = ctor;

      VariableTree pvt = mt.getParameters().get(0);
      VariableElement pve =
          (VariableElement) resolver.trees.getElement(new TreePath(mtPath, pvt));
      String pName = ctx.registerLocalVar(pve);
      String pType = resolver.goType(pve.asType(), ctx.imports);

      out.append("func ")
          .append(newFn)
          .append('(')
          .append(pName)
          .append(' ')
          .append(pType)
          .append(") Node {\n");
      if (hasClinit) {
        out.append("\tclinit_Node()\n");
      }
      out.append("\ttok_ := tokenToByte(").append(pName).append(")\n");
      out.append("\traw_, rec_ := ").append(nsPref).append("AllocNodeFastRec()\n");
      out.append("\tthis_ := Node(raw_)\n");
      out.append("\trec_.Token = tok_\n");
      out.append("\trec_.LinenoCharno = -1\n");
      out.append("\treturn this_\n}\n\n");

      out.append("func ")
          .append(initFn)
          .append("(this_ Node, ")
          .append(pName)
          .append(' ')
          .append(pType)
          .append(") {\n\t_ = ")
          .append(pName)
          .append("\n\tthis_.nodeSet_token(")
          .append(pName)
          .append(")\n}\n\n");
      fileCtx.syncTempSeqFrom(ctx);
      fileCtx.deferredLocalTypes.addAll(ctx.deferredLocalTypes);
      return;
    }

    if (!stmts.isEmpty()
        && stmts.get(0) instanceof ExpressionStatementTree est
        && est.getExpression() instanceof MethodInvocationTree call0
        && isNodeCtorChainingCall(call0)
        && call0.getArguments().size() == 1) {
      TreePath bodyPath = new TreePath(mtPath, body);
      TreePath call0Path = new TreePath(new TreePath(bodyPath, est), call0);
      Element targetEl = resolver.trees.getElement(call0Path);
      if (targetEl instanceof ExecutableElement targetCtor
          && targetCtor.getEnclosingElement() instanceof TypeElement targetOwner
          && isBaseTokenOnlyNodeCtor(targetOwner, targetCtor)
          && !isTargetedByAnySiblingThisCall(nodeTe, ctor)) {
        String ns = resolver.runtime.nodestoreAlias(fileCtx.imports);
        String nsPref = ns.isEmpty() ? "" : (ns + ".");
        String newFn = resolver.getConstructorGoName(ctor);
        EmitContext ctx = fileCtx.copyForMember(nodeTe, false);
        ctx.currentExec = ctor;

        List<String> sigParams = new ArrayList<>();
        StringBuilder unusedGuards = new StringBuilder();
        for (VariableTree pvt : mt.getParameters()) {
          VariableElement pve =
              (VariableElement) resolver.trees.getElement(new TreePath(mtPath, pvt));
          String pName = ctx.registerLocalVar(pve);
          sigParams.add(pName + " " + resolver.goType(pve.asType(), ctx.imports));
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
        out.append(unusedGuards);
        String tokGo =
            exprEmitter.emitAs(
                call0.getArguments().get(0),
                call0Path,
                resolver.safeErasure(targetCtor.getParameters().get(0).asType()),
                false,
                ctx);
        out.append("\ttok_ := tokenToByte(").append(tokGo).append(")\n");
        out.append("\traw_, rec_ := ").append(nsPref).append("AllocNodeFastRec()\n");
        out.append("\trec_.Token = tok_\n");
        out.append("\trec_.LinenoCharno = -1\n");
        if (kindFlag != 0) {
          out.append("\trec_.KindFlags = ").append(kindFlag).append('\n');
        }
        if (tryEmitFastLeafNodeCtorBody(mt, mtPath, bodyPath, stmts, nodeTe, ctor, kindFlag, nsPref, out, ctx)) {
          out.append("\treturn Node(raw_)\n}\n\n");
          fileCtx.syncTempSeqFrom(ctx);
          fileCtx.deferredLocalTypes.addAll(ctx.deferredLocalTypes);
          return;
        }
        out.append("\tthis_ := Node(raw_)\n");
        int nextStmtIdx = 1;
        if (stmts.size() > 1
            && tryEmitDirectNodeFirstAssign(stmts.get(1), bodyPath, nodeTe, nsPref, out, ctx)) {
          nextStmtIdx = 2;
        }
        if (nextStmtIdx >= stmts.size()) {
          out.append("\t_ = this_\n");
        } else {
          for (int i = nextStmtIdx; i < stmts.size(); i++) {
            StatementTree s = stmts.get(i);
            stmtEmitter.emitStmt(s, new TreePath(bodyPath, s), 1, null, out, ctx);
          }
        }
        out.append("\treturn this_\n}\n\n");
        fileCtx.syncTempSeqFrom(ctx);
        fileCtx.deferredLocalTypes.addAll(ctx.deferredLocalTypes);
        return;
      }
    }

    methodEmitter.emitConstructor(mt, mtPath, nodeTe, hasClinit, Map.of(), out, fileCtx);
  }

  private boolean tryEmitFastLeafNodeCtorBody(
      MethodTree mt,
      TreePath mtPath,
      TreePath bodyPath,
      List<? extends StatementTree> stmts,
      TypeElement nodeTe,
      ExecutableElement ctor,
      int kindFlag,
      String nsPref,
      StringBuilder out,
      EmitContext ctx) {
    if (stmts.size() == 1) {
      // e.g. `private StringNode(Token token) { super(token); }`
      return true;
    }
    String rt = rtPrefix(ctx);
    // 1. Multi-child `Node(Token, Node[, Node[, Node]])`
    if (kindFlag == 0 && isMultiChildDetachedNodeCtor(stmts, ctor)) {
      int arity = ctor.getParameters().size();
      List<String> childVars = new ArrayList<>(arity - 1);
      for (int i = 1; i < arity; i++) {
        VariableElement pve =
            (VariableElement) resolver.trees.getElement(new TreePath(mtPath, mt.getParameters().get(i)));
        childVars.add(ctx.lookupLocalVar(pve));
      }
      if (arity == 2) {
        String c0 = childVars.get(0);
        out.append("\trec_.First = ").append(nsPref).append("Node(").append(c0).append(")\n");
        out.append("\tcrec_ := checkDetachedRec_(").append(c0).append(")\n");
        out.append("\tcrec_.Previous = ").append(nsPref).append("Node(").append(c0).append(")\n");
        out.append("\tcrec_.Parent = raw_\n");
        return true;
      }
      if (arity == 3) {
        String c0 = childVars.get(0);
        String c1 = childVars.get(1);
        out.append("\trec_.First = ").append(nsPref).append("Node(").append(c0).append(")\n");
        out.append("\tlrec_ := checkDetachedRec_(").append(c0).append(")\n");
        out.append("\tlrec_.Next = ").append(nsPref).append("Node(").append(c1).append(")\n");
        out.append("\tlrec_.Previous = ").append(nsPref).append("Node(").append(c1).append(")\n");
        out.append("\tlrec_.Parent = raw_\n");
        out.append("\trrec_ := checkDetachedRec_(").append(c1).append(")\n");
        out.append("\trrec_.Previous = ").append(nsPref).append("Node(").append(c0).append(")\n");
        out.append("\trrec_.Parent = raw_\n");
        return true;
      }
      if (arity == 4) {
        String c0 = childVars.get(0);
        String c1 = childVars.get(1);
        String c2 = childVars.get(2);
        out.append("\trec_.First = ").append(nsPref).append("Node(").append(c0).append(")\n");
        out.append("\tlrec_ := checkDetachedRec_(").append(c0).append(")\n");
        out.append("\tlrec_.Next = ").append(nsPref).append("Node(").append(c1).append(")\n");
        out.append("\tlrec_.Previous = ").append(nsPref).append("Node(").append(c2).append(")\n");
        out.append("\tlrec_.Parent = raw_\n");
        out.append("\tmrec_ := checkDetachedRec_(").append(c1).append(")\n");
        out.append("\tmrec_.Next = ").append(nsPref).append("Node(").append(c2).append(")\n");
        out.append("\tmrec_.Previous = ").append(nsPref).append("Node(").append(c0).append(")\n");
        out.append("\tmrec_.Parent = raw_\n");
        out.append("\trrec_ := checkDetachedRec_(").append(c2).append(")\n");
        out.append("\trrec_.Previous = ").append(nsPref).append("Node(").append(c1).append(")\n");
        out.append("\trrec_.Parent = raw_\n");
        return true;
      }
    }
    // 2. `StringNode(Token, String)` and `StringNode(Token, LazyInternedStringList, int)`
    if (kindFlag == 3
        && stmts.size() == 2
        && stmts.get(1) instanceof ExpressionStatementTree est1
        && est1.getExpression() instanceof MethodInvocationTree call1) {
      String mName = extractSelfNodeCallSimpleName(call1);
      TreePath call1Path = new TreePath(new TreePath(bodyPath, est1), call1);
      if ("setString".equals(mName) && call1.getArguments().size() == 1) {
        String strGo =
            exprEmitter.emitAs(
                call1.getArguments().get(0),
                call1Path,
                resolver.elements.getTypeElement("java.lang.String").asType(),
                false,
                ctx);
        out.append("\ts_ := ").append(strGo).append('\n');
        out.append("\tif s_ == ").append(rt).append("NullString {\n\t\t").append(rt).append("ThrowNPE()\n\t}\n");
        out.append("\trec_.Payload = uint32(").append(nsPref).append("Intern(s_))\n");
        return true;
      }
      if ("setStringFromStringPool".equals(mName) && call1.getArguments().size() == 2) {
        Element targetM = resolver.trees.getElement(call1Path);
        if (targetM instanceof ExecutableElement ee && ee.getParameters().size() == 2) {
          String poolGo =
              exprEmitter.emitAs(
                  call1.getArguments().get(0),
                  call1Path,
                  resolver.safeErasure(ee.getParameters().get(0).asType()),
                  false,
                  ctx);
          String offGo =
              exprEmitter.emitAs(
                  call1.getArguments().get(1),
                  call1Path,
                  resolver.types.getPrimitiveType(TypeKind.INT),
                  false,
                  ctx);
          out.append("\trec_.Payload = uint32(")
              .append(poolGo)
              .append(".InternID_(")
              .append(offGo)
              .append("))\n");
          return true;
        }
      }
    }
    // 3. `NumberNode(double)`
    if (kindFlag == 1
        && stmts.size() == 2
        && stmts.get(1) instanceof ExpressionStatementTree est1
        && est1.getExpression() instanceof MethodInvocationTree call1) {
      String mName = extractSelfNodeCallSimpleName(call1);
      TreePath call1Path = new TreePath(new TreePath(bodyPath, est1), call1);
      if ("setDouble".equals(mName) && call1.getArguments().size() == 1) {
        String numGo =
            exprEmitter.emitAs(
                call1.getArguments().get(0),
                call1Path,
                resolver.types.getPrimitiveType(TypeKind.DOUBLE),
                false,
                ctx);
        out.append("\tnv_ := ").append(numGo).append('\n');
        out.append("\tif !(nv_ == 0 && 1/nv_ > 0) {\n\t\tNode(raw_).SetDouble(nv_)\n\t}\n");
        return true;
      }
    }
    return false;
  }

  private boolean isMultiChildDetachedNodeCtor(
      List<? extends StatementTree> stmts, ExecutableElement ctor) {
    int arity = ctor.getParameters().size();
    int expectedStmts =
        switch (arity) {
          case 2 -> 5;
          case 3 -> 9;
          case 4 -> 13;
          default -> -1;
        };
    if (stmts.size() != expectedStmts) {
      return false;
    }
    for (int i = 1; i < arity; i++) {
      if (!resolver.isNodeType(ctor.getParameters().get(i).asType())) {
        return false;
      }
    }
    if (!(stmts.get(2) instanceof ExpressionStatementTree est2)
        || !(est2.getExpression() instanceof MethodInvocationTree call2)
        || !(call2.getMethodSelect() instanceof MemberSelectTree ms2)
        || !ms2.getIdentifier().contentEquals("checkDetached")) {
      return false;
    }
    return true;
  }

  private static String extractSelfNodeCallSimpleName(MethodInvocationTree call) {
    ExpressionTree sel = call.getMethodSelect();
    if (sel instanceof IdentifierTree id) {
      return id.getName().toString();
    }
    if (sel instanceof MemberSelectTree ms) {
      ExpressionTree q = ms.getExpression();
      while (q instanceof ParenthesizedTree pt) {
        q = pt.getExpression();
      }
      if (q instanceof IdentifierTree qid && qid.getName().contentEquals("this")) {
        return ms.getIdentifier().toString();
      }
    }
    return null;
  }

  private boolean isBaseTokenOnlyNodeCtor(TypeElement owner, ExecutableElement ctor) {
    if (!resolver.isNodeClass(owner)
        || resolver.getNodeKindFlag(owner) != 0
        || ctor.getParameters().size() != 1) {
      return false;
    }
    TypeMirror p0 = resolver.safeErasure(ctor.getParameters().get(0).asType());
    if (!(p0 instanceof DeclaredType dt)
        || !(dt.asElement() instanceof TypeElement pTe)
        || !pTe.getQualifiedName().contentEquals("com.google.javascript.rhino.Token")) {
      return false;
    }
    TreePath ctorPath = resolver.trees.getPath(ctor);
    MethodTree cmt =
        (ctorPath != null && ctorPath.getLeaf() instanceof MethodTree m)
            ? m
            : resolver.trees.getTree(ctor);
    if (cmt == null || cmt.getBody() == null) {
      return false;
    }
    List<? extends StatementTree> cStmts = cmt.getBody().getStatements();
    int assignIdx = 0;
    if (cStmts.size() == 2
        && cStmts.get(0) instanceof ExpressionStatementTree est0
        && est0.getExpression() instanceof MethodInvocationTree call0
        && call0.getArguments().isEmpty()
        && isNodeCtorChainingCall(call0)) {
      assignIdx = 1;
    } else if (cStmts.size() != 1) {
      return false;
    }
    StatementTree st = cStmts.get(assignIdx);
    if (!(st instanceof ExpressionStatementTree est)
        || !(est.getExpression() instanceof AssignmentTree at)
        || at.getKind() != Tree.Kind.ASSIGNMENT) {
      return false;
    }
    ExpressionTree lhs = unwrapParensExpr(at.getVariable());
    ExpressionTree rhs = unwrapParensExpr(at.getExpression());
    if (!(rhs instanceof IdentifierTree rhsId)
        || !rhsId.getName().contentEquals(ctor.getParameters().get(0).getSimpleName())) {
      return false;
    }
    if (lhs instanceof MemberSelectTree ms) {
      ExpressionTree q = unwrapParensExpr(ms.getExpression());
      return q instanceof IdentifierTree qid
          && qid.getName().contentEquals("this")
          && ms.getIdentifier().contentEquals("token");
    }
    return lhs instanceof IdentifierTree lid && lid.getName().contentEquals("token");
  }

  private boolean isTargetedByAnySiblingThisCall(TypeElement owner, ExecutableElement candidate) {
    TreePath ownerPath = resolver.trees.getPath(owner);
    if (ownerPath == null || !(ownerPath.getLeaf() instanceof ClassTree ownerCt)) {
      return true;
    }
    for (Tree m : ownerCt.getMembers()) {
      if (m instanceof MethodTree mt && isConstructorTree(mt) && mt.getBody() != null) {
        List<? extends StatementTree> stmts = mt.getBody().getStatements();
        if (!stmts.isEmpty()
            && stmts.get(0) instanceof ExpressionStatementTree est
            && est.getExpression() instanceof MethodInvocationTree call0
            && call0.getMethodSelect() instanceof IdentifierTree id
            && id.getName().contentEquals("this")) {
          TreePath mtPath = new TreePath(ownerPath, mt);
          TreePath bodyPath = new TreePath(mtPath, mt.getBody());
          TreePath callPath = new TreePath(new TreePath(bodyPath, est), call0);
          if (candidate.equals(resolver.trees.getElement(callPath))) {
            return true;
          }
        }
      }
    }
    return false;
  }

  private static boolean isNodeCtorChainingCall(MethodInvocationTree call) {
    ExpressionTree sel = call.getMethodSelect();
    if (sel instanceof IdentifierTree id) {
      String n = id.getName().toString();
      return n.equals("this") || n.equals("super");
    }
    return sel instanceof MemberSelectTree ms && ms.getIdentifier().contentEquals("super");
  }

  private boolean tryEmitDirectNodeFirstAssign(
      StatementTree st,
      TreePath bodyPath,
      TypeElement nodeTe,
      String nsPref,
      StringBuilder out,
      EmitContext ctx) {
    if (!(st instanceof ExpressionStatementTree est)
        || !(est.getExpression() instanceof AssignmentTree at)
        || at.getKind() != Tree.Kind.ASSIGNMENT) {
      return false;
    }
    ExpressionTree lhs = unwrapParensExpr(at.getVariable());
    boolean isThisFirst = false;
    if (lhs instanceof MemberSelectTree ms) {
      ExpressionTree q = unwrapParensExpr(ms.getExpression());
      isThisFirst =
          q instanceof IdentifierTree qid
              && qid.getName().contentEquals("this")
              && ms.getIdentifier().contentEquals("first");
    } else if (lhs instanceof IdentifierTree lid && lid.getName().contentEquals("first")) {
      TreePath lhsPath = new TreePath(new TreePath(new TreePath(bodyPath, est), at), at.getVariable());
      Element lhsEl = resolver.trees.getElement(lhsPath);
      isThisFirst = lhsEl instanceof VariableElement ve && ve.getKind() == ElementKind.FIELD;
    }
    if (!isThisFirst) {
      return false;
    }
    TreePath atPath = new TreePath(new TreePath(bodyPath, est), at);
    String firstGo = exprEmitter.emitAs(at.getExpression(), atPath, nodeTe.asType(), false, ctx);
    out.append("\trec_.First = ").append(nsPref).append("Node(").append(firstGo).append(")\n");
    return true;
  }

  private void emitNodeExtraAccessor(
      StringBuilder out, String fieldName, String goType, String defaultExpr) {
    out.append("func (this_ Node) nodeGet_").append(fieldName).append("() ").append(goType).append(" {\n");
    out.append("\tif this_ == 0 { panic(nodeNullPtrEx) }\n");
    out.append("\tif ex := getNodeExtra(this_); ex != nil {\n\t\treturn ex.").append(fieldName).append("\n\t}\n");
    out.append("\treturn ").append(defaultExpr).append("\n}\n");
    out.append("func (this_ Node) nodeSet_").append(fieldName).append("(v ").append(goType).append(") {\n");
    out.append("\tif this_ == 0 { panic(nodeNullPtrEx) }\n");
    out.append("\tensureNodeExtra(this_).").append(fieldName).append(" = v\n}\n");
  }

  private String findNodeFieldGoType(
      TypeElement nodeTe, String fieldName, String fallback, EmitContext fileCtx) {
    VariableElement ve = findInstanceFieldByName(nodeTe, fieldName);
    return ve != null ? resolver.goType(ve.asType(), fileCtx.imports) : fallback;
  }

  private String findNodeSubFieldGoType(
      ClassTree nodeCt,
      TreePath nodeCtPath,
      String subSimpleName,
      String fieldName,
      String fallback,
      EmitContext fileCtx) {
    for (Tree m : nodeCt.getMembers()) {
      if (m instanceof ClassTree subCt && subCt.getSimpleName().contentEquals(subSimpleName)) {
        Element el = resolver.trees.getElement(new TreePath(nodeCtPath, subCt));
        if (el instanceof TypeElement subTe) {
          return findNodeFieldGoType(subTe, fieldName, fallback, fileCtx);
        }
      }
    }
    return fallback;
  }

  // =========================================================================
  // Shared Helpers
  // =========================================================================

  private Map<VariableElement, String> buildCapturedFieldMap(TypeElement te) {
    Map<VariableElement, String> map = new LinkedHashMap<>();
    List<VariableElement> caps = resolver.getCapturedVars(te);
    for (int i = 0; i < caps.size(); i++) {
      VariableElement cap = caps.get(i);
      map.put(cap, "cap_" + i + "_" + GoNames.cleanChars(cap.getSimpleName().toString()));
    }
    return map;
  }

  private void emitMarkerMethodsForType(
      TypeElement te, String structName, boolean includeSelfMarker, StringBuilder out) {
    String recvType = resolver.isValueStructClass(te) ? structName : ("*" + structName);
    Set<String> emittedMarkers = new HashSet<>();
    if (includeSelfMarker) {
      String selfMarker = resolver.getInterfaceMarkerMethod(te);
      if (emittedMarkers.add(selfMarker)) {
        out.append("func (").append(recvType).append(") ").append(selfMarker).append("() {}\n");
      }
    }
    for (TypeElement anc : resolver.getAllAncestorTypes(te)) {
      if (anc.getQualifiedName().contentEquals("java.lang.Number")) {
        if (emittedMarkers.add("Implements_j2gort_Number")) {
          out.append("func (").append(recvType).append(") Implements_j2gort_Number() {}\n");
        }
        continue;
      }
      if (resolver.runtime.mapsToAny(anc) || resolver.runtime.isShimException(anc)) {
        continue;
      }
      boolean needsMarker =
          resolver.isSourceType(anc)
              ? (anc.getKind() == ElementKind.INTERFACE || resolver.isSubclassed(anc))
              : (anc.getKind() == ElementKind.INTERFACE || resolver.runtime.isSubclassedExternal(anc));
      if (needsMarker) {
        String marker = resolver.getInterfaceMarkerMethod(anc);
        if (emittedMarkers.add(marker)) {
          out.append("func (").append(recvType).append(") ").append(marker).append("() {}\n");
        }
      }
    }
  }

  private Set<String> collectEmbeddedAncestorMethodGoNames(TypeElement te) {
    Set<String> names = new HashSet<>();
    for (TypeElement anc : resolver.getAllAncestorTypes(te)) {
      if (resolver.runtime.mapsToAny(anc)) {
        continue;
      }
      if (resolver.runtime.isShimException(anc)) {
        names.add("GetMessage");
        names.add("GetCause");
        names.add("Error");
        continue;
      }
      if (!resolver.isSourceType(anc)
          && !resolver.getReferencedExternalTypes().contains(anc)) {
        continue;
      }
      for (Element m : anc.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)
            && !ee.getModifiers().contains(Modifier.PRIVATE)) {
          TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
          if (sig != null && (resolver.isSourceType(anc) || resolver.isNeededExternalSig(sig))) {
            names.add(sig.goName);
          }
        }
        if (resolver.isSourceType(anc)
            && resolver.isSubclassed(anc)
            && m instanceof VariableElement ve
            && ve.getKind() == ElementKind.FIELD
            && !ve.getModifiers().contains(Modifier.STATIC)) {
          String fGo = resolver.getFieldGoName(ve);
          if (resolver.isFieldCrossPkgRead(ve)) names.add("FieldGet_" + fGo);
          if (resolver.isFieldCrossPkgWrite(ve)) names.add("FieldSet_" + fGo);
        }
      }
    }
    return names;
  }

  private String formatInterfaceMethodSpec(
      TypeElement ifaceOrSubclassedTe,
      TypeResolver.MethodSig sig,
      PackagePlanner.ImportTracker imports) {
    TypeMirror siteTm = ifaceOrSubclassedTe != null ? ifaceOrSubclassedTe.asType() : null;
    StringBuilder sb = new StringBuilder();
    sb.append(sig.goName).append('(');
    for (int i = 0; i < sig.canonParamTypes.size(); i++) {
      if (i > 0) sb.append(", ");
      sb.append(resolver.goType(resolver.resolveSigParamType(sig, i, siteTm), imports));
    }
    sb.append(')');
    TypeMirror retTm = resolver.resolveSigReturnType(sig, siteTm);
    if (retTm != null) {
      String ret = resolver.goType(retTm, imports);
      if (!ret.isEmpty()) {
        sb.append(' ').append(ret);
      }
    }
    return sb.toString();
  }

  private boolean isConstructorTree(MethodTree mt) {
    return mt.getReturnType() == null || mt.getName().contentEquals("<init>");
  }

  private VariableElement findInstanceFieldByName(TypeElement te, String simpleName) {
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

  private boolean implementsInterface(TypeElement te, String ifaceQualName) {
    if (te == null) {
      return false;
    }
    if (te.getQualifiedName().contentEquals(ifaceQualName)) {
      return true;
    }
    for (TypeElement anc : resolver.getAllAncestorTypes(te)) {
      if (anc.getQualifiedName().contentEquals(ifaceQualName)) {
        return true;
      }
    }
    return false;
  }

  private record NodePredInfo(
      TypeMirror nodeTm, boolean hasGuavaPred, boolean hasJavaPred, String primarySamGoName) {}

  private NodePredInfo inspectNodePredicateImpl(TypeElement te) {
    if (te == null
        || te.getKind() == ElementKind.INTERFACE
        || resolver.isNodeClass(te)
        || !(te.asType() instanceof DeclaredType subDt)) {
      return null;
    }
    TypeElement gPredTe = resolver.elements.getTypeElement("com.google.common.base.Predicate");
    TypeElement jPredTe = resolver.elements.getTypeElement("java.util.function.Predicate");
    DeclaredType gSup = gPredTe != null ? resolver.findSuperType(subDt, gPredTe) : null;
    DeclaredType jSup = jPredTe != null ? resolver.findSuperType(subDt, jPredTe) : null;
    TypeMirror nodeTm = null;
    if (gSup != null && gSup.getTypeArguments().size() == 1) {
      TypeMirror ta0 = gSup.getTypeArguments().get(0);
      if (ta0 instanceof DeclaredType tadt
          && tadt.asElement() instanceof TypeElement tate
          && resolver.isNodeClass(tate)) {
        nodeTm = ta0;
      }
    }
    if (nodeTm == null && jSup != null && jSup.getTypeArguments().size() == 1) {
      TypeMirror ta0 = jSup.getTypeArguments().get(0);
      if (ta0 instanceof DeclaredType tadt
          && tadt.asElement() instanceof TypeElement tate
          && resolver.isNodeClass(tate)) {
        nodeTm = ta0;
      }
    }
    if (nodeTm == null) {
      return null;
    }
    boolean hasGuava = gSup != null;
    boolean hasJava = jSup != null;
    String primarySamGoName;
    if (hasGuava && gPredTe != null) {
      ExecutableElement sam = resolver.findSamMethod(gPredTe);
      TypeResolver.MethodSig samSig = sam != null ? resolver.getInstanceMethodSig(sam) : null;
      primarySamGoName = samSig != null ? samSig.goName : "Apply";
    } else if (hasJava && jPredTe != null) {
      ExecutableElement sam = resolver.findSamMethod(jPredTe);
      TypeResolver.MethodSig samSig = sam != null ? resolver.getInstanceMethodSig(sam) : null;
      primarySamGoName = samSig != null ? samSig.goName : "Test";
    } else {
      return null;
    }
    return new NodePredInfo(nodeTm, hasGuava, hasJava, primarySamGoName);
  }

  private void emitNodePredicateMethodAndBridge(
      MethodTree mt,
      TreePath mtPath,
      TypeElement typeEl,
      ExecutableElement ee,
      TypeResolver.MethodSig sig,
      NodePredInfo nodePred,
      boolean hasEnumOverrideHook,
      Map<VariableElement, String> capturedFields,
      StringBuilder out,
      EmitContext fileCtx) {
    EmitContext ctx = fileCtx.copyForMember(typeEl, false);
    ctx.currentExec = ee;
    ctx.capturedVarFields.putAll(capturedFields);

    String goTypeName = resolver.planner.getGoTypeName(typeEl);
    List<TypeParameterElement> effTypeParams = resolver.getEffectiveTypeParameters(typeEl);
    String useParams = resolver.formatGoTypeParamUse(effTypeParams);
    String rt = rtPrefix(ctx);

    String recvSpec;
    String nilGuard;
    if (resolver.isValueEnum(typeEl)) {
      recvSpec = "(this_ " + goTypeName + useParams + ")";
      nilGuard = "\tif this_ == 0 {\n\t\t" + rt + "ThrowNPE()\n\t}\n";
    } else if (resolver.isValueStructClass(typeEl)) {
      recvSpec = "(this_ " + goTypeName + useParams + ")";
      String sentinelField = resolver.getValueStructEnumSentinelFieldGoName(typeEl);
      nilGuard =
          resolver.hasValueStructValidFlag(typeEl)
              ? ("\tif !this_.valid_ {\n\t\t" + rt + "ThrowNPE()\n\t}\n")
              : (sentinelField != null
                  ? ("\tif this_." + sentinelField + " == 0 {\n\t\t" + rt + "ThrowNPE()\n\t}\n")
                  : (resolver.isNullableValueStructClass(typeEl)
                      ? ("\tif this_ == (" + goTypeName + useParams + "{}) {\n\t\t" + rt + "ThrowNPE()\n\t}\n")
                      : ""));
    } else if (resolver.isSubclassed(typeEl)) {
      recvSpec = "(this_ *" + goTypeName + "Base" + useParams + ")";
      nilGuard = "\tif this_ == nil {\n\t\t" + rt + "ThrowNPE()\n\t}\n";
    } else {
      recvSpec = "(this_ *" + goTypeName + useParams + ")";
      nilGuard = "\tif this_ == nil {\n\t\t" + rt + "ThrowNPE()\n\t}\n";
    }

    VariableTree pvt = mt.getParameters().get(0);
    VariableElement pve = (VariableElement) resolver.trees.getElement(new TreePath(mtPath, pvt));
    TypeMirror actualTm = resolver.goEffectiveType(pve.asType());
    String userVar = ctx.registerLocalVar(pve);
    String rawNodeArg = ctx.nextTemp("rawNode");
    String nodeGo = resolver.goType(nodePred.nodeTm(), ctx.imports);
    String conv =
        exprEmitter.convert(
            nodeGo + "(" + rawNodeArg + ")", nodePred.nodeTm(), actualTm, false, ctx);

    out.append("// Java method: ")
        .append(typeEl.getSimpleName())
        .append('.')
        .append(ee.getSimpleName())
        .append(formatJavaParamsComment(ee))
        .append('\n');
    out.append("func ")
        .append(recvSpec)
        .append(" TestNode_(")
        .append(rawNodeArg)
        .append(" uint32) bool {\n");
    out.append(nilGuard);

    if (ee.getModifiers().contains(Modifier.SYNCHRONIZED)) {
      if (resolver.hasFastInstanceSyncMutex(typeEl)) {
        out.append("\tthis_._syncMu.Lock()\n");
        out.append("\tdefer this_._syncMu.Unlock()\n");
      } else {
        out.append("\tunlock_ := ").append(rt).append("LockMonitor(this_)\n");
        out.append("\tdefer unlock_()\n");
      }
    }

    if (hasEnumOverrideHook) {
      out.append("\tif this_.fn_TestNode_ != nil {\n\t\treturn this_.fn_TestNode_(this_, ")
          .append(rawNodeArg)
          .append(")\n\t}\n\treturn this_.baseEnum_TestNode_(")
          .append(rawNodeArg)
          .append(")\n}\n\n");
      out.append("func ")
          .append(recvSpec)
          .append(" baseEnum_TestNode_(")
          .append(rawNodeArg)
          .append(" uint32) bool {\n");
    }

    out.append('\t').append(userVar).append(" := ").append(conv).append('\n');
    out.append("\t_ = ").append(userVar).append('\n');

    ctx.expectedReturnType = resolver.types.getPrimitiveType(TypeKind.BOOLEAN);
    ctx.returnWidenedToAny = false;

    boolean isAbstract = ee.getModifiers().contains(Modifier.ABSTRACT);
    boolean isNative = ee.getModifiers().contains(Modifier.NATIVE);
    if (isAbstract || isNative || mt.getBody() == null) {
      out.append("\tpanic(\"abstract or native method invoked\")\n");
    } else {
      out.append(
          stmtEmitter.emitBlockStatements(
              mt.getBody(), new TreePath(mtPath, mt.getBody()), 1, ctx));
      if (!MethodEmitter.endsWithGoTerminatingReturnOrThrow(mt.getBody())) {
        out.append("\tpanic(\"unreachable\")\n");
      }
    }
    out.append("}\n\n");

    TypeMirror canonRetTm = resolver.resolveSigReturnType(sig, typeEl.asType());
    String canonRetGo = canonRetTm != null ? resolver.goType(canonRetTm, ctx.imports) : "bool";
    out.append("func ")
        .append(recvSpec)
        .append(' ')
        .append(sig.goName)
        .append("(arg_ any) ")
        .append(canonRetGo)
        .append(" {\n")
        .append(nilGuard)
        .append("\treturn this_.TestNode_(")
        .append(rt)
        .append("CastNode[uint32](arg_))\n}\n\n");

    fileCtx.syncTempSeqFrom(ctx);
    fileCtx.deferredLocalTypes.addAll(ctx.deferredLocalTypes);
  }

  private String emitEnumConstantNodePredicateFunc(
      MethodTree mt,
      TreePath mtPath,
      TypeElement enumTe,
      NodePredInfo nodePred,
      EmitContext fileCtx) {
    ExecutableElement ee = (ExecutableElement) resolver.trees.getElement(mtPath);
    String enumGoName = resolver.planner.getGoTypeName(enumTe);

    EmitContext ctx = fileCtx.copyForMember(enumTe, false);
    ctx.currentExec = ee;

    VariableTree pvt = mt.getParameters().get(0);
    VariableElement pve = (VariableElement) resolver.trees.getElement(new TreePath(mtPath, pvt));
    TypeMirror actualTm = resolver.goEffectiveType(pve.asType());
    String userVar = ctx.registerLocalVar(pve);
    String rawNodeArg = ctx.nextTemp("rawNode");
    String nodeGo = resolver.goType(nodePred.nodeTm(), ctx.imports);
    String conv =
        exprEmitter.convert(
            nodeGo + "(" + rawNodeArg + ")", nodePred.nodeTm(), actualTm, false, ctx);

    StringBuilder sb = new StringBuilder();
    sb.append("func(this_ *")
        .append(enumGoName)
        .append(", ")
        .append(rawNodeArg)
        .append(" uint32) bool {\n")
        .append('\t')
        .append(userVar)
        .append(" := ")
        .append(conv)
        .append('\n')
        .append("\t_ = ")
        .append(userVar)
        .append('\n');

    ctx.expectedReturnType = resolver.types.getPrimitiveType(TypeKind.BOOLEAN);
    ctx.returnWidenedToAny = false;

    if (mt.getBody() != null) {
      sb.append(
          stmtEmitter.emitBlockStatements(
              mt.getBody(), new TreePath(mtPath, mt.getBody()), 1, ctx));
      if (!MethodEmitter.endsWithGoTerminatingReturnOrThrow(mt.getBody())) {
        sb.append("\tpanic(\"unreachable\")\n");
      }
    } else {
      sb.append("\treturn false\n");
    }
    sb.append('}');

    fileCtx.syncTempSeqFrom(ctx);
    fileCtx.deferredLocalTypes.addAll(ctx.deferredLocalTypes);
    return sb.toString();
  }

  private void emitEnumInheritedNodePredicateDispatcher(
      TypeElement enumTe,
      String goName,
      TypeResolver.MethodSig sig,
      StringBuilder out,
      EmitContext fileCtx) {
    String rt = rtPrefix(fileCtx);
    TypeMirror enumDispRetTm = resolver.resolveSigReturnType(sig, enumTe.asType());
    String retGo = enumDispRetTm != null ? resolver.goType(enumDispRetTm, fileCtx.imports) : "bool";

    out.append("func (this_ *")
        .append(goName)
        .append(") TestNode_(rawNode_ uint32) bool {\n\tif this_ == nil {\n\t\t")
        .append(rt)
        .append(
            "ThrowNPE()\n\t}\n\tif this_.fn_TestNode_ != nil {\n\t\treturn this_.fn_TestNode_(this_, rawNode_)\n\t}\n\treturn this_.baseEnum_TestNode_(rawNode_)\n}\n\n");

    out.append("func (this_ *")
        .append(goName)
        .append(") baseEnum_TestNode_(rawNode_ uint32) bool {\n\t_ = rawNode_\n\tpanic(\"abstract or native method invoked\")\n}\n\n");

    out.append("func (this_ *")
        .append(goName)
        .append(") ")
        .append(sig.goName)
        .append("(arg_ any) ")
        .append(retGo)
        .append(" {\n\tif this_ == nil {\n\t\t")
        .append(rt)
        .append("ThrowNPE()\n\t}\n\treturn this_.TestNode_(")
        .append(rt)
        .append("CastNode[uint32](arg_))\n}\n\n");
  }

  private void emitFallbackOrAbstractNodePredicateMethod(
      TypeElement te,
      String goName,
      String useParams,
      String structWithParams,
      TypeElement superTe,
      NodePredInfo nodePred,
      Set<String> emittedMethodGoNames,
      StringBuilder out,
      EmitContext fileCtx) {
    TypeElement curSup = superTe;
    while (curSup != null && resolver.isSourceType(curSup) && !resolver.runtime.mapsToAny(curSup)) {
      if (inspectNodePredicateImpl(curSup) != null) {
        return;
      }
      curSup = resolver.getSuperclassElement(curSup);
    }

    String rt = rtPrefix(fileCtx);
    String recvSpec;
    String nilGuard;
    if (resolver.isValueStructClass(te)) {
      recvSpec = "(this_ " + structWithParams + ")";
      String sentinelField = resolver.getValueStructEnumSentinelFieldGoName(te);
      nilGuard =
          resolver.hasValueStructValidFlag(te)
              ? ("\tif !this_.valid_ {\n\t\t" + rt + "ThrowNPE()\n\t}\n")
              : (sentinelField != null
                  ? ("\tif this_." + sentinelField + " == 0 {\n\t\t" + rt + "ThrowNPE()\n\t}\n")
                  : (resolver.isNullableValueStructClass(te)
                      ? ("\tif this_ == (" + goName + useParams + "{}) {\n\t\t" + rt + "ThrowNPE()\n\t}\n")
                      : ""));
    } else {
      recvSpec = "(this_ *" + structWithParams + ")";
      nilGuard = "\tif this_ == nil {\n\t\t" + rt + "ThrowNPE()\n\t}\n";
    }

    String samGo = nodePred.primarySamGoName();
    boolean hasConcreteSam =
        emittedMethodGoNames.contains(samGo) || hasConcreteInheritedMethod(te, samGo);
    if (hasConcreteSam) {
      String nodeGo = resolver.goType(nodePred.nodeTm(), fileCtx.imports);
      out.append("func ")
          .append(recvSpec)
          .append(" TestNode_(rawNode_ uint32) bool {\n")
          .append(nilGuard)
          .append("\treturn this_.")
          .append(samGo)
          .append('(')
          .append(rt)
          .append("BoxNode(")
          .append(nodeGo)
          .append("(rawNode_)))\n}\n\n");
    } else {
      out.append("func ")
          .append(recvSpec)
          .append(" TestNode_(rawNode_ uint32) bool {\n")
          .append(nilGuard)
          .append("\t_ = rawNode_\n\tpanic(\"abstract or native method invoked\")\n}\n\n");
    }
  }

  private boolean hasConcreteInheritedMethod(TypeElement te, String methodGoName) {
    TypeElement curSup = resolver.getSuperclassElement(te);
    while (curSup != null && !resolver.runtime.mapsToAny(curSup)) {
      for (Element m : curSup.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)
            && !ee.getModifiers().contains(Modifier.PRIVATE)
            && !ee.getModifiers().contains(Modifier.ABSTRACT)) {
          TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
          if (sig != null && sig.goName.equals(methodGoName)) {
            return true;
          }
        }
      }
      curSup = resolver.getSuperclassElement(curSup);
    }
    for (TypeElement anc : resolver.getAllAncestorTypes(te)) {
      if (anc.getKind() == ElementKind.INTERFACE && !resolver.runtime.mapsToAny(anc)) {
        for (Element m : anc.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && ee.getModifiers().contains(Modifier.DEFAULT)) {
            TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
            if (sig != null
                && sig.goName.equals(methodGoName)
                && (resolver.isSourceType(anc) || resolver.isNeededExternalSig(sig))) {
              return true;
            }
          }
        }
      }
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
