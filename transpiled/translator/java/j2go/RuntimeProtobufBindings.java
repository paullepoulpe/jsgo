package j2go;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.TreeMap;
import javax.lang.model.element.Element;
import javax.lang.model.element.ElementKind;
import javax.lang.model.element.ExecutableElement;
import javax.lang.model.element.Modifier;
import javax.lang.model.element.TypeElement;
import javax.lang.model.element.VariableElement;
import javax.lang.model.type.TypeMirror;

/**
 * Clean-room runtime bindings from {@code com.google.protobuf.*} and generated Closure protobuf
 * Java classes to {@code transpiled/shims/protobuf}.
 */
public final class RuntimeProtobufBindings {

  private RuntimeProtobufBindings() {}

  private static final Map<String, String> PROTO_MSG_TYPES =
      Map.ofEntries(
          Map.entry(
              "com.google.debugging.sourcemap.proto.Mapping.OriginalMapping",
              "Mapping_OriginalMapping"),
          Map.entry(
              "com.google.javascript.jscomp.AccessorSummaryEntryProto",
              "AccessorSummaryEntryProto"),
          Map.entry(
              "com.google.javascript.jscomp.AccessorSummaryProto", "AccessorSummaryProto"),
          Map.entry("com.google.javascript.jscomp.ChunkProto", "ChunkProto"),
          Map.entry("com.google.javascript.jscomp.ConformanceConfig", "ConformanceConfig"),
          Map.entry(
              "com.google.javascript.jscomp.JSCompilerStateProto", "JSCompilerStateProto"),
          Map.entry("com.google.javascript.jscomp.Requirement", "Requirement"),
          Map.entry(
              "com.google.javascript.jscomp.RequirementScopeEntry", "RequirementScopeEntry"),
          Map.entry("com.google.javascript.jscomp.UniqueIdProto", "UniqueIdProto"),
          Map.entry(
              "com.google.javascript.jscomp.VariableMapEntryProto", "VariableMapEntryProto"),
          Map.entry(
              "com.google.javascript.jscomp.instrumentation.reporter.proto.FileProfile",
              "FileProfile"),
          Map.entry(
              "com.google.javascript.jscomp.instrumentation.reporter.proto.InstrumentationPoint",
              "InstrumentationPoint"),
          Map.entry(
              "com.google.javascript.jscomp.instrumentation.reporter.proto.InstrumentationPointStats",
              "InstrumentationPointStats"),
          Map.entry(
              "com.google.javascript.jscomp.instrumentation.reporter.proto.ReportProfile",
              "ReportProfile"),
          Map.entry("com.google.javascript.jscomp.serialization.AstNode", "AstNode"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.ExternsSummary", "ExternsSummary"),
          Map.entry("com.google.javascript.jscomp.serialization.LazyAst", "LazyAst"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.ObjectTypeProto", "ObjectTypeProto"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.OptimizationJsdoc",
              "OptimizationJsdoc"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.SourceFilePool", "SourceFilePool"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.SourceFileProto", "SourceFileProto"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.SourceFileProto.FileOnDisk",
              "SourceFileProto_FileOnDisk"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.SourceFileProto.ZipEntryOnDisk",
              "SourceFileProto_ZipEntryOnDisk"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.StringPoolProto", "StringPoolProto"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.SubtypingEdge", "SubtypingEdge"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.TemplateStringValue",
              "TemplateStringValue"),
          Map.entry("com.google.javascript.jscomp.serialization.TypePool", "TypePool"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.TypePool.DebugInfo",
              "TypePool_DebugInfo"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.TypePool.DebugInfo.Mismatch",
              "TypePool_DebugInfo_Mismatch"),
          Map.entry("com.google.javascript.jscomp.serialization.TypeProto", "TypeProto"),
          Map.entry("com.google.javascript.jscomp.serialization.TypedAst", "TypedAst"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.TypedAst.List", "TypedAst_List"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.UnionTypeProto", "UnionTypeProto"));

  private static final Map<String, String> PROTO_ENUM_TYPES =
      Map.ofEntries(
          Map.entry(
              "com.google.debugging.sourcemap.proto.Mapping.OriginalMapping.Precision",
              "Mapping_OriginalMapping_Precision"),
          Map.entry(
              "com.google.javascript.jscomp.ConformanceConfig.LibraryLevelNonAllowlistedConformanceViolationsBehavior",
              "ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior"),
          Map.entry("com.google.javascript.jscomp.FeatureProto", "FeatureProto"),
          Map.entry("com.google.javascript.jscomp.LifeCycleStageProto", "LifeCycleStageProto"),
          Map.entry(
              "com.google.javascript.jscomp.PropertyAccessKindProto",
              "PropertyAccessKindProto"),
          Map.entry(
              "com.google.javascript.jscomp.Requirement.Severity", "Requirement_Severity"),
          Map.entry("com.google.javascript.jscomp.Requirement.Type", "Requirement_Type"),
          Map.entry(
              "com.google.javascript.jscomp.Requirement.TypeMatchingStrategy",
              "Requirement_TypeMatchingStrategy"),
          Map.entry(
              "com.google.javascript.jscomp.RequirementScopeEntry.Reason",
              "RequirementScopeEntry_Reason"),
          Map.entry(
              "com.google.javascript.jscomp.instrumentation.reporter.proto.InstrumentationPoint.Type",
              "InstrumentationPoint_Type"),
          Map.entry(
              "com.google.javascript.jscomp.instrumentation.reporter.proto.InstrumentationPointStats.Presence",
              "InstrumentationPointStats_Presence"),
          Map.entry("com.google.javascript.jscomp.serialization.JsdocTag", "JsdocTag"),
          Map.entry("com.google.javascript.jscomp.serialization.NodeKind", "NodeKind"),
          Map.entry("com.google.javascript.jscomp.serialization.NodeProperty", "NodeProperty"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.PrimitiveType", "PrimitiveType"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.SourceFileProto.LoaderCase",
              "SourceFileProto_LoaderCase"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.SourceFileProto.SourceKind",
              "SourceFileProto_SourceKind"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.TypeProto.KindCase",
              "TypeProto_KindCase"),
          Map.entry(
              "com.google.javascript.jscomp.serialization.AstNode.ValueCase",
              "AstNode_ValueCase"));

  private static String baseMsgQn(String qn) {
    if (qn.endsWith(".Builder")) {
      String base = qn.substring(0, qn.length() - ".Builder".length());
      if (PROTO_MSG_TYPES.containsKey(base)) {
        return base;
      }
    }
    return PROTO_MSG_TYPES.containsKey(qn) ? qn : null;
  }

  private static boolean hasToBuilderMethod(String msgGo) {
    return msgGo.equals("Requirement")
        || msgGo.equals("JSCompilerStateProto")
        || msgGo.equals("Mapping_OriginalMapping");
  }

  private static boolean isExtTypeReferenced(RuntimeBindings rb, String qn) {
    TypeElement te = rb.resolver().elements.getTypeElement(qn);
    return te != null && rb.resolver().getReferencedExternalTypes().contains(te);
  }

  private static boolean isMethodNeededOnChain(
      RuntimeBindings rb, TypeElement te, String simpleName, int arity) {
    List<TypeElement> chain = new ArrayList<>();
    chain.add(te);
    chain.addAll(rb.resolver().getAllAncestorTypes(te));
    for (TypeElement cur : chain) {
      for (Element enc : cur.getEnclosedElements()) {
        if (enc instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)
            && !ee.getModifiers().contains(Modifier.PRIVATE)
            && ee.getSimpleName().contentEquals(simpleName)
            && ee.getParameters().size() == arity) {
          TypeResolver.MethodSig sig = rb.resolver().getInstanceMethodSig(ee);
          if (sig != null && rb.resolver().isNeededExternalSig(sig)) {
            return true;
          }
        }
      }
    }
    return false;
  }

  private static boolean shouldCacheIndexedList(
      RuntimeBindings rb, TypeElement te, String idxGetter, String listGetter, String elemQn) {
    if (!isExtTypeReferenced(rb, elemQn)) {
      return false;
    }
    return isMethodNeededOnChain(rb, te, idxGetter, 1)
        || isMethodNeededOnChain(rb, te, listGetter, 0);
  }

  private static boolean isStringListGetter(String goMethod) {
    return goMethod.equals("GetInputIdsList")
        || goMethod.equals("GetExportedNamesList")
        || goMethod.equals("GetCssNamesList")
        || goMethod.equals("GetExternsList")
        || goMethod.equals("GetInjectedLibrariesList")
        || goMethod.equals("GetWhitelistList")
        || goMethod.equals("GetWhitelistRegexpList")
        || goMethod.equals("GetAllowlistList")
        || goMethod.equals("GetAllowlistRegexpList")
        || goMethod.equals("GetOnlyApplyToList")
        || goMethod.equals("GetOnlyApplyToRegexpList")
        || goMethod.equals("GetValueList")
        || goMethod.equals("GetConfigFileList")
        || goMethod.equals("GetPrefixList")
        || goMethod.equals("GetRegexpList")
        || goMethod.equals("GetRuntimeLibraryToInjectList");
  }

  private static boolean isInt32ListGetter(String goMethod) {
    return goMethod.equals("GetDependenciesList")
        || goMethod.equals("GetPropNamePtrList")
        || goMethod.equals("GetInstanceTypeList")
        || goMethod.equals("GetPrototypeList")
        || goMethod.equals("GetInvolvedColorList")
        || goMethod.equals("GetUnionMemberList");
  }

  private static String subMsgListElemExt(String goMethod) {
    return switch (goMethod) {
      case "GetAccessorsList" -> "Ext_com_google_javascript_jscomp_AccessorSummaryEntryProto";
      case "GetRequirementList" -> "Ext_com_google_javascript_jscomp_Requirement";
      case "GetChunksList" -> "Ext_com_google_javascript_jscomp_ChunkProto";
      case "GetUniqueIdSupplierList" -> "Ext_com_google_javascript_jscomp_UniqueIdProto";
      case "GetStringMapList", "GetInstrumentationMappingList" ->
          "Ext_com_google_javascript_jscomp_VariableMapEntryProto";
      case "GetWhitelistEntryList", "GetAllowlistEntryList" ->
          "Ext_com_google_javascript_jscomp_RequirementScopeEntry";
      case "GetInstrumentationPointsStatsList" ->
          "Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats";
      case "GetFileProfileList" ->
          "Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_FileProfile";
      case "GetSourceFileList" ->
          "Ext_com_google_javascript_jscomp_serialization_SourceFileProto";
      case "GetTypeList" -> "Ext_com_google_javascript_jscomp_serialization_TypeProto";
      case "GetDisambiguationEdgesList" ->
          "Ext_com_google_javascript_jscomp_serialization_SubtypingEdge";
      case "GetMismatchList" ->
          "Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Mismatch";
      case "GetExternAstList", "GetCodeAstList" ->
          "Ext_com_google_javascript_jscomp_serialization_LazyAst";
      case "GetChildList" -> "Ext_com_google_javascript_jscomp_serialization_AstNode";
      case "GetTypedAstsList" -> "Ext_com_google_javascript_jscomp_serialization_TypedAst";
      default -> null;
    };
  }

  private static boolean isOtherListGetter(String msgGo, String goMethod) {
    if (msgGo.equals("StringPoolProto") && goMethod.equals("GetStringsList")) {
      return true;
    }
    if (isInt32ListGetter(goMethod)) {
      return true;
    }
    if (msgGo.equals("JSCompilerStateProto") && goMethod.equals("GetAllowableFeaturesList")) {
      return true;
    }
    if (msgGo.equals("OptimizationJsdoc") && goMethod.equals("GetKindList")) {
      return true;
    }
    return subMsgListElemExt(goMethod) != null;
  }

  private static boolean isDefaultSingletonSubMsgGetter(String msgGo, TypeResolver.MethodSig sig) {
    String m = sig.simpleName;
    if (sig.canonParamTypes.size() != 0) {
      return false;
    }
    if (!m.startsWith("get") || m.equals("getDefaultInstanceForType")) {
      return false;
    }
    if (m.endsWith("List") || m.endsWith("Count") || m.endsWith("Builder")) {
      return false;
    }
    if (msgGo.equals("LazyAst") && m.equals("getScript")) {
      return false;
    }
    if (msgGo.equals("ObjectTypeProto") && m.equals("getUuid")) {
      return false;
    }
    TypeElement retTe = elementOfReturn(sig);
    if (retTe == null) {
      return false;
    }
    return PROTO_MSG_TYPES.containsKey(retTe.getQualifiedName().toString());
  }

  private static boolean isCacheableSubMsgGetter(String msgGo, TypeResolver.MethodSig sig) {
    // AstNode is instantiated 35M+ times and visited once during deserialization; omitting per-node
    // atomic.Pointer fields keeps every Ext_..._AstNode at 16 bytes while unset sub-messages still
    // return shared protoDefSub_AstNode_* singletons in 0 allocations.
    if (msgGo.equals("AstNode")) {
      return false;
    }
    return isDefaultSingletonSubMsgGetter(msgGo, sig);
  }

