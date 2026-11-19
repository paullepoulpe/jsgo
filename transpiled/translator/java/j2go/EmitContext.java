package j2go;

import com.sun.source.tree.ClassTree;
import com.sun.source.util.TreePath;
import java.util.ArrayDeque;
import java.util.ArrayList;
import java.util.Deque;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import javax.lang.model.element.ExecutableElement;
import javax.lang.model.element.TypeElement;
import javax.lang.model.element.VariableElement;
import javax.lang.model.type.TypeMirror;

/**
 * Per-file and per-method state shared across {@link ClassEmitter}, {@link MethodEmitter},
 * {@link StmtEmitter}, and {@link ExprEmitter}.
 */
public final class EmitContext {
  public final TypeResolver resolver;
  public final PackagePlanner.UnitPlan unitPlan;
  public final PackagePlanner.ImportTracker imports;
  public final List<DeferredLocalType> deferredLocalTypes = new ArrayList<>();
  private final List<String> hoistedVarDecls;
  private final List<String> hoistedInitStmts;
  private final int[] hoistedSeq;
  private final String unitSingletonStem;

  public TypeElement currentType;
  public ExecutableElement currentExec;
  public TypeMirror expectedReturnType;
  public boolean returnWidenedToAny;
  public boolean inStaticContext;
  public String receiverName = "this_";

  // Switch expression target stack (`yield expr` / arrow result inside a SwitchExpressionTree)
  public record SwitchExprTarget(String assignVar, TypeMirror targetType, boolean forceAny) {}

  public final Deque<SwitchExprTarget> switchExprTargetStack = new ArrayDeque<>();

  // Try-catch wrapper return/control-flow propagation state
  public boolean inTryWrapper;
  public boolean tryNeedsReturnProp;
  public String tryRetValVar;
  public String tryDidReturnVar;
  public String tryFlowVar; // non-null if break/continue out of enclosing loop inside try IIFE

  // Active loop and switch depth inside current IIFE (0 means outer loop/switch is outside current try IIFE)
  public int loopDepthInTryIife;
  public int switchDepthInTryIife;

  // Local variable Go names within current method/initializer
  private final Map<VariableElement, String> localNames = new HashMap<>();
  private final Set<String> usedIdentNames = new HashSet<>();
  public final Set<VariableElement> predeclaredPatternVars = new HashSet<>();
  // Captured outer local variables accessible on `this_` when inside a local/anonymous class
  public final Map<VariableElement, String> capturedVarFields = new HashMap<>();

  private int tempSeq = 0;

  public record DeferredLocalType(
      ClassTree tree,
      TreePath path,
      TypeElement element,
      Map<VariableElement, String> enclosingCapturedFields) {}

  public EmitContext(
      TypeResolver resolver,
      PackagePlanner.UnitPlan unitPlan,
      PackagePlanner.ImportTracker imports) {
    this(
        resolver,
        unitPlan,
        imports,
        new ArrayList<>(),
        new ArrayList<>(),
        new int[1],
        computeUnitSingletonStem(unitPlan));
  }

  private EmitContext(
      TypeResolver resolver,
      PackagePlanner.UnitPlan unitPlan,
      PackagePlanner.ImportTracker imports,
      List<String> hoistedVarDecls,
      List<String> hoistedInitStmts,
      int[] hoistedSeq,
      String unitSingletonStem) {
    this.resolver = resolver;
    this.unitPlan = unitPlan;
    this.imports = imports;
    this.hoistedVarDecls = hoistedVarDecls;
    this.hoistedInitStmts = hoistedInitStmts;
    this.hoistedSeq = hoistedSeq;
    this.unitSingletonStem = unitSingletonStem;
  }

  private static String computeUnitSingletonStem(PackagePlanner.UnitPlan unitPlan) {
    String fn = unitPlan != null && unitPlan.goFileName != null ? unitPlan.goFileName : "unit";
    if (fn.endsWith(".go")) {
      fn = fn.substring(0, fn.length() - 3);
    }
    return GoNames.cleanChars(fn);
  }

  public EmitContext copyForMember(TypeElement type, boolean isStatic) {
    EmitContext sub =
        new EmitContext(
            resolver,
            unitPlan,
            imports,
            hoistedVarDecls,
            hoistedInitStmts,
            hoistedSeq,
            unitSingletonStem);
    sub.currentType = type;
    sub.inStaticContext = isStatic;
    sub.tempSeq = this.tempSeq;
    return sub;
  }

  public EmitContext copyForLambda() {
    EmitContext sub =
        new EmitContext(
            resolver,
            unitPlan,
            imports,
            hoistedVarDecls,
            hoistedInitStmts,
            hoistedSeq,
            unitSingletonStem);
    sub.currentType = this.currentType;
    sub.currentExec = this.currentExec;
    sub.inStaticContext = this.inStaticContext;
    sub.receiverName = this.receiverName;
    sub.tempSeq = this.tempSeq;
    return sub;
  }

  public String hoistLambdaSingleton(String goType, String initExpr) {
    String name = "lambdaSingleton_" + unitSingletonStem + "_" + (++hoistedSeq[0]);
    hoistedVarDecls.add("var " + name + " " + goType + "\n");
    hoistedInitStmts.add("\t" + name + " = " + initExpr + "\n");
    return name;
  }

  public void appendHoistedDecls(StringBuilder out) {
    if (hoistedVarDecls.isEmpty()) {
      return;
    }
    out.append('\n');
    for (String vd : hoistedVarDecls) {
      out.append(vd);
    }
    out.append("func init() {\n");
    for (String st : hoistedInitStmts) {
      out.append(st);
    }
    out.append("}\n");
  }

  public void syncTempSeqFrom(EmitContext sub) {
    if (sub.tempSeq > this.tempSeq) {
      this.tempSeq = sub.tempSeq;
    }
  }

  public void copyLocalsFrom(EmitContext parent) {
    this.localNames.putAll(parent.localNames);
    this.usedIdentNames.addAll(parent.usedIdentNames);
  }

  public String nextTemp(String prefix) {
    return prefix + "_" + (++tempSeq);
  }

  public String registerLocalVar(VariableElement ve) {
    String existing = localNames.get(ve);
    if (existing != null) {
      return existing;
    }
    String base = GoNames.localIdent(ve.getSimpleName().toString());
    if (base.equals(receiverName)) {
      base = base + "_v";
    }
    String cand = base;
    int idx = 2;
    while (!usedIdentNames.add(cand)) {
      cand = base + "_" + idx++;
    }
    localNames.put(ve, cand);
    return cand;
  }

  public String lookupLocalVar(VariableElement ve) {
    String local = localNames.get(ve);
    if (local != null) {
      return local;
    }
    String capField = capturedVarFields.get(ve);
    return capField != null ? (receiverName + "." + capField) : null;
  }

  /** Expression representing Java `this` as a value of `currentType`. */
  public String selfValueExpr() {
    if (currentType != null && resolver.isSubclassed(currentType)) {
      if (resolver.getSamePkgReusableSelfRoot(currentType) != null) {
        return receiverName + ".self_" + resolver.planner.getGoTypeName(currentType) + "_()";
      }
      return receiverName + ".self_";
    }
    return receiverName;
  }

  public void deferLocalOrAnonClass(ClassTree ct, TreePath path, TypeElement te) {
    deferredLocalTypes.add(
        new DeferredLocalType(ct, path, te, new HashMap<>(this.capturedVarFields)));
  }
}
