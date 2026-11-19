package j2go;

import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import javax.lang.model.element.AnnotationMirror;
import javax.lang.model.element.AnnotationValue;
import javax.lang.model.element.Element;
import javax.lang.model.element.ElementKind;
import javax.lang.model.element.ExecutableElement;
import javax.lang.model.element.Modifier;
import javax.lang.model.element.TypeElement;
import javax.lang.model.element.VariableElement;
import javax.lang.model.type.DeclaredType;
import javax.lang.model.type.TypeKind;
import javax.lang.model.type.TypeMirror;

/**
 * Clean-room mechanical emitter for {@code org.kohsuke.args4j} {@code @Option} and
 * {@code @Argument} field bindings and the {@code j2gort} args4j runtime support.
 */
public final class Args4jEmitter {

  private Args4jEmitter() {}

  private static final String OPT_ANN = "org.kohsuke.args4j.Option";
  private static final String ARG_ANN = "org.kohsuke.args4j.Argument";

  public static void emitArgs4jRegistrationIfNeeded(
      TypeElement te,
      String structName,
      TypeResolver resolver,
      StringBuilder out,
      EmitContext fileCtx) {
    List<FieldOpt> opts = collectOptions(te, resolver);
    if (opts.isEmpty()) {
      return;
    }

    String rt = resolver.runtime.j2gortAlias(fileCtx.imports);
    String pfx = rt.isEmpty() ? "" : (rt + ".");

    out.append("func (this_ *")
        .append(structName)
        .append(") RegisterArgs4jOptions(p *")
        .append(pfx)
        .append("Ext_org_kohsuke_args4j_CmdLineParser) {\n");
    out.append("\t_ = this_\n\t_ = p\n");

    int idx = 0;
    for (FieldOpt fo : opts) {
      emitSingleOptionRegistration(fo, idx++, pfx, resolver, out, fileCtx);
    }

    out.append("}\n\n");
  }

  private record FieldOpt(
      TypeElement owner,
      VariableElement field,
      boolean isArgument,
      String name,
      List<String> aliases,
      String usage,
      String metaVar,
      boolean required,
      boolean help,
      boolean hidden,
      Boolean explicitMultiValued,
      int index,
      TypeMirror handlerType) {}

  private static List<FieldOpt> collectOptions(TypeElement leafTe, TypeResolver resolver) {
    List<TypeElement> hierarchy = new ArrayList<>();
    TypeElement cur = leafTe;
    while (cur != null && resolver.isSourceType(cur)) {
      hierarchy.add(0, cur);
      cur = resolver.getSuperclassElement(cur);
    }
    List<FieldOpt> result = new ArrayList<>();
    for (TypeElement te : hierarchy) {
      for (Element enc : te.getEnclosedElements()) {
        if (enc instanceof VariableElement ve && ve.getKind() == ElementKind.FIELD) {
          for (AnnotationMirror am : ve.getAnnotationMirrors()) {
            Element annEl = am.getAnnotationType().asElement();
            if (!(annEl instanceof TypeElement annTe)) continue;
            String qn = annTe.getQualifiedName().toString();
            if (qn.equals(OPT_ANN) || qn.equals(ARG_ANN)) {
              result.add(parseAnnotation(te, ve, qn.equals(ARG_ANN), am));
            }
          }
        }
      }
    }
    return result;
  }

  private static FieldOpt parseAnnotation(
      TypeElement owner, VariableElement field, boolean isArgument, AnnotationMirror am) {
    String name = "";
    List<String> aliases = new ArrayList<>();
    String usage = "";
    String metaVar = "";
    boolean required = false;
    boolean help = false;
    boolean hidden = false;
    Boolean explicitMultiValued = null;
    int index = 0;
    TypeMirror handlerType = null;

    for (Map.Entry<? extends ExecutableElement, ? extends AnnotationValue> e :
        am.getElementValues().entrySet()) {
      String k = e.getKey().getSimpleName().toString();
      Object val = e.getValue().getValue();
      switch (k) {
        case "name" -> name = val != null ? val.toString() : "";
        case "aliases" -> {
          if (val instanceof List<?> list) {
            for (Object item : list) {
              Object iv = item instanceof AnnotationValue av ? av.getValue() : item;
              if (iv != null) aliases.add(iv.toString());
            }
          }
        }
        case "usage" -> usage = val != null ? val.toString() : "";
        case "metaVar" -> metaVar = val != null ? val.toString() : "";
        case "required" -> required = Boolean.TRUE.equals(val);
        case "help" -> help = Boolean.TRUE.equals(val);
        case "hidden" -> hidden = Boolean.TRUE.equals(val);
        case "multiValued" -> {
          if (val instanceof Boolean b) explicitMultiValued = b;
        }
        case "index" -> {
          if (val instanceof Integer n) index = n;
        }
        case "handler" -> {
          if (val instanceof TypeMirror tm) handlerType = tm;
        }
        default -> {}
      }
    }
    return new FieldOpt(
        owner,
        field,
        isArgument,
        name,
        aliases,
        usage,
        metaVar,
        required,
        help,
        help || hidden ? hidden : false,
        explicitMultiValued,
        index,
        handlerType);
  }

