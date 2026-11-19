package j2go;

import com.sun.source.tree.AssignmentTree;
import com.sun.source.tree.BinaryTree;
import com.sun.source.tree.BlockTree;
import com.sun.source.tree.ClassTree;
import com.sun.source.tree.CompilationUnitTree;
import com.sun.source.tree.ConditionalExpressionTree;
import com.sun.source.tree.DoWhileLoopTree;
import com.sun.source.tree.EnhancedForLoopTree;
import com.sun.source.tree.ExpressionStatementTree;
import com.sun.source.tree.ExpressionTree;
import com.sun.source.tree.ForLoopTree;
import com.sun.source.tree.IdentifierTree;
import com.sun.source.tree.IfTree;
import com.sun.source.tree.LambdaExpressionTree;
import com.sun.source.tree.LiteralTree;
import com.sun.source.tree.MemberReferenceTree;
import com.sun.source.tree.MemberSelectTree;
import com.sun.source.tree.MethodInvocationTree;
import com.sun.source.tree.MethodTree;
import com.sun.source.tree.NewClassTree;
import com.sun.source.tree.ParenthesizedTree;
import com.sun.source.tree.ReturnTree;
import com.sun.source.tree.StatementTree;
import com.sun.source.tree.SynchronizedTree;
import com.sun.source.tree.ThrowTree;
import com.sun.source.tree.Tree;
import com.sun.source.tree.UnaryTree;
import com.sun.source.tree.VariableTree;
import com.sun.source.tree.WhileLoopTree;
import com.sun.source.util.JavacTask;
import com.sun.source.util.TreePath;
import com.sun.source.util.TreePathScanner;
import com.sun.source.util.Trees;
import java.util.ArrayDeque;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Comparator;
import java.util.Deque;
import java.util.HashMap;
import java.util.HashSet;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import javax.lang.model.element.Element;
import javax.lang.model.element.ElementKind;
import javax.lang.model.element.ExecutableElement;
import javax.lang.model.element.Modifier;
import javax.lang.model.element.NestingKind;
import javax.lang.model.element.TypeElement;
import javax.lang.model.element.TypeParameterElement;
import javax.lang.model.element.VariableElement;
import javax.lang.model.type.ArrayType;
import javax.lang.model.type.DeclaredType;
import javax.lang.model.type.ExecutableType;
import javax.lang.model.type.IntersectionType;
import javax.lang.model.type.TypeKind;
import javax.lang.model.type.TypeMirror;
import javax.lang.model.type.TypeVariable;
import javax.lang.model.type.UnionType;
import javax.lang.model.type.WildcardType;
import javax.lang.model.util.Elements;
import javax.lang.model.util.Types;

/**
 * Resolves Java AST types, inheritance hierarchies, Union-Find canonical virtual method
 * signatures, overloaded method/constructor names, field shadowing, cross-package field
 * accessors, functional interfaces, and local/anonymous captured variables.
 */
public final class TypeResolver {

  public static final String NODE_QUAL_NAME = "com.google.javascript.rhino.Node";

  private static final Set<String> NODE_SUBCLASS_QUAL_NAMES =
      Set.of(
          "com.google.javascript.rhino.Node.NumberNode",
          "com.google.javascript.rhino.Node.BigIntNode",
          "com.google.javascript.rhino.Node.StringNode",
          "com.google.javascript.rhino.Node.TemplateLiteralSubstringNode");

  private static final Set<String> VALUE_STRUCT_QUAL_NAMES =
      Set.of(
          "com.google.javascript.jscomp.parsing.parser.util.SourcePosition",
          "com.google.javascript.jscomp.parsing.parser.util.SourceRange",
          "com.google.javascript.rhino.jstype.FunctionType.Parameter",
          "com.google.javascript.rhino.jstype.JSType.TypePair",
          "com.google.javascript.rhino.JSTypeExpression",
          "com.google.javascript.jscomp.parsing.parser.Token",
          "com.google.javascript.jscomp.RemoveUnusedCode.Continuation",
          "com.google.javascript.jscomp.parsing.JsDocInfoParser.ExtractionInfo",
          "com.google.javascript.jscomp.parsing.JsDocInfoParser.ExtendedTypeInfo",
          "com.google.javascript.jscomp.colors.ColorId");

  private static final Set<String> NULLABLE_VALUE_STRUCT_QUAL_NAMES =
      Set.of(
          "com.google.javascript.rhino.jstype.FunctionType.Parameter",
          "com.google.javascript.rhino.jstype.JSType.TypePair",
          "com.google.javascript.rhino.JSTypeExpression",
          "com.google.javascript.jscomp.parsing.parser.Token",
          "com.google.javascript.jscomp.RemoveUnusedCode.Continuation",
          "com.google.javascript.jscomp.parsing.JsDocInfoParser.ExtractionInfo",
          "com.google.javascript.jscomp.parsing.JsDocInfoParser.ExtendedTypeInfo",
          "com.google.javascript.jscomp.colors.ColorId");

  private static final Set<String> VALID_FLAG_VALUE_STRUCT_QUAL_NAMES =
      Set.of(
          "com.google.javascript.rhino.jstype.JSType.TypePair",
          "com.google.javascript.rhino.JSTypeExpression",
          "com.google.javascript.jscomp.colors.ColorId");

  private static final Set<String> NATIVE_GENERIC_QUAL_NAMES =
      Set.of(
          "com.google.javascript.jscomp.base.LinkedIdentityHashMap",
          "com.google.javascript.jscomp.base.LinkedIdentityHashMap.Entry",
          "com.google.javascript.jscomp.base.LinkedIdentityHashSet",
          "com.google.javascript.jscomp.SimpleMultimap",
          "com.google.javascript.jscomp.SimpleTable",
          "com.google.javascript.jscomp.SimpleTable.SimpleCell",
          "com.google.javascript.rhino.PMap",
          "com.google.javascript.rhino.PMap.Reconciler",
          "com.google.javascript.rhino.HamtPMap",
          "com.google.javascript.rhino.HamtPMap.Holder",
          "com.google.javascript.rhino.HamtPMap.KeyIter",
          "com.google.javascript.rhino.HamtPMap.ValueIter",
          "com.google.javascript.jscomp.graph.SubGraph",
          "com.google.javascript.jscomp.graph.AdjacencyGraph",
          "com.google.javascript.jscomp.graph.GraphNode",
          "com.google.javascript.jscomp.graph.Graph",
          "com.google.javascript.jscomp.graph.Graph.GraphEdge",
          "com.google.javascript.jscomp.graph.Graph.SimpleSubGraph",
          "com.google.javascript.jscomp.graph.DiGraph",
          "com.google.javascript.jscomp.graph.DiGraph.DiGraphNode",
          "com.google.javascript.jscomp.graph.DiGraph.DiGraphEdge",
          "com.google.javascript.jscomp.graph.LinkedDirectedGraph",
          "com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode",
          "com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphEdge",
          "com.google.javascript.jscomp.graph.UndiGraph",
          "com.google.javascript.jscomp.graph.UndiGraph.UndiGraphNode",
          "com.google.javascript.jscomp.graph.UndiGraph.UndiGraphEdge",
          "com.google.javascript.jscomp.graph.LinkedUndirectedGraph",
          "com.google.javascript.jscomp.graph.LinkedUndirectedGraph.LinkedUndirectedGraphNode",
          "com.google.javascript.jscomp.graph.LinkedUndirectedGraph.LinkedUndirectedGraphEdge",
          "com.google.javascript.jscomp.graph.DominatorTree",
          "com.google.javascript.jscomp.graph.FixedPointGraphTraversal",
          "com.google.javascript.jscomp.graph.FixedPointGraphTraversal.EdgeCallback",
          "com.google.javascript.jscomp.graph.GraphReachability",
          "com.google.javascript.jscomp.graph.GraphReachability.EdgeTuple",
          "com.google.javascript.jscomp.graph.CheckPathsBetweenNodes",
          "com.google.javascript.jscomp.graph.LowestCommonAncestorFinder",
          "com.google.javascript.jscomp.graph.LowestCommonAncestorFinder.Factory",
          "com.google.javascript.jscomp.graph.GraphColoring",
          "com.google.javascript.jscomp.graph.GraphColoring.GreedyGraphColoring",
          "com.google.javascript.jscomp.graph.UnionFind",
          "com.google.javascript.jscomp.graph.StandardUnionFind",
          "com.google.javascript.jscomp.graph.StandardUnionFind.Node",
          "com.google.javascript.jscomp.ControlFlowGraph",
          "com.google.javascript.jscomp.DataFlowAnalysis",
          "com.google.javascript.jscomp.DataFlowAnalysis.FlowJoiner",
          "com.google.javascript.jscomp.DataFlowAnalysis.FlowBrancher",
          "com.google.javascript.jscomp.DataFlowAnalysis.UniqueQueue");

  private static final Set<String> RESERVED_METHOD_NAMES =
      Set.of("String", "Equals", "HashCode", "Error");

  /** Canonical Go signature of a virtual method slot unified across overrides. */
  public static final class MethodSig {
    public String goName;
    public final String simpleName;
    public TypeElement canonOwner; // Root generic owner whose type params appear in canon*Types
    public final List<TypeMirror> canonParamTypes; // null entry means widened to Go `any`
    public TypeMirror canonReturnType; // null means void; Object TypeMirror if widened to `any`
    public boolean returnWidenedToAny;
    public final boolean[] paramWidenedToAny;

    MethodSig(
        String simpleName,
        int arity,
        TypeElement canonOwner,
        TypeMirror initRet,
        List<TypeMirror> initParams) {
      this.simpleName = simpleName;
      this.canonOwner = canonOwner;
      this.canonParamTypes = new ArrayList<>(initParams);
      this.canonReturnType = initRet;
      this.paramWidenedToAny = new boolean[arity];
    }
  }

  public final JavacTask task;
  public final Trees trees;
  public final Elements elements;
  public final Types types;
  public final PackagePlanner planner;
  public final RuntimeBindings runtime;

  private final TypeElement objectElement;
  private final TypeElement throwableElement;
  private final TypeElement enumElement;

  private final Set<TypeElement> subclassedTypes = new LinkedHashSet<>();
  private final Set<TypeElement> referencedExternalTypes = new LinkedHashSet<>();
  private final Set<ExecutableElement> referencedExternalExecs = new LinkedHashSet<>();
  private final Set<VariableElement> referencedExternalFields = new LinkedHashSet<>();
  private final Set<TypeElement> lambdaTargetInterfaces = new LinkedHashSet<>();

  // Anonymous class NewClassTree -> resolved anonymous TypeElement
  private final Map<ClassTree, TypeElement> classTreeElements = new LinkedHashMap<>();
  private final Map<TypeElement, ClassTree> typeElementClassTrees = new LinkedHashMap<>();
  // Enum constant VariableElement with class body -> anonymous TypeElement
  private final Map<VariableElement, TypeElement> enumConstantBodyTypes = new LinkedHashMap<>();
  private final Set<TypeElement> enumConstantAnonTypes = new HashSet<>();
  // In-tree enums lowered to `uint8` (`0 = null, 1..N = ordinal + 1`) and their constant fields
  public static final Object ENUM_NULL_STRING = new Object();
  private final Set<TypeElement> valueEnumTypes = new HashSet<>();
  private final Map<TypeElement, Map<VariableElement, List<Object>>> valueEnumFieldTables =
      new LinkedHashMap<>();
  // Direct non-overridden virtual call slots on subclassed base classes (`(base, sig)`)
  private final Map<TypeElement, Set<MethodSig>> directBaseDispatchSigs = new HashMap<>();
  // Subclassed source base classes that may emit `this_.self_`
  private final Set<TypeElement> mayHaveSelfFieldTypes = new HashSet<>();
  // Ancestor base classes whose `self_` field may be reused by same-package subclasses
  private final Set<TypeElement> reusableSelfAncestors = new HashSet<>();
  // Non-reentrant `static synchronized` and leaf `synchronized` / `synchronized (this)` types
  private final Set<TypeElement> fastStaticSyncTypes = new HashSet<>();
  private final Set<TypeElement> fastInstanceSyncTypes = new HashSet<>();

  // Captured local variables for local & anonymous classes
  private final Map<TypeElement, List<VariableElement>> capturedVars = new LinkedHashMap<>();

  // Field naming and cross-package access tracking
  private final Set<VariableElement> crossPkgReadFields = new LinkedHashSet<>();
  private final Set<VariableElement> crossPkgWriteFields = new LinkedHashSet<>();
  private final Map<VariableElement, String> fieldGoNames = new LinkedHashMap<>();
  private final Map<TypeElement, String> valueStructEnumSentinelFields = new HashMap<>();

  // Constructor and method Go names / unified virtual slots
  private final Map<ExecutableElement, String> constructorGoNames = new LinkedHashMap<>();
  private final Map<ExecutableElement, String> staticMethodGoNames = new LinkedHashMap<>();
  private final Map<ExecutableElement, MethodSig> instanceMethodSigs = new LinkedHashMap<>();
  private final Set<MethodSig> neededExternalSigs = new LinkedHashSet<>();
  // Source generic interfaces whose k-th type parameter is instantiated with `?` somewhere
  private final Map<TypeElement, Set<Integer>> wildcardedInterfaceParams = new HashMap<>();

  public record GenericScope(TypeElement ownerType, Map<Element, String> typeParamNames) {}

  private final ThreadLocal<Deque<GenericScope>> genericScopeStack =
      ThreadLocal.withInitial(ArrayDeque::new);

  public TypeResolver(JavacTask task, PackagePlanner planner) {
    this.task = task;
    this.trees = Trees.instance(task);
    this.elements = task.getElements();
    this.types = task.getTypes();
    this.planner = planner;
    this.objectElement = elements.getTypeElement("java.lang.Object");
    this.throwableElement = elements.getTypeElement("java.lang.Throwable");
    this.enumElement = elements.getTypeElement("java.lang.Enum");
    this.runtime = new RuntimeBindings(this);
  }

  public TypeElement getObjectElement() {
    return objectElement;
  }

  public TypeElement getThrowableElement() {
    return throwableElement;
  }

  public boolean isSourceType(TypeElement te) {
    return planner.isSourceType(te);
  }

  public boolean isNodeClass(TypeElement te) {
    if (te == null) {
      return false;
    }
    String qn = te.getQualifiedName().toString();
    return qn.equals(NODE_QUAL_NAME) || NODE_SUBCLASS_QUAL_NAMES.contains(qn);
  }

  public int getNodeKindFlag(TypeElement te) {
    return switch (te.getQualifiedName().toString()) {
      case "com.google.javascript.rhino.Node.NumberNode" -> 1;
      case "com.google.javascript.rhino.Node.BigIntNode" -> 2;
      case "com.google.javascript.rhino.Node.StringNode" -> 3;
      case "com.google.javascript.rhino.Node.TemplateLiteralSubstringNode" -> 4;
      default -> 0;
    };
  }

  public boolean isNodeType(TypeMirror tm) {
    if (tm == null) {
      return false;
    }
    TypeMirror er = safeErasure(tm);
    if (er instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      return isNodeClass(te);
    }
    return false;
  }

  public static boolean isValueStructClass(String qn) {
    return qn != null && VALUE_STRUCT_QUAL_NAMES.contains(qn);
  }

  public boolean isValueStructClass(TypeElement te) {
    return te != null
        && isValueStructClass(te.getQualifiedName().toString())
        && !subclassedTypes.contains(te)
        && !te.getModifiers().contains(Modifier.ABSTRACT);
  }

  public boolean isValueStructType(TypeMirror tm) {
    if (tm == null) {
      return false;
    }
    TypeMirror er = safeErasure(tm);
    if (er instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      return isValueStructClass(te);
    }
    return false;
  }

  public static boolean isNullableValueStructClass(String qn) {
    return qn != null && NULLABLE_VALUE_STRUCT_QUAL_NAMES.contains(qn);
  }

  public boolean isNullableValueStructClass(TypeElement te) {
    return te != null
        && isNullableValueStructClass(te.getQualifiedName().toString())
        && !subclassedTypes.contains(te)
        && !te.getModifiers().contains(Modifier.ABSTRACT);
  }

  public boolean isNullableValueStructType(TypeMirror tm) {
    if (tm == null) {
      return false;
    }
    TypeMirror er = safeErasure(tm);
    if (er instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      return isNullableValueStructClass(te);
    }
    return false;
  }

  public boolean hasValueStructValidFlag(TypeElement te) {
    return isNullableValueStructClass(te)
        && VALID_FLAG_VALUE_STRUCT_QUAL_NAMES.contains(te.getQualifiedName().toString());
  }

  public String getValueStructEnumSentinelFieldGoName(TypeElement te) {
    return te != null ? valueStructEnumSentinelFields.get(te) : null;
  }

  public boolean hasValueStructSentinelNilCheck(TypeElement te) {
    return hasValueStructValidFlag(te) || getValueStructEnumSentinelFieldGoName(te) != null;
  }

  public boolean hasValueStructSentinelNilCheck(TypeMirror tm) {
    if (tm == null) {
      return false;
    }
    TypeMirror er = safeErasure(tm);
    if (er instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      return hasValueStructSentinelNilCheck(te);
    }
    return false;
  }

  public boolean isByteStringClass(TypeElement te) {
    return te != null && te.getQualifiedName().contentEquals("com.google.protobuf.ByteString");
  }

  public boolean isByteStringType(TypeMirror tm) {
    if (tm == null) {
      return false;
    }
    TypeMirror er = safeErasure(tm);
    if (er instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      return isByteStringClass(te);
    }
    return false;
  }

  public boolean isNativeGenericClass(TypeElement te) {
    return te != null
        && isSourceType(te)
        && (te.getKind() == ElementKind.CLASS || te.getKind() == ElementKind.INTERFACE)
        && !te.getTypeParameters().isEmpty()
        && !isNodeClass(te)
        && !isValueStructClass(te)
        && NATIVE_GENERIC_QUAL_NAMES.contains(te.getQualifiedName().toString());
  }

  public boolean isGenericStaticMethod(ExecutableElement ee) {
    return ee != null
        && ee.getKind() == ElementKind.METHOD
        && ee.getModifiers().contains(Modifier.STATIC)
        && !ee.getTypeParameters().isEmpty()
        && ee.getEnclosingElement() instanceof TypeElement owner
        && isNativeGenericClass(owner);
  }

  public List<TypeParameterElement> getEffectiveTypeParameters(TypeElement te) {
    if (te == null || !isSourceType(te) || enumConstantAnonTypes.contains(te)) {
      return List.of();
    }
    List<List<? extends TypeParameterElement>> chunks = new ArrayList<>();
    if (te.getNestingKind() == NestingKind.LOCAL
        || te.getNestingKind() == NestingKind.ANONYMOUS
        || hasOuterInstance(te)) {
      Element cur = te.getEnclosingElement();
      while (cur != null) {
        if (cur instanceof ExecutableElement ee) {
          if (ee.getModifiers().contains(Modifier.STATIC)) {
            if (isGenericStaticMethod(ee)) {
              chunks.add(0, ee.getTypeParameters());
            }
            break;
          }
        } else if (cur instanceof VariableElement ve) {
          if (ve.getModifiers().contains(Modifier.STATIC)) {
            break;
          }
        } else if (cur instanceof TypeElement encTe) {
          if (isNativeGenericClass(encTe)) {
            chunks.add(0, encTe.getTypeParameters());
          }
          if (!hasOuterInstance(encTe)) {
            break;
          }
        }
        cur = cur.getEnclosingElement();
      }
    }
    if (isNativeGenericClass(te)) {
      chunks.add(te.getTypeParameters());
    }
    if (chunks.isEmpty()) {
      return List.of();
    }
    Set<TypeParameterElement> out = new LinkedHashSet<>();
    for (List<? extends TypeParameterElement> chunk : chunks) {
      out.addAll(chunk);
    }
    return new ArrayList<>(out);
  }

  public boolean hasEffectiveTypeParameters(TypeElement te) {
    return !getEffectiveTypeParameters(te).isEmpty();
  }

  public void pushTypeParamScope(TypeElement ownerType, List<? extends TypeParameterElement> tpes) {
    Map<Element, String> map = new LinkedHashMap<>();
    Set<String> used = new HashSet<>();
    if (tpes != null) {
      for (TypeParameterElement tpe : tpes) {
        String base = GoNames.cleanChars(tpe.getSimpleName().toString());
        if (base.isEmpty()) {
          base = "T";
        }
        String cand = base;
        int idx = 2;
        while (!used.add(cand)) {
          cand = base + "_" + idx++;
        }
        map.put(tpe, cand);
      }
    }
    genericScopeStack.get().push(new GenericScope(ownerType, map));
  }

  public void popTypeParamScope() {
    Deque<GenericScope> stack = genericScopeStack.get();
    if (!stack.isEmpty()) {
      stack.pop();
    }
  }

  public String lookupActiveTypeParam(Element tpe) {
    if (tpe == null) {
      return null;
    }
    GenericScope scope = genericScopeStack.get().peek();
    return scope != null ? scope.typeParamNames().get(tpe) : null;
  }

  public boolean hasActiveTypeParameters() {
    GenericScope scope = genericScopeStack.get().peek();
    return scope != null && !scope.typeParamNames().isEmpty();
  }

  public boolean isActiveTypeVariable(TypeMirror tm) {
    return tm instanceof TypeVariable tv && lookupActiveTypeParam(tv.asElement()) != null;
  }

