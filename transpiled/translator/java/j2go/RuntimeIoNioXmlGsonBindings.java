package j2go;

import javax.lang.model.element.ExecutableElement;
import javax.lang.model.element.TypeElement;
import javax.lang.model.element.VariableElement;

/**
 * Clean-room runtime bindings for {@code java.io.*}, {@code java.nio.*}, {@code java.net.URI},
 * {@code java.util.zip.*}, {@code java.util.regex.*}, {@code java.text.*}, {@code java.util.logging.*},
 * {@code java.lang.management.*}, {@code com.google.common.{base,escape,hash,io,net,primitives}.*},
 * {@code com.google.gson.*}, {@code javax.xml.*}/{@code org.xml.sax.*}, and
 * {@code org.kohsuke.args4j.*}.
 */
public final class RuntimeIoNioXmlGsonBindings {

  private RuntimeIoNioXmlGsonBindings() {}

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
    if (qn.equals("java.io.File")) {
      if (fn.equals("separator") || fn.equals("Separator")) {
        out.append("var ").append(varGoName).append(" string = \"/\"\n");
        return true;
      }
      if (fn.equals("separatorChar") || fn.equals("SeparatorChar")) {
        out.append("var ").append(varGoName).append(" uint16 = uint16('/')\n");
        return true;
      }
      if (fn.equals("pathSeparator") || fn.equals("PathSeparator")) {
        out.append("var ").append(varGoName).append(" string = \":\"\n");
        return true;
      }
      if (fn.equals("pathSeparatorChar") || fn.equals("PathSeparatorChar")) {
        out.append("var ").append(varGoName).append(" uint16 = uint16(':')\n");
        return true;
      }
    }
    if (qn.equals("java.lang.System")) {
      if (fn.equals("out") || fn.equals("Out")) {
        out.append("var ")
            .append(varGoName)
            .append(' ')
            .append(fieldType)
            .append(" = wrapPrintStream(javaio.System_Out)\n");
        return true;
      }
      if (fn.equals("err") || fn.equals("Err")) {
        out.append("var ")
            .append(varGoName)
            .append(' ')
            .append(fieldType)
            .append(" = wrapPrintStream(javaio.System_Err)\n");
        return true;
      }
      if (fn.equals("in") || fn.equals("In")) {
        out.append("var ")
            .append(varGoName)
            .append(' ')
            .append(fieldType)
            .append(" = wrapStdinStream()\n");
        return true;
      }
    }
    if (qn.equals("java.nio.charset.StandardCharsets")) {
      if (fn.equals("UTF_8")) {
        out.append("var ")
            .append(varGoName)
            .append(' ')
            .append(fieldType)
            .append(" = wrapCharsetByName(\"UTF-8\")\n");
        return true;
      }
      if (fn.equals("US_ASCII")) {
        out.append("var ")
            .append(varGoName)
            .append(' ')
            .append(fieldType)
            .append(" = wrapCharsetByName(\"US-ASCII\")\n");
        return true;
      }
      if (fn.equals("ISO_8859_1")) {
        out.append("var ")
            .append(varGoName)
            .append(' ')
            .append(fieldType)
            .append(" = wrapCharsetByName(\"ISO-8859-1\")\n");
        return true;
      }
    }
    if (qn.equals("java.nio.charset.CodingErrorAction")) {
      out.append("var ")
          .append(varGoName)
          .append(' ')
          .append(fieldType)
          .append(" = newCodingErrorAction(")
          .append(GoNames.goStringLiteral(fn))
          .append(")\n");
      return true;
    }
    if (qn.equals("java.util.Locale")) {
      out.append("var ")
          .append(varGoName)
          .append(' ')
          .append(fieldType)
          .append(" = &Ext_java_util_Locale{Core_: ")
          .append(GoNames.goStringLiteral(fn))
          .append("}\n");
      return true;
    }
    if (qn.equals("java.util.logging.Level")) {
      out.append("var ")
          .append(varGoName)
          .append(' ')
          .append(fieldType)
          .append(" = &Ext_java_util_logging_Level{Core_: ")
          .append(GoNames.goStringLiteral(fn))
          .append("}\n");
      return true;
    }
    return false;
  }

  public static boolean tryEmitConstructorBody(
      RuntimeBindings rb,
      TypeElement owner,
      ExecutableElement ctor,
      PackagePlanner.ImportTracker imports,
      StringBuilder out) {
    String qn = owner.getQualifiedName().toString();
    int n = ctor.getParameters().size();
    String p0 = n >= 1 ? rb.goTypeInJ2gort(ctor.getParameters().get(0).asType(), imports) : "";

    switch (qn) {
      case "java.io.File" -> {
        if (n == 1 && p0.equals("string")) {
          out.append("\tobj.Core_ = javaio.NewFile(CheckNotNullString(arg0))\n");
          return true;
        }
        if (n == 2 && p0.contains("Ext_java_io_File")) {
          out.append("\tobj_f := unwrapFile(arg0)\n\tobj.Core_ = javaio.NewFileWithParent(obj_f, CheckNotNullString(arg1))\n");
          return true;
        }
        if (n == 2 && p0.equals("string")) {
          out.append("\tif arg0 == NullString { obj.Core_ = javaio.NewFile(CheckNotNullString(arg1)) } else { obj.Core_ = javaio.NewFileParentString(arg0, CheckNotNullString(arg1)) }\n");
          return true;
        }
      }
      case "java.io.FileInputStream" -> {
        if (n == 1 && p0.equals("string")) {
          out.append("\tfis, err := javaio.NewFileInputStream(CheckNotNullString(arg0))\n\tioErr(err)\n\tobj.Core_ = fis\n");
          return true;
        }
        if (n == 1 && p0.contains("Ext_java_io_File")) {
          out.append("\tfis, err := javaio.NewFileInputStreamFromFile(CheckNotNullPtr(unwrapFile(arg0)))\n\tioErr(err)\n\tobj.Core_ = fis\n");
          return true;
        }
      }
      case "java.io.FileOutputStream" -> {
        if (n == 1 && p0.equals("string")) {
          out.append("\tfos, err := javaio.NewFileOutputStream(CheckNotNullString(arg0))\n\tioErr(err)\n\tobj.Core_ = fos\n");
          return true;
        }
        if (n == 1 && p0.contains("Ext_java_io_File")) {
          out.append("\tfos, err := javaio.NewFileOutputStreamFromFile(CheckNotNullPtr(unwrapFile(arg0)))\n\tioErr(err)\n\tobj.Core_ = fos\n");
          return true;
        }
      }
      case "java.io.BufferedInputStream" -> {
        if (n >= 1) {
          out.append("\tCheckNotNull(arg0)\n\tobj.Core_ = javaio.NewBufferedInputStream(asInputStream(arg0))\n");
          return true;
        }
      }
      case "java.io.BufferedOutputStream" -> {
        if (n >= 1) {
          out.append("\tCheckNotNull(arg0)\n\tobj.Core_ = newBufferedExtOutputStream(arg0)\n");
          return true;
        }
      }
      case "java.io.ByteArrayOutputStream" -> {
        out.append("\tobj.Core_ = javaio.NewByteArrayOutputStream()\n");
        return true;
      }
      case "java.io.FilterInputStream", "java.io.FilterOutputStream" -> {
        if (n >= 1) {
          out.append("\tobj.Core_ = arg0\n");
          return true;
        }
      }
      case "java.io.InputStreamReader" -> {
        if (n == 1) {
          out.append("\tCheckNotNull(arg0)\n\tobj.Core_ = javaio.NewInputStreamReader(asInputStream(arg0))\n");
          return true;
        }
        if (n == 2) {
          String p1 = rb.goTypeInJ2gort(ctor.getParameters().get(1).asType(), imports);
          if (p1.equals("string")) {
            out.append("\tCheckNotNull(arg0)\n\tobj.Core_ = newInputStreamReaderCharsetName(arg0, CheckNotNullString(arg1))\n");
            return true;
          }
          out.append("\tCheckNotNull(arg0)\n\tobj.Core_ = newInputStreamReaderCharset(arg0, arg1)\n");
          return true;
        }
      }
      case "java.io.OutputStreamWriter" -> {
        if (n == 1) {
          out.append("\tCheckNotNull(arg0)\n\tobj.Core_ = javaio.NewOutputStreamWriter(asOutputStream(arg0))\n");
          return true;
        }
        if (n == 2) {
          String p1 = rb.goTypeInJ2gort(ctor.getParameters().get(1).asType(), imports);
          if (p1.contains("CharsetEncoder")) {
            out.append("\tCheckNotNull(arg0)\n\tobj.Core_ = newOutputStreamWriterEncoder(arg0, arg1)\n");
            return true;
          }
          if (p1.equals("string")) {
            out.append("\tCheckNotNull(arg0)\n\tobj.Core_ = newOutputStreamWriterCharsetName(arg0, CheckNotNullString(arg1))\n");
            return true;
          }
          out.append("\tCheckNotNull(arg0)\n\tobj.Core_ = newOutputStreamWriterCharset(arg0, arg1)\n");
          return true;
        }
      }
      case "java.io.BufferedReader" -> {
        if (n >= 1) {
          out.append("\tCheckNotNull(arg0)\n\tobj.Core_ = javaio.NewBufferedReader(asReader(arg0))\n");
          return true;
        }
      }
      case "java.io.BufferedWriter" -> {
        if (n >= 1) {
          out.append("\tCheckNotNull(arg0)\n\tobj.Core_ = javaio.NewBufferedWriter(asWriter(arg0))\n");
          return true;
        }
      }
      case "java.io.StringReader" -> {
        if (n == 1) {
          out.append("\tobj.Core_ = javaio.NewStringReader(CheckNotNullString(arg0))\n");
          return true;
        }
      }
      case "java.io.StringWriter" -> {
        out.append("\tobj.Core_ = javaio.NewStringWriter()\n");
        return true;
      }
      case "java.io.PrintWriter" -> {
        if (n >= 1 && p0.contains("Writer")) {
          out.append("\tCheckNotNull(arg0)\n\tobj.Core_ = asWriter(arg0)\n");
          return true;
        }
        if (n >= 1 && p0.contains("OutputStream")) {
          out.append("\tCheckNotNull(arg0)\n\tobj.Core_ = javaio.NewOutputStreamWriter(asOutputStream(arg0))\n");
          return true;
        }
      }
      case "java.io.PrintStream" -> {
        if (n >= 1 && p0.equals("string")) {
          out.append("\tfos, err := javaio.NewFileOutputStream(CheckNotNullString(arg0))\n\tioErr(err)\n\tobj.Core_ = javaio.NewPrintStream(fos)\n");
          return true;
        }
        if (n >= 1 && p0.contains("Ext_java_io_File")) {
          out.append("\tfos, err := javaio.NewFileOutputStreamFromFile(CheckNotNullPtr(unwrapFile(arg0)))\n\tioErr(err)\n\tobj.Core_ = javaio.NewPrintStream(fos)\n");
          return true;
        }
        if (n >= 1 && p0.contains("OutputStream")) {
          out.append("\tCheckNotNull(arg0)\n\tobj.Core_ = javaio.NewPrintStream(asOutputStream(arg0))\n");
          return true;
        }
      }
      case "java.util.zip.GZIPInputStream" -> {
        if (n >= 1) {
          out.append("\tCheckNotNull(arg0)\n\tobj.Core_ = newGzipInputStream(arg0)\n");
          return true;
        }
      }
      case "java.util.zip.GZIPOutputStream" -> {
        if (n >= 1) {
          out.append("\tinitGzipOutputStream(&obj.Ext_java_util_zip_DeflaterOutputStreamBase, arg0)\n");
          return true;
        }
      }
      case "java.util.zip.ZipFile" -> {
        if (n >= 1 && p0.contains("Ext_java_io_File")) {
          out.append("\tobj.Core_ = newZipFileFromFile(arg0)\n");
          return true;
        }
        if (n >= 1 && p0.equals("string")) {
          out.append("\tobj.Core_ = newZipFileFromName(CheckNotNullString(arg0))\n");
          return true;
        }
      }
      case "java.util.zip.ZipInputStream" -> {
        if (n == 1) {
          out.append("\tCheckNotNull(arg0)\n\tobj.Core_ = newZipInputStream(arg0, nil)\n");
          return true;
        }
        if (n == 2) {
          out.append("\tCheckNotNull(arg0)\n\tobj.Core_ = newZipInputStream(arg0, arg1)\n");
          return true;
        }
      }
      case "java.net.URI" -> {
        if (n == 1 && p0.equals("string")) {
          out.append("\tobj.Core_ = newJavaURI(CheckNotNullString(arg0))\n");
          return true;
        }
      }
      case "java.nio.file.SimpleFileVisitor" -> {
        return true;
      }
      case "com.google.common.escape.Escaper", "com.google.common.escape.UnicodeEscaper" -> {
        return true;
      }
      case "com.google.common.net.PercentEscaper" -> {
        if (n == 2) {
          out.append("\tobj.Core_ = newPercentEscaper(CheckNotNullString(arg0), arg1)\n");
          return true;
        }
      }
      case "com.google.gson.Gson" -> {
        out.append("\tobj.Core_ = newDefaultGson()\n");
        return true;
      }
      case "com.google.gson.GsonBuilder" -> {
        out.append("\tobj.Core_ = newGsonBuilderCore()\n");
        return true;
      }
      case "com.google.gson.JsonArray" -> {
        out.append("\tobj.Core_ = newGsonJsonArray()\n");
        return true;
      }
      case "com.google.gson.JsonObject" -> {
        out.append("\tobj.Core_ = newGsonJsonObject()\n");
        return true;
      }
      case "com.google.gson.JsonPrimitive" -> {
        if (n == 1) {
          out.append("\tobj.Core_ = newGsonJsonPrimitive(arg0)\n");
          return true;
        }
      }
      case "com.google.gson.TypeAdapter", "com.google.gson.reflect.TypeToken" -> {
        return true;
      }
      case "com.google.gson.stream.JsonReader" -> {
        if (n == 1) {
          out.append("\tobj.Core_ = newGsonJsonReader(arg0)\n");
          return true;
        }
      }
      case "com.google.gson.stream.JsonWriter" -> {
        if (n == 1) {
          out.append("\tobj.Core_ = newGsonJsonWriter(arg0)\n");
          return true;
        }
      }
      case "org.xml.sax.InputSource" -> {
        if (n == 1) {
          out.append("\tobj.Core_ = arg0\n");
          return true;
        }
      }
      case "org.kohsuke.args4j.CmdLineParser" -> {
        if (n >= 1) {
          out.append("\tinitArgs4jParser(obj, arg0)\n");
          return true;
        }
      }
      case "org.kohsuke.args4j.spi.OptionHandler",
          "org.kohsuke.args4j.spi.OneArgumentOptionHandler",
          "org.kohsuke.args4j.spi.StringOptionHandler",
          "org.kohsuke.args4j.spi.IntOptionHandler" -> {
        if (n == 3) {
          out.append("\tobj.Core_ = arg0\n\tobj.F_Option = arg1\n\tobj.F_Setter = arg2\n");
          return true;
        }
      }
      case "java.lang.Enum" -> {
        if (n == 2) {
          out.append("\tobj.Core_ = &EnumBase{Name_: arg0, Ordinal_: arg1}\n");
          return true;
        }
      }
      case "java.lang.Record" -> {
        return true;
      }
      case "java.lang.Thread" -> {
        if (n == 4) {
          out.append("\tobj.Core_ = arg1\n");
          return true;
        }
        if (n == 1 && p0.contains("Runnable")) {
          out.append("\tobj.Core_ = arg0\n");
          return true;
        }
        return true;
      }
      case "java.lang.ThreadLocal" -> {
        out.append("\tobj.Core_ = &jThreadLocal{}\n");
        return true;
      }
      case "java.lang.management.MemoryUsage" -> {
        if (n == 4) {
          out.append("\tobj.Core_ = &jMemoryUsage{init: arg0, used: arg1, committed: arg2, max: arg3}\n");
          return true;
        }
      }
      case "java.util.StringTokenizer" -> {
        if (n == 2) {
          out.append("\tobj.Core_ = newStringTokenizer(CheckNotNullString(arg0), CheckNotNullString(arg1))\n");
          return true;
        }
        if (n == 1) {
          out.append("\tobj.Core_ = newStringTokenizer(CheckNotNullString(arg0), \" \\t\\n\\r\\f\")\n");
          return true;
        }
      }
      default -> {}
    }
    return false;
  }

  public static boolean tryEmitInitBody(
      RuntimeBindings rb,
      TypeElement owner,
      ExecutableElement ctor,
      PackagePlanner.ImportTracker imports,
      StringBuilder out) {
    String qn = owner.getQualifiedName().toString();
    int n = ctor.getParameters().size();
    switch (qn) {
      case "java.io.FilterInputStream", "java.io.FilterOutputStream" -> {
        if (n >= 1) {
          out.append("\tthis_.Core_ = arg0\n");
          String fieldName = qn.endsWith("InputStream") ? "in" : "out";
          for (var ve : rb.resolver().getReferencedExternalFields()) {
            if (ve.getEnclosingElement().equals(owner)
                && ve.getSimpleName().toString().equals(fieldName)) {
              out.append("\tthis_.F_")
                  .append(GoNames.exported(fieldName))
                  .append(" = arg0\n");
              break;
            }
          }
          return true;
        }
      }
      case "java.util.zip.GZIPOutputStream" -> {
        if (n >= 1) {
          out.append("\tinitGzipOutputStream(&this_.Ext_java_util_zip_DeflaterOutputStreamBase, arg0)\n");
          return true;
        }
      }
      case "org.kohsuke.args4j.CmdLineException" -> {
        if (n == 2) {
          String p1 = rb.goTypeInJ2gort(ctor.getParameters().get(1).asType(), imports);
          if (p1.equals("string")) {
            out.append("\t_ = arg0\n\tjavalang.InitThrowable(this_, \"org.kohsuke.args4j.CmdLineException\", arg1, nil)\n");
          } else {
            out.append("\t_ = arg0\n\tmsg := \"\"\n\tif arg1 != nil { msg = arg1.GetMessage() }\n\tjavalang.InitThrowable(this_, \"org.kohsuke.args4j.CmdLineException\", msg, arg1)\n");
          }
          return true;
        }
      }
      case "org.kohsuke.args4j.spi.OptionHandler",
          "org.kohsuke.args4j.spi.OneArgumentOptionHandler",
          "org.kohsuke.args4j.spi.StringOptionHandler",
          "org.kohsuke.args4j.spi.IntOptionHandler" -> {
        if (n == 3) {
          out.append("\tthis_.Core_ = arg0\n\tthis_.F_Option = arg1\n\tthis_.F_Setter = arg2\n");
          return true;
        }
      }
      case "java.lang.Enum" -> {
        if (n == 2) {
          out.append("\tthis_.Core_ = &EnumBase{Name_: arg0, Ordinal_: arg1}\n");
          return true;
        }
      }
      case "java.lang.Record",
          "java.nio.file.SimpleFileVisitor",
          "com.google.common.escape.Escaper",
          "com.google.common.escape.UnicodeEscaper",
          "com.google.gson.TypeAdapter",
          "com.google.gson.reflect.TypeToken" -> {
        return true;
      }
      default -> {}
    }
    if (rb.resolver().isThrowable(owner)) {
      out.append("\tInitUserThrowable(this_, ").append(GoNames.goStringLiteral(qn));
      for (int i = 0; i < n; i++) {
        out.append(", arg").append(i);
      }
      out.append(")\n");
      return true;
    }
    return false;
  }

  public static boolean tryEmitStaticMethodBody(
      RuntimeBindings rb,
      TypeElement owner,
      ExecutableElement m,
      String retType,
      PackagePlanner.ImportTracker imports,
      StringBuilder out) {
    String qn = owner.getQualifiedName().toString();
    String simple = m.getSimpleName().toString();
    int arity = m.getParameters().size();
    String p0 = arity >= 1 ? rb.goTypeInJ2gort(m.getParameters().get(0).asType(), imports) : "";

    switch (qn) {
      case "com.google.common.base.Ascii" -> {
        if (simple.equals("equalsIgnoreCase") && arity == 2) {
          out.append("\treturn guavaAsciiEqualsIgnoreCase(arg0, arg1)\n");
          return true;
        }
        if (simple.equals("isUpperCase") && arity == 1) {
          out.append("\treturn arg0 >= 'A' && arg0 <= 'Z'\n");
          return true;
        }
        if (simple.equals("isLowerCase") && arity == 1) {
          out.append("\treturn arg0 >= 'a' && arg0 <= 'z'\n");
          return true;
        }
        if (simple.equals("toLowerCase") && arity == 1 && p0.equals("string")) {
          out.append("\treturn javalang.StringToLowerCaseASCII(CheckNotNullString(arg0))\n");
          return true;
        }
        if (simple.equals("toUpperCase") && arity == 1 && p0.equals("string")) {
          out.append("\treturn javalang.StringToUpperCaseASCII(CheckNotNullString(arg0))\n");
          return true;
        }
      }
      case "com.google.common.base.CharMatcher" -> {
        if (simple.equals("anyOf") && arity == 1) {
          out.append("\treturn wrapCharMatcherAnyOf(arg0)\n");
          return true;
        }
        if (simple.equals("inRange") && arity == 2) {
          out.append("\treturn wrapCharMatcherInRange(arg0, arg1)\n");
          return true;
        }
        if (simple.equals("is") && arity == 1) {
          out.append("\treturn wrapCharMatcherIs(arg0)\n");
          return true;
        }
        if (simple.equals("isNot") && arity == 1) {
          out.append("\treturn wrapCharMatcherIsNot(arg0)\n");
          return true;
        }
        if (simple.equals("javaLetterOrDigit") && arity == 0) {
          out.append("\treturn wrapCharMatcherJavaLetterOrDigit()\n");
          return true;
        }
        if (simple.equals("whitespace") && arity == 0) {
          out.append("\treturn wrapCharMatcherWhitespace()\n");
          return true;
        }
        if (simple.equals("ascii") && arity == 0) {
          out.append("\treturn wrapCharMatcherAscii()\n");
          return true;
        }
      }
      case "com.google.common.base.Functions" -> {
        if (simple.equals("identity") && arity == 0) {
          out.append("\treturn newGuavaIdentityFunction()\n");
          return true;
        }
      }
      case "com.google.common.base.Predicates" -> {
        if (simple.equals("alwaysTrue") && arity == 0) {
          out.append("\treturn guavaAlwaysTruePred_\n");
          return true;
        }
        if (simple.equals("alwaysFalse") && arity == 0) {
          out.append("\treturn guavaAlwaysFalsePred_\n");
          return true;
        }
        if (simple.equals("not") && arity == 1) {
          out.append("\tCheckNotNull(arg0)\n\treturn newGuavaNotPredicate(arg0)\n");
          return true;
        }
        if (simple.equals("equalTo") && arity == 1) {
          out.append("\ttarget := ToAny(arg0)\n\treturn newGuavaPredicate(func(x any) bool { return javalang.ObjectsEquals(ToAny(x), target) })\n");
          return true;
        }
      }
      case "com.google.common.base.Suppliers" -> {
        if (simple.equals("memoize") && arity == 1) {
          out.append("\treturn newMemoizingSupplier(CheckNotNull(arg0))\n");
          return true;
        }
      }
      case "com.google.common.base.Joiner" -> {
        if (simple.equals("on") && arity == 1) {
          if (p0.equals("uint16")) {
            out.append("\treturn newGuavaJoiner(javalang.StringFromChar(arg0))\n");
          } else {
            out.append("\treturn newGuavaJoiner(CheckNotNullString(arg0))\n");
          }
          return true;
        }
      }
      case "com.google.common.base.Splitter" -> {
        if (simple.equals("on") && arity == 1) {
          if (p0.equals("uint16")) {
            out.append("\treturn newGuavaSplitterChar(arg0)\n");
          } else {
            out.append("\treturn newGuavaSplitterStr(CheckNotNullString(arg0))\n");
          }
          return true;
        }
      }
      case "com.google.common.base.MoreObjects" -> {
        if (simple.equals("firstNonNull") && arity == 2) {
          out.append("\tif a := ToAny(arg0); a != nil { return a }\n\treturn CheckNotNull(ToAny(arg1))\n");
          return true;
        }
        if (simple.equals("toStringHelper") && arity == 1) {
          if (p0.contains("Class")) {
            out.append("\treturn newToStringHelper(CheckNotNullPtr(arg0).GetSimpleName())\n");
          } else {
            out.append("\treturn newToStringHelper(GetObjectClass(arg0).GetSimpleName())\n");
          }
          return true;
        }
      }
      case "com.google.common.base.Preconditions" -> {
        if (simple.equals("checkArgument")) {
          if (arity == 1) {
            out.append("\tCheckArgument(arg0)\n");
          } else if (arity == 2) {
            out.append("\tCheckArgumentMsg(arg0, arg1)\n");
          } else if (arity == 3
              && rb.goTypeInJ2gort(m.getParameters().get(2).asType(), imports).equals("[]any")) {
            out.append("\tCheckArgumentFmt(arg0, arg1, arg2...)\n");
          } else {
            out.append("\tCheckArgumentFmt(arg0, arg1");
            for (int i = 2; i < arity; i++) out.append(", arg").append(i);
            out.append(")\n");
          }
          return true;
        }
        if (simple.equals("checkState")) {
          if (arity == 1) {
            out.append("\tCheckState(arg0)\n");
          } else if (arity == 2) {
            out.append("\tCheckStateMsg(arg0, arg1)\n");
          } else if (arity == 3
              && rb.goTypeInJ2gort(m.getParameters().get(2).asType(), imports).equals("[]any")) {
            out.append("\tCheckStateFmt(arg0, arg1, arg2...)\n");
          } else {
            out.append("\tCheckStateFmt(arg0, arg1");
            for (int i = 2; i < arity; i++) out.append(", arg").append(i);
            out.append(")\n");
          }
          return true;
        }
        if (simple.equals("checkNotNull")) {
          if (arity == 1) {
            out.append("\treturn CheckNotNull(arg0)\n");
          } else if (arity == 2) {
            out.append("\tif arg0 == nil || javalang.IsNil(arg0) { panic(javalang.NewNullPointerExceptionMsg(StrVal(arg1))) }\n\tif s, ok := arg0.(string); ok && s == NullString { panic(javalang.NewNullPointerExceptionMsg(StrVal(arg1))) }\n\treturn arg0\n");
          } else {
            out.append("\tif arg0 == nil || javalang.IsNil(arg0) || (func() bool { s, ok := arg0.(string); return ok && s == NullString }()) {\n\t\tpanic(javalang.NewNullPointerExceptionMsg(LenientFormat(arg1");
            for (int i = 2; i < arity; i++) out.append(", arg").append(i);
            out.append(")))\n\t}\n\treturn arg0\n");
          }
          return true;
        }
      }
      case "com.google.common.base.Strings" -> {
        if (simple.equals("isNullOrEmpty") && arity == 1) {
          out.append("\treturn arg0 == NullString || len(arg0) == 0\n");
          return true;
        }
        if (simple.equals("nullToEmpty") && arity == 1) {
          out.append("\tif arg0 == NullString { return \"\" }\n\treturn arg0\n");
          return true;
        }
        if (simple.equals("emptyToNull") && arity == 1) {
          out.append("\tif arg0 == NullString || len(arg0) == 0 { return NullString }\n\treturn arg0\n");
          return true;
        }
        if (simple.equals("lenientFormat") && arity == 2) {
          out.append("\tif arg1 == nil { return LenientFormat(arg0, nil) }\n\treturn LenientFormat(arg0, arg1...)\n");
          return true;
        }
        if (simple.equals("padStart") && arity == 3) {
          out.append("\treturn guavaPadStart(CheckNotNullString(arg0), arg1, arg2)\n");
          return true;
        }
      }
      case "com.google.common.base.Throwables" -> {
        if (simple.equals("throwIfUnchecked") && arity == 1) {
          out.append("\tCheckNotNull(arg0)\n\tif _, ok := arg0.(javalang.IRuntimeException); ok { panic(arg0) }\n\tif _, ok := arg0.(javalang.IError); ok { panic(arg0) }\n");
          return true;
        }
      }
      case "com.google.common.hash.Hashing" -> {
        if (simple.equals("farmHashFingerprint64") && arity == 0) {
          out.append("\treturn newFarmHashFunction()\n");
          return true;
        }
        if (simple.equals("murmur3_128") && arity == 0) {
          out.append("\treturn newMurmur3_128HashFunction()\n");
          return true;
        }
      }
      case "com.google.common.io.BaseEncoding" -> {
        if (simple.equals("base64") && arity == 0) {
          out.append("\treturn newGuavaBase64Encoding()\n");
          return true;
        }
      }
      case "com.google.common.io.ByteStreams" -> {
        if (simple.equals("limit") && arity == 2) {
          out.append("\treturn limitExtInputStream(CheckNotNull(arg0), arg1)\n");
          return true;
        }
        if (simple.equals("toByteArray") && arity == 1) {
          out.append("\treturn readAllExtInputStreamBytes(CheckNotNull(arg0))\n");
          return true;
        }
      }
      case "com.google.common.io.CharSource" -> {
        if (simple.equals("wrap") && arity == 1) {
          out.append("\treturn wrapStringCharSource(StrVal(CheckNotNull(arg0)))\n");
          return true;
        }
      }
      case "com.google.common.io.CharStreams" -> {
        if (simple.equals("toString") && arity == 1) {
          out.append("\treturn readAllFromReadable(CheckNotNull(arg0))\n");
          return true;
        }
        if (simple.equals("readLines") && arity == 1) {
          out.append("\treturn readLinesFromReadable(CheckNotNull(arg0))\n");
          return true;
        }
      }
      case "com.google.common.io.Files" -> {
        if (simple.equals("asCharSink") && arity == 3) {
          out.append("\treturn newGuavaFileCharSink(CheckNotNullPtr(arg0), arg1, arg2)\n");
          return true;
        }
        if (simple.equals("asCharSource") && arity == 2) {
          out.append("\treturn newGuavaFileCharSource(CheckNotNullPtr(arg0), arg1)\n");
          return true;
        }
        if (simple.equals("newReader") && arity == 2) {
          out.append("\treturn newGuavaFileBufferedReader(CheckNotNullPtr(arg0), arg1)\n");
          return true;
        }
        if (simple.equals("toByteArray") && arity == 1) {
          out.append("\treturn guavaFileToByteArray(CheckNotNullPtr(arg0))\n");
          return true;
        }
        if (simple.equals("write") && arity == 2) {
          out.append("\tguavaFileWriteBytes(CheckNotNullSlice(arg0), CheckNotNullPtr(arg1))\n");
          return true;
        }
      }
      case "com.google.common.primitives.Chars" -> {
        if (simple.equals("asList") && arity == 1) {
          out.append("\tCheckNotNullSlice(arg0)\n\titems := make([]any, len(arg0))\n\tfor i, c := range arg0 { items[i] = ValueOfChar(c) }\n\treturn wrapArrayList(items)\n");
          return true;
        }
        if (simple.equals("toArray") && arity == 1) {
          out.append("\tCheckNotNull(arg0)\n\tvar res []uint16\n\tfor it := Iterate(arg0); it.HasNext(); {\n\t\tres = append(res, UnboxChar(Cast[*BoxChar](it.Next())))\n\t}\n\treturn res\n");
          return true;
        }
      }
      case "java.lang.Runtime" -> {
        if (simple.equals("getRuntime") && arity == 0) {
          out.append("\treturn &Ext_java_lang_Runtime{}\n");
          return true;
        }
      }
      case "java.lang.String" -> {
        if (simple.equals("format") && arity == 2) {
          out.append("\treturn StringFormat(CheckNotNullString(arg0), arg1...)\n");
          return true;
        }
        if (simple.equals("format") && arity == 3) {
          out.append("\t_ = arg0\n\treturn StringFormat(CheckNotNullString(arg1), arg2...)\n");
          return true;
        }
        if (simple.equals("join") && arity == 2) {
          if (rb.goTypeInJ2gort(m.getParameters().get(1).asType(), imports).startsWith("[]")) {
            out.append("\tdelim := StrVal(CheckNotNull(arg0))\n\tCheckNotNullSlice(arg1)\n\tparts := make([]string, len(arg1))\n\tfor i, x := range arg1 { parts[i] = StrVal(x) }\n\treturn strings.Join(parts, delim)\n");
          } else {
            out.append("\tdelim := StrVal(CheckNotNull(arg0))\n\tCheckNotNull(arg1)\n\tvar parts []string\n\tfor it := Iterate(arg1); it.HasNext(); { parts = append(parts, StrVal(it.Next())) }\n\treturn strings.Join(parts, delim)\n");
          }
          return true;
        }
        if (simple.equals("valueOf") && arity == 1) {
          switch (p0) {
            case "int32", "int64" -> {
              out.append("\treturn strconv.FormatInt(int64(arg0), 10)\n");
              return true;
            }
            case "float64" -> {
              out.append("\treturn javalang.DoubleToString(arg0)\n");
              return true;
            }
            case "float32" -> {
              out.append("\treturn javalang.FloatToString(arg0)\n");
              return true;
            }
            case "bool" -> {
              out.append("\treturn BoolToStr(arg0)\n");
              return true;
            }
            case "string" -> {
              out.append("\treturn StrValString(arg0)\n");
              return true;
            }
            case "uint16" -> {
              out.append("\treturn javalang.StringFromChar(arg0)\n");
              return true;
            }
            case "[]uint16" -> {
              out.append("\treturn javalang.StringFromUTF16(CheckNotNullSlice(arg0))\n");
              return true;
            }
            default -> {
              out.append("\treturn StrVal(arg0)\n");
              return true;
            }
          }
        }
        if (simple.equals("valueOf") && arity == 3) {
          out.append("\tCheckNotNullSlice(arg0)\n\tif arg1 < 0 || arg2 < 0 || int(arg1)+int(arg2) > len(arg0) {\n\t\tpanic(javalang.NewStringIndexOutOfBoundsExceptionForRange(int(arg1), int(arg1+arg2), len(arg0)))\n\t}\n\treturn javalang.StringFromUTF16(arg0[arg1 : arg1+arg2])\n");
          return true;
        }
      }
      case "java.lang.System" -> {
        if (simple.equals("arraycopy") && arity == 5) {
          out.append("\tArraycopy(arg0, arg1, arg2, arg3, arg4)\n");
          return true;
        }
        if (simple.equals("exit") && arity == 1) {
          out.append("\tjavaio.SystemExit(arg0)\n");
          return true;
        }
        if (simple.equals("getProperty") && arity == 1) {
          out.append("\treturn javaSystemGetProperty(CheckNotNullString(arg0), NullString)\n");
          return true;
        }
        if (simple.equals("getProperty") && arity == 2) {
          out.append("\treturn javaSystemGetProperty(CheckNotNullString(arg0), arg1)\n");
          return true;
        }
        if (simple.equals("getenv") && arity == 1) {
          out.append("\treturn javaSystemGetenv(CheckNotNullString(arg0))\n");
          return true;
        }
      }
      case "java.lang.Thread" -> {
        if (simple.equals("currentThread") && arity == 0) {
          out.append("\treturn currentJavaThread()\n");
          return true;
        }
        if (simple.equals("interrupted") && arity == 0) {
          out.append("\treturn false\n");
          return true;
        }
      }
      case "java.lang.management.ManagementFactory" -> {
        if (simple.equals("getCompilationMXBean") && arity == 0) {
          out.append("\treturn &jCompilationMXBean{}\n");
          return true;
        }
        if (simple.equals("getGarbageCollectorMXBeans") && arity == 0) {
          out.append("\treturn wrapArrayList([]any{&jGarbageCollectorMXBean{}})\n");
          return true;
        }
        if (simple.equals("getMemoryMXBean") && arity == 0) {
          out.append("\treturn &jMemoryMXBean{}\n");
          return true;
        }
        if (simple.equals("getMemoryPoolMXBeans") && arity == 0) {
          out.append("\treturn wrapArrayList([]any{&jMemoryPoolMXBean{}})\n");
          return true;
        }
      }
      case "java.nio.charset.Charset" -> {
        if (simple.equals("forName") && arity == 1) {
          out.append("\treturn charsetForName(CheckNotNullString(arg0))\n");
          return true;
        }
        if (simple.equals("isSupported") && arity == 1) {
          out.append("\treturn charsetIsSupported(CheckNotNullString(arg0))\n");
          return true;
        }
        if (simple.equals("defaultCharset") && arity == 0) {
          out.append("\treturn wrapCharsetByName(\"UTF-8\")\n");
          return true;
        }
      }
      case "java.nio.file.FileSystems" -> {
        if (simple.equals("getDefault") && arity == 0) {
          out.append("\treturn newDefaultFileSystem()\n");
          return true;
        }
      }
      case "java.nio.file.Path", "java.nio.file.Paths" -> {
        if ((simple.equals("of") || simple.equals("get")) && arity == 2 && p0.equals("string")) {
          out.append("\treturn newExtPath(CheckNotNullString(arg0), arg1)\n");
          return true;
        }
        if ((simple.equals("of") || simple.equals("get")) && arity == 1 && p0.equals("string")) {
          out.append("\treturn newExtPath(CheckNotNullString(arg0), nil)\n");
          return true;
        }
      }
      case "java.nio.file.Files" -> {
        if (simple.equals("createDirectories") && arity == 2) {
          out.append("\treturn nioFilesCreateDirectories(CheckNotNull(arg0))\n");
          return true;
        }
        if (simple.equals("getLastModifiedTime") && arity == 2) {
          out.append("\treturn nioFilesGetLastModifiedTime(CheckNotNull(arg0))\n");
          return true;
        }
        if (simple.equals("newBufferedReader") && arity == 2) {
          out.append("\treturn nioFilesNewBufferedReader(CheckNotNull(arg0), arg1)\n");
          return true;
        }
        if (simple.equals("newBufferedWriter") && arity == 3) {
          out.append("\treturn nioFilesNewBufferedWriter(CheckNotNull(arg0), arg1, arg2)\n");
          return true;
        }
        if (simple.equals("newOutputStream") && arity == 2) {
          out.append("\treturn nioFilesNewOutputStream(CheckNotNull(arg0), arg1)\n");
          return true;
        }
        if (simple.equals("readString") && arity == 2) {
          out.append("\treturn nioFilesReadString(CheckNotNull(arg0), arg1)\n");
          return true;
        }
        if (simple.equals("readString") && arity == 1) {
          out.append("\treturn nioFilesReadString(CheckNotNull(arg0), nil)\n");
          return true;
        }
        if (simple.equals("walkFileTree") && arity == 4) {
          out.append("\treturn nioFilesWalkFileTree(CheckNotNull(arg0), arg1, arg2, CheckNotNull(arg3))\n");
          return true;
        }
      }
      case "java.text.MessageFormat" -> {
        if (simple.equals("format") && arity == 2) {
          out.append("\treturn javaMessageFormat(CheckNotNullString(arg0), arg1)\n");
          return true;
        }
      }
      case "java.util.Objects" -> {
        if (simple.equals("equals") && arity == 2) {
          out.append("\treturn javalang.ObjectsEquals(ToAny(arg0), ToAny(arg1))\n");
          return true;
        }
        if (simple.equals("requireNonNull") && arity == 1) {
          out.append("\treturn CheckNotNull(arg0)\n");
          return true;
        }
        if (simple.equals("requireNonNull") && arity == 2) {
          out.append("\tif arg0 == nil || javalang.IsNil(arg0) || (func() bool { s, ok := arg0.(string); return ok && s == NullString }()) {\n\t\tpanic(javalang.NewNullPointerExceptionMsg(StrVal(arg1)))\n\t}\n\treturn arg0\n");
          return true;
        }
        if (simple.equals("hashCode") && arity == 1) {
          out.append("\treturn javalang.ObjectHashCode(ToAny(arg0))\n");
          return true;
        }
        if (simple.equals("hash") && arity == 1) {
          out.append("\treturn ObjectsHash(arg0...)\n");
          return true;
        }
      }
      case "java.util.ResourceBundle" -> {
        if (simple.equals("getBundle") && arity >= 1) {
          out.append("\treturn wrapResourceBundle(CheckNotNullString(arg0))\n");
          return true;
        }
      }
      case "java.util.logging.Level" -> {
        if (simple.equals("parse") && arity == 1) {
          out.append("\treturn &Ext_java_util_logging_Level{Core_: CheckNotNullString(arg0)}\n");
          return true;
        }
      }
      case "java.util.logging.Logger" -> {
        if (simple.equals("getLogger") && arity == 1) {
          out.append("\treturn getJavaLogger(CheckNotNullString(arg0))\n");
          return true;
        }
      }
      case "java.util.regex.Pattern" -> {
        if (simple.equals("compile") && arity == 1) {
          out.append("\treturn wrapPatternCompile(CheckNotNullString(arg0), 0)\n");
          return true;
        }
        if (simple.equals("compile") && arity == 2) {
          out.append("\treturn wrapPatternCompile(CheckNotNullString(arg0), int(arg1))\n");
          return true;
        }
        if (simple.equals("quote") && arity == 1) {
          out.append("\treturn javaregex.Pattern_Quote(CheckNotNullString(arg0))\n");
          return true;
        }
        if (simple.equals("matches") && arity == 2) {
          out.append("\treturn javaregex.Pattern_Matches(CheckNotNullString(arg0), StrVal(CheckNotNull(arg1)))\n");
          return true;
        }
      }
      case "java.util.regex.Matcher" -> {
        if (simple.equals("quoteReplacement") && arity == 1) {
          out.append("\treturn javaregex.Matcher_QuoteReplacement(CheckNotNullString(arg0))\n");
          return true;
        }
      }
      case "javax.xml.parsers.SAXParserFactory" -> {
        if (simple.equals("newInstance") && arity == 0) {
          out.append("\tf := &Ext_javax_xml_parsers_SAXParserFactoryBase{}\n\tf.BindExtSelf_(f)\n\treturn f\n");
          return true;
        }
      }
      default -> {}
    }
    return false;
  }

  public static String tryBuildInstanceMethodBody(
      RuntimeBindings rb,
      TypeElement te,
      TypeResolver.MethodSig sig,
      String retType,
      PackagePlanner.ImportTracker imports) {
    String qn = te.getQualifiedName().toString();
    String m = sig.simpleName;
    int arity = sig.canonParamTypes.size();
    String p0 = arity >= 1 ? rb.goTypeInJ2gort(sig.canonParamTypes.get(0), imports) : "";

    // StringBuilder / AbstractStringBuilder / StringBuffer
    if (qn.equals("java.lang.AbstractStringBuilder")
        || qn.equals("java.lang.StringBuilder")
        || qn.equals("java.lang.StringBuffer")) {
      String retSelf =
          qn.equals("java.lang.AbstractStringBuilder")
              ? "\tif this_.Self_ != nil { return this_.Self_ }\n\treturn this_\n"
              : "\treturn this_\n";
      if (m.equals("append") && arity == 1) {
        if (p0.equals("string")) {
          return "\tthis_.sbCore().appendString(arg0)\n" + retSelf;
        }
        if (p0.equals("uint16")) {
          return "\tthis_.sbCore().appendChar(arg0)\n" + retSelf;
        }
        if (p0.equals("[]uint16")) {
          return "\tthis_.sbCore().appendString(javalang.StringFromUTF16(CheckNotNullSlice(arg0)))\n" + retSelf;
        }
        if (p0.equals("int32")) {
          return "\tthis_.sbCore().appendInt32(arg0)\n" + retSelf;
        }
        if (p0.equals("int64")) {
          return "\tthis_.sbCore().appendInt64(arg0)\n" + retSelf;
        }
        if (p0.equals("bool")) {
          return "\tthis_.sbCore().appendBool(arg0)\n" + retSelf;
        }
        if (p0.equals("float64")) {
          return "\tthis_.sbCore().appendASCII(javalang.DoubleToString(arg0))\n" + retSelf;
        }
        if (p0.equals("float32")) {
          return "\tthis_.sbCore().appendASCII(javalang.FloatToString(arg0))\n" + retSelf;
        }
        return "\tthis_.sbCore().appendString(StrVal(arg0))\n" + retSelf;
      }
      if (m.equals("append") && arity == 3) {
        if (p0.equals("[]uint16")) {
          return "\tCheckNotNullSlice(arg0)\n\tthis_.sbCore().appendString(javalang.StringFromUTF16(arg0[arg1 : arg1+arg2]))\n" + retSelf;
        }
        String sExpr = p0.equals("string") ? "StrValString(arg0)" : "StrVal(arg0)";
        return "\ts := " + sExpr + "\n\tthis_.sbCore().appendString(JStrSub(s, arg1, arg2))\n" + retSelf;
      }
      if (m.equals("appendCodePoint") && arity == 1) {
        return "\tthis_.sbCore().appendCodePoint(arg0)\n" + retSelf;
      }
      if (m.equals("insert") && arity == 2) {
        String p1Type = rb.goTypeInJ2gort(sig.canonParamTypes.get(1), imports);
        if (p1Type.equals("uint16")) {
          return "\tthis_.sbCore().insertChar(arg0, arg1)\n" + retSelf;
        }
        if (p1Type.equals("string")) {
          return "\tthis_.sbCore().insertString(arg0, arg1)\n" + retSelf;
        }
        return "\tthis_.sbCore().insertString(arg0, StrVal(arg1))\n" + retSelf;
      }
      if (m.equals("setCharAt") && arity == 2) {
        return "\tthis_.sbCore().setCharAt(arg0, arg1)\n";
      }
      if (m.equals("deleteCharAt") && arity == 1) {
        return "\tthis_.sbCore().deleteCharAt(arg0)\n" + retSelf;
      }
      if (m.equals("delete") && arity == 2) {
        return "\tthis_.sbCore().deleteRange(arg0, arg1)\n" + retSelf;
      }
      if (m.equals("replace") && arity == 3) {
        return "\tthis_.sbCore().replaceRange(arg0, arg1, arg2)\n" + retSelf;
      }
      if (m.equals("reverse") && arity == 0) {
        return "\tthis_.sbCore().reverse()\n" + retSelf;
      }
      if (m.equals("substring") && arity == 1) {
        return "\treturn JStrSubFrom(this_.sbCore().str(), arg0)\n";
      }
      if (m.equals("substring") && arity == 2) {
        return "\treturn JStrSub(this_.sbCore().str(), arg0, arg1)\n";
      }
      if (m.equals("indexOf") && arity == 1) {
        return "\treturn JStrIndexOf(this_.sbCore().str(), arg0)\n";
      }
      if (m.equals("indexOf") && arity == 2) {
        return "\treturn JStrIndexOfFrom(this_.sbCore().str(), arg0, arg1)\n";
      }
      if (m.equals("repeat") && arity == 2) {
        if (arg1CheckNeeded()) {
          if (p0.equals("int32")) {
            return "\tif arg1 < 0 { panic(javalang.NewIllegalArgumentException(\"count < 0\")) }\n\tfor i := int32(0); i < arg1; i++ { this_.sbCore().appendCodePoint(arg0) }\n"
                + retSelf;
          }
          String sExpr = p0.equals("string") ? "StrValString(arg0)" : "StrVal(arg0)";
          return "\tif arg1 < 0 { panic(javalang.NewIllegalArgumentException(\"count < 0\")) }\n\tunit := "
              + sExpr
              + "\n\tfor i := int32(0); i < arg1; i++ { this_.sbCore().appendString(unit) }\n"
              + retSelf;
        }
      }
      if (m.equals("length") && arity == 0) {
        return "\treturn this_.sbCore().length()\n";
      }
      if (m.equals("isEmpty") && arity == 0) {
        return "\treturn this_.sbCore().length() == 0\n";
      }
      if (m.equals("charAt") && arity == 1) {
        return "\treturn this_.sbCore().charAt(arg0)\n";
      }
      if (m.equals("toString") && arity == 0) {
        return "\treturn this_.sbCore().str()\n";
      }
      if (m.equals("setLength") && arity == 1) {
        return "\tthis_.sbCore().setLength(arg0)\n";
      }
      if (m.equals("compareTo") && arity == 1) {
        return "\tother := CheckNotNullPtr(Cast[*Ext_java_lang_StringBuilder](arg0))\n\treturn JStrCompareTo(this_.sbCore().str(), other.sbCore().str())\n";
      }
    }

    // External Enum methods (.equals / .compareTo on any external enum + Ext_java_lang_EnumBase)
    if (te.getKind() == javax.lang.model.element.ElementKind.ENUM) {
      if (m.equals("equals") && arity == 1) {
        return "\treturn RefEquals(this_, arg0)\n";
      }
      if (m.equals("compareTo") && arity == 1) {
        return "\tCheckNotNull(arg0)\n\treturn EnumCompareTo(this_, arg0)\n";
      }
    }

    if (qn.equals("java.lang.Enum")) {
      if (m.equals("name") && arity == 0) return "\treturn this_.Core_.(*EnumBase).Name_\n";
      if (m.equals("ordinal") && arity == 0) return "\treturn this_.Core_.(*EnumBase).Ordinal_\n";
      if (m.equals("toString") && arity == 0) return "\treturn this_.Core_.(*EnumBase).Name_\n";
      if (m.equals("equals") && arity == 1) {
        return "\tif this_.Self_ != nil { return RefEquals(this_.Self_, arg0) }\n\treturn RefEquals(this_, arg0)\n";
      }
      if (m.equals("hashCode") && arity == 0) return "\treturn this_.Core_.(*EnumBase).HashCode()\n";
      if (m.equals("compareTo") && arity == 1) {
        return "\tCheckNotNull(arg0)\n\treturn this_.Core_.(*EnumBase).Ordinal_ - EnumOrdinal(arg0)\n";
      }
    }

    if (qn.equals("java.lang.Record")) {
      if (m.equals("equals") && arity == 1) {
        return "\tif this_.Self_ != nil && !RefEquals(this_.Self_, this_) { return this_.Self_.Equals(arg0) }\n\treturn RefEquals(this_, arg0)\n";
      }
      if (m.equals("hashCode") && arity == 0) {
        return "\tif this_.Self_ != nil && !RefEquals(this_.Self_, this_) { return this_.Self_.HashCode() }\n\treturn javalang.IdentityHashCode(this_)\n";
      }
      if (m.equals("toString") && arity == 0) {
        return "\tif this_.Self_ != nil && !RefEquals(this_.Self_, this_) { return this_.Self_.String() }\n\treturn \"Record[]\"\n";
      }
    }

    // Annotation value stubs (only accessed if inspected via reflection; harmless constants)
    if (te.getKind() == javax.lang.model.element.ElementKind.ANNOTATION_TYPE) {
      if (retType.equals("string")) return "\treturn \"\"\n";
      if (retType.startsWith("[]")) return "\treturn nil\n";
      if (retType.equals("bool")) return "\treturn false\n";
      if (retType.startsWith("*")) return "\treturn nil\n";
    }

    switch (qn) {
      case "com.google.common.base.CaseFormat" -> {
        if (m.equals("to") && arity == 2) {
          return "\treturn guavaCaseFormatConvert(this_.Name_, CheckNotNullPtr(arg0).Name_, CheckNotNullString(arg1))\n";
        }
      }
      case "com.google.common.base.StandardSystemProperty" -> {
        if (m.equals("value") && arity == 0) {
          return "\treturn guavaStandardSystemPropertyValue(this_.Name_)\n";
        }
      }
      case "com.google.common.base.CharMatcher" -> {
        if (m.equals("or") && arity == 1) return "\treturn charMatcherOr(this_, CheckNotNull(arg0))\n";
        if (m.equals("matchesAllOf") && arity == 1) return "\treturn charMatcherMatchesAllOf(this_, CheckNotNull(arg0))\n";
        if (m.equals("lastIndexIn") && arity == 1) return "\treturn charMatcherLastIndexIn(this_, CheckNotNull(arg0))\n";
        if (m.equals("trimFrom") && arity == 1) return "\treturn charMatcherTrimFrom(this_, CheckNotNull(arg0))\n";
        if (m.equals("collapseFrom") && arity == 2) return "\treturn charMatcherCollapseFrom(this_, CheckNotNull(arg0), arg1)\n";
        if ((m.equals("apply") || m.equals("test")) && arity == 1) return "\treturn charMatcherMatchesBoxed(this_, arg0)\n";
        if (m.equals("equals") && arity == 1) return "\treturn RefEquals(this_, arg0)\n";
      }
      case "com.google.common.base.Joiner" -> {
        if (m.equals("appendTo") && arity == 2) {
          return "\tCheckNotNullPtr(arg0)\n\targ0.sbCore().appendString(guavaJoinerJoinIter(this_, CheckNotNull(arg1)))\n\treturn arg0\n";
        }
        if (m.equals("join") && arity == 1) {
          if (p0.startsWith("[]")) {
            return "\treturn guavaJoinerJoinSlice(this_, CheckNotNullSlice(arg0))\n";
          }
          return "\treturn guavaJoinerJoinIter(this_, CheckNotNull(arg0))\n";
        }
        if (m.equals("join") && arity == 3) {
          return "\treturn guavaJoinerJoinVarargs(this_, arg0, arg1, arg2)\n";
        }
      }
      case "com.google.common.base.MoreObjects.ToStringHelper" -> {
        if (m.equals("omitNullValues") && arity == 0) {
          return "\tthis_.toStringHelperOmitNulls()\n\treturn this_\n";
        }
        if (m.equals("add") && arity == 2) {
          return "\tthis_.toStringHelperAdd(CheckNotNullString(arg0), arg1)\n\treturn this_\n";
        }
        if (m.equals("toString") && arity == 0) {
          return "\treturn this_.toStringHelperBuild()\n";
        }
      }
      case "com.google.common.base.Optional", "java.util.Optional" -> {
        if ((m.equals("isPresent") || m.equals("isEmpty")) && arity == 0) {
          return m.equals("isPresent") ? "\treturn this_.Core_ != nil\n" : "\treturn this_.Core_ == nil\n";
        }
        if ((m.equals("get") || m.equals("orElseThrow")) && arity == 0) {
          return "\tif this_.Core_ == nil { panic(javalang.NewNoSuchElementException(\"Optional.get() on empty\")) }\n\treturn this_.Core_\n";
        }
        if (m.equals("orElseThrow") && arity == 1) {
          return "\tif this_.Core_ != nil { return this_.Core_ }\n\tex := CheckNotNull(arg0).Get()\n\tpanic(CheckNotNull(ex))\n";
        }
        if (m.equals("ifPresent") && arity == 1) {
          return "\tCheckNotNull(arg0)\n\tif this_.Core_ != nil { arg0.Accept(this_.Core_) }\n";
        }
        if ((m.equals("orElse") || m.equals("or")) && arity == 1 && retType.equals("any")) {
          return "\tif this_.Core_ != nil { return this_.Core_ }\n\treturn ToAny(arg0)\n";
        }
        if (m.equals("or") && arity == 1 && retType.contains("Optional")) {
          return "\tCheckNotNull(arg0)\n\tif this_.Core_ != nil { return this_ }\n\treturn arg0\n";
        }
        if (m.equals("orNull") && arity == 0) {
          return "\treturn this_.Core_\n";
        }
        if (m.equals("transform") && arity == 1) {
          return "\tCheckNotNull(arg0)\n\tif this_.Core_ == nil { return &Ext_com_google_common_base_OptionalBase{} }\n\tr := CheckNotNull(ToAny(arg0.Apply(this_.Core_)))\n\tout := &Ext_com_google_common_base_OptionalBase{Core_: r}\n\tout.BindExtSelf_(out)\n\treturn out\n";
        }
        if (m.equals("equals") && arity == 1) {
          return "\tif o, ok := arg0.(interface{ getCore() any }); ok && !javalang.IsNil(arg0) {\n\t\treturn javalang.ObjectsEquals(this_.Core_, o.getCore())\n\t}\n\treturn false\n";
        }
        if (m.equals("hashCode") && arity == 0) {
          return "\tif this_.Core_ == nil { return 0x79a31aac - 0x20407160 }\n\treturn 0x598df91c + javalang.ObjectHashCode(this_.Core_)\n";
        }
      }
      case "com.google.common.base.Splitter" -> {
        if (m.equals("omitEmptyStrings") && arity == 0) return "\treturn splitterOmitEmpty(this_)\n";
        if (m.equals("trimResults") && arity == 0) return "\treturn splitterTrimResults(this_)\n";
        if (m.equals("limit") && arity == 1) return "\treturn splitterLimit(this_, arg0)\n";
        if (m.equals("split") && arity == 1) return "\treturn splitterSplit(this_, CheckNotNull(arg0))\n";
        if (m.equals("splitToList") && arity == 1) return "\treturn splitterSplitToList(this_, CheckNotNull(arg0))\n";
      }
      case "com.google.common.escape.Escaper", "com.google.common.escape.UnicodeEscaper" -> {
        if (m.equals("escape") && arity == 1) {
          return "\tif this_.Self_ != nil && !RefEquals(this_.Self_, this_) { return this_.Self_.Escape_1_string(arg0) }\n\treturn escapePercentCore(this_.Core_, CheckNotNullString(arg0))\n";
        }
        if (m.equals("asFunction") && arity == 0) {
          return "\tself := Ext_com_google_common_escape_Escaper(this_)\n\tif this_.Self_ != nil { self = this_.Self_ }\n\treturn newEscaperFunction(self)\n";
        }
      }
      case "com.google.common.net.PercentEscaper" -> {
        if (m.equals("escape") && arity == 1) {
          return "\treturn escapePercentCore(this_.Core_, CheckNotNullString(arg0))\n";
        }
        if (m.equals("asFunction") && arity == 0) {
          return "\treturn newEscaperFunction(this_)\n";
        }
      }
      case "com.google.common.hash.HashCode" -> {
        if (m.equals("asInt") && arity == 0) return "\treturn hashCodeAsInt(this_.Core_)\n";
        if (m.equals("asLong") && arity == 0) return "\treturn hashCodeAsLong(this_.Core_)\n";
      }
      case "com.google.common.io.BaseEncoding" -> {
        if (m.equals("encode") && arity == 1) return "\treturn base64Encode(CheckNotNullSlice(arg0))\n";
        if (m.equals("decode") && arity == 1) return "\treturn base64Decode(StrVal(CheckNotNull(arg0)))\n";
      }
      case "com.google.common.io.CharSink" -> {
        if (m.equals("write") && arity == 1) return "\tcharSinkWrite(this_.Core_, StrVal(CheckNotNull(arg0)))\n";
      }
      case "com.google.common.io.CharSource" -> {
        if (m.equals("openStream") && arity == 0) return "\treturn charSourceOpenStream(this_.Core_)\n";
        if (m.equals("read") && arity == 0) return "\treturn charSourceRead(this_.Core_)\n";
      }
      case "com.google.gson.Gson" -> {
        if (m.equals("toJson") && arity == 1) {
          return "\treturn gsonToJson(this_, arg0)\n";
        }
        if (m.equals("toJson") && arity == 3) {
          return "\tgsonToJsonWriter(this_, arg0, arg2)\n";
        }
        if (m.equals("fromJson") && arity == 2) {
          if (rb.goTypeInJ2gort(sig.canonParamTypes.get(1), imports).contains("Class")) {
            return "\treturn gsonFromJsonClass(this_, arg0, arg1)\n";
          }
          return "\treturn gsonFromJsonType(this_, arg0, arg1)\n";
        }
      }
      case "com.google.gson.GsonBuilder" -> {
        if (m.equals("setPrettyPrinting") && arity == 0) {
          return "\tgsonBuilderSetPretty(this_)\n\treturn this_\n";
        }
        if (m.equals("disableHtmlEscaping") && arity == 0) {
          return "\tgsonBuilderDisableHtml(this_)\n\treturn this_\n";
        }
        if (m.equals("registerTypeHierarchyAdapter") && arity == 2) {
          return "\tgsonBuilderRegisterHierarchyAdapter(this_, arg0, arg1)\n\treturn this_\n";
        }
        if (m.equals("create") && arity == 0) {
          return "\treturn gsonBuilderCreate(this_)\n";
        }
      }
      case "com.google.gson.JsonElement",
          "com.google.gson.JsonArray",
          "com.google.gson.JsonObject",
          "com.google.gson.JsonPrimitive",
          "com.google.gson.JsonNull" -> {
        return buildJsonElementMethodBody(qn, m, arity);
      }
      case "com.google.gson.TypeAdapter" -> {
        if (m.equals("write") && arity == 2) {
          return "\tif this_.Self_ != nil && !RefEquals(this_.Self_, this_) { this_.Self_.Write(arg0, arg1); return }\n\tpanic(\"UNIMPLEMENTED_EXT: TypeAdapter.write\")\n";
        }
        if (m.equals("read") && arity == 1) {
          return "\tif this_.Self_ != nil && !RefEquals(this_.Self_, this_) { return this_.Self_.Read(arg0) }\n\tpanic(\"UNIMPLEMENTED_EXT: TypeAdapter.read\")\n";
        }
      }
      case "com.google.gson.reflect.TypeToken" -> {
        if (m.equals("getType") && arity == 0) {
          return "\treturn this_.Core_\n";
        }
      }
      case "com.google.gson.stream.JsonReader" -> {
        if (m.equals("beginArray") && arity == 0) return "\tgsonReader(this_).BeginArray()\n";
        if (m.equals("endArray") && arity == 0) return "\tgsonReader(this_).EndArray()\n";
        if (m.equals("beginObject") && arity == 0) return "\tgsonReader(this_).BeginObject()\n";
        if (m.equals("endObject") && arity == 0) return "\tgsonReader(this_).EndObject()\n";
        if (m.equals("hasNext") && arity == 0) return "\treturn gsonReader(this_).HasNext()\n";
        if (m.equals("nextName") && arity == 0) return "\treturn gsonReader(this_).NextName()\n";
        if (m.equals("nextString") && arity == 0) return "\treturn gsonReader(this_).NextString()\n";
        if (m.equals("skipValue") && arity == 0) return "\tgsonReader(this_).SkipValue()\n";
        if (m.equals("close") && arity == 0) return "\tgsonReader(this_).Close()\n";
      }
      case "com.google.gson.stream.JsonWriter" -> {
        if (m.equals("setHtmlSafe") && arity == 1) return "\tgsonWriter(this_).SetHtmlSafe(arg0)\n";
        if (m.equals("beginArray") && arity == 0) return "\tgsonWriter(this_).BeginArray()\n\treturn this_\n";
        if (m.equals("endArray") && arity == 0) return "\tgsonWriter(this_).EndArray()\n\treturn this_\n";
        if (m.equals("beginObject") && arity == 0) return "\tgsonWriter(this_).BeginObject()\n\treturn this_\n";
        if (m.equals("endObject") && arity == 0) return "\tgsonWriter(this_).EndObject()\n\treturn this_\n";
        if (m.equals("name") && arity == 1) return "\tgsonWriter(this_).Name(CheckNotNullString(arg0))\n\treturn this_\n";
        if (m.equals("value") && arity == 1 && p0.equals("string")) {
          return "\tif arg0 == NullString { gsonWriter(this_).NullValue() } else { gsonWriter(this_).ValueString(arg0) }\n\treturn this_\n";
        }
        if (m.equals("value") && arity == 1 && p0.equals("int64")) {
          return "\tgsonWriter(this_).ValueLong(arg0)\n\treturn this_\n";
        }
        if (m.equals("jsonValue") && arity == 1) {
          return "\tif arg0 == NullString { gsonWriter(this_).NullValue() } else { gsonWriter(this_).JsonValue(arg0) }\n\treturn this_\n";
        }
        if (m.equals("flush") && arity == 0) return "\tgsonWriter(this_).Flush()\n";
        if (m.equals("close") && arity == 0) return "\tgsonWriter(this_).Close()\n";
      }
      case "java.io.File" -> {
        if (m.equals("getName") && arity == 0) return "\treturn unwrapFile(this_).GetName()\n";
        if (m.equals("getParent") && arity == 0) {
          return "\tp, ok := unwrapFile(this_).Parent()\n\tif !ok { return NullString }\n\treturn p\n";
        }
        if (m.equals("getParentFile") && arity == 0) {
          return "\tpf := unwrapFile(this_).GetParentFile()\n\tif pf == nil { return nil }\n\treturn wrapFile(pf)\n";
        }
        if (m.equals("getAbsolutePath") && arity == 0) return "\treturn unwrapFile(this_).GetAbsolutePath()\n";
        if (m.equals("isDirectory") && arity == 0) return "\treturn unwrapFile(this_).IsDirectory()\n";
        if (m.equals("isFile") && arity == 0) return "\treturn unwrapFile(this_).IsFile()\n";
        if (m.equals("exists") && arity == 0) return "\treturn unwrapFile(this_).Exists()\n";
        if (m.equals("createNewFile") && arity == 0) {
          return "\tok, err := unwrapFile(this_).CreateNewFile()\n\tioErr(err)\n\treturn ok\n";
        }
        if (m.equals("delete") && arity == 0) return "\treturn unwrapFile(this_).Delete()\n";
        if (m.equals("mkdirs") && arity == 0) return "\treturn unwrapFile(this_).Mkdirs()\n";
        if (m.equals("listFiles") && arity == 0) {
          return "\tfiles := unwrapFile(this_).ListFiles()\n\tif files == nil { return nil }\n\tout := make([]*Ext_java_io_File, len(files))\n\tfor i, f := range files { out[i] = wrapFile(f) }\n\treturn out\n";
        }
        if (m.equals("compareTo") && arity == 1) {
          return "\treturn int32(unwrapFile(this_).CompareTo(CheckNotNullPtr(unwrapFile(Cast[*Ext_java_io_File](arg0)))))\n";
        }
        if (m.equals("equals") && arity == 1) {
          return "\to, ok := arg0.(*Ext_java_io_File)\n\treturn ok && o != nil && unwrapFile(this_).Equals(unwrapFile(o))\n";
        }
        if (m.equals("hashCode") && arity == 0) return "\treturn unwrapFile(this_).HashCode()\n";
        if (m.equals("toString") && arity == 0) return "\treturn unwrapFile(this_).String()\n";
        if (m.equals("toPath") && arity == 0) return "\treturn wrapPath(javanio.FileToPath(unwrapFile(this_)))\n";
      }
      case "java.io.ByteArrayInputStream" -> {
        if (m.equals("close") && arity == 0) {
          return "";
        }
      }
      case "java.io.InputStream",
          "java.io.FilterInputStream",
          "java.io.BufferedInputStream",
          "java.io.FileInputStream",
          "java.util.zip.InflaterInputStream",
          "java.util.zip.GZIPInputStream" -> {
        if (m.equals("close") && arity == 0) {
          return "\textInputStreamClose(this_.Core_)\n";
        }
      }
      case "java.io.ObjectInputStream" -> {
        if (m.equals("defaultReadObject") && arity == 0) return "";
        if (m.equals("close") && arity == 0) return "\textInputStreamClose(this_.Core_)\n";
      }
      case "java.io.OutputStream",
          "java.io.FilterOutputStream",
          "java.io.BufferedOutputStream",
          "java.io.FileOutputStream" -> {
        if (m.equals("write") && arity == 1 && p0.equals("int32")) {
          return "\textOutputStreamWriteByte(this_.Core_, byte(arg0))\n";
        }
        if (m.equals("write") && arity == 1 && p0.equals("[]int8")) {
          return "\textOutputStreamWriteSlice(this_.Core_, CheckNotNullSlice(arg0))\n";
        }
        if (m.equals("flush") && arity == 0) {
          return "\textOutputStreamFlush(this_.Core_)\n";
        }
        if (m.equals("close") && arity == 0) {
          if (qn.equals("java.io.FilterOutputStream")) {
            return "\tif this_.Self_ != nil { this_.Self_.Flush() } else { this_.Flush() }\n\textOutputStreamClose(this_.Core_)\n";
          }
          return "\textOutputStreamClose(this_.Core_)\n";
        }
      }
      case "java.io.ByteArrayOutputStream" -> {
        if (m.equals("write") && arity == 1 && p0.equals("int32")) {
          return "\t_ = this_.Core_.(*javaio.ByteArrayOutputStream).WriteByte(byte(arg0))\n";
        }
        if (m.equals("write") && arity == 1 && p0.equals("[]int8")) {
          return "\t_, _ = this_.Core_.(*javaio.ByteArrayOutputStream).Write(int8ToBytes(CheckNotNullSlice(arg0)))\n";
        }
        if (m.equals("toByteArray") && arity == 0) {
          return "\treturn bytesToInt8(this_.Core_.(*javaio.ByteArrayOutputStream).ToByteArray())\n";
        }
        if (m.equals("size") && arity == 0) {
          return "\treturn int32(this_.Core_.(*javaio.ByteArrayOutputStream).Size())\n";
        }
        if (m.equals("toString") && arity == 1) {
          return "\tcs := unwrapCharset(CheckNotNull(arg0))\n\treturn cs.Decode(this_.Core_.(*javaio.ByteArrayOutputStream).ToByteArray())\n";
        }
        if (m.equals("toString") && arity == 0) {
          return "\treturn this_.Core_.(*javaio.ByteArrayOutputStream).ToString()\n";
        }
        if (m.equals("flush") && arity == 0) return "";
        if (m.equals("close") && arity == 0) return "";
      }
      case "java.io.PrintStream" -> {
        if (m.equals("flush") && arity == 0) return "\tunwrapPrintStream(this_).Flush()\n";
        if (m.equals("close") && arity == 0) return "\tunwrapPrintStream(this_).Close()\n";
        if (m.equals("write") && arity == 1 && p0.equals("int32")) {
          return "\t_, _ = unwrapPrintStream(this_).Write([]byte{byte(arg0)})\n";
        }
        if (m.equals("write") && arity == 1 && p0.equals("[]int8")) {
          return "\t_, _ = unwrapPrintStream(this_).Write(int8ToBytes(CheckNotNullSlice(arg0)))\n";
        }
        if (m.equals("print") && arity == 1) {
          return "\tunwrapPrintStream(this_).Print(StrVal(arg0))\n";
        }
        if (m.equals("println") && arity == 0) {
          return "\tunwrapPrintStream(this_).Println(\"\")\n";
        }
        if (m.equals("println") && arity == 1) {
          return "\tunwrapPrintStream(this_).Println(StrVal(arg0))\n";
        }
        if ((m.equals("printf") || m.equals("format")) && arity == 2) {
          return "\tunwrapPrintStream(this_).Print(StringFormat(CheckNotNullString(arg0), arg1...))\n\treturn this_\n";
        }
        if ((m.equals("printf") || m.equals("format")) && arity == 3) {
          return "\t_ = arg0\n\tunwrapPrintStream(this_).Print(StringFormat(CheckNotNullString(arg1), arg2...))\n\treturn this_\n";
        }
        if (m.equals("append") && arity == 1 && p0.equals("uint16")) {
          return "\tunwrapPrintStream(this_).Print(javalang.StringFromChar(arg0))\n\treturn this_\n";
        }
        if (m.equals("append") && arity == 1) {
          return "\tunwrapPrintStream(this_).Print(StrVal(arg0))\n\treturn this_\n";
        }
        if (m.equals("append") && arity == 3) {
          return "\ts := StrVal(arg0)\n\tunwrapPrintStream(this_).Print(JStrSub(s, arg1, arg2))\n\treturn this_\n";
        }
      }
      case "java.io.Reader",
          "java.io.BufferedReader",
          "java.io.InputStreamReader",
          "java.io.StringReader" -> {
        if (m.equals("read") && arity == 0) {
          return "\treturn extReaderReadOne(this_.Core_)\n";
        }
        if (m.equals("readLine") && arity == 0) {
          return "\treturn extBufferedReaderReadLine(this_.Core_)\n";
        }
        if (m.equals("close") && arity == 0) {
          return "\textReaderClose(this_.Core_)\n";
        }
      }
      case "java.io.Writer",
          "java.io.BufferedWriter",
          "java.io.OutputStreamWriter",
          "java.io.StringWriter",
          "java.io.PrintWriter" -> {
        String retReceiver =
            qn.equals("java.io.Writer")
                ? "\tif this_.Self_ != nil { return this_.Self_ }\n\treturn this_\n"
                : "\treturn this_\n";
        if (m.equals("write") && arity == 1 && p0.equals("int32")) {
          return "\textWriterWriteString(this_.Core_, javalang.StringFromChar(uint16(arg0)))\n";
        }
        if (m.equals("write") && arity == 1 && p0.equals("string")) {
          return "\textWriterWriteString(this_.Core_, CheckNotNullString(arg0))\n";
        }
        if (m.equals("append") && arity == 1 && p0.equals("uint16")) {
          return "\textWriterWriteString(this_.Core_, javalang.StringFromChar(arg0))\n" + retReceiver;
        }
        if (m.equals("append") && arity == 1) {
          return "\textWriterWriteString(this_.Core_, StrVal(arg0))\n" + retReceiver;
        }
        if (m.equals("append") && arity == 3) {
          return "\ts := StrVal(arg0)\n\textWriterWriteString(this_.Core_, JStrSub(s, arg1, arg2))\n" + retReceiver;
        }
        if (m.equals("flush") && arity == 0) {
          return "\textWriterFlush(this_.Core_)\n";
        }
        if (m.equals("close") && arity == 0) {
          return "\textWriterClose(this_.Core_)\n";
        }
        if (m.equals("toString") && arity == 0 && qn.equals("java.io.StringWriter")) {
          return "\treturn this_.Core_.(*javaio.StringWriter).String()\n";
        }
      }
      case "java.util.zip.Deflater" -> {
        if (m.equals("setLevel") && arity == 1) {
          return "\tzipDeflaterSetLevel(this_.Core_, int(arg0))\n";
        }
      }
      case "java.util.zip.DeflaterOutputStream", "java.util.zip.GZIPOutputStream" -> {
        if (m.equals("write") && arity == 1 && p0.equals("int32")) {
          return "\textOutputStreamWriteByte(this_.Core_, byte(arg0))\n";
        }
        if (m.equals("write") && arity == 1 && p0.equals("[]int8")) {
          return "\textOutputStreamWriteSlice(this_.Core_, CheckNotNullSlice(arg0))\n";
        }
        if (m.equals("flush") && arity == 0) {
          return "\textOutputStreamFlush(this_.Core_)\n";
        }
        if (m.equals("finish") && arity == 0) {
          return "\tgzipOutputStreamFinish(this_.Core_)\n";
        }
        if (m.equals("close") && arity == 0) {
          return "\textOutputStreamClose(this_.Core_)\n";
        }
      }
      case "java.util.zip.ZipEntry" -> {
        if (m.equals("getName") && arity == 0) return "\treturn this_.Core_.(*javazip.ZipEntry).GetName()\n";
        if (m.equals("getSize") && arity == 0) return "\treturn this_.Core_.(*javazip.ZipEntry).GetSize()\n";
        if (m.equals("isDirectory") && arity == 0) return "\treturn this_.Core_.(*javazip.ZipEntry).IsDirectory()\n";
        if (m.equals("toString") && arity == 0) return "\treturn this_.Core_.(*javazip.ZipEntry).String()\n";
      }
      case "java.util.zip.ZipFile" -> {
        if (m.equals("getEntry") && arity == 1) {
          return "\te := this_.Core_.(*javazip.ZipFile).GetEntry(CheckNotNullString(arg0))\n\tif e == nil { return nil }\n\treturn &Ext_java_util_zip_ZipEntry{Core_: e}\n";
        }
        if (m.equals("getInputStream") && arity == 1) {
          return "\tin, err := this_.Core_.(*javazip.ZipFile).GetInputStream(CheckNotNullPtr(arg0).Core_.(*javazip.ZipEntry))\n\tioErr(err)\n\treturn wrapInputStream(in)\n";
        }
        if (m.equals("close") && arity == 0) {
          return "\tioErr(this_.Core_.(*javazip.ZipFile).Close())\n";
        }
      }
      case "java.util.zip.ZipInputStream" -> {
        if (m.equals("getNextEntry") && arity == 0) {
          return "\te, err := this_.Core_.(*javazip.ZipInputStream).GetNextEntry()\n\tioErr(err)\n\tif e == nil { return nil }\n\treturn &Ext_java_util_zip_ZipEntry{Core_: e}\n";
        }
        if (m.equals("close") && arity == 0) {
          return "\tioErr(this_.Core_.(*javazip.ZipInputStream).Close())\n";
        }
      }
      case "java.net.URI" -> {
        if (m.equals("getScheme") && arity == 0) {
          return "\ts, ok := this_.Core_.(*javanio.URI).GetScheme()\n\tif !ok { return NullString }\n\treturn s\n";
        }
        if (m.equals("isAbsolute") && arity == 0) return "\treturn this_.Core_.(*javanio.URI).IsAbsolute()\n";
        if (m.equals("getPath") && arity == 0) {
          return "\tp, ok := this_.Core_.(*javanio.URI).GetPath()\n\tif !ok { return NullString }\n\treturn p\n";
        }
        if (m.equals("equals") && arity == 1) {
          return "\to, ok := arg0.(*Ext_java_net_URI)\n\treturn ok && o != nil && this_.Core_.(*javanio.URI).Equals(o.Core_.(*javanio.URI))\n";
        }
        if (m.equals("compareTo") && arity == 1) {
          return "\to := CheckNotNullPtr(Cast[*Ext_java_net_URI](arg0))\n\treturn JStrCompareTo(this_.Core_.(*javanio.URI).String(), o.Core_.(*javanio.URI).String())\n";
        }
        if (m.equals("hashCode") && arity == 0) return "\treturn this_.Core_.(*javanio.URI).HashCode()\n";
        if (m.equals("toString") && arity == 0) return "\treturn this_.Core_.(*javanio.URI).String()\n";
      }
      case "java.nio.charset.Charset" -> {
        if (m.equals("name") && arity == 0) return "\treturn unwrapCharset(this_).Name()\n";
        if (m.equals("newEncoder") && arity == 0) return "\treturn wrapCharsetEncoder(unwrapCharset(this_).NewEncoder())\n";
        if (m.equals("equals") && arity == 1) {
          return "\tif o, ok := arg0.(Ext_java_nio_charset_Charset); ok && !javalang.IsNil(arg0) {\n\t\treturn unwrapCharset(this_).Equals(unwrapCharset(o))\n\t}\n\treturn false\n";
        }
        if (m.equals("compareTo") && arity == 1) {
          return "\to := CheckNotNull(AsIface[Ext_java_nio_charset_Charset](arg0))\n\treturn JStrCompareTo(strings.ToLower(unwrapCharset(this_).Name()), strings.ToLower(unwrapCharset(o).Name()))\n";
        }
        if (m.equals("hashCode") && arity == 0) return "\treturn unwrapCharset(this_).HashCode()\n";
        if (m.equals("toString") && arity == 0) return "\treturn unwrapCharset(this_).String()\n";
      }
      case "java.nio.charset.CharsetEncoder" -> {
        if (m.equals("onMalformedInput") && arity == 1) {
          return "\tunwrapCharsetEncoder(this_).OnMalformedInput(unwrapCodingAction(arg0))\n\treturn this_\n";
        }
        if (m.equals("onUnmappableCharacter") && arity == 1) {
          return "\tunwrapCharsetEncoder(this_).OnUnmappableCharacter(unwrapCodingAction(arg0))\n\treturn this_\n";
        }
        if (m.equals("canEncode") && arity == 1 && p0.equals("uint16")) {
          return "\treturn unwrapCharsetEncoder(this_).CanEncodeChar(arg0)\n";
        }
      }
      case "java.nio.file.FileSystem" -> {
        if (m.equals("close") && arity == 0) return "";
        if (m.equals("getPath") && arity == 2) return "\treturn newExtPath(CheckNotNullString(arg0), arg1)\n";
        if (m.equals("getPathMatcher") && arity == 1) return "\treturn newExtPathMatcher(CheckNotNullString(arg0))\n";
      }
      case "java.nio.file.SimpleFileVisitor" -> {
        if (m.equals("visitFile") && arity == 2) {
          return "\tif this_.Self_ != nil && !RefEquals(this_.Self_, this_) { return this_.Self_.VisitFile(arg0, arg1) }\n\treturn Ext_java_nio_file_FileVisitResult_F_CONTINUE\n";
        }
        if (m.equals("visitFileFailed") && arity == 2) {
          return "\tif this_.Self_ != nil && !RefEquals(this_.Self_, this_) { return this_.Self_.VisitFileFailed(arg0, arg1) }\n\tpanic(CheckNotNull(arg1))\n";
        }
      }
      case "java.nio.file.attribute.FileTime" -> {
        if (m.equals("equals") && arity == 1) {
          return "\to, ok := arg0.(*Ext_java_nio_file_attribute_FileTime)\n\treturn ok && o != nil && this_.Core_.(*javanio.FileTime).Equals(o.Core_.(*javanio.FileTime))\n";
        }
        if (m.equals("compareTo") && arity == 1) {
          return "\to := CheckNotNullPtr(Cast[*Ext_java_nio_file_attribute_FileTime](arg0))\n\treturn int32(this_.Core_.(*javanio.FileTime).CompareTo(o.Core_.(*javanio.FileTime)))\n";
        }
        if (m.equals("hashCode") && arity == 0) return "\treturn this_.Core_.(*javanio.FileTime).HashCode()\n";
        if (m.equals("toString") && arity == 0) return "\treturn this_.Core_.(*javanio.FileTime).String()\n";
      }
      case "java.lang.Runtime" -> {
        if (m.equals("freeMemory") && arity == 0) return "\treturn runtimeFreeMemory()\n";
        if (m.equals("totalMemory") && arity == 0) return "\treturn runtimeTotalMemory()\n";
        if (m.equals("maxMemory") && arity == 0) return "\treturn runtimeMaxMemory()\n";
      }
      case "java.lang.Thread" -> {
        if (m.equals("run") && arity == 0) {
          return "\tif r, ok := this_.Core_.(Ext_java_lang_Runnable); ok && r != nil { r.Run() }\n";
        }
        if (m.equals("setDaemon") && arity == 1) return "";
      }
      case "java.lang.ThreadLocal" -> {
        if (m.equals("get") && arity == 0) return "\treturn this_.Core_.(*jThreadLocal).get()\n";
        if (m.equals("set") && arity == 1) return "\tthis_.Core_.(*jThreadLocal).set(arg0)\n";
        if (m.equals("remove") && arity == 0) return "\tthis_.Core_.(*jThreadLocal).remove()\n";
      }
      case "java.lang.management.MemoryUsage" -> {
        if (m.equals("getUsed") && arity == 0) return "\treturn this_.Core_.(*jMemoryUsage).used\n";
        if (m.equals("getMax") && arity == 0) return "\treturn this_.Core_.(*jMemoryUsage).max\n";
      }
      case "java.util.ResourceBundle" -> {
        if (m.equals("getString") && arity == 1) {
          return "\treturn this_.Core_.(*shimsres.ResourceBundle).GetString(CheckNotNullString(arg0))\n";
        }
      }
      case "java.util.StringTokenizer" -> {
        if ((m.equals("hasMoreTokens") || m.equals("hasMoreElements")) && arity == 0) {
          return "\treturn this_.Core_.(*jStringTokenizer).hasMoreTokens()\n";
        }
        if (m.equals("nextToken") && arity == 0) {
          return "\treturn this_.Core_.(*jStringTokenizer).nextToken()\n";
        }
      }
      case "java.util.logging.Level" -> {
        if (m.equals("getName") && arity == 0) {
          return "\tif s, ok := this_.Core_.(string); ok { return s }\n\treturn \"INFO\"\n";
        }
      }
      case "java.util.logging.Logger" -> {
        if (m.equals("setLevel") && arity == 1) return "\tthis_.Core_ = arg0\n";
        if (m.equals("isLoggable") && arity == 1) return "\treturn loggerIsLoggable(this_, arg0)\n";
        if ((m.equals("log")
            || m.equals("severe")
            || m.equals("warning")
            || m.equals("info")
            || m.equals("fine")
            || m.equals("finer")
            || m.equals("finest"))) {
          return "";
        }
      }
      case "java.util.regex.Pattern" -> {
        if ((m.equals("pattern") || m.equals("toString")) && arity == 0) {
          return "\treturn this_.Core_.(*javaregex.Pattern).Pattern()\n";
        }
        if (m.equals("matcher") && arity == 1) {
          return "\ts := StrVal(CheckNotNull(arg0))\n\tm := &Ext_java_util_regex_Matcher{Core_: this_.Core_.(*javaregex.Pattern).Matcher(s)}\n\tm.BindExtSelf_(m)\n\treturn m\n";
        }
      }
      case "java.util.regex.Matcher" -> {
        if (m.equals("reset") && arity == 1) {
          return "\tthis_.Core_.(*javaregex.Matcher).ResetInput(StrVal(CheckNotNull(arg0)))\n\treturn this_\n";
        }
        if (m.equals("reset") && arity == 0) {
          return "\tthis_.Core_.(*javaregex.Matcher).Reset()\n\treturn this_\n";
        }
        if (m.equals("start") && arity == 0) return "\treturn int32(this_.Core_.(*javaregex.Matcher).Start())\n";
        if (m.equals("start") && arity == 1) return "\treturn int32(this_.Core_.(*javaregex.Matcher).StartGroup(int(arg0)))\n";
        if (m.equals("end") && arity == 0) return "\treturn int32(this_.Core_.(*javaregex.Matcher).End())\n";
        if (m.equals("end") && arity == 1) return "\treturn int32(this_.Core_.(*javaregex.Matcher).EndGroup(int(arg0)))\n";
        if (m.equals("group") && arity == 0) return "\treturn this_.Core_.(*javaregex.Matcher).Group0()\n";
        if (m.equals("group") && arity == 1 && p0.equals("int32")) {
          return "\ts, ok := this_.Core_.(*javaregex.Matcher).Group(int(arg0))\n\tif !ok { return NullString }\n\treturn s\n";
        }
        if (m.equals("group") && arity == 1 && p0.equals("string")) {
          return "\ts, ok := this_.Core_.(*javaregex.Matcher).GroupNamed(CheckNotNullString(arg0))\n\tif !ok { return NullString }\n\treturn s\n";
        }
        if (m.equals("matches") && arity == 0) return "\treturn this_.Core_.(*javaregex.Matcher).Matches()\n";
        if (m.equals("find") && arity == 0) return "\treturn this_.Core_.(*javaregex.Matcher).Find()\n";
        if (m.equals("find") && arity == 1) return "\treturn this_.Core_.(*javaregex.Matcher).FindFrom(int(arg0))\n";
        if (m.equals("lookingAt") && arity == 0) return "\treturn this_.Core_.(*javaregex.Matcher).LookingAt()\n";
        if (m.equals("replaceAll") && arity == 1 && p0.equals("string")) return "\treturn this_.Core_.(*javaregex.Matcher).ReplaceAll(CheckNotNullString(arg0))\n";
        if (m.equals("replaceFirst") && arity == 1 && p0.equals("string")) return "\treturn this_.Core_.(*javaregex.Matcher).ReplaceFirst(CheckNotNullString(arg0))\n";
        if (m.equals("region") && arity == 2) {
          return "\tthis_.Core_.(*javaregex.Matcher).Region(int(arg0), int(arg1))\n\treturn this_\n";
        }
        if (m.equals("regionStart") && arity == 0) return "\treturn int32(this_.Core_.(*javaregex.Matcher).RegionStart())\n";
        if (m.equals("regionEnd") && arity == 0) return "\treturn int32(this_.Core_.(*javaregex.Matcher).RegionEnd())\n";
        if (m.equals("hitEnd") && arity == 0) return "\treturn this_.Core_.(*javaregex.Matcher).HitEnd()\n";
      }
      case "javax.xml.parsers.SAXParserFactory" -> {
        if (m.equals("setValidating") || m.equals("setXIncludeAware") || m.equals("setFeature")) return "";
        if (m.equals("newSAXParser") && arity == 0) {
          return "\tp := &Ext_javax_xml_parsers_SAXParserBase{}\n\tp.BindExtSelf_(p)\n\treturn p\n";
        }
      }
      case "javax.xml.parsers.SAXParser" -> {
        if (m.equals("getXMLReader") && arity == 0) {
          return "\treturn &jSaxXMLReader{}\n";
        }
      }
      case "org.xml.sax.SAXException" -> {
        if (m.equals("toString") && arity == 0) {
          return "\treturn this_.Error()\n";
        }
      }
      case "org.kohsuke.args4j.CmdLineParser" -> {
        if (m.equals("setUsageWidth") && arity == 1) {
          return "\tthis_.args4jCore().usageWidth = int(arg0)\n";
        }
        if (m.equals("parseArgument") && arity == 1 && p0.equals("[]string")) {
          return "\targs4jParseTokens(this_, CheckNotNullSlice(arg0))\n";
        }
        if (m.equals("parseArgument") && arity == 1) {
          return "\tvar toks []string\n\tfor it := Iterate(CheckNotNull(arg0)); it.HasNext(); { toks = append(toks, CastString(it.Next())) }\n\targs4jParseTokens(this_, toks)\n";
        }
        if (m.equals("printUsage") && arity == 3) {
          return "\t_ = arg1\n\targs4jPrintUsage(this_, asWriter(arg0), arg2)\n";
        }
      }
      case "org.kohsuke.args4j.NamedOptionDef" -> {
        if (m.equals("name") && arity == 0) return "\treturn getOptionDefData(this_.Core_).name\n";
        if (m.equals("hidden") && arity == 0) return "\treturn getOptionDefData(this_.Core_).hidden\n";
      }
      case "org.kohsuke.args4j.OptionDef" -> {
        if (m.equals("hidden") && arity == 0) return "\treturn getOptionDefData(this_.Core_).hidden\n";
      }
      case "org.kohsuke.args4j.spi.IntOptionHandler" -> {
        if (m.equals("parseArguments") && arity == 1) {
          return "\tp, _ := this_.Core_.(*Ext_org_kohsuke_args4j_CmdLineParser)\n\ts := arg0.GetParameter(0)\n\tv, err := strconv.ParseInt(strings.TrimSpace(s), 10, 32)\n\tif err != nil {\n\t\tpanic(NewArgs4jCmdLineException(p, fmt.Sprintf(\"%q is not a valid number\", s)))\n\t}\n\tthis_.F_Setter.AddValue(ValueOfInt(int32(v)))\n\treturn 1\n";
        }
        if (m.equals("getDefaultMetaVariable") && arity == 0) {
          return "\treturn \"N\"\n";
        }
      }
      case "org.kohsuke.args4j.spi.OptionHandler",
          "org.kohsuke.args4j.spi.OneArgumentOptionHandler",
          "org.kohsuke.args4j.spi.StringOptionHandler" -> {
        if (m.equals("parseArguments") && arity == 1) {
          return "\tp, _ := this_.Core_.(*Ext_org_kohsuke_args4j_CmdLineParser)\n\t_ = p\n\ts := arg0.GetParameter(0)\n\tthis_.F_Setter.AddValue(BoxString(s))\n\treturn 1\n";
        }
        if (m.equals("getDefaultMetaVariable") && arity == 0) {
          return "\treturn \"VAL\"\n";
        }
      }
      case "org.kohsuke.args4j.spi.FieldSetter" -> {
        if (m.equals("getType") && arity == 0) return "\treturn this_.Core_.(*args4jFieldSetterData).typ\n";
        if (m.equals("isMultiValued") && arity == 0) return "\treturn this_.Core_.(*args4jFieldSetterData).multiValued\n";
        if (m.equals("asFieldSetter") && arity == 0) return "\treturn this_\n";
        if (m.equals("asAnnotatedElement") && arity == 0) return "\treturn nil\n";
        if (m.equals("addValue") && arity == 1) return "\tthis_.Core_.(*args4jFieldSetterData).addFn(arg0)\n";
      }
      default -> {}
    }

    // Throwable .toString() fallback on any external Throwable subclass
    if (rb.resolver().isThrowable(te) && m.equals("toString") && arity == 0) {
      return "\treturn this_.Error()\n";
    }

    return null;
  }

  private static boolean arg1CheckNeeded() {
    return true;
  }

  private static String buildJsonElementMethodBody(String qn, String m, int arity) {
    if (qn.equals("com.google.gson.JsonArray")) {
      if (m.equals("add") && arity == 1) return "\tgsonJsonArrayAdd(this_, arg0)\n";
      if (m.equals("size") && arity == 0) return "\treturn int32(unwrapJsonArray(this_).Size())\n";
      if (m.equals("get") && arity == 1) return "\treturn wrapJsonElement(unwrapJsonArray(this_).Get(int(arg0)))\n";
      if (m.equals("iterator") && arity == 0) return "\treturn gsonJsonArrayIterator(this_)\n";
      if (m.equals("forEach") && arity == 1) return "\tgsonJsonArrayForEach(this_, CheckNotNull(arg0))\n";
    }
    if (qn.equals("com.google.gson.JsonObject")) {
      if (m.equals("add") && arity == 2) return "\tgsonJsonObjectAdd(this_, CheckNotNullString(arg0), arg1)\n";
      if (m.equals("has") && arity == 1) return "\treturn unwrapJsonObject(this_).Has(CheckNotNullString(arg0))\n";
      if (m.equals("get") && arity == 1) return "\treturn wrapJsonElement(unwrapJsonObject(this_).Get(CheckNotNullString(arg0)))\n";
      if (m.equals("entrySet") && arity == 0) return "\treturn gsonJsonObjectEntrySet(this_)\n";
    }
    if (qn.equals("com.google.gson.JsonPrimitive")) {
      if (m.equals("isString") && arity == 0) return "\treturn unwrapJsonPrimitive(this_).IsString()\n";
    }
    if (m.equals("isJsonPrimitive") && arity == 0) return "\treturn unwrapJsonElement(this_).IsJsonPrimitive()\n";
    if (m.equals("isJsonNull") && arity == 0) return "\treturn unwrapJsonElement(this_).IsJsonNull()\n";
    if (m.equals("getAsJsonObject") && arity == 0) return "\treturn wrapJsonObject(unwrapJsonElement(this_).GetAsJsonObject())\n";
    if (m.equals("getAsJsonArray") && arity == 0) return "\treturn wrapJsonArray(unwrapJsonElement(this_).GetAsJsonArray())\n";
    if (m.equals("getAsJsonPrimitive") && arity == 0) return "\treturn wrapJsonPrimitive(unwrapJsonElement(this_).GetAsJsonPrimitive())\n";
    if (m.equals("getAsString") && arity == 0) return "\treturn unwrapJsonElement(this_).GetAsString()\n";
    if (m.equals("getAsInt") && arity == 0) return "\treturn unwrapJsonElement(this_).GetAsInt()\n";
    if (m.equals("toString") && arity == 0) return "\treturn unwrapJsonElement(this_).String()\n";
    if (m.equals("equals") && arity == 1) {
      return "\tif o, ok := arg0.(Ext_com_google_gson_JsonElement); ok && !javalang.IsNil(arg0) {\n\t\treturn javalang.ObjectsEquals(unwrapJsonElement(this_), unwrapJsonElement(o))\n\t}\n\treturn false\n";
    }
    if (m.equals("hashCode") && arity == 0) return "\treturn javalang.ObjectHashCode(unwrapJsonElement(this_))\n";
    return null;
  }

  public static String buildIoNioXmlGsonRuntimeGo(RuntimeBindings rb) {
    return buildIoRtPart1() + buildIoRtPart2();
  }

  private static String buildIoRtPart1() {
    return """
        // Clean-room IO, NIO, Zip, Regex, Guava, Gson, XML/SAX, and JDK runtime helpers for j2gort.
        package j2gort

        import (
        \t"bytes"
        \t"encoding/xml"
        \t"fmt"
        \t"io"
        \t"os"
        \t"runtime"
        \t"strconv"
        \t"strings"
        \t"sync"
        \t"unicode/utf16"

        \t"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
        \t"github.com/paullepoulpe/jsgo/transpiled/shims/guava"
        \t"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
        \t"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
        \t"github.com/paullepoulpe/jsgo/transpiled/shims/javanio"
        \t"github.com/paullepoulpe/jsgo/transpiled/shims/javaregex"
        \t"github.com/paullepoulpe/jsgo/transpiled/shims/javazip"
        \tshimsres "github.com/paullepoulpe/jsgo/transpiled/shims/resources"
        )

        var (
        \t_ = bytes.NewReader
        \t_ = fmt.Sprintf
        \t_ = strconv.Itoa
        \t_ = javaregex.Pattern_Quote
        \t_ = shimsres.ResourceBundle_GetBundle
        )

        func int8ToBytes(s []int8) []byte {
        \tout := make([]byte, len(s))
        \tfor i, v := range s {
        \t\tout[i] = byte(v)
        \t}
        \treturn out
        }

        func bytesToInt8(b []byte) []int8 {
        \tout := make([]int8, len(b))
        \tfor i, v := range b {
        \t\tout[i] = int8(v)
        \t}
        \treturn out
        }

        func ioErr(err error) {
        \tif err == nil {
        \t\treturn
        \t}
        \tif t, ok := err.(javalang.Throwable); ok {
        \t\tpanic(t)
        \t}
        \tswitch err.(type) {
        \tcase *javanio.CharacterCodingException, *javanio.MalformedInputException, *javanio.UnmappableCharacterException:
        \t\tex := &Ext_java_nio_charset_CharacterCodingExceptionBase{}
        \t\tex.BindExtSelf_(ex)
        \t\tInitUserThrowable(ex, "java.nio.charset.CharacterCodingException", err.Error())
        \t\tpanic(ex)
        \t}
        \tif os.IsNotExist(err) {
        \t\tpanic(javaio.NewFileNotFoundException(err.Error()))
        \t}
        \tpanic(javaio.NewIOException(err.Error()))
        }

        // --- File & PrintStream Wrappers ---

        func unwrapFile(f *Ext_java_io_File) *javaio.File {
        \tif f == nil {
        \t\treturn nil
        \t}
        \tif jf, ok := f.Core_.(*javaio.File); ok {
        \t\treturn jf
        \t}
        \treturn nil
        }

        func wrapFile(jf *javaio.File) *Ext_java_io_File {
        \tif jf == nil {
        \t\treturn nil
        \t}
        \tf := &Ext_java_io_File{Core_: jf}
        \tf.BindExtSelf_(f)
        \treturn f
        }

        func wrapPrintStream(ps *javaio.PrintStream) *Ext_java_io_PrintStream {
        \tp := &Ext_java_io_PrintStream{}
        \tp.Core_ = ps
        \tp.BindExtSelf_(p)
        \treturn p
        }

        func unwrapPrintStream(p *Ext_java_io_PrintStream) *javaio.PrintStream {
        \tCheckNotNullPtr(p)
        \tif ps, ok := p.Core_.(*javaio.PrintStream); ok && ps != nil {
        \t\treturn ps
        \t}
        \treturn javaio.System_Out
        }

        func wrapStdinStream() Ext_java_io_InputStream {
        \treturn wrapInputStream(javaio.NewReaderInputStream(os.Stdin))
        }

        func wrapInputStream(in javaio.InputStream) Ext_java_io_InputStream {
        \tif in == nil {
        \t\treturn nil
        \t}
        \ts := &Ext_java_io_InputStreamBase{Core_: in}
        \ts.BindExtSelf_(s)
        \treturn s
        }

        func wrapOutputStream(out javaio.OutputStream) Ext_java_io_OutputStream {
        \tif out == nil {
        \t\treturn nil
        \t}
        \ts := &Ext_java_io_OutputStreamBase{Core_: out}
        \ts.BindExtSelf_(s)
        \treturn s
        }

        // --- Polymorphic InputStream / OutputStream / Reader / Writer Bridges ---

        type extInputStreamAdapter struct {
        \text   Ext_java_io_InputStream
        \tinner javaio.InputStream
        }

        func (a *extInputStreamAdapter) Read(p []byte) (int, error) {
        \treturn a.inner.Read(p)
        }
        func (a *extInputStreamAdapter) ReadByte() (byte, error) {
        \treturn a.inner.ReadByte()
        }
        func (a *extInputStreamAdapter) Skip(n int64) (int64, error) {
        \treturn a.inner.Skip(n)
        }
        func (a *extInputStreamAdapter) Available() int {
        \treturn a.inner.Available()
        }
        func (a *extInputStreamAdapter) Close() error {
        \ta.ext.Close_0()
        \treturn nil
        }

        func asInputStream(ext Ext_java_io_InputStream) javaio.InputStream {
        \tif ext == nil || javalang.IsNil(ext) {
        \t\tpanic(javalang.NewNullPointerException())
        \t}
        \tbase := ext.GetExt_java_io_InputStreamBase()
        \tvar underlying javaio.InputStream
        \tswitch c := base.Core_.(type) {
        \tcase javaio.InputStream:
        \t\tunderlying = c
        \tcase []byte:
        \t\tbais := javaio.NewByteArrayInputStream(c)
        \t\tbase.Core_ = bais
        \t\tunderlying = bais
        \tcase Ext_java_io_InputStream:
        \t\tunderlying = asInputStream(c)
        \tcase io.Reader:
        \t\trs := javaio.NewReaderInputStream(c)
        \t\tbase.Core_ = rs
        \t\tunderlying = rs
        \tdefault:
        \t\tpanic(fmt.Sprintf("UNIMPLEMENTED_EXT: unsupported InputStream Core_ %T", base.Core_))
        \t}
        \t// Always route Close() through ext.Close_0() so NonClosingInputStream is respected!
        \treturn &extInputStreamAdapter{ext: ext, inner: underlying}
        }

        type extOutputStreamAdapter struct {
        \text   Ext_java_io_OutputStream
        \tinner javaio.OutputStream
        }

        func (a *extOutputStreamAdapter) Write(p []byte) (int, error) {
        \tif a.inner != nil {
        \t\treturn a.inner.Write(p)
        \t}
        \ta.ext.Write_1_slice_int8(bytesToInt8(p))
        \treturn len(p), nil
        }
        func (a *extOutputStreamAdapter) WriteByte(c byte) error {
        \tif a.inner != nil {
        \t\treturn a.inner.WriteByte(c)
        \t}
        \ta.ext.Write_1_int32(int32(c))
        \treturn nil
        }
        func (a *extOutputStreamAdapter) Flush() error {
        \tif a.inner != nil {
        \t\treturn a.inner.Flush()
        \t}
        \ta.ext.Flush()
        \treturn nil
        }
        func (a *extOutputStreamAdapter) Close() error {
        \ta.ext.Close_0()
        \treturn nil
        }

        type jByteSinkCore interface {
        \tWrite([]byte) (int, error)
        \tWriteByte(byte) error
        }

        type jByteSinkOutputAdapter struct {
        \tsink jByteSinkCore
        }

        func (a *jByteSinkOutputAdapter) Write(p []byte) (int, error) { return a.sink.Write(p) }
        func (a *jByteSinkOutputAdapter) WriteByte(c byte) error      { return a.sink.WriteByte(c) }
        func (a *jByteSinkOutputAdapter) Flush() error                { return nil }
        func (a *jByteSinkOutputAdapter) Close() error                { return nil }

        func unwrapOutputStream(ext Ext_java_io_OutputStream) (javaio.OutputStream, bool) {
        \tbase := ext.GetExt_java_io_OutputStreamBase()
        \tif base == nil {
        \t\treturn nil, false
        \t}
        \tci, ok := any(ext).(interface{ JavaClassInfo_() (string, string) })
        \tif !ok {
        \t\treturn nil, false
        \t}
        \tbinName, _ := ci.JavaClassInfo_()
        \tswitch binName {
        \tcase "java.io.OutputStream",
        \t\t"java.io.ByteArrayOutputStream",
        \t\t"java.io.BufferedOutputStream",
        \t\t"java.io.FileOutputStream",
        \t\t"java.io.PrintStream",
        \t\t"java.util.zip.DeflaterOutputStream",
        \t\t"java.util.zip.GZIPOutputStream",
        \t\t"com.google.javascript.jscomp.serialization.FastGzipOutputStream":
        \t\tif os, ok := base.Core_.(javaio.OutputStream); ok && os != nil {
        \t\t\treturn os, true
        \t\t}
        \tcase "java.io.FilterOutputStream":
        \t\tif innerExt, ok := base.Core_.(Ext_java_io_OutputStream); ok && innerExt != nil && !javalang.IsNil(innerExt) {
        \t\t\treturn unwrapOutputStream(innerExt)
        \t\t}
        \t\tif os, ok := base.Core_.(javaio.OutputStream); ok && os != nil {
        \t\t\treturn os, true
        \t\t}
        \tcase "com.google.javascript.jscomp.Compiler$NonClosingOutputStream":
        \t\tif innerExt, ok := base.Core_.(Ext_java_io_OutputStream); ok && innerExt != nil && !javalang.IsNil(innerExt) {
        \t\t\tif os, _ := unwrapOutputStream(innerExt); os != nil {
        \t\t\t\treturn os, false
        \t\t\t}
        \t\t}
        \t\tif os, ok := base.Core_.(javaio.OutputStream); ok && os != nil {
        \t\t\treturn os, false
        \t\t}
        \tcase "com.google.protobuf.ByteString$Output":
        \t\tif os, ok := base.Core_.(javaio.OutputStream); ok && os != nil {
        \t\t\treturn os, true
        \t\t}
        \t\tif sink, ok := base.Core_.(jByteSinkCore); ok && sink != nil {
        \t\t\treturn &jByteSinkOutputAdapter{sink: sink}, true
        \t\t}
        \t}
        \treturn nil, false
        }

        func asOutputStream(ext Ext_java_io_OutputStream) javaio.OutputStream {
        \tif ext == nil || javalang.IsNil(ext) {
        \t\tpanic(javalang.NewNullPointerException())
        \t}
        \tif inner, directClose := unwrapOutputStream(ext); inner != nil {
        \t\tif directClose {
        \t\t\treturn inner
        \t\t}
        \t\treturn &extOutputStreamAdapter{ext: ext, inner: inner}
        \t}
        \treturn &extOutputStreamAdapter{ext: ext}
        }

        func newBufferedExtOutputStream(target Ext_java_io_OutputStream) javaio.OutputStream {
        \treturn javaio.NewBufferedOutputStream(asOutputStream(target))
        }

        func extInputStreamClose(core any) {
        \tswitch c := core.(type) {
        \tcase nil:
        \t\treturn
        \tcase Ext_java_io_InputStream:
        \t\tc.Close_0()
        \tcase javaio.InputStream:
        \t\tioErr(c.Close())
        \tcase io.Closer:
        \t\tioErr(c.Close())
        \t}
        }

        func extOutputStreamWriteByte(core any, b byte) {
        \tswitch c := core.(type) {
        \tcase Ext_java_io_OutputStream:
        \t\tc.Write_1_int32(int32(b))
        \tcase javaio.OutputStream:
        \t\tioErr(c.WriteByte(b))
        \tcase io.Writer:
        \t\t_, err := c.Write([]byte{b})
        \t\tioErr(err)
        \tdefault:
        \t\tpanic(fmt.Sprintf("UNIMPLEMENTED_EXT: OutputStream.write(byte) on %T", core))
        \t}
        }

        func extOutputStreamWriteSlice(core any, s []int8) {
        \tswitch c := core.(type) {
        \tcase Ext_java_io_OutputStream:
        \t\tc.Write_1_slice_int8(s)
        \tcase javaio.OutputStream:
        \t\t_, err := c.Write(int8ToBytes(s))
        \t\tioErr(err)
        \tcase io.Writer:
        \t\t_, err := c.Write(int8ToBytes(s))
        \t\tioErr(err)
        \tdefault:
        \t\tpanic(fmt.Sprintf("UNIMPLEMENTED_EXT: OutputStream.write(byte[]) on %T", core))
        \t}
        }

        func extOutputStreamFlush(core any) {
        \tswitch c := core.(type) {
        \tcase nil:
        \t\treturn
        \tcase Ext_java_io_OutputStream:
        \t\tc.Flush()
        \tcase javaio.OutputStream:
        \t\tioErr(c.Flush())
        \t}
        }

        func extOutputStreamClose(core any) {
        \tswitch c := core.(type) {
        \tcase nil:
        \t\treturn
        \tcase Ext_java_io_OutputStream:
        \t\tc.Close_0()
        \tcase javaio.OutputStream:
        \t\tioErr(c.Close())
        \tcase io.Closer:
        \t\tioErr(c.Close())
        \t}
        }

        func asReader(ext Ext_java_io_Reader) javaio.Reader {
        \tif ext == nil || javalang.IsNil(ext) {
        \t\tpanic(javalang.NewNullPointerException())
        \t}
        \tbase := ext.GetExt_java_io_ReaderBase()
        \tif r, ok := base.Core_.(javaio.Reader); ok && r != nil {
        \t\treturn r
        \t}
        \tpanic(fmt.Sprintf("UNIMPLEMENTED_EXT: unsupported Reader Core_ %T", base.Core_))
        }

        type extWriterAdapter struct {
        \text Ext_java_io_Writer
        }

        func (a *extWriterAdapter) Write(s string) {
        \ta.ext.Write_1_string(s)
        }
        func (a *extWriterAdapter) AppendCharSequence(s string) {
        \ta.ext.Write_1_string(s)
        }
        func (a *extWriterAdapter) AppendCharacter(c uint16) {
        \ta.ext.Write_1_string(javalang.StringFromChar(c))
        }
        func (a *extWriterAdapter) Flush() {
        \ta.ext.Flush()
        }
        func (a *extWriterAdapter) Close() {
        \ta.ext.Close_0()
        }

        func asWriter(ext Ext_java_io_Writer) javaio.Writer {
        \tif ext == nil || javalang.IsNil(ext) {
        \t\tpanic(javalang.NewNullPointerException())
        \t}
        \tbase := ext.GetExt_java_io_WriterBase()
        \tif w, ok := base.Core_.(javaio.Writer); ok && w != nil {
        \t\tif base.Self_ == nil || RefEquals(base.Self_, ext) {
        \t\t\treturn w
        \t\t}
        \t}
        \treturn &extWriterAdapter{ext: ext}
        }

        func extReaderReadOne(core any) int32 {
        \tr, ok := core.(javaio.Reader)
        \tif !ok || r == nil {
        \t\treturn -1
        \t}
        \tvar buf [1]uint16
        \tn := r.Read(buf[:], 0, 1)
        \tif n <= 0 {
        \t\treturn -1
        \t}
        \treturn int32(buf[0])
        }

        func extBufferedReaderReadLine(core any) string {
        \tbr, ok := core.(*javaio.BufferedReader)
        \tif !ok || br == nil {
        \t\treturn NullString
        \t}
        \tline, err := br.ReadLine()
        \tif err == io.EOF {
        \t\treturn NullString
        \t}
        \tioErr(err)
        \treturn line
        }

        func extReaderClose(core any) {
        \tif r, ok := core.(javaio.Reader); ok && r != nil {
        \t\tioErr(r.Close())
        \t}
        }

        func extWriterWriteString(core any, s string) {
        \tw, ok := core.(javaio.Writer)
        \tif !ok || w == nil {
        \t\tpanic("UNIMPLEMENTED_EXT: Writer Core_ is nil")
        \t}
        \tw.Write(s)
        }

        func extWriterFlush(core any) {
        \tif w, ok := core.(javaio.Writer); ok && w != nil {
        \t\tw.Flush()
        \t}
        }

        func extWriterClose(core any) {
        \tif w, ok := core.(javaio.Writer); ok && w != nil {
        \t\tw.Close()
        \t}
        }

        func newInputStreamReaderCharset(in Ext_java_io_InputStream, cs Ext_java_nio_charset_Charset) javaio.Reader {
        \treturn javanio.NewInputStreamReaderWithCharset(asInputStream(in), unwrapCharset(cs))
        }

        func newInputStreamReaderCharsetName(in Ext_java_io_InputStream, csName string) javaio.Reader {
        \treturn javanio.NewInputStreamReaderWithCharset(asInputStream(in), unwrapCharset(charsetForName(csName)))
        }

        func newOutputStreamWriterCharset(out Ext_java_io_OutputStream, cs Ext_java_nio_charset_Charset) javaio.Writer {
        \treturn javanio.NewOutputStreamWriterWithCharset(asOutputStream(out), unwrapCharset(cs))
        }

        func newOutputStreamWriterCharsetName(out Ext_java_io_OutputStream, csName string) javaio.Writer {
        \treturn javanio.NewOutputStreamWriterWithCharset(asOutputStream(out), unwrapCharset(charsetForName(csName)))
        }

        func newOutputStreamWriterEncoder(out Ext_java_io_OutputStream, enc Ext_java_nio_charset_CharsetEncoder) javaio.Writer {
        \treturn javanio.NewOutputStreamWriterWithEncoder(asOutputStream(out), unwrapCharsetEncoder(CheckNotNull(enc)))
        }

        // --- Zip & GZIP Helpers ---

        func newGzipInputStream(in Ext_java_io_InputStream) *javazip.GZIPInputStream {
        \tgz, err := javazip.NewGZIPInputStream(asInputStream(in))
        \tioErr(err)
        \treturn gz
        }

        type gzipDeflaterBridge struct {
        \tout   io.Writer
        \tgz    *javazip.GZIPOutputStream
        \towner *Ext_java_util_zip_DeflaterOutputStreamBase
        }

        func initGzipOutputStream(base *Ext_java_util_zip_DeflaterOutputStreamBase, out Ext_java_io_OutputStream) {
        \tw := asOutputStream(CheckNotNull(out))
        \tgz, err := javazip.NewGZIPOutputStream(w)
        \tioErr(err)
        \tbase.Core_ = gz
        \tbridge := &gzipDeflaterBridge{out: w, gz: gz, owner: base}
        \tbase.F_Def = &Ext_java_util_zip_Deflater{Core_: bridge}
        }

        func zipDeflaterSetLevel(core any, level int) {
        \tif b, ok := core.(*gzipDeflaterBridge); ok && b != nil {
        \t\tgz, err := javazip.NewGZIPOutputStreamLevel(b.out, level)
        \t\tioErr(err)
        \t\tb.gz = gz
        \t\tb.owner.Core_ = gz
        \t}
        }

        func gzipOutputStreamFinish(core any) {
        \tif gz, ok := core.(*javazip.GZIPOutputStream); ok && gz != nil {
        \t\tioErr(gz.Finish())
        \t}
        }

        func newZipFileFromFile(f *Ext_java_io_File) *javazip.ZipFile {
        \tzf, err := javazip.NewZipFileFromFile(CheckNotNullPtr(unwrapFile(f)))
        \tioErr(err)
        \treturn zf
        }

        func newZipFileFromName(name string) *javazip.ZipFile {
        \tzf, err := javazip.NewZipFile(name)
        \tioErr(err)
        \treturn zf
        }

        func newZipInputStream(in Ext_java_io_InputStream, cs Ext_java_nio_charset_Charset) *javazip.ZipInputStream {
        \tif cs == nil || javalang.IsNil(cs) {
        \t\treturn javazip.NewZipInputStream(asInputStream(in))
        \t}
        \treturn javazip.NewZipInputStreamCharset(asInputStream(in), unwrapCharset(cs))
        }

        // --- NIO Charset, Path, Files, FileSystem, URI ---

        func wrapCharset(cs *javanio.Charset) Ext_java_nio_charset_Charset {
        \tif cs == nil {
        \t\treturn nil
        \t}
        \tc := &Ext_java_nio_charset_CharsetBase{Core_: cs}
        \tc.BindExtSelf_(c)
        \treturn c
        }

        func wrapCharsetByName(name string) Ext_java_nio_charset_Charset {
        \treturn wrapCharset(javanio.Charset_ForName(name))
        }

        func unwrapCharset(cs Ext_java_nio_charset_Charset) *javanio.Charset {
        \tif cs == nil || javalang.IsNil(cs) {
        \t\treturn javanio.Charset_DefaultCharset()
        \t}
        \tif c, ok := cs.GetExt_java_nio_charset_CharsetBase().Core_.(*javanio.Charset); ok && c != nil {
        \t\treturn c
        \t}
        \treturn javanio.Charset_DefaultCharset()
        }

        func NewStringFromBytes(b []int8, off, length int32, csArg any) string {
        \tif b == nil {
        \t\tThrowNPE()
        \t}
        \tif off < 0 || length < 0 || int(off) > len(b)-int(length) {
        \t\tpanic(javalang.NewStringIndexOutOfBoundsException("String index out of range"))
        \t}
        \tsub := b[off : off+length]
        \traw := make([]byte, len(sub))
        \tfor i, v := range sub {
        \t\traw[i] = byte(v)
        \t}
        \tif csArg == nil {
        \t\treturn javanio.Charset_DefaultCharset().Decode(raw)
        \t}
        \tswitch x := csArg.(type) {
        \tcase Ext_java_nio_charset_Charset:
        \t\tif javalang.IsNil(x) {
        \t\t\tThrowNPE()
        \t\t}
        \t\treturn unwrapCharset(x).Decode(raw)
        \tcase string:
        \t\tif x == NullString {
        \t\t\tThrowNPE()
        \t\t}
        \t\treturn unwrapCharset(charsetForName(x)).Decode(raw)
        \tdefault:
        \t\ts := CastString(csArg)
        \t\tif s == NullString {
        \t\t\tThrowNPE()
        \t\t}
        \t\treturn unwrapCharset(charsetForName(s)).Decode(raw)
        \t}
        }

        func charsetForName(name string) Ext_java_nio_charset_Charset {
        \tvar cs *javanio.Charset
        \tfunc() {
        \t\tdefer func() {
        \t\t\tif r := recover(); r != nil {
        \t\t\t\tswitch e := r.(type) {
        \t\t\t\tcase *javanio.IllegalCharsetNameException:
        \t\t\t\t\tpanic(javalang.NewIllegalArgumentException(e.Error()))
        \t\t\t\tcase *javanio.UnsupportedCharsetException:
        \t\t\t\t\tpanic(javalang.NewIllegalArgumentException(e.Error()))
        \t\t\t\tdefault:
        \t\t\t\t\tpanic(r)
        \t\t\t\t}
        \t\t\t}
        \t\t}()
        \t\tcs = javanio.Charset_ForName(name)
        \t}()
        \treturn wrapCharset(cs)
        }

        func charsetIsSupported(name string) bool {
        \tvar ok bool
        \tfunc() {
        \t\tdefer func() {
        \t\t\tif r := recover(); r != nil {
        \t\t\t\tswitch e := r.(type) {
        \t\t\t\tcase *javanio.IllegalCharsetNameException:
        \t\t\t\t\tpanic(javalang.NewIllegalArgumentException(e.Error()))
        \t\t\t\tdefault:
        \t\t\t\t\tpanic(r)
        \t\t\t\t}
        \t\t\t}
        \t\t}()
        \t\tok = javanio.Charset_IsSupported(name)
        \t}()
        \treturn ok
        }

        func newCodingErrorAction(name string) *Ext_java_nio_charset_CodingErrorAction {
        \tact := javanio.CodingErrorAction_REPLACE
        \tswitch name {
        \tcase "IGNORE":
        \t\tact = javanio.CodingErrorAction_IGNORE
        \tcase "REPORT":
        \t\tact = javanio.CodingErrorAction_REPORT
        \t}
        \treturn &Ext_java_nio_charset_CodingErrorAction{Core_: act}
        }

        func unwrapCodingAction(a *Ext_java_nio_charset_CodingErrorAction) javanio.CodingErrorAction {
        \tCheckNotNullPtr(a)
        \tif act, ok := a.Core_.(javanio.CodingErrorAction); ok {
        \t\treturn act
        \t}
        \treturn javanio.CodingErrorAction_REPLACE
        }

        func wrapCharsetEncoder(enc *javanio.CharsetEncoder) Ext_java_nio_charset_CharsetEncoder {
        \te := &Ext_java_nio_charset_CharsetEncoderBase{Core_: enc}
        \te.BindExtSelf_(e)
        \treturn e
        }

        func unwrapCharsetEncoder(enc Ext_java_nio_charset_CharsetEncoder) *javanio.CharsetEncoder {
        \treturn enc.GetExt_java_nio_charset_CharsetEncoderBase().Core_.(*javanio.CharsetEncoder)
        }

        type jNioPath struct {
        \tp *javanio.Path
        }

        func (*jNioPath) Implements_j2gort_Ext_java_nio_file_Path()      {}
        func (*jNioPath) Implements_j2gort_Ext_java_lang_Iterable()      {}
        func (*jNioPath) Implements_j2gort_Ext_java_nio_file_Watchable() {}

        func wrapPath(p *javanio.Path) Ext_java_nio_file_Path {
        \tif p == nil {
        \t\treturn nil
        \t}
        \treturn &jNioPath{p: p}
        }

        func unwrapPath(ep Ext_java_nio_file_Path) *javanio.Path {
        \tCheckNotNull(ep)
        \treturn ep.(*jNioPath).p
        }

        func newExtPath(first string, more []string) Ext_java_nio_file_Path {
        \tfor _, m := range more {
        \t\tCheckNotNullString(m)
        \t}
        \tvar p *javanio.Path
        \tfunc() {
        \t\tdefer func() {
        \t\t\tif r := recover(); r != nil {
        \t\t\t\tif ipe, ok := r.(*javanio.InvalidPathException); ok {
        \t\t\t\t\tpanic(javalang.NewIllegalArgumentException(ipe.Error()))
        \t\t\t\t}
        \t\t\t\tpanic(r)
        \t\t\t}
        \t\t}()
        \t\tp = javanio.Path_Of(first, more...)
        \t}()
        \treturn wrapPath(p)
        }

        func (jp *jNioPath) GetParent() Ext_java_nio_file_Path {
        \treturn wrapPath(jp.p.GetParent())
        }
        func (jp *jNioPath) Normalize() Ext_java_nio_file_Path {
        \treturn wrapPath(jp.p.Normalize())
        }
        func (jp *jNioPath) Resolve_1_Ext_java_nio_file_Path(other Ext_java_nio_file_Path) Ext_java_nio_file_Path {
        \treturn wrapPath(jp.p.Resolve(unwrapPath(other)))
        }
        func (jp *jNioPath) Resolve_1_string(other string) Ext_java_nio_file_Path {
        \treturn wrapPath(jp.p.ResolveString(CheckNotNullString(other)))
        }
        func (jp *jNioPath) ResolveSibling_1_string(other string) Ext_java_nio_file_Path {
        \treturn wrapPath(jp.p.ResolveSiblingString(CheckNotNullString(other)))
        }
        func (jp *jNioPath) Relativize(other Ext_java_nio_file_Path) Ext_java_nio_file_Path {
        \treturn wrapPath(jp.p.Relativize(unwrapPath(other)))
        }
        func (jp *jNioPath) ToAbsolutePath() Ext_java_nio_file_Path {
        \treturn wrapPath(jp.p.ToAbsolutePath())
        }
        func (jp *jNioPath) ToFile() *Ext_java_io_File {
        \treturn wrapFile(jp.p.ToFile())
        }
        func (jp *jNioPath) CompareTo(other any) int32 {
        \to := unwrapPath(AsIface[Ext_java_nio_file_Path](CheckNotNull(other)))
        \treturn int32(jp.p.CompareTo(o))
        }
        func (jp *jNioPath) Equals(other any) bool {
        \tif o, ok := other.(*jNioPath); ok && o != nil {
        \t\treturn jp.p.Equals(o.p)
        \t}
        \treturn false
        }
        func (jp *jNioPath) HashCode() int32 { return jp.p.HashCode() }
        func (jp *jNioPath) String() string  { return jp.p.String() }
        func (jp *jNioPath) Iterator() any {
        \tn := jp.p.GetNameCount()
        \titems := make([]any, n)
        \tfor i := 0; i < n; i++ {
        \t\titems[i] = wrapPath(jp.p.GetName(i))
        \t}
        \treturn newSliceObjIter(items)
        }
        func (jp *jNioPath) ForEach(action Ext_java_util_function_Consumer) {
        \tExtDefault_Iterable_ForEach(jp, action)
        }

        type jPathMatcher struct {
        \tpm javanio.PathMatcher
        }

        func (*jPathMatcher) Implements_j2gort_Ext_java_nio_file_PathMatcher() {}
        func (m *jPathMatcher) Matches(p Ext_java_nio_file_Path) bool {
        \treturn m.pm.Matches(unwrapPath(p))
        }

        func newExtPathMatcher(syntaxAndPattern string) Ext_java_nio_file_PathMatcher {
        \treturn &jPathMatcher{pm: javanio.FileSystems_GetDefault().GetPathMatcher(syntaxAndPattern)}
        }

        func newDefaultFileSystem() Ext_java_nio_file_FileSystem {
        \tfs := &Ext_java_nio_file_FileSystemBase{Core_: javanio.FileSystems_GetDefault()}
        \tfs.BindExtSelf_(fs)
        \treturn fs
        }

        func convertOpenOptions(opts []Ext_java_nio_file_OpenOption) []javanio.StandardOpenOption {
        \tvar out []javanio.StandardOpenOption
        \tfor _, o := range opts {
        \t\tif so, ok := o.(*Ext_java_nio_file_StandardOpenOption); ok && so != nil {
        \t\t\tswitch so.Name_ {
        \t\t\tcase "READ":
        \t\t\t\tout = append(out, javanio.StandardOpenOption_READ)
        \t\t\tcase "WRITE":
        \t\t\t\tout = append(out, javanio.StandardOpenOption_WRITE)
        \t\t\tcase "APPEND":
        \t\t\t\tout = append(out, javanio.StandardOpenOption_APPEND)
        \t\t\tcase "TRUNCATE_EXISTING":
        \t\t\t\tout = append(out, javanio.StandardOpenOption_TRUNCATE_EXISTING)
        \t\t\tcase "CREATE":
        \t\t\t\tout = append(out, javanio.StandardOpenOption_CREATE)
        \t\t\tcase "CREATE_NEW":
        \t\t\t\tout = append(out, javanio.StandardOpenOption_CREATE_NEW)
        \t\t\t}
        \t\t}
        \t}
        \treturn out
        }

        func nioFilesCreateDirectories(dir Ext_java_nio_file_Path) Ext_java_nio_file_Path {
        \tp, err := javanio.Files_CreateDirectories(unwrapPath(dir))
        \tioErr(err)
        \treturn wrapPath(p)
        }

        func nioFilesGetLastModifiedTime(p Ext_java_nio_file_Path) *Ext_java_nio_file_attribute_FileTime {
        \tft, err := javanio.Files_GetLastModifiedTime(unwrapPath(p))
        \tioErr(err)
        \treturn &Ext_java_nio_file_attribute_FileTime{Core_: ft}
        }

        func nioFilesNewBufferedReader(p Ext_java_nio_file_Path, cs Ext_java_nio_charset_Charset) *Ext_java_io_BufferedReader {
        \tbr, err := javanio.Files_NewBufferedReader(unwrapPath(p), unwrapCharset(cs))
        \tioErr(err)
        \tr := &Ext_java_io_BufferedReader{}
        \tr.Core_ = br
        \tr.BindExtSelf_(r)
        \treturn r
        }

        func nioFilesNewBufferedWriter(p Ext_java_nio_file_Path, cs Ext_java_nio_charset_Charset, opts []Ext_java_nio_file_OpenOption) *Ext_java_io_BufferedWriter {
        \tbw, err := javanio.Files_NewBufferedWriter(unwrapPath(p), unwrapCharset(cs), convertOpenOptions(opts)...)
        \tioErr(err)
        \tw := &Ext_java_io_BufferedWriter{}
        \tw.Core_ = bw
        \tw.BindExtSelf_(w)
        \treturn w
        }

        func nioFilesNewOutputStream(p Ext_java_nio_file_Path, opts []Ext_java_nio_file_OpenOption) Ext_java_io_OutputStream {
        \tfos, err := javanio.Files_NewOutputStream(unwrapPath(p), convertOpenOptions(opts)...)
        \tioErr(err)
        \treturn wrapOutputStream(fos)
        }

        func nioFilesReadString(p Ext_java_nio_file_Path, cs Ext_java_nio_charset_Charset) string {
        \ts, err := javanio.Files_ReadStringCharset(unwrapPath(p), unwrapCharset(cs))
        \tioErr(err)
        \treturn s
        }

        type jBasicFileAttrs struct {
        \ta *javanio.BasicFileAttributes
        }

        func (*jBasicFileAttrs) Implements_j2gort_Ext_java_nio_file_attribute_BasicFileAttributes() {}

        type jFileVisitorAdapter struct {
        \tv Ext_java_nio_file_FileVisitor
        }

        func mapVisitResult(r *Ext_java_nio_file_FileVisitResult) javanio.FileVisitResult {
        \tif r == nil {
        \t\treturn javanio.FileVisitResult_CONTINUE
        \t}
        \tswitch r.Name_ {
        \tcase "TERMINATE":
        \t\treturn javanio.FileVisitResult_TERMINATE
        \tcase "SKIP_SUBTREE":
        \t\treturn javanio.FileVisitResult_SKIP_SUBTREE
        \tcase "SKIP_SIBLINGS":
        \t\treturn javanio.FileVisitResult_SKIP_SIBLINGS
        \tdefault:
        \t\treturn javanio.FileVisitResult_CONTINUE
        \t}
        }

        func (a *jFileVisitorAdapter) PreVisitDirectory(_ *javanio.Path, _ *javanio.BasicFileAttributes) (javanio.FileVisitResult, error) {
        \treturn javanio.FileVisitResult_CONTINUE, nil
        }

        func (a *jFileVisitorAdapter) VisitFile(file *javanio.Path, attrs *javanio.BasicFileAttributes) (javanio.FileVisitResult, error) {
        \tr := a.v.VisitFile(wrapPath(file), &jBasicFileAttrs{a: attrs})
        \treturn mapVisitResult(r), nil
        }

        func (a *jFileVisitorAdapter) VisitFileFailed(file *javanio.Path, exc error) (javanio.FileVisitResult, error) {
        \tvar ioEx javaio.IIOException
        \tif e, ok := exc.(javaio.IIOException); ok {
        \t\tioEx = e
        \t} else if exc != nil {
        \t\tioEx = javaio.NewIOException(exc.Error())
        \t}
        \tr := a.v.VisitFileFailed(wrapPath(file), ioEx)
        \treturn mapVisitResult(r), nil
        }

        func (a *jFileVisitorAdapter) PostVisitDirectory(_ *javanio.Path, exc error) (javanio.FileVisitResult, error) {
        \tif exc != nil {
        \t\treturn javanio.FileVisitResult_TERMINATE, exc
        \t}
        \treturn javanio.FileVisitResult_CONTINUE, nil
        }

        func nioFilesWalkFileTree(start Ext_java_nio_file_Path, options Ext_java_util_Set, maxDepth int32, visitor Ext_java_nio_file_FileVisitor) Ext_java_nio_file_Path {
        \tfollow := false
        \tif options != nil && !javalang.IsNil(options) {
        \t\tfor it := Iterate(options); it.HasNext(); {
        \t\t\tif opt, ok := it.Next().(*Ext_java_nio_file_FileVisitOption); ok && opt != nil && opt.Name_ == "FOLLOW_LINKS" {
        \t\t\t\tfollow = true
        \t\t\t}
        \t\t}
        \t}
        \t_, err := javanio.Files_WalkFileTree(unwrapPath(start), follow, int(maxDepth), &jFileVisitorAdapter{v: visitor})
        \tioErr(err)
        \treturn start
        }

        func newJavaURI(str string) *javanio.URI {
        \tu, err := javanio.NewURI(str)
        \tif err != nil {
        \t\tex := &Ext_java_net_URISyntaxExceptionBase{}
        \t\tex.BindExtSelf_(ex)
        \t\tInitUserThrowable(ex, "java.net.URISyntaxException", err.Error())
        \t\tpanic(ex)
        \t}
        \treturn u
        }

        // --- Guava Base, Escape, Hash, IO Helpers ---

        func guavaAsciiEqualsIgnoreCase(a, b any) bool {
        \tif a == b {
        \t\treturn true
        \t}
        \tif a == nil || b == nil || javalang.IsNil(a) || javalang.IsNil(b) {
        \t\treturn false
        \t}
        \tsa := StrVal(a)
        \tsb := StrVal(b)
        \tif sa == NullString || sb == NullString {
        \t\treturn sa == sb
        \t}
        \treturn guava.Ascii_EqualsIgnoreCase(sa, sb)
        }

        func wrapCharMatcher(cm *guava.CharMatcher) Ext_com_google_common_base_CharMatcher {
        \tobj := &Ext_com_google_common_base_CharMatcherBase{Core_: cm}
        \tobj.BindExtSelf_(obj)
        \treturn obj
        }

        func unwrapCharMatcher(cm Ext_com_google_common_base_CharMatcher) *guava.CharMatcher {
        \treturn cm.GetExt_com_google_common_base_CharMatcherBase().Core_.(*guava.CharMatcher)
        }

        func wrapCharMatcherAnyOf(seq any) Ext_com_google_common_base_CharMatcher {
        \treturn wrapCharMatcher(guava.CharMatcher_AnyOf(StrVal(CheckNotNull(seq))))
        }
        func wrapCharMatcherInRange(lo, hi uint16) Ext_com_google_common_base_CharMatcher {
        \treturn wrapCharMatcher(guava.CharMatcher_InRange(lo, hi))
        }
        func wrapCharMatcherIs(ch uint16) Ext_com_google_common_base_CharMatcher {
        \treturn wrapCharMatcher(guava.CharMatcher_Is(ch))
        }
        func wrapCharMatcherIsNot(ch uint16) Ext_com_google_common_base_CharMatcher {
        \treturn wrapCharMatcher(guava.CharMatcher_IsNot(ch))
        }
        func wrapCharMatcherJavaLetterOrDigit() Ext_com_google_common_base_CharMatcher {
        \treturn wrapCharMatcher(guava.CharMatcher_JavaLetterOrDigit())
        }
        func wrapCharMatcherWhitespace() Ext_com_google_common_base_CharMatcher {
        \treturn wrapCharMatcher(guava.CharMatcher_Whitespace())
        }
        func wrapCharMatcherAscii() Ext_com_google_common_base_CharMatcher {
        \treturn wrapCharMatcher(guava.CharMatcher_Ascii())
        }

        func charMatcherOr(a *Ext_com_google_common_base_CharMatcherBase, b Ext_com_google_common_base_CharMatcher) Ext_com_google_common_base_CharMatcher {
        \treturn wrapCharMatcher(a.Core_.(*guava.CharMatcher).Or(unwrapCharMatcher(b)))
        }
        func charMatcherMatchesAllOf(a *Ext_com_google_common_base_CharMatcherBase, seq any) bool {
        \treturn a.Core_.(*guava.CharMatcher).MatchesAllOf(StrVal(seq))
        }
        func charMatcherLastIndexIn(a *Ext_com_google_common_base_CharMatcherBase, seq any) int32 {
        \treturn int32(a.Core_.(*guava.CharMatcher).LastIndexIn(StrVal(seq)))
        }
        func charMatcherTrimFrom(a *Ext_com_google_common_base_CharMatcherBase, seq any) string {
        \treturn a.Core_.(*guava.CharMatcher).TrimFrom(StrVal(seq))
        }
        func charMatcherCollapseFrom(a *Ext_com_google_common_base_CharMatcherBase, seq any, repl uint16) string {
        \treturn a.Core_.(*guava.CharMatcher).CollapseFrom(StrVal(seq), repl)
        }
        func charMatcherMatchesBoxed(a *Ext_com_google_common_base_CharMatcherBase, ch any) bool {
        \tc := UnboxChar(Cast[*BoxChar](ch))
        \treturn a.Core_.(*guava.CharMatcher).Matches(c)
        }

        type jGuavaFunc struct {
        \tfn func(any) any
        }

        func (*jGuavaFunc) Implements_j2gort_Ext_com_google_common_base_Function() {}
        func (*jGuavaFunc) Implements_j2gort_Ext_java_util_function_Function()     {}
        func (f *jGuavaFunc) Apply(x any) any                                      { return f.fn(x) }
        func (f *jGuavaFunc) Equals(o any) bool                                    { return RefEquals(f, o) }
        func (f *jGuavaFunc) AndThen(after Ext_java_util_function_Function) Ext_java_util_function_Function {
        \treturn ExtDefault_Function_AndThen(f, after)
        }
        func (f *jGuavaFunc) Compose(before Ext_java_util_function_Function) Ext_java_util_function_Function {
        \treturn ExtDefault_Function_Compose(f, before)
        }

        func newGuavaIdentityFunction() Ext_com_google_common_base_Function {
        \treturn &jGuavaFunc{fn: func(x any) any { return x }}
        }

        type jGuavaPred struct {
        \tfn     func(any) bool
        \tnodeFn func(uint32) bool
        }

        var (
        \tguavaAlwaysTruePred_  Ext_com_google_common_base_Predicate = &jGuavaPred{fn: func(_ any) bool { return true }, nodeFn: func(_ uint32) bool { return true }}
        \tguavaAlwaysFalsePred_ Ext_com_google_common_base_Predicate = &jGuavaPred{fn: func(_ any) bool { return false }, nodeFn: func(_ uint32) bool { return false }}
        )

        func (*jGuavaPred) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (*jGuavaPred) Implements_j2gort_Ext_java_util_function_Predicate()     {}
        func (*jGuavaPred) HasIdentityEquals_()                                     {}
        func (p *jGuavaPred) Apply(x any) bool                                      { return p.fn(x) }
        func (p *jGuavaPred) Test(x any) bool                                       { return p.fn(x) }
        func (p *jGuavaPred) TestNode_(n uint32) bool {
        \tif p == nil {
        \t\tThrowNPE()
        \t}
        \tif p.nodeFn != nil {
        \t\treturn p.nodeFn(n)
        \t}
        \treturn p.fn(boxNodeById(n))
        }
        func (p *jGuavaPred) Equals(o any) bool                                     { return RefEquals(p, o) }
        func (p *jGuavaPred) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \treturn ExtDefault_Predicate_And(p, o)
        }
        func (p *jGuavaPred) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \treturn ExtDefault_Predicate_Or(p, o)
        }
        func (p *jGuavaPred) Negate() Ext_java_util_function_Predicate {
        \treturn ExtDefault_Predicate_Negate(p)
        }

        func newGuavaPredicate(fn func(any) bool) Ext_com_google_common_base_Predicate {
        \treturn &jGuavaPred{fn: fn}
        }

        func newGuavaNotPredicate(target Ext_com_google_common_base_Predicate) Ext_com_google_common_base_Predicate {
        \treturn &jGuavaPred{
        \t\tfn:     func(x any) bool { return !target.Apply(x) },
        \t\tnodeFn: func(n uint32) bool { return !PredTestNode(target, n) },
        \t}
        }

        type jMemoSupplier struct {
        \tonce     sync.Once
        \tdelegate Ext_com_google_common_base_Supplier
        \tval      any
        }

        func (*jMemoSupplier) Implements_j2gort_Ext_com_google_common_base_Supplier() {}
        func (*jMemoSupplier) Implements_j2gort_Ext_java_util_function_Supplier()     {}
        func (s *jMemoSupplier) Get() any {
        \ts.once.Do(func() { s.val = s.delegate.Get() })
        \treturn s.val
        }

        func newMemoizingSupplier(d Ext_com_google_common_base_Supplier) Ext_com_google_common_base_Supplier {
        \treturn &jMemoSupplier{delegate: d}
        }

        func newGuavaJoiner(sep string) *Ext_com_google_common_base_Joiner {
        \tj := &Ext_com_google_common_base_Joiner{Core_: guava.Joiner_On(sep)}
        \tj.BindExtSelf_(j)
        \treturn j
        }

        func guavaJoinerJoinIter(j *Ext_com_google_common_base_Joiner, iter Ext_java_lang_Iterable) string {
        \tif ssi, ok := any(iter).(StrSliceIterable); ok {
        \t\tif ss, ok2 := ssi.IterStrSlice_(); ok2 {
        \t\t\tfor _, s := range ss {
        \t\t\t\tCheckNotNullString(s)
        \t\t\t}
        \t\t\treturn j.Core_.(*guava.Joiner).JoinStrings(ss)
        \t\t}
        \t}
        \tvar parts []string
        \tfor it := Iterate(iter); it.HasNext(); {
        \t\tv := it.Next()
        \t\tif v == nil || javalang.IsNil(v) {
        \t\t\tpanic(javalang.NewNullPointerException())
        \t\t}
        \t\tparts = append(parts, StrVal(v))
        \t}
        \treturn j.Core_.(*guava.Joiner).JoinStrings(parts)
        }

        func guavaJoinerJoinSlice(j *Ext_com_google_common_base_Joiner, parts []any) string {
        \tss := make([]string, len(parts))
        \tfor i, v := range parts {
        \t\tif v == nil || javalang.IsNil(v) {
        \t\t\tpanic(javalang.NewNullPointerException())
        \t\t}
        \t\tss[i] = StrVal(v)
        \t}
        \treturn j.Core_.(*guava.Joiner).JoinStrings(ss)
        }

        func guavaJoinerJoinVarargs(j *Ext_com_google_common_base_Joiner, first, second any, rest []any) string {
        \tall := make([]any, 0, 2+len(rest))
        \tall = append(all, first, second)
        \tall = append(all, rest...)
        \treturn guavaJoinerJoinSlice(j, all)
        }

        func newGuavaSplitterChar(sep uint16) *Ext_com_google_common_base_Splitter {
        \ts := &Ext_com_google_common_base_Splitter{Core_: guava.Splitter_OnChar(sep)}
        \ts.BindExtSelf_(s)
        \treturn s
        }

        func newGuavaSplitterStr(sep string) *Ext_com_google_common_base_Splitter {
        \ts := &Ext_com_google_common_base_Splitter{Core_: guava.Splitter_On(sep)}
        \ts.BindExtSelf_(s)
        \treturn s
        }

        func splitterOmitEmpty(s *Ext_com_google_common_base_Splitter) *Ext_com_google_common_base_Splitter {
        \tout := &Ext_com_google_common_base_Splitter{Core_: s.Core_.(*guava.Splitter).OmitEmptyStrings()}
        \tout.BindExtSelf_(out)
        \treturn out
        }

        func splitterTrimResults(s *Ext_com_google_common_base_Splitter) *Ext_com_google_common_base_Splitter {
        \tout := &Ext_com_google_common_base_Splitter{Core_: s.Core_.(*guava.Splitter).TrimResults()}
        \tout.BindExtSelf_(out)
        \treturn out
        }

        func splitterLimit(s *Ext_com_google_common_base_Splitter, lim int32) *Ext_com_google_common_base_Splitter {
        \tout := &Ext_com_google_common_base_Splitter{Core_: s.Core_.(*guava.Splitter).Limit(int(lim))}
        \tout.BindExtSelf_(out)
        \treturn out
        }

        func SplitterSplitToListStr(s *Ext_com_google_common_base_Splitter, seq string) Ext_java_util_List {
        \tCheckNotNullPtr(s)
        \tCheckNotNullString(seq)
        \tparts := s.Core_.(*guava.Splitter).SplitToList(seq)
        \treturn wrapImmStrSliceOwn(parts)
        }

        func SplitterSplitStr(s *Ext_com_google_common_base_Splitter, seq string) Ext_java_lang_Iterable {
        \treturn SplitterSplitToListStr(s, seq)
        }

        func splitterSplitToList(s *Ext_com_google_common_base_Splitter, seq any) Ext_java_util_List {
        \treturn SplitterSplitToListStr(s, StrVal(seq))
        }

        func splitterSplit(s *Ext_com_google_common_base_Splitter, seq any) Ext_java_lang_Iterable {
        \treturn SplitterSplitToListStr(s, StrVal(seq))
        }

        type jToStringHelperData struct {
        \tname      string
        \tomitNulls bool
        \tfields    []struct {
        \t\tk string
        \t\tv any
        \t}
        }

        func newToStringHelper(name string) *Ext_com_google_common_base_MoreObjects_ToStringHelper {
        \th := &Ext_com_google_common_base_MoreObjects_ToStringHelper{Core_: &jToStringHelperData{name: name}}
        \th.BindExtSelf_(h)
        \treturn h
        }

        func (h *Ext_com_google_common_base_MoreObjects_ToStringHelper) toStringHelperOmitNulls() {
        \th.Core_.(*jToStringHelperData).omitNulls = true
        }

        func (h *Ext_com_google_common_base_MoreObjects_ToStringHelper) toStringHelperAdd(name string, val any) {
        \td := h.Core_.(*jToStringHelperData)
        \td.fields = append(d.fields, struct {
        \t\tk string
        \t\tv any
        \t}{k: name, v: ToAny(val)})
        }

        func (h *Ext_com_google_common_base_MoreObjects_ToStringHelper) toStringHelperBuild() string {
        \td := h.Core_.(*jToStringHelperData)
        \tvar sb strings.Builder
        \tsb.WriteString(d.name)
        \tsb.WriteByte('{')
        \tfirst := true
        \tfor _, f := range d.fields {
        \t\tif d.omitNulls && f.v == nil {
        \t\t\tcontinue
        \t\t}
        \t\tif !first {
        \t\t\tsb.WriteString(", ")
        \t\t}
        \t\tfirst = false
        \t\tsb.WriteString(f.k)
        \t\tsb.WriteByte('=')
        \t\tsb.WriteString(StrVal(f.v))
        \t}
        \tsb.WriteByte('}')
        \treturn sb.String()
        }

        func guavaPadStart(s string, minLen int32, padChar uint16) string {
        \treturn guava.Strings_PadStart(s, int(minLen), padChar)
        }

        func splitCaseWords(format, s string) []string {
        \tswitch format {
        \tcase "LOWER_HYPHEN":
        \t\treturn strings.Split(s, "-")
        \tcase "LOWER_UNDERSCORE", "UPPER_UNDERSCORE":
        \t\treturn strings.Split(s, "_")
        \tcase "LOWER_CAMEL", "UPPER_CAMEL":
        \t\tif s == "" {
        \t\t\treturn []string{""}
        \t\t}
        \t\tvar words []string
        \t\tstart := 0
        \t\tfor i := 1; i < len(s); i++ {
        \t\t\tif s[i] >= 'A' && s[i] <= 'Z' {
        \t\t\t\twords = append(words, s[start:i])
        \t\t\t\tstart = i
        \t\t\t}
        \t\t}
        \t\twords = append(words, s[start:])
        \t\treturn words
        \tdefault:
        \t\treturn []string{s}
        \t}
        }

        func capAsciiWord(w string) string {
        \tif w == "" {
        \t\treturn ""
        \t}
        \tlow := javalang.StringToLowerCaseASCII(w)
        \tfirst := low[0]
        \tif first >= 'a' && first <= 'z' {
        \t\tfirst = first - 'a' + 'A'
        \t}
        \treturn string(first) + low[1:]
        }

        func guavaCaseFormatConvert(srcFmt, dstFmt, s string) string {
        \tif srcFmt == dstFmt || s == "" {
        \t\treturn s
        \t}
        \twords := splitCaseWords(srcFmt, s)
        \tswitch dstFmt {
        \tcase "LOWER_HYPHEN":
        \t\tout := make([]string, len(words))
        \t\tfor i, w := range words {
        \t\t\tout[i] = javalang.StringToLowerCaseASCII(w)
        \t\t}
        \t\treturn strings.Join(out, "-")
        \tcase "LOWER_UNDERSCORE":
        \t\tout := make([]string, len(words))
        \t\tfor i, w := range words {
        \t\t\tout[i] = javalang.StringToLowerCaseASCII(w)
        \t\t}
        \t\treturn strings.Join(out, "_")
        \tcase "UPPER_UNDERSCORE":
        \t\tout := make([]string, len(words))
        \t\tfor i, w := range words {
        \t\t\tout[i] = javalang.StringToUpperCaseASCII(w)
        \t\t}
        \t\treturn strings.Join(out, "_")
        \tcase "UPPER_CAMEL":
        \t\tvar sb strings.Builder
        \t\tfor _, w := range words {
        \t\t\tsb.WriteString(capAsciiWord(w))
        \t\t}
        \t\treturn sb.String()
        \tcase "LOWER_CAMEL":
        \t\tvar sb strings.Builder
        \t\tfor i, w := range words {
        \t\t\tif i == 0 {
        \t\t\t\tsb.WriteString(javalang.StringToLowerCaseASCII(w))
        \t\t\t} else {
        \t\t\t\tsb.WriteString(capAsciiWord(w))
        \t\t\t}
        \t\t}
        \t\treturn sb.String()
        \tdefault:
        \t\treturn s
        \t}
        }

        func guavaStandardSystemPropertyValue(propName string) string {
        \tswitch propName {
        \tcase "USER_DIR":
        \t\twd, err := os.Getwd()
        \t\tif err != nil {
        \t\t\treturn "."
        \t\t}
        \t\treturn wd
        \tcase "FILE_SEPARATOR":
        \t\treturn "/"
        \tcase "PATH_SEPARATOR":
        \t\treturn ":"
        \tcase "LINE_SEPARATOR":
        \t\treturn "\\n"
        \tcase "JAVA_IO_TMPDIR":
        \t\treturn os.TempDir()
        \tcase "USER_HOME":
        \t\th, _ := os.UserHomeDir()
        \t\treturn h
        \tdefault:
        \t\treturn NullString
        \t}
        }
        """;
  }

  private static String buildIoRtPart2() {
    return """

        func newPercentEscaper(safeChars string, plusForSpace bool) *guava.PercentEscaper {
        \treturn guava.NewPercentEscaper(safeChars, plusForSpace)
        }

        func escapePercentCore(core any, s string) string {
        \tif pe, ok := core.(*guava.PercentEscaper); ok && pe != nil {
        \t\treturn pe.Escape(s)
        \t}
        \tpanic("UNIMPLEMENTED_EXT: Escaper.escape without override")
        }

        func newEscaperFunction(esc Ext_com_google_common_escape_Escaper) Ext_com_google_common_base_Function {
        \treturn &jGuavaFunc{fn: func(x any) any {
        \t\tif x == nil || javalang.IsNil(x) {
        \t\t\treturn nil
        \t\t}
        \t\treturn BoxString(esc.Escape_1_string(CastString(x)))
        \t}}
        }

        type jHashFunction struct {
        \thf guava.HashFunction
        }

        func (*jHashFunction) Implements_j2gort_Ext_com_google_common_hash_HashFunction() {}
        func (f *jHashFunction) NewHasher_0() Ext_com_google_common_hash_Hasher {
        \treturn &jHasher{h: f.hf.NewHasher()}
        }

        type jHasher struct {
        \th guava.Hasher
        }

        func (*jHasher) Implements_j2gort_Ext_com_google_common_hash_Hasher()        {}
        func (*jHasher) Implements_j2gort_Ext_com_google_common_hash_PrimitiveSink() {}
        func (hr *jHasher) PutInt(v int32) any {
        \thr.h.PutInt(v)
        \treturn hr
        }
        func (hr *jHasher) PutLong(v int64) any {
        \thr.h.PutLong(v)
        \treturn hr
        }
        func (hr *jHasher) PutUnencodedChars(cs any) any {
        \thr.h.PutUnencodedChars(StrVal(CheckNotNull(cs)))
        \treturn hr
        }
        func (hr *jHasher) PutString(cs any, charset Ext_java_nio_charset_Charset) any {
        \thr.h.PutString(StrVal(CheckNotNull(cs)), unwrapCharset(charset))
        \treturn hr
        }
        func (hr *jHasher) Hash() Ext_com_google_common_hash_HashCode {
        \thc := hr.h.Hash()
        \tobj := &Ext_com_google_common_hash_HashCodeBase{Core_: hc}
        \tobj.BindExtSelf_(obj)
        \treturn obj
        }

        func newFarmHashFunction() Ext_com_google_common_hash_HashFunction {
        \treturn &jHashFunction{hf: guava.Hashing_FarmHashFingerprint64()}
        }

        func newMurmur3_128HashFunction() Ext_com_google_common_hash_HashFunction {
        \treturn &jHashFunction{hf: guava.Hashing_Murmur3_128()}
        }

        func hashCodeAsInt(core any) int32 {
        \treturn core.(guava.HashCode).AsInt()
        }

        func hashCodeAsLong(core any) int64 {
        \treturn core.(guava.HashCode).AsLong()
        }

        func newGuavaBase64Encoding() Ext_com_google_common_io_BaseEncoding {
        \tb := &Ext_com_google_common_io_BaseEncodingBase{Core_: guava.BaseEncoding_Base64()}
        \tb.BindExtSelf_(b)
        \treturn b
        }

        func base64Encode(b []int8) string {
        \treturn guava.BaseEncoding_Base64().Encode(int8ToBytes(b))
        }

        func base64Decode(s string) []int8 {
        \traw, err := guava.BaseEncoding_Base64().DecodeChecked(s)
        \tif err != nil {
        \t\tpanic(javalang.NewIllegalArgumentException(err.Error()))
        \t}
        \treturn bytesToInt8(raw)
        }

        func limitExtInputStream(in Ext_java_io_InputStream, limit int64) Ext_java_io_InputStream {
        \tif limit < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException("limit must be non-negative"))
        \t}
        \tlr := io.LimitReader(asInputStream(in), limit)
        \treturn wrapInputStream(javaio.NewReaderInputStream(lr))
        }

        func readAllExtInputStreamBytes(in Ext_java_io_InputStream) []int8 {
        \tb, err := guava.ByteStreams_ToByteArray(asInputStream(in))
        \tioErr(err)
        \treturn bytesToInt8(b)
        }

        type jCharSourceSpec struct {
        \tcontent string
        \tfile    *javaio.File
        \tcs      *javanio.Charset
        }

        func wrapStringCharSource(s string) Ext_com_google_common_io_CharSource {
        \tobj := &Ext_com_google_common_io_CharSourceBase{Core_: &jCharSourceSpec{content: s}}
        \tobj.BindExtSelf_(obj)
        \treturn obj
        }

        func newGuavaFileCharSource(f *Ext_java_io_File, cs Ext_java_nio_charset_Charset) Ext_com_google_common_io_CharSource {
        \tobj := &Ext_com_google_common_io_CharSourceBase{
        \t\tCore_: &jCharSourceSpec{file: unwrapFile(f), cs: unwrapCharset(cs)},
        \t}
        \tobj.BindExtSelf_(obj)
        \treturn obj
        }

        func charSourceRead(core any) string {
        \tspec := core.(*jCharSourceSpec)
        \tif spec.file == nil {
        \t\treturn spec.content
        \t}
        \tb, err := os.ReadFile(spec.file.GetPath())
        \tioErr(err)
        \treturn spec.cs.Decode(b)
        }

        func charSourceOpenStream(core any) Ext_java_io_Reader {
        \tspec := core.(*jCharSourceSpec)
        \tif spec.file == nil {
        \t\tr := &Ext_java_io_StringReader{}
        \t\tr.Core_ = javaio.NewStringReader(spec.content)
        \t\tr.BindExtSelf_(r)
        \t\treturn r
        \t}
        \tfis, err := javaio.NewFileInputStreamFromFile(spec.file)
        \tioErr(err)
        \tr := &Ext_java_io_InputStreamReader{}
        \tr.Core_ = javanio.NewInputStreamReaderWithCharset(fis, spec.cs)
        \tr.BindExtSelf_(r)
        \treturn r
        }

        type jCharSinkSpec struct {
        \tfile   *javaio.File
        \tcs     *javanio.Charset
        \tappend bool
        }

        func newGuavaFileCharSink(f *Ext_java_io_File, cs Ext_java_nio_charset_Charset, modes []*Ext_com_google_common_io_FileWriteMode) Ext_com_google_common_io_CharSink {
        \tapp := false
        \tfor _, m := range modes {
        \t\tif m != nil && m.Name_ == "APPEND" {
        \t\t\tapp = true
        \t\t}
        \t}
        \tobj := &Ext_com_google_common_io_CharSinkBase{
        \t\tCore_: &jCharSinkSpec{file: unwrapFile(f), cs: unwrapCharset(cs), append: app},
        \t}
        \tobj.BindExtSelf_(obj)
        \treturn obj
        }

        func charSinkWrite(core any, content string) {
        \tspec := core.(*jCharSinkSpec)
        \tdata := spec.cs.Encode(content)
        \tflags := os.O_WRONLY | os.O_CREATE
        \tif spec.append {
        \t\tflags |= os.O_APPEND
        \t} else {
        \t\tflags |= os.O_TRUNC
        \t}
        \tf, err := os.OpenFile(spec.file.GetPath(), flags, 0o666)
        \tioErr(err)
        \t_, wErr := f.Write(data)
        \tcErr := f.Close()
        \tioErr(wErr)
        \tioErr(cErr)
        }

        func newGuavaFileBufferedReader(f *Ext_java_io_File, cs Ext_java_nio_charset_Charset) *Ext_java_io_BufferedReader {
        \tfis, err := javaio.NewFileInputStreamFromFile(unwrapFile(f))
        \tioErr(err)
        \tbr := &Ext_java_io_BufferedReader{}
        \tbr.Core_ = javaio.NewBufferedReader(javanio.NewInputStreamReaderWithCharset(fis, unwrapCharset(cs)))
        \tbr.BindExtSelf_(br)
        \treturn br
        }

        func guavaFileToByteArray(f *Ext_java_io_File) []int8 {
        \tb, err := guava.Files_ToByteArray(unwrapFile(f))
        \tioErr(err)
        \treturn bytesToInt8(b)
        }

        func guavaFileWriteBytes(b []int8, f *Ext_java_io_File) {
        \tioErr(guava.Files_WriteBytes(int8ToBytes(b), unwrapFile(f)))
        }

        func readAllFromReadable(r Ext_java_lang_Readable) string {
        \tif rd, ok := r.(Ext_java_io_Reader); ok {
        \t\treturn guava.CharStreams_ToString(asReader(rd))
        \t}
        \tpanic(fmt.Sprintf("UNIMPLEMENTED_EXT: CharStreams.toString on %T", r))
        }

        func readLinesFromReadable(r Ext_java_lang_Readable) Ext_java_util_List {
        \trd, ok := r.(Ext_java_io_Reader)
        \tif !ok {
        \t\tpanic(fmt.Sprintf("UNIMPLEMENTED_EXT: CharStreams.readLines on %T", r))
        \t}
        \tlines := guava.CharStreams_ReadLines(asReader(rd))
        \titems := make([]any, len(lines))
        \tfor i, ln := range lines {
        \t\titems[i] = BoxString(ln)
        \t}
        \treturn wrapArrayList(items)
        }

        // --- Gson & JsonElement Helpers + GsonWritable Implementations ---

        type jGsonAdapterBridge struct {
        \text Ext_com_google_gson_TypeAdapter
        }

        func (b *jGsonAdapterBridge) Write(out *gson.JsonWriter, val any) {
        \tjw := &Ext_com_google_gson_stream_JsonWriter{Core_: out}
        \tjw.BindExtSelf_(jw)
        \tb.ext.Write(jw, val)
        }

        type jGsonBuilderCore struct {
        \tb *gson.GsonBuilder
        }

        func newGsonBuilderCore() *jGsonBuilderCore {
        \treturn &jGsonBuilderCore{b: gson.NewGsonBuilder()}
        }

        func newDefaultGson() *gson.Gson {
        \treturn gson.NewGson()
        }

        func gsonBuilderSetPretty(gb *Ext_com_google_gson_GsonBuilder) {
        \tgb.Core_.(*jGsonBuilderCore).b.SetPrettyPrinting()
        }

        func gsonBuilderDisableHtml(gb *Ext_com_google_gson_GsonBuilder) {
        \tgb.Core_.(*jGsonBuilderCore).b.DisableHtmlEscaping()
        }

        func gsonBuilderRegisterHierarchyAdapter(gb *Ext_com_google_gson_GsonBuilder, baseClass *javalang.Class, adapter any) {
        \tCheckNotNullPtr(baseClass)
        \tta := AsIface[Ext_com_google_gson_TypeAdapter](CheckNotNull(adapter))
        \tgb.Core_.(*jGsonBuilderCore).b.RegisterTypeHierarchyAdapter(
        \t\tfunc(v any) bool { return CheckTypePredicate(baseClass, v) },
        \t\t&jGsonAdapterBridge{ext: ta},
        \t)
        }

        func gsonBuilderCreate(gb *Ext_com_google_gson_GsonBuilder) *Ext_com_google_gson_Gson {
        \tg := &Ext_com_google_gson_Gson{Core_: gb.Core_.(*jGsonBuilderCore).b.Create()}
        \tg.BindExtSelf_(g)
        \treturn g
        }

        func unwrapGson(g *Ext_com_google_gson_Gson) *gson.Gson {
        \tCheckNotNullPtr(g)
        \tif gs, ok := g.Core_.(*gson.Gson); ok && gs != nil {
        \t\treturn gs
        \t}
        \treturn gson.NewGson()
        }

        func gsonToJson(g *Ext_com_google_gson_Gson, v any) string {
        \tif v == nil || javalang.IsNil(v) {
        \t\treturn "null"
        \t}
        \tif je, ok := v.(Ext_com_google_gson_JsonElement); ok {
        \t\treturn unwrapGson(g).ToJson(unwrapJsonElement(je))
        \t}
        \treturn unwrapGson(g).ToJson(ToAny(v))
        }

        func gsonToJsonWriter(g *Ext_com_google_gson_Gson, v any, w *Ext_com_google_gson_stream_JsonWriter) {
        \tjw := gsonWriter(CheckNotNullPtr(w))
        \tif v == nil || javalang.IsNil(v) {
        \t\tjw.NullValue()
        \t\treturn
        \t}
        \tif je, ok := v.(Ext_com_google_gson_JsonElement); ok {
        \t\tgson.Streams_Write(unwrapJsonElement(je), jw)
        \t\treturn
        \t}
        \tunwrapGson(g).ToJsonWriter(ToAny(v), jw)
        }

        func wrapGsonPanic(fn func() any) (res any) {
        \tdefer func() {
        \t\tif r := recover(); r != nil {
        \t\t\tif gp, ok := r.(gson.IJsonParseException); ok {
        \t\t\t\tex := &Ext_com_google_gson_JsonParseExceptionBase{}
        \t\t\t\tex.BindExtSelf_(ex)
        \t\t\t\tInitUserThrowable(ex, "com.google.gson.JsonParseException", gp.GetMessage(), gp)
        \t\t\t\tpanic(ex)
        \t\t\t}
        \t\t\tpanic(r)
        \t\t}
        \t}()
        \treturn fn()
        }

        func gsonFromJsonClass(g *Ext_com_google_gson_Gson, jsonStr string, cls *javalang.Class) any {
        \tif jsonStr == NullString {
        \t\treturn nil
        \t}
        \tCheckNotNullPtr(cls)
        \tgs := unwrapGson(g)
        \treturn wrapGsonPanic(func() any {
        \t\tswitch cls.GetName() {
        \t\tcase "com.google.gson.JsonObject":
        \t\t\tjo := gs.FromJsonObject(jsonStr)
        \t\t\tif jo == nil {
        \t\t\t\treturn nil
        \t\t\t}
        \t\t\treturn wrapJsonObject(jo)
        \t\tcase "com.google.gson.JsonArray":
        \t\t\tja := gs.FromJsonArray(jsonStr)
        \t\t\tif ja == nil {
        \t\t\t\treturn nil
        \t\t\t}
        \t\t\treturn wrapJsonArray(ja)
        \t\tcase "com.google.gson.JsonElement":
        \t\t\tje := gs.FromJsonElement(jsonStr)
        \t\t\tif je == nil {
        \t\t\t\treturn nil
        \t\t\t}
        \t\t\treturn wrapJsonElement(je)
        \t\tdefault:
        \t\t\tpanic("UNIMPLEMENTED_EXT: Gson.fromJson(String, " + cls.GetName() + ")")
        \t\t}
        \t})
        }

        func gsonFromJsonType(g *Ext_com_google_gson_Gson, jsonStr string, typ any) any {
        \tif jsonStr == NullString {
        \t\treturn nil
        \t}
        \tif cls, ok := typ.(*javalang.Class); ok {
        \t\treturn gsonFromJsonClass(g, jsonStr, cls)
        \t}
        \tts := fmt.Sprint(typ)
        \tgs := unwrapGson(g)
        \treturn wrapGsonPanic(func() any {
        \t\tif strings.Contains(ts, "Map<java.lang.String, java.lang.Long>") || strings.Contains(ts, "Map<String, Long>") {
        \t\t\tm := gs.FromJsonStringLongMap(jsonStr)
        \t\t\tif m == nil {
        \t\t\t\treturn nil
        \t\t\t}
        \t\t\tmc := NewJMapCore(false)
        \t\t\tfor k, v := range m.All() {
        \t\t\t\tmc.Put(BoxString(k), ValueOfLong(v))
        \t\t\t}
        \t\t\treturn wrapLinkedHashMap(mc)
        \t\t}
        \t\tif strings.Contains(ts, "List<java.lang.String>") || strings.Contains(ts, "List<String>") {
        \t\t\tl := gs.FromJsonStringList(jsonStr)
        \t\t\tif l == nil {
        \t\t\t\treturn nil
        \t\t\t}
        \t\t\ts := l.ToSlice()
        \t\t\titems := make([]any, len(s))
        \t\t\tfor i, str := range s {
        \t\t\t\titems[i] = BoxString(str)
        \t\t\t}
        \t\t\treturn wrapArrayList(items)
        \t\t}
        \t\tpanic("UNIMPLEMENTED_EXT: Gson.fromJson(String, Type=" + ts + ")")
        \t})
        }

        func newGsonJsonArray() *gson.JsonArray   { return gson.NewJsonArray() }
        func newGsonJsonObject() *gson.JsonObject { return gson.NewJsonObject() }
        func newGsonJsonPrimitive(v any) *gson.JsonPrimitive {
        \tif s, ok := v.(string); ok {
        \t\treturn gson.NewJsonPrimitiveString(CheckNotNullString(s))
        \t}
        \tif b, ok := v.(bool); ok {
        \t\treturn gson.NewJsonPrimitiveBoolean(b)
        \t}
        \tif bb, ok := v.(*BoxBoolean); ok {
        \t\treturn gson.NewJsonPrimitiveBoolean(UnboxBoolean(bb))
        \t}
        \tif n, ok := v.(Number); ok {
        \t\treturn gson.NewJsonPrimitiveLong(n.LongValue())
        \t}
        \treturn gson.NewJsonPrimitiveString(StrVal(v))
        }

        func unwrapJsonElement(je Ext_com_google_gson_JsonElement) gson.JsonElement {
        \tif je == nil || javalang.IsNil(je) {
        \t\treturn gson.JsonNull_INSTANCE
        \t}
        \tif e, ok := je.GetExt_com_google_gson_JsonElementBase().Core_.(gson.JsonElement); ok && e != nil {
        \t\treturn e
        \t}
        \treturn gson.JsonNull_INSTANCE
        }

        func unwrapJsonArray(ja *Ext_com_google_gson_JsonArray) *gson.JsonArray {
        \tCheckNotNullPtr(ja)
        \treturn ja.Core_.(*gson.JsonArray)
        }

        func unwrapJsonObject(jo *Ext_com_google_gson_JsonObject) *gson.JsonObject {
        \tCheckNotNullPtr(jo)
        \treturn jo.Core_.(*gson.JsonObject)
        }

        func unwrapJsonPrimitive(jp *Ext_com_google_gson_JsonPrimitive) *gson.JsonPrimitive {
        \tCheckNotNullPtr(jp)
        \treturn jp.Core_.(*gson.JsonPrimitive)
        }

        func wrapJsonObject(jo *gson.JsonObject) *Ext_com_google_gson_JsonObject {
        \tif jo == nil {
        \t\treturn nil
        \t}
        \to := &Ext_com_google_gson_JsonObject{}
        \to.Core_ = jo
        \to.BindExtSelf_(o)
        \treturn o
        }

        func wrapJsonArray(ja *gson.JsonArray) *Ext_com_google_gson_JsonArray {
        \tif ja == nil {
        \t\treturn nil
        \t}
        \ta := &Ext_com_google_gson_JsonArray{}
        \ta.Core_ = ja
        \ta.BindExtSelf_(a)
        \treturn a
        }

        func wrapJsonPrimitive(jp *gson.JsonPrimitive) *Ext_com_google_gson_JsonPrimitive {
        \tif jp == nil {
        \t\treturn nil
        \t}
        \tp := &Ext_com_google_gson_JsonPrimitive{}
        \tp.Core_ = jp
        \tp.BindExtSelf_(p)
        \treturn p
        }

        func wrapJsonElement(je gson.JsonElement) Ext_com_google_gson_JsonElement {
        \tif je == nil {
        \t\treturn nil
        \t}
        \tswitch x := je.(type) {
        \tcase *gson.JsonObject:
        \t\treturn wrapJsonObject(x)
        \tcase *gson.JsonArray:
        \t\treturn wrapJsonArray(x)
        \tcase *gson.JsonPrimitive:
        \t\treturn wrapJsonPrimitive(x)
        \tcase *gson.JsonNull:
        \t\tn := &Ext_com_google_gson_JsonNull{}
        \t\tn.Core_ = x
        \t\tn.BindExtSelf_(n)
        \t\treturn n
        \tdefault:
        \t\tb := &Ext_com_google_gson_JsonElementBase{Core_: je}
        \t\tb.BindExtSelf_(b)
        \t\treturn b
        \t}
        }

        func gsonJsonArrayAdd(ja *Ext_com_google_gson_JsonArray, elem Ext_com_google_gson_JsonElement) {
        \tunwrapJsonArray(ja).Add(unwrapJsonElement(elem))
        }

        func gsonJsonArrayIterator(ja *Ext_com_google_gson_JsonArray) Ext_java_util_Iterator {
        \ta := unwrapJsonArray(ja)
        \tn := a.Size()
        \titems := make([]any, n)
        \tfor i := 0; i < n; i++ {
        \t\titems[i] = wrapJsonElement(a.Get(i))
        \t}
        \treturn newSliceObjIter(items)
        }

        func gsonJsonArrayForEach(ja *Ext_com_google_gson_JsonArray, action Ext_java_util_function_Consumer) {
        \ta := unwrapJsonArray(ja)
        \tn := a.Size()
        \tfor i := 0; i < n; i++ {
        \t\taction.Accept(wrapJsonElement(a.Get(i)))
        \t}
        }

        func gsonJsonObjectAdd(jo *Ext_com_google_gson_JsonObject, key string, val Ext_com_google_gson_JsonElement) {
        \tunwrapJsonObject(jo).Add(key, unwrapJsonElement(val))
        }

        func gsonJsonObjectEntrySet(jo *Ext_com_google_gson_JsonObject) Ext_java_util_Set {
        \to := unwrapJsonObject(jo)
        \tmc := NewJMapCore(false)
        \tfor k, v := range o.EntrySet().All() {
        \t\tmc.Put(BoxString(k), wrapJsonElement(v))
        \t}
        \treturn newLiveMapEntrySet(mc)
        }

        func newGsonJsonReader(r Ext_java_io_Reader) *gson.JsonReader {
        \treturn gson.NewJsonReader(asReader(r))
        }

        func newGsonJsonWriter(w Ext_java_io_Writer) *gson.JsonWriter {
        \treturn gson.NewJsonWriter(asWriter(w))
        }

        func gsonReader(r *Ext_com_google_gson_stream_JsonReader) *gson.JsonReader {
        \tCheckNotNullPtr(r)
        \treturn r.Core_.(*gson.JsonReader)
        }

        func gsonWriter(w *Ext_com_google_gson_stream_JsonWriter) *gson.JsonWriter {
        \tCheckNotNullPtr(w)
        \treturn w.Core_.(*gson.JsonWriter)
        }

        func writeGsonAny(out *gson.JsonWriter, v any) {
        \tv = ToAny(v)
        \tif v == nil {
        \t\tout.NullValue()
        \t\treturn
        \t}
        \tswitch x := v.(type) {
        \tcase string:
        \t\tif x == NullString {
        \t\t\tout.NullValue()
        \t\t} else {
        \t\t\tout.ValueString(x)
        \t\t}
        \tcase bool:
        \t\tout.ValueBoolean(x)
        \tcase *BoxBoolean:
        \t\tout.ValueBoolean(x.Val)
        \tcase *BoxInt:
        \t\tout.ValueInt(x.Val)
        \tcase *BoxLong:
        \t\tout.ValueLong(x.Val)
        \tcase *BoxByte:
        \t\tout.ValueInt(int32(x.Val))
        \tcase *BoxShort:
        \t\tout.ValueInt(int32(x.Val))
        \tcase *BoxFloat:
        \t\tout.ValueDouble(float64(x.Val))
        \tcase *BoxDouble:
        \t\tout.ValueDouble(x.Val)
        \tcase *BoxChar:
        \t\tout.ValueString(javalang.StringFromChar(x.Val))
        \tcase EnumAccessor:
        \t\tout.ValueString(x.GetEnumBase().Name_)
        \tcase Ext_com_google_gson_JsonElement:
        \t\tgson.Streams_Write(unwrapJsonElement(x), out)
        \tcase gson.GsonWritable:
        \t\tx.WriteGson(out)
        \tcase Ext_java_util_Map:
        \t\tout.BeginObject()
        \t\tfor it := Iterate(x.EntrySet()); it.HasNext(); {
        \t\t\te := AsIface[Ext_java_util_Map_Entry](it.Next())
        \t\t\tout.Name(StrVal(e.GetKey()))
        \t\t\twriteGsonAny(out, e.GetValue())
        \t\t}
        \t\tout.EndObject()
        \tcase Ext_java_lang_Iterable:
        \t\tout.BeginArray()
        \t\tfor it := Iterate(x); it.HasNext(); {
        \t\t\twriteGsonAny(out, it.Next())
        \t\t}
        \t\tout.EndArray()
        \tdefault:
        \t\t// Struct reflection via default Gson
        \t\tout.JsonValue(gson.NewGson().ToJson(v))
        \t}
        }

        func (b *BoxBoolean) WriteGson(out *gson.JsonWriter) { out.ValueBoolean(b.Val) }
        func (b *BoxByte) WriteGson(out *gson.JsonWriter)    { out.ValueInt(int32(b.Val)) }
        func (b *BoxShort) WriteGson(out *gson.JsonWriter)   { out.ValueInt(int32(b.Val)) }
        func (b *BoxChar) WriteGson(out *gson.JsonWriter)    { out.ValueString(javalang.StringFromChar(b.Val)) }
        func (b *BoxInt) WriteGson(out *gson.JsonWriter)     { out.ValueInt(b.Val) }
        func (b *BoxLong) WriteGson(out *gson.JsonWriter)    { out.ValueLong(b.Val) }
        func (b *BoxFloat) WriteGson(out *gson.JsonWriter)   { out.ValueDouble(float64(b.Val)) }
        func (b *BoxDouble) WriteGson(out *gson.JsonWriter)  { out.ValueDouble(b.Val) }
        func (e *EnumBase) WriteGson(out *gson.JsonWriter)   { out.ValueString(e.Name_) }

        func (l *JListCore) WriteGson(out *gson.JsonWriter) {
        \tout.BeginArray()
        \tfor _, item := range l.Items {
        \t\twriteGsonAny(out, item)
        \t}
        \tout.EndArray()
        }

        func (m *JMapCore) WriteGson(out *gson.JsonWriter) {
        \tout.BeginObject()
        \tm.ForEachEntry(func(k, v any) bool {
        \t\tout.Name(StrVal(k))
        \t\twriteGsonAny(out, v)
        \t\treturn true
        \t})
        \tout.EndObject()
        }

        func (b *Ext_java_util_AbstractCollectionBase) WriteGson(out *gson.JsonWriter) {
        \tvar target Ext_java_lang_Iterable = b
        \tif b.Self_ != nil {
        \t\ttarget = b.Self_
        \t}
        \tout.BeginArray()
        \tfor it := Iterate(target); it.HasNext(); {
        \t\twriteGsonAny(out, it.Next())
        \t}
        \tout.EndArray()
        }

        func (b *Ext_java_util_AbstractMapBase) WriteGson(out *gson.JsonWriter) {
        \tvar target Ext_java_util_Map = b
        \tif b.Self_ != nil {
        \t\ttarget = b.Self_
        \t}
        \tout.BeginObject()
        \tfor it := Iterate(target.EntrySet()); it.HasNext(); {
        \t\te := AsIface[Ext_java_util_Map_Entry](it.Next())
        \t\tout.Name(StrVal(e.GetKey()))
        \t\twriteGsonAny(out, e.GetValue())
        \t}
        \tout.EndObject()
        }

        func (b *Ext_com_google_common_collect_ImmutableCollectionBase) WriteGson(out *gson.JsonWriter) {
        \tvar target Ext_java_lang_Iterable = b
        \tif b.Self_ != nil {
        \t\ttarget = b.Self_
        \t}
        \tout.BeginArray()
        \tfor it := Iterate(target); it.HasNext(); {
        \t\twriteGsonAny(out, it.Next())
        \t}
        \tout.EndArray()
        }

        func (b *Ext_com_google_common_collect_ImmutableMapBase) WriteGson(out *gson.JsonWriter) {
        \tvar target Ext_java_util_Map = b
        \tif b.Self_ != nil {
        \t\ttarget = b.Self_
        \t}
        \tout.BeginObject()
        \tfor it := Iterate(target.EntrySet()); it.HasNext(); {
        \t\te := AsIface[Ext_java_util_Map_Entry](it.Next())
        \t\tout.Name(StrVal(e.GetKey()))
        \t\twriteGsonAny(out, e.GetValue())
        \t}
        \tout.EndObject()
        }

        // --- Clean-Room SAX XML Reader Adapter (over Go encoding/xml) ---

        type jSaxAttrs struct {
        \tattrs []xml.Attr
        }

        func (*jSaxAttrs) Implements_j2gort_Ext_org_xml_sax_Attributes() {}

        func (a *jSaxAttrs) GetValue_1_string(qName string) string {
        \tfor _, attr := range a.attrs {
        \t\tfull := attr.Name.Local
        \t\tif attr.Name.Space != "" && !strings.HasPrefix(attr.Name.Space, "http") {
        \t\t\tfull = attr.Name.Space + ":" + attr.Name.Local
        \t\t}
        \t\tif full == qName || attr.Name.Local == qName {
        \t\t\treturn attr.Value
        \t\t}
        \t}
        \treturn NullString
        }

        type jSaxXMLReader struct {
        \tentityResolver Ext_org_xml_sax_EntityResolver
        \tcontentHandler Ext_org_xml_sax_ContentHandler
        }

        func (*jSaxXMLReader) Implements_j2gort_Ext_org_xml_sax_XMLReader() {}

        func (xr *jSaxXMLReader) SetEntityResolver(er Ext_org_xml_sax_EntityResolver) {
        \txr.entityResolver = er
        }

        func (xr *jSaxXMLReader) SetContentHandler(ch Ext_org_xml_sax_ContentHandler) {
        \txr.contentHandler = ch
        }

        func throwSAXException(msg string) {
        \tex := &Ext_org_xml_sax_SAXExceptionBase{}
        \tex.BindExtSelf_(ex)
        \tInitUserThrowable(ex, "org.xml.sax.SAXException", msg)
        \tpanic(ex)
        }

        func (xr *jSaxXMLReader) Parse_1_Ext_org_xml_sax_InputSource(src *Ext_org_xml_sax_InputSource) {
        \tCheckNotNullPtr(src)
        \tvar raw []byte
        \tswitch c := src.Core_.(type) {
        \tcase Ext_java_io_InputStream:
        \t\tb, err := io.ReadAll(asInputStream(c))
        \t\tioErr(err)
        \t\traw = b
        \tcase Ext_java_io_Reader:
        \t\ts := guava.CharStreams_ToString(asReader(c))
        \t\traw = []byte(s)
        \tdefault:
        \t\tthrowSAXException("Unsupported InputSource")
        \t}

        \tdec := xml.NewDecoder(bytes.NewReader(raw))
        \tdec.Strict = false
        \t// Note: custom entities in XTB files if any
        \tdec.Entity = xml.HTMLEntity

        \tif xr.contentHandler != nil {
        \t\txr.contentHandler.StartDocument()
        \t}
        \tfor {
        \t\ttok, err := dec.Token()
        \t\tif err == io.EOF {
        \t\t\tbreak
        \t\t}
        \t\tif err != nil {
        \t\t\tthrowSAXException(err.Error())
        \t\t}
        \t\tif xr.contentHandler == nil {
        \t\t\tcontinue
        \t\t}
        \t\tswitch t := tok.(type) {
        \t\tcase xml.StartElement:
        \t\t\txr.contentHandler.StartElement(t.Name.Space, t.Name.Local, t.Name.Local, &jSaxAttrs{attrs: t.Attr})
        \t\tcase xml.EndElement:
        \t\t\txr.contentHandler.EndElement(t.Name.Space, t.Name.Local, t.Name.Local)
        \t\tcase xml.CharData:
        \t\t\tu16 := utf16.Encode([]rune(string(t)))
        \t\t\tif len(u16) > 0 {
        \t\t\t\txr.contentHandler.Characters(u16, 0, int32(len(u16)))
        \t\t\t}
        \t\tcase xml.ProcInst:
        \t\t\txr.contentHandler.ProcessingInstruction(t.Target, string(t.Inst))
        \t\t}
        \t}
        \tif xr.contentHandler != nil {
        \t\txr.contentHandler.EndDocument()
        \t}
        }

        // --- Regex, MessageFormat, ResourceBundle, StringTokenizer, Logging, JMX, System ---

        func wrapPatternCompile(regex string, flags int) *Ext_java_util_regex_Pattern {
        \tp := &Ext_java_util_regex_Pattern{Core_: javaregex.Pattern_CompileFlags(regex, flags)}
        \tp.BindExtSelf_(p)
        \treturn p
        }

        func javaMessageFormat(pattern string, args []any) string {
        \tvar sb strings.Builder
        \tinQuote := false
        \tn := len(pattern)
        \tfor i := 0; i < n; i++ {
        \t\tc := pattern[i]
        \t\tif c == '\\'' {
        \t\t\tif i+1 < n && pattern[i+1] == '\\'' {
        \t\t\t\tsb.WriteByte('\\'')
        \t\t\t\ti++
        \t\t\t\tcontinue
        \t\t\t}
        \t\t\tinQuote = !inQuote
        \t\t\tcontinue
        \t\t}
        \t\tif !inQuote && c == '{' {
        \t\t\tend := strings.IndexByte(pattern[i+1:], '}')
        \t\t\tif end > 0 {
        \t\t\t\tspec := pattern[i+1 : i+1+end]
        \t\t\t\tif comma := strings.IndexByte(spec, ','); comma >= 0 {
        \t\t\t\t\tspec = spec[:comma]
        \t\t\t\t}
        \t\t\t\tif idx, err := strconv.Atoi(strings.TrimSpace(spec)); err == nil && idx >= 0 {
        \t\t\t\t\tif idx < len(args) {
        \t\t\t\t\t\tsb.WriteString(StrVal(args[idx]))
        \t\t\t\t\t} else {
        \t\t\t\t\t\tsb.WriteByte('{')
        \t\t\t\t\t\tsb.WriteString(pattern[i+1 : i+1+end])
        \t\t\t\t\t\tsb.WriteByte('}')
        \t\t\t\t\t}
        \t\t\t\t\ti += end + 1
        \t\t\t\t\tcontinue
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t\tsb.WriteByte(c)
        \t}
        \treturn sb.String()
        }

        func wrapResourceBundle(baseName string) Ext_java_util_ResourceBundle {
        \trb := shimsres.ResourceBundle_GetBundle(baseName)
        \tobj := &Ext_java_util_ResourceBundleBase{Core_: rb}
        \tobj.BindExtSelf_(obj)
        \treturn obj
        }

        type jStringTokenizer struct {
        \tstr    string
        \tdelims string
        \tpos    int
        }

        func newStringTokenizer(str, delims string) *jStringTokenizer {
        \treturn &jStringTokenizer{str: str, delims: delims, pos: 0}
        }

        func (st *jStringTokenizer) isDelim(b byte) bool {
        \treturn strings.IndexByte(st.delims, b) >= 0
        }

        func (st *jStringTokenizer) skipDelims(p int) int {
        \tfor p < len(st.str) && st.isDelim(st.str[p]) {
        \t\tp++
        \t}
        \treturn p
        }

        func (st *jStringTokenizer) hasMoreTokens() bool {
        \tst.pos = st.skipDelims(st.pos)
        \treturn st.pos < len(st.str)
        }

        func (st *jStringTokenizer) nextToken() string {
        \tst.pos = st.skipDelims(st.pos)
        \tif st.pos >= len(st.str) {
        \t\tpanic(javalang.NewNoSuchElementException("StringTokenizer"))
        \t}
        \tstart := st.pos
        \tfor st.pos < len(st.str) && !st.isDelim(st.str[st.pos]) {
        \t\tst.pos++
        \t}
        \treturn st.str[start:st.pos]
        }

        var loggerCache sync.Map // string -> *Ext_java_util_logging_Logger

        func getJavaLogger(name string) *Ext_java_util_logging_Logger {
        \tif v, ok := loggerCache.Load(name); ok {
        \t\treturn v.(*Ext_java_util_logging_Logger)
        \t}
        \tl := &Ext_java_util_logging_Logger{Core_: Ext_java_util_logging_Level_F_OFF}
        \tl.BindExtSelf_(l)
        \tactual, _ := loggerCache.LoadOrStore(name, l)
        \treturn actual.(*Ext_java_util_logging_Logger)
        }

        func loggerIsLoggable(l *Ext_java_util_logging_Logger, lvl *Ext_java_util_logging_Level) bool {
        \tif l == nil || lvl == nil {
        \t\treturn false
        \t}
        \tcur, _ := l.Core_.(*Ext_java_util_logging_Level)
        \tif cur == nil || cur.GetName() == "OFF" || lvl.GetName() == "OFF" {
        \t\treturn false
        \t}
        \treturn false
        }

        var (
        \tmainJavaThreadOnce sync.Once
        \tmainJavaThread     *Ext_java_lang_Thread
        \tgoroutineThreads   sync.Map
        )

        func currentJavaThread() *Ext_java_lang_Thread {
        \tgid := currentGoroutineID()
        \tif gid == 1 {
        \t\tmainJavaThreadOnce.Do(func() {
        \t\t\tt := &Ext_java_lang_Thread{}
        \t\t\tt.BindExtSelf_(t)
        \t\t\tmainJavaThread = t
        \t\t})
        \t\treturn mainJavaThread
        \t}
        \tif v, ok := goroutineThreads.Load(gid); ok {
        \t\treturn v.(*Ext_java_lang_Thread)
        \t}
        \tt := &Ext_java_lang_Thread{}
        \tt.BindExtSelf_(t)
        \tactual, _ := goroutineThreads.LoadOrStore(gid, t)
        \treturn actual.(*Ext_java_lang_Thread)
        }

        type jThreadLocal struct {
        \tmu   sync.Mutex
        \tvals map[uint64]any
        }

        func (tl *jThreadLocal) get() any {
        \tgid := currentGoroutineID()
        \ttl.mu.Lock()
        \tdefer tl.mu.Unlock()
        \tif tl.vals == nil {
        \t\treturn nil
        \t}
        \treturn tl.vals[gid]
        }

        func (tl *jThreadLocal) set(v any) {
        \tgid := currentGoroutineID()
        \ttl.mu.Lock()
        \tdefer tl.mu.Unlock()
        \tif tl.vals == nil {
        \t\ttl.vals = make(map[uint64]any)
        \t}
        \ttl.vals[gid] = ToAny(v)
        }

        func (tl *jThreadLocal) remove() {
        \tgid := currentGoroutineID()
        \ttl.mu.Lock()
        \tdefer tl.mu.Unlock()
        \tdelete(tl.vals, gid)
        }

        type jMemoryUsage struct {
        \tinit, used, committed, max int64
        }

        func runtimeFreeMemory() int64 {
        \tvar ms runtime.MemStats
        \truntime.ReadMemStats(&ms)
        \tif ms.HeapIdle > ms.HeapReleased {
        \t\treturn int64(ms.HeapIdle - ms.HeapReleased)
        \t}
        \treturn 0
        }

        func runtimeTotalMemory() int64 {
        \tvar ms runtime.MemStats
        \truntime.ReadMemStats(&ms)
        \treturn int64(ms.HeapSys)
        }

        func runtimeMaxMemory() int64 {
        \treturn 4 * 1024 * 1024 * 1024
        }

        type jCompilationMXBean struct{}

        func (*jCompilationMXBean) Implements_j2gort_Ext_java_lang_management_CompilationMXBean()     {}
        func (*jCompilationMXBean) Implements_j2gort_Ext_java_lang_management_PlatformManagedObject() {}
        func (*jCompilationMXBean) GetName() string                                                   { return "GoCompiler" }
        func (*jCompilationMXBean) GetTotalCompilationTime() int64                                    { return 0 }

        type jGarbageCollectorMXBean struct{}

        func (*jGarbageCollectorMXBean) Implements_j2gort_Ext_java_lang_management_GarbageCollectorMXBean() {}
        func (*jGarbageCollectorMXBean) Implements_j2gort_Ext_java_lang_management_MemoryManagerMXBean()    {}
        func (*jGarbageCollectorMXBean) Implements_j2gort_Ext_java_lang_management_PlatformManagedObject()  {}
        func (*jGarbageCollectorMXBean) GetName() string                                                    { return "GoGC" }
        func (*jGarbageCollectorMXBean) GetCollectionCount() int64 {
        \tvar ms runtime.MemStats
        \truntime.ReadMemStats(&ms)
        \treturn int64(ms.NumGC)
        }
        func (*jGarbageCollectorMXBean) GetCollectionTime() int64 {
        \tvar ms runtime.MemStats
        \truntime.ReadMemStats(&ms)
        \treturn int64(ms.PauseTotalNs / 1_000_000)
        }

        type jMemoryMXBean struct{}

        func (*jMemoryMXBean) Implements_j2gort_Ext_java_lang_management_MemoryMXBean()          {}
        func (*jMemoryMXBean) Implements_j2gort_Ext_java_lang_management_PlatformManagedObject() {}
        func (*jMemoryMXBean) GetHeapMemoryUsage() *Ext_java_lang_management_MemoryUsage {
        \tvar ms runtime.MemStats
        \truntime.ReadMemStats(&ms)
        \treturn &Ext_java_lang_management_MemoryUsage{
        \t\tCore_: &jMemoryUsage{init: 0, used: int64(ms.HeapAlloc), committed: int64(ms.HeapSys), max: runtimeMaxMemory()},
        \t}
        }
        func (*jMemoryMXBean) GetNonHeapMemoryUsage() *Ext_java_lang_management_MemoryUsage {
        \tvar ms runtime.MemStats
        \truntime.ReadMemStats(&ms)
        \treturn &Ext_java_lang_management_MemoryUsage{
        \t\tCore_: &jMemoryUsage{init: 0, used: int64(ms.StackInuse), committed: int64(ms.StackSys), max: -1},
        \t}
        }

        type jMemoryPoolMXBean struct{}

        func (*jMemoryPoolMXBean) Implements_j2gort_Ext_java_lang_management_MemoryPoolMXBean()      {}
        func (*jMemoryPoolMXBean) Implements_j2gort_Ext_java_lang_management_PlatformManagedObject() {}
        func (*jMemoryPoolMXBean) GetName() string                                                   { return "GoHeap" }
        func (*jMemoryPoolMXBean) GetUsage() *Ext_java_lang_management_MemoryUsage {
        \treturn (&jMemoryMXBean{}).GetHeapMemoryUsage()
        }
        func (*jMemoryPoolMXBean) GetPeakUsage() *Ext_java_lang_management_MemoryUsage {
        \treturn (&jMemoryMXBean{}).GetHeapMemoryUsage()
        }

        func javaSystemGetProperty(key, def string) string {
        \tswitch key {
        \tcase "line.separator":
        \t\treturn "\\n"
        \tcase "file.separator":
        \t\treturn "/"
        \tcase "path.separator":
        \t\treturn ":"
        \tcase "user.dir":
        \t\twd, err := os.Getwd()
        \t\tif err != nil {
        \t\t\treturn "."
        \t\t}
        \t\treturn wd
        \tcase "user.home":
        \t\th, err := os.UserHomeDir()
        \t\tif err == nil {
        \t\t\treturn h
        \t\t}
        \tcase "java.io.tmpdir":
        \t\treturn os.TempDir()
        \t}
        \tif val, ok := os.LookupEnv("JSGO_PROP_" + strings.ReplaceAll(key, ".", "_")); ok {
        \t\treturn val
        \t}
        \treturn def
        }

        func javaSystemGetenv(key string) string {
        \tval, ok := os.LookupEnv(key)
        \tif !ok {
        \t\treturn NullString
        \t}
        \treturn val
        }

        func ThrowablePrintStackTrace(t javalang.Throwable, dest any) {
        \tCheckNotNull(t)
        \tmsg := t.Error()
        \tif dest == nil || javalang.IsNil(dest) {
        \t\tjavaio.System_Err.Println(msg)
        \t\treturn
        \t}
        \tswitch d := dest.(type) {
        \tcase *Ext_java_io_PrintStream:
        \t\tunwrapPrintStream(d).Println(msg)
        \tcase *Ext_java_io_PrintWriter:
        \t\textWriterWriteString(d.Core_, msg+"\\n")
        \t\textWriterFlush(d.Core_)
        \tdefault:
        \t\tjavaio.System_Err.Println(msg)
        \t}
        }

        func (x *extInputStreamAdapter) IsJavaNil_() bool   { return x == nil }
        func (x *extOutputStreamAdapter) IsJavaNil_() bool  { return x == nil }
        func (x *extWriterAdapter) IsJavaNil_() bool        { return x == nil }
        func (x *gzipDeflaterBridge) IsJavaNil_() bool      { return x == nil }
        func (x *jNioPath) IsJavaNil_() bool                { return x == nil }
        func (x *jPathMatcher) IsJavaNil_() bool            { return x == nil }
        func (x *jBasicFileAttrs) IsJavaNil_() bool         { return x == nil }
        func (x *jFileVisitorAdapter) IsJavaNil_() bool     { return x == nil }
        func (x *jGuavaFunc) IsJavaNil_() bool              { return x == nil }
        func (x *jGuavaPred) IsJavaNil_() bool              { return x == nil }
        func (x *jMemoSupplier) IsJavaNil_() bool           { return x == nil }
        func (x *jToStringHelperData) IsJavaNil_() bool     { return x == nil }
        func (x *jHashFunction) IsJavaNil_() bool           { return x == nil }
        func (x *jHasher) IsJavaNil_() bool                 { return x == nil }
        func (x *jCharSourceSpec) IsJavaNil_() bool         { return x == nil }
        func (x *jCharSinkSpec) IsJavaNil_() bool           { return x == nil }
        func (x *jGsonAdapterBridge) IsJavaNil_() bool      { return x == nil }
        func (x *jGsonBuilderCore) IsJavaNil_() bool        { return x == nil }
        func (x *jSaxAttrs) IsJavaNil_() bool               { return x == nil }
        func (x *jSaxXMLReader) IsJavaNil_() bool           { return x == nil }
        func (x *jStringTokenizer) IsJavaNil_() bool        { return x == nil }
        func (x *jThreadLocal) IsJavaNil_() bool            { return x == nil }
        func (x *jMemoryUsage) IsJavaNil_() bool            { return x == nil }
        func (x *jCompilationMXBean) IsJavaNil_() bool      { return x == nil }
        func (x *jGarbageCollectorMXBean) IsJavaNil_() bool { return x == nil }
        func (x *jMemoryMXBean) IsJavaNil_() bool           { return x == nil }
        func (x *jMemoryPoolMXBean) IsJavaNil_() bool       { return x == nil }
        """;
  }
}