  private static void emitSingleOptionRegistration(
      FieldOpt fo,
      int idx,
      String pfx,
      TypeResolver resolver,
      StringBuilder out,
      EmitContext fileCtx) {
    VariableElement ve = fo.field;
    TypeMirror ft = ve.asType();
    TypeMirror rawFt = resolver.types.erasure(ft);
    String rawQn = rawFt.toString();

    boolean isList =
        rawQn.equals("java.util.List")
            || rawQn.equals("java.util.ArrayList")
            || rawQn.equals("java.util.Collection");
    boolean isMap = rawQn.equals("java.util.Map") || rawQn.equals("java.util.HashMap") || rawQn.equals("java.util.LinkedHashMap");
    boolean multiValued =
        fo.explicitMultiValued != null ? fo.explicitMultiValued : (isList || isMap);

    String fieldTargetExpr = buildFieldAccessExpr(fo, resolver, fileCtx);

    TypeMirror valueType = ft;
    if (isList && ft instanceof DeclaredType dt && !dt.getTypeArguments().isEmpty()) {
      valueType = dt.getTypeArguments().get(0);
    }

    String optVar = "opt_" + idx + "_";
    String setterVar = "set_" + idx + "_";
    out.append("\t{\n");

    // 1. OptionDef / NamedOptionDef
    String aliasesLit = "nil";
    if (!fo.aliases.isEmpty()) {
      List<String> quoted = new ArrayList<>();
      for (String a : fo.aliases) quoted.add(GoNames.goStringLiteral(a));
      aliasesLit = "[]string{" + String.join(", ", quoted) + "}";
    }
    if (fo.isArgument) {
      out.append("\t\t")
          .append(optVar)
          .append(" := ")
          .append(pfx)
          .append("NewArgs4jArgumentDef(")
          .append(GoNames.goStringLiteral(fo.usage))
          .append(", ")
          .append(GoNames.goStringLiteral(fo.metaVar))
          .append(", ")
          .append(fo.required)
          .append(", ")
          .append(fo.hidden)
          .append(", ")
          .append(multiValued)
          .append(", ")
          .append(fo.index)
          .append(")\n");
    } else {
      out.append("\t\t")
          .append(optVar)
          .append(" := ")
          .append(pfx)
          .append("NewArgs4jNamedOptionDef(")
          .append(GoNames.goStringLiteral(fo.name))
          .append(", ")
          .append(aliasesLit)
          .append(", ")
          .append(GoNames.goStringLiteral(fo.usage))
          .append(", ")
          .append(GoNames.goStringLiteral(fo.metaVar))
          .append(", ")
          .append(fo.required)
          .append(", ")
          .append(fo.help)
          .append(", ")
          .append(fo.hidden)
          .append(", ")
          .append(multiValued)
          .append(")\n");
    }

    // 2. Setter / FieldSetter
    String classLit = emitClassLiteral(valueType, pfx, resolver);
    out.append("\t\t")
        .append(setterVar)
        .append(" := ")
        .append(pfx)
        .append("NewArgs4jFieldSetter(")
        .append(multiValued)
        .append(", ")
        .append(classLit)
        .append(", func(val_ any) {\n");
    emitSetterBody(fieldTargetExpr, ft, valueType, isList, isMap, pfx, resolver, out, fileCtx);
    out.append("\t\t}, func() any {\n");
    emitGetterBody(fieldTargetExpr, ft, multiValued, pfx, resolver, out);
    out.append("\t\t})\n");

    // 3. OptionHandler (custom source handler or built-in runtime handler)
    String handlerExpr =
        buildHandlerExpr(
            fo, valueType, isMap, optVar, setterVar, pfx, resolver, fileCtx);
    out.append("\t\t")
        .append(pfx)
        .append("Args4jRegisterHandler(p, ")
        .append(optVar)
        .append(", ")
        .append(handlerExpr)
        .append(", ")
        .append(fo.isArgument)
        .append(")\n");
    out.append("\t}\n");
  }

  private static String buildFieldAccessExpr(
      FieldOpt fo, TypeResolver resolver, EmitContext fileCtx) {
    VariableElement ve = fo.field;
    String fGo = resolver.getFieldGoName(ve);
    if (ve.getModifiers().contains(Modifier.STATIC)) {
      String ownerGo = resolver.planner.getGoTypeName(fo.owner);
      String ownerPkg = resolver.planner.getGoPkg(fo.owner);
      String al = fileCtx.imports.ext(resolver.planner.getGoImportPath(ownerPkg), ownerPkg);
      String pref = al.isEmpty() ? "" : (al + ".");
      return pref + "static_" + ownerGo + "_" + fGo;
    }
    return "this_." + fGo;
  }

  private static String emitClassLiteral(
      TypeMirror tm, String pfx, TypeResolver resolver) {
    TypeMirror er = resolver.types.erasure(tm);
    if (er.getKind().isPrimitive()) {
      return pfx + "ClassLiteral(" + GoNames.goStringLiteral(er.toString()) + ", " + GoNames.goStringLiteral(er.toString()) + ")";
    }
    if (er instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      String qn = te.getQualifiedName().toString();
      String sn = te.getSimpleName().toString();
      return pfx + "ClassLiteral(" + GoNames.goStringLiteral(qn) + ", " + GoNames.goStringLiteral(sn) + ")";
    }
    return pfx + "ClassLiteral(\"java.lang.Object\", \"Object\")";
  }

  private static void emitSetterBody(
      String target,
      TypeMirror fieldType,
      TypeMirror valueType,
      boolean isList,
      boolean isMap,
      String pfx,
      TypeResolver resolver,
      StringBuilder out,
      EmitContext fileCtx) {
    if (isMap) {
      out.append("\t\t\tif ").append(target).append(" == nil {\n");
      out.append("\t\t\t\t")
          .append(target)
          .append(" = ")
          .append(pfx)
          .append("NewArgs4jStringMap()\n");
      out.append("\t\t\t}\n");
      out.append("\t\t\tk_, v_ := ").append(pfx).append("Args4jSplitMapEntry(p, val_)\n");
      out.append("\t\t\t").append(pfx).append("Args4jMapPut(").append(target).append(", k_, v_)\n");
      return;
    }

    String convertedVal = convertRawArgValue("val_", valueType, pfx, resolver, fileCtx);
    if (isList) {
      out.append("\t\t\tif ").append(target).append(" == nil {\n");
      out.append("\t\t\t\t")
          .append(target)
          .append(" = ")
          .append(pfx)
          .append("NewArgs4jList()\n");
      out.append("\t\t\t}\n");
      String boxedElem = boxForCollection(convertedVal, valueType, pfx, resolver);
      out.append("\t\t\t").append(pfx).append("Args4jListAdd(").append(target).append(", ").append(boxedElem).append(")\n");
      return;
    }

    out.append("\t\t\t").append(target).append(" = ").append(convertedVal).append('\n');
  }