  public boolean isActiveComparableTypeVariable(TypeMirror tm) {
    return tm instanceof TypeVariable tv
        && tv.asElement() instanceof TypeParameterElement tpe
        && lookupActiveTypeParam(tpe) != null
        && "comparable".equals(getGoTypeParamConstraint(tpe));
  }

  public boolean isSameActiveComparableTypeVariable(TypeMirror a, TypeMirror b) {
    return a instanceof TypeVariable tva
        && b instanceof TypeVariable tvb
        && tva.asElement().equals(tvb.asElement())
        && isActiveComparableTypeVariable(a);
  }

  public TypeElement currentGenericScopeOwner() {
    GenericScope scope = genericScopeStack.get().peek();
    return scope != null ? scope.ownerType() : null;
  }

  private static final Set<String> FIRST_TYPE_PARAM_COMPARABLE_QUAL_NAMES =
      Set.of(
          "com.google.javascript.jscomp.base.LinkedIdentityHashMap",
          "com.google.javascript.jscomp.base.LinkedIdentityHashMap.Entry",
          "com.google.javascript.jscomp.base.LinkedIdentityHashSet",
          "com.google.javascript.jscomp.SimpleMultimap",
          "com.google.javascript.jscomp.graph.SubGraph",
          "com.google.javascript.jscomp.graph.AdjacencyGraph",
          "com.google.javascript.jscomp.graph.GraphNode",
          "com.google.javascript.jscomp.graph.Graph",
          "com.google.javascript.jscomp.graph.Graph.GraphEdge",
          "com.google.javascript.jscomp.graph.Graph.SimpleSubGraph",
          "com.google.javascript.jscomp.graph.DiGraph",
          "com.google.javascript.jscomp.graph.DiGraph.DiGraphNode",
          "com.google.javascript.jscomp.graph.DiGraph.DiGraphEdge",
          "com.google.javascript.jscomp.graph.LinkedDirectedGraph",
          "com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphNode",
          "com.google.javascript.jscomp.graph.LinkedDirectedGraph.LinkedDiGraphEdge",
          "com.google.javascript.jscomp.graph.UndiGraph",
          "com.google.javascript.jscomp.graph.UndiGraph.UndiGraphNode",
          "com.google.javascript.jscomp.graph.UndiGraph.UndiGraphEdge",
          "com.google.javascript.jscomp.graph.LinkedUndirectedGraph",
          "com.google.javascript.jscomp.graph.LinkedUndirectedGraph.LinkedUndirectedGraphNode",
          "com.google.javascript.jscomp.graph.LinkedUndirectedGraph.LinkedUndirectedGraphEdge",
          "com.google.javascript.jscomp.graph.DominatorTree",
          "com.google.javascript.jscomp.graph.FixedPointGraphTraversal",
          "com.google.javascript.jscomp.graph.FixedPointGraphTraversal.EdgeCallback",
          "com.google.javascript.jscomp.graph.GraphReachability",
          "com.google.javascript.jscomp.graph.GraphReachability.EdgeTuple",
          "com.google.javascript.jscomp.graph.CheckPathsBetweenNodes",
          "com.google.javascript.jscomp.graph.LowestCommonAncestorFinder",
          "com.google.javascript.jscomp.graph.LowestCommonAncestorFinder.Factory",
          "com.google.javascript.jscomp.graph.GraphColoring",
          "com.google.javascript.jscomp.graph.GraphColoring.GreedyGraphColoring",
          "com.google.javascript.jscomp.graph.UnionFind",
          "com.google.javascript.jscomp.graph.StandardUnionFind",
          "com.google.javascript.jscomp.graph.StandardUnionFind.Node",
          "com.google.javascript.jscomp.ControlFlowGraph",
          "com.google.javascript.jscomp.DataFlowAnalysis",
          "com.google.javascript.jscomp.DataFlowAnalysis.UniqueQueue",
          "com.google.javascript.rhino.HamtPMap.Holder");

  private static final Set<String> FIRST_TWO_TYPE_PARAMS_COMPARABLE_QUAL_NAMES =
      Set.of(
          "com.google.javascript.jscomp.SimpleTable",
          "com.google.javascript.jscomp.SimpleTable.SimpleCell",
          "com.google.javascript.rhino.PMap",
          "com.google.javascript.rhino.PMap.Reconciler",
          "com.google.javascript.rhino.HamtPMap",
          "com.google.javascript.rhino.HamtPMap.KeyIter",
          "com.google.javascript.rhino.HamtPMap.ValueIter");

  public String getGoTypeParamConstraint(TypeParameterElement tpe) {
    if (tpe == null) {
      return "any";
    }
    Element genericOwner = tpe.getGenericElement();
    TypeElement ownerTe = null;
    int idx = -1;
    if (genericOwner instanceof TypeElement te) {
      ownerTe = te;
      idx = te.getTypeParameters().indexOf(tpe);
    } else if (genericOwner instanceof ExecutableElement ee
        && ee.getEnclosingElement() instanceof TypeElement te) {
      ownerTe = te;
      idx = ee.getTypeParameters().indexOf(tpe);
    }
    if (ownerTe != null && idx >= 0) {
      String qn = ownerTe.getQualifiedName().toString();
      if (idx == 0 && FIRST_TYPE_PARAM_COMPARABLE_QUAL_NAMES.contains(qn)) {
        return "comparable";
      }
      if ((idx == 0 || idx == 1) && FIRST_TWO_TYPE_PARAMS_COMPARABLE_QUAL_NAMES.contains(qn)) {
        return "comparable";
      }
    }
    return "any";
  }

  public String formatGoTypeParamDecl(List<? extends TypeParameterElement> tpes) {
    if (tpes == null || tpes.isEmpty()) {
      return "";
    }
    List<String> parts = new ArrayList<>(tpes.size());
    for (TypeParameterElement tpe : tpes) {
      String name = lookupActiveTypeParam(tpe);
      if (name == null) {
        name = GoNames.cleanChars(tpe.getSimpleName().toString());
      }
      parts.add(name + " " + getGoTypeParamConstraint(tpe));
    }
    return "[" + String.join(", ", parts) + "]";
  }

  public String formatGoTypeParamUse(List<? extends TypeParameterElement> tpes) {
    if (tpes == null || tpes.isEmpty()) {
      return "";
    }
    List<String> parts = new ArrayList<>(tpes.size());
    for (TypeParameterElement tpe : tpes) {
      String name = lookupActiveTypeParam(tpe);
      parts.add(name != null ? name : "any");
    }
    return "[" + String.join(", ", parts) + "]";
  }

  public static boolean isAllUnboundedWildcards(DeclaredType dt) {
    if (dt == null || dt.getTypeArguments().isEmpty()) {
      return false;
    }
    for (TypeMirror ta : dt.getTypeArguments()) {
      if (!(ta instanceof WildcardType wt)
          || wt.getExtendsBound() != null
          || wt.getSuperBound() != null) {
        return false;
      }
    }
    return true;
  }

  public DeclaredType findSuperType(DeclaredType subType, TypeElement targetAnc) {
    if (subType == null || targetAnc == null) {
      return null;
    }
    if (targetAnc.equals(subType.asElement())) {
      return subType;
    }
    ArrayDeque<DeclaredType> queue = new ArrayDeque<>();
    Set<Element> visited = new HashSet<>();
    queue.add(subType);
    visited.add(subType.asElement());
    while (!queue.isEmpty()) {
      DeclaredType cur = queue.poll();
      for (TypeMirror supTm : types.directSupertypes(cur)) {
        if (supTm instanceof DeclaredType supDt) {
          Element supEl = supDt.asElement();
          if (targetAnc.equals(supEl)) {
            return supDt;
          }
          if (visited.add(supEl)) {
            queue.add(supDt);
          }
        }
      }
    }
    return null;
  }

  public String formatDeclaredTypeArgs(
      TypeMirror tm, TypeElement expectedOwner, PackagePlanner.ImportTracker imports) {
    if (expectedOwner == null || !hasEffectiveTypeParameters(expectedOwner)) {
      return "";
    }
    List<TypeParameterElement> effParams = getEffectiveTypeParameters(expectedOwner);
    TypeMirror unwrapped = tm;
    if (unwrapped instanceof WildcardType wt && wt.getExtendsBound() != null) {
      unwrapped = wt.getExtendsBound();
    }
    if (unwrapped instanceof TypeVariable tv && tv.getUpperBound() != null) {
      unwrapped = tv.getUpperBound();
    }
    if (unwrapped instanceof DeclaredType udt && !expectedOwner.equals(udt.asElement())) {
      DeclaredType sup = findSuperType(udt, expectedOwner);
      if (sup != null) {
        unwrapped = sup;
      }
    }
    DeclaredType matchedDt =
        (unwrapped instanceof DeclaredType dt && expectedOwner.equals(dt.asElement())) ? dt : null;
    if (matchedDt != null
        && isAllUnboundedWildcards(matchedDt)
        && expectedOwner.equals(currentGenericScopeOwner())
        && effParams.size() == expectedOwner.getTypeParameters().size()) {
      return formatGoTypeParamUse(effParams);
    }
    List<String> args = new ArrayList<>(effParams.size());
    for (TypeParameterElement tpe : effParams) {
      TypeMirror resolvedArg = null;
      boolean matchedRawOwner = false;
      if (tpe.getGenericElement() instanceof TypeElement paramOwnerTe && matchedDt != null) {
        int idx = paramOwnerTe.getTypeParameters().indexOf(tpe);
        if (idx >= 0) {
          TypeMirror cur = matchedDt;
          while (cur instanceof DeclaredType curDt) {
            if (paramOwnerTe.equals(curDt.asElement())) {
              if (idx < curDt.getTypeArguments().size()) {
                resolvedArg = curDt.getTypeArguments().get(idx);
              } else {
                matchedRawOwner = true;
              }
              break;
            }
            cur = curDt.getEnclosingType();
          }
        }
      }
      if (resolvedArg != null) {
        if (resolvedArg instanceof WildcardType wt
            && wt.getExtendsBound() == null
            && wt.getSuperBound() == null
            && tpe.getGenericElement().equals(currentGenericScopeOwner())) {
          String inScope = lookupActiveTypeParam(tpe);
          args.add(inScope != null ? inScope : "any");
        } else {
          String g = goType(resolvedArg, imports);
          args.add(g.isEmpty() ? "any" : g);
        }
      } else if (matchedRawOwner) {
        args.add("any");
      } else {
        String inScope = lookupActiveTypeParam(tpe);
        args.add(inScope != null ? inScope : "any");
      }
    }
    return "[" + String.join(", ", args) + "]";
  }

  public String formatSuperTypeArgs(
      TypeMirror superOrIfaceTm, TypeElement superOrIfaceTe, PackagePlanner.ImportTracker imports) {
    if (superOrIfaceTe == null
        || !isSourceType(superOrIfaceTe)
        || !hasEffectiveTypeParameters(superOrIfaceTe)) {
      return "";
    }
    if (superOrIfaceTm != null) {
      return formatDeclaredTypeArgs(superOrIfaceTm, superOrIfaceTe, imports);
    }
    return formatGoTypeParamUse(getEffectiveTypeParameters(superOrIfaceTe));
  }

  public String goEmbeddedSuperStructType(
      TypeMirror superclassTm, TypeElement superTe, PackagePlanner.ImportTracker imports) {
    String base = goStructTypeForElement(superTe, imports);
    return base + formatSuperTypeArgs(superclassTm, superTe, imports);
  }

  public TypeMirror substituteTypeParams(
      TypeMirror tm,
      List<? extends TypeParameterElement> params,
      List<? extends TypeMirror> args) {
    if (tm == null || params == null || params.isEmpty()) {
      return tm;
    }
    return switch (tm.getKind()) {
      case TYPEVAR -> {
        TypeVariable tv = (TypeVariable) tm;
        Element el = tv.asElement();
        for (int i = 0; i < params.size() && i < args.size(); i++) {
          if (params.get(i).equals(el)) {
            yield args.get(i);
          }
        }
        yield tm;
      }
      case WILDCARD -> {
        WildcardType wt = (WildcardType) tm;
        if (wt.getExtendsBound() != null) {
          yield substituteTypeParams(wt.getExtendsBound(), params, args);
        }
        yield objectElement.asType();
      }
      case ARRAY -> {
        ArrayType at = (ArrayType) tm;
        TypeMirror subComp = substituteTypeParams(at.getComponentType(), params, args);
        yield subComp == at.getComponentType() ? at : types.getArrayType(subComp);
      }
      case DECLARED -> {
        DeclaredType dt = (DeclaredType) tm;
        List<? extends TypeMirror> tas = dt.getTypeArguments();
        if (tas.isEmpty() || !(dt.asElement() instanceof TypeElement dte)) {
          yield dt;
        }
        if (dte.getTypeParameters().size() != tas.size()) {
          yield dt;
        }
        TypeMirror[] subTas = new TypeMirror[tas.size()];
        boolean changed = false;
        for (int i = 0; i < tas.size(); i++) {
          TypeMirror orig = tas.get(i);
          TypeMirror sub = substituteTypeParams(orig, params, args);
          if (sub.getKind() == TypeKind.WILDCARD) {
            WildcardType wt = (WildcardType) sub;
            sub = wt.getExtendsBound() != null ? wt.getExtendsBound() : objectElement.asType();
          }
          subTas[i] = sub;
          if (sub != orig) {
            changed = true;
          }
        }
        yield changed ? types.getDeclaredType(dte, subTas) : dt;
      }
      default -> tm;
    };
  }

  public TypeMirror substituteOwnerTypeArgs(
      TypeMirror memberTm, TypeElement owner, TypeMirror receiverOrSiteTm) {
    if (memberTm == null || owner == null || !hasEffectiveTypeParameters(owner)) {
      return memberTm;
    }
    if (!isNativeGenericClass(owner)) {
      List<TypeParameterElement> effTpes = getEffectiveTypeParameters(owner);
      List<TypeMirror> effArgs = new ArrayList<>(effTpes.size());
      TypeMirror curSite = receiverOrSiteTm;
      if (curSite instanceof WildcardType wt && wt.getExtendsBound() != null) {
        curSite = wt.getExtendsBound();
      }
      if (curSite instanceof TypeVariable tv && tv.getUpperBound() != null) {
        curSite = tv.getUpperBound();
      }
      for (TypeParameterElement tpe : effTpes) {
        if (lookupActiveTypeParam(tpe) != null) {
          effArgs.add(tpe.asType());
        } else if (tpe.getGenericElement() instanceof TypeElement encOwner
            && isNativeGenericClass(encOwner)
            && curSite instanceof DeclaredType cdt) {
          TypeMirror encSub = substituteOwnerTypeArgs(tpe.asType(), encOwner, cdt);
          effArgs.add(encSub != null ? encSub : objectElement.asType());
        } else {
          effArgs.add(objectElement.asType());
        }
      }
      return substituteTypeParams(memberTm, effTpes, effArgs);
    }
    List<? extends TypeParameterElement> tpes = owner.getTypeParameters();
    int n = tpes.size();
    List<TypeMirror> actualArgs = new ArrayList<>(n);
    TypeMirror unwrapped = receiverOrSiteTm;
    if (unwrapped instanceof WildcardType wt && wt.getExtendsBound() != null) {
      unwrapped = wt.getExtendsBound();
    }
    if (unwrapped instanceof TypeVariable tv && tv.getUpperBound() != null) {
      unwrapped = tv.getUpperBound();
    }
    if (unwrapped instanceof DeclaredType udt && !owner.equals(udt.asElement())) {
      DeclaredType sup = findSuperType(udt, owner);
      if (sup != null) {
        unwrapped = sup;
      }
    }
    if (unwrapped instanceof DeclaredType dt
        && owner.equals(dt.asElement())
        && dt.getTypeArguments().size() == n) {
      if (isAllUnboundedWildcards(dt) && owner.equals(currentGenericScopeOwner())) {
        for (TypeParameterElement tpe : tpes) {
          actualArgs.add(tpe.asType());
        }
      } else {
        for (TypeMirror ta : dt.getTypeArguments()) {
          if (ta instanceof WildcardType wt) {
            actualArgs.add(
                wt.getExtendsBound() != null ? wt.getExtendsBound() : objectElement.asType());
          } else {
            actualArgs.add(ta);
          }
        }
      }
    } else {
      for (int i = 0; i < n; i++) {
        actualArgs.add(objectElement.asType());
      }
    }
    return substituteTypeParams(memberTm, tpes, actualArgs);
  }

  public TypeMirror resolveSigParamType(MethodSig sig, int idx, TypeMirror siteTm) {
    TypeMirror base = sig.canonParamTypes.get(idx);
    if (sig.paramWidenedToAny[idx] || sig.canonOwner == null || siteTm == null) {
      return goEffectiveType(base);
    }
    return goEffectiveType(substituteOwnerTypeArgs(base, sig.canonOwner, siteTm));
  }

  public TypeMirror resolveSigReturnType(MethodSig sig, TypeMirror siteTm) {
    if (sig.canonReturnType == null) {
      return null;
    }
    if (sig.returnWidenedToAny || sig.canonOwner == null || siteTm == null) {
      return goEffectiveType(sig.canonReturnType);
    }
    return goEffectiveType(substituteOwnerTypeArgs(sig.canonReturnType, sig.canonOwner, siteTm));
  }

  public List<TypeMirror> inferGenericStaticMethodTypeArgs(
      ExecutableElement ee,
      List<? extends TypeMirror> explicitAstTypeArgs,
      List<? extends TypeMirror> actualValueArgTypes,
      TypeMirror callSiteReturnType) {
    List<? extends TypeParameterElement> tpes = ee.getTypeParameters();
    int n = tpes.size();
    if (explicitAstTypeArgs != null && explicitAstTypeArgs.size() == n) {
      return new ArrayList<>(explicitAstTypeArgs);
    }
    Map<Element, TypeMirror> inferred = new HashMap<>();
    if (actualValueArgTypes != null) {
      int fixed =
          ee.isVarArgs() && !ee.getParameters().isEmpty()
              ? ee.getParameters().size() - 1
              : ee.getParameters().size();
      for (int i = 0; i < fixed && i < actualValueArgTypes.size(); i++) {
        unifyGenericTypePattern(
            ee.getParameters().get(i).asType(), actualValueArgTypes.get(i), inferred);
      }
    }
    if (callSiteReturnType != null) {
      unifyGenericTypePattern(ee.getReturnType(), callSiteReturnType, inferred);
    }
    List<TypeMirror> out = new ArrayList<>(n);
    for (TypeParameterElement tpe : tpes) {
      TypeMirror m = inferred.get(tpe);
      if (m == null
          || m.getKind() == TypeKind.VOID
          || m.getKind() == TypeKind.NONE
          || m.getKind() == TypeKind.NULL) {
        m = objectElement.asType();
      }
      out.add(m);
    }
    return out;
  }

  private void unifyGenericTypePattern(
      TypeMirror pattern, TypeMirror actual, Map<Element, TypeMirror> out) {
    if (pattern == null
        || actual == null
        || actual.getKind() == TypeKind.VOID
        || actual.getKind() == TypeKind.NONE
        || actual.getKind() == TypeKind.NULL) {
      return;
    }
    if (pattern instanceof TypeVariable tv) {
      out.putIfAbsent(tv.asElement(), actual);
      return;
    }
    if (pattern instanceof WildcardType wt && wt.getExtendsBound() != null) {
      if (actual instanceof WildcardType awt && awt.getExtendsBound() != null) {
        unifyGenericTypePattern(wt.getExtendsBound(), awt.getExtendsBound(), out);
      } else {
        unifyGenericTypePattern(wt.getExtendsBound(), actual, out);
      }
      return;
    }
    if (pattern instanceof ArrayType pat && actual instanceof ArrayType aat) {
      unifyGenericTypePattern(pat.getComponentType(), aat.getComponentType(), out);
      return;
    }
    if (pattern instanceof DeclaredType pdt
        && actual instanceof DeclaredType adt
        && pdt.asElement().equals(adt.asElement())
        && !adt.getTypeArguments().isEmpty()
        && pdt.getTypeArguments().size() == adt.getTypeArguments().size()) {
      for (int i = 0; i < pdt.getTypeArguments().size(); i++) {
        unifyGenericTypePattern(
            pdt.getTypeArguments().get(i), adt.getTypeArguments().get(i), out);
      }
    }
  }

  public boolean isValueEnum(TypeElement te) {
    return te != null && valueEnumTypes.contains(te);
  }

  public Map<VariableElement, List<Object>> getValueEnumFieldTable(TypeElement te) {
    return valueEnumFieldTables.getOrDefault(te, Map.of());
  }

  public boolean canDispatchDirectlyOnBase(TypeElement curType, MethodSig sig) {
    if (curType == null || sig == null) {
      return false;
    }
    Set<MethodSig> direct = directBaseDispatchSigs.get(curType);
    return direct != null && direct.contains(sig);
  }

  public boolean mayHaveSelfField(TypeElement te) {
    return te != null && mayHaveSelfFieldTypes.contains(te);
  }

  public boolean isReusableSelfAncestor(TypeElement te) {
    return te != null && reusableSelfAncestors.contains(te);
  }