  /**
   * Returns a deterministic map of {@code goMethod -> retQn} for singular sub-message getters
   * emitted on an immutable protobuf message wrapper (used for default sub-message singletons).
   */
  private static Map<String, String> collectDefaultSubMsgGetters(
      RuntimeBindings rb, TypeElement msgTe, String msgGo) {
    Map<String, String> result = new LinkedHashMap<>();
    Set<String> seenGoNames = new HashSet<>();
    List<TypeElement> chain = new ArrayList<>();
    chain.add(msgTe);
    chain.addAll(rb.resolver().getAllAncestorTypes(msgTe));
    for (TypeElement cur : chain) {
      for (Element enc : cur.getEnclosedElements()) {
        if (enc instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)
            && !ee.getModifiers().contains(Modifier.PRIVATE)) {
          TypeResolver.MethodSig sig = rb.resolver().getInstanceMethodSig(ee);
          if (sig != null && rb.resolver().isNeededExternalSig(sig) && seenGoNames.add(sig.goName)) {
            if (isDefaultSingletonSubMsgGetter(msgGo, sig)) {
              TypeElement retTe = elementOfReturn(sig);
              String retQn = retTe.getQualifiedName().toString();
              if (isExtTypeReferenced(rb, retQn)) {
                result.putIfAbsent(GoNames.exported(sig.simpleName), retQn);
              }
            }
          }
        }
      }
    }
    return result;
  }

  /**
   * Emits value-embedded {@code inlineRaw} + {@code builtMsg} on {@code .Builder} structs and
   * {@code atomic.Pointer} wrapper caches on immutable protobuf message structs.
   */
  public static void emitExtraStructFields(
      RuntimeBindings rb,
      TypeElement te,
      PackagePlanner.ImportTracker imports,
      StringBuilder out) {
    String qn = te.getQualifiedName().toString();
    if (qn.endsWith(".Builder")) {
      String baseQn = baseMsgQn(qn);
      if (baseQn != null) {
        String msgGo = PROTO_MSG_TYPES.get(baseQn);
        String pbPkg =
            imports.ext(
                "github.com/paullepoulpe/jsgo/transpiled/shims/protobuf", "protobuf");
        if (isExtTypeReferenced(rb, baseQn)) {
          String msgExtStruct = "Ext_" + GoNames.cleanChars(baseQn.replace('.', '_'));
          out.append("\tbuiltMsg ").append(msgExtStruct).append(" `gson:\"-\"`\n");
        }
        out.append("\tinlineRaw ").append(pbPkg).append('.').append(msgGo).append(" `gson:\"-\"`\n");
        out.append("\tRaw_ *").append(pbPkg).append('.').append(msgGo).append(" `gson:\"-\"`\n");
        out.append("\tbuilt_ bool `gson:\"-\"`\n");
        if (isExtTypeReferenced(rb, baseQn)) {
          out.append("\tusedBuiltMsg_ bool `gson:\"-\"`\n");
        }
        if (qn.equals("com.google.javascript.jscomp.serialization.AstNode.Builder")) {
          String selfStruct = "Ext_" + GoNames.cleanChars(qn.replace('.', '_'));
          out.append("\tpoolChildHead_ *").append(selfStruct).append(" `gson:\"-\"`\n");
          out.append("\tpoolNextSibling_ *").append(selfStruct).append(" `gson:\"-\"`\n");
          out.append("\tfromPool_ bool `gson:\"-\"`\n");
          out.append("\tlinked_ bool `gson:\"-\"`\n");
          out.append("\tpinned_ bool `gson:\"-\"`\n");
        }
      }
      return;
    }

    String msgGo = PROTO_MSG_TYPES.get(qn);
    if (msgGo == null) {
      return;
    }

    Set<String> emittedFields = new LinkedHashSet<>();
    if (msgGo.equals("TypePool")
        && shouldCacheIndexedList(
            rb,
            te,
            "getType",
            "getTypeList",
            "com.google.javascript.jscomp.serialization.TypeProto")) {
      String atomicPkg = imports.std("sync/atomic");
      emittedFields.add("cache_GetTypeList");
      out.append("\tcache_GetTypeList ")
          .append(atomicPkg)
          .append(".Pointer[jProtoList] `gson:\"-\"`\n");
    }
    if (msgGo.equals("JSCompilerStateProto")
        && shouldCacheIndexedList(
            rb,
            te,
            "getChunks",
            "getChunksList",
            "com.google.javascript.jscomp.ChunkProto")) {
      String atomicPkg = imports.std("sync/atomic");
      emittedFields.add("cache_GetChunksList");
      out.append("\tcache_GetChunksList ")
          .append(atomicPkg)
          .append(".Pointer[jProtoList] `gson:\"-\"`\n");
    }

    Set<String> seenGoNames = new HashSet<>();
    List<TypeElement> chain = new ArrayList<>();
    chain.add(te);
    chain.addAll(rb.resolver().getAllAncestorTypes(te));
    for (TypeElement cur : chain) {
      for (Element enc : cur.getEnclosedElements()) {
        if (enc instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)
            && !ee.getModifiers().contains(Modifier.PRIVATE)) {
          TypeResolver.MethodSig sig = rb.resolver().getInstanceMethodSig(ee);
          if (sig == null || !rb.resolver().isNeededExternalSig(sig) || !seenGoNames.add(sig.goName)) {
            continue;
          }
          String m = sig.simpleName;
          int arity = sig.canonParamTypes.size();
          String goMethod = GoNames.exported(m);
          if (arity == 0 && m.startsWith("get") && m.endsWith("List")) {
            if (isStringListGetter(goMethod)) {
              if (emittedFields.add("cache_" + goMethod)) {
                String atomicPkg = imports.std("sync/atomic");
                out.append("\tcache_")
                    .append(goMethod)
                    .append(' ')
                    .append(atomicPkg)
                    .append(".Pointer[jProtoStringList] `gson:\"-\"`\n");
              }
            } else if (msgGo.equals("StringPoolProto") && goMethod.equals("GetStringsList")) {
              if (emittedFields.add("cache_" + goMethod)) {
                String atomicPkg = imports.std("sync/atomic");
                out.append("\tcache_")
                    .append(goMethod)
                    .append(' ')
                    .append(atomicPkg)
                    .append(".Pointer[jProtoByteStringList] `gson:\"-\"`\n");
              }
            } else if (isOtherListGetter(msgGo, goMethod)) {
              if (emittedFields.add("cache_" + goMethod)) {
                String atomicPkg = imports.std("sync/atomic");
                out.append("\tcache_")
                    .append(goMethod)
                    .append(' ')
                    .append(atomicPkg)
                    .append(".Pointer[jProtoList] `gson:\"-\"`\n");
              }
            }
          } else if (isCacheableSubMsgGetter(msgGo, sig)) {
            TypeElement retTe = elementOfReturn(sig);
            String retQn = retTe.getQualifiedName().toString();
            if (isExtTypeReferenced(rb, retQn)) {
              String retStruct = "Ext_" + GoNames.cleanChars(retQn.replace('.', '_'));
              if (emittedFields.add("cache_" + goMethod)) {
                String atomicPkg = imports.std("sync/atomic");
                out.append("\tcache_")
                    .append(goMethod)
                    .append(' ')
                    .append(atomicPkg)
                    .append(".Pointer[")
                    .append(retStruct)
                    .append("] `gson:\"-\"`\n");
              }
            }
          }
        }
      }
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
    if (qn.equals("com.google.protobuf.ByteString") && fn.equals("EMPTY")) {
      out.append("var ")
          .append(varGoName)
          .append(' ')
          .append(fieldType)
          .append(" = wrapByteString(protobuf.ByteString_EMPTY)\n");
      return true;
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

    if (qn.equals("com.google.protobuf.WireFormat")) {
      if (simple.equals("getTagFieldNumber") && arity == 1) {
        out.append("\treturn int32(uint32(arg0) >> 3)\n}\n\n");
        return true;
      }
      if (simple.equals("getTagWireType") && arity == 1) {
        out.append("\treturn arg0 & 7\n}\n\n");
        return true;
      }
    }

    String msgGo = PROTO_MSG_TYPES.get(qn);
    if (msgGo != null) {
      if (simple.equals("newBuilder") && arity == 0) {
        if (qn.equals("com.google.javascript.jscomp.serialization.AstNode")) {
          TypeElement bldTe = rb.resolver().elements.getTypeElement(qn + ".Builder");
          if (bldTe != null && rb.resolver().getReferencedExternalTypes().contains(bldTe)) {
            out.append("\treturn allocAstNodeBuilder_()\n}\n\n");
            return true;
          }
        }
        String bStruct = retType.startsWith("*") ? retType.substring(1) : retType;
        out.append("\tb := &").append(bStruct).append("{}\n");
        out.append("\tb.Raw_ = &b.inlineRaw\n");
        out.append("\tb.Core_ = b.Raw_\n");
        out.append("\treturn b\n}\n\n");
        return true;
      }
      if (simple.equals("newBuilder") && arity == 1) {
        String bStruct = retType.startsWith("*") ? retType.substring(1) : retType;
        out.append("\tCheckNotNullPtr(arg0)\n");
        if (hasToBuilderMethod(msgGo)) {
          out.append("\traw := arg0.p_").append(msgGo).append("().ToBuilder()\n");
          out.append("\tb := &").append(bStruct).append("{Raw_: raw}\n");
          out.append("\tb.Core_ = raw\n");
          out.append("\treturn b\n}\n\n");
        } else {
          out.append("\tb := &").append(bStruct).append("{}\n");
          out.append("\tb.inlineRaw = *arg0.p_").append(msgGo).append("()\n");
          out.append("\tb.Raw_ = &b.inlineRaw\n");
          out.append("\tb.Core_ = b.Raw_\n");
          out.append("\treturn b\n}\n\n");
        }
        return true;
      }
      if (simple.equals("getDefaultInstance") && arity == 0) {
        out.append("\treturn protoDefault_").append(msgGo).append("_\n}\n\n");
        return true;
      }
      if (qn.equals("com.google.javascript.jscomp.serialization.AstNode")
          && simple.equals("parseFrom")
          && arity == 2) {
        String mStruct = retType.startsWith("*") ? retType.substring(1) : retType;
        out.append("\t_ = arg1\n");
        out.append("\tn, err := protobuf.AstNode_ParseFromCodedInput(unwrapCodedInputStream(arg0), nil)\n");
        out.append("\tcheckProtoErr(err)\n");
        out.append("\tm := &").append(mStruct).append("{}\n");
        out.append("\tm.Core_ = n\n");
        out.append("\treturn m\n}\n\n");
        return true;
      }
      if (qn.equals("com.google.javascript.jscomp.JSCompilerStateProto")
          && simple.equals("parseDelimitedFrom")
          && arity == 1) {
        String mStruct = retType.startsWith("*") ? retType.substring(1) : retType;
        out.append("\tp := runProtoOp(func() *protobuf.JSCompilerStateProto {\n");
        out.append("\t\treturn protobuf.JSCompilerStateProto_ParseDelimitedFrom(asInputStream(arg0))\n");
        out.append("\t})\n");
        out.append("\tif p == nil {\n\t\treturn nil\n\t}\n");
        out.append("\tm := &").append(mStruct).append("{}\n");
        out.append("\tm.Core_ = p\n");
        out.append("\treturn m\n}\n\n");
        return true;
      }
    }

    String enumGo = PROTO_ENUM_TYPES.get(qn);
    if (enumGo != null) {
      if (simple.equals("forNumber") && arity == 1) {
        if (enumGo.equals("NodeProperty")) {
          out.append("\tp := protobuf.NodeProperty_ForNumber(arg0)\n");
          out.append("\tif p == nil {\n\t\treturn nil\n\t}\n");
          out.append("\treturn lookupProtoEnum_NodeProperty(int32(*p))\n}\n\n");
          return true;
        }
      }
    }

    if (qn.equals("com.google.protobuf.ByteString")) {
      if (simple.equals("copyFrom") && arity == 1) {
        String p0 = rb.goTypeInJ2gort(m.getParameters().get(0).asType(), imports);
        if (p0.equals("[]int8")) {
          out.append("\tCheckNotNullSlice(arg0)\n");
          out.append("\treturn wrapByteString(protobuf.ByteString_CopyFrom(int8SliceToBytes(arg0)))\n}\n\n");
          return true;
        }
      }
      if (simple.equals("copyFrom") && arity == 3) {
        out.append("\tCheckNotNullSlice(arg0)\n");
        out.append("\traw := int8SliceToBytes(arg0)\n");
        out.append("\treturn wrapByteString(protobuf.ByteString_CopyFrom(raw[arg1 : arg1+arg2]))\n}\n\n");
        return true;
      }
      if (simple.equals("copyFromUtf8") && arity == 1) {
        out.append("\tCheckNotNullString(arg0)\n");
        out.append("\treturn wrapByteString(protobuf.ByteString_CopyFromUtf8(arg0))\n}\n\n");
        return true;
      }
      if (simple.equals("newOutput") && arity == 1) {
        out.append("\tout := &Ext_com_google_protobuf_ByteString_Output{}\n");
        out.append("\tout.BindExtSelf_(out)\n");
        out.append("\tout.Core_ = protobuf.ByteString_NewOutput(int(arg0))\n");
        out.append("\treturn out\n}\n\n");
        return true;
      }
      if (simple.equals("newOutput") && arity == 0) {
        out.append("\tout := &Ext_com_google_protobuf_ByteString_Output{}\n");
        out.append("\tout.BindExtSelf_(out)\n");
        out.append("\tout.Core_ = protobuf.ByteString_NewOutput(128)\n");
        out.append("\treturn out\n}\n\n");
        return true;
      }
    }

    if (qn.equals("com.google.protobuf.CodedInputStream")) {
      if (simple.equals("newInstance") && arity == 1) {
        String p0 = rb.goTypeInJ2gort(m.getParameters().get(0).asType(), imports);
        if (p0.equals("[]int8")) {
          out.append("\tCheckNotNullSlice(arg0)\n");
          out.append("\treturn wrapCodedInputStream(protobuf.CodedInputStream_NewInstanceBytes(int8SliceToBytes(arg0)))\n}\n\n");
          return true;
        }
        out.append("\tCheckNotNull(arg0)\n");
        out.append("\tcis, err := protobuf.CodedInputStream_NewInstance(asInputStream(arg0))\n");
        out.append("\tcheckProtoErr(err)\n");
        out.append("\treturn wrapCodedInputStream(cis)\n}\n\n");
        return true;
      }
    }

    if ((qn.equals("com.google.protobuf.ExtensionRegistry")
            || qn.equals("com.google.protobuf.ExtensionRegistryLite"))
        && simple.equals("getEmptyRegistry")
        && arity == 0) {
      String rStruct = retType.startsWith("*") ? retType.substring(1) : retType;
      out.append("\tr := &").append(rStruct).append("{}\n");
      out.append("\tr.Core_ = protobuf.ExtensionRegistry_GetEmptyRegistry()\n");
      out.append("\treturn r\n}\n\n");
      return true;
    }

    if (qn.equals("com.google.protobuf.TextFormat")) {
      if (simple.equals("merge") && arity == 2) {
        out.append("\ttxt := readTextFormatInput(arg0)\n");
        out.append("\tCheckNotNull(arg1)\n");
        out.append("\ttarget := unwrapProtoMutCore(arg1).(protobuf.ConformanceTextFormatTarget)\n");
        out.append("\tif pe := protobuf.TextFormat_MergeErr(txt, target); pe != nil {\n");
        out.append("\t\tpanic(javaio.NewIOException(pe.Error()))\n");
        out.append("\t}\n}\n\n");
        return true;
      }
      if (simple.equals("printer") && arity == 0) {
        out.append("\tp := &Ext_com_google_protobuf_TextFormat_Printer{}\n");
        out.append("\tp.Core_ = protobuf.TextFormat_PrinterInstance()\n");
        out.append("\treturn p\n}\n\n");
        return true;
      }
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

    // 1. Protobuf enums
    String enumGo = PROTO_ENUM_TYPES.get(qn);
    if (enumGo != null) {
      if (m.equals("getNumber") && arity == 0) {
        return "\tCheckNotNullPtr(this_)\n\treturn protobuf." + enumGo + "(this_.Ordinal()).GetNumber()\n";
      }
      if (m.equals("equals") && arity == 1) {
        return "\treturn any(this_) == arg0\n";
      }
      if (m.equals("compareTo") && arity == 1) {
        return "\treturn EnumCompareTo(this_, arg0)\n";
      }
    }

    // 2. Core com.google.protobuf classes
    if (qn.equals("com.google.protobuf.ByteString")) {
      if (m.equals("size") && arity == 0) return "\treturn int32(this_.bs().Size())\n";
      if (m.equals("isEmpty") && arity == 0) return "\treturn this_.bs().IsEmpty()\n";
      if (m.equals("byteAt") && arity == 1) return "\treturn int8(this_.bs().ByteAt(int(arg0)))\n";
      if (m.equals("toByteArray") && arity == 0) return "\treturn bytesToInt8Slice(this_.bs().ToByteArray())\n";
      if (m.equals("toStringUtf8") && arity == 0) return "\treturn this_.bs().ToStringUtf8()\n";
      if (m.equals("newCodedInput") && arity == 0) return "\treturn wrapCodedInputStream(this_.bs().NewCodedInput())\n";
      if (m.equals("equals") && arity == 1) return "\treturn byteStringEquals(this_.bs(), arg0)\n";
      if (m.equals("hashCode") && arity == 0) return "\treturn this_.bs().HashCode()\n";
      if (m.equals("toString") && arity == 0) return "\treturn this_.bs().String()\n";
      if (m.equals("iterator") && arity == 0) return "\treturn newByteStringIter(this_.bs())\n";
      if (m.equals("forEach") && arity == 1) return "\tbyteStringForEach(this_.bs(), arg0)\n";
    }

    if (qn.equals("com.google.protobuf.ByteString.Output")) {
      if (m.equals("write") && arity == 1) {
        String p0 = rb.goTypeInJ2gort(sig.canonParamTypes.get(0), imports);
        if (p0.equals("int32")) {
          return "\tthis_.Core_.(*protobuf.ByteString_Output).WriteByteInt(int(arg0))\n";
        }
        if (p0.equals("[]int8")) {
          return "\t_, _ = this_.Core_.(*protobuf.ByteString_Output).Write(int8SliceToBytes(CheckNotNullSlice(arg0)))\n";
        }
      }
      if (m.equals("write") && arity == 3) {
        return "\traw := int8SliceToBytes(CheckNotNullSlice(arg0))\n\t_, _ = this_.Core_.(*protobuf.ByteString_Output).Write(raw[arg1 : arg1+arg2])\n";
      }
      if (m.equals("size") && arity == 0) {
        return "\treturn int32(this_.Core_.(*protobuf.ByteString_Output).Size())\n";
      }
      if (m.equals("toByteString") && arity == 0) {
        return "\treturn wrapByteString(this_.Core_.(*protobuf.ByteString_Output).ToByteString())\n";
      }
      if ((m.equals("flush") || m.equals("close")) && arity == 0) {
        return "";
      }
    }

    if (qn.equals("com.google.protobuf.CodedInputStream")) {
      if (m.equals("readTag") && arity == 0) {
        return "\ttag, err := this_.cis().ReadTag()\n\tcheckProtoErr(err)\n\treturn tag\n";
      }
      if (m.equals("readMessage") && arity == 2) {
        return "\t_ = arg1\n\tCheckNotNull(arg0)\n\tmb := unwrapProtoMutCore(arg0).(protobuf.MessageBuilder)\n\terr := this_.cis().ReadMessage(mb, nil)\n\tcheckProtoErr(err)\n";
      }
      if (m.equals("setRecursionLimit") && arity == 1) {
        return "\treturn int32(this_.cis().SetRecursionLimit(int(arg0)))\n";
      }
      if (m.equals("resetSizeCounter") && arity == 0) {
        return "\tthis_.cis().ResetSizeCounter()\n";
      }
      if (m.equals("isAtEnd") && arity == 0) {
        return "\treturn this_.cis().IsAtEnd()\n";
      }
    }

    if (qn.equals("com.google.protobuf.TextFormat.Printer")) {
      if (m.equals("printToString") && arity == 1) {
        return "\tif arg0 == nil || javalang.IsNil(arg0) { return \"\" }\n\treturn this_.Core_.(*protobuf.TextFormat_Printer).PrintToString(unwrapProtoReadCore(arg0))\n";
      }
    }

    if (qn.equals("com.google.protobuf.Descriptors.FieldDescriptor")
        || qn.equals("com.google.protobuf.Descriptors.GenericDescriptor")) {
      if (m.equals("getName") && arity == 0) {
        return "\treturn this_.Core_.(protoFieldDesc).name\n";
      }
      if (m.equals("compareTo") && arity == 1) {
        return "\ta := this_.Core_.(protoFieldDesc)\n\tb := arg0.(*Ext_com_google_protobuf_Descriptors_FieldDescriptor).Core_.(protoFieldDesc)\n\tif a.idx < b.idx { return -1 }\n\tif a.idx > b.idx { return 1 }\n\treturn 0\n";
      }
    }

    // Abstract / Generated message base classes
    if (qn.startsWith("com.google.protobuf.AbstractMessage")
        || qn.startsWith("com.google.protobuf.GeneratedMessage")) {
      if (m.equals("toByteString") && arity == 0) {
        return "\treturn protoMsgToByteString(unwrapProtoReadCore(this_))\n";
      }
      if (m.equals("writeTo") && arity == 1) {
        return "\tprotoMsgWriteTo(unwrapProtoReadCore(this_), arg0)\n";
      }
      if (m.equals("writeDelimitedTo") && arity == 1) {
        return "\tprotoMsgWriteDelimitedTo(unwrapProtoReadCore(this_), arg0)\n";
      }
      if (m.equals("equals") && arity == 1) {
        return "\treturn protoMsgEquals(this_, arg0)\n";
      }
      if (m.equals("getAllFields") && arity == 0) {
        return "\treturn protoGetAllFields(unwrapProtoReadCore(this_))\n";
      }
      if ((m.equals("clear") || m.equals("build") || m.equals("toBuilder")) && arity == 0) {
        return "\tpanic(\"UNIMPLEMENTED_EXT: abstract proto base method \" + "
            + GoNames.goStringLiteral(qn + "." + m)
            + ")\n";
      }
    }

    // 3. Concrete Protobuf Messages & Builders
    String baseQn = baseMsgQn(qn);
    if (baseQn != null) {
      boolean isBuilder = qn.endsWith(".Builder");
      String msgGo = PROTO_MSG_TYPES.get(baseQn);
      String msgExtStruct = "Ext_" + GoNames.cleanChars(baseQn.replace('.', '_'));
      String bldExtStruct = msgExtStruct + "_Builder";
      String pExpr = "this_.p_" + msgGo + "()";
      String pMutExpr = "this_.pMut_" + msgGo + "()";

      if (m.equals("toByteString") && arity == 0) {
        if (!isBuilder
            && qn.equals("com.google.javascript.jscomp.serialization.AstNode")
            && isExtTypeReferenced(rb, qn + ".Builder")) {
          return "\tbs := protoMsgToByteString(" + pExpr + ")\n"
              + "\tif b, ok := this_.Core_.(*" + bldExtStruct + "); ok && b.fromPool_ && !b.linked_ && !b.pinned_ {\n"
              + "\t\tthis_.Core_ = nil\n"
              + "\t\trecycleAstNodeBuilderTree_(b)\n"
              + "\t}\n"
              + "\treturn bs\n";
        }
        return "\treturn protoMsgToByteString(" + pExpr + ")\n";
      }
      if (m.equals("toByteArray") && arity == 0) {
        return "\treturn bytesToInt8Slice(protobuf.MessageToByteArray(" + pExpr + "))\n";
      }
      if (m.equals("writeTo") && arity == 1) {
        return "\tprotoMsgWriteTo(" + pExpr + ", arg0)\n";
      }
      if (m.equals("writeDelimitedTo") && arity == 1) {
        return "\tprotoMsgWriteDelimitedTo(" + pExpr + ", arg0)\n";
      }
      if (m.equals("equals") && arity == 1) {
        return "\treturn protoMsgEquals(this_, arg0)\n";
      }
      if (m.equals("hashCode") && arity == 0) {
        return "\tif h, ok := any(" + pExpr + ").(interface{ HashCode() int32 }); ok { return h.HashCode() }\n\treturn javalang.IdentityHashCode(this_)\n";
      }
      if (m.equals("toString") && arity == 0) {
        return "\tif s, ok := any(" + pExpr + ").(interface{ String() string }); ok { return s.String() }\n\treturn "
            + GoNames.goStringLiteral(qn)
            + "\n";
      }
      if (m.equals("getAllFields") && arity == 0) {
        return "\treturn protoGetAllFields(" + pExpr + ")\n";
      }
      if (m.equals("getDefaultInstanceForType") && arity == 0 && !isBuilder) {
        return "\tif this_ == nil { ThrowNPE() }\n\treturn protoDefault_" + msgGo + "_\n";
      }
      if (m.equals("clear") && arity == 0 && isBuilder) {
        String detachCond =
            isExtTypeReferenced(rb, baseQn) ? "(this_.built_ || this_.usedBuiltMsg_)" : "this_.built_";
        String pinAstNode =
            qn.equals("com.google.javascript.jscomp.serialization.AstNode.Builder")
                ? ("\t\tthis_.pinned_ = true\n"
                    + (isExtTypeReferenced(rb, baseQn)
                        ? "\t\tif this_.usedBuiltMsg_ && this_.builtMsg.Core_ == this_ { this_.builtMsg.Core_ = this_.Raw_ }\n"
                        : ""))
                : "";
        return "\tif this_ == nil { ThrowNPE() }\n"
            + "\tif " + detachCond + " {\n"
            + pinAstNode
            + "\t\tthis_.Raw_ = &protobuf." + msgGo + "{}\n"
            + "\t\tthis_.Core_ = this_.Raw_\n"
            + "\t\tthis_.built_ = false\n"
            + "\t} else {\n"
            + "\t\t*" + pMutExpr + " = protobuf." + msgGo + "{}\n"
            + "\t}\n"
            + "\treturn this_\n";
      }
      if (m.equals("build") && arity == 0 && isBuilder) {
        if (isExtTypeReferenced(rb, baseQn)) {
          return "\treturn this_.buildMsg_" + msgGo + "()\n";
        }
        return "\tout := &" + msgExtStruct + "{}\n\tout.Core_ = this_.buildRaw_" + msgGo + "()\n\treturn out\n";
      }
      if (m.equals("toBuilder") && arity == 0 && !isBuilder) {
        if (hasToBuilderMethod(msgGo)) {
          return "\traw := " + pExpr + ".ToBuilder()\n"
              + "\tout := &" + bldExtStruct + "{Raw_: raw}\n"
              + "\tout.Core_ = raw\n"
              + "\treturn out\n";
        }
        String pinAst =
            qn.equals("com.google.javascript.jscomp.serialization.AstNode")
                ? "\tif b, ok := this_.Core_.(*" + bldExtStruct + "); ok { b.pinned_ = true }\n"
                : "";
        return "\tsrc := " + pExpr + "\n"
            + pinAst
            + "\tout := &" + bldExtStruct + "{}\n"
            + "\tout.inlineRaw = *src\n"
            + "\tout.Raw_ = &out.inlineRaw\n"
            + "\tout.Core_ = out.Raw_\n"
            + "\treturn out\n";
      }

      // Special ByteString getters where Go shim GetX returns []byte and GetXByteString returns ByteString
      if (msgGo.equals("LazyAst") && m.equals("getScript") && arity == 0) {
        return "\treturn wrapByteString(" + pExpr + ".GetScriptByteString())\n";
      }
      if (msgGo.equals("ObjectTypeProto") && m.equals("getUuid") && arity == 0) {
        return "\treturn wrapByteString(" + pExpr + ".GetUuidByteString())\n";
      }
      if (msgGo.equals("StringPoolProto") && m.equals("getStringsList") && arity == 0) {
        if (!isBuilder) {
          return "\tCheckNotNullPtr(this_)\n"
              + "\tif p := this_.cache_GetStringsList.Load(); p != nil { return p }\n"
              + "\traw := " + pExpr + ".GetStringsByteStringList()\n"
              + "\tif len(raw) == 0 {\n"
              + "\t\tthis_.cache_GetStringsList.CompareAndSwap(nil, emptyProtoByteStringList_)\n"
              + "\t\treturn emptyProtoByteStringList_\n"
              + "\t}\n"
              + "\tl := newProtoByteStringList(raw)\n"
              + "\tif this_.cache_GetStringsList.CompareAndSwap(nil, l) { return l }\n"
              + "\treturn this_.cache_GetStringsList.Load()\n";
        }
        return "\treturn newProtoByteStringList(" + pExpr + ".GetStringsByteStringList())\n";
      }

      // Count & indexed getters where Go shim returns/takes `int`
      if (m.endsWith("Count") && m.startsWith("get") && arity == 0 && retType.equals("int32")) {
        if (m.equals("getChildCount")
            || m.equals("getKindCount")
            || m.equals("getOwnPropertyCount")
            || m.equals("getUnionMemberCount")
            || m.equals("getTypeCount")
            || msgGo.equals("Requirement")) {
          return "\treturn int32(" + pExpr + "." + GoNames.exported(m) + "())\n";
        }
        String listMethod = GoNames.exported(m.substring(0, m.length() - 5) + "List");
        return "\treturn int32(len(" + pExpr + "." + listMethod + "()))\n";
      }
      if (m.equals("getOwnProperty") && arity == 1 && retType.equals("int32")) {
        return "\treturn " + pExpr + ".GetOwnProperty(int(arg0))\n";
      }
      if (!isBuilder
          && msgGo.equals("TypePool")
          && m.equals("getType")
          && arity == 1
          && shouldCacheIndexedList(
              rb,
              te,
              "getType",
              "getTypeList",
              "com.google.javascript.jscomp.serialization.TypeProto")) {
        return "\tl := this_.getOrBuildTypeList_()\n"
            + "\tCheckElementIndex(arg0, int32(len(l.items)))\n"
            + "\treturn l.items[arg0].(*Ext_com_google_javascript_jscomp_serialization_TypeProto)\n";
      }
      if (!isBuilder
          && msgGo.equals("JSCompilerStateProto")
          && m.equals("getChunks")
          && arity == 1
          && shouldCacheIndexedList(
              rb,
              te,
              "getChunks",
              "getChunksList",
              "com.google.javascript.jscomp.ChunkProto")) {
        return "\tl := this_.getOrBuildChunksList_()\n"
            + "\tCheckElementIndex(arg0, int32(len(l.items)))\n"
            + "\treturn l.items[arg0].(*Ext_com_google_javascript_jscomp_ChunkProto)\n";
      }
      if ((m.equals("getChild") || m.equals("getType") || m.equals("getChunks")) && arity == 1) {
        String rStruct = retType.startsWith("*") ? retType.substring(1) : retType;
        if (!isBuilder
            && qn.equals("com.google.javascript.jscomp.serialization.AstNode")
            && m.equals("getChild")) {
          String pinAst =
              isExtTypeReferenced(rb, qn + ".Builder")
                  ? "\tif b, ok := this_.Core_.(*" + bldExtStruct + "); ok { b.pinned_ = true }\n"
                  : "";
          return "\tif this_ == nil { ThrowNPE() }\n"
              + "\tif r, ok := this_.Core_.(*protobuf.AstNode); ok {\n"
              + "\t\tif cIdx, cache := r.SlabChildFast(int(arg0)); cache != nil {\n"
              + "\t\t\treturn &cache.(*astNodeExtSlabCache_).wraps[cIdx]\n"
              + "\t\t}\n"
              + "\t\treturn getAstNodeChildSlow_(r, int(arg0))\n"
              + "\t}\n"
              + "\tv := " + pExpr + ".GetChild(int(arg0))\n"
              + pinAst
              + "\tif v == nil { return nil }\n\tout := &" + rStruct + "{}\n\tout.Core_ = v\n\treturn out\n";
        }
        String pinAst =
            (!isBuilder
                    && qn.equals("com.google.javascript.jscomp.serialization.AstNode")
                    && isExtTypeReferenced(rb, qn + ".Builder"))
                ? "\tif b, ok := this_.Core_.(*" + bldExtStruct + "); ok { b.pinned_ = true }\n"
                : "";
        return "\tv := " + pExpr + "." + GoNames.exported(m) + "(int(arg0))\n"
            + pinAst
            + "\tif v == nil { return nil }\n\tout := &" + rStruct + "{}\n\tout.Core_ = v\n\treturn out\n";
      }

      // Shared indexed/list getters on TypePool & JSCompilerStateProto
      if (!isBuilder
          && msgGo.equals("TypePool")
          && m.equals("getTypeList")
          && arity == 0
          && shouldCacheIndexedList(
              rb,
              te,
              "getType",
              "getTypeList",
              "com.google.javascript.jscomp.serialization.TypeProto")) {
        return "\treturn this_.getOrBuildTypeList_()\n";
      }
      if (!isBuilder
          && msgGo.equals("JSCompilerStateProto")
          && m.equals("getChunksList")
          && arity == 0
          && shouldCacheIndexedList(
              rb,
              te,
              "getChunks",
              "getChunksList",
              "com.google.javascript.jscomp.ChunkProto")) {
        return "\treturn this_.getOrBuildChunksList_()\n";
      }

      // Repeated List getters (`get...List()`)
      if (m.startsWith("get") && m.endsWith("List") && arity == 0) {
        return buildProtoListGetterBody(msgGo, pExpr, GoNames.exported(m), !isBuilder);
      }

      // `addAll...` on builders (note: guard against field names like `addAllowableFeatures`!)
      if (isBuilder && m.startsWith("addAll") && arity == 1) {
        String addAllBody = buildProtoAddAllBody(msgGo, pMutExpr, GoNames.exported(m));
        if (addAllBody != null) {
          return addAllBody;
        }
      }

      // Zero-arg Clear<Field> on builders
      if (isBuilder && m.startsWith("clear") && m.length() > 5 && arity == 0) {
        return "\t" + pMutExpr + "." + GoNames.exported(m) + "()\n\treturn this_\n";
      }

      // Zero-arg scalar/enum/sub-message/sub-builder getters (`get...`, `has...`, `add...Builder`)
      if (arity == 0) {
        String goMethod = GoNames.exported(m);
        if (retType.equals("bool")
            || retType.equals("int32")
            || retType.equals("int64")
            || retType.equals("float64")
            || retType.equals("string")) {
          return "\treturn " + pExpr + "." + goMethod + "()\n";
        }
        if (retType.startsWith("*Ext_")) {
          String retStruct = retType.substring(1);
          TypeElement retTe = elementOfReturn(sig);
          if (retTe != null) {
            String retQn = retTe.getQualifiedName().toString();
            String retEnumGo = PROTO_ENUM_TYPES.get(retQn);
            if (retEnumGo != null && isExtTypeReferenced(rb, retQn)) {
              return "\tv := "
                  + pExpr
                  + "."
                  + goMethod
                  + "()\n\treturn lookupProtoEnum_"
                  + retEnumGo
                  + "(int32(v))\n";
            }
            if (PROTO_ENUM_TYPES.containsKey(retQn)) {
              return "\tv := "
                  + pExpr
                  + "."
                  + goMethod
                  + "()\n\treturn lookupExtEnumByOrdinal["
                  + retStruct
                  + "]("
                  + GoNames.goStringLiteral(retQn)
                  + ", int32(v))\n";
            }
            if (!isBuilder && isCacheableSubMsgGetter(msgGo, sig) && isExtTypeReferenced(rb, retQn)) {
              return "\tif this_ == nil { ThrowNPE() }\n"
                  + "\tif p := this_.cache_" + goMethod + ".Load(); p != nil { return p }\n"
                  + "\tv := " + pExpr + "." + goMethod + "()\n"
                  + "\tif v == nil { return nil }\n"
                  + "\tif def := protoDefSub_" + msgGo + "_" + goMethod + "; def != nil && def.Core_ == v { return def }\n"
                  + "\tout := &" + retStruct + "{}\n"
                  + "\tout.Core_ = v\n"
                  + "\tif this_.cache_" + goMethod + ".CompareAndSwap(nil, out) { return out }\n"
                  + "\treturn this_.cache_" + goMethod + ".Load()\n";
            }
            if (!isBuilder && isDefaultSingletonSubMsgGetter(msgGo, sig) && isExtTypeReferenced(rb, retQn)) {
              return "\tif this_ == nil { ThrowNPE() }\n"
                  + "\tv := " + pExpr + "." + goMethod + "()\n"
                  + "\tif v == nil { return nil }\n"
                  + "\tif def := protoDefSub_" + msgGo + "_" + goMethod + "; def != nil && def.Core_ == v { return def }\n"
                  + "\tout := &" + retStruct + "{}\n"
                  + "\tout.Core_ = v\n"
                  + "\treturn out\n";
            }
            if (baseMsgQn(retQn) != null) {
              if (isBuilder && retQn.endsWith(".Builder")) {
                return "\tv := " + pMutExpr + "." + goMethod + "()\n"
                    + "\tif v == nil { return nil }\n"
                    + "\tout := &" + retStruct + "{Raw_: v}\n"
                    + "\tout.Core_ = v\n"
                    + "\treturn out\n";
              }
              return "\tv := " + pExpr + "." + goMethod + "()\n\tif v == nil { return nil }\n\tout := &" + retStruct + "{}\n\tout.Core_ = v\n\treturn out\n";
            }
          }
        }
      }

      // 1-arg setters & adders on builders (`set...`, `add...`)
      if (isBuilder && arity == 1 && (m.startsWith("set") || m.startsWith("add"))) {
        String goMethod = GoNames.exported(m);
        TypeMirror p0Tm = sig.canonParamTypes.get(0);
        String p0Go = rb.goTypeInJ2gort(p0Tm, imports);
        if (p0Go.equals("bool")
            || p0Go.equals("int32")
            || p0Go.equals("int64")
            || p0Go.equals("float64")
            || p0Go.equals("string")) {
          if (p0Go.equals("string")) {
            return "\tCheckNotNullString(arg0)\n\t" + pMutExpr + "." + goMethod + "(arg0)\n\treturn this_\n";
          }
          return "\t" + pMutExpr + "." + goMethod + "(arg0)\n\treturn this_\n";
        }
        if (p0Go.equals("Ext_com_google_protobuf_ByteString")) {
          return "\traw := unwrapByteString(arg0)\n\t" + pMutExpr + "." + goMethod + "(raw)\n\treturn this_\n";
        }
        TypeElement p0Te = elementOfType(p0Tm);
        if (p0Te != null) {
          String p0Qn = p0Te.getQualifiedName().toString();
          String p0EnumGo = PROTO_ENUM_TYPES.get(p0Qn);
          if (p0EnumGo != null) {
            return "\tCheckNotNull(arg0)\n\t"
                + pMutExpr
                + "."
                + goMethod
                + "(protobuf."
                + p0EnumGo
                + "(arg0.Ordinal()))\n\treturn this_\n";
          }
          String p0BaseMsg = baseMsgQn(p0Qn);
          if (p0BaseMsg != null) {
            String subGo = PROTO_MSG_TYPES.get(p0BaseMsg);
            if (qn.equals("com.google.javascript.jscomp.serialization.AstNode.Builder")
                && m.equals("addChild")) {
              if (p0Qn.equals("com.google.javascript.jscomp.serialization.AstNode")) {
                return "\tCheckNotNullPtr(arg0)\n"
                    + "\tif cb, ok := arg0.Core_.(*" + bldExtStruct + "); ok && cb.fromPool_ {\n"
                    + "\t\tif cb.linked_ || cb == this_ || !this_.fromPool_ || this_.pinned_ || this_.Raw_ != &this_.inlineRaw {\n"
                    + "\t\t\tcb.pinned_ = true\n"
                    + "\t\t} else {\n"
                    + "\t\t\tcb.linked_ = true\n"
                    + "\t\t\tcb.poolNextSibling_ = this_.poolChildHead_\n"
                    + "\t\t\tthis_.poolChildHead_ = cb\n"
                    + "\t\t}\n"
                    + "\t}\n"
                    + "\t" + pMutExpr + ".AddChild(arg0.buildRaw_AstNode())\n\treturn this_\n";
              }
              if (p0Qn.equals("com.google.javascript.jscomp.serialization.AstNode.Builder")) {
                return "\tCheckNotNullPtr(arg0)\n"
                    + "\trawChild := arg0.buildRaw_AstNode()\n"
                    + "\tif arg0.fromPool_ {\n"
                    + "\t\tif arg0.linked_ || arg0 == this_ || !this_.fromPool_ || this_.pinned_ || this_.Raw_ != &this_.inlineRaw || rawChild != &arg0.inlineRaw {\n"
                    + "\t\t\targ0.pinned_ = true\n"
                    + "\t\t} else {\n"
                    + "\t\t\targ0.linked_ = true\n"
                    + "\t\t\targ0.poolNextSibling_ = this_.poolChildHead_\n"
                    + "\t\t\tthis_.poolChildHead_ = arg0\n"
                    + "\t\t}\n"
                    + "\t}\n"
                    + "\t" + pMutExpr + ".AddChild(rawChild)\n\treturn this_\n";
              }
            }
            return "\tCheckNotNullPtr(arg0)\n\t"
                + pMutExpr
                + "."
                + goMethod
                + "(arg0.buildRaw_"
                + subGo
                + "())\n\treturn this_\n";
          }
        }
      }
    }

    return null;
  }

  private static String buildProtoListGetterBody(
      String msgGo, String pExpr, String goMethod, boolean cached) {
    // 1. Repeated string fields -> ProtocolStringList
    if (isStringListGetter(goMethod)) {
      if (!cached) {
        return "\treturn newProtoStringList(" + pExpr + "." + goMethod + "())\n";
      }
      return "\tCheckNotNullPtr(this_)\n"
          + "\tif p := this_.cache_" + goMethod + ".Load(); p != nil { return p }\n"
          + "\traw := " + pExpr + "." + goMethod + "()\n"
          + "\tif len(raw) == 0 {\n"
          + "\t\tthis_.cache_" + goMethod + ".CompareAndSwap(nil, emptyProtoStringList_)\n"
          + "\t\treturn emptyProtoStringList_\n"
          + "\t}\n"
          + "\tl := newProtoStringList(raw)\n"
          + "\tif this_.cache_" + goMethod + ".CompareAndSwap(nil, l) { return l }\n"
          + "\treturn this_.cache_" + goMethod + ".Load()\n";
    }

    // 2. Repeated int32 fields -> List<Integer>
    if (isInt32ListGetter(goMethod)) {
      if (!cached) {
        return "\treturn wrapProtoInt32List(" + pExpr + "." + goMethod + "())\n";
      }
      return "\tCheckNotNullPtr(this_)\n"
          + "\tif p := this_.cache_" + goMethod + ".Load(); p != nil { return p }\n"
          + "\traw := " + pExpr + "." + goMethod + "()\n"
          + "\tif len(raw) == 0 {\n"
          + "\t\tthis_.cache_" + goMethod + ".CompareAndSwap(nil, emptyProtoList_)\n"
          + "\t\treturn emptyProtoList_\n"
          + "\t}\n"
          + "\tl := wrapProtoInt32List(raw)\n"
          + "\tif this_.cache_" + goMethod + ".CompareAndSwap(nil, l) { return l }\n"
          + "\treturn this_.cache_" + goMethod + ".Load()\n";
    }

    // 3. Repeated enum fields
    if (msgGo.equals("JSCompilerStateProto") && goMethod.equals("GetAllowableFeaturesList")) {
      String buildItems =
          "\titems := make([]any, len(raw))\n"
              + "\tfor i, e := range raw { items[i] = lookupExtEnumByOrdinal[Ext_com_google_javascript_jscomp_FeatureProto](\"com.google.javascript.jscomp.FeatureProto\", int32(e)) }\n";
      if (!cached) {
        return "\traw := " + pExpr + ".GetAllowableFeaturesList()\n" + buildItems + "\treturn newProtoList(items)\n";
      }
      return "\tCheckNotNullPtr(this_)\n"
          + "\tif p := this_.cache_" + goMethod + ".Load(); p != nil { return p }\n"
          + "\traw := " + pExpr + ".GetAllowableFeaturesList()\n"
          + "\tif len(raw) == 0 {\n"
          + "\t\tthis_.cache_" + goMethod + ".CompareAndSwap(nil, emptyProtoList_)\n"
          + "\t\treturn emptyProtoList_\n"
          + "\t}\n"
          + buildItems
          + "\tl := newProtoList(items)\n"
          + "\tif this_.cache_" + goMethod + ".CompareAndSwap(nil, l) { return l }\n"
          + "\treturn this_.cache_" + goMethod + ".Load()\n";
    }
    if (msgGo.equals("OptimizationJsdoc") && goMethod.equals("GetKindList")) {
      String buildItems =
          "\titems := make([]any, len(raw))\n"
              + "\tfor i, e := range raw { items[i] = lookupExtEnumByOrdinal[Ext_com_google_javascript_jscomp_serialization_JsdocTag](\"com.google.javascript.jscomp.serialization.JsdocTag\", int32(e)) }\n";
      if (!cached) {
        return "\traw := " + pExpr + ".GetKindList()\n" + buildItems + "\treturn newProtoList(items)\n";
      }
      return "\tCheckNotNullPtr(this_)\n"
          + "\tif p := this_.cache_" + goMethod + ".Load(); p != nil { return p }\n"
          + "\traw := " + pExpr + ".GetKindList()\n"
          + "\tif len(raw) == 0 {\n"
          + "\t\tthis_.cache_" + goMethod + ".CompareAndSwap(nil, emptyProtoList_)\n"
          + "\t\treturn emptyProtoList_\n"
          + "\t}\n"
          + buildItems
          + "\tl := newProtoList(items)\n"
          + "\tif this_.cache_" + goMethod + ".CompareAndSwap(nil, l) { return l }\n"
          + "\treturn this_.cache_" + goMethod + ".Load()\n";
    }

    // 4. Repeated sub-message fields (single contiguous value slice instead of N heap pointers)
    String elemExt = subMsgListElemExt(goMethod);
    if (elemExt != null) {
      String buildItems =
          "\twrappers := make([]" + elemExt + ", len(raw))\n"
              + "\titems := make([]any, len(raw))\n"
              + "\tfor i, x := range raw { wrappers[i].Core_ = x; items[i] = &wrappers[i] }\n";
      if (!cached) {
        return "\traw := " + pExpr + "." + goMethod + "()\n\tif len(raw) == 0 { return emptyProtoList_ }\n" + buildItems + "\treturn newProtoList(items)\n";
      }
      return "\tif this_ == nil { ThrowNPE() }\n"
          + "\tif p := this_.cache_" + goMethod + ".Load(); p != nil { return p }\n"
          + "\traw := " + pExpr + "." + goMethod + "()\n"
          + "\tif len(raw) == 0 {\n"
          + "\t\tthis_.cache_" + goMethod + ".CompareAndSwap(nil, emptyProtoList_)\n"
          + "\t\treturn emptyProtoList_\n"
          + "\t}\n"
          + buildItems
          + "\tl := newProtoList(items)\n"
          + "\tif this_.cache_" + goMethod + ".CompareAndSwap(nil, l) { return l }\n"
          + "\treturn this_.cache_" + goMethod + ".Load()\n";
    }
    return null;
  }

  private static String buildProtoAddAllBody(String msgGo, String pExpr, String goMethod) {
    if (goMethod.equals("AddAllExportedNames")
        || goMethod.equals("AddAllCssNames")
        || goMethod.equals("AddAllInjectedLibraries")
        || goMethod.equals("AddAllWhitelist")
        || goMethod.equals("AddAllWhitelistRegexp")
        || goMethod.equals("AddAllAllowlist")
        || goMethod.equals("AddAllAllowlistRegexp")
        || goMethod.equals("AddAllOnlyApplyTo")
        || goMethod.equals("AddAllOnlyApplyToRegexp")
        || goMethod.equals("AddAllValue")
        || goMethod.equals("AddAllConfigFile")
        || goMethod.equals("AddAllRuntimeLibraryToInject")) {
      return "\t" + pExpr + "." + goMethod + "(unwrapProtoStringIterable(arg0))\n\treturn this_\n";
    }
    if (goMethod.equals("AddAllInstanceType")
        || goMethod.equals("AddAllPrototype")
        || goMethod.equals("AddAllOwnProperty")
        || goMethod.equals("AddAllInvolvedColor")
        || goMethod.equals("AddAllUnionMember")) {
      return "\t" + pExpr + "." + goMethod + "(unwrapProtoInt32Iterable(arg0))\n\treturn this_\n";
    }
    String subGo =
        switch (goMethod) {
          case "AddAllUniqueIdSupplier" -> "UniqueIdProto";
          case "AddAllStringMap", "AddAllInstrumentationMapping" -> "VariableMapEntryProto";
          case "AddAllWhitelistEntry", "AddAllAllowlistEntry" -> "RequirementScopeEntry";
          case "AddAllInstrumentationPointsStats" -> "InstrumentationPointStats";
          case "AddAllFileProfile" -> "FileProfile";
          case "AddAllSourceFile" -> "SourceFileProto";
          default -> null;
        };
    if (subGo != null) {
      return "\t" + pExpr + "." + goMethod + "(unwrapProtoMsgIterable[*protobuf." + subGo + "](arg0))\n\treturn this_\n";
    }
    return null;
  }

  private static String defaultRawExpr(String msgGo) {
    return switch (msgGo) {
      case "Mapping_OriginalMapping", "OptimizationJsdoc", "TypePool" ->
          "protobuf." + msgGo + "_GetDefaultInstance()";
      case "TemplateStringValue" -> "(&protobuf.AstNode{}).GetTemplateStringValue()";
      case "StringPoolProto" -> "(&protobuf.TypedAst{}).GetStringPool()";
      case "SourceFilePool" -> "(&protobuf.TypedAst{}).GetSourceFilePool()";
      case "ExternsSummary" -> "(&protobuf.TypedAst{}).GetExternsSummary()";
      case "TypePool_DebugInfo" -> "(&protobuf.TypePool{}).GetDebugInfo()";
      case "ObjectTypeProto" -> "(&protobuf.TypeProto{}).GetObject()";
      case "UnionTypeProto" -> "(&protobuf.TypeProto{}).GetUnion()";
      case "SourceFileProto_FileOnDisk" -> "(&protobuf.SourceFileProto{}).GetFileOnDisk()";
      case "SourceFileProto_ZipEntryOnDisk" -> "(&protobuf.SourceFileProto{}).GetZipEntry()";
      case "AccessorSummaryProto" -> "(&protobuf.JSCompilerStateProto{}).GetAccessorSummary()";
      default -> "&protobuf." + msgGo + "{}";
    };
  }

  private static TypeElement elementOfReturn(TypeResolver.MethodSig sig) {
    return sig.canonReturnType != null ? elementOfType(sig.canonReturnType) : null;
  }

  private static TypeElement elementOfType(TypeMirror tm) {
    if (tm instanceof javax.lang.model.type.DeclaredType dt
        && dt.asElement() instanceof TypeElement te) {
      return te;
    }
    return null;
  }

  public static String buildRuntimeProtobufGo(RuntimeBindings rb) {
    StringBuilder sb = new StringBuilder(24576);
    sb.append(
        """
        // Clean-room j2gort bridge to transpiled/shims/protobuf.
        package j2gort

        import (
        \t"strings"
        \t"sync"
        \t"sync/atomic"

        \t"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
        \t"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
        \t"github.com/paullepoulpe/jsgo/transpiled/shims/protobuf"
        )

        var _ = (*sync.Pool)(nil)
        var _ = (*atomic.Bool)(nil)

        type protoFieldDesc struct {
        \tname string
        \tidx  int32
        }

        func int8SliceToBytes(s []int8) []byte {
        \tif s == nil {
        \t\treturn nil
        \t}
        \tout := make([]byte, len(s))
        \tfor i, v := range s {
        \t\tout[i] = byte(v)
        \t}
        \treturn out
        }

        func bytesToInt8Slice(b []byte) []int8 {
        \tif b == nil {
        \t\treturn nil
        \t}
        \tout := make([]int8, len(b))
        \tfor i, v := range b {
        \t\tout[i] = int8(v)
        \t}
        \treturn out
        }

        func wrapByteString(bs protobuf.ByteString) Ext_com_google_protobuf_ByteString {
        \treturn Ext_com_google_protobuf_ByteString{raw: bs, valid: true}
        }

        func (b Ext_com_google_protobuf_ByteString) bs() protobuf.ByteString {
        \tif !b.valid {
        \t\tThrowNPE()
        \t}
        \treturn b.raw
        }

        func unwrapByteString(bs Ext_com_google_protobuf_ByteString) protobuf.ByteString {
        \tif !bs.valid {
        \t\tThrowNPE()
        \t}
        \treturn bs.raw
        }

        func BoxByteString(bs Ext_com_google_protobuf_ByteString) any {
        \tif !bs.valid {
        \t\treturn nil
        \t}
        \treturn bs
        }

        func CheckNotNullByteString(bs Ext_com_google_protobuf_ByteString) Ext_com_google_protobuf_ByteString {
        \tif !bs.valid {
        \t\tThrowNPE()
        \t}
        \treturn bs
        }

        func ByteStringRefEquals(a, b Ext_com_google_protobuf_ByteString) bool {
        \tif !a.valid || !b.valid {
        \t\treturn !a.valid && !b.valid
        \t}
        \tla := a.raw.Size()
        \tlb := b.raw.Size()
        \tif la != lb {
        \t\treturn false
        \t}
        \tif la == 0 {
        \t\treturn true
        \t}
        \treturn &a.raw.BytesView()[0] == &b.raw.BytesView()[0]
        }

        func ByteStringObjEquals(a, b Ext_com_google_protobuf_ByteString) bool {
        \tif !a.valid || !b.valid {
        \t\treturn !a.valid && !b.valid
        \t}
        \treturn a.raw.Equals(b.raw)
        }

        func ByteStringObjHash(a Ext_com_google_protobuf_ByteString) int32 {
        \tif !a.valid {
        \t\treturn 0
        \t}
        \treturn a.raw.HashCode()
        }

        func byteStringEquals(a protobuf.ByteString, other any) bool {
        \tif other == nil {
        \t\treturn false
        \t}
        \tif obs, ok := other.(Ext_com_google_protobuf_ByteString); ok {
        \t\treturn obs.valid && a.Equals(obs.raw)
        \t}
        \treturn false
        }

        func newByteStringIter(bs protobuf.ByteString) Ext_java_util_Iterator {
        \traw := bs.BytesView()
        \titems := make([]any, len(raw))
        \tfor i, v := range raw {
        \t\titems[i] = ValueOfByte(int8(v))
        \t}
        \treturn newSliceObjIter(items)
        }

        func byteStringForEach(bs protobuf.ByteString, consumer Ext_java_util_function_Consumer) {
        \tCheckNotNull(consumer)
        \tfor _, v := range bs.BytesView() {
        \t\tconsumer.Accept(ValueOfByte(int8(v)))
        \t}
        }

        func wrapCodedInputStream(c *protobuf.CodedInputStream) Ext_com_google_protobuf_CodedInputStream {
        \tif c == nil {
        \t\treturn nil
        \t}
        \tw := &Ext_com_google_protobuf_CodedInputStreamBase{Core_: c}
        \tw.Self_ = w
        \treturn w
        }

        func (c *Ext_com_google_protobuf_CodedInputStreamBase) cis() *protobuf.CodedInputStream {
        \tCheckNotNullPtr(c)
        \treturn c.Core_.(*protobuf.CodedInputStream)
        }

        func unwrapCodedInputStream(c Ext_com_google_protobuf_CodedInputStream) *protobuf.CodedInputStream {
        \tCheckNotNull(c)
        \treturn c.GetExt_com_google_protobuf_CodedInputStreamBase().cis()
        }

        func wrapInvalidProtoEx(msg string) Ext_com_google_protobuf_InvalidProtocolBufferException {
        \tex := &Ext_com_google_protobuf_InvalidProtocolBufferExceptionBase{}
        \tex.Self_ = ex
        \tInitUserThrowable(ex, "com.google.protobuf.InvalidProtocolBufferException", msg)
        \treturn ex
        }

        func checkProtoErr(err error) {
        \tif err == nil {
        \t\treturn
        \t}
        \tif ipbe, ok := err.(*protobuf.InvalidProtocolBufferException); ok {
        \t\tpanic(wrapInvalidProtoEx(ipbe.GetMessage()))
        \t}
        \tif th, ok := err.(javalang.Throwable); ok {
        \t\tpanic(th)
        \t}
        \tpanic(javaio.NewIOException(err.Error()))
        }

        func runProtoOp[R any](fn func() R) (res R) {
        \tdefer func() {
        \t\tif rec := recover(); rec != nil {
        \t\t\tif ipbe, ok := rec.(*protobuf.InvalidProtocolBufferException); ok {
        \t\t\t\tpanic(wrapInvalidProtoEx(ipbe.GetMessage()))
        \t\t\t}
        \t\t\tif pe, ok := rec.(*protobuf.TextFormat_ParseException); ok {
        \t\t\t\tpanic(javaio.NewIOException(pe.Error()))
        \t\t\t}
        \t\t\tpanic(rec)
        \t\t}
        \t}()
        \treturn fn()
        }

        func lookupExtEnumByOrdinal[T any](className string, ord int32) *T {
        \tvals := LookupEnumConstantsByName(className)
        \tif int(ord) >= 0 && int(ord) < len(vals) {
        \t\treturn vals[ord].(*T)
        \t}
        \tif len(vals) > 0 {
        \t\treturn vals[len(vals)-1].(*T)
        \t}
        \treturn nil
        }

        // Immutable random-access List wrapper for protobuf repeated fields.
        type jProtoList struct {
        \tExt_java_util_AbstractListBase
        \titems []any
        }

        func (l *jProtoList) Size() int32       { return int32(len(l.items)) }
        func (l *jProtoList) IsEmpty() bool     { return len(l.items) == 0 }
        func (l *jProtoList) Get(i int32) any   { CheckElementIndex(i, int32(len(l.items))); return l.items[i] }
        func (l *jProtoList) Iterator() any     { return newSliceObjIter(l.items) }
        func (l *jProtoList) IterSlice_() []any { return l.items }
        func (l *jProtoList) ToArray_0() []any {
        \tif len(l.items) == 0 {
        \t\treturn []any{}
        \t}
        \treturn CloneSlice(l.items)
        }

        type jProtoStringList struct {
        \tjProtoList
        }

        func (*jProtoStringList) Implements_j2gort_Ext_com_google_protobuf_ProtocolStringList() {}

        type jProtoByteStringList struct {
        \tExt_java_util_AbstractListBase
        \traw []protobuf.ByteString
        }

        func (l *jProtoByteStringList) Size() int32   { return int32(len(l.raw)) }
        func (l *jProtoByteStringList) IsEmpty() bool { return len(l.raw) == 0 }
        func (l *jProtoByteStringList) Get(i int32) any {
        \tCheckElementIndex(i, int32(len(l.raw)))
        \treturn Ext_com_google_protobuf_ByteString{raw: l.raw[i], valid: true}
        }
        func (l *jProtoByteStringList) toAnySlice() []any {
        \tif len(l.raw) == 0 {
        \t\treturn []any{}
        \t}
        \tout := make([]any, len(l.raw))
        \tfor i, bs := range l.raw {
        \t\tout[i] = Ext_com_google_protobuf_ByteString{raw: bs, valid: true}
        \t}
        \treturn out
        }
        func (l *jProtoByteStringList) Iterator() any     { return newSliceObjIter(l.toAnySlice()) }
        func (l *jProtoByteStringList) IterSlice_() []any { return l.toAnySlice() }
        func (l *jProtoByteStringList) ToArray_0() []any  { return l.toAnySlice() }

        func ListGetByteString(l any, i int32) Ext_com_google_protobuf_ByteString {
        \tif l == nil {
        \t\tThrowNPE()
        \t}
        \tif bl, ok := l.(*jProtoByteStringList); ok {
        \t\tif bl == nil {
        \t\t\tThrowNPE()
        \t\t}
        \t\tCheckElementIndex(i, int32(len(bl.raw)))
        \t\treturn Ext_com_google_protobuf_ByteString{raw: bl.raw[i], valid: true}
        \t}
        \treturn ListGetVal[Ext_com_google_protobuf_ByteString](l, i)
        }

        var emptyProtoList_ *jProtoList
        var emptyProtoStringList_ *jProtoStringList
        var emptyProtoByteStringList_ *jProtoByteStringList

        func init() {
        \tel := &jProtoList{}
        \tel.BindExtSelf_(el)
        \temptyProtoList_ = el
        \tes := &jProtoStringList{}
        \tes.BindExtSelf_(es)
        \temptyProtoStringList_ = es
        \teb := &jProtoByteStringList{}
        \teb.BindExtSelf_(eb)
        \temptyProtoByteStringList_ = eb
        }

        func newProtoByteStringList(raw []protobuf.ByteString) *jProtoByteStringList {
        \tif len(raw) == 0 {
        \t\treturn emptyProtoByteStringList_
        \t}
        \tl := &jProtoByteStringList{raw: raw}
        \tl.BindExtSelf_(l)
        \treturn l
        }

        func newProtoList(items []any) *jProtoList {
        \tif len(items) == 0 {
        \t\treturn emptyProtoList_
        \t}
        \tl := &jProtoList{items: items}
        \tl.BindExtSelf_(l)
        \treturn l
        }

        func newProtoStringList(strs []string) *jProtoStringList {
        \tif len(strs) == 0 {
        \t\treturn emptyProtoStringList_
        \t}
        \titems := make([]any, len(strs))
        \tfor i, s := range strs {
        \t\titems[i] = s
        \t}
        \tl := &jProtoStringList{}
        \tl.items = items
        \tl.BindExtSelf_(l)
        \treturn l
        }

        func wrapProtoStringList(strs []string) Ext_com_google_protobuf_ProtocolStringList {
        \treturn newProtoStringList(strs)
        }

        func wrapProtoInt32List(nums []int32) *jProtoList {
        \tif len(nums) == 0 {
        \t\treturn emptyProtoList_
        \t}
        \titems := make([]any, len(nums))
        \tfor i, n := range nums {
        \t\titems[i] = ValueOfInt(n)
        \t}
        \treturn newProtoList(items)
        }

        func unwrapProtoReadCore(v any) any {
        \tif pb, ok := v.(interface{ protoBuildCore_() any }); ok {
        \t\treturn pb.protoBuildCore_()
        \t}
        \tif gc, ok := v.(interface{ getCore() any }); ok {
        \t\treturn gc.getCore()
        \t}
        \treturn v
        }

        func unwrapProtoMutCore(v any) any {
        \tif pm, ok := v.(interface{ protoMutCore_() any }); ok {
        \t\treturn pm.protoMutCore_()
        \t}
        \treturn v.(interface{ getCore() any }).getCore()
        }

        func unwrapProtoStringIterable(it Ext_java_lang_Iterable) []string {
        \tCheckNotNull(it)
        \tvar out []string
        \tfor iter := NewIterCursor(it); iter.HasNext(); {
        \t\tv := CheckNotNull(iter.Next())
        \t\tout = append(out, v.(string))
        \t}
        \treturn out
        }

        func unwrapProtoInt32Iterable(it Ext_java_lang_Iterable) []int32 {
        \tCheckNotNull(it)
        \tvar out []int32
        \tfor iter := NewIterCursor(it); iter.HasNext(); {
        \t\tv := CheckNotNull(iter.Next())
        \t\tout = append(out, v.(*BoxInt).Val)
        \t}
        \treturn out
        }

        func unwrapProtoMsgIterable[T any](it Ext_java_lang_Iterable) []T {
        \tCheckNotNull(it)
        \tvar out []T
        \tfor iter := NewIterCursor(it); iter.HasNext(); {
        \t\tv := CheckNotNull(iter.Next())
        \t\tif pb, ok := v.(interface{ protoBuildCore_() any }); ok {
        \t\t\tout = append(out, pb.protoBuildCore_().(T))
        \t\t} else {
        \t\t\tc := v.(interface{ getCore() any }).getCore().(T)
        \t\t\tout = append(out, c)
        \t\t}
        \t}
        \treturn out
        }

        func protoMsgToByteString(core any) Ext_com_google_protobuf_ByteString {
        \tif ps, ok := core.(protobuf.ProtoSerializable); ok {
        \t\treturn wrapByteString(protobuf.MessageToByteString(ps))
        \t}
        \tif tba, ok := core.(interface{ ToByteArray() []byte }); ok {
        \t\treturn wrapByteString(protobuf.ByteString_FromBytesUnsafeWrap(tba.ToByteArray()))
        \t}
        \treturn wrapByteString(protobuf.ByteString_EMPTY)
        }

        func protoMsgWriteTo(core any, out Ext_java_io_OutputStream) {
        \tCheckNotNull(out)
        \tw := asOutputStream(out)
        \tif ps, ok := core.(protobuf.ProtoSerializable); ok {
        \t\tprotobuf.WriteMessageTo(ps, w)
        \t\treturn
        \t}
        \tif tba, ok := core.(interface{ ToByteArray() []byte }); ok {
        \t\tif _, err := w.Write(tba.ToByteArray()); err != nil {
        \t\t\tcheckProtoErr(err)
        \t\t}
        \t}
        }

        func protoMsgWriteDelimitedTo(core any, out Ext_java_io_OutputStream) {
        \tCheckNotNull(out)
        \tw := asOutputStream(out)
        \tif ps, ok := core.(protobuf.ProtoSerializable); ok {
        \t\tprotobuf.WriteMessageDelimitedTo(ps, w)
        \t\treturn
        \t}
        \tif tba, ok := core.(interface{ ToByteArray() []byte }); ok {
        \t\traw := tba.ToByteArray()
        \t\tvar hdr [10]byte
        \t\tn := 0
        \t\tv := uint64(len(raw))
        \t\tfor v >= 0x80 {
        \t\t\thdr[n] = byte(v) | 0x80
        \t\t\tv >>= 7
        \t\t\tn++
        \t\t}
        \t\thdr[n] = byte(v)
        \t\tn++
        \t\tif _, err := w.Write(hdr[:n]); err != nil {
        \t\t\tcheckProtoErr(err)
        \t\t}
        \t\tif _, err := w.Write(raw); err != nil {
        \t\t\tcheckProtoErr(err)
        \t\t}
        \t}
        }

        func protoMsgEquals(self any, other any) bool {
        \tif self == other {
        \t\treturn true
        \t}
        \tif other == nil || javalang.IsNil(other) {
        \t\treturn false
        \t}
        \tr1 := unwrapProtoReadCore(self)
        \tr2 := unwrapProtoReadCore(other)
        \tif r1 == self || r2 == other {
        \t\treturn false
        \t}
        \tif r1 == r2 {
        \t\treturn true
        \t}
        \tif eq, ok := r1.(interface{ Equals(any) bool }); ok {
        \t\treturn eq.Equals(r2)
        \t}
        \tps1, s1 := r1.(protobuf.ProtoSerializable)
        \tps2, s2 := r2.(protobuf.ProtoSerializable)
        \tif s1 && s2 {
        \t\treturn string(protobuf.MessageToByteArray(ps1)) == string(protobuf.MessageToByteArray(ps2))
        \t}
        \treturn false
        }

        func protoGetAllFields(core any) Ext_java_util_Map {
        \tmc := NewJMapCore(false)
        \tif req, ok := core.(*protobuf.Requirement); ok && req != nil {
        \t\tfor i, name := range req.PopulatedFieldNames() {
        \t\t\tfd := &Ext_com_google_protobuf_Descriptors_FieldDescriptor{}
        \t\t\tfd.Core_ = protoFieldDesc{name: name, idx: int32(i)}
        \t\t\tmc.Put(fd, Boolean_TRUE)
        \t\t}
        \t}
        \treturn wrapImmutableMap(mc)
        }

        func readTextFormatInput(in any) string {
        \tCheckNotNull(in)
        \tif s, ok := in.(string); ok {
        \t\treturn s
        \t}
        \tif rd, ok := in.(Ext_java_io_Reader); ok {
        \t\treturn readAllFromReadable(rd)
        \t}
        \tvar sb strings.Builder
        \t_ = sb
        \treturn StrVal(in)
        }

        func (x *protoFieldDesc) IsJavaNil_() bool       { return x == nil }
        func (x *jProtoList) IsJavaNil_() bool           { return x == nil }
        func (x *jProtoStringList) IsJavaNil_() bool     { return x == nil }
        func (x *jProtoByteStringList) IsJavaNil_() bool { return x == nil }
        """);

    Map<String, String> sortedProtoTypes = new TreeMap<>(PROTO_MSG_TYPES);
    Map<String, String> sortedEnumTypes = new TreeMap<>(PROTO_ENUM_TYPES);
    StringBuilder initBody = new StringBuilder(4096);

    // 0. Emit fast atomic-cached ordinal lookup helpers for referenced protobuf enums.
    for (Map.Entry<String, String> e : sortedEnumTypes.entrySet()) {
      String enumQn = e.getKey();
      String enumGo = e.getValue();
      TypeElement enumTe = rb.resolver().elements.getTypeElement(enumQn);
      if (enumTe != null && rb.resolver().getReferencedExternalTypes().contains(enumTe)) {
        String eStruct = "Ext_" + GoNames.cleanChars(enumQn.replace('.', '_'));
        sb.append("var protoEnumCache_").append(enumGo).append("_ atomic.Pointer[[]*").append(eStruct).append("]\n");
        sb.append("func lookupProtoEnum_").append(enumGo).append("(ord int32) *").append(eStruct).append(" {\n")
            .append("\ts := protoEnumCache_").append(enumGo).append("_.Load()\n")
            .append("\tif s == nil {\n")
            .append("\t\traw := LookupEnumConstantsByName(").append(GoNames.goStringLiteral(enumQn)).append(")\n")
            .append("\t\ttyped := make([]*").append(eStruct).append(", len(raw))\n")
            .append("\t\tfor i, v := range raw { typed[i] = v.(*").append(eStruct).append(") }\n")
            .append("\t\tprotoEnumCache_").append(enumGo).append("_.Store(&typed)\n")
            .append("\t\ts = &typed\n")
            .append("\t}\n")
            .append("\tvals := *s\n")
            .append("\tif int(ord) >= 0 && int(ord) < len(vals) { return vals[ord] }\n")
            .append("\tif len(vals) > 0 { return vals[len(vals)-1] }\n")
            .append("\treturn nil\n}\n");
      }
    }

    // 1. Emit default instance vars, sub-message default vars, and typed helpers on referenced proto types.
    for (Map.Entry<String, String> e : sortedProtoTypes.entrySet()) {
      String qn = e.getKey();
      String msgGo = e.getValue();
      TypeElement msgTe = rb.resolver().elements.getTypeElement(qn);
      boolean msgReferenced =
          msgTe != null && rb.resolver().getReferencedExternalTypes().contains(msgTe);
      if (msgReferenced) {
        String st = "Ext_" + GoNames.cleanChars(qn.replace('.', '_'));
        sb.append("var protoDefault_").append(msgGo).append("_ *").append(st).append('\n');
        sb.append("var protoDefaultRaw_").append(msgGo).append("_ *protobuf.").append(msgGo).append('\n');
        initBody
            .append("\tprotoDefaultRaw_")
            .append(msgGo)
            .append("_ = ")
            .append(defaultRawExpr(msgGo))
            .append("\n\tprotoDefault_")
            .append(msgGo)
            .append("_ = &")
            .append(st)
            .append("{}\n\tprotoDefault_")
            .append(msgGo)
            .append("_.Core_ = protoDefaultRaw_")
            .append(msgGo)
            .append("_\n");

        sb.append("func (m *")
            .append(st)
            .append(") p_")
            .append(msgGo)
            .append("() *protobuf.")
            .append(msgGo)
            .append(" {\n")
            .append("\tif m == nil { ThrowNPE() }\n")
            .append("\tif m.Core_ == nil { return protoDefaultRaw_")
            .append(msgGo)
            .append("_ }\n");
        if (qn.equals("com.google.javascript.jscomp.serialization.AstNode")
            && isExtTypeReferenced(rb, qn + ".Builder")) {
          String bst = "Ext_" + GoNames.cleanChars((qn + ".Builder").replace('.', '_'));
          sb.append("\tif r, ok := m.Core_.(*protobuf.AstNode); ok { return r }\n")
              .append("\treturn &m.Core_.(*")
              .append(bst)
              .append(").inlineRaw\n}\n");
        } else {
          sb.append("\treturn m.Core_.(*protobuf.")
              .append(msgGo)
              .append(")\n}\n");
        }
        sb.append("func (m *")
            .append(st)
            .append(") buildRaw_")
            .append(msgGo)
            .append("() *protobuf.")
            .append(msgGo)
            .append(" { return m.p_")
            .append(msgGo)
            .append("() }\n");
        sb.append("func (m *")
            .append(st)
            .append(") protoBuildCore_() any { return m.buildRaw_")
            .append(msgGo)
            .append("() }\n");
        if (qn.equals("com.google.javascript.jscomp.serialization.AstNode")) {
          sb.append("type astNodeExtSlabCache_ struct {\n")
              .append("\twraps       []")
              .append(st)
              .append("\n\tfirstHandle *protobuf.AstNode\n")
              .append("\tinitLen     int\n}\n")
              .append("func initAstNodeExtSlabCache_(prev any, root *protobuf.AstNode, handles []protobuf.AstNode) any {\n")
              .append("\tc, _ := prev.(*astNodeExtSlabCache_)\n")
              .append("\tif c == nil {\n\t\tc = &astNodeExtSlabCache_{}\n\t}\n")
              .append("\ttotal := 1 + len(handles)\n")
              .append("\tif cap(c.wraps) < total {\n")
              .append("\t\tc.wraps = make([]")
              .append(st)
              .append(", total, max(total, cap(c.wraps)*2))\n")
              .append("\t\tc.firstHandle = nil\n\t\tc.initLen = 0\n")
              .append("\t} else {\n\t\tc.wraps = c.wraps[:total]\n\t}\n")
              .append("\tif len(handles) > 0 && c.firstHandle != &handles[0] {\n")
              .append("\t\tc.firstHandle = &handles[0]\n\t\tc.initLen = 0\n\t}\n")
              .append("\tif c.initLen == 0 {\n\t\tc.wraps[0].Core_ = root\n\t\tc.initLen = 1\n\t}\n")
              .append("\tfor i := c.initLen; i < total; i++ {\n\t\tc.wraps[i].Core_ = &handles[i-1]\n\t}\n")
              .append("\tif c.initLen < total {\n\t\tc.initLen = total\n\t}\n")
              .append("\treturn c\n}\n")
              .append("//go:noinline\n")
              .append("func getAstNodeChildSlow_(r *protobuf.AstNode, idx int) *")
              .append(st)
              .append(" {\n")
              .append("\tif cIdx, cache, ok := r.SlabChildInit(idx, initAstNodeExtSlabCache_); ok {\n")
              .append("\t\treturn &cache.(*astNodeExtSlabCache_).wraps[cIdx]\n\t}\n")
              .append("\tv := r.GetChild(idx)\n")
              .append("\tif v == nil {\n\t\treturn nil\n\t}\n")
              .append("\tout := &")
              .append(st)
              .append("{}\n\tout.Core_ = v\n\treturn out\n}\n")
              .append("func (this_ *")
              .append(st)
              .append(") Release_() {\n")
              .append("\tif this_ == nil {\n\t\treturn\n\t}\n")
              .append("\tif r, ok := this_.Core_.(*protobuf.AstNode); ok && r != nil {\n")
              .append("\t\tthis_.Core_ = nil\n\t\tr.ReleaseSlab()\n\t}\n}\n");
        }

        if (msgGo.equals("TypePool")
            && shouldCacheIndexedList(
                rb,
                msgTe,
                "getType",
                "getTypeList",
                "com.google.javascript.jscomp.serialization.TypeProto")) {
          sb.append("func (this_ *")
              .append(st)
              .append(") getOrBuildTypeList_() *jProtoList {\n")
              .append("\tif this_ == nil { ThrowNPE() }\n")
              .append("\tif p := this_.cache_GetTypeList.Load(); p != nil { return p }\n")
              .append("\traw := this_.p_TypePool().GetTypeList()\n")
              .append("\tif len(raw) == 0 {\n")
              .append("\t\tthis_.cache_GetTypeList.CompareAndSwap(nil, emptyProtoList_)\n")
              .append("\t\treturn emptyProtoList_\n")
              .append("\t}\n")
              .append("\twrappers := make([]Ext_com_google_javascript_jscomp_serialization_TypeProto, len(raw))\n")
              .append("\titems := make([]any, len(raw))\n")
              .append("\tfor i, x := range raw {\n")
              .append("\t\twrappers[i].Core_ = x\n")
              .append("\t\titems[i] = &wrappers[i]\n")
              .append("\t}\n")
              .append("\tl := newProtoList(items)\n")
              .append("\tif this_.cache_GetTypeList.CompareAndSwap(nil, l) { return l }\n")
              .append("\treturn this_.cache_GetTypeList.Load()\n}\n");
        }
        if (msgGo.equals("JSCompilerStateProto")
            && shouldCacheIndexedList(
                rb,
                msgTe,
                "getChunks",
                "getChunksList",
                "com.google.javascript.jscomp.ChunkProto")) {
          sb.append("func (this_ *")
              .append(st)
              .append(") getOrBuildChunksList_() *jProtoList {\n")
              .append("\tif this_ == nil { ThrowNPE() }\n")
              .append("\tif p := this_.cache_GetChunksList.Load(); p != nil { return p }\n")
              .append("\traw := this_.p_JSCompilerStateProto().GetChunksList()\n")
              .append("\tif len(raw) == 0 {\n")
              .append("\t\tthis_.cache_GetChunksList.CompareAndSwap(nil, emptyProtoList_)\n")
              .append("\t\treturn emptyProtoList_\n")
              .append("\t}\n")
              .append("\twrappers := make([]Ext_com_google_javascript_jscomp_ChunkProto, len(raw))\n")
              .append("\titems := make([]any, len(raw))\n")
              .append("\tfor i, x := range raw {\n")
              .append("\t\twrappers[i].Core_ = x\n")
              .append("\t\titems[i] = &wrappers[i]\n")
              .append("\t}\n")
              .append("\tl := newProtoList(items)\n")
              .append("\tif this_.cache_GetChunksList.CompareAndSwap(nil, l) { return l }\n")
              .append("\treturn this_.cache_GetChunksList.Load()\n}\n");
        }
      }

      TypeElement bldTe = rb.resolver().elements.getTypeElement(qn + ".Builder");
      if (bldTe != null && rb.resolver().getReferencedExternalTypes().contains(bldTe)) {
        String bst = "Ext_" + GoNames.cleanChars((qn + ".Builder").replace('.', '_'));
        String mst = "Ext_" + GoNames.cleanChars(qn.replace('.', '_'));
        sb.append("func (b *")
            .append(bst)
            .append(") p_")
            .append(msgGo)
            .append("() *protobuf.")
            .append(msgGo)
            .append(" {\n")
            .append("\tif b == nil { ThrowNPE() }\n")
            .append("\tif b.Raw_ != nil { return b.Raw_ }\n")
            .append("\tif b.Core_ != nil {\n")
            .append("\t\tr := b.Core_.(*protobuf.")
            .append(msgGo)
            .append(")\n")
            .append("\t\tb.Raw_ = r\n")
            .append("\t\treturn r\n")
            .append("\t}\n")
            .append("\tb.Raw_ = &b.inlineRaw\n")
            .append("\tb.Core_ = b.Raw_\n")
            .append("\treturn b.Raw_\n}\n");

        sb.append("func (b *")
            .append(bst)
            .append(") pMut_")
            .append(msgGo)
            .append("() *protobuf.")
            .append(msgGo)
            .append(" {\n")
            .append("\tif b == nil { ThrowNPE() }\n")
            .append("\tif !b.built_ && b.Raw_ != nil { return b.Raw_ }\n")
            .append("\tif b.built_ {\n");
        if (qn.equals("com.google.javascript.jscomp.serialization.AstNode")) {
          sb.append("\t\tb.pinned_ = true\n");
          if (msgReferenced) {
            sb.append("\t\tif b.usedBuiltMsg_ && b.builtMsg.Core_ == b { b.builtMsg.Core_ = b.Raw_ }\n");
          }
        }
        if (hasToBuilderMethod(msgGo)) {
          sb.append("\t\tb.Raw_ = b.Raw_.ToBuilder()\n");
        } else {
          sb.append("\t\tcp := *b.Raw_\n\t\tb.Raw_ = &cp\n");
        }
        sb.append("\t\tb.Core_ = b.Raw_\n")
            .append("\t\tb.built_ = false\n")
            .append("\t\treturn b.Raw_\n")
            .append("\t}\n")
            .append("\treturn b.p_")
            .append(msgGo)
            .append("()\n}\n");

        sb.append("func (b *")
            .append(bst)
            .append(") buildRaw_")
            .append(msgGo)
            .append("() *protobuf.")
            .append(msgGo)
            .append(" {\n")
            .append("\tp := b.p_")
            .append(msgGo)
            .append("()\n")
            .append("\tout := p.Build()\n")
            .append("\tif out == p { b.built_ = true }\n")
            .append("\treturn out\n}\n");
        if (msgReferenced) {
          sb.append("func (b *")
              .append(bst)
              .append(") buildMsg_")
              .append(msgGo)
              .append("() *")
              .append(mst)
              .append(" {\n")
              .append("\traw := b.buildRaw_")
              .append(msgGo)
              .append("()\n");
          if (qn.equals("com.google.javascript.jscomp.serialization.AstNode")) {
            sb.append("\tif b.fromPool_ && !b.pinned_ && raw == &b.inlineRaw {\n")
                .append("\t\tif !b.usedBuiltMsg_ {\n")
                .append("\t\t\tb.usedBuiltMsg_ = true\n")
                .append("\t\t\tb.builtMsg.Core_ = b\n")
                .append("\t\t\treturn &b.builtMsg\n")
                .append("\t\t}\n")
                .append("\t\tif b.builtMsg.Core_ == b {\n")
                .append("\t\t\treturn &b.builtMsg\n")
                .append("\t\t}\n")
                .append("\t\tb.pinned_ = true\n")
                .append("\t}\n");
          }
          sb.append("\tif !b.usedBuiltMsg_ {\n")
              .append("\t\tb.usedBuiltMsg_ = true\n")
              .append("\t\tb.builtMsg.Core_ = raw\n")
              .append("\t\treturn &b.builtMsg\n")
              .append("\t}\n")
              .append("\tif b.builtMsg.Core_ == raw {\n")
              .append("\t\treturn &b.builtMsg\n")
              .append("\t}\n")
              .append("\tout := &")
              .append(mst)
              .append("{}\n")
              .append("\tout.Core_ = raw\n")
              .append("\treturn out\n}\n");
        }
        sb.append("func (b *")
            .append(bst)
            .append(") protoMutCore_() any { return b.pMut_")
            .append(msgGo)
            .append("() }\n");
        sb.append("func (b *")
            .append(bst)
            .append(") protoBuildCore_() any { return b.buildRaw_")
            .append(msgGo)
            .append("() }\n");
        if (qn.equals("com.google.javascript.jscomp.serialization.AstNode")) {
          sb.append("type astNodeBuilderBatch struct {\n")
              .append("\tidx  atomic.Uint32\n")
              .append("\tdone atomic.Uint32\n")
              .append("\tptrs [128]*")
              .append(bst)
              .append("\n}\n")
              .append("type astNodeBuilderChunk struct {\n")
              .append("\titems [128]")
              .append(bst)
              .append("\n}\n")
              .append("var (\n\tastNodeBuilderBatch_     atomic.Pointer[astNodeBuilderBatch]\n\tastNodeBuilderPoolMu_    sync.Mutex\n\tastNodeBuilderCurBatch_  *astNodeBuilderBatch\n\tastNodeBuilderBatchFree_ []*astNodeBuilderBatch\n\tastNodeBuilderRetired_   []*astNodeBuilderBatch\n\tastNodeBuilderFree_      []*")
              .append(bst)
              .append("\n)\n")
              .append("func allocAstNodeBuilder_() *")
              .append(bst)
              .append(" {\n")
              .append("\tif batch := astNodeBuilderBatch_.Load(); batch != nil {\n")
              .append("\t\tif i := batch.idx.Add(1) - 1; i < 128 {\n")
              .append("\t\t\tb := batch.ptrs[i]\n")
              .append("\t\t\tbatch.done.Add(1)\n")
              .append("\t\t\treturn b\n")
              .append("\t\t}\n")
              .append("\t}\n")
              .append("\treturn allocAstNodeBuilderSlow_()\n}\n")
              .append("func allocAstNodeBuilderSlow_() *")
              .append(bst)
              .append(" {\n")
              .append("\tastNodeBuilderPoolMu_.Lock()\n")
              .append("\tif batch := astNodeBuilderBatch_.Load(); batch != nil {\n")
              .append("\t\tif i := batch.idx.Add(1) - 1; i < 128 {\n")
              .append("\t\t\tb := batch.ptrs[i]\n")
              .append("\t\t\tbatch.done.Add(1)\n")
              .append("\t\t\tastNodeBuilderPoolMu_.Unlock()\n")
              .append("\t\t\treturn b\n")
              .append("\t\t}\n")
              .append("\t}\n")
              .append("\tif old := astNodeBuilderCurBatch_; old != nil {\n")
              .append("\t\tastNodeBuilderCurBatch_ = nil\n")
              .append("\t\tastNodeBuilderBatch_.Store(nil)\n")
              .append("\t\tif old.done.Load() == 128 {\n")
              .append("\t\t\tclear(old.ptrs[:])\n")
              .append("\t\t\tif len(astNodeBuilderBatchFree_) < 16 {\n")
              .append("\t\t\t\tastNodeBuilderBatchFree_ = append(astNodeBuilderBatchFree_, old)\n")
              .append("\t\t\t}\n")
              .append("\t\t} else {\n")
              .append("\t\t\tastNodeBuilderRetired_ = append(astNodeBuilderRetired_, old)\n")
              .append("\t\t}\n")
              .append("\t}\n")
              .append("\tif len(astNodeBuilderRetired_) > 0 {\n")
              .append("\t\tw := 0\n")
              .append("\t\tfor _, rb := range astNodeBuilderRetired_ {\n")
              .append("\t\t\tif rb.done.Load() == 128 {\n")
              .append("\t\t\t\tclear(rb.ptrs[:])\n")
              .append("\t\t\t\tif len(astNodeBuilderBatchFree_) < 16 {\n")
              .append("\t\t\t\t\tastNodeBuilderBatchFree_ = append(astNodeBuilderBatchFree_, rb)\n")
              .append("\t\t\t\t}\n")
              .append("\t\t\t} else {\n")
              .append("\t\t\t\tastNodeBuilderRetired_[w] = rb\n")
              .append("\t\t\t\tw++\n")
              .append("\t\t\t}\n")
              .append("\t\t}\n")
              .append("\t\tclear(astNodeBuilderRetired_[w:])\n")
              .append("\t\tastNodeBuilderRetired_ = astNodeBuilderRetired_[:w]\n")
              .append("\t}\n")
              .append("\tvar nb *astNodeBuilderBatch\n")
              .append("\tif nbf := len(astNodeBuilderBatchFree_); nbf > 0 {\n")
              .append("\t\tnb = astNodeBuilderBatchFree_[nbf-1]\n")
              .append("\t\tastNodeBuilderBatchFree_[nbf-1] = nil\n")
              .append("\t\tastNodeBuilderBatchFree_ = astNodeBuilderBatchFree_[:nbf-1]\n")
              .append("\t} else {\n")
              .append("\t\tnb = &astNodeBuilderBatch{}\n")
              .append("\t\tnb.idx.Store(128)\n")
              .append("\t}\n")
              .append("\tnFree := len(astNodeBuilderFree_)\n")
              .append("\tif nFree >= 128 {\n")
              .append("\t\tstart := nFree - 128\n")
              .append("\t\tcopy(nb.ptrs[:], astNodeBuilderFree_[start:])\n")
              .append("\t\tclear(astNodeBuilderFree_[start:])\n")
              .append("\t\tastNodeBuilderFree_ = astNodeBuilderFree_[:start]\n")
              .append("\t} else {\n")
              .append("\t\tcopied := 0\n")
              .append("\t\tif nFree > 0 {\n")
              .append("\t\t\tcopied = copy(nb.ptrs[:], astNodeBuilderFree_)\n")
              .append("\t\t\tclear(astNodeBuilderFree_)\n")
              .append("\t\t\tastNodeBuilderFree_ = astNodeBuilderFree_[:0]\n")
              .append("\t\t}\n")
              .append("\t\tnc := &astNodeBuilderChunk{}\n")
              .append("\t\tneed := 128 - copied\n")
              .append("\t\tfor i := 0; i < need; i++ {\n")
              .append("\t\t\tb := &nc.items[i]\n")
              .append("\t\t\tb.Raw_ = &b.inlineRaw\n")
              .append("\t\t\tb.Core_ = b.Raw_\n")
              .append("\t\t\tb.fromPool_ = true\n")
              .append("\t\t\tnb.ptrs[copied+i] = b\n")
              .append("\t\t}\n")
              .append("\t\tfor i := need; i < 128; i++ {\n")
              .append("\t\t\tb := &nc.items[i]\n")
              .append("\t\t\tb.Raw_ = &b.inlineRaw\n")
              .append("\t\t\tb.Core_ = b.Raw_\n")
              .append("\t\t\tb.fromPool_ = true\n")
              .append("\t\t\tastNodeBuilderFree_ = append(astNodeBuilderFree_, b)\n")
              .append("\t\t}\n")
              .append("\t}\n")
              .append("\tret := nb.ptrs[0]\n")
              .append("\tnb.done.Store(1)\n")
              .append("\tnb.idx.Store(1)\n")
              .append("\tastNodeBuilderCurBatch_ = nb\n")
              .append("\tastNodeBuilderBatch_.Store(nb)\n")
              .append("\tastNodeBuilderPoolMu_.Unlock()\n")
              .append("\treturn ret\n}\n")
              .append("func recycleAstNodeBuilderTree_(root *")
              .append(bst)
              .append(") {\n")
              .append("\tif root == nil || root.pinned_ {\n")
              .append("\t\treturn\n")
              .append("\t}\n")
              .append("\troot.poolNextSibling_ = nil\n")
              .append("\tworkHead := root\n")
              .append("\tvar freeHead *")
              .append(bst)
              .append("\n")
              .append("\tfor workHead != nil {\n")
              .append("\t\tcur := workHead\n")
              .append("\t\tworkHead = cur.poolNextSibling_\n")
              .append("\t\tcur.poolNextSibling_ = nil\n")
              .append("\t\tif cur.pinned_ {\n")
              .append("\t\t\tcontinue\n")
              .append("\t\t}\n")
              .append("\t\tchildren := cur.poolChildHead_\n")
              .append("\t\tcur.poolChildHead_ = nil\n")
              .append("\t\tfor c := children; c != nil; {\n")
              .append("\t\t\tnext := c.poolNextSibling_\n")
              .append("\t\t\tc.poolNextSibling_ = nil\n")
              .append("\t\t\tc.linked_ = false\n")
              .append("\t\t\tif !c.pinned_ {\n")
              .append("\t\t\t\tc.poolNextSibling_ = workHead\n")
              .append("\t\t\t\tworkHead = c\n")
              .append("\t\t\t}\n")
              .append("\t\t\tc = next\n")
              .append("\t\t}\n")
              .append("\t\tcur.inlineRaw.ResetForPool()\n")
              .append("\t\tcur.Raw_ = &cur.inlineRaw\n")
              .append("\t\tcur.Core_ = cur.Raw_\n")
              .append("\t\tcur.builtMsg.Core_ = nil\n")
              .append("\t\tcur.built_ = false\n")
              .append("\t\tcur.usedBuiltMsg_ = false\n")
              .append("\t\tcur.linked_ = false\n")
              .append("\t\tcur.pinned_ = false\n")
              .append("\t\tcur.poolChildHead_ = freeHead\n")
              .append("\t\tfreeHead = cur\n")
              .append("\t}\n")
              .append("\tif freeHead == nil {\n")
              .append("\t\treturn\n")
              .append("\t}\n")
              .append("\tastNodeBuilderPoolMu_.Lock()\n")
              .append("\tfor freeHead != nil {\n")
              .append("\t\tn := freeHead\n")
              .append("\t\tfreeHead = n.poolChildHead_\n")
              .append("\t\tn.poolChildHead_ = nil\n")
              .append("\t\tif len(astNodeBuilderFree_) < 262144 {\n")
              .append("\t\t\tastNodeBuilderFree_ = append(astNodeBuilderFree_, n)\n")
              .append("\t\t}\n")
              .append("\t}\n")
              .append("\tastNodeBuilderPoolMu_.Unlock()\n}\n");
          initBody.append("\tregisterCompPoolReset(func() {\n")
              .append("\t\tastNodeBuilderPoolMu_.Lock()\n")
              .append("\t\tastNodeBuilderBatch_.Store(nil)\n")
              .append("\t\tastNodeBuilderCurBatch_ = nil\n")
              .append("\t\tastNodeBuilderBatchFree_ = nil\n")
              .append("\t\tastNodeBuilderRetired_ = nil\n")
              .append("\t\tastNodeBuilderFree_ = nil\n")
              .append("\t\tastNodeBuilderPoolMu_.Unlock()\n")
              .append("\t})\n");
        }
      }
    }

    // 2. Emit default sub-message wrapper singletons after all `protoDefault_<Msg>_` are initialized.
    for (Map.Entry<String, String> e : sortedProtoTypes.entrySet()) {
      String qn = e.getKey();
      String msgGo = e.getValue();
      TypeElement msgTe = rb.resolver().elements.getTypeElement(qn);
      if (msgTe == null || !rb.resolver().getReferencedExternalTypes().contains(msgTe)) {
        continue;
      }
      Map<String, String> cachedSubMsgs = collectDefaultSubMsgGetters(rb, msgTe, msgGo);
      for (Map.Entry<String, String> sub : cachedSubMsgs.entrySet()) {
        String goMethod = sub.getKey();
        String retQn = sub.getValue();
        String retMsgGo = PROTO_MSG_TYPES.get(retQn);
        String retStruct = "Ext_" + GoNames.cleanChars(retQn.replace('.', '_'));
        String varName = "protoDefSub_" + msgGo + "_" + goMethod;
        sb.append("var ").append(varName).append(" *").append(retStruct).append('\n');
        initBody
            .append("\tif def := (&protobuf.")
            .append(msgGo)
            .append("{}).")
            .append(goMethod)
            .append("(); def != nil {\n")
            .append("\t\tif protoDefault_")
            .append(retMsgGo)
            .append("_ != nil && protoDefault_")
            .append(retMsgGo)
            .append("_.Core_ == def {\n")
            .append("\t\t\t")
            .append(varName)
            .append(" = protoDefault_")
            .append(retMsgGo)
            .append("_\n")
            .append("\t\t} else {\n")
            .append("\t\t\tw := &")
            .append(retStruct)
            .append("{}\n")
            .append("\t\t\tw.Core_ = def\n")
            .append("\t\t\t")
            .append(varName)
            .append(" = w\n")
            .append("\t\t}\n")
            .append("\t}\n");
      }
    }

    if (initBody.length() > 0) {
      sb.append("func init() {\n").append(initBody).append("}\n");
    }

    return sb.toString();
  }
}