  private static void emitGetterBody(
      String target,
      TypeMirror fieldType,
      boolean multiValued,
      String pfx,
      TypeResolver resolver,
      StringBuilder out) {
    if (multiValued) {
      out.append("\t\t\treturn nil\n");
      return;
    }
    TypeKind k = fieldType.getKind();
    if (k == TypeKind.BOOLEAN) {
      out.append("\t\t\treturn ").append(pfx).append("ValueOfBoolean(").append(target).append(")\n");
      return;
    }
    if (k == TypeKind.INT) {
      out.append("\t\t\treturn ").append(pfx).append("ValueOfInt(").append(target).append(")\n");
      return;
    }
    if (k == TypeKind.LONG) {
      out.append("\t\t\treturn ").append(pfx).append("ValueOfLong(").append(target).append(")\n");
      return;
    }
    if (k == TypeKind.DOUBLE) {
      out.append("\t\t\treturn ").append(pfx).append("ValueOfDouble(").append(target).append(")\n");
      return;
    }
    if (k == TypeKind.FLOAT) {
      out.append("\t\t\treturn ").append(pfx).append("ValueOfFloat(").append(target).append(")\n");
      return;
    }
    if (resolver.isStringType(fieldType)) {
      out.append("\t\t\tif ")
          .append(target)
          .append(" == ")
          .append(pfx)
          .append("NullString { return nil }\n");
      out.append("\t\t\treturn ").append(pfx).append("BoxString(").append(target).append(")\n");
      return;
    }
    out.append("\t\t\treturn ").append(pfx).append("ToAny(").append(target).append(")\n");
  }

  private static String convertRawArgValue(
      String valExpr,
      TypeMirror targetType,
      String pfx,
      TypeResolver resolver,
      EmitContext fileCtx) {
    TypeKind k = targetType.getKind();
    if (k == TypeKind.BOOLEAN) return pfx + "Args4jAsBool(p, " + valExpr + ")";
    if (k == TypeKind.INT) return pfx + "Args4jAsInt(p, " + valExpr + ")";
    if (k == TypeKind.LONG) return pfx + "Args4jAsLong(p, " + valExpr + ")";
    if (k == TypeKind.DOUBLE) return pfx + "Args4jAsDouble(p, " + valExpr + ")";
    if (k == TypeKind.FLOAT) return "float32(" + pfx + "Args4jAsDouble(p, " + valExpr + "))";

    if (targetType instanceof DeclaredType dt && dt.asElement() instanceof TypeElement te) {
      String qn = te.getQualifiedName().toString();
      switch (qn) {
        case "java.lang.String" -> {
          return pfx + "Args4jAsString(" + valExpr + ")";
        }
        case "java.lang.Boolean" -> {
          return pfx + "ValueOfBoolean(" + pfx + "Args4jAsBool(p, " + valExpr + "))";
        }
        case "java.lang.Integer" -> {
          return pfx + "ValueOfInt(" + pfx + "Args4jAsInt(p, " + valExpr + "))";
        }
        case "java.lang.Long" -> {
          return pfx + "ValueOfLong(" + pfx + "Args4jAsLong(p, " + valExpr + "))";
        }
        case "java.lang.Double" -> {
          return pfx + "ValueOfDouble(" + pfx + "Args4jAsDouble(p, " + valExpr + "))";
        }
        default -> {}
      }
      if (te.getKind() == ElementKind.ENUM) {
        String goType = resolver.goTypeForElement(te, fileCtx.imports);
        String fqnLit = GoNames.goStringLiteral(qn);
        if (resolver.isSourceType(te)) {
          String enumPkg = resolver.planner.getGoPkg(te);
          String al = fileCtx.imports.ext(resolver.planner.getGoImportPath(enumPkg), enumPkg);
          String ep = al.isEmpty() ? "" : (al + ".");
          String enumGo = resolver.planner.getGoTypeName(te);
          // Ensure enum clinit has run so RegisterEnumConstants has populated!
          return "func() "
              + goType
              + " { _ = "
              + ep
              + enumGo
              + "_Values(); return "
              + pfx
              + "Cast["
              + goType
              + "]("
              + pfx
              + "Args4jParseEnum(p, "
              + fqnLit
              + ", "
              + valExpr
              + ")) }()";
        } else {
          return pfx + "Cast[" + goType + "](" + pfx + "Args4jParseEnum(p, " + fqnLit + ", " + valExpr + "))";
        }
      }
    }
    String goT = resolver.goType(targetType, fileCtx.imports);
    return pfx + "Cast[" + goT + "](" + valExpr + ")";
  }

  private static String boxForCollection(
      String expr, TypeMirror elemType, String pfx, TypeResolver resolver) {
    if (resolver.isStringType(elemType)) {
      return pfx + "BoxString(" + expr + ")";
    }
    if (elemType.getKind() == TypeKind.BOOLEAN) {
      return pfx + "ValueOfBoolean(" + expr + ")";
    }
    if (elemType.getKind() == TypeKind.INT) {
      return pfx + "ValueOfInt(" + expr + ")";
    }
    return pfx + "ToAny(" + expr + ")";
  }

