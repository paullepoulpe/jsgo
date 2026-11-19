package j2go;

import java.util.Set;
import javax.lang.model.element.ExecutableElement;
import javax.lang.model.element.TypeElement;
import javax.lang.model.element.VariableElement;

/**
 * Clean-room runtime bindings for {@code java.util.*}, {@code java.util.stream.*}, {@code
 * java.util.function.*}, {@code java.util.concurrent.*}, {@code com.google.common.collect.*}, and
 * {@code com.google.common.cache.*}.
 *
 * <p>Guarantees 100% deterministic iteration across all hash, sorted, enum, and identity maps,
 * sets, multimaps, multisets, and tables, with full Java {@code equals}/{@code hashCode} and live
 * view semantics.
 */
public final class RuntimeCollectionsBindings {

  private RuntimeCollectionsBindings() {}

  private static final Set<String> CONCRETE_LIST_TYPES =
      Set.of(
          "java.util.ArrayList",
          "java.util.LinkedList",
          "java.util.concurrent.CopyOnWriteArrayList",
          "java.util.concurrent.LinkedBlockingQueue",
          "com.google.common.collect.ImmutableList");

  private static final Set<String> CONCRETE_SET_TYPES =
      Set.of(
          "java.util.HashSet",
          "java.util.LinkedHashSet",
          "java.util.TreeSet",
          "java.util.EnumSet",
          "com.google.common.collect.ImmutableSet",
          "com.google.common.collect.ImmutableSet.CachingAsList",
          "com.google.common.collect.ImmutableSortedSet");

  private static final Set<String> CONCRETE_MAP_TYPES =
      Set.of(
          "java.util.HashMap",
          "java.util.LinkedHashMap",
          "java.util.TreeMap",
          "java.util.EnumMap",
          "java.util.IdentityHashMap",
          "com.google.common.collect.ImmutableMap",
          "com.google.common.collect.ImmutableSortedMap");

  private static final Set<String> MULTIMAP_TYPES =
      Set.of(
          "com.google.common.collect.AbstractMultimap",
          "com.google.common.collect.AbstractMapBasedMultimap",
          "com.google.common.collect.AbstractListMultimap",
          "com.google.common.collect.AbstractSetMultimap",
          "com.google.common.collect.AbstractSortedSetMultimap",
          "com.google.common.collect.AbstractSortedKeySortedSetMultimap",
          "com.google.common.collect.ArrayListMultimap",
          "com.google.common.collect.ArrayListMultimapGwtSerializationDependencies",
          "com.google.common.collect.LinkedListMultimap",
          "com.google.common.collect.HashMultimap",
          "com.google.common.collect.HashMultimapGwtSerializationDependencies",
          "com.google.common.collect.LinkedHashMultimap",
          "com.google.common.collect.LinkedHashMultimapGwtSerializationDependencies",
          "com.google.common.collect.TreeMultimap",
          "com.google.common.collect.BaseImmutableMultimap",
          "com.google.common.collect.ImmutableMultimap",
          "com.google.common.collect.ImmutableListMultimap",
          "com.google.common.collect.ImmutableSetMultimap");

  private static final Set<String> MULTISET_TYPES =
      Set.of(
          "com.google.common.collect.AbstractMultiset",
          "com.google.common.collect.AbstractMapBasedMultiset",
          "com.google.common.collect.AbstractSortedMultiset",
          "com.google.common.collect.HashMultiset",
          "com.google.common.collect.LinkedHashMultiset",
          "com.google.common.collect.TreeMultiset",
          "com.google.common.collect.ImmutableMultiset",
          "com.google.common.collect.ImmutableMultisetGwtSerializationDependencies");

  private static final Set<String> TABLE_TYPES =
      Set.of(
          "com.google.common.collect.AbstractTable",
          "com.google.common.collect.StandardTable",
          "com.google.common.collect.HashBasedTable",
          "com.google.common.collect.ImmutableTable");

  public static boolean isConcreteList(String qn) {
    return CONCRETE_LIST_TYPES.contains(qn);
  }

  public static boolean isConcreteSet(String qn) {
    return CONCRETE_SET_TYPES.contains(qn);
  }

  public static boolean isConcreteMap(String qn) {
    return CONCRETE_MAP_TYPES.contains(qn);
  }

  public static boolean skipCoreField(String qn) {
    return qn.equals("java.util.BitSet");
  }

  public static void emitExtraStructFields(TypeElement te, StringBuilder out) {
    if (te.getQualifiedName().contentEquals("java.util.BitSet")) {
      out.append("\twords []uint64 `gson:\"-\"`\n");
      out.append("\tbuf [8]uint64 `gson:\"-\"`\n");
      out.append("\thi [8]int32 `gson:\"-\"`\n");
      out.append("\tinitW int32 `gson:\"-\"`\n");
      out.append("\thiLen int32 `gson:\"-\"`\n");
    }
  }

  private static String selfIfaceFor(String qn, String extName) {
    return switch (qn) {
      case "java.util.ArrayList",
          "java.util.EnumSet",
          "java.util.HashSet",
          "com.google.common.collect.ImmutableCollection",
          "com.google.common.collect.ImmutableList",
          "com.google.common.collect.ImmutableSet",
          "com.google.common.collect.ImmutableSet.CachingAsList",
          "com.google.common.collect.ImmutableSortedSet" ->
          "Ext_java_util_AbstractCollection";
      case "java.util.HashMap", "java.util.EnumMap" -> "Ext_java_util_AbstractMap";
      case "com.google.common.collect.ImmutableSortedMap" ->
          "Ext_com_google_common_collect_ImmutableMap";
      default -> extName;
    };
  }

  public static void emitFastMapSetMethods(
      RuntimeBindings rb,
      TypeElement te,
      String structName,
      boolean subclassed,
      String extName,
      StringBuilder out) {
    String qn = te.getQualifiedName().toString();
    String selfIface = selfIfaceFor(qn, extName);
    if (qn.equals("com.google.common.collect.Ordering")) {
      out.append("func (this_ *")
          .append(structName)
          .append(") isNatOrder_() bool { return this_ != nil && this_.Core_ == jNatOrderCore{} }\n");
      return;
    }

    if (qn.equals("java.util.ArrayDeque")) {
      out.append("func (this_ *").append(structName).append(") fastDequeCore_() *JDequeCore {\n")
          .append("\tif this_ == nil { return nil }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return nil }\n");
      }
      out.append("\treturn this_.dc()\n}\n");
      return;
    }

    if (qn.equals("java.util.LinkedList")) {
      out.append("func (this_ *").append(structName).append(") fastLinkedDequeCore_() *JListCore {\n")
          .append("\tif this_ == nil { return nil }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return nil }\n");
      }
      out.append("\treturn this_.lc()\n}\n");
      return;
    }

    if (qn.equals("java.util.ArrayList")) {
      String listGetGo = rb.findExtMethodGoName("java.util.List", "get", 1);
      out.append("func (this_ *").append(structName).append(") IterSlice_() []any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return iterToSliceDirect(this_.Self_) }\n");
      }
      out.append("\treturn this_.lc().Items\n}\n");
      out.append("func (this_ *").append(structName).append(") ListGetStr_(i int32) string {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return CastString(this_.Self_.(Ext_java_util_List).").append(listGetGo).append("(i)) }\n");
      }
      out.append("\treturn CastString(this_.lc().Get(i))\n}\n");
      return;
    }

    if (qn.equals("com.google.common.collect.ImmutableList")) {
      out.append("func (this_ *").append(structName).append(") IterSlice_() []any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n")
          .append("\treturn this_.items_\n}\n")
          .append("func (this_ *").append(structName).append(") ListGetStr_(i int32) string {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n")
          .append("\tCheckElementIndex(i, int32(len(this_.items_)))\n")
          .append("\treturn CastString(this_.items_[i])\n}\n")
          .append("func (this_ *").append(structName).append(") SetContainsStr_(k string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n")
          .append("\tif k == NullString { return false }\n")
          .append("\tfor _, e := range this_.items_ {\n\t\tif es, ok := e.(string); ok && es == k { return true }\n\t}\n")
          .append("\treturn false\n}\n")
          .append("func (this_ *").append(structName).append(") SetAddStr_(string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n")
          .append("func (this_ *").append(structName).append(") SetRemoveStr_(string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n")
          .append("func (this_ *").append(structName).append(") SetContainsNode_(u uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\tif u == 0 { return false }\n")
          .append("\tfor _, e := range this_.items_ {\n\t\tif singletonElemMatchesNode(e, u) { return true }\n\t}\n")
          .append("\treturn false\n}\n")
          .append("func (this_ *").append(structName).append(") SetAddNode_(uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n")
          .append("func (this_ *").append(structName).append(") SetRemoveNode_(uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      return;
    }

    if (qn.equals("com.google.common.collect.ImmutableCollection")) {
      out.append("func (this_ *").append(structName).append(") IterSlice_() []any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n")
          .append("\tif lc, ok := this_.Core_.(*JListCore); ok { return lc.Items }\n")
          .append("\tif sc, ok := this_.Core_.(*JSetCore); ok { return sc.LiveKeys() }\n")
          .append("\tif mc, ok := this_.Core_.(*JMapCore); ok { return mc.LiveKeys() }\n")
          .append("\tif ms, ok := this_.Core_.(*jMultisetCore); ok { return multisetOp_toSlice(ms) }\n")
          .append("\treturn iterToSliceDirect(this_.Self_)\n}\n");
      return;
    }

    if (qn.equals("com.google.common.collect.ImmutableSet")
        || qn.equals("com.google.common.collect.ImmutableSet.CachingAsList")) {
      String colContainsGo = rb.findExtMethodGoName("java.util.Collection", "contains", 1);
      String colAddGo = rb.findExtMethodGoName("java.util.Collection", "add", 1);
      String colRemoveGo = rb.findExtMethodGoName("java.util.Collection", "remove", 1);

      if (qn.equals("com.google.common.collect.ImmutableSet")) {
        out.append("func (this_ *").append(structName).append(") IterStrSlice_() ([]string, bool) {\n")
            .append("\tif this_ == nil { ThrowNPE() }\n")
            .append("\tif this_.strs_ != nil { return this_.strs_, true }\n")
            .append("\treturn nil, false\n}\n")
            .append("func (this_ *").append(structName).append(") IterSlice_() []any {\n")
            .append("\tif this_ == nil { ThrowNPE() }\n")
            .append("\tif this_.elems_ != nil { return this_.elems_ }\n")
            .append("\tif this_.strs_ != nil { return strsToAnySlice(this_.strs_) }\n")
            .append("\treturn iterToSliceDirect(this_.Self_)\n}\n");
      }

      out.append("func (this_ *").append(structName).append(") SetContainsStr_(k string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colContainsGo).append("(BoxString(k)) }\n");
      }
      out.append("\treturn this_.immSetContainsStr(k)\n}\n");

      out.append("func (this_ *").append(structName).append(") SetAddStr_(k string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colAddGo).append("(BoxString(k)) }\n");
      }
      out.append("\t_ = k\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");

      out.append("func (this_ *").append(structName).append(") SetRemoveStr_(k string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colRemoveGo).append("(BoxString(k)) }\n");
      }
      out.append("\t_ = k\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");

      out.append("func (this_ *").append(structName).append(") SetContainsNode_(n uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colContainsGo).append("(boxNodeById(n)) }\n");
      }
      out.append("\treturn this_.immSetContainsNode(n)\n}\n");

      out.append("func (this_ *").append(structName).append(") SetAddNode_(n uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colAddGo).append("(boxNodeById(n)) }\n");
      }
      out.append("\t_ = n\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");

      out.append("func (this_ *").append(structName).append(") SetRemoveNode_(n uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colRemoveGo).append("(boxNodeById(n)) }\n");
      }
      out.append("\t_ = n\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      return;
    }

    if (qn.equals("java.util.EnumSet")) {
      String colContainsGo = rb.findExtMethodGoName("java.util.Collection", "contains", 1);
      String colAddGo = rb.findExtMethodGoName("java.util.Collection", "add", 1);
      String colRemoveGo = rb.findExtMethodGoName("java.util.Collection", "remove", 1);

      out.append("func (this_ *").append(structName).append(") fastEnumSet_() *jEnumSetCore {\n")
          .append("\tif this_ == nil { return nil }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return nil }\n");
      }
      out.append("\treturn &this_.inlineEnumSet\n}\n");

      out.append("func (this_ *").append(structName).append(") SetContainsStr_(k string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colContainsGo).append("(BoxString(k)) }\n");
      }
      out.append("\t_ = k\n\treturn false\n}\n");

      out.append("func (this_ *").append(structName).append(") SetAddStr_(k string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colAddGo).append("(BoxString(k)) }\n");
      }
      out.append("\tCheckNotNullString(k)\n\tpanic(javalang.NewClassCastException(\"java.lang.String cannot be cast to java.lang.Enum\"))\n}\n");

      out.append("func (this_ *").append(structName).append(") SetRemoveStr_(k string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colRemoveGo).append("(BoxString(k)) }\n");
      }
      out.append("\t_ = k\n\treturn false\n}\n");

      out.append("func (this_ *").append(structName).append(") SetContainsNode_(n uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colContainsGo).append("(boxNodeById(n)) }\n");
      }
      out.append("\t_ = n\n\treturn false\n}\n");

      out.append("func (this_ *").append(structName).append(") SetAddNode_(n uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colAddGo).append("(boxNodeById(n)) }\n");
      }
      out.append("\tif n == 0 { ThrowNPE() }\n\tpanic(javalang.NewClassCastException(\"Node cannot be cast to java.lang.Enum\"))\n}\n");

      out.append("func (this_ *").append(structName).append(") SetRemoveNode_(n uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colRemoveGo).append("(boxNodeById(n)) }\n");
      }
      out.append("\t_ = n\n\treturn false\n}\n");
      return;
    }

    if (qn.equals("java.util.HashSet") || qn.equals("java.util.LinkedHashSet")) {
      String colContainsGo = rb.findExtMethodGoName("java.util.Collection", "contains", 1);
      String colAddGo = rb.findExtMethodGoName("java.util.Collection", "add", 1);
      String colRemoveGo = rb.findExtMethodGoName("java.util.Collection", "remove", 1);
      String mapContainsKeyGo = rb.findExtMethodGoName("java.util.Map", "containsKey", 1);
      String mapPutGo = rb.findExtMethodGoName("java.util.Map", "put", 2);
      String mapRemoveGo = rb.findExtMethodGoName("java.util.Map", "remove", 1);

      out.append("func (this_ *").append(structName).append(") fastSetJCore_() *JSetCore {\n")
          .append("\tif this_ == nil || this_.Core_ != nil { return nil }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return nil }\n");
      }
      out.append("\treturn &this_.inlineSet\n}\n");

      out.append("func (this_ *").append(structName).append(") SetContainsStr_(k string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colContainsGo).append("(BoxString(k)) }\n");
      }
      out.append("\tif mb := this_.mapBack(); mb != nil { return mb.").append(mapContainsKeyGo).append("(BoxString(k)) }\n")
          .append("\treturn this_.inlineSet.ContainsStr(k)\n}\n");

      out.append("func (this_ *").append(structName).append(") SetAddStr_(k string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colAddGo).append("(BoxString(k)) }\n");
      }
      out.append("\tif mb := this_.mapBack(); mb != nil { return mb.").append(mapPutGo).append("(BoxString(k), ValueOfBoolean(true)) == nil }\n")
          .append("\treturn this_.inlineSet.AddStr(k)\n}\n");

      out.append("func (this_ *").append(structName).append(") SetRemoveStr_(k string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colRemoveGo).append("(BoxString(k)) }\n");
      }
      out.append("\tif mb := this_.mapBack(); mb != nil { return mb.").append(mapRemoveGo).append("(BoxString(k)) != nil }\n")
          .append("\treturn this_.inlineSet.RemoveStr(k)\n}\n");

      out.append("func (this_ *").append(structName).append(") SetContainsNode_(n uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colContainsGo).append("(boxNodeById(n)) }\n");
      }
      out.append("\tif mb := this_.mapBack(); mb != nil { return mb.").append(mapContainsKeyGo).append("(boxNodeById(n)) }\n")
          .append("\treturn this_.inlineSet.ContainsNode(n)\n}\n");

      out.append("func (this_ *").append(structName).append(") SetAddNode_(n uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colAddGo).append("(boxNodeById(n)) }\n");
      }
      out.append("\tif mb := this_.mapBack(); mb != nil { return mb.").append(mapPutGo).append("(boxNodeById(n), ValueOfBoolean(true)) == nil }\n")
          .append("\treturn this_.inlineSet.AddNode(n)\n}\n");

      out.append("func (this_ *").append(structName).append(") SetRemoveNode_(n uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colRemoveGo).append("(boxNodeById(n)) }\n");
      }
      out.append("\tif mb := this_.mapBack(); mb != nil { return mb.").append(mapRemoveGo).append("(boxNodeById(n)) != nil }\n")
          .append("\treturn this_.inlineSet.RemoveNode(n)\n}\n");
      return;
    }

    if (CONCRETE_SET_TYPES.contains(qn)) {
      boolean isImm = qn.startsWith("com.google.");
      String colContainsGo = rb.findExtMethodGoName("java.util.Collection", "contains", 1);
      String colAddGo = rb.findExtMethodGoName("java.util.Collection", "add", 1);
      String colRemoveGo = rb.findExtMethodGoName("java.util.Collection", "remove", 1);

      out.append("func (this_ *").append(structName).append(") fastSetCore_() *JMapCore {\n")
          .append("\tif this_ == nil { return nil }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return nil }\n");
      }
      out.append("\treturn this_.mc()\n}\n");

      out.append("func (this_ *").append(structName).append(") SetContainsStr_(k string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colContainsGo).append("(BoxString(k)) }\n");
      }
      out.append("\treturn this_.mc().ContainsStr(k)\n}\n");

      out.append("func (this_ *").append(structName).append(") SetAddStr_(k string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colAddGo).append("(BoxString(k)) }\n");
      }
      if (isImm) {
        out.append("\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      } else {
        out.append("\treturn this_.mc().AddKeyStr(k)\n}\n");
      }

      out.append("func (this_ *").append(structName).append(") SetRemoveStr_(k string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colRemoveGo).append("(BoxString(k)) }\n");
      }
      if (isImm) {
        out.append("\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      } else {
        out.append("\treturn this_.mc().RemoveKeyStr(k)\n}\n");
      }

      out.append("func (this_ *").append(structName).append(") SetContainsNode_(n uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colContainsGo).append("(boxNodeById(n)) }\n");
      }
      out.append("\treturn this_.mc().ContainsNode(n)\n}\n");

      out.append("func (this_ *").append(structName).append(") SetAddNode_(n uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colAddGo).append("(boxNodeById(n)) }\n");
      }
      if (isImm) {
        out.append("\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      } else {
        out.append("\treturn this_.mc().AddKeyNode(n)\n}\n");
      }

      out.append("func (this_ *").append(structName).append(") SetRemoveNode_(n uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(colRemoveGo).append("(boxNodeById(n)) }\n");
      }
      if (isImm) {
        out.append("\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      } else {
        out.append("\treturn this_.mc().RemoveKeyNode(n)\n}\n");
      }
      return;
    }

    if (qn.equals("java.util.EnumMap")) {
      String mapGetGo = rb.findExtMethodGoName("java.util.Map", "get", 1);
      String mapGetOrDefGo = rb.findExtMethodGoName("java.util.Map", "getOrDefault", 2);
      String mapContainsKeyGo = rb.findExtMethodGoName("java.util.Map", "containsKey", 1);
      String mapPutGo = rb.findExtMethodGoName("java.util.Map", "put", 2);
      String mapRemoveGo = rb.findExtMethodGoName("java.util.Map", "remove", 1);
      String mapComputeIfAbsentGo = rb.findExtMethodGoName("java.util.Map", "computeIfAbsent", 2);

      out.append("func (this_ *").append(structName).append(") fastEnumMap_() *jEnumMapCore {\n")
          .append("\tif this_ == nil { return nil }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return nil }\n");
      }
      out.append("\treturn &this_.inlineEnumMap\n}\n");

      out.append("func (this_ *").append(structName).append(") MapGetStr_(k string) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapGetGo).append("(BoxString(k)) }\n");
      }
      out.append("\t_ = k\n\treturn nil\n}\n");

      out.append("func (this_ *").append(structName).append(") MapGetOrDefaultStr_(k string, def any) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapGetOrDefGo).append("(BoxString(k), ToAny(def)) }\n");
      }
      out.append("\t_ = k\n\treturn ToAny(def)\n}\n");

      out.append("func (this_ *").append(structName).append(") MapContainsStr_(k string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapContainsKeyGo).append("(BoxString(k)) }\n");
      }
      out.append("\t_ = k\n\treturn false\n}\n");

      out.append("func (this_ *").append(structName).append(") MapPutStr_(k string, v any) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapPutGo).append("(BoxString(k), ToAny(v)) }\n");
      }
      out.append("\tCheckNotNullString(k)\n\t_ = v\n\tpanic(javalang.NewClassCastException(\"java.lang.String cannot be cast to java.lang.Enum\"))\n}\n");

      out.append("func (this_ *").append(structName).append(") MapRemoveStr_(k string) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapRemoveGo).append("(BoxString(k)) }\n");
      }
      out.append("\t_ = k\n\treturn nil\n}\n");

      out.append("func (this_ *").append(structName).append(") MapGetNode_(n uint32) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapGetGo).append("(boxNodeById(n)) }\n");
      }
      out.append("\t_ = n\n\treturn nil\n}\n");

      out.append("func (this_ *").append(structName).append(") MapGetOrDefaultNode_(n uint32, def any) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapGetOrDefGo).append("(boxNodeById(n), ToAny(def)) }\n");
      }
      out.append("\t_ = n\n\treturn ToAny(def)\n}\n");

      out.append("func (this_ *").append(structName).append(") MapContainsNode_(n uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapContainsKeyGo).append("(boxNodeById(n)) }\n");
      }
      out.append("\t_ = n\n\treturn false\n}\n");

      out.append("func (this_ *").append(structName).append(") MapPutNode_(n uint32, v any) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapPutGo).append("(boxNodeById(n), ToAny(v)) }\n");
      }
      out.append("\tif n == 0 { ThrowNPE() }\n\t_ = v\n\tpanic(javalang.NewClassCastException(\"Node cannot be cast to java.lang.Enum\"))\n}\n");

      out.append("func (this_ *").append(structName).append(") MapRemoveNode_(n uint32) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapRemoveGo).append("(boxNodeById(n)) }\n");
      }
      out.append("\t_ = n\n\treturn nil\n}\n");

      out.append("func (this_ *").append(structName).append(") MapComputeIfAbsentStr_(k string, fn Ext_java_util_function_Function) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapComputeIfAbsentGo).append("(BoxString(k), fn) }\n");
      }
      out.append("\tCheckNotNull(fn)\n\tCheckNotNullString(k)\n\tpanic(javalang.NewClassCastException(\"java.lang.String cannot be cast to java.lang.Enum\"))\n}\n");

      out.append("func (this_ *").append(structName).append(") MapComputeIfAbsentNode_(n uint32, fn Ext_java_util_function_Function) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapComputeIfAbsentGo).append("(boxNodeById(n), fn) }\n");
      }
      out.append("\tCheckNotNull(fn)\n\tif n == 0 { ThrowNPE() }\n\tpanic(javalang.NewClassCastException(\"Node cannot be cast to java.lang.Enum\"))\n}\n");

      out.append("func (this_ *").append(structName).append(") MapTryGetComputeStr_(k string) (any, bool) {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\t_ = k\n\treturn nil, false\n}\n");

      out.append("func (this_ *").append(structName).append(") MapTryGetComputeNode_(n uint32) (any, bool) {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\t_ = n\n\treturn nil, false\n}\n");

      out.append("func (this_ *").append(structName).append(") MapTryGetComputeAny_(k any) (any, bool) {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return nil, false }\n");
      }
      out.append("\treturn this_.inlineEnumMap.tryGetNonNull(k)\n}\n");
      return;
    }

    if (CONCRETE_MAP_TYPES.contains(qn)) {
      boolean isImm = qn.startsWith("com.google.");
      String mapGetGo = rb.findExtMethodGoName("java.util.Map", "get", 1);
      String mapGetOrDefGo = rb.findExtMethodGoName("java.util.Map", "getOrDefault", 2);
      String mapContainsKeyGo = rb.findExtMethodGoName("java.util.Map", "containsKey", 1);
      String mapPutGo = rb.findExtMethodGoName("java.util.Map", "put", 2);
      String mapRemoveGo = rb.findExtMethodGoName("java.util.Map", "remove", 1);
      String mapComputeIfAbsentGo = rb.findExtMethodGoName("java.util.Map", "computeIfAbsent", 2);

      out.append("func (this_ *").append(structName).append(") fastMapCore_() *JMapCore {\n")
          .append("\tif this_ == nil { return nil }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return nil }\n");
      }
      out.append("\treturn this_.mc()\n}\n");

      out.append("func (this_ *").append(structName).append(") MapGetStr_(k string) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapGetGo).append("(BoxString(k)) }\n");
      }
      out.append("\treturn this_.mc().GetStr(k)\n}\n");

      out.append("func (this_ *").append(structName).append(") MapGetOrDefaultStr_(k string, def any) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapGetOrDefGo).append("(BoxString(k), ToAny(def)) }\n");
      }
      out.append("\treturn this_.mc().GetOrDefaultStr(k, def)\n}\n");

      out.append("func (this_ *").append(structName).append(") MapContainsStr_(k string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapContainsKeyGo).append("(BoxString(k)) }\n");
      }
      out.append("\treturn this_.mc().ContainsStr(k)\n}\n");

      out.append("func (this_ *").append(structName).append(") MapPutStr_(k string, v any) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapPutGo).append("(BoxString(k), ToAny(v)) }\n");
      }
      if (isImm) {
        out.append("\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      } else {
        out.append("\treturn this_.mc().PutStr(k, v)\n}\n");
      }

      out.append("func (this_ *").append(structName).append(") MapRemoveStr_(k string) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapRemoveGo).append("(BoxString(k)) }\n");
      }
      if (isImm) {
        out.append("\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      } else {
        out.append("\treturn this_.mc().RemoveStr(k)\n}\n");
      }

      out.append("func (this_ *").append(structName).append(") MapGetNode_(n uint32) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapGetGo).append("(boxNodeById(n)) }\n");
      }
      out.append("\treturn this_.mc().GetNode(n)\n}\n");

      out.append("func (this_ *").append(structName).append(") MapGetOrDefaultNode_(n uint32, def any) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapGetOrDefGo).append("(boxNodeById(n), ToAny(def)) }\n");
      }
      out.append("\treturn this_.mc().GetOrDefaultNode(n, def)\n}\n");

      out.append("func (this_ *").append(structName).append(") MapContainsNode_(n uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapContainsKeyGo).append("(boxNodeById(n)) }\n");
      }
      out.append("\treturn this_.mc().ContainsNode(n)\n}\n");

      out.append("func (this_ *").append(structName).append(") MapPutNode_(n uint32, v any) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapPutGo).append("(boxNodeById(n), ToAny(v)) }\n");
      }
      if (isImm) {
        out.append("\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      } else {
        out.append("\treturn this_.mc().PutNode(n, v)\n}\n");
      }

      out.append("func (this_ *").append(structName).append(") MapRemoveNode_(n uint32) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapRemoveGo).append("(boxNodeById(n)) }\n");
      }
      if (isImm) {
        out.append("\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      } else {
        out.append("\treturn this_.mc().RemoveNode(n)\n}\n");
      }

      out.append("func (this_ *").append(structName).append(") MapComputeIfAbsentStr_(k string, fn Ext_java_util_function_Function) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapComputeIfAbsentGo).append("(BoxString(k), fn) }\n");
      }
      if (isImm) {
        out.append("\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      } else {
        out.append("\treturn mapCoreComputeIfAbsentStr(this_.mc(), k, fn)\n}\n");
      }

      out.append("func (this_ *").append(structName).append(") MapComputeIfAbsentNode_(n uint32, fn Ext_java_util_function_Function) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return this_.Self_.").append(mapComputeIfAbsentGo).append("(boxNodeById(n), fn) }\n");
      }
      if (isImm) {
        out.append("\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      } else {
        out.append("\treturn mapCoreComputeIfAbsentNode(this_.mc(), n, fn)\n}\n");
      }

      out.append("func (this_ *").append(structName).append(") MapTryGetComputeStr_(k string) (any, bool) {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (isImm) {
        out.append("\t_ = k\n\treturn nil, false\n}\n");
      } else {
        if (subclassed) {
          out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
              .append("(this_) { return nil, false }\n");
        }
        out.append("\treturn mapCoreTryGetComputeStr(this_.mc(), k)\n}\n");
      }

      out.append("func (this_ *").append(structName).append(") MapTryGetComputeNode_(n uint32) (any, bool) {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (isImm) {
        out.append("\t_ = n\n\treturn nil, false\n}\n");
      } else {
        if (subclassed) {
          out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
              .append("(this_) { return nil, false }\n");
        }
        out.append("\treturn mapCoreTryGetComputeNode(this_.mc(), n)\n}\n");
      }

      out.append("func (this_ *").append(structName).append(") MapTryGetComputeAny_(k any) (any, bool) {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n");
      if (isImm) {
        out.append("\t_ = k\n\treturn nil, false\n}\n");
      } else {
        if (subclassed) {
          out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
              .append("(this_) { return nil, false }\n");
        }
        out.append("\treturn mapCoreTryGetComputeAny(this_.mc(), k)\n}\n");
      }
      return;
    }

    if (qn.equals("java.util.concurrent.ConcurrentHashMap") && !subclassed) {
      out.append("func (this_ *").append(structName).append(") MapGetStr_(k string) any { return syncMapGetStr(this_, k) }\n")
          .append("func (this_ *").append(structName).append(") MapGetOrDefaultStr_(k string, def any) any { return syncMapGetOrDefaultStr(this_, k, def) }\n")
          .append("func (this_ *").append(structName).append(") MapContainsStr_(k string) bool { return syncMapContainsStr(this_, k) }\n")
          .append("func (this_ *").append(structName).append(") MapPutStr_(k string, v any) any { return syncMapPutStr(this_, k, v) }\n")
          .append("func (this_ *").append(structName).append(") MapRemoveStr_(k string) any { return syncMapRemoveStr(this_, k) }\n")
          .append("func (this_ *").append(structName).append(") MapGetNode_(n uint32) any { return syncMapGetNode(this_, n) }\n")
          .append("func (this_ *").append(structName).append(") MapGetOrDefaultNode_(n uint32, def any) any { return syncMapGetOrDefaultNode(this_, n, def) }\n")
          .append("func (this_ *").append(structName).append(") MapContainsNode_(n uint32) bool { return syncMapContainsNode(this_, n) }\n")
          .append("func (this_ *").append(structName).append(") MapPutNode_(n uint32, v any) any { return syncMapPutNode(this_, n, v) }\n")
          .append("func (this_ *").append(structName).append(") MapRemoveNode_(n uint32) any { return syncMapRemoveNode(this_, n) }\n")
          .append("func (this_ *").append(structName).append(") MapComputeIfAbsentStr_(k string, fn Ext_java_util_function_Function) any { return syncMapComputeIfAbsentStr(this_, k, fn) }\n")
          .append("func (this_ *").append(structName).append(") MapComputeIfAbsentNode_(n uint32, fn Ext_java_util_function_Function) any { return syncMapComputeIfAbsentNode(this_, n, fn) }\n")
          .append("func (this_ *").append(structName).append(") MapTryGetComputeStr_(k string) (any, bool) { return syncMapTryGetComputeStr(this_, k) }\n")
          .append("func (this_ *").append(structName).append(") MapTryGetComputeNode_(n uint32) (any, bool) { return syncMapTryGetComputeNode(this_, n) }\n")
          .append("func (this_ *").append(structName).append(") MapTryGetComputeAny_(k any) (any, bool) { return syncMapTryGetComputeAny(this_, k) }\n");
      return;
    }

    if (qn.equals("com.google.common.collect.HashBiMap")
        || qn.equals("com.google.common.collect.ImmutableBiMap")) {
      boolean isImmBi = qn.equals("com.google.common.collect.ImmutableBiMap");
      out.append("func (this_ *").append(structName).append(") fastMapCore_() *JMapCore {\n")
          .append("\tif this_ == nil { return nil }\n");
      if (subclassed) {
        out.append("\tif this_.Self_ != nil && this_.Self_ != ").append(selfIface)
            .append("(this_) { return nil }\n");
      }
      out.append("\tif bc, ok := this_.Core_.(*jBiMapCore); ok { return bc.fwd }\n\treturn nil\n}\n");
      out.append("func (this_ *").append(structName).append(") MapGetStr_(k string) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.Core_.(*jBiMapCore).fwd.GetStr(k)\n}\n")
          .append("func (this_ *").append(structName).append(") MapGetOrDefaultStr_(k string, def any) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.Core_.(*jBiMapCore).fwd.GetOrDefaultStr(k, def)\n}\n")
          .append("func (this_ *").append(structName).append(") MapContainsStr_(k string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.Core_.(*jBiMapCore).fwd.ContainsStr(k)\n}\n")
          .append("func (this_ *").append(structName).append(") MapPutStr_(k string, v any) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn biMapPutStr(this_.Core_.(*jBiMapCore), k, v)\n}\n")
          .append("func (this_ *").append(structName).append(") MapRemoveStr_(k string) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn biMapRemoveStr(this_.Core_.(*jBiMapCore), k)\n}\n")
          .append("func (this_ *").append(structName).append(") MapGetNode_(n uint32) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.Core_.(*jBiMapCore).fwd.GetNode(n)\n}\n")
          .append("func (this_ *").append(structName).append(") MapGetOrDefaultNode_(n uint32, def any) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.Core_.(*jBiMapCore).fwd.GetOrDefaultNode(n, def)\n}\n")
          .append("func (this_ *").append(structName).append(") MapContainsNode_(n uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.Core_.(*jBiMapCore).fwd.ContainsNode(n)\n}\n")
          .append("func (this_ *").append(structName).append(") MapPutNode_(n uint32, v any) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn biMapOp_put_2(this_.Core_.(*jBiMapCore), boxNodeById(n), v)\n}\n")
          .append("func (this_ *").append(structName).append(") MapRemoveNode_(n uint32) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn biMapOp_remove_1(this_.Core_.(*jBiMapCore), boxNodeById(n))\n}\n")
          .append("func (this_ *").append(structName).append(") MapComputeIfAbsentStr_(k string, fn Ext_java_util_function_Function) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn biMapComputeIfAbsentStr(this_.Core_.(*jBiMapCore), k, fn)\n}\n")
          .append("func (this_ *").append(structName).append(") MapComputeIfAbsentNode_(n uint32, fn Ext_java_util_function_Function) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn biMapComputeIfAbsentNode(this_.Core_.(*jBiMapCore), n, fn)\n}\n");
      if (isImmBi) {
        out.append("func (this_ *").append(structName).append(") MapTryGetComputeStr_(k string) (any, bool) {\n")
            .append("\tif this_ == nil { ThrowNPE() }\n\t_ = k\n\treturn nil, false\n}\n")
            .append("func (this_ *").append(structName).append(") MapTryGetComputeNode_(n uint32) (any, bool) {\n")
            .append("\tif this_ == nil { ThrowNPE() }\n\t_ = n\n\treturn nil, false\n}\n")
            .append("func (this_ *").append(structName).append(") MapTryGetComputeAny_(k any) (any, bool) {\n")
            .append("\tif this_ == nil { ThrowNPE() }\n\t_ = k\n\treturn nil, false\n}\n");
      } else {
        out.append("func (this_ *").append(structName).append(") MapTryGetComputeStr_(k string) (any, bool) {\n")
            .append("\tif this_ == nil { ThrowNPE() }\n\tif v := this_.Core_.(*jBiMapCore).fwd.GetStr(k); v != nil { return v, true }\n\treturn nil, false\n}\n")
            .append("func (this_ *").append(structName).append(") MapTryGetComputeNode_(n uint32) (any, bool) {\n")
            .append("\tif this_ == nil { ThrowNPE() }\n\tif v := this_.Core_.(*jBiMapCore).fwd.GetNode(n); v != nil { return v, true }\n\treturn nil, false\n}\n")
            .append("func (this_ *").append(structName).append(") MapTryGetComputeAny_(k any) (any, bool) {\n")
            .append("\tif this_ == nil { ThrowNPE() }\n\tif v := this_.Core_.(*jBiMapCore).fwd.Get(k); v != nil { return v, true }\n\treturn nil, false\n}\n");
      }
      return;
    }

    if (MULTIMAP_TYPES.contains(qn)) {
      out.append("func (this_ *").append(structName).append(") MultimapContainsKeyStr_(k string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.getCore().(*jMultimapCore).containsKeyStr(k)\n}\n")
          .append("func (this_ *").append(structName).append(") MultimapContainsKeyNode_(u uint32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.getCore().(*jMultimapCore).containsKeyNode(u)\n}\n")
          .append("func (this_ *").append(structName).append(") MultimapContainsEntryStr_(k string, v any) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.getCore().(*jMultimapCore).containsEntryStr(k, v)\n}\n")
          .append("func (this_ *").append(structName).append(") MultimapContainsEntryNode_(u uint32, v any) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.getCore().(*jMultimapCore).containsEntryNode(u, v)\n}\n")
          .append("func (this_ *").append(structName).append(") MultimapGetStr_(k string) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.getCore().(*jMultimapCore).getStr(k)\n}\n")
          .append("func (this_ *").append(structName).append(") MultimapGetNode_(u uint32) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.getCore().(*jMultimapCore).getNode(u)\n}\n")
          .append("func (this_ *").append(structName).append(") MultimapPutStr_(k string, v any) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.getCore().(*jMultimapCore).putStr(k, v)\n}\n")
          .append("func (this_ *").append(structName).append(") MultimapPutNode_(u uint32, v any) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.getCore().(*jMultimapCore).putNode(u, v)\n}\n")
          .append("func (this_ *").append(structName).append(") MultimapRemoveStr_(k string, v any) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.getCore().(*jMultimapCore).removeStr(k, v)\n}\n")
          .append("func (this_ *").append(structName).append(") MultimapRemoveNode_(u uint32, v any) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.getCore().(*jMultimapCore).removeNode(u, v)\n}\n")
          .append("func (this_ *").append(structName).append(") MultimapRemoveAllStr_(k string) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.getCore().(*jMultimapCore).removeAllStr(k)\n}\n")
          .append("func (this_ *").append(structName).append(") MultimapRemoveAllNode_(u uint32) any {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.getCore().(*jMultimapCore).removeAllNode(u)\n}\n")
          .append("func (this_ *").append(structName).append(") MultimapPutAllStr_(k string, vals Ext_java_lang_Iterable) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.getCore().(*jMultimapCore).putAllStr(k, vals)\n}\n")
          .append("func (this_ *").append(structName).append(") MultimapPutAllNode_(u uint32, vals Ext_java_lang_Iterable) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn this_.getCore().(*jMultimapCore).putAllNode(u, vals)\n}\n");
      return;
    }

    if (MULTISET_TYPES.contains(qn)) {
      out.append("func (this_ *").append(structName).append(") MultisetCountStr_(s string) int32 {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn multisetCountStr(this_.getCore().(*jMultisetCore), s)\n}\n")
          .append("func (this_ *").append(structName).append(") MultisetContainsStr_(s string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn multisetCountStr(this_.getCore().(*jMultisetCore), s) > 0\n}\n")
          .append("func (this_ *").append(structName).append(") MultisetAddStr1_(s string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\tmultisetAddStr2(this_.getCore().(*jMultisetCore), s, 1)\n\treturn true\n}\n")
          .append("func (this_ *").append(structName).append(") MultisetAddStr2_(s string, n int32) int32 {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn multisetAddStr2(this_.getCore().(*jMultisetCore), s, n)\n}\n")
          .append("func (this_ *").append(structName).append(") MultisetRemoveStr1_(s string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn multisetRemoveStr1(this_.getCore().(*jMultisetCore), s)\n}\n")
          .append("func (this_ *").append(structName).append(") MultisetSetCountStr3_(s string, oldCount, newCount int32) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn multisetSetCountStr3(this_.getCore().(*jMultisetCore), s, oldCount, newCount)\n}\n")
          .append("func (this_ *").append(structName).append(") SetContainsStr_(s string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn multisetCountStr(this_.getCore().(*jMultisetCore), s) > 0\n}\n")
          .append("func (this_ *").append(structName).append(") SetAddStr_(s string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\tmultisetAddStr2(this_.getCore().(*jMultisetCore), s, 1)\n\treturn true\n}\n")
          .append("func (this_ *").append(structName).append(") SetRemoveStr_(s string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\treturn multisetRemoveStr1(this_.getCore().(*jMultisetCore), s)\n}\n");
      return;
    }

    if (TABLE_TYPES.contains(qn)) {
      boolean isImmTbl = qn.equals("com.google.common.collect.ImmutableTable");
      out.append("func (this_ *").append(structName).append(") fastTableCore_() (*jTableCore, bool) {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n")
          .append("\tc, _ := this_.getCore().(*jTableCore)\n")
          .append("\treturn c, ").append(isImmTbl).append("\n}\n");
    }

    if (qn.equals("com.google.common.collect.Sets.SetView")) {
      out.append("func (this_ *").append(structName).append(") SetContainsStr_(k string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n")
          .append("\tc := this_.Core_.(*jSetViewCore)\n")
          .append("\tswitch c.mode {\n")
          .append("\tcase 0:\n\t\treturn SetContainsStr(c.s1, k) || SetContainsStr(c.s2, k)\n")
          .append("\tcase 1:\n\t\treturn SetContainsStr(c.s1, k) && SetContainsStr(c.s2, k)\n")
          .append("\tdefault:\n\t\treturn SetContainsStr(c.s1, k) && !SetContainsStr(c.s2, k)\n\t}\n}\n")
          .append("func (this_ *").append(structName).append(") SetAddStr_(_ string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n")
          .append("func (this_ *").append(structName).append(") SetRemoveStr_(_ string) bool {\n")
          .append("\tif this_ == nil { ThrowNPE() }\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
    }
  }

  public static boolean tryEmitStaticField(
      RuntimeBindings rb,
      TypeElement owner,
      VariableElement ve,
      String varGoName,
      String fieldType,
      StringBuilder out) {
    String qn = owner.getQualifiedName().toString();
    String fn = ve.getSimpleName().toString();
    if (fn.equals("class") || fn.equals("Class")) {
      out.append("var ")
          .append(varGoName)
          .append(' ')
          .append(fieldType)
          .append(" = ClassLiteral(")
          .append(GoNames.goStringLiteral(qn))
          .append(", ")
          .append(GoNames.goStringLiteral(owner.getSimpleName().toString()))
          .append(")\n");
      return true;
    }
    return false;
  }

  public static boolean tryEmitConstructorInit(
      RuntimeBindings rb,
      TypeElement owner,
      ExecutableElement ctor,
      PackagePlanner.ImportTracker imports,
      String targetVar,
      StringBuilder out) {
    String qn = owner.getQualifiedName().toString();
    int arity = ctor.getParameters().size();
    String p0 =
        arity >= 1 ? rb.goTypeInJ2gort(ctor.getParameters().get(0).asType(), imports) : "";

    if (qn.equals("java.util.ArrayDeque")) {
      out.append("\tdc := &").append(targetVar).append(".inlineCore\n");
      if (arity == 1) {
        if (p0.equals("int32")) {
          out.append("\tdc.InitCap(arg0)\n");
        } else {
          out.append("\tfor it := NewIterCursor(arg0); it.HasNext(); { dc.AddLast(it.Next()) }\n");
        }
      } else {
        out.append("\t_ = dc\n");
      }
      return true;
    }

    if (CONCRETE_LIST_TYPES.contains(qn)) {
      if (qn.equals("java.util.ArrayList")) {
        if (arity == 1) {
          if (p0.equals("int32")) {
            out.append('\t').append(targetVar).append(".inlineCore.InitCap(arg0)\n");
          } else {
            out.append("\tlistCoreAddAllIterable(&").append(targetVar).append(".inlineCore, arg0)\n");
          }
        }
        return true;
      }
      if (qn.equals("com.google.common.collect.ImmutableList")) {
        out.append('\t').append(targetVar).append(".items_ = emptyAnySlice_\n");
        return true;
      }
      out.append("\tlc := NewJListCore()\n");
      if (arity == 1 && !p0.equals("int32")) {
        out.append("\tlistCoreAddAllIterable(lc, arg0)\n");
      }
      out.append('\t').append(targetVar).append(".Core_ = lc\n");
      return true;
    }

    if (qn.equals("java.util.EnumSet")) {
      return true;
    }

    if (CONCRETE_SET_TYPES.contains(qn)) {
      if (qn.equals("java.util.TreeSet")) {
        if (arity == 1 && p0.contains("Comparator")) {
          out.append("\tinitSortedCoreCmp(&").append(targetVar).append(".inlineCore, arg0)\n");
        } else {
          out.append('\t').append(targetVar).append(".inlineCore.mode = jMapModeSortedNat\n");
          if (arity == 1 && !p0.equals("int32")) {
            out.append("\tmapCoreAddAllIterable(&").append(targetVar).append(".inlineCore, arg0)\n");
          }
        }
        return true;
      }
      if (qn.equals("com.google.common.collect.ImmutableSortedSet")) {
        out.append('\t').append(targetVar).append(".elems_ = emptyAnySlice_\n");
        if (arity == 1 && p0.contains("Comparator")) {
          out.append('\t').append(targetVar).append(".sortedCore_ = newSortedCoreFromCmp(arg0)\n");
        } else {
          out.append('\t').append(targetVar).append(".sortedCore_ = NewSortedJMapCore(nil)\n");
        }
        return true;
      }
      if (qn.equals("java.util.HashSet") || qn.equals("java.util.LinkedHashSet")) {
        if (arity == 2 && p0.equals("int32")) {
          out.append("\tinitSetCapLoad(&").append(targetVar).append(".inlineSet, arg0, arg1)\n");
        } else if (arity == 1 && p0.equals("int32")) {
          out.append('\t').append(targetVar).append(".inlineSet.InitExpectedSize(arg0)\n");
        } else if (arity == 1 && !p0.equals("float32")) {
          out.append("\tsetCoreAddAllIterable(&").append(targetVar).append(".inlineSet, arg0)\n");
        }
        return true;
      }
      if (qn.equals("com.google.common.collect.ImmutableSet")
          || qn.equals("com.google.common.collect.ImmutableSet.CachingAsList")) {
        out.append('\t').append(targetVar).append(".elems_ = emptyAnySlice_\n");
        return true;
      }
      out.append("\tmc := NewJMapCore(false)\n");
      if (arity == 1 && !p0.equals("int32") && !p0.equals("float32")) {
        out.append("\tmapCoreAddAllIterable(mc, arg0)\n");
      }
      out.append('\t').append(targetVar).append(".Core_ = mc\n");
      return true;
    }

    if (qn.equals("java.util.EnumMap")) {
      if (arity == 1) {
        if (p0.contains("Class")) {
          out.append('\t').append(targetVar).append(".inlineEnumMap.initClass(arg0)\n");
        } else if (p0.contains("EnumMap")) {
          out.append('\t').append(targetVar).append(".inlineEnumMap.initCopyEnumMap(arg0)\n");
        } else if (p0.contains("Map")) {
          out.append('\t').append(targetVar).append(".inlineEnumMap.initCopyMap(arg0)\n");
        }
      }
      return true;
    }

    if (CONCRETE_MAP_TYPES.contains(qn)) {
      if (qn.equals("java.util.TreeMap")
          || qn.equals("com.google.common.collect.ImmutableSortedMap")) {
        if (arity == 1 && p0.contains("Comparator")) {
          out.append("\tinitSortedCoreCmp(&").append(targetVar).append(".inlineCore, arg0)\n");
        } else {
          out.append('\t').append(targetVar).append(".inlineCore.mode = jMapModeSortedNat\n");
          if (arity == 1 && p0.contains("Map")) {
            out.append("\tcopyExtMapInto(&").append(targetVar).append(".inlineCore, arg0)\n");
          }
        }
        return true;
      }
      boolean id = qn.equals("java.util.IdentityHashMap");
      if (qn.equals("java.util.HashMap")
          || qn.equals("java.util.LinkedHashMap")
          || qn.equals("com.google.common.collect.ImmutableMap")
          || id) {
        if (id) {
          out.append('\t').append(targetVar).append(".inlineCore.Identity = true\n");
          out.append('\t').append(targetVar).append(".inlineCore.genNonCustom = true\n");
          if (arity == 1 && p0.equals("int32")) {
            out.append('\t').append(targetVar).append(".inlineCore.InitExpectedSize(arg0)\n");
          } else if (arity == 1 && p0.contains("Map")) {
            out.append("\tcopyExtMapInto(&").append(targetVar).append(".inlineCore, arg0)\n");
          }
          return true;
        }
        if (arity == 3 && p0.equals("int32")) {
          out.append('\t').append(targetVar).append(".inlineCore.InitAccessOrder(arg0, arg1, arg2)\n");
        } else if (arity == 2 && p0.equals("int32")) {
          out.append('\t').append(targetVar).append(".inlineCore.InitAccessOrder(arg0, arg1, false)\n");
        } else if (arity == 1 && p0.equals("int32")) {
          out.append('\t').append(targetVar).append(".inlineCore.InitExpectedSize(arg0)\n");
        } else if (arity == 1 && p0.contains("Map")) {
          out.append("\tcopyExtMapInto(&").append(targetVar).append(".inlineCore, arg0)\n");
        }
        return true;
      }
      out.append("\tmc := NewJMapCore(").append(id).append(")\n");
      if (arity == 1 && p0.contains("Map")) {
        out.append("\tcopyExtMapInto(mc, arg0)\n");
      }
      out.append('\t').append(targetVar).append(".Core_ = mc\n");
      return true;
    }

    switch (qn) {
      case "java.util.AbstractCollection",
          "java.util.AbstractSet",
          "java.util.AbstractList",
          "java.util.AbstractSequentialList",
          "java.util.AbstractQueue",
          "java.util.AbstractMap",
          "com.google.common.collect.Maps.IteratorBasedAbstractMap",
          "com.google.common.collect.Ordering",
          "com.google.common.collect.UnmodifiableIterator",
          "com.google.common.cache.CacheLoader",
          "java.util.concurrent.AbstractExecutorService" -> {
        return true;
      }
      case "com.google.common.collect.AbstractIterator" -> {
        out.append('\t').append(targetVar).append(".Core_ = &jAbstractIterState{}\n");
        return true;
      }
      case "java.util.AbstractMap.SimpleImmutableEntry" -> {
        if (arity == 2) {
          out.append('\t')
              .append(targetVar)
              .append(".Core_ = &jSimpleEntry{k: ToAny(arg0), v: ToAny(arg1), imm: true}\n");
          return true;
        }
      }
      case "com.google.common.collect.ImmutableList.Builder" -> {
        if (arity == 1 && p0.equals("int32")) {
          out.append("\tInitImmListBuilderCap(").append(targetVar).append(", arg0)\n");
        }
        return true;
      }
      case "com.google.common.collect.ImmutableCollection.Builder" -> {
        out.append('\t').append(targetVar).append(".Core_ = NewJListCore()\n");
        return true;
      }
      case "com.google.common.collect.ImmutableSet.Builder" -> {
        if (arity == 1 && p0.equals("int32")) {
          out.append('\t').append(targetVar).append(".inlineSet.InitExpectedSize(arg0)\n");
        }
        return true;
      }
      case "com.google.common.collect.ImmutableSortedSet.Builder" -> {
        if (arity == 1 && p0.contains("Comparator")) {
          out.append('\t').append(targetVar).append(".sortedCore_ = newSortedCoreFromCmp(arg0)\n");
        } else {
          out.append('\t').append(targetVar).append(".sortedCore_ = NewSortedJMapCore(nil)\n");
        }
        return true;
      }
      case "com.google.common.collect.ImmutableMap.Builder",
          "com.google.common.collect.ImmutableBiMap.Builder" -> {
        out.append("\tb := ").append(targetVar).append(".imbCore()\n");
        if (arity == 1 && p0.equals("int32")) {
          out.append("\tb.mc.InitExpectedSize(arg0)\n");
        } else {
          out.append("\t_ = b\n");
        }
        return true;
      }
      case "com.google.common.collect.ImmutableSortedMap.Builder" -> {
        if (arity == 1 && p0.contains("Comparator")) {
          out.append("\tinitSortedCoreCmp(&").append(targetVar).append(".inlineCore, arg0)\n");
        } else {
          out.append('\t').append(targetVar).append(".inlineCore.mode = jMapModeSortedNat\n");
        }
        out.append('\t').append(targetVar).append(".inlineBuilder.mc = &").append(targetVar).append(".inlineCore\n");
        out.append('\t').append(targetVar).append(".inlineBuilder.sorted = true\n");
        out.append('\t').append(targetVar).append(".Core_ = &").append(targetVar).append(".inlineBuilder\n");
        return true;
      }
      case "com.google.common.collect.ImmutableMultimap.Builder",
          "com.google.common.collect.ImmutableListMultimap.Builder" -> {
        out.append('\t').append(targetVar).append(".Core_ = newMultimapBuilderCore(false)\n");
        return true;
      }
      case "com.google.common.collect.ImmutableSetMultimap.Builder" -> {
        out.append('\t').append(targetVar).append(".Core_ = newMultimapBuilderCore(true)\n");
        return true;
      }
      case "com.google.common.collect.ImmutableMultiset.Builder" -> {
        out.append('\t').append(targetVar).append(".Core_ = newMultisetCore(false, nil)\n");
        return true;
      }
      case "java.util.BitSet" -> {
        if (arity == 0) {
          out.append("\tInitExt_java_util_BitSet_0(").append(targetVar).append(")\n");
          return true;
        }
        if (arity == 1 && p0.equals("int32")) {
          out.append("\tInitExt_java_util_BitSet_1(").append(targetVar).append(", arg0)\n");
          return true;
        }
      }
      case "java.util.PriorityQueue" -> {
        if (arity == 0) {
          out.append('\t').append(targetVar).append(".Core_ = newPriorityQueue(CompareValues)\n");
          return true;
        }
        if (arity == 1 && p0.equals("int32")) {
          out.append('\t').append(targetVar).append(".Core_ = newPriorityQueue(CompareValues)\n");
          return true;
        }
        if (arity == 1 && p0.contains("Comparator")) {
          out.append('\t')
              .append(targetVar)
              .append(".Core_ = newPriorityQueue(comparatorFuncOrNatural(arg0))\n");
          return true;
        }
        if (arity == 2 && p0.equals("int32")) {
          out.append('\t')
              .append(targetVar)
              .append(".Core_ = newPriorityQueue(comparatorFuncOrNatural(arg1))\n");
          return true;
        }
        if (arity == 1) {
          out.append("\tpq := newPriorityQueue(CompareValues)\n");
          out.append("\tfor it := Iterate(arg0); it.HasNext(); { pq.offer(it.Next()) }\n");
          out.append('\t').append(targetVar).append(".Core_ = pq\n");
          return true;
        }
      }
      case "java.util.concurrent.ConcurrentHashMap" -> {
        out.append('\t').append(targetVar).append(".Core_ = &jSyncMapCore{mc: NewJMapCore(false)}\n");
        return true;
      }
      case "java.util.concurrent.ThreadPoolExecutor" -> {
        if (arity >= 1 && p0.equals("int32")) {
          out.append('\t').append(targetVar).append(".Core_ = newThreadPoolExecutor(arg0)\n");
        } else {
          out.append('\t').append(targetVar).append(".Core_ = &jDirectExecutor{}\n");
        }
        return true;
      }
      case "java.util.concurrent.atomic.AtomicIntegerArray" -> {
        if (arity == 1 && p0.equals("int32")) {
          out.append('\t')
              .append(targetVar)
              .append(".Core_ = &jAtomicIntArray{vals: make([]int32, arg0)}\n");
          return true;
        }
      }
      case "java.util.Random" -> {
        if (arity == 0) {
          out.append('\t').append(targetVar).append(".Core_ = javautil.NewRandom()\n");
          return true;
        }
        if (arity == 1 && p0.equals("int64")) {
          out.append('\t').append(targetVar).append(".Core_ = javautil.NewRandomSeed(arg0)\n");
          return true;
        }
      }
      default -> {}
    }
    return false;
  }

  public static boolean tryEmitStaticMethod(
      RuntimeBindings rb,
      TypeElement owner,
      ExecutableElement m,
      PackagePlanner.ImportTracker imports,
      String retType,
      StringBuilder out) {
    String qn = owner.getQualifiedName().toString();
    String simple = m.getSimpleName().toString();
    int arity = m.getParameters().size();
    String p0 =
        arity >= 1 ? rb.goTypeInJ2gort(m.getParameters().get(0).asType(), imports) : "";

    // --- ImmutableList / java.util.List ---
    if (qn.equals("com.google.common.collect.ImmutableList") || qn.equals("java.util.List")) {
      if (simple.equals("of") && retType.contains("List")) {
        if (arity == 0) {
          out.append("\treturn emptyImmutableList_\n}\n\n");
          return true;
        }
        boolean hasVarArgs =
            m.isVarArgs()
                && rb.goTypeInJ2gort(m.getParameters().get(arity - 1).asType(), imports)
                    .startsWith("[]");
        if (!hasVarArgs) {
          if (arity >= 1 && arity <= 8) {
            out.append("\treturn newImmList").append(arity).append('(');
            for (int i = 0; i < arity; i++) {
              if (i > 0) out.append(", ");
              out.append("arg").append(i);
            }
            out.append(")\n}\n\n");
            return true;
          }
          out.append("\titems := make([]any, ").append(arity).append(")\n");
          for (int i = 0; i < arity; i++) {
            out.append("\titems[")
                .append(i)
                .append("] = CheckNotNull(ToAny(arg")
                .append(i)
                .append("))\n");
          }
          out.append("\treturn wrapImmutableListOwn(items)\n}\n\n");
          return true;
        }
        out.append("\trest := CheckNotNullSlice(arg").append(arity - 1).append(")\n");
        out.append("\tn := ").append(arity - 1).append(" + len(rest)\n");
        out.append("\tif n == 0 { return emptyImmutableList_ }\n");
        out.append("\titems := make([]any, 0, n)\n");
        for (int i = 0; i < arity - 1; i++) {
          out.append("\titems = append(items, CheckNotNull(ToAny(arg").append(i).append(")))\n");
        }
        out.append("\tfor _, v := range rest { items = append(items, CheckNotNull(ToAny(v))) }\n");
        out.append("\treturn wrapImmutableListOwn(items)\n}\n\n");
        return true;
      }
      if (simple.equals("copyOf") && arity == 1) {
        if (p0.equals("[]any")) {
          out.append("\treturn immListCopyOfSlice(arg0)\n}\n\n");
        } else if (p0.equals("[]string")) {
          out.append("\tCheckNotNullSlice(arg0)\n\treturn wrapImmStrSliceCopy(arg0)\n}\n\n");
        } else if (p0.startsWith("[]")) {
          out.append(
              "\tCheckNotNullSlice(arg0)\n\tswitch len(arg0) {\n\tcase 0: return emptyImmutableList_\n\tcase 1: return newImmList1(arg0[0])\n\tcase 2: return newImmList2(arg0[0], arg0[1])\n\tcase 3: return newImmList3(arg0[0], arg0[1], arg0[2])\n\tcase 4: return newImmList4(arg0[0], arg0[1], arg0[2], arg0[3])\n\tcase 5: return newImmList5(arg0[0], arg0[1], arg0[2], arg0[3], arg0[4])\n\tcase 6: return newImmList6(arg0[0], arg0[1], arg0[2], arg0[3], arg0[4], arg0[5])\n\tcase 7: return newImmList7(arg0[0], arg0[1], arg0[2], arg0[3], arg0[4], arg0[5], arg0[6])\n\tcase 8: return newImmList8(arg0[0], arg0[1], arg0[2], arg0[3], arg0[4], arg0[5], arg0[6], arg0[7])\n\t}\n");
          out.append("\titems := make([]any, len(arg0))\n\tfor i, v := range arg0 { items[i] = CheckNotNull(ToAny(v)) }\n\treturn wrapImmutableListOwn(items)\n}\n\n");
        } else {
          out.append("\treturn immListCopyOfIterable(arg0)\n}\n\n");
        }
        return true;
      }
      if ((simple.equals("builder") && arity == 0)
          || (simple.equals("builderWithExpectedSize") && arity == 1)) {
        out.append("\tb := &Ext_com_google_common_collect_ImmutableList_Builder{}\n");
        if (arity == 1) {
          out.append("\tInitImmListBuilderCap(b, arg0)\n");
        }
        out.append("\treturn b\n}\n\n");
        return true;
      }
      if (simple.equals("sortedCopyOf") && arity == 2) {
        out.append("\treturn immListSortedCopyOf(arg0, arg1)\n}\n\n");
        return true;
      }
      if (simple.equals("toImmutableList") && arity == 0) {
        out.append("\treturn collectorToImmutableList()\n}\n\n");
        return true;
      }
    }

    // --- ImmutableSet / java.util.Set ---
    if (qn.equals("com.google.common.collect.ImmutableSet") || qn.equals("java.util.Set")) {
      boolean isGuava = qn.startsWith("com.google.");
      if (simple.equals("of") && retType.contains("Set")) {
        if (arity == 0) {
          out.append("\treturn emptyImmutableSet_\n}\n\n");
          return true;
        }
        boolean hasVarArgs =
            m.isVarArgs()
                && rb.goTypeInJ2gort(m.getParameters().get(arity - 1).asType(), imports)
                    .startsWith("[]");
        if (!hasVarArgs && arity == 1) {
          out.append("\treturn newImmSet1(arg0)\n}\n\n");
          return true;
        }
        if (!hasVarArgs && arity == 2) {
          out.append(isGuava ? "\treturn newImmSet2(arg0, arg1)\n}\n\n" : "\treturn newJavaSet2(arg0, arg1)\n}\n\n");
          return true;
        }
        if (!hasVarArgs) {
          out.append("\titems := []any{");
          for (int i = 0; i < arity; i++) {
            if (i > 0) out.append(", ");
            out.append("ToAny(arg").append(i).append(')');
          }
          out.append("}\n\treturn immSetFromCheckedSlice(")
              .append(!isGuava)
              .append(", items)\n}\n\n");
          return true;
        }
        out.append("\trest := CheckNotNullSlice(arg").append(arity - 1).append(")\n");
        out.append("\titems := make([]any, 0, ").append(arity - 1).append("+len(rest))\n");
        for (int i = 0; i < arity - 1; i++) {
          out.append("\titems = append(items, ToAny(arg").append(i).append("))\n");
        }
        out.append("\tfor _, v := range rest { items = append(items, ToAny(v)) }\n");
        out.append("\treturn immSetFromCheckedSlice(").append(!isGuava).append(", items)\n}\n\n");
        return true;
      }
      if (simple.equals("copyOf") && arity == 1) {
        if (p0.equals("[]any")) {
          out.append("\treturn immSetCopyOfSlice(arg0)\n}\n\n");
        } else if (p0.equals("[]string")) {
          out.append("\tCheckNotNullSlice(arg0)\n\treturn immSetFromStrSlice(false, arg0)\n}\n\n");
        } else if (p0.startsWith("[]")) {
          out.append("\tCheckNotNullSlice(arg0)\n\tif len(arg0) == 0 { return emptyImmutableSet_ }\n");
          out.append("\tif len(arg0) == 1 { return newImmSet1(arg0[0]) }\n");
          out.append("\titems := make([]any, len(arg0))\n\tfor i, v := range arg0 { items[i] = CheckNotNull(ToAny(v)) }\n\treturn wrapImmutableSet(items)\n}\n\n");
        } else {
          out.append("\treturn immSetCopyOfIterable(arg0)\n}\n\n");
        }
        return true;
      }
      if ((simple.equals("builder") && arity == 0)
          || (simple.equals("builderWithExpectedSize") && arity == 1)) {
        out.append("\tb := &Ext_com_google_common_collect_ImmutableSet_BuilderBase{}\n\tb.BindExtSelf_(b)\n");
        if (arity == 1) {
          out.append("\tb.inlineSet.InitExpectedSize(arg0)\n");
        }
        out.append("\treturn b\n}\n\n");
        return true;
      }
      if (simple.equals("toImmutableSet") && arity == 0) {
        out.append("\treturn collectorToImmutableSet()\n}\n\n");
        return true;
      }
    }

    // --- ImmutableSortedSet ---
    if (qn.equals("com.google.common.collect.ImmutableSortedSet")) {
      if (simple.equals("of")) {
        if (arity == 0) {
          out.append("\treturn emptyImmutableSortedSet_\n}\n\n");
          return true;
        }
        out.append("\tmc := NewSortedJMapCore(nil)\n");
        for (int i = 0; i < arity; i++) {
          String pt = rb.goTypeInJ2gort(m.getParameters().get(i).asType(), imports);
          if (pt.startsWith("[]") && i == arity - 1 && m.isVarArgs()) {
            out.append("\tfor _, v := range CheckNotNullSlice(arg")
                .append(i)
                .append(") { mc.Put(CheckNotNull(ToAny(v)), true) }\n");
          } else {
            out.append("\tmc.Put(CheckNotNull(ToAny(arg").append(i).append(")), true)\n");
          }
        }
        out.append("\treturn wrapImmutableSortedSet(mc)\n}\n\n");
        return true;
      }
      if (simple.equals("naturalOrder") && arity == 0) {
        out.append("\tb := &Ext_com_google_common_collect_ImmutableSortedSet_Builder{}\n\tb.BindExtSelf_(b)\n\tb.sortedCore_ = NewSortedJMapCore(nil)\n\treturn b\n}\n\n");
        return true;
      }
      if ((simple.equals("copyOf") || simple.equals("copyOfSorted")) && arity == 1) {
        out.append("\treturn immSortedSetCopyOf(nil, arg0)\n}\n\n");
        return true;
      }
      if (simple.equals("copyOf") && arity == 2) {
        out.append("\treturn immSortedSetCopyOf(arg0, arg1)\n}\n\n");
        return true;
      }
      if (simple.equals("toImmutableSortedSet") && arity == 1) {
        out.append("\treturn collectorToImmutableSortedSet(arg0)\n}\n\n");
        return true;
      }
    }

    // --- ImmutableMap / java.util.Map ---
    if (qn.equals("com.google.common.collect.ImmutableMap") || qn.equals("java.util.Map")) {
      boolean isGuava = qn.startsWith("com.google.");
      if (simple.equals("of") && retType.contains("Map")) {
        if (arity == 0) {
          out.append("\treturn emptyImmutableMap_\n}\n\n");
          return true;
        }
        out.append("\tobj := newImmutableMapShell()\n\tmc := &obj.inlineCore\n");
        int pairs = arity / 2;
        out.append("\tmc.InitExpectedSize(").append(pairs).append(")\n");
        for (int i = 0; i + 1 < arity; i += 2) {
          out.append("\tk").append(i).append(" := CheckNotNull(ToAny(arg").append(i).append("))\n");
          out.append("\tv")
              .append(i)
              .append(" := CheckNotNull(ToAny(arg")
              .append(i + 1)
              .append("))\n");
          if (!isGuava) {
            out.append("\tif mc.ContainsKey(k")
                .append(i)
                .append(") { panic(javalang.NewIllegalArgumentException(\"duplicate key\")) }\n");
          } else {
            out.append("\tif mc.ContainsKey(k")
                .append(i)
                .append(") { panic(javalang.NewIllegalArgumentException(\"Multiple entries with same key: \" + StrVal(k")
                .append(i)
                .append("))) }\n");
          }
          out.append("\tmc.Put(k").append(i).append(", v").append(i).append(")\n");
        }
        out.append("\tmc.sealImmutable()\n\treturn obj\n}\n\n");
        return true;
      }
      if (simple.equals("copyOf") && arity == 1) {
        if (p0.contains("Iterable")) {
          out.append("\treturn immMapCopyOfEntries(arg0)\n}\n\n");
        } else {
          out.append("\treturn immMapCopyOf(arg0)\n}\n\n");
        }
        return true;
      }
      if ((simple.equals("builder") && arity == 0)
          || (simple.equals("builderWithExpectedSize") && arity == 1)) {
        out.append("\tb := &Ext_com_google_common_collect_ImmutableMap_BuilderBase{}\n\tb.BindExtSelf_(b)\n\timb := b.imbCore()\n");
        if (arity == 1) {
          out.append("\timb.mc.InitExpectedSize(arg0)\n");
        } else {
          out.append("\t_ = imb\n");
        }
        out.append("\treturn b\n}\n\n");
        return true;
      }
      if (simple.equals("toImmutableMap") && arity == 2) {
        out.append("\treturn collectorToImmutableMap(arg0, arg1)\n}\n\n");
        return true;
      }
      if (simple.equals("entry") && arity == 2) {
        out.append("\treturn newSimpleImmEntry(CheckNotNull(ToAny(arg0)), CheckNotNull(ToAny(arg1)))\n}\n\n");
        return true;
      }
    }

    // --- ImmutableSortedMap ---
    if (qn.equals("com.google.common.collect.ImmutableSortedMap")) {
      if (simple.equals("of") && arity == 0) {
        out.append("\treturn emptyImmutableSortedMap_\n}\n\n");
        return true;
      }
      if (simple.equals("copyOf") && arity == 1) {
        out.append("\treturn immSortedMapCopyOf(arg0)\n}\n\n");
        return true;
      }
      if (simple.equals("naturalOrder") && arity == 0) {
        out.append("\tb := &Ext_com_google_common_collect_ImmutableSortedMap_Builder{}\n\tb.BindExtSelf_(b)\n\tb.inlineCore.mode = jMapModeSortedNat\n\tb.inlineBuilder.mc = &b.inlineCore\n\tb.inlineBuilder.sorted = true\n\tb.Core_ = &b.inlineBuilder\n\treturn b\n}\n\n");
        return true;
      }
      if (simple.equals("toImmutableSortedMap") && arity == 3) {
        out.append("\treturn collectorToImmutableSortedMap(arg0, arg1, arg2)\n}\n\n");
        return true;
      }
    }

    // --- BiMap ---
    if (qn.equals("com.google.common.collect.HashBiMap") && simple.equals("create") && arity == 0) {
      out.append("\treturn newHashBiMap()\n}\n\n");
      return true;
    }
    if (qn.equals("com.google.common.collect.ImmutableBiMap")) {
      if (simple.equals("of") && arity == 0) {
        out.append("\treturn newImmutableBiMapFrom(nil)\n}\n\n");
        return true;
      }
      if (simple.equals("copyOf") && arity == 1) {
        out.append("\treturn newImmutableBiMapFrom(arg0)\n}\n\n");
        return true;
      }
      if ((simple.equals("builder") && arity == 0)
          || (simple.equals("builderWithExpectedSize") && arity == 1)) {
        out.append("\tb := &Ext_com_google_common_collect_ImmutableBiMap_Builder{}\n\tb.BindExtSelf_(b)\n\timb := b.imbCore()\n");
        if (arity == 1) {
          out.append("\timb.mc.InitExpectedSize(arg0)\n");
        } else {
          out.append("\t_ = imb\n");
        }
        out.append("\treturn b\n}\n\n");
        return true;
      }
    }

    // --- Multimap factories & builders ---
    switch (qn) {
      case "com.google.common.collect.ArrayListMultimap",
          "com.google.common.collect.LinkedListMultimap" -> {
        if (simple.equals("create") && arity == 0) {
          String st = retType.startsWith("*") ? retType.substring(1) : retType;
          out.append("\tm := &").append(st).append("{}\n\tm.Core_ = newMultimapCore(false, false, nil, nil)\n\treturn m\n}\n\n");
          return true;
        }
      }
      case "com.google.common.collect.HashMultimap",
          "com.google.common.collect.LinkedHashMultimap" -> {
        if (simple.equals("create") && arity == 0) {
          String st = retType.startsWith("*") ? retType.substring(1) : retType;
          out.append("\tm := &").append(st).append("{}\n\tm.Core_ = newMultimapCore(true, false, nil, nil)\n\treturn m\n}\n\n");
          return true;
        }
      }
      case "com.google.common.collect.TreeMultimap" -> {
        if (simple.equals("create") && arity == 0) {
          out.append("\tm := &Ext_com_google_common_collect_TreeMultimap{}\n\tm.Core_ = newTreeMultimapCore()\n\treturn m\n}\n\n");
          return true;
        }
      }
      case "com.google.common.collect.ImmutableMultimap" -> {
        if (simple.equals("builder") && arity == 0) {
          out.append("\tb := &Ext_com_google_common_collect_ImmutableMultimap_BuilderBase{}\n\tb.BindExtSelf_(b)\n\tb.Core_ = newMultimapBuilderCore(false)\n\treturn b\n}\n\n");
          return true;
        }
      }
      case "com.google.common.collect.ImmutableListMultimap" -> {
        if (simple.equals("builder") && arity == 0) {
          out.append("\tb := &Ext_com_google_common_collect_ImmutableListMultimap_Builder{}\n\tb.Core_ = newMultimapBuilderCore(false)\n\treturn b\n}\n\n");
          return true;
        }
      }
      case "com.google.common.collect.ImmutableSetMultimap" -> {
        if (simple.equals("of") && arity == 0) {
          out.append("\treturn emptyImmSetMultimap_\n}\n\n");
          return true;
        }
        if (simple.equals("builder") && arity == 0) {
          out.append("\tb := &Ext_com_google_common_collect_ImmutableSetMultimap_Builder{}\n\tb.Core_ = newMultimapBuilderCore(true)\n\treturn b\n}\n\n");
          return true;
        }
        if (simple.equals("copyOf") && arity == 1) {
          out.append("\treturn immSetMultimapCopyOf(arg0)\n}\n\n");
          return true;
        }
      }
      case "com.google.common.collect.MultimapBuilder" -> {
        if ((simple.equals("hashKeys") || simple.equals("linkedHashKeys")) && arity == 0) {
          out.append("\tb := &Ext_com_google_common_collect_MultimapBuilder_MultimapBuilderWithKeysBase{}\n\tb.BindExtSelf_(b)\n\treturn b\n}\n\n");
          return true;
        }
      }
      case "com.google.common.collect.Multimaps" -> {
        if (simple.equals("asMap") && arity == 1) {
          out.append("\tCheckNotNull(arg0)\n\treturn AsIface[Ext_java_util_Map](arg0.AsMap())\n}\n\n");
          return true;
        }
        if (simple.equals("index") && arity == 2) {
          out.append("\treturn multimapsIndex(arg0, arg1)\n}\n\n");
          return true;
        }
      }
      default -> {}
    }

    // --- Multiset factories ---
    switch (qn) {
      case "com.google.common.collect.HashMultiset",
          "com.google.common.collect.LinkedHashMultiset" -> {
        if (simple.equals("create") && arity == 0) {
          String st = retType.startsWith("*") ? retType.substring(1) : retType;
          out.append("\tm := &").append(st).append("{}\n\tm.Core_ = newMultisetCore(false, nil)\n\treturn m\n}\n\n");
          return true;
        }
      }
      case "com.google.common.collect.TreeMultiset" -> {
        if (simple.equals("create") && arity == 0) {
          out.append("\tm := &Ext_com_google_common_collect_TreeMultiset{}\n\tm.Core_ = newTreeMultisetCore()\n\treturn m\n}\n\n");
          return true;
        }
      }
      case "com.google.common.collect.ImmutableMultiset" -> {
        if (simple.equals("of") && arity == 0) {
          out.append("\treturn wrapImmutableMultiset(newMultisetCore(true, nil))\n}\n\n");
          return true;
        }
        if (simple.equals("builder") && arity == 0) {
          out.append("\tb := &Ext_com_google_common_collect_ImmutableMultiset_Builder{}\n\tb.Core_ = newMultisetCore(false, nil)\n\treturn b\n}\n\n");
          return true;
        }
        if (simple.equals("toImmutableMultiset") && arity == 0) {
          out.append("\treturn collectorToImmutableMultiset()\n}\n\n");
          return true;
        }
      }
      default -> {}
    }

    // --- Table, Range, RangeSet ---
    if (qn.equals("com.google.common.collect.HashBasedTable") && simple.equals("create") && arity == 0) {
      out.append("\tt := &Ext_com_google_common_collect_HashBasedTable{}\n\tt.Core_ = newTableCore(nil)\n\treturn t\n}\n\n");
      return true;
    }
    if (qn.equals("com.google.common.collect.Tables") && simple.equals("newCustomTable") && arity == 2) {
      out.append("\treturn newCustomTable(arg0, arg1)\n}\n\n");
      return true;
    }
    if (qn.equals("com.google.common.collect.Range") && simple.equals("closed") && arity == 2) {
      out.append("\tCheckNotNull(arg0); CheckNotNull(arg1)\n\tr := &Ext_com_google_common_collect_Range{}\n\tr.Core_ = &jRangeCore{lower: ToAny(arg0), upper: ToAny(arg1)}\n\treturn r\n}\n\n");
      return true;
    }
    if (qn.equals("com.google.common.collect.TreeRangeSet") && simple.equals("create") && arity == 0) {
      out.append("\trs := &Ext_com_google_common_collect_TreeRangeSet{}\n\trs.Core_ = &jRangeSetCore{}\n\treturn rs\n}\n\n");
      return true;
    }

    // --- EnumSet ---
    if (qn.equals("java.util.EnumSet")) {
      if (simple.equals("noneOf") && arity == 1) {
        out.append("\treturn enumSetNoneOf(arg0)\n}\n\n");
        return true;
      }
      if (simple.equals("allOf") && arity == 1) {
        out.append("\treturn enumSetAllOf(arg0)\n}\n\n");
        return true;
      }
      if (simple.equals("of")) {
        out.append("\tobj := newEnumSetShell()\n\tes := &obj.inlineEnumSet\n");
        for (int i = 0; i < arity; i++) {
          String pt = rb.goTypeInJ2gort(m.getParameters().get(i).asType(), imports);
          if (pt.startsWith("[]") && i == arity - 1 && m.isVarArgs()) {
            out.append("\tfor _, v := range CheckNotNullSlice(arg")
                .append(i)
                .append(") { es.Add(CheckNotNull(ToAny(v))) }\n");
          } else {
            out.append("\tes.Add(CheckNotNull(ToAny(arg").append(i).append(")))\n");
          }
        }
        out.append("\treturn obj\n}\n\n");
        return true;
      }
      if (simple.equals("copyOf") && arity == 1) {
        out.append("\treturn enumSetCopyOf(arg0)\n}\n\n");
        return true;
      }
      if (simple.equals("complementOf") && arity == 1) {
        out.append("\treturn enumSetComplementOf(arg0)\n}\n\n");
        return true;
      }
      if (simple.equals("range") && arity == 2) {
        out.append("\treturn enumSetRange(ToAny(arg0), ToAny(arg1))\n}\n\n");
        return true;
      }
    }

    // --- Guava Sets / Maps / Lists / Iterables / Iterators / Ordering / ComparisonChain ---
    switch (qn) {
      case "com.google.common.collect.Sets" -> {
        if (simple.equals("newIdentityHashSet") && arity == 0) {
          out.append("\tobj := &Ext_java_util_HashSetBase{}\n\tobj.BindExtSelf_(obj)\n\tobj.inlineSet.identity = true\n\treturn obj\n}\n\n");
          return true;
        }
        if (simple.equals("newHashSet") && arity == 0) {
          out.append("\tobj := &Ext_java_util_HashSetBase{}\n\tobj.BindExtSelf_(obj)\n\treturn obj\n}\n\n");
          return true;
        }
        if (simple.equals("newHashSet") && arity == 1) {
          out.append("\tobj := &Ext_java_util_HashSetBase{}\n\tobj.BindExtSelf_(obj)\n");
          if (p0.startsWith("[]")) {
            out.append("\tCheckNotNullSlice(arg0)\n\tobj.inlineSet.InitExpectedSize(int32(len(arg0)))\n\tfor _, x := range arg0 { obj.inlineSet.Add(ToAny(x)) }\n");
          } else if (p0.contains("Iterator")) {
            out.append("\tCheckNotNull(arg0)\n\tfor arg0.HasNext() { obj.inlineSet.Add(arg0.Next()) }\n");
          } else {
            out.append("\tsetCoreAddAllIterable(&obj.inlineSet, arg0)\n");
          }
          out.append("\treturn obj\n}\n\n");
          return true;
        }
        if (simple.equals("newHashSetWithExpectedSize") && arity == 1) {
          out.append("\tobj := &Ext_java_util_HashSetBase{}\n\tobj.BindExtSelf_(obj)\n\tobj.inlineSet.InitExpectedSize(arg0)\n\treturn obj\n}\n\n");
          return true;
        }
        if (simple.equals("newLinkedHashSet") && arity == 0) {
          TypeElement lhsTe = rb.resolver().elements.getTypeElement("java.util.LinkedHashSet");
          String st = lhsTe != null ? rb.goStructType(lhsTe, imports) : "Ext_java_util_LinkedHashSet";
          out.append("\tobj := &").append(st).append("{}\n\tobj.BindExtSelf_(obj)\n\treturn obj\n}\n\n");
          return true;
        }
        if (simple.equals("newLinkedHashSet") && arity == 1) {
          TypeElement lhsTe = rb.resolver().elements.getTypeElement("java.util.LinkedHashSet");
          String st = lhsTe != null ? rb.goStructType(lhsTe, imports) : "Ext_java_util_LinkedHashSet";
          out.append("\tobj := &").append(st).append("{}\n\tobj.BindExtSelf_(obj)\n\tsetCoreAddAllIterable(&obj.inlineSet, arg0)\n\treturn obj\n}\n\n");
          return true;
        }
        if (simple.equals("newLinkedHashSetWithExpectedSize") && arity == 1) {
          TypeElement lhsTe = rb.resolver().elements.getTypeElement("java.util.LinkedHashSet");
          String st = lhsTe != null ? rb.goStructType(lhsTe, imports) : "Ext_java_util_LinkedHashSet";
          out.append("\tobj := &").append(st).append("{}\n\tobj.BindExtSelf_(obj)\n\tobj.inlineSet.InitExpectedSize(arg0)\n\treturn obj\n}\n\n");
          return true;
        }
        if (simple.equals("newSetFromMap") && arity == 1) {
          out.append("\treturn newSetFromMap(arg0)\n}\n\n");
          return true;
        }
        if (simple.equals("union") && arity == 2) {
          out.append("\treturn newSetView(arg0, arg1, 0)\n}\n\n");
          return true;
        }
        if (simple.equals("intersection") && arity == 2) {
          out.append("\treturn newSetView(arg0, arg1, 1)\n}\n\n");
          return true;
        }
        if (simple.equals("difference") && arity == 2) {
          out.append("\treturn newSetView(arg0, arg1, 2)\n}\n\n");
          return true;
        }
        if (simple.equals("immutableEnumSet") && arity == 1) {
          out.append("\treturn setsImmutableEnumSetIterable(arg0)\n}\n\n");
          return true;
        }
        if (simple.equals("immutableEnumSet") && arity == 2) {
          out.append("\tvar es jEnumSetCore\n\tes.Add(CheckNotNull(ToAny(arg0)))\n\tfor _, e := range CheckNotNullSlice(arg1) { es.Add(CheckNotNull(ToAny(e))) }\n\treturn immSetFromEnumSet(&es)\n}\n\n");
          return true;
        }
      }
      case "com.google.common.collect.Maps" -> {
        if (simple.equals("newHashMap") && arity == 0) {
          out.append("\tobj := &Ext_java_util_HashMapBase{}\n\tobj.BindExtSelf_(obj)\n\treturn obj\n}\n\n");
          return true;
        }
        if (simple.equals("newHashMap") && arity == 1) {
          out.append("\tobj := &Ext_java_util_HashMapBase{}\n\tobj.BindExtSelf_(obj)\n\tcopyExtMapInto(&obj.inlineCore, arg0)\n\treturn obj\n}\n\n");
          return true;
        }
        if (simple.equals("newLinkedHashMap") && arity == 0) {
          TypeElement lhmTe = rb.resolver().elements.getTypeElement("java.util.LinkedHashMap");
          String st = lhmTe != null ? rb.goStructType(lhmTe, imports) : "Ext_java_util_LinkedHashMap";
          out.append("\tobj := &").append(st).append("{}\n\tobj.BindExtSelf_(obj)\n\treturn obj\n}\n\n");
          return true;
        }
        if (simple.equals("newLinkedHashMap") && arity == 1) {
          TypeElement lhmTe = rb.resolver().elements.getTypeElement("java.util.LinkedHashMap");
          String st = lhmTe != null ? rb.goStructType(lhmTe, imports) : "Ext_java_util_LinkedHashMap";
          out.append("\tobj := &").append(st).append("{}\n\tobj.BindExtSelf_(obj)\n\tcopyExtMapInto(&obj.inlineCore, arg0)\n\treturn obj\n}\n\n");
          return true;
        }
        if (simple.equals("newTreeMap") && arity == 0) {
          TypeElement tmTe = rb.resolver().elements.getTypeElement("java.util.TreeMap");
          String st = tmTe != null ? rb.goStructType(tmTe, imports) : "Ext_java_util_TreeMap";
          out.append("\tobj := &").append(st).append("{}\n\tobj.BindExtSelf_(obj)\n\tobj.inlineCore.mode = jMapModeSortedNat\n\treturn obj\n}\n\n");
          return true;
        }
        if (simple.equals("newTreeMap") && arity == 1 && p0.contains("Comparator")) {
          TypeElement tmTe = rb.resolver().elements.getTypeElement("java.util.TreeMap");
          String st = tmTe != null ? rb.goStructType(tmTe, imports) : "Ext_java_util_TreeMap";
          out.append("\tobj := &").append(st).append("{}\n\tobj.BindExtSelf_(obj)\n\tinitSortedCoreCmp(&obj.inlineCore, arg0)\n\treturn obj\n}\n\n");
          return true;
        }
        if (simple.equals("newLinkedHashMapWithExpectedSize") && arity == 1) {
          TypeElement lhmTe = rb.resolver().elements.getTypeElement("java.util.LinkedHashMap");
          String st = lhmTe != null ? rb.goStructType(lhmTe, imports) : "Ext_java_util_LinkedHashMap";
          out.append("\tobj := &").append(st).append("{}\n\tobj.BindExtSelf_(obj)\n\tobj.inlineCore.InitExpectedSize(arg0)\n\treturn obj\n}\n\n");
          return true;
        }
        if (simple.equals("newHashMapWithExpectedSize") && arity == 1) {
          out.append("\tobj := &Ext_java_util_HashMapBase{}\n\tobj.BindExtSelf_(obj)\n\tobj.inlineCore.InitExpectedSize(arg0)\n\treturn obj\n}\n\n");
          return true;
        }
        if (simple.equals("newIdentityHashMap") && arity == 0) {
          out.append("\tobj := &Ext_java_util_IdentityHashMap{}\n\tobj.BindExtSelf_(obj)\n\tobj.inlineCore.Identity = true\n\tobj.inlineCore.genNonCustom = true\n\treturn obj\n}\n\n");
          return true;
        }
        if (simple.equals("immutableEnumMap") && arity == 1) {
          out.append("\treturn mapsImmutableEnumMap(arg0)\n}\n\n");
          return true;
        }
      }
      case "com.google.common.collect.Lists" -> {
        if (simple.equals("charactersOf") && arity == 1) {
          out.append("\tCheckNotNullString(arg0)\n\tu := javalang.StringToUTF16(arg0)\n\tif len(u) == 0 { return emptyImmutableList_ }\n\tif len(u) == 1 { return newImmList1(ValueOfChar(u[0])) }\n\titems := make([]any, len(u))\n\tfor i, c := range u { items[i] = ValueOfChar(c) }\n\treturn wrapImmutableListOwn(items)\n}\n\n");
          return true;
        }
        if (simple.equals("newArrayList") && arity == 0) {
          out.append("\treturn wrapArrayList(nil)\n}\n\n");
          return true;
        }
        if (simple.equals("newArrayList") && arity == 1) {
          TypeElement alTe = rb.resolver().elements.getTypeElement("java.util.ArrayList");
          String alSt = alTe != null ? rb.goStructType(alTe, imports) : "Ext_java_util_ArrayListBase";
          if (p0.startsWith("[]")) {
            out.append("\tCheckNotNullSlice(arg0)\n\tobj := &")
                .append(alSt)
                .append("{}\n\tobj.BindExtSelf_(obj)\n\tobj.inlineCore.InitCap(int32(len(arg0)))\n\tfor _, x := range arg0 { obj.inlineCore.Add(ToAny(x)) }\n\treturn obj\n}\n\n");
          } else if (p0.contains("Iterator")) {
            out.append("\tCheckNotNull(arg0)\n\tobj := &")
                .append(alSt)
                .append("{}\n\tobj.BindExtSelf_(obj)\n\tfor arg0.HasNext() { obj.inlineCore.Add(arg0.Next()) }\n\treturn obj\n}\n\n");
          } else {
            out.append("\tobj := &")
                .append(alSt)
                .append("{}\n\tobj.BindExtSelf_(obj)\n\tlistCoreAddAllIterable(&obj.inlineCore, arg0)\n\treturn obj\n}\n\n");
          }
          return true;
        }
        if (simple.equals("reverse") && arity == 1) {
          out.append("\treturn newReversedListView(arg0)\n}\n\n");
          return true;
        }
        if (simple.equals("transform") && arity == 2) {
          out.append("\treturn newTransformedListView(arg0, arg1)\n}\n\n");
          return true;
        }
      }
      case "com.google.common.collect.Iterables" -> {
        switch (simple) {
          case "addAll" -> {
            out.append("\treturn iterablesAddAll(arg0, arg1)\n}\n\n");
            return true;
          }
          case "concat" -> {
            if (arity == 1) {
              out.append("\treturn iterablesConcatMany(arg0)\n}\n\n");
              return true;
            }
            if (arity == 2) {
              out.append("\treturn iterablesConcatList([]Ext_java_lang_Iterable{arg0, arg1})\n}\n\n");
              return true;
            }
            if (arity == 3) {
              out.append("\treturn iterablesConcatList([]Ext_java_lang_Iterable{arg0, arg1, arg2})\n}\n\n");
              return true;
            }
          }
          case "contains" -> {
            out.append("\treturn iterablesContains(arg0, arg1)\n}\n\n");
            return true;
          }
          case "filter" -> {
            out.append("\treturn iterablesFilter(arg0, arg1)\n}\n\n");
            return true;
          }
          case "get" -> {
            out.append("\treturn iterablesGet(arg0, arg1)\n}\n\n");
            return true;
          }
          case "getFirst" -> {
            out.append("\treturn iterablesGetFirst(arg0, arg1)\n}\n\n");
            return true;
          }
          case "getLast" -> {
            out.append("\treturn iterablesGetLast(arg0)\n}\n\n");
            return true;
          }
          case "getOnlyElement" -> {
            out.append("\treturn iterablesGetOnlyElement(arg0)\n}\n\n");
            return true;
          }
          case "isEmpty" -> {
            out.append("\tCheckNotNull(arg0)\n\tif c, ok := any(arg0).(Ext_java_util_Collection); ok { return c.IsEmpty() }\n\treturn !Iterate(arg0).HasNext()\n}\n\n");
            return true;
          }
          case "indexOf" -> {
            if (arity == 2) {
              out.append("\tCheckNotNull(arg0); CheckNotNull(arg1)\n\tidx := int32(0)\n\tfor it := Iterate(arg0); it.HasNext(); idx++ {\n\t\tif arg1.Apply(it.Next()) { return idx }\n\t}\n\treturn -1\n}\n\n");
              return true;
            }
          }
          case "size" -> {
            if (arity == 1) {
              out.append("\tCheckNotNull(arg0)\n\tif c, ok := any(arg0).(Ext_java_util_Collection); ok { return c.Size() }\n\tn := int32(0)\n\tfor it := Iterate(arg0); it.HasNext(); { _ = it.Next(); n++ }\n\treturn n\n}\n\n");
              return true;
            }
          }
          case "skip" -> {
            out.append("\treturn iterablesSkip(arg0, arg1)\n}\n\n");
            return true;
          }
          case "transform" -> {
            out.append("\treturn iterablesTransform(arg0, arg1)\n}\n\n");
            return true;
          }
          case "unmodifiableIterable" -> {
            out.append("\tCheckNotNull(arg0)\n\treturn arg0\n}\n\n");
            return true;
          }
          default -> {}
        }
      }
      case "com.google.common.collect.Iterators" -> {
        if (simple.equals("filter") && arity == 2) {
          out.append("\treturn iteratorsFilter(arg0, arg1)\n}\n\n");
          return true;
        }
        if (simple.equals("peekingIterator") && arity == 1) {
          out.append("\treturn newPeekingIter(arg0)\n}\n\n");
          return true;
        }
      }
      case "com.google.common.collect.Ordering" -> {
        if (simple.equals("natural") && arity == 0) {
          out.append("\to := &Ext_com_google_common_collect_OrderingBase{}\n\to.BindExtSelf_(o)\n\to.Core_ = jNatOrderCore{}\n\treturn o\n}\n\n");
          return true;
        }
      }
      case "com.google.common.collect.ComparisonChain" -> {
        if (simple.equals("start") && arity == 0) {
          out.append("\treturn comparisonChainActive\n}\n\n");
          return true;
        }
      }
      case "com.google.common.collect.Interners" -> {
        if ((simple.equals("newWeakInterner") || simple.equals("newStrongInterner")) && arity == 0) {
          out.append("\treturn newSimpleInterner()\n}\n\n");
          return true;
        }
      }
      case "com.google.common.primitives.Chars" -> {
        if (simple.equals("asList") && arity == 1) {
          out.append("\tCheckNotNullSlice(arg0)\n\titems := make([]any, len(arg0))\n\tfor i, c := range arg0 { items[i] = ValueOfChar(c) }\n\treturn wrapArrayListOwn(items)\n}\n\n");
          return true;
        }
        if (simple.equals("toArray") && arity == 1) {
          out.append("\tCheckNotNull(arg0)\n\tout := make([]uint16, 0, arg0.Size())\n\tfor it := Iterate(arg0); it.HasNext(); { out = append(out, UnboxChar(Cast[*BoxChar](it.Next()))) }\n\treturn out\n}\n\n");
          return true;
        }
      }
      default -> {}
    }

    // --- java.util.Comparator / Map.Entry / Function / Predicate ---
    if (qn.equals("java.util.Comparator")) {
      if (simple.equals("naturalOrder") && arity == 0) {
        out.append("\treturn jNaturalOrderCmp\n}\n\n");
        return true;
      }
      if (simple.equals("comparing") && arity == 1) {
        out.append("\tCheckNotNull(arg0)\n\treturn newComparatorFromFunc(func(a, b any) int32 { return CompareValues(arg0.Apply(a), arg0.Apply(b)) })\n}\n\n");
        return true;
      }
      if (simple.equals("comparing") && arity == 2) {
        out.append("\tCheckNotNull(arg0); CheckNotNull(arg1)\n\treturn newComparatorFromFunc(func(a, b any) int32 { return arg1.Compare(arg0.Apply(a), arg0.Apply(b)) })\n}\n\n");
        return true;
      }
      if (simple.equals("comparingInt") && arity == 1) {
        out.append("\tCheckNotNull(arg0)\n\treturn newComparatorFromFunc(func(a, b any) int32 {\n\t\tx, y := arg0.ApplyAsInt(a), arg0.ApplyAsInt(b)\n\t\tif x < y { return -1 }\n\t\tif x > y { return 1 }\n\t\treturn 0\n\t})\n}\n\n");
        return true;
      }
      if (simple.equals("comparingLong") && arity == 1) {
        out.append("\tCheckNotNull(arg0)\n\treturn newComparatorFromFunc(func(a, b any) int32 {\n\t\tx, y := arg0.ApplyAsLong(a), arg0.ApplyAsLong(b)\n\t\tif x < y { return -1 }\n\t\tif x > y { return 1 }\n\t\treturn 0\n\t})\n}\n\n");
        return true;
      }
      if (simple.equals("reverseOrder") && arity == 0) {
        out.append("\treturn jReverseOrderCmp\n}\n\n");
        return true;
      }
    }
    if (qn.equals("java.util.Map.Entry") && simple.equals("comparingByKey") && arity == 0) {
      out.append("\treturn newComparatorFromFunc(func(a, b any) int32 {\n\t\tea := Cast[Ext_java_util_Map_Entry](a)\n\t\teb := Cast[Ext_java_util_Map_Entry](b)\n\t\treturn CompareValues(ea.GetKey(), eb.GetKey())\n\t})\n}\n\n");
      return true;
    }
    if (qn.equals("java.util.function.Function") && simple.equals("identity") && arity == 0) {
      out.append("\treturn &jFuncIdentity{}\n}\n\n");
      return true;
    }
    if (qn.equals("java.util.function.Predicate") && simple.equals("isEqual") && arity == 1) {
      out.append("\ttarget := ToAny(arg0)\n\treturn &jPredFunc{fn: func(x any) bool { return javalang.ObjectsEquals(target, ToAny(x)) }}\n}\n\n");
      return true;
    }

    // --- java.util.Optional, com.google.common.base.Optional, com.google.common.base.Ascii ---
    if (qn.equals("com.google.common.base.Optional")) {
      if (simple.equals("absent") && arity == 0) {
        out.append("\to := &Ext_com_google_common_base_OptionalBase{}\n\to.BindExtSelf_(o)\n\treturn o\n}\n\n");
        return true;
      }
      if (simple.equals("of") && arity == 1) {
        out.append("\to := &Ext_com_google_common_base_OptionalBase{Core_: CheckNotNull(ToAny(arg0))}\n\to.BindExtSelf_(o)\n\treturn o\n}\n\n");
        return true;
      }
      if (simple.equals("fromNullable") && arity == 1) {
        out.append("\tv := ToAny(arg0)\n\tif v == nil || javalang.IsNil(v) {\n\t\to := &Ext_com_google_common_base_OptionalBase{}\n\t\to.BindExtSelf_(o)\n\t\treturn o\n\t}\n\to := &Ext_com_google_common_base_OptionalBase{Core_: v}\n\to.BindExtSelf_(o)\n\treturn o\n}\n\n");
        return true;
      }
    }
    if (qn.equals("java.util.Optional")) {
      if (simple.equals("empty") && arity == 0) {
        out.append("\treturn &Ext_java_util_Optional{}\n}\n\n");
        return true;
      }
      if (simple.equals("of") && arity == 1) {
        out.append("\treturn &Ext_java_util_Optional{Core_: CheckNotNull(ToAny(arg0))}\n}\n\n");
        return true;
      }
      if (simple.equals("ofNullable") && arity == 1) {
        out.append("\tv := ToAny(arg0)\n\tif v == nil || javalang.IsNil(v) { return &Ext_java_util_Optional{} }\n\treturn &Ext_java_util_Optional{Core_: v}\n}\n\n");
        return true;
      }
    }
    if (qn.equals("com.google.common.base.Ascii")) {
      if (simple.equals("toLowerCase") && arity == 1 && p0.equals("uint16")) {
        out.append("\tif arg0 >= 'A' && arg0 <= 'Z' { return arg0 + 32 }\n\treturn arg0\n}\n\n");
        return true;
      }
      if (simple.equals("toUpperCase") && arity == 1 && p0.equals("uint16")) {
        out.append("\tif arg0 >= 'a' && arg0 <= 'z' { return arg0 - 32 }\n\treturn arg0\n}\n\n");
        return true;
      }
      if (simple.equals("toLowerCase") && arity == 1 && p0.equals("string")) {
        out.append("\treturn strings.ToLower(CheckNotNullString(arg0))\n}\n\n");
        return true;
      }
      if (simple.equals("toUpperCase") && arity == 1 && p0.equals("string")) {
        out.append("\treturn strings.ToUpper(CheckNotNullString(arg0))\n}\n\n");
        return true;
      }
      if (simple.equals("equalsIgnoreCase") && arity == 2) {
        out.append("\treturn strings.EqualFold(StrVal(CheckNotNull(arg0)), StrVal(CheckNotNull(arg1)))\n}\n\n");
        return true;
      }
    }

    // --- java.util.Arrays & java.util.Collections ---
    if (qn.equals("java.util.Arrays")) {
      switch (simple) {
        case "binarySearch" -> {
          if (arity == 2 && p0.equals("[]int32")) {
            out.append("\tCheckNotNullSlice(arg0)\n\tlo, hi := 0, len(arg0)-1\n\tfor lo <= hi {\n\t\tmid := int(uint(lo+hi) >> 1)\n\t\tv := arg0[mid]\n\t\tif v < arg1 { lo = mid + 1 } else if v > arg1 { hi = mid - 1 } else { return int32(mid) }\n\t}\n\treturn -int32(lo + 1)\n}\n\n");
            return true;
          }
        }
        case "asList" -> {
          out.append("\treturn newBackedArrayList(arg0)\n}\n\n");
          return true;
        }
        case "copyOf" -> {
          if (arity == 2 && p0.startsWith("[]")) {
            out.append("\tif arg0 == nil { ThrowNPE() }\n\tout := make(")
                .append(p0)
                .append(", arg1)\n\tcopy(out, arg0)\n\treturn out\n}\n\n");
            return true;
          }
        }
        case "copyOfRange" -> {
          if (arity == 3 && p0.startsWith("[]")) {
            out.append("\tif arg0 == nil { ThrowNPE() }\n\tif arg1 > arg2 { panic(javalang.NewIllegalArgumentException(\"fromIndex > toIndex\")) }\n\tif arg1 < 0 || int(arg1) > len(arg0) { panic(javalang.NewArrayIndexOutOfBoundsException(\"\")) }\n\tout := make(")
                .append(p0)
                .append(", arg2-arg1)\n\tcopy(out, arg0[arg1:])\n\treturn out\n}\n\n");
            return true;
          }
        }
        case "equals" -> {
          if (arity == 2 && p0.equals("[]int32")) {
            out.append("\tif arg0 == nil || arg1 == nil { return arg0 == nil && arg1 == nil }\n\tif len(arg0) != len(arg1) { return false }\n\tfor i := range arg0 { if arg0[i] != arg1[i] { return false } }\n\treturn true\n}\n\n");
            return true;
          }
        }
        case "fill" -> {
          if (arity == 2 && p0.startsWith("[]")) {
            out.append("\tif arg0 == nil { ThrowNPE() }\n");
            if (p0.equals("[]bool")) {
              out.append("\tif !arg1 { clear(arg0); return }\n");
            } else if (p0.equals("[]int32")) {
              out.append("\tif arg1 == 0 { clear(arg0); return }\n");
            }
            String valExpr = p0.equals("[]any") ? "ToAny(arg1)" : "arg1";
            out.append("\tfor i := range arg0 { arg0[i] = ").append(valExpr).append(" }\n}\n\n");
            return true;
          }
          if (arity == 4 && p0.startsWith("[]")) {
            out.append("\tif arg0 == nil { ThrowNPE() }\n\tCheckPositionIndexes(arg1, arg2, int32(len(arg0)))\n");
            if (p0.equals("[]bool")) {
              out.append("\tif !arg3 { clear(arg0[arg1:arg2]); return }\n");
            } else if (p0.equals("[]int32")) {
              out.append("\tif arg3 == 0 { clear(arg0[arg1:arg2]); return }\n");
            }
            String valExpr = p0.equals("[]any") ? "ToAny(arg3)" : "arg3";
            out.append("\tfor i := arg1; i < arg2; i++ { arg0[i] = ").append(valExpr).append(" }\n}\n\n");
            return true;
          }
        }
        case "hashCode" -> {
          if (arity == 1 && p0.equals("[]any")) {
            out.append("\tif arg0 == nil { return 0 }\n\tvar h int32 = 1\n\tfor _, e := range arg0 { h = 31*h + javalang.ObjectHashCode(ToAny(e)) }\n\treturn h\n}\n\n");
            return true;
          }
        }
        case "setAll" -> {
          if (arity == 2 && p0.equals("[]any")) {
            out.append("\tCheckNotNullSlice(arg0); CheckNotNull(arg1)\n\tfor i := range arg0 { arg0[i] = ToAny(arg1.Apply(int32(i))) }\n}\n\n");
            return true;
          }
        }
        case "sort" -> {
          if (arity == 1 && p0.equals("[]int32")) {
            out.append("\tCheckNotNullSlice(arg0)\n\tsort.Slice(arg0, func(i, j int) bool { return arg0[i] < arg0[j] })\n}\n\n");
            return true;
          }
          if (arity == 1 && p0.equals("[]int64")) {
            out.append("\tCheckNotNullSlice(arg0)\n\tsort.Slice(arg0, func(i, j int) bool { return arg0[i] < arg0[j] })\n}\n\n");
            return true;
          }
          if (arity == 1 && p0.equals("[]int8")) {
            out.append("\tCheckNotNullSlice(arg0)\n\tsort.Slice(arg0, func(i, j int) bool { return arg0[i] < arg0[j] })\n}\n\n");
            return true;
          }
          if (arity == 1 && p0.equals("[]any")) {
            out.append("\tCheckNotNullSlice(arg0)\n\tsort.SliceStable(arg0, func(i, j int) bool { return CompareValues(arg0[i], arg0[j]) < 0 })\n}\n\n");
            return true;
          }
          if (arity == 2 && p0.equals("[]any")) {
            out.append("\tCheckNotNullSlice(arg0)\n\tcmp := comparatorFuncOrNatural(arg1)\n\tsort.SliceStable(arg0, func(i, j int) bool { return cmp(arg0[i], arg0[j]) < 0 })\n}\n\n");
            return true;
          }
        }
        case "toString" -> {
          if (arity == 1 && p0.equals("[]uint16")) {
            out.append("\tif arg0 == nil { return \"null\" }\n\tvar sb strings.Builder\n\tsb.WriteByte('[')\n\tfor i, c := range arg0 {\n\t\tif i > 0 { sb.WriteString(\", \") }\n\t\tsb.WriteString(javalang.StringFromChar(c))\n\t}\n\tsb.WriteByte(']')\n\treturn sb.String()\n}\n\n");
            return true;
          }
          if (arity == 1 && p0.startsWith("[]")) {
            out.append("\treturn SliceFormatJavaString(arg0)\n}\n\n");
            return true;
          }
        }
        default -> {}
      }
    }

    if (qn.equals("java.util.Collections")) {
      switch (simple) {
        case "addAll" -> {
          if (arity == 2) {
            out.append("\tCheckNotNull(arg0)\n\tif arg1 == nil { ThrowNPE() }\n\tchanged := false\n\tfor _, e := range arg1 {\n\t\tif arg0.Add_1(ToAny(e)) { changed = true }\n\t}\n\treturn changed\n}\n\n");
            return true;
          }
        }
        case "emptyIterator" -> {
          out.append("\treturn newSliceObjIter(nil)\n}\n\n");
          return true;
        }
        case "emptyList" -> {
          out.append("\treturn emptyImmutableList_\n}\n\n");
          return true;
        }
        case "emptySet" -> {
          out.append("\treturn emptyImmutableSet_\n}\n\n");
          return true;
        }
        case "emptyMap" -> {
          out.append("\treturn emptyImmutableMap_\n}\n\n");
          return true;
        }
        case "singletonList" -> {
          out.append("\treturn newSingletonList(arg0)\n}\n\n");
          return true;
        }
        case "singleton" -> {
          out.append("\treturn newSingletonSet(arg0)\n}\n\n");
          return true;
        }
        case "singletonMap" -> {
          out.append("\tobj := newImmutableMapShell()\n\tobj.inlineCore.InitExpectedSize(1)\n\tobj.inlineCore.Put(arg0, arg1)\n\tobj.inlineCore.sealImmutable()\n\treturn obj\n}\n\n");
          return true;
        }
        case "newSetFromMap" -> {
          out.append("\treturn newSetFromMap(arg0)\n}\n\n");
          return true;
        }
        case "nCopies" -> {
          out.append("\tif arg0 < 0 { panic(javalang.NewIllegalArgumentException(\"List length < 0\")) }\n\tif arg0 == 0 { return wrapArrayList(nil) }\n\titems := make([]any, arg0)\n\tv := ToAny(arg1)\n\tfor i := range items { items[i] = v }\n\treturn wrapArrayListOwn(items)\n}\n\n");
          return true;
        }
        case "fill" -> {
          out.append("\tCheckNotNull(arg0)\n\tn := arg0.Size()\n\tv := ToAny(arg1)\n\tfor i := int32(0); i < n; i++ { arg0.Set(i, v) }\n}\n\n");
          return true;
        }
        case "reverse" -> {
          out.append("\tCheckNotNull(arg0)\n\tfor i, j := int32(0), arg0.Size()-1; i < j; i, j = i+1, j-1 {\n\t\ttmp := arg0.Get(i)\n\t\targ0.Set(i, arg0.Get(j))\n\t\targ0.Set(j, tmp)\n\t}\n}\n\n");
          return true;
        }
        case "reverseOrder" -> {
          if (arity == 0) {
            out.append("\treturn jReverseOrderCmp\n}\n\n");
            return true;
          }
          if (arity == 1) {
            out.append("\tif isJavaNilAny(arg0) { return jReverseOrderCmp }\n\treturn ExtDefault_Comparator_Reversed(arg0)\n}\n\n");
            return true;
          }
        }
        case "shuffle" -> {
          out.append("\tcollectionsShuffle(arg0, arg1)\n}\n\n");
          return true;
        }
        case "sort" -> {
          if (arity == 1) {
            out.append("\tCheckNotNull(arg0)\n\targ0.Sort(nil)\n}\n\n");
            return true;
          }
          if (arity == 2) {
            out.append("\tCheckNotNull(arg0)\n\targ0.Sort(arg1)\n}\n\n");
            return true;
          }
        }
        case "synchronizedList",
            "unmodifiableCollection",
            "unmodifiableList",
            "unmodifiableMap",
            "unmodifiableSet",
            "unmodifiableSortedSet" -> {
          out.append("\tCheckNotNull(arg0)\n\treturn arg0\n}\n\n");
          return true;
        }
        default -> {}
      }
    }

    // --- Streams & Collectors ---
    switch (qn) {
      case "java.util.stream.Stream" -> {
        if (simple.equals("builder") && arity == 0) {
          out.append("\treturn &jStreamBuilder{}\n}\n\n");
          return true;
        }
        if (simple.equals("concat") && arity == 2) {
          out.append("\treturn streamConcat(arg0, arg1)\n}\n\n");
          return true;
        }
        if (simple.equals("of") && arity == 1) {
          if (p0.startsWith("[]")) {
            out.append("\treturn streamFromSlice(arg0)\n}\n\n");
          } else {
            out.append("\treturn streamFromSlice([]any{ToAny(arg0)})\n}\n\n");
          }
          return true;
        }
      }
      case "java.util.stream.IntStream" -> {
        if (simple.equals("range") && arity == 2) {
          out.append("\treturn &jIntStream{start: arg0, end: arg1}\n}\n\n");
          return true;
        }
      }
      case "java.util.stream.Collectors" -> {
        switch (simple) {
          case "toList" -> {
            out.append("\treturn collectorToList()\n}\n\n");
            return true;
          }
          case "toCollection" -> {
            out.append("\treturn collectorToCollection(arg0)\n}\n\n");
            return true;
          }
          case "joining" -> {
            out.append("\treturn collectorJoining(arg0)\n}\n\n");
            return true;
          }
          case "groupingBy" -> {
            out.append("\treturn collectorGroupingBy(arg0)\n}\n\n");
            return true;
          }
          case "toMap" -> {
            out.append("\treturn collectorToMap(arg0, arg1)\n}\n\n");
            return true;
          }
          default -> {}
        }
      }
      case "com.google.common.collect.Streams" -> {
        if (simple.equals("stream") && arity == 1) {
          out.append("\treturn streamFromIterable(arg0)\n}\n\n");
          return true;
        }
        if (simple.equals("mapWithIndex") && arity == 2) {
          out.append("\treturn streamMapWithIndex(arg0, arg1)\n}\n\n");
          return true;
        }
      }
      default -> {}
    }

    // --- Cache & Concurrent ---
    if (qn.equals("com.google.common.cache.CacheBuilder")
        && (simple.equals("newBuilder") || simple.equals("from"))) {
      out.append("\tb := &Ext_com_google_common_cache_CacheBuilder{}\n\tb.Core_ = &jCacheBuilderCore{maxSize: -1}\n\treturn b\n}\n\n");
      return true;
    }
    if (qn.equals("com.google.common.cache.CacheLoader") && simple.equals("from") && arity == 1) {
      out.append("\tCheckNotNull(arg0)\n\treturn &jFuncCacheLoader{fn: arg0}\n}\n\n");
      return true;
    }
    if (qn.equals("java.util.concurrent.Executors")
        && simple.equals("newSingleThreadExecutor")
        && arity == 1) {
      out.append("\t_ = arg0\n\treturn &jDirectExecutor{}\n}\n\n");
      return true;
    }
    if (qn.equals("com.google.common.util.concurrent.MoreExecutors")
        && simple.equals("listeningDecorator")
        && arity == 1) {
      out.append("\tCheckNotNull(arg0)\n\treturn &jListeningExecWrapper{exec: arg0}\n}\n\n");
      return true;
    }
    if (qn.equals("com.google.common.util.concurrent.Futures")
        && simple.equals("allAsList")
        && arity == 1) {
      out.append("\treturn futuresAllAsList(arg0)\n}\n\n");
      return true;
    }

    return false;
  }

  public static String tryBuildInstanceMethodBody(
      RuntimeBindings rb,
      TypeElement te,
      TypeResolver.MethodSig sig,
      PackagePlanner.ImportTracker imports,
      String retType) {
    String qn = te.getQualifiedName().toString();
    String m = sig.simpleName;
    int arity = sig.canonParamTypes.size();
    String p0 =
        arity >= 1 ? rb.goTypeInJ2gort(sig.canonParamTypes.get(0), imports) : "";

    // 0. java.util.ArrayDeque backed by O(1) circular ring buffer *JDequeCore
    if (qn.equals("java.util.ArrayDeque")) {
      switch (m) {
        case "size" -> {
          if (arity == 0) return "\treturn this_.dc().Size()\n";
        }
        case "isEmpty" -> {
          if (arity == 0) return "\treturn this_.dc().IsEmpty()\n";
        }
        case "addFirst", "push" -> {
          if (arity == 1) return "\tthis_.dc().AddFirst(arg0)\n";
        }
        case "offerFirst" -> {
          if (arity == 1) return "\treturn this_.dc().OfferFirst(arg0)\n";
        }
        case "addLast" -> {
          if (arity == 1) return "\tthis_.dc().AddLast(arg0)\n";
        }
        case "add", "offer", "offerLast" -> {
          if (arity == 1) {
            return retType.isEmpty()
                ? "\tthis_.dc().AddLast(arg0)\n"
                : "\treturn this_.dc().OfferLast(arg0)\n";
          }
        }
        case "removeFirst", "pop" -> {
          if (arity == 0) return "\treturn this_.dc().RemoveFirst()\n";
        }
        case "remove" -> {
          if (arity == 0) return "\treturn this_.dc().RemoveFirst()\n";
          if (arity == 1) return "\treturn this_.dc().RemoveFirstOccurrence(arg0)\n";
        }
        case "poll", "pollFirst" -> {
          if (arity == 0) return "\treturn this_.dc().PollFirst()\n";
        }
        case "removeLast" -> {
          if (arity == 0) return "\treturn this_.dc().RemoveLast()\n";
        }
        case "pollLast" -> {
          if (arity == 0) return "\treturn this_.dc().PollLast()\n";
        }
        case "getFirst", "element" -> {
          if (arity == 0) return "\treturn this_.dc().GetFirst()\n";
        }
        case "peek", "peekFirst" -> {
          if (arity == 0) return "\treturn this_.dc().PeekFirst()\n";
        }
        case "getLast" -> {
          if (arity == 0) return "\treturn this_.dc().GetLast()\n";
        }
        case "peekLast" -> {
          if (arity == 0) return "\treturn this_.dc().PeekLast()\n";
        }
        case "removeFirstOccurrence" -> {
          if (arity == 1) return "\treturn this_.dc().RemoveFirstOccurrence(arg0)\n";
        }
        case "removeLastOccurrence" -> {
          if (arity == 1) return "\treturn this_.dc().RemoveLastOccurrence(arg0)\n";
        }
        case "contains" -> {
          if (arity == 1) return "\treturn this_.dc().Contains(arg0)\n";
        }
        case "containsAll" -> {
          if (arity == 1) return "\treturn colContainsAll(this_.dc().Contains, arg0)\n";
        }
        case "iterator" -> {
          if (arity == 0) return "\treturn newDequeIter(this_.dc())\n";
        }
        case "descendingIterator" -> {
          if (arity == 0) return "\treturn newDequeDescIter(this_.dc())\n";
        }
        case "toArray" -> {
          if (arity == 0) return "\treturn this_.dc().ToSlice()\n";
          if (arity == 1) return "\treturn sliceToArray1(this_.dc().ToSlice(), arg0)\n";
        }
        case "addAll" -> {
          if (arity == 1) {
            return "\tCheckNotNull(arg0)\n\tadded := false\n\tfor it := Iterate(arg0); it.HasNext(); { this_.dc().AddLast(it.Next()); added = true }\n\t"
                + (retType.equals("bool")
                    ? "return added\n"
                    : (retType.isEmpty() ? "" : "return this_\n"));
          }
        }
        case "removeAll" -> {
          if (arity == 1) {
            return "\tCheckNotNull(arg0)\n\treturn this_.dc().RemoveIf(func(x any) bool { return arg0.Contains(x) })\n";
          }
        }
        case "retainAll" -> {
          if (arity == 1) {
            return "\tCheckNotNull(arg0)\n\treturn this_.dc().RemoveIf(func(x any) bool { return !arg0.Contains(x) })\n";
          }
        }
        case "removeIf" -> {
          if (arity == 1) {
            return "\tCheckNotNull(arg0)\n\treturn this_.dc().RemoveIf(func(x any) bool { return arg0.Test(x) })\n";
          }
        }
        case "forEach" -> {
          if (arity == 1) {
            return "\tCheckNotNull(arg0)\n\tthis_.dc().ForEach(func(x any) { arg0.Accept(x) })\n";
          }
        }
        case "stream" -> {
          if (arity == 0) return "\treturn streamFromSlice(this_.dc().ToSlice())\n";
        }
        case "clear" -> {
          if (arity == 0) return "\tthis_.dc().Clear()\n";
        }
        case "clone" -> {
          if (arity == 0) {
            String st = rb.goStructType(te, imports);
            return "\tout := &"
                + st
                + "{}\n\tout.BindExtSelf_(out)\n\tthis_.dc().CloneInto(out.dc())\n\treturn out\n";
          }
        }
        case "equals" -> {
          if (arity == 1) return "\treturn RefEquals(this_, arg0)\n";
        }
        case "hashCode" -> {
          if (arity == 0) return "\treturn javalang.IdentityHashCode(this_)\n";
        }
        case "toString" -> {
          if (arity == 0) return "\treturn sliceJavaString(this_.dc().ToSlice())\n";
        }
        default -> {}
      }
    }

    // 0b. com.google.common.collect.ImmutableList backed directly by flat items_ []any / singleArr_ [1]any
    if (qn.equals("com.google.common.collect.ImmutableList")) {
      switch (m) {
        case "size" -> {
          if (arity == 0) return "\treturn int32(len(this_.items_))\n";
        }
        case "isEmpty" -> {
          if (arity == 0) return "\treturn len(this_.items_) == 0\n";
        }
        case "get" -> {
          if (arity == 1) {
            return "\tCheckElementIndex(arg0, int32(len(this_.items_)))\n\treturn this_.items_[arg0]\n";
          }
        }
        case "getFirst", "element" -> {
          if (arity == 0) {
            return "\tif len(this_.items_) == 0 { panic(javalang.NewNoSuchElementException(\"empty\")) }\n\treturn this_.items_[0]\n";
          }
        }
        case "getLast" -> {
          if (arity == 0) {
            return "\tif len(this_.items_) == 0 { panic(javalang.NewNoSuchElementException(\"empty\")) }\n\treturn this_.items_[len(this_.items_)-1]\n";
          }
        }
        case "contains" -> {
          if (arity == 1) return "\treturn immSliceContains(this_.items_, arg0)\n";
        }
        case "indexOf" -> {
          if (arity == 1) return "\treturn immSliceIndexOf(this_.items_, arg0)\n";
        }
        case "lastIndexOf" -> {
          if (arity == 1) return "\treturn immSliceLastIndexOf(this_.items_, arg0)\n";
        }
        case "containsAll" -> {
          if (arity == 1) {
            return "\treturn colContainsAll(func(x any) bool { return immSliceContains(this_.items_, x) }, arg0)\n";
          }
        }
        case "iterator" -> {
          if (arity == 0) return "\treturn newSliceObjIter(this_.items_)\n";
        }
        case "listIterator" -> {
          if (arity == 0) return "\treturn newImmListIter(this_.items_, 0)\n";
          if (arity == 1) return "\treturn newImmListIter(this_.items_, arg0)\n";
        }
        case "subList" -> {
          if (arity == 2) {
            return "\tn := int32(len(this_.items_))\n\tCheckPositionIndexes(arg0, arg1, n)\n\tif arg0 == 0 && arg1 == n { return this_.selfImmList() }\n\treturn wrapImmutableList(this_.items_[arg0:arg1])\n";
          }
        }
        case "toArray" -> {
          if (arity == 0) {
            return "\tout := make([]any, len(this_.items_))\n\tcopy(out, this_.items_)\n\treturn out\n";
          }
          if (arity == 1) return "\treturn sliceToArray1(this_.items_, arg0)\n";
        }
        case "forEach" -> {
          if (arity == 1) {
            return "\tCheckNotNull(arg0)\n\tfor _, x := range this_.items_ { arg0.Accept(x) }\n";
          }
        }
        case "stream" -> {
          if (arity == 0) return "\treturn streamFromSlice(this_.items_)\n";
        }
        case "asList" -> {
          if (arity == 0) return "\treturn this_.selfImmList()\n";
        }
        case "reverse" -> {
          if (arity == 0) {
            return "\tsrc := this_.items_\n\tif len(src) <= 1 { return this_.selfImmList() }\n\trev := make([]any, len(src))\n\tfor i, v := range src { rev[len(src)-1-i] = v }\n\treturn wrapImmutableListOwn(rev)\n";
          }
        }
        case "equals" -> {
          if (arity == 1) return "\treturn listEqualsSlice(this_.items_, arg0)\n";
        }
        case "hashCode" -> {
          if (arity == 0) return "\treturn listHashSlice(this_.items_)\n";
        }
        case "toString" -> {
          if (arity == 0) return "\treturn sliceJavaString(this_.items_)\n";
        }
        case "set",
            "add",
            "addFirst",
            "addLast",
            "remove",
            "removeFirst",
            "removeLast",
            "clear",
            "addAll",
            "removeAll",
            "retainAll",
            "removeIf",
            "sort",
            "replaceAll" -> {
          return "\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n";
        }
        default -> {}
      }
    }

    // 1. Concrete List / Deque / Queue types backed by *JListCore
    if (CONCRETE_LIST_TYPES.contains(qn)) {
      switch (m) {
        case "size" -> {
          if (arity == 0) return "\treturn this_.lc().Size()\n";
        }
        case "isEmpty" -> {
          if (arity == 0) return "\treturn this_.lc().IsEmpty()\n";
        }
        case "get" -> {
          if (arity == 1) return "\treturn this_.lc().Get(arg0)\n";
        }
        case "set" -> {
          if (arity == 2) return "\treturn this_.lc().Set(arg0, arg1)\n";
        }
        case "add", "addLast", "offer", "offerLast" -> {
          if (arity == 1) {
            return retType.isEmpty()
                ? "\tthis_.lc().Add(arg0)\n"
                : (retType.equals("bool")
                    ? "\treturn this_.lc().Add(arg0)\n"
                    : "\tthis_.lc().Add(arg0)\n\treturn this_\n");
          }
          if (arity == 2 && m.equals("add")) {
            return "\tthis_.lc().AddAt(arg0, arg1)\n";
          }
        }
        case "push", "addFirst" -> {
          if (arity == 1) return "\tthis_.lc().AddAt(0, arg0)\n";
        }
        case "offerFirst" -> {
          if (arity == 1) return "\tthis_.lc().AddAt(0, arg0)\n\treturn true\n";
        }
        case "pop", "removeFirst" -> {
          if (arity == 0) {
            return "\tif this_.lc().IsEmpty() { panic(javalang.NewNoSuchElementException(\"empty\")) }\n\treturn this_.lc().RemoveAt(0)\n";
          }
        }
        case "remove" -> {
          if (arity == 0) {
            return "\tif this_.lc().IsEmpty() { panic(javalang.NewNoSuchElementException(\"empty\")) }\n\treturn this_.lc().RemoveAt(0)\n";
          }
          if (arity == 1) {
            return p0.equals("int32")
                ? "\treturn this_.lc().RemoveAt(arg0)\n"
                : "\treturn this_.lc().RemoveVal(arg0)\n";
          }
        }
        case "removeFirstOccurrence" -> {
          if (arity == 1) return "\treturn this_.lc().RemoveVal(arg0)\n";
        }
        case "removeLastOccurrence" -> {
          if (arity == 1) return "\treturn listCoreRemoveLastVal(this_.lc(), arg0)\n";
        }
        case "poll", "pollFirst" -> {
          if (arity == 0) {
            return "\tif this_.lc().IsEmpty() { return nil }\n\treturn this_.lc().RemoveAt(0)\n";
          }
        }
        case "pollLast" -> {
          if (arity == 0) {
            return "\tif this_.lc().IsEmpty() { return nil }\n\treturn this_.lc().RemoveAt(this_.lc().Size() - 1)\n";
          }
        }
        case "peek", "peekFirst" -> {
          if (arity == 0) {
            return "\tif this_.lc().IsEmpty() { return nil }\n\treturn this_.lc().Get(0)\n";
          }
        }
        case "getFirst", "element" -> {
          if (arity == 0) {
            return "\tif this_.lc().IsEmpty() { panic(javalang.NewNoSuchElementException(\"empty\")) }\n\treturn this_.lc().Get(0)\n";
          }
        }
        case "peekLast" -> {
          if (arity == 0) {
            return "\tif this_.lc().IsEmpty() { return nil }\n\treturn this_.lc().Get(this_.lc().Size() - 1)\n";
          }
        }
        case "getLast" -> {
          if (arity == 0) {
            return "\tif this_.lc().IsEmpty() { panic(javalang.NewNoSuchElementException(\"empty\")) }\n\treturn this_.lc().Get(this_.lc().Size() - 1)\n";
          }
        }
        case "removeLast" -> {
          if (arity == 0) {
            return "\tif this_.lc().IsEmpty() { panic(javalang.NewNoSuchElementException(\"empty\")) }\n\treturn this_.lc().RemoveAt(this_.lc().Size() - 1)\n";
          }
        }
        case "contains" -> {
          if (arity == 1) return "\treturn this_.lc().Contains(arg0)\n";
        }
        case "indexOf" -> {
          if (arity == 1) return "\treturn this_.lc().IndexOf(arg0)\n";
        }
        case "lastIndexOf" -> {
          if (arity == 1) return "\treturn this_.lc().LastIndexOf(arg0)\n";
        }
        case "clear" -> {
          if (arity == 0) return "\tthis_.lc().Clear()\n";
        }
        case "iterator" -> {
          if (arity == 0) return "\treturn newListIter(this_.lc())\n";
        }
        case "descendingIterator" -> {
          if (arity == 0) return "\treturn newDescendingListIter(this_.lc())\n";
        }
        case "listIterator" -> {
          if (arity == 0) return "\treturn newFullListIter(this_.lc(), 0)\n";
          if (arity == 1) return "\treturn newFullListIter(this_.lc(), arg0)\n";
        }
        case "subList" -> {
          if (arity == 2) return "\treturn newSubListView(this_.lc(), arg0, arg1)\n";
        }
        case "addAll" -> {
          if (arity == 1) {
            return "\tadded := listCoreAddAllIterable(this_.lc(), arg0)\n\t"
                + (retType.equals("bool") ? "return added\n" : (retType.isEmpty() ? "_ = added\n" : "_ = added\n\treturn this_\n"));
          }
        }
        case "containsAll" -> {
          if (arity == 1) return "\treturn colContainsAll(this_.lc().Contains, arg0)\n";
        }
        case "removeAll" -> {
          if (arity == 1) return "\treturn listCoreRemoveAll(this_.lc(), arg0)\n";
        }
        case "retainAll" -> {
          if (arity == 1) {
            return "\tCheckNotNull(arg0)\n\treturn listCoreRemoveIf(this_.lc(), func(x any) bool { return !arg0.Contains(x) })\n";
          }
        }
        case "removeIf" -> {
          if (arity == 1) {
            return "\tCheckNotNull(arg0)\n\treturn listCoreRemoveIf(this_.lc(), func(x any) bool { return arg0.Test(x) })\n";
          }
        }
        case "toArray" -> {
          if (arity == 0) return "\treturn this_.lc().ToArray()\n";
          if (arity == 1) return "\treturn sliceToArray1(this_.lc().Items, arg0)\n";
        }
        case "sort" -> {
          if (arity == 1) return "\tlistCoreSort(this_.lc(), arg0)\n";
        }
        case "forEach" -> {
          if (arity == 1) {
            return "\tCheckNotNull(arg0)\n\tfor _, x := range this_.lc().Items { arg0.Accept(x) }\n";
          }
        }
        case "stream" -> {
          if (arity == 0) return "\treturn streamFromSlice(this_.lc().ToArray())\n";
        }
        case "asList" -> {
          if (arity == 0) return "\treturn wrapImmutableList(this_.lc().Items)\n";
        }
        case "clone" -> {
          if (arity == 0) {
            String st = rb.goStructType(te, imports);
            if (qn.equals("java.util.ArrayList")) {
              return "\tout := &"
                  + st
                  + "{}\n\tout.BindExtSelf_(out)\n\tout.inlineCore.InitFromSlice(this_.lc().Items)\n\treturn out\n";
            }
            return "\tout := &"
                + st
                + "{}\n\tout.BindExtSelf_(out)\n\tlc := NewJListCore()\n\tlc.InitFromSlice(this_.lc().Items)\n\tout.Core_ = lc\n\treturn out\n";
          }
        }
        case "equals" -> {
          if (arity == 1) {
            if (qn.equals("java.util.concurrent.LinkedBlockingQueue")) {
              return "\treturn RefEquals(this_, arg0)\n";
            }
            return "\treturn listEqualsSlice(this_.lc().Items, arg0)\n";
          }
        }
        case "hashCode" -> {
          if (arity == 0) {
            if (qn.equals("java.util.concurrent.LinkedBlockingQueue")) {
              return "\treturn javalang.IdentityHashCode(this_)\n";
            }
            return "\treturn listHashSlice(this_.lc().Items)\n";
          }
        }
        case "toString" -> {
          if (arity == 0) return "\treturn sliceJavaString(this_.lc().Items)\n";
        }
        default -> {}
      }
    }

    // 2a. java.util.EnumSet backed by bitmask jEnumSetCore
    if (qn.equals("java.util.EnumSet")) {
      switch (m) {
        case "size" -> {
          if (arity == 0) return "\treturn this_.esCore().Size()\n";
        }
        case "isEmpty" -> {
          if (arity == 0) return "\treturn this_.esCore().Size() == 0\n";
        }
        case "contains" -> {
          if (arity == 1) return "\treturn this_.esCore().Contains(arg0)\n";
        }
        case "add" -> {
          if (arity == 1) return "\treturn this_.esCore().Add(arg0)\n";
        }
        case "remove" -> {
          if (arity == 1) return "\treturn this_.esCore().Remove(arg0)\n";
        }
        case "containsAll" -> {
          if (arity == 1) return "\treturn this_.esCore().ContainsAll(arg0)\n";
        }
        case "addAll" -> {
          if (arity == 1) return "\treturn this_.esCore().AddAll(arg0)\n";
        }
        case "removeAll" -> {
          if (arity == 1) return "\treturn this_.esCore().RemoveAll(arg0)\n";
        }
        case "retainAll" -> {
          if (arity == 1) return "\treturn this_.esCore().RetainAll(arg0)\n";
        }
        case "removeIf" -> {
          if (arity == 1) return "\treturn this_.esCore().RemoveIf(arg0)\n";
        }
        case "clear" -> {
          if (arity == 0) return "\tthis_.esCore().Clear()\n";
        }
        case "iterator" -> {
          if (arity == 0) return "\treturn this_.esCore().Iterator()\n";
        }
        case "toArray" -> {
          if (arity == 0) return "\treturn this_.esCore().ToArray()\n";
          if (arity == 1) return "\treturn sliceToArray1(this_.esCore().ToArray(), arg0)\n";
        }
        case "forEach" -> {
          if (arity == 1) return "\tthis_.esCore().ForEachConsumer(arg0)\n";
        }
        case "stream" -> {
          if (arity == 0) return "\treturn streamFromSlice(this_.esCore().ToArray())\n";
        }
        case "clone" -> {
          if (arity == 0) return "\treturn enumSetClone(this_.esCore())\n";
        }
        case "equals" -> {
          if (arity == 1) return "\treturn this_.esCore().Equals(arg0)\n";
        }
        case "hashCode" -> {
          if (arity == 0) return "\treturn this_.esCore().HashCode()\n";
        }
        case "toString" -> {
          if (arity == 0) return "\treturn this_.esCore().String()\n";
        }
        default -> {}
      }
    }

    // 2b. Flat ImmutableSet & ImmutableSet.CachingAsList backed by elems_ / strs_ / optional big_ *JSetCore
    if (qn.equals("com.google.common.collect.ImmutableSet")
        || qn.equals("com.google.common.collect.ImmutableSet.CachingAsList")) {
      switch (m) {
        case "size" -> {
          if (arity == 0) {
            return "\tif this_.strs_ != nil { return int32(len(this_.strs_)) }\n\treturn int32(len(this_.elems_))\n";
          }
        }
        case "isEmpty" -> {
          if (arity == 0) return "\treturn len(this_.elems_) == 0 && len(this_.strs_) == 0\n";
        }
        case "first", "getFirst" -> {
          if (arity == 0) {
            return "\tif this_.strs_ != nil {\n\t\tif len(this_.strs_) == 0 { panic(javalang.NewNoSuchElementException(\"empty\")) }\n\t\treturn this_.strs_[0]\n\t}\n\tif len(this_.elems_) == 0 { panic(javalang.NewNoSuchElementException(\"empty\")) }\n\treturn this_.elems_[0]\n";
          }
        }
        case "last", "getLast" -> {
          if (arity == 0) {
            return "\tif this_.strs_ != nil {\n\t\tif len(this_.strs_) == 0 { panic(javalang.NewNoSuchElementException(\"empty\")) }\n\t\treturn this_.strs_[len(this_.strs_)-1]\n\t}\n\tif len(this_.elems_) == 0 { panic(javalang.NewNoSuchElementException(\"empty\")) }\n\treturn this_.elems_[len(this_.elems_)-1]\n";
          }
        }
        case "contains" -> {
          if (arity == 1) return "\treturn this_.immSetContains(arg0)\n";
        }
        case "containsAll" -> {
          if (arity == 1) return "\treturn colContainsAll(this_.immSetContains, arg0)\n";
        }
        case "iterator" -> {
          if (arity == 0) {
            return "\tif this_.strs_ != nil { return newStrSliceIter(this_.strs_) }\n\treturn newSliceObjIter(this_.elems_)\n";
          }
        }
        case "toArray" -> {
          if (arity == 0) {
            return "\tif this_.strs_ != nil { return strsToAnySlice(this_.strs_) }\n\tout := make([]any, len(this_.elems_))\n\tcopy(out, this_.elems_)\n\treturn out\n";
          }
          if (arity == 1) {
            return "\tif this_.strs_ != nil { return sliceToArray1(strsToAnySlice(this_.strs_), arg0) }\n\treturn sliceToArray1(this_.elems_, arg0)\n";
          }
        }
        case "asList" -> {
          if (arity == 0) {
            return "\tif this_.strs_ != nil { return wrapImmStrSliceCopy(this_.strs_) }\n\treturn wrapImmutableList(this_.elems_)\n";
          }
        }
        case "forEach" -> {
          if (arity == 1) {
            return "\tCheckNotNull(arg0)\n\tif this_.strs_ != nil {\n\t\tfor _, x := range this_.strs_ { arg0.Accept(x) }\n\t\treturn\n\t}\n\tfor _, x := range this_.elems_ { arg0.Accept(x) }\n";
          }
        }
        case "stream" -> {
          if (arity == 0) {
            return "\tif this_.strs_ != nil { return streamFromSlice(strsToAnySlice(this_.strs_)) }\n\treturn streamFromSlice(this_.elems_)\n";
          }
        }
        case "equals" -> {
          if (arity == 1) return "\treturn immSetEquals(this_.GetExt_com_google_common_collect_ImmutableSetBase(), arg0)\n";
        }
        case "hashCode" -> {
          if (arity == 0) {
            return "\tif this_.strs_ != nil { return immStrSetHash(this_.strs_) }\n\treturn immSetHash(this_.elems_)\n";
          }
        }
        case "toString" -> {
          if (arity == 0) {
            return "\tif this_.strs_ != nil { return strSliceJavaString(this_.strs_) }\n\treturn sliceJavaString(this_.elems_)\n";
          }
        }
        case "add", "remove", "addAll", "removeAll", "retainAll", "removeIf", "clear" -> {
          return "\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n";
        }
        default -> {}
      }
    }

    // 2c. HashSet & LinkedHashSet backed by dedicated JSetCore (with optional mapBack() for Sets.newSetFromMap)
    if (qn.equals("java.util.HashSet") || qn.equals("java.util.LinkedHashSet")) {
      String hsBase = qn.equals("java.util.HashSet") ? "this_" : "&this_.Ext_java_util_HashSetBase";
      switch (m) {
        case "size" -> {
          if (arity == 0) {
            return "\tif mb := this_.mapBack(); mb != nil { return mb.Size() }\n\treturn this_.inlineSet.Size()\n";
          }
        }
        case "isEmpty" -> {
          if (arity == 0) {
            return "\tif mb := this_.mapBack(); mb != nil { return mb.Size() == 0 }\n\treturn this_.inlineSet.IsEmpty()\n";
          }
        }
        case "contains" -> {
          if (arity == 1) {
            return "\tif mb := this_.mapBack(); mb != nil { return mb.ContainsKey(arg0) }\n\treturn this_.inlineSet.Contains(arg0)\n";
          }
        }
        case "add" -> {
          if (arity == 1) {
            return "\tif mb := this_.mapBack(); mb != nil { return mb.Put_2(arg0, ValueOfBoolean(true)) == nil }\n\treturn this_.inlineSet.Add(arg0)\n";
          }
        }
        case "remove" -> {
          if (arity == 1) {
            return "\tif mb := this_.mapBack(); mb != nil { return mb.Remove_1(arg0) != nil }\n\treturn this_.inlineSet.Remove(arg0)\n";
          }
        }
        case "addFirst" -> {
          if (arity == 1) return "\tthis_.inlineSet.AddFirst(arg0)\n";
        }
        case "addLast" -> {
          if (arity == 1) return "\tthis_.inlineSet.AddLast(arg0)\n";
        }
        case "first", "getFirst" -> {
          if (arity == 0) return "\treturn this_.inlineSet.FirstKey()\n";
        }
        case "last", "getLast" -> {
          if (arity == 0) return "\treturn this_.inlineSet.LastKey()\n";
        }
        case "removeFirst" -> {
          if (arity == 0) return "\treturn this_.inlineSet.RemoveFirstKey()\n";
        }
        case "removeLast" -> {
          if (arity == 0) return "\treturn this_.inlineSet.RemoveLastKey()\n";
        }
        case "addAll" -> {
          if (arity == 1) return "\treturn hashSetAddAll(" + hsBase + ", arg0)\n";
        }
        case "containsAll" -> {
          if (arity == 1) return "\treturn hashSetContainsAll(" + hsBase + ", arg0)\n";
        }
        case "removeAll" -> {
          if (arity == 1) return "\treturn hashSetRemoveAll(" + hsBase + ", arg0)\n";
        }
        case "retainAll" -> {
          if (arity == 1) return "\treturn hashSetRetainAll(" + hsBase + ", arg0)\n";
        }
        case "removeIf" -> {
          if (arity == 1) return "\treturn hashSetRemoveIf(" + hsBase + ", arg0)\n";
        }
        case "clear" -> {
          if (arity == 0) {
            return "\tif mb := this_.mapBack(); mb != nil { Cast[Ext_java_util_Collection](mb.KeySet_0()).Clear(); return }\n\tthis_.inlineSet.Clear()\n";
          }
        }
        case "iterator" -> {
          if (arity == 0) return "\treturn newHashSetIter(" + hsBase + ")\n";
        }
        case "toArray" -> {
          if (arity == 0) return "\treturn hashSetToArray(" + hsBase + ")\n";
          if (arity == 1) return "\treturn sliceToArray1(hashSetToArray(" + hsBase + "), arg0)\n";
        }
        case "forEach" -> {
          if (arity == 1) return "\thashSetForEach(" + hsBase + ", arg0)\n";
        }
        case "stream" -> {
          if (arity == 0) return "\treturn streamFromSlice(hashSetToArray(" + hsBase + "))\n";
        }
        case "clone" -> {
          if (arity == 0) {
            String st = rb.goStructType(te, imports);
            return "\tout := &"
                + st
                + "{}\n\tout.BindExtSelf_(out)\n\tcopySetCoreInto(&out.inlineSet, &this_.inlineSet)\n\treturn out\n";
          }
        }
        case "equals" -> {
          if (arity == 1) return "\treturn hashSetEquals(" + hsBase + ", arg0)\n";
        }
        case "hashCode" -> {
          if (arity == 0) return "\treturn hashSetHash(" + hsBase + ")\n";
        }
        case "toString" -> {
          if (arity == 0) return "\treturn sliceJavaString(hashSetToArray(" + hsBase + "))\n";
        }
        default -> {}
      }
    }

    // 2. Sorted Set types backed by *JMapCore (TreeSet, ImmutableSortedSet)
    if (CONCRETE_SET_TYPES.contains(qn)) {
      switch (m) {
        case "size" -> {
          if (arity == 0) return "\treturn this_.mc().Size()\n";
        }
        case "isEmpty" -> {
          if (arity == 0) return "\treturn this_.mc().IsEmpty()\n";
        }
        case "contains" -> {
          if (arity == 1) return "\treturn this_.mc().ContainsKey(arg0)\n";
        }
        case "add" -> {
          if (arity == 1) {
            if (qn.startsWith("com.google.")) {
              return "\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n";
            }
            return "\treturn this_.mc().AddKey(arg0)\n";
          }
        }
        case "addFirst" -> {
          if (arity == 1) return "\tthis_.mc().AddFirstKey(arg0, true)\n";
        }
        case "addLast" -> {
          if (arity == 1) return "\tthis_.mc().AddLastKey(arg0, true)\n";
        }
        case "first", "getFirst" -> {
          if (arity == 0) return "\treturn this_.mc().FirstKey()\n";
        }
        case "last", "getLast" -> {
          if (arity == 0) return "\treturn this_.mc().LastKey()\n";
        }
        case "removeFirst" -> {
          if (arity == 0) return "\treturn this_.mc().RemoveFirstKey()\n";
        }
        case "removeLast" -> {
          if (arity == 0) return "\treturn this_.mc().RemoveLastKey()\n";
        }
        case "pollFirst" -> {
          if (arity == 0) {
            return qn.startsWith("com.google.")
                ? "\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n"
                : "\tif this_.mc().IsEmpty() { return nil }\n\treturn this_.mc().RemoveFirstKey()\n";
          }
        }
        case "pollLast" -> {
          if (arity == 0) {
            return qn.startsWith("com.google.")
                ? "\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n"
                : "\tif this_.mc().IsEmpty() { return nil }\n\treturn this_.mc().RemoveLastKey()\n";
          }
        }
        case "lower" -> {
          if (arity == 1) return "\treturn sortedCoreNavKey(this_.mc(), arg0, -2)\n";
        }
        case "floor" -> {
          if (arity == 1) return "\treturn sortedCoreNavKey(this_.mc(), arg0, -1)\n";
        }
        case "ceiling" -> {
          if (arity == 1) return "\treturn sortedCoreNavKey(this_.mc(), arg0, 1)\n";
        }
        case "higher" -> {
          if (arity == 1) return "\treturn sortedCoreNavKey(this_.mc(), arg0, 2)\n";
        }
        case "addAll" -> {
          if (arity == 1) {
            if (qn.startsWith("com.google.")) {
              return "\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n";
            }
            return "\treturn mapCoreAddAllIterable(this_.mc(), arg0)\n";
          }
        }
        case "remove" -> {
          if (arity == 1) {
            if (qn.startsWith("com.google.")) {
              return "\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n";
            }
            return "\treturn this_.mc().RemoveKey(arg0)\n";
          }
        }
        case "containsAll" -> {
          if (arity == 1) return "\treturn colContainsAll(this_.mc().ContainsKey, arg0)\n";
        }
        case "removeAll" -> {
          if (arity == 1) return "\treturn mapCoreRemoveAllKeys(this_.mc(), arg0)\n";
        }
        case "retainAll" -> {
          if (arity == 1) {
            if (qn.startsWith("com.google.")) {
              return "\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n";
            }
            return "\treturn mapCoreRetainAllKeys(this_.mc(), arg0)\n";
          }
        }
        case "removeIf" -> {
          if (arity == 1) {
            if (qn.startsWith("com.google.")) {
              return "\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n";
            }
            return "\tCheckNotNull(arg0)\n\treturn mapCoreRemoveIfKeys(this_.mc(), arg0.Test)\n";
          }
        }
        case "clear" -> {
          if (arity == 0) return "\tthis_.mc().Clear()\n";
        }
        case "iterator" -> {
          if (arity == 0) return "\treturn newMapKeyIter(this_.mc())\n";
        }
        case "toArray" -> {
          if (arity == 0) return "\treturn this_.mc().LiveKeys()\n";
          if (arity == 1) return "\treturn sliceToArray1(this_.mc().LiveKeys(), arg0)\n";
        }
        case "asList" -> {
          if (arity == 0) return "\treturn wrapImmutableListOwn(this_.mc().LiveKeys())\n";
        }
        case "descendingSet" -> {
          if (arity == 0) return "\treturn setCoreDescending(this_.mc(), " + qn.startsWith("com.google.") + ")\n";
        }
        case "tailSet" -> {
          if (arity == 2) return "\treturn setCoreTailSet(this_.mc(), arg0, arg1, " + qn.startsWith("com.google.") + ")\n";
        }
        case "forEach" -> {
          if (arity == 1) {
            if (qn.startsWith("com.google.")) {
              return "\tCheckNotNull(arg0)\n\tthis_.mc().ForEachEntry(func(k, _ any) bool { arg0.Accept(k); return true })\n";
            }
            return "\tCheckNotNull(arg0)\n\tfor _, k := range this_.mc().LiveKeys() { arg0.Accept(k) }\n";
          }
        }
        case "stream" -> {
          if (arity == 0) return "\treturn streamFromSlice(this_.mc().LiveKeys())\n";
        }
        case "clone" -> {
          if (arity == 0 && !qn.startsWith("com.google.")) {
            String st = rb.goStructType(te, imports);
            if (qn.equals("java.util.TreeSet")) {
              return "\tout := &"
                  + st
                  + "{}\n\tout.BindExtSelf_(out)\n\tcopySortedCoreInto(&out.inlineCore, this_.mc())\n\treturn out\n";
            }
            return "\tout := &"
                + st
                + "{}\n\tout.BindExtSelf_(out)\n\tcopyHashCoreInto(&out.inlineCore, this_.mc())\n\treturn out\n";
          }
        }
        case "equals" -> {
          if (arity == 1) return "\treturn setCoreEquals(this_.mc(), arg0)\n";
        }
        case "hashCode" -> {
          if (arity == 0) {
            return qn.startsWith("com.google.")
                ? "\treturn this_.mc().ImmutableSetHash()\n"
                : "\treturn setCoreHash(this_.mc())\n";
          }
        }
        case "toString" -> {
          if (arity == 0) return "\treturn sliceJavaString(this_.mc().LiveKeys())\n";
        }
        default -> {}
      }
    }

    // 3a. java.util.EnumMap backed by slice/inline bitmask jEnumMapCore
    if (qn.equals("java.util.EnumMap")) {
      switch (m) {
        case "size" -> {
          if (arity == 0) return "\treturn this_.emCore().Size()\n";
        }
        case "isEmpty" -> {
          if (arity == 0) return "\treturn this_.emCore().Size() == 0\n";
        }
        case "containsKey" -> {
          if (arity == 1) return "\treturn this_.emCore().ContainsKey(arg0)\n";
        }
        case "containsValue" -> {
          if (arity == 1) return "\treturn this_.emCore().ContainsValue(arg0)\n";
        }
        case "get" -> {
          if (arity == 1) return "\treturn this_.emCore().Get(arg0)\n";
        }
        case "getOrDefault" -> {
          if (arity == 2) return "\treturn this_.emCore().GetOrDefault(arg0, arg1)\n";
        }
        case "put" -> {
          if (arity == 2) return "\treturn this_.emCore().Put(arg0, arg1)\n";
        }
        case "putIfAbsent" -> {
          if (arity == 2) return "\treturn this_.emCore().PutIfAbsent(arg0, arg1)\n";
        }
        case "remove" -> {
          if (arity == 1) return "\treturn this_.emCore().Remove(arg0)\n";
          if (arity == 2) return "\treturn this_.emCore().Remove2(arg0, arg1)\n";
        }
        case "putAll" -> {
          if (arity == 1) return "\tthis_.emCore().PutAll(arg0)\n";
        }
        case "clear" -> {
          if (arity == 0) return "\tthis_.emCore().Clear()\n";
        }
        case "replace" -> {
          if (arity == 2) return "\treturn this_.emCore().Replace2(arg0, arg1)\n";
          if (arity == 3) return "\treturn this_.emCore().Replace3(arg0, arg1, arg2)\n";
        }
        case "replaceAll" -> {
          if (arity == 1) return "\tthis_.emCore().ReplaceAll(arg0)\n";
        }
        case "computeIfAbsent" -> {
          if (arity == 2) return "\treturn this_.emCore().ComputeIfAbsent(arg0, arg1)\n";
        }
        case "computeIfPresent" -> {
          if (arity == 2) return "\treturn this_.emCore().ComputeIfPresent(arg0, arg1)\n";
        }
        case "compute" -> {
          if (arity == 2) return "\treturn this_.emCore().Compute(arg0, arg1)\n";
        }
        case "merge" -> {
          if (arity == 3) return "\treturn this_.emCore().Merge(arg0, arg1, arg2)\n";
        }
        case "forEach" -> {
          if (arity == 1) return "\tthis_.emCore().ForEachBiConsumer(arg0)\n";
        }
        case "keySet" -> {
          if (arity == 0) return "\treturn newLiveEnumMapKeySet(this_.emCore())\n";
        }
        case "values" -> {
          if (arity == 0) return "\treturn newLiveEnumMapValues(this_.emCore())\n";
        }
        case "entrySet" -> {
          if (arity == 0) return "\treturn newLiveEnumMapEntrySet(this_.emCore())\n";
        }
        case "clone" -> {
          if (arity == 0) {
            String st = rb.goStructType(te, imports);
            return "\tout := &"
                + st
                + "{}\n\tout.BindExtSelf_(out)\n\tout.emCore().initCopyCore(this_.emCore())\n\treturn out\n";
          }
        }
        case "equals" -> {
          if (arity == 1) return "\treturn this_.emCore().Equals(arg0)\n";
        }
        case "hashCode" -> {
          if (arity == 0) return "\treturn this_.emCore().HashCode()\n";
        }
        case "toString" -> {
          if (arity == 0) return "\treturn this_.emCore().String()\n";
        }
        default -> {}
      }
    }

    // 3. Concrete Map types backed by *JMapCore
    if (CONCRETE_MAP_TYPES.contains(qn)) {
      boolean isImm = qn.startsWith("com.google.");
      switch (m) {
        case "size" -> {
          if (arity == 0) return "\treturn this_.mc().Size()\n";
        }
        case "isEmpty" -> {
          if (arity == 0) return "\treturn this_.mc().IsEmpty()\n";
        }
        case "containsKey" -> {
          if (arity == 1) return "\treturn this_.mc().ContainsKey(arg0)\n";
        }
        case "containsValue" -> {
          if (arity == 1) return "\treturn this_.mc().ContainsValue(arg0)\n";
        }
        case "get" -> {
          if (arity == 1) return "\treturn this_.mc().Get(arg0)\n";
        }
        case "getOrDefault" -> {
          if (arity == 2) return "\treturn this_.mc().GetOrDefault(arg0, arg1)\n";
        }
        case "put", "putIfAbsent", "remove", "putAll", "clear", "replace", "replaceAll", "computeIfAbsent", "computeIfPresent", "compute", "merge" -> {
          if (isImm) {
            return "\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n";
          }
          switch (m) {
            case "put" -> {
              if (arity == 2) return "\treturn this_.mc().Put(arg0, arg1)\n";
            }
            case "putIfAbsent" -> {
              if (arity == 2) return "\treturn this_.mc().PutIfAbsent(arg0, arg1)\n";
            }
            case "remove" -> {
              if (arity == 1) return "\treturn this_.mc().Remove(arg0)\n";
              if (arity == 2) return "\treturn mapCoreRemove2(this_.mc(), arg0, arg1)\n";
            }
            case "putAll" -> {
              if (arity == 1) return "\tcopyExtMapInto(this_.mc(), arg0)\n";
            }
            case "clear" -> {
              if (arity == 0) return "\tthis_.mc().Clear()\n";
            }
            case "replace" -> {
              if (arity == 2) return "\treturn mapCoreReplace2(this_.mc(), arg0, arg1)\n";
              if (arity == 3) return "\treturn mapCoreReplace3(this_.mc(), arg0, arg1, arg2)\n";
            }
            case "replaceAll" -> {
              if (arity == 1) return "\tmapCoreReplaceAll(this_.mc(), arg0)\n";
            }
            case "computeIfAbsent" -> {
              if (arity == 2) return "\treturn mapCoreComputeIfAbsent(this_.mc(), arg0, arg1)\n";
            }
            case "computeIfPresent" -> {
              if (arity == 2) return "\treturn mapCoreComputeIfPresent(this_.mc(), arg0, arg1)\n";
            }
            case "compute" -> {
              if (arity == 2) return "\treturn mapCoreCompute(this_.mc(), arg0, arg1)\n";
            }
            case "merge" -> {
              if (arity == 3) return "\treturn mapCoreMerge(this_.mc(), arg0, arg1, arg2)\n";
            }
          }
        }
        case "forEach" -> {
          if (arity == 1) return "\tmapCoreForEach(this_.mc(), arg0)\n";
        }
        case "firstKey" -> {
          if (arity == 0) return "\treturn this_.mc().FirstKey()\n";
        }
        case "lastKey" -> {
          if (arity == 0) return "\treturn this_.mc().LastKey()\n";
        }
        case "lowerKey" -> {
          if (arity == 1) return "\treturn sortedCoreNavKey(this_.mc(), arg0, -2)\n";
        }
        case "floorKey" -> {
          if (arity == 1) return "\treturn sortedCoreNavKey(this_.mc(), arg0, -1)\n";
        }
        case "ceilingKey" -> {
          if (arity == 1) return "\treturn sortedCoreNavKey(this_.mc(), arg0, 1)\n";
        }
        case "higherKey" -> {
          if (arity == 1) return "\treturn sortedCoreNavKey(this_.mc(), arg0, 2)\n";
        }
        case "firstEntry" -> {
          if (arity == 0) return "\treturn mapCoreFirstEntry(this_.mc())\n";
        }
        case "lastEntry" -> {
          if (arity == 0) return "\treturn mapCoreLastEntry(this_.mc())\n";
        }
        case "lowerEntry" -> {
          if (arity == 1) return "\treturn sortedCoreNavEntry(this_.mc(), arg0, -2)\n";
        }
        case "floorEntry" -> {
          if (arity == 1) return "\treturn sortedCoreNavEntry(this_.mc(), arg0, -1)\n";
        }
        case "ceilingEntry" -> {
          if (arity == 1) return "\treturn sortedCoreNavEntry(this_.mc(), arg0, 1)\n";
        }
        case "higherEntry" -> {
          if (arity == 1) return "\treturn sortedCoreNavEntry(this_.mc(), arg0, 2)\n";
        }
        case "pollFirstEntry" -> {
          if (arity == 0) {
            return isImm
                ? "\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n"
                : "\treturn mapCorePollFirstEntry(this_.mc())\n";
          }
        }
        case "pollLastEntry" -> {
          if (arity == 0) {
            return isImm
                ? "\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n"
                : "\treturn mapCorePollLastEntry(this_.mc())\n";
          }
        }
        case "putFirst" -> {
          if (arity == 2) {
            return isImm
                ? "\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n"
                : "\told := this_.mc().Get(arg0)\n\tthis_.mc().AddFirstKey(arg0, arg1)\n\treturn old\n";
          }
        }
        case "putLast" -> {
          if (arity == 2) {
            return isImm
                ? "\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n"
                : "\told := this_.mc().Get(arg0)\n\tthis_.mc().AddLastKey(arg0, arg1)\n\treturn old\n";
          }
        }
        case "keySet" -> {
          if (arity == 0) {
            if (isImm) {
              return qn.equals("com.google.common.collect.ImmutableSortedMap")
                  ? "\treturn this_.cachedSortedKeySet_()\n"
                  : "\treturn this_.cachedKeySet_()\n";
            }
            return "\treturn newLiveMapKeySet(this_.mc())\n";
          }
        }
        case "values" -> {
          if (arity == 0) {
            return isImm
                ? "\treturn this_.cachedValues_()\n"
                : "\treturn newLiveMapValues(this_.mc())\n";
          }
        }
        case "entrySet" -> {
          if (arity == 0) {
            return isImm
                ? "\treturn this_.cachedEntrySet_()\n"
                : "\treturn newLiveMapEntrySet(this_.mc())\n";
          }
        }
        case "clone" -> {
          if (arity == 0 && !isImm) {
            String st = rb.goStructType(te, imports);
            if (qn.equals("java.util.TreeMap")) {
              return "\tout := &"
                  + st
                  + "{}\n\tout.BindExtSelf_(out)\n\tcopySortedCoreInto(&out.inlineCore, this_.mc())\n\treturn out\n";
            }
            return "\tout := &"
                + st
                + "{}\n\tout.BindExtSelf_(out)\n\tcopyHashCoreInto(&out.inlineCore, this_.mc())\n\treturn out\n";
          }
        }
        case "equals" -> {
          if (arity == 1) return "\treturn mapCoreEquals(this_.mc(), arg0)\n";
        }
        case "hashCode" -> {
          if (arity == 0) return "\treturn mapCoreHash(this_.mc())\n";
        }
        case "toString" -> {
          if (arity == 0) return "\treturn mapCoreString(this_.mc())\n";
        }
        default -> {}
      }
    }

    // 4. ConcurrentHashMap (mutex-protected *jSyncMapCore)
    if (qn.equals("java.util.concurrent.ConcurrentHashMap")) {
      String call =
          "syncMapOp_"
              + m
              + "_"
              + arity
              + "(this_"
              + (arity > 0 ? ", " + joinArgs(arity) : "")
              + ")";
      return retType.isEmpty() ? "\t" + call + "\n" : "\treturn " + call + "\n";
    }

    // 5. BiMap (HashBiMap & ImmutableBiMap)
    if (qn.equals("com.google.common.collect.HashBiMap")
        || qn.equals("com.google.common.collect.ImmutableBiMap")) {
      String call = "biMapOp_" + m + "_" + arity + "(this_.Core_.(*jBiMapCore)" + (arity > 0 ? ", " + joinArgs(arity) : "") + ")";
      return retType.isEmpty() ? "\t" + call + "\n" : "\treturn " + call + "\n";
    }

    // 6. Multimap family
    if (MULTIMAP_TYPES.contains(qn)) {
      String selfArg =
          (qn.equals("com.google.common.collect.ImmutableListMultimap")
                  || qn.equals("com.google.common.collect.ImmutableSetMultimap")
                  || qn.equals("com.google.common.collect.ImmutableMultimap"))
              ? (qn.contains("Set") ? "true" : "false")
              : "false";
      if (m.equals("inverse") && arity == 0) {
        return "\treturn multimapInverse(this_.getCore().(*jMultimapCore), " + selfArg + ")\n";
      }
      String call =
          "multimapOp_"
              + m
              + "_"
              + arity
              + "(this_.getCore().(*jMultimapCore)"
              + (arity > 0 ? ", " + joinArgs(arity) : "")
              + ")";
      return retType.isEmpty() ? "\t" + call + "\n" : "\treturn " + call + "\n";
    }

    // 7. Multiset family
    if (MULTISET_TYPES.contains(qn)) {
      String call =
          "multisetOp_"
              + m
              + "_"
              + arity
              + "(this_.getCore().(*jMultisetCore)"
              + (arity > 0 ? ", " + joinArgs(arity) : "")
              + ")";
      return retType.isEmpty() ? "\t" + call + "\n" : "\treturn " + call + "\n";
    }

    // 8. Table family
    if (TABLE_TYPES.contains(qn)) {
      if (qn.equals("com.google.common.collect.ImmutableTable")
          && (m.equals("put") || m.equals("clear") || m.equals("remove"))) {
        return "\tpanic(javalang.NewUnsupportedOperationException(\"ImmutableTable\"))\n";
      }
      String call =
          "tableOp_"
              + m
              + "_"
              + arity
              + "(this_.getCore().(*jTableCore)"
              + (arity > 0 ? ", " + joinArgs(arity) : "")
              + ")";
      return retType.isEmpty() ? "\t" + call + "\n" : "\treturn " + call + "\n";
    }

    // 9. Builders (ImmutableCollection/List/Set/SortedSet/Map/SortedMap/Multimap/Multiset & MultimapBuilder)
    switch (qn) {
      case "com.google.common.collect.ImmutableList.Builder" -> {
        if (m.equals("add") && arity == 1) {
          return p0.startsWith("[]")
              ? "\tfor _, x := range CheckNotNullSlice(arg0) { this_.AddVoid_(x) }\n\treturn this_\n"
              : "\tthis_.AddVoid_(arg0)\n\treturn this_\n";
        }
        if (m.equals("addAll") && arity == 1) {
          return "\tthis_.AddAllVoid_(arg0)\n\treturn this_\n";
        }
        if (m.equals("build") && arity == 0) {
          return "\treturn this_.BuildImmList_()\n";
        }
      }
      case "com.google.common.collect.ImmutableCollection.Builder" -> {
        if (m.equals("add") && arity == 1) {
          return p0.startsWith("[]")
              ? "\tfor _, x := range CheckNotNullSlice(arg0) { this_.lc().Add(CheckNotNull(ToAny(x))) }\n\treturn this_\n"
              : "\tthis_.lc().Add(CheckNotNull(ToAny(arg0)))\n\treturn this_\n";
        }
        if (m.equals("addAll") && arity == 1) {
          return "\tlistCoreAddAllNonNull(this_.lc(), arg0)\n\treturn this_\n";
        }
        if (m.equals("build") && arity == 0) {
          return "\treturn wrapImmutableList(this_.lc().Items)\n";
        }
      }
      case "com.google.common.collect.ImmutableSet.Builder" -> {
        if (m.equals("add") && arity == 1) {
          return p0.startsWith("[]")
              ? "\tfor _, x := range CheckNotNullSlice(arg0) { this_.sc().Add(CheckNotNull(ToAny(x))) }\n\treturn this_\n"
              : "\tthis_.sc().Add(CheckNotNull(ToAny(arg0)))\n\treturn this_\n";
        }
        if (m.equals("addAll") && arity == 1) {
          return "\tsetCoreAddAllNonNull(this_.sc(), arg0)\n\treturn this_\n";
        }
        if (m.equals("build") && arity == 0) {
          return "\treturn immSetFromSetCoreCopy(this_.sc())\n";
        }
      }
      case "com.google.common.collect.ImmutableSortedSet.Builder" -> {
        if (m.equals("add") && arity == 1) {
          return p0.startsWith("[]")
              ? "\tfor _, x := range CheckNotNullSlice(arg0) { this_.sortedCore_.Put(CheckNotNull(ToAny(x)), true) }\n\treturn this_\n"
              : "\tthis_.sortedCore_.Put(CheckNotNull(ToAny(arg0)), true)\n\treturn this_\n";
        }
        if (m.equals("addAll") && arity == 1) {
          return "\tmapCoreAddAllNonNull(this_.sortedCore_, arg0)\n\treturn this_\n";
        }
        if (m.equals("build") && arity == 0) {
          return "\treturn immSortedSetFromCoreCopy(this_.sortedCore_)\n";
        }
      }
      case "com.google.common.collect.ImmutableBiMap.Builder" -> {
        if (m.equals("put") && arity == 2) {
          return "\timmMapBuilderPut(this_.imbCore(), arg0, arg1)\n\treturn this_\n";
        }
        if (m.equals("putAll") && arity == 1) {
          return "\timmMapBuilderPutAll(this_.imbCore(), arg0)\n\treturn this_\n";
        }
        if ((m.equals("build") || m.equals("buildOrThrow")) && arity == 0) {
          return "\treturn immBiMapBuilderBuild(this_.imbCore(), true)\n";
        }
        if (m.equals("buildKeepingLast") && arity == 0) {
          return "\treturn immBiMapBuilderBuild(this_.imbCore(), false)\n";
        }
      }
      case "com.google.common.collect.ImmutableMap.Builder",
          "com.google.common.collect.ImmutableSortedMap.Builder" -> {
        if (m.equals("put") && arity == 2) {
          return "\timmMapBuilderPut(this_.imbCore(), arg0, arg1)\n\treturn this_\n";
        }
        if (m.equals("putAll") && arity == 1) {
          return "\timmMapBuilderPutAll(this_.imbCore(), arg0)\n\treturn this_\n";
        }
        if ((m.equals("build") || m.equals("buildOrThrow")) && arity == 0) {
          return "\treturn immMapBuilderBuild(this_.imbCore(), true)\n";
        }
        if (m.equals("buildKeepingLast") && arity == 0) {
          return "\treturn immMapBuilderBuild(this_.imbCore(), false)\n";
        }
      }
      case "com.google.common.collect.ImmutableMultimap.Builder",
          "com.google.common.collect.ImmutableListMultimap.Builder",
          "com.google.common.collect.ImmutableSetMultimap.Builder" -> {
        if (m.equals("put") && arity == 2) {
          return "\tmultimapBuilderPut(this_.getCore().(*jMultimapBuilderCore), arg0, arg1)\n\treturn this_\n";
        }
        if (m.equals("putAll") && arity == 2) {
          return "\tmultimapBuilderPutAll(this_.getCore().(*jMultimapBuilderCore), arg0, arg1)\n\treturn this_\n";
        }
        if (m.equals("orderKeysBy") && arity == 1) {
          return "\tthis_.getCore().(*jMultimapBuilderCore).keyCmp = comparatorFunc(CheckNotNull(arg0))\n\treturn this_\n";
        }
        if (m.equals("build") && arity == 0) {
          return "\treturn multimapBuilderBuild(this_.getCore().(*jMultimapBuilderCore))\n";
        }
      }
      case "com.google.common.collect.ImmutableMultiset.Builder" -> {
        if (m.equals("add") && arity == 1) {
          return p0.startsWith("[]")
              ? "\tc := this_.Core_.(*jMultisetCore)\n\tfor _, x := range arg0 { multisetOp_add_2(c, CheckNotNull(ToAny(x)), 1) }\n\treturn this_\n"
              : "\tmultisetOp_add_2(this_.Core_.(*jMultisetCore), CheckNotNull(ToAny(arg0)), 1)\n\treturn this_\n";
        }
        if (m.equals("addAll") && arity == 1) {
          return "\tc := this_.Core_.(*jMultisetCore)\n\tfor it := Iterate(arg0); it.HasNext(); { multisetOp_add_2(c, CheckNotNull(it.Next()), 1) }\n\treturn this_\n";
        }
        if (m.equals("build") && arity == 0) {
          return "\treturn wrapImmutableMultiset(cloneMultisetCore(this_.Core_.(*jMultisetCore), true))\n";
        }
      }
      case "com.google.common.collect.MultimapBuilder.MultimapBuilderWithKeys" -> {
        if (m.equals("arrayListValues") && arity == 0) {
          return "\tb := &Ext_com_google_common_collect_MultimapBuilder_ListMultimapBuilderBase{}\n\tb.BindExtSelf_(b)\n\tb.Core_ = int32(0)\n\treturn b\n";
        }
        if ((m.equals("hashSetValues") || m.equals("linkedHashSetValues")) && arity == 0) {
          return "\tb := &Ext_com_google_common_collect_MultimapBuilder_SetMultimapBuilderBase{}\n\tb.BindExtSelf_(b)\n\tb.Core_ = int32(1)\n\treturn b\n";
        }
      }
      case "com.google.common.collect.MultimapBuilder",
          "com.google.common.collect.MultimapBuilder.ListMultimapBuilder",
          "com.google.common.collect.MultimapBuilder.SetMultimapBuilder" -> {
        if (m.equals("build") && arity == 0) {
          boolean isSet = qn.endsWith("SetMultimapBuilder");
          return isSet
              ? "\tm := &Ext_com_google_common_collect_LinkedHashMultimap{}\n\tm.Core_ = newMultimapCore(true, false, nil, nil)\n\treturn m\n"
              : "\tm := &Ext_com_google_common_collect_ArrayListMultimap{}\n\tm.Core_ = newMultimapCore(false, false, nil, nil)\n\treturn m\n";
        }
      }
      default -> {}
    }

    // 10. Abstract delegation bases (AbstractCollection, ImmutableCollection, AbstractSet, Sets.SetView, AbstractList, AbstractQueue, AbstractMap)
    if (qn.equals("java.util.AbstractCollection")
        || qn.equals("com.google.common.collect.ImmutableCollection")
        || qn.equals("java.util.AbstractSet")
        || qn.equals("com.google.common.collect.Sets.SetView")) {
      boolean isSet = qn.endsWith("Set") || qn.endsWith("SetView");
      switch (m) {
        case "immutableCopy" -> {
          if (arity == 0) return "\treturn wrapImmutableSet(colToSlice(this_.Self_))\n";
        }
        case "asList" -> {
          if (arity == 0) return "\treturn wrapImmutableList(colToSlice(this_.Self_))\n";
        }
        case "iterator" -> {
          if (arity == 0) {
            if (qn.equals("com.google.common.collect.Sets.SetView")) {
              return "\treturn newSliceObjIter(this_.Core_.(*jSetViewCore).compute())\n";
            }
            return "\tpanic(javalang.NewUnsupportedOperationException(\"AbstractCollection.iterator\"))\n";
          }
        }
        case "size" -> {
          if (arity == 0) {
            if (qn.equals("com.google.common.collect.Sets.SetView")) {
              return "\treturn int32(len(this_.Core_.(*jSetViewCore).compute()))\n";
            }
            return "\tpanic(javalang.NewUnsupportedOperationException(\"AbstractCollection.size\"))\n";
          }
        }
        case "isEmpty" -> {
          if (arity == 0) return "\treturn this_.Self_.Size() == 0\n";
        }
        case "contains" -> {
          if (arity == 1) {
            if (qn.equals("com.google.common.collect.Sets.SetView")) {
              return "\treturn this_.Core_.(*jSetViewCore).contains(arg0)\n";
            }
            return "\treturn abstractColContains(this_.Self_, arg0)\n";
          }
        }
        case "toArray" -> {
          if (arity == 0) return "\treturn colToSlice(this_.Self_)\n";
          if (arity == 1) return "\treturn sliceToArray1(colToSlice(this_.Self_), arg0)\n";
        }
        case "add" -> {
          if (arity == 1) return "\tpanic(javalang.NewUnsupportedOperationException(\"add\"))\n";
        }
        case "remove" -> {
          if (arity == 1) return "\treturn abstractColRemove(this_.Self_, arg0)\n";
        }
        case "containsAll" -> {
          if (arity == 1) return "\treturn colContainsAll(this_.Self_.Contains, arg0)\n";
        }
        case "addAll" -> {
          if (arity == 1) return "\treturn abstractColAddAll(this_.Self_, arg0)\n";
        }
        case "removeAll" -> {
          if (arity == 1) return "\treturn abstractColRemoveAll(this_.Self_, arg0)\n";
        }
        case "retainAll" -> {
          if (arity == 1) return "\treturn abstractColRetainAll(this_.Self_, arg0)\n";
        }
        case "removeIf" -> {
          if (arity == 1) return "\treturn abstractColRemoveIf(this_.Self_, arg0)\n";
        }
        case "clear" -> {
          if (arity == 0) return "\tabstractColClear(this_.Self_)\n";
        }
        case "forEach" -> {
          if (arity == 1) return "\tExtDefault_Iterable_ForEach(this_.Self_, arg0)\n";
        }
        case "stream" -> {
          if (arity == 0) return "\treturn streamFromIterable(this_.Self_)\n";
        }
        case "toString" -> {
          if (arity == 0) return "\treturn sliceJavaString(colToSlice(this_.Self_))\n";
        }
        case "equals" -> {
          if (arity == 1) {
            return isSet
                ? "\treturn abstractSetEquals(this_.Self_, arg0)\n"
                : "\treturn RefEquals(this_.Self_, arg0)\n";
          }
        }
        case "hashCode" -> {
          if (arity == 0) {
            return isSet
                ? "\treturn abstractSetHash(this_.Self_)\n"
                : "\treturn javalang.IdentityHashCode(this_.Self_)\n";
          }
        }
        default -> {}
      }
    }

    if (qn.equals("java.util.AbstractList") || qn.equals("java.util.AbstractSequentialList")) {
      if (m.equals("remove") && arity == 1 && !p0.equals("int32")) {
        return "\treturn abstractColRemove(this_.Self_, arg0)\n";
      }
      String call =
          "abstractListOp_"
              + m
              + "_"
              + arity
              + "(this_.selfList()"
              + (arity > 0 ? ", " + joinArgs(arity) : "")
              + ")";
      return retType.isEmpty() ? "\t" + call + "\n" : "\treturn " + call + "\n";
    }

    if (qn.equals("java.util.AbstractQueue")) {
      String call =
          "abstractQueueOp_"
              + m
              + "_"
              + arity
              + "(this_.selfQueue()"
              + (arity > 0 ? ", " + joinArgs(arity) : "")
              + ")";
      return retType.isEmpty() ? "\t" + call + "\n" : "\treturn " + call + "\n";
    }

    if (qn.equals("java.util.AbstractMap")
        || qn.equals("com.google.common.collect.Maps.IteratorBasedAbstractMap")) {
      String call =
          "abstractMapOp_"
              + m
              + "_"
              + arity
              + "(this_.Self_"
              + (arity > 0 ? ", " + joinArgs(arity) : "")
              + ")";
      return retType.isEmpty() ? "\t" + call + "\n" : "\treturn " + call + "\n";
    }

    // 11. Other collection/util classes: SimpleImmutableEntry, PriorityQueue, BitSet, AbstractIterator, UnmodifiableIterator, Ordering, ComparisonChain, Range, RangeSet, Cache, ThreadPoolExecutor, AtomicIntegerArray
    switch (qn) {
      case "java.util.AbstractMap.SimpleImmutableEntry" -> {
        if (m.equals("getKey") && arity == 0) return "\treturn this_.Core_.(*jSimpleEntry).k\n";
        if (m.equals("getValue") && arity == 0) return "\treturn this_.Core_.(*jSimpleEntry).v\n";
        if (m.equals("setValue") && arity == 1) {
          return "\tpanic(javalang.NewUnsupportedOperationException(\"SimpleImmutableEntry.setValue\"))\n";
        }
        if (m.equals("equals") && arity == 1) {
          return "\treturn mapEntryEquals(this_.Self_, arg0)\n";
        }
        if (m.equals("hashCode") && arity == 0) {
          return "\treturn mapEntryHash(this_.Self_)\n";
        }
        if (m.equals("toString") && arity == 0) {
          return "\te := this_.Core_.(*jSimpleEntry)\n\treturn StrVal(e.k) + \"=\" + StrVal(e.v)\n";
        }
      }
      case "java.util.PriorityQueue" -> {
        String call =
            "pqOp_"
                + m
                + "_"
                + arity
                + "(this_.Core_.(*jPriorityQueue)"
                + (arity > 0 ? ", " + joinArgs(arity) : "")
                + ")";
        return retType.isEmpty() ? "\t" + call + "\n" : "\treturn " + call + "\n";
      }
      case "java.util.BitSet" -> {
        String call =
            "bitSetOp_"
                + sig.goName
                + "(this_"
                + (arity > 0 ? ", " + joinArgs(arity) : "")
                + ")";
        return retType.isEmpty() ? "\t" + call + "\n" : "\treturn " + call + "\n";
      }
      case "com.google.common.collect.AbstractIterator" -> {
        if (m.equals("endOfData") && arity == 0) {
          return "\tthis_.Core_.(*jAbstractIterState).state = 2\n\treturn nil\n";
        }
        if (m.equals("hasNext") && arity == 0) {
          return "\treturn abstractIterHasNext(this_)\n";
        }
        if (m.equals("next") && arity == 0) {
          return "\treturn abstractIterNext(this_)\n";
        }
        if (m.equals("remove") && arity == 0) {
          return "\tpanic(javalang.NewUnsupportedOperationException(\"remove\"))\n";
        }
      }
      case "com.google.common.collect.UnmodifiableIterator",
          "com.google.common.collect.UnmodifiableListIterator" -> {
        if (m.equals("remove") && arity == 0) {
          return "\tpanic(javalang.NewUnsupportedOperationException(\"remove\"))\n";
        }
      }
      case "com.google.common.collect.Ordering" -> {
        if (m.equals("reverse") && arity == 0) {
          return "\tself := this_.Self_\n\treturn newOrdering(func(a, b any) int32 { return self.Compare(b, a) })\n";
        }
        if (m.equals("nullsFirst") && arity == 0) {
          return "\tself := this_.Self_\n\treturn newOrdering(func(a, b any) int32 {\n\t\ta, b = ToAny(a), ToAny(b)\n\t\tif a == nil && b == nil { return 0 }\n\t\tif a == nil { return -1 }\n\t\tif b == nil { return 1 }\n\t\treturn self.Compare(a, b)\n\t})\n";
        }
        if (m.equals("compare") && arity == 2) {
          return "\tif this_.Core_ == (jNatOrderCore{}) { return CompareValues(arg0, arg1) }\n\tif fn, ok := this_.Core_.(func(any, any) int32); ok && fn != nil { return fn(arg0, arg1) }\n\tpanic(\"UNIMPLEMENTED_EXT: Ordering.compare\")\n";
        }
        if (m.equals("sortedCopy") && arity == 1) {
          return "\treturn immListSortedCopyOf(this_.Self_, arg0)\n";
        }
        if (m.equals("reversed") && arity == 0) {
          return "\treturn ExtDefault_Comparator_Reversed(this_.Self_)\n";
        }
        if (m.equals("thenComparing") && arity == 1) {
          return p0.contains("Comparator")
              ? "\treturn ExtDefault_Comparator_ThenComparing(this_.Self_, arg0)\n"
              : "\treturn ExtDefault_Comparator_ThenComparingKey(this_.Self_, arg0)\n";
        }
        if (m.equals("thenComparingInt") && arity == 1) {
          return "\treturn ExtDefault_Comparator_ThenComparingInt(this_.Self_, arg0)\n";
        }
        if (m.equals("thenComparingLong") && arity == 1) {
          return "\treturn ExtDefault_Comparator_ThenComparingLong(this_.Self_, arg0)\n";
        }
      }
      case "com.google.common.collect.ComparisonChain" -> {
        if (m.equals("compare") && arity == 2) {
          return p0.equals("int32")
              ? "\treturn this_.Core_.(*jComparisonChain).cmpInt(arg0, arg1)\n"
              : "\treturn this_.Core_.(*jComparisonChain).cmpAny(arg0, arg1)\n";
        }
        if (m.equals("result") && arity == 0) {
          return "\treturn this_.Core_.(*jComparisonChain).res\n";
        }
      }
      case "com.google.common.collect.Range" -> {
        if ((m.equals("apply") || m.equals("test") || m.equals("contains")) && arity == 1) {
          return "\trc := this_.Core_.(*jRangeCore)\n\tv := CheckNotNull(ToAny(arg0))\n\treturn CompareValues(v, rc.lower) >= 0 && CompareValues(v, rc.upper) <= 0\n";
        }
        if (m.equals("equals") && arity == 1) {
          return "\to, ok := arg0.(*Ext_com_google_common_collect_Range)\n\tif !ok || o == nil { return false }\n\tr1, r2 := this_.Core_.(*jRangeCore), o.Core_.(*jRangeCore)\n\treturn javalang.ObjectsEquals(r1.lower, r2.lower) && javalang.ObjectsEquals(r1.upper, r2.upper)\n";
        }
      }
      case "com.google.common.collect.TreeRangeSet",
          "com.google.common.collect.AbstractRangeSet",
          "com.google.common.collect.ImmutableRangeSet" -> {
        if (m.equals("add") && arity == 1) {
          if (qn.equals("com.google.common.collect.ImmutableRangeSet")) {
            return "\tpanic(javalang.NewUnsupportedOperationException(\"ImmutableRangeSet.add\"))\n";
          }
          return "\tCheckNotNullPtr(arg0)\n\trc := arg0.Core_.(*jRangeCore)\n\tif CompareValues(rc.lower, rc.upper) <= 0 {\n\t\trs := this_.getCore().(*jRangeSetCore)\n\t\trs.ranges = append(rs.ranges, rc)\n\t}\n";
        }
        if (m.equals("contains") && arity == 1) {
          return "\tv := CheckNotNull(ToAny(arg0))\n\tfor _, r := range this_.getCore().(*jRangeSetCore).ranges {\n\t\tif CompareValues(v, r.lower) >= 0 && CompareValues(v, r.upper) <= 0 { return true }\n\t}\n\treturn false\n";
        }
        if (m.equals("isEmpty") && arity == 0) {
          return "\treturn len(this_.getCore().(*jRangeSetCore).ranges) == 0\n";
        }
      }
      case "com.google.common.cache.CacheBuilder" -> {
        if (m.equals("maximumSize") && arity == 1) {
          return "\tthis_.Core_.(*jCacheBuilderCore).maxSize = arg0\n\treturn this_\n";
        }
        if (m.equals("removalListener") && arity == 1) {
          return "\tthis_.Core_.(*jCacheBuilderCore).listener = CheckNotNull(arg0)\n\treturn this_\n";
        }
        if (m.equals("build") && arity == 1) {
          return "\treturn newLoadingCache(this_.Core_.(*jCacheBuilderCore), CheckNotNull(arg0))\n";
        }
      }
      case "com.google.common.cache.RemovalNotification" -> {
        if (m.equals("getKey") && arity == 0) return "\treturn this_.Core_.(*jSimpleEntry).k\n";
        if (m.equals("getValue") && arity == 0) return "\treturn this_.Core_.(*jSimpleEntry).v\n";
        if (m.equals("setValue") && arity == 1) {
          return "\tpanic(javalang.NewUnsupportedOperationException(\"RemovalNotification.setValue\"))\n";
        }
      }
      case "java.util.concurrent.ThreadPoolExecutor",
          "java.util.concurrent.AbstractExecutorService" -> {
        if (m.equals("shutdown") && arity == 0) {
          return "\tif ex, ok := this_.Core_.(Ext_java_util_concurrent_ExecutorService); ok { ex.Shutdown() }\n";
        }
        if (m.equals("close") && arity == 0) {
          return "\tif ex, ok := this_.Core_.(Ext_java_util_concurrent_ExecutorService); ok { ex.Close_0() }\n";
        }
        if (m.equals("submit") && arity == 1) {
          return "\tif ex, ok := this_.Core_.(Ext_java_util_concurrent_ExecutorService); ok { return ex.Submit_1_Ext_java_util_concurrent_Callable(arg0) }\n\treturn execSubmitCallable(arg0)\n";
        }
      }
      case "java.util.concurrent.atomic.AtomicIntegerArray" -> {
        if (m.equals("get") && arity == 1) {
          return "\treturn this_.Core_.(*jAtomicIntArray).get(arg0)\n";
        }
        if (m.equals("updateAndGet") && arity == 2) {
          return "\treturn this_.Core_.(*jAtomicIntArray).updateAndGet(arg0, arg1)\n";
        }
      }
      case "java.util.Random" -> {
        if (m.equals("nextInt") && arity == 0) return "\treturn this_.Core_.(*javautil.Random).NextInt()\n";
        if (m.equals("nextInt") && arity == 1) return "\treturn this_.Core_.(*javautil.Random).NextIntBound(arg0)\n";
        if (m.equals("nextLong") && arity == 0) return "\treturn this_.Core_.(*javautil.Random).NextLong()\n";
        if (m.equals("nextBoolean") && arity == 0) return "\treturn this_.Core_.(*javautil.Random).NextBoolean()\n";
        if (m.equals("setSeed") && arity == 1) return "\tthis_.Core_.(*javautil.Random).SetSeed(arg0)\n";
      }
      default -> {}
    }

    return null;
  }

  private static String joinArgs(int arity) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < arity; i++) {
      if (i > 0) sb.append(", ");
      sb.append("arg").append(i);
    }
    return sb.toString();
  }

  /**
   * Emits {@code j2gort/ext_collections_rt.go} containing all deterministic runtime implementations
   * for collections, live views, multimaps, multisets, tables, streams, collectors, caches, and
   * functional default methods.
   */
  public static String buildCollectionsRuntimeGo(RuntimeBindings rb) {
    return buildColRtPart1()
        + buildColRtPart1a()
        + buildColRtMapCursors(rb)
        + buildColRtPart1b()
        + buildColRtPart1c()
        + buildColRtPart2()
        + buildColRtPart3()
        + buildColRtLoweringHelpers()
        + buildColRtOrderedMapHelpers()
        + buildColRtPlainMapSetHelpers()
        + buildColRtCustomOrderedMapHelpers()
        + buildColRtIdentityAndCustomSetHelpers()
        + buildColRtDequeHelpers()
        + buildColRtNodeDequeHelpers()
        + buildColRtSlabAndTypedListHelpers()
        + buildColRtStringBuilderHelpers();
  }

  private static String buildColRtPart1() {
    return """
        // Auto-generated clean-room deterministic collections, Guava, and Stream runtime for j2gort.
        package j2gort

        import (
        \t"fmt"
        \t"math/bits"
        \t"reflect"
        \t"sort"
        \t"strconv"
        \t"strings"
        \t"sync"
        \t"sync/atomic"
        \t"unicode/utf16"

        \t"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
        \t"github.com/paullepoulpe/jsgo/transpiled/shims/javautil"
        \t"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
        )

        var _ = javautil.NewBitSet
        var _ = nodestore.Node(0)
        var _ = bits.TrailingZeros64
        var _ = reflect.TypeFor[int]
        var _ = (*atomic.Bool)(nil)
        var _ = strconv.AppendInt
        var _ = utf16.DecodeRune

        var (
        \tcompPoolResetMu    sync.Mutex
        \tcompPoolResets     []func()
        \truntimePoolPurgers []func()
        )

        func registerCompPoolReset(fn func()) {
        \tcompPoolResetMu.Lock()
        \tcompPoolResets = append(compPoolResets, fn)
        \tcompPoolResetMu.Unlock()
        }

        func registerRuntimePoolPurger(fn func()) {
        \tcompPoolResetMu.Lock()
        \truntimePoolPurgers = append(runtimePoolPurgers, fn)
        \tcompPoolResetMu.Unlock()
        }

        func PurgeAllRuntimePools() {
        \tTrimStaleSlabPools(0)
        \tcompPoolResetMu.Lock()
        \tfns := runtimePoolPurgers
        \tcompPoolResetMu.Unlock()
        \tfor _, fn := range fns {
        \t\tfn()
        \t}
        }

        func FlushCompilationPools() {
        \tcompPoolResetMu.Lock()
        \tfns := compPoolResets
        \tcompPoolResetMu.Unlock()
        \tfor _, fn := range fns {
        \t\tfn()
        \t}
        \tslabGlobalPooledBytes.Store(0)
        }

        func ResetCompilationPools() {
        \tFlushCompilationPools()
        }

        func init() {
        \tnodestore.RegisterStoreResetHook(FlushCompilationPools)
        \tRegisterTypePredicate("com.google.common.collect.Multimap", "com.google.common.collect.Multimap", func(v any) bool {
        \t\t_, ok := v.(Ext_com_google_common_collect_Multimap)
        \t\treturn ok
        \t})
        \tRegisterTypePredicate("java.util.Collection", "java.util.Collection", func(v any) bool {
        \t\t_, ok := v.(Ext_java_util_Collection)
        \t\treturn ok
        \t})
        \tRegisterTypePredicate("java.util.List", "java.util.List", func(v any) bool {
        \t\t_, ok := v.(Ext_java_util_List)
        \t\treturn ok
        \t})
        \tRegisterTypePredicate("java.util.Set", "java.util.Set", func(v any) bool {
        \t\t_, ok := v.(Ext_java_util_Set)
        \t\treturn ok
        \t})
        \tRegisterTypePredicate("java.util.Map", "java.util.Map", func(v any) bool {
        \t\t_, ok := v.(Ext_java_util_Map)
        \t\treturn ok
        \t})
        }

        // --- External Interface Default Methods ---

        func ExtDefault_Iterable_ForEach(this_ Ext_java_lang_Iterable, action Ext_java_util_function_Consumer) {
        \tCheckNotNull(action)
        \tfor it := NewIterCursor(this_); it.HasNext(); {
        \t\taction.Accept(it.Next())
        \t}
        }

        func ExtDefault_Iterator_ForEachRemaining(this_ Ext_java_util_Iterator, action Ext_java_util_function_Consumer) {
        \tCheckNotNull(action)
        \tfor this_.HasNext() {
        \t\taction.Accept(this_.Next())
        \t}
        }

        type jNatOrderCore struct{}

        type jComparatorFunc struct {
        \tfn  func(a, b any) int32
        \tnat bool
        }

        var jNaturalOrderCmp Ext_java_util_Comparator = &jComparatorFunc{fn: CompareValues, nat: true}
        var jReverseOrderCmp Ext_java_util_Comparator = &jComparatorFunc{fn: func(a, b any) int32 { return CompareValues(b, a) }}

        func newComparatorFromFunc(fn func(a, b any) int32) Ext_java_util_Comparator {
        \treturn &jComparatorFunc{fn: fn}
        }

        func (c *jComparatorFunc) isNatOrder_() bool { return c != nil && c.nat }
        func (*jComparatorFunc) Implements_j2gort_Ext_java_util_Comparator() {}
        func (c *jComparatorFunc) Compare(a, b any) int32 { return c.fn(a, b) }
        func (c *jComparatorFunc) Reversed() Ext_java_util_Comparator {
        \treturn ExtDefault_Comparator_Reversed(c)
        }
        func (c *jComparatorFunc) ThenComparing_1_Ext_java_util_Comparator(other Ext_java_util_Comparator) Ext_java_util_Comparator {
        \treturn ExtDefault_Comparator_ThenComparing(c, other)
        }
        func (c *jComparatorFunc) ThenComparing_1_Ext_java_util_function_Function(ke Ext_java_util_function_Function) Ext_java_util_Comparator {
        \treturn ExtDefault_Comparator_ThenComparingKey(c, ke)
        }
        func (c *jComparatorFunc) ThenComparingInt(ke Ext_java_util_function_ToIntFunction) Ext_java_util_Comparator {
        \treturn ExtDefault_Comparator_ThenComparingInt(c, ke)
        }
        func (c *jComparatorFunc) ThenComparingLong(ke Ext_java_util_function_ToLongFunction) Ext_java_util_Comparator {
        \treturn ExtDefault_Comparator_ThenComparingLong(c, ke)
        }

        func isNaturalOrderComparator(cmp Ext_java_util_Comparator) bool {
        \tif isJavaNilAny(cmp) {
        \t\treturn true
        \t}
        \tif m, ok := any(cmp).(interface{ isNatOrder_() bool }); ok && m.isNatOrder_() {
        \t\treturn true
        \t}
        \treturn false
        }

        func initSortedCore(mc *JMapCore, cmp func(a, b any) int32) *JMapCore {
        \tif cmp == nil {
        \t\tmc.mode = jMapModeSortedNat
        \t\tmc.ext = nil
        \t} else {
        \t\tmc.mode = jMapModeSortedCmp
        \t\tmc.ext = &jSortedMapState{cmp: cmp}
        \t}
        \treturn mc
        }

        func initSortedCoreCmp(mc *JMapCore, cmp Ext_java_util_Comparator) *JMapCore {
        \tif isNaturalOrderComparator(cmp) {
        \t\treturn initSortedCore(mc, nil)
        \t}
        \treturn initSortedCore(mc, comparatorFunc(cmp))
        }

        func initSetCapLoad(sc *JSetCore, initialCap int32, loadFactor float32) {
        \tif initialCap < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException("Illegal initial capacity: " + fmt.Sprint(initialCap)))
        \t}
        \tif loadFactor <= 0 || loadFactor != loadFactor {
        \t\tpanic(javalang.NewIllegalArgumentException("Illegal load factor: " + fmt.Sprint(loadFactor)))
        \t}
        \tsc.InitExpectedSize(initialCap)
        }

        func newSetFromMap(m Ext_java_util_Map) Ext_java_util_Set {
        \tCheckNotNull(m)
        \tif m.Size() != 0 {
        \t\tpanic(javalang.NewIllegalArgumentException("Map is non-empty"))
        \t}
        \tobj := &Ext_java_util_HashSetBase{}
        \tobj.BindExtSelf_(obj)
        \tobj.Core_ = m
        \treturn obj
        }

        func newSortedCoreFromCmp(cmp Ext_java_util_Comparator) *JMapCore {
        \tif isNaturalOrderComparator(cmp) {
        \t\treturn NewSortedJMapCore(nil)
        \t}
        \treturn NewSortedJMapCore(comparatorFunc(cmp))
        }

        func sortedCoreNav(mc *JMapCore, k any, op int) (any, any, bool) {
        \tk = ToAny(k)
        \tif mc == nil || mc.Live == 0 || mc.ext == nil {
        \t\tif mc != nil && mc.mode == jMapModeSortedNat && k == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\treturn nil, nil, false
        \t}
        \tst := mc.ext.(*jSortedMapState)
        \tif st.natOrder && k == nil {
        \t\tThrowNPE()
        \t}
        \tbest := int32(0)
        \tcur := st.root
        \tfor cur != 0 {
        \t\tc := st.avlCmpAny(k, cur)
        \t\tif c == 0 {
        \t\t\tif op == -1 || op == 1 {
        \t\t\t\treturn st.treeNodeKey(cur), st.tree[cur-1].val, true
        \t\t\t}
        \t\t\tif op < 0 {
        \t\t\t\tcur = st.tree[cur-1].left
        \t\t\t} else {
        \t\t\t\tcur = st.tree[cur-1].right
        \t\t\t}
        \t\t} else if c < 0 {
        \t\t\tif op > 0 {
        \t\t\t\tbest = cur
        \t\t\t}
        \t\t\tcur = st.tree[cur-1].left
        \t\t} else {
        \t\t\tif op < 0 {
        \t\t\t\tbest = cur
        \t\t\t}
        \t\t\tcur = st.tree[cur-1].right
        \t\t}
        \t}
        \tif best == 0 {
        \t\treturn nil, nil, false
        \t}
        \treturn st.treeNodeKey(best), st.tree[best-1].val, true
        }

        func sortedCoreNavKey(mc *JMapCore, k any, op int) any {
        \trk, _, ok := sortedCoreNav(mc, k, op)
        \tif !ok {
        \t\treturn nil
        \t}
        \treturn rk
        }

        func sortedCoreNavEntry(mc *JMapCore, k any, op int) any {
        \trk, rv, ok := sortedCoreNav(mc, k, op)
        \tif !ok {
        \t\treturn nil
        \t}
        \treturn &jSimpleEntry{k: rk, v: rv, imm: true}
        }

        func mapCoreFirstEntry(mc *JMapCore) any {
        \tif mc == nil || mc.Live == 0 {
        \t\treturn nil
        \t}
        \tk, v := mc.FirstEntry()
        \treturn &jSimpleEntry{k: k, v: v, imm: true}
        }

        func mapCoreLastEntry(mc *JMapCore) any {
        \tif mc == nil || mc.Live == 0 {
        \t\treturn nil
        \t}
        \tk := mc.LastKey()
        \treturn &jSimpleEntry{k: k, v: mc.Get(k), imm: true}
        }

        func mapCorePollFirstEntry(mc *JMapCore) any {
        \tif mc == nil || mc.Live == 0 {
        \t\treturn nil
        \t}
        \tk, v := mc.FirstEntry()
        \tmc.Remove(k)
        \treturn &jSimpleEntry{k: k, v: v, imm: true}
        }

        func mapCorePollLastEntry(mc *JMapCore) any {
        \tif mc == nil || mc.Live == 0 {
        \t\treturn nil
        \t}
        \tk := mc.LastKey()
        \tv := mc.Remove(k)
        \treturn &jSimpleEntry{k: k, v: v, imm: true}
        }

        func comparatorFunc(cmp Ext_java_util_Comparator) func(a, b any) int32 {
        \tif isNaturalOrderComparator(cmp) {
        \t\treturn CompareValues
        \t}
        \treturn func(a, b any) int32 { return cmp.Compare(a, b) }
        }

        func comparatorFuncOrNatural(cmp Ext_java_util_Comparator) func(a, b any) int32 {
        \treturn comparatorFunc(cmp)
        }

        func ExtDefault_Comparator_Reversed(this_ Ext_java_util_Comparator) Ext_java_util_Comparator {
        \treturn newComparatorFromFunc(func(a, b any) int32 { return this_.Compare(b, a) })
        }

        func ExtDefault_Comparator_ThenComparing(this_ Ext_java_util_Comparator, other Ext_java_util_Comparator) Ext_java_util_Comparator {
        \tCheckNotNull(other)
        \treturn newComparatorFromFunc(func(a, b any) int32 {
        \t\tr := this_.Compare(a, b)
        \t\tif r != 0 {
        \t\t\treturn r
        \t\t}
        \t\treturn other.Compare(a, b)
        \t})
        }

        func ExtDefault_Comparator_ThenComparingKeyCmp(this_ Ext_java_util_Comparator, ke Ext_java_util_function_Function, kc Ext_java_util_Comparator) Ext_java_util_Comparator {
        \tCheckNotNull(ke)
        \tCheckNotNull(kc)
        \treturn newComparatorFromFunc(func(a, b any) int32 {
        \t\tr := this_.Compare(a, b)
        \t\tif r != 0 {
        \t\t\treturn r
        \t\t}
        \t\treturn kc.Compare(ke.Apply(a), ke.Apply(b))
        \t})
        }

        func ExtDefault_Comparator_ThenComparingKey(this_ Ext_java_util_Comparator, ke Ext_java_util_function_Function) Ext_java_util_Comparator {
        \tCheckNotNull(ke)
        \treturn newComparatorFromFunc(func(a, b any) int32 {
        \t\tr := this_.Compare(a, b)
        \t\tif r != 0 {
        \t\t\treturn r
        \t\t}
        \t\treturn CompareValues(ke.Apply(a), ke.Apply(b))
        \t})
        }

        func ExtDefault_Comparator_ThenComparingInt(this_ Ext_java_util_Comparator, ke Ext_java_util_function_ToIntFunction) Ext_java_util_Comparator {
        \tCheckNotNull(ke)
        \treturn newComparatorFromFunc(func(a, b any) int32 {
        \t\tr := this_.Compare(a, b)
        \t\tif r != 0 {
        \t\t\treturn r
        \t\t}
        \t\tx, y := ke.ApplyAsInt(a), ke.ApplyAsInt(b)
        \t\tif x < y { return -1 }
        \t\tif x > y { return 1 }
        \t\treturn 0
        \t})
        }

        func ExtDefault_Comparator_ThenComparingLong(this_ Ext_java_util_Comparator, ke Ext_java_util_function_ToLongFunction) Ext_java_util_Comparator {
        \tCheckNotNull(ke)
        \treturn newComparatorFromFunc(func(a, b any) int32 {
        \t\tr := this_.Compare(a, b)
        \t\tif r != 0 {
        \t\t\treturn r
        \t\t}
        \t\tx, y := ke.ApplyAsLong(a), ke.ApplyAsLong(b)
        \t\tif x < y { return -1 }
        \t\tif x > y { return 1 }
        \t\treturn 0
        \t})
        }

        func ExtDefault_Comparator_ThenComparingDouble(this_ Ext_java_util_Comparator, ke any) Ext_java_util_Comparator {
        \tCheckNotNull(ke)
        \tdf := ke.(interface{ ApplyAsDouble(any) float64 })
        \treturn newComparatorFromFunc(func(a, b any) int32 {
        \t\tr := this_.Compare(a, b)
        \t\tif r != 0 {
        \t\t\treturn r
        \t\t}
        \t\treturn int32(javalang.DoubleCompare(df.ApplyAsDouble(a), df.ApplyAsDouble(b)))
        \t})
        }

        type jPredFunc struct {
        \tfn     func(any) bool
        \tnodeFn func(uint32) bool
        }

        func (*jPredFunc) Implements_j2gort_Ext_java_util_function_Predicate() {}
        func (*jPredFunc) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (*jPredFunc) HasIdentityEquals_() {}
        func (p *jPredFunc) Test(x any) bool   { return p.fn(x) }
        func (p *jPredFunc) Apply(x any) bool  { return p.fn(x) }
        func (p *jPredFunc) TestNode_(n uint32) bool {
        \tif p == nil {
        \t\tThrowNPE()
        \t}
        \tif p.nodeFn != nil {
        \t\treturn p.nodeFn(n)
        \t}
        \treturn p.fn(boxNodeById(n))
        }
        func (p *jPredFunc) Equals(o any) bool { return RefEquals(p, o) }

        func ExtDefault_Predicate_Negate(this_ Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \treturn &jPredFunc{
        \t\tfn:     func(x any) bool { return !this_.Test(x) },
        \t\tnodeFn: func(n uint32) bool { return !PredTestNode(this_, n) },
        \t}
        }

        func ExtDefault_Predicate_And(this_ Ext_java_util_function_Predicate, other Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNull(other)
        \treturn &jPredFunc{
        \t\tfn:     func(x any) bool { return this_.Test(x) && other.Test(x) },
        \t\tnodeFn: func(n uint32) bool { return PredTestNode(this_, n) && PredTestNode(other, n) },
        \t}
        }

        func ExtDefault_Predicate_Or(this_ Ext_java_util_function_Predicate, other Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNull(other)
        \treturn &jPredFunc{
        \t\tfn:     func(x any) bool { return this_.Test(x) || other.Test(x) },
        \t\tnodeFn: func(n uint32) bool { return PredTestNode(this_, n) || PredTestNode(other, n) },
        \t}
        }

        type jFuncIdentity struct{}

        func (*jFuncIdentity) Implements_j2gort_Ext_java_util_function_Function() {}
        func (*jFuncIdentity) Implements_j2gort_Ext_com_google_common_base_Function() {}
        func (*jFuncIdentity) Apply(x any) any       { return x }
        func (f *jFuncIdentity) Equals(o any) bool   { _, ok := o.(*jFuncIdentity); return ok }

        type jFuncWrap struct {
        \tfn func(any) any
        }

        func (*jFuncWrap) Implements_j2gort_Ext_java_util_function_Function() {}
        func (*jFuncWrap) Implements_j2gort_Ext_com_google_common_base_Function() {}
        func (f *jFuncWrap) Apply(x any) any   { return f.fn(x) }
        func (f *jFuncWrap) Equals(o any) bool { return RefEquals(f, o) }

        func ExtDefault_Function_AndThen(this_ Ext_java_util_function_Function, after Ext_java_util_function_Function) Ext_java_util_function_Function {
        \tCheckNotNull(after)
        \treturn &jFuncWrap{fn: func(x any) any { return after.Apply(this_.Apply(x)) }}
        }

        func ExtDefault_Function_Compose(this_ Ext_java_util_function_Function, before Ext_java_util_function_Function) Ext_java_util_function_Function {
        \tCheckNotNull(before)
        \treturn &jFuncWrap{fn: func(x any) any { return this_.Apply(before.Apply(x)) }}
        }

        type jConsumerWrap struct {
        \tfn func(any)
        }

        func (*jConsumerWrap) Implements_j2gort_Ext_java_util_function_Consumer() {}
        func (c *jConsumerWrap) Accept(x any) { c.fn(x) }

        func ExtDefault_Consumer_AndThen(this_ Ext_java_util_function_Consumer, after Ext_java_util_function_Consumer) Ext_java_util_function_Consumer {
        \tCheckNotNull(after)
        \treturn &jConsumerWrap{fn: func(x any) { this_.Accept(x); after.Accept(x) }}
        }

        // --- Core Collection Equality, Hash, ToArray, String &Iterators ---

        func sliceToArray1(items []any, dst []any) []any {
        \tCheckNotNullSlice(dst)
        \tn := len(items)
        \tif len(dst) < n {
        \t\tout := make([]any, n)
        \t\tcopy(out, items)
        \t\treturn out
        \t}
        \tcopy(dst, items)
        \tif len(dst) > n {
        \t\tdst[n] = nil
        \t}
        \treturn dst
        }

        type jFastSetCoreProvider interface {
        \tfastSetCore_() *JMapCore
        }

        type jFastJSetCoreProvider interface {
        \tfastSetJCore_() *JSetCore
        }

        type jFastMapCoreProvider interface {
        \tfastMapCore_() *JMapCore
        }

        type jFastEnumSetProvider interface {
        \tfastEnumSet_() *jEnumSetCore
        }

        type jFastEnumMapProvider interface {
        \tfastEnumMap_() *jEnumMapCore
        }

        func immSliceContains(items []any, target any) bool {
        \tif len(items) == 0 {
        \t\treturn false
        \t}
        \tt := ToAny(target)
        \tif t == nil {
        \t\treturn false
        \t}
        \tfor _, x := range items {
        \t\tif elemJavaEquals(x, t, false) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func immSliceIndexOf(items []any, target any) int32 {
        \tif len(items) == 0 {
        \t\treturn -1
        \t}
        \tt := ToAny(target)
        \tif t == nil {
        \t\treturn -1
        \t}
        \tfor i, x := range items {
        \t\tif elemJavaEquals(x, t, false) {
        \t\t\treturn int32(i)
        \t\t}
        \t}
        \treturn -1
        }

        func immSliceLastIndexOf(items []any, target any) int32 {
        \tif len(items) == 0 {
        \t\treturn -1
        \t}
        \tt := ToAny(target)
        \tif t == nil {
        \t\treturn -1
        \t}
        \tfor i := len(items) - 1; i >= 0; i-- {
        \t\tif elemJavaEquals(items[i], t, false) {
        \t\t\treturn int32(i)
        \t\t}
        \t}
        \treturn -1
        }

        func colContainsAll(containsFn func(any) bool, other Ext_java_util_Collection) bool {
        \tCheckNotNull(other)
        \tif tv, ok := any(other).(jTypedViewCursorSource); ok {
        \t\tif n, ok2 := tv.typedViewLen_(); ok2 {
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\tif !containsFn(tv.typedViewAt_(i)) {
        \t\t\t\t\treturn false
        \t\t\t\t}
        \t\t\t}
        \t\t\treturn true
        \t\t}
        \t}
        \tif si, ok := any(other).(SliceIterable); ok {
        \t\tfor _, x := range si.IterSlice_() {
        \t\t\tif !containsFn(x) {
        \t\t\t\treturn false
        \t\t\t}
        \t\t}
        \t\treturn true
        \t}
        \tif jsp, ok := any(other).(jFastJSetCoreProvider); ok {
        \t\tif sc := jsp.fastSetJCore_(); sc != nil {
        \t\t\tall := true
        \t\t\tsc.ForEach(func(k any) bool {
        \t\t\t\tif !containsFn(k) {
        \t\t\t\t\tall = false
        \t\t\t\t\treturn false
        \t\t\t\t}
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn all
        \t\t}
        \t}
        \tif sp, ok := any(other).(jFastSetCoreProvider); ok {
        \t\tif mc := sp.fastSetCore_(); mc != nil {
        \t\t\tall := true
        \t\t\tmc.ForEachEntry(func(k, _ any) bool {
        \t\t\t\tif !containsFn(k) {
        \t\t\t\t\tall = false
        \t\t\t\t\treturn false
        \t\t\t\t}
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn all
        \t\t}
        \t}
        \tif ep, ok := any(other).(jFastEnumSetProvider); ok {
        \t\tif es := ep.fastEnumSet_(); es != nil {
        \t\t\treturn es.ForEach(func(k any) bool { return containsFn(k) })
        \t\t}
        \t}
        \tfor it := NewIterCursor(other); it.HasNext(); {
        \t\tif !containsFn(it.Next()) {
        \t\t\treturn false
        \t\t}
        \t}
        \treturn true
        }

        func listCoreAddAllIterable(lc *JListCore, src any) bool {
        \tCheckNotNull(src)
        \tif tv, ok := src.(jTypedViewCursorSource); ok {
        \t\tif n, ok2 := tv.typedViewLen_(); ok2 {
        \t\t\tif n == 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tlc.InitCap(lc.Size() + int32(n))
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\tlc.Add(tv.typedViewAt_(i))
        \t\t\t}
        \t\t\treturn true
        \t\t}
        \t}
        \tif si, ok := src.(SliceIterable); ok {
        \t\titems := si.IterSlice_()
        \t\tif len(items) == 0 {
        \t\t\treturn false
        \t\t}
        \t\tlc.InitCap(lc.Size() + int32(len(items)))
        \t\tfor _, x := range items {
        \t\t\tlc.Add(x)
        \t\t}
        \t\treturn true
        \t}
        \tif mv, ok := src.(*jLiveMapValues); ok {
        \t\tif mv == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\tif mv.mc.Size() == 0 {
        \t\t\treturn false
        \t\t}
        \t\tlc.InitCap(lc.Size() + mv.mc.Size())
        \t\tmv.mc.ForEachEntry(func(_, v any) bool {
        \t\t\tlc.Add(v)
        \t\t\treturn true
        \t\t})
        \t\treturn true
        \t}
        \tif mk, ok := src.(*jLiveMapKeySet); ok {
        \t\tif mk == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\tif mk.mc.Size() == 0 {
        \t\t\treturn false
        \t\t}
        \t\tlc.InitCap(lc.Size() + mk.mc.Size())
        \t\tmk.mc.ForEachEntry(func(k, _ any) bool {
        \t\t\tlc.Add(k)
        \t\t\treturn true
        \t\t})
        \t\treturn true
        \t}
        \tif jsp, ok := src.(jFastJSetCoreProvider); ok {
        \t\tif sc := jsp.fastSetJCore_(); sc != nil {
        \t\t\tif sc.Size() == 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tlc.InitCap(lc.Size() + sc.Size())
        \t\t\tsc.ForEach(func(k any) bool {
        \t\t\t\tlc.Add(k)
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn true
        \t\t}
        \t}
        \tif sp, ok := src.(jFastSetCoreProvider); ok {
        \t\tif mc := sp.fastSetCore_(); mc != nil {
        \t\t\tif mc.Size() == 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tlc.InitCap(lc.Size() + mc.Size())
        \t\t\tmc.ForEachEntry(func(k, _ any) bool {
        \t\t\t\tlc.Add(k)
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn true
        \t\t}
        \t}
        \tif ep, ok := src.(jFastEnumSetProvider); ok {
        \t\tif es := ep.fastEnumSet_(); es != nil {
        \t\t\tif es.size == 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tlc.InitCap(lc.Size() + es.size)
        \t\t\tes.ForEach(func(k any) bool {
        \t\t\t\tlc.Add(k)
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn true
        \t\t}
        \t}
        \tif c, ok := src.(Ext_java_util_Collection); ok {
        \t\tn := c.Size()
        \t\tif n == 0 {
        \t\t\treturn false
        \t\t}
        \t\tlc.InitCap(lc.Size() + n)
        \t}
        \tadded := false
        \tfor it := NewIterCursor(src); it.HasNext(); {
        \t\tlc.Add(it.Next())
        \t\tadded = true
        \t}
        \treturn added
        }

        func listCoreAddAllNonNull(lc *JListCore, src any) {
        \tCheckNotNull(src)
        \tif ssi, ok := src.(StrSliceIterable); ok {
        \t\tif strs, ok2 := ssi.IterStrSlice_(); ok2 {
        \t\t\tif len(strs) == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tlc.InitCap(lc.Size() + int32(len(strs)))
        \t\t\tfor _, s := range strs {
        \t\t\t\tif s == NullString {
        \t\t\t\t\tThrowNPE()
        \t\t\t\t}
        \t\t\t\tlc.Add(s)
        \t\t\t}
        \t\t\treturn
        \t\t}
        \t}
        \tif tv, ok := src.(jTypedViewCursorSource); ok {
        \t\tif n, ok2 := tv.typedViewLen_(); ok2 {
        \t\t\tif n == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tlc.InitCap(lc.Size() + int32(n))
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\tlc.Add(CheckNotNull(ToAny(tv.typedViewAt_(i))))
        \t\t\t}
        \t\t\treturn
        \t\t}
        \t}
        \tif si, ok := src.(SliceIterable); ok {
        \t\titems := si.IterSlice_()
        \t\tif len(items) == 0 {
        \t\t\treturn
        \t\t}
        \t\tlc.InitCap(lc.Size() + int32(len(items)))
        \t\tfor _, x := range items {
        \t\t\tlc.Add(CheckNotNull(ToAny(x)))
        \t\t}
        \t\treturn
        \t}
        \tif mv, ok := src.(*jLiveMapValues); ok {
        \t\tif mv == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\tif mv.mc.Size() == 0 {
        \t\t\treturn
        \t\t}
        \t\tlc.InitCap(lc.Size() + mv.mc.Size())
        \t\tmv.mc.ForEachEntry(func(_, v any) bool {
        \t\t\tlc.Add(CheckNotNull(v))
        \t\t\treturn true
        \t\t})
        \t\treturn
        \t}
        \tif mk, ok := src.(*jLiveMapKeySet); ok {
        \t\tif mk == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\tif mk.mc.Size() == 0 {
        \t\t\treturn
        \t\t}
        \t\tlc.InitCap(lc.Size() + mk.mc.Size())
        \t\tmk.mc.ForEachEntry(func(k, _ any) bool {
        \t\t\tlc.Add(CheckNotNull(k))
        \t\t\treturn true
        \t\t})
        \t\treturn
        \t}
        \tif jsp, ok := src.(jFastJSetCoreProvider); ok {
        \t\tif sc := jsp.fastSetJCore_(); sc != nil {
        \t\t\tif sc.Size() == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tlc.InitCap(lc.Size() + sc.Size())
        \t\t\tsc.ForEach(func(k any) bool {
        \t\t\t\tlc.Add(CheckNotNull(k))
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn
        \t\t}
        \t}
        \tif sp, ok := src.(jFastSetCoreProvider); ok {
        \t\tif mc := sp.fastSetCore_(); mc != nil {
        \t\t\tif mc.Size() == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tlc.InitCap(lc.Size() + mc.Size())
        \t\t\tmc.ForEachEntry(func(k, _ any) bool {
        \t\t\t\tlc.Add(CheckNotNull(k))
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn
        \t\t}
        \t}
        \tif ep, ok := src.(jFastEnumSetProvider); ok {
        \t\tif es := ep.fastEnumSet_(); es != nil {
        \t\t\tif es.size == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tlc.InitCap(lc.Size() + es.size)
        \t\t\tes.ForEach(func(k any) bool {
        \t\t\t\tlc.Add(k)
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn
        \t\t}
        \t}
        \tif c, ok := src.(Ext_java_util_Collection); ok {
        \t\tlc.InitCap(lc.Size() + c.Size())
        \t}
        \tfor it := NewIterCursor(src); it.HasNext(); {
        \t\tlc.Add(CheckNotNull(it.Next()))
        \t}
        }

        func listCoreRemoveAll(lc *JListCore, other Ext_java_util_Collection) bool {
        \tCheckNotNull(other)
        \treturn listCoreRemoveIf(lc, func(x any) bool { return other.Contains(x) })
        }

        func listCoreRemoveIf(lc *JListCore, pred func(any) bool) bool {
        \tchanged := false
        \tw := 0
        \tfor _, x := range lc.Items {
        \t\tif pred(x) {
        \t\t\tchanged = true
        \t\t} else {
        \t\t\tlc.Items[w] = x
        \t\t\tw++
        \t\t}
        \t}
        \tfor i := w; i < len(lc.Items); i++ {
        \t\tlc.Items[i] = nil
        \t}
        \tlc.Items = lc.Items[:w]
        \treturn changed
        }

        func listCoreRemoveLastVal(lc *JListCore, v any) bool {
        \tidx := lc.LastIndexOf(v)
        \tif idx < 0 {
        \t\treturn false
        \t}
        \tlc.RemoveAt(idx)
        \treturn true
        }

        func listCoreSort(lc *JListCore, cmp Ext_java_util_Comparator) {
        \tfn := comparatorFuncOrNatural(cmp)
        \tsort.SliceStable(lc.Items, func(i, j int) bool {
        \t\treturn fn(lc.Items[i], lc.Items[j]) < 0
        \t})
        }

        func listEqualsSlice(items []any, other any) bool {
        \tlst, ok := other.(Ext_java_util_List)
        \tif !ok || lst == nil || javalang.IsNil(lst) {
        \t\treturn false
        \t}
        \tif int32(len(items)) != lst.Size() {
        \t\treturn false
        \t}
        \tit := Iterate(lst)
        \tfor _, x := range items {
        \t\tif !it.HasNext() || !javalang.ObjectsEquals(ToAny(x), ToAny(it.Next())) {
        \t\t\treturn false
        \t\t}
        \t}
        \treturn true
        }

        func listHashSlice(items []any) int32 {
        \th := int32(1)
        \tfor _, x := range items {
        \t\th = 31*h + javalang.ObjectHashCode(ToAny(x))
        \t}
        \treturn h
        }

        func sliceJavaString(items []any) string {
        \tvar sb strings.Builder
        \tsb.WriteByte('[')
        \tfor i, x := range items {
        \t\tif i > 0 {
        \t\t\tsb.WriteString(", ")
        \t\t}
        \t\tsb.WriteString(StrVal(x))
        \t}
        \tsb.WriteByte(']')
        \treturn sb.String()
        }

        func copySetCoreInto(dst *JSetCore, src *JSetCore) *JSetCore {
        \tdst.identity = src.identity
        \tdst.InitExpectedSize(src.Size())
        \tsrc.ForEach(func(k any) bool {
        \t\tdst.Add(k)
        \t\treturn true
        \t})
        \treturn dst
        }

        func setCoreAddAllIterable(sc *JSetCore, src any) bool {
        \tCheckNotNull(src)
        \tif tv, ok := src.(jTypedViewCursorSource); ok {
        \t\tif n, ok2 := tv.typedViewLen_(); ok2 {
        \t\t\tif n == 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tsc.InitExpectedSize(sc.Size() + int32(n))
        \t\t\tchanged := false
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\tif sc.Add(tv.typedViewAt_(i)) {
        \t\t\t\t\tchanged = true
        \t\t\t\t}
        \t\t\t}
        \t\t\treturn changed
        \t\t}
        \t}
        \tif si, ok := src.(SliceIterable); ok {
        \t\titems := si.IterSlice_()
        \t\tif len(items) == 0 {
        \t\t\treturn false
        \t\t}
        \t\tsc.InitExpectedSize(sc.Size() + int32(len(items)))
        \t\tchanged := false
        \t\tfor _, x := range items {
        \t\t\tif sc.Add(x) {
        \t\t\t\tchanged = true
        \t\t\t}
        \t\t}
        \t\treturn changed
        \t}
        \tif jsp, ok := src.(jFastJSetCoreProvider); ok {
        \t\tif osc := jsp.fastSetJCore_(); osc != nil && osc != sc {
        \t\t\tif osc.Size() == 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tsc.InitExpectedSize(sc.Size() + osc.Size())
        \t\t\tchanged := false
        \t\t\tosc.ForEach(func(k any) bool {
        \t\t\t\tif sc.Add(k) {
        \t\t\t\t\tchanged = true
        \t\t\t\t}
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn changed
        \t\t}
        \t}
        \tif sp, ok := src.(jFastSetCoreProvider); ok {
        \t\tif smc := sp.fastSetCore_(); smc != nil {
        \t\t\tif smc.Size() == 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tsc.InitExpectedSize(sc.Size() + smc.Size())
        \t\t\tchanged := false
        \t\t\tsmc.ForEachEntry(func(k, _ any) bool {
        \t\t\t\tif sc.Add(k) {
        \t\t\t\t\tchanged = true
        \t\t\t\t}
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn changed
        \t\t}
        \t}
        \tif ep, ok := src.(jFastEnumSetProvider); ok {
        \t\tif es := ep.fastEnumSet_(); es != nil {
        \t\t\tif es.size == 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tsc.InitExpectedSize(sc.Size() + es.size)
        \t\t\tchanged := false
        \t\t\tes.ForEach(func(k any) bool {
        \t\t\t\tif sc.Add(k) {
        \t\t\t\t\tchanged = true
        \t\t\t\t}
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn changed
        \t\t}
        \t}
        \tif c, ok := src.(Ext_java_util_Collection); ok {
        \t\tn := c.Size()
        \t\tif n == 0 {
        \t\t\treturn false
        \t\t}
        \t\tsc.InitExpectedSize(sc.Size() + n)
        \t}
        \tchanged := false
        \tfor it := NewIterCursor(src); it.HasNext(); {
        \t\tif sc.Add(it.Next()) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func setCoreAddAllNonNull(sc *JSetCore, src any) {
        \tCheckNotNull(src)
        \tif ssi, ok := src.(StrSliceIterable); ok {
        \t\tif strs, ok2 := ssi.IterStrSlice_(); ok2 {
        \t\t\tif len(strs) == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tsc.InitExpectedSize(sc.Size() + int32(len(strs)))
        \t\t\tfor _, s := range strs {
        \t\t\t\tif s == NullString {
        \t\t\t\t\tThrowNPE()
        \t\t\t\t}
        \t\t\t\tsc.AddStr(s)
        \t\t\t}
        \t\t\treturn
        \t\t}
        \t}
        \tif tv, ok := src.(jTypedViewCursorSource); ok {
        \t\tif n, ok2 := tv.typedViewLen_(); ok2 {
        \t\t\tif n == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tsc.InitExpectedSize(sc.Size() + int32(n))
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\tsc.Add(CheckNotNull(ToAny(tv.typedViewAt_(i))))
        \t\t\t}
        \t\t\treturn
        \t\t}
        \t}
        \tif si, ok := src.(SliceIterable); ok {
        \t\titems := si.IterSlice_()
        \t\tif len(items) == 0 {
        \t\t\treturn
        \t\t}
        \t\tsc.InitExpectedSize(sc.Size() + int32(len(items)))
        \t\tfor _, x := range items {
        \t\t\tsc.Add(CheckNotNull(ToAny(x)))
        \t\t}
        \t\treturn
        \t}
        \tif jsp, ok := src.(jFastJSetCoreProvider); ok {
        \t\tif osc := jsp.fastSetJCore_(); osc != nil && osc != sc {
        \t\t\tif osc.Size() == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tsc.InitExpectedSize(sc.Size() + osc.Size())
        \t\t\tosc.ForEach(func(k any) bool {
        \t\t\t\tsc.Add(CheckNotNull(k))
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn
        \t\t}
        \t}
        \tif sp, ok := src.(jFastSetCoreProvider); ok {
        \t\tif smc := sp.fastSetCore_(); smc != nil {
        \t\t\tif smc.Size() == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tsc.InitExpectedSize(sc.Size() + smc.Size())
        \t\t\tsmc.ForEachEntry(func(k, _ any) bool {
        \t\t\t\tsc.Add(CheckNotNull(k))
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn
        \t\t}
        \t}
        \tif ep, ok := src.(jFastEnumSetProvider); ok {
        \t\tif es := ep.fastEnumSet_(); es != nil {
        \t\t\tif es.size == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tsc.InitExpectedSize(sc.Size() + es.size)
        \t\t\tes.ForEach(func(k any) bool {
        \t\t\t\tsc.Add(k)
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn
        \t\t}
        \t}
        \tif c, ok := src.(Ext_java_util_Collection); ok {
        \t\tsc.InitExpectedSize(sc.Size() + c.Size())
        \t}
        \tfor it := NewIterCursor(src); it.HasNext(); {
        \t\tsc.Add(CheckNotNull(it.Next()))
        \t}
        }

        func hashSetAddAll(hs *Ext_java_util_HashSetBase, src any) bool {
        \tif mb := hs.mapBack(); mb != nil {
        \t\tCheckNotNull(src)
        \t\tchanged := false
        \t\tfor it := Iterate(src); it.HasNext(); {
        \t\t\tif mb.Put_2(it.Next(), ValueOfBoolean(true)) == nil {
        \t\t\t\tchanged = true
        \t\t\t}
        \t\t}
        \t\treturn changed
        \t}
        \treturn setCoreAddAllIterable(&hs.inlineSet, src)
        }

        func hashSetContainsAll(hs *Ext_java_util_HashSetBase, other Ext_java_util_Collection) bool {
        \tif mb := hs.mapBack(); mb != nil {
        \t\treturn colContainsAll(mb.ContainsKey, other)
        \t}
        \treturn colContainsAll(hs.inlineSet.Contains, other)
        }

        func hashSetRemoveAll(hs *Ext_java_util_HashSetBase, other Ext_java_util_Collection) bool {
        \tCheckNotNull(other)
        \tif mb := hs.mapBack(); mb != nil {
        \t\tchanged := false
        \t\tfor it := Iterate(other); it.HasNext(); {
        \t\t\tif mb.Remove_1(it.Next()) != nil {
        \t\t\t\tchanged = true
        \t\t\t}
        \t\t}
        \t\treturn changed
        \t}
        \tsc := &hs.inlineSet
        \tif sc.Size() == 0 {
        \t\treturn false
        \t}
        \tif tv, ok := any(other).(jTypedViewCursorSource); ok {
        \t\tif n, ok2 := tv.typedViewLen_(); ok2 {
        \t\t\tchanged := false
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\tif sc.Remove(tv.typedViewAt_(i)) {
        \t\t\t\t\tchanged = true
        \t\t\t\t}
        \t\t\t}
        \t\t\treturn changed
        \t\t}
        \t}
        \tif si, ok := any(other).(SliceIterable); ok {
        \t\tchanged := false
        \t\tfor _, x := range si.IterSlice_() {
        \t\t\tif sc.Remove(x) {
        \t\t\t\tchanged = true
        \t\t\t}
        \t\t}
        \t\treturn changed
        \t}
        \tif jsp, ok := any(other).(jFastJSetCoreProvider); ok {
        \t\tif osc := jsp.fastSetJCore_(); osc != nil && osc != sc {
        \t\t\tchanged := false
        \t\t\tosc.ForEach(func(k any) bool {
        \t\t\t\tif sc.Remove(k) {
        \t\t\t\t\tchanged = true
        \t\t\t\t}
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn changed
        \t\t}
        \t}
        \tchanged := false
        \tfor it := Iterate(other); it.HasNext(); {
        \t\tif sc.Remove(it.Next()) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func hashSetRetainAll(hs *Ext_java_util_HashSetBase, other Ext_java_util_Collection) bool {
        \tCheckNotNull(other)
        \tif mb := hs.mapBack(); mb != nil {
        \t\treturn abstractColRetainAll(Cast[Ext_java_util_Set](mb.KeySet_0()), other)
        \t}
        \tsc := &hs.inlineSet
        \tchanged := false
        \tfor _, k := range sc.LiveKeys() {
        \t\tif !other.Contains(k) && sc.Remove(k) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func hashSetRemoveIf(hs *Ext_java_util_HashSetBase, pred Ext_java_util_function_Predicate) bool {
        \tCheckNotNull(pred)
        \tif mb := hs.mapBack(); mb != nil {
        \t\treturn abstractColRemoveIf(Cast[Ext_java_util_Set](mb.KeySet_0()), pred)
        \t}
        \tsc := &hs.inlineSet
        \tchanged := false
        \tfor _, k := range sc.LiveKeys() {
        \t\tif pred.Test(k) && sc.Remove(k) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func hashSetToArray(hs *Ext_java_util_HashSetBase) []any {
        \tif mb := hs.mapBack(); mb != nil {
        \t\treturn colToSlice(Cast[Ext_java_util_Set](mb.KeySet_0()))
        \t}
        \treturn hs.inlineSet.LiveKeys()
        }

        func hashSetForEach(hs *Ext_java_util_HashSetBase, action Ext_java_util_function_Consumer) {
        \tCheckNotNull(action)
        \tif mb := hs.mapBack(); mb != nil {
        \t\tfor it := Iterate(mb.KeySet_0()); it.HasNext(); {
        \t\t\taction.Accept(it.Next())
        \t\t}
        \t\treturn
        \t}
        \tfor _, k := range hs.inlineSet.LiveKeys() {
        \t\taction.Accept(k)
        \t}
        }

        func hashSetEquals(hs *Ext_java_util_HashSetBase, other any) bool {
        \tst, ok := other.(Ext_java_util_Set)
        \tif !ok || st == nil || javalang.IsNil(st) {
        \t\treturn false
        \t}
        \tif mb := hs.mapBack(); mb != nil {
        \t\tif mb.Size() != st.Size() {
        \t\t\treturn false
        \t\t}
        \t\treturn colContainsAll(mb.ContainsKey, st)
        \t}
        \tsc := &hs.inlineSet
        \tif sc.Size() != st.Size() {
        \t\treturn false
        \t}
        \tif sc.Size() == 0 {
        \t\treturn true
        \t}
        \tif jsp, ok := other.(jFastJSetCoreProvider); ok {
        \t\tif osc := jsp.fastSetJCore_(); osc != nil {
        \t\t\tif sc == osc {
        \t\t\t\treturn true
        \t\t\t}
        \t\t\tall := true
        \t\t\tosc.ForEach(func(k any) bool {
        \t\t\t\tif !sc.Contains(k) {
        \t\t\t\t\tall = false
        \t\t\t\t\treturn false
        \t\t\t\t}
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn all
        \t\t}
        \t}
        \treturn colContainsAll(sc.Contains, st)
        }

        func hashSetHash(hs *Ext_java_util_HashSetBase) int32 {
        \tif mb := hs.mapBack(); mb != nil {
        \t\treturn javalang.ObjectHashCode(mb.KeySet_0())
        \t}
        \treturn hs.inlineSet.SetHash()
        }

        type jSetIter struct {
        \tsc      *JSetCore
        \tidx     int
        \tlastVal any
        \thasLast bool
        }

        func newHashSetIter(hs *Ext_java_util_HashSetBase) Ext_java_util_Iterator {
        \tif mb := hs.mapBack(); mb != nil {
        \t\treturn Cast[Ext_java_util_Iterator](Cast[Ext_java_util_Set](mb.KeySet_0()).Iterator())
        \t}
        \tsc := &hs.inlineSet
        \tif sc.Size() == 0 {
        \t\treturn emptyMapIter_
        \t}
        \tstart := int(sc.firstIdx)
        \tit := &jSetIter{sc: sc, idx: start}
        \tit.advance()
        \tif it.idx != start && !sc.imm {
        \t\tsc.firstIdx = int32(it.idx)
        \t}
        \treturn it
        }

        func (*jSetIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jSetIter) sliceLen() int {
        \tswitch it.sc.mode {
        \tcase jSetModeNode:
        \t\treturn len(it.sc.nodes)
        \tcase jSetModeStr:
        \t\treturn len(it.sc.ext.(*jSetStrState).strs)
        \tcase jSetModeAny:
        \t\treturn len(it.sc.ext.(*jSetAnyState).items)
        \tdefault:
        \t\treturn 0
        \t}
        }
        func (it *jSetIter) advance() {
        \tswitch it.sc.mode {
        \tcase jSetModeNode:
        \t\tfor it.idx < len(it.sc.nodes) && it.sc.nodes[it.idx] == jTombstoneNode {
        \t\t\tit.idx++
        \t\t}
        \tcase jSetModeStr:
        \t\tstrs := it.sc.ext.(*jSetStrState).strs
        \t\tfor it.idx < len(strs) && strs[it.idx] == jTombstoneStr {
        \t\t\tit.idx++
        \t\t}
        \tcase jSetModeAny:
        \t\titems := it.sc.ext.(*jSetAnyState).items
        \t\tfor it.idx < len(items) && items[it.idx] == jTombstoneObj {
        \t\t\tit.idx++
        \t\t}
        \t}
        }
        func (it *jSetIter) HasNext() bool {
        \tit.advance()
        \treturn it.idx < it.sliceLen()
        }
        func (it *jSetIter) Next() any {
        \tit.advance()
        \tif it.idx >= it.sliceLen() {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tvar v any
        \tswitch it.sc.mode {
        \tcase jSetModeNode:
        \t\tv = boxNodeById(it.sc.nodes[it.idx])
        \tcase jSetModeStr:
        \t\ts := it.sc.ext.(*jSetStrState).strs[it.idx]
        \t\tif s != NullString {
        \t\t\tv = s
        \t\t}
        \tcase jSetModeAny:
        \t\tv = it.sc.ext.(*jSetAnyState).items[it.idx]
        \t}
        \tit.idx++
        \tit.lastVal = v
        \tit.hasLast = true
        \tit.advance()
        \treturn v
        }
        func (it *jSetIter) Remove() {
        \tif !it.hasLast {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \tif it.sc.imm {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tit.sc.Remove(it.lastVal)
        \tit.hasLast = false
        }

        const jImmEnumMaxBits = 256

        type jEnumOrdAccessor interface {
        \tEnumOrdinal_() int32
        \tJavaClassInfo_() (string, string)
        }

        type immSetIdx struct {
        \tstrIdx      map[string]struct{}
        \tidentIdx    map[any]struct{}
        \tnodeIdx     map[uint32]struct{}
        \tbig         *JSetCore
        \tenumClass   string
        \tenumBits    [4]uint64
        \thasEnumBits bool
        }

        func tryInitImmEnumSingle(obj *Ext_com_google_common_collect_ImmutableSetBase, e0 any) {
        \tif e0 == nil {
        \t\treturn
        \t}
        \trt0 := reflect.TypeOf(e0)
        \tif rt0.Kind() == reflect.Uint8 {
        \t\tu := uint8(reflect.ValueOf(e0).Uint())
        \t\tif u != 0 {
        \t\t\ttab := new([256]bool)
        \t\t\ttab[u] = true
        \t\t\tobj.u8Tab_ = tab
        \t\t\tobj.u8PtrNil_ = reflect.Zero(reflect.PointerTo(rt0)).Interface()
        \t\t\tobj.hasU8Bits_ = true
        \t\t}
        \t}
        }

        func tryInitImmEnumMulti(obj *Ext_com_google_common_collect_ImmutableSetBase, elems []any) bool {
        \tif len(elems) == 0 || elems[0] == nil {
        \t\treturn false
        \t}
        \trt0 := reflect.TypeOf(elems[0])
        \tif rt0.Kind() == reflect.Uint8 {
        \t\tu0 := uint8(reflect.ValueOf(elems[0]).Uint())
        \t\tif u0 == 0 {
        \t\t\treturn false
        \t\t}
        \t\ttab := new([256]bool)
        \t\ttab[u0] = true
        \t\tfor i := 1; i < len(elems); i++ {
        \t\t\tei := elems[i]
        \t\t\tif ei == nil || reflect.TypeOf(ei) != rt0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tu := uint8(reflect.ValueOf(ei).Uint())
        \t\t\tif u == 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\ttab[u] = true
        \t\t}
        \t\tobj.u8Tab_ = tab
        \t\tobj.u8PtrNil_ = reflect.Zero(reflect.PointerTo(rt0)).Interface()
        \t\tobj.hasU8Bits_ = true
        \t\treturn true
        \t}
        \tif len(elems) <= 4 {
        \t\treturn false
        \t}
        \te0, ok := elems[0].(jEnumOrdAccessor)
        \tif !ok {
        \t\treturn false
        \t}
        \tord0 := e0.EnumOrdinal_()
        \tif ord0 < 0 || ord0 >= jImmEnumMaxBits {
        \t\treturn false
        \t}
        \tcls0, _ := e0.JavaClassInfo_()
        \tvar maskBits [4]uint64
        \tmaskBits[uint(ord0)>>6] |= uint64(1) << (uint(ord0) & 63)
        \tfor i := 1; i < len(elems); i++ {
        \t\tei := elems[i]
        \t\tif ei == nil {
        \t\t\treturn false
        \t\t}
        \t\teai, ok2 := ei.(jEnumOrdAccessor)
        \t\tif !ok2 {
        \t\t\treturn false
        \t\t}
        \t\tclsi, _ := eai.JavaClassInfo_()
        \t\tif clsi != cls0 {
        \t\t\treturn false
        \t\t}
        \t\tord := eai.EnumOrdinal_()
        \t\tif ord < 0 || ord >= jImmEnumMaxBits {
        \t\t\treturn false
        \t\t}
        \t\tmaskBits[uint(ord)>>6] |= uint64(1) << (uint(ord) & 63)
        \t}
        \tobj.idx_.Store(&immSetIdx{
        \t\tenumClass:   cls0,
        \t\tenumBits:    maskBits,
        \t\thasEnumBits: true,
        \t})
        \treturn true
        }

        func isAllFastImmSetElems(elems []any) bool {
        \tallStr := true
        \tfor _, e := range elems {
        \t\tstr, ok := e.(string)
        \t\tif !ok || str == NullString {
        \t\t\tallStr = false
        \t\t\tbreak
        \t\t}
        \t}
        \tif allStr {
        \t\treturn true
        \t}
        \tallNode := true
        \tfor _, e := range elems {
        \t\tif e == nil {
        \t\t\tallNode = false
        \t\t\tbreak
        \t\t}
        \t\tif rn, ok := e.(jRawNodeProvider); ok {
        \t\t\tif rn.RawNodeId_() == 0 {
        \t\t\t\tallNode = false
        \t\t\t\tbreak
        \t\t\t}
        \t\t} else if nd, ok := e.(nodestore.Node); ok {
        \t\t\tif nd == 0 {
        \t\t\t\tallNode = false
        \t\t\t\tbreak
        \t\t\t}
        \t\t} else {
        \t\t\tallNode = false
        \t\t\tbreak
        \t\t}
        \t}
        \tif allNode {
        \t\treturn true
        \t}
        \tfor _, e := range elems {
        \t\tif !isFastIdentElem(e) {
        \t\t\treturn false
        \t\t}
        \t}
        \treturn true
        }

        func finishImmSetMulti(elems []any, scOwned *JSetCore) *Ext_com_google_common_collect_ImmutableSetBase {
        \tobj := &Ext_com_google_common_collect_ImmutableSetBase{elems_: elems[:len(elems):len(elems)], u8Tab_: &emptyU8Tab_}
        \tobj.BindExtSelf_(obj)
        \tif tryInitImmEnumMulti(obj, elems) {
        \t\treturn obj
        \t}
        \tif len(elems) > 8 && scOwned != nil && !isAllFastImmSetElems(elems) {
        \t\tscOwned.sealForImmutable()
        \t\tobj.idx_.Store(&immSetIdx{big: scOwned})
        \t}
        \treturn obj
        }

        func singletonElemMatchesNode(e any, n uint32) bool {
        \tif e == nil || n == 0 {
        \t\treturn false
        \t}
        \tif rn, ok := e.(jRawNodeProvider); ok {
        \t\treturn rn.RawNodeId_() == n
        \t}
        \tif nd, ok := e.(nodestore.Node); ok {
        \t\treturn uint32(nd) == n
        \t}
        \treturn false
        }

        func colSetEqualsSingleton(elem any, other any) bool {
        \tif other == nil {
        \t\treturn false
        \t}
        \tst, ok := other.(Ext_java_util_Set)
        \tif !ok || st == nil || javalang.IsNil(st) {
        \t\treturn false
        \t}
        \tif st.Size() != 1 {
        \t\treturn false
        \t}
        \tif si, ok := other.(SliceIterable); ok {
        \t\ts := si.IterSlice_()
        \t\treturn len(s) == 1 && jSetElemEquals(elem, s[0])
        \t}
        \tit := Iterate(st)
        \treturn it.HasNext() && jSetElemEquals(elem, it.Next())
        }

        func strsToAnySlice(strs []string) []any {
        \tif len(strs) == 0 {
        \t\treturn []any{}
        \t}
        \tout := make([]any, len(strs))
        \tfor i, s := range strs {
        \t\tif s != NullString {
        \t\t\tout[i] = s
        \t\t}
        \t}
        \treturn out
        }

        func immStrSliceContains(strs []string, arg0 any) bool {
        \ts, ok := arg0.(string)
        \tif !ok || s == NullString {
        \t\treturn false
        \t}
        \tfor _, x := range strs {
        \t\tif x == s {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func immStrSliceIndexOf(strs []string, arg0 any) int32 {
        \ts, ok := arg0.(string)
        \tif !ok || s == NullString {
        \t\treturn -1
        \t}
        \tfor i, x := range strs {
        \t\tif x == s {
        \t\t\treturn int32(i)
        \t\t}
        \t}
        \treturn -1
        }

        func immStrSliceLastIndexOf(strs []string, arg0 any) int32 {
        \ts, ok := arg0.(string)
        \tif !ok || s == NullString {
        \t\treturn -1
        \t}
        \tfor i := len(strs) - 1; i >= 0; i-- {
        \t\tif strs[i] == s {
        \t\t\treturn int32(i)
        \t\t}
        \t}
        \treturn -1
        }

        func listEqualsStrSlice(strs []string, other any) bool {
        \tif other == nil {
        \t\treturn false
        \t}
        \tlst, ok := other.(Ext_java_util_List)
        \tif !ok || lst == nil || javalang.IsNil(lst) {
        \t\treturn false
        \t}
        \tn := len(strs)
        \tif int32(n) != lst.Size() {
        \t\treturn false
        \t}
        \tif n == 0 {
        \t\treturn true
        \t}
        \tif ssi, ok2 := other.(StrSliceIterable); ok2 {
        \t\tif os, ok3 := ssi.IterStrSlice_(); ok3 {
        \t\t\tif len(os) != n {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\tif strs[i] != os[i] {
        \t\t\t\t\treturn false
        \t\t\t\t}
        \t\t\t}
        \t\t\treturn true
        \t\t}
        \t}
        \tif si, ok2 := other.(SliceIterable); ok2 {
        \t\tos := si.IterSlice_()
        \t\tif len(os) != n {
        \t\t\treturn false
        \t\t}
        \t\tfor i := 0; i < n; i++ {
        \t\t\tes, ok3 := os[i].(string)
        \t\t\tif !ok3 || es != strs[i] {
        \t\t\t\treturn false
        \t\t\t}
        \t\t}
        \t\treturn true
        \t}
        \tit := Iterate(lst)
        \tfor i := 0; i < n; i++ {
        \t\tif !it.HasNext() {
        \t\t\treturn false
        \t\t}
        \t\tes, ok2 := it.Next().(string)
        \t\tif !ok2 || es != strs[i] {
        \t\t\treturn false
        \t\t}
        \t}
        \treturn !it.HasNext()
        }

        func listHashStrSlice(strs []string) int32 {
        \th := int32(1)
        \tfor _, s := range strs {
        \t\tif s == NullString {
        \t\t\th = 31 * h
        \t\t} else {
        \t\t\th = 31*h + javalang.StringHashCode(s)
        \t\t}
        \t}
        \treturn h
        }

        func strSliceJavaString(strs []string) string {
        \tif len(strs) == 0 {
        \t\treturn "[]"
        \t}
        \tif len(strs) == 1 {
        \t\tif strs[0] == NullString {
        \t\t\treturn "[null]"
        \t\t}
        \t\treturn "[" + strs[0] + "]"
        \t}
        \tvar sb strings.Builder
        \tsb.WriteByte('[')
        \tfor i, s := range strs {
        \t\tif i > 0 {
        \t\t\tsb.WriteString(", ")
        \t\t}
        \t\tif s == NullString {
        \t\t\tsb.WriteString("null")
        \t\t} else {
        \t\t\tsb.WriteString(s)
        \t\t}
        \t}
        \tsb.WriteByte(']')
        \treturn sb.String()
        }

        func immStrSetHash(strs []string) int32 {
        \tvar h int32
        \tfor _, s := range strs {
        \t\tif s != NullString {
        \t\t\th += javalang.StringHashCode(s)
        \t\t}
        \t}
        \treturn h
        }

        func (s *Ext_com_google_common_collect_ImmutableSetBase) getOrBuildIdx() *immSetIdx {
        \tif idx := s.idx_.Load(); idx != nil {
        \t\treturn idx
        \t}
        \treturn s.buildIdxSlow()
        }

        //go:noinline
        func (s *Ext_com_google_common_collect_ImmutableSetBase) buildIdxSlow() *immSetIdx {
        \tbuilt := &immSetIdx{}
        \tif len(s.strs_) > 4 {
        \t\tidx := make(map[string]struct{}, len(s.strs_))
        \t\tfor _, str := range s.strs_ {
        \t\t\tidx[str] = struct{}{}
        \t\t}
        \t\tbuilt.strIdx = idx
        \t} else if len(s.elems_) > 4 {
        \t\tallStr := true
        \t\tfor _, e := range s.elems_ {
        \t\t\tstr, ok := e.(string)
        \t\t\tif !ok || str == NullString {
        \t\t\t\tallStr = false
        \t\t\t\tbreak
        \t\t\t}
        \t\t}
        \t\tif allStr {
        \t\t\tidx := make(map[string]struct{}, len(s.elems_))
        \t\t\tfor _, e := range s.elems_ {
        \t\t\t\tidx[e.(string)] = struct{}{}
        \t\t\t}
        \t\t\tbuilt.strIdx = idx
        \t\t} else {
        \t\t\tallNode := true
        \t\t\tfor _, e := range s.elems_ {
        \t\t\t\tif e == nil {
        \t\t\t\t\tallNode = false
        \t\t\t\t\tbreak
        \t\t\t\t}
        \t\t\t\tif rn, ok := e.(jRawNodeProvider); ok {
        \t\t\t\t\tif rn.RawNodeId_() == 0 {
        \t\t\t\t\t\tallNode = false
        \t\t\t\t\t\tbreak
        \t\t\t\t\t}
        \t\t\t\t} else if nd, ok := e.(nodestore.Node); ok {
        \t\t\t\t\tif nd == 0 {
        \t\t\t\t\t\tallNode = false
        \t\t\t\t\t\tbreak
        \t\t\t\t\t}
        \t\t\t\t} else {
        \t\t\t\t\tallNode = false
        \t\t\t\t\tbreak
        \t\t\t\t}
        \t\t\t}
        \t\t\tif allNode {
        \t\t\t\tidx := make(map[uint32]struct{}, len(s.elems_))
        \t\t\t\tfor _, e := range s.elems_ {
        \t\t\t\t\tif rn, ok := e.(jRawNodeProvider); ok {
        \t\t\t\t\t\tidx[rn.RawNodeId_()] = struct{}{}
        \t\t\t\t\t} else {
        \t\t\t\t\t\tidx[uint32(e.(nodestore.Node))] = struct{}{}
        \t\t\t\t\t}
        \t\t\t\t}
        \t\t\t\tbuilt.nodeIdx = idx
        \t\t\t} else {
        \t\t\t\tallIdent := true
        \t\t\t\tfor _, e := range s.elems_ {
        \t\t\t\t\tif !isFastIdentElem(e) {
        \t\t\t\t\t\tallIdent = false
        \t\t\t\t\t\tbreak
        \t\t\t\t\t}
        \t\t\t\t}
        \t\t\t\tif allIdent {
        \t\t\t\t\tidx := make(map[any]struct{}, len(s.elems_))
        \t\t\t\t\tfor _, e := range s.elems_ {
        \t\t\t\t\t\tidx[e] = struct{}{}
        \t\t\t\t\t}
        \t\t\t\t\tbuilt.identIdx = idx
        \t\t\t\t} else if len(s.elems_) > 8 {
        \t\t\t\t\tbig := &JSetCore{}
        \t\t\t\t\tbig.InitExpectedSize(int32(len(s.elems_)))
        \t\t\t\t\tfor _, e := range s.elems_ {
        \t\t\t\t\t\tbig.Add(e)
        \t\t\t\t\t}
        \t\t\t\t\tbig.sealForImmutable()
        \t\t\t\t\tbuilt.big = big
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t}
        \tif s.idx_.CompareAndSwap(nil, built) {
        \t\treturn built
        \t}
        \treturn s.idx_.Load()
        }

        func (s *Ext_com_google_common_collect_ImmutableSetBase) immSetContainsStr(k string) bool {
        \tif k == NullString || s.hasU8Bits_ {
        \t\treturn false
        \t}
        \tif s.strs_ != nil {
        \t\tif len(s.strs_) == 0 {
        \t\t\treturn false
        \t\t}
        \t\tif len(s.strs_) > 4 {
        \t\t\tidx := s.getOrBuildIdx()
        \t\t\t_, hit := idx.strIdx[k]
        \t\t\treturn hit
        \t\t}
        \t\tfor _, es := range s.strs_ {
        \t\t\tif es == k {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \t}
        \telems := s.elems_
        \tif len(elems) == 0 {
        \t\treturn false
        \t}
        \tif len(elems) <= 4 {
        \t\tfor _, e := range elems {
        \t\t\tif es, ok := e.(string); ok && es == k {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \t}
        \tidx := s.getOrBuildIdx()
        \tif idx.hasEnumBits {
        \t\treturn false
        \t}
        \tif idx.strIdx != nil {
        \t\t_, hit := idx.strIdx[k]
        \t\treturn hit
        \t}
        \tif idx.nodeIdx != nil || idx.identIdx != nil {
        \t\treturn false
        \t}
        \tif idx.big != nil {
        \t\treturn idx.big.ContainsStr(k)
        \t}
        \tfor _, e := range elems {
        \t\tif es, ok := e.(string); ok && es == k {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (s *Ext_com_google_common_collect_ImmutableSetBase) immSetContainsNode(n uint32) bool {
        \tif n == 0 || s.strs_ != nil || s.hasU8Bits_ {
        \t\treturn false
        \t}
        \telems := s.elems_
        \tif len(elems) == 0 {
        \t\treturn false
        \t}
        \tif len(elems) <= 4 {
        \t\tfor _, e := range elems {
        \t\t\tif singletonElemMatchesNode(e, n) {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \t}
        \tidx := s.getOrBuildIdx()
        \tif idx.hasEnumBits {
        \t\treturn false
        \t}
        \tif idx.nodeIdx != nil {
        \t\t_, hit := idx.nodeIdx[n]
        \t\treturn hit
        \t}
        \tif idx.strIdx != nil || idx.identIdx != nil {
        \t\treturn false
        \t}
        \tif idx.big != nil {
        \t\treturn idx.big.ContainsNode(n)
        \t}
        \tfor _, e := range elems {
        \t\tif singletonElemMatchesNode(e, n) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func immSmallSetContains(elems []any, v any) bool {
        \tswitch x := v.(type) {
        \tcase string:
        \t\tif x == NullString {
        \t\t\treturn false
        \t\t}
        \t\tfor _, e := range elems {
        \t\t\tif es, ok := e.(string); ok && es == x {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \tcase jIdentityKey:
        \t\tif x.IsJavaNil_() {
        \t\t\treturn false
        \t\t}
        \t\tfor _, e := range elems {
        \t\t\tif e == v {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \tcase jEnumOrdAccessor:
        \t\tif x.EnumOrdinal_() < 0 {
        \t\t\treturn false
        \t\t}
        \t\tfor _, e := range elems {
        \t\t\tif e == v {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \tcase *BoxInt:
        \t\tif x == nil {
        \t\t\treturn false
        \t\t}
        \t\txv := x.Val
        \t\tfor _, e := range elems {
        \t\t\tif eb, ok := e.(*BoxInt); ok {
        \t\t\t\tif eb != nil && eb.Val == xv {
        \t\t\t\t\treturn true
        \t\t\t\t}
        \t\t\t} else if ei, ok2 := e.(int32); ok2 && ei == xv {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \tcase int32:
        \t\tfor _, e := range elems {
        \t\t\tif ei, ok := e.(int32); ok {
        \t\t\t\tif ei == x {
        \t\t\t\t\treturn true
        \t\t\t\t}
        \t\t\t} else if eb, ok2 := e.(*BoxInt); ok2 && eb != nil && eb.Val == x {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \tcase jRawNodeProvider:
        \t\tuid := x.RawNodeId_()
        \t\tif uid == 0 {
        \t\t\treturn false
        \t\t}
        \t\tfor _, e := range elems {
        \t\t\tif singletonElemMatchesNode(e, uid) {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \tcase nodestore.Node:
        \t\tif x == 0 {
        \t\t\treturn false
        \t\t}
        \t\tuid := uint32(x)
        \t\tfor _, e := range elems {
        \t\t\tif singletonElemMatchesNode(e, uid) {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \tcase FastCustomKey:
        \t\tif x.IsJavaNil_() {
        \t\t\treturn false
        \t\t}
        \t\tfor _, e := range elems {
        \t\t\tif x.Equals(e) {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \tcase IdentityComparable:
        \t\tif nc, ok := v.(jJavaNilChecker); ok && nc.IsJavaNil_() {
        \t\t\treturn false
        \t\t}
        \t\tfor _, e := range elems {
        \t\t\tif e == v {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \tcase *javalang.Class:
        \t\tif x == nil {
        \t\t\treturn false
        \t\t}
        \t\tfor _, e := range elems {
        \t\t\tif e == v {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \t}
        \tv = ToAny(v)
        \tif v == nil {
        \t\treturn false
        \t}
        \tfor _, e := range elems {
        \t\tif jSetElemEquals(e, v) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (b *Ext_com_google_common_collect_ImmutableSetBase) ContainsU8_(u uint8) bool {
        \tif b != nil && b.u8Tab_ != nil {
        \t\treturn b.u8Tab_[u]
        \t}
        \treturn b.immSetContainsU8Fallback(u)
        }

        //go:noinline
        func (b *Ext_com_google_common_collect_ImmutableSetBase) immSetContainsU8Fallback(u uint8) bool {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tif b.hasU8Bits_ && b.u8Tab_ != nil {
        \t\treturn b.u8Tab_[u]
        \t}
        \treturn false
        }

        func (s *Ext_com_google_common_collect_ImmutableSetBase) immSetContains(v any) bool {
        \tif s.hasU8Bits_ {
        \t\tif s.u8PtrNil_ == nil || v == nil {
        \t\t\treturn false
        \t\t}
        \t\tif ea, ok := v.(jEnumOrdAccessor); ok {
        \t\t\tu := ea.EnumOrdinal_() + 1
        \t\t\tif u <= 0 || u >= 256 || !s.u8Tab_[u] {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\treturn reflect.TypeOf(v) == reflect.TypeOf(s.elems_[0])
        \t\t}
        \t\treturn false
        \t}
        \tif s.strs_ != nil {
        \t\tif len(s.strs_) == 0 || v == nil {
        \t\t\treturn false
        \t\t}
        \t\tstr, ok := v.(string)
        \t\tif !ok || str == NullString {
        \t\t\treturn false
        \t\t}
        \t\tif len(s.strs_) > 4 {
        \t\t\tidx := s.getOrBuildIdx()
        \t\t\t_, hit := idx.strIdx[str]
        \t\t\treturn hit
        \t\t}
        \t\tfor _, es := range s.strs_ {
        \t\t\tif es == str {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \t}
        \telems := s.elems_
        \tif len(elems) == 0 || v == nil {
        \t\treturn false
        \t}
        \tif len(elems) <= 4 {
        \t\treturn immSmallSetContains(elems, v)
        \t}
        \tidx := s.idx_.Load()
        \tif idx == nil {
        \t\tidx = s.buildIdxSlow()
        \t}
        \tif idx.hasEnumBits {
        \t\tea, ok := v.(jEnumOrdAccessor)
        \t\tif !ok {
        \t\t\treturn false
        \t\t}
        \t\tord := ea.EnumOrdinal_()
        \t\tif ord < 0 || ord >= jImmEnumMaxBits || (idx.enumBits[uint(ord)>>6]&(uint64(1)<<(uint(ord)&63))) == 0 {
        \t\t\treturn false
        \t\t}
        \t\tcls, _ := ea.JavaClassInfo_()
        \t\treturn cls == idx.enumClass
        \t}
        \tif idx.strIdx != nil {
        \t\tstr, ok := v.(string)
        \t\tif !ok || str == NullString {
        \t\t\treturn false
        \t\t}
        \t\t_, hit := idx.strIdx[str]
        \t\treturn hit
        \t}
        \tif idx.nodeIdx != nil {
        \t\tif rn, ok := v.(jRawNodeProvider); ok {
        \t\t\tu := rn.RawNodeId_()
        \t\t\tif u == 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\t_, hit := idx.nodeIdx[u]
        \t\t\treturn hit
        \t\t}
        \t\tif nd, ok := v.(nodestore.Node); ok {
        \t\t\tif nd == 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\t_, hit := idx.nodeIdx[uint32(nd)]
        \t\t\treturn hit
        \t\t}
        \t\treturn false
        \t}
        \tif idx.identIdx != nil {
        \t\tif !isFastIdentElem(v) {
        \t\t\treturn false
        \t\t}
        \t\t_, hit := idx.identIdx[v]
        \t\treturn hit
        \t}
        \tif idx.big != nil {
        \t\treturn idx.big.Contains(v)
        \t}
        \treturn immSmallSetContains(elems, v)
        }

        func findSetCoreEqVal[K FastEqValKey](sc *JSetCore, k K, h int32, nonNull, isIdent bool) bool {
        \tif sc.live == 0 {
        \t\treturn false
        \t}
        \tif !nonNull {
        \t\tswitch sc.mode {
        \t\tcase jSetModeNode:
        \t\t\treturn sc.ContainsNode(0)
        \t\tcase jSetModeStr:
        \t\t\treturn sc.ContainsStr(NullString)
        \t\tcase jSetModeAny:
        \t\t\treturn sc.containsAnyNormalized(nil)
        \t\tdefault:
        \t\t\treturn false
        \t\t}
        \t}
        \tif sc.mode != jSetModeAny {
        \t\treturn false
        \t}
        \tst := sc.ext.(*jSetAnyState)
        \tif st.allIdent {
        \t\tif !isIdent {
        \t\t\treturn false
        \t\t}
        \t\tfor _, it := range st.items {
        \t\t\tif ek, ok := it.(K); ok && ek == k {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \t}
        \tif g := st.gen; g != nil && len(g.buckets) > 0 {
        \t\tslot := smearHash(h) & uint32(len(g.buckets)-1)
        \t\tfor cur := g.buckets[slot]; cur != 0; cur = g.hashNext[cur-1] {
        \t\t\ti := cur - 1
        \t\t\tif g.hashes[i] == h {
        \t\t\t\tif ek, ok := st.items[i].(K); ok && ek == k {
        \t\t\t\t\treturn true
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t\treturn false
        \t}
        \tfor _, it := range st.items {
        \t\tif ek, ok := it.(K); ok && ek == k {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func fastSetContainsEqVal[K FastEqValKey](s any, k K) (bool, bool) {
        \tif isb, ok := s.(*Ext_com_google_common_collect_ImmutableSetBase); ok {
        \t\tif isb == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\tif isb.Self_ != nil && isb.Self_ != Ext_java_util_AbstractCollection(isb) {
        \t\t\treturn false, false
        \t\t}
        \t\tif isb.strs_ != nil || isb.hasU8Bits_ || len(isb.elems_) == 0 {
        \t\t\treturn false, true
        \t\t}
        \t\th, nonNull, isIdent := k.FastValHash_()
        \t\tif !nonNull {
        \t\t\treturn false, true
        \t\t}
        \t\telems := isb.elems_
        \t\tif len(elems) <= 8 {
        \t\t\tfor _, e := range elems {
        \t\t\t\tif ek, ok2 := e.(K); ok2 && ek == k {
        \t\t\t\t\treturn true, true
        \t\t\t\t}
        \t\t\t}
        \t\t\treturn false, true
        \t\t}
        \t\tidx := isb.getOrBuildIdx()
        \t\tif idx.hasEnumBits || idx.strIdx != nil || idx.nodeIdx != nil {
        \t\t\treturn false, true
        \t\t}
        \t\tif idx.big != nil {
        \t\t\treturn findSetCoreEqVal(idx.big, k, h, true, isIdent), true
        \t\t}
        \t\tfor _, e := range elems {
        \t\t\tif ek, ok2 := e.(K); ok2 && ek == k {
        \t\t\t\treturn true, true
        \t\t\t}
        \t\t}
        \t\treturn false, true
        \t}
        \tif jsp, ok := s.(jFastJSetCoreProvider); ok {
        \t\tif sc := jsp.fastSetJCore_(); sc != nil {
        \t\t\th, nonNull, isIdent := k.FastValHash_()
        \t\t\treturn findSetCoreEqVal(sc, k, h, nonNull, isIdent), true
        \t\t}
        \t}
        \tif sp, ok := s.(jFastSetCoreProvider); ok {
        \t\tif mc := sp.fastSetCore_(); mc != nil {
        \t\t\t_, found := findMapCoreEqVal(mc, k)
        \t\t\treturn found, true
        \t\t}
        \t}
        \treturn false, false
        }

        func immSetEquals(s *Ext_com_google_common_collect_ImmutableSetBase, other any) bool {
        \tif other == nil {
        \t\treturn false
        \t}
        \tif os, ok := other.(*Ext_com_google_common_collect_ImmutableSetBase); ok && os != nil {
        \t\tif os == s {
        \t\t\treturn true
        \t\t}
        \t\tif s.hasU8Bits_ && os.hasU8Bits_ {
        \t\t\treturn s.u8PtrNil_ == os.u8PtrNil_ && len(s.elems_) == len(os.elems_) && *s.u8Tab_ == *os.u8Tab_
        \t\t}
        \t}
        \tst, ok := other.(Ext_java_util_Set)
        \tif !ok || st == nil || javalang.IsNil(st) {
        \t\treturn false
        \t}
        \tn := int32(len(s.elems_))
        \tif s.strs_ != nil {
        \t\tn = int32(len(s.strs_))
        \t}
        \tif n != st.Size() {
        \t\treturn false
        \t}
        \tif n == 0 {
        \t\treturn true
        \t}
        \tif s.strs_ != nil {
        \t\tif ssi, ok2 := other.(StrSliceIterable); ok2 {
        \t\t\tif os, ok3 := ssi.IterStrSlice_(); ok3 {
        \t\t\t\tfor _, x := range os {
        \t\t\t\t\tif !s.immSetContainsStr(x) {
        \t\t\t\t\t\treturn false
        \t\t\t\t\t}
        \t\t\t\t}
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t}
        \treturn colContainsAll(s.immSetContains, st)
        }

        func immSetHash(elems []any) int32 {
        \tvar h int32
        \tfor _, e := range elems {
        \t\th += elemJavaHashCode(e, false)
        \t}
        \treturn h
        }
        """;
  }

  private static String buildColRtPart1a() {
    return """

        func setCoreEquals(mc *JMapCore, other any) bool {
        \tst, ok := other.(Ext_java_util_Set)
        \tif !ok || st == nil || javalang.IsNil(st) {
        \t\treturn false
        \t}
        \tif mc.Size() != st.Size() {
        \t\treturn false
        \t}
        \tif mc.Size() == 0 {
        \t\treturn true
        \t}
        \treturn colContainsAll(mc.ContainsKey, st)
        }

        func setCoreHash(mc *JMapCore) int32 {
        \treturn mc.SetHash()
        }

        func mapCoreAddAllIterable(mc *JMapCore, src any) bool {
        \tCheckNotNull(src)
        \tif tv, ok := src.(jTypedViewCursorSource); ok {
        \t\tif n, ok2 := tv.typedViewLen_(); ok2 {
        \t\t\tif n == 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tmc.InitExpectedSize(mc.Size() + int32(n))
        \t\t\tchanged := false
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\tif mc.AddKey(tv.typedViewAt_(i)) {
        \t\t\t\t\tchanged = true
        \t\t\t\t}
        \t\t\t}
        \t\t\treturn changed
        \t\t}
        \t}
        \tif si, ok := src.(SliceIterable); ok {
        \t\titems := si.IterSlice_()
        \t\tif len(items) == 0 {
        \t\t\treturn false
        \t\t}
        \t\tmc.InitExpectedSize(mc.Size() + int32(len(items)))
        \t\tchanged := false
        \t\tfor _, x := range items {
        \t\t\tif mc.AddKey(x) {
        \t\t\t\tchanged = true
        \t\t\t}
        \t\t}
        \t\treturn changed
        \t}
        \tif jsp, ok := src.(jFastJSetCoreProvider); ok {
        \t\tif sc := jsp.fastSetJCore_(); sc != nil {
        \t\t\tif sc.Size() == 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tmc.InitExpectedSize(mc.Size() + sc.Size())
        \t\t\tchanged := false
        \t\t\tsc.ForEach(func(k any) bool {
        \t\t\t\tif mc.AddKey(k) {
        \t\t\t\t\tchanged = true
        \t\t\t\t}
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn changed
        \t\t}
        \t}
        \tif sp, ok := src.(jFastSetCoreProvider); ok {
        \t\tif smc := sp.fastSetCore_(); smc != nil {
        \t\t\tif smc.Size() == 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tmc.InitExpectedSize(mc.Size() + smc.Size())
        \t\t\tchanged := false
        \t\t\tsmc.ForEachEntry(func(k, _ any) bool {
        \t\t\t\tif mc.AddKey(k) {
        \t\t\t\t\tchanged = true
        \t\t\t\t}
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn changed
        \t\t}
        \t}
        \tif ep, ok := src.(jFastEnumSetProvider); ok {
        \t\tif es := ep.fastEnumSet_(); es != nil {
        \t\t\tif es.size == 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tmc.InitExpectedSize(mc.Size() + es.size)
        \t\t\tchanged := false
        \t\t\tes.ForEach(func(k any) bool {
        \t\t\t\tif mc.AddKey(k) {
        \t\t\t\t\tchanged = true
        \t\t\t\t}
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn changed
        \t\t}
        \t}
        \tif c, ok := src.(Ext_java_util_Collection); ok {
        \t\tn := c.Size()
        \t\tif n == 0 {
        \t\t\treturn false
        \t\t}
        \t\tmc.InitExpectedSize(mc.Size() + n)
        \t}
        \tchanged := false
        \tfor it := NewIterCursor(src); it.HasNext(); {
        \t\tif mc.AddKey(it.Next()) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func mapCoreAddAllNonNull(mc *JMapCore, src any) {
        \tCheckNotNull(src)
        \tif tv, ok := src.(jTypedViewCursorSource); ok {
        \t\tif n, ok2 := tv.typedViewLen_(); ok2 {
        \t\t\tif n == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tmc.InitExpectedSize(mc.Size() + int32(n))
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\tmc.Put(CheckNotNull(ToAny(tv.typedViewAt_(i))), true)
        \t\t\t}
        \t\t\treturn
        \t\t}
        \t}
        \tif si, ok := src.(SliceIterable); ok {
        \t\titems := si.IterSlice_()
        \t\tif len(items) == 0 {
        \t\t\treturn
        \t\t}
        \t\tmc.InitExpectedSize(mc.Size() + int32(len(items)))
        \t\tfor _, x := range items {
        \t\t\tmc.Put(CheckNotNull(ToAny(x)), true)
        \t\t}
        \t\treturn
        \t}
        \tif jsp, ok := src.(jFastJSetCoreProvider); ok {
        \t\tif sc := jsp.fastSetJCore_(); sc != nil {
        \t\t\tif sc.Size() == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tmc.InitExpectedSize(mc.Size() + sc.Size())
        \t\t\tsc.ForEach(func(k any) bool {
        \t\t\t\tmc.Put(CheckNotNull(k), true)
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn
        \t\t}
        \t}
        \tif sp, ok := src.(jFastSetCoreProvider); ok {
        \t\tif smc := sp.fastSetCore_(); smc != nil {
        \t\t\tif smc.Size() == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tmc.InitExpectedSize(mc.Size() + smc.Size())
        \t\t\tsmc.ForEachEntry(func(k, _ any) bool {
        \t\t\t\tmc.Put(CheckNotNull(k), true)
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn
        \t\t}
        \t}
        \tif ep, ok := src.(jFastEnumSetProvider); ok {
        \t\tif es := ep.fastEnumSet_(); es != nil {
        \t\t\tif es.size == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tmc.InitExpectedSize(mc.Size() + es.size)
        \t\t\tes.ForEach(func(k any) bool {
        \t\t\t\tmc.Put(k, true)
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn
        \t\t}
        \t}
        \tif c, ok := src.(Ext_java_util_Collection); ok {
        \t\tmc.InitExpectedSize(mc.Size() + c.Size())
        \t}
        \tfor it := NewIterCursor(src); it.HasNext(); {
        \t\tmc.Put(CheckNotNull(it.Next()), true)
        \t}
        }

        func mapCoreRemoveAllKeys(mc *JMapCore, other Ext_java_util_Collection) bool {
        \tCheckNotNull(other)
        \tif mc.immutable {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tif mc.Size() == 0 {
        \t\treturn false
        \t}
        \tif tv, ok := any(other).(jTypedViewCursorSource); ok {
        \t\tif n, ok2 := tv.typedViewLen_(); ok2 {
        \t\t\tchanged := false
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\tif mc.RemoveKey(tv.typedViewAt_(i)) {
        \t\t\t\t\tchanged = true
        \t\t\t\t}
        \t\t\t}
        \t\t\treturn changed
        \t\t}
        \t}
        \tif si, ok := any(other).(SliceIterable); ok {
        \t\tchanged := false
        \t\tfor _, x := range si.IterSlice_() {
        \t\t\tif mc.RemoveKey(x) {
        \t\t\t\tchanged = true
        \t\t\t}
        \t\t}
        \t\treturn changed
        \t}
        \tif jsp, ok := any(other).(jFastJSetCoreProvider); ok {
        \t\tif sc := jsp.fastSetJCore_(); sc != nil {
        \t\t\tchanged := false
        \t\t\tsc.ForEach(func(k any) bool {
        \t\t\t\tif mc.RemoveKey(k) {
        \t\t\t\t\tchanged = true
        \t\t\t\t}
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn changed
        \t\t}
        \t}
        \tif sp, ok := any(other).(jFastSetCoreProvider); ok {
        \t\tif omc := sp.fastSetCore_(); omc != nil && omc != mc {
        \t\t\tchanged := false
        \t\t\tomc.ForEachEntry(func(k, _ any) bool {
        \t\t\t\tif mc.RemoveKey(k) {
        \t\t\t\t\tchanged = true
        \t\t\t\t}
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn changed
        \t\t}
        \t}
        \tif ep, ok := any(other).(jFastEnumSetProvider); ok {
        \t\tif es := ep.fastEnumSet_(); es != nil {
        \t\t\tchanged := false
        \t\t\tes.ForEach(func(k any) bool {
        \t\t\t\tif mc.RemoveKey(k) {
        \t\t\t\t\tchanged = true
        \t\t\t\t}
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn changed
        \t\t}
        \t}
        \tchanged := false
        \tfor it := Iterate(other); it.HasNext(); {
        \t\tif mc.RemoveKey(it.Next()) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func mapCoreRetainAllKeys(mc *JMapCore, other Ext_java_util_Collection) bool {
        \tCheckNotNull(other)
        \tif mc.immutable {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tchanged := false
        \tfor _, k := range mc.LiveKeys() {
        \t\tif !other.Contains(k) && mc.RemoveKey(k) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func mapCoreRemoveIfKeys(mc *JMapCore, pred func(any) bool) bool {
        \tif mc.immutable {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tchanged := false
        \tfor _, k := range mc.LiveKeys() {
        \t\tif pred(k) && mc.RemoveKey(k) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func copyHashCoreInto(dst *JMapCore, src *JMapCore) *JMapCore {
        \tdst.Identity = src.Identity
        \tdst.genNonCustom = src.Identity
        \tdst.accessOrd = src.accessOrd
        \tif dst.CopyFromSameMode(src) {
        \t\treturn dst
        \t}
        \tdst.InitExpectedSize(src.Size())
        \tsrc.ForEachEntry(func(k, v any) bool {
        \t\tdst.Put(k, v)
        \t\treturn true
        \t})
        \treturn dst
        }

        func copySortedCoreInto(dst *JMapCore, src *JMapCore) *JMapCore {
        \tif src.IsNaturalOrder() {
        \t\tinitSortedCore(dst, nil)
        \t} else {
        \t\tinitSortedCore(dst, src.SortedCmp())
        \t}
        \tsrc.ForEachEntry(func(k, v any) bool {
        \t\tdst.Put(k, v)
        \t\treturn true
        \t})
        \treturn dst
        }

        func cloneJMapCore(src *JMapCore) *JMapCore {
        \tvar dst *JMapCore
        \tif src.IsNaturalOrder() {
        \t\tdst = NewSortedJMapCore(nil)
        \t} else if src.IsSorted() {
        \t\tdst = NewSortedJMapCore(src.SortedCmp())
        \t} else {
        \t\tdst = NewJMapCore(src.Identity)
        \t\tdst.accessOrd = src.accessOrd
        \t\tif dst.CopyFromSameMode(src) {
        \t\t\treturn dst
        \t\t}
        \t\tdst.InitExpectedSize(src.Size())
        \t}
        \tsrc.ForEachEntry(func(k, v any) bool {
        \t\tdst.Put(k, v)
        \t\treturn true
        \t})
        \treturn dst
        }

        func copyExtMapInto(dst *JMapCore, src Ext_java_util_Map) {
        \tCheckNotNull(src)
        \tif mp, ok := any(src).(jFastMapCoreProvider); ok {
        \t\tif smc := mp.fastMapCore_(); smc != nil {
        \t\t\tif smc == dst || smc.Size() == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tif dst.Size() == 0 && dst.CopyFromSameMode(smc) {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tdst.InitExpectedSize(dst.Size() + smc.Size())
        \t\t\tsmc.ForEachEntry(func(k, v any) bool {
        \t\t\t\tdst.Put(k, v)
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn
        \t\t}
        \t}
        \tif ep, ok := any(src).(jFastEnumMapProvider); ok {
        \t\tif em := ep.fastEnumMap_(); em != nil {
        \t\t\tif em.size == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tdst.InitExpectedSize(dst.Size() + em.size)
        \t\t\tem.ForEachEntry(func(k, v any) bool {
        \t\t\t\tdst.Put(k, v)
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn
        \t\t}
        \t}
        \tif es := src.EntrySet(); es != nil && !javalang.IsNil(es) {
        \t\tfor it := Iterate(es); it.HasNext(); {
        \t\t\te := Cast[Ext_java_util_Map_Entry](it.Next())
        \t\t\tdst.Put(e.GetKey(), e.GetValue())
        \t\t}
        \t}
        }

        func copyExtMapIntoNonNull(dst *JMapCore, src Ext_java_util_Map) {
        \tCheckNotNull(src)
        \tif mp, ok := any(src).(jFastMapCoreProvider); ok {
        \t\tif smc := mp.fastMapCore_(); smc != nil {
        \t\t\tif smc.Size() == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tdst.InitExpectedSize(dst.Size() + smc.Size())
        \t\t\tsmc.ForEachEntry(func(k, v any) bool {
        \t\t\t\tdst.Put(CheckNotNull(ToAny(k)), CheckNotNull(ToAny(v)))
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn
        \t\t}
        \t}
        \tif ep, ok := any(src).(jFastEnumMapProvider); ok {
        \t\tif em := ep.fastEnumMap_(); em != nil {
        \t\t\tif em.size == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tdst.InitExpectedSize(dst.Size() + em.size)
        \t\t\tem.ForEachEntry(func(k, v any) bool {
        \t\t\t\tdst.Put(CheckNotNull(ToAny(k)), CheckNotNull(ToAny(v)))
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn
        \t\t}
        \t}
        \tif es := src.EntrySet(); es != nil && !javalang.IsNil(es) {
        \t\tfor it := Iterate(es); it.HasNext(); {
        \t\t\te := Cast[Ext_java_util_Map_Entry](it.Next())
        \t\t\tdst.Put(CheckNotNull(ToAny(e.GetKey())), CheckNotNull(ToAny(e.GetValue())))
        \t\t}
        \t}
        }

        func populateImmSortedSetElems(obj *Ext_com_google_common_collect_ImmutableSortedSetBase, mc *JMapCore) {
        \tn := mc.Size()
        \tif n == 0 {
        \t\tobj.elems_ = emptyAnySlice_
        \t} else if n == 1 {
        \t\tobj.singleArr_[0] = mc.FirstKey()
        \t\tobj.elems_ = obj.singleArr_[:1:1]
        \t} else {
        \t\telems := mc.LiveKeys()
        \t\tobj.elems_ = elems[:len(elems):len(elems)]
        \t}
        }

        func setCoreDescending(mc *JMapCore, imm bool) any {
        \torigCmp := mc.SortedCmp()
        \tif origCmp == nil {
        \t\torigCmp = CompareValues
        \t}
        \trevCmp := func(a, b any) int32 { return origCmp(b, a) }
        \tkeys := mc.LiveKeys()
        \tif imm {
        \t\tout := NewSortedJMapCore(revCmp)
        \t\tfor i := len(keys) - 1; i >= 0; i-- {
        \t\t\tout.Put(keys[i], true)
        \t\t}
        \t\tout.sealImmutable()
        \t\tobj := &Ext_com_google_common_collect_ImmutableSortedSetBase{sortedCore_: out}
        \t\tobj.BindExtSelf_(obj)
        \t\tpopulateImmSortedSetElems(obj, out)
        \t\treturn obj
        \t}
        \tts := &Ext_java_util_TreeSet{}
        \tts.BindExtSelf_(ts)
        \tout := initSortedCore(&ts.inlineCore, revCmp)
        \tfor i := len(keys) - 1; i >= 0; i-- {
        \t\tout.Put(keys[i], true)
        \t}
        \treturn ts
        }

        func setCoreTailSet(mc *JMapCore, fromElem any, inclusive bool, imm bool) any {
        \tcmp := mc.SortedCmp()
        \tnat := mc.IsNaturalOrder()
        \tif cmp == nil {
        \t\tcmp = CompareValues
        \t}
        \tf := ToAny(fromElem)
        \tif imm {
        \t\tvar out *JMapCore
        \t\tif nat {
        \t\t\tout = NewSortedJMapCore(nil)
        \t\t} else {
        \t\t\tout = NewSortedJMapCore(cmp)
        \t\t}
        \t\tfor _, k := range mc.LiveKeys() {
        \t\t\tc := cmp(k, f)
        \t\t\tif c > 0 || (inclusive && c == 0) {
        \t\t\t\tout.Put(k, true)
        \t\t\t}
        \t\t}
        \t\treturn wrapImmutableSortedSet(out)
        \t}
        \tts := &Ext_java_util_TreeSet{}
        \tts.BindExtSelf_(ts)
        \tvar out *JMapCore
        \tif nat {
        \t\tout = initSortedCore(&ts.inlineCore, nil)
        \t} else {
        \t\tout = initSortedCore(&ts.inlineCore, cmp)
        \t}
        \tfor _, k := range mc.LiveKeys() {
        \t\tc := cmp(k, f)
        \t\tif c > 0 || (inclusive && c == 0) {
        \t\t\tout.Put(k, true)
        \t\t}
        \t}
        \treturn ts
        }

        // --- Iterators (ListIterator, Descending, MapKeyIter, SubList, Arrays.asList) ---

        type jFullListIter struct {
        \tlc      *JListCore
        \titems   []any
        \timm     bool
        \tcursor  int32
        \tlastRet int32
        }

        func newFullListIter(lc *JListCore, index int32) Ext_java_util_ListIterator {
        \tCheckPositionIndex(index, lc.Size())
        \treturn &jFullListIter{lc: lc, cursor: index, lastRet: -1}
        }

        func newImmListIter(items []any, index int32) Ext_java_util_ListIterator {
        \tCheckPositionIndex(index, int32(len(items)))
        \treturn &jFullListIter{items: items, imm: true, cursor: index, lastRet: -1}
        }

        func (*jFullListIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (*jFullListIter) Implements_j2gort_Ext_java_util_ListIterator() {}
        func (*jFullListIter) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator() {}
        func (*jFullListIter) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableListIterator() {}
        func (*jFullListIter) GetExt_com_google_common_collect_UnmodifiableIteratorBase() *Ext_com_google_common_collect_UnmodifiableIteratorBase {
        \treturn nil
        }
        func (*jFullListIter) GetExt_com_google_common_collect_UnmodifiableListIteratorBase() *Ext_com_google_common_collect_UnmodifiableListIteratorBase {
        \treturn nil
        }

        func (it *jFullListIter) listSize() int32 {
        \tif it.imm {
        \t\treturn int32(len(it.items))
        \t}
        \treturn it.lc.Size()
        }
        func (it *jFullListIter) HasNext() bool { return it.cursor < it.listSize() }
        func (it *jFullListIter) Next() any {
        \tif it.cursor >= it.listSize() {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tvar v any
        \tif it.imm {
        \t\tv = it.items[it.cursor]
        \t} else {
        \t\tv = it.lc.Get(it.cursor)
        \t}
        \tit.lastRet = it.cursor
        \tit.cursor++
        \treturn v
        }
        func (it *jFullListIter) HasPrevious() bool { return it.cursor > 0 }
        func (it *jFullListIter) Previous() any {
        \tif it.cursor <= 0 {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tit.cursor--
        \tit.lastRet = it.cursor
        \tif it.imm {
        \t\treturn it.items[it.cursor]
        \t}
        \treturn it.lc.Get(it.cursor)
        }
        func (it *jFullListIter) NextIndex() int32     { return it.cursor }
        func (it *jFullListIter) PreviousIndex() int32 { return it.cursor - 1 }
        func (it *jFullListIter) Remove() {
        \tif it.imm {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tif it.lastRet < 0 {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \tit.lc.RemoveAt(it.lastRet)
        \tif it.lastRet < it.cursor {
        \t\tit.cursor--
        \t}
        \tit.lastRet = -1
        }
        func (it *jFullListIter) Set(x any) {
        \tif it.imm {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tif it.lastRet < 0 {
        \t\tpanic(javalang.NewIllegalStateException("set before next"))
        \t}
        \tit.lc.Set(it.lastRet, x)
        }
        func (it *jFullListIter) Add(x any) {
        \tif it.imm {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tit.lc.AddAt(it.cursor, x)
        \tit.cursor++
        \tit.lastRet = -1
        }

        type jDescListIter struct {
        \tlc      *JListCore
        \tcursor  int32
        \tlastRet int32
        }

        func newDescendingListIter(lc *JListCore) Ext_java_util_Iterator {
        \treturn &jDescListIter{lc: lc, cursor: lc.Size() - 1, lastRet: -1}
        }

        func (*jDescListIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jDescListIter) HasNext() bool { return it.cursor >= 0 }
        func (it *jDescListIter) Next() any {
        \tif it.cursor < 0 {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tv := it.lc.Get(it.cursor)
        \tit.lastRet = it.cursor
        \tit.cursor--
        \treturn v
        }
        func (it *jDescListIter) Remove() {
        \tif it.lastRet < 0 {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \tit.lc.RemoveAt(it.lastRet)
        \tit.lastRet = -1
        }

        type jDequeIter struct {
        \tdc      *JDequeCore
        \toff     int32
        \tlastOff int32
        }

        func newDequeIter(dc *JDequeCore) Ext_java_util_Iterator {
        \tif dc == nil || dc.size == 0 {
        \t\treturn emptyMapIter_
        \t}
        \treturn &jDequeIter{dc: dc, off: 0, lastOff: -1}
        }

        func (*jDequeIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jDequeIter) HasNext() bool { return it.off < it.dc.size }
        func (it *jDequeIter) Next() any {
        \tif it.off >= it.dc.size {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tv := it.dc.elemAt((it.dc.head + it.off) & it.dc.mask)
        \tit.lastOff = it.off
        \tit.off++
        \treturn v
        }
        func (it *jDequeIter) Remove() {
        \tif it.lastOff < 0 {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \tit.dc.deleteAtOffset(it.lastOff)
        \tit.off = it.lastOff
        \tit.lastOff = -1
        }

        type jDequeDescIter struct {
        \tdc      *JDequeCore
        \toff     int32
        \tlastOff int32
        }

        func newDequeDescIter(dc *JDequeCore) Ext_java_util_Iterator {
        \tif dc == nil || dc.size == 0 {
        \t\treturn emptyMapIter_
        \t}
        \treturn &jDequeDescIter{dc: dc, off: dc.size - 1, lastOff: -1}
        }

        func (*jDequeDescIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jDequeDescIter) HasNext() bool { return it.off >= 0 && it.off < it.dc.size }
        func (it *jDequeDescIter) Next() any {
        \tif it.off < 0 || it.off >= it.dc.size {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tv := it.dc.elemAt((it.dc.head + it.off) & it.dc.mask)
        \tit.lastOff = it.off
        \tit.off--
        \treturn v
        }
        func (it *jDequeDescIter) Remove() {
        \tif it.lastOff < 0 {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \tit.dc.deleteAtOffset(it.lastOff)
        \tit.lastOff = -1
        }

        type jEmptyMapIter struct{}

        var emptyMapIter_ = &jEmptyMapIter{}

        func (*jEmptyMapIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (*jEmptyMapIter) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator() {}
        func (*jEmptyMapIter) GetExt_com_google_common_collect_UnmodifiableIteratorBase() *Ext_com_google_common_collect_UnmodifiableIteratorBase {
        \treturn nil
        }
        func (*jEmptyMapIter) HasNext() bool { return false }
        func (*jEmptyMapIter) Next() any {
        \tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        }
        func (*jEmptyMapIter) Remove() {
        \tpanic(javalang.NewIllegalStateException("remove before next"))
        }

        type jMapCursor struct {
        \tmc       *JMapCore
        \tmode     uint8
        \tidx      int
        \tcurPos   int32
        \tlastSlot int32
        \tlastKey  any
        \thasLast  bool
        \tkeys     []any
        \tstrKeys  []string
        \tvals     []any
        }

        func initMapCursor(mc *JMapCore, needVals bool) (jMapCursor, bool) {
        \tif mc == nil || mc.Live == 0 {
        \t\treturn jMapCursor{}, false
        \t}
        \tswitch mc.mode {
        \tcase jMapModeNullOnly, jMapModeNode1:
        \t\treturn jMapCursor{mc: mc, mode: mc.mode, idx: 0, lastSlot: -1}, true
        \tcase jMapModeStr, jMapModeNode, jMapModePtr:
        \t\tstart := int(mc.firstIdx)
        \t\tc := jMapCursor{mc: mc, mode: mc.mode, idx: start, lastSlot: -1}
        \t\tc.advanceSlice()
        \t\tif c.idx != start && !mc.immutable {
        \t\t\tmc.firstIdx = int32(c.idx)
        \t\t}
        \t\treturn c, c.hasSliceNext()
        \tcase jMapModeGeneral:
        \t\tg := mc.ext.(*jGenMapState)
        \t\treturn jMapCursor{mc: mc, mode: jMapModeGeneral, curPos: g.head, lastSlot: -1}, g.head != 0
        \tcase jMapModeSortedNat, jMapModeSortedCmp:
        \t\tkeys := make([]any, 0, mc.Live)
        \t\tvar vals []any
        \t\tif needVals {
        \t\t\tvals = make([]any, 0, mc.Live)
        \t\t}
        \t\tmc.ForEachEntry(func(k, v any) bool {
        \t\t\tkeys = append(keys, k)
        \t\t\tif needVals {
        \t\t\t\tvals = append(vals, v)
        \t\t\t}
        \t\t\treturn true
        \t\t})
        \t\treturn jMapCursor{mc: mc, mode: mc.mode, lastSlot: -1, keys: keys, vals: vals}, len(keys) > 0
        \tdefault:
        \t\treturn jMapCursor{}, false
        \t}
        }

        func initMapValCursor(mc *JMapCore) (jMapCursor, bool) {
        \tif mc == nil || mc.Live == 0 {
        \t\treturn jMapCursor{}, false
        \t}
        \tif mc.mode == jMapModeSortedNat || mc.mode == jMapModeSortedCmp {
        \t\tvals := make([]any, 0, mc.Live)
        \t\tmc.ForEachEntry(func(_, v any) bool {
        \t\t\tvals = append(vals, v)
        \t\t\treturn true
        \t\t})
        \t\treturn jMapCursor{mc: mc, mode: mc.mode, lastSlot: -1, vals: vals}, len(vals) > 0
        \t}
        \treturn initMapCursor(mc, false)
        }

        func (c *jMapCursor) hasSliceNext() bool {
        \tswitch c.mode {
        \tcase jMapModeStr:
        \t\treturn c.idx < len(c.mc.strEntries)
        \tcase jMapModeNode:
        \t\treturn c.idx < len(c.mc.ext.(*jNodeMapState).entries)
        \tcase jMapModePtr:
        \t\treturn c.idx < len(c.mc.ext.(*jPtrMapState).entries)
        \tdefault:
        \t\treturn false
        \t}
        }

        func (c *jMapCursor) advanceSlice() {
        \tswitch c.mode {
        \tcase jMapModeStr:
        \t\tfor c.idx < len(c.mc.strEntries) && c.mc.strEntries[c.idx].key == jTombstoneStr {
        \t\t\tc.idx++
        \t\t}
        \tcase jMapModeNode:
        \t\tes := c.mc.ext.(*jNodeMapState).entries
        \t\tfor c.idx < len(es) && !es[c.idx].active {
        \t\t\tc.idx++
        \t\t}
        \tcase jMapModePtr:
        \t\tes := c.mc.ext.(*jPtrMapState).entries
        \t\tfor c.idx < len(es) && es[c.idx].key == jTombstoneObj {
        \t\t\tc.idx++
        \t\t}
        \t}
        }

        func (c *jMapCursor) hasNext() bool {
        \tswitch c.mode {
        \tcase jMapModeNullOnly:
        \t\treturn c.idx == 0 && c.mc.mode == jMapModeNullOnly
        \tcase jMapModeNode1:
        \t\treturn c.idx == 0 && c.mc.mode == jMapModeNode1
        \tcase jMapModeGeneral:
        \t\treturn c.curPos != 0
        \tcase jMapModeSortedNat, jMapModeSortedCmp:
        \t\treturn c.idx < len(c.keys) || c.idx < len(c.strKeys)
        \tdefault:
        \t\tc.advanceSlice()
        \t\treturn c.hasSliceNext()
        \t}
        }

        func (c *jMapCursor) hasValNext() bool {
        \tswitch c.mode {
        \tcase jMapModeNullOnly:
        \t\treturn c.idx == 0 && c.mc.mode == jMapModeNullOnly
        \tcase jMapModeNode1:
        \t\treturn c.idx == 0 && c.mc.mode == jMapModeNode1
        \tcase jMapModeGeneral:
        \t\treturn c.curPos != 0
        \tcase jMapModeSortedNat, jMapModeSortedCmp:
        \t\treturn c.idx < len(c.vals)
        \tdefault:
        \t\tc.advanceSlice()
        \t\treturn c.hasSliceNext()
        \t}
        }

        func (c *jMapCursor) nextVal() any {
        \tswitch c.mode {
        \tcase jMapModeNullOnly:
        \t\tif c.idx > 0 || c.mc.mode != jMapModeNullOnly {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tc.idx = 1
        \t\treturn c.mc.ext
        \tcase jMapModeNode1:
        \t\tif c.idx > 0 || c.mc.mode != jMapModeNode1 {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tc.idx = 1
        \t\treturn c.mc.ext
        \tcase jMapModeStr:
        \t\tc.advanceSlice()
        \t\tif c.idx >= len(c.mc.strEntries) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tv := c.mc.strEntries[c.idx].val
        \t\tc.idx++
        \t\tc.advanceSlice()
        \t\treturn v
        \tcase jMapModeNode:
        \t\tns := c.mc.ext.(*jNodeMapState)
        \t\tc.advanceSlice()
        \t\tif c.idx >= len(ns.entries) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tv := ns.entries[c.idx].val
        \t\tc.idx++
        \t\tc.advanceSlice()
        \t\treturn v
        \tcase jMapModePtr:
        \t\tps := c.mc.ext.(*jPtrMapState)
        \t\tc.advanceSlice()
        \t\tif c.idx >= len(ps.entries) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tv := ps.entries[c.idx].val
        \t\tc.idx++
        \t\tc.advanceSlice()
        \t\treturn v
        \tcase jMapModeGeneral:
        \t\tif c.curPos == 0 {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tg := c.mc.ext.(*jGenMapState)
        \t\tpos := c.curPos
        \t\tv := g.entries[pos-1].Val
        \t\tc.curPos = g.entries[pos-1].next()
        \t\treturn v
        \tdefault:
        \t\tif c.idx >= len(c.vals) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tv := c.vals[c.idx]
        \t\tc.idx++
        \t\treturn v
        \t}
        }

        func (c *jMapCursor) nextKey() any {
        \tswitch c.mode {
        \tcase jMapModeNullOnly:
        \t\tif c.idx > 0 || c.mc.mode != jMapModeNullOnly {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tc.idx = 1
        \t\treturn nil
        \tcase jMapModeNode1:
        \t\tif c.idx > 0 || c.mc.mode != jMapModeNode1 {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tc.idx = 1
        \t\treturn boxNodeById(uint32(c.mc.firstIdx))
        \tcase jMapModeStr:
        \t\tc.advanceSlice()
        \t\tif c.idx >= len(c.mc.strEntries) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tek := c.mc.strEntries[c.idx].key
        \t\tc.idx++
        \t\tc.advanceSlice()
        \t\treturn BoxString(ek)
        \tcase jMapModeNode:
        \t\tns := c.mc.ext.(*jNodeMapState)
        \t\tc.advanceSlice()
        \t\tif c.idx >= len(ns.entries) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tek := ns.entries[c.idx].key
        \t\tc.idx++
        \t\tc.advanceSlice()
        \t\treturn boxNodeById(ek)
        \tcase jMapModePtr:
        \t\tps := c.mc.ext.(*jPtrMapState)
        \t\tc.advanceSlice()
        \t\tif c.idx >= len(ps.entries) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tek := ps.entries[c.idx].key
        \t\tc.idx++
        \t\tc.advanceSlice()
        \t\treturn ek
        \tcase jMapModeGeneral:
        \t\tif c.curPos == 0 {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tg := c.mc.ext.(*jGenMapState)
        \t\tpos := c.curPos
        \t\tk := g.entryKey(pos)
        \t\tc.curPos = g.entries[pos-1].next()
        \t\treturn k
        \tdefault:
        \t\tif c.idx >= len(c.keys) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tk := c.keys[c.idx]
        \t\tc.idx++
        \t\treturn k
        \t}
        }

        func avlAppendStrKeys(s *jSortedMapState, cur int32, out []string) []string {
        \tif cur == 0 {
        \t\treturn out
        \t}
        \tn := &s.tree[cur-1]
        \tif n.left != 0 {
        \t\tout = avlAppendStrKeys(s, n.left, out)
        \t}
        \tif n.isStr {
        \t\tout = append(out, n.strKey)
        \t} else {
        \t\tout = append(out, CastString(n.key))
        \t}
        \tif n.right != 0 {
        \t\tout = avlAppendStrKeys(s, n.right, out)
        \t}
        \treturn out
        }

        func initMapKeyStrCursor(mc *JMapCore) (jMapCursor, bool) {
        \tif mc == nil || mc.Live == 0 {
        \t\treturn jMapCursor{}, false
        \t}
        \tif mc.mode == jMapModeSortedNat || mc.mode == jMapModeSortedCmp {
        \t\tst := mc.ext.(*jSortedMapState)
        \t\tstrKeys := avlAppendStrKeys(st, st.root, make([]string, 0, mc.Live))
        \t\treturn jMapCursor{mc: mc, mode: mc.mode, lastSlot: -1, strKeys: strKeys}, len(strKeys) > 0
        \t}
        \treturn initMapCursor(mc, false)
        }

        func (c *jMapCursor) nextKeyStr() string {
        \tswitch c.mode {
        \tcase jMapModeNullOnly:
        \t\tif c.idx > 0 || c.mc.mode != jMapModeNullOnly {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tc.idx = 1
        \t\treturn NullString
        \tcase jMapModeStr:
        \t\tc.advanceSlice()
        \t\tif c.idx >= len(c.mc.strEntries) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tek := c.mc.strEntries[c.idx].key
        \t\tc.idx++
        \t\tc.advanceSlice()
        \t\treturn ek
        \tcase jMapModeGeneral:
        \t\tif c.curPos == 0 {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tg := c.mc.ext.(*jGenMapState)
        \t\tpos := c.curPos
        \t\te := &g.entries[pos-1]
        \t\tc.curPos = e.next()
        \t\tif e.kind() == jKeyNull {
        \t\t\treturn NullString
        \t\t}
        \t\tif s, ok := e.Key.(string); ok {
        \t\t\treturn s
        \t\t}
        \t\treturn CastString(g.entryKey(pos))
        \tdefault:
        \t\tif c.strKeys != nil {
        \t\t\tif c.idx >= len(c.strKeys) {
        \t\t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t\t}
        \t\t\ts := c.strKeys[c.idx]
        \t\t\tc.idx++
        \t\t\treturn s
        \t\t}
        \t\treturn CastString(c.nextKey())
        \t}
        }

        func (c *jMapCursor) nextKeyNode() uint32 {
        \tswitch c.mode {
        \tcase jMapModeNullOnly:
        \t\tif c.idx > 0 || c.mc.mode != jMapModeNullOnly {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tc.idx = 1
        \t\treturn 0
        \tcase jMapModeNode1:
        \t\tif c.idx > 0 || c.mc.mode != jMapModeNode1 {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tc.idx = 1
        \t\treturn uint32(c.mc.firstIdx)
        \tcase jMapModeNode:
        \t\tns := c.mc.ext.(*jNodeMapState)
        \t\tc.advanceSlice()
        \t\tif c.idx >= len(ns.entries) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tu := ns.entries[c.idx].key
        \t\tc.idx++
        \t\tc.advanceSlice()
        \t\treturn u
        \tcase jMapModeGeneral:
        \t\tif c.curPos == 0 {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tg := c.mc.ext.(*jGenMapState)
        \t\tpos := c.curPos
        \t\te := &g.entries[pos-1]
        \t\tc.curPos = e.next()
        \t\tif e.kind() == jKeyNull {
        \t\t\treturn 0
        \t\t}
        \t\tif e.kind() == jKeyNode {
        \t\t\treturn e.u32OrHash
        \t\t}
        \t\treturn CastNode[uint32](g.entryKey(pos))
        \tdefault:
        \t\treturn CastNode[uint32](c.nextKey())
        \t}
        }

        func (c *jMapCursor) nextKV() (k any, v any, slot int32) {
        \tswitch c.mode {
        \tcase jMapModeNullOnly:
        \t\tif c.idx > 0 || c.mc.mode != jMapModeNullOnly {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tk = nil
        \t\tv = c.mc.ext
        \t\tslot = 0
        \t\tc.idx = 1
        \t\tc.lastSlot = 0
        \t\tc.lastKey = nil
        \t\tc.hasLast = true
        \t\treturn
        \tcase jMapModeNode1:
        \t\tif c.idx > 0 || c.mc.mode != jMapModeNode1 {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tk = boxNodeById(uint32(c.mc.firstIdx))
        \t\tv = c.mc.ext
        \t\tslot = 0
        \t\tc.idx = 1
        \t\tc.lastSlot = 0
        \t\tc.lastKey = k
        \t\tc.hasLast = true
        \t\treturn
        \tcase jMapModeStr:
        \t\tc.advanceSlice()
        \t\tif c.idx >= len(c.mc.strEntries) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\te := c.mc.strEntries[c.idx]
        \t\tif e.key != NullString {
        \t\t\tk = e.key
        \t\t}
        \t\tv = e.val
        \t\tslot = int32(c.idx)
        \t\tc.idx++
        \t\tc.lastSlot = slot
        \t\tc.lastKey = k
        \t\tc.hasLast = true
        \t\tc.advanceSlice()
        \t\treturn
        \tcase jMapModeNode:
        \t\tns := c.mc.ext.(*jNodeMapState)
        \t\tc.advanceSlice()
        \t\tif c.idx >= len(ns.entries) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\te := ns.entries[c.idx]
        \t\tk = boxNodeById(e.key)
        \t\tv = e.val
        \t\tslot = int32(c.idx)
        \t\tc.idx++
        \t\tc.lastSlot = slot
        \t\tc.lastKey = k
        \t\tc.hasLast = true
        \t\tc.advanceSlice()
        \t\treturn
        \tcase jMapModePtr:
        \t\tps := c.mc.ext.(*jPtrMapState)
        \t\tc.advanceSlice()
        \t\tif c.idx >= len(ps.entries) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\te := ps.entries[c.idx]
        \t\tk = e.key
        \t\tv = e.val
        \t\tslot = int32(c.idx)
        \t\tc.idx++
        \t\tc.lastSlot = slot
        \t\tc.lastKey = k
        \t\tc.hasLast = true
        \t\tc.advanceSlice()
        \t\treturn
        \tcase jMapModeGeneral:
        \t\tif c.curPos == 0 {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tg := c.mc.ext.(*jGenMapState)
        \t\tpos := c.curPos
        \t\tc.lastSlot = pos
        \t\tc.hasLast = true
        \t\tk = g.entryKey(pos)
        \t\tv = g.entries[pos-1].Val
        \t\tslot = pos
        \t\tc.curPos = g.entries[pos-1].next()
        \t\treturn
        \tdefault:
        \t\tif c.idx >= len(c.keys) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tk = c.keys[c.idx]
        \t\tif c.vals != nil {
        \t\t\tv = c.vals[c.idx]
        \t\t}
        \t\tslot = -1
        \t\tc.idx++
        \t\tc.lastSlot = -1
        \t\tc.lastKey = k
        \t\tc.hasLast = true
        \t\treturn
        \t}
        }

        func (c *jMapCursor) remove() {
        \tif !c.hasLast {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \tif c.mc.immutable {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tc.hasLast = false
        \tif c.mode == jMapModeGeneral && c.mc.mode == jMapModeGeneral && c.lastSlot > 0 {
        \t\tc.mc.removeGenPos(c.lastSlot)
        \t\tc.lastSlot = -1
        \t\treturn
        \t}
        \tc.mc.Remove(c.lastKey)
        \tc.lastSlot = -1
        }

        type jMapKeyIter struct {
        \tcur jMapCursor
        }

        func newMapKeyIter(mc *JMapCore) Ext_java_util_Iterator {
        \tc, ok := initMapCursor(mc, false)
        \tif !ok {
        \t\treturn emptyMapIter_
        \t}
        \treturn &jMapKeyIter{cur: c}
        }

        func (*jMapKeyIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (*jMapKeyIter) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator() {}
        func (*jMapKeyIter) GetExt_com_google_common_collect_UnmodifiableIteratorBase() *Ext_com_google_common_collect_UnmodifiableIteratorBase {
        \treturn nil
        }

        func (it *jMapKeyIter) HasNext() bool { return it.cur.hasNext() }
        func (it *jMapKeyIter) Next() any {
        \tk, _, _ := it.cur.nextKV()
        \treturn k
        }
        func (it *jMapKeyIter) Remove() { it.cur.remove() }

        // Live SubList view
        type jSubListView struct {
        \tExt_java_util_AbstractListBase
        \tparent *JListCore
        \tfrom   int32
        \tto     int32
        }

        func newSubListView(parent *JListCore, from, to int32) Ext_java_util_List {
        \tCheckPositionIndexes(from, to, parent.Size())
        \tv := &jSubListView{parent: parent, from: from, to: to}
        \tv.BindExtSelf_(v)
        \treturn v
        }

        func (v *jSubListView) Size() int32 { return v.to - v.from }
        func (v *jSubListView) Get(i int32) any {
        \tCheckElementIndex(i, v.Size())
        \treturn v.parent.Get(v.from + i)
        }
        func (v *jSubListView) Set(i int32, x any) any {
        \tCheckElementIndex(i, v.Size())
        \treturn v.parent.Set(v.from+i, x)
        }
        func (v *jSubListView) Add_2(i int32, x any) {
        \tCheckPositionIndex(i, v.Size())
        \tv.parent.AddAt(v.from+i, x)
        \tv.to++
        }
        func (v *jSubListView) Remove_1_int32(i int32) any {
        \tCheckElementIndex(i, v.Size())
        \told := v.parent.RemoveAt(v.from + i)
        \tv.to--
        \treturn old
        }
        func (v *jSubListView) Clear() {
        \tif v.from == v.to {
        \t\treturn
        \t}
        \tCheckPositionIndexes(v.from, v.to, v.parent.Size())
        \ts := v.parent.Items
        \tnewLen := len(s) - int(v.to-v.from)
        \tcopy(s[v.from:], s[v.to:])
        \tclear(s[newLen:])
        \tv.parent.Items = s[:newLen]
        \tv.to = v.from
        }

        // Backed Arrays.asList view
        type jBackedArrayList struct {
        \tExt_java_util_AbstractListBase
        \tarr []any
        }

        func newBackedArrayList(arr []any) Ext_java_util_List {
        \tCheckNotNullSlice(arr)
        \tv := &jBackedArrayList{arr: arr}
        \tv.BindExtSelf_(v)
        \treturn v
        }

        func (v *jBackedArrayList) Size() int32 { return int32(len(v.arr)) }
        func (v *jBackedArrayList) Get(i int32) any {
        \tCheckElementIndex(i, int32(len(v.arr)))
        \treturn ToAny(v.arr[i])
        }
        func (v *jBackedArrayList) Set(i int32, x any) any {
        \tCheckElementIndex(i, int32(len(v.arr)))
        \told := ToAny(v.arr[i])
        \tv.arr[i] = ToAny(x)
        \treturn old
        }

        // Lists.reverse & Lists.transform live views
        type jReversedListView struct {
        \tExt_java_util_AbstractListBase
        \tbase Ext_java_util_List
        }

        func newReversedListView(base Ext_java_util_List) Ext_java_util_List {
        \tCheckNotNull(base)
        \tv := &jReversedListView{base: base}
        \tv.BindExtSelf_(v)
        \treturn v
        }

        func (v *jReversedListView) Size() int32 { return v.base.Size() }
        func (v *jReversedListView) Get(i int32) any {
        \tn := v.base.Size()
        \tCheckElementIndex(i, n)
        \treturn v.base.Get(n - 1 - i)
        }
        func (v *jReversedListView) Set(i int32, x any) any {
        \tn := v.base.Size()
        \tCheckElementIndex(i, n)
        \treturn v.base.Set(n-1-i, x)
        }
        func (v *jReversedListView) Add_2(i int32, x any) {
        \tn := v.base.Size()
        \tCheckPositionIndex(i, n)
        \tv.base.Add_2(n-i, x)
        }
        func (v *jReversedListView) Remove_1_int32(i int32) any {
        \tn := v.base.Size()
        \tCheckElementIndex(i, n)
        \treturn v.base.Remove_1_int32(n - 1 - i)
        }

        type jTransformedListView struct {
        \tExt_java_util_AbstractListBase
        \tbase Ext_java_util_List
        \tfn   Ext_com_google_common_base_Function
        }

        func newTransformedListView(base Ext_java_util_List, fn Ext_com_google_common_base_Function) Ext_java_util_List {
        \tCheckNotNull(base)
        \tCheckNotNull(fn)
        \tv := &jTransformedListView{base: base, fn: fn}
        \tv.BindExtSelf_(v)
        \treturn v
        }

        func (v *jTransformedListView) Size() int32     { return v.base.Size() }
        func (v *jTransformedListView) Get(i int32) any { return v.fn.Apply(v.base.Get(i)) }
        func (v *jTransformedListView) Remove_1_int32(i int32) any {
        \treturn v.fn.Apply(v.base.Remove_1_int32(i))
        }
        func (v *jTransformedListView) Clear() { v.base.Clear() }

        // --- Map Operations & Live Views (keySet, values, entrySet) ---

        func mapCoreRemove2(mc *JMapCore, k, v any) bool {
        \tif !mc.ContainsKey(k) || !javalang.ObjectsEquals(mc.Get(k), ToAny(v)) {
        \t\treturn false
        \t}
        \tmc.Remove(k)
        \treturn true
        }

        func mapCoreReplace2(mc *JMapCore, k, newVal any) any {
        \tif !mc.ContainsKey(k) {
        \t\treturn nil
        \t}
        \treturn mc.Put(k, newVal)
        }

        func mapCoreReplace3(mc *JMapCore, k, oldVal, newVal any) bool {
        \tif !mc.ContainsKey(k) || !javalang.ObjectsEquals(mc.Get(k), ToAny(oldVal)) {
        \t\treturn false
        \t}
        \tmc.Put(k, newVal)
        \treturn true
        }

        func mapCoreReplaceAll(mc *JMapCore, fn Ext_java_util_function_BiFunction) {
        \tCheckNotNull(fn)
        \tmc.TransformValues(func(k, v any) any {
        \t\treturn fn.Apply(k, v)
        \t})
        }

        func mapCoreForEach(mc *JMapCore, action Ext_java_util_function_BiConsumer) {
        \tCheckNotNull(action)
        \tmc.ForEachEntry(func(k, v any) bool {
        \t\taction.Accept(k, v)
        \t\treturn true
        \t})
        }

        func mapCoreComputeIfAbsent(mc *JMapCore, k any, fn Ext_java_util_function_Function) any {
        \tCheckNotNull(fn)
        \tif v := mc.Get(k); v != nil {
        \t\treturn v
        \t}
        \tnewVal := ToAny(fn.Apply(ToAny(k)))
        \tif newVal != nil {
        \t\tmc.Put(k, newVal)
        \t}
        \treturn newVal
        }

        func mapCoreComputeIfAbsentStr(mc *JMapCore, k string, fn Ext_java_util_function_Function) any {
        \tCheckNotNull(fn)
        \tif k == NullString && mc.IsNaturalOrder() {
        \t\tThrowNPE()
        \t}
        \tif v := mc.GetStr(k); v != nil {
        \t\treturn v
        \t}
        \tnewVal := ToAny(fn.Apply(BoxString(k)))
        \tif newVal != nil {
        \t\tmc.PutStr(k, newVal)
        \t}
        \treturn newVal
        }

        func mapCoreComputeIfAbsentNode(mc *JMapCore, u uint32, fn Ext_java_util_function_Function) any {
        \tCheckNotNull(fn)
        \tif u == 0 && mc.IsNaturalOrder() {
        \t\tThrowNPE()
        \t}
        \tif v := mc.GetNode(u); v != nil {
        \t\treturn v
        \t}
        \tnewVal := ToAny(fn.Apply(boxNodeById(u)))
        \tif newVal != nil {
        \t\tmc.PutNode(u, newVal)
        \t}
        \treturn newVal
        }

        func mapCoreTryGetComputeStr(mc *JMapCore, k string) (any, bool) {
        \tif k == NullString && mc.IsNaturalOrder() {
        \t\treturn nil, false
        \t}
        \tif v := mc.GetStr(k); v != nil {
        \t\treturn v, true
        \t}
        \treturn nil, false
        }

        func mapCoreTryGetComputeNode(mc *JMapCore, u uint32) (any, bool) {
        \tif u == 0 && mc.IsNaturalOrder() {
        \t\treturn nil, false
        \t}
        \tif v := mc.GetNode(u); v != nil {
        \t\treturn v, true
        \t}
        \treturn nil, false
        }

        func mapCoreTryGetComputeAny(mc *JMapCore, k any) (any, bool) {
        \tif (k == nil || isJavaNilAny(k)) && mc.IsNaturalOrder() {
        \t\treturn nil, false
        \t}
        \tif v := mc.Get(k); v != nil {
        \t\treturn v, true
        \t}
        \treturn nil, false
        }

        func mapCoreComputeIfPresent(mc *JMapCore, k any, fn Ext_java_util_function_BiFunction) any {
        \tCheckNotNull(fn)
        \toldVal := mc.Get(k)
        \tif oldVal == nil {
        \t\treturn nil
        \t}
        \tnewVal := ToAny(fn.Apply(ToAny(k), oldVal))
        \tif newVal == nil {
        \t\tmc.Remove(k)
        \t\treturn nil
        \t}
        \tmc.Put(k, newVal)
        \treturn newVal
        }

        func mapCoreCompute(mc *JMapCore, k any, fn Ext_java_util_function_BiFunction) any {
        \tCheckNotNull(fn)
        \toldVal := mc.Get(k)
        \thaveOld := mc.ContainsKey(k)
        \tnewVal := ToAny(fn.Apply(ToAny(k), oldVal))
        \tif newVal == nil {
        \t\tif haveOld {
        \t\t\tmc.Remove(k)
        \t\t}
        \t\treturn nil
        \t}
        \tmc.Put(k, newVal)
        \treturn newVal
        }

        func mapCoreMerge(mc *JMapCore, k, val any, fn Ext_java_util_function_BiFunction) any {
        \tCheckNotNull(fn)
        \tv := CheckNotNull(ToAny(val))
        \toldVal := mc.Get(k)
        \tvar newVal any
        \tif oldVal == nil {
        \t\tnewVal = v
        \t} else {
        \t\tnewVal = ToAny(fn.Apply(oldVal, v))
        \t}
        \tif newVal == nil {
        \t\tmc.Remove(k)
        \t} else {
        \t\tmc.Put(k, newVal)
        \t}
        \treturn newVal
        }

        func mapCoreEquals(mc *JMapCore, other any) bool {
        \tmp, ok := other.(Ext_java_util_Map)
        \tif !ok || mp == nil || isJavaNilAny(mp) {
        \t\treturn false
        \t}
        \tif mc.Size() != mp.Size() {
        \t\treturn false
        \t}
        \tif mc.Size() == 0 {
        \t\treturn true
        \t}
        \tif fmp, ok2 := other.(jFastMapCoreProvider); ok2 {
        \t\tif omc := fmp.fastMapCore_(); omc != nil {
        \t\t\tif mc == omc {
        \t\t\t\treturn true
        \t\t\t}
        \t\t\tif mc.mode == jMapModeStr && omc.mode == jMapModeStr {
        \t\t\t\tfor i := int(mc.firstIdx); i < len(mc.strEntries); i++ {
        \t\t\t\t\te := mc.strEntries[i]
        \t\t\t\t\tif e.key == jTombstoneStr {
        \t\t\t\t\t\tcontinue
        \t\t\t\t\t}
        \t\t\t\t\tv2, ok3 := omc.getStr(e.key, false)
        \t\t\t\t\tif !ok3 || !javalang.ObjectsEquals(e.val, v2) {
        \t\t\t\t\t\treturn false
        \t\t\t\t\t}
        \t\t\t\t}
        \t\t\t\treturn true
        \t\t\t}
        \t\t\tif mc.mode == jMapModeNode1 {
        \t\t\t\tv2, ok3 := omc.getNode(uint32(mc.firstIdx), false)
        \t\t\t\treturn ok3 && javalang.ObjectsEquals(mc.ext, v2)
        \t\t\t}
        \t\t\tif mc.mode == jMapModeNode && omc.mode == jMapModeNode {
        \t\t\t\tns := mc.ext.(*jNodeMapState)
        \t\t\t\tfor i := int(mc.firstIdx); i < len(ns.entries); i++ {
        \t\t\t\t\te := ns.entries[i]
        \t\t\t\t\tif !e.active {
        \t\t\t\t\t\tcontinue
        \t\t\t\t\t}
        \t\t\t\t\tv2, ok3 := omc.getNode(e.key, false)
        \t\t\t\t\tif !ok3 || !javalang.ObjectsEquals(e.val, v2) {
        \t\t\t\t\t\treturn false
        \t\t\t\t\t}
        \t\t\t\t}
        \t\t\t\treturn true
        \t\t\t}
        \t\t\teq := true
        \t\t\tmc.ForEachEntry(func(k, v any) bool {
        \t\t\t\tv2, ok3 := omc.GetEntry(k)
        \t\t\t\tif !ok3 || !javalang.ObjectsEquals(v, v2) {
        \t\t\t\t\teq = false
        \t\t\t\t\treturn false
        \t\t\t\t}
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn eq
        \t\t}
        \t}
        \teq := true
        \tmc.ForEachEntry(func(k, v any) bool {
        \t\tif v == nil {
        \t\t\tif mp.Get(k) != nil || !mp.ContainsKey(k) {
        \t\t\t\teq = false
        \t\t\t\treturn false
        \t\t\t}
        \t\t} else if !javalang.ObjectsEquals(v, ToAny(mp.Get(k))) {
        \t\t\teq = false
        \t\t\treturn false
        \t\t}
        \t\treturn true
        \t})
        \treturn eq
        }

        func mapCoreHash(mc *JMapCore) int32 {
        \th := int32(0)
        \tmc.ForEachEntry(func(k, v any) bool {
        \t\tif mc.Identity {
        \t\t\th += javalang.IdentityHashCode(k) ^ javalang.IdentityHashCode(v)
        \t\t} else {
        \t\t\th += javalang.ObjectHashCode(k) ^ javalang.ObjectHashCode(v)
        \t\t}
        \t\treturn true
        \t})
        \treturn h
        }

        func mapCoreString(mc *JMapCore) string {
        \tvar sb strings.Builder
        \tsb.WriteByte('{')
        \tfirst := true
        \tmc.ForEachEntry(func(k, v any) bool {
        \t\tif !first {
        \t\t\tsb.WriteString(", ")
        \t\t}
        \t\tfirst = false
        \t\tsb.WriteString(StrVal(k))
        \t\tsb.WriteByte('=')
        \t\tsb.WriteString(StrVal(v))
        \t\treturn true
        \t})
        \tsb.WriteByte('}')
        \treturn sb.String()
        }

        type jSimpleEntry struct {
        \tk   any
        \tv   any
        \tmc  *JMapCore
        \timm bool
        }

        func (*jSimpleEntry) Implements_j2gort_Ext_java_util_Map_Entry() {}
        func (e *jSimpleEntry) GetKey() any   { return e.k }
        func (e *jSimpleEntry) GetValue() any { return e.v }
        func (e *jSimpleEntry) SetValue(nv any) any {
        \tif e.imm {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable entry"))
        \t}
        \told := e.v
        \te.v = ToAny(nv)
        \tif e.mc != nil {
        \t\te.mc.Put(e.k, e.v)
        \t}
        \treturn old
        }
        func (e *jSimpleEntry) Equals(o any) bool { return mapEntryEquals(e, o) }
        func (e *jSimpleEntry) HashCode() int32   { return mapEntryHash(e) }
        func (e *jSimpleEntry) String() string    { return StrVal(e.k) + "=" + StrVal(e.v) }

        func newSimpleImmEntry(k, v any) Ext_java_util_AbstractMap_SimpleImmutableEntry {
        \te := &Ext_java_util_AbstractMap_SimpleImmutableEntryBase{}
        \te.BindExtSelf_(e)
        \te.Core_ = &jSimpleEntry{k: ToAny(k), v: ToAny(v), imm: true}
        \treturn e
        }

        func mapEntryEquals(self Ext_java_util_Map_Entry, o any) bool {
        \tother, ok := o.(Ext_java_util_Map_Entry)
        \tif !ok || other == nil || javalang.IsNil(other) {
        \t\treturn false
        \t}
        \treturn javalang.ObjectsEquals(ToAny(self.GetKey()), ToAny(other.GetKey())) &&
        \t\tjavalang.ObjectsEquals(ToAny(self.GetValue()), ToAny(other.GetValue()))
        }

        func mapEntryHash(self Ext_java_util_Map_Entry) int32 {
        \treturn javalang.ObjectHashCode(ToAny(self.GetKey())) ^ javalang.ObjectHashCode(ToAny(self.GetValue()))
        }

        type jLiveMapKeySet struct {
        \tExt_java_util_AbstractSetBase
        \tmc *JMapCore
        }

        func newLiveMapKeySet(mc *JMapCore) Ext_java_util_Set {
        \ts := &jLiveMapKeySet{mc: mc}
        \ts.BindExtSelf_(s)
        \treturn s
        }
        func (s *jLiveMapKeySet) Size() int32             { return s.mc.Size() }
        func (s *jLiveMapKeySet) Contains(k any) bool     { return s.mc.ContainsKey(k) }
        func (s *jLiveMapKeySet) Iterator() any           { return newMapKeyIter(s.mc) }
        func (s *jLiveMapKeySet) Remove_1_any(k any) bool { return s.mc.RemoveKey(k) }
        func (s *jLiveMapKeySet) Clear()                  { s.mc.Clear() }
        func (s *jLiveMapKeySet) SetContainsStr_(k string) bool {
        \tif s == nil { ThrowNPE() }
        \treturn s.mc.ContainsStr(k)
        }
        func (s *jLiveMapKeySet) SetAddStr_(_ string) bool {
        \tif s == nil { ThrowNPE() }
        \tpanic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
        }
        func (s *jLiveMapKeySet) SetRemoveStr_(k string) bool {
        \tif s == nil { ThrowNPE() }
        \treturn s.mc.RemoveKeyStr(k)
        }
        func (s *jLiveMapKeySet) SetContainsNode_(n uint32) bool {
        \tif s == nil { ThrowNPE() }
        \treturn s.mc.ContainsNode(n)
        }
        func (s *jLiveMapKeySet) SetAddNode_(_ uint32) bool {
        \tif s == nil { ThrowNPE() }
        \tpanic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
        }
        func (s *jLiveMapKeySet) SetRemoveNode_(n uint32) bool {
        \tif s == nil { ThrowNPE() }
        \treturn s.mc.RemoveKeyNode(n)
        }

        type jLiveMapValues struct {
        \tExt_java_util_AbstractCollectionBase
        \tmc *JMapCore
        }

        func newLiveMapValues(mc *JMapCore) Ext_java_util_Collection {
        \tv := &jLiveMapValues{mc: mc}
        \tv.BindExtSelf_(v)
        \treturn v
        }
        func (v *jLiveMapValues) Size() int32         { return v.mc.Size() }
        func (v *jLiveMapValues) Contains(x any) bool { return v.mc.ContainsValue(x) }
        func (v *jLiveMapValues) Clear()              { v.mc.Clear() }
        func (v *jLiveMapValues) Iterator() any {
        \tc, ok := initMapCursor(v.mc, true)
        \tif !ok {
        \t\treturn emptyMapIter_
        \t}
        \treturn &jMapValIter{cur: c}
        }

        type jMapValIter struct {
        \tcur jMapCursor
        }

        func (*jMapValIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jMapValIter) HasNext() bool { return it.cur.hasNext() }
        func (it *jMapValIter) Next() any {
        \t_, v, _ := it.cur.nextKV()
        \treturn v
        }
        func (it *jMapValIter) Remove() { it.cur.remove() }

        type jLiveMapEntrySet struct {
        \tExt_java_util_AbstractSetBase
        \tmc *JMapCore
        }

        func newLiveMapEntrySet(mc *JMapCore) Ext_java_util_Set {
        \ts := &jLiveMapEntrySet{mc: mc}
        \ts.BindExtSelf_(s)
        \treturn s
        }
        func (s *jLiveMapEntrySet) Size() int32 { return s.mc.Size() }
        func (s *jLiveMapEntrySet) Clear()      { s.mc.Clear() }
        func (s *jLiveMapEntrySet) Contains(o any) bool {
        \te, ok := o.(Ext_java_util_Map_Entry)
        \tif !ok || e == nil || javalang.IsNil(e) {
        \t\treturn false
        \t}
        \tk := e.GetKey()
        \treturn s.mc.ContainsKey(k) && javalang.ObjectsEquals(s.mc.Get(k), ToAny(e.GetValue()))
        }
        func (s *jLiveMapEntrySet) Remove_1_any(o any) bool {
        \te, ok := o.(Ext_java_util_Map_Entry)
        \tif !ok || e == nil || javalang.IsNil(e) {
        \t\treturn false
        \t}
        \tk := e.GetKey()
        \tif !s.mc.ContainsKey(k) || !javalang.ObjectsEquals(s.mc.Get(k), ToAny(e.GetValue())) {
        \t\treturn false
        \t}
        \ts.mc.Remove(k)
        \treturn true
        }
        func (s *jLiveMapEntrySet) Iterator() any {
        \tc, ok := initMapCursor(s.mc, true)
        \tif !ok {
        \t\treturn emptyMapIter_
        \t}
        \treturn &jMapEntryIter{cur: c}
        }

        type jLiveMapEntry struct {
        \tmc   *JMapCore
        \tmode uint8
        \tslot int32
        \tk    any
        \tv    any
        }

        func (*jLiveMapEntry) Implements_j2gort_Ext_java_util_Map_Entry() {}
        func (e *jLiveMapEntry) GetKey() any   { return e.k }
        func (e *jLiveMapEntry) GetValue() any { return e.v }
        func (e *jLiveMapEntry) SetValue(nv any) any {
        \tif e.mc.immutable {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable entry"))
        \t}
        \tnv = ToAny(nv)
        \told := e.v
        \te.v = nv
        \tswitch e.mode {
        \tcase jMapModeNullOnly:
        \t\tif e.mc.mode == jMapModeNullOnly && e.slot == 0 {
        \t\t\te.mc.ext = nv
        \t\t\treturn old
        \t\t}
        \tcase jMapModeNode1:
        \t\tif e.mc.mode == jMapModeNode1 && e.slot == 0 {
        \t\t\te.mc.ext = nv
        \t\t\treturn old
        \t\t}
        \tcase jMapModeStr:
        \t\tif e.mc.mode == jMapModeStr && e.slot >= 0 && int(e.slot) < len(e.mc.strEntries) && e.mc.strEntries[e.slot].key != jTombstoneStr {
        \t\t\te.mc.strEntries[e.slot].val = nv
        \t\t\treturn old
        \t\t}
        \tcase jMapModeNode:
        \t\tif e.mc.mode == jMapModeNode && e.slot >= 0 {
        \t\t\tns := e.mc.ext.(*jNodeMapState)
        \t\t\tif int(e.slot) < len(ns.entries) && ns.entries[e.slot].active {
        \t\t\t\tns.entries[e.slot].val = nv
        \t\t\t\treturn old
        \t\t\t}
        \t\t}
        \tcase jMapModePtr:
        \t\tif e.mc.mode == jMapModePtr && e.slot >= 0 {
        \t\t\tps := e.mc.ext.(*jPtrMapState)
        \t\t\tif int(e.slot) < len(ps.entries) && ps.entries[e.slot].key != jTombstoneObj {
        \t\t\t\tps.entries[e.slot].val = nv
        \t\t\t\treturn old
        \t\t\t}
        \t\t}
        \tcase jMapModeGeneral:
        \t\tif e.mc.mode == jMapModeGeneral && e.slot > 0 {
        \t\t\tg := e.mc.ext.(*jGenMapState)
        \t\t\tif int(e.slot) <= len(g.entries) {
        \t\t\t\tg.entries[e.slot-1].Val = nv
        \t\t\t\treturn old
        \t\t\t}
        \t\t}
        \t}
        \te.mc.Put(e.k, nv)
        \treturn old
        }
        func (e *jLiveMapEntry) Equals(o any) bool { return mapEntryEquals(e, o) }
        func (e *jLiveMapEntry) HashCode() int32   { return mapEntryHash(e) }
        func (e *jLiveMapEntry) String() string    { return StrVal(e.k) + "=" + StrVal(e.v) }

        type jMapEntryIter struct {
        \tcur jMapCursor
        }

        func (*jMapEntryIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jMapEntryIter) HasNext() bool { return it.cur.hasNext() }
        func (it *jMapEntryIter) Next() any {
        \tk, v, slot := it.cur.nextKV()
        \tif slot < 0 {
        \t\treturn &jSimpleEntry{k: k, v: v, mc: it.cur.mc, imm: it.cur.mc.immutable}
        \t}
        \treturn &jLiveMapEntry{mc: it.cur.mc, mode: it.cur.mode, slot: slot, k: k, v: v}
        }
        func (it *jMapEntryIter) Remove() { it.cur.remove() }
        """;
  }

  private static String buildColRtMapCursors(RuntimeBindings rb) {
    String mapValuesGo = rb.findExtMethodGoName("java.util.Map", "values", 0);
    String mapKeySetGo = rb.findExtMethodGoName("java.util.Map", "keySet", 0);
    String mapEntrySetGo = rb.findExtMethodGoName("java.util.Map", "entrySet", 0);
    return """

        type JMapValCursor struct {
        \tcur          jMapCursor
        \tem           *jEnumMapCore
        \tfallbackIter SimpleIter
        }

        func NewMapValCursor(m any) JMapValCursor {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif mp, ok := m.(jFastMapCoreProvider); ok {
        \t\tif mc := mp.fastMapCore_(); mc != nil {
        \t\t\tc, _ := initMapValCursor(mc)
        \t\t\treturn JMapValCursor{cur: c}
        \t\t}
        \t}
        \tif ep, ok := m.(jFastEnumMapProvider); ok {
        \t\tif em := ep.fastEnumMap_(); em != nil {
        \t\t\tif em.size == 0 {
        \t\t\t\treturn JMapValCursor{}
        \t\t\t}
        \t\t\treturn JMapValCursor{em: em, cur: jMapCursor{idx: em.nextSetBit(0)}}
        \t\t}
        \t}
        \tif isJavaNilAny(m) {
        \t\tThrowNPE()
        \t}
        \treturn JMapValCursor{fallbackIter: Iterate(m.(Ext_java_util_Map)."""
        + mapValuesGo
        + "())}\n}\n\n"
        + """
        func (c *JMapValCursor) HasNext() bool {
        \tif c.cur.mc != nil {
        \t\treturn c.cur.hasValNext()
        \t}
        \tif c.em != nil {
        \t\treturn c.cur.idx >= 0
        \t}
        \tif c.fallbackIter != nil {
        \t\treturn c.fallbackIter.HasNext()
        \t}
        \treturn false
        }

        func (c *JMapValCursor) Next() any {
        \tif c.cur.mc != nil {
        \t\treturn c.cur.nextVal()
        \t}
        \tif c.em != nil {
        \t\tif c.cur.idx < 0 {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tord := c.cur.idx
        \t\tc.cur.idx = c.em.nextSetBit(ord + 1)
        \t\treturn c.em.vals[ord]
        \t}
        \tif c.fallbackIter != nil {
        \t\treturn c.fallbackIter.Next()
        \t}
        \tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        }

        type JMapKeyCursor struct {
        \tcur          jMapCursor
        \tem           *jEnumMapCore
        \tfallbackIter SimpleIter
        }

        func NewMapKeyCursor(m any) JMapKeyCursor {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif mp, ok := m.(jFastMapCoreProvider); ok {
        \t\tif mc := mp.fastMapCore_(); mc != nil {
        \t\t\tc, _ := initMapCursor(mc, false)
        \t\t\treturn JMapKeyCursor{cur: c}
        \t\t}
        \t}
        \tif ep, ok := m.(jFastEnumMapProvider); ok {
        \t\tif em := ep.fastEnumMap_(); em != nil {
        \t\t\tif em.size == 0 {
        \t\t\t\treturn JMapKeyCursor{}
        \t\t\t}
        \t\t\treturn JMapKeyCursor{em: em, cur: jMapCursor{idx: em.nextSetBit(0)}}
        \t\t}
        \t}
        \tif isJavaNilAny(m) {
        \t\tThrowNPE()
        \t}
        \treturn JMapKeyCursor{fallbackIter: Iterate(m.(Ext_java_util_Map)."""
        + mapKeySetGo
        + "())}\n}\n\n"
        + """
        func NewMapKeyStrCursor(m any) JMapKeyCursor {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif mp, ok := m.(jFastMapCoreProvider); ok {
        \t\tif mc := mp.fastMapCore_(); mc != nil {
        \t\t\tc, _ := initMapKeyStrCursor(mc)
        \t\t\treturn JMapKeyCursor{cur: c}
        \t\t}
        \t}
        \treturn NewMapKeyCursor(m)
        }

        func (c *JMapKeyCursor) HasNext() bool {
        \tif c.cur.mc != nil {
        \t\treturn c.cur.hasNext()
        \t}
        \tif c.em != nil {
        \t\treturn c.cur.idx >= 0
        \t}
        \tif c.fallbackIter != nil {
        \t\treturn c.fallbackIter.HasNext()
        \t}
        \treturn false
        }

        func (c *JMapKeyCursor) Next() any {
        \tif c.cur.mc != nil {
        \t\treturn c.cur.nextKey()
        \t}
        \tif c.em != nil {
        \t\tif c.cur.idx < 0 {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tord := c.cur.idx
        \t\tc.cur.idx = c.em.nextSetBit(ord + 1)
        \t\treturn c.em.universe[ord]
        \t}
        \tif c.fallbackIter != nil {
        \t\treturn c.fallbackIter.Next()
        \t}
        \tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        }

        func (c *JMapKeyCursor) NextStr() string {
        \tif c.cur.mc != nil {
        \t\treturn c.cur.nextKeyStr()
        \t}
        \treturn CastString(c.Next())
        }

        func (c *JMapKeyCursor) NextNode() uint32 {
        \tif c.cur.mc != nil {
        \t\treturn c.cur.nextKeyNode()
        \t}
        \treturn CastNode[uint32](c.Next())
        }

        type JMapEntryCursor struct {
        \tcur          jMapCursor
        \ttrm          *jTableRowMap
        \tem           *jEnumMapCore
        \tfallbackIter SimpleIter
        \tcurEntry     Ext_java_util_Map_Entry
        \tcurVal       any
        \tcurKeyAny    any
        \tcurKeyStr    string
        \tcurNode      uint32
        \tcurSlot      int32
        \tcurKind      uint8 // 0=any(curKeyAny), 1=str(curKeyStr), 2=node(curNode), 3=int32(curNode)
        \timm          bool
        \tsingleReady  bool
        }

        func NewMapEntryCursor(m any) JMapEntryCursor {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif mp, ok := m.(jFastMapCoreProvider); ok {
        \t\tif mc := mp.fastMapCore_(); mc != nil {
        \t\t\tc, _ := initMapCursor(mc, true)
        \t\t\t_, isImmBi := m.(*Ext_com_google_common_collect_ImmutableBiMapBase)
        \t\t\treturn JMapEntryCursor{cur: c, curSlot: -1, imm: mc.immutable || isImmBi}
        \t\t}
        \t}
        \tif trm, ok := m.(*jTableRowMap); ok {
        \t\tif trm == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\te, _ := trm.resolveEntry()
        \t\tif e == nil {
        \t\t\treturn JMapEntryCursor{}
        \t\t}
        \t\tif e.hasSingle {
        \t\t\treturn JMapEntryCursor{trm: trm, singleReady: true, curSlot: -1}
        \t\t}
        \t\tc, _ := initMapCursor(e.multi, true)
        \t\treturn JMapEntryCursor{cur: c, curSlot: -1}
        \t}
        \tif ep, ok := m.(jFastEnumMapProvider); ok {
        \t\tif em := ep.fastEnumMap_(); em != nil {
        \t\t\tif em.size == 0 {
        \t\t\t\treturn JMapEntryCursor{}
        \t\t\t}
        \t\t\treturn JMapEntryCursor{em: em, cur: jMapCursor{idx: em.nextSetBit(0)}, curSlot: -1}
        \t\t}
        \t}
        \tif isJavaNilAny(m) {
        \t\tThrowNPE()
        \t}
        \treturn JMapEntryCursor{fallbackIter: Iterate(m.(Ext_java_util_Map)."""
        + mapEntrySetGo
        + "()), curSlot: -1}\n}\n\n"
        + """
        func (c *JMapEntryCursor) NextEntry() bool {
        \tif c.singleReady {
        \t\tc.singleReady = false
        \t\te, _ := c.trm.resolveEntry()
        \t\tif e == nil {
        \t\t\treturn false
        \t\t}
        \t\tif e.hasSingle {
        \t\t\tc.curKind = 0
        \t\t\tc.curKeyAny = e.col0
        \t\t\tc.curVal = e.val0
        \t\t\tc.curSlot = 0
        \t\t\treturn true
        \t\t}
        \t\tmc, _ := initMapCursor(e.multi, true)
        \t\tc.cur = mc
        \t}
        \tif c.cur.mc != nil {
        \t\tswitch c.cur.mode {
        \t\tcase jMapModeNullOnly:
        \t\t\tif c.cur.idx > 0 || c.cur.mc.mode != jMapModeNullOnly {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tc.cur.idx = 1
        \t\t\tc.curSlot = 0
        \t\t\tc.curKind = 0
        \t\t\tc.curKeyAny = nil
        \t\t\tc.curVal = c.cur.mc.ext
        \t\t\treturn true
        \t\tcase jMapModeNode1:
        \t\t\tif c.cur.idx > 0 || c.cur.mc.mode != jMapModeNode1 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tc.cur.idx = 1
        \t\t\tc.curSlot = 0
        \t\t\tc.curKind = 2
        \t\t\tc.curNode = uint32(c.cur.mc.firstIdx)
        \t\t\tc.curVal = c.cur.mc.ext
        \t\t\treturn true
        \t\tcase jMapModeStr:
        \t\t\tc.cur.advanceSlice()
        \t\t\tif c.cur.idx >= len(c.cur.mc.strEntries) {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tslot := c.cur.idx
        \t\t\te := c.cur.mc.strEntries[slot]
        \t\t\tc.cur.idx = slot + 1
        \t\t\tc.cur.advanceSlice()
        \t\t\tc.curSlot = int32(slot)
        \t\t\tc.curKind = 1
        \t\t\tc.curKeyStr = e.key
        \t\t\tc.curVal = e.val
        \t\t\treturn true
        \t\tcase jMapModeNode:
        \t\t\tns := c.cur.mc.ext.(*jNodeMapState)
        \t\t\tc.cur.advanceSlice()
        \t\t\tif c.cur.idx >= len(ns.entries) {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tslot := c.cur.idx
        \t\t\te := ns.entries[slot]
        \t\t\tc.cur.idx = slot + 1
        \t\t\tc.cur.advanceSlice()
        \t\t\tc.curSlot = int32(slot)
        \t\t\tc.curKind = 2
        \t\t\tc.curNode = e.key
        \t\t\tc.curVal = e.val
        \t\t\treturn true
        \t\tcase jMapModePtr:
        \t\t\tps := c.cur.mc.ext.(*jPtrMapState)
        \t\t\tc.cur.advanceSlice()
        \t\t\tif c.cur.idx >= len(ps.entries) {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tslot := c.cur.idx
        \t\t\te := ps.entries[slot]
        \t\t\tc.cur.idx = slot + 1
        \t\t\tc.cur.advanceSlice()
        \t\t\tc.curSlot = int32(slot)
        \t\t\tc.curKind = 0
        \t\t\tc.curKeyAny = e.key
        \t\t\tc.curVal = e.val
        \t\t\treturn true
        \t\tcase jMapModeGeneral:
        \t\t\tif c.cur.curPos == 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tg := c.cur.mc.ext.(*jGenMapState)
        \t\t\tpos := c.cur.curPos
        \t\t\tge := &g.entries[pos-1]
        \t\t\tc.cur.curPos = ge.next()
        \t\t\tc.curSlot = pos
        \t\t\tc.curVal = ge.Val
        \t\t\tswitch ge.kind() {
        \t\t\tcase jKeyNull:
        \t\t\t\tc.curKind = 0
        \t\t\t\tc.curKeyAny = nil
        \t\t\tcase jKeyStr:
        \t\t\t\tc.curKind = 1
        \t\t\t\tc.curKeyStr = ge.Key.(string)
        \t\t\t\tc.curKeyAny = ge.Key
        \t\t\tcase jKeyNode:
        \t\t\t\tc.curKind = 2
        \t\t\t\tc.curNode = ge.u32OrHash
        \t\t\tcase jKeyInt:
        \t\t\t\tc.curKind = 3
        \t\t\t\tc.curNode = ge.u32OrHash
        \t\t\t\tc.curKeyAny = ge.Key
        \t\t\tdefault:
        \t\t\t\tc.curKind = 0
        \t\t\t\tc.curKeyAny = ge.Key
        \t\t\t}
        \t\t\treturn true
        \t\tdefault:
        \t\t\tif c.cur.idx >= len(c.cur.keys) {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\ti := c.cur.idx
        \t\t\tc.cur.idx = i + 1
        \t\t\tc.curSlot = -1
        \t\t\tc.curKind = 0
        \t\t\tc.curKeyAny = c.cur.keys[i]
        \t\t\tc.curVal = c.cur.vals[i]
        \t\t\treturn true
        \t\t}
        \t}
        \tif c.em != nil {
        \t\tif c.cur.idx < 0 {
        \t\t\treturn false
        \t\t}
        \t\tord := c.cur.idx
        \t\tc.cur.idx = c.em.nextSetBit(ord + 1)
        \t\tc.curSlot = int32(ord)
        \t\tc.curKind = 0
        \t\tc.curKeyAny = c.em.universe[ord]
        \t\tc.curVal = c.em.vals[ord]
        \t\treturn true
        \t}
        \tif c.fallbackIter != nil {
        \t\tif !c.fallbackIter.HasNext() {
        \t\t\treturn false
        \t\t}
        \t\tc.curEntry = Cast[Ext_java_util_Map_Entry](c.fallbackIter.Next())
        \t\treturn true
        \t}
        \treturn false
        }

        func (c *JMapEntryCursor) GetKey() any {
        \tif c.fallbackIter != nil {
        \t\treturn c.curEntry.GetKey()
        \t}
        \tswitch c.curKind {
        \tcase 1:
        \t\treturn BoxString(c.curKeyStr)
        \tcase 2:
        \t\treturn boxNodeById(c.curNode)
        \tcase 3:
        \t\tif c.curKeyAny != nil {
        \t\t\treturn c.curKeyAny
        \t\t}
        \t\treturn ValueOfInt(int32(c.curNode))
        \tdefault:
        \t\treturn c.curKeyAny
        \t}
        }

        func (c *JMapEntryCursor) KeyStr() string {
        \tif c.fallbackIter != nil {
        \t\treturn CastString(c.curEntry.GetKey())
        \t}
        \tif c.curKind == 1 {
        \t\treturn c.curKeyStr
        \t}
        \tif c.curKind == 0 && c.curKeyAny == nil {
        \t\treturn NullString
        \t}
        \treturn CastString(c.GetKey())
        }

        func (c *JMapEntryCursor) KeyNode() uint32 {
        \tif c.fallbackIter != nil {
        \t\treturn CastNode[uint32](c.curEntry.GetKey())
        \t}
        \tif c.curKind == 2 {
        \t\treturn c.curNode
        \t}
        \tif c.curKind == 0 && c.curKeyAny == nil {
        \t\treturn 0
        \t}
        \treturn CastNode[uint32](c.GetKey())
        }

        func (c *JMapEntryCursor) GetValue() any {
        \tif c.fallbackIter != nil {
        \t\treturn c.curEntry.GetValue()
        \t}
        \treturn c.curVal
        }

        func (c *JMapEntryCursor) SetValue(nv any) any {
        \tif c.fallbackIter != nil {
        \t\treturn c.curEntry.SetValue(ToAny(nv))
        \t}
        \tif c.imm || (c.cur.mc != nil && c.cur.mc.immutable) {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable entry"))
        \t}
        \tif c.trm != nil && c.cur.mc == nil {
        \t\tnv = CheckNotNull(ToAny(nv))
        \t\told := c.trm.putImpl(c.curKeyAny, nv)
        \t\tc.curVal = nv
        \t\treturn old
        \t}
        \tnv = ToAny(nv)
        \told := c.curVal
        \tc.curVal = nv
        \tif c.em != nil {
        \t\tif c.curSlot >= 0 && c.em.hasBit(int(c.curSlot)) {
        \t\t\tc.em.vals[c.curSlot] = nv
        \t\t\treturn old
        \t\t}
        \t\tc.em.Put(c.curKeyAny, nv)
        \t\treturn old
        \t}
        \tswitch c.cur.mode {
        \tcase jMapModeNullOnly:
        \t\tif c.cur.mc.mode == jMapModeNullOnly && c.curSlot == 0 {
        \t\t\tc.cur.mc.ext = nv
        \t\t\treturn old
        \t\t}
        \tcase jMapModeNode1:
        \t\tif c.cur.mc.mode == jMapModeNode1 && c.curSlot == 0 {
        \t\t\tc.cur.mc.ext = nv
        \t\t\treturn old
        \t\t}
        \tcase jMapModeStr:
        \t\tif c.cur.mc.mode == jMapModeStr && c.curSlot >= 0 && int(c.curSlot) < len(c.cur.mc.strEntries) && c.cur.mc.strEntries[c.curSlot].key != jTombstoneStr {
        \t\t\tc.cur.mc.strEntries[c.curSlot].val = nv
        \t\t\treturn old
        \t\t}
        \tcase jMapModeNode:
        \t\tif c.cur.mc.mode == jMapModeNode && c.curSlot >= 0 {
        \t\t\tns := c.cur.mc.ext.(*jNodeMapState)
        \t\t\tif int(c.curSlot) < len(ns.entries) && ns.entries[c.curSlot].active {
        \t\t\t\tns.entries[c.curSlot].val = nv
        \t\t\t\treturn old
        \t\t\t}
        \t\t}
        \tcase jMapModePtr:
        \t\tif c.cur.mc.mode == jMapModePtr && c.curSlot >= 0 {
        \t\t\tps := c.cur.mc.ext.(*jPtrMapState)
        \t\t\tif int(c.curSlot) < len(ps.entries) && ps.entries[c.curSlot].key != jTombstoneObj {
        \t\t\t\tps.entries[c.curSlot].val = nv
        \t\t\t\treturn old
        \t\t\t}
        \t\t}
        \tcase jMapModeGeneral:
        \t\tif c.cur.mc.mode == jMapModeGeneral && c.curSlot > 0 {
        \t\t\tg := c.cur.mc.ext.(*jGenMapState)
        \t\t\tif int(c.curSlot) <= len(g.entries) && g.entries[c.curSlot-1].kind() != jKeyFree {
        \t\t\t\tg.entries[c.curSlot-1].Val = nv
        \t\t\t\treturn old
        \t\t\t}
        \t\t}
        \t}
        \tc.cur.mc.Put(c.GetKey(), nv)
        \treturn old
        }

        type inlineStepIterable interface {
        \titerInit_() (pos int32, has bool)
        \titerNext_(pos int32) (elem any, nextPos int32)
        }

        type JIterCursor struct {
        \tslice        []any
        \tstrs         []string
        \ttv           jTypedViewCursorSource
        \tstep         inlineStepIterable
        \tmc           *JMapCore
        \tsc           *JSetCore
        \tdc           *JDequeCore
        \tem           *jEnumMapCore
        \tfallbackIter SimpleIter
        \tidx          int
        \tcurPos       int32
        \tstepPos      int32
        \tkind         uint8 // 0=empty, 1=slice, 2=mapVal, 3=mapKey, 4=jset, 5=deque, 6=enumVal, 7=enumKey, 8=fallback, 9=strSlice, 10=typedView, 11=inlineStep
        \tmapMode      uint8
        }

        func initIterMapCursor(mc *JMapCore, vals bool) JIterCursor {
        \tif mc == nil || mc.Live == 0 {
        \t\treturn JIterCursor{}
        \t}
        \tk := uint8(3)
        \tif vals {
        \t\tk = 2
        \t}
        \tswitch mc.mode {
        \tcase jMapModeNullOnly, jMapModeNode1:
        \t\treturn JIterCursor{mc: mc, kind: k, mapMode: mc.mode}
        \tcase jMapModeStr, jMapModeNode, jMapModePtr:
        \t\tstart := int(mc.firstIdx)
        \t\tc := JIterCursor{mc: mc, idx: start, kind: k, mapMode: mc.mode}
        \t\tc.advanceMapSlice()
        \t\tif c.idx != start && !mc.immutable {
        \t\t\tmc.firstIdx = int32(c.idx)
        \t\t}
        \t\treturn c
        \tcase jMapModeGeneral:
        \t\tg := mc.ext.(*jGenMapState)
        \t\tif g.head == 0 {
        \t\t\treturn JIterCursor{}
        \t\t}
        \t\treturn JIterCursor{mc: mc, curPos: g.head, kind: k, mapMode: jMapModeGeneral}
        \tcase jMapModeSortedNat, jMapModeSortedCmp:
        \t\tvar s []any
        \t\tif vals {
        \t\t\ts = mc.LiveValues()
        \t\t} else {
        \t\t\ts = mc.LiveKeys()
        \t\t}
        \t\tif len(s) == 0 {
        \t\t\treturn JIterCursor{}
        \t\t}
        \t\treturn JIterCursor{slice: s, kind: 1}
        \tdefault:
        \t\treturn JIterCursor{}
        \t}
        }

        func NewIterCursor(coll any) JIterCursor {
        \tif coll == nil {
        \t\tThrowNPE()
        \t}
        \tif ssi, ok := coll.(StrSliceIterable); ok {
        \t\tif nc, ok2 := coll.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
        \t\t\tThrowNPE()
        \t\t}
        \t\tif ss, ok2 := ssi.IterStrSlice_(); ok2 {
        \t\t\tif len(ss) == 0 {
        \t\t\t\treturn JIterCursor{}
        \t\t\t}
        \t\t\treturn JIterCursor{strs: ss, kind: 9}
        \t\t}
        \t}
        \tif tv, ok := coll.(jTypedViewCursorSource); ok {
        \t\tif nc, ok2 := coll.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
        \t\t\tThrowNPE()
        \t\t}
        \t\tif n, ok2 := tv.typedViewLen_(); ok2 {
        \t\t\tif n == 0 {
        \t\t\t\treturn JIterCursor{}
        \t\t\t}
        \t\t\treturn JIterCursor{tv: tv, kind: 10}
        \t\t}
        \t}
        \tif si, ok := coll.(SliceIterable); ok {
        \t\tif nc, ok2 := coll.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
        \t\t\tThrowNPE()
        \t\t}
        \t\ts := si.IterSlice_()
        \t\tif len(s) == 0 {
        \t\t\treturn JIterCursor{}
        \t\t}
        \t\treturn JIterCursor{slice: s, kind: 1}
        \t}
        \tif st, ok := coll.(inlineStepIterable); ok {
        \t\tpos, has := st.iterInit_()
        \t\tif !has {
        \t\t\treturn JIterCursor{}
        \t\t}
        \t\treturn JIterCursor{step: st, stepPos: pos, kind: 11}
        \t}
        \tif mv, ok := coll.(*jLiveMapValues); ok {
        \t\tif mv == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\treturn initIterMapCursor(mv.mc, true)
        \t}
        \tif rv, ok := coll.(jTableRowValues); ok {
        \t\tif rv.m == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\te, _ := rv.m.resolveEntry()
        \t\tif e == nil {
        \t\t\treturn JIterCursor{}
        \t\t}
        \t\treturn initIterMapCursor(e.multi, true)
        \t}
        \tif mk, ok := coll.(*jLiveMapKeySet); ok {
        \t\tif mk == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\treturn initIterMapCursor(mk.mc, false)
        \t}
        \tif ev, ok := coll.(*jLiveEnumMapValues); ok {
        \t\tif ev == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\tif ev.em == nil || ev.em.size == 0 {
        \t\t\treturn JIterCursor{}
        \t\t}
        \t\treturn JIterCursor{em: ev.em, idx: ev.em.nextSetBit(0), kind: 6}
        \t}
        \tif ek, ok := coll.(*jLiveEnumMapKeySet); ok {
        \t\tif ek == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\tif ek.em == nil || ek.em.size == 0 {
        \t\t\treturn JIterCursor{}
        \t\t}
        \t\treturn JIterCursor{em: ek.em, idx: ek.em.nextSetBit(0), kind: 7}
        \t}
        \tif jsp, ok := coll.(jFastJSetCoreProvider); ok {
        \t\tif sc := jsp.fastSetJCore_(); sc != nil {
        \t\t\tif sc.Size() == 0 {
        \t\t\t\treturn JIterCursor{}
        \t\t\t}
        \t\t\tstart := int(sc.firstIdx)
        \t\t\tc := JIterCursor{sc: sc, idx: start, kind: 4}
        \t\t\tc.advanceJSet()
        \t\t\tif c.idx != start && !sc.imm {
        \t\t\t\tsc.firstIdx = int32(c.idx)
        \t\t\t}
        \t\t\treturn c
        \t\t}
        \t}
        \tif sp, ok := coll.(jFastSetCoreProvider); ok {
        \t\tif mc := sp.fastSetCore_(); mc != nil {
        \t\t\treturn initIterMapCursor(mc, false)
        \t\t}
        \t}
        \tif dp, ok := coll.(jFastDequeCoreProvider); ok {
        \t\tif dc := dp.fastDequeCore_(); dc != nil {
        \t\t\tif dc.size == 0 {
        \t\t\t\treturn JIterCursor{}
        \t\t\t}
        \t\t\treturn JIterCursor{dc: dc, kind: 5}
        \t\t}
        \t}
        \tif isJavaNilAny(coll) {
        \t\tThrowNPE()
        \t}
        \tit := Iterate(coll)
        \tif ssi, ok := it.(StrSliceIterable); ok {
        \t\tif ss, ok2 := ssi.IterStrSlice_(); ok2 {
        \t\t\tif len(ss) == 0 {
        \t\t\t\treturn JIterCursor{}
        \t\t\t}
        \t\t\treturn JIterCursor{strs: ss, kind: 9}
        \t\t}
        \t}
        \tif tv, ok := it.(jTypedViewCursorSource); ok {
        \t\tif n, ok2 := tv.typedViewLen_(); ok2 {
        \t\t\tif n == 0 {
        \t\t\t\treturn JIterCursor{}
        \t\t\t}
        \t\t\treturn JIterCursor{tv: tv, kind: 10}
        \t\t}
        \t}
        \treturn JIterCursor{fallbackIter: it, kind: 8}
        }

        func (c *JIterCursor) advanceMapSlice() {
        \tswitch c.mapMode {
        \tcase jMapModeStr:
        \t\tfor c.idx < len(c.mc.strEntries) && c.mc.strEntries[c.idx].key == jTombstoneStr {
        \t\t\tc.idx++
        \t\t}
        \tcase jMapModeNode:
        \t\tes := c.mc.ext.(*jNodeMapState).entries
        \t\tfor c.idx < len(es) && !es[c.idx].active {
        \t\t\tc.idx++
        \t\t}
        \tcase jMapModePtr:
        \t\tes := c.mc.ext.(*jPtrMapState).entries
        \t\tfor c.idx < len(es) && es[c.idx].key == jTombstoneObj {
        \t\t\tc.idx++
        \t\t}
        \t}
        }

        func (c *JIterCursor) jsetLen() int {
        \tswitch c.sc.mode {
        \tcase jSetModeNode:
        \t\treturn len(c.sc.nodes)
        \tcase jSetModeStr:
        \t\treturn len(c.sc.ext.(*jSetStrState).strs)
        \tcase jSetModeAny:
        \t\treturn len(c.sc.ext.(*jSetAnyState).items)
        \tdefault:
        \t\treturn 0
        \t}
        }

        func (c *JIterCursor) advanceJSet() {
        \tswitch c.sc.mode {
        \tcase jSetModeNode:
        \t\tfor c.idx < len(c.sc.nodes) && c.sc.nodes[c.idx] == jTombstoneNode {
        \t\t\tc.idx++
        \t\t}
        \tcase jSetModeStr:
        \t\tstrs := c.sc.ext.(*jSetStrState).strs
        \t\tfor c.idx < len(strs) && strs[c.idx] == jTombstoneStr {
        \t\t\tc.idx++
        \t\t}
        \tcase jSetModeAny:
        \t\titems := c.sc.ext.(*jSetAnyState).items
        \t\tfor c.idx < len(items) && items[c.idx] == jTombstoneObj {
        \t\t\tc.idx++
        \t\t}
        \t}
        }

        func (c *JIterCursor) HasNext() bool {
        \tswitch c.kind {
        \tcase 0:
        \t\treturn false
        \tcase 1:
        \t\treturn c.idx < len(c.slice)
        \tcase 9:
        \t\treturn c.idx < len(c.strs)
        \tcase 10:
        \t\tn, _ := c.tv.typedViewLen_()
        \t\treturn c.idx < n
        \tcase 2, 3:
        \t\tswitch c.mapMode {
        \t\tcase jMapModeNullOnly:
        \t\t\treturn c.idx == 0 && c.mc.mode == jMapModeNullOnly
        \t\tcase jMapModeNode1:
        \t\t\treturn c.idx == 0 && c.mc.mode == jMapModeNode1
        \t\tcase jMapModeGeneral:
        \t\t\treturn c.curPos != 0
        \t\tcase jMapModeStr:
        \t\t\tc.advanceMapSlice()
        \t\t\treturn c.idx < len(c.mc.strEntries)
        \t\tcase jMapModeNode:
        \t\t\tc.advanceMapSlice()
        \t\t\treturn c.idx < len(c.mc.ext.(*jNodeMapState).entries)
        \t\tcase jMapModePtr:
        \t\t\tc.advanceMapSlice()
        \t\t\treturn c.idx < len(c.mc.ext.(*jPtrMapState).entries)
        \t\tdefault:
        \t\t\treturn false
        \t\t}
        \tcase 4:
        \t\tc.advanceJSet()
        \t\treturn c.idx < c.jsetLen()
        \tcase 5:
        \t\treturn int32(c.idx) < c.dc.size
        \tcase 6, 7:
        \t\treturn c.idx >= 0
        \tcase 8:
        \t\treturn c.fallbackIter.HasNext()
        \tcase 11:
        \t\treturn c.stepPos >= 0
        \tdefault:
        \t\treturn false
        \t}
        }

        func (c *JIterCursor) Next() any {
        \tswitch c.kind {
        \tcase 1:
        \t\tif c.idx >= len(c.slice) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tv := c.slice[c.idx]
        \t\tc.idx++
        \t\treturn v
        \tcase 2:
        \t\tswitch c.mapMode {
        \t\tcase jMapModeNullOnly:
        \t\t\tif c.idx > 0 || c.mc.mode != jMapModeNullOnly {
        \t\t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t\t}
        \t\t\tc.idx = 1
        \t\t\treturn c.mc.ext
        \t\tcase jMapModeNode1:
        \t\t\tif c.idx > 0 || c.mc.mode != jMapModeNode1 {
        \t\t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t\t}
        \t\t\tc.idx = 1
        \t\t\treturn c.mc.ext
        \t\tcase jMapModeStr:
        \t\t\tc.advanceMapSlice()
        \t\t\tif c.idx >= len(c.mc.strEntries) {
        \t\t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t\t}
        \t\t\tv := c.mc.strEntries[c.idx].val
        \t\t\tc.idx++
        \t\t\treturn v
        \t\tcase jMapModeNode:
        \t\t\tns := c.mc.ext.(*jNodeMapState)
        \t\t\tc.advanceMapSlice()
        \t\t\tif c.idx >= len(ns.entries) {
        \t\t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t\t}
        \t\t\tv := ns.entries[c.idx].val
        \t\t\tc.idx++
        \t\t\treturn v
        \t\tcase jMapModePtr:
        \t\t\tps := c.mc.ext.(*jPtrMapState)
        \t\t\tc.advanceMapSlice()
        \t\t\tif c.idx >= len(ps.entries) {
        \t\t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t\t}
        \t\t\tv := ps.entries[c.idx].val
        \t\t\tc.idx++
        \t\t\treturn v
        \t\tcase jMapModeGeneral:
        \t\t\tif c.curPos == 0 {
        \t\t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t\t}
        \t\t\tg := c.mc.ext.(*jGenMapState)
        \t\t\tpos := c.curPos
        \t\t\tv := g.entries[pos-1].Val
        \t\t\tc.curPos = g.entries[pos-1].next()
        \t\t\treturn v
        \t\t}
        \tcase 3:
        \t\tswitch c.mapMode {
        \t\tcase jMapModeNullOnly:
        \t\t\tif c.idx > 0 || c.mc.mode != jMapModeNullOnly {
        \t\t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t\t}
        \t\t\tc.idx = 1
        \t\t\treturn nil
        \t\tcase jMapModeNode1:
        \t\t\tif c.idx > 0 || c.mc.mode != jMapModeNode1 {
        \t\t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t\t}
        \t\t\tc.idx = 1
        \t\t\treturn boxNodeById(uint32(c.mc.firstIdx))
        \t\tcase jMapModeStr:
        \t\t\tc.advanceMapSlice()
        \t\t\tif c.idx >= len(c.mc.strEntries) {
        \t\t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t\t}
        \t\t\tek := c.mc.strEntries[c.idx].key
        \t\t\tc.idx++
        \t\t\treturn BoxString(ek)
        \t\tcase jMapModeNode:
        \t\t\tns := c.mc.ext.(*jNodeMapState)
        \t\t\tc.advanceMapSlice()
        \t\t\tif c.idx >= len(ns.entries) {
        \t\t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t\t}
        \t\t\tek := ns.entries[c.idx].key
        \t\t\tc.idx++
        \t\t\treturn boxNodeById(ek)
        \t\tcase jMapModePtr:
        \t\t\tps := c.mc.ext.(*jPtrMapState)
        \t\t\tc.advanceMapSlice()
        \t\t\tif c.idx >= len(ps.entries) {
        \t\t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t\t}
        \t\t\tek := ps.entries[c.idx].key
        \t\t\tc.idx++
        \t\t\treturn ek
        \t\tcase jMapModeGeneral:
        \t\t\tif c.curPos == 0 {
        \t\t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t\t}
        \t\t\tg := c.mc.ext.(*jGenMapState)
        \t\t\tpos := c.curPos
        \t\t\tk := g.entryKey(pos)
        \t\t\tc.curPos = g.entries[pos-1].next()
        \t\t\treturn k
        \t\t}
        \tcase 4:
        \t\tc.advanceJSet()
        \t\tif c.idx >= c.jsetLen() {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tvar v any
        \t\tswitch c.sc.mode {
        \t\tcase jSetModeNode:
        \t\t\tv = boxNodeById(c.sc.nodes[c.idx])
        \t\tcase jSetModeStr:
        \t\t\ts := c.sc.ext.(*jSetStrState).strs[c.idx]
        \t\t\tif s != NullString {
        \t\t\t\tv = s
        \t\t\t}
        \t\tcase jSetModeAny:
        \t\t\tv = c.sc.ext.(*jSetAnyState).items[c.idx]
        \t\t}
        \t\tc.idx++
        \t\treturn v
        \tcase 5:
        \t\tif int32(c.idx) >= c.dc.size {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tv := c.dc.elemAt((c.dc.head + int32(c.idx)) & c.dc.mask)
        \t\tc.idx++
        \t\treturn v
        \tcase 6:
        \t\tif c.idx < 0 {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tord := c.idx
        \t\tc.idx = c.em.nextSetBit(ord + 1)
        \t\treturn c.em.vals[ord]
        \tcase 7:
        \t\tif c.idx < 0 {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tord := c.idx
        \t\tc.idx = c.em.nextSetBit(ord + 1)
        \t\treturn c.em.universe[ord]
        \tcase 8:
        \t\treturn c.fallbackIter.Next()
        \tcase 9:
        \t\tif c.idx >= len(c.strs) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\ts := c.strs[c.idx]
        \t\tc.idx++
        \t\tif s == NullString {
        \t\t\treturn nil
        \t\t}
        \t\treturn s
        \tcase 10:
        \t\tn, _ := c.tv.typedViewLen_()
        \t\tif c.idx >= n {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tv := c.tv.typedViewAt_(c.idx)
        \t\tc.idx++
        \t\treturn v
        \tcase 11:
        \t\tv, np := c.step.iterNext_(c.stepPos)
        \t\tc.stepPos = np
        \t\treturn v
        \t}
        \tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        }

        func (c *JIterCursor) NextStr() string {
        \tif c.kind == 9 {
        \t\tif c.idx >= len(c.strs) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\ts := c.strs[c.idx]
        \t\tc.idx++
        \t\treturn s
        \t}
        \tif c.kind == 4 && c.sc.mode == jSetModeStr {
        \t\tstrs := c.sc.ext.(*jSetStrState).strs
        \t\tfor c.idx < len(strs) && strs[c.idx] == jTombstoneStr {
        \t\t\tc.idx++
        \t\t}
        \t\tif c.idx >= len(strs) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\ts := strs[c.idx]
        \t\tc.idx++
        \t\treturn s
        \t}
        \tif c.kind == 3 && c.mapMode == jMapModeStr {
        \t\tfor c.idx < len(c.mc.strEntries) && c.mc.strEntries[c.idx].key == jTombstoneStr {
        \t\t\tc.idx++
        \t\t}
        \t\tif c.idx >= len(c.mc.strEntries) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\ts := c.mc.strEntries[c.idx].key
        \t\tc.idx++
        \t\treturn s
        \t}
        \treturn CastString(c.Next())
        }

        func (c *JIterCursor) NextNode() uint32 {
        \tif c.kind == 10 {
        \t\tn, _ := c.tv.typedViewLen_()
        \t\tif c.idx >= n {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tif u, ok := c.tv.typedViewNodeAt_(c.idx); ok {
        \t\t\tc.idx++
        \t\t\treturn u
        \t\t}
        \t}
        \tif c.kind == 4 && c.sc.mode == jSetModeNode {
        \t\tnodes := c.sc.nodes
        \t\tfor c.idx < len(nodes) && nodes[c.idx] == jTombstoneNode {
        \t\t\tc.idx++
        \t\t}
        \t\tif c.idx >= len(nodes) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tu := nodes[c.idx]
        \t\tc.idx++
        \t\treturn u
        \t}
        \tif c.kind == 3 && c.mapMode == jMapModeNode1 {
        \t\tif c.idx > 0 || c.mc.mode != jMapModeNode1 {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tc.idx = 1
        \t\treturn uint32(c.mc.firstIdx)
        \t}
        \tif c.kind == 3 && c.mapMode == jMapModeNode {
        \t\tes := c.mc.ext.(*jNodeMapState).entries
        \t\tfor c.idx < len(es) && !es[c.idx].active {
        \t\t\tc.idx++
        \t\t}
        \t\tif c.idx >= len(es) {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tu := es[c.idx].key
        \t\tc.idx++
        \t\treturn u
        \t}
        \tif c.kind == 5 && c.dc.nodes != nil {
        \t\tif int32(c.idx) >= c.dc.size {
        \t\t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t\t}
        \t\tu := c.dc.nodes[(c.dc.head+int32(c.idx))&c.dc.mask]
        \t\tc.idx++
        \t\treturn u
        \t}
        \treturn CastNode[uint32](c.Next())
        }
        """;
  }

  private static String buildColRtPart1b() {
    return """

        func newJavaSet2(e0, e1 any) Ext_com_google_common_collect_ImmutableSet {
        \tif s0, ok0 := e0.(string); ok0 {
        \t\tif s1, ok1 := e1.(string); ok1 {
        \t\t\tif s0 == NullString || s1 == NullString {
        \t\t\t\tThrowNPE()
        \t\t\t}
        \t\t\tif s0 == s1 {
        \t\t\t\tpanic(javalang.NewIllegalArgumentException("duplicate element"))
        \t\t\t}
        \t\t\treturn ImmSetOfStr2(s0, s1)
        \t\t}
        \t}
        \te0 = CheckNotNull(ToAny(e0))
        \te1 = CheckNotNull(ToAny(e1))
        \tif elemJavaEquals(e0, e1, false) {
        \t\tpanic(javalang.NewIllegalArgumentException("duplicate element"))
        \t}
        \tc := &immAnySet2_{}
        \tc.buf[0] = e0
        \tc.buf[1] = e1
        \tc.set.elems_ = c.buf[:2:2]
        \tc.set.u8Tab_ = &emptyU8Tab_
        \tc.set.BindExtSelf_(&c.set)
        \ttryInitImmEnumMulti(&c.set, c.set.elems_)
        \treturn &c.set
        }

        func immSetFromCheckedSlice(dupThrows bool, items []any) Ext_com_google_common_collect_ImmutableSet {
        \tif len(items) == 0 {
        \t\treturn emptyImmutableSet_
        \t}
        \tif len(items) == 1 {
        \t\treturn newImmSet1(items[0])
        \t}
        \tif len(items) <= 8 {
        \t\tvar buf [8]any
        \t\tn := 0
        \t\tallStr := true
        \t\tfor _, raw := range items {
        \t\t\tv := CheckNotNull(ToAny(raw))
        \t\t\tif _, ok := v.(string); !ok {
        \t\t\t\tallStr = false
        \t\t\t}
        \t\t\tdup := false
        \t\t\tfor j := 0; j < n; j++ {
        \t\t\t\tif elemJavaEquals(buf[j], v, false) {
        \t\t\t\t\tif dupThrows {
        \t\t\t\t\t\tpanic(javalang.NewIllegalArgumentException("duplicate element"))
        \t\t\t\t\t}
        \t\t\t\t\tdup = true
        \t\t\t\t\tbreak
        \t\t\t\t}
        \t\t\t}
        \t\t\tif !dup {
        \t\t\t\tbuf[n] = v
        \t\t\t\tn++
        \t\t\t}
        \t\t}
        \t\tif n == 1 {
        \t\t\treturn newImmSet1(buf[0])
        \t\t}
        \t\tif allStr {
        \t\t\tvar sbuf [8]string
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\tsbuf[i] = buf[i].(string)
        \t\t\t}
        \t\t\treturn finishImmStrSetSmall(sbuf[:n])
        \t\t}
        \t\treturn finishImmSetSmall(buf[:n])
        \t}
        \tsc := &JSetCore{}
        \tsc.InitExpectedSize(int32(len(items)))
        \tfor _, raw := range items {
        \t\tv := CheckNotNull(ToAny(raw))
        \t\tif !sc.Add(v) && dupThrows {
        \t\t\tpanic(javalang.NewIllegalArgumentException("duplicate element"))
        \t\t}
        \t}
        \treturn wrapImmutableSetFromSetCoreOwn(sc)
        }

        func wrapImmutableMapEntrySet(mc *JMapCore) Ext_com_google_common_collect_ImmutableSet {
        \tif mc == nil || mc.Size() == 0 {
        \t\treturn emptyImmutableSet_
        \t}
        \tif mc.Size() == 1 {
        \t\tvar e0 any
        \t\tmc.ForEachEntry(func(k, v any) bool {
        \t\t\te0 = &jSimpleEntry{k: k, v: v, imm: true}
        \t\t\treturn false
        \t\t})
        \t\treturn newImmSet1(e0)
        \t}
        \telems := make([]any, 0, mc.Size())
        \tmc.ForEachEntry(func(k, v any) bool {
        \t\telems = append(elems, &jSimpleEntry{k: k, v: v, imm: true})
        \t\treturn true
        \t})
        \treturn finishImmSetMulti(elems, nil)
        }

        var emptyImmutableSortedSet_ Ext_com_google_common_collect_ImmutableSortedSet
        var emptyImmutableSortedMap_ *Ext_com_google_common_collect_ImmutableSortedMap
        func init() {
        \t{
        \t\tmc := NewSortedJMapCore(nil)
        \t\tmc.sealImmutable()
        \t\tobj := &Ext_com_google_common_collect_ImmutableSortedSetBase{sortedCore_: mc}
        \t\tobj.elems_ = emptyAnySlice_
        \t\tobj.BindExtSelf_(obj)
        \t\temptyImmutableSortedSet_ = obj
        \t}
        \t{
        \t\tobj := &Ext_com_google_common_collect_ImmutableSortedMap{}
        \t\tobj.BindExtSelf_(obj)
        \t\tobj.inlineCore.mode = jMapModeSortedNat
        \t\tobj.inlineCore.sealImmutable()
        \t\temptyImmutableSortedMap_ = obj
        \t}
        }

        func wrapImmutableSortedSet(mc *JMapCore) Ext_com_google_common_collect_ImmutableSortedSet {
        \tif mc == nil || (mc.Size() == 0 && mc.IsNaturalOrder()) {
        \t\treturn emptyImmutableSortedSet_
        \t}
        \tmc.sealImmutable()
        \tobj := &Ext_com_google_common_collect_ImmutableSortedSetBase{sortedCore_: mc}
        \tobj.BindExtSelf_(obj)
        \tpopulateImmSortedSetElems(obj, mc)
        \treturn obj
        }

        func wrapImmutableSortedMap(mc *JMapCore) *Ext_com_google_common_collect_ImmutableSortedMap {
        \tif mc == nil || (mc.Size() == 0 && mc.IsNaturalOrder()) {
        \t\treturn emptyImmutableSortedMap_
        \t}
        \tmc.sealImmutable()
        \tobj := &Ext_com_google_common_collect_ImmutableSortedMap{}
        \tobj.BindExtSelf_(obj)
        \tobj.Core_ = mc
        \treturn obj
        }

        func immSetFromStrSlice(dupThrows bool, src []string) Ext_com_google_common_collect_ImmutableSet {
        \tif len(src) == 0 {
        \t\treturn emptyImmutableSet_
        \t}
        \tif len(src) == 1 {
        \t\treturn ImmSetOfStr1(src[0])
        \t}
        \tif len(src) <= 8 {
        \t\tvar buf [8]string
        \t\tn := 0
        \t\tfor _, s := range src {
        \t\t\tif s == NullString {
        \t\t\t\tThrowNPE()
        \t\t\t}
        \t\t\tdup := false
        \t\t\tfor j := 0; j < n; j++ {
        \t\t\t\tif buf[j] == s {
        \t\t\t\t\tif dupThrows {
        \t\t\t\t\t\tpanic(javalang.NewIllegalArgumentException("duplicate element"))
        \t\t\t\t\t}
        \t\t\t\t\tdup = true
        \t\t\t\t\tbreak
        \t\t\t\t}
        \t\t\t}
        \t\t\tif !dup {
        \t\t\t\tbuf[n] = s
        \t\t\t\tn++
        \t\t\t}
        \t\t}
        \t\tif n == 1 {
        \t\t\treturn ImmSetOfStr1(buf[0])
        \t\t}
        \t\treturn finishImmStrSetSmall(buf[:n])
        \t}
        \tsc := &JSetCore{}
        \tsc.InitExpectedSize(int32(len(src)))
        \tfor _, s := range src {
        \t\tif s == NullString {
        \t\t\tThrowNPE()
        \t\t}
        \t\tif !sc.AddStr(s) && dupThrows {
        \t\t\tpanic(javalang.NewIllegalArgumentException("duplicate element"))
        \t\t}
        \t}
        \treturn wrapImmutableSetFromSetCoreOwn(sc)
        }

        func immListCopyOfSlice(src []any) Ext_com_google_common_collect_ImmutableList {
        \tCheckNotNullSlice(src)
        \tswitch len(src) {
        \tcase 0:
        \t\treturn emptyImmutableList_
        \tcase 1:
        \t\treturn newImmList1(src[0])
        \tcase 2:
        \t\tif s0, ok0 := src[0].(string); ok0 {
        \t\t\tif s1, ok1 := src[1].(string); ok1 {
        \t\t\t\treturn ImmListOfStr2(s0, s1)
        \t\t\t}
        \t\t}
        \t\treturn newImmList2(src[0], src[1])
        \tcase 3:
        \t\tif s0, ok0 := src[0].(string); ok0 {
        \t\t\tif s1, ok1 := src[1].(string); ok1 {
        \t\t\t\tif s2, ok2 := src[2].(string); ok2 {
        \t\t\t\t\treturn ImmListOfStr3(s0, s1, s2)
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t\treturn newImmList3(src[0], src[1], src[2])
        \tcase 4:
        \t\tif s0, ok0 := src[0].(string); ok0 {
        \t\t\tif s1, ok1 := src[1].(string); ok1 {
        \t\t\t\tif s2, ok2 := src[2].(string); ok2 {
        \t\t\t\t\tif s3, ok3 := src[3].(string); ok3 {
        \t\t\t\t\t\treturn ImmListOfStr4(s0, s1, s2, s3)
        \t\t\t\t\t}
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t\treturn newImmList4(src[0], src[1], src[2], src[3])
        \tcase 5:
        \t\treturn newImmList5(src[0], src[1], src[2], src[3], src[4])
        \tcase 6:
        \t\treturn newImmList6(src[0], src[1], src[2], src[3], src[4], src[5])
        \tcase 7:
        \t\treturn newImmList7(src[0], src[1], src[2], src[3], src[4], src[5], src[6])
        \tcase 8:
        \t\treturn newImmList8(src[0], src[1], src[2], src[3], src[4], src[5], src[6], src[7])
        \tdefault:
        \t\tif _, ok0 := src[0].(string); ok0 {
        \t\t\tallStr := true
        \t\t\tfor _, v := range src {
        \t\t\t\ts, ok := v.(string)
        \t\t\t\tif !ok {
        \t\t\t\t\tallStr = false
        \t\t\t\t\tbreak
        \t\t\t\t}
        \t\t\t\tif s == NullString {
        \t\t\t\t\tThrowNPE()
        \t\t\t\t}
        \t\t\t}
        \t\t\tif allStr {
        \t\t\t\tstrs := make([]string, len(src))
        \t\t\t\tfor i, v := range src {
        \t\t\t\t\tstrs[i] = v.(string)
        \t\t\t\t}
        \t\t\t\treturn &jImmStrList{strs: strs}
        \t\t\t}
        \t\t}
        \t\titems := make([]any, len(src))
        \t\tfor i, v := range src {
        \t\t\titems[i] = CheckNotNull(ToAny(v))
        \t\t}
        \t\treturn wrapImmutableListOwn(items)
        \t}
        }

        func immListCopyOfIterable(src any) Ext_com_google_common_collect_ImmutableList {
        \tCheckNotNull(src)
        \tif il, ok := src.(Ext_com_google_common_collect_ImmutableList); ok && !isJavaNilAny(il) {
        \t\treturn il
        \t}
        \tif ssi, ok := src.(StrSliceIterable); ok {
        \t\tif ss, ok2 := ssi.IterStrSlice_(); ok2 {
        \t\t\treturn wrapImmStrSliceCopy(ss)
        \t\t}
        \t}
        \tif si, ok := src.(SliceIterable); ok {
        \t\ts := si.IterSlice_()
        \t\tif len(s) == 0 {
        \t\t\treturn emptyImmutableList_
        \t\t}
        \t\treturn immListCopyOfSlice(s)
        \t}
        \tif c, ok := src.(Ext_java_util_Collection); ok && c.IsEmpty() {
        \t\treturn emptyImmutableList_
        \t}
        \tvar b Ext_com_google_common_collect_ImmutableList_Builder
        \tb.AddAllVoid_(src)
        \treturn b.BuildImmList_()
        }

        func immSetCopyOfSlice(src []any) Ext_com_google_common_collect_ImmutableSet {
        \tCheckNotNullSlice(src)
        \treturn immSetFromCheckedSlice(false, src)
        }

        func immSetCopyOfIterable(src any) Ext_com_google_common_collect_ImmutableSet {
        \tCheckNotNull(src)
        \tif is, ok := src.(Ext_com_google_common_collect_ImmutableSet); ok && !isJavaNilAny(is) {
        \t\treturn is
        \t}
        \tif ep, ok := src.(jFastEnumSetProvider); ok {
        \t\tif es := ep.fastEnumSet_(); es != nil {
        \t\t\treturn immSetFromEnumSet(es)
        \t\t}
        \t}
        \tif ssi, ok := src.(StrSliceIterable); ok {
        \t\tif ss, ok2 := ssi.IterStrSlice_(); ok2 {
        \t\t\treturn immSetFromStrSlice(false, ss)
        \t\t}
        \t}
        \tif si, ok := src.(SliceIterable); ok {
        \t\treturn immSetFromCheckedSlice(false, si.IterSlice_())
        \t}
        \tif jsp, ok := src.(jFastJSetCoreProvider); ok {
        \t\tif osc := jsp.fastSetJCore_(); osc != nil && !osc.identity {
        \t\t\treturn immSetFromSetCoreCopyChecked(osc)
        \t\t}
        \t}
        \tif sp, ok := src.(jFastSetCoreProvider); ok {
        \t\tif smc := sp.fastSetCore_(); smc != nil && !smc.Identity {
        \t\t\tif smc.Size() == 0 {
        \t\t\t\treturn emptyImmutableSet_
        \t\t\t}
        \t\t\tif smc.Size() == 1 {
        \t\t\t\treturn newImmSet1(smc.FirstKey())
        \t\t\t}
        \t\t\treturn immSetFromCheckedSlice(false, smc.LiveKeys())
        \t\t}
        \t}
        \tif c, ok := src.(Ext_java_util_Collection); ok && c.IsEmpty() {
        \t\treturn emptyImmutableSet_
        \t}
        \tsc := &JSetCore{}
        \tsetCoreAddAllNonNull(sc, src)
        \treturn wrapImmutableSetFromSetCoreOwn(sc)
        }

        func immSetFromSetCoreCopyChecked(src *JSetCore) Ext_com_google_common_collect_ImmutableSet {
        \tif src == nil || src.Size() == 0 {
        \t\treturn emptyImmutableSet_
        \t}
        \tn := int(src.Size())
        \tif src.mode == jSetModeStr {
        \t\tst := src.ext.(*jSetStrState)
        \t\tif n == 1 {
        \t\t\tfor i := int(src.firstIdx); i < len(st.strs); i++ {
        \t\t\t\ts := st.strs[i]
        \t\t\t\tif s == jTombstoneStr {
        \t\t\t\t\tcontinue
        \t\t\t\t}
        \t\t\t\tif s == NullString {
        \t\t\t\t\tThrowNPE()
        \t\t\t\t}
        \t\t\t\treturn ImmSetOfStr1(s)
        \t\t\t}
        \t\t\treturn emptyImmutableSet_
        \t\t}
        \t\tif n <= 8 {
        \t\t\tvar sbuf [8]string
        \t\t\tw := 0
        \t\t\tfor i := int(src.firstIdx); i < len(st.strs); i++ {
        \t\t\t\ts := st.strs[i]
        \t\t\t\tif s == jTombstoneStr {
        \t\t\t\t\tcontinue
        \t\t\t\t}
        \t\t\t\tif s == NullString {
        \t\t\t\t\tThrowNPE()
        \t\t\t\t}
        \t\t\t\tsbuf[w] = s
        \t\t\t\tw++
        \t\t\t}
        \t\t\tif w == 1 {
        \t\t\t\treturn ImmSetOfStr1(sbuf[0])
        \t\t\t}
        \t\t\treturn finishImmStrSetSmall(sbuf[:w])
        \t\t}
        \t\tstrs := make([]string, 0, n)
        \t\tfor i := int(src.firstIdx); i < len(st.strs); i++ {
        \t\t\ts := st.strs[i]
        \t\t\tif s == jTombstoneStr {
        \t\t\t\tcontinue
        \t\t\t}
        \t\t\tif s == NullString {
        \t\t\t\tThrowNPE()
        \t\t\t}
        \t\t\tstrs = append(strs, s)
        \t\t}
        \t\tif len(strs) == 1 {
        \t\t\treturn ImmSetOfStr1(strs[0])
        \t\t}
        \t\treturn finishImmStrSetOwn(strs, nil)
        \t}
        \tif n == 1 {
        \t\treturn newImmSet1(src.FirstKey())
        \t}
        \tif n <= 8 {
        \t\tvar buf [8]any
        \t\tw := 0
        \t\tsrc.ForEach(func(k any) bool {
        \t\t\tbuf[w] = CheckNotNull(k)
        \t\t\tw++
        \t\t\treturn true
        \t\t})
        \t\tif w == 1 {
        \t\t\treturn newImmSet1(buf[0])
        \t\t}
        \t\treturn finishImmSetSmall(buf[:w])
        \t}
        \telems := make([]any, 0, n)
        \tsrc.ForEach(func(k any) bool {
        \t\telems = append(elems, CheckNotNull(k))
        \t\treturn true
        \t})
        \treturn finishImmSetMulti(elems, nil)
        }

        func immSetFromSetCoreCopy(src *JSetCore) Ext_com_google_common_collect_ImmutableSet {
        \treturn immSetFromSetCoreCopyChecked(src)
        }

        func immSortedSetFromCoreCopy(src *JMapCore) Ext_com_google_common_collect_ImmutableSortedSet {
        \tif src == nil || (src.Size() == 0 && src.IsNaturalOrder()) {
        \t\treturn emptyImmutableSortedSet_
        \t}
        \tdst := cloneJMapCore(src)
        \treturn wrapImmutableSortedSet(dst)
        }

        func immMapCopyOf(src Ext_java_util_Map) Ext_com_google_common_collect_ImmutableMap {
        \tCheckNotNull(src)
        \tif im, ok := any(src).(Ext_com_google_common_collect_ImmutableMap); ok && !isJavaNilAny(im) {
        \t\treturn im
        \t}
        \tif src.IsEmpty() {
        \t\treturn emptyImmutableMap_
        \t}
        \tobj := newImmutableMapShell()
        \tmc := &obj.inlineCore
        \tmc.InitExpectedSize(src.Size())
        \tcopyExtMapIntoNonNull(mc, src)
        \tmc.sealImmutable()
        \treturn obj
        }

        func immMapCopyOfEntries(entries any) Ext_com_google_common_collect_ImmutableMap {
        \tCheckNotNull(entries)
        \tobj := newImmutableMapShell()
        \tmc := &obj.inlineCore
        \tfor it := Iterate(entries); it.HasNext(); {
        \t\te := Cast[Ext_java_util_Map_Entry](CheckNotNull(it.Next()))
        \t\tk := CheckNotNull(ToAny(e.GetKey()))
        \t\tv := CheckNotNull(ToAny(e.GetValue()))
        \t\tif mc.ContainsKey(k) {
        \t\t\tpanic(javalang.NewIllegalArgumentException("Multiple entries with same key: " + StrVal(k)))
        \t\t}
        \t\tmc.Put(k, v)
        \t}
        \tif mc.IsEmpty() {
        \t\treturn emptyImmutableMap_
        \t}
        \tmc.sealImmutable()
        \treturn obj
        }

        func immSortedMapCopyOf(src Ext_java_util_Map) *Ext_com_google_common_collect_ImmutableSortedMap {
        \tCheckNotNull(src)
        \tif ism, ok := any(src).(*Ext_com_google_common_collect_ImmutableSortedMap); ok && ism != nil && ism.mc().IsNaturalOrder() {
        \t\treturn ism
        \t}
        \tif src.IsEmpty() {
        \t\treturn emptyImmutableSortedMap_
        \t}
        \tobj := &Ext_com_google_common_collect_ImmutableSortedMap{}
        \tobj.BindExtSelf_(obj)
        \tmc := initSortedCore(&obj.inlineCore, nil)
        \tcopyExtMapIntoNonNull(mc, src)
        \tmc.sealImmutable()
        \treturn obj
        }

        type jImmMapBuilderCore struct {
        \tmc     *JMapCore
        \thasDup bool
        \tdupKey any
        \tsorted bool
        }

        func immMapBuilderPut(b *jImmMapBuilderCore, k, v any) {
        \tk = CheckNotNull(ToAny(k))
        \tv = CheckNotNull(ToAny(v))
        \tif b.mc.ContainsKey(k) {
        \t\tb.hasDup = true
        \t\tb.dupKey = k
        \t}
        \tb.mc.Put(k, v)
        }

        func immMapBuilderPutAll(b *jImmMapBuilderCore, m Ext_java_util_Map) {
        \tCheckNotNull(m)
        \tif mp, ok := any(m).(jFastMapCoreProvider); ok {
        \t\tif smc := mp.fastMapCore_(); smc != nil {
        \t\t\tb.mc.InitExpectedSize(b.mc.Size() + smc.Size())
        \t\t\tsmc.ForEachEntry(func(k, v any) bool {
        \t\t\t\timmMapBuilderPut(b, k, v)
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn
        \t\t}
        \t}
        \tif ep, ok := any(m).(jFastEnumMapProvider); ok {
        \t\tif em := ep.fastEnumMap_(); em != nil {
        \t\t\tb.mc.InitExpectedSize(b.mc.Size() + em.size)
        \t\t\tem.ForEachEntry(func(k, v any) bool {
        \t\t\t\timmMapBuilderPut(b, k, v)
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn
        \t\t}
        \t}
        \tfor it := Iterate(m.EntrySet()); it.HasNext(); {
        \t\te := Cast[Ext_java_util_Map_Entry](it.Next())
        \t\timmMapBuilderPut(b, e.GetKey(), e.GetValue())
        \t}
        }

        func immMapBuilderBuild(b *jImmMapBuilderCore, throwOnDup bool) any {
        \tif throwOnDup && b.hasDup {
        \t\tpanic(javalang.NewIllegalArgumentException("Multiple entries with same key: " + StrVal(b.dupKey)))
        \t}
        \tif b.sorted {
        \t\tif b.mc.Size() == 0 && b.mc.IsNaturalOrder() {
        \t\t\treturn emptyImmutableSortedMap_
        \t\t}
        \t\tobj := &Ext_com_google_common_collect_ImmutableSortedMap{}
        \t\tobj.BindExtSelf_(obj)
        \t\tdst := copySortedCoreInto(&obj.inlineCore, b.mc)
        \t\tdst.sealImmutable()
        \t\treturn obj
        \t}
        \tif b.mc.Size() == 0 {
        \t\treturn emptyImmutableMap_
        \t}
        \tobj := newImmutableMapShell()
        \tdst := copyHashCoreInto(&obj.inlineCore, b.mc)
        \tdst.sealImmutable()
        \treturn obj
        }

        func immBiMapBuilderBuild(b *jImmMapBuilderCore, throwOnDup bool) any {
        \tif throwOnDup && b.hasDup {
        \t\tpanic(javalang.NewIllegalArgumentException("Multiple entries with same key: " + StrVal(b.dupKey)))
        \t}
        \treturn newImmutableBiMapFrom(wrapLinkedHashMap(cloneJMapCore(b.mc)))
        }

        func immListSortedCopyOf(cmp Ext_java_util_Comparator, it Ext_java_lang_Iterable) Ext_com_google_common_collect_ImmutableList {
        \tCheckNotNull(it)
        \tvar items []any
        \tfor iter := NewIterCursor(it); iter.HasNext(); {
        \t\titems = append(items, CheckNotNull(iter.Next()))
        \t}
        \tif len(items) == 0 {
        \t\treturn emptyImmutableList_
        \t}
        \tfn := comparatorFuncOrNatural(cmp)
        \tsort.SliceStable(items, func(i, j int) bool { return fn(items[i], items[j]) < 0 })
        \treturn wrapImmutableListOwn(items)
        }

        func immSortedSetCopyOf(cmp Ext_java_util_Comparator, col any) Ext_com_google_common_collect_ImmutableSortedSet {
        \tCheckNotNull(col)
        \tnat := isNaturalOrderComparator(cmp)
        \tif nat {
        \t\tif iss, ok := col.(Ext_com_google_common_collect_ImmutableSortedSet); ok && !isJavaNilAny(iss) {
        \t\t\tif sp, ok2 := col.(jFastSetCoreProvider); ok2 {
        \t\t\t\tif smc := sp.fastSetCore_(); smc != nil && smc.IsNaturalOrder() {
        \t\t\t\t\treturn iss
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t\tif c, ok := col.(Ext_java_util_Collection); ok && c.IsEmpty() {
        \t\t\treturn emptyImmutableSortedSet_
        \t\t}
        \t}
        \tmc := newSortedCoreFromCmp(cmp)
        \tmapCoreAddAllNonNull(mc, col)
        \treturn wrapImmutableSortedSet(mc)
        }

        // --- Bitmask EnumSet & Slice/Inline EnumMap ---

        func resolveEnumUniverseFromElem(v any) []any {
        \tif cip, ok := v.(interface{ JavaClassInfo_() (string, string) }); ok {
        \t\tbinName, _ := cip.JavaClassInfo_()
        \t\tif u := LookupEnumConstantsByName(binName); u != nil {
        \t\t\treturn u
        \t\t}
        \t}
        \tpanic(javalang.NewClassCastException(fmt.Sprintf("%T cannot be cast to java.lang.Enum", v)))
        }

        type jEnumSetCore struct {
        \tuniverse []any
        \tbits0    uint64
        \tbits1    uint64
        \tmore     []uint64
        \tsize     int32
        }

        func (es *jEnumSetCore) Size() int32 { return es.size }

        func (es *jEnumSetCore) hasBit(ord int) bool {
        \tw := ord >> 6
        \tmask := uint64(1) << uint(ord&63)
        \tswitch w {
        \tcase 0:
        \t\treturn (es.bits0 & mask) != 0
        \tcase 1:
        \t\treturn (es.bits1 & mask) != 0
        \tdefault:
        \t\tmw := w - 2
        \t\treturn mw < len(es.more) && (es.more[mw]&mask) != 0
        \t}
        }

        func (es *jEnumSetCore) setOrd(ord int) bool {
        \tw := ord >> 6
        \tmask := uint64(1) << uint(ord&63)
        \tswitch w {
        \tcase 0:
        \t\tif (es.bits0 & mask) != 0 {
        \t\t\treturn false
        \t\t}
        \t\tes.bits0 |= mask
        \t\tes.size++
        \t\treturn true
        \tcase 1:
        \t\tif (es.bits1 & mask) != 0 {
        \t\t\treturn false
        \t\t}
        \t\tes.bits1 |= mask
        \t\tes.size++
        \t\treturn true
        \tdefault:
        \t\tmw := w - 2
        \t\tif mw >= len(es.more) {
        \t\t\tneed := ((len(es.universe) + 63) >> 6) - 2
        \t\t\tif mw+1 > need {
        \t\t\t\tneed = mw + 1
        \t\t\t}
        \t\t\tnext := make([]uint64, need)
        \t\t\tcopy(next, es.more)
        \t\t\tes.more = next
        \t\t}
        \t\tif (es.more[mw] & mask) != 0 {
        \t\t\treturn false
        \t\t}
        \t\tes.more[mw] |= mask
        \t\tes.size++
        \t\treturn true
        \t}
        }

        func (es *jEnumSetCore) clearOrd(ord int) bool {
        \tw := ord >> 6
        \tmask := uint64(1) << uint(ord&63)
        \tswitch w {
        \tcase 0:
        \t\tif (es.bits0 & mask) == 0 {
        \t\t\treturn false
        \t\t}
        \t\tes.bits0 &^= mask
        \t\tes.size--
        \t\treturn true
        \tcase 1:
        \t\tif (es.bits1 & mask) == 0 {
        \t\t\treturn false
        \t\t}
        \t\tes.bits1 &^= mask
        \t\tes.size--
        \t\treturn true
        \tdefault:
        \t\tmw := w - 2
        \t\tif mw >= len(es.more) || (es.more[mw]&mask) == 0 {
        \t\t\treturn false
        \t\t}
        \t\tes.more[mw] &^= mask
        \t\tes.size--
        \t\treturn true
        \t}
        }

        func (es *jEnumSetCore) nextSetBit(fromOrd int) int {
        \tif fromOrd < 0 {
        \t\tfromOrd = 0
        \t}
        \tw := fromOrd >> 6
        \tshift := uint(fromOrd & 63)
        \tif w == 0 {
        \t\tword := es.bits0 & (^uint64(0) << shift)
        \t\tif word != 0 {
        \t\t\treturn bits.TrailingZeros64(word)
        \t\t}
        \t\tw = 1
        \t\tshift = 0
        \t}
        \tif w == 1 {
        \t\tword := es.bits1 & (^uint64(0) << shift)
        \t\tif word != 0 {
        \t\t\treturn 64 + bits.TrailingZeros64(word)
        \t\t}
        \t\tw = 2
        \t\tshift = 0
        \t}
        \tfor mw := w - 2; mw < len(es.more); mw++ {
        \t\tword := es.more[mw] & (^uint64(0) << shift)
        \t\tshift = 0
        \t\tif word != 0 {
        \t\t\treturn ((mw + 2) << 6) + bits.TrailingZeros64(word)
        \t\t}
        \t}
        \treturn -1
        }

        func (es *jEnumSetCore) recomputeSize() {
        \tn := bits.OnesCount64(es.bits0) + bits.OnesCount64(es.bits1)
        \tfor _, w := range es.more {
        \t\tn += bits.OnesCount64(w)
        \t}
        \tes.size = int32(n)
        }

        func (es *jEnumSetCore) ordIfValid(v any) (int, bool) {
        \tv = ToAny(v)
        \tif v == nil {
        \t\treturn -1, false
        \t}
        \tea, ok := v.(EnumAccessor)
        \tif !ok {
        \t\treturn -1, false
        \t}
        \teb := ea.GetEnumBase()
        \tif eb == nil {
        \t\treturn -1, false
        \t}
        \tif es.universe == nil {
        \t\treturn -1, false
        \t}
        \tord := int(eb.Ordinal_)
        \tif ord >= 0 && ord < len(es.universe) && es.universe[ord] == v {
        \t\treturn ord, true
        \t}
        \treturn -1, false
        }

        func (es *jEnumSetCore) mustOrd(v any) int {
        \tv = CheckNotNull(ToAny(v))
        \tea, ok := v.(EnumAccessor)
        \tif !ok {
        \t\tpanic(javalang.NewClassCastException(fmt.Sprintf("%T cannot be cast to java.lang.Enum", v)))
        \t}
        \teb := ea.GetEnumBase()
        \tif eb == nil {
        \t\tThrowNPE()
        \t}
        \tif es.universe == nil {
        \t\tes.universe = resolveEnumUniverseFromElem(v)
        \t}
        \tord := int(eb.Ordinal_)
        \tif ord >= 0 && ord < len(es.universe) && es.universe[ord] == v {
        \t\treturn ord
        \t}
        \tpanic(javalang.NewClassCastException("enum type mismatch in EnumSet"))
        }

        func (es *jEnumSetCore) Contains(v any) bool {
        \tif es.size == 0 {
        \t\treturn false
        \t}
        \tord, ok := es.ordIfValid(v)
        \treturn ok && es.hasBit(ord)
        }

        func (es *jEnumSetCore) Add(v any) bool {
        \treturn es.setOrd(es.mustOrd(v))
        }

        func (es *jEnumSetCore) Remove(v any) bool {
        \tif es.size == 0 {
        \t\treturn false
        \t}
        \tord, ok := es.ordIfValid(v)
        \tif !ok {
        \t\treturn false
        \t}
        \treturn es.clearOrd(ord)
        }

        func (es *jEnumSetCore) Clear() {
        \tes.bits0 = 0
        \tes.bits1 = 0
        \tif len(es.more) > 0 {
        \t\tclear(es.more)
        \t}
        \tes.size = 0
        }

        func (es *jEnumSetCore) sameUniverse(other *jEnumSetCore) bool {
        \tif len(es.universe) != len(other.universe) {
        \t\treturn false
        \t}
        \tif len(es.universe) == 0 {
        \t\treturn es.universe != nil && other.universe != nil
        \t}
        \treturn es.universe[0] == other.universe[0]
        }

        func (es *jEnumSetCore) initCopy(src *jEnumSetCore) {
        \tes.universe = src.universe
        \tes.bits0 = src.bits0
        \tes.bits1 = src.bits1
        \tif len(src.more) > 0 {
        \t\tes.more = make([]uint64, len(src.more))
        \t\tcopy(es.more, src.more)
        \t}
        \tes.size = src.size
        }

        func (es *jEnumSetCore) initAll(u []any) {
        \tes.universe = u
        \tn := len(u)
        \tes.size = int32(n)
        \tif n == 0 {
        \t\treturn
        \t}
        \tif n <= 64 {
        \t\tif n == 64 {
        \t\t\tes.bits0 = ^uint64(0)
        \t\t} else {
        \t\t\tes.bits0 = (uint64(1) << uint(n)) - 1
        \t\t}
        \t\treturn
        \t}
        \tes.bits0 = ^uint64(0)
        \tif n <= 128 {
        \t\tif n == 128 {
        \t\t\tes.bits1 = ^uint64(0)
        \t\t} else {
        \t\t\tes.bits1 = (uint64(1) << uint(n-64)) - 1
        \t\t}
        \t\treturn
        \t}
        \tes.bits1 = ^uint64(0)
        \tmWords := ((n + 63) >> 6) - 2
        \tes.more = make([]uint64, mWords)
        \tfor i := 0; i < mWords-1; i++ {
        \t\tes.more[i] = ^uint64(0)
        \t}
        \trem := n & 63
        \tif rem == 0 {
        \t\tes.more[mWords-1] = ^uint64(0)
        \t} else {
        \t\tes.more[mWords-1] = (uint64(1) << uint(rem)) - 1
        \t}
        }

        func (es *jEnumSetCore) ForEach(fn func(any) bool) bool {
        \tif es.size == 0 {
        \t\treturn true
        \t}
        \tfor ord := es.nextSetBit(0); ord >= 0; ord = es.nextSetBit(ord + 1) {
        \t\tif !fn(es.universe[ord]) {
        \t\t\treturn false
        \t\t}
        \t}
        \treturn true
        }

        func (es *jEnumSetCore) ForEachConsumer(action Ext_java_util_function_Consumer) {
        \tCheckNotNull(action)
        \tfor ord := es.nextSetBit(0); ord >= 0; ord = es.nextSetBit(ord + 1) {
        \t\taction.Accept(es.universe[ord])
        \t}
        }

        func (es *jEnumSetCore) ToArray() []any {
        \tif es.size == 0 {
        \t\treturn []any{}
        \t}
        \tout := make([]any, 0, es.size)
        \tfor ord := es.nextSetBit(0); ord >= 0; ord = es.nextSetBit(ord + 1) {
        \t\tout = append(out, es.universe[ord])
        \t}
        \treturn out
        }

        func (es *jEnumSetCore) ContainsAll(other Ext_java_util_Collection) bool {
        \tCheckNotNull(other)
        \tif ep, ok := any(other).(jFastEnumSetProvider); ok {
        \t\tif oes := ep.fastEnumSet_(); oes != nil {
        \t\t\tif oes.size == 0 {
        \t\t\t\treturn true
        \t\t\t}
        \t\t\tif es.size == 0 || !es.sameUniverse(oes) {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tif (oes.bits0&^es.bits0) != 0 || (oes.bits1&^es.bits1) != 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tfor i, ow := range oes.more {
        \t\t\t\tvar ew uint64
        \t\t\t\tif i < len(es.more) {
        \t\t\t\t\tew = es.more[i]
        \t\t\t\t}
        \t\t\t\tif (ow &^ ew) != 0 {
        \t\t\t\t\treturn false
        \t\t\t\t}
        \t\t\t}
        \t\t\treturn true
        \t\t}
        \t}
        \treturn colContainsAll(es.Contains, other)
        }

        func (es *jEnumSetCore) AddAll(other Ext_java_util_Collection) bool {
        \tCheckNotNull(other)
        \tif ep, ok := any(other).(jFastEnumSetProvider); ok {
        \t\tif oes := ep.fastEnumSet_(); oes != nil {
        \t\t\tif oes.size == 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tif es.universe == nil {
        \t\t\t\tes.universe = oes.universe
        \t\t\t} else if !es.sameUniverse(oes) {
        \t\t\t\tpanic(javalang.NewClassCastException("enum type mismatch in EnumSet.addAll"))
        \t\t\t}
        \t\t\toldSize := es.size
        \t\t\tes.bits0 |= oes.bits0
        \t\t\tes.bits1 |= oes.bits1
        \t\t\tif len(oes.more) > 0 {
        \t\t\t\tif len(es.more) < len(oes.more) {
        \t\t\t\t\tnext := make([]uint64, len(oes.more))
        \t\t\t\t\tcopy(next, es.more)
        \t\t\t\t\tes.more = next
        \t\t\t\t}
        \t\t\t\tfor i, ow := range oes.more {
        \t\t\t\t\tes.more[i] |= ow
        \t\t\t\t}
        \t\t\t}
        \t\t\tes.recomputeSize()
        \t\t\treturn es.size != oldSize
        \t\t}
        \t}
        \tchanged := false
        \tfor it := Iterate(other); it.HasNext(); {
        \t\tif es.Add(it.Next()) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func (es *jEnumSetCore) RemoveAll(other Ext_java_util_Collection) bool {
        \tCheckNotNull(other)
        \tif es.size == 0 {
        \t\treturn false
        \t}
        \tif ep, ok := any(other).(jFastEnumSetProvider); ok {
        \t\tif oes := ep.fastEnumSet_(); oes != nil {
        \t\t\tif oes.size == 0 || !es.sameUniverse(oes) {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\toldSize := es.size
        \t\t\tes.bits0 &^= oes.bits0
        \t\t\tes.bits1 &^= oes.bits1
        \t\t\tfor i := 0; i < len(es.more) && i < len(oes.more); i++ {
        \t\t\t\tes.more[i] &^= oes.more[i]
        \t\t\t}
        \t\t\tes.recomputeSize()
        \t\t\treturn es.size != oldSize
        \t\t}
        \t}
        \tchanged := false
        \tfor it := Iterate(other); it.HasNext(); {
        \t\tif es.Remove(it.Next()) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func (es *jEnumSetCore) RetainAll(other Ext_java_util_Collection) bool {
        \tCheckNotNull(other)
        \tif es.size == 0 {
        \t\treturn false
        \t}
        \tif ep, ok := any(other).(jFastEnumSetProvider); ok {
        \t\tif oes := ep.fastEnumSet_(); oes != nil {
        \t\t\tif oes.size == 0 || !es.sameUniverse(oes) {
        \t\t\t\tes.Clear()
        \t\t\t\treturn true
        \t\t\t}
        \t\t\toldSize := es.size
        \t\t\tes.bits0 &= oes.bits0
        \t\t\tes.bits1 &= oes.bits1
        \t\t\tfor i := range es.more {
        \t\t\t\tif i < len(oes.more) {
        \t\t\t\t\tes.more[i] &= oes.more[i]
        \t\t\t\t} else {
        \t\t\t\t\tes.more[i] = 0
        \t\t\t\t}
        \t\t\t}
        \t\t\tes.recomputeSize()
        \t\t\treturn es.size != oldSize
        \t\t}
        \t}
        \tchanged := false
        \tfor ord := es.nextSetBit(0); ord >= 0; ord = es.nextSetBit(ord + 1) {
        \t\tif !other.Contains(es.universe[ord]) {
        \t\t\tes.clearOrd(ord)
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func (es *jEnumSetCore) RemoveIf(pred Ext_java_util_function_Predicate) bool {
        \tCheckNotNull(pred)
        \tif es.size == 0 {
        \t\treturn false
        \t}
        \tchanged := false
        \tfor ord := es.nextSetBit(0); ord >= 0; ord = es.nextSetBit(ord + 1) {
        \t\tif pred.Test(es.universe[ord]) {
        \t\t\tes.clearOrd(ord)
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func (es *jEnumSetCore) Equals(other any) bool {
        \tst, ok := other.(Ext_java_util_Set)
        \tif !ok || st == nil || javalang.IsNil(st) {
        \t\treturn false
        \t}
        \tif ep, ok := other.(jFastEnumSetProvider); ok {
        \t\tif oes := ep.fastEnumSet_(); oes != nil {
        \t\t\tif es == oes {
        \t\t\t\treturn true
        \t\t\t}
        \t\t\tif es.size != oes.size {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tif es.size == 0 {
        \t\t\t\treturn true
        \t\t\t}
        \t\t\tif !es.sameUniverse(oes) {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tif es.bits0 != oes.bits0 || es.bits1 != oes.bits1 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tmaxM := max(len(es.more), len(oes.more))
        \t\t\tfor i := 0; i < maxM; i++ {
        \t\t\t\tvar a, b uint64
        \t\t\t\tif i < len(es.more) {
        \t\t\t\t\ta = es.more[i]
        \t\t\t\t}
        \t\t\t\tif i < len(oes.more) {
        \t\t\t\t\tb = oes.more[i]
        \t\t\t\t}
        \t\t\t\tif a != b {
        \t\t\t\t\treturn false
        \t\t\t\t}
        \t\t\t}
        \t\t\treturn true
        \t\t}
        \t}
        \tif es.size != st.Size() {
        \t\treturn false
        \t}
        \tif es.size == 0 {
        \t\treturn true
        \t}
        \treturn es.ContainsAll(st)
        }

        func (es *jEnumSetCore) HashCode() int32 {
        \tif es.size == 0 {
        \t\treturn 0
        \t}
        \th := int32(0)
        \tfor ord := es.nextSetBit(0); ord >= 0; ord = es.nextSetBit(ord + 1) {
        \t\th += es.universe[ord].(EnumAccessor).GetEnumBase().HashCode()
        \t}
        \treturn h
        }

        func (es *jEnumSetCore) String() string {
        \tif es.size == 0 {
        \t\treturn "[]"
        \t}
        \tvar sb strings.Builder
        \tsb.WriteByte('[')
        \tfirst := true
        \tfor ord := es.nextSetBit(0); ord >= 0; ord = es.nextSetBit(ord + 1) {
        \t\tif !first {
        \t\t\tsb.WriteString(", ")
        \t\t}
        \t\tfirst = false
        \t\tsb.WriteString(StrVal(es.universe[ord]))
        \t}
        \tsb.WriteByte(']')
        \treturn sb.String()
        }

        type jEnumSetIter struct {
        \tes      *jEnumSetCore
        \tnextOrd int
        \tlastOrd int
        }

        func (*jEnumSetIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jEnumSetIter) HasNext() bool                          { return it.nextOrd >= 0 }
        func (it *jEnumSetIter) Next() any {
        \tif it.nextOrd < 0 {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tord := it.nextOrd
        \tit.lastOrd = ord
        \tit.nextOrd = it.es.nextSetBit(ord + 1)
        \treturn it.es.universe[ord]
        }
        func (it *jEnumSetIter) Remove() {
        \tif it.lastOrd < 0 {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \tit.es.clearOrd(it.lastOrd)
        \tit.lastOrd = -1
        }

        func (es *jEnumSetCore) Iterator() Ext_java_util_Iterator {
        \tif es.size == 0 {
        \t\treturn emptyMapIter_
        \t}
        \treturn &jEnumSetIter{es: es, nextOrd: es.nextSetBit(0), lastOrd: -1}
        }

        func enumSetNoneOf(cls *javalang.Class) Ext_java_util_EnumSet {
        \tCheckNotNullPtr(cls)
        \tu := LookupEnumConstants(cls)
        \tobj := newEnumSetShell()
        \tobj.inlineEnumSet.universe = u
        \treturn obj
        }

        func enumSetAllOf(cls *javalang.Class) Ext_java_util_EnumSet {
        \tCheckNotNullPtr(cls)
        \tu := LookupEnumConstants(cls)
        \tif u == nil {
        \t\tpanic(javalang.NewClassCastException(cls.GetName() + " not an enum"))
        \t}
        \tobj := newEnumSetShell()
        \tobj.inlineEnumSet.initAll(u)
        \treturn obj
        }

        func enumSetClone(src *jEnumSetCore) Ext_java_util_EnumSet {
        \tobj := newEnumSetShell()
        \tobj.inlineEnumSet.initCopy(src)
        \treturn obj
        }

        func enumSetCopyOf(src any) Ext_java_util_EnumSet {
        \tCheckNotNull(src)
        \tif ep, ok := src.(jFastEnumSetProvider); ok {
        \t\tif oes := ep.fastEnumSet_(); oes != nil {
        \t\t\treturn enumSetClone(oes)
        \t\t}
        \t}
        \tobj := newEnumSetShell()
        \tes := &obj.inlineEnumSet
        \taddedAny := false
        \tfor it := Iterate(src); it.HasNext(); {
        \t\tes.Add(CheckNotNull(it.Next()))
        \t\taddedAny = true
        \t}
        \tif !addedAny && es.universe == nil {
        \t\tpanic(javalang.NewIllegalArgumentException("Collection is empty"))
        \t}
        \treturn obj
        }

        func enumSetComplementOf(src Ext_java_util_EnumSet) Ext_java_util_EnumSet {
        \tCheckNotNull(src)
        \tobj := newEnumSetShell()
        \tdst := &obj.inlineEnumSet
        \tif ep, ok := any(src).(jFastEnumSetProvider); ok {
        \t\tif ses := ep.fastEnumSet_(); ses != nil && ses.universe != nil {
        \t\t\tdst.initAll(ses.universe)
        \t\t\tdst.bits0 &^= ses.bits0
        \t\t\tdst.bits1 &^= ses.bits1
        \t\t\tfor i := 0; i < len(dst.more) && i < len(ses.more); i++ {
        \t\t\t\tdst.more[i] &^= ses.more[i]
        \t\t\t}
        \t\t\tdst.size -= ses.size
        \t\t\treturn obj
        \t\t}
        \t}
        \treturn obj
        }

        func enumSetRange(from, to any) Ext_java_util_EnumSet {
        \tCheckNotNull(from)
        \tCheckNotNull(to)
        \tobj := newEnumSetShell()
        \tes := &obj.inlineEnumSet
        \tfOrd := es.mustOrd(from)
        \ttOrd := es.mustOrd(to)
        \tif fOrd > tOrd {
        \t\tpanic(javalang.NewIllegalArgumentException(StrVal(from) + " > " + StrVal(to)))
        \t}
        \tfor ord := fOrd; ord <= tOrd; ord++ {
        \t\tes.setOrd(ord)
        \t}
        \treturn obj
        }

        func immSetFromEnumSet(es *jEnumSetCore) Ext_com_google_common_collect_ImmutableSet {
        \tif es == nil || es.size == 0 {
        \t\treturn emptyImmutableSet_
        \t}
        \tif es.size == 1 {
        \t\tord := es.nextSetBit(0)
        \t\treturn newImmSet1(es.universe[ord])
        \t}
        \telems := make([]any, 0, es.size)
        \tfor ord := es.nextSetBit(0); ord >= 0; ord = es.nextSetBit(ord + 1) {
        \t\telems = append(elems, es.universe[ord])
        \t}
        \treturn finishImmSetMulti(elems, nil)
        }

        func setsImmutableEnumSetIterable(src Ext_java_lang_Iterable) Ext_com_google_common_collect_ImmutableSet {
        \tCheckNotNull(src)
        \tif ep, ok := any(src).(jFastEnumSetProvider); ok {
        \t\tif es := ep.fastEnumSet_(); es != nil {
        \t\t\treturn immSetFromEnumSet(es)
        \t\t}
        \t}
        \tvar es jEnumSetCore
        \tfor it := Iterate(src); it.HasNext(); {
        \t\tes.Add(CheckNotNull(it.Next()))
        \t}
        \treturn immSetFromEnumSet(&es)
        }

        type jEnumMapCore struct {
        \tuniverse   []any
        \tvals       []any
        \tinlineVals [4]any
        \tbits0      uint64
        \tbits1      uint64
        \tmore       []uint64
        \tsize       int32
        }

        func (em *jEnumMapCore) Size() int32 { return em.size }

        func (em *jEnumMapCore) initClass(cls *javalang.Class) {
        \tCheckNotNullPtr(cls)
        \tem.universe = LookupEnumConstants(cls)
        }

        func (em *jEnumMapCore) ensureVals() {
        \tif em.vals != nil {
        \t\treturn
        \t}
        \tn := len(em.universe)
        \tif n <= len(em.inlineVals) {
        \t\tem.vals = em.inlineVals[:n]
        \t} else {
        \t\tem.vals = make([]any, n)
        \t}
        }

        func (em *jEnumMapCore) initCopyCore(src *jEnumMapCore) {
        \tem.universe = src.universe
        \tem.bits0 = src.bits0
        \tem.bits1 = src.bits1
        \tif len(src.more) > 0 {
        \t\tem.more = make([]uint64, len(src.more))
        \t\tcopy(em.more, src.more)
        \t}
        \tem.size = src.size
        \tif src.size > 0 && src.vals != nil {
        \t\tem.ensureVals()
        \t\tcopy(em.vals, src.vals)
        \t}
        }

        func (em *jEnumMapCore) initCopyEnumMap(src any) {
        \tCheckNotNull(src)
        \tif ep, ok := src.(jFastEnumMapProvider); ok {
        \t\tif sem := ep.fastEnumMap_(); sem != nil {
        \t\t\tem.initCopyCore(sem)
        \t\t\treturn
        \t\t}
        \t}
        \tif m, ok := src.(Ext_java_util_Map); ok {
        \t\tem.initCopyMap(m)
        \t}
        }

        func (em *jEnumMapCore) initCopyMap(src Ext_java_util_Map) {
        \tCheckNotNull(src)
        \tif ep, ok := any(src).(jFastEnumMapProvider); ok {
        \t\tif sem := ep.fastEnumMap_(); sem != nil {
        \t\t\tem.initCopyCore(sem)
        \t\t\treturn
        \t\t}
        \t}
        \tif src.IsEmpty() {
        \t\tpanic(javalang.NewIllegalArgumentException("Specified map is empty"))
        \t}
        \tem.PutAll(src)
        }

        func (em *jEnumMapCore) hasBit(ord int) bool {
        \tw := ord >> 6
        \tmask := uint64(1) << uint(ord&63)
        \tswitch w {
        \tcase 0:
        \t\treturn (em.bits0 & mask) != 0
        \tcase 1:
        \t\treturn (em.bits1 & mask) != 0
        \tdefault:
        \t\tmw := w - 2
        \t\treturn mw < len(em.more) && (em.more[mw]&mask) != 0
        \t}
        }
        """;
  }

  private static String buildColRtPart1c() {
    return """

        func (em *jEnumMapCore) setBit(ord int) bool {
        \tw := ord >> 6
        \tmask := uint64(1) << uint(ord&63)
        \tswitch w {
        \tcase 0:
        \t\tif (em.bits0 & mask) != 0 {
        \t\t\treturn false
        \t\t}
        \t\tem.bits0 |= mask
        \t\tem.size++
        \t\treturn true
        \tcase 1:
        \t\tif (em.bits1 & mask) != 0 {
        \t\t\treturn false
        \t\t}
        \t\tem.bits1 |= mask
        \t\tem.size++
        \t\treturn true
        \tdefault:
        \t\tmw := w - 2
        \t\tif mw >= len(em.more) {
        \t\t\tneed := ((len(em.universe) + 63) >> 6) - 2
        \t\t\tif mw+1 > need {
        \t\t\t\tneed = mw + 1
        \t\t\t}
        \t\t\tnext := make([]uint64, need)
        \t\t\tcopy(next, em.more)
        \t\t\tem.more = next
        \t\t}
        \t\tif (em.more[mw] & mask) != 0 {
        \t\t\treturn false
        \t\t}
        \t\tem.more[mw] |= mask
        \t\tem.size++
        \t\treturn true
        \t}
        }

        func (em *jEnumMapCore) clearBit(ord int) bool {
        \tw := ord >> 6
        \tmask := uint64(1) << uint(ord&63)
        \tswitch w {
        \tcase 0:
        \t\tif (em.bits0 & mask) == 0 {
        \t\t\treturn false
        \t\t}
        \t\tem.bits0 &^= mask
        \t\tem.size--
        \t\tif ord < len(em.vals) {
        \t\t\tem.vals[ord] = nil
        \t\t}
        \t\treturn true
        \tcase 1:
        \t\tif (em.bits1 & mask) == 0 {
        \t\t\treturn false
        \t\t}
        \t\tem.bits1 &^= mask
        \t\tem.size--
        \t\tif ord < len(em.vals) {
        \t\t\tem.vals[ord] = nil
        \t\t}
        \t\treturn true
        \tdefault:
        \t\tmw := w - 2
        \t\tif mw >= len(em.more) || (em.more[mw]&mask) == 0 {
        \t\t\treturn false
        \t\t}
        \t\tem.more[mw] &^= mask
        \t\tem.size--
        \t\tif ord < len(em.vals) {
        \t\t\tem.vals[ord] = nil
        \t\t}
        \t\treturn true
        \t}
        }

        func (em *jEnumMapCore) nextSetBit(fromOrd int) int {
        \tif fromOrd < 0 {
        \t\tfromOrd = 0
        \t}
        \tw := fromOrd >> 6
        \tshift := uint(fromOrd & 63)
        \tif w == 0 {
        \t\tword := em.bits0 & (^uint64(0) << shift)
        \t\tif word != 0 {
        \t\t\treturn bits.TrailingZeros64(word)
        \t\t}
        \t\tw = 1
        \t\tshift = 0
        \t}
        \tif w == 1 {
        \t\tword := em.bits1 & (^uint64(0) << shift)
        \t\tif word != 0 {
        \t\t\treturn 64 + bits.TrailingZeros64(word)
        \t\t}
        \t\tw = 2
        \t\tshift = 0
        \t}
        \tfor mw := w - 2; mw < len(em.more); mw++ {
        \t\tword := em.more[mw] & (^uint64(0) << shift)
        \t\tshift = 0
        \t\tif word != 0 {
        \t\t\treturn ((mw + 2) << 6) + bits.TrailingZeros64(word)
        \t\t}
        \t}
        \treturn -1
        }

        func (em *jEnumMapCore) ordIfValid(k any) (int, bool) {
        \tk = ToAny(k)
        \tif k == nil || em.universe == nil {
        \t\treturn -1, false
        \t}
        \tea, ok := k.(EnumAccessor)
        \tif !ok {
        \t\treturn -1, false
        \t}
        \teb := ea.GetEnumBase()
        \tif eb == nil {
        \t\treturn -1, false
        \t}
        \tord := int(eb.Ordinal_)
        \tif ord >= 0 && ord < len(em.universe) && em.universe[ord] == k {
        \t\treturn ord, true
        \t}
        \treturn -1, false
        }

        func (em *jEnumMapCore) mustOrd(k any) int {
        \tk = CheckNotNull(ToAny(k))
        \tea, ok := k.(EnumAccessor)
        \tif !ok {
        \t\tpanic(javalang.NewClassCastException(fmt.Sprintf("%T cannot be cast to java.lang.Enum", k)))
        \t}
        \teb := ea.GetEnumBase()
        \tif eb == nil {
        \t\tThrowNPE()
        \t}
        \tif em.universe == nil {
        \t\tem.universe = resolveEnumUniverseFromElem(k)
        \t}
        \tord := int(eb.Ordinal_)
        \tif ord >= 0 && ord < len(em.universe) && em.universe[ord] == k {
        \t\treturn ord
        \t}
        \tpanic(javalang.NewClassCastException("enum key type mismatch in EnumMap"))
        }

        func (em *jEnumMapCore) ContainsKey(k any) bool {
        \tif em.size == 0 {
        \t\treturn false
        \t}
        \tord, ok := em.ordIfValid(k)
        \treturn ok && em.hasBit(ord)
        }

        func (em *jEnumMapCore) ContainsValue(v any) bool {
        \tif em.size == 0 {
        \t\treturn false
        \t}
        \tv = ToAny(v)
        \tfor ord := em.nextSetBit(0); ord >= 0; ord = em.nextSetBit(ord + 1) {
        \t\tif javalang.ObjectsEquals(em.vals[ord], v) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (em *jEnumMapCore) Get(k any) any {
        \tif em.size == 0 {
        \t\treturn nil
        \t}
        \tord, ok := em.ordIfValid(k)
        \tif !ok || !em.hasBit(ord) {
        \t\treturn nil
        \t}
        \treturn em.vals[ord]
        }

        func (em *jEnumMapCore) GetOrDefault(k, def any) any {
        \tif em.size > 0 {
        \t\tif ord, ok := em.ordIfValid(k); ok && em.hasBit(ord) {
        \t\t\treturn em.vals[ord]
        \t\t}
        \t}
        \treturn ToAny(def)
        }

        func (em *jEnumMapCore) Put(k, v any) any {
        \tord := em.mustOrd(k)
        \tem.ensureVals()
        \tv = ToAny(v)
        \tif em.hasBit(ord) {
        \t\told := em.vals[ord]
        \t\tem.vals[ord] = v
        \t\treturn old
        \t}
        \tem.setBit(ord)
        \tem.vals[ord] = v
        \treturn nil
        }

        func (em *jEnumMapCore) PutIfAbsent(k, v any) any {
        \tord := em.mustOrd(k)
        \tem.ensureVals()
        \tif em.hasBit(ord) {
        \t\tif cur := em.vals[ord]; cur != nil {
        \t\t\treturn cur
        \t\t}
        \t\tem.vals[ord] = ToAny(v)
        \t\treturn nil
        \t}
        \tem.setBit(ord)
        \tem.vals[ord] = ToAny(v)
        \treturn nil
        }

        func (em *jEnumMapCore) Remove(k any) any {
        \tif em.size == 0 {
        \t\treturn nil
        \t}
        \tord, ok := em.ordIfValid(k)
        \tif !ok || !em.hasBit(ord) {
        \t\treturn nil
        \t}
        \told := em.vals[ord]
        \tem.clearBit(ord)
        \treturn old
        }

        func (em *jEnumMapCore) Remove2(k, v any) bool {
        \tif em.size == 0 {
        \t\treturn false
        \t}
        \tord, ok := em.ordIfValid(k)
        \tif !ok || !em.hasBit(ord) || !javalang.ObjectsEquals(em.vals[ord], ToAny(v)) {
        \t\treturn false
        \t}
        \tem.clearBit(ord)
        \treturn true
        }

        func (em *jEnumMapCore) PutAll(src Ext_java_util_Map) {
        \tCheckNotNull(src)
        \tif ep, ok := any(src).(jFastEnumMapProvider); ok {
        \t\tif sem := ep.fastEnumMap_(); sem != nil {
        \t\t\tif sem.size == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tsem.ForEachEntry(func(k, v any) bool {
        \t\t\t\tem.Put(k, v)
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn
        \t\t}
        \t}
        \tif mp, ok := any(src).(jFastMapCoreProvider); ok {
        \t\tif smc := mp.fastMapCore_(); smc != nil {
        \t\t\tsmc.ForEachEntry(func(k, v any) bool {
        \t\t\t\tem.Put(k, v)
        \t\t\t\treturn true
        \t\t\t})
        \t\t\treturn
        \t\t}
        \t}
        \tif es := src.EntrySet(); es != nil && !javalang.IsNil(es) {
        \t\tfor it := Iterate(es); it.HasNext(); {
        \t\t\te := Cast[Ext_java_util_Map_Entry](it.Next())
        \t\t\tem.Put(e.GetKey(), e.GetValue())
        \t\t}
        \t}
        }

        func (em *jEnumMapCore) Clear() {
        \tem.bits0 = 0
        \tem.bits1 = 0
        \tif len(em.more) > 0 {
        \t\tclear(em.more)
        \t}
        \tif len(em.vals) > 0 {
        \t\tclear(em.vals)
        \t}
        \tem.size = 0
        }

        func (em *jEnumMapCore) Replace2(k, newVal any) any {
        \tif em.size == 0 {
        \t\treturn nil
        \t}
        \tord, ok := em.ordIfValid(k)
        \tif !ok || !em.hasBit(ord) {
        \t\treturn nil
        \t}
        \told := em.vals[ord]
        \tem.vals[ord] = ToAny(newVal)
        \treturn old
        }

        func (em *jEnumMapCore) Replace3(k, oldVal, newVal any) bool {
        \tif em.size == 0 {
        \t\treturn false
        \t}
        \tord, ok := em.ordIfValid(k)
        \tif !ok || !em.hasBit(ord) || !javalang.ObjectsEquals(em.vals[ord], ToAny(oldVal)) {
        \t\treturn false
        \t}
        \tem.vals[ord] = ToAny(newVal)
        \treturn true
        }

        func (em *jEnumMapCore) ReplaceAll(fn Ext_java_util_function_BiFunction) {
        \tCheckNotNull(fn)
        \tfor ord := em.nextSetBit(0); ord >= 0; ord = em.nextSetBit(ord + 1) {
        \t\tem.vals[ord] = ToAny(fn.Apply(em.universe[ord], em.vals[ord]))
        \t}
        }

        func (em *jEnumMapCore) tryGetNonNull(k any) (any, bool) {
        \tord, ok := em.ordIfValid(k)
        \tif ok && em.hasBit(ord) {
        \t\tif cur := em.vals[ord]; cur != nil {
        \t\t\treturn cur, true
        \t\t}
        \t}
        \treturn nil, false
        }

        func (em *jEnumMapCore) ComputeIfAbsent(k any, fn Ext_java_util_function_Function) any {
        \tCheckNotNull(fn)
        \tord := em.mustOrd(k)
        \tif em.hasBit(ord) {
        \t\tif cur := em.vals[ord]; cur != nil {
        \t\t\treturn cur
        \t\t}
        \t}
        \tnewVal := ToAny(fn.Apply(ToAny(k)))
        \tif newVal != nil {
        \t\tem.ensureVals()
        \t\tem.setBit(ord)
        \t\tem.vals[ord] = newVal
        \t}
        \treturn newVal
        }

        func (em *jEnumMapCore) ComputeIfPresent(k any, fn Ext_java_util_function_BiFunction) any {
        \tCheckNotNull(fn)
        \tord, ok := em.ordIfValid(k)
        \tif !ok || !em.hasBit(ord) || em.vals[ord] == nil {
        \t\treturn nil
        \t}
        \tnewVal := ToAny(fn.Apply(ToAny(k), em.vals[ord]))
        \tif newVal == nil {
        \t\tem.clearBit(ord)
        \t\treturn nil
        \t}
        \tem.vals[ord] = newVal
        \treturn newVal
        }

        func (em *jEnumMapCore) Compute(k any, fn Ext_java_util_function_BiFunction) any {
        \tCheckNotNull(fn)
        \tord := em.mustOrd(k)
        \thaveOld := em.hasBit(ord)
        \tvar oldVal any
        \tif haveOld {
        \t\toldVal = em.vals[ord]
        \t}
        \tnewVal := ToAny(fn.Apply(ToAny(k), oldVal))
        \tif newVal == nil {
        \t\tif haveOld {
        \t\t\tem.clearBit(ord)
        \t\t}
        \t\treturn nil
        \t}
        \tem.ensureVals()
        \tem.setBit(ord)
        \tem.vals[ord] = newVal
        \treturn newVal
        }

        func (em *jEnumMapCore) Merge(k, val any, fn Ext_java_util_function_BiFunction) any {
        \tCheckNotNull(fn)
        \tv := CheckNotNull(ToAny(val))
        \tord := em.mustOrd(k)
        \thaveOld := em.hasBit(ord)
        \tvar oldVal any
        \tif haveOld {
        \t\toldVal = em.vals[ord]
        \t}
        \tvar newVal any
        \tif oldVal == nil {
        \t\tnewVal = v
        \t} else {
        \t\tnewVal = ToAny(fn.Apply(oldVal, v))
        \t}
        \tif newVal == nil {
        \t\tif haveOld {
        \t\t\tem.clearBit(ord)
        \t\t}
        \t\treturn nil
        \t}
        \tem.ensureVals()
        \tem.setBit(ord)
        \tem.vals[ord] = newVal
        \treturn newVal
        }

        func (em *jEnumMapCore) ForEachEntry(fn func(k, v any) bool) {
        \tif em.size == 0 {
        \t\treturn
        \t}
        \tfor ord := em.nextSetBit(0); ord >= 0; ord = em.nextSetBit(ord + 1) {
        \t\tif !fn(em.universe[ord], em.vals[ord]) {
        \t\t\treturn
        \t\t}
        \t}
        }

        func (em *jEnumMapCore) ForEachBiConsumer(action Ext_java_util_function_BiConsumer) {
        \tCheckNotNull(action)
        \tfor ord := em.nextSetBit(0); ord >= 0; ord = em.nextSetBit(ord + 1) {
        \t\taction.Accept(em.universe[ord], em.vals[ord])
        \t}
        }

        func (em *jEnumMapCore) Equals(other any) bool {
        \tmp, ok := other.(Ext_java_util_Map)
        \tif !ok || mp == nil || isJavaNilAny(mp) {
        \t\treturn false
        \t}
        \tif em.size != mp.Size() {
        \t\treturn false
        \t}
        \tif em.size == 0 {
        \t\treturn true
        \t}
        \tfor ord := em.nextSetBit(0); ord >= 0; ord = em.nextSetBit(ord + 1) {
        \t\tk := em.universe[ord]
        \t\tv := em.vals[ord]
        \t\tif v == nil {
        \t\t\tif mp.Get(k) != nil || !mp.ContainsKey(k) {
        \t\t\t\treturn false
        \t\t\t}
        \t\t} else if !javalang.ObjectsEquals(v, ToAny(mp.Get(k))) {
        \t\t\treturn false
        \t\t}
        \t}
        \treturn true
        }

        func (em *jEnumMapCore) HashCode() int32 {
        \tif em.size == 0 {
        \t\treturn 0
        \t}
        \th := int32(0)
        \tfor ord := em.nextSetBit(0); ord >= 0; ord = em.nextSetBit(ord + 1) {
        \t\th += em.universe[ord].(EnumAccessor).GetEnumBase().HashCode() ^ javalang.ObjectHashCode(em.vals[ord])
        \t}
        \treturn h
        }

        func (em *jEnumMapCore) String() string {
        \tif em.size == 0 {
        \t\treturn "{}"
        \t}
        \tvar sb strings.Builder
        \tsb.WriteByte('{')
        \tfirst := true
        \tfor ord := em.nextSetBit(0); ord >= 0; ord = em.nextSetBit(ord + 1) {
        \t\tif !first {
        \t\t\tsb.WriteString(", ")
        \t\t}
        \t\tfirst = false
        \t\tsb.WriteString(StrVal(em.universe[ord]))
        \t\tsb.WriteByte('=')
        \t\tsb.WriteString(StrVal(em.vals[ord]))
        \t}
        \tsb.WriteByte('}')
        \treturn sb.String()
        }

        type jLiveEnumMapKeySet struct {
        \tExt_java_util_AbstractSetBase
        \tem *jEnumMapCore
        }

        func newLiveEnumMapKeySet(em *jEnumMapCore) Ext_java_util_Set {
        \ts := &jLiveEnumMapKeySet{em: em}
        \ts.BindExtSelf_(s)
        \treturn s
        }
        func (s *jLiveEnumMapKeySet) Size() int32         { return s.em.size }
        func (s *jLiveEnumMapKeySet) Contains(k any) bool { return s.em.ContainsKey(k) }
        func (s *jLiveEnumMapKeySet) Remove_1_any(k any) bool {
        \tif s.em.size == 0 {
        \t\treturn false
        \t}
        \tord, ok := s.em.ordIfValid(k)
        \treturn ok && s.em.clearBit(ord)
        }
        func (s *jLiveEnumMapKeySet) Clear() { s.em.Clear() }
        func (s *jLiveEnumMapKeySet) Iterator() any {
        \tif s.em == nil || s.em.size == 0 {
        \t\treturn emptyMapIter_
        \t}
        \treturn &jEnumMapKeyIter{em: s.em, nextOrd: s.em.nextSetBit(0), lastOrd: -1}
        }
        func (s *jLiveEnumMapKeySet) SetContainsStr_(_ string) bool {
        \tif s == nil { ThrowNPE() }
        \treturn false
        }
        func (s *jLiveEnumMapKeySet) SetAddStr_(_ string) bool {
        \tif s == nil { ThrowNPE() }
        \tpanic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
        }
        func (s *jLiveEnumMapKeySet) SetRemoveStr_(_ string) bool {
        \tif s == nil { ThrowNPE() }
        \treturn false
        }
        func (s *jLiveEnumMapKeySet) SetContainsNode_(_ uint32) bool {
        \tif s == nil { ThrowNPE() }
        \treturn false
        }
        func (s *jLiveEnumMapKeySet) SetAddNode_(_ uint32) bool {
        \tif s == nil { ThrowNPE() }
        \tpanic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
        }
        func (s *jLiveEnumMapKeySet) SetRemoveNode_(_ uint32) bool {
        \tif s == nil { ThrowNPE() }
        \treturn false
        }

        type jEnumMapKeyIter struct {
        \tem      *jEnumMapCore
        \tnextOrd int
        \tlastOrd int
        }

        func (*jEnumMapKeyIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jEnumMapKeyIter) HasNext() bool                          { return it.nextOrd >= 0 }
        func (it *jEnumMapKeyIter) Next() any {
        \tif it.nextOrd < 0 {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tord := it.nextOrd
        \tit.lastOrd = ord
        \tit.nextOrd = it.em.nextSetBit(ord + 1)
        \treturn it.em.universe[ord]
        }
        func (it *jEnumMapKeyIter) Remove() {
        \tif it.lastOrd < 0 {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \tit.em.clearBit(it.lastOrd)
        \tit.lastOrd = -1
        }

        type jLiveEnumMapValues struct {
        \tExt_java_util_AbstractCollectionBase
        \tem *jEnumMapCore
        }

        func newLiveEnumMapValues(em *jEnumMapCore) Ext_java_util_Collection {
        \tv := &jLiveEnumMapValues{em: em}
        \tv.BindExtSelf_(v)
        \treturn v
        }
        func (v *jLiveEnumMapValues) Size() int32         { return v.em.size }
        func (v *jLiveEnumMapValues) Contains(x any) bool { return v.em.ContainsValue(x) }
        func (v *jLiveEnumMapValues) Clear()              { v.em.Clear() }
        func (v *jLiveEnumMapValues) Iterator() any {
        \tif v.em == nil || v.em.size == 0 {
        \t\treturn emptyMapIter_
        \t}
        \treturn &jEnumMapValIter{em: v.em, nextOrd: v.em.nextSetBit(0), lastOrd: -1}
        }

        type jEnumMapValIter struct {
        \tem      *jEnumMapCore
        \tnextOrd int
        \tlastOrd int
        }

        func (*jEnumMapValIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jEnumMapValIter) HasNext() bool                          { return it.nextOrd >= 0 }
        func (it *jEnumMapValIter) Next() any {
        \tif it.nextOrd < 0 {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tord := it.nextOrd
        \tit.lastOrd = ord
        \tit.nextOrd = it.em.nextSetBit(ord + 1)
        \treturn it.em.vals[ord]
        }
        func (it *jEnumMapValIter) Remove() {
        \tif it.lastOrd < 0 {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \tit.em.clearBit(it.lastOrd)
        \tit.lastOrd = -1
        }

        type jLiveEnumMapEntrySet struct {
        \tExt_java_util_AbstractSetBase
        \tem *jEnumMapCore
        }

        func newLiveEnumMapEntrySet(em *jEnumMapCore) Ext_java_util_Set {
        \ts := &jLiveEnumMapEntrySet{em: em}
        \ts.BindExtSelf_(s)
        \treturn s
        }
        func (s *jLiveEnumMapEntrySet) Size() int32 { return s.em.size }
        func (s *jLiveEnumMapEntrySet) Clear()      { s.em.Clear() }
        func (s *jLiveEnumMapEntrySet) Contains(o any) bool {
        \te, ok := o.(Ext_java_util_Map_Entry)
        \tif !ok || e == nil || javalang.IsNil(e) {
        \t\treturn false
        \t}
        \tord, ok := s.em.ordIfValid(e.GetKey())
        \treturn ok && s.em.hasBit(ord) && javalang.ObjectsEquals(s.em.vals[ord], ToAny(e.GetValue()))
        }
        func (s *jLiveEnumMapEntrySet) Remove_1_any(o any) bool {
        \te, ok := o.(Ext_java_util_Map_Entry)
        \tif !ok || e == nil || javalang.IsNil(e) {
        \t\treturn false
        \t}
        \tord, ok := s.em.ordIfValid(e.GetKey())
        \tif !ok || !s.em.hasBit(ord) || !javalang.ObjectsEquals(s.em.vals[ord], ToAny(e.GetValue())) {
        \t\treturn false
        \t}
        \ts.em.clearBit(ord)
        \treturn true
        }
        func (s *jLiveEnumMapEntrySet) Iterator() any {
        \tif s.em == nil || s.em.size == 0 {
        \t\treturn emptyMapIter_
        \t}
        \treturn &jEnumMapEntryIter{em: s.em, nextOrd: s.em.nextSetBit(0), lastOrd: -1}
        }

        type jEnumMapEntry struct {
        \tem  *jEnumMapCore
        \tord int
        \tk   any
        \tv   any
        }

        func (*jEnumMapEntry) Implements_j2gort_Ext_java_util_Map_Entry() {}
        func (e *jEnumMapEntry) GetKey() any                              { return e.k }
        func (e *jEnumMapEntry) GetValue() any {
        \tif e.em.hasBit(e.ord) {
        \t\treturn e.em.vals[e.ord]
        \t}
        \treturn e.v
        }
        func (e *jEnumMapEntry) SetValue(nv any) any {
        \tnv = ToAny(nv)
        \tif e.em.hasBit(e.ord) {
        \t\told := e.em.vals[e.ord]
        \t\te.em.vals[e.ord] = nv
        \t\te.v = nv
        \t\treturn old
        \t}
        \told := e.v
        \te.v = nv
        \treturn old
        }
        func (e *jEnumMapEntry) Equals(o any) bool { return mapEntryEquals(e, o) }
        func (e *jEnumMapEntry) HashCode() int32   { return mapEntryHash(e) }
        func (e *jEnumMapEntry) String() string    { return StrVal(e.GetKey()) + "=" + StrVal(e.GetValue()) }

        type jEnumMapEntryIter struct {
        \tem      *jEnumMapCore
        \tnextOrd int
        \tlastOrd int
        }

        func (*jEnumMapEntryIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jEnumMapEntryIter) HasNext() bool                          { return it.nextOrd >= 0 }
        func (it *jEnumMapEntryIter) Next() any {
        \tif it.nextOrd < 0 {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tord := it.nextOrd
        \tit.lastOrd = ord
        \tit.nextOrd = it.em.nextSetBit(ord + 1)
        \treturn &jEnumMapEntry{em: it.em, ord: ord, k: it.em.universe[ord], v: it.em.vals[ord]}
        }
        func (it *jEnumMapEntryIter) Remove() {
        \tif it.lastOrd < 0 {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \tit.em.clearBit(it.lastOrd)
        \tit.lastOrd = -1
        }

        func mapsImmutableEnumMap(src Ext_java_util_Map) Ext_com_google_common_collect_ImmutableMap {
        \tCheckNotNull(src)
        \tif src.IsEmpty() {
        \t\treturn emptyImmutableMap_
        \t}
        \tif ep, ok := any(src).(jFastEnumMapProvider); ok {
        \t\tif sem := ep.fastEnumMap_(); sem != nil {
        \t\t\tobj := newImmutableMapShell()
        \t\t\tmc := &obj.inlineCore
        \t\t\tmc.InitExpectedSize(sem.size)
        \t\t\tsem.ForEachEntry(func(k, v any) bool {
        \t\t\t\tmc.Put(CheckNotNull(k), CheckNotNull(v))
        \t\t\t\treturn true
        \t\t\t})
        \t\t\tmc.sealImmutable()
        \t\t\treturn obj
        \t\t}
        \t}
        \tvar tmp jEnumMapCore
        \tif es := src.EntrySet(); es != nil && !javalang.IsNil(es) {
        \t\tfor it := Iterate(es); it.HasNext(); {
        \t\t\te := Cast[Ext_java_util_Map_Entry](it.Next())
        \t\t\ttmp.Put(CheckNotNull(ToAny(e.GetKey())), CheckNotNull(ToAny(e.GetValue())))
        \t\t}
        \t}
        \tobj := newImmutableMapShell()
        \tmc := &obj.inlineCore
        \tmc.InitExpectedSize(tmp.size)
        \ttmp.ForEachEntry(func(k, v any) bool {
        \t\tmc.Put(k, v)
        \t\treturn true
        \t})
        \tmc.sealImmutable()
        \treturn obj
        }

        // --- Abstract Delegation Helpers ---

        func iterToSliceDirect(col Ext_java_util_Collection) []any {
        \tif col == nil || javalang.IsNil(col) { ThrowNPE() }
        \tout := make([]any, 0, col.Size())
        \tit := Cast[Ext_java_util_Iterator](col.Iterator())
        \tfor it.HasNext() {
        \t\tout = append(out, ToAny(it.Next()))
        \t}
        \treturn out
        }

        func colToSlice(col Ext_java_util_Collection) []any {
        \tif tv, ok := any(col).(jTypedViewCursorSource); ok {
        \t\tif n, ok2 := tv.typedViewLen_(); ok2 {
        \t\t\tout := make([]any, n)
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\tout[i] = ToAny(tv.typedViewAt_(i))
        \t\t\t}
        \t\t\treturn out
        \t\t}
        \t}
        \tif si, ok := any(col).(SliceIterable); ok {
        \t\tsrc := si.IterSlice_()
        \t\tout := make([]any, len(src))
        \t\tfor i, x := range src {
        \t\t\tout[i] = ToAny(x)
        \t\t}
        \t\treturn out
        \t}
        \tout := make([]any, 0, col.Size())
        \tfor it := NewIterCursor(col); it.HasNext(); {
        \t\tout = append(out, ToAny(it.Next()))
        \t}
        \treturn out
        }

        func abstractColContains(col Ext_java_util_Collection, target any) bool {
        \tt := ToAny(target)
        \tfor it := Iterate(col); it.HasNext(); {
        \t\tif javalang.ObjectsEquals(t, ToAny(it.Next())) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func abstractColRemove(col Ext_java_util_Collection, target any) bool {
        \tt := ToAny(target)
        \tit := Cast[Ext_java_util_Iterator](col.Iterator())
        \tfor it.HasNext() {
        \t\tif javalang.ObjectsEquals(t, ToAny(it.Next())) {
        \t\t\tit.Remove()
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func abstractColAddAll(col Ext_java_util_Collection, other Ext_java_util_Collection) bool {
        \tCheckNotNull(other)
        \tmodified := false
        \tfor it := Iterate(other); it.HasNext(); {
        \t\tif col.Add_1(it.Next()) {
        \t\t\tmodified = true
        \t\t}
        \t}
        \treturn modified
        }

        func abstractColRemoveAll(col Ext_java_util_Collection, other Ext_java_util_Collection) bool {
        \tCheckNotNull(other)
        \tmodified := false
        \tit := Cast[Ext_java_util_Iterator](col.Iterator())
        \tfor it.HasNext() {
        \t\tif other.Contains(it.Next()) {
        \t\t\tit.Remove()
        \t\t\tmodified = true
        \t\t}
        \t}
        \treturn modified
        }

        func abstractColRetainAll(col Ext_java_util_Collection, other Ext_java_util_Collection) bool {
        \tCheckNotNull(other)
        \tmodified := false
        \tit := Cast[Ext_java_util_Iterator](col.Iterator())
        \tfor it.HasNext() {
        \t\tif !other.Contains(it.Next()) {
        \t\t\tit.Remove()
        \t\t\tmodified = true
        \t\t}
        \t}
        \treturn modified
        }

        func abstractColRemoveIf(col Ext_java_util_Collection, pred Ext_java_util_function_Predicate) bool {
        \tCheckNotNull(pred)
        \tmodified := false
        \tit := Cast[Ext_java_util_Iterator](col.Iterator())
        \tfor it.HasNext() {
        \t\tif pred.Test(it.Next()) {
        \t\t\tit.Remove()
        \t\t\tmodified = true
        \t\t}
        \t}
        \treturn modified
        }

        func abstractColClear(col Ext_java_util_Collection) {
        \tit := Cast[Ext_java_util_Iterator](col.Iterator())
        \tfor it.HasNext() {
        \t\t_ = it.Next()
        \t\tit.Remove()
        \t}
        }

        func abstractSetEquals(self Ext_java_util_Collection, other any) bool {
        \tif RefEquals(self, other) {
        \t\treturn true
        \t}
        \tst, ok := other.(Ext_java_util_Set)
        \tif !ok || st == nil || javalang.IsNil(st) {
        \t\treturn false
        \t}
        \tif self.Size() != st.Size() {
        \t\treturn false
        \t}
        \treturn self.ContainsAll(st)
        }

        func abstractSetHash(self Ext_java_util_Collection) int32 {
        \th := int32(0)
        \tfor it := Iterate(self); it.HasNext(); {
        \t\th += javalang.ObjectHashCode(ToAny(it.Next()))
        \t}
        \treturn h
        }

        // AbstractList delegation
        type jAbstractListIter struct {
        \tlst     Ext_java_util_List
        \tcursor  int32
        \tlastRet int32
        }

        func (*jAbstractListIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (*jAbstractListIter) Implements_j2gort_Ext_java_util_ListIterator() {}
        func (it *jAbstractListIter) HasNext() bool { return it.cursor < it.lst.Size() }
        func (it *jAbstractListIter) Next() any {
        \tif it.cursor >= it.lst.Size() {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tv := it.lst.Get(it.cursor)
        \tit.lastRet = it.cursor
        \tit.cursor++
        \treturn v
        }
        func (it *jAbstractListIter) Remove() {
        \tif it.lastRet < 0 {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \tit.lst.Remove_1_int32(it.lastRet)
        \tif it.lastRet < it.cursor {
        \t\tit.cursor--
        \t}
        \tit.lastRet = -1
        }

        func abstractListOp_size_0(l Ext_java_util_List) int32 {
        \tpanic(javalang.NewUnsupportedOperationException("AbstractList.size"))
        }
        func abstractListOp_get_1(l Ext_java_util_List, i int32) any {
        \tpanic(javalang.NewUnsupportedOperationException("AbstractList.get"))
        }
        func abstractListOp_set_2(l Ext_java_util_List, i int32, v any) any {
        \tpanic(javalang.NewUnsupportedOperationException("AbstractList.set"))
        }
        func abstractListOp_add_2(l Ext_java_util_List, i int32, v any) {
        \tpanic(javalang.NewUnsupportedOperationException("AbstractList.add"))
        }
        func abstractListOp_remove_1(l Ext_java_util_List, arg any) any {
        \tif idx, ok := arg.(int32); ok {
        \t\t_ = idx
        \t\tpanic(javalang.NewUnsupportedOperationException("AbstractList.remove(int)"))
        \t}
        \treturn abstractColRemove(l, arg)
        }
        func abstractListOp_add_1(l Ext_java_util_List, v any) bool {
        \tl.Add_2(l.Size(), v)
        \treturn true
        }
        func abstractListOp_isEmpty_0(l Ext_java_util_List) bool       { return l.Size() == 0 }
        func abstractListOp_contains_1(l Ext_java_util_List, v any) bool { return l.IndexOf_1(v) >= 0 }
        func abstractListOp_indexOf_1(l Ext_java_util_List, v any) int32 {
        \tn := l.Size()
        \tt := ToAny(v)
        \tfor i := int32(0); i < n; i++ {
        \t\tif javalang.ObjectsEquals(t, ToAny(l.Get(i))) {
        \t\t\treturn i
        \t\t}
        \t}
        \treturn -1
        }
        func abstractListOp_lastIndexOf_1(l Ext_java_util_List, v any) int32 {
        \tt := ToAny(v)
        \tfor i := l.Size() - 1; i >= 0; i-- {
        \t\tif javalang.ObjectsEquals(t, ToAny(l.Get(i))) {
        \t\t\treturn i
        \t\t}
        \t}
        \treturn -1
        }
        func abstractListOp_clear_0(l Ext_java_util_List)              { abstractColClear(l) }
        func abstractListOp_iterator_0(l Ext_java_util_List) any       { return &jAbstractListIter{lst: l, lastRet: -1} }
        func abstractListOp_listIterator_0(l Ext_java_util_List) any   { return abstractListOp_listIterator_1(l, 0) }
        func abstractListOp_listIterator_1(l Ext_java_util_List, i int32) any {
        \tCheckPositionIndex(i, l.Size())
        \treturn &jAbstractListIter{lst: l, cursor: i, lastRet: -1}
        }
        func abstractListOp_subList_2(l Ext_java_util_List, from, to int32) any {
        \tCheckPositionIndexes(from, to, l.Size())
        \titems := make([]any, to-from)
        \tfor i := from; i < to; i++ {
        \t\titems[i-from] = l.Get(i)
        \t}
        \treturn wrapArrayList(items)
        }
        func abstractListOp_toArray_0(l Ext_java_util_List) []any                { return colToSlice(l) }
        func abstractListOp_toArray_1(l Ext_java_util_List, dst []any) []any     { return sliceToArray1(colToSlice(l), dst) }
        func abstractListOp_containsAll_1(l Ext_java_util_List, c Ext_java_util_Collection) bool {
        \treturn colContainsAll(l.Contains, c)
        }
        func abstractListOp_addAll_1(l Ext_java_util_List, c Ext_java_util_Collection) bool {
        \treturn abstractColAddAll(l, c)
        }
        func abstractListOp_removeAll_1(l Ext_java_util_List, c Ext_java_util_Collection) bool {
        \treturn abstractColRemoveAll(l, c)
        }
        func abstractListOp_retainAll_1(l Ext_java_util_List, c Ext_java_util_Collection) bool {
        \treturn abstractColRetainAll(l, c)
        }
        func abstractListOp_removeIf_1(l Ext_java_util_List, p Ext_java_util_function_Predicate) bool {
        \treturn abstractColRemoveIf(l, p)
        }
        func abstractListOp_equals_1(l Ext_java_util_List, o any) bool { return listEqualsSlice(colToSlice(l), o) }
        func abstractListOp_hashCode_0(l Ext_java_util_List) int32     { return listHashSlice(colToSlice(l)) }
        func abstractListOp_toString_0(l Ext_java_util_List) string    { return sliceJavaString(colToSlice(l)) }
        func abstractListOp_sort_1(l Ext_java_util_List, cmp Ext_java_util_Comparator) {
        \ts := colToSlice(l)
        \tfn := comparatorFuncOrNatural(cmp)
        \tsort.SliceStable(s, func(i, j int) bool { return fn(s[i], s[j]) < 0 })
        \tfor i, v := range s {
        \t\tl.Set(int32(i), v)
        \t}
        }
        func abstractListOp_addFirst_1(l Ext_java_util_List, v any) { l.Add_2(0, v) }
        func abstractListOp_addLast_1(l Ext_java_util_List, v any)  { l.Add_2(l.Size(), v) }
        func abstractListOp_getFirst_0(l Ext_java_util_List) any {
        \tif l.IsEmpty() { panic(javalang.NewNoSuchElementException("empty")) }
        \treturn l.Get(0)
        }
        func abstractListOp_getLast_0(l Ext_java_util_List) any {
        \tif l.IsEmpty() { panic(javalang.NewNoSuchElementException("empty")) }
        \treturn l.Get(l.Size() - 1)
        }
        func abstractListOp_removeFirst_0(l Ext_java_util_List) any {
        \tif l.IsEmpty() { panic(javalang.NewNoSuchElementException("empty")) }
        \treturn l.Remove_1_int32(0)
        }
        func abstractListOp_removeLast_0(l Ext_java_util_List) any {
        \tif l.IsEmpty() { panic(javalang.NewNoSuchElementException("empty")) }
        \treturn l.Remove_1_int32(l.Size() - 1)
        }
        func abstractListOp_stream_0(l Ext_java_util_List) Ext_java_util_stream_Stream {
        \treturn streamFromIterable(l)
        }
        func abstractListOp_forEach_1(l Ext_java_util_List, c Ext_java_util_function_Consumer) {
        \tExtDefault_Iterable_ForEach(l, c)
        }

        // AbstractQueue delegation
        func abstractQueueOp_size_0(q Ext_java_util_Queue) int32     { panic(javalang.NewUnsupportedOperationException("Queue.size")) }
        func abstractQueueOp_iterator_0(q Ext_java_util_Queue) any   { panic(javalang.NewUnsupportedOperationException("Queue.iterator")) }
        func abstractQueueOp_poll_0(q Ext_java_util_Queue) any       { panic(javalang.NewUnsupportedOperationException("Queue.poll")) }
        func abstractQueueOp_peek_0(q Ext_java_util_Queue) any       { panic(javalang.NewUnsupportedOperationException("Queue.peek")) }
        func abstractQueueOp_add_1(q Ext_java_util_Queue, v any) bool {
        \tpanic(javalang.NewUnsupportedOperationException("Queue.add"))
        }
        func abstractQueueOp_remove_0(q Ext_java_util_Queue) any {
        \tx := q.Poll_0()
        \tif x != nil { return x }
        \tpanic(javalang.NewNoSuchElementException("empty queue"))
        }
        func abstractQueueOp_element_0(q Ext_java_util_Queue) any {
        \tx := q.Peek()
        \tif x != nil { return x }
        \tpanic(javalang.NewNoSuchElementException("empty queue"))
        }
        func abstractQueueOp_remove_1(q Ext_java_util_Queue, v any) bool { return abstractColRemove(q, v) }
        func abstractQueueOp_retainAll_1(q Ext_java_util_Queue, c Ext_java_util_Collection) bool {
        \treturn abstractColRetainAll(q, c)
        }
        func abstractQueueOp_removeIf_1(q Ext_java_util_Queue, p Ext_java_util_function_Predicate) bool {
        \treturn abstractColRemoveIf(q, p)
        }
        func abstractQueueOp_clear_0(q Ext_java_util_Queue)              { for q.Poll_0() != nil {} }
        func abstractQueueOp_addAll_1(q Ext_java_util_Queue, c Ext_java_util_Collection) bool {
        \treturn abstractColAddAll(q, c)
        }
        func abstractQueueOp_isEmpty_0(q Ext_java_util_Queue) bool              { return q.Size() == 0 }
        func abstractQueueOp_contains_1(q Ext_java_util_Queue, v any) bool      { return abstractColContains(q, v) }
        func abstractQueueOp_toArray_0(q Ext_java_util_Queue) []any             { return colToSlice(q) }
        func abstractQueueOp_toArray_1(q Ext_java_util_Queue, d []any) []any    { return sliceToArray1(colToSlice(q), d) }
        func abstractQueueOp_containsAll_1(q Ext_java_util_Queue, c Ext_java_util_Collection) bool {
        \treturn colContainsAll(q.Contains, c)
        }
        func abstractQueueOp_removeAll_1(q Ext_java_util_Queue, c Ext_java_util_Collection) bool {
        \treturn abstractColRemoveAll(q, c)
        }
        func abstractQueueOp_toString_0(q Ext_java_util_Queue) string  { return sliceJavaString(colToSlice(q)) }
        func abstractQueueOp_equals_1(q Ext_java_util_Queue, o any) bool { return RefEquals(q, o) }
        func abstractQueueOp_hashCode_0(q Ext_java_util_Queue) int32   { return javalang.IdentityHashCode(q) }
        func abstractQueueOp_stream_0(q Ext_java_util_Queue) Ext_java_util_stream_Stream {
        \treturn streamFromIterable(q)
        }
        func abstractQueueOp_forEach_1(q Ext_java_util_Queue, c Ext_java_util_function_Consumer) {
        \tExtDefault_Iterable_ForEach(q, c)
        }

        // AbstractMap delegation (via EntrySet())
        func abstractMapOp_entrySet_0(m Ext_java_util_Map) any {
        \tpanic(javalang.NewUnsupportedOperationException("AbstractMap.entrySet"))
        }
        func abstractMapOp_size_0(m Ext_java_util_Map) int32 {
        \treturn Cast[Ext_java_util_Set](m.EntrySet()).Size()
        }
        func abstractMapOp_isEmpty_0(m Ext_java_util_Map) bool { return m.Size() == 0 }
        func abstractMapOp_containsKey_1(m Ext_java_util_Map, k any) bool {
        \tt := ToAny(k)
        \tfor it := Iterate(m.EntrySet()); it.HasNext(); {
        \t\te := Cast[Ext_java_util_Map_Entry](it.Next())
        \t\tif javalang.ObjectsEquals(t, ToAny(e.GetKey())) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }
        func abstractMapOp_containsValue_1(m Ext_java_util_Map, v any) bool {
        \tt := ToAny(v)
        \tfor it := Iterate(m.EntrySet()); it.HasNext(); {
        \t\te := Cast[Ext_java_util_Map_Entry](it.Next())
        \t\tif javalang.ObjectsEquals(t, ToAny(e.GetValue())) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }
        func abstractMapOp_get_1(m Ext_java_util_Map, k any) any {
        \tt := ToAny(k)
        \tfor it := Iterate(m.EntrySet()); it.HasNext(); {
        \t\te := Cast[Ext_java_util_Map_Entry](it.Next())
        \t\tif javalang.ObjectsEquals(t, ToAny(e.GetKey())) {
        \t\t\treturn e.GetValue()
        \t\t}
        \t}
        \treturn nil
        }
        func abstractMapOp_put_2(m Ext_java_util_Map, k, v any) any {
        \tpanic(javalang.NewUnsupportedOperationException("AbstractMap.put"))
        }
        func abstractMapOp_remove_1(m Ext_java_util_Map, k any) any {
        \tt := ToAny(k)
        \tit := Cast[Ext_java_util_Iterator](Cast[Ext_java_util_Set](m.EntrySet()).Iterator())
        \tfor it.HasNext() {
        \t\te := Cast[Ext_java_util_Map_Entry](it.Next())
        \t\tif javalang.ObjectsEquals(t, ToAny(e.GetKey())) {
        \t\t\told := e.GetValue()
        \t\t\tit.Remove()
        \t\t\treturn old
        \t\t}
        \t}
        \treturn nil
        }
        func abstractMapOp_putAll_1(m Ext_java_util_Map, other Ext_java_util_Map) {
        \tCheckNotNull(other)
        \tfor it := Iterate(other.EntrySet()); it.HasNext(); {
        \t\te := Cast[Ext_java_util_Map_Entry](it.Next())
        \t\tm.Put_2(e.GetKey(), e.GetValue())
        \t}
        }
        func abstractMapOp_clear_0(m Ext_java_util_Map) {
        \tCast[Ext_java_util_Set](m.EntrySet()).Clear()
        }
        func abstractMapOp_keySet_0(m Ext_java_util_Map) any {
        \tvar keys []any
        \tfor it := Iterate(m.EntrySet()); it.HasNext(); {
        \t\tkeys = append(keys, Cast[Ext_java_util_Map_Entry](it.Next()).GetKey())
        \t}
        \treturn wrapLinkedHashSet(keys)
        }
        func abstractMapOp_values_0(m Ext_java_util_Map) any {
        \tvar vals []any
        \tfor it := Iterate(m.EntrySet()); it.HasNext(); {
        \t\tvals = append(vals, Cast[Ext_java_util_Map_Entry](it.Next()).GetValue())
        \t}
        \treturn wrapArrayList(vals)
        }
        func abstractMapOp_equals_1(m Ext_java_util_Map, o any) bool {
        \tmp, ok := o.(Ext_java_util_Map)
        \tif !ok || mp == nil || javalang.IsNil(mp) { return false }
        \tif m.Size() != mp.Size() { return false }
        \tfor it := Iterate(m.EntrySet()); it.HasNext(); {
        \t\te := Cast[Ext_java_util_Map_Entry](it.Next())
        \t\tk, v := e.GetKey(), e.GetValue()
        \t\tif v == nil {
        \t\t\tif mp.Get(k) != nil || !mp.ContainsKey(k) { return false }
        \t\t} else if !javalang.ObjectsEquals(ToAny(v), ToAny(mp.Get(k))) {
        \t\t\treturn false
        \t\t}
        \t}
        \treturn true
        }
        func abstractMapOp_hashCode_0(m Ext_java_util_Map) int32 {
        \th := int32(0)
        \tfor it := Iterate(m.EntrySet()); it.HasNext(); {
        \t\th += mapEntryHash(Cast[Ext_java_util_Map_Entry](it.Next()))
        \t}
        \treturn h
        }
        func abstractMapOp_toString_0(m Ext_java_util_Map) string {
        \tvar sb strings.Builder
        \tsb.WriteByte('{')
        \tfirst := true
        \tfor it := Iterate(m.EntrySet()); it.HasNext(); {
        \t\tif !first {
        \t\t\tsb.WriteString(", ")
        \t\t}
        \t\tfirst = false
        \t\te := Cast[Ext_java_util_Map_Entry](it.Next())
        \t\tk, v := e.GetKey(), e.GetValue()
        \t\tif k == any(m) {
        \t\t\tsb.WriteString("(this Map)")
        \t\t} else {
        \t\t\tsb.WriteString(StrVal(k))
        \t\t}
        \t\tsb.WriteByte('=')
        \t\tif v == any(m) {
        \t\t\tsb.WriteString("(this Map)")
        \t\t} else {
        \t\t\tsb.WriteString(StrVal(v))
        \t\t}
        \t}
        \tsb.WriteByte('}')
        \treturn sb.String()
        }
        func abstractMapOp_getOrDefault_2(m Ext_java_util_Map, k, def any) any {
        \tif v := m.Get(k); v != nil || m.ContainsKey(k) {
        \t\treturn v
        \t}
        \treturn ToAny(def)
        }
        func abstractMapOp_forEach_1(m Ext_java_util_Map, fn Ext_java_util_function_BiConsumer) {
        \tCheckNotNull(fn)
        \tfor it := Iterate(m.EntrySet()); it.HasNext(); {
        \t\te := Cast[Ext_java_util_Map_Entry](it.Next())
        \t\tfn.Accept(e.GetKey(), e.GetValue())
        \t}
        }
        func abstractMapOp_replaceAll_1(m Ext_java_util_Map, fn Ext_java_util_function_BiFunction) {
        \tCheckNotNull(fn)
        \tfor it := Iterate(m.EntrySet()); it.HasNext(); {
        \t\te := Cast[Ext_java_util_Map_Entry](it.Next())
        \t\te.SetValue(fn.Apply(e.GetKey(), e.GetValue()))
        \t}
        }
        func abstractMapOp_putIfAbsent_2(m Ext_java_util_Map, k, v any) any {
        \tif cur := m.Get(k); cur != nil {
        \t\treturn cur
        \t}
        \treturn m.Put_2(k, v)
        }
        func abstractMapOp_replace_3(m Ext_java_util_Map, k, oldV, newV any) bool {
        \tif !m.ContainsKey(k) || !javalang.ObjectsEquals(ToAny(m.Get(k)), ToAny(oldV)) {
        \t\treturn false
        \t}
        \tm.Put_2(k, newV)
        \treturn true
        }
        func abstractMapOp_computeIfAbsent_2(m Ext_java_util_Map, k any, fn Ext_java_util_function_Function) any {
        \tCheckNotNull(fn)
        \tif v := m.Get(k); v != nil { return v }
        \tnv := ToAny(fn.Apply(ToAny(k)))
        \tif nv != nil { m.Put_2(k, nv) }
        \treturn nv
        }
        func abstractMapOp_compute_2(m Ext_java_util_Map, k any, fn Ext_java_util_function_BiFunction) any {
        \tCheckNotNull(fn)
        \toldV := m.Get(k)
        \tnv := ToAny(fn.Apply(ToAny(k), oldV))
        \tif nv == nil {
        \t\tif m.ContainsKey(k) { m.Remove_1(k) }
        \t\treturn nil
        \t}
        \tm.Put_2(k, nv)
        \treturn nv
        }
        func abstractMapOp_merge_3(m Ext_java_util_Map, k, v any, fn Ext_java_util_function_BiFunction) any {
        \tCheckNotNull(fn)
        \tv = CheckNotNull(ToAny(v))
        \toldV := m.Get(k)
        \tvar nv any
        \tif oldV == nil { nv = v } else { nv = ToAny(fn.Apply(oldV, v)) }
        \tif nv == nil { m.Remove_1(k) } else { m.Put_2(k, nv) }
        \treturn nv
        }

        // --- Sets.SetView ---

        type jSetViewCore struct {
        \ts1   Ext_java_util_Set
        \ts2   Ext_java_util_Set
        \tmode int // 0=union, 1=intersection, 2=difference
        }

        func newSetView(s1, s2 Ext_java_util_Set, mode int) Ext_com_google_common_collect_Sets_SetView {
        \tCheckNotNull(s1)
        \tCheckNotNull(s2)
        \tv := &Ext_com_google_common_collect_Sets_SetViewBase{}
        \tv.BindExtSelf_(v)
        \tv.Core_ = &jSetViewCore{s1: s1, s2: s2, mode: mode}
        \treturn v
        }

        func (c *jSetViewCore) contains(x any) bool {
        \tswitch c.mode {
        \tcase 0:
        \t\treturn c.s1.Contains(x) || c.s2.Contains(x)
        \tcase 1:
        \t\treturn c.s1.Contains(x) && c.s2.Contains(x)
        \tdefault:
        \t\treturn c.s1.Contains(x) && !c.s2.Contains(x)
        \t}
        }

        func (c *jSetViewCore) compute() []any {
        \tvar sc JSetCore
        \tswitch c.mode {
        \tcase 0:
        \t\tfor it := Iterate(c.s1); it.HasNext(); { sc.Add(it.Next()) }
        \t\tfor it := Iterate(c.s2); it.HasNext(); { sc.Add(it.Next()) }
        \tcase 1:
        \t\tfor it := Iterate(c.s1); it.HasNext(); {
        \t\t\tx := it.Next()
        \t\t\tif c.s2.Contains(x) { sc.Add(x) }
        \t\t}
        \tdefault:
        \t\tfor it := Iterate(c.s1); it.HasNext(); {
        \t\t\tx := it.Next()
        \t\t\tif !c.s2.Contains(x) { sc.Add(x) }
        \t\t}
        \t}
        \treturn sc.LiveKeys()
        }

        // --- ConcurrentHashMap ---

        type jSyncMapCore struct {
        \tmu sync.Mutex
        \tmc *JMapCore
        }

        func syncMapCore(m *Ext_java_util_concurrent_ConcurrentHashMap) *jSyncMapCore {
        \tCheckNotNullPtr(m)
        \tif m.Core_ == nil {
        \t\tm.Core_ = &jSyncMapCore{mc: NewJMapCore(false)}
        \t}
        \treturn m.Core_.(*jSyncMapCore)
        }

        func syncMapOp_size_0(m *Ext_java_util_concurrent_ConcurrentHashMap, _ ...any) int32 {
        \tc := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.Size()
        }
        func syncMapOp_isEmpty_0(m *Ext_java_util_concurrent_ConcurrentHashMap, _ ...any) bool {
        \tc := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.IsEmpty()
        }
        func syncMapOp_get_1(m *Ext_java_util_concurrent_ConcurrentHashMap, k any) any {
        \tCheckNotNull(k); c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.Get(k)
        }
        func syncMapOp_getOrDefault_2(m *Ext_java_util_concurrent_ConcurrentHashMap, k, def any) any {
        \tCheckNotNull(k); c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.GetOrDefault(k, def)
        }
        func syncMapOp_containsKey_1(m *Ext_java_util_concurrent_ConcurrentHashMap, k any) bool {
        \tCheckNotNull(k); c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.ContainsKey(k)
        }
        func syncMapOp_containsValue_1(m *Ext_java_util_concurrent_ConcurrentHashMap, v any) bool {
        \tCheckNotNull(v); c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.ContainsValue(v)
        }
        func syncMapOp_put_2(m *Ext_java_util_concurrent_ConcurrentHashMap, k, v any) any {
        \tCheckNotNull(k); CheckNotNull(v); c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.Put(k, v)
        }
        func syncMapOp_putIfAbsent_2(m *Ext_java_util_concurrent_ConcurrentHashMap, k, v any) any {
        \tCheckNotNull(k); CheckNotNull(v); c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.PutIfAbsent(k, v)
        }
        func syncMapOp_putAll_1(m *Ext_java_util_concurrent_ConcurrentHashMap, src Ext_java_util_Map) {
        \tCheckNotNull(src); c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); copyExtMapInto(c.mc, src)
        }
        func syncMapOp_remove_1(m *Ext_java_util_concurrent_ConcurrentHashMap, k any) any {
        \tCheckNotNull(k); c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.Remove(k)
        }
        func syncMapOp_clear_0(m *Ext_java_util_concurrent_ConcurrentHashMap, _ ...any) {
        \tc := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); c.mc.Clear()
        }
        func syncMapOp_keySet_0(m *Ext_java_util_concurrent_ConcurrentHashMap, _ ...any) any {
        \tc := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return newLiveMapKeySet(c.mc)
        }
        func syncMapOp_values_0(m *Ext_java_util_concurrent_ConcurrentHashMap, _ ...any) any {
        \tc := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return newLiveMapValues(c.mc)
        }
        func syncMapOp_entrySet_0(m *Ext_java_util_concurrent_ConcurrentHashMap, _ ...any) any {
        \tc := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return newLiveMapEntrySet(c.mc)
        }
        func syncMapOp_hashCode_0(m *Ext_java_util_concurrent_ConcurrentHashMap, _ ...any) int32 {
        \tc := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return mapCoreHash(c.mc)
        }
        func syncMapOp_toString_0(m *Ext_java_util_concurrent_ConcurrentHashMap, _ ...any) string {
        \tc := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return mapCoreString(c.mc)
        }
        func syncMapOp_equals_1(m *Ext_java_util_concurrent_ConcurrentHashMap, o any) bool {
        \tc := syncMapCore(m); c.mu.Lock(); snapshot := cloneJMapCore(c.mc); c.mu.Unlock(); return mapCoreEquals(snapshot, o)
        }
        func syncMapOp_replace_3(m *Ext_java_util_concurrent_ConcurrentHashMap, k, ov, nv any) bool {
        \tc := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return mapCoreReplace3(c.mc, k, ov, nv)
        }
        func syncMapOp_forEach_1(m *Ext_java_util_concurrent_ConcurrentHashMap, fn Ext_java_util_function_BiConsumer) {
        \tc := syncMapCore(m); c.mu.Lock(); snapshot := cloneJMapCore(c.mc); c.mu.Unlock(); mapCoreForEach(snapshot, fn)
        }
        func syncMapOp_replaceAll_1(m *Ext_java_util_concurrent_ConcurrentHashMap, fn Ext_java_util_function_BiFunction) {
        \tc := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); mapCoreReplaceAll(c.mc, fn)
        }
        func syncMapOp_computeIfAbsent_2(m *Ext_java_util_concurrent_ConcurrentHashMap, k any, fn Ext_java_util_function_Function) any {
        \tCheckNotNull(k); CheckNotNull(fn); c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return mapCoreComputeIfAbsent(c.mc, k, fn)
        }
        func syncMapComputeIfAbsentStr(m *Ext_java_util_concurrent_ConcurrentHashMap, k string, fn Ext_java_util_function_Function) any {
        \tCheckNotNullString(k); CheckNotNull(fn); c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock()
        \tif v := c.mc.GetStr(k); v != nil { return v }
        \tnv := ToAny(fn.Apply(BoxString(k)))
        \tif nv != nil { c.mc.PutStr(k, nv) }
        \treturn nv
        }
        func syncMapComputeIfAbsentNode(m *Ext_java_util_concurrent_ConcurrentHashMap, n uint32, fn Ext_java_util_function_Function) any {
        \tif n == 0 { ThrowNPE() }; CheckNotNull(fn); c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock()
        \tif v := c.mc.GetNode(n); v != nil { return v }
        \tnv := ToAny(fn.Apply(boxNodeById(n)))
        \tif nv != nil { c.mc.PutNode(n, nv) }
        \treturn nv
        }
        func syncMapTryGetComputeStr(m *Ext_java_util_concurrent_ConcurrentHashMap, k string) (any, bool) {
        \tif m == nil { ThrowNPE() }
        \tif k == NullString { return nil, false }
        \tc := syncMapCore(m); c.mu.Lock(); v := c.mc.GetStr(k); c.mu.Unlock()
        \tif v != nil { return v, true }
        \treturn nil, false
        }
        func syncMapTryGetComputeNode(m *Ext_java_util_concurrent_ConcurrentHashMap, n uint32) (any, bool) {
        \tif m == nil { ThrowNPE() }
        \tif n == 0 { return nil, false }
        \tc := syncMapCore(m); c.mu.Lock(); v := c.mc.GetNode(n); c.mu.Unlock()
        \tif v != nil { return v, true }
        \treturn nil, false
        }
        func syncMapTryGetComputeAny(m *Ext_java_util_concurrent_ConcurrentHashMap, k any) (any, bool) {
        \tif m == nil { ThrowNPE() }
        \tif k == nil || isJavaNilAny(k) { return nil, false }
        \tc := syncMapCore(m); c.mu.Lock(); v := c.mc.Get(k); c.mu.Unlock()
        \tif v != nil { return v, true }
        \treturn nil, false
        }
        func syncMapOp_compute_2(m *Ext_java_util_concurrent_ConcurrentHashMap, k any, fn Ext_java_util_function_BiFunction) any {
        \tc := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return mapCoreCompute(c.mc, k, fn)
        }
        func syncMapOp_merge_3(m *Ext_java_util_concurrent_ConcurrentHashMap, k, v any, fn Ext_java_util_function_BiFunction) any {
        \tc := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return mapCoreMerge(c.mc, k, v, fn)
        }
        func syncMapGetStr(m *Ext_java_util_concurrent_ConcurrentHashMap, k string) any {
        \tCheckNotNullString(k); c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.GetStr(k)
        }
        func syncMapGetOrDefaultStr(m *Ext_java_util_concurrent_ConcurrentHashMap, k string, def any) any {
        \tCheckNotNullString(k); c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.GetOrDefaultStr(k, def)
        }
        func syncMapContainsStr(m *Ext_java_util_concurrent_ConcurrentHashMap, k string) bool {
        \tCheckNotNullString(k); c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.ContainsStr(k)
        }
        func syncMapPutStr(m *Ext_java_util_concurrent_ConcurrentHashMap, k string, v any) any {
        \tCheckNotNullString(k); CheckNotNull(v); c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.PutStr(k, v)
        }
        func syncMapRemoveStr(m *Ext_java_util_concurrent_ConcurrentHashMap, k string) any {
        \tCheckNotNullString(k); c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.RemoveStr(k)
        }
        func syncMapGetNode(m *Ext_java_util_concurrent_ConcurrentHashMap, n uint32) any {
        \tif n == 0 { ThrowNPE() }; c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.GetNode(n)
        }
        func syncMapGetOrDefaultNode(m *Ext_java_util_concurrent_ConcurrentHashMap, n uint32, def any) any {
        \tif n == 0 { ThrowNPE() }; c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.GetOrDefaultNode(n, def)
        }
        func syncMapContainsNode(m *Ext_java_util_concurrent_ConcurrentHashMap, n uint32) bool {
        \tif n == 0 { ThrowNPE() }; c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.ContainsNode(n)
        }
        func syncMapPutNode(m *Ext_java_util_concurrent_ConcurrentHashMap, n uint32, v any) any {
        \tif n == 0 { ThrowNPE() }; CheckNotNull(v); c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.PutNode(n, v)
        }
        func syncMapRemoveNode(m *Ext_java_util_concurrent_ConcurrentHashMap, n uint32) any {
        \tif n == 0 { ThrowNPE() }; c := syncMapCore(m); c.mu.Lock(); defer c.mu.Unlock(); return c.mc.RemoveNode(n)
        }
        """;
  }

  private static String buildColRtPart2() {
    return buildColRtPart2a() + buildColRtPart2b();
  }

  private static String buildColRtPart2a() {
    return """

        // --- BiMap (HashBiMap & ImmutableBiMap) ---

        type jBiMapCore struct {
        \tfwd           *JMapCore
        \trev           *JMapCore
        \tinverse       any
        \timmutable     bool
        \tkeySetCache   atomic.Pointer[Ext_com_google_common_collect_ImmutableSetBase]
        \tvaluesCache   atomic.Pointer[Ext_com_google_common_collect_ImmutableSetBase]
        \tentrySetCache atomic.Pointer[Ext_com_google_common_collect_ImmutableSetBase]
        }

        func newHashBiMap() *Ext_com_google_common_collect_HashBiMap {
        \tfwd := NewJMapCore(false)
        \trev := NewJMapCore(false)
        \tm1 := &Ext_com_google_common_collect_HashBiMap{}
        \tm2 := &Ext_com_google_common_collect_HashBiMap{}
        \tc1 := &jBiMapCore{fwd: fwd, rev: rev, inverse: m2}
        \tc2 := &jBiMapCore{fwd: rev, rev: fwd, inverse: m1}
        \tm1.Core_ = c1
        \tm2.Core_ = c2
        \treturn m1
        }

        func newImmutableBiMapFrom(src Ext_java_util_Map) Ext_com_google_common_collect_ImmutableBiMap {
        \tfwd := NewJMapCore(false)
        \trev := NewJMapCore(false)
        \tif src != nil && !javalang.IsNil(src) {
        \t\tfor it := Iterate(src.EntrySet()); it.HasNext(); {
        \t\t\te := Cast[Ext_java_util_Map_Entry](it.Next())
        \t\t\tk := CheckNotNull(ToAny(e.GetKey()))
        \t\t\tv := CheckNotNull(ToAny(e.GetValue()))
        \t\t\tif fwd.ContainsKey(k) || rev.ContainsKey(v) {
        \t\t\t\tpanic(javalang.NewIllegalArgumentException("duplicate key or value in ImmutableBiMap"))
        \t\t\t}
        \t\t\tfwd.Put(k, v)
        \t\t\trev.Put(v, k)
        \t\t}
        \t}
        \tm1 := &Ext_com_google_common_collect_ImmutableBiMapBase{}
        \tm1.BindExtSelf_(m1)
        \tm2 := &Ext_com_google_common_collect_ImmutableBiMapBase{}
        \tm2.BindExtSelf_(m2)
        \tm1.Core_ = &jBiMapCore{fwd: fwd, rev: rev, inverse: m2, immutable: true}
        \tm2.Core_ = &jBiMapCore{fwd: rev, rev: fwd, inverse: m1, immutable: true}
        \treturn m1
        }

        func biMapOp_size_0(b *jBiMapCore) int32            { return b.fwd.Size() }
        func biMapOp_toString_0(b *jBiMapCore) string       { return mapCoreString(b.fwd) }
        func biMapOp_isEmpty_0(b *jBiMapCore) bool          { return b.fwd.IsEmpty() }
        func biMapOp_containsKey_1(b *jBiMapCore, k any) bool   { return b.fwd.ContainsKey(k) }
        func biMapOp_containsValue_1(b *jBiMapCore, v any) bool { return b.rev.ContainsKey(v) }
        func biMapOp_get_1(b *jBiMapCore, k any) any        { return b.fwd.Get(k) }
        func biMapOp_getOrDefault_2(b *jBiMapCore, k, d any) any { return b.fwd.GetOrDefault(k, d) }
        func biMapOp_inverse_0(b *jBiMapCore) any           { return b.inverse }
        func biMapOp_put_2(b *jBiMapCore, k, v any) any {
        \tif b.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tk, v = ToAny(k), ToAny(v)
        \tif b.fwd.ContainsKey(k) && javalang.ObjectsEquals(b.fwd.Get(k), v) {
        \t\treturn v
        \t}
        \tif b.rev.ContainsKey(v) {
        \t\tpanic(javalang.NewIllegalArgumentException("value already present: " + StrVal(v)))
        \t}
        \tif b.fwd.ContainsKey(k) {
        \t\toldV := b.fwd.Get(k)
        \t\tb.rev.Remove(oldV)
        \t}
        \told := b.fwd.Put(k, v)
        \tb.rev.Put(v, k)
        \treturn old
        }
        func biMapOp_putIfAbsent_2(b *jBiMapCore, k, v any) any {
        \tif b.fwd.ContainsKey(k) && b.fwd.Get(k) != nil { return b.fwd.Get(k) }
        \treturn biMapOp_put_2(b, k, v)
        }
        func biMapPutStr(b *jBiMapCore, k string, v any) any {
        \tif b.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tv = ToAny(v)
        \tif b.fwd.ContainsStr(k) && javalang.ObjectsEquals(b.fwd.GetStr(k), v) {
        \t\treturn v
        \t}
        \tif b.rev.ContainsKey(v) {
        \t\tpanic(javalang.NewIllegalArgumentException("value already present: " + StrVal(v)))
        \t}
        \tif b.fwd.ContainsStr(k) {
        \t\toldV := b.fwd.GetStr(k)
        \t\tif os, ok := oldV.(string); ok {
        \t\t\tb.rev.RemoveStr(os)
        \t\t} else {
        \t\t\tb.rev.Remove(oldV)
        \t\t}
        \t}
        \told := b.fwd.PutStr(k, v)
        \tbk := BoxString(k)
        \tif vs, ok := v.(string); ok {
        \t\tb.rev.PutStr(vs, bk)
        \t} else {
        \t\tb.rev.Put(v, bk)
        \t}
        \treturn old
        }
        func biMapRemoveStr(b *jBiMapCore, k string) any {
        \tif b.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tif !b.fwd.ContainsStr(k) { return nil }
        \toldV := b.fwd.RemoveStr(k)
        \tif os, ok := oldV.(string); ok {
        \t\tb.rev.RemoveStr(os)
        \t} else {
        \t\tb.rev.Remove(oldV)
        \t}
        \treturn oldV
        }
        func biMapOp_remove_1(b *jBiMapCore, k any) any {
        \tif b.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tif !b.fwd.ContainsKey(k) { return nil }
        \toldV := b.fwd.Remove(k)
        \tb.rev.Remove(oldV)
        \treturn oldV
        }
        func biMapOp_clear_0(b *jBiMapCore) {
        \tif b.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tb.fwd.Clear()
        \tb.rev.Clear()
        }
        func biMapOp_putAll_1(b *jBiMapCore, src Ext_java_util_Map) {
        \tCheckNotNull(src)
        \tfor it := Iterate(src.EntrySet()); it.HasNext(); {
        \t\te := Cast[Ext_java_util_Map_Entry](it.Next())
        \t\tbiMapOp_put_2(b, e.GetKey(), e.GetValue())
        \t}
        }
        func biMapOp_keySet_0(b *jBiMapCore) any {
        \tif b.immutable {
        \t\tif b.fwd.IsEmpty() { return emptyImmutableSet_ }
        \t\tif s := b.keySetCache.Load(); s != nil { return s }
        \t\ts := wrapImmutableSetCore(b.fwd).(*Ext_com_google_common_collect_ImmutableSetBase)
        \t\tif b.keySetCache.CompareAndSwap(nil, s) { return s }
        \t\treturn b.keySetCache.Load()
        \t}
        \treturn newLiveMapKeySet(b.fwd)
        }
        func biMapOp_values_0(b *jBiMapCore) any {
        \tif b.immutable {
        \t\tif b.fwd.IsEmpty() { return emptyImmutableSet_ }
        \t\tif s := b.valuesCache.Load(); s != nil { return s }
        \t\ts := wrapImmutableSet(b.fwd.LiveValues()).(*Ext_com_google_common_collect_ImmutableSetBase)
        \t\tif b.valuesCache.CompareAndSwap(nil, s) { return s }
        \t\treturn b.valuesCache.Load()
        \t}
        \treturn newLiveMapKeySet(b.rev)
        }
        func biMapOp_entrySet_0(b *jBiMapCore) any {
        \tif b.immutable {
        \t\tif b.fwd.IsEmpty() { return emptyImmutableSet_ }
        \t\tif s := b.entrySetCache.Load(); s != nil { return s }
        \t\ts := wrapImmutableMapEntrySet(b.fwd).(*Ext_com_google_common_collect_ImmutableSetBase)
        \t\tif b.entrySetCache.CompareAndSwap(nil, s) { return s }
        \t\treturn b.entrySetCache.Load()
        \t}
        \treturn newLiveMapEntrySet(b.fwd)
        }
        func biMapOp_equals_1(b *jBiMapCore, o any) bool { return mapCoreEquals(b.fwd, o) }
        func biMapOp_hashCode_0(b *jBiMapCore) int32     { return mapCoreHash(b.fwd) }
        func biMapOp_forEach_1(b *jBiMapCore, fn Ext_java_util_function_BiConsumer) { mapCoreForEach(b.fwd, fn) }
        func biMapOp_replaceAll_1(b *jBiMapCore, fn Ext_java_util_function_BiFunction) {
        \tif b.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tfor _, k := range b.fwd.LiveKeys() {
        \t\toldV := b.fwd.Get(k)
        \t\tbiMapOp_put_2(b, k, fn.Apply(k, oldV))
        \t}
        }
        func biMapOp_replace_3(b *jBiMapCore, k, ov, nv any) bool {
        \tif !b.fwd.ContainsKey(k) || !javalang.ObjectsEquals(b.fwd.Get(k), ToAny(ov)) { return false }
        \tbiMapOp_put_2(b, k, nv)
        \treturn true
        }
        func biMapOp_computeIfAbsent_2(b *jBiMapCore, k any, fn Ext_java_util_function_Function) any {
        \tif b.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tCheckNotNull(fn)
        \tif v := b.fwd.Get(k); v != nil { return v }
        \tnv := ToAny(fn.Apply(ToAny(k)))
        \tif nv != nil { biMapOp_put_2(b, k, nv) }
        \treturn nv
        }
        func biMapComputeIfAbsentStr(b *jBiMapCore, k string, fn Ext_java_util_function_Function) any {
        \tif b.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tCheckNotNull(fn)
        \tif v := b.fwd.GetStr(k); v != nil { return v }
        \tbk := BoxString(k)
        \tnv := ToAny(fn.Apply(bk))
        \tif nv != nil { biMapOp_put_2(b, bk, nv) }
        \treturn nv
        }
        func biMapComputeIfAbsentNode(b *jBiMapCore, n uint32, fn Ext_java_util_function_Function) any {
        \tif b.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tCheckNotNull(fn)
        \tif v := b.fwd.GetNode(n); v != nil { return v }
        \tbn := boxNodeById(n)
        \tnv := ToAny(fn.Apply(bn))
        \tif nv != nil { biMapOp_put_2(b, bn, nv) }
        \treturn nv
        }
        func biMapOp_compute_2(b *jBiMapCore, k any, fn Ext_java_util_function_BiFunction) any {
        \tCheckNotNull(fn)
        \toldV := b.fwd.Get(k)
        \tnv := ToAny(fn.Apply(ToAny(k), oldV))
        \tif nv == nil { biMapOp_remove_1(b, k); return nil }
        \tbiMapOp_put_2(b, k, nv)
        \treturn nv
        }
        func biMapOp_merge_3(b *jBiMapCore, k, v any, fn Ext_java_util_function_BiFunction) any {
        \tCheckNotNull(fn); v = CheckNotNull(ToAny(v))
        \toldV := b.fwd.Get(k)
        \tvar nv any
        \tif oldV == nil { nv = v } else { nv = ToAny(fn.Apply(oldV, v)) }
        \tif nv == nil { biMapOp_remove_1(b, k) } else { biMapOp_put_2(b, k, nv) }
        \treturn nv
        }

        // --- Multimap (ListMultimap, SetMultimap, SortedSetMultimap, Immutable*Multimap) ---

        type jMultimapCore struct {
        \tisSet     bool
        \timmutable bool
        \tnatVal    bool
        \tkeyCmp    func(a, b any) int32
        \tvalCmp    func(a, b any) int32
        \tkeys      *JMapCore // key -> *JListCore (if !isSet) or *JMapCore (if isSet)
        \ttotal     int32
        }

        func newMultimapCore(isSet, immutable bool, keyCmp, valCmp func(a, b any) int32) *jMultimapCore {
        \tvar km *JMapCore
        \tif keyCmp != nil {
        \t\tkm = NewSortedJMapCore(keyCmp)
        \t} else {
        \t\tkm = NewJMapCore(false)
        \t}
        \treturn &jMultimapCore{isSet: isSet, immutable: immutable, keyCmp: keyCmp, valCmp: valCmp, keys: km}
        }

        func newTreeMultimapCore() *jMultimapCore {
        \treturn &jMultimapCore{
        \t\tisSet:  true,
        \t\tnatVal: true,
        \t\tkeyCmp: CompareValues,
        \t\tvalCmp: CompareValues,
        \t\tkeys:   NewSortedJMapCore(nil),
        \t}
        }

        func (mm *jMultimapCore) newBucket() any {
        \tif mm.isSet {
        \t\tif mm.natVal {
        \t\t\treturn NewSortedJMapCore(nil)
        \t\t}
        \t\tif mm.valCmp != nil {
        \t\t\treturn NewSortedJMapCore(mm.valCmp)
        \t\t}
        \t\treturn NewJMapCore(false)
        \t}
        \treturn NewJListCore()
        }

        type jMultimapNilElem struct{}

        var jMultimapNilElem_ any = jMultimapNilElem{}

        func mmWrapSingle(v any) any {
        \tif v == nil {
        \t\treturn jMultimapNilElem_
        \t}
        \treturn v
        }

        func mmUnwrapSingle(b any) any {
        \tif _, isNil := b.(jMultimapNilElem); isNil {
        \t\treturn nil
        \t}
        \treturn b
        }

        func (mm *jMultimapCore) sealImmutableBuckets() {
        \tmm.keys.TransformValues(func(_, b any) any {
        \t\tif mm.isSet {
        \t\t\treturn wrapImmutableSetCore(b.(*JMapCore))
        \t\t}
        \t\tif lc, ok := b.(*JListCore); ok {
        \t\t\treturn wrapImmutableListOwn(lc.Items)
        \t\t}
        \t\treturn newImmList1(mmUnwrapSingle(b))
        \t})
        \tmm.keys.sealImmutable()
        \tmm.immutable = true
        }

        func multimapOp_size_0(mm *jMultimapCore) int32           { return mm.total }
        func multimapOp_isEmpty_0(mm *jMultimapCore) bool         { return mm.total == 0 }
        func multimapOp_containsKey_1(mm *jMultimapCore, k any) bool { return mm.keys.ContainsKey(k) }
        func (mm *jMultimapCore) containsKeyStr(k string) bool     { return mm.keys.ContainsStr(k) }
        func (mm *jMultimapCore) containsKeyNode(u uint32) bool    { return mm.keys.ContainsNode(u) }
        func multimapBucketContains(mm *jMultimapCore, b any, v any) bool {
        \tif b == nil { return false }
        \tif mm.immutable {
        \t\tif mm.isSet { return b.(Ext_com_google_common_collect_ImmutableSet).Contains(v) }
        \t\treturn b.(Ext_com_google_common_collect_ImmutableList).Contains(v)
        \t}
        \tif mm.isSet { return b.(*JMapCore).ContainsKey(v) }
        \tif lc, ok := b.(*JListCore); ok { return lc.Contains(v) }
        \treturn jSetElemEquals(mmUnwrapSingle(b), v)
        }
        func multimapOp_containsEntry_2(mm *jMultimapCore, k, v any) bool {
        \treturn multimapBucketContains(mm, mm.keys.Get(k), v)
        }
        func (mm *jMultimapCore) containsEntryStr(k string, v any) bool {
        \treturn multimapBucketContains(mm, mm.keys.GetStr(k), v)
        }
        func (mm *jMultimapCore) containsEntryNode(u uint32, v any) bool {
        \treturn multimapBucketContains(mm, mm.keys.GetNode(u), v)
        }
        func multimapOp_put_2(mm *jMultimapCore, k, v any) bool {
        \tif mm.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \treturn mm.putInternal(k, v)
        }
        func (mm *jMultimapCore) putSetBucketVal(b any, v any) bool {
        \tsm := b.(*JMapCore)
        \tif sm.ContainsKey(v) { return false }
        \tsm.Put(v, true)
        \tmm.total++
        \treturn true
        }
        func (mm *jMultimapCore) putInternal(k, v any) bool {
        \tif s, ok := k.(string); ok {
        \t\treturn mm.putStrInternal(s, v)
        \t}
        \tk, v = ToAny(k), ToAny(v)
        \tb := mm.keys.Get(k)
        \tif mm.isSet {
        \t\tif b == nil {
        \t\t\tb = mm.newBucket()
        \t\t\tmm.keys.Put(k, b)
        \t\t}
        \t\treturn mm.putSetBucketVal(b, v)
        \t}
        \tif b == nil {
        \t\tmm.keys.Put(k, mmWrapSingle(v))
        \t\tmm.total++
        \t\treturn true
        \t}
        \tif lc, ok := b.(*JListCore); ok {
        \t\tlc.Add(v)
        \t\tmm.total++
        \t\treturn true
        \t}
        \tlc := &JListCore{Items: make([]any, 2, 4)}
        \tlc.Items[0] = mmUnwrapSingle(b)
        \tlc.Items[1] = v
        \tmm.keys.Put(k, lc)
        \tmm.total++
        \treturn true
        }
        func (mm *jMultimapCore) putStrInternal(k string, v any) bool {
        \tv = ToAny(v)
        \tb := mm.keys.GetStr(k)
        \tif mm.isSet {
        \t\tif b == nil {
        \t\t\tb = mm.newBucket()
        \t\t\tmm.keys.PutStr(k, b)
        \t\t}
        \t\treturn mm.putSetBucketVal(b, v)
        \t}
        \tif b == nil {
        \t\tmm.keys.PutStr(k, mmWrapSingle(v))
        \t\tmm.total++
        \t\treturn true
        \t}
        \tif lc, ok := b.(*JListCore); ok {
        \t\tlc.Add(v)
        \t\tmm.total++
        \t\treturn true
        \t}
        \tlc := &JListCore{Items: make([]any, 2, 4)}
        \tlc.Items[0] = mmUnwrapSingle(b)
        \tlc.Items[1] = v
        \tmm.keys.PutStr(k, lc)
        \tmm.total++
        \treturn true
        }
        func (mm *jMultimapCore) putNodeInternal(u uint32, v any) bool {
        \tv = ToAny(v)
        \tb := mm.keys.GetNode(u)
        \tif mm.isSet {
        \t\tif b == nil {
        \t\t\tb = mm.newBucket()
        \t\t\tmm.keys.PutNode(u, b)
        \t\t}
        \t\treturn mm.putSetBucketVal(b, v)
        \t}
        \tif b == nil {
        \t\tmm.keys.PutNode(u, mmWrapSingle(v))
        \t\tmm.total++
        \t\treturn true
        \t}
        \tif lc, ok := b.(*JListCore); ok {
        \t\tlc.Add(v)
        \t\tmm.total++
        \t\treturn true
        \t}
        \tlc := &JListCore{Items: make([]any, 2, 4)}
        \tlc.Items[0] = mmUnwrapSingle(b)
        \tlc.Items[1] = v
        \tmm.keys.PutNode(u, lc)
        \tmm.total++
        \treturn true
        }
        func (mm *jMultimapCore) putStr(k string, v any) bool {
        \tif mm.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \treturn mm.putStrInternal(k, v)
        }
        func (mm *jMultimapCore) putNode(u uint32, v any) bool {
        \tif mm.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \treturn mm.putNodeInternal(u, v)
        }
        func multimapOp_remove_2(mm *jMultimapCore, k, v any) bool {
        \tif mm.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tb := mm.keys.Get(k)
        \tif b == nil { return false }
        \tif mm.isSet {
        \t\tsm := b.(*JMapCore)
        \t\tif !sm.ContainsKey(v) { return false }
        \t\tsm.Remove(v)
        \t\tmm.total--
        \t\tif sm.IsEmpty() { mm.keys.Remove(k) }
        \t\treturn true
        \t}
        \tif lc, ok := b.(*JListCore); ok {
        \t\tif !lc.RemoveVal(v) { return false }
        \t\tmm.total--
        \t\tif lc.IsEmpty() { mm.keys.Remove(k) }
        \t\treturn true
        \t}
        \tif !jSetElemEquals(mmUnwrapSingle(b), v) { return false }
        \tmm.keys.Remove(k)
        \tmm.total--
        \treturn true
        }
        func (mm *jMultimapCore) removeStr(k string, v any) bool {
        \tif mm.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tb := mm.keys.GetStr(k)
        \tif b == nil { return false }
        \tif mm.isSet {
        \t\tsm := b.(*JMapCore)
        \t\tif !sm.ContainsKey(v) { return false }
        \t\tsm.Remove(v)
        \t\tmm.total--
        \t\tif sm.IsEmpty() { mm.keys.RemoveStr(k) }
        \t\treturn true
        \t}
        \tif lc, ok := b.(*JListCore); ok {
        \t\tif !lc.RemoveVal(v) { return false }
        \t\tmm.total--
        \t\tif lc.IsEmpty() { mm.keys.RemoveStr(k) }
        \t\treturn true
        \t}
        \tif !jSetElemEquals(mmUnwrapSingle(b), v) { return false }
        \tmm.keys.RemoveStr(k)
        \tmm.total--
        \treturn true
        }
        func (mm *jMultimapCore) removeNode(u uint32, v any) bool {
        \tif mm.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tb := mm.keys.GetNode(u)
        \tif b == nil { return false }
        \tif mm.isSet {
        \t\tsm := b.(*JMapCore)
        \t\tif !sm.ContainsKey(v) { return false }
        \t\tsm.Remove(v)
        \t\tmm.total--
        \t\tif sm.IsEmpty() { mm.keys.RemoveNode(u) }
        \t\treturn true
        \t}
        \tif lc, ok := b.(*JListCore); ok {
        \t\tif !lc.RemoveVal(v) { return false }
        \t\tmm.total--
        \t\tif lc.IsEmpty() { mm.keys.RemoveNode(u) }
        \t\treturn true
        \t}
        \tif !jSetElemEquals(mmUnwrapSingle(b), v) { return false }
        \tmm.keys.RemoveNode(u)
        \tmm.total--
        \treturn true
        }
        func multimapOp_putAll_2(mm *jMultimapCore, k any, vals Ext_java_lang_Iterable) bool {
        \tif mm.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tCheckNotNull(vals)
        \tchanged := false
        \tfor it := NewIterCursor(vals); it.HasNext(); {
        \t\tif mm.putInternal(k, it.Next()) { changed = true }
        \t}
        \treturn changed
        }
        func (mm *jMultimapCore) putAllStr(k string, vals Ext_java_lang_Iterable) bool {
        \tif mm.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tCheckNotNull(vals)
        \tchanged := false
        \tfor it := NewIterCursor(vals); it.HasNext(); {
        \t\tif mm.putStrInternal(k, it.Next()) { changed = true }
        \t}
        \treturn changed
        }
        func (mm *jMultimapCore) putAllNode(u uint32, vals Ext_java_lang_Iterable) bool {
        \tif mm.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tCheckNotNull(vals)
        \tchanged := false
        \tfor it := NewIterCursor(vals); it.HasNext(); {
        \t\tif mm.putNodeInternal(u, it.Next()) { changed = true }
        \t}
        \treturn changed
        }
        func multimapOp_putAll_1(mm *jMultimapCore, other Ext_com_google_common_collect_Multimap) bool {
        \tif mm.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tCheckNotNull(other)
        \tchanged := false
        \tasMap := Cast[Ext_java_util_Map](other.AsMap())
        \tfor it := Iterate(asMap.EntrySet()); it.HasNext(); {
        \t\te := Cast[Ext_java_util_Map_Entry](it.Next())
        \t\tcol := Cast[Ext_java_lang_Iterable](e.GetValue())
        \t\tif multimapOp_putAll_2(mm, e.GetKey(), col) { changed = true }
        \t}
        \treturn changed
        }
        func (mm *jMultimapCore) finishRemoveAllBucket(b any) any {
        \tif b == nil {
        \t\tif mm.isSet {
        \t\t\tif mm.natVal {
        \t\t\t\tts := &Ext_java_util_TreeSet{}; ts.BindExtSelf_(ts); ts.Core_ = NewSortedJMapCore(nil); return ts
        \t\t\t}
        \t\t\tif mm.valCmp != nil {
        \t\t\t\tts := &Ext_java_util_TreeSet{}; ts.BindExtSelf_(ts); ts.Core_ = NewSortedJMapCore(mm.valCmp); return ts
        \t\t\t}
        \t\t\treturn wrapLinkedHashSet(nil)
        \t\t}
        \t\treturn wrapArrayList(nil)
        \t}
        \tif mm.isSet {
        \t\tsm := b.(*JMapCore)
        \t\tmm.total -= sm.Size()
        \t\tif mm.valCmp != nil {
        \t\t\tts := &Ext_java_util_TreeSet{}; ts.BindExtSelf_(ts); ts.Core_ = sm; return ts
        \t\t}
        \t\treturn wrapLinkedHashSet(sm.LiveKeys())
        \t}
        \tif lc, ok := b.(*JListCore); ok {
        \t\tmm.total -= lc.Size()
        \t\treturn wrapArrayList(lc.Items)
        \t}
        \tmm.total--
        \treturn wrapArrayListOwn([]any{mmUnwrapSingle(b)})
        }
        func multimapOp_removeAll_1(mm *jMultimapCore, k any) any {
        \tif mm.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \treturn mm.finishRemoveAllBucket(mm.keys.Remove(k))
        }
        func (mm *jMultimapCore) removeAllStr(k string) any {
        \tif mm.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \treturn mm.finishRemoveAllBucket(mm.keys.RemoveStr(k))
        }
        func (mm *jMultimapCore) removeAllNode(u uint32) any {
        \tif mm.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \treturn mm.finishRemoveAllBucket(mm.keys.RemoveNode(u))
        }
        func multimapOp_clear_0(mm *jMultimapCore) {
        \tif mm.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tmm.keys.Clear()
        \tmm.total = 0
        }
        func (mm *jMultimapCore) finishImmGet(b any) any {
        \tif b != nil {
        \t\treturn b
        \t}
        \tif mm.isSet {
        \t\treturn emptyImmutableSet_
        \t}
        \treturn emptyImmutableList_
        }
        func multimapOp_get_1(mm *jMultimapCore, k any) any {
        \tif s, ok := k.(string); ok {
        \t\treturn mm.getStr(s)
        \t}
        \tk = ToAny(k)
        \tif mm.immutable {
        \t\treturn mm.finishImmGet(mm.keys.Get(k))
        \t}
        \tif mm.isSet {
        \t\tv := &jMultimapLiveSet{mm: mm, key: k}
        \t\tv.BindExtSelf_(v)
        \t\treturn v
        \t}
        \tv := &jMultimapLiveList{mm: mm, key: k}
        \tv.BindExtSelf_(v)
        \treturn v
        }
        func (mm *jMultimapCore) getStr(k string) any {
        \tif mm.immutable {
        \t\treturn mm.finishImmGet(mm.keys.GetStr(k))
        \t}
        \tif mm.isSet {
        \t\tv := &jMultimapLiveSet{mm: mm, strKey: k, keyKind: 1}
        \t\tv.BindExtSelf_(v)
        \t\treturn v
        \t}
        \tv := &jMultimapLiveList{mm: mm, strKey: k, keyKind: 1}
        \tv.BindExtSelf_(v)
        \treturn v
        }
        func (mm *jMultimapCore) getNode(u uint32) any {
        \tif mm.immutable {
        \t\treturn mm.finishImmGet(mm.keys.GetNode(u))
        \t}
        \tif mm.isSet {
        \t\tv := &jMultimapLiveSet{mm: mm, nodeKey: u, keyKind: 2}
        \t\tv.BindExtSelf_(v)
        \t\treturn v
        \t}
        \tv := &jMultimapLiveList{mm: mm, nodeKey: u, keyKind: 2}
        \tv.BindExtSelf_(v)
        \treturn v
        }
        func multimapOp_keySet_0(mm *jMultimapCore) any {
        \tif mm.immutable { return wrapImmutableSetCore(mm.keys) }
        \tif mm.keyCmp != nil {
        \t\tts := &Ext_java_util_TreeSet{}
        \t\tts.Core_ = cloneJMapCore(mm.keys)
        \t\treturn ts
        \t}
        \tv := &jMultimapLiveKeySet{mm: mm}
        \tv.BindExtSelf_(v)
        \treturn v
        }
        func multimapOp_values_0(mm *jMultimapCore) any {
        \tout := make([]any, 0, mm.total)
        \tif mm.immutable {
        \t\tmm.keys.ForEachEntry(func(_, b any) bool {
        \t\t\tfor it := NewIterCursor(b); it.HasNext(); {
        \t\t\t\tout = append(out, it.Next())
        \t\t\t}
        \t\t\treturn true
        \t\t})
        \t\treturn wrapImmutableListOwn(out)
        \t}
        \tmm.keys.ForEachEntry(func(_, b any) bool {
        \t\tif mm.isSet {
        \t\t\tout = append(out, b.(*JMapCore).LiveKeys()...)
        \t\t} else if lc, ok := b.(*JListCore); ok {
        \t\t\tout = append(out, lc.Items...)
        \t\t} else {
        \t\t\tout = append(out, mmUnwrapSingle(b))
        \t\t}
        \t\treturn true
        \t})
        \treturn wrapArrayList(out)
        }
        func multimapOp_asMap_0(mm *jMultimapCore) any {
        \tif mm.immutable {
        \t\tout := NewJMapCore(false)
        \t\tmm.keys.ForEachEntry(func(k, b any) bool {
        \t\t\tout.Put(k, b)
        \t\t\treturn true
        \t\t})
        \t\treturn wrapImmutableMap(out)
        \t}
        \tv := &jMultimapAsMap{mm: mm}
        \tv.BindExtSelf_(v)
        \treturn v
        }

        func multimapInverse(mm *jMultimapCore, isSet bool) any {
        \tinv := newMultimapCore(isSet, false, nil, nil)
        \tif mm.immutable {
        \t\tmm.keys.ForEachEntry(func(k, b any) bool {
        \t\t\tfor it := NewIterCursor(b); it.HasNext(); {
        \t\t\t\tinv.putInternal(it.Next(), k)
        \t\t\t}
        \t\t\treturn true
        \t\t})
        \t} else {
        \t\tmm.keys.ForEachEntry(func(k, b any) bool {
        \t\t\tif mm.isSet {
        \t\t\t\tb.(*JMapCore).ForEachEntry(func(v, _ any) bool {
        \t\t\t\t\tinv.putInternal(v, k)
        \t\t\t\t\treturn true
        \t\t\t\t})
        \t\t\t} else if lc, ok := b.(*JListCore); ok {
        \t\t\t\tfor _, v := range lc.Items {
        \t\t\t\t\tinv.putInternal(v, k)
        \t\t\t\t}
        \t\t\t} else {
        \t\t\t\tinv.putInternal(mmUnwrapSingle(b), k)
        \t\t\t}
        \t\t\treturn true
        \t\t})
        \t}
        \tinv.sealImmutableBuckets()
        \tif isSet {
        \t\tobj := &Ext_com_google_common_collect_ImmutableSetMultimap{}
        \t\tobj.Core_ = inv
        \t\treturn obj
        \t}
        \tobj := &Ext_com_google_common_collect_ImmutableListMultimap{}
        \tobj.Core_ = inv
        \treturn obj
        }

        type jMultimapLiveList struct {
        \tExt_java_util_AbstractListBase
        \tmm      *jMultimapCore
        \tkey     any
        \tstrKey  string
        \tnodeKey uint32
        \tkeyKind uint8
        }

        func (l *jMultimapLiveList) getBucket() any {
        \tswitch l.keyKind {
        \tcase 1:
        \t\treturn l.mm.keys.GetStr(l.strKey)
        \tcase 2:
        \t\treturn l.mm.keys.GetNode(l.nodeKey)
        \tdefault:
        \t\treturn l.mm.keys.Get(l.key)
        \t}
        }
        func (l *jMultimapLiveList) putBucket(b any) {
        \tswitch l.keyKind {
        \tcase 1:
        \t\tl.mm.keys.PutStr(l.strKey, b)
        \tcase 2:
        \t\tl.mm.keys.PutNode(l.nodeKey, b)
        \tdefault:
        \t\tl.mm.keys.Put(l.key, b)
        \t}
        }
        func (l *jMultimapLiveList) removeBucket() {
        \tswitch l.keyKind {
        \tcase 1:
        \t\tl.mm.keys.RemoveStr(l.strKey)
        \tcase 2:
        \t\tl.mm.keys.RemoveNode(l.nodeKey)
        \tdefault:
        \t\tl.mm.keys.Remove(l.key)
        \t}
        }

        func (l *jMultimapLiveList) typedViewLen_() (int, bool) {
        \tif l == nil { ThrowNPE() }
        \treturn int(l.Size()), true
        }
        func (l *jMultimapLiveList) typedViewAt_(i int) any { return l.Get(int32(i)) }
        func (*jMultimapLiveList) typedViewNodeAt_(int) (uint32, bool) { return 0, false }
        func (l *jMultimapLiveList) Size() int32 {
        \tb := l.getBucket()
        \tif b == nil { return 0 }
        \tif lc, ok := b.(*JListCore); ok { return lc.Size() }
        \treturn 1
        }
        func (l *jMultimapLiveList) Get(i int32) any {
        \tb := l.getBucket()
        \tif b == nil { panic(javalang.NewIndexOutOfBoundsExceptionForIndex(int(i), 0)) }
        \tif lc, ok := b.(*JListCore); ok { return lc.Get(i) }
        \tCheckElementIndex(i, 1)
        \treturn mmUnwrapSingle(b)
        }
        func (l *jMultimapLiveList) Set(i int32, v any) any {
        \tb := l.getBucket()
        \tif b == nil { panic(javalang.NewIndexOutOfBoundsExceptionForIndex(int(i), 0)) }
        \tif lc, ok := b.(*JListCore); ok { return lc.Set(i, v) }
        \tCheckElementIndex(i, 1)
        \told := mmUnwrapSingle(b)
        \tl.putBucket(mmWrapSingle(ToAny(v)))
        \treturn old
        }
        func (l *jMultimapLiveList) Add_2(i int32, v any) {
        \tb := l.getBucket()
        \tif b == nil {
        \t\tif i != 0 { panic(javalang.NewIndexOutOfBoundsExceptionForIndex(int(i), 0)) }
        \t\tl.putBucket(mmWrapSingle(ToAny(v)))
        \t\tl.mm.total++
        \t\treturn
        \t}
        \tif lc, ok := b.(*JListCore); ok {
        \t\tlc.AddAt(i, v)
        \t\tl.mm.total++
        \t\treturn
        \t}
        \tCheckPositionIndex(i, 1)
        \tv0 := mmUnwrapSingle(b)
        \tv1 := ToAny(v)
        \tlc := &JListCore{Items: make([]any, 2, 4)}
        \tif i == 0 {
        \t\tlc.Items[0], lc.Items[1] = v1, v0
        \t} else {
        \t\tlc.Items[0], lc.Items[1] = v0, v1
        \t}
        \tl.putBucket(lc)
        \tl.mm.total++
        }
        func (l *jMultimapLiveList) Remove_1_int32(i int32) any {
        \tb := l.getBucket()
        \tif b == nil { panic(javalang.NewIndexOutOfBoundsExceptionForIndex(int(i), 0)) }
        \tif lc, ok := b.(*JListCore); ok {
        \t\told := lc.RemoveAt(i)
        \t\tl.mm.total--
        \t\tif lc.IsEmpty() { l.removeBucket() }
        \t\treturn old
        \t}
        \tCheckElementIndex(i, 1)
        \told := mmUnwrapSingle(b)
        \tl.removeBucket()
        \tl.mm.total--
        \treturn old
        }

        type jMultimapLiveSet struct {
        \tExt_java_util_AbstractSetBase
        \tmm      *jMultimapCore
        \tkey     any
        \tstrKey  string
        \tnodeKey uint32
        \tkeyKind uint8
        }

        func (s *jMultimapLiveSet) getBucket() any {
        \tswitch s.keyKind {
        \tcase 1:
        \t\treturn s.mm.keys.GetStr(s.strKey)
        \tcase 2:
        \t\treturn s.mm.keys.GetNode(s.nodeKey)
        \tdefault:
        \t\treturn s.mm.keys.Get(s.key)
        \t}
        }
        func (s *jMultimapLiveSet) putVal(v any) bool {
        \tswitch s.keyKind {
        \tcase 1:
        \t\treturn s.mm.putStr(s.strKey, v)
        \tcase 2:
        \t\treturn s.mm.putNode(s.nodeKey, v)
        \tdefault:
        \t\treturn multimapOp_put_2(s.mm, s.key, v)
        \t}
        }
        func (s *jMultimapLiveSet) remVal(v any) bool {
        \tswitch s.keyKind {
        \tcase 1:
        \t\treturn s.mm.removeStr(s.strKey, v)
        \tcase 2:
        \t\treturn s.mm.removeNode(s.nodeKey, v)
        \tdefault:
        \t\treturn multimapOp_remove_2(s.mm, s.key, v)
        \t}
        }

        func (*jMultimapLiveSet) Implements_j2gort_Ext_java_util_SequencedCollection() {}
        func (*jMultimapLiveSet) Implements_j2gort_Ext_java_util_SequencedSet() {}
        func (*jMultimapLiveSet) Implements_j2gort_Ext_java_util_SortedSet() {}
        func (*jMultimapLiveSet) Implements_j2gort_Ext_java_util_NavigableSet() {}
        func (s *jMultimapLiveSet) Size() int32 {
        \tb := s.getBucket()
        \tif b == nil { return 0 }
        \treturn b.(*JMapCore).Size()
        }
        func (s *jMultimapLiveSet) Contains(v any) bool {
        \tb := s.getBucket()
        \treturn b != nil && b.(*JMapCore).ContainsKey(v)
        }
        func (s *jMultimapLiveSet) SetContainsStr_(k string) bool {
        \tb := s.getBucket()
        \treturn b != nil && b.(*JMapCore).ContainsStr(k)
        }
        func (s *jMultimapLiveSet) SetAddStr_(k string) bool {
        \tif s.mm.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tb := s.getBucket()
        \tif b == nil {
        \t\tsm := s.mm.newBucket().(*JMapCore)
        \t\tsm.PutStr(k, true)
        \t\tswitch s.keyKind {
        \t\tcase 1:
        \t\t\ts.mm.keys.PutStr(s.strKey, sm)
        \t\tcase 2:
        \t\t\ts.mm.keys.PutNode(s.nodeKey, sm)
        \t\tdefault:
        \t\t\ts.mm.keys.Put(s.key, sm)
        \t\t}
        \t\ts.mm.total++
        \t\treturn true
        \t}
        \tsm := b.(*JMapCore)
        \tif sm.ContainsStr(k) { return false }
        \tsm.PutStr(k, true)
        \ts.mm.total++
        \treturn true
        }
        func (s *jMultimapLiveSet) SetRemoveStr_(k string) bool {
        \tif s.mm.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tb := s.getBucket()
        \tif b == nil { return false }
        \tsm := b.(*JMapCore)
        \tif !sm.ContainsStr(k) { return false }
        \tsm.RemoveStr(k)
        \ts.mm.total--
        \tif sm.IsEmpty() {
        \t\tswitch s.keyKind {
        \t\tcase 1:
        \t\t\ts.mm.keys.RemoveStr(s.strKey)
        \t\tcase 2:
        \t\t\ts.mm.keys.RemoveNode(s.nodeKey)
        \t\tdefault:
        \t\t\ts.mm.keys.Remove(s.key)
        \t\t}
        \t}
        \treturn true
        }
        func (s *jMultimapLiveSet) Add_1(v any) bool           { return s.putVal(v) }
        func (s *jMultimapLiveSet) Remove_1_any(v any) bool    { return s.remVal(v) }
        func (s *jMultimapLiveSet) Clear() {
        \tswitch s.keyKind {
        \tcase 1:
        \t\t_ = s.mm.removeAllStr(s.strKey)
        \tcase 2:
        \t\t_ = s.mm.removeAllNode(s.nodeKey)
        \tdefault:
        \t\t_ = multimapOp_removeAll_1(s.mm, s.key)
        \t}
        }
        func (s *jMultimapLiveSet) AddFirst(v any)             { _ = s.putVal(v) }
        func (s *jMultimapLiveSet) AddLast(v any)              { _ = s.putVal(v) }
        func (s *jMultimapLiveSet) GetFirst() any {
        \tb := s.getBucket()
        \tif b == nil { panic(javalang.NewNoSuchElementException("empty")) }
        \treturn b.(*JMapCore).FirstKey()
        }
        func (s *jMultimapLiveSet) GetLast() any {
        \tb := s.getBucket()
        \tif b == nil { panic(javalang.NewNoSuchElementException("empty")) }
        \treturn b.(*JMapCore).LastKey()
        }
        func (s *jMultimapLiveSet) RemoveFirst() any {
        \tx := s.GetFirst()
        \ts.Remove_1_any(x)
        \treturn x
        }
        func (s *jMultimapLiveSet) RemoveLast() any {
        \tx := s.GetLast()
        \ts.Remove_1_any(x)
        \treturn x
        }
        func (s *jMultimapLiveSet) DescendingSet() any {
        \tb := s.getBucket()
        \tif b == nil { return setCoreDescending(NewSortedJMapCore(nil), false) }
        \treturn setCoreDescending(b.(*JMapCore), false)
        }
        func (s *jMultimapLiveSet) TailSet_2(from any, incl bool) any {
        \tb := s.getBucket()
        \tif b == nil { return setCoreTailSet(NewSortedJMapCore(nil), from, incl, false) }
        \treturn setCoreTailSet(b.(*JMapCore), from, incl, false)
        }
        func (s *jMultimapLiveSet) Iterator() any {
        \tb := s.getBucket()
        \tif b == nil { return newSliceObjIter(nil) }
        \treturn &jMultimapSetIter{mm: s.mm, key: s.key, strKey: s.strKey, nodeKey: s.nodeKey, keyKind: s.keyKind, vals: b.(*JMapCore).LiveKeys()}
        }

        type jMultimapSetIter struct {
        \tmm      *jMultimapCore
        \tkey     any
        \tstrKey  string
        \tnodeKey uint32
        \tvals    []any
        \tidx     int
        \tlastVal any
        \tkeyKind uint8
        \thasLast bool
        }

        func (*jMultimapSetIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jMultimapSetIter) HasNext() bool { return it.idx < len(it.vals) }
        func (it *jMultimapSetIter) Next() any {
        \tif it.idx >= len(it.vals) { panic(javalang.NewNoSuchElementException("exhausted")) }
        \tv := it.vals[it.idx]
        \tit.idx++
        \tit.lastVal = v
        \tit.hasLast = true
        \treturn v
        }
        func (it *jMultimapSetIter) Remove() {
        \tif !it.hasLast { panic(javalang.NewIllegalStateException("remove before next")) }
        \tswitch it.keyKind {
        \tcase 1:
        \t\tit.mm.removeStr(it.strKey, it.lastVal)
        \tcase 2:
        \t\tit.mm.removeNode(it.nodeKey, it.lastVal)
        \tdefault:
        \t\tmultimapOp_remove_2(it.mm, it.key, it.lastVal)
        \t}
        \tit.hasLast = false
        }

        type jMultimapLiveKeySet struct {
        \tExt_java_util_AbstractSetBase
        \tmm *jMultimapCore
        }

        func (s *jMultimapLiveKeySet) Size() int32             { return s.mm.keys.Size() }
        func (s *jMultimapLiveKeySet) Contains(k any) bool     { return s.mm.keys.ContainsKey(k) }
        func (s *jMultimapLiveKeySet) SetContainsStr_(k string) bool { return s.mm.keys.ContainsStr(k) }
        func (s *jMultimapLiveKeySet) SetAddStr_(_ string) bool {
        \tpanic(javalang.NewUnsupportedOperationException("keySet.add"))
        }
        func (s *jMultimapLiveKeySet) SetRemoveStr_(k string) bool {
        \tif !s.mm.keys.ContainsStr(k) { return false }
        \t_ = s.mm.removeAllStr(k)
        \treturn true
        }
        func (s *jMultimapLiveKeySet) Clear()                  { multimapOp_clear_0(s.mm) }
        func (s *jMultimapLiveKeySet) Remove_1_any(k any) bool {
        \tif !s.mm.keys.ContainsKey(k) { return false }
        \t_ = multimapOp_removeAll_1(s.mm, k)
        \treturn true
        }
        func (s *jMultimapLiveKeySet) Iterator() any {
        \treturn &jMultimapKeyIter{mm: s.mm, keys: s.mm.keys.LiveKeys()}
        }

        type jMultimapKeyIter struct {
        \tmm      *jMultimapCore
        \tkeys    []any
        \tidx     int
        \tlastKey any
        \thasLast bool
        }

        func (*jMultimapKeyIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jMultimapKeyIter) HasNext() bool { return it.idx < len(it.keys) }
        func (it *jMultimapKeyIter) Next() any {
        \tif it.idx >= len(it.keys) { panic(javalang.NewNoSuchElementException("exhausted")) }
        \tk := it.keys[it.idx]
        \tit.idx++
        \tit.lastKey = k
        \tit.hasLast = true
        \treturn k
        }
        func (it *jMultimapKeyIter) Remove() {
        \tif !it.hasLast { panic(javalang.NewIllegalStateException("remove before next")) }
        \t_ = multimapOp_removeAll_1(it.mm, it.lastKey)
        \tit.hasLast = false
        }

        type jMultimapAsMap struct {
        \tExt_java_util_AbstractMapBase
        \tmm *jMultimapCore
        }

        func (m *jMultimapAsMap) Size() int32             { return m.mm.keys.Size() }
        func (m *jMultimapAsMap) ContainsKey(k any) bool  { return m.mm.keys.ContainsKey(k) }
        func (m *jMultimapAsMap) MapGetStr_(k string) any {
        \tif !m.mm.keys.ContainsStr(k) { return nil }
        \treturn m.mm.getStr(k)
        }
        func (m *jMultimapAsMap) MapGetOrDefaultStr_(k string, def any) any {
        \tif !m.mm.keys.ContainsStr(k) { return ToAny(def) }
        \treturn m.mm.getStr(k)
        }
        func (m *jMultimapAsMap) MapContainsStr_(k string) bool { return m.mm.keys.ContainsStr(k) }
        func (m *jMultimapAsMap) MapPutStr_(_ string, _ any) any {
        \tpanic(javalang.NewUnsupportedOperationException("asMap.put"))
        }
        func (m *jMultimapAsMap) MapRemoveStr_(k string) any {
        \tif !m.mm.keys.ContainsStr(k) { return nil }
        \treturn m.mm.removeAllStr(k)
        }
        func (m *jMultimapAsMap) Get(k any) any {
        \tif !m.mm.keys.ContainsKey(k) { return nil }
        \treturn multimapOp_get_1(m.mm, k)
        }
        func (m *jMultimapAsMap) Remove_1(k any) any {
        \tif !m.mm.keys.ContainsKey(k) { return nil }
        \treturn multimapOp_removeAll_1(m.mm, k)
        }
        func (m *jMultimapAsMap) Clear() { multimapOp_clear_0(m.mm) }
        func (m *jMultimapAsMap) EntrySet() any {
        \ts := &jMultimapAsMapEntries{mm: m.mm}
        \ts.BindExtSelf_(s)
        \treturn s
        }

        type jMultimapAsMapEntries struct {
        \tExt_java_util_AbstractSetBase
        \tmm *jMultimapCore
        }

        func (s *jMultimapAsMapEntries) Size() int32 { return s.mm.keys.Size() }
        func (s *jMultimapAsMapEntries) Clear()      { multimapOp_clear_0(s.mm) }
        func (s *jMultimapAsMapEntries) Iterator() any {
        \tvar items []any
        \tfor _, k := range s.mm.keys.LiveKeys() {
        \t\titems = append(items, &jSimpleEntry{k: k, v: multimapOp_get_1(s.mm, k), imm: true})
        \t}
        \treturn &jMultimapAsMapIter{mm: s.mm, entries: items}
        }

        type jMultimapAsMapIter struct {
        \tmm      *jMultimapCore
        \tentries []any
        \tidx     int
        \tlastKey any
        \thasLast bool
        }

        func (*jMultimapAsMapIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jMultimapAsMapIter) HasNext() bool { return it.idx < len(it.entries) }
        func (it *jMultimapAsMapIter) Next() any {
        \tif it.idx >= len(it.entries) { panic(javalang.NewNoSuchElementException("exhausted")) }
        \te := it.entries[it.idx].(*jSimpleEntry)
        \tit.idx++
        \tit.lastKey = e.k
        \tit.hasLast = true
        \treturn e
        }
        func (it *jMultimapAsMapIter) Remove() {
        \tif !it.hasLast { panic(javalang.NewIllegalStateException("remove before next")) }
        \t_ = multimapOp_removeAll_1(it.mm, it.lastKey)
        \tit.hasLast = false
        }

        type jMultimapBuilderEntry struct {
        \tstrKey string
        \tkey    any
        \tval    any
        \tisStr  bool
        }

        type jMultimapBuilderCore struct {
        \tisSet   bool
        \tentries []jMultimapBuilderEntry
        \tkeyCmp  func(a, b any) int32
        }

        func newMultimapBuilderCore(isSet bool) *jMultimapBuilderCore {
        \treturn &jMultimapBuilderCore{isSet: isSet}
        }

        func multimapBuilderPut(b *jMultimapBuilderCore, k, v any) {
        \tif s, ok := k.(string); ok {
        \t\tif s == NullString { ThrowNPE() }
        \t\tb.entries = append(b.entries, jMultimapBuilderEntry{strKey: s, val: CheckNotNull(ToAny(v)), isStr: true})
        \t\treturn
        \t}
        \tb.entries = append(b.entries, jMultimapBuilderEntry{key: CheckNotNull(ToAny(k)), val: CheckNotNull(ToAny(v))})
        }

        func MultimapBuilderPutStr(b any, k string, v any) any {
        \tif b == nil || k == NullString { ThrowNPE() }
        \tv = CheckNotNull(ToAny(v))
        \tif cp, ok := b.(interface{ getCore() any }); ok {
        \t\tif isJavaNilAny(b) { ThrowNPE() }
        \t\tbc := cp.getCore().(*jMultimapBuilderCore)
        \t\tbc.entries = append(bc.entries, jMultimapBuilderEntry{strKey: k, val: v, isStr: true})
        \t\treturn b
        \t}
        \tif isJavaNilAny(b) { ThrowNPE() }
        \tpanic(javalang.NewUnsupportedOperationException(fmt.Sprintf("MultimapBuilderPutStr unsupported on %T", b)))
        }

        func multimapBuilderPutAll(b *jMultimapBuilderCore, k any, vals Ext_java_lang_Iterable) {
        \tif s, ok := k.(string); ok {
        \t\tif s == NullString { ThrowNPE() }
        \t\tCheckNotNull(vals)
        \t\tfor it := NewIterCursor(vals); it.HasNext(); {
        \t\t\tb.entries = append(b.entries, jMultimapBuilderEntry{strKey: s, val: CheckNotNull(it.Next()), isStr: true})
        \t\t}
        \t\treturn
        \t}
        \tk = CheckNotNull(ToAny(k))
        \tCheckNotNull(vals)
        \tfor it := NewIterCursor(vals); it.HasNext(); {
        \t\tb.entries = append(b.entries, jMultimapBuilderEntry{key: k, val: CheckNotNull(it.Next())})
        \t}
        }

        func multimapBuilderBuild(b *jMultimapBuilderCore) any {
        \tmm := newMultimapCore(b.isSet, false, b.keyCmp, nil)
        \tfor _, e := range b.entries {
        \t\tif e.isStr {
        \t\t\tmm.putStrInternal(e.strKey, e.val)
        \t\t} else {
        \t\t\tmm.putInternal(e.key, e.val)
        \t\t}
        \t}
        \tmm.sealImmutableBuckets()
        \tif b.isSet {
        \t\tobj := &Ext_com_google_common_collect_ImmutableSetMultimap{}
        \t\tobj.Core_ = mm
        \t\treturn obj
        \t}
        \tobj := &Ext_com_google_common_collect_ImmutableListMultimap{}
        \tobj.Core_ = mm
        \treturn obj
        }

        func immSetMultimapCopyOf(src Ext_com_google_common_collect_Multimap) *Ext_com_google_common_collect_ImmutableSetMultimap {
        \tCheckNotNull(src)
        \tif ism, ok := any(src).(*Ext_com_google_common_collect_ImmutableSetMultimap); ok && ism != nil {
        \t\treturn ism
        \t}
        \tmm := newMultimapCore(true, false, nil, nil)
        \tasMap := Cast[Ext_java_util_Map](src.AsMap())
        \tfor it := Iterate(asMap.EntrySet()); it.HasNext(); {
        \t\te := Cast[Ext_java_util_Map_Entry](it.Next())
        \t\tk := CheckNotNull(e.GetKey())
        \t\tfor vit := NewIterCursor(e.GetValue()); vit.HasNext(); {
        \t\t\tmm.putInternal(k, CheckNotNull(vit.Next()))
        \t\t}
        \t}
        \tmm.sealImmutableBuckets()
        \tobj := &Ext_com_google_common_collect_ImmutableSetMultimap{}
        \tobj.Core_ = mm
        \treturn obj
        }

        func multimapsIndex(vals Ext_java_lang_Iterable, keyFn Ext_com_google_common_base_Function) *Ext_com_google_common_collect_ImmutableListMultimap {
        \tCheckNotNull(vals)
        \tCheckNotNull(keyFn)
        \tmm := newMultimapCore(false, false, nil, nil)
        \tfor it := NewIterCursor(vals); it.HasNext(); {
        \t\tv := CheckNotNull(it.Next())
        \t\tk := CheckNotNull(ToAny(keyFn.Apply(v)))
        \t\tmm.putInternal(k, v)
        \t}
        \tmm.sealImmutableBuckets()
        \tobj := &Ext_com_google_common_collect_ImmutableListMultimap{}
        \tobj.Core_ = mm
        \treturn obj
        }

        // --- Multiset ---

        type jMultisetCore struct {
        \tcounts    *JMapCore // elem -> int32
        \ttotal     int32
        \timmutable bool
        \tcmp       func(a, b any) int32
        }

        func newMultisetCore(immutable bool, cmp func(a, b any) int32) *jMultisetCore {
        \tvar cm *JMapCore
        \tif cmp != nil {
        \t\tcm = NewSortedJMapCore(cmp)
        \t} else {
        \t\tcm = NewJMapCore(false)
        \t}
        \treturn &jMultisetCore{counts: cm, immutable: immutable, cmp: cmp}
        }

        func newTreeMultisetCore() *jMultisetCore {
        \treturn &jMultisetCore{counts: NewSortedJMapCore(nil), cmp: CompareValues}
        }

        func cloneMultisetCore(src *jMultisetCore, immutable bool) *jMultisetCore {
        \tdst := newMultisetCore(immutable, src.cmp)
        \tdst.counts = cloneJMapCore(src.counts)
        \tdst.total = src.total
        \treturn dst
        }

        func wrapImmutableMultiset(c *jMultisetCore) Ext_com_google_common_collect_ImmutableMultiset {
        \tobj := &Ext_com_google_common_collect_ImmutableMultisetBase{}
        \tobj.BindExtSelf_(obj)
        \tobj.Core_ = c
        \treturn obj
        }

        type jMultisetEntry struct {
        \telem  any
        \tcount int32
        }

        func (*jMultisetEntry) Implements_j2gort_Ext_com_google_common_collect_Multiset_Entry() {}
        func (e *jMultisetEntry) GetElement() any { return e.elem }
        func (e *jMultisetEntry) GetCount() int32 { return e.count }
        func (e *jMultisetEntry) Equals(o any) bool {
        \tother, ok := o.(Ext_com_google_common_collect_Multiset_Entry)
        \treturn ok && other != nil && !javalang.IsNil(other) &&
        \t\te.count == other.GetCount() && javalang.ObjectsEquals(e.elem, ToAny(other.GetElement()))
        }
        func (e *jMultisetEntry) HashCode() int32 {
        \treturn javalang.ObjectHashCode(e.elem) ^ e.count
        }
        func (e *jMultisetEntry) String() string {
        \tif e.count == 1 { return StrVal(e.elem) }
        \treturn fmt.Sprintf("%s x %d", StrVal(e.elem), e.count)
        }

        func multisetCountStr(c *jMultisetCore, s string) int32 {
        \tv := c.counts.GetStr(s)
        \tif v == nil { return 0 }
        \treturn v.(int32)
        }
        func multisetAddStr2(c *jMultisetCore, s string, n int32) int32 {
        \tif c.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tCheckArgument(n >= 0)
        \told := multisetCountStr(c, s)
        \tif n > 0 {
        \t\tc.counts.PutStr(s, old+n)
        \t\tc.total += n
        \t}
        \treturn old
        }
        func multisetRemoveStr1(c *jMultisetCore, s string) bool {
        \tif c.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \told := multisetCountStr(c, s)
        \tif old == 0 { return false }
        \tif old == 1 {
        \t\tc.counts.RemoveStr(s)
        \t} else {
        \t\tc.counts.PutStr(s, old-1)
        \t}
        \tc.total--
        \treturn true
        }
        func multisetSetCountStr3(c *jMultisetCore, s string, oldCount, newCount int32) bool {
        \tif c.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tCheckArgument(oldCount >= 0 && newCount >= 0)
        \tcur := multisetCountStr(c, s)
        \tif cur != oldCount { return false }
        \tif newCount == 0 {
        \t\tif cur > 0 { c.counts.RemoveStr(s) }
        \t} else {
        \t\tc.counts.PutStr(s, newCount)
        \t}
        \tc.total += newCount - cur
        \treturn true
        }
        func multisetOp_size_0(c *jMultisetCore) int32   { return c.total }
        func multisetOp_isEmpty_0(c *jMultisetCore) bool { return c.total == 0 }
        func multisetOp_count_1(c *jMultisetCore, e any) int32 {
        \tv := c.counts.Get(e)
        \tif v == nil { return 0 }
        \treturn v.(int32)
        }
        func multisetOp_contains_1(c *jMultisetCore, e any) bool { return multisetOp_count_1(c, e) > 0 }
        func multisetOp_add_2(c *jMultisetCore, e any, n int32) int32 {
        \tif c.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tCheckArgument(n >= 0)
        \told := multisetOp_count_1(c, e)
        \tif n > 0 {
        \t\tc.counts.Put(e, old+n)
        \t\tc.total += n
        \t}
        \treturn old
        }
        func multisetOp_add_1(c *jMultisetCore, e any) bool {
        \tmultisetOp_add_2(c, e, 1)
        \treturn true
        }
        func multisetOp_remove_1(c *jMultisetCore, e any) bool {
        \tif c.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \told := multisetOp_count_1(c, e)
        \tif old == 0 { return false }
        \tif old == 1 {
        \t\tc.counts.Remove(e)
        \t} else {
        \t\tc.counts.Put(e, old-1)
        \t}
        \tc.total--
        \treturn true
        }
        func multisetOp_setCount_3(c *jMultisetCore, e any, oldCount, newCount int32) bool {
        \tif c.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tCheckArgument(oldCount >= 0 && newCount >= 0)
        \tcur := multisetOp_count_1(c, e)
        \tif cur != oldCount { return false }
        \tif newCount == 0 {
        \t\tif cur > 0 { c.counts.Remove(e) }
        \t} else {
        \t\tc.counts.Put(e, newCount)
        \t}
        \tc.total += newCount - cur
        \treturn true
        }
        func multisetOp_addAll_1(c *jMultisetCore, col Ext_java_util_Collection) bool {
        \tCheckNotNull(col)
        \tchanged := false
        \tfor it := Iterate(col); it.HasNext(); {
        \t\tmultisetOp_add_1(c, it.Next())
        \t\tchanged = true
        \t}
        \treturn changed
        }
        func multisetOp_removeAll_1(c *jMultisetCore, col Ext_java_util_Collection) bool {
        \tif c.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tCheckNotNull(col)
        \tchanged := false
        \tfor it := Iterate(col); it.HasNext(); {
        \t\tx := it.Next()
        \t\tif cnt := multisetOp_count_1(c, x); cnt > 0 {
        \t\t\tc.counts.Remove(x)
        \t\t\tc.total -= cnt
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }
        func multisetOp_clear_0(c *jMultisetCore) {
        \tif c.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tc.counts.Clear()
        \tc.total = 0
        }
        func multisetOp_retainAll_1(c *jMultisetCore, col Ext_java_util_Collection) bool {
        \tCheckNotNull(col)
        \tif c.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tchanged := false
        \tfor _, x := range c.counts.LiveKeys() {
        \t\tif !col.Contains(x) {
        \t\t\tcnt := multisetOp_count_1(c, x)
        \t\t\tc.counts.Remove(x)
        \t\t\tc.total -= cnt
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }
        func multisetOp_removeIf_1(c *jMultisetCore, pred Ext_java_util_function_Predicate) bool {
        \tCheckNotNull(pred)
        \tif c.immutable { panic(javalang.NewUnsupportedOperationException("immutable")) }
        \tchanged := false
        \tfor _, x := range c.counts.LiveKeys() {
        \t\tif pred.Test(x) {
        \t\t\tcnt := multisetOp_count_1(c, x)
        \t\t\tc.counts.Remove(x)
        \t\t\tc.total -= cnt
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }
        func multisetOp_elementSet_0(c *jMultisetCore) any {
        \tif c.immutable { return wrapImmutableSet(c.counts.LiveKeys()) }
        \tif c.cmp != nil {
        \t\tts := &Ext_java_util_TreeSet{}; ts.BindExtSelf_(ts); ts.Core_ = cloneJMapCore(c.counts); return ts
        \t}
        \treturn newLiveMapKeySet(c.counts)
        }
        func multisetOp_entrySet_0(c *jMultisetCore) any {
        \tentries := make([]any, 0, c.counts.Size())
        \tc.counts.ForEachEntry(func(k, v any) bool {
        \t\tentries = append(entries, &jMultisetEntry{elem: k, count: v.(int32)})
        \t\treturn true
        \t})
        \tif c.immutable { return wrapImmutableSet(entries) }
        \treturn wrapLinkedHashSet(entries)
        }
        func multisetOp_toSlice(c *jMultisetCore) []any {
        \tout := make([]any, 0, c.total)
        \tc.counts.ForEachEntry(func(e, v any) bool {
        \t\tn := v.(int32)
        \t\tfor j := int32(0); j < n; j++ {
        \t\t\tout = append(out, e)
        \t\t}
        \t\treturn true
        \t})
        \treturn out
        }
        func multisetOp_iterator_0(c *jMultisetCore) any {
        \treturn &jMultisetIter{c: c, items: multisetOp_toSlice(c)}
        }
        func multisetOp_asList_0(c *jMultisetCore) Ext_com_google_common_collect_ImmutableList {
        \treturn wrapImmutableList(multisetOp_toSlice(c))
        }
        func multisetOp_toArray_0(c *jMultisetCore) []any             { return multisetOp_toSlice(c) }
        func multisetOp_toArray_1(c *jMultisetCore, d []any) []any    { return sliceToArray1(multisetOp_toSlice(c), d) }
        func multisetOp_containsAll_1(c *jMultisetCore, col Ext_java_util_Collection) bool {
        \treturn colContainsAll(func(x any) bool { return multisetOp_contains_1(c, x) }, col)
        }
        func multisetOp_forEach_1(c *jMultisetCore, fn Ext_java_util_function_Consumer) {
        \tCheckNotNull(fn)
        \tfor _, x := range multisetOp_toSlice(c) { fn.Accept(x) }
        }
        func multisetOp_stream_0(c *jMultisetCore) Ext_java_util_stream_Stream {
        \treturn streamFromSlice(multisetOp_toSlice(c))
        }
        func multisetOp_equals_1(c *jMultisetCore, o any) bool {
        \tms, ok := o.(Ext_com_google_common_collect_Multiset)
        \tif !ok || ms == nil || javalang.IsNil(ms) { return false }
        \tif c.total != ms.Size() { return false }
        \tes := Cast[Ext_java_util_Set](ms.EntrySet())
        \tif c.counts.Size() != es.Size() { return false }
        \tfor it := Iterate(es); it.HasNext(); {
        \t\te := Cast[Ext_com_google_common_collect_Multiset_Entry](it.Next())
        \t\tif multisetOp_count_1(c, e.GetElement()) != e.GetCount() { return false }
        \t}
        \treturn true
        }
        func multisetOp_hashCode_0(c *jMultisetCore) int32 {
        \th := int32(0)
        \tc.counts.ForEachEntry(func(k, v any) bool {
        \t\th += javalang.ObjectHashCode(k) ^ v.(int32)
        \t\treturn true
        \t})
        \treturn h
        }
        func multisetOp_toString_0(c *jMultisetCore) string {
        \tvar sb strings.Builder
        \tsb.WriteByte('[')
        \tfirst := true
        \tc.counts.ForEachEntry(func(k, v any) bool {
        \t\tif !first { sb.WriteString(", ") }
        \t\tfirst = false
        \t\te := &jMultisetEntry{elem: k, count: v.(int32)}
        \t\tsb.WriteString(e.String())
        \t\treturn true
        \t})
        \tsb.WriteByte(']')
        \treturn sb.String()
        }

        type jMultisetIter struct {
        \tc       *jMultisetCore
        \titems   []any
        \tidx     int
        \tlastVal any
        \thasLast bool
        }

        func (*jMultisetIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (*jMultisetIter) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator() {}
        func (*jMultisetIter) GetExt_com_google_common_collect_UnmodifiableIteratorBase() *Ext_com_google_common_collect_UnmodifiableIteratorBase {
        \treturn nil
        }
        func (it *jMultisetIter) HasNext() bool { return it.idx < len(it.items) }
        func (it *jMultisetIter) Next() any {
        \tif it.idx >= len(it.items) { panic(javalang.NewNoSuchElementException("exhausted")) }
        \tv := it.items[it.idx]
        \tit.idx++
        \tit.lastVal = v
        \tit.hasLast = true
        \treturn v
        }
        func (it *jMultisetIter) Remove() {
        \tif !it.hasLast { panic(javalang.NewIllegalStateException("remove before next")) }
        \tmultisetOp_remove_1(it.c, it.lastVal)
        \tit.hasLast = false
        }
        """;
  }

  private static String buildColRtPart2b() {
    return """

        // --- Table (HashBasedTable, StandardTable, Tables.newCustomTable) ---

        type jFastTableCoreProvider interface {
        \tfastTableCore_() (*jTableCore, bool)
        }

        type jTableRowEntry struct {
        \trowKey    any
        \trowKeyI64 int64
        \tcol0      any
        \tval0      any
        \tmulti     *JMapCore
        \thashNext  int32
        \thash      uint32
        \thasSingle bool
        \tdeleted   bool
        \tisColorId bool
        }

        func (e *jTableRowEntry) getRowKey() any {
        \tif e.isColorId {
        \t\tif e.rowKey == nil {
        \t\t\te.rowKey = boxColorId(e.rowKeyI64)
        \t\t}
        \t\treturn e.rowKey
        \t}
        \treturn e.rowKey
        }

        func (e *jTableRowEntry) matchesAnyRowKey(r any, fk FastCustomKey, h uint32) bool {
        \tif e.hash != h {
        \t\treturn false
        \t}
        \tif e.isColorId {
        \t\tif rc, ok := r.(jRawColorIdProvider); ok {
        \t\t\traw, valid := rc.RawColorId_()
        \t\t\treturn valid && raw == e.rowKeyI64
        \t\t}
        \t\treturn false
        \t}
        \tif fk != nil {
        \t\treturn fk.Equals(e.rowKey)
        \t}
        \treturn javalang.ObjectsEquals(e.rowKey, r)
        }

        func (e *jTableRowEntry) isEmptyRow() bool {
        \treturn e == nil || e.deleted || (!e.hasSingle && (e.multi == nil || e.multi.IsEmpty()))
        }

        type jTableCore struct {
        \tentries    []jTableRowEntry
        \tbuckets    []int32
        \tstrIdx     map[string]int32
        \trowFactory func() *JMapCore
        \tliveRows   int32
        \tstrLive    int32
        \tnonStrLive int32
        \tmodVer     uint32
        }

        func newTableCore(rowFactory func() *JMapCore) *jTableCore {
        \treturn &jTableCore{rowFactory: rowFactory, modVer: 1}
        }

        func newCustomTable(backingMap Ext_java_util_Map, factory Ext_com_google_common_base_Supplier) Ext_com_google_common_collect_Table {
        \tCheckNotNull(backingMap)
        \tCheckNotNull(factory)
        \trf := func() *JMapCore {
        \t\tm := factory.Get()
        \t\tif mh, ok := m.(interface{ mc() *JMapCore }); ok {
        \t\t\treturn mh.mc()
        \t\t}
        \t\treturn NewJMapCore(false)
        \t}
        \tc := newTableCore(rf)
        \tif !backingMap.IsEmpty() {
        \t\tfor cur := NewMapEntryCursor(backingMap); cur.NextEntry(); {
        \t\t\tr := cur.GetKey()
        \t\t\tif rm, ok := cur.GetValue().(Ext_java_util_Map); ok && rm != nil && !javalang.IsNil(rm) {
        \t\t\t\tfor rc := NewMapEntryCursor(rm); rc.NextEntry(); {
        \t\t\t\t\ttableOp_put_3(c, r, rc.GetKey(), rc.GetValue())
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t}
        \tst := &Ext_com_google_common_collect_StandardTableBase{}
        \tst.BindExtSelf_(st)
        \tst.Core_ = c
        \treturn st
        }

        func (t *jTableCore) newMultiRow() *JMapCore {
        \tif t.rowFactory != nil {
        \t\treturn t.rowFactory()
        \t}
        \treturn NewJMapCore(false)
        }

        func (t *jTableCore) rebuildBuckets(capB int) {
        \tnb := make([]int32, capB)
        \tmask := uint32(capB - 1)
        \tfor i := range t.entries {
        \t\te := &t.entries[i]
        \t\tif e.deleted {
        \t\t\tcontinue
        \t\t}
        \t\tif !e.isColorId {
        \t\t\tif _, isStr := e.rowKey.(string); isStr {
        \t\t\t\tcontinue
        \t\t\t}
        \t\t}
        \t\tslot := smearHash(int32(e.hash)) & mask
        \t\te.hashNext = nb[slot]
        \t\tnb[slot] = int32(i + 1)
        \t}
        \tt.buckets = nb
        }

        func (t *jTableCore) rebuildStrIdx() {
        \tm := make(map[string]int32, max(16, int(t.strLive)*2))
        \tfor i := range t.entries {
        \t\te := &t.entries[i]
        \t\tif e.deleted || e.isColorId {
        \t\t\tcontinue
        \t\t}
        \t\tif s, isStr := e.rowKey.(string); isStr {
        \t\t\tm[s] = int32(i)
        \t\t}
        \t}
        \tt.strIdx = m
        }

        func (t *jTableCore) findRowIdxColorId(raw int64) int {
        \tif len(t.entries) == 0 {
        \t\treturn -1
        \t}
        \th := uint32(int32(raw))
        \tif len(t.buckets) > 0 {
        \t\tslot := smearHash(int32(h)) & uint32(len(t.buckets)-1)
        \t\tfor pos := t.buckets[slot]; pos != 0; {
        \t\t\te := &t.entries[pos-1]
        \t\t\tif e.hash == h {
        \t\t\t\tif e.isColorId {
        \t\t\t\t\tif e.rowKeyI64 == raw {
        \t\t\t\t\t\treturn int(pos - 1)
        \t\t\t\t\t}
        \t\t\t\t} else if rc, ok := e.rowKey.(jRawColorIdProvider); ok {
        \t\t\t\t\tif eraw, valid := rc.RawColorId_(); valid && eraw == raw {
        \t\t\t\t\t\treturn int(pos - 1)
        \t\t\t\t\t}
        \t\t\t\t}
        \t\t\t}
        \t\t\tpos = e.hashNext
        \t\t}
        \t\treturn -1
        \t}
        \tfor i := range t.entries {
        \t\te := &t.entries[i]
        \t\tif !e.deleted && e.hash == h {
        \t\t\tif e.isColorId {
        \t\t\t\tif e.rowKeyI64 == raw {
        \t\t\t\t\treturn i
        \t\t\t\t}
        \t\t\t} else if rc, ok := e.rowKey.(jRawColorIdProvider); ok {
        \t\t\t\tif eraw, valid := rc.RawColorId_(); valid && eraw == raw {
        \t\t\t\t\treturn i
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t}
        \treturn -1
        }

        func (t *jTableCore) findRowIdx(r any) int {
        \tif len(t.entries) == 0 || r == nil {
        \t\treturn -1
        \t}
        \tif s, isStr := r.(string); isStr {
        \t\tif s == NullString {
        \t\t\treturn -1
        \t\t}
        \t\tif t.strIdx != nil {
        \t\t\tif idx, ok := t.strIdx[s]; ok {
        \t\t\t\treturn int(idx)
        \t\t\t}
        \t\t\treturn -1
        \t\t}
        \t\tfor i := range t.entries {
        \t\t\te := &t.entries[i]
        \t\t\tif !e.deleted && !e.isColorId {
        \t\t\t\tif es, ok := e.rowKey.(string); ok && es == s {
        \t\t\t\t\treturn i
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t\treturn -1
        \t}
        \tif rc, ok := r.(jRawColorIdProvider); ok {
        \t\traw, valid := rc.RawColorId_()
        \t\tif !valid {
        \t\t\treturn -1
        \t\t}
        \t\treturn t.findRowIdxColorId(raw)
        \t}
        \tif fk, ok := r.(FastCustomKey); ok {
        \t\tif fk.IsJavaNil_() {
        \t\t\treturn -1
        \t\t}
        \t\th := uint32(fk.HashCode())
        \t\tif len(t.buckets) > 0 {
        \t\t\tslot := smearHash(int32(h)) & uint32(len(t.buckets)-1)
        \t\t\tfor pos := t.buckets[slot]; pos != 0; {
        \t\t\t\te := &t.entries[pos-1]
        \t\t\t\tif e.matchesAnyRowKey(r, fk, h) {
        \t\t\t\t\treturn int(pos - 1)
        \t\t\t\t}
        \t\t\t\tpos = e.hashNext
        \t\t\t}
        \t\t\treturn -1
        \t\t}
        \t\tfor i := range t.entries {
        \t\t\te := &t.entries[i]
        \t\t\tif !e.deleted && e.matchesAnyRowKey(r, fk, h) {
        \t\t\t\treturn i
        \t\t\t}
        \t\t}
        \t\treturn -1
        \t}
        \tif javalang.IsNil(r) {
        \t\treturn -1
        \t}
        \th := uint32(elemJavaHashCode(r, false))
        \tif len(t.buckets) > 0 {
        \t\tslot := smearHash(int32(h)) & uint32(len(t.buckets)-1)
        \t\tfor pos := t.buckets[slot]; pos != 0; {
        \t\t\te := &t.entries[pos-1]
        \t\t\tif e.matchesAnyRowKey(r, nil, h) {
        \t\t\t\treturn int(pos - 1)
        \t\t\t}
        \t\t\tpos = e.hashNext
        \t\t}
        \t\treturn -1
        \t}
        \tfor i := range t.entries {
        \t\te := &t.entries[i]
        \t\tif !e.deleted && e.matchesAnyRowKey(r, nil, h) {
        \t\t\treturn i
        \t\t}
        \t}
        \treturn -1
        }

        func (t *jTableCore) insertNewRowColorId(raw int64, c, v any) {
        \tif t.entries == nil {
        \t\tt.entries = make([]jTableRowEntry, 0, 4)
        \t}
        \tidx := len(t.entries)
        \th := uint32(int32(raw))
        \tif t.rowFactory != nil {
        \t\tmc := t.rowFactory()
        \t\tmc.Put(c, v)
        \t\tt.entries = append(t.entries, jTableRowEntry{rowKeyI64: raw, isColorId: true, multi: mc, hash: h})
        \t} else {
        \t\tt.entries = append(t.entries, jTableRowEntry{rowKeyI64: raw, isColorId: true, col0: c, val0: v, hash: h, hasSingle: true})
        \t}
        \tt.liveRows++
        \tt.nonStrLive++
        \tif len(t.buckets) == 0 {
        \t\tif t.nonStrLive >= 8 {
        \t\t\tt.rebuildBuckets(16)
        \t\t}
        \t\treturn
        \t}
        \tif int(t.nonStrLive)*4 >= len(t.buckets)*3 {
        \t\tt.rebuildBuckets(len(t.buckets) * 2)
        \t\treturn
        \t}
        \tslot := smearHash(int32(h)) & uint32(len(t.buckets)-1)
        \tt.entries[idx].hashNext = t.buckets[slot]
        \tt.buckets[slot] = int32(idx + 1)
        }

        func (t *jTableCore) insertNewRow(r, c, v any) {
        \tif rc, ok := r.(jRawColorIdProvider); ok {
        \t\tif raw, valid := rc.RawColorId_(); valid {
        \t\t\tt.insertNewRowColorId(raw, c, v)
        \t\t\treturn
        \t\t}
        \t}
        \tif t.entries == nil {
        \t\tt.entries = make([]jTableRowEntry, 0, 4)
        \t}
        \tidx := len(t.entries)
        \tif s, isStr := r.(string); isStr {
        \t\tif t.rowFactory != nil {
        \t\t\tmc := t.rowFactory()
        \t\t\tmc.Put(c, v)
        \t\t\tt.entries = append(t.entries, jTableRowEntry{rowKey: r, multi: mc})
        \t\t} else {
        \t\t\tt.entries = append(t.entries, jTableRowEntry{rowKey: r, col0: c, val0: v, hasSingle: true})
        \t\t}
        \t\tt.liveRows++
        \t\tt.strLive++
        \t\tif t.strIdx != nil {
        \t\t\tt.strIdx[s] = int32(idx)
        \t\t} else if t.strLive >= 8 {
        \t\t\tt.rebuildStrIdx()
        \t\t}
        \t\treturn
        \t}
        \tvar h uint32
        \tif fk, ok := r.(FastCustomKey); ok {
        \t\th = uint32(fk.HashCode())
        \t} else {
        \t\th = uint32(elemJavaHashCode(r, false))
        \t}
        \tif t.rowFactory != nil {
        \t\tmc := t.rowFactory()
        \t\tmc.Put(c, v)
        \t\tt.entries = append(t.entries, jTableRowEntry{rowKey: r, multi: mc, hash: h})
        \t} else {
        \t\tt.entries = append(t.entries, jTableRowEntry{rowKey: r, col0: c, val0: v, hash: h, hasSingle: true})
        \t}
        \tt.liveRows++
        \tt.nonStrLive++
        \tif len(t.buckets) == 0 {
        \t\tif t.nonStrLive >= 8 {
        \t\t\tt.rebuildBuckets(16)
        \t\t}
        \t\treturn
        \t}
        \tif int(t.nonStrLive)*4 >= len(t.buckets)*3 {
        \t\tt.rebuildBuckets(len(t.buckets) * 2)
        \t\treturn
        \t}
        \tslot := smearHash(int32(h)) & uint32(len(t.buckets)-1)
        \tt.entries[idx].hashNext = t.buckets[slot]
        \tt.buckets[slot] = int32(idx + 1)
        }

        func (t *jTableCore) deleteRowAt(idx int) {
        \tif idx < 0 || idx >= len(t.entries) {
        \t\treturn
        \t}
        \te := &t.entries[idx]
        \tif e.deleted {
        \t\treturn
        \t}
        \tif s, isStr := e.rowKey.(string); isStr {
        \t\tif t.strIdx != nil {
        \t\t\tdelete(t.strIdx, s)
        \t\t}
        \t\tt.strLive--
        \t} else {
        \t\tif len(t.buckets) > 0 {
        \t\t\tpos := int32(idx + 1)
        \t\t\tslot := smearHash(int32(e.hash)) & uint32(len(t.buckets)-1)
        \t\t\tvar prev int32
        \t\t\tfor cur := t.buckets[slot]; cur != 0; {
        \t\t\t\tif cur == pos {
        \t\t\t\t\tif prev == 0 {
        \t\t\t\t\t\tt.buckets[slot] = e.hashNext
        \t\t\t\t\t} else {
        \t\t\t\t\t\tt.entries[prev-1].hashNext = e.hashNext
        \t\t\t\t\t}
        \t\t\t\t\tbreak
        \t\t\t\t}
        \t\t\t\tprev = cur
        \t\t\t\tcur = t.entries[cur-1].hashNext
        \t\t\t}
        \t\t}
        \t\tt.nonStrLive--
        \t}
        \t*e = jTableRowEntry{deleted: true}
        \tt.liveRows--
        \tif t.liveRows == 0 {
        \t\tt.entries = t.entries[:0]
        \t\tif len(t.buckets) > 0 {
        \t\t\tclear(t.buckets)
        \t\t}
        \t\tif t.strIdx != nil {
        \t\t\tclear(t.strIdx)
        \t\t}
        \t\tt.modVer++
        \t\treturn
        \t}
        \tif len(t.entries) >= 32 && int(t.liveRows)*2 < len(t.entries) {
        \t\tw := 0
        \t\tfor i := range t.entries {
        \t\t\tif !t.entries[i].deleted {
        \t\t\t\tif w != i {
        \t\t\t\t\tt.entries[w] = t.entries[i]
        \t\t\t\t}
        \t\t\t\tw++
        \t\t\t}
        \t\t}
        \t\tclear(t.entries[w:])
        \t\tt.entries = t.entries[:w]
        \t\tt.modVer++
        \t\tif t.strIdx != nil {
        \t\t\tt.rebuildStrIdx()
        \t\t}
        \t\tif len(t.buckets) > 0 {
        \t\t\tt.rebuildBuckets(len(t.buckets))
        \t\t}
        \t}
        }

        func tableOp_size_0(t *jTableCore) int32 {
        \tvar n int32
        \tfor i := range t.entries {
        \t\te := &t.entries[i]
        \t\tif e.deleted {
        \t\t\tcontinue
        \t\t}
        \t\tif e.hasSingle {
        \t\t\tn++
        \t\t} else if e.multi != nil {
        \t\t\tn += e.multi.Size()
        \t\t}
        \t}
        \treturn n
        }

        func tableOp_isEmpty_0(t *jTableCore) bool {
        \tif t.liveRows == 0 {
        \t\treturn true
        \t}
        \tfor i := range t.entries {
        \t\te := &t.entries[i]
        \t\tif !e.isEmptyRow() {
        \t\t\treturn false
        \t\t}
        \t}
        \treturn true
        }

        func tableOp_containsRow_1(t *jTableCore, r any) bool {
        \tidx := t.findRowIdx(r)
        \treturn idx >= 0 && !t.entries[idx].isEmptyRow()
        }

        func tableOp_containsColumn_1(t *jTableCore, c any) bool {
        \tif c == nil || javalang.IsNil(c) {
        \t\treturn false
        \t}
        \tfor i := range t.entries {
        \t\te := &t.entries[i]
        \t\tif e.deleted {
        \t\t\tcontinue
        \t\t}
        \t\tif e.hasSingle {
        \t\t\tif javalang.ObjectsEquals(e.col0, c) {
        \t\t\t\treturn true
        \t\t\t}
        \t\t} else if e.multi != nil && e.multi.ContainsKey(c) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func tableOp_containsValue_1(t *jTableCore, v any) bool {
        \tif v == nil || javalang.IsNil(v) {
        \t\treturn false
        \t}
        \tfor i := range t.entries {
        \t\te := &t.entries[i]
        \t\tif e.deleted {
        \t\t\tcontinue
        \t\t}
        \t\tif e.hasSingle {
        \t\t\tif javalang.ObjectsEquals(e.val0, v) {
        \t\t\t\treturn true
        \t\t\t}
        \t\t} else if e.multi != nil && e.multi.ContainsValue(v) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func tableOp_contains_2(t *jTableCore, r, c any) bool {
        \tif c == nil || javalang.IsNil(c) {
        \t\treturn false
        \t}
        \tidx := t.findRowIdx(r)
        \tif idx < 0 {
        \t\treturn false
        \t}
        \te := &t.entries[idx]
        \tif e.hasSingle {
        \t\treturn javalang.ObjectsEquals(e.col0, c)
        \t}
        \treturn e.multi != nil && e.multi.ContainsKey(c)
        }

        func tableOp_get_2(t *jTableCore, r, c any) any {
        \tif c == nil || javalang.IsNil(c) {
        \t\treturn nil
        \t}
        \tidx := t.findRowIdx(r)
        \tif idx < 0 {
        \t\treturn nil
        \t}
        \te := &t.entries[idx]
        \tif e.hasSingle {
        \t\tif javalang.ObjectsEquals(e.col0, c) {
        \t\t\treturn e.val0
        \t\t}
        \t\treturn nil
        \t}
        \tif e.multi == nil {
        \t\treturn nil
        \t}
        \treturn e.multi.Get(c)
        }

        func tableOp_put_3(t *jTableCore, r, c, v any) any {
        \tr = CheckNotNull(ToAny(r))
        \tc = CheckNotNull(ToAny(c))
        \tv = CheckNotNull(ToAny(v))
        \tidx := t.findRowIdx(r)
        \tif idx < 0 {
        \t\tt.insertNewRow(r, c, v)
        \t\treturn nil
        \t}
        \te := &t.entries[idx]
        \tif e.hasSingle {
        \t\tif javalang.ObjectsEquals(e.col0, c) {
        \t\t\told := e.val0
        \t\t\te.val0 = v
        \t\t\treturn old
        \t\t}
        \t\tmc := t.newMultiRow()
        \t\tmc.Put(e.col0, e.val0)
        \t\told := mc.Put(c, v)
        \t\te.col0 = nil
        \t\te.val0 = nil
        \t\te.hasSingle = false
        \t\te.multi = mc
        \t\treturn old
        \t}
        \tif e.multi == nil {
        \t\tif t.rowFactory == nil {
        \t\t\te.col0 = c
        \t\t\te.val0 = v
        \t\t\te.hasSingle = true
        \t\t\treturn nil
        \t\t}
        \t\te.multi = t.newMultiRow()
        \t}
        \treturn e.multi.Put(c, v)
        }

        func tableOp_remove_2(t *jTableCore, r, c any) any {
        \tif c == nil || javalang.IsNil(c) {
        \t\treturn nil
        \t}
        \tidx := t.findRowIdx(r)
        \tif idx < 0 {
        \t\treturn nil
        \t}
        \te := &t.entries[idx]
        \tif e.hasSingle {
        \t\tif !javalang.ObjectsEquals(e.col0, c) {
        \t\t\treturn nil
        \t\t}
        \t\told := e.val0
        \t\tt.deleteRowAt(idx)
        \t\treturn old
        \t}
        \tif e.multi == nil {
        \t\treturn nil
        \t}
        \told := e.multi.Remove(c)
        \tif e.multi.IsEmpty() {
        \t\tt.deleteRowAt(idx)
        \t}
        \treturn old
        }

        func tableOp_clear_0(t *jTableCore) {
        \tif len(t.entries) == 0 {
        \t\treturn
        \t}
        \tclear(t.entries)
        \tt.entries = t.entries[:0]
        \tif len(t.buckets) > 0 {
        \t\tclear(t.buckets)
        \t}
        \tif t.strIdx != nil {
        \t\tclear(t.strIdx)
        \t}
        \tt.liveRows = 0
        \tt.strLive = 0
        \tt.nonStrLive = 0
        \tt.modVer++
        }

        func tableOp_rowKeySet_0(t *jTableCore) Ext_java_util_Set {
        \ts := &jTableRowKeySet{t: t}
        \ts.BindExtSelf_(s)
        \treturn s
        }

        func tableOp_columnKeySet_0(t *jTableCore) Ext_java_util_Set {
        \tvar sc JSetCore
        \tfor i := range t.entries {
        \t\te := &t.entries[i]
        \t\tif e.deleted {
        \t\t\tcontinue
        \t\t}
        \t\tif e.hasSingle {
        \t\t\tsc.Add(e.col0)
        \t\t} else if e.multi != nil {
        \t\t\te.multi.ForEachEntry(func(ck, _ any) bool {
        \t\t\t\tsc.Add(ck)
        \t\t\t\treturn true
        \t\t\t})
        \t\t}
        \t}
        \tobj := &Ext_java_util_LinkedHashSet{Ext_java_util_HashSetBase: Ext_java_util_HashSetBase{inlineSet: sc}}
        \tobj.BindExtSelf_(obj)
        \treturn obj
        }

        func tableOp_values_0(t *jTableCore) Ext_java_util_Collection {
        \tvar vals []any
        \tfor i := range t.entries {
        \t\te := &t.entries[i]
        \t\tif e.deleted {
        \t\t\tcontinue
        \t\t}
        \t\tif e.hasSingle {
        \t\t\tvals = append(vals, e.val0)
        \t\t} else if e.multi != nil {
        \t\t\te.multi.ForEachEntry(func(_, cv any) bool {
        \t\t\t\tvals = append(vals, cv)
        \t\t\t\treturn true
        \t\t\t})
        \t\t}
        \t}
        \treturn wrapArrayListOwn(vals)
        }

        func tableOp_row_1(t *jTableCore, r any) Ext_java_util_Map {
        \tr = CheckNotNull(ToAny(r))
        \tif rc, ok := r.(jRawColorIdProvider); ok {
        \t\tif raw, valid := rc.RawColorId_(); valid {
        \t\t\treturn t.rowColorId(raw)
        \t\t}
        \t}
        \tidx := t.findRowIdx(r)
        \treturn &jTableRowMap{t: t, r: r, hintIdx: int32(idx), hintVer: t.modVer}
        }

        func (t *jTableCore) containsRowColorId(raw int64) bool {
        \tidx := t.findRowIdxColorId(raw)
        \treturn idx >= 0 && !t.entries[idx].isEmptyRow()
        }

        func (t *jTableCore) containsColorId(raw int64, c any) bool {
        \tif c == nil || javalang.IsNil(c) {
        \t\treturn false
        \t}
        \tidx := t.findRowIdxColorId(raw)
        \tif idx < 0 {
        \t\treturn false
        \t}
        \te := &t.entries[idx]
        \tif e.hasSingle {
        \t\treturn javalang.ObjectsEquals(e.col0, c)
        \t}
        \treturn e.multi != nil && e.multi.ContainsKey(c)
        }

        func (t *jTableCore) getColorId(raw int64, c any) any {
        \tif c == nil || javalang.IsNil(c) {
        \t\treturn nil
        \t}
        \tidx := t.findRowIdxColorId(raw)
        \tif idx < 0 {
        \t\treturn nil
        \t}
        \te := &t.entries[idx]
        \tif e.hasSingle {
        \t\tif javalang.ObjectsEquals(e.col0, c) {
        \t\t\treturn e.val0
        \t\t}
        \t\treturn nil
        \t}
        \tif e.multi == nil {
        \t\treturn nil
        \t}
        \treturn e.multi.Get(c)
        }

        func (t *jTableCore) putColorId(raw int64, c, v any) any {
        \tc = CheckNotNull(ToAny(c))
        \tv = CheckNotNull(ToAny(v))
        \tidx := t.findRowIdxColorId(raw)
        \tif idx < 0 {
        \t\tt.insertNewRowColorId(raw, c, v)
        \t\treturn nil
        \t}
        \te := &t.entries[idx]
        \tif e.hasSingle {
        \t\tif javalang.ObjectsEquals(e.col0, c) {
        \t\t\told := e.val0
        \t\t\te.val0 = v
        \t\t\treturn old
        \t\t}
        \t\tmc := t.newMultiRow()
        \t\tmc.Put(e.col0, e.val0)
        \t\told := mc.Put(c, v)
        \t\te.col0 = nil
        \t\te.val0 = nil
        \t\te.hasSingle = false
        \t\te.multi = mc
        \t\treturn old
        \t}
        \tif e.multi == nil {
        \t\tif t.rowFactory == nil {
        \t\t\te.col0 = c
        \t\t\te.val0 = v
        \t\t\te.hasSingle = true
        \t\t\treturn nil
        \t\t}
        \t\te.multi = t.newMultiRow()
        \t}
        \treturn e.multi.Put(c, v)
        }

        func (t *jTableCore) removeColorId(raw int64, c any) any {
        \tif c == nil || javalang.IsNil(c) {
        \t\treturn nil
        \t}
        \tidx := t.findRowIdxColorId(raw)
        \tif idx < 0 {
        \t\treturn nil
        \t}
        \te := &t.entries[idx]
        \tif e.hasSingle {
        \t\tif !javalang.ObjectsEquals(e.col0, c) {
        \t\t\treturn nil
        \t\t}
        \t\told := e.val0
        \t\tt.deleteRowAt(idx)
        \t\treturn old
        \t}
        \tif e.multi == nil {
        \t\treturn nil
        \t}
        \told := e.multi.Remove(c)
        \tif e.multi.IsEmpty() {
        \t\tt.deleteRowAt(idx)
        \t}
        \treturn old
        }

        func (t *jTableCore) rowColorId(raw int64) Ext_java_util_Map {
        \tidx := t.findRowIdxColorId(raw)
        \treturn &jTableRowMap{t: t, rI64: raw, isColorId: true, hintIdx: int32(idx), hintVer: t.modVer}
        }

        func tableOp_rowMap_0(t *jTableCore) Ext_java_util_Map {
        \tm := &jTableRowMapOuter{t: t}
        \tm.BindExtSelf_(m)
        \treturn m
        }

        type jTableRowKeySet struct {
        \tExt_java_util_AbstractSetBase
        \tt *jTableCore
        }

        func (s *jTableRowKeySet) nextLiveIdx(from int) int {
        \tentries := s.t.entries
        \tfor i := from; i < len(entries); i++ {
        \t\tif !entries[i].isEmptyRow() {
        \t\t\treturn i
        \t\t}
        \t}
        \treturn -1
        }

        func (s *jTableRowKeySet) iterInit_() (int32, bool) {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \ti := s.nextLiveIdx(0)
        \tif i < 0 {
        \t\treturn -1, false
        \t}
        \treturn int32(i), true
        }

        func (s *jTableRowKeySet) iterNext_(pos int32) (any, int32) {
        \tr := s.t.entries[pos].getRowKey()
        \tnext := s.nextLiveIdx(int(pos) + 1)
        \treturn r, int32(next)
        }

        func (s *jTableRowKeySet) Size() int32 {
        \tvar n int32
        \tfor i := range s.t.entries {
        \t\tif !s.t.entries[i].isEmptyRow() {
        \t\t\tn++
        \t\t}
        \t}
        \treturn n
        }

        func (s *jTableRowKeySet) IsEmpty() bool       { return tableOp_isEmpty_0(s.t) }
        func (s *jTableRowKeySet) Contains(r any) bool { return tableOp_containsRow_1(s.t, r) }
        func (s *jTableRowKeySet) Clear()              { tableOp_clear_0(s.t) }
        func (s *jTableRowKeySet) Remove_1_any(r any) bool {
        \tidx := s.t.findRowIdx(r)
        \tif idx < 0 || s.t.entries[idx].isEmptyRow() {
        \t\treturn false
        \t}
        \ts.t.deleteRowAt(idx)
        \treturn true
        }
        func (s *jTableRowKeySet) Iterator() any {
        \tvar keys []any
        \tfor i := range s.t.entries {
        \t\tif !s.t.entries[i].isEmptyRow() {
        \t\t\tkeys = append(keys, s.t.entries[i].getRowKey())
        \t\t}
        \t}
        \treturn &jTableRowKeyIter{t: s.t, keys: keys}
        }

        type jTableRowKeyIter struct {
        \tt       *jTableCore
        \tkeys    []any
        \tidx     int
        \tlastKey any
        \thasLast bool
        }

        func (*jTableRowKeyIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jTableRowKeyIter) HasNext() bool { return it.idx < len(it.keys) }
        func (it *jTableRowKeyIter) Next() any {
        \tif it.idx >= len(it.keys) {
        \t\tpanic(javalang.NewNoSuchElementException("exhausted"))
        \t}
        \tk := it.keys[it.idx]
        \tit.idx++
        \tit.lastKey = k
        \tit.hasLast = true
        \treturn k
        }
        func (it *jTableRowKeyIter) Remove() {
        \tif !it.hasLast {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \tif idx := it.t.findRowIdx(it.lastKey); idx >= 0 {
        \t\tit.t.deleteRowAt(idx)
        \t}
        \tit.hasLast = false
        }

        type jTableRowMapOuter struct {
        \tExt_java_util_AbstractMapBase
        \tt *jTableCore
        }

        func (m *jTableRowMapOuter) Size() int32 {
        \tvar n int32
        \tfor i := range m.t.entries {
        \t\tif !m.t.entries[i].isEmptyRow() {
        \t\t\tn++
        \t\t}
        \t}
        \treturn n
        }
        func (m *jTableRowMapOuter) ContainsKey(r any) bool { return tableOp_containsRow_1(m.t, r) }
        func (m *jTableRowMapOuter) Get(r any) any {
        \tif !tableOp_containsRow_1(m.t, r) {
        \t\treturn nil
        \t}
        \treturn tableOp_row_1(m.t, r)
        }
        func (m *jTableRowMapOuter) Remove_1(r any) any {
        \tidx := m.t.findRowIdx(r)
        \tif idx < 0 || m.t.entries[idx].isEmptyRow() {
        \t\treturn nil
        \t}
        \te := &m.t.entries[idx]
        \tcopyCore := NewJMapCore(false)
        \tif e.hasSingle {
        \t\tcopyCore.Put(e.col0, e.val0)
        \t} else if e.multi != nil {
        \t\te.multi.ForEachEntry(func(ck, cv any) bool {
        \t\t\tcopyCore.Put(ck, cv)
        \t\t\treturn true
        \t\t})
        \t}
        \tm.t.deleteRowAt(idx)
        \tlhm := &Ext_java_util_LinkedHashMap{Ext_java_util_HashMapBase: Ext_java_util_HashMapBase{inlineCore: *copyCore}}
        \tlhm.BindExtSelf_(lhm)
        \treturn lhm
        }
        func (m *jTableRowMapOuter) Clear() { tableOp_clear_0(m.t) }
        func (m *jTableRowMapOuter) EntrySet() any {
        \tvar ents []any
        \tfor i := range m.t.entries {
        \t\te := &m.t.entries[i]
        \t\tif !e.isEmptyRow() {
        \t\t\tr := e.getRowKey()
        \t\t\tents = append(ents, &jSimpleEntry{k: r, v: tableOp_row_1(m.t, r)})
        \t\t}
        \t}
        \treturn wrapLinkedHashSet(ents)
        }

        type jTableRowMap struct {
        \tt         *jTableCore
        \tr         any
        \trI64      int64
        \thintIdx   int32
        \thintVer   uint32
        \tisColorId bool
        }

        func (m *jTableRowMap) resolveEntry() (*jTableRowEntry, int) {
        \tt := m.t
        \tif m.hintIdx >= 0 && m.hintVer == t.modVer && int(m.hintIdx) < len(t.entries) {
        \t\te := &t.entries[m.hintIdx]
        \t\tif !e.deleted {
        \t\t\treturn e, int(m.hintIdx)
        \t\t}
        \t\tm.hintIdx = -1
        \t}
        \tvar idx int
        \tif m.isColorId {
        \t\tidx = t.findRowIdxColorId(m.rI64)
        \t} else {
        \t\tidx = t.findRowIdx(m.r)
        \t}
        \tif idx < 0 {
        \t\treturn nil, -1
        \t}
        \tm.hintIdx = int32(idx)
        \tm.hintVer = t.modVer
        \treturn &t.entries[idx], idx
        }

        func (m *jTableRowMap) putImpl(c, v any) any {
        \tif m.isColorId {
        \t\treturn m.t.putColorId(m.rI64, c, v)
        \t}
        \treturn tableOp_put_3(m.t, m.r, c, v)
        }

        func (m *jTableRowMap) removeImpl(c any) any {
        \tif m.isColorId {
        \t\treturn m.t.removeColorId(m.rI64, c)
        \t}
        \treturn tableOp_remove_2(m.t, m.r, c)
        }

        func (m *jTableRowMap) sizeImpl() int32 {
        \te, _ := m.resolveEntry()
        \tif e == nil {
        \t\treturn 0
        \t}
        \tif e.hasSingle {
        \t\treturn 1
        \t}
        \tif e.multi != nil {
        \t\treturn e.multi.Size()
        \t}
        \treturn 0
        }

        func (m *jTableRowMap) isEmptyImpl() bool {
        \te, _ := m.resolveEntry()
        \treturn e.isEmptyRow()
        }

        func (m *jTableRowMap) containsKeyImpl(c any) bool {
        \tif c == nil || javalang.IsNil(c) {
        \t\treturn false
        \t}
        \te, _ := m.resolveEntry()
        \tif e == nil {
        \t\treturn false
        \t}
        \tif e.hasSingle {
        \t\treturn javalang.ObjectsEquals(e.col0, c)
        \t}
        \treturn e.multi != nil && e.multi.ContainsKey(c)
        }
        func (m *jTableRowMap) MapContainsStr_(c string) bool {
        \tif m == nil { ThrowNPE() }
        \tif c == NullString { return false }
        \te, _ := m.resolveEntry()
        \tif e == nil { return false }
        \tif e.hasSingle {
        \t\tcs, ok := e.col0.(string)
        \t\treturn ok && cs == c
        \t}
        \treturn e.multi != nil && e.multi.ContainsStr(c)
        }
        func (m *jTableRowMap) MapGetStr_(c string) any {
        \tif m == nil { ThrowNPE() }
        \tif c == NullString { return nil }
        \te, _ := m.resolveEntry()
        \tif e == nil { return nil }
        \tif e.hasSingle {
        \t\tif cs, ok := e.col0.(string); ok && cs == c { return e.val0 }
        \t\treturn nil
        \t}
        \tif e.multi != nil { return e.multi.GetStr(c) }
        \treturn nil
        }
        func (m *jTableRowMap) MapGetOrDefaultStr_(c string, def any) any {
        \tif v := m.MapGetStr_(c); v != nil { return v }
        \treturn ToAny(def)
        }
        func (m *jTableRowMap) MapPutStr_(c string, v any) any {
        \tif m == nil { ThrowNPE() }
        \treturn m.putImpl(BoxString(c), v)
        }
        func (m *jTableRowMap) MapRemoveStr_(c string) any {
        \tif m == nil { ThrowNPE() }
        \tif c == NullString { return nil }
        \tif !m.MapContainsStr_(c) { return nil }
        \treturn m.removeImpl(BoxString(c))
        }

        func (m *jTableRowMap) containsValueImpl(v any) bool {
        \tif v == nil || javalang.IsNil(v) {
        \t\treturn false
        \t}
        \te, _ := m.resolveEntry()
        \tif e == nil {
        \t\treturn false
        \t}
        \tif e.hasSingle {
        \t\treturn javalang.ObjectsEquals(e.val0, v)
        \t}
        \treturn e.multi != nil && e.multi.ContainsValue(v)
        }

        func (m *jTableRowMap) getImpl(c any) any {
        \tif c == nil || javalang.IsNil(c) {
        \t\treturn nil
        \t}
        \te, _ := m.resolveEntry()
        \tif e == nil {
        \t\treturn nil
        \t}
        \tif e.hasSingle {
        \t\tif javalang.ObjectsEquals(e.col0, c) {
        \t\t\treturn e.val0
        \t\t}
        \t\treturn nil
        \t}
        \tif e.multi == nil {
        \t\treturn nil
        \t}
        \treturn e.multi.Get(c)
        }

        func (m *jTableRowMap) getOrDefaultImpl(c, def any) any {
        \tv := m.getImpl(c)
        \tif v != nil {
        \t\treturn v
        \t}
        \treturn ToAny(def)
        }

        func (m *jTableRowMap) clearImpl() {
        \t_, idx := m.resolveEntry()
        \tif idx >= 0 {
        \t\tm.t.deleteRowAt(idx)
        \t\tm.hintIdx = -1
        \t}
        }

        func (m *jTableRowMap) forEachImpl(action Ext_java_util_function_BiConsumer) {
        \tCheckNotNull(action)
        \te, _ := m.resolveEntry()
        \tif e == nil {
        \t\treturn
        \t}
        \tif e.hasSingle {
        \t\taction.Accept(e.col0, e.val0)
        \t\treturn
        \t}
        \tif e.multi != nil {
        \t\tmapCoreForEach(e.multi, action)
        \t}
        }

        func (m *jTableRowMap) keySetImpl() Ext_java_util_Set {
        \ts := &jTableRowKeyView{m: m}
        \ts.BindExtSelf_(s)
        \treturn s
        }

        func (m *jTableRowMap) entrySetImpl() Ext_java_util_Set {
        \ts := &jTableRowEntrySet{m: m}
        \ts.BindExtSelf_(s)
        \treturn s
        }

        func (m *jTableRowMap) equalsImpl(other any) bool {
        \tif other == m {
        \t\treturn true
        \t}
        \tom, ok := other.(Ext_java_util_Map)
        \tif !ok || om == nil || javalang.IsNil(om) {
        \t\treturn false
        \t}
        \te, _ := m.resolveEntry()
        \tif e.isEmptyRow() {
        \t\treturn om.IsEmpty()
        \t}
        \tif e.hasSingle {
        \t\treturn om.Size() == 1 && om.ContainsKey(e.col0) && javalang.ObjectsEquals(e.val0, om.Get(e.col0))
        \t}
        \treturn mapCoreEquals(e.multi, other)
        }

        func (m *jTableRowMap) hashCodeImpl() int32 {
        \te, _ := m.resolveEntry()
        \tif e.isEmptyRow() {
        \t\treturn 0
        \t}
        \tif e.hasSingle {
        \t\treturn elemJavaHashCode(e.col0, false) ^ elemJavaHashCode(e.val0, false)
        \t}
        \treturn mapCoreHash(e.multi)
        }

        func (m *jTableRowMap) stringImpl() string {
        \te, _ := m.resolveEntry()
        \tif e.isEmptyRow() {
        \t\treturn "{}"
        \t}
        \tif e.hasSingle {
        \t\treturn "{" + StrVal(e.col0) + "=" + StrVal(e.val0) + "}"
        \t}
        \treturn mapCoreString(e.multi)
        }

        type jTableRowKeyView struct {
        \tExt_java_util_AbstractSetBase
        \tm *jTableRowMap
        }

        func (s *jTableRowKeyView) Size() int32             { return s.m.sizeImpl() }
        func (s *jTableRowKeyView) IsEmpty() bool           { return s.m.isEmptyImpl() }
        func (s *jTableRowKeyView) Contains(c any) bool     { return s.m.containsKeyImpl(c) }
        func (s *jTableRowKeyView) Clear()                  { s.m.clearImpl() }
        func (s *jTableRowKeyView) Remove_1_any(c any) bool { return s.m.removeImpl(c) != nil }
        func (s *jTableRowKeyView) Iterator() any {
        \te, _ := s.m.resolveEntry()
        \tif e.isEmptyRow() {
        \t\treturn emptyMapIter_
        \t}
        \tif e.hasSingle {
        \t\treturn &jTableRowSingleIter{m: s.m, col: e.col0, val: e.val0, mode: 0}
        \t}
        \treturn newMapKeyIter(e.multi)
        }

        type jTableRowEntrySet struct {
        \tExt_java_util_AbstractSetBase
        \tm *jTableRowMap
        }

        func (s *jTableRowEntrySet) Size() int32   { return s.m.sizeImpl() }
        func (s *jTableRowEntrySet) IsEmpty() bool { return s.m.isEmptyImpl() }
        func (s *jTableRowEntrySet) Clear()        { s.m.clearImpl() }
        func (s *jTableRowEntrySet) Contains(o any) bool {
        \tent, ok := o.(Ext_java_util_Map_Entry)
        \tif !ok || ent == nil || javalang.IsNil(ent) {
        \t\treturn false
        \t}
        \tk := ent.GetKey()
        \tv := s.m.getImpl(k)
        \treturn v != nil && javalang.ObjectsEquals(v, ToAny(ent.GetValue()))
        }
        func (s *jTableRowEntrySet) Remove_1_any(o any) bool {
        \tent, ok := o.(Ext_java_util_Map_Entry)
        \tif !ok || ent == nil || javalang.IsNil(ent) {
        \t\treturn false
        \t}
        \tk := ent.GetKey()
        \tv := s.m.getImpl(k)
        \tif v == nil || !javalang.ObjectsEquals(v, ToAny(ent.GetValue())) {
        \t\treturn false
        \t}
        \ts.m.removeImpl(k)
        \treturn true
        }
        func (s *jTableRowEntrySet) Iterator() any {
        \te, _ := s.m.resolveEntry()
        \tif e.isEmptyRow() {
        \t\treturn emptyMapIter_
        \t}
        \tif e.hasSingle {
        \t\treturn &jTableRowSingleIter{m: s.m, col: e.col0, val: e.val0, mode: 2}
        \t}
        \tc, ok := initMapCursor(e.multi, true)
        \tif !ok {
        \t\treturn emptyMapIter_
        \t}
        \treturn &jMapEntryIter{cur: c}
        }

        type jTableRowSingleIter struct {
        \tm       *jTableRowMap
        \tcol     any
        \tval     any
        \tmode    uint8 // 0=key, 1=val, 2=entry
        \tdone    bool
        \tremoved bool
        }

        func (*jTableRowSingleIter) Implements_j2gort_Ext_java_util_Iterator()  {}
        func (*jTableRowSingleIter) Implements_j2gort_Ext_java_util_Map_Entry() {}
        func (it *jTableRowSingleIter) GetKey() any                             { return it.col }
        func (it *jTableRowSingleIter) GetValue() any                           { return it.val }
        func (it *jTableRowSingleIter) SetValue(nv any) any {
        \tnv = CheckNotNull(ToAny(nv))
        \told := it.m.putImpl(it.col, nv)
        \tit.val = nv
        \treturn old
        }
        func (it *jTableRowSingleIter) Equals(o any) bool { return mapEntryEquals(it, o) }
        func (it *jTableRowSingleIter) HashCode() int32   { return mapEntryHash(it) }
        func (it *jTableRowSingleIter) String() string    { return StrVal(it.col) + "=" + StrVal(it.val) }
        func (it *jTableRowSingleIter) HasNext() bool     { return !it.done }
        func (it *jTableRowSingleIter) Next() any {
        \tif it.done {
        \t\tpanic(javalang.NewNoSuchElementException("exhausted"))
        \t}
        \tit.done = true
        \tswitch it.mode {
        \tcase 0:
        \t\treturn it.col
        \tcase 1:
        \t\treturn it.val
        \tdefault:
        \t\treturn it
        \t}
        }
        func (it *jTableRowSingleIter) Remove() {
        \tif !it.done || it.removed {
        \t\tpanic(javalang.NewIllegalStateException("remove"))
        \t}
        \tit.removed = true
        \tit.m.removeImpl(it.col)
        }

        type jTableRowValues struct {
        \tm *jTableRowMap
        }

        func (v jTableRowValues) IsJavaNil_() bool { return v.m == nil }
        func (v jTableRowValues) typedViewLen_() (int, bool) {
        \tif v.m == nil {
        \t\tThrowNPE()
        \t}
        \te, _ := v.m.resolveEntry()
        \tif e == nil || e.deleted {
        \t\treturn 0, true
        \t}
        \tif e.hasSingle {
        \t\treturn 1, true
        \t}
        \tif e.multi == nil || e.multi.IsEmpty() {
        \t\treturn 0, true
        \t}
        \treturn 0, false
        }
        func (v jTableRowValues) typedViewAt_(int) any {
        \te, _ := v.m.resolveEntry()
        \tif e == nil || !e.hasSingle {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \treturn e.val0
        }
        func (jTableRowValues) typedViewNodeAt_(int) (uint32, bool) { return 0, false }
        func (v jTableRowValues) iteratorImpl() any {
        \te, _ := v.m.resolveEntry()
        \tif e.isEmptyRow() {
        \t\treturn emptyMapIter_
        \t}
        \tif e.hasSingle {
        \t\treturn &jTableRowSingleIter{m: v.m, col: e.col0, val: e.val0, mode: 1}
        \t}
        \tc, ok := initMapCursor(e.multi, true)
        \tif !ok {
        \t\treturn emptyMapIter_
        \t}
        \treturn &jMapValIter{cur: c}
        }
        func (v jTableRowValues) toArrayImpl() []any {
        \te, _ := v.m.resolveEntry()
        \tif e.isEmptyRow() {
        \t\treturn []any{}
        \t}
        \tif e.hasSingle {
        \t\treturn []any{e.val0}
        \t}
        \tvar out []any
        \tif e.multi != nil {
        \t\tout = make([]any, 0, e.multi.Size())
        \t\te.multi.ForEachEntry(func(_, cv any) bool {
        \t\t\tout = append(out, cv)
        \t\t\treturn true
        \t\t})
        \t}
        \treturn out
        }
        func (v jTableRowValues) forEachImpl(action Ext_java_util_function_Consumer) {
        \tCheckNotNull(action)
        \te, _ := v.m.resolveEntry()
        \tif e.isEmptyRow() {
        \t\treturn
        \t}
        \tif e.hasSingle {
        \t\taction.Accept(e.val0)
        \t\treturn
        \t}
        \tif e.multi != nil {
        \t\te.multi.ForEachEntry(func(_, cv any) bool {
        \t\t\taction.Accept(cv)
        \t\t\treturn true
        \t\t})
        \t}
        }

        // --- PriorityQueue ---

        type jPriorityQueue struct {
        \titems []any
        \tcmp   func(a, b any) int32
        }

        func newPriorityQueue(cmp func(a, b any) int32) *jPriorityQueue {
        \tif cmp == nil { cmp = CompareValues }
        \treturn &jPriorityQueue{cmp: cmp}
        }

        func (pq *jPriorityQueue) offer(v any) bool {
        \tv = CheckNotNull(ToAny(v))
        \tpq.items = append(pq.items, v)
        \tpq.siftUp(len(pq.items) - 1)
        \treturn true
        }

        func (pq *jPriorityQueue) siftUp(k int) {
        \tx := pq.items[k]
        \tfor k > 0 {
        \t\tp := (k - 1) >> 1
        \t\tif pq.cmp(x, pq.items[p]) >= 0 { break }
        \t\tpq.items[k] = pq.items[p]
        \t\tk = p
        \t}
        \tpq.items[k] = x
        }

        func (pq *jPriorityQueue) siftDown(k int) {
        \tn := len(pq.items)
        \tx := pq.items[k]
        \thalf := n >> 1
        \tfor k < half {
        \t\tc := (k << 1) + 1
        \t\tright := c + 1
        \t\tif right < n && pq.cmp(pq.items[c], pq.items[right]) > 0 {
        \t\t\tc = right
        \t\t}
        \t\tif pq.cmp(x, pq.items[c]) <= 0 { break }
        \t\tpq.items[k] = pq.items[c]
        \t\tk = c
        \t}
        \tpq.items[k] = x
        }

        func (pq *jPriorityQueue) poll() any {
        \tn := len(pq.items)
        \tif n == 0 { return nil }
        \tres := pq.items[0]
        \tlast := pq.items[n-1]
        \tpq.items[n-1] = nil
        \tpq.items = pq.items[:n-1]
        \tif n > 1 {
        \t\tpq.items[0] = last
        \t\tpq.siftDown(0)
        \t}
        \treturn res
        }

        func pqOp_size_0(pq *jPriorityQueue) int32         { return int32(len(pq.items)) }
        func pqOp_isEmpty_0(pq *jPriorityQueue) bool       { return len(pq.items) == 0 }
        func pqOp_add_1(pq *jPriorityQueue, v any) bool    { return pq.offer(v) }
        func pqOp_peek_0(pq *jPriorityQueue) any {
        \tif len(pq.items) == 0 { return nil }
        \treturn pq.items[0]
        }
        func pqOp_poll_0(pq *jPriorityQueue) any { return pq.poll() }
        func pqOp_remove_0(pq *jPriorityQueue) any {
        \tif len(pq.items) == 0 { panic(javalang.NewNoSuchElementException("empty")) }
        \treturn pq.poll()
        }
        func pqOp_element_0(pq *jPriorityQueue) any {
        \tif len(pq.items) == 0 { panic(javalang.NewNoSuchElementException("empty")) }
        \treturn pq.items[0]
        }
        func pqOp_retainAll_1(pq *jPriorityQueue, c Ext_java_util_Collection) bool {
        \tCheckNotNull(c)
        \tchanged := false
        \tfor i := len(pq.items) - 1; i >= 0; i-- {
        \t\tif !c.Contains(pq.items[i]) {
        \t\t\tpq.removeAt(i)
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }
        func pqOp_removeIf_1(pq *jPriorityQueue, pred Ext_java_util_function_Predicate) bool {
        \tCheckNotNull(pred)
        \tchanged := false
        \tfor i := len(pq.items) - 1; i >= 0; i-- {
        \t\tif pred.Test(pq.items[i]) {
        \t\t\tpq.removeAt(i)
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }
        func pqOp_remove_1(pq *jPriorityQueue, v any) bool {
        \tv = ToAny(v)
        \tfor i, x := range pq.items {
        \t\tif javalang.ObjectsEquals(x, v) {
        \t\t\tpq.removeAt(i)
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }
        """;
  }

  private static String buildColRtPart3() {
    return """
        func (pq *jPriorityQueue) removeAt(i int) {
        \tn := len(pq.items) - 1
        \tif i == n {
        \t\tpq.items[n] = nil
        \t\tpq.items = pq.items[:n]
        \t\treturn
        \t}
        \tmoved := pq.items[n]
        \tpq.items[n] = nil
        \tpq.items = pq.items[:n]
        \tpq.items[i] = moved
        \tpq.siftDown(i)
        \tif pq.items[i] == moved {
        \t\tpq.siftUp(i)
        \t}
        }
        func pqOp_contains_1(pq *jPriorityQueue, v any) bool {
        \tv = ToAny(v)
        \tfor _, x := range pq.items {
        \t\tif javalang.ObjectsEquals(x, v) { return true }
        \t}
        \treturn false
        }
        func pqOp_clear_0(pq *jPriorityQueue)                      { pq.items = nil }
        func pqOp_toArray_0(pq *jPriorityQueue) []any              { return CloneSlice(pq.items) }
        func pqOp_toArray_1(pq *jPriorityQueue, d []any) []any     { return sliceToArray1(pq.items, d) }
        func pqOp_iterator_0(pq *jPriorityQueue) any               { return newSliceObjIter(CloneSlice(pq.items)) }
        func pqOp_addAll_1(pq *jPriorityQueue, c Ext_java_util_Collection) bool {
        \tCheckNotNull(c)
        \tadded := false
        \tfor it := Iterate(c); it.HasNext(); { pq.offer(it.Next()); added = true }
        \treturn added
        }
        func pqOp_containsAll_1(pq *jPriorityQueue, c Ext_java_util_Collection) bool {
        \treturn colContainsAll(func(x any) bool { return pqOp_contains_1(pq, x) }, c)
        }
        func pqOp_removeAll_1(pq *jPriorityQueue, c Ext_java_util_Collection) bool {
        \tCheckNotNull(c)
        \tvar kept []any
        \tchanged := false
        \tfor _, x := range pq.items {
        \t\tif c.Contains(x) { changed = true } else { kept = append(kept, x) }
        \t}
        \tif changed {
        \t\tpq.items = nil
        \t\tfor _, x := range kept { pq.offer(x) }
        \t}
        \treturn changed
        }
        func pqOp_forEach_1(pq *jPriorityQueue, fn Ext_java_util_function_Consumer) {
        \tCheckNotNull(fn)
        \tfor _, x := range CloneSlice(pq.items) { fn.Accept(x) }
        }
        func pqOp_stream_0(pq *jPriorityQueue) Ext_java_util_stream_Stream {
        \treturn streamFromSlice(CloneSlice(pq.items))
        }
        func pqOp_toString_0(pq *jPriorityQueue) string       { return sliceJavaString(pq.items) }
        func pqOp_equals_1(pq *jPriorityQueue, o any) bool    { return RefEquals(pq, o) }
        func pqOp_hashCode_0(pq *jPriorityQueue) int32        { return javalang.IdentityHashCode(pq) }

        // --- BitSet (inline [8]uint64 = 512 dense bits + [8]int32 = 8 sparse high bits inside Ext_java_util_BitSet) ---

        const (
        \tbsInlineWords = 8
        \tbsInlineBits  = bsInlineWords * 64
        \tbsSparseCap   = 8
        )

        //go:noinline
        func bitSetThrowNegIdx(bitIndex int, name string) {
        \tpanic(javalang.NewIndexOutOfBoundsException(name + " < 0: " + strconv.Itoa(bitIndex)))
        }

        //go:noinline
        func bitSetThrowNegSize(nbits int32) {
        \tpanic(javalang.NewNegativeArraySizeException("nbits < 0: " + strconv.Itoa(int(nbits))))
        }

        //go:noinline
        func bitSetThrowBadRange(fromIndex, toIndex int) {
        \tpanic(javalang.NewIndexOutOfBoundsException("fromIndex: " + strconv.Itoa(fromIndex) + " > toIndex: " + strconv.Itoa(toIndex)))
        }

        func bitSetCheckIndex(bitIndex int, name string) {
        \tif bitIndex < 0 {
        \t\tbitSetThrowNegIdx(bitIndex, name)
        \t}
        }

        func bitSetCheckRange(fromIndex, toIndex int) {
        \tif fromIndex < 0 {
        \t\tbitSetThrowNegIdx(fromIndex, "fromIndex")
        \t}
        \tif toIndex < 0 {
        \t\tbitSetThrowNegIdx(toIndex, "toIndex")
        \t}
        \tif fromIndex > toIndex {
        \t\tbitSetThrowBadRange(fromIndex, toIndex)
        \t}
        }

        func InitExt_java_util_BitSet_0(b *Ext_java_util_BitSet) {
        \t*b = Ext_java_util_BitSet{initW: 1}
        }

        func InitExt_java_util_BitSet_1(b *Ext_java_util_BitSet, nbits int32) {
        \tif nbits < 0 {
        \t\tbitSetThrowNegSize(nbits)
        \t}
        \t*b = Ext_java_util_BitSet{initW: int32((int64(nbits) + 63) >> 6)}
        }

        func bitSetInitCap(b *Ext_java_util_BitSet, nbits int32) {
        \tInitExt_java_util_BitSet_1(b, nbits)
        }

        func (b *Ext_java_util_BitSet) bsHiContains(idx int32) bool {
        \tn := int(b.hiLen)
        \tfor i := 0; i < n; i++ {
        \t\tv := b.hi[i]
        \t\tif v >= idx {
        \t\t\treturn v == idx
        \t\t}
        \t}
        \treturn false
        }

        func (b *Ext_java_util_BitSet) bsHiInsert(idx int32) bool {
        \tn := int(b.hiLen)
        \tpos := 0
        \tfor pos < n && b.hi[pos] < idx {
        \t\tpos++
        \t}
        \tif pos < n && b.hi[pos] == idx {
        \t\treturn true
        \t}
        \tif n >= bsSparseCap {
        \t\treturn false
        \t}
        \tcopy(b.hi[pos+1:n+1], b.hi[pos:n])
        \tb.hi[pos] = idx
        \tb.hiLen = int32(n + 1)
        \treturn true
        }

        func (b *Ext_java_util_BitSet) bsHiRemove(idx int32) bool {
        \tn := int(b.hiLen)
        \tfor i := 0; i < n; i++ {
        \t\tv := b.hi[i]
        \t\tif v == idx {
        \t\t\tcopy(b.hi[i:n-1], b.hi[i+1:n])
        \t\t\tb.hi[n-1] = 0
        \t\t\tb.hiLen = int32(n - 1)
        \t\t\treturn true
        \t\t}
        \t\tif v > idx {
        \t\t\treturn false
        \t\t}
        \t}
        \treturn false
        }

        func (b *Ext_java_util_BitSet) bsWordAt(i int) uint64 {
        \tif b.words == nil {
        \t\tif uint(i) < bsInlineWords {
        \t\t\treturn b.buf[i]
        \t\t}
        \t\treturn 0
        \t}
        \tif uint(i) < uint(len(b.words)) {
        \t\treturn b.words[i]
        \t}
        \treturn 0
        }

        //go:noinline
        func (b *Ext_java_util_BitSet) bsPromoteToHeap(minWords int) {
        \tif b.hiLen > 0 {
        \t\thw := int(b.hi[b.hiLen-1]>>6) + 1
        \t\tif hw > minWords {
        \t\t\tminWords = hw
        \t\t}
        \t}
        \tif minWords < bsInlineWords {
        \t\tminWords = bsInlineWords
        \t}
        \tnewCap := max(minWords, max(16, int(b.initW)))
        \tw := make([]uint64, minWords, newCap)
        \tcopy(w[:bsInlineWords], b.buf[:])
        \tb.buf = [bsInlineWords]uint64{}
        \tfor i := int32(0); i < b.hiLen; i++ {
        \t\thb := int(b.hi[i])
        \t\tb.hi[i] = 0
        \t\tw[hb>>6] |= uint64(1) << (uint(hb) & 63)
        \t}
        \tb.hiLen = 0
        \tb.words = w
        }

        func (b *Ext_java_util_BitSet) bsMaterialize() {
        \tif b.hiLen > 0 && b.words == nil {
        \t\tb.bsPromoteToHeap(int(b.hi[b.hiLen-1]>>6) + 1)
        \t}
        }

        //go:noinline
        func (b *Ext_java_util_BitSet) bsEnsure(wordIdx int) {
        \tif b.words == nil {
        \t\tb.bsPromoteToHeap(wordIdx + 1)
        \t\treturn
        \t}
        \tif wordIdx < len(b.words) {
        \t\treturn
        \t}
        \tif wordIdx < cap(b.words) {
        \t\tprev := len(b.words)
        \t\tb.words = b.words[:wordIdx+1]
        \t\tclear(b.words[prev:])
        \t\treturn
        \t}
        \tnewCap := max(wordIdx+1, max(2*cap(b.words), int(b.initW)))
        \tnewWords := make([]uint64, wordIdx+1, newCap)
        \tcopy(newWords, b.words)
        \tb.words = newWords
        }

        func (b *Ext_java_util_BitSet) bsUsedWords() int {
        \tif b.words == nil {
        \t\tfor n := bsInlineWords; n > 0; n-- {
        \t\t\tif b.buf[n-1] != 0 {
        \t\t\t\treturn n
        \t\t\t}
        \t\t}
        \t\treturn 0
        \t}
        \tn := len(b.words)
        \tfor n > 0 && b.words[n-1] == 0 {
        \t\tn--
        \t}
        \treturn n
        }

        func (b *Ext_java_util_BitSet) bsTrim() {
        \tif b.words != nil {
        \t\tb.words = b.words[:b.bsUsedWords()]
        \t}
        }

        func (b *Ext_java_util_BitSet) bsLength() int {
        \tif b.hiLen > 0 {
        \t\treturn int(b.hi[b.hiLen-1]) + 1
        \t}
        \tn := b.bsUsedWords()
        \tif n == 0 {
        \t\treturn 0
        \t}
        \tlast := b.bsWordAt(n - 1)
        \treturn 64*(n-1) + (64 - bits.LeadingZeros64(last))
        }

        func (b *Ext_java_util_BitSet) bsSize() int {
        \tw := max(b.bsUsedWords(), int(b.initW))
        \tif cap(b.words) > bsInlineWords {
        \t\tw = max(w, cap(b.words))
        \t}
        \tif b.hiLen > 0 {
        \t\tw = max(w, int(b.hi[b.hiLen-1]>>6)+1)
        \t}
        \treturn w * 64
        }

        func (b *Ext_java_util_BitSet) bsSet(bitIndex int) {
        \tif bitIndex < 0 {
        \t\tbitSetThrowNegIdx(bitIndex, "bitIndex")
        \t}
        \tw := bitIndex >> 6
        \tmask := uint64(1) << (uint(bitIndex) & 63)
        \tif b.words == nil {
        \t\tif w < bsInlineWords {
        \t\t\tb.buf[w] |= mask
        \t\t\treturn
        \t\t}
        \t\tif b.bsHiInsert(int32(bitIndex)) {
        \t\t\treturn
        \t\t}
        \t\tb.bsPromoteToHeap(w + 1)
        \t\tb.words[w] |= mask
        \t\treturn
        \t}
        \tif w >= len(b.words) {
        \t\tb.bsEnsure(w)
        \t}
        \tb.words[w] |= mask
        }

        func (b *Ext_java_util_BitSet) bsClearBit(bitIndex int) {
        \tif bitIndex < 0 {
        \t\tbitSetThrowNegIdx(bitIndex, "bitIndex")
        \t}
        \tw := bitIndex >> 6
        \tif b.words == nil {
        \t\tif w < bsInlineWords {
        \t\t\tb.buf[w] &^= uint64(1) << (uint(bitIndex) & 63)
        \t\t} else if b.hiLen > 0 {
        \t\t\tb.bsHiRemove(int32(bitIndex))
        \t\t}
        \t\treturn
        \t}
        \tif w < len(b.words) {
        \t\tb.words[w] &^= uint64(1) << (uint(bitIndex) & 63)
        \t\tb.bsTrim()
        \t}
        }

        func (b *Ext_java_util_BitSet) bsSetRange(fromIndex, toIndex int) {
        \tbitSetCheckRange(fromIndex, toIndex)
        \tif fromIndex == toIndex {
        \t\treturn
        \t}
        \tif toIndex-1 == fromIndex {
        \t\tb.bsSet(fromIndex)
        \t\treturn
        \t}
        \tif b.words == nil && fromIndex >= bsInlineBits && toIndex-fromIndex <= bsSparseCap {
        \t\tfor i := fromIndex; i < toIndex; i++ {
        \t\t\tb.bsSet(i)
        \t\t}
        \t\treturn
        \t}
        \tstartWord := fromIndex >> 6
        \tendWord := (toIndex - 1) >> 6
        \tfirstWordMask := ^uint64(0) << (uint(fromIndex) & 63)
        \tlastWordMask := ^uint64(0) >> (63 - (uint(toIndex-1) & 63))
        \tif b.words == nil && endWord < bsInlineWords {
        \t\tif startWord == endWord {
        \t\t\tb.buf[startWord] |= firstWordMask & lastWordMask
        \t\t} else {
        \t\t\tb.buf[startWord] |= firstWordMask
        \t\t\tfor i := startWord + 1; i < endWord; i++ {
        \t\t\t\tb.buf[i] = ^uint64(0)
        \t\t\t}
        \t\t\tb.buf[endWord] |= lastWordMask
        \t\t}
        \t\treturn
        \t}
        \tb.bsEnsure(endWord)
        \tif startWord == endWord {
        \t\tb.words[startWord] |= firstWordMask & lastWordMask
        \t} else {
        \t\tb.words[startWord] |= firstWordMask
        \t\tfor i := startWord + 1; i < endWord; i++ {
        \t\t\tb.words[i] = ^uint64(0)
        \t\t}
        \t\tb.words[endWord] |= lastWordMask
        \t}
        }

        func (b *Ext_java_util_BitSet) bsClearRange(fromIndex, toIndex int) {
        \tbitSetCheckRange(fromIndex, toIndex)
        \tif fromIndex == toIndex {
        \t\treturn
        \t}
        \tstartWord := fromIndex >> 6
        \tendWord := (toIndex - 1) >> 6
        \tfirstWordMask := ^uint64(0) << (uint(fromIndex) & 63)
        \tlastWordMask := ^uint64(0) >> (63 - (uint(toIndex-1) & 63))
        \tif b.words == nil {
        \t\tif b.hiLen > 0 && toIndex > bsInlineBits {
        \t\t\twIdx := 0
        \t\t\tn := int(b.hiLen)
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\thb := b.hi[i]
        \t\t\t\tif int(hb) < fromIndex || int(hb) >= toIndex {
        \t\t\t\t\tb.hi[wIdx] = hb
        \t\t\t\t\twIdx++
        \t\t\t\t}
        \t\t\t}
        \t\t\tfor i := wIdx; i < n; i++ {
        \t\t\t\tb.hi[i] = 0
        \t\t\t}
        \t\t\tb.hiLen = int32(wIdx)
        \t\t}
        \t\tif startWord >= bsInlineWords {
        \t\t\treturn
        \t\t}
        \t\tif endWord >= bsInlineWords {
        \t\t\tendWord = bsInlineWords - 1
        \t\t\tlastWordMask = ^uint64(0)
        \t\t}
        \t\tif startWord == endWord {
        \t\t\tb.buf[startWord] &^= firstWordMask & lastWordMask
        \t\t} else {
        \t\t\tb.buf[startWord] &^= firstWordMask
        \t\t\tfor i := startWord + 1; i < endWord; i++ {
        \t\t\t\tb.buf[i] = 0
        \t\t\t}
        \t\t\tb.buf[endWord] &^= lastWordMask
        \t\t}
        \t\treturn
        \t}
        \tif startWord >= len(b.words) {
        \t\treturn
        \t}
        \tif endWord >= len(b.words) {
        \t\tendWord = len(b.words) - 1
        \t\tlastWordMask = ^uint64(0)
        \t}
        \tif startWord == endWord {
        \t\tb.words[startWord] &^= firstWordMask & lastWordMask
        \t} else {
        \t\tb.words[startWord] &^= firstWordMask
        \t\tfor i := startWord + 1; i < endWord; i++ {
        \t\t\tb.words[i] = 0
        \t\t}
        \t\tb.words[endWord] &^= lastWordMask
        \t}
        \tb.bsTrim()
        }

        func (b *Ext_java_util_BitSet) bsNextSetBit(fromIndex int) int {
        \tif fromIndex < 0 {
        \t\tbitSetThrowNegIdx(fromIndex, "fromIndex")
        \t}
        \tu := fromIndex >> 6
        \tif b.words == nil {
        \t\tif u < bsInlineWords {
        \t\t\tword := b.buf[u] & (^uint64(0) << (uint(fromIndex) & 63))
        \t\t\tfor {
        \t\t\t\tif word != 0 {
        \t\t\t\t\treturn u*64 + bits.TrailingZeros64(word)
        \t\t\t\t}
        \t\t\t\tu++
        \t\t\t\tif u >= bsInlineWords {
        \t\t\t\t\tbreak
        \t\t\t\t}
        \t\t\t\tword = b.buf[u]
        \t\t\t}
        \t\t}
        \t\tn := int(b.hiLen)
        \t\tfor i := 0; i < n; i++ {
        \t\t\tif hb := int(b.hi[i]); hb >= fromIndex {
        \t\t\t\treturn hb
        \t\t\t}
        \t\t}
        \t\treturn -1
        \t}
        \tif u >= len(b.words) {
        \t\treturn -1
        \t}
        \tword := b.words[u] & (^uint64(0) << (uint(fromIndex) & 63))
        \tfor {
        \t\tif word != 0 {
        \t\t\treturn u*64 + bits.TrailingZeros64(word)
        \t\t}
        \t\tu++
        \t\tif u >= len(b.words) {
        \t\t\treturn -1
        \t\t}
        \t\tword = b.words[u]
        \t}
        }

        func bitSetOp_ToByteArray(b *Ext_java_util_BitSet) []int8 {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tlength := b.bsLength()
        \tif length == 0 {
        \t\treturn []int8{}
        \t}
        \tnBytes := (length + 7) / 8
        \tout := make([]int8, nBytes)
        \tfor i := b.bsNextSetBit(0); i >= 0; i = b.bsNextSetBit(i + 1) {
        \t\tout[i/8] |= int8(uint8(1 << uint(i%8)))
        \t}
        \treturn out
        }
        func bitSetOp_ToLongArray(b *Ext_java_util_BitSet) []int64 {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tif b.hiLen > 0 {
        \t\thw := int(b.hi[b.hiLen-1]) >> 6
        \t\tout := make([]int64, hw+1)
        \t\tn := b.bsUsedWords()
        \t\tfor i := 0; i < n; i++ {
        \t\t\tout[i] = int64(b.bsWordAt(i))
        \t\t}
        \t\tfor i := int32(0); i < b.hiLen; i++ {
        \t\t\thb := int(b.hi[i])
        \t\t\tout[hb>>6] |= int64(uint64(1) << (uint(hb) & 63))
        \t\t}
        \t\treturn out
        \t}
        \tn := b.bsUsedWords()
        \tif n == 0 {
        \t\treturn []int64{}
        \t}
        \tout := make([]int64, n)
        \tfor i := 0; i < n; i++ {
        \t\tout[i] = int64(b.bsWordAt(i))
        \t}
        \treturn out
        }
        func bitSetOp_Flip_1(b *Ext_java_util_BitSet, bitIndex int32) {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \ti := int(bitIndex)
        \tif i < 0 {
        \t\tbitSetThrowNegIdx(i, "bitIndex")
        \t}
        \tw := i >> 6
        \tmask := uint64(1) << (uint(i) & 63)
        \tif b.words == nil {
        \t\tif w < bsInlineWords {
        \t\t\tb.buf[w] ^= mask
        \t\t\treturn
        \t\t}
        \t\tif b.bsHiRemove(bitIndex) {
        \t\t\treturn
        \t\t}
        \t\tif b.bsHiInsert(bitIndex) {
        \t\t\treturn
        \t\t}
        \t}
        \tb.bsEnsure(w)
        \tb.words[w] ^= mask
        \tb.bsTrim()
        }
        func bitSetOp_Flip_2(b *Ext_java_util_BitSet, from, to int32) {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tfromIndex, toIndex := int(from), int(to)
        \tbitSetCheckRange(fromIndex, toIndex)
        \tif fromIndex == toIndex {
        \t\treturn
        \t}
        \tif toIndex-1 == fromIndex {
        \t\tbitSetOp_Flip_1(b, from)
        \t\treturn
        \t}
        \tif b.words == nil && fromIndex >= bsInlineBits && toIndex-fromIndex <= bsSparseCap {
        \t\tfor idx := from; idx < to; idx++ {
        \t\t\tbitSetOp_Flip_1(b, idx)
        \t\t}
        \t\treturn
        \t}
        \tstartWord := fromIndex >> 6
        \tendWord := (toIndex - 1) >> 6
        \tfirstWordMask := ^uint64(0) << (uint(fromIndex) & 63)
        \tlastWordMask := ^uint64(0) >> (63 - (uint(toIndex-1) & 63))
        \tif b.words == nil && endWord < bsInlineWords {
        \t\tif startWord == endWord {
        \t\t\tb.buf[startWord] ^= firstWordMask & lastWordMask
        \t\t} else {
        \t\t\tb.buf[startWord] ^= firstWordMask
        \t\t\tfor i := startWord + 1; i < endWord; i++ {
        \t\t\t\tb.buf[i] = ^b.buf[i]
        \t\t\t}
        \t\t\tb.buf[endWord] ^= lastWordMask
        \t\t}
        \t\treturn
        \t}
        \tb.bsEnsure(endWord)
        \tif startWord == endWord {
        \t\tb.words[startWord] ^= firstWordMask & lastWordMask
        \t} else {
        \t\tb.words[startWord] ^= firstWordMask
        \t\tfor i := startWord + 1; i < endWord; i++ {
        \t\t\tb.words[i] = ^b.words[i]
        \t\t}
        \t\tb.words[endWord] ^= lastWordMask
        \t}
        \tb.bsTrim()
        }
        func bitSetOp_Set_1(b *Ext_java_util_BitSet, i int32) {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tb.bsSet(int(i))
        }
        func bitSetOp_Set_2_int32_bool(b *Ext_java_util_BitSet, i int32, v bool) {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tif v {
        \t\tb.bsSet(int(i))
        \t} else {
        \t\tb.bsClearBit(int(i))
        \t}
        }
        func bitSetOp_Set_2_int32_int32(b *Ext_java_util_BitSet, f, t int32) {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tb.bsSetRange(int(f), int(t))
        }
        func bitSetOp_Set_3(b *Ext_java_util_BitSet, f, t int32, v bool) {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tif v {
        \t\tb.bsSetRange(int(f), int(t))
        \t} else {
        \t\tb.bsClearRange(int(f), int(t))
        \t}
        }
        func bitSetOp_Clear(b *Ext_java_util_BitSet) {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tif b.words == nil {
        \t\tb.buf = [bsInlineWords]uint64{}
        \t} else {
        \t\tb.words = b.words[:0]
        \t}
        \tif b.hiLen > 0 {
        \t\tclear(b.hi[:b.hiLen])
        \t\tb.hiLen = 0
        \t}
        }
        func bitSetOp_Clear_0(b *Ext_java_util_BitSet) { bitSetOp_Clear(b) }
        func bitSetOp_Clear_1(b *Ext_java_util_BitSet, i int32) {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tb.bsClearBit(int(i))
        }
        func bitSetOp_Clear_2(b *Ext_java_util_BitSet, f, t int32) {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tb.bsClearRange(int(f), int(t))
        }
        func bitSetOp_Get_1(b *Ext_java_util_BitSet, i int32) bool {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tidx := int(i)
        \tif idx < 0 {
        \t\tbitSetThrowNegIdx(idx, "bitIndex")
        \t}
        \tw := idx >> 6
        \tif b.words == nil {
        \t\tif w < bsInlineWords {
        \t\t\treturn (b.buf[w] & (uint64(1) << (uint(idx) & 63))) != 0
        \t\t}
        \t\treturn b.hiLen > 0 && b.bsHiContains(i)
        \t}
        \treturn w < len(b.words) && (b.words[w]&(uint64(1)<<(uint(idx)&63))) != 0
        }
        func bitSetOp_Get(b *Ext_java_util_BitSet, i int32) bool { return bitSetOp_Get_1(b, i) }
        func bitSetOp_Get_2(b *Ext_java_util_BitSet, f, t int32) *Ext_java_util_BitSet {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tfromIndex, toIndex := int(f), int(t)
        \tbitSetCheckRange(fromIndex, toIndex)
        \tdst := &Ext_java_util_BitSet{initW: 1}
        \tdst.BindExtSelf_(dst)
        \tlength := b.bsLength()
        \tif length <= fromIndex || fromIndex == toIndex {
        \t\treturn dst
        \t}
        \tif toIndex > length {
        \t\ttoIndex = length
        \t}
        \tfor i := b.bsNextSetBit(fromIndex); i >= 0 && i < toIndex; i = b.bsNextSetBit(i + 1) {
        \t\tdst.bsSet(i - fromIndex)
        \t}
        \treturn dst
        }
        func bitSetOp_NextSetBit(b *Ext_java_util_BitSet, f int32) int32 {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \treturn int32(b.bsNextSetBit(int(f)))
        }
        func bitSetOp_NextClearBit(b *Ext_java_util_BitSet, f int32) int32 {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tfromIndex := int(f)
        \tif fromIndex < 0 {
        \t\tbitSetThrowNegIdx(fromIndex, "fromIndex")
        \t}
        \tu := fromIndex >> 6
        \tnLim := len(b.words)
        \tif b.words == nil {
        \t\tnLim = bsInlineWords
        \t}
        \tvar cand int
        \tif u >= nLim {
        \t\tcand = fromIndex
        \t} else {
        \t\tword := ^b.bsWordAt(u) & (^uint64(0) << (uint(fromIndex) & 63))
        \t\tfor {
        \t\t\tif word != 0 {
        \t\t\t\tcand = u*64 + bits.TrailingZeros64(word)
        \t\t\t\tbreak
        \t\t\t}
        \t\t\tu++
        \t\t\tif u >= nLim {
        \t\t\t\tcand = nLim * 64
        \t\t\t\tbreak
        \t\t\t}
        \t\t\tword = ^b.bsWordAt(u)
        \t\t}
        \t}
        \tif b.hiLen > 0 && cand >= bsInlineBits {
        \t\tn := int(b.hiLen)
        \t\tfor i := 0; i < n; i++ {
        \t\t\thb := int(b.hi[i])
        \t\t\tif hb == cand {
        \t\t\t\tcand++
        \t\t\t} else if hb > cand {
        \t\t\t\tbreak
        \t\t\t}
        \t\t}
        \t}
        \treturn int32(cand)
        }
        func bitSetOp_PreviousSetBit(b *Ext_java_util_BitSet, f int32) int32 {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tfromIndex := int(f)
        \tif fromIndex < 0 {
        \t\tif fromIndex == -1 {
        \t\t\treturn -1
        \t\t}
        \t\tbitSetThrowNegIdx(fromIndex, "fromIndex")
        \t}
        \tif b.hiLen > 0 {
        \t\tfor i := int(b.hiLen) - 1; i >= 0; i-- {
        \t\t\tif hb := int(b.hi[i]); hb <= fromIndex {
        \t\t\t\treturn int32(hb)
        \t\t\t}
        \t\t}
        \t}
        \tnUsed := b.bsUsedWords()
        \tif nUsed == 0 {
        \t\treturn -1
        \t}
        \tw := fromIndex >> 6
        \tif w >= nUsed {
        \t\tlast := b.bsWordAt(nUsed - 1)
        \t\treturn int32(64*(nUsed-1) + (64 - bits.LeadingZeros64(last)) - 1)
        \t}
        \tword := b.bsWordAt(w) & (^uint64(0) >> (63 - (uint(fromIndex) & 63)))
        \tfor {
        \t\tif word != 0 {
        \t\t\treturn int32((w+1)*64 - 1 - bits.LeadingZeros64(word))
        \t\t}
        \t\tif w == 0 {
        \t\t\treturn -1
        \t\t}
        \t\tw--
        \t\tword = b.bsWordAt(w)
        \t}
        }
        func bitSetOp_PreviousClearBit(b *Ext_java_util_BitSet, f int32) int32 {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tfromIndex := int(f)
        \tif fromIndex < 0 {
        \t\tif fromIndex == -1 {
        \t\t\treturn -1
        \t\t}
        \t\tbitSetThrowNegIdx(fromIndex, "fromIndex")
        \t}
        \tif b.hiLen > 0 && fromIndex >= bsInlineBits {
        \t\tfor i := int(b.hiLen) - 1; i >= 0; i-- {
        \t\t\thb := int(b.hi[i])
        \t\t\tif hb == fromIndex {
        \t\t\t\tfromIndex--
        \t\t\t} else if hb < fromIndex {
        \t\t\t\tbreak
        \t\t\t}
        \t\t}
        \t\tif fromIndex >= bsInlineBits {
        \t\t\treturn int32(fromIndex)
        \t\t}
        \t}
        \tnUsed := b.bsUsedWords()
        \tw := fromIndex >> 6
        \tif w >= nUsed {
        \t\treturn int32(fromIndex)
        \t}
        \tword := ^b.bsWordAt(w) & (^uint64(0) >> (63 - (uint(fromIndex) & 63)))
        \tfor {
        \t\tif word != 0 {
        \t\t\treturn int32((w+1)*64 - 1 - bits.LeadingZeros64(word))
        \t\t}
        \t\tif w == 0 {
        \t\t\treturn -1
        \t\t}
        \t\tw--
        \t\tword = ^b.bsWordAt(w)
        \t}
        }
        func bitSetOp_Length(b *Ext_java_util_BitSet) int32 {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \treturn int32(b.bsLength())
        }
        func bitSetOp_IsEmpty(b *Ext_java_util_BitSet) bool {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tif b.hiLen > 0 {
        \t\treturn false
        \t}
        \treturn b.bsUsedWords() == 0
        }
        func bitSetOp_Intersects(b, o *Ext_java_util_BitSet) bool {
        \tif b == nil || o == nil {
        \t\tThrowNPE()
        \t}
        \tn := min(b.bsUsedWords(), o.bsUsedWords())
        \tfor i := 0; i < n; i++ {
        \t\tif b.bsWordAt(i)&o.bsWordAt(i) != 0 {
        \t\t\treturn true
        \t\t}
        \t}
        \tfor i := int32(0); i < b.hiLen; i++ {
        \t\tif bitSetOp_Get_1(o, b.hi[i]) {
        \t\t\treturn true
        \t\t}
        \t}
        \tif b.words != nil {
        \t\tfor i := int32(0); i < o.hiLen; i++ {
        \t\t\tif bitSetOp_Get_1(b, o.hi[i]) {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t}
        \treturn false
        }
        func bitSetOp_Cardinality(b *Ext_java_util_BitSet) int32 {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tcnt := int(b.hiLen)
        \tif b.words == nil {
        \t\tfor i := 0; i < bsInlineWords; i++ {
        \t\t\tcnt += bits.OnesCount64(b.buf[i])
        \t\t}
        \t} else {
        \t\tfor _, w := range b.words {
        \t\t\tcnt += bits.OnesCount64(w)
        \t\t}
        \t}
        \treturn int32(cnt)
        }
        func bitSetOp_And(b, o *Ext_java_util_BitSet) {
        \tif b == nil || o == nil {
        \t\tThrowNPE()
        \t}
        \tif b == o {
        \t\treturn
        \t}
        \tif b.words == nil {
        \t\tfor i := 0; i < bsInlineWords; i++ {
        \t\t\tb.buf[i] &= o.bsWordAt(i)
        \t\t}
        \t\tif b.hiLen > 0 {
        \t\t\twIdx := 0
        \t\t\tn := int(b.hiLen)
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\thb := b.hi[i]
        \t\t\t\tif bitSetOp_Get_1(o, hb) {
        \t\t\t\t\tb.hi[wIdx] = hb
        \t\t\t\t\twIdx++
        \t\t\t\t}
        \t\t\t}
        \t\t\tfor i := wIdx; i < n; i++ {
        \t\t\t\tb.hi[i] = 0
        \t\t\t}
        \t\t\tb.hiLen = int32(wIdx)
        \t\t}
        \t\treturn
        \t}
        \tvar keepHi [bsSparseCap]int32
        \tkeepCnt := 0
        \tfor i := int32(0); i < o.hiLen; i++ {
        \t\thb := o.hi[i]
        \t\tif bitSetOp_Get_1(b, hb) {
        \t\t\tkeepHi[keepCnt] = hb
        \t\t\tkeepCnt++
        \t\t}
        \t}
        \tn := min(len(b.words), o.bsUsedWords())
        \tfor i := 0; i < n; i++ {
        \t\tb.words[i] &= o.bsWordAt(i)
        \t}
        \tb.words = b.words[:n]
        \tb.bsTrim()
        \tfor i := 0; i < keepCnt; i++ {
        \t\tb.bsSet(int(keepHi[i]))
        \t}
        }
        func bitSetOp_Or(b, o *Ext_java_util_BitSet) {
        \tif b == nil || o == nil {
        \t\tThrowNPE()
        \t}
        \tif b == o {
        \t\treturn
        \t}
        \tnO := o.bsUsedWords()
        \tif nO > 0 {
        \t\tif b.words == nil && nO <= bsInlineWords {
        \t\t\tfor i := 0; i < nO; i++ {
        \t\t\t\tb.buf[i] |= o.bsWordAt(i)
        \t\t\t}
        \t\t} else {
        \t\t\tb.bsEnsure(nO - 1)
        \t\t\tfor i := 0; i < nO; i++ {
        \t\t\t\tb.words[i] |= o.bsWordAt(i)
        \t\t\t}
        \t\t}
        \t}
        \tfor i := int32(0); i < o.hiLen; i++ {
        \t\tb.bsSet(int(o.hi[i]))
        \t}
        }
        func bitSetOp_Xor(b, o *Ext_java_util_BitSet) {
        \tif b == nil || o == nil {
        \t\tThrowNPE()
        \t}
        \tif b == o {
        \t\tbitSetOp_Clear(b)
        \t\treturn
        \t}
        \tnO := o.bsUsedWords()
        \tif nO > 0 {
        \t\tif b.words == nil && nO <= bsInlineWords {
        \t\t\tfor i := 0; i < nO; i++ {
        \t\t\t\tb.buf[i] ^= o.bsWordAt(i)
        \t\t\t}
        \t\t} else {
        \t\t\tb.bsEnsure(nO - 1)
        \t\t\tfor i := 0; i < nO; i++ {
        \t\t\t\tb.words[i] ^= o.bsWordAt(i)
        \t\t\t}
        \t\t\tb.bsTrim()
        \t\t}
        \t}
        \tfor i := int32(0); i < o.hiLen; i++ {
        \t\tbitSetOp_Flip_1(b, o.hi[i])
        \t}
        }
        func bitSetOp_AndNot(b, o *Ext_java_util_BitSet) {
        \tif b == nil || o == nil {
        \t\tThrowNPE()
        \t}
        \tif b == o {
        \t\tbitSetOp_Clear(b)
        \t\treturn
        \t}
        \tif b.words == nil {
        \t\tfor i := 0; i < bsInlineWords; i++ {
        \t\t\tb.buf[i] &^= o.bsWordAt(i)
        \t\t}
        \t\tif b.hiLen > 0 {
        \t\t\twIdx := 0
        \t\t\tn := int(b.hiLen)
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\thb := b.hi[i]
        \t\t\t\tif !bitSetOp_Get_1(o, hb) {
        \t\t\t\t\tb.hi[wIdx] = hb
        \t\t\t\t\twIdx++
        \t\t\t\t}
        \t\t\t}
        \t\t\tfor i := wIdx; i < n; i++ {
        \t\t\t\tb.hi[i] = 0
        \t\t\t}
        \t\t\tb.hiLen = int32(wIdx)
        \t\t}
        \t\treturn
        \t}
        \tn := min(len(b.words), o.bsUsedWords())
        \tfor i := 0; i < n; i++ {
        \t\tb.words[i] &^= o.bsWordAt(i)
        \t}
        \tfor i := int32(0); i < o.hiLen; i++ {
        \t\tb.bsClearBit(int(o.hi[i]))
        \t}
        \tb.bsTrim()
        }
        func bitSetOp_HashCode(b *Ext_java_util_BitSet) int32 {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \th := int64(1234)
        \tfor i := int(b.hiLen) - 1; i >= 0; {
        \t\thw := int(b.hi[i]) >> 6
        \t\tvar hword uint64
        \t\tfor i >= 0 && (int(b.hi[i])>>6) == hw {
        \t\t\thword |= uint64(1) << (uint(b.hi[i]) & 63)
        \t\t\ti--
        \t\t}
        \t\th ^= int64(hword) * int64(hw+1)
        \t}
        \tn := b.bsUsedWords()
        \tfor i := n - 1; i >= 0; i-- {
        \t\th ^= int64(b.bsWordAt(i)) * int64(i+1)
        \t}
        \treturn int32((h >> 32) ^ h)
        }
        func bitSetOp_Size(b *Ext_java_util_BitSet) int32 {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \treturn int32(b.bsSize())
        }
        func bitSetEqualsInlineVsHeap(inl, hp *Ext_java_util_BitSet) bool {
        \tnHp := hp.bsUsedWords()
        \tmaxHW := int(inl.hi[inl.hiLen-1]) >> 6
        \tif nHp != maxHW+1 {
        \t\treturn false
        \t}
        \tfor i := 0; i < bsInlineWords; i++ {
        \t\tif inl.buf[i] != hp.bsWordAt(i) {
        \t\t\treturn false
        \t\t}
        \t}
        \tcurW := bsInlineWords
        \tn := int(inl.hiLen)
        \tfor idx := 0; idx < n; {
        \t\thw := int(inl.hi[idx]) >> 6
        \t\tfor curW < hw {
        \t\t\tif hp.words[curW] != 0 {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tcurW++
        \t\t}
        \t\tvar hword uint64
        \t\tfor idx < n && (int(inl.hi[idx])>>6) == hw {
        \t\t\thword |= uint64(1) << (uint(inl.hi[idx]) & 63)
        \t\t\tidx++
        \t\t}
        \t\tif hp.words[hw] != hword {
        \t\t\treturn false
        \t\t}
        \t\tcurW = hw + 1
        \t}
        \treturn true
        }
        func bitSetEqualsPtr(b, ob *Ext_java_util_BitSet) bool {
        \tif b.words == nil && ob.words == nil {
        \t\treturn b.buf == ob.buf && b.hiLen == ob.hiLen && b.hi == ob.hi
        \t}
        \tif b.hiLen == 0 && ob.hiLen == 0 {
        \t\tnB := b.bsUsedWords()
        \t\tnO := ob.bsUsedWords()
        \t\tif nB != nO {
        \t\t\treturn false
        \t\t}
        \t\tfor i := 0; i < nB; i++ {
        \t\t\tif b.bsWordAt(i) != ob.bsWordAt(i) {
        \t\t\t\treturn false
        \t\t\t}
        \t\t}
        \t\treturn true
        \t}
        \tif b.hiLen > 0 {
        \t\treturn bitSetEqualsInlineVsHeap(b, ob)
        \t}
        \treturn bitSetEqualsInlineVsHeap(ob, b)
        }
        func bitSetOp_Equals(b *Ext_java_util_BitSet, o any) bool {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tob, ok := o.(*Ext_java_util_BitSet)
        \tif !ok || ob == nil {
        \t\treturn false
        \t}
        \treturn b == ob || bitSetEqualsPtr(b, ob)
        }
        func bitSetOp_Clone(b *Ext_java_util_BitSet) any {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tn := b.bsUsedWords()
        \tc := &Ext_java_util_BitSet{
        \t\tinitW: b.initW,
        \t\thiLen: b.hiLen,
        \t\thi:    b.hi,
        \t}
        \tc.BindExtSelf_(c)
        \tif n <= bsInlineWords {
        \t\tif b.words == nil {
        \t\t\tc.buf = b.buf
        \t\t} else {
        \t\t\tcopy(c.buf[:], b.words[:n])
        \t\t}
        \t} else {
        \t\tcapW := max(n, int(b.initW))
        \t\tc.words = make([]uint64, n, capW)
        \t\tcopy(c.words, b.words[:n])
        \t}
        \treturn c
        }
        func bitSetOp_String(b *Ext_java_util_BitSet) string {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tvar sb strings.Builder
        \tsb.WriteByte('{')
        \tfirst := true
        \tfor i := b.bsNextSetBit(0); i >= 0; i = b.bsNextSetBit(i + 1) {
        \t\tif !first {
        \t\t\tsb.WriteString(", ")
        \t\t}
        \t\tfirst = false
        \t\tsb.WriteString(strconv.Itoa(i))
        \t}
        \tsb.WriteByte('}')
        \treturn sb.String()
        }

        // --- Iterables, Iterators, AbstractIterator, PeekingIterator, Ordering, ComparisonChain, Range ---

        func iterablesAddAll(dst Ext_java_util_Collection, src Ext_java_lang_Iterable) bool {
        \tCheckNotNull(dst)
        \tCheckNotNull(src)
        \tadded := false
        \tfor it := Iterate(src); it.HasNext(); {
        \t\tif dst.Add_1(it.Next()) { added = true }
        \t}
        \treturn added
        }

        type jFuncIterable struct {
        \tmakeIter func() Ext_java_util_Iterator
        }

        func (*jFuncIterable) Implements_j2gort_Ext_java_lang_Iterable() {}
        func (f *jFuncIterable) Iterator() any { return f.makeIter() }
        func (f *jFuncIterable) ForEach(c Ext_java_util_function_Consumer) {
        \tExtDefault_Iterable_ForEach(f, c)
        }
        func (f *jFuncIterable) String() string {
        \tvar items []any
        \tfor it := f.makeIter(); it.HasNext(); { items = append(items, it.Next()) }
        \treturn sliceJavaString(items)
        }

        func iterablesConcatList(list []Ext_java_lang_Iterable) Ext_java_lang_Iterable {
        \tfor _, x := range list { CheckNotNull(x) }
        \treturn &jFuncIterable{makeIter: func() Ext_java_util_Iterator {
        \t\treturn &jConcatIter{outer: newSliceObjIter(toAnySliceIterables(list))}
        \t}}
        }

        func toAnySliceIterables(list []Ext_java_lang_Iterable) []any {
        \tout := make([]any, len(list))
        \tfor i, v := range list { out[i] = v }
        \treturn out
        }

        func iterablesConcatMany(outer Ext_java_lang_Iterable) Ext_java_lang_Iterable {
        \tCheckNotNull(outer)
        \treturn &jFuncIterable{makeIter: func() Ext_java_util_Iterator {
        \t\treturn &jConcatIter{outer: Iterate(outer)}
        \t}}
        }

        type jConcatIter struct {
        \touter SimpleIter
        \tcur   SimpleIter
        }

        func (*jConcatIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jConcatIter) HasNext() bool {
        \tfor it.cur == nil || !it.cur.HasNext() {
        \t\tif !it.outer.HasNext() { return false }
        \t\tnextIt := CheckNotNull(it.outer.Next())
        \t\tit.cur = Iterate(nextIt)
        \t}
        \treturn true
        }
        func (it *jConcatIter) Next() any {
        \tif !it.HasNext() { panic(javalang.NewNoSuchElementException("exhausted")) }
        \treturn it.cur.Next()
        }
        func (it *jConcatIter) Remove() {
        \tif r, ok := it.cur.(Ext_java_util_Iterator); ok {
        \t\tr.Remove()
        \t\treturn
        \t}
        \tpanic(javalang.NewUnsupportedOperationException("remove"))
        }

        func iterablesContains(it Ext_java_lang_Iterable, target any) bool {
        \tCheckNotNull(it)
        \tif c, ok := any(it).(Ext_java_util_Collection); ok { return c.Contains(target) }
        \tt := ToAny(target)
        \tfor iter := Iterate(it); iter.HasNext(); {
        \t\tif javalang.ObjectsEquals(t, ToAny(iter.Next())) { return true }
        \t}
        \treturn false
        }

        func iterablesFilter(it Ext_java_lang_Iterable, pred Ext_com_google_common_base_Predicate) Ext_java_lang_Iterable {
        \tCheckNotNull(it)
        \tCheckNotNull(pred)
        \treturn &jFuncIterable{makeIter: func() Ext_java_util_Iterator {
        \t\treturn iteratorsFilterSimple(Iterate(it), pred)
        \t}}
        }

        func iteratorsFilter(it Ext_java_util_Iterator, pred Ext_com_google_common_base_Predicate) Ext_com_google_common_collect_UnmodifiableIterator {
        \tCheckNotNull(it)
        \tCheckNotNull(pred)
        \treturn iteratorsFilterSimple(it, pred)
        }

        type jFilterIter struct {
        \tsrc      SimpleIter
        \tpred     Ext_com_google_common_base_Predicate
        \thaveNext bool
        \tnextVal  any
        }

        func iteratorsFilterSimple(src SimpleIter, pred Ext_com_google_common_base_Predicate) *jFilterIter {
        \treturn &jFilterIter{src: src, pred: pred}
        }

        func (*jFilterIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (*jFilterIter) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator() {}
        func (*jFilterIter) GetExt_com_google_common_collect_UnmodifiableIteratorBase() *Ext_com_google_common_collect_UnmodifiableIteratorBase {
        \treturn nil
        }
        func (it *jFilterIter) HasNext() bool {
        \tif it.haveNext { return true }
        \tfor it.src.HasNext() {
        \t\tv := it.src.Next()
        \t\tif it.pred.Apply(v) {
        \t\t\tit.nextVal = v
        \t\t\tit.haveNext = true
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }
        func (it *jFilterIter) Next() any {
        \tif !it.HasNext() { panic(javalang.NewNoSuchElementException("exhausted")) }
        \tit.haveNext = false
        \tv := it.nextVal
        \tit.nextVal = nil
        \treturn v
        }
        func (*jFilterIter) Remove() { panic(javalang.NewUnsupportedOperationException("remove")) }

        func iterablesGet(it Ext_java_lang_Iterable, idx int32) any {
        \tCheckNotNull(it)
        \tif idx < 0 { panic(javalang.NewIndexOutOfBoundsException("negative index")) }
        \tif l, ok := any(it).(Ext_java_util_List); ok { return l.Get(idx) }
        \ti := int32(0)
        \tfor iter := Iterate(it); iter.HasNext(); i++ {
        \t\tv := iter.Next()
        \t\tif i == idx { return v }
        \t}
        \tpanic(javalang.NewIndexOutOfBoundsExceptionForIndex(int(idx), int(i)))
        }

        func iterablesGetFirst(it Ext_java_lang_Iterable, def any) any {
        \tCheckNotNull(it)
        \titer := NewIterCursor(it)
        \tif iter.HasNext() { return iter.Next() }
        \treturn ToAny(def)
        }

        func iterablesGetLast(it Ext_java_lang_Iterable) any {
        \tCheckNotNull(it)
        \tif l, ok := any(it).(Ext_java_util_List); ok {
        \t\tif l.IsEmpty() { panic(javalang.NewNoSuchElementException("empty")) }
        \t\treturn l.Get(l.Size() - 1)
        \t}
        \titer := NewIterCursor(it)
        \tif !iter.HasNext() { panic(javalang.NewNoSuchElementException("empty")) }
        \tlast := iter.Next()
        \tfor iter.HasNext() { last = iter.Next() }
        \treturn last
        }

        func iterablesGetOnlyElement(it Ext_java_lang_Iterable) any {
        \tCheckNotNull(it)
        \titer := NewIterCursor(it)
        \tif !iter.HasNext() { panic(javalang.NewNoSuchElementException("empty")) }
        \tfirst := iter.Next()
        \tif iter.HasNext() {
        \t\tpanic(javalang.NewIllegalArgumentException("expected one element but was multiple"))
        \t}
        \treturn first
        }

        func iterablesSkip(it Ext_java_lang_Iterable, n int32) Ext_java_lang_Iterable {
        \tCheckNotNull(it)
        \tCheckArgument(n >= 0)
        \treturn &jFuncIterable{makeIter: func() Ext_java_util_Iterator {
        \t\titer := Iterate(it)
        \t\tfor i := int32(0); i < n && iter.HasNext(); i++ { _ = iter.Next() }
        \t\treturn &jSimpleIterAdapter{it: iter}
        \t}}
        }

        func iterablesTransform(it Ext_java_lang_Iterable, fn Ext_com_google_common_base_Function) Ext_java_lang_Iterable {
        \tCheckNotNull(it)
        \tCheckNotNull(fn)
        \treturn &jFuncIterable{makeIter: func() Ext_java_util_Iterator {
        \t\treturn &jTransformIter{src: Iterate(it), fn: fn}
        \t}}
        }

        type jSimpleIterAdapter struct {
        \tit SimpleIter
        }

        func (*jSimpleIterAdapter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (a *jSimpleIterAdapter) HasNext() bool { return a.it.HasNext() }
        func (a *jSimpleIterAdapter) Next() any     { return a.it.Next() }
        func (a *jSimpleIterAdapter) Remove() {
        \tif r, ok := a.it.(Ext_java_util_Iterator); ok { r.Remove(); return }
        \tpanic(javalang.NewUnsupportedOperationException("remove"))
        }

        type jTransformIter struct {
        \tsrc SimpleIter
        \tfn  Ext_com_google_common_base_Function
        }

        func (*jTransformIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (t *jTransformIter) HasNext() bool { return t.src.HasNext() }
        func (t *jTransformIter) Next() any     { return t.fn.Apply(t.src.Next()) }
        func (t *jTransformIter) Remove() {
        \tif r, ok := t.src.(Ext_java_util_Iterator); ok { r.Remove(); return }
        \tpanic(javalang.NewUnsupportedOperationException("remove"))
        }

        type jPeekingIter struct {
        \tsrc       Ext_java_util_Iterator
        \thasPeeked bool
        \tpeeked    any
        }

        func newPeekingIter(src Ext_java_util_Iterator) Ext_com_google_common_collect_PeekingIterator {
        \tCheckNotNull(src)
        \tif p, ok := src.(Ext_com_google_common_collect_PeekingIterator); ok { return p }
        \treturn &jPeekingIter{src: src}
        }

        func (*jPeekingIter) Implements_j2gort_Ext_java_util_Iterator() {}
        func (*jPeekingIter) Implements_j2gort_Ext_com_google_common_collect_PeekingIterator() {}
        func (p *jPeekingIter) HasNext() bool { return p.hasPeeked || p.src.HasNext() }
        func (p *jPeekingIter) Next() any {
        \tif !p.hasPeeked { return p.src.Next() }
        \tv := p.peeked
        \tp.hasPeeked = false
        \tp.peeked = nil
        \treturn v
        }
        func (p *jPeekingIter) Peek() any {
        \tif !p.hasPeeked {
        \t\tp.peeked = p.src.Next()
        \t\tp.hasPeeked = true
        \t}
        \treturn p.peeked
        }
        func (p *jPeekingIter) Remove() {
        \tif p.hasPeeked { panic(javalang.NewIllegalStateException("cannot remove after peek")) }
        \tp.src.Remove()
        }

        type jAbstractIterState struct {
        \tstate   int // 0=NOT_READY, 1=READY, 2=DONE, 3=FAILED
        \tnextVal any
        }

        func abstractIterHasNext(b *Ext_com_google_common_collect_AbstractIteratorBase) bool {
        \tif b.Core_ == nil { b.Core_ = &jAbstractIterState{} }
        \tst := b.Core_.(*jAbstractIterState)
        \tif st.state == 3 { panic(javalang.NewIllegalStateException("failed iterator")) }
        \tif st.state == 2 { return false }
        \tif st.state == 1 { return true }
        \tst.state = 3
        \tst.nextVal = b.Self_.ComputeNext()
        \tif st.state != 2 {
        \t\tst.state = 1
        \t\treturn true
        \t}
        \treturn false
        }

        func abstractIterNext(b *Ext_com_google_common_collect_AbstractIteratorBase) any {
        \tif !b.Self_.HasNext() { panic(javalang.NewNoSuchElementException("exhausted")) }
        \tst := b.Core_.(*jAbstractIterState)
        \tst.state = 0
        \tv := st.nextVal
        \tst.nextVal = nil
        \treturn v
        }

        func newOrdering(cmp func(a, b any) int32) Ext_com_google_common_collect_Ordering {
        \to := &Ext_com_google_common_collect_OrderingBase{}
        \to.BindExtSelf_(o)
        \to.Core_ = cmp
        \treturn o
        }

        type jComparisonChain struct {
        \tres int32
        }

        var (
        \tcomparisonChainActive = newComparisonChainObj(0)
        \tcomparisonChainLess   = newComparisonChainObj(-1)
        \tcomparisonChainGreater = newComparisonChainObj(1)
        )

        func newComparisonChainObj(res int32) Ext_com_google_common_collect_ComparisonChain {
        \tc := &Ext_com_google_common_collect_ComparisonChainBase{}
        \tc.BindExtSelf_(c)
        \tc.Core_ = &jComparisonChain{res: res}
        \treturn c
        }

        func (c *jComparisonChain) cmpInt(a, b int32) Ext_com_google_common_collect_ComparisonChain {
        \tif c.res != 0 {
        \t\tif c.res < 0 { return comparisonChainLess }
        \t\treturn comparisonChainGreater
        \t}
        \tif a < b { return comparisonChainLess }
        \tif a > b { return comparisonChainGreater }
        \treturn comparisonChainActive
        }

        func (c *jComparisonChain) cmpAny(a, b any) Ext_com_google_common_collect_ComparisonChain {
        \tif c.res != 0 {
        \t\tif c.res < 0 { return comparisonChainLess }
        \t\treturn comparisonChainGreater
        \t}
        \tr := CompareValues(a, b)
        \tif r < 0 { return comparisonChainLess }
        \tif r > 0 { return comparisonChainGreater }
        \treturn comparisonChainActive
        }

        type jRangeCore struct {
        \tlower any
        \tupper any
        }

        type jRangeSetCore struct {
        \tranges []*jRangeCore
        }

        func collectionsShuffle(lst Ext_java_util_List, rnd *Ext_java_util_Random) {
        \tCheckNotNull(lst)
        \tCheckNotNullPtr(rnd)
        \tif rnd.Core_ == nil { rnd.Core_ = javautil.NewRandom() }
        \tr := rnd.Core_.(*javautil.Random)
        \tfor i := lst.Size(); i > 1; i-- {
        \t\tj := r.NextIntBound(i)
        \t\ttmp := lst.Get(i - 1)
        \t\tlst.Set(i-1, lst.Get(j))
        \t\tlst.Set(j, tmp)
        \t}
        }

        // --- Streams, IntStream, Collectors ---

        type jStream struct {
        \tnext func() (any, bool)
        }

        func streamFromSlice(items []any) Ext_java_util_stream_Stream {
        \tidx := 0
        \treturn &jStream{next: func() (any, bool) {
        \t\tif idx >= len(items) { return nil, false }
        \t\tv := items[idx]
        \t\tidx++
        \t\treturn ToAny(v), true
        \t}}
        }

        func streamFromIterable(it Ext_java_lang_Iterable) Ext_java_util_stream_Stream {
        \tCheckNotNull(it)
        \tif si, ok := any(it).(SliceIterable); ok {
        \t\treturn streamFromSlice(si.IterSlice_())
        \t}
        \tvar cur JIterCursor
        \tinited := false
        \treturn &jStream{next: func() (any, bool) {
        \t\tif !inited {
        \t\t\tcur = NewIterCursor(it)
        \t\t\tinited = true
        \t\t}
        \t\tif !cur.HasNext() { return nil, false }
        \t\treturn ToAny(cur.Next()), true
        \t}}
        }

        func streamConcat(a, b Ext_java_util_stream_Stream) Ext_java_util_stream_Stream {
        \tCheckNotNull(a)
        \tCheckNotNull(b)
        \ts1 := a.(*jStream)
        \ts2 := b.(*jStream)
        \tdone1 := false
        \treturn &jStream{next: func() (any, bool) {
        \t\tif !done1 {
        \t\t\tif v, ok := s1.next(); ok { return v, true }
        \t\t\tdone1 = true
        \t\t}
        \t\treturn s2.next()
        \t}}
        }

        func streamMapWithIndex(s Ext_java_util_stream_Stream, fn Ext_com_google_common_collect_Streams_FunctionWithIndex) Ext_java_util_stream_Stream {
        \tCheckNotNull(s)
        \tCheckNotNull(fn)
        \tsrc := s.(*jStream)
        \tvar idx int64
        \treturn &jStream{next: func() (any, bool) {
        \t\tv, ok := src.next()
        \t\tif !ok { return nil, false }
        \t\ti := idx
        \t\tidx++
        \t\treturn ToAny(fn.Apply(v, i)), true
        \t}}
        }

        func (*jStream) Implements_j2gort_Ext_java_lang_AutoCloseable() {}
        func (*jStream) Implements_j2gort_Ext_java_util_stream_BaseStream() {}
        func (*jStream) Implements_j2gort_Ext_java_util_stream_Stream() {}
        func (*jStream) Close_0() {}

        func (s *jStream) drain() []any {
        \tvar out []any
        \tfor {
        \t\tv, ok := s.next()
        \t\tif !ok { break }
        \t\tout = append(out, v)
        \t}
        \treturn out
        }

        func (s *jStream) Filter(pred Ext_java_util_function_Predicate) Ext_java_util_stream_Stream {
        \tCheckNotNull(pred)
        \treturn &jStream{next: func() (any, bool) {
        \t\tfor {
        \t\t\tv, ok := s.next()
        \t\t\tif !ok { return nil, false }
        \t\t\tif pred.Test(v) { return v, true }
        \t\t}
        \t}}
        }

        func (s *jStream) Map(fn Ext_java_util_function_Function) Ext_java_util_stream_Stream {
        \tCheckNotNull(fn)
        \treturn &jStream{next: func() (any, bool) {
        \t\tv, ok := s.next()
        \t\tif !ok { return nil, false }
        \t\treturn ToAny(fn.Apply(v)), true
        \t}}
        }

        func (s *jStream) FlatMap(fn Ext_java_util_function_Function) Ext_java_util_stream_Stream {
        \tCheckNotNull(fn)
        \tvar cur *jStream
        \treturn &jStream{next: func() (any, bool) {
        \t\tfor {
        \t\t\tif cur != nil {
        \t\t\t\tif v, ok := cur.next(); ok { return v, true }
        \t\t\t\tcur = nil
        \t\t\t}
        \t\t\touter, ok := s.next()
        \t\t\tif !ok { return nil, false }
        \t\t\tsub := fn.Apply(outer)
        \t\t\tif sub != nil && !javalang.IsNil(sub) {
        \t\t\t\tcur = sub.(*jStream)
        \t\t\t}
        \t\t}
        \t}}
        }

        func (s *jStream) Distinct() Ext_java_util_stream_Stream {
        \tseen := NewJMapCore(false)
        \treturn &jStream{next: func() (any, bool) {
        \t\tfor {
        \t\t\tv, ok := s.next()
        \t\t\tif !ok { return nil, false }
        \t\t\tif !seen.ContainsKey(v) {
        \t\t\t\tseen.Put(v, true)
        \t\t\t\treturn v, true
        \t\t\t}
        \t\t}
        \t}}
        }

        func (s *jStream) Sorted_0() Ext_java_util_stream_Stream {
        \treturn s.Sorted_1(nil)
        }

        func (s *jStream) Sorted_1(cmp Ext_java_util_Comparator) Ext_java_util_stream_Stream {
        \tfn := comparatorFuncOrNatural(cmp)
        \tvar buf []any
        \tsorted := false
        \tidx := 0
        \treturn &jStream{next: func() (any, bool) {
        \t\tif !sorted {
        \t\t\tbuf = s.drain()
        \t\t\tsort.SliceStable(buf, func(i, j int) bool { return fn(buf[i], buf[j]) < 0 })
        \t\t\tsorted = true
        \t\t}
        \t\tif idx >= len(buf) { return nil, false }
        \t\tv := buf[idx]
        \t\tidx++
        \t\treturn v, true
        \t}}
        }

        func (s *jStream) Peek(action Ext_java_util_function_Consumer) Ext_java_util_stream_Stream {
        \tCheckNotNull(action)
        \treturn &jStream{next: func() (any, bool) {
        \t\tv, ok := s.next()
        \t\tif ok { action.Accept(v) }
        \t\treturn v, ok
        \t}}
        }

        func (s *jStream) Limit(maxSize int64) Ext_java_util_stream_Stream {
        \tCheckArgument(maxSize >= 0)
        \tvar count int64
        \treturn &jStream{next: func() (any, bool) {
        \t\tif count >= maxSize { return nil, false }
        \t\tv, ok := s.next()
        \t\tif !ok { return nil, false }
        \t\tcount++
        \t\treturn v, true
        \t}}
        }

        func (s *jStream) Skip(n int64) Ext_java_util_stream_Stream {
        \tCheckArgument(n >= 0)
        \tskipped := false
        \treturn &jStream{next: func() (any, bool) {
        \t\tif !skipped {
        \t\t\tfor i := int64(0); i < n; i++ {
        \t\t\t\tif _, ok := s.next(); !ok {
        \t\t\t\t\tskipped = true
        \t\t\t\t\treturn nil, false
        \t\t\t\t}
        \t\t\t}
        \t\t\tskipped = true
        \t\t}
        \t\treturn s.next()
        \t}}
        }

        func (s *jStream) ForEach(action Ext_java_util_function_Consumer) {
        \tCheckNotNull(action)
        \tfor {
        \t\tv, ok := s.next()
        \t\tif !ok { break }
        \t\taction.Accept(v)
        \t}
        }

        func (s *jStream) ForEachOrdered(action Ext_java_util_function_Consumer) {
        \ts.ForEach(action)
        }

        func (s *jStream) ToArray_1(gen Ext_java_util_function_IntFunction) []any {
        \tCheckNotNull(gen)
        \titems := s.drain()
        \tarr := gen.Apply(int32(len(items)))
        \tif out, ok := arr.([]any); ok {
        \t\tcopy(out, items)
        \t\treturn out
        \t}
        \treturn items
        }

        func (s *jStream) Reduce_2(identity any, op Ext_java_util_function_BinaryOperator) any {
        \tCheckNotNull(op)
        \tacc := ToAny(identity)
        \tfor {
        \t\tv, ok := s.next()
        \t\tif !ok { break }
        \t\tacc = ToAny(op.Apply(acc, v))
        \t}
        \treturn acc
        }

        func (s *jStream) Collect_1(collector Ext_java_util_stream_Collector) any {
        \tCheckNotNull(collector)
        \treturn collector.(*jCollector).fn(s.drain())
        }

        func (s *jStream) AnyMatch(pred Ext_java_util_function_Predicate) bool {
        \tCheckNotNull(pred)
        \tfor {
        \t\tv, ok := s.next()
        \t\tif !ok { return false }
        \t\tif pred.Test(v) { return true }
        \t}
        }

        func (s *jStream) AllMatch(pred Ext_java_util_function_Predicate) bool {
        \tCheckNotNull(pred)
        \tfor {
        \t\tv, ok := s.next()
        \t\tif !ok { return true }
        \t\tif !pred.Test(v) { return false }
        \t}
        }

        func (s *jStream) FindFirst() *Ext_java_util_Optional {
        \tv, ok := s.next()
        \topt := &Ext_java_util_Optional{}
        \tif ok {
        \t\topt.Core_ = CheckNotNull(v)
        \t}
        \treturn opt
        }

        type jStreamBuilder struct {
        \titems []any
        }

        func (*jStreamBuilder) Implements_j2gort_Ext_java_util_function_Consumer() {}
        func (*jStreamBuilder) Implements_j2gort_Ext_java_util_stream_Stream_Builder() {}
        func (b *jStreamBuilder) Accept(v any) { b.items = append(b.items, ToAny(v)) }
        func (b *jStreamBuilder) Add(v any) Ext_java_util_stream_Stream_Builder {
        \tb.Accept(v)
        \treturn b
        }
        func (b *jStreamBuilder) Build() Ext_java_util_stream_Stream {
        \treturn streamFromSlice(b.items)
        }

        type jIntStream struct {
        \tstart int32
        \tend   int32
        }

        func (*jIntStream) Implements_j2gort_Ext_java_lang_AutoCloseable() {}
        func (*jIntStream) Implements_j2gort_Ext_java_util_stream_BaseStream() {}
        func (*jIntStream) Implements_j2gort_Ext_java_util_stream_IntStream() {}
        func (*jIntStream) Close_0() {}
        func (is *jIntStream) MapToObj(fn Ext_java_util_function_IntFunction) Ext_java_util_stream_Stream {
        \tCheckNotNull(fn)
        \tcur := is.start
        \treturn &jStream{next: func() (any, bool) {
        \t\tif cur >= is.end { return nil, false }
        \t\ti := cur
        \t\tcur++
        \t\treturn ToAny(fn.Apply(i)), true
        \t}}
        }

        type jCollector struct {
        \tfn func(items []any) any
        }

        func (*jCollector) Implements_j2gort_Ext_java_util_stream_Collector() {}

        func collectorToList() Ext_java_util_stream_Collector {
        \treturn &jCollector{fn: func(items []any) any { return wrapArrayListOwn(items) }}
        }

        func collectorToImmutableList() Ext_java_util_stream_Collector {
        \treturn &jCollector{fn: func(items []any) any {
        \t\tfor _, x := range items { CheckNotNull(x) }
        \t\treturn wrapImmutableListOwn(items)
        \t}}
        }

        func collectorToImmutableSet() Ext_java_util_stream_Collector {
        \treturn &jCollector{fn: func(items []any) any {
        \t\tfor _, x := range items { CheckNotNull(x) }
        \t\treturn wrapImmutableSet(items)
        \t}}
        }

        func collectorToImmutableSortedSet(cmp Ext_java_util_Comparator) Ext_java_util_stream_Collector {
        \treturn &jCollector{fn: func(items []any) any {
        \t\tif len(items) == 0 && isNaturalOrderComparator(cmp) {
        \t\t\treturn emptyImmutableSortedSet_
        \t\t}
        \t\tmc := newSortedCoreFromCmp(cmp)
        \t\tfor _, x := range items { mc.Put(CheckNotNull(x), true) }
        \t\treturn wrapImmutableSortedSet(mc)
        \t}}
        }

        func collectorToImmutableMultiset() Ext_java_util_stream_Collector {
        \treturn &jCollector{fn: func(items []any) any {
        \t\tc := newMultisetCore(false, nil)
        \t\tfor _, x := range items { multisetOp_add_2(c, CheckNotNull(x), 1) }
        \t\tc.immutable = true
        \t\treturn wrapImmutableMultiset(c)
        \t}}
        }

        func collectorToImmutableMap(kFn, vFn Ext_java_util_function_Function) Ext_java_util_stream_Collector {
        \tCheckNotNull(kFn)
        \tCheckNotNull(vFn)
        \treturn &jCollector{fn: func(items []any) any {
        \t\tif len(items) == 0 {
        \t\t\treturn emptyImmutableMap_
        \t\t}
        \t\tobj := newImmutableMapShell()
        \t\tmc := &obj.inlineCore
        \t\tmc.InitExpectedSize(int32(len(items)))
        \t\tfor _, x := range items {
        \t\t\tk := CheckNotNull(ToAny(kFn.Apply(x)))
        \t\t\tv := CheckNotNull(ToAny(vFn.Apply(x)))
        \t\t\tif mc.ContainsKey(k) {
        \t\t\t\tpanic(javalang.NewIllegalArgumentException("Multiple entries with same key: " + StrVal(k)))
        \t\t\t}
        \t\t\tmc.Put(k, v)
        \t\t}
        \t\tmc.sealImmutable()
        \t\treturn obj
        \t}}
        }

        func collectorToImmutableSortedMap(cmp Ext_java_util_Comparator, kFn, vFn Ext_java_util_function_Function) Ext_java_util_stream_Collector {
        \tCheckNotNull(cmp)
        \tCheckNotNull(kFn)
        \tCheckNotNull(vFn)
        \treturn &jCollector{fn: func(items []any) any {
        \t\tif len(items) == 0 && isNaturalOrderComparator(cmp) {
        \t\t\treturn emptyImmutableSortedMap_
        \t\t}
        \t\tobj := &Ext_com_google_common_collect_ImmutableSortedMap{}
        \t\tobj.BindExtSelf_(obj)
        \t\tmc := initSortedCoreCmp(&obj.inlineCore, cmp)
        \t\tfor _, x := range items {
        \t\t\tk := CheckNotNull(ToAny(kFn.Apply(x)))
        \t\t\tv := CheckNotNull(ToAny(vFn.Apply(x)))
        \t\t\tif mc.ContainsKey(k) {
        \t\t\t\tpanic(javalang.NewIllegalArgumentException("Multiple entries with same key: " + StrVal(k)))
        \t\t\t}
        \t\t\tmc.Put(k, v)
        \t\t}
        \t\tmc.sealImmutable()
        \t\treturn obj
        \t}}
        }

        func collectorToCollection(sup Ext_java_util_function_Supplier) Ext_java_util_stream_Collector {
        \tCheckNotNull(sup)
        \treturn &jCollector{fn: func(items []any) any {
        \t\tcol := Cast[Ext_java_util_Collection](sup.Get())
        \t\tfor _, x := range items { col.Add_1(x) }
        \t\treturn col
        \t}}
        }

        func collectorJoining(delim any) Ext_java_util_stream_Collector {
        \td := StrVal(CheckNotNull(delim))
        \treturn &jCollector{fn: func(items []any) any {
        \t\tvar sb strings.Builder
        \t\tfor i, x := range items {
        \t\t\tif i > 0 { sb.WriteString(d) }
        \t\t\tsb.WriteString(StrVal(x))
        \t\t}
        \t\treturn sb.String()
        \t}}
        }

        func collectorGroupingBy(classifier Ext_java_util_function_Function) Ext_java_util_stream_Collector {
        \tCheckNotNull(classifier)
        \treturn &jCollector{fn: func(items []any) any {
        \t\tmc := NewJMapCore(false)
        \t\tfor _, x := range items {
        \t\t\tk := CheckNotNull(ToAny(classifier.Apply(x)))
        \t\t\tlst := mc.Get(k)
        \t\t\tif lst == nil {
        \t\t\t\tlst = wrapArrayList(nil)
        \t\t\t\tmc.Put(k, lst)
        \t\t\t}
        \t\t\tlst.(Ext_java_util_List).Add_1(x)
        \t\t}
        \t\treturn wrapLinkedHashMap(mc)
        \t}}
        }

        func collectorToMap(kFn, vFn Ext_java_util_function_Function) Ext_java_util_stream_Collector {
        \tCheckNotNull(kFn)
        \tCheckNotNull(vFn)
        \treturn &jCollector{fn: func(items []any) any {
        \t\tmc := NewJMapCore(false)
        \t\tfor _, x := range items {
        \t\t\tk := ToAny(kFn.Apply(x))
        \t\t\tv := CheckNotNull(ToAny(vFn.Apply(x)))
        \t\t\tif mc.ContainsKey(k) {
        \t\t\t\tpanic(javalang.NewIllegalStateException("Duplicate key " + StrVal(k)))
        \t\t\t}
        \t\t\tmc.Put(k, v)
        \t\t}
        \t\treturn wrapLinkedHashMap(mc)
        \t}}
        }

        // --- Cache & Executor / Future / AtomicIntegerArray ---

        type jCacheBuilderCore struct {
        \tmaxSize  int64
        \tlistener Ext_com_google_common_cache_RemovalListener
        }

        type jFuncCacheLoader struct {
        \tExt_com_google_common_cache_CacheLoaderBase
        \tfn Ext_com_google_common_base_Function
        }

        func (l *jFuncCacheLoader) Load(k any) any { return l.fn.Apply(k) }

        type jLoadingCache struct {
        \tmu       sync.Mutex
        \tmc       *JMapCore
        \tmaxSize  int64
        \tlistener Ext_com_google_common_cache_RemovalListener
        \tloader   Ext_com_google_common_cache_CacheLoader
        }

        func newLoadingCache(cfg *jCacheBuilderCore, loader Ext_com_google_common_cache_CacheLoader) Ext_com_google_common_cache_LoadingCache {
        \treturn &jLoadingCache{
        \t\tmc:       NewJMapCore(false),
        \t\tmaxSize:  cfg.maxSize,
        \t\tlistener: cfg.listener,
        \t\tloader:   loader,
        \t}
        }

        func (*jLoadingCache) Implements_j2gort_Ext_com_google_common_cache_Cache() {}
        func (*jLoadingCache) Implements_j2gort_Ext_com_google_common_cache_LoadingCache() {}
        func (*jLoadingCache) Implements_j2gort_Ext_java_util_function_Function() {}
        func (*jLoadingCache) Implements_j2gort_Ext_com_google_common_base_Function() {}
        func (c *jLoadingCache) Equals(o any) bool { return RefEquals(c, o) }
        func (c *jLoadingCache) Apply(k any) any   { return c.GetUnchecked(k) }

        func (c *jLoadingCache) GetUnchecked(k any) any {
        \tk = CheckNotNull(ToAny(k))
        \tc.mu.Lock()
        \tif c.mc.ContainsKey(k) {
        \t\tv := c.mc.Get(k)
        \t\tc.mu.Unlock()
        \t\treturn v
        \t}
        \tc.mu.Unlock()
        \tif c.maxSize == 0 {
        \t\tv := CheckNotNull(ToAny(c.loader.Load(k)))
        \t\tif c.listener != nil {
        \t\t\tn := &Ext_com_google_common_cache_RemovalNotification{}
        \t\t\tn.BindExtSelf_(n)
        \t\t\tn.Core_ = &jSimpleEntry{k: k, v: v, imm: true}
        \t\t\tc.listener.OnRemoval(n)
        \t\t}
        \t\treturn v
        \t}
        \tv := CheckNotNull(ToAny(c.loader.Load(k)))
        \tvar evicted *jSimpleEntry
        \tc.mu.Lock()
        \tif c.maxSize > 0 && int64(c.mc.Size()) >= c.maxSize && !c.mc.ContainsKey(k) {
        \t\tek := c.mc.FirstKey()
        \t\tev := c.mc.Remove(ek)
        \t\tif c.listener != nil {
        \t\t\tevicted = &jSimpleEntry{k: ek, v: ev, imm: true}
        \t\t}
        \t}
        \tc.mc.Put(k, v)
        \tc.mu.Unlock()
        \tif evicted != nil {
        \t\tn := &Ext_com_google_common_cache_RemovalNotification{}
        \t\tn.BindExtSelf_(n)
        \t\tn.Core_ = evicted
        \t\tc.listener.OnRemoval(n)
        \t}
        \treturn v
        }

        type jDirectExecutor struct{}

        func (*jDirectExecutor) Implements_j2gort_Ext_java_util_concurrent_Executor() {}
        func (*jDirectExecutor) Implements_j2gort_Ext_java_lang_AutoCloseable() {}
        func (*jDirectExecutor) Implements_j2gort_Ext_java_util_concurrent_ExecutorService() {}
        func (*jDirectExecutor) Implements_j2gort_Ext_com_google_common_util_concurrent_ListeningExecutorService() {}
        func (*jDirectExecutor) Shutdown() {}
        func (*jDirectExecutor) Close_0()  {}
        func (*jDirectExecutor) Submit_1_Ext_java_util_concurrent_Callable(task Ext_java_util_concurrent_Callable) any {
        \treturn execSubmitCallable(task)
        }

        type jThreadPoolExecutor struct {
        \ttasks     chan func()
        \twg        sync.WaitGroup
        \tcloseOnce sync.Once
        }

        func newThreadPoolExecutor(threads int32) Ext_java_util_concurrent_ExecutorService {
        \tif threads <= 1 {
        \t\treturn &jDirectExecutor{}
        \t}
        \tn := int(threads)
        \tbuf := n * 16
        \tif buf < 64 {
        \t\tbuf = 64
        \t}
        \tp := &jThreadPoolExecutor{tasks: make(chan func(), buf)}
        \tp.wg.Add(n)
        \tfor i := 0; i < n; i++ {
        \t\tgo func() {
        \t\t\tdefer p.wg.Done()
        \t\t\tfor fn := range p.tasks {
        \t\t\t\tfn()
        \t\t\t}
        \t\t}()
        \t}
        \treturn p
        }

        func (*jThreadPoolExecutor) Implements_j2gort_Ext_java_util_concurrent_Executor() {}
        func (*jThreadPoolExecutor) Implements_j2gort_Ext_java_lang_AutoCloseable() {}
        func (*jThreadPoolExecutor) Implements_j2gort_Ext_java_util_concurrent_ExecutorService() {}
        func (*jThreadPoolExecutor) Implements_j2gort_Ext_com_google_common_util_concurrent_ListeningExecutorService() {}
        func (p *jThreadPoolExecutor) Shutdown() {
        \tp.closeOnce.Do(func() {
        \t\tclose(p.tasks)
        \t})
        }
        func (p *jThreadPoolExecutor) Close_0() {
        \tp.Shutdown()
        \tp.wg.Wait()
        }
        func (p *jThreadPoolExecutor) Submit_1_Ext_java_util_concurrent_Callable(task Ext_java_util_concurrent_Callable) any {
        \tCheckNotNull(task)
        \tf := &jAsyncFuture{done: make(chan struct{})}
        \tnodestore.SuspendSingleOwnerAlloc()
        \tp.tasks <- func() {
        \t\tdefer func() {
        \t\t\tif r := recover(); r != nil {
        \t\t\t\tf.err = javalang.RecoverException(r)
        \t\t\t}
        \t\t\tnodestore.ResumeSingleOwnerAlloc()
        \t\t\tclose(f.done)
        \t\t}()
        \t\tf.val = ToAny(task.Call())
        \t}
        \treturn f
        }

        type jAsyncFuture struct {
        \tdone chan struct{}
        \tval  any
        \terr  javalang.Throwable
        }

        func (*jAsyncFuture) Implements_j2gort_Ext_java_util_concurrent_Future() {}
        func (*jAsyncFuture) Implements_j2gort_Ext_com_google_common_util_concurrent_ListenableFuture() {}
        func (f *jAsyncFuture) Get_0() any {
        \t<-f.done
        \tif f.err != nil {
        \t\tex := &Ext_java_util_concurrent_ExecutionExceptionBase{}
        \t\tex.BindExtSelf_(ex)
        \t\tInitUserThrowable(ex, "java.util.concurrent.ExecutionException", f.err)
        \t\tpanic(ex)
        \t}
        \treturn f.val
        }
        func (f *jAsyncFuture) Get_2(_ int64, _ *Ext_java_util_concurrent_TimeUnit) any {
        \treturn f.Get_0()
        }

        type jListeningExecWrapper struct {
        \texec Ext_java_util_concurrent_ExecutorService
        }

        func (*jListeningExecWrapper) Implements_j2gort_Ext_java_util_concurrent_Executor() {}
        func (*jListeningExecWrapper) Implements_j2gort_Ext_java_lang_AutoCloseable() {}
        func (*jListeningExecWrapper) Implements_j2gort_Ext_java_util_concurrent_ExecutorService() {}
        func (*jListeningExecWrapper) Implements_j2gort_Ext_com_google_common_util_concurrent_ListeningExecutorService() {}
        func (w *jListeningExecWrapper) Shutdown() { w.exec.Shutdown() }
        func (w *jListeningExecWrapper) Close_0()  { w.exec.Close_0() }
        func (w *jListeningExecWrapper) Submit_1_Ext_java_util_concurrent_Callable(task Ext_java_util_concurrent_Callable) any {
        \treturn w.exec.Submit_1_Ext_java_util_concurrent_Callable(task)
        }

        type jCompletedFuture struct {
        \tval any
        \terr javalang.Throwable
        }

        func (*jCompletedFuture) Implements_j2gort_Ext_java_util_concurrent_Future() {}
        func (*jCompletedFuture) Implements_j2gort_Ext_com_google_common_util_concurrent_ListenableFuture() {}
        func (f *jCompletedFuture) Get_0() any {
        \tif f.err != nil {
        \t\tex := &Ext_java_util_concurrent_ExecutionExceptionBase{}
        \t\tex.BindExtSelf_(ex)
        \t\tInitUserThrowable(ex, "java.util.concurrent.ExecutionException", f.err)
        \t\tpanic(ex)
        \t}
        \treturn f.val
        }
        func (f *jCompletedFuture) Get_2(_ int64, _ *Ext_java_util_concurrent_TimeUnit) any {
        \treturn f.Get_0()
        }

        func execSubmitCallable(task Ext_java_util_concurrent_Callable) Ext_com_google_common_util_concurrent_ListenableFuture {
        \tCheckNotNull(task)
        \tvar res any
        \tvar th javalang.Throwable
        \tfunc() {
        \t\tif nodestore.TryEnterSingleOwnerAlloc() {
        \t\t\tdefer nodestore.ExitSingleOwnerAlloc()
        \t\t}
        \t\tdefer func() {
        \t\t\tif r := recover(); r != nil {
        \t\t\t\tth = javalang.RecoverException(r)
        \t\t\t}
        \t\t}()
        \t\tres = ToAny(task.Call())
        \t}()
        \treturn &jCompletedFuture{val: res, err: th}
        }

        func futuresAllAsList(futures Ext_java_lang_Iterable) Ext_com_google_common_util_concurrent_ListenableFuture {
        \tCheckNotNull(futures)
        \tvar results []any
        \tfor it := Iterate(futures); it.HasNext(); {
        \t\tf := Cast[Ext_java_util_concurrent_Future](it.Next())
        \t\tvar v any
        \t\tvar th javalang.Throwable
        \t\tfunc() {
        \t\t\tdefer func() {
        \t\t\t\tif r := recover(); r != nil {
        \t\t\t\t\tth = javalang.RecoverException(r)
        \t\t\t\t}
        \t\t\t}()
        \t\t\tv = f.Get_0()
        \t\t}()
        \t\tif th != nil {
        \t\t\tif ee, ok := th.(Ext_java_util_concurrent_ExecutionException); ok {
        \t\t\t\treturn &jCompletedFuture{err: ee.GetCause()}
        \t\t\t}
        \t\t\treturn &jCompletedFuture{err: th}
        \t\t}
        \t\tresults = append(results, v)
        \t}
        \treturn &jCompletedFuture{val: wrapArrayList(results)}
        }

        type jAtomicIntArray struct {
        \tmu   sync.Mutex
        \tvals []int32
        }

        func (a *jAtomicIntArray) get(idx int32) int32 {
        \ta.mu.Lock()
        \tdefer a.mu.Unlock()
        \tCheckElementIndex(idx, int32(len(a.vals)))
        \treturn a.vals[idx]
        }

        func (a *jAtomicIntArray) updateAndGet(idx int32, op Ext_java_util_function_IntUnaryOperator) int32 {
        \tCheckNotNull(op)
        \ta.mu.Lock()
        \tdefer a.mu.Unlock()
        \tCheckElementIndex(idx, int32(len(a.vals)))
        \tnv := op.ApplyAsInt(a.vals[idx])
        \ta.vals[idx] = nv
        \treturn nv
        }

        func (x *JListCore) IsJavaNil_() bool            { return x == nil }
        func (x *JDequeCore) IsJavaNil_() bool           { return x == nil }
        func (x *JSetCore) IsJavaNil_() bool             { return x == nil }
        func (x *JMapCore) IsJavaNil_() bool             { return x == nil }
        func (x *jSetIter) IsJavaNil_() bool             { return x == nil }
        func (x *jComparatorFunc) IsJavaNil_() bool      { return x == nil }
        func (x *jPredFunc) IsJavaNil_() bool            { return x == nil }
        func (x *jFuncIdentity) IsJavaNil_() bool        { return x == nil }
        func (x *jFuncWrap) IsJavaNil_() bool            { return x == nil }
        func (x *jConsumerWrap) IsJavaNil_() bool        { return x == nil }
        func (x *jFullListIter) IsJavaNil_() bool        { return x == nil }
        func (x *jDescListIter) IsJavaNil_() bool        { return x == nil }
        func (x *jDequeIter) IsJavaNil_() bool           { return x == nil }
        func (x *jDequeDescIter) IsJavaNil_() bool       { return x == nil }
        func (x *jMapKeyIter) IsJavaNil_() bool          { return x == nil }
        func (x *jSubListView) IsJavaNil_() bool         { return x == nil }
        func (x *jBackedArrayList) IsJavaNil_() bool     { return x == nil }
        func (x *jReversedListView) IsJavaNil_() bool    { return x == nil }
        func (x *jTransformedListView) IsJavaNil_() bool { return x == nil }
        func (x *jSimpleEntry) IsJavaNil_() bool         { return x == nil }
        func (x *jLiveMapKeySet) IsJavaNil_() bool       { return x == nil }
        func (x *jLiveMapValues) IsJavaNil_() bool       { return x == nil }
        func (x *jMapValIter) IsJavaNil_() bool          { return x == nil }
        func (x *jLiveMapEntrySet) IsJavaNil_() bool     { return x == nil }
        func (x *jMapEntryIter) IsJavaNil_() bool        { return x == nil }
        func (x *jImmMapBuilderCore) IsJavaNil_() bool   { return x == nil }
        func (x *jAbstractListIter) IsJavaNil_() bool    { return x == nil }
        func (x *jSetViewCore) IsJavaNil_() bool         { return x == nil }
        func (x *jSyncMapCore) IsJavaNil_() bool         { return x == nil }
        func (x *jBiMapCore) IsJavaNil_() bool           { return x == nil }
        func (x *jMultimapCore) IsJavaNil_() bool        { return x == nil }
        func (x *jMultimapLiveList) IsJavaNil_() bool    { return x == nil }
        func (x *jMultimapLiveSet) IsJavaNil_() bool     { return x == nil }
        func (x *jMultimapSetIter) IsJavaNil_() bool     { return x == nil }
        func (x *jMultimapLiveKeySet) IsJavaNil_() bool  { return x == nil }
        func (x *jMultimapKeyIter) IsJavaNil_() bool     { return x == nil }
        func (x *jMultimapAsMap) IsJavaNil_() bool       { return x == nil }
        func (x *jMultimapAsMapEntries) IsJavaNil_() bool { return x == nil }
        func (x *jMultimapAsMapIter) IsJavaNil_() bool   { return x == nil }
        func (x *jMultimapBuilderCore) IsJavaNil_() bool { return x == nil }
        func (x *jMultisetCore) IsJavaNil_() bool        { return x == nil }
        func (x *jMultisetEntry) IsJavaNil_() bool       { return x == nil }
        func (x *jMultisetIter) IsJavaNil_() bool        { return x == nil }
        func (x *jTableCore) IsJavaNil_() bool           { return x == nil }
        func (x *jTableRowMap) IsJavaNil_() bool         { return x == nil }
        func (x *jTableRowKeySet) IsJavaNil_() bool      { return x == nil }
        func (x *jTableRowKeyIter) IsJavaNil_() bool     { return x == nil }
        func (x *jTableRowMapOuter) IsJavaNil_() bool    { return x == nil }
        func (x *jTableRowKeyView) IsJavaNil_() bool     { return x == nil }
        func (x *jTableRowEntrySet) IsJavaNil_() bool    { return x == nil }
        func (x *jTableRowSingleIter) IsJavaNil_() bool  { return x == nil }
        func (x *jPriorityQueue) IsJavaNil_() bool       { return x == nil }
        func (x *jFuncIterable) IsJavaNil_() bool        { return x == nil }
        func (x *jConcatIter) IsJavaNil_() bool          { return x == nil }
        func (x *jFilterIter) IsJavaNil_() bool          { return x == nil }
        func (x *jSimpleIterAdapter) IsJavaNil_() bool   { return x == nil }
        func (x *jTransformIter) IsJavaNil_() bool       { return x == nil }
        func (x *jPeekingIter) IsJavaNil_() bool         { return x == nil }
        func (x *jAbstractIterState) IsJavaNil_() bool   { return x == nil }
        func (x *jComparisonChain) IsJavaNil_() bool     { return x == nil }
        func (x *jRangeCore) IsJavaNil_() bool           { return x == nil }
        func (x *jRangeSetCore) IsJavaNil_() bool        { return x == nil }
        func (x *jStream) IsJavaNil_() bool              { return x == nil }
        func (x *jStreamBuilder) IsJavaNil_() bool       { return x == nil }
        func (x *jIntStream) IsJavaNil_() bool           { return x == nil }
        func (x *jCollector) IsJavaNil_() bool           { return x == nil }
        func (x *jCacheBuilderCore) IsJavaNil_() bool    { return x == nil }
        func (x *jFuncCacheLoader) IsJavaNil_() bool     { return x == nil }
        func (x *jLoadingCache) IsJavaNil_() bool        { return x == nil }
        func (x *jEmptyMapIter) IsJavaNil_() bool        { return x == nil }
        func (x *jLiveMapEntry) IsJavaNil_() bool        { return x == nil }
        func (x *jDirectExecutor) IsJavaNil_() bool      { return x == nil }
        func (x *jThreadPoolExecutor) IsJavaNil_() bool  { return x == nil }
        func (x *jListeningExecWrapper) IsJavaNil_() bool { return x == nil }
        func (x *jCompletedFuture) IsJavaNil_() bool     { return x == nil }
        func (x *jAsyncFuture) IsJavaNil_() bool         { return x == nil }
        func (x *jAtomicIntArray) IsJavaNil_() bool      { return x == nil }
        func (x *jEnumSetCore) IsJavaNil_() bool         { return x == nil }
        func (x *jEnumSetIter) IsJavaNil_() bool         { return x == nil }
        func (x *jEnumMapCore) IsJavaNil_() bool         { return x == nil }
        func (x *jLiveEnumMapKeySet) IsJavaNil_() bool   { return x == nil }
        func (x *jEnumMapKeyIter) IsJavaNil_() bool      { return x == nil }
        func (x *jLiveEnumMapValues) IsJavaNil_() bool   { return x == nil }
        func (x *jEnumMapValIter) IsJavaNil_() bool      { return x == nil }
        func (x *jLiveEnumMapEntrySet) IsJavaNil_() bool { return x == nil }
        func (x *jEnumMapEntry) IsJavaNil_() bool        { return x == nil }
        func (x *jEnumMapEntryIter) IsJavaNil_() bool    { return x == nil }
        """;
  }

  private static String buildColRtLoweringHelpers() {
    return """

        // --- Post-Pass Native Slice, OrderedSet, and LocalStringBuilder Runtime Helpers ---

        func boxJavaElem(v any) any {
        \tswitch x := v.(type) {
        \tcase int32:
        \t\treturn ValueOfInt(x)
        \tcase int64:
        \t\treturn ValueOfLong(x)
        \tcase bool:
        \t\treturn ValueOfBoolean(x)
        \tcase float64:
        \t\treturn ValueOfDouble(x)
        \tcase float32:
        \t\treturn ValueOfFloat(x)
        \tcase uint16:
        \t\treturn ValueOfChar(x)
        \tcase int8:
        \t\treturn ValueOfByte(x)
        \tcase int16:
        \t\treturn ValueOfShort(x)
        \tcase uint32:
        \t\treturn boxNodeById(x)
        \tdefault:
        \t\tif rn, ok := v.(jRawNodeProvider); ok {
        \t\t\tif rn.RawNodeId_() == 0 {
        \t\t\t\treturn nil
        \t\t\t}
        \t\t\treturn v
        \t\t}
        \t\treturn ToAny(v)
        \t}
        }

        func unboxJavaElem[T any](v any) T {
        \tvar zero T
        \tswitch any(zero).(type) {
        \tcase string:
        \t\treturn any(CastString(v)).(T)
        \tcase int32:
        \t\treturn any(UnboxInt(CastPtr[BoxInt](v))).(T)
        \tcase int64:
        \t\treturn any(UnboxLong(CastPtr[BoxLong](v))).(T)
        \tcase bool:
        \t\treturn any(UnboxBoolean(CastPtr[BoxBoolean](v))).(T)
        \tcase float64:
        \t\treturn any(UnboxDouble(CastPtr[BoxDouble](v))).(T)
        \tcase float32:
        \t\treturn any(UnboxFloat(CastPtr[BoxFloat](v))).(T)
        \tcase uint16:
        \t\treturn any(UnboxChar(CastPtr[BoxChar](v))).(T)
        \tcase int8:
        \t\treturn any(UnboxByte(CastPtr[BoxByte](v))).(T)
        \tcase int16:
        \t\treturn any(UnboxShort(CastPtr[BoxShort](v))).(T)
        \tcase uint32:
        \t\treturn any(CastNode[uint32](v)).(T)
        \tdefault:
        \t\treturn Cast[T](v)
        \t}
        }

        func MakeNonNilSlice[T any](cap int32) []T {
        \tif cap < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal Capacity: %d", cap)))
        \t}
        \tif cap == 0 {
        \t\treturn []T{}
        \t}
        \treturn make([]T, 0, cap)
        }

        func CheckNotNullElem[T any](v T) T {
        \tif any(v) == nil {
        \t\tThrowNPE()
        \t}
        \tswitch x := any(v).(type) {
        \tcase string:
        \t\tif x == NullString {
        \t\t\tThrowNPE()
        \t\t}
        \t\treturn v
        \tcase int32, int64, bool, float64, float32, uint16, int8, int16, uint8, int, uint64:
        \t\treturn v
        \tdefault:
        \t\tswitch nullValKindOf[T]() {
        \t\tcase kindNeverNull:
        \t\t\treturn v
        \t\tcase kindZeroIsNull:
        \t\t\tvar zero T
        \t\t\tif any(v) == any(zero) {
        \t\t\t\tThrowNPE()
        \t\t\t}
        \t\t\treturn v
        \t\tcase kindNullStringType:
        \t\t\tif any(v).(string) == NullString {
        \t\t\t\tThrowNPE()
        \t\t\t}
        \t\t\treturn v
        \t\tdefault:
        \t\t\tif isJavaNullMapValSlow(v) {
        \t\t\t\tThrowNPE()
        \t\t\t}
        \t\t\treturn v
        \t\t}
        \t}
        }

        func ImmSliceCloneNonNull[T any](s []T) []T {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif len(s) == 0 {
        \t\treturn []T{}
        \t}
        \tout := make([]T, len(s))
        \tswitch nullValKindOf[T]() {
        \tcase kindNeverNull:
        \t\tcopy(out, s)
        \t\treturn out
        \tcase kindZeroIsNull:
        \t\tvar zero T
        \t\tfor i, v := range s {
        \t\t\tif any(v) == any(zero) {
        \t\t\t\tThrowNPE()
        \t\t\t}
        \t\t\tout[i] = v
        \t\t}
        \t\treturn out
        \tcase kindNullStringType:
        \t\tfor _, sv := range any(s).([]string) {
        \t\t\tif sv == NullString {
        \t\t\t\tThrowNPE()
        \t\t\t}
        \t\t}
        \t\tcopy(out, s)
        \t\treturn out
        \tdefault:
        \t\tfor i, v := range s {
        \t\t\tif isJavaNullMapValSlow(v) {
        \t\t\t\tThrowNPE()
        \t\t\t}
        \t\t\tout[i] = v
        \t\t}
        \t\treturn out
        \t}
        }

        type immTypedSliceView[T any] struct {
        \ts []T
        }

        var _ Ext_com_google_common_collect_ImmutableList = (*immTypedSliceView[int32])(nil)

        func (l *immTypedSliceView[T]) IsJavaNil_() bool { return l == nil }
        func (*immTypedSliceView[T]) JavaClassInfo_() (string, string) {
        \treturn "com.google.common.collect.ImmutableList", "ImmutableList"
        }
        func (l *immTypedSliceView[T]) typedViewLen_() (int, bool) {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \treturn len(l.s), true
        }
        func (l *immTypedSliceView[T]) typedViewAt_(i int) any {
        \treturn boxJavaElem(l.s[i])
        }
        func (l *immTypedSliceView[T]) typedViewNodeAt_(i int) (uint32, bool) {
        \treturn typedElemRawNode(l.s[i])
        }
        func (l *immTypedSliceView[T]) IterSlice_() []any {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif len(l.s) == 0 {
        \t\treturn emptyAnySlice_
        \t}
        \tout := make([]any, len(l.s))
        \tfor i, v := range l.s {
        \t\tout[i] = boxJavaElem(v)
        \t}
        \treturn out
        }
        func (l *immTypedSliceView[T]) ListGetStr_(i int32) string {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tCheckElementIndex(i, int32(len(l.s)))
        \treturn CastString(boxJavaElem(l.s[i]))
        }
        func (l *immTypedSliceView[T]) SetContainsStr_(k string) bool {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif k == NullString {
        \t\treturn false
        \t}
        \treturn SliceContains(l.s, k)
        }
        func (l *immTypedSliceView[T]) SetAddStr_(string) bool {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }
        func (l *immTypedSliceView[T]) SetRemoveStr_(string) bool {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }
        func (l *immTypedSliceView[T]) SetContainsNode_(u uint32) bool {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif u == 0 {
        \t\treturn false
        \t}
        \tfor _, v := range l.s {
        \t\tif singletonElemMatchesNode(boxJavaElem(v), u) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }
        func (l *immTypedSliceView[T]) SetAddNode_(uint32) bool {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }
        func (l *immTypedSliceView[T]) SetRemoveNode_(uint32) bool {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }
        func (l *immTypedSliceView[T]) Equals(arg0 any) bool {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif any(l) == arg0 {
        \t\treturn true
        \t}
        \treturn ImmSliceEquals(l.s, arg0)
        }
        func (l *immTypedSliceView[T]) HashCode() int32 {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \treturn ImmSliceHashCode(l.s)
        }
        func (l *immTypedSliceView[T]) String() string {
        \tif l == nil {
        \t\treturn "null"
        \t}
        \tif len(l.s) == 0 {
        \t\treturn "[]"
        \t}
        \tvar sb strings.Builder
        \tsb.WriteByte('[')
        \tfor i, v := range l.s {
        \t\tif i > 0 {
        \t\t\tsb.WriteString(", ")
        \t\t}
        \t\tsb.WriteString(StrVal(boxJavaElem(v)))
        \t}
        \tsb.WriteByte(']')
        \treturn sb.String()
        }
        func (*immTypedSliceView[T]) Implements_j2gort_Ext_com_google_common_collect_ImmutableList()       {}
        func (*immTypedSliceView[T]) Implements_j2gort_Ext_com_google_common_collect_ImmutableCollection() {}
        func (*immTypedSliceView[T]) Implements_j2gort_Ext_java_util_List()                                {}
        func (*immTypedSliceView[T]) Implements_j2gort_Ext_java_util_RandomAccess()                        {}
        func (*immTypedSliceView[T]) Implements_j2gort_Ext_java_util_AbstractCollection()                  {}
        func (*immTypedSliceView[T]) Implements_j2gort_Ext_java_util_SequencedCollection()                 {}
        func (*immTypedSliceView[T]) Implements_j2gort_Ext_java_util_Collection()                          {}
        func (*immTypedSliceView[T]) Implements_j2gort_Ext_java_lang_Iterable()                            {}
        func (*immTypedSliceView[T]) Implements_j2gort_Ext_java_io_Serializable()                          {}
        func (*immTypedSliceView[T]) GetExt_com_google_common_collect_ImmutableListBase() *Ext_com_google_common_collect_ImmutableListBase {
        \treturn nil
        }
        func (*immTypedSliceView[T]) GetExt_com_google_common_collect_ImmutableCollectionBase() *Ext_com_google_common_collect_ImmutableCollectionBase {
        \treturn nil
        }
        func (*immTypedSliceView[T]) GetExt_java_util_AbstractCollectionBase() *Ext_java_util_AbstractCollectionBase {
        \treturn nil
        }
        func (l *immTypedSliceView[T]) Iterator() any {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \treturn SliceIterator(l.s)
        }
        func (l *immTypedSliceView[T]) ListIterator_1(arg0 int32) any {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \treturn newImmListIter(l.IterSlice_(), arg0)
        }
        func (l *immTypedSliceView[T]) ForEach(arg0 Ext_java_util_function_Consumer) {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(arg0)
        \tfor _, v := range l.s {
        \t\targ0.Accept(boxJavaElem(v))
        \t}
        }
        func (l *immTypedSliceView[T]) IndexOf_1(arg0 any) int32     { return SliceIndexOf(l.s, arg0) }
        func (l *immTypedSliceView[T]) LastIndexOf_1(arg0 any) int32 { return SliceLastIndexOf(l.s, arg0) }
        func (l *immTypedSliceView[T]) Contains(arg0 any) bool       { return SliceContains(l.s, arg0) }
        func (l *immTypedSliceView[T]) SubList(arg0 int32, arg1 int32) any {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tn := int32(len(l.s))
        \tCheckPositionIndexes(arg0, arg1, n)
        \tif arg0 == 0 && arg1 == n {
        \t\treturn l
        \t}
        \treturn ImmListFromTypedSlice(l.s[arg0:arg1:arg1])
        }
        func (l *immTypedSliceView[T]) Set(int32, any) any {
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }
        func (l *immTypedSliceView[T]) Add_2(int32, any) {
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }
        func (l *immTypedSliceView[T]) Remove_1_int32(int32) any {
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }
        func (l *immTypedSliceView[T]) Sort(Ext_java_util_Comparator) {
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }
        func (l *immTypedSliceView[T]) AsList() Ext_com_google_common_collect_ImmutableList { return l }
        func (l *immTypedSliceView[T]) Reverse() Ext_com_google_common_collect_ImmutableList {
        \tif len(l.s) <= 1 {
        \t\treturn l
        \t}
        \treturn ImmListFromTypedSlice(ImmSliceReverse(l.s))
        }
        func (l *immTypedSliceView[T]) ToArray_1_slice_any(arg0 []any) []any {
        \treturn sliceToArray1(l.IterSlice_(), arg0)
        }
        func (l *immTypedSliceView[T]) Add_1(any) bool {
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }
        func (l *immTypedSliceView[T]) Remove_1_any(any) bool {
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }
        func (l *immTypedSliceView[T]) AddAll_1(Ext_java_util_Collection) bool {
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }
        func (l *immTypedSliceView[T]) RemoveAll(Ext_java_util_Collection) bool {
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }
        func (l *immTypedSliceView[T]) Clear() {
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }
        func (l *immTypedSliceView[T]) Size() int32   { return int32(len(l.s)) }
        func (l *immTypedSliceView[T]) IsEmpty() bool { return len(l.s) == 0 }
        func (l *immTypedSliceView[T]) ContainsAll(arg0 Ext_java_util_Collection) bool {
        \treturn colContainsAll(func(x any) bool { return SliceContains(l.s, x) }, arg0)
        }
        func (l *immTypedSliceView[T]) Get(arg0 int32) any {
        \tCheckElementIndex(arg0, int32(len(l.s)))
        \treturn boxJavaElem(l.s[arg0])
        }
        func (l *immTypedSliceView[T]) AddFirst(any) {
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }
        func (l *immTypedSliceView[T]) AddLast(any) {
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }
        func (l *immTypedSliceView[T]) GetFirst() any {
        \tif len(l.s) == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \treturn boxJavaElem(l.s[0])
        }
        func (l *immTypedSliceView[T]) GetLast() any {
        \tif len(l.s) == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \treturn boxJavaElem(l.s[len(l.s)-1])
        }
        func (l *immTypedSliceView[T]) RemoveFirst() any {
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }
        func (l *immTypedSliceView[T]) RemoveLast() any {
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }
        func (l *immTypedSliceView[T]) Stream() Ext_java_util_stream_Stream {
        \treturn streamFromSlice(l.IterSlice_())
        }

        func emptyNonNilSlice[T any]() []T {
        \treturn make([]T, 0)
        }

        func ImmListFromTypedSlice[T any](s []T) Ext_com_google_common_collect_ImmutableList {
        \tif s == nil {
        \t\treturn nil
        \t}
        \tif len(s) == 0 {
        \t\treturn emptyImmutableList_
        \t}
        \tif ss, ok := any(s).([]string); ok {
        \t\treturn &jImmStrList{strs: ss[:len(ss):len(ss)]}
        \t}
        \treturn &immTypedSliceView[T]{s: s[:len(s):len(s)]}
        }

        func ImmListFromCompSlice[T comparable](s []T) Ext_com_google_common_collect_ImmutableList {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif len(s) == 0 {
        \t\treturn emptyImmutableList_
        \t}
        \tvar zero T
        \tcs := make([]T, len(s))
        \tfor i, v := range s {
        \t\tif v == zero {
        \t\t\tThrowNPE()
        \t\t}
        \t\tcs[i] = v
        \t}
        \treturn &immTypedSliceView[T]{s: cs[:len(cs):len(cs)]}
        }

        func ImmListFromSlice[T any](s []T) Ext_com_google_common_collect_ImmutableList {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif len(s) == 0 {
        \t\treturn emptyImmutableList_
        \t}
        \tif ss, ok := any(s).([]string); ok {
        \t\treturn wrapImmStrSliceCopy(ss)
        \t}
        \tcs := ImmSliceCloneNonNull(s)
        \treturn &immTypedSliceView[T]{s: cs[:len(cs):len(cs)]}
        }

        func unboxNonNullImmElem[T any](v any) T {
        \tif v == nil {
        \t\tThrowNPE()
        \t}
        \tif res, ok := v.(T); ok {
        \t\treturn res
        \t}
        \treturn unboxJavaElem[T](CheckNotNull(v))
        }

        func ImmListCopyToSlice[T any](src any) []T {
        \tCheckNotNull(src)
        \tif tv, ok := src.(*immTypedSliceView[T]); ok {
        \t\tif len(tv.s) == 0 {
        \t\t\treturn emptyNonNilSlice[T]()
        \t\t}
        \t\treturn tv.s[:len(tv.s):len(tv.s)]
        \t}
        \tif al, ok := src.(*jTypedArrayList[T]); ok && al != nil && al.lc == nil {
        \t\tif len(al.s) == 0 {
        \t\t\treturn emptyNonNilSlice[T]()
        \t\t}
        \t\treturn ImmSliceCloneNonNull(al.s)
        \t}
        \tif src == emptyImmutableList_ {
        \t\treturn emptyNonNilSlice[T]()
        \t}
        \tif se, ok := src.(jSingleNonNullElem); ok {
        \t\treturn []T{unboxNonNullImmElem[T](se.singleNonNullElem_())}
        \t}
        \tswitch l := src.(type) {
        \tcase *jImmList1:
        \t\tif l == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\treturn []T{unboxNonNullImmElem[T](l.elem[0])}
        \tcase *jImmList2:
        \t\tif l == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\treturn []T{unboxNonNullImmElem[T](l.elem[0]), unboxNonNullImmElem[T](l.elem[1])}
        \tcase *jSingletonList:
        \t\tif l == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\treturn []T{unboxJavaElem[T](CheckNotNull(l.elem[0]))}
        \t}
        \tif ssi, ok := src.(StrSliceIterable); ok {
        \t\tif ss, ok2 := ssi.IterStrSlice_(); ok2 {
        \t\t\tif len(ss) == 0 {
        \t\t\t\treturn emptyNonNilSlice[T]()
        \t\t\t}
        \t\t\tif _, isStr := any((*T)(nil)).(*string); isStr {
        \t\t\t\tif l, ok3 := src.(*jImmStrList); ok3 && l != nil {
        \t\t\t\t\treturn any(l.strs[:len(l.strs):len(l.strs)]).([]T)
        \t\t\t\t}
        \t\t\t\tout := make([]string, len(ss))
        \t\t\t\tfor i, str := range ss {
        \t\t\t\t\tif str == NullString {
        \t\t\t\t\t\tThrowNPE()
        \t\t\t\t\t}
        \t\t\t\t\tout[i] = str
        \t\t\t\t}
        \t\t\t\treturn any(out).([]T)
        \t\t\t}
        \t\t\tout := make([]T, len(ss))
        \t\t\tfor i, str := range ss {
        \t\t\t\tif str == NullString {
        \t\t\t\t\tThrowNPE()
        \t\t\t\t}
        \t\t\t\tout[i] = unboxJavaElem[T](str)
        \t\t\t}
        \t\t\treturn out
        \t\t}
        \t}
        \tif tv, ok := src.(jTypedViewCursorSource); ok {
        \t\tif n, ok2 := tv.typedViewLen_(); ok2 {
        \t\t\tif n == 0 {
        \t\t\t\treturn emptyNonNilSlice[T]()
        \t\t\t}
        \t\t\tout := make([]T, n)
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\tout[i] = unboxJavaElem[T](CheckNotNull(tv.typedViewAt_(i)))
        \t\t\t}
        \t\t\treturn out
        \t\t}
        \t}
        \tif si, ok := src.(SliceIterable); ok {
        \t\titems := si.IterSlice_()
        \t\tif len(items) == 0 {
        \t\t\treturn emptyNonNilSlice[T]()
        \t\t}
        \t\tout := make([]T, len(items))
        \t\tfor i, raw := range items {
        \t\t\tout[i] = unboxJavaElem[T](CheckNotNull(raw))
        \t\t}
        \t\treturn out
        \t}
        \tif c, ok := src.(Ext_java_util_Collection); ok {
        \t\tsz := int(c.Size())
        \t\tif sz == 0 {
        \t\t\treturn emptyNonNilSlice[T]()
        \t\t}
        \t\tout := make([]T, 0, sz)
        \t\tfor cur := NewIterCursor(src); cur.HasNext(); {
        \t\t\tout = append(out, unboxJavaElem[T](CheckNotNull(cur.Next())))
        \t\t}
        \t\treturn out
        \t}
        \tout := emptyNonNilSlice[T]()
        \tfor cur := NewIterCursor(src); cur.HasNext(); {
        \t\tout = append(out, unboxJavaElem[T](CheckNotNull(cur.Next())))
        \t}
        \treturn out
        }

        func ImmListCopyArrayToSlice[T any](arr []any) []T {
        \tCheckNotNullSlice(arr)
        \tif len(arr) == 0 {
        \t\treturn emptyNonNilSlice[T]()
        \t}
        \tout := make([]T, len(arr))
        \tfor i, raw := range arr {
        \t\tout[i] = unboxJavaElem[T](CheckNotNull(raw))
        \t}
        \treturn out
        }

        func ImmListToSliceNullable[T any](imm Ext_com_google_common_collect_ImmutableList) []T {
        \tif imm == nil {
        \t\treturn nil
        \t}
        \tif tv, ok := any(imm).(*immTypedSliceView[T]); ok {
        \t\tif tv == nil {
        \t\t\treturn nil
        \t\t}
        \t\tif len(tv.s) == 0 {
        \t\t\treturn emptyNonNilSlice[T]()
        \t\t}
        \t\treturn tv.s[:len(tv.s):len(tv.s)]
        \t}
        \tif imm == emptyImmutableList_ {
        \t\treturn emptyNonNilSlice[T]()
        \t}
        \tif isJavaNilAny(imm) {
        \t\treturn nil
        \t}
        \tif imm.IsEmpty() {
        \t\treturn emptyNonNilSlice[T]()
        \t}
        \treturn ImmListCopyToSlice[T](imm)
        }

        func ImmListOfVarargsToSlice[T any](prefix []T, rest []any) []T {
        \tCheckNotNullSlice(rest)
        \tn := len(prefix) + len(rest)
        \tif n == 0 {
        \t\treturn []T{}
        \t}
        \tout := make([]T, n)
        \tif len(prefix) > 0 {
        \t\tswitch nullValKindOf[T]() {
        \t\tcase kindNeverNull:
        \t\t\tcopy(out, prefix)
        \t\tcase kindZeroIsNull:
        \t\t\tvar zero T
        \t\t\tfor i, v := range prefix {
        \t\t\t\tif any(v) == any(zero) {
        \t\t\t\t\tThrowNPE()
        \t\t\t\t}
        \t\t\t\tout[i] = v
        \t\t\t}
        \t\tcase kindNullStringType:
        \t\t\tfor _, sv := range any(prefix).([]string) {
        \t\t\t\tif sv == NullString {
        \t\t\t\t\tThrowNPE()
        \t\t\t\t}
        \t\t\t}
        \t\t\tcopy(out, prefix)
        \t\tdefault:
        \t\t\tfor i, v := range prefix {
        \t\t\t\tif isJavaNullMapValSlow(v) {
        \t\t\t\t\tThrowNPE()
        \t\t\t\t}
        \t\t\t\tout[i] = v
        \t\t\t}
        \t\t}
        \t}
        \tfor i, raw := range rest {
        \t\tout[len(prefix)+i] = unboxJavaElem[T](CheckNotNull(raw))
        \t}
        \treturn out
        }

        func SliceBuilderAdd[T any](b *[]T, v T) *[]T {
        \tif b == nil || *b == nil {
        \t\tThrowNPE()
        \t}
        \tif isJavaNullMapVal(v) {
        \t\tThrowNPE()
        \t}
        \t*b = append(*b, v)
        \treturn b
        }

        func AppendImmCollection[T any](dst []T, src any) []T {
        \tif dst == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(src)
        \tif tv, ok := src.(*immTypedSliceView[T]); ok {
        \t\treturn append(dst, tv.s...)
        \t}
        \tif al, ok := src.(*jTypedArrayList[T]); ok && al != nil && al.lc == nil {
        \t\tif len(al.s) == 0 {
        \t\t\treturn dst
        \t\t}
        \t\treturn append(dst, ImmSliceCloneNonNull(al.s)...)
        \t}
        \tif ssi, ok := src.(StrSliceIterable); ok {
        \t\tif ss, ok2 := ssi.IterStrSlice_(); ok2 {
        \t\t\tif dstStr, ok3 := any(&dst).(*[]string); ok3 {
        \t\t\t\tfor _, str := range ss {
        \t\t\t\t\tif str == NullString {
        \t\t\t\t\t\tThrowNPE()
        \t\t\t\t\t}
        \t\t\t\t}
        \t\t\t\t*dstStr = append(*dstStr, ss...)
        \t\t\t\treturn dst
        \t\t\t}
        \t\t\tfor _, str := range ss {
        \t\t\t\tif str == NullString {
        \t\t\t\t\tThrowNPE()
        \t\t\t\t}
        \t\t\t\tdst = append(dst, unboxJavaElem[T](str))
        \t\t\t}
        \t\t\treturn dst
        \t\t}
        \t}
        \tif si, ok := src.(SliceIterable); ok {
        \t\tfor _, raw := range si.IterSlice_() {
        \t\t\tdst = append(dst, unboxJavaElem[T](CheckNotNull(raw)))
        \t\t}
        \t\treturn dst
        \t}
        \tfor cur := NewIterCursor(src); cur.HasNext(); {
        \t\tdst = append(dst, unboxJavaElem[T](CheckNotNull(cur.Next())))
        \t}
        \treturn dst
        }

        func SliceBuilderAddAll[T any](b *[]T, src any) *[]T {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \t*b = AppendImmCollection(*b, src)
        \treturn b
        }

        func SliceBuilderAddAllSlice[T any](b *[]T, src []T) *[]T {
        \tif b == nil || *b == nil || src == nil {
        \t\tThrowNPE()
        \t}
        \t*b = append(*b, src...)
        \treturn b
        }

        func ImmListCopyNodeToSlice[N ~uint32](src any) []N {
        \tCheckNotNull(src)
        \tif tv, ok := src.(*immTypedSliceView[N]); ok {
        \t\tif len(tv.s) == 0 {
        \t\t\treturn emptyNonNilSlice[N]()
        \t\t}
        \t\treturn tv.s[:len(tv.s):len(tv.s)]
        \t}
        \tif src == emptyImmutableList_ {
        \t\treturn emptyNonNilSlice[N]()
        \t}
        \tif si, ok := src.(SliceIterable); ok {
        \t\titems := si.IterSlice_()
        \t\tif len(items) == 0 {
        \t\t\treturn emptyNonNilSlice[N]()
        \t\t}
        \t\tout := make([]N, len(items))
        \t\tfor i, raw := range items {
        \t\t\tout[i] = CastNode[N](CheckNotNull(raw))
        \t\t}
        \t\treturn out
        \t}
        \tout := emptyNonNilSlice[N]()
        \tfor cur := NewIterCursor(src); cur.HasNext(); {
        \t\tout = append(out, N(cur.NextNode()))
        \t}
        \treturn out
        }

        func ImmListToNodeSliceNullable[N ~uint32](imm Ext_com_google_common_collect_ImmutableList) []N {
        \tif imm == nil {
        \t\treturn nil
        \t}
        \tif tv, ok := any(imm).(*immTypedSliceView[N]); ok {
        \t\tif tv == nil {
        \t\t\treturn nil
        \t\t}
        \t\tif len(tv.s) == 0 {
        \t\t\treturn emptyNonNilSlice[N]()
        \t\t}
        \t\treturn tv.s[:len(tv.s):len(tv.s)]
        \t}
        \tif imm == emptyImmutableList_ {
        \t\treturn emptyNonNilSlice[N]()
        \t}
        \tif isJavaNilAny(imm) {
        \t\treturn nil
        \t}
        \tif imm.IsEmpty() {
        \t\treturn emptyNonNilSlice[N]()
        \t}
        \treturn ImmListCopyNodeToSlice[N](imm)
        }

        func AppendImmNodeCollection[N ~uint32](dst []N, src any) []N {
        \tif dst == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(src)
        \tif tv, ok := src.(*immTypedSliceView[N]); ok {
        \t\treturn append(dst, tv.s...)
        \t}
        \tfor cur := NewIterCursor(src); cur.HasNext(); {
        \t\tdst = append(dst, N(cur.NextNode()))
        \t}
        \treturn dst
        }

        func SliceBuilderAddAllNode[N ~uint32](b *[]N, src any) *[]N {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \t*b = AppendImmNodeCollection(*b, src)
        \treturn b
        }

        func ImmSliceSize[T any](s []T) int32 {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn int32(len(s))
        }

        func ImmSliceIsEmpty[T any](s []T) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn len(s) == 0
        }

        func ImmSliceGet[T any](s []T, idx int32) T {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tCheckElementIndex(idx, int32(len(s)))
        \treturn s[idx]
        }

        func ListGetVal[T any](l any, i int32) T {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif tv, ok := l.(*immTypedSliceView[T]); ok {
        \t\tif tv == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\tCheckElementIndex(i, int32(len(tv.s)))
        \t\treturn tv.s[i]
        \t}
        \tif al, ok := l.(*jTypedArrayList[T]); ok {
        \t\tif al == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\tif al.lc != nil {
        \t\t\treturn unboxJavaElem[T](al.lc.Get(i))
        \t\t}
        \t\tCheckElementIndex(i, int32(len(al.s)))
        \t\treturn al.s[i]
        \t}
        \tif isJavaNilAny(l) {
        \t\tThrowNPE()
        \t}
        \treturn unboxJavaElem[T](l.(Ext_java_util_List).Get(i))
        }

        type TypedIterCursor[T any] struct {
        \ts    []T
        \tidx  int
        \trest JIterCursor
        \traw  bool
        }

        func NewTypedIterCursor[T any](coll any) TypedIterCursor[T] {
        \tif coll == nil {
        \t\tThrowNPE()
        \t}
        \tif tv, ok := coll.(*immTypedSliceView[T]); ok {
        \t\tif tv == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\treturn TypedIterCursor[T]{s: tv.s}
        \t}
        \tif al, ok := coll.(*jTypedArrayList[T]); ok {
        \t\tif al == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\tif al.lc == nil {
        \t\t\treturn TypedIterCursor[T]{s: al.s}
        \t\t}
        \t}
        \tcur := NewIterCursor(coll)
        \tif !cur.HasNext() {
        \t\treturn TypedIterCursor[T]{}
        \t}
        \treturn TypedIterCursor[T]{rest: cur, raw: true}
        }

        func (c *TypedIterCursor[T]) HasNext() bool {
        \tif !c.raw {
        \t\treturn c.idx < len(c.s)
        \t}
        \treturn c.rest.HasNext()
        }

        func (c *TypedIterCursor[T]) Next() T {
        \tif !c.raw {
        \t\tv := c.s[c.idx]
        \t\tc.idx++
        \t\treturn v
        \t}
        \treturn unboxJavaElem[T](c.rest.Next())
        }

        func SliceContains[T any](s []T, target any) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif target == nil || isJavaNilAny(target) {
        \t\treturn false
        \t}
        \tt := ToAny(target)
        \tif ss, ok := any(s).([]string); ok {
        \t\tts, ok2 := t.(string)
        \t\tif !ok2 || ts == NullString {
        \t\t\treturn false
        \t\t}
        \t\tfor _, x := range ss {
        \t\t\tif x == ts {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \t}
        \tfor _, v := range s {
        \t\tif javalang.ObjectsEquals(boxJavaElem(v), t) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func SliceIndexOf[T any](s []T, target any) int32 {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif target == nil || isJavaNilAny(target) {
        \t\treturn -1
        \t}
        \tt := ToAny(target)
        \tfor i, v := range s {
        \t\tif javalang.ObjectsEquals(boxJavaElem(v), t) {
        \t\t\treturn int32(i)
        \t\t}
        \t}
        \treturn -1
        }

        func SliceLastIndexOf[T any](s []T, target any) int32 {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif target == nil || isJavaNilAny(target) {
        \t\treturn -1
        \t}
        \tt := ToAny(target)
        \tfor i := len(s) - 1; i >= 0; i-- {
        \t\tif javalang.ObjectsEquals(boxJavaElem(s[i]), t) {
        \t\t\treturn int32(i)
        \t\t}
        \t}
        \treturn -1
        }

        func ImmSliceSubList[T any](s []T, fromIdx, toIdx int32) []T {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tCheckPositionIndexes(fromIdx, toIdx, int32(len(s)))
        \tif fromIdx == toIdx {
        \t\treturn []T{}
        \t}
        \treturn s[fromIdx:toIdx]
        }

        func ImmSliceReverse[T any](s []T) []T {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif len(s) == 0 {
        \t\treturn []T{}
        \t}
        \tif len(s) == 1 {
        \t\treturn s
        \t}
        \trev := make([]T, len(s))
        \tfor i, v := range s {
        \t\trev[len(s)-1-i] = v
        \t}
        \treturn rev
        }

        func ImmSliceEqualsTyped[T any](a, b []T) bool {
        \tif a == nil {
        \t\tThrowNPE()
        \t}
        \tif b == nil || len(a) != len(b) {
        \t\treturn false
        \t}
        \tfor i := range a {
        \t\tif !javaMapValsEqual(a[i], b[i], false) {
        \t\t\treturn false
        \t\t}
        \t}
        \treturn true
        }

        func ImmSliceEquals[T any](s []T, other any) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif other == nil || isJavaNilAny(other) {
        \t\treturn false
        \t}
        \tif tv, ok := other.(*immTypedSliceView[T]); ok {
        \t\treturn ImmSliceEqualsTyped(s, tv.s)
        \t}
        \tif al, ok := other.(*jTypedArrayList[T]); ok && al != nil && al.lc == nil {
        \t\treturn ImmSliceEqualsTyped(s, al.s)
        \t}
        \tlst, ok := other.(Ext_java_util_List)
        \tif !ok {
        \t\treturn false
        \t}
        \tif int32(len(s)) != lst.Size() {
        \t\treturn false
        \t}
        \tif _, isRA := other.(Ext_java_util_RandomAccess); isRA {
        \t\tfor i, v := range s {
        \t\t\tif !javalang.ObjectsEquals(boxJavaElem(v), lst.Get(int32(i))) {
        \t\t\t\treturn false
        \t\t\t}
        \t\t}
        \t\treturn true
        \t}
        \tit := Cast[Ext_java_util_Iterator](lst.Iterator())
        \tfor _, v := range s {
        \t\tif !it.HasNext() || !javalang.ObjectsEquals(boxJavaElem(v), it.Next()) {
        \t\t\treturn false
        \t\t}
        \t}
        \treturn !it.HasNext()
        }

        func ImmSliceHashCode[T any](s []T) int32 {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tvar h int32 = 1
        \tfor _, v := range s {
        \t\th = 31*h + javalang.ObjectHashCode(boxJavaElem(v))
        \t}
        \treturn h
        }

        func SliceUnmodifiableList[T any](s []T) Ext_java_util_List {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif len(s) == 0 {
        \t\treturn emptyImmutableList_
        \t}
        \tcp := make([]T, len(s))
        \tcopy(cp, s)
        \thasNull := false
        \tswitch nullValKindOf[T]() {
        \tcase kindNeverNull:
        \tcase kindNullStringType:
        \t\tfor _, v := range any(cp).([]string) {
        \t\t\tif v == NullString {
        \t\t\t\thasNull = true
        \t\t\t\tbreak
        \t\t\t}
        \t\t}
        \tcase kindZeroIsNull:
        \t\tvar zero T
        \t\tfor _, v := range cp {
        \t\t\tif any(v) == any(zero) {
        \t\t\t\thasNull = true
        \t\t\t\tbreak
        \t\t\t}
        \t\t}
        \tdefault:
        \t\tfor _, v := range cp {
        \t\t\tif isJavaNullMapValSlow(v) {
        \t\t\t\thasNull = true
        \t\t\t\tbreak
        \t\t\t}
        \t\t}
        \t}
        \tif hasNull {
        \t\tv := &unmodSlicePtrView[T]{p: &cp}
        \t\tv.BindExtSelf_(v)
        \t\treturn v
        \t}
        \treturn &immTypedSliceView[T]{s: cp[:len(cp):len(cp)]}
        }

        func SliceToJavaList[T any](s []T) Ext_java_util_List {
        \tif s == nil {
        \t\treturn nil
        \t}
        \titems := make([]any, len(s))
        \tfor i, v := range s {
        \t\titems[i] = boxJavaElem(v)
        \t}
        \treturn wrapArrayListOwn(items)
        }

        type jEmptyImmSliceIter struct{}

        func (*jEmptyImmSliceIter) IsJavaNil_() bool                                                     { return false }
        func (*jEmptyImmSliceIter) Implements_j2gort_Ext_java_util_Iterator()                            {}
        func (*jEmptyImmSliceIter) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator() {}
        func (*jEmptyImmSliceIter) GetExt_com_google_common_collect_UnmodifiableIteratorBase() *Ext_com_google_common_collect_UnmodifiableIteratorBase {
        \treturn nil
        }
        func (*jEmptyImmSliceIter) HasNext() bool { return false }
        func (*jEmptyImmSliceIter) Next() any {
        \tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        }
        func (*jEmptyImmSliceIter) Remove() {
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }
        func (*jEmptyImmSliceIter) ForEachRemaining(action Ext_java_util_function_Consumer) {
        \tCheckNotNull(action)
        }
        func (*jEmptyImmSliceIter) typedViewLen_() (int, bool)       { return 0, true }
        func (*jEmptyImmSliceIter) typedViewAt_(int) any             { return nil }
        func (*jEmptyImmSliceIter) typedViewNodeAt_(int) (uint32, bool) { return 0, false }

        var emptyImmSliceIter_ Ext_com_google_common_collect_UnmodifiableIterator = &jEmptyImmSliceIter{}

        type immTypedSliceIter[T any] struct {
        \ts   []T
        \tpos int
        }

        var _ Ext_com_google_common_collect_UnmodifiableIterator = (*immTypedSliceIter[int32])(nil)

        func (it *immTypedSliceIter[T]) IsJavaNil_() bool                                                     { return it == nil }
        func (*immTypedSliceIter[T]) Implements_j2gort_Ext_java_util_Iterator()                               {}
        func (*immTypedSliceIter[T]) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator()   {}
        func (*immTypedSliceIter[T]) GetExt_com_google_common_collect_UnmodifiableIteratorBase() *Ext_com_google_common_collect_UnmodifiableIteratorBase {
        \treturn nil
        }
        func (it *immTypedSliceIter[T]) HasNext() bool {
        \treturn it.pos < len(it.s)
        }
        func (it *immTypedSliceIter[T]) Next() any {
        \tif it.pos >= len(it.s) {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tv := it.s[it.pos]
        \tit.pos++
        \treturn boxJavaElem(v)
        }
        func (it *immTypedSliceIter[T]) NextTyped() T {
        \tif it.pos >= len(it.s) {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tv := it.s[it.pos]
        \tit.pos++
        \treturn v
        }
        func (it *immTypedSliceIter[T]) Remove() {
        \tpanic(javalang.NewUnsupportedOperationException("immutable"))
        }

        func SliceIterator[T any](s []T) Ext_com_google_common_collect_UnmodifiableIterator {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif len(s) == 0 {
        \t\treturn emptyImmSliceIter_
        \t}
        \treturn &immTypedSliceIter[T]{s: s}
        }

        type TypedSliceIter[E any] struct {
        \ts   []E
        \tidx int
        }

        var _ Ext_java_util_Iterator = (*TypedSliceIter[int32])(nil)

        func (it *TypedSliceIter[E]) IsJavaNil_() bool                                                     { return it == nil }
        func (*TypedSliceIter[E]) Implements_j2gort_Ext_java_util_Iterator()                               {}
        func (*TypedSliceIter[E]) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator()   {}
        func (*TypedSliceIter[E]) GetExt_com_google_common_collect_UnmodifiableIteratorBase() *Ext_com_google_common_collect_UnmodifiableIteratorBase {
        \treturn nil
        }
        func (it *TypedSliceIter[E]) HasNext() bool {
        \treturn it != nil && it.idx < len(it.s)
        }
        func (it *TypedSliceIter[E]) Next() any {
        \tif it == nil || it.idx >= len(it.s) {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tv := it.s[it.idx]
        \tit.idx++
        \treturn boxJavaElem(v)
        }
        func (it *TypedSliceIter[E]) NextTyped() E {
        \tif it == nil || it.idx >= len(it.s) {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tv := it.s[it.idx]
        \tit.idx++
        \treturn v
        }
        func (it *TypedSliceIter[E]) Remove() {
        \tpanic(javalang.NewUnsupportedOperationException("remove"))
        }
        func (it *TypedSliceIter[E]) ForEachRemaining(action Ext_java_util_function_Consumer) {
        \tCheckNotNull(action)
        \tif it == nil {
        \t\treturn
        \t}
        \tfor it.idx < len(it.s) {
        \t\tv := it.s[it.idx]
        \t\tit.idx++
        \t\taction.Accept(boxJavaElem(v))
        \t}
        }
        func (it *TypedSliceIter[E]) typedViewLen_() (int, bool) {
        \tif it == nil {
        \t\treturn 0, true
        \t}
        \treturn len(it.s) - it.idx, true
        }
        func (it *TypedSliceIter[E]) typedViewAt_(i int) any {
        \treturn boxJavaElem(it.s[it.idx+i])
        }
        func (it *TypedSliceIter[E]) typedViewNodeAt_(i int) (uint32, bool) {
        \treturn typedElemRawNode(it.s[it.idx+i])
        }
        func (it *TypedSliceIter[E]) IterStrSlice_() ([]string, bool) {
        \tif it == nil {
        \t\treturn nil, false
        \t}
        \tif ss, ok := any(it.s).([]string); ok {
        \t\treturn ss[it.idx:], true
        \t}
        \treturn nil, false
        }

        func NewTypedSliceIter[E any](s []E) Ext_java_util_Iterator {
        \tif len(s) == 0 {
        \t\treturn emptyImmSliceIter_
        \t}
        \treturn &TypedSliceIter[E]{s: s}
        }

        func SliceGetFirst[T any](s []T) T {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif len(s) == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \treturn s[0]
        }

        func SliceGetLast[T any](s []T) T {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif len(s) == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \treturn s[len(s)-1]
        }

        func SliceRemoveAt[T any](s []T, idx int32) []T {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tCheckElementIndex(idx, int32(len(s)))
        \ti := int(idx)
        \tvar zero T
        \tcopy(s[i:], s[i+1:])
        \ts[len(s)-1] = zero
        \treturn s[:len(s)-1]
        }

        func SliceRemoveAtElem[T any](sp *[]T, idx int32) T {
        \ts := *sp
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tCheckElementIndex(idx, int32(len(s)))
        \ti := int(idx)
        \tremoved := s[i]
        \tvar zero T
        \tif i+1 < len(s) {
        \t\tcopy(s[i:], s[i+1:])
        \t}
        \ts[len(s)-1] = zero
        \t*sp = s[:len(s)-1]
        \treturn removed
        }

        func AppendStrCol(dst []string, col any) []string {
        \tif dst == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(col)
        \tif ssi, ok := col.(StrSliceIterable); ok {
        \t\tif ss, ok2 := ssi.IterStrSlice_(); ok2 {
        \t\t\treturn append(dst, ss...)
        \t\t}
        \t}
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\tdst = append(dst, cur.NextStr())
        \t}
        \treturn dst
        }

        func AppendNodeCol[N ~uint32](dst []N, col any) []N {
        \tif dst == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(col)
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\tdst = append(dst, N(cur.NextNode()))
        \t}
        \treturn dst
        }

        func AppendCollection[T any](dst []T, col any) []T {
        \tif dst == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(col)
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\tdst = append(dst, unboxJavaElem[T](cur.Next()))
        \t}
        \treturn dst
        }

        func SliceInsertAt[T any](s []T, idx int32, elem T) []T {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tCheckPositionIndex(idx, int32(len(s)))
        \ti := int(idx)
        \tif i == len(s) {
        \t\treturn append(s, elem)
        \t}
        \tvar zero T
        \ts = append(s, zero)
        \tcopy(s[i+1:], s[i:len(s)-1])
        \ts[i] = elem
        \treturn s
        }

        func SliceReverseInPlace[T any](s []T) []T {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tfor i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
        \t\ts[i], s[j] = s[j], s[i]
        \t}
        \treturn s
        }

        func CloneNonNilSlice[T any](s []T) []T {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif len(s) == 0 {
        \t\treturn []T{}
        \t}
        \tout := make([]T, len(s))
        \tcopy(out, s)
        \treturn out
        }

        func CopyStrColToSlice(col any) []string {
        \tCheckNotNull(col)
        \tif ssi, ok := col.(StrSliceIterable); ok {
        \t\tif ss, ok2 := ssi.IterStrSlice_(); ok2 {
        \t\t\tif len(ss) == 0 {
        \t\t\t\treturn []string{}
        \t\t\t}
        \t\t\tout := make([]string, len(ss))
        \t\t\tcopy(out, ss)
        \t\t\treturn out
        \t\t}
        \t}
        \tif tv, ok := col.(jTypedViewCursorSource); ok {
        \t\tif n, ok2 := tv.typedViewLen_(); ok2 {
        \t\t\tif n == 0 {
        \t\t\t\treturn []string{}
        \t\t\t}
        \t\t\tout := make([]string, n)
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\tout[i] = CastString(tv.typedViewAt_(i))
        \t\t\t}
        \t\t\treturn out
        \t\t}
        \t}
        \tif si, ok := col.(SliceIterable); ok {
        \t\titems := si.IterSlice_()
        \t\tif len(items) == 0 {
        \t\t\treturn []string{}
        \t\t}
        \t\tout := make([]string, len(items))
        \t\tfor i, raw := range items {
        \t\t\tout[i] = CastString(raw)
        \t\t}
        \t\treturn out
        \t}
        \tif c, ok := col.(Ext_java_util_Collection); ok {
        \t\tsz := int(c.Size())
        \t\tif sz == 0 {
        \t\t\treturn []string{}
        \t\t}
        \t\tout := make([]string, 0, sz)
        \t\tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\t\tout = append(out, cur.NextStr())
        \t\t}
        \t\treturn out
        \t}
        \tout := []string{}
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\tout = append(out, cur.NextStr())
        \t}
        \treturn out
        }

        func CopyNodeColToSlice[N ~uint32](col any) []N {
        \tCheckNotNull(col)
        \tif tv, ok := col.(jTypedViewCursorSource); ok {
        \t\tif n, ok2 := tv.typedViewLen_(); ok2 {
        \t\t\tif n == 0 {
        \t\t\t\treturn []N{}
        \t\t\t}
        \t\t\tout := make([]N, n)
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\tif u, ok3 := tv.typedViewNodeAt_(i); ok3 {
        \t\t\t\t\tout[i] = N(u)
        \t\t\t\t} else {
        \t\t\t\t\tout[i] = CastNode[N](tv.typedViewAt_(i))
        \t\t\t\t}
        \t\t\t}
        \t\t\treturn out
        \t\t}
        \t}
        \tif si, ok := col.(SliceIterable); ok {
        \t\titems := si.IterSlice_()
        \t\tif len(items) == 0 {
        \t\t\treturn []N{}
        \t\t}
        \t\tout := make([]N, len(items))
        \t\tfor i, raw := range items {
        \t\t\tout[i] = CastNode[N](raw)
        \t\t}
        \t\treturn out
        \t}
        \tif c, ok := col.(Ext_java_util_Collection); ok {
        \t\tsz := int(c.Size())
        \t\tif sz == 0 {
        \t\t\treturn []N{}
        \t\t}
        \t\tout := make([]N, 0, sz)
        \t\tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\t\tout = append(out, N(cur.NextNode()))
        \t\t}
        \t\treturn out
        \t}
        \tout := []N{}
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\tout = append(out, N(cur.NextNode()))
        \t}
        \treturn out
        }

        func CopyCollectionToSlice[T any](col any) []T {
        \tCheckNotNull(col)
        \tif tv, ok := col.(*immTypedSliceView[T]); ok {
        \t\treturn CloneNonNilSlice(tv.s)
        \t}
        \tif al, ok := col.(*jTypedArrayList[T]); ok && al != nil && al.lc == nil {
        \t\treturn CloneNonNilSlice(al.s)
        \t}
        \tif se, ok := col.(jSingleNonNullElem); ok {
        \t\treturn []T{unboxNonNullImmElem[T](se.singleNonNullElem_())}
        \t}
        \tswitch l := col.(type) {
        \tcase *jImmList1:
        \t\tif l == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\treturn []T{unboxNonNullImmElem[T](l.elem[0])}
        \tcase *jImmList2:
        \t\tif l == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\treturn []T{unboxNonNullImmElem[T](l.elem[0]), unboxNonNullImmElem[T](l.elem[1])}
        \tcase *jSingletonList:
        \t\tif l == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\treturn []T{unboxJavaElem[T](l.elem[0])}
        \t}
        \tif ssi, ok := col.(StrSliceIterable); ok {
        \t\tif ss, ok2 := ssi.IterStrSlice_(); ok2 {
        \t\t\tif len(ss) == 0 {
        \t\t\t\treturn []T{}
        \t\t\t}
        \t\t\tout := make([]T, len(ss))
        \t\t\tif outStr, ok3 := any(&out).(*[]string); ok3 {
        \t\t\t\tcopy(*outStr, ss)
        \t\t\t\treturn out
        \t\t\t}
        \t\t\tfor i, str := range ss {
        \t\t\t\tout[i] = unboxJavaElem[T](str)
        \t\t\t}
        \t\t\treturn out
        \t\t}
        \t}
        \tif tv, ok := col.(jTypedViewCursorSource); ok {
        \t\tif n, ok2 := tv.typedViewLen_(); ok2 {
        \t\t\tif n == 0 {
        \t\t\t\treturn []T{}
        \t\t\t}
        \t\t\tout := make([]T, n)
        \t\t\tfor i := 0; i < n; i++ {
        \t\t\t\tout[i] = unboxJavaElem[T](tv.typedViewAt_(i))
        \t\t\t}
        \t\t\treturn out
        \t\t}
        \t}
        \tif si, ok := col.(SliceIterable); ok {
        \t\titems := si.IterSlice_()
        \t\tif len(items) == 0 {
        \t\t\treturn []T{}
        \t\t}
        \t\tout := make([]T, len(items))
        \t\tfor i, raw := range items {
        \t\t\tout[i] = unboxJavaElem[T](raw)
        \t\t}
        \t\treturn out
        \t}
        \tif c, ok := col.(Ext_java_util_Collection); ok {
        \t\tsz := int(c.Size())
        \t\tif sz == 0 {
        \t\t\treturn []T{}
        \t\t}
        \t\tout := make([]T, 0, sz)
        \t\tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\t\tout = append(out, unboxJavaElem[T](cur.Next()))
        \t\t}
        \t\treturn out
        \t}
        \tout := []T{}
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\tout = append(out, unboxJavaElem[T](cur.Next()))
        \t}
        \treturn out
        }

        func CopyIteratorToSlice[T any](it Ext_java_util_Iterator) []T {
        \tCheckNotNull(it)
        \tout := []T{}
        \tfor it.HasNext() {
        \t\tout = append(out, unboxJavaElem[T](it.Next()))
        \t}
        \treturn out
        }

        func SliceFormatJavaString[T any](s []T) string {
        \tif s == nil {
        \t\treturn "null"
        \t}
        \tif len(s) == 0 {
        \t\treturn "[]"
        \t}
        \tvar sb strings.Builder
        \tsb.WriteByte('[')
        \tfor i, v := range s {
        \t\tif i > 0 {
        \t\t\tsb.WriteString(", ")
        \t\t}
        \t\tsb.WriteString(StrVal(boxJavaElem(v)))
        \t}
        \tsb.WriteByte(']')
        \treturn sb.String()
        }

        func sliceTimSortBy[T any](s []T, less func(a, b T) bool) {
        \t// Stable insertion + merge sort matching Java TimSort stability.
        \tn := len(s)
        \tif n < 2 {
        \t\treturn
        \t}
        \tif n <= 24 {
        \t\tfor i := 1; i < n; i++ {
        \t\t\tv := s[i]
        \t\t\tj := i
        \t\t\tfor j > 0 && less(v, s[j-1]) {
        \t\t\t\ts[j] = s[j-1]
        \t\t\t\tj--
        \t\t\t}
        \t\t\ts[j] = v
        \t\t}
        \t\treturn
        \t}
        \tbuf := make([]T, n)
        \tvar sortRange func(lo, hi int)
        \tsortRange = func(lo, hi int) {
        \t\tif hi-lo <= 16 {
        \t\t\tfor i := lo + 1; i < hi; i++ {
        \t\t\t\tv := s[i]
        \t\t\t\tj := i
        \t\t\t\tfor j > lo && less(v, s[j-1]) {
        \t\t\t\t\ts[j] = s[j-1]
        \t\t\t\t\tj--
        \t\t\t\t}
        \t\t\t\ts[j] = v
        \t\t\t}
        \t\t\treturn
        \t\t}
        \t\tmid := lo + (hi-lo)/2
        \t\tsortRange(lo, mid)
        \t\tsortRange(mid, hi)
        \t\tif !less(s[mid], s[mid-1]) {
        \t\t\treturn
        \t\t}
        \t\tcopy(buf[lo:hi], s[lo:hi])
        \t\ti, j, k := lo, mid, lo
        \t\tfor i < mid && j < hi {
        \t\t\tif less(buf[j], buf[i]) {
        \t\t\t\ts[k] = buf[j]
        \t\t\t\tj++
        \t\t\t} else {
        \t\t\t\ts[k] = buf[i]
        \t\t\t\ti++
        \t\t\t}
        \t\t\tk++
        \t\t}
        \t\tif i < mid {
        \t\t\tcopy(s[k:hi], buf[i:mid])
        \t\t}
        \t}
        \tsortRange(0, n)
        }

        func SliceSortNatural[T any](s []T) {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif len(s) == 0 {
        \t\treturn
        \t}
        \tvar zero T
        \tif _, isStr := any(zero).(string); isStr {
        \t\tss := any(s).([]string)
        \t\tfor _, v := range ss {
        \t\t\tif v == NullString {
        \t\t\t\tThrowNPE()
        \t\t\t}
        \t\t}
        \t\tsliceTimSortBy(ss, func(a, b string) bool { return JStrCompareTo(a, b) < 0 })
        \t\treturn
        \t}
        \tif nullValKindOf[T]() != kindNeverNull {
        \t\tfor _, v := range s {
        \t\t\tif isJavaNullMapVal(v) {
        \t\t\t\tThrowNPE()
        \t\t\t}
        \t\t}
        \t}
        \tif len(s) == 1 {
        \t\treturn
        \t}
        \tsliceTimSortBy(s, func(a, b T) bool { return CompareValues(boxJavaElem(a), boxJavaElem(b)) < 0 })
        }

        func SliceSortComparator[T any](s []T, cmp Ext_java_util_Comparator) {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif isJavaNullMapVal(cmp) {
        \t\tSliceSortNatural(s)
        \t\treturn
        \t}
        \tif len(s) < 2 {
        \t\treturn
        \t}
        \tsliceTimSortBy(s, func(a, b T) bool { return cmp.Compare(boxJavaElem(a), boxJavaElem(b)) < 0 })
        }

        type jSliceIterableView[T any] struct {
        \ts []T
        }

        func (*jSliceIterableView[T]) Implements_j2gort_Ext_java_lang_Iterable() {}
        func (v *jSliceIterableView[T]) IsJavaNil_() bool                        { return v == nil }
        func (v *jSliceIterableView[T]) IterStrSlice_() ([]string, bool) {
        \tif ss, ok := any(v.s).([]string); ok {
        \t\treturn ss, true
        \t}
        \treturn nil, false
        }
        func (v *jSliceIterableView[T]) IterSlice_() []any {
        \tif len(v.s) == 0 {
        \t\treturn emptyAnySlice_
        \t}
        \tout := make([]any, len(v.s))
        \tfor i, x := range v.s {
        \t\tout[i] = boxJavaElem(x)
        \t}
        \treturn out
        }
        func (v *jSliceIterableView[T]) Iterator() any { return SliceIterator(v.s) }
        func (v *jSliceIterableView[T]) ForEach(c Ext_java_util_function_Consumer) {
        \tCheckNotNull(c)
        \tfor _, x := range v.s {
        \t\tc.Accept(boxJavaElem(x))
        \t}
        }
        func (v *jSliceIterableView[T]) String() string { return SliceFormatJavaString(v.s) }

        func SliceReadOnlyIterable[T any](s []T) Ext_java_lang_Iterable {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn &jSliceIterableView[T]{s: s}
        }

        type unmodSlicePtrView[T any] struct {
        \tExt_java_util_AbstractListBase
        \tp *[]T
        }

        type unmodSlicePtrIter[T any] struct {
        \tp   *[]T
        \tpos int
        }

        func (*unmodSlicePtrIter[T]) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *unmodSlicePtrIter[T]) IsJavaNil_() bool                       { return it == nil }
        func (it *unmodSlicePtrIter[T]) HasNext() bool                          { return it.pos < len(*it.p) }
        func (it *unmodSlicePtrIter[T]) Next() any {
        \ts := *it.p
        \tif it.pos >= len(s) {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tv := s[it.pos]
        \tit.pos++
        \treturn boxJavaElem(v)
        }
        func (it *unmodSlicePtrIter[T]) Remove() {
        \tpanic(javalang.NewUnsupportedOperationException("unmodifiable"))
        }

        func (*unmodSlicePtrView[T]) Implements_j2gort_Ext_java_util_RandomAccess() {}
        func (v *unmodSlicePtrView[T]) IsJavaNil_() bool                            { return v == nil }
        func (v *unmodSlicePtrView[T]) Size() int32                                 { return int32(len(*v.p)) }
        func (v *unmodSlicePtrView[T]) IsEmpty() bool                               { return len(*v.p) == 0 }
        func (v *unmodSlicePtrView[T]) Get(i int32) any {
        \ts := *v.p
        \tCheckElementIndex(i, int32(len(s)))
        \treturn boxJavaElem(s[i])
        }
        func (v *unmodSlicePtrView[T]) Contains(x any) bool {
        \tt := ToAny(x)
        \tfor _, e := range *v.p {
        \t\tif javalang.ObjectsEquals(boxJavaElem(e), t) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }
        func (v *unmodSlicePtrView[T]) IndexOf_1(x any) int32 {
        \tt := ToAny(x)
        \tfor i, e := range *v.p {
        \t\tif javalang.ObjectsEquals(boxJavaElem(e), t) {
        \t\t\treturn int32(i)
        \t\t}
        \t}
        \treturn -1
        }
        func (v *unmodSlicePtrView[T]) LastIndexOf_1(x any) int32 {
        \tt := ToAny(x)
        \ts := *v.p
        \tfor i := len(s) - 1; i >= 0; i-- {
        \t\tif javalang.ObjectsEquals(boxJavaElem(s[i]), t) {
        \t\t\treturn int32(i)
        \t\t}
        \t}
        \treturn -1
        }
        func (v *unmodSlicePtrView[T]) IterSlice_() []any {
        \ts := *v.p
        \tif len(s) == 0 {
        \t\treturn emptyAnySlice_
        \t}
        \tout := make([]any, len(s))
        \tfor i, e := range s {
        \t\tout[i] = boxJavaElem(e)
        \t}
        \treturn out
        }
        func (v *unmodSlicePtrView[T]) Iterator() any {
        \tif len(*v.p) == 0 {
        \t\treturn emptyImmSliceIter_
        \t}
        \treturn &unmodSlicePtrIter[T]{p: v.p}
        }
        func (v *unmodSlicePtrView[T]) ForEach(c Ext_java_util_function_Consumer) {
        \tCheckNotNull(c)
        \tfor _, e := range *v.p {
        \t\tc.Accept(boxJavaElem(e))
        \t}
        }
        func (v *unmodSlicePtrView[T]) Equals(o any) bool   { return ImmSliceEquals(*v.p, o) }
        func (v *unmodSlicePtrView[T]) HashCode() int32     { return ImmSliceHashCode(*v.p) }
        func (v *unmodSlicePtrView[T]) String() string      { return SliceFormatJavaString(*v.p) }
        func (v *unmodSlicePtrView[T]) Clear()              { panic(javalang.NewUnsupportedOperationException("unmodifiable")) }
        func (v *unmodSlicePtrView[T]) Add_1(any) bool      { panic(javalang.NewUnsupportedOperationException("unmodifiable")) }
        func (v *unmodSlicePtrView[T]) Remove_1_any(any) bool { panic(javalang.NewUnsupportedOperationException("unmodifiable")) }

        func SliceUnmodifiableListLive[T any](p *[]T) Ext_java_util_List {
        \tif p == nil || *p == nil {
        \t\tThrowNPE()
        \t}
        \tv := &unmodSlicePtrView[T]{p: p}
        \tv.BindExtSelf_(v)
        \treturn v
        }

        const orderedMapIndexThreshold = 12
        const orderedSetIndexThreshold = 8

        // OrderedSet is a slim value struct preserving Java LinkedHashSet insertion order
        // with zero inline array bloat and lazy backing-slice allocation on first Add.
        type OrderedSet[T comparable] struct {
        \torder    []T
        \tidx      map[T]int32
        \tholes    int32
        \tfirstIdx int32
        }

        type OrderedSetCursor[T comparable] struct {
        \ts   *OrderedSet[T]
        \tpos int32
        }

        func NewOrderedSet[T comparable]() OrderedSet[T] {
        \treturn OrderedSet[T]{}
        }

        func NewOrderedSetPtr[T comparable]() *OrderedSet[T] {
        \treturn &OrderedSet[T]{}
        }

        func NewOrderedSetCap[T comparable](cap int32) OrderedSet[T] {
        \tif cap < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", cap)))
        \t}
        \tvar s OrderedSet[T]
        \tif cap > 4 {
        \t\ts.order = make([]T, 0, cap)
        \t}
        \treturn s
        }

        func NewOrderedSetPtrCap[T comparable](cap int32) *OrderedSet[T] {
        \ts := NewOrderedSetCap[T](cap)
        \treturn &s
        }

        func OrderedSetFromSet[T comparable](src OrderedSet[T]) OrderedSet[T] {
        \tn := src.Size()
        \tif n == 0 {
        \t\treturn OrderedSet[T]{}
        \t}
        \tvar dst OrderedSet[T]
        \tdst.order = make([]T, 0, max(4, int(n)))
        \tdst.AddAllSet(src)
        \treturn dst
        }

        func OrderedSetPtrFromSet[T comparable](src *OrderedSet[T]) *OrderedSet[T] {
        \tif src == nil {
        \t\tThrowNPE()
        \t}
        \tdst := OrderedSetFromSet(*src)
        \treturn &dst
        }

        func OrderedSetFromStrCol(col any) OrderedSet[string] {
        \tvar s OrderedSet[string]
        \tOrderedSetAddAllStr(&s, col)
        \treturn s
        }

        func OrderedSetPtrFromStrCol(col any) *OrderedSet[string] {
        \ts := OrderedSetFromStrCol(col)
        \treturn &s
        }

        func OrderedSetFromNodeCol[N ~uint32](col any) OrderedSet[N] {
        \tvar s OrderedSet[N]
        \tOrderedSetAddAllNode(&s, col)
        \treturn s
        }

        func OrderedSetPtrFromNodeCol[N ~uint32](col any) *OrderedSet[N] {
        \ts := OrderedSetFromNodeCol[N](col)
        \treturn &s
        }

        func OrderedSetFromCol[T comparable](col any) OrderedSet[T] {
        \tvar s OrderedSet[T]
        \tOrderedSetAddAllCol(&s, col)
        \treturn s
        }

        func OrderedSetPtrFromCol[T comparable](col any) *OrderedSet[T] {
        \ts := OrderedSetFromCol[T](col)
        \treturn &s
        }

        func (s *OrderedSet[T]) Init() {}

        func (s *OrderedSet[T]) InitCap(cap int32) {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif cap < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", cap)))
        \t}
        \tif cap > 4 && s.order == nil {
        \t\ts.order = make([]T, 0, cap)
        \t}
        }

        func (s *OrderedSet[T]) IsNil() bool     { return s == nil }
        func (s *OrderedSet[T]) IsJavaNil_() bool { return s == nil }

        func (s OrderedSet[T]) Size() int32 {
        \treturn int32(len(s.order)) - s.holes
        }

        func (s OrderedSet[T]) IsEmpty() bool {
        \treturn int32(len(s.order)) == s.holes
        }

        func (s OrderedSet[T]) Contains(v T) bool {
        \tif int32(len(s.order)) == s.holes {
        \t\treturn false
        \t}
        \tif s.idx != nil {
        \t\treturn s.idx[v] != 0
        \t}
        \tfor _, x := range s.order {
        \t\tif x == v {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (s OrderedSet[T]) GetFirst() T {
        \tif int32(len(s.order)) == s.holes {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \treturn s.order[s.firstIdx]
        }

        func (s *OrderedSet[T]) Add(v T) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif s.order == nil {
        \t\ts.order = make([]T, 1, 4)
        \t\ts.order[0] = v
        \t\treturn true
        \t}
        \tif s.idx != nil {
        \t\tif s.idx[v] != 0 {
        \t\t\treturn false
        \t\t}
        \t\ts.order = append(s.order, v)
        \t\ts.idx[v] = int32(len(s.order))
        \t\treturn true
        \t}
        \tfor _, x := range s.order {
        \t\tif x == v {
        \t\t\treturn false
        \t\t}
        \t}
        \ts.order = append(s.order, v)
        \tif len(s.order) > orderedSetIndexThreshold {
        \t\tm := make(map[T]int32, 16)
        \t\tfor i, x := range s.order {
        \t\t\tm[x] = int32(i + 1)
        \t\t}
        \t\ts.idx = m
        \t}
        \treturn true
        }

        func (s *OrderedSet[T]) Remove(v T) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(s.order)) == s.holes {
        \t\treturn false
        \t}
        \tif s.idx == nil {
        \t\tfor i, x := range s.order {
        \t\t\tif x == v {
        \t\t\t\tvar zero T
        \t\t\t\tcopy(s.order[i:], s.order[i+1:])
        \t\t\t\tlast := len(s.order) - 1
        \t\t\t\ts.order[last] = zero
        \t\t\t\ts.order = s.order[:last]
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \t}
        \tpos1 := s.idx[v]
        \tif pos1 == 0 {
        \t\treturn false
        \t}
        \tdelete(s.idx, v)
        \tvar zero T
        \tslot := pos1 - 1
        \ts.order[slot] = zero
        \ts.holes++
        \tlive := int32(len(s.order)) - s.holes
        \tif live == 0 {
        \t\ts.order = s.order[:0]
        \t\ts.idx = nil
        \t\ts.holes = 0
        \t\ts.firstIdx = 0
        \t} else if len(s.order) > 16 && int(live)*2 < len(s.order) {
        \t\ts.compact()
        \t} else if slot == s.firstIdx {
        \t\tfor s.firstIdx++; int(s.firstIdx) < len(s.order); s.firstIdx++ {
        \t\t\tif s.idx[s.order[s.firstIdx]] == s.firstIdx+1 {
        \t\t\t\tbreak
        \t\t\t}
        \t\t}
        \t}
        \treturn true
        }

        func (s *OrderedSet[T]) compact() {
        \tw := 0
        \tfor i := int(s.firstIdx); i < len(s.order); i++ {
        \t\tv := s.order[i]
        \t\tif s.idx[v] == int32(i+1) {
        \t\t\ts.order[w] = v
        \t\t\ts.idx[v] = int32(w + 1)
        \t\t\tw++
        \t\t}
        \t}
        \tclear(s.order[w:])
        \ts.order = s.order[:w]
        \ts.holes = 0
        \ts.firstIdx = 0
        }

        func (s *OrderedSet[T]) Clear() {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif s.order != nil {
        \t\tclear(s.order)
        \t\ts.order = s.order[:0]
        \t}
        \ts.idx = nil
        \ts.holes = 0
        \ts.firstIdx = 0
        }

        func (s OrderedSet[T]) EqualsSet(other OrderedSet[T]) bool {
        \tif s.Size() != other.Size() {
        \t\treturn false
        \t}
        \tif s.IsEmpty() {
        \t\treturn true
        \t}
        \tif s.holes > 0 {
        \t\tfor i := int(s.firstIdx); i < len(s.order); i++ {
        \t\t\tv := s.order[i]
        \t\t\tif s.idx[v] == int32(i+1) && !other.Contains(v) {
        \t\t\t\treturn false
        \t\t\t}
        \t\t}
        \t\treturn true
        \t}
        \tfor _, v := range s.order {
        \t\tif !other.Contains(v) {
        \t\t\treturn false
        \t\t}
        \t}
        \treturn true
        }

        func (s *OrderedSet[T]) EqualsSetPtr(other *OrderedSet[T]) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif other == nil {
        \t\treturn false
        \t}
        \treturn s.EqualsSet(*other)
        }

        func (s OrderedSet[T]) HashCode() int32 {
        \tif s.IsEmpty() {
        \t\treturn 0
        \t}
        \tvar h int32
        \tif s.holes > 0 {
        \t\tfor i := int(s.firstIdx); i < len(s.order); i++ {
        \t\t\tv := s.order[i]
        \t\t\tif s.idx[v] == int32(i+1) {
        \t\t\t\th += javalang.ObjectHashCode(boxJavaElem(v))
        \t\t\t}
        \t\t}
        \t\treturn h
        \t}
        \tfor _, v := range s.order {
        \t\th += javalang.ObjectHashCode(boxJavaElem(v))
        \t}
        \treturn h
        }

        func (s *OrderedSet[T]) AddAllSet(other OrderedSet[T]) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif other.IsEmpty() {
        \t\treturn false
        \t}
        \tchanged := false
        \tif other.holes > 0 {
        \t\tfor i := int(other.firstIdx); i < len(other.order); i++ {
        \t\t\tv := other.order[i]
        \t\t\tif other.idx[v] == int32(i+1) && s.Add(v) {
        \t\t\t\tchanged = true
        \t\t\t}
        \t\t}
        \t\treturn changed
        \t}
        \tfor _, v := range other.order {
        \t\tif s.Add(v) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func (s *OrderedSet[T]) AddAllSetPtr(other *OrderedSet[T]) bool {
        \tif s == nil || other == nil {
        \t\tThrowNPE()
        \t}
        \treturn s.AddAllSet(*other)
        }

        func OrderedSetAddAllStr(s *OrderedSet[string], col any) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(col)
        \tchanged := false
        \tif ssi, ok := col.(StrSliceIterable); ok {
        \t\tif ss, ok2 := ssi.IterStrSlice_(); ok2 {
        \t\t\tfor _, k := range ss {
        \t\t\t\tif s.Add(k) {
        \t\t\t\t\tchanged = true
        \t\t\t\t}
        \t\t\t}
        \t\t\treturn changed
        \t\t}
        \t}
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\tif s.Add(cur.NextStr()) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func OrderedSetAddAllNode[N ~uint32](s *OrderedSet[N], col any) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(col)
        \tchanged := false
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\tif s.Add(N(cur.NextNode())) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func OrderedSetAddAllCol[T comparable](s *OrderedSet[T], col any) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(col)
        \tchanged := false
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\tif s.Add(unboxJavaElem[T](cur.Next())) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func (s *OrderedSet[T]) Iter() OrderedSetCursor[T] {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn OrderedSetCursor[T]{s: s, pos: s.firstIdx}
        }

        func (c *OrderedSetCursor[T]) HasNext() bool {
        \ts := c.s
        \tif int32(len(s.order)) == s.holes {
        \t\treturn false
        \t}
        \tif s.holes > 0 {
        \t\treturn c.advanceDead()
        \t}
        \treturn int(c.pos) < len(s.order)
        }

        func (c *OrderedSetCursor[T]) advanceDead() bool {
        \ts := c.s
        \titems := s.order
        \tidx := s.idx
        \tfor i := int(c.pos); i < len(items); i++ {
        \t\tif idx[items[i]] == int32(i+1) {
        \t\t\tc.pos = int32(i)
        \t\t\treturn true
        \t\t}
        \t}
        \tc.pos = int32(len(items))
        \treturn false
        }

        func (c *OrderedSetCursor[T]) Next() T {
        \ti := c.pos
        \tc.pos = i + 1
        \treturn c.s.order[i]
        }

        func (s *OrderedSet[T]) String() string {
        \tif s == nil {
        \t\treturn "null"
        \t}
        \tif s.IsEmpty() {
        \t\treturn "[]"
        \t}
        \tvar sb strings.Builder
        \tsb.WriteByte('[')
        \tfirst := true
        \tfor cur := s.Iter(); cur.HasNext(); {
        \t\tif !first {
        \t\t\tsb.WriteString(", ")
        \t\t}
        \t\tfirst = false
        \t\tsb.WriteString(StrVal(boxJavaElem(cur.Next())))
        \t}
        \tsb.WriteByte(']')
        \treturn sb.String()
        }

        func (s OrderedSet[T]) ToImmutableSet() Ext_com_google_common_collect_ImmutableSet {
        \tlive := s.Size()
        \tif live == 0 {
        \t\treturn emptyImmutableSet_
        \t}
        \tvar zero T
        \tif _, isStr := any(zero).(string); isStr {
        \t\tif live == 1 {
        \t\t\treturn ImmSetOfStr1(any(s.GetFirst()).(string))
        \t\t}
        \t\tstrs := make([]string, 0, live)
        \t\tif s.holes > 0 {
        \t\t\tfor i, v := range s.order {
        \t\t\t\tif s.idx[v] == int32(i+1) {
        \t\t\t\t\tstr := any(v).(string)
        \t\t\t\t\tif str == NullString {
        \t\t\t\t\t\tThrowNPE()
        \t\t\t\t\t}
        \t\t\t\t\tstrs = append(strs, str)
        \t\t\t\t}
        \t\t\t}
        \t\t} else {
        \t\t\tfor _, v := range s.order {
        \t\t\t\tstr := any(v).(string)
        \t\t\t\tif str == NullString {
        \t\t\t\t\tThrowNPE()
        \t\t\t\t}
        \t\t\t\tstrs = append(strs, str)
        \t\t\t}
        \t\t}
        \t\treturn finishImmStrSetOwn(strs, nil)
        \t}
        \tif live == 1 {
        \t\treturn newImmSet1(boxJavaElem(s.GetFirst()))
        \t}
        \tif live <= 8 {
        \t\telems := make([]any, 0, live)
        \t\tif s.holes > 0 {
        \t\t\tfor i, v := range s.order {
        \t\t\t\tif s.idx[v] == int32(i+1) {
        \t\t\t\t\telems = append(elems, CheckNotNull(boxJavaElem(v)))
        \t\t\t\t}
        \t\t\t}
        \t\t} else {
        \t\t\tfor _, v := range s.order {
        \t\t\t\telems = append(elems, CheckNotNull(boxJavaElem(v)))
        \t\t\t}
        \t\t}
        \t\treturn finishImmSetMulti(elems, nil)
        \t}
        \tsc := &JSetCore{}
        \tsc.InitExpectedSize(live)
        \tif s.holes > 0 {
        \t\tfor i, v := range s.order {
        \t\t\tif s.idx[v] == int32(i+1) {
        \t\t\t\tsc.Add(CheckNotNull(boxJavaElem(v)))
        \t\t\t}
        \t\t}
        \t} else {
        \t\tfor _, v := range s.order {
        \t\t\tsc.Add(CheckNotNull(boxJavaElem(v)))
        \t\t}
        \t}
        \treturn wrapImmutableSetFromSetCoreOwn(sc)
        }

        func (s OrderedSet[T]) ToJavaSet() *Ext_java_util_LinkedHashSet {
        \tlive := s.Size()
        \tobj := &Ext_java_util_LinkedHashSet{}
        \tobj.BindExtSelf_(obj)
        \tobj.inlineSet.InitExpectedSize(live)
        \tif live == 0 {
        \t\treturn obj
        \t}
        \tvar zero T
        \tif _, isStr := any(zero).(string); isStr {
        \t\tif s.holes > 0 {
        \t\t\tfor i, v := range s.order {
        \t\t\t\tif s.idx[v] == int32(i+1) {
        \t\t\t\t\tobj.inlineSet.AddStr(any(v).(string))
        \t\t\t\t}
        \t\t\t}
        \t\t} else {
        \t\t\tfor _, v := range s.order {
        \t\t\t\tobj.inlineSet.AddStr(any(v).(string))
        \t\t\t}
        \t\t}
        \t\treturn obj
        \t}
        \tif s.holes > 0 {
        \t\tfor i, v := range s.order {
        \t\t\tif s.idx[v] == int32(i+1) {
        \t\t\t\tobj.inlineSet.Add(boxJavaElem(v))
        \t\t\t}
        \t\t}
        \t} else {
        \t\tfor _, v := range s.order {
        \t\t\tobj.inlineSet.Add(boxJavaElem(v))
        \t\t}
        \t}
        \treturn obj
        }

        func ArrayDequeFromOrderedSet[T comparable](s OrderedSet[T]) *Ext_java_util_ArrayDeque {
        \tlive := s.Size()
        \tobj := &Ext_java_util_ArrayDeque{}
        \tobj.BindExtSelf_(obj)
        \tdc := &obj.inlineCore
        \tdc.InitCap(live)
        \tif live == 0 {
        \t\treturn obj
        \t}
        \tif s.holes > 0 {
        \t\tfor i, v := range s.order {
        \t\t\tif s.idx[v] == int32(i+1) {
        \t\t\t\tdc.AddLast(boxJavaElem(v))
        \t\t\t}
        \t\t}
        \t} else {
        \t\tfor _, v := range s.order {
        \t\t\tdc.AddLast(boxJavaElem(v))
        \t\t}
        \t}
        \treturn obj
        }

        func ArrayDequeFromOrderedSetPtr[T comparable](s *OrderedSet[T]) *Ext_java_util_ArrayDeque {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn ArrayDequeFromOrderedSet(*s)
        }

        // LocalOrderedSet is a stack-only OrderedSet with [4]T inline storage so standalone stack
        // local sets with <= 4 items allocate zero heap memory while struct fields remain slim.
        type LocalOrderedSet[T comparable] struct {
        \tinline   [4]T
        \titems    []T
        \tidx      map[T]int32
        \tlive     int32
        \tfirstIdx int32
        \tlenInl   uint8
        }

        type LocalOrderedSetCursor[T comparable] struct {
        \ts   *LocalOrderedSet[T]
        \tpos int32
        }

        func NewLocalOrderedSet[T comparable]() LocalOrderedSet[T] {
        \treturn LocalOrderedSet[T]{}
        }

        func NewLocalOrderedSetCap[T comparable](cap int32) LocalOrderedSet[T] {
        \tif cap < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", cap)))
        \t}
        \tvar s LocalOrderedSet[T]
        \tif cap > 4 {
        \t\ts.items = make([]T, 0, cap)
        \t}
        \treturn s
        }

        func LocalOrderedSetFromStrCol(col any) LocalOrderedSet[string] {
        \tvar s LocalOrderedSet[string]
        \tLocalOrderedSetAddAllStr(&s, col)
        \treturn s
        }

        func LocalOrderedSetFromNodeCol[N ~uint32](col any) LocalOrderedSet[N] {
        \tvar s LocalOrderedSet[N]
        \tLocalOrderedSetAddAllNode(&s, col)
        \treturn s
        }

        func LocalOrderedSetFromCol[T comparable](col any) LocalOrderedSet[T] {
        \tvar s LocalOrderedSet[T]
        \tLocalOrderedSetAddAllCol(&s, col)
        \treturn s
        }

        func (s LocalOrderedSet[T]) Size() int32    { return s.live }
        func (s LocalOrderedSet[T]) IsEmpty() bool  { return s.live == 0 }
        func (s LocalOrderedSet[T]) Contains(v T) bool {
        \tif s.live == 0 {
        \t\treturn false
        \t}
        \tif s.idx != nil {
        \t\treturn s.idx[v] != 0
        \t}
        \tif s.items != nil {
        \t\tfor _, x := range s.items {
        \t\t\tif x == v {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \t}
        \tfor i := uint8(0); i < s.lenInl; i++ {
        \t\tif s.inline[i] == v {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (s LocalOrderedSet[T]) GetFirst() T {
        \tif s.live == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \tif s.items == nil {
        \t\treturn s.inline[0]
        \t}
        \treturn s.items[s.firstIdx]
        }

        func (s *LocalOrderedSet[T]) Add(v T) bool {
        \tif s.items == nil {
        \t\tfor i := uint8(0); i < s.lenInl; i++ {
        \t\t\tif s.inline[i] == v {
        \t\t\t\treturn false
        \t\t\t}
        \t\t}
        \t\tif s.lenInl < 4 {
        \t\t\ts.inline[s.lenInl] = v
        \t\t\ts.lenInl++
        \t\t\ts.live++
        \t\t\treturn true
        \t\t}
        \t\titems := make([]T, 5, 8)
        \t\tcopy(items, s.inline[:4])
        \t\titems[4] = v
        \t\tvar z [4]T
        \t\ts.inline = z
        \t\ts.lenInl = 0
        \t\ts.items = items
        \t\ts.live = 5
        \t\treturn true
        \t}
        \tif s.idx != nil {
        \t\tif s.idx[v] != 0 {
        \t\t\treturn false
        \t\t}
        \t\ts.items = append(s.items, v)
        \t\ts.idx[v] = int32(len(s.items))
        \t\ts.live++
        \t\treturn true
        \t}
        \tfor _, x := range s.items {
        \t\tif x == v {
        \t\t\treturn false
        \t\t}
        \t}
        \ts.items = append(s.items, v)
        \ts.live++
        \tif len(s.items) > orderedSetIndexThreshold {
        \t\tm := make(map[T]int32, 16)
        \t\tfor i, x := range s.items {
        \t\t\tm[x] = int32(i + 1)
        \t\t}
        \t\ts.idx = m
        \t}
        \treturn true
        }

        func (s *LocalOrderedSet[T]) Remove(v T) bool {
        \tif s.live == 0 {
        \t\treturn false
        \t}
        \tif s.items == nil {
        \t\tfor i := uint8(0); i < s.lenInl; i++ {
        \t\t\tif s.inline[i] == v {
        \t\t\t\tvar zero T
        \t\t\t\tcopy(s.inline[i:], s.inline[i+1:s.lenInl])
        \t\t\t\ts.lenInl--
        \t\t\t\ts.inline[s.lenInl] = zero
        \t\t\t\ts.live--
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \t}
        \tif s.idx == nil {
        \t\tfor i, x := range s.items {
        \t\t\tif x == v {
        \t\t\t\tvar zero T
        \t\t\t\tcopy(s.items[i:], s.items[i+1:])
        \t\t\t\tlast := len(s.items) - 1
        \t\t\t\ts.items[last] = zero
        \t\t\t\ts.items = s.items[:last]
        \t\t\t\ts.live--
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \t}
        \tpos1 := s.idx[v]
        \tif pos1 == 0 {
        \t\treturn false
        \t}
        \tdelete(s.idx, v)
        \tvar zero T
        \tslot := pos1 - 1
        \ts.items[slot] = zero
        \ts.live--
        \tif s.live == 0 {
        \t\ts.items = s.items[:0]
        \t\ts.idx = nil
        \t\ts.firstIdx = 0
        \t} else if len(s.items) > 16 && int(s.live)*2 < len(s.items) {
        \t\tw := 0
        \t\tfor i := int(s.firstIdx); i < len(s.items); i++ {
        \t\t\tx := s.items[i]
        \t\t\tif s.idx[x] == int32(i+1) {
        \t\t\t\ts.items[w] = x
        \t\t\t\ts.idx[x] = int32(w + 1)
        \t\t\t\tw++
        \t\t\t}
        \t\t}
        \t\tclear(s.items[w:])
        \t\ts.items = s.items[:w]
        \t\ts.firstIdx = 0
        \t} else if slot == s.firstIdx {
        \t\tfor s.firstIdx++; int(s.firstIdx) < len(s.items); s.firstIdx++ {
        \t\t\tif s.idx[s.items[s.firstIdx]] == s.firstIdx+1 {
        \t\t\t\tbreak
        \t\t\t}
        \t\t}
        \t}
        \treturn true
        }

        func (s *LocalOrderedSet[T]) Clear() {
        \tvar z [4]T
        \ts.inline = z
        \ts.lenInl = 0
        \tif s.items != nil {
        \t\tclear(s.items)
        \t\ts.items = s.items[:0]
        \t}
        \ts.idx = nil
        \ts.live = 0
        \ts.firstIdx = 0
        }

        func (s *LocalOrderedSet[T]) HashCode() int32 {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn s.toOrderedSetView().HashCode()
        }

        func (s *LocalOrderedSet[T]) String() string {
        \tif s == nil {
        \t\treturn "null"
        \t}
        \tv := s.toOrderedSetView()
        \treturn (&v).String()
        }

        func LocalOrderedSetAddAllStr(s *LocalOrderedSet[string], col any) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(col)
        \tchanged := false
        \tif ssi, ok := col.(StrSliceIterable); ok {
        \t\tif ss, ok2 := ssi.IterStrSlice_(); ok2 {
        \t\t\tfor _, k := range ss {
        \t\t\t\tif s.Add(k) {
        \t\t\t\t\tchanged = true
        \t\t\t\t}
        \t\t\t}
        \t\t\treturn changed
        \t\t}
        \t}
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\tif s.Add(cur.NextStr()) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func LocalOrderedSetAddAllNode[N ~uint32](s *LocalOrderedSet[N], col any) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(col)
        \tchanged := false
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\tif s.Add(N(cur.NextNode())) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func LocalOrderedSetAddAllCol[T comparable](s *LocalOrderedSet[T], col any) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(col)
        \tchanged := false
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\tif s.Add(unboxJavaElem[T](cur.Next())) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func (s *LocalOrderedSet[T]) Iter() LocalOrderedSetCursor[T] {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn LocalOrderedSetCursor[T]{s: s, pos: s.firstIdx}
        }

        func (c *LocalOrderedSetCursor[T]) HasNext() bool {
        \ts := c.s
        \tif s.live == 0 {
        \t\treturn false
        \t}
        \tif s.items == nil {
        \t\treturn c.pos < int32(s.lenInl)
        \t}
        \tif s.idx != nil && int(s.live) < len(s.items) {
        \t\tfor i := int(c.pos); i < len(s.items); i++ {
        \t\t\tif s.idx[s.items[i]] == int32(i+1) {
        \t\t\t\tc.pos = int32(i)
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\tc.pos = int32(len(s.items))
        \t\treturn false
        \t}
        \treturn int(c.pos) < len(s.items)
        }

        func (c *LocalOrderedSetCursor[T]) Next() T {
        \ts := c.s
        \ti := c.pos
        \tc.pos = i + 1
        \tif s.items == nil {
        \t\treturn s.inline[i]
        \t}
        \treturn s.items[i]
        }

        func (s *LocalOrderedSet[T]) toOrderedSetView() OrderedSet[T] {
        \tif s.live == 0 {
        \t\treturn OrderedSet[T]{}
        \t}
        \tif s.items == nil {
        \t\treturn OrderedSet[T]{order: s.inline[:s.lenInl]}
        \t}
        \treturn OrderedSet[T]{order: s.items, idx: s.idx, holes: int32(len(s.items)) - s.live, firstIdx: s.firstIdx}
        }

        func (s *LocalOrderedSet[T]) ToImmutableSet() Ext_com_google_common_collect_ImmutableSet {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn s.toOrderedSetView().ToImmutableSet()
        }

        func ArrayDequeFromLocalOrderedSet[T comparable](s *LocalOrderedSet[T]) *Ext_java_util_ArrayDeque {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn ArrayDequeFromOrderedSet(s.toOrderedSetView())
        }
        """;
  }

  private static String buildColRtOrderedMapHelpers() {
    return """

        type jIsNiler interface {
        \tIsNil() bool
        }

        const (
        \tkindNullFallback   uint8 = 0
        \tkindNeverNull      uint8 = 1
        \tkindZeroIsNull     uint8 = 2
        \tkindNullStringType uint8 = 3
        )

        var nullValKindMap sync.Map

        func computeNullValKind[V any](rt reflect.Type) uint8 {
        \tif rt == nil || rt.Kind() == reflect.Interface {
        \t\treturn kindNullFallback
        \t}
        \tif rt == reflect.TypeFor[string]() {
        \t\treturn kindNullStringType
        \t}
        \tvar zero V
        \tif nc, ok := any(zero).(jJavaNilChecker); ok {
        \t\tif !nc.IsJavaNil_() {
        \t\t\treturn kindNeverNull
        \t\t}
        \t\tif rt.Comparable() {
        \t\t\treturn kindZeroIsNull
        \t\t}
        \t\treturn kindNullFallback
        \t}
        \tif rn, ok := any(zero).(jRawNodeProvider); ok {
        \t\tif rn.RawNodeId_() != 0 {
        \t\t\treturn kindNeverNull
        \t\t}
        \t\tif rt.Comparable() {
        \t\t\treturn kindZeroIsNull
        \t\t}
        \t\treturn kindNullFallback
        \t}
        \tif rt.Implements(reflect.TypeFor[jIsNiler]()) {
        \t\treturn kindNullFallback
        \t}
        \tswitch rt.Kind() {
        \tcase reflect.Bool,
        \t\treflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
        \t\treflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint64, reflect.Uintptr,
        \t\treflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
        \t\treflect.Struct, reflect.Array:
        \t\treturn kindNeverNull
        \tcase reflect.Uint32, reflect.Pointer, reflect.UnsafePointer:
        \t\treturn kindZeroIsNull
        \t}
        \treturn kindNullFallback
        }

        func nullValKindOf[V any]() uint8 {
        \trt := reflect.TypeFor[V]()
        \tif k, ok := nullValKindMap.Load(rt); ok {
        \t\treturn k.(uint8)
        \t}
        \tk := computeNullValKind[V](rt)
        \tnullValKindMap.Store(rt, k)
        \treturn k
        }

        func isJavaNullMapVal[V any](v V) bool {
        \tif any(v) == nil {
        \t\treturn true
        \t}
        \tswitch x := any(v).(type) {
        \tcase string:
        \t\treturn x == NullString
        \tcase int32, int64, bool, float64, float32, uint16, int8, int16, uint8, int, uint64:
        \t\treturn false
        \tdefault:
        \t\tswitch nullValKindOf[V]() {
        \t\tcase kindNeverNull:
        \t\t\treturn false
        \t\tcase kindZeroIsNull:
        \t\t\tvar zero V
        \t\t\treturn any(v) == any(zero)
        \t\tcase kindNullStringType:
        \t\t\treturn any(v).(string) == NullString
        \t\tdefault:
        \t\t\treturn isJavaNullMapValSlow(v)
        \t\t}
        \t}
        }

        //go:noinline
        func isJavaNullMapValSlow[V any](v V) bool {
        \tvar zero V
        \tswitch any(zero).(type) {
        \tcase nil:
        \t\tx := any(v)
        \t\tif x == nil {
        \t\t\treturn true
        \t\t}
        \t\tif nc, ok := x.(jJavaNilChecker); ok {
        \t\t\treturn nc.IsJavaNil_()
        \t\t}
        \t\tif rn, ok := x.(jRawNodeProvider); ok {
        \t\t\treturn rn.RawNodeId_() == 0
        \t\t}
        \t\tif nd, ok := x.(nodestore.Node); ok {
        \t\t\treturn uint32(nd) == 0
        \t\t}
        \t\treturn isJavaNilAny(x)
        \tcase string:
        \t\treturn any(v).(string) == NullString
        \tcase int32, int64, bool, float64, float32, uint16, int8, int16:
        \t\treturn false
        \tcase uint32:
        \t\treturn any(v).(uint32) == 0
        \tcase nodestore.Node:
        \t\treturn uint32(any(v).(nodestore.Node)) == 0
        \tcase jJavaNilChecker:
        \t\treturn any(v).(jJavaNilChecker).IsJavaNil_()
        \tcase jRawNodeProvider:
        \t\treturn any(v).(jRawNodeProvider).RawNodeId_() == 0
        \tdefault:
        \t\treturn isJavaNilAny(any(v))
        \t}
        }

        func javaMapValsEqual[V any](a, b V, identity bool) bool {
        \tba, bb := boxJavaElem(a), boxJavaElem(b)
        \tif identity {
        \t\treturn ba == bb
        \t}
        \treturn javalang.ObjectsEquals(ba, bb)
        }

        func javaMapEntryHash[K, V any](k K, v V, identity bool) int32 {
        \tbk, bv := boxJavaElem(k), boxJavaElem(v)
        \tif identity {
        \t\treturn IdentityHash(bk) ^ IdentityHash(bv)
        \t}
        \treturn ObjectHash(bk) ^ ObjectHash(bv)
        }

        // OrderedMap is a slim 40-byte struct preserving Java LinkedHashMap/HashMap/IdentityHashMap order
        // for Go-comparable identity/primitive/string/Node keys with zero any-boxing and lazy backing allocation.
        type OrderedMapEntry[K comparable, V any] struct {
        \tKey K
        \tVal V
        }

        type OrderedMap[K comparable, V any] struct {
        \tentries  []OrderedMapEntry[K, V]
        \tidx      map[K]int32
        \tholes    int32
        \tfirstIdx int32
        }

        const (
        \tordMapPoolSmallMax = 8192
        \tordMapPoolMedMax   = 1024
        \tordMapPoolLargeMax = 16
        \tordMapPoolHugeMax  = 4
        \tordMapPtrSmallMax  = 16384
        \tordMapPtrLargeMax  = 64
        \tordMapPtrHugeMax   = 4
        )

        type orderedMapBuf[K comparable, V any] struct {
        \tentries []OrderedMapEntry[K, V]
        \tidx     map[K]int32
        }

        type orderedMapTypePool[K comparable, V any] struct {
        \tmu    sync.Mutex
        \tsmall []orderedMapBuf[K, V]
        \tmed   []orderedMapBuf[K, V]
        \tlarge []orderedMapBuf[K, V]
        \thuge  []orderedMapBuf[K, V]
        }

        type orderedMapPtrTypePool[K comparable, V any] struct {
        \tmu        sync.Mutex
        \tsmallFree []*OrderedMap[K, V]
        \tlargeFree []*OrderedMap[K, V]
        \thugeFree  []*OrderedMap[K, V]
        }

        var (
        \torderedMapPools    sync.Map
        \torderedMapPtrPools sync.Map
        )

        func getOrderedMapPool[K comparable, V any]() *orderedMapTypePool[K, V] {
        \trt := reflect.TypeFor[orderedMapBuf[K, V]]()
        \tif v, ok := orderedMapPools.Load(rt); ok {
        \t\treturn v.(*orderedMapTypePool[K, V])
        \t}
        \tp := &orderedMapTypePool[K, V]{}
        \tactual, loaded := orderedMapPools.LoadOrStore(rt, p)
        \tif !loaded {
        \t\tpurge := func() {
        \t\t\tp.mu.Lock()
        \t\t\tp.small = nil
        \t\t\tp.med = nil
        \t\t\tp.large = nil
        \t\t\tp.huge = nil
        \t\t\tp.mu.Unlock()
        \t\t}
        \t\tregisterCompPoolReset(purge)
        \t\tregisterRuntimePoolPurger(purge)
        \t}
        \treturn actual.(*orderedMapTypePool[K, V])
        }

        func getOrderedMapPtrPool[K comparable, V any]() *orderedMapPtrTypePool[K, V] {
        \trt := reflect.TypeFor[OrderedMap[K, V]]()
        \tif v, ok := orderedMapPtrPools.Load(rt); ok {
        \t\treturn v.(*orderedMapPtrTypePool[K, V])
        \t}
        \tp := &orderedMapPtrTypePool[K, V]{}
        \tactual, loaded := orderedMapPtrPools.LoadOrStore(rt, p)
        \tif !loaded {
        \t\tpurge := func() {
        \t\t\tp.mu.Lock()
        \t\t\tp.smallFree = nil
        \t\t\tp.largeFree = nil
        \t\t\tp.hugeFree = nil
        \t\t\tp.mu.Unlock()
        \t\t}
        \t\tregisterCompPoolReset(purge)
        \t\tregisterRuntimePoolPurger(purge)
        \t}
        \treturn actual.(*orderedMapPtrTypePool[K, V])
        }

        func AcquireOrderedMap[K comparable, V any]() OrderedMap[K, V] {
        \tp := getOrderedMapPool[K, V]()
        \tp.mu.Lock()
        \tif n := len(p.small); n > 0 {
        \t\tb := p.small[n-1]
        \t\tp.small[n-1] = orderedMapBuf[K, V]{}
        \t\tp.small = p.small[:n-1]
        \t\tp.mu.Unlock()
        \t\treturn OrderedMap[K, V]{entries: b.entries[:0], idx: b.idx}
        \t}
        \tp.mu.Unlock()
        \treturn OrderedMap[K, V]{}
        }

        func AcquireOrderedMapCap[K comparable, V any](capHint int32) OrderedMap[K, V] {
        \tif capHint < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", capHint)))
        \t}
        \tif capHint <= 64 {
        \t\tm := AcquireOrderedMap[K, V]()
        \t\tif cap(m.entries) == 0 && capHint > 8 {
        \t\t\tif capHint >= 64 {
        \t\t\t\tm.entries = AcquirePooledSlice[OrderedMapEntry[K, V]](int(capHint))
        \t\t\t} else {
        \t\t\t\tm.entries = make([]OrderedMapEntry[K, V], 0, capHint)
        \t\t\t}
        \t\t}
        \t\tif m.idx == nil && int(capHint) > orderedMapIndexThreshold {
        \t\t\tm.idx = make(map[K]int32, capHint)
        \t\t}
        \t\treturn m
        \t}
        \tp := getOrderedMapPool[K, V]()
        \tif capHint < 1024 {
        \t\tp.mu.Lock()
        \t\tif n := len(p.med); n > 0 {
        \t\t\tb := p.med[n-1]
        \t\t\tp.med[n-1] = orderedMapBuf[K, V]{}
        \t\t\tp.med = p.med[:n-1]
        \t\t\tp.mu.Unlock()
        \t\t\treturn OrderedMap[K, V]{entries: b.entries[:0], idx: b.idx}
        \t\t}
        \t\tif n := len(p.large); n > 0 {
        \t\t\tb := p.large[n-1]
        \t\t\tp.large[n-1] = orderedMapBuf[K, V]{}
        \t\t\tp.large = p.large[:n-1]
        \t\t\tp.mu.Unlock()
        \t\t\treturn OrderedMap[K, V]{entries: b.entries[:0], idx: b.idx}
        \t\t}
        \t\tp.mu.Unlock()
        \t\treturn NewOrderedMapCap[K, V](capHint)
        \t}
        \tp.mu.Lock()
        \tif n := len(p.huge); n > 0 {
        \t\tb := p.huge[n-1]
        \t\tp.huge[n-1] = orderedMapBuf[K, V]{}
        \t\tp.huge = p.huge[:n-1]
        \t\tp.mu.Unlock()
        \t\treturn OrderedMap[K, V]{entries: b.entries[:0], idx: b.idx}
        \t}
        \tif capHint <= 65536 {
        \t\tif n := len(p.large); n > 0 {
        \t\t\tb := p.large[n-1]
        \t\t\tp.large[n-1] = orderedMapBuf[K, V]{}
        \t\t\tp.large = p.large[:n-1]
        \t\t\tp.mu.Unlock()
        \t\t\treturn OrderedMap[K, V]{entries: b.entries[:0], idx: b.idx}
        \t\t}
        \t\tif capHint <= 2048 {
        \t\t\tif n := len(p.med); n > 0 {
        \t\t\t\tb := p.med[n-1]
        \t\t\t\tp.med[n-1] = orderedMapBuf[K, V]{}
        \t\t\t\tp.med = p.med[:n-1]
        \t\t\t\tp.mu.Unlock()
        \t\t\t\treturn OrderedMap[K, V]{entries: b.entries[:0], idx: b.idx}
        \t\t\t}
        \t\t}
        \t}
        \tp.mu.Unlock()
        \treturn NewOrderedMapCap[K, V](capHint)
        }

        func AcquireOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] {
        \tp := getOrderedMapPtrPool[K, V]()
        \tp.mu.Lock()
        \tif n := len(p.smallFree); n > 0 {
        \t\tm := p.smallFree[n-1]
        \t\tp.smallFree[n-1] = nil
        \t\tp.smallFree = p.smallFree[:n-1]
        \t\tp.mu.Unlock()
        \t\treturn m
        \t}
        \tp.mu.Unlock()
        \treturn &OrderedMap[K, V]{}
        }

        func AcquireOrderedMapPtrCap[K comparable, V any](capHint int32) *OrderedMap[K, V] {
        \tif capHint < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", capHint)))
        \t}
        \tp := getOrderedMapPtrPool[K, V]()
        \tvar m *OrderedMap[K, V]
        \tp.mu.Lock()
        \tif capHint >= 1024 {
        \t\tif n := len(p.hugeFree); n > 0 {
        \t\t\tm = p.hugeFree[n-1]
        \t\t\tp.hugeFree[n-1] = nil
        \t\t\tp.hugeFree = p.hugeFree[:n-1]
        \t\t} else if n := len(p.largeFree); n > 0 {
        \t\t\tm = p.largeFree[n-1]
        \t\t\tp.largeFree[n-1] = nil
        \t\t\tp.largeFree = p.largeFree[:n-1]
        \t\t} else if n := len(p.smallFree); n > 0 {
        \t\t\tm = p.smallFree[n-1]
        \t\t\tp.smallFree[n-1] = nil
        \t\t\tp.smallFree = p.smallFree[:n-1]
        \t\t}
        \t} else if capHint > 256 {
        \t\tif n := len(p.largeFree); n > 0 {
        \t\t\tm = p.largeFree[n-1]
        \t\t\tp.largeFree[n-1] = nil
        \t\t\tp.largeFree = p.largeFree[:n-1]
        \t\t} else if n := len(p.hugeFree); n > 0 {
        \t\t\tm = p.hugeFree[n-1]
        \t\t\tp.hugeFree[n-1] = nil
        \t\t\tp.hugeFree = p.hugeFree[:n-1]
        \t\t}
        \t} else {
        \t\tif n := len(p.smallFree); n > 0 {
        \t\t\tm = p.smallFree[n-1]
        \t\t\tp.smallFree[n-1] = nil
        \t\t\tp.smallFree = p.smallFree[:n-1]
        \t\t}
        \t}
        \tp.mu.Unlock()
        \tif m == nil {
        \t\tm = &OrderedMap[K, V]{}
        \t}
        \tif capHint >= 64 && cap(m.entries) < 64 {
        \t\tm.entries = AcquirePooledSlice[OrderedMapEntry[K, V]](int(capHint))
        \t} else if cap(m.entries) == 0 && capHint > 0 {
        \t\tm.entries = make([]OrderedMapEntry[K, V], 0, capHint)
        \t}
        \tif m.idx == nil && int(capHint) > orderedMapIndexThreshold {
        \t\tm.idx = make(map[K]int32, capHint)
        \t}
        \treturn m
        }

        func ReleaseOrderedMap[K comparable, V any](m *OrderedMap[K, V]) {
        \tif m == nil {
        \t\treturn
        \t}
        \tc := cap(m.entries)
        \tif c > 4194304 {
        \t\tes := m.entries
        \t\tclear(es)
        \t\tm.entries = nil
        \t\tm.idx = nil
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\tif c >= 64 && c <= pooledSliceMaxCap {
        \t\t\tReleasePooledSlice(es)
        \t\t}
        \t\treturn
        \t}
        \tm.Clear()
        \tp := getOrderedMapPtrPool[K, V]()
        \tkept := false
        \tif c > 65536 {
        \t\tp.mu.Lock()
        \t\tif len(p.hugeFree) < ordMapPtrHugeMax {
        \t\t\tp.hugeFree = append(p.hugeFree, m)
        \t\t\tkept = true
        \t\t}
        \t\tp.mu.Unlock()
        \t} else if c > 256 {
        \t\tp.mu.Lock()
        \t\tif len(p.largeFree) < ordMapPtrLargeMax {
        \t\t\tp.largeFree = append(p.largeFree, m)
        \t\t\tkept = true
        \t\t}
        \t\tp.mu.Unlock()
        \t} else {
        \t\tp.mu.Lock()
        \t\tif len(p.smallFree) < ordMapPtrSmallMax {
        \t\t\tp.smallFree = append(p.smallFree, m)
        \t\t\tkept = true
        \t\t}
        \t\tp.mu.Unlock()
        \t}
        \tif !kept && c >= 64 && c <= pooledSliceMaxCap {
        \t\tes := m.entries
        \t\tm.entries = nil
        \t\tReleasePooledSlice(es)
        \t}
        }

        func ReleaseOrderedMapPtr[K comparable, V any](m *OrderedMap[K, V]) {
        \tReleaseOrderedMap(m)
        }

        func (m *OrderedMap[K, V]) ReleaseToPool() {
        \tif m == nil {
        \t\treturn
        \t}
        \tes := m.entries
        \tidx := m.idx
        \tholes := m.holes
        \tfirstIdx := m.firstIdx
        \t*m = OrderedMap[K, V]{}
        \tc := cap(es)
        \tif c == 0 && idx == nil {
        \t\treturn
        \t}
        \tif c > 4194304 {
        \t\tif c >= 64 && c <= pooledSliceMaxCap {
        \t\t\tReleasePooledSlice(es)
        \t\t}
        \t\treturn
        \t}
        \tn := len(es)
        \tif idx != nil && len(idx) != 0 {
        \t\tlive := n - int(holes)
        \t\tif live > 0 && live <= 32 && n-int(firstIdx) <= 64 && c >= 128 {
        \t\t\tfor i := int(firstIdx); i < n; i++ {
        \t\t\t\tdelete(idx, es[i].Key)
        \t\t\t}
        \t\t\tif len(idx) != 0 {
        \t\t\t\tclear(idx)
        \t\t\t}
        \t\t} else {
        \t\t\tclear(idx)
        \t\t}
        \t}
        \tif n > 0 {
        \t\tclear(es)
        \t\tes = es[:0]
        \t}
        \tp := getOrderedMapPool[K, V]()
        \tkept := false
        \tif c <= 64 {
        \t\tp.mu.Lock()
        \t\tif len(p.small) < ordMapPoolSmallMax {
        \t\t\tp.small = append(p.small, orderedMapBuf[K, V]{entries: es, idx: idx})
        \t\t\tkept = true
        \t\t}
        \t\tp.mu.Unlock()
        \t} else if c <= 2048 {
        \t\tp.mu.Lock()
        \t\tif len(p.med) < ordMapPoolMedMax {
        \t\t\tp.med = append(p.med, orderedMapBuf[K, V]{entries: es, idx: idx})
        \t\t\tkept = true
        \t\t}
        \t\tp.mu.Unlock()
        \t} else if c <= 65536 {
        \t\tp.mu.Lock()
        \t\tif len(p.large) < ordMapPoolLargeMax {
        \t\t\tp.large = append(p.large, orderedMapBuf[K, V]{entries: es, idx: idx})
        \t\t\tkept = true
        \t\t}
        \t\tp.mu.Unlock()
        \t} else {
        \t\tp.mu.Lock()
        \t\tif len(p.huge) < ordMapPoolHugeMax {
        \t\t\tp.huge = append(p.huge, orderedMapBuf[K, V]{entries: es, idx: idx})
        \t\t\tkept = true
        \t\t}
        \t\tp.mu.Unlock()
        \t}
        \tif !kept && c >= 64 && c <= pooledSliceMaxCap {
        \t\tReleasePooledSlice(es)
        \t}
        }

        func (m *OrderedMap[K, V]) isEntryActive(i int) bool {
        \tif m.holes == 0 {
        \t\treturn true
        \t}
        \tk := m.entries[i].Key
        \treturn m.idx[k] == int32(i+1) || k != k
        }

        func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
        \treturn OrderedMap[K, V]{}
        }

        func NewOrderedMapPtr[K comparable, V any]() *OrderedMap[K, V] {
        \treturn &OrderedMap[K, V]{}
        }

        func NewOrderedMapCap[K comparable, V any](cap int32) OrderedMap[K, V] {
        \tif cap < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", cap)))
        \t}
        \tvar m OrderedMap[K, V]
        \tif cap >= 64 {
        \t\tm.entries = AcquirePooledSlice[OrderedMapEntry[K, V]](int(cap))
        \t} else if cap > 8 {
        \t\tm.entries = make([]OrderedMapEntry[K, V], 0, cap)
        \t}
        \tif int(cap) > orderedMapIndexThreshold {
        \t\tm.idx = make(map[K]int32, cap)
        \t}
        \treturn m
        }

        func NewOrderedMapPtrCap[K comparable, V any](cap int32) *OrderedMap[K, V] {
        \tm := NewOrderedMapCap[K, V](cap)
        \treturn &m
        }

        func OrderedMapFromPtr[K comparable, V any](src *OrderedMap[K, V]) OrderedMap[K, V] {
        \tif src == nil {
        \t\tThrowNPE()
        \t}
        \tlive := int32(len(src.entries)) - src.holes
        \tif live == 0 {
        \t\treturn OrderedMap[K, V]{}
        \t}
        \tcapHint := int(live)
        \tif capHint < 4 {
        \t\tcapHint = 4
        \t}
        \tvar es []OrderedMapEntry[K, V]
        \tif capHint >= 64 {
        \t\tes = AcquirePooledSlice[OrderedMapEntry[K, V]](capHint)[:live]
        \t} else {
        \t\tes = make([]OrderedMapEntry[K, V], live, capHint)
        \t}
        \tif src.holes == 0 {
        \t\tcopy(es, src.entries)
        \t} else {
        \t\tw := 0
        \t\tfor i := int(src.firstIdx); i < len(src.entries); i++ {
        \t\t\tif src.isEntryActive(i) {
        \t\t\t\tes[w] = src.entries[i]
        \t\t\t\tw++
        \t\t\t}
        \t\t}
        \t}
        \tdst := OrderedMap[K, V]{entries: es}
        \tif len(es) > orderedMapIndexThreshold {
        \t\tidx := make(map[K]int32, max(16, len(es)*2))
        \t\tfor i := range es {
        \t\t\tidx[es[i].Key] = int32(i + 1)
        \t\t}
        \t\tdst.idx = idx
        \t}
        \treturn dst
        }

        func OrderedMapFrom[K comparable, V any](src OrderedMap[K, V]) OrderedMap[K, V] {
        \treturn OrderedMapFromPtr(&src)
        }

        func OrderedMapPtrFrom[K comparable, V any](src *OrderedMap[K, V]) *OrderedMap[K, V] {
        \tdst := OrderedMapFromPtr(src)
        \treturn &dst
        }

        func OrderedMapFromJavaMap[K comparable, V any](src any) OrderedMap[K, V] {
        \tCheckNotNull(src)
        \tmp := src.(Ext_java_util_Map)
        \tm := NewOrderedMapCap[K, V](mp.Size())
        \tfor cur := NewMapEntryCursor(mp); cur.NextEntry(); {
        \t\tm.Put(unboxJavaElem[K](cur.GetKey()), unboxJavaElem[V](cur.GetValue()))
        \t}
        \treturn m
        }

        func OrderedMapPtrFromJavaMap[K comparable, V any](src any) *OrderedMap[K, V] {
        \tm := OrderedMapFromJavaMap[K, V](src)
        \treturn &m
        }

        func OrderedMapFromNodeKeyJavaMap[K ~uint32, V any](src any) OrderedMap[K, V] {
        \tCheckNotNull(src)
        \tmp := src.(Ext_java_util_Map)
        \tm := NewOrderedMapCap[K, V](mp.Size())
        \tfor cur := NewMapEntryCursor(mp); cur.NextEntry(); {
        \t\tm.Put(K(cur.KeyNode()), unboxJavaElem[V](cur.GetValue()))
        \t}
        \treturn m
        }

        func OrderedMapPtrFromNodeKeyJavaMap[K ~uint32, V any](src any) *OrderedMap[K, V] {
        \tm := OrderedMapFromNodeKeyJavaMap[K, V](src)
        \treturn &m
        }

        func (m *OrderedMap[K, V]) IsNil() bool     { return m == nil }
        func (m *OrderedMap[K, V]) IsJavaNil_() bool { return m == nil }

        func (m *OrderedMap[K, V]) Size() int32 {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \treturn int32(len(m.entries)) - m.holes
        }

        func (m *OrderedMap[K, V]) IsEmpty() bool {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \treturn int32(len(m.entries)) == m.holes
        }

        func (m *OrderedMap[K, V]) GetOk(k K) (V, bool) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tvar zero V
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn zero, false
        \t}
        \tif m.idx != nil {
        \t\tpos1 := m.idx[k]
        \t\tif pos1 == 0 {
        \t\t\treturn zero, false
        \t\t}
        \t\treturn m.entries[pos1-1].Val, true
        \t}
        \tfor i := range m.entries {
        \t\tif m.entries[i].Key == k {
        \t\t\treturn m.entries[i].Val, true
        \t\t}
        \t}
        \treturn zero, false
        }

        func (m *OrderedMap[K, V]) Get(k K) V {
        \tv, _ := m.GetOk(k)
        \treturn v
        }

        func (m *OrderedMap[K, V]) GetOrDefault(k K, def V) V {
        \tif v, ok := m.GetOk(k); ok {
        \t\treturn v
        \t}
        \treturn def
        }

        func (m *OrderedMap[K, V]) ContainsKey(k K) bool {
        \t_, ok := m.GetOk(k)
        \treturn ok
        }

        func (m *OrderedMap[K, V]) appendEntry(e OrderedMapEntry[K, V]) {
        \tes := m.entries
        \tif len(es) < cap(es) {
        \t\tm.entries = append(es, e)
        \t\treturn
        \t}
        \tm.appendEntryGrow(e)
        }

        func (m *OrderedMap[K, V]) appendEntryGrow(e OrderedMapEntry[K, V]) {
        \told := m.entries
        \toldLen := len(old)
        \toldCap := cap(old)
        \tnewCap := oldCap * 2
        \tif newCap < 4 {
        \t\tnewCap = 4
        \t}
        \tif newCap <= oldLen {
        \t\tnewCap = oldLen + 1
        \t}
        \tif newCap < 64 {
        \t\tnext := make([]OrderedMapEntry[K, V], oldLen+1, newCap)
        \t\tcopy(next, old)
        \t\tnext[oldLen] = e
        \t\tm.entries = next
        \t\treturn
        \t}
        \tnext := AcquirePooledSlice[OrderedMapEntry[K, V]](newCap)
        \tnext = next[:oldLen+1]
        \tcopy(next, old)
        \tnext[oldLen] = e
        \tif oldCap >= 64 {
        \t\tReleasePooledSlice(old)
        \t}
        \tm.entries = next
        }

        func (m *OrderedMap[K, V]) PutOk(k K, v V) (V, bool) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tvar zero V
        \tif m.idx != nil {
        \t\tif pos1 := m.idx[k]; pos1 != 0 {
        \t\t\told := m.entries[pos1-1].Val
        \t\t\tm.entries[pos1-1].Val = v
        \t\t\treturn old, true
        \t\t}
        \t\tm.appendEntry(OrderedMapEntry[K, V]{Key: k, Val: v})
        \t\tm.idx[k] = int32(len(m.entries))
        \t\treturn zero, false
        \t}
        \tif m.entries == nil {
        \t\tes := make([]OrderedMapEntry[K, V], 1, 4)
        \t\tes[0] = OrderedMapEntry[K, V]{Key: k, Val: v}
        \t\tm.entries = es
        \t\treturn zero, false
        \t}
        \tfor i := range m.entries {
        \t\tif m.entries[i].Key == k {
        \t\t\told := m.entries[i].Val
        \t\t\tm.entries[i].Val = v
        \t\t\treturn old, true
        \t\t}
        \t}
        \tm.appendEntry(OrderedMapEntry[K, V]{Key: k, Val: v})
        \tif len(m.entries) > orderedMapIndexThreshold {
        \t\tidx := make(map[K]int32, max(24, cap(m.entries)))
        \t\tfor i := range m.entries {
        \t\t\tidx[m.entries[i].Key] = int32(i + 1)
        \t\t}
        \t\tm.idx = idx
        \t}
        \treturn zero, false
        }

        func (m *OrderedMap[K, V]) Put(k K, v V) V {
        \told, _ := m.PutOk(k, v)
        \treturn old
        }

        func (m *OrderedMap[K, V]) PutIfAbsent(k K, v V) V {
        \tcur, ok := m.GetOk(k)
        \tif !ok || isJavaNullMapVal(cur) {
        \t\tm.PutOk(k, v)
        \t}
        \treturn cur
        }

        func (m *OrderedMap[K, V]) RemoveOk(k K) (V, bool) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tvar zero V
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn zero, false
        \t}
        \tif m.idx == nil {
        \t\tfor i := range m.entries {
        \t\t\tif m.entries[i].Key == k {
        \t\t\t\told := m.entries[i].Val
        \t\t\t\tcopy(m.entries[i:], m.entries[i+1:])
        \t\t\t\tlast := len(m.entries) - 1
        \t\t\t\tm.entries[last] = OrderedMapEntry[K, V]{}
        \t\t\t\tm.entries = m.entries[:last]
        \t\t\t\treturn old, true
        \t\t\t}
        \t\t}
        \t\treturn zero, false
        \t}
        \tpos1 := m.idx[k]
        \tif pos1 == 0 {
        \t\treturn zero, false
        \t}
        \tdelete(m.idx, k)
        \tslot := pos1 - 1
        \told := m.entries[slot].Val
        \tm.entries[slot] = OrderedMapEntry[K, V]{}
        \tm.holes++
        \tlive := int32(len(m.entries)) - m.holes
        \tif live == 0 {
        \t\tm.entries = m.entries[:0]
        \t\tif m.idx != nil {
        \t\t\tclear(m.idx)
        \t\t}
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t} else {
        \t\tif slot == m.firstIdx {
        \t\t\tfor int(m.firstIdx) < len(m.entries) && !m.isEntryActive(int(m.firstIdx)) {
        \t\t\t\tm.firstIdx++
        \t\t\t}
        \t\t}
        \t\tif len(m.entries) > 16 && int(live)*2 < len(m.entries) {
        \t\t\tm.compact()
        \t\t}
        \t}
        \treturn old, true
        }

        func (m *OrderedMap[K, V]) Remove(k K) V {
        \told, _ := m.RemoveOk(k)
        \treturn old
        }

        func (m *OrderedMap[K, V]) compact() {
        \tw := 0
        \tfor i := int(m.firstIdx); i < len(m.entries); i++ {
        \t\tk := m.entries[i].Key
        \t\tif m.idx[k] == int32(i+1) || k != k {
        \t\t\tm.entries[w] = m.entries[i]
        \t\t\tif k == k {
        \t\t\t\tm.idx[k] = int32(w + 1)
        \t\t\t}
        \t\t\tw++
        \t\t}
        \t}
        \tclear(m.entries[w:])
        \tm.entries = m.entries[:w]
        \tm.holes = 0
        \tm.firstIdx = 0
        }

        func (m *OrderedMap[K, V]) Clear() {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tn := len(m.entries)
        \tif n == 0 {
        \t\tif len(m.idx) != 0 {
        \t\t\tclear(m.idx)
        \t\t}
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\treturn
        \t}
        \tif m.idx != nil && len(m.idx) != 0 {
        \t\tlive := n - int(m.holes)
        \t\tif live > 0 && live <= 32 && n-int(m.firstIdx) <= 64 && cap(m.entries) >= 128 {
        \t\t\tfor i := int(m.firstIdx); i < n; i++ {
        \t\t\t\tdelete(m.idx, m.entries[i].Key)
        \t\t\t}
        \t\t\tif len(m.idx) != 0 {
        \t\t\t\tclear(m.idx)
        \t\t\t}
        \t\t} else {
        \t\t\tclear(m.idx)
        \t\t}
        \t}
        \tclear(m.entries)
        \tm.entries = m.entries[:0]
        \tm.holes = 0
        \tm.firstIdx = 0
        }

        func (m *OrderedMap[K, V]) PutAllMapPtr(other *OrderedMap[K, V]) {
        \tif m == nil || other == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(other.entries)) == other.holes {
        \t\treturn
        \t}
        \tif cap(m.entries) == 0 && m.idx == nil {
        \t\t*m = OrderedMapFromPtr(other)
        \t\treturn
        \t}
        \tfor i := int(other.firstIdx); i < len(other.entries); i++ {
        \t\tif other.isEntryActive(i) {
        \t\t\tm.PutOk(other.entries[i].Key, other.entries[i].Val)
        \t\t}
        \t}
        }

        func (m *OrderedMap[K, V]) PutAllMap(other OrderedMap[K, V]) {
        \tm.PutAllMapPtr(&other)
        }

        func (m *OrderedMap[K, V]) PutAllJavaMap(src any) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(src)
        \tmp := src.(Ext_java_util_Map)
        \tfor cur := NewMapEntryCursor(mp); cur.NextEntry(); {
        \t\tm.PutOk(unboxJavaElem[K](cur.GetKey()), unboxJavaElem[V](cur.GetValue()))
        \t}
        }

        func OrderedMapPtrGetStr[K comparable](m *OrderedMap[K, string], k K) string {
        \tif v, ok := m.GetOk(k); ok {
        \t\treturn v
        \t}
        \treturn NullString
        }

        func OrderedMapGetStr[K comparable](m OrderedMap[K, string], k K) string {
        \treturn OrderedMapPtrGetStr(&m, k)
        }

        func OrderedMapPutStr[K comparable](m *OrderedMap[K, string], k K, v string) string {
        \tif old, ok := m.PutOk(k, v); ok {
        \t\treturn old
        \t}
        \treturn NullString
        }

        func OrderedMapRemoveStr[K comparable](m *OrderedMap[K, string], k K) string {
        \tif old, ok := m.RemoveOk(k); ok {
        \t\treturn old
        \t}
        \treturn NullString
        }

        func OrderedMapPutIfAbsentStr[K comparable](m *OrderedMap[K, string], k K, v string) string {
        \tcur, ok := m.GetOk(k)
        \tif !ok || cur == NullString {
        \t\tm.PutOk(k, v)
        \t\treturn NullString
        \t}
        \treturn cur
        }

        func OrderedMapComputeIfAbsent[K comparable, V comparable](m *OrderedMap[K, V], k K, nullV V, fn func(K) V) V {
        \tif m == nil || fn == nil {
        \t\tThrowNPE()
        \t}
        \tif cur, ok := m.GetOk(k); ok && cur != nullV {
        \t\treturn cur
        \t}
        \tnv := fn(k)
        \tif nv != nullV {
        \t\tm.PutOk(k, nv)
        \t}
        \treturn nv
        }

        func OrderedMapComputeIfAbsentFn[K comparable, V comparable](m *OrderedMap[K, V], k K, nullV V, fn Ext_java_util_function_Function) V {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(fn)
        \tif cur, ok := m.GetOk(k); ok && cur != nullV {
        \t\treturn cur
        \t}
        \traw := fn.Apply(boxJavaElem(k))
        \tif raw == nil || isJavaNilAny(raw) {
        \t\treturn nullV
        \t}
        \tnv := unboxJavaElem[V](raw)
        \tif nv != nullV {
        \t\tm.PutOk(k, nv)
        \t}
        \treturn nv
        }

        type OrderedMapKeyCursor[K comparable, V any] struct {
        \tm   *OrderedMap[K, V]
        \tpos int32
        }

        func (m *OrderedMap[K, V]) IterKeys() OrderedMapKeyCursor[K, V] {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn OrderedMapKeyCursor[K, V]{m: m}
        \t}
        \treturn OrderedMapKeyCursor[K, V]{m: m, pos: m.firstIdx}
        }

        func (c *OrderedMapKeyCursor[K, V]) HasNext() bool {
        \tm := c.m
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn false
        \t}
        \tfor int(c.pos) < len(m.entries) {
        \t\tif m.isEntryActive(int(c.pos)) {
        \t\t\treturn true
        \t\t}
        \t\tc.pos++
        \t}
        \treturn false
        }

        func (c *OrderedMapKeyCursor[K, V]) Next() K {
        \tm := c.m
        \ti := c.pos
        \tc.pos = i + 1
        \treturn m.entries[i].Key
        }

        type OrderedMapValCursor[K comparable, V any] struct {
        \tm   *OrderedMap[K, V]
        \tpos int32
        }

        func (m *OrderedMap[K, V]) IterValues() OrderedMapValCursor[K, V] {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn OrderedMapValCursor[K, V]{m: m}
        \t}
        \treturn OrderedMapValCursor[K, V]{m: m, pos: m.firstIdx}
        }

        func (c *OrderedMapValCursor[K, V]) HasNext() bool {
        \tm := c.m
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn false
        \t}
        \tfor int(c.pos) < len(m.entries) {
        \t\tif m.isEntryActive(int(c.pos)) {
        \t\t\treturn true
        \t\t}
        \t\tc.pos++
        \t}
        \treturn false
        }

        func (c *OrderedMapValCursor[K, V]) Next() V {
        \tm := c.m
        \ti := c.pos
        \tc.pos = i + 1
        \treturn m.entries[i].Val
        }

        type OrderedMapEntryCursor[K comparable, V any] struct {
        \tm        *OrderedMap[K, V]
        \tpos      int32
        \tlastSlot int32
        }

        func (m *OrderedMap[K, V]) IterEntries() OrderedMapEntryCursor[K, V] {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn OrderedMapEntryCursor[K, V]{m: m, lastSlot: -1}
        \t}
        \treturn OrderedMapEntryCursor[K, V]{m: m, pos: m.firstIdx, lastSlot: -1}
        }

        func (c *OrderedMapEntryCursor[K, V]) NextEntry() bool {
        \tm := c.m
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn false
        \t}
        \tfor int(c.pos) < len(m.entries) {
        \t\tif m.isEntryActive(int(c.pos)) {
        \t\t\tc.lastSlot = c.pos
        \t\t\tc.pos++
        \t\t\treturn true
        \t\t}
        \t\tc.pos++
        \t}
        \treturn false
        }

        func (c *OrderedMapEntryCursor[K, V]) Key() K {
        \treturn c.m.entries[c.lastSlot].Key
        }

        func (c *OrderedMapEntryCursor[K, V]) Value() V {
        \treturn c.m.entries[c.lastSlot].Val
        }

        func (c *OrderedMapEntryCursor[K, V]) SetValue(nv V) V {
        \told := c.m.entries[c.lastSlot].Val
        \tc.m.entries[c.lastSlot].Val = nv
        \treturn old
        }

        func (m *OrderedMap[K, V]) ForEach(fn func(K, V)) {
        \tif m == nil || fn == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn
        \t}
        \tfor i := int(m.firstIdx); i < len(m.entries); i++ {
        \t\tif m.isEntryActive(i) {
        \t\t\tfn(m.entries[i].Key, m.entries[i].Val)
        \t\t}
        \t}
        }

        func (m *OrderedMap[K, V]) ForEachBiConsumer(action Ext_java_util_function_BiConsumer) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(action)
        \tm.ForEach(func(k K, v V) {
        \t\taction.Accept(boxJavaElem(k), boxJavaElem(v))
        \t})
        }

        func OrderedMapPtrEquals[K comparable, V comparable](a, b *OrderedMap[K, V], identityVal bool) bool {
        \tif a == nil {
        \t\tThrowNPE()
        \t}
        \tif b == nil {
        \t\treturn false
        \t}
        \taLive := int32(len(a.entries)) - a.holes
        \tbLive := int32(len(b.entries)) - b.holes
        \tif aLive != bLive {
        \t\treturn false
        \t}
        \tif aLive == 0 {
        \t\treturn true
        \t}
        \tfor i := int(a.firstIdx); i < len(a.entries); i++ {
        \t\tif a.isEntryActive(i) {
        \t\t\tv2, ok := b.GetOk(a.entries[i].Key)
        \t\t\tif !ok || (a.entries[i].Val != v2 && (identityVal || !javaMapValsEqual(a.entries[i].Val, v2, false))) {
        \t\t\t\treturn false
        \t\t\t}
        \t\t}
        \t}
        \treturn true
        }

        func OrderedMapEquals[K comparable, V comparable](a, b OrderedMap[K, V], identityVal bool) bool {
        \treturn OrderedMapPtrEquals(&a, &b, identityVal)
        }

        func (m *OrderedMap[K, V]) EqualsMapPtr(other *OrderedMap[K, V], identityVal bool) bool {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif other == nil {
        \t\treturn false
        \t}
        \tmLive := int32(len(m.entries)) - m.holes
        \toLive := int32(len(other.entries)) - other.holes
        \tif mLive != oLive {
        \t\treturn false
        \t}
        \tif mLive == 0 {
        \t\treturn true
        \t}
        \tfor i := int(m.firstIdx); i < len(m.entries); i++ {
        \t\tif m.isEntryActive(i) {
        \t\t\tv2, ok := other.GetOk(m.entries[i].Key)
        \t\t\tif !ok || !javaMapValsEqual(m.entries[i].Val, v2, identityVal) {
        \t\t\t\treturn false
        \t\t\t}
        \t\t}
        \t}
        \treturn true
        }

        func (m *OrderedMap[K, V]) EqualsMap(other OrderedMap[K, V], identityVal bool) bool {
        \treturn m.EqualsMapPtr(&other, identityVal)
        }

        func (m *OrderedMap[K, V]) HashCodeMap(identity bool) int32 {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tvar h int32
        \tm.ForEach(func(k K, v V) {
        \t\th += javaMapEntryHash(k, v, identity)
        \t})
        \treturn h
        }

        func (m *OrderedMap[K, V]) String() string {
        \tif m == nil {
        \t\treturn "null"
        \t}
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn "{}"
        \t}
        \tvar sb strings.Builder
        \tsb.WriteByte('{')
        \tfirst := true
        \tm.ForEach(func(k K, v V) {
        \t\tif !first {
        \t\t\tsb.WriteString(", ")
        \t\t}
        \t\tfirst = false
        \t\tsb.WriteString(StrVal(boxJavaElem(k)))
        \t\tsb.WriteByte('=')
        \t\tsb.WriteString(StrVal(boxJavaElem(v)))
        \t})
        \tsb.WriteByte('}')
        \treturn sb.String()
        }

        func (m *OrderedMap[K, V]) ToImmutableMap() Ext_com_google_common_collect_ImmutableMap {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tlive := int32(len(m.entries)) - m.holes
        \tif live == 0 {
        \t\treturn emptyImmutableMap_
        \t}
        \tobj := newImmutableMapShell()
        \tmc := &obj.inlineCore
        \tmc.InitExpectedSize(live)
        \tm.ForEach(func(k K, v V) {
        \t\tmc.Put(CheckNotNull(boxJavaElem(k)), CheckNotNull(boxJavaElem(v)))
        \t})
        \tmc.sealImmutable()
        \treturn obj
        }

        func (m *OrderedMap[K, V]) KeysToImmutableSet() Ext_com_google_common_collect_ImmutableSet {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tlive := int32(len(m.entries)) - m.holes
        \tif live == 0 {
        \t\treturn emptyImmutableSet_
        \t}
        \ts := NewOrderedSetCap[K](live)
        \tm.ForEach(func(k K, _ V) {
        \t\ts.Add(k)
        \t})
        \treturn s.ToImmutableSet()
        }

        func (m *OrderedMap[K, V]) ValuesToImmutableList() Ext_com_google_common_collect_ImmutableList {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tlive := int32(len(m.entries)) - m.holes
        \tif live == 0 {
        \t\treturn emptyImmutableList_
        \t}
        \titems := make([]V, 0, live)
        \tm.ForEach(func(_ K, v V) {
        \t\titems = append(items, v)
        \t})
        \treturn ImmListFromSlice(items)
        }

        func (m *OrderedMap[K, V]) ToJavaLinkedHashMap() *Ext_java_util_LinkedHashMap {
        \tif m == nil {
        \t\treturn nil
        \t}
        \tlive := int32(len(m.entries)) - m.holes
        \tobj := &Ext_java_util_LinkedHashMap{}
        \tobj.BindExtSelf_(obj)
        \tmc := &obj.inlineCore
        \tmc.InitExpectedSize(live)
        \tm.ForEach(func(k K, v V) {
        \t\tmc.Put(boxJavaElem(k), boxJavaElem(v))
        \t})
        \treturn obj
        }

        func (m *OrderedMap[K, V]) ToJavaMap(identity bool) Ext_java_util_Map {
        \tif m == nil {
        \t\treturn nil
        \t}
        \tif identity {
        \t\tlive := int32(len(m.entries)) - m.holes
        \t\tobj := NewExt_java_util_IdentityHashMap()
        \t\tmc := &obj.inlineCore
        \t\tmc.InitExpectedSize(live)
        \t\tm.ForEach(func(k K, v V) {
        \t\t\tmc.Put(boxJavaElem(k), boxJavaElem(v))
        \t\t})
        \t\treturn obj
        \t}
        \treturn m.ToJavaLinkedHashMap()
        }

        type jOrdMapLiveValues[K comparable, V any] struct {
        \tExt_java_util_AbstractCollectionBase
        \tm *OrderedMap[K, V]
        }

        func (v *jOrdMapLiveValues[K, V]) IsJavaNil_() bool { return v == nil }
        func (v *jOrdMapLiveValues[K, V]) Size() int32 {
        \tif v == nil || v.m == nil {
        \t\tThrowNPE()
        \t}
        \treturn v.m.Size()
        }
        func (v *jOrdMapLiveValues[K, V]) IsEmpty() bool {
        \tif v == nil || v.m == nil {
        \t\tThrowNPE()
        \t}
        \treturn v.m.IsEmpty()
        }
        func (v *jOrdMapLiveValues[K, V]) Clear() {
        \tif v == nil || v.m == nil {
        \t\tThrowNPE()
        \t}
        \tv.m.Clear()
        }
        func (v *jOrdMapLiveValues[K, V]) Contains(x any) bool {
        \tif v == nil || v.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(v.m.entries)) == v.m.holes {
        \t\treturn false
        \t}
        \ttarget := ToAny(x)
        \tfor cur := v.m.IterValues(); cur.HasNext(); {
        \t\tif javalang.ObjectsEquals(boxJavaElem(cur.Next()), target) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }
        func ordMapNextActive[K comparable, V any](m *OrderedMap[K, V], from int32) int32 {
        \tn := int32(len(m.entries))
        \tif n == m.holes || from < 0 || from >= n {
        \t\treturn -1
        \t}
        \tif m.holes == 0 {
        \t\treturn from
        \t}
        \tfor i := int(from); i < int(n); i++ {
        \t\tif m.isEntryActive(i) {
        \t\t\treturn int32(i)
        \t\t}
        \t}
        \treturn -1
        }

        func (v *jOrdMapLiveValues[K, V]) iterInit_() (int32, bool) {
        \tif v == nil || v.m == nil {
        \t\tThrowNPE()
        \t}
        \tp := ordMapNextActive(v.m, v.m.firstIdx)
        \treturn p, p >= 0
        }
        func (v *jOrdMapLiveValues[K, V]) iterNext_(pos int32) (any, int32) {
        \tp := ordMapNextActive(v.m, pos)
        \tif p < 0 {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \treturn boxJavaElem(v.m.entries[p].Val), ordMapNextActive(v.m, p+1)
        }
        func (v *jOrdMapLiveValues[K, V]) Iterator() any {
        \tif v == nil || v.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(v.m.entries)) == v.m.holes {
        \t\treturn emptyMapIter_
        \t}
        \treturn &jOrdMapValIter[K, V]{m: v.m, pos: v.m.firstIdx, lastKeyValid: false}
        }

        type jOrdMapValIter[K comparable, V any] struct {
        \tm            *OrderedMap[K, V]
        \tpos          int32
        \tlastKey      K
        \tlastKeyValid bool
        }

        func (*jOrdMapValIter[K, V]) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jOrdMapValIter[K, V]) IsJavaNil_() bool                       { return it == nil }
        func (it *jOrdMapValIter[K, V]) HasNext() bool {
        \tif int32(len(it.m.entries)) == it.m.holes {
        \t\treturn false
        \t}
        \tfor int(it.pos) < len(it.m.entries) {
        \t\tif it.m.isEntryActive(int(it.pos)) {
        \t\t\treturn true
        \t\t}
        \t\tit.pos++
        \t}
        \treturn false
        }
        func (it *jOrdMapValIter[K, V]) Next() any {
        \tif !it.HasNext() {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \ti := it.pos
        \tit.pos = i + 1
        \tit.lastKey = it.m.entries[i].Key
        \tit.lastKeyValid = true
        \treturn boxJavaElem(it.m.entries[i].Val)
        }
        func (it *jOrdMapValIter[K, V]) Remove() {
        \tif !it.lastKeyValid {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \twasShifted := it.m.idx == nil
        \tit.m.RemoveOk(it.lastKey)
        \tif wasShifted && it.pos > 0 {
        \t\tit.pos--
        \t}
        \tit.lastKeyValid = false
        }

        func (m *OrderedMap[K, V]) LiveValues() Ext_java_util_Collection {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tv := &jOrdMapLiveValues[K, V]{m: m}
        \tv.BindExtSelf_(v)
        \treturn v
        }

        type jOrdMapLiveKeySet[K comparable, V any] struct {
        \tExt_java_util_AbstractSetBase
        \tm *OrderedMap[K, V]
        }

        func (s *jOrdMapLiveKeySet[K, V]) IsJavaNil_() bool { return s == nil }
        func (s *jOrdMapLiveKeySet[K, V]) Size() int32 {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \treturn s.m.Size()
        }
        func (s *jOrdMapLiveKeySet[K, V]) IsEmpty() bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \treturn s.m.IsEmpty()
        }
        func (s *jOrdMapLiveKeySet[K, V]) Clear() {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \ts.m.Clear()
        }
        func (s *jOrdMapLiveKeySet[K, V]) Contains(x any) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(s.m.entries)) == s.m.holes {
        \t\treturn false
        \t}
        \tk, ok := tryUnboxOrderedMapKey[K](x)
        \tif !ok {
        \t\treturn false
        \t}
        \treturn s.m.ContainsKey(k)
        }
        func (s *jOrdMapLiveKeySet[K, V]) Remove_1_any(x any) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(s.m.entries)) == s.m.holes {
        \t\treturn false
        \t}
        \tk, ok := tryUnboxOrderedMapKey[K](x)
        \tif !ok {
        \t\treturn false
        \t}
        \t_, removed := s.m.RemoveOk(k)
        \treturn removed
        }
        func (s *jOrdMapLiveKeySet[K, V]) SetContainsStr_(k string) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(s.m.entries)) == s.m.holes {
        \t\treturn false
        \t}
        \tvar zero K
        \tif any(zero) == nil {
        \t\tif k == NullString {
        \t\t\treturn s.m.ContainsKey(zero)
        \t\t}
        \t\tuk, ok := any(k).(K)
        \t\treturn ok && s.m.ContainsKey(uk)
        \t}
        \tif _, isStr := any(zero).(string); isStr {
        \t\treturn s.m.ContainsKey(any(k).(K))
        \t}
        \tif k == NullString {
        \t\tif uk, ok := tryUnboxOrderedMapKey[K](nil); ok {
        \t\t\treturn s.m.ContainsKey(uk)
        \t\t}
        \t}
        \treturn false
        }
        func (s *jOrdMapLiveKeySet[K, V]) SetAddStr_(_ string) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tpanic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
        }
        func (s *jOrdMapLiveKeySet[K, V]) SetRemoveStr_(k string) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(s.m.entries)) == s.m.holes {
        \t\treturn false
        \t}
        \tvar zero K
        \tif any(zero) == nil {
        \t\tif k == NullString {
        \t\t\t_, removed := s.m.RemoveOk(zero)
        \t\t\treturn removed
        \t\t}
        \t\tuk, ok := any(k).(K)
        \t\tif !ok {
        \t\t\treturn false
        \t\t}
        \t\t_, removed := s.m.RemoveOk(uk)
        \t\treturn removed
        \t}
        \tif _, isStr := any(zero).(string); isStr {
        \t\t_, removed := s.m.RemoveOk(any(k).(K))
        \t\treturn removed
        \t}
        \tif k == NullString {
        \t\tif uk, ok := tryUnboxOrderedMapKey[K](nil); ok {
        \t\t\t_, removed := s.m.RemoveOk(uk)
        \t\t\treturn removed
        \t\t}
        \t}
        \treturn false
        }
        func (s *jOrdMapLiveKeySet[K, V]) SetContainsNode_(u uint32) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(s.m.entries)) == s.m.holes {
        \t\treturn false
        \t}
        \tif u == 0 {
        \t\tuk, ok := tryUnboxOrderedMapKey[K](nil)
        \t\treturn ok && s.m.ContainsKey(uk)
        \t}
        \tvar zero K
        \tif any(zero) == nil {
        \t\tuk, ok := boxNodeById(u).(K)
        \t\treturn ok && s.m.ContainsKey(uk)
        \t}
        \tswitch any(zero).(type) {
        \tcase uint32:
        \t\treturn s.m.ContainsKey(any(u).(K))
        \tcase nodestore.Node:
        \t\treturn s.m.ContainsKey(any(nodestore.Node(u)).(K))
        \tcase jRawNodeProvider:
        \t\tif globalNodeBoxer != nil {
        \t\t\tif uk, ok := globalNodeBoxer(u).(K); ok {
        \t\t\t\treturn s.m.ContainsKey(uk)
        \t\t\t}
        \t\t}
        \t\tuk, ok := tryUnboxOrderedMapKey[K](u)
        \t\treturn ok && s.m.ContainsKey(uk)
        \tdefault:
        \t\treturn false
        \t}
        }
        func (s *jOrdMapLiveKeySet[K, V]) SetAddNode_(_ uint32) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tpanic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
        }
        func (s *jOrdMapLiveKeySet[K, V]) SetRemoveNode_(u uint32) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(s.m.entries)) == s.m.holes {
        \t\treturn false
        \t}
        \tif u == 0 {
        \t\tuk, ok := tryUnboxOrderedMapKey[K](nil)
        \t\tif !ok {
        \t\t\treturn false
        \t\t}
        \t\t_, removed := s.m.RemoveOk(uk)
        \t\treturn removed
        \t}
        \tvar zero K
        \tif any(zero) == nil {
        \t\tuk, ok := boxNodeById(u).(K)
        \t\tif !ok {
        \t\t\treturn false
        \t\t}
        \t\t_, removed := s.m.RemoveOk(uk)
        \t\treturn removed
        \t}
        \tswitch any(zero).(type) {
        \tcase uint32:
        \t\t_, removed := s.m.RemoveOk(any(u).(K))
        \t\treturn removed
        \tcase nodestore.Node:
        \t\t_, removed := s.m.RemoveOk(any(nodestore.Node(u)).(K))
        \t\treturn removed
        \tcase jRawNodeProvider:
        \t\tif globalNodeBoxer != nil {
        \t\t\tif uk, ok := globalNodeBoxer(u).(K); ok {
        \t\t\t\t_, removed := s.m.RemoveOk(uk)
        \t\t\t\treturn removed
        \t\t\t}
        \t\t}
        \t\tuk, ok := tryUnboxOrderedMapKey[K](u)
        \t\tif !ok {
        \t\t\treturn false
        \t\t}
        \t\t_, removed := s.m.RemoveOk(uk)
        \t\treturn removed
        \tdefault:
        \t\treturn false
        \t}
        }
        func (s *jOrdMapLiveKeySet[K, V]) iterInit_() (int32, bool) {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tp := ordMapNextActive(s.m, s.m.firstIdx)
        \treturn p, p >= 0
        }
        func (s *jOrdMapLiveKeySet[K, V]) iterNext_(pos int32) (any, int32) {
        \tp := ordMapNextActive(s.m, pos)
        \tif p < 0 {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \treturn boxJavaElem(s.m.entries[p].Key), ordMapNextActive(s.m, p+1)
        }
        func (s *jOrdMapLiveKeySet[K, V]) Iterator() any {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(s.m.entries)) == s.m.holes {
        \t\treturn emptyMapIter_
        \t}
        \treturn &jOrdMapKeyIter[K, V]{m: s.m, pos: s.m.firstIdx}
        }

        type jOrdMapKeyIter[K comparable, V any] struct {
        \tm            *OrderedMap[K, V]
        \tpos          int32
        \tlastKey      K
        \tlastKeyValid bool
        }

        func (*jOrdMapKeyIter[K, V]) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jOrdMapKeyIter[K, V]) IsJavaNil_() bool                       { return it == nil }
        func (it *jOrdMapKeyIter[K, V]) HasNext() bool {
        \tif int32(len(it.m.entries)) == it.m.holes {
        \t\treturn false
        \t}
        \tfor int(it.pos) < len(it.m.entries) {
        \t\tif it.m.isEntryActive(int(it.pos)) {
        \t\t\treturn true
        \t\t}
        \t\tit.pos++
        \t}
        \treturn false
        }
        func (it *jOrdMapKeyIter[K, V]) Next() any {
        \tif !it.HasNext() {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \ti := it.pos
        \tit.pos = i + 1
        \tit.lastKey = it.m.entries[i].Key
        \tit.lastKeyValid = true
        \treturn boxJavaElem(it.lastKey)
        }
        func (it *jOrdMapKeyIter[K, V]) Remove() {
        \tif !it.lastKeyValid {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \twasShifted := it.m.idx == nil
        \tit.m.RemoveOk(it.lastKey)
        \tif wasShifted && it.pos > 0 {
        \t\tit.pos--
        \t}
        \tit.lastKeyValid = false
        }

        func (m *OrderedMap[K, V]) LiveKeySet() Ext_java_util_Set {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \ts := &jOrdMapLiveKeySet[K, V]{m: m}
        \ts.BindExtSelf_(s)
        \treturn s
        }

        func (m *OrderedMap[K, V]) FirstValue() V {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tfor i := int(m.firstIdx); i < len(m.entries); i++ {
        \t\tif m.isEntryActive(i) {
        \t\t\treturn m.entries[i].Val
        \t\t}
        \t}
        \tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        }

        func OrderedMapGetNodeAny[K ~uint32, V any](m *OrderedMap[K, V], rawKey any) any {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tvar k K
        \tif rawKey != nil && !isJavaNilAny(rawKey) {
        \t\tif nk, ok := rawKey.(K); ok {
        \t\t\tk = nk
        \t\t} else if rn, ok := rawKey.(jRawNodeProvider); ok {
        \t\t\tk = K(rn.RawNodeId_())
        \t\t} else if nd, ok := rawKey.(nodestore.Node); ok {
        \t\t\tk = K(nd)
        \t\t} else if u, ok := rawKey.(uint32); ok {
        \t\t\tk = K(u)
        \t\t} else {
        \t\t\treturn nil
        \t\t}
        \t}
        \tv, found := m.GetOk(k)
        \tif !found {
        \t\treturn nil
        \t}
        \treturn boxJavaElem(v)
        }

        var boxedNullStringAny_ any = NullString

        func tryUnboxOrderedMapKey[K comparable](raw any) (K, bool) {
        \tvar zero K
        \tif any(zero) == nil {
        \t\tif raw == nil || isJavaNilAny(raw) {
        \t\t\treturn zero, true
        \t\t}
        \t\tif !reflect.TypeOf(raw).Comparable() {
        \t\t\treturn zero, false
        \t\t}
        \t\tk, ok := raw.(K)
        \t\treturn k, ok
        \t}
        \tif raw != nil {
        \t\tif k, ok := raw.(K); ok {
        \t\t\treturn k, true
        \t\t}
        \t}
        \tif raw == nil || isJavaNilAny(raw) {
        \t\tswitch any(zero).(type) {
        \t\tcase string:
        \t\t\treturn boxedNullStringAny_.(K), true
        \t\tcase int32, int64, bool, float64, float32, uint16, int8, int16:
        \t\t\treturn zero, false
        \t\tdefault:
        \t\t\treturn zero, true
        \t\t}
        \t}
        \tswitch any(zero).(type) {
        \tcase string:
        \t\treturn zero, false
        \tcase int32:
        \t\tif b, ok := raw.(*BoxInt); ok && b != nil {
        \t\t\treturn any(b.Val).(K), true
        \t\t}
        \t\treturn zero, false
        \tcase int64:
        \t\tif b, ok := raw.(*BoxLong); ok && b != nil {
        \t\t\treturn any(b.Val).(K), true
        \t\t}
        \t\treturn zero, false
        \tcase bool:
        \t\tif b, ok := raw.(*BoxBoolean); ok && b != nil {
        \t\t\treturn any(b.Val).(K), true
        \t\t}
        \t\treturn zero, false
        \tcase float64:
        \t\tif b, ok := raw.(*BoxDouble); ok && b != nil {
        \t\t\treturn any(b.Val).(K), true
        \t\t}
        \t\treturn zero, false
        \tcase float32:
        \t\tif b, ok := raw.(*BoxFloat); ok && b != nil {
        \t\t\treturn any(b.Val).(K), true
        \t\t}
        \t\treturn zero, false
        \tcase uint16:
        \t\tif b, ok := raw.(*BoxChar); ok && b != nil {
        \t\t\treturn any(b.Val).(K), true
        \t\t}
        \t\treturn zero, false
        \tcase int8:
        \t\tif b, ok := raw.(*BoxByte); ok && b != nil {
        \t\t\treturn any(b.Val).(K), true
        \t\t}
        \t\treturn zero, false
        \tcase int16:
        \t\tif b, ok := raw.(*BoxShort); ok && b != nil {
        \t\t\treturn any(b.Val).(K), true
        \t\t}
        \t\treturn zero, false
        \tcase uint32:
        \t\tif rn, ok := raw.(jRawNodeProvider); ok {
        \t\t\treturn any(rn.RawNodeId_()).(K), true
        \t\t}
        \t\tif nd, ok := raw.(nodestore.Node); ok {
        \t\t\treturn any(uint32(nd)).(K), true
        \t\t}
        \t\treturn zero, false
        \tcase nodestore.Node:
        \t\tif rn, ok := raw.(jRawNodeProvider); ok {
        \t\t\treturn any(nodestore.Node(rn.RawNodeId_())).(K), true
        \t\t}
        \t\tif u, ok := raw.(uint32); ok {
        \t\t\treturn any(nodestore.Node(u)).(K), true
        \t\t}
        \t\treturn zero, false
        \tcase jRawNodeProvider:
        \t\tvar u uint32
        \t\tvar haveU bool
        \t\tif rn, ok := raw.(jRawNodeProvider); ok {
        \t\t\tu = rn.RawNodeId_()
        \t\t\thaveU = true
        \t\t} else if nd, ok := raw.(nodestore.Node); ok {
        \t\t\tu = uint32(nd)
        \t\t\thaveU = true
        \t\t} else if rawU, ok := raw.(uint32); ok {
        \t\t\tu = rawU
        \t\t\thaveU = true
        \t\t}
        \t\tif !haveU {
        \t\t\treturn zero, false
        \t\t}
        \t\tif u == 0 {
        \t\t\treturn zero, true
        \t\t}
        \t\tif globalNodeBoxer != nil {
        \t\t\tif b, ok := globalNodeBoxer(u).(K); ok {
        \t\t\t\treturn b, true
        \t\t\t}
        \t\t}
        \t\trt := reflect.TypeFor[K]()
        \t\tif rt.Kind() == reflect.Uint32 {
        \t\t\treturn reflect.ValueOf(u).Convert(rt).Interface().(K), true
        \t\t}
        \t\treturn zero, false
        \tdefault:
        \t\treturn zero, false
        \t}
        }

        func OrderedMapGetAny[K comparable, V any](m *OrderedMap[K, V], rawKey any) any {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tk, ok := tryUnboxOrderedMapKey[K](rawKey)
        \tif !ok {
        \t\treturn nil
        \t}
        \tv, found := m.GetOk(k)
        \tif !found {
        \t\treturn nil
        \t}
        \treturn boxJavaElem(v)
        }
        """;
  }

  private static String buildColRtCustomOrderedMapHelpers() {
    return """

        // CustomOrderedMap is a slim 64-byte struct preserving Java LinkedHashMap/HashMap order
        // for FastCustomKey types (such as *Var, AbstractVar, ScopedName) with lazy slice and bucket allocation.
        type FastCustomKeyComparable interface {
        \tcomparable
        \tFastCustomKey
        }

        type CustomMapEntry[K FastCustomKeyComparable, V any] struct {
        \tKey      K
        \tVal      V
        \tHash     int32
        \tHashNext int32
        }

        type CustomOrderedMap[K FastCustomKeyComparable, V any] struct {
        \tentries  []CustomMapEntry[K, V]
        \tbuckets  []int32
        \tnullSlot int32
        \tholes    int32
        \tfirstIdx int32
        }

        func isCustomKeyNull[K FastCustomKeyComparable](k K) bool {
        \tvar zero K
        \treturn k == zero || k.IsJavaNil_()
        }

        const (
        \tcustMapPoolSmallMax = 2048
        \tcustMapPoolMedMax   = 256
        \tcustMapPoolLargeMax = 2
        \tcustMapPoolHugeMax  = 2
        )

        type customOrderedMapBuf[K FastCustomKeyComparable, V any] struct {
        \tentries []CustomMapEntry[K, V]
        \tbuckets []int32
        }

        type customOrderedMapTypePool[K FastCustomKeyComparable, V any] struct {
        \tmu    sync.Mutex
        \tsmall []customOrderedMapBuf[K, V]
        \tmed   []customOrderedMapBuf[K, V]
        \tlarge []customOrderedMapBuf[K, V]
        \thuge  []customOrderedMapBuf[K, V]
        }

        var customOrderedMapPools sync.Map

        func getCustomOrderedMapPool[K FastCustomKeyComparable, V any]() *customOrderedMapTypePool[K, V] {
        \trt := reflect.TypeFor[customOrderedMapBuf[K, V]]()
        \tif v, ok := customOrderedMapPools.Load(rt); ok {
        \t\treturn v.(*customOrderedMapTypePool[K, V])
        \t}
        \tp := &customOrderedMapTypePool[K, V]{}
        \tactual, loaded := customOrderedMapPools.LoadOrStore(rt, p)
        \tif !loaded {
        \t\tpurge := func() {
        \t\t\tp.mu.Lock()
        \t\t\tp.small = nil
        \t\t\tp.med = nil
        \t\t\tp.large = nil
        \t\t\tp.huge = nil
        \t\t\tp.mu.Unlock()
        \t\t}
        \t\tregisterCompPoolReset(purge)
        \t\tregisterRuntimePoolPurger(purge)
        \t}
        \treturn actual.(*customOrderedMapTypePool[K, V])
        }

        func AcquireCustomOrderedMap[K FastCustomKeyComparable, V any]() CustomOrderedMap[K, V] {
        \tp := getCustomOrderedMapPool[K, V]()
        \tp.mu.Lock()
        \tif n := len(p.small); n > 0 {
        \t\tb := p.small[n-1]
        \t\tp.small[n-1] = customOrderedMapBuf[K, V]{}
        \t\tp.small = p.small[:n-1]
        \t\tp.mu.Unlock()
        \t\treturn CustomOrderedMap[K, V]{entries: b.entries[:0]}
        \t}
        \tp.mu.Unlock()
        \treturn CustomOrderedMap[K, V]{}
        }

        func AcquireCustomOrderedMapCap[K FastCustomKeyComparable, V any](capHint int32) CustomOrderedMap[K, V] {
        \tif capHint < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", capHint)))
        \t}
        \tif capHint <= 16 {
        \t\tm := AcquireCustomOrderedMap[K, V]()
        \t\tif cap(m.entries) == 0 && capHint > 8 {
        \t\t\tm.entries = make([]CustomMapEntry[K, V], 0, capHint)
        \t\t}
        \t\treturn m
        \t}
        \tp := getCustomOrderedMapPool[K, V]()
        \tif capHint <= 512 {
        \t\tp.mu.Lock()
        \t\tif n := len(p.med); n > 0 {
        \t\t\tb := p.med[n-1]
        \t\t\tp.med[n-1] = customOrderedMapBuf[K, V]{}
        \t\t\tp.med = p.med[:n-1]
        \t\t\tp.mu.Unlock()
        \t\t\treturn CustomOrderedMap[K, V]{entries: b.entries[:0], buckets: b.buckets}
        \t\t}
        \t\tp.mu.Unlock()
        \t\treturn NewCustomOrderedMapCap[K, V](capHint)
        \t}
        \tp.mu.Lock()
        \tif capHint >= 1024 {
        \t\tif n := len(p.huge); n > 0 {
        \t\t\tb := p.huge[n-1]
        \t\t\tp.huge[n-1] = customOrderedMapBuf[K, V]{}
        \t\t\tp.huge = p.huge[:n-1]
        \t\t\tp.mu.Unlock()
        \t\t\treturn CustomOrderedMap[K, V]{entries: b.entries[:0], buckets: b.buckets}
        \t\t}
        \t}
        \tif n := len(p.large); n > 0 {
        \t\tb := p.large[n-1]
        \t\tp.large[n-1] = customOrderedMapBuf[K, V]{}
        \t\tp.large = p.large[:n-1]
        \t\tp.mu.Unlock()
        \t\treturn CustomOrderedMap[K, V]{entries: b.entries[:0], buckets: b.buckets}
        \t}
        \tp.mu.Unlock()
        \treturn NewCustomOrderedMapCap[K, V](capHint)
        }

        func AcquireCustomOrderedMapPtr[K FastCustomKeyComparable, V any]() *CustomOrderedMap[K, V] {
        \tm := AcquireCustomOrderedMap[K, V]()
        \treturn &m
        }

        func AcquireCustomOrderedMapPtrCap[K FastCustomKeyComparable, V any](capHint int32) *CustomOrderedMap[K, V] {
        \tm := AcquireCustomOrderedMapCap[K, V](capHint)
        \treturn &m
        }

        func ReleaseCustomOrderedMap[K FastCustomKeyComparable, V any](m *CustomOrderedMap[K, V]) {
        \tif m != nil {
        \t\tm.ReleaseToPool()
        \t}
        }

        func ReleaseCustomOrderedMapPtr[K FastCustomKeyComparable, V any](m *CustomOrderedMap[K, V]) {
        \tif m != nil {
        \t\tm.ReleaseToPool()
        \t}
        }

        func (m *CustomOrderedMap[K, V]) ReleaseToPool() {
        \tif m == nil {
        \t\treturn
        \t}
        \tes := m.entries
        \tbk := m.buckets
        \tnullSlot := m.nullSlot
        \tfirstIdx := m.firstIdx
        \t*m = CustomOrderedMap[K, V]{}
        \tc := cap(es)
        \tbl := len(bk)
        \tbc := cap(bk)
        \tif c == 0 && bc == 0 {
        \t\treturn
        \t}
        \tif c > 4194304 || bl > 8388608 {
        \t\tif c >= 64 && c <= pooledSliceMaxCap {
        \t\t\tReleasePooledSlice(es)
        \t\t}
        \t\tif bc >= 64 && bc <= pooledSliceMaxCap {
        \t\t\tReleasePooledSlice(bk)
        \t\t}
        \t\treturn
        \t}
        \tn := len(es)
        \tif bk != nil && n > 0 {
        \t\tif n*4 < bl {
        \t\t\tmask := uint32(bl - 1)
        \t\t\tfor i := int(firstIdx); i < n; i++ {
        \t\t\t\te := &es[i]
        \t\t\t\tif e.HashNext >= 0 && int32(i+1) != nullSlot {
        \t\t\t\t\tbk[smearHash(e.Hash)&mask] = 0
        \t\t\t\t}
        \t\t\t}
        \t\t} else {
        \t\t\tclear(bk)
        \t\t}
        \t}
        \tif n > 0 {
        \t\tclear(es)
        \t\tes = es[:0]
        \t}
        \tp := getCustomOrderedMapPool[K, V]()
        \tkept := false
        \tif bk == nil && c <= 32 {
        \t\tp.mu.Lock()
        \t\tif len(p.small) < custMapPoolSmallMax {
        \t\t\tp.small = append(p.small, customOrderedMapBuf[K, V]{entries: es})
        \t\t\tkept = true
        \t\t}
        \t\tp.mu.Unlock()
        \t} else if c > 262144 || bl > 524288 {
        \t\tp.mu.Lock()
        \t\tif len(p.huge) < custMapPoolHugeMax {
        \t\t\tp.huge = append(p.huge, customOrderedMapBuf[K, V]{entries: es, buckets: bk})
        \t\t\tkept = true
        \t\t}
        \t\tp.mu.Unlock()
        \t} else if c <= 512 && bl <= 1024 {
        \t\tp.mu.Lock()
        \t\tif len(p.med) < custMapPoolMedMax {
        \t\t\tp.med = append(p.med, customOrderedMapBuf[K, V]{entries: es, buckets: bk})
        \t\t\tkept = true
        \t\t}
        \t\tp.mu.Unlock()
        \t} else {
        \t\tp.mu.Lock()
        \t\tif len(p.large) < custMapPoolLargeMax {
        \t\t\tp.large = append(p.large, customOrderedMapBuf[K, V]{entries: es, buckets: bk})
        \t\t\tkept = true
        \t\t}
        \t\tp.mu.Unlock()
        \t}
        \tif !kept {
        \t\tif c >= 64 && c <= pooledSliceMaxCap {
        \t\t\tReleasePooledSlice(es)
        \t\t}
        \t\tif bc >= 64 && bc <= pooledSliceMaxCap {
        \t\t\tReleasePooledSlice(bk)
        \t\t}
        \t}
        }

        func NewCustomOrderedMap[K FastCustomKeyComparable, V any]() CustomOrderedMap[K, V] {
        \treturn CustomOrderedMap[K, V]{}
        }

        func NewCustomOrderedMapPtr[K FastCustomKeyComparable, V any]() *CustomOrderedMap[K, V] {
        \treturn &CustomOrderedMap[K, V]{}
        }

        func NewCustomOrderedMapCap[K FastCustomKeyComparable, V any](cap int32) CustomOrderedMap[K, V] {
        \tif cap < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", cap)))
        \t}
        \tvar m CustomOrderedMap[K, V]
        \tif cap >= 64 {
        \t\tm.entries = AcquirePooledSlice[CustomMapEntry[K, V]](int(cap))
        \t} else if cap > 8 {
        \t\tm.entries = make([]CustomMapEntry[K, V], 0, cap)
        \t}
        \treturn m
        }

        func NewCustomOrderedMapPtrCap[K FastCustomKeyComparable, V any](cap int32) *CustomOrderedMap[K, V] {
        \tif cap < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", cap)))
        \t}
        \tm := &CustomOrderedMap[K, V]{}
        \tif cap >= 64 {
        \t\tm.entries = AcquirePooledSlice[CustomMapEntry[K, V]](int(cap))
        \t} else if cap > 8 {
        \t\tm.entries = make([]CustomMapEntry[K, V], 0, cap)
        \t}
        \treturn m
        }

        func CustomOrderedMapFromPtr[K FastCustomKeyComparable, V any](src *CustomOrderedMap[K, V]) CustomOrderedMap[K, V] {
        \tif src == nil {
        \t\tThrowNPE()
        \t}
        \tlive := int32(len(src.entries)) - src.holes
        \tif live == 0 {
        \t\treturn CustomOrderedMap[K, V]{}
        \t}
        \tcapHint := int(live)
        \tif capHint < 4 {
        \t\tcapHint = 4
        \t}
        \tvar es []CustomMapEntry[K, V]
        \tif capHint >= 64 {
        \t\tes = AcquirePooledSlice[CustomMapEntry[K, V]](capHint)[:live]
        \t} else {
        \t\tes = make([]CustomMapEntry[K, V], live, capHint)
        \t}
        \tif src.holes == 0 {
        \t\tcopy(es, src.entries)
        \t\tdst := CustomOrderedMap[K, V]{entries: es, nullSlot: src.nullSlot}
        \t\tif src.buckets != nil && len(es) > orderedMapIndexThreshold {
        \t\t\tvar b []int32
        \t\t\tif len(src.buckets) >= 64 {
        \t\t\t\tb = AcquirePooledSlice[int32](len(src.buckets))[:len(src.buckets)]
        \t\t\t} else {
        \t\t\t\tb = make([]int32, len(src.buckets))
        \t\t\t}
        \t\t\tcopy(b, src.buckets)
        \t\t\tdst.buckets = b
        \t\t} else if len(es) <= orderedMapIndexThreshold {
        \t\t\tfor i := range es {
        \t\t\t\tes[i].HashNext = 0
        \t\t\t}
        \t\t}
        \t\treturn dst
        \t}
        \tdst := CustomOrderedMap[K, V]{entries: es}
        \tw := 0
        \tfor i := int(src.firstIdx); i < len(src.entries); i++ {
        \t\te := src.entries[i]
        \t\tif e.HashNext >= 0 {
        \t\t\te.HashNext = 0
        \t\t\tif int32(i+1) == src.nullSlot {
        \t\t\t\tdst.nullSlot = int32(w + 1)
        \t\t\t}
        \t\t\tes[w] = e
        \t\t\tw++
        \t\t}
        \t}
        \tif len(es) > orderedMapIndexThreshold {
        \t\tbCap := 16
        \t\tfor (bCap*3)/4 < len(es) {
        \t\t\tbCap *= 2
        \t\t}
        \t\tdst.rehash(bCap)
        \t}
        \treturn dst
        }

        func CustomOrderedMapFrom[K FastCustomKeyComparable, V any](src CustomOrderedMap[K, V]) CustomOrderedMap[K, V] {
        \treturn CustomOrderedMapFromPtr(&src)
        }

        func CustomOrderedMapPtrFrom[K FastCustomKeyComparable, V any](src *CustomOrderedMap[K, V]) *CustomOrderedMap[K, V] {
        \tdst := CustomOrderedMapFromPtr(src)
        \treturn &dst
        }

        func CustomOrderedMapFromJavaMap[K FastCustomKeyComparable, V any](src any) CustomOrderedMap[K, V] {
        \tCheckNotNull(src)
        \tmp := src.(Ext_java_util_Map)
        \tm := NewCustomOrderedMapCap[K, V](mp.Size())
        \tfor cur := NewMapEntryCursor(mp); cur.NextEntry(); {
        \t\tm.PutOk(unboxJavaElem[K](cur.GetKey()), unboxJavaElem[V](cur.GetValue()))
        \t}
        \treturn m
        }

        func CustomOrderedMapPtrFromJavaMap[K FastCustomKeyComparable, V any](src any) *CustomOrderedMap[K, V] {
        \tm := CustomOrderedMapFromJavaMap[K, V](src)
        \treturn &m
        }

        func (m *CustomOrderedMap[K, V]) IsNil() bool     { return m == nil }
        func (m *CustomOrderedMap[K, V]) IsJavaNil_() bool { return m == nil }

        func (m *CustomOrderedMap[K, V]) Size() int32 {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \treturn int32(len(m.entries)) - m.holes
        }

        func (m *CustomOrderedMap[K, V]) IsEmpty() bool {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \treturn int32(len(m.entries)) == m.holes
        }

        func (m *CustomOrderedMap[K, V]) GetOk(k K) (V, bool) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tvar zero V
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn zero, false
        \t}
        \tif isCustomKeyNull(k) {
        \t\tif m.nullSlot == 0 {
        \t\t\treturn zero, false
        \t\t}
        \t\treturn m.entries[m.nullSlot-1].Val, true
        \t}
        \tif m.buckets != nil {
        \t\th := k.HashCode()
        \t\tslot := smearHash(h) & uint32(len(m.buckets)-1)
        \t\tfor pos := m.buckets[slot]; pos != 0; {
        \t\t\te := &m.entries[pos-1]
        \t\t\tif e.Hash == h && (e.Key == k || e.Key.Equals(k)) {
        \t\t\t\treturn e.Val, true
        \t\t\t}
        \t\t\tpos = e.HashNext
        \t\t}
        \t\treturn zero, false
        \t}
        \tfor i := range m.entries {
        \t\tif m.entries[i].Key == k && m.entries[i].HashNext >= 0 && int32(i+1) != m.nullSlot {
        \t\t\treturn m.entries[i].Val, true
        \t\t}
        \t}
        \tfor i := range m.entries {
        \t\te := &m.entries[i]
        \t\tif e.HashNext >= 0 && int32(i+1) != m.nullSlot && e.Key.Equals(k) {
        \t\t\treturn e.Val, true
        \t\t}
        \t}
        \treturn zero, false
        }

        func (m *CustomOrderedMap[K, V]) Get(k K) V {
        \tv, _ := m.GetOk(k)
        \treturn v
        }

        func (m *CustomOrderedMap[K, V]) GetOrDefault(k K, def V) V {
        \tif v, ok := m.GetOk(k); ok {
        \t\treturn v
        \t}
        \treturn def
        }

        func (m *CustomOrderedMap[K, V]) ContainsKey(k K) bool {
        \t_, ok := m.GetOk(k)
        \treturn ok
        }

        func (m *CustomOrderedMap[K, V]) rehash(newCap int) {
        \toldBk := m.buckets
        \tneedHash := oldBk == nil
        \tvar b []int32
        \tif newCap >= 64 {
        \t\tb = AcquirePooledSlice[int32](newCap)[:newCap]
        \t} else {
        \t\tb = make([]int32, newCap)
        \t}
        \tmask := uint32(newCap - 1)
        \tfor i := int(m.firstIdx); i < len(m.entries); i++ {
        \t\te := &m.entries[i]
        \t\tif e.HashNext < 0 {
        \t\t\tcontinue
        \t\t}
        \t\tif int32(i+1) == m.nullSlot {
        \t\t\te.HashNext = 0
        \t\t\tcontinue
        \t\t}
        \t\tif needHash {
        \t\t\te.Hash = e.Key.HashCode()
        \t\t}
        \t\tslot := smearHash(e.Hash) & mask
        \t\te.HashNext = b[slot]
        \t\tb[slot] = int32(i + 1)
        \t}
        \tm.buckets = b
        \tif cap(oldBk) >= 64 {
        \t\tReleasePooledSlice(oldBk)
        \t}
        }

        func (m *CustomOrderedMap[K, V]) appendEntry(e CustomMapEntry[K, V]) {
        \tes := m.entries
        \tif len(es) < cap(es) {
        \t\tm.entries = append(es, e)
        \t\treturn
        \t}
        \tm.appendEntryGrow(e)
        }

        func (m *CustomOrderedMap[K, V]) appendEntryGrow(e CustomMapEntry[K, V]) {
        \told := m.entries
        \toldLen := len(old)
        \toldCap := cap(old)
        \tnewCap := oldCap * 2
        \tif newCap < 4 {
        \t\tnewCap = 4
        \t}
        \tif newCap <= oldLen {
        \t\tnewCap = oldLen + 1
        \t}
        \tif newCap < 64 {
        \t\tnext := make([]CustomMapEntry[K, V], oldLen+1, newCap)
        \t\tcopy(next, old)
        \t\tnext[oldLen] = e
        \t\tm.entries = next
        \t\treturn
        \t}
        \tnext := AcquirePooledSlice[CustomMapEntry[K, V]](newCap)
        \tnext = next[:oldLen+1]
        \tcopy(next, old)
        \tnext[oldLen] = e
        \tif oldCap >= 64 {
        \t\tReleasePooledSlice(old)
        \t}
        \tm.entries = next
        }

        func (m *CustomOrderedMap[K, V]) PutOk(k K, v V) (V, bool) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tvar zeroK K
        \tvar zeroV V
        \tif isCustomKeyNull(k) {
        \t\tif m.nullSlot != 0 {
        \t\t\told := m.entries[m.nullSlot-1].Val
        \t\t\tm.entries[m.nullSlot-1].Val = v
        \t\t\treturn old, true
        \t\t}
        \t\tne := CustomMapEntry[K, V]{Key: zeroK, Val: v}
        \t\tif m.entries == nil {
        \t\t\tes := make([]CustomMapEntry[K, V], 1, 4)
        \t\t\tes[0] = ne
        \t\t\tm.entries = es
        \t\t\tm.nullSlot = 1
        \t\t\treturn zeroV, false
        \t\t}
        \t\tm.appendEntry(ne)
        \t\tm.nullSlot = int32(len(m.entries))
        \t\tif m.buckets == nil && len(m.entries) > orderedMapIndexThreshold {
        \t\t\tm.rehash(32)
        \t\t}
        \t\treturn zeroV, false
        \t}
        \tif m.buckets != nil {
        \t\th := k.HashCode()
        \t\tslot := smearHash(h) & uint32(len(m.buckets)-1)
        \t\tfor pos := m.buckets[slot]; pos != 0; {
        \t\t\te := &m.entries[pos-1]
        \t\t\tif e.Hash == h && (e.Key == k || e.Key.Equals(k)) {
        \t\t\t\told := e.Val
        \t\t\t\te.Val = v
        \t\t\t\treturn old, true
        \t\t\t}
        \t\t\tpos = e.HashNext
        \t\t}
        \t\tm.appendEntry(CustomMapEntry[K, V]{Key: k, Val: v, Hash: h})
        \t\tlive := int32(len(m.entries)) - m.holes
        \t\tif int(live) > (len(m.buckets)*3)/4 {
        \t\t\tm.rehash(len(m.buckets) * 2)
        \t\t} else {
        \t\t\tpos := int32(len(m.entries))
        \t\t\tm.entries[pos-1].HashNext = m.buckets[slot]
        \t\t\tm.buckets[slot] = pos
        \t\t}
        \t\treturn zeroV, false
        \t}
        \tif m.entries == nil {
        \t\tes := make([]CustomMapEntry[K, V], 1, 4)
        \t\tes[0] = CustomMapEntry[K, V]{Key: k, Val: v}
        \t\tm.entries = es
        \t\treturn zeroV, false
        \t}
        \tfor i := range m.entries {
        \t\tif m.entries[i].Key == k && m.entries[i].HashNext >= 0 && int32(i+1) != m.nullSlot {
        \t\t\told := m.entries[i].Val
        \t\t\tm.entries[i].Val = v
        \t\t\treturn old, true
        \t\t}
        \t}
        \tfor i := range m.entries {
        \t\te := &m.entries[i]
        \t\tif e.HashNext >= 0 && int32(i+1) != m.nullSlot && e.Key.Equals(k) {
        \t\t\told := e.Val
        \t\t\te.Val = v
        \t\t\treturn old, true
        \t\t}
        \t}
        \tm.appendEntry(CustomMapEntry[K, V]{Key: k, Val: v})
        \tif len(m.entries) > orderedMapIndexThreshold {
        \t\tm.rehash(32)
        \t}
        \treturn zeroV, false
        }

        func (m *CustomOrderedMap[K, V]) Put(k K, v V) V {
        \told, _ := m.PutOk(k, v)
        \treturn old
        }

        func (m *CustomOrderedMap[K, V]) PutIfAbsent(k K, v V) V {
        \tcur, ok := m.GetOk(k)
        \tif !ok || isJavaNullMapVal(cur) {
        \t\tm.PutOk(k, v)
        \t}
        \treturn cur
        }

        func (m *CustomOrderedMap[K, V]) finishBucketedRemove(slotIdx int32) {
        \tm.holes++
        \tlive := int32(len(m.entries)) - m.holes
        \tif live == 0 {
        \t\tm.entries = m.entries[:0]
        \t\tif len(m.buckets) > 8388608 {
        \t\t\tm.buckets = nil
        \t\t}
        \t\tm.nullSlot = 0
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\treturn
        \t}
        \tif slotIdx == m.firstIdx {
        \t\tfor int(m.firstIdx) < len(m.entries) && m.entries[m.firstIdx].HashNext < 0 {
        \t\t\tm.firstIdx++
        \t\t}
        \t}
        \tif len(m.entries) > 16 && int(live)*2 < len(m.entries) {
        \t\tm.compact()
        \t}
        }

        func (m *CustomOrderedMap[K, V]) RemoveOk(k K) (V, bool) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tvar zeroV V
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn zeroV, false
        \t}
        \tif isCustomKeyNull(k) {
        \t\tif m.nullSlot == 0 {
        \t\t\treturn zeroV, false
        \t\t}
        \t\tpos := m.nullSlot
        \t\tm.nullSlot = 0
        \t\tif m.buckets == nil {
        \t\t\ti := int(pos - 1)
        \t\t\told := m.entries[i].Val
        \t\t\tcopy(m.entries[i:], m.entries[i+1:])
        \t\t\tlast := len(m.entries) - 1
        \t\t\tm.entries[last] = CustomMapEntry[K, V]{}
        \t\t\tm.entries = m.entries[:last]
        \t\t\treturn old, true
        \t\t}
        \t\told := m.entries[pos-1].Val
        \t\tm.entries[pos-1] = CustomMapEntry[K, V]{HashNext: -1}
        \t\tm.finishBucketedRemove(pos - 1)
        \t\treturn old, true
        \t}
        \tif m.buckets == nil {
        \t\tfor i := range m.entries {
        \t\t\tif m.entries[i].Key == k && m.entries[i].HashNext >= 0 && int32(i+1) != m.nullSlot {
        \t\t\t\told := m.entries[i].Val
        \t\t\t\tcopy(m.entries[i:], m.entries[i+1:])
        \t\t\t\tlast := len(m.entries) - 1
        \t\t\t\tm.entries[last] = CustomMapEntry[K, V]{}
        \t\t\t\tm.entries = m.entries[:last]
        \t\t\t\tif m.nullSlot > int32(i+1) {
        \t\t\t\t\tm.nullSlot--
        \t\t\t\t}
        \t\t\t\treturn old, true
        \t\t\t}
        \t\t}
        \t\tfor i := range m.entries {
        \t\t\te := &m.entries[i]
        \t\t\tif e.HashNext >= 0 && int32(i+1) != m.nullSlot && e.Key.Equals(k) {
        \t\t\t\told := e.Val
        \t\t\t\tcopy(m.entries[i:], m.entries[i+1:])
        \t\t\t\tlast := len(m.entries) - 1
        \t\t\t\tm.entries[last] = CustomMapEntry[K, V]{}
        \t\t\t\tm.entries = m.entries[:last]
        \t\t\t\tif m.nullSlot > int32(i+1) {
        \t\t\t\t\tm.nullSlot--
        \t\t\t\t}
        \t\t\t\treturn old, true
        \t\t\t}
        \t\t}
        \t\treturn zeroV, false
        \t}
        \th := k.HashCode()
        \tslot := smearHash(h) & uint32(len(m.buckets)-1)
        \tvar prev int32
        \tfor pos := m.buckets[slot]; pos != 0; {
        \t\te := &m.entries[pos-1]
        \t\tnext := e.HashNext
        \t\tif e.Hash == h && (e.Key == k || e.Key.Equals(k)) {
        \t\t\tif prev == 0 {
        \t\t\t\tm.buckets[slot] = next
        \t\t\t} else {
        \t\t\t\tm.entries[prev-1].HashNext = next
        \t\t\t}
        \t\t\told := e.Val
        \t\t\tm.entries[pos-1] = CustomMapEntry[K, V]{HashNext: -1}
        \t\t\tm.finishBucketedRemove(pos - 1)
        \t\t\treturn old, true
        \t\t}
        \t\tprev = pos
        \t\tpos = next
        \t}
        \treturn zeroV, false
        }

        func (m *CustomOrderedMap[K, V]) Remove(k K) V {
        \told, _ := m.RemoveOk(k)
        \treturn old
        }

        func (m *CustomOrderedMap[K, V]) compact() {
        \tif m.buckets != nil {
        \t\tclear(m.buckets)
        \t}
        \tmask := uint32(len(m.buckets) - 1)
        \toldNull := m.nullSlot
        \tvar newNull int32
        \tw := 0
        \tfor i := int(m.firstIdx); i < len(m.entries); i++ {
        \t\te := m.entries[i]
        \t\tif e.HashNext < 0 {
        \t\t\tcontinue
        \t\t}
        \t\te.HashNext = 0
        \t\tif int32(i+1) == oldNull {
        \t\t\tnewNull = int32(w + 1)
        \t\t} else if m.buckets != nil {
        \t\t\tslot := smearHash(e.Hash) & mask
        \t\t\te.HashNext = m.buckets[slot]
        \t\t\tm.buckets[slot] = int32(w + 1)
        \t\t}
        \t\tm.entries[w] = e
        \t\tw++
        \t}
        \tclear(m.entries[w:])
        \tm.entries = m.entries[:w]
        \tm.nullSlot = newNull
        \tm.holes = 0
        \tm.firstIdx = 0
        }

        func (m *CustomOrderedMap[K, V]) Clear() {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tn := len(m.entries)
        \tif n == 0 && m.nullSlot == 0 {
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\treturn
        \t}
        \tif m.buckets != nil {
        \t\tif len(m.buckets) > 8388608 {
        \t\t\tm.buckets = nil
        \t\t} else if n*4 < len(m.buckets) {
        \t\t\tmask := uint32(len(m.buckets) - 1)
        \t\t\tfor i := int(m.firstIdx); i < n; i++ {
        \t\t\t\te := &m.entries[i]
        \t\t\t\tif e.HashNext >= 0 && int32(i+1) != m.nullSlot {
        \t\t\t\t\tm.buckets[smearHash(e.Hash)&mask] = 0
        \t\t\t\t}
        \t\t\t}
        \t\t} else {
        \t\t\tclear(m.buckets)
        \t\t}
        \t}
        \tif n > 0 {
        \t\tclear(m.entries)
        \t\tm.entries = m.entries[:0]
        \t}
        \tm.nullSlot = 0
        \tm.holes = 0
        \tm.firstIdx = 0
        }

        func (m *CustomOrderedMap[K, V]) PutAllMapPtr(other *CustomOrderedMap[K, V]) {
        \tif m == nil || other == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(other.entries)) == other.holes {
        \t\treturn
        \t}
        \tif len(m.entries) == 0 && m.buckets == nil {
        \t\t*m = CustomOrderedMapFromPtr(other)
        \t\treturn
        \t}
        \tfor i := int(other.firstIdx); i < len(other.entries); i++ {
        \t\tif other.entries[i].HashNext >= 0 {
        \t\t\tm.PutOk(other.entries[i].Key, other.entries[i].Val)
        \t\t}
        \t}
        }

        func (m *CustomOrderedMap[K, V]) PutAllMap(other CustomOrderedMap[K, V]) {
        \tm.PutAllMapPtr(&other)
        }

        func (m *CustomOrderedMap[K, V]) PutAllJavaMap(src any) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(src)
        \tmp := src.(Ext_java_util_Map)
        \tfor cur := NewMapEntryCursor(mp); cur.NextEntry(); {
        \t\tm.PutOk(unboxJavaElem[K](cur.GetKey()), unboxJavaElem[V](cur.GetValue()))
        \t}
        }

        func CustomOrderedMapPtrGetStr[K FastCustomKeyComparable](m *CustomOrderedMap[K, string], k K) string {
        \tif v, ok := m.GetOk(k); ok {
        \t\treturn v
        \t}
        \treturn NullString
        }

        func CustomOrderedMapGetStr[K FastCustomKeyComparable](m CustomOrderedMap[K, string], k K) string {
        \treturn CustomOrderedMapPtrGetStr(&m, k)
        }

        func CustomOrderedMapPutStr[K FastCustomKeyComparable](m *CustomOrderedMap[K, string], k K, v string) string {
        \tif old, ok := m.PutOk(k, v); ok {
        \t\treturn old
        \t}
        \treturn NullString
        }

        func CustomOrderedMapRemoveStr[K FastCustomKeyComparable](m *CustomOrderedMap[K, string], k K) string {
        \tif old, ok := m.RemoveOk(k); ok {
        \t\treturn old
        \t}
        \treturn NullString
        }

        func CustomOrderedMapPutIfAbsentStr[K FastCustomKeyComparable](m *CustomOrderedMap[K, string], k K, v string) string {
        \tcur, ok := m.GetOk(k)
        \tif !ok || cur == NullString {
        \t\tm.PutOk(k, v)
        \t\treturn NullString
        \t}
        \treturn cur
        }

        func CustomOrderedMapComputeIfAbsent[K FastCustomKeyComparable, V comparable](m *CustomOrderedMap[K, V], k K, nullV V, fn func(K) V) V {
        \tif m == nil || fn == nil {
        \t\tThrowNPE()
        \t}
        \tif cur, ok := m.GetOk(k); ok && cur != nullV {
        \t\treturn cur
        \t}
        \tnv := fn(k)
        \tif nv != nullV {
        \t\tm.PutOk(k, nv)
        \t}
        \treturn nv
        }

        func CustomOrderedMapComputeIfAbsentFn[K FastCustomKeyComparable, V comparable](m *CustomOrderedMap[K, V], k K, nullV V, fn Ext_java_util_function_Function) V {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(fn)
        \tif cur, ok := m.GetOk(k); ok && cur != nullV {
        \t\treturn cur
        \t}
        \traw := fn.Apply(boxJavaElem(k))
        \tif raw == nil || isJavaNilAny(raw) {
        \t\treturn nullV
        \t}
        \tnv := unboxJavaElem[V](raw)
        \tif nv != nullV {
        \t\tm.PutOk(k, nv)
        \t}
        \treturn nv
        }

        type CustomOrderedMapKeyCursor[K FastCustomKeyComparable, V any] struct {
        \tm   *CustomOrderedMap[K, V]
        \tpos int32
        }

        func (m *CustomOrderedMap[K, V]) IterKeys() CustomOrderedMapKeyCursor[K, V] {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn CustomOrderedMapKeyCursor[K, V]{m: m}
        \t}
        \treturn CustomOrderedMapKeyCursor[K, V]{m: m, pos: m.firstIdx}
        }

        func (c *CustomOrderedMapKeyCursor[K, V]) HasNext() bool {
        \tm := c.m
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn false
        \t}
        \tfor int(c.pos) < len(m.entries) {
        \t\tif m.entries[c.pos].HashNext >= 0 {
        \t\t\treturn true
        \t\t}
        \t\tc.pos++
        \t}
        \treturn false
        }

        func (c *CustomOrderedMapKeyCursor[K, V]) Next() K {
        \tm := c.m
        \ti := c.pos
        \tc.pos = i + 1
        \treturn m.entries[i].Key
        }

        type CustomOrderedMapValCursor[K FastCustomKeyComparable, V any] struct {
        \tm   *CustomOrderedMap[K, V]
        \tpos int32
        }

        func (m *CustomOrderedMap[K, V]) IterValues() CustomOrderedMapValCursor[K, V] {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn CustomOrderedMapValCursor[K, V]{m: m}
        \t}
        \treturn CustomOrderedMapValCursor[K, V]{m: m, pos: m.firstIdx}
        }

        func (c *CustomOrderedMapValCursor[K, V]) HasNext() bool {
        \tm := c.m
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn false
        \t}
        \tfor int(c.pos) < len(m.entries) {
        \t\tif m.entries[c.pos].HashNext >= 0 {
        \t\t\treturn true
        \t\t}
        \t\tc.pos++
        \t}
        \treturn false
        }

        func (c *CustomOrderedMapValCursor[K, V]) Next() V {
        \tm := c.m
        \ti := c.pos
        \tc.pos = i + 1
        \treturn m.entries[i].Val
        }

        type CustomOrderedMapEntryCursor[K FastCustomKeyComparable, V any] struct {
        \tm        *CustomOrderedMap[K, V]
        \tpos      int32
        \tlastSlot int32
        }

        func (m *CustomOrderedMap[K, V]) IterEntries() CustomOrderedMapEntryCursor[K, V] {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn CustomOrderedMapEntryCursor[K, V]{m: m, lastSlot: -1}
        \t}
        \treturn CustomOrderedMapEntryCursor[K, V]{m: m, pos: m.firstIdx, lastSlot: -1}
        }

        func (c *CustomOrderedMapEntryCursor[K, V]) NextEntry() bool {
        \tm := c.m
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn false
        \t}
        \tfor int(c.pos) < len(m.entries) {
        \t\tif m.entries[c.pos].HashNext >= 0 {
        \t\t\tc.lastSlot = c.pos
        \t\t\tc.pos++
        \t\t\treturn true
        \t\t}
        \t\tc.pos++
        \t}
        \treturn false
        }

        func (c *CustomOrderedMapEntryCursor[K, V]) Key() K {
        \treturn c.m.entries[c.lastSlot].Key
        }

        func (c *CustomOrderedMapEntryCursor[K, V]) Value() V {
        \treturn c.m.entries[c.lastSlot].Val
        }

        func (c *CustomOrderedMapEntryCursor[K, V]) SetValue(nv V) V {
        \told := c.m.entries[c.lastSlot].Val
        \tc.m.entries[c.lastSlot].Val = nv
        \treturn old
        }

        func (m *CustomOrderedMap[K, V]) ForEach(fn func(K, V)) {
        \tif m == nil || fn == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn
        \t}
        \tfor i := int(m.firstIdx); i < len(m.entries); i++ {
        \t\tif m.entries[i].HashNext >= 0 {
        \t\t\tfn(m.entries[i].Key, m.entries[i].Val)
        \t\t}
        \t}
        }

        func (m *CustomOrderedMap[K, V]) ForEachBiConsumer(action Ext_java_util_function_BiConsumer) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(action)
        \tm.ForEach(func(k K, v V) {
        \t\taction.Accept(boxJavaElem(k), boxJavaElem(v))
        \t})
        }

        func CustomOrderedMapPtrEquals[K FastCustomKeyComparable, V comparable](a, b *CustomOrderedMap[K, V]) bool {
        \tif a == nil {
        \t\tThrowNPE()
        \t}
        \tif b == nil {
        \t\treturn false
        \t}
        \taLive := int32(len(a.entries)) - a.holes
        \tbLive := int32(len(b.entries)) - b.holes
        \tif aLive != bLive {
        \t\treturn false
        \t}
        \tif aLive == 0 {
        \t\treturn true
        \t}
        \tfor i := int(a.firstIdx); i < len(a.entries); i++ {
        \t\tif a.entries[i].HashNext >= 0 {
        \t\t\tv2, ok := b.GetOk(a.entries[i].Key)
        \t\t\tif !ok || (a.entries[i].Val != v2 && !javaMapValsEqual(a.entries[i].Val, v2, false)) {
        \t\t\t\treturn false
        \t\t\t}
        \t\t}
        \t}
        \treturn true
        }

        func CustomOrderedMapEquals[K FastCustomKeyComparable, V comparable](a, b CustomOrderedMap[K, V]) bool {
        \treturn CustomOrderedMapPtrEquals(&a, &b)
        }

        func (m *CustomOrderedMap[K, V]) EqualsMapPtr(other *CustomOrderedMap[K, V]) bool {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif other == nil {
        \t\treturn false
        \t}
        \tmLive := int32(len(m.entries)) - m.holes
        \toLive := int32(len(other.entries)) - other.holes
        \tif mLive != oLive {
        \t\treturn false
        \t}
        \tif mLive == 0 {
        \t\treturn true
        \t}
        \tfor i := int(m.firstIdx); i < len(m.entries); i++ {
        \t\tif m.entries[i].HashNext >= 0 {
        \t\t\tv2, ok := other.GetOk(m.entries[i].Key)
        \t\t\tif !ok || !javaMapValsEqual(m.entries[i].Val, v2, false) {
        \t\t\t\treturn false
        \t\t\t}
        \t\t}
        \t}
        \treturn true
        }

        func (m *CustomOrderedMap[K, V]) EqualsMap(other CustomOrderedMap[K, V]) bool {
        \treturn m.EqualsMapPtr(&other)
        }

        func (m *CustomOrderedMap[K, V]) HashCodeMap(_ bool) int32 {
        \treturn m.HashCode()
        }

        func (m *CustomOrderedMap[K, V]) HashCode() int32 {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn 0
        \t}
        \tvar h int32
        \tfor i := int(m.firstIdx); i < len(m.entries); i++ {
        \t\te := &m.entries[i]
        \t\tif e.HashNext >= 0 {
        \t\t\tkh := e.Hash
        \t\t\tif m.buckets == nil && int32(i+1) != m.nullSlot {
        \t\t\t\tkh = e.Key.HashCode()
        \t\t\t}
        \t\t\th += kh ^ javalang.ObjectHashCode(boxJavaElem(e.Val))
        \t\t}
        \t}
        \treturn h
        }

        func (m *CustomOrderedMap[K, V]) String() string {
        \tif m == nil {
        \t\treturn "null"
        \t}
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn "{}"
        \t}
        \tvar sb strings.Builder
        \tsb.WriteByte('{')
        \tfirst := true
        \tm.ForEach(func(k K, v V) {
        \t\tif !first {
        \t\t\tsb.WriteString(", ")
        \t\t}
        \t\tfirst = false
        \t\tsb.WriteString(StrVal(boxJavaElem(k)))
        \t\tsb.WriteByte('=')
        \t\tsb.WriteString(StrVal(boxJavaElem(v)))
        \t})
        \tsb.WriteByte('}')
        \treturn sb.String()
        }

        func (m *CustomOrderedMap[K, V]) ToImmutableMap() Ext_com_google_common_collect_ImmutableMap {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tlive := int32(len(m.entries)) - m.holes
        \tif live == 0 {
        \t\treturn emptyImmutableMap_
        \t}
        \tobj := newImmutableMapShell()
        \tmc := &obj.inlineCore
        \tmc.InitExpectedSize(live)
        \tm.ForEach(func(k K, v V) {
        \t\tmc.Put(CheckNotNull(boxJavaElem(k)), CheckNotNull(boxJavaElem(v)))
        \t})
        \tmc.sealImmutable()
        \treturn obj
        }

        func (m *CustomOrderedMap[K, V]) KeysToImmutableSet() Ext_com_google_common_collect_ImmutableSet {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tlive := int32(len(m.entries)) - m.holes
        \tif live == 0 {
        \t\treturn emptyImmutableSet_
        \t}
        \tsc := &JSetCore{}
        \tsc.InitExpectedSize(live)
        \tm.ForEach(func(k K, _ V) {
        \t\tsc.Add(CheckNotNull(boxJavaElem(k)))
        \t})
        \treturn wrapImmutableSetFromSetCoreOwn(sc)
        }

        func (m *CustomOrderedMap[K, V]) ValuesToImmutableList() Ext_com_google_common_collect_ImmutableList {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tlive := int32(len(m.entries)) - m.holes
        \tif live == 0 {
        \t\treturn emptyImmutableList_
        \t}
        \titems := make([]V, 0, live)
        \tm.ForEach(func(_ K, v V) {
        \t\titems = append(items, v)
        \t})
        \treturn ImmListFromSlice(items)
        }

        func (m *CustomOrderedMap[K, V]) ToJavaLinkedHashMap() *Ext_java_util_LinkedHashMap {
        \tif m == nil {
        \t\treturn nil
        \t}
        \tlive := int32(len(m.entries)) - m.holes
        \tobj := &Ext_java_util_LinkedHashMap{}
        \tobj.BindExtSelf_(obj)
        \tmc := &obj.inlineCore
        \tmc.InitExpectedSize(live)
        \tm.ForEach(func(k K, v V) {
        \t\tmc.Put(boxJavaElem(k), boxJavaElem(v))
        \t})
        \treturn obj
        }

        func (m *CustomOrderedMap[K, V]) ToJavaMap(_ bool) Ext_java_util_Map {
        \treturn m.ToJavaLinkedHashMap()
        }

        type jCustMapLiveValues[K FastCustomKeyComparable, V any] struct {
        \tExt_java_util_AbstractCollectionBase
        \tm *CustomOrderedMap[K, V]
        }

        func (v *jCustMapLiveValues[K, V]) IsJavaNil_() bool { return v == nil }
        func (v *jCustMapLiveValues[K, V]) Size() int32 {
        \tif v == nil || v.m == nil {
        \t\tThrowNPE()
        \t}
        \treturn v.m.Size()
        }
        func (v *jCustMapLiveValues[K, V]) IsEmpty() bool {
        \tif v == nil || v.m == nil {
        \t\tThrowNPE()
        \t}
        \treturn v.m.IsEmpty()
        }
        func (v *jCustMapLiveValues[K, V]) Clear() {
        \tif v == nil || v.m == nil {
        \t\tThrowNPE()
        \t}
        \tv.m.Clear()
        }
        func (v *jCustMapLiveValues[K, V]) Contains(x any) bool {
        \tif v == nil || v.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(v.m.entries)) == v.m.holes {
        \t\treturn false
        \t}
        \ttarget := ToAny(x)
        \tfor cur := v.m.IterValues(); cur.HasNext(); {
        \t\tif javalang.ObjectsEquals(boxJavaElem(cur.Next()), target) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }
        func custMapNextActive[K FastCustomKeyComparable, V any](m *CustomOrderedMap[K, V], from int32) int32 {
        \tn := int32(len(m.entries))
        \tif n == m.holes || from < 0 || from >= n {
        \t\treturn -1
        \t}
        \tif m.holes == 0 {
        \t\treturn from
        \t}
        \tfor i := int(from); i < int(n); i++ {
        \t\tif m.entries[i].HashNext >= 0 {
        \t\t\treturn int32(i)
        \t\t}
        \t}
        \treturn -1
        }

        func (v *jCustMapLiveValues[K, V]) iterInit_() (int32, bool) {
        \tif v == nil || v.m == nil {
        \t\tThrowNPE()
        \t}
        \tp := custMapNextActive(v.m, v.m.firstIdx)
        \treturn p, p >= 0
        }
        func (v *jCustMapLiveValues[K, V]) iterNext_(pos int32) (any, int32) {
        \tp := custMapNextActive(v.m, pos)
        \tif p < 0 {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \treturn boxJavaElem(v.m.entries[p].Val), custMapNextActive(v.m, p+1)
        }
        func (v *jCustMapLiveValues[K, V]) Iterator() any {
        \tif v == nil || v.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(v.m.entries)) == v.m.holes {
        \t\treturn emptyMapIter_
        \t}
        \treturn &jCustMapValIter[K, V]{m: v.m, pos: v.m.firstIdx}
        }

        type jCustMapValIter[K FastCustomKeyComparable, V any] struct {
        \tm            *CustomOrderedMap[K, V]
        \tpos          int32
        \tlastKey      K
        \tlastKeyValid bool
        }

        func (*jCustMapValIter[K, V]) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jCustMapValIter[K, V]) IsJavaNil_() bool                       { return it == nil }
        func (it *jCustMapValIter[K, V]) HasNext() bool {
        \tif int32(len(it.m.entries)) == it.m.holes {
        \t\treturn false
        \t}
        \tfor int(it.pos) < len(it.m.entries) {
        \t\tif it.m.entries[it.pos].HashNext >= 0 {
        \t\t\treturn true
        \t\t}
        \t\tit.pos++
        \t}
        \treturn false
        }
        func (it *jCustMapValIter[K, V]) Next() any {
        \tif !it.HasNext() {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \ti := it.pos
        \tit.pos = i + 1
        \tit.lastKey = it.m.entries[i].Key
        \tit.lastKeyValid = true
        \treturn boxJavaElem(it.m.entries[i].Val)
        }
        func (it *jCustMapValIter[K, V]) Remove() {
        \tif !it.lastKeyValid {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \twasShifted := it.m.buckets == nil
        \tit.m.RemoveOk(it.lastKey)
        \tif wasShifted && it.pos > 0 {
        \t\tit.pos--
        \t}
        \tit.lastKeyValid = false
        }

        func (m *CustomOrderedMap[K, V]) LiveValues() Ext_java_util_Collection {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tv := &jCustMapLiveValues[K, V]{m: m}
        \tv.BindExtSelf_(v)
        \treturn v
        }

        type jCustMapLiveKeySet[K FastCustomKeyComparable, V any] struct {
        \tExt_java_util_AbstractSetBase
        \tm *CustomOrderedMap[K, V]
        }

        func tryUnboxCustomMapKey[K FastCustomKeyComparable](raw any) (K, bool) {
        \tvar zero K
        \tif raw == nil {
        \t\treturn zero, true
        \t}
        \tif any(zero) == nil {
        \t\tif isJavaNilAny(raw) {
        \t\t\treturn zero, true
        \t\t}
        \t\tif !reflect.TypeOf(raw).Comparable() {
        \t\t\treturn zero, false
        \t\t}
        \t\tk, ok := raw.(K)
        \t\treturn k, ok
        \t}
        \tif k, ok := raw.(K); ok {
        \t\tif isCustomKeyNull(k) {
        \t\t\treturn zero, true
        \t\t}
        \t\treturn k, true
        \t}
        \tif isJavaNilAny(raw) {
        \t\treturn zero, true
        \t}
        \treturn zero, false
        }

        func (s *jCustMapLiveKeySet[K, V]) IsJavaNil_() bool { return s == nil }
        func (s *jCustMapLiveKeySet[K, V]) Size() int32 {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \treturn s.m.Size()
        }
        func (s *jCustMapLiveKeySet[K, V]) IsEmpty() bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \treturn s.m.IsEmpty()
        }
        func (s *jCustMapLiveKeySet[K, V]) Clear() {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \ts.m.Clear()
        }
        func (s *jCustMapLiveKeySet[K, V]) Contains(x any) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(s.m.entries)) == s.m.holes {
        \t\treturn false
        \t}
        \tk, ok := tryUnboxCustomMapKey[K](x)
        \tif !ok {
        \t\treturn false
        \t}
        \treturn s.m.ContainsKey(k)
        }
        func (s *jCustMapLiveKeySet[K, V]) Remove_1_any(x any) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(s.m.entries)) == s.m.holes {
        \t\treturn false
        \t}
        \tk, ok := tryUnboxCustomMapKey[K](x)
        \tif !ok {
        \t\treturn false
        \t}
        \t_, removed := s.m.RemoveOk(k)
        \treturn removed
        }
        func (s *jCustMapLiveKeySet[K, V]) SetContainsStr_(k string) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif k == NullString {
        \t\treturn s.m.nullSlot != 0
        \t}
        \treturn false
        }
        func (s *jCustMapLiveKeySet[K, V]) SetAddStr_(_ string) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tpanic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
        }
        func (s *jCustMapLiveKeySet[K, V]) SetRemoveStr_(k string) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif k == NullString {
        \t\tvar zero K
        \t\t_, removed := s.m.RemoveOk(zero)
        \t\treturn removed
        \t}
        \treturn false
        }
        func (s *jCustMapLiveKeySet[K, V]) SetContainsNode_(u uint32) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif u == 0 {
        \t\treturn s.m.nullSlot != 0
        \t}
        \tif int32(len(s.m.entries)) == s.m.holes {
        \t\treturn false
        \t}
        \tvar zero K
        \tif any(zero) == nil {
        \t\tif uk, ok := tryUnboxCustomMapKey[K](boxNodeById(u)); ok {
        \t\t\treturn s.m.ContainsKey(uk)
        \t\t}
        \t}
        \treturn false
        }
        func (s *jCustMapLiveKeySet[K, V]) SetAddNode_(_ uint32) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tpanic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
        }
        func (s *jCustMapLiveKeySet[K, V]) SetRemoveNode_(u uint32) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif u == 0 {
        \t\tvar zero K
        \t\t_, removed := s.m.RemoveOk(zero)
        \t\treturn removed
        \t}
        \tif int32(len(s.m.entries)) == s.m.holes {
        \t\treturn false
        \t}
        \tvar zero K
        \tif any(zero) == nil {
        \t\tif uk, ok := tryUnboxCustomMapKey[K](boxNodeById(u)); ok {
        \t\t\t_, removed := s.m.RemoveOk(uk)
        \t\t\treturn removed
        \t\t}
        \t}
        \treturn false
        }
        func (s *jCustMapLiveKeySet[K, V]) iterInit_() (int32, bool) {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tp := custMapNextActive(s.m, s.m.firstIdx)
        \treturn p, p >= 0
        }
        func (s *jCustMapLiveKeySet[K, V]) iterNext_(pos int32) (any, int32) {
        \tp := custMapNextActive(s.m, pos)
        \tif p < 0 {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \treturn boxJavaElem(s.m.entries[p].Key), custMapNextActive(s.m, p+1)
        }
        func (s *jCustMapLiveKeySet[K, V]) Iterator() any {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(s.m.entries)) == s.m.holes {
        \t\treturn emptyMapIter_
        \t}
        \treturn &jCustMapKeyIter[K, V]{m: s.m, pos: s.m.firstIdx}
        }

        type jCustMapKeyIter[K FastCustomKeyComparable, V any] struct {
        \tm            *CustomOrderedMap[K, V]
        \tpos          int32
        \tlastKey      K
        \tlastKeyValid bool
        }

        func (*jCustMapKeyIter[K, V]) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jCustMapKeyIter[K, V]) IsJavaNil_() bool                       { return it == nil }
        func (it *jCustMapKeyIter[K, V]) HasNext() bool {
        \tif int32(len(it.m.entries)) == it.m.holes {
        \t\treturn false
        \t}
        \tfor int(it.pos) < len(it.m.entries) {
        \t\tif it.m.entries[it.pos].HashNext >= 0 {
        \t\t\treturn true
        \t\t}
        \t\tit.pos++
        \t}
        \treturn false
        }
        func (it *jCustMapKeyIter[K, V]) Next() any {
        \tif !it.HasNext() {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \ti := it.pos
        \tit.pos = i + 1
        \tit.lastKey = it.m.entries[i].Key
        \tit.lastKeyValid = true
        \treturn boxJavaElem(it.lastKey)
        }
        func (it *jCustMapKeyIter[K, V]) Remove() {
        \tif !it.lastKeyValid {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \twasShifted := it.m.buckets == nil
        \tit.m.RemoveOk(it.lastKey)
        \tif wasShifted && it.pos > 0 {
        \t\tit.pos--
        \t}
        \tit.lastKeyValid = false
        }

        func (m *CustomOrderedMap[K, V]) LiveKeySet() Ext_java_util_Set {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \ts := &jCustMapLiveKeySet[K, V]{m: m}
        \ts.BindExtSelf_(s)
        \treturn s
        }

        func (m *CustomOrderedMap[K, V]) FirstValue() V {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tfor i := int(m.firstIdx); i < len(m.entries); i++ {
        \t\tif m.entries[i].HashNext >= 0 {
        \t\t\treturn m.entries[i].Val
        \t\t}
        \t}
        \tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        }

        func CustomOrderedMapGetAny[K FastCustomKeyComparable, V any](m *CustomOrderedMap[K, V], rawKey any) any {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(m.entries)) == m.holes {
        \t\treturn nil
        \t}
        \tk, ok := tryUnboxCustomMapKey[K](rawKey)
        \tif !ok {
        \t\treturn nil
        \t}
        \tv, found := m.GetOk(k)
        \tif !found {
        \t\treturn nil
        \t}
        \treturn boxJavaElem(v)
        }
        """;
  }

  private static String buildColRtIdentityAndCustomSetHelpers() {
    return """

        // IdentityOrderedMap is a slim 40-byte insertion-ordered map always using Go == reference identity
        // with lazy slice allocation and orderedMapIndexThreshold = 12 before allocating the Go map index.
        type IdentityOrderedMap[K comparable, V any] struct {
        \traw OrderedMap[K, V]
        }

        func normIdentKey[K comparable](k K) K {
        \tvar z K
        \tif any(z) == nil {
        \t\tif any(k) == nil || isJavaNilAny(k) {
        \t\t\treturn z
        \t\t}
        \t}
        \treturn k
        }

        func NewIdentityOrderedMap[K comparable, V any]() IdentityOrderedMap[K, V] {
        \treturn IdentityOrderedMap[K, V]{}
        }

        func AcquireIdentityOrderedMap[K comparable, V any]() IdentityOrderedMap[K, V] {
        \treturn IdentityOrderedMap[K, V]{raw: AcquireOrderedMap[K, V]()}
        }

        func AcquireIdentityOrderedMapCap[K comparable, V any](capHint int32) IdentityOrderedMap[K, V] {
        \treturn IdentityOrderedMap[K, V]{raw: AcquireOrderedMapCap[K, V](capHint)}
        }

        func AcquireIdentityOrderedMapPtr[K comparable, V any]() *IdentityOrderedMap[K, V] {
        \tm := AcquireIdentityOrderedMap[K, V]()
        \treturn &m
        }

        func AcquireIdentityOrderedMapPtrCap[K comparable, V any](capHint int32) *IdentityOrderedMap[K, V] {
        \tm := AcquireIdentityOrderedMapCap[K, V](capHint)
        \treturn &m
        }

        func (m *IdentityOrderedMap[K, V]) ReleaseToPool() {
        \tif m == nil {
        \t\treturn
        \t}
        \tm.raw.ReleaseToPool()
        }

        func NewIdentityOrderedMapPtr[K comparable, V any]() *IdentityOrderedMap[K, V] {
        \treturn &IdentityOrderedMap[K, V]{}
        }

        func NewIdentityOrderedMapCap[K comparable, V any](cap int32) IdentityOrderedMap[K, V] {
        \treturn IdentityOrderedMap[K, V]{raw: NewOrderedMapCap[K, V](cap)}
        }

        func NewIdentityOrderedMapPtrCap[K comparable, V any](cap int32) *IdentityOrderedMap[K, V] {
        \tm := NewIdentityOrderedMapCap[K, V](cap)
        \treturn &m
        }

        func IdentityOrderedMapFromPtr[K comparable, V any](src *IdentityOrderedMap[K, V]) IdentityOrderedMap[K, V] {
        \tif src == nil {
        \t\tThrowNPE()
        \t}
        \treturn IdentityOrderedMap[K, V]{raw: OrderedMapFromPtr(&src.raw)}
        }

        func IdentityOrderedMapFrom[K comparable, V any](src IdentityOrderedMap[K, V]) IdentityOrderedMap[K, V] {
        \treturn IdentityOrderedMapFromPtr(&src)
        }

        func IdentityOrderedMapPtrFrom[K comparable, V any](src *IdentityOrderedMap[K, V]) *IdentityOrderedMap[K, V] {
        \tdst := IdentityOrderedMapFromPtr(src)
        \treturn &dst
        }

        func IdentityOrderedMapFromJavaMap[K comparable, V any](src any) IdentityOrderedMap[K, V] {
        \tCheckNotNull(src)
        \tmp := src.(Ext_java_util_Map)
        \tm := NewIdentityOrderedMapCap[K, V](mp.Size())
        \tfor cur := NewMapEntryCursor(mp); cur.NextEntry(); {
        \t\tm.PutOk(unboxJavaElem[K](cur.GetKey()), unboxJavaElem[V](cur.GetValue()))
        \t}
        \treturn m
        }

        func IdentityOrderedMapPtrFromJavaMap[K comparable, V any](src any) *IdentityOrderedMap[K, V] {
        \tm := IdentityOrderedMapFromJavaMap[K, V](src)
        \treturn &m
        }

        func (m *IdentityOrderedMap[K, V]) IsNil() bool     { return m == nil }
        func (m *IdentityOrderedMap[K, V]) IsJavaNil_() bool { return m == nil }
        func (m *IdentityOrderedMap[K, V]) Size() int32 {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \treturn m.raw.Size()
        }
        func (m *IdentityOrderedMap[K, V]) IsEmpty() bool {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \treturn m.raw.IsEmpty()
        }
        func (m *IdentityOrderedMap[K, V]) GetOk(k K) (V, bool) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \treturn m.raw.GetOk(normIdentKey(k))
        }
        func (m *IdentityOrderedMap[K, V]) Get(k K) V {
        \tv, _ := m.GetOk(k)
        \treturn v
        }
        func (m *IdentityOrderedMap[K, V]) GetOrDefault(k K, def V) V {
        \tif v, ok := m.GetOk(k); ok {
        \t\treturn v
        \t}
        \treturn def
        }
        func (m *IdentityOrderedMap[K, V]) ContainsKey(k K) bool {
        \t_, ok := m.GetOk(k)
        \treturn ok
        }
        func (m *IdentityOrderedMap[K, V]) PutOk(k K, v V) (V, bool) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \treturn m.raw.PutOk(normIdentKey(k), v)
        }
        func (m *IdentityOrderedMap[K, V]) Put(k K, v V) V {
        \told, _ := m.PutOk(k, v)
        \treturn old
        }
        func (m *IdentityOrderedMap[K, V]) PutIfAbsent(k K, v V) V {
        \tcur, ok := m.GetOk(k)
        \tif !ok || isJavaNullMapVal(cur) {
        \t\tm.PutOk(k, v)
        \t}
        \treturn cur
        }
        func (m *IdentityOrderedMap[K, V]) RemoveOk(k K) (V, bool) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \treturn m.raw.RemoveOk(normIdentKey(k))
        }
        func (m *IdentityOrderedMap[K, V]) Remove(k K) V {
        \told, _ := m.RemoveOk(k)
        \treturn old
        }
        func (m *IdentityOrderedMap[K, V]) Clear() {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tm.raw.Clear()
        }
        func (m *IdentityOrderedMap[K, V]) PutAllMapPtr(other *IdentityOrderedMap[K, V]) {
        \tif m == nil || other == nil {
        \t\tThrowNPE()
        \t}
        \tm.raw.PutAllMapPtr(&other.raw)
        }
        func (m *IdentityOrderedMap[K, V]) PutAllMap(other IdentityOrderedMap[K, V]) {
        \tm.PutAllMapPtr(&other)
        }
        func (m *IdentityOrderedMap[K, V]) PutAllJavaMap(src any) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(src)
        \tmp := src.(Ext_java_util_Map)
        \tfor cur := NewMapEntryCursor(mp); cur.NextEntry(); {
        \t\tm.PutOk(unboxJavaElem[K](cur.GetKey()), unboxJavaElem[V](cur.GetValue()))
        \t}
        }

        func IdentityOrderedMapPtrGetStr[K comparable](m *IdentityOrderedMap[K, string], k K) string {
        \tif v, ok := m.GetOk(k); ok {
        \t\treturn v
        \t}
        \treturn NullString
        }
        func IdentityOrderedMapGetStr[K comparable](m IdentityOrderedMap[K, string], k K) string {
        \treturn IdentityOrderedMapPtrGetStr(&m, k)
        }
        func IdentityOrderedMapPutStr[K comparable](m *IdentityOrderedMap[K, string], k K, v string) string {
        \tif old, ok := m.PutOk(k, v); ok {
        \t\treturn old
        \t}
        \treturn NullString
        }
        func IdentityOrderedMapRemoveStr[K comparable](m *IdentityOrderedMap[K, string], k K) string {
        \tif old, ok := m.RemoveOk(k); ok {
        \t\treturn old
        \t}
        \treturn NullString
        }
        func IdentityOrderedMapPutIfAbsentStr[K comparable](m *IdentityOrderedMap[K, string], k K, v string) string {
        \tcur, ok := m.GetOk(k)
        \tif !ok || cur == NullString {
        \t\tm.PutOk(k, v)
        \t\treturn NullString
        \t}
        \treturn cur
        }
        func IdentityOrderedMapComputeIfAbsent[K comparable, V comparable](m *IdentityOrderedMap[K, V], k K, nullV V, fn func(K) V) V {
        \tif m == nil || fn == nil {
        \t\tThrowNPE()
        \t}
        \tnk := normIdentKey(k)
        \treturn OrderedMapComputeIfAbsent(&m.raw, nk, nullV, fn)
        }
        func IdentityOrderedMapComputeIfAbsentFn[K comparable, V comparable](m *IdentityOrderedMap[K, V], k K, nullV V, fn Ext_java_util_function_Function) V {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tnk := normIdentKey(k)
        \treturn OrderedMapComputeIfAbsentFn(&m.raw, nk, nullV, fn)
        }

        func (m *IdentityOrderedMap[K, V]) IterKeys() OrderedMapKeyCursor[K, V] {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \treturn m.raw.IterKeys()
        }
        func (m *IdentityOrderedMap[K, V]) IterValues() OrderedMapValCursor[K, V] {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \treturn m.raw.IterValues()
        }
        func (m *IdentityOrderedMap[K, V]) IterEntries() OrderedMapEntryCursor[K, V] {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \treturn m.raw.IterEntries()
        }
        func (m *IdentityOrderedMap[K, V]) ForEach(fn func(K, V)) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tm.raw.ForEach(fn)
        }
        func (m *IdentityOrderedMap[K, V]) ForEachBiConsumer(action Ext_java_util_function_BiConsumer) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tm.raw.ForEachBiConsumer(action)
        }
        func IdentityOrderedMapPtrEquals[K comparable, V comparable](a, b *IdentityOrderedMap[K, V]) bool {
        \tif a == nil {
        \t\tThrowNPE()
        \t}
        \tif b == nil {
        \t\treturn false
        \t}
        \treturn OrderedMapPtrEquals(&a.raw, &b.raw, true)
        }
        func IdentityOrderedMapEquals[K comparable, V comparable](a, b IdentityOrderedMap[K, V]) bool {
        \treturn IdentityOrderedMapPtrEquals(&a, &b)
        }
        func (m *IdentityOrderedMap[K, V]) EqualsMapPtr(other *IdentityOrderedMap[K, V]) bool {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif other == nil {
        \t\treturn false
        \t}
        \treturn m.raw.EqualsMapPtr(&other.raw, true)
        }
        func (m *IdentityOrderedMap[K, V]) EqualsMap(other IdentityOrderedMap[K, V]) bool {
        \treturn m.EqualsMapPtr(&other)
        }
        func (m *IdentityOrderedMap[K, V]) HashCode() int32 {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \treturn m.raw.HashCodeMap(true)
        }
        func (m *IdentityOrderedMap[K, V]) HashCodeMap(_ bool) int32 {
        \treturn m.HashCode()
        }
        func (m *IdentityOrderedMap[K, V]) String() string {
        \tif m == nil {
        \t\treturn "null"
        \t}
        \treturn m.raw.String()
        }
        func (m *IdentityOrderedMap[K, V]) ValuesToImmutableList() Ext_com_google_common_collect_ImmutableList {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \treturn m.raw.ValuesToImmutableList()
        }
        func (m *IdentityOrderedMap[K, V]) ToJavaIdentityHashMap() *Ext_java_util_IdentityHashMap {
        \tif m == nil {
        \t\treturn nil
        \t}
        \treturn m.raw.ToJavaMap(true).(*Ext_java_util_IdentityHashMap)
        }
        func (m *IdentityOrderedMap[K, V]) ToJavaMap(_ bool) Ext_java_util_Map {
        \tif m == nil {
        \t\treturn nil
        \t}
        \treturn m.raw.ToJavaMap(true)
        }
        func (m *IdentityOrderedMap[K, V]) LiveValues() Ext_java_util_Collection {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \treturn m.raw.LiveValues()
        }

        type jIdentMapLiveKeySet[K comparable, V any] struct {
        \tExt_java_util_AbstractSetBase
        \tm *IdentityOrderedMap[K, V]
        }

        func (s *jIdentMapLiveKeySet[K, V]) IsJavaNil_() bool { return s == nil }
        func (s *jIdentMapLiveKeySet[K, V]) Size() int32 {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \treturn s.m.Size()
        }
        func (s *jIdentMapLiveKeySet[K, V]) IsEmpty() bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \treturn s.m.IsEmpty()
        }
        func (s *jIdentMapLiveKeySet[K, V]) Clear() {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \ts.m.Clear()
        }
        func (s *jIdentMapLiveKeySet[K, V]) Contains(x any) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(s.m.raw.entries)) == s.m.raw.holes {
        \t\treturn false
        \t}
        \tk, ok := tryUnboxOrderedMapKey[K](x)
        \tif !ok {
        \t\treturn false
        \t}
        \treturn s.m.raw.ContainsKey(normIdentKey(k))
        }
        func (s *jIdentMapLiveKeySet[K, V]) Remove_1_any(x any) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(s.m.raw.entries)) == s.m.raw.holes {
        \t\treturn false
        \t}
        \tk, ok := tryUnboxOrderedMapKey[K](x)
        \tif !ok {
        \t\treturn false
        \t}
        \t_, removed := s.m.raw.RemoveOk(normIdentKey(k))
        \treturn removed
        }
        func (s *jIdentMapLiveKeySet[K, V]) SetContainsStr_(k string) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(s.m.raw.entries)) == s.m.raw.holes {
        \t\treturn false
        \t}
        \tvar zero K
        \tif any(zero) == nil {
        \t\tif k == NullString {
        \t\t\treturn s.m.raw.ContainsKey(zero)
        \t\t}
        \t\tuk, ok := any(k).(K)
        \t\treturn ok && s.m.raw.ContainsKey(uk)
        \t}
        \tif _, isStr := any(zero).(string); isStr {
        \t\treturn s.m.raw.ContainsKey(any(k).(K))
        \t}
        \tif k == NullString {
        \t\tif uk, ok := tryUnboxOrderedMapKey[K](nil); ok {
        \t\t\treturn s.m.raw.ContainsKey(uk)
        \t\t}
        \t}
        \treturn false
        }
        func (s *jIdentMapLiveKeySet[K, V]) SetAddStr_(_ string) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tpanic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
        }
        func (s *jIdentMapLiveKeySet[K, V]) SetRemoveStr_(k string) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(s.m.raw.entries)) == s.m.raw.holes {
        \t\treturn false
        \t}
        \tvar zero K
        \tif any(zero) == nil {
        \t\tif k == NullString {
        \t\t\t_, removed := s.m.raw.RemoveOk(zero)
        \t\t\treturn removed
        \t\t}
        \t\tuk, ok := any(k).(K)
        \t\tif !ok {
        \t\t\treturn false
        \t\t}
        \t\t_, removed := s.m.raw.RemoveOk(uk)
        \t\treturn removed
        \t}
        \tif _, isStr := any(zero).(string); isStr {
        \t\t_, removed := s.m.raw.RemoveOk(any(k).(K))
        \t\treturn removed
        \t}
        \tif k == NullString {
        \t\tif uk, ok := tryUnboxOrderedMapKey[K](nil); ok {
        \t\t\t_, removed := s.m.raw.RemoveOk(uk)
        \t\t\treturn removed
        \t\t}
        \t}
        \treturn false
        }
        func (s *jIdentMapLiveKeySet[K, V]) SetContainsNode_(u uint32) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(s.m.raw.entries)) == s.m.raw.holes {
        \t\treturn false
        \t}
        \tif u == 0 {
        \t\tuk, ok := tryUnboxOrderedMapKey[K](nil)
        \t\treturn ok && s.m.raw.ContainsKey(uk)
        \t}
        \tvar zero K
        \tif any(zero) == nil {
        \t\tuk, ok := boxNodeById(u).(K)
        \t\treturn ok && s.m.raw.ContainsKey(uk)
        \t}
        \tswitch any(zero).(type) {
        \tcase uint32:
        \t\treturn s.m.raw.ContainsKey(any(u).(K))
        \tcase nodestore.Node:
        \t\treturn s.m.raw.ContainsKey(any(nodestore.Node(u)).(K))
        \tcase jRawNodeProvider:
        \t\tif globalNodeBoxer != nil {
        \t\t\tif uk, ok := globalNodeBoxer(u).(K); ok {
        \t\t\t\treturn s.m.raw.ContainsKey(uk)
        \t\t\t}
        \t\t}
        \t\tuk, ok := tryUnboxOrderedMapKey[K](u)
        \t\treturn ok && s.m.raw.ContainsKey(uk)
        \tdefault:
        \t\treturn false
        \t}
        }
        func (s *jIdentMapLiveKeySet[K, V]) SetAddNode_(_ uint32) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tpanic(javalang.NewUnsupportedOperationException("Map.keySet.add"))
        }
        func (s *jIdentMapLiveKeySet[K, V]) SetRemoveNode_(u uint32) bool {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(s.m.raw.entries)) == s.m.raw.holes {
        \t\treturn false
        \t}
        \tif u == 0 {
        \t\tuk, ok := tryUnboxOrderedMapKey[K](nil)
        \t\tif !ok {
        \t\t\treturn false
        \t\t}
        \t\t_, removed := s.m.raw.RemoveOk(uk)
        \t\treturn removed
        \t}
        \tvar zero K
        \tif any(zero) == nil {
        \t\tuk, ok := boxNodeById(u).(K)
        \t\tif !ok {
        \t\t\treturn false
        \t\t}
        \t\t_, removed := s.m.raw.RemoveOk(uk)
        \t\treturn removed
        \t}
        \tswitch any(zero).(type) {
        \tcase uint32:
        \t\t_, removed := s.m.raw.RemoveOk(any(u).(K))
        \t\treturn removed
        \tcase nodestore.Node:
        \t\t_, removed := s.m.raw.RemoveOk(any(nodestore.Node(u)).(K))
        \t\treturn removed
        \tcase jRawNodeProvider:
        \t\tif globalNodeBoxer != nil {
        \t\t\tif uk, ok := globalNodeBoxer(u).(K); ok {
        \t\t\t\t_, removed := s.m.raw.RemoveOk(uk)
        \t\t\t\treturn removed
        \t\t\t}
        \t\t}
        \t\tuk, ok := tryUnboxOrderedMapKey[K](u)
        \t\tif !ok {
        \t\t\treturn false
        \t\t}
        \t\t_, removed := s.m.raw.RemoveOk(uk)
        \t\treturn removed
        \tdefault:
        \t\treturn false
        \t}
        }
        func (s *jIdentMapLiveKeySet[K, V]) iterInit_() (int32, bool) {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tp := ordMapNextActive(&s.m.raw, s.m.raw.firstIdx)
        \treturn p, p >= 0
        }
        func (s *jIdentMapLiveKeySet[K, V]) iterNext_(pos int32) (any, int32) {
        \tp := ordMapNextActive(&s.m.raw, pos)
        \tif p < 0 {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \treturn boxJavaElem(s.m.raw.entries[p].Key), ordMapNextActive(&s.m.raw, p+1)
        }
        func (s *jIdentMapLiveKeySet[K, V]) Iterator() any {
        \tif s == nil || s.m == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(s.m.raw.entries)) == s.m.raw.holes {
        \t\treturn emptyMapIter_
        \t}
        \treturn &jOrdMapKeyIter[K, V]{m: &s.m.raw, pos: s.m.raw.firstIdx}
        }
        func (m *IdentityOrderedMap[K, V]) LiveKeySet() Ext_java_util_Set {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \ts := &jIdentMapLiveKeySet[K, V]{m: m}
        \ts.BindExtSelf_(s)
        \treturn s
        }

        func (m *IdentityOrderedMap[K, V]) FirstValue() V {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \treturn m.raw.FirstValue()
        }

        func IdentityOrderedMapGetAny[K comparable, V any](m *IdentityOrderedMap[K, V], rawKey any) any {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tk, ok := tryUnboxOrderedMapKey[K](rawKey)
        \tif !ok {
        \t\treturn nil
        \t}
        \tv, found := m.raw.GetOk(normIdentKey(k))
        \tif !found {
        \t\treturn nil
        \t}
        \treturn boxJavaElem(v)
        }

        // IdentityOrderedSet is a slim 40-byte insertion-ordered set always using Go == reference identity
        // with lazy slice allocation and orderedSetIndexThreshold = 8 before allocating the Go map index.
        type IdentityOrderedSet[T comparable] struct {
        \traw OrderedSet[T]
        }

        func NewIdentityOrderedSet[T comparable]() IdentityOrderedSet[T] {
        \treturn IdentityOrderedSet[T]{}
        }
        func NewIdentityOrderedSetPtr[T comparable]() *IdentityOrderedSet[T] {
        \treturn &IdentityOrderedSet[T]{}
        }
        func NewIdentityOrderedSetCap[T comparable](cap int32) IdentityOrderedSet[T] {
        \treturn IdentityOrderedSet[T]{raw: NewOrderedSetCap[T](cap)}
        }
        func NewIdentityOrderedSetPtrCap[T comparable](cap int32) *IdentityOrderedSet[T] {
        \ts := NewIdentityOrderedSetCap[T](cap)
        \treturn &s
        }
        func IdentityOrderedSetFromSet[T comparable](src IdentityOrderedSet[T]) IdentityOrderedSet[T] {
        \treturn IdentityOrderedSet[T]{raw: OrderedSetFromSet(src.raw)}
        }
        func IdentityOrderedSetPtrFromSet[T comparable](src *IdentityOrderedSet[T]) *IdentityOrderedSet[T] {
        \tif src == nil {
        \t\tThrowNPE()
        \t}
        \tdst := IdentityOrderedSetFromSet(*src)
        \treturn &dst
        }
        func IdentityOrderedSetFromNodeCol[N ~uint32](col any) IdentityOrderedSet[N] {
        \treturn IdentityOrderedSet[N]{raw: OrderedSetFromNodeCol[N](col)}
        }
        func IdentityOrderedSetPtrFromNodeCol[N ~uint32](col any) *IdentityOrderedSet[N] {
        \ts := IdentityOrderedSetFromNodeCol[N](col)
        \treturn &s
        }
        func IdentityOrderedSetFromCol[T comparable](col any) IdentityOrderedSet[T] {
        \tvar s IdentityOrderedSet[T]
        \tIdentityOrderedSetAddAllCol(&s, col)
        \treturn s
        }
        func IdentityOrderedSetPtrFromCol[T comparable](col any) *IdentityOrderedSet[T] {
        \ts := IdentityOrderedSetFromCol[T](col)
        \treturn &s
        }
        func (s *IdentityOrderedSet[T]) Init() {}
        func (s *IdentityOrderedSet[T]) InitCap(cap int32) {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \ts.raw.InitCap(cap)
        }
        func (s *IdentityOrderedSet[T]) IsNil() bool       { return s == nil }
        func (s *IdentityOrderedSet[T]) IsJavaNil_() bool   { return s == nil }
        func (s IdentityOrderedSet[T]) Size() int32         { return s.raw.Size() }
        func (s IdentityOrderedSet[T]) IsEmpty() bool       { return s.raw.IsEmpty() }
        func (s IdentityOrderedSet[T]) Contains(v T) bool   { return s.raw.Contains(normIdentKey(v)) }
        func (s IdentityOrderedSet[T]) GetFirst() T         { return s.raw.GetFirst() }
        func (s *IdentityOrderedSet[T]) Add(v T) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn s.raw.Add(normIdentKey(v))
        }
        func (s *IdentityOrderedSet[T]) Remove(v T) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn s.raw.Remove(normIdentKey(v))
        }
        func (s *IdentityOrderedSet[T]) Clear() {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \ts.raw.Clear()
        }
        func (s *IdentityOrderedSet[T]) AddAllSet(other IdentityOrderedSet[T]) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn s.raw.AddAllSet(other.raw)
        }
        func (s *IdentityOrderedSet[T]) AddAllSetPtr(other *IdentityOrderedSet[T]) bool {
        \tif s == nil || other == nil {
        \t\tThrowNPE()
        \t}
        \treturn s.raw.AddAllSet(other.raw)
        }
        func IdentityOrderedSetAddAllNode[N ~uint32](s *IdentityOrderedSet[N], col any) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn OrderedSetAddAllNode(&s.raw, col)
        }
        func IdentityOrderedSetAddAllCol[T comparable](s *IdentityOrderedSet[T], col any) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(col)
        \tchanged := false
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\tif s.Add(unboxJavaElem[T](cur.Next())) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }
        func (s *IdentityOrderedSet[T]) Iter() OrderedSetCursor[T] {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn s.raw.Iter()
        }
        func (s IdentityOrderedSet[T]) EqualsSet(other IdentityOrderedSet[T]) bool {
        \treturn s.raw.EqualsSet(other.raw)
        }
        func (s *IdentityOrderedSet[T]) EqualsSetPtr(other *IdentityOrderedSet[T]) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif other == nil {
        \t\treturn false
        \t}
        \treturn s.raw.EqualsSet(other.raw)
        }
        func (s IdentityOrderedSet[T]) HashCode() int32 {
        \tif s.IsEmpty() {
        \t\treturn 0
        \t}
        \tvar h int32
        \tif s.raw.holes > 0 {
        \t\tfor i, v := range s.raw.order {
        \t\t\tif s.raw.idx[v] == int32(i+1) {
        \t\t\t\th += IdentityHash(boxJavaElem(v))
        \t\t\t}
        \t\t}
        \t\treturn h
        \t}
        \tfor _, v := range s.raw.order {
        \t\th += IdentityHash(boxJavaElem(v))
        \t}
        \treturn h
        }
        func (s *IdentityOrderedSet[T]) String() string {
        \tif s == nil {
        \t\treturn "null"
        \t}
        \treturn (&s.raw).String()
        }
        func (s IdentityOrderedSet[T]) ToJavaSet() Ext_java_util_Set {
        \tout := &Ext_java_util_HashSetBase{}
        \tout.BindExtSelf_(out)
        \tout.inlineSet.identity = true
        \tif s.IsEmpty() {
        \t\treturn out
        \t}
        \tfor cur := s.Iter(); cur.HasNext(); {
        \t\tout.Add_1(boxJavaElem(cur.Next()))
        \t}
        \treturn out
        }
        func ArrayDequeFromIdentityOrderedSet[T comparable](s IdentityOrderedSet[T]) *Ext_java_util_ArrayDeque {
        \treturn ArrayDequeFromOrderedSet(s.raw)
        }
        func ArrayDequeFromIdentityOrderedSetPtr[T comparable](s *IdentityOrderedSet[T]) *Ext_java_util_ArrayDeque {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn ArrayDequeFromOrderedSet(s.raw)
        }
        func DequeFromIdentityOrderedSet[T comparable](s IdentityOrderedSet[T]) Deque[T] {
        \treturn DequeFromOrderedSet(s.raw)
        }
        func DequePtrFromIdentityOrderedSet[T comparable](s IdentityOrderedSet[T]) *Deque[T] {
        \treturn DequePtrFromOrderedSet(s.raw)
        }
        func DequeFromIdentityOrderedSetPtr[T comparable](s *IdentityOrderedSet[T]) Deque[T] {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn DequeFromOrderedSet(s.raw)
        }
        func DequePtrFromIdentityOrderedSetPtr[T comparable](s *IdentityOrderedSet[T]) *Deque[T] {
        \td := DequeFromIdentityOrderedSetPtr(s)
        \treturn &d
        }

        // CustomOrderedSet is a 64-byte insertion-ordered set backed by CustomOrderedMap[T, struct{}]
        // for custom .Equals/.HashCode types (JSType, ObjectType, *Var, AbstractVar, ScopedName, etc.)
        // with lazy slice allocation and orderedMapIndexThreshold = 12 before bucket allocation.
        type CustomOrderedSet[T FastCustomKeyComparable] struct {
        \tm CustomOrderedMap[T, struct{}]
        }

        type CustomOrderedSetCursor[T FastCustomKeyComparable] struct {
        \traw CustomOrderedMapKeyCursor[T, struct{}]
        }

        func NewCustomOrderedSet[T FastCustomKeyComparable]() CustomOrderedSet[T] {
        \treturn CustomOrderedSet[T]{}
        }
        func NewCustomOrderedSetPtr[T FastCustomKeyComparable]() *CustomOrderedSet[T] {
        \treturn &CustomOrderedSet[T]{}
        }
        func NewCustomOrderedSetCap[T FastCustomKeyComparable](cap int32) CustomOrderedSet[T] {
        \treturn CustomOrderedSet[T]{m: NewCustomOrderedMapCap[T, struct{}](cap)}
        }
        func NewCustomOrderedSetPtrCap[T FastCustomKeyComparable](cap int32) *CustomOrderedSet[T] {
        \ts := NewCustomOrderedSetCap[T](cap)
        \treturn &s
        }
        func CustomOrderedSetFromSet[T FastCustomKeyComparable](src CustomOrderedSet[T]) CustomOrderedSet[T] {
        \treturn CustomOrderedSet[T]{m: CustomOrderedMapFromPtr(&src.m)}
        }
        func CustomOrderedSetPtrFromSet[T FastCustomKeyComparable](src *CustomOrderedSet[T]) *CustomOrderedSet[T] {
        \tif src == nil {
        \t\tThrowNPE()
        \t}
        \tdst := CustomOrderedSetFromSet(*src)
        \treturn &dst
        }
        func CustomOrderedSetFromCol[T FastCustomKeyComparable](col any) CustomOrderedSet[T] {
        \tvar s CustomOrderedSet[T]
        \tCustomOrderedSetAddAllCol(&s, col)
        \treturn s
        }
        func CustomOrderedSetPtrFromCol[T FastCustomKeyComparable](col any) *CustomOrderedSet[T] {
        \ts := CustomOrderedSetFromCol[T](col)
        \treturn &s
        }
        func (s *CustomOrderedSet[T]) Init() {}
        func (s *CustomOrderedSet[T]) InitCap(cap int32) {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif cap < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", cap)))
        \t}
        \tif cap > 8 && s.m.entries == nil {
        \t\tif cap >= 64 {
        \t\t\ts.m.entries = AcquirePooledSlice[CustomMapEntry[T, struct{}]](int(cap))
        \t\t} else {
        \t\t\ts.m.entries = make([]CustomMapEntry[T, struct{}], 0, cap)
        \t\t}
        \t}
        }
        func (s *CustomOrderedSet[T]) IsNil() bool     { return s == nil }
        func (s *CustomOrderedSet[T]) IsJavaNil_() bool { return s == nil }
        func (s CustomOrderedSet[T]) Size() int32      { return int32(len(s.m.entries)) - s.m.holes }
        func (s CustomOrderedSet[T]) IsEmpty() bool    { return int32(len(s.m.entries)) == s.m.holes }
        func (s CustomOrderedSet[T]) Contains(v T) bool {
        \t_, ok := s.m.GetOk(v)
        \treturn ok
        }
        func (s CustomOrderedSet[T]) GetFirst() T {
        \tif int32(len(s.m.entries)) == s.m.holes {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \tfor i := int(s.m.firstIdx); i < len(s.m.entries); i++ {
        \t\tif s.m.entries[i].HashNext >= 0 {
        \t\t\treturn s.m.entries[i].Key
        \t\t}
        \t}
        \tpanic(javalang.NewNoSuchElementException("empty"))
        }
        func (s *CustomOrderedSet[T]) Add(v T) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \t_, existed := s.m.PutOk(v, struct{}{})
        \treturn !existed
        }
        func (s *CustomOrderedSet[T]) Remove(v T) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \t_, existed := s.m.RemoveOk(v)
        \treturn existed
        }
        func (s *CustomOrderedSet[T]) Clear() {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \ts.m.Clear()
        }
        func (s *CustomOrderedSet[T]) AddAllSet(other CustomOrderedSet[T]) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif int32(len(other.m.entries)) == other.m.holes {
        \t\treturn false
        \t}
        \tchanged := false
        \tfor i := int(other.m.firstIdx); i < len(other.m.entries); i++ {
        \t\tif other.m.entries[i].HashNext >= 0 && s.Add(other.m.entries[i].Key) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }
        func (s *CustomOrderedSet[T]) AddAllSetPtr(other *CustomOrderedSet[T]) bool {
        \tif s == nil || other == nil {
        \t\tThrowNPE()
        \t}
        \treturn s.AddAllSet(*other)
        }
        func CustomOrderedSetAddAllCol[T FastCustomKeyComparable](s *CustomOrderedSet[T], col any) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(col)
        \tchanged := false
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\tif s.Add(unboxJavaElem[T](cur.Next())) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }
        func (s *CustomOrderedSet[T]) Iter() CustomOrderedSetCursor[T] {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn CustomOrderedSetCursor[T]{raw: s.m.IterKeys()}
        }
        func (c *CustomOrderedSetCursor[T]) HasNext() bool { return c.raw.HasNext() }
        func (c *CustomOrderedSetCursor[T]) Next() T       { return c.raw.Next() }

        func (s CustomOrderedSet[T]) EqualsSet(other CustomOrderedSet[T]) bool {
        \tif s.Size() != other.Size() {
        \t\treturn false
        \t}
        \tif s.IsEmpty() {
        \t\treturn true
        \t}
        \tfor i := int(s.m.firstIdx); i < len(s.m.entries); i++ {
        \t\tif s.m.entries[i].HashNext >= 0 && !other.Contains(s.m.entries[i].Key) {
        \t\t\treturn false
        \t\t}
        \t}
        \treturn true
        }
        func (s *CustomOrderedSet[T]) EqualsSetPtr(other *CustomOrderedSet[T]) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif other == nil {
        \t\treturn false
        \t}
        \treturn s.EqualsSet(*other)
        }
        func (s CustomOrderedSet[T]) HashCode() int32 {
        \tif s.IsEmpty() {
        \t\treturn 0
        \t}
        \tvar h int32
        \tfor i := int(s.m.firstIdx); i < len(s.m.entries); i++ {
        \t\te := &s.m.entries[i]
        \t\tif e.HashNext >= 0 && int32(i+1) != s.m.nullSlot {
        \t\t\th += e.Hash
        \t\t}
        \t}
        \treturn h
        }
        func (s *CustomOrderedSet[T]) String() string {
        \tif s == nil {
        \t\treturn "null"
        \t}
        \tif s.IsEmpty() {
        \t\treturn "[]"
        \t}
        \tvar sb strings.Builder
        \tsb.WriteByte('[')
        \tfirst := true
        \tfor i := int(s.m.firstIdx); i < len(s.m.entries); i++ {
        \t\tif s.m.entries[i].HashNext >= 0 {
        \t\t\tif !first {
        \t\t\t\tsb.WriteString(", ")
        \t\t\t}
        \t\t\tfirst = false
        \t\t\tsb.WriteString(StrVal(boxJavaElem(s.m.entries[i].Key)))
        \t\t}
        \t}
        \tsb.WriteByte(']')
        \treturn sb.String()
        }
        func (s CustomOrderedSet[T]) ToImmutableSet() Ext_com_google_common_collect_ImmutableSet {
        \treturn s.m.KeysToImmutableSet()
        }
        func (s CustomOrderedSet[T]) ToJavaSet() *Ext_java_util_LinkedHashSet {
        \tlive := s.Size()
        \tobj := &Ext_java_util_LinkedHashSet{}
        \tobj.BindExtSelf_(obj)
        \tobj.inlineSet.InitExpectedSize(live)
        \tif live == 0 {
        \t\treturn obj
        \t}
        \tfor i := int(s.m.firstIdx); i < len(s.m.entries); i++ {
        \t\tif s.m.entries[i].HashNext >= 0 {
        \t\t\tobj.inlineSet.Add(boxJavaElem(s.m.entries[i].Key))
        \t\t}
        \t}
        \treturn obj
        }
        func ArrayDequeFromCustomOrderedSet[T FastCustomKeyComparable](s CustomOrderedSet[T]) *Ext_java_util_ArrayDeque {
        \tlive := s.Size()
        \tobj := &Ext_java_util_ArrayDeque{}
        \tobj.BindExtSelf_(obj)
        \tdc := &obj.inlineCore
        \tdc.InitCap(live)
        \tif live == 0 {
        \t\treturn obj
        \t}
        \tfor i := int(s.m.firstIdx); i < len(s.m.entries); i++ {
        \t\tif s.m.entries[i].HashNext >= 0 {
        \t\t\tdc.AddLast(boxJavaElem(s.m.entries[i].Key))
        \t\t}
        \t}
        \treturn obj
        }
        func ArrayDequeFromCustomOrderedSetPtr[T FastCustomKeyComparable](s *CustomOrderedSet[T]) *Ext_java_util_ArrayDeque {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn ArrayDequeFromCustomOrderedSet(*s)
        }
        func DequeFromCustomOrderedSet[T FastCustomKeyComparable](s CustomOrderedSet[T]) Deque[T] {
        \tlive := s.Size()
        \tif live == 0 {
        \t\treturn Deque[T]{}
        \t}
        \tvar dst Deque[T]
        \tdst.InitCap(live)
        \tfor i := int(s.m.firstIdx); i < len(s.m.entries); i++ {
        \t\tif s.m.entries[i].HashNext >= 0 {
        \t\t\tdst.AddLast(s.m.entries[i].Key)
        \t\t}
        \t}
        \treturn dst
        }
        func DequePtrFromCustomOrderedSet[T FastCustomKeyComparable](s CustomOrderedSet[T]) *Deque[T] {
        \td := DequeFromCustomOrderedSet(s)
        \treturn &d
        }
        func DequeFromCustomOrderedSetPtr[T FastCustomKeyComparable](s *CustomOrderedSet[T]) Deque[T] {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn DequeFromCustomOrderedSet(*s)
        }
        func DequePtrFromCustomOrderedSetPtr[T FastCustomKeyComparable](s *CustomOrderedSet[T]) *Deque[T] {
        \td := DequeFromCustomOrderedSetPtr(s)
        \treturn &d
        }
        """;
  }

  private static String buildColRtDequeHelpers() {
    return """

        // Deque is a compact generic power-of-2 ring-buffer deque matching Java ArrayDeque
        // with a 4-element inline buffer (0 heap allocations for 0..4 elements) and zero interface boxing.
        type Deque[T any] struct {
        \tbuf    []T
        \thead   int32
        \tcount  int32
        \tinline [4]T
        }

        type DequeCursor[T any] struct {
        \td   *Deque[T]
        \toff int32
        }

        func dequeNullZero[T any]() T {
        \tvar zero T
        \tif _, isStr := any(zero).(string); isStr {
        \t\treturn any(NullString).(T)
        \t}
        \treturn zero
        }

        func (d *Deque[T]) capMask() int32 {
        \tif d.buf == nil {
        \t\treturn 3
        \t}
        \treturn int32(len(d.buf) - 1)
        }

        func (d *Deque[T]) slotGet(phys int32) T {
        \tif d.buf == nil {
        \t\treturn d.inline[phys]
        \t}
        \treturn d.buf[phys]
        }

        func (d *Deque[T]) slotSet(phys int32, v T) {
        \tif d.buf == nil {
        \t\td.inline[phys] = v
        \t} else {
        \t\td.buf[phys] = v
        \t}
        }

        func NewDeque[T any]() Deque[T] {
        \treturn Deque[T]{}
        }

        func NewDequePtr[T any]() *Deque[T] {
        \treturn &Deque[T]{}
        }

        func NewDequeCap[T any](initialCap int32) Deque[T] {
        \tvar d Deque[T]
        \td.InitCap(initialCap)
        \treturn d
        }

        func NewDequePtrCap[T any](initialCap int32) *Deque[T] {
        \td := NewDequeCap[T](initialCap)
        \treturn &d
        }

        func DequeFromDeque[T any](src Deque[T]) Deque[T] {
        \tif src.count == 0 {
        \t\treturn Deque[T]{}
        \t}
        \tif src.count <= 4 {
        \t\tvar dst Deque[T]
        \t\tmask := src.capMask()
        \t\tfor i := int32(0); i < src.count; i++ {
        \t\t\tdst.inline[i] = src.slotGet((src.head + i) & mask)
        \t\t}
        \t\tdst.count = src.count
        \t\treturn dst
        \t}
        \tvar dst Deque[T]
        \tdst.InitCap(src.count)
        \tfirstPart := int32(len(src.buf)) - src.head
        \tif firstPart >= src.count {
        \t\tcopy(dst.buf, src.buf[src.head:src.head+src.count])
        \t} else {
        \t\tcopy(dst.buf, src.buf[src.head:])
        \t\tcopy(dst.buf[firstPart:], src.buf[:src.count-firstPart])
        \t}
        \tdst.count = src.count
        \treturn dst
        }

        func DequePtrFromDeque[T any](src *Deque[T]) *Deque[T] {
        \tif src == nil {
        \t\tThrowNPE()
        \t}
        \tdst := DequeFromDeque(*src)
        \treturn &dst
        }

        func DequeFromSlice[T any](src []T) Deque[T] {
        \tif src == nil {
        \t\tThrowNPE()
        \t}
        \tif len(src) == 0 {
        \t\treturn Deque[T]{}
        \t}
        \tvar dst Deque[T]
        \tdst.InitCap(int32(len(src)))
        \tfor _, v := range src {
        \t\tdst.AddLast(v)
        \t}
        \treturn dst
        }

        func DequePtrFromSlice[T any](src []T) *Deque[T] {
        \td := DequeFromSlice(src)
        \treturn &d
        }

        func DequeFromOrderedSet[T comparable](s OrderedSet[T]) Deque[T] {
        \tlive := s.Size()
        \tif live == 0 {
        \t\treturn Deque[T]{}
        \t}
        \tvar dst Deque[T]
        \tdst.InitCap(live)
        \tif s.holes > 0 {
        \t\tfor i, v := range s.order {
        \t\t\tif s.idx[v] == int32(i+1) {
        \t\t\t\tdst.AddLast(v)
        \t\t\t}
        \t\t}
        \t} else {
        \t\tfor _, v := range s.order {
        \t\t\tdst.AddLast(v)
        \t\t}
        \t}
        \treturn dst
        }

        func DequePtrFromOrderedSet[T comparable](s OrderedSet[T]) *Deque[T] {
        \td := DequeFromOrderedSet(s)
        \treturn &d
        }

        func DequeFromOrderedSetPtr[T comparable](s *OrderedSet[T]) Deque[T] {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn DequeFromOrderedSet(*s)
        }

        func DequePtrFromOrderedSetPtr[T comparable](s *OrderedSet[T]) *Deque[T] {
        \td := DequeFromOrderedSetPtr(s)
        \treturn &d
        }

        func DequeFromLocalOrderedSet[T comparable](s *LocalOrderedSet[T]) Deque[T] {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn DequeFromOrderedSet(s.toOrderedSetView())
        }

        func DequePtrFromLocalOrderedSet[T comparable](s *LocalOrderedSet[T]) *Deque[T] {
        \td := DequeFromLocalOrderedSet(s)
        \treturn &d
        }

        func DequeFromCol[T any](col any) Deque[T] {
        \tCheckNotNull(col)
        \tvar d Deque[T]
        \td.AddAllCol(col)
        \treturn d
        }

        func DequePtrFromCol[T any](col any) *Deque[T] {
        \td := DequeFromCol[T](col)
        \treturn &d
        }

        func (d *Deque[T]) Init() {}

        func (d *Deque[T]) InitCap(initialCap int32) {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif initialCap <= 4 || d.buf != nil {
        \t\treturn
        \t}
        \tc := int32(8)
        \tlimit := initialCap
        \tif limit > (1 << 20) {
        \t\tlimit = 1 << 20
        \t}
        \tfor c < limit {
        \t\tc <<= 1
        \t}
        \tnb := make([]T, c)
        \tfor i := int32(0); i < d.count; i++ {
        \t\tnb[i] = d.inline[(d.head+i)&3]
        \t}
        \td.inline = [4]T{}
        \td.buf = nb
        \td.head = 0
        }

        func (d *Deque[T]) grow() {
        \tif len(d.buf) == 0 {
        \t\tnb := make([]T, 8)
        \t\tfor i := int32(0); i < d.count; i++ {
        \t\t\tnb[i] = d.inline[(d.head+i)&3]
        \t\t}
        \t\td.inline = [4]T{}
        \t\td.buf = nb
        \t\td.head = 0
        \t\treturn
        \t}
        \toldCap := int32(len(d.buf))
        \tif oldCap >= (1 << 30) {
        \t\tpanic(javalang.NewIllegalStateException("Sorry, deque too big"))
        \t}
        \tnewCap := oldCap << 1
        \tnb := make([]T, newCap)
        \tif d.count > 0 {
        \t\tfirstPart := oldCap - d.head
        \t\tif firstPart >= d.count {
        \t\t\tcopy(nb, d.buf[d.head:d.head+d.count])
        \t\t} else {
        \t\t\tcopy(nb, d.buf[d.head:])
        \t\t\tcopy(nb[firstPart:], d.buf[:d.count-firstPart])
        \t\t}
        \t}
        \td.buf = nb
        \td.head = 0
        }

        func (d *Deque[T]) IsNil() bool     { return d == nil }
        func (d *Deque[T]) IsJavaNil_() bool { return d == nil }

        func (d *Deque[T]) Size() int32 {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \treturn d.count
        }

        func (d *Deque[T]) IsEmpty() bool {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \treturn d.count == 0
        }

        func (d *Deque[T]) At(i int32) T {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tCheckElementIndex(i, d.count)
        \treturn d.slotGet((d.head + i) & d.capMask())
        }

        func (d *Deque[T]) AddFirst(v T) {
        \tif d == nil || isJavaNullMapVal(v) {
        \t\tThrowNPE()
        \t}
        \tif d.buf == nil {
        \t\tif d.count < 4 {
        \t\t\td.head = (d.head - 1) & 3
        \t\t\td.inline[d.head] = v
        \t\t\td.count++
        \t\t\treturn
        \t\t}
        \t\td.grow()
        \t} else if int(d.count) == len(d.buf) {
        \t\td.grow()
        \t}
        \td.head = (d.head - 1) & int32(len(d.buf)-1)
        \td.buf[d.head] = v
        \td.count++
        }

        func (d *Deque[T]) Push(v T) {
        \td.AddFirst(v)
        }

        func (d *Deque[T]) OfferFirst(v T) bool {
        \td.AddFirst(v)
        \treturn true
        }

        func (d *Deque[T]) AddLast(v T) {
        \tif d == nil || isJavaNullMapVal(v) {
        \t\tThrowNPE()
        \t}
        \tif d.buf == nil {
        \t\tif d.count < 4 {
        \t\t\td.inline[(d.head+d.count)&3] = v
        \t\t\td.count++
        \t\t\treturn
        \t\t}
        \t\td.grow()
        \t} else if int(d.count) == len(d.buf) {
        \t\td.grow()
        \t}
        \td.buf[(d.head+d.count)&int32(len(d.buf)-1)] = v
        \td.count++
        }

        func (d *Deque[T]) Add(v T) bool {
        \td.AddLast(v)
        \treturn true
        }

        func (d *Deque[T]) Offer(v T) bool {
        \td.AddLast(v)
        \treturn true
        }

        func (d *Deque[T]) OfferLast(v T) bool {
        \td.AddLast(v)
        \treturn true
        }

        func (d *Deque[T]) RemoveFirst() T {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \tvar zero T
        \tif d.buf == nil {
        \t\tv := d.inline[d.head]
        \t\td.inline[d.head] = zero
        \t\td.head = (d.head + 1) & 3
        \t\td.count--
        \t\treturn v
        \t}
        \tv := d.buf[d.head]
        \td.buf[d.head] = zero
        \td.head = (d.head + 1) & int32(len(d.buf)-1)
        \td.count--
        \treturn v
        }

        func (d *Deque[T]) Pop() T {
        \treturn d.RemoveFirst()
        }

        func (d *Deque[T]) RemoveLast() T {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \tvar zero T
        \tif d.buf == nil {
        \t\tidx := (d.head + d.count - 1) & 3
        \t\tv := d.inline[idx]
        \t\td.inline[idx] = zero
        \t\td.count--
        \t\treturn v
        \t}
        \tidx := (d.head + d.count - 1) & int32(len(d.buf)-1)
        \tv := d.buf[idx]
        \td.buf[idx] = zero
        \td.count--
        \treturn v
        }

        func (d *Deque[T]) GetFirst() T {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \treturn d.slotGet(d.head)
        }

        func (d *Deque[T]) Element() T {
        \treturn d.GetFirst()
        }

        func (d *Deque[T]) GetLast() T {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \treturn d.slotGet((d.head + d.count - 1) & d.capMask())
        }

        func (d *Deque[T]) PollFirst() T {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\treturn dequeNullZero[T]()
        \t}
        \tvar zero T
        \tif d.buf == nil {
        \t\tv := d.inline[d.head]
        \t\td.inline[d.head] = zero
        \t\td.head = (d.head + 1) & 3
        \t\td.count--
        \t\treturn v
        \t}
        \tv := d.buf[d.head]
        \td.buf[d.head] = zero
        \td.head = (d.head + 1) & int32(len(d.buf)-1)
        \td.count--
        \treturn v
        }

        func (d *Deque[T]) Poll() T {
        \treturn d.PollFirst()
        }

        func (d *Deque[T]) PollLast() T {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\treturn dequeNullZero[T]()
        \t}
        \tvar zero T
        \tif d.buf == nil {
        \t\tidx := (d.head + d.count - 1) & 3
        \t\tv := d.inline[idx]
        \t\td.inline[idx] = zero
        \t\td.count--
        \t\treturn v
        \t}
        \tidx := (d.head + d.count - 1) & int32(len(d.buf)-1)
        \tv := d.buf[idx]
        \td.buf[idx] = zero
        \td.count--
        \treturn v
        }

        func (d *Deque[T]) PeekFirst() T {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\treturn dequeNullZero[T]()
        \t}
        \treturn d.slotGet(d.head)
        }

        func (d *Deque[T]) Peek() T {
        \treturn d.PeekFirst()
        }

        func (d *Deque[T]) PeekLast() T {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\treturn dequeNullZero[T]()
        \t}
        \treturn d.slotGet((d.head + d.count - 1) & d.capMask())
        }

        func (d *Deque[T]) PollFirstUnbox() T {
        \tif d == nil || d.count == 0 {
        \t\tThrowNPE()
        \t}
        \treturn d.PollFirst()
        }

        func (d *Deque[T]) PollLastUnbox() T {
        \tif d == nil || d.count == 0 {
        \t\tThrowNPE()
        \t}
        \treturn d.PollLast()
        }

        func (d *Deque[T]) PeekFirstUnbox() T {
        \tif d == nil || d.count == 0 {
        \t\tThrowNPE()
        \t}
        \treturn d.dHead()
        }

        func (d *Deque[T]) dHead() T {
        \treturn d.slotGet(d.head)
        }

        func (d *Deque[T]) PeekLastUnbox() T {
        \tif d == nil || d.count == 0 {
        \t\tThrowNPE()
        \t}
        \treturn d.slotGet((d.head + d.count - 1) & d.capMask())
        }

        func (d *Deque[T]) Clear() {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\treturn
        \t}
        \tif d.buf == nil {
        \t\td.inline = [4]T{}
        \t\td.head = 0
        \t\td.count = 0
        \t\treturn
        \t}
        \tvar zero T
        \tmask := int32(len(d.buf) - 1)
        \tfor i := int32(0); i < d.count; i++ {
        \t\td.buf[(d.head+i)&mask] = zero
        \t}
        \td.head = 0
        \td.count = 0
        }

        func (d *Deque[T]) deleteAtOffset(off int32) {
        \tmask := d.capMask()
        \tvar zero T
        \tif off < (d.count >> 1) {
        \t\tfor j := off; j > 0; j-- {
        \t\t\td.slotSet((d.head+j)&mask, d.slotGet((d.head+j-1)&mask))
        \t\t}
        \t\td.slotSet(d.head, zero)
        \t\td.head = (d.head + 1) & mask
        \t\td.count--
        \t\treturn
        \t}
        \tlast := d.count - 1
        \tfor j := off; j < last; j++ {
        \t\td.slotSet((d.head+j)&mask, d.slotGet((d.head+j+1)&mask))
        \t}
        \td.slotSet((d.head+last)&mask, zero)
        \td.count--
        }

        func (d *Deque[T]) Contains(v T) bool {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 || isJavaNullMapVal(v) {
        \t\treturn false
        \t}
        \tmask := d.capMask()
        \tfor i := int32(0); i < d.count; i++ {
        \t\tif javaMapValsEqual(d.slotGet((d.head+i)&mask), v, false) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (d *Deque[T]) ContainsAny(target any) bool {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 || target == nil || isJavaNilAny(target) {
        \t\treturn false
        \t}
        \tt := ToAny(target)
        \tmask := d.capMask()
        \tfor i := int32(0); i < d.count; i++ {
        \t\tif javalang.ObjectsEquals(boxJavaElem(d.slotGet((d.head+i)&mask)), t) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (d *Deque[T]) RemoveFirstOccurrence(v T) bool {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 || isJavaNullMapVal(v) {
        \t\treturn false
        \t}
        \tmask := d.capMask()
        \tfor i := int32(0); i < d.count; i++ {
        \t\tif javaMapValsEqual(d.slotGet((d.head+i)&mask), v, false) {
        \t\t\td.deleteAtOffset(i)
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (d *Deque[T]) RemoveFirstOccurrenceAny(target any) bool {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 || target == nil || isJavaNilAny(target) {
        \t\treturn false
        \t}
        \tt := ToAny(target)
        \tmask := d.capMask()
        \tfor i := int32(0); i < d.count; i++ {
        \t\tif javalang.ObjectsEquals(boxJavaElem(d.slotGet((d.head+i)&mask)), t) {
        \t\t\td.deleteAtOffset(i)
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (d *Deque[T]) RemoveLastOccurrence(v T) bool {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 || isJavaNullMapVal(v) {
        \t\treturn false
        \t}
        \tmask := d.capMask()
        \tfor i := d.count - 1; i >= 0; i-- {
        \t\tif javaMapValsEqual(d.slotGet((d.head+i)&mask), v, false) {
        \t\t\td.deleteAtOffset(i)
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (d *Deque[T]) RemoveLastOccurrenceAny(target any) bool {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 || target == nil || isJavaNilAny(target) {
        \t\treturn false
        \t}
        \tt := ToAny(target)
        \tmask := d.capMask()
        \tfor i := d.count - 1; i >= 0; i-- {
        \t\tif javalang.ObjectsEquals(boxJavaElem(d.slotGet((d.head+i)&mask)), t) {
        \t\t\td.deleteAtOffset(i)
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (d *Deque[T]) AddAllDeque(other *Deque[T]) bool {
        \tif d == nil || other == nil {
        \t\tThrowNPE()
        \t}
        \tn := other.count
        \tif n == 0 {
        \t\treturn false
        \t}
        \tif other == d {
        \t\tsnap := d.IterSlice()
        \t\tfor _, v := range snap {
        \t\t\td.AddLast(v)
        \t\t}
        \t\treturn true
        \t}
        \tmask := other.capMask()
        \tfor i := int32(0); i < n; i++ {
        \t\td.AddLast(other.slotGet((other.head + i) & mask))
        \t}
        \treturn true
        }

        func (d *Deque[T]) AddAllSlice(s []T) bool {
        \tif d == nil || s == nil {
        \t\tThrowNPE()
        \t}
        \tif len(s) == 0 {
        \t\treturn false
        \t}
        \tfor _, v := range s {
        \t\td.AddLast(v)
        \t}
        \treturn true
        }

        func (d *Deque[T]) AddAllCol(col any) bool {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(col)
        \tif tv, ok := col.(*immTypedSliceView[T]); ok {
        \t\treturn d.AddAllSlice(tv.s)
        \t}
        \tadded := false
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\td.AddLast(unboxJavaElem[T](CheckNotNull(cur.Next())))
        \t\tadded = true
        \t}
        \treturn added
        }

        func (d *Deque[T]) Iter() DequeCursor[T] {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \treturn DequeCursor[T]{d: d}
        }

        func (c *DequeCursor[T]) HasNext() bool {
        \treturn c.off < c.d.count
        }

        func (c *DequeCursor[T]) Next() T {
        \ti := c.off
        \tc.off = i + 1
        \treturn c.d.slotGet((c.d.head + i) & c.d.capMask())
        }

        func (d *Deque[T]) IterSlice() []T {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\treturn nil
        \t}
        \tout := make([]T, d.count)
        \tif d.buf == nil {
        \t\tfor i := int32(0); i < d.count; i++ {
        \t\t\tout[i] = d.inline[(d.head+i)&3]
        \t\t}
        \t\treturn out
        \t}
        \tfirstPart := int32(len(d.buf)) - d.head
        \tif firstPart >= d.count {
        \t\tcopy(out, d.buf[d.head:d.head+d.count])
        \t} else {
        \t\tcopy(out, d.buf[d.head:])
        \t\tcopy(out[firstPart:], d.buf[:d.count-firstPart])
        \t}
        \treturn out
        }

        func (d *Deque[T]) String() string {
        \tif d == nil {
        \t\treturn "null"
        \t}
        \tif d.count == 0 {
        \t\treturn "[]"
        \t}
        \tvar sb strings.Builder
        \tsb.WriteByte('[')
        \tmask := d.capMask()
        \tfor i := int32(0); i < d.count; i++ {
        \t\tif i > 0 {
        \t\t\tsb.WriteString(", ")
        \t\t}
        \t\tsb.WriteString(StrVal(boxJavaElem(d.slotGet((d.head + i) & mask))))
        \t}
        \tsb.WriteByte(']')
        \treturn sb.String()
        }

        func DequeJoinStrings(d *Deque[string], delim any) string {
        \tCheckNotNull(delim)
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tsep := StrVal(delim)
        \tif d.count == 0 {
        \t\treturn ""
        \t}
        \tmask := d.capMask()
        \tif d.count == 1 {
        \t\ts := d.slotGet(d.head)
        \t\tif s == NullString {
        \t\t\treturn "null"
        \t\t}
        \t\treturn s
        \t}
        \tvar sb strings.Builder
        \tfor i := int32(0); i < d.count; i++ {
        \t\tif i > 0 {
        \t\t\tsb.WriteString(sep)
        \t\t}
        \t\ts := d.slotGet((d.head + i) & mask)
        \t\tif s == NullString {
        \t\t\tsb.WriteString("null")
        \t\t} else {
        \t\t\tsb.WriteString(s)
        \t\t}
        \t}
        \treturn sb.String()
        }

        type typedDequeIter[T any] struct {
        \td       *Deque[T]
        \toff     int32
        \tlastOff int32
        }

        func (*typedDequeIter[T]) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *typedDequeIter[T]) IsJavaNil_() bool                       { return it == nil }
        func (it *typedDequeIter[T]) HasNext() bool                          { return it.off < it.d.count }
        func (it *typedDequeIter[T]) Next() any {
        \tif it.off >= it.d.count {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tv := it.d.slotGet((it.d.head + it.off) & it.d.capMask())
        \tit.lastOff = it.off
        \tit.off++
        \treturn boxJavaElem(v)
        }
        func (it *typedDequeIter[T]) NextTyped() T {
        \tif it.off >= it.d.count {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tv := it.d.slotGet((it.d.head + it.off) & it.d.capMask())
        \tit.lastOff = it.off
        \tit.off++
        \treturn v
        }
        func (it *typedDequeIter[T]) Remove() {
        \tif it.lastOff < 0 {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \tit.d.deleteAtOffset(it.lastOff)
        \tit.off = it.lastOff
        \tit.lastOff = -1
        }

        func DequeIterator[T any](d *Deque[T]) Ext_java_util_Iterator {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\treturn emptyMapIter_
        \t}
        \treturn &typedDequeIter[T]{d: d, off: 0, lastOff: -1}
        }

        type typedDequeDescIter[T any] struct {
        \td       *Deque[T]
        \toff     int32
        \tlastOff int32
        }

        func (*typedDequeDescIter[T]) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *typedDequeDescIter[T]) IsJavaNil_() bool                       { return it == nil }
        func (it *typedDequeDescIter[T]) HasNext() bool                          { return it.off >= 0 && it.off < it.d.count }
        func (it *typedDequeDescIter[T]) Next() any {
        \tif it.off < 0 || it.off >= it.d.count {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tv := it.d.slotGet((it.d.head + it.off) & it.d.capMask())
        \tit.lastOff = it.off
        \tit.off--
        \treturn boxJavaElem(v)
        }
        func (it *typedDequeDescIter[T]) Remove() {
        \tif it.lastOff < 0 {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \tit.d.deleteAtOffset(it.lastOff)
        \tit.lastOff = -1
        }

        func DequeDescendingIterator[T any](d *Deque[T]) Ext_java_util_Iterator {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\treturn emptyMapIter_
        \t}
        \treturn &typedDequeDescIter[T]{d: d, off: d.count - 1, lastOff: -1}
        }
        """;
  }

  private static String buildColRtNodeDequeHelpers() {
    return """

        // NodeDeque is an unboxed ~uint32 Node specialization of Deque[N] with 1-instruction (v == 0)
        // null checks and zero interface boxing across all ring-buffer operations.
        type NodeDeque[N ~uint32] struct {
        \tbuf    []N
        \thead   int32
        \tcount  int32
        \tinline [4]N
        }

        type NodeDequeCursor[N ~uint32] struct {
        \td   *NodeDeque[N]
        \toff int32
        }

        func (d *NodeDeque[N]) capMask() int32 {
        \tif d.buf == nil {
        \t\treturn 3
        \t}
        \treturn int32(len(d.buf) - 1)
        }

        func (d *NodeDeque[N]) slotGet(phys int32) N {
        \tif d.buf == nil {
        \t\treturn d.inline[phys]
        \t}
        \treturn d.buf[phys]
        }

        func (d *NodeDeque[N]) slotSet(phys int32, v N) {
        \tif d.buf == nil {
        \t\td.inline[phys] = v
        \t} else {
        \t\td.buf[phys] = v
        \t}
        }

        func NewNodeDeque[N ~uint32]() NodeDeque[N] {
        \treturn NodeDeque[N]{}
        }

        func NewNodeDequePtr[N ~uint32]() *NodeDeque[N] {
        \treturn &NodeDeque[N]{}
        }

        func NewNodeDequeCap[N ~uint32](initialCap int32) NodeDeque[N] {
        \tvar d NodeDeque[N]
        \td.InitCap(initialCap)
        \treturn d
        }

        func NewNodeDequePtrCap[N ~uint32](initialCap int32) *NodeDeque[N] {
        \td := NewNodeDequeCap[N](initialCap)
        \treturn &d
        }

        func NodeDequeFromNodeDeque[N ~uint32](src NodeDeque[N]) NodeDeque[N] {
        \tif src.count == 0 {
        \t\treturn NodeDeque[N]{}
        \t}
        \tif src.count <= 4 {
        \t\tvar dst NodeDeque[N]
        \t\tmask := src.capMask()
        \t\tfor i := int32(0); i < src.count; i++ {
        \t\t\tdst.inline[i] = src.slotGet((src.head + i) & mask)
        \t\t}
        \t\tdst.count = src.count
        \t\treturn dst
        \t}
        \tvar dst NodeDeque[N]
        \tdst.InitCap(src.count)
        \tfirstPart := int32(len(src.buf)) - src.head
        \tif firstPart >= src.count {
        \t\tcopy(dst.buf, src.buf[src.head:src.head+src.count])
        \t} else {
        \t\tcopy(dst.buf, src.buf[src.head:])
        \t\tcopy(dst.buf[firstPart:], src.buf[:src.count-firstPart])
        \t}
        \tdst.count = src.count
        \treturn dst
        }

        func NodeDequeFromDeque[N ~uint32](src NodeDeque[N]) NodeDeque[N] {
        \treturn NodeDequeFromNodeDeque(src)
        }

        func NodeDequePtrFromNodeDeque[N ~uint32](src *NodeDeque[N]) *NodeDeque[N] {
        \tif src == nil {
        \t\tThrowNPE()
        \t}
        \tdst := NodeDequeFromNodeDeque(*src)
        \treturn &dst
        }

        func NodeDequePtrFromDeque[N ~uint32](src *NodeDeque[N]) *NodeDeque[N] {
        \treturn NodeDequePtrFromNodeDeque(src)
        }

        func NodeDequeFromSlice[N ~uint32](src []N) NodeDeque[N] {
        \tif src == nil {
        \t\tThrowNPE()
        \t}
        \tif len(src) == 0 {
        \t\treturn NodeDeque[N]{}
        \t}
        \tvar dst NodeDeque[N]
        \tdst.InitCap(int32(len(src)))
        \tfor _, v := range src {
        \t\tdst.AddLast(v)
        \t}
        \treturn dst
        }

        func NodeDequePtrFromSlice[N ~uint32](src []N) *NodeDeque[N] {
        \td := NodeDequeFromSlice(src)
        \treturn &d
        }

        func NodeDequeFromOrderedSet[N ~uint32](s OrderedSet[N]) NodeDeque[N] {
        \tlive := s.Size()
        \tif live == 0 {
        \t\treturn NodeDeque[N]{}
        \t}
        \tvar dst NodeDeque[N]
        \tdst.InitCap(live)
        \tif s.holes > 0 {
        \t\tfor i, v := range s.order {
        \t\t\tif s.idx[v] == int32(i+1) {
        \t\t\t\tdst.AddLast(v)
        \t\t\t}
        \t\t}
        \t} else {
        \t\tfor _, v := range s.order {
        \t\t\tdst.AddLast(v)
        \t\t}
        \t}
        \treturn dst
        }

        func NodeDequePtrFromOrderedSet[N ~uint32](s OrderedSet[N]) *NodeDeque[N] {
        \td := NodeDequeFromOrderedSet(s)
        \treturn &d
        }

        func NodeDequeFromOrderedSetPtr[N ~uint32](s *OrderedSet[N]) NodeDeque[N] {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn NodeDequeFromOrderedSet(*s)
        }

        func NodeDequePtrFromOrderedSetPtr[N ~uint32](s *OrderedSet[N]) *NodeDeque[N] {
        \td := NodeDequeFromOrderedSetPtr(s)
        \treturn &d
        }

        func NodeDequeFromLocalOrderedSet[N ~uint32](s *LocalOrderedSet[N]) NodeDeque[N] {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn NodeDequeFromOrderedSet(s.toOrderedSetView())
        }

        func NodeDequePtrFromLocalOrderedSet[N ~uint32](s *LocalOrderedSet[N]) *NodeDeque[N] {
        \td := NodeDequeFromLocalOrderedSet(s)
        \treturn &d
        }

        func NodeDequeFromIdentityOrderedSet[N ~uint32](s IdentityOrderedSet[N]) NodeDeque[N] {
        \treturn NodeDequeFromOrderedSet(s.raw)
        }

        func NodeDequePtrFromIdentityOrderedSet[N ~uint32](s IdentityOrderedSet[N]) *NodeDeque[N] {
        \treturn NodeDequePtrFromOrderedSet(s.raw)
        }

        func NodeDequeFromIdentityOrderedSetPtr[N ~uint32](s *IdentityOrderedSet[N]) NodeDeque[N] {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn NodeDequeFromOrderedSet(s.raw)
        }

        func NodeDequePtrFromIdentityOrderedSetPtr[N ~uint32](s *IdentityOrderedSet[N]) *NodeDeque[N] {
        \td := NodeDequeFromIdentityOrderedSetPtr(s)
        \treturn &d
        }

        func NodeDequeFromCol[N ~uint32](col any) NodeDeque[N] {
        \tCheckNotNull(col)
        \tvar d NodeDeque[N]
        \td.AddAllCol(col)
        \treturn d
        }

        func NodeDequePtrFromCol[N ~uint32](col any) *NodeDeque[N] {
        \td := NodeDequeFromCol[N](col)
        \treturn &d
        }

        func (d *NodeDeque[N]) Init() {}

        func (d *NodeDeque[N]) InitCap(initialCap int32) {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif initialCap <= 4 || d.buf != nil {
        \t\treturn
        \t}
        \tc := int32(8)
        \tlimit := initialCap
        \tif limit > (1 << 20) {
        \t\tlimit = 1 << 20
        \t}
        \tfor c < limit {
        \t\tc <<= 1
        \t}
        \tnb := make([]N, c)
        \tfor i := int32(0); i < d.count; i++ {
        \t\tnb[i] = d.inline[(d.head+i)&3]
        \t}
        \td.inline = [4]N{}
        \td.buf = nb
        \td.head = 0
        }

        func (d *NodeDeque[N]) grow() {
        \tif len(d.buf) == 0 {
        \t\tnb := make([]N, 8)
        \t\tfor i := int32(0); i < d.count; i++ {
        \t\t\tnb[i] = d.inline[(d.head+i)&3]
        \t\t}
        \t\td.inline = [4]N{}
        \t\td.buf = nb
        \t\td.head = 0
        \t\treturn
        \t}
        \toldCap := int32(len(d.buf))
        \tif oldCap >= (1 << 30) {
        \t\tpanic(javalang.NewIllegalStateException("Sorry, deque too big"))
        \t}
        \tnewCap := oldCap << 1
        \tnb := make([]N, newCap)
        \tif d.count > 0 {
        \t\tfirstPart := oldCap - d.head
        \t\tif firstPart >= d.count {
        \t\t\tcopy(nb, d.buf[d.head:d.head+d.count])
        \t\t} else {
        \t\t\tcopy(nb, d.buf[d.head:])
        \t\t\tcopy(nb[firstPart:], d.buf[:d.count-firstPart])
        \t\t}
        \t}
        \td.buf = nb
        \td.head = 0
        }

        func (d *NodeDeque[N]) IsNil() bool     { return d == nil }
        func (d *NodeDeque[N]) IsJavaNil_() bool { return d == nil }

        func (d *NodeDeque[N]) Size() int32 {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \treturn d.count
        }

        func (d *NodeDeque[N]) IsEmpty() bool {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \treturn d.count == 0
        }

        func (d *NodeDeque[N]) At(i int32) N {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tCheckElementIndex(i, d.count)
        \treturn d.slotGet((d.head + i) & d.capMask())
        }

        func (d *NodeDeque[N]) AddFirst(v N) {
        \tif d == nil || v == 0 {
        \t\tThrowNPE()
        \t}
        \tif d.buf == nil {
        \t\tif d.count < 4 {
        \t\t\td.head = (d.head - 1) & 3
        \t\t\td.inline[d.head] = v
        \t\t\td.count++
        \t\t\treturn
        \t\t}
        \t\td.grow()
        \t} else if int(d.count) == len(d.buf) {
        \t\td.grow()
        \t}
        \td.head = (d.head - 1) & int32(len(d.buf)-1)
        \td.buf[d.head] = v
        \td.count++
        }

        func (d *NodeDeque[N]) Push(v N) {
        \td.AddFirst(v)
        }

        func (d *NodeDeque[N]) OfferFirst(v N) bool {
        \td.AddFirst(v)
        \treturn true
        }

        func (d *NodeDeque[N]) AddLast(v N) {
        \tif d == nil || v == 0 {
        \t\tThrowNPE()
        \t}
        \tif d.buf == nil {
        \t\tif d.count < 4 {
        \t\t\td.inline[(d.head+d.count)&3] = v
        \t\t\td.count++
        \t\t\treturn
        \t\t}
        \t\td.grow()
        \t} else if int(d.count) == len(d.buf) {
        \t\td.grow()
        \t}
        \td.buf[(d.head+d.count)&int32(len(d.buf)-1)] = v
        \td.count++
        }

        func (d *NodeDeque[N]) Add(v N) bool {
        \td.AddLast(v)
        \treturn true
        }

        func (d *NodeDeque[N]) Offer(v N) bool {
        \td.AddLast(v)
        \treturn true
        }

        func (d *NodeDeque[N]) OfferLast(v N) bool {
        \td.AddLast(v)
        \treturn true
        }

        func (d *NodeDeque[N]) RemoveFirst() N {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \tif d.buf == nil {
        \t\tv := d.inline[d.head]
        \t\td.inline[d.head] = 0
        \t\td.head = (d.head + 1) & 3
        \t\td.count--
        \t\treturn v
        \t}
        \tv := d.buf[d.head]
        \td.buf[d.head] = 0
        \td.head = (d.head + 1) & int32(len(d.buf)-1)
        \td.count--
        \treturn v
        }

        func (d *NodeDeque[N]) Pop() N {
        \treturn d.RemoveFirst()
        }

        func (d *NodeDeque[N]) RemoveLast() N {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \tif d.buf == nil {
        \t\tidx := (d.head + d.count - 1) & 3
        \t\tv := d.inline[idx]
        \t\td.inline[idx] = 0
        \t\td.count--
        \t\treturn v
        \t}
        \tidx := (d.head + d.count - 1) & int32(len(d.buf)-1)
        \tv := d.buf[idx]
        \td.buf[idx] = 0
        \td.count--
        \treturn v
        }

        func (d *NodeDeque[N]) GetFirst() N {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \treturn d.slotGet(d.head)
        }

        func (d *NodeDeque[N]) Element() N {
        \treturn d.GetFirst()
        }

        func (d *NodeDeque[N]) GetLast() N {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \treturn d.slotGet((d.head + d.count - 1) & d.capMask())
        }

        func (d *NodeDeque[N]) PollFirst() N {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\treturn 0
        \t}
        \tif d.buf == nil {
        \t\tv := d.inline[d.head]
        \t\td.inline[d.head] = 0
        \t\td.head = (d.head + 1) & 3
        \t\td.count--
        \t\treturn v
        \t}
        \tv := d.buf[d.head]
        \td.buf[d.head] = 0
        \td.head = (d.head + 1) & int32(len(d.buf)-1)
        \td.count--
        \treturn v
        }

        func (d *NodeDeque[N]) Poll() N {
        \treturn d.PollFirst()
        }

        func (d *NodeDeque[N]) PollLast() N {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\treturn 0
        \t}
        \tif d.buf == nil {
        \t\tidx := (d.head + d.count - 1) & 3
        \t\tv := d.inline[idx]
        \t\td.inline[idx] = 0
        \t\td.count--
        \t\treturn v
        \t}
        \tidx := (d.head + d.count - 1) & int32(len(d.buf)-1)
        \tv := d.buf[idx]
        \td.buf[idx] = 0
        \td.count--
        \treturn v
        }

        func (d *NodeDeque[N]) PeekFirst() N {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\treturn 0
        \t}
        \treturn d.slotGet(d.head)
        }

        func (d *NodeDeque[N]) Peek() N {
        \treturn d.PeekFirst()
        }

        func (d *NodeDeque[N]) PeekLast() N {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\treturn 0
        \t}
        \treturn d.slotGet((d.head + d.count - 1) & d.capMask())
        }

        func (d *NodeDeque[N]) PollFirstUnbox() N {
        \tif d == nil || d.count == 0 {
        \t\tThrowNPE()
        \t}
        \treturn d.PollFirst()
        }

        func (d *NodeDeque[N]) PollLastUnbox() N {
        \tif d == nil || d.count == 0 {
        \t\tThrowNPE()
        \t}
        \treturn d.PollLast()
        }

        func (d *NodeDeque[N]) PeekFirstUnbox() N {
        \tif d == nil || d.count == 0 {
        \t\tThrowNPE()
        \t}
        \treturn d.slotGet(d.head)
        }

        func (d *NodeDeque[N]) PeekLastUnbox() N {
        \tif d == nil || d.count == 0 {
        \t\tThrowNPE()
        \t}
        \treturn d.slotGet((d.head + d.count - 1) & d.capMask())
        }

        func (d *NodeDeque[N]) Clear() {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\treturn
        \t}
        \tif d.buf == nil {
        \t\td.inline = [4]N{}
        \t\td.head = 0
        \t\td.count = 0
        \t\treturn
        \t}
        \tmask := int32(len(d.buf) - 1)
        \tfor i := int32(0); i < d.count; i++ {
        \t\td.buf[(d.head+i)&mask] = 0
        \t}
        \td.head = 0
        \td.count = 0
        }

        func (d *NodeDeque[N]) deleteAtOffset(off int32) {
        \tmask := d.capMask()
        \tif off < (d.count >> 1) {
        \t\tfor j := off; j > 0; j-- {
        \t\t\td.slotSet((d.head+j)&mask, d.slotGet((d.head+j-1)&mask))
        \t\t}
        \t\td.slotSet(d.head, 0)
        \t\td.head = (d.head + 1) & mask
        \t\td.count--
        \t\treturn
        \t}
        \tlast := d.count - 1
        \tfor j := off; j < last; j++ {
        \t\td.slotSet((d.head+j)&mask, d.slotGet((d.head+j+1)&mask))
        \t}
        \td.slotSet((d.head+last)&mask, 0)
        \td.count--
        }

        func (d *NodeDeque[N]) Contains(v N) bool {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 || v == 0 {
        \t\treturn false
        \t}
        \tmask := d.capMask()
        \tfor i := int32(0); i < d.count; i++ {
        \t\tif d.slotGet((d.head+i)&mask) == v {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (d *NodeDeque[N]) ContainsAny(target any) bool {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 || target == nil {
        \t\treturn false
        \t}
        \tu, ok := extractDequeNodeId(target)
        \tif !ok || u == 0 {
        \t\treturn false
        \t}
        \t needle := N(u)
        \tmask := d.capMask()
        \tfor i := int32(0); i < d.count; i++ {
        \t\tif d.slotGet((d.head+i)&mask) == needle {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (d *NodeDeque[N]) Remove(v N) bool {
        \treturn d.RemoveFirstOccurrence(v)
        }

        func (d *NodeDeque[N]) RemoveFirstOccurrence(v N) bool {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 || v == 0 {
        \t\treturn false
        \t}
        \tmask := d.capMask()
        \tfor i := int32(0); i < d.count; i++ {
        \t\tif d.slotGet((d.head+i)&mask) == v {
        \t\t\td.deleteAtOffset(i)
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (d *NodeDeque[N]) RemoveFirstOccurrenceAny(target any) bool {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 || target == nil {
        \t\treturn false
        \t}
        \tu, ok := extractDequeNodeId(target)
        \tif !ok || u == 0 {
        \t\treturn false
        \t}
        \tneedle := N(u)
        \tmask := d.capMask()
        \tfor i := int32(0); i < d.count; i++ {
        \t\tif d.slotGet((d.head+i)&mask) == needle {
        \t\t\td.deleteAtOffset(i)
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (d *NodeDeque[N]) RemoveLastOccurrence(v N) bool {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 || v == 0 {
        \t\treturn false
        \t}
        \tmask := d.capMask()
        \tfor i := d.count - 1; i >= 0; i-- {
        \t\tif d.slotGet((d.head+i)&mask) == v {
        \t\t\td.deleteAtOffset(i)
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (d *NodeDeque[N]) RemoveLastOccurrenceAny(target any) bool {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 || target == nil {
        \t\treturn false
        \t}
        \tu, ok := extractDequeNodeId(target)
        \tif !ok || u == 0 {
        \t\treturn false
        \t}
        \tneedle := N(u)
        \tmask := d.capMask()
        \tfor i := d.count - 1; i >= 0; i-- {
        \t\tif d.slotGet((d.head+i)&mask) == needle {
        \t\t\td.deleteAtOffset(i)
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (d *NodeDeque[N]) AddAllNodeDeque(other *NodeDeque[N]) bool {
        \tif d == nil || other == nil {
        \t\tThrowNPE()
        \t}
        \tn := other.count
        \tif n == 0 {
        \t\treturn false
        \t}
        \tif other == d {
        \t\tsnap := d.IterSlice()
        \t\tfor _, v := range snap {
        \t\t\td.AddLast(v)
        \t\t}
        \t\treturn true
        \t}
        \tmask := other.capMask()
        \tfor i := int32(0); i < n; i++ {
        \t\td.AddLast(other.slotGet((other.head + i) & mask))
        \t}
        \treturn true
        }

        func (d *NodeDeque[N]) AddAllDeque(other *NodeDeque[N]) bool {
        \treturn d.AddAllNodeDeque(other)
        }

        func (d *NodeDeque[N]) AddAllSlice(s []N) bool {
        \tif d == nil || s == nil {
        \t\tThrowNPE()
        \t}
        \tif len(s) == 0 {
        \t\treturn false
        \t}
        \tfor _, v := range s {
        \t\td.AddLast(v)
        \t}
        \treturn true
        }

        func (d *NodeDeque[N]) AddAllCol(col any) bool {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(col)
        \tif tv, ok := col.(*immTypedSliceView[N]); ok {
        \t\treturn d.AddAllSlice(tv.s)
        \t}
        \tadded := false
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\td.AddLast(N(cur.NextNode()))
        \t\tadded = true
        \t}
        \treturn added
        }

        func (d *NodeDeque[N]) Iter() NodeDequeCursor[N] {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \treturn NodeDequeCursor[N]{d: d}
        }

        func (c *NodeDequeCursor[N]) HasNext() bool {
        \treturn c.off < c.d.count
        }

        func (c *NodeDequeCursor[N]) Next() N {
        \ti := c.off
        \tc.off = i + 1
        \treturn c.d.slotGet((c.d.head + i) & c.d.capMask())
        }

        func (d *NodeDeque[N]) IterSlice() []N {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\treturn nil
        \t}
        \tout := make([]N, d.count)
        \tif d.buf == nil {
        \t\tfor i := int32(0); i < d.count; i++ {
        \t\t\tout[i] = d.inline[(d.head+i)&3]
        \t\t}
        \t\treturn out
        \t}
        \tfirstPart := int32(len(d.buf)) - d.head
        \tif firstPart >= d.count {
        \t\tcopy(out, d.buf[d.head:d.head+d.count])
        \t} else {
        \t\tcopy(out, d.buf[d.head:])
        \t\tcopy(out[firstPart:], d.buf[:d.count-firstPart])
        \t}
        \treturn out
        }

        func (d *NodeDeque[N]) ToArray_0() []any {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\treturn []any{}
        \t}
        \tout := make([]any, d.count)
        \tmask := d.capMask()
        \tfor i := int32(0); i < d.count; i++ {
        \t\tout[i] = boxNodeById(uint32(d.slotGet((d.head + i) & mask)))
        \t}
        \treturn out
        }

        func (d *NodeDeque[N]) Iterator() Ext_java_util_Iterator {
        \treturn NodeDequeIterator(d)
        }

        func (d *NodeDeque[N]) DescendingIterator() Ext_java_util_Iterator {
        \treturn NodeDequeDescendingIterator(d)
        }

        func (d *NodeDeque[N]) String() string {
        \tif d == nil {
        \t\treturn "null"
        \t}
        \tif d.count == 0 {
        \t\treturn "[]"
        \t}
        \tvar sb strings.Builder
        \tsb.WriteByte('[')
        \tmask := d.capMask()
        \tfor i := int32(0); i < d.count; i++ {
        \t\tif i > 0 {
        \t\t\tsb.WriteString(", ")
        \t\t}
        \t\tsb.WriteString(StrVal(boxNodeById(uint32(d.slotGet((d.head + i) & mask)))))
        \t}
        \tsb.WriteByte(']')
        \treturn sb.String()
        }

        type typedNodeDequeIter[N ~uint32] struct {
        \td       *NodeDeque[N]
        \toff     int32
        \tlastOff int32
        }

        func (*typedNodeDequeIter[N]) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *typedNodeDequeIter[N]) IsJavaNil_() bool                       { return it == nil }
        func (it *typedNodeDequeIter[N]) HasNext() bool                          { return it.off < it.d.count }
        func (it *typedNodeDequeIter[N]) Next() any {
        \tif it.off >= it.d.count {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tv := it.d.slotGet((it.d.head + it.off) & it.d.capMask())
        \tit.lastOff = it.off
        \tit.off++
        \treturn boxNodeById(uint32(v))
        }
        func (it *typedNodeDequeIter[N]) NextTyped() N {
        \tif it.off >= it.d.count {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tv := it.d.slotGet((it.d.head + it.off) & it.d.capMask())
        \tit.lastOff = it.off
        \tit.off++
        \treturn v
        }
        func (it *typedNodeDequeIter[N]) Remove() {
        \tif it.lastOff < 0 {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \tit.d.deleteAtOffset(it.lastOff)
        \tit.off = it.lastOff
        \tit.lastOff = -1
        }

        func NodeDequeIterator[N ~uint32](d *NodeDeque[N]) Ext_java_util_Iterator {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\treturn emptyMapIter_
        \t}
        \treturn &typedNodeDequeIter[N]{d: d, off: 0, lastOff: -1}
        }

        type typedNodeDequeDescIter[N ~uint32] struct {
        \td       *NodeDeque[N]
        \toff     int32
        \tlastOff int32
        }

        func (*typedNodeDequeDescIter[N]) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *typedNodeDequeDescIter[N]) IsJavaNil_() bool                       { return it == nil }
        func (it *typedNodeDequeDescIter[N]) HasNext() bool                          { return it.off >= 0 && it.off < it.d.count }
        func (it *typedNodeDequeDescIter[N]) Next() any {
        \tif it.off < 0 || it.off >= it.d.count {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tv := it.d.slotGet((it.d.head + it.off) & it.d.capMask())
        \tit.lastOff = it.off
        \tit.off--
        \treturn boxNodeById(uint32(v))
        }
        func (it *typedNodeDequeDescIter[N]) Remove() {
        \tif it.lastOff < 0 {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \tit.d.deleteAtOffset(it.lastOff)
        \tit.lastOff = -1
        }

        func NodeDequeDescendingIterator[N ~uint32](d *NodeDeque[N]) Ext_java_util_Iterator {
        \tif d == nil {
        \t\tThrowNPE()
        \t}
        \tif d.count == 0 {
        \t\treturn emptyMapIter_
        \t}
        \treturn &typedNodeDequeDescIter[N]{d: d, off: d.count - 1, lastOff: -1}
        }
        """;
  }

  private static String buildColRtSlabAndTypedListHelpers() {
    return """

        const (
        \tslabSmallChunkSize = 16
        \tslabChunkSize      = 256
        \tslabPoolMaxBytes   = 384 << 20
        \tslabSmallPoolMaxB  = 32 << 20
        \tslabGlobalMaxBytes = 1024 << 20
        \tslabTrimMinBytes   = 96 << 20
        \tslabMinEvictBytes  = 256 << 10
        )

        var (
        \tslabGlobalPooledBytes atomic.Int64
        \tslabTrimLastClock     atomic.Uint64
        )

        // Slab is a compact 32-byte pure bump-pointer arena allocator for package-owned AST/graph leaf structs
        // whose lifetime matches their owner and are not pooled across runs.
        type Slab[T any] struct {
        \tbuf     []T
        \tnextCap int32
        }

        func (s *Slab[T]) Alloc() *T {
        \tif i := len(s.buf); i < cap(s.buf) {
        \t\tb := s.buf[:i+1]
        \t\ts.buf = b
        \t\treturn &b[i]
        \t}
        \treturn s.allocSlow()
        }

        func (s *Slab[T]) allocSlow() *T {
        \tn := s.nextCap
        \tif n == 0 {
        \t\tn = 16
        \t\ts.nextCap = 32
        \t} else if n < 256 {
        \t\ts.nextCap = n * 2
        \t} else {
        \t\ts.nextCap = 256
        \t}
        \tc := make([]T, 1, n)
        \ts.buf = c
        \treturn &c[0]
        }

        func (s *Slab[T]) Release() {
        \tif s == nil {
        \t\treturn
        \t}
        \tclear(s.buf)
        \ts.buf = nil
        \ts.nextCap = 0
        }

        // PooledSlab is a 48-byte two-tier pooled chunk arena allocator using a pooled 16-item
        // slabSmallChunk[T] for small per-script owners before promoting to pooled 128-item slabChunk[T]
        // chunks for large whole-AST owners, with zero unpooled bootstrap slices.
        type PooledSlab[T any] struct {
        \tbuf   []T
        \tsmall *slabSmallChunk[T]
        \thead  *slabChunk[T]
        \tpool  *slabChunkPool[T]
        }

        func (s *PooledSlab[T]) Alloc() *T {
        \tif i := len(s.buf); i < cap(s.buf) {
        \t\tb := s.buf[:i+1]
        \t\ts.buf = b
        \t\treturn &b[i]
        \t}
        \treturn s.allocSlow()
        }

        func (s *PooledSlab[T]) allocSlow() *T {
        \tp := s.pool
        \tif p == nil {
        \t\tp = getSlabChunkPool[T]()
        \t\ts.pool = p
        \t}
        \tif s.small == nil && s.head == nil {
        \t\tsc := p.allocSmall()
        \t\ts.small = sc
        \t\ts.buf = sc.items[:1:slabSmallChunkSize]
        \t\treturn &sc.items[0]
        \t}
        \tc := p.alloc()
        \tc.next = s.head
        \ts.head = c
        \ts.buf = c.items[:1:slabChunkSize]
        \treturn &c.items[0]
        }

        func (s *PooledSlab[T]) Reset() {
        \tif s == nil {
        \t\treturn
        \t}
        \tsc := s.small
        \thead := s.head
        \tif head == nil {
        \t\tif sc == nil {
        \t\t\ts.buf = nil
        \t\t\treturn
        \t\t}
        \t\tif n := len(s.buf); n > 0 {
        \t\t\tclear(sc.items[:n])
        \t\t}
        \t\ts.buf = sc.items[:0:slabSmallChunkSize]
        \t\treturn
        \t}
        \tn := len(s.buf)
        \ts.head = nil
        \tif sc != nil {
        \t\tclear(sc.items[:])
        \t\ts.buf = sc.items[:0:slabSmallChunkSize]
        \t} else {
        \t\ts.buf = nil
        \t}
        \tif n > 0 {
        \t\tclear(head.items[:n])
        \t}
        \tfor c := head.next; c != nil; c = c.next {
        \t\tclear(c.items[:])
        \t}
        \tp := s.pool
        \tif p == nil {
        \t\tp = getSlabChunkPool[T]()
        \t\ts.pool = p
        \t}
        \tp.releaseList(head)
        }

        func (s *PooledSlab[T]) Release() {
        \tif s == nil {
        \t\treturn
        \t}
        \tsc := s.small
        \thead := s.head
        \tif sc == nil && head == nil {
        \t\ts.buf = nil
        \t\treturn
        \t}
        \tn := len(s.buf)
        \ts.buf = nil
        \ts.small = nil
        \ts.head = nil
        \tp := s.pool
        \tif p == nil {
        \t\tp = getSlabChunkPool[T]()
        \t\ts.pool = p
        \t}
        \tif head == nil {
        \t\tif n > 0 {
        \t\t\tclear(sc.items[:n])
        \t\t}
        \t\tp.releaseSmall(sc)
        \t\treturn
        \t}
        \tif sc != nil {
        \t\tclear(sc.items[:])
        \t\tp.releaseSmall(sc)
        \t}
        \tif n > 0 {
        \t\tclear(head.items[:n])
        \t}
        \tfor c := head.next; c != nil; c = c.next {
        \t\tclear(c.items[:])
        \t}
        \tp.releaseList(head)
        }

        func (s *PooledSlab[T]) ReleaseWith(cleanup func(*T)) {
        \tif s == nil {
        \t\treturn
        \t}
        \tif cleanup == nil {
        \t\ts.Release()
        \t\treturn
        \t}
        \tsc := s.small
        \thead := s.head
        \tif sc == nil && head == nil {
        \t\ts.buf = nil
        \t\treturn
        \t}
        \tn := len(s.buf)
        \ts.buf = nil
        \ts.small = nil
        \ts.head = nil
        \tp := s.pool
        \tif p == nil {
        \t\tp = getSlabChunkPool[T]()
        \t\ts.pool = p
        \t}
        \tif head == nil {
        \t\tfor i := 0; i < n; i++ {
        \t\t\tcleanup(&sc.items[i])
        \t\t}
        \t\tif n > 0 {
        \t\t\tclear(sc.items[:n])
        \t\t}
        \t\tp.releaseSmall(sc)
        \t\treturn
        \t}
        \tif sc != nil {
        \t\tfor i := 0; i < slabSmallChunkSize; i++ {
        \t\t\tcleanup(&sc.items[i])
        \t\t}
        \t\tclear(sc.items[:])
        \t\tp.releaseSmall(sc)
        \t}
        \tfor i := 0; i < n; i++ {
        \t\tcleanup(&head.items[i])
        \t}
        \tif n > 0 {
        \t\tclear(head.items[:n])
        \t}
        \tfor c := head.next; c != nil; c = c.next {
        \t\tfor i := 0; i < slabChunkSize; i++ {
        \t\t\tcleanup(&c.items[i])
        \t\t}
        \t\tclear(c.items[:])
        \t}
        \tp.releaseList(head)
        }

        const (
        \tpooledSliceMinShift        = 3
        \tpooledSliceMinCap          = 1 << pooledSliceMinShift // 8
        \tpooledSliceMaxShift        = 22
        \tpooledSliceMaxCap          = 1 << pooledSliceMaxShift               // 4194304
        \tpooledSliceTiers           = pooledSliceMaxShift - pooledSliceMinShift + 1 // 20
        \tpooledSliceTierBudgetBytes = 4 << 20                                // 4 MB max per tier per type
        \tpooledSliceTierMaxCount    = 65536
        )

        type pooledSliceTier[T any] struct {
        \tmu      sync.Mutex
        \tfree    [][]T
        \tmaxKeep int
        }

        type pooledSliceTypePool[T any] struct {
        \ttiers [pooledSliceTiers]pooledSliceTier[T]
        }

        var pooledSlicePools sync.Map

        func newPooledSliceTypePool[T any]() *pooledSliceTypePool[T] {
        \tp := &pooledSliceTypePool[T]{}
        \telemBytes := max(1, int(reflect.TypeFor[T]().Size()))
        \tfor i := 0; i < pooledSliceTiers; i++ {
        \t\ttierCap := pooledSliceMinCap << i
        \t\tsliceBytes := tierCap * elemBytes
        \t\tvar maxKeep int
        \t\tif sliceBytes <= 262144 {
        \t\t\tmaxKeep = min(pooledSliceTierMaxCount, max(2, pooledSliceTierBudgetBytes/sliceBytes))
        \t\t} else {
        \t\t\tmaxKeep = 2
        \t\t}
        \t\tp.tiers[i].maxKeep = maxKeep
        \t}
        \treturn p
        }

        func (p *pooledSliceTypePool[T]) purge() {
        \tfor i := 0; i < pooledSliceTiers; i++ {
        \t\tt := &p.tiers[i]
        \t\tt.mu.Lock()
        \t\tt.free = nil
        \t\tt.mu.Unlock()
        \t}
        }

        func getPooledSlicePool[T any]() *pooledSliceTypePool[T] {
        \trt := reflect.TypeFor[T]()
        \tif v, ok := pooledSlicePools.Load(rt); ok {
        \t\treturn v.(*pooledSliceTypePool[T])
        \t}
        \tp := newPooledSliceTypePool[T]()
        \tactual, loaded := pooledSlicePools.LoadOrStore(rt, p)
        \tif !loaded {
        \t\tpurge := func() {
        \t\t\tp.purge()
        \t\t}
        \t\tregisterCompPoolReset(purge)
        \t\tregisterRuntimePoolPurger(purge)
        \t}
        \treturn actual.(*pooledSliceTypePool[T])
        }

        func AcquirePooledSlice[T any](capHint int) []T {
        \tif capHint <= 0 {
        \t\treturn nil
        \t}
        \tif capHint > pooledSliceMaxCap {
        \t\treturn make([]T, 0, capHint)
        \t}
        \ttierIdx := 0
        \tallocCap := pooledSliceMinCap
        \tif capHint > pooledSliceMinCap {
        \t\tshift := bits.Len(uint(capHint - 1))
        \t\ttierIdx = shift - pooledSliceMinShift
        \t\tallocCap = 1 << shift
        \t}
        \tt := &getPooledSlicePool[T]().tiers[tierIdx]
        \tt.mu.Lock()
        \tif n := len(t.free); n > 0 {
        \t\ts := t.free[n-1]
        \t\tt.free[n-1] = nil
        \t\tt.free = t.free[:n-1]
        \t\tt.mu.Unlock()
        \t\treturn s[:0]
        \t}
        \tt.mu.Unlock()
        \treturn make([]T, 0, allocCap)
        }

        func AppendPooledSlice[T any](s *[]T, v T) {
        \tcur := *s
        \tif len(cur) < cap(cur) {
        \t\t*s = append(cur, v)
        \t\treturn
        \t}
        \tappendPooledSliceGrow(s, v)
        }

        func appendPooledSliceGrow[T any](s *[]T, v T) {
        \told := *s
        \toldLen := len(old)
        \tnewCap := cap(old) * 2
        \tif newCap < pooledSliceMinCap {
        \t\tnewCap = pooledSliceMinCap
        \t}
        \tif newCap <= oldLen {
        \t\tnewCap = oldLen + 1
        \t}
        \tnext := AcquirePooledSlice[T](newCap)
        \tnext = next[:oldLen+1]
        \tcopy(next, old)
        \tnext[oldLen] = v
        \tif cap(old) >= pooledSliceMinCap {
        \t\tReleasePooledSlice(old)
        \t}
        \t*s = next
        }

        func ReleasePooledSlice[T any](s []T) {
        \tc := cap(s)
        \tif c < pooledSliceMinCap || c > pooledSliceMaxCap {
        \t\treturn
        \t}
        \tclear(s[:c])
        \ttierIdx := bits.Len(uint(c)) - 1 - pooledSliceMinShift
        \tt := &getPooledSlicePool[T]().tiers[tierIdx]
        \tt.mu.Lock()
        \tif len(t.free) < t.maxKeep {
        \t\tt.free = append(t.free, s[:0])
        \t}
        \tt.mu.Unlock()
        }

        type jTypedViewCursorSource interface {
        \ttypedViewLen_() (int, bool)
        \ttypedViewAt_(i int) any
        \ttypedViewNodeAt_(i int) (uint32, bool)
        }

        func typedElemRawNode[T any](v T) (uint32, bool) {
        \tvar zero T
        \tswitch any(zero).(type) {
        \tcase uint32:
        \t\treturn any(v).(uint32), true
        \tcase nodestore.Node:
        \t\treturn uint32(any(v).(nodestore.Node)), true
        \t}
        \tif rn, ok := any(v).(jRawNodeProvider); ok {
        \t\treturn rn.RawNodeId_(), true
        \t}
        \treturn 0, false
        }

        // jTypedArrayList wraps a native []T slice as a full mutable Ext_java_util_ArrayList without
        // boxing elements into []any unless untyped mutation/subList fallback is invoked.
        type jTypedArrayList[T any] struct {
        \ts  []T
        \tlc *JListCore
        }

        var _ Ext_java_util_ArrayList = (*jTypedArrayList[int32])(nil)

        func WrapTypedArrayList[T any](s []T) Ext_java_util_ArrayList {
        \tif s == nil {
        \t\treturn nil
        \t}
        \treturn &jTypedArrayList[T]{s: s}
        }

        func (l *jTypedArrayList[T]) materialize() *JListCore {
        \tif l.lc == nil {
        \t\titems := make([]any, len(l.s))
        \t\tfor i, v := range l.s {
        \t\t\titems[i] = boxJavaElem(v)
        \t\t}
        \t\tl.s = nil
        \t\tl.lc = &JListCore{Items: items}
        \t}
        \treturn l.lc
        }

        func tryUnboxTypedArrayListElem[T any](raw any) (T, bool) {
        \tvar zero T
        \tif raw == nil || isJavaNilAny(raw) {
        \t\tswitch any(zero).(type) {
        \t\tcase string:
        \t\t\treturn any(NullString).(T), true
        \t\tcase int32, int64, bool, float64, float32, uint16, int8, int16:
        \t\t\treturn zero, false
        \t\tdefault:
        \t\t\treturn zero, true
        \t\t}
        \t}
        \tswitch any(zero).(type) {
        \tcase string:
        \t\tif s, ok := raw.(string); ok {
        \t\t\treturn any(s).(T), true
        \t\t}
        \t\treturn zero, false
        \tcase int32:
        \t\tif b, ok := raw.(*BoxInt); ok && b != nil {
        \t\t\treturn any(b.Val).(T), true
        \t\t}
        \t\tif v, ok := raw.(int32); ok {
        \t\t\treturn any(v).(T), true
        \t\t}
        \t\treturn zero, false
        \tcase int64:
        \t\tif b, ok := raw.(*BoxLong); ok && b != nil {
        \t\t\treturn any(b.Val).(T), true
        \t\t}
        \t\tif v, ok := raw.(int64); ok {
        \t\t\treturn any(v).(T), true
        \t\t}
        \t\treturn zero, false
        \tcase bool:
        \t\tif b, ok := raw.(*BoxBoolean); ok && b != nil {
        \t\t\treturn any(b.Val).(T), true
        \t\t}
        \t\tif v, ok := raw.(bool); ok {
        \t\t\treturn any(v).(T), true
        \t\t}
        \t\treturn zero, false
        \tcase float64:
        \t\tif b, ok := raw.(*BoxDouble); ok && b != nil {
        \t\t\treturn any(b.Val).(T), true
        \t\t}
        \t\tif v, ok := raw.(float64); ok {
        \t\t\treturn any(v).(T), true
        \t\t}
        \t\treturn zero, false
        \tcase float32:
        \t\tif b, ok := raw.(*BoxFloat); ok && b != nil {
        \t\t\treturn any(b.Val).(T), true
        \t\t}
        \t\tif v, ok := raw.(float32); ok {
        \t\t\treturn any(v).(T), true
        \t\t}
        \t\treturn zero, false
        \tcase uint16:
        \t\tif b, ok := raw.(*BoxChar); ok && b != nil {
        \t\t\treturn any(b.Val).(T), true
        \t\t}
        \t\tif v, ok := raw.(uint16); ok {
        \t\t\treturn any(v).(T), true
        \t\t}
        \t\treturn zero, false
        \tcase int8:
        \t\tif b, ok := raw.(*BoxByte); ok && b != nil {
        \t\t\treturn any(b.Val).(T), true
        \t\t}
        \t\tif v, ok := raw.(int8); ok {
        \t\t\treturn any(v).(T), true
        \t\t}
        \t\treturn zero, false
        \tcase int16:
        \t\tif b, ok := raw.(*BoxShort); ok && b != nil {
        \t\t\treturn any(b.Val).(T), true
        \t\t}
        \t\tif v, ok := raw.(int16); ok {
        \t\t\treturn any(v).(T), true
        \t\t}
        \t\treturn zero, false
        \tcase uint32:
        \t\tif rn, ok := raw.(jRawNodeProvider); ok {
        \t\t\treturn any(rn.RawNodeId_()).(T), true
        \t\t}
        \t\tif nd, ok := raw.(nodestore.Node); ok {
        \t\t\treturn any(uint32(nd)).(T), true
        \t\t}
        \t\tif u, ok := raw.(uint32); ok {
        \t\t\treturn any(u).(T), true
        \t\t}
        \t\treturn zero, false
        \tcase nodestore.Node:
        \t\tif rn, ok := raw.(jRawNodeProvider); ok {
        \t\t\treturn any(nodestore.Node(rn.RawNodeId_())).(T), true
        \t\t}
        \t\tif nd, ok := raw.(nodestore.Node); ok {
        \t\t\treturn any(nd).(T), true
        \t\t}
        \t\tif u, ok := raw.(uint32); ok {
        \t\t\treturn any(nodestore.Node(u)).(T), true
        \t\t}
        \t\treturn zero, false
        \tcase jRawNodeProvider:
        \t\tif v, ok := raw.(T); ok {
        \t\t\treturn v, true
        \t\t}
        \t\tvar u uint32
        \t\tvar haveU bool
        \t\tif rn, ok := raw.(jRawNodeProvider); ok {
        \t\t\tu = rn.RawNodeId_()
        \t\t\thaveU = true
        \t\t} else if nd, ok := raw.(nodestore.Node); ok {
        \t\t\tu = uint32(nd)
        \t\t\thaveU = true
        \t\t} else if rawU, ok := raw.(uint32); ok {
        \t\t\tu = rawU
        \t\t\thaveU = true
        \t\t}
        \t\tif !haveU {
        \t\t\treturn zero, false
        \t\t}
        \t\tif u == 0 {
        \t\t\treturn zero, true
        \t\t}
        \t\tif globalNodeBoxer != nil {
        \t\t\tif b, ok := globalNodeBoxer(u).(T); ok {
        \t\t\t\treturn b, true
        \t\t\t}
        \t\t}
        \t\trt := reflect.TypeFor[T]()
        \t\tif rt.Kind() == reflect.Uint32 {
        \t\t\treturn reflect.ValueOf(u).Convert(rt).Interface().(T), true
        \t\t}
        \t\treturn zero, false
        \tdefault:
        \t\tif v, ok := raw.(T); ok {
        \t\t\treturn v, true
        \t\t}
        \t\treturn zero, false
        \t}
        }

        func (l *jTypedArrayList[T]) IsJavaNil_() bool { return l == nil }
        func (*jTypedArrayList[T]) JavaClassInfo_() (string, string) {
        \treturn "java.util.ArrayList", "ArrayList"
        }
        func (*jTypedArrayList[T]) Implements_j2gort_Ext_java_util_ArrayList()          {}
        func (*jTypedArrayList[T]) Implements_j2gort_Ext_java_util_AbstractList()       {}
        func (*jTypedArrayList[T]) Implements_j2gort_Ext_java_util_List()               {}
        func (*jTypedArrayList[T]) Implements_j2gort_Ext_java_util_RandomAccess()       {}
        func (*jTypedArrayList[T]) Implements_j2gort_Ext_java_util_AbstractCollection() {}
        func (*jTypedArrayList[T]) Implements_j2gort_Ext_java_util_SequencedCollection() {}
        func (*jTypedArrayList[T]) Implements_j2gort_Ext_java_util_Collection()         {}
        func (*jTypedArrayList[T]) Implements_j2gort_Ext_java_lang_Iterable()           {}
        func (*jTypedArrayList[T]) Implements_j2gort_Ext_java_io_Serializable()         {}
        func (*jTypedArrayList[T]) GetExt_java_util_ArrayListBase() *Ext_java_util_ArrayListBase {
        \treturn nil
        }
        func (*jTypedArrayList[T]) GetExt_java_util_AbstractListBase() *Ext_java_util_AbstractListBase {
        \treturn nil
        }
        func (*jTypedArrayList[T]) GetExt_java_util_AbstractCollectionBase() *Ext_java_util_AbstractCollectionBase {
        \treturn nil
        }
        func (l *jTypedArrayList[T]) typedViewLen_() (int, bool) {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc != nil {
        \t\treturn 0, false
        \t}
        \treturn len(l.s), true
        }
        func (l *jTypedArrayList[T]) typedViewAt_(i int) any {
        \treturn boxJavaElem(l.s[i])
        }
        func (l *jTypedArrayList[T]) typedViewNodeAt_(i int) (uint32, bool) {
        \treturn typedElemRawNode(l.s[i])
        }
        func (l *jTypedArrayList[T]) IterSlice_() []any {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc != nil {
        \t\treturn l.lc.Items
        \t}
        \tif len(l.s) == 0 {
        \t\treturn emptyAnySlice_
        \t}
        \tout := make([]any, len(l.s))
        \tfor i, v := range l.s {
        \t\tout[i] = boxJavaElem(v)
        \t}
        \treturn out
        }
        func (l *jTypedArrayList[T]) ListGetStr_(i int32) string {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc != nil {
        \t\treturn CastString(l.lc.Get(i))
        \t}
        \tCheckElementIndex(i, int32(len(l.s)))
        \treturn CastString(boxJavaElem(l.s[i]))
        }
        func (l *jTypedArrayList[T]) Size() int32 {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc != nil {
        \t\treturn l.lc.Size()
        \t}
        \treturn int32(len(l.s))
        }
        func (l *jTypedArrayList[T]) IsEmpty() bool {
        \treturn l.Size() == 0
        }
        func (l *jTypedArrayList[T]) Get(i int32) any {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc != nil {
        \t\treturn l.lc.Get(i)
        \t}
        \tCheckElementIndex(i, int32(len(l.s)))
        \treturn boxJavaElem(l.s[i])
        }
        func (l *jTypedArrayList[T]) GetFirst() any {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc != nil {
        \t\tif l.lc.IsEmpty() {
        \t\t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t\t}
        \t\treturn l.lc.Get(0)
        \t}
        \tif len(l.s) == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \treturn boxJavaElem(l.s[0])
        }
        func (l *jTypedArrayList[T]) GetLast() any {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc != nil {
        \t\tif l.lc.IsEmpty() {
        \t\t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t\t}
        \t\treturn l.lc.Get(l.lc.Size() - 1)
        \t}
        \tif len(l.s) == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \treturn boxJavaElem(l.s[len(l.s)-1])
        }
        func (l *jTypedArrayList[T]) Set(i int32, x any) any {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc == nil {
        \t\tif v, ok := tryUnboxTypedArrayListElem[T](x); ok {
        \t\t\tCheckElementIndex(i, int32(len(l.s)))
        \t\t\told := boxJavaElem(l.s[i])
        \t\t\tl.s[i] = v
        \t\t\treturn old
        \t\t}
        \t}
        \treturn l.materialize().Set(i, x)
        }
        func (l *jTypedArrayList[T]) Add_1(x any) bool {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc == nil {
        \t\tif v, ok := tryUnboxTypedArrayListElem[T](x); ok {
        \t\t\tl.s = append(l.s, v)
        \t\t\treturn true
        \t\t}
        \t}
        \treturn l.materialize().Add(x)
        }
        func (l *jTypedArrayList[T]) Add_2(i int32, x any) {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc == nil {
        \t\tif v, ok := tryUnboxTypedArrayListElem[T](x); ok {
        \t\t\tl.s = SliceInsertAt(l.s, i, v)
        \t\t\treturn
        \t\t}
        \t}
        \tl.materialize().AddAt(i, x)
        }
        func (l *jTypedArrayList[T]) AddFirst(x any) {
        \tl.Add_2(0, x)
        }
        func (l *jTypedArrayList[T]) AddLast(x any) {
        \tl.Add_1(x)
        }
        func (l *jTypedArrayList[T]) Remove_1_int32(i int32) any {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc != nil {
        \t\treturn l.lc.RemoveAt(i)
        \t}
        \treturn boxJavaElem(SliceRemoveAtElem(&l.s, i))
        }
        func (l *jTypedArrayList[T]) RemoveFirst() any {
        \tif l.Size() == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \treturn l.Remove_1_int32(0)
        }
        func (l *jTypedArrayList[T]) RemoveLast() any {
        \tn := l.Size()
        \tif n == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \treturn l.Remove_1_int32(n - 1)
        }
        func (l *jTypedArrayList[T]) Remove_1_any(x any) bool {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc != nil {
        \t\treturn l.lc.RemoveVal(x)
        \t}
        \tt := ToAny(x)
        \tfor i, v := range l.s {
        \t\tif javalang.ObjectsEquals(boxJavaElem(v), t) {
        \t\t\tSliceRemoveAtElem(&l.s, int32(i))
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }
        func (l *jTypedArrayList[T]) Clear() {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc != nil {
        \t\tl.lc.Clear()
        \t\treturn
        \t}
        \tclear(l.s)
        \tl.s = l.s[:0]
        }
        func (l *jTypedArrayList[T]) Contains(x any) bool {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc != nil {
        \t\treturn l.lc.Contains(x)
        \t}
        \tt := ToAny(x)
        \tfor _, v := range l.s {
        \t\tif javalang.ObjectsEquals(boxJavaElem(v), t) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }
        func (l *jTypedArrayList[T]) ContainsAll(col Ext_java_util_Collection) bool {
        \treturn colContainsAll(l.Contains, col)
        }
        func (l *jTypedArrayList[T]) IndexOf_1(x any) int32 {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc != nil {
        \t\treturn l.lc.IndexOf(x)
        \t}
        \tt := ToAny(x)
        \tfor i, v := range l.s {
        \t\tif javalang.ObjectsEquals(boxJavaElem(v), t) {
        \t\t\treturn int32(i)
        \t\t}
        \t}
        \treturn -1
        }
        func (l *jTypedArrayList[T]) LastIndexOf_1(x any) int32 {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc != nil {
        \t\treturn l.lc.LastIndexOf(x)
        \t}
        \tt := ToAny(x)
        \tfor i := len(l.s) - 1; i >= 0; i-- {
        \t\tif javalang.ObjectsEquals(boxJavaElem(l.s[i]), t) {
        \t\t\treturn int32(i)
        \t\t}
        \t}
        \treturn -1
        }
        func (l *jTypedArrayList[T]) AddAll_1(col Ext_java_util_Collection) bool {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(col)
        \tif l.lc != nil {
        \t\treturn listCoreAddAllIterable(l.lc, col)
        \t}
        \tadded := false
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\tl.Add_1(cur.Next())
        \t\tadded = true
        \t}
        \treturn added
        }
        func (l *jTypedArrayList[T]) RemoveAll(col Ext_java_util_Collection) bool {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \treturn listCoreRemoveAll(l.materialize(), col)
        }
        func (l *jTypedArrayList[T]) Sort(cmp Ext_java_util_Comparator) {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tlistCoreSort(l.materialize(), cmp)
        }
        func (l *jTypedArrayList[T]) SubList(fromIdx, toIdx int32) any {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \treturn newSubListView(l.materialize(), fromIdx, toIdx)
        }
        func (l *jTypedArrayList[T]) ListIterator_1(idx int32) any {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \treturn newFullListIter(l.materialize(), idx)
        }
        func (l *jTypedArrayList[T]) Iterator() any {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc != nil {
        \t\treturn newListIter(l.lc)
        \t}
        \treturn &jTypedArrayListIter[T]{l: l, lastIdx: -1}
        }
        func (l *jTypedArrayList[T]) ForEach(action Ext_java_util_function_Consumer) {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(action)
        \tif l.lc != nil {
        \t\tfor _, x := range l.lc.Items {
        \t\t\taction.Accept(x)
        \t\t}
        \t\treturn
        \t}
        \tfor _, v := range l.s {
        \t\taction.Accept(boxJavaElem(v))
        \t}
        }
        func (l *jTypedArrayList[T]) ToArray_1_slice_any(dst []any) []any {
        \treturn sliceToArray1(l.IterSlice_(), dst)
        }
        func (l *jTypedArrayList[T]) Stream() Ext_java_util_stream_Stream {
        \treturn streamFromSlice(l.IterSlice_())
        }
        func (l *jTypedArrayList[T]) Equals(other any) bool {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc != nil {
        \t\treturn listEqualsSlice(l.lc.Items, other)
        \t}
        \tif any(l) == other {
        \t\treturn true
        \t}
        \treturn ImmSliceEquals(l.s, other)
        }
        func (l *jTypedArrayList[T]) HashCode() int32 {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif l.lc != nil {
        \t\treturn listHashSlice(l.lc.Items)
        \t}
        \treturn ImmSliceHashCode(l.s)
        }
        func (l *jTypedArrayList[T]) String() string {
        \tif l == nil {
        \t\treturn "null"
        \t}
        \tif l.lc != nil {
        \t\treturn sliceJavaString(l.lc.Items)
        \t}
        \tif len(l.s) == 0 {
        \t\treturn "[]"
        \t}
        \tvar sb strings.Builder
        \tsb.WriteByte('[')
        \tfor i, v := range l.s {
        \t\tif i > 0 {
        \t\t\tsb.WriteString(", ")
        \t\t}
        \t\tsb.WriteString(StrVal(boxJavaElem(v)))
        \t}
        \tsb.WriteByte(']')
        \treturn sb.String()
        }

        type jTypedArrayListIter[T any] struct {
        \tl       *jTypedArrayList[T]
        \tpos     int
        \tlastIdx int
        }

        func (*jTypedArrayListIter[T]) Implements_j2gort_Ext_java_util_Iterator() {}
        func (it *jTypedArrayListIter[T]) IsJavaNil_() bool                       { return it == nil }
        func (it *jTypedArrayListIter[T]) HasNext() bool                          { return int32(it.pos) < it.l.Size() }
        func (it *jTypedArrayListIter[T]) Next() any {
        \tif int32(it.pos) >= it.l.Size() {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tv := it.l.Get(int32(it.pos))
        \tit.lastIdx = it.pos
        \tit.pos++
        \treturn v
        }
        func (it *jTypedArrayListIter[T]) Remove() {
        \tif it.lastIdx < 0 {
        \t\tpanic(javalang.NewIllegalStateException("remove before next"))
        \t}
        \tit.l.Remove_1_int32(int32(it.lastIdx))
        \tit.pos = it.lastIdx
        \tit.lastIdx = -1
        }
        """;
  }

  private static String buildColRtStringBuilderHelpers() {
    return """

        // LocalStringBuilder is a stack-allocated append-only StringBuilder backed by strings.Builder
        // with exact Java UTF-16 surrogate-pair stitching and WTF-8 lone-surrogate preservation.
        type LocalStringBuilder struct {
        \tbuf         strings.Builder
        \tpendingHigh uint16
        }

        func (b *LocalStringBuilder) InitCap(cap int32) *LocalStringBuilder {
        \tif cap < 0 {
        \t\tpanic(javalang.NewRuntimeExceptionMsg("NegativeArraySizeException"))
        \t}
        \tif cap > 0 && cap <= (1<<24) {
        \t\tb.buf.Grow(int(cap))
        \t}
        \treturn b
        }

        func (b *LocalStringBuilder) InitStr(s string) *LocalStringBuilder {
        \tCheckNotNullString(s)
        \treturn b.AppendStr(s)
        }

        func (b *LocalStringBuilder) flushPendingHigh() {
        \tif b.pendingHigh != 0 {
        \t\tch := b.pendingHigh
        \t\tb.pendingHigh = 0
        \t\tb.buf.WriteByte(byte(0xE0 | (ch >> 12)))
        \t\tb.buf.WriteByte(byte(0x80 | ((ch >> 6) & 0x3F)))
        \t\tb.buf.WriteByte(byte(0x80 | (ch & 0x3F)))
        \t}
        }

        func (b *LocalStringBuilder) AppendStr(raw string) *LocalStringBuilder {
        \ts := raw
        \tif s == NullString {
        \t\ts = "null"
        \t}
        \tif len(s) == 0 {
        \t\treturn b
        \t}
        \tif b.pendingHigh == 0 && s[len(s)-1] < 0x80 {
        \t\tb.buf.WriteString(s)
        \t\treturn b
        \t}
        \tif b.pendingHigh != 0 {
        \t\tif len(s) >= 3 && s[0] == 0xED && s[1] >= 0xB0 && s[1] <= 0xBF && (s[2]&0xC0) == 0x80 {
        \t\t\tlow := (uint16(s[0]&0x0F) << 12) | (uint16(s[1]&0x3F) << 6) | uint16(s[2]&0x3F)
        \t\t\tr := utf16.DecodeRune(rune(b.pendingHigh), rune(low))
        \t\t\tb.pendingHigh = 0
        \t\t\tb.buf.WriteRune(r)
        \t\t\ts = s[3:]
        \t\t\tif len(s) == 0 {
        \t\t\t\treturn b
        \t\t\t}
        \t\t} else {
        \t\t\tb.flushPendingHigh()
        \t\t}
        \t}
        \tn := len(s)
        \tif n >= 3 && s[n-3] == 0xED && s[n-2] >= 0xA0 && s[n-2] <= 0xAF && (s[n-1]&0xC0) == 0x80 {
        \t\thigh := (uint16(s[n-3]&0x0F) << 12) | (uint16(s[n-2]&0x3F) << 6) | uint16(s[n-1]&0x3F)
        \t\tb.buf.WriteString(s[:n-3])
        \t\tb.pendingHigh = high
        \t\treturn b
        \t}
        \tb.buf.WriteString(s)
        \treturn b
        }

        func (b *LocalStringBuilder) AppendChar(ch uint16) *LocalStringBuilder {
        \tif b.pendingHigh == 0 && ch < 0x80 {
        \t\tb.buf.WriteByte(byte(ch))
        \t\treturn b
        \t}
        \tif b.pendingHigh != 0 {
        \t\tif ch >= 0xDC00 && ch <= 0xDFFF {
        \t\t\tr := utf16.DecodeRune(rune(b.pendingHigh), rune(ch))
        \t\t\tb.pendingHigh = 0
        \t\t\tb.buf.WriteRune(r)
        \t\t\treturn b
        \t\t}
        \t\tb.flushPendingHigh()
        \t}
        \tif ch >= 0xD800 && ch <= 0xDBFF {
        \t\tb.pendingHigh = ch
        \t\treturn b
        \t}
        \tif ch < 0x80 {
        \t\tb.buf.WriteByte(byte(ch))
        \t} else if ch < 0x800 {
        \t\tb.buf.WriteByte(byte(0xC0 | (ch >> 6)))
        \t\tb.buf.WriteByte(byte(0x80 | (ch & 0x3F)))
        \t} else {
        \t\tb.buf.WriteByte(byte(0xE0 | (ch >> 12)))
        \t\tb.buf.WriteByte(byte(0x80 | ((ch >> 6) & 0x3F)))
        \t\tb.buf.WriteByte(byte(0x80 | (ch & 0x3F)))
        \t}
        \treturn b
        }

        func (b *LocalStringBuilder) AppendChars(chars []uint16) *LocalStringBuilder {
        \tCheckNotNullSlice(chars)
        \tfor _, ch := range chars {
        \t\tb.AppendChar(ch)
        \t}
        \treturn b
        }

        func (b *LocalStringBuilder) AppendSub(csq any, start, end int32) *LocalStringBuilder {
        \treturn b.AppendStr(JStrSub(StrVal(csq), start, end))
        }

        func (b *LocalStringBuilder) AppendCodePoint(cp int32) *LocalStringBuilder {
        \tif cp >= 0 && cp <= 0xFFFF {
        \t\treturn b.AppendChar(uint16(cp))
        \t}
        \tif cp > 0xFFFF && cp <= 0x10FFFF {
        \t\tb.flushPendingHigh()
        \t\tb.buf.WriteRune(rune(cp))
        \t\treturn b
        \t}
        \tpanic(javalang.NewIllegalArgumentException("Not a valid Unicode code point"))
        }

        func (b *LocalStringBuilder) AppendInt(v int32) *LocalStringBuilder {
        \treturn b.AppendLong(int64(v))
        }

        func (b *LocalStringBuilder) AppendLong(v int64) *LocalStringBuilder {
        \tb.flushPendingHigh()
        \tvar tmp [20]byte
        \tb.buf.Write(strconv.AppendInt(tmp[:0], v, 10))
        \treturn b
        }

        func (b *LocalStringBuilder) AppendBool(v bool) *LocalStringBuilder {
        \tb.flushPendingHigh()
        \tif v {
        \t\tb.buf.WriteString("true")
        \t} else {
        \t\tb.buf.WriteString("false")
        \t}
        \treturn b
        }

        func (b *LocalStringBuilder) AppendDouble(v float64) *LocalStringBuilder {
        \tb.flushPendingHigh()
        \tb.buf.WriteString(javalang.DoubleToString(v))
        \treturn b
        }

        func (b *LocalStringBuilder) AppendFloat(v float32) *LocalStringBuilder {
        \tb.flushPendingHigh()
        \tb.buf.WriteString(javalang.FloatToString(v))
        \treturn b
        }

        func (b *LocalStringBuilder) AppendAny(v any) *LocalStringBuilder {
        \treturn b.AppendStr(StrVal(v))
        }

        func (b *LocalStringBuilder) IsEmpty() bool {
        \treturn b.buf.Len() == 0 && b.pendingHigh == 0
        }

        func (b *LocalStringBuilder) Reset() {
        \tb.buf.Reset()
        \tb.pendingHigh = 0
        }

        func (b *LocalStringBuilder) String() string {
        \tif b.pendingHigh == 0 {
        \t\treturn b.buf.String()
        \t}
        \ts := b.buf.String()
        \tch := b.pendingHigh
        \tout := make([]byte, len(s)+3)
        \tcopy(out, s)
        \tout[len(s)] = byte(0xE0 | (ch >> 12))
        \tout[len(s)+1] = byte(0x80 | ((ch >> 6) & 0x3F))
        \tout[len(s)+2] = byte(0x80 | (ch & 0x3F))
        \treturn string(out)
        }

        type slabSmallChunk[T any] struct {
        \tnext  *slabSmallChunk[T]
        \titems [slabSmallChunkSize]T
        }

        type slabChunk[T any] struct {
        \tnext  *slabChunk[T]
        \titems [slabChunkSize]T
        }

        type slabEvictablePool interface {
        \tlastTouchVal() uint64
        \tpooledBytesApprox() int64
        \tevictBytes(maxBytes int64) int64
        }

        type slabChunkPool[T any] struct {
        \tmu           sync.Mutex
        \tfree         []*slabChunk[T]
        \tsmallFree    []*slabSmallChunk[T]
        \tpooledBytes  atomic.Int64
        \tlastTouch    atomic.Uint64
        \tfreeLen      atomic.Int32
        \tsmallLen     atomic.Int32
        \tmaxFree      int
        \tmaxSmallFree int
        \tchunkBytes   int64
        \tsmallBytes   int64
        }

        var (
        \tslabClock      atomic.Uint64
        \tslabChunkPools sync.Map
        )

        func (p *slabChunkPool[T]) lastTouchVal() uint64 {
        \treturn p.lastTouch.Load()
        }

        func (p *slabChunkPool[T]) pooledBytesApprox() int64 {
        \treturn p.pooledBytes.Load()
        }

        func (p *slabChunkPool[T]) evictBytes(maxBytes int64) int64 {
        \tif maxBytes <= 0 || p.pooledBytes.Load() <= 0 {
        \t\treturn 0
        \t}
        \tp.mu.Lock()
        \tvar freed int64
        \tif nFree := len(p.free); nFree > 0 {
        \t\trem := maxBytes
        \t\tvar drop int
        \t\tif rem >= int64(nFree)*p.chunkBytes {
        \t\t\tdrop = nFree
        \t\t} else {
        \t\t\tdrop = int((rem + p.chunkBytes - 1) / p.chunkBytes)
        \t\t}
        \t\tif drop > 0 {
        \t\t\tnewLen := nFree - drop
        \t\t\tclear(p.free[newLen:])
        \t\t\tif newLen == 0 {
        \t\t\t\tp.free = nil
        \t\t\t} else {
        \t\t\t\tp.free = p.free[:newLen]
        \t\t\t}
        \t\t\tp.freeLen.Store(int32(newLen))
        \t\t\tfreed += int64(drop) * p.chunkBytes
        \t\t}
        \t}
        \tif freed < maxBytes {
        \t\tif nSmall := len(p.smallFree); nSmall > 0 {
        \t\t\trem := maxBytes - freed
        \t\t\tvar drop int
        \t\t\tif rem >= int64(nSmall)*p.smallBytes {
        \t\t\t\tdrop = nSmall
        \t\t\t} else {
        \t\t\t\tdrop = int((rem + p.smallBytes - 1) / p.smallBytes)
        \t\t\t}
        \t\t\tif drop > 0 {
        \t\t\t\tnewLen := nSmall - drop
        \t\t\t\tclear(p.smallFree[newLen:])
        \t\t\t\tif newLen == 0 {
        \t\t\t\t\tp.smallFree = nil
        \t\t\t\t} else {
        \t\t\t\t\tp.smallFree = p.smallFree[:newLen]
        \t\t\t\t}
        \t\t\t\tp.smallLen.Store(int32(newLen))
        \t\t\t\tfreed += int64(drop) * p.smallBytes
        \t\t\t}
        \t\t}
        \t}
        \tif freed > 0 {
        \t\tp.pooledBytes.Add(-freed)
        \t\tslabGlobalPooledBytes.Add(-freed)
        \t}
        \tp.mu.Unlock()
        \treturn freed
        }

        type slabEvictCand struct {
        \tpool  slabEvictablePool
        \ttouch uint64
        \tbytes int64
        }

        func slabEvictColdestPools(exclude slabEvictablePool, needBytes int64) int64 {
        \tif needBytes <= 0 {
        \t\treturn 0
        \t}
        \tvar cands []slabEvictCand
        \tslabChunkPools.Range(func(_, v any) bool {
        \t\tep, ok := v.(slabEvictablePool)
        \t\tif !ok || ep == exclude {
        \t\t\treturn true
        \t\t}
        \t\tif b := ep.pooledBytesApprox(); b > 0 {
        \t\t\tcands = append(cands, slabEvictCand{
        \t\t\t\tpool:  ep,
        \t\t\t\ttouch: ep.lastTouchVal(),
        \t\t\t\tbytes: b,
        \t\t\t})
        \t\t}
        \t\treturn true
        \t})
        \tif len(cands) == 0 {
        \t\treturn 0
        \t}
        \tsort.Slice(cands, func(i, j int) bool {
        \t\treturn cands[i].touch < cands[j].touch
        \t})
        \tvar freed int64
        \tfor i := range cands {
        \t\tstillNeeded := needBytes - freed
        \t\tif stillNeeded <= 0 {
        \t\t\tbreak
        \t\t}
        \t\tfreed += cands[i].pool.evictBytes(stillNeeded)
        \t}
        \treturn freed
        }

        func TrimStaleSlabPools(keepRecent int) {
        \tconst maxInt64 = int64(^uint64(0) >> 1)
        \tif keepRecent <= 0 {
        \t\tif slabGlobalPooledBytes.Load() <= 0 {
        \t\t\treturn
        \t\t}
        \t\tslabChunkPools.Range(func(_, v any) bool {
        \t\t\tif ep, ok := v.(slabEvictablePool); ok && ep.pooledBytesApprox() > 0 {
        \t\t\t\tep.evictBytes(maxInt64)
        \t\t\t}
        \t\t\treturn true
        \t\t})
        \t\tslabTrimLastClock.Store(slabClock.Load())
        \t\treturn
        \t}
        \tcurTotal := slabGlobalPooledBytes.Load()
        \tif curTotal < slabTrimMinBytes {
        \t\treturn
        \t}
        \tcurClock := slabClock.Load()
        \tif curClock == slabTrimLastClock.Load() && curTotal <= slabGlobalMaxBytes {
        \t\treturn
        \t}
        \tslabTrimLastClock.Store(curClock)
        \tvar cands []slabEvictCand
        \tslabChunkPools.Range(func(_, v any) bool {
        \t\tif ep, ok := v.(slabEvictablePool); ok {
        \t\t\tif b := ep.pooledBytesApprox(); b > 0 {
        \t\t\t\tcands = append(cands, slabEvictCand{
        \t\t\t\t\tpool:  ep,
        \t\t\t\t\ttouch: ep.lastTouchVal(),
        \t\t\t\t\tbytes: b,
        \t\t\t\t})
        \t\t\t}
        \t\t}
        \t\treturn true
        \t})
        \tif len(cands) == 0 {
        \t\treturn
        \t}
        \tsort.Slice(cands, func(i, j int) bool {
        \t\treturn cands[i].touch > cands[j].touch
        \t})
        \tvar keptBytes int64
        \twarmLimit := min(keepRecent, len(cands))
        \tfor i := 0; i < warmLimit; i++ {
        \t\tremBudget := int64(slabGlobalMaxBytes) - keptBytes
        \t\tif remBudget <= 0 {
        \t\t\tcands[i].pool.evictBytes(maxInt64)
        \t\t\tcontinue
        \t\t}
        \t\tif cands[i].bytes > remBudget {
        \t\t\tcands[i].pool.evictBytes(cands[i].bytes - remBudget)
        \t\t\tkeptBytes = slabGlobalMaxBytes
        \t\t} else {
        \t\t\tkeptBytes += cands[i].bytes
        \t\t}
        \t}
        \tfor i := warmLimit; i < len(cands); i++ {
        \t\tcands[i].pool.evictBytes(maxInt64)
        \t}
        }

        func getSlabChunkPool[T any]() *slabChunkPool[T] {
        \trt := reflect.TypeFor[T]()
        \tif v, ok := slabChunkPools.Load(rt); ok {
        \t\treturn v.(*slabChunkPool[T])
        \t}
        \tchunkBytes := int64(max(1, int(reflect.TypeFor[slabChunk[T]]().Size())))
        \tsmallBytes := int64(max(1, int(reflect.TypeFor[slabSmallChunk[T]]().Size())))
        \tmaxFree := max(16, int(int64(slabPoolMaxBytes)/chunkBytes))
        \tmaxSmallFree := max(64, int(int64(slabSmallPoolMaxB)/smallBytes))
        \tp := &slabChunkPool[T]{
        \t\tmaxFree:      maxFree,
        \t\tmaxSmallFree: maxSmallFree,
        \t\tchunkBytes:   chunkBytes,
        \t\tsmallBytes:   smallBytes,
        \t}
        \tp.lastTouch.Store(slabClock.Add(1))
        \tactual, loaded := slabChunkPools.LoadOrStore(rt, p)
        \tif !loaded {
        \t\tpurge := func() {
        \t\t\tp.evictBytes(int64(^uint64(0) >> 1))
        \t\t\tp.mu.Lock()
        \t\t\tif len(p.free) == 0 {
        \t\t\t\tp.free = nil
        \t\t\t}
        \t\t\tif len(p.smallFree) == 0 {
        \t\t\t\tp.smallFree = nil
        \t\t\t}
        \t\t\tp.mu.Unlock()
        \t\t}
        \t\tregisterCompPoolReset(purge)
        \t\tregisterRuntimePoolPurger(purge)
        \t}
        \treturn actual.(*slabChunkPool[T])
        }

        func (p *slabChunkPool[T]) alloc() *slabChunk[T] {
        \tp.lastTouch.Store(slabClock.Add(1))
        \tp.mu.Lock()
        \tif n := len(p.free); n > 0 {
        \t\tc := p.free[n-1]
        \t\tp.free[n-1] = nil
        \t\tp.free = p.free[:n-1]
        \t\tp.freeLen.Store(int32(n - 1))
        \t\tp.pooledBytes.Add(-p.chunkBytes)
        \t\tslabGlobalPooledBytes.Add(-p.chunkBytes)
        \t\tp.mu.Unlock()
        \t\treturn c
        \t}
        \tp.mu.Unlock()
        \treturn new(slabChunk[T])
        }

        func (p *slabChunkPool[T]) allocSmall() *slabSmallChunk[T] {
        \tp.lastTouch.Store(slabClock.Add(1))
        \tp.mu.Lock()
        \tif n := len(p.smallFree); n > 0 {
        \t\tsc := p.smallFree[n-1]
        \t\tp.smallFree[n-1] = nil
        \t\tp.smallFree = p.smallFree[:n-1]
        \t\tp.smallLen.Store(int32(n - 1))
        \t\tp.pooledBytes.Add(-p.smallBytes)
        \t\tslabGlobalPooledBytes.Add(-p.smallBytes)
        \t\tp.mu.Unlock()
        \t\treturn sc
        \t}
        \tp.mu.Unlock()
        \treturn new(slabSmallChunk[T])
        }

        func (p *slabChunkPool[T]) releaseSmall(sc *slabSmallChunk[T]) {
        \tif sc == nil {
        \t\treturn
        \t}
        \tsc.next = nil
        \tp.lastTouch.Store(slabClock.Add(1))
        \tsb := p.smallBytes
        \tif int(p.smallLen.Load()) < p.maxSmallFree {
        \t\tif curGlobal := slabGlobalPooledBytes.Load(); curGlobal+sb > slabGlobalMaxBytes {
        \t\t\tneed := (curGlobal + sb) - slabGlobalMaxBytes
        \t\t\tslabEvictColdestPools(p, max(need, int64(slabMinEvictBytes)))
        \t\t}
        \t}
        \tp.mu.Lock()
        \tif len(p.smallFree) < p.maxSmallFree && slabGlobalPooledBytes.Load()+sb <= slabGlobalMaxBytes {
        \t\tif tot := slabGlobalPooledBytes.Add(sb); tot <= slabGlobalMaxBytes {
        \t\t\tp.smallFree = append(p.smallFree, sc)
        \t\t\tp.smallLen.Store(int32(len(p.smallFree)))
        \t\t\tp.pooledBytes.Add(sb)
        \t\t\tp.mu.Unlock()
        \t\t\treturn
        \t\t} else {
        \t\t\tslabGlobalPooledBytes.Add(-sb)
        \t\t}
        \t}
        \tp.mu.Unlock()
        }

        func (p *slabChunkPool[T]) releaseList(first *slabChunk[T]) {
        \tif first == nil {
        \t\treturn
        \t}
        \tp.lastTouch.Store(slabClock.Add(1))
        \tcb := p.chunkBytes
        \tif maxRoom := p.maxFree - int(p.freeLen.Load()); maxRoom > 0 {
        \t\twantCnt := 0
        \t\tfor c := first; c != nil && wantCnt < maxRoom; c = c.next {
        \t\t\twantCnt++
        \t\t}
        \t\twantBytes := int64(wantCnt) * cb
        \t\tif curGlobal := slabGlobalPooledBytes.Load(); curGlobal+wantBytes > slabGlobalMaxBytes {
        \t\t\tneed := (curGlobal + wantBytes) - slabGlobalMaxBytes
        \t\t\tslabEvictColdestPools(p, max(need, int64(slabMinEvictBytes)))
        \t\t}
        \t}
        \tp.mu.Lock()
        \tfree := p.free
        \troom := p.maxFree - len(free)
        \tif room > 0 {
        \t\trem := int64(slabGlobalMaxBytes) - slabGlobalPooledBytes.Load()
        \t\tif rem <= 0 {
        \t\t\troom = 0
        \t\t} else if gr := int(rem / cb); gr < room {
        \t\t\troom = gr
        \t\t}
        \t}
        \tvar added int64
        \tc := first
        \tfor c != nil && room > 0 {
        \t\tcn := c.next
        \t\tc.next = nil
        \t\tfree = append(free, c)
        \t\tc = cn
        \t\troom--
        \t\tadded++
        \t}
        \tif added > 0 {
        \t\taddBytes := added * cb
        \t\tif tot := slabGlobalPooledBytes.Add(addBytes); tot > slabGlobalMaxBytes {
        \t\t\tdrop := int((tot - slabGlobalMaxBytes + cb - 1) / cb)
        \t\t\tif drop > int(added) {
        \t\t\t\tdrop = int(added)
        \t\t\t}
        \t\t\tif drop > 0 {
        \t\t\t\tnewLen := len(free) - drop
        \t\t\t\tclear(free[newLen:])
        \t\t\t\tfree = free[:newLen]
        \t\t\t\tdropBytes := int64(drop) * cb
        \t\t\t\tslabGlobalPooledBytes.Add(-dropBytes)
        \t\t\t\taddBytes -= dropBytes
        \t\t\t}
        \t\t}
        \t\tif addBytes > 0 {
        \t\t\tp.pooledBytes.Add(addBytes)
        \t\t}
        \t}
        \tp.free = free
        \tp.freeLen.Store(int32(len(free)))
        \tp.mu.Unlock()
        \tfor c != nil {
        \t\tcn := c.next
        \t\tc.next = nil
        \t\tc = cn
        \t}
        }

        func allocSlabChunk[T any]() *slabChunk[T] {
        \treturn getSlabChunkPool[T]().alloc()
        }
        """;
  }

  private static String buildColRtPlainMapSetHelpers() {
    return """

        // PlainMap is an un-ordered value struct backed solely by a native Go map[K]V (zero entries slice)
        // for non-iterated HashMap lowering.
        type PlainMap[K comparable, V any] struct {
        \tm map[K]V
        }

        func NewPlainMap[K comparable, V any](capHint int) *PlainMap[K, V] {
        \tif capHint < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", capHint)))
        \t}
        \tif capHint > 0 {
        \t\treturn &PlainMap[K, V]{m: make(map[K]V, capHint)}
        \t}
        \treturn &PlainMap[K, V]{}
        }

        func NewPlainMapVal[K comparable, V any]() PlainMap[K, V] {
        \treturn PlainMap[K, V]{}
        }

        func NewPlainMapCap[K comparable, V any](capHint int32) PlainMap[K, V] {
        \tif capHint < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", capHint)))
        \t}
        \tif capHint > 0 {
        \t\treturn PlainMap[K, V]{m: make(map[K]V, capHint)}
        \t}
        \treturn PlainMap[K, V]{}
        }

        func NewPlainMapPtr[K comparable, V any]() *PlainMap[K, V] {
        \treturn &PlainMap[K, V]{}
        }

        func NewPlainMapPtrCap[K comparable, V any](capHint int32) *PlainMap[K, V] {
        \treturn NewPlainMap[K, V](int(capHint))
        }

        func (m *PlainMap[K, V]) Init() {}

        func (m *PlainMap[K, V]) InitCap(capHint int) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif capHint < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", capHint)))
        \t}
        \tif capHint > 0 && m.m == nil {
        \t\tm.m = make(map[K]V, capHint)
        \t}
        }

        func (m *PlainMap[K, V]) IsNil() bool     { return m == nil }
        func (m *PlainMap[K, V]) IsJavaNil_() bool { return m == nil }

        func (m *PlainMap[K, V]) Size() int32 {
        \tif m == nil {
        \t\treturn 0
        \t}
        \treturn int32(len(m.m))
        }

        func (m *PlainMap[K, V]) IsEmpty() bool {
        \treturn m == nil || len(m.m) == 0
        }

        func (m *PlainMap[K, V]) ContainsKey(k K) bool {
        \tif m == nil || len(m.m) == 0 {
        \t\treturn false
        \t}
        \t_, ok := m.m[k]
        \treturn ok
        }

        func (m *PlainMap[K, V]) Get(k K) V {
        \tif m == nil || len(m.m) == 0 {
        \t\tvar zero V
        \t\treturn zero
        \t}
        \treturn m.m[k]
        }

        func (m *PlainMap[K, V]) GetOk(k K) (V, bool) {
        \tif m == nil || len(m.m) == 0 {
        \t\tvar zero V
        \t\treturn zero, false
        \t}
        \tv, ok := m.m[k]
        \treturn v, ok
        }

        func (m *PlainMap[K, V]) GetOrDefault(k K, def V) V {
        \tif m == nil || len(m.m) == 0 {
        \t\treturn def
        \t}
        \tif v, ok := m.m[k]; ok {
        \t\treturn v
        \t}
        \treturn def
        }

        func (m *PlainMap[K, V]) Put(k K, v V) V {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif m.m == nil {
        \t\tm.m = make(map[K]V)
        \t\tm.m[k] = v
        \t\tvar zero V
        \t\treturn zero
        \t}
        \told := m.m[k]
        \tm.m[k] = v
        \treturn old
        }

        func (m *PlainMap[K, V]) PutOk(k K, v V) (V, bool) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif m.m == nil {
        \t\tm.m = make(map[K]V)
        \t\tm.m[k] = v
        \t\tvar zero V
        \t\treturn zero, false
        \t}
        \told, existed := m.m[k]
        \tm.m[k] = v
        \treturn old, existed
        }

        func (m *PlainMap[K, V]) PutVoid(k K, v V) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif m.m == nil {
        \t\tm.m = make(map[K]V)
        \t}
        \tm.m[k] = v
        }

        func (m *PlainMap[K, V]) PutIfAbsentOk(k K, v V) (V, bool) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif m.m == nil {
        \t\tm.m = make(map[K]V)
        \t\tm.m[k] = v
        \t\tvar zero V
        \t\treturn zero, false
        \t}
        \tcur, ok := m.m[k]
        \tif !ok || isJavaNullMapVal(cur) {
        \t\tm.m[k] = v
        \t\treturn cur, false
        \t}
        \treturn cur, true
        }

        func (m *PlainMap[K, V]) PutIfAbsent(k K, v V) V {
        \tcur, _ := m.PutIfAbsentOk(k, v)
        \treturn cur
        }

        func (m *PlainMap[K, V]) RemoveOk(k K) (V, bool) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tif len(m.m) == 0 {
        \t\tvar zero V
        \t\treturn zero, false
        \t}
        \told, existed := m.m[k]
        \tif existed {
        \t\tdelete(m.m, k)
        \t}
        \treturn old, existed
        }

        func (m *PlainMap[K, V]) Remove(k K) V {
        \told, _ := m.RemoveOk(k)
        \treturn old
        }

        func (m *PlainMap[K, V]) Clear() {
        \tif m != nil && len(m.m) > 0 {
        \t\tclear(m.m)
        \t}
        }

        func PlainMapComputeIfAbsentFn[K comparable, V any](m *PlainMap[K, V], k K, fn func(K) V) V {
        \tif m == nil || fn == nil {
        \t\tThrowNPE()
        \t}
        \tif cur, ok := m.GetOk(k); ok && !isJavaNullMapVal(cur) {
        \t\treturn cur
        \t}
        \tnv := fn(k)
        \tif !isJavaNullMapVal(nv) {
        \t\tm.PutVoid(k, nv)
        \t}
        \treturn nv
        }

        func PlainMapComputeIfAbsent[K comparable, V comparable](m *PlainMap[K, V], k K, nullV V, fn func(K) V) V {
        \tif m == nil || fn == nil {
        \t\tThrowNPE()
        \t}
        \tif cur, ok := m.GetOk(k); ok && cur != nullV {
        \t\treturn cur
        \t}
        \tnv := fn(k)
        \tif nv != nullV {
        \t\tm.PutVoid(k, nv)
        \t}
        \treturn nv
        }

        func PlainMapComputeIfAbsentExtFn[K comparable, V comparable](m *PlainMap[K, V], k K, nullV V, fn Ext_java_util_function_Function) V {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(fn)
        \tif cur, ok := m.GetOk(k); ok && cur != nullV {
        \t\treturn cur
        \t}
        \traw := fn.Apply(boxJavaElem(k))
        \tif raw == nil || isJavaNilAny(raw) {
        \t\treturn nullV
        \t}
        \tnv := unboxJavaElem[V](raw)
        \tif nv != nullV {
        \t\tm.PutVoid(k, nv)
        \t}
        \treturn nv
        }

        func PlainMapPtrGetStr[K comparable](m *PlainMap[K, string], k K) string {
        \tif v, ok := m.GetOk(k); ok {
        \t\treturn v
        \t}
        \treturn NullString
        }

        func PlainMapGetStr[K comparable](m PlainMap[K, string], k K) string {
        \treturn PlainMapPtrGetStr(&m, k)
        }

        func PlainMapPutStr[K comparable](m *PlainMap[K, string], k K, v string) string {
        \tif old, ok := m.PutOk(k, v); ok {
        \t\treturn old
        \t}
        \treturn NullString
        }

        func PlainMapRemoveStr[K comparable](m *PlainMap[K, string], k K) string {
        \tif old, ok := m.RemoveOk(k); ok {
        \t\treturn old
        \t}
        \treturn NullString
        }

        func PlainMapPutIfAbsentStr[K comparable](m *PlainMap[K, string], k K, v string) string {
        \tcur, ok := m.GetOk(k)
        \tif !ok || cur == NullString {
        \t\tm.PutVoid(k, v)
        \t\treturn NullString
        \t}
        \treturn cur
        }

        func PlainMapFromJavaMap[K comparable, V any](src any) PlainMap[K, V] {
        \tCheckNotNull(src)
        \tmp := src.(Ext_java_util_Map)
        \tm := NewPlainMapCap[K, V](mp.Size())
        \tfor cur := NewMapEntryCursor(mp); cur.NextEntry(); {
        \t\tm.PutVoid(unboxJavaElem[K](cur.GetKey()), unboxJavaElem[V](cur.GetValue()))
        \t}
        \treturn m
        }

        func PlainMapPtrFromJavaMap[K comparable, V any](src any) *PlainMap[K, V] {
        \tm := PlainMapFromJavaMap[K, V](src)
        \treturn &m
        }

        func PlainMapFromNodeKeyJavaMap[K ~uint32, V any](src any) PlainMap[K, V] {
        \tCheckNotNull(src)
        \tmp := src.(Ext_java_util_Map)
        \tm := NewPlainMapCap[K, V](mp.Size())
        \tfor cur := NewMapEntryCursor(mp); cur.NextEntry(); {
        \t\tm.PutVoid(K(cur.KeyNode()), unboxJavaElem[V](cur.GetValue()))
        \t}
        \treturn m
        }

        func PlainMapPtrFromNodeKeyJavaMap[K ~uint32, V any](src any) *PlainMap[K, V] {
        \tm := PlainMapFromNodeKeyJavaMap[K, V](src)
        \treturn &m
        }

        func (m *PlainMap[K, V]) PutAllJavaMap(src any) {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(src)
        \tmp := src.(Ext_java_util_Map)
        \tfor cur := NewMapEntryCursor(mp); cur.NextEntry(); {
        \t\tm.PutVoid(unboxJavaElem[K](cur.GetKey()), unboxJavaElem[V](cur.GetValue()))
        \t}
        }

        func PlainMapGetNodeAny[K ~uint32, V any](m *PlainMap[K, V], rawKey any) any {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tvar k K
        \tif rawKey != nil && !isJavaNilAny(rawKey) {
        \t\tif nk, ok := rawKey.(K); ok {
        \t\t\tk = nk
        \t\t} else if rn, ok := rawKey.(jRawNodeProvider); ok {
        \t\t\tk = K(rn.RawNodeId_())
        \t\t} else if nd, ok := rawKey.(nodestore.Node); ok {
        \t\t\tk = K(nd)
        \t\t} else if u, ok := rawKey.(uint32); ok {
        \t\t\tk = K(u)
        \t\t} else {
        \t\t\treturn nil
        \t\t}
        \t}
        \tv, found := m.GetOk(k)
        \tif !found {
        \t\treturn nil
        \t}
        \treturn boxJavaElem(v)
        }

        func PlainMapGetAny[K comparable, V any](m *PlainMap[K, V], rawKey any) any {
        \tif m == nil {
        \t\tThrowNPE()
        \t}
        \tk, ok := tryUnboxOrderedMapKey[K](rawKey)
        \tif !ok {
        \t\treturn nil
        \t}
        \tv, found := m.GetOk(k)
        \tif !found {
        \t\treturn nil
        \t}
        \treturn boxJavaElem(v)
        }

        // PlainSet is an un-ordered value struct backed solely by a native Go map[T]struct{}
        // for non-iterated HashSet lowering.
        type PlainSet[T comparable] struct {
        \tm map[T]struct{}
        }

        func NewPlainSet[T comparable](capHint int) *PlainSet[T] {
        \tif capHint < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", capHint)))
        \t}
        \tif capHint > 0 {
        \t\treturn &PlainSet[T]{m: make(map[T]struct{}, capHint)}
        \t}
        \treturn &PlainSet[T]{}
        }

        func NewPlainSetVal[T comparable]() PlainSet[T] {
        \treturn PlainSet[T]{}
        }

        func NewPlainSetCap[T comparable](capHint int32) PlainSet[T] {
        \tif capHint < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", capHint)))
        \t}
        \tif capHint > 0 {
        \t\treturn PlainSet[T]{m: make(map[T]struct{}, capHint)}
        \t}
        \treturn PlainSet[T]{}
        }

        func NewPlainSetPtr[T comparable]() *PlainSet[T] {
        \treturn &PlainSet[T]{}
        }

        func NewPlainSetPtrCap[T comparable](capHint int32) *PlainSet[T] {
        \treturn NewPlainSet[T](int(capHint))
        }

        func (s *PlainSet[T]) Init() {}

        func (s *PlainSet[T]) InitCap(capHint int) {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif capHint < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException(fmt.Sprintf("Illegal initial capacity: %d", capHint)))
        \t}
        \tif capHint > 0 && s.m == nil {
        \t\ts.m = make(map[T]struct{}, capHint)
        \t}
        }

        func (s *PlainSet[T]) IsNil() bool     { return s == nil }
        func (s *PlainSet[T]) IsJavaNil_() bool { return s == nil }

        func (s *PlainSet[T]) Size() int32 {
        \tif s == nil {
        \t\treturn 0
        \t}
        \treturn int32(len(s.m))
        }

        func (s *PlainSet[T]) IsEmpty() bool {
        \treturn s == nil || len(s.m) == 0
        }

        func (s *PlainSet[T]) Contains(v T) bool {
        \tif s == nil || len(s.m) == 0 {
        \t\treturn false
        \t}
        \t_, ok := s.m[v]
        \treturn ok
        }

        func (s *PlainSet[T]) Add(v T) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif s.m == nil {
        \t\ts.m = make(map[T]struct{})
        \t\ts.m[v] = struct{}{}
        \t\treturn true
        \t}
        \tif _, ok := s.m[v]; ok {
        \t\treturn false
        \t}
        \ts.m[v] = struct{}{}
        \treturn true
        }

        func (s *PlainSet[T]) AddVoid(v T) {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif s.m == nil {
        \t\ts.m = make(map[T]struct{})
        \t}
        \ts.m[v] = struct{}{}
        }

        func (s *PlainSet[T]) Remove(v T) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tif len(s.m) == 0 {
        \t\treturn false
        \t}
        \tif _, ok := s.m[v]; ok {
        \t\tdelete(s.m, v)
        \t\treturn true
        \t}
        \treturn false
        }

        func (s *PlainSet[T]) Clear() {
        \tif s != nil && len(s.m) > 0 {
        \t\tclear(s.m)
        \t}
        }

        func PlainSetFromStrCol(col any) PlainSet[string] {
        \tvar s PlainSet[string]
        \tPlainSetAddAllStr(&s, col)
        \treturn s
        }

        func PlainSetPtrFromStrCol(col any) *PlainSet[string] {
        \ts := PlainSetFromStrCol(col)
        \treturn &s
        }

        func PlainSetFromNodeCol[N ~uint32](col any) PlainSet[N] {
        \tvar s PlainSet[N]
        \tPlainSetAddAllNode(&s, col)
        \treturn s
        }

        func PlainSetPtrFromNodeCol[N ~uint32](col any) *PlainSet[N] {
        \ts := PlainSetFromNodeCol[N](col)
        \treturn &s
        }

        func PlainSetFromCol[T comparable](col any) PlainSet[T] {
        \tvar s PlainSet[T]
        \tPlainSetAddAllCol(&s, col)
        \treturn s
        }

        func PlainSetPtrFromCol[T comparable](col any) *PlainSet[T] {
        \ts := PlainSetFromCol[T](col)
        \treturn &s
        }

        func PlainSetAddAllStr(s *PlainSet[string], col any) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(col)
        \tchanged := false
        \tif ssi, ok := col.(StrSliceIterable); ok {
        \t\tif ss, ok2 := ssi.IterStrSlice_(); ok2 {
        \t\t\tfor i := 0; i < len(ss); i++ {
        \t\t\t\tif s.Add(ss[i]) {
        \t\t\t\t\tchanged = true
        \t\t\t\t}
        \t\t\t}
        \t\t\treturn changed
        \t\t}
        \t}
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\tif s.Add(cur.NextStr()) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func PlainSetAddAllNode[N ~uint32](s *PlainSet[N], col any) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(col)
        \tchanged := false
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\tif s.Add(N(cur.NextNode())) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }

        func PlainSetAddAllCol[T comparable](s *PlainSet[T], col any) bool {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \tCheckNotNull(col)
        \tchanged := false
        \tfor cur := NewIterCursor(col); cur.HasNext(); {
        \t\tif s.Add(unboxJavaElem[T](cur.Next())) {
        \t\t\tchanged = true
        \t\t}
        \t}
        \treturn changed
        }
        """;
  }
}