  public TypeElement getSamePkgReusableSelfRoot(TypeElement curType) {
    if (curType == null
        || !isSubclassed(curType)
        || !isSourceType(curType)
        || curType.getKind() != ElementKind.CLASS
        || isNodeClass(curType)
        || isValueStructClass(curType)) {
      return null;
    }
    String curPkg = planner.getGoPkg(curType);
    TypeElement root = null;
    TypeElement sup = getSuperclassElement(curType);
    while (sup != null
        && isSourceType(sup)
        && !runtime.mapsToAny(sup)
        && curPkg.equals(planner.getGoPkg(sup))) {
      if (reusableSelfAncestors.contains(sup)) {
        root = sup;
      }
      sup = getSuperclassElement(sup);
    }
    return root;
  }

  private List<TypeElement> collectReusableSelfAncestorsRootFirst(TypeElement curType) {
    String curPkg = planner.getGoPkg(curType);
    List<TypeElement> chain = new ArrayList<>();
    TypeElement sup = getSuperclassElement(curType);
    while (sup != null
        && isSourceType(sup)
        && !runtime.mapsToAny(sup)
        && curPkg.equals(planner.getGoPkg(sup))) {
      if (reusableSelfAncestors.contains(sup)) {
        chain.add(sup);
      }
      sup = getSuperclassElement(sup);
    }
    Collections.reverse(chain);
    return chain;
  }

  public TypeElement findReusableSelfAncestor(
      TypeElement curType, TypeMirror targetType, boolean forceAny) {
    if (curType == null
        || !isSubclassed(curType)
        || !isSourceType(curType)
        || curType.getKind() != ElementKind.CLASS
        || isNodeClass(curType)
        || isValueStructClass(curType)
        || !(curType.asType() instanceof DeclaredType curDt)) {
      return null;
    }
    List<TypeElement> candidates = collectReusableSelfAncestorsRootFirst(curType);
    if (candidates.isEmpty()) {
      return null;
    }
    if (forceAny || isObjectOrErasedAny(targetType)) {
      return candidates.get(0);
    }
    if (targetType == null) {
      return null;
    }
    TypeMirror dstEff = goEffectiveType(targetType);
    if (!(dstEff instanceof DeclaredType dstDt)
        || !(dstDt.asElement() instanceof TypeElement dstTe)
        || isActiveTypeVariable(dstEff)) {
      return null;
    }
    String dstKey = canonicalGoTypeKey(dstEff);
    for (TypeElement anc : candidates) {
      if (anc.equals(dstTe)) {
        DeclaredType ancSite = findSuperType(curDt, anc);
        if (ancSite != null
            && canonicalGoTypeKey(goEffectiveType(ancSite)).equals(dstKey)) {
          return anc;
        }
      }
    }
    return null;
  }

  public TypeElement findReusableSelfAncestorForCall(
      TypeElement curType, TypeElement owner, MethodSig sig) {
    if (curType == null
        || owner == null
        || sig == null
        || !isSubclassed(curType)
        || !isSourceType(curType)
        || curType.getKind() != ElementKind.CLASS
        || isNodeClass(curType)
        || isValueStructClass(curType)
        || runtime.mapsToAny(owner)
        || runtime.isShimException(owner)
        || !(isSourceType(owner) || isNeededExternalSig(sig))
        || !(curType.asType() instanceof DeclaredType curDt)) {
      return null;
    }
    List<TypeElement> candidates = collectReusableSelfAncestorsRootFirst(curType);
    if (candidates.isEmpty()) {
      return null;
    }
    TypeMirror ownerEr = safeErasure(owner.asType());
    for (TypeElement anc : candidates) {
      if (types.isSubtype(safeErasure(anc.asType()), ownerEr)) {
        DeclaredType ancSite = findSuperType(curDt, anc);
        if (ancSite != null && sameResolvedMethodSigAtSites(sig, curDt, ancSite)) {
          return anc;
        }
      }
    }
    return null;
  }

  private boolean sameResolvedMethodSigAtSites(
      MethodSig sig, TypeMirror siteA, TypeMirror siteB) {
    for (int i = 0; i < sig.canonParamTypes.size(); i++) {
      TypeMirror pa = resolveSigParamType(sig, i, siteA);
      TypeMirror pb = resolveSigParamType(sig, i, siteB);
      if (!canonicalGoTypeKey(pa).equals(canonicalGoTypeKey(pb))) {
        return false;
      }
    }
    if (sig.canonReturnType != null) {
      TypeMirror ra = resolveSigReturnType(sig, siteA);
      TypeMirror rb = resolveSigReturnType(sig, siteB);
      if (!canonicalGoTypeKey(ra).equals(canonicalGoTypeKey(rb))) {
        return false;
      }
    }
    return true;
  }

  public boolean needsBindSelf(TypeElement te) {
    if (te == null || isValueStructClass(te) || isNodeClass(te) || isValueEnum(te)) {
      return false;
    }
    if (isSourceType(te)) {
      if (isSubclassed(te) && mayHaveSelfField(te)) {
        return true;
      }
      TypeElement sup = getSuperclassElement(te);
      if (sup == null || runtime.mapsToAny(sup)) {
        return false;
      }
      if (isSourceType(sup)) {
        return needsBindSelf(sup);
      }
      return !runtime.isShimException(sup) && runtime.isSubclassedExternal(sup);
    }
    return !runtime.isShimException(te) && runtime.isSubclassedExternal(te);
  }

  public boolean isValueEnumType(TypeMirror tm) {
    if (tm == null) {
      return false;
    }
    TypeMirror er = safeErasure(tm);
    if (er instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      return isValueEnum(te);
    }
    return false;
  }

  public boolean hasFastStaticSyncMutex(TypeElement te) {
    return te != null && fastStaticSyncTypes.contains(te);
  }

  public boolean hasFastInstanceSyncMutex(TypeElement te) {
    return te != null && fastInstanceSyncTypes.contains(te);
  }

  public boolean isStringType(TypeMirror tm) {
    if (tm == null) {
      return false;
    }
    TypeMirror er = safeErasure(tm);
    if (er instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      return te.getQualifiedName().contentEquals("java.lang.String");
    }
    return false;
  }

  public boolean isObjectOrErasedAny(TypeMirror tm) {
    TypeMirror eff = goEffectiveType(tm);
    if (eff == null || isActiveTypeVariable(eff)) {
      return false;
    }
    TypeMirror er = safeErasure(eff);
    if (er instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      return runtime.mapsToAny(te);
    }
    return false;
  }

  public boolean isThrowable(TypeElement te) {
    if (te == null || throwableElement == null) {
      return false;
    }
    return types.isSubtype(safeErasure(te.asType()), safeErasure(throwableElement.asType()));
  }

  public boolean isThrowableType(TypeMirror tm) {
    if (tm == null || throwableElement == null) {
      return false;
    }
    if (tm instanceof UnionType) {
      return true;
    }
    TypeMirror er = safeErasure(tm);
    if (er.getKind() == TypeKind.DECLARED) {
      return types.isSubtype(er, safeErasure(throwableElement.asType()));
    }
    return false;
  }

  public boolean isSubclassed(TypeElement te) {
    if (isNodeClass(te) || isValueStructClass(te)) {
      return false;
    }
    if (te.getKind() == ElementKind.ENUM
        || te.getKind() == ElementKind.RECORD
        || te.getKind() == ElementKind.INTERFACE
        || te.getKind() == ElementKind.ANNOTATION_TYPE) {
      return false;
    }
    if (isSourceType(te)) {
      return subclassedTypes.contains(te) || te.getModifiers().contains(Modifier.ABSTRACT);
    }
    return runtime.isSubclassedExternal(te);
  }

  public boolean isGoInterfaceType(TypeMirror tm) {
    if (tm == null) {
      return false;
    }
    TypeMirror er = safeErasure(tm);
    if (er instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      if (isNodeClass(te)
          || isValueStructClass(te)
          || te.getQualifiedName().contentEquals("java.lang.String")) {
        return false;
      }
      if (runtime.mapsToAny(te)) {
        return true;
      }
      if (!isSourceType(te)) {
        return runtime.isGoInterfaceExternal(te);
      }
      return te.getKind() == ElementKind.INTERFACE || isSubclassed(te);
    }
    return false;
  }

  public boolean isEnumConstantAnonType(TypeElement te) {
    return enumConstantAnonTypes.contains(te);
  }

  public TypeElement getEnumConstantBodyType(VariableElement enumConst) {
    return enumConstantBodyTypes.get(enumConst);
  }

  public TypeElement getClassTreeElement(ClassTree ct) {
    return classTreeElements.get(ct);
  }

  public List<VariableElement> getCapturedVars(TypeElement localOrAnon) {
    return capturedVars.getOrDefault(localOrAnon, List.of());
  }

  /** True if {@code inner} is a non-static inner class holding an implicit {@code outer_} reference. */
  public boolean hasOuterInstance(TypeElement te) {
    if (te == null
        || isNodeClass(te)
        || enumConstantAnonTypes.contains(te)) {
      return false;
    }
    if (te.getKind() == ElementKind.INTERFACE
        || te.getKind() == ElementKind.ENUM
        || te.getKind() == ElementKind.RECORD
        || te.getKind() == ElementKind.ANNOTATION_TYPE) {
      return false;
    }
    if (te.getModifiers().contains(Modifier.STATIC)) {
      return false;
    }
    if (te.getNestingKind() == NestingKind.MEMBER) {
      Element enc = te.getEnclosingElement();
      // Members of interfaces/records are implicitly static in Java
      if (enc instanceof TypeElement encTe) {
        return encTe.getKind() != ElementKind.INTERFACE && encTe.getKind() != ElementKind.RECORD;
      }
      return false;
    }
    if (te.getNestingKind() == NestingKind.LOCAL || te.getNestingKind() == NestingKind.ANONYMOUS) {
      // Check if declared in a static context (static method, static field, static init, interface/record)
      Element cur = te.getEnclosingElement();
      while (cur != null) {
        if (cur instanceof ExecutableElement ee) {
          return !ee.getModifiers().contains(Modifier.STATIC);
        }
        if (cur instanceof VariableElement ve) {
          return !ve.getModifiers().contains(Modifier.STATIC);
        }
        if (cur instanceof TypeElement encTe) {
          // Instance initializer or instance field without VariableElement parent
          return encTe.getKind() != ElementKind.INTERFACE;
        }
        cur = cur.getEnclosingElement();
      }
    }
    return false;
  }

  public TypeElement getEnclosingInstanceType(TypeElement inner) {
    if (!hasOuterInstance(inner)) {
      return null;
    }
    Element cur = inner.getEnclosingElement();
    while (cur != null) {
      if (cur instanceof TypeElement te && !enumConstantAnonTypes.contains(te)) {
        return te;
      }
      cur = cur.getEnclosingElement();
    }
    return null;
  }

  public Set<TypeElement> getReferencedExternalTypes() {
    return Collections.unmodifiableSet(referencedExternalTypes);
  }

  public Set<ExecutableElement> getReferencedExternalExecs() {
    return Collections.unmodifiableSet(referencedExternalExecs);
  }

  public Set<VariableElement> getReferencedExternalFields() {
    return Collections.unmodifiableSet(referencedExternalFields);
  }

  public void recordExternalField(VariableElement ve) {
    referencedExternalFields.add(ve);
    recordTypeMirror(ve.asType());
  }

  public Set<TypeElement> getLambdaTargetInterfaces() {
    return Collections.unmodifiableSet(lambdaTargetInterfaces);
  }

  public boolean isNeededExternalSig(MethodSig sig) {
    return sig != null && neededExternalSigs.contains(sig);
  }

  public boolean isFieldCrossPkgRead(VariableElement ve) {
    return crossPkgReadFields.contains(ve);
  }

  public boolean isFieldCrossPkgWrite(VariableElement ve) {
    return crossPkgWriteFields.contains(ve);
  }

  public String getFieldGoName(VariableElement ve) {
    String n = fieldGoNames.get(ve);
    if (n != null) {
      return n;
    }
    return GoNames.unexported(ve.getSimpleName().toString());
  }

  public String getConstructorGoName(ExecutableElement ctor) {
    String n = constructorGoNames.get(ctor);
    if (n != null) {
      return n;
    }
    TypeElement te = (TypeElement) ctor.getEnclosingElement();
    return "New" + getTypeGoBaseName(te);
  }

  public String getStaticMethodGoName(ExecutableElement m) {
    String n = staticMethodGoNames.get(m);
    if (n != null) {
      return n;
    }
    registerExternalExec(m);
    n = staticMethodGoNames.get(m);
    if (n != null) {
      return n;
    }
    TypeElement te = (TypeElement) m.getEnclosingElement();
    return staticMethodPrefix(te, GoNames.exported(m.getSimpleName().toString()));
  }

  public MethodSig getInstanceMethodSig(ExecutableElement m) {
    MethodSig sig = instanceMethodSigs.get(m);
    if (sig != null) {
      return sig;
    }
    // Fallback for late-discovered external methods
    registerExternalExec(m);
    return instanceMethodSigs.get(m);
  }

  public String getTypeGoBaseName(TypeElement te) {
    if (isSourceType(te)) {
      return planner.getGoTypeName(te);
    }
    return runtime.extGoTypeName(te);
  }

  /** Marker method name on interfaces and subclassed class interfaces so empty interfaces never match arbitrary types. */
  public String getInterfaceMarkerMethod(TypeElement te) {
    if (isSourceType(te)) {
      return "Implements_" + planner.getGoPkg(te) + "_" + planner.getGoTypeName(te);
    }
    return "Implements_j2gort_" + runtime.extGoTypeName(te);
  }

  /** Exported parameterless marker method on generic source structs for wildcard `instanceof` and reflection. */
  public String getGenericClassMarkerMethod(TypeElement te) {
    return "IsGenericInst_" + planner.getGoPkg(te) + "_" + planner.getGoTypeName(te) + "_";
  }

  public TypeMirror safeErasure(TypeMirror tm) {
    if (tm == null) {
      return null;
    }
    return switch (tm.getKind()) {
      case VOID, NONE, NULL, ERROR, PACKAGE, EXECUTABLE, OTHER -> tm;
      case BOOLEAN, BYTE, SHORT, INT, LONG, CHAR, FLOAT, DOUBLE -> tm;
      case UNION -> {
        UnionType ut = (UnionType) tm;
        yield throwableElement != null ? throwableElement.asType() : ut.getAlternatives().get(0);
      }
      case INTERSECTION -> {
        IntersectionType it = (IntersectionType) tm;
        yield safeErasure(it.getBounds().get(0));
      }
      case WILDCARD -> {
        WildcardType wt = (WildcardType) tm;
        if (wt.getExtendsBound() != null) {
          yield safeErasure(wt.getExtendsBound());
        }
        yield objectElement.asType();
      }
      case TYPEVAR -> {
        TypeVariable tv = (TypeVariable) tm;
        TypeMirror ub = tv.getUpperBound();
        yield ub != null ? safeErasure(ub) : objectElement.asType();
      }
      default -> types.erasure(tm);
    };
  }

  public TypeMirror goEffectiveType(TypeMirror tm) {
    if (tm == null) {
      return null;
    }
    return switch (tm.getKind()) {
      case VOID, NONE, NULL, ERROR, PACKAGE, EXECUTABLE, OTHER -> tm;
      case BOOLEAN, BYTE, SHORT, INT, LONG, CHAR, FLOAT, DOUBLE -> tm;
      case UNION -> {
        UnionType ut = (UnionType) tm;
        yield throwableElement != null ? throwableElement.asType() : ut.getAlternatives().get(0);
      }
      case INTERSECTION -> {
        IntersectionType it = (IntersectionType) tm;
        yield goEffectiveType(it.getBounds().get(0));
      }
      case WILDCARD -> {
        WildcardType wt = (WildcardType) tm;
        if (wt.getExtendsBound() != null) {
          yield goEffectiveType(wt.getExtendsBound());
        }
        yield objectElement.asType();
      }
      case TYPEVAR -> {
        TypeVariable tv = (TypeVariable) tm;
        if (isActiveTypeVariable(tv)) {
          yield tv;
        }
        TypeMirror ub = tv.getUpperBound();
        yield ub != null ? goEffectiveType(ub) : objectElement.asType();
      }
      case ARRAY -> {
        ArrayType at = (ArrayType) tm;
        yield types.getArrayType(goEffectiveType(at.getComponentType()));
      }
      case DECLARED -> {
        DeclaredType dt = (DeclaredType) tm;
        if (dt.asElement() instanceof TypeElement te) {
          if (isNativeGenericClass(te)) {
            int n = te.getTypeParameters().size();
            TypeMirror[] normArgs = new TypeMirror[n];
            if (dt.getTypeArguments().size() == n) {
              if (isAllUnboundedWildcards(dt) && te.equals(currentGenericScopeOwner())) {
                for (int i = 0; i < n; i++) {
                  normArgs[i] = te.getTypeParameters().get(i).asType();
                }
              } else {
                for (int i = 0; i < n; i++) {
                  TypeMirror eff = goEffectiveType(dt.getTypeArguments().get(i));
                  normArgs[i] =
                      (eff == null
                              || eff.getKind() == TypeKind.VOID
                              || eff.getKind() == TypeKind.NONE
                              || eff.getKind() == TypeKind.NULL)
                          ? objectElement.asType()
                          : eff;
                }
              }
            } else {
              for (int i = 0; i < n; i++) {
                normArgs[i] = objectElement.asType();
              }
            }
            DeclaredType normEnc =
                (hasOuterInstance(te) && dt.getEnclosingType() instanceof DeclaredType encDt)
                    ? (DeclaredType) goEffectiveType(encDt)
                    : null;
            yield normEnc != null
                ? types.getDeclaredType(normEnc, te, normArgs)
                : types.getDeclaredType(te, normArgs);
          }
          if (hasEffectiveTypeParameters(te)) {
            if (hasOuterInstance(te) && dt.getEnclosingType() instanceof DeclaredType encDt) {
              yield types.getDeclaredType((DeclaredType) goEffectiveType(encDt), te);
            }
            yield dt;
          }
        }
        yield types.erasure(tm);
      }
      default -> types.erasure(tm);
    };
  }

  /** Renders a Java {@link TypeMirror} as a Go type in the context of {@code imports}. */
  public String goType(TypeMirror tm, PackagePlanner.ImportTracker imports) {
    if (tm == null || tm.getKind() == TypeKind.VOID || tm.getKind() == TypeKind.NONE) {
      return "";
    }
    if (tm.getKind() == TypeKind.NULL) {
      return "any";
    }
    TypeMirror er = goEffectiveType(tm);
    return switch (er.getKind()) {
      case BOOLEAN -> "bool";
      case BYTE -> "int8";
      case SHORT -> "int16";
      case CHAR -> "uint16";
      case INT -> "int32";
      case LONG -> "int64";
      case FLOAT -> "float32";
      case DOUBLE -> "float64";
      case ARRAY -> "[]" + goType(((ArrayType) er).getComponentType(), imports);
      case TYPEVAR -> {
        String tvName = lookupActiveTypeParam(((TypeVariable) er).asElement());
        yield tvName != null ? tvName : "any";
      }
      case DECLARED -> {
        DeclaredType dt = (DeclaredType) er;
        TypeElement te = (TypeElement) dt.asElement();
        if (hasEffectiveTypeParameters(te)) {
          yield goTypeForElementBase(te, imports) + formatDeclaredTypeArgs(dt, te, imports);
        }
        yield goTypeForElement(te, imports);
      }
      default -> "any";
    };
  }

  public String goTypeForElement(TypeElement te, PackagePlanner.ImportTracker imports) {
    String base = goTypeForElementBase(te, imports);
    if (te != null && !base.equals("any") && hasEffectiveTypeParameters(te)) {
      return base + formatGoTypeParamUse(getEffectiveTypeParameters(te));
    }
    return base;
  }

  public String goTypeForElementBase(TypeElement te, PackagePlanner.ImportTracker imports) {
    if (te == null) {
      return "any";
    }
    if (isSourceType(te)) {
      if (te.getKind() == ElementKind.ANNOTATION_TYPE) {
        return "any";
      }
      if (enumConstantAnonTypes.contains(te)) {
        TypeElement sup = getSuperclassElement(te);
        if (sup != null) {
          return goTypeForElementBase(sup, imports);
        }
      }
      String targetPkg = planner.getGoPkg(te);
      String prefix = "";
      String alias = imports.ext(planner.getGoImportPath(targetPkg), targetPkg);
      if (!alias.isEmpty()) {
        prefix = alias + ".";
      }
      if (isNodeClass(te)) {
        return prefix + "Node";
      }
      String goName = planner.getGoTypeName(te);
      if (te.getKind() == ElementKind.INTERFACE
          || isSubclassed(te)
          || isValueStructClass(te)
          || isValueEnum(te)) {
        return prefix + goName;
      }
      return "*" + prefix + goName;
    }
    recordExternalType(te);
    return runtime.goType(te, imports);
  }

  /** Package-qualified concrete struct name (without leading `*`) for `new T(...)` or base embedding. */
  public String goStructTypeForElement(TypeElement te, PackagePlanner.ImportTracker imports) {
    if (isSourceType(te)) {
      String targetPkg = planner.getGoPkg(te);
      String alias = imports.ext(planner.getGoImportPath(targetPkg), targetPkg);
      String prefix = alias.isEmpty() ? "" : (alias + ".");
      String goName = planner.getGoTypeName(te);
      return isSubclassed(te) ? (prefix + goName + "Base") : (prefix + goName);
    }
    recordExternalType(te);
    return runtime.goStructType(te, imports);
  }