  private static String buildHandlerExpr(
      FieldOpt fo,
      TypeMirror valueType,
      boolean isMap,
      String optVar,
      String setterVar,
      String pfx,
      TypeResolver resolver,
      EmitContext fileCtx) {
    if (fo.handlerType != null
        && fo.handlerType instanceof DeclaredType dt
        && dt.asElement() instanceof TypeElement hTe) {
      String hQn = hTe.getQualifiedName().toString();
      if (!hQn.equals("org.kohsuke.args4j.spi.OptionHandler")) {
        if (resolver.isSourceType(hTe)) {
          for (Element m : hTe.getEnclosedElements()) {
            if (m instanceof ExecutableElement ctor
                && ctor.getKind() == ElementKind.CONSTRUCTOR
                && ctor.getParameters().size() == 3) {
              String hPkg = resolver.planner.getGoPkg(hTe);
              String al = fileCtx.imports.ext(resolver.planner.getGoImportPath(hPkg), hPkg);
              String hp = al.isEmpty() ? "" : (al + ".");
              String ctorName = resolver.getConstructorGoName(ctor);
              return pfx
                  + "AsIface["
                  + pfx
                  + "Ext_org_kohsuke_args4j_spi_OptionHandler]("
                  + hp
                  + ctorName
                  + "(p, "
                  + optVar
                  + ", "
                  + setterVar
                  + "))";
            }
          }
        }
        if (hQn.equals("org.kohsuke.args4j.spi.IntOptionHandler")) {
          return pfx + "NewArgs4jBuiltinHandler(p, " + optVar + ", " + setterVar + ", \"int\", \"\")";
        }
        if (hQn.equals("org.kohsuke.args4j.spi.BooleanOptionHandler")) {
          return pfx + "NewArgs4jBuiltinHandler(p, " + optVar + ", " + setterVar + ", \"boolean\", \"\")";
        }
        if (hQn.equals("org.kohsuke.args4j.spi.ExplicitBooleanOptionHandler")) {
          return pfx + "NewArgs4jBuiltinHandler(p, " + optVar + ", " + setterVar + ", \"explicit_boolean\", \"\")";
        }
        if (hQn.equals("org.kohsuke.args4j.spi.StringOptionHandler")) {
          return pfx + "NewArgs4jBuiltinHandler(p, " + optVar + ", " + setterVar + ", \"string\", \"\")";
        }
        if (hQn.equals("org.kohsuke.args4j.spi.MapOptionHandler")) {
          return pfx + "NewArgs4jBuiltinHandler(p, " + optVar + ", " + setterVar + ", \"map\", \"\")";
        }
      }
    }

    if (isMap) {
      return pfx + "NewArgs4jBuiltinHandler(p, " + optVar + ", " + setterVar + ", \"map\", \"\")";
    }
    TypeKind k = valueType.getKind();
    if (k == TypeKind.BOOLEAN) {
      return pfx + "NewArgs4jBuiltinHandler(p, " + optVar + ", " + setterVar + ", \"boolean\", \"\")";
    }
    if (k == TypeKind.INT || k == TypeKind.LONG) {
      return pfx + "NewArgs4jBuiltinHandler(p, " + optVar + ", " + setterVar + ", \"int\", \"\")";
    }
    if (k == TypeKind.DOUBLE || k == TypeKind.FLOAT) {
      return pfx + "NewArgs4jBuiltinHandler(p, " + optVar + ", " + setterVar + ", \"double\", \"\")";
    }
    if (valueType instanceof DeclaredType dt && dt.asElement() instanceof TypeElement vTe) {
      String qn = vTe.getQualifiedName().toString();
      switch (qn) {
        case "java.lang.Boolean" -> {
          return pfx + "NewArgs4jBuiltinHandler(p, " + optVar + ", " + setterVar + ", \"boolean\", \"\")";
        }
        case "java.lang.Integer", "java.lang.Long" -> {
          return pfx + "NewArgs4jBuiltinHandler(p, " + optVar + ", " + setterVar + ", \"int\", \"\")";
        }
        case "java.lang.Double", "java.lang.Float" -> {
          return pfx + "NewArgs4jBuiltinHandler(p, " + optVar + ", " + setterVar + ", \"double\", \"\")";
        }
        default -> {}
      }
      if (vTe.getKind() == ElementKind.ENUM) {
        String ensureInit = "";
        if (resolver.isSourceType(vTe)) {
          String enumPkg = resolver.planner.getGoPkg(vTe);
          String al = fileCtx.imports.ext(resolver.planner.getGoImportPath(enumPkg), enumPkg);
          String ep = al.isEmpty() ? "" : (al + ".");
          String enumGo = resolver.planner.getGoTypeName(vTe);
          ensureInit = "_ = " + ep + enumGo + "_Values(); ";
        }
        return "func() "
            + pfx
            + "Ext_org_kohsuke_args4j_spi_OptionHandler { "
            + ensureInit
            + "return "
            + pfx
            + "NewArgs4jBuiltinHandler(p, "
            + optVar
            + ", "
            + setterVar
            + ", \"enum\", "
            + GoNames.goStringLiteral(qn)
            + ") }()";
      }
    }
    return pfx + "NewArgs4jBuiltinHandler(p, " + optVar + ", " + setterVar + ", \"string\", \"\")";
  }

