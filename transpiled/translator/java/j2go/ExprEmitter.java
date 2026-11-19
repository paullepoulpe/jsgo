package j2go;

import com.sun.source.tree.ArrayAccessTree;
import com.sun.source.tree.AssignmentTree;
import com.sun.source.tree.BinaryTree;
import com.sun.source.tree.BindingPatternTree;
import com.sun.source.tree.ClassTree;
import com.sun.source.tree.CompoundAssignmentTree;
import com.sun.source.tree.ConditionalExpressionTree;
import com.sun.source.tree.ExpressionStatementTree;
import com.sun.source.tree.ExpressionTree;
import com.sun.source.tree.IdentifierTree;
import com.sun.source.tree.InstanceOfTree;
import com.sun.source.tree.LambdaExpressionTree;
import com.sun.source.tree.LiteralTree;
import com.sun.source.tree.MemberReferenceTree;
import com.sun.source.tree.MemberSelectTree;
import com.sun.source.tree.MethodInvocationTree;
import com.sun.source.tree.NewArrayTree;
import com.sun.source.tree.NewClassTree;
import com.sun.source.tree.ParenthesizedTree;
import com.sun.source.tree.PatternTree;
import com.sun.source.tree.SwitchExpressionTree;
import com.sun.source.tree.Tree;
import com.sun.source.tree.TypeCastTree;
import com.sun.source.tree.UnaryTree;
import com.sun.source.tree.VariableTree;
import com.sun.source.util.TreePath;
import com.sun.source.util.TreePathScanner;
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
import javax.lang.model.element.TypeElement;
import javax.lang.model.element.TypeParameterElement;
import javax.lang.model.element.VariableElement;
import javax.lang.model.type.ArrayType;
import javax.lang.model.type.DeclaredType;
import javax.lang.model.type.ExecutableType;
import javax.lang.model.type.TypeKind;
import javax.lang.model.type.TypeMirror;
import javax.lang.model.type.TypeVariable;
import javax.lang.model.type.WildcardType;

/**
 * Translates Java AST {@link ExpressionTree} nodes into type-checked Go expressions.
 */
public final class ExprEmitter {

  private static final Set<String> NODE_FIELDS =
      Set.of(
          "token",
          "next",
          "previous",
          "first",
          "parent",
          "linenoCharno",
          "length",
          "changeTime",
          "jstypeOrColor",
          "originalName",
          "propListHead",
          "sourceFile",
          "jsDocInfo",
          "typeBeforeCast",
          "propFlags",
          "number",
          "bigint",
          "str",
          "cooked",
          "raw");

  private final TypeResolver resolver;
  private StmtEmitter stmtEmitter;
  private final ThreadLocal<Map<VariableElement, String>> activeEntryCursors =
      ThreadLocal.withInitial(HashMap::new);

  public ExprEmitter(TypeResolver resolver) {
    this.resolver = resolver;
  }

  public void setStmtEmitter(StmtEmitter stmtEmitter) {
    this.stmtEmitter = stmtEmitter;
  }

  void pushMapEntryCursor(VariableElement ve, String curTmp) {
    activeEntryCursors.get().put(ve, curTmp);
  }

  void popMapEntryCursor(VariableElement ve) {
    activeEntryCursors.get().remove(ve);
  }

  public String emitCheckNotNull(
      ExpressionTree tree, String emittedGo, TypeMirror tm, EmitContext ctx) {
    while (tree instanceof ParenthesizedTree pt) {
      tree = pt.getExpression();
    }
    if (tree != null) {
      switch (tree.getKind()) {
        case STRING_LITERAL, NEW_CLASS, NEW_ARRAY, LAMBDA_EXPRESSION, MEMBER_REFERENCE -> {
          return emittedGo;
        }
        case IDENTIFIER -> {
          String nm = ((IdentifierTree) tree).getName().toString();
          if (nm.equals("this") || nm.equals("super")) {
            return emittedGo;
          }
        }
        default -> {}
      }
    }
    String rt = rtPrefix(ctx);
    if (emittedGo.startsWith(rt + "ValueOf")) {
      return emittedGo;
    }
    if (tm != null) {
      if (resolver.isActiveComparableTypeVariable(tm)) {
        String goTv = resolver.goType(tm, ctx.imports);
        String tmp = ctx.nextTemp("tv");
        return "func() "
            + goTv
            + " { "
            + tmp
            + " := "
            + emittedGo
            + "; if _, ok_ := any((*"
            + goTv
            + ")(nil)).(*string); ok_ { if "
            + tmp
            + " == any("
            + rt
            + "NullString).("
            + goTv
            + ") { "
            + rt
            + "ThrowNPE() }; return "
            + tmp
            + " }; var z_ "
            + goTv
            + "; if "
            + tmp
            + " == z_ { "
            + rt
            + "ThrowNPE() }; return "
            + tmp
            + " }()";
      }
      if (resolver.isActiveTypeVariable(tm)) {
        String goTv = resolver.goType(tm, ctx.imports);
        String tmp = ctx.nextTemp("tv");
        return "func() "
            + goTv
            + " { "
            + tmp
            + " := "
            + emittedGo
            + "; _ = "
            + rt
            + "CheckNotNull("
            + rt
            + "ToAny("
            + tmp
            + ")); return "
            + tmp
            + " }()";
      }
      TypeMirror er = resolver.safeErasure(tm);
      if (er.getKind().isPrimitive()) {
        return emittedGo;
      }
      if (resolver.isStringType(er)) {
        return rt + "CheckNotNullString(" + emittedGo + ")";
      }
      if (resolver.isNodeType(er)) {
        return rt + "CheckNotNullNode(" + emittedGo + ")";
      }
      if (resolver.isByteStringType(er)) {
        return rt + "CheckNotNullByteString(" + emittedGo + ")";
      }
      if (resolver.isValueStructType(er)) {
        if (resolver.isNullableValueStructType(er) && !emittedGo.startsWith("New")) {
          return rt + "CheckNotNullVal(" + emittedGo + ")";
        }
        return emittedGo;
      }
      if (resolver.isValueEnumType(er)) {
        return rt + "CheckNotNullEnum(" + emittedGo + ")";
      }
      String goT = resolver.goType(tm, ctx.imports);
      if (goT.startsWith("*")) {
        return rt + "CheckNotNullPtr(" + emittedGo + ")";
      }
      if (goT.startsWith("[]")) {
        return rt + "CheckNotNullSlice(" + emittedGo + ")";
      }
      if (isSourceJavaNilIface(er)) {
        if (emittedGo.equals("nil")) {
          return rt + "CheckNotNullIface[" + goT + "](nil)";
        }
        return rt + "CheckNotNullIface(" + emittedGo + ")";
      }
    }
    return rt + "CheckNotNull(" + emittedGo + ")";
  }

  private boolean isSourceJavaNilIface(TypeMirror tm) {
    if (tm == null) {
      return false;
    }
    TypeMirror er = resolver.safeErasure(tm);
    if (er instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      if (!resolver.isSourceType(te)
          || te.getKind() == ElementKind.ANNOTATION_TYPE
          || resolver.runtime.mapsToAny(te)
          || resolver.isNodeClass(te)) {
        return false;
      }
      return te.getKind() == ElementKind.INTERFACE || resolver.isSubclassed(te);
    }
    return false;
  }

  private boolean isSourceRefAssignableToIface(TypeMirror tm) {
    if (tm == null) {
      return false;
    }
    TypeMirror er = resolver.safeErasure(tm);
    if (er instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      return resolver.isSourceType(te)
          && te.getKind() != ElementKind.ANNOTATION_TYPE
          && !resolver.runtime.mapsToAny(te)
          && !resolver.isNodeClass(te)
          && !resolver.isValueStructClass(te)
          && !resolver.isValueEnum(te);
    }
    return false;
  }

  private boolean isSimpleLocalOrParamIdent(TreePath exprPath, String emittedGo) {
    if (!isSimpleGoIdent(emittedGo)) {
      return false;
    }
    TreePath tp = exprPath;
    while (tp != null && tp.getLeaf() instanceof ParenthesizedTree pt) {
      tp = new TreePath(tp, pt.getExpression());
    }
    if (tp == null || !(tp.getLeaf() instanceof IdentifierTree id)) {
      return false;
    }
    String nm = id.getName().toString();
    if (nm.equals("this") || nm.equals("super")) {
      return false;
    }
    Element el = resolver.trees.getElement(tp);
    if (el instanceof VariableElement ve) {
      return switch (ve.getKind()) {
        case LOCAL_VARIABLE,
            PARAMETER,
            EXCEPTION_PARAMETER,
            RESOURCE_VARIABLE,
            BINDING_VARIABLE ->
            true;
        default -> false;
      };
    }
    return false;
  }

  private static boolean isSimpleGoIdent(String s) {
    if (s == null || s.isEmpty()) {
      return false;
    }
    if (!Character.isJavaIdentifierStart(s.charAt(0))) {
      return false;
    }
    for (int i = 1; i < s.length(); i++) {
      char c = s.charAt(i);
      if (!Character.isLetterOrDigit(c) && c != '_') {
        return false;
      }
    }
    return true;
  }

  private String checkNotNullStringExpr(ExpressionTree tree, String emittedGo, EmitContext ctx) {
    while (tree instanceof ParenthesizedTree pt) {
      tree = pt.getExpression();
    }
    if (tree != null && tree.getKind() == Tree.Kind.STRING_LITERAL) {
      return emittedGo;
    }
    return rtPrefix(ctx) + "CheckNotNullString(" + emittedGo + ")";
  }

  boolean isSubtypeOfQualified(TypeElement te, String targetQualifiedName) {
    if (te == null) {
      return false;
    }
    if (te.getQualifiedName().contentEquals(targetQualifiedName)) {
      return true;
    }
    TypeElement targetTe = resolver.elements.getTypeElement(targetQualifiedName);
    if (targetTe == null) {
      return false;
    }
    return resolver.types.isSubtype(
        resolver.safeErasure(te.asType()), resolver.safeErasure(targetTe.asType()));
  }

  public TypeMirror typeOf(TreePath path) {
    return resolver.trees.getTypeMirror(path);
  }

  public String emit(ExpressionTree expr, TreePath parentPath, EmitContext ctx) {
    if (expr == null) {
      return "";
    }
    TreePath path = new TreePath(parentPath, expr);
    return emitPath(path, ctx);
  }

  public String emitAs(
      ExpressionTree expr,
      TreePath parentPath,
      TypeMirror targetType,
      boolean forceAny,
      EmitContext ctx) {
    if (isDirectCurrentThisExpr(expr, parentPath, ctx)
        && resolver.isSubclassed(ctx.currentType)) {
      TypeElement anc =
          resolver.findReusableSelfAncestor(ctx.currentType, targetType, forceAny);
      if (anc != null) {
        String ancSelf =
            ctx.receiverName + "." + resolver.planner.getGoTypeName(anc) + "Base.self_";
        TypeMirror ancSiteTm =
            ctx.currentType.asType() instanceof DeclaredType curDt
                ? resolver.findSuperType(curDt, anc)
                : anc.asType();
        if (ancSiteTm == null) {
          ancSiteTm = anc.asType();
        }
        return convert(ancSelf, ancSiteTm, targetType, forceAny, ctx);
      }
    }
    TreePath path = new TreePath(parentPath, expr);
    String raw = emitPath(path, ctx);
    TypeMirror srcType = typeOf(path);
    return convert(raw, srcType, targetType, forceAny, ctx);
  }

  private String emitPathOrReusableSelf(
      TreePath p, TypeMirror targetType, boolean forceAny, EmitContext ctx) {
    if (p != null
        && p.getLeaf() instanceof ExpressionTree et
        && isDirectCurrentThisExpr(et, p.getParentPath(), ctx)
        && resolver.isSubclassed(ctx.currentType)) {
      TypeElement anc =
          resolver.findReusableSelfAncestor(ctx.currentType, targetType, forceAny);
      if (anc != null) {
        return ctx.receiverName + "." + resolver.planner.getGoTypeName(anc) + "Base.self_";
      }
    }
    return emitPath(p, ctx);
  }

  private String emitImplicitSelfAsAny(EmitContext ctx) {
    if (ctx.currentType != null
        && !ctx.inStaticContext
        && resolver.isSubclassed(ctx.currentType)) {
      TypeElement anc =
          resolver.findReusableSelfAncestor(
              ctx.currentType, resolver.getObjectElement().asType(), true);
      if (anc != null) {
        return ctx.receiverName + "." + resolver.planner.getGoTypeName(anc) + "Base.self_";
      }
    }
    return ctx.currentType != null
        ? convert(
            ctx.selfValueExpr(),
            ctx.currentType.asType(),
            resolver.getObjectElement().asType(),
            true,
            ctx)
        : ctx.selfValueExpr();
  }

  public String emitPath(TreePath path, EmitContext ctx) {
    Tree leaf = path.getLeaf();
    return switch (leaf.getKind()) {
      case PARENTHESIZED -> "(" + emit(((ParenthesizedTree) leaf).getExpression(), path, ctx) + ")";
      case INT_LITERAL,
          LONG_LITERAL,
          FLOAT_LITERAL,
          DOUBLE_LITERAL,
          BOOLEAN_LITERAL,
          CHAR_LITERAL,
          STRING_LITERAL,
          NULL_LITERAL ->
          emitLiteral((LiteralTree) leaf, typeOf(path), ctx);
      case IDENTIFIER -> emitIdentifier((IdentifierTree) leaf, path, ctx);
      case MEMBER_SELECT -> emitMemberSelect((MemberSelectTree) leaf, path, ctx);
      case ARRAY_ACCESS -> emitArrayAccess((ArrayAccessTree) leaf, path, ctx);
      case METHOD_INVOCATION -> emitMethodInvocation((MethodInvocationTree) leaf, path, ctx);
      case NEW_CLASS -> emitNewClass((NewClassTree) leaf, path, ctx);
      case NEW_ARRAY -> emitNewArray((NewArrayTree) leaf, path, ctx);
      case TYPE_CAST -> emitTypeCast((TypeCastTree) leaf, path, ctx);
      case INSTANCE_OF -> emitInstanceOf((InstanceOfTree) leaf, path, ctx);
      case CONDITIONAL_EXPRESSION ->
          emitConditional((ConditionalExpressionTree) leaf, path, ctx);
      case SWITCH_EXPRESSION -> stmtEmitter.emitSwitchExpression((SwitchExpressionTree) leaf, path, ctx);
      case LAMBDA_EXPRESSION -> emitLambda((LambdaExpressionTree) leaf, path, ctx);
      case MEMBER_REFERENCE -> emitMemberReference((MemberReferenceTree) leaf, path, ctx);
      case ASSIGNMENT -> emitAssignmentExpr((AssignmentTree) leaf, path, ctx);
      case MULTIPLY_ASSIGNMENT,
          DIVIDE_ASSIGNMENT,
          REMAINDER_ASSIGNMENT,
          PLUS_ASSIGNMENT,
          MINUS_ASSIGNMENT,
          LEFT_SHIFT_ASSIGNMENT,
          RIGHT_SHIFT_ASSIGNMENT,
          UNSIGNED_RIGHT_SHIFT_ASSIGNMENT,
          AND_ASSIGNMENT,
          XOR_ASSIGNMENT,
          OR_ASSIGNMENT ->
          emitCompoundAssignmentExpr((CompoundAssignmentTree) leaf, path, ctx);
      case UNARY_PLUS,
          UNARY_MINUS,
          BITWISE_COMPLEMENT,
          LOGICAL_COMPLEMENT,
          PREFIX_INCREMENT,
          PREFIX_DECREMENT,
          POSTFIX_INCREMENT,
          POSTFIX_DECREMENT ->
          emitUnary((UnaryTree) leaf, path, ctx);
      default -> {
        if (leaf instanceof BinaryTree bt) {
          yield emitBinary(bt, path, ctx);
        }
        yield "nil /* unsupported expr: " + leaf.getKind() + " */";
      }
    };
  }

  // --- Universal Type Conversion Helper ---

  public String convert(
      String expr,
      TypeMirror srcType,
      TypeMirror dstType,
      boolean dstForceAny,
      EmitContext ctx) {
    String rt = rtPrefix(ctx);
    if (dstForceAny || resolver.isObjectOrErasedAny(dstType)) {
      if (srcType == null || srcType.getKind() == TypeKind.NULL) {
        return "nil";
      }
      if (resolver.isActiveTypeVariable(srcType)) {
        return rt + "ToAny(" + expr + ")";
      }
      TypeMirror ser = resolver.goEffectiveType(srcType);
      return switch (ser.getKind()) {
        case BOOLEAN -> rt + "ValueOfBoolean(" + expr + ")";
        case BYTE -> rt + "ValueOfByte(" + expr + ")";
        case SHORT -> rt + "ValueOfShort(" + expr + ")";
        case CHAR -> rt + "ValueOfChar(" + expr + ")";
        case INT -> rt + "ValueOfInt(" + expr + ")";
        case LONG -> rt + "ValueOfLong(" + expr + ")";
        case FLOAT -> rt + "ValueOfFloat(" + expr + ")";
        case DOUBLE -> rt + "ValueOfDouble(" + expr + ")";
        default -> {
          if (resolver.isStringType(ser)) {
            yield rt + "BoxString(" + expr + ")";
          }
          if (resolver.isNodeType(ser)) {
            yield rt + "BoxNode(" + expr + ")";
          }
          if (resolver.isByteStringType(ser)) {
            yield rt + "BoxByteString(" + expr + ")";
          }
          if (resolver.isValueEnumType(ser) || resolver.isNullableValueStructType(ser)) {
            yield rt + "BoxNullableVal(" + expr + ")";
          }
          if (resolver.isValueStructType(ser) || resolver.isGoInterfaceType(ser)) {
            yield expr;
          }
          yield emitToAny(rt, resolver.goType(ser, ctx.imports), expr);
        }
      };
    }

    if (dstType == null || dstType.getKind() == TypeKind.VOID || dstType.getKind() == TypeKind.NONE) {
      return expr;
    }

    TypeMirror der = resolver.goEffectiveType(dstType);
    TypeMirror ser = srcType == null ? null : resolver.goEffectiveType(srcType);

    if (ser != null && ser.getKind() == TypeKind.NULL) {
      return resolver.zeroValue(der, ctx.imports);
    }

    if (resolver.isActiveTypeVariable(der)) {
      String dstGo = resolver.goType(der, ctx.imports);
      String srcGo = ser != null ? resolver.goType(ser, ctx.imports) : "any";
      if (dstGo.equals(srcGo)) {
        return expr;
      }
      String anyExpr =
          srcGo.equals("any")
              ? expr
              : convert(expr, srcType, resolver.getObjectElement().asType(), true, ctx);
      String tmp = ctx.nextTemp("tv");
      return "func() "
          + dstGo
          + " { "
          + tmp
          + " := "
          + anyExpr
          + "; if "
          + tmp
          + " == nil { return "
          + resolver.zeroValue(der, ctx.imports)
          + " }; return "
          + rt
          + "Cast["
          + dstGo
          + "]("
          + tmp
          + ") }()";
    }

    // Target is primitive
    if (der.getKind().isPrimitive()) {
      String unboxed = unboxIfNeeded(expr, ser, der.getKind(), ctx);
      TypeKind sk = (ser != null && ser.getKind().isPrimitive()) ? ser.getKind() : der.getKind();
      return convertPrimitive(unboxed, sk, der.getKind(), ctx);
    }

    // Target is java.lang.String
    if (resolver.isStringType(der)) {
      if (ser != null && resolver.isStringType(ser)) {
        return expr;
      }
      return rt + "CastString(" + expr + ")";
    }

    // Target is rhino.Node
    if (resolver.isNodeType(der)) {
      if (ser != null && resolver.isNodeType(ser)) {
        String dstGo = resolver.goType(der, ctx.imports);
        String srcGo = resolver.goType(ser, ctx.imports);
        return dstGo.equals(srcGo) ? expr : (dstGo + "(" + expr + ")");
      }
      String dstGo = resolver.goType(der, ctx.imports);
      return rt + "CastNode[" + dstGo + "](" + expr + ")";
    }

    // Target is boxed primitive (Integer, Boolean, Double, etc.)
    if (der instanceof DeclaredType ddt
        && ddt.asElement() instanceof TypeElement dte
        && resolver.runtime.isBoxedPrimitive(dte)) {
      if (ser != null && ser.getKind().isPrimitive()) {
        return boxPrimitiveToWrapper(expr, ser.getKind(), dte, ctx);
      }
      String dstGo = resolver.goType(der, ctx.imports);
      String srcGo = ser != null ? resolver.goType(ser, ctx.imports) : "any";
      if (dstGo.equals(srcGo)) {
        return expr;
      }
      return emitCast(rt, dstGo, expr);
    }

    String dstGo = resolver.goType(der, ctx.imports);
    String srcGo = ser != null ? resolver.goType(ser, ctx.imports) : "any";

    // Target is Go interface (source interface, subclassed class, external interface, Number, Throwable)
    if (resolver.isGoInterfaceType(der)) {
      boolean boxedScalar = false;
      if (ser != null && ser.getKind().isPrimitive()) {
        expr = convert(expr, ser, resolver.getObjectElement().asType(), true, ctx);
        boxedScalar = true;
      } else if (ser != null && resolver.isStringType(ser)) {
        expr = rt + "BoxString(" + expr + ")";
        boxedScalar = true;
      } else if (ser != null && resolver.isNodeType(ser)) {
        expr = rt + "BoxNode(" + expr + ")";
        boxedScalar = true;
      } else if (ser != null && resolver.isByteStringType(ser)) {
        if (dstGo.equals("any")) {
          return rt + "BoxByteString(" + expr + ")";
        }
        if (der instanceof DeclaredType ddt
            && ddt.asElement() instanceof TypeElement dte
            && !resolver.runtime.mapsToAny(dte)
            && !resolver.runtime.isShimException(dte)
            && resolver.types.isAssignable(
                resolver.safeErasure(ser), resolver.safeErasure(der))) {
          String tmp = ctx.nextTemp("bs");
          return "func() "
              + dstGo
              + " { "
              + tmp
              + " := "
              + expr
              + "; if !"
              + tmp
              + ".IsJavaNil_() { return "
              + tmp
              + " }; return nil }()";
        }
        expr = rt + "BoxByteString(" + expr + ")";
        boxedScalar = true;
      } else if (ser != null && resolver.isNullableValueStructType(ser)) {
        if (dstGo.equals("any")) {
          return rt + "BoxNullableVal(" + expr + ")";
        }
        String tmp = ctx.nextTemp("vs");
        if (der instanceof DeclaredType ddt
            && ddt.asElement() instanceof TypeElement dte
            && !resolver.runtime.mapsToAny(dte)
            && !resolver.runtime.isShimException(dte)
            && resolver.types.isAssignable(
                resolver.safeErasure(ser), resolver.safeErasure(der))) {
          return "func() "
              + dstGo
              + " { "
              + tmp
              + " := "
              + expr
              + "; if "
              + tmp
              + " != ("
              + srcGo
              + "{}) { return "
              + tmp
              + " }; return nil }()";
        }
        expr = rt + "BoxNullableVal(" + expr + ")";
        boxedScalar = true;
      } else if (ser != null && resolver.isValueStructType(ser)) {
        boxedScalar = true;
      }
      if (dstGo.equals(srcGo) && resolver.isGoInterfaceType(ser)) {
        return expr;
      }
      if (dstGo.equals("any")) {
        return boxedScalar ? expr : emitToAny(rt, srcGo, expr);
      }
      if (!boxedScalar) {
        String fastIface = tryEmitToIface(expr, ser, srcGo, der, dstGo);
        if (fastIface != null) {
          return fastIface;
        }
        if (resolver.isActiveTypeVariable(ser)) {
          expr = rt + "ToAny(" + expr + ")";
        }
      }
      return rt + "AsIface[" + dstGo + "](" + expr + ")";
    }

    // Target is concrete pointer `*Foo` or slice `[]T`
    if (dstGo.equals(srcGo)) {
      return expr;
    }
    // Java array covariance: `String[]` -> `Object[]`
    if (der instanceof ArrayType dat && ser instanceof ArrayType sat) {
      String elemGo = resolver.goType(dat.getComponentType(), ctx.imports);
      String tmp = ctx.nextTemp("arr");
      String out = ctx.nextTemp("out");
      String v = ctx.nextTemp("el");
      return "func() "
          + dstGo
          + " { "
          + tmp
          + " := "
          + expr
          + "; if "
          + tmp
          + " == nil { return nil }; "
          + out
          + " := make("
          + dstGo
          + ", len("
          + tmp
          + ")); for i_, "
          + v
          + " := range "
          + tmp
          + " { "
          + out
          + "[i_] = "
          + convert(v, sat.getComponentType(), dat.getComponentType(), false, ctx)
          + " }; return "
          + out
          + " }()";
    }
    return emitCast(rt, dstGo, expr);
  }

  private String emitToAny(String rt, String srcGo, String expr) {
    if (expr.equals("nil")) {
      return "nil";
    }
    if (srcGo != null) {
      if (srcGo.startsWith("*")) {
        return rt + "PtrToAny(" + expr + ")";
      }
      if (srcGo.startsWith("[]")) {
        return rt + "SliceToAny(" + expr + ")";
      }
    }
    return rt + "ToAny(" + expr + ")";
  }

  private String emitCast(String rt, String dstGo, String expr) {
    if (dstGo.startsWith("*")) {
      return rt + "CastPtr[" + dstGo.substring(1) + "](" + expr + ")";
    }
    return rt + "Cast[" + dstGo + "](" + expr + ")";
  }

  private String tryEmitToIface(
      String expr, TypeMirror ser, String srcGo, TypeMirror der, String dstGo) {
    if (expr.equals("nil")) {
      return "nil";
    }
    if (srcGo == null
        || !(ser instanceof DeclaredType sdt)
        || !(sdt.asElement() instanceof TypeElement ste)
        || (!srcGo.startsWith("*") && !resolver.isValueEnum(ste))
        || !resolver.isSourceType(ste)
        || resolver.isNodeClass(ste)) {
      return null;
    }
    if (!(der instanceof DeclaredType ddt) || !(ddt.asElement() instanceof TypeElement dte)) {
      return null;
    }
    if (resolver.runtime.mapsToAny(dte) || resolver.runtime.isShimException(dte)) {
      return null;
    }
    if (!resolver.isSourceType(dte)
        && !dte.getQualifiedName().contentEquals("java.lang.Number")
        && !resolver.getReferencedExternalTypes().contains(dte)) {
      return null;
    }
    if (!resolver.types.isAssignable(ser, der)) {
      return null;
    }
    int bracket = dstGo.indexOf('[');
    int headEnd = bracket >= 0 ? bracket : dstGo.length();
    int dot = dstGo.substring(0, headEnd).indexOf('.');
    String fn =
        dot >= 0
            ? (dstGo.substring(0, dot + 1) + "ToIface_" + dstGo.substring(dot + 1))
            : ("ToIface_" + dstGo);
    return fn + "(" + expr + ")";
  }

  private String boxPrimitiveToWrapper(
      String expr, TypeKind srcKind, TypeElement dstWrapper, EmitContext ctx) {
    String rt = rtPrefix(ctx);
    return switch (dstWrapper.getQualifiedName().toString()) {
      case "java.lang.Boolean" -> rt + "ValueOfBoolean(" + expr + ")";
      case "java.lang.Byte" -> rt + "ValueOfByte(" + convertPrimitive(expr, srcKind, TypeKind.BYTE, ctx) + ")";
      case "java.lang.Short" -> rt + "ValueOfShort(" + convertPrimitive(expr, srcKind, TypeKind.SHORT, ctx) + ")";
      case "java.lang.Character" -> rt + "ValueOfChar(" + convertPrimitive(expr, srcKind, TypeKind.CHAR, ctx) + ")";
      case "java.lang.Integer" -> rt + "ValueOfInt(" + convertPrimitive(expr, srcKind, TypeKind.INT, ctx) + ")";
      case "java.lang.Long" -> rt + "ValueOfLong(" + convertPrimitive(expr, srcKind, TypeKind.LONG, ctx) + ")";
      case "java.lang.Float" -> rt + "ValueOfFloat(" + convertPrimitive(expr, srcKind, TypeKind.FLOAT, ctx) + ")";
      case "java.lang.Double" -> rt + "ValueOfDouble(" + convertPrimitive(expr, srcKind, TypeKind.DOUBLE, ctx) + ")";
      default -> expr;
    };
  }

  private String unboxIfNeeded(
      String expr, TypeMirror srcErased, TypeKind targetPrim, EmitContext ctx) {
    if (srcErased == null || srcErased.getKind().isPrimitive()) {
      return expr;
    }
    String rt = rtPrefix(ctx);
    if (srcErased instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      String qn = te.getQualifiedName().toString();
      return switch (qn) {
        case "java.lang.Boolean" -> rt + "UnboxBoolean(" + expr + ")";
        case "java.lang.Byte" ->
            convertPrimitive(rt + "UnboxByte(" + expr + ")", TypeKind.BYTE, targetPrim, ctx);
        case "java.lang.Short" ->
            convertPrimitive(rt + "UnboxShort(" + expr + ")", TypeKind.SHORT, targetPrim, ctx);
        case "java.lang.Character" ->
            convertPrimitive(rt + "UnboxChar(" + expr + ")", TypeKind.CHAR, targetPrim, ctx);
        case "java.lang.Integer" ->
            convertPrimitive(rt + "UnboxInt(" + expr + ")", TypeKind.INT, targetPrim, ctx);
        case "java.lang.Long" ->
            convertPrimitive(rt + "UnboxLong(" + expr + ")", TypeKind.LONG, targetPrim, ctx);
        case "java.lang.Float" ->
            convertPrimitive(rt + "UnboxFloat(" + expr + ")", TypeKind.FLOAT, targetPrim, ctx);
        case "java.lang.Double" ->
            convertPrimitive(rt + "UnboxDouble(" + expr + ")", TypeKind.DOUBLE, targetPrim, ctx);
        default ->
            switch (targetPrim) {
              case BOOLEAN -> rt + "UnboxBoolean(" + rt + "Cast[*" + rt + "BoxBoolean](" + expr + "))";
              case CHAR -> rt + "UnboxChar(" + rt + "Cast[*" + rt + "BoxChar](" + expr + "))";
              default ->
                  convertPrimitive(
                      rt + "UnboxToDouble(" + expr + ")", TypeKind.DOUBLE, targetPrim, ctx);
            };
      };
    }
    return expr;
  }

  private String convertPrimitive(
      String expr, TypeKind from, TypeKind to, EmitContext ctx) {
    if (from == to) {
      return expr;
    }
    if ((from == TypeKind.DOUBLE || from == TypeKind.FLOAT)
        && (to == TypeKind.INT
            || to == TypeKind.LONG
            || to == TypeKind.SHORT
            || to == TypeKind.BYTE
            || to == TypeKind.CHAR)) {
      String jl = jlPrefix(ctx);
      String dExpr = from == TypeKind.FLOAT ? ("float64(" + expr + ")") : expr;
      if (to == TypeKind.LONG) {
        return jl + "D2L(" + dExpr + ")";
      }
      String i32 = jl + "D2I(" + dExpr + ")";
      return switch (to) {
        case INT -> i32;
        case SHORT -> "int16(" + i32 + ")";
        case BYTE -> "int8(" + i32 + ")";
        case CHAR -> "uint16(" + i32 + ")";
        default -> i32;
      };
    }
    String goTarget =
        switch (to) {
          case BOOLEAN -> "bool";
          case BYTE -> "int8";
          case SHORT -> "int16";
          case CHAR -> "uint16";
          case INT -> "int32";
          case LONG -> "int64";
          case FLOAT -> "float32";
          case DOUBLE -> "float64";
          default -> "";
        };
    return goTarget.isEmpty() ? expr : (goTarget + "(" + expr + ")");
  }

  // --- Literals ---

  private String emitLiteral(LiteralTree lit, TypeMirror tm, EmitContext ctx) {
    Object val = lit.getValue();
    return switch (lit.getKind()) {
      case NULL_LITERAL -> resolver.zeroValue(tm, ctx.imports);
      case BOOLEAN_LITERAL -> Boolean.TRUE.equals(val) ? "true" : "false";
      case INT_LITERAL -> "int32(" + val + ")";
      case LONG_LITERAL -> {
        long lv = ((Number) val).longValue();
        yield lv == Long.MIN_VALUE ? "int64(-9223372036854775808)" : ("int64(" + lv + ")");
      }
      case FLOAT_LITERAL -> {
        float fv = ((Number) val).floatValue();
        if (Float.isNaN(fv)) {
          ctx.imports.std("math");
          yield "float32(math.NaN())";
        }
        if (Float.isInfinite(fv)) {
          ctx.imports.std("math");
          yield fv > 0 ? "float32(math.Inf(1))" : "float32(math.Inf(-1))";
        }
        yield "float32(" + fv + ")";
      }
      case DOUBLE_LITERAL -> {
        double dv = ((Number) val).doubleValue();
        if (Double.isNaN(dv)) {
          ctx.imports.std("math");
          yield "math.NaN()";
        }
        if (Double.isInfinite(dv)) {
          ctx.imports.std("math");
          yield dv > 0 ? "math.Inf(1)" : "math.Inf(-1)";
        }
        yield "float64(" + dv + ")";
      }
      case CHAR_LITERAL -> "uint16(" + (int) ((Character) val) + ")";
      case STRING_LITERAL -> GoNames.goStringLiteral((String) val);
      default -> "nil";
    };
  }

  // --- Identifiers & Member Selects ---

  private String emitIdentifier(IdentifierTree id, TreePath path, EmitContext ctx) {
    String name = id.getName().toString();
    if (name.equals("this")) {
      return ctx.selfValueExpr();
    }
    if (name.equals("super")) {
      return ctx.receiverName;
    }
    Element el = resolver.trees.getElement(path);
    if (el instanceof VariableElement ve) {
      return switch (ve.getKind()) {
        case LOCAL_VARIABLE,
            PARAMETER,
            EXCEPTION_PARAMETER,
            RESOURCE_VARIABLE,
            BINDING_VARIABLE -> {
          String loc = ctx.lookupLocalVar(ve);
          yield loc != null ? loc : GoNames.localIdent(name);
        }
        case ENUM_CONSTANT, FIELD -> emitFieldRead(null, null, ve, typeOf(path), ctx);
        default -> GoNames.localIdent(name);
      };
    }
    return GoNames.localIdent(name);
  }

  private String emitMemberSelect(MemberSelectTree ms, TreePath path, EmitContext ctx) {
    String memberName = ms.getIdentifier().toString();
    TreePath exprPath = new TreePath(path, ms.getExpression());

    if (memberName.equals("class")) {
      TypeMirror targetTm = typeOf(exprPath);
      String qual = targetTm.toString();
      String simple = qual.substring(qual.lastIndexOf('.') + 1);
      if (targetTm instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
        qual = te.getQualifiedName().toString();
        simple = te.getSimpleName().toString();
      }
      return rtPrefix(ctx)
          + "ClassLiteral("
          + GoNames.goStringLiteral(qual)
          + ", "
          + GoNames.goStringLiteral(simple)
          + ")";
    }

    if (memberName.equals("this")) {
      Element targetEl = resolver.trees.getElement(exprPath);
      if (targetEl instanceof TypeElement targetType) {
        return resolveImplicitThisValue(targetType, ctx);
      }
      return ctx.selfValueExpr();
    }

    if (memberName.equals("super")) {
      Element targetEl = resolver.trees.getElement(exprPath);
      if (targetEl instanceof TypeElement targetType) {
        return resolveImplicitThisStruct(targetType, ctx);
      }
      return ctx.receiverName;
    }

    TypeMirror recvType = typeOf(exprPath);
    if (memberName.equals("length") && recvType != null && recvType.getKind() == TypeKind.ARRAY) {
      return "int32(len(" + emitPath(exprPath, ctx) + "))";
    }

    Element el = resolver.trees.getElement(path);
    if (el instanceof VariableElement ve
        && (ve.getKind() == ElementKind.FIELD || ve.getKind() == ElementKind.ENUM_CONSTANT)) {
      boolean isStatic =
          ve.getModifiers().contains(Modifier.STATIC) || ve.getKind() == ElementKind.ENUM_CONSTANT;
      if (!isStatic && isDirectCurrentReceiverExpr(ms.getExpression(), path, ctx)) {
        return emitFieldRead(null, recvType, ve, typeOf(path), ctx);
      }
      String recvGo = isStatic ? null : emitPath(exprPath, ctx);
      return emitFieldRead(recvGo, recvType, ve, typeOf(path), ctx);
    }

    return emitPath(exprPath, ctx) + "." + GoNames.exported(memberName);
  }

  private boolean isDirectCurrentReceiverExpr(
      ExpressionTree expr, TreePath parentPath, EmitContext ctx) {
    if (ctx.currentType == null || ctx.inStaticContext) {
      return false;
    }
    while (expr instanceof ParenthesizedTree pt) {
      parentPath = new TreePath(parentPath, pt);
      expr = pt.getExpression();
    }
    if (expr instanceof IdentifierTree id) {
      String nm = id.getName().toString();
      return nm.equals("this") || nm.equals("super");
    }
    if (expr instanceof MemberSelectTree ms) {
      String mem = ms.getIdentifier().toString();
      if (mem.equals("this") || mem.equals("super")) {
        Element targetEl = resolver.trees.getElement(new TreePath(parentPath, ms.getExpression()));
        return targetEl instanceof TypeElement targetType && targetType.equals(ctx.currentType);
      }
    }
    return false;
  }

  private boolean isDirectCurrentThisExpr(
      ExpressionTree expr, TreePath parentPath, EmitContext ctx) {
    if (ctx.currentType == null || ctx.inStaticContext) {
      return false;
    }
    while (expr instanceof ParenthesizedTree pt) {
      parentPath = new TreePath(parentPath, pt);
      expr = pt.getExpression();
    }
    if (expr instanceof IdentifierTree id) {
      return id.getName().contentEquals("this");
    }
    if (expr instanceof MemberSelectTree ms && ms.getIdentifier().contentEquals("this")) {
      Element targetEl = resolver.trees.getElement(new TreePath(parentPath, ms.getExpression()));
      return targetEl instanceof TypeElement targetType && targetType.equals(ctx.currentType);
    }
    return false;
  }

  /** Resolves implicit `this` -> `this_.outer_...` as a value of `targetOwner`. */
  public String resolveImplicitThisValue(TypeElement targetOwner, EmitContext ctx) {
    if (ctx.currentType == null) {
      return ctx.receiverName;
    }
    TypeElement cur = ctx.currentType;
    String expr = ctx.selfValueExpr();
    int guard = 0;
    while (cur != null && guard++ < 16) {
      if (isSubtypeOf(cur, targetOwner)) {
        return expr;
      }
      if (!resolver.hasOuterInstance(cur)) {
        break;
      }
      TypeElement nextOuter = resolver.getEnclosingInstanceType(cur);
      String baseAccess;
      if (cur.equals(ctx.currentType)) {
        baseAccess = ctx.receiverName + ".outer_";
      } else if (resolver.isSubclassed(cur)) {
        baseAccess = expr + ".Get" + resolver.planner.getGoTypeName(cur) + "Base().outer_";
      } else {
        baseAccess = expr + ".outer_";
      }
      expr = baseAccess;
      cur = nextOuter;
    }
    return expr;
  }

  private String resolveImplicitThisStruct(TypeElement targetOwner, EmitContext ctx) {
    if (ctx.currentType != null
        && !ctx.inStaticContext
        && isSubtypeOf(ctx.currentType, targetOwner)) {
      if (targetOwner.equals(ctx.currentType)) {
        return ctx.receiverName;
      }
      if (resolver.isSubclassed(targetOwner)) {
        return ctx.receiverName + ".Get" + resolver.planner.getGoTypeName(targetOwner) + "Base()";
      }
      return ctx.receiverName;
    }
    String val = resolveImplicitThisValue(targetOwner, ctx);
    if (resolver.isSubclassed(targetOwner) && !val.equals(ctx.receiverName)) {
      return val + ".Get" + resolver.planner.getGoTypeName(targetOwner) + "Base()";
    }
    if (val.equals(ctx.selfValueExpr())) {
      return ctx.receiverName;
    }
    return val;
  }

  private TypeMirror resolveImplicitOwnerSiteType(TypeElement targetOwner, EmitContext ctx) {
    if (ctx.currentType != null) {
      TypeElement cur = ctx.currentType;
      int guard = 0;
      while (cur != null && guard++ < 16) {
        if (isSubtypeOf(cur, targetOwner)) {
          return cur.asType();
        }
        if (!resolver.hasOuterInstance(cur)) {
          break;
        }
        cur = resolver.getEnclosingInstanceType(cur);
      }
    }
    return targetOwner != null ? targetOwner.asType() : null;
  }

  private boolean isSubtypeOf(TypeElement sub, TypeElement sup) {
    if (sub.equals(sup)) {
      return true;
    }
    if (resolver.isNodeClass(sub) && resolver.isNodeClass(sup)) {
      return true;
    }
    return resolver.types.isSubtype(
        resolver.safeErasure(sub.asType()), resolver.safeErasure(sup.asType()));
  }

  private String emitFieldRead(
      String recvExpr,
      TypeMirror recvType,
      VariableElement ve,
      TypeMirror expectedType,
      EmitContext ctx) {
    TypeElement owner = (TypeElement) ve.getEnclosingElement();
    boolean isStatic =
        ve.getModifiers().contains(Modifier.STATIC) || ve.getKind() == ElementKind.ENUM_CONSTANT;

    if (!resolver.isSourceType(owner)) {
      String extRead = emitExternalFieldRead(recvExpr, owner, ve, isStatic, ctx);
      return !isStatic ? narrowGenericFieldRead(extRead, ve.asType(), expectedType, ctx) : extRead;
    }

    if (isStatic) {
      Object constVal = ve.getConstantValue();
      if (constVal != null && ve.getModifiers().contains(Modifier.FINAL)) {
        return emitConstantValue(constVal, ve.asType(), ctx);
      }
      String ownerPkg = resolver.planner.getGoPkg(owner);
      String alias = ctx.imports.ext(resolver.planner.getGoImportPath(ownerPkg), ownerPkg);
      String prefix = alias.isEmpty() ? "" : (alias + ".");
      String getter =
          "StaticGet_"
              + resolver.planner.getGoTypeName(owner)
              + "_"
              + GoNames.exported(ve.getSimpleName().toString());
      return prefix + getter + "()";
    }

    boolean directCurrentStructRecv =
        recvExpr == null
            && ctx.currentType != null
            && !ctx.inStaticContext
            && isSubtypeOf(ctx.currentType, owner);
    String target =
        recvExpr != null
            ? recvExpr
            : (directCurrentStructRecv ? ctx.receiverName : resolveImplicitThisValue(owner, ctx));

    // Special lowering for Node.java's 10 fields and 4 inner classes' fields
    if (resolver.isNodeClass(owner) && NODE_FIELDS.contains(ve.getSimpleName().toString())) {
      if (recvExpr == null) {
        target = ctx.receiverName;
      }
      return target + ".nodeGet_" + ve.getSimpleName() + "()";
    }

    String ownerPkg = resolver.planner.getGoPkg(owner);
    String fieldGo = resolver.getFieldGoName(ve);
    String rawRead;
    if (resolver.isValueEnum(owner)) {
      String recvTarget = isSimpleGoIdent(target) ? target : ("(" + target + ")");
      rawRead = recvTarget + ".FieldGet_" + fieldGo + "()";
    } else if (!ownerPkg.equals(ctx.unitPlan.goPkg)) {
      rawRead = target + ".FieldGet_" + fieldGo + "()";
    } else if (resolver.isSubclassed(owner)) {
      if (recvExpr == null && owner.equals(ctx.currentType)) {
        rawRead = ctx.receiverName + "." + fieldGo;
      } else {
        rawRead = target + ".Get" + resolver.planner.getGoTypeName(owner) + "Base()." + fieldGo;
      }
    } else if (recvExpr == null && owner.equals(ctx.currentType)) {
      rawRead = ctx.receiverName + "." + fieldGo;
    } else {
      rawRead = target + "." + fieldGo;
    }
    TypeMirror actualFieldTm;
    if (resolver.isNativeGenericClass(owner)) {
      if (recvExpr != null) {
        actualFieldTm = resolver.substituteOwnerTypeArgs(ve.asType(), owner, recvType);
      } else if (!owner.equals(ctx.currentType)) {
        actualFieldTm =
            resolver.substituteOwnerTypeArgs(
                ve.asType(), owner, resolveImplicitOwnerSiteType(owner, ctx));
      } else {
        actualFieldTm = ve.asType();
      }
    } else {
      actualFieldTm = ve.asType();
    }
    return narrowGenericFieldRead(rawRead, actualFieldTm, expectedType, ctx);
  }

  private String narrowGenericFieldRead(
      String rawRead, TypeMirror fieldDeclaredTm, TypeMirror expectedType, EmitContext ctx) {
    if (expectedType == null) {
      return rawRead;
    }
    TypeMirror rawTm = resolver.goEffectiveType(fieldDeclaredTm);
    TypeMirror expTm = resolver.goEffectiveType(expectedType);
    if (!resolver.goType(rawTm, ctx.imports).equals(resolver.goType(expTm, ctx.imports))) {
      return convert(rawRead, rawTm, expTm, false, ctx);
    }
    return rawRead;
  }

  private String emitExternalFieldRead(
      String recvExpr,
      TypeElement owner,
      VariableElement ve,
      boolean isStatic,
      EmitContext ctx) {
    String qn = owner.getQualifiedName().toString();
    String fn = ve.getSimpleName().toString();

    if (isStatic) {
      Object constVal = ve.getConstantValue();
      if (constVal != null) {
        return emitConstantValue(constVal, ve.asType(), ctx);
      }
      if (qn.equals("java.lang.System")) {
        if (fn.equals("out") || fn.equals("err")) {
          // Expose via j2gort static field or javaio
          resolver.recordExternalType(owner);
        }
      }
      if (qn.equals("java.lang.Boolean")) {
        if (fn.equals("TRUE")) return rtPrefix(ctx) + "Boolean_TRUE";
        if (fn.equals("FALSE")) return rtPrefix(ctx) + "Boolean_FALSE";
      }
      if (qn.equals("java.lang.Double")) {
        ctx.imports.std("math");
        if (fn.equals("NaN")) return "math.NaN()";
        if (fn.equals("POSITIVE_INFINITY")) return "math.Inf(1)";
        if (fn.equals("NEGATIVE_INFINITY")) return "math.Inf(-1)";
      }
      if (qn.equals("java.lang.Float")) {
        ctx.imports.std("math");
        if (fn.equals("NaN")) return "float32(math.NaN())";
        if (fn.equals("POSITIVE_INFINITY")) return "float32(math.Inf(1))";
        if (fn.equals("NEGATIVE_INFINITY")) return "float32(math.Inf(-1))";
      }
      if (qn.equals("java.math.BigInteger")) {
        // Handled via ext_bindings static field
      }
      String rt = rtPrefix(ctx);
      return rt + resolver.runtime.getExternalStaticFieldGoName(owner, ve);
    }

    String target = recvExpr != null ? recvExpr : ctx.receiverName;
    if (resolver.runtime.isSubclassedExternal(owner)) {
      String extName = resolver.runtime.extGoTypeName(owner);
      return target + ".Get" + extName + "Base().F_" + GoNames.exported(fn);
    }
    return target + ".F_" + GoNames.exported(fn);
  }

  public String emitConstantValue(Object constVal, TypeMirror tm, EmitContext ctx) {
    if (constVal instanceof String s) {
      return GoNames.goStringLiteral(s);
    }
    if (constVal instanceof Character c) {
      return "uint16(" + (int) c + ")";
    }
    if (constVal instanceof Boolean b) {
      return b ? "true" : "false";
    }
    if (constVal instanceof Double d) {
      if (d.isNaN()) {
        ctx.imports.std("math");
        return "math.NaN()";
      }
      if (d.isInfinite()) {
        ctx.imports.std("math");
        return d > 0 ? "math.Inf(1)" : "math.Inf(-1)";
      }
      return "float64(" + d + ")";
    }
    if (constVal instanceof Float f) {
      if (f.isNaN()) {
        ctx.imports.std("math");
        return "float32(math.NaN())";
      }
      if (f.isInfinite()) {
        ctx.imports.std("math");
        return f > 0 ? "float32(math.Inf(1))" : "float32(math.Inf(-1))";
      }
      return "float32(" + f + ")";
    }
    if (constVal instanceof Long l) {
      return l == Long.MIN_VALUE ? "int64(-9223372036854775808)" : ("int64(" + l + ")");
    }
    String goT = resolver.goType(tm, ctx.imports);
    return goT + "(" + constVal + ")";
  }

  private String emitArrayAccess(ArrayAccessTree aa, TreePath path, EmitContext ctx) {
    TreePath arrPath = new TreePath(path, aa.getExpression());
    String arr = emitPath(arrPath, ctx);
    String idx = emitAs(aa.getIndex(), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
    String raw = arr + "[" + idx + "]";
    TypeMirror arrTm = resolver.goEffectiveType(typeOf(arrPath));
    TypeMirror expTm = typeOf(path);
    if (arrTm instanceof ArrayType at && expTm != null) {
      return narrowGenericFieldRead(raw, at.getComponentType(), expTm, ctx);
    }
    return raw;
  }

  // --- Assignments & Increment/Decrement ---

  public String emitAssignmentStatement(AssignmentTree at, TreePath path, EmitContext ctx) {
    TreePath lhsPath = new TreePath(path, at.getVariable());
    TypeMirror lhsType = typeOf(lhsPath);
    String rhs = emitAs(at.getExpression(), path, lhsType, false, ctx);
    return writeLhs(lhsPath, rhs, ctx);
  }

  private String emitAssignmentExpr(AssignmentTree at, TreePath path, EmitContext ctx) {
    TreePath lhsPath = new TreePath(path, at.getVariable());
    TypeMirror lhsType = typeOf(lhsPath);
    String goType = resolver.goType(lhsType, ctx.imports);
    String valTmp = ctx.nextTemp("asgn");
    String rhs = emitAs(at.getExpression(), path, lhsType, false, ctx);
    String writeStmt = writeLhs(lhsPath, valTmp, ctx);
    return "func() "
        + goType
        + " { var "
        + valTmp
        + " "
        + goType
        + " = "
        + rhs
        + "; "
        + writeStmt
        + "; return "
        + valTmp
        + " }()";
  }

  public String emitCompoundAssignmentStatement(
      CompoundAssignmentTree cat, TreePath path, EmitContext ctx) {
    return buildCompoundAssignCore(cat, path, false, ctx);
  }

  private String emitCompoundAssignmentExpr(
      CompoundAssignmentTree cat, TreePath path, EmitContext ctx) {
    return buildCompoundAssignCore(cat, path, true, ctx);
  }

  private String buildCompoundAssignCore(
      CompoundAssignmentTree cat, TreePath path, boolean asExpr, EmitContext ctx) {
    TreePath lhsPath = new TreePath(path, cat.getVariable());
    TreePath rhsPath = new TreePath(path, cat.getExpression());
    TypeMirror lhsType = typeOf(lhsPath);
    TypeMirror rhsType = typeOf(rhsPath);
    String lhsGoType = resolver.goType(lhsType, ctx.imports);

    LhsHandle h = prepareLhsHandle(lhsPath, ctx);
    String curVal = h.readExpr;

    String computed;
    if (cat.getKind() == Tree.Kind.PLUS_ASSIGNMENT && resolver.isStringType(lhsType)) {
      List<TreePath> rhsParts = new ArrayList<>();
      collectStringConcatOperandPaths(rhsPath, rhsParts);
      StringBuilder catExpr =
          new StringBuilder("(").append(rtPrefix(ctx)).append("StrValString(").append(curVal).append(')');
      for (TreePath rp : rhsParts) {
        catExpr.append(" + ").append(emitStringConcatOperand(rp, ctx));
      }
      catExpr.append(')');
      computed = catExpr.toString();
    } else {
      String rhsRaw = emitPath(rhsPath, ctx);
      Tree.Kind binKind = compoundToBinaryKind(cat.getKind());
      computed = computeNumericOrBoolBinary(binKind, curVal, lhsType, rhsRaw, rhsType, lhsType, ctx);
      computed = convert(computed, promotedBinaryType(binKind, lhsType, rhsType), lhsType, false, ctx);
    }

    if (!asExpr) {
      StringBuilder sb = new StringBuilder();
      if (!h.setupStmt.isEmpty()) {
        sb.append(h.setupStmt).append("; ");
      }
      sb.append(h.writeWith(computed));
      return sb.toString();
    }

    String vTmp = ctx.nextTemp("cv");
    StringBuilder iife = new StringBuilder("func() ").append(lhsGoType).append(" { ");
    if (!h.setupStmt.isEmpty()) {
      iife.append(h.setupStmt).append("; ");
    }
    iife.append(vTmp)
        .append(" := ")
        .append(computed)
        .append("; ")
        .append(h.writeWith(vTmp))
        .append("; return ")
        .append(vTmp)
        .append(" }()");
    return iife.toString();
  }

  public String emitUnaryStatement(UnaryTree ut, TreePath path, EmitContext ctx) {
    Tree.Kind k = ut.getKind();
    if (k == Tree.Kind.PREFIX_INCREMENT
        || k == Tree.Kind.POSTFIX_INCREMENT
        || k == Tree.Kind.PREFIX_DECREMENT
        || k == Tree.Kind.POSTFIX_DECREMENT) {
      return buildIncDecCore(ut, path, false, ctx);
    }
    return "_ = " + emitUnary(ut, path, ctx);
  }

  private String emitUnary(UnaryTree ut, TreePath path, EmitContext ctx) {
    Tree.Kind k = ut.getKind();
    if (k == Tree.Kind.PREFIX_INCREMENT
        || k == Tree.Kind.POSTFIX_INCREMENT
        || k == Tree.Kind.PREFIX_DECREMENT
        || k == Tree.Kind.POSTFIX_DECREMENT) {
      return buildIncDecCore(ut, path, true, ctx);
    }

    TreePath subPath = new TreePath(path, ut.getExpression());
    TypeMirror subType = typeOf(subPath);
    TypeMirror resType = typeOf(path);

    if (k == Tree.Kind.LOGICAL_COMPLEMENT) {
      String sub = emitAs(ut.getExpression(), path, resolver.types.getPrimitiveType(TypeKind.BOOLEAN), false, ctx);
      return "!(" + sub + ")";
    }

    // Unary numeric promotion (+, -, ~)
    TypeKind rk = resolver.safeErasure(resType).getKind();
    String promoted = emitAs(ut.getExpression(), path, resType, false, ctx);
    return switch (k) {
      case UNARY_PLUS -> "(" + promoted + ")";
      case UNARY_MINUS -> "-(" + promoted + ")";
      case BITWISE_COMPLEMENT -> "^(" + promoted + ")";
      default -> promoted;
    };
  }

  private String buildIncDecCore(UnaryTree ut, TreePath path, boolean asExpr, EmitContext ctx) {
    TreePath lhsPath = new TreePath(path, ut.getExpression());
    TypeMirror lhsType = typeOf(lhsPath);
    String lhsGoType = resolver.goType(lhsType, ctx.imports);
    boolean isInc =
        ut.getKind() == Tree.Kind.PREFIX_INCREMENT || ut.getKind() == Tree.Kind.POSTFIX_INCREMENT;
    boolean isPost =
        ut.getKind() == Tree.Kind.POSTFIX_INCREMENT || ut.getKind() == Tree.Kind.POSTFIX_DECREMENT;

    LhsHandle h = prepareLhsHandle(lhsPath, ctx);
    String unboxedCur = unboxIfNeeded(h.readExpr, resolver.safeErasure(lhsType), primKindForInc(lhsType), ctx);
    String step = unboxedCur + (isInc ? " + 1" : " - 1");
    String nextVal = convert(step, resolver.types.getPrimitiveType(primKindForInc(lhsType)), lhsType, false, ctx);

    if (!asExpr) {
      StringBuilder sb = new StringBuilder();
      if (!h.setupStmt.isEmpty()) {
        sb.append(h.setupStmt).append("; ");
      }
      sb.append(h.writeWith(nextVal));
      return sb.toString();
    }

    String oldTmp = ctx.nextTemp("old");
    String newTmp = ctx.nextTemp("nw");
    String nextFromOld =
        convert(
            unboxIfNeeded(oldTmp, resolver.safeErasure(lhsType), primKindForInc(lhsType), ctx)
                + (isInc ? " + 1" : " - 1"),
            resolver.types.getPrimitiveType(primKindForInc(lhsType)),
            lhsType,
            false,
            ctx);
    StringBuilder iife = new StringBuilder("func() ").append(lhsGoType).append(" { ");
    if (!h.setupStmt.isEmpty()) {
      iife.append(h.setupStmt).append("; ");
    }
    iife.append(oldTmp)
        .append(" := ")
        .append(h.readExpr)
        .append("; ")
        .append(newTmp)
        .append(" := ")
        .append(nextFromOld)
        .append("; ")
        .append(h.writeWith(newTmp))
        .append("; return ")
        .append(isPost ? oldTmp : newTmp)
        .append(" }()");
    return iife.toString();
  }

  private TypeKind primKindForInc(TypeMirror tm) {
    TypeMirror er = resolver.safeErasure(tm);
    if (er.getKind().isPrimitive()) {
      return er.getKind();
    }
    if (er instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      return switch (te.getQualifiedName().toString()) {
        case "java.lang.Long" -> TypeKind.LONG;
        case "java.lang.Double" -> TypeKind.DOUBLE;
        case "java.lang.Float" -> TypeKind.FLOAT;
        case "java.lang.Character" -> TypeKind.CHAR;
        case "java.lang.Short" -> TypeKind.SHORT;
        case "java.lang.Byte" -> TypeKind.BYTE;
        default -> TypeKind.INT;
      };
    }
    return TypeKind.INT;
  }

  private record LhsHandle(String setupStmt, String readExpr, java.util.function.Function<String, String> writer) {
    String writeWith(String valExpr) {
      return writer.apply(valExpr);
    }
  }

  private LhsHandle prepareLhsHandle(TreePath lhsPath, EmitContext ctx) {
    Tree lhs = lhsPath.getLeaf();
    while (lhs instanceof ParenthesizedTree pt) {
      lhsPath = new TreePath(lhsPath, pt.getExpression());
      lhs = lhsPath.getLeaf();
    }
    if (lhs instanceof ArrayAccessTree aa) {
      String aTmp = ctx.nextTemp("a");
      String iTmp = ctx.nextTemp("i");
      String aExpr = emit(aa.getExpression(), lhsPath, ctx);
      String iExpr =
          emitAs(aa.getIndex(), lhsPath, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
      String setup = aTmp + " := " + aExpr + "; " + iTmp + " := " + iExpr;
      String access = aTmp + "[" + iTmp + "]";
      return new LhsHandle(setup, access, v -> access + " = " + v);
    }

    Element el = resolver.trees.getElement(lhsPath);
    if (el instanceof VariableElement ve) {
      if (ve.getKind() != ElementKind.FIELD && ve.getKind() != ElementKind.ENUM_CONSTANT) {
        String loc = ctx.lookupLocalVar(ve);
        String name = loc != null ? loc : GoNames.localIdent(ve.getSimpleName().toString());
        return new LhsHandle("", name, v -> name + " = " + v);
      }

      TypeElement owner = (TypeElement) ve.getEnclosingElement();
      boolean isStatic = ve.getModifiers().contains(Modifier.STATIC);
      if (isStatic) {
        if (!resolver.isSourceType(owner)) {
          String extVar = rtPrefix(ctx) + resolver.runtime.getExternalStaticFieldGoName(owner, ve);
          return new LhsHandle("", extVar, v -> extVar + " = " + v);
        }
        String ownerPkg = resolver.planner.getGoPkg(owner);
        String alias = ctx.imports.ext(resolver.planner.getGoImportPath(ownerPkg), ownerPkg);
        String prefix = alias.isEmpty() ? "" : (alias + ".");
        String getter =
            prefix
                + "StaticGet_"
                + resolver.planner.getGoTypeName(owner)
                + "_"
                + GoNames.exported(ve.getSimpleName().toString())
                + "()";
        String setterName =
            prefix
                + "StaticSet_"
                + resolver.planner.getGoTypeName(owner)
                + "_"
                + GoNames.exported(ve.getSimpleName().toString());
        return new LhsHandle("", getter, v -> setterName + "(" + v + ")");
      }

      // Instance field
      String setup = "";
      String recvVal;
      boolean directCurrentOwnerField = false;
      if (lhs instanceof MemberSelectTree ms) {
        if (isDirectCurrentReceiverExpr(ms.getExpression(), lhsPath, ctx)) {
          recvVal = ctx.receiverName;
          directCurrentOwnerField = owner.equals(ctx.currentType);
        } else {
          String rTmp = ctx.nextTemp("r");
          String rRaw = emit(ms.getExpression(), lhsPath, ctx);
          setup = rTmp + " := " + rRaw;
          recvVal = rTmp;
        }
      } else if (ctx.currentType != null
          && !ctx.inStaticContext
          && isSubtypeOf(ctx.currentType, owner)) {
        recvVal = ctx.receiverName;
        directCurrentOwnerField = owner.equals(ctx.currentType);
      } else {
        recvVal =
            resolver.isNodeClass(owner)
                ? ctx.receiverName
                : resolveImplicitThisValue(owner, ctx);
      }

      TypeMirror lhsTm = typeOf(lhsPath);
      TypeMirror actualFieldTm;
      if (resolver.isNativeGenericClass(owner)) {
        if (lhs instanceof MemberSelectTree ms
            && !isDirectCurrentReceiverExpr(ms.getExpression(), lhsPath, ctx)) {
          actualFieldTm =
              resolver.substituteOwnerTypeArgs(
                  ve.asType(), owner, typeOf(new TreePath(lhsPath, ms.getExpression())));
        } else if (!owner.equals(ctx.currentType)) {
          actualFieldTm =
              resolver.substituteOwnerTypeArgs(
                  ve.asType(), owner, resolveImplicitOwnerSiteType(owner, ctx));
        } else {
          actualFieldTm = ve.asType();
        }
      } else {
        actualFieldTm = ve.asType();
      }
      TypeMirror fieldErasedTm = resolver.goEffectiveType(actualFieldTm);
      boolean forceFieldAny = resolver.isObjectOrErasedAny(fieldErasedTm);
      boolean genericFieldMismatch =
          lhsTm != null
              && !resolver.goType(fieldErasedTm, ctx.imports).equals(resolver.goType(lhsTm, ctx.imports));

      if (!resolver.isSourceType(owner)) {
        String fieldAccess =
            resolver.runtime.isSubclassedExternal(owner)
                ? (recvVal
                    + ".Get"
                    + resolver.runtime.extGoTypeName(owner)
                    + "Base().F_"
                    + GoNames.exported(ve.getSimpleName().toString()))
                : (recvVal + ".F_" + GoNames.exported(ve.getSimpleName().toString()));
        if (genericFieldMismatch) {
          return new LhsHandle(
              setup,
              convert(fieldAccess, fieldErasedTm, lhsTm, false, ctx),
              v -> fieldAccess + " = " + convert(v, lhsTm, fieldErasedTm, forceFieldAny, ctx));
        }
        return new LhsHandle(setup, fieldAccess, v -> fieldAccess + " = " + v);
      }

      if (resolver.isNodeClass(owner) && NODE_FIELDS.contains(ve.getSimpleName().toString())) {
        String fn = ve.getSimpleName().toString();
        return new LhsHandle(
            setup,
            recvVal + ".nodeGet_" + fn + "()",
            v -> recvVal + ".nodeSet_" + fn + "(" + v + ")");
      }

      String ownerPkg = resolver.planner.getGoPkg(owner);
      String fieldGo = resolver.getFieldGoName(ve);
      if (!ownerPkg.equals(ctx.unitPlan.goPkg)) {
        String rawGet = recvVal + ".FieldGet_" + fieldGo + "()";
        if (genericFieldMismatch) {
          return new LhsHandle(
              setup,
              convert(rawGet, fieldErasedTm, lhsTm, false, ctx),
              v ->
                  recvVal
                      + ".FieldSet_"
                      + fieldGo
                      + "("
                      + convert(v, lhsTm, fieldErasedTm, forceFieldAny, ctx)
                      + ")");
        }
        return new LhsHandle(
            setup,
            rawGet,
            v -> recvVal + ".FieldSet_" + fieldGo + "(" + v + ")");
      }
      String directField;
      if (resolver.isSubclassed(owner)) {
        directField =
            directCurrentOwnerField
                ? (ctx.receiverName + "." + fieldGo)
                : (recvVal + ".Get" + resolver.planner.getGoTypeName(owner) + "Base()." + fieldGo);
      } else {
        directField =
            directCurrentOwnerField
                ? (ctx.receiverName + "." + fieldGo)
                : (recvVal + "." + fieldGo);
      }
      if (genericFieldMismatch) {
        return new LhsHandle(
            setup,
            convert(directField, fieldErasedTm, lhsTm, false, ctx),
            v -> directField + " = " + convert(v, lhsTm, fieldErasedTm, forceFieldAny, ctx));
      }
      return new LhsHandle(setup, directField, v -> directField + " = " + v);
    }

    String raw = emitPath(lhsPath, ctx);
    return new LhsHandle("", raw, v -> raw + " = " + v);
  }

  private String writeLhs(TreePath lhsPath, String rhsGoExpr, EmitContext ctx) {
    LhsHandle h = prepareLhsHandle(lhsPath, ctx);
    if (h.setupStmt.isEmpty()) {
      return h.writeWith(rhsGoExpr);
    }
    return "{ " + h.setupStmt + "; " + h.writeWith(rhsGoExpr) + " }";
  }

  public boolean isPureSideEffectFreeAssignLhs(TreePath lhsPath, EmitContext ctx) {
    Tree lhs = lhsPath.getLeaf();
    while (lhs instanceof ParenthesizedTree pt) {
      lhsPath = new TreePath(lhsPath, pt.getExpression());
      lhs = lhsPath.getLeaf();
    }
    Element el = resolver.trees.getElement(lhsPath);
    if (!(el instanceof VariableElement ve)) {
      return false;
    }
    if (ve.getKind() != ElementKind.FIELD && ve.getKind() != ElementKind.ENUM_CONSTANT) {
      return lhs instanceof IdentifierTree;
    }
    if (ve.getModifiers().contains(Modifier.STATIC)) {
      if (lhs instanceof IdentifierTree) {
        return true;
      }
      if (lhs instanceof MemberSelectTree ms) {
        Element qEl = resolver.trees.getElement(new TreePath(lhsPath, ms.getExpression()));
        return qEl instanceof TypeElement;
      }
      return false;
    }
    if (lhs instanceof IdentifierTree) {
      return ctx.currentType != null && !ctx.inStaticContext;
    }
    if (lhs instanceof MemberSelectTree ms) {
      return isDirectCurrentReceiverExpr(ms.getExpression(), lhsPath, ctx);
    }
    return false;
  }

  public String emitPureLhsWrite(TreePath lhsPath, String rhsGoExpr, EmitContext ctx) {
    return writeLhs(lhsPath, rhsGoExpr, ctx);
  }

  // --- Binary Expressions ---

  private String emitBinary(BinaryTree bt, TreePath path, EmitContext ctx) {
    Tree.Kind kind = bt.getKind();
    TypeMirror resType = typeOf(path);

    // String concatenation `+` -> native Go `+` (`runtime.concatstrings`, zero `[]any` slice / `any` boxing!)
    if (kind == Tree.Kind.PLUS && resolver.isStringType(resType)) {
      List<TreePath> partPaths = new ArrayList<>();
      collectStringConcatOperandPaths(path, partPaths);
      List<String> parts = new ArrayList<>(partPaths.size());
      for (TreePath pp : partPaths) {
        parts.add(emitStringConcatOperand(pp, ctx));
      }
      return "(" + String.join(" + ", parts) + ")";
    }

    TreePath lp = new TreePath(path, bt.getLeftOperand());
    TreePath rp = new TreePath(path, bt.getRightOperand());
    TypeMirror lt = typeOf(lp);
    TypeMirror rt = typeOf(rp);

    if (kind == Tree.Kind.CONDITIONAL_AND || kind == Tree.Kind.CONDITIONAL_OR) {
      TypeMirror boolTm = resolver.types.getPrimitiveType(TypeKind.BOOLEAN);
      String l = emitAs(bt.getLeftOperand(), path, boolTm, false, ctx);
      String r = emitAs(bt.getRightOperand(), path, boolTm, false, ctx);
      String op = kind == Tree.Kind.CONDITIONAL_AND ? " && " : " || ";
      return "(" + l + op + r + ")";
    }

    if (kind == Tree.Kind.EQUAL_TO || kind == Tree.Kind.NOT_EQUAL_TO) {
      return emitEquality(kind == Tree.Kind.EQUAL_TO, lp, lt, rp, rt, ctx);
    }

    String lRaw = emitPath(lp, ctx);
    String rRaw = emitPath(rp, ctx);
    return computeNumericOrBoolBinary(kind, lRaw, lt, rRaw, rt, resType, ctx);
  }

  private void collectStringConcatOperandPaths(TreePath path, List<TreePath> out) {
    Tree leaf = path.getLeaf();
    if (leaf instanceof ParenthesizedTree pt) {
      collectStringConcatOperandPaths(new TreePath(path, pt.getExpression()), out);
      return;
    }
    if (leaf instanceof BinaryTree bt
        && bt.getKind() == Tree.Kind.PLUS
        && resolver.isStringType(typeOf(path))) {
      collectStringConcatOperandPaths(new TreePath(path, bt.getLeftOperand()), out);
      collectStringConcatOperandPaths(new TreePath(path, bt.getRightOperand()), out);
      return;
    }
    out.add(path);
  }

  private String emitStringConcatOperand(TreePath partPath, EmitContext ctx) {
    while (partPath.getLeaf() instanceof ParenthesizedTree pt) {
      partPath = new TreePath(partPath, pt.getExpression());
    }
    Tree leaf = partPath.getLeaf();
    String rt = rtPrefix(ctx);

    if (leaf instanceof LiteralTree lit) {
      Object val = lit.getValue();
      switch (lit.getKind()) {
        case NULL_LITERAL -> {
          return "\"null\"";
        }
        case STRING_LITERAL -> {
          return GoNames.goStringLiteral((String) val);
        }
        case BOOLEAN_LITERAL -> {
          return Boolean.TRUE.equals(val) ? "\"true\"" : "\"false\"";
        }
        case INT_LITERAL, LONG_LITERAL -> {
          return GoNames.goStringLiteral(String.valueOf(val));
        }
        case CHAR_LITERAL -> {
          char ch = (Character) val;
          if (!Character.isSurrogate(ch)) {
            return GoNames.goStringLiteral(String.valueOf(ch));
          }
          return rt + "StrValChar(uint16(" + (int) ch + "))";
        }
        default -> {}
      }
    }

    if (leaf instanceof IdentifierTree || leaf instanceof MemberSelectTree) {
      boolean hasInstanceReceiver = false;
      if (leaf instanceof MemberSelectTree ms) {
        Element qualEl = resolver.trees.getElement(new TreePath(partPath, ms.getExpression()));
        if (!(qualEl instanceof TypeElement) && !ms.getExpression().toString().equals("this")
            && !ms.getExpression().toString().equals("super")) {
          hasInstanceReceiver = true;
        }
      }
      if (!hasInstanceReceiver) {
        Element el = resolver.trees.getElement(partPath);
        if (el instanceof VariableElement ve
            && ve.getModifiers().contains(Modifier.STATIC)
            && ve.getModifiers().contains(Modifier.FINAL)
            && ve.getConstantValue() instanceof String constStr) {
          return GoNames.goStringLiteral(constStr);
        }
      }
    }

    TypeMirror rawTm = typeOf(partPath);
    TypeMirror tm = resolver.safeErasure(rawTm);
    ExpressionTree exprLeaf = (ExpressionTree) leaf;
    TreePath parentPath = partPath.getParentPath();

    if (resolver.isStringType(tm)) {
      return rt + "StrValString(" + emitAs(exprLeaf, parentPath, tm, false, ctx) + ")";
    }
    return switch (tm.getKind()) {
      case INT -> rt + "StrValInt32(" + emitAs(exprLeaf, parentPath, tm, false, ctx) + ")";
      case BYTE, SHORT ->
          rt + "StrValInt32(int32(" + emitAs(exprLeaf, parentPath, tm, false, ctx) + "))";
      case LONG -> rt + "StrValInt64(" + emitAs(exprLeaf, parentPath, tm, false, ctx) + ")";
      case BOOLEAN -> rt + "StrValBool(" + emitAs(exprLeaf, parentPath, tm, false, ctx) + ")";
      case CHAR -> rt + "StrValChar(" + emitAs(exprLeaf, parentPath, tm, false, ctx) + ")";
      case DOUBLE ->
          jlPrefix(ctx) + "DoubleToString(" + emitAs(exprLeaf, parentPath, tm, false, ctx) + ")";
      case FLOAT ->
          jlPrefix(ctx) + "FloatToString(" + emitAs(exprLeaf, parentPath, tm, false, ctx) + ")";
      default ->
          resolver.isActiveTypeVariable(rawTm)
              ? (rt + "StrVal(" + rt + "ToAny(" + emitPath(partPath, ctx) + "))")
              : (rt + "StrVal(" + emitPath(partPath, ctx) + ")");
    };
  }

  private String emitEquality(
      boolean isEq,
      TreePath lp,
      TypeMirror lt,
      TreePath rp,
      TypeMirror rt,
      EmitContext ctx) {
    String op = isEq ? " == " : " != ";
    TypeMirror ler = resolver.safeErasure(lt);
    TypeMirror rer = resolver.safeErasure(rt);

    // Numeric or boolean equality (if either side is primitive, JLS unboxes the other!)
    if (ler.getKind().isPrimitive() || rer.getKind().isPrimitive()) {
      if (ler.getKind() == TypeKind.BOOLEAN || rer.getKind() == TypeKind.BOOLEAN) {
        TypeMirror btm = resolver.types.getPrimitiveType(TypeKind.BOOLEAN);
        return "("
            + convert(emitPath(lp, ctx), lt, btm, false, ctx)
            + op
            + convert(emitPath(rp, ctx), rt, btm, false, ctx)
            + ")";
      }
      TypeMirror prom = promotedNumericPair(ler, rer);
      return "("
          + convert(emitPath(lp, ctx), lt, prom, false, ctx)
          + op
          + convert(emitPath(rp, ctx), rt, prom, false, ctx)
          + ")";
    }

    // Node equality (fast uint32 handle equality!)
    boolean lNode = resolver.isNodeType(ler);
    boolean rNode = resolver.isNodeType(rer);
    if ((lNode && (rNode || rer.getKind() == TypeKind.NULL))
        || (rNode && ler.getKind() == TypeKind.NULL)) {
      TypeMirror nodeTm = lNode ? ler : rer;
      String l = ler.getKind() == TypeKind.NULL ? "0" : convert(emitPath(lp, ctx), lt, nodeTm, false, ctx);
      String r = rer.getKind() == TypeKind.NULL ? "0" : convert(emitPath(rp, ctx), rt, nodeTm, false, ctx);
      return "(" + l + op + r + ")";
    }

    // String vs null or String vs String
    if (resolver.isStringType(ler) && rer.getKind() == TypeKind.NULL) {
      return "(" + emitPath(lp, ctx) + op + rtPrefix(ctx) + "NullString)";
    }
    if (ler.getKind() == TypeKind.NULL && resolver.isStringType(rer)) {
      return "(" + emitPath(rp, ctx) + op + rtPrefix(ctx) + "NullString)";
    }
    if (resolver.isStringType(ler) && resolver.isStringType(rer)) {
      return "(" + emitPath(lp, ctx) + op + emitPath(rp, ctx) + ")";
    }

    // Fast Node Token ordinal comparison (`n.getToken() == Token.X`, `this.token == Token.X`, `n.getToken() == null`, `n.token == m.token`)
    TreePath lu = unwrapParens(lp);
    TreePath ru = unwrapParens(rp);
    String lNodeTok = tryEmitNodeForTokenRead(lu, ctx);
    if (lNodeTok != null) {
      int rOrd = tryGetStaticTokenEnumOrdinal(ru);
      if (rOrd >= 0 && rOrd <= 254) {
        return "(" + lNodeTok + ".RawToken_()" + op + "uint8(" + rOrd + "))";
      }
      if (rOrd >= 0) {
        return "(" + lNodeTok + ".TokenOrd_()" + op + "int32(" + rOrd + "))";
      }
      if (ru.getLeaf().getKind() == Tree.Kind.NULL_LITERAL) {
        return "(" + lNodeTok + ".RawToken_()" + (isEq ? " == uint8(0xFF))" : " != uint8(0xFF))");
      }
      String rNodeTok = tryEmitNodeForTokenRead(ru, ctx);
      if (rNodeTok != null) {
        return "(" + lNodeTok + ".RawToken_()" + op + rNodeTok + ".RawToken_())";
      }
    } else {
      int lOrd = tryGetStaticTokenEnumOrdinal(lu);
      if (lOrd >= 0) {
        String rNodeTok = tryEmitNodeForTokenRead(ru, ctx);
        if (rNodeTok != null) {
          if (lOrd <= 254) {
            return "(" + rNodeTok + ".RawToken_()" + op + "uint8(" + lOrd + "))";
          }
          return "(" + rNodeTok + ".TokenOrd_()" + op + "int32(" + lOrd + "))";
        }
      } else if (lu.getLeaf().getKind() == Tree.Kind.NULL_LITERAL) {
        String rNodeTok = tryEmitNodeForTokenRead(ru, ctx);
        if (rNodeTok != null) {
          return "(" + rNodeTok + ".RawToken_()" + (isEq ? " == uint8(0xFF))" : " != uint8(0xFF))");
        }
      }
    }

    // Null check on pointer, slice, value struct, or interface
    if (rer.getKind() == TypeKind.NULL || ler.getKind() == TypeKind.NULL) {
      TreePath nonNullPath = rer.getKind() == TypeKind.NULL ? lp : rp;
      TypeMirror nonNullRaw = rer.getKind() == TypeKind.NULL ? lt : rt;
      TypeMirror nonNullEr = rer.getKind() == TypeKind.NULL ? ler : rer;
      String x = emitPath(nonNullPath, ctx);
      if (resolver.isActiveComparableTypeVariable(nonNullRaw)) {
        String goTv = resolver.goType(nonNullRaw, ctx.imports);
        return "("
            + x
            + op
            + "func() "
            + goTv
            + " { if _, ok_ := any((*"
            + goTv
            + ")(nil)).(*string); ok_ { return any("
            + rtPrefix(ctx)
            + "NullString).("
            + goTv
            + ") }; var z_ "
            + goTv
            + "; return z_ }())";
      }
      if (resolver.isActiveTypeVariable(nonNullRaw)) {
        return "(" + rtPrefix(ctx) + "ToAny(" + x + ")" + op + "nil)";
      }
      if (resolver.isValueEnumType(nonNullEr)) {
        return "(" + x + op + "0)";
      }
      if (resolver.isByteStringType(nonNullEr)) {
        return isEq ? "(" + x + ").IsJavaNil_()" : "!(" + x + ").IsJavaNil_()";
      }
      if (resolver.isValueStructType(nonNullEr)) {
        String zero = "(" + resolver.goType(nonNullRaw, ctx.imports) + "{})";
        return "(" + x + op + zero + ")";
      }
      if (isSourceJavaNilIface(nonNullEr)) {
        if (isSimpleLocalOrParamIdent(nonNullPath, x)) {
          return isEq
              ? "(" + x + " == nil || " + x + ".IsJavaNil_())"
              : "(" + x + " != nil && !" + x + ".IsJavaNil_())";
        }
        String goIface = resolver.goType(nonNullRaw, ctx.imports);
        String call = rtPrefix(ctx) + "IfaceEquals[" + goIface + "](" + x + ", nil)";
        return isEq ? call : ("!(" + call + ")");
      }
      return "(" + x + op + "nil)";
    }

    String lGo = resolver.goType(lt, ctx.imports);
    String rGo = resolver.goType(rt, ctx.imports);

    // Fast source Go interface identity equality (`IfaceEquals[I](l, r)`)
    TypeMirror ltEff = resolver.goEffectiveType(lt);
    TypeMirror rtEff = resolver.goEffectiveType(rt);
    if (isSourceJavaNilIface(rer)
        && (lGo.equals(rGo)
            || (isSourceRefAssignableToIface(ler) && resolver.types.isAssignable(ltEff, rtEff)))) {
      String lExpr = emitPathOrReusableSelf(lp, rtEff, false, ctx);
      String rExpr = emitPathOrReusableSelf(rp, rtEff, false, ctx);
      String eqCall = rtPrefix(ctx) + "IfaceEquals[" + rGo + "](" + lExpr + ", " + rExpr + ")";
      return isEq ? eqCall : ("!(" + eqCall + ")");
    }
    if (isSourceJavaNilIface(ler)
        && isSourceRefAssignableToIface(rer)
        && resolver.types.isAssignable(rtEff, ltEff)) {
      String lExpr = emitPathOrReusableSelf(lp, ltEff, false, ctx);
      String rExpr = emitPathOrReusableSelf(rp, ltEff, false, ctx);
      String eqCall = rtPrefix(ctx) + "IfaceEquals[" + lGo + "](" + lExpr + ", " + rExpr + ")";
      return isEq ? eqCall : ("!(" + eqCall + ")");
    }

    if (resolver.isByteStringType(ler) && resolver.isByteStringType(rer)) {
      String eqCall =
          rtPrefix(ctx) + "ByteStringRefEquals(" + emitPath(lp, ctx) + ", " + emitPath(rp, ctx) + ")";
      return isEq ? eqCall : ("!(" + eqCall + ")");
    }

    // Same concrete struct pointer, value struct, or value enum type
    if ((lGo.startsWith("*") || resolver.isValueStructType(ler) || resolver.isValueEnumType(ler))
        && lGo.equals(rGo)) {
      return "(" + emitPath(lp, ctx) + op + emitPath(rp, ctx) + ")";
    }

    // Same Go slice type (`T[] == T[]` / `T[] != T[]`)
    if (ler.getKind() == TypeKind.ARRAY
        && rer.getKind() == TypeKind.ARRAY
        && lGo.startsWith("[]")
        && lGo.equals(rGo)) {
      String eqCall =
          rtPrefix(ctx) + "SliceRefEquals(" + emitPath(lp, ctx) + ", " + emitPath(rp, ctx) + ")";
      return isEq ? eqCall : ("!(" + eqCall + ")");
    }

    // General reference identity across interfaces, Object, arrays, etc.
    TypeMirror objTm = resolver.getObjectElement().asType();
    String lAny = emitPathOrReusableSelf(lp, objTm, true, ctx);
    String rAny = emitPathOrReusableSelf(rp, objTm, true, ctx);
    String eqCall = rtPrefix(ctx) + "RefEquals(" + lAny + ", " + rAny + ")";
    return isEq ? eqCall : ("!(" + eqCall + ")");
  }

  private static TreePath unwrapParens(TreePath tp) {
    while (tp != null && tp.getLeaf() instanceof ParenthesizedTree pt) {
      tp = new TreePath(tp, pt.getExpression());
    }
    return tp;
  }

  /**
   * If {@code exprPath} is a read of {@code Node.getToken()} or {@code Node.token}, emits and
   * returns the Go {@code Node} receiver expression; otherwise returns {@code null}.
   */
  String tryEmitNodeForTokenRead(TreePath exprPath, EmitContext ctx) {
    TreePath p = unwrapParens(exprPath);
    if (p == null) {
      return null;
    }
    Tree leaf = p.getLeaf();
    if (leaf instanceof MethodInvocationTree mi) {
      if (!mi.getArguments().isEmpty()) {
        return null;
      }
      Element el = resolver.trees.getElement(p);
      if (!(el instanceof ExecutableElement ee)
          || ee.getModifiers().contains(Modifier.STATIC)
          || !ee.getParameters().isEmpty()
          || !ee.getSimpleName().contentEquals("getToken")
          || !(ee.getEnclosingElement() instanceof TypeElement owner)
          || !resolver.isNodeClass(owner)) {
        return null;
      }
      ExpressionTree sel = mi.getMethodSelect();
      if (sel instanceof MemberSelectTree ms) {
        ExpressionTree recvTree = ms.getExpression();
        if (recvTree instanceof IdentifierTree id && id.getName().contentEquals("super")) {
          return ctx.receiverName;
        }
        if (recvTree instanceof MemberSelectTree rms
            && rms.getIdentifier().contentEquals("super")) {
          return null;
        }
        return emitAs(recvTree, new TreePath(p, ms), owner.asType(), false, ctx);
      }
      if (sel instanceof IdentifierTree id) {
        if (id.getName().contentEquals("super")) {
          return null;
        }
        return resolveImplicitThisValue(owner, ctx);
      }
      return null;
    }

    if (leaf instanceof MemberSelectTree ms) {
      if (!ms.getIdentifier().contentEquals("token")) {
        return null;
      }
      Element el = resolver.trees.getElement(p);
      if (!(el instanceof VariableElement ve)
          || ve.getKind() != ElementKind.FIELD
          || ve.getModifiers().contains(Modifier.STATIC)
          || !(ve.getEnclosingElement() instanceof TypeElement owner)
          || !resolver.isNodeClass(owner)) {
        return null;
      }
      ExpressionTree recvTree = ms.getExpression();
      if (recvTree instanceof IdentifierTree id && id.getName().contentEquals("super")) {
        return ctx.receiverName;
      }
      if (recvTree instanceof MemberSelectTree rms && rms.getIdentifier().contentEquals("super")) {
        return null;
      }
      return emitAs(recvTree, p, owner.asType(), false, ctx);
    }

    if (leaf instanceof IdentifierTree id) {
      if (!id.getName().contentEquals("token")) {
        return null;
      }
      Element el = resolver.trees.getElement(p);
      if (!(el instanceof VariableElement ve)
          || ve.getKind() != ElementKind.FIELD
          || ve.getModifiers().contains(Modifier.STATIC)
          || !(ve.getEnclosingElement() instanceof TypeElement owner)
          || !resolver.isNodeClass(owner)) {
        return null;
      }
      return resolveImplicitThisValue(owner, ctx);
    }

    return null;
  }

  private int tryGetStaticTokenEnumOrdinal(TreePath exprPath) {
    TreePath p = unwrapParens(exprPath);
    if (p == null) {
      return -1;
    }
    Tree leaf = p.getLeaf();
    if (leaf instanceof MemberSelectTree ms) {
      Element qualEl = resolver.trees.getElement(new TreePath(p, ms.getExpression()));
      if (!(qualEl instanceof TypeElement)) {
        return -1;
      }
    } else if (!(leaf instanceof IdentifierTree)) {
      return -1;
    }
    Element el = resolver.trees.getElement(p);
    if (!(el instanceof VariableElement ve)
        || ve.getKind() != ElementKind.ENUM_CONSTANT
        || !(ve.getEnclosingElement() instanceof TypeElement ete)
        || !ete.getQualifiedName().contentEquals("com.google.javascript.rhino.Token")) {
      return -1;
    }
    int ord = 0;
    for (Element enc : ete.getEnclosedElements()) {
      if (enc instanceof VariableElement eve && eve.getKind() == ElementKind.ENUM_CONSTANT) {
        if (eve.equals(ve)) {
          return ord;
        }
        ord++;
      }
    }
    return -1;
  }

  private String computeNumericOrBoolBinary(
      Tree.Kind kind,
      String lRaw,
      TypeMirror lt,
      String rRaw,
      TypeMirror rt,
      TypeMirror resType,
      EmitContext ctx) {
    TypeMirror ler = resolver.safeErasure(lt);
    TypeMirror rer = resolver.safeErasure(rt);

    // Boolean `&`, `|`, `^`
    if ( isBooleanOrBoxedBoolean(ler) && isBooleanOrBoxedBoolean(rer)
        && (kind == Tree.Kind.AND || kind == Tree.Kind.OR || kind == Tree.Kind.XOR)) {
      TypeMirror btm = resolver.types.getPrimitiveType(TypeKind.BOOLEAN);
      String l = convert(lRaw, lt, btm, false, ctx);
      String r = convert(rRaw, rt, btm, false, ctx);
      return switch (kind) {
        case AND -> {
          String lt2 = ctx.nextTemp("ba");
          String rt2 = ctx.nextTemp("bb");
          yield "func() bool { " + lt2 + " := " + l + "; " + rt2 + " := " + r + "; return " + lt2 + " && " + rt2 + " }()";
        }
        case OR -> {
          String lt2 = ctx.nextTemp("ba");
          String rt2 = ctx.nextTemp("bb");
          yield "func() bool { " + lt2 + " := " + l + "; " + rt2 + " := " + r + "; return " + lt2 + " || " + rt2 + " }()";
        }
        case XOR -> "((" + l + ") != (" + r + "))";
        default -> l;
      };
    }

    // Shifts `<<`, `>>`, `>>>`
    if (kind == Tree.Kind.LEFT_SHIFT
        || kind == Tree.Kind.RIGHT_SHIFT
        || kind == Tree.Kind.UNSIGNED_RIGHT_SHIFT) {
      boolean isLong = unboxedPrimKind(ler) == TypeKind.LONG;
      TypeMirror lTarget =
          resolver.types.getPrimitiveType(isLong ? TypeKind.LONG : TypeKind.INT);
      TypeMirror rTarget = resolver.types.getPrimitiveType(TypeKind.INT);
      String l = convert(lRaw, lt, lTarget, false, ctx);
      String r = convert(rRaw, rt, rTarget, false, ctx);
      String jl = jlPrefix(ctx);
      if (isLong) {
        return switch (kind) {
          case LEFT_SHIFT -> jl + "Shl64(" + l + ", int(" + r + "))";
          case RIGHT_SHIFT -> jl + "Shr64(" + l + ", int(" + r + "))";
          default -> jl + "Ushr64(" + l + ", int(" + r + "))";
        };
      }
      return switch (kind) {
        case LEFT_SHIFT -> jl + "Shl32(" + l + ", int(" + r + "))";
        case RIGHT_SHIFT -> jl + "Shr32(" + l + ", int(" + r + "))";
        default -> jl + "Ushr32(" + l + ", int(" + r + "))";
      };
    }

    // Standard binary numeric promotion
    TypeMirror prom = promotedNumericPair(ler, rer);
    String l = convert(lRaw, lt, prom, false, ctx);
    String r = convert(rRaw, rt, prom, false, ctx);

    if (kind == Tree.Kind.REMAINDER
        && (prom.getKind() == TypeKind.DOUBLE || prom.getKind() == TypeKind.FLOAT)) {
      ctx.imports.std("math");
      if (prom.getKind() == TypeKind.FLOAT) {
        return "float32(math.Mod(float64(" + l + "), float64(" + r + ")))";
      }
      return "math.Mod(" + l + ", " + r + ")";
    }

    String op =
        switch (kind) {
          case PLUS -> " + ";
          case MINUS -> " - ";
          case MULTIPLY -> " * ";
          case DIVIDE -> " / ";
          case REMAINDER -> " % ";
          case AND -> " & ";
          case OR -> " | ";
          case XOR -> " ^ ";
          case LESS_THAN -> " < ";
          case LESS_THAN_EQUAL -> " <= ";
          case GREATER_THAN -> " > ";
          case GREATER_THAN_EQUAL -> " >= ";
          default -> " + ";
        };
    return "(" + l + op + r + ")";
  }

  private boolean isBooleanOrBoxedBoolean(TypeMirror tm) {
    if (tm.getKind() == TypeKind.BOOLEAN) return true;
    return tm instanceof DeclaredType dt
        && dt.asElement() instanceof TypeElement te
        && te.getQualifiedName().contentEquals("java.lang.Boolean");
  }

  private TypeKind unboxedPrimKind(TypeMirror tm) {
    if (tm.getKind().isPrimitive()) return tm.getKind();
    if (tm instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      return switch (te.getQualifiedName().toString()) {
        case "java.lang.Double" -> TypeKind.DOUBLE;
        case "java.lang.Float" -> TypeKind.FLOAT;
        case "java.lang.Long" -> TypeKind.LONG;
        case "java.lang.Character" -> TypeKind.CHAR;
        case "java.lang.Short" -> TypeKind.SHORT;
        case "java.lang.Byte" -> TypeKind.BYTE;
        case "java.lang.Boolean" -> TypeKind.BOOLEAN;
        default -> TypeKind.INT;
      };
    }
    return TypeKind.INT;
  }

  private TypeMirror promotedNumericPair(TypeMirror a, TypeMirror b) {
    TypeKind ka = unboxedPrimKind(a);
    TypeKind kb = unboxedPrimKind(b);
    if (ka == TypeKind.DOUBLE || kb == TypeKind.DOUBLE) {
      return resolver.types.getPrimitiveType(TypeKind.DOUBLE);
    }
    if (ka == TypeKind.FLOAT || kb == TypeKind.FLOAT) {
      return resolver.types.getPrimitiveType(TypeKind.FLOAT);
    }
    if (ka == TypeKind.LONG || kb == TypeKind.LONG) {
      return resolver.types.getPrimitiveType(TypeKind.LONG);
    }
    return resolver.types.getPrimitiveType(TypeKind.INT);
  }

  private TypeMirror promotedBinaryType(Tree.Kind binKind, TypeMirror lt, TypeMirror rt) {
    if (binKind == Tree.Kind.LEFT_SHIFT
        || binKind == Tree.Kind.RIGHT_SHIFT
        || binKind == Tree.Kind.UNSIGNED_RIGHT_SHIFT) {
      boolean isLong = unboxedPrimKind(resolver.safeErasure(lt)) == TypeKind.LONG;
      return resolver.types.getPrimitiveType(isLong ? TypeKind.LONG : TypeKind.INT);
    }
    TypeMirror ler = resolver.safeErasure(lt);
    if (isBooleanOrBoxedBoolean(ler)) {
      return resolver.types.getPrimitiveType(TypeKind.BOOLEAN);
    }
    return promotedNumericPair(ler, resolver.safeErasure(rt));
  }

  private Tree.Kind compoundToBinaryKind(Tree.Kind ck) {
    return switch (ck) {
      case MULTIPLY_ASSIGNMENT -> Tree.Kind.MULTIPLY;
      case DIVIDE_ASSIGNMENT -> Tree.Kind.DIVIDE;
      case REMAINDER_ASSIGNMENT -> Tree.Kind.REMAINDER;
      case PLUS_ASSIGNMENT -> Tree.Kind.PLUS;
      case MINUS_ASSIGNMENT -> Tree.Kind.MINUS;
      case LEFT_SHIFT_ASSIGNMENT -> Tree.Kind.LEFT_SHIFT;
      case RIGHT_SHIFT_ASSIGNMENT -> Tree.Kind.RIGHT_SHIFT;
      case UNSIGNED_RIGHT_SHIFT_ASSIGNMENT -> Tree.Kind.UNSIGNED_RIGHT_SHIFT;
      case AND_ASSIGNMENT -> Tree.Kind.AND;
      case XOR_ASSIGNMENT -> Tree.Kind.XOR;
      case OR_ASSIGNMENT -> Tree.Kind.OR;
      default -> Tree.Kind.PLUS;
    };
  }

  // --- Casts, InstanceOf, Conditional ---

  private String emitTypeCast(TypeCastTree tc, TreePath path, EmitContext ctx) {
    TypeMirror targetTm = typeOf(path);
    ExpressionTree inner = tc.getExpression();
    while (inner instanceof ParenthesizedTree pt) {
      inner = pt.getExpression();
    }
    if (inner instanceof NewArrayTree na
        && (na.getInitializers() != null || na.getDimensions().size() == 1)
        && resolver.goEffectiveType(targetTm) instanceof ArrayType at) {
      TreePath naPath = new TreePath(path, tc.getExpression());
      while (naPath.getLeaf() instanceof ParenthesizedTree pt) {
        naPath = new TreePath(naPath, pt.getExpression());
      }
      TypeMirror compTm = at.getComponentType();
      String arrGo = resolver.goType(at, ctx.imports);
      if (na.getInitializers() != null) {
        List<String> elems = new ArrayList<>();
        for (ExpressionTree init : na.getInitializers()) {
          elems.add(emitAs(init, naPath, compTm, false, ctx));
        }
        return arrGo + "{" + String.join(", ", elems) + "}";
      }
      String lenExpr =
          emitAs(
              na.getDimensions().get(0),
              naPath,
              resolver.types.getPrimitiveType(TypeKind.INT),
              false,
              ctx);
      if (resolver.isStringType(compTm)) {
        return rtPrefix(ctx) + "NewStringArray(" + lenExpr + ")";
      }
      if (resolver.isActiveTypeVariable(compTm)) {
        String compGo = resolver.goType(compTm, ctx.imports);
        String sTmp = ctx.nextTemp("arr");
        return "func() "
            + arrGo
            + " { "
            + sTmp
            + " := make("
            + arrGo
            + ", "
            + lenExpr
            + "); if _, ok_ := any((*"
            + compGo
            + ")(nil)).(*string); ok_ { ns_ := any("
            + rtPrefix(ctx)
            + "NullString).("
            + compGo
            + "); for i_ := range "
            + sTmp
            + " { "
            + sTmp
            + "[i_] = ns_ } }; return "
            + sTmp
            + " }()";
      }
      return "make(" + arrGo + ", " + lenExpr + ")";
    }
    if (targetTm != null) {
      TypeKind tk = targetTm.getKind();
      if (tk == TypeKind.BYTE || tk == TypeKind.SHORT || tk == TypeKind.CHAR || tk == TypeKind.INT) {
        Long cv = tryEvalIntegralConst(new TreePath(path, tc.getExpression()));
        if (cv != null) {
          return switch (tk) {
            case BYTE -> "int8(" + cv.byteValue() + ")";
            case SHORT -> "int16(" + cv.shortValue() + ")";
            case CHAR -> "uint16(" + (int) (char) cv.longValue() + ")";
            case INT -> "int32(" + cv.intValue() + ")";
            default -> emitAs(tc.getExpression(), path, targetTm, false, ctx);
          };
        }
      }
    }
    return emitAs(tc.getExpression(), path, targetTm, false, ctx);
  }

  private Long tryEvalIntegralConst(TreePath exprPath) {
    if (exprPath == null) {
      return null;
    }
    Tree leaf = exprPath.getLeaf();
    if (leaf instanceof ParenthesizedTree pt) {
      return tryEvalIntegralConst(new TreePath(exprPath, pt.getExpression()));
    }
    if (leaf instanceof LiteralTree lt) {
      Object v = lt.getValue();
      if (v instanceof Number n
          && (lt.getKind() == Tree.Kind.INT_LITERAL || lt.getKind() == Tree.Kind.LONG_LITERAL)) {
        return n.longValue();
      }
      if (v instanceof Character c && lt.getKind() == Tree.Kind.CHAR_LITERAL) {
        return (long) c.charValue();
      }
      return null;
    }
    if (leaf instanceof IdentifierTree || leaf instanceof MemberSelectTree) {
      Element el = resolver.trees.getElement(exprPath);
      if (el instanceof VariableElement ve
          && ve.getModifiers().contains(Modifier.FINAL)
          && ve.getConstantValue() != null) {
        Object cv = ve.getConstantValue();
        if (cv instanceof Integer || cv instanceof Long || cv instanceof Short || cv instanceof Byte) {
          return ((Number) cv).longValue();
        }
        if (cv instanceof Character c) {
          return (long) c.charValue();
        }
      }
      return null;
    }
    if (leaf instanceof UnaryTree ut) {
      Long sub = tryEvalIntegralConst(new TreePath(exprPath, ut.getExpression()));
      if (sub == null) {
        return null;
      }
      TypeMirror tm = typeOf(exprPath);
      boolean isLong = tm != null && tm.getKind() == TypeKind.LONG;
      return switch (ut.getKind()) {
        case UNARY_PLUS -> isLong ? sub : (long) sub.intValue();
        case UNARY_MINUS -> isLong ? -sub : (long) -sub.intValue();
        case BITWISE_COMPLEMENT -> isLong ? ~sub : (long) ~sub.intValue();
        default -> null;
      };
    }
    if (leaf instanceof BinaryTree bt) {
      Long l = tryEvalIntegralConst(new TreePath(exprPath, bt.getLeftOperand()));
      Long r = tryEvalIntegralConst(new TreePath(exprPath, bt.getRightOperand()));
      if (l == null || r == null) {
        return null;
      }
      TypeMirror tm = typeOf(exprPath);
      boolean isLong = tm != null && tm.getKind() == TypeKind.LONG;
      return switch (bt.getKind()) {
        case PLUS -> isLong ? (l + r) : (long) (l.intValue() + r.intValue());
        case MINUS -> isLong ? (l - r) : (long) (l.intValue() - r.intValue());
        case MULTIPLY -> isLong ? (l * r) : (long) (l.intValue() * r.intValue());
        case AND -> isLong ? (l & r) : (long) (l.intValue() & r.intValue());
        case OR -> isLong ? (l | r) : (long) (l.intValue() | r.intValue());
        case XOR -> isLong ? (l ^ r) : (long) (l.intValue() ^ r.intValue());
        case LEFT_SHIFT -> isLong ? (l << r.intValue()) : (long) (l.intValue() << r.intValue());
        case RIGHT_SHIFT -> isLong ? (l >> r.intValue()) : (long) (l.intValue() >> r.intValue());
        case UNSIGNED_RIGHT_SHIFT ->
            isLong ? (l >>> r.intValue()) : (long) (l.intValue() >>> r.intValue());
        default -> null;
      };
    }
    if (leaf instanceof TypeCastTree innerTc) {
      TypeMirror innerTm = typeOf(exprPath);
      if (innerTm != null) {
        Long sub = tryEvalIntegralConst(new TreePath(exprPath, innerTc.getExpression()));
        if (sub != null) {
          return switch (innerTm.getKind()) {
            case BYTE -> (long) sub.byteValue();
            case SHORT -> (long) sub.shortValue();
            case CHAR -> (long) (char) sub.longValue();
            case INT -> (long) sub.intValue();
            case LONG -> sub;
            default -> null;
          };
        }
      }
    }
    return null;
  }

  public void emitPatternTryCastAssign(
      InstanceOfTree io,
      TreePath ioPath,
      String okVar,
      boolean declareOkVar,
      String indentTabs,
      StringBuilder out,
      EmitContext ctx) {
    TreePath exprPath = new TreePath(ioPath, io.getExpression());
    TypeMirror exprTm = typeOf(exprPath);
    String exprGo = emitPath(exprPath, ctx);
    String rt = rtPrefix(ctx);
    BindingPatternTree bpt = (BindingPatternTree) io.getPattern();
    Element bindEl =
        resolver.trees.getElement(new TreePath(new TreePath(ioPath, bpt), bpt.getVariable()));
    VariableElement ve = (VariableElement) bindEl;
    String varName = ctx.registerLocalVar(ve);
    TypeMirror patTm = ve.asType();

    if (declareOkVar) {
      out.append(indentTabs).append("var ").append(okVar).append(" bool\n");
    }
    if (resolver.isStringType(patTm)) {
      String exprAny = convert(exprGo, exprTm, resolver.getObjectElement().asType(), true, ctx);
      out.append(indentTabs)
          .append(varName)
          .append(", ")
          .append(okVar)
          .append(" = ")
          .append(rt)
          .append("TryCastString(")
          .append(exprAny)
          .append(")\n");
      return;
    }
    if (patTm instanceof DeclaredType dt
        && dt.asElement() instanceof TypeElement te
        && resolver.isNodeClass(te)
        && resolver.getNodeKindFlag(te) > 0) {
      int kf = resolver.getNodeKindFlag(te);
      String ns = nsPrefix(ctx);
      if (resolver.isNodeType(exprTm)) {
        String nTmp = ctx.nextTemp("n");
        out.append(indentTabs)
            .append("if ")
            .append(nTmp)
            .append(" := ")
            .append(exprGo)
            .append("; ")
            .append(nTmp)
            .append(" != 0 && ")
            .append(ns)
            .append("Node(")
            .append(nTmp)
            .append(").KindFlags() == ")
            .append(kf)
            .append(" {\n")
            .append(indentTabs)
            .append('\t')
            .append(varName)
            .append(" = ")
            .append(nTmp)
            .append('\n')
            .append(indentTabs)
            .append('\t')
            .append(okVar)
            .append(" = true\n")
            .append(indentTabs)
            .append("} else {\n")
            .append(indentTabs)
            .append('\t')
            .append(okVar)
            .append(" = false\n")
            .append(indentTabs)
            .append("}\n");
      } else {
        String exprAny = convert(exprGo, exprTm, resolver.getObjectElement().asType(), true, ctx);
        String nodeGo = resolver.goType(patTm, ctx.imports);
        String nTmp = ctx.nextTemp("n");
        String okTmp = ctx.nextTemp("ok");
        out.append(indentTabs)
            .append("if ")
            .append(nTmp)
            .append(", ")
            .append(okTmp)
            .append(" := ")
            .append(exprAny)
            .append(".(")
            .append(nodeGo)
            .append("); ")
            .append(okTmp)
            .append(" && ")
            .append(nTmp)
            .append(" != 0 && ")
            .append(ns)
            .append("Node(")
            .append(nTmp)
            .append(").KindFlags() == ")
            .append(kf)
            .append(" {\n")
            .append(indentTabs)
            .append('\t')
            .append(varName)
            .append(" = ")
            .append(nTmp)
            .append('\n')
            .append(indentTabs)
            .append('\t')
            .append(okVar)
            .append(" = true\n")
            .append(indentTabs)
            .append("} else {\n")
            .append(indentTabs)
            .append('\t')
            .append(okVar)
            .append(" = false\n")
            .append(indentTabs)
            .append("}\n");
      }
      return;
    }
    String exprAny = convert(exprGo, exprTm, resolver.getObjectElement().asType(), true, ctx);
    String patGo = resolver.goType(patTm, ctx.imports);
    out.append(indentTabs)
        .append(varName)
        .append(", ")
        .append(okVar)
        .append(" = ")
        .append(rt)
        .append("TryCast[")
        .append(patGo)
        .append("](")
        .append(exprAny)
        .append(")\n");
  }

  private String emitInstanceOf(InstanceOfTree io, TreePath path, EmitContext ctx) {
    TreePath exprPath = new TreePath(path, io.getExpression());
    TypeMirror exprTm = typeOf(exprPath);
    String exprGo = emitPath(exprPath, ctx);
    String exprAny = convert(exprGo, exprTm, resolver.getObjectElement().asType(), true, ctx);
    String rt = rtPrefix(ctx);

    PatternTree pat = io.getPattern();
    if (pat instanceof BindingPatternTree bpt) {
      Element bindEl = resolver.trees.getElement(new TreePath(new TreePath(path, pat), bpt.getVariable()));
      VariableElement ve = (VariableElement) bindEl;
      String varName = ctx.registerLocalVar(ve);
      TypeMirror patTm = ve.asType();

      if (resolver.isStringType(patTm)) {
        return rt + "InstanceOfStringBind(" + exprAny + ", &" + varName + ")";
      }
      if (patTm instanceof DeclaredType dt
          && dt.asElement() instanceof TypeElement te
          && resolver.isNodeClass(te)
          && resolver.getNodeKindFlag(te) > 0) {
        int kf = resolver.getNodeKindFlag(te);
        String ns = nsPrefix(ctx);
        String nodeGo = resolver.goType(patTm, ctx.imports);
        return "func() bool { if n_, ok_ := "
            + exprAny
            + ".("
            + nodeGo
            + "); ok_ && n_ != 0 && "
            + ns
            + "Node(n_).KindFlags() == "
            + kf
            + " { "
            + varName
            + " = n_; return true }; return false }()";
      }
      if (patTm instanceof DeclaredType dt
          && dt.asElement() instanceof TypeElement te
          && resolver.isNullableValueStructClass(te)) {
        String patGo = resolver.goType(patTm, ctx.imports);
        return "func() bool { if v_, ok_ := "
            + exprAny
            + ".("
            + patGo
            + "); ok_ && !v_.IsJavaNil_() { "
            + varName
            + " = v_; return true }; return false }()";
      }
      String patGo = resolver.goType(patTm, ctx.imports);
      return rt + "InstanceOfBind[" + patGo + "](" + exprAny + ", &" + varName + ")";
    }

    Tree typeTree = io.getType();
    TypeMirror checkTm = typeOf(new TreePath(path, typeTree));
    if (resolver.isStringType(checkTm)) {
      return rt + "InstanceOfString(" + exprAny + ")";
    }
    if (checkTm instanceof DeclaredType dt
        && dt.asElement() instanceof TypeElement te
        && resolver.isNodeClass(te)
        && resolver.getNodeKindFlag(te) > 0) {
      int kf = resolver.getNodeKindFlag(te);
      String ns = nsPrefix(ctx);
      if (resolver.isNodeType(exprTm)) {
        return "(" + exprGo + " != 0 && " + ns + "Node(" + exprGo + ").KindFlags() == " + kf + ")";
      }
      String checkGo = resolver.goType(checkTm, ctx.imports);
      return "func() bool { n_, ok_ := "
          + exprAny
          + ".("
          + checkGo
          + "); return ok_ && n_ != 0 && "
          + ns
          + "Node(n_).KindFlags() == "
          + kf
          + " }()";
    }
    if (checkTm instanceof DeclaredType dt
        && dt.asElement() instanceof TypeElement te
        && resolver.isSourceType(te)
        && resolver.hasEffectiveTypeParameters(te)) {
      String marker =
          (te.getKind() == ElementKind.INTERFACE || resolver.isSubclassed(te))
              ? resolver.getInterfaceMarkerMethod(te)
              : resolver.getGenericClassMarkerMethod(te);
      return "func() bool { v_, ok_ := "
          + exprAny
          + ".(interface { "
          + marker
          + "(); IsJavaNil_() bool }); return ok_ && !v_.IsJavaNil_() }()";
    }
    String checkGo = resolver.goType(checkTm, ctx.imports);
    return rt + "InstanceOf[" + checkGo + "](" + exprAny + ")";
  }

  private String emitConditional(ConditionalExpressionTree ce, TreePath path, EmitContext ctx) {
    TypeMirror resTm = typeOf(path);
    String resGo = resolver.goType(resTm, ctx.imports);
    String cond =
        emitAs(ce.getCondition(), path, resolver.types.getPrimitiveType(TypeKind.BOOLEAN), false, ctx);
    String tVal = emitAs(ce.getTrueExpression(), path, resTm, false, ctx);
    String fVal = emitAs(ce.getFalseExpression(), path, resTm, false, ctx);
    return "func() " + resGo + " { if " + cond + " { return " + tVal + " }; return " + fVal + " }()";
  }

  // --- Arrays & Object Instantiation ---

  private String emitNewArray(NewArrayTree na, TreePath path, EmitContext ctx) {
    TypeMirror arrTm = typeOf(path);
    ArrayType at = (ArrayType) resolver.goEffectiveType(arrTm);
    TypeMirror compTm = at.getComponentType();
    String arrGo = resolver.goType(at, ctx.imports);

    if (na.getInitializers() != null) {
      List<String> elems = new ArrayList<>();
      for (ExpressionTree init : na.getInitializers()) {
        elems.add(emitAs(init, path, compTm, false, ctx));
      }
      return arrGo + "{" + String.join(", ", elems) + "}";
    }

    List<? extends ExpressionTree> dims = na.getDimensions();
    if (dims.size() == 1) {
      String lenExpr =
          emitAs(dims.get(0), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
      if (resolver.isStringType(compTm)) {
        return rtPrefix(ctx) + "NewStringArray(" + lenExpr + ")";
      }
      return "make(" + arrGo + ", " + lenExpr + ")";
    }

    if (dims.size() == 2) {
      String d1 =
          emitAs(dims.get(0), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
      String d2 =
          emitAs(dims.get(1), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
      ArrayType innerAt = (ArrayType) compTm;
      TypeMirror leafComp = innerAt.getComponentType();
      String leafGo = resolver.goType(leafComp, ctx.imports);
      String initFn =
          resolver.isStringType(leafComp)
              ? ("func() string { return " + rtPrefix(ctx) + "NullString }")
              : "nil";
      return rtPrefix(ctx) + "Make2D[" + leafGo + "](" + d1 + ", " + d2 + ", " + initFn + ")";
    }

    String d0 =
        emitAs(dims.get(0), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
    return "make(" + arrGo + ", " + d0 + ")";
  }

  private String emitNewClass(NewClassTree nc, TreePath path, EmitContext ctx) {
    Element el = resolver.trees.getElement(path);
    ExecutableElement ctor = (el instanceof ExecutableElement ee) ? ee : null;

    // Anonymous class `new Base(args) { ... }`
    if (nc.getClassBody() != null) {
      ClassTree body = nc.getClassBody();
      TreePath bodyPath = new TreePath(path, body);
      TypeElement anonTe = resolver.getClassTreeElement(body);
      if (anonTe == null && resolver.trees.getElement(bodyPath) instanceof TypeElement t) {
        anonTe = t;
      }
      ctx.deferLocalOrAnonClass(body, bodyPath, anonTe);

      List<String> callArgs = new ArrayList<>();
      if (resolver.hasOuterInstance(anonTe)) {
        TypeElement outerTe = resolver.getEnclosingInstanceType(anonTe);
        callArgs.add(resolveImplicitThisValue(outerTe, ctx));
      }
      TypeElement anonSuperTe = anonTe != null ? resolver.getSuperclassElement(anonTe) : null;
      if (anonSuperTe != null
          && resolver.isSourceType(anonSuperTe)
          && resolver.hasOuterInstance(anonSuperTe)) {
        TypeElement superOuterTe = resolver.getEnclosingInstanceType(anonSuperTe);
        if (nc.getEnclosingExpression() != null) {
          TypeMirror superOuterTm =
              (anonTe.getSuperclass() instanceof DeclaredType sdt
                      && sdt.getEnclosingType() instanceof DeclaredType sedt)
                  ? sedt
                  : superOuterTe.asType();
          callArgs.add(emitAs(nc.getEnclosingExpression(), path, superOuterTm, false, ctx));
        } else {
          callArgs.add(resolveImplicitThisValue(superOuterTe, ctx));
        }
      }
      for (VariableElement cap : resolver.getCapturedVars(anonTe)) {
        String loc = ctx.lookupLocalVar(cap);
        callArgs.add(loc != null ? loc : GoNames.localIdent(cap.getSimpleName().toString()));
      }
      if (ctor != null) {
        TypeMirror anonSuperSiteTm =
            (anonTe != null && anonSuperTe != null) ? anonTe.getSuperclass() : null;
        appendFormattedExecArgs(
            nc.getArguments(),
            path,
            ctor,
            (ExecutableType) ctor.asType(),
            null,
            callArgs,
            ctx,
            anonSuperSiteTm,
            null);
      } else {
        for (ExpressionTree arg : nc.getArguments()) {
          callArgs.add(emit(arg, path, ctx));
        }
      }
      String anonGoName = resolver.planner.getGoTypeName(anonTe);
      String anonUseParams =
          resolver.formatGoTypeParamUse(resolver.getEffectiveTypeParameters(anonTe));
      return "New" + anonGoName + anonUseParams + "(" + String.join(", ", callArgs) + ")";
    }

    if (ctor == null) {
      return "nil";
    }

    TypeElement owner = (TypeElement) ctor.getEnclosingElement();
    String qn = owner.getQualifiedName().toString();

    // Special `new String(...)` constructors
    if (qn.equals("java.lang.String")) {
      return emitNewString(nc, path, ctor, ctx);
    }

    // Shim exceptions (`new IllegalStateException(...)`, `new IOException(...)`, etc.)
    if (resolver.runtime.isShimException(owner)) {
      List<String> args = new ArrayList<>();
      args.add(GoNames.goStringLiteral(qn));
      for (ExpressionTree arg : nc.getArguments()) {
        args.add(emitAs(arg, path, resolver.getObjectElement().asType(), true, ctx));
      }
      String exIface = resolver.goTypeForElement(owner, ctx.imports);
      String rt = rtPrefix(ctx);
      return rt + "AsIface[" + exIface + "](" + rt + "ConstructShimThrowable(" + String.join(", ", args) + "))";
    }

    // Boxed primitive constructors (rare)
    if (resolver.runtime.isBoxedPrimitive(owner) && nc.getArguments().size() == 1) {
      return emitAs(nc.getArguments().get(0), path, owner.asType(), false, ctx);
    }

    TypeMirror newSiteTm = typeOf(path);
    List<String> callArgs = new ArrayList<>();
    if (resolver.isSourceType(owner)) {
      if (resolver.hasOuterInstance(owner)) {
        TypeElement outerTe = resolver.getEnclosingInstanceType(owner);
        if (nc.getEnclosingExpression() != null) {
          TypeMirror outerSiteTm =
              (newSiteTm instanceof DeclaredType ndt
                      && ndt.getEnclosingType() instanceof DeclaredType encDt)
                  ? encDt
                  : outerTe.asType();
          callArgs.add(emitAs(nc.getEnclosingExpression(), path, outerSiteTm, false, ctx));
        } else {
          callArgs.add(resolveImplicitThisValue(outerTe, ctx));
        }
      }
      for (VariableElement cap : resolver.getCapturedVars(owner)) {
        String loc = ctx.lookupLocalVar(cap);
        callArgs.add(loc != null ? loc : GoNames.localIdent(cap.getSimpleName().toString()));
      }
    }

    appendFormattedExecArgs(
        nc.getArguments(),
        path,
        ctor,
        (ExecutableType) ctor.asType(),
        null,
        callArgs,
        ctx,
        newSiteTm,
        null);

    String ctorFn = resolver.getConstructorGoName(ctor);
    String ctorTypeArgs =
        resolver.hasEffectiveTypeParameters(owner)
            ? resolver.formatDeclaredTypeArgs(newSiteTm, owner, ctx.imports)
            : "";
    if (resolver.isSourceType(owner)) {
      String targetPkg = resolver.planner.getGoPkg(owner);
      String alias = ctx.imports.ext(resolver.planner.getGoImportPath(targetPkg), targetPkg);
      String prefix = alias.isEmpty() ? "" : (alias + ".");
      return prefix + ctorFn + ctorTypeArgs + "(" + String.join(", ", callArgs) + ")";
    } else {
      return rtPrefix(ctx) + ctorFn + ctorTypeArgs + "(" + String.join(", ", callArgs) + ")";
    }
  }

  void emitNewClassInitArgs(
      NewClassTree nc,
      TreePath ncPath,
      ExecutableElement ctor,
      List<String> outArgs,
      EmitContext ctx) {
    appendFormattedExecArgs(
        nc.getArguments(),
        ncPath,
        ctor,
        (ExecutableType) ctor.asType(),
        null,
        outArgs,
        ctx,
        typeOf(ncPath),
        null);
  }

  private String emitNewString(
      NewClassTree nc, TreePath path, ExecutableElement ctor, EmitContext ctx) {
    List<? extends ExpressionTree> args = nc.getArguments();
    if (args.isEmpty()) {
      return "\"\"";
    }
    if (args.size() == 1) {
      TypeMirror p0 = resolver.safeErasure(ctor.getParameters().get(0).asType());
      if (p0 instanceof ArrayType at && at.getComponentType().getKind() == TypeKind.CHAR) {
        return jlPrefix(ctx) + "StringFromUTF16(" + emit(args.get(0), path, ctx) + ")";
      }
      if (p0 instanceof ArrayType at && at.getComponentType().getKind() == TypeKind.BYTE) {
        String bTmp = ctx.nextTemp("b");
        return "func() string { "
            + bTmp
            + " := "
            + emit(args.get(0), path, ctx)
            + "; return "
            + rtPrefix(ctx)
            + "NewStringFromBytes("
            + bTmp
            + ", 0, int32(len("
            + bTmp
            + ")), nil) }()";
      }
      return rtPrefix(ctx) + "StrVal(" + emit(args.get(0), path, ctx) + ")";
    }
    if (args.size() == 2) {
      TypeMirror p0 = resolver.safeErasure(ctor.getParameters().get(0).asType());
      if (p0 instanceof ArrayType at && at.getComponentType().getKind() == TypeKind.BYTE) {
        String bTmp = ctx.nextTemp("b");
        String csExpr = emit(args.get(1), path, ctx);
        return "func() string { "
            + bTmp
            + " := "
            + emit(args.get(0), path, ctx)
            + "; return "
            + rtPrefix(ctx)
            + "NewStringFromBytes("
            + bTmp
            + ", 0, int32(len("
            + bTmp
            + ")), "
            + csExpr
            + ") }()";
      }
    }
    if (args.size() == 3) {
      TypeMirror p0 = resolver.safeErasure(ctor.getParameters().get(0).asType());
      if (p0 instanceof ArrayType at && at.getComponentType().getKind() == TypeKind.CHAR) {
        String arr = emit(args.get(0), path, ctx);
        String off = emitAs(args.get(1), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
        String cnt = emitAs(args.get(2), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
        String aTmp = ctx.nextTemp("ca");
        String oTmp = ctx.nextTemp("co");
        String cTmp = ctx.nextTemp("cc");
        return "func() string { "
            + aTmp
            + " := "
            + arr
            + "; "
            + oTmp
            + " := "
            + off
            + "; "
            + cTmp
            + " := "
            + cnt
            + "; return "
            + jlPrefix(ctx)
            + "StringFromUTF16("
            + aTmp
            + "["
            + oTmp
            + " : "
            + oTmp
            + "+"
            + cTmp
            + "]) }()";
      }
      if (p0 instanceof ArrayType at && at.getComponentType().getKind() == TypeKind.INT) {
        String arr = emit(args.get(0), path, ctx);
        String off = emitAs(args.get(1), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
        String cnt = emitAs(args.get(2), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
        return jlPrefix(ctx) + "StringFromCodePoints(" + arr + ", " + off + ", " + cnt + ")";
      }
      if (p0 instanceof ArrayType at && at.getComponentType().getKind() == TypeKind.BYTE) {
        String arr = emit(args.get(0), path, ctx);
        String off = emitAs(args.get(1), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
        String cnt = emitAs(args.get(2), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
        return rtPrefix(ctx) + "NewStringFromBytes(" + arr + ", " + off + ", " + cnt + ", nil)";
      }
    }
    if (args.size() == 4) {
      TypeMirror p0 = resolver.safeErasure(ctor.getParameters().get(0).asType());
      if (p0 instanceof ArrayType at && at.getComponentType().getKind() == TypeKind.BYTE) {
        String arr = emit(args.get(0), path, ctx);
        String off = emitAs(args.get(1), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
        String cnt = emitAs(args.get(2), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
        String cs = emit(args.get(3), path, ctx);
        return rtPrefix(ctx) + "NewStringFromBytes(" + arr + ", " + off + ", " + cnt + ", " + cs + ")";
      }
    }
    return rtPrefix(ctx) + "StrVal(" + emit(args.get(0), path, ctx) + ")";
  }

  // --- Method Invocations ---

  private String emitMethodInvocation(MethodInvocationTree mi, TreePath path, EmitContext ctx) {
    ExpressionTree sel = mi.getMethodSelect();
    Element el = resolver.trees.getElement(path);
    if (!(el instanceof ExecutableElement ee)) {
      return "nil";
    }

    TypeElement owner = (TypeElement) ee.getEnclosingElement();
    boolean isStatic = ee.getModifiers().contains(Modifier.STATIC);
    boolean isSuperCall = false;
    TreePath recvPath = null;
    TypeMirror recvType = null;

    if (sel instanceof MemberSelectTree ms) {
      ExpressionTree recvTree = ms.getExpression();
      recvPath = new TreePath(new TreePath(path, sel), recvTree);
      recvType = typeOf(recvPath);
      if (recvTree instanceof IdentifierTree id && id.getName().contentEquals("super")) {
        isSuperCall = true;
      } else if (recvTree instanceof MemberSelectTree rms
          && rms.getIdentifier().contentEquals("super")) {
        isSuperCall = true;
      }
    }

    if (!isStatic && !isSuperCall && recvPath != null) {
      Map<VariableElement, String> entryCursors = activeEntryCursors.get();
      if (!entryCursors.isEmpty()) {
        TreePath strippedRecv = unwrapParens(recvPath);
        if (strippedRecv.getLeaf() instanceof IdentifierTree
            && resolver.trees.getElement(strippedRecv) instanceof VariableElement recvVe) {
          String curTmp = entryCursors.get(recvVe);
          if (curTmp != null) {
            String mName = ee.getSimpleName().toString();
            TypeMirror callSiteRet = typeOf(path);
            TypeMirror objTm = resolver.getObjectElement().asType();
            boolean isStmt = isDirectExprStmt(path);
            if (mName.equals("getKey") && mi.getArguments().isEmpty()) {
              if (!isStmt && resolver.isStringType(callSiteRet)) {
                return curTmp + ".KeyStr()";
              }
              if (!isStmt && resolver.isNodeType(callSiteRet)) {
                String nodeGo = resolver.goType(resolver.safeErasure(callSiteRet), ctx.imports);
                return nodeGo + "(" + curTmp + ".KeyNode())";
              }
              String raw = curTmp + ".GetKey()";
              return (!isStmt && callSiteRet.getKind() != TypeKind.VOID)
                  ? convert(raw, objTm, callSiteRet, false, ctx)
                  : raw;
            }
            if (mName.equals("getValue") && mi.getArguments().isEmpty()) {
              String raw = curTmp + ".GetValue()";
              return (!isStmt && callSiteRet.getKind() != TypeKind.VOID)
                  ? convert(raw, objTm, callSiteRet, false, ctx)
                  : raw;
            }
            if (mName.equals("setValue") && mi.getArguments().size() == 1) {
              String nvGo = emitAs(mi.getArguments().get(0), path, objTm, true, ctx);
              String raw = curTmp + ".SetValue(" + nvGo + ")";
              return (!isStmt && callSiteRet.getKind() != TypeKind.VOID)
                  ? convert(raw, objTm, callSiteRet, false, ctx)
                  : raw;
            }
          }
        }
      }
    }

    // Check intrinsic bindings for JDK / Guava / Enum / Object / String / Preconditions / Math
    String intrinsic = tryEmitIntrinsicCall(mi, path, ee, owner, isStatic, isSuperCall, recvPath, ctx);
    if (intrinsic != null) {
      return intrinsic;
    }

    TypeMirror callSiteRet = typeOf(path);

    if (isStatic) {
      String mName = ee.getSimpleName().toString();
      if (owner.getQualifiedName().contentEquals("java.util.Arrays")
          && mName.equals("copyOf")
          && mi.getArguments().size() == 2) {
        TreePath arg0Path = new TreePath(path, mi.getArguments().get(0));
        TypeMirror srcArrEff = resolver.goEffectiveType(typeOf(arg0Path));
        TypeMirror dstArrEff = resolver.goEffectiveType(callSiteRet);
        if (srcArrEff instanceof ArrayType sat && dstArrEff instanceof ArrayType dat) {
          String srcGo = resolver.goType(sat, ctx.imports);
          String dstGo = resolver.goType(dat, ctx.imports);
          if (srcGo.equals(dstGo)) {
            TypeMirror compTm = dat.getComponentType();
            String sTmp = ctx.nextTemp("arrSrc");
            String nTmp = ctx.nextTemp("arrLen");
            String oTmp = ctx.nextTemp("arrOut");
            String arg0Go = emitPath(arg0Path, ctx);
            String lenGo =
                emitAs(
                    mi.getArguments().get(1),
                    path,
                    resolver.types.getPrimitiveType(TypeKind.INT),
                    false,
                    ctx);
            String rt = rtPrefix(ctx);
            StringBuilder cb = new StringBuilder();
            cb.append("func() ").append(dstGo).append(" { ");
            cb.append(sTmp).append(" := ").append(arg0Go).append("; ");
            cb.append(nTmp).append(" := ").append(lenGo).append("; ");
            cb.append(rt).append("CheckNotNullSlice(").append(sTmp).append("); ");
            cb.append("if ").append(nTmp).append(" < 0 { panic(")
                .append(rt)
                .append("ConstructShimThrowable(\"java.lang.NegativeArraySizeException\")) }; ");
            if (resolver.isStringType(compTm)) {
              cb.append(oTmp).append(" := ").append(rt).append("NewStringArray(").append(nTmp).append("); ");
            } else if (resolver.isActiveTypeVariable(compTm)) {
              String compGo = resolver.goType(compTm, ctx.imports);
              cb.append(oTmp).append(" := make(").append(dstGo).append(", ").append(nTmp).append("); ");
              cb.append("if int(").append(nTmp).append(") > len(").append(sTmp)
                  .append(") { if _, ok_ := any((*").append(compGo).append(")(nil)).(*string); ok_ { ns_ := any(")
                  .append(rt).append("NullString).(").append(compGo).append("); for i_ := len(")
                  .append(sTmp).append("); i_ < int(").append(nTmp).append("); i_++ { ")
                  .append(oTmp).append("[i_] = ns_ } } }; ");
            } else {
              cb.append(oTmp).append(" := make(").append(dstGo).append(", ").append(nTmp).append("); ");
            }
            cb.append("copy(").append(oTmp).append(", ").append(sTmp).append("); return ").append(oTmp).append(" }()");
            return cb.toString();
          }
        }
      }
      if (owner.getQualifiedName().contentEquals("java.util.Arrays")
          && mName.equals("copyOfRange")
          && mi.getArguments().size() == 3) {
        TreePath arg0Path = new TreePath(path, mi.getArguments().get(0));
        TypeMirror srcArrEff = resolver.goEffectiveType(typeOf(arg0Path));
        TypeMirror dstArrEff = resolver.goEffectiveType(callSiteRet);
        if (srcArrEff instanceof ArrayType sat && dstArrEff instanceof ArrayType dat) {
          String srcGo = resolver.goType(sat, ctx.imports);
          String dstGo = resolver.goType(dat, ctx.imports);
          if (srcGo.equals(dstGo)) {
            TypeMirror compTm = dat.getComponentType();
            String sTmp = ctx.nextTemp("arrSrc");
            String fTmp = ctx.nextTemp("arrFrom");
            String tTmp = ctx.nextTemp("arrTo");
            String nTmp = ctx.nextTemp("arrLen");
            String oTmp = ctx.nextTemp("arrOut");
            String arg0Go = emitPath(arg0Path, ctx);
            String fromGo =
                emitAs(
                    mi.getArguments().get(1),
                    path,
                    resolver.types.getPrimitiveType(TypeKind.INT),
                    false,
                    ctx);
            String toGo =
                emitAs(
                    mi.getArguments().get(2),
                    path,
                    resolver.types.getPrimitiveType(TypeKind.INT),
                    false,
                    ctx);
            String rt = rtPrefix(ctx);
            String jl = jlPrefix(ctx);
            StringBuilder cb = new StringBuilder();
            cb.append("func() ").append(dstGo).append(" { ");
            cb.append(sTmp).append(" := ").append(arg0Go).append("; ");
            cb.append(fTmp).append(" := ").append(fromGo).append("; ");
            cb.append(tTmp).append(" := ").append(toGo).append("; ");
            cb.append(rt).append("CheckNotNullSlice(").append(sTmp).append("); ");
            cb.append("if ").append(fTmp).append(" > ").append(tTmp)
                .append(" { panic(").append(jl).append("NewIllegalArgumentException(\"fromIndex > toIndex\")) }; ");
            cb.append("if ").append(fTmp).append(" < 0 || int(").append(fTmp).append(") > len(").append(sTmp)
                .append(") { panic(").append(jl).append("NewArrayIndexOutOfBoundsException(\"\")) }; ");
            cb.append(nTmp).append(" := ").append(tTmp).append(" - ").append(fTmp).append("; ");
            if (resolver.isStringType(compTm)) {
              cb.append(oTmp).append(" := ").append(rt).append("NewStringArray(").append(nTmp).append("); ");
            } else if (resolver.isActiveTypeVariable(compTm)) {
              String compGo = resolver.goType(compTm, ctx.imports);
              cb.append(oTmp).append(" := make(").append(dstGo).append(", ").append(nTmp).append("); ");
              cb.append("if int(").append(tTmp).append(") > len(").append(sTmp)
                  .append(") { if _, ok_ := any((*").append(compGo).append(")(nil)).(*string); ok_ { ns_ := any(")
                  .append(rt).append("NullString).(").append(compGo).append("); for i_ := len(")
                  .append(sTmp).append(") - int(").append(fTmp).append("); i_ < int(").append(nTmp).append("); i_++ { ")
                  .append(oTmp).append("[i_] = ns_ } } }; ");
            } else {
              cb.append(oTmp).append(" := make(").append(dstGo).append(", ").append(nTmp).append("); ");
            }
            cb.append("copy(").append(oTmp).append(", ").append(sTmp).append("[").append(fTmp).append(":]); return ").append(oTmp).append(" }()");
            return cb.toString();
          }
        }
      }
      if (owner.getQualifiedName().contentEquals("java.util.Arrays")
          && mName.equals("toString")
          && mi.getArguments().size() == 1) {
        TreePath arg0Path = new TreePath(path, mi.getArguments().get(0));
        TypeMirror srcArrEff = resolver.goEffectiveType(typeOf(arg0Path));
        if (srcArrEff instanceof ArrayType sat
            && sat.getComponentType().getKind() != TypeKind.CHAR) {
          return rtPrefix(ctx) + "SliceFormatJavaString(" + emitPath(arg0Path, ctx) + ")";
        }
      }
      if (owner.getQualifiedName().contentEquals("java.util.Arrays")
          && mName.equals("fill")
          && (mi.getArguments().size() == 2 || mi.getArguments().size() == 4)
          && !ee.getParameters().isEmpty()) {
        TreePath arg0Path = new TreePath(path, mi.getArguments().get(0));
        TypeMirror arg0Eff = resolver.goEffectiveType(typeOf(arg0Path));
        if (arg0Eff instanceof ArrayType sat) {
          TypeMirror compTm = sat.getComponentType();
          String arg0Go = emitPath(arg0Path, ctx);
          TypeMirror intTm = resolver.types.getPrimitiveType(TypeKind.INT);
          if (mi.getArguments().size() == 2) {
            String valGo = emitAs(mi.getArguments().get(1), path, compTm, false, ctx);
            if (valGo.equals("nil")) {
              return rtPrefix(ctx) + "ArraysFillZero(" + arg0Go + ")";
            }
            return rtPrefix(ctx) + "ArraysFill(" + arg0Go + ", " + valGo + ")";
          }
          String fromGo = emitAs(mi.getArguments().get(1), path, intTm, false, ctx);
          String toGo = emitAs(mi.getArguments().get(2), path, intTm, false, ctx);
          String valGo = emitAs(mi.getArguments().get(3), path, compTm, false, ctx);
          if (valGo.equals("nil")) {
            return rtPrefix(ctx)
                + "ArraysFillRangeZero("
                + arg0Go
                + ", "
                + fromGo
                + ", "
                + toGo
                + ")";
          }
          return rtPrefix(ctx)
              + "ArraysFillRange("
              + arg0Go
              + ", "
              + fromGo
              + ", "
              + toGo
              + ", "
              + valGo
              + ")";
        }
      }
      if (owner.getQualifiedName().contentEquals("java.util.Arrays")
          && (mName.equals("setAll")
              || mName.equals("parallelSetAll")
              || mName.equals("fill")
              || mName.equals("sort")
              || mName.equals("parallelSort"))
          && !mi.getArguments().isEmpty()
          && !ee.getParameters().isEmpty()) {
        TypeMirror p0Erased = resolver.safeErasure(ee.getParameters().get(0).asType());
        TreePath arg0Path = new TreePath(path, mi.getArguments().get(0));
        TypeMirror arg0Erased = resolver.goEffectiveType(typeOf(arg0Path));
        if (p0Erased instanceof ArrayType dat && arg0Erased instanceof ArrayType sat) {
          String dstSliceGo = resolver.goType(dat, ctx.imports);
          String srcSliceGo = resolver.goType(sat, ctx.imports);
          if (!dstSliceGo.equals(srcSliceGo)) {
            String srcTmp = ctx.nextTemp("arrMut");
            String boxTmp = ctx.nextTemp("arrBox");
            String vTmp = ctx.nextTemp("el");
            String arg0Go = emitPath(arg0Path, ctx);
            List<String> callArgs = new ArrayList<>();
            callArgs.add(boxTmp);
            for (int i = 1; i < mi.getArguments().size(); i++) {
              TypeMirror tp = resolver.safeErasure(ee.getParameters().get(i).asType());
              callArgs.add(emitAs(mi.getArguments().get(i), path, tp, false, ctx));
            }
            String fnName = resolver.getStaticMethodGoName(ee);
            return "func() { "
                + srcTmp
                + " := "
                + arg0Go
                + "; var "
                + boxTmp
                + " "
                + dstSliceGo
                + "; if "
                + srcTmp
                + " != nil { "
                + boxTmp
                + " = make("
                + dstSliceGo
                + ", len("
                + srcTmp
                + ")); for i_, "
                + vTmp
                + " := range "
                + srcTmp
                + " { "
                + boxTmp
                + "[i_] = "
                + convert(vTmp, sat.getComponentType(), dat.getComponentType(), false, ctx)
                + " } }; "
                + rtPrefix(ctx)
                + fnName
                + "("
                + String.join(", ", callArgs)
                + "); for i_, "
                + vTmp
                + " := range "
                + boxTmp
                + " { "
                + srcTmp
                + "[i_] = "
                + convert(vTmp, dat.getComponentType(), sat.getComponentType(), false, ctx)
                + " } }()";
          }
        }
      }
      List<TypeMirror> inferredStaticArgs = null;
      String staticTypeArgSuffix = "";
      if (resolver.isGenericStaticMethod(ee)) {
        List<TypeMirror> explicitTas = new ArrayList<>();
        for (Tree ta : mi.getTypeArguments()) {
          explicitTas.add(typeOf(new TreePath(path, ta)));
        }
        List<TypeMirror> actualArgTms = new ArrayList<>(mi.getArguments().size());
        for (ExpressionTree argTree : mi.getArguments()) {
          actualArgTms.add(typeOf(new TreePath(path, argTree)));
        }
        inferredStaticArgs =
            resolver.inferGenericStaticMethodTypeArgs(
                ee, explicitTas, actualArgTms, callSiteRet);
        List<String> goTas = new ArrayList<>(inferredStaticArgs.size());
        for (TypeMirror ta : inferredStaticArgs) {
          String gta = resolver.goType(ta, ctx.imports);
          goTas.add(gta.isEmpty() ? "any" : gta);
        }
        staticTypeArgSuffix = "[" + String.join(", ", goTas) + "]";
      }
      List<String> args = new ArrayList<>();
      appendFormattedExecArgs(
          mi.getArguments(),
          path,
          ee,
          (ExecutableType) ee.asType(),
          null,
          args,
          ctx,
          null,
          inferredStaticArgs);
      String fnName = resolver.getStaticMethodGoName(ee) + staticTypeArgSuffix;
      String rawCall;
      if (resolver.isSourceType(owner)) {
        String ownerPkg = resolver.planner.getGoPkg(owner);
        String alias = ctx.imports.ext(resolver.planner.getGoImportPath(ownerPkg), ownerPkg);
        String prefix = alias.isEmpty() ? "" : (alias + ".");
        rawCall = prefix + fnName + "(" + String.join(", ", args) + ")";
      } else {
        rawCall = rtPrefix(ctx) + fnName + "(" + String.join(", ", args) + ")";
      }
      if (callSiteRet.getKind() != TypeKind.VOID) {
        TypeMirror rawRetTm =
            inferredStaticArgs != null
                ? resolver.substituteTypeParams(
                    ee.getReturnType(), ee.getTypeParameters(), inferredStaticArgs)
                : ee.getReturnType();
        return convert(rawCall, resolver.goEffectiveType(rawRetTm), callSiteRet, false, ctx);
      }
      return rawCall;
    }

    // Instance / interface method call
    TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
    ExecutableType execType = resolver.resolveExecutableType(recvType, ee);
    TypeMirror effectiveRecvTm = recvType != null ? recvType : resolveImplicitOwnerSiteType(owner, ctx);
    List<String> args = new ArrayList<>();
    appendFormattedExecArgs(
        mi.getArguments(), path, ee, execType, sig, args, ctx, effectiveRecvTm, null);

    String recvGo;
    boolean isDirectOwnPrivateCall =
        !isSuperCall
            && ee.getModifiers().contains(Modifier.PRIVATE)
            && owner.equals(ctx.currentType)
            && (recvPath == null
                || isDirectCurrentReceiverExpr(
                    (ExpressionTree) recvPath.getLeaf(), recvPath.getParentPath(), ctx));
    if (isSuperCall) {
      if (owner.getQualifiedName().contentEquals("java.lang.Object")) {
        String mName = ee.getSimpleName().toString();
        String selfObj = emitImplicitSelfAsAny(ctx);
        if (mName.equals("equals") && mi.getArguments().size() == 1) {
          String a0 =
              emitAs(
                  mi.getArguments().get(0),
                  path,
                  resolver.getObjectElement().asType(),
                  true,
                  ctx);
          return rtPrefix(ctx) + "RefEquals(" + selfObj + ", " + a0 + ")";
        }
        if (mName.equals("hashCode") && mi.getArguments().isEmpty()) {
          return jlPrefix(ctx) + "IdentityHashCode(" + selfObj + ")";
        }
        if (mName.equals("toString") && mi.getArguments().isEmpty()) {
          return jlPrefix(ctx)
              + "ObjectToString("
              + rtPrefix(ctx)
              + "GetObjectClass("
              + selfObj
              + ").GetName(), "
              + selfObj
              + ")";
        }
      }
      if (owner.getKind() == ElementKind.INTERFACE && resolver.isSourceType(owner)) {
        String ownerPkg = resolver.planner.getGoPkg(owner);
        String alias = ctx.imports.ext(resolver.planner.getGoImportPath(ownerPkg), ownerPkg);
        String prefix = alias.isEmpty() ? "" : (alias + ".");
        TypeMirror ownerSiteTm = resolveImplicitOwnerSiteType(owner, ctx);
        TypeMirror ifaceSiteTm = ownerSiteTm;
        if (ownerSiteTm instanceof DeclaredType odt && !owner.equals(odt.asElement())) {
          DeclaredType sup = resolver.findSuperType(odt, owner);
          if (sup != null) {
            ifaceSiteTm = sup;
          }
        }
        String ifaceGo = resolver.goType(ifaceSiteTm, ctx.imports);
        String selfIface = rtPrefix(ctx) + "AsIface[" + ifaceGo + "](" + resolveImplicitThisValue(owner, ctx) + ")";
        List<String> allArgs = new ArrayList<>();
        allArgs.add(selfIface);
        allArgs.addAll(args);
        String defaultTypeArgs = resolver.formatSuperTypeArgs(ifaceSiteTm, owner, ctx.imports);
        String rawCall =
            prefix
                + "IfaceDefault_"
                + resolver.planner.getGoTypeName(owner)
                + "_"
                + sig.goName
                + defaultTypeArgs
                + "("
                + String.join(", ", allArgs)
                + ")";
        TypeMirror defRetTm = resolver.resolveSigReturnType(sig, ifaceSiteTm);
        return (callSiteRet.getKind() != TypeKind.VOID && defRetTm != null)
            ? convert(rawCall, defRetTm, callSiteRet, false, ctx)
            : rawCall;
      }
      if (ctx.currentType != null && resolver.isValueEnum(ctx.currentType)) {
        String mName = ee.getSimpleName().toString();
        if ((mName.equals("name") || mName.equals("toString")) && args.isEmpty()) {
          return ctx.receiverName + ".Name()";
        }
        if (mName.equals("ordinal") && args.isEmpty()) {
          return ctx.receiverName + ".Ordinal()";
        }
        if (mName.equals("hashCode") && args.isEmpty()) {
          return ctx.receiverName + ".HashCode()";
        }
        if (mName.equals("equals") && args.size() == 1) {
          return ctx.receiverName + ".Equals(" + args.get(0) + ")";
        }
        if (mName.equals("compareTo") && args.size() == 1) {
          return ctx.receiverName + ".CompareTo(" + args.get(0) + ")";
        }
        if (mName.equals("getDeclaringClass") && args.isEmpty()) {
          return rtPrefix(ctx) + "GetObjectClass(" + ctx.receiverName + ")";
        }
      }
      if (ctx.currentType != null && ctx.currentType.getKind() == ElementKind.ENUM) {
        String rawCall =
            ctx.receiverName + ".baseEnum_" + sig.goName + "(" + String.join(", ", args) + ")";
        TypeMirror enumRetTm = resolver.resolveSigReturnType(sig, ctx.currentType.asType());
        return enumRetTm != null
            ? convert(rawCall, enumRetTm, callSiteRet, false, ctx)
            : rawCall;
      }
      if (resolver.isNodeClass(ctx.currentType) && resolver.getNodeKindFlag(ctx.currentType) > 0) {
        // Node's 4 private inner classes calling `super.isEquivalentTo(...)`
        String rawCall =
            ctx.receiverName + ".base_" + sig.goName + "(" + String.join(", ", args) + ")";
        return sig.canonReturnType != null
            ? convert(rawCall, sig.canonReturnType, callSiteRet, false, ctx)
            : rawCall;
      }
      TypeElement targetClass = ctx.currentType;
      String structBase = ctx.receiverName;
      if (sel instanceof MemberSelectTree ms
          && ms.getExpression() instanceof MemberSelectTree rms
          && rms.getIdentifier().contentEquals("super")) {
        TreePath qPath = new TreePath(new TreePath(new TreePath(path, sel), rms), rms.getExpression());
        Element qEl = resolver.trees.getElement(qPath);
        if (qEl instanceof TypeElement qTe && qTe.getKind() == ElementKind.CLASS && !qTe.equals(ctx.currentType)) {
          targetClass = qTe;
          structBase = resolveImplicitThisStruct(qTe, ctx);
        }
      }
      TypeElement superTe = resolver.getSuperclassElement(targetClass);
      if (superTe != null) {
        String superStruct = resolver.goStructTypeForElement(superTe, ctx.imports);
        int dot = superStruct.lastIndexOf('.');
        String embeddedField = dot >= 0 ? superStruct.substring(dot + 1) : superStruct;
        recvGo = structBase + "." + embeddedField;
      } else {
        recvGo = structBase;
      }
    } else if (isDirectOwnPrivateCall) {
      recvGo = ctx.receiverName;
    } else if (ctx.currentType != null
        && !ctx.inStaticContext
        && owner.getKind() == ElementKind.CLASS
        && resolver.isSourceType(owner)
        && isSubtypeOf(ctx.currentType, owner)
        && resolver.canDispatchDirectlyOnBase(ctx.currentType, sig)
        && (recvPath == null
            || isDirectCurrentReceiverExpr(
                (ExpressionTree) recvPath.getLeaf(), recvPath.getParentPath(), ctx))) {
      recvGo = ctx.receiverName;
    } else if (ctx.currentType != null
        && !ctx.inStaticContext
        && resolver.isSubclassed(ctx.currentType)
        && !ee.getModifiers().contains(Modifier.PRIVATE)
        && isSubtypeOf(ctx.currentType, owner)
        && (recvPath == null
            || isDirectCurrentThisExpr(
                (ExpressionTree) recvPath.getLeaf(), recvPath.getParentPath(), ctx))) {
      TypeElement reusableSelfAnc =
          resolver.findReusableSelfAncestorForCall(ctx.currentType, owner, sig);
      if (reusableSelfAnc != null) {
        recvGo =
            ctx.receiverName
                + "."
                + resolver.planner.getGoTypeName(reusableSelfAnc)
                + "Base.self_";
        if (ctx.currentType.asType() instanceof DeclaredType curDt) {
          DeclaredType ancSite = resolver.findSuperType(curDt, reusableSelfAnc);
          if (ancSite != null) {
            effectiveRecvTm = ancSite;
          }
        }
      } else {
        recvGo = recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
      }
    } else if (recvPath != null) {
      if (ee.getModifiers().contains(Modifier.PRIVATE)
          && resolver.isSourceType(owner)
          && !resolver.isNativeGenericClass(owner)) {
        recvGo =
            emitAs(
                (ExpressionTree) recvPath.getLeaf(),
                recvPath.getParentPath(),
                owner.asType(),
                false,
                ctx);
      } else {
        recvGo = emitPath(recvPath, ctx);
      }
    } else {
      // Unqualified instance call: dispatch on dynamic `self_` (or `outer_`) unless private
      recvGo = resolveImplicitThisValue(owner, ctx);
    }

    String rawCall;
    if (owner.getKind() == ElementKind.INTERFACE && ee.getModifiers().contains(Modifier.PRIVATE)) {
      String ownerPkg = resolver.planner.getGoPkg(owner);
      String alias = ctx.imports.ext(resolver.planner.getGoImportPath(ownerPkg), ownerPkg);
      String prefix = alias.isEmpty() ? "" : (alias + ".");
      List<String> allArgs = new ArrayList<>();
      allArgs.add(recvGo);
      allArgs.addAll(args);
      String privTypeArgs = resolver.formatSuperTypeArgs(effectiveRecvTm, owner, ctx.imports);
      rawCall =
          prefix
              + "IfaceDefault_"
              + resolver.planner.getGoTypeName(owner)
              + "_"
              + sig.goName
              + privTypeArgs
              + "("
              + String.join(", ", allArgs)
              + ")";
    } else if (ee.getModifiers().contains(Modifier.PRIVATE)
        && resolver.isSourceType(owner)
        && resolver.isSubclassed(owner)
        && !isDirectOwnPrivateCall) {
      rawCall =
          recvGo
              + ".Get"
              + resolver.planner.getGoTypeName(owner)
              + "Base()."
              + sig.goName
              + "("
              + String.join(", ", args)
              + ")";
    } else {
      rawCall = recvGo + "." + sig.goName + "(" + String.join(", ", args) + ")";
    }
    if (callSiteRet.getKind() != TypeKind.VOID && sig.canonReturnType != null) {
      TypeMirror actualRetTm = resolver.resolveSigReturnType(sig, effectiveRecvTm);
      return convert(rawCall, actualRetTm, callSiteRet, false, ctx);
    }
    return rawCall;
  }

  void appendFormattedExecArgs(
      List<? extends ExpressionTree> argTrees,
      TreePath callParentPath,
      ExecutableElement ee,
      ExecutableType execType,
      TypeResolver.MethodSig virtualSig,
      List<String> outArgs,
      EmitContext ctx) {
    appendFormattedExecArgs(
        argTrees, callParentPath, ee, execType, virtualSig, outArgs, ctx, null, null);
  }

  private void appendFormattedExecArgs(
      List<? extends ExpressionTree> argTrees,
      TreePath callParentPath,
      ExecutableElement ee,
      ExecutableType execType,
      TypeResolver.MethodSig virtualSig,
      List<String> outArgs,
      EmitContext ctx,
      TypeMirror ownerSiteTm,
      List<TypeMirror> genericMethodTypeArgs) {
    int paramCount = ee.getParameters().size();
    boolean varArgs = ee.isVarArgs() && paramCount > 0;
    int fixedCount = varArgs ? (paramCount - 1) : paramCount;
    TypeElement eeOwner = ee.getEnclosingElement() instanceof TypeElement t ? t : null;

    for (int i = 0; i < fixedCount && i < argTrees.size(); i++) {
      boolean forceAny = virtualSig != null && virtualSig.paramWidenedToAny[i];
      TypeMirror targetParam;
      if (virtualSig != null && (genericMethodTypeArgs == null || ee.getTypeParameters().isEmpty())) {
        targetParam = resolver.resolveSigParamType(virtualSig, i, ownerSiteTm);
      } else {
        TypeMirror baseParam =
            virtualSig != null
                ? virtualSig.canonParamTypes.get(i)
                : ee.getParameters().get(i).asType();
        if (genericMethodTypeArgs != null && !ee.getTypeParameters().isEmpty()) {
          baseParam =
              resolver.substituteTypeParams(
                  baseParam, ee.getTypeParameters(), genericMethodTypeArgs);
        } else if (ownerSiteTm != null
            && eeOwner != null
            && resolver.isNativeGenericClass(eeOwner)
            && !forceAny) {
          baseParam = resolver.substituteOwnerTypeArgs(baseParam, eeOwner, ownerSiteTm);
        }
        targetParam = resolver.goEffectiveType(baseParam);
      }
      outArgs.add(emitAs(argTrees.get(i), callParentPath, targetParam, forceAny, ctx));
    }

    if (varArgs) {
      boolean lastForceAny = virtualSig != null && virtualSig.paramWidenedToAny[paramCount - 1];
      TypeMirror effLastParam;
      if (virtualSig != null && (genericMethodTypeArgs == null || ee.getTypeParameters().isEmpty())) {
        effLastParam = resolver.resolveSigParamType(virtualSig, paramCount - 1, ownerSiteTm);
      } else {
        TypeMirror lastParamTm =
            virtualSig != null
                ? virtualSig.canonParamTypes.get(paramCount - 1)
                : ee.getParameters().get(paramCount - 1).asType();
        if (genericMethodTypeArgs != null && !ee.getTypeParameters().isEmpty()) {
          lastParamTm =
              resolver.substituteTypeParams(
                  lastParamTm, ee.getTypeParameters(), genericMethodTypeArgs);
        } else if (ownerSiteTm != null
            && eeOwner != null
            && resolver.isNativeGenericClass(eeOwner)
            && !lastForceAny) {
          lastParamTm = resolver.substituteOwnerTypeArgs(lastParamTm, eeOwner, ownerSiteTm);
        }
        effLastParam = resolver.goEffectiveType(lastParamTm);
      }
      ArrayType varArrType = (ArrayType) effLastParam;
      TypeMirror compTm = varArrType.getComponentType();

      // Check if single array argument passed directly
      if (argTrees.size() == paramCount) {
        TreePath lastArgPath = new TreePath(callParentPath, argTrees.get(paramCount - 1));
        TypeMirror lastArgTm = typeOf(lastArgPath);
        if (lastArgTm.getKind() == TypeKind.ARRAY) {
          outArgs.add(emitAs(argTrees.get(paramCount - 1), callParentPath, varArrType, false, ctx));
          return;
        }
        if (lastArgTm.getKind() == TypeKind.NULL && !resolver.isObjectOrErasedAny(compTm)) {
          outArgs.add("nil");
          return;
        }
      }

      List<String> packed = new ArrayList<>();
      for (int i = fixedCount; i < argTrees.size(); i++) {
        packed.add(emitAs(argTrees.get(i), callParentPath, compTm, false, ctx));
      }
      String sliceGo = resolver.goType(varArrType, ctx.imports);
      outArgs.add(sliceGo + "{" + String.join(", ", packed) + "}");
    }
  }

  // --- Intrinsic Fast-Path Mappings to transpiled/shims/* & j2gort ---

  private String tryEmitIntrinsicCall(
      MethodInvocationTree mi,
      TreePath path,
      ExecutableElement ee,
      TypeElement owner,
      boolean isStatic,
      boolean isSuperCall,
      TreePath recvPath,
      EmitContext ctx) {
    if (isSuperCall) {
      return null;
    }
    String qn = owner.getQualifiedName().toString();
    String m = ee.getSimpleName().toString();
    List<? extends ExpressionTree> args = mi.getArguments();
    int arity = args.size();
    String jl = jlPrefix(ctx);
    String rt = rtPrefix(ctx);

    // 1. Source or External Enum synthesized `values()` and `valueOf(String)` + `name()` / `ordinal()`
    if (owner.getKind() == ElementKind.ENUM && isStatic) {
      if (resolver.isSourceType(owner)) {
        String ownerPkg = resolver.planner.getGoPkg(owner);
        String alias = ctx.imports.ext(resolver.planner.getGoImportPath(ownerPkg), ownerPkg);
        String prefix = alias.isEmpty() ? "" : (alias + ".");
        String enumGo = resolver.planner.getGoTypeName(owner);
        if (m.equals("values") && arity == 0) {
          return prefix + enumGo + "_Values()";
        }
        if (m.equals("valueOf")
            && arity == 1
            && resolver.isStringType(ee.getParameters().get(0).asType())) {
          return prefix
              + enumGo
              + "_ValueOf("
              + emitAs(args.get(0), path, ee.getParameters().get(0).asType(), false, ctx)
              + ")";
        }
      }
    }
    if (qn.equals("java.lang.Enum") && !isStatic) {
      TypeMirror recvErased =
          recvPath != null
              ? resolver.safeErasure(typeOf(recvPath))
              : (ctx.currentType != null ? resolver.safeErasure(ctx.currentType.asType()) : null);
      String recv = recvPath != null ? emitPath(recvPath, ctx) : ctx.selfValueExpr();
      if (recvErased != null && resolver.isValueEnumType(recvErased)) {
        ExpressionTree recvTree =
            recvPath != null && recvPath.getLeaf() instanceof ExpressionTree et ? et : null;
        if (m.equals("name") && arity == 0) return "(" + recv + ").Name()";
        if (m.equals("ordinal") && arity == 0) return "(" + recv + ").Ordinal()";
        if (m.equals("compareTo") && arity == 1) {
          TreePath arg0Path = new TreePath(path, args.get(0));
          TypeMirror arg0Erased = resolver.safeErasure(typeOf(arg0Path));
          if (resolver.types.isSameType(recvErased, arg0Erased)) {
            String aGo = emitPath(arg0Path, ctx);
            if (isSimpleLocalOrParamIdent(arg0Path, aGo)) {
              return "((" + recv + ").Ordinal() - (" + aGo + ").Ordinal())";
            }
            String rTmp = ctx.nextTemp("re");
            String aTmp = ctx.nextTemp("ae");
            return "func() int32 { "
                + rTmp
                + " := "
                + recv
                + "; "
                + aTmp
                + " := "
                + aGo
                + "; return "
                + rTmp
                + ".Ordinal() - "
                + aTmp
                + ".Ordinal() }()";
          }
          return "("
              + recv
              + ").CompareTo("
              + emitAs(args.get(0), path, resolver.getObjectElement().asType(), true, ctx)
              + ")";
        }
        if (m.equals("toString") && arity == 0) {
          return emitCheckNotNull(recvTree, recv, recvErased, ctx) + ".String()";
        }
        if (m.equals("equals") && arity == 1) {
          TreePath arg0Path = new TreePath(path, args.get(0));
          TypeMirror arg0Erased = resolver.safeErasure(typeOf(arg0Path));
          if (arg0Erased.getKind() == TypeKind.NULL) {
            return "func() bool { _ = "
                + emitCheckNotNull(recvTree, recv, recvErased, ctx)
                + "; return false }()";
          }
          if (resolver.types.isSameType(recvErased, arg0Erased)) {
            String aGo = emitPath(arg0Path, ctx);
            if (isSimpleLocalOrParamIdent(arg0Path, aGo)) {
              return "(" + emitCheckNotNull(recvTree, recv, recvErased, ctx) + " == " + aGo + ")";
            }
            String rTmp = ctx.nextTemp("re");
            String aTmp = ctx.nextTemp("ae");
            return "func() bool { "
                + rTmp
                + " := "
                + recv
                + "; "
                + aTmp
                + " := "
                + aGo
                + "; if "
                + rTmp
                + " == 0 { "
                + rt
                + "ThrowNPE() }; return "
                + rTmp
                + " == "
                + aTmp
                + " }()";
          }
          return emitCheckNotNull(recvTree, recv, recvErased, ctx)
              + ".Equals("
              + emitAs(args.get(0), path, resolver.getObjectElement().asType(), true, ctx)
              + ")";
        }
        if (m.equals("hashCode") && arity == 0) {
          return emitCheckNotNull(recvTree, recv, recvErased, ctx) + ".HashCode()";
        }
      }
      if (m.equals("name") && arity == 0) return rt + "EnumName(" + recv + ")";
      if (m.equals("ordinal") && arity == 0) return rt + "EnumOrdinal(" + recv + ")";
      if (m.equals("compareTo") && arity == 1) {
        return rt + "EnumCompareTo(" + recv + ", " + emit(args.get(0), path, ctx) + ")";
      }
      if (m.equals("toString") && arity == 0) return recv + ".String()";
      if (m.equals("equals") && arity == 1) {
        return rt + "RefEquals(" + recv + ", " + emit(args.get(0), path, ctx) + ")";
      }
      if (m.equals("hashCode") && arity == 0) return recv + ".HashCode()";
    }

    // 2a. Value-struct `equals(Object)` overrides whose semantics match Go `==` (e.g. ColorId)
    if (!isStatic
        && !isSuperCall
        && m.equals("equals")
        && arity == 1
        && ee.getParameters().size() == 1
        && ee.getReturnType().getKind() == TypeKind.BOOLEAN
        && resolver.safeErasure(ee.getParameters().get(0).asType()) instanceof DeclaredType eqPdt
        && eqPdt.asElement() instanceof TypeElement eqPte
        && eqPte.getQualifiedName().contentEquals("java.lang.Object")) {
      TypeMirror recvErased =
          recvPath != null
              ? resolver.safeErasure(typeOf(recvPath))
              : (ctx.currentType != null ? resolver.safeErasure(ctx.currentType.asType()) : null);
      if (recvErased instanceof DeclaredType rdt
          && rdt.asElement() instanceof TypeElement rte
          && resolver.isValueStructClass(rte)
          && valueStructGoEqMatchesJavaEquals(rte)) {
        TreePath arg0Path = new TreePath(path, args.get(0));
        TypeMirror arg0Erased = resolver.safeErasure(typeOf(arg0Path));
        if (arg0Erased.getKind() == TypeKind.NULL
            || (arg0Erased instanceof DeclaredType adt && rte.equals(adt.asElement()))) {
          String rGo = recvPath != null ? emitPath(recvPath, ctx) : ctx.selfValueExpr();
          ExpressionTree recvTree =
              recvPath != null && recvPath.getLeaf() instanceof ExpressionTree et ? et : null;
          if (arg0Erased.getKind() == TypeKind.NULL) {
            return "func() bool { _ = "
                + emitCheckNotNull(recvTree, rGo, recvErased, ctx)
                + "; return false }()";
          }
          String aGo = emitPath(arg0Path, ctx);
          if (!resolver.isNullableValueStructClass(rte)
              || isSimpleLocalOrParamIdent(arg0Path, aGo)) {
            return "(" + emitCheckNotNull(recvTree, rGo, recvErased, ctx) + " == " + aGo + ")";
          }
          String rTmp = ctx.nextTemp("rv");
          String aTmp = ctx.nextTemp("av");
          return "func() bool { "
              + rTmp
              + " := "
              + rGo
              + "; "
              + aTmp
              + " := "
              + aGo
              + "; _ = "
              + emitCheckNotNull(null, rTmp, recvErased, ctx)
              + "; return "
              + rTmp
              + " == "
              + aTmp
              + " }()";
        }
      }
    }

    // 2. java.lang.Object methods (`equals`, `hashCode`, `toString`, `getClass`)
    if (qn.equals("java.lang.Object") && !isStatic) {
      TypeMirror recvErased =
          recvPath != null
              ? resolver.safeErasure(typeOf(recvPath))
              : (ctx.currentType != null ? resolver.safeErasure(ctx.currentType.asType()) : null);
      if (m.equals("equals") && arity == 1 && recvErased != null) {
        TreePath arg0Path = new TreePath(path, args.get(0));
        TypeMirror arg0Erased = resolver.safeErasure(typeOf(arg0Path));
        if (resolver.isNodeType(recvErased)
            && (resolver.isNodeType(arg0Erased) || arg0Erased.getKind() == TypeKind.NULL)) {
          String rNode = recvPath != null ? emitPath(recvPath, ctx) : ctx.selfValueExpr();
          String aNode =
              arg0Erased.getKind() == TypeKind.NULL
                  ? "0"
                  : convert(emitPath(arg0Path, ctx), typeOf(arg0Path), recvErased, false, ctx);
          if (arg0Erased.getKind() == TypeKind.NULL || isSimpleLocalOrParamIdent(arg0Path, aNode)) {
            return "(" + rt + "CheckNotNullNode(" + rNode + ") == " + aNode + ")";
          }
          String rTmp = ctx.nextTemp("rn");
          String aTmp = ctx.nextTemp("an");
          return "func() bool { "
              + rTmp
              + " := "
              + rNode
              + "; "
              + aTmp
              + " := "
              + aNode
              + "; return "
              + rt
              + "CheckNotNullNode("
              + rTmp
              + ") == "
              + aTmp
              + " }()";
        }
        if (recvErased instanceof DeclaredType rdt
            && rdt.asElement() instanceof TypeElement rte
            && (rte.getKind() == ElementKind.ENUM
                || hasIdentityEquals(rte)
                || valueStructGoEqMatchesJavaEquals(rte))
            && (arg0Erased.getKind() == TypeKind.NULL
                || (arg0Erased instanceof DeclaredType adt && rte.equals(adt.asElement())))) {
          String rGo = recvPath != null ? emitPath(recvPath, ctx) : ctx.selfValueExpr();
          if (resolver.isValueStructClass(rte)) {
            ExpressionTree recvTree =
                recvPath != null && recvPath.getLeaf() instanceof ExpressionTree et ? et : null;
            if (arg0Erased.getKind() == TypeKind.NULL) {
              return "func() bool { _ = "
                  + emitCheckNotNull(recvTree, rGo, recvErased, ctx)
                  + "; return false }()";
            }
            String aGo = emitPath(arg0Path, ctx);
            if (!resolver.isNullableValueStructClass(rte)
                || isSimpleLocalOrParamIdent(arg0Path, aGo)) {
              return "(" + emitCheckNotNull(recvTree, rGo, recvErased, ctx) + " == " + aGo + ")";
            }
            String rTmp = ctx.nextTemp("rv");
            String aTmp = ctx.nextTemp("av");
            return "func() bool { "
                + rTmp
                + " := "
                + rGo
                + "; "
                + aTmp
                + " := "
                + aGo
                + "; _ = "
                + emitCheckNotNull(null, rTmp, recvErased, ctx)
                + "; return "
                + rTmp
                + " == "
                + aTmp
                + " }()";
          }
          if (resolver.isValueEnum(rte)) {
            ExpressionTree recvTree =
                recvPath != null && recvPath.getLeaf() instanceof ExpressionTree et ? et : null;
            if (arg0Erased.getKind() == TypeKind.NULL) {
              return "func() bool { _ = "
                  + emitCheckNotNull(recvTree, rGo, recvErased, ctx)
                  + "; return false }()";
            }
            String aGo = emitPath(arg0Path, ctx);
            if (isSimpleLocalOrParamIdent(arg0Path, aGo)) {
              return "(" + emitCheckNotNull(recvTree, rGo, recvErased, ctx) + " == " + aGo + ")";
            }
            String rTmp = ctx.nextTemp("re");
            String aTmp = ctx.nextTemp("ae");
            return "func() bool { "
                + rTmp
                + " := "
                + rGo
                + "; "
                + aTmp
                + " := "
                + aGo
                + "; if "
                + rTmp
                + " == 0 { "
                + rt
                + "ThrowNPE() }; return "
                + rTmp
                + " == "
                + aTmp
                + " }()";
          }
          String aGo = arg0Erased.getKind() == TypeKind.NULL ? "nil" : emitPath(arg0Path, ctx);
          if (arg0Erased.getKind() == TypeKind.NULL || isSimpleLocalOrParamIdent(arg0Path, aGo)) {
            return "(" + rt + "CheckNotNullPtr(" + rGo + ") == " + aGo + ")";
          }
          String rTmp = ctx.nextTemp("rp");
          String aTmp = ctx.nextTemp("ap");
          return "func() bool { "
              + rTmp
              + " := "
              + rGo
              + "; "
              + aTmp
              + " := "
              + aGo
              + "; return "
              + rt
              + "CheckNotNullPtr("
              + rTmp
              + ") == "
              + aTmp
              + " }()";
        }
        TypeMirror rawRecv =
            recvPath != null
                ? typeOf(recvPath)
                : (ctx.currentType != null ? ctx.currentType.asType() : null);
        TypeMirror rawArg0 = typeOf(arg0Path);
        if (resolver.isSameActiveComparableTypeVariable(rawRecv, rawArg0)) {
          String goTv = resolver.goType(rawRecv, ctx.imports);
          String rGo = recvPath != null ? emitPath(recvPath, ctx) : ctx.selfValueExpr();
          String aGo = emitPath(arg0Path, ctx);
          String rTmp = ctx.nextTemp("rk");
          String aTmp = ctx.nextTemp("ak");
          return "func() bool { "
              + rTmp
              + " := "
              + rGo
              + "; "
              + aTmp
              + " := "
              + aGo
              + "; if _, ok_ := any((*"
              + goTv
              + ")(nil)).(*string); ok_ { _ = "
              + rt
              + "CheckNotNullString(any("
              + rTmp
              + ").(string)); return "
              + rTmp
              + " == "
              + aTmp
              + " }; var z_ "
              + goTv
              + "; if "
              + rTmp
              + " == z_ { "
              + rt
              + "ThrowNPE() }; if "
              + rTmp
              + " == "
              + aTmp
              + " { return true }; return "
              + jl
              + "ObjectsEquals("
              + rt
              + "ToAny("
              + rTmp
              + "), "
              + rt
              + "ToAny("
              + aTmp
              + ")) }()";
        }
      }
      if (m.equals("hashCode") && arity == 0 && recvErased != null) {
        TypeMirror rawRecv =
            recvPath != null
                ? typeOf(recvPath)
                : (ctx.currentType != null ? ctx.currentType.asType() : null);
        if (resolver.isActiveTypeVariable(rawRecv)) {
          String goTv = resolver.goType(rawRecv, ctx.imports);
          String rGo = recvPath != null ? emitPath(recvPath, ctx) : ctx.selfValueExpr();
          String rTmp = ctx.nextTemp("hk");
          return "func() int32 { "
              + rTmp
              + " := "
              + rGo
              + "; if _, ok_ := any((*"
              + goTv
              + ")(nil)).(*string); ok_ { return "
              + jl
              + "StringHashCode("
              + rt
              + "CheckNotNullString(any("
              + rTmp
              + ").(string))) }; return "
              + jl
              + "ObjectHashCode("
              + rt
              + "CheckNotNull("
              + rt
              + "ToAny("
              + rTmp
              + "))) }()";
        }
        if (resolver.isNodeType(recvErased)) {
          String rNode = recvPath != null ? emitPath(recvPath, ctx) : ctx.selfValueExpr();
          return "int32(" + rt + "CheckNotNullNode(" + rNode + "))";
        }
        if (recvErased instanceof DeclaredType rdt && rdt.asElement() instanceof TypeElement rte) {
          if (resolver.isValueStructClass(rte)) {
            String rGo = recvPath != null ? emitPath(recvPath, ctx) : ctx.selfValueExpr();
            ExpressionTree recvTree =
                recvPath != null && recvPath.getLeaf() instanceof ExpressionTree et ? et : null;
            if (hasConcreteValueStructHashCodeMethod(rte)) {
              return "(" + emitCheckNotNull(recvTree, rGo, recvErased, ctx) + ").HashCode()";
            }
            if (valueStructHasIdentityEquals(rte)) {
              return "(" + emitCheckNotNull(recvTree, rGo, recvErased, ctx) + ").IdentityHash_()";
            }
          }
          if (resolver.isValueEnum(rte)) {
            String rGo = recvPath != null ? emitPath(recvPath, ctx) : ctx.selfValueExpr();
            ExpressionTree recvTree =
                recvPath != null && recvPath.getLeaf() instanceof ExpressionTree et ? et : null;
            return emitCheckNotNull(recvTree, rGo, recvErased, ctx) + ".HashCode()";
          }
          if (rte.getKind() == ElementKind.ENUM
              || (rte.getKind() == ElementKind.RECORD && resolver.isSourceType(rte))
              || hasConcreteGoHashCodeMethod(rte)) {
            String rGo = recvPath != null ? emitPath(recvPath, ctx) : ctx.selfValueExpr();
            return rt + "CheckNotNullPtr(" + rGo + ").HashCode()";
          }
          if (hasIdentityEquals(rte)) {
            String rGo = recvPath != null ? emitPath(recvPath, ctx) : ctx.selfValueExpr();
            return rt + "IdentityHash(" + rt + "CheckNotNullPtr(" + rGo + "))";
          }
        }
      }
      String recv =
          recvPath != null
              ? emitAs(
                  (ExpressionTree) recvPath.getLeaf(),
                  recvPath.getParentPath(),
                  resolver.getObjectElement().asType(),
                  true,
                  ctx)
              : emitImplicitSelfAsAny(ctx);
      if (m.equals("equals") && arity == 1) {
        String a0 = emitAs(args.get(0), path, resolver.getObjectElement().asType(), true, ctx);
        return jl + "ObjectsEquals(" + rt + "CheckNotNull(" + recv + "), " + a0 + ")";
      }
      if (m.equals("hashCode") && arity == 0) {
        return jl + "ObjectHashCode(" + rt + "CheckNotNull(" + recv + "))";
      }
      if (m.equals("toString") && arity == 0) {
        return jl + "ToJavaString(" + rt + "CheckNotNull(" + recv + "))";
      }
      if (m.equals("getClass") && arity == 0) {
        return rt + "GetObjectClass(" + recv + ")";
      }
    }

    // 2b. java.lang.Class intrinsics
    if (qn.equals("java.lang.Class")) {
      if (isStatic) {
        if (m.equals("forName") && arity >= 1) {
          return rt + "ClassForName(" + emit(args.get(0), path, ctx) + ")";
        }
      } else {
        String recv = recvPath != null ? emitPath(recvPath, ctx) : ctx.selfValueExpr();
        if (m.equals("getName") && arity == 0) return rt + "CheckNotNullPtr(" + recv + ").GetName()";
        if (m.equals("getSimpleName") && arity == 0) return rt + "CheckNotNullPtr(" + recv + ").GetSimpleName()";
        if (m.equals("getCanonicalName") && arity == 0) return rt + "CheckNotNullPtr(" + recv + ").GetName()";
        if (m.equals("desiredAssertionStatus") && arity == 0) {
          return "func() bool { _ = " + recv + "; return false }()";
        }
        if (m.equals("cast") && arity == 1) {
          return "func() "
              + resolver.goType(typeOf(path), ctx.imports)
              + " { _ = "
              + rt
              + "CheckNotNullPtr("
              + recv
              + "); return "
              + emitAs(args.get(0), path, typeOf(path), false, ctx)
              + " }()";
        }
        if (m.equals("isInstance") && arity == 1) {
          return "func() bool { _ = "
              + rt
              + "CheckNotNullPtr("
              + recv
              + "); return "
              + emitAs(args.get(0), path, resolver.getObjectElement().asType(), true, ctx)
              + " != nil }()";
        }
        if (m.equals("getResource") && arity == 1) {
          String retGo = resolver.goType(ee.getReturnType(), ctx.imports);
          TypeElement urlTe = resolver.elements.getTypeElement("java.net.URL");
          String structGo =
              urlTe != null ? resolver.runtime.goStructType(urlTe, ctx.imports) : "struct{}";
          return "func() "
              + retGo
              + " { if "
              + rt
              + "ClassResourceExists("
              + recv
              + ", "
              + emit(args.get(0), path, ctx)
              + ") { return &"
              + structGo
              + "{} }; return nil }()";
        }
        if (m.equals("getResourceAsStream") && arity == 1) {
          String retGo = resolver.goType(ee.getReturnType(), ctx.imports);
          TypeElement isTe = resolver.elements.getTypeElement("java.io.InputStream");
          String structGo =
              isTe != null ? resolver.runtime.goStructType(isTe, ctx.imports) : "struct{}";
          String selfAssign =
              (isTe != null && resolver.runtime.isSubclassedExternal(isTe)) ? "s_.Self_ = s_; " : "";
          return "func() "
              + retGo
              + " { if b_, ok_ := "
              + rt
              + "ClassResourceBytes("
              + recv
              + ", "
              + emit(args.get(0), path, ctx)
              + "); ok_ { s_ := &"
              + structGo
              + "{Core_: b_}; "
              + selfAssign
              + "return s_ }; return nil }()";
        }
      }
    }

    // Also intercept `.getClass()` on ANY type!
    if (m.equals("getClass") && arity == 0 && !isStatic) {
      String recv =
          recvPath != null
              ? emitAs(
                  (ExpressionTree) recvPath.getLeaf(),
                  recvPath.getParentPath(),
                  resolver.getObjectElement().asType(),
                  true,
                  ctx)
              : emitImplicitSelfAsAny(ctx);
      return rt + "GetObjectClass(" + recv + ")";
    }

    // 3. java.lang.String instance & static methods
    if (qn.equals("java.lang.String")) {
      return emitStringIntrinsic(mi, path, ee, isStatic, recvPath, ctx);
    }

    // 4. java.lang.CharSequence & java.lang.Comparable (mapped to `any`)
    if (qn.equals("java.lang.CharSequence") && !isStatic) {
      String recv = recvPath != null ? emitPath(recvPath, ctx) : ctx.selfValueExpr();
      if (m.equals("length") && arity == 0) return rt + "CharSeqLength(" + recv + ")";
      if (m.equals("charAt") && arity == 1) {
        return rt
            + "CharSeqCharAt("
            + recv
            + ", "
            + emitAs(args.get(0), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx)
            + ")";
      }
      if (m.equals("toString") && arity == 0) return rt + "StrVal(" + rt + "CheckNotNull(" + recv + "))";
      if (m.equals("isEmpty") && arity == 0) return "(" + rt + "CharSeqLength(" + recv + ") == 0)";
    }
    if (qn.equals("java.lang.Comparable") && !isStatic && m.equals("compareTo") && arity == 1) {
      TypeMirror recvErased =
          recvPath != null
              ? resolver.safeErasure(typeOf(recvPath))
              : (ctx.currentType != null ? resolver.safeErasure(ctx.currentType.asType()) : null);
      String recv = recvPath != null ? emitPath(recvPath, ctx) : ctx.selfValueExpr();
      if (recvErased != null && resolver.isStringType(recvErased)) {
        TypeMirror strTm = resolver.elements.getTypeElement("java.lang.String").asType();
        return rt + "JStrCompareTo(" + recv + ", " + emitAs(args.get(0), path, strTm, false, ctx) + ")";
      }
      return rt + "CompareValues(" + recv + ", " + emit(args.get(0), path, ctx) + ")";
    }

    // 5. java.util.Objects & com.google.common.base.Objects
    if (qn.equals("java.util.Objects") && isStatic) {
      if (m.equals("equals") && arity == 2) {
        return emitInlineObjectsEquals(args.get(0), args.get(1), path, ctx);
      }
      if (m.equals("hashCode") && arity == 1 && !ee.isVarArgs()) {
        return emitInlineSingleObjectHashCode(args.get(0), path, ctx);
      }
      if (m.equals("hash") && ee.isVarArgs()) {
        return emitInlineObjectsHash(args, path, ee, ctx);
      }
      if (m.equals("requireNonNull") && arity >= 1) {
        TypeMirror retTm = typeOf(path);
        String v = emitAs(args.get(0), path, retTm, false, ctx);
        return emitCheckNotNull(args.get(0), v, retTm, ctx);
      }
      if (m.equals("isNull") && arity == 1) {
        return emitInlineNullCheck(true, args.get(0), path, ctx);
      }
      if (m.equals("nonNull") && arity == 1) {
        return emitInlineNullCheck(false, args.get(0), path, ctx);
      }
      if (m.equals("toString") && arity == 1) {
        return rt + "StrVal(" + emitAs(args.get(0), path, resolver.getObjectElement().asType(), true, ctx) + ")";
      }
    }
    if (qn.equals("com.google.common.base.Objects") && isStatic) {
      if (m.equals("equal") && arity == 2) {
        return emitInlineObjectsEquals(args.get(0), args.get(1), path, ctx);
      }
      if (m.equals("hashCode") && ee.isVarArgs()) {
        return emitInlineObjectsHash(args, path, ee, ctx);
      }
    }

    // 6. com.google.common.base.Preconditions -> transpiled/shims/guava
    if (qn.equals("com.google.common.base.Preconditions") && isStatic) {
      return emitPreconditionsIntrinsic(mi, path, ee, ctx);
    }

    // 7. java.lang.Math & StrictMath
    if ((qn.equals("java.lang.Math") || qn.equals("java.lang.StrictMath")) && isStatic) {
      String mathCall = emitMathIntrinsic(mi, path, ee, ctx);
      if (mathCall != null) {
        return mathCall;
      }
    }

    // 8. Primitive wrapper helpers (Integer, Long, Double, Float, Character, Boolean)
    String primHelper = emitWrapperIntrinsic(mi, path, ee, qn, m, isStatic, recvPath, ctx);
    if (primHelper != null) {
      return primHelper;
    }

    // 9. Throwable methods (`getMessage`, `getCause`, `initCause`, `printStackTrace`, `addSuppressed`)
    if (resolver.runtime.isShimException(owner) && !isStatic) {
      String recv = recvPath != null ? emitPath(recvPath, ctx) : ctx.selfValueExpr();
      if (m.equals("getMessage") && arity == 0) return rt + "ThrowableGetMessage(" + recv + ")";
      if (m.equals("getCause") && arity == 0) return recv + ".GetCause()";
      if (m.equals("toString") && arity == 0) return recv + ".Error()";
      if (m.equals("addSuppressed") && arity == 1) {
        return "func() { _ = " + recv + "; _ = " + emit(args.get(0), path, ctx) + " }()";
      }
      if (m.equals("printStackTrace")) {
        if (arity == 0) {
          return rt + "ThrowablePrintStackTrace(" + recv + ", nil)";
        }
        if (arity == 1) {
          return rt + "ThrowablePrintStackTrace(" + recv + ", " + emit(args.get(0), path, ctx) + ")";
        }
      }
    }

    // 10. System intrinsics
    if (qn.equals("java.lang.System") && isStatic) {
      if (m.equals("identityHashCode") && arity == 1) {
        return rt
            + "IdentityHash("
            + emitAs(args.get(0), path, resolver.getObjectElement().asType(), true, ctx)
            + ")";
      }
      if (m.equals("currentTimeMillis") && arity == 0) return jl + "System_CurrentTimeMillis()";
      if (m.equals("nanoTime") && arity == 0) return jl + "System_NanoTime()";
      if (m.equals("lineSeparator") && arity == 0) return "\"\\n\"";
      if (m.equals("arraycopy") && arity == 5) {
        TreePath sPath = new TreePath(path, args.get(0));
        TreePath dPath = new TreePath(path, args.get(2));
        TypeMirror sEff = resolver.goEffectiveType(typeOf(sPath));
        TypeMirror dEff = resolver.goEffectiveType(typeOf(dPath));
        if (sEff instanceof ArrayType && dEff instanceof ArrayType) {
          String sGo = resolver.goType(sEff, ctx.imports);
          String dGo = resolver.goType(dEff, ctx.imports);
          if (sGo.equals(dGo)) {
            String sTmp = ctx.nextTemp("acSrc");
            String spTmp = ctx.nextTemp("acSP");
            String dTmp = ctx.nextTemp("acDst");
            String dpTmp = ctx.nextTemp("acDP");
            String lnTmp = ctx.nextTemp("acLen");
            String sExpr = emitPath(sPath, ctx);
            String spExpr =
                emitAs(
                    args.get(1),
                    path,
                    resolver.types.getPrimitiveType(TypeKind.INT),
                    false,
                    ctx);
            String dExpr = emitPath(dPath, ctx);
            String dpExpr =
                emitAs(
                    args.get(3),
                    path,
                    resolver.types.getPrimitiveType(TypeKind.INT),
                    false,
                    ctx);
            String lnExpr =
                emitAs(
                    args.get(4),
                    path,
                    resolver.types.getPrimitiveType(TypeKind.INT),
                    false,
                    ctx);
            return "func() { "
                + sTmp
                + " := "
                + sExpr
                + "; "
                + spTmp
                + " := "
                + spExpr
                + "; "
                + dTmp
                + " := "
                + dExpr
                + "; "
                + dpTmp
                + " := "
                + dpExpr
                + "; "
                + lnTmp
                + " := "
                + lnExpr
                + "; "
                + rt
                + "CheckNotNullSlice("
                + sTmp
                + "); "
                + rt
                + "CheckNotNullSlice("
                + dTmp
                + "); if "
                + spTmp
                + " < 0 || "
                + dpTmp
                + " < 0 || "
                + lnTmp
                + " < 0 || int64("
                + spTmp
                + ")+int64("
                + lnTmp
                + ") > int64(len("
                + sTmp
                + ")) || int64("
                + dpTmp
                + ")+int64("
                + lnTmp
                + ") > int64(len("
                + dTmp
                + ")) { panic("
                + rt
                + "ConstructShimThrowable(\"java.lang.ArrayIndexOutOfBoundsException\")) }; copy("
                + dTmp
                + "[int("
                + dpTmp
                + "):int("
                + dpTmp
                + ")+int("
                + lnTmp
                + ")], "
                + sTmp
                + "[int("
                + spTmp
                + "):int("
                + spTmp
                + ")+int("
                + lnTmp
                + ")]) }()";
          }
        }
        return rt
            + "Arraycopy("
            + emitAs(args.get(0), path, resolver.getObjectElement().asType(), true, ctx)
            + ", "
            + emitAs(args.get(1), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx)
            + ", "
            + emitAs(args.get(2), path, resolver.getObjectElement().asType(), true, ctx)
            + ", "
            + emitAs(args.get(3), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx)
            + ", "
            + emitAs(args.get(4), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx)
            + ")";
      }
    }

    // 11. Array `.clone()` (`byte[].clone()`, `Prop[].clone()`, etc.)
    if (m.equals("clone")
        && arity == 0
        && recvPath != null
        && typeOf(recvPath).getKind() == TypeKind.ARRAY) {
      return rt + "CloneSlice(" + emitPath(recvPath, ctx) + ")";
    }

    // 12. JSCompObjects.identical(a, b) -> zero-alloc direct Go `==` when types match!
    if (qn.equals("com.google.javascript.jscomp.base.JSCompObjects")
        && isStatic
        && m.equals("identical")
        && arity == 2) {
      TreePath p0 = new TreePath(path, args.get(0));
      TreePath p1 = new TreePath(path, args.get(1));
      return emitEquality(true, p0, typeOf(p0), p1, typeOf(p1), ctx);
    }

    // 13. Unboxed Predicate<Node> apply/test -> j2gort.PredTestNode
    if (!isStatic
        && arity == 1
        && ee.getParameters().size() == 1
        && ee.getReturnType().getKind() == TypeKind.BOOLEAN
        && !ee.getModifiers().contains(Modifier.PRIVATE)
        && ((qn.equals("com.google.common.base.Predicate")
                && (m.equals("apply") || m.equals("test")))
            || (qn.equals("java.util.function.Predicate") && m.equals("test")))) {
      TreePath arg0Path = new TreePath(path, args.get(0));
      if (resolver.isNodeType(typeOf(arg0Path))) {
        String recvGo =
            recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
        String arg0Go = emitPath(arg0Path, ctx);
        return rt + "PredTestNode(" + recvGo + ", uint32(" + arg0Go + "))";
      }
    }

    // 13b. Static Guava ImmutableList / ImmutableSet zero-box & zero-escape intrinsics
    if (isStatic) {
      if (qn.equals("com.google.common.collect.ImmutableList")) {
        if (m.equals("builder")
            && arity == 0
            && !isDirectExprStmt(path)
            && resolver.safeErasure(typeOf(path)) instanceof DeclaredType bDt
            && bDt.asElement() instanceof TypeElement bTe
            && bTe.getQualifiedName().contentEquals("com.google.common.collect.ImmutableList.Builder")) {
          return "(&" + rt + "Ext_com_google_common_collect_ImmutableList_Builder{})";
        }
        if (m.equals("builderWithExpectedSize")
            && arity == 1
            && ee.getParameters().size() == 1
            && ee.getParameters().get(0).asType().getKind() == TypeKind.INT
            && resolver.safeErasure(typeOf(path)) instanceof DeclaredType bDt
            && bDt.asElement() instanceof TypeElement bTe
            && bTe.getQualifiedName().contentEquals("com.google.common.collect.ImmutableList.Builder")) {
          String capGo =
              emitAs(args.get(0), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
          return rt
              + "InitImmListBuilderCap(&"
              + rt
              + "Ext_com_google_common_collect_ImmutableList_Builder{}, "
              + capGo
              + ")";
        }
        if (m.equals("of")
            && !ee.isVarArgs()
            && arity >= 1
            && arity <= 4
            && resolver.safeErasure(typeOf(path)) instanceof DeclaredType ilDt
            && ilDt.asElement() instanceof TypeElement ilTe
            && ilTe.getQualifiedName().contentEquals("com.google.common.collect.ImmutableList")) {
          boolean allStr = true;
          for (int i = 0; i < arity; i++) {
            if (!resolver.isStringType(typeOf(new TreePath(path, args.get(i))))) {
              allStr = false;
              break;
            }
          }
          if (allStr) {
            StringBuilder sb = new StringBuilder(rt).append("ImmListOfStr").append(arity).append('(');
            for (int i = 0; i < arity; i++) {
              if (i > 0) sb.append(", ");
              sb.append(emit(args.get(i), path, ctx));
            }
            return sb.append(')').toString();
          }
          if (arity == 1) {
            TreePath arg0Path = new TreePath(path, args.get(0));
            TypeMirror arg0Tm = typeOf(arg0Path);
            if (arg0Tm != null && arg0Tm.getKind() != TypeKind.NULL) {
              String arg0GoType =
                  resolver.goType(arg0Tm, new PackagePlanner.ImportTracker(""));
              if (arg0GoType.startsWith("*")) {
                return rt + "ImmListOfPtr(" + emitPath(arg0Path, ctx) + ")";
              }
            }
          }
        }
        if (m.equals("copyOf")
            && !ee.isVarArgs()
            && arity == 1
            && ee.getParameters().size() == 1
            && ee.getParameters().get(0).asType() instanceof ArrayType
            && resolver.safeErasure(typeOf(path)) instanceof DeclaredType ilCopyDt
            && ilCopyDt.asElement() instanceof TypeElement ilCopyTe
            && ilCopyTe.getQualifiedName().contentEquals("com.google.common.collect.ImmutableList")) {
          TreePath arg0Path = new TreePath(path, args.get(0));
          TypeMirror arg0Eff = resolver.goEffectiveType(typeOf(arg0Path));
          if (arg0Eff instanceof ArrayType aat
              && resolver.isNullableValueStructType(aat.getComponentType())) {
            return rt + "ImmListFromCompSlice(" + emitPath(arg0Path, ctx) + ")";
          }
        }
      } else if (qn.equals("com.google.common.collect.ImmutableSet")
          && m.equals("of")
          && !ee.isVarArgs()
          && (arity == 1 || arity == 2)
          && resolver.safeErasure(typeOf(path)) instanceof DeclaredType isDt
          && isDt.asElement() instanceof TypeElement isTe
          && isTe.getQualifiedName().contentEquals("com.google.common.collect.ImmutableSet")) {
        boolean allStr = true;
        for (int i = 0; i < arity; i++) {
          if (!resolver.isStringType(typeOf(new TreePath(path, args.get(i))))) {
            allStr = false;
            break;
          }
        }
        if (allStr) {
          StringBuilder sb = new StringBuilder(rt).append("ImmSetOfStr").append(arity).append('(');
          for (int i = 0; i < arity; i++) {
            if (i > 0) sb.append(", ");
            sb.append(emit(args.get(i), path, ctx));
          }
          return sb.append(')').toString();
        }
      } else if (qn.equals("java.util.Collections")
          && m.equals("singletonList")
          && arity == 1
          && resolver.safeErasure(typeOf(path)) instanceof DeclaredType slDt
          && slDt.asElement() instanceof TypeElement slTe
          && slTe.getQualifiedName().contentEquals("java.util.List")) {
        TreePath arg0Path = new TreePath(path, args.get(0));
        TypeMirror arg0Tm = typeOf(arg0Path);
        if (arg0Tm != null && arg0Tm.getKind() != TypeKind.NULL) {
          String arg0GoType = resolver.goType(arg0Tm, ctx.imports);
          if (arg0GoType.startsWith("*")) {
            return rt + "SingletonListPtr(" + emitPath(arg0Path, ctx) + ")";
          }
        }
      }
    }

    // 13c. Instance intrinsics for Interner<String>, ImmutableList.Builder, ImmutableSet.Builder,
    //      Immutable*Multimap.Builder, and List<String>.get(int)
    if (!isStatic && !ee.getModifiers().contains(Modifier.PRIVATE)) {
      if (m.equals("intern")
          && arity == 1
          && ee.getParameters().size() == 1
          && isSubtypeOfQualified(owner, "com.google.common.collect.Interner")) {
        TreePath arg0Path = new TreePath(path, args.get(0));
        if (resolver.isStringType(typeOf(arg0Path))
            && (resolver.isStringType(typeOf(path)) || isDirectExprStmt(path))) {
          TypeElement internerTe =
              resolver.elements.getTypeElement("com.google.common.collect.Interner");
          String recvGo =
              recvPath != null
                  ? (internerTe != null
                      ? convert(
                          emitPath(recvPath, ctx), typeOf(recvPath), internerTe.asType(), false, ctx)
                      : emitPath(recvPath, ctx))
                  : resolveImplicitThisValue(owner, ctx);
          return recvGo + ".InternStr(" + emitPath(arg0Path, ctx) + ")";
        }
      }

      if (qn.equals("com.google.common.collect.ImmutableList.Builder")) {
        boolean isStmt = isDirectExprStmt(path);
        boolean retIsBuilder =
            !isStmt
                && resolver.safeErasure(typeOf(path)) instanceof DeclaredType rDt
                && rDt.asElement() instanceof TypeElement rTe
                && rTe.getQualifiedName()
                    .contentEquals("com.google.common.collect.ImmutableList.Builder");
        if (m.equals("add")
            && arity == 1
            && !ee.isVarArgs()
            && ee.getParameters().size() == 1
            && ee.getParameters().get(0).asType().getKind() != TypeKind.ARRAY
            && (isStmt || retIsBuilder)) {
          String recvGo =
              recvPath != null
                  ? convert(emitPath(recvPath, ctx), typeOf(recvPath), owner.asType(), false, ctx)
                  : resolveImplicitThisValue(owner, ctx);
          TreePath arg0Path = new TreePath(path, args.get(0));
          if (resolver.isStringType(typeOf(arg0Path))) {
            String sGo = emitPath(arg0Path, ctx);
            return recvGo + (isStmt ? ".AddStrVoid_(" : ".AddStrSelf_(") + sGo + ")";
          }
          String vGo = emitAs(args.get(0), path, resolver.getObjectElement().asType(), true, ctx);
          return recvGo + (isStmt ? ".AddVoid_(" : ".AddSelf_(") + vGo + ")";
        }
        if (m.equals("addAll")
            && arity == 1
            && ee.getParameters().size() == 1
            && (isStmt || retIsBuilder)
            && resolver.safeErasure(ee.getParameters().get(0).asType()) instanceof DeclaredType pDt
            && pDt.asElement() instanceof TypeElement pTe
            && pTe.getQualifiedName().contentEquals("java.lang.Iterable")) {
          String recvGo =
              recvPath != null
                  ? convert(emitPath(recvPath, ctx), typeOf(recvPath), owner.asType(), false, ctx)
                  : resolveImplicitThisValue(owner, ctx);
          String itGo = emitAs(args.get(0), path, ee.getParameters().get(0).asType(), false, ctx);
          return recvGo + (isStmt ? ".AddAllVoid_(" : ".AddAllSelf_(") + itGo + ")";
        }
        if (m.equals("build")
            && arity == 0
            && ee.getParameters().isEmpty()
            && (isStmt
                || (resolver.safeErasure(typeOf(path)) instanceof DeclaredType ilDt
                    && ilDt.asElement() instanceof TypeElement ilTe
                    && ilTe.getQualifiedName()
                        .contentEquals("com.google.common.collect.ImmutableList")))) {
          String recvGo =
              recvPath != null
                  ? convert(emitPath(recvPath, ctx), typeOf(recvPath), owner.asType(), false, ctx)
                  : resolveImplicitThisValue(owner, ctx);
          return recvGo + ".BuildImmList_()";
        }
      }

      if (qn.equals("com.google.common.collect.ImmutableSet.Builder")
          && m.equals("add")
          && arity == 1
          && !ee.isVarArgs()
          && ee.getParameters().size() == 1
          && ee.getParameters().get(0).asType().getKind() != TypeKind.ARRAY) {
        TreePath arg0Path = new TreePath(path, args.get(0));
        if (resolver.isStringType(typeOf(arg0Path))) {
          boolean isStmt = isDirectExprStmt(path);
          boolean retIsSetBuilder =
              !isStmt
                  && resolver.safeErasure(typeOf(path)) instanceof DeclaredType rDt
                  && rDt.asElement() instanceof TypeElement rTe
                  && rTe.getQualifiedName()
                      .contentEquals("com.google.common.collect.ImmutableSet.Builder");
          if (isStmt || retIsSetBuilder) {
            String recvGo =
                recvPath != null
                    ? convert(emitPath(recvPath, ctx), typeOf(recvPath), owner.asType(), false, ctx)
                    : resolveImplicitThisValue(owner, ctx);
            String sGo = emitPath(arg0Path, ctx);
            return rt
                + (isStmt ? "ImmSetBuilderAddStrVoid(" : "ImmSetBuilderAddStr(")
                + recvGo
                + ", "
                + sGo
                + ")";
          }
        }
      }

      if ((qn.equals("com.google.common.collect.ImmutableMultimap.Builder")
              || qn.equals("com.google.common.collect.ImmutableListMultimap.Builder")
              || qn.equals("com.google.common.collect.ImmutableSetMultimap.Builder"))
          && m.equals("put")
          && arity == 2
          && ee.getParameters().size() == 2) {
        TreePath arg0Path = new TreePath(path, args.get(0));
        if (resolver.isStringType(typeOf(arg0Path))) {
          String recvGo =
              recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
          String kGo = emitPath(arg0Path, ctx);
          TypeMirror objTm = resolver.getObjectElement().asType();
          String vGo = emitAs(args.get(1), path, objTm, true, ctx);
          String raw = rt + "MultimapBuilderPutStr(" + recvGo + ", " + kGo + ", " + vGo + ")";
          TypeMirror callRet = typeOf(path);
          return (!isDirectExprStmt(path) && callRet.getKind() != TypeKind.VOID)
              ? convert(raw, objTm, callRet, false, ctx)
              : raw;
        }
      }

      if (m.equals("get")
          && arity == 1
          && ee.getParameters().size() == 1
          && ee.getParameters().get(0).asType().getKind() == TypeKind.INT
          && isSubtypeOfQualified(owner, "java.util.List")) {
        TypeMirror callRet = typeOf(path);
        if (resolver.isStringType(callRet)) {
          String recvGo =
              recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
          String idxGo =
              emitAs(args.get(0), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
          return rt + "ListGetStr(" + recvGo + ", " + idxGo + ")";
        }
        if (resolver.isByteStringType(callRet)
            && !isDirectExprStmt(path)
            && (owner.getQualifiedName().contentEquals("java.util.List")
                || isNonLoweredCollectionExpr(recvPath))) {
          String recvGo =
              recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
          String idxGo =
              emitAs(args.get(0), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
          return rt + "ListGetByteString(" + recvGo + ", " + idxGo + ")";
        }
        if (!isDirectExprStmt(path)
            && (resolver.isValueStructType(callRet) || resolver.isValueEnumType(callRet))
            && !resolver.isActiveTypeVariable(callRet)
            && isNonLoweredCollectionExpr(recvPath)) {
          String elemGo = resolver.goType(resolver.safeErasure(callRet), ctx.imports);
          String recvGo =
              recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
          String idxGo =
              emitAs(args.get(0), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
          return rt + "ListGetVal[" + elemGo + "](" + recvGo + ", " + idxGo + ")";
        }
      }
    }

    // 14. Zero-boxing String / Node fast paths for Map, Set, Multimap, Multiset, ImmutableCollection, and Deque<Node>
    if (!isStatic
        && arity <= 3
        && !ee.getModifiers().contains(Modifier.PRIVATE)) {
      String fastColl =
          tryEmitFastKeyedCollectionCall(mi, path, ee, owner, m, arity, recvPath, rt, ctx);
      if (fastColl != null) {
        return fastColl;
      }
    }

    return null;
  }

  public boolean isNonLoweredCollectionExpr(TreePath path) {
    if (path == null) {
      return false;
    }
    TreePath stripped = unwrapParens(path);
    Tree leaf = stripped.getLeaf();
    if (leaf instanceof IdentifierTree) {
      Element el = resolver.trees.getElement(stripped);
      return el instanceof VariableElement ve && ve.getKind() == ElementKind.PARAMETER;
    }
    if (leaf instanceof MethodInvocationTree) {
      Element el = resolver.trees.getElement(stripped);
      if (el instanceof ExecutableElement mee) {
        return mee.getModifiers().contains(Modifier.STATIC) || !mee.getParameters().isEmpty();
      }
    }
    return false;
  }

  private static boolean isDirectExprStmt(TreePath path) {
    TreePath cur = path.getParentPath();
    while (cur != null && cur.getLeaf() instanceof ParenthesizedTree) {
      cur = cur.getParentPath();
    }
    return cur != null && cur.getLeaf() instanceof ExpressionStatementTree;
  }

  private String tryEmitFastKeyedCollectionCall(
      MethodInvocationTree mi,
      TreePath path,
      ExecutableElement ee,
      TypeElement owner,
      String m,
      int arity,
      TreePath recvPath,
      String rt,
      EmitContext ctx) {
    List<? extends ExpressionTree> args = mi.getArguments();
    boolean isStmt = isDirectExprStmt(path);
    TypeMirror callRet = typeOf(path);

    if (isSubtypeOfQualified(owner, "java.util.Deque")) {
      if (arity == 0 && ee.getParameters().isEmpty() && resolver.isNodeType(callRet)) {
        String dequeFn =
            switch (m) {
              case "pop", "removeFirst", "remove" -> "DequePopNode(";
              case "removeLast" -> "DequeRemoveLastNode(";
              case "poll", "pollFirst" -> "DequePollFirstNode(";
              case "peek", "peekFirst" -> "DequePeekFirstNode(";
              case "peekLast" -> "DequePeekLastNode(";
              default -> null;
            };
        if (dequeFn != null) {
          String recvExpr =
              recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
          String rawCall = rt + dequeFn + recvExpr + ")";
          return isStmt ? rawCall : (resolver.goType(callRet, ctx.imports) + "(" + rawCall + ")");
        }
      } else if (arity == 0
          && ee.getParameters().isEmpty()
          && isStmt
          && isColorIdType(resolver.safeErasure(callRet))) {
        String dequeVoidFn =
            switch (m) {
              case "pop", "removeFirst", "remove" -> "DequePopVoid(";
              case "removeLast" -> "DequeRemoveLastVoid(";
              default -> null;
            };
        if (dequeVoidFn != null) {
          String recvExpr =
              recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
          return rt + dequeVoidFn + recvExpr + ")";
        }
      } else if (arity == 1 && ee.getParameters().size() == 1) {
        TreePath arg0Path = new TreePath(path, args.get(0));
        TypeMirror arg0Er = resolver.safeErasure(typeOf(arg0Path));
        if (resolver.isNodeType(arg0Er)) {
          if (m.equals("push") && ee.getReturnType().getKind() == TypeKind.VOID) {
            String recvExpr =
                recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
            String nGo = emitPath(arg0Path, ctx);
            return rt + "DequePushNode(" + recvExpr + ", uint32(" + nGo + "))";
          }
          if ((m.equals("addLast") && ee.getReturnType().getKind() == TypeKind.VOID)
              || ((m.equals("add") || m.equals("offerLast"))
                  && ee.getReturnType().getKind() == TypeKind.BOOLEAN)) {
            String recvExpr =
                recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
            String nGo = emitPath(arg0Path, ctx);
            return rt + "DequeAddLastNode(" + recvExpr + ", uint32(" + nGo + "))";
          }
        } else if (isColorIdType(arg0Er)) {
          if ((m.equals("push") || m.equals("addFirst"))
              && ee.getReturnType().getKind() == TypeKind.VOID) {
            String recvExpr =
                recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
            String cGo = emitPath(arg0Path, ctx);
            return rt + "DequePushColorId(" + recvExpr + ", " + cGo + ")";
          }
          if (m.equals("addLast") && ee.getReturnType().getKind() == TypeKind.VOID) {
            String recvExpr =
                recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
            String cGo = emitPath(arg0Path, ctx);
            return rt + "DequeAddLastColorIdVoid(" + recvExpr + ", " + cGo + ")";
          }
          if ((m.equals("add") || m.equals("offerLast"))
              && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
            String recvExpr =
                recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
            String cGo = emitPath(arg0Path, ctx);
            return rt + "DequeAddLastColorIdBool(" + recvExpr + ", " + cGo + ")";
          }
        }
      }
    }

    if (arity == 0) {
      return null;
    }

    switch (m) {
      case "get",
          "containsKey",
          "containsEntry",
          "containsRow",
          "row",
          "put",
          "putAll",
          "remove",
          "removeAll",
          "getOrDefault",
          "computeIfAbsent",
          "contains",
          "add",
          "count",
          "setCount" -> {}
      default -> {
        return null;
      }
    }
    TreePath arg0Path = new TreePath(path, args.get(0));
    TypeMirror arg0Tm = resolver.safeErasure(typeOf(arg0Path));
    boolean strKey = resolver.isStringType(arg0Tm);
    boolean nodeKey = !strKey && resolver.isNodeType(arg0Tm);
    TypeMirror objTm = resolver.getObjectElement().asType();

    if (isColorIdType(arg0Tm)
        && isSubtypeOfQualified(owner, "com.google.common.collect.Table")) {
      String recvExpr =
          recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
      String rGo = emitPath(arg0Path, ctx);
      if (m.equals("containsRow")
          && arity == 1
          && ee.getParameters().size() == 1
          && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
        return rt + "TableContainsRowColorId(" + recvExpr + ", " + rGo + ")";
      }
      if (m.equals("contains")
          && arity == 2
          && ee.getParameters().size() == 2
          && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
        String cGo = emitAs(args.get(1), path, objTm, true, ctx);
        return rt + "TableContainsColorId(" + recvExpr + ", " + rGo + ", " + cGo + ")";
      }
      if (m.equals("get") && arity == 2 && ee.getParameters().size() == 2) {
        String cGo = emitAs(args.get(1), path, objTm, true, ctx);
        String raw = rt + "TableGetColorId(" + recvExpr + ", " + rGo + ", " + cGo + ")";
        return (!isStmt && callRet.getKind() != TypeKind.VOID)
            ? convert(raw, objTm, callRet, false, ctx)
            : raw;
      }
      if (m.equals("put") && arity == 3 && ee.getParameters().size() == 3) {
        String cGo = emitAs(args.get(1), path, objTm, true, ctx);
        String vGo = emitAs(args.get(2), path, objTm, true, ctx);
        String raw =
            rt + "TablePutColorId(" + recvExpr + ", " + rGo + ", " + cGo + ", " + vGo + ")";
        return (!isStmt && callRet.getKind() != TypeKind.VOID)
            ? convert(raw, objTm, callRet, false, ctx)
            : raw;
      }
      if (m.equals("remove") && arity == 2 && ee.getParameters().size() == 2) {
        String cGo = emitAs(args.get(1), path, objTm, true, ctx);
        String raw = rt + "TableRemoveColorId(" + recvExpr + ", " + rGo + ", " + cGo + ")";
        return (!isStmt && callRet.getKind() != TypeKind.VOID)
            ? convert(raw, objTm, callRet, false, ctx)
            : raw;
      }
      if (m.equals("row") && arity == 1 && ee.getParameters().size() == 1) {
        TypeElement mapTe = resolver.elements.getTypeElement("java.util.Map");
        if (mapTe != null) {
          String raw = rt + "TableRowColorId(" + recvExpr + ", " + rGo + ")";
          return (!isStmt && callRet.getKind() != TypeKind.VOID)
              ? convert(raw, mapTe.asType(), callRet, false, ctx)
              : raw;
        }
      }
    }

    if (!strKey
        && !nodeKey
        && arg0Tm instanceof DeclaredType eqKdt
        && eqKdt.asElement() instanceof TypeElement eqKte
        && resolver.isValueStructClass(eqKte)
        && valueStructGoEqMatchesJavaEquals(eqKte)
        && !resolver.hasEffectiveTypeParameters(eqKte)) {
      if (isSubtypeOfQualified(owner, "com.google.common.collect.ImmutableMap")) {
        String recvExpr =
            recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
        String kGo = emitPath(arg0Path, ctx);
        if (m.equals("get") && arity == 1 && ee.getParameters().size() == 1) {
          String raw = rt + "MapGetEqVal(" + recvExpr + ", " + kGo + ")";
          return (!isStmt && callRet.getKind() != TypeKind.VOID)
              ? convert(raw, objTm, callRet, false, ctx)
              : raw;
        }
        if (m.equals("containsKey")
            && arity == 1
            && ee.getParameters().size() == 1
            && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
          return rt + "MapContainsEqVal(" + recvExpr + ", " + kGo + ")";
        }
        if (m.equals("getOrDefault") && arity == 2 && ee.getParameters().size() == 2) {
          String defGo = emitAs(args.get(1), path, objTm, true, ctx);
          String raw = rt + "MapGetOrDefaultEqVal(" + recvExpr + ", " + kGo + ", " + defGo + ")";
          return (!isStmt && callRet.getKind() != TypeKind.VOID)
              ? convert(raw, objTm, callRet, false, ctx)
              : raw;
        }
      } else if (m.equals("contains")
          && arity == 1
          && ee.getParameters().size() == 1
          && ee.getReturnType().getKind() == TypeKind.BOOLEAN
          && isSubtypeOfQualified(owner, "com.google.common.collect.ImmutableCollection")) {
        String recvExpr =
            recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
        String kGo = emitPath(arg0Path, ctx);
        return rt + "SetContainsEqVal(" + recvExpr + ", " + kGo + ")";
      }
    }

    if (m.equals("computeIfAbsent")
        && arity == 2
        && ee.getParameters().size() == 2
        && isSubtypeOfQualified(owner, "java.util.Map")
        && resolver.safeErasure(ee.getParameters().get(1).asType()) instanceof DeclaredType fnDt
        && fnDt.asElement() instanceof TypeElement fnTe
        && fnTe.getQualifiedName().contentEquals("java.util.function.Function")) {
      boolean isInlineLambda =
          unwrapParens(new TreePath(path, args.get(1))).getLeaf() instanceof LambdaExpressionTree;
      if (strKey || nodeKey) {
        String recvExpr =
            recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
        String kGo = emitPath(arg0Path, ctx);
        String computeFn = strKey ? "MapComputeIfAbsentStr(" : "MapComputeIfAbsentNode(";
        if (!isInlineLambda) {
          String fnGo = emitAs(args.get(1), path, fnDt, false, ctx);
          String raw = rt + computeFn + recvExpr + ", " + kGo + ", " + fnGo + ")";
          return (!isStmt && callRet.getKind() != TypeKind.VOID)
              ? convert(raw, objTm, callRet, false, ctx)
              : raw;
        }
        String tryFn = strKey ? "MapTryGetComputeStr(" : "MapTryGetComputeNode(";
        String mTmp = ctx.nextTemp("m");
        String kTmp = ctx.nextTemp("k");
        String vTmp = ctx.nextTemp("v");
        String okTmp = ctx.nextTemp("ok");
        String fnGo = emitAs(args.get(1), path, fnDt, false, ctx);
        String slowCall = rt + computeFn + mTmp + ", " + kTmp + ", " + fnGo + ")";
        if (isStmt || callRet.getKind() == TypeKind.VOID) {
          return "func() { "
              + mTmp
              + " := "
              + recvExpr
              + "; "
              + kTmp
              + " := "
              + kGo
              + "; if _, "
              + okTmp
              + " := "
              + rt
              + tryFn
              + mTmp
              + ", "
              + kTmp
              + "); !"
              + okTmp
              + " { "
              + slowCall
              + " } }()";
        }
        String retGo = resolver.goType(callRet, ctx.imports);
        return "func() "
            + retGo
            + " { "
            + mTmp
            + " := "
            + recvExpr
            + "; "
            + kTmp
            + " := "
            + kGo
            + "; if "
            + vTmp
            + ", "
            + okTmp
            + " := "
            + rt
            + tryFn
            + mTmp
            + ", "
            + kTmp
            + "); "
            + okTmp
            + " { return "
            + convert(vTmp, objTm, callRet, false, ctx)
            + " }; return "
            + convert(slowCall, objTm, callRet, false, ctx)
            + " }()";
      }
      if (isInlineLambda) {
        TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
        if (sig != null
            && sig.canonParamTypes.size() == 2
            && resolver.goType(sig.canonParamTypes.get(0), ctx.imports).equals("any")) {
          String recvExpr =
              recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
          String kGo = emitAs(args.get(0), path, objTm, true, ctx);
          String mTmp = ctx.nextTemp("m");
          String kTmp = ctx.nextTemp("k");
          String vTmp = ctx.nextTemp("v");
          String okTmp = ctx.nextTemp("ok");
          String fnGo = emitAs(args.get(1), path, fnDt, false, ctx);
          String slowCall = mTmp + "." + sig.goName + "(" + kTmp + ", " + fnGo + ")";
          if (isStmt || callRet.getKind() == TypeKind.VOID) {
            return "func() { "
                + mTmp
                + " := "
                + recvExpr
                + "; "
                + kTmp
                + " := "
                + kGo
                + "; if _, "
                + okTmp
                + " := "
                + rt
                + "MapTryGetComputeAny("
                + mTmp
                + ", "
                + kTmp
                + "); !"
                + okTmp
                + " { "
                + slowCall
                + " } }()";
          }
          String retGo = resolver.goType(callRet, ctx.imports);
          TypeMirror slowRetTm = sig.canonReturnType != null ? sig.canonReturnType : objTm;
          return "func() "
              + retGo
              + " { "
              + mTmp
              + " := "
              + recvExpr
              + "; "
              + kTmp
              + " := "
              + kGo
              + "; if "
              + vTmp
              + ", "
              + okTmp
              + " := "
              + rt
              + "MapTryGetComputeAny("
              + mTmp
              + ", "
              + kTmp
              + "); "
              + okTmp
              + " { return "
              + convert(vTmp, objTm, callRet, false, ctx)
              + " }; return "
              + convert(slowCall, slowRetTm, callRet, false, ctx)
              + " }()";
        }
      }
      return null;
    }

    if (!strKey && !nodeKey) {
      return null;
    }

    if (isSubtypeOfQualified(owner, "java.util.Map")) {
      String recvExpr = recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
      String kGo = emitPath(arg0Path, ctx);
      if (m.equals("get") && arity == 1) {
        String raw = rt + (strKey ? "MapGetStr(" : "MapGetNode(") + recvExpr + ", " + kGo + ")";
        return (!isStmt && callRet.getKind() != TypeKind.VOID)
            ? convert(raw, objTm, callRet, false, ctx)
            : raw;
      }
      if (m.equals("containsKey") && arity == 1 && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
        return rt + (strKey ? "MapContainsKeyStr(" : "MapContainsKeyNode(") + recvExpr + ", " + kGo + ")";
      }
      if (m.equals("put") && arity == 2 && ee.getParameters().size() == 2) {
        String vGo = emitAs(args.get(1), path, objTm, true, ctx);
        String raw =
            rt + (strKey ? "MapPutStr(" : "MapPutNode(") + recvExpr + ", " + kGo + ", " + vGo + ")";
        return (!isStmt && callRet.getKind() != TypeKind.VOID)
            ? convert(raw, objTm, callRet, false, ctx)
            : raw;
      }
      if (m.equals("remove") && arity == 1 && ee.getParameters().size() == 1) {
        String raw = rt + (strKey ? "MapRemoveStr(" : "MapRemoveNode(") + recvExpr + ", " + kGo + ")";
        return (!isStmt && callRet.getKind() != TypeKind.VOID)
            ? convert(raw, objTm, callRet, false, ctx)
            : raw;
      }
      if (m.equals("getOrDefault") && arity == 2 && ee.getParameters().size() == 2) {
        String defGo = emitAs(args.get(1), path, objTm, true, ctx);
        String raw =
            rt
                + (strKey ? "MapGetOrDefaultStr(" : "MapGetOrDefaultNode(")
                + recvExpr
                + ", "
                + kGo
                + ", "
                + defGo
                + ")";
        return (!isStmt && callRet.getKind() != TypeKind.VOID)
            ? convert(raw, objTm, callRet, false, ctx)
            : raw;
      }
      return null;
    }

    if (arity == 1
        && ee.getParameters().size() == 1
        && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
      if (isSubtypeOfQualified(owner, "java.util.Set")) {
        if (m.equals("contains") || m.equals("add") || m.equals("remove")) {
          String recvExpr =
              recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
          String kGo = emitPath(arg0Path, ctx);
          String fn =
              switch (m) {
                case "contains" -> strKey ? "SetContainsStr(" : "SetContainsNode(";
                case "add" -> strKey ? "SetAddStr(" : "SetAddNode(";
                default -> strKey ? "SetRemoveStr(" : "SetRemoveNode(";
              };
          return rt + fn + recvExpr + ", " + kGo + ")";
        }
      } else if (m.equals("contains")
          && isSubtypeOfQualified(owner, "com.google.common.collect.ImmutableCollection")) {
        String recvExpr =
            recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
        String kGo = emitPath(arg0Path, ctx);
        return rt + (strKey ? "SetContainsStr(" : "SetContainsNode(") + recvExpr + ", " + kGo + ")";
      }
    }

    if (owner.getQualifiedName().contentEquals("com.google.javascript.jscomp.SimpleMultimap")
        && resolver.isSourceType(owner)
        && !resolver.isSubclassed(owner)) {
      String recvExpr;
      if (recvPath != null) {
        TypeMirror rTm = typeOf(recvPath);
        TypeMirror targetRecvTm = rTm;
        if (targetRecvTm instanceof TypeVariable tv && tv.getUpperBound() != null) {
          targetRecvTm = tv.getUpperBound();
        } else if (targetRecvTm instanceof WildcardType wt && wt.getExtendsBound() != null) {
          targetRecvTm = wt.getExtendsBound();
        }
        recvExpr = convert(emitPath(recvPath, ctx), rTm, targetRecvTm, false, ctx);
      } else {
        recvExpr = resolveImplicitThisValue(owner, ctx);
      }
      String kGo =
          strKey ? emitPath(arg0Path, ctx) : ("uint32(" + emitPath(arg0Path, ctx) + ")");
      if (m.equals("containsKey")
          && arity == 1
          && ee.getParameters().size() == 1
          && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
        return recvExpr + (strKey ? ".ContainsKeyStr_(" : ".ContainsKeyNode_(") + kGo + ")";
      }
      if (m.equals("containsEntry")
          && arity == 2
          && ee.getParameters().size() == 2
          && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
        String vGo = emitAs(args.get(1), path, objTm, true, ctx);
        return recvExpr
            + (strKey ? ".ContainsEntryStr_(" : ".ContainsEntryNode_(")
            + kGo
            + ", "
            + vGo
            + ")";
      }
      if (m.equals("put")
          && arity == 2
          && ee.getParameters().size() == 2
          && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
        String vGo = emitAs(args.get(1), path, objTm, true, ctx);
        return recvExpr + (strKey ? ".PutStr_(" : ".PutNode_(") + kGo + ", " + vGo + ")";
      }
      if (m.equals("remove")
          && arity == 2
          && ee.getParameters().size() == 2
          && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
        String vGo = emitAs(args.get(1), path, objTm, true, ctx);
        return recvExpr + (strKey ? ".RemoveStr_(" : ".RemoveNode_(") + kGo + ", " + vGo + ")";
      }
      if ((m.equals("get") || m.equals("removeAll"))
          && arity == 1
          && ee.getParameters().size() == 1) {
        String fn =
            m.equals("get")
                ? (strKey ? ".GetStr_(" : ".GetNode_(")
                : (strKey ? ".RemoveAllStr_(" : ".RemoveAllNode_(");
        String raw = recvExpr + fn + kGo + ")";
        TypeElement listTe = resolver.elements.getTypeElement("java.util.List");
        return (!isStmt && callRet.getKind() != TypeKind.VOID && listTe != null)
            ? convert(raw, listTe.asType(), callRet, false, ctx)
            : raw;
      }
      if (m.equals("putAll")
          && arity == 2
          && ee.getParameters().size() == 2
          && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
        TypeElement iterTe = resolver.elements.getTypeElement("java.lang.Iterable");
        if (iterTe != null
            && resolver.safeErasure(ee.getParameters().get(1).asType()) instanceof DeclaredType pDt
            && pDt.asElement().equals(iterTe)) {
          String valsGo = emitAs(args.get(1), path, iterTe.asType(), false, ctx);
          return recvExpr + (strKey ? ".PutAllStr_(" : ".PutAllNode_(") + kGo + ", " + valsGo + ")";
        }
      }
      return null;
    }

    if (isSubtypeOfQualified(owner, "com.google.common.collect.Multimap")) {
      String recvExpr =
          recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
      String kGo = emitPath(arg0Path, ctx);
      if (m.equals("containsKey")
          && arity == 1
          && ee.getParameters().size() == 1
          && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
        return rt
            + (strKey ? "MultimapContainsKeyStr(" : "MultimapContainsKeyNode(")
            + recvExpr
            + ", "
            + kGo
            + ")";
      }
      if (m.equals("containsEntry")
          && arity == 2
          && ee.getParameters().size() == 2
          && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
        String vGo = emitAs(args.get(1), path, objTm, true, ctx);
        return rt
            + (strKey ? "MultimapContainsEntryStr(" : "MultimapContainsEntryNode(")
            + recvExpr
            + ", "
            + kGo
            + ", "
            + vGo
            + ")";
      }
      if (m.equals("get") && arity == 1 && ee.getParameters().size() == 1) {
        String raw =
            rt + (strKey ? "MultimapGetStr(" : "MultimapGetNode(") + recvExpr + ", " + kGo + ")";
        return (!isStmt && callRet.getKind() != TypeKind.VOID)
            ? convert(raw, objTm, callRet, false, ctx)
            : raw;
      }
      if (m.equals("put")
          && arity == 2
          && ee.getParameters().size() == 2
          && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
        String vGo = emitAs(args.get(1), path, objTm, true, ctx);
        return rt
            + (strKey ? "MultimapPutStr(" : "MultimapPutNode(")
            + recvExpr
            + ", "
            + kGo
            + ", "
            + vGo
            + ")";
      }
      if (m.equals("remove")
          && arity == 2
          && ee.getParameters().size() == 2
          && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
        String vGo = emitAs(args.get(1), path, objTm, true, ctx);
        return rt
            + (strKey ? "MultimapRemoveStr(" : "MultimapRemoveNode(")
            + recvExpr
            + ", "
            + kGo
            + ", "
            + vGo
            + ")";
      }
      if (m.equals("removeAll") && arity == 1 && ee.getParameters().size() == 1) {
        String raw =
            rt
                + (strKey ? "MultimapRemoveAllStr(" : "MultimapRemoveAllNode(")
                + recvExpr
                + ", "
                + kGo
                + ")";
        return (!isStmt && callRet.getKind() != TypeKind.VOID)
            ? convert(raw, objTm, callRet, false, ctx)
            : raw;
      }
      if (m.equals("putAll")
          && arity == 2
          && ee.getParameters().size() == 2
          && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
        TypeElement iterTe = resolver.elements.getTypeElement("java.lang.Iterable");
        if (iterTe != null
            && resolver.types.isSubtype(
                resolver.safeErasure(ee.getParameters().get(1).asType()),
                resolver.safeErasure(iterTe.asType()))) {
          String valsGo = emitAs(args.get(1), path, iterTe.asType(), false, ctx);
          return rt
              + (strKey ? "MultimapPutAllStr(" : "MultimapPutAllNode(")
              + recvExpr
              + ", "
              + kGo
              + ", "
              + valsGo
              + ")";
        }
      }
      return null;
    }

    if (strKey && isSubtypeOfQualified(owner, "com.google.common.collect.Multiset")) {
      String recvExpr =
          recvPath != null ? emitPath(recvPath, ctx) : resolveImplicitThisValue(owner, ctx);
      String kGo = emitPath(arg0Path, ctx);
      if (m.equals("count")
          && arity == 1
          && ee.getParameters().size() == 1
          && ee.getReturnType().getKind() == TypeKind.INT) {
        return rt + "MultisetCountStr(" + recvExpr + ", " + kGo + ")";
      }
      if (m.equals("contains")
          && arity == 1
          && ee.getParameters().size() == 1
          && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
        return rt + "MultisetContainsStr(" + recvExpr + ", " + kGo + ")";
      }
      if (m.equals("add")
          && arity == 1
          && ee.getParameters().size() == 1
          && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
        return rt + "MultisetAddStr(" + recvExpr + ", " + kGo + ")";
      }
      if (m.equals("add")
          && arity == 2
          && ee.getParameters().size() == 2
          && ee.getParameters().get(1).asType().getKind() == TypeKind.INT
          && ee.getReturnType().getKind() == TypeKind.INT) {
        String nGo =
            emitAs(args.get(1), path, resolver.types.getPrimitiveType(TypeKind.INT), false, ctx);
        return rt + "MultisetAddCountStr(" + recvExpr + ", " + kGo + ", " + nGo + ")";
      }
      if (m.equals("remove")
          && arity == 1
          && ee.getParameters().size() == 1
          && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
        return rt + "MultisetRemoveStr(" + recvExpr + ", " + kGo + ")";
      }
      if (m.equals("setCount")
          && arity == 3
          && ee.getParameters().size() == 3
          && ee.getParameters().get(1).asType().getKind() == TypeKind.INT
          && ee.getParameters().get(2).asType().getKind() == TypeKind.INT
          && ee.getReturnType().getKind() == TypeKind.BOOLEAN) {
        TypeMirror intTm = resolver.types.getPrimitiveType(TypeKind.INT);
        String oldGo = emitAs(args.get(1), path, intTm, false, ctx);
        String newGo = emitAs(args.get(2), path, intTm, false, ctx);
        return rt + "MultisetSetCount3Str(" + recvExpr + ", " + kGo + ", " + oldGo + ", " + newGo + ")";
      }
      return null;
    }

    return null;
  }

  private String emitStringIntrinsic(
      MethodInvocationTree mi,
      TreePath path,
      ExecutableElement ee,
      boolean isStatic,
      TreePath recvPath,
      EmitContext ctx) {
    String m = ee.getSimpleName().toString();
    List<? extends ExpressionTree> args = mi.getArguments();
    int arity = args.size();
    String jl = jlPrefix(ctx);
    String rt = rtPrefix(ctx);
    TypeMirror intTm = resolver.types.getPrimitiveType(TypeKind.INT);
    TypeMirror charTm = resolver.types.getPrimitiveType(TypeKind.CHAR);

    if (isStatic) {
      if (m.equals("valueOf") && arity == 1) {
        TypeMirror p0 = resolver.safeErasure(ee.getParameters().get(0).asType());
        if (p0.getKind() == TypeKind.DOUBLE) {
          TypeMirror dblTm = resolver.types.getPrimitiveType(TypeKind.DOUBLE);
          return jl + "DoubleToString(" + emitAs(args.get(0), path, dblTm, false, ctx) + ")";
        }
        if (p0.getKind() == TypeKind.FLOAT) {
          TypeMirror fltTm = resolver.types.getPrimitiveType(TypeKind.FLOAT);
          return jl + "FloatToString(" + emitAs(args.get(0), path, fltTm, false, ctx) + ")";
        }
        if (p0.getKind() == TypeKind.INT) {
          ctx.imports.std("strconv");
          return "strconv.FormatInt(int64(" + emitAs(args.get(0), path, intTm, false, ctx) + "), 10)";
        }
        if (p0.getKind() == TypeKind.LONG) {
          TypeMirror longTm = resolver.types.getPrimitiveType(TypeKind.LONG);
          ctx.imports.std("strconv");
          return "strconv.FormatInt(" + emitAs(args.get(0), path, longTm, false, ctx) + ", 10)";
        }
        if (p0.getKind() == TypeKind.BOOLEAN) {
          TypeMirror boolTm = resolver.types.getPrimitiveType(TypeKind.BOOLEAN);
          return rt + "BoolToStr(" + emitAs(args.get(0), path, boolTm, false, ctx) + ")";
        }
        if (p0.getKind() == TypeKind.CHAR) {
          return jl + "StringFromChar(" + emitAs(args.get(0), path, charTm, false, ctx) + ")";
        }
        if (p0 instanceof ArrayType at && at.getComponentType().getKind() == TypeKind.CHAR) {
          return jl + "StringFromUTF16(" + emit(args.get(0), path, ctx) + ")";
        }
        TypeMirror argTm = resolver.safeErasure(typeOf(new TreePath(path, args.get(0))));
        if (resolver.isStringType(argTm)) {
          return rt + "StrValString(" + emit(args.get(0), path, ctx) + ")";
        }
        return rt + "StrVal(" + emit(args.get(0), path, ctx) + ")";
      }
      return null;
    }

    ExpressionTree recvTree = recvPath != null ? (ExpressionTree) recvPath.getLeaf() : null;
    String recv = recvPath != null ? emitPath(recvPath, ctx) : ctx.selfValueExpr();
    if (m.equals("length") && arity == 0) return rt + "JStrLen(" + recv + ")";
    if (m.equals("isEmpty") && arity == 0) {
      return "(" + checkNotNullStringExpr(recvTree, recv, ctx) + " == \"\")";
    }
    if (m.equals("charAt") && arity == 1) {
      return rt + "JStrCharAt(" + recv + ", " + emitAs(args.get(0), path, intTm, false, ctx) + ")";
    }
    if (m.equals("codePointAt") && arity == 1) {
      return rt + "JStrCodePointAt(" + recv + ", " + emitAs(args.get(0), path, intTm, false, ctx) + ")";
    }
    if (m.equals("substring") && arity == 1) {
      return rt + "JStrSubFrom(" + recv + ", " + emitAs(args.get(0), path, intTm, false, ctx) + ")";
    }
    if (m.equals("substring") && arity == 2) {
      return rt
          + "JStrSub("
          + recv
          + ", "
          + emitAs(args.get(0), path, intTm, false, ctx)
          + ", "
          + emitAs(args.get(1), path, intTm, false, ctx)
          + ")";
    }
    if (m.equals("subSequence") && arity == 2) {
      return rt
          + "JStrSub("
          + recv
          + ", "
          + emitAs(args.get(0), path, intTm, false, ctx)
          + ", "
          + emitAs(args.get(1), path, intTm, false, ctx)
          + ")";
    }
    if (m.equals("indexOf") && arity == 1) {
      TypeMirror p0 = resolver.safeErasure(ee.getParameters().get(0).asType());
      if (p0.getKind() == TypeKind.INT) {
        return rt + "JStrIndexOfCp(" + recv + ", " + emitAs(args.get(0), path, intTm, false, ctx) + ")";
      }
      return rt + "JStrIndexOf(" + recv + ", " + emit(args.get(0), path, ctx) + ")";
    }
    if (m.equals("indexOf") && arity == 2) {
      TypeMirror p0 = resolver.safeErasure(ee.getParameters().get(0).asType());
      if (p0.getKind() == TypeKind.INT) {
        return rt
            + "JStrIndexOfCpFrom("
            + recv
            + ", "
            + emitAs(args.get(0), path, intTm, false, ctx)
            + ", "
            + emitAs(args.get(1), path, intTm, false, ctx)
            + ")";
      }
      return rt
          + "JStrIndexOfFrom("
          + recv
          + ", "
          + emit(args.get(0), path, ctx)
          + ", "
          + emitAs(args.get(1), path, intTm, false, ctx)
          + ")";
    }
    if (m.equals("lastIndexOf") && arity == 1) {
      TypeMirror p0 = resolver.safeErasure(ee.getParameters().get(0).asType());
      if (p0.getKind() == TypeKind.INT) {
        return rt + "JStrLastIndexOfCp(" + recv + ", " + emitAs(args.get(0), path, intTm, false, ctx) + ")";
      }
      return rt + "JStrLastIndexOf(" + recv + ", " + emit(args.get(0), path, ctx) + ")";
    }
    if (m.equals("lastIndexOf") && arity == 2) {
      TypeMirror p0 = resolver.safeErasure(ee.getParameters().get(0).asType());
      if (p0.getKind() == TypeKind.INT) {
        return rt
            + "JStrLastIndexOfCpFrom("
            + recv
            + ", "
            + emitAs(args.get(0), path, intTm, false, ctx)
            + ", "
            + emitAs(args.get(1), path, intTm, false, ctx)
            + ")";
      }
      return rt
          + "JStrLastIndexOfFrom("
          + recv
          + ", "
          + emit(args.get(0), path, ctx)
          + ", "
          + emitAs(args.get(1), path, intTm, false, ctx)
          + ")";
    }
    if (m.equals("startsWith") && arity == 1) {
      return jl
          + "StringStartsWith("
          + checkNotNullStringExpr(recvTree, recv, ctx)
          + ", "
          + checkNotNullStringExpr(args.get(0), emit(args.get(0), path, ctx), ctx)
          + ")";
    }
    if (m.equals("startsWith") && arity == 2) {
      return rt
          + "JStrStartsWithOffset("
          + recv
          + ", "
          + emit(args.get(0), path, ctx)
          + ", "
          + emitAs(args.get(1), path, intTm, false, ctx)
          + ")";
    }
    if (m.equals("endsWith") && arity == 1) {
      return jl
          + "StringEndsWith("
          + checkNotNullStringExpr(recvTree, recv, ctx)
          + ", "
          + checkNotNullStringExpr(args.get(0), emit(args.get(0), path, ctx), ctx)
          + ")";
    }
    if (m.equals("contains") && arity == 1) {
      TypeMirror arg0Tm = resolver.safeErasure(typeOf(new TreePath(path, args.get(0))));
      if (resolver.isStringType(arg0Tm)) {
        return "("
            + rt
            + "JStrIndexOf("
            + recv
            + ", "
            + checkNotNullStringExpr(args.get(0), emit(args.get(0), path, ctx), ctx)
            + ") >= 0)";
      }
      return "(" + rt + "JStrIndexOf(" + recv + ", " + rt + "StrVal(" + rt + "CheckNotNull(" + emit(args.get(0), path, ctx) + "))) >= 0)";
    }
    if (m.equals("equals") && arity == 1) {
      TreePath arg0Path = new TreePath(path, args.get(0));
      TypeMirror arg0Tm = resolver.safeErasure(typeOf(arg0Path));
      if (resolver.isStringType(arg0Tm)) {
        String a0 = emitPath(arg0Path, ctx);
        if (args.get(0).getKind() == Tree.Kind.STRING_LITERAL
            || isSimpleLocalOrParamIdent(arg0Path, a0)) {
          return "(" + checkNotNullStringExpr(recvTree, recv, ctx) + " == " + a0 + ")";
        }
        String r0 = ctx.nextTemp("sr");
        String t0 = ctx.nextTemp("se");
        return "func() bool { "
            + r0
            + " := "
            + recv
            + "; "
            + t0
            + " := "
            + a0
            + "; return "
            + rt
            + "CheckNotNullString("
            + r0
            + ") == "
            + t0
            + " }()";
      }
      return rt + "StringEquals(" + recv + ", " + emitAs(args.get(0), path, resolver.getObjectElement().asType(), true, ctx) + ")";
    }
    if (m.equals("equalsIgnoreCase") && arity == 1) {
      ctx.imports.std("strings");
      return "strings.EqualFold(" + checkNotNullStringExpr(recvTree, recv, ctx) + ", " + emit(args.get(0), path, ctx) + ")";
    }
    if (m.equals("compareTo") && arity == 1) {
      TypeMirror strTm = resolver.elements.getTypeElement("java.lang.String").asType();
      return rt + "JStrCompareTo(" + recv + ", " + emitAs(args.get(0), path, strTm, false, ctx) + ")";
    }
    if (m.equals("hashCode") && arity == 0) {
      return jl + "StringHashCode(" + checkNotNullStringExpr(recvTree, recv, ctx) + ")";
    }
    if (m.equals("toString") && arity == 0) {
      return checkNotNullStringExpr(recvTree, recv, ctx);
    }
    if (m.equals("intern") && arity == 0) {
      return checkNotNullStringExpr(recvTree, recv, ctx);
    }
    if (m.equals("trim") && arity == 0) {
      return jl + "StringTrim(" + checkNotNullStringExpr(recvTree, recv, ctx) + ")";
    }
    if (m.equals("replace") && arity == 2) {
      TypeMirror p0 = resolver.safeErasure(ee.getParameters().get(0).asType());
      if (p0.getKind() == TypeKind.CHAR) {
        return rt
            + "JStrReplaceChar("
            + recv
            + ", "
            + emitAs(args.get(0), path, charTm, false, ctx)
            + ", "
            + emitAs(args.get(1), path, charTm, false, ctx)
            + ")";
      }
      TypeMirror a0Tm = resolver.safeErasure(typeOf(new TreePath(path, args.get(0))));
      TypeMirror a1Tm = resolver.safeErasure(typeOf(new TreePath(path, args.get(1))));
      String s0 =
          resolver.isStringType(a0Tm)
              ? checkNotNullStringExpr(args.get(0), emit(args.get(0), path, ctx), ctx)
              : (rt + "StrVal(" + rt + "CheckNotNull(" + emit(args.get(0), path, ctx) + "))");
      String s1 =
          resolver.isStringType(a1Tm)
              ? checkNotNullStringExpr(args.get(1), emit(args.get(1), path, ctx), ctx)
              : (rt + "StrVal(" + rt + "CheckNotNull(" + emit(args.get(1), path, ctx) + "))");
      return jl + "StringReplace(" + recv + ", " + s0 + ", " + s1 + ")";
    }
    if (m.equals("toLowerCase") && arity <= 1) {
      return jl + "StringToLowerCaseASCII(" + recv + ")";
    }
    if (m.equals("toUpperCase") && arity <= 1) {
      return jl + "StringToUpperCaseASCII(" + recv + ")";
    }
    if (m.equals("toCharArray") && arity == 0) {
      return jl + "StringToUTF16(" + recv + ")";
    }
    if (m.equals("repeat") && arity == 1) {
      ctx.imports.std("strings");
      return "strings.Repeat(" + recv + ", int(" + emitAs(args.get(0), path, intTm, false, ctx) + "))";
    }
    if (m.equals("split") && arity == 1) {
      return rt + "JStrSplit(" + recv + ", " + emit(args.get(0), path, ctx) + ")";
    }
    if (m.equals("split") && arity == 2) {
      return rt
          + "JStrSplitLimit("
          + recv
          + ", "
          + emit(args.get(0), path, ctx)
          + ", "
          + emitAs(args.get(1), path, intTm, false, ctx)
          + ")";
    }
    if (m.equals("matches") && arity == 1) {
      return rt + "JStrMatches(" + recv + ", " + emit(args.get(0), path, ctx) + ")";
    }
    if (m.equals("replaceAll") && arity == 2) {
      return rt
          + "JStrReplaceAll("
          + recv
          + ", "
          + emit(args.get(0), path, ctx)
          + ", "
          + emit(args.get(1), path, ctx)
          + ")";
    }
    if (m.equals("replaceFirst") && arity == 2) {
      return rt
          + "JStrReplaceFirst("
          + recv
          + ", "
          + emit(args.get(0), path, ctx)
          + ", "
          + emit(args.get(1), path, ctx)
          + ")";
    }
    if (m.equals("regionMatches") && arity == 4) {
      return rt
          + "JStrRegionMatches("
          + recv
          + ", false, "
          + emitAs(args.get(0), path, intTm, false, ctx)
          + ", "
          + emit(args.get(1), path, ctx)
          + ", "
          + emitAs(args.get(2), path, intTm, false, ctx)
          + ", "
          + emitAs(args.get(3), path, intTm, false, ctx)
          + ")";
    }
    if (m.equals("regionMatches") && arity == 5) {
      return rt
          + "JStrRegionMatches("
          + recv
          + ", "
          + emitAs(args.get(0), path, resolver.types.getPrimitiveType(TypeKind.BOOLEAN), false, ctx)
          + ", "
          + emitAs(args.get(1), path, intTm, false, ctx)
          + ", "
          + emit(args.get(2), path, ctx)
          + ", "
          + emitAs(args.get(3), path, intTm, false, ctx)
          + ", "
          + emitAs(args.get(4), path, intTm, false, ctx)
          + ")";
    }
    if (m.equals("getBytes") && arity <= 1) {
      if (arity == 1) {
        return "func() []int8 { _ = "
            + emit(args.get(0), path, ctx)
            + "; return "
            + rt
            + "JStrGetBytes("
            + recv
            + ") }()";
      }
      return rt + "JStrGetBytes(" + recv + ")";
    }
    if (m.equals("concat") && arity == 1) {
      return "("
          + checkNotNullStringExpr(recvTree, recv, ctx)
          + " + "
          + checkNotNullStringExpr(args.get(0), emit(args.get(0), path, ctx), ctx)
          + ")";
    }
    if ((m.equals("strip") || m.equals("stripLeading") || m.equals("stripTrailing")) && arity == 0) {
      ctx.imports.std("strings");
      return "strings.TrimSpace(" + checkNotNullStringExpr(recvTree, recv, ctx) + ")";
    }
    if (m.equals("isBlank") && arity == 0) {
      ctx.imports.std("strings");
      return "(strings.TrimSpace(" + checkNotNullStringExpr(recvTree, recv, ctx) + ") == \"\")";
    }
    if (m.equals("compareToIgnoreCase") && arity == 1) {
      return rt
          + "JStrCompareTo("
          + jl
          + "StringToLowerCaseASCII("
          + checkNotNullStringExpr(recvTree, recv, ctx)
          + "), "
          + jl
          + "StringToLowerCaseASCII("
          + checkNotNullStringExpr(args.get(0), emit(args.get(0), path, ctx), ctx)
          + "))";
    }
    if (m.equals("contentEquals") && arity == 1) {
      TypeMirror arg0Tm = resolver.safeErasure(typeOf(new TreePath(path, args.get(0))));
      if (resolver.isStringType(arg0Tm)) {
        return "("
            + checkNotNullStringExpr(recvTree, recv, ctx)
            + " == "
            + checkNotNullStringExpr(args.get(0), emit(args.get(0), path, ctx), ctx)
            + ")";
      }
      return "(" + checkNotNullStringExpr(recvTree, recv, ctx) + " == " + rt + "StrVal(" + rt + "CheckNotNull(" + emit(args.get(0), path, ctx) + ")))";
    }
    return null;
  }

  public boolean hasIdentityEquals(TypeElement te) {
    if (te == null
        || resolver.isSubclassed(te)
        || te.getKind() != ElementKind.CLASS
        || resolver.isNodeClass(te)
        || resolver.isValueStructClass(te)
        || !resolver.isSourceType(te)) {
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
      if (!(supTm instanceof DeclaredType supDt)
          || !(supDt.asElement() instanceof TypeElement ste)) {
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

  private static boolean isColorIdType(TypeMirror er) {
    return er instanceof DeclaredType dt
        && dt.asElement() instanceof TypeElement te
        && te.getQualifiedName().contentEquals("com.google.javascript.jscomp.colors.ColorId");
  }

  public boolean valueStructGoEqMatchesJavaEquals(TypeElement te) {
    return valueStructHasIdentityEquals(te)
        || (te != null
            && te.getQualifiedName()
                .contentEquals("com.google.javascript.jscomp.colors.ColorId"));
  }

  public boolean valueStructHasIdentityEquals(TypeElement te) {
    if (te == null
        || resolver.isSubclassed(te)
        || te.getKind() != ElementKind.CLASS
        || resolver.isNodeClass(te)
        || !resolver.isValueStructClass(te)
        || !resolver.isSourceType(te)) {
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
      if (!(supTm instanceof DeclaredType supDt)
          || !(supDt.asElement() instanceof TypeElement ste)) {
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

  private boolean hasConcreteValueStructHashCodeMethod(TypeElement te) {
    if (te == null
        || resolver.isSubclassed(te)
        || te.getKind() != ElementKind.CLASS
        || !resolver.isValueStructClass(te)
        || !resolver.isSourceType(te)) {
      return false;
    }
    for (Element m : te.getEnclosedElements()) {
      if (m instanceof ExecutableElement ee
          && ee.getKind() == ElementKind.METHOD
          && !ee.getModifiers().contains(Modifier.STATIC)
          && !ee.getModifiers().contains(Modifier.ABSTRACT)
          && !ee.getModifiers().contains(Modifier.PRIVATE)
          && ee.getParameters().isEmpty()) {
        TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
        if (sig != null
            && sig.goName.equals("HashCode")
            && sig.canonParamTypes.isEmpty()
            && sig.canonReturnType != null
            && sig.canonReturnType.getKind() == TypeKind.INT
            && !sig.returnWidenedToAny) {
          return true;
        }
      }
    }
    return false;
  }

  private boolean hasConcreteGoHashCodeMethod(TypeElement te) {
    if (te == null
        || resolver.isSubclassed(te)
        || te.getKind() != ElementKind.CLASS
        || resolver.isNodeClass(te)
        || resolver.isValueStructClass(te)
        || !resolver.isSourceType(te)) {
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
            && !ee.getModifiers().contains(Modifier.STATIC)
            && !ee.getModifiers().contains(Modifier.ABSTRACT)
            && !ee.getModifiers().contains(Modifier.PRIVATE)
            && ee.getParameters().isEmpty()) {
          TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
          if (sig != null
              && sig.goName.equals("HashCode")
              && sig.canonParamTypes.isEmpty()
              && sig.canonReturnType != null
              && sig.canonReturnType.getKind() == TypeKind.INT
              && !sig.returnWidenedToAny) {
            return true;
          }
        }
      }
      TypeMirror supTm = cur.getSuperclass();
      if (!(supTm instanceof DeclaredType supDt)
          || !(supDt.asElement() instanceof TypeElement ste)) {
        break;
      }
      cur = ste;
    }
    return false;
  }

  private boolean ifaceHasGoHashCodeMethod(TypeElement te) {
    if (te == null || !resolver.isSourceType(te) || resolver.runtime.mapsToAny(te)) {
      return false;
    }
    List<TypeElement> toCheck = new ArrayList<>();
    toCheck.add(te);
    for (TypeElement anc : resolver.getAllAncestorTypes(te)) {
      if (resolver.isSourceType(anc) && !resolver.runtime.mapsToAny(anc)) {
        toCheck.add(anc);
      }
    }
    for (TypeElement cur : toCheck) {
      for (Element m : cur.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)
            && !ee.getModifiers().contains(Modifier.PRIVATE)
            && ee.getParameters().isEmpty()) {
          TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
          if (sig != null
              && sig.goName.equals("HashCode")
              && sig.canonParamTypes.isEmpty()
              && sig.canonReturnType != null
              && sig.canonReturnType.getKind() == TypeKind.INT
              && !sig.returnWidenedToAny) {
            return true;
          }
        }
      }
    }
    return false;
  }

  private boolean hasConcreteGoEqualsMethod(TypeElement te) {
    if (te == null
        || resolver.isSubclassed(te)
        || te.getKind() != ElementKind.CLASS
        || resolver.isNodeClass(te)
        || resolver.isValueStructClass(te)
        || !resolver.isSourceType(te)) {
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
            && !ee.getModifiers().contains(Modifier.STATIC)
            && !ee.getModifiers().contains(Modifier.ABSTRACT)
            && !ee.getModifiers().contains(Modifier.PRIVATE)
            && ee.getParameters().size() == 1) {
          TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
          if (sig != null
              && sig.goName.equals("Equals")
              && sig.canonParamTypes.size() == 1
              && sig.canonParamTypes.get(0) instanceof DeclaredType pdt
              && pdt.asElement() instanceof TypeElement pte
              && pte.getQualifiedName().contentEquals("java.lang.Object")
              && sig.canonReturnType != null
              && sig.canonReturnType.getKind() == TypeKind.BOOLEAN
              && !sig.returnWidenedToAny) {
            return true;
          }
        }
      }
      TypeMirror supTm = cur.getSuperclass();
      if (!(supTm instanceof DeclaredType supDt)
          || !(supDt.asElement() instanceof TypeElement ste)) {
        break;
      }
      cur = ste;
    }
    return false;
  }

  private boolean ifaceHasGoEqualsMethod(TypeElement te) {
    if (te == null || !resolver.isSourceType(te) || resolver.runtime.mapsToAny(te)) {
      return false;
    }
    List<TypeElement> toCheck = new ArrayList<>();
    toCheck.add(te);
    for (TypeElement anc : resolver.getAllAncestorTypes(te)) {
      if (resolver.isSourceType(anc) && !resolver.runtime.mapsToAny(anc)) {
        toCheck.add(anc);
      }
    }
    for (TypeElement cur : toCheck) {
      for (Element m : cur.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)
            && !ee.getModifiers().contains(Modifier.PRIVATE)
            && ee.getParameters().size() == 1) {
          TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
          if (sig != null
              && sig.goName.equals("Equals")
              && sig.canonParamTypes.size() == 1
              && sig.canonParamTypes.get(0) instanceof DeclaredType pdt
              && pdt.asElement() instanceof TypeElement pte
              && pte.getQualifiedName().contentEquals("java.lang.Object")
              && sig.canonReturnType != null
              && sig.canonReturnType.getKind() == TypeKind.BOOLEAN
              && !sig.returnWidenedToAny) {
            return true;
          }
        }
      }
    }
    return false;
  }

  private String emitInlineNullCheck(
      boolean isNull, ExpressionTree arg, TreePath callPath, EmitContext ctx) {
    TreePath p0 = new TreePath(callPath, arg);
    TypeMirror raw = typeOf(p0);
    TypeMirror er = resolver.safeErasure(raw);
    if (er.getKind() == TypeKind.NULL) {
      return isNull ? "true" : "false";
    }
    if (er.getKind().isPrimitive()) {
      return "func() bool { _ = "
          + emitPath(p0, ctx)
          + "; return "
          + (isNull ? "false" : "true")
          + " }()";
    }
    return emitEquality(isNull, p0, raw, null, resolver.types.getNullType(), ctx);
  }

  private String emitInlineObjectsEquals(
      ExpressionTree aTree, ExpressionTree bTree, TreePath callPath, EmitContext ctx) {
    TreePath p0 = unwrapBoxValueOf(new TreePath(callPath, aTree));
    TreePath p1 = unwrapBoxValueOf(new TreePath(callPath, bTree));
    TypeMirror raw0 = typeOf(p0);
    TypeMirror raw1 = typeOf(p1);
    TypeMirror t0 = resolver.safeErasure(raw0);
    TypeMirror t1 = resolver.safeErasure(raw1);

    if (t0.getKind() == TypeKind.NULL && t1.getKind() == TypeKind.NULL) {
      return "true";
    }
    if (t0.getKind() == TypeKind.NULL || t1.getKind() == TypeKind.NULL) {
      TreePath nonNullP = t0.getKind() == TypeKind.NULL ? p1 : p0;
      TypeMirror nonNullEr = t0.getKind() == TypeKind.NULL ? t1 : t0;
      if (nonNullEr.getKind().isPrimitive()) {
        return "func() bool { _ = " + emitPath(nonNullP, ctx) + "; return false }()";
      }
      return emitEquality(true, p0, raw0, p1, raw1, ctx);
    }
    if (t0.getKind() == t1.getKind()
        && t0.getKind().isPrimitive()
        && t0.getKind() != TypeKind.FLOAT
        && t0.getKind() != TypeKind.DOUBLE) {
      return "(" + emitPath(p0, ctx) + " == " + emitPath(p1, ctx) + ")";
    }
    if (resolver.isStringType(t0) && resolver.isStringType(t1)) {
      return "(" + emitPath(p0, ctx) + " == " + emitPath(p1, ctx) + ")";
    }
    if (resolver.isNodeType(t0) && resolver.isNodeType(t1)) {
      String l = convert(emitPath(p0, ctx), raw0, t0, false, ctx);
      String r = convert(emitPath(p1, ctx), raw1, t0, false, ctx);
      return "(" + l + " == " + r + ")";
    }
    if (resolver.isByteStringType(t0) && resolver.isByteStringType(t1)) {
      return rtPrefix(ctx) + "ByteStringObjEquals(" + emitPath(p0, ctx) + ", " + emitPath(p1, ctx) + ")";
    }
    if (resolver.isValueEnumType(t0) && resolver.types.isSameType(t0, t1)) {
      return "(" + emitPath(p0, ctx) + " == " + emitPath(p1, ctx) + ")";
    }
    if (resolver.isValueStructType(t0)
        && t0 instanceof DeclaredType d0s
        && t1 instanceof DeclaredType d1s
        && d0s.asElement() instanceof TypeElement te0s
        && te0s.equals(d1s.asElement())) {
      if (valueStructGoEqMatchesJavaEquals(te0s)) {
        return "(" + emitPath(p0, ctx) + " == " + emitPath(p1, ctx) + ")";
      }
      if (resolver.isNullableValueStructType(t0)) {
        String goT = resolver.goType(raw0, ctx.imports);
        String aTmp = ctx.nextTemp("va");
        String bTmp = ctx.nextTemp("vb");
        return "func() bool { "
            + aTmp
            + " := "
            + emitPath(p0, ctx)
            + "; "
            + bTmp
            + " := "
            + emitPath(p1, ctx)
            + "; if "
            + aTmp
            + " == "
            + bTmp
            + " { return true }; if "
            + aTmp
            + " == ("
            + goT
            + "{}) || "
            + bTmp
            + " == ("
            + goT
            + "{}) { return false }; return "
            + aTmp
            + ".Equals("
            + bTmp
            + ") }()";
      }
    }
    if (t0 instanceof DeclaredType d0
        && d0.asElement() instanceof TypeElement te0
        && t1 instanceof DeclaredType d1
        && d1.asElement() instanceof TypeElement te1) {
      if (te0.getKind() == ElementKind.ENUM && te1.getKind() == ElementKind.ENUM) {
        String lGo = resolver.goType(t0, ctx.imports);
        String rGo = resolver.goType(t1, ctx.imports);
        if (lGo.startsWith("*") && lGo.equals(rGo)) {
          return "(" + emitPath(p0, ctx) + " == " + emitPath(p1, ctx) + ")";
        }
      }
      String lGo = resolver.goType(raw0, ctx.imports);
      String rGo = resolver.goType(raw1, ctx.imports);
      if (te0.equals(te1) && lGo.equals(rGo) && hasIdentityEquals(te0)) {
        return "(" + emitPath(p0, ctx) + " == " + emitPath(p1, ctx) + ")";
      }
      if (te0.equals(te1)
          && lGo.startsWith("*")
          && lGo.equals(rGo)
          && ((te0.getKind() == ElementKind.RECORD && resolver.isSourceType(te0))
              || resolver.runtime.isBoxedPrimitive(te0)
              || hasConcreteGoEqualsMethod(te0))) {
        return rtPrefix(ctx) + "ObjEqCustomPtr(" + emitPath(p0, ctx) + ", " + emitPath(p1, ctx) + ")";
      }
      TypeMirror raw0Eff = resolver.goEffectiveType(raw0);
      TypeMirror raw1Eff = resolver.goEffectiveType(raw1);
      if (isSourceJavaNilIface(t1)
          && (lGo.equals(rGo)
              || (isSourceRefAssignableToIface(t0) && resolver.types.isAssignable(raw0Eff, raw1Eff)))) {
        String fn = ifaceHasGoEqualsMethod(te1) ? "ObjEqIface" : "ObjEqNilIface";
        return rtPrefix(ctx) + fn + "[" + rGo + "](" + emitPath(p0, ctx) + ", " + emitPath(p1, ctx) + ")";
      }
      if (isSourceJavaNilIface(t0)
          && isSourceRefAssignableToIface(t1)
          && resolver.types.isAssignable(raw1Eff, raw0Eff)) {
        String fn = ifaceHasGoEqualsMethod(te0) ? "ObjEqIface" : "ObjEqNilIface";
        return rtPrefix(ctx) + fn + "[" + lGo + "](" + emitPath(p0, ctx) + ", " + emitPath(p1, ctx) + ")";
      }
    }
    if (resolver.isSameActiveComparableTypeVariable(raw0, raw1)) {
      String goTv = resolver.goType(raw0, ctx.imports);
      String rTmp = ctx.nextTemp("oa");
      String aTmp = ctx.nextTemp("ob");
      return "func() bool { "
          + rTmp
          + " := "
          + emitPath(p0, ctx)
          + "; "
          + aTmp
          + " := "
          + emitPath(p1, ctx)
          + "; if "
          + rTmp
          + " == "
          + aTmp
          + " { return true }; if _, ok_ := any((*"
          + goTv
          + ")(nil)).(*string); ok_ { return false }; return "
          + jlPrefix(ctx)
          + "ObjectsEquals("
          + rtPrefix(ctx)
          + "ToAny("
          + rTmp
          + "), "
          + rtPrefix(ctx)
          + "ToAny("
          + aTmp
          + ")) }()";
    }

    String a0 = emitAs(aTree, callPath, resolver.getObjectElement().asType(), true, ctx);
    String a1 = emitAs(bTree, callPath, resolver.getObjectElement().asType(), true, ctx);
    return jlPrefix(ctx) + "ObjectsEquals(" + a0 + ", " + a1 + ")";
  }

  private TreePath unwrapBoxValueOf(TreePath argPath) {
    TreePath cur = argPath;
    while (cur.getLeaf() instanceof ParenthesizedTree pt) {
      cur = new TreePath(cur, pt.getExpression());
    }
    if (cur.getLeaf() instanceof MethodInvocationTree mi && mi.getArguments().size() == 1) {
      Element el = resolver.trees.getElement(cur);
      if (el instanceof ExecutableElement ee
          && ee.getModifiers().contains(Modifier.STATIC)
          && ee.getSimpleName().contentEquals("valueOf")
          && ee.getParameters().size() == 1
          && ee.getParameters().get(0).asType().getKind().isPrimitive()
          && ee.getEnclosingElement() instanceof TypeElement owner
          && resolver.runtime.isBoxedPrimitive(owner)) {
        TreePath sub = new TreePath(cur, mi.getArguments().get(0));
        TypeMirror subTm = typeOf(sub);
        if (subTm != null
            && subTm.getKind() == ee.getParameters().get(0).asType().getKind()) {
          return sub;
        }
      }
    }
    return argPath;
  }

  private static boolean isStringLiteralTree(Tree tree) {
    while (tree instanceof ParenthesizedTree pt) {
      tree = pt.getExpression();
    }
    return tree != null && tree.getKind() == Tree.Kind.STRING_LITERAL;
  }

  private String emitInlineSingleObjectHashCode(
      ExpressionTree argTree, TreePath callPath, EmitContext ctx) {
    TreePath argPath = unwrapBoxValueOf(new TreePath(callPath, argTree));
    TypeMirror rawTm = typeOf(argPath);
    TypeMirror er = resolver.safeErasure(rawTm);
    String rt = rtPrefix(ctx);
    String jl = jlPrefix(ctx);

    if (er.getKind() == TypeKind.NULL) {
      return "int32(0)";
    }
    if (er.getKind() == TypeKind.INT) {
      return emitPath(argPath, ctx);
    }
    if (er.getKind() == TypeKind.BYTE
        || er.getKind() == TypeKind.SHORT
        || er.getKind() == TypeKind.CHAR) {
      return "int32(" + emitPath(argPath, ctx) + ")";
    }
    if (er.getKind() == TypeKind.BOOLEAN) {
      return "func() int32 { if "
          + emitPath(argPath, ctx)
          + " { return 1231 }; return 1237 }()";
    }
    if (er.getKind() == TypeKind.LONG) {
      String tmp = ctx.nextTemp("hl");
      return "func() int32 { "
          + tmp
          + " := uint64("
          + emitPath(argPath, ctx)
          + "); return int32("
          + tmp
          + " ^ ("
          + tmp
          + " >> 32)) }()";
    }
    if (resolver.isStringType(er)) {
      String s = emitPath(argPath, ctx);
      if (isStringLiteralTree(argPath.getLeaf())) {
        return jl + "StringHashCode(" + s + ")";
      }
      String tmp = ctx.nextTemp("hs");
      return "func() int32 { "
          + tmp
          + " := "
          + s
          + "; if "
          + tmp
          + " == "
          + rt
          + "NullString { return 0 }; return "
          + jl
          + "StringHashCode("
          + tmp
          + ") }()";
    }
    if (resolver.isNodeType(er)) {
      return "int32(" + emitPath(argPath, ctx) + ")";
    }
    if (resolver.isValueEnumType(er)) {
      return "(" + emitPath(argPath, ctx) + ").HashCode()";
    }
    if (resolver.isByteStringType(er)) {
      return rt + "ByteStringObjHash(" + emitPath(argPath, ctx) + ")";
    }
    if (er instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      if (resolver.isNullableValueStructClass(te) && hasConcreteValueStructHashCodeMethod(te)) {
        String goT = resolver.goType(rawTm, ctx.imports);
        String tmp = ctx.nextTemp("hv");
        return "func() int32 { "
            + tmp
            + " := "
            + emitPath(argPath, ctx)
            + "; if "
            + tmp
            + " == ("
            + goT
            + "{}) { return 0 }; return "
            + tmp
            + ".HashCode() }()";
      }
      if (resolver.isValueStructClass(te)) {
        if (hasConcreteValueStructHashCodeMethod(te)) {
          return "(" + emitPath(argPath, ctx) + ").HashCode()";
        }
        if (valueStructHasIdentityEquals(te)) {
          return "(" + emitPath(argPath, ctx) + ").IdentityHash_()";
        }
      }
      if ((te.getKind() == ElementKind.ENUM && resolver.goType(er, ctx.imports).startsWith("*"))
          || (te.getKind() == ElementKind.RECORD && resolver.isSourceType(te))
          || resolver.runtime.isBoxedPrimitive(te)
          || hasConcreteGoHashCodeMethod(te)) {
        String tmp = ctx.nextTemp("hp");
        return "func() int32 { "
            + tmp
            + " := "
            + emitPath(argPath, ctx)
            + "; if "
            + tmp
            + " == nil { return 0 }; return "
            + tmp
            + ".HashCode() }()";
      }
      if (hasIdentityEquals(te)) {
        String tmp = ctx.nextTemp("hp");
        return "func() int32 { "
            + tmp
            + " := "
            + emitPath(argPath, ctx)
            + "; if "
            + tmp
            + " == nil { return 0 }; return "
            + rt
            + "IdentityHash("
            + tmp
            + ") }()";
      }
      if (isSourceJavaNilIface(er) && ifaceHasGoHashCodeMethod(te)) {
        String tmp = ctx.nextTemp("hi");
        return "func() int32 { "
            + tmp
            + " := "
            + emitPath(argPath, ctx)
            + "; if "
            + tmp
            + " == nil || "
            + tmp
            + ".IsJavaNil_() { return 0 }; return "
            + tmp
            + ".HashCode() }()";
      }
      if (resolver.isGoInterfaceType(er)) {
        return rt + "ObjectHash(" + emitPath(argPath, ctx) + ")";
      }
    }
    if (resolver.isActiveTypeVariable(rawTm)) {
      String goTv = resolver.goType(rawTm, ctx.imports);
      String rTmp = ctx.nextTemp("ho");
      return "func() int32 { "
          + rTmp
          + " := "
          + emitPath(argPath, ctx)
          + "; if _, ok_ := any((*"
          + goTv
          + ")(nil)).(*string); ok_ { s_ := any("
          + rTmp
          + ").(string); if s_ == "
          + rt
          + "NullString { return 0 }; return "
          + jl
          + "StringHashCode(s_) }; return "
          + rt
          + "ObjectHash("
          + rt
          + "ToAny("
          + rTmp
          + ")) }()";
    }
    return rt
        + "ObjectHash("
        + emitAs(argTree, callPath, resolver.getObjectElement().asType(), true, ctx)
        + ")";
  }

  private String emitInlineObjectsHash(
      List<? extends ExpressionTree> args,
      TreePath callPath,
      ExecutableElement ee,
      EmitContext ctx) {
    if (args.isEmpty()) {
      return "int32(1)";
    }
    if (args.size() == 1) {
      TreePath arg0Path = new TreePath(callPath, args.get(0));
      TypeMirror arg0Tm = typeOf(arg0Path);
      TypeMirror param0Tm = resolver.safeErasure(ee.getParameters().get(0).asType());
      if (arg0Tm.getKind() == TypeKind.NULL) {
        return "int32(0)";
      }
      if (arg0Tm.getKind() == TypeKind.ARRAY
          && resolver.types.isAssignable(resolver.safeErasure(arg0Tm), param0Tm)) {
        return rtPrefix(ctx)
            + "ObjectsHash("
            + emitAs(args.get(0), callPath, param0Tm, false, ctx)
            + "...)";
      }
    }

    String hVar = ctx.nextTemp("oh");
    StringBuilder sb = new StringBuilder("func() int32 { var ");
    sb.append(hVar).append(" int32 = 1; ");
    for (ExpressionTree arg : args) {
      TreePath argPath = unwrapBoxValueOf(new TreePath(callPath, arg));
      TypeMirror rawTm = typeOf(argPath);
      String valGo =
          rawTm.getKind() == TypeKind.NULL ? "nil" : emitPath(argPath, ctx);
      appendTypedHashAccum(
          sb, hVar, valGo, rawTm, isStringLiteralTree(argPath.getLeaf()), ctx);
      sb.append("; ");
    }
    sb.append("return ").append(hVar).append(" }()");
    return sb.toString();
  }

  public void appendTypedHashAccum(
      StringBuilder sb,
      String hVar,
      String emittedTypedExpr,
      TypeMirror rawTm,
      boolean isNonNullStringLit,
      EmitContext ctx) {
    TypeMirror er = resolver.safeErasure(rawTm);
    String rt = rtPrefix(ctx);
    String jl = jlPrefix(ctx);

    if (er.getKind() == TypeKind.NULL) {
      sb.append(hVar).append(" = 31 * ").append(hVar);
      return;
    }
    switch (er.getKind()) {
      case BOOLEAN -> {
        sb.append("if ")
            .append(emittedTypedExpr)
            .append(" { ")
            .append(hVar)
            .append(" = 31*")
            .append(hVar)
            .append(" + 1231 } else { ")
            .append(hVar)
            .append(" = 31*")
            .append(hVar)
            .append(" + 1237 }");
        return;
      }
      case INT -> {
        sb.append(hVar).append(" = 31*").append(hVar).append(" + (").append(emittedTypedExpr).append(')');
        return;
      }
      case BYTE, SHORT, CHAR -> {
        sb.append(hVar).append(" = 31*").append(hVar).append(" + int32(").append(emittedTypedExpr).append(')');
        return;
      }
      case LONG -> {
        String tmp = ctx.nextTemp("hl");
        sb.append("{ ")
            .append(tmp)
            .append(" := uint64(")
            .append(emittedTypedExpr)
            .append("); ")
            .append(hVar)
            .append(" = 31*")
            .append(hVar)
            .append(" + int32(")
            .append(tmp)
            .append(" ^ (")
            .append(tmp)
            .append(" >> 32)) }");
        return;
      }
      case FLOAT -> {
        ctx.imports.std("math");
        String tmp = ctx.nextTemp("hf");
        sb.append("{ ")
            .append(tmp)
            .append(" := ")
            .append(emittedTypedExpr)
            .append("; if math.IsNaN(float64(")
            .append(tmp)
            .append(")) { ")
            .append(hVar)
            .append(" = 31*")
            .append(hVar)
            .append(" + 0x7fc00000 } else { ")
            .append(hVar)
            .append(" = 31*")
            .append(hVar)
            .append(" + int32(math.Float32bits(")
            .append(tmp)
            .append(")) } }");
        return;
      }
      case DOUBLE -> {
        ctx.imports.std("math");
        String tmp = ctx.nextTemp("hd");
        String uTmp = ctx.nextTemp("hu");
        sb.append("{ ")
            .append(tmp)
            .append(" := ")
            .append(emittedTypedExpr)
            .append("; if math.IsNaN(")
            .append(tmp)
            .append(") { ")
            .append(hVar)
            .append(" = 31*")
            .append(hVar)
            .append(" + 0x7ff80000 } else { ")
            .append(uTmp)
            .append(" := math.Float64bits(")
            .append(tmp)
            .append("); ")
            .append(hVar)
            .append(" = 31*")
            .append(hVar)
            .append(" + int32(")
            .append(uTmp)
            .append(" ^ (")
            .append(uTmp)
            .append(" >> 32)) } }");
        return;
      }
      default -> {}
    }

    if (resolver.isStringType(er)) {
      if (isNonNullStringLit) {
        sb.append(hVar)
            .append(" = 31*")
            .append(hVar)
            .append(" + ")
            .append(jl)
            .append("StringHashCode(")
            .append(emittedTypedExpr)
            .append(')');
      } else {
        String tmp = ctx.nextTemp("hs");
        sb.append("{ ")
            .append(tmp)
            .append(" := ")
            .append(emittedTypedExpr)
            .append("; ")
            .append(hVar)
            .append(" = 31 * ")
            .append(hVar)
            .append("; if ")
            .append(tmp)
            .append(" != ")
            .append(rt)
            .append("NullString { ")
            .append(hVar)
            .append(" += ")
            .append(jl)
            .append("StringHashCode(")
            .append(tmp)
            .append(") } }");
      }
      return;
    }

    if (resolver.isNodeType(er)) {
      sb.append(hVar).append(" = 31*").append(hVar).append(" + int32(").append(emittedTypedExpr).append(')');
      return;
    }

    if (resolver.isValueEnumType(er)) {
      sb.append(hVar)
          .append(" = 31*")
          .append(hVar)
          .append(" + (")
          .append(emittedTypedExpr)
          .append(").HashCode()");
      return;
    }

    if (resolver.isByteStringType(er)) {
      sb.append(hVar)
          .append(" = 31*")
          .append(hVar)
          .append(" + ")
          .append(rt)
          .append("ByteStringObjHash(")
          .append(emittedTypedExpr)
          .append(')');
      return;
    }

    if (er instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      if (resolver.isNullableValueStructClass(te) && hasConcreteValueStructHashCodeMethod(te)) {
        String goT = resolver.goType(rawTm, ctx.imports);
        String tmp = ctx.nextTemp("hv");
        sb.append("{ ")
            .append(tmp)
            .append(" := ")
            .append(emittedTypedExpr)
            .append("; ")
            .append(hVar)
            .append(" = 31 * ")
            .append(hVar)
            .append("; if ")
            .append(tmp)
            .append(" != (")
            .append(goT)
            .append("{}) { ")
            .append(hVar)
            .append(" += ")
            .append(tmp)
            .append(".HashCode() } }");
        return;
      }
      if (resolver.isValueStructClass(te)) {
        if (hasConcreteValueStructHashCodeMethod(te)) {
          sb.append(hVar)
              .append(" = 31*")
              .append(hVar)
              .append(" + (")
              .append(emittedTypedExpr)
              .append(").HashCode()");
          return;
        }
        if (valueStructHasIdentityEquals(te)) {
          sb.append(hVar)
              .append(" = 31*")
              .append(hVar)
              .append(" + (")
              .append(emittedTypedExpr)
              .append(").IdentityHash_()");
          return;
        }
      }
      if ((te.getKind() == ElementKind.ENUM && resolver.goType(er, ctx.imports).startsWith("*"))
          || (te.getKind() == ElementKind.RECORD && resolver.isSourceType(te))
          || resolver.runtime.isBoxedPrimitive(te)
          || hasConcreteGoHashCodeMethod(te)) {
        String tmp = ctx.nextTemp("hp");
        sb.append("{ ")
            .append(tmp)
            .append(" := ")
            .append(emittedTypedExpr)
            .append("; ")
            .append(hVar)
            .append(" = 31 * ")
            .append(hVar)
            .append("; if ")
            .append(tmp)
            .append(" != nil { ")
            .append(hVar)
            .append(" += ")
            .append(tmp)
            .append(".HashCode() } }");
        return;
      }
      if (hasIdentityEquals(te)) {
        String tmp = ctx.nextTemp("hp");
        sb.append("{ ")
            .append(tmp)
            .append(" := ")
            .append(emittedTypedExpr)
            .append("; ")
            .append(hVar)
            .append(" = 31 * ")
            .append(hVar)
            .append("; if ")
            .append(tmp)
            .append(" != nil { ")
            .append(hVar)
            .append(" += ")
            .append(rt)
            .append("IdentityHash(")
            .append(tmp)
            .append(") } }");
        return;
      }
      if (isSourceJavaNilIface(er) && ifaceHasGoHashCodeMethod(te)) {
        String tmp = ctx.nextTemp("hi");
        sb.append("{ ")
            .append(tmp)
            .append(" := ")
            .append(emittedTypedExpr)
            .append("; ")
            .append(hVar)
            .append(" = 31 * ")
            .append(hVar)
            .append("; if ")
            .append(tmp)
            .append(" != nil && !")
            .append(tmp)
            .append(".IsJavaNil_() { ")
            .append(hVar)
            .append(" += ")
            .append(tmp)
            .append(".HashCode() } }");
        return;
      }
      if (resolver.isGoInterfaceType(er)) {
        sb.append(hVar)
            .append(" = 31*")
            .append(hVar)
            .append(" + ")
            .append(rt)
            .append("ObjectHash(")
            .append(emittedTypedExpr)
            .append(')');
        return;
      }
    }

    String boxed =
        convert(emittedTypedExpr, er, resolver.getObjectElement().asType(), true, ctx);
    sb.append(hVar)
        .append(" = 31*")
        .append(hVar)
        .append(" + ")
        .append(rt)
        .append("ObjectHash(")
        .append(boxed)
        .append(')');
  }

  public String tryEmitPreconditionsStmt(TreePath exprPath, EmitContext ctx) {
    if (!(exprPath.getLeaf() instanceof MethodInvocationTree mi)) {
      return null;
    }
    Element el = resolver.trees.getElement(exprPath);
    if (!(el instanceof ExecutableElement ee)
        || !ee.getModifiers().contains(Modifier.STATIC)
        || !(ee.getEnclosingElement() instanceof TypeElement owner)
        || !owner.getQualifiedName().contentEquals("com.google.common.base.Preconditions")) {
      return null;
    }
    String m = ee.getSimpleName().toString();
    if (!m.equals("checkState") && !m.equals("checkArgument")) {
      return null;
    }
    List<? extends ExpressionTree> args = mi.getArguments();
    int arity = args.size();
    if (arity == 0) {
      return null;
    }

    String fn = m.equals("checkState") ? "CheckState" : "CheckArgument";
    String rt = rtPrefix(ctx);
    TypeMirror boolTm = resolver.types.getPrimitiveType(TypeKind.BOOLEAN);

    if (arity == 1) {
      String cond = emitAs(args.get(0), exprPath, boolTm, false, ctx);
      return "if !(" + cond + ") { " + rt + fn + "(false) }";
    }

    Set<VariableElement> nonNullVars = new HashSet<>();
    collectGuaranteedNonNullVarsWhenTrue(new TreePath(exprPath, args.get(0)), nonNullVars);
    for (int i = 1; i < arity; i++) {
      if (!isPureDiagnosticArg(new TreePath(exprPath, args.get(i)), nonNullVars)) {
        return null;
      }
    }

    String cond = emitAs(args.get(0), exprPath, boolTm, false, ctx);
    TreePath arg1Path = new TreePath(exprPath, args.get(1));
    if (arity == 2) {
      if (resolver.isStringType(typeOf(arg1Path))) {
        return "if !("
            + cond
            + ") { "
            + rt
            + fn
            + "Str(false, "
            + emitPath(arg1Path, ctx)
            + ") }";
      }
      return "if !("
          + cond
          + ") { "
          + rt
          + fn
          + "Msg(false, "
          + emitAs(args.get(1), exprPath, resolver.getObjectElement().asType(), true, ctx)
          + ") }";
    }

    List<String> fmtArgs = new ArrayList<>();
    for (int i = 2; i < arity; i++) {
      fmtArgs.add(emitAs(args.get(i), exprPath, resolver.getObjectElement().asType(), true, ctx));
    }
    String tmplExpr =
        resolver.isStringType(typeOf(arg1Path))
            ? (rt + "StrValString(" + emitPath(arg1Path, ctx) + ")")
            : (rt + "StrVal(" + emitPath(arg1Path, ctx) + ")");
    return "if !("
        + cond
        + ") { "
        + rt
        + fn
        + "Fmt(false, "
        + tmplExpr
        + ", "
        + String.join(", ", fmtArgs)
        + ") }";
  }

  private void collectGuaranteedNonNullVarsWhenTrue(
      TreePath condPath, Set<VariableElement> out) {
    while (condPath.getLeaf() instanceof ParenthesizedTree pt) {
      condPath = new TreePath(condPath, pt.getExpression());
    }
    Tree leaf = condPath.getLeaf();
    if (leaf instanceof BinaryTree bt) {
      TreePath lp = new TreePath(condPath, bt.getLeftOperand());
      TreePath rp = new TreePath(condPath, bt.getRightOperand());
      if (bt.getKind() == Tree.Kind.CONDITIONAL_AND) {
        collectGuaranteedNonNullVarsWhenTrue(lp, out);
        collectGuaranteedNonNullVarsWhenTrue(rp, out);
        return;
      }
      if (bt.getKind() == Tree.Kind.CONDITIONAL_OR) {
        collectUnconditionallyDereferencedVars(lp, out);
        return;
      }
      if (bt.getKind() == Tree.Kind.NOT_EQUAL_TO) {
        if (typeOf(rp).getKind() == TypeKind.NULL) {
          addVarIfIdentifier(lp, out);
        } else if (typeOf(lp).getKind() == TypeKind.NULL) {
          addVarIfIdentifier(rp, out);
        }
      }
      collectUnconditionallyDereferencedVars(lp, out);
      collectUnconditionallyDereferencedVars(rp, out);
      return;
    }
    collectUnconditionallyDereferencedVars(condPath, out);
  }

  private void collectUnconditionallyDereferencedVars(
      TreePath exprPath, Set<VariableElement> out) {
    while (exprPath.getLeaf() instanceof ParenthesizedTree pt) {
      exprPath = new TreePath(exprPath, pt.getExpression());
    }
    Tree leaf = exprPath.getLeaf();
    if (leaf instanceof UnaryTree ut && ut.getKind() == Tree.Kind.LOGICAL_COMPLEMENT) {
      collectUnconditionallyDereferencedVars(new TreePath(exprPath, ut.getExpression()), out);
    } else if (leaf instanceof MemberSelectTree ms) {
      Element el = resolver.trees.getElement(exprPath);
      if (el instanceof VariableElement ve
          && ve.getKind() == ElementKind.FIELD
          && !ve.getModifiers().contains(Modifier.STATIC)) {
        TreePath recv = new TreePath(exprPath, ms.getExpression());
        addVarIfIdentifier(recv, out);
        collectUnconditionallyDereferencedVars(recv, out);
      }
    } else if (leaf instanceof MethodInvocationTree mi) {
      Element el = resolver.trees.getElement(exprPath);
      if (el instanceof ExecutableElement ee) {
        if (!ee.getModifiers().contains(Modifier.STATIC)
            && mi.getMethodSelect() instanceof MemberSelectTree ms) {
          TreePath recv =
              new TreePath(new TreePath(exprPath, mi.getMethodSelect()), ms.getExpression());
          addVarIfIdentifier(recv, out);
          collectUnconditionallyDereferencedVars(recv, out);
        }
        for (ExpressionTree a : mi.getArguments()) {
          collectUnconditionallyDereferencedVars(new TreePath(exprPath, a), out);
        }
      }
    } else if (leaf instanceof BinaryTree bt
        && bt.getKind() != Tree.Kind.CONDITIONAL_AND
        && bt.getKind() != Tree.Kind.CONDITIONAL_OR) {
      collectUnconditionallyDereferencedVars(new TreePath(exprPath, bt.getLeftOperand()), out);
      collectUnconditionallyDereferencedVars(new TreePath(exprPath, bt.getRightOperand()), out);
    }
  }

  private void addVarIfIdentifier(TreePath path, Set<VariableElement> out) {
    while (path.getLeaf() instanceof ParenthesizedTree pt) {
      path = new TreePath(path, pt.getExpression());
    }
    if (path.getLeaf() instanceof IdentifierTree) {
      Element el = resolver.trees.getElement(path);
      if (el instanceof VariableElement ve) {
        out.add(ve);
      }
    }
  }

  private static boolean hasNullableAnnotation(Element el, TypeMirror tm) {
    for (var am : el.getAnnotationMirrors()) {
      if (am.getAnnotationType().asElement().getSimpleName().contentEquals("Nullable")) {
        return true;
      }
    }
    if (tm != null) {
      for (var am : tm.getAnnotationMirrors()) {
        if (am.getAnnotationType().asElement().getSimpleName().contentEquals("Nullable")) {
          return true;
        }
      }
    }
    return false;
  }

  private boolean isKnownNonNullDiagnosticRecv(
      TreePath recvPath, Set<VariableElement> nonNullVars) {
    while (recvPath.getLeaf() instanceof ParenthesizedTree pt) {
      recvPath = new TreePath(recvPath, pt.getExpression());
    }
    Tree leaf = recvPath.getLeaf();
    if (leaf.getKind() == Tree.Kind.STRING_LITERAL) {
      return true;
    }
    if (leaf instanceof IdentifierTree id) {
      String nm = id.getName().toString();
      if (nm.equals("this") || nm.equals("super")) {
        return true;
      }
      Element el = resolver.trees.getElement(recvPath);
      if (el instanceof VariableElement ve) {
        if (ve.getKind() == ElementKind.ENUM_CONSTANT) {
          return true;
        }
        if (nonNullVars.contains(ve)) {
          return true;
        }
        return !hasNullableAnnotation(ve, ve.asType());
      }
      return false;
    }
    if (leaf instanceof MemberSelectTree ms) {
      if (ms.getIdentifier().contentEquals("this") || ms.getIdentifier().contentEquals("class")) {
        return true;
      }
      Element el = resolver.trees.getElement(recvPath);
      if (el instanceof VariableElement ve) {
        if (ve.getKind() == ElementKind.ENUM_CONSTANT) {
          return true;
        }
        if (ve.getModifiers().contains(Modifier.STATIC)) {
          return !hasNullableAnnotation(ve, ve.asType());
        }
        if (ve.getKind() == ElementKind.FIELD) {
          return isKnownNonNullDiagnosticRecv(new TreePath(recvPath, ms.getExpression()), nonNullVars)
              && !hasNullableAnnotation(ve, ve.asType());
        }
      }
    }
    return false;
  }

  private boolean isPureDiagnosticArg(TreePath argPath, Set<VariableElement> nonNullVars) {
    while (argPath.getLeaf() instanceof ParenthesizedTree pt) {
      argPath = new TreePath(argPath, pt.getExpression());
    }
    Tree leaf = argPath.getLeaf();
    return switch (leaf.getKind()) {
      case INT_LITERAL,
          LONG_LITERAL,
          FLOAT_LITERAL,
          DOUBLE_LITERAL,
          BOOLEAN_LITERAL,
          CHAR_LITERAL,
          STRING_LITERAL,
          NULL_LITERAL,
          IDENTIFIER ->
          true;
      case MEMBER_SELECT -> {
        MemberSelectTree ms = (MemberSelectTree) leaf;
        String id = ms.getIdentifier().toString();
        if (id.equals("class") || id.equals("this") || id.equals("super")) {
          yield true;
        }
        TreePath sub = new TreePath(argPath, ms.getExpression());
        TypeMirror subTm = typeOf(sub);
        if (id.equals("length") && subTm != null && subTm.getKind() == TypeKind.ARRAY) {
          yield isKnownNonNullDiagnosticRecv(sub, nonNullVars);
        }
        Element el = resolver.trees.getElement(argPath);
        if (el instanceof VariableElement ve) {
          if (ve.getModifiers().contains(Modifier.STATIC)
              || ve.getKind() == ElementKind.ENUM_CONSTANT) {
            yield true;
          }
          if (ve.getKind() == ElementKind.FIELD) {
            yield isKnownNonNullDiagnosticRecv(sub, nonNullVars);
          }
        }
        yield false;
      }
      case TYPE_CAST -> {
        TypeCastTree tc = (TypeCastTree) leaf;
        TreePath sub = new TreePath(argPath, tc.getExpression());
        TypeMirror srcTm = resolver.safeErasure(typeOf(sub));
        TypeMirror dstTm = resolver.safeErasure(typeOf(argPath));
        if (srcTm != null
            && dstTm != null
            && ((srcTm.getKind().isPrimitive() && dstTm.getKind().isPrimitive())
                || (!srcTm.getKind().isPrimitive()
                    && !dstTm.getKind().isPrimitive()
                    && resolver.types.isAssignable(srcTm, dstTm)))) {
          yield isPureDiagnosticArg(sub, nonNullVars);
        }
        yield false;
      }
      case UNARY_PLUS, UNARY_MINUS, BITWISE_COMPLEMENT, LOGICAL_COMPLEMENT -> {
        UnaryTree ut = (UnaryTree) leaf;
        TreePath sub = new TreePath(argPath, ut.getExpression());
        TypeMirror subTm = typeOf(sub);
        yield subTm != null
            && subTm.getKind().isPrimitive()
            && isPureDiagnosticArg(sub, nonNullVars);
      }
      case PLUS,
          MINUS,
          MULTIPLY,
          AND,
          OR,
          XOR,
          EQUAL_TO,
          NOT_EQUAL_TO,
          LESS_THAN,
          LESS_THAN_EQUAL,
          GREATER_THAN,
          GREATER_THAN_EQUAL,
          CONDITIONAL_AND,
          CONDITIONAL_OR -> {
        BinaryTree bt = (BinaryTree) leaf;
        TreePath lp = new TreePath(argPath, bt.getLeftOperand());
        TreePath rp = new TreePath(argPath, bt.getRightOperand());
        if (bt.getKind() == Tree.Kind.PLUS && resolver.isStringType(typeOf(argPath))) {
          yield isPureDiagnosticArg(lp, nonNullVars) && isPureDiagnosticArg(rp, nonNullVars);
        }
        TypeMirror lTm = typeOf(lp);
        TypeMirror rTm = typeOf(rp);
        yield lTm != null
            && rTm != null
            && lTm.getKind().isPrimitive()
            && rTm.getKind().isPrimitive()
            && isPureDiagnosticArg(lp, nonNullVars)
            && isPureDiagnosticArg(rp, nonNullVars);
      }
      case METHOD_INVOCATION -> {
        MethodInvocationTree mi = (MethodInvocationTree) leaf;
        Element el = resolver.trees.getElement(argPath);
        if (!(el instanceof ExecutableElement ee)
            || !(ee.getEnclosingElement() instanceof TypeElement owner)) {
          yield false;
        }
        String mName = ee.getSimpleName().toString();
        String ownerQn = owner.getQualifiedName().toString();
        if (ee.getModifiers().contains(Modifier.STATIC)) {
          if (mi.getArguments().size() == 1
              && ee.getParameters().size() == 1
              && ((mName.equals("valueOf")
                      && (resolver.runtime.isBoxedPrimitive(owner)
                          || ownerQn.equals("java.lang.String"))
                      && (ee.getParameters().get(0).asType().getKind().isPrimitive()
                          || ownerQn.equals("java.lang.String")))
                  || (ownerQn.equals("java.util.Objects") && mName.equals("toString")))) {
            TreePath a0 = new TreePath(argPath, mi.getArguments().get(0));
            if (ee.getParameters().get(0).asType().getKind().isPrimitive()) {
              yield typeOf(a0).getKind().isPrimitive() && isPureDiagnosticArg(a0, nonNullVars);
            }
            yield isPureDiagnosticArg(a0, nonNullVars);
          }
          yield false;
        }
        if (!mi.getArguments().isEmpty()) {
          yield false;
        }
        boolean pureGetterName =
            ((mName.startsWith("get") && mName.length() > 3 && !mName.startsWith("getAnd"))
                || (mName.startsWith("is") && mName.length() > 2)
                || (mName.startsWith("has") && mName.length() > 3)
                || mName.equals("size")
                || mName.equals("length")
                || mName.equals("name")
                || mName.equals("ordinal")
                || mName.equals("optMessage"));
        if (!pureGetterName) {
          yield false;
        }
        if (mi.getMethodSelect() instanceof IdentifierTree) {
          yield true;
        }
        if (mi.getMethodSelect() instanceof MemberSelectTree ms) {
          TreePath recv = new TreePath(new TreePath(argPath, ms), ms.getExpression());
          yield isKnownNonNullDiagnosticRecv(recv, nonNullVars);
        }
        yield false;
      }
      default -> false;
    };
  }

  private String emitPreconditionsIntrinsic(
      MethodInvocationTree mi, TreePath path, ExecutableElement ee, EmitContext ctx) {
    String m = ee.getSimpleName().toString();
    List<? extends ExpressionTree> args = mi.getArguments();
    int arity = args.size();
    String rt = rtPrefix(ctx);
    TypeMirror boolTm = resolver.types.getPrimitiveType(TypeKind.BOOLEAN);
    TypeMirror intTm = resolver.types.getPrimitiveType(TypeKind.INT);

    if (m.equals("checkState") || m.equals("checkArgument")) {
      if (arity >= 2) {
        String lazyStmt = tryEmitPreconditionsStmt(path, ctx);
        if (lazyStmt != null) {
          return "func() { " + lazyStmt + " }()";
        }
      }
      String fn = m.equals("checkState") ? "CheckState" : "CheckArgument";
      String cond = emitAs(args.get(0), path, boolTm, false, ctx);
      if (arity == 1) {
        return rt + fn + "(" + cond + ")";
      }
      if (arity == 2) {
        return rt
            + fn
            + "Msg("
            + cond
            + ", "
            + emitAs(args.get(1), path, resolver.getObjectElement().asType(), true, ctx)
            + ")";
      }
      List<String> fmtArgs = new ArrayList<>();
      for (int i = 2; i < arity; i++) {
        fmtArgs.add(emitAs(args.get(i), path, resolver.getObjectElement().asType(), true, ctx));
      }
      return rt
          + fn
          + "Fmt("
          + cond
          + ", "
          + rt
          + "StrVal("
          + emit(args.get(1), path, ctx)
          + "), "
          + String.join(", ", fmtArgs)
          + ")";
    }

    if (m.equals("checkNotNull") && arity >= 1) {
      TypeMirror retTm = typeOf(path);
      String val = emitAs(args.get(0), path, retTm, false, ctx);
      return emitCheckNotNull(args.get(0), val, retTm, ctx);
    }

    if (m.equals("checkElementIndex") && arity >= 2) {
      String idx = emitAs(args.get(0), path, intTm, false, ctx);
      String sz = emitAs(args.get(1), path, intTm, false, ctx);
      return rt + "CheckElementIndex(" + idx + ", " + sz + ")";
    }

    if (m.equals("checkPositionIndex") && arity >= 2) {
      String idx = emitAs(args.get(0), path, intTm, false, ctx);
      String sz = emitAs(args.get(1), path, intTm, false, ctx);
      return rt + "CheckPositionIndex(" + idx + ", " + sz + ")";
    }

    if (m.equals("checkPositionIndexes") && arity == 3) {
      return rt
          + "CheckPositionIndexes("
          + emitAs(args.get(0), path, intTm, false, ctx)
          + ", "
          + emitAs(args.get(1), path, intTm, false, ctx)
          + ", "
          + emitAs(args.get(2), path, intTm, false, ctx)
          + ")";
    }
    return null;
  }

  private String emitMathIntrinsic(
      MethodInvocationTree mi, TreePath path, ExecutableElement ee, EmitContext ctx) {
    String m = ee.getSimpleName().toString();
    List<? extends ExpressionTree> args = mi.getArguments();
    int arity = args.size();
    TypeMirror retTm = typeOf(path);
    TypeMirror dblTm = resolver.types.getPrimitiveType(TypeKind.DOUBLE);

    if ((m.equals("min") || m.equals("max")) && arity == 2) {
      String a = emitAs(args.get(0), path, retTm, false, ctx);
      String b = emitAs(args.get(1), path, retTm, false, ctx);
      if (retTm.getKind() == TypeKind.DOUBLE || retTm.getKind() == TypeKind.FLOAT) {
        ctx.imports.std("math");
        String call =
            "math."
                + (m.equals("min") ? "Min" : "Max")
                + "(float64("
                + a
                + "), float64("
                + b
                + "))";
        return retTm.getKind() == TypeKind.FLOAT ? ("float32(" + call + ")") : call;
      }
      return m + "(" + a + ", " + b + ")";
    }
    if (m.equals("abs") && arity == 1) {
      String a = emitAs(args.get(0), path, retTm, false, ctx);
      if (retTm.getKind() == TypeKind.DOUBLE || retTm.getKind() == TypeKind.FLOAT) {
        ctx.imports.std("math");
        String call = "math.Abs(float64(" + a + "))";
        return retTm.getKind() == TypeKind.FLOAT ? ("float32(" + call + ")") : call;
      }
      String tmp = ctx.nextTemp("ab");
      String goT = resolver.goType(retTm, ctx.imports);
      return "func() " + goT + " { " + tmp + " := " + a + "; if " + tmp + " < 0 { return -" + tmp + " }; return " + tmp + " }()";
    }
    if (arity == 1
        && Set.of("floor", "ceil", "sqrt", "sin", "cos", "tan", "log", "log10", "exp", "rint")
            .contains(m)) {
      ctx.imports.std("math");
      String goFn =
          switch (m) {
            case "floor" -> "Floor";
            case "ceil" -> "Ceil";
            case "sqrt" -> "Sqrt";
            case "sin" -> "Sin";
            case "cos" -> "Cos";
            case "tan" -> "Tan";
            case "log" -> "Log";
            case "log10" -> "Log10";
            case "exp" -> "Exp";
            case "rint" -> "RoundToEven";
            default -> "Floor";
          };
      return "math." + goFn + "(" + emitAs(args.get(0), path, dblTm, false, ctx) + ")";
    }
    if (m.equals("pow") && arity == 2) {
      ctx.imports.std("math");
      return "math.Pow("
          + emitAs(args.get(0), path, dblTm, false, ctx)
          + ", "
          + emitAs(args.get(1), path, dblTm, false, ctx)
          + ")";
    }
    if (m.equals("atan2") && arity == 2) {
      ctx.imports.std("math");
      return "math.Atan2("
          + emitAs(args.get(0), path, dblTm, false, ctx)
          + ", "
          + emitAs(args.get(1), path, dblTm, false, ctx)
          + ")";
    }
    return null;
  }

  private String emitWrapperIntrinsic(
      MethodInvocationTree mi,
      TreePath path,
      ExecutableElement ee,
      String qn,
      String m,
      boolean isStatic,
      TreePath recvPath,
      EmitContext ctx) {
    List<? extends ExpressionTree> args = mi.getArguments();
    int arity = args.size();
    String rt = rtPrefix(ctx);
    String jl = jlPrefix(ctx);
    TypeMirror intTm = resolver.types.getPrimitiveType(TypeKind.INT);
    TypeMirror longTm = resolver.types.getPrimitiveType(TypeKind.LONG);
    TypeMirror dblTm = resolver.types.getPrimitiveType(TypeKind.DOUBLE);

    if (!isStatic && resolver.runtime.BOXED_TYPES_CONTAINS(qn)) {
      String recv = recvPath != null ? emitPath(recvPath, ctx) : ctx.selfValueExpr();
      return switch (m) {
        case "booleanValue" -> recv + ".BooleanValue()";
        case "byteValue" -> recv + ".ByteValue()";
        case "shortValue" -> recv + ".ShortValue()";
        case "charValue" -> recv + ".CharValue()";
        case "intValue" -> recv + ".IntValue()";
        case "longValue" -> recv + ".LongValue()";
        case "floatValue" -> recv + ".FloatValue()";
        case "doubleValue" -> recv + ".DoubleValue()";
        case "equals" ->
            arity == 1
                ? (recv + ".Equals(" + emitAs(args.get(0), path, resolver.getObjectElement().asType(), true, ctx) + ")")
                : null;
        case "hashCode" -> arity == 0 ? (recv + ".HashCode()") : null;
        case "toString" -> arity == 0 ? (recv + ".String()") : null;
        default -> null;
      };
    }

    if (qn.equals("java.lang.Number") && !isStatic) {
      String recv = recvPath != null ? emitPath(recvPath, ctx) : ctx.selfValueExpr();
      return switch (m) {
        case "byteValue" -> recv + ".ByteValue()";
        case "shortValue" -> recv + ".ShortValue()";
        case "intValue" -> recv + ".IntValue()";
        case "longValue" -> recv + ".LongValue()";
        case "floatValue" -> recv + ".FloatValue()";
        case "doubleValue" -> recv + ".DoubleValue()";
        default -> null;
      };
    }

    if (!isStatic) {
      return null;
    }

    if (qn.equals("java.lang.Integer")) {
      if (m.equals("valueOf") && arity == 1) {
        TypeMirror p0 = resolver.safeErasure(ee.getParameters().get(0).asType());
        if (p0.getKind() == TypeKind.INT) {
          return rt + "ValueOfInt(" + emitAs(args.get(0), path, intTm, false, ctx) + ")";
        }
        return rt + "ValueOfInt(" + rt + "ParseInt(" + emit(args.get(0), path, ctx) + ", 10))";
      }
      if (m.equals("parseInt") && arity == 1) {
        return rt + "ParseInt(" + emit(args.get(0), path, ctx) + ", 10)";
      }
      if (m.equals("parseInt") && arity == 2) {
        return rt + "ParseInt(" + emit(args.get(0), path, ctx) + ", " + emitAs(args.get(1), path, intTm, false, ctx) + ")";
      }
      if (m.equals("toHexString") && arity == 1) {
        return jl + "IntegerToHexString(" + emitAs(args.get(0), path, intTm, false, ctx) + ")";
      }
      if (m.equals("toString") && arity == 1) {
        ctx.imports.std("strconv");
        return "strconv.FormatInt(int64(" + emitAs(args.get(0), path, intTm, false, ctx) + "), 10)";
      }
      if (m.equals("bitCount") && arity == 1) {
        return rt + "IntBitCount(" + emitAs(args.get(0), path, intTm, false, ctx) + ")";
      }
      if (m.equals("numberOfLeadingZeros") && arity == 1) {
        return rt + "IntLeadingZeros(" + emitAs(args.get(0), path, intTm, false, ctx) + ")";
      }
      if (m.equals("numberOfTrailingZeros") && arity == 1) {
        return rt + "IntTrailingZeros(" + emitAs(args.get(0), path, intTm, false, ctx) + ")";
      }
      if (m.equals("lowestOneBit") && arity == 1) {
        return rt + "IntLowestOneBit(" + emitAs(args.get(0), path, intTm, false, ctx) + ")";
      }
      if (m.equals("highestOneBit") && arity == 1) {
        return rt + "IntHighestOneBit(" + emitAs(args.get(0), path, intTm, false, ctx) + ")";
      }
      if (m.equals("compare") && arity == 2) {
        return rt
            + "CompareValues("
            + rt
            + "ValueOfInt("
            + emitAs(args.get(0), path, intTm, false, ctx)
            + "), "
            + rt
            + "ValueOfInt("
            + emitAs(args.get(1), path, intTm, false, ctx)
            + "))";
      }
      if (m.equals("hashCode") && arity == 1) {
        return emitAs(args.get(0), path, intTm, false, ctx);
      }
    }

    if (qn.equals("java.lang.Long")) {
      if (m.equals("valueOf") && arity == 1 && ee.getParameters().get(0).asType().getKind() == TypeKind.LONG) {
        return rt + "ValueOfLong(" + emitAs(args.get(0), path, longTm, false, ctx) + ")";
      }
      if (m.equals("parseLong") && arity == 1) {
        return rt + "ParseLong(" + emit(args.get(0), path, ctx) + ", 10)";
      }
      if (m.equals("parseLong") && arity == 2) {
        return rt + "ParseLong(" + emit(args.get(0), path, ctx) + ", " + emitAs(args.get(1), path, intTm, false, ctx) + ")";
      }
      if (m.equals("toHexString") && arity == 1) {
        return jl + "LongToHexString(" + emitAs(args.get(0), path, longTm, false, ctx) + ")";
      }
      if (m.equals("toString") && arity == 1) {
        ctx.imports.std("strconv");
        return "strconv.FormatInt(" + emitAs(args.get(0), path, longTm, false, ctx) + ", 10)";
      }
      if (m.equals("hashCode") && arity == 1) {
        return rt + "ValueOfLong(" + emitAs(args.get(0), path, longTm, false, ctx) + ").HashCode()";
      }
    }

    if (qn.equals("java.lang.Float")) {
      TypeMirror fltTm = resolver.types.getPrimitiveType(TypeKind.FLOAT);
      if (m.equals("valueOf") && arity == 1 && ee.getParameters().get(0).asType().getKind() == TypeKind.FLOAT) {
        return rt + "ValueOfFloat(" + emitAs(args.get(0), path, fltTm, false, ctx) + ")";
      }
      if (m.equals("isNaN") && arity == 1) {
        ctx.imports.std("math");
        return "math.IsNaN(float64(" + emitAs(args.get(0), path, fltTm, false, ctx) + "))";
      }
      if (m.equals("isInfinite") && arity == 1) {
        ctx.imports.std("math");
        return "math.IsInf(float64(" + emitAs(args.get(0), path, fltTm, false, ctx) + "), 0)";
      }
      if (m.equals("isFinite") && arity == 1) {
        ctx.imports.std("math");
        String tmp = ctx.nextTemp("ff");
        return "func() bool { "
            + tmp
            + " := float64("
            + emitAs(args.get(0), path, fltTm, false, ctx)
            + "); return !math.IsNaN("
            + tmp
            + ") && !math.IsInf("
            + tmp
            + ", 0) }()";
      }
      if (m.equals("toString") && arity == 1) {
        return jl + "FloatToString(" + emitAs(args.get(0), path, fltTm, false, ctx) + ")";
      }
      if (m.equals("hashCode") && arity == 1) {
        return rt + "ValueOfFloat(" + emitAs(args.get(0), path, fltTm, false, ctx) + ").HashCode()";
      }
    }

    if (qn.equals("java.lang.Double")) {
      if (m.equals("valueOf") && arity == 1 && ee.getParameters().get(0).asType().getKind() == TypeKind.DOUBLE) {
        return rt + "ValueOfDouble(" + emitAs(args.get(0), path, dblTm, false, ctx) + ")";
      }
      if (m.equals("parseDouble") && arity == 1) {
        return rt + "ParseDouble(" + emit(args.get(0), path, ctx) + ")";
      }
      if (m.equals("isNaN") && arity == 1) {
        ctx.imports.std("math");
        return "math.IsNaN(" + emitAs(args.get(0), path, dblTm, false, ctx) + ")";
      }
      if (m.equals("isInfinite") && arity == 1) {
        ctx.imports.std("math");
        return "math.IsInf(" + emitAs(args.get(0), path, dblTm, false, ctx) + ", 0)";
      }
      if (m.equals("isFinite") && arity == 1) {
        ctx.imports.std("math");
        String tmp = ctx.nextTemp("df");
        return "func() bool { " + tmp + " := " + emitAs(args.get(0), path, dblTm, false, ctx) + "; return !math.IsNaN(" + tmp + ") && !math.IsInf(" + tmp + ", 0) }()";
      }
      if ((m.equals("doubleToLongBits") || m.equals("doubleToRawLongBits")) && arity == 1) {
        ctx.imports.std("math");
        return "int64(math.Float64bits(" + emitAs(args.get(0), path, dblTm, false, ctx) + "))";
      }
      if (m.equals("longBitsToDouble") && arity == 1) {
        ctx.imports.std("math");
        return "math.Float64frombits(uint64(" + emitAs(args.get(0), path, longTm, false, ctx) + "))";
      }
      if (m.equals("toString") && arity == 1) {
        return jl + "DoubleToString(" + emitAs(args.get(0), path, dblTm, false, ctx) + ")";
      }
      if (m.equals("hashCode") && arity == 1) {
        return rt + "ValueOfDouble(" + emitAs(args.get(0), path, dblTm, false, ctx) + ").HashCode()";
      }
    }

    if (qn.equals("java.lang.Boolean")) {
      TypeMirror boolTm = resolver.types.getPrimitiveType(TypeKind.BOOLEAN);
      if (m.equals("valueOf") && arity == 1 && ee.getParameters().get(0).asType().getKind() == TypeKind.BOOLEAN) {
        return rt + "ValueOfBoolean(" + emitAs(args.get(0), path, boolTm, false, ctx) + ")";
      }
      if (m.equals("toString") && arity == 1 && ee.getParameters().get(0).asType().getKind() == TypeKind.BOOLEAN) {
        return rt + "BoolToStr(" + emitAs(args.get(0), path, boolTm, false, ctx) + ")";
      }
      if (m.equals("hashCode") && arity == 1) {
        return rt + "ValueOfBoolean(" + emitAs(args.get(0), path, boolTm, false, ctx) + ").HashCode()";
      }
    }

    if (qn.equals("java.lang.Character")) {
      TypeMirror charTm = resolver.types.getPrimitiveType(TypeKind.CHAR);
      if (m.equals("valueOf") && arity == 1) {
        return rt + "ValueOfChar(" + emitAs(args.get(0), path, charTm, false, ctx) + ")";
      }
      if (m.equals("toString") && arity == 1 && ee.getParameters().get(0).asType().getKind() == TypeKind.CHAR) {
        return jl + "StringFromChar(" + emitAs(args.get(0), path, charTm, false, ctx) + ")";
      }
      if (m.equals("charCount") && arity == 1) {
        String cp = emitAs(args.get(0), path, intTm, false, ctx);
        String t = ctx.nextTemp("cp");
        return "func() int32 { " + t + " := " + cp + "; if " + t + " >= 0x10000 { return 2 }; return 1 }()";
      }
    }
    return null;
  }

  // --- Lambdas & Method References ---

  private String emitLambda(LambdaExpressionTree lam, TreePath path, EmitContext ctx) {
    TypeMirror fiTm = typeOf(path);
    TypeMirror erFi = resolver.safeErasure(fiTm);
    if (!(erFi instanceof DeclaredType fiDt) || !(fiDt.asElement() instanceof TypeElement fiTe)) {
      return "nil";
    }
    ExecutableElement sam = resolver.findSamMethod(fiTe);
    if (sam == null) {
      return "nil";
    }
    TypeResolver.MethodSig samSig = resolver.getInstanceMethodSig(sam);
    ExecutableType instSam = resolver.resolveExecutableType(fiTm, sam);

    // Build lambda function signature matching `NewLambda_<FI>`
    boolean stateless = isStatelessLambda(lam, path);
    EmitContext lamCtx = ctx.copyForLambda();
    if (!stateless) {
      lamCtx.capturedVarFields.putAll(ctx.capturedVarFields);
      // Copy outer local variable mappings so lambda body sees enclosing locals directly via Go closure!
      copyLocalsIntoLambda(ctx, lamCtx);
    }

    TypeMirror nodePredArgTm =
        lam.getParameters().size() == 1 ? extractNodePredicateTypeArg(fiTm, fiTe, samSig) : null;

    StringBuilder fnHeader = new StringBuilder("func(");
    StringBuilder unboxPreamble = new StringBuilder();
    if (nodePredArgTm != null) {
      String rawArg = lamCtx.nextTemp("rawNode");
      fnHeader.append(rawArg).append(" uint32) bool");
      TreePath pPath = new TreePath(path, lam.getParameters().get(0));
      VariableElement pVe = (VariableElement) resolver.trees.getElement(pPath);
      TypeMirror actualParamTm = pVe.asType();
      String actualParamGo = resolver.goType(actualParamTm, ctx.imports);
      String nodeGoType = resolver.goType(nodePredArgTm, ctx.imports);
      String userVar = lamCtx.registerLocalVar(pVe);
      String nodeExpr = nodeGoType + "(" + rawArg + ")";
      String bindExpr =
          actualParamGo.equals(nodeGoType)
              ? nodeExpr
              : convert(nodeExpr, nodePredArgTm, actualParamTm, false, lamCtx);
      unboxPreamble.append('\t').append(userVar).append(" := ").append(bindExpr).append('\n');
      unboxPreamble.append("\t_ = ").append(userVar).append('\n');
    } else {
      for (int i = 0; i < lam.getParameters().size(); i++) {
        if (i > 0) fnHeader.append(", ");
        TreePath pPath = new TreePath(path, lam.getParameters().get(i));
        VariableElement pVe = (VariableElement) resolver.trees.getElement(pPath);
        TypeMirror canonParamTm = resolver.resolveSigParamType(samSig, i, fiTm);
        String canonParamGo = resolver.goType(canonParamTm, ctx.imports);
        TypeMirror actualParamTm = pVe.asType();
        String actualParamGo = resolver.goType(actualParamTm, ctx.imports);

        String userVar = lamCtx.registerLocalVar(pVe);
        if (canonParamGo.equals(actualParamGo) && !samSig.paramWidenedToAny[i]) {
          fnHeader.append(userVar).append(' ').append(canonParamGo);
          unboxPreamble.append("\t_ = ").append(userVar).append('\n');
        } else {
          String rawArg = lamCtx.nextTemp("lp");
          fnHeader.append(rawArg).append(' ').append(canonParamGo);
          String conv = convert(rawArg, canonParamTm, actualParamTm, false, lamCtx);
          unboxPreamble.append('\t').append(userVar).append(" := ").append(conv).append('\n');
          unboxPreamble.append("\t_ = ").append(userVar).append('\n');
        }
      }
      fnHeader.append(')');
      TypeMirror samRetTm = resolver.resolveSigReturnType(samSig, fiTm);
      String retGo = samRetTm != null ? resolver.goType(samRetTm, ctx.imports) : "";
      if (!retGo.isEmpty()) {
        fnHeader.append(' ').append(retGo);
      }
    }

    TypeMirror resolvedSamRet = resolver.resolveSigReturnType(samSig, fiTm);
    lamCtx.expectedReturnType = resolvedSamRet;
    lamCtx.returnWidenedToAny = samSig.returnWidenedToAny;

    StringBuilder body = new StringBuilder();
    body.append(unboxPreamble);
    if (lam.getBodyKind() == LambdaExpressionTree.BodyKind.EXPRESSION) {
      ExpressionTree exprBody = (ExpressionTree) lam.getBody();
      if (resolvedSamRet == null) {
        body.append('\t')
            .append(stmtEmitter.emitExprAsStmt(exprBody, path, lamCtx))
            .append('\n');
      } else {
        String val = emitAs(exprBody, path, resolvedSamRet, samSig.returnWidenedToAny, lamCtx);
        body.append("\treturn ").append(val).append('\n');
      }
    } else {
      com.sun.source.tree.BlockTree bt = (com.sun.source.tree.BlockTree) lam.getBody();
      body.append(stmtEmitter.emitBlockStatements(bt, new TreePath(path, bt), 1, lamCtx));
      if (resolvedSamRet != null && !MethodEmitter.endsWithGoTerminatingReturnOrThrow(bt)) {
        body.append("\tpanic(\"unreachable\")\n");
      }
    }

    ctx.syncTempSeqFrom(lamCtx);
    ctx.deferredLocalTypes.addAll(lamCtx.deferredLocalTypes);

    String ctorPrefix =
        nodePredArgTm != null ? nodeLambdaCtorPrefix(fiTe, ctx) : lambdaCtorPrefix(fiTe, fiTm, ctx);
    String lamExpr = ctorPrefix + "(" + fnHeader + " {\n" + body + "})";
    if (stateless
        && !resolver.hasActiveTypeParameters()
        && !containsGoIdent(lamExpr, ctx.receiverName)) {
      String fiGo = resolver.goType(fiTm, ctx.imports);
      return ctx.hoistLambdaSingleton(fiGo, lamExpr);
    }
    return lamExpr;
  }

  private boolean isStatelessLambda(LambdaExpressionTree lam, TreePath lamPath) {
    Set<VariableElement> declaredInside = new HashSet<>();
    boolean[] capturesState = {false};

    new TreePathScanner<Void, Void>() {
      @Override
      public Void scan(Tree tree, Void p) {
        if (capturesState[0] || tree == null) {
          return null;
        }
        return super.scan(tree, p);
      }

      @Override
      public Void visitClass(ClassTree node, Void p) {
        capturesState[0] = true;
        return null;
      }

      @Override
      public Void visitVariable(VariableTree node, Void p) {
        Element el = resolver.trees.getElement(getCurrentPath());
        if (el instanceof VariableElement ve) {
          declaredInside.add(ve);
        }
        return super.visitVariable(node, p);
      }

      @Override
      public Void visitBindingPattern(BindingPatternTree node, Void p) {
        TreePath vp = new TreePath(getCurrentPath(), node.getVariable());
        Element el = resolver.trees.getElement(vp);
        if (el instanceof VariableElement ve) {
          declaredInside.add(ve);
        }
        return super.visitBindingPattern(node, p);
      }

      @Override
      public Void visitIdentifier(IdentifierTree node, Void p) {
        String name = node.getName().toString();
        if (name.equals("this") || name.equals("super")) {
          capturesState[0] = true;
          return null;
        }
        Element el = resolver.trees.getElement(getCurrentPath());
        if (el instanceof VariableElement ve) {
          switch (ve.getKind()) {
            case LOCAL_VARIABLE,
                PARAMETER,
                EXCEPTION_PARAMETER,
                RESOURCE_VARIABLE,
                BINDING_VARIABLE -> {
              if (!declaredInside.contains(ve)) {
                capturesState[0] = true;
              }
            }
            case FIELD -> {
              if (!ve.getModifiers().contains(Modifier.STATIC)) {
                capturesState[0] = true;
              }
            }
            default -> {}
          }
        }
        return null;
      }

      @Override
      public Void visitMemberSelect(MemberSelectTree node, Void p) {
        String member = node.getIdentifier().toString();
        if (member.equals("this") || member.equals("super")) {
          capturesState[0] = true;
          return null;
        }
        return super.visitMemberSelect(node, p);
      }

      @Override
      public Void visitMethodInvocation(MethodInvocationTree node, Void p) {
        if (node.getMethodSelect() instanceof IdentifierTree) {
          Element el = resolver.trees.getElement(getCurrentPath());
          if (el instanceof ExecutableElement ee && !ee.getModifiers().contains(Modifier.STATIC)) {
            capturesState[0] = true;
            return null;
          }
        }
        for (Tree ta : node.getTypeArguments()) {
          scan(ta, p);
        }
        if (!(node.getMethodSelect() instanceof IdentifierTree)) {
          scan(node.getMethodSelect(), p);
        }
        for (ExpressionTree arg : node.getArguments()) {
          scan(arg, p);
        }
        return null;
      }

      @Override
      public Void visitNewClass(NewClassTree node, Void p) {
        if (node.getClassBody() != null) {
          capturesState[0] = true;
          return null;
        }
        Element el = resolver.trees.getElement(getCurrentPath());
        if (el instanceof ExecutableElement ctor
            && ctor.getEnclosingElement() instanceof TypeElement owner
            && resolver.isSourceType(owner)) {
          if ((node.getEnclosingExpression() == null && resolver.hasOuterInstance(owner))
              || !resolver.getCapturedVars(owner).isEmpty()) {
            capturesState[0] = true;
            return null;
          }
        }
        return super.visitNewClass(node, p);
      }

      @Override
      public Void visitMemberReference(MemberReferenceTree node, Void p) {
        if (node.getMode() == MemberReferenceTree.ReferenceMode.NEW) {
          Element el = resolver.trees.getElement(getCurrentPath());
          if (el instanceof ExecutableElement ctor
              && ctor.getEnclosingElement() instanceof TypeElement owner
              && resolver.isSourceType(owner)
              && (resolver.hasOuterInstance(owner) || !resolver.getCapturedVars(owner).isEmpty())) {
            capturesState[0] = true;
            return null;
          }
        }
        return super.visitMemberReference(node, p);
      }
    }.scan(lamPath, null);

    return !capturesState[0];
  }

  private static boolean containsGoIdent(String code, String ident) {
    if (ident == null || ident.isEmpty()) {
      return false;
    }
    int idx = 0;
    int len = ident.length();
    while ((idx = code.indexOf(ident, idx)) >= 0) {
      boolean leftOk =
          idx == 0
              || (!Character.isLetterOrDigit(code.charAt(idx - 1)) && code.charAt(idx - 1) != '_');
      int end = idx + len;
      boolean rightOk =
          end == code.length()
              || (!Character.isLetterOrDigit(code.charAt(end)) && code.charAt(end) != '_');
      if (leftOk && rightOk) {
        return true;
      }
      idx = end;
    }
    return false;
  }

  private String emitMemberReference(MemberReferenceTree mr, TreePath path, EmitContext ctx) {
    TypeMirror fiTm = typeOf(path);
    TypeMirror erFi = resolver.safeErasure(fiTm);
    if (!(erFi instanceof DeclaredType fiDt) || !(fiDt.asElement() instanceof TypeElement fiTe)) {
      return "nil";
    }
    ExecutableElement sam = resolver.findSamMethod(fiTe);
    if (sam == null) {
      return "nil";
    }
    ExecutableType instSam = resolver.resolveExecutableType(fiTm, sam);
    TypeResolver.MethodSig samSig = resolver.getInstanceMethodSig(sam);
    Element refEl = resolver.trees.getElement(path);
    TreePath qualPath = new TreePath(path, mr.getQualifierExpression());
    Element qualEl = resolver.trees.getElement(qualPath);
    TypeMirror qualTm = typeOf(qualPath);

    int samArity = samSig.canonParamTypes.size();
    TypeMirror nodePredArgTm =
        samArity == 1 ? extractNodePredicateTypeArg(fiTm, fiTe, samSig) : null;
    StringBuilder fnSig = new StringBuilder("func(");
    List<String> rawParamNames = new ArrayList<>(samArity);
    List<TypeMirror> rawParamTypes = new ArrayList<>(samArity);
    if (nodePredArgTm != null) {
      fnSig.append("rawNode_ uint32) bool");
      String nodeGoType = resolver.goType(nodePredArgTm, ctx.imports);
      rawParamNames.add(nodeGoType + "(rawNode_)");
      rawParamTypes.add(nodePredArgTm);
    } else {
      for (int i = 0; i < samArity; i++) {
        if (i > 0) fnSig.append(", ");
        String pn = "mrp" + i + "_";
        TypeMirror pTm = resolver.resolveSigParamType(samSig, i, fiTm);
        rawParamNames.add(pn);
        rawParamTypes.add(pTm);
        fnSig.append(pn).append(' ').append(resolver.goType(pTm, ctx.imports));
      }
      fnSig.append(')');
      TypeMirror mrSamRetTm = resolver.resolveSigReturnType(samSig, fiTm);
      String retGo = mrSamRetTm != null ? resolver.goType(mrSamRetTm, ctx.imports) : "";
      if (!retGo.isEmpty()) {
        fnSig.append(' ').append(retGo);
      }
    }

    TypeMirror resolvedMrSamRet = resolver.resolveSigReturnType(samSig, fiTm);
    String ctorPrefix =
        nodePredArgTm != null ? nodeLambdaCtorPrefix(fiTe, ctx) : lambdaCtorPrefix(fiTe, fiTm, ctx);
    String fiGo = resolver.goType(fiTm, ctx.imports);

    // Case 1: Array constructor ref `T[]::new`
    if (mr.getMode() == MemberReferenceTree.ReferenceMode.NEW
        && qualTm != null
        && qualTm.getKind() == TypeKind.ARRAY) {
      ArrayType at = (ArrayType) resolver.goEffectiveType(qualTm);
      String sliceGo = resolver.goType(at, ctx.imports);
      String sz =
          convert(
              rawParamNames.get(0),
              rawParamTypes.get(0),
              resolver.types.getPrimitiveType(TypeKind.INT),
              false,
              ctx);
      String alloc =
          resolver.isStringType(at.getComponentType())
              ? (rtPrefix(ctx) + "NewStringArray(" + sz + ")")
              : ("make(" + sliceGo + ", " + sz + ")");
      String convRet =
          convert(alloc, at, resolvedMrSamRet, samSig.returnWidenedToAny, ctx);
      String mrArrExpr = ctorPrefix + "(" + fnSig + " { return " + convRet + " })";
      return !resolver.hasActiveTypeParameters()
          ? ctx.hoistLambdaSingleton(fiGo, mrArrExpr)
          : mrArrExpr;
    }

    if (!(refEl instanceof ExecutableElement targetEe)) {
      String z =
          resolvedMrSamRet != null ? resolver.zeroValue(resolvedMrSamRet, ctx.imports) : "";
      return ctorPrefix + "(" + fnSig + " { " + (z.isEmpty() ? "" : ("return " + z)) + " })";
    }

    TypeElement targetOwner = (TypeElement) targetEe.getEnclosingElement();

    // Case 2: Class constructor ref `Foo::new`
    if (mr.getMode() == MemberReferenceTree.ReferenceMode.NEW) {
      boolean hasOuter = resolver.isSourceType(targetOwner) && resolver.hasOuterInstance(targetOwner);
      boolean hasCaptured =
          resolver.isSourceType(targetOwner) && !resolver.getCapturedVars(targetOwner).isEmpty();
      List<String> cArgs = new ArrayList<>();
      if (hasOuter) {
        TypeElement outerTe = resolver.getEnclosingInstanceType(targetOwner);
        cArgs.add(resolveImplicitThisValue(outerTe, ctx));
      }
      if (hasCaptured) {
        for (VariableElement cap : resolver.getCapturedVars(targetOwner)) {
          String loc = ctx.lookupLocalVar(cap);
          cArgs.add(loc != null ? loc : GoNames.localIdent(cap.getSimpleName().toString()));
        }
      }
      TypeMirror mrConstructedTm =
          resolver.hasEffectiveTypeParameters(targetOwner)
              ? instSam.getReturnType()
              : targetOwner.asType();
      appendMrArgs(rawParamNames, rawParamTypes, 0, targetEe, null, cArgs, ctx, mrConstructedTm, null);
      String fn = resolver.getConstructorGoName(targetEe);
      if (resolver.hasEffectiveTypeParameters(targetOwner)) {
        fn += resolver.formatDeclaredTypeArgs(mrConstructedTm, targetOwner, ctx.imports);
      }
      String call;
      if (resolver.runtime.isShimException(targetOwner)) {
        List<String> exArgs = new ArrayList<>();
        exArgs.add(GoNames.goStringLiteral(targetOwner.getQualifiedName().toString()));
        for (int i = 0; i < samArity; i++) {
          exArgs.add(
              convert(
                  rawParamNames.get(i),
                  rawParamTypes.get(i),
                  resolver.getObjectElement().asType(),
                  true,
                  ctx));
        }
        String exIface = resolver.goTypeForElement(targetOwner, ctx.imports);
        call =
            rtPrefix(ctx)
                + "AsIface["
                + exIface
                + "]("
                + rtPrefix(ctx)
                + "ConstructShimThrowable("
                + String.join(", ", exArgs)
                + "))";
      } else if (resolver.isSourceType(targetOwner)) {
        String pkg = resolver.planner.getGoPkg(targetOwner);
        String alias = ctx.imports.ext(resolver.planner.getGoImportPath(pkg), pkg);
        call = (alias.isEmpty() ? "" : alias + ".") + fn + "(" + String.join(", ", cArgs) + ")";
      } else {
        call = rtPrefix(ctx) + fn + "(" + String.join(", ", cArgs) + ")";
      }
      String retStmt =
          resolvedMrSamRet != null
              ? ("return "
                  + convert(
                      call,
                      mrConstructedTm,
                      resolvedMrSamRet,
                      samSig.returnWidenedToAny,
                      ctx))
              : ("_ = " + call);
      String mrExpr = ctorPrefix + "(" + fnSig + " { " + retStmt + " })";
      if (!hasOuter && !hasCaptured && !resolver.hasActiveTypeParameters()) {
        return ctx.hoistLambdaSingleton(fiGo, mrExpr);
      }
      return mrExpr;
    }

    // Case 3: Static method ref `Foo::staticBar`
    if (targetEe.getModifiers().contains(Modifier.STATIC)) {
      List<TypeMirror> mrStaticTas = null;
      if (resolver.isGenericStaticMethod(targetEe)) {
        List<TypeMirror> explicitTas = new ArrayList<>();
        if (mr.getTypeArguments() != null) {
          for (Tree ta : mr.getTypeArguments()) {
            explicitTas.add(typeOf(new TreePath(path, ta)));
          }
        }
        mrStaticTas =
            resolver.inferGenericStaticMethodTypeArgs(
                targetEe, explicitTas, instSam.getParameterTypes(), instSam.getReturnType());
      }
      String call =
          buildMrStaticOrIntrinsicCall(
              targetOwner, targetEe, rawParamNames, rawParamTypes, 0, mrStaticTas, ctx);
      TypeMirror staticRetTm =
          mrStaticTas != null
              ? resolver.substituteTypeParams(
                  targetEe.getReturnType(), targetEe.getTypeParameters(), mrStaticTas)
              : targetEe.getReturnType();
      String retStmt =
          resolvedMrSamRet != null
              ? ("return "
                  + convert(
                      call,
                      resolver.goEffectiveType(staticRetTm),
                      resolvedMrSamRet,
                      samSig.returnWidenedToAny,
                      ctx))
              : call;
      String mrExpr = ctorPrefix + "(" + fnSig + " { " + retStmt + " })";
      if (qualEl instanceof TypeElement && !resolver.hasActiveTypeParameters()) {
        return ctx.hoistLambdaSingleton(fiGo, mrExpr);
      }
      return mrExpr;
    }

    // Case 4: Unbound instance method ref `ClassName::instanceMethod` (`mrp0_` is receiver!)
    if (qualEl instanceof TypeElement) {
      TypeMirror instRecvTm =
          (resolver.hasEffectiveTypeParameters(targetOwner) && !instSam.getParameterTypes().isEmpty())
              ? instSam.getParameterTypes().get(0)
              : targetOwner.asType();
      String recvExpr =
          convert(
              rawParamNames.get(0),
              rawParamTypes.get(0),
              instRecvTm,
              false,
              ctx);
      String call =
          buildMrInstanceCall(
              targetOwner, targetEe, recvExpr, instRecvTm, rawParamNames, rawParamTypes, 1, ctx);
      TypeResolver.MethodSig tSig = resolver.getInstanceMethodSig(targetEe);
      TypeMirror callRetTm =
          tSig != null && tSig.canonReturnType != null
              ? resolver.resolveSigReturnType(tSig, instRecvTm)
              : resolver.goEffectiveType(targetEe.getReturnType());
      String retStmt =
          resolvedMrSamRet != null
              ? ("return "
                  + convert(call, callRetTm, resolvedMrSamRet, samSig.returnWidenedToAny, ctx))
              : call;
      String mrExpr = ctorPrefix + "(" + fnSig + " { " + retStmt + " })";
      return !resolver.hasActiveTypeParameters()
          ? ctx.hoistLambdaSingleton(fiGo, mrExpr)
          : mrExpr;
    }

    // Case 5: Bound instance method ref `expr::instanceMethod` (evaluate receiver ONCE outside lambda!)
    String boundTmp = ctx.nextTemp("bound");
    String qualGo = emitPath(qualPath, ctx);
    String call =
        buildMrInstanceCall(
            targetOwner, targetEe, boundTmp, qualTm, rawParamNames, rawParamTypes, 0, ctx);
    TypeResolver.MethodSig tSig = resolver.getInstanceMethodSig(targetEe);
    TypeMirror callRetTm =
        tSig != null && tSig.canonReturnType != null
            ? resolver.resolveSigReturnType(tSig, qualTm)
            : resolver.goEffectiveType(targetEe.getReturnType());
    String retStmt =
        resolvedMrSamRet != null
            ? ("return "
                + convert(call, callRetTm, resolvedMrSamRet, samSig.returnWidenedToAny, ctx))
            : call;
    return "func() "
        + fiGo
        + " { "
        + boundTmp
        + " := "
        + qualGo
        + "; return "
        + ctorPrefix
        + "("
        + fnSig
        + " { "
        + retStmt
        + " }) }()";
  }

  private String emitMrObjectsEquals(
      List<String> rawParams,
      List<TypeMirror> rawParamTypes,
      int startIdx,
      EmitContext ctx) {
    TypeMirror raw0 = rawParamTypes.get(startIdx);
    TypeMirror raw1 = rawParamTypes.get(startIdx + 1);
    TypeMirror t0 = resolver.safeErasure(raw0);
    TypeMirror t1 = resolver.safeErasure(raw1);
    String p0 = rawParams.get(startIdx);
    String p1 = rawParams.get(startIdx + 1);
    if (t0.getKind() == t1.getKind()
        && t0.getKind().isPrimitive()
        && t0.getKind() != TypeKind.FLOAT
        && t0.getKind() != TypeKind.DOUBLE) {
      return "(" + p0 + " == " + p1 + ")";
    }
    if (resolver.isStringType(t0) && resolver.isStringType(t1)) {
      return "(" + p0 + " == " + p1 + ")";
    }
    if (resolver.isNodeType(t0) && resolver.isNodeType(t1)) {
      String l = convert(p0, raw0, t0, false, ctx);
      String r = convert(p1, raw1, t0, false, ctx);
      return "(" + l + " == " + r + ")";
    }
    if (resolver.isValueEnumType(t0) && resolver.types.isSameType(t0, t1)) {
      return "(" + p0 + " == " + p1 + ")";
    }
    if (resolver.isByteStringType(t0) && resolver.isByteStringType(t1)) {
      return rtPrefix(ctx) + "ByteStringObjEquals(" + p0 + ", " + p1 + ")";
    }
    if (resolver.isValueStructType(t0)
        && t0 instanceof DeclaredType d0s
        && t1 instanceof DeclaredType d1s
        && d0s.asElement() instanceof TypeElement te0s
        && te0s.equals(d1s.asElement())
        && valueStructGoEqMatchesJavaEquals(te0s)) {
      return "(" + p0 + " == " + p1 + ")";
    }
    String a0 = convert(p0, raw0, resolver.getObjectElement().asType(), true, ctx);
    String a1 = convert(p1, raw1, resolver.getObjectElement().asType(), true, ctx);
    return jlPrefix(ctx) + "ObjectsEquals(" + a0 + ", " + a1 + ")";
  }

  private String buildMrStaticOrIntrinsicCall(
      TypeElement owner,
      ExecutableElement ee,
      List<String> rawParams,
      List<TypeMirror> rawParamTypes,
      int startIdx,
      List<TypeMirror> genericMethodTypeArgs,
      EmitContext ctx) {
    String qn = owner.getQualifiedName().toString();
    String m = ee.getSimpleName().toString();
    int arity = rawParams.size() - startIdx;
    String rt = rtPrefix(ctx);

    if (qn.equals("java.util.Objects")) {
      if (m.equals("equals") && arity == 2) {
        return emitMrObjectsEquals(rawParams, rawParamTypes, startIdx, ctx);
      }
      if (m.equals("nonNull") && arity == 1) {
        if (resolver.isNodeType(rawParamTypes.get(startIdx))
            || resolver.isValueEnumType(rawParamTypes.get(startIdx))) {
          return "(" + rawParams.get(startIdx) + " != 0)";
        }
        if (resolver.isByteStringType(rawParamTypes.get(startIdx))) {
          return "!(" + rawParams.get(startIdx) + ").IsJavaNil_()";
        }
        if (resolver.isValueStructType(rawParamTypes.get(startIdx))) {
          return "("
              + rawParams.get(startIdx)
              + " != ("
              + resolver.goType(rawParamTypes.get(startIdx), ctx.imports)
              + "{}))";
        }
        return "("
            + convert(
                rawParams.get(startIdx),
                rawParamTypes.get(startIdx),
                resolver.getObjectElement().asType(),
                true,
                ctx)
            + " != nil)";
      }
      if (m.equals("isNull") && arity == 1) {
        if (resolver.isNodeType(rawParamTypes.get(startIdx))
            || resolver.isValueEnumType(rawParamTypes.get(startIdx))) {
          return "(" + rawParams.get(startIdx) + " == 0)";
        }
        if (resolver.isByteStringType(rawParamTypes.get(startIdx))) {
          return "(" + rawParams.get(startIdx) + ").IsJavaNil_()";
        }
        if (resolver.isValueStructType(rawParamTypes.get(startIdx))) {
          return "("
              + rawParams.get(startIdx)
              + " == ("
              + resolver.goType(rawParamTypes.get(startIdx), ctx.imports)
              + "{}))";
        }
        return "("
            + convert(
                rawParams.get(startIdx),
                rawParamTypes.get(startIdx),
                resolver.getObjectElement().asType(),
                true,
                ctx)
            + " == nil)";
      }
      if (m.equals("requireNonNull") && arity == 1) {
        return emitCheckNotNull(null, rawParams.get(startIdx), rawParamTypes.get(startIdx), ctx);
      }
    }
    if (qn.equals("com.google.common.base.Objects") && m.equals("equal") && arity == 2) {
      return emitMrObjectsEquals(rawParams, rawParamTypes, startIdx, ctx);
    }
    if (qn.equals("java.lang.String") && m.equals("valueOf") && arity == 1) {
      TypeMirror pTm = resolver.safeErasure(rawParamTypes.get(startIdx));
      if (resolver.isStringType(pTm)) {
        return rt + "StrValString(" + rawParams.get(startIdx) + ")";
      }
      if (pTm.getKind() == TypeKind.DOUBLE) {
        return jlPrefix(ctx) + "DoubleToString(" + rawParams.get(startIdx) + ")";
      }
      if (pTm.getKind() == TypeKind.FLOAT) {
        return jlPrefix(ctx) + "FloatToString(" + rawParams.get(startIdx) + ")";
      }
      if (pTm.getKind() == TypeKind.INT) {
        ctx.imports.std("strconv");
        return "strconv.FormatInt(int64(" + rawParams.get(startIdx) + "), 10)";
      }
      if (pTm.getKind() == TypeKind.LONG) {
        ctx.imports.std("strconv");
        return "strconv.FormatInt(" + rawParams.get(startIdx) + ", 10)";
      }
      if (pTm.getKind() == TypeKind.BOOLEAN) {
        return rt + "BoolToStr(" + rawParams.get(startIdx) + ")";
      }
      return rt + "StrVal(" + rawParams.get(startIdx) + ")";
    }
    if (qn.equals("com.google.common.base.Preconditions") && m.equals("checkNotNull") && arity == 1) {
      return emitCheckNotNull(null, rawParams.get(startIdx), rawParamTypes.get(startIdx), ctx);
    }
    if (qn.equals("java.lang.Integer") && m.equals("parseInt") && arity == 1) {
      return rt
          + "ParseInt("
          + convert(
              rawParams.get(startIdx),
              rawParamTypes.get(startIdx),
              resolver.elements.getTypeElement("java.lang.String").asType(),
              false,
              ctx)
          + ", 10)";
    }

    List<String> args = new ArrayList<>();
    appendMrArgs(rawParams, rawParamTypes, startIdx, ee, null, args, ctx, null, genericMethodTypeArgs);
    String fn = resolver.getStaticMethodGoName(ee);
    if (genericMethodTypeArgs != null) {
      List<String> goTas = new ArrayList<>(genericMethodTypeArgs.size());
      for (TypeMirror ta : genericMethodTypeArgs) {
        String gta = resolver.goType(ta, ctx.imports);
        goTas.add(gta.isEmpty() ? "any" : gta);
      }
      fn += "[" + String.join(", ", goTas) + "]";
    }
    if (resolver.isSourceType(owner)) {
      String pkg = resolver.planner.getGoPkg(owner);
      String alias = ctx.imports.ext(resolver.planner.getGoImportPath(pkg), pkg);
      return (alias.isEmpty() ? "" : alias + ".") + fn + "(" + String.join(", ", args) + ")";
    }
    return rt + fn + "(" + String.join(", ", args) + ")";
  }

  private String buildMrInstanceCall(
      TypeElement owner,
      ExecutableElement ee,
      String recvExpr,
      TypeMirror recvTm,
      List<String> rawParams,
      List<TypeMirror> rawParamTypes,
      int startIdx,
      EmitContext ctx) {
    String qn = owner.getQualifiedName().toString();
    String m = ee.getSimpleName().toString();
    int arity = rawParams.size() - startIdx;
    String jl = jlPrefix(ctx);
    String rt = rtPrefix(ctx);

    if (m.equals("equals")
        && arity == 1
        && ee.getParameters().size() == 1
        && ee.getReturnType().getKind() == TypeKind.BOOLEAN
        && resolver.safeErasure(recvTm) instanceof DeclaredType rdt
        && rdt.asElement() instanceof TypeElement rte
        && resolver.isValueStructClass(rte)
        && valueStructGoEqMatchesJavaEquals(rte)
        && resolver.safeErasure(rawParamTypes.get(startIdx)) instanceof DeclaredType adt
        && rte.equals(adt.asElement())) {
      return "(" + emitCheckNotNull(null, recvExpr, recvTm, ctx) + " == " + rawParams.get(startIdx) + ")";
    }
    if (qn.equals("java.lang.Object")) {
      if (m.equals("toString") && arity == 0) return jl + "ToJavaString(" + rt + "ToAny(" + recvExpr + "))";
      if (m.equals("hashCode") && arity == 0) return jl + "ObjectHashCode(" + rt + "ToAny(" + recvExpr + "))";
      if (m.equals("equals") && arity == 1) {
        return jl
            + "ObjectsEquals("
            + rt
            + "ToAny("
            + recvExpr
            + "), "
            + convert(
                rawParams.get(startIdx),
                rawParamTypes.get(startIdx),
                resolver.getObjectElement().asType(),
                true,
                ctx)
            + ")";
      }
    }
    if (qn.equals("java.lang.String")) {
      if (m.equals("toString") && arity == 0) return rt + "CheckNotNullString(" + recvExpr + ")";
      if (m.equals("hashCode") && arity == 0) return jl + "StringHashCode(" + rt + "CheckNotNullString(" + recvExpr + "))";
      if (m.equals("isEmpty") && arity == 0) return "(" + rt + "CheckNotNullString(" + recvExpr + ") == \"\")";
      if (m.equals("length") && arity == 0) return rt + "JStrLen(" + recvExpr + ")";
      if (m.equals("trim") && arity == 0) return jl + "StringTrim(" + rt + "CheckNotNullString(" + recvExpr + "))";
      if (m.equals("toLowerCase") && arity == 0) return jl + "StringToLowerCaseASCII(" + rt + "CheckNotNullString(" + recvExpr + "))";
      if (m.equals("toUpperCase") && arity == 0) return jl + "StringToUpperCaseASCII(" + rt + "CheckNotNullString(" + recvExpr + "))";
      if (m.equals("compareTo") && arity == 1) {
        String a0 =
            convert(
                rawParams.get(startIdx),
                rawParamTypes.get(startIdx),
                resolver.elements.getTypeElement("java.lang.String").asType(),
                false,
                ctx);
        return rt + "JStrCompareTo(" + recvExpr + ", " + a0 + ")";
      }
      if (m.equals("concat") && arity == 1) {
        String a0 =
            convert(
                rawParams.get(startIdx),
                rawParamTypes.get(startIdx),
                resolver.elements.getTypeElement("java.lang.String").asType(),
                false,
                ctx);
        return "(" + rt + "CheckNotNullString(" + recvExpr + ") + " + rt + "CheckNotNullString(" + a0 + "))";
      }
      if (m.equals("startsWith") && arity == 1) {
        String a0 =
            convert(
                rawParams.get(startIdx),
                rawParamTypes.get(startIdx),
                resolver.elements.getTypeElement("java.lang.String").asType(),
                false,
                ctx);
        return jl + "StringStartsWith(" + rt + "CheckNotNullString(" + recvExpr + "), " + rt + "CheckNotNullString(" + a0 + "))";
      }
      if (m.equals("endsWith") && arity == 1) {
        String a0 =
            convert(
                rawParams.get(startIdx),
                rawParamTypes.get(startIdx),
                resolver.elements.getTypeElement("java.lang.String").asType(),
                false,
                ctx);
        return jl + "StringEndsWith(" + rt + "CheckNotNullString(" + recvExpr + "), " + rt + "CheckNotNullString(" + a0 + "))";
      }
      if (m.equals("contains") && arity == 1) {
        String a0 =
            convert(
                rawParams.get(startIdx),
                rawParamTypes.get(startIdx),
                resolver.getObjectElement().asType(),
                true,
                ctx);
        return "(" + rt + "JStrIndexOf(" + recvExpr + ", " + rt + "StrVal(" + a0 + ")) >= 0)";
      }
      if (m.equals("equals") && arity == 1) {
        return rt
            + "StringEquals("
            + recvExpr
            + ", "
            + convert(
                rawParams.get(startIdx),
                rawParamTypes.get(startIdx),
                resolver.getObjectElement().asType(),
                true,
                ctx)
            + ")";
      }
    }
    if (qn.equals("java.lang.Enum")) {
      if (m.equals("name") && arity == 0) return rt + "EnumName(" + recvExpr + ")";
      if (m.equals("ordinal") && arity == 0) return rt + "EnumOrdinal(" + recvExpr + ")";
    }
    if (arity == 1
        && ee.getParameters().size() == 1
        && ee.getReturnType().getKind() == TypeKind.BOOLEAN
        && !ee.getModifiers().contains(Modifier.PRIVATE)
        && resolver.isNodeType(rawParamTypes.get(startIdx))) {
      if ((qn.equals("com.google.common.base.Predicate") && (m.equals("apply") || m.equals("test")))
          || (qn.equals("java.util.function.Predicate") && m.equals("test"))) {
        return rt + "PredTestNode(" + recvExpr + ", uint32(" + rawParams.get(startIdx) + "))";
      }
      if (m.equals("contains")
          && (isSubtypeOfQualified(owner, "java.util.Set")
              || isSubtypeOfQualified(owner, "com.google.common.collect.ImmutableCollection"))) {
        return rt + "SetContainsNode(" + recvExpr + ", " + rawParams.get(startIdx) + ")";
      }
    }

    TypeResolver.MethodSig tSig = resolver.getInstanceMethodSig(ee);
    List<String> args = new ArrayList<>();
    appendMrArgs(rawParams, rawParamTypes, startIdx, ee, tSig, args, ctx, recvTm, null);
    if (ee.getModifiers().contains(Modifier.PRIVATE)
        && resolver.isSourceType(owner)
        && resolver.isSubclassed(owner)) {
      return recvExpr
          + ".Get"
          + resolver.planner.getGoTypeName(owner)
          + "Base()."
          + tSig.goName
          + "("
          + String.join(", ", args)
          + ")";
    }
    return recvExpr + "." + tSig.goName + "(" + String.join(", ", args) + ")";
  }

  private void appendMrArgs(
      List<String> rawParams,
      List<TypeMirror> rawParamTypes,
      int startIdx,
      ExecutableElement targetEe,
      TypeResolver.MethodSig targetSig,
      List<String> outArgs,
      EmitContext ctx,
      TypeMirror ownerSiteTm,
      List<TypeMirror> genericMethodTypeArgs) {
    int targetParamCount = targetEe.getParameters().size();
    boolean varArgs = targetEe.isVarArgs() && targetParamCount > 0;
    int fixed = varArgs ? (targetParamCount - 1) : targetParamCount;
    TypeElement eeOwner = targetEe.getEnclosingElement() instanceof TypeElement t ? t : null;

    for (int i = 0; i < fixed && (startIdx + i) < rawParams.size(); i++) {
      TypeMirror srcTm = rawParamTypes.get(startIdx + i);
      boolean forceAny = targetSig != null && targetSig.paramWidenedToAny[i];
      TypeMirror dstTm;
      if (targetSig != null && (genericMethodTypeArgs == null || targetEe.getTypeParameters().isEmpty())) {
        dstTm = resolver.resolveSigParamType(targetSig, i, ownerSiteTm);
      } else {
        TypeMirror baseDst =
            targetSig != null
                ? targetSig.canonParamTypes.get(i)
                : targetEe.getParameters().get(i).asType();
        if (genericMethodTypeArgs != null && !targetEe.getTypeParameters().isEmpty()) {
          baseDst =
              resolver.substituteTypeParams(
                  baseDst, targetEe.getTypeParameters(), genericMethodTypeArgs);
        } else if (ownerSiteTm != null
            && eeOwner != null
            && resolver.isNativeGenericClass(eeOwner)
            && !forceAny) {
          baseDst = resolver.substituteOwnerTypeArgs(baseDst, eeOwner, ownerSiteTm);
        }
        dstTm = resolver.goEffectiveType(baseDst);
      }
      outArgs.add(convert(rawParams.get(startIdx + i), srcTm, dstTm, forceAny, ctx));
    }

    if (varArgs) {
      boolean lastForceAny = targetSig != null && targetSig.paramWidenedToAny[targetParamCount - 1];
      TypeMirror effLast;
      if (targetSig != null && (genericMethodTypeArgs == null || targetEe.getTypeParameters().isEmpty())) {
        effLast = resolver.resolveSigParamType(targetSig, targetParamCount - 1, ownerSiteTm);
      } else {
        TypeMirror lastParam =
            targetSig != null
                ? targetSig.canonParamTypes.get(targetParamCount - 1)
                : targetEe.getParameters().get(targetParamCount - 1).asType();
        if (genericMethodTypeArgs != null && !targetEe.getTypeParameters().isEmpty()) {
          lastParam =
              resolver.substituteTypeParams(
                  lastParam, targetEe.getTypeParameters(), genericMethodTypeArgs);
        } else if (ownerSiteTm != null
            && eeOwner != null
            && resolver.isNativeGenericClass(eeOwner)
            && !lastForceAny) {
          lastParam = resolver.substituteOwnerTypeArgs(lastParam, eeOwner, ownerSiteTm);
        }
        effLast = resolver.goEffectiveType(lastParam);
      }
      ArrayType at = (ArrayType) effLast;
      if ((rawParams.size() - startIdx) == targetParamCount
          && rawParamTypes.get(startIdx + fixed).getKind() == TypeKind.ARRAY) {
        outArgs.add(
            convert(
                rawParams.get(startIdx + fixed),
                rawParamTypes.get(startIdx + fixed),
                at,
                false,
                ctx));
      } else {
        List<String> packed = new ArrayList<>();
        for (int j = startIdx + fixed; j < rawParams.size(); j++) {
          packed.add(
              convert(
                  rawParams.get(j),
                  rawParamTypes.get(j),
                  at.getComponentType(),
                  false,
                  ctx));
        }
        outArgs.add(resolver.goType(at, ctx.imports) + "{" + String.join(", ", packed) + "}");
      }
    }
  }

  private TypeMirror extractNodePredicateTypeArg(
      TypeMirror fiTm, TypeElement fiTe, TypeResolver.MethodSig samSig) {
    if (fiTe == null || resolver.isSourceType(fiTe)) {
      return null;
    }
    String qn = fiTe.getQualifiedName().toString();
    if (!qn.equals("com.google.common.base.Predicate")
        && !qn.equals("java.util.function.Predicate")) {
      return null;
    }
    if (samSig == null
        || samSig.canonParamTypes.size() != 1
        || samSig.canonReturnType == null
        || samSig.canonReturnType.getKind() != TypeKind.BOOLEAN
        || samSig.returnWidenedToAny) {
      return null;
    }
    if (!(fiTm instanceof DeclaredType dt) || dt.getTypeArguments().size() != 1) {
      return null;
    }
    TypeMirror ta0 = dt.getTypeArguments().get(0);
    if (ta0 instanceof DeclaredType tadt
        && tadt.asElement() instanceof TypeElement tate
        && resolver.isNodeClass(tate)) {
      return ta0;
    }
    return null;
  }

  private String nodeLambdaCtorPrefix(TypeElement fiTe, EmitContext ctx) {
    return rtPrefix(ctx) + "NewNodeLambda_" + resolver.runtime.extGoTypeName(fiTe);
  }

  private String lambdaCtorPrefix(TypeElement fiTe, TypeMirror fiTm, EmitContext ctx) {
    if (resolver.isSourceType(fiTe)) {
      String pkg = resolver.planner.getGoPkg(fiTe);
      String alias = ctx.imports.ext(resolver.planner.getGoImportPath(pkg), pkg);
      String typeArgs = resolver.formatSuperTypeArgs(fiTm, fiTe, ctx.imports);
      return (alias.isEmpty() ? "" : alias + ".")
          + "NewLambda_"
          + resolver.planner.getGoTypeName(fiTe)
          + typeArgs;
    }
    return rtPrefix(ctx) + "NewLambda_" + resolver.runtime.extGoTypeName(fiTe);
  }

  private void copyLocalsIntoLambda(EmitContext src, EmitContext dst) {
    // Copy local variable table via helper on EmitContext
    dst.copyLocalsFrom(src);
  }

  private String rtPrefix(EmitContext ctx) {
    String a = resolver.runtime.j2gortAlias(ctx.imports);
    return a.isEmpty() ? "" : (a + ".");
  }

  private String jlPrefix(EmitContext ctx) {
    String a = resolver.runtime.javalangAlias(ctx.imports);
    return a.isEmpty() ? "" : (a + ".");
  }

  private String gvPrefix(EmitContext ctx) {
    String a = resolver.runtime.guavaAlias(ctx.imports);
    return a.isEmpty() ? "" : (a + ".");
  }

  private String nsPrefix(EmitContext ctx) {
    String a = resolver.runtime.nodestoreAlias(ctx.imports);
    return a.isEmpty() ? "" : (a + ".");
  }
}