  /** Canonical package-independent Go type key used for Union-Find and overload suffixes. */
  public String canonicalGoTypeKey(TypeMirror tm) {
    if (tm == null || tm.getKind() == TypeKind.VOID || tm.getKind() == TypeKind.NONE) {
      return "void";
    }
    TypeMirror er = goEffectiveType(tm);
    return switch (er.getKind()) {
      case BOOLEAN -> "bool";
      case BYTE -> "int8";
      case SHORT -> "int16";
      case CHAR -> "uint16";
      case INT -> "int32";
      case LONG -> "int64";
      case FLOAT -> "float32";
      case DOUBLE -> "float64";
      case ARRAY -> "slice_" + canonicalGoTypeKey(((ArrayType) er).getComponentType());
      case TYPEVAR -> {
        Element tvEl = ((TypeVariable) er).asElement();
        String activeName = lookupActiveTypeParam(tvEl);
        yield "typevar_" + (activeName != null ? activeName : tvEl.getSimpleName());
      }
      case DECLARED -> {
        DeclaredType dt = (DeclaredType) er;
        TypeElement te = (TypeElement) dt.asElement();
        if (isNodeClass(te)) {
          yield "rhino.Node";
        }
        if (isSourceType(te)) {
          if (te.getKind() == ElementKind.ANNOTATION_TYPE) {
            yield "any";
          }
          String baseKey = planner.getGoPkg(te) + "." + planner.getGoTypeName(te);
          if (isNativeGenericClass(te) && !dt.getTypeArguments().isEmpty()) {
            List<String> argKeys = new ArrayList<>(dt.getTypeArguments().size());
            for (TypeMirror ta : dt.getTypeArguments()) {
              argKeys.add(canonicalGoTypeKey(ta));
            }
            yield baseKey + "[" + String.join(",", argKeys) + "]";
          }
          yield baseKey;
        }
        yield runtime.canonicalTypeKey(te);
      }
      default -> "any";
    };
  }

  /** Zero / default Go value expression for a Java {@link TypeMirror}. */
  public String zeroValue(TypeMirror tm, PackagePlanner.ImportTracker imports) {
    if (tm == null || tm.getKind() == TypeKind.VOID || tm.getKind() == TypeKind.NONE) {
      return "";
    }
    if (isActiveTypeVariable(tm)) {
      String tvGo = goType(tm, imports);
      String rt = runtime.j2gortAlias(imports);
      String rtPref = rt.isEmpty() ? "" : (rt + ".");
      return "func() "
          + tvGo
          + " { if _, ok_ := any((*"
          + tvGo
          + ")(nil)).(*string); ok_ { return any("
          + rtPref
          + "NullString).("
          + tvGo
          + ") }; var z_ "
          + tvGo
          + "; return z_ }()";
    }
    TypeMirror er = safeErasure(tm);
    return switch (er.getKind()) {
      case BOOLEAN -> "false";
      case BYTE -> "int8(0)";
      case SHORT -> "int16(0)";
      case CHAR -> "uint16(0)";
      case INT -> "int32(0)";
      case LONG -> "int64(0)";
      case FLOAT -> "float32(0)";
      case DOUBLE -> "float64(0)";
      case ARRAY -> "nil";
      case DECLARED -> {
        TypeElement te = (TypeElement) ((DeclaredType) er).asElement();
        if (te.getQualifiedName().contentEquals("java.lang.String")) {
          String rt = runtime.j2gortAlias(imports);
          yield (rt.isEmpty() ? "" : rt + ".") + "NullString";
        }
        if (isNodeClass(te) || isValueEnum(te)) {
          yield goTypeForElement(te, imports) + "(0)";
        }
        if (isValueStructClass(te) || isByteStringClass(te)) {
          yield "(" + goTypeForElement(te, imports) + "{})";
        }
        yield "nil";
      }
      default -> "nil";
    };
  }

  public TypeElement getSuperclassElement(TypeElement te) {
    TypeMirror sc = te.getSuperclass();
    if (sc instanceof DeclaredType dt && dt.asElement() instanceof TypeElement ste) {
      String qn = ste.getQualifiedName().toString();
      if (!qn.equals("java.lang.Object")
          && !qn.equals("java.lang.Enum")
          && !qn.equals("java.lang.Record")) {
        return ste;
      }
    }
    return null;
  }

  public List<TypeElement> getAllAncestorTypes(TypeElement te) {
    List<TypeElement> out = new ArrayList<>();
    Set<TypeElement> visited = new LinkedHashSet<>();
    Deque<TypeElement> q = new ArrayDeque<>();
    q.add(te);
    visited.add(te);
    while (!q.isEmpty()) {
      TypeElement cur = q.poll();
      if (cur != te) {
        out.add(cur);
      }
      if (cur.getSuperclass() instanceof DeclaredType dt && dt.asElement() instanceof TypeElement ste) {
        if (!ste.getQualifiedName().contentEquals("java.lang.Object") && visited.add(ste)) {
          q.add(ste);
        }
      }
      for (TypeMirror im : cur.getInterfaces()) {
        if (im instanceof DeclaredType dt && dt.asElement() instanceof TypeElement ite) {
          if (visited.add(ite)) {
            q.add(ite);
          }
        }
      }
    }
    return out;
  }

  public boolean isSyntheticArrayElement(TypeElement te) {
    if (te == null) {
      return false;
    }
    return te.getQualifiedName().contentEquals("Array")
        || te.asType().getKind() == TypeKind.ARRAY;
  }

  public void recordExternalType(TypeElement te) {
    if (te == null || isSourceType(te) || isSyntheticArrayElement(te)) {
      return;
    }
    if (referencedExternalTypes.add(te)) {
      for (TypeElement anc : getAllAncestorTypes(te)) {
        if (!isSourceType(anc) && !isSyntheticArrayElement(anc)) {
          referencedExternalTypes.add(anc);
        }
      }
    }
  }

  /** Finds the single abstract method (JLS §9.8) of a functional interface {@code fi}. */
  public ExecutableElement findSamMethod(TypeElement fi) {
    List<ExecutableElement> candidates = new ArrayList<>();
    for (Element m : elements.getAllMembers(fi)) {
      if (m.getKind() != ElementKind.METHOD || !(m instanceof ExecutableElement ee)) {
        continue;
      }
      if (ee.getModifiers().contains(Modifier.STATIC)
          || ee.getModifiers().contains(Modifier.DEFAULT)
          || !ee.getModifiers().contains(Modifier.ABSTRACT)) {
        continue;
      }
      if (isPublicObjectMethod(ee)) {
        continue;
      }
      candidates.add(ee);
    }
    if (candidates.size() <= 1) {
      return candidates.isEmpty() ? null : candidates.get(0);
    }
    // Filter out methods overridden by another candidate or default method on fi
    for (ExecutableElement c : candidates) {
      boolean shadowedByDefault = false;
      for (Element m : elements.getAllMembers(fi)) {
        if (m instanceof ExecutableElement other
            && other != c
            && other.getModifiers().contains(Modifier.DEFAULT)
            && elements.overrides(other, c, fi)) {
          shadowedByDefault = true;
          break;
        }
      }
      if (!shadowedByDefault) {
        return c;
      }
    }
    return candidates.get(0);
  }

  private boolean isPublicObjectMethod(ExecutableElement ee) {
    String name = ee.getSimpleName().toString();
    int arity = ee.getParameters().size();
    if (name.equals("toString") && arity == 0) return true;
    if (name.equals("hashCode") && arity == 0) return true;
    if (name.equals("equals") && arity == 1) {
      TypeMirror p = safeErasure(ee.getParameters().get(0).asType());
      return p instanceof DeclaredType dt
          && dt.asElement() instanceof TypeElement te
          && te.getQualifiedName().contentEquals("java.lang.Object");
    }
    return false;
  }

  /** Full whole-program type, hierarchy, field, capture, and Union-Find analysis. */
  public void analyze(List<CompilationUnitTree> units) {
    int[] anonCounter = {0};
    Map<TypeElement, Set<TypeElement>> localInstantiations = new HashMap<>();
    Set<TypeElement> externallySynchronizedTypes = new HashSet<>();

    // Step 1: Discover local/anonymous classes, enum constant bodies, external refs, cross-pkg fields, captures
    for (PackagePlanner.UnitPlan plan : planner.getUnitPlans()) {
      CompilationUnitTree cu = plan.unit;
      String curGoPkg = plan.goPkg;

      new TreePathScanner<Void, TypeElement>() {
        private TypeElement currentTop;
        private final Deque<Set<VariableElement>> localDeclaredScopes = new ArrayDeque<>();
        private final Deque<TypeElement> localOrAnonStack = new ArrayDeque<>();
        private final Map<TypeElement, Set<VariableElement>> capturedSets = new HashMap<>();

        @Override
        public Void visitClass(ClassTree node, TypeElement enclosingType) {
          Element el = trees.getElement(getCurrentPath());
          TypeElement te = (el instanceof TypeElement t) ? t : null;
          if (te != null) {
            classTreeElements.put(node, te);
            typeElementClassTrees.put(te, node);
            if (te.getNestingKind() == NestingKind.TOP_LEVEL) {
              currentTop = te;
            } else if (te.getNestingKind() == NestingKind.LOCAL
                || te.getNestingKind() == NestingKind.ANONYMOUS) {
              // Check if this is an enum constant body
              Tree parentTree = getCurrentPath().getParentPath().getLeaf();
              boolean isEnumConstBody = false;
              if (parentTree instanceof NewClassTree) {
                Tree gp = getCurrentPath().getParentPath().getParentPath().getLeaf();
                if (gp instanceof VariableTree vt) {
                  Element vel = trees.getElement(getCurrentPath().getParentPath().getParentPath());
                  if (vel instanceof VariableElement ve && ve.getKind() == ElementKind.ENUM_CONSTANT) {
                    enumConstantBodyTypes.put(ve, te);
                    enumConstantAnonTypes.add(te);
                    isEnumConstBody = true;
                  }
                }
              }
              if (currentTop != null) {
                planner.registerLocalOrAnonType(
                    te, currentTop, ++anonCounter[0], te.getNestingKind() == NestingKind.ANONYMOUS);
              }
              if (!isEnumConstBody) {
                localOrAnonStack.push(te);
                capturedSets.put(te, new LinkedHashSet<>());
              }
            }
          }
          localDeclaredScopes.push(new HashSet<>());
          Void r = super.visitClass(node, te != null ? te : enclosingType);
          localDeclaredScopes.pop();
          if (te != null && !localOrAnonStack.isEmpty() && localOrAnonStack.peek() == te) {
            localOrAnonStack.pop();
            capturedVars.put(te, new ArrayList<>(capturedSets.get(te)));
          }
          return r;
        }

        @Override
        public Void visitMethod(MethodTree node, TypeElement enclosingType) {
          localDeclaredScopes.push(new HashSet<>());
          Void r = super.visitMethod(node, enclosingType);
          localDeclaredScopes.pop();
          return r;
        }

        @Override
        public Void visitVariable(VariableTree node, TypeElement enclosingType) {
          Element el = trees.getElement(getCurrentPath());
          if (el instanceof VariableElement ve) {
            if (ve.getKind() == ElementKind.LOCAL_VARIABLE
                || ve.getKind() == ElementKind.PARAMETER
                || ve.getKind() == ElementKind.EXCEPTION_PARAMETER
                || ve.getKind() == ElementKind.RESOURCE_VARIABLE
                || ve.getKind() == ElementKind.BINDING_VARIABLE) {
              if (!localDeclaredScopes.isEmpty()) {
                localDeclaredScopes.peek().add(ve);
              }
            }
            recordTypeMirror(ve.asType());
          }
          return super.visitVariable(node, enclosingType);
        }

        @Override
        public Void visitLambdaExpression(LambdaExpressionTree node, TypeElement enclosingType) {
          TypeMirror tm = trees.getTypeMirror(getCurrentPath());
          recordLambdaTarget(tm);
          localDeclaredScopes.push(new HashSet<>());
          Void r = super.visitLambdaExpression(node, enclosingType);
          localDeclaredScopes.pop();
          return r;
        }

        @Override
        public Void visitMemberReference(MemberReferenceTree node, TypeElement enclosingType) {
          TypeMirror tm = trees.getTypeMirror(getCurrentPath());
          recordLambdaTarget(tm);
          Element refEl = trees.getElement(getCurrentPath());
          recordElement(refEl, curGoPkg, false);
          return super.visitMemberReference(node, enclosingType);
        }

        @Override
        public Void visitIdentifier(IdentifierTree node, TypeElement enclosingType) {
          Element el = trees.getElement(getCurrentPath());
          if (el instanceof VariableElement ve) {
            if (ve.getKind() == ElementKind.LOCAL_VARIABLE
                || ve.getKind() == ElementKind.PARAMETER
                || ve.getKind() == ElementKind.EXCEPTION_PARAMETER
                || ve.getKind() == ElementKind.RESOURCE_VARIABLE
                || ve.getKind() == ElementKind.BINDING_VARIABLE) {
              // If declared outside a currently active local/anonymous class, mark as captured!
              for (TypeElement activeLocalAnon : localOrAnonStack) {
                Element varOwner = ve.getEnclosingElement();
                if (!isEnclosedWithin(varOwner, activeLocalAnon)) {
                  capturedSets.get(activeLocalAnon).add(ve);
                }
              }
            }
          }
          recordElement(el, curGoPkg, isWriteTarget(getCurrentPath()));
          return super.visitIdentifier(node, enclosingType);
        }

        @Override
        public Void visitMemberSelect(MemberSelectTree node, TypeElement enclosingType) {
          Element el = trees.getElement(getCurrentPath());
          recordElement(el, curGoPkg, isWriteTarget(getCurrentPath()));
          return super.visitMemberSelect(node, enclosingType);
        }

        @Override
        public Void visitMethodInvocation(MethodInvocationTree node, TypeElement enclosingType) {
          Element el = trees.getElement(getCurrentPath());
          recordElement(el, curGoPkg, false);
          recordTypeMirror(trees.getTypeMirror(getCurrentPath()));
          return super.visitMethodInvocation(node, enclosingType);
        }

        @Override
        public Void visitNewClass(NewClassTree node, TypeElement enclosingType) {
          Element el = trees.getElement(getCurrentPath());
          if (el instanceof ExecutableElement ctor
              && ctor.getEnclosingElement() instanceof TypeElement targetTe
              && targetTe.getNestingKind() == NestingKind.LOCAL) {
            for (TypeElement activeLocalAnon : localOrAnonStack) {
              localInstantiations
                  .computeIfAbsent(activeLocalAnon, k -> new LinkedHashSet<>())
                  .add(targetTe);
            }
          }
          recordElement(el, curGoPkg, false);
          recordTypeMirror(trees.getTypeMirror(getCurrentPath()));
          return super.visitNewClass(node, enclosingType);
        }

        @Override
        public Void visitSynchronized(SynchronizedTree node, TypeElement enclosingType) {
          TreePath exprPath = new TreePath(getCurrentPath(), node.getExpression());
          while (exprPath.getLeaf() instanceof ParenthesizedTree pt) {
            exprPath = new TreePath(exprPath, pt.getExpression());
          }
          Tree leaf = exprPath.getLeaf();
          if (!(leaf instanceof IdentifierTree id && id.getName().contentEquals("this"))) {
            if (leaf instanceof MemberSelectTree ms && ms.getIdentifier().contentEquals("class")) {
              Element clsEl = trees.getElement(new TreePath(exprPath, ms.getExpression()));
              if (clsEl instanceof TypeElement clsTe) {
                externallySynchronizedTypes.add(clsTe);
              }
            }
            TypeMirror lockTm = safeErasure(trees.getTypeMirror(exprPath));
            if (lockTm instanceof DeclaredType dt && dt.asElement() instanceof TypeElement lockTe) {
              externallySynchronizedTypes.add(lockTe);
            }
          }
          return super.visitSynchronized(node, enclosingType);
        }
      }.scan(cu, null);
    }

    // Propagate transitive captures when a local/anonymous class instantiates a sibling/enclosing local class
    boolean capChanged = true;
    while (capChanged) {
      capChanged = false;
      for (Map.Entry<TypeElement, Set<TypeElement>> e : localInstantiations.entrySet()) {
        TypeElement caller = e.getKey();
        List<VariableElement> callerCaps = capturedVars.computeIfAbsent(caller, k -> new ArrayList<>());
        for (TypeElement callee : e.getValue()) {
          for (VariableElement v : capturedVars.getOrDefault(callee, List.of())) {
            if (!isEnclosedWithin(v.getEnclosingElement(), caller) && !callerCaps.contains(v)) {
              callerCaps.add(v);
              capChanged = true;
            }
          }
        }
      }
    }

    // Seed core JDK collection & closeable types used by runtime helpers BEFORE Union-Find
    for (String coreQn :
        List.of(
            "java.lang.AutoCloseable",
            "java.lang.Iterable",
            "java.util.Iterator",
            "java.util.Collection",
            "java.util.List",
            "java.util.Set",
            "java.util.Map",
            "java.util.ArrayList",
            "java.util.LinkedHashSet",
            "java.util.LinkedHashMap",
            "java.util.function.IntUnaryOperator")) {
      TypeElement coreTe = elements.getTypeElement(coreQn);
      if (coreTe != null) {
        recordExternalType(coreTe);
      }
    }

    // Step 2: Compute subclassedTypes and collect all ancestor external types
    for (TypeElement te : new ArrayList<>(planner.getAllSourceTypes())) {
      if (enumConstantAnonTypes.contains(te)) {
        continue;
      }
      if (te.getModifiers().contains(Modifier.ABSTRACT) && te.getKind() == ElementKind.CLASS) {
        subclassedTypes.add(te);
      }
      TypeElement sup = getSuperclassElement(te);
      if (sup != null) {
        if (isSourceType(sup) && !isNodeClass(sup)) {
          subclassedTypes.add(sup);
        } else if (!isSourceType(sup)) {
          recordExternalType(sup);
          runtime.markSubclassedExternal(sup);
        }
      }
      for (TypeElement anc : getAllAncestorTypes(te)) {
        if (!isSourceType(anc)) {
          recordExternalType(anc);
        }
      }
      // Also record parameter/return/field types of all members
      for (Element m : te.getEnclosedElements()) {
        if (m instanceof VariableElement ve) {
          recordTypeMirror(ve.asType());
        } else if (m instanceof ExecutableElement ee) {
          recordTypeMirror(ee.getReturnType());
          for (VariableElement p : ee.getParameters()) {
            recordTypeMirror(p.asType());
          }
        }
      }
    }

    // Step 2b: Identify fieldless in-tree value enums (`uint8`) and non-reentrant `sync.Mutex` classes
    computeValueEnumTypes();
    computeFastSyncTypes(externallySynchronizedTypes);

    // Step 3: Assign field Go names, disambiguating shadowed fields across inheritance
    computeFieldNames();
    computeValueStructEnumSentinels();

    // Step 4: Assign constructor, static method, and Union-Find virtual method signatures
    computeExecutableSignatures();

    // Step 5: Mark which unified MethodSigs are needed on external interfaces/classes
    computeNeededExternalSigs();

    // Step 6: Compute non-overridden direct base dispatch slots & unused `self_` / `BindSelf_` elision
    computeDirectBaseDispatchAndSelfTypes();
  }