  public static String buildRuntimeArgs4jGo() {
    return """
        // Clean-room org.kohsuke.args4j runtime implementation for j2gort.
        package j2gort

        import (
        \t"fmt"
        \t"os"
        \t"sort"
        \t"strconv"
        \t"strings"

        \t"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
        \t"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
        )

        type args4jParserCore struct {
        \tbean       any
        \toptions    []args4jHandlerEntry
        \targuments  []args4jHandlerEntry
        \tbyName     map[string]Ext_org_kohsuke_args4j_spi_OptionHandler
        \tusageWidth int
        }

        type args4jHandlerEntry struct {
        \tdef     Ext_org_kohsuke_args4j_OptionDef
        \thandler Ext_org_kohsuke_args4j_spi_OptionHandler
        }

        type args4jOptionBinder interface {
        \tRegisterArgs4jOptions(p *Ext_org_kohsuke_args4j_CmdLineParser)
        }

        func initArgs4jParser(p *Ext_org_kohsuke_args4j_CmdLineParser, bean any) {
        \tcore := &args4jParserCore{
        \t\tbean:       bean,
        \t\tbyName:     make(map[string]Ext_org_kohsuke_args4j_spi_OptionHandler),
        \t\tusageWidth: 80,
        \t}
        \tp.Core_ = core
        \tif binder, ok := bean.(args4jOptionBinder); ok && !javalang.IsNil(bean) {
        \t\tbinder.RegisterArgs4jOptions(p)
        \t}
        }

        func NewArgs4jStringMap() Ext_java_util_Map {
        \treturn wrapLinkedHashMap(NewJMapCore(false))
        }

        func NewArgs4jList() Ext_java_util_List {
        \treturn wrapArrayList(nil)
        }

        func Args4jMapPut(m Ext_java_util_Map, k, v string) {
        \tif mc, ok := any(m).(interface{ mc() *JMapCore }); ok {
        \t\tmc.mc().Put(BoxString(k), BoxString(v))
        \t}
        }

        func Args4jListAdd(l Ext_java_util_List, elem any) {
        \tif lc, ok := any(l).(interface{ lc() *JListCore }); ok {
        \t\tlc.lc().Add(elem)
        \t}
        }

        func (p *Ext_org_kohsuke_args4j_CmdLineParser) args4jCore() *args4jParserCore {
        \tif p == nil {
        \t\treturn nil
        \t}
        \tif c, ok := p.Core_.(*args4jParserCore); ok && c != nil {
        \t\treturn c
        \t}
        \tc := &args4jParserCore{byName: make(map[string]Ext_org_kohsuke_args4j_spi_OptionHandler), usageWidth: 80}
        \tp.Core_ = c
        \treturn c
        }

        func NewArgs4jCmdLineException(p *Ext_org_kohsuke_args4j_CmdLineParser, msg string) Ext_org_kohsuke_args4j_CmdLineException {
        \t_ = p
        \tex := &Ext_org_kohsuke_args4j_CmdLineExceptionBase{}
        \tex.Self_ = ex
        \tjavalang.InitThrowable(ex, "org.kohsuke.args4j.CmdLineException", msg, nil)
        \treturn ex
        }

        func NewArgs4jCmdLineExceptionCause(p *Ext_org_kohsuke_args4j_CmdLineParser, cause javalang.Throwable) Ext_org_kohsuke_args4j_CmdLineException {
        \t_ = p
        \tex := &Ext_org_kohsuke_args4j_CmdLineExceptionBase{}
        \tex.Self_ = ex
        \tmsg := ""
        \tif cause != nil {
        \t\tmsg = cause.GetMessage()
        \t}
        \tjavalang.InitThrowable(ex, "org.kohsuke.args4j.CmdLineException", msg, cause)
        \treturn ex
        }

        type args4jOptionDefData struct {
        \tisNamed     bool
        \tname        string
        \taliases     []string
        \tusage       string
        \tmetaVar     string
        \trequired    bool
        \thelp        bool
        \thidden      bool
        \tmultiValued bool
        \tindex       int
        }

        func NewArgs4jNamedOptionDef(
        \tname string,
        \taliases []string,
        \tusage string,
        \tmetaVar string,
        \trequired bool,
        \thelp bool,
        \thidden bool,
        \tmultiValued bool,
        ) *Ext_org_kohsuke_args4j_NamedOptionDef {
        \td := &args4jOptionDefData{
        \t\tisNamed:     true,
        \t\tname:        name,
        \t\taliases:     aliases,
        \t\tusage:       usage,
        \t\tmetaVar:     metaVar,
        \t\trequired:    required,
        \t\thelp:        help,
        \t\thidden:      hidden,
        \t\tmultiValued: multiValued,
        \t}
        \tobj := &Ext_org_kohsuke_args4j_NamedOptionDef{}
        \tobj.Core_ = d
        \tobj.Ext_org_kohsuke_args4j_OptionDefBase.Core_ = d
        \tobj.Ext_org_kohsuke_args4j_OptionDefBase.Self_ = obj
        \treturn obj
        }

        func NewArgs4jArgumentDef(
        \tusage string,
        \tmetaVar string,
        \trequired bool,
        \thidden bool,
        \tmultiValued bool,
        \tindex int,
        ) Ext_org_kohsuke_args4j_OptionDef {
        \td := &args4jOptionDefData{
        \t\tisNamed:     false,
        \t\tusage:       usage,
        \t\tmetaVar:     metaVar,
        \t\trequired:    required,
        \t\thidden:      hidden,
        \t\tmultiValued: multiValued,
        \t\tindex:       index,
        \t}
        \tobj := &Ext_org_kohsuke_args4j_OptionDefBase{Core_: d}
        \tobj.Self_ = obj
        \treturn obj
        }

        func getOptionDefData(core any) *args4jOptionDefData {
        \tif d, ok := core.(*args4jOptionDefData); ok && d != nil {
        \t\treturn d
        \t}
        \treturn &args4jOptionDefData{}
        }

        type args4jFieldSetterData struct {
        \tmultiValued bool
        \ttyp         *javalang.Class
        \taddFn       func(any)
        \tgetFn       func() any
        }

        func NewArgs4jFieldSetter(multiValued bool, typ *javalang.Class, addFn func(any), getFn func() any) *Ext_org_kohsuke_args4j_spi_FieldSetter {
        \td := &args4jFieldSetterData{
        \t\tmultiValued: multiValued,
        \t\ttyp:         typ,
        \t\taddFn:       addFn,
        \t\tgetFn:       getFn,
        \t}
        \treturn &Ext_org_kohsuke_args4j_spi_FieldSetter{Core_: d}
        }

        type args4jParametersImpl struct {
        \tp       *Ext_org_kohsuke_args4j_CmdLineParser
        \toptName string
        \targs    []string
        \tpos     int
        }

        func (*args4jParametersImpl) Implements_j2gort_Ext_org_kohsuke_args4j_spi_Parameters() {}

        func (pr *args4jParametersImpl) GetParameter(idx int32) string {
        \ti := pr.pos + int(idx)
        \tif i < 0 || i >= len(pr.args) {
        \t\tif pr.optName != "" {
        \t\t\tpanic(NewArgs4jCmdLineException(pr.p, fmt.Sprintf("Option %q takes an operand", pr.optName)))
        \t\t}
        \t\tpanic(NewArgs4jCmdLineException(pr.p, "Option takes an operand"))
        \t}
        \treturn pr.args[i]
        }

        func (pr *args4jParametersImpl) Size() int32 {
        \treturn int32(len(pr.args) - pr.pos)
        }

        type args4jBuiltinHandler struct {
        \tExt_org_kohsuke_args4j_spi_OptionHandlerBase
        \tkind    string
        \tenumFqn string
        }

        func NewArgs4jBuiltinHandler(
        \tp *Ext_org_kohsuke_args4j_CmdLineParser,
        \topt Ext_org_kohsuke_args4j_OptionDef,
        \tsetter Ext_org_kohsuke_args4j_spi_Setter,
        \tkind string,
        \tenumFqn string,
        ) Ext_org_kohsuke_args4j_spi_OptionHandler {
        \th := &args4jBuiltinHandler{kind: kind, enumFqn: enumFqn}
        \th.Self_ = h
        \tInitExt_org_kohsuke_args4j_spi_OptionHandler(&h.Ext_org_kohsuke_args4j_spi_OptionHandlerBase, p, opt, setter)
        \treturn h
        }

        func (h *args4jBuiltinHandler) ParseArguments(params Ext_org_kohsuke_args4j_spi_Parameters) int32 {
        \tswitch h.kind {
        \tcase "boolean":
        \t\th.F_Setter.AddValue(ValueOfBoolean(true))
        \t\treturn 0
        \tcase "explicit_boolean":
        \t\tvar tok string
        \t\thasTok := false
        \t\tfunc() {
        \t\t\tdefer func() { _ = recover() }()
        \t\t\ttok = params.GetParameter(0)
        \t\t\thasTok = true
        \t\t}()
        \t\tif !hasTok {
        \t\t\th.F_Setter.AddValue(ValueOfBoolean(true))
        \t\t\treturn 0
        \t\t}
        \t\tlow := strings.ToLower(tok)
        \t\tswitch low {
        \t\tcase "true", "on", "yes", "1":
        \t\t\th.F_Setter.AddValue(ValueOfBoolean(true))
        \t\t\treturn 1
        \t\tcase "false", "off", "no", "0":
        \t\t\th.F_Setter.AddValue(ValueOfBoolean(false))
        \t\t\treturn 1
        \t\tdefault:
        \t\t\th.F_Setter.AddValue(ValueOfBoolean(true))
        \t\t\treturn 0
        \t\t}
        \tdefault:
        \t\ts := params.GetParameter(0)
        \t\th.F_Setter.AddValue(BoxString(s))
        \t\treturn 1
        \t}
        }

        func (h *args4jBuiltinHandler) GetDefaultMetaVariable() string {
        \tswitch h.kind {
        \tcase "boolean":
        \t\treturn NullString
        \tcase "explicit_boolean":
        \t\treturn "[true | false]"
        \tcase "int", "double":
        \t\treturn "N"
        \tcase "map":
        \t\treturn "key=value"
        \tcase "enum":
        \t\tconsts := LookupEnumConstantsByName(h.enumFqn)
        \t\tif len(consts) > 0 {
        \t\t\tnames := make([]string, 0, len(consts))
        \t\t\tfor _, c := range consts {
        \t\t\t\tif n, ok := c.(interface{ Name() string }); ok {
        \t\t\t\t\tnames = append(names, n.Name())
        \t\t\t\t}
        \t\t\t}
        \t\t\tif len(names) > 0 {
        \t\t\t\treturn "[" + strings.Join(names, " | ") + "]"
        \t\t\t}
        \t\t}
        \t\treturn "VAL"
        \tdefault:
        \t\treturn "VAL"
        \t}
        }

        func Args4jRegisterHandler(
        \tp *Ext_org_kohsuke_args4j_CmdLineParser,
        \tdef Ext_org_kohsuke_args4j_OptionDef,
        \thandler Ext_org_kohsuke_args4j_spi_OptionHandler,
        \tisArgument bool,
        ) {
        \tc := p.args4jCore()
        \tif isArgument {
        \t\tc.arguments = append(c.arguments, args4jHandlerEntry{def: def, handler: handler})
        \t\treturn
        \t}
        \tc.options = append(c.options, args4jHandlerEntry{def: def, handler: handler})
        \tif nd, ok := def.(*Ext_org_kohsuke_args4j_NamedOptionDef); ok && nd != nil {
        \t\tc.byName[nd.Name()] = handler
        \t\td := getOptionDefData(nd.Core_)
        \t\tfor _, al := range d.aliases {
        \t\t\tc.byName[al] = handler
        \t\t}
        \t}
        }

        func Args4jAsString(v any) string {
        \tif v == nil {
        \t\treturn NullString
        \t}
        \tif s, ok := v.(string); ok {
        \t\treturn s
        \t}
        \treturn StrVal(v)
        }

        func Args4jAsBool(p *Ext_org_kohsuke_args4j_CmdLineParser, v any) bool {
        \tswitch x := v.(type) {
        \tcase bool:
        \t\treturn x
        \tcase *BoxBoolean:
        \t\treturn x != nil && x.Val
        \t}
        \ts := strings.ToLower(strings.TrimSpace(Args4jAsString(v)))
        \tswitch s {
        \tcase "true", "on", "yes", "1":
        \t\treturn true
        \tcase "false", "off", "no", "0":
        \t\treturn false
        \tdefault:
        \t\tpanic(NewArgs4jCmdLineException(p, fmt.Sprintf("%q is not a valid boolean value", s)))
        \t}
        }

        func Args4jAsInt(p *Ext_org_kohsuke_args4j_CmdLineParser, v any) int32 {
        \tswitch x := v.(type) {
        \tcase int32:
        \t\treturn x
        \tcase *BoxInt:
        \t\tif x != nil {
        \t\t\treturn x.Val
        \t\t}
        \t}
        \ts := Args4jAsString(v)
        \tn, err := strconv.ParseInt(strings.TrimSpace(s), 10, 32)
        \tif err != nil {
        \t\tpanic(NewArgs4jCmdLineException(p, fmt.Sprintf("%q is not a valid number", s)))
        \t}
        \treturn int32(n)
        }

        func Args4jAsLong(p *Ext_org_kohsuke_args4j_CmdLineParser, v any) int64 {
        \tswitch x := v.(type) {
        \tcase int64:
        \t\treturn x
        \tcase *BoxLong:
        \t\tif x != nil {
        \t\t\treturn x.Val
        \t\t}
        \t}
        \ts := Args4jAsString(v)
        \tn, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
        \tif err != nil {
        \t\tpanic(NewArgs4jCmdLineException(p, fmt.Sprintf("%q is not a valid number", s)))
        \t}
        \treturn n
        }

        func Args4jAsDouble(p *Ext_org_kohsuke_args4j_CmdLineParser, v any) float64 {
        \tswitch x := v.(type) {
        \tcase float64:
        \t\treturn x
        \tcase *BoxDouble:
        \t\tif x != nil {
        \t\t\treturn x.Val
        \t\t}
        \t}
        \ts := Args4jAsString(v)
        \tn, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
        \tif err != nil {
        \t\tpanic(NewArgs4jCmdLineException(p, fmt.Sprintf("%q is not a valid number", s)))
        \t}
        \treturn n
        }

        func Args4jParseEnum(p *Ext_org_kohsuke_args4j_CmdLineParser, fqn string, v any) any {
        \tif v != nil {
        \t\tif _, isStr := v.(string); !isStr {
        \t\t\treturn v
        \t\t}
        \t}
        \ts := Args4jAsString(v)
        \tnorm := strings.ReplaceAll(s, "-", "_")
        \tconsts := LookupEnumConstantsByName(fqn)
        \tfor _, c := range consts {
        \t\tif n, ok := c.(interface{ Name() string }); ok && n.Name() == s {
        \t\t\treturn c
        \t\t}
        \t}
        \tfor _, c := range consts {
        \t\tif n, ok := c.(interface{ Name() string }); ok && (strings.EqualFold(n.Name(), s) || strings.EqualFold(n.Name(), norm)) {
        \t\t\treturn c
        \t\t}
        \t}
        \tpanic(NewArgs4jCmdLineException(p, fmt.Sprintf("%q is not a valid value", s)))
        }

        func Args4jSplitMapEntry(p *Ext_org_kohsuke_args4j_CmdLineParser, v any) (string, string) {
        \ts := Args4jAsString(v)
        \tidx := strings.IndexByte(s, '=')
        \tif idx < 0 {
        \t\tpanic(NewArgs4jCmdLineException(p, "An argument for setting a Map must contain a \\"=\\""))
        \t}
        \treturn s[:idx], s[idx+1:]
        }

        func args4jParseTokens(p *Ext_org_kohsuke_args4j_CmdLineParser, rawTokens []string) {
        \tc := p.args4jCore()
        \ttokens := make([]string, 0, len(rawTokens))
        \tfor _, t := range rawTokens {
        \t\tif strings.HasPrefix(t, "@") && len(t) > 1 {
        \t\t\tfilePath := t[1:]
        \t\t\tdata, err := os.ReadFile(filePath)
        \t\t\tif err != nil {
        \t\t\t\tpanic(NewArgs4jCmdLineException(p, "No such file: "+filePath))
        \t\t\t}
        \t\t\tfor _, line := range strings.Split(strings.ReplaceAll(string(data), "\\r\\n", "\\n"), "\\n") {
        \t\t\t\ttrimmed := strings.TrimSpace(line)
        \t\t\t\tif trimmed == "" || strings.HasPrefix(trimmed, "#") {
        \t\t\t\t\tcontinue
        \t\t\t\t}
        \t\t\t\ttokens = append(tokens, trimmed)
        \t\t\t}
        \t\t} else {
        \t\t\ttokens = append(tokens, t)
        \t\t}
        \t}

        \tendOfOpts := false
        \ti := 0
        \targHandlerIdx := 0
        \tfor i < len(tokens) {
        \t\ttok := tokens[i]
        \t\tif !endOfOpts && tok == "--" {
        \t\t\tendOfOpts = true
        \t\t\ti++
        \t\t\tcontinue
        \t\t}
        \t\tif !endOfOpts && strings.HasPrefix(tok, "-") && len(tok) > 1 {
        \t\t\t// Look up exact option name first, then --key=value split
        \t\t\tif h, ok := c.byName[tok]; ok {
        \t\t\t\tpr := &args4jParametersImpl{p: p, optName: tok, args: tokens, pos: i + 1}
        \t\t\t\tconsumed := int(h.ParseArguments(pr))
        \t\t\t\ti += 1 + consumed
        \t\t\t\tcontinue
        \t\t\t}
        \t\t\tif eq := strings.IndexByte(tok, '='); eq > 0 {
        \t\t\t\toptKey := tok[:eq]
        \t\t\t\toptVal := tok[eq+1:]
        \t\t\t\tif len(optVal) >= 2 && ((optVal[0] == '\\'' && optVal[len(optVal)-1] == '\\'') || (optVal[0] == '"' && optVal[len(optVal)-1] == '"')) {
        \t\t\t\t\toptVal = optVal[1 : len(optVal)-1]
        \t\t\t\t}
        \t\t\t\tif h, ok := c.byName[optKey]; ok {
        \t\t\t\t\tpr := &args4jParametersImpl{p: p, optName: optKey, args: []string{optVal}, pos: 0}
        \t\t\t\t\t_ = h.ParseArguments(pr)
        \t\t\t\t\ti++
        \t\t\t\t\tcontinue
        \t\t\t\t}
        \t\t\t}
        \t\t\tpanic(NewArgs4jCmdLineException(p, fmt.Sprintf("%q is not a valid option", tok)))
        \t\t}

        \t\t// Positional argument
        \t\tif len(c.arguments) == 0 {
        \t\t\tpanic(NewArgs4jCmdLineException(p, fmt.Sprintf("No argument is allowed: %s", tok)))
        \t\t}
        \t\tidx := argHandlerIdx
        \t\tif idx >= len(c.arguments) {
        \t\t\tidx = len(c.arguments) - 1
        \t\t}
        \t\tent := c.arguments[idx]
        \t\tpr := &args4jParametersImpl{p: p, args: tokens, pos: i}
        \t\tconsumed := int(ent.handler.ParseArguments(pr))
        \t\tif consumed <= 0 {
        \t\t\tconsumed = 1
        \t\t}
        \t\ti += consumed
        \t\tif !getOptionDefData(ent.def.GetExt_org_kohsuke_args4j_OptionDefBase().Core_).multiValued {
        \t\t\targHandlerIdx++
        \t\t}
        \t}
        }

        func args4jDefaultSuffix(h Ext_org_kohsuke_args4j_spi_OptionHandler) string {
        \tif h == nil {
        \t\treturn ""
        \t}
        \tsetter := h.GetExt_org_kohsuke_args4j_spi_OptionHandlerBase().F_Setter
        \tfs, ok := setter.(*Ext_org_kohsuke_args4j_spi_FieldSetter)
        \tif !ok || fs == nil {
        \t\treturn ""
        \t}
        \td, ok := fs.Core_.(*args4jFieldSetterData)
        \tif !ok || d == nil || d.multiValued || d.getFn == nil {
        \t\treturn ""
        \t}
        \tval := d.getFn()
        \tif javalang.IsNil(val) {
        \t\treturn ""
        \t}
        \treturn " (default: " + StrVal(val) + ")"
        }

        func args4jOptionHead(ent args4jHandlerEntry) (string, *args4jOptionDefData, bool) {
        \tnd, ok := ent.def.(*Ext_org_kohsuke_args4j_NamedOptionDef)
        \tif !ok || nd == nil {
        \t\treturn "", nil, false
        \t}
        \td := getOptionDefData(nd.Core_)
        \tif d.usage == "" {
        \t\treturn "", d, false
        \t}
        \tvar hb strings.Builder
        \thb.WriteString(" ")
        \thb.WriteString(d.name)
        \tif len(d.aliases) > 0 {
        \t\thb.WriteString(" (")
        \t\thb.WriteString(strings.Join(d.aliases, ", "))
        \t\thb.WriteString(")")
        \t}
        \tmeta := d.metaVar
        \tif meta == "" {
        \t\tmv := ent.handler.GetDefaultMetaVariable()
        \t\tif mv != NullString && mv != "" {
        \t\t\tmeta = mv
        \t\t}
        \t}
        \tif meta != "" {
        \t\thb.WriteString(" ")
        \t\thb.WriteString(meta)
        \t}
        \treturn hb.String(), d, true
        }

        func args4jWrapLines(s string, w int, indentCont bool) []string {
        \tif w <= 0 {
        \t\treturn strings.Split(s, "\\n")
        \t}
        \tvar out []string
        \tfor _, para := range strings.Split(s, "\\n") {
        \t\tcur := para
        \t\tfor len(cur) > w {
        \t\t\tidx := strings.LastIndexByte(cur[:w], ' ')
        \t\t\tif idx <= w*3/5 {
        \t\t\t\tout = append(out, cur[:w])
        \t\t\t\trest := strings.TrimPrefix(cur[w:], " ")
        \t\t\t\tif indentCont {
        \t\t\t\t\tcur = " " + rest
        \t\t\t\t} else {
        \t\t\t\t\tcur = rest
        \t\t\t\t}
        \t\t\t} else {
        \t\t\t\tout = append(out, cur[:idx])
        \t\t\t\tif indentCont {
        \t\t\t\t\tcur = " " + cur[idx+1:]
        \t\t\t\t} else {
        \t\t\t\t\tcur = cur[idx+1:]
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t\tout = append(out, cur)
        \t}
        \treturn out
        }

        func args4jPrintUsage(
        \tp *Ext_org_kohsuke_args4j_CmdLineParser,
        \tw javaio.Writer,
        \tfilter Ext_org_kohsuke_args4j_OptionHandlerFilter,
        ) {
        \tif w == nil {
        \t\treturn
        \t}
        \tc := p.args4jCore()
        \topts := make([]args4jHandlerEntry, len(c.options))
        \tcopy(opts, c.options)
        \tsort.SliceStable(opts, func(i, j int) bool {
        \t\tni, nj := "", ""
        \t\tif ndi, ok := opts[i].def.(*Ext_org_kohsuke_args4j_NamedOptionDef); ok && ndi != nil {
        \t\t\tni = getOptionDefData(ndi.Core_).name
        \t\t}
        \t\tif ndj, ok := opts[j].def.(*Ext_org_kohsuke_args4j_NamedOptionDef); ok && ndj != nil {
        \t\t\tnj = getOptionDefData(ndj.Core_).name
        \t\t}
        \t\treturn ni < nj
        \t})
        \tmaxHead := 0
        \tfor _, ent := range opts {
        \t\tif head, _, ok := args4jOptionHead(ent); ok {
        \t\t\tif len(head) > maxHead {
        \t\t\t\tmaxHead = len(head)
        \t\t\t}
        \t\t}
        \t}
        \twidth := c.usageWidth
        \tif width <= 0 {
        \t\twidth = 80
        \t}
        \theadCol := maxHead
        \tif headCol > width/2-1 {
        \t\theadCol = width/2 - 1
        \t}
        \tusageCol := width - headCol - 3
        \tfor _, ent := range opts {
        \t\tif filter != nil && !filter.Select(ent.handler) {
        \t\t\tcontinue
        \t\t}
        \t\thead, d, ok := args4jOptionHead(ent)
        \t\tif !ok {
        \t\t\tcontinue
        \t\t}
        \t\theadLines := args4jWrapLines(head, headCol, true)
        \t\tusageLines := args4jWrapLines(d.usage+args4jDefaultSuffix(ent.handler), usageCol, false)
        \t\tnLines := len(headLines)
        \t\tif len(usageLines) > nLines {
        \t\t\tnLines = len(usageLines)
        \t\t}
        \t\tvar lb strings.Builder
        \t\tfor li := 0; li < nLines; li++ {
        \t\t\thl := ""
        \t\t\tif li < len(headLines) {
        \t\t\t\thl = headLines[li]
        \t\t\t}
        \t\t\tul := ""
        \t\t\tif li < len(usageLines) {
        \t\t\t\tul = usageLines[li]
        \t\t\t}
        \t\t\tlb.WriteString(hl)
        \t\t\tif li == 0 {
        \t\t\t\tif len(hl) < headCol {
        \t\t\t\t\tlb.WriteString(strings.Repeat(" ", headCol-len(hl)))
        \t\t\t\t}
        \t\t\t\tlb.WriteString(" : ")
        \t\t\t\tlb.WriteString(ul)
        \t\t\t} else if ul != "" {
        \t\t\t\tif len(hl) < headCol {
        \t\t\t\t\tlb.WriteString(strings.Repeat(" ", headCol-len(hl)))
        \t\t\t\t}
        \t\t\t\tlb.WriteString("   ")
        \t\t\t\tlb.WriteString(ul)
        \t\t\t}
        \t\t\tlb.WriteString("\\n")
        \t\t}
        \t\tw.Write(lb.String())
        \t}
        }

        func (x *args4jParserCore) IsJavaNil_() bool      { return x == nil }
        func (x *args4jHandlerEntry) IsJavaNil_() bool    { return x == nil }
        func (x *args4jOptionDefData) IsJavaNil_() bool   { return x == nil }
        func (x *args4jFieldSetterData) IsJavaNil_() bool { return x == nil }
        func (x *args4jParametersImpl) IsJavaNil_() bool  { return x == nil }
        func (x *args4jBuiltinHandler) IsJavaNil_() bool  { return x == nil }
        """;
  }
}