  private void computeNeededExternalSigs() {
    for (ExecutableElement ee : referencedExternalExecs) {
      if (ee.getKind() == ElementKind.METHOD && !ee.getModifiers().contains(Modifier.STATIC)) {
        MethodSig sig = instanceMethodSigs.get(ee);
        if (sig != null) {
          neededExternalSigs.add(sig);
        }
      }
    }
    for (TypeElement srcTe : planner.getAllSourceTypes()) {
      for (Element m : srcTe.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)) {
          MethodSig sig = instanceMethodSigs.get(ee);
          if (sig != null) {
            neededExternalSigs.add(sig);
          }
        }
      }
    }
    for (TypeElement fi : lambdaTargetInterfaces) {
      ExecutableElement sam = findSamMethod(fi);
      if (sam != null) {
        MethodSig sig = instanceMethodSigs.get(sam);
        if (sig != null) {
          neededExternalSigs.add(sig);
        }
      }
    }
    // Always include essential runtime helper methods on Iterable / Iterator / AutoCloseable / Collections
    Set<String> coreHelperMethods =
        Set.of(
            "java.lang.Iterable#iterator",
            "java.util.Iterator#hasNext",
            "java.util.Iterator#next",
            "java.util.Iterator#remove",
            "java.lang.AutoCloseable#close",
            "java.util.Collection#size",
            "java.util.Collection#isEmpty",
            "java.util.Collection#contains",
            "java.util.Collection#add",
            "java.util.Collection#remove",
            "java.util.Collection#clear",
            "java.util.List#get",
            "java.util.Map#size",
            "java.util.Map#get",
            "java.util.Map#getOrDefault",
            "java.util.Map#put",
            "java.util.Map#remove",
            "java.util.Map#containsKey",
            "java.util.Map#keySet",
            "java.util.Map#values");
    for (TypeElement extTe : new ArrayList<>(referencedExternalTypes)) {
      String qn = extTe.getQualifiedName().toString();
      if (extTe.getKind() == ElementKind.INTERFACE && qn.startsWith("java.util.function.")) {
        ExecutableElement sam = findSamMethod(extTe);
        if (sam != null) {
          MethodSig sig = instanceMethodSigs.get(sam);
          if (sig != null) {
            neededExternalSigs.add(sig);
          }
        }
      }
      for (Element m : extTe.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)) {
          if (coreHelperMethods.contains(qn + "#" + ee.getSimpleName())
              && !(qn.equals("java.util.Map")
                  && ee.getSimpleName().contentEquals("remove")
                  && ee.getParameters().size() != 1)) {
            MethodSig sig = instanceMethodSigs.get(ee);
            if (sig != null) {
              neededExternalSigs.add(sig);
            }
          }
        }
      }
    }
    for (MethodSig sig : neededExternalSigs) {
      if (sig.canonReturnType != null) {
        recordTypeMirror(sig.canonReturnType);
      }
      for (TypeMirror pm : sig.canonParamTypes) {
        recordTypeMirror(pm);
      }
    }
  }

  private void computeValueEnumTypes() {
    for (TypeElement te : planner.getAllSourceTypes()) {
      if (te.getKind() != ElementKind.ENUM
          || !isSourceType(te)
          || enumConstantAnonTypes.contains(te)
          || !te.getInterfaces().isEmpty()
          || te.getModifiers().contains(Modifier.ABSTRACT)
          || hasOuterInstance(te)
          || !getCapturedVars(te).isEmpty()) {
        continue;
      }
      ClassTree ct = typeElementClassTrees.get(te);
      TreePath ctPath = trees.getPath(te);
      if (ct == null || ctPath == null) {
        continue;
      }
      int enumConstCount = 0;
      boolean eligible = true;
      List<VariableElement> instanceFields = new ArrayList<>();
      for (Element enc : te.getEnclosedElements()) {
        if (enc instanceof VariableElement ve) {
          if (ve.getKind() == ElementKind.ENUM_CONSTANT) {
            enumConstCount++;
            if (enumConstantBodyTypes.containsKey(ve)) {
              eligible = false;
              break;
            }
          } else if (ve.getKind() == ElementKind.FIELD
              && !ve.getModifiers().contains(Modifier.STATIC)) {
            if (!ve.getModifiers().contains(Modifier.FINAL)
                || (!ve.asType().getKind().isPrimitive() && !isStringType(ve.asType()))) {
              eligible = false;
              break;
            }
            instanceFields.add(ve);
          }
        } else if (enc instanceof ExecutableElement ee && ee.getKind() == ElementKind.METHOD) {
          if (ee.getModifiers().contains(Modifier.ABSTRACT)
              || (!ee.getModifiers().contains(Modifier.STATIC)
                  && ee.getModifiers().contains(Modifier.SYNCHRONIZED))) {
            eligible = false;
            break;
          }
        }
      }
      int maxConsts = instanceFields.isEmpty() ? 254 : 250;
      if (!eligible || enumConstCount < 1 || enumConstCount > maxConsts) {
        continue;
      }
      for (Tree member : ct.getMembers()) {
        if (member instanceof BlockTree bt && !bt.isStatic() && !bt.getStatements().isEmpty()) {
          eligible = false;
          break;
        }
      }
      if (!eligible) {
        continue;
      }
      if (instanceFields.isEmpty()) {
        for (Tree member : ct.getMembers()) {
          if (member instanceof MethodTree mt
              && (mt.getReturnType() == null || mt.getName().contentEquals("<init>"))) {
            if (!mt.getParameters().isEmpty() || !isTrivialEmptyEnumCtorBody(mt.getBody())) {
              eligible = false;
              break;
            }
          }
          if (member instanceof VariableTree vt) {
            Element mel = trees.getElement(new TreePath(ctPath, vt));
            if (mel instanceof VariableElement ve
                && ve.getKind() == ElementKind.ENUM_CONSTANT
                && vt.getInitializer() instanceof NewClassTree nct
                && (nct.getClassBody() != null || !nct.getArguments().isEmpty())) {
              eligible = false;
              break;
            }
          }
        }
        if (eligible) {
          valueEnumTypes.add(te);
        }
      } else {
        Map<VariableElement, List<Object>> table =
            tryEvaluateConstantEnumFields(te, ct, ctPath, instanceFields);
        if (table != null) {
          valueEnumTypes.add(te);
          valueEnumFieldTables.put(te, table);
        }
      }
    }
  }

  private Map<VariableElement, List<Object>> tryEvaluateConstantEnumFields(
      TypeElement enumTe,
      ClassTree ct,
      TreePath ctPath,
      List<VariableElement> instanceFields) {
    Map<ExecutableElement, MethodTree> ctorTrees = new LinkedHashMap<>();
    Map<ExecutableElement, TreePath> ctorPaths = new LinkedHashMap<>();
    Map<VariableElement, Object> declDefaults = new LinkedHashMap<>();
    List<VariableTree> enumConstTrees = new ArrayList<>();

    for (Tree member : ct.getMembers()) {
      if (member instanceof MethodTree mt
          && (mt.getReturnType() == null || mt.getName().contentEquals("<init>"))) {
        TreePath mtPath = new TreePath(ctPath, mt);
        if (!(trees.getElement(mtPath) instanceof ExecutableElement ctorEe)
            || ctorEe.isVarArgs()
            || mt.getBody() == null) {
          return null;
        }
        for (VariableElement pve : ctorEe.getParameters()) {
          if (!pve.asType().getKind().isPrimitive() && !isStringType(pve.asType())) {
            return null;
          }
        }
        ctorTrees.put(ctorEe, mt);
        ctorPaths.put(ctorEe, mtPath);
      } else if (member instanceof VariableTree vt) {
        TreePath vtPath = new TreePath(ctPath, vt);
        Element el = trees.getElement(vtPath);
        if (el instanceof VariableElement ve) {
          if (ve.getKind() == ElementKind.ENUM_CONSTANT) {
            if (vt.getInitializer() instanceof NewClassTree nct && nct.getClassBody() != null) {
              return null;
            }
            enumConstTrees.add(vt);
          } else if (instanceFields.contains(ve) && vt.getInitializer() != null) {
            Object val =
                evalSimpleEnumConstExpr(
                    new TreePath(vtPath, vt.getInitializer()), Map.of(), ve.asType());
            if (val == null) {
              return null;
            }
            declDefaults.put(ve, val);
          }
        }
      }
    }

    Map<VariableElement, List<Object>> result = new LinkedHashMap<>();
    for (VariableElement f : instanceFields) {
      result.put(f, new ArrayList<>(enumConstTrees.size()));
    }

    for (VariableTree vt : enumConstTrees) {
      TreePath vtPath = new TreePath(ctPath, vt);
      ExecutableElement targetCtor = null;
      List<Object> argVals = new ArrayList<>();

      if (vt.getInitializer() instanceof NewClassTree nct) {
        TreePath nctPath = new TreePath(vtPath, nct);
        if (!(trees.getElement(nctPath) instanceof ExecutableElement cee)
            || !ctorTrees.containsKey(cee)
            || nct.getArguments().size() != cee.getParameters().size()) {
          return null;
        }
        targetCtor = cee;
        for (int i = 0; i < nct.getArguments().size(); i++) {
          Object av =
              evalSimpleEnumConstExpr(
                  new TreePath(nctPath, nct.getArguments().get(i)),
                  Map.of(),
                  cee.getParameters().get(i).asType());
          if (av == null) {
            return null;
          }
          argVals.add(av);
        }
      } else {
        for (ExecutableElement cee : ctorTrees.keySet()) {
          if (cee.getParameters().isEmpty()) {
            targetCtor = cee;
            break;
          }
        }
        if (targetCtor == null && ctorTrees.isEmpty() && declDefaults.size() == instanceFields.size()) {
          for (VariableElement f : instanceFields) {
            result.get(f).add(declDefaults.get(f));
          }
          continue;
        }
        if (targetCtor == null) {
          return null;
        }
      }

      Map<VariableElement, Object> fieldEnv =
          evalEnumCtorBody(
              targetCtor, argVals, ctorTrees, ctorPaths, instanceFields, declDefaults, 0);
      if (fieldEnv == null || !fieldEnv.keySet().containsAll(instanceFields)) {
        return null;
      }
      for (VariableElement f : instanceFields) {
        result.get(f).add(fieldEnv.get(f));
      }
    }
    return result;
  }

  private Map<VariableElement, Object> evalEnumCtorBody(
      ExecutableElement ctorEe,
      List<Object> argVals,
      Map<ExecutableElement, MethodTree> ctorTrees,
      Map<ExecutableElement, TreePath> ctorPaths,
      List<VariableElement> instanceFields,
      Map<VariableElement, Object> declDefaults,
      int depth) {
    if (depth > 8) {
      return null;
    }
    MethodTree mt = ctorTrees.get(ctorEe);
    TreePath mtPath = ctorPaths.get(ctorEe);
    if (mt == null || mtPath == null || mt.getBody() == null) {
      return null;
    }
    Map<VariableElement, Object> paramEnv = new HashMap<>();
    for (int i = 0; i < ctorEe.getParameters().size(); i++) {
      paramEnv.put(ctorEe.getParameters().get(i), argVals.get(i));
    }

    TreePath bodyPath = new TreePath(mtPath, mt.getBody());
    List<? extends StatementTree> stmts = mt.getBody().getStatements();

    // Check for `this(...)` delegation as first statement
    if (!stmts.isEmpty()
        && stmts.get(0) instanceof ExpressionStatementTree est0
        && est0.getExpression() instanceof MethodInvocationTree call0
        && call0.getMethodSelect() instanceof IdentifierTree id0
        && id0.getName().contentEquals("this")) {
      for (int i = 1; i < stmts.size(); i++) {
        if (stmts.get(i).getKind() != Tree.Kind.EMPTY_STATEMENT) {
          return null;
        }
      }
      TreePath call0Path = new TreePath(new TreePath(bodyPath, est0), call0);
      if (!(trees.getElement(call0Path) instanceof ExecutableElement nextCtor)
          || !ctorTrees.containsKey(nextCtor)
          || call0.getArguments().size() != nextCtor.getParameters().size()) {
        return null;
      }
      List<Object> nextArgs = new ArrayList<>(call0.getArguments().size());
      for (int i = 0; i < call0.getArguments().size(); i++) {
        Object av =
            evalSimpleEnumConstExpr(
                new TreePath(call0Path, call0.getArguments().get(i)),
                paramEnv,
                nextCtor.getParameters().get(i).asType());
        if (av == null) {
          return null;
        }
        nextArgs.add(av);
      }
      return evalEnumCtorBody(
          nextCtor, nextArgs, ctorTrees, ctorPaths, instanceFields, declDefaults, depth + 1);
    }

    Map<VariableElement, Object> fieldEnv = new LinkedHashMap<>(declDefaults);
    for (int i = 0; i < stmts.size(); i++) {
      StatementTree st = stmts.get(i);
      if (st.getKind() == Tree.Kind.EMPTY_STATEMENT) {
        continue;
      }
      if (i == 0
          && st instanceof ExpressionStatementTree estS
          && estS.getExpression() instanceof MethodInvocationTree callS
          && callS.getArguments().isEmpty()
          && callS.getMethodSelect() instanceof IdentifierTree idS
          && idS.getName().contentEquals("super")) {
        continue;
      }
      if (!(st instanceof ExpressionStatementTree est)
          || !(est.getExpression() instanceof AssignmentTree at)
          || at.getKind() != Tree.Kind.ASSIGNMENT) {
        return null;
      }
      TreePath atPath = new TreePath(new TreePath(bodyPath, est), at);
      TreePath lhsPath = unwrapParensPathLocal(new TreePath(atPath, at.getVariable()));
      Tree lhsLeaf = lhsPath.getLeaf();
      if (lhsLeaf instanceof MemberSelectTree ms) {
        TreePath qPath = unwrapParensPathLocal(new TreePath(lhsPath, ms.getExpression()));
        if (!(qPath.getLeaf() instanceof IdentifierTree qid)
            || !qid.getName().contentEquals("this")) {
          return null;
        }
      } else if (!(lhsLeaf instanceof IdentifierTree)) {
        return null;
      }
      if (!(trees.getElement(lhsPath) instanceof VariableElement lhsVe)
          || !instanceFields.contains(lhsVe)) {
        return null;
      }
      Object rhsVal =
          evalSimpleEnumConstExpr(
              new TreePath(atPath, at.getExpression()), paramEnv, lhsVe.asType());
      if (rhsVal == null) {
        return null;
      }
      fieldEnv.put(lhsVe, rhsVal);
    }
    return fieldEnv;
  }

  private Object evalSimpleEnumConstExpr(
      TreePath exprPath, Map<VariableElement, Object> paramEnv, TypeMirror targetType) {
    exprPath = unwrapParensPathLocal(exprPath);
    if (exprPath == null) {
      return null;
    }
    Tree leaf = exprPath.getLeaf();
    if (leaf instanceof LiteralTree lt) {
      if (lt.getKind() == Tree.Kind.NULL_LITERAL) {
        return isStringType(targetType) ? ENUM_NULL_STRING : null;
      }
      Object v = lt.getValue();
      if (v instanceof Float f && (f.isNaN() || f.isInfinite())) {
        return null;
      }
      if (v instanceof Double d && (d.isNaN() || d.isInfinite())) {
        return null;
      }
      return coerceEnumConstValue(v, targetType);
    }
    if (leaf instanceof IdentifierTree || leaf instanceof MemberSelectTree) {
      if (leaf instanceof MemberSelectTree ms) {
        Element qEl = trees.getElement(new TreePath(exprPath, ms.getExpression()));
        if (!(qEl instanceof TypeElement)) {
          return null;
        }
      }
      Element el = trees.getElement(exprPath);
      if (el instanceof VariableElement ve) {
        if (paramEnv.containsKey(ve)) {
          return coerceEnumConstValue(paramEnv.get(ve), targetType);
        }
        if (ve.getModifiers().contains(Modifier.STATIC)
            && ve.getModifiers().contains(Modifier.FINAL)
            && ve.getConstantValue() != null) {
          return coerceEnumConstValue(ve.getConstantValue(), targetType);
        }
      }
      return null;
    }
    if (leaf instanceof UnaryTree ut) {
      TypeMirror innerTm = trees.getTypeMirror(exprPath);
      Object sub =
          evalSimpleEnumConstExpr(
              new TreePath(exprPath, ut.getExpression()),
              paramEnv,
              innerTm != null ? innerTm : targetType);
      if (sub == null || sub == ENUM_NULL_STRING) {
        return null;
      }
      Object res =
          switch (ut.getKind()) {
            case UNARY_PLUS -> sub;
            case UNARY_MINUS -> {
              if (sub instanceof Double d) yield -d;
              if (sub instanceof Float f) yield -f;
              if (sub instanceof Long l) yield -l;
              if (sub instanceof Number n) yield -n.intValue();
              if (sub instanceof Character c) yield -((int) c);
              yield null;
            }
            case BITWISE_COMPLEMENT -> {
              if (sub instanceof Long l) yield ~l;
              if (sub instanceof Number n) yield ~n.intValue();
              if (sub instanceof Character c) yield ~((int) c);
              yield null;
            }
            case LOGICAL_COMPLEMENT -> (sub instanceof Boolean b) ? !b : null;
            default -> null;
          };
      return coerceEnumConstValue(res, targetType);
    }
    return null;
  }

  private Object coerceEnumConstValue(Object val, TypeMirror targetType) {
    if (val == null || targetType == null) {
      return null;
    }
    if (isStringType(targetType)) {
      return (val == ENUM_NULL_STRING || val instanceof String) ? val : null;
    }
    if (val == ENUM_NULL_STRING) {
      return null;
    }
    long numL = 0;
    double numD = 0;
    boolean isNum = false;
    if (val instanceof Number n) {
      numL = n.longValue();
      numD = n.doubleValue();
      isNum = true;
    } else if (val instanceof Character c) {
      numL = c;
      numD = c;
      isNum = true;
    }
    return switch (targetType.getKind()) {
      case BOOLEAN -> (val instanceof Boolean) ? val : null;
      case BYTE -> isNum ? (byte) numL : null;
      case SHORT -> isNum ? (short) numL : null;
      case CHAR -> isNum ? (char) numL : null;
      case INT -> isNum ? (int) numL : null;
      case LONG -> isNum ? numL : null;
      case FLOAT -> isNum ? (float) numD : null;
      case DOUBLE -> isNum ? numD : null;
      default -> null;
    };
  }

  private static TreePath unwrapParensPathLocal(TreePath tp) {
    while (tp != null && tp.getLeaf() instanceof ParenthesizedTree pt) {
      tp = new TreePath(tp, pt.getExpression());
    }
    return tp;
  }

  private void computeDirectBaseDispatchAndSelfTypes() {
    List<TypeElement> allSourceTypes = new ArrayList<>(planner.getAllSourceTypes());
    Map<MethodSig, Set<TypeElement>> declaringTypesBySig = new HashMap<>();
    for (TypeElement t : allSourceTypes) {
      if (enumConstantAnonTypes.contains(t)) {
        continue;
      }
      for (Element enc : t.getEnclosedElements()) {
        if (enc instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)
            && !ee.getModifiers().contains(Modifier.PRIVATE)) {
          MethodSig sig = instanceMethodSigs.get(ee);
          if (sig != null) {
            declaringTypesBySig.computeIfAbsent(sig, k -> new HashSet<>()).add(t);
          }
        }
      }
    }

    Map<TypeElement, Set<TypeElement>> properSubtypes = new HashMap<>();
    for (TypeElement base : allSourceTypes) {
      if (isSubclassed(base)) {
        properSubtypes.put(base, new HashSet<>());
      }
    }
    for (TypeElement t : allSourceTypes) {
      if (enumConstantAnonTypes.contains(t)) {
        continue;
      }
      for (TypeElement anc : getAllAncestorTypes(t)) {
        Set<TypeElement> subs = properSubtypes.get(anc);
        if (subs != null) {
          subs.add(t);
        }
      }
    }

    for (Map.Entry<TypeElement, Set<TypeElement>> e : properSubtypes.entrySet()) {
      TypeElement base = e.getKey();
      Set<TypeElement> subs = e.getValue();
      Set<MethodSig> directSigs = new HashSet<>();

      // Collect concrete non-abstract methods declared on `base` or its source superclasses
      Map<MethodSig, Boolean> concreteFromClassChain = new LinkedHashMap<>();
      TypeElement cur = base;
      while (cur != null && isSourceType(cur)) {
        for (Element enc : cur.getEnclosedElements()) {
          if (enc instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && !ee.getModifiers().contains(Modifier.STATIC)
              && !ee.getModifiers().contains(Modifier.PRIVATE)) {
            MethodSig sig = instanceMethodSigs.get(ee);
            if (sig != null && !concreteFromClassChain.containsKey(sig)) {
              concreteFromClassChain.put(sig, !ee.getModifiers().contains(Modifier.ABSTRACT));
            }
          }
        }
        cur = getSuperclassElement(cur);
      }

      for (Map.Entry<MethodSig, Boolean> ce : concreteFromClassChain.entrySet()) {
        if (!ce.getValue()) {
          continue;
        }
        MethodSig sig = ce.getKey();
        boolean overriddenBelow = false;
        for (TypeElement declT : declaringTypesBySig.getOrDefault(sig, Set.of())) {
          if (subs.contains(declT)) {
            overriddenBelow = true;
            break;
          }
        }
        if (!overriddenBelow) {
          directSigs.add(sig);
        }
      }
      if (!directSigs.isEmpty()) {
        directBaseDispatchSigs.put(base, directSigs);
      }
    }

    // Now determine which subclassed source classes may emit `this_.self_`
    for (TypeElement base : properSubtypes.keySet()) {
      if (classMayEmitSelfRef(base, false)) {
        mayHaveSelfFieldTypes.add(base);
      }
    }

    // In top-down superclass-before-subclass order, identify ancestor classes that genuinely
    // retain their own `self_` field AND have at least one subclassed descendant in the same Go
    // package, so subclasses can reuse `this_.<AncBase>.self_` without emitting redundant `self_`
    // fields on intermediate base structs (such as PrototypeObjectTypeBase, D_AutoValue_ColorBase).
    List<TypeElement> sortedBases = new ArrayList<>(mayHaveSelfFieldTypes);
    sortedBases.sort(Comparator.comparingInt(this::sourceSuperclassDepth));
    for (TypeElement base : sortedBases) {
      if (base.getKind() != ElementKind.CLASS
          || isNodeClass(base)
          || isValueStructClass(base)) {
        continue;
      }
      if (getSamePkgReusableSelfRoot(base) != null) {
        continue;
      }
      String basePkg = planner.getGoPkg(base);
      boolean hasSubclassedSamePkgSub = false;
      for (TypeElement sub : properSubtypes.getOrDefault(base, Set.of())) {
        if (!sub.equals(base)
            && properSubtypes.containsKey(sub)
            && basePkg.equals(planner.getGoPkg(sub))) {
          hasSubclassedSamePkgSub = true;
          break;
        }
      }
      if (!hasSubclassedSamePkgSub) {
        continue;
      }
      List<TypeParameterElement> effTpes = getEffectiveTypeParameters(base);
      if (!effTpes.isEmpty()) {
        pushTypeParamScope(base, effTpes);
      }
      try {
        if (classMayEmitSelfRef(base, true)) {
          reusableSelfAncestors.add(base);
        }
      } finally {
        if (!effTpes.isEmpty()) {
          popTypeParamScope();
        }
      }
    }
  }

  private int sourceSuperclassDepth(TypeElement te) {
    int d = 0;
    TypeElement cur = getSuperclassElement(te);
    while (cur != null && isSourceType(cur) && !runtime.mapsToAny(cur) && d < 64) {
      d++;
      cur = getSuperclassElement(cur);
    }
    return d;
  }

  private boolean classMayEmitSelfRef(TypeElement base, boolean allowAncestorReuse) {
    // 1. Check if `emitInheritedInterfaceDefaultForwarders` will emit a forwarder on `base`
    Set<String> covered = new HashSet<>();
    TypeElement curClass = base;
    while (curClass != null && !runtime.mapsToAny(curClass)) {
      if (runtime.isShimException(curClass)) {
        covered.add("GetMessage");
        covered.add("GetCause");
        covered.add("Error");
      } else if (isSourceType(curClass)) {
        for (Element m : curClass.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && !ee.getModifiers().contains(Modifier.STATIC)
              && !ee.getModifiers().contains(Modifier.PRIVATE)
              && !ee.getModifiers().contains(Modifier.ABSTRACT)) {
            MethodSig sig = instanceMethodSigs.get(ee);
            if (sig != null) {
              covered.add(sig.goName);
            }
          }
        }
        if (curClass != base) {
          // Any default/Object-fallback method already emitted on a source superclass is promoted!
          for (TypeElement supAnc : getAllAncestorTypes(curClass)) {
            if (supAnc.getKind() == ElementKind.INTERFACE && !runtime.mapsToAny(supAnc)) {
              for (Element m : supAnc.getEnclosedElements()) {
                if (m instanceof ExecutableElement ee
                    && ee.getKind() == ElementKind.METHOD
                    && !ee.getModifiers().contains(Modifier.STATIC)
                    && !ee.getModifiers().contains(Modifier.PRIVATE)) {
                  MethodSig sig = instanceMethodSigs.get(ee);
                  if (sig != null
                      && (ee.getModifiers().contains(Modifier.DEFAULT)
                          || isPublicObjectMethodFromSig(sig))
                      && (isSourceType(supAnc) || isNeededExternalSig(sig))) {
                    covered.add(sig.goName);
                  }
                }
              }
            }
          }
        }
      } else {
        List<TypeElement> extChain = new ArrayList<>();
        extChain.add(curClass);
        extChain.addAll(getAllAncestorTypes(curClass));
        for (TypeElement extAnc : extChain) {
          if (runtime.mapsToAny(extAnc) || runtime.isShimException(extAnc)) {
            continue;
          }
          for (Element m : extAnc.getEnclosedElements()) {
            if (m instanceof ExecutableElement ee
                && ee.getKind() == ElementKind.METHOD
                && !ee.getModifiers().contains(Modifier.STATIC)
                && !ee.getModifiers().contains(Modifier.PRIVATE)) {
              MethodSig sig = instanceMethodSigs.get(ee);
              if (sig != null && isNeededExternalSig(sig)) {
                covered.add(sig.goName);
              }
            }
          }
        }
      }
      curClass = getSuperclassElement(curClass);
    }

    for (TypeElement anc : getAllAncestorTypes(base)) {
      if (anc.getKind() == ElementKind.INTERFACE && !runtime.mapsToAny(anc)) {
        for (Element m : anc.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && !ee.getModifiers().contains(Modifier.STATIC)
              && !ee.getModifiers().contains(Modifier.PRIVATE)) {
            MethodSig sig = instanceMethodSigs.get(ee);
            if (sig != null
                && (ee.getModifiers().contains(Modifier.DEFAULT)
                    || isPublicObjectMethodFromSig(sig))
                && (isSourceType(anc) || isNeededExternalSig(sig))
                && !covered.contains(sig.goName)) {
              if (allowAncestorReuse) {
                if (!ee.getModifiers().contains(Modifier.DEFAULT)
                    && isPublicObjectMethodFromSig(sig)
                    && findReusableSelfAncestor(base, objectElement.asType(), true) != null) {
                  continue;
                }
                if (ee.getModifiers().contains(Modifier.DEFAULT)
                    && isSourceType(anc)
                    && base.asType() instanceof DeclaredType baseDt) {
                  DeclaredType ifaceSite = findSuperType(baseDt, anc);
                  if (ifaceSite != null
                      && findReusableSelfAncestor(base, ifaceSite, false) != null) {
                    continue;
                  }
                }
              }
              return true;
            }
          }
        }
      }
    }

    // 2. Scan `base`'s own AST members (excluding nested ClassTrees) for any selfValueExpr() trigger
    TreePath basePath = trees.getPath(base);
    ClassTree baseCt = typeElementClassTrees.get(base);
    if (basePath == null || baseCt == null) {
      return true;
    }
    boolean[] mayUseSelf = {false};
    for (Tree member : baseCt.getMembers()) {
      if (member instanceof ClassTree) {
        continue;
      }
      if (member instanceof BlockTree bt && bt.isStatic()) {
        continue;
      }
      if (member instanceof VariableTree vt) {
        Element vEl = trees.getElement(new TreePath(basePath, vt));
        if (vEl instanceof VariableElement ve && ve.getModifiers().contains(Modifier.STATIC)) {
          continue;
        }
      }
      if (member instanceof MethodTree mt) {
        Element mEl = trees.getElement(new TreePath(basePath, mt));
        if (mEl instanceof ExecutableElement ee && ee.getModifiers().contains(Modifier.STATIC)) {
          continue;
        }
      }

      new TreePathScanner<Void, Void>() {
        @Override
        public Void visitClass(ClassTree node, Void p) {
          // Anonymous class in non-static context captures `this_.self_` as `outer_`!
          TypeElement nTe = classTreeElements.get(node);
          if (nTe != null && hasOuterInstance(nTe)) {
            TypeElement outerTe = getEnclosingInstanceType(nTe);
            if (outerTe != null && outerTe.equals(base)) {
              mayUseSelf[0] = true;
            }
          }
          return null;
        }

        @Override
        public Void visitIdentifier(IdentifierTree node, Void p) {
          if (mayUseSelf[0]) {
            return null;
          }
          if (node.getName().contentEquals("this")) {
            TreePath curPath = getCurrentPath();
            TreePath pPath = curPath.getParentPath();
            while (pPath != null && pPath.getLeaf() instanceof ParenthesizedTree) {
              curPath = pPath;
              pPath = pPath.getParentPath();
            }
            if (pPath != null) {
              Tree pLeaf = pPath.getLeaf();
              if (pLeaf instanceof MethodInvocationTree mit
                  && mit.getMethodSelect() == curPath.getLeaf()) {
                // `this(...)` constructor chaining call
                return null;
              }
              if (pLeaf instanceof MemberSelectTree ms && ms.getExpression() == curPath.getLeaf()) {
                Element selEl = trees.getElement(pPath);
                if (selEl instanceof VariableElement ve
                    && ve.getKind() == ElementKind.FIELD
                    && !ve.getModifiers().contains(Modifier.STATIC)) {
                  // `this.field` accesses `this_` directly
                  return null;
                }
                if (selEl instanceof ExecutableElement ee
                    && ee.getKind() == ElementKind.METHOD
                    && !ee.getModifiers().contains(Modifier.STATIC)) {
                  if (ee.getModifiers().contains(Modifier.PRIVATE)
                      && ee.getEnclosingElement().equals(base)) {
                    return null;
                  }
                  if (ee.getEnclosingElement() instanceof TypeElement mOwner
                      && isSourceType(mOwner)
                      && mOwner.getKind() == ElementKind.CLASS
                      && canDispatchDirectlyOnBase(base, instanceMethodSigs.get(ee))) {
                    return null;
                  }
                  if (allowAncestorReuse
                      && ee.getEnclosingElement() instanceof TypeElement mOwner
                      && !(mOwner.getKind() == ElementKind.INTERFACE
                          && ee.getModifiers().contains(Modifier.PRIVATE))
                      && findReusableSelfAncestorForCall(base, mOwner, instanceMethodSigs.get(ee))
                          != null) {
                    return null;
                  }
                }
              }
              if (allowAncestorReuse && isThisExprSatisfiedByAncestor(base, curPath, pPath)) {
                return null;
              }
            }
            mayUseSelf[0] = true;
            return null;
          }
          return super.visitIdentifier(node, p);
        }

        @Override
        public Void visitMemberSelect(MemberSelectTree node, Void p) {
          if (mayUseSelf[0]) {
            return null;
          }
          if (node.getIdentifier().contentEquals("this")) {
            mayUseSelf[0] = true;
            return null;
          }
          return super.visitMemberSelect(node, p);
        }

        @Override
        public Void visitMethodInvocation(MethodInvocationTree node, Void p) {
          if (mayUseSelf[0]) {
            return null;
          }
          if (node.getMethodSelect() instanceof IdentifierTree id) {
            String nm = id.getName().toString();
            if (!nm.equals("this") && !nm.equals("super")) {
              Element el = trees.getElement(getCurrentPath());
              if (el instanceof ExecutableElement ee
                  && !ee.getModifiers().contains(Modifier.STATIC)) {
                if (!ee.getModifiers().contains(Modifier.PRIVATE)
                    || !ee.getEnclosingElement().equals(base)) {
                  boolean directBase =
                      ee.getEnclosingElement() instanceof TypeElement mOwner
                          && isSourceType(mOwner)
                          && mOwner.getKind() == ElementKind.CLASS
                          && types.isSubtype(safeErasure(base.asType()), safeErasure(mOwner.asType()))
                          && canDispatchDirectlyOnBase(base, instanceMethodSigs.get(ee));
                  boolean ancSelf =
                      allowAncestorReuse
                          && ee.getEnclosingElement() instanceof TypeElement mOwner2
                          && !(mOwner2.getKind() == ElementKind.INTERFACE
                              && ee.getModifiers().contains(Modifier.PRIVATE))
                          && types.isSubtype(
                              safeErasure(base.asType()), safeErasure(mOwner2.asType()))
                          && findReusableSelfAncestorForCall(
                                  base, mOwner2, instanceMethodSigs.get(ee))
                              != null;
                  if (!directBase && !ancSelf) {
                    mayUseSelf[0] = true;
                  }
                }
              }
            }
          }
          return super.visitMethodInvocation(node, p);
        }

        @Override
        public Void visitNewClass(NewClassTree node, Void p) {
          if (mayUseSelf[0]) {
            return null;
          }
          if (node.getEnclosingExpression() == null) {
            Element el = trees.getElement(getCurrentPath());
            if (el instanceof ExecutableElement ctor
                && ctor.getEnclosingElement() instanceof TypeElement targetOwner
                && isSourceType(targetOwner)
                && hasOuterInstance(targetOwner)) {
              TypeElement outerTe = getEnclosingInstanceType(targetOwner);
              if (outerTe != null && outerTe.equals(base)) {
                mayUseSelf[0] = true;
              }
            }
          }
          return super.visitNewClass(node, p);
        }
      }.scan(new TreePath(basePath, member), null);
      if (mayUseSelf[0]) {
        return true;
      }
    }
    return false;
  }

  private boolean isThisExprSatisfiedByAncestor(
      TypeElement base, TreePath curPath, TreePath pPath) {
    Tree pLeaf = pPath.getLeaf();
    Tree curLeaf = curPath.getLeaf();
    if (pLeaf instanceof VariableTree vt && vt.getInitializer() == curLeaf) {
      Element vEl = trees.getElement(pPath);
      return vEl instanceof VariableElement ve
          && findReusableSelfAncestor(base, ve.asType(), false) != null;
    }
    if (pLeaf instanceof AssignmentTree at
        && at.getKind() == Tree.Kind.ASSIGNMENT
        && at.getExpression() == curLeaf) {
      TypeMirror lhsTm = trees.getTypeMirror(new TreePath(pPath, at.getVariable()));
      return lhsTm != null && findReusableSelfAncestor(base, lhsTm, false) != null;
    }
    if (pLeaf instanceof ReturnTree rt && rt.getExpression() == curLeaf) {
      TreePath enc = pPath.getParentPath();
      while (enc != null) {
        Tree encLeaf = enc.getLeaf();
        if (encLeaf instanceof LambdaExpressionTree || encLeaf instanceof ClassTree) {
          return false;
        }
        if (encLeaf instanceof MethodTree) {
          Element mEl = trees.getElement(enc);
          if (mEl instanceof ExecutableElement ee) {
            MethodSig mSig = instanceMethodSigs.get(ee);
            TypeMirror retTm =
                (mSig != null && mSig.canonReturnType != null)
                    ? resolveSigReturnType(mSig, base.asType())
                    : ee.getReturnType();
            boolean forceAny = mSig != null && mSig.returnWidenedToAny;
            return findReusableSelfAncestor(base, retTm, forceAny) != null;
          }
          return false;
        }
        enc = enc.getParentPath();
      }
      return false;
    }
    if (pLeaf instanceof MethodInvocationTree mit) {
      int idx = mit.getArguments().indexOf(curLeaf);
      if (idx >= 0) {
        Element callEl = trees.getElement(pPath);
        if (callEl instanceof ExecutableElement ee
            && !ee.isVarArgs()
            && idx < ee.getParameters().size()) {
          if (ee.getModifiers().contains(Modifier.STATIC)
              && ee.getSimpleName().contentEquals("identityHashCode")
              && ee.getEnclosingElement() instanceof TypeElement sysTe
              && sysTe.getQualifiedName().contentEquals("java.lang.System")) {
            return findReusableSelfAncestor(base, objectElement.asType(), true) != null;
          }
          TypeMirror recvTm = null;
          if (mit.getMethodSelect() instanceof MemberSelectTree cms) {
            recvTm =
                trees.getTypeMirror(
                    new TreePath(new TreePath(pPath, cms), cms.getExpression()));
          }
          ExecutableType execType = resolveExecutableType(recvTm, ee);
          TypeMirror paramTm =
              execType != null && idx < execType.getParameterTypes().size()
                  ? execType.getParameterTypes().get(idx)
                  : ee.getParameters().get(idx).asType();
          return findReusableSelfAncestor(base, paramTm, false) != null;
        }
      }
      return false;
    }
    if (pLeaf instanceof NewClassTree nct && nct.getClassBody() == null) {
      int idx = nct.getArguments().indexOf(curLeaf);
      if (idx >= 0) {
        Element ctorEl = trees.getElement(pPath);
        TypeMirror newTm = trees.getTypeMirror(pPath);
        if (ctorEl instanceof ExecutableElement ctor
            && !ctor.isVarArgs()
            && idx < ctor.getParameters().size()) {
          ExecutableType execType = resolveExecutableType(newTm, ctor);
          TypeMirror paramTm =
              execType != null && idx < execType.getParameterTypes().size()
                  ? execType.getParameterTypes().get(idx)
                  : ctor.getParameters().get(idx).asType();
          return findReusableSelfAncestor(base, paramTm, false) != null;
        }
      }
      return false;
    }
    if (pLeaf instanceof BinaryTree bt
        && (bt.getKind() == Tree.Kind.EQUAL_TO || bt.getKind() == Tree.Kind.NOT_EQUAL_TO)) {
      ExpressionTree other =
          bt.getLeftOperand() == curLeaf ? bt.getRightOperand() : bt.getLeftOperand();
      TypeMirror otherTm = trees.getTypeMirror(new TreePath(pPath, other));
      if (otherTm != null && otherTm.getKind() != TypeKind.NULL) {
        if (isObjectOrErasedAny(otherTm)) {
          return findReusableSelfAncestor(base, objectElement.asType(), true) != null;
        }
        return findReusableSelfAncestor(base, goEffectiveType(otherTm), false) != null;
      }
    }
    return false;
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
      if (i == 0
          && st instanceof ExpressionStatementTree est
          && est.getExpression() instanceof MethodInvocationTree call0
          && call0.getArguments().isEmpty()
          && call0.getMethodSelect() instanceof IdentifierTree id
          && id.getName().contentEquals("super")) {
        continue;
      }
      return false;
    }
    return true;
  }

  private void computeFastSyncTypes(Set<TypeElement> externallySynchronizedTypes) {
    for (TypeElement te : planner.getAllSourceTypes()) {
      if (enumConstantAnonTypes.contains(te) || externallySynchronizedTypes.contains(te)) {
        continue;
      }
      TreePath ctPath = trees.getPath(te);
      if (ctPath == null || !(ctPath.getLeaf() instanceof ClassTree ct)) {
        continue;
      }

      // 1. Check `static synchronized` methods on `te`
      Set<ExecutableElement> staticMethods = new HashSet<>();
      Set<ExecutableElement> staticSyncMethods = new HashSet<>();
      for (Element enc : te.getEnclosedElements()) {
        if (enc instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && ee.getModifiers().contains(Modifier.STATIC)) {
          staticMethods.add(ee);
          if (ee.getModifiers().contains(Modifier.SYNCHRONIZED)) {
            staticSyncMethods.add(ee);
          }
        }
      }
      if (!staticSyncMethods.isEmpty()) {
        boolean[] staticHasSyncBlock = {false};
        Map<ExecutableElement, Set<ExecutableElement>> staticEdges = new HashMap<>();
        for (Tree member : ct.getMembers()) {
          if (member instanceof BlockTree bt && bt.isStatic()) {
            new TreePathScanner<Void, Void>() {
              @Override
              public Void visitClass(ClassTree node, Void p) {
                return null;
              }

              @Override
              public Void visitSynchronized(SynchronizedTree node, Void p) {
                staticHasSyncBlock[0] = true;
                return null;
              }
            }.scan(new TreePath(ctPath, bt), null);
          } else if (member instanceof MethodTree mt) {
            TreePath mtPath = new TreePath(ctPath, mt);
            if (trees.getElement(mtPath) instanceof ExecutableElement ee
                && staticMethods.contains(ee)) {
              Set<ExecutableElement> callees = new HashSet<>();
              staticEdges.put(ee, callees);
              new TreePathScanner<Void, Void>() {
                @Override
                public Void visitClass(ClassTree node, Void p) {
                  return null;
                }

                @Override
                public Void visitSynchronized(SynchronizedTree node, Void p) {
                  staticHasSyncBlock[0] = true;
                  return null;
                }

                @Override
                public Void visitMethodInvocation(MethodInvocationTree node, Void p) {
                  if (trees.getElement(getCurrentPath()) instanceof ExecutableElement target
                      && staticMethods.contains(target)) {
                    callees.add(target);
                  }
                  return super.visitMethodInvocation(node, p);
                }

                @Override
                public Void visitMemberReference(MemberReferenceTree node, Void p) {
                  if (trees.getElement(getCurrentPath()) instanceof ExecutableElement target
                      && staticMethods.contains(target)) {
                    callees.add(target);
                  }
                  return super.visitMemberReference(node, p);
                }
              }.scan(mtPath, null);
            }
          }
        }
        if (!staticHasSyncBlock[0] && !hasCycleOrSyncToSyncReachability(staticSyncMethods, staticEdges)) {
          fastStaticSyncTypes.add(te);
        }
      }

      // 2. Check leaf non-reentrant instance `synchronized` methods & `synchronized (this)` blocks
      if (te.getKind() != ElementKind.CLASS
          || isSubclassed(te)
          || isNodeClass(te)
          || isValueStructClass(te)
          || isValueEnum(te)
          || te.getSimpleName().toString().startsWith("AutoValue_")) {
        continue;
      }
      if (!ancestorsFreeOfInstanceSynchronization(te, externallySynchronizedTypes)) {
        continue;
      }

      Set<ExecutableElement> ownAndAncInstMethods = collectHierarchyInstanceMethods(te);
      Set<ExecutableElement> syncInstMethods = new HashSet<>();
      Map<ExecutableElement, Set<ExecutableElement>> instEdges = new HashMap<>();
      boolean[] instDisqualified = {false};

      for (Tree member : ct.getMembers()) {
        if (member instanceof BlockTree bt && !bt.isStatic()) {
          new TreePathScanner<Void, Void>() {
            @Override
            public Void visitClass(ClassTree node, Void p) {
              return null;
            }

            @Override
            public Void visitSynchronized(SynchronizedTree node, Void p) {
              instDisqualified[0] = true;
              return null;
            }
          }.scan(new TreePath(ctPath, bt), null);
        } else if (member instanceof MethodTree mt) {
          TreePath mtPath = new TreePath(ctPath, mt);
          if (!(trees.getElement(mtPath) instanceof ExecutableElement ee)) {
            continue;
          }
          if (ee.getModifiers().contains(Modifier.STATIC)) {
            continue;
          }
          boolean isCtor = ee.getKind() == ElementKind.CONSTRUCTOR;
          boolean isMethodSync = !isCtor && ee.getModifiers().contains(Modifier.SYNCHRONIZED);
          if (isMethodSync) {
            syncInstMethods.add(ee);
          }
          Set<ExecutableElement> callees = new HashSet<>();
          if (!isCtor) {
            instEdges.put(ee, callees);
          }
          int[] syncBlockCount = {0};
          new TreePathScanner<Void, int[]>() {
            @Override
            public Void visitClass(ClassTree node, int[] d) {
              return null;
            }

            @Override
            public Void visitLambdaExpression(LambdaExpressionTree node, int[] d) {
              return super.visitLambdaExpression(node, new int[] {d[0], d[1], d[2] + 1});
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
            public Void visitSynchronized(SynchronizedTree node, int[] d) {
              syncBlockCount[0]++;
              if (isCtor || isMethodSync || d[0] > 0 || d[1] > 0 || d[2] > 0 || syncBlockCount[0] > 1) {
                instDisqualified[0] = true;
                return null;
              }
              ExpressionTree lockExpr = node.getExpression();
              while (lockExpr instanceof ParenthesizedTree pt) {
                lockExpr = pt.getExpression();
              }
              if (!(lockExpr instanceof IdentifierTree id && id.getName().contentEquals("this"))) {
                instDisqualified[0] = true;
                return null;
              }
              syncInstMethods.add(ee);
              return super.visitSynchronized(node, new int[] {d[0], d[1] + 1, d[2]});
            }

            @Override
            public Void visitMethodInvocation(MethodInvocationTree node, int[] d) {
              if (trees.getElement(getCurrentPath()) instanceof ExecutableElement target) {
                String tName = target.getSimpleName().toString();
                if (tName.equals("wait") || tName.equals("notify") || tName.equals("notifyAll")) {
                  instDisqualified[0] = true;
                }
                if (ownAndAncInstMethods.contains(target)) {
                  callees.add(target);
                }
              }
              return super.visitMethodInvocation(node, d);
            }

            @Override
            public Void visitMemberReference(MemberReferenceTree node, int[] d) {
              if (trees.getElement(getCurrentPath()) instanceof ExecutableElement target
                  && ownAndAncInstMethods.contains(target)) {
                callees.add(target);
              }
              return super.visitMemberReference(node, d);
            }
          }.scan(mtPath, new int[] {0, 0, 0});
        }
      }

      if (!instDisqualified[0]
          && !syncInstMethods.isEmpty()
          && !hasCycleOrSyncToSyncReachability(syncInstMethods, instEdges)) {
        fastInstanceSyncTypes.add(te);
      }
    }
  }

  private boolean ancestorsFreeOfInstanceSynchronization(
      TypeElement te, Set<TypeElement> externallySynchronizedTypes) {
    TypeElement cur = getSuperclassElement(te);
    while (cur != null) {
      if (!isSourceType(cur) || externallySynchronizedTypes.contains(cur)) {
        return false;
      }
      TreePath ancPath = trees.getPath(cur);
      if (ancPath == null) {
        return false;
      }
      for (Element m : cur.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && !ee.getModifiers().contains(Modifier.STATIC)
            && ee.getModifiers().contains(Modifier.SYNCHRONIZED)) {
          return false;
        }
      }
      boolean[] hasSyncStmt = {false};
      new TreePathScanner<Void, Void>() {
        @Override
        public Void visitClass(ClassTree node, Void p) {
          if (getCurrentPath() == ancPath) {
            return super.visitClass(node, p);
          }
          return null;
        }

        @Override
        public Void visitSynchronized(SynchronizedTree node, Void p) {
          hasSyncStmt[0] = true;
          return null;
        }
      }.scan(ancPath, null);
      if (hasSyncStmt[0]) {
        return false;
      }
      cur = getSuperclassElement(cur);
    }
    return true;
  }

  private Set<ExecutableElement> collectHierarchyInstanceMethods(TypeElement te) {
    Set<ExecutableElement> out = new HashSet<>();
    TypeElement cur = te;
    while (cur != null && isSourceType(cur)) {
      for (Element m : cur.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)) {
          out.add(ee);
        }
      }
      cur = getSuperclassElement(cur);
    }
    return out;
  }

  private static boolean hasCycleOrSyncToSyncReachability(
      Set<ExecutableElement> syncRoots, Map<ExecutableElement, Set<ExecutableElement>> edges) {
    for (ExecutableElement start : syncRoots) {
      Set<ExecutableElement> visited = new HashSet<>();
      Deque<ExecutableElement> q = new ArrayDeque<>(edges.getOrDefault(start, Set.of()));
      while (!q.isEmpty()) {
        ExecutableElement cur = q.poll();
        if (syncRoots.contains(cur)) {
          return true;
        }
        if (visited.add(cur)) {
          q.addAll(edges.getOrDefault(cur, Set.of()));
        }
      }
    }
    return false;
  }

  private boolean isEnclosedWithin(Element inner, Element outer) {
    Element cur = inner;
    while (cur != null) {
      if (cur.equals(outer)) {
        return true;
      }
      cur = cur.getEnclosingElement();
    }
    return false;
  }

  private boolean isWriteTarget(TreePath path) {
    Tree leaf = path.getLeaf();
    TreePath parent = path.getParentPath();
    if (parent == null) {
      return false;
    }
    Tree p = parent.getLeaf();
    return switch (p.getKind()) {
      case ASSIGNMENT,
          MULTIPLY_ASSIGNMENT,
          DIVIDE_ASSIGNMENT,
          REMAINDER_ASSIGNMENT,
          PLUS_ASSIGNMENT,
          MINUS_ASSIGNMENT,
          LEFT_SHIFT_ASSIGNMENT,
          RIGHT_SHIFT_ASSIGNMENT,
          UNSIGNED_RIGHT_SHIFT_ASSIGNMENT,
          AND_ASSIGNMENT,
          XOR_ASSIGNMENT,
          OR_ASSIGNMENT -> {
        if (p instanceof com.sun.source.tree.AssignmentTree at) {
          yield at.getVariable() == leaf;
        }
        if (p instanceof com.sun.source.tree.CompoundAssignmentTree cat) {
          yield cat.getVariable() == leaf;
        }
        yield false;
      }
      case PREFIX_INCREMENT, PREFIX_DECREMENT, POSTFIX_INCREMENT, POSTFIX_DECREMENT -> true;
      default -> false;
    };
  }

  private void recordLambdaTarget(TypeMirror tm) {
    if (tm == null) {
      return;
    }
    TypeMirror er = safeErasure(tm);
    if (er instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      if (te.getKind() == ElementKind.INTERFACE) {
        lambdaTargetInterfaces.add(te);
        if (!isSourceType(te)) {
          recordExternalType(te);
        }
      }
    }
  }

  private void recordTypeMirror(TypeMirror tm) {
    if (tm == null) {
      return;
    }
    if (tm instanceof ArrayType rawAt) {
      recordTypeMirror(rawAt.getComponentType());
      return;
    }
    if (tm instanceof WildcardType wt) {
      if (wt.getExtendsBound() != null) {
        recordTypeMirror(wt.getExtendsBound());
      }
      if (wt.getSuperBound() != null) {
        recordTypeMirror(wt.getSuperBound());
      }
      return;
    }
    if (tm instanceof DeclaredType rawDt
        && rawDt.asElement() instanceof TypeElement rawTe
        && isNativeGenericClass(rawTe)) {
      if (rawTe.getKind() == ElementKind.INTERFACE
          && rawDt.getTypeArguments().size() == rawTe.getTypeParameters().size()) {
        for (int k = 0; k < rawDt.getTypeArguments().size(); k++) {
          if (rawDt.getTypeArguments().get(k) instanceof WildcardType) {
            wildcardedInterfaceParams.computeIfAbsent(rawTe, x -> new HashSet<>()).add(k);
          }
        }
        for (TypeElement anc : getAllAncestorTypes(rawTe)) {
          if (anc.getKind() == ElementKind.INTERFACE && isNativeGenericClass(anc)) {
            DeclaredType supDt = findSuperType(rawDt, anc);
            if (supDt != null
                && supDt.getTypeArguments().size() == anc.getTypeParameters().size()) {
              for (int ak = 0; ak < supDt.getTypeArguments().size(); ak++) {
                if (supDt.getTypeArguments().get(ak) instanceof WildcardType) {
                  wildcardedInterfaceParams.computeIfAbsent(anc, x -> new HashSet<>()).add(ak);
                }
              }
            }
          }
        }
      }
      for (TypeMirror ta : rawDt.getTypeArguments()) {
        recordTypeMirror(ta);
      }
    }
    TypeMirror er = safeErasure(tm);
    if (er instanceof ArrayType at) {
      recordTypeMirror(at.getComponentType());
    } else if (er instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      if (!isSourceType(te)) {
        recordExternalType(te);
      }
    }
  }

  private void recordElement(Element el, String curGoPkg, boolean isWrite) {
    if (el == null) {
      return;
    }
    if (el instanceof TypeElement te) {
      if (!isSourceType(te)) {
        recordExternalType(te);
      }
      return;
    }
    Element enc = el.getEnclosingElement();
    if (!(enc instanceof TypeElement owner) || isSyntheticArrayElement(owner)) {
      return;
    }
    if (isSourceType(owner)) {
      if (el instanceof VariableElement ve && ve.getKind() == ElementKind.FIELD) {
        if (!ve.getModifiers().contains(Modifier.STATIC)) {
          String ownerGoPkg = planner.getGoPkg(owner);
          if (!ownerGoPkg.equals(curGoPkg)) {
            crossPkgReadFields.add(ve);
            if (isWrite) {
              crossPkgWriteFields.add(ve);
            }
          }
        }
      }
    } else {
      recordExternalType(owner);
      if (el instanceof ExecutableElement ee) {
        referencedExternalExecs.add(ee);
        recordTypeMirror(ee.getReturnType());
        for (VariableElement p : ee.getParameters()) {
          recordTypeMirror(p.asType());
        }
      } else if (el instanceof VariableElement ve) {
        referencedExternalFields.add(ve);
        recordTypeMirror(ve.asType());
      }
    }
  }

  private void computeFieldNames() {
    for (TypeElement te : planner.getAllSourceTypes()) {
      Set<String> ancestorFieldNames = new HashSet<>();
      TypeElement cur = getSuperclassElement(te);
      while (cur != null && isSourceType(cur)) {
        for (Element m : cur.getEnclosedElements()) {
          if (m instanceof VariableElement ve
              && ve.getKind() == ElementKind.FIELD
              && !ve.getModifiers().contains(Modifier.STATIC)) {
            ancestorFieldNames.add(ve.getSimpleName().toString());
          }
        }
        cur = getSuperclassElement(cur);
      }

      for (Element m : te.getEnclosedElements()) {
        if (m instanceof VariableElement ve && ve.getKind() == ElementKind.FIELD) {
          if (ve.getModifiers().contains(Modifier.STATIC)) {
            continue;
          }
          String raw = ve.getSimpleName().toString();
          if (ancestorFieldNames.contains(raw)) {
            fieldGoNames.put(ve, "f_" + planner.getGoTypeName(te) + "_" + GoNames.cleanChars(raw));
          } else {
            fieldGoNames.put(ve, GoNames.unexported(raw));
          }
        }
      }
    }
  }

  private void computeValueStructEnumSentinels() {
    for (TypeElement te : planner.getAllSourceTypes()) {
      if (!isNullableValueStructClass(te) || hasValueStructValidFlag(te)) {
        continue;
      }
      ClassTree ct = typeElementClassTrees.get(te);
      TreePath ctPath = trees.getPath(te);
      if (ct == null || ctPath == null) {
        continue;
      }
      List<MethodTree> ctors = new ArrayList<>();
      Map<VariableElement, VariableTree> fieldTrees = new LinkedHashMap<>();
      for (Tree member : ct.getMembers()) {
        if (member instanceof MethodTree mt
            && (mt.getReturnType() == null || mt.getName().contentEquals("<init>"))) {
          ctors.add(mt);
        } else if (member instanceof VariableTree vt) {
          Element el = trees.getElement(new TreePath(ctPath, vt));
          if (el instanceof VariableElement ve
              && ve.getKind() == ElementKind.FIELD
              && !ve.getModifiers().contains(Modifier.STATIC)) {
            fieldTrees.put(ve, vt);
          }
        }
      }
      for (Map.Entry<VariableElement, VariableTree> entry : fieldTrees.entrySet()) {
        VariableElement ve = entry.getKey();
        VariableTree vt = entry.getValue();
        if (!ve.getModifiers().contains(Modifier.FINAL) || !isValueEnumType(ve.asType())) {
          continue;
        }
        boolean guaranteedNonNull = false;
        if (vt.getInitializer() != null) {
          TreePath initPath =
              new TreePath(new TreePath(ctPath, vt), vt.getInitializer());
          guaranteedNonNull = isGuaranteedNonNullEnumExpr(initPath, Set.of());
        } else if (!ctors.isEmpty()) {
          guaranteedNonNull = true;
          for (MethodTree ctorMt : ctors) {
            if (!ctorGuaranteesNonNullEnumField(new TreePath(ctPath, ctorMt), ctorMt, ve)) {
              guaranteedNonNull = false;
              break;
            }
          }
        }
        if (guaranteedNonNull) {
          valueStructEnumSentinelFields.put(te, getFieldGoName(ve));
          break;
        }
      }
    }
  }

  private boolean ctorGuaranteesNonNullEnumField(
      TreePath mtPath, MethodTree mt, VariableElement targetField) {
    if (mt.getBody() == null) {
      return false;
    }
    List<? extends StatementTree> stmts = mt.getBody().getStatements();
    if (stmts.isEmpty()) {
      return false;
    }
    TreePath bodyPath = new TreePath(mtPath, mt.getBody());
    if (stmts.get(0) instanceof ExpressionStatementTree est0
        && est0.getExpression() instanceof MethodInvocationTree call0
        && call0.getMethodSelect() instanceof IdentifierTree id0
        && id0.getName().contentEquals("this")) {
      return true;
    }
    Set<Element> nullCheckedVars = new HashSet<>();
    for (StatementTree st : stmts) {
      TreePath stPath = new TreePath(bodyPath, st);
      if (st instanceof ExpressionStatementTree est) {
        TreePath exprPath = unwrapParensTreePath(new TreePath(stPath, est.getExpression()));
        Tree expr = exprPath.getLeaf();
        if (expr instanceof AssignmentTree at && at.getKind() == Tree.Kind.ASSIGNMENT) {
          TreePath lhsPath = unwrapParensTreePath(new TreePath(exprPath, at.getVariable()));
          if (targetField.equals(trees.getElement(lhsPath))) {
            TreePath rhsPath = unwrapParensTreePath(new TreePath(exprPath, at.getExpression()));
            return isGuaranteedNonNullEnumExpr(rhsPath, nullCheckedVars);
          }
        } else if (expr instanceof MethodInvocationTree mit
            && isStaticNullCheckCall(exprPath)
            && !mit.getArguments().isEmpty()) {
          TreePath arg0Path = unwrapParensTreePath(new TreePath(exprPath, mit.getArguments().get(0)));
          Element argEl = trees.getElement(arg0Path);
          if (argEl != null) {
            nullCheckedVars.add(argEl);
          }
        }
      } else if (st instanceof IfTree ift && ift.getElseStatement() == null) {
        Element guarded = extractIfNullThrowVar(stPath, ift);
        if (guarded != null) {
          nullCheckedVars.add(guarded);
        }
      }
    }
    return false;
  }

  private boolean isGuaranteedNonNullEnumExpr(
      TreePath exprPath, Set<Element> nullCheckedVars) {
    exprPath = unwrapParensTreePath(exprPath);
    Tree leaf = exprPath.getLeaf();
    Element el = trees.getElement(exprPath);
    if (el instanceof VariableElement ve) {
      if (ve.getKind() == ElementKind.ENUM_CONSTANT) {
        return true;
      }
      if ((ve.getKind() == ElementKind.PARAMETER || ve.getKind() == ElementKind.LOCAL_VARIABLE)
          && ve.getModifiers().contains(Modifier.FINAL)
          && nullCheckedVars.contains(ve)) {
        return true;
      }
      if (ve.getKind() == ElementKind.PARAMETER && nullCheckedVars.contains(ve)) {
        return true;
      }
    }
    if (leaf instanceof MethodInvocationTree && isStaticNullCheckCall(exprPath)) {
      return true;
    }
    if (leaf instanceof ConditionalExpressionTree cet) {
      return isGuaranteedNonNullEnumExpr(
              new TreePath(exprPath, cet.getTrueExpression()), nullCheckedVars)
          && isGuaranteedNonNullEnumExpr(
              new TreePath(exprPath, cet.getFalseExpression()), nullCheckedVars);
    }
    return false;
  }

  private boolean isStaticNullCheckCall(TreePath callPath) {
    Element el = trees.getElement(callPath);
    if (el instanceof ExecutableElement ee
        && ee.getModifiers().contains(Modifier.STATIC)
        && !ee.getParameters().isEmpty()
        && ee.getEnclosingElement() instanceof TypeElement owner) {
      String oqn = owner.getQualifiedName().toString();
      String mn = ee.getSimpleName().toString();
      return (oqn.equals("com.google.common.base.Preconditions") && mn.equals("checkNotNull"))
          || (oqn.equals("java.util.Objects") && mn.equals("requireNonNull"));
    }
    return false;
  }

  private Element extractIfNullThrowVar(TreePath ifPath, IfTree ift) {
    StatementTree thenSt = ift.getThenStatement();
    boolean throwsAlways = thenSt instanceof ThrowTree;
    if (!throwsAlways && thenSt instanceof BlockTree bt && !bt.getStatements().isEmpty()) {
      throwsAlways = bt.getStatements().get(0) instanceof ThrowTree;
    }
    if (!throwsAlways) {
      return null;
    }
    TreePath condPath = unwrapParensTreePath(new TreePath(ifPath, ift.getCondition()));
    if (!(condPath.getLeaf() instanceof BinaryTree bin) || bin.getKind() != Tree.Kind.EQUAL_TO) {
      return null;
    }
    TreePath lp = unwrapParensTreePath(new TreePath(condPath, bin.getLeftOperand()));
    TreePath rp = unwrapParensTreePath(new TreePath(condPath, bin.getRightOperand()));
    TreePath varPath =
        rp.getLeaf().getKind() == Tree.Kind.NULL_LITERAL
            ? lp
            : (lp.getLeaf().getKind() == Tree.Kind.NULL_LITERAL ? rp : null);
    if (varPath == null || !(varPath.getLeaf() instanceof IdentifierTree)) {
      return null;
    }
    return trees.getElement(varPath);
  }

  private static TreePath unwrapParensTreePath(TreePath tp) {
    while (tp != null && tp.getLeaf() instanceof ParenthesizedTree pt) {
      tp = new TreePath(tp, pt.getExpression());
    }
    return tp;
  }

  private void computeExecutableSignatures() {
    Set<TypeElement> allTypes = new LinkedHashSet<>();
    allTypes.addAll(planner.getAllSourceTypes());
    allTypes.addAll(referencedExternalTypes);

    // Collect all constructors, static methods, and instance/interface methods
    List<ExecutableElement> allInstanceMethods = new ArrayList<>();
    Map<TypeElement, List<ExecutableElement>> declaredInstanceMethods = new LinkedHashMap<>();

    for (TypeElement te : allTypes) {
      List<ExecutableElement> ctors = new ArrayList<>();
      List<ExecutableElement> statics = new ArrayList<>();
      List<ExecutableElement> insts = new ArrayList<>();

      for (Element enc : te.getEnclosedElements()) {
        if (!(enc instanceof ExecutableElement ee)) {
          continue;
        }
        if (ee.getKind() == ElementKind.CONSTRUCTOR) {
          ctors.add(ee);
        } else if (ee.getKind() == ElementKind.METHOD) {
          if (ee.getModifiers().contains(Modifier.STATIC)) {
            statics.add(ee);
          } else {
            insts.add(ee);
            allInstanceMethods.add(ee);
          }
        }
      }
      declaredInstanceMethods.put(te, insts);
      assignConstructorNames(te, ctors);
      assignStaticMethodNames(te, statics);
    }

    // Union-Find over all instance/interface ExecutableElements
    Map<ExecutableElement, ExecutableElement> parent = new LinkedHashMap<>();
    for (ExecutableElement m : allInstanceMethods) {
      parent.put(m, m);
    }

    for (TypeElement sub : allTypes) {
      List<TypeElement> ancestors = getAllAncestorTypes(sub);
      List<ExecutableElement> subMethods = declaredInstanceMethods.getOrDefault(sub, List.of());

      // 1) Link subMethods with ancestor methods they override
      for (ExecutableElement mSub : subMethods) {
        if (mSub.getModifiers().contains(Modifier.PRIVATE)) {
          continue;
        }
        String name = mSub.getSimpleName().toString();
        for (TypeElement anc : ancestors) {
          for (ExecutableElement mAnc : declaredInstanceMethods.getOrDefault(anc, List.of())) {
            if (!mAnc.getSimpleName().contentEquals(name)) {
              continue;
            }
            if (elements.overrides(mSub, mAnc, sub) || sameErasedJavaParams(mSub, mAnc)) {
              ufUnion(parent, mSub, mAnc);
            }
          }
        }
      }

      // 2) Link inherited ancestor methods that share override-equivalent signatures on `sub`
      //    (e.g., interface C extends A, B or class K extends Base implements Iface)
      Map<String, List<ExecutableElement>> byName = new HashMap<>();
      for (ExecutableElement m : subMethods) {
        if (!m.getModifiers().contains(Modifier.PRIVATE)) {
          byName.computeIfAbsent(m.getSimpleName().toString(), k -> new ArrayList<>()).add(m);
        }
      }
      for (TypeElement anc : ancestors) {
        for (ExecutableElement m : declaredInstanceMethods.getOrDefault(anc, List.of())) {
          if (!m.getModifiers().contains(Modifier.PRIVATE)) {
            byName.computeIfAbsent(m.getSimpleName().toString(), k -> new ArrayList<>()).add(m);
          }
        }
      }
      for (List<ExecutableElement> group : byName.values()) {
        if (group.size() <= 1) {
          continue;
        }
        for (int i = 0; i < group.size(); i++) {
          ExecutableElement a = group.get(i);
          for (int j = i + 1; j < group.size(); j++) {
            ExecutableElement b = group.get(j);
            if (sameErasedJavaParams(a, b)) {
              ufUnion(parent, a, b);
            }
          }
        }
      }
    }

    // Build MethodSig per Union-Find root
    Map<ExecutableElement, List<ExecutableElement>> groups = new LinkedHashMap<>();
    for (ExecutableElement m : allInstanceMethods) {
      ExecutableElement root = ufFind(parent, m);
      groups.computeIfAbsent(root, k -> new ArrayList<>()).add(m);
    }

    Map<ExecutableElement, MethodSig> rootSigs = new LinkedHashMap<>();
    for (Map.Entry<ExecutableElement, List<ExecutableElement>> e : groups.entrySet()) {
      ExecutableElement root = e.getKey();
      List<ExecutableElement> members = e.getValue();
      String simple = root.getSimpleName().toString();
      int arity = root.getParameters().size();

      ExecutableElement canonMember = findCanonicalGenericRootMember(members);
      TypeElement canonOwner =
          canonMember != null ? (TypeElement) canonMember.getEnclosingElement() : null;
      ExecutableElement seedExec = canonMember != null ? canonMember : root;

      if (canonOwner != null) {
        pushTypeParamScope(canonOwner, getEffectiveTypeParameters(canonOwner));
      }
      List<TypeMirror> initParams = new ArrayList<>(arity);
      TypeMirror initRet;
      try {
        for (int i = 0; i < arity; i++) {
          initParams.add(goEffectiveType(seedExec.getParameters().get(i).asType()));
        }
        initRet =
            seedExec.getReturnType().getKind() == TypeKind.VOID
                ? null
                : goEffectiveType(seedExec.getReturnType());
      } finally {
        if (canonOwner != null) {
          popTypeParamScope();
        }
      }

      MethodSig sig = new MethodSig(simple, arity, canonOwner, initRet, initParams);

      // Unify canonical Go parameter and return types across all members in the slot
      for (ExecutableElement m : members) {
        TypeElement mOwner =
            m.getEnclosingElement() instanceof TypeElement t ? t : null;
        List<TypeParameterElement> mEffTpes =
            mOwner != null ? getEffectiveTypeParameters(mOwner) : List.of();
        if (!mEffTpes.isEmpty()) {
          pushTypeParamScope(mOwner, mEffTpes);
        }
        try {
          TypeMirror mSiteTm = mOwner != null ? mOwner.asType() : null;
          for (int i = 0; i < arity; i++) {
            if (!sig.paramWidenedToAny[i]) {
              TypeMirror expectedPm = resolveSigParamType(sig, i, mSiteTm);
              TypeMirror actualPm = goEffectiveType(m.getParameters().get(i).asType());
              if (!canonicalGoTypeKey(expectedPm).equals(canonicalGoTypeKey(actualPm))) {
                sig.paramWidenedToAny[i] = true;
                sig.canonParamTypes.set(i, objectElement.asType());
              }
            }
          }
          if (sig.canonReturnType != null && !sig.returnWidenedToAny) {
            TypeMirror expectedRm = resolveSigReturnType(sig, mSiteTm);
            TypeMirror actualRm = goEffectiveType(m.getReturnType());
            if (!canonicalGoTypeKey(expectedRm).equals(canonicalGoTypeKey(actualRm))) {
              boolean leafImplCovariantOk = false;
              if (canonOwner != null
                  && mOwner != null
                  && mOwner.getKind() == ElementKind.CLASS
                  && !isSubclassed(mOwner)
                  && isGoInterfaceType(expectedRm)
                  && expectedRm instanceof DeclaredType edt
                  && edt.asElement() instanceof TypeElement ete
                  && actualRm instanceof DeclaredType adt) {
                DeclaredType asExp = findSuperType(adt, ete);
                if (asExp != null
                    && canonicalGoTypeKey(expectedRm).equals(canonicalGoTypeKey(goEffectiveType(asExp)))) {
                  leafImplCovariantOk = true;
                }
              }
              if (!leafImplCovariantOk) {
                sig.returnWidenedToAny = true;
                sig.canonReturnType = objectElement.asType();
              }
            }
          }
        } finally {
          if (!mEffTpes.isEmpty()) {
            popTypeParamScope();
          }
        }
      }

      // If `canonOwner` is a generic interface instantiated with `?` on parameter `k` anywhere,
      // widen any method return/param type whose Go type mentions `k` so concrete implementations
      // structurally satisfy the wildcard-instantiated Go interface.
      if (canonOwner != null && canonOwner.getKind() == ElementKind.INTERFACE && isNativeGenericClass(canonOwner)) {
        Set<Integer> wildIdxs = wildcardedInterfaceParams.get(canonOwner);
        if (wildIdxs != null) {
          List<? extends TypeParameterElement> ownerTpes = canonOwner.getTypeParameters();
          for (int k : wildIdxs) {
            if (k >= 0 && k < ownerTpes.size()) {
              TypeParameterElement wildTpe = ownerTpes.get(k);
              if (sig.canonReturnType != null
                  && !sig.returnWidenedToAny
                  && mentionsGoTypeParam(sig.canonReturnType, wildTpe)) {
                sig.returnWidenedToAny = true;
                sig.canonReturnType = objectElement.asType();
              }
              for (int i = 0; i < arity; i++) {
                if (!sig.paramWidenedToAny[i]
                    && mentionsGoTypeParam(sig.canonParamTypes.get(i), wildTpe)) {
                  sig.paramWidenedToAny[i] = true;
                  sig.canonParamTypes.set(i, objectElement.asType());
                }
              }
            }
          }
        }
      }

      rootSigs.put(root, sig);
      for (ExecutableElement m : members) {
        instanceMethodSigs.put(m, sig);
      }
    }

    // Assign collision-free Go method names per MethodSig across all types
    // First find which MethodSigs co-exist on any single TypeElement (declared + ancestors)
    Map<MethodSig, Set<MethodSig>> coexistingSameName = new LinkedHashMap<>();
    for (MethodSig s : rootSigs.values()) {
      coexistingSameName.put(s, new LinkedHashSet<>());
    }

    for (TypeElement te : allTypes) {
      Map<String, Set<MethodSig>> visibleByName = new HashMap<>();
      for (ExecutableElement m : declaredInstanceMethods.getOrDefault(te, List.of())) {
        MethodSig s = instanceMethodSigs.get(m);
        visibleByName.computeIfAbsent(s.simpleName, k -> new LinkedHashSet<>()).add(s);
      }
      for (TypeElement anc : getAllAncestorTypes(te)) {
        for (ExecutableElement m : declaredInstanceMethods.getOrDefault(anc, List.of())) {
          if (!m.getModifiers().contains(Modifier.PRIVATE)) {
            MethodSig s = instanceMethodSigs.get(m);
            visibleByName.computeIfAbsent(s.simpleName, k -> new LinkedHashSet<>()).add(s);
          }
        }
      }
      for (Set<MethodSig> set : visibleByName.values()) {
        if (set.size() > 1) {
          for (MethodSig a : set) {
            for (MethodSig b : set) {
              if (a != b) {
                coexistingSameName.get(a).add(b);
              }
            }
          }
        }
      }
    }

    // Assign Go method names
    int slotSeq = 0;
    for (Map.Entry<ExecutableElement, MethodSig> e : rootSigs.entrySet()) {
      ExecutableElement rep = e.getKey();
      MethodSig sig = e.getValue();
      slotSeq++;

      // Special canonical JLS Object methods
      if (isPublicObjectMethod(rep)) {
        sig.goName =
            switch (sig.simpleName) {
              case "toString" -> "String";
              case "hashCode" -> "HashCode";
              case "equals" -> "Equals";
              default -> GoNames.exported(sig.simpleName);
            };
        for (int i = 0; i < sig.paramWidenedToAny.length; i++) {
          sig.paramWidenedToAny[i] = true;
          sig.canonParamTypes.set(i, objectElement.asType());
        }
        continue;
      }

      String baseExp = GoNames.exported(sig.simpleName);
      if (RESERVED_METHOD_NAMES.contains(baseExp)
          || baseExp.startsWith("Implements_")
          || baseExp.startsWith("FieldGet_")
          || baseExp.startsWith("FieldSet_")
          || baseExp.startsWith("Init_")) {
        baseExp = "M_" + baseExp;
      }

      // Private instance methods that aren't overridden still must not collide on their declaring type
      if (rep.getModifiers().contains(Modifier.PRIVATE)) {
        TypeElement owner = (TypeElement) rep.getEnclosingElement();
        sig.goName =
            "priv_" + GoNames.cleanChars(owner.getSimpleName().toString()) + "_" + sig.simpleName + "_" + slotSeq;
        continue;
      }

      Set<MethodSig> peers = coexistingSameName.get(sig);
      if (peers.isEmpty()) {
        sig.goName = baseExp;
      } else {
        boolean sameArityCollision = false;
        int arity = sig.canonParamTypes.size();
        for (MethodSig peer : peers) {
          if (peer.canonParamTypes.size() == arity && !isPublicObjectMethodFromSig(peer)) {
            sameArityCollision = true;
            break;
          }
        }
        if (!sameArityCollision) {
          sig.goName = baseExp + "_" + arity;
        } else {
          String suffix = paramListTag(sig.canonParamTypes);
          String cand = baseExp + "_" + arity + "_" + suffix;
          // Ensure uniqueness against peers
          for (MethodSig peer : peers) {
            if (cand.equals(peer.goName)) {
              cand = cand + "_" + slotSeq;
              break;
            }
          }
          sig.goName = cand;
        }
      }
    }
  }

  private void registerExternalExec(ExecutableElement ee) {
    TypeElement owner = (TypeElement) ee.getEnclosingElement();
    if (isSyntheticArrayElement(owner)) {
      return;
    }
    recordExternalType(owner);
    referencedExternalExecs.add(ee);
    if (ee.getKind() == ElementKind.CONSTRUCTOR) {
      List<ExecutableElement> ctors = new ArrayList<>();
      for (Element m : owner.getEnclosedElements()) {
        if (m instanceof ExecutableElement c && c.getKind() == ElementKind.CONSTRUCTOR) {
          ctors.add(c);
        }
      }
      assignConstructorNames(owner, ctors);
    } else if (ee.getModifiers().contains(Modifier.STATIC)) {
      List<ExecutableElement> statics = new ArrayList<>();
      for (Element m : owner.getEnclosedElements()) {
        if (m instanceof ExecutableElement s
            && s.getKind() == ElementKind.METHOD
            && s.getModifiers().contains(Modifier.STATIC)) {
          statics.add(s);
        }
      }
      assignStaticMethodNames(owner, statics);
    } else {
      int arity = ee.getParameters().size();
      List<TypeMirror> params = new ArrayList<>(arity);
      for (VariableElement p : ee.getParameters()) {
        params.add(safeErasure(p.asType()));
      }
      TypeMirror ret =
          ee.getReturnType().getKind() == TypeKind.VOID ? null : safeErasure(ee.getReturnType());
      MethodSig sig = new MethodSig(ee.getSimpleName().toString(), arity, null, ret, params);
      sig.goName = GoNames.exported(sig.simpleName);
      instanceMethodSigs.put(ee, sig);
    }
  }

  private ExecutableElement findCanonicalGenericRootMember(List<ExecutableElement> members) {
    if (members.size() == 1) {
      ExecutableElement only = members.get(0);
      if (only.getEnclosingElement() instanceof TypeElement ro && hasEffectiveTypeParameters(ro)) {
        return only;
      }
      return null;
    }
    Set<TypeElement> owners = new LinkedHashSet<>();
    for (ExecutableElement m : members) {
      if (!(m.getEnclosingElement() instanceof TypeElement te)) {
        return null;
      }
      owners.add(te);
    }
    for (ExecutableElement m : members) {
      TypeElement cand = (TypeElement) m.getEnclosingElement();
      if (!isNativeGenericClass(cand)) {
        continue;
      }
      boolean isAncestorOfAll = true;
      TypeMirror candErased = safeErasure(cand.asType());
      for (TypeElement other : owners) {
        if (!other.equals(cand) && !types.isSubtype(safeErasure(other.asType()), candErased)) {
          isAncestorOfAll = false;
          break;
        }
      }
      if (isAncestorOfAll) {
        return m;
      }
    }
    return null;
  }

  private boolean mentionsGoTypeParam(TypeMirror tm, TypeParameterElement targetTpe) {
    if (tm == null) {
      return false;
    }
    return switch (tm.getKind()) {
      case TYPEVAR -> targetTpe.equals(((TypeVariable) tm).asElement());
      case ARRAY -> mentionsGoTypeParam(((ArrayType) tm).getComponentType(), targetTpe);
      case WILDCARD -> {
        WildcardType wt = (WildcardType) tm;
        yield wt.getExtendsBound() != null && mentionsGoTypeParam(wt.getExtendsBound(), targetTpe);
      }
      case DECLARED -> {
        DeclaredType dt = (DeclaredType) tm;
        if (dt.asElement() instanceof TypeElement dte && isNativeGenericClass(dte)) {
          for (TypeMirror ta : dt.getTypeArguments()) {
            if (mentionsGoTypeParam(ta, targetTpe)) {
              yield true;
            }
          }
        }
        yield false;
      }
      default -> false;
    };
  }

  private boolean isPublicObjectMethodFromSig(MethodSig s) {
    if (s.simpleName.equals("toString") && s.canonParamTypes.isEmpty()) return true;
    if (s.simpleName.equals("hashCode") && s.canonParamTypes.isEmpty()) return true;
    if (s.simpleName.equals("equals") && s.canonParamTypes.size() == 1) {
      return canonicalGoTypeKey(s.canonParamTypes.get(0)).equals("any");
    }
    return false;
  }

  private void assignConstructorNames(TypeElement te, List<ExecutableElement> ctors) {
    String typeBase = getTypeGoBaseName(te);
    if (ctors.size() <= 1) {
      for (ExecutableElement c : ctors) {
        constructorGoNames.put(c, "New" + typeBase);
      }
      return;
    }
    Map<Integer, List<ExecutableElement>> byArity = new HashMap<>();
    for (ExecutableElement c : ctors) {
      byArity.computeIfAbsent(c.getParameters().size(), k -> new ArrayList<>()).add(c);
    }
    Set<String> used = new HashSet<>();
    ctors.sort(Comparator.comparing(this::execSortKey));
    for (ExecutableElement c : ctors) {
      int arity = c.getParameters().size();
      String cand;
      if (byArity.get(arity).size() == 1) {
        cand = arity == 0 ? ("New" + typeBase) : ("New" + typeBase + "_" + arity);
      } else {
        List<TypeMirror> pts = new ArrayList<>(arity);
        for (VariableElement p : c.getParameters()) {
          pts.add(safeErasure(p.asType()));
        }
        cand = "New" + typeBase + "_" + arity + "_" + paramListTag(pts);
      }
      int idx = 2;
      String unique = cand;
      while (!used.add(unique)) {
        unique = cand + "_" + idx++;
      }
      constructorGoNames.put(c, unique);
    }
  }

  private String staticMethodPrefix(TypeElement te, String exp) {
    String typeBase = getTypeGoBaseName(te);
    String cand = typeBase + "_" + exp;
    for (Element enc : te.getEnclosedElements()) {
      if (enc instanceof TypeElement sub) {
        String subClean = GoNames.cleanChars(sub.getSimpleName().toString());
        if (subClean.equals(exp) || (subClean + "Base").equals(exp)) {
          return typeBase + "_S_" + exp;
        }
      }
    }
    if (isSourceType(te)) {
      String pkg = planner.getGoPkg(te);
      for (TypeElement peer : planner.getAllSourceTypes()) {
        if (pkg.equals(planner.getGoPkg(peer))) {
          String peerName = planner.getGoTypeName(peer);
          if (peerName.equals(cand) || (peerName + "Base").equals(cand)) {
            return typeBase + "_S_" + exp;
          }
        }
      }
    } else {
      for (TypeElement peer : referencedExternalTypes) {
        String peerName = runtime.extGoTypeName(peer);
        if (peerName.equals(cand) || (peerName + "Base").equals(cand)) {
          return typeBase + "_S_" + exp;
        }
      }
    }
    return cand;
  }

  private void assignStaticMethodNames(TypeElement te, List<ExecutableElement> statics) {
    String typeBase = getTypeGoBaseName(te);
    boolean isSourceEnum = isSourceType(te) && te.getKind() == ElementKind.ENUM;
    Map<String, List<ExecutableElement>> bySimple = new LinkedHashMap<>();
    for (ExecutableElement m : statics) {
      bySimple.computeIfAbsent(m.getSimpleName().toString(), k -> new ArrayList<>()).add(m);
    }
    for (Map.Entry<String, List<ExecutableElement>> e : bySimple.entrySet()) {
      String exp = GoNames.exported(e.getKey());
      String baseCand = staticMethodPrefix(te, exp);
      List<ExecutableElement> group = e.getValue();
      if (group.size() == 1) {
        staticMethodGoNames.put(group.get(0), baseCand);
      } else {
        Map<Integer, Integer> arityCounts = new HashMap<>();
        for (ExecutableElement m : group) {
          arityCounts.merge(m.getParameters().size(), 1, Integer::sum);
        }
        Set<String> used = new HashSet<>();
        group.sort(Comparator.comparing(this::execSortKey));
        for (ExecutableElement m : group) {
          int arity = m.getParameters().size();
          String cand;
          if (arityCounts.get(arity) == 1) {
            cand = baseCand + "_" + arity;
          } else {
            List<TypeMirror> pts = new ArrayList<>(arity);
            for (VariableElement p : m.getParameters()) {
              pts.add(safeErasure(p.asType()));
            }
            cand = baseCand + "_" + arity + "_" + paramListTag(pts);
          }
          int idx = 2;
          String unique = cand;
          while (!used.add(unique)) {
            unique = cand + "_" + idx++;
          }
          staticMethodGoNames.put(m, unique);
        }
      }
    }
    if (isSourceEnum) {
      for (ExecutableElement m : statics) {
        String mn = m.getSimpleName().toString();
        if (mn.equals("values") && m.getParameters().isEmpty()) {
          staticMethodGoNames.put(m, typeBase + "_Values");
        } else if (mn.equals("valueOf")
            && m.getParameters().size() == 1
            && isStringType(m.getParameters().get(0).asType())) {
          staticMethodGoNames.put(m, typeBase + "_ValueOf");
        }
      }
    }
  }

  private String execSortKey(ExecutableElement ee) {
    StringBuilder sb = new StringBuilder();
    sb.append(ee.getParameters().size()).append(':');
    for (VariableElement p : ee.getParameters()) {
      sb.append(safeErasure(p.asType()).toString()).append(',');
    }
    return sb.toString();
  }

  private String paramListTag(List<TypeMirror> params) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < params.size(); i++) {
      if (i > 0) {
        sb.append('_');
      }
      String key = canonicalGoTypeKey(params.get(i));
      int dot = key.lastIndexOf('.');
      String shortName = dot >= 0 ? key.substring(dot + 1) : key;
      sb.append(GoNames.cleanChars(shortName));
    }
    return sb.toString();
  }

  private boolean sameErasedJavaParams(ExecutableElement a, ExecutableElement b) {
    if (a.getParameters().size() != b.getParameters().size()) {
      return false;
    }
    for (int i = 0; i < a.getParameters().size(); i++) {
      TypeMirror pa = safeErasure(a.getParameters().get(i).asType());
      TypeMirror pb = safeErasure(b.getParameters().get(i).asType());
      if (!types.isSameType(pa, pb)) {
        return false;
      }
    }
    return true;
  }

  private static ExecutableElement ufFind(
      Map<ExecutableElement, ExecutableElement> parent, ExecutableElement x) {
    ExecutableElement p = parent.get(x);
    if (p == x) {
      return x;
    }
    ExecutableElement r = ufFind(parent, p);
    parent.put(x, r);
    return r;
  }

  private static void ufUnion(
      Map<ExecutableElement, ExecutableElement> parent,
      ExecutableElement a,
      ExecutableElement b) {
    ExecutableElement ra = ufFind(parent, a);
    ExecutableElement rb = ufFind(parent, b);
    if (ra != rb) {
      parent.put(rb, ra);
    }
  }

  /** Computes the instantiated ExecutableType at a call/reference site when receiver has DeclaredType. */
  public ExecutableType resolveExecutableType(TypeMirror receiverType, ExecutableElement ee) {
    if (receiverType != null) {
      TypeMirror actual = receiverType;
      if (actual instanceof WildcardType wt && wt.getExtendsBound() != null) {
        actual = wt.getExtendsBound();
      }
      if (actual instanceof TypeVariable tv && tv.getUpperBound() != null) {
        actual = tv.getUpperBound();
      }
      if (actual instanceof DeclaredType dt) {
        try {
          Element asMemberType = dt.asElement();
          if (asMemberType instanceof TypeElement recvTe
              && types.isSubtype(
                  safeErasure(recvTe.asType()), safeErasure(ee.getEnclosingElement().asType()))) {
            return (ExecutableType) types.asMemberOf(dt, ee);
          }
        } catch (IllegalArgumentException ignored) {
          // Fall back to ee.asType()
        }
      }
    }
    return (ExecutableType) ee.asType();
  }
}
