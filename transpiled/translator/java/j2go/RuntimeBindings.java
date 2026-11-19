package j2go;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.Comparator;
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
import javax.lang.model.element.TypeElement;
import javax.lang.model.element.VariableElement;
import javax.lang.model.type.TypeKind;
import javax.lang.model.type.TypeMirror;

/**
 * Maps JDK / Guava / Protobuf / Gson / Args4j / XML external types, constructors, methods, and
 * constants to {@code transpiled/shims/*} and {@code j2gort}, and emits the self-contained
 * {@code <outDir>/j2gort/} bridge package.
 *
 * <p>All unhandled external constructors, static methods, instance methods, and lambda adapter
 * methods fail fast at runtime with {@code panic("UNIMPLEMENTED_EXT: ...")} instead of silently
 * returning zero values.
 */
public final class RuntimeBindings {

  public static final String SHIM_BASE = "github.com/paullepoulpe/jsgo/transpiled/shims";
  public static final String PKG_JAVALANG = SHIM_BASE + "/javalang";
  public static final String PKG_JAVAUTIL = SHIM_BASE + "/javautil";
  public static final String PKG_JAVAIO = SHIM_BASE + "/javaio";
  public static final String PKG_GUAVA = SHIM_BASE + "/guava";
  public static final String PKG_NODESTORE = SHIM_BASE + "/nodestore";

  private static final Set<String> ANY_TYPES =
      Set.of(
          "java.lang.Object",
          "java.lang.CharSequence",
          "java.io.Serializable",
          "java.lang.Cloneable",
          "java.lang.Comparable",
          "java.lang.Void",
          "java.lang.reflect.Type",
          "java.lang.annotation.Annotation");

  // Subclassed collection/map types that share the root Self_ field from AbstractCollectionBase,
  // AbstractMapBase, or ImmutableMapBase instead of duplicating a 16-byte Self_ interface at every level.
  private static final Set<String> INHERITED_SELF_COLLECTION_TYPES =
      Set.of(
          "java.util.AbstractSet",
          "java.util.HashSet",
          "java.util.EnumSet",
          "java.util.AbstractList",
          "java.util.AbstractSequentialList",
          "java.util.ArrayList",
          "java.util.AbstractQueue",
          "com.google.common.collect.ImmutableCollection",
          "com.google.common.collect.ImmutableList",
          "com.google.common.collect.ImmutableSet",
          "com.google.common.collect.ImmutableSet.CachingAsList",
          "com.google.common.collect.ImmutableSortedSetFauxverideShim",
          "com.google.common.collect.ImmutableSortedSet",
          "com.google.common.collect.Sets.SetView",
          "java.util.HashMap",
          "com.google.common.collect.Maps.IteratorBasedAbstractMap",
          "com.google.common.collect.ImmutableSortedMapFauxverideShim",
          "com.google.common.collect.ImmutableSortedMap");

  // External exception classes mapped directly to transpiled/shims/{javalang,javaio}
  private record ShimExInfo(String pkgPath, String pkgAlias, String ifaceName, String structName) {}

  private static final Map<String, ShimExInfo> SHIM_EXCEPTIONS = new LinkedHashMap<>();

  static {
    putEx("java.lang.Throwable", PKG_JAVALANG, "javalang", "Throwable", "ThrowableBase");
    putEx("java.lang.Exception", PKG_JAVALANG, "javalang", "IException", "JavaException");
    putEx(
        "java.lang.RuntimeException",
        PKG_JAVALANG,
        "javalang",
        "IRuntimeException",
        "RuntimeException");
    putEx(
        "java.lang.IllegalArgumentException",
        PKG_JAVALANG,
        "javalang",
        "IIllegalArgumentException",
        "IllegalArgumentException");
    putEx(
        "java.lang.IllegalStateException",
        PKG_JAVALANG,
        "javalang",
        "IIllegalStateException",
        "IllegalStateException");
    putEx(
        "java.lang.NullPointerException",
        PKG_JAVALANG,
        "javalang",
        "INullPointerException",
        "NullPointerException");
    putEx(
        "java.lang.IndexOutOfBoundsException",
        PKG_JAVALANG,
        "javalang",
        "IIndexOutOfBoundsException",
        "IndexOutOfBoundsException");
    putEx(
        "java.lang.ArrayIndexOutOfBoundsException",
        PKG_JAVALANG,
        "javalang",
        "IArrayIndexOutOfBoundsException",
        "ArrayIndexOutOfBoundsException");
    putEx(
        "java.lang.StringIndexOutOfBoundsException",
        PKG_JAVALANG,
        "javalang",
        "IStringIndexOutOfBoundsException",
        "StringIndexOutOfBoundsException");
    putEx(
        "java.lang.UnsupportedOperationException",
        PKG_JAVALANG,
        "javalang",
        "IUnsupportedOperationException",
        "UnsupportedOperationException");
    putEx(
        "java.lang.ClassCastException",
        PKG_JAVALANG,
        "javalang",
        "IClassCastException",
        "ClassCastException");
    putEx(
        "java.util.NoSuchElementException",
        PKG_JAVALANG,
        "javalang",
        "INoSuchElementException",
        "NoSuchElementException");
    putEx(
        "java.lang.ArithmeticException",
        PKG_JAVALANG,
        "javalang",
        "IArithmeticException",
        "ArithmeticException");
    putEx(
        "java.lang.NumberFormatException",
        PKG_JAVALANG,
        "javalang",
        "INumberFormatException",
        "NumberFormatException");
    putEx("java.lang.Error", PKG_JAVALANG, "javalang", "IError", "JavaLangError");
    putEx(
        "java.lang.AssertionError",
        PKG_JAVALANG,
        "javalang",
        "IAssertionError",
        "AssertionError");
    putEx("java.io.IOException", PKG_JAVAIO, "javaio", "IIOException", "IOException");
    putEx(
        "java.io.FileNotFoundException",
        PKG_JAVAIO,
        "javaio",
        "IFileNotFoundException",
        "FileNotFoundException");
    putEx("java.io.EOFException", PKG_JAVAIO, "javaio", "IEOFException", "EOFException");
  }

  private static void putEx(String qn, String pkgPath, String alias, String iface, String struct) {
    SHIM_EXCEPTIONS.put(qn, new ShimExInfo(pkgPath, alias, iface, struct));
  }

  private static final Map<String, String> BOXED_TYPES =
      Map.of(
          "java.lang.Boolean", "BoxBoolean",
          "java.lang.Byte", "BoxByte",
          "java.lang.Short", "BoxShort",
          "java.lang.Character", "BoxChar",
          "java.lang.Integer", "BoxInt",
          "java.lang.Long", "BoxLong",
          "java.lang.Float", "BoxFloat",
          "java.lang.Double", "BoxDouble");

  // Core external types referenced directly by j2gort helper Go files.
  private static final List<String> CORE_RUNTIME_EXT_TYPES =
      List.of(
          "java.lang.AutoCloseable",
          "java.lang.Iterable",
          "java.lang.Runnable",
          "java.lang.Runtime",
          "java.lang.StringBuilder",
          "java.lang.Thread",
          "java.lang.ThreadLocal",
          "java.lang.management.CompilationMXBean",
          "java.lang.management.GarbageCollectorMXBean",
          "java.lang.management.MemoryMXBean",
          "java.lang.management.MemoryManagerMXBean",
          "java.lang.management.MemoryPoolMXBean",
          "java.lang.management.MemoryUsage",
          "java.lang.management.PlatformManagedObject",
          "java.lang.reflect.AnnotatedElement",
          "java.math.BigInteger",
          "java.net.URI",
          "java.net.URISyntaxException",
          "java.io.BufferedInputStream",
          "java.io.BufferedOutputStream",
          "java.io.BufferedReader",
          "java.io.BufferedWriter",
          "java.io.ByteArrayInputStream",
          "java.io.ByteArrayOutputStream",
          "java.io.Closeable",
          "java.io.File",
          "java.io.FileInputStream",
          "java.io.FileOutputStream",
          "java.io.FilterInputStream",
          "java.io.FilterOutputStream",
          "java.io.Flushable",
          "java.io.InputStream",
          "java.io.InputStreamReader",
          "java.io.OutputStream",
          "java.io.OutputStreamWriter",
          "java.io.PrintStream",
          "java.io.PrintWriter",
          "java.io.Reader",
          "java.io.StringReader",
          "java.io.StringWriter",
          "java.io.Writer",
          "java.nio.charset.CharacterCodingException",
          "java.nio.charset.Charset",
          "java.nio.charset.CharsetEncoder",
          "java.nio.charset.CodingErrorAction",
          "java.nio.file.FileSystem",
          "java.nio.file.FileVisitOption",
          "java.nio.file.FileVisitResult",
          "java.nio.file.FileVisitor",
          "java.nio.file.OpenOption",
          "java.nio.file.Path",
          "java.nio.file.PathMatcher",
          "java.nio.file.StandardOpenOption",
          "java.nio.file.Watchable",
          "java.nio.file.attribute.BasicFileAttributes",
          "java.nio.file.attribute.FileTime",
          "java.util.AbstractCollection",
          "java.util.AbstractList",
          "java.util.AbstractMap",
          "java.util.AbstractMap.SimpleImmutableEntry",
          "java.util.AbstractQueue",
          "java.util.AbstractSet",
          "java.util.ArrayDeque",
          "java.util.ArrayList",
          "java.util.BitSet",
          "java.util.Collection",
          "java.util.Comparator",
          "java.util.Deque",
          "java.util.EnumMap",
          "java.util.EnumSet",
          "java.util.HashMap",
          "java.util.HashSet",
          "java.util.IdentityHashMap",
          "java.util.Iterator",
          "java.util.LinkedHashMap",
          "java.util.LinkedHashSet",
          "java.util.LinkedList",
          "java.util.List",
          "java.util.ListIterator",
          "java.util.Locale",
          "java.util.Map",
          "java.util.Map.Entry",
          "java.util.NavigableMap",
          "java.util.NavigableSet",
          "java.util.Optional",
          "java.util.PriorityQueue",
          "java.util.Queue",
          "java.util.Random",
          "java.util.ResourceBundle",
          "java.util.SequencedCollection",
          "java.util.SequencedMap",
          "java.util.SequencedSet",
          "java.util.Set",
          "java.util.SortedMap",
          "java.util.SortedSet",
          "java.util.Spliterator",
          "java.util.StringTokenizer",
          "java.util.TreeMap",
          "java.util.TreeSet",
          "java.util.concurrent.Callable",
          "java.util.concurrent.ConcurrentHashMap",
          "java.util.concurrent.ExecutionException",
          "java.util.concurrent.Executor",
          "java.util.concurrent.ExecutorService",
          "java.util.concurrent.Future",
          "java.util.concurrent.TimeUnit",
          "java.util.concurrent.atomic.AtomicIntegerArray",
          "java.util.function.BiConsumer",
          "java.util.function.BiFunction",
          "java.util.function.BinaryOperator",
          "java.util.function.Consumer",
          "java.util.function.Function",
          "java.util.function.IntFunction",
          "java.util.function.IntPredicate",
          "java.util.function.IntUnaryOperator",
          "java.util.function.Predicate",
          "java.util.function.Supplier",
          "java.util.function.ToDoubleFunction",
          "java.util.function.ToIntFunction",
          "java.util.function.ToLongFunction",
          "java.util.function.UnaryOperator",
          "java.util.logging.Level",
          "java.util.logging.Logger",
          "java.util.regex.MatchResult",
          "java.util.regex.Matcher",
          "java.util.regex.Pattern",
          "java.util.stream.BaseStream",
          "java.util.stream.Collector",
          "java.util.stream.IntStream",
          "java.util.stream.Stream",
          "java.util.stream.Stream.Builder",
          "java.util.zip.Deflater",
          "java.util.zip.DeflaterOutputStream",
          "java.util.zip.GZIPInputStream",
          "java.util.zip.GZIPOutputStream",
          "java.util.zip.ZipEntry",
          "java.util.zip.ZipFile",
          "java.util.zip.ZipInputStream",
          "javax.xml.parsers.ParserConfigurationException",
          "javax.xml.parsers.SAXParser",
          "javax.xml.parsers.SAXParserFactory",
          "org.xml.sax.Attributes",
          "org.xml.sax.ContentHandler",
          "org.xml.sax.EntityResolver",
          "org.xml.sax.InputSource",
          "org.xml.sax.SAXException",
          "org.xml.sax.XMLReader",
          "org.kohsuke.args4j.CmdLineException",
          "org.kohsuke.args4j.CmdLineParser",
          "org.kohsuke.args4j.NamedOptionDef",
          "org.kohsuke.args4j.OptionDef",
          "org.kohsuke.args4j.OptionHandlerFilter",
          "org.kohsuke.args4j.spi.FieldSetter",
          "org.kohsuke.args4j.spi.OneArgumentOptionHandler",
          "org.kohsuke.args4j.spi.OptionHandler",
          "org.kohsuke.args4j.spi.Parameters",
          "org.kohsuke.args4j.spi.Setter",
          "org.kohsuke.args4j.spi.StringOptionHandler",
          "com.google.common.base.CaseFormat",
          "com.google.common.base.CharMatcher",
          "com.google.common.base.Function",
          "com.google.common.base.Joiner",
          "com.google.common.base.MoreObjects.ToStringHelper",
          "com.google.common.base.Optional",
          "com.google.common.base.Predicate",
          "com.google.common.base.Splitter",
          "com.google.common.base.Supplier",
          "com.google.common.cache.Cache",
          "com.google.common.cache.CacheBuilder",
          "com.google.common.cache.CacheLoader",
          "com.google.common.cache.LoadingCache",
          "com.google.common.cache.RemovalListener",
          "com.google.common.cache.RemovalNotification",
          "com.google.common.collect.AbstractIterator",
          "com.google.common.collect.ArrayListMultimap",
          "com.google.common.collect.BiMap",
          "com.google.common.collect.BoundType",
          "com.google.common.collect.ComparisonChain",
          "com.google.common.collect.HashBasedTable",
          "com.google.common.collect.HashBiMap",
          "com.google.common.collect.HashMultimap",
          "com.google.common.collect.HashMultiset",
          "com.google.common.collect.ImmutableBiMap",
          "com.google.common.collect.ImmutableBiMap.Builder",
          "com.google.common.collect.ImmutableCollection",
          "com.google.common.collect.ImmutableCollection.Builder",
          "com.google.common.collect.ImmutableList",
          "com.google.common.collect.ImmutableList.Builder",
          "com.google.common.collect.ImmutableListMultimap",
          "com.google.common.collect.ImmutableListMultimap.Builder",
          "com.google.common.collect.ImmutableMap",
          "com.google.common.collect.ImmutableMap.Builder",
          "com.google.common.collect.ImmutableMultimap",
          "com.google.common.collect.ImmutableMultimap.Builder",
          "com.google.common.collect.ImmutableMultiset",
          "com.google.common.collect.ImmutableRangeSet",
          "com.google.common.collect.ImmutableSet",
          "com.google.common.collect.ImmutableSet.Builder",
          "com.google.common.collect.ImmutableSetMultimap",
          "com.google.common.collect.ImmutableSetMultimap.Builder",
          "com.google.common.collect.ImmutableSortedMap",
          "com.google.common.collect.ImmutableSortedMap.Builder",
          "com.google.common.collect.ImmutableSortedSet",
          "com.google.common.collect.ImmutableSortedSet.Builder",
          "com.google.common.collect.ImmutableTable",
          "com.google.common.collect.ImmutableTable.Builder",
          "com.google.common.collect.Interner",
          "com.google.common.collect.LinkedHashMultimap",
          "com.google.common.collect.LinkedHashMultiset",
          "com.google.common.collect.ListMultimap",
          "com.google.common.collect.Multimap",
          "com.google.common.collect.MultimapBuilder.ListMultimapBuilder",
          "com.google.common.collect.MultimapBuilder.MultimapBuilderWithKeys",
          "com.google.common.collect.MultimapBuilder.SetMultimapBuilder",
          "com.google.common.collect.Multiset",
          "com.google.common.collect.Multiset.Entry",
          "com.google.common.collect.Ordering",
          "com.google.common.collect.PeekingIterator",
          "com.google.common.collect.Range",
          "com.google.common.collect.RangeSet",
          "com.google.common.collect.SetMultimap",
          "com.google.common.collect.Sets.SetView",
          "com.google.common.collect.SortedSetMultimap",
          "com.google.common.collect.Streams.FunctionWithIndex",
          "com.google.common.collect.Table",
          "com.google.common.collect.Table.Cell",
          "com.google.common.collect.TreeMultimap",
          "com.google.common.collect.TreeMultiset",
          "com.google.common.collect.UnmodifiableIterator",
          "com.google.common.collect.UnmodifiableListIterator",
          "com.google.common.escape.Escaper",
          "com.google.common.hash.HashCode",
          "com.google.common.hash.HashFunction",
          "com.google.common.hash.Hasher",
          "com.google.common.hash.PrimitiveSink",
          "com.google.common.io.BaseEncoding",
          "com.google.common.io.CharSink",
          "com.google.common.io.CharSource",
          "com.google.common.io.FileWriteMode",
          "com.google.common.net.PercentEscaper",
          "com.google.common.util.concurrent.ListenableFuture",
          "com.google.common.util.concurrent.ListeningExecutorService",
          "com.google.gson.Gson",
          "com.google.gson.GsonBuilder",
          "com.google.gson.JsonArray",
          "com.google.gson.JsonElement",
          "com.google.gson.JsonIOException",
          "com.google.gson.JsonNull",
          "com.google.gson.JsonObject",
          "com.google.gson.JsonParseException",
          "com.google.gson.JsonPrimitive",
          "com.google.gson.JsonSyntaxException",
          "com.google.gson.TypeAdapter",
          "com.google.gson.reflect.TypeToken",
          "com.google.gson.stream.JsonReader",
          "com.google.gson.stream.JsonToken",
          "com.google.gson.stream.JsonWriter",
          "com.google.protobuf.ByteString",
          "com.google.protobuf.ByteString.Output",
          "com.google.protobuf.CodedInputStream",
          "com.google.protobuf.Descriptors.FieldDescriptor",
          "com.google.protobuf.ExtensionRegistry",
          "com.google.protobuf.ExtensionRegistryLite",
          "com.google.protobuf.InvalidProtocolBufferException",
          "com.google.protobuf.ProtocolStringList",
          "com.google.protobuf.TextFormat.Printer",
          "com.google.javascript.jscomp.RequirementScopeEntry.Builder");

  private final TypeResolver resolver;
  private final Set<TypeElement> subclassedExternals = new LinkedHashSet<>();

  public RuntimeBindings(TypeResolver resolver) {
    this.resolver = resolver;
  }

  public TypeResolver resolver() {
    return resolver;
  }

  public String j2gortImportPath() {
    return resolver.planner.getGoModuleBase() + "/j2gort";
  }

  public String j2gortAlias(PackagePlanner.ImportTracker imports) {
    return imports.ext(j2gortImportPath(), "j2gort");
  }

  public String javalangAlias(PackagePlanner.ImportTracker imports) {
    return imports.ext(PKG_JAVALANG, "javalang");
  }

  public String guavaAlias(PackagePlanner.ImportTracker imports) {
    return imports.ext(PKG_GUAVA, "guava");
  }

  public String nodestoreAlias(PackagePlanner.ImportTracker imports) {
    return imports.ext(PKG_NODESTORE, "nodestore");
  }

  public boolean mapsToAny(TypeElement te) {
    return ANY_TYPES.contains(te.getQualifiedName().toString());
  }

  public boolean isShimException(TypeElement te) {
    return SHIM_EXCEPTIONS.containsKey(te.getQualifiedName().toString());
  }

  public boolean isBoxedPrimitive(TypeElement te) {
    return BOXED_TYPES.containsKey(te.getQualifiedName().toString());
  }

  public boolean BOXED_TYPES_CONTAINS(String qn) {
    return BOXED_TYPES.containsKey(qn);
  }

  public String getBoxedTypeName(TypeElement te) {
    return BOXED_TYPES.get(te.getQualifiedName().toString());
  }

  public void markSubclassedExternal(TypeElement te) {
    if (te != null && !isShimException(te) && !mapsToAny(te)) {
      subclassedExternals.add(te);
      TypeElement sup = resolver.getSuperclassElement(te);
      if (sup != null && !resolver.isSourceType(sup)) {
        markSubclassedExternal(sup);
      }
    }
  }

  public static boolean isValueStructExternal(String qn) {
    return "com.google.protobuf.ByteString".equals(qn);
  }

  public boolean isValueStructExternal(TypeElement te) {
    return te != null && isValueStructExternal(te.getQualifiedName().toString());
  }

  public boolean isSubclassedExternal(TypeElement te) {
    String qn = te.getQualifiedName().toString();
    if (SHIM_EXCEPTIONS.containsKey(qn)) {
      return true;
    }
    if (te.getKind() == ElementKind.ENUM
        || BOXED_TYPES.containsKey(qn)
        || isValueStructExternal(qn)
        || qn.equals("java.lang.Class")
        || qn.equals("java.lang.String")) {
      return false;
    }
    if (te.getModifiers().contains(Modifier.ABSTRACT) && te.getKind() == ElementKind.CLASS) {
      return true;
    }
    if (resolver.isThrowable(te)) {
      return true;
    }
    // External base classes explicitly embedded by runtime helpers or Closure subclasses
    if (qn.equals("java.util.ArrayList")
        || qn.equals("java.util.HashSet")
        || qn.equals("java.util.HashMap")
        || qn.equals("com.google.common.collect.ImmutableSet.Builder")
        || qn.equals("com.google.common.collect.ImmutableMap.Builder")
        || qn.equals("com.google.common.collect.ImmutableMultimap.Builder")
        || qn.equals("java.util.zip.DeflaterOutputStream")
        || qn.equals("java.util.zip.GZIPOutputStream")
        || qn.equals("java.io.FilterInputStream")
        || qn.equals("java.io.FilterOutputStream")
        || qn.equals("org.kohsuke.args4j.OptionDef")
        || qn.equals("org.kohsuke.args4j.spi.OptionHandler")
        || qn.equals("org.kohsuke.args4j.spi.OneArgumentOptionHandler")
        || qn.equals("org.kohsuke.args4j.spi.StringOptionHandler")
        || qn.equals("com.google.gson.reflect.TypeToken")) {
      return true;
    }
    for (TypeElement other : resolver.getReferencedExternalTypes()) {
      if (other != te) {
        TypeElement sup = resolver.getSuperclassElement(other);
        if (te.equals(sup)) {
          return true;
        }
      }
    }
    return subclassedExternals.contains(te);
  }

  public boolean isGoInterfaceExternal(TypeElement te) {
    String qn = te.getQualifiedName().toString();
    if (ANY_TYPES.contains(qn)
        || SHIM_EXCEPTIONS.containsKey(qn)
        || qn.equals("java.lang.Number")) {
      return true;
    }
    if (qn.equals("java.lang.String")
        || qn.equals("java.lang.Class")
        || isValueStructExternal(qn)
        || BOXED_TYPES.containsKey(qn)) {
      return false;
    }
    if (te.getKind() == ElementKind.INTERFACE) {
      return true;
    }
    return isSubclassedExternal(te);
  }

  public String extGoTypeName(TypeElement te) {
    String qn = te.getQualifiedName().toString();
    if (BOXED_TYPES.containsKey(qn)) {
      return BOXED_TYPES.get(qn);
    }
    if (qn.equals("java.lang.Number")) {
      return "Number";
    }
    return "Ext_" + GoNames.cleanChars(qn.replace('.', '_'));
  }

  public String canonicalTypeKey(TypeElement te) {
    String qn = te.getQualifiedName().toString();
    if (ANY_TYPES.contains(qn)) return "any";
    if (qn.equals("java.lang.String")) return "string";
    if (qn.equals("java.lang.Class")) return "javalang.Class";
    if (qn.equals("java.lang.Number")) return "j2gort.Number";
    if (BOXED_TYPES.containsKey(qn)) return "j2gort." + BOXED_TYPES.get(qn);
    ShimExInfo ex = SHIM_EXCEPTIONS.get(qn);
    if (ex != null) {
      return ex.pkgAlias + "." + ex.ifaceName;
    }
    return "j2gort." + extGoTypeName(te);
  }

  public String goType(TypeElement te, PackagePlanner.ImportTracker imports) {
    String qn = te.getQualifiedName().toString();
    if (ANY_TYPES.contains(qn)) {
      return "any";
    }
    if (qn.equals("java.lang.String")) {
      return "string";
    }
    if (qn.equals("java.lang.Class")) {
      String jl = javalangAlias(imports);
      return "*" + (jl.isEmpty() ? "" : jl + ".") + "Class";
    }
    ShimExInfo ex = SHIM_EXCEPTIONS.get(qn);
    if (ex != null) {
      String al = imports.ext(ex.pkgPath, ex.pkgAlias);
      return (al.isEmpty() ? "" : al + ".") + ex.ifaceName;
    }
    String rt = j2gortAlias(imports);
    String prefix = rt.isEmpty() ? "" : (rt + ".");
    if (qn.equals("java.lang.Number")) {
      return prefix + "Number";
    }
    if (BOXED_TYPES.containsKey(qn)) {
      return "*" + prefix + BOXED_TYPES.get(qn);
    }
    String extName = extGoTypeName(te);
    if (te.getKind() == ElementKind.INTERFACE || isSubclassedExternal(te) || isValueStructExternal(qn)) {
      return prefix + extName;
    }
    return "*" + prefix + extName;
  }

  public String goStructType(TypeElement te, PackagePlanner.ImportTracker imports) {
    String qn = te.getQualifiedName().toString();
    ShimExInfo ex = SHIM_EXCEPTIONS.get(qn);
    if (ex != null) {
      String al = imports.ext(ex.pkgPath, ex.pkgAlias);
      return (al.isEmpty() ? "" : al + ".") + ex.structName;
    }
    String rt = j2gortAlias(imports);
    String prefix = rt.isEmpty() ? "" : (rt + ".");
    if (qn.equals("java.lang.Number")) {
      return prefix + "NumberBase";
    }
    if (BOXED_TYPES.containsKey(qn)) {
      return prefix + BOXED_TYPES.get(qn);
    }
    String extName = extGoTypeName(te);
    return isSubclassedExternal(te) ? (prefix + extName + "Base") : (prefix + extName);
  }

  /**
   * Emits `<outDir>/j2gort/{runtime,ext_bindings,ext_args4j,ext_protobuf,ext_collections_rt,ext_io_nio_xml_gson}.go`.
   */
  public void emitRuntimePackage(Path outDir) throws IOException {
    Path rtDir = outDir.resolve("j2gort");
    Files.createDirectories(rtDir);
    String extBindingsGo = buildExtBindingsGo();
    Files.writeString(rtDir.resolve("runtime.go"), buildCoreRuntimeGo(), StandardCharsets.UTF_8);
    Files.writeString(rtDir.resolve("ext_bindings.go"), extBindingsGo, StandardCharsets.UTF_8);
    Files.writeString(
        rtDir.resolve("ext_args4j.go"),
        Args4jEmitter.buildRuntimeArgs4jGo(),
        StandardCharsets.UTF_8);
    Files.writeString(
        rtDir.resolve("ext_protobuf.go"),
        RuntimeProtobufBindings.buildRuntimeProtobufGo(this),
        StandardCharsets.UTF_8);
    Files.writeString(
        rtDir.resolve("ext_collections_rt.go"),
        RuntimeCollectionsBindings.buildCollectionsRuntimeGo(this),
        StandardCharsets.UTF_8);
    Files.writeString(
        rtDir.resolve("ext_io_nio_xml_gson.go"),
        RuntimeIoNioXmlGsonBindings.buildIoNioXmlGsonRuntimeGo(this),
        StandardCharsets.UTF_8);
  }

  private String buildCoreRuntimeGo() {
    return buildCoreRuntimeGoPart1()
        + buildCoreRuntimeGoPart1Predicates()
        + buildCoreRuntimeGoPart1Tail()
        + buildCoreRuntimeGoPart1B()
        + buildCoreRuntimeGoPart2ListDeque()
        + buildCoreRuntimeGoPart2Set()
        + buildCoreRuntimeGoPart2Map()
        + buildCoreRuntimeGoPart2Gen()
        + buildCoreRuntimeGoPart2MapOps();
  }

  private String buildCoreRuntimeGoPart1() {
    return """
        // Clean-room runtime support bridge for j2go mechanical Java-to-Go output.
        package j2gort

        import (
        \t"fmt"
        \t"math"
        \t"math/big"
        \t"math/bits"
        \t"reflect"
        \t"regexp"
        \t"runtime"
        \t"sort"
        \t"strconv"
        \t"strings"
        \t"sync"
        \t"sync/atomic"
        \t"unicode"
        \t"unicode/utf16"
        \t"unicode/utf8"
        \t"unique"

        \t"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
        \t"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
        \t"github.com/paullepoulpe/jsgo/transpiled/shims/javaregex"
        \t"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
        \tshimsres "github.com/paullepoulpe/jsgo/transpiled/shims/resources"
        )

        var (
        \t_ = javaio.System_Out
        \t_ = javaregex.StringSplit
        \t_ = shimsres.ResourceExists
        \t_ = nodestore.NilNode
        \t_ = big.NewInt
        \t_ = regexp.Compile
        \t_ = sort.SliceStable
        \t_ = atomic.LoadUint32
        \t_ = unicode.IsDigit
        \t_ = utf16.DecodeRune
        \t_ = utf8.DecodeRune
        \t_ = unique.Make[string]
        )

        func LenientFormat(template string, args ...any) string {
        \tif template == NullString { template = "null" }
        \tvar sb strings.Builder
        \ti, argIdx := 0, 0
        \tfor i < len(template) {
        \t\tp := strings.Index(template[i:], "%s")
        \t\tif p < 0 || argIdx >= len(args) {
        \t\t\tsb.WriteString(template[i:])
        \t\t\ti = len(template)
        \t\t\tbreak
        \t\t}
        \t\tsb.WriteString(template[i : i+p])
        \t\tsb.WriteString(StrVal(args[argIdx]))
        \t\targIdx++
        \t\ti += p + 2
        \t}
        \tif argIdx < len(args) {
        \t\tsb.WriteString(" [")
        \t\tfor j := argIdx; j < len(args); j++ {
        \t\t\tif j > argIdx { sb.WriteString(", ") }
        \t\t\tsb.WriteString(StrVal(args[j]))
        \t\t}
        \t\tsb.WriteByte(']')
        \t}
        \treturn sb.String()
        }

        func normFmtArg(v any) any {
        \tif v == nil {
        \t\treturn nil
        \t}
        \tif s, ok := v.(string); ok {
        \t\tif s == NullString {
        \t\t\treturn nil
        \t\t}
        \t\treturn s
        \t}
        \tif rn, ok := v.(jRawNodeProvider); ok {
        \t\tif rn.RawNodeId_() == 0 {
        \t\t\treturn nil
        \t\t}
        \t\treturn v
        \t}
        \treturn ToAny(v)
        }

        //go:noinline
        func checkStateFmtFail(tmpl string, args ...any) {
        \tfor i, a := range args {
        \t\targs[i] = normFmtArg(a)
        \t}
        \tpanic(javalang.NewIllegalStateException(LenientFormat(tmpl, args...)))
        }

        //go:noinline
        func checkArgFmtFail(tmpl string, args ...any) {
        \tfor i, a := range args {
        \t\targs[i] = normFmtArg(a)
        \t}
        \tpanic(javalang.NewIllegalArgumentException(LenientFormat(tmpl, args...)))
        }

        func CheckState(b bool) {
        \tif !b { panic(javalang.NewIllegalStateException("")) }
        }
        func CheckStateStr(b bool, msg string) {
        \tif !b { panic(javalang.NewIllegalStateException(StrValString(msg))) }
        }
        func CheckStateMsg(b bool, msg any) {
        \tif !b { panic(javalang.NewIllegalStateException(StrVal(msg))) }
        }
        func CheckStateFmt(b bool, tmpl string, args ...any) {
        \tif !b { panic(javalang.NewIllegalStateException(LenientFormat(tmpl, args...))) }
        }
        func CheckStateFmt1[A any](b bool, tmpl string, a0 A) {
        \tif !b { checkStateFmtFail(tmpl, a0) }
        }
        func CheckStateFmt2[A, B any](b bool, tmpl string, a0 A, a1 B) {
        \tif !b { checkStateFmtFail(tmpl, a0, a1) }
        }
        func CheckStateFmt3[A, B, C any](b bool, tmpl string, a0 A, a1 B, a2 C) {
        \tif !b { checkStateFmtFail(tmpl, a0, a1, a2) }
        }
        func CheckStateFmt4[A, B, C, D any](b bool, tmpl string, a0 A, a1 B, a2 C, a3 D) {
        \tif !b { checkStateFmtFail(tmpl, a0, a1, a2, a3) }
        }
        func CheckArgument(b bool) {
        \tif !b { panic(javalang.NewIllegalArgumentException("")) }
        }
        func CheckArgumentStr(b bool, msg string) {
        \tif !b { panic(javalang.NewIllegalArgumentException(StrValString(msg))) }
        }
        func CheckArgumentMsg(b bool, msg any) {
        \tif !b { panic(javalang.NewIllegalArgumentException(StrVal(msg))) }
        }
        func CheckArgumentFmt(b bool, tmpl string, args ...any) {
        \tif !b { panic(javalang.NewIllegalArgumentException(LenientFormat(tmpl, args...))) }
        }
        func CheckArgumentFmt1[A any](b bool, tmpl string, a0 A) {
        \tif !b { checkArgFmtFail(tmpl, a0) }
        }
        func CheckArgumentFmt2[A, B any](b bool, tmpl string, a0 A, a1 B) {
        \tif !b { checkArgFmtFail(tmpl, a0, a1) }
        }
        func CheckArgumentFmt3[A, B, C any](b bool, tmpl string, a0 A, a1 B, a2 C) {
        \tif !b { checkArgFmtFail(tmpl, a0, a1, a2) }
        }
        func CheckArgumentFmt4[A, B, C, D any](b bool, tmpl string, a0 A, a1 B, a2 C, a3 D) {
        \tif !b { checkArgFmtFail(tmpl, a0, a1, a2, a3) }
        }
        func CheckElementIndex(index, size int32) int32 {
        \tif index < 0 || index >= size {
        \t\tpanic(javalang.NewIndexOutOfBoundsExceptionForIndex(int(index), int(size)))
        \t}
        \treturn index
        }
        func CheckPositionIndex(index, size int32) int32 {
        \tif index < 0 || index > size {
        \t\tpanic(javalang.NewIndexOutOfBoundsExceptionForIndex(int(index), int(size)))
        \t}
        \treturn index
        }
        func CheckPositionIndexes(start, end, size int32) {
        \tif start < 0 || end < start || end > size {
        \t\tpanic(javalang.NewIndexOutOfBoundsException(fmt.Sprintf("start=%d end=%d size=%d", start, end, size)))
        \t}
        }
        func ArraysFill[T any](a []T, val T) {
        \tif a == nil {
        \t\tThrowNPE()
        \t}
        \tfor i := range a {
        \t\ta[i] = val
        \t}
        }
        func ArraysFillZero[T any](a []T) {
        \tif a == nil {
        \t\tThrowNPE()
        \t}
        \tclear(a)
        }
        func ArraysFillRange[T any](a []T, from, to int32, val T) {
        \tif a == nil {
        \t\tThrowNPE()
        \t}
        \tCheckPositionIndexes(from, to, int32(len(a)))
        \tfor i := from; i < to; i++ {
        \t\ta[i] = val
        \t}
        }
        func ArraysFillRangeZero[T any](a []T, from, to int32) {
        \tif a == nil {
        \t\tThrowNPE()
        \t}
        \tCheckPositionIndexes(from, to, int32(len(a)))
        \tclear(a[from:to])
        }

        // NullString is a WTF-8/UTF-8 impossible byte sentinel representing Java `(String) null`
        // inside a value-typed Go `string` with zero heap allocation.
        const NullString = "\\xff\\x00j2go_null"

        func IsNullString(s string) bool {
        \treturn s == NullString
        }

        type jJavaNilChecker interface {
        \tIsJavaNil_() bool
        }

        func isJavaNilAny(v any) bool {
        \tif v == nil {
        \t\treturn true
        \t}
        \tif nc, ok := v.(jJavaNilChecker); ok {
        \t\treturn nc.IsJavaNil_()
        \t}
        \tswitch x := v.(type) {
        \tcase string:
        \t\treturn x == NullString
        \tcase nodestore.Node:
        \t\treturn x == 0
        \tcase []any:
        \t\treturn x == nil
        \tcase []string:
        \t\treturn x == nil
        \tcase []byte:
        \t\treturn x == nil
        \tcase []int8:
        \t\treturn x == nil
        \tcase []int32:
        \t\treturn x == nil
        \tcase []uint16:
        \t\treturn x == nil
        \tcase bool, int32, int64, uint16, int8, int16, float32, float64:
        \t\treturn false
        \t}
        \treturn javalang.IsNil(v)
        }

        var asciiBoxedStrings_ [128]any
        var asciiInternStrings_ [128]string

        func init() {
        \tfor i := 0; i < 128; i++ {
        \t\ts := string(byte(i))
        \t\tasciiInternStrings_[i] = s
        \t\tasciiBoxedStrings_[i] = s
        \t}
        }

        func BoxString(s string) any {
        \tif len(s) == 1 && s[0] < 128 {
        \t\treturn asciiBoxedStrings_[s[0]]
        \t}
        \tif s == NullString {
        \t\treturn nil
        \t}
        \treturn s
        }

        func CastString(v any) string {
        \tif v == nil {
        \t\treturn NullString
        \t}
        \tif s, ok := v.(string); ok {
        \t\treturn s
        \t}
        \tif isJavaNilAny(v) {
        \t\treturn NullString
        \t}
        \tpanic(javalang.NewClassCastException(fmt.Sprintf("cannot cast %T to java.lang.String", v)))
        }

        type jRawNodeProvider interface {
        \tRawNodeId_() uint32
        }

        var globalNodeBoxer func(uint32) any

        func RegisterNodeBoxer(fn func(uint32) any) {
        \tglobalNodeBoxer = fn
        }

        func RegisterNodeChunkFiller(_ func(uint32, *[4096]any)) {}

        func BoxNode[N ~uint32](n N) any {
        \tif n == 0 {
        \t\treturn nil
        \t}
        \treturn n
        }

        func boxNodeById(u uint32) any {
        \tif u == 0 {
        \t\treturn nil
        \t}
        \tif globalNodeBoxer != nil {
        \t\treturn globalNodeBoxer(u)
        \t}
        \treturn nodestore.Node(u)
        }

        var globalColorIdBoxer func(int64) any

        func RegisterColorIdBoxer(fn func(int64) any) {
        \tglobalColorIdBoxer = fn
        }

        func boxColorId(raw int64) any {
        \tif globalColorIdBoxer != nil {
        \t\treturn globalColorIdBoxer(raw)
        \t}
        \treturn raw
        }

        type RawColorIdKey interface {
        \tcomparable
        \tRawColorId_() (int64, bool)
        }

        type jRawColorIdProvider interface {
        \tRawColorId_() (int64, bool)
        }

        type FastEqValKey interface {
        \tcomparable
        \tFastValHash_() (int32, bool, bool)
        }

        type jSingleNonNullElem interface {
        \tsingleNonNullElem_() any
        }

        func CastNode[N ~uint32](v any) N {
        \tif v == nil {
        \t\treturn 0
        \t}
        \tif n, ok := v.(N); ok {
        \t\treturn n
        \t}
        \tif rn, ok := v.(jRawNodeProvider); ok {
        \t\treturn N(rn.RawNodeId_())
        \t}
        \tif n2, ok := v.(nodestore.Node); ok {
        \t\treturn N(n2)
        \t}
        \tif isJavaNilAny(v) {
        \t\treturn 0
        \t}
        \tpanic(javalang.NewClassCastException(fmt.Sprintf("cannot cast %T to Node", v)))
        }

        // IdentityComparable marks generated Java classes that inherit Object.equals and
        // Object.hashCode (guaranteed comparable concrete struct pointers).
        type IdentityComparable interface {
        \tHasIdentityEquals_()
        }

        // SliceIterable allows zero-allocation enhanced-for iteration over backing slices.
        type SliceIterable interface {
        \tIterSlice_() []any
        }

        // StrSliceIterable allows zero-box enhanced-for iteration over unboxed string slices.
        type StrSliceIterable interface {
        \tIterStrSlice_() ([]string, bool)
        }

        type FastStrList interface {
        \tListGetStr_(int32) string
        }

        func IterSlice(v any) ([]any, bool) {
        \tif v == nil {
        \t\treturn nil, false
        \t}
        \tif si, ok := v.(SliceIterable); ok {
        \t\tif nc, ok2 := v.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
        \t\t\treturn nil, false
        \t\t}
        \t\treturn si.IterSlice_(), true
        \t}
        \treturn nil, false
        }

        func IfaceEquals[I interface{ IsJavaNil_() bool }](a, b I) bool {
        \tif any(a) == any(b) {
        \t\treturn true
        \t}
        \tif any(a) == nil {
        \t\treturn b.IsJavaNil_()
        \t}
        \tif any(b) == nil {
        \t\treturn a.IsJavaNil_()
        \t}
        \treturn a.IsJavaNil_() && b.IsJavaNil_()
        }

        func ObjEqStr(a, b string) bool { return a == b }

        func ObjEqNode[N ~uint32](a, b N) bool { return a == b }

        func ObjEqVal[T comparable](a, b T) bool { return a == b }

        func ObjEqPtrIdent[T any](a, b *T) bool { return a == b }

        func ObjEqCustomPtr[T any, P interface {
        \t*T
        \tEquals(any) bool
        }](a, b P) bool {
        \tif a == b {
        \t\treturn true
        \t}
        \tif a == nil || b == nil {
        \t\treturn false
        \t}
        \treturn a.Equals(b)
        }

        func ObjEqIface[T interface {
        \tcomparable
        \tIsJavaNil_() bool
        \tEquals(any) bool
        }](a, b T) bool {
        \tif a == b {
        \t\treturn true
        \t}
        \tvar z T
        \taNil := a == z || a.IsJavaNil_()
        \tbNil := b == z || b.IsJavaNil_()
        \tif aNil || bNil {
        \t\treturn aNil && bNil
        \t}
        \treturn a.Equals(b)
        }

        func ObjEqNilIface[T interface {
        \tcomparable
        \tIsJavaNil_() bool
        }](a, b T) bool {
        \tif a == b {
        \t\treturn true
        \t}
        \tvar z T
        \taNil := a == z || a.IsJavaNil_()
        \tbNil := b == z || b.IsJavaNil_()
        \tif aNil || bNil {
        \t\treturn aNil && bNil
        \t}
        \tif eq, ok := any(a).(jEqualer); ok {
        \t\treturn eq.Equals(b)
        \t}
        \treturn false
        }

        func CheckNotNullIface[I interface{ IsJavaNil_() bool }](v I) I {
        \tif any(v) == nil || v.IsJavaNil_() {
        \t\tThrowNPE()
        \t}
        \treturn v
        }

        // AsIface converts a concrete pointer or interface to target interface I while turning
        // typed nil pointers into untyped Go `nil` (preventing Go's typed-nil interface trap).
        func AsIface[I any](v any) I {
        \tvar zero I
        \tif v == nil {
        \t\treturn zero
        \t}
        \tif res, ok := v.(I); ok {
        \t\tif nc, ok2 := v.(jJavaNilChecker); ok2 {
        \t\t\tif nc.IsJavaNil_() {
        \t\t\t\treturn zero
        \t\t\t}
        \t\t\treturn res
        \t\t}
        \t\tif isJavaNilAny(v) {
        \t\t\treturn zero
        \t\t}
        \t\treturn res
        \t}
        \tif isJavaNilAny(v) {
        \t\treturn zero
        \t}
        \tpanic(javalang.NewClassCastException(fmt.Sprintf("cannot convert %T to %T", v, zero)))
        }

        func ToAny(v any) any {
        \tif v == nil {
        \t\treturn nil
        \t}
        \tif nc, ok := v.(jJavaNilChecker); ok {
        \t\tif nc.IsJavaNil_() {
        \t\t\treturn nil
        \t\t}
        \t\treturn v
        \t}
        \tswitch x := v.(type) {
        \tcase string:
        \t\tif x == NullString {
        \t\t\treturn nil
        \t\t}
        \t\treturn v
        \tcase nodestore.Node:
        \t\tif x == 0 {
        \t\t\treturn nil
        \t\t}
        \t\treturn v
        \tcase bool, int32, int64, uint16, int8, int16, float32, float64:
        \t\treturn v
        \tcase []any:
        \t\tif x == nil {
        \t\t\treturn nil
        \t\t}
        \t\treturn v
        \tcase []string:
        \t\tif x == nil {
        \t\t\treturn nil
        \t\t}
        \t\treturn v
        \tcase []byte:
        \t\tif x == nil {
        \t\t\treturn nil
        \t\t}
        \t\treturn v
        \tcase []int8:
        \t\tif x == nil {
        \t\t\treturn nil
        \t\t}
        \t\treturn v
        \tcase []int32:
        \t\tif x == nil {
        \t\t\treturn nil
        \t\t}
        \t\treturn v
        \tcase []uint16:
        \t\tif x == nil {
        \t\t\treturn nil
        \t\t}
        \t\treturn v
        \t}
        \tif javalang.IsNil(v) {
        \t\treturn nil
        \t}
        \treturn v
        }

        func PtrToAny[T any](p *T) any {
        \tif p == nil {
        \t\treturn nil
        \t}
        \treturn p
        }

        func SliceToAny[E any](s []E) any {
        \tif s == nil {
        \t\treturn nil
        \t}
        \treturn s
        }

        func BoxNullableVal[T comparable](v T) any {
        \tvar zero T
        \tif v == zero {
        \t\treturn nil
        \t}
        \treturn v
        }

        // Cast performs a Java checked reference cast `(T) v`, returning zero for nil and
        // throwing ClassCastException on mismatch.
        func Cast[T any](v any) T {
        \tvar zero T
        \tif v == nil {
        \t\treturn zero
        \t}
        \tif res, ok := v.(T); ok {
        \t\tif nc, ok2 := v.(jJavaNilChecker); ok2 {
        \t\t\tif nc.IsJavaNil_() {
        \t\t\t\treturn zero
        \t\t\t}
        \t\t\treturn res
        \t\t}
        \t\tif isJavaNilAny(v) {
        \t\t\treturn zero
        \t\t}
        \t\treturn res
        \t}
        \tif isJavaNilAny(v) {
        \t\treturn zero
        \t}
        \tpanic(javalang.NewClassCastException(fmt.Sprintf("cannot cast %T to %T", v, zero)))
        }

        func CastPtr[T any](v any) *T {
        \tif v == nil {
        \t\treturn nil
        \t}
        \tif res, ok := v.(*T); ok {
        \t\treturn res
        \t}
        \tif isJavaNilAny(v) {
        \t\treturn nil
        \t}
        \tvar zero *T
        \tpanic(javalang.NewClassCastException(fmt.Sprintf("cannot cast %T to %T", v, zero)))
        }

        func UncheckedCast[T any](v any) T {
        \tvar zero T
        \tif v == nil {
        \t\treturn zero
        \t}
        \tif res, ok := v.(T); ok {
        \t\tif nc, ok2 := v.(jJavaNilChecker); ok2 {
        \t\t\tif nc.IsJavaNil_() {
        \t\t\t\treturn zero
        \t\t\t}
        \t\t\treturn res
        \t\t}
        \t\tif isJavaNilAny(v) {
        \t\t\treturn zero
        \t\t}
        \t\treturn res
        \t}
        \treturn zero
        }

        func TryCast[T any](v any) (T, bool) {
        \tvar zero T
        \tif v == nil {
        \t\treturn zero, false
        \t}
        \tres, ok := v.(T)
        \tif !ok {
        \t\treturn zero, false
        \t}
        \tif nc, ok2 := v.(jJavaNilChecker); ok2 {
        \t\tif nc.IsJavaNil_() {
        \t\t\treturn zero, false
        \t\t}
        \t\treturn res, true
        \t}
        \tif isJavaNilAny(v) {
        \t\treturn zero, false
        \t}
        \treturn res, true
        }

        func TryCastString(v any) (string, bool) {
        \tif v == nil {
        \t\treturn NullString, false
        \t}
        \tif s, ok := v.(string); ok && s != NullString {
        \t\treturn s, true
        \t}
        \treturn NullString, false
        }

        func InstanceOf[T any](v any) bool {
        \tif v == nil {
        \t\treturn false
        \t}
        \tif _, ok := v.(T); !ok {
        \t\treturn false
        \t}
        \tif nc, ok2 := v.(jJavaNilChecker); ok2 {
        \t\treturn !nc.IsJavaNil_()
        \t}
        \treturn !isJavaNilAny(v)
        }

        func InstanceOfBind[T any](v any, out *T) bool {
        \tif v == nil {
        \t\treturn false
        \t}
        \tt, ok := v.(T)
        \tif !ok {
        \t\treturn false
        \t}
        \tif nc, ok2 := v.(jJavaNilChecker); ok2 {
        \t\tif nc.IsJavaNil_() {
        \t\t\treturn false
        \t\t}
        \t} else if isJavaNilAny(v) {
        \t\treturn false
        \t}
        \t*out = t
        \treturn true
        }

        func InstanceOfString(v any) bool {
        \tif v == nil {
        \t\treturn false
        \t}
        \ts, ok := v.(string)
        \treturn ok && s != NullString
        }

        func InstanceOfStringBind(v any, out *string) bool {
        \tif v == nil {
        \t\treturn false
        \t}
        \tif s, ok := v.(string); ok && s != NullString {
        \t\t*out = s
        \t\treturn true
        \t}
        \treturn false
        }

        // RefEquals implements Java reference `==` across interfaces, pointers, slices, and boxed values.
        func RefEquals(a, b any) bool {
        \tif a == nil && b == nil {
        \t\treturn true
        \t}
        \tif a == nil {
        \t\treturn isJavaNilAny(b)
        \t}
        \tif b == nil {
        \t\treturn isJavaNilAny(a)
        \t}
        \tif ba, ok := a.(Ext_com_google_protobuf_ByteString); ok {
        \t\tif bb, ok2 := b.(Ext_com_google_protobuf_ByteString); ok2 {
        \t\t\treturn ByteStringRefEquals(ba, bb)
        \t\t}
        \t\treturn !ba.valid && isJavaNilAny(b)
        \t}
        \tif bb, ok := b.(Ext_com_google_protobuf_ByteString); ok {
        \t\treturn !bb.valid && isJavaNilAny(a)
        \t}
        \tncA, okA := a.(jJavaNilChecker)
        \tncB, okB := b.(jJavaNilChecker)
        \tif okA && okB {
        \t\tna := ncA.IsJavaNil_()
        \t\tnb := ncB.IsJavaNil_()
        \t\tif na || nb {
        \t\t\treturn na && nb
        \t\t}
        \t\tta := reflect.TypeOf(a)
        \t\ttb := reflect.TypeOf(b)
        \t\tif ta == tb {
        \t\t\tif ta.Comparable() {
        \t\t\t\treturn a == b
        \t\t\t}
        \t\t\treturn false
        \t\t}
        \t\tif ta.Kind() == reflect.Ptr && tb.Kind() == reflect.Ptr {
        \t\t\treturn reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
        \t\t}
        \t\treturn false
        \t}
        \tna := isJavaNilAny(a)
        \tnb := isJavaNilAny(b)
        \tif na || nb {
        \t\treturn na && nb
        \t}
        \tswitch xa := a.(type) {
        \tcase string:
        \t\txb, ok := b.(string)
        \t\treturn ok && xa == xb
        \tcase nodestore.Node:
        \t\txb, ok := b.(nodestore.Node)
        \t\treturn ok && xa == xb
        \t}
        \tta := reflect.TypeOf(a)
        \ttb := reflect.TypeOf(b)
        \tif ta == tb {
        \t\tif ta.Kind() == reflect.Slice {
        \t\t\tra := reflect.ValueOf(a)
        \t\t\trb := reflect.ValueOf(b)
        \t\t\treturn ra.Pointer() == rb.Pointer() && ra.Len() == rb.Len()
        \t\t}
        \t\tif ta.Comparable() {
        \t\t\treturn a == b
        \t\t}
        \t\treturn false
        \t}
        \tif ta.Kind() == reflect.Ptr && tb.Kind() == reflect.Ptr {
        \t\treturn reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
        \t}
        \treturn false
        }

        func SliceRefEquals[T any](a, b []T) bool {
        \tif (a == nil) != (b == nil) || len(a) != len(b) || cap(a) != cap(b) {
        \t\treturn false
        \t}
        \tif a == nil {
        \t\treturn true
        \t}
        \tif cap(a) > 0 {
        \t\treturn &a[:1][0] == &b[:1][0]
        \t}
        \treturn true
        }

        func ThrowNPE() {
        \tpanic(javalang.NewNullPointerException())
        }

        func CheckNotNullString(s string) string {
        \tif s == NullString {
        \t\tThrowNPE()
        \t}
        \treturn s
        }

        func CheckNotNullNode[N ~uint32](n N) N {
        \tif n == 0 {
        \t\tThrowNPE()
        \t}
        \treturn n
        }

        func CheckNotNullEnum[E ~uint8](v E) E {
        \tif v == 0 {
        \t\tThrowNPE()
        \t}
        \treturn v
        }

        func CheckNotNullVal[T comparable](v T) T {
        \tvar zero T
        \tif v == zero {
        \t\tThrowNPE()
        \t}
        \treturn v
        }

        func CheckNotNullPtr[T any](p *T) *T {
        \tif p == nil {
        \t\tThrowNPE()
        \t}
        \treturn p
        }

        func CheckNotNullSlice[T any](s []T) []T {
        \tif s == nil {
        \t\tThrowNPE()
        \t}
        \treturn s
        }

        func CheckNotNull[T any](v T) T {
        \tif any(v) == nil {
        \t\tThrowNPE()
        \t}
        \tif nc, ok := any(v).(jJavaNilChecker); ok {
        \t\tif nc.IsJavaNil_() {
        \t\t\tThrowNPE()
        \t\t}
        \t\treturn v
        \t}
        \tif isJavaNilAny(v) {
        \t\tThrowNPE()
        \t}
        \treturn v
        }

        func EnumSwitchOrd[E interface {
        \tcomparable
        \tOrdinal() int32
        }](e E) int32 {
        \tvar zero E
        \tif e == zero {
        \t\treturn -1
        \t}
        \treturn e.Ordinal()
        }

        func BoolToStr(b bool) string {
        \tif b {
        \t\treturn "true"
        \t}
        \treturn "false"
        }

        func StrValString(s string) string {
        \tif s == NullString {
        \t\treturn "null"
        \t}
        \treturn s
        }

        func StrValInt32(v int32) string {
        \treturn strconv.FormatInt(int64(v), 10)
        }

        func StrValInt64(v int64) string {
        \treturn strconv.FormatInt(v, 10)
        }

        func StrValBool(v bool) string {
        \tif v {
        \t\treturn "true"
        \t}
        \treturn "false"
        }

        func StrValChar(v uint16) string {
        \tif v < 0x80 {
        \t\treturn string(rune(v))
        \t}
        \treturn javalang.StringFromChar(v)
        }

        func StrVal(v any) string {
        \tif v == nil {
        \t\treturn "null"
        \t}
        \tswitch x := v.(type) {
        \tcase string:
        \t\tif x == NullString {
        \t\t\treturn "null"
        \t\t}
        \t\treturn x
        \tcase float64:
        \t\treturn javalang.DoubleToString(x)
        \tcase float32:
        \t\treturn javalang.FloatToString(x)
        \tcase int32:
        \t\treturn strconv.FormatInt(int64(x), 10)
        \tcase int64:
        \t\treturn strconv.FormatInt(x, 10)
        \tcase bool:
        \t\tif x {
        \t\t\treturn "true"
        \t\t}
        \t\treturn "false"
        \tcase uint16:
        \t\treturn javalang.StringFromChar(x)
        \tcase int8:
        \t\treturn strconv.FormatInt(int64(x), 10)
        \tcase int16:
        \t\treturn strconv.FormatInt(int64(x), 10)
        \tcase int:
        \t\treturn strconv.Itoa(x)
        \tcase *BoxDouble:
        \t\tif x == nil {
        \t\t\treturn "null"
        \t\t}
        \t\treturn javalang.DoubleToString(x.Val)
        \tcase *BoxFloat:
        \t\tif x == nil {
        \t\t\treturn "null"
        \t\t}
        \t\treturn javalang.FloatToString(x.Val)
        \tcase *BoxInt:
        \t\tif x == nil {
        \t\t\treturn "null"
        \t\t}
        \t\treturn strconv.FormatInt(int64(x.Val), 10)
        \tcase *BoxLong:
        \t\tif x == nil {
        \t\t\treturn "null"
        \t\t}
        \t\treturn strconv.FormatInt(x.Val, 10)
        \tcase *BoxBoolean:
        \t\tif x == nil {
        \t\t\treturn "null"
        \t\t}
        \t\tif x.Val {
        \t\t\treturn "true"
        \t\t}
        \t\treturn "false"
        \tcase *BoxChar:
        \t\tif x == nil {
        \t\t\treturn "null"
        \t\t}
        \t\treturn javalang.StringFromChar(x.Val)
        \tcase *BoxByte:
        \t\tif x == nil {
        \t\t\treturn "null"
        \t\t}
        \t\treturn strconv.FormatInt(int64(x.Val), 10)
        \tcase *BoxShort:
        \t\tif x == nil {
        \t\t\treturn "null"
        \t\t}
        \t\treturn strconv.FormatInt(int64(x.Val), 10)
        \t}
        \tif isJavaNilAny(v) {
        \t\treturn "null"
        \t}
        \treturn javalang.ToJavaString(v)
        }

        func StrConcat(parts ...any) string {
        \tn := len(parts)
        \tif n == 0 {
        \t\treturn ""
        \t}
        \tif n == 1 {
        \t\treturn StrVal(parts[0])
        \t}
        \tvar stackBuf [16]string
        \tvar strs []string
        \tif n <= len(stackBuf) {
        \t\tstrs = stackBuf[:n]
        \t} else {
        \t\tstrs = make([]string, n)
        \t}
        \ttotal := 0
        \tfor i, p := range parts {
        \t\ts := StrVal(p)
        \t\tstrs[i] = s
        \t\ttotal += len(s)
        \t}
        \tif total == 0 {
        \t\treturn ""
        \t}
        \tvar sb strings.Builder
        \tsb.Grow(total)
        \tfor _, s := range strs {
        \t\tsb.WriteString(s)
        \t}
        \treturn sb.String()
        }

        func NewStringArray(n int32) []string {
        \tif n < 0 {
        \t\tpanic(javalang.NewRuntimeExceptionMsg("NegativeArraySizeException"))
        \t}
        \ts := make([]string, n)
        \tfor i := range s {
        \t\ts[i] = NullString
        \t}
        \treturn s
        }

        func Make2D[T any](d1, d2 int32, initElem func() T) [][]T {
        \tout := make([][]T, d1)
        \tfor i := int32(0); i < d1; i++ {
        \t\trow := make([]T, d2)
        \t\tif initElem != nil {
        \t\t\tfor j := int32(0); j < d2; j++ {
        \t\t\t\trow[j] = initElem()
        \t\t\t}
        \t\t}
        \t\tout[i] = row
        \t}
        \treturn out
        }

        func CloneSlice[T any](s []T) []T {
        \tif s == nil {
        \t\treturn nil
        \t}
        \tout := make([]T, len(s))
        \tcopy(out, s)
        \treturn out
        }

        // --- Goroutine ID & Reentrant Monitor for Java `synchronized` ---

        func currentGoroutineID() uint64 {
        \tvar buf [64]byte
        \tn := runtime.Stack(buf[:], false)
        \t// Format: "goroutine <id> [..."
        \tconst prefix = "goroutine "
        \tif n <= len(prefix) {
        \t\treturn 1
        \t}
        \tvar id uint64
        \tfor i := len(prefix); i < n; i++ {
        \t\tb := buf[i]
        \t\tif b < '0' || b > '9' {
        \t\t\tbreak
        \t\t}
        \t\tid = id*10 + uint64(b-'0')
        \t}
        \treturn id
        }

        type reentrantMonitor struct {
        \tstateMu sync.Mutex
        \tlockCh  chan struct{}
        \towner   uint64
        \tdepth   int
        }

        func newReentrantMonitor() *reentrantMonitor {
        \tch := make(chan struct{}, 1)
        \tch <- struct{}{}
        \treturn &reentrantMonitor{lockCh: ch}
        }

        func (rm *reentrantMonitor) lock() {
        \tgid := currentGoroutineID()
        \trm.stateMu.Lock()
        \tif rm.depth > 0 && rm.owner == gid {
        \t\trm.depth++
        \t\trm.stateMu.Unlock()
        \t\treturn
        \t}
        \trm.stateMu.Unlock()
        \t<-rm.lockCh
        \trm.stateMu.Lock()
        \trm.owner = gid
        \trm.depth = 1
        \trm.stateMu.Unlock()
        }

        func (rm *reentrantMonitor) unlock() {
        \trm.stateMu.Lock()
        \trm.depth--
        \tif rm.depth == 0 {
        \t\trm.owner = 0
        \t\trm.stateMu.Unlock()
        \t\trm.lockCh <- struct{}{}
        \t\treturn
        \t}
        \trm.stateMu.Unlock()
        }

        var monitorMap sync.Map // key -> *reentrantMonitor

        func LockMonitor(obj any) func() {
        \tif obj == nil || javalang.IsNil(obj) {
        \t\tpanic(javalang.NewNullPointerException())
        \t}
        \trv := reflect.ValueOf(obj)
        \tvar key any = obj
        \tif rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Slice {
        \t\tkey = rv.Pointer()
        \t}
        \tm, ok := monitorMap.Load(key)
        \tif !ok {
        \t\tm, _ = monitorMap.LoadOrStore(key, newReentrantMonitor())
        \t}
        \trm := m.(*reentrantMonitor)
        \trm.lock()
        \treturn rm.unlock
        }

        func CloseResource(res any) {
        \tif res == nil || javalang.IsNil(res) {
        \t\treturn
        \t}
        \tif c, ok := res.(interface{ Close() }); ok {
        \t\tc.Close()
        \t} else if c2, ok := res.(interface{ Close_0() }); ok {
        \t\tc2.Close_0()
        \t} else if c3, ok := res.(interface{ Close() error }); ok {
        \t\t_ = c3.Close()
        \t}
        }

        type ClinitState struct {
        \tmu    sync.Mutex
        \towner uint64
        \twait  chan struct{}
        }

        func EnterClinit(st *ClinitState, done *uint32) bool {
        \tgid := currentGoroutineID()
        \tst.mu.Lock()
        \tfor {
        \t\tif atomic.LoadUint32(done) != 0 {
        \t\t\tst.mu.Unlock()
        \t\t\treturn false
        \t\t}
        \t\tif st.owner == gid {
        \t\t\tst.mu.Unlock()
        \t\t\treturn false
        \t\t}
        \t\tif st.owner == 0 {
        \t\t\tst.owner = gid
        \t\t\tst.wait = make(chan struct{})
        \t\t\tst.mu.Unlock()
        \t\t\treturn true
        \t\t}
        \t\tch := st.wait
        \t\tst.mu.Unlock()
        \t\t<-ch
        \t\tst.mu.Lock()
        \t}
        }

        func ExitClinit(st *ClinitState, done *uint32) {
        \tst.mu.Lock()
        \tatomic.StoreUint32(done, 1)
        \tch := st.wait
        \tst.owner = 0
        \tst.wait = nil
        \tst.mu.Unlock()
        \tif ch != nil {
        \t\tclose(ch)
        \t}
        }

        // --- Boxed Primitives & Number ---

        type Number interface {
        \tIntValue() int32
        \tLongValue() int64
        \tFloatValue() float32
        \tDoubleValue() float64
        \tByteValue() int8
        \tShortValue() int16
        \tImplements_j2gort_Number()
        }

        func ToIface_Number[T interface {
        \tNumber
        \tcomparable
        }](v T) Number {
        \tvar z T
        \tif v == z {
        \t\treturn nil
        \t}
        \treturn v
        }

        type NumberBase struct{}
        func (n *NumberBase) IsJavaNil_() bool   { return n == nil }
        func (*NumberBase) Implements_j2gort_Number() {}
        func (*NumberBase) IntValue() int32     { return 0 }
        func (*NumberBase) LongValue() int64    { return 0 }
        func (*NumberBase) FloatValue() float32 { return 0 }
        func (*NumberBase) DoubleValue() float64 { return 0 }
        func (*NumberBase) ByteValue() int8     { return 0 }
        func (*NumberBase) ShortValue() int16   { return 0 }
        func (*NumberBase) BindExtSelf_(any)    {}
        func InitExt_java_lang_Number(*NumberBase) {}

        type BoxBoolean struct{ Val bool }
        func (b *BoxBoolean) IsJavaNil_() bool   { return b == nil }
        func (b *BoxBoolean) BooleanValue() bool { return CheckNotNullPtr(b).Val }
        func (b *BoxBoolean) Equals(o any) bool {
        \tx, ok := o.(*BoxBoolean)
        \treturn ok && x != nil && x.Val == b.Val
        }
        func (b *BoxBoolean) HashCode() int32 {
        \tif b.Val { return 1231 }
        \treturn 1237
        }
        func (b *BoxBoolean) String() string {
        \tif b.Val { return "true" }
        \treturn "false"
        }
        func (*BoxBoolean) JavaClassInfo_() (string, string) { return "java.lang.Boolean", "Boolean" }

        var (
        \tBoolean_TRUE  = &BoxBoolean{Val: true}
        \tBoolean_FALSE = &BoxBoolean{Val: false}
        )

        func ValueOfBoolean(v bool) *BoxBoolean {
        \tif v { return Boolean_TRUE }
        \treturn Boolean_FALSE
        }
        func UnboxBoolean(b *BoxBoolean) bool { return CheckNotNullPtr(b).Val }

        type BoxByte struct{ Val int8 }
        func (b *BoxByte) IsJavaNil_() bool   { return b == nil }
        func (b *BoxByte) Implements_j2gort_Number() {}
        func (b *BoxByte) ByteValue() int8     { return CheckNotNullPtr(b).Val }
        func (b *BoxByte) ShortValue() int16   { return int16(CheckNotNullPtr(b).Val) }
        func (b *BoxByte) IntValue() int32     { return int32(CheckNotNullPtr(b).Val) }
        func (b *BoxByte) LongValue() int64    { return int64(CheckNotNullPtr(b).Val) }
        func (b *BoxByte) FloatValue() float32 { return float32(CheckNotNullPtr(b).Val) }
        func (b *BoxByte) DoubleValue() float64 { return float64(CheckNotNullPtr(b).Val) }
        func (b *BoxByte) Equals(o any) bool   { x, ok := o.(*BoxByte); return ok && x != nil && x.Val == b.Val }
        func (b *BoxByte) HashCode() int32     { return int32(b.Val) }
        func (b *BoxByte) String() string      { return strconv.Itoa(int(b.Val)) }
        func (*BoxByte) JavaClassInfo_() (string, string) { return "java.lang.Byte", "Byte" }
        var (
        \tboxByteCache_  [256]BoxByte
        \tboxShortCache_ [256]BoxShort
        \tboxCharCache_  [128]BoxChar
        \tboxIntCache_   [1152]BoxInt // [-128..1023]
        \tboxLongCache_  [256]BoxLong
        )

        func init() {
        \tfor i := 0; i < 256; i++ {
        \t\tboxByteCache_[i].Val = int8(i - 128)
        \t\tboxShortCache_[i].Val = int16(i - 128)
        \t\tboxLongCache_[i].Val = int64(i - 128)
        \t}
        \tfor i := 0; i < 128; i++ {
        \t\tboxCharCache_[i].Val = uint16(i)
        \t}
        \tfor i := 0; i < len(boxIntCache_); i++ {
        \t\tboxIntCache_[i].Val = int32(i - 128)
        \t}
        }

        func ValueOfByte(v int8) *BoxByte   { return &boxByteCache_[int(v)+128] }
        func UnboxByte(b *BoxByte) int8     { return CheckNotNullPtr(b).Val }

        type BoxShort struct{ Val int16 }
        func (b *BoxShort) IsJavaNil_() bool   { return b == nil }
        func (b *BoxShort) Implements_j2gort_Number() {}
        func (b *BoxShort) ByteValue() int8     { return int8(CheckNotNullPtr(b).Val) }
        func (b *BoxShort) ShortValue() int16   { return CheckNotNullPtr(b).Val }
        func (b *BoxShort) IntValue() int32     { return int32(CheckNotNullPtr(b).Val) }
        func (b *BoxShort) LongValue() int64    { return int64(CheckNotNullPtr(b).Val) }
        func (b *BoxShort) FloatValue() float32 { return float32(CheckNotNullPtr(b).Val) }
        func (b *BoxShort) DoubleValue() float64 { return float64(CheckNotNullPtr(b).Val) }
        func (b *BoxShort) Equals(o any) bool   { x, ok := o.(*BoxShort); return ok && x != nil && x.Val == b.Val }
        func (b *BoxShort) HashCode() int32     { return int32(b.Val) }
        func (b *BoxShort) String() string      { return strconv.Itoa(int(b.Val)) }
        func (*BoxShort) JavaClassInfo_() (string, string) { return "java.lang.Short", "Short" }
        func ValueOfShort(v int16) *BoxShort {
        \tif v >= -128 && v <= 127 {
        \t\treturn &boxShortCache_[int(v)+128]
        \t}
        \treturn &BoxShort{Val: v}
        }
        func UnboxShort(b *BoxShort) int16 { return CheckNotNullPtr(b).Val }

        type BoxChar struct{ Val uint16 }
        func (b *BoxChar) IsJavaNil_() bool { return b == nil }
        func (b *BoxChar) CharValue() uint16 { return CheckNotNullPtr(b).Val }
        func (b *BoxChar) Equals(o any) bool { x, ok := o.(*BoxChar); return ok && x != nil && x.Val == b.Val }
        func (b *BoxChar) HashCode() int32   { return int32(b.Val) }
        func (b *BoxChar) String() string    { return javalang.StringFromChar(b.Val) }
        func (*BoxChar) JavaClassInfo_() (string, string) { return "java.lang.Character", "Character" }
        func ValueOfChar(v uint16) *BoxChar {
        \tif v < 128 {
        \t\treturn &boxCharCache_[v]
        \t}
        \treturn &BoxChar{Val: v}
        }
        func UnboxChar(b *BoxChar) uint16 { return CheckNotNullPtr(b).Val }

        type BoxInt struct{ Val int32 }
        func (b *BoxInt) IsJavaNil_() bool   { return b == nil }
        func (b *BoxInt) Implements_j2gort_Number() {}
        func (b *BoxInt) ByteValue() int8     { return int8(CheckNotNullPtr(b).Val) }
        func (b *BoxInt) ShortValue() int16   { return int16(CheckNotNullPtr(b).Val) }
        func (b *BoxInt) IntValue() int32     { return CheckNotNullPtr(b).Val }
        func (b *BoxInt) LongValue() int64    { return int64(CheckNotNullPtr(b).Val) }
        func (b *BoxInt) FloatValue() float32 { return float32(CheckNotNullPtr(b).Val) }
        func (b *BoxInt) DoubleValue() float64 { return float64(CheckNotNullPtr(b).Val) }
        func (b *BoxInt) Equals(o any) bool   { x, ok := o.(*BoxInt); return ok && x != nil && x.Val == b.Val }
        func (b *BoxInt) HashCode() int32     { return b.Val }
        func (b *BoxInt) String() string      { return strconv.FormatInt(int64(b.Val), 10) }
        func (*BoxInt) JavaClassInfo_() (string, string) { return "java.lang.Integer", "Integer" }
        func ValueOfInt(v int32) *BoxInt {
        \tif v >= -128 && v < 1024 {
        \t\treturn &boxIntCache_[v+128]
        \t}
        \treturn &BoxInt{Val: v}
        }
        func UnboxInt(b *BoxInt) int32 { return CheckNotNullPtr(b).Val }

        type BoxLong struct{ Val int64 }
        func (b *BoxLong) IsJavaNil_() bool   { return b == nil }
        func (b *BoxLong) Implements_j2gort_Number() {}
        func (b *BoxLong) ByteValue() int8     { return int8(CheckNotNullPtr(b).Val) }
        func (b *BoxLong) ShortValue() int16   { return int16(CheckNotNullPtr(b).Val) }
        func (b *BoxLong) IntValue() int32     { return int32(CheckNotNullPtr(b).Val) }
        func (b *BoxLong) LongValue() int64    { return CheckNotNullPtr(b).Val }
        func (b *BoxLong) FloatValue() float32 { return float32(CheckNotNullPtr(b).Val) }
        func (b *BoxLong) DoubleValue() float64 { return float64(CheckNotNullPtr(b).Val) }
        func (b *BoxLong) Equals(o any) bool   { x, ok := o.(*BoxLong); return ok && x != nil && x.Val == b.Val }
        func (b *BoxLong) HashCode() int32     { return int32(uint64(b.Val) ^ (uint64(b.Val) >> 32)) }
        func (b *BoxLong) String() string      { return strconv.FormatInt(b.Val, 10) }
        func (*BoxLong) JavaClassInfo_() (string, string) { return "java.lang.Long", "Long" }
        func ValueOfLong(v int64) *BoxLong {
        \tif v >= -128 && v <= 127 {
        \t\treturn &boxLongCache_[int(v)+128]
        \t}
        \treturn &BoxLong{Val: v}
        }
        func UnboxLong(b *BoxLong) int64       { return CheckNotNullPtr(b).Val }

        type BoxFloat struct{ Val float32 }
        func (b *BoxFloat) IsJavaNil_() bool   { return b == nil }
        func (b *BoxFloat) Implements_j2gort_Number() {}
        func (b *BoxFloat) ByteValue() int8     { return int8(javalang.D2I(float64(CheckNotNullPtr(b).Val))) }
        func (b *BoxFloat) ShortValue() int16   { return int16(javalang.D2I(float64(CheckNotNullPtr(b).Val))) }
        func (b *BoxFloat) IntValue() int32     { return javalang.D2I(float64(CheckNotNullPtr(b).Val)) }
        func (b *BoxFloat) LongValue() int64    { return javalang.D2L(float64(CheckNotNullPtr(b).Val)) }
        func (b *BoxFloat) FloatValue() float32 { return CheckNotNullPtr(b).Val }
        func (b *BoxFloat) DoubleValue() float64 { return float64(CheckNotNullPtr(b).Val) }
        func (b *BoxFloat) Equals(o any) bool {
        \tx, ok := o.(*BoxFloat)
        \tif !ok || x == nil { return false }
        \tif math.IsNaN(float64(b.Val)) && math.IsNaN(float64(x.Val)) { return true }
        \treturn math.Float32bits(b.Val) == math.Float32bits(x.Val)
        }
        func (b *BoxFloat) HashCode() int32 {
        \tif math.IsNaN(float64(b.Val)) { return 0x7fc00000 }
        \treturn int32(math.Float32bits(b.Val))
        }
        func (b *BoxFloat) String() string    { return javalang.FloatToString(b.Val) }
        func (b *BoxFloat) IsNaN() bool       { return math.IsNaN(float64(CheckNotNullPtr(b).Val)) }
        func (b *BoxFloat) IsInfinite() bool  { return math.IsInf(float64(CheckNotNullPtr(b).Val), 0) }
        func (*BoxFloat) JavaClassInfo_() (string, string) { return "java.lang.Float", "Float" }
        func ValueOfFloat(v float32) *BoxFloat { return &BoxFloat{Val: v} }
        func UnboxFloat(b *BoxFloat) float32   { return CheckNotNullPtr(b).Val }

        type BoxDouble struct{ Val float64 }
        func (b *BoxDouble) IsJavaNil_() bool   { return b == nil }
        func (b *BoxDouble) Implements_j2gort_Number() {}
        func (b *BoxDouble) ByteValue() int8     { return int8(javalang.D2I(CheckNotNullPtr(b).Val)) }
        func (b *BoxDouble) ShortValue() int16   { return int16(javalang.D2I(CheckNotNullPtr(b).Val)) }
        func (b *BoxDouble) IntValue() int32     { return javalang.D2I(CheckNotNullPtr(b).Val) }
        func (b *BoxDouble) LongValue() int64    { return javalang.D2L(CheckNotNullPtr(b).Val) }
        func (b *BoxDouble) FloatValue() float32 { return float32(CheckNotNullPtr(b).Val) }
        func (b *BoxDouble) DoubleValue() float64 { return CheckNotNullPtr(b).Val }
        func (b *BoxDouble) Equals(o any) bool {
        \tx, ok := o.(*BoxDouble)
        \tif !ok || x == nil { return false }
        \tif math.IsNaN(b.Val) && math.IsNaN(x.Val) { return true }
        \treturn math.Float64bits(b.Val) == math.Float64bits(x.Val)
        }
        func (b *BoxDouble) HashCode() int32 {
        \tif math.IsNaN(b.Val) { return 0x7ff80000 }
        \tu := math.Float64bits(b.Val)
        \treturn int32(u ^ (u >> 32))
        }
        func (b *BoxDouble) String() string     { return javalang.DoubleToString(b.Val) }
        func (b *BoxDouble) IsNaN() bool        { return math.IsNaN(CheckNotNullPtr(b).Val) }
        func (b *BoxDouble) IsInfinite() bool   { return math.IsInf(CheckNotNullPtr(b).Val, 0) }
        func (*BoxDouble) JavaClassInfo_() (string, string) { return "java.lang.Double", "Double" }
        func ValueOfDouble(v float64) *BoxDouble { return &BoxDouble{Val: v} }
        func UnboxDouble(b *BoxDouble) float64   { return CheckNotNullPtr(b).Val }

        func UnboxToDouble(v any) float64 {
        \tif n, ok := v.(Number); ok && !isJavaNilAny(n) {
        \t\treturn n.DoubleValue()
        \t}
        \tif c, ok := v.(*BoxChar); ok && c != nil {
        \t\treturn float64(c.Val)
        \t}
        \tpanic(javalang.NewNullPointerException())
        }

        // --- Enum & Type Predicate Registries ---

        type EnumBase struct {
        \tName_     string                         `gson:"-"`
        \tOrdinal_  int32                          `gson:"-"`
        \tequalPred atomic.Pointer[predEqualEnum_] `gson:"-"`
        }

        func (e *EnumBase) IsJavaNil_() bool { return e == nil }
        func (e *EnumBase) Name() string   { return e.Name_ }
        func (e *EnumBase) Ordinal() int32 { return e.Ordinal_ }
        func (e *EnumBase) String() string { return e.Name_ }
        func (e *EnumBase) HashCode() int32 {
        \treturn javalang.StringHashCode(e.Name_)*31 + e.Ordinal_
        }
        func (e *EnumBase) Equals(other any) bool {
        \tif isJavaNilAny(other) {
        \t\treturn false
        \t}
        \tif oa, ok := other.(EnumAccessor); ok {
        \t\treturn e == oa.GetEnumBase()
        \t}
        \treturn false
        }
        func (e *EnumBase) CompareTo(other any) int32 {
        \tif isJavaNilAny(other) {
        \t\tpanic(javalang.NewNullPointerException())
        \t}
        \toa := other.(EnumAccessor).GetEnumBase()
        \tif e.Ordinal_ < oa.Ordinal_ {
        \t\treturn -1
        \t}
        \tif e.Ordinal_ > oa.Ordinal_ {
        \t\treturn 1
        \t}
        \treturn 0
        }

        type EnumAccessor interface {
        \tGetEnumBase() *EnumBase
        }

        func (e *EnumBase) GetEnumBase() *EnumBase { return e }
        func (e *EnumBase) EnumOrdinal_() int32 {
        \tif e == nil {
        \t\treturn -1
        \t}
        \treturn e.Ordinal_
        }
        func (*EnumBase) HasIdentityEquals_()                                 {}
        func (*EnumBase) Implements_j2gort_Ext_java_lang_Enum()               {}
        func (*EnumBase) Implements_j2gort_Ext_java_lang_constant_Constable() {}
        func (e *EnumBase) GetExt_java_lang_EnumBase() *Ext_java_lang_EnumBase {
        \treturn &Ext_java_lang_EnumBase{Core_: e}
        }

        func EnumName(v any) string {
        \tif isJavaNilAny(v) { panic(javalang.NewNullPointerException()) }
        \treturn v.(EnumAccessor).GetEnumBase().Name_
        }

        func EnumOrdinal(v any) int32 {
        \tif isJavaNilAny(v) { panic(javalang.NewNullPointerException()) }
        \treturn v.(EnumAccessor).GetEnumBase().Ordinal_
        }

        func EnumCompareTo(a, b any) int32 {
        \treturn EnumOrdinal(a) - EnumOrdinal(b)
        }

        var (
        \tenumRegistry   sync.Map // string -> []any
        \tenumLoaderMap  sync.Map // string -> func()
        \ttypePredMap    sync.Map // string -> func(any) bool
        )

        func RegisterEnumConstants(className string, vals []any) {
        \tenumRegistry.Store(className, vals)
        }

        func RegisterEnumConstantsLoader(qn, binName string, loader func()) {
        \tenumLoaderMap.Store(qn, loader)
        \tif binName != qn {
        \t\tenumLoaderMap.Store(binName, loader)
        \t}
        }

        func LookupEnumConstantsByName(name string) []any {
        \tif name == "" {
        \t\treturn nil
        \t}
        \tif v, ok := enumRegistry.Load(name); ok {
        \t\treturn v.([]any)
        \t}
        \tif loader, ok := enumLoaderMap.Load(name); ok {
        \t\tloader.(func())()
        \t}
        \tif v, ok := enumRegistry.Load(name); ok {
        \t\treturn v.([]any)
        \t}
        \treturn nil
        }

        func LookupEnumConstants(cls *javalang.Class) []any {
        \tif cls == nil {
        \t\treturn nil
        \t}
        \treturn LookupEnumConstantsByName(cls.GetName())
        }

        func RegisterTypePredicate(qn, binName string, pred func(any) bool) {
        \ttypePredMap.Store(qn, pred)
        \tif binName != qn {
        \t\ttypePredMap.Store(binName, pred)
        \t}
        }

        func CheckTypePredicate(cls *javalang.Class, v any) bool {
        \tif cls == nil || v == nil || javalang.IsNil(v) {
        \t\treturn false
        \t}
        \tif pred, ok := typePredMap.Load(cls.GetName()); ok {
        \t\treturn pred.(func(any) bool)(v)
        \t}
        \treturn GetObjectClass(v).GetName() == cls.GetName()
        }
        """;
  }

  private String buildCoreRuntimeGoPart1Predicates() {
    return """

        // --- Zero-Closure Guava & JDK Concrete Predicates ---

        type jEqualTargetProvider interface {
        \tequalTarget_() any
        }

        type predAlwaysTrueType_ struct{ _ byte }

        func (*predAlwaysTrueType_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (*predAlwaysTrueType_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
        func (*predAlwaysTrueType_) HasIdentityEquals_()                                     {}
        func (p *predAlwaysTrueType_) IsJavaNil_() bool                                      { return p == nil }
        func (p *predAlwaysTrueType_) Apply(_ any) bool                                      { CheckNotNullPtr(p); return true }
        func (p *predAlwaysTrueType_) Test(_ any) bool                                       { CheckNotNullPtr(p); return true }
        func (p *predAlwaysTrueType_) TestNode_(_ uint32) bool                               { CheckNotNullPtr(p); return true }
        func (p *predAlwaysTrueType_) Equals(o any) bool                                     { CheckNotNullPtr(p); return p == o }
        func (p *predAlwaysTrueType_) HashCode() int32                                       { CheckNotNullPtr(p); return 0x54525545 }
        func (p *predAlwaysTrueType_) String() string                                        { CheckNotNullPtr(p); return "Predicates.alwaysTrue()" }
        func (p *predAlwaysTrueType_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn CheckNotNull(o)
        }
        func (p *predAlwaysTrueType_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \tCheckNotNull(o)
        \treturn p
        }
        func (p *predAlwaysTrueType_) Negate() Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn predAlwaysFalseIface_
        }

        type predAlwaysFalseType_ struct{ _ byte }

        func (*predAlwaysFalseType_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (*predAlwaysFalseType_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
        func (*predAlwaysFalseType_) HasIdentityEquals_()                                     {}
        func (p *predAlwaysFalseType_) IsJavaNil_() bool                                      { return p == nil }
        func (p *predAlwaysFalseType_) Apply(_ any) bool                                      { CheckNotNullPtr(p); return false }
        func (p *predAlwaysFalseType_) Test(_ any) bool                                       { CheckNotNullPtr(p); return false }
        func (p *predAlwaysFalseType_) TestNode_(_ uint32) bool                               { CheckNotNullPtr(p); return false }
        func (p *predAlwaysFalseType_) Equals(o any) bool                                     { CheckNotNullPtr(p); return p == o }
        func (p *predAlwaysFalseType_) HashCode() int32                                       { CheckNotNullPtr(p); return 0x46414c53 }
        func (p *predAlwaysFalseType_) String() string                                        { CheckNotNullPtr(p); return "Predicates.alwaysFalse()" }
        func (p *predAlwaysFalseType_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \tCheckNotNull(o)
        \treturn p
        }
        func (p *predAlwaysFalseType_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn CheckNotNull(o)
        }
        func (p *predAlwaysFalseType_) Negate() Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn predAlwaysTrueIface_
        }

        type predIsNullType_ struct{ _ byte }

        func (*predIsNullType_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (*predIsNullType_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
        func (*predIsNullType_) HasIdentityEquals_()                                     {}
        func (p *predIsNullType_) IsJavaNil_() bool                                      { return p == nil }
        func (p *predIsNullType_) Apply(x any) bool                                      { CheckNotNullPtr(p); return x == nil || isJavaNilAny(x) }
        func (p *predIsNullType_) Test(x any) bool                                       { CheckNotNullPtr(p); return x == nil || isJavaNilAny(x) }
        func (p *predIsNullType_) TestNode_(n uint32) bool                               { CheckNotNullPtr(p); return n == 0 }
        func (p *predIsNullType_) Equals(o any) bool                                     { CheckNotNullPtr(p); return p == o }
        func (p *predIsNullType_) HashCode() int32                                       { CheckNotNullPtr(p); return 0 }
        func (p *predIsNullType_) String() string                                        { CheckNotNullPtr(p); return "Predicates.isNull()" }
        func (p *predIsNullType_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predAnd2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predIsNullType_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predOr2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predIsNullType_) Negate() Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn predNotNullIface_
        }

        type predNotNullType_ struct{ _ byte }

        func (*predNotNullType_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (*predNotNullType_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
        func (*predNotNullType_) HasIdentityEquals_()                                     {}
        func (p *predNotNullType_) IsJavaNil_() bool                                      { return p == nil }
        func (p *predNotNullType_) Apply(x any) bool                                      { CheckNotNullPtr(p); return x != nil && !isJavaNilAny(x) }
        func (p *predNotNullType_) Test(x any) bool                                       { CheckNotNullPtr(p); return x != nil && !isJavaNilAny(x) }
        func (p *predNotNullType_) TestNode_(n uint32) bool                               { CheckNotNullPtr(p); return n != 0 }
        func (p *predNotNullType_) Equals(o any) bool                                     { CheckNotNullPtr(p); return p == o }
        func (p *predNotNullType_) HashCode() int32                                       { CheckNotNullPtr(p); return 0x4e4e554c }
        func (p *predNotNullType_) String() string                                        { CheckNotNullPtr(p); return "Predicates.notNull()" }
        func (p *predNotNullType_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predAnd2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predNotNullType_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predOr2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predNotNullType_) Negate() Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn predIsNullIface_
        }

        var (
        \tpredAlwaysTrueSingleton_  predAlwaysTrueType_
        \tpredAlwaysFalseSingleton_ predAlwaysFalseType_
        \tpredIsNullSingleton_      predIsNullType_
        \tpredNotNullSingleton_     predNotNullType_

        \tpredAlwaysTrueIface_  Ext_com_google_common_base_Predicate = &predAlwaysTrueSingleton_
        \tpredAlwaysFalseIface_ Ext_com_google_common_base_Predicate = &predAlwaysFalseSingleton_
        \tpredIsNullIface_      Ext_com_google_common_base_Predicate = &predIsNullSingleton_
        \tpredNotNullIface_     Ext_com_google_common_base_Predicate = &predNotNullSingleton_

        \tpredAlwaysTrueAny_ any = &predAlwaysTrueSingleton_
        )

        type predEqualNode_ uint32

        func (predEqualNode_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (predEqualNode_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
        func (predEqualNode_) HasIdentityEquals_()                                     {}
        func (predEqualNode_) IsJavaNil_() bool                                        { return false }
        func (pn predEqualNode_) equalTarget_() any                                    { return boxNodeById(uint32(pn)) }
        func (pn predEqualNode_) Apply(x any) bool {
        \tif x == nil {
        \t\treturn false
        \t}
        \tif rn, ok := x.(jRawNodeProvider); ok {
        \t\treturn rn.RawNodeId_() == uint32(pn)
        \t}
        \tif n, ok := x.(nodestore.Node); ok {
        \t\treturn uint32(n) == uint32(pn)
        \t}
        \treturn false
        }
        func (pn predEqualNode_) Test(x any) bool         { return pn.Apply(x) }
        func (pn predEqualNode_) TestNode_(n uint32) bool { return uint32(pn) == n && n != 0 }
        func (pn predEqualNode_) Equals(o any) bool {
        \tif op, ok := o.(predEqualNode_); ok {
        \t\treturn op == pn
        \t}
        \tif op, ok := o.(jEqualTargetProvider); ok && !isJavaNilAny(o) {
        \t\treturn pn.Apply(op.equalTarget_())
        \t}
        \treturn false
        }
        func (pn predEqualNode_) HashCode() int32 { return int32(pn) }
        func (pn predEqualNode_) String() string {
        \treturn "Predicates.equalTo(" + javalang.ToJavaString(boxNodeById(uint32(pn))) + ")"
        }
        func (pn predEqualNode_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \treturn &predAnd2_{p1: pn, p2: CheckNotNull(o)}
        }
        func (pn predEqualNode_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \treturn &predOr2_{p1: pn, p2: CheckNotNull(o)}
        }
        func (pn predEqualNode_) Negate() Ext_java_util_function_Predicate {
        \treturn &predNot_{inner: pn}
        }

        type predEqualStr_ string

        var predEqualEmptyStrIface_ Ext_com_google_common_base_Predicate = predEqualStr_("")

        func (predEqualStr_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (predEqualStr_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
        func (predEqualStr_) HasIdentityEquals_()                                     {}
        func (predEqualStr_) IsJavaNil_() bool                                        { return false }
        func (ps predEqualStr_) equalTarget_() any                                    { return string(ps) }
        func (ps predEqualStr_) Apply(x any) bool {
        \ts, ok := x.(string)
        \treturn ok && s == string(ps)
        }
        func (ps predEqualStr_) Test(x any) bool {
        \ts, ok := x.(string)
        \treturn ok && s == string(ps)
        }
        func (ps predEqualStr_) TestNode_(_ uint32) bool { return false }
        func (ps predEqualStr_) Equals(o any) bool {
        \tif op, ok := o.(predEqualStr_); ok {
        \t\treturn op == ps
        \t}
        \tif op, ok := o.(jEqualTargetProvider); ok && !isJavaNilAny(o) {
        \t\treturn ps.Apply(op.equalTarget_())
        \t}
        \treturn false
        }
        func (ps predEqualStr_) HashCode() int32 { return javalang.StringHashCode(string(ps)) }
        func (ps predEqualStr_) String() string  { return "Predicates.equalTo(" + string(ps) + ")" }
        func (ps predEqualStr_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \treturn &predAnd2_{p1: ps, p2: CheckNotNull(o)}
        }
        func (ps predEqualStr_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \treturn &predOr2_{p1: ps, p2: CheckNotNull(o)}
        }
        func (ps predEqualStr_) Negate() Ext_java_util_function_Predicate {
        \treturn &predNot_{inner: ps}
        }

        type predEqualEnum_ struct {
        \teb *EnumBase
        }

        func (*predEqualEnum_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (*predEqualEnum_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
        func (*predEqualEnum_) HasIdentityEquals_()                                     {}
        func (p *predEqualEnum_) IsJavaNil_() bool                                      { return p == nil }
        func (p *predEqualEnum_) equalTarget_() any                                     { return p.eb }
        func (p *predEqualEnum_) Apply(x any) bool {
        \tCheckNotNullPtr(p)
        \tif x == nil {
        \t\treturn false
        \t}
        \tif oa, ok := x.(EnumAccessor); ok {
        \t\tif nc, ok2 := x.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
        \t\t\treturn false
        \t\t}
        \t\treturn oa.GetEnumBase() == p.eb
        \t}
        \treturn false
        }
        func (p *predEqualEnum_) Test(x any) bool         { return p.Apply(x) }
        func (p *predEqualEnum_) TestNode_(_ uint32) bool { CheckNotNullPtr(p); return false }
        func (p *predEqualEnum_) Equals(o any) bool {
        \tCheckNotNullPtr(p)
        \tif op, ok := o.(*predEqualEnum_); ok {
        \t\treturn op != nil && op.eb == p.eb
        \t}
        \tif op, ok := o.(jEqualTargetProvider); ok && !isJavaNilAny(o) {
        \t\treturn p.Apply(op.equalTarget_())
        \t}
        \treturn false
        }
        func (p *predEqualEnum_) HashCode() int32 { CheckNotNullPtr(p); return p.eb.HashCode() }
        func (p *predEqualEnum_) String() string  { CheckNotNullPtr(p); return "Predicates.equalTo(" + p.eb.Name_ + ")" }
        func (p *predEqualEnum_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predAnd2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predEqualEnum_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predOr2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predEqualEnum_) Negate() Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predNot_{inner: p}
        }

        type predEqualIdentity_ struct {
        \ttarget any
        }

        func (*predEqualIdentity_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (*predEqualIdentity_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
        func (p *predEqualIdentity_) IsJavaNil_() bool                                      { return p == nil }
        func (p *predEqualIdentity_) equalTarget_() any                                     { return p.target }
        func (p *predEqualIdentity_) Apply(x any) bool                                      { CheckNotNullPtr(p); return x == p.target }
        func (p *predEqualIdentity_) Test(x any) bool                                       { CheckNotNullPtr(p); return x == p.target }
        func (p *predEqualIdentity_) TestNode_(_ uint32) bool                               { CheckNotNullPtr(p); return false }
        func (p *predEqualIdentity_) Equals(o any) bool {
        \tCheckNotNullPtr(p)
        \tif op, ok := o.(jEqualTargetProvider); ok && !isJavaNilAny(o) {
        \t\treturn op.equalTarget_() == p.target
        \t}
        \treturn false
        }
        func (p *predEqualIdentity_) HashCode() int32 { CheckNotNullPtr(p); return javalang.ObjectHashCode(p.target) }
        func (p *predEqualIdentity_) String() string {
        \tCheckNotNullPtr(p)
        \treturn "Predicates.equalTo(" + javalang.ToJavaString(p.target) + ")"
        }
        func (p *predEqualIdentity_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predAnd2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predEqualIdentity_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predOr2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predEqualIdentity_) Negate() Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predNot_{inner: p}
        }

        type predEqualAny_ struct {
        \ttarget any
        }

        func (*predEqualAny_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (*predEqualAny_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
        func (p *predEqualAny_) IsJavaNil_() bool                                      { return p == nil }
        func (p *predEqualAny_) equalTarget_() any                                     { return p.target }
        func (p *predEqualAny_) Apply(x any) bool {
        \tCheckNotNullPtr(p)
        \tif x == nil {
        \t\treturn false
        \t}
        \treturn javalang.ObjectsEquals(p.target, ToAny(x))
        }
        func (p *predEqualAny_) Test(x any) bool { return p.Apply(x) }
        func (p *predEqualAny_) TestNode_(n uint32) bool {
        \tCheckNotNullPtr(p)
        \tif n == 0 {
        \t\treturn false
        \t}
        \treturn javalang.ObjectsEquals(p.target, boxNodeById(n))
        }
        func (p *predEqualAny_) Equals(o any) bool {
        \tCheckNotNullPtr(p)
        \tif op, ok := o.(jEqualTargetProvider); ok && !isJavaNilAny(o) {
        \t\treturn javalang.ObjectsEquals(p.target, op.equalTarget_())
        \t}
        \treturn false
        }
        func (p *predEqualAny_) HashCode() int32 { CheckNotNullPtr(p); return javalang.ObjectHashCode(p.target) }
        func (p *predEqualAny_) String() string {
        \tCheckNotNullPtr(p)
        \treturn "Predicates.equalTo(" + javalang.ToJavaString(p.target) + ")"
        }
        func (p *predEqualAny_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predAnd2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predEqualAny_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predOr2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predEqualAny_) Negate() Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predNot_{inner: p}
        }

        var (
        \tpredEqualBoolTrue_  Ext_com_google_common_base_Predicate = &predEqualAny_{target: Boolean_TRUE}
        \tpredEqualBoolFalse_ Ext_com_google_common_base_Predicate = &predEqualAny_{target: Boolean_FALSE}
        )

        func guavaPredEqualTo(raw any) Ext_com_google_common_base_Predicate {
        \tif raw == nil || isJavaNilAny(raw) {
        \t\treturn predIsNullIface_
        \t}
        \tif ea, ok := raw.(EnumAccessor); ok {
        \t\tif eb := ea.GetEnumBase(); eb != nil {
        \t\t\tif p := eb.equalPred.Load(); p != nil {
        \t\t\t\treturn p
        \t\t\t}
        \t\t\tnp := &predEqualEnum_{eb: eb}
        \t\t\tif eb.equalPred.CompareAndSwap(nil, np) {
        \t\t\t\treturn np
        \t\t\t}
        \t\t\treturn eb.equalPred.Load()
        \t\t}
        \t}
        \tif rn, ok := raw.(jRawNodeProvider); ok {
        \t\tid := rn.RawNodeId_()
        \t\tif id == 0 {
        \t\t\treturn predIsNullIface_
        \t\t}
        \t\treturn predEqualNode_(id)
        \t}
        \tswitch v := raw.(type) {
        \tcase nodestore.Node:
        \t\tif v == 0 {
        \t\t\treturn predIsNullIface_
        \t\t}
        \t\treturn predEqualNode_(v)
        \tcase string:
        \t\tif v == "" {
        \t\t\treturn predEqualEmptyStrIface_
        \t\t}
        \t\treturn predEqualStr_(v)
        \tcase *BoxBoolean:
        \t\tif v.Val {
        \t\t\treturn predEqualBoolTrue_
        \t\t}
        \t\treturn predEqualBoolFalse_
        \t}
        \tif _, ok := raw.(IdentityComparable); ok {
        \t\treturn &predEqualIdentity_{target: raw}
        \t}
        \treturn &predEqualAny_{target: ToAny(raw)}
        }

        type predNot_ struct {
        \tinner Ext_java_util_function_Predicate
        }

        func (*predNot_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (*predNot_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
        func (p *predNot_) IsJavaNil_() bool                                      { return p == nil }
        func (p *predNot_) Apply(x any) bool                                      { CheckNotNullPtr(p); return !p.inner.Test(x) }
        func (p *predNot_) Test(x any) bool                                       { CheckNotNullPtr(p); return !p.inner.Test(x) }
        func (p *predNot_) TestNode_(n uint32) bool                               { CheckNotNullPtr(p); return !PredTestNode(p.inner, n) }
        func (p *predNot_) Equals(o any) bool {
        \tCheckNotNullPtr(p)
        \top, ok := o.(*predNot_)
        \treturn ok && op != nil && javalang.ObjectsEquals(p.inner, op.inner)
        }
        func (p *predNot_) HashCode() int32 { CheckNotNullPtr(p); return ^javalang.ObjectHashCode(p.inner) }
        func (p *predNot_) String() string  { CheckNotNullPtr(p); return "Predicates.not(" + javalang.ToJavaString(p.inner) + ")" }
        func (p *predNot_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predAnd2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predNot_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predOr2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predNot_) Negate() Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn p.inner
        }

        func guavaPredNot(target Ext_com_google_common_base_Predicate) Ext_com_google_common_base_Predicate {
        \tCheckNotNull(target)
        \treturn &predNot_{inner: target}
        }

        func jdkPredNot(target Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNull(target)
        \treturn &predNot_{inner: target}
        }

        type predAnd2_ struct {
        \tp1 Ext_java_util_function_Predicate
        \tp2 Ext_java_util_function_Predicate
        }

        func (*predAnd2_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (*predAnd2_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
        func (p *predAnd2_) IsJavaNil_() bool                                      { return p == nil }
        func (p *predAnd2_) Apply(x any) bool                                      { CheckNotNullPtr(p); return p.p1.Test(x) && p.p2.Test(x) }
        func (p *predAnd2_) Test(x any) bool                                       { CheckNotNullPtr(p); return p.p1.Test(x) && p.p2.Test(x) }
        func (p *predAnd2_) TestNode_(n uint32) bool                               { CheckNotNullPtr(p); return PredTestNode(p.p1, n) && PredTestNode(p.p2, n) }
        func (p *predAnd2_) Equals(o any) bool {
        \tCheckNotNullPtr(p)
        \top, ok := o.(*predAnd2_)
        \treturn ok && op != nil && javalang.ObjectsEquals(p.p1, op.p1) && javalang.ObjectsEquals(p.p2, op.p2)
        }
        func (p *predAnd2_) HashCode() int32 {
        \tCheckNotNullPtr(p)
        \treturn javalang.ObjectHashCode(p.p1)*31 + javalang.ObjectHashCode(p.p2) + 0x12472c2c
        }
        func (p *predAnd2_) String() string {
        \tCheckNotNullPtr(p)
        \treturn "Predicates.and(" + javalang.ToJavaString(p.p1) + "," + javalang.ToJavaString(p.p2) + ")"
        }
        func (p *predAnd2_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predAnd2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predAnd2_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predOr2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predAnd2_) Negate() Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predNot_{inner: p}
        }

        func guavaPredAnd2(a, b Ext_com_google_common_base_Predicate) Ext_com_google_common_base_Predicate {
        \tCheckNotNull(a)
        \tCheckNotNull(b)
        \treturn &predAnd2_{p1: a, p2: b}
        }

        type predOr2_ struct {
        \tp1 Ext_java_util_function_Predicate
        \tp2 Ext_java_util_function_Predicate
        }

        func (*predOr2_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (*predOr2_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
        func (p *predOr2_) IsJavaNil_() bool                                      { return p == nil }
        func (p *predOr2_) Apply(x any) bool                                      { CheckNotNullPtr(p); return p.p1.Test(x) || p.p2.Test(x) }
        func (p *predOr2_) Test(x any) bool                                       { CheckNotNullPtr(p); return p.p1.Test(x) || p.p2.Test(x) }
        func (p *predOr2_) TestNode_(n uint32) bool                               { CheckNotNullPtr(p); return PredTestNode(p.p1, n) || PredTestNode(p.p2, n) }
        func (p *predOr2_) Equals(o any) bool {
        \tCheckNotNullPtr(p)
        \top, ok := o.(*predOr2_)
        \treturn ok && op != nil && javalang.ObjectsEquals(p.p1, op.p1) && javalang.ObjectsEquals(p.p2, op.p2)
        }
        func (p *predOr2_) HashCode() int32 {
        \tCheckNotNullPtr(p)
        \treturn javalang.ObjectHashCode(p.p1)*31 + javalang.ObjectHashCode(p.p2) + 0x053c91cf
        }
        func (p *predOr2_) String() string {
        \tCheckNotNullPtr(p)
        \treturn "Predicates.or(" + javalang.ToJavaString(p.p1) + "," + javalang.ToJavaString(p.p2) + ")"
        }
        func (p *predOr2_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predAnd2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predOr2_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predOr2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predOr2_) Negate() Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predNot_{inner: p}
        }

        func guavaPredOr2(a, b Ext_com_google_common_base_Predicate) Ext_com_google_common_base_Predicate {
        \tCheckNotNull(a)
        \tCheckNotNull(b)
        \treturn &predOr2_{p1: a, p2: b}
        }

        type predAndSlice_ struct {
        \tpreds []Ext_java_util_function_Predicate
        }

        func (*predAndSlice_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (*predAndSlice_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
        func (p *predAndSlice_) IsJavaNil_() bool                                      { return p == nil }
        func (p *predAndSlice_) Apply(x any) bool {
        \tCheckNotNullPtr(p)
        \tfor _, q := range p.preds {
        \t\tif !q.Test(x) {
        \t\t\treturn false
        \t\t}
        \t}
        \treturn true
        }
        func (p *predAndSlice_) Test(x any) bool { return p.Apply(x) }
        func (p *predAndSlice_) TestNode_(n uint32) bool {
        \tCheckNotNullPtr(p)
        \tfor _, q := range p.preds {
        \t\tif !PredTestNode(q, n) {
        \t\t\treturn false
        \t\t}
        \t}
        \treturn true
        }
        func (p *predAndSlice_) Equals(o any) bool {
        \tCheckNotNullPtr(p)
        \top, ok := o.(*predAndSlice_)
        \tif !ok || op == nil || len(p.preds) != len(op.preds) {
        \t\treturn false
        \t}
        \tfor i := range p.preds {
        \t\tif !javalang.ObjectsEquals(p.preds[i], op.preds[i]) {
        \t\t\treturn false
        \t\t}
        \t}
        \treturn true
        }
        func (p *predAndSlice_) HashCode() int32 {
        \tCheckNotNullPtr(p)
        \th := int32(1)
        \tfor _, q := range p.preds {
        \t\th = 31*h + javalang.ObjectHashCode(q)
        \t}
        \treturn h + 0x12472c2c
        }
        func (p *predAndSlice_) String() string { CheckNotNullPtr(p); return "Predicates.and(...)" }
        func (p *predAndSlice_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predAnd2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predAndSlice_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predOr2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predAndSlice_) Negate() Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predNot_{inner: p}
        }

        type predOrSlice_ struct {
        \tpreds []Ext_java_util_function_Predicate
        }

        func (*predOrSlice_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (*predOrSlice_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
        func (p *predOrSlice_) IsJavaNil_() bool                                      { return p == nil }
        func (p *predOrSlice_) Apply(x any) bool {
        \tCheckNotNullPtr(p)
        \tfor _, q := range p.preds {
        \t\tif q.Test(x) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }
        func (p *predOrSlice_) Test(x any) bool { return p.Apply(x) }
        func (p *predOrSlice_) TestNode_(n uint32) bool {
        \tCheckNotNullPtr(p)
        \tfor _, q := range p.preds {
        \t\tif PredTestNode(q, n) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }
        func (p *predOrSlice_) Equals(o any) bool {
        \tCheckNotNullPtr(p)
        \top, ok := o.(*predOrSlice_)
        \tif !ok || op == nil || len(p.preds) != len(op.preds) {
        \t\treturn false
        \t}
        \tfor i := range p.preds {
        \t\tif !javalang.ObjectsEquals(p.preds[i], op.preds[i]) {
        \t\t\treturn false
        \t\t}
        \t}
        \treturn true
        }
        func (p *predOrSlice_) HashCode() int32 {
        \tCheckNotNullPtr(p)
        \th := int32(1)
        \tfor _, q := range p.preds {
        \t\th = 31*h + javalang.ObjectHashCode(q)
        \t}
        \treturn h + 0x053c91cf
        }
        func (p *predOrSlice_) String() string { CheckNotNullPtr(p); return "Predicates.or(...)" }
        func (p *predOrSlice_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predAnd2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predOrSlice_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predOr2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predOrSlice_) Negate() Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predNot_{inner: p}
        }

        func collectPredIter_(src Ext_java_lang_Iterable) []Ext_java_util_function_Predicate {
        \tCheckNotNull(src)
        \tif si, ok := src.(SliceIterable); ok {
        \t\traw := si.IterSlice_()
        \t\tout := make([]Ext_java_util_function_Predicate, len(raw))
        \t\tfor i, x := range raw {
        \t\t\tout[i] = AsIface[Ext_java_util_function_Predicate](CheckNotNull(x))
        \t\t}
        \t\treturn out
        \t}
        \tvar out []Ext_java_util_function_Predicate
        \tfor it := NewIterCursor(src); it.HasNext(); {
        \t\tout = append(out, AsIface[Ext_java_util_function_Predicate](CheckNotNull(it.Next())))
        \t}
        \treturn out
        }

        func guavaPredAndFromList_(elems []Ext_java_util_function_Predicate) Ext_com_google_common_base_Predicate {
        \tswitch len(elems) {
        \tcase 0:
        \t\treturn predAlwaysTrueIface_
        \tcase 1:
        \t\tif gp, ok := elems[0].(Ext_com_google_common_base_Predicate); ok {
        \t\t\treturn gp
        \t\t}
        \tcase 2:
        \t\treturn &predAnd2_{p1: elems[0], p2: elems[1]}
        \t}
        \treturn &predAndSlice_{preds: elems}
        }

        func guavaPredOrFromList_(elems []Ext_java_util_function_Predicate) Ext_com_google_common_base_Predicate {
        \tswitch len(elems) {
        \tcase 0:
        \t\treturn predAlwaysFalseIface_
        \tcase 1:
        \t\tif gp, ok := elems[0].(Ext_com_google_common_base_Predicate); ok {
        \t\t\treturn gp
        \t\t}
        \tcase 2:
        \t\treturn &predOr2_{p1: elems[0], p2: elems[1]}
        \t}
        \treturn &predOrSlice_{preds: elems}
        }

        func guavaPredAndIter(src Ext_java_lang_Iterable) Ext_com_google_common_base_Predicate {
        \treturn guavaPredAndFromList_(collectPredIter_(src))
        }

        func guavaPredOrIter(src Ext_java_lang_Iterable) Ext_com_google_common_base_Predicate {
        \treturn guavaPredOrFromList_(collectPredIter_(src))
        }

        func guavaPredAndSlice[E any](src []E) Ext_com_google_common_base_Predicate {
        \tCheckNotNullSlice(src)
        \telems := make([]Ext_java_util_function_Predicate, len(src))
        \tfor i, x := range src {
        \t\telems[i] = AsIface[Ext_java_util_function_Predicate](CheckNotNull(x))
        \t}
        \treturn guavaPredAndFromList_(elems)
        }

        func guavaPredOrSlice[E any](src []E) Ext_com_google_common_base_Predicate {
        \tCheckNotNullSlice(src)
        \telems := make([]Ext_java_util_function_Predicate, len(src))
        \tfor i, x := range src {
        \t\telems[i] = AsIface[Ext_java_util_function_Predicate](CheckNotNull(x))
        \t}
        \treturn guavaPredOrFromList_(elems)
        }

        type predInCollection_ struct {
        \tcol Ext_java_util_Collection
        }

        func (*predInCollection_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (*predInCollection_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
        func (p *predInCollection_) IsJavaNil_() bool                                      { return p == nil }
        func (p *predInCollection_) Apply(x any) (matched bool) {
        \tCheckNotNullPtr(p)
        \tdefer func() {
        \t\tif r := recover(); r != nil {
        \t\t\tswitch r.(type) {
        \t\t\tcase *javalang.NullPointerException, *javalang.ClassCastException:
        \t\t\t\tmatched = false
        \t\t\tdefault:
        \t\t\t\tpanic(r)
        \t\t\t}
        \t\t}
        \t}()
        \treturn predColContainsAny_(p.col, ToAny(x))
        }
        func (p *predInCollection_) Test(x any) bool { return p.Apply(x) }
        func (p *predInCollection_) TestNode_(n uint32) bool {
        \tCheckNotNullPtr(p)
        \tif fs, ok := p.col.(FastNodeSet); ok {
        \t\treturn fs.SetContainsNode_(n)
        \t}
        \treturn p.Apply(boxNodeById(n))
        }
        func (p *predInCollection_) Equals(o any) bool {
        \tCheckNotNullPtr(p)
        \top, ok := o.(*predInCollection_)
        \treturn ok && op != nil && javalang.ObjectsEquals(p.col, op.col)
        }
        func (p *predInCollection_) HashCode() int32 { CheckNotNullPtr(p); return javalang.ObjectHashCode(p.col) }
        func (p *predInCollection_) String() string {
        \tCheckNotNullPtr(p)
        \treturn "Predicates.in(" + javalang.ToJavaString(p.col) + ")"
        }
        func (p *predInCollection_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predAnd2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predInCollection_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predOr2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predInCollection_) Negate() Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predNot_{inner: p}
        }

        func guavaPredInCollection(col Ext_java_util_Collection) Ext_com_google_common_base_Predicate {
        \tCheckNotNull(col)
        \treturn &predInCollection_{col: col}
        }

        type predInstanceOf_ struct {
        \tcls *javalang.Class
        }

        func (*predInstanceOf_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (*predInstanceOf_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
        func (p *predInstanceOf_) IsJavaNil_() bool                                      { return p == nil }
        func (p *predInstanceOf_) Apply(x any) bool                                      { CheckNotNullPtr(p); return CheckTypePredicate(p.cls, x) }
        func (p *predInstanceOf_) Test(x any) bool                                       { CheckNotNullPtr(p); return CheckTypePredicate(p.cls, x) }
        func (p *predInstanceOf_) TestNode_(n uint32) bool {
        \tCheckNotNullPtr(p)
        \tif n == 0 {
        \t\treturn false
        \t}
        \treturn CheckTypePredicate(p.cls, boxNodeById(n))
        }
        func (p *predInstanceOf_) Equals(o any) bool {
        \tCheckNotNullPtr(p)
        \top, ok := o.(*predInstanceOf_)
        \treturn ok && op != nil && p.cls == op.cls
        }
        func (p *predInstanceOf_) HashCode() int32 { CheckNotNullPtr(p); return javalang.ObjectHashCode(p.cls) }
        func (p *predInstanceOf_) String() string {
        \tCheckNotNullPtr(p)
        \treturn "Predicates.instanceOf(" + p.cls.GetName() + ")"
        }
        func (p *predInstanceOf_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predAnd2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predInstanceOf_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predOr2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predInstanceOf_) Negate() Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predNot_{inner: p}
        }

        func guavaPredInstanceOf(cls *javalang.Class) Ext_com_google_common_base_Predicate {
        \tCheckNotNullPtr(cls)
        \treturn &predInstanceOf_{cls: cls}
        }

        type predCompose_ struct {
        \tpred Ext_com_google_common_base_Predicate
        \tfn   Ext_com_google_common_base_Function
        }

        func (*predCompose_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
        func (*predCompose_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
        func (p *predCompose_) IsJavaNil_() bool                                      { return p == nil }
        func (p *predCompose_) Apply(x any) bool                                      { CheckNotNullPtr(p); return p.pred.Apply(p.fn.Apply(x)) }
        func (p *predCompose_) Test(x any) bool                                       { CheckNotNullPtr(p); return p.pred.Apply(p.fn.Apply(x)) }
        func (p *predCompose_) TestNode_(n uint32) bool {
        \tCheckNotNullPtr(p)
        \treturn p.pred.Apply(p.fn.Apply(boxNodeById(n)))
        }
        func (p *predCompose_) Equals(o any) bool {
        \tCheckNotNullPtr(p)
        \top, ok := o.(*predCompose_)
        \treturn ok && op != nil && javalang.ObjectsEquals(p.fn, op.fn) && javalang.ObjectsEquals(p.pred, op.pred)
        }
        func (p *predCompose_) HashCode() int32 {
        \tCheckNotNullPtr(p)
        \treturn javalang.ObjectHashCode(p.fn) ^ javalang.ObjectHashCode(p.pred)
        }
        func (p *predCompose_) String() string {
        \tCheckNotNullPtr(p)
        \treturn javalang.ToJavaString(p.pred) + "(" + javalang.ToJavaString(p.fn) + ")"
        }
        func (p *predCompose_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predAnd2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predCompose_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predOr2_{p1: p, p2: CheckNotNull(o)}
        }
        func (p *predCompose_) Negate() Ext_java_util_function_Predicate {
        \tCheckNotNullPtr(p)
        \treturn &predNot_{inner: p}
        }

        func guavaPredCompose(pred Ext_com_google_common_base_Predicate, fn Ext_com_google_common_base_Function) Ext_com_google_common_base_Predicate {
        \tCheckNotNull(pred)
        \tCheckNotNull(fn)
        \treturn &predCompose_{pred: pred, fn: fn}
        }
        """;
  }

  private String buildCoreRuntimeGoPart1Tail() {
    return """

        // --- Class Literal Cache ---

        var classCache sync.Map // string -> *javalang.Class

        func ClassLiteral(qualName, simpleName string) *javalang.Class {
        \tif c, ok := classCache.Load(qualName); ok {
        \t\treturn c.(*javalang.Class)
        \t}
        \tc := javalang.NewClass(qualName, simpleName)
        \tactual, _ := classCache.LoadOrStore(qualName, c)
        \treturn actual.(*javalang.Class)
        }

        func GetObjectClass(v any) *javalang.Class {
        \tif v == nil || javalang.IsNil(v) {
        \t\tpanic(javalang.NewNullPointerException())
        \t}
        \tif s, ok := v.(string); ok {
        \t\tif s == NullString {
        \t\t\tpanic(javalang.NewNullPointerException())
        \t\t}
        \t\treturn ClassLiteral("java.lang.String", "String")
        \t}
        \tif n, ok := v.(nodestore.Node); ok {
        \t\tif n == 0 {
        \t\t\tpanic(javalang.NewNullPointerException())
        \t\t}
        \t\tswitch n.KindFlags() {
        \t\tcase 1:
        \t\t\treturn ClassLiteral("com.google.javascript.rhino.Node$NumberNode", "NumberNode")
        \t\tcase 2:
        \t\t\treturn ClassLiteral("com.google.javascript.rhino.Node$BigIntNode", "BigIntNode")
        \t\tcase 3:
        \t\t\treturn ClassLiteral("com.google.javascript.rhino.Node$StringNode", "StringNode")
        \t\tcase 4:
        \t\t\treturn ClassLiteral("com.google.javascript.rhino.Node$TemplateLiteralSubstringNode", "TemplateLiteralSubstringNode")
        \t\tdefault:
        \t\t\treturn ClassLiteral("com.google.javascript.rhino.Node", "Node")
        \t\t}
        \t}
        \tif ci, ok := v.(interface{ JavaClassInfo_() (string, string) }); ok {
        \t\tqn, sn := ci.JavaClassInfo_()
        \t\treturn ClassLiteral(qn, sn)
        \t}
        \tt := reflect.TypeOf(v)
        \treturn ClassLiteral(t.String(), t.Name())
        }

        // --- Throwable Helpers ---

        func InitUserThrowable(t javalang.Throwable, className string, args ...any) {
        \tmsg := ""
        \thasMsg := false
        \tvar cause javalang.Throwable
        \tfor _, a := range args {
        \t\tif s, ok := a.(string); ok {
        \t\t\tif s != NullString {
        \t\t\t\tmsg = s
        \t\t\t\thasMsg = true
        \t\t\t}
        \t\t} else if th, ok := a.(javalang.Throwable); ok {
        \t\t\tcause = th
        \t\t\tif !hasMsg && th != nil && !javalang.IsNil(th) {
        \t\t\t\tmsg = th.Error()
        \t\t\t\thasMsg = true
        \t\t\t}
        \t\t} else if a != nil && !javalang.IsNil(a) {
        \t\t\tmsg = StrVal(a)
        \t\t\thasMsg = true
        \t\t}
        \t}
        \tif hasMsg {
        \t\tjavalang.InitThrowable(t, className, msg, cause)
        \t} else {
        \t\tjavalang.InitThrowable(t, className, "", cause)
        \t}
        }

        func ConstructShimThrowable(qualName string, args ...any) javalang.Throwable {
        \tvar t javalang.Throwable
        \tswitch qualName {
        \tcase "java.lang.Throwable", "java.lang.Exception":
        \t\tt = &javalang.JavaException{}
        \tcase "java.lang.RuntimeException":
        \t\tt = &javalang.RuntimeException{}
        \tcase "java.lang.IllegalArgumentException":
        \t\tt = &javalang.IllegalArgumentException{}
        \tcase "java.lang.IllegalStateException":
        \t\tt = &javalang.IllegalStateException{}
        \tcase "java.lang.NullPointerException":
        \t\tt = &javalang.NullPointerException{}
        \tcase "java.lang.IndexOutOfBoundsException":
        \t\tt = &javalang.IndexOutOfBoundsException{}
        \tcase "java.lang.ArrayIndexOutOfBoundsException":
        \t\tt = &javalang.ArrayIndexOutOfBoundsException{}
        \tcase "java.lang.StringIndexOutOfBoundsException":
        \t\tt = &javalang.StringIndexOutOfBoundsException{}
        \tcase "java.lang.UnsupportedOperationException":
        \t\tt = &javalang.UnsupportedOperationException{}
        \tcase "java.lang.ClassCastException":
        \t\tt = &javalang.ClassCastException{}
        \tcase "java.util.NoSuchElementException":
        \t\tt = &javalang.NoSuchElementException{}
        \tcase "java.lang.ArithmeticException":
        \t\tt = &javalang.ArithmeticException{}
        \tcase "java.lang.NumberFormatException":
        \t\tt = &javalang.NumberFormatException{}
        \tcase "java.lang.Error":
        \t\tt = &javalang.JavaLangError{}
        \tcase "java.lang.AssertionError":
        \t\tt = &javalang.AssertionError{}
        \tcase "java.io.IOException":
        \t\tt = &javaio.IOException{}
        \tcase "java.io.FileNotFoundException":
        \t\tt = &javaio.FileNotFoundException{}
        \tcase "java.io.EOFException":
        \t\tt = &javaio.EOFException{}
        \tdefault:
        \t\tt = &javalang.RuntimeException{}
        \t}
        \tInitUserThrowable(t, qualName, args...)
        \treturn t
        }

        func ThrowableGetMessage(t javalang.Throwable) string {
        \tif t == nil || javalang.IsNil(t) {
        \t\tpanic(javalang.NewNullPointerException())
        \t}
        \treturn t.GetMessage()
        }

        // --- String & CharSequence Intrinsics ---

        func StringEquals(s string, other any) bool {
        \tif s == NullString {
        \t\tpanic(javalang.NewNullPointerException())
        \t}
        \to, ok := other.(string)
        \treturn ok && o != NullString && s == o
        }

        func StringFormat(format string, args ...any) string {
        \tCheckNotNull(format)
        \tvar out strings.Builder
        \tout.Grow(len(format) + 16)
        \tseqIdx := 0
        \tlastIdx := 0
        \tfor i := 0; i < len(format); {
        \t\tif format[i] != '%' {
        \t\t\tj := strings.IndexByte(format[i:], '%')
        \t\t\tif j < 0 {
        \t\t\t\tout.WriteString(format[i:])
        \t\t\t\tbreak
        \t\t\t}
        \t\t\tout.WriteString(format[i : i+j])
        \t\t\ti += j
        \t\t}
        \t\ti++ // consume '%'
        \t\tif i >= len(format) {
        \t\t\tout.WriteByte('%')
        \t\t\tbreak
        \t\t}
        \t\tif format[i] == '%' {
        \t\t\tout.WriteByte('%')
        \t\t\ti++
        \t\t\tcontinue
        \t\t}
        \t\tif format[i] == 'n' {
        \t\t\tout.WriteByte('\\n')
        \t\t\ti++
        \t\t\tcontinue
        \t\t}
        \t\t// Optional [argument_index$]
        \t\texplicitIdx := -1
        \t\tusePrev := false
        \t\tsave := i
        \t\tif format[i] >= '1' && format[i] <= '9' {
        \t\t\tn := 0
        \t\t\tfor i < len(format) && format[i] >= '0' && format[i] <= '9' {
        \t\t\t\tn = n*10 + int(format[i]-'0')
        \t\t\t\ti++
        \t\t\t}
        \t\t\tif i < len(format) && format[i] == '$' {
        \t\t\t\texplicitIdx = n - 1
        \t\t\t\ti++
        \t\t\t} else {
        \t\t\t\ti = save
        \t\t\t}
        \t\t}
        \t\t// Flags
        \t\tvar flags strings.Builder
        \t\tleftAlign := false
        \t\tcommaFlag := false
        \t\tparenFlag := false
        \t\tfor i < len(format) {
        \t\t\tc := format[i]
        \t\t\tif c == '<' {
        \t\t\t\tusePrev = true
        \t\t\t\ti++
        \t\t\t} else if c == '-' {
        \t\t\t\tleftAlign = true
        \t\t\t\tflags.WriteByte(c)
        \t\t\t\ti++
        \t\t\t} else if c == ',' {
        \t\t\t\tcommaFlag = true
        \t\t\t\ti++
        \t\t\t} else if c == '(' {
        \t\t\t\tparenFlag = true
        \t\t\t\ti++
        \t\t\t} else if c == '#' || c == '+' || c == ' ' || c == '0' {
        \t\t\t\tflags.WriteByte(c)
        \t\t\t\ti++
        \t\t\t} else {
        \t\t\t\tbreak
        \t\t\t}
        \t\t}
        \t\t_ = commaFlag
        \t\t_ = parenFlag
        \t\t// Width
        \t\twidth := -1
        \t\twidthStart := i
        \t\tif i < len(format) && format[i] >= '1' && format[i] <= '9' {
        \t\t\tw := 0
        \t\t\tfor i < len(format) && format[i] >= '0' && format[i] <= '9' {
        \t\t\t\tw = w*10 + int(format[i]-'0')
        \t\t\t\ti++
        \t\t\t}
        \t\t\twidth = w
        \t\t}
        \t\twidthStr := format[widthStart:i]
        \t\t// Precision
        \t\tprec := -1
        \t\tprecStr := ""
        \t\tif i < len(format) && format[i] == '.' {
        \t\t\tpStart := i
        \t\t\ti++
        \t\t\tp := 0
        \t\t\tfor i < len(format) && format[i] >= '0' && format[i] <= '9' {
        \t\t\t\tp = p*10 + int(format[i]-'0')
        \t\t\t\ti++
        \t\t\t}
        \t\t\tprec = p
        \t\t\tprecStr = format[pStart:i]
        \t\t}
        \t\tif i >= len(format) {
        \t\t\tbreak
        \t\t}
        \t\tconvCh := format[i]
        \t\ti++
        \t\tvar argIdx int
        \t\tif explicitIdx >= 0 {
        \t\t\targIdx = explicitIdx
        \t\t} else if usePrev {
        \t\t\targIdx = lastIdx
        \t\t} else {
        \t\t\targIdx = seqIdx
        \t\t\tseqIdx++
        \t\t}
        \t\tlastIdx = argIdx
        \t\tvar arg any
        \t\tif argIdx >= 0 && argIdx < len(args) {
        \t\t\targ = args[argIdx]
        \t\t}
        \t\tpadAndWrite := func(s string) {
        \t\t\tif prec >= 0 && javalang.StringLength(s) > prec {
        \t\t\t\ts = javalang.StringSubstring(s, 0, prec)
        \t\t\t}
        \t\t\tsLen := javalang.StringLength(s)
        \t\t\tif width > sLen {
        \t\t\t\tpad := strings.Repeat(" ", width-sLen)
        \t\t\t\tif leftAlign {
        \t\t\t\t\tout.WriteString(s)
        \t\t\t\t\tout.WriteString(pad)
        \t\t\t\t} else {
        \t\t\t\t\tout.WriteString(pad)
        \t\t\t\t\tout.WriteString(s)
        \t\t\t\t}
        \t\t\t} else {
        \t\t\t\tout.WriteString(s)
        \t\t\t}
        \t\t}
        \t\tswitch convCh {
        \t\tcase 's', 'S':
        \t\t\ts := StrVal(arg)
        \t\t\tif convCh == 'S' {
        \t\t\t\ts = strings.ToUpper(s)
        \t\t\t}
        \t\t\tpadAndWrite(s)
        \t\tcase 'b', 'B':
        \t\t\tbv := false
        \t\t\tif arg != nil && !javalang.IsNil(arg) {
        \t\t\t\tswitch x := arg.(type) {
        \t\t\t\tcase bool:
        \t\t\t\t\tbv = x
        \t\t\t\tcase *BoxBoolean:
        \t\t\t\t\tbv = x != nil && x.Val
        \t\t\t\tdefault:
        \t\t\t\t\tbv = true
        \t\t\t\t}
        \t\t\t}
        \t\t\ts := "false"
        \t\t\tif bv {
        \t\t\t\ts = "true"
        \t\t\t}
        \t\t\tif convCh == 'B' {
        \t\t\t\ts = strings.ToUpper(s)
        \t\t\t}
        \t\t\tpadAndWrite(s)
        \t\tcase 'h', 'H':
        \t\t\ts := "null"
        \t\t\tif arg != nil && !javalang.IsNil(arg) {
        \t\t\t\ts = strconv.FormatUint(uint64(uint32(javalang.ObjectHashCode(arg))), 16)
        \t\t\t\tif convCh == 'H' {
        \t\t\t\t\ts = strings.ToUpper(s)
        \t\t\t\t}
        \t\t\t}
        \t\t\tpadAndWrite(s)
        \t\tcase 'c', 'C':
        \t\t\tif arg == nil || javalang.IsNil(arg) {
        \t\t\t\tpadAndWrite("null")
        \t\t\t} else {
        \t\t\t\tvar r rune
        \t\t\t\tswitch x := arg.(type) {
        \t\t\t\tcase uint16:
        \t\t\t\t\tr = rune(x)
        \t\t\t\tcase *BoxChar:
        \t\t\t\t\tr = rune(x.Val)
        \t\t\t\tcase int32:
        \t\t\t\t\tr = rune(x)
        \t\t\t\tcase *BoxInt:
        \t\t\t\t\tr = rune(x.Val)
        \t\t\t\tcase int8:
        \t\t\t\t\tr = rune(uint8(x))
        \t\t\t\tcase *BoxByte:
        \t\t\t\t\tr = rune(uint8(x.Val))
        \t\t\t\tcase int16:
        \t\t\t\t\tr = rune(uint16(x))
        \t\t\t\tcase *BoxShort:
        \t\t\t\t\tr = rune(uint16(x.Val))
        \t\t\t\t}
        \t\t\t\ts := string(r)
        \t\t\t\tif convCh == 'C' {
        \t\t\t\t\ts = strings.ToUpper(s)
        \t\t\t\t}
        \t\t\t\tprec = -1
        \t\t\t\tpadAndWrite(s)
        \t\t\t}
        \t\tcase 'd', 'o', 'x', 'X':
        \t\t\tif arg == nil || javalang.IsNil(arg) {
        \t\t\t\tpadAndWrite("null")
        \t\t\t} else {
        \t\t\t\tgoVerb := "%" + flags.String() + widthStr + string(convCh)
        \t\t\t\tvar num any = arg
        \t\t\t\tif convCh == 'd' {
        \t\t\t\t\tswitch x := arg.(type) {
        \t\t\t\t\tcase *BoxInt:
        \t\t\t\t\t\tnum = x.Val
        \t\t\t\t\tcase *BoxLong:
        \t\t\t\t\t\tnum = x.Val
        \t\t\t\t\tcase *BoxShort:
        \t\t\t\t\t\tnum = x.Val
        \t\t\t\t\tcase *BoxByte:
        \t\t\t\t\t\tnum = x.Val
        \t\t\t\t\tcase Number:
        \t\t\t\t\t\tnum = x.LongValue()
        \t\t\t\t\t}
        \t\t\t\t} else {
        \t\t\t\t\tswitch x := arg.(type) {
        \t\t\t\t\tcase int8:
        \t\t\t\t\t\tnum = uint8(x)
        \t\t\t\t\tcase *BoxByte:
        \t\t\t\t\t\tnum = uint8(x.Val)
        \t\t\t\t\tcase int16:
        \t\t\t\t\t\tnum = uint16(x)
        \t\t\t\t\tcase *BoxShort:
        \t\t\t\t\t\tnum = uint16(x.Val)
        \t\t\t\t\tcase int32:
        \t\t\t\t\t\tnum = uint32(x)
        \t\t\t\t\tcase *BoxInt:
        \t\t\t\t\t\tnum = uint32(x.Val)
        \t\t\t\t\tcase int64:
        \t\t\t\t\t\tnum = uint64(x)
        \t\t\t\t\tcase *BoxLong:
        \t\t\t\t\t\tnum = uint64(x.Val)
        \t\t\t\t\tcase Number:
        \t\t\t\t\t\tnum = uint64(x.LongValue())
        \t\t\t\t\t}
        \t\t\t\t}
        \t\t\t\tout.WriteString(fmt.Sprintf(goVerb, num))
        \t\t\t}
        \t\tcase 'e', 'E', 'f', 'g', 'G', 'a', 'A':
        \t\t\tif arg == nil || javalang.IsNil(arg) {
        \t\t\t\tpadAndWrite("null")
        \t\t\t} else {
        \t\t\t\tvar fv float64
        \t\t\t\tswitch x := arg.(type) {
        \t\t\t\tcase float64:
        \t\t\t\t\tfv = x
        \t\t\t\tcase float32:
        \t\t\t\t\tfv = float64(x)
        \t\t\t\tcase *BoxDouble:
        \t\t\t\t\tfv = x.Val
        \t\t\t\tcase *BoxFloat:
        \t\t\t\t\tfv = float64(x.Val)
        \t\t\t\tcase Number:
        \t\t\t\t\tfv = x.DoubleValue()
        \t\t\t\t}
        \t\t\t\tif math.IsNaN(fv) {
        \t\t\t\t\tprec = -1
        \t\t\t\t\tpadAndWrite("NaN")
        \t\t\t\t} else if math.IsInf(fv, 1) {
        \t\t\t\t\tprec = -1
        \t\t\t\t\tif convCh == 'E' || convCh == 'G' || convCh == 'A' {
        \t\t\t\t\t\tpadAndWrite("INFINITY")
        \t\t\t\t\t} else {
        \t\t\t\t\t\tpadAndWrite("Infinity")
        \t\t\t\t\t}
        \t\t\t\t} else if math.IsInf(fv, -1) {
        \t\t\t\t\tprec = -1
        \t\t\t\t\tif convCh == 'E' || convCh == 'G' || convCh == 'A' {
        \t\t\t\t\t\tpadAndWrite("-INFINITY")
        \t\t\t\t\t} else {
        \t\t\t\t\t\tpadAndWrite("-Infinity")
        \t\t\t\t\t}
        \t\t\t\t} else {
        \t\t\t\t\tvCh := convCh
        \t\t\t\t\tif vCh == 'a' {
        \t\t\t\t\t\tvCh = 'x'
        \t\t\t\t\t} else if vCh == 'A' {
        \t\t\t\t\t\tvCh = 'X'
        \t\t\t\t\t}
        \t\t\t\t\tgoVerb := "%" + flags.String() + widthStr + precStr + string(vCh)
        \t\t\t\t\tout.WriteString(fmt.Sprintf(goVerb, fv))
        \t\t\t\t}
        \t\t\t}
        \t\tdefault:
        \t\t\tpadAndWrite(StrVal(arg))
        \t\t}
        \t}
        \treturn out.String()
        }

        func JStrLen(s string) int32 {
        \tif len(s) < 48 {
        \t\tfor i := 0; i < len(s); i++ {
        \t\t\tif s[i] >= 0x80 {
        \t\t\t\tif s == NullString {
        \t\t\t\t\tThrowNPE()
        \t\t\t\t}
        \t\t\t\treturn int32(javalang.StringLength(s))
        \t\t\t}
        \t\t}
        \t\treturn int32(len(s))
        \t}
        \tif s == NullString {
        \t\tThrowNPE()
        \t}
        \treturn int32(javalang.StringLength(s))
        }
        func JStrCharAt(s string, idx int32) uint16 {
        \tCheckNotNullString(s)
        \treturn javalang.StringCharAt(s, int(idx))
        }
        func JStrCodePointAt(s string, idx int32) int32 {
        \tCheckNotNullString(s)
        \treturn javalang.StringCodePointAt(s, int(idx))
        }
        func JStrSub(s string, begin, end int32) string {
        \tCheckNotNullString(s)
        \treturn javalang.StringSubstring(s, int(begin), int(end))
        }
        func JStrSubFrom(s string, begin int32) string {
        \tCheckNotNullString(s)
        \treturn javalang.StringSubstringFrom(s, int(begin))
        }
        func JStrIndexOf(s, sub string) int32 {
        \tCheckNotNullString(s); CheckNotNullString(sub)
        \treturn int32(javalang.StringIndexOf(s, sub))
        }
        func JStrIndexOfFrom(s, sub string, from int32) int32 {
        \tCheckNotNullString(s); CheckNotNullString(sub)
        \treturn int32(javalang.StringIndexOfFrom(s, sub, int(from)))
        }
        func JStrIndexOfCp(s string, cp int32) int32 {
        \tCheckNotNullString(s)
        \tif cp >= 0 && cp <= 0xFFFF {
        \t\treturn int32(javalang.StringIndexOfChar(s, uint16(cp)))
        \t}
        \tif cp > 0xFFFF && cp <= 0x10FFFF {
        \t\treturn JStrIndexOf(s, string(rune(cp)))
        \t}
        \treturn -1
        }
        func JStrIndexOfCpFrom(s string, cp, from int32) int32 {
        \tCheckNotNullString(s)
        \tif cp >= 0 && cp <= 0xFFFF {
        \t\treturn JStrIndexOfFrom(s, javalang.StringFromChar(uint16(cp)), from)
        \t}
        \tif cp > 0xFFFF && cp <= 0x10FFFF {
        \t\treturn JStrIndexOfFrom(s, string(rune(cp)), from)
        \t}
        \treturn -1
        }
        func JStrLastIndexOf(s, sub string) int32 {
        \tCheckNotNullString(s); CheckNotNullString(sub)
        \treturn int32(javalang.StringLastIndexOf(s, sub))
        }
        func JStrLastIndexOfFrom(s, sub string, from int32) int32 {
        \tCheckNotNullString(s); CheckNotNullString(sub)
        \treturn int32(javalang.StringLastIndexOfFrom(s, sub, int(from)))
        }
        func JStrLastIndexOfCp(s string, cp int32) int32 {
        \tCheckNotNullString(s)
        \tif cp >= 0 && cp <= 0xFFFF {
        \t\treturn int32(javalang.StringLastIndexOfCharFrom(s, uint16(cp), javalang.StringLength(s)))
        \t}
        \tif cp > 0xFFFF && cp <= 0x10FFFF {
        \t\treturn JStrLastIndexOf(s, string(rune(cp)))
        \t}
        \treturn -1
        }
        func JStrLastIndexOfCpFrom(s string, cp, from int32) int32 {
        \tCheckNotNullString(s)
        \tif cp >= 0 && cp <= 0xFFFF {
        \t\treturn int32(javalang.StringLastIndexOfCharFrom(s, uint16(cp), int(from)))
        \t}
        \tif cp > 0xFFFF && cp <= 0x10FFFF {
        \t\treturn JStrLastIndexOfFrom(s, string(rune(cp)), from)
        \t}
        \treturn -1
        }
        func JStrStartsWithOffset(s, prefix string, offset int32) bool {
        \tCheckNotNullString(s); CheckNotNullString(prefix)
        \treturn javalang.StringRegionMatches(s, int(offset), prefix, 0, javalang.StringLength(prefix))
        }
        func JStrReplaceChar(s string, oldCh, newCh uint16) string {
        \tCheckNotNullString(s)
        \tif oldCh == newCh {
        \t\treturn s
        \t}
        \treturn javalang.StringReplace(s, javalang.StringFromChar(oldCh), javalang.StringFromChar(newCh))
        }
        func JStrCompareTo(a, b string) int32 {
        \tCheckNotNullString(a); CheckNotNullString(b)
        \treturn int32(javalang.StringCompareTo(a, b))
        }
        func JStrSplit(s, regex string) []string {
        \tCheckNotNullString(s); CheckNotNullString(regex)
        \treturn javaregex.StringSplit(s, regex)
        }
        func JStrSplitLimit(s, regex string, limit int32) []string {
        \tCheckNotNullString(s); CheckNotNullString(regex)
        \treturn javaregex.StringSplitLimit(s, regex, int(limit))
        }
        func JStrMatches(s, regex string) bool {
        \tCheckNotNullString(s); CheckNotNullString(regex)
        \treturn javaregex.StringMatches(s, regex)
        }
        func JStrReplaceAll(s, regex, repl string) string {
        \tCheckNotNullString(s); CheckNotNullString(regex); CheckNotNullString(repl)
        \treturn javaregex.StringReplaceAll(s, regex, repl)
        }
        func JStrReplaceFirst(s, regex, repl string) string {
        \tCheckNotNullString(s); CheckNotNullString(regex); CheckNotNullString(repl)
        \treturn javaregex.StringReplaceFirst(s, regex, repl)
        }
        func JStrRegionMatches(s string, ignoreCase bool, toffset int32, other string, ooffset, length int32) bool {
        \tCheckNotNullString(s); CheckNotNullString(other)
        \tif !ignoreCase {
        \t\treturn javalang.StringRegionMatches(s, int(toffset), other, int(ooffset), int(length))
        \t}
        \tif toffset < 0 || ooffset < 0 || length < 0 {
        \t\treturn false
        \t}
        \tsLen := javalang.StringLength(s)
        \toLen := javalang.StringLength(other)
        \tif int(toffset)+int(length) > sLen || int(ooffset)+int(length) > oLen {
        \t\treturn false
        \t}
        \ta := javalang.StringSubstring(s, int(toffset), int(toffset+length))
        \tb := javalang.StringSubstring(other, int(ooffset), int(ooffset+length))
        \treturn strings.EqualFold(a, b)
        }
        func JStrGetBytes(s string) []int8 {
        \tCheckNotNullString(s)
        \tout := make([]int8, len(s))
        \tfor i := 0; i < len(s); i++ {
        \t\tout[i] = int8(s[i])
        \t}
        \treturn out
        }
        """;
  }

  private String buildCoreRuntimeGoPart1B() {
    return """

        // --- O(1) UTF-16 Length & WTF-8 StringBuilder Core ---

        type jStringBuilderCore struct {
        \tbuf         []byte
        \tcachedStr   string
        \tutf16Len    int32  // >= 0 exact UTF-16 code unit count; -1 dirty
        \tpendingHigh uint16 // 0 = none; 0xD800..0xDBFF = trailing lone high surrogate held out of buf
        \thasCache    bool
        }

        func (c *jStringBuilderCore) invalidateCache() {
        \tif c.hasCache {
        \t\tc.cachedStr = ""
        \t\tc.hasCache = false
        \t}
        }

        func (c *jStringBuilderCore) grow(capHint int) {
        \tif capHint > 0 && capHint <= (1<<24) && cap(c.buf) < capHint {
        \t\tnb := make([]byte, len(c.buf), capHint)
        \t\tcopy(nb, c.buf)
        \t\tc.buf = nb
        \t}
        }

        func (c *jStringBuilderCore) reset() {
        \tc.buf = c.buf[:0]
        \tc.cachedStr = ""
        \tc.hasCache = false
        \tc.utf16Len = 0
        \tc.pendingHigh = 0
        }

        func (c *jStringBuilderCore) resetTo(s string) {
        \tc.reset()
        \tc.appendString(s)
        }

        func (c *jStringBuilderCore) flushPendingHigh() {
        \tif c.pendingHigh != 0 {
        \t\tch := c.pendingHigh
        \t\tc.pendingHigh = 0
        \t\tc.invalidateCache()
        \t\tc.buf = append(c.buf,
        \t\t\tbyte(0xE0|(ch>>12)),
        \t\t\tbyte(0x80|((ch>>6)&0x3F)),
        \t\t\tbyte(0x80|(ch&0x3F)))
        \t}
        }

        func (c *jStringBuilderCore) str() string {
        \tif len(c.buf) == 0 && c.pendingHigh == 0 {
        \t\treturn ""
        \t}
        \tif c.pendingHigh == 0 {
        \t\tif !c.hasCache {
        \t\t\tc.cachedStr = string(c.buf)
        \t\t\tc.hasCache = true
        \t\t}
        \t\treturn c.cachedStr
        \t}
        \tch := c.pendingHigh
        \tout := make([]byte, len(c.buf)+3)
        \tcopy(out, c.buf)
        \tout[len(c.buf)] = byte(0xE0 | (ch >> 12))
        \tout[len(c.buf)+1] = byte(0x80 | ((ch >> 6) & 0x3F))
        \tout[len(c.buf)+2] = byte(0x80 | (ch & 0x3F))
        \treturn string(out)
        }

        func (c *jStringBuilderCore) length() int32 {
        \tif c.utf16Len < 0 {
        \t\tn := int32(javalang.StringLength(string(c.buf)))
        \t\tif c.pendingHigh != 0 {
        \t\t\tn++
        \t\t}
        \t\tc.utf16Len = n
        \t}
        \treturn c.utf16Len
        }

        func (c *jStringBuilderCore) appendChar(ch uint16) {
        \tif c.hasCache {
        \t\tc.cachedStr = ""
        \t\tc.hasCache = false
        \t}
        \tif c.pendingHigh == 0 && ch < 0x80 {
        \t\tc.buf = append(c.buf, byte(ch))
        \t\tif c.utf16Len >= 0 {
        \t\t\tc.utf16Len++
        \t\t}
        \t\treturn
        \t}
        \tif c.utf16Len >= 0 {
        \t\tc.utf16Len++
        \t}
        \tif c.pendingHigh != 0 {
        \t\tif ch >= 0xDC00 && ch <= 0xDFFF {
        \t\t\tr := utf16.DecodeRune(rune(c.pendingHigh), rune(ch))
        \t\t\tc.pendingHigh = 0
        \t\t\tvar tmp [4]byte
        \t\t\tn := utf8.EncodeRune(tmp[:], r)
        \t\t\tc.buf = append(c.buf, tmp[:n]...)
        \t\t\treturn
        \t\t}
        \t\tc.flushPendingHigh()
        \t}
        \tif ch >= 0xD800 && ch <= 0xDBFF {
        \t\tc.pendingHigh = ch
        \t\treturn
        \t}
        \tif ch < 0x80 {
        \t\tc.buf = append(c.buf, byte(ch))
        \t} else if ch < 0x800 {
        \t\tc.buf = append(c.buf, byte(0xC0|(ch>>6)), byte(0x80|(ch&0x3F)))
        \t} else {
        \t\tc.buf = append(c.buf,
        \t\t\tbyte(0xE0|(ch>>12)),
        \t\t\tbyte(0x80|((ch>>6)&0x3F)),
        \t\t\tbyte(0x80|(ch&0x3F)))
        \t}
        }

        func (c *jStringBuilderCore) appendString(raw string) {
        \ts := StrValString(raw)
        \tif len(s) == 0 {
        \t\treturn
        \t}
        \tif c.hasCache {
        \t\tc.cachedStr = ""
        \t\tc.hasCache = false
        \t}
        \tif c.pendingHigh == 0 && len(s) < 48 {
        \t\tascii := true
        \t\tfor i := 0; i < len(s); i++ {
        \t\t\tif s[i] >= 0x80 {
        \t\t\t\tascii = false
        \t\t\t\tbreak
        \t\t\t}
        \t\t}
        \t\tif ascii {
        \t\t\tc.buf = append(c.buf, s...)
        \t\t\tif c.utf16Len >= 0 {
        \t\t\t\tc.utf16Len += int32(len(s))
        \t\t\t}
        \t\t\treturn
        \t\t}
        \t}
        \tif c.utf16Len >= 0 {
        \t\tc.utf16Len += JStrLen(s)
        \t}
        \tif c.pendingHigh != 0 {
        \t\tif len(s) >= 3 && s[0] == 0xED && s[1] >= 0xB0 && s[1] <= 0xBF && (s[2]&0xC0) == 0x80 {
        \t\t\tlow := (uint16(s[0]&0x0F) << 12) | (uint16(s[1]&0x3F) << 6) | uint16(s[2]&0x3F)
        \t\t\tr := utf16.DecodeRune(rune(c.pendingHigh), rune(low))
        \t\t\tc.pendingHigh = 0
        \t\t\tvar tmp [4]byte
        \t\t\tn := utf8.EncodeRune(tmp[:], r)
        \t\t\tc.buf = append(c.buf, tmp[:n]...)
        \t\t\ts = s[3:]
        \t\t\tif len(s) == 0 {
        \t\t\t\treturn
        \t\t\t}
        \t\t} else {
        \t\t\tc.flushPendingHigh()
        \t\t}
        \t}
        \tn := len(s)
        \tif n >= 3 && s[n-3] == 0xED && s[n-2] >= 0xA0 && s[n-2] <= 0xAF && (s[n-1]&0xC0) == 0x80 {
        \t\thigh := (uint16(s[n-3]&0x0F) << 12) | (uint16(s[n-2]&0x3F) << 6) | uint16(s[n-1]&0x3F)
        \t\tc.buf = append(c.buf, s[:n-3]...)
        \t\tc.pendingHigh = high
        \t\treturn
        \t}
        \tc.buf = append(c.buf, s...)
        }

        func (c *jStringBuilderCore) appendASCII(s string) {
        \tif len(s) == 0 {
        \t\treturn
        \t}
        \tif c.hasCache {
        \t\tc.cachedStr = ""
        \t\tc.hasCache = false
        \t}
        \tif c.pendingHigh != 0 {
        \t\tc.flushPendingHigh()
        \t}
        \tc.buf = append(c.buf, s...)
        \tif c.utf16Len >= 0 {
        \t\tc.utf16Len += int32(len(s))
        \t}
        }

        func (c *jStringBuilderCore) appendInt64(v int64) {
        \tif c.hasCache {
        \t\tc.cachedStr = ""
        \t\tc.hasCache = false
        \t}
        \tif c.pendingHigh != 0 {
        \t\tc.flushPendingHigh()
        \t}
        \tvar tmp [20]byte
        \tb := strconv.AppendInt(tmp[:0], v, 10)
        \tc.buf = append(c.buf, b...)
        \tif c.utf16Len >= 0 {
        \t\tc.utf16Len += int32(len(b))
        \t}
        }

        func (c *jStringBuilderCore) appendInt32(v int32) {
        \tc.appendInt64(int64(v))
        }

        func (c *jStringBuilderCore) appendBool(v bool) {
        \tif v {
        \t\tc.appendASCII("true")
        \t} else {
        \t\tc.appendASCII("false")
        \t}
        }

        func (c *jStringBuilderCore) appendCodePoint(cp int32) {
        \tif cp >= 0 && cp <= 0xFFFF {
        \t\tc.appendChar(uint16(cp))
        \t\treturn
        \t}
        \tif cp > 0xFFFF && cp <= 0x10FFFF {
        \t\tif c.hasCache {
        \t\t\tc.cachedStr = ""
        \t\t\tc.hasCache = false
        \t\t}
        \t\tif c.pendingHigh != 0 {
        \t\t\tc.flushPendingHigh()
        \t\t}
        \t\tvar tmp [4]byte
        \t\tn := utf8.EncodeRune(tmp[:], rune(cp))
        \t\tc.buf = append(c.buf, tmp[:n]...)
        \t\tif c.utf16Len >= 0 {
        \t\t\tc.utf16Len += 2
        \t\t}
        \t\treturn
        \t}
        \tpanic(javalang.NewIllegalArgumentException("Not a valid Unicode code point"))
        }

        // byteOffsetOfUTF16 returns the byte offset in c.buf for UTF-16 index uIdx (0 <= uIdx <= effN),
        // or -1 if uIdx splits a 4-byte astral surrogate pair or encounters non-WTF-8 raw bytes.
        func (c *jStringBuilderCore) byteOffsetOfUTF16(uIdx, effN int32) int {
        \tif uIdx == 0 {
        \t\treturn 0
        \t}
        \tif uIdx == effN {
        \t\treturn len(c.buf)
        \t}
        \tif int(effN) == len(c.buf) {
        \t\treturn int(uIdx)
        \t}
        \tif uIdx > effN/2 {
        \t\tb := len(c.buf)
        \t\trem := effN - uIdx
        \t\tfor rem > 0 {
        \t\t\tif b <= 0 {
        \t\t\t\treturn -1
        \t\t\t}
        \t\t\tb--
        \t\t\tsteps := 1
        \t\t\tfor b > 0 && steps < 4 && (c.buf[b]&0xC0) == 0x80 {
        \t\t\t\tb--
        \t\t\t\tsteps++
        \t\t\t}
        \t\t\tlead := c.buf[b]
        \t\t\tif lead < 0x80 {
        \t\t\t\tif steps != 1 {
        \t\t\t\t\treturn -1
        \t\t\t\t}
        \t\t\t\trem--
        \t\t\t} else if lead >= 0xC2 && lead < 0xE0 && steps == 2 {
        \t\t\t\trem--
        \t\t\t} else if lead >= 0xE0 && lead < 0xF0 && steps == 3 {
        \t\t\t\trem--
        \t\t\t} else if lead >= 0xF0 && lead < 0xF5 && steps == 4 {
        \t\t\t\tif rem < 2 {
        \t\t\t\t\treturn -1
        \t\t\t\t}
        \t\t\t\trem -= 2
        \t\t\t} else {
        \t\t\t\treturn -1
        \t\t\t}
        \t\t}
        \t\treturn b
        \t}
        \tb := 0
        \trem := uIdx
        \tfor rem > 0 {
        \t\tif b >= len(c.buf) {
        \t\t\treturn -1
        \t\t}
        \t\tlead := c.buf[b]
        \t\tif lead < 0x80 {
        \t\t\tb++
        \t\t\trem--
        \t\t} else if lead >= 0xC2 && lead < 0xE0 && b+2 <= len(c.buf) {
        \t\t\tb += 2
        \t\t\trem--
        \t\t} else if lead >= 0xE0 && lead < 0xF0 && b+3 <= len(c.buf) {
        \t\t\tb += 3
        \t\t\trem--
        \t\t} else if lead >= 0xF0 && lead < 0xF5 && b+4 <= len(c.buf) {
        \t\t\tif rem < 2 {
        \t\t\t\treturn -1
        \t\t\t}
        \t\t\tb += 4
        \t\t\trem -= 2
        \t\t} else {
        \t\t\treturn -1
        \t\t}
        \t}
        \treturn b
        }

        func (c *jStringBuilderCore) charAt(idx int32) uint16 {
        \tn := c.length()
        \tif idx < 0 || idx >= n {
        \t\tpanic(javalang.NewStringIndexOutOfBoundsExceptionForIndex(int(idx), int(n)))
        \t}
        \tif c.pendingHigh != 0 && idx == n-1 {
        \t\treturn c.pendingHigh
        \t}
        \teffN := n
        \tif c.pendingHigh != 0 {
        \t\teffN--
        \t}
        \tif int(effN) == len(c.buf) {
        \t\treturn uint16(c.buf[idx])
        \t}
        \tif idx == effN-1 && len(c.buf) > 0 && c.buf[len(c.buf)-1] < 0x80 {
        \t\treturn uint16(c.buf[len(c.buf)-1])
        \t}
        \toff := c.byteOffsetOfUTF16(idx, effN)
        \tif off >= 0 && off < len(c.buf) && c.buf[off] < 0x80 {
        \t\treturn uint16(c.buf[off])
        \t}
        \treturn javalang.StringCharAt(c.str(), int(idx))
        }

        func (c *jStringBuilderCore) insertChar(idx int32, ch uint16) {
        \tn := c.length()
        \tCheckPositionIndex(idx, n)
        \tif idx == n {
        \t\tc.appendChar(ch)
        \t\treturn
        \t}
        \tif ch < 0x80 {
        \t\teffN := n
        \t\tif c.pendingHigh != 0 {
        \t\t\teffN--
        \t\t}
        \t\tif idx <= effN {
        \t\t\toff := c.byteOffsetOfUTF16(idx, effN)
        \t\t\tif off >= 0 {
        \t\t\t\tc.invalidateCache()
        \t\t\t\tc.buf = append(c.buf, 0)
        \t\t\t\tcopy(c.buf[off+1:], c.buf[off:len(c.buf)-1])
        \t\t\t\tc.buf[off] = byte(ch)
        \t\t\t\tc.utf16Len = n + 1
        \t\t\t\treturn
        \t\t\t}
        \t\t}
        \t}
        \ts := c.str()
        \tc.resetTo(JStrSub(s, 0, idx))
        \tc.appendChar(ch)
        \tc.appendString(JStrSubFrom(s, idx))
        }

        func (c *jStringBuilderCore) insertString(idx int32, raw string) {
        \tsIns := StrValString(raw)
        \tn := c.length()
        \tCheckPositionIndex(idx, n)
        \tif len(sIns) == 0 {
        \t\treturn
        \t}
        \tif idx == n {
        \t\tc.appendString(sIns)
        \t\treturn
        \t}
        \tif len(sIns) < 64 {
        \t\tascii := true
        \t\tfor i := 0; i < len(sIns); i++ {
        \t\t\tif sIns[i] >= 0x80 {
        \t\t\t\tascii = false
        \t\t\t\tbreak
        \t\t\t}
        \t\t}
        \t\tif ascii {
        \t\t\teffN := n
        \t\t\tif c.pendingHigh != 0 {
        \t\t\t\teffN--
        \t\t\t}
        \t\t\tif idx <= effN {
        \t\t\t\toff := c.byteOffsetOfUTF16(idx, effN)
        \t\t\t\tif off >= 0 {
        \t\t\t\t\tc.invalidateCache()
        \t\t\t\t\tk := len(sIns)
        \t\t\t\t\toldLen := len(c.buf)
        \t\t\t\t\tc.buf = append(c.buf, sIns...)
        \t\t\t\t\tcopy(c.buf[off+k:], c.buf[off:oldLen])
        \t\t\t\t\tcopy(c.buf[off:], sIns)
        \t\t\t\t\tc.utf16Len = n + int32(k)
        \t\t\t\t\treturn
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t}
        \ts := c.str()
        \tc.resetTo(JStrSub(s, 0, idx))
        \tc.appendString(sIns)
        \tc.appendString(JStrSubFrom(s, idx))
        }

        func (c *jStringBuilderCore) setCharAt(idx int32, ch uint16) {
        \tn := c.length()
        \tCheckElementIndex(idx, n)
        \tif ch < 0x80 && c.pendingHigh == 0 && int(n) == len(c.buf) {
        \t\tc.invalidateCache()
        \t\tc.buf[idx] = byte(ch)
        \t\treturn
        \t}
        \ts := c.str()
        \tc.resetTo(JStrSub(s, 0, idx))
        \tc.appendChar(ch)
        \tc.appendString(JStrSubFrom(s, idx+1))
        }

        func (c *jStringBuilderCore) deleteCharAt(idx int32) {
        \tn := c.length()
        \tCheckElementIndex(idx, n)
        \tif c.pendingHigh != 0 && idx == n-1 {
        \t\tc.pendingHigh = 0
        \t\tc.utf16Len = n - 1
        \t\treturn
        \t}
        \tif c.pendingHigh == 0 && int(n) == len(c.buf) {
        \t\tc.invalidateCache()
        \t\tcopy(c.buf[idx:], c.buf[idx+1:])
        \t\tc.buf = c.buf[:len(c.buf)-1]
        \t\tc.utf16Len = n - 1
        \t\treturn
        \t}
        \ts := c.str()
        \tc.resetTo(JStrSub(s, 0, idx))
        \tc.appendString(JStrSubFrom(s, idx+1))
        }

        func (c *jStringBuilderCore) deleteRange(start, end int32) {
        \tn := c.length()
        \tif start < 0 || start > n || start > end {
        \t\tpanic(javalang.NewStringIndexOutOfBoundsExceptionForRange(int(start), int(end), int(n)))
        \t}
        \tif end > n {
        \t\tend = n
        \t}
        \tif start == end {
        \t\treturn
        \t}
        \tif start == 0 && end == n {
        \t\tc.reset()
        \t\treturn
        \t}
        \tif c.pendingHigh == 0 && int(n) == len(c.buf) {
        \t\tc.invalidateCache()
        \t\tcopy(c.buf[start:], c.buf[end:])
        \t\tc.buf = c.buf[:len(c.buf)-int(end-start)]
        \t\tc.utf16Len = n - (end - start)
        \t\treturn
        \t}
        \ts := c.str()
        \tc.resetTo(JStrSub(s, 0, start))
        \tc.appendString(JStrSubFrom(s, end))
        }

        func (c *jStringBuilderCore) replaceRange(start, end int32, raw string) {
        \trpl := CheckNotNullString(raw)
        \tn := c.length()
        \tif start < 0 || start > n || start > end {
        \t\tpanic(javalang.NewStringIndexOutOfBoundsExceptionForRange(int(start), int(end), int(n)))
        \t}
        \tif end > n {
        \t\tend = n
        \t}
        \ts := c.str()
        \tc.resetTo(JStrSub(s, 0, start))
        \tc.appendString(rpl)
        \tc.appendString(JStrSubFrom(s, end))
        }

        func (c *jStringBuilderCore) reverse() {
        \tn := c.length()
        \tif n <= 1 {
        \t\treturn
        \t}
        \tif c.pendingHigh == 0 && int(n) == len(c.buf) {
        \t\tc.invalidateCache()
        \t\tfor i, j := 0, len(c.buf)-1; i < j; i, j = i+1, j-1 {
        \t\t\tc.buf[i], c.buf[j] = c.buf[j], c.buf[i]
        \t\t}
        \t\treturn
        \t}
        \ts := c.str()
        \tu := make([]uint16, n)
        \tfor i := int32(0); i < n; i++ {
        \t\tu[i] = javalang.StringCharAt(s, int(i))
        \t}
        \t// Reverse UTF-16 code units while preserving valid surrogate pairs per Java StringBuilder.reverse()
        \thasSurrogates := false
        \tfor i, j := int32(0), n-1; i < j; i, j = i+1, j-1 {
        \t\tci, cj := u[i], u[j]
        \t\tif (ci >= 0xD800 && ci <= 0xDFFF) || (cj >= 0xD800 && cj <= 0xDFFF) {
        \t\t\thasSurrogates = true
        \t\t}
        \t\tu[i], u[j] = cj, ci
        \t}
        \tif len(u) == 1 && u[0] >= 0xD800 && u[0] <= 0xDFFF {
        \t\thasSurrogates = true
        \t}
        \tif hasSurrogates {
        \t\tfor i := int32(0); i < n-1; i++ {
        \t\t\tc2 := u[i]
        \t\t\tif c2 >= 0xDC00 && c2 <= 0xDFFF {
        \t\t\t\tc1 := u[i+1]
        \t\t\t\tif c1 >= 0xD800 && c1 <= 0xDBFF {
        \t\t\t\t\tu[i] = c1
        \t\t\t\t\tu[i+1] = c2
        \t\t\t\t\ti++
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t}
        \tc.resetTo(javalang.StringFromUTF16(u))
        }

        func (c *jStringBuilderCore) setLength(newLen int32) {
        \tif newLen < 0 {
        \t\tpanic(javalang.NewStringIndexOutOfBoundsExceptionForIndex(int(newLen), 0))
        \t}
        \tif newLen == 0 {
        \t\tc.reset()
        \t\treturn
        \t}
        \tcur := c.length()
        \tif newLen == cur {
        \t\treturn
        \t}
        \tif newLen > cur {
        \t\tc.invalidateCache()
        \t\tif c.pendingHigh != 0 {
        \t\t\tc.flushPendingHigh()
        \t\t}
        \t\tadd := int(newLen - cur)
        \t\tc.buf = append(c.buf, make([]byte, add)...)
        \t\tc.utf16Len = newLen
        \t\treturn
        \t}
        \tif c.pendingHigh == 0 && int(cur) == len(c.buf) {
        \t\tc.invalidateCache()
        \t\tc.buf = c.buf[:newLen]
        \t\tc.utf16Len = newLen
        \t\treturn
        \t}
        \ts := c.str()
        \tc.resetTo(JStrSub(s, 0, newLen))
        }
        func StringFromInt8Slice(b []int8) string {
        \tCheckNotNull(b)
        \traw := make([]byte, len(b))
        \tfor i, v := range b {
        \t\traw[i] = byte(v)
        \t}
        \treturn string(raw)
        }
        func StringFromInt8Range(b []int8, off, length int32) string {
        \tCheckNotNull(b)
        \tif off < 0 || length < 0 || int(off)+int(length) > len(b) {
        \t\tpanic(javalang.NewStringIndexOutOfBoundsExceptionForRange(int(off), int(off+length), len(b)))
        \t}
        \traw := make([]byte, int(length))
        \tfor i := 0; i < int(length); i++ {
        \t\traw[i] = byte(b[int(off)+i])
        \t}
        \treturn string(raw)
        }
        func Arraycopy(src any, srcPos int32, dest any, destPos int32, length int32) {
        \tCheckNotNull(src)
        \tCheckNotNull(dest)
        \tif srcPos < 0 || destPos < 0 || length < 0 {
        \t\tpanic(javalang.NewArrayIndexOutOfBoundsException("arraycopy negative index"))
        \t}
        \tsp, dp, n := int(srcPos), int(destPos), int(length)
        \tvs := reflect.ValueOf(src)
        \tvd := reflect.ValueOf(dest)
        \tif sp+n > vs.Len() || dp+n > vd.Len() {
        \t\tpanic(javalang.NewArrayIndexOutOfBoundsException("arraycopy out of bounds"))
        \t}
        \tif n == 0 {
        \t\treturn
        \t}
        \tif vs.Type() == vd.Type() {
        \t\treflect.Copy(vd.Slice(dp, dp+n), vs.Slice(sp, sp+n))
        \t\treturn
        \t}
        \tfor i := 0; i < n; i++ {
        \t\tvd.Index(dp + i).Set(vs.Index(sp + i))
        \t}
        }
        func ClassPackageName(c *javalang.Class) string {
        \tCheckNotNull(c)
        \tn := c.GetName()
        \tif idx := strings.LastIndexByte(n, '.'); idx >= 0 {
        \t\treturn n[:idx]
        \t}
        \treturn ""
        }
        func ClassResourceExists(c *javalang.Class, path string) bool {
        \tCheckNotNull(c); CheckNotNull(path)
        \treturn shimsres.ResourceExists(ClassPackageName(c), path)
        }
        func ClassResourceBytes(c *javalang.Class, path string) ([]byte, bool) {
        \tCheckNotNull(c); CheckNotNull(path)
        \treturn shimsres.GetResourceBytes(ClassPackageName(c), path)
        }
        func ClassForName(name string) *javalang.Class {
        \tCheckNotNull(name)
        \tsimple := name
        \tif idx := strings.LastIndexAny(name, ".$"); idx >= 0 {
        \t\tsimple = name[idx+1:]
        \t}
        \treturn ClassLiteral(name, simple)
        }

        func CharSeqLength(cs any) int32 {
        \tif cs == nil || javalang.IsNil(cs) {
        \t\tpanic(javalang.NewNullPointerException())
        \t}
        \tif s, ok := cs.(string); ok {
        \t\tif s == NullString { panic(javalang.NewNullPointerException()) }
        \t\treturn int32(javalang.StringLength(s))
        \t}
        \tif l, ok := cs.(interface{ Length() int32 }); ok {
        \t\treturn l.Length()
        \t}
        \treturn int32(javalang.StringLength(StrVal(cs)))
        }

        func CharSeqCharAt(cs any, idx int32) uint16 {
        \tif cs == nil || javalang.IsNil(cs) {
        \t\tpanic(javalang.NewNullPointerException())
        \t}
        \tif s, ok := cs.(string); ok {
        \t\tif s == NullString { panic(javalang.NewNullPointerException()) }
        \t\treturn javalang.StringCharAt(s, int(idx))
        \t}
        \tif c, ok := cs.(interface{ CharAt(int32) uint16 }); ok {
        \t\treturn c.CharAt(idx)
        \t}
        \treturn javalang.StringCharAt(StrVal(cs), int(idx))
        }

        func CompareValues(a, b any) int32 {
        \ta = CheckNotNull(ToAny(a))
        \tb = CheckNotNull(ToAny(b))
        \tif sa, ok := a.(string); ok {
        \t\treturn int32(javalang.StringCompareTo(sa, b.(string)))
        \t}
        \tif ea, ok := a.(EnumAccessor); ok {
        \t\treturn ea.GetEnumBase().Ordinal_ - b.(EnumAccessor).GetEnumBase().Ordinal_
        \t}
        \tif ia, ok := a.(*BoxInt); ok {
        \t\tib := b.(*BoxInt)
        \t\tif ia.Val < ib.Val { return -1 }
        \t\tif ia.Val > ib.Val { return 1 }
        \t\treturn 0
        \t}
        \tif la, ok := a.(*BoxLong); ok {
        \t\tlb := b.(*BoxLong)
        \t\tif la.Val < lb.Val { return -1 }
        \t\tif la.Val > lb.Val { return 1 }
        \t\treturn 0
        \t}
        \tif da, ok := a.(*BoxDouble); ok {
        \t\tdb := b.(*BoxDouble)
        \t\treturn int32(javalang.DoubleCompare(da.Val, db.Val))
        \t}
        \tif fa, ok := a.(*BoxFloat); ok {
        \t\tfb := b.(*BoxFloat)
        \t\treturn int32(javalang.DoubleCompare(float64(fa.Val), float64(fb.Val)))
        \t}
        \tif ba, ok := a.(*BoxBoolean); ok {
        \t\tbb := b.(*BoxBoolean)
        \t\tif ba.Val == bb.Val { return 0 }
        \t\tif ba.Val { return 1 }
        \t\treturn -1
        \t}
        \tif cha, ok := a.(*BoxChar); ok {
        \t\tchb := b.(*BoxChar)
        \t\treturn int32(cha.Val) - int32(chb.Val)
        \t}
        \tif ca, ok := a.(interface{ CompareTo(any) int32 }); ok {
        \t\treturn ca.CompareTo(b)
        \t}
        \tpanic(javalang.NewClassCastException(fmt.Sprintf("%T is not Comparable", a)))
        }

        // --- Math & Numeric Intrinsics ---

        func ParseInt(s string, radix int32) int32 {
        \tCheckNotNull(s)
        \tv, err := strconv.ParseInt(strings.TrimPrefix(s, "+"), int(radix), 32)
        \tif err != nil {
        \t\tpanic(javalang.NewNumberFormatException(s))
        \t}
        \treturn int32(v)
        }

        func ParseLong(s string, radix int32) int64 {
        \tCheckNotNull(s)
        \tv, err := strconv.ParseInt(strings.TrimPrefix(s, "+"), int(radix), 64)
        \tif err != nil {
        \t\tpanic(javalang.NewNumberFormatException(s))
        \t}
        \treturn v
        }

        func ParseDouble(s string) float64 {
        \tCheckNotNullString(s)
        \treturn javalang.DoubleParseDouble(s)
        }

        func IntBitCount(x int32) int32       { return int32(bits.OnesCount32(uint32(x))) }
        func IntLeadingZeros(x int32) int32   { return int32(bits.LeadingZeros32(uint32(x))) }
        func IntTrailingZeros(x int32) int32  { return int32(bits.TrailingZeros32(uint32(x))) }
        func IntLowestOneBit(x int32) int32   { return x & (-x) }
        func IntHighestOneBit(x int32) int32 {
        \tif x == 0 {
        \t\treturn 0
        \t}
        \treturn int32(uint32(1) << (31 - bits.LeadingZeros32(uint32(x))))
        }
        func LongBitCount(x int64) int32      { return int32(bits.OnesCount64(uint64(x))) }
        func LongLeadingZeros(x int64) int32  { return int32(bits.LeadingZeros64(uint64(x))) }
        func LongTrailingZeros(x int64) int32 { return int32(bits.TrailingZeros64(uint64(x))) }
        func LongLowestOneBit(x int64) int64  { return x & (-x) }
        func LongHighestOneBit(x int64) int64 {
        \tif x == 0 {
        \t\treturn 0
        \t}
        \treturn int64(uint64(1) << (63 - bits.LeadingZeros64(uint64(x))))
        }

        func StringHash(s string) int32 {
        \tif s == NullString {
        \t\treturn 0
        \t}
        \treturn javalang.StringHashCode(s)
        }

        func BoolHash(b bool) int32 {
        \tif b {
        \t\treturn 1231
        \t}
        \treturn 1237
        }

        func LongHash(x int64) int32 {
        \tu := uint64(x)
        \treturn int32(u ^ (u >> 32))
        }

        type IdentityHasher interface {
        \tIdentityHash_() int32
        }

        var nextIdentHash_ atomic.Uint32

        //go:noinline
        func AssignIdentHash(p *atomic.Int32) int32 {
        \tfor {
        \t\th := int32(nextIdentHash_.Add(0x9e3779b9))
        \t\tif h == 0 {
        \t\t\tcontinue
        \t\t}
        \t\tif p.CompareAndSwap(0, h) {
        \t\t\treturn h
        \t\t}
        \t\treturn p.Load()
        \t}
        }

        func IdentityHash(v any) int32 {
        \tif v == nil {
        \t\treturn 0
        \t}
        \tif ea, ok := v.(EnumAccessor); ok {
        \t\teb := ea.GetEnumBase()
        \t\tif eb == nil {
        \t\t\treturn 0
        \t\t}
        \t\treturn eb.HashCode()
        \t}
        \tif nc, ok := v.(jJavaNilChecker); ok && nc.IsJavaNil_() {
        \t\treturn 0
        \t}
        \tif ih, ok := v.(IdentityHasher); ok {
        \t\treturn ih.IdentityHash_()
        \t}
        \treturn javalang.IdentityHashCode(v)
        }

        func ObjectHash(v any) int32 {
        \tif v == nil {
        \t\treturn 0
        \t}
        \tif fk, ok := v.(FastCustomKey); ok {
        \t\tif fk.IsJavaNil_() {
        \t\t\treturn 0
        \t\t}
        \t\treturn fk.HashCode()
        \t}
        \tif _, hasHC := v.(jHashCoder); !hasHC {
        \t\tif ih, ok := v.(IdentityHasher); ok {
        \t\t\tif nc, ok2 := v.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
        \t\t\t\treturn 0
        \t\t\t}
        \t\t\treturn ih.IdentityHash_()
        \t\t}
        \t}
        \tif isJavaNilAny(v) {
        \t\treturn 0
        \t}
        \treturn javalang.ObjectHashCode(v)
        }

        func ObjectsHash(vals ...any) int32 {
        \tif vals == nil {
        \t\treturn 0
        \t}
        \th := int32(1)
        \tfor _, v := range vals {
        \t\th = 31*h + ObjectHash(ToAny(v))
        \t}
        \treturn h
        }
        """;
  }

  private String buildCoreRuntimeGoPart2ListDeque() {
    return """

        // --- Deterministic Slim Zero-Bloat Collections Core ---

        type JListCore struct {
        \tItems   []any
        \tcapHint int32
        }

        func NewJListCore(items ...any) *JListCore {
        \tc := &JListCore{}
        \tn := len(items)
        \tif n == 0 {
        \t\treturn c
        \t}
        \ts := make([]any, n)
        \tfor i, v := range items {
        \t\ts[i] = ToAny(v)
        \t}
        \tc.Items = s
        \treturn c
        }

        func (l *JListCore) InitFromSlice(items []any) {
        \tn := len(items)
        \tif n == 0 {
        \t\tl.Items = nil
        \t\treturn
        \t}
        \tout := make([]any, n)
        \tcopy(out, items)
        \tl.Items = out
        }

        func (l *JListCore) InitOwnSlice(items []any) {
        \tif len(items) == 0 {
        \t\tl.Items = nil
        \t\treturn
        \t}
        \tl.Items = items
        }

        func (l *JListCore) ToArray() []any {
        \tif len(l.Items) == 0 {
        \t\treturn []any{}
        \t}
        \treturn CloneSlice(l.Items)
        }

        func (l *JListCore) InitCap(initialCap int32) {
        \tif initialCap < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException("Illegal Capacity: " + fmt.Sprint(initialCap)))
        \t}
        \tif initialCap <= 0 {
        \t\treturn
        \t}
        \tif initialCap > 1024 {
        \t\tinitialCap = 1024
        \t}
        \tif l.Items == nil {
        \t\tif initialCap > l.capHint {
        \t\t\tl.capHint = initialCap
        \t\t}
        \t\treturn
        \t}
        \tif int(initialCap) > cap(l.Items) {
        \t\tnext := make([]any, len(l.Items), initialCap)
        \t\tcopy(next, l.Items)
        \t\tl.Items = next
        \t}
        }

        func (l *JListCore) Size() int32   { return int32(len(l.Items)) }
        func (l *JListCore) IsEmpty() bool { return len(l.Items) == 0 }
        func (l *JListCore) Get(i int32) any {
        \tif i < 0 || int(i) >= len(l.Items) {
        \t\tpanic(javalang.NewIndexOutOfBoundsExceptionForIndex(int(i), len(l.Items)))
        \t}
        \treturn l.Items[i]
        }
        func (l *JListCore) Set(i int32, v any) any {
        \told := l.Get(i)
        \tl.Items[i] = ToAny(v)
        \treturn old
        }
        func (l *JListCore) Add(v any) bool {
        \tv = ToAny(v)
        \tif cap(l.Items) == 0 {
        \t\tc := int(l.capHint)
        \t\tif c < 4 {
        \t\t\tc = 4
        \t\t}
        \t\tl.Items = make([]any, 1, c)
        \t\tl.Items[0] = v
        \t\treturn true
        \t}
        \tl.Items = append(l.Items, v)
        \treturn true
        }
        func (l *JListCore) AddAt(i int32, v any) {
        \tif i < 0 || int(i) > len(l.Items) {
        \t\tpanic(javalang.NewIndexOutOfBoundsExceptionForIndex(int(i), len(l.Items)))
        \t}
        \tv = ToAny(v)
        \tif cap(l.Items) == 0 {
        \t\tc := int(l.capHint)
        \t\tif c < 4 {
        \t\t\tc = 4
        \t\t}
        \t\tl.Items = make([]any, 1, c)
        \t\tl.Items[0] = v
        \t\treturn
        \t}
        \tl.Items = append(l.Items, nil)
        \tcopy(l.Items[int(i)+1:], l.Items[int(i):])
        \tl.Items[i] = v
        }
        func (l *JListCore) RemoveAt(i int32) any {
        \told := l.Get(i)
        \tn := len(l.Items)
        \tif n == 1 {
        \t\tl.Items[0] = nil
        \t\tl.Items = l.Items[:0]
        \t\treturn old
        \t}
        \tif i == 0 && n > 8 {
        \t\tl.Items[0] = nil
        \t\tl.Items = l.Items[1:]
        \t\tif cap(l.Items) > 64 && len(l.Items)*4 <= cap(l.Items) {
        \t\t\tnext := make([]any, len(l.Items), len(l.Items)*2)
        \t\t\tcopy(next, l.Items)
        \t\t\tl.Items = next
        \t\t}
        \t\treturn old
        \t}
        \tlast := n - 1
        \tcopy(l.Items[int(i):], l.Items[int(i)+1:])
        \tl.Items[last] = nil
        \tl.Items = l.Items[:last]
        \treturn old
        }
        func (l *JListCore) RemoveVal(v any) bool {
        \tv = ToAny(v)
        \tfor i, x := range l.Items {
        \t\tif javalang.ObjectsEquals(x, v) {
        \t\t\tl.RemoveAt(int32(i))
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }
        func (l *JListCore) Contains(v any) bool {
        \treturn l.IndexOf(v) >= 0
        }
        func (l *JListCore) IndexOf(v any) int32 {
        \tv = ToAny(v)
        \tfor i, x := range l.Items {
        \t\tif javalang.ObjectsEquals(x, v) {
        \t\t\treturn int32(i)
        \t\t}
        \t}
        \treturn -1
        }
        func (l *JListCore) LastIndexOf(v any) int32 {
        \tv = ToAny(v)
        \tfor i := len(l.Items) - 1; i >= 0; i-- {
        \t\tif javalang.ObjectsEquals(l.Items[i], v) {
        \t\t\treturn int32(i)
        \t\t}
        \t}
        \treturn -1
        }
        func (l *JListCore) Clear() {
        \tif l.Items != nil {
        \t\tclear(l.Items)
        \t\tl.Items = l.Items[:0]
        \t}
        }

        // JDequeCore is an O(1) power-of-2 circular ring buffer backing java.util.ArrayDeque,
        // with a homogeneous noscan `nodes []uint32` fast mode for AST Node worklists.
        type JDequeCore struct {
        \tbuf      []any
        \tnodes    []uint32
        \tcolorIds []int64
        \thead     int32
        \tmask     int32 // int32(len(buf)-1), int32(len(nodes)-1), or int32(len(colorIds)-1)
        \tsize     int32
        \tcapHint  int32
        }

        func checkDequeElem(v any) any {
        \tif v == nil {
        \t\tThrowNPE()
        \t}
        \tif nc, ok := v.(jJavaNilChecker); ok {
        \t\tif nc.IsJavaNil_() {
        \t\t\tThrowNPE()
        \t\t}
        \t\treturn v
        \t}
        \tif s, ok := v.(string); ok {
        \t\tif s == NullString {
        \t\t\tThrowNPE()
        \t\t}
        \t\treturn v
        \t}
        \tv = ToAny(v)
        \tif v == nil {
        \t\tThrowNPE()
        \t}
        \treturn v
        }

        func extractDequeNodeId(v any) (uint32, bool) {
        \tif rn, ok := v.(jRawNodeProvider); ok {
        \t\tif isJavaNilAny(v) {
        \t\t\treturn 0, false
        \t\t}
        \t\tu := rn.RawNodeId_()
        \t\treturn u, u != 0
        \t}
        \tif n, ok := v.(nodestore.Node); ok {
        \t\treturn uint32(n), n != 0
        \t}
        \tif u, ok := v.(uint32); ok {
        \t\treturn u, u != 0
        \t}
        \treturn 0, false
        }

        func (d *JDequeCore) InitCap(initialCap int32) {
        \tif initialCap <= 0 {
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
        \tif d.buf == nil && d.nodes == nil && d.colorIds == nil {
        \t\tif c > d.capHint {
        \t\t\td.capHint = c
        \t\t}
        \t}
        }

        func (d *JDequeCore) growColorIds() {
        \tif d.colorIds == nil {
        \t\tc := d.capHint
        \t\tif c < 8 {
        \t\t\tc = 8
        \t\t}
        \t\td.colorIds = make([]int64, c)
        \t\td.mask = c - 1
        \t\treturn
        \t}
        \toldCap := int32(len(d.colorIds))
        \tnewCap := oldCap << 1
        \tnb := make([]int64, newCap)
        \tfirstPart := oldCap - d.head
        \tcopy(nb[:firstPart], d.colorIds[d.head:])
        \tif d.head > 0 {
        \t\tcopy(nb[firstPart:], d.colorIds[:d.head])
        \t}
        \td.colorIds = nb
        \td.head = 0
        \td.mask = newCap - 1
        }

        func (d *JDequeCore) upgradeColorIdsToBuf() {
        \tif d.size == 0 {
        \t\tif int32(len(d.colorIds)) > d.capHint {
        \t\t\td.capHint = int32(len(d.colorIds))
        \t\t}
        \t\td.colorIds = nil
        \t\td.head = 0
        \t\td.mask = 0
        \t\treturn
        \t}
        \tcapLen := int32(len(d.colorIds))
        \tnb := make([]any, capLen)
        \tfor i := int32(0); i < d.size; i++ {
        \t\tnb[i] = boxColorId(d.colorIds[(d.head+i)&d.mask])
        \t}
        \td.colorIds = nil
        \td.buf = nb
        \td.head = 0
        \td.mask = capLen - 1
        }

        func (d *JDequeCore) grow() {
        \tif d.buf == nil {
        \t\tc := d.capHint
        \t\tif c < 8 {
        \t\t\tc = 8
        \t\t}
        \t\td.buf = make([]any, c)
        \t\td.mask = c - 1
        \t\treturn
        \t}
        \toldCap := int32(len(d.buf))
        \tnewCap := oldCap << 1
        \tnb := make([]any, newCap)
        \tfirstPart := oldCap - d.head
        \tcopy(nb[:firstPart], d.buf[d.head:])
        \tif d.head > 0 {
        \t\tcopy(nb[firstPart:], d.buf[:d.head])
        \t}
        \td.buf = nb
        \td.head = 0
        \td.mask = newCap - 1
        }

        func (d *JDequeCore) growNodes() {
        \tif d.nodes == nil {
        \t\tc := d.capHint
        \t\tif c < 8 {
        \t\t\tc = 8
        \t\t}
        \t\td.nodes = make([]uint32, c)
        \t\td.mask = c - 1
        \t\treturn
        \t}
        \toldCap := int32(len(d.nodes))
        \tnewCap := oldCap << 1
        \tnb := make([]uint32, newCap)
        \tfirstPart := oldCap - d.head
        \tcopy(nb[:firstPart], d.nodes[d.head:])
        \tif d.head > 0 {
        \t\tcopy(nb[firstPart:], d.nodes[:d.head])
        \t}
        \td.nodes = nb
        \td.head = 0
        \td.mask = newCap - 1
        }

        func (d *JDequeCore) upgradeNodesToBuf() {
        \tif d.size == 0 {
        \t\tif int32(len(d.nodes)) > d.capHint {
        \t\t\td.capHint = int32(len(d.nodes))
        \t\t}
        \t\td.nodes = nil
        \t\td.head = 0
        \t\td.mask = 0
        \t\treturn
        \t}
        \tcapLen := int32(len(d.nodes))
        \tnb := make([]any, capLen)
        \tfor i := int32(0); i < d.size; i++ {
        \t\tnb[i] = boxNodeById(d.nodes[(d.head+i)&d.mask])
        \t}
        \td.nodes = nil
        \td.buf = nb
        \td.head = 0
        \td.mask = capLen - 1
        }

        func (d *JDequeCore) elemAt(idx int32) any {
        \tif d.nodes != nil {
        \t\treturn boxNodeById(d.nodes[idx])
        \t}
        \tif d.colorIds != nil {
        \t\treturn boxColorId(d.colorIds[idx])
        \t}
        \treturn d.buf[idx]
        }

        func (d *JDequeCore) Size() int32   { return d.size }
        func (d *JDequeCore) IsEmpty() bool { return d.size == 0 }

        func (d *JDequeCore) AddFirstColorId(raw int64) {
        \tif d.nodes != nil {
        \t\td.upgradeNodesToBuf()
        \t}
        \tif d.buf != nil {
        \t\tif int(d.size) == len(d.buf) {
        \t\t\td.grow()
        \t\t}
        \t\td.head = (d.head - 1) & d.mask
        \t\td.buf[d.head] = boxColorId(raw)
        \t\td.size++
        \t\treturn
        \t}
        \tif int(d.size) == len(d.colorIds) {
        \t\td.growColorIds()
        \t}
        \td.head = (d.head - 1) & d.mask
        \td.colorIds[d.head] = raw
        \td.size++
        }

        func (d *JDequeCore) AddLastColorId(raw int64) {
        \tif d.nodes != nil {
        \t\td.upgradeNodesToBuf()
        \t}
        \tif d.buf != nil {
        \t\tif int(d.size) == len(d.buf) {
        \t\t\td.grow()
        \t\t}
        \t\td.buf[(d.head+d.size)&d.mask] = boxColorId(raw)
        \t\td.size++
        \t\treturn
        \t}
        \tif int(d.size) == len(d.colorIds) {
        \t\td.growColorIds()
        \t}
        \td.colorIds[(d.head+d.size)&d.mask] = raw
        \td.size++
        }

        func (d *JDequeCore) RemoveFirstVoid() {
        \tif d.size == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \tif d.nodes != nil {
        \t\td.nodes[d.head] = 0
        \t\td.head = (d.head + 1) & d.mask
        \t\td.size--
        \t\treturn
        \t}
        \tif d.colorIds != nil {
        \t\td.head = (d.head + 1) & d.mask
        \t\td.size--
        \t\treturn
        \t}
        \td.buf[d.head] = nil
        \td.head = (d.head + 1) & d.mask
        \td.size--
        }

        func (d *JDequeCore) RemoveLastVoid() {
        \tif d.size == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \tidx := (d.head + d.size - 1) & d.mask
        \tif d.nodes != nil {
        \t\td.nodes[idx] = 0
        \t\td.size--
        \t\treturn
        \t}
        \tif d.colorIds != nil {
        \t\td.size--
        \t\treturn
        \t}
        \td.buf[idx] = nil
        \td.size--
        }

        func (d *JDequeCore) AddFirstNode(n uint32) {
        \tif n == 0 {
        \t\tThrowNPE()
        \t}
        \tif d.colorIds != nil {
        \t\td.upgradeColorIdsToBuf()
        \t}
        \tif d.buf != nil {
        \t\tif int(d.size) == len(d.buf) {
        \t\t\td.grow()
        \t\t}
        \t\td.head = (d.head - 1) & d.mask
        \t\td.buf[d.head] = boxNodeById(n)
        \t\td.size++
        \t\treturn
        \t}
        \tif int(d.size) == len(d.nodes) {
        \t\td.growNodes()
        \t}
        \td.head = (d.head - 1) & d.mask
        \td.nodes[d.head] = n
        \td.size++
        }

        func (d *JDequeCore) AddLastNode(n uint32) {
        \tif n == 0 {
        \t\tThrowNPE()
        \t}
        \tif d.colorIds != nil {
        \t\td.upgradeColorIdsToBuf()
        \t}
        \tif d.buf != nil {
        \t\tif int(d.size) == len(d.buf) {
        \t\t\td.grow()
        \t\t}
        \t\td.buf[(d.head+d.size)&d.mask] = boxNodeById(n)
        \t\td.size++
        \t\treturn
        \t}
        \tif int(d.size) == len(d.nodes) {
        \t\td.growNodes()
        \t}
        \td.nodes[(d.head+d.size)&d.mask] = n
        \td.size++
        }

        func (d *JDequeCore) PopNode() uint32 {
        \tif d.size == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \tif d.nodes != nil {
        \t\tu := d.nodes[d.head]
        \t\td.nodes[d.head] = 0
        \t\td.head = (d.head + 1) & d.mask
        \t\td.size--
        \t\treturn u
        \t}
        \treturn CastNode[uint32](d.RemoveFirst())
        }

        func (d *JDequeCore) RemoveLastNode() uint32 {
        \tif d.size == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \tif d.nodes != nil {
        \t\tidx := (d.head + d.size - 1) & d.mask
        \t\tu := d.nodes[idx]
        \t\td.nodes[idx] = 0
        \t\td.size--
        \t\treturn u
        \t}
        \treturn CastNode[uint32](d.RemoveLast())
        }

        func (d *JDequeCore) PollFirstNode() uint32 {
        \tif d.size == 0 {
        \t\treturn 0
        \t}
        \tif d.nodes != nil {
        \t\tu := d.nodes[d.head]
        \t\td.nodes[d.head] = 0
        \t\td.head = (d.head + 1) & d.mask
        \t\td.size--
        \t\treturn u
        \t}
        \treturn CastNode[uint32](d.PollFirst())
        }

        func (d *JDequeCore) PeekFirstNode() uint32 {
        \tif d.size == 0 {
        \t\treturn 0
        \t}
        \tif d.nodes != nil {
        \t\treturn d.nodes[d.head]
        \t}
        \treturn CastNode[uint32](d.PeekFirst())
        }

        func (d *JDequeCore) PeekLastNode() uint32 {
        \tif d.size == 0 {
        \t\treturn 0
        \t}
        \tif d.nodes != nil {
        \t\treturn d.nodes[(d.head+d.size-1)&d.mask]
        \t}
        \treturn CastNode[uint32](d.PeekLast())
        }

        func (d *JDequeCore) AddLast(v any) {
        \tv = checkDequeElem(v)
        \tif d.buf == nil {
        \t\tif d.colorIds == nil {
        \t\t\tif u, ok := extractDequeNodeId(v); ok {
        \t\t\t\td.AddLastNode(u)
        \t\t\t\treturn
        \t\t\t}
        \t\t}
        \t\tif d.nodes == nil {
        \t\t\tif rc, ok := v.(jRawColorIdProvider); ok {
        \t\t\t\tif raw, valid := rc.RawColorId_(); valid {
        \t\t\t\t\td.AddLastColorId(raw)
        \t\t\t\t\treturn
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t\tif d.nodes != nil {
        \t\t\td.upgradeNodesToBuf()
        \t\t}
        \t\tif d.colorIds != nil {
        \t\t\td.upgradeColorIdsToBuf()
        \t\t}
        \t}
        \tif int(d.size) == len(d.buf) {
        \t\td.grow()
        \t}
        \td.buf[(d.head+d.size)&d.mask] = v
        \td.size++
        }

        func (d *JDequeCore) OfferLast(v any) bool {
        \td.AddLast(v)
        \treturn true
        }

        func (d *JDequeCore) Add(v any) bool {
        \td.AddLast(v)
        \treturn true
        }

        func (d *JDequeCore) AddFirst(v any) {
        \tv = checkDequeElem(v)
        \tif d.buf == nil {
        \t\tif d.colorIds == nil {
        \t\t\tif u, ok := extractDequeNodeId(v); ok {
        \t\t\t\td.AddFirstNode(u)
        \t\t\t\treturn
        \t\t\t}
        \t\t}
        \t\tif d.nodes == nil {
        \t\t\tif rc, ok := v.(jRawColorIdProvider); ok {
        \t\t\t\tif raw, valid := rc.RawColorId_(); valid {
        \t\t\t\t\td.AddFirstColorId(raw)
        \t\t\t\t\treturn
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t\tif d.nodes != nil {
        \t\t\td.upgradeNodesToBuf()
        \t\t}
        \t\tif d.colorIds != nil {
        \t\t\td.upgradeColorIdsToBuf()
        \t\t}
        \t}
        \tif int(d.size) == len(d.buf) {
        \t\td.grow()
        \t}
        \td.head = (d.head - 1) & d.mask
        \td.buf[d.head] = v
        \td.size++
        }

        func (d *JDequeCore) OfferFirst(v any) bool {
        \td.AddFirst(v)
        \treturn true
        }

        func (d *JDequeCore) Push(v any) {
        \td.AddFirst(v)
        }

        func (d *JDequeCore) PollFirst() any {
        \tif d.size == 0 {
        \t\treturn nil
        \t}
        \tif d.nodes != nil {
        \t\tu := d.nodes[d.head]
        \t\td.nodes[d.head] = 0
        \t\td.head = (d.head + 1) & d.mask
        \t\td.size--
        \t\treturn boxNodeById(u)
        \t}
        \tif d.colorIds != nil {
        \t\traw := d.colorIds[d.head]
        \t\td.head = (d.head + 1) & d.mask
        \t\td.size--
        \t\treturn boxColorId(raw)
        \t}
        \tv := d.buf[d.head]
        \td.buf[d.head] = nil
        \td.head = (d.head + 1) & d.mask
        \td.size--
        \treturn v
        }

        func (d *JDequeCore) RemoveFirst() any {
        \tif d.size == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \tif d.nodes != nil {
        \t\tu := d.nodes[d.head]
        \t\td.nodes[d.head] = 0
        \t\td.head = (d.head + 1) & d.mask
        \t\td.size--
        \t\treturn boxNodeById(u)
        \t}
        \tif d.colorIds != nil {
        \t\traw := d.colorIds[d.head]
        \t\td.head = (d.head + 1) & d.mask
        \t\td.size--
        \t\treturn boxColorId(raw)
        \t}
        \tv := d.buf[d.head]
        \td.buf[d.head] = nil
        \td.head = (d.head + 1) & d.mask
        \td.size--
        \treturn v
        }

        func (d *JDequeCore) Pop() any {
        \treturn d.RemoveFirst()
        }

        func (d *JDequeCore) PollLast() any {
        \tif d.size == 0 {
        \t\treturn nil
        \t}
        \tidx := (d.head + d.size - 1) & d.mask
        \tif d.nodes != nil {
        \t\tu := d.nodes[idx]
        \t\td.nodes[idx] = 0
        \t\td.size--
        \t\treturn boxNodeById(u)
        \t}
        \tif d.colorIds != nil {
        \t\traw := d.colorIds[idx]
        \t\td.size--
        \t\treturn boxColorId(raw)
        \t}
        \tv := d.buf[idx]
        \td.buf[idx] = nil
        \td.size--
        \treturn v
        }

        func (d *JDequeCore) RemoveLast() any {
        \tif d.size == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \tidx := (d.head + d.size - 1) & d.mask
        \tif d.nodes != nil {
        \t\tu := d.nodes[idx]
        \t\td.nodes[idx] = 0
        \t\td.size--
        \t\treturn boxNodeById(u)
        \t}
        \tif d.colorIds != nil {
        \t\traw := d.colorIds[idx]
        \t\td.size--
        \t\treturn boxColorId(raw)
        \t}
        \tv := d.buf[idx]
        \td.buf[idx] = nil
        \td.size--
        \treturn v
        }

        func (d *JDequeCore) PeekFirst() any {
        \tif d.size == 0 {
        \t\treturn nil
        \t}
        \treturn d.elemAt(d.head)
        }

        func (d *JDequeCore) GetFirst() any {
        \tif d.size == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \treturn d.elemAt(d.head)
        }

        func (d *JDequeCore) PeekLast() any {
        \tif d.size == 0 {
        \t\treturn nil
        \t}
        \treturn d.elemAt((d.head + d.size - 1) & d.mask)
        }

        func (d *JDequeCore) GetLast() any {
        \tif d.size == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty"))
        \t}
        \treturn d.elemAt((d.head + d.size - 1) & d.mask)
        }

        func (d *JDequeCore) Clear() {
        \tif d.nodes != nil || d.colorIds != nil {
        \t\td.head = 0
        \t\td.size = 0
        \t\treturn
        \t}
        \tfor i := int32(0); i < d.size; i++ {
        \t\td.buf[(d.head+i)&d.mask] = nil
        \t}
        \td.head = 0
        \td.size = 0
        }

        func (d *JDequeCore) ForEach(fn func(any)) {
        \tif d.nodes != nil {
        \t\tfor i := int32(0); i < d.size; i++ {
        \t\t\tfn(boxNodeById(d.nodes[(d.head+i)&d.mask]))
        \t\t}
        \t\treturn
        \t}
        \tif d.colorIds != nil {
        \t\tfor i := int32(0); i < d.size; i++ {
        \t\t\tfn(boxColorId(d.colorIds[(d.head+i)&d.mask]))
        \t\t}
        \t\treturn
        \t}
        \tfor i := int32(0); i < d.size; i++ {
        \t\tfn(d.buf[(d.head+i)&d.mask])
        \t}
        }

        func (d *JDequeCore) RemoveIf(pred func(any) bool) bool {
        \tif d.size == 0 {
        \t\treturn false
        \t}
        \tif d.colorIds != nil {
        \t\td.upgradeColorIdsToBuf()
        \t}
        \toldSize := d.size
        \twriteOff := int32(0)
        \tif d.nodes != nil {
        \t\tfor readOff := int32(0); readOff < oldSize; readOff++ {
        \t\t\tu := d.nodes[(d.head+readOff)&d.mask]
        \t\t\tif !pred(boxNodeById(u)) {
        \t\t\t\tif writeOff != readOff {
        \t\t\t\t\td.nodes[(d.head+writeOff)&d.mask] = u
        \t\t\t\t}
        \t\t\t\twriteOff++
        \t\t\t}
        \t\t}
        \t\tif writeOff == oldSize {
        \t\t\treturn false
        \t\t}
        \t\td.size = writeOff
        \t\treturn true
        \t}
        \tfor readOff := int32(0); readOff < oldSize; readOff++ {
        \t\tv := d.buf[(d.head+readOff)&d.mask]
        \t\tif !pred(v) {
        \t\t\tif writeOff != readOff {
        \t\t\t\td.buf[(d.head+writeOff)&d.mask] = v
        \t\t\t}
        \t\t\twriteOff++
        \t\t}
        \t}
        \tif writeOff == oldSize {
        \t\treturn false
        \t}
        \tfor i := writeOff; i < oldSize; i++ {
        \t\td.buf[(d.head+i)&d.mask] = nil
        \t}
        \td.size = writeOff
        \treturn true
        }

        func (d *JDequeCore) deleteAtOffset(off int32) {
        \tif d.colorIds != nil {
        \t\td.upgradeColorIdsToBuf()
        \t}
        \tif d.nodes != nil {
        \t\tif off < (d.size >> 1) {
        \t\t\tfor j := off; j > 0; j-- {
        \t\t\t\td.nodes[(d.head+j)&d.mask] = d.nodes[(d.head+j-1)&d.mask]
        \t\t\t}
        \t\t\td.nodes[d.head] = 0
        \t\t\td.head = (d.head + 1) & d.mask
        \t\t\td.size--
        \t\t\treturn
        \t\t}
        \t\tlast := d.size - 1
        \t\tfor j := off; j < last; j++ {
        \t\t\td.nodes[(d.head+j)&d.mask] = d.nodes[(d.head+j+1)&d.mask]
        \t\t}
        \t\td.nodes[(d.head+last)&d.mask] = 0
        \t\td.size--
        \t\treturn
        \t}
        \tif off < (d.size >> 1) {
        \t\tfor j := off; j > 0; j-- {
        \t\t\td.buf[(d.head+j)&d.mask] = d.buf[(d.head+j-1)&d.mask]
        \t\t}
        \t\td.buf[d.head] = nil
        \t\td.head = (d.head + 1) & d.mask
        \t\td.size--
        \t\treturn
        \t}
        \tlast := d.size - 1
        \tfor j := off; j < last; j++ {
        \t\td.buf[(d.head+j)&d.mask] = d.buf[(d.head+j+1)&d.mask]
        \t}
        \td.buf[(d.head+last)&d.mask] = nil
        \td.size--
        }

        func (d *JDequeCore) Contains(v any) bool {
        \tif v == nil || d.size == 0 {
        \t\treturn false
        \t}
        \tif d.colorIds != nil {
        \t\trc, ok := v.(jRawColorIdProvider)
        \t\tif !ok {
        \t\t\treturn false
        \t\t}
        \t\traw, valid := rc.RawColorId_()
        \t\tif !valid {
        \t\t\treturn false
        \t\t}
        \t\tfor i := int32(0); i < d.size; i++ {
        \t\t\tif d.colorIds[(d.head+i)&d.mask] == raw {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \t}
        \tif d.nodes != nil {
        \t\tu, ok := extractDequeNodeId(v)
        \t\tif !ok {
        \t\t\treturn false
        \t\t}
        \t\tfor i := int32(0); i < d.size; i++ {
        \t\t\tif d.nodes[(d.head+i)&d.mask] == u {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \t}
        \tv = ToAny(v)
        \tif v == nil {
        \t\treturn false
        \t}
        \tfor i := int32(0); i < d.size; i++ {
        \t\tif javalang.ObjectsEquals(d.buf[(d.head+i)&d.mask], v) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (d *JDequeCore) RemoveFirstOccurrence(v any) bool {
        \tif v == nil || d.size == 0 {
        \t\treturn false
        \t}
        \tif d.colorIds != nil {
        \t\td.upgradeColorIdsToBuf()
        \t}
        \tif d.nodes != nil {
        \t\tu, ok := extractDequeNodeId(v)
        \t\tif !ok {
        \t\t\treturn false
        \t\t}
        \t\tfor i := int32(0); i < d.size; i++ {
        \t\t\tif d.nodes[(d.head+i)&d.mask] == u {
        \t\t\t\td.deleteAtOffset(i)
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \t}
        \tv = ToAny(v)
        \tif v == nil {
        \t\treturn false
        \t}
        \tfor i := int32(0); i < d.size; i++ {
        \t\tif javalang.ObjectsEquals(d.buf[(d.head+i)&d.mask], v) {
        \t\t\td.deleteAtOffset(i)
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (d *JDequeCore) RemoveLastOccurrence(v any) bool {
        \tif v == nil || d.size == 0 {
        \t\treturn false
        \t}
        \tif d.colorIds != nil {
        \t\td.upgradeColorIdsToBuf()
        \t}
        \tif d.nodes != nil {
        \t\tu, ok := extractDequeNodeId(v)
        \t\tif !ok {
        \t\t\treturn false
        \t\t}
        \t\tfor i := d.size - 1; i >= 0; i-- {
        \t\t\tif d.nodes[(d.head+i)&d.mask] == u {
        \t\t\t\td.deleteAtOffset(i)
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \t}
        \tv = ToAny(v)
        \tif v == nil {
        \t\treturn false
        \t}
        \tfor i := d.size - 1; i >= 0; i-- {
        \t\tif javalang.ObjectsEquals(d.buf[(d.head+i)&d.mask], v) {
        \t\t\td.deleteAtOffset(i)
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (d *JDequeCore) CloneInto(dst *JDequeCore) {
        \tif d.size == 0 {
        \t\treturn
        \t}
        \tif d.colorIds != nil {
        \t\tc := int32(len(d.colorIds))
        \t\tdst.colorIds = make([]int64, c)
        \t\tdst.mask = c - 1
        \t\tfor i := int32(0); i < d.size; i++ {
        \t\t\tdst.colorIds[i] = d.colorIds[(d.head+i)&d.mask]
        \t\t}
        \t\tdst.size = d.size
        \t\treturn
        \t}
        \tif d.nodes != nil {
        \t\tc := int32(len(d.nodes))
        \t\tdst.nodes = make([]uint32, c)
        \t\tdst.mask = c - 1
        \t\tfor i := int32(0); i < d.size; i++ {
        \t\t\tdst.nodes[i] = d.nodes[(d.head+i)&d.mask]
        \t\t}
        \t\tdst.size = d.size
        \t\treturn
        \t}
        \tc := int32(len(d.buf))
        \tdst.buf = make([]any, c)
        \tdst.mask = c - 1
        \tfor i := int32(0); i < d.size; i++ {
        \t\tdst.buf[i] = d.buf[(d.head+i)&d.mask]
        \t}
        \tdst.size = d.size
        }

        func (d *JDequeCore) ToSlice() []any {
        \tif d.size == 0 {
        \t\treturn []any{}
        \t}
        \tout := make([]any, d.size)
        \tif d.colorIds != nil {
        \t\tfor i := int32(0); i < d.size; i++ {
        \t\t\tout[i] = boxColorId(d.colorIds[(d.head+i)&d.mask])
        \t\t}
        \t\treturn out
        \t}
        \tif d.nodes != nil {
        \t\tfor i := int32(0); i < d.size; i++ {
        \t\t\tout[i] = boxNodeById(d.nodes[(d.head+i)&d.mask])
        \t\t}
        \t\treturn out
        \t}
        \tcapLen := int32(len(d.buf))
        \tif d.head+d.size <= capLen {
        \t\tcopy(out, d.buf[d.head:d.head+d.size])
        \t} else {
        \t\tfirstPart := capLen - d.head
        \t\tcopy(out[:firstPart], d.buf[d.head:])
        \t\tcopy(out[firstPart:], d.buf[:d.size-firstPart])
        \t}
        \treturn out
        }
        """;
  }

  private String buildCoreRuntimeGoPart2Set() {
    return """

        const (
        \tjTombstoneStr         = "\\xff\\x00j2go_tomb"
        \tjTombstoneNode uint32 = 0xFFFFFFFF

        \tjSetModeEmpty uint8 = 0
        \tjSetModeNode  uint8 = 1
        \tjSetModeStr   uint8 = 2
        \tjSetModeAny   uint8 = 3

        \tjSmallLimit     int = 8
        \tjStrSmallLimit  int = 13
        \tjNodeSmallLimit int = 16
        )

        type jTombstoneType struct{}

        var (
        \tjTombstoneObj  any   = jTombstoneType{}
        \temptyAnySlice_ []any = []any{}
        )

        type jIdentityKey interface {
        \tHasIdentityEquals_()
        \tIsJavaNil_() bool
        }

        type FastCustomKey interface {
        \tHashCode() int32
        \tEquals(any) bool
        \tIsJavaNil_() bool
        }

        type jSourceObjPtr interface {
        \tBindSelf_(any)
        \tIsJavaNil_() bool
        }

        type jHashCoder interface {
        \tHashCode() int32
        }

        type jEqualer interface {
        \tEquals(any) bool
        }

        func smearHash(h int32) uint32 {
        \tu := uint32(h)
        \treturn u ^ (u >> 16)
        }

        func elemJavaHashCode(v any, identity bool) int32 {
        \tif v == nil {
        \t\treturn 0
        \t}
        \tif ea, ok := v.(EnumAccessor); ok {
        \t\teb := ea.GetEnumBase()
        \t\tif eb == nil {
        \t\t\treturn 0
        \t\t}
        \t\treturn eb.HashCode()
        \t}
        \tif identity {
        \t\tif nc, ok := v.(jJavaNilChecker); ok && nc.IsJavaNil_() {
        \t\t\treturn 0
        \t\t}
        \t\tif ih, ok := v.(IdentityHasher); ok {
        \t\t\treturn ih.IdentityHash_()
        \t\t}
        \t\treturn javalang.IdentityHashCode(v)
        \t}
        \tif ik, ok := v.(jIdentityKey); ok {
        \t\tif ik.IsJavaNil_() {
        \t\t\treturn 0
        \t\t}
        \t\tif ih, ok := v.(IdentityHasher); ok {
        \t\t\treturn ih.IdentityHash_()
        \t\t}
        \t\treturn javalang.IdentityHashCode(v)
        \t}
        \tif rn, ok := v.(jRawNodeProvider); ok {
        \t\treturn int32(rn.RawNodeId_())
        \t}
        \tif fk, ok := v.(FastCustomKey); ok {
        \t\tif fk.IsJavaNil_() {
        \t\t\treturn 0
        \t\t}
        \t\treturn fk.HashCode()
        \t}
        \tif s, ok := v.(string); ok {
        \t\tif s == NullString {
        \t\t\treturn 0
        \t\t}
        \t\treturn javalang.StringHashCode(s)
        \t}
        \tif nd, ok := v.(nodestore.Node); ok {
        \t\treturn int32(nd)
        \t}
        \tif _, hasHC := v.(jHashCoder); !hasHC {
        \t\tif ih, ok := v.(IdentityHasher); ok {
        \t\t\tif nc, ok2 := v.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
        \t\t\t\treturn 0
        \t\t\t}
        \t\t\treturn ih.IdentityHash_()
        \t\t}
        \t}
        \treturn javalang.ObjectHashCode(v)
        }

        func jSetElemEquals(a, b any) bool {
        \tif b == nil {
        \t\tif a == nil {
        \t\t\treturn true
        \t\t}
        \t\tif sa, ok := a.(string); ok {
        \t\t\treturn sa == NullString
        \t\t}
        \t\tif na, ok := a.(jJavaNilChecker); ok {
        \t\t\treturn na.IsJavaNil_()
        \t\t}
        \t\treturn isJavaNilAny(a)
        \t}
        \tif sb, ok := b.(string); ok {
        \t\tif sb == NullString {
        \t\t\treturn a == nil || isJavaNilAny(a)
        \t\t}
        \t\tsa, ok2 := a.(string)
        \t\treturn ok2 && sa == sb
        \t}
        \tif ib, ok := b.(jIdentityKey); ok {
        \t\tif ib.IsJavaNil_() {
        \t\t\treturn a == nil || isJavaNilAny(a)
        \t\t}
        \t\treturn a == b
        \t}
        \tif _, ok := b.(IdentityComparable); ok {
        \t\tif nb, ok2 := b.(jJavaNilChecker); ok2 && nb.IsJavaNil_() {
        \t\t\treturn a == nil || isJavaNilAny(a)
        \t\t}
        \t\treturn a == b
        \t}
        \tif _, ok := b.(EnumAccessor); ok {
        \t\tif nb, ok2 := b.(jJavaNilChecker); ok2 && nb.IsJavaNil_() {
        \t\t\treturn a == nil || isJavaNilAny(a)
        \t\t}
        \t\treturn a == b
        \t}
        \tif rnb, ok := b.(jRawNodeProvider); ok {
        \t\tub := rnb.RawNodeId_()
        \t\tif ub == 0 {
        \t\t\treturn a == nil || isJavaNilAny(a)
        \t\t}
        \t\tif rna, ok2 := a.(jRawNodeProvider); ok2 {
        \t\t\treturn rna.RawNodeId_() == ub
        \t\t}
        \t\tif nda, ok2 := a.(nodestore.Node); ok2 {
        \t\t\treturn uint32(nda) == ub
        \t\t}
        \t\treturn false
        \t}
        \tif ib, ok := b.(*BoxInt); ok {
        \t\tif ib == nil {
        \t\t\treturn a == nil || isJavaNilAny(a)
        \t\t}
        \t\tif ia, ok2 := a.(*BoxInt); ok2 {
        \t\t\treturn ia != nil && ia.Val == ib.Val
        \t\t}
        \t\tif iv, ok2 := a.(int32); ok2 {
        \t\t\treturn iv == ib.Val
        \t\t}
        \t\treturn false
        \t}
        \tif ivb, ok := b.(int32); ok {
        \t\tif iva, ok2 := a.(int32); ok2 {
        \t\t\treturn iva == ivb
        \t\t}
        \t\tif ia, ok2 := a.(*BoxInt); ok2 {
        \t\t\treturn ia != nil && ia.Val == ivb
        \t\t}
        \t\treturn false
        \t}
        \tif fkb, ok := b.(FastCustomKey); ok {
        \t\tif fkb.IsJavaNil_() {
        \t\t\treturn a == nil || isJavaNilAny(a)
        \t\t}
        \t\tif a == nil || isJavaNilAny(a) {
        \t\t\treturn false
        \t\t}
        \t\treturn a == b || fkb.Equals(a)
        \t}
        \tif a == nil {
        \t\treturn isJavaNilAny(b)
        \t}
        \tif ia, ok := a.(jIdentityKey); ok {
        \t\tif ia.IsJavaNil_() {
        \t\t\treturn isJavaNilAny(b)
        \t\t}
        \t\treturn a == b
        \t}
        \tif _, ok := a.(IdentityComparable); ok {
        \t\tif na, ok2 := a.(jJavaNilChecker); ok2 && na.IsJavaNil_() {
        \t\t\treturn isJavaNilAny(b)
        \t\t}
        \t\treturn a == b
        \t}
        \tif _, ok := a.(EnumAccessor); ok {
        \t\tif na, ok2 := a.(jJavaNilChecker); ok2 && na.IsJavaNil_() {
        \t\t\treturn isJavaNilAny(b)
        \t\t}
        \t\treturn a == b
        \t}
        \tif rna, ok := a.(jRawNodeProvider); ok {
        \t\tua := rna.RawNodeId_()
        \t\tif ua == 0 {
        \t\t\treturn isJavaNilAny(b)
        \t\t}
        \t\tif ndb, ok2 := b.(nodestore.Node); ok2 {
        \t\t\treturn ua == uint32(ndb)
        \t\t}
        \t\treturn false
        \t}
        \tif ia, ok := a.(*BoxInt); ok {
        \t\tif ia == nil {
        \t\t\treturn isJavaNilAny(b)
        \t\t}
        \t\tif iv, ok2 := b.(int32); ok2 {
        \t\t\treturn ia.Val == iv
        \t\t}
        \t\treturn false
        \t}
        \tif fka, ok := a.(FastCustomKey); ok {
        \t\tif fka.IsJavaNil_() {
        \t\t\treturn isJavaNilAny(b)
        \t\t}
        \t\tif isJavaNilAny(b) {
        \t\t\treturn false
        \t\t}
        \t\treturn a == b || fka.Equals(b)
        \t}
        \tif nda, ok := a.(nodestore.Node); ok {
        \t\tif nda == 0 {
        \t\t\treturn isJavaNilAny(b)
        \t\t}
        \t\tif ndb, ok2 := b.(nodestore.Node); ok2 {
        \t\t\treturn nda == ndb
        \t\t}
        \t\treturn false
        \t}
        \treturn javalang.ObjectsEquals(a, b)
        }

        func elemJavaEquals(a, b any, identity bool) bool {
        \tif identity {
        \t\treturn RefEquals(a, b)
        \t}
        \treturn jSetElemEquals(a, b)
        }

        func fastSetKey(k any, identity bool) (any, bool) {
        \tif k == nil {
        \t\treturn nil, true
        \t}
        \tif ik, ok := k.(jIdentityKey); ok {
        \t\tif ik.IsJavaNil_() {
        \t\t\treturn nil, true
        \t\t}
        \t\treturn k, true
        \t}
        \tif identity {
        \t\tif sp, ok := k.(jSourceObjPtr); ok {
        \t\t\tif sp.IsJavaNil_() {
        \t\t\t\treturn nil, true
        \t\t\t}
        \t\t\treturn k, true
        \t\t}
        \t\tif _, ok := k.(EnumAccessor); ok {
        \t\t\tif nc, ok2 := k.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
        \t\t\t\treturn nil, true
        \t\t\t}
        \t\t\treturn k, true
        \t\t}
        \t\tif rn, ok := k.(jRawNodeProvider); ok {
        \t\t\tu := rn.RawNodeId_()
        \t\t\tif u == 0 {
        \t\t\t\treturn nil, true
        \t\t\t}
        \t\t\treturn nodestore.Node(u), true
        \t\t}
        \t\tif nd, ok := k.(nodestore.Node); ok {
        \t\t\tif nd == 0 {
        \t\t\t\treturn nil, true
        \t\t\t}
        \t\t\treturn nd, true
        \t\t}
        \t\treturn nil, false
        \t}
        \tif _, ok := k.(IdentityComparable); ok {
        \t\tif nc, ok2 := k.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
        \t\t\treturn nil, true
        \t\t}
        \t\treturn k, true
        \t}
        \tif _, ok := k.(EnumAccessor); ok {
        \t\tif nc, ok2 := k.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
        \t\t\treturn nil, true
        \t\t}
        \t\treturn k, true
        \t}
        \tswitch x := k.(type) {
        \tcase string:
        \t\tif x == NullString {
        \t\t\treturn nil, true
        \t\t}
        \t\treturn x, true
        \tcase nodestore.Node:
        \t\tif x == 0 {
        \t\t\treturn nil, true
        \t\t}
        \t\treturn x, true
        \tcase *BoxInt:
        \t\tif x == nil {
        \t\t\treturn nil, true
        \t\t}
        \t\treturn x.Val, true
        \tcase int32:
        \t\treturn x, true
        \t}
        \tif rn, ok := k.(jRawNodeProvider); ok {
        \t\tu := rn.RawNodeId_()
        \t\tif u == 0 {
        \t\t\treturn nil, true
        \t\t}
        \t\treturn nodestore.Node(u), true
        \t}
        \tif sp, ok := k.(jSourceObjPtr); ok {
        \t\tif sp.IsJavaNil_() {
        \t\t\treturn nil, true
        \t\t}
        \t\tif _, hasEq := k.(javalang.Equaler); !hasEq {
        \t\t\tif _, hasHc := k.(jHashCoder); !hasHc {
        \t\t\t\treturn k, true
        \t\t\t}
        \t\t}
        \t}
        \treturn nil, false
        }

        func fastSetMapKey(k any) (any, bool) {
        \treturn fastSetKey(k, false)
        }

        // --- Dedicated Compact JSetCore (HashSet, LinkedHashSet, ImmutableSet.Builder, >8 ImmutableSet) ---

        type jSetStrState struct {
        \tstrs   []string
        \tstrIdx map[string]int32
        }

        type jSetStrState4 struct {
        \tstate jSetStrState
        \tbuf   [4]string
        }

        type jSetGenBuckets struct {
        \thashes   []int32
        \thashNext []int32
        \tbuckets  []int32
        \tgenLive  int32
        }

        type jSetAnyState struct {
        \titems    []any
        \tptrIdx   map[any]int32
        \tgen      *jSetGenBuckets
        \treads    uint16
        \tallIdent bool
        }

        func isNonNilIdentityElem(k any) bool {
        \tif k == nil {
        \t\treturn false
        \t}
        \tik, ok := k.(jIdentityKey)
        \treturn ok && !ik.IsJavaNil_()
        }

        func isFastIdentElem(v any) bool {
        \tif v == nil {
        \t\treturn false
        \t}
        \tif ik, ok := v.(jIdentityKey); ok {
        \t\treturn !ik.IsJavaNil_()
        \t}
        \tif _, ok := v.(IdentityComparable); ok {
        \t\tif nc, ok2 := v.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
        \t\t\treturn false
        \t\t}
        \t\treturn true
        \t}
        \tif _, ok := v.(EnumAccessor); ok {
        \t\tif nc, ok2 := v.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
        \t\t\treturn false
        \t\t}
        \t\treturn true
        \t}
        \tif cls, ok := v.(*javalang.Class); ok {
        \t\treturn cls != nil
        \t}
        \tif globalNodeBoxer == nil {
        \t\tif nd, ok := v.(nodestore.Node); ok {
        \t\t\treturn nd != 0
        \t\t}
        \t}
        \tif sp, ok := v.(jSourceObjPtr); ok {
        \t\tif sp.IsJavaNil_() {
        \t\t\treturn false
        \t\t}
        \t\tif _, hasEq := v.(javalang.Equaler); !hasEq {
        \t\t\tif _, hasHc := v.(jHashCoder); !hasHc {
        \t\t\t\treturn reflect.TypeOf(v).Comparable()
        \t\t\t}
        \t\t}
        \t}
        \treturn false
        }

        type jSetAnyState4 struct {
        \tstate jSetAnyState
        \tbuf   [4]any
        }

        type JSetCore struct {
        \tnodes    []uint32
        \text      any
        \tlive     int32
        \tholes    int32
        \tfirstIdx int32
        \tnodeBuf  [1]uint32
        \tcapHint  int16
        \tmode     uint8
        \tidentity bool
        \timm      bool
        }

        func NewJSetCore(identity bool) *JSetCore {
        \treturn &JSetCore{identity: identity}
        }

        func (s *JSetCore) Size() int32   { return s.live }
        func (s *JSetCore) IsEmpty() bool { return s.live == 0 }

        func (s *JSetCore) InitExpectedSize(expected int32) {
        \tif expected < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException("Illegal initial capacity: " + fmt.Sprint(expected)))
        \t}
        \tif expected > 4 && s.mode == jSetModeEmpty {
        \t\tif expected > 1024 {
        \t\t\ts.capHint = 1024
        \t\t} else {
        \t\t\ts.capHint = int16(expected)
        \t\t}
        \t}
        }

        func (s *JSetCore) initCap() int {
        \tif s.capHint > 4 {
        \t\treturn min(int(s.capHint), 256)
        \t}
        \treturn 4
        }

        func (s *JSetCore) compactNodes() {
        \tif s.holes == 0 {
        \t\ts.firstIdx = 0
        \t\treturn
        \t}
        \tvar idxMap map[uint32]int32
        \tif s.ext != nil {
        \t\tidxMap = s.ext.(map[uint32]int32)
        \t}
        \tw := 0
        \tfor i := int(s.firstIdx); i < len(s.nodes); i++ {
        \t\tu := s.nodes[i]
        \t\tif u != jTombstoneNode {
        \t\t\ts.nodes[w] = u
        \t\t\tif idxMap != nil {
        \t\t\t\tidxMap[u] = int32(w)
        \t\t\t}
        \t\t\tw++
        \t\t}
        \t}
        \ts.nodes = s.nodes[:w]
        \ts.holes = 0
        \ts.firstIdx = 0
        }

        func (s *JSetCore) compactStrs(st *jSetStrState) {
        \tif s.holes == 0 {
        \t\ts.firstIdx = 0
        \t\treturn
        \t}
        \tw := 0
        \tfor i := int(s.firstIdx); i < len(st.strs); i++ {
        \t\tstr := st.strs[i]
        \t\tif str != jTombstoneStr {
        \t\t\tst.strs[w] = str
        \t\t\tif st.strIdx != nil {
        \t\t\t\tst.strIdx[str] = int32(w)
        \t\t\t}
        \t\t\tw++
        \t\t}
        \t}
        \tclear(st.strs[w:])
        \tst.strs = st.strs[:w]
        \ts.holes = 0
        \ts.firstIdx = 0
        }

        func (s *JSetCore) rebuildSetGenBuckets(st *jSetAnyState) {
        \tg := st.gen
        \tif g == nil {
        \t\treturn
        \t}
        \tcapB := 16
        \tfor capB < int(g.genLive)*2 {
        \t\tcapB <<= 1
        \t}
        \tif len(g.buckets) < capB {
        \t\tg.buckets = make([]int32, capB)
        \t} else {
        \t\tclear(g.buckets)
        \t}
        \tif cap(g.hashes) < len(st.items) {
        \t\tg.hashes = make([]int32, len(st.items))
        \t\tg.hashNext = make([]int32, len(st.items))
        \t} else {
        \t\tg.hashes = g.hashes[:len(st.items)]
        \t\tg.hashNext = g.hashNext[:len(st.items)]
        \t\tclear(g.hashNext)
        \t}
        \tmask := uint32(len(g.buckets) - 1)
        \tvar cnt int32
        \tfor i, it := range st.items {
        \t\tif it == jTombstoneObj {
        \t\t\tcontinue
        \t\t}
        \t\tif _, fast := fastSetKey(it, s.identity); fast {
        \t\t\tcontinue
        \t\t}
        \t\th := elemJavaHashCode(it, s.identity)
        \t\tg.hashes[i] = h
        \t\tslot := smearHash(h) & mask
        \t\tg.hashNext[i] = g.buckets[slot]
        \t\tg.buckets[slot] = int32(i + 1)
        \t\tcnt++
        \t}
        \tg.genLive = cnt
        }

        func (s *JSetCore) compactAny(st *jSetAnyState) {
        \tif s.holes == 0 {
        \t\ts.firstIdx = 0
        \t\treturn
        \t}
        \tw := 0
        \tfor i := int(s.firstIdx); i < len(st.items); i++ {
        \t\tit := st.items[i]
        \t\tif it != jTombstoneObj {
        \t\t\tst.items[w] = it
        \t\t\tif st.ptrIdx != nil {
        \t\t\t\tif fk, ok := fastSetKey(it, s.identity); ok {
        \t\t\t\t\tst.ptrIdx[fk] = int32(w)
        \t\t\t\t}
        \t\t\t}
        \t\t\tw++
        \t\t}
        \t}
        \tclear(st.items[w:])
        \tst.items = st.items[:w]
        \ts.holes = 0
        \ts.firstIdx = 0
        \tif st.gen != nil {
        \t\ts.rebuildSetGenBuckets(st)
        \t}
        }

        func (s *JSetCore) materializeAnyIdx(st *jSetAnyState) {
        \tvar needPtr, needGen int32
        \tfor _, it := range st.items {
        \t\tif it == jTombstoneObj {
        \t\t\tcontinue
        \t\t}
        \t\tif _, ok := fastSetKey(it, s.identity); ok {
        \t\t\tneedPtr++
        \t\t} else {
        \t\t\tneedGen++
        \t\t}
        \t}
        \tif needPtr > 0 {
        \t\tif st.ptrIdx == nil {
        \t\t\tst.ptrIdx = make(map[any]int32, int(needPtr)+4)
        \t\t}
        \t\tfor i, it := range st.items {
        \t\t\tif it == jTombstoneObj {
        \t\t\t\tcontinue
        \t\t\t}
        \t\t\tif fk, ok := fastSetKey(it, s.identity); ok {
        \t\t\t\tst.ptrIdx[fk] = int32(i)
        \t\t\t}
        \t\t}
        \t}
        \tif needGen > 0 {
        \t\tif st.gen == nil {
        \t\t\tst.gen = &jSetGenBuckets{genLive: needGen}
        \t\t}
        \t\ts.rebuildSetGenBuckets(st)
        \t}
        }

        func newSetAnyState(capHint int) *jSetAnyState {
        \tif capHint <= 4 {
        \t\ts4 := &jSetAnyState4{}
        \t\ts4.state.items = s4.buf[:0]
        \t\ts4.state.allIdent = true
        \t\treturn &s4.state
        \t}
        \treturn &jSetAnyState{items: make([]any, 0, capHint), allIdent: true}
        }

        func (s *JSetCore) upgradeToAny() {
        \tswitch s.mode {
        \tcase jSetModeEmpty:
        \t\ts.ext = newSetAnyState(s.initCap())
        \t\ts.mode = jSetModeAny
        \t\ts.firstIdx = 0
        \tcase jSetModeNode:
        \t\tst := newSetAnyState(max(int(s.live)+2, s.initCap()))
        \t\tst.allIdent = false
        \t\tfor i := int(s.firstIdx); i < len(s.nodes); i++ {
        \t\t\tif u := s.nodes[i]; u != jTombstoneNode {
        \t\t\t\tst.items = append(st.items, boxNodeById(u))
        \t\t\t}
        \t\t}
        \t\ts.nodes = nil
        \t\ts.holes = 0
        \t\ts.firstIdx = 0
        \t\ts.ext = st
        \t\ts.mode = jSetModeAny
        \t\tif len(st.items) > jSmallLimit {
        \t\t\ts.materializeAnyIdx(st)
        \t\t}
        \tcase jSetModeStr:
        \t\told := s.ext.(*jSetStrState)
        \t\tst := newSetAnyState(max(int(s.live)+2, s.initCap()))
        \t\tst.allIdent = false
        \t\tfor i := int(s.firstIdx); i < len(old.strs); i++ {
        \t\t\tstr := old.strs[i]
        \t\t\tif str == jTombstoneStr {
        \t\t\t\tcontinue
        \t\t\t}
        \t\t\tif str == NullString {
        \t\t\t\tst.items = append(st.items, nil)
        \t\t\t} else {
        \t\t\t\tst.items = append(st.items, str)
        \t\t\t}
        \t\t}
        \t\ts.holes = 0
        \t\ts.firstIdx = 0
        \t\ts.ext = st
        \t\ts.mode = jSetModeAny
        \t\tif len(st.items) > jSmallLimit {
        \t\t\ts.materializeAnyIdx(st)
        \t\t}
        \t}
        }

        func (s *JSetCore) findGenIdx(st *jSetAnyState, k any, h int32) int32 {
        \tg := st.gen
        \tif g == nil || len(g.buckets) == 0 {
        \t\treturn -1
        \t}
        \tslot := smearHash(h) & uint32(len(g.buckets)-1)
        \tif !s.identity {
        \t\tif _, isId := k.(jIdentityKey); !isId {
        \t\t\tif fk, ok := k.(FastCustomKey); ok {
        \t\t\t\tfor cur := g.buckets[slot]; cur != 0; cur = g.hashNext[cur-1] {
        \t\t\t\t\ti := cur - 1
        \t\t\t\t\tif g.hashes[i] == h {
        \t\t\t\t\t\tit := st.items[i]
        \t\t\t\t\t\tif it == k || fk.Equals(it) {
        \t\t\t\t\t\t\treturn i
        \t\t\t\t\t\t}
        \t\t\t\t\t}
        \t\t\t\t}
        \t\t\t\treturn -1
        \t\t\t}
        \t\t}
        \t}
        \tfor cur := g.buckets[slot]; cur != 0; cur = g.hashNext[cur-1] {
        \t\ti := cur - 1
        \t\tif g.hashes[i] == h {
        \t\t\tit := st.items[i]
        \t\t\tif s.identity {
        \t\t\t\tif RefEquals(it, k) {
        \t\t\t\t\treturn i
        \t\t\t\t}
        \t\t\t} else if javalang.ObjectsEquals(it, k) {
        \t\t\t\treturn i
        \t\t\t}
        \t\t}
        \t}
        \treturn -1
        }

        func (s *JSetCore) unlinkGenIdx(st *jSetAnyState, idx int32, h int32) {
        \tg := st.gen
        \tif g == nil || len(g.buckets) == 0 {
        \t\treturn
        \t}
        \tslot := smearHash(h) & uint32(len(g.buckets)-1)
        \ttarget := idx + 1
        \tvar prev int32
        \tfor cur := g.buckets[slot]; cur != 0; cur = g.hashNext[cur-1] {
        \t\tif cur == target {
        \t\t\tnext := g.hashNext[cur-1]
        \t\t\tif prev == 0 {
        \t\t\t\tg.buckets[slot] = next
        \t\t\t} else {
        \t\t\t\tg.hashNext[prev-1] = next
        \t\t\t}
        \t\t\tg.genLive--
        \t\t\treturn
        \t\t}
        \t\tprev = cur
        \t}
        }

        func (s *JSetCore) containsAllIdent(st *jSetAnyState, k any) bool {
        \tif k == nil {
        \t\treturn false
        \t}
        \tif st.ptrIdx != nil {
        \t\tif isFastIdentElem(k) {
        \t\t\t_, ok := st.ptrIdx[k]
        \t\t\treturn ok
        \t\t}
        \t\tswitch k.(type) {
        \t\tcase string, *BoxInt, int32, int64, FastCustomKey, jRawNodeProvider, nodestore.Node:
        \t\t\treturn false
        \t\t}
        \t\tif !reflect.TypeOf(k).Comparable() {
        \t\t\treturn false
        \t\t}
        \t\t_, ok := st.ptrIdx[k]
        \t\treturn ok
        \t}
        \tif s.live >= 2 && !s.imm {
        \t\tst.reads++
        \t\tif st.reads >= 32 {
        \t\t\ts.materializeAnyIdx(st)
        \t\t\tif isFastIdentElem(k) {
        \t\t\t\t_, ok := st.ptrIdx[k]
        \t\t\t\treturn ok
        \t\t\t}
        \t\t\tswitch k.(type) {
        \t\t\tcase string, *BoxInt, int32, int64, FastCustomKey, jRawNodeProvider, nodestore.Node:
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tif !reflect.TypeOf(k).Comparable() {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\t_, ok := st.ptrIdx[k]
        \t\t\treturn ok
        \t\t}
        \t}
        \tfor _, it := range st.items {
        \t\tif it == k {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (s *JSetCore) containsAnyNormalized(k any) bool {
        \tst := s.ext.(*jSetAnyState)
        \tif st.allIdent {
        \t\treturn s.containsAllIdent(st, k)
        \t}
        \tif st.ptrIdx != nil || st.gen != nil {
        \t\tif fk, ok := fastSetKey(k, s.identity); ok {
        \t\t\tif st.ptrIdx == nil {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\t_, found := st.ptrIdx[fk]
        \t\t\treturn found
        \t\t}
        \t\th := elemJavaHashCode(k, s.identity)
        \t\treturn s.findGenIdx(st, k, h) >= 0
        \t}
        \tif k != nil {
        \t\tif ik, ok := k.(jIdentityKey); ok {
        \t\t\tif !ik.IsJavaNil_() {
        \t\t\t\tfor _, it := range st.items {
        \t\t\t\t\tif it == k {
        \t\t\t\t\t\treturn true
        \t\t\t\t\t}
        \t\t\t\t}
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tk = nil
        \t\t} else if !s.identity {
        \t\t\tif _, isRN := k.(jRawNodeProvider); !isRN {
        \t\t\t\tif _, isBI := k.(*BoxInt); !isBI {
        \t\t\t\t\tif fk, ok := k.(FastCustomKey); ok && !fk.IsJavaNil_() {
        \t\t\t\t\t\tfor _, it := range st.items {
        \t\t\t\t\t\t\tif it != jTombstoneObj && (it == k || fk.Equals(it)) {
        \t\t\t\t\t\t\t\treturn true
        \t\t\t\t\t\t\t}
        \t\t\t\t\t\t}
        \t\t\t\t\t\treturn false
        \t\t\t\t\t}
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t}
        \tfor _, it := range st.items {
        \t\tif it == jTombstoneObj {
        \t\t\tcontinue
        \t\t}
        \t\tif s.identity {
        \t\t\tif RefEquals(it, k) {
        \t\t\t\treturn true
        \t\t\t}
        \t\t} else if jSetElemEquals(it, k) {
        \t\t\treturn true
        \t\t}
        \t}
        \treturn false
        }

        func (s *JSetCore) addAnyNormalized(k any) bool {
        \tst := s.ext.(*jSetAnyState)
        \tif st.ptrIdx != nil || st.gen != nil {
        \t\tif fk, ok := fastSetKey(k, s.identity); ok {
        \t\t\tif st.ptrIdx == nil {
        \t\t\t\tst.ptrIdx = make(map[any]int32, 8)
        \t\t\t} else if _, exists := st.ptrIdx[fk]; exists {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tif int(s.holes)*2 >= len(st.items) && len(st.items) >= jSmallLimit {
        \t\t\t\ts.compactAny(st)
        \t\t\t}
        \t\t\tif st.allIdent && !isNonNilIdentityElem(k) {
        \t\t\t\tst.allIdent = false
        \t\t\t}
        \t\t\tidx := int32(len(st.items))
        \t\t\tst.items = append(st.items, k)
        \t\t\tst.ptrIdx[fk] = idx
        \t\t\tif st.gen != nil {
        \t\t\t\tst.gen.hashes = append(st.gen.hashes, 0)
        \t\t\t\tst.gen.hashNext = append(st.gen.hashNext, 0)
        \t\t\t}
        \t\t\ts.live++
        \t\t\treturn true
        \t\t}
        \t\th := elemJavaHashCode(k, s.identity)
        \t\tif st.gen != nil && s.findGenIdx(st, k, h) >= 0 {
        \t\t\treturn false
        \t\t}
        \t\tif int(s.holes)*2 >= len(st.items) && len(st.items) >= jSmallLimit {
        \t\t\ts.compactAny(st)
        \t\t}
        \t\tst.allIdent = false
        \t\tidx := int32(len(st.items))
        \t\tst.items = append(st.items, k)
        \t\ts.live++
        \t\tif st.gen == nil {
        \t\t\tst.gen = &jSetGenBuckets{genLive: 1}
        \t\t\ts.rebuildSetGenBuckets(st)
        \t\t\treturn true
        \t\t}
        \t\tst.gen.genLive++
        \t\tif int(st.gen.genLive) >= len(st.gen.buckets) {
        \t\t\ts.rebuildSetGenBuckets(st)
        \t\t} else {
        \t\t\tslot := smearHash(h) & uint32(len(st.gen.buckets)-1)
        \t\t\tst.gen.hashes = append(st.gen.hashes, h)
        \t\t\tst.gen.hashNext = append(st.gen.hashNext, st.gen.buckets[slot])
        \t\t\tst.gen.buckets[slot] = idx + 1
        \t\t}
        \t\treturn true
        \t}
        \tif k != nil {
        \t\tif ik, ok := k.(jIdentityKey); ok {
        \t\t\tif !ik.IsJavaNil_() {
        \t\t\t\tfor _, it := range st.items {
        \t\t\t\t\tif it == k {
        \t\t\t\t\t\treturn false
        \t\t\t\t\t}
        \t\t\t\t}
        \t\t\t\tif s.holes > 0 && int(s.holes)*2 >= len(st.items) {
        \t\t\t\t\ts.compactAny(st)
        \t\t\t\t}
        \t\t\t\tst.items = append(st.items, k)
        \t\t\t\ts.live++
        \t\t\t\tif len(st.items) > jSmallLimit {
        \t\t\t\t\ts.materializeAnyIdx(st)
        \t\t\t\t}
        \t\t\t\treturn true
        \t\t\t}
        \t\t\tk = nil
        \t\t} else if !s.identity {
        \t\t\tif _, isRN := k.(jRawNodeProvider); !isRN {
        \t\t\t\tif _, isBI := k.(*BoxInt); !isBI {
        \t\t\t\t\tif fk, ok := k.(FastCustomKey); ok && !fk.IsJavaNil_() {
        \t\t\t\t\t\tfor _, it := range st.items {
        \t\t\t\t\t\t\tif it != jTombstoneObj && (it == k || fk.Equals(it)) {
        \t\t\t\t\t\t\t\treturn false
        \t\t\t\t\t\t\t}
        \t\t\t\t\t\t}
        \t\t\t\t\t\tif s.holes > 0 && int(s.holes)*2 >= len(st.items) {
        \t\t\t\t\t\t\ts.compactAny(st)
        \t\t\t\t\t\t}
        \t\t\t\t\t\tst.allIdent = false
        \t\t\t\t\t\tst.items = append(st.items, k)
        \t\t\t\t\t\ts.live++
        \t\t\t\t\t\tif len(st.items) > jSmallLimit {
        \t\t\t\t\t\t\ts.materializeAnyIdx(st)
        \t\t\t\t\t\t}
        \t\t\t\t\t\treturn true
        \t\t\t\t\t}
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t}
        \tfor _, it := range st.items {
        \t\tif it == jTombstoneObj {
        \t\t\tcontinue
        \t\t}
        \t\tif s.identity {
        \t\t\tif RefEquals(it, k) {
        \t\t\t\treturn false
        \t\t\t}
        \t\t} else if jSetElemEquals(it, k) {
        \t\t\treturn false
        \t\t}
        \t}
        \tif s.holes > 0 && int(s.holes)*2 >= len(st.items) {
        \t\ts.compactAny(st)
        \t}
        \tst.allIdent = false
        \tst.items = append(st.items, k)
        \ts.live++
        \tif len(st.items) > jSmallLimit {
        \t\ts.materializeAnyIdx(st)
        \t}
        \treturn true
        }

        func (s *JSetCore) removeAnyAt(idx int) bool {
        \tst := s.ext.(*jSetAnyState)
        \tit := st.items[idx]
        \tif it == jTombstoneObj {
        \t\treturn false
        \t}
        \tif st.ptrIdx != nil {
        \t\tif fk, ok := fastSetKey(it, s.identity); ok {
        \t\t\tdelete(st.ptrIdx, fk)
        \t\t} else if st.gen != nil {
        \t\t\ts.unlinkGenIdx(st, int32(idx), st.gen.hashes[idx])
        \t\t}
        \t} else if st.gen != nil {
        \t\ts.unlinkGenIdx(st, int32(idx), st.gen.hashes[idx])
        \t}
        \tst.items[idx] = jTombstoneObj
        \ts.live--
        \tif s.live == 0 {
        \t\tclear(st.items)
        \t\tst.items = st.items[:0]
        \t\tif st.gen != nil {
        \t\t\tst.gen.hashes = st.gen.hashes[:0]
        \t\t\tst.gen.hashNext = st.gen.hashNext[:0]
        \t\t}
        \t\tst.reads = 0
        \t\tst.allIdent = true
        \t\ts.holes = 0
        \t\ts.firstIdx = 0
        \t} else if idx == len(st.items)-1 {
        \t\tst.items = st.items[:idx]
        \t\tif st.gen != nil && int(idx) < len(st.gen.hashes) {
        \t\t\tst.gen.hashes = st.gen.hashes[:idx]
        \t\t\tst.gen.hashNext = st.gen.hashNext[:idx]
        \t\t}
        \t} else {
        \t\ts.holes++
        \t\tif idx == int(s.firstIdx) {
        \t\t\ts.firstIdx++
        \t\t}
        \t}
        \treturn true
        }

        func (s *JSetCore) removeAnyNormalized(k any) bool {
        \tst := s.ext.(*jSetAnyState)
        \tif st.ptrIdx != nil || st.gen != nil {
        \t\tif fk, ok := fastSetKey(k, s.identity); ok {
        \t\t\tif st.ptrIdx == nil {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tidx, found := st.ptrIdx[fk]
        \t\t\tif !found {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\treturn s.removeAnyAt(int(idx))
        \t\t}
        \t\th := elemJavaHashCode(k, s.identity)
        \t\tidx := s.findGenIdx(st, k, h)
        \t\tif idx < 0 {
        \t\t\treturn false
        \t\t}
        \t\treturn s.removeAnyAt(int(idx))
        \t}
        \tif k != nil {
        \t\tif ik, ok := k.(jIdentityKey); ok {
        \t\t\tif !ik.IsJavaNil_() {
        \t\t\t\tfor i, it := range st.items {
        \t\t\t\t\tif it == k {
        \t\t\t\t\t\treturn s.removeAnyAt(i)
        \t\t\t\t\t}
        \t\t\t\t}
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tk = nil
        \t\t} else if !s.identity {
        \t\t\tif _, isRN := k.(jRawNodeProvider); !isRN {
        \t\t\t\tif _, isBI := k.(*BoxInt); !isBI {
        \t\t\t\t\tif fk, ok := k.(FastCustomKey); ok && !fk.IsJavaNil_() {
        \t\t\t\t\t\tfor i, it := range st.items {
        \t\t\t\t\t\t\tif it != jTombstoneObj && (it == k || fk.Equals(it)) {
        \t\t\t\t\t\t\t\treturn s.removeAnyAt(i)
        \t\t\t\t\t\t\t}
        \t\t\t\t\t\t}
        \t\t\t\t\t\treturn false
        \t\t\t\t\t}
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t}
        \tfor i, it := range st.items {
        \t\tif it == jTombstoneObj {
        \t\t\tcontinue
        \t\t}
        \t\tif s.identity {
        \t\t\tif RefEquals(it, k) {
        \t\t\t\treturn s.removeAnyAt(i)
        \t\t\t}
        \t\t} else if jSetElemEquals(it, k) {
        \t\t\treturn s.removeAnyAt(i)
        \t\t}
        \t}
        \treturn false
        }

        func (s *JSetCore) ContainsNode(u uint32) bool {
        \tif s.live == 0 {
        \t\treturn false
        \t}
        \tswitch s.mode {
        \tcase jSetModeNode:
        \t\tif u == jTombstoneNode {
        \t\t\treturn false
        \t\t}
        \t\tif s.ext != nil {
        \t\t\t_, ok := s.ext.(map[uint32]int32)[u]
        \t\t\treturn ok
        \t\t}
        \t\tfor _, x := range s.nodes {
        \t\t\tif x == u {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \tcase jSetModeStr:
        \t\tif u == 0 {
        \t\t\treturn s.ContainsStr(NullString)
        \t\t}
        \t\treturn false
        \tcase jSetModeAny:
        \t\tif u == 0 {
        \t\t\treturn s.containsAnyNormalized(nil)
        \t\t}
        \t\treturn s.containsAnyNormalized(boxNodeById(u))
        \t}
        \treturn false
        }

        func (s *JSetCore) AddNode(u uint32) bool {
        \tif s.imm {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tif s.identity || u == jTombstoneNode {
        \t\tif s.mode != jSetModeAny {
        \t\t\ts.upgradeToAny()
        \t\t}
        \t\treturn s.addAnyNormalized(boxNodeById(u))
        \t}
        \tswitch s.mode {
        \tcase jSetModeEmpty:
        \t\tif s.capHint > 4 {
        \t\t\ts.nodes = make([]uint32, 1, min(int(s.capHint), 256))
        \t\t\ts.nodes[0] = u
        \t\t} else {
        \t\t\ts.nodeBuf[0] = u
        \t\t\ts.nodes = s.nodeBuf[:1]
        \t\t}
        \t\ts.mode = jSetModeNode
        \t\ts.live = 1
        \t\treturn true
        \tcase jSetModeNode:
        \t\tif s.ext != nil {
        \t\t\tidxMap := s.ext.(map[uint32]int32)
        \t\t\tif _, ok := idxMap[u]; ok {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tif int(s.holes)*2 >= len(s.nodes) && len(s.nodes) >= jNodeSmallLimit {
        \t\t\t\ts.compactNodes()
        \t\t\t}
        \t\t\tidxMap[u] = int32(len(s.nodes))
        \t\t\ts.nodes = append(s.nodes, u)
        \t\t\ts.live++
        \t\t\treturn true
        \t\t}
        \t\tfor _, x := range s.nodes {
        \t\t\tif x == u {
        \t\t\t\treturn false
        \t\t\t}
        \t\t}
        \t\tif s.holes > 0 && int(s.holes)*2 >= len(s.nodes) {
        \t\t\ts.compactNodes()
        \t\t}
        \t\tif len(s.nodes) >= jNodeSmallLimit {
        \t\t\tidxMap := make(map[uint32]int32, max(len(s.nodes)+4, s.initCap()))
        \t\t\tfor i, x := range s.nodes {
        \t\t\t\tif x != jTombstoneNode {
        \t\t\t\t\tidxMap[x] = int32(i)
        \t\t\t\t}
        \t\t\t}
        \t\t\tidxMap[u] = int32(len(s.nodes))
        \t\t\ts.ext = idxMap
        \t\t}
        \t\tif len(s.nodes) == 1 && cap(s.nodes) == 1 {
        \t\t\tnb := make([]uint32, 2, max(8, s.initCap()))
        \t\t\tnb[0] = s.nodes[0]
        \t\t\tnb[1] = u
        \t\t\ts.nodes = nb
        \t\t} else {
        \t\t\ts.nodes = append(s.nodes, u)
        \t\t}
        \t\ts.live++
        \t\treturn true
        \tcase jSetModeStr:
        \t\tif u == 0 {
        \t\t\treturn s.AddStr(NullString)
        \t\t}
        \t\ts.upgradeToAny()
        \t\treturn s.addAnyNormalized(boxNodeById(u))
        \tdefault:
        \t\tif u == 0 {
        \t\t\treturn s.addAnyNormalized(nil)
        \t\t}
        \t\treturn s.addAnyNormalized(boxNodeById(u))
        \t}
        }

        func (s *JSetCore) removeNodeAt(idx int) bool {
        \tu := s.nodes[idx]
        \tif u == jTombstoneNode {
        \t\treturn false
        \t}
        \tif s.ext != nil {
        \t\tdelete(s.ext.(map[uint32]int32), u)
        \t}
        \ts.nodes[idx] = jTombstoneNode
        \ts.live--
        \tif s.live == 0 {
        \t\ts.nodes = s.nodes[:0]
        \t\ts.holes = 0
        \t\ts.firstIdx = 0
        \t} else if idx == len(s.nodes)-1 {
        \t\ts.nodes = s.nodes[:idx]
        \t} else {
        \t\ts.holes++
        \t\tif idx == int(s.firstIdx) {
        \t\t\ts.firstIdx++
        \t\t}
        \t}
        \treturn true
        }

        func (s *JSetCore) RemoveNode(u uint32) bool {
        \tif s.imm {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tif s.live == 0 {
        \t\treturn false
        \t}
        \tswitch s.mode {
        \tcase jSetModeNode:
        \t\tif u == jTombstoneNode {
        \t\t\treturn false
        \t\t}
        \t\tif s.ext != nil {
        \t\t\tidx, ok := s.ext.(map[uint32]int32)[u]
        \t\t\tif !ok {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\treturn s.removeNodeAt(int(idx))
        \t\t}
        \t\tfor i, x := range s.nodes {
        \t\t\tif x == u {
        \t\t\t\treturn s.removeNodeAt(i)
        \t\t\t}
        \t\t}
        \t\treturn false
        \tcase jSetModeStr:
        \t\tif u == 0 {
        \t\t\treturn s.RemoveStr(NullString)
        \t\t}
        \t\treturn false
        \tcase jSetModeAny:
        \t\tif u == 0 {
        \t\t\treturn s.removeAnyNormalized(nil)
        \t\t}
        \t\treturn s.removeAnyNormalized(boxNodeById(u))
        \t}
        \treturn false
        }

        func (s *JSetCore) ContainsStr(k string) bool {
        \tif s.live == 0 {
        \t\treturn false
        \t}
        \tswitch s.mode {
        \tcase jSetModeStr:
        \t\tif k == jTombstoneStr {
        \t\t\treturn false
        \t\t}
        \t\tst := s.ext.(*jSetStrState)
        \t\tif st.strIdx != nil {
        \t\t\t_, ok := st.strIdx[k]
        \t\t\treturn ok
        \t\t}
        \t\tfor _, x := range st.strs {
        \t\t\tif x == k {
        \t\t\t\treturn true
        \t\t\t}
        \t\t}
        \t\treturn false
        \tcase jSetModeNode:
        \t\tif k == NullString {
        \t\t\treturn s.ContainsNode(0)
        \t\t}
        \t\treturn false
        \tcase jSetModeAny:
        \t\tif k == NullString {
        \t\t\treturn s.containsAnyNormalized(nil)
        \t\t}
        \t\treturn s.containsAnyNormalized(k)
        \t}
        \treturn false
        }

        func (s *JSetCore) AddStr(k string) bool {
        \tif s.imm {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tif s.identity || k == jTombstoneStr {
        \t\tif s.mode != jSetModeAny {
        \t\t\ts.upgradeToAny()
        \t\t}
        \t\tif k == NullString {
        \t\t\treturn s.addAnyNormalized(nil)
        \t\t}
        \t\treturn s.addAnyNormalized(k)
        \t}
        \tswitch s.mode {
        \tcase jSetModeEmpty:
        \t\tvar st *jSetStrState
        \t\tif s.capHint > 4 {
        \t\t\tst = &jSetStrState{strs: make([]string, 1, min(int(s.capHint), 256))}
        \t\t\tst.strs[0] = k
        \t\t} else {
        \t\t\ts4 := &jSetStrState4{}
        \t\t\ts4.buf[0] = k
        \t\t\tst = &s4.state
        \t\t\tst.strs = s4.buf[:1]
        \t\t}
        \t\ts.ext = st
        \t\ts.mode = jSetModeStr
        \t\ts.live = 1
        \t\treturn true
        \tcase jSetModeStr:
        \t\tst := s.ext.(*jSetStrState)
        \t\tif st.strIdx != nil {
        \t\t\tif _, ok := st.strIdx[k]; ok {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\tif int(s.holes)*2 >= len(st.strs) && len(st.strs) >= jStrSmallLimit {
        \t\t\t\ts.compactStrs(st)
        \t\t\t}
        \t\t\tst.strIdx[k] = int32(len(st.strs))
        \t\t\tst.strs = append(st.strs, k)
        \t\t\ts.live++
        \t\t\treturn true
        \t\t}
        \t\tfor _, x := range st.strs {
        \t\t\tif x == k {
        \t\t\t\treturn false
        \t\t\t}
        \t\t}
        \t\tif s.holes > 0 && int(s.holes)*2 >= len(st.strs) {
        \t\t\ts.compactStrs(st)
        \t\t}
        \t\tif len(st.strs) >= jStrSmallLimit {
        \t\t\tst.strIdx = make(map[string]int32, max(len(st.strs)+1, s.initCap()))
        \t\t\tfor i, x := range st.strs {
        \t\t\t\tif x != jTombstoneStr {
        \t\t\t\t\tst.strIdx[x] = int32(i)
        \t\t\t\t}
        \t\t\t}
        \t\t\tst.strIdx[k] = int32(len(st.strs))
        \t\t}
        \t\tst.strs = append(st.strs, k)
        \t\ts.live++
        \t\treturn true
        \tcase jSetModeNode:
        \t\tif k == NullString {
        \t\t\treturn s.AddNode(0)
        \t\t}
        \t\ts.upgradeToAny()
        \t\treturn s.addAnyNormalized(k)
        \tdefault:
        \t\tif k == NullString {
        \t\t\treturn s.addAnyNormalized(nil)
        \t\t}
        \t\treturn s.addAnyNormalized(k)
        \t}
        }

        func (s *JSetCore) removeStrAt(idx int) bool {
        \tst := s.ext.(*jSetStrState)
        \tk := st.strs[idx]
        \tif k == jTombstoneStr {
        \t\treturn false
        \t}
        \tif st.strIdx != nil {
        \t\tdelete(st.strIdx, k)
        \t}
        \tst.strs[idx] = jTombstoneStr
        \ts.live--
        \tif s.live == 0 {
        \t\tclear(st.strs)
        \t\tst.strs = st.strs[:0]
        \t\ts.holes = 0
        \t\ts.firstIdx = 0
        \t} else if idx == len(st.strs)-1 {
        \t\tst.strs = st.strs[:idx]
        \t} else {
        \t\ts.holes++
        \t\tif idx == int(s.firstIdx) {
        \t\t\ts.firstIdx++
        \t\t}
        \t}
        \treturn true
        }

        func (s *JSetCore) RemoveStr(k string) bool {
        \tif s.imm {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tif s.live == 0 {
        \t\treturn false
        \t}
        \tswitch s.mode {
        \tcase jSetModeStr:
        \t\tif k == jTombstoneStr {
        \t\t\treturn false
        \t\t}
        \t\tst := s.ext.(*jSetStrState)
        \t\tif st.strIdx != nil {
        \t\t\tidx, ok := st.strIdx[k]
        \t\t\tif !ok {
        \t\t\t\treturn false
        \t\t\t}
        \t\t\treturn s.removeStrAt(int(idx))
        \t\t}
        \t\tfor i, x := range st.strs {
        \t\t\tif x == k {
        \t\t\t\treturn s.removeStrAt(i)
        \t\t\t}
        \t\t}
        \t\treturn false
        \tcase jSetModeNode:
        \t\tif k == NullString {
        \t\t\treturn s.RemoveNode(0)
        \t\t}
        \t\treturn false
        \tcase jSetModeAny:
        \t\tif k == NullString {
        \t\t\treturn s.removeAnyNormalized(nil)
        \t\t}
        \t\treturn s.removeAnyNormalized(k)
        \t}
        \treturn false
        }

        func (s *JSetCore) Contains(k any) bool {
        \tif s.live == 0 {
        \t\treturn false
        \t}
        \tif s.mode == jSetModeAny {
        \t\tst := s.ext.(*jSetAnyState)
        \t\tif st.allIdent {
        \t\t\treturn s.containsAllIdent(st, k)
        \t\t}
        \t}
        \tif !s.identity {
        \t\tif rn, ok := k.(jRawNodeProvider); ok {
        \t\t\treturn s.ContainsNode(rn.RawNodeId_())
        \t\t}
        \t\tif str, ok := k.(string); ok {
        \t\t\treturn s.ContainsStr(str)
        \t\t}
        \t\tif nd, ok := k.(nodestore.Node); ok {
        \t\t\treturn s.ContainsNode(uint32(nd))
        \t\t}
        \t}
        \tk = ToAny(k)
        \tif k == nil {
        \t\tswitch s.mode {
        \t\tcase jSetModeNode:
        \t\t\treturn s.ContainsNode(0)
        \t\tcase jSetModeStr:
        \t\t\treturn s.ContainsStr(NullString)
        \t\tcase jSetModeAny:
        \t\t\treturn s.containsAnyNormalized(nil)
        \t\tdefault:
        \t\t\treturn false
        \t\t}
        \t}
        \tif s.mode != jSetModeAny {
        \t\treturn false
        \t}
        \treturn s.containsAnyNormalized(k)
        }

        func (s *JSetCore) Add(k any) bool {
        \tif s.imm {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tif !s.identity {
        \t\tif rn, ok := k.(jRawNodeProvider); ok {
        \t\t\treturn s.AddNode(rn.RawNodeId_())
        \t\t}
        \t\tif str, ok := k.(string); ok {
        \t\t\treturn s.AddStr(str)
        \t\t}
        \t\tif nd, ok := k.(nodestore.Node); ok {
        \t\t\treturn s.AddNode(uint32(nd))
        \t\t}
        \t}
        \tk = ToAny(k)
        \tif k == nil {
        \t\tswitch s.mode {
        \t\tcase jSetModeNode:
        \t\t\treturn s.AddNode(0)
        \t\tcase jSetModeStr:
        \t\t\treturn s.AddStr(NullString)
        \t\tcase jSetModeEmpty:
        \t\t\ts.upgradeToAny()
        \t\t\treturn s.addAnyNormalized(nil)
        \t\tdefault:
        \t\t\treturn s.addAnyNormalized(nil)
        \t\t}
        \t}
        \tif s.mode != jSetModeAny {
        \t\ts.upgradeToAny()
        \t}
        \treturn s.addAnyNormalized(k)
        }

        func (s *JSetCore) Remove(k any) bool {
        \tif s.imm {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tif s.live == 0 {
        \t\treturn false
        \t}
        \tif !s.identity {
        \t\tif rn, ok := k.(jRawNodeProvider); ok {
        \t\t\treturn s.RemoveNode(rn.RawNodeId_())
        \t\t}
        \t\tif str, ok := k.(string); ok {
        \t\t\treturn s.RemoveStr(str)
        \t\t}
        \t\tif nd, ok := k.(nodestore.Node); ok {
        \t\t\treturn s.RemoveNode(uint32(nd))
        \t\t}
        \t}
        \tk = ToAny(k)
        \tif k == nil {
        \t\tswitch s.mode {
        \t\tcase jSetModeNode:
        \t\t\treturn s.RemoveNode(0)
        \t\tcase jSetModeStr:
        \t\t\treturn s.RemoveStr(NullString)
        \t\tcase jSetModeAny:
        \t\t\treturn s.removeAnyNormalized(nil)
        \t\tdefault:
        \t\t\treturn false
        \t\t}
        \t}
        \tif s.mode != jSetModeAny {
        \t\treturn false
        \t}
        \treturn s.removeAnyNormalized(k)
        }

        func (s *JSetCore) Clear() {
        \tif s.imm {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tif s.live == 0 && s.holes == 0 && s.firstIdx == 0 {
        \t\treturn
        \t}
        \tswitch s.mode {
        \tcase jSetModeNode:
        \t\ts.nodes = s.nodes[:0]
        \t\tif s.ext != nil {
        \t\t\tclear(s.ext.(map[uint32]int32))
        \t\t}
        \tcase jSetModeStr:
        \t\tst := s.ext.(*jSetStrState)
        \t\tclear(st.strs)
        \t\tst.strs = st.strs[:0]
        \t\tif st.strIdx != nil {
        \t\t\tclear(st.strIdx)
        \t\t}
        \tcase jSetModeAny:
        \t\tst := s.ext.(*jSetAnyState)
        \t\tclear(st.items)
        \t\tst.items = st.items[:0]
        \t\tif st.ptrIdx != nil {
        \t\t\tclear(st.ptrIdx)
        \t\t}
        \t\tif st.gen != nil {
        \t\t\tst.gen.hashes = st.gen.hashes[:0]
        \t\t\tst.gen.hashNext = st.gen.hashNext[:0]
        \t\t\tclear(st.gen.buckets)
        \t\t\tst.gen.genLive = 0
        \t\t}
        \t\tst.reads = 0
        \t\tst.allIdent = true
        \t}
        \ts.live = 0
        \ts.holes = 0
        \ts.firstIdx = 0
        }

        func (s *JSetCore) ForEach(fn func(any) bool) {
        \tif s == nil || s.live == 0 {
        \t\treturn
        \t}
        \tstart := int(s.firstIdx)
        \tswitch s.mode {
        \tcase jSetModeNode:
        \t\tfor i := start; i < len(s.nodes); i++ {
        \t\t\tif u := s.nodes[i]; u != jTombstoneNode {
        \t\t\t\tif !fn(boxNodeById(u)) {
        \t\t\t\t\treturn
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \tcase jSetModeStr:
        \t\tst := s.ext.(*jSetStrState)
        \t\tfor i := start; i < len(st.strs); i++ {
        \t\t\tstr := st.strs[i]
        \t\t\tif str == jTombstoneStr {
        \t\t\t\tcontinue
        \t\t\t}
        \t\t\tvar v any
        \t\t\tif str != NullString {
        \t\t\t\tv = str
        \t\t\t}
        \t\t\tif !fn(v) {
        \t\t\t\treturn
        \t\t\t}
        \t\t}
        \tcase jSetModeAny:
        \t\tst := s.ext.(*jSetAnyState)
        \t\tfor i := start; i < len(st.items); i++ {
        \t\t\tif it := st.items[i]; it != jTombstoneObj {
        \t\t\t\tif !fn(it) {
        \t\t\t\t\treturn
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t}
        }

        func (s *JSetCore) LiveKeys() []any {
        \tif s == nil || s.live == 0 {
        \t\treturn []any{}
        \t}
        \tout := make([]any, 0, s.live)
        \ts.ForEach(func(k any) bool {
        \t\tout = append(out, k)
        \t\treturn true
        \t})
        \treturn out
        }

        func (s *JSetCore) SetHash() int32 {
        \tif s == nil || s.live == 0 {
        \t\treturn 0
        \t}
        \tvar h int32
        \tstart := int(s.firstIdx)
        \tswitch s.mode {
        \tcase jSetModeNode:
        \t\tfor i := start; i < len(s.nodes); i++ {
        \t\t\tif u := s.nodes[i]; u != jTombstoneNode {
        \t\t\t\th += int32(u)
        \t\t\t}
        \t\t}
        \tcase jSetModeStr:
        \t\tst := s.ext.(*jSetStrState)
        \t\tfor i := start; i < len(st.strs); i++ {
        \t\t\tif str := st.strs[i]; str != jTombstoneStr && str != NullString {
        \t\t\t\th += javalang.StringHashCode(str)
        \t\t\t}
        \t\t}
        \tcase jSetModeAny:
        \t\tst := s.ext.(*jSetAnyState)
        \t\tfor i := start; i < len(st.items); i++ {
        \t\t\tif it := st.items[i]; it != jTombstoneObj {
        \t\t\t\th += elemJavaHashCode(it, s.identity)
        \t\t\t}
        \t\t}
        \t}
        \treturn h
        }

        func (s *JSetCore) FirstKey() any {
        \tif s == nil || s.live == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty set"))
        \t}
        \tstart := int(s.firstIdx)
        \tswitch s.mode {
        \tcase jSetModeNode:
        \t\tfor i := start; i < len(s.nodes); i++ {
        \t\t\tif u := s.nodes[i]; u != jTombstoneNode {
        \t\t\t\tif i != start && !s.imm {
        \t\t\t\t\ts.firstIdx = int32(i)
        \t\t\t\t}
        \t\t\t\treturn boxNodeById(u)
        \t\t\t}
        \t\t}
        \tcase jSetModeStr:
        \t\tst := s.ext.(*jSetStrState)
        \t\tfor i := start; i < len(st.strs); i++ {
        \t\t\tif str := st.strs[i]; str != jTombstoneStr {
        \t\t\t\tif i != start && !s.imm {
        \t\t\t\t\ts.firstIdx = int32(i)
        \t\t\t\t}
        \t\t\t\tif str == NullString {
        \t\t\t\t\treturn nil
        \t\t\t\t}
        \t\t\t\treturn str
        \t\t\t}
        \t\t}
        \tcase jSetModeAny:
        \t\tst := s.ext.(*jSetAnyState)
        \t\tfor i := start; i < len(st.items); i++ {
        \t\t\tif it := st.items[i]; it != jTombstoneObj {
        \t\t\t\tif i != start && !s.imm {
        \t\t\t\t\ts.firstIdx = int32(i)
        \t\t\t\t}
        \t\t\t\treturn it
        \t\t\t}
        \t\t}
        \t}
        \tpanic(javalang.NewNoSuchElementException("empty set"))
        }

        func (s *JSetCore) LastKey() any {
        \tif s == nil || s.live == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty set"))
        \t}
        \tswitch s.mode {
        \tcase jSetModeNode:
        \t\tfor i := len(s.nodes) - 1; i >= 0; i-- {
        \t\t\tif s.nodes[i] != jTombstoneNode {
        \t\t\t\treturn boxNodeById(s.nodes[i])
        \t\t\t}
        \t\t}
        \tcase jSetModeStr:
        \t\tst := s.ext.(*jSetStrState)
        \t\tfor i := len(st.strs) - 1; i >= 0; i-- {
        \t\t\tif str := st.strs[i]; str != jTombstoneStr {
        \t\t\t\tif str == NullString {
        \t\t\t\t\treturn nil
        \t\t\t\t}
        \t\t\t\treturn str
        \t\t\t}
        \t\t}
        \tcase jSetModeAny:
        \t\tst := s.ext.(*jSetAnyState)
        \t\tfor i := len(st.items) - 1; i >= 0; i-- {
        \t\t\tif it := st.items[i]; it != jTombstoneObj {
        \t\t\t\treturn it
        \t\t\t}
        \t\t}
        \t}
        \tpanic(javalang.NewNoSuchElementException("empty set"))
        }

        func (s *JSetCore) RemoveFirstKey() any {
        \tif s == nil || s.live == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty set"))
        \t}
        \tif s.imm {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tstart := int(s.firstIdx)
        \tswitch s.mode {
        \tcase jSetModeNode:
        \t\tfor i := start; i < len(s.nodes); i++ {
        \t\t\tif u := s.nodes[i]; u != jTombstoneNode {
        \t\t\t\ts.firstIdx = int32(i)
        \t\t\t\ts.removeNodeAt(i)
        \t\t\t\treturn boxNodeById(u)
        \t\t\t}
        \t\t}
        \tcase jSetModeStr:
        \t\tst := s.ext.(*jSetStrState)
        \t\tfor i := start; i < len(st.strs); i++ {
        \t\t\tif str := st.strs[i]; str != jTombstoneStr {
        \t\t\t\ts.firstIdx = int32(i)
        \t\t\t\ts.removeStrAt(i)
        \t\t\t\tif str == NullString {
        \t\t\t\t\treturn nil
        \t\t\t\t}
        \t\t\t\treturn str
        \t\t\t}
        \t\t}
        \tcase jSetModeAny:
        \t\tst := s.ext.(*jSetAnyState)
        \t\tfor i := start; i < len(st.items); i++ {
        \t\t\tif it := st.items[i]; it != jTombstoneObj {
        \t\t\t\ts.firstIdx = int32(i)
        \t\t\t\ts.removeAnyAt(i)
        \t\t\t\treturn it
        \t\t\t}
        \t\t}
        \t}
        \tpanic(javalang.NewNoSuchElementException("empty set"))
        }

        func (s *JSetCore) RemoveLastKey() any {
        \tif s == nil || s.live == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty set"))
        \t}
        \tswitch s.mode {
        \tcase jSetModeNode:
        \t\tfor i := len(s.nodes) - 1; i >= 0; i-- {
        \t\t\tif u := s.nodes[i]; u != jTombstoneNode {
        \t\t\t\ts.removeNodeAt(i)
        \t\t\t\treturn boxNodeById(u)
        \t\t\t}
        \t\t}
        \tcase jSetModeStr:
        \t\tst := s.ext.(*jSetStrState)
        \t\tfor i := len(st.strs) - 1; i >= 0; i-- {
        \t\t\tif str := st.strs[i]; str != jTombstoneStr {
        \t\t\t\ts.removeStrAt(i)
        \t\t\t\tif str == NullString {
        \t\t\t\t\treturn nil
        \t\t\t\t}
        \t\t\t\treturn str
        \t\t\t}
        \t\t}
        \tcase jSetModeAny:
        \t\tst := s.ext.(*jSetAnyState)
        \t\tfor i := len(st.items) - 1; i >= 0; i-- {
        \t\t\tif it := st.items[i]; it != jTombstoneObj {
        \t\t\t\ts.removeAnyAt(i)
        \t\t\t\treturn it
        \t\t\t}
        \t\t}
        \t}
        \tpanic(javalang.NewNoSuchElementException("empty set"))
        }

        func (s *JSetCore) AddLast(k any) {
        \tif s.imm {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \ts.Remove(k)
        \ts.Add(k)
        }

        func (s *JSetCore) AddFirst(k any) {
        \tif s.imm {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tif s.live == 0 {
        \t\ts.Add(k)
        \t\treturn
        \t}
        \ts.Remove(k)
        \trest := s.LiveKeys()
        \ts.Clear()
        \ts.Add(k)
        \tfor _, x := range rest {
        \t\ts.Add(x)
        \t}
        }

        func (s *JSetCore) sealForImmutable() {
        \tif s == nil || s.imm {
        \t\treturn
        \t}
        \tswitch s.mode {
        \tcase jSetModeNode:
        \t\ts.compactNodes()
        \t\tif s.ext == nil && len(s.nodes) > jNodeSmallLimit {
        \t\t\tidxMap := make(map[uint32]int32, len(s.nodes))
        \t\t\tfor i, u := range s.nodes {
        \t\t\t\tidxMap[u] = int32(i)
        \t\t\t}
        \t\t\ts.ext = idxMap
        \t\t}
        \tcase jSetModeStr:
        \t\tst := s.ext.(*jSetStrState)
        \t\ts.compactStrs(st)
        \t\tif st.strIdx == nil && len(st.strs) > jStrSmallLimit {
        \t\t\tst.strIdx = make(map[string]int32, len(st.strs))
        \t\t\tfor i, str := range st.strs {
        \t\t\t\tst.strIdx[str] = int32(i)
        \t\t\t}
        \t\t}
        \tcase jSetModeAny:
        \t\tst := s.ext.(*jSetAnyState)
        \t\ts.compactAny(st)
        \t\tif st.ptrIdx == nil && st.gen == nil && len(st.items) > jSmallLimit {
        \t\t\ts.materializeAnyIdx(st)
        \t\t}
        \t}
        \ts.firstIdx = 0
        \ts.imm = true
        }
        """;
  }

  private String buildCoreRuntimeGoPart2Map() {
    return """

        // --- OrderedMap-Style Slim JMapCore (56 bytes, 24B/32B homogeneous entries, 0 inline bloat) ---

        const (
        \tjMapModeEmpty     uint8 = 0
        \tjMapModeStr       uint8 = 1
        \tjMapModeNode      uint8 = 2
        \tjMapModePtr       uint8 = 3
        \tjMapModeGeneral   uint8 = 4
        \tjMapModeSortedNat uint8 = 5
        \tjMapModeSortedCmp uint8 = 6
        \tjMapModeNullOnly  uint8 = 7
        \tjMapModeNode1     uint8 = 8

        \tjKeyGeneral   uint8  = 0
        \tjKeyNull      uint8  = 1
        \tjKeyStr       uint8  = 2
        \tjKeyNode      uint8  = 3
        \tjKeyInt       uint8  = 4
        \tjKeyPtr       uint8  = 5
        \tjKeyFree      uint8  = 7
        \tjEntryPosMask uint32 = 0x1fffffff
        )

        type jStrMapEntry struct {
        \tkey string
        \tval any
        }

        type jNodeMapEntry struct {
        \tkey    uint32
        \tactive bool
        \tval    any
        }

        type jNodeMapState struct {
        \tentries []jNodeMapEntry
        \tnodeIdx map[uint32]int32
        }

        type jNodeMapState4 struct {
        \tstate jNodeMapState
        \tbuf   [4]jNodeMapEntry
        }

        type jNodeMapState8 struct {
        \tstate jNodeMapState
        \tbuf   [8]jNodeMapEntry
        }

        func newNodeMapState2(capHint int16) *jNodeMapState {
        \tif capHint > 8 {
        \t\treturn &jNodeMapState{entries: make([]jNodeMapEntry, 2, min(int(capHint), 256))}
        \t}
        \tif capHint >= 1 && capHint <= 4 {
        \t\ts4 := &jNodeMapState4{}
        \t\ts4.state.entries = s4.buf[:2]
        \t\treturn &s4.state
        \t}
        \ts8 := &jNodeMapState8{}
        \ts8.state.entries = s8.buf[:2]
        \treturn &s8.state
        }

        type jPtrMapEntry struct {
        \tkey any
        \tval any
        }

        type jPtrMapState struct {
        \tentries []jPtrMapEntry
        \tptrIdx  map[any]int32
        }

        type jPtrMapState4 struct {
        \tstate jPtrMapState
        \tbuf   [4]jPtrMapEntry
        }

        type jTreeNode struct {
        \tleft   int32
        \tright  int32
        \theight int8
        \tisStr  bool
        \tstrKey string
        \tkey    any
        \tval    any
        }

        type jSortedMapState struct {
        \tcmp      func(a, b any) int32
        \ttree     []jTreeNode
        \troot     int32
        \tfreeTree int32
        \tnatOrder bool
        }

        type JMapEntry struct {
        \tprev        int32
        \tnextAndKind uint32
        \thashNext    int32
        \tu32OrHash   uint32
        \tKey         any
        \tVal         any
        }

        func makeNextKind(n int32, k uint8) uint32 {
        \treturn (uint32(k) << 29) | (uint32(n) & jEntryPosMask)
        }

        func (e *JMapEntry) next() int32 {
        \treturn int32(e.nextAndKind & jEntryPosMask)
        }

        func (e *JMapEntry) setNext(n int32) {
        \te.nextAndKind = (e.nextAndKind &^ jEntryPosMask) | (uint32(n) & jEntryPosMask)
        }

        func (e *JMapEntry) kind() uint8 {
        \treturn uint8(e.nextAndKind >> 29)
        }

        type jGenMapIdx struct {
        \tstrIdx    map[string]int32
        \tu64Idx    map[uint64]int32
        \tptrIdx    map[any]int32
        \tbuckets   []int32
        \totherLive int32
        }

        type jGenMapState struct {
        \tentries  []JMapEntry
        \tidx      *jGenMapIdx
        \thead     int32
        \ttail     int32
        \tfreeHead int32
        \tnullPos  int32
        }

        type jGenMapState1 struct {
        \tstate  jGenMapState
        \tinline [1]JMapEntry
        }

        type jGenMapState2 struct {
        \tstate  jGenMapState
        \tinline [2]JMapEntry
        }

        type jGenMapState4 struct {
        \tstate  jGenMapState
        \tinline [4]JMapEntry
        }

        func newGenMapState(capHint int16) *jGenMapState {
        \tif capHint > 4 {
        \t\treturn &jGenMapState{entries: make([]JMapEntry, 0, min(int(capHint), 256))}
        \t}
        \tif capHint == 1 {
        \t\ts1 := &jGenMapState1{}
        \t\ts1.state.entries = s1.inline[:0]
        \t\treturn &s1.state
        \t}
        \tif capHint == 2 {
        \t\ts2 := &jGenMapState2{}
        \t\ts2.state.entries = s2.inline[:0]
        \t\treturn &s2.state
        \t}
        \ts4 := &jGenMapState4{}
        \ts4.state.entries = s4.inline[:0]
        \treturn &s4.state
        }

        func newGenMapStateSized(n int) *jGenMapState {
        \tswitch {
        \tcase n <= 1:
        \t\ts1 := &jGenMapState1{}
        \t\ts1.state.entries = s1.inline[:n]
        \t\treturn &s1.state
        \tcase n == 2:
        \t\ts2 := &jGenMapState2{}
        \t\ts2.state.entries = s2.inline[:n]
        \t\treturn &s2.state
        \tcase n <= 4:
        \t\ts4 := &jGenMapState4{}
        \t\ts4.state.entries = s4.inline[:n]
        \t\treturn &s4.state
        \tdefault:
        \t\treturn &jGenMapState{entries: make([]JMapEntry, n)}
        \t}
        }

        type JMapCore struct {
        \tstrEntries   []jStrMapEntry
        \text          any
        \tLive         int32
        \tholes        int32
        \tfirstIdx     int32
        \tcapHint      int16
        \tmode         uint8
        \tIdentity     bool
        \timmutable    bool
        \taccessOrd    bool
        \tgenNonCustom bool
        }

        type FastComputeMap interface {
        \tMapComputeIfAbsentStr_(string, Ext_java_util_function_Function) any
        \tMapComputeIfAbsentNode_(uint32, Ext_java_util_function_Function) any
        \tMapTryGetComputeStr_(string) (any, bool)
        \tMapTryGetComputeNode_(uint32) (any, bool)
        \tMapTryGetComputeAny_(any) (any, bool)
        }

        type FastStrMap interface {
        \tMapGetStr_(string) any
        \tMapGetOrDefaultStr_(string, any) any
        \tMapContainsStr_(string) bool
        \tMapPutStr_(string, any) any
        \tMapRemoveStr_(string) any
        }

        type FastStrSet interface {
        \tSetContainsStr_(string) bool
        \tSetAddStr_(string) bool
        \tSetRemoveStr_(string) bool
        }

        type FastNodeMap interface {
        \tMapGetNode_(uint32) any
        \tMapGetOrDefaultNode_(uint32, any) any
        \tMapContainsNode_(uint32) bool
        \tMapPutNode_(uint32, any) any
        \tMapRemoveNode_(uint32) any
        }

        type FastNodeSet interface {
        \tSetContainsNode_(uint32) bool
        \tSetAddNode_(uint32) bool
        \tSetRemoveNode_(uint32) bool
        }

        type FastMultimap interface {
        \tMultimapContainsKeyStr_(string) bool
        \tMultimapContainsKeyNode_(uint32) bool
        \tMultimapContainsEntryStr_(string, any) bool
        \tMultimapContainsEntryNode_(uint32, any) bool
        \tMultimapGetStr_(string) any
        \tMultimapGetNode_(uint32) any
        \tMultimapPutStr_(string, any) bool
        \tMultimapPutNode_(uint32, any) bool
        \tMultimapRemoveStr_(string, any) bool
        \tMultimapRemoveNode_(uint32, any) bool
        \tMultimapRemoveAllStr_(string) any
        \tMultimapRemoveAllNode_(uint32) any
        \tMultimapPutAllStr_(string, Ext_java_lang_Iterable) bool
        \tMultimapPutAllNode_(uint32, Ext_java_lang_Iterable) bool
        }

        type FastMultiset interface {
        \tMultisetCountStr_(string) int32
        \tMultisetContainsStr_(string) bool
        \tMultisetAddStr1_(string) bool
        \tMultisetAddStr2_(string, int32) int32
        \tMultisetRemoveStr1_(string) bool
        \tMultisetSetCountStr3_(string, int32, int32) bool
        }

        type NodePredicate interface {
        \tTestNode_(uint32) bool
        }

        type jFastDequeCoreProvider interface {
        \tfastDequeCore_() *JDequeCore
        }

        type jLinkedDequeCoreProvider interface {
        \tfastLinkedDequeCore_() *JListCore
        }

        func NewJMapCore(identity bool) *JMapCore {
        \treturn &JMapCore{Identity: identity, genNonCustom: identity}
        }

        func NewSortedJMapCore(cmp func(a, b any) int32) *JMapCore {
        \tm := &JMapCore{}
        \tif cmp == nil {
        \t\tm.mode = jMapModeSortedNat
        \t} else {
        \t\tm.mode = jMapModeSortedCmp
        \t\tm.ext = &jSortedMapState{cmp: cmp}
        \t}
        \treturn m
        }

        func (m *JMapCore) IsSorted() bool {
        \treturn m != nil && (m.mode == jMapModeSortedNat || m.mode == jMapModeSortedCmp)
        }

        func (m *JMapCore) SortedCmp() func(a, b any) int32 {
        \tif m == nil {
        \t\treturn nil
        \t}
        \tif m.mode == jMapModeSortedNat {
        \t\treturn CompareValues
        \t}
        \tif m.mode == jMapModeSortedCmp {
        \t\treturn m.ext.(*jSortedMapState).cmp
        \t}
        \treturn nil
        }

        func (m *JMapCore) IsNaturalOrder() bool {
        \treturn m != nil && m.mode == jMapModeSortedNat
        }

        func (m *JMapCore) ensureSortedState() *jSortedMapState {
        \tif m.ext != nil {
        \t\treturn m.ext.(*jSortedMapState)
        \t}
        \tst := &jSortedMapState{cmp: CompareValues, natOrder: true}
        \tm.ext = st
        \treturn st
        }

        func (m *JMapCore) InitExpectedSize(expected int32) {
        \tif expected < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException("Illegal initial capacity: " + fmt.Sprint(expected)))
        \t}
        \tif expected > 0 && (m.mode == jMapModeEmpty || m.mode == jMapModeNullOnly || m.mode == jMapModeNode1) {
        \t\tif expected > 4096 {
        \t\t\texpected = 4096
        \t\t}
        \t\tif int16(expected) > m.capHint {
        \t\t\tm.capHint = int16(expected)
        \t\t}
        \t}
        }

        func (m *JMapCore) InitAccessOrder(initialCap int32, loadFactor float32, accessOrder bool) {
        \tif initialCap < 0 {
        \t\tpanic(javalang.NewIllegalArgumentException("Illegal initial capacity: " + fmt.Sprint(initialCap)))
        \t}
        \tif loadFactor <= 0 || math.IsNaN(float64(loadFactor)) {
        \t\tpanic(javalang.NewIllegalArgumentException("Illegal load factor: " + fmt.Sprint(loadFactor)))
        \t}
        \tm.InitExpectedSize(initialCap)
        \tm.accessOrd = accessOrder
        }

        func (m *JMapCore) initCap() int {
        \tif m.capHint > 0 {
        \t\treturn min(int(m.capHint), 256)
        \t}
        \treturn 4
        }

        func (m *JMapCore) strInitCap() int {
        \tif m.capHint > 0 {
        \t\treturn min(int(m.capHint), 256)
        \t}
        \treturn 2
        }

        func (m *JMapCore) Size() int32   { return m.Live }
        func (m *JMapCore) IsEmpty() bool { return m.Live == 0 }

        // --- AVL Tree Helpers (m.IsSorted(); 1-based node ids, 0 = nil) ---

        func (s *jSortedMapState) treeNodeKey(id int32) any {
        \tn := &s.tree[id-1]
        \tif n.key != nil {
        \t\treturn n.key
        \t}
        \tif n.isStr {
        \t\treturn n.strKey
        \t}
        \treturn nil
        }

        func (s *jSortedMapState) avlHeight(id int32) int8 {
        \tif id == 0 {
        \t\treturn 0
        \t}
        \treturn s.tree[id-1].height
        }

        func (s *jSortedMapState) avlUpdateHeight(id int32) {
        \tlh := s.avlHeight(s.tree[id-1].left)
        \trh := s.avlHeight(s.tree[id-1].right)
        \tif lh > rh {
        \t\ts.tree[id-1].height = lh + 1
        \t} else {
        \t\ts.tree[id-1].height = rh + 1
        \t}
        }

        func (s *jSortedMapState) avlBalance(id int32) int8 {
        \treturn s.avlHeight(s.tree[id-1].left) - s.avlHeight(s.tree[id-1].right)
        }

        func (s *jSortedMapState) avlRotateRight(y int32) int32 {
        \tx := s.tree[y-1].left
        \tt2 := s.tree[x-1].right
        \ts.tree[x-1].right = y
        \ts.tree[y-1].left = t2
        \ts.avlUpdateHeight(y)
        \ts.avlUpdateHeight(x)
        \treturn x
        }

        func (s *jSortedMapState) avlRotateLeft(x int32) int32 {
        \ty := s.tree[x-1].right
        \tt2 := s.tree[y-1].left
        \ts.tree[y-1].left = x
        \ts.tree[x-1].right = t2
        \ts.avlUpdateHeight(x)
        \ts.avlUpdateHeight(y)
        \treturn y
        }

        func (s *jSortedMapState) avlRebalance(id int32) int32 {
        \ts.avlUpdateHeight(id)
        \tbf := s.avlBalance(id)
        \tif bf > 1 {
        \t\tif s.avlBalance(s.tree[id-1].left) < 0 {
        \t\t\ts.tree[id-1].left = s.avlRotateLeft(s.tree[id-1].left)
        \t\t}
        \t\treturn s.avlRotateRight(id)
        \t}
        \tif bf < -1 {
        \t\tif s.avlBalance(s.tree[id-1].right) > 0 {
        \t\t\ts.tree[id-1].right = s.avlRotateRight(s.tree[id-1].right)
        \t\t}
        \t\treturn s.avlRotateLeft(id)
        \t}
        \treturn id
        }

        func (s *jSortedMapState) avlAllocNode(isStr bool, strKey string, key, val any) int32 {
        \tif s.freeTree != 0 {
        \t\tid := s.freeTree
        \t\ts.freeTree = s.tree[id-1].left
        \t\ts.tree[id-1] = jTreeNode{height: 1, isStr: isStr, strKey: strKey, key: key, val: val}
        \t\treturn id
        \t}
        \tnode := jTreeNode{height: 1, isStr: isStr, strKey: strKey, key: key, val: val}
        \tif s.tree == nil {
        \t\ts.tree = make([]jTreeNode, 1, 4)
        \t\ts.tree[0] = node
        \t\treturn 1
        \t}
        \ts.tree = append(s.tree, node)
        \treturn int32(len(s.tree))
        }

        func (s *jSortedMapState) avlFreeNode(id int32) {
        \ts.tree[id-1] = jTreeNode{left: s.freeTree}
        \ts.freeTree = id
        }

        func (s *jSortedMapState) avlCmpStr(k string, id int32) int32 {
        \tn := &s.tree[id-1]
        \tif s.natOrder && n.isStr {
        \t\treturn int32(javalang.StringCompareTo(k, n.strKey))
        \t}
        \treturn s.cmp(k, s.treeNodeKey(id))
        }

        func (s *jSortedMapState) avlCmpAny(k any, id int32) int32 {
        \tn := &s.tree[id-1]
        \tif s.natOrder && n.isStr {
        \t\tif str, ok := k.(string); ok {
        \t\t\treturn int32(javalang.StringCompareTo(str, n.strKey))
        \t\t}
        \t}
        \treturn s.cmp(k, s.treeNodeKey(id))
        }

        func (s *jSortedMapState) avlFindStr(k string) int32 {
        \tcur := s.root
        \tfor cur != 0 {
        \t\tn := &s.tree[cur-1]
        \t\tvar c int32
        \t\tif s.natOrder && n.isStr {
        \t\t\tc = int32(javalang.StringCompareTo(k, n.strKey))
        \t\t} else {
        \t\t\tc = s.cmp(k, s.treeNodeKey(cur))
        \t\t}
        \t\tif c == 0 {
        \t\t\treturn cur
        \t\t}
        \t\tif c < 0 {
        \t\t\tcur = n.left
        \t\t} else {
        \t\t\tcur = n.right
        \t\t}
        \t}
        \treturn 0
        }

        func (s *jSortedMapState) avlFindAny(k any) int32 {
        \tcur := s.root
        \tfor cur != 0 {
        \t\tc := s.avlCmpAny(k, cur)
        \t\tif c == 0 {
        \t\t\treturn cur
        \t\t}
        \t\tif c < 0 {
        \t\t\tcur = s.tree[cur-1].left
        \t\t} else {
        \t\t\tcur = s.tree[cur-1].right
        \t\t}
        \t}
        \treturn 0
        }

        func (m *JMapCore) avlInsertStr(s *jSortedMapState, cur int32, k string, v any) (int32, any, bool) {
        \tif cur == 0 {
        \t\tif s.natOrder {
        \t\t\t_ = javalang.StringCompareTo(k, k)
        \t\t} else {
        \t\t\t_ = s.cmp(k, k)
        \t\t}
        \t\tid := s.avlAllocNode(true, k, nil, v)
        \t\tm.Live++
        \t\treturn id, nil, true
        \t}
        \tc := s.avlCmpStr(k, cur)
        \tif c == 0 {
        \t\told := s.tree[cur-1].val
        \t\ts.tree[cur-1].val = v
        \t\treturn cur, old, false
        \t}
        \tif c < 0 {
        \t\tnl, old, ins := m.avlInsertStr(s, s.tree[cur-1].left, k, v)
        \t\ts.tree[cur-1].left = nl
        \t\tif ins {
        \t\t\treturn s.avlRebalance(cur), nil, true
        \t\t}
        \t\treturn cur, old, false
        \t}
        \tnr, old, ins := m.avlInsertStr(s, s.tree[cur-1].right, k, v)
        \ts.tree[cur-1].right = nr
        \tif ins {
        \t\treturn s.avlRebalance(cur), nil, true
        \t}
        \treturn cur, old, false
        }

        func (m *JMapCore) avlInsertAny(s *jSortedMapState, cur int32, k, v any) (int32, any, bool) {
        \tif cur == 0 {
        \t\t_ = s.cmp(k, k)
        \t\tif str, ok := k.(string); ok {
        \t\t\tid := s.avlAllocNode(true, str, nil, v)
        \t\t\tm.Live++
        \t\t\treturn id, nil, true
        \t\t}
        \t\tid := s.avlAllocNode(false, "", k, v)
        \t\tm.Live++
        \t\treturn id, nil, true
        \t}
        \tc := s.avlCmpAny(k, cur)
        \tif c == 0 {
        \t\told := s.tree[cur-1].val
        \t\ts.tree[cur-1].val = v
        \t\treturn cur, old, false
        \t}
        \tif c < 0 {
        \t\tnl, old, ins := m.avlInsertAny(s, s.tree[cur-1].left, k, v)
        \t\ts.tree[cur-1].left = nl
        \t\tif ins {
        \t\t\treturn s.avlRebalance(cur), nil, true
        \t\t}
        \t\treturn cur, old, false
        \t}
        \tnr, old, ins := m.avlInsertAny(s, s.tree[cur-1].right, k, v)
        \ts.tree[cur-1].right = nr
        \tif ins {
        \t\treturn s.avlRebalance(cur), nil, true
        \t}
        \treturn cur, old, false
        }

        func (s *jSortedMapState) avlExtractMin(cur int32) (int32, int32) {
        \tif s.tree[cur-1].left == 0 {
        \t\treturn s.tree[cur-1].right, cur
        \t}
        \tnl, minId := s.avlExtractMin(s.tree[cur-1].left)
        \ts.tree[cur-1].left = nl
        \treturn s.avlRebalance(cur), minId
        }

        func (s *jSortedMapState) avlExtractMax(cur int32) (int32, int32) {
        \tif s.tree[cur-1].right == 0 {
        \t\treturn s.tree[cur-1].left, cur
        \t}
        \tnr, maxId := s.avlExtractMax(s.tree[cur-1].right)
        \ts.tree[cur-1].right = nr
        \treturn s.avlRebalance(cur), maxId
        }

        func (m *JMapCore) avlUnlinkRoot(s *jSortedMapState, cur int32) int32 {
        \tl := s.tree[cur-1].left
        \tr := s.tree[cur-1].right
        \ts.avlFreeNode(cur)
        \tm.Live--
        \tif l == 0 {
        \t\treturn r
        \t}
        \tif r == 0 {
        \t\treturn l
        \t}
        \tnr, succ := s.avlExtractMin(r)
        \ts.tree[succ-1].left = l
        \ts.tree[succ-1].right = nr
        \treturn s.avlRebalance(succ)
        }

        func (m *JMapCore) avlDeleteStr(s *jSortedMapState, cur int32, k string) (int32, any, bool) {
        \tif cur == 0 {
        \t\treturn 0, nil, false
        \t}
        \tc := s.avlCmpStr(k, cur)
        \tif c == 0 {
        \t\told := s.tree[cur-1].val
        \t\treturn m.avlUnlinkRoot(s, cur), old, true
        \t}
        \tif c < 0 {
        \t\tnl, old, del := m.avlDeleteStr(s, s.tree[cur-1].left, k)
        \t\tif !del {
        \t\t\treturn cur, nil, false
        \t\t}
        \t\ts.tree[cur-1].left = nl
        \t\treturn s.avlRebalance(cur), old, true
        \t}
        \tnr, old, del := m.avlDeleteStr(s, s.tree[cur-1].right, k)
        \tif !del {
        \t\treturn cur, nil, false
        \t}
        \ts.tree[cur-1].right = nr
        \treturn s.avlRebalance(cur), old, true
        }

        func (m *JMapCore) avlDeleteAny(s *jSortedMapState, cur int32, k any) (int32, any, bool) {
        \tif cur == 0 {
        \t\treturn 0, nil, false
        \t}
        \tc := s.avlCmpAny(k, cur)
        \tif c == 0 {
        \t\told := s.tree[cur-1].val
        \t\treturn m.avlUnlinkRoot(s, cur), old, true
        \t}
        \tif c < 0 {
        \t\tnl, old, del := m.avlDeleteAny(s, s.tree[cur-1].left, k)
        \t\tif !del {
        \t\t\treturn cur, nil, false
        \t\t}
        \t\ts.tree[cur-1].left = nl
        \t\treturn s.avlRebalance(cur), old, true
        \t}
        \tnr, old, del := m.avlDeleteAny(s, s.tree[cur-1].right, k)
        \tif !del {
        \t\treturn cur, nil, false
        \t}
        \ts.tree[cur-1].right = nr
        \treturn s.avlRebalance(cur), old, true
        }

        func (s *jSortedMapState) avlForEach(cur int32, fn func(k, v any) bool) bool {
        \tif cur == 0 {
        \t\treturn true
        \t}
        \tl := s.tree[cur-1].left
        \tif l != 0 && !s.avlForEach(l, fn) {
        \t\treturn false
        \t}
        \tk := s.treeNodeKey(cur)
        \tv := s.tree[cur-1].val
        \tr := s.tree[cur-1].right
        \tif !fn(k, v) {
        \t\treturn false
        \t}
        \tif r != 0 && !s.avlForEach(r, fn) {
        \t\treturn false
        \t}
        \treturn true
        }

        func (s *jSortedMapState) avlTransformValues(cur int32, fn func(k, v any) any) {
        \tif cur == 0 {
        \t\treturn
        \t}
        \ts.avlTransformValues(s.tree[cur-1].left, fn)
        \ts.tree[cur-1].val = ToAny(fn(s.treeNodeKey(cur), s.tree[cur-1].val))
        \ts.avlTransformValues(s.tree[cur-1].right, fn)
        }

        func (s *jSortedMapState) avlSetHash(cur int32) int32 {
        \tif cur == 0 {
        \t\treturn 0
        \t}
        \tn := &s.tree[cur-1]
        \tvar kh int32
        \tif n.isStr {
        \t\tkh = javalang.StringHashCode(n.strKey)
        \t} else {
        \t\tkh = elemJavaHashCode(n.key, false)
        \t}
        \treturn s.avlSetHash(n.left) + kh + s.avlSetHash(n.right)
        }

        // --- Homogeneous Mode Compaction & Promotion Helpers ---

        func (m *JMapCore) compactStr() {
        \tif m.holes == 0 {
        \t\tm.firstIdx = 0
        \t\treturn
        \t}
        \tvar idxMap map[string]int32
        \tif m.ext != nil {
        \t\tidxMap = m.ext.(map[string]int32)
        \t}
        \tw := 0
        \tfor i := int(m.firstIdx); i < len(m.strEntries); i++ {
        \t\te := m.strEntries[i]
        \t\tif e.key != jTombstoneStr {
        \t\t\tm.strEntries[w] = e
        \t\t\tif idxMap != nil {
        \t\t\t\tidxMap[e.key] = int32(w)
        \t\t\t}
        \t\t\tw++
        \t\t}
        \t}
        \tclear(m.strEntries[w:])
        \tm.strEntries = m.strEntries[:w]
        \tm.holes = 0
        \tm.firstIdx = 0
        }

        func (m *JMapCore) compactNode(ns *jNodeMapState) {
        \tif m.holes == 0 {
        \t\tm.firstIdx = 0
        \t\treturn
        \t}
        \tw := 0
        \tfor i := int(m.firstIdx); i < len(ns.entries); i++ {
        \t\te := ns.entries[i]
        \t\tif e.active {
        \t\t\tns.entries[w] = e
        \t\t\tif ns.nodeIdx != nil {
        \t\t\t\tns.nodeIdx[e.key] = int32(w)
        \t\t\t}
        \t\t\tw++
        \t\t}
        \t}
        \tclear(ns.entries[w:])
        \tns.entries = ns.entries[:w]
        \tm.holes = 0
        \tm.firstIdx = 0
        }

        func (m *JMapCore) compactPtr(ps *jPtrMapState) {
        \tif m.holes == 0 {
        \t\tm.firstIdx = 0
        \t\treturn
        \t}
        \tw := 0
        \tfor i := int(m.firstIdx); i < len(ps.entries); i++ {
        \t\te := ps.entries[i]
        \t\tif e.key != jTombstoneObj {
        \t\t\tps.entries[w] = e
        \t\t\tif ps.ptrIdx != nil {
        \t\t\t\tps.ptrIdx[e.key] = int32(w)
        \t\t\t}
        \t\t\tw++
        \t\t}
        \t}
        \tclear(ps.entries[w:])
        \tps.entries = ps.entries[:w]
        \tm.holes = 0
        \tm.firstIdx = 0
        }

        func (m *JMapCore) removeStrAt(idx int) any {
        \te := &m.strEntries[idx]
        \tif e.key == jTombstoneStr {
        \t\treturn nil
        \t}
        \told := e.val
        \tif m.ext != nil {
        \t\tdelete(m.ext.(map[string]int32), e.key)
        \t}
        \t*e = jStrMapEntry{key: jTombstoneStr}
        \tm.Live--
        \tif m.Live == 0 {
        \t\tclear(m.strEntries)
        \t\tm.strEntries = m.strEntries[:0]
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t} else if idx == len(m.strEntries)-1 {
        \t\tm.strEntries = m.strEntries[:idx]
        \t} else {
        \t\tm.holes++
        \t\tif idx == int(m.firstIdx) {
        \t\t\tm.firstIdx++
        \t\t}
        \t}
        \treturn old
        }

        func (m *JMapCore) removeNodeAt(idx int) any {
        \tns := m.ext.(*jNodeMapState)
        \te := &ns.entries[idx]
        \tif !e.active {
        \t\treturn nil
        \t}
        \told := e.val
        \tif ns.nodeIdx != nil {
        \t\tdelete(ns.nodeIdx, e.key)
        \t}
        \t*e = jNodeMapEntry{}
        \tm.Live--
        \tif m.Live == 0 {
        \t\tclear(ns.entries)
        \t\tns.entries = ns.entries[:0]
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t} else if idx == len(ns.entries)-1 {
        \t\tns.entries = ns.entries[:idx]
        \t} else {
        \t\tm.holes++
        \t\tif idx == int(m.firstIdx) {
        \t\t\tm.firstIdx++
        \t\t}
        \t}
        \treturn old
        }

        func (m *JMapCore) removePtrAt(idx int) any {
        \tps := m.ext.(*jPtrMapState)
        \te := &ps.entries[idx]
        \tif e.key == jTombstoneObj {
        \t\treturn nil
        \t}
        \told := e.val
        \tif ps.ptrIdx != nil {
        \t\tdelete(ps.ptrIdx, e.key)
        \t}
        \t*e = jPtrMapEntry{key: jTombstoneObj}
        \tm.Live--
        \tif m.Live == 0 {
        \t\tclear(ps.entries)
        \t\tps.entries = ps.entries[:0]
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t} else if idx == len(ps.entries)-1 {
        \t\tps.entries = ps.entries[:idx]
        \t} else {
        \t\tm.holes++
        \t\tif idx == int(m.firstIdx) {
        \t\t\tm.firstIdx++
        \t\t}
        \t}
        \treturn old
        }
        """;
  }

  private String buildCoreRuntimeGoPart2Gen() {
    return """

        // --- General / Mixed / LRU Access-Order State on *JMapCore ---

        func (g *jGenMapState) ensureIdx() *jGenMapIdx {
        \tif g.idx == nil {
        \t\tg.idx = &jGenMapIdx{}
        \t}
        \treturn g.idx
        }

        func (g *jGenMapState) entryKey(pos int32) any {
        \te := &g.entries[pos-1]
        \tif e.Key != nil {
        \t\treturn e.Key
        \t}
        \tswitch e.kind() {
        \tcase jKeyNull:
        \t\treturn nil
        \tcase jKeyNode:
        \t\treturn boxNodeById(e.u32OrHash)
        \tcase jKeyInt:
        \t\treturn ValueOfInt(int32(e.u32OrHash))
        \tdefault:
        \t\treturn e.Key
        \t}
        }

        func (m *JMapCore) allocGenEntry(g *jGenMapState, entry JMapEntry) int32 {
        \tvar pos int32
        \tif g.freeHead != 0 {
        \t\tpos = g.freeHead
        \t\tg.freeHead = g.entries[pos-1].next()
        \t\tg.entries[pos-1] = entry
        \t} else {
        \t\tif len(g.entries) == cap(g.entries) {
        \t\t\tnewCap := cap(g.entries) * 2
        \t\t\tif newCap < 4 {
        \t\t\t\tnewCap = 4
        \t\t\t}
        \t\t\tif hint := min(int(m.capHint), 256); hint > newCap {
        \t\t\t\tnewCap = hint
        \t\t\t}
        \t\t\tne := make([]JMapEntry, len(g.entries)+1, newCap)
        \t\t\tcopy(ne, g.entries)
        \t\t\tne[len(g.entries)] = entry
        \t\t\tg.entries = ne
        \t\t} else {
        \t\t\tg.entries = append(g.entries, entry)
        \t\t}
        \t\tpos = int32(len(g.entries))
        \t}
        \tg.entries[pos-1].prev = g.tail
        \tg.entries[pos-1].setNext(0)
        \tif g.tail != 0 {
        \t\tg.entries[g.tail-1].setNext(pos)
        \t} else {
        \t\tg.head = pos
        \t}
        \tg.tail = pos
        \tm.Live++
        \treturn pos
        }

        func (g *jGenMapState) moveToTail(pos int32) {
        \tif pos == 0 || pos == g.tail {
        \t\treturn
        \t}
        \tprev := g.entries[pos-1].prev
        \tnext := g.entries[pos-1].next()
        \tif prev != 0 {
        \t\tg.entries[prev-1].setNext(next)
        \t} else {
        \t\tg.head = next
        \t}
        \tg.entries[next-1].prev = prev
        \tg.entries[pos-1].prev = g.tail
        \tg.entries[pos-1].setNext(0)
        \tg.entries[g.tail-1].setNext(pos)
        \tg.tail = pos
        }

        func (g *jGenMapState) moveTailToFront() {
        \tpos := g.tail
        \tif pos == 0 || pos == g.head {
        \t\treturn
        \t}
        \tprev := g.entries[pos-1].prev
        \tg.entries[prev-1].setNext(0)
        \tg.tail = prev
        \tg.entries[pos-1].prev = 0
        \tg.entries[pos-1].setNext(g.head)
        \tg.entries[g.head-1].prev = pos
        \tg.head = pos
        }

        func (m *JMapCore) hashGeneral(k any) int32 {
        \treturn elemJavaHashCode(k, m.Identity)
        }

        func (m *JMapCore) eqGeneral(a, b any) bool {
        \tif m.Identity {
        \t\treturn RefEquals(a, b)
        \t}
        \treturn javalang.ObjectsEquals(a, b)
        }

        func (m *JMapCore) materializeGenBuckets(g *jGenMapState) *jGenMapIdx {
        \tidx := g.ensureIdx()
        \tif len(idx.buckets) == 0 {
        \t\tcapB := 16
        \t\tfor capB*3 <= int(m.capHint)*4 && capB < 2048 {
        \t\t\tcapB <<= 1
        \t\t}
        \t\tidx.buckets = make([]int32, capB)
        \t\tmask := uint32(capB - 1)
        \t\tvar cnt int32
        \t\tfor pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
        \t\t\te := &g.entries[pos-1]
        \t\t\tif e.kind() == jKeyGeneral {
        \t\t\t\tslot := smearHash(int32(e.u32OrHash)) & mask
        \t\t\t\te.hashNext = idx.buckets[slot]
        \t\t\t\tidx.buckets[slot] = pos
        \t\t\t\tcnt++
        \t\t\t}
        \t\t}
        \t\tidx.otherLive = cnt
        \t\treturn idx
        \t}
        \tif int(idx.otherLive)*4 < len(idx.buckets)*3 {
        \t\treturn idx
        \t}
        \tnewCap := len(idx.buckets) * 2
        \tnb := make([]int32, newCap)
        \tmask := uint32(newCap - 1)
        \tfor _, headPos := range idx.buckets {
        \t\tfor pos := headPos; pos != 0; {
        \t\t\te := &g.entries[pos-1]
        \t\t\tnextHash := e.hashNext
        \t\t\tslot := smearHash(int32(e.u32OrHash)) & mask
        \t\t\te.hashNext = nb[slot]
        \t\t\tnb[slot] = pos
        \t\t\tpos = nextHash
        \t\t}
        \t}
        \tidx.buckets = nb
        \treturn idx
        }

        func findGenFastCustom(g *jGenMapState, k any, fk FastCustomKey, h int32) int32 {
        \tif idx := g.idx; idx != nil && len(idx.buckets) > 0 {
        \t\tentries := g.entries
        \t\tslot := smearHash(h) & uint32(len(idx.buckets)-1)
        \t\tfor pos := idx.buckets[slot]; pos != 0; {
        \t\t\te := &entries[pos-1]
        \t\t\tif int32(e.u32OrHash) == h && (e.Key == k || fk.Equals(e.Key)) {
        \t\t\t\treturn pos
        \t\t\t}
        \t\t\tpos = e.hashNext
        \t\t}
        \t\treturn 0
        \t}
        \tentries := g.entries
        \tuh := uint32(h)
        \tfor i := range entries {
        \t\te := &entries[i]
        \t\tif e.u32OrHash == uh && e.kind() == jKeyGeneral && (e.Key == k || fk.Equals(e.Key)) {
        \t\t\treturn int32(i + 1)
        \t\t}
        \t}
        \treturn 0
        }

        func (m *JMapCore) findGenPos(g *jGenMapState, k any, h int32) int32 {
        \tif idx := g.idx; idx != nil && len(idx.buckets) > 0 {
        \t\tentries := g.entries
        \t\tslot := smearHash(h) & uint32(len(idx.buckets)-1)
        \t\tfor pos := idx.buckets[slot]; pos != 0; {
        \t\t\te := &entries[pos-1]
        \t\t\tif int32(e.u32OrHash) == h && m.eqGeneral(e.Key, k) {
        \t\t\t\treturn pos
        \t\t\t}
        \t\t\tpos = e.hashNext
        \t\t}
        \t\treturn 0
        \t}
        \tentries := g.entries
        \tuh := uint32(h)
        \tfor i := range entries {
        \t\te := &entries[i]
        \t\tif e.u32OrHash == uh && e.kind() == jKeyGeneral && m.eqGeneral(e.Key, k) {
        \t\t\treturn int32(i + 1)
        \t\t}
        \t}
        \treturn 0
        }

        func (g *jGenMapState) unlinkGenBucket(pos int32, h int32) {
        \tif g.idx == nil || len(g.idx.buckets) == 0 {
        \t\treturn
        \t}
        \tslot := smearHash(h) & uint32(len(g.idx.buckets)-1)
        \tvar prev int32
        \tfor cur := g.idx.buckets[slot]; cur != 0; cur = g.entries[cur-1].hashNext {
        \t\tif cur == pos {
        \t\t\tnext := g.entries[cur-1].hashNext
        \t\t\tif prev == 0 {
        \t\t\t\tg.idx.buckets[slot] = next
        \t\t\t} else {
        \t\t\t\tg.entries[prev-1].hashNext = next
        \t\t\t}
        \t\t\tg.idx.otherLive--
        \t\t\treturn
        \t\t}
        \t\tprev = cur
        \t}
        }

        func (m *JMapCore) removeGenPos(pos int32) any {
        \tg := m.ext.(*jGenMapState)
        \te := &g.entries[pos-1]
        \told := e.Val
        \tswitch e.kind() {
        \tcase jKeyNull:
        \t\tg.nullPos = 0
        \tcase jKeyStr:
        \t\tif g.idx != nil && g.idx.strIdx != nil {
        \t\t\tdelete(g.idx.strIdx, e.Key.(string))
        \t\t}
        \tcase jKeyNode:
        \t\tif g.idx != nil && g.idx.u64Idx != nil {
        \t\t\tdelete(g.idx.u64Idx, (uint64(1)<<32)|uint64(e.u32OrHash))
        \t\t}
        \tcase jKeyInt:
        \t\tif g.idx != nil && g.idx.u64Idx != nil {
        \t\t\tdelete(g.idx.u64Idx, (uint64(2)<<32)|uint64(e.u32OrHash))
        \t\t}
        \tcase jKeyPtr:
        \t\tif g.idx != nil && g.idx.ptrIdx != nil {
        \t\t\tdelete(g.idx.ptrIdx, e.Key)
        \t\t}
        \tdefault:
        \t\tg.unlinkGenBucket(pos, int32(e.u32OrHash))
        \t}
        \tprev, next := e.prev, e.next()
        \tif prev != 0 {
        \t\tg.entries[prev-1].setNext(next)
        \t} else {
        \t\tg.head = next
        \t}
        \tif next != 0 {
        \t\tg.entries[next-1].prev = prev
        \t} else {
        \t\tg.tail = prev
        \t}
        \tg.entries[pos-1] = JMapEntry{nextAndKind: makeNextKind(g.freeHead, jKeyFree)}
        \tg.freeHead = pos
        \tm.Live--
        \tif m.Live == 0 {
        \t\tg.head = 0
        \t\tg.tail = 0
        \t\tg.freeHead = 0
        \t\tg.nullPos = 0
        \t\tclear(g.entries)
        \t\tg.entries = g.entries[:0]
        \t}
        \treturn old
        }

        func (g *jGenMapState) materializeStrIdx() *jGenMapIdx {
        \tidx := g.ensureIdx()
        \tidx.strIdx = make(map[string]int32, 16)
        \tfor pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
        \t\te := &g.entries[pos-1]
        \t\tif e.kind() == jKeyStr {
        \t\t\tidx.strIdx[e.Key.(string)] = pos
        \t\t}
        \t}
        \treturn idx
        }

        func (g *jGenMapState) materializeU64Idx() *jGenMapIdx {
        \tidx := g.ensureIdx()
        \tidx.u64Idx = make(map[uint64]int32, 16)
        \tfor pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
        \t\te := &g.entries[pos-1]
        \t\tswitch e.kind() {
        \t\tcase jKeyNode:
        \t\t\tidx.u64Idx[(uint64(1)<<32)|uint64(e.u32OrHash)] = pos
        \t\tcase jKeyInt:
        \t\t\tidx.u64Idx[(uint64(2)<<32)|uint64(e.u32OrHash)] = pos
        \t\t}
        \t}
        \treturn idx
        }

        func (g *jGenMapState) materializePtrIdx() *jGenMapIdx {
        \tidx := g.ensureIdx()
        \tidx.ptrIdx = make(map[any]int32, 16)
        \tfor pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
        \t\te := &g.entries[pos-1]
        \t\tif e.kind() == jKeyPtr {
        \t\t\tidx.ptrIdx[e.Key] = pos
        \t\t}
        \t}
        \treturn idx
        }

        func (m *JMapCore) putGenNull(g *jGenMapState, v any) (any, bool) {
        \tif g.nullPos != 0 {
        \t\told := g.entries[g.nullPos-1].Val
        \t\tg.entries[g.nullPos-1].Val = v
        \t\tif m.accessOrd {
        \t\t\tg.moveToTail(g.nullPos)
        \t\t}
        \t\treturn old, false
        \t}
        \tg.nullPos = m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyNull), Val: v})
        \treturn nil, true
        }

        func (m *JMapCore) putGenStr(g *jGenMapState, k string, boxedKey any, v any) (any, bool) {
        \tm.genNonCustom = true
        \tif !m.Identity {
        \t\tif boxedKey == nil {
        \t\t\tboxedKey = k
        \t\t}
        \t\tif g.idx != nil && g.idx.strIdx != nil {
        \t\t\tif pos := g.idx.strIdx[k]; pos != 0 {
        \t\t\t\told := g.entries[pos-1].Val
        \t\t\t\tg.entries[pos-1].Val = v
        \t\t\t\tif m.accessOrd {
        \t\t\t\t\tg.moveToTail(pos)
        \t\t\t\t}
        \t\t\t\treturn old, false
        \t\t\t}
        \t\t\tpos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyStr), Key: boxedKey, Val: v})
        \t\t\tg.idx.strIdx[k] = pos
        \t\t\treturn nil, true
        \t\t}
        \t\tfor pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
        \t\t\te := &g.entries[pos-1]
        \t\t\tif e.kind() == jKeyStr && e.Key.(string) == k {
        \t\t\t\told := e.Val
        \t\t\t\te.Val = v
        \t\t\t\tif m.accessOrd {
        \t\t\t\t\tg.moveToTail(pos)
        \t\t\t\t}
        \t\t\t\treturn old, false
        \t\t\t}
        \t\t}
        \t\tif int(m.Live) < jSmallLimit {
        \t\t\tm.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyStr), Key: boxedKey, Val: v})
        \t\t\treturn nil, true
        \t\t}
        \t\tidx := g.materializeStrIdx()
        \t\tpos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyStr), Key: boxedKey, Val: v})
        \t\tidx.strIdx[k] = pos
        \t\treturn nil, true
        \t}
        \tif boxedKey != nil {
        \t\treturn m.putGenOther(g, boxedKey, v)
        \t}
        \treturn m.putGenOther(g, k, v)
        }

        func (m *JMapCore) putGenNode(g *jGenMapState, u uint32, v any) (any, bool) {
        \tm.genNonCustom = true
        \tif g.idx != nil && g.idx.u64Idx != nil {
        \t\tuk := (uint64(1) << 32) | uint64(u)
        \t\tif pos := g.idx.u64Idx[uk]; pos != 0 {
        \t\t\told := g.entries[pos-1].Val
        \t\t\tg.entries[pos-1].Val = v
        \t\t\tif m.accessOrd {
        \t\t\t\tg.moveToTail(pos)
        \t\t\t}
        \t\t\treturn old, false
        \t\t}
        \t\tpos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyNode), u32OrHash: u, Val: v})
        \t\tg.idx.u64Idx[uk] = pos
        \t\treturn nil, true
        \t}
        \tfor pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
        \t\te := &g.entries[pos-1]
        \t\tif e.kind() == jKeyNode && e.u32OrHash == u {
        \t\t\told := e.Val
        \t\t\te.Val = v
        \t\t\tif m.accessOrd {
        \t\t\t\tg.moveToTail(pos)
        \t\t\t}
        \t\t\treturn old, false
        \t\t}
        \t}
        \tif int(m.Live) < jSmallLimit {
        \t\tm.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyNode), u32OrHash: u, Val: v})
        \t\treturn nil, true
        \t}
        \tidx := g.materializeU64Idx()
        \tpos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyNode), u32OrHash: u, Val: v})
        \tidx.u64Idx[(uint64(1)<<32)|uint64(u)] = pos
        \treturn nil, true
        }

        func (m *JMapCore) putGenInt32(g *jGenMapState, iv int32, boxedKey any, v any) (any, bool) {
        \tm.genNonCustom = true
        \tu := uint32(iv)
        \tif g.idx != nil && g.idx.u64Idx != nil {
        \t\tuk := (uint64(2) << 32) | uint64(u)
        \t\tif pos := g.idx.u64Idx[uk]; pos != 0 {
        \t\t\told := g.entries[pos-1].Val
        \t\t\tg.entries[pos-1].Val = v
        \t\t\tif m.accessOrd {
        \t\t\t\tg.moveToTail(pos)
        \t\t\t}
        \t\t\treturn old, false
        \t\t}
        \t\tpos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyInt), u32OrHash: u, Key: boxedKey, Val: v})
        \t\tg.idx.u64Idx[uk] = pos
        \t\treturn nil, true
        \t}
        \tfor pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
        \t\te := &g.entries[pos-1]
        \t\tif e.kind() == jKeyInt && e.u32OrHash == u {
        \t\t\told := e.Val
        \t\t\te.Val = v
        \t\t\tif m.accessOrd {
        \t\t\t\tg.moveToTail(pos)
        \t\t\t}
        \t\t\treturn old, false
        \t\t}
        \t}
        \tif int(m.Live) < jSmallLimit {
        \t\tm.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyInt), u32OrHash: u, Key: boxedKey, Val: v})
        \t\treturn nil, true
        \t}
        \tidx := g.materializeU64Idx()
        \tpos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyInt), u32OrHash: u, Key: boxedKey, Val: v})
        \tidx.u64Idx[(uint64(2)<<32)|uint64(u)] = pos
        \treturn nil, true
        }

        func (m *JMapCore) putGenPtr(g *jGenMapState, k, v any) (any, bool) {
        \tm.genNonCustom = true
        \tif g.idx != nil && g.idx.ptrIdx != nil {
        \t\tif pos := g.idx.ptrIdx[k]; pos != 0 {
        \t\t\told := g.entries[pos-1].Val
        \t\t\tg.entries[pos-1].Val = v
        \t\t\tif m.accessOrd {
        \t\t\t\tg.moveToTail(pos)
        \t\t\t}
        \t\t\treturn old, false
        \t\t}
        \t\tpos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyPtr), Key: k, Val: v})
        \t\tg.idx.ptrIdx[k] = pos
        \t\treturn nil, true
        \t}
        \tfor pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
        \t\te := &g.entries[pos-1]
        \t\tif e.kind() == jKeyPtr && e.Key == k {
        \t\t\told := e.Val
        \t\t\te.Val = v
        \t\t\tif m.accessOrd {
        \t\t\t\tg.moveToTail(pos)
        \t\t\t}
        \t\t\treturn old, false
        \t\t}
        \t}
        \tif int(m.Live) < jSmallLimit {
        \t\tm.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyPtr), Key: k, Val: v})
        \t\treturn nil, true
        \t}
        \tidx := g.materializePtrIdx()
        \tpos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyPtr), Key: k, Val: v})
        \tidx.ptrIdx[k] = pos
        \treturn nil, true
        }

        func (m *JMapCore) putGenFastCustom(g *jGenMapState, k any, fk FastCustomKey, v any) (any, bool) {
        \th := fk.HashCode()
        \tif pos := findGenFastCustom(g, k, fk, h); pos != 0 {
        \t\told := g.entries[pos-1].Val
        \t\tg.entries[pos-1].Val = v
        \t\tif m.accessOrd {
        \t\t\tg.moveToTail(pos)
        \t\t}
        \t\treturn old, false
        \t}
        \tif int(m.Live) < jSmallLimit && (g.idx == nil || len(g.idx.buckets) == 0) {
        \t\tm.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyGeneral), u32OrHash: uint32(h), Key: k, Val: v})
        \t\treturn nil, true
        \t}
        \tidx := m.materializeGenBuckets(g)
        \tslot := smearHash(h) & uint32(len(idx.buckets)-1)
        \tpos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyGeneral), u32OrHash: uint32(h), hashNext: idx.buckets[slot], Key: k, Val: v})
        \tidx.buckets[slot] = pos
        \tidx.otherLive++
        \treturn nil, true
        }

        func (m *JMapCore) putGenOther(g *jGenMapState, k, v any) (any, bool) {
        \th := m.hashGeneral(k)
        \tif pos := m.findGenPos(g, k, h); pos != 0 {
        \t\told := g.entries[pos-1].Val
        \t\tg.entries[pos-1].Val = v
        \t\tif m.accessOrd {
        \t\t\tg.moveToTail(pos)
        \t\t}
        \t\treturn old, false
        \t}
        \tif int(m.Live) < jSmallLimit && (g.idx == nil || len(g.idx.buckets) == 0) {
        \t\tm.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyGeneral), u32OrHash: uint32(h), Key: k, Val: v})
        \t\treturn nil, true
        \t}
        \tidx := m.materializeGenBuckets(g)
        \tslot := smearHash(h) & uint32(len(idx.buckets)-1)
        \tpos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyGeneral), u32OrHash: uint32(h), hashNext: idx.buckets[slot], Key: k, Val: v})
        \tidx.buckets[slot] = pos
        \tidx.otherLive++
        \treturn nil, true
        }

        func newGenMapStateCap(need int) *jGenMapState {
        \tif need <= 4 {
        \t\ts4 := &jGenMapState4{}
        \t\ts4.state.entries = s4.inline[:0]
        \t\treturn &s4.state
        \t}
        \treturn &jGenMapState{entries: make([]JMapEntry, 0, need)}
        }

        func (m *JMapCore) upgradeToGeneral() *jGenMapState {
        \tswitch m.mode {
        \tcase jMapModeGeneral:
        \t\treturn m.ext.(*jGenMapState)
        \tcase jMapModeEmpty:
        \t\tg := newGenMapState(m.capHint)
        \t\tm.ext = g
        \t\tm.mode = jMapModeGeneral
        \t\tm.firstIdx = 0
        \t\tm.genNonCustom = m.Identity
        \t\treturn g
        \tcase jMapModeNullOnly:
        \t\tnullVal := m.ext
        \t\tc := m.capHint
        \t\tif c > 0 && c < 2 {
        \t\t\tc = 2
        \t\t}
        \t\tg := newGenMapState(c)
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\tm.Live = 0
        \t\tm.ext = g
        \t\tm.mode = jMapModeGeneral
        \t\tm.genNonCustom = m.Identity
        \t\tm.putGenNull(g, nullVal)
        \t\treturn g
        \tcase jMapModeNode1:
        \t\tu1 := uint32(m.firstIdx)
        \t\tv1 := m.ext
        \t\tc := m.capHint
        \t\tif c > 0 && c < 2 {
        \t\t\tc = 2
        \t\t}
        \t\tg := newGenMapState(c)
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\tm.Live = 0
        \t\tm.ext = g
        \t\tm.mode = jMapModeGeneral
        \t\tm.genNonCustom = m.Identity
        \t\tm.putGenNode(g, u1, v1)
        \t\treturn g
        \tcase jMapModeStr:
        \t\told := m.strEntries
        \t\tstart := int(m.firstIdx)
        \t\tg := newGenMapStateCap(max(int(m.Live)+2, m.initCap()))
        \t\tm.strEntries = nil
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\tm.Live = 0
        \t\tm.ext = g
        \t\tm.mode = jMapModeGeneral
        \t\tm.genNonCustom = m.Identity
        \t\tfor i := start; i < len(old); i++ {
        \t\t\te := old[i]
        \t\t\tif e.key == jTombstoneStr {
        \t\t\t\tcontinue
        \t\t\t}
        \t\t\tif e.key == NullString {
        \t\t\t\tm.putGenNull(g, e.val)
        \t\t\t} else {
        \t\t\t\tm.putGenStr(g, e.key, nil, e.val)
        \t\t\t}
        \t\t}
        \t\treturn g
        \tcase jMapModeNode:
        \t\tns := m.ext.(*jNodeMapState)
        \t\tstart := int(m.firstIdx)
        \t\tg := newGenMapStateCap(max(int(m.Live)+2, m.initCap()))
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\tm.Live = 0
        \t\tm.ext = g
        \t\tm.mode = jMapModeGeneral
        \t\tm.genNonCustom = m.Identity
        \t\tfor i := start; i < len(ns.entries); i++ {
        \t\t\te := ns.entries[i]
        \t\t\tif !e.active {
        \t\t\t\tcontinue
        \t\t\t}
        \t\t\tif e.key == 0 {
        \t\t\t\tm.putGenNull(g, e.val)
        \t\t\t} else {
        \t\t\t\tm.putGenNode(g, e.key, e.val)
        \t\t\t}
        \t\t}
        \t\treturn g
        \tcase jMapModePtr:
        \t\tps := m.ext.(*jPtrMapState)
        \t\tstart := int(m.firstIdx)
        \t\tg := newGenMapStateCap(max(int(m.Live)+2, m.initCap()))
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\tm.Live = 0
        \t\tm.ext = g
        \t\tm.mode = jMapModeGeneral
        \t\tm.genNonCustom = m.Identity
        \t\tfor i := start; i < len(ps.entries); i++ {
        \t\t\te := ps.entries[i]
        \t\t\tif e.key == jTombstoneObj {
        \t\t\t\tcontinue
        \t\t\t}
        \t\t\tif e.key == nil {
        \t\t\t\tm.putGenNull(g, e.val)
        \t\t\t} else {
        \t\t\t\tm.putGenPtr(g, e.key, e.val)
        \t\t\t}
        \t\t}
        \t\treturn g
        \t}
        \tpanic("upgradeToGeneral called on sorted JMapCore")
        }

        // --- Fast Typed String, Node, Ptr & General Dispatch on *JMapCore ---

        func (m *JMapCore) getGenNull(g *jGenMapState, touch bool) (any, bool) {
        \tif g.nullPos != 0 {
        \t\tif touch && m.accessOrd {
        \t\t\tg.moveToTail(g.nullPos)
        \t\t}
        \t\treturn g.entries[g.nullPos-1].Val, true
        \t}
        \treturn nil, false
        }

        func (m *JMapCore) getStr(k string, touch bool) (any, bool) {
        \tif m.Live == 0 {
        \t\treturn nil, false
        \t}
        \tswitch m.mode {
        \tcase jMapModeStr:
        \t\tif k == jTombstoneStr {
        \t\t\treturn nil, false
        \t\t}
        \t\tif m.ext != nil {
        \t\t\tif idx, ok := m.ext.(map[string]int32)[k]; ok {
        \t\t\t\treturn m.strEntries[idx].val, true
        \t\t\t}
        \t\t\treturn nil, false
        \t\t}
        \t\tfor i := range m.strEntries {
        \t\t\tif m.strEntries[i].key == k {
        \t\t\t\treturn m.strEntries[i].val, true
        \t\t\t}
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeNode:
        \t\tif k == NullString {
        \t\t\treturn m.getNode(0, touch)
        \t\t}
        \t\treturn nil, false
        \tcase jMapModePtr:
        \t\tif k == NullString {
        \t\t\treturn m.getPtr(nil)
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeNullOnly:
        \t\tif k == NullString {
        \t\t\treturn m.ext, true
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeNode1:
        \t\treturn nil, false
        \tcase jMapModeSortedNat, jMapModeSortedCmp:
        \t\tif k == NullString {
        \t\t\treturn m.getSortedNull()
        \t\t}
        \t\tst := m.ext.(*jSortedMapState)
        \t\tid := st.avlFindStr(k)
        \t\tif id == 0 {
        \t\t\treturn nil, false
        \t\t}
        \t\treturn st.tree[id-1].val, true
        \tcase jMapModeGeneral:
        \t\tg := m.ext.(*jGenMapState)
        \t\tif k == NullString {
        \t\t\treturn m.getGenNull(g, touch)
        \t\t}
        \t\tif !m.Identity {
        \t\t\tif g.idx != nil && g.idx.strIdx != nil {
        \t\t\t\tif pos := g.idx.strIdx[k]; pos != 0 {
        \t\t\t\t\tif touch && m.accessOrd {
        \t\t\t\t\t\tg.moveToTail(pos)
        \t\t\t\t\t}
        \t\t\t\t\treturn g.entries[pos-1].Val, true
        \t\t\t\t}
        \t\t\t\treturn nil, false
        \t\t\t}
        \t\t\tfor pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
        \t\t\t\te := &g.entries[pos-1]
        \t\t\t\tif e.kind() == jKeyStr && e.Key.(string) == k {
        \t\t\t\t\tif touch && m.accessOrd {
        \t\t\t\t\t\tg.moveToTail(pos)
        \t\t\t\t\t}
        \t\t\t\t\treturn e.Val, true
        \t\t\t\t}
        \t\t\t}
        \t\t\treturn nil, false
        \t\t}
        \t\th := m.hashGeneral(k)
        \t\tif pos := m.findGenPos(g, k, h); pos != 0 {
        \t\t\tif touch && m.accessOrd {
        \t\t\t\tg.moveToTail(pos)
        \t\t\t}
        \t\t\treturn g.entries[pos-1].Val, true
        \t\t}
        \t\treturn nil, false
        \t}
        \treturn nil, false
        }

        func (m *JMapCore) getSortedNull() (any, bool) {
        \tif m.mode == jMapModeSortedCmp && m.ext != nil && m.Live > 0 {
        \t\tst := m.ext.(*jSortedMapState)
        \t\tid := st.avlFindAny(nil)
        \t\tif id != 0 {
        \t\t\treturn st.tree[id-1].val, true
        \t\t}
        \t}
        \treturn nil, false
        }

        func (m *JMapCore) ContainsStr(k string) bool {
        \t_, ok := m.getStr(k, false)
        \treturn ok
        }

        func (m *JMapCore) GetStr(k string) any {
        \tv, _ := m.getStr(k, true)
        \treturn v
        }

        func (m *JMapCore) GetOrDefaultStr(k string, def any) any {
        \tif v, ok := m.getStr(k, true); ok {
        \t\treturn v
        \t}
        \treturn ToAny(def)
        }

        func (m *JMapCore) putStrBoxed(k string, boxedKey any, v any) (any, bool) {
        \tif m.immutable {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tv = ToAny(v)
        \tswitch m.mode {
        \tcase jMapModeEmpty:
        \t\tif k == NullString && !m.accessOrd {
        \t\t\tm.mode = jMapModeNullOnly
        \t\t\tm.ext = v
        \t\t\tm.Live = 1
        \t\t\tm.firstIdx = 0
        \t\t\treturn nil, true
        \t\t}
        \t\tif m.Identity || m.accessOrd || k == jTombstoneStr {
        \t\t\tg := m.upgradeToGeneral()
        \t\t\tif k == NullString {
        \t\t\t\treturn m.putGenNull(g, v)
        \t\t\t}
        \t\t\treturn m.putGenStr(g, k, boxedKey, v)
        \t\t}
        \t\tm.strEntries = make([]jStrMapEntry, 1, m.strInitCap())
        \t\tm.strEntries[0] = jStrMapEntry{key: k, val: v}
        \t\tm.mode = jMapModeStr
        \t\tm.Live = 1
        \t\treturn nil, true
        \tcase jMapModeNullOnly:
        \t\tif k == NullString {
        \t\t\told := m.ext
        \t\t\tm.ext = v
        \t\t\treturn old, false
        \t\t}
        \t\tif m.Identity || m.accessOrd || k == jTombstoneStr {
        \t\t\tg := m.upgradeToGeneral()
        \t\t\treturn m.putGenStr(g, k, boxedKey, v)
        \t\t}
        \t\tnullVal := m.ext
        \t\tm.strEntries = make([]jStrMapEntry, 2, max(m.initCap(), 2))
        \t\tm.strEntries[0] = jStrMapEntry{key: NullString, val: nullVal}
        \t\tm.strEntries[1] = jStrMapEntry{key: k, val: v}
        \t\tm.ext = nil
        \t\tm.mode = jMapModeStr
        \t\tm.Live = 2
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\treturn nil, true
        \tcase jMapModeStr:
        \t\tif k == jTombstoneStr {
        \t\t\tg := m.upgradeToGeneral()
        \t\t\treturn m.putGenStr(g, k, boxedKey, v)
        \t\t}
        \t\tif m.ext != nil {
        \t\t\tidxMap := m.ext.(map[string]int32)
        \t\t\tif idx, ok := idxMap[k]; ok {
        \t\t\t\told := m.strEntries[idx].val
        \t\t\t\tm.strEntries[idx].val = v
        \t\t\t\treturn old, false
        \t\t\t}
        \t\t\tif int(m.holes)*2 >= len(m.strEntries) && len(m.strEntries) >= jStrSmallLimit {
        \t\t\t\tm.compactStr()
        \t\t\t}
        \t\t\tidxMap[k] = int32(len(m.strEntries))
        \t\t\tm.strEntries = append(m.strEntries, jStrMapEntry{key: k, val: v})
        \t\t\tm.Live++
        \t\t\treturn nil, true
        \t\t}
        \t\tfor i := range m.strEntries {
        \t\t\tif m.strEntries[i].key == k {
        \t\t\t\told := m.strEntries[i].val
        \t\t\t\tm.strEntries[i].val = v
        \t\t\t\treturn old, false
        \t\t\t}
        \t\t}
        \t\tif m.holes > 0 && int(m.holes)*2 >= len(m.strEntries) {
        \t\t\tm.compactStr()
        \t\t}
        \t\tif len(m.strEntries) >= jStrSmallLimit {
        \t\t\tidxMap := make(map[string]int32, max(len(m.strEntries)+1, min(int(m.capHint), 256)))
        \t\t\tfor i := range m.strEntries {
        \t\t\t\tif ek := m.strEntries[i].key; ek != jTombstoneStr {
        \t\t\t\t\tidxMap[ek] = int32(i)
        \t\t\t\t}
        \t\t\t}
        \t\t\tidxMap[k] = int32(len(m.strEntries))
        \t\t\tm.ext = idxMap
        \t\t}
        \t\tm.strEntries = append(m.strEntries, jStrMapEntry{key: k, val: v})
        \t\tm.Live++
        \t\treturn nil, true
        \tcase jMapModeSortedNat, jMapModeSortedCmp:
        \t\tif k == NullString {
        \t\t\tif m.mode == jMapModeSortedNat {
        \t\t\t\tThrowNPE()
        \t\t\t}
        \t\t\tst := m.ensureSortedState()
        \t\t\tnr, old, ins := m.avlInsertAny(st, st.root, nil, v)
        \t\t\tst.root = nr
        \t\t\treturn old, ins
        \t\t}
        \t\tst := m.ensureSortedState()
        \t\tnr, old, ins := m.avlInsertStr(st, st.root, k, v)
        \t\tst.root = nr
        \t\treturn old, ins
        \tcase jMapModeNode1:
        \t\tif k == NullString && !m.Identity && !m.accessOrd {
        \t\t\treturn m.putNodeInternal(0, v)
        \t\t}
        \t\tg := m.upgradeToGeneral()
        \t\tif k == NullString {
        \t\t\treturn m.putGenNull(g, v)
        \t\t}
        \t\treturn m.putGenStr(g, k, boxedKey, v)
        \tcase jMapModeNode:
        \t\tif k == NullString {
        \t\t\treturn m.putNodeInternal(0, v)
        \t\t}
        \t\tg := m.upgradeToGeneral()
        \t\treturn m.putGenStr(g, k, boxedKey, v)
        \tcase jMapModePtr:
        \t\tif k == NullString {
        \t\t\treturn m.putPtr(nil, v)
        \t\t}
        \t\tg := m.upgradeToGeneral()
        \t\treturn m.putGenStr(g, k, boxedKey, v)
        \tdefault:
        \t\tg := m.ext.(*jGenMapState)
        \t\tif k == NullString {
        \t\t\treturn m.putGenNull(g, v)
        \t\t}
        \t\treturn m.putGenStr(g, k, boxedKey, v)
        \t}
        }

        func (m *JMapCore) PutStr(k string, v any) any {
        \told, _ := m.putStrBoxed(k, nil, v)
        \treturn old
        }

        func (m *JMapCore) AddKeyStr(k string) bool {
        \tif _, ok := m.getStr(k, false); ok {
        \t\treturn false
        \t}
        \t_, ins := m.putStrBoxed(k, nil, true)
        \treturn ins
        }

        func (m *JMapCore) removeStrInternal(k string) (any, bool) {
        \tif m.immutable {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tif m.Live == 0 {
        \t\treturn nil, false
        \t}
        \tswitch m.mode {
        \tcase jMapModeStr:
        \t\tif k == jTombstoneStr {
        \t\t\treturn nil, false
        \t\t}
        \t\tif m.ext != nil {
        \t\t\tidx, ok := m.ext.(map[string]int32)[k]
        \t\t\tif !ok {
        \t\t\t\treturn nil, false
        \t\t\t}
        \t\t\treturn m.removeStrAt(int(idx)), true
        \t\t}
        \t\tfor i := range m.strEntries {
        \t\t\tif m.strEntries[i].key == k {
        \t\t\t\treturn m.removeStrAt(i), true
        \t\t\t}
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeNode:
        \t\tif k == NullString {
        \t\t\treturn m.removeNodeInternal(0)
        \t\t}
        \t\treturn nil, false
        \tcase jMapModePtr:
        \t\tif k == NullString {
        \t\t\treturn m.removePtr(nil)
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeNullOnly:
        \t\tif k == NullString {
        \t\t\told := m.ext
        \t\t\tm.ext = nil
        \t\t\tm.mode = jMapModeEmpty
        \t\t\tm.Live = 0
        \t\t\tm.holes = 0
        \t\t\tm.firstIdx = 0
        \t\t\treturn old, true
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeNode1:
        \t\treturn nil, false
        \tcase jMapModeSortedNat, jMapModeSortedCmp:
        \t\tst := m.ext.(*jSortedMapState)
        \t\tif k == NullString {
        \t\t\tif m.mode == jMapModeSortedNat {
        \t\t\t\treturn nil, false
        \t\t\t}
        \t\t\tnr, old, del := m.avlDeleteAny(st, st.root, nil)
        \t\t\tst.root = nr
        \t\t\treturn old, del
        \t\t}
        \t\tnr, old, del := m.avlDeleteStr(st, st.root, k)
        \t\tst.root = nr
        \t\treturn old, del
        \tcase jMapModeGeneral:
        \t\tg := m.ext.(*jGenMapState)
        \t\tif k == NullString {
        \t\t\tif g.nullPos == 0 {
        \t\t\t\treturn nil, false
        \t\t\t}
        \t\t\treturn m.removeGenPos(g.nullPos), true
        \t\t}
        \t\tif !m.Identity {
        \t\t\tif g.idx != nil && g.idx.strIdx != nil {
        \t\t\t\tpos := g.idx.strIdx[k]
        \t\t\t\tif pos == 0 {
        \t\t\t\t\treturn nil, false
        \t\t\t\t}
        \t\t\t\treturn m.removeGenPos(pos), true
        \t\t\t}
        \t\t\tfor pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
        \t\t\t\te := &g.entries[pos-1]
        \t\t\t\tif e.kind() == jKeyStr && e.Key.(string) == k {
        \t\t\t\t\treturn m.removeGenPos(pos), true
        \t\t\t\t}
        \t\t\t}
        \t\t\treturn nil, false
        \t\t}
        \t\th := m.hashGeneral(k)
        \t\tpos := m.findGenPos(g, k, h)
        \t\tif pos == 0 {
        \t\t\treturn nil, false
        \t\t}
        \t\treturn m.removeGenPos(pos), true
        \t}
        \treturn nil, false
        }

        func (m *JMapCore) RemoveStr(k string) any {
        \told, _ := m.removeStrInternal(k)
        \treturn old
        }

        func (m *JMapCore) RemoveKeyStr(k string) bool {
        \t_, del := m.removeStrInternal(k)
        \treturn del
        }

        func (m *JMapCore) getNode(u uint32, touch bool) (any, bool) {
        \tif m.Live == 0 {
        \t\treturn nil, false
        \t}
        \tswitch m.mode {
        \tcase jMapModeNode1:
        \t\tif u != 0 && uint32(m.firstIdx) == u {
        \t\t\treturn m.ext, true
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeNode:
        \t\tns := m.ext.(*jNodeMapState)
        \t\tif ns.nodeIdx != nil {
        \t\t\tif idx, ok := ns.nodeIdx[u]; ok {
        \t\t\t\treturn ns.entries[idx].val, true
        \t\t\t}
        \t\t\treturn nil, false
        \t\t}
        \t\tfor i := range ns.entries {
        \t\t\tif ns.entries[i].active && ns.entries[i].key == u {
        \t\t\t\treturn ns.entries[i].val, true
        \t\t\t}
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeStr:
        \t\tif u == 0 {
        \t\t\treturn m.getStr(NullString, false)
        \t\t}
        \t\treturn nil, false
        \tcase jMapModePtr:
        \t\tif u == 0 {
        \t\t\treturn m.getPtr(nil)
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeNullOnly:
        \t\tif u == 0 {
        \t\t\treturn m.ext, true
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeSortedNat, jMapModeSortedCmp:
        \t\tif u == 0 {
        \t\t\treturn m.getSortedNull()
        \t\t}
        \t\treturn m.getAny(boxNodeById(u), false)
        \tcase jMapModeGeneral:
        \t\tg := m.ext.(*jGenMapState)
        \t\tif u == 0 {
        \t\t\treturn m.getGenNull(g, touch)
        \t\t}
        \t\tif g.idx != nil && g.idx.u64Idx != nil {
        \t\t\tuk := (uint64(1) << 32) | uint64(u)
        \t\t\tif pos := g.idx.u64Idx[uk]; pos != 0 {
        \t\t\t\tif touch && m.accessOrd {
        \t\t\t\t\tg.moveToTail(pos)
        \t\t\t\t}
        \t\t\t\treturn g.entries[pos-1].Val, true
        \t\t\t}
        \t\t\treturn nil, false
        \t\t}
        \t\tfor pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
        \t\t\te := &g.entries[pos-1]
        \t\t\tif e.kind() == jKeyNode && e.u32OrHash == u {
        \t\t\t\tif touch && m.accessOrd {
        \t\t\t\t\tg.moveToTail(pos)
        \t\t\t\t}
        \t\t\t\treturn e.Val, true
        \t\t\t}
        \t\t}
        \t\treturn nil, false
        \t}
        \treturn nil, false
        }

        func (m *JMapCore) ContainsNode(u uint32) bool {
        \t_, ok := m.getNode(u, false)
        \treturn ok
        }

        func (m *JMapCore) GetNode(u uint32) any {
        \tv, _ := m.getNode(u, true)
        \treturn v
        }

        func (m *JMapCore) GetOrDefaultNode(u uint32, def any) any {
        \tif v, ok := m.getNode(u, true); ok {
        \t\treturn v
        \t}
        \treturn ToAny(def)
        }

        func (m *JMapCore) putNodeInternal(u uint32, v any) (any, bool) {
        \tif m.immutable {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tv = ToAny(v)
        \tswitch m.mode {
        \tcase jMapModeEmpty:
        \t\tif m.accessOrd {
        \t\t\tg := m.upgradeToGeneral()
        \t\t\tif u == 0 {
        \t\t\t\treturn m.putGenNull(g, v)
        \t\t\t}
        \t\t\treturn m.putGenNode(g, u, v)
        \t\t}
        \t\tif u == 0 {
        \t\t\tm.mode = jMapModeNullOnly
        \t\t\tm.ext = v
        \t\t\tm.Live = 1
        \t\t\tm.firstIdx = 0
        \t\t\treturn nil, true
        \t\t}
        \t\tm.mode = jMapModeNode1
        \t\tm.firstIdx = int32(u)
        \t\tm.ext = v
        \t\tm.Live = 1
        \t\treturn nil, true
        \tcase jMapModeNode1:
        \t\tu1 := uint32(m.firstIdx)
        \t\tif u == u1 {
        \t\t\told := m.ext
        \t\t\tm.ext = v
        \t\t\treturn old, false
        \t\t}
        \t\tif u == 0 && m.Identity {
        \t\t\tg := m.upgradeToGeneral()
        \t\t\treturn m.putGenNull(g, v)
        \t\t}
        \t\tv1 := m.ext
        \t\tns := newNodeMapState2(m.capHint)
        \t\tns.entries[0] = jNodeMapEntry{key: u1, active: true, val: v1}
        \t\tns.entries[1] = jNodeMapEntry{key: u, active: true, val: v}
        \t\tm.ext = ns
        \t\tm.mode = jMapModeNode
        \t\tm.Live = 2
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\treturn nil, true
        \tcase jMapModeNullOnly:
        \t\tif u == 0 {
        \t\t\told := m.ext
        \t\t\tm.ext = v
        \t\t\treturn old, false
        \t\t}
        \t\tif m.Identity || m.accessOrd {
        \t\t\tg := m.upgradeToGeneral()
        \t\t\treturn m.putGenNode(g, u, v)
        \t\t}
        \t\tnullVal := m.ext
        \t\tns := newNodeMapState2(m.capHint)
        \t\tns.entries[0] = jNodeMapEntry{key: 0, active: true, val: nullVal}
        \t\tns.entries[1] = jNodeMapEntry{key: u, active: true, val: v}
        \t\tm.ext = ns
        \t\tm.mode = jMapModeNode
        \t\tm.Live = 2
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\treturn nil, true
        \tcase jMapModeNode:
        \t\tns := m.ext.(*jNodeMapState)
        \t\tif ns.nodeIdx != nil {
        \t\t\tif idx, ok := ns.nodeIdx[u]; ok {
        \t\t\t\told := ns.entries[idx].val
        \t\t\t\tns.entries[idx].val = v
        \t\t\t\treturn old, false
        \t\t\t}
        \t\t\tif int(m.holes)*2 >= len(ns.entries) && len(ns.entries) >= jNodeSmallLimit {
        \t\t\t\tm.compactNode(ns)
        \t\t\t}
        \t\t\tns.nodeIdx[u] = int32(len(ns.entries))
        \t\t\tns.entries = append(ns.entries, jNodeMapEntry{key: u, active: true, val: v})
        \t\t\tm.Live++
        \t\t\treturn nil, true
        \t\t}
        \t\tfor i := range ns.entries {
        \t\t\tif ns.entries[i].active && ns.entries[i].key == u {
        \t\t\t\told := ns.entries[i].val
        \t\t\t\tns.entries[i].val = v
        \t\t\t\treturn old, false
        \t\t\t}
        \t\t}
        \t\tif m.holes > 0 && int(m.holes)*2 >= len(ns.entries) {
        \t\t\tm.compactNode(ns)
        \t\t}
        \t\tif len(ns.entries) >= jNodeSmallLimit {
        \t\t\tns.nodeIdx = make(map[uint32]int32, max(len(ns.entries)+4, min(int(m.capHint), 256)))
        \t\t\tfor i := range ns.entries {
        \t\t\t\tif ns.entries[i].active {
        \t\t\t\t\tns.nodeIdx[ns.entries[i].key] = int32(i)
        \t\t\t\t}
        \t\t\t}
        \t\t\tns.nodeIdx[u] = int32(len(ns.entries))
        \t\t}
        \t\tns.entries = append(ns.entries, jNodeMapEntry{key: u, active: true, val: v})
        \t\tm.Live++
        \t\treturn nil, true
        \tcase jMapModeStr:
        \t\tif u == 0 {
        \t\t\treturn m.putStrBoxed(NullString, nil, v)
        \t\t}
        \t\tg := m.upgradeToGeneral()
        \t\treturn m.putGenNode(g, u, v)
        \tcase jMapModePtr:
        \t\tif u == 0 {
        \t\t\treturn m.putPtr(nil, v)
        \t\t}
        \t\tg := m.upgradeToGeneral()
        \t\treturn m.putGenNode(g, u, v)
        \tcase jMapModeSortedNat, jMapModeSortedCmp:
        \t\treturn m.putAnyInternal(boxNodeById(u), v)
        \tdefault:
        \t\tg := m.ext.(*jGenMapState)
        \t\tif u == 0 {
        \t\t\treturn m.putGenNull(g, v)
        \t\t}
        \t\treturn m.putGenNode(g, u, v)
        \t}
        }

        func (m *JMapCore) PutNode(u uint32, v any) any {
        \told, _ := m.putNodeInternal(u, v)
        \treturn old
        }

        func (m *JMapCore) AddKeyNode(u uint32) bool {
        \tif _, ok := m.getNode(u, false); ok {
        \t\treturn false
        \t}
        \t_, ins := m.putNodeInternal(u, true)
        \treturn ins
        }

        func (m *JMapCore) removeNodeInternal(u uint32) (any, bool) {
        \tif m.immutable {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tif m.Live == 0 {
        \t\treturn nil, false
        \t}
        \tswitch m.mode {
        \tcase jMapModeNode1:
        \t\tif u != 0 && uint32(m.firstIdx) == u {
        \t\t\told := m.ext
        \t\t\tm.ext = nil
        \t\t\tm.mode = jMapModeEmpty
        \t\t\tm.Live = 0
        \t\t\tm.holes = 0
        \t\t\tm.firstIdx = 0
        \t\t\treturn old, true
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeNode:
        \t\tns := m.ext.(*jNodeMapState)
        \t\tif ns.nodeIdx != nil {
        \t\t\tidx, ok := ns.nodeIdx[u]
        \t\t\tif !ok {
        \t\t\t\treturn nil, false
        \t\t\t}
        \t\t\treturn m.removeNodeAt(int(idx)), true
        \t\t}
        \t\tfor i := range ns.entries {
        \t\t\tif ns.entries[i].active && ns.entries[i].key == u {
        \t\t\t\treturn m.removeNodeAt(i), true
        \t\t\t}
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeStr:
        \t\tif u == 0 {
        \t\t\treturn m.removeStrInternal(NullString)
        \t\t}
        \t\treturn nil, false
        \tcase jMapModePtr:
        \t\tif u == 0 {
        \t\t\treturn m.removePtr(nil)
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeNullOnly:
        \t\tif u == 0 {
        \t\t\told := m.ext
        \t\t\tm.ext = nil
        \t\t\tm.mode = jMapModeEmpty
        \t\t\tm.Live = 0
        \t\t\tm.holes = 0
        \t\t\tm.firstIdx = 0
        \t\t\treturn old, true
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeSortedNat, jMapModeSortedCmp:
        \t\treturn m.removeAnyInternal(boxNodeById(u))
        \tcase jMapModeGeneral:
        \t\tg := m.ext.(*jGenMapState)
        \t\tif u == 0 {
        \t\t\tif g.nullPos == 0 {
        \t\t\t\treturn nil, false
        \t\t\t}
        \t\t\treturn m.removeGenPos(g.nullPos), true
        \t\t}
        \t\tif g.idx != nil && g.idx.u64Idx != nil {
        \t\t\tuk := (uint64(1) << 32) | uint64(u)
        \t\t\tpos := g.idx.u64Idx[uk]
        \t\t\tif pos == 0 {
        \t\t\t\treturn nil, false
        \t\t\t}
        \t\t\treturn m.removeGenPos(pos), true
        \t\t}
        \t\tfor pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
        \t\t\te := &g.entries[pos-1]
        \t\t\tif e.kind() == jKeyNode && e.u32OrHash == u {
        \t\t\t\treturn m.removeGenPos(pos), true
        \t\t\t}
        \t\t}
        \t\treturn nil, false
        \t}
        \treturn nil, false
        }

        func (m *JMapCore) RemoveNode(u uint32) any {
        \told, _ := m.removeNodeInternal(u)
        \treturn old
        }

        func (m *JMapCore) RemoveKeyNode(u uint32) bool {
        \t_, del := m.removeNodeInternal(u)
        \treturn del
        }
        """;
  }

  private String buildCoreRuntimeGoPart2MapOps() {
    return """

        func (m *JMapCore) getPtr(k any) (any, bool) {
        \tps := m.ext.(*jPtrMapState)
        \tif ps.ptrIdx != nil {
        \t\tif idx, ok := ps.ptrIdx[k]; ok {
        \t\t\treturn ps.entries[idx].val, true
        \t\t}
        \t\treturn nil, false
        \t}
        \tfor i := range ps.entries {
        \t\tif ps.entries[i].key == k {
        \t\t\treturn ps.entries[i].val, true
        \t\t}
        \t}
        \treturn nil, false
        }

        func (m *JMapCore) putPtr(k, v any) (any, bool) {
        \tv = ToAny(v)
        \tif m.mode == jMapModeEmpty {
        \t\tif m.accessOrd {
        \t\t\tg := m.upgradeToGeneral()
        \t\t\tif k == nil {
        \t\t\t\treturn m.putGenNull(g, v)
        \t\t\t}
        \t\t\treturn m.putGenPtr(g, k, v)
        \t\t}
        \t\tif k == nil {
        \t\t\tm.mode = jMapModeNullOnly
        \t\t\tm.ext = v
        \t\t\tm.Live = 1
        \t\t\tm.firstIdx = 0
        \t\t\treturn nil, true
        \t\t}
        \t\tvar ps *jPtrMapState
        \t\tif m.capHint > 4 {
        \t\t\tps = &jPtrMapState{entries: make([]jPtrMapEntry, 1, min(int(m.capHint), 256))}
        \t\t} else {
        \t\t\ts4 := &jPtrMapState4{}
        \t\t\tps = &s4.state
        \t\t\tps.entries = s4.buf[:1]
        \t\t}
        \t\tps.entries[0] = jPtrMapEntry{key: k, val: v}
        \t\tm.ext = ps
        \t\tm.mode = jMapModePtr
        \t\tm.Live = 1
        \t\treturn nil, true
        \t}
        \tif m.mode == jMapModeNullOnly {
        \t\tif k == nil {
        \t\t\told := m.ext
        \t\t\tm.ext = v
        \t\t\treturn old, false
        \t\t}
        \t\tif m.accessOrd {
        \t\t\tg := m.upgradeToGeneral()
        \t\t\treturn m.putGenPtr(g, k, v)
        \t\t}
        \t\tnullVal := m.ext
        \t\tvar ps *jPtrMapState
        \t\tif m.capHint > 4 {
        \t\t\tps = &jPtrMapState{entries: make([]jPtrMapEntry, 2, min(int(m.capHint), 256))}
        \t\t} else {
        \t\t\ts4 := &jPtrMapState4{}
        \t\t\tps = &s4.state
        \t\t\tps.entries = s4.buf[:2]
        \t\t}
        \t\tps.entries[0] = jPtrMapEntry{key: nil, val: nullVal}
        \t\tps.entries[1] = jPtrMapEntry{key: k, val: v}
        \t\tm.ext = ps
        \t\tm.mode = jMapModePtr
        \t\tm.Live = 2
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\treturn nil, true
        \t}
        \tif m.mode != jMapModePtr {
        \t\tg := m.upgradeToGeneral()
        \t\tif k == nil {
        \t\t\treturn m.putGenNull(g, v)
        \t\t}
        \t\treturn m.putGenPtr(g, k, v)
        \t}
        \tps := m.ext.(*jPtrMapState)
        \tif ps.ptrIdx != nil {
        \t\tif idx, ok := ps.ptrIdx[k]; ok {
        \t\t\told := ps.entries[idx].val
        \t\t\tps.entries[idx].val = v
        \t\t\treturn old, false
        \t\t}
        \t\tif int(m.holes)*2 >= len(ps.entries) && len(ps.entries) >= jSmallLimit {
        \t\t\tm.compactPtr(ps)
        \t\t}
        \t\tps.ptrIdx[k] = int32(len(ps.entries))
        \t\tps.entries = append(ps.entries, jPtrMapEntry{key: k, val: v})
        \t\tm.Live++
        \t\treturn nil, true
        \t}
        \tfor i := range ps.entries {
        \t\tif ps.entries[i].key == k {
        \t\t\told := ps.entries[i].val
        \t\t\tps.entries[i].val = v
        \t\t\treturn old, false
        \t\t}
        \t}
        \tif m.holes > 0 && int(m.holes)*2 >= len(ps.entries) {
        \t\tm.compactPtr(ps)
        \t}
        \tif len(ps.entries) >= jSmallLimit {
        \t\tps.ptrIdx = make(map[any]int32, max(len(ps.entries)+4, min(int(m.capHint), 256)))
        \t\tfor i := range ps.entries {
        \t\t\tif ek := ps.entries[i].key; ek != jTombstoneObj {
        \t\t\t\tps.ptrIdx[ek] = int32(i)
        \t\t\t}
        \t\t}
        \t\tps.ptrIdx[k] = int32(len(ps.entries))
        \t}
        \tps.entries = append(ps.entries, jPtrMapEntry{key: k, val: v})
        \tm.Live++
        \treturn nil, true
        }

        func (m *JMapCore) removePtr(k any) (any, bool) {
        \tps := m.ext.(*jPtrMapState)
        \tif ps.ptrIdx != nil {
        \t\tidx, ok := ps.ptrIdx[k]
        \t\tif !ok {
        \t\t\treturn nil, false
        \t\t}
        \t\treturn m.removePtrAt(int(idx)), true
        \t}
        \tfor i := range ps.entries {
        \t\tif ps.entries[i].key == k {
        \t\t\treturn m.removePtrAt(i), true
        \t\t}
        \t}
        \treturn nil, false
        }

        func isIdentityOrEnumKey(k any) bool {
        \treturn isFastPtrMapKey(k, false)
        }

        func isFastPtrMapKey(k any, identity bool) bool {
        \tif _, ok := k.(IdentityComparable); ok {
        \t\treturn true
        \t}
        \tif _, ok := k.(EnumAccessor); ok {
        \t\treturn true
        \t}
        \tif _, ok := k.(jSourceObjPtr); ok {
        \t\tif identity {
        \t\t\treturn true
        \t\t}
        \t\tif _, hasEq := k.(javalang.Equaler); hasEq {
        \t\t\treturn false
        \t\t}
        \t\t_, hasHc := k.(jHashCoder)
        \t\treturn !hasHc
        \t}
        \treturn false
        }

        func (m *JMapCore) getAny(k any, touch bool) (any, bool) {
        \tif m.Live == 0 {
        \t\treturn nil, false
        \t}
        \tswitch m.mode {
        \tcase jMapModeEmpty:
        \t\treturn nil, false
        \tcase jMapModeStr:
        \t\tif s, ok := k.(string); ok {
        \t\t\treturn m.getStr(s, false)
        \t\t}
        \t\tif k == nil || isJavaNilAny(k) {
        \t\t\treturn m.getStr(NullString, false)
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeNode1:
        \t\tif rn, ok := k.(jRawNodeProvider); ok {
        \t\t\tu := rn.RawNodeId_()
        \t\t\tif u != 0 && uint32(m.firstIdx) == u {
        \t\t\t\treturn m.ext, true
        \t\t\t}
        \t\t\treturn nil, false
        \t\t}
        \t\tif nd, ok := k.(nodestore.Node); ok {
        \t\t\tif nd != 0 && uint32(m.firstIdx) == uint32(nd) {
        \t\t\t\treturn m.ext, true
        \t\t\t}
        \t\t\treturn nil, false
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeNode:
        \t\tif rn, ok := k.(jRawNodeProvider); ok {
        \t\t\treturn m.getNode(rn.RawNodeId_(), false)
        \t\t}
        \t\tif nd, ok := k.(nodestore.Node); ok {
        \t\t\treturn m.getNode(uint32(nd), false)
        \t\t}
        \t\tif k == nil || isJavaNilAny(k) {
        \t\t\treturn m.getNode(0, false)
        \t\t}
        \t\treturn nil, false
        \tcase jMapModePtr:
        \t\tif k == nil {
        \t\t\treturn m.getPtr(nil)
        \t\t}
        \t\tif ik, ok := k.(jIdentityKey); ok {
        \t\t\tif ik.IsJavaNil_() {
        \t\t\t\treturn m.getPtr(nil)
        \t\t\t}
        \t\t\treturn m.getPtr(k)
        \t\t}
        \t\tk = ToAny(k)
        \t\tif k == nil {
        \t\t\treturn m.getPtr(nil)
        \t\t}
        \t\tif isFastPtrMapKey(k, m.Identity) {
        \t\t\treturn m.getPtr(k)
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeGeneral:
        \t\tg := m.ext.(*jGenMapState)
        \t\tif k == nil {
        \t\t\treturn m.getGenNull(g, touch)
        \t\t}
        \t\tif ik, ok := k.(jIdentityKey); ok {
        \t\t\tif ik.IsJavaNil_() {
        \t\t\t\treturn m.getGenNull(g, touch)
        \t\t\t}
        \t\t\treturn m.getGenPtrKey(g, k, touch)
        \t\t}
        \t\tif rn, ok := k.(jRawNodeProvider); ok {
        \t\t\treturn m.getNode(rn.RawNodeId_(), touch)
        \t\t}
        \t\tif !m.Identity {
        \t\t\tswitch x := k.(type) {
        \t\t\tcase string:
        \t\t\t\treturn m.getStr(x, touch)
        \t\t\tcase nodestore.Node:
        \t\t\t\treturn m.getNode(uint32(x), touch)
        \t\t\tcase *BoxInt:
        \t\t\t\tif x == nil {
        \t\t\t\t\treturn m.getGenNull(g, touch)
        \t\t\t\t}
        \t\t\t\treturn m.getGenInt32(g, x.Val, touch)
        \t\t\tcase int32:
        \t\t\t\treturn m.getGenInt32(g, x, touch)
        \t\t\tcase FastCustomKey:
        \t\t\t\tif x.IsJavaNil_() {
        \t\t\t\t\treturn m.getGenNull(g, touch)
        \t\t\t\t}
        \t\t\t\tif pos := findGenFastCustom(g, k, x, x.HashCode()); pos != 0 {
        \t\t\t\t\tif touch && m.accessOrd {
        \t\t\t\t\t\tg.moveToTail(pos)
        \t\t\t\t\t}
        \t\t\t\t\treturn g.entries[pos-1].Val, true
        \t\t\t\t}
        \t\t\t\treturn nil, false
        \t\t\t}
        \t\t}
        \t\tk = ToAny(k)
        \t\tif k == nil {
        \t\t\treturn m.getGenNull(g, touch)
        \t\t}
        \t\tif isFastPtrMapKey(k, m.Identity) {
        \t\t\treturn m.getGenPtrKey(g, k, touch)
        \t\t}
        \t\th := m.hashGeneral(k)
        \t\tif pos := m.findGenPos(g, k, h); pos != 0 {
        \t\t\tif touch && m.accessOrd {
        \t\t\t\tg.moveToTail(pos)
        \t\t\t}
        \t\t\treturn g.entries[pos-1].Val, true
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeNullOnly:
        \t\tif k == nil || isJavaNilAny(k) {
        \t\t\treturn m.ext, true
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeSortedNat, jMapModeSortedCmp:
        \t\tif s, ok := k.(string); ok {
        \t\t\tif s == NullString {
        \t\t\t\treturn m.getSortedNull()
        \t\t\t}
        \t\t\tst := m.ext.(*jSortedMapState)
        \t\t\tid := st.avlFindStr(s)
        \t\t\tif id == 0 {
        \t\t\t\treturn nil, false
        \t\t\t}
        \t\t\treturn st.tree[id-1].val, true
        \t\t}
        \t\tk = ToAny(k)
        \t\tif k == nil {
        \t\t\treturn m.getSortedNull()
        \t\t}
        \t\tst := m.ext.(*jSortedMapState)
        \t\tid := st.avlFindAny(k)
        \t\tif id == 0 {
        \t\t\treturn nil, false
        \t\t}
        \t\treturn st.tree[id-1].val, true
        \t}
        \treturn nil, false
        }

        func (m *JMapCore) getNullKey(touch bool) (any, bool) {
        \tswitch m.mode {
        \tcase jMapModeNullOnly:
        \t\treturn m.ext, true
        \tcase jMapModeStr:
        \t\treturn m.getStr(NullString, false)
        \tcase jMapModeNode:
        \t\treturn m.getNode(0, false)
        \tcase jMapModePtr:
        \t\treturn m.getPtr(nil)
        \tcase jMapModeGeneral:
        \t\treturn m.getGenNull(m.ext.(*jGenMapState), touch)
        \tcase jMapModeSortedNat, jMapModeSortedCmp:
        \t\treturn m.getSortedNull()
        \t}
        \treturn nil, false
        }

        func (m *JMapCore) getGenPtrKey(g *jGenMapState, k any, touch bool) (any, bool) {
        \tif g.idx != nil && g.idx.ptrIdx != nil {
        \t\tif pos := g.idx.ptrIdx[k]; pos != 0 {
        \t\t\tif touch && m.accessOrd {
        \t\t\t\tg.moveToTail(pos)
        \t\t\t}
        \t\t\treturn g.entries[pos-1].Val, true
        \t\t}
        \t\treturn nil, false
        \t}
        \tfor pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
        \t\te := &g.entries[pos-1]
        \t\tif e.kind() == jKeyPtr && e.Key == k {
        \t\t\tif touch && m.accessOrd {
        \t\t\t\tg.moveToTail(pos)
        \t\t\t}
        \t\t\treturn e.Val, true
        \t\t}
        \t}
        \treturn nil, false
        }

        func (m *JMapCore) getGenInt32(g *jGenMapState, iv int32, touch bool) (any, bool) {
        \tu := uint32(iv)
        \tif g.idx != nil && g.idx.u64Idx != nil {
        \t\tuk := (uint64(2) << 32) | uint64(u)
        \t\tif pos := g.idx.u64Idx[uk]; pos != 0 {
        \t\t\tif touch && m.accessOrd {
        \t\t\t\tg.moveToTail(pos)
        \t\t\t}
        \t\t\treturn g.entries[pos-1].Val, true
        \t\t}
        \t\treturn nil, false
        \t}
        \tfor pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
        \t\te := &g.entries[pos-1]
        \t\tif e.kind() == jKeyInt && e.u32OrHash == u {
        \t\t\tif touch && m.accessOrd {
        \t\t\t\tg.moveToTail(pos)
        \t\t\t}
        \t\t\treturn e.Val, true
        \t\t}
        \t}
        \treturn nil, false
        }

        func (m *JMapCore) putNullKey(v any) (any, bool) {
        \tv = ToAny(v)
        \tswitch m.mode {
        \tcase jMapModeEmpty:
        \t\tif m.accessOrd {
        \t\t\tg := m.upgradeToGeneral()
        \t\t\treturn m.putGenNull(g, v)
        \t\t}
        \t\tm.mode = jMapModeNullOnly
        \t\tm.ext = v
        \t\tm.Live = 1
        \t\tm.firstIdx = 0
        \t\treturn nil, true
        \tcase jMapModeNullOnly:
        \t\told := m.ext
        \t\tm.ext = v
        \t\treturn old, false
        \tcase jMapModeStr:
        \t\treturn m.putStrBoxed(NullString, nil, v)
        \tcase jMapModeNode1, jMapModeNode:
        \t\treturn m.putNodeInternal(0, v)
        \tcase jMapModePtr:
        \t\treturn m.putPtr(nil, v)
        \tcase jMapModeGeneral:
        \t\treturn m.putGenNull(m.ext.(*jGenMapState), v)
        \tcase jMapModeSortedNat:
        \t\tThrowNPE()
        \t\treturn nil, false
        \tdefault:
        \t\tst := m.ensureSortedState()
        \t\tnr, old, ins := m.avlInsertAny(st, st.root, nil, v)
        \t\tst.root = nr
        \t\treturn old, ins
        \t}
        }

        func (m *JMapCore) putAnyInternal(k, v any) (any, bool) {
        \tif m.immutable {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tif k == nil {
        \t\treturn m.putNullKey(v)
        \t}
        \tif !m.IsSorted() {
        \t\tif ik, ok := k.(jIdentityKey); ok {
        \t\t\tif ik.IsJavaNil_() {
        \t\t\t\treturn m.putNullKey(v)
        \t\t\t}
        \t\t\treturn m.putPtr(k, v)
        \t\t}
        \t\tif rn, ok := k.(jRawNodeProvider); ok {
        \t\t\treturn m.putNodeInternal(rn.RawNodeId_(), v)
        \t\t}
        \t\tif !m.Identity {
        \t\t\tswitch x := k.(type) {
        \t\t\tcase string:
        \t\t\t\treturn m.putStrBoxed(x, k, v)
        \t\t\tcase nodestore.Node:
        \t\t\t\treturn m.putNodeInternal(uint32(x), v)
        \t\t\tcase *BoxInt:
        \t\t\t\tif x == nil {
        \t\t\t\t\treturn m.putNullKey(v)
        \t\t\t\t}
        \t\t\t\tg := m.upgradeToGeneral()
        \t\t\t\treturn m.putGenInt32(g, x.Val, x, ToAny(v))
        \t\t\tcase int32:
        \t\t\t\tg := m.upgradeToGeneral()
        \t\t\t\treturn m.putGenInt32(g, x, nil, ToAny(v))
        \t\t\tcase FastCustomKey:
        \t\t\t\tif x.IsJavaNil_() {
        \t\t\t\t\treturn m.putNullKey(v)
        \t\t\t\t}
        \t\t\t\tvar g *jGenMapState
        \t\t\t\tif m.mode == jMapModeGeneral {
        \t\t\t\t\tg = m.ext.(*jGenMapState)
        \t\t\t\t} else {
        \t\t\t\t\tg = m.upgradeToGeneral()
        \t\t\t\t}
        \t\t\t\treturn m.putGenFastCustom(g, k, x, ToAny(v))
        \t\t\t}
        \t\t}
        \t} else if iv, ok := k.(int32); ok {
        \t\tk = ValueOfInt(iv)
        \t}
        \tk = ToAny(k)
        \tif k == nil {
        \t\treturn m.putNullKey(v)
        \t}
        \tif m.IsSorted() {
        \t\tv = ToAny(v)
        \t\tst := m.ensureSortedState()
        \t\tnr, old, ins := m.avlInsertAny(st, st.root, k, v)
        \t\tst.root = nr
        \t\treturn old, ins
        \t}
        \tif isFastPtrMapKey(k, m.Identity) {
        \t\treturn m.putPtr(k, v)
        \t}
        \tg := m.upgradeToGeneral()
        \treturn m.putGenOther(g, k, ToAny(v))
        }

        func (m *JMapCore) removeNullKey() (any, bool) {
        \tswitch m.mode {
        \tcase jMapModeNullOnly:
        \t\told := m.ext
        \t\tm.ext = nil
        \t\tm.mode = jMapModeEmpty
        \t\tm.Live = 0
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\treturn old, true
        \tcase jMapModeStr:
        \t\treturn m.removeStrInternal(NullString)
        \tcase jMapModeNode:
        \t\treturn m.removeNodeInternal(0)
        \tcase jMapModePtr:
        \t\treturn m.removePtr(nil)
        \tcase jMapModeGeneral:
        \t\tg := m.ext.(*jGenMapState)
        \t\tif g.nullPos == 0 {
        \t\t\treturn nil, false
        \t\t}
        \t\treturn m.removeGenPos(g.nullPos), true
        \tcase jMapModeSortedCmp:
        \t\tif m.Live == 0 {
        \t\t\treturn nil, false
        \t\t}
        \t\tst := m.ext.(*jSortedMapState)
        \t\tnr, old, del := m.avlDeleteAny(st, st.root, nil)
        \t\tst.root = nr
        \t\treturn old, del
        \t}
        \treturn nil, false
        }

        func (m *JMapCore) removeAnyInternal(k any) (any, bool) {
        \tif m.immutable {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tif m.Live == 0 {
        \t\treturn nil, false
        \t}
        \tswitch m.mode {
        \tcase jMapModeEmpty:
        \t\treturn nil, false
        \tcase jMapModeStr:
        \t\tif s, ok := k.(string); ok {
        \t\t\treturn m.removeStrInternal(s)
        \t\t}
        \t\tif k == nil || isJavaNilAny(k) {
        \t\t\treturn m.removeStrInternal(NullString)
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeNode1:
        \t\tif rn, ok := k.(jRawNodeProvider); ok {
        \t\t\treturn m.removeNodeInternal(rn.RawNodeId_())
        \t\t}
        \t\tif nd, ok := k.(nodestore.Node); ok {
        \t\t\treturn m.removeNodeInternal(uint32(nd))
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeNode:
        \t\tif rn, ok := k.(jRawNodeProvider); ok {
        \t\t\treturn m.removeNodeInternal(rn.RawNodeId_())
        \t\t}
        \t\tif nd, ok := k.(nodestore.Node); ok {
        \t\t\treturn m.removeNodeInternal(uint32(nd))
        \t\t}
        \t\tif k == nil || isJavaNilAny(k) {
        \t\t\treturn m.removeNodeInternal(0)
        \t\t}
        \t\treturn nil, false
        \tcase jMapModePtr:
        \t\tif k == nil {
        \t\t\treturn m.removePtr(nil)
        \t\t}
        \t\tif ik, ok := k.(jIdentityKey); ok {
        \t\t\tif ik.IsJavaNil_() {
        \t\t\t\treturn m.removePtr(nil)
        \t\t\t}
        \t\t\treturn m.removePtr(k)
        \t\t}
        \t\tk = ToAny(k)
        \t\tif k == nil {
        \t\t\treturn m.removePtr(nil)
        \t\t}
        \t\tif isFastPtrMapKey(k, m.Identity) {
        \t\t\treturn m.removePtr(k)
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeGeneral:
        \t\tg := m.ext.(*jGenMapState)
        \t\tif k == nil {
        \t\t\tif g.nullPos == 0 {
        \t\t\t\treturn nil, false
        \t\t\t}
        \t\t\treturn m.removeGenPos(g.nullPos), true
        \t\t}
        \t\tif ik, ok := k.(jIdentityKey); ok {
        \t\t\tif ik.IsJavaNil_() {
        \t\t\t\tif g.nullPos == 0 {
        \t\t\t\t\treturn nil, false
        \t\t\t\t}
        \t\t\t\treturn m.removeGenPos(g.nullPos), true
        \t\t\t}
        \t\t\treturn m.removeGenPtrKey(g, k)
        \t\t}
        \t\tif rn, ok := k.(jRawNodeProvider); ok {
        \t\t\treturn m.removeNodeInternal(rn.RawNodeId_())
        \t\t}
        \t\tif !m.Identity {
        \t\t\tswitch x := k.(type) {
        \t\t\tcase string:
        \t\t\t\treturn m.removeStrInternal(x)
        \t\t\tcase nodestore.Node:
        \t\t\t\treturn m.removeNodeInternal(uint32(x))
        \t\t\tcase *BoxInt:
        \t\t\t\tif x == nil {
        \t\t\t\t\tif g.nullPos == 0 {
        \t\t\t\t\t\treturn nil, false
        \t\t\t\t\t}
        \t\t\t\t\treturn m.removeGenPos(g.nullPos), true
        \t\t\t\t}
        \t\t\t\treturn m.removeGenInt32(g, x.Val)
        \t\t\tcase int32:
        \t\t\t\treturn m.removeGenInt32(g, x)
        \t\t\tcase FastCustomKey:
        \t\t\t\tif x.IsJavaNil_() {
        \t\t\t\t\tif g.nullPos == 0 {
        \t\t\t\t\t\treturn nil, false
        \t\t\t\t\t}
        \t\t\t\t\treturn m.removeGenPos(g.nullPos), true
        \t\t\t\t}
        \t\t\t\tpos := findGenFastCustom(g, k, x, x.HashCode())
        \t\t\t\tif pos == 0 {
        \t\t\t\t\treturn nil, false
        \t\t\t\t}
        \t\t\t\treturn m.removeGenPos(pos), true
        \t\t\t}
        \t\t}
        \t\tk = ToAny(k)
        \t\tif k == nil {
        \t\t\tif g.nullPos == 0 {
        \t\t\t\treturn nil, false
        \t\t\t}
        \t\t\treturn m.removeGenPos(g.nullPos), true
        \t\t}
        \t\tif isFastPtrMapKey(k, m.Identity) {
        \t\t\treturn m.removeGenPtrKey(g, k)
        \t\t}
        \t\th := m.hashGeneral(k)
        \t\tpos := m.findGenPos(g, k, h)
        \t\tif pos == 0 {
        \t\t\treturn nil, false
        \t\t}
        \t\treturn m.removeGenPos(pos), true
        \tcase jMapModeNullOnly:
        \t\tif k == nil || isJavaNilAny(k) {
        \t\t\treturn m.removeNullKey()
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeSortedNat, jMapModeSortedCmp:
        \t\tif s, ok := k.(string); ok {
        \t\t\treturn m.removeStrInternal(s)
        \t\t}
        \t\tif iv, ok := k.(int32); ok {
        \t\t\tk = ValueOfInt(iv)
        \t\t} else {
        \t\t\tk = ToAny(k)
        \t\t}
        \t\tif k == nil {
        \t\t\treturn m.removeNullKey()
        \t\t}
        \t\tst := m.ext.(*jSortedMapState)
        \t\tnr, old, del := m.avlDeleteAny(st, st.root, k)
        \t\tst.root = nr
        \t\treturn old, del
        \t}
        \treturn nil, false
        }

        func (m *JMapCore) removeGenPtrKey(g *jGenMapState, k any) (any, bool) {
        \tif g.idx != nil && g.idx.ptrIdx != nil {
        \t\tpos := g.idx.ptrIdx[k]
        \t\tif pos == 0 {
        \t\t\treturn nil, false
        \t\t}
        \t\treturn m.removeGenPos(pos), true
        \t}
        \tfor pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
        \t\te := &g.entries[pos-1]
        \t\tif e.kind() == jKeyPtr && e.Key == k {
        \t\t\treturn m.removeGenPos(pos), true
        \t\t}
        \t}
        \treturn nil, false
        }

        func (m *JMapCore) removeGenInt32(g *jGenMapState, iv int32) (any, bool) {
        \tu := uint32(iv)
        \tif g.idx != nil && g.idx.u64Idx != nil {
        \t\tuk := (uint64(2) << 32) | uint64(u)
        \t\tpos := g.idx.u64Idx[uk]
        \t\tif pos == 0 {
        \t\t\treturn nil, false
        \t\t}
        \t\treturn m.removeGenPos(pos), true
        \t}
        \tfor pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
        \t\te := &g.entries[pos-1]
        \t\tif e.kind() == jKeyInt && e.u32OrHash == u {
        \t\t\treturn m.removeGenPos(pos), true
        \t\t}
        \t}
        \treturn nil, false
        }

        func (m *JMapCore) GetEntry(k any) (any, bool) {
        \treturn m.getAny(k, false)
        }

        func (m *JMapCore) CopyFromSameMode(src *JMapCore) bool {
        \tif m == nil || src == nil || m == src || m.immutable || m.mode != jMapModeEmpty || m.Live != 0 {
        \t\treturn false
        \t}
        \tif m.IsSorted() || src.IsSorted() || m.Identity != src.Identity || (src.accessOrd && !m.accessOrd) {
        \t\treturn false
        \t}
        \tif src.Live == 0 {
        \t\treturn true
        \t}
        \tswitch src.mode {
        \tcase jMapModeNullOnly:
        \t\tif m.accessOrd {
        \t\t\treturn false
        \t\t}
        \t\tm.mode = jMapModeNullOnly
        \t\tm.ext = src.ext
        \t\tm.Live = 1
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\treturn true
        \tcase jMapModeNode1:
        \t\tif m.accessOrd {
        \t\t\treturn false
        \t\t}
        \t\tm.mode = jMapModeNode1
        \t\tm.firstIdx = src.firstIdx
        \t\tm.ext = src.ext
        \t\tm.Live = 1
        \t\tm.holes = 0
        \t\treturn true
        \tcase jMapModeStr:
        \t\tif m.accessOrd {
        \t\t\treturn false
        \t\t}
        \t\tif src.holes == 0 && src.firstIdx == 0 {
        \t\t\tn := len(src.strEntries)
        \t\t\tm.strEntries = make([]jStrMapEntry, n)
        \t\t\tcopy(m.strEntries, src.strEntries)
        \t\t\tif n > jStrSmallLimit {
        \t\t\t\tdstIdx := make(map[string]int32, n)
        \t\t\t\tfor i := range m.strEntries {
        \t\t\t\t\tdstIdx[m.strEntries[i].key] = int32(i)
        \t\t\t\t}
        \t\t\t\tm.ext = dstIdx
        \t\t\t}
        \t\t} else {
        \t\t\tlive := int(src.Live)
        \t\t\tentries := make([]jStrMapEntry, live)
        \t\t\tvar dstIdx map[string]int32
        \t\t\tif live > jStrSmallLimit {
        \t\t\t\tdstIdx = make(map[string]int32, live)
        \t\t\t}
        \t\t\tw := 0
        \t\t\tfor i := int(src.firstIdx); i < len(src.strEntries); i++ {
        \t\t\t\te := src.strEntries[i]
        \t\t\t\tif e.key != jTombstoneStr {
        \t\t\t\t\tif dstIdx != nil {
        \t\t\t\t\t\tdstIdx[e.key] = int32(w)
        \t\t\t\t\t}
        \t\t\t\t\tentries[w] = e
        \t\t\t\t\tw++
        \t\t\t\t}
        \t\t\t}
        \t\t\tm.strEntries = entries
        \t\t\tif dstIdx != nil {
        \t\t\t\tm.ext = dstIdx
        \t\t\t}
        \t\t}
        \t\tm.mode = jMapModeStr
        \t\tm.Live = src.Live
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\treturn true
        \tcase jMapModeNode:
        \t\tif m.accessOrd {
        \t\t\treturn false
        \t\t}
        \t\tsns := src.ext.(*jNodeMapState)
        \t\tn := int(src.Live)
        \t\tif n == 1 {
        \t\t\tfor i := int(src.firstIdx); i < len(sns.entries); i++ {
        \t\t\t\tif e := sns.entries[i]; e.active {
        \t\t\t\t\tif e.key == 0 {
        \t\t\t\t\t\tm.mode = jMapModeNullOnly
        \t\t\t\t\t\tm.firstIdx = 0
        \t\t\t\t\t} else {
        \t\t\t\t\t\tm.mode = jMapModeNode1
        \t\t\t\t\t\tm.firstIdx = int32(e.key)
        \t\t\t\t\t}
        \t\t\t\t\tm.ext = e.val
        \t\t\t\t\tm.Live = 1
        \t\t\t\t\tm.holes = 0
        \t\t\t\t\treturn true
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t\tvar dns *jNodeMapState
        \t\tif n <= 4 {
        \t\t\ts4 := &jNodeMapState4{}
        \t\t\tdns = &s4.state
        \t\t\tdns.entries = s4.buf[:n]
        \t\t} else if n <= 8 {
        \t\t\ts8 := &jNodeMapState8{}
        \t\t\tdns = &s8.state
        \t\t\tdns.entries = s8.buf[:n]
        \t\t} else {
        \t\t\tdns = &jNodeMapState{entries: make([]jNodeMapEntry, n)}
        \t\t}
        \t\tif src.holes == 0 && src.firstIdx == 0 && len(sns.entries) == n {
        \t\t\tcopy(dns.entries, sns.entries)
        \t\t} else {
        \t\t\tw := 0
        \t\t\tfor i := int(src.firstIdx); i < len(sns.entries); i++ {
        \t\t\t\tif e := sns.entries[i]; e.active {
        \t\t\t\t\tdns.entries[w] = e
        \t\t\t\t\tw++
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t\tif n > jNodeSmallLimit {
        \t\t\tdstIdx := make(map[uint32]int32, n+4)
        \t\t\tfor i := range dns.entries {
        \t\t\t\tdstIdx[dns.entries[i].key] = int32(i)
        \t\t\t}
        \t\t\tdns.nodeIdx = dstIdx
        \t\t}
        \t\tm.ext = dns
        \t\tm.mode = jMapModeNode
        \t\tm.Live = src.Live
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\treturn true
        \tcase jMapModePtr:
        \t\tif m.accessOrd {
        \t\t\treturn false
        \t\t}
        \t\tsps := src.ext.(*jPtrMapState)
        \t\tn := int(src.Live)
        \t\tvar dps *jPtrMapState
        \t\tif n <= 4 {
        \t\t\ts4 := &jPtrMapState4{}
        \t\t\tdps = &s4.state
        \t\t\tdps.entries = s4.buf[:n]
        \t\t} else {
        \t\t\tdps = &jPtrMapState{entries: make([]jPtrMapEntry, n)}
        \t\t}
        \t\tif src.holes == 0 && src.firstIdx == 0 && len(sps.entries) == n {
        \t\t\tcopy(dps.entries, sps.entries)
        \t\t} else {
        \t\t\tw := 0
        \t\t\tfor i := int(src.firstIdx); i < len(sps.entries); i++ {
        \t\t\t\tif e := sps.entries[i]; e.key != jTombstoneObj {
        \t\t\t\t\tdps.entries[w] = e
        \t\t\t\t\tw++
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t\tif n > jSmallLimit {
        \t\t\tdstIdx := make(map[any]int32, n+4)
        \t\t\tfor i := range dps.entries {
        \t\t\t\tdstIdx[dps.entries[i].key] = int32(i)
        \t\t\t}
        \t\t\tdps.ptrIdx = dstIdx
        \t\t}
        \t\tm.ext = dps
        \t\tm.mode = jMapModePtr
        \t\tm.Live = src.Live
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\treturn true
        \tcase jMapModeGeneral:
        \t\tsg := src.ext.(*jGenMapState)
        \t\tif sg.freeHead == 0 {
        \t\t\tn := len(sg.entries)
        \t\t\tdg := newGenMapStateSized(n)
        \t\t\tentries := dg.entries
        \t\t\tcopy(entries, sg.entries)
        \t\t\tdg.head = sg.head
        \t\t\tdg.tail = sg.tail
        \t\t\tdg.nullPos = sg.nullPos
        \t\t\tif sg.idx != nil && n > jSmallLimit {
        \t\t\t\tsidx := sg.idx
        \t\t\t\tdidx := &jGenMapIdx{otherLive: sidx.otherLive}
        \t\t\t\tif len(sidx.buckets) > 0 {
        \t\t\t\t\tdidx.buckets = make([]int32, len(sidx.buckets))
        \t\t\t\t\tcopy(didx.buckets, sidx.buckets)
        \t\t\t\t}
        \t\t\t\thasMapIdx := false
        \t\t\t\tif sidx.strIdx != nil {
        \t\t\t\t\tdidx.strIdx = make(map[string]int32, len(sidx.strIdx))
        \t\t\t\t\thasMapIdx = true
        \t\t\t\t}
        \t\t\t\tif sidx.u64Idx != nil {
        \t\t\t\t\tdidx.u64Idx = make(map[uint64]int32, len(sidx.u64Idx))
        \t\t\t\t\thasMapIdx = true
        \t\t\t\t}
        \t\t\t\tif sidx.ptrIdx != nil {
        \t\t\t\t\tdidx.ptrIdx = make(map[any]int32, len(sidx.ptrIdx))
        \t\t\t\t\thasMapIdx = true
        \t\t\t\t}
        \t\t\t\tif hasMapIdx {
        \t\t\t\t\tfor i := range entries {
        \t\t\t\t\t\te := &entries[i]
        \t\t\t\t\t\tpos := int32(i + 1)
        \t\t\t\t\t\tswitch e.kind() {
        \t\t\t\t\t\tcase jKeyStr:
        \t\t\t\t\t\t\tif didx.strIdx != nil {
        \t\t\t\t\t\t\t\tdidx.strIdx[e.Key.(string)] = pos
        \t\t\t\t\t\t\t}
        \t\t\t\t\t\tcase jKeyNode:
        \t\t\t\t\t\t\tif didx.u64Idx != nil {
        \t\t\t\t\t\t\t\tdidx.u64Idx[(uint64(1)<<32)|uint64(e.u32OrHash)] = pos
        \t\t\t\t\t\t\t}
        \t\t\t\t\t\tcase jKeyInt:
        \t\t\t\t\t\t\tif didx.u64Idx != nil {
        \t\t\t\t\t\t\t\tdidx.u64Idx[(uint64(2)<<32)|uint64(e.u32OrHash)] = pos
        \t\t\t\t\t\t\t}
        \t\t\t\t\t\tcase jKeyPtr:
        \t\t\t\t\t\t\tif didx.ptrIdx != nil {
        \t\t\t\t\t\t\t\tdidx.ptrIdx[e.Key] = pos
        \t\t\t\t\t\t\t}
        \t\t\t\t\t\t}
        \t\t\t\t\t}
        \t\t\t\t}
        \t\t\t\tdg.idx = didx
        \t\t\t} else if sg.idx != nil {
        \t\t\t\tfor i := range entries {
        \t\t\t\t\tentries[i].hashNext = 0
        \t\t\t\t}
        \t\t\t}
        \t\t\tm.ext = dg
        \t\t\tm.mode = jMapModeGeneral
        \t\t\tm.Live = src.Live
        \t\t\tm.holes = 0
        \t\t\tm.firstIdx = 0
        \t\t\tm.genNonCustom = src.genNonCustom
        \t\t\treturn true
        \t\t}
        \t\tlive := int(src.Live)
        \t\tdg := newGenMapStateSized(live)
        \t\tentries := dg.entries
        \t\tdg.head = 1
        \t\tdg.tail = src.Live
        \t\tvar didx *jGenMapIdx
        \t\tvar bucketMask uint32
        \t\tif sg.idx != nil && live > jSmallLimit {
        \t\t\tsidx := sg.idx
        \t\t\tdidx = &jGenMapIdx{}
        \t\t\tif len(sidx.buckets) > 0 && sidx.otherLive > 0 {
        \t\t\t\tcapB := 16
        \t\t\t\tfor capB*3 <= int(sidx.otherLive)*4 && capB < 2048 {
        \t\t\t\t\tcapB <<= 1
        \t\t\t\t}
        \t\t\t\tdidx.buckets = make([]int32, capB)
        \t\t\t\tbucketMask = uint32(capB - 1)
        \t\t\t}
        \t\t\tif sidx.strIdx != nil {
        \t\t\t\tdidx.strIdx = make(map[string]int32, len(sidx.strIdx))
        \t\t\t}
        \t\t\tif sidx.u64Idx != nil {
        \t\t\t\tdidx.u64Idx = make(map[uint64]int32, len(sidx.u64Idx))
        \t\t\t}
        \t\t\tif sidx.ptrIdx != nil {
        \t\t\t\tdidx.ptrIdx = make(map[any]int32, len(sidx.ptrIdx))
        \t\t\t}
        \t\t\tdg.idx = didx
        \t\t}
        \t\ti := 0
        \t\tfor pos := sg.head; pos != 0; pos = sg.entries[pos-1].next() {
        \t\t\tse := sg.entries[pos-1]
        \t\t\tkind := se.kind()
        \t\t\tnewPos := int32(i + 1)
        \t\t\tvar nextPos int32
        \t\t\tif i+1 < live {
        \t\t\t\tnextPos = int32(i + 2)
        \t\t\t}
        \t\t\tde := JMapEntry{
        \t\t\t\tprev:        int32(i),
        \t\t\t\tnextAndKind: makeNextKind(nextPos, kind),
        \t\t\t\tu32OrHash:   se.u32OrHash,
        \t\t\t\tKey:         se.Key,
        \t\t\t\tVal:         se.Val,
        \t\t\t}
        \t\t\tswitch kind {
        \t\t\tcase jKeyNull:
        \t\t\t\tdg.nullPos = newPos
        \t\t\tcase jKeyStr:
        \t\t\t\tif didx != nil && didx.strIdx != nil {
        \t\t\t\t\tdidx.strIdx[se.Key.(string)] = newPos
        \t\t\t\t}
        \t\t\tcase jKeyNode:
        \t\t\t\tif didx != nil && didx.u64Idx != nil {
        \t\t\t\t\tdidx.u64Idx[(uint64(1)<<32)|uint64(se.u32OrHash)] = newPos
        \t\t\t\t}
        \t\t\tcase jKeyInt:
        \t\t\t\tif didx != nil && didx.u64Idx != nil {
        \t\t\t\t\tdidx.u64Idx[(uint64(2)<<32)|uint64(se.u32OrHash)] = newPos
        \t\t\t\t}
        \t\t\tcase jKeyPtr:
        \t\t\t\tif didx != nil && didx.ptrIdx != nil {
        \t\t\t\t\tdidx.ptrIdx[se.Key] = newPos
        \t\t\t\t}
        \t\t\tcase jKeyGeneral:
        \t\t\t\tif didx != nil && len(didx.buckets) > 0 {
        \t\t\t\t\tslot := smearHash(int32(se.u32OrHash)) & bucketMask
        \t\t\t\t\tde.hashNext = didx.buckets[slot]
        \t\t\t\t\tdidx.buckets[slot] = newPos
        \t\t\t\t\tdidx.otherLive++
        \t\t\t\t}
        \t\t\t}
        \t\t\tentries[i] = de
        \t\t\ti++
        \t\t}
        \t\tm.ext = dg
        \t\tm.mode = jMapModeGeneral
        \t\tm.Live = src.Live
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\tm.genNonCustom = src.genNonCustom
        \t\treturn true
        \t}
        \treturn false
        }

        func (m *JMapCore) ContainsKey(k any) bool {
        \t_, ok := m.getAny(k, false)
        \treturn ok
        }

        func (m *JMapCore) ContainsValue(v any) bool {
        \tv = ToAny(v)
        \tfound := false
        \tm.ForEachEntry(func(_, ev any) bool {
        \t\tif javalang.ObjectsEquals(ev, v) {
        \t\t\tfound = true
        \t\t\treturn false
        \t\t}
        \t\treturn true
        \t})
        \treturn found
        }

        func (m *JMapCore) Get(k any) any {
        \tv, _ := m.getAny(k, true)
        \treturn v
        }

        func findMapCoreEqVal[K FastEqValKey](m *JMapCore, k K) (any, bool) {
        \tif m.Live == 0 {
        \t\treturn nil, false
        \t}
        \th, nonNull, isIdent := k.FastValHash_()
        \tif !nonNull {
        \t\treturn m.getNullKey(false)
        \t}
        \tswitch m.mode {
        \tcase jMapModeGeneral:
        \t\tg := m.ext.(*jGenMapState)
        \t\tif idx := g.idx; idx != nil && len(idx.buckets) > 0 {
        \t\t\tentries := g.entries
        \t\t\tslot := smearHash(h) & uint32(len(idx.buckets)-1)
        \t\t\tfor pos := idx.buckets[slot]; pos != 0; {
        \t\t\t\te := &entries[pos-1]
        \t\t\t\tif int32(e.u32OrHash) == h {
        \t\t\t\t\tif ek, ok := e.Key.(K); ok && ek == k {
        \t\t\t\t\t\tif m.accessOrd {
        \t\t\t\t\t\t\tg.moveToTail(pos)
        \t\t\t\t\t\t}
        \t\t\t\t\t\treturn e.Val, true
        \t\t\t\t\t}
        \t\t\t\t}
        \t\t\t\tpos = e.hashNext
        \t\t\t}
        \t\t\treturn nil, false
        \t\t}
        \t\tentries := g.entries
        \t\tuh := uint32(h)
        \t\tfor i := range entries {
        \t\t\te := &entries[i]
        \t\t\tif e.u32OrHash == uh && e.kind() == jKeyGeneral {
        \t\t\t\tif ek, ok := e.Key.(K); ok && ek == k {
        \t\t\t\t\tif m.accessOrd {
        \t\t\t\t\t\tg.moveToTail(int32(i + 1))
        \t\t\t\t\t}
        \t\t\t\t\treturn e.Val, true
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t\treturn nil, false
        \tcase jMapModePtr:
        \t\tif !isIdent {
        \t\t\treturn nil, false
        \t\t}
        \t\tps := m.ext.(*jPtrMapState)
        \t\tfor i := range ps.entries {
        \t\t\tif ek, ok := ps.entries[i].key.(K); ok && ek == k {
        \t\t\t\treturn ps.entries[i].val, true
        \t\t\t}
        \t\t}
        \t\treturn nil, false
        \tcase jMapModeEmpty, jMapModeStr, jMapModeNode, jMapModeNode1, jMapModeNullOnly:
        \t\treturn nil, false
        \tdefault:
        \t\treturn m.getAny(k, false)
        \t}
        }

        func (m *JMapCore) GetOrDefault(k, def any) any {
        \tif v, ok := m.getAny(k, true); ok {
        \t\treturn v
        \t}
        \treturn ToAny(def)
        }

        func (m *JMapCore) Put(k, v any) any {
        \told, _ := m.putAnyInternal(k, v)
        \treturn old
        }

        func (m *JMapCore) AddKey(k any) bool {
        \tif _, ok := m.getAny(k, false); ok {
        \t\treturn false
        \t}
        \t_, ins := m.putAnyInternal(k, true)
        \treturn ins
        }

        func (m *JMapCore) PutIfAbsent(k, v any) any {
        \tif cur, ok := m.getAny(k, true); ok && cur != nil {
        \t\treturn cur
        \t}
        \treturn m.Put(k, v)
        }

        func (m *JMapCore) Remove(k any) any {
        \told, _ := m.removeAnyInternal(k)
        \treturn old
        }

        func (m *JMapCore) RemoveKey(k any) bool {
        \t_, del := m.removeAnyInternal(k)
        \treturn del
        }

        func (m *JMapCore) FirstEntry() (any, any) {
        \tif m == nil || m.Live == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty set/map"))
        \t}
        \tstart := int(m.firstIdx)
        \tswitch m.mode {
        \tcase jMapModeSortedNat, jMapModeSortedCmp:
        \t\tst := m.ext.(*jSortedMapState)
        \t\tcur := st.root
        \t\tfor st.tree[cur-1].left != 0 {
        \t\t\tcur = st.tree[cur-1].left
        \t\t}
        \t\treturn st.treeNodeKey(cur), st.tree[cur-1].val
        \tcase jMapModeStr:
        \t\tfor i := start; i < len(m.strEntries); i++ {
        \t\t\tif ek := m.strEntries[i].key; ek != jTombstoneStr {
        \t\t\t\tif i != start && !m.immutable {
        \t\t\t\t\tm.firstIdx = int32(i)
        \t\t\t\t}
        \t\t\t\tif ek == NullString {
        \t\t\t\t\treturn nil, m.strEntries[i].val
        \t\t\t\t}
        \t\t\t\treturn ek, m.strEntries[i].val
        \t\t\t}
        \t\t}
        \tcase jMapModeNode:
        \t\tns := m.ext.(*jNodeMapState)
        \t\tfor i := start; i < len(ns.entries); i++ {
        \t\t\tif ns.entries[i].active {
        \t\t\t\tif i != start && !m.immutable {
        \t\t\t\t\tm.firstIdx = int32(i)
        \t\t\t\t}
        \t\t\t\treturn boxNodeById(ns.entries[i].key), ns.entries[i].val
        \t\t\t}
        \t\t}
        \tcase jMapModePtr:
        \t\tps := m.ext.(*jPtrMapState)
        \t\tfor i := start; i < len(ps.entries); i++ {
        \t\t\tif ek := ps.entries[i].key; ek != jTombstoneObj {
        \t\t\t\tif i != start && !m.immutable {
        \t\t\t\t\tm.firstIdx = int32(i)
        \t\t\t\t}
        \t\t\t\treturn ek, ps.entries[i].val
        \t\t\t}
        \t\t}
        \tcase jMapModeGeneral:
        \t\tg := m.ext.(*jGenMapState)
        \t\treturn g.entryKey(g.head), g.entries[g.head-1].Val
        \tcase jMapModeNullOnly:
        \t\treturn nil, m.ext
        \tcase jMapModeNode1:
        \t\treturn boxNodeById(uint32(m.firstIdx)), m.ext
        \t}
        \tpanic(javalang.NewNoSuchElementException("empty set/map"))
        }

        func (m *JMapCore) FirstKey() any {
        \tk, _ := m.FirstEntry()
        \treturn k
        }

        func (m *JMapCore) LastKey() any {
        \tif m == nil || m.Live == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty set/map"))
        \t}
        \tswitch m.mode {
        \tcase jMapModeSortedNat, jMapModeSortedCmp:
        \t\tst := m.ext.(*jSortedMapState)
        \t\tcur := st.root
        \t\tfor st.tree[cur-1].right != 0 {
        \t\t\tcur = st.tree[cur-1].right
        \t\t}
        \t\treturn st.treeNodeKey(cur)
        \tcase jMapModeStr:
        \t\tfor i := len(m.strEntries) - 1; i >= 0; i-- {
        \t\t\tif ek := m.strEntries[i].key; ek != jTombstoneStr {
        \t\t\t\tif ek == NullString {
        \t\t\t\t\treturn nil
        \t\t\t\t}
        \t\t\t\treturn ek
        \t\t\t}
        \t\t}
        \tcase jMapModeNode:
        \t\tns := m.ext.(*jNodeMapState)
        \t\tfor i := len(ns.entries) - 1; i >= 0; i-- {
        \t\t\tif ns.entries[i].active {
        \t\t\t\treturn boxNodeById(ns.entries[i].key)
        \t\t\t}
        \t\t}
        \tcase jMapModePtr:
        \t\tps := m.ext.(*jPtrMapState)
        \t\tfor i := len(ps.entries) - 1; i >= 0; i-- {
        \t\t\tif ek := ps.entries[i].key; ek != jTombstoneObj {
        \t\t\t\treturn ek
        \t\t\t}
        \t\t}
        \tcase jMapModeGeneral:
        \t\tg := m.ext.(*jGenMapState)
        \t\treturn g.entryKey(g.tail)
        \tcase jMapModeNullOnly:
        \t\treturn nil
        \tcase jMapModeNode1:
        \t\treturn boxNodeById(uint32(m.firstIdx))
        \t}
        \tpanic(javalang.NewNoSuchElementException("empty set/map"))
        }

        func (m *JMapCore) RemoveFirstKey() any {
        \tif m == nil || m.Live == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty set/map"))
        \t}
        \tif m.immutable {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tif m.IsSorted() {
        \t\tst := m.ext.(*jSortedMapState)
        \t\tnr, minId := st.avlExtractMin(st.root)
        \t\tst.root = nr
        \t\tk := st.treeNodeKey(minId)
        \t\tst.avlFreeNode(minId)
        \t\tm.Live--
        \t\treturn k
        \t}
        \tstart := int(m.firstIdx)
        \tswitch m.mode {
        \tcase jMapModeStr:
        \t\tfor i := start; i < len(m.strEntries); i++ {
        \t\t\tif ek := m.strEntries[i].key; ek != jTombstoneStr {
        \t\t\t\tm.firstIdx = int32(i)
        \t\t\t\tm.removeStrAt(i)
        \t\t\t\tif ek == NullString {
        \t\t\t\t\treturn nil
        \t\t\t\t}
        \t\t\t\treturn ek
        \t\t\t}
        \t\t}
        \tcase jMapModeNode:
        \t\tns := m.ext.(*jNodeMapState)
        \t\tfor i := start; i < len(ns.entries); i++ {
        \t\t\tif ns.entries[i].active {
        \t\t\t\tk := boxNodeById(ns.entries[i].key)
        \t\t\t\tm.firstIdx = int32(i)
        \t\t\t\tm.removeNodeAt(i)
        \t\t\t\treturn k
        \t\t\t}
        \t\t}
        \tcase jMapModePtr:
        \t\tps := m.ext.(*jPtrMapState)
        \t\tfor i := start; i < len(ps.entries); i++ {
        \t\t\tif ek := ps.entries[i].key; ek != jTombstoneObj {
        \t\t\t\tm.firstIdx = int32(i)
        \t\t\t\tm.removePtrAt(i)
        \t\t\t\treturn ek
        \t\t\t}
        \t\t}
        \tcase jMapModeGeneral:
        \t\tg := m.ext.(*jGenMapState)
        \t\tk := g.entryKey(g.head)
        \t\tm.removeGenPos(g.head)
        \t\treturn k
        \tcase jMapModeNullOnly:
        \t\tm.ext = nil
        \t\tm.mode = jMapModeEmpty
        \t\tm.Live = 0
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\treturn nil
        \tcase jMapModeNode1:
        \t\tk := boxNodeById(uint32(m.firstIdx))
        \t\tm.ext = nil
        \t\tm.mode = jMapModeEmpty
        \t\tm.Live = 0
        \t\tm.holes = 0
        \t\tm.firstIdx = 0
        \t\treturn k
        \t}
        \tpanic(javalang.NewNoSuchElementException("empty set/map"))
        }

        func (m *JMapCore) RemoveLastKey() any {
        \tif m == nil || m.Live == 0 {
        \t\tpanic(javalang.NewNoSuchElementException("empty set/map"))
        \t}
        \tif m.IsSorted() {
        \t\tst := m.ext.(*jSortedMapState)
        \t\tnr, maxId := st.avlExtractMax(st.root)
        \t\tst.root = nr
        \t\tk := st.treeNodeKey(maxId)
        \t\tst.avlFreeNode(maxId)
        \t\tm.Live--
        \t\treturn k
        \t}
        \tk := m.LastKey()
        \tm.Remove(k)
        \treturn k
        }

        func (m *JMapCore) AddFirstKey(k, v any) {
        \tif m.IsSorted() {
        \t\t_, _ = m.removeAnyInternal(k)
        \t\t_, _ = m.putAnyInternal(k, v)
        \t\treturn
        \t}
        \tg := m.upgradeToGeneral()
        \t_, _ = m.removeAnyInternal(k)
        \t_, _ = m.putAnyInternal(k, v)
        \tg.moveTailToFront()
        }

        func (m *JMapCore) AddLastKey(k, v any) {
        \t_, _ = m.removeAnyInternal(k)
        \t_, _ = m.putAnyInternal(k, v)
        }

        func (m *JMapCore) Clear() {
        \tif m.immutable {
        \t\tpanic(javalang.NewUnsupportedOperationException("immutable"))
        \t}
        \tif m.Live == 0 && m.holes == 0 && m.firstIdx == 0 {
        \t\treturn
        \t}
        \tswitch m.mode {
        \tcase jMapModeNullOnly, jMapModeNode1:
        \t\tm.ext = nil
        \t\tm.mode = jMapModeEmpty
        \tcase jMapModeStr:
        \t\tclear(m.strEntries)
        \t\tm.strEntries = m.strEntries[:0]
        \t\tif m.ext != nil {
        \t\t\tclear(m.ext.(map[string]int32))
        \t\t}
        \tcase jMapModeNode:
        \t\tns := m.ext.(*jNodeMapState)
        \t\tclear(ns.entries)
        \t\tns.entries = ns.entries[:0]
        \t\tif ns.nodeIdx != nil {
        \t\t\tclear(ns.nodeIdx)
        \t\t}
        \tcase jMapModePtr:
        \t\tps := m.ext.(*jPtrMapState)
        \t\tclear(ps.entries)
        \t\tps.entries = ps.entries[:0]
        \t\tif ps.ptrIdx != nil {
        \t\t\tclear(ps.ptrIdx)
        \t\t}
        \tcase jMapModeGeneral:
        \t\tg := m.ext.(*jGenMapState)
        \t\tg.head = 0
        \t\tg.tail = 0
        \t\tg.freeHead = 0
        \t\tg.nullPos = 0
        \t\tif g.entries != nil {
        \t\t\tclear(g.entries)
        \t\t\tg.entries = g.entries[:0]
        \t\t}
        \t\tif g.idx != nil {
        \t\t\tif g.idx.strIdx != nil {
        \t\t\t\tclear(g.idx.strIdx)
        \t\t\t}
        \t\t\tif g.idx.u64Idx != nil {
        \t\t\t\tclear(g.idx.u64Idx)
        \t\t\t}
        \t\t\tif g.idx.ptrIdx != nil {
        \t\t\t\tclear(g.idx.ptrIdx)
        \t\t\t}
        \t\t\tif len(g.idx.buckets) > 0 {
        \t\t\t\tclear(g.idx.buckets)
        \t\t\t}
        \t\t\tg.idx.otherLive = 0
        \t\t}
        \t\tm.genNonCustom = m.Identity
        \tcase jMapModeSortedNat, jMapModeSortedCmp:
        \t\tif m.ext != nil {
        \t\t\tst := m.ext.(*jSortedMapState)
        \t\t\tif st.tree != nil {
        \t\t\t\tclear(st.tree)
        \t\t\t\tst.tree = st.tree[:0]
        \t\t\t}
        \t\t\tst.root = 0
        \t\t\tst.freeTree = 0
        \t\t}
        \t}
        \tm.Live = 0
        \tm.holes = 0
        \tm.firstIdx = 0
        }

        func (m *JMapCore) sealImmutable() {
        \tif m == nil || m.immutable {
        \t\treturn
        \t}
        \tswitch m.mode {
        \tcase jMapModeStr:
        \t\tm.compactStr()
        \t\tm.firstIdx = 0
        \tcase jMapModeNode:
        \t\tm.compactNode(m.ext.(*jNodeMapState))
        \t\tm.firstIdx = 0
        \tcase jMapModePtr:
        \t\tm.compactPtr(m.ext.(*jPtrMapState))
        \t\tm.firstIdx = 0
        \tcase jMapModeNode1:
        \tdefault:
        \t\tm.firstIdx = 0
        \t}
        \tm.immutable = true
        }

        func (m *JMapCore) SetHash() int32 {
        \tif m == nil || m.Live == 0 {
        \t\treturn 0
        \t}
        \tif m.IsSorted() {
        \t\treturn m.ext.(*jSortedMapState).avlSetHash(m.ext.(*jSortedMapState).root)
        \t}
        \tvar h int32
        \tm.ForEachEntry(func(k, _ any) bool {
        \t\th += elemJavaHashCode(k, m.Identity)
        \t\treturn true
        \t})
        \treturn h
        }

        func (m *JMapCore) ImmutableSetHash() int32 {
        \treturn m.SetHash()
        }

        func (m *JMapCore) ForEachEntry(fn func(k, v any) bool) {
        \tif m == nil || m.Live == 0 {
        \t\treturn
        \t}
        \tstart := int(m.firstIdx)
        \tswitch m.mode {
        \tcase jMapModeNullOnly:
        \t\tfn(nil, m.ext)
        \tcase jMapModeNode1:
        \t\tfn(boxNodeById(uint32(m.firstIdx)), m.ext)
        \tcase jMapModeStr:
        \t\tfor i := start; i < len(m.strEntries); i++ {
        \t\t\te := m.strEntries[i]
        \t\t\tif e.key == jTombstoneStr {
        \t\t\t\tcontinue
        \t\t\t}
        \t\t\tvar k any
        \t\t\tif e.key != NullString {
        \t\t\t\tk = e.key
        \t\t\t}
        \t\t\tif !fn(k, e.val) {
        \t\t\t\treturn
        \t\t\t}
        \t\t}
        \tcase jMapModeNode:
        \t\tns := m.ext.(*jNodeMapState)
        \t\tfor i := start; i < len(ns.entries); i++ {
        \t\t\te := ns.entries[i]
        \t\t\tif !e.active {
        \t\t\t\tcontinue
        \t\t\t}
        \t\t\tif !fn(boxNodeById(e.key), e.val) {
        \t\t\t\treturn
        \t\t\t}
        \t\t}
        \tcase jMapModePtr:
        \t\tps := m.ext.(*jPtrMapState)
        \t\tfor i := start; i < len(ps.entries); i++ {
        \t\t\te := ps.entries[i]
        \t\t\tif e.key == jTombstoneObj {
        \t\t\t\tcontinue
        \t\t\t}
        \t\t\tif !fn(e.key, e.val) {
        \t\t\t\treturn
        \t\t\t}
        \t\t}
        \tcase jMapModeGeneral:
        \t\tg := m.ext.(*jGenMapState)
        \t\tfor pos := g.head; pos != 0; {
        \t\t\tnextPos := g.entries[pos-1].next()
        \t\t\tk := g.entryKey(pos)
        \t\t\tv := g.entries[pos-1].Val
        \t\t\tif !fn(k, v) {
        \t\t\t\treturn
        \t\t\t}
        \t\t\tpos = nextPos
        \t\t}
        \tcase jMapModeSortedNat, jMapModeSortedCmp:
        \t\tst := m.ext.(*jSortedMapState)
        \t\tst.avlForEach(st.root, fn)
        \t}
        }

        func (m *JMapCore) TransformValues(fn func(k, v any) any) {
        \tif m == nil || m.Live == 0 {
        \t\treturn
        \t}
        \tstart := int(m.firstIdx)
        \tswitch m.mode {
        \tcase jMapModeNullOnly:
        \t\tm.ext = ToAny(fn(nil, m.ext))
        \tcase jMapModeNode1:
        \t\tm.ext = ToAny(fn(boxNodeById(uint32(m.firstIdx)), m.ext))
        \tcase jMapModeStr:
        \t\tfor i := start; i < len(m.strEntries); i++ {
        \t\t\te := &m.strEntries[i]
        \t\t\tif e.key == jTombstoneStr {
        \t\t\t\tcontinue
        \t\t\t}
        \t\t\tvar k any
        \t\t\tif e.key != NullString {
        \t\t\t\tk = e.key
        \t\t\t}
        \t\t\te.val = ToAny(fn(k, e.val))
        \t\t}
        \tcase jMapModeNode:
        \t\tns := m.ext.(*jNodeMapState)
        \t\tfor i := start; i < len(ns.entries); i++ {
        \t\t\te := &ns.entries[i]
        \t\t\tif !e.active {
        \t\t\t\tcontinue
        \t\t\t}
        \t\t\te.val = ToAny(fn(boxNodeById(e.key), e.val))
        \t\t}
        \tcase jMapModePtr:
        \t\tps := m.ext.(*jPtrMapState)
        \t\tfor i := start; i < len(ps.entries); i++ {
        \t\t\te := &ps.entries[i]
        \t\t\tif e.key == jTombstoneObj {
        \t\t\t\tcontinue
        \t\t\t}
        \t\t\te.val = ToAny(fn(e.key, e.val))
        \t\t}
        \tcase jMapModeGeneral:
        \t\tg := m.ext.(*jGenMapState)
        \t\tfor pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
        \t\t\tk := g.entryKey(pos)
        \t\t\tg.entries[pos-1].Val = ToAny(fn(k, g.entries[pos-1].Val))
        \t\t}
        \tcase jMapModeSortedNat, jMapModeSortedCmp:
        \t\tst := m.ext.(*jSortedMapState)
        \t\tst.avlTransformValues(st.root, fn)
        \t}
        }

        func (m *JMapCore) LiveKeys() []any {
        \tif m == nil || m.Live == 0 {
        \t\treturn []any{}
        \t}
        \tout := make([]any, 0, m.Live)
        \tm.ForEachEntry(func(k, _ any) bool {
        \t\tout = append(out, k)
        \t\treturn true
        \t})
        \treturn out
        }

        func (m *JMapCore) LiveValues() []any {
        \tif m == nil || m.Live == 0 {
        \t\treturn []any{}
        \t}
        \tout := make([]any, 0, m.Live)
        \tm.ForEachEntry(func(_, v any) bool {
        \t\tout = append(out, v)
        \t\treturn true
        \t})
        \treturn out
        }

        // --- Simple Slice Iterator for `for-each` Loops ---

        type SimpleIter interface {
        \tHasNext() bool
        \tNext() any
        }

        type sliceIter struct {
        \titems []any
        \tidx   int
        }

        func (s *sliceIter) IsJavaNil_() bool { return s == nil }
        func (s *sliceIter) HasNext() bool    { return s.idx < len(s.items) }
        func (s *sliceIter) Next() any {
        \tif s.idx >= len(s.items) {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tv := s.items[s.idx]
        \ts.idx++
        \treturn v
        }

        type strSliceIter struct {
        \tstrs []string
        \tidx  int
        }

        func (s *strSliceIter) HasNext() bool { return s.idx < len(s.strs) }
        func (s *strSliceIter) Next() any {
        \tif s.idx >= len(s.strs) {
        \t\tpanic(javalang.NewNoSuchElementException("iterator exhausted"))
        \t}
        \tv := s.strs[s.idx]
        \ts.idx++
        \tif v == NullString {
        \t\treturn nil
        \t}
        \treturn v
        }

        func Iterate(iterable any) SimpleIter {
        \tif iterable == nil {
        \t\tpanic(javalang.NewNullPointerException())
        \t}
        \tif ssi, ok := iterable.(StrSliceIterable); ok {
        \t\tif nc, ok2 := iterable.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
        \t\t\tpanic(javalang.NewNullPointerException())
        \t\t}
        \t\tif strs, ok2 := ssi.IterStrSlice_(); ok2 {
        \t\t\treturn &strSliceIter{strs: strs}
        \t\t}
        \t}
        \tif tv, ok := iterable.(jTypedViewCursorSource); ok {
        \t\tif nc, ok2 := iterable.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
        \t\t\tpanic(javalang.NewNullPointerException())
        \t\t}
        \t\tif _, ok2 := tv.typedViewLen_(); ok2 {
        \t\t\tif itProvider, ok3 := iterable.(interface{ Iterator() any }); ok3 {
        \t\t\t\tif si, ok4 := itProvider.Iterator().(SimpleIter); ok4 {
        \t\t\t\t\treturn si
        \t\t\t\t}
        \t\t\t}
        \t\t}
        \t}
        \tif si, ok := iterable.(SliceIterable); ok {
        \t\tif nc, ok2 := iterable.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
        \t\t\tpanic(javalang.NewNullPointerException())
        \t\t}
        \t\treturn &sliceIter{items: si.IterSlice_()}
        \t}
        \tif isJavaNilAny(iterable) {
        \t\tpanic(javalang.NewNullPointerException())
        \t}
        \tif itProvider, ok := iterable.(interface{ Iterator() Ext_java_util_Iterator }); ok {
        \t\treturn itProvider.Iterator()
        \t}
        \tif itProvider, ok := iterable.(interface{ Iterator_0() Ext_java_util_Iterator }); ok {
        \t\treturn itProvider.Iterator_0()
        \t}
        \tif itProvider, ok := iterable.(interface{ Iterator() any }); ok {
        \t\tif si, ok := itProvider.Iterator().(SimpleIter); ok {
        \t\t\treturn si
        \t\t}
        \t}
        \tif itProvider, ok := iterable.(interface{ Iterator_0() any }); ok {
        \t\tif si, ok := itProvider.Iterator_0().(SimpleIter); ok {
        \t\t\treturn si
        \t\t}
        \t}
        \trv := reflect.ValueOf(iterable)
        \tif rv.Kind() == reflect.Slice {
        \t\titems := make([]any, rv.Len())
        \t\tfor i := 0; i < rv.Len(); i++ {
        \t\t\titems[i] = rv.Index(i).Interface()
        \t\t}
        \t\treturn &sliceIter{items: items}
        \t}
        \tpanic(javalang.NewUnsupportedOperationException(fmt.Sprintf("Iterate unsupported on %T", iterable)))
        }
        """;
  }

  private void recordProtectedExtFields() {
    for (String qn :
        List.of(
            "java.util.zip.DeflaterOutputStream", "org.kohsuke.args4j.spi.OptionHandler")) {
      TypeElement te = resolver.elements.getTypeElement(qn);
      if (te != null) {
        for (Element enc : te.getEnclosedElements()) {
          if (enc instanceof VariableElement ve && !ve.getModifiers().contains(Modifier.STATIC)) {
            String fn = ve.getSimpleName().toString();
            if (fn.equals("def") || fn.equals("option") || fn.equals("setter")) {
              resolver.recordExternalField(ve);
            }
          }
        }
      }
    }
  }

  private String buildExtBindingsGo() {
    PackagePlanner.ImportTracker imports = new PackagePlanner.ImportTracker("j2gort");
    StringBuilder body = new StringBuilder(131072);

    // Ensure all core external types & protected fields referenced by j2gort helpers exist.
    for (String coreQn : CORE_RUNTIME_EXT_TYPES) {
      TypeElement coreTe = resolver.elements.getTypeElement(coreQn);
      if (coreTe != null) {
        resolver.recordExternalType(coreTe);
      }
    }
    recordProtectedExtFields();

    // Ensure every referenced external enum records ALL its enum constants so EnumSet.allOf
    // and CLI flag enum lookup work even for constants not referenced directly in Java source.
    for (TypeElement te : new ArrayList<>(resolver.getReferencedExternalTypes())) {
      if (te.getKind() == ElementKind.ENUM) {
        for (Element enc : te.getEnclosedElements()) {
          if (enc instanceof VariableElement ve && ve.getKind() == ElementKind.ENUM_CONSTANT) {
            resolver.recordExternalField(ve);
          }
        }
      }
    }

    List<TypeElement> extTypes = new ArrayList<>(resolver.getReferencedExternalTypes());
    extTypes.sort(Comparator.comparing(t -> t.getQualifiedName().toString()));

    // Emit external types
    for (TypeElement te : extTypes) {
      if (mapsToAny(te)
          || isShimException(te)
          || isBoxedPrimitive(te)
          || te.getQualifiedName().contentEquals("java.lang.String")
          || te.getQualifiedName().contentEquals("java.lang.Class")
          || te.getQualifiedName().contentEquals("java.lang.Number")) {
        continue;
      }
      emitExternalType(te, imports, body);
    }

    // Emit static methods & constructors referenced on ALL external types
    Set<ExecutableElement> emittedExecs = new HashSet<>();
    List<ExecutableElement> extExecs = new ArrayList<>(resolver.getReferencedExternalExecs());
    extExecs.sort(
        Comparator.comparing(
            (ExecutableElement e) ->
                ((TypeElement) e.getEnclosingElement()).getQualifiedName().toString()
                    + "#"
                    + e.getSimpleName()
                    + "#"
                    + e.getParameters().size()
                    + "#"
                    + e.getParameters()));

    for (ExecutableElement ee : extExecs) {
      TypeElement owner = (TypeElement) ee.getEnclosingElement();
      if (ee.getKind() == ElementKind.CONSTRUCTOR) {
        if (!isShimException(owner)
            && !isBoxedPrimitive(owner)
            && !mapsToAny(owner)
            && !owner.getQualifiedName().contentEquals("java.lang.String")
            && emittedExecs.add(ee)) {
          emitExternalConstructor(owner, ee, imports, body);
        }
      } else if (ee.getModifiers().contains(Modifier.STATIC)) {
        if (emittedExecs.add(ee)) {
          emitExternalStaticMethod(owner, ee, imports, body);
        }
      }
    }

    // Emit referenced external static fields / enum constants
    List<VariableElement> extFields = new ArrayList<>(resolver.getReferencedExternalFields());
    extFields.sort(
        Comparator.comparing(
            v ->
                ((TypeElement) v.getEnclosingElement()).getQualifiedName().toString()
                    + "#"
                    + v.getSimpleName()));
    Set<String> emittedVars = new HashSet<>();
    for (VariableElement ve : extFields) {
      TypeElement owner = (TypeElement) ve.getEnclosingElement();
      if (ve.getModifiers().contains(Modifier.STATIC) || ve.getKind() == ElementKind.ENUM_CONSTANT) {
        String varGoName = getExternalStaticFieldGoName(owner, ve);
        if (emittedVars.add(varGoName)) {
          emitExternalStaticField(owner, ve, varGoName, imports, body);
        }
      }
    }

    // Emit enum constant registrations for all referenced external enums
    StringBuilder enumRegInit = new StringBuilder();
    for (TypeElement te : extTypes) {
      if (te.getKind() == ElementKind.ENUM) {
        List<String> constVars = new ArrayList<>();
        for (Element enc : te.getEnclosedElements()) {
          if (enc instanceof VariableElement ve && ve.getKind() == ElementKind.ENUM_CONSTANT) {
            constVars.add(getExternalStaticFieldGoName(te, ve));
          }
        }
        String qn = te.getQualifiedName().toString();
        String binName = resolver.elements.getBinaryName(te).toString();
        String sliceLit = "[]any{" + String.join(", ", constVars) + "}";
        enumRegInit
            .append("\tRegisterEnumConstants(")
            .append(GoNames.goStringLiteral(qn))
            .append(", ")
            .append(sliceLit)
            .append(")\n");
        if (!binName.equals(qn)) {
          enumRegInit
              .append("\tRegisterEnumConstants(")
              .append(GoNames.goStringLiteral(binName))
              .append(", ")
              .append(sliceLit)
              .append(")\n");
        }
      }
    }
    if (!enumRegInit.isEmpty()) {
      body.append("func init() {\n").append(enumRegInit).append("}\n\n");
    }

    // Emit lambda adapters for external functional interfaces
    List<TypeElement> lambdaIfaces = new ArrayList<>(resolver.getLambdaTargetInterfaces());
    lambdaIfaces.sort(Comparator.comparing(t -> t.getQualifiedName().toString()));
    for (TypeElement fi : lambdaIfaces) {
      if (!resolver.isSourceType(fi) && !mapsToAny(fi)) {
        emitExternalLambdaAdapter(fi, imports, body);
      }
    }
    emitNodePredicateAdapters(imports, body);

    // Emit core iterator & collection wrapper helpers
    emitCollectionHelpers(imports, body);

    StringBuilder full = new StringBuilder(body.length() + 1024);
    full.append("// Auto-generated external type & method bindings for j2gort.\n");
    full.append("package j2gort\n\n");
    full.append("import (\n");
    full.append("\t\"fmt\"\n");
    full.append("\t\"math\"\n");
    full.append("\t\"math/big\"\n");
    full.append("\t\"math/bits\"\n");
    full.append("\t\"reflect\"\n");
    full.append("\t\"regexp\"\n");
    full.append("\t\"sort\"\n");
    full.append("\t\"strconv\"\n");
    full.append("\t\"strings\"\n");
    full.append("\t\"sync\"\n");
    full.append("\t\"sync/atomic\"\n");
    full.append("\t\"unique\"\n\n");
    full.append("\t\"github.com/paullepoulpe/jsgo/transpiled/shims/gson\"\n");
    full.append("\t\"github.com/paullepoulpe/jsgo/transpiled/shims/guava\"\n");
    full.append("\t\"github.com/paullepoulpe/jsgo/transpiled/shims/javaio\"\n");
    full.append("\t\"github.com/paullepoulpe/jsgo/transpiled/shims/javalang\"\n");
    full.append("\t\"github.com/paullepoulpe/jsgo/transpiled/shims/javamath\"\n");
    full.append("\t\"github.com/paullepoulpe/jsgo/transpiled/shims/javanio\"\n");
    full.append("\t\"github.com/paullepoulpe/jsgo/transpiled/shims/javaregex\"\n");
    full.append("\t\"github.com/paullepoulpe/jsgo/transpiled/shims/javautil\"\n");
    full.append("\t\"github.com/paullepoulpe/jsgo/transpiled/shims/javazip\"\n");
    full.append("\t\"github.com/paullepoulpe/jsgo/transpiled/shims/protobuf\"\n");
    full.append("\tshimsres \"github.com/paullepoulpe/jsgo/transpiled/shims/resources\"\n");
    full.append(")\n\n");
    full.append(
        """
        var (
        \t_ = fmt.Sprint
        \t_ = math.NaN
        \t_ = big.NewInt
        \t_ = bits.Reverse32
        \t_ = reflect.TypeOf
        \t_ = regexp.Compile
        \t_ = sort.SliceStable
        \t_ = strconv.Itoa
        \t_ = strings.Join
        \t_ = (*sync.Mutex)(nil)
        \t_ = (*atomic.Bool)(nil)
        \t_ = unique.Make[string]
        \t_ = gson.NewGson
        \t_ = guava.Ascii_ToLowerCase
        \t_ = javaio.System_Out
        \t_ = javalang.ObjectsEquals
        \t_ = javamath.BigInteger_ZERO
        \t_ = javanio.StandardCharsets_UTF_8
        \t_ = javaregex.Pattern_Quote
        \t_ = javautil.NewBitSet
        \t_ = javazip.Deflater_DEFAULT_COMPRESSION
        \t_ = protobuf.ByteString_EMPTY
        \t_ = shimsres.ResourceExists
        )

        """);
    full.append(body);
    return full.toString();
  }

  public String getExternalStaticFieldGoName(TypeElement owner, VariableElement ve) {
    return extGoTypeName(owner) + "_F_" + GoNames.exported(ve.getSimpleName().toString());
  }

  public String goTypeInJ2gort(TypeMirror tm, PackagePlanner.ImportTracker imports) {
    if (tm == null || tm.getKind() == TypeKind.VOID || tm.getKind() == TypeKind.NONE) {
      return "";
    }
    TypeMirror er = resolver.safeErasure(tm);
    if (er.getKind() == TypeKind.ARRAY) {
      return "[]"
          + goTypeInJ2gort(((javax.lang.model.type.ArrayType) er).getComponentType(), imports);
    }
    if (er.getKind() == TypeKind.DECLARED) {
      TypeElement te = (TypeElement) ((javax.lang.model.type.DeclaredType) er).asElement();
      if (resolver.isSourceType(te)) {
        // External library signatures cannot directly reference user packages
        return "any";
      }
      if (!mapsToAny(te)
          && !isShimException(te)
          && !isBoxedPrimitive(te)
          && !te.getQualifiedName().contentEquals("java.lang.String")
          && !te.getQualifiedName().contentEquals("java.lang.Class")
          && !te.getQualifiedName().contentEquals("java.lang.Number")
          && !resolver.getReferencedExternalTypes().contains(te)) {
        return "any";
      }
      return goType(te, imports);
    }
    return resolver.goType(er, imports);
  }

  private void emitExternalType(
      TypeElement te, PackagePlanner.ImportTracker imports, StringBuilder out) {
    String extName = extGoTypeName(te);
    String marker = resolver.getInterfaceMarkerMethod(te);
    String qn = te.getQualifiedName().toString();
    String binName = resolver.elements.getBinaryName(te).toString();
    String simpleName = te.getSimpleName().toString();

    if (te.getKind() == ElementKind.INTERFACE) {
      out.append("// External interface ").append(qn).append('\n');
      out.append("type ").append(extName).append(" interface {\n");
      out.append('\t').append(marker).append("()\n");
      for (TypeMirror im : te.getInterfaces()) {
        if (im instanceof javax.lang.model.type.DeclaredType dt
            && dt.asElement() instanceof TypeElement ste
            && !mapsToAny(ste)
            && resolver.getReferencedExternalTypes().contains(ste)) {
          out.append('\t').append(goTypeInJ2gort(ste.asType(), imports)).append('\n');
        }
      }
      Set<String> emittedMethodNames = new HashSet<>();
      for (Element m : te.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)
            && !ee.getModifiers().contains(Modifier.PRIVATE)) {
          TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
          if (sig != null
              && resolver.isNeededExternalSig(sig)
              && emittedMethodNames.add(sig.goName)) {
            if (!isMethodInheritedFromEmbeddedExtIface(te, sig.goName)) {
              out.append('\t').append(sig.goName).append(formatExtSig(sig, imports)).append('\n');
            }
          }
        }
      }
      if (qn.equals("com.google.common.collect.Interner")) {
        out.append("\tInternStr(string) string\n");
      }
      if (qn.equals("java.util.Comparator")) {
        out.append("\tIsJavaNil_() bool\n");
      }
      out.append("}\n");
      out.append("func ToIface_")
          .append(extName)
          .append("[T interface {\n\t")
          .append(extName)
          .append("\n\tcomparable\n}](v T) ")
          .append(extName)
          .append(" {\n\tvar z T\n\tif v == z {\n\t\treturn nil\n\t}\n\treturn v\n}\n");
      out.append("func init() {\n\tRegisterTypePredicate(")
          .append(GoNames.goStringLiteral(qn))
          .append(", ")
          .append(GoNames.goStringLiteral(binName))
          .append(", func(v_ any) bool { return InstanceOf[")
          .append(extName)
          .append("](v_) })\n}\n\n");
      return;
    }

    if (te.getKind() == ElementKind.ENUM) {
      out.append("// External enum ").append(qn).append('\n');
      out.append("type ").append(extName).append(" struct {\n\tEnumBase `gson:\"-\"`\n}\n");
      out.append("func (e *")
          .append(extName)
          .append(") IsJavaNil_() bool { return e == nil }\n");
      out.append("func (e *")
          .append(extName)
          .append(") GetEnumBase() *EnumBase {\n\tif e == nil {\n\t\treturn nil\n\t}\n\treturn &e.EnumBase\n}\n");
      out.append("func (e *")
          .append(extName)
          .append(") EnumOrdinal_() int32 {\n\tif e == nil {\n\t\treturn -1\n\t}\n\treturn e.EnumBase.Ordinal_\n}\n");
      out.append("func (*")
          .append(extName)
          .append(") JavaClassInfo_() (string, string) { return ")
          .append(GoNames.goStringLiteral(binName))
          .append(", ")
          .append(GoNames.goStringLiteral(simpleName))
          .append(" }\n");
      out.append("func (*").append(extName).append(") BindExtSelf_(any) {}\n");
      for (TypeElement anc : resolver.getAllAncestorTypes(te)) {
        if (anc.getKind() == ElementKind.INTERFACE && !mapsToAny(anc)) {
          out.append("func (*")
              .append(extName)
              .append(") ")
              .append(resolver.getInterfaceMarkerMethod(anc))
              .append("() {}\n");
        }
      }
      out.append("func init() {\n\tRegisterTypePredicate(")
          .append(GoNames.goStringLiteral(qn))
          .append(", ")
          .append(GoNames.goStringLiteral(binName))
          .append(", func(v_ any) bool { return InstanceOf[*")
          .append(extName)
          .append("](v_) })\n}\n\n");
      emitExternalClassMethods(te, extName, imports, out);
      return;
    }

    if (qn.equals("com.google.protobuf.ByteString")) {
      String pbPkg =
          imports.ext("github.com/paullepoulpe/jsgo/transpiled/shims/protobuf", "protobuf");
      out.append("// External value-struct ").append(qn).append('\n');
      out.append("type ").append(extName).append(" struct {\n");
      out.append("\traw ").append(pbPkg).append(".ByteString `gson:\"-\"`\n");
      out.append("\tvalid bool `gson:\"-\"`\n");
      out.append("}\n\n");
      out.append("func (b ").append(extName).append(") IsJavaNil_() bool { return !b.valid }\n");
      out.append("func (")
          .append(extName)
          .append(") JavaClassInfo_() (string, string) { return ")
          .append(GoNames.goStringLiteral(binName))
          .append(", ")
          .append(GoNames.goStringLiteral(simpleName))
          .append(" }\n");
      out.append("func init() {\n\tRegisterTypePredicate(")
          .append(GoNames.goStringLiteral(qn))
          .append(", ")
          .append(GoNames.goStringLiteral(binName))
          .append(", func(v_ any) bool { return InstanceOf[")
          .append(extName)
          .append("](v_) })\n}\n");
      for (TypeElement anc : resolver.getAllAncestorTypes(te)) {
        if (!mapsToAny(anc)
            && !isShimException(anc)
            && !anc.getQualifiedName().contentEquals("java.lang.Number")) {
          if (anc.getKind() == ElementKind.INTERFACE || isSubclassedExternal(anc)) {
            out.append("func (")
                .append(extName)
                .append(") ")
                .append(resolver.getInterfaceMarkerMethod(anc))
                .append("() {}\n");
          }
        }
      }
      emitExternalValueStructMethods(te, extName, imports, out);
      return;
    }

    boolean subclassed = isSubclassedExternal(te);
    String structName = subclassed ? (extName + "Base") : extName;
    if (subclassed) {
      out.append("// External subclassable class ").append(qn).append('\n');
      out.append("type ").append(extName).append(" interface {\n");
      out.append('\t').append(marker).append("()\n");
      out.append("\tGet").append(extName).append("Base() *").append(structName).append('\n');
      TypeElement sup = resolver.getSuperclassElement(te);
      if (sup != null && !mapsToAny(sup)) {
        out.append('\t').append(goTypeInJ2gort(sup.asType(), imports)).append('\n');
      } else if (resolver.isThrowable(te)) {
        out.append("\tjavalang.IRuntimeException\n");
      }
      for (TypeMirror im : te.getInterfaces()) {
        if (im instanceof javax.lang.model.type.DeclaredType dt
            && dt.asElement() instanceof TypeElement ite
            && !mapsToAny(ite)
            && resolver.getReferencedExternalTypes().contains(ite)) {
          out.append('\t').append(goTypeInJ2gort(ite.asType(), imports)).append('\n');
        }
      }
      Set<String> allMethods = collectAllExtInstanceMethods(te, imports);
      for (String mLine : allMethods) {
        out.append('\t').append(mLine).append('\n');
      }
      out.append("}\n");
      out.append("func ToIface_")
          .append(extName)
          .append("[T interface {\n\t")
          .append(extName)
          .append("\n\tcomparable\n}](v T) ")
          .append(extName)
          .append(" {\n\tvar z T\n\tif v == z {\n\t\treturn nil\n\t}\n\treturn v\n}\n\n");
    }

    boolean needsSelfField =
        subclassed
            && !INHERITED_SELF_COLLECTION_TYPES.contains(qn)
            && !qn.startsWith("com.google.protobuf.AbstractMessage")
            && !qn.startsWith("com.google.protobuf.GeneratedMessage");
    out.append("type ").append(structName).append(" struct {\n");
    if (needsSelfField) {
      out.append("\tSelf_ ").append(extName).append(" `gson:\"-\"`\n");
    }
    TypeElement sup = resolver.getSuperclassElement(te);
    boolean skipEmbeddedSuperStruct =
        qn.equals("com.google.common.collect.ImmutableList")
            || qn.equals("com.google.common.collect.ImmutableList.Builder");
    String supEmbeddedField = null;
    if (sup != null && !mapsToAny(sup) && !skipEmbeddedSuperStruct) {
      String supStruct = goStructType(sup, imports);
      out.append('\t').append(supStruct).append(" `gson:\"-\"`\n");
      int dot = supStruct.lastIndexOf('.');
      supEmbeddedField = dot >= 0 ? supStruct.substring(dot + 1) : supStruct;
    } else if (resolver.isThrowable(te)) {
      out.append("\tjavalang.RuntimeException `gson:\"-\"`\n");
    }
    boolean isRootExt =
        !skipEmbeddedSuperStruct
            && !qn.equals("java.util.BitSet")
            && (sup == null
                || mapsToAny(sup)
                || isShimException(sup)
                || sup.getQualifiedName().contentEquals("java.lang.Number"));
    if (isRootExt) {
      out.append("\tCore_ any `gson:\"-\"`\n");
    }
    if (qn.equals("java.lang.AbstractStringBuilder")) {
      out.append("\tjSBCore_ jStringBuilderCore `gson:\"-\"`\n");
    }
    if (qn.equals("java.util.ArrayDeque")) {
      out.append("\tinlineCore JDequeCore `gson:\"-\"`\n");
    } else if (qn.equals("java.util.ArrayList")) {
      out.append("\tinlineCore JListCore `gson:\"-\"`\n");
    } else if (qn.equals("com.google.common.collect.ImmutableList.Builder")) {
      out.append("\tanyBuf [8]any `gson:\"-\"`\n");
      out.append("\toverflow []any `gson:\"-\"`\n");
      out.append("\tstrBuf [8]string `gson:\"-\"`\n");
      out.append("\tstrOver []string `gson:\"-\"`\n");
      out.append("\tn int32 `gson:\"-\"`\n");
      out.append("\tcapHint int32 `gson:\"-\"`\n");
      out.append("\tmode uint8 `gson:\"-\"`\n");
      out.append("\taliased bool `gson:\"-\"`\n");
    } else if (qn.equals("com.google.common.collect.ImmutableList")) {
      out.append("\titems_ []any `gson:\"-\"`\n");
    } else if (qn.equals("java.util.HashSet")
        || qn.equals("com.google.common.collect.ImmutableSet.Builder")) {
      out.append("\tinlineSet JSetCore `gson:\"-\"`\n");
    } else if (qn.equals("com.google.common.collect.ImmutableSet")) {
      String atomicPkg = imports.std("sync/atomic");
      out.append("\telems_ []any `gson:\"-\"`\n");
      out.append("\tsingleArr_ [1]any `gson:\"-\"`\n");
      out.append("\tstrs_ []string `gson:\"-\"`\n");
      out.append("\tsingleStr_ [1]string `gson:\"-\"`\n");
      out.append("\tidx_ ").append(atomicPkg).append(".Pointer[immSetIdx] `gson:\"-\"`\n");
    } else if (qn.equals("com.google.common.collect.ImmutableSortedSet")
        || qn.equals("com.google.common.collect.ImmutableSortedSet.Builder")) {
      out.append("\tsortedCore_ *JMapCore `gson:\"-\"`\n");
    } else if (qn.equals("com.google.common.collect.ImmutableMap")) {
      out.append("\tinlineCore JMapCore `gson:\"-\"`\n");
      TypeElement isTe = resolver.elements.getTypeElement("com.google.common.collect.ImmutableSet");
      TypeElement ilTe = resolver.elements.getTypeElement("com.google.common.collect.ImmutableList");
      if (isTe != null && ilTe != null) {
        String atomicPkg = imports.std("sync/atomic");
        String isSt = goStructType(isTe, imports);
        String ilSt = goStructType(ilTe, imports);
        out.append("\tkeySetCache_ ").append(atomicPkg).append(".Pointer[").append(isSt).append("] `gson:\"-\"`\n");
        out.append("\tvaluesCache_ ").append(atomicPkg).append(".Pointer[").append(ilSt).append("] `gson:\"-\"`\n");
        out.append("\tentrySetCache_ ").append(atomicPkg).append(".Pointer[").append(isSt).append("] `gson:\"-\"`\n");
      }
    } else if (qn.equals("com.google.common.collect.ImmutableSortedMap")) {
      TypeElement issTe = resolver.elements.getTypeElement("com.google.common.collect.ImmutableSortedSet");
      if (issTe != null && resolver.getReferencedExternalTypes().contains(issTe)) {
        String atomicPkg = imports.std("sync/atomic");
        String issSt = goStructType(issTe, imports);
        out.append("\tsortedKeySetCache_ ").append(atomicPkg).append(".Pointer[").append(issSt).append("] `gson:\"-\"`\n");
      }
    } else if (qn.equals("java.util.HashMap")
        || qn.equals("java.util.IdentityHashMap")
        || qn.equals("java.util.TreeMap")
        || qn.equals("java.util.TreeSet")) {
      out.append("\tinlineCore JMapCore `gson:\"-\"`\n");
    } else if (qn.equals("com.google.common.collect.ImmutableMap.Builder")) {
      out.append("\tinlineBuilder jImmMapBuilderCore `gson:\"-\"`\n");
      out.append("\tinlineCore JMapCore `gson:\"-\"`\n");
    } else if (qn.equals("java.util.EnumSet")) {
      out.append("\tinlineEnumSet jEnumSetCore `gson:\"-\"`\n");
    } else if (qn.equals("java.util.EnumMap")) {
      out.append("\tinlineEnumMap jEnumMapCore `gson:\"-\"`\n");
    } else if (qn.equals("java.util.BitSet")) {
      RuntimeCollectionsBindings.emitExtraStructFields(te, out);
    }
    RuntimeProtobufBindings.emitExtraStructFields(this, te, imports, out);
    for (VariableElement ve : resolver.getReferencedExternalFields()) {
      if (ve.getEnclosingElement().equals(te) && !ve.getModifiers().contains(Modifier.STATIC)) {
        out.append("\tF_")
            .append(GoNames.exported(ve.getSimpleName().toString()))
            .append(' ')
            .append(goTypeInJ2gort(ve.asType(), imports))
            .append('\n');
      }
    }
    out.append("}\n\n");

    out.append("func (b *")
        .append(structName)
        .append(") IsJavaNil_() bool { return b == nil }\n");
    out.append("func (*")
        .append(structName)
        .append(") JavaClassInfo_() (string, string) { return ")
        .append(GoNames.goStringLiteral(binName))
        .append(", ")
        .append(GoNames.goStringLiteral(simpleName))
        .append(" }\n");
    out.append("func (b *").append(structName).append(") BindExtSelf_(s any) {\n\t_ = s\n");
    if (needsSelfField) {
      out.append("\tb.Self_ = s.(").append(extName).append(")\n");
    }
    if (sup != null
        && supEmbeddedField != null
        && !mapsToAny(sup)
        && !isShimException(sup)
        && !sup.getQualifiedName().contentEquals("java.lang.Number")) {
      out.append("\tb.").append(supEmbeddedField).append(".BindExtSelf_(s)\n");
    }
    out.append("}\n");

    String checkType = subclassed ? extName : ("*" + extName);
    out.append("func init() {\n\tRegisterTypePredicate(")
        .append(GoNames.goStringLiteral(qn))
        .append(", ")
        .append(GoNames.goStringLiteral(binName))
        .append(", func(v_ any) bool { return InstanceOf[")
        .append(checkType)
        .append("](v_) })\n}\n");

    boolean overridesLc =
        qn.equals("java.util.ArrayList")
            || qn.equals("com.google.common.collect.ImmutableList.Builder");
    boolean overridesMc =
        qn.equals("java.util.HashMap")
            || qn.equals("com.google.common.collect.ImmutableMap")
            || qn.equals("java.util.IdentityHashMap")
            || qn.equals("java.util.TreeMap")
            || qn.equals("java.util.TreeSet")
            || qn.equals("com.google.common.collect.ImmutableSortedSet");

    if (isRootExt) {
      out.append("func (b *").append(structName).append(") getCore() any { return b.Core_ }\n");
      out.append(
          "func (b *")
              .append(structName)
              .append(
                  ") sb() *strings.Builder {\n\tif b.Core_ == nil { b.Core_ = &strings.Builder{} }\n\treturn b.Core_.(*strings.Builder)\n}\n");
      if (!overridesLc) {
        out.append(
            "func (b *")
                .append(structName)
                .append(
                    ") lc() *JListCore {\n\tif b.Core_ == nil { b.Core_ = NewJListCore() }\n\treturn b.Core_.(*JListCore)\n}\n");
      }
      if (!overridesMc) {
        out.append(
            "func (b *")
                .append(structName)
                .append(
                    ") mc() *JMapCore {\n\tif b.Core_ == nil { b.Core_ = NewJMapCore(false) }\n\treturn b.Core_.(*JMapCore)\n}\n");
      }
    }

    if (qn.equals("java.util.AbstractList")) {
      out.append("func (b *")
          .append(structName)
          .append(") selfList() Ext_java_util_List {\n\tif b.Self_ != nil { return b.Self_.(Ext_java_util_List) }\n\treturn b\n}\n");
    } else if (qn.equals("java.util.AbstractQueue")) {
      out.append("func (b *")
          .append(structName)
          .append(") selfQueue() Ext_java_util_Queue {\n\tif b.Self_ != nil { return b.Self_.(Ext_java_util_Queue) }\n\treturn b\n}\n");
    } else if (qn.equals("com.google.common.collect.ImmutableList")) {
      out.append("func (b *")
          .append(structName)
          .append(") selfImmList() Ext_com_google_common_collect_ImmutableList { return b }\n");
    }

    if (qn.equals("java.lang.AbstractStringBuilder")) {
      out.append("func (b *")
          .append(structName)
          .append(") sbCore() *jStringBuilderCore { return &b.jSBCore_ }\n");
    }
    if (qn.equals("java.util.ArrayDeque")) {
      out.append("func (b *")
          .append(structName)
          .append(") dc() *JDequeCore {\n\tif b.Core_ == nil { return &b.inlineCore }\n\treturn b.Core_.(*JDequeCore)\n}\n");
    } else if (qn.equals("java.util.ArrayList")) {
      out.append("func (b *")
          .append(structName)
          .append(") lc() *JListCore {\n\tif b.Core_ == nil { return &b.inlineCore }\n\treturn b.Core_.(*JListCore)\n}\n");
    } else if (qn.equals("java.util.HashSet")) {
      out.append("func (b *")
          .append(structName)
          .append(") sc() *JSetCore { return &b.inlineSet }\n");
      out.append("func (b *")
          .append(structName)
          .append(") mapBack() Ext_java_util_Map {\n\tif m, ok := b.Core_.(Ext_java_util_Map); ok { return m }\n\treturn nil\n}\n");
    } else if (qn.equals("com.google.common.collect.ImmutableSet.Builder")) {
      out.append("func (b *")
          .append(structName)
          .append(") sc() *JSetCore { return &b.inlineSet }\n");
    } else if (qn.equals("java.util.HashMap")
        || qn.equals("com.google.common.collect.ImmutableMap")) {
      out.append("func (b *")
          .append(structName)
          .append(") mc() *JMapCore {\n\tif b.Core_ == nil { return &b.inlineCore }\n\treturn b.Core_.(*JMapCore)\n}\n");
    } else if (qn.equals("java.util.IdentityHashMap")) {
      out.append("func (b *")
          .append(structName)
          .append(") mc() *JMapCore {\n\tif b.Core_ == nil { b.inlineCore.Identity = true; b.inlineCore.genNonCustom = true; return &b.inlineCore }\n\treturn b.Core_.(*JMapCore)\n}\n");
    } else if (qn.equals("java.util.TreeMap") || qn.equals("java.util.TreeSet")) {
      out.append("func (b *")
          .append(structName)
          .append(") mc() *JMapCore {\n\tif b.Core_ == nil {\n\t\tif b.inlineCore.mode == jMapModeEmpty {\n\t\t\tb.inlineCore.mode = jMapModeSortedNat\n\t\t}\n\t\treturn &b.inlineCore\n\t}\n\treturn b.Core_.(*JMapCore)\n}\n");
    } else if (qn.equals("com.google.common.collect.ImmutableSortedSet")) {
      out.append("func (b *")
          .append(structName)
          .append(") mc() *JMapCore {\n\tif b.sortedCore_ == nil {\n\t\tb.sortedCore_ = NewSortedJMapCore(nil)\n\t}\n\treturn b.sortedCore_\n}\n");
    } else if (qn.equals("com.google.common.collect.ImmutableMap.Builder")) {
      out.append("func (b *")
          .append(structName)
          .append(") imbCore() *jImmMapBuilderCore {\n\tif b.Core_ == nil {\n\t\tb.inlineBuilder.mc = &b.inlineCore\n\t\tb.Core_ = &b.inlineBuilder\n\t\treturn &b.inlineBuilder\n\t}\n\treturn b.Core_.(*jImmMapBuilderCore)\n}\n");
    } else if (qn.equals("java.util.EnumSet")) {
      out.append("func (b *")
          .append(structName)
          .append(") esCore() *jEnumSetCore { return &b.inlineEnumSet }\n");
    } else if (qn.equals("java.util.EnumMap")) {
      out.append("func (b *")
          .append(structName)
          .append(") emCore() *jEnumMapCore { return &b.inlineEnumMap }\n");
    }

    if (subclassed) {
      out.append("func (b *").append(structName).append(") ").append(marker).append("() {}\n");
      out.append("func (b *")
          .append(structName)
          .append(") Get")
          .append(extName)
          .append("Base() *")
          .append(structName)
          .append(" { return b }\n");
    }
    for (TypeElement anc : resolver.getAllAncestorTypes(te)) {
      if (!mapsToAny(anc)
          && !isShimException(anc)
          && !anc.getQualifiedName().contentEquals("java.lang.Number")) {
        if (anc.getKind() == ElementKind.INTERFACE || isSubclassedExternal(anc)) {
          out.append("func (b *")
              .append(structName)
              .append(") ")
              .append(resolver.getInterfaceMarkerMethod(anc))
              .append("() {}\n");
          if (skipEmbeddedSuperStruct
              && anc.getKind() == ElementKind.CLASS
              && isSubclassedExternal(anc)
              && resolver.getReferencedExternalTypes().contains(anc)) {
            String ancExt = extGoTypeName(anc);
            String ancSt = goStructType(anc, imports);
            out.append("func (b *")
                .append(structName)
                .append(") Get")
                .append(ancExt)
                .append("Base() *")
                .append(ancSt)
                .append(" { return nil }\n");
          }
        }
      } else if (anc.getQualifiedName().contentEquals("java.lang.Number")) {
        out.append("func (b *").append(structName).append(") Implements_j2gort_Number() {}\n");
      }
    }
    emitExternalClassMethods(te, structName, imports, out);
    RuntimeCollectionsBindings.emitFastMapSetMethods(
        this, te, structName, subclassed, extName, out);
  }

  private boolean isMethodInheritedFromEmbeddedExtIface(TypeElement iface, String goMethodName) {
    for (TypeElement anc : resolver.getAllAncestorTypes(iface)) {
      if (anc.getKind() == ElementKind.INTERFACE
          && !mapsToAny(anc)
          && resolver.getReferencedExternalTypes().contains(anc)) {
        for (Element m : anc.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && !ee.getModifiers().contains(Modifier.STATIC)
              && !ee.getModifiers().contains(Modifier.PRIVATE)) {
            TypeResolver.MethodSig s = resolver.getInstanceMethodSig(ee);
            if (s != null && resolver.isNeededExternalSig(s) && goMethodName.equals(s.goName)) {
              return true;
            }
          }
        }
      }
    }
    return false;
  }

  private Set<String> collectAllExtInstanceMethods(
      TypeElement te, PackagePlanner.ImportTracker imports) {
    Set<String> lines = new LinkedHashSet<>();
    Set<String> seenNames = new HashSet<>();
    for (TypeElement anc : resolver.getAllAncestorTypes(te)) {
      if (!mapsToAny(anc)
          && (isShimException(anc) || resolver.getReferencedExternalTypes().contains(anc))) {
        for (Element m : anc.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && !ee.getModifiers().contains(Modifier.STATIC)
              && !ee.getModifiers().contains(Modifier.PRIVATE)) {
            TypeResolver.MethodSig s = resolver.getInstanceMethodSig(ee);
            if (s != null && resolver.isNeededExternalSig(s)) {
              seenNames.add(s.goName);
            }
          }
        }
      }
    }
    if (resolver.isThrowable(te)) {
      seenNames.add("GetMessage");
      seenNames.add("GetCause");
      seenNames.add("Error");
    }
    for (Element m : te.getEnclosedElements()) {
      if (m instanceof ExecutableElement ee
          && ee.getKind() == ElementKind.METHOD
          && !ee.getModifiers().contains(Modifier.STATIC)
          && !ee.getModifiers().contains(Modifier.PRIVATE)) {
        TypeResolver.MethodSig s = resolver.getInstanceMethodSig(ee);
        if (s != null && resolver.isNeededExternalSig(s) && seenNames.add(s.goName)) {
          lines.add(s.goName + formatExtSig(s, imports));
        }
      }
    }
    String teQn = te.getQualifiedName().toString();
    if ((teQn.equals("java.util.ArrayList")
            || teQn.equals("com.google.common.collect.ImmutableCollection"))
        && seenNames.add("IterSlice_")) {
      lines.add("IterSlice_() []any");
    }
    return lines;
  }

  private void emitExternalClassMethods(
      TypeElement te, String structName, PackagePlanner.ImportTracker imports, StringBuilder out) {
    Set<String> emitted = new HashSet<>();
    if (resolver.isThrowable(te)) {
      emitted.add("GetMessage");
      emitted.add("GetCause");
      emitted.add("Error");
    }
    if (te.getKind() == ElementKind.ENUM) {
      emitted.add("Name");
      emitted.add("Ordinal");
      emitted.add("String");
      emitted.add("HashCode");
      emitted.add("Equals");
      emitted.add("CompareTo");
    }

    List<TypeElement> chain = new ArrayList<>();
    chain.add(te);
    chain.addAll(resolver.getAllAncestorTypes(te));

    for (TypeElement cur : chain) {
      if (mapsToAny(cur) || isShimException(cur)) {
        continue;
      }
      for (Element m : cur.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)
            && !ee.getModifiers().contains(Modifier.PRIVATE)) {
          TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
          if (sig != null && resolver.isNeededExternalSig(sig) && emitted.add(sig.goName)) {
            out.append("func (this_ *")
                .append(structName)
                .append(") ")
                .append(sig.goName)
                .append(formatExtSigWithParamNames(sig, imports))
                .append(" {\n");
            out.append(buildExtInstanceMethodBody(te, sig, imports));
            out.append("}\n");
          }
        }
      }
    }
    out.append('\n');
  }

  private void emitExternalValueStructMethods(
      TypeElement te, String structName, PackagePlanner.ImportTracker imports, StringBuilder out) {
    Set<String> emitted = new HashSet<>();
    List<TypeElement> chain = new ArrayList<>();
    chain.add(te);
    chain.addAll(resolver.getAllAncestorTypes(te));

    for (TypeElement cur : chain) {
      if (mapsToAny(cur) || isShimException(cur)) {
        continue;
      }
      for (Element m : cur.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)
            && !ee.getModifiers().contains(Modifier.PRIVATE)) {
          TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
          if (sig != null && resolver.isNeededExternalSig(sig) && emitted.add(sig.goName)) {
            out.append("func (this_ ")
                .append(structName)
                .append(") ")
                .append(sig.goName)
                .append(formatExtSigWithParamNames(sig, imports))
                .append(" {\n");
            out.append(buildExtInstanceMethodBody(te, sig, imports));
            out.append("}\n");
          }
        }
      }
    }
    if (te.getQualifiedName().contentEquals("com.google.protobuf.ByteString")) {
      if (emitted.add("Equals")) {
        out.append("func (this_ ")
            .append(structName)
            .append(") Equals(arg0 any) bool {\n\treturn byteStringEquals(this_.bs(), arg0)\n}\n");
      }
      if (emitted.add("HashCode")) {
        out.append("func (this_ ")
            .append(structName)
            .append(") HashCode() int32 {\n\treturn this_.bs().HashCode()\n}\n");
      }
      if (emitted.add("String")) {
        out.append("func (this_ ")
            .append(structName)
            .append(") String() string {\n\treturn this_.bs().String()\n}\n");
      }
    }
    out.append('\n');
  }

  private String formatExtSig(TypeResolver.MethodSig sig, PackagePlanner.ImportTracker imports) {
    StringBuilder sb = new StringBuilder("(");
    for (int i = 0; i < sig.canonParamTypes.size(); i++) {
      if (i > 0) sb.append(", ");
      sb.append(goTypeInJ2gort(sig.canonParamTypes.get(i), imports));
    }
    sb.append(')');
    if (sig.canonReturnType != null) {
      String rt = goTypeInJ2gort(sig.canonReturnType, imports);
      if (!rt.isEmpty()) {
        sb.append(' ').append(rt);
      }
    }
    return sb.toString();
  }

  public String formatExtSigWithParamNames(
      TypeResolver.MethodSig sig, PackagePlanner.ImportTracker imports) {
    StringBuilder sb = new StringBuilder("(");
    for (int i = 0; i < sig.canonParamTypes.size(); i++) {
      if (i > 0) sb.append(", ");
      sb.append("arg")
          .append(i)
          .append(' ')
          .append(goTypeInJ2gort(sig.canonParamTypes.get(i), imports));
    }
    sb.append(')');
    if (sig.canonReturnType != null) {
      String rt = goTypeInJ2gort(sig.canonReturnType, imports);
      if (!rt.isEmpty()) {
        sb.append(' ').append(rt);
      }
    }
    return sb.toString();
  }

  private String buildExtInstanceMethodBody(
      TypeElement te, TypeResolver.MethodSig sig, PackagePlanner.ImportTracker imports) {
    String qn = te.getQualifiedName().toString();
    String m = sig.simpleName;
    int arity = sig.canonParamTypes.size();
    String retType =
        sig.canonReturnType != null ? goTypeInJ2gort(sig.canonReturnType, imports) : "";

    // 1. Protobuf messages, builders, enums, ByteString, CodedInputStream, Descriptors
    String protoBody =
        RuntimeProtobufBindings.tryBuildInstanceMethodBody(this, te, sig, imports, retType);
    if (protoBody != null) {
      return protoBody;
    }

    // 2. Collections, Streams, Multimaps, Multisets, Tables, Ranges, Caches, Executors
    String colBody =
        RuntimeCollectionsBindings.tryBuildInstanceMethodBody(this, te, sig, imports, retType);
    if (colBody != null) {
      return colBody;
    }

    // 3. IO, NIO, Zip, Regex, Guava Base/Escape/Hash/IO, Gson, XML/SAX, Args4j, JDK misc
    String ioBody =
        RuntimeIoNioXmlGsonBindings.tryBuildInstanceMethodBody(this, te, sig, retType, imports);
    if (ioBody != null) {
      return ioBody;
    }

    // 4. java.math.BigInteger backed by transpiled/shims/javamath
    if (qn.equals("java.math.BigInteger")) {
      if ((m.equals("add")
              || m.equals("subtract")
              || m.equals("multiply")
              || m.equals("divide")
              || m.equals("remainder")
              || m.equals("mod")
              || m.equals("and")
              || m.equals("or")
              || m.equals("xor"))
          && arity == 1) {
        return "\treturn wrapBigInteger(this_.bi()." + GoNames.exported(m) + "(arg0.bi()))\n";
      }
      if ((m.equals("negate") || m.equals("abs") || m.equals("not")) && arity == 0) {
        return "\treturn wrapBigInteger(this_.bi()." + GoNames.exported(m) + "())\n";
      }
      if (m.equals("divideAndRemainder") && arity == 1) {
        return "\tqr := this_.bi().DivideAndRemainder(arg0.bi())\n\treturn []*Ext_java_math_BigInteger{wrapBigInteger(qr[0]), wrapBigInteger(qr[1])}\n";
      }
      if ((m.equals("pow") || m.equals("shiftLeft") || m.equals("shiftRight")) && arity == 1) {
        return "\treturn wrapBigInteger(this_.bi()." + GoNames.exported(m) + "(int(arg0)))\n";
      }
      if (m.equals("signum") && arity == 0) return "\treturn int32(this_.bi().Signum())\n";
      if (m.equals("bitLength") && arity == 0) return "\treturn int32(this_.bi().BitLength())\n";
      if (m.equals("compareTo") && arity == 1) {
        String p0 = goTypeInJ2gort(sig.canonParamTypes.get(0), imports);
        String otherExpr =
            p0.equals("any") ? "Cast[*Ext_java_math_BigInteger](arg0).bi()" : "arg0.bi()";
        return "\treturn int32(this_.bi().CompareTo(" + otherExpr + "))\n";
      }
      if (m.equals("equals") && arity == 1) {
        return "\to, ok := arg0.(*Ext_java_math_BigInteger)\n\treturn ok && o != nil && this_ != nil && this_.bi().Equals(o.bi())\n";
      }
      if (m.equals("hashCode") && arity == 0) return "\treturn this_.bi().HashCode()\n";
      if (m.equals("toString") && arity == 0) return "\treturn this_.bi().String()\n";
      if (m.equals("toString") && arity == 1)
        return "\treturn this_.bi().ToStringRadix(int(arg0))\n";
      if (m.equals("toByteArray") && arity == 0) {
        return "\traw := this_.bi().ToByteArray()\n\tout := make([]int8, len(raw))\n\tfor i, v := range raw { out[i] = int8(v) }\n\treturn out\n";
      }
      if (m.equals("intValue") && arity == 0) return "\treturn this_.bi().IntValue()\n";
      if (m.equals("intValueExact") && arity == 0) return "\treturn this_.bi().IntValueExact()\n";
      if (m.equals("longValue") && arity == 0) return "\treturn this_.bi().LongValue()\n";
      if (m.equals("longValueExact") && arity == 0) return "\treturn this_.bi().LongValueExact()\n";
      if (m.equals("doubleValue") && arity == 0) return "\treturn this_.bi().DoubleValue()\n";
      if (m.equals("floatValue") && arity == 0)
        return "\treturn float32(this_.bi().DoubleValue())\n";
      if (m.equals("byteValue") && arity == 0) return "\treturn int8(this_.bi().IntValue())\n";
      if (m.equals("shortValue") && arity == 0) return "\treturn int16(this_.bi().IntValue())\n";
    }

    // 5. Generic Enum compareTo / equals fallback
    if (te.getKind() == ElementKind.ENUM) {
      if (m.equals("compareTo") && arity == 1) {
        return "\treturn EnumCompareTo(this_, arg0)\n";
      }
      if (m.equals("equals") && arity == 1) {
        return "\treturn RefEquals(this_, arg0)\n";
      }
    }

    // 6. Generic Object equals / hashCode / toString on external wrapper classes if not specialized
    if (m.equals("equals") && arity == 1) {
      return "\treturn RefEquals(this_, arg0)\n";
    }
    if (m.equals("hashCode") && arity == 0) {
      return "\treturn javalang.IdentityHashCode(this_)\n";
    }
    if (m.equals("toString") && arity == 0) {
      return "\treturn " + GoNames.goStringLiteral(qn) + "\n";
    }

    // Fail-fast loud panic on any unwired external instance method!
    return "\tpanic("
        + GoNames.goStringLiteral("UNIMPLEMENTED_EXT: " + qn + "." + sig.goName)
        + ")\n";
  }

  private void emitExternalConstructor(
      TypeElement owner,
      ExecutableElement ctor,
      PackagePlanner.ImportTracker imports,
      StringBuilder out) {
    String fnName = resolver.getConstructorGoName(ctor);
    String retType = goTypeInJ2gort(owner.asType(), imports);
    String structName = goStructType(owner, imports);
    String qn = owner.getQualifiedName().toString();
    int arity = ctor.getParameters().size();

    out.append("func ").append(fnName).append('(');
    for (int i = 0; i < arity; i++) {
      if (i > 0) out.append(", ");
      out.append("arg")
          .append(i)
          .append(' ')
          .append(goTypeInJ2gort(ctor.getParameters().get(i).asType(), imports));
    }
    out.append(") ").append(retType).append(" {\n");
    out.append("\tobj := &").append(structName).append("{}\n");
    out.append("\tobj.BindExtSelf_(obj)\n");

    boolean handled = false;
    if (RuntimeCollectionsBindings.tryEmitConstructorInit(this, owner, ctor, imports, "obj", out)) {
      handled = true;
    } else if (RuntimeIoNioXmlGsonBindings.tryEmitConstructorBody(
        this, owner, ctor, imports, out)) {
      handled = true;
    } else if (qn.equals("java.lang.StringBuilder") || qn.equals("java.lang.StringBuffer")) {
      if (arity == 0) {
        handled = true;
      } else if (arity == 1) {
        String p0 = goTypeInJ2gort(ctor.getParameters().get(0).asType(), imports);
        if (p0.equals("int32")) {
          out.append("\tif arg0 < 0 { panic(javalang.NewRuntimeExceptionMsg(\"NegativeArraySizeException\")) }\n\tobj.sbCore().grow(int(arg0))\n");
          handled = true;
        } else if (p0.equals("string")) {
          out.append("\tobj.sbCore().appendString(CheckNotNullString(arg0))\n");
          handled = true;
        } else {
          out.append("\tobj.sbCore().appendString(StrVal(CheckNotNull(arg0)))\n");
          handled = true;
        }
      }
    } else if (qn.equals("java.math.BigInteger")) {
      if (arity == 1) {
        String p0 = goTypeInJ2gort(ctor.getParameters().get(0).asType(), imports);
        if (p0.equals("string")) {
          out.append("\tobj.Core_ = javamath.NewBigInteger(CheckNotNull(arg0))\n");
          handled = true;
        } else if (p0.equals("[]int8")) {
          out.append(
              "\tCheckNotNull(arg0)\n\traw := make([]byte, len(arg0))\n\tfor i, v := range arg0 { raw[i] = byte(v) }\n\tobj.Core_ = javamath.NewBigIntegerBytes(raw)\n");
          handled = true;
        }
      } else if (arity == 2) {
        String p0 = goTypeInJ2gort(ctor.getParameters().get(0).asType(), imports);
        if (p0.equals("string")) {
          out.append("\tobj.Core_ = javamath.NewBigIntegerRadix(CheckNotNull(arg0), int(arg1))\n");
          handled = true;
        }
      }
    } else if (resolver.isThrowable(owner)) {
      out.append("\tInitUserThrowable(obj, ").append(GoNames.goStringLiteral(qn));
      for (int i = 0; i < arity; i++) {
        out.append(", arg").append(i);
      }
      out.append(")\n");
      handled = true;
    }

    if (!handled) {
      out.append("\tpanic(")
          .append(GoNames.goStringLiteral("UNIMPLEMENTED_EXT: new " + qn + "/" + arity))
          .append(")\n");
    } else {
      out.append("\treturn obj\n");
    }
    out.append("}\n\n");

    // Emit InitExt_ helper for external superclasses extended by Closure classes
    if (isSubclassedExternal(owner)) {
      String initName = "Init" + fnName.substring(3);
      out.append("func ").append(initName).append("(this_ *").append(structName);
      for (int i = 0; i < arity; i++) {
        out.append(", arg")
            .append(i)
            .append(' ')
            .append(goTypeInJ2gort(ctor.getParameters().get(i).asType(), imports));
      }
      out.append(") {\n");
      boolean initHandled = false;
      if (RuntimeCollectionsBindings.tryEmitConstructorInit(
          this, owner, ctor, imports, "this_", out)) {
        initHandled = true;
      } else if (RuntimeIoNioXmlGsonBindings.tryEmitInitBody(this, owner, ctor, imports, out)) {
        initHandled = true;
      } else if (arity == 0 && owner.getModifiers().contains(Modifier.ABSTRACT)) {
        // No-arg super() constructor on abstract base class (e.g. AbstractSet, InputStream, Writer)
        initHandled = true;
      } else if (resolver.isThrowable(owner)) {
        out.append("\tInitUserThrowable(this_, ").append(GoNames.goStringLiteral(qn));
        for (int i = 0; i < arity; i++) {
          out.append(", arg").append(i);
        }
        out.append(")\n");
        initHandled = true;
      }
      if (!initHandled) {
        out.append("\tpanic(")
            .append(GoNames.goStringLiteral("UNIMPLEMENTED_EXT: super " + qn + "/" + arity))
            .append(")\n");
      }
      out.append("}\n\n");
    }
  }

  private void emitExternalStaticMethod(
      TypeElement owner,
      ExecutableElement m,
      PackagePlanner.ImportTracker imports,
      StringBuilder out) {
    String fnName = resolver.getStaticMethodGoName(m);
    String retType = goTypeInJ2gort(m.getReturnType(), imports);
    String qn = owner.getQualifiedName().toString();
    String simple = m.getSimpleName().toString();
    int arity = m.getParameters().size();

    out.append("func ").append(fnName).append('(');
    for (int i = 0; i < arity; i++) {
      if (i > 0) out.append(", ");
      out.append("arg")
          .append(i)
          .append(' ')
          .append(goTypeInJ2gort(m.getParameters().get(i).asType(), imports));
    }
    out.append(')');
    if (!retType.isEmpty()) {
      out.append(' ').append(retType);
    }
    out.append(" {\n");

    // Note: RuntimeProtobufBindings.tryEmitStaticMethod and RuntimeCollectionsBindings.tryEmitStaticMethod
    // append the closing "}\n\n" themselves when returning true.
    if (tryEmitPredicateStaticMethod(qn, simple, arity, m, imports, out)) {
      return;
    }
    if (RuntimeProtobufBindings.tryEmitStaticMethod(this, owner, m, imports, retType, out)) {
      return;
    }
    if (RuntimeCollectionsBindings.tryEmitStaticMethod(this, owner, m, imports, retType, out)) {
      return;
    }
    if (RuntimeIoNioXmlGsonBindings.tryEmitStaticMethodBody(
        this, owner, m, retType, imports, out)) {
      out.append("}\n\n");
      return;
    }

    if (qn.equals("java.math.BigInteger") && simple.equals("valueOf") && arity == 1) {
      out.append("\treturn wrapBigInteger(javamath.BigInteger_ValueOf(arg0))\n}\n\n");
      return;
    }

    String p0Type = arity >= 1 ? goTypeInJ2gort(m.getParameters().get(0).asType(), imports) : "";
    if (tryEmitJavaLangStaticMethod(qn, simple, arity, p0Type, out)) {
      out.append("}\n\n");
      return;
    }

    if (qn.equals("com.google.common.collect.Interners")
        && (simple.equals("newWeakInterner") || simple.equals("newStrongInterner"))
        && arity == 0) {
      out.append(
          simple.equals("newWeakInterner")
              ? "\treturn newWeakSimpleInterner()\n}\n\n"
              : "\treturn newSimpleInterner()\n}\n\n");
      return;
    }

    if (owner.getKind() == ElementKind.ENUM) {
      if (simple.equals("values") && arity == 0) {
        String elemType = retType.startsWith("[]") ? retType.substring(2) : retType;
        out.append("\traw := LookupEnumConstantsByName(")
            .append(GoNames.goStringLiteral(qn))
            .append(")\n\tout := make(")
            .append(retType)
            .append(", len(raw))\n\tfor i, v := range raw { out[i] = v.(")
            .append(elemType)
            .append(") }\n\treturn out\n}\n\n");
        return;
      }
      if (simple.equals("valueOf") && arity == 1) {
        out.append("\tCheckNotNull(arg0)\n\tfor _, v := range LookupEnumConstantsByName(")
            .append(GoNames.goStringLiteral(qn))
            .append(") {\n\t\tif EnumName(v) == arg0 { return v.(")
            .append(retType)
            .append(") }\n\t}\n\tpanic(javalang.NewIllegalArgumentException(\"No enum constant ")
            .append(qn)
            .append(".\" + arg0))\n}\n\n");
        return;
      }
    }

    out.append("\tpanic(")
        .append(GoNames.goStringLiteral("UNIMPLEMENTED_EXT: " + qn + "." + fnName))
        .append(")\n}\n\n");
  }

  private boolean tryEmitPredicateStaticMethod(
      String qn,
      String simple,
      int arity,
      ExecutableElement m,
      PackagePlanner.ImportTracker imports,
      StringBuilder out) {
    if (qn.equals("com.google.common.base.Predicates")) {
      switch (simple) {
        case "alwaysTrue" -> {
          if (arity == 0) {
            out.append("\treturn predAlwaysTrueIface_\n}\n\n");
            return true;
          }
        }
        case "alwaysFalse" -> {
          if (arity == 0) {
            out.append("\treturn predAlwaysFalseIface_\n}\n\n");
            return true;
          }
        }
        case "isNull" -> {
          if (arity == 0) {
            out.append("\treturn predIsNullIface_\n}\n\n");
            return true;
          }
        }
        case "notNull" -> {
          if (arity == 0) {
            out.append("\treturn predNotNullIface_\n}\n\n");
            return true;
          }
        }
        case "equalTo" -> {
          if (arity == 1) {
            out.append("\treturn guavaPredEqualTo(arg0)\n}\n\n");
            return true;
          }
        }
        case "not" -> {
          if (arity == 1) {
            out.append("\treturn guavaPredNot(arg0)\n}\n\n");
            return true;
          }
        }
        case "in" -> {
          if (arity == 1) {
            out.append("\treturn guavaPredInCollection(arg0)\n}\n\n");
            return true;
          }
        }
        case "and" -> {
          if (arity == 2) {
            out.append("\treturn guavaPredAnd2(arg0, arg1)\n}\n\n");
            return true;
          }
          if (arity == 1) {
            String p0 = goTypeInJ2gort(m.getParameters().get(0).asType(), imports);
            if (p0.startsWith("[]")) {
              out.append("\treturn guavaPredAndSlice(arg0)\n}\n\n");
            } else {
              out.append("\treturn guavaPredAndIter(arg0)\n}\n\n");
            }
            return true;
          }
        }
        case "or" -> {
          if (arity == 2) {
            out.append("\treturn guavaPredOr2(arg0, arg1)\n}\n\n");
            return true;
          }
          if (arity == 1) {
            String p0 = goTypeInJ2gort(m.getParameters().get(0).asType(), imports);
            if (p0.startsWith("[]")) {
              out.append("\treturn guavaPredOrSlice(arg0)\n}\n\n");
            } else {
              out.append("\treturn guavaPredOrIter(arg0)\n}\n\n");
            }
            return true;
          }
        }
        case "instanceOf" -> {
          if (arity == 1) {
            out.append("\treturn guavaPredInstanceOf(arg0)\n}\n\n");
            return true;
          }
        }
        case "compose" -> {
          if (arity == 2) {
            out.append("\treturn guavaPredCompose(arg0, arg1)\n}\n\n");
            return true;
          }
        }
        default -> {}
      }
    } else if (qn.equals("java.util.function.Predicate")) {
      if (simple.equals("isEqual") && arity == 1) {
        out.append("\treturn guavaPredEqualTo(arg0)\n}\n\n");
        return true;
      }
      if (simple.equals("not") && arity == 1) {
        out.append("\treturn jdkPredNot(arg0)\n}\n\n");
        return true;
      }
    }
    return false;
  }

  private boolean tryEmitJavaLangStaticMethod(
      String qn, String simple, int arity, String p0Type, StringBuilder out) {
    switch (qn) {
      case "java.lang.Boolean" -> {
        if (simple.equals("valueOf") && arity == 1 && p0Type.equals("bool")) {
          out.append("\treturn ValueOfBoolean(arg0)\n");
          return true;
        }
        if (simple.equals("parseBoolean") && arity == 1) {
          out.append("\treturn strings.EqualFold(arg0, \"true\") && arg0 != NullString\n");
          return true;
        }
      }
      case "java.lang.Byte" -> {
        if (simple.equals("toUnsignedInt") && arity == 1) {
          out.append("\treturn int32(uint8(arg0))\n");
          return true;
        }
        if (simple.equals("valueOf") && arity == 1 && p0Type.equals("int8")) {
          out.append("\treturn ValueOfByte(arg0)\n");
          return true;
        }
      }
      case "java.lang.Character" -> {
        if (simple.equals("digit") && arity == 2) {
          out.append("\treturn int32(javalang.CharacterDigit(uint16(arg0), int(arg1)))\n");
          return true;
        }
        if (simple.equals("getType") && arity == 1) {
          out.append("\treturn javalang.CharacterGetType(int(arg0))\n");
          return true;
        }
        if (simple.equals("isDigit") && arity == 1 && p0Type.equals("uint16")) {
          out.append("\treturn javalang.CharacterIsDigit(arg0)\n");
          return true;
        }
        if (simple.equals("isIdentifierIgnorable") && arity == 1 && p0Type.equals("uint16")) {
          out.append("\treturn javalang.CharacterIsIdentifierIgnorable(arg0)\n");
          return true;
        }
        if (simple.equals("isJavaIdentifierPart") && arity == 1 && p0Type.equals("uint16")) {
          out.append("\treturn javalang.CharacterIsJavaIdentifierPart(arg0)\n");
          return true;
        }
        if (simple.equals("isJavaIdentifierStart") && arity == 1 && p0Type.equals("uint16")) {
          out.append("\treturn javalang.CharacterIsJavaIdentifierStart(arg0)\n");
          return true;
        }
        if (simple.equals("isLetter") && arity == 1 && p0Type.equals("uint16")) {
          out.append("\treturn javalang.CharacterIsLetter(arg0)\n");
          return true;
        }
        if (simple.equals("isLetterOrDigit") && arity == 1 && p0Type.equals("uint16")) {
          out.append("\treturn javalang.CharacterIsLetterOrDigit(arg0)\n");
          return true;
        }
        if (simple.equals("isLowerCase") && arity == 1 && p0Type.equals("uint16")) {
          out.append("\treturn javalang.CharacterIsLowerCase(arg0)\n");
          return true;
        }
        if (simple.equals("isUpperCase") && arity == 1 && p0Type.equals("uint16")) {
          out.append("\treturn javalang.CharacterIsUpperCase(arg0)\n");
          return true;
        }
        if (simple.equals("isSupplementaryCodePoint") && arity == 1) {
          out.append("\treturn javalang.CharacterIsSupplementaryCodePoint(int(arg0))\n");
          return true;
        }
        if (simple.equals("isWhitespace") && arity == 1) {
          out.append("\treturn javalang.CharacterIsWhitespaceCodePoint(int(arg0))\n");
          return true;
        }
        if (simple.equals("toChars") && arity == 1) {
          out.append("\treturn javalang.CharacterToChars(int(arg0))\n");
          return true;
        }
        if (simple.equals("toString") && arity == 1 && p0Type.equals("uint16")) {
          out.append("\treturn javalang.StringFromChar(arg0)\n");
          return true;
        }
        if (simple.equals("toUpperCase") && arity == 1 && p0Type.equals("uint16")) {
          out.append("\treturn javalang.CharacterToUpperCase(arg0)\n");
          return true;
        }
        if (simple.equals("toLowerCase") && arity == 1 && p0Type.equals("uint16")) {
          out.append("\treturn javalang.CharacterToLowerCase(arg0)\n");
          return true;
        }
        if (simple.equals("valueOf") && arity == 1 && p0Type.equals("uint16")) {
          out.append("\treturn ValueOfChar(arg0)\n");
          return true;
        }
      }
      case "java.lang.Double" -> {
        if (simple.equals("compare") && arity == 2) {
          out.append("\treturn int32(javalang.DoubleCompare(arg0, arg1))\n");
          return true;
        }
        if (simple.equals("doubleToLongBits") && arity == 1) {
          out.append("\tif math.IsNaN(arg0) { return 0x7ff8000000000000 }\n\treturn int64(math.Float64bits(arg0))\n");
          return true;
        }
        if (simple.equals("doubleToRawLongBits") && arity == 1) {
          out.append("\treturn int64(math.Float64bits(arg0))\n");
          return true;
        }
        if (simple.equals("isFinite") && arity == 1) {
          out.append("\treturn javalang.DoubleIsFinite(arg0)\n");
          return true;
        }
        if (simple.equals("isInfinite") && arity == 1) {
          out.append("\treturn javalang.DoubleIsInfinite(arg0)\n");
          return true;
        }
        if (simple.equals("isNaN") && arity == 1) {
          out.append("\treturn javalang.DoubleIsNaN(arg0)\n");
          return true;
        }
        if (simple.equals("longBitsToDouble") && arity == 1) {
          out.append("\treturn math.Float64frombits(uint64(arg0))\n");
          return true;
        }
        if (simple.equals("parseDouble") && arity == 1) {
          out.append("\treturn javalang.DoubleParseDouble(CheckNotNull(arg0))\n");
          return true;
        }
        if (simple.equals("toString") && arity == 1) {
          out.append("\treturn javalang.DoubleToString(arg0)\n");
          return true;
        }
        if (simple.equals("valueOf") && arity == 1 && p0Type.equals("float64")) {
          out.append("\treturn ValueOfDouble(arg0)\n");
          return true;
        }
      }
      case "java.lang.Float" -> {
        if (simple.equals("compare") && arity == 2) {
          out.append("\treturn int32(javalang.DoubleCompare(float64(arg0), float64(arg1)))\n");
          return true;
        }
        if (simple.equals("floatToIntBits") && arity == 1) {
          out.append("\tif math.IsNaN(float64(arg0)) { return 0x7fc00000 }\n\treturn int32(math.Float32bits(arg0))\n");
          return true;
        }
        if (simple.equals("floatToRawIntBits") && arity == 1) {
          out.append("\treturn int32(math.Float32bits(arg0))\n");
          return true;
        }
        if (simple.equals("intBitsToFloat") && arity == 1) {
          out.append("\treturn math.Float32frombits(uint32(arg0))\n");
          return true;
        }
        if (simple.equals("isFinite") && arity == 1) {
          out.append("\treturn javalang.DoubleIsFinite(float64(arg0))\n");
          return true;
        }
        if (simple.equals("isInfinite") && arity == 1) {
          out.append("\treturn javalang.DoubleIsInfinite(float64(arg0))\n");
          return true;
        }
        if (simple.equals("isNaN") && arity == 1) {
          out.append("\treturn javalang.DoubleIsNaN(float64(arg0))\n");
          return true;
        }
        if (simple.equals("toString") && arity == 1) {
          out.append("\treturn javalang.FloatToString(arg0)\n");
          return true;
        }
        if (simple.equals("valueOf") && arity == 1 && p0Type.equals("float32")) {
          out.append("\treturn ValueOfFloat(arg0)\n");
          return true;
        }
      }
      case "java.lang.Integer" -> {
        if (simple.equals("bitCount") && arity == 1) {
          out.append("\treturn IntBitCount(arg0)\n");
          return true;
        }
        if (simple.equals("compare") && arity == 2) {
          out.append("\tif arg0 < arg1 { return -1 }\n\tif arg0 > arg1 { return 1 }\n\treturn 0\n");
          return true;
        }
        if (simple.equals("lowestOneBit") && arity == 1) {
          out.append("\treturn IntLowestOneBit(arg0)\n");
          return true;
        }
        if (simple.equals("highestOneBit") && arity == 1) {
          out.append("\treturn IntHighestOneBit(arg0)\n");
          return true;
        }
        if (simple.equals("numberOfLeadingZeros") && arity == 1) {
          out.append("\treturn IntLeadingZeros(arg0)\n");
          return true;
        }
        if (simple.equals("numberOfTrailingZeros") && arity == 1) {
          out.append("\treturn IntTrailingZeros(arg0)\n");
          return true;
        }
        if (simple.equals("parseInt") && arity == 1) {
          out.append("\treturn javalang.IntegerParseInt(CheckNotNull(arg0))\n");
          return true;
        }
        if (simple.equals("parseInt") && arity == 2) {
          out.append("\treturn javalang.IntegerParseIntRadix(CheckNotNull(arg0), int(arg1))\n");
          return true;
        }
        if (simple.equals("parseInt") && arity == 4) {
          out.append("\ts := StrVal(CheckNotNull(arg0))\n\treturn javalang.IntegerParseIntRadix(JStrSub(s, arg1, arg2), int(arg3))\n");
          return true;
        }
        if (simple.equals("reverse") && arity == 1) {
          out.append("\treturn int32(bits.Reverse32(uint32(arg0)))\n");
          return true;
        }
        if (simple.equals("toHexString") && arity == 1) {
          out.append("\treturn javalang.IntegerToHexString(arg0)\n");
          return true;
        }
        if (simple.equals("toString") && arity == 1) {
          out.append("\treturn strconv.FormatInt(int64(arg0), 10)\n");
          return true;
        }
        if (simple.equals("toString") && arity == 2) {
          out.append("\tr := int(arg1)\n\tif r < 2 || r > 36 { r = 10 }\n\treturn strconv.FormatInt(int64(arg0), r)\n");
          return true;
        }
        if (simple.equals("toUnsignedLong") && arity == 1) {
          out.append("\treturn int64(uint32(arg0))\n");
          return true;
        }
        if (simple.equals("valueOf") && arity == 1 && p0Type.equals("int32")) {
          out.append("\treturn ValueOfInt(arg0)\n");
          return true;
        }
      }
      case "java.lang.Long" -> {
        if (simple.equals("bitCount") && arity == 1) {
          out.append("\treturn LongBitCount(arg0)\n");
          return true;
        }
        if (simple.equals("lowestOneBit") && arity == 1) {
          out.append("\treturn LongLowestOneBit(arg0)\n");
          return true;
        }
        if (simple.equals("highestOneBit") && arity == 1) {
          out.append("\treturn LongHighestOneBit(arg0)\n");
          return true;
        }
        if (simple.equals("numberOfLeadingZeros") && arity == 1) {
          out.append("\treturn LongLeadingZeros(arg0)\n");
          return true;
        }
        if (simple.equals("numberOfTrailingZeros") && arity == 1) {
          out.append("\treturn LongTrailingZeros(arg0)\n");
          return true;
        }
        if (simple.equals("parseLong") && arity == 1) {
          out.append("\treturn javalang.LongParseLong(CheckNotNull(arg0))\n");
          return true;
        }
        if (simple.equals("parseLong") && arity == 2) {
          out.append("\treturn javalang.LongParseLongRadix(CheckNotNull(arg0), int(arg1))\n");
          return true;
        }
        if (simple.equals("toHexString") && arity == 1) {
          out.append("\treturn javalang.LongToHexString(arg0)\n");
          return true;
        }
        if (simple.equals("toString") && arity == 1) {
          out.append("\treturn strconv.FormatInt(arg0, 10)\n");
          return true;
        }
        if (simple.equals("toString") && arity == 2) {
          out.append("\tr := int(arg1)\n\tif r < 2 || r > 36 { r = 10 }\n\treturn strconv.FormatInt(arg0, r)\n");
          return true;
        }
        if (simple.equals("valueOf") && arity == 1 && p0Type.equals("int64")) {
          out.append("\treturn ValueOfLong(arg0)\n");
          return true;
        }
      }
      case "java.lang.Math" -> {
        if (simple.equals("abs") && arity == 1) {
          if (p0Type.equals("float64")) {
            out.append("\treturn math.Abs(arg0)\n");
            return true;
          }
          if (p0Type.equals("int32") || p0Type.equals("int64")) {
            out.append("\tif arg0 < 0 { return -arg0 }\n\treturn arg0\n");
            return true;
          }
        }
        if (simple.equals("cbrt") && arity == 1) {
          out.append("\treturn math.Cbrt(arg0)\n");
          return true;
        }
        if (simple.equals("ceil") && arity == 1) {
          out.append("\treturn math.Ceil(arg0)\n");
          return true;
        }
        if (simple.equals("floor") && arity == 1) {
          out.append("\treturn math.Floor(arg0)\n");
          return true;
        }
        if (simple.equals("max") && arity == 2) {
          if (p0Type.equals("float64")) {
            out.append("\tif math.IsNaN(arg0) || math.IsNaN(arg1) { return math.NaN() }\n\tif arg0 == 0 && arg1 == 0 && math.Signbit(arg0) { return arg1 }\n\tif arg0 >= arg1 { return arg0 }\n\treturn arg1\n");
            return true;
          }
          if (p0Type.equals("int32") || p0Type.equals("int64")) {
            out.append("\tif arg0 >= arg1 { return arg0 }\n\treturn arg1\n");
            return true;
          }
        }
        if (simple.equals("min") && arity == 2) {
          if (p0Type.equals("float64")) {
            out.append("\tif math.IsNaN(arg0) || math.IsNaN(arg1) { return math.NaN() }\n\tif arg0 == 0 && arg1 == 0 && math.Signbit(arg1) { return arg1 }\n\tif arg0 <= arg1 { return arg0 }\n\treturn arg1\n");
            return true;
          }
          if (p0Type.equals("int32") || p0Type.equals("int64")) {
            out.append("\tif arg0 <= arg1 { return arg0 }\n\treturn arg1\n");
            return true;
          }
        }
        if (simple.equals("pow") && arity == 2) {
          out.append("\treturn math.Pow(arg0, arg1)\n");
          return true;
        }
        if (simple.equals("round") && arity == 1 && p0Type.equals("float64")) {
          out.append("\tif math.IsNaN(arg0) || arg0 == 0.49999999999999994 { return 0 }\n\tif arg0 <= float64(math.MinInt64) { return math.MinInt64 }\n\tif arg0 >= float64(math.MaxInt64) { return math.MaxInt64 }\n\treturn int64(math.Floor(arg0 + 0.5))\n");
          return true;
        }
        if (simple.equals("signum") && arity == 1 && p0Type.equals("float64")) {
          out.append("\tif math.IsNaN(arg0) || arg0 == 0 { return arg0 }\n\tif arg0 < 0 { return -1.0 }\n\treturn 1.0\n");
          return true;
        }
        if (simple.equals("toIntExact") && arity == 1) {
          out.append("\tv := int32(arg0)\n\tif int64(v) != arg0 { panic(javalang.NewArithmeticException(\"integer overflow\")) }\n\treturn v\n");
          return true;
        }
      }
      case "java.lang.System" -> {
        if (simple.equals("currentTimeMillis") && arity == 0) {
          out.append("\treturn javalang.System_CurrentTimeMillis()\n");
          return true;
        }
        if (simple.equals("nanoTime") && arity == 0) {
          out.append("\treturn javalang.System_NanoTime()\n");
          return true;
        }
        if (simple.equals("identityHashCode") && arity == 1) {
          out.append("\treturn IdentityHash(arg0)\n");
          return true;
        }
      }
      default -> {}
    }
    return false;
  }

  private void emitExternalStaticField(
      TypeElement owner,
      VariableElement ve,
      String varGoName,
      PackagePlanner.ImportTracker imports,
      StringBuilder out) {
    Object constVal = ve.getConstantValue();
    String fieldType = goTypeInJ2gort(ve.asType(), imports);
    if (ve.getKind() == ElementKind.ENUM_CONSTANT) {
      int ord = 0;
      for (Element enc : owner.getEnclosedElements()) {
        if (enc.getKind() == ElementKind.ENUM_CONSTANT) {
          if (enc.equals(ve)) break;
          ord++;
        }
      }
      out.append("var ")
          .append(varGoName)
          .append(" = &")
          .append(extGoTypeName(owner))
          .append("{EnumBase: EnumBase{Name_: ")
          .append(GoNames.goStringLiteral(ve.getSimpleName().toString()))
          .append(", Ordinal_: ")
          .append(ord)
          .append("}}\n");
      return;
    }
    if (constVal != null) {
      out.append("var ")
          .append(varGoName)
          .append(' ')
          .append(fieldType)
          .append(" = ")
          .append(formatConstLiteral(constVal, fieldType))
          .append('\n');
      return;
    }
    if (RuntimeProtobufBindings.tryEmitStaticField(this, owner, ve, varGoName, fieldType, out)) {
      return;
    }
    if (RuntimeCollectionsBindings.tryEmitStaticField(this, owner, ve, varGoName, fieldType, out)) {
      return;
    }
    if (RuntimeIoNioXmlGsonBindings.tryEmitStaticField(
        this, owner, ve, varGoName, fieldType, out)) {
      return;
    }
    if (owner.getQualifiedName().contentEquals("java.math.BigInteger")) {
      String fn = ve.getSimpleName().toString();
      if (fn.equals("ZERO") || fn.equals("ONE") || fn.equals("TWO") || fn.equals("TEN")) {
        out.append("var ")
            .append(varGoName)
            .append(' ')
            .append(fieldType)
            .append(" = wrapBigInteger(javamath.BigInteger_")
            .append(fn)
            .append(")\n");
        return;
      }
    }
    String qn = owner.getQualifiedName().toString();
    String fn = ve.getSimpleName().toString();
    out.append("var ")
        .append(varGoName)
        .append(' ')
        .append(fieldType)
        .append(" = func() ")
        .append(fieldType)
        .append(" { panic(")
        .append(GoNames.goStringLiteral("UNIMPLEMENTED_EXT_FIELD: " + qn + "." + fn))
        .append(") }()\n");
  }

  private String formatConstLiteral(Object val, String goType) {
    if (val instanceof String s) return GoNames.goStringLiteral(s);
    if (val instanceof Character c) return "uint16(" + (int) c + ")";
    if (val instanceof Boolean b) return b.toString();
    if (val instanceof Double d) {
      if (d.isNaN()) return "math.NaN()";
      if (d.isInfinite()) return d > 0 ? "math.Inf(1)" : "math.Inf(-1)";
      return "float64(" + d + ")";
    }
    if (val instanceof Float f) {
      if (f.isNaN()) return "float32(math.NaN())";
      if (f.isInfinite()) return f > 0 ? "float32(math.Inf(1))" : "float32(math.Inf(-1))";
      return "float32(" + f + ")";
    }
    if (val instanceof Long l) {
      if (l == Long.MIN_VALUE) return "int64(-9223372036854775808)";
      return "int64(" + l + ")";
    }
    return goType + "(" + val + ")";
  }

  private void emitExternalLambdaAdapter(
      TypeElement fi, PackagePlanner.ImportTracker imports, StringBuilder out) {
    ExecutableElement sam = resolver.findSamMethod(fi);
    if (sam == null) {
      return;
    }
    TypeResolver.MethodSig samSig = resolver.getInstanceMethodSig(sam);
    String extName = extGoTypeName(fi);
    String adapterStruct = "lambda_" + extName;
    String fnSig = "func" + formatExtSig(samSig, imports);

    out.append("type ")
        .append(adapterStruct)
        .append(" struct {\n\tfn ")
        .append(fnSig)
        .append("\n}\n");
    out.append("func NewLambda_")
        .append(extName)
        .append("(fn ")
        .append(fnSig)
        .append(") ")
        .append(extName)
        .append(" {\n");
    out.append("\treturn &").append(adapterStruct).append("{fn: fn}\n}\n");
    out.append("func (l *")
        .append(adapterStruct)
        .append(") IsJavaNil_() bool { return l == nil }\n");
    out.append("func (l *")
        .append(adapterStruct)
        .append(") ")
        .append(resolver.getInterfaceMarkerMethod(fi))
        .append("() {}\n");
    for (TypeElement anc : resolver.getAllAncestorTypes(fi)) {
      if (anc.getKind() == ElementKind.INTERFACE && !mapsToAny(anc)) {
        out.append("func (l *")
            .append(adapterStruct)
            .append(") ")
            .append(resolver.getInterfaceMarkerMethod(anc))
            .append("() {}\n");
      }
    }

    Set<String> emitted = new HashSet<>();
    emitted.add(samSig.goName);
    out.append("func (l *")
        .append(adapterStruct)
        .append(") ")
        .append(samSig.goName)
        .append(formatExtSigWithParamNames(samSig, imports))
        .append(" {\n\t");
    if (samSig.canonReturnType != null) {
      out.append("return ");
    }
    out.append("l.fn(");
    for (int i = 0; i < samSig.canonParamTypes.size(); i++) {
      if (i > 0) out.append(", ");
      out.append("arg").append(i);
    }
    out.append(")\n}\n");
    if (fi.getQualifiedName().contentEquals("com.google.common.collect.Interner")) {
      out.append("func (l *")
          .append(adapterStruct)
          .append(") InternStr(s string) string {\n\tif s == NullString { ThrowNPE() }\n\treturn CastString(l.")
          .append(samSig.goName)
          .append("(BoxString(s)))\n}\n");
    }

    List<TypeElement> allIfaces = new ArrayList<>();
    allIfaces.add(fi);
    allIfaces.addAll(resolver.getAllAncestorTypes(fi));
    for (TypeElement iface : allIfaces) {
      if (iface.getKind() != ElementKind.INTERFACE || mapsToAny(iface)) continue;
      String ifaceQn = iface.getQualifiedName().toString();
      for (Element m : iface.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)
            && !ee.getModifiers().contains(Modifier.PRIVATE)) {
          TypeResolver.MethodSig s = resolver.getInstanceMethodSig(ee);
          if (s != null && resolver.isNeededExternalSig(s) && emitted.add(s.goName)) {
            out.append("func (l *")
                .append(adapterStruct)
                .append(") ")
                .append(s.goName)
                .append(formatExtSigWithParamNames(s, imports))
                .append(" {\n");
            int arity = s.canonParamTypes.size();
            String p0 = arity >= 1 ? goTypeInJ2gort(s.canonParamTypes.get(0), imports) : "";
            if (s.simpleName.equals("test")
                && samSig.simpleName.equals("apply")
                && arity == 1
                && samSig.canonParamTypes.size() == 1) {
              out.append("\treturn l.").append(samSig.goName).append("(arg0)\n}\n");
            } else if (s.simpleName.equals("apply")
                && samSig.simpleName.equals("test")
                && arity == 1
                && samSig.canonParamTypes.size() == 1) {
              out.append("\treturn l.").append(samSig.goName).append("(arg0)\n}\n");
            } else if (s.simpleName.equals("equals") && arity == 1) {
              out.append("\treturn RefEquals(l, arg0)\n}\n");
            } else if (s.simpleName.equals("hashCode") && arity == 0) {
              out.append("\treturn javalang.IdentityHashCode(l)\n}\n");
            } else if (s.simpleName.equals("toString") && arity == 0) {
              out.append("\treturn javalang.ObjectToString(")
                  .append(GoNames.goStringLiteral(adapterStruct))
                  .append(", l)\n}\n");
            } else if (ifaceQn.equals("java.lang.Iterable")
                && s.simpleName.equals("forEach")
                && arity == 1) {
              out.append("\tExtDefault_Iterable_ForEach(l, arg0)\n}\n");
            } else if (ifaceQn.equals("java.util.Comparator")) {
              switch (s.simpleName) {
                case "reversed" -> out.append("\treturn ExtDefault_Comparator_Reversed(l)\n}\n");
                case "thenComparing" -> {
                  if (arity == 1 && p0.contains("Comparator")) {
                    out.append("\treturn ExtDefault_Comparator_ThenComparing(l, arg0)\n}\n");
                  } else if (arity == 1) {
                    out.append("\treturn ExtDefault_Comparator_ThenComparingKey(l, arg0)\n}\n");
                  } else if (arity == 2) {
                    out.append(
                        "\treturn ExtDefault_Comparator_ThenComparingKeyCmp(l, arg0, arg1)\n}\n");
                  } else {
                    out.append("\tpanic(")
                        .append(
                            GoNames.goStringLiteral(
                                "UNIMPLEMENTED_EXT: lambda "
                                    + fi.getQualifiedName()
                                    + "."
                                    + s.goName))
                        .append(")\n}\n");
                  }
                }
                case "thenComparingInt" ->
                    out.append("\treturn ExtDefault_Comparator_ThenComparingInt(l, arg0)\n}\n");
                case "thenComparingLong" ->
                    out.append("\treturn ExtDefault_Comparator_ThenComparingLong(l, arg0)\n}\n");
                case "thenComparingDouble" ->
                    out.append("\treturn ExtDefault_Comparator_ThenComparingDouble(l, arg0)\n}\n");
                default ->
                    out.append("\tpanic(")
                        .append(
                            GoNames.goStringLiteral(
                                "UNIMPLEMENTED_EXT: lambda "
                                    + fi.getQualifiedName()
                                    + "."
                                    + s.goName))
                        .append(")\n}\n");
              }
            } else if (ifaceQn.equals("java.util.function.Predicate")
                || ifaceQn.equals("com.google.common.base.Predicate")) {
              switch (s.simpleName) {
                case "negate" -> out.append("\treturn ExtDefault_Predicate_Negate(l)\n}\n");
                case "and" -> out.append("\treturn ExtDefault_Predicate_And(l, arg0)\n}\n");
                case "or" -> out.append("\treturn ExtDefault_Predicate_Or(l, arg0)\n}\n");
                default ->
                    out.append("\tpanic(")
                        .append(
                            GoNames.goStringLiteral(
                                "UNIMPLEMENTED_EXT: lambda "
                                    + fi.getQualifiedName()
                                    + "."
                                    + s.goName))
                        .append(")\n}\n");
              }
            } else if (ifaceQn.equals("java.util.function.Function")
                || ifaceQn.equals("com.google.common.base.Function")
                || ifaceQn.equals("java.util.function.UnaryOperator")) {
              switch (s.simpleName) {
                case "andThen" -> out.append("\treturn ExtDefault_Function_AndThen(l, arg0)\n}\n");
                case "compose" -> out.append("\treturn ExtDefault_Function_Compose(l, arg0)\n}\n");
                default ->
                    out.append("\tpanic(")
                        .append(
                            GoNames.goStringLiteral(
                                "UNIMPLEMENTED_EXT: lambda "
                                    + fi.getQualifiedName()
                                    + "."
                                    + s.goName))
                        .append(")\n}\n");
              }
            } else if (ifaceQn.equals("java.util.function.Consumer")) {
              if (s.simpleName.equals("andThen") && arity == 1) {
                out.append("\treturn ExtDefault_Consumer_AndThen(l, arg0)\n}\n");
              } else {
                out.append("\tpanic(")
                    .append(
                        GoNames.goStringLiteral(
                            "UNIMPLEMENTED_EXT: lambda "
                                + fi.getQualifiedName()
                                + "."
                                + s.goName))
                    .append(")\n}\n");
              }
            } else {
              out.append("\tpanic(")
                  .append(
                      GoNames.goStringLiteral(
                          "UNIMPLEMENTED_EXT: lambda " + fi.getQualifiedName() + "." + s.goName))
                  .append(")\n}\n");
            }
          }
        }
      }
    }
    out.append('\n');
  }

  private void emitNodePredicateAdapters(
      PackagePlanner.ImportTracker imports, StringBuilder out) {
    TypeElement gPred = resolver.elements.getTypeElement("com.google.common.base.Predicate");
    TypeElement jPred = resolver.elements.getTypeElement("java.util.function.Predicate");
    boolean hasGPred = gPred != null && resolver.getReferencedExternalTypes().contains(gPred);
    boolean hasJPred = jPred != null && resolver.getReferencedExternalTypes().contains(jPred);
    if (hasGPred) {
      emitSingleNodePredicateAdapter(gPred, imports, out);
    }
    if (hasJPred) {
      emitSingleNodePredicateAdapter(jPred, imports, out);
    }
    String gApplyGo = findExtMethodGoName("com.google.common.base.Predicate", "apply", 1);
    String jTestGo = findExtMethodGoName("java.util.function.Predicate", "test", 1);
    String colContainsGo = findExtMethodGoName("java.util.Collection", "contains", 1);
    TypeElement colElem = resolver.elements.getTypeElement("java.util.Collection");
    TypeResolver.MethodSig colContainsSig = null;
    if (colElem != null) {
      for (Element enc : colElem.getEnclosedElements()) {
        if (enc instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && ee.getSimpleName().contentEquals("contains")
            && ee.getParameters().size() == 1) {
          colContainsSig = resolver.getInstanceMethodSig(ee);
          break;
        }
      }
    }
    out.append("func predColContainsAny_(col Ext_java_util_Collection, v any) bool {\n");
    if (colContainsSig != null && resolver.isNeededExternalSig(colContainsSig)) {
      out.append("\treturn col.").append(colContainsGo).append("(v)\n}\n\n");
    } else {
      out.append("\treturn any(col).(interface{ ")
          .append(colContainsGo)
          .append("(any) bool }).")
          .append(colContainsGo)
          .append("(v)\n}\n\n");
    }

    out.append("func PredTestNode(p any, n uint32) bool {\n")
        .append("\tif p == predAlwaysTrueAny_ {\n\t\treturn true\n\t}\n");
    if (hasGPred) {
      out.append("\tif p == guavaAlwaysTruePred_ {\n\t\treturn true\n\t}\n");
    }
    out.append("\tif p == nil {\n\t\tThrowNPE()\n\t}\n")
        .append("\tif np, ok := p.(NodePredicate); ok {\n")
        .append("\t\tif nc, ok2 := p.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {\n\t\t\tThrowNPE()\n\t\t}\n")
        .append("\t\treturn np.TestNode_(n)\n\t}\n")
        .append("\tif isJavaNilAny(p) {\n\t\tThrowNPE()\n\t}\n");
    if (hasGPred) {
      String gExt = extGoTypeName(gPred);
      ExecutableElement gSam = resolver.findSamMethod(gPred);
      TypeResolver.MethodSig gSig = gSam != null ? resolver.getInstanceMethodSig(gSam) : null;
      out.append("\tif gp, ok := p.(").append(gExt).append("); ok {\n");
      if (gSig != null && resolver.isNeededExternalSig(gSig)) {
        out.append("\t\treturn gp.").append(gApplyGo).append("(boxNodeById(n))\n\t}\n");
      } else {
        out.append("\t\treturn any(gp).(interface{ ")
            .append(gApplyGo)
            .append("(any) bool }).")
            .append(gApplyGo)
            .append("(boxNodeById(n))\n\t}\n");
      }
    }
    if (hasJPred) {
      String jExt = extGoTypeName(jPred);
      ExecutableElement jSam = resolver.findSamMethod(jPred);
      TypeResolver.MethodSig jSig = jSam != null ? resolver.getInstanceMethodSig(jSam) : null;
      out.append("\tif jp, ok := p.(").append(jExt).append("); ok {\n");
      if (jSig != null && resolver.isNeededExternalSig(jSig)) {
        out.append("\t\treturn jp.").append(jTestGo).append("(boxNodeById(n))\n\t}\n");
      } else {
        out.append("\t\treturn any(jp).(interface{ ")
            .append(jTestGo)
            .append("(any) bool }).")
            .append(jTestGo)
            .append("(boxNodeById(n))\n\t}\n");
      }
    }
    out.append("\tpanic(javalang.NewClassCastException(fmt.Sprintf(\"cannot cast %T to Predicate\", p)))\n}\n\n");
  }

  private void emitSingleNodePredicateAdapter(
      TypeElement fi, PackagePlanner.ImportTracker imports, StringBuilder out) {
    ExecutableElement sam = resolver.findSamMethod(fi);
    if (sam == null) {
      return;
    }
    TypeResolver.MethodSig samSig = resolver.getInstanceMethodSig(sam);
    String extName = extGoTypeName(fi);
    String adapterStruct = "nodeLambda_" + extName;

    out.append("type ").append(adapterStruct).append(" struct {\n\tfn func(uint32) bool\n}\n");
    out.append("func NewNodeLambda_")
        .append(extName)
        .append("(fn func(uint32) bool) ")
        .append(extName)
        .append(" {\n\treturn &")
        .append(adapterStruct)
        .append("{fn: fn}\n}\n");
    out.append("func (l *").append(adapterStruct).append(") IsJavaNil_() bool { return l == nil }\n");
    out.append("func (l *").append(adapterStruct).append(") HasIdentityEquals_() {}\n");
    out.append("func (l *")
        .append(adapterStruct)
        .append(") TestNode_(n uint32) bool {\n\tif l == nil {\n\t\tThrowNPE()\n\t}\n\treturn l.fn(n)\n}\n");
    out.append("func (l *")
        .append(adapterStruct)
        .append(") ")
        .append(resolver.getInterfaceMarkerMethod(fi))
        .append("() {}\n");
    for (TypeElement anc : resolver.getAllAncestorTypes(fi)) {
      if (anc.getKind() == ElementKind.INTERFACE && !mapsToAny(anc)) {
        out.append("func (l *")
            .append(adapterStruct)
            .append(") ")
            .append(resolver.getInterfaceMarkerMethod(anc))
            .append("() {}\n");
      }
    }

    Set<String> emitted = new HashSet<>();
    emitted.add(samSig.goName);
    out.append("func (l *")
        .append(adapterStruct)
        .append(") ")
        .append(samSig.goName)
        .append(formatExtSigWithParamNames(samSig, imports))
        .append(" {\n\treturn l.fn(CastNode[uint32](arg0))\n}\n");

    List<TypeElement> allIfaces = new ArrayList<>();
    allIfaces.add(fi);
    allIfaces.addAll(resolver.getAllAncestorTypes(fi));
    for (TypeElement iface : allIfaces) {
      if (iface.getKind() != ElementKind.INTERFACE || mapsToAny(iface)) continue;
      String ifaceQn = iface.getQualifiedName().toString();
      for (Element m : iface.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)
            && !ee.getModifiers().contains(Modifier.PRIVATE)) {
          TypeResolver.MethodSig s = resolver.getInstanceMethodSig(ee);
          if (s == null || !resolver.isNeededExternalSig(s) || !emitted.add(s.goName)) continue;
          int arity = s.canonParamTypes.size();
          out.append("func (l *")
              .append(adapterStruct)
              .append(") ")
              .append(s.goName)
              .append(formatExtSigWithParamNames(s, imports))
              .append(" {\n");
          if (s.simpleName.equals("test") && samSig.simpleName.equals("apply") && arity == 1) {
            out.append("\treturn l.").append(samSig.goName).append("(arg0)\n}\n");
          } else if (s.simpleName.equals("apply") && samSig.simpleName.equals("test") && arity == 1) {
            out.append("\treturn l.").append(samSig.goName).append("(arg0)\n}\n");
          } else if (s.simpleName.equals("equals") && arity == 1) {
            out.append("\treturn RefEquals(l, arg0)\n}\n");
          } else if (s.simpleName.equals("hashCode") && arity == 0) {
            out.append("\treturn javalang.IdentityHashCode(l)\n}\n");
          } else if (s.simpleName.equals("toString") && arity == 0) {
            out.append("\treturn javalang.ObjectToString(")
                .append(GoNames.goStringLiteral(adapterStruct))
                .append(", l)\n}\n");
          } else if (ifaceQn.equals("java.util.function.Predicate")
              || ifaceQn.equals("com.google.common.base.Predicate")) {
            switch (s.simpleName) {
              case "negate" -> out.append("\treturn ExtDefault_Predicate_Negate(l)\n}\n");
              case "and" -> out.append("\treturn ExtDefault_Predicate_And(l, arg0)\n}\n");
              case "or" -> out.append("\treturn ExtDefault_Predicate_Or(l, arg0)\n}\n");
              default ->
                  out.append("\tpanic(")
                      .append(
                          GoNames.goStringLiteral(
                              "UNIMPLEMENTED_EXT: nodeLambda " + fi.getQualifiedName() + "." + s.goName))
                      .append(")\n}\n");
            }
          } else {
            out.append("\tpanic(")
                .append(
                    GoNames.goStringLiteral(
                        "UNIMPLEMENTED_EXT: nodeLambda " + fi.getQualifiedName() + "." + s.goName))
                .append(")\n}\n");
          }
        }
      }
    }
    out.append('\n');
  }

  public String findExtMethodGoName(String ifaceQn, String methodSimpleName, int arity) {
    TypeElement te = resolver.elements.getTypeElement(ifaceQn);
    if (te != null) {
      for (Element m : te.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)
            && ee.getSimpleName().contentEquals(methodSimpleName)
            && ee.getParameters().size() == arity) {
          TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
          if (sig != null) {
            return sig.goName;
          }
        }
      }
    }
    return GoNames.exported(methodSimpleName);
  }

  private void emitCollectionHelpers(PackagePlanner.ImportTracker imports, StringBuilder out) {
    TypeElement iterTe = resolver.elements.getTypeElement("java.util.Iterator");
    if (iterTe != null) {
      out.append("type jListIterObj struct {\n\tlc *JListCore\n\titems []any\n\tstrs []string\n\tidx int32\n}\n");
      out.append("func (it *jListIterObj) IsJavaNil_() bool { return it == nil }\n");
      out.append(
          "func newListIter(lc *JListCore) Ext_java_util_Iterator {\n\tif lc == nil || lc.IsEmpty() { return emptyMapIter_ }\n\treturn &jListIterObj{lc: lc, idx: 0}\n}\n");
      out.append(
          "func newSliceObjIter(items []any) Ext_java_util_Iterator {\n\tif len(items) == 0 { return emptyMapIter_ }\n\treturn &jListIterObj{items: items, idx: 0}\n}\n");
      out.append(
          "func newStrSliceIter(strs []string) Ext_java_util_Iterator {\n\tif len(strs) == 0 { return emptyMapIter_ }\n\treturn &jListIterObj{strs: strs, idx: 0}\n}\n");
      out.append("func (it *jListIterObj) Implements_j2gort_Ext_java_util_Iterator() {}\n");
      out.append("type jSingleObjIter struct {\n\telem any\n\tdone bool\n}\n");
      out.append("func (it *jSingleObjIter) IsJavaNil_() bool { return it == nil }\n");
      out.append(
          "func newSingleObjIter(v any) Ext_java_util_Iterator {\n\treturn &jSingleObjIter{elem: v}\n}\n");
      out.append(
          "func newSingleStrIter(s string) Ext_java_util_Iterator {\n\tif s == NullString { return &jSingleObjIter{} }\n\treturn &jSingleObjIter{elem: s}\n}\n");
      out.append("func (it *jSingleObjIter) Implements_j2gort_Ext_java_util_Iterator() {}\n");
      TypeElement unmodItTe =
          resolver.elements.getTypeElement("com.google.common.collect.UnmodifiableIterator");
      if (unmodItTe != null && resolver.getReferencedExternalTypes().contains(unmodItTe)) {
        out.append(
            "func (it *jListIterObj) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator() {}\n");
        out.append(
            "func (it *jSingleObjIter) Implements_j2gort_Ext_com_google_common_collect_UnmodifiableIterator() {}\n");
        if (isSubclassedExternal(unmodItTe)) {
          out.append(
              "func (it *jListIterObj) GetExt_com_google_common_collect_UnmodifiableIteratorBase() *Ext_com_google_common_collect_UnmodifiableIteratorBase { return nil }\n");
          out.append(
              "func (it *jSingleObjIter) GetExt_com_google_common_collect_UnmodifiableIteratorBase() *Ext_com_google_common_collect_UnmodifiableIteratorBase { return nil }\n");
        }
      }
      for (Element m : iterTe.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)) {
          TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
          if (sig == null || !resolver.isNeededExternalSig(sig)) continue;
          out.append("func (it *jListIterObj) ")
              .append(sig.goName)
              .append(formatExtSigWithParamNames(sig, imports))
              .append(" {\n");
          switch (sig.simpleName) {
            case "hasNext" ->
                out.append(
                    "\tif it.lc != nil { return it.idx < it.lc.Size() }\n\tif it.strs != nil { return int(it.idx) < len(it.strs) }\n\treturn int(it.idx) < len(it.items)\n");
            case "next" ->
                out.append(
                    "\tif it.lc != nil {\n\t\tif it.idx >= it.lc.Size() { panic(javalang.NewNoSuchElementException(\"iterator exhausted\")) }\n\t\tv := it.lc.Get(it.idx)\n\t\tit.idx++\n\t\treturn v\n\t}\n\tif it.strs != nil {\n\t\tif int(it.idx) >= len(it.strs) { panic(javalang.NewNoSuchElementException(\"iterator exhausted\")) }\n\t\ts := it.strs[it.idx]\n\t\tit.idx++\n\t\tif s == NullString { return nil }\n\t\treturn s\n\t}\n\tif int(it.idx) >= len(it.items) { panic(javalang.NewNoSuchElementException(\"iterator exhausted\")) }\n\tv := it.items[it.idx]\n\tit.idx++\n\treturn v\n");
            case "remove" ->
                out.append(
                    "\tif it.lc == nil { panic(javalang.NewUnsupportedOperationException(\"remove\")) }\n\tif it.idx <= 0 { panic(javalang.NewIllegalStateException(\"remove before next\")) }\n\tit.idx--\n\tit.lc.RemoveAt(it.idx)\n");
            case "forEachRemaining" ->
                out.append("\tExtDefault_Iterator_ForEachRemaining(it, arg0)\n");
            default ->
                out.append("\tpanic(")
                    .append(GoNames.goStringLiteral("UNIMPLEMENTED_EXT: jListIterObj." + sig.goName))
                    .append(")\n");
          }
          out.append("}\n");
          out.append("func (it *jSingleObjIter) ")
              .append(sig.goName)
              .append(formatExtSigWithParamNames(sig, imports))
              .append(" {\n");
          switch (sig.simpleName) {
            case "hasNext" -> out.append("\treturn !it.done\n");
            case "next" ->
                out.append(
                    "\tif it.done { panic(javalang.NewNoSuchElementException(\"iterator exhausted\")) }\n\tit.done = true\n\treturn it.elem\n");
            case "remove" ->
                out.append("\tpanic(javalang.NewUnsupportedOperationException(\"remove\"))\n");
            case "forEachRemaining" ->
                out.append("\tExtDefault_Iterator_ForEachRemaining(it, arg0)\n");
            default ->
                out.append("\tpanic(")
                    .append(GoNames.goStringLiteral("UNIMPLEMENTED_EXT: jSingleObjIter." + sig.goName))
                    .append(")\n");
          }
          out.append("}\n");
        }
      }
      out.append('\n');
    }

    TypeElement alTe = resolver.elements.getTypeElement("java.util.ArrayList");
    if (alTe != null) {
      String st = goStructType(alTe, imports);
      String alGo = goTypeInJ2gort(alTe.asType(), imports);
      out.append("func wrapArrayList(items []any) ").append(alGo).append(" {\n");
      out.append("\tobj := &").append(st).append("{}\n");
      out.append("\tobj.BindExtSelf_(obj)\n");
      out.append("\tobj.inlineCore.InitFromSlice(items)\n\treturn obj\n}\n\n");
      out.append("func wrapArrayListOwn(items []any) ").append(alGo).append(" {\n");
      out.append("\tobj := &").append(st).append("{}\n");
      out.append("\tobj.BindExtSelf_(obj)\n");
      out.append("\tobj.inlineCore.InitOwnSlice(items)\n\treturn obj\n}\n\n");
    }

    TypeElement lhsTe = resolver.elements.getTypeElement("java.util.LinkedHashSet");
    if (lhsTe != null) {
      String st = goStructType(lhsTe, imports);
      out.append("func wrapLinkedHashSet(items []any) ")
          .append(goTypeInJ2gort(lhsTe.asType(), imports))
          .append(" {\n");
      out.append("\tobj := &").append(st).append("{}\n");
      out.append("\tobj.BindExtSelf_(obj)\n");
      out.append(
          "\tobj.inlineSet.InitExpectedSize(int32(len(items)))\n\tfor _, x := range items { obj.inlineSet.Add(x) }\n\treturn obj\n}\n\n");
    }

    TypeElement lhmTe = resolver.elements.getTypeElement("java.util.LinkedHashMap");
    if (lhmTe != null) {
      String st = goStructType(lhmTe, imports);
      out.append("func wrapLinkedHashMap(mc *JMapCore) ")
          .append(goTypeInJ2gort(lhmTe.asType(), imports))
          .append(" {\n");
      out.append("\tobj := &").append(st).append("{}\n");
      out.append("\tobj.BindExtSelf_(obj)\n");
      out.append("\tif mc != nil { obj.Core_ = mc }\n\treturn obj\n}\n\n");
    }

    TypeElement ilTe = resolver.elements.getTypeElement("com.google.common.collect.ImmutableList");
    boolean hasImmList = ilTe != null && resolver.getReferencedExternalTypes().contains(ilTe);

    TypeElement listTe = resolver.elements.getTypeElement("java.util.List");
    if (listTe != null && resolver.getReferencedExternalTypes().contains(listTe)) {
      String listGo = goTypeInJ2gort(listTe.asType(), imports);
      out.append("type jSingletonList struct {\n\telem [1]any\n}\n");
      out.append("var emptyElemSingletonList_ = &jSingletonList{}\n");
      out.append("func newSingletonList(v any) ")
          .append(listGo)
          .append(
              " {\n\tv = ToAny(v)\n\tif v == nil { return emptyElemSingletonList_ }\n\treturn &jSingletonList{elem: [1]any{v}}\n}\n");
      out.append("func (l *jSingletonList) IsJavaNil_() bool { return l == nil }\n");
      out.append(
          "func (*jSingletonList) JavaClassInfo_() (string, string) { return \"java.util.Collections$SingletonList\", \"SingletonList\" }\n");
      out.append(
          "func (l *jSingletonList) IterSlice_() []any {\n\tif l == nil { ThrowNPE() }\n\treturn l.elem[:1:1]\n}\n");
      out.append(
          "func (l *jSingletonList) Equals(arg0 any) bool { return listEqualsSlice(l.elem[:1:1], arg0) }\n");
      out.append(
          "func (l *jSingletonList) HashCode() int32 { return 31 + javalang.ObjectHashCode(l.elem[0]) }\n");
      out.append("func (l *jSingletonList) String() string { return \"[\" + StrVal(l.elem[0]) + \"]\" }\n");
      out.append("type jSingletonPtr[T any] struct {\n\tp *T\n}\n");
      out.append("func SingletonListPtr[T any](p *T) ")
          .append(listGo)
          .append(" {\n\tif p == nil { return emptyElemSingletonList_ }\n\treturn jSingletonPtr[T]{p: p}\n}\n");
      out.append("func (l jSingletonPtr[T]) IsJavaNil_() bool { return l.p == nil }\n");
      out.append(
          "func (jSingletonPtr[T]) JavaClassInfo_() (string, string) { return \"java.util.Collections$SingletonList\", \"SingletonList\" }\n");
      out.append(
          "func (l jSingletonPtr[T]) typedViewLen_() (int, bool) {\n\tif l.p == nil { ThrowNPE() }\n\treturn 1, true\n}\n");
      out.append("func (l jSingletonPtr[T]) typedViewAt_(int) any { return l.p }\n");
      out.append("func (jSingletonPtr[T]) typedViewNodeAt_(int) (uint32, bool) { return 0, false }\n");
      out.append(
          "func (l jSingletonPtr[T]) Equals(arg0 any) bool { var a [1]any = [1]any{l.p}; return listEqualsSlice(a[:], arg0) }\n");
      out.append(
          "func (l jSingletonPtr[T]) HashCode() int32 { return 31 + javalang.ObjectHashCode(l.p) }\n");
      out.append("func (l jSingletonPtr[T]) String() string { return \"[\" + StrVal(l.p) + \"]\" }\n");
      Set<String> listMarkers = new LinkedHashSet<>();
      listMarkers.add(resolver.getInterfaceMarkerMethod(listTe));
      for (TypeElement anc : resolver.getAllAncestorTypes(listTe)) {
        if (anc.getKind() == ElementKind.INTERFACE && !mapsToAny(anc)) {
          listMarkers.add(resolver.getInterfaceMarkerMethod(anc));
        }
      }
      listMarkers.add("Implements_j2gort_Ext_java_util_RandomAccess");
      listMarkers.add("Implements_j2gort_Ext_java_io_Serializable");
      for (String mk : listMarkers) {
        out.append("func (*jSingletonList) ").append(mk).append("() {}\n");
        out.append("func (jSingletonPtr[T]) ").append(mk).append("() {}\n");
      }
      Set<String> emittedList = new HashSet<>(Set.of("Equals", "HashCode", "String"));
      List<TypeElement> listChain = new ArrayList<>();
      listChain.add(listTe);
      listChain.addAll(resolver.getAllAncestorTypes(listTe));
      String intFnApply = findExtMethodGoName("java.util.function.IntFunction", "apply", 1);
      String consumerAccept = findExtMethodGoName("java.util.function.Consumer", "accept", 1);
      String emptyListExpr = hasImmList ? "emptyImmutableList_" : "wrapArrayList(nil)";
      for (TypeElement cur : listChain) {
        if (cur.getKind() != ElementKind.INTERFACE || mapsToAny(cur)) continue;
        for (Element m : cur.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && !ee.getModifiers().contains(Modifier.STATIC)
              && !ee.getModifiers().contains(Modifier.PRIVATE)) {
            TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
            if (sig == null || !resolver.isNeededExternalSig(sig) || !emittedList.add(sig.goName)) {
              continue;
            }
            int arity = sig.canonParamTypes.size();
            String p0 = arity >= 1 ? goTypeInJ2gort(sig.canonParamTypes.get(0), imports) : "";
            out.append("func (l *jSingletonList) ")
                .append(sig.goName)
                .append(formatExtSigWithParamNames(sig, imports))
                .append(" {\n");
            switch (sig.simpleName) {
              case "size" -> out.append("\treturn 1\n");
              case "isEmpty" -> out.append("\treturn false\n");
              case "get" -> out.append("\tCheckElementIndex(arg0, 1)\n\treturn l.elem[0]\n");
              case "getFirst", "getLast" -> out.append("\treturn l.elem[0]\n");
              case "contains" -> out.append("\treturn jSetElemEquals(l.elem[0], arg0)\n");
              case "indexOf", "lastIndexOf" ->
                  out.append("\tif jSetElemEquals(l.elem[0], arg0) { return 0 }\n\treturn -1\n");
              case "containsAll" ->
                  out.append(
                      "\treturn colContainsAll(func(x any) bool { return jSetElemEquals(l.elem[0], x) }, arg0)\n");
              case "iterator" -> out.append("\treturn newSingleObjIter(l.elem[0])\n");
              case "listIterator" -> {
                if (arity == 0) {
                  out.append("\treturn newImmListIter(l.elem[:1:1], 0)\n");
                } else {
                  out.append("\treturn newImmListIter(l.elem[:1:1], arg0)\n");
                }
              }
              case "subList" -> {
                out.append("\tCheckPositionIndexes(arg0, arg1, 1)\n\tif arg0 == arg1 { return ")
                    .append(emptyListExpr)
                    .append(" }\n\treturn l\n");
              }
              case "reversed" -> out.append("\treturn l\n");
              case "toArray" -> {
                if (arity == 0) {
                  out.append("\treturn []any{l.elem[0]}\n");
                } else if (p0.startsWith("[]")) {
                  out.append("\treturn sliceToArray1(l.elem[:1:1], arg0)\n");
                } else {
                  out.append("\tCheckNotNull(arg0)\n\treturn sliceToArray1(l.elem[:1:1], Cast[[]any](arg0.")
                      .append(intFnApply)
                      .append("(1)))\n");
                }
              }
              case "forEach" ->
                  out.append("\tCheckNotNull(arg0)\n\targ0.")
                      .append(consumerAccept)
                      .append("(l.elem[0])\n");
              case "stream", "parallelStream" ->
                  out.append("\treturn streamFromSlice(l.elem[:1:1])\n");
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
                  "replaceAll" ->
                  out.append("\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n");
              default ->
                  out.append("\tpanic(")
                      .append(GoNames.goStringLiteral("UNIMPLEMENTED_EXT: jSingletonList." + sig.goName))
                      .append(")\n");
            }
            out.append("}\n");
            out.append("func (l jSingletonPtr[T]) ")
                .append(sig.goName)
                .append(formatExtSigWithParamNames(sig, imports))
                .append(" {\n");
            switch (sig.simpleName) {
              case "size" -> out.append("\treturn 1\n");
              case "isEmpty" -> out.append("\treturn false\n");
              case "get" -> out.append("\tCheckElementIndex(arg0, 1)\n\treturn l.p\n");
              case "getFirst", "getLast" -> out.append("\treturn l.p\n");
              case "contains" -> out.append("\treturn jSetElemEquals(l.p, arg0)\n");
              case "indexOf", "lastIndexOf" ->
                  out.append("\tif jSetElemEquals(l.p, arg0) { return 0 }\n\treturn -1\n");
              case "containsAll" ->
                  out.append(
                      "\treturn colContainsAll(func(x any) bool { return jSetElemEquals(l.p, x) }, arg0)\n");
              case "iterator" -> out.append("\treturn newSingleObjIter(l.p)\n");
              case "listIterator" -> {
                if (arity == 0) {
                  out.append("\treturn newImmListIter([]any{l.p}, 0)\n");
                } else {
                  out.append("\treturn newImmListIter([]any{l.p}, arg0)\n");
                }
              }
              case "subList" -> {
                out.append("\tCheckPositionIndexes(arg0, arg1, 1)\n\tif arg0 == arg1 { return ")
                    .append(emptyListExpr)
                    .append(" }\n\treturn l\n");
              }
              case "reversed" -> out.append("\treturn l\n");
              case "toArray" -> {
                if (arity == 0) {
                  out.append("\treturn []any{l.p}\n");
                } else if (p0.startsWith("[]")) {
                  out.append("\treturn sliceToArray1([]any{l.p}, arg0)\n");
                } else {
                  out.append("\tCheckNotNull(arg0)\n\treturn sliceToArray1([]any{l.p}, Cast[[]any](arg0.")
                      .append(intFnApply)
                      .append("(1)))\n");
                }
              }
              case "forEach" ->
                  out.append("\tCheckNotNull(arg0)\n\targ0.")
                      .append(consumerAccept)
                      .append("(l.p)\n");
              case "stream", "parallelStream" ->
                  out.append("\treturn streamFromSlice([]any{l.p})\n");
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
                  "replaceAll" ->
                  out.append("\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n");
              default ->
                  out.append("\tpanic(")
                      .append(GoNames.goStringLiteral("UNIMPLEMENTED_EXT: jSingletonPtr." + sig.goName))
                      .append(")\n");
            }
            out.append("}\n");
          }
        }
      }
      out.append('\n');
    }

    TypeElement setTe = resolver.elements.getTypeElement("java.util.Set");
    if (setTe != null && resolver.getReferencedExternalTypes().contains(setTe)) {
      String setGo = goTypeInJ2gort(setTe.asType(), imports);
      out.append("type jSingletonSet struct {\n\telem [1]any\n}\n");
      out.append("func newSingletonSet(v any) ")
          .append(setGo)
          .append(" {\n\treturn &jSingletonSet{elem: [1]any{ToAny(v)}}\n}\n");
      out.append("func (s *jSingletonSet) IsJavaNil_() bool { return s == nil }\n");
      out.append(
          "func (*jSingletonSet) JavaClassInfo_() (string, string) { return \"java.util.Collections$SingletonSet\", \"SingletonSet\" }\n");
      out.append(
          "func (s *jSingletonSet) IterSlice_() []any {\n\tif s == nil { ThrowNPE() }\n\treturn s.elem[:1:1]\n}\n");
      out.append(
          "func (s *jSingletonSet) SetContainsStr_(k string) bool {\n\tif s == nil { ThrowNPE() }\n\tif k == NullString { return s.elem[0] == nil }\n\tes, ok := s.elem[0].(string)\n\treturn ok && es == k\n}\n");
      out.append(
          "func (s *jSingletonSet) SetAddStr_(k string) bool {\n\tif s == nil { ThrowNPE() }\n\t_ = k\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      out.append(
          "func (s *jSingletonSet) SetRemoveStr_(k string) bool {\n\tif s == nil { ThrowNPE() }\n\t_ = k\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      out.append(
          "func (s *jSingletonSet) SetContainsNode_(n uint32) bool {\n\tif s == nil { ThrowNPE() }\n\tif n == 0 { return s.elem[0] == nil }\n\treturn singletonElemMatchesNode(s.elem[0], n)\n}\n");
      out.append(
          "func (s *jSingletonSet) SetAddNode_(n uint32) bool {\n\tif s == nil { ThrowNPE() }\n\t_ = n\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      out.append(
          "func (s *jSingletonSet) SetRemoveNode_(n uint32) bool {\n\tif s == nil { ThrowNPE() }\n\t_ = n\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      out.append(
          "func (s *jSingletonSet) Equals(arg0 any) bool { return colSetEqualsSingleton(s.elem[0], arg0) }\n");
      out.append(
          "func (s *jSingletonSet) HashCode() int32 { return javalang.ObjectHashCode(s.elem[0]) }\n");
      out.append("func (s *jSingletonSet) String() string { return \"[\" + StrVal(s.elem[0]) + \"]\" }\n");
      Set<String> setMarkers = new LinkedHashSet<>();
      setMarkers.add(resolver.getInterfaceMarkerMethod(setTe));
      for (TypeElement anc : resolver.getAllAncestorTypes(setTe)) {
        if (anc.getKind() == ElementKind.INTERFACE && !mapsToAny(anc)) {
          setMarkers.add(resolver.getInterfaceMarkerMethod(anc));
        }
      }
      setMarkers.add("Implements_j2gort_Ext_java_io_Serializable");
      for (String mk : setMarkers) {
        out.append("func (*jSingletonSet) ").append(mk).append("() {}\n");
      }
      Set<String> emittedSet = new HashSet<>(Set.of("Equals", "HashCode", "String"));
      List<TypeElement> setChain = new ArrayList<>();
      setChain.add(setTe);
      setChain.addAll(resolver.getAllAncestorTypes(setTe));
      String intFnApply = findExtMethodGoName("java.util.function.IntFunction", "apply", 1);
      String consumerAccept = findExtMethodGoName("java.util.function.Consumer", "accept", 1);
      for (TypeElement cur : setChain) {
        if (cur.getKind() != ElementKind.INTERFACE || mapsToAny(cur)) continue;
        for (Element m : cur.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && !ee.getModifiers().contains(Modifier.STATIC)
              && !ee.getModifiers().contains(Modifier.PRIVATE)) {
            TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
            if (sig == null || !resolver.isNeededExternalSig(sig) || !emittedSet.add(sig.goName)) {
              continue;
            }
            int arity = sig.canonParamTypes.size();
            String p0 = arity >= 1 ? goTypeInJ2gort(sig.canonParamTypes.get(0), imports) : "";
            out.append("func (s *jSingletonSet) ")
                .append(sig.goName)
                .append(formatExtSigWithParamNames(sig, imports))
                .append(" {\n");
            switch (sig.simpleName) {
              case "size" -> out.append("\treturn 1\n");
              case "isEmpty" -> out.append("\treturn false\n");
              case "contains" -> out.append("\treturn jSetElemEquals(s.elem[0], arg0)\n");
              case "containsAll" ->
                  out.append(
                      "\treturn colContainsAll(func(x any) bool { return jSetElemEquals(s.elem[0], x) }, arg0)\n");
              case "iterator" -> out.append("\treturn newSingleObjIter(s.elem[0])\n");
              case "toArray" -> {
                if (arity == 0) {
                  out.append("\treturn []any{s.elem[0]}\n");
                } else if (p0.startsWith("[]")) {
                  out.append("\treturn sliceToArray1(s.elem[:1:1], arg0)\n");
                } else {
                  out.append("\tCheckNotNull(arg0)\n\treturn sliceToArray1(s.elem[:1:1], Cast[[]any](arg0.")
                      .append(intFnApply)
                      .append("(1)))\n");
                }
              }
              case "forEach" ->
                  out.append("\tCheckNotNull(arg0)\n\targ0.")
                      .append(consumerAccept)
                      .append("(s.elem[0])\n");
              case "stream", "parallelStream" ->
                  out.append("\treturn streamFromSlice(s.elem[:1:1])\n");
              case "add",
                  "remove",
                  "clear",
                  "addAll",
                  "removeAll",
                  "retainAll",
                  "removeIf" ->
                  out.append("\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n");
              default ->
                  out.append("\tpanic(")
                      .append(GoNames.goStringLiteral("UNIMPLEMENTED_EXT: jSingletonSet." + sig.goName))
                      .append(")\n");
            }
            out.append("}\n");
          }
        }
      }
      out.append('\n');
    }

    TypeElement mapTeForRow = resolver.elements.getTypeElement("java.util.Map");
    if (mapTeForRow != null && resolver.getReferencedExternalTypes().contains(mapTeForRow)) {
      Set<String> mapMarkers = new LinkedHashSet<>();
      mapMarkers.add(resolver.getInterfaceMarkerMethod(mapTeForRow));
      for (TypeElement anc : resolver.getAllAncestorTypes(mapTeForRow)) {
        if (anc.getKind() == ElementKind.INTERFACE && !mapsToAny(anc)) {
          mapMarkers.add(resolver.getInterfaceMarkerMethod(anc));
        }
      }
      TypeElement absMapTe = resolver.elements.getTypeElement("java.util.AbstractMap");
      if (absMapTe != null && resolver.getReferencedExternalTypes().contains(absMapTe)) {
        mapMarkers.add(resolver.getInterfaceMarkerMethod(absMapTe));
        if (isSubclassedExternal(absMapTe)) {
          String absSt = goStructType(absMapTe, imports);
          out.append("func (*jTableRowMap) Get").append(absSt).append("() *").append(absSt).append(" { return nil }\n");
        }
      }
      for (String mk : mapMarkers) {
        out.append("func (*jTableRowMap) ").append(mk).append("() {}\n");
      }
      Set<String> emittedRowMap = new HashSet<>(Set.of("Equals", "HashCode", "String"));
      out.append("func (m *jTableRowMap) Equals(arg0 any) bool { return m.equalsImpl(arg0) }\n");
      out.append("func (m *jTableRowMap) HashCode() int32 { return m.hashCodeImpl() }\n");
      out.append("func (m *jTableRowMap) String() string { return m.stringImpl() }\n");
      List<TypeElement> mapChain = new ArrayList<>();
      mapChain.add(mapTeForRow);
      mapChain.addAll(resolver.getAllAncestorTypes(mapTeForRow));
      for (TypeElement cur : mapChain) {
        if (cur.getKind() != ElementKind.INTERFACE || mapsToAny(cur)) continue;
        for (Element m : cur.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && !ee.getModifiers().contains(Modifier.STATIC)
              && !ee.getModifiers().contains(Modifier.PRIVATE)) {
            TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
            if (sig == null || !resolver.isNeededExternalSig(sig) || !emittedRowMap.add(sig.goName)) {
              continue;
            }
            int arity = sig.canonParamTypes.size();
            out.append("func (m *jTableRowMap) ")
                .append(sig.goName)
                .append(formatExtSigWithParamNames(sig, imports))
                .append(" {\n");
            switch (sig.simpleName) {
              case "size" -> out.append("\treturn m.sizeImpl()\n");
              case "isEmpty" -> out.append("\treturn m.isEmptyImpl()\n");
              case "containsKey" -> out.append("\treturn m.containsKeyImpl(arg0)\n");
              case "containsValue" -> out.append("\treturn m.containsValueImpl(arg0)\n");
              case "get" -> out.append("\treturn m.getImpl(arg0)\n");
              case "getOrDefault" -> out.append("\treturn m.getOrDefaultImpl(arg0, arg1)\n");
              case "put" -> out.append("\treturn m.putImpl(arg0, arg1)\n");
              case "putIfAbsent" -> out.append("\treturn abstractMapOp_putIfAbsent_2(m, arg0, arg1)\n");
              case "remove" -> {
                if (arity == 1) {
                  out.append("\treturn m.removeImpl(arg0)\n");
                } else {
                  out.append(
                      "\tif !m.containsKeyImpl(arg0) || !javalang.ObjectsEquals(ToAny(m.getImpl(arg0)), ToAny(arg1)) { return false }\n\tm.removeImpl(arg0)\n\treturn true\n");
                }
              }
              case "putAll" -> out.append("\tabstractMapOp_putAll_1(m, arg0)\n");
              case "clear" -> out.append("\tm.clearImpl()\n");
              case "keySet" -> out.append("\treturn m.keySetImpl()\n");
              case "values" -> out.append("\treturn jTableRowValues{m: m}\n");
              case "entrySet" -> out.append("\treturn m.entrySetImpl()\n");
              case "forEach" -> out.append("\tm.forEachImpl(arg0)\n");
              case "replaceAll" -> out.append("\tabstractMapOp_replaceAll_1(m, arg0)\n");
              case "replace" -> {
                if (arity == 2) {
                  out.append(
                      "\tif !m.containsKeyImpl(arg0) { return nil }\n\treturn m.putImpl(arg0, arg1)\n");
                } else {
                  out.append("\treturn abstractMapOp_replace_3(m, arg0, arg1, arg2)\n");
                }
              }
              case "computeIfAbsent" -> out.append("\treturn abstractMapOp_computeIfAbsent_2(m, arg0, arg1)\n");
              case "compute" -> out.append("\treturn abstractMapOp_compute_2(m, arg0, arg1)\n");
              case "merge" -> out.append("\treturn abstractMapOp_merge_3(m, arg0, arg1, arg2)\n");
              default ->
                  out.append("\tpanic(")
                      .append(GoNames.goStringLiteral("UNIMPLEMENTED_EXT: jTableRowMap." + sig.goName))
                      .append(")\n");
            }
            out.append("}\n");
          }
        }
      }
      out.append('\n');
    }

    TypeElement colTeForRow = resolver.elements.getTypeElement("java.util.Collection");
    if (colTeForRow != null && resolver.getReferencedExternalTypes().contains(colTeForRow)) {
      Set<String> colMarkers = new LinkedHashSet<>();
      colMarkers.add(resolver.getInterfaceMarkerMethod(colTeForRow));
      for (TypeElement anc : resolver.getAllAncestorTypes(colTeForRow)) {
        if (anc.getKind() == ElementKind.INTERFACE && !mapsToAny(anc)) {
          colMarkers.add(resolver.getInterfaceMarkerMethod(anc));
        }
      }
      for (String mk : colMarkers) {
        out.append("func (jTableRowValues) ").append(mk).append("() {}\n");
      }
      String intFnApply = findExtMethodGoName("java.util.function.IntFunction", "apply", 1);
      Set<String> emittedRowVals = new HashSet<>(Set.of("Equals", "HashCode", "String"));
      out.append("func (v jTableRowValues) Equals(arg0 any) bool { return RefEquals(v, arg0) }\n");
      out.append("func (v jTableRowValues) HashCode() int32 { return javalang.IdentityHashCode(v) }\n");
      out.append("func (v jTableRowValues) String() string { return sliceJavaString(v.toArrayImpl()) }\n");
      List<TypeElement> colChain = new ArrayList<>();
      colChain.add(colTeForRow);
      colChain.addAll(resolver.getAllAncestorTypes(colTeForRow));
      for (TypeElement cur : colChain) {
        if (cur.getKind() != ElementKind.INTERFACE || mapsToAny(cur)) continue;
        for (Element m : cur.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && !ee.getModifiers().contains(Modifier.STATIC)
              && !ee.getModifiers().contains(Modifier.PRIVATE)) {
            TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
            if (sig == null || !resolver.isNeededExternalSig(sig) || !emittedRowVals.add(sig.goName)) {
              continue;
            }
            int arity = sig.canonParamTypes.size();
            String p0 = arity >= 1 ? goTypeInJ2gort(sig.canonParamTypes.get(0), imports) : "";
            out.append("func (v jTableRowValues) ")
                .append(sig.goName)
                .append(formatExtSigWithParamNames(sig, imports))
                .append(" {\n");
            switch (sig.simpleName) {
              case "size" -> out.append("\treturn v.m.sizeImpl()\n");
              case "isEmpty" -> out.append("\treturn v.m.isEmptyImpl()\n");
              case "contains" -> out.append("\treturn v.m.containsValueImpl(arg0)\n");
              case "containsAll" -> out.append("\treturn colContainsAll(v.m.containsValueImpl, arg0)\n");
              case "iterator" -> out.append("\treturn v.iteratorImpl()\n");
              case "toArray" -> {
                if (arity == 0) {
                  out.append("\treturn v.toArrayImpl()\n");
                } else if (p0.startsWith("[]")) {
                  out.append("\treturn sliceToArray1(v.toArrayImpl(), arg0)\n");
                } else {
                  out.append(
                      "\tCheckNotNull(arg0)\n\tarr := v.toArrayImpl()\n\treturn sliceToArray1(arr, Cast[[]any](arg0.")
                      .append(intFnApply)
                      .append("(int32(len(arr)))))\n");
                }
              }
              case "forEach" -> out.append("\tv.forEachImpl(arg0)\n");
              case "stream", "parallelStream" -> out.append("\treturn streamFromIterable(v)\n");
              case "clear" -> out.append("\tv.m.clearImpl()\n");
              case "remove" -> out.append("\treturn abstractColRemove(v, arg0)\n");
              case "removeAll" -> out.append("\treturn abstractColRemoveAll(v, arg0)\n");
              case "retainAll" -> out.append("\treturn abstractColRetainAll(v, arg0)\n");
              case "removeIf" -> out.append("\treturn abstractColRemoveIf(v, arg0)\n");
              case "add", "addAll" ->
                  out.append("\tpanic(javalang.NewUnsupportedOperationException(\"Map.values.add\"))\n");
              default ->
                  out.append("\tpanic(")
                      .append(GoNames.goStringLiteral("UNIMPLEMENTED_EXT: jTableRowValues." + sig.goName))
                      .append(")\n");
            }
            out.append("}\n");
          }
        }
      }
      out.append('\n');
    }

    if (hasImmList) {
      emitImmListVariantsAndBuilder(ilTe, imports, out);
    }

    TypeElement ismmTe = resolver.elements.getTypeElement("com.google.common.collect.ImmutableSetMultimap");
    if (ismmTe != null && resolver.getReferencedExternalTypes().contains(ismmTe)) {
      String st = goStructType(ismmTe, imports);
      out.append("var emptyImmSetMultimap_ *").append(st).append("\n");
      out.append("func init() {\n\tmm := newMultimapCore(true, true, nil, nil)\n\tobj := &")
          .append(st)
          .append("{}\n\tobj.BindExtSelf_(obj)\n\tobj.Core_ = mm\n\temptyImmSetMultimap_ = obj\n}\n\n");
    }

    TypeElement isTe = resolver.elements.getTypeElement("com.google.common.collect.ImmutableSet");
    if (isTe != null && resolver.getReferencedExternalTypes().contains(isTe)) {
      String st = goStructType(isTe, imports);
      String isGo = goTypeInJ2gort(isTe.asType(), imports);
      out.append("var emptyImmutableSet_ ").append(isGo).append("\n");
      out.append("func init() {\n\tobj := &").append(st).append("{elems_: emptyAnySlice_}\n\tobj.BindExtSelf_(obj)\n\temptyImmutableSet_ = obj\n}\n\n");
      out.append("func newImmutableSetShell() *").append(st).append(" {\n");
      out.append("\tobj := &").append(st).append("{}\n");
      out.append("\tobj.BindExtSelf_(obj)\n");
      out.append("\treturn obj\n}\n\n");
      out.append("func ImmSetOfStr1(s0 string) ").append(isGo).append(" {\n");
      out.append("\tif s0 == NullString { ThrowNPE() }\n");
      out.append("\tobj := &").append(st).append("{}\n\tobj.BindExtSelf_(obj)\n");
      out.append("\tobj.singleStr_[0] = s0\n\tobj.strs_ = obj.singleStr_[:1:1]\n\treturn obj\n}\n\n");
      out.append("func ImmSetOfStr2(s0, s1 string) ").append(isGo).append(" {\n");
      out.append("\tif s0 == NullString || s1 == NullString { ThrowNPE() }\n");
      out.append("\tif s0 == s1 { return ImmSetOfStr1(s0) }\n");
      out.append("\tobj := &").append(st).append("{strs_: []string{s0, s1}}\n\tobj.BindExtSelf_(obj)\n\treturn obj\n}\n\n");
      out.append("func finishImmStrSetOwn(strs []string, _ *JSetCore) *").append(st).append(" {\n");
      out.append("\tobj := &").append(st).append("{strs_: strs[:len(strs):len(strs)]}\n\tobj.BindExtSelf_(obj)\n");
      out.append("\treturn obj\n}\n\n");
      out.append("func finishImmSetSingleRaw(e0 any) *").append(st).append(" {\n");
      out.append("\tobj := &").append(st).append("{}\n");
      out.append("\tobj.BindExtSelf_(obj)\n");
      out.append("\tobj.singleArr_[0] = e0\n");
      out.append("\tobj.elems_ = obj.singleArr_[:1:1]\n");
      out.append("\ttryInitImmEnumSingle(obj, e0)\n");
      out.append("\treturn obj\n}\n\n");
      out.append("func newImmSet1(e0 any) ").append(isGo).append(" {\n");
      out.append("\tif s0, ok := e0.(string); ok {\n\t\treturn ImmSetOfStr1(s0)\n\t}\n");
      out.append("\te0 = ToAny(e0)\n\tif e0 == nil { ThrowNPE() }\n");
      out.append("\treturn finishImmSetSingleRaw(e0)\n}\n\n");
      out.append("func newImmSet2(e0, e1 any) ").append(isGo).append(" {\n");
      out.append("\tif s0, ok0 := e0.(string); ok0 {\n\t\tif s1, ok1 := e1.(string); ok1 {\n\t\t\treturn ImmSetOfStr2(s0, s1)\n\t\t}\n\t}\n");
      out.append("\te0 = ToAny(e0)\n\tif e0 == nil { ThrowNPE() }\n");
      out.append("\te1 = ToAny(e1)\n\tif e1 == nil { ThrowNPE() }\n");
      out.append("\tif elemJavaEquals(e0, e1, false) {\n\t\treturn finishImmSetSingleRaw(e0)\n\t}\n");
      out.append("\treturn finishImmSetMulti([]any{e0, e1}, nil)\n}\n\n");
      out.append("func wrapImmutableSetFromSetCoreOwn(sc *JSetCore) ").append(isGo).append(" {\n");
      out.append("\tif sc == nil || sc.Size() == 0 {\n\t\treturn emptyImmutableSet_\n\t}\n");
      out.append("\tif sc.mode == jSetModeStr {\n\t\tsc.sealForImmutable()\n\t\tst := sc.ext.(*jSetStrState)\n\t\thasNull := false\n\t\tfor _, s := range st.strs { if s == NullString { hasNull = true; break } }\n\t\tif !hasNull {\n\t\t\tif len(st.strs) == 1 { return ImmSetOfStr1(st.strs[0]) }\n\t\t\treturn finishImmStrSetOwn(st.strs, sc)\n\t\t}\n\t}\n");
      out.append("\tif sc.Size() == 1 {\n\t\treturn finishImmSetSingleRaw(sc.FirstKey())\n\t}\n");
      out.append("\treturn finishImmSetMulti(sc.LiveKeys(), sc)\n}\n\n");
      out.append("func wrapImmutableSetCore(mc *JMapCore) ").append(isGo).append(" {\n");
      out.append("\tif mc == nil || mc.Size() == 0 {\n\t\treturn emptyImmutableSet_\n\t}\n");
      out.append("\tif mc.mode == jMapModeStr {\n\t\thasNull := false\n\t\tfor i := int(mc.firstIdx); i < len(mc.strEntries); i++ {\n\t\t\tk := mc.strEntries[i].key\n\t\t\tif k == NullString { hasNull = true; break }\n\t\t}\n\t\tif !hasNull {\n\t\t\tif mc.Size() == 1 {\n\t\t\t\tfor i := int(mc.firstIdx); i < len(mc.strEntries); i++ {\n\t\t\t\t\tif k := mc.strEntries[i].key; k != jTombstoneStr { return ImmSetOfStr1(k) }\n\t\t\t\t}\n\t\t\t}\n\t\t\tstrs := make([]string, 0, mc.Size())\n\t\t\tfor i := int(mc.firstIdx); i < len(mc.strEntries); i++ {\n\t\t\t\tif k := mc.strEntries[i].key; k != jTombstoneStr { strs = append(strs, k) }\n\t\t\t}\n\t\t\treturn finishImmStrSetOwn(strs, nil)\n\t\t}\n\t}\n");
      out.append("\tif mc.Size() == 1 {\n\t\treturn finishImmSetSingleRaw(mc.FirstKey())\n\t}\n");
      out.append("\treturn finishImmSetMulti(mc.LiveKeys(), nil)\n}\n\n");
      out.append("func wrapImmutableSet(items []any) ").append(isGo).append(" {\n");
      out.append("\tswitch len(items) {\n\tcase 0:\n\t\treturn emptyImmutableSet_\n");
      out.append("\tcase 1:\n\t\tif s0, ok := items[0].(string); ok && s0 != NullString { return ImmSetOfStr1(s0) }\n\t\treturn finishImmSetSingleRaw(ToAny(items[0]))\n\t}\n");
      out.append("\tif len(items) <= 8 {\n\t\tvar buf [8]any\n\t\tn := 0\n\t\tallStr := true\n\t\tfor _, raw := range items {\n\t\t\tv := ToAny(raw)\n\t\t\tif _, ok := v.(string); !ok { allStr = false }\n\t\t\tdup := false\n\t\t\tfor j := 0; j < n; j++ {\n\t\t\t\tif elemJavaEquals(buf[j], v, false) { dup = true; break }\n\t\t\t}\n\t\t\tif !dup { buf[n] = v; n++ }\n\t\t}\n");
      out.append("\t\tif allStr {\n\t\t\tif n == 1 { return ImmSetOfStr1(buf[0].(string)) }\n\t\t\tstrs := make([]string, n)\n\t\t\tfor i := 0; i < n; i++ { strs[i] = buf[i].(string) }\n\t\t\treturn finishImmStrSetOwn(strs, nil)\n\t\t}\n");
      out.append("\t\tif n == 1 {\n\t\t\treturn finishImmSetSingleRaw(buf[0])\n\t\t}\n");
      out.append("\t\telems := make([]any, n)\n\t\tcopy(elems, buf[:n])\n\t\treturn finishImmSetMulti(elems, nil)\n\t}\n");
      out.append("\tsc := &JSetCore{}\n\tsc.InitExpectedSize(int32(len(items)))\n\tfor _, x := range items { sc.Add(x) }\n\treturn wrapImmutableSetFromSetCoreOwn(sc)\n}\n\n");
      TypeElement isbTe = resolver.elements.getTypeElement("com.google.common.collect.ImmutableSet.Builder");
      TypeElement issbTe = resolver.elements.getTypeElement("com.google.common.collect.ImmutableSortedSet.Builder");
      if (isbTe != null && resolver.getReferencedExternalTypes().contains(isbTe)) {
        String isbGo = goTypeInJ2gort(isbTe.asType(), imports);
        String isbSt = goStructType(isbTe, imports);
        String isbAddGo = "Add";
        for (Element m : isbTe.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && !ee.isVarArgs()
              && ee.getSimpleName().contentEquals("add")
              && ee.getParameters().size() == 1) {
            TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
            if (sig != null) {
              isbAddGo = sig.goName;
              break;
            }
          }
        }
        boolean hasIssb = issbTe != null && resolver.getReferencedExternalTypes().contains(issbTe);
        String issbSt = hasIssb ? goStructType(issbTe, imports) : "";
        out.append("func ImmSetBuilderAddStrVoid(b ").append(isbGo).append(", s string) {\n")
            .append("\tif b == nil || s == NullString { ThrowNPE() }\n")
            .append("\tif sb, ok := b.(*").append(isbSt).append("); ok {\n\t\tif sb == nil { ThrowNPE() }\n\t\tsb.inlineSet.AddStr(s)\n\t\treturn\n\t}\n");
        if (hasIssb) {
          out.append("\tif ssb, ok := b.(*").append(issbSt).append("); ok {\n\t\tif ssb == nil { ThrowNPE() }\n\t\tssb.sortedCore_.PutStr(s, true)\n\t\treturn\n\t}\n");
        }
        out.append("\tif isJavaNilAny(b) { ThrowNPE() }\n")
            .append("\tif ab, ok := any(b).(interface{ ").append(isbAddGo).append("(any) ").append(isbGo).append(" }); ok { ab.").append(isbAddGo).append("(BoxString(s)); return }\n")
            .append("\tif ab, ok := any(b).(interface{ ").append(isbAddGo).append("(any) any }); ok { ab.").append(isbAddGo).append("(BoxString(s)) }\n}\n")
            .append("func ImmSetBuilderAddStr(b ").append(isbGo).append(", s string) ").append(isbGo).append(" {\n")
            .append("\tImmSetBuilderAddStrVoid(b, s)\n\treturn b\n}\n\n");
      }
    }

    TypeElement imTe = resolver.elements.getTypeElement("com.google.common.collect.ImmutableMap");
    if (imTe != null && resolver.getReferencedExternalTypes().contains(imTe)) {
      String st = goStructType(imTe, imports);
      String imGo = goTypeInJ2gort(imTe.asType(), imports);
      out.append("var emptyImmutableMap_ ").append(imGo).append("\n");
      out.append("func init() {\n\tobj := &").append(st).append("{}\n\tobj.BindExtSelf_(obj)\n\tobj.inlineCore.sealImmutable()\n\temptyImmutableMap_ = obj\n}\n\n");
      out.append("func newImmutableMapShell() *").append(st).append(" {\n");
      out.append("\tobj := &").append(st).append("{}\n");
      out.append("\tobj.BindExtSelf_(obj)\n");
      out.append("\treturn obj\n}\n\n");
      out.append("func wrapImmutableMap(src *JMapCore) ").append(imGo).append(" {\n");
      out.append("\tif src == nil || src.Size() == 0 {\n\t\treturn emptyImmutableMap_\n\t}\n");
      out.append("\tobj := newImmutableMapShell()\n");
      out.append(
          "\tmc := &obj.inlineCore\n\tif !mc.CopyFromSameMode(src) {\n\t\tmc.InitExpectedSize(src.Size())\n\t\tsrc.ForEachEntry(func(k, v any) bool { mc.Put(k, v); return true })\n\t}\n\tmc.sealImmutable()\n\treturn obj\n}\n");
      if (isTe != null && hasImmList) {
        String isSt = goStructType(isTe, imports);
        String isGo = goTypeInJ2gort(isTe.asType(), imports);
        String ilSt = goStructType(ilTe, imports);
        String ilGo = goTypeInJ2gort(ilTe.asType(), imports);
        out.append("func (b *").append(st).append(") cachedKeySet_() ").append(isGo).append(" {\n")
            .append("\tmc := b.mc()\n")
            .append("\tif mc.IsEmpty() {\n\t\treturn emptyImmutableSet_\n\t}\n")
            .append("\tif s := b.keySetCache_.Load(); s != nil {\n\t\treturn s\n\t}\n")
            .append("\ts := any(wrapImmutableSetCore(mc)).(*").append(isSt).append(")\n")
            .append("\tif b.keySetCache_.CompareAndSwap(nil, s) {\n\t\treturn s\n\t}\n")
            .append("\treturn b.keySetCache_.Load()\n}\n")
            .append("func (b *").append(st).append(") cachedValues_() ").append(ilGo).append(" {\n")
            .append("\tmc := b.mc()\n")
            .append("\tif mc.IsEmpty() {\n\t\treturn emptyImmutableList_\n\t}\n")
            .append("\tif v := b.valuesCache_.Load(); v != nil {\n\t\treturn v\n\t}\n")
            .append("\tv := newImmListBaseOwn(mc.LiveValues())\n")
            .append("\tif b.valuesCache_.CompareAndSwap(nil, v) {\n\t\treturn v\n\t}\n")
            .append("\treturn b.valuesCache_.Load()\n}\n")
            .append("func (b *").append(st).append(") cachedEntrySet_() ").append(isGo).append(" {\n")
            .append("\tmc := b.mc()\n")
            .append("\tif mc.IsEmpty() {\n\t\treturn emptyImmutableSet_\n\t}\n")
            .append("\tif s := b.entrySetCache_.Load(); s != nil {\n\t\treturn s\n\t}\n")
            .append("\ts := any(wrapImmutableMapEntrySet(mc)).(*").append(isSt).append(")\n")
            .append("\tif b.entrySetCache_.CompareAndSwap(nil, s) {\n\t\treturn s\n\t}\n")
            .append("\treturn b.entrySetCache_.Load()\n}\n\n");
      } else {
        out.append('\n');
      }
      TypeElement ismTe = resolver.elements.getTypeElement("com.google.common.collect.ImmutableSortedMap");
      TypeElement issTe = resolver.elements.getTypeElement("com.google.common.collect.ImmutableSortedSet");
      if (ismTe != null && resolver.getReferencedExternalTypes().contains(ismTe)) {
        String ismSt = goStructType(ismTe, imports);
        if (issTe != null && resolver.getReferencedExternalTypes().contains(issTe)) {
          String issSt = goStructType(issTe, imports);
          String issGo = goTypeInJ2gort(issTe.asType(), imports);
          out.append("func (b *").append(ismSt).append(") cachedSortedKeySet_() ").append(issGo).append(" {\n")
              .append("\tmc := b.mc()\n")
              .append("\tif mc.IsEmpty() {\n\t\treturn emptyImmutableSortedSet_\n\t}\n")
              .append("\tif s := b.sortedKeySetCache_.Load(); s != nil {\n\t\treturn s\n\t}\n")
              .append("\ts := any(wrapImmutableSortedSet(mc)).(*").append(issSt).append(")\n")
              .append("\tif b.sortedKeySetCache_.CompareAndSwap(nil, s) {\n\t\treturn s\n\t}\n")
              .append("\treturn b.sortedKeySetCache_.Load()\n}\n\n");
        } else if (isTe != null && hasImmList) {
          String isGo = goTypeInJ2gort(isTe.asType(), imports);
          out.append("func (b *").append(ismSt).append(") cachedSortedKeySet_() ").append(isGo).append(" {\n")
              .append("\treturn b.cachedKeySet_()\n}\n\n");
        }
      }
    }

    TypeElement esTe = resolver.elements.getTypeElement("java.util.EnumSet");
    if (esTe != null && resolver.getReferencedExternalTypes().contains(esTe)) {
      String st = goStructType(esTe, imports);
      String esGo = goTypeInJ2gort(esTe.asType(), imports);
      out.append("func newEnumSetShell() *").append(st).append(" {\n");
      out.append("\tobj := &").append(st).append("{}\n");
      out.append("\tobj.BindExtSelf_(obj)\n");
      out.append("\treturn obj\n}\n\n");
      out.append("func wrapEnumSet(items []any) ").append(esGo).append(" {\n");
      out.append("\tobj := newEnumSetShell()\n");
      out.append("\tfor _, x := range items { obj.inlineEnumSet.Add(x) }\n\treturn obj\n}\n\n");
    }

    TypeElement intTe = resolver.elements.getTypeElement("com.google.common.collect.Interner");
    if (intTe != null && resolver.getReferencedExternalTypes().contains(intTe)) {
      String internerGo = goTypeInJ2gort(intTe.asType(), imports);
      out.append(
          "const (\n"
              + "\tjWeakInternShards = 64\n"
              + "\tjWeakInternSets = 16\n"
              + "\tjWeakInternWays = 4\n"
              + "\tjWeakInternMaxStrLen = 128\n"
              + "\tjInternShardThreshold = 256\n"
              + ")\n");
      out.append(
          "type jWeakInternShard struct {\n"
              + "\tmu sync.RWMutex\n"
              + "\tsets [jWeakInternSets][jWeakInternWays]string\n"
              + "}\n");
      out.append("type jInternStrShard struct {\n\tmu sync.RWMutex\n\tm map[string]string\n}\n");
      out.append(
          "type simpleInterner struct {\n"
              + "\tweak bool\n"
              + "\tstrMu sync.RWMutex\n"
              + "\tstrs map[string]string\n"
              + "\tstrShards atomic.Pointer[[64]jInternStrShard]\n"
              + "\tweakShards atomic.Pointer[[jWeakInternShards]jWeakInternShard]\n"
              + "\tmu sync.Mutex\n"
              + "\tmc *JMapCore\n"
              + "}\n");
      out.append("func (s *simpleInterner) IsJavaNil_() bool { return s == nil }\n");
      out.append("func newSimpleInterner() ")
          .append(internerGo)
          .append(" {\n\treturn &simpleInterner{}\n}\n");
      out.append("func newWeakSimpleInterner() ")
          .append(internerGo)
          .append(" {\n\treturn &simpleInterner{weak: true}\n}\n");
      out.append(
          "func (s *simpleInterner) Implements_j2gort_Ext_com_google_common_collect_Interner() {}\n");
      out.append("func internStrHash(s string) uint32 {\n")
          .append("\tn := len(s)\n\th := uint32(2166136261) ^ uint32(n)\n")
          .append("\tif n <= 32 {\n\t\tfor i := 0; i < n; i++ { h = (h ^ uint32(s[i])) * 16777619 }\n\t} else {\n")
          .append("\t\tstep := n >> 4\n\t\tfor i := 0; i < n; i += step { h = (h ^ uint32(s[i])) * 16777619 }\n\t\th = (h ^ uint32(s[n-1])) * 16777619\n\t}\n")
          .append("\treturn h ^ (h >> 13)\n}\n");
      out.append(
          "func (s *simpleInterner) initWeakShards() *[jWeakInternShards]jWeakInternShard {\n"
              + "\tws := new([jWeakInternShards]jWeakInternShard)\n"
              + "\tif s.weakShards.CompareAndSwap(nil, ws) {\n"
              + "\t\treturn ws\n"
              + "\t}\n"
              + "\treturn s.weakShards.Load()\n"
              + "}\n");
      out.append(
          "func internIntoShard(sh *jInternStrShard, str string) string {\n"
              + "\tsh.mu.RLock()\n\tv, ok := sh.m[str]\n\tsh.mu.RUnlock()\n\tif ok { return v }\n"
              + "\tsh.mu.Lock()\n\tif sh.m == nil {\n\t\tsh.m = make(map[string]string, 32)\n\t} else if v, ok = sh.m[str]; ok {\n\t\tsh.mu.Unlock()\n\t\treturn v\n\t}\n"
              + "\tsh.m[str] = str\n\tsh.mu.Unlock()\n\treturn str\n"
              + "}\n");
      out.append(
          "func (s *simpleInterner) InternStr(str string) string {\n"
              + "\tif s == nil || str == NullString { ThrowNPE() }\n"
              + "\tif len(str) == 0 { return \"\" }\n"
              + "\tif len(str) == 1 && str[0] < 128 { return asciiInternStrings_[str[0]] }\n"
              + "\tif s.weak {\n"
              + "\t\tif len(str) > jWeakInternMaxStrLen { return str }\n"
              + "\t\th := internStrHash(str)\n"
              + "\t\tws := s.weakShards.Load()\n"
              + "\t\tif ws == nil { ws = s.initWeakShards() }\n"
              + "\t\tsh := &ws[h&(jWeakInternShards-1)]\n"
              + "\t\tsetIdx := (h >> 6) & (jWeakInternSets - 1)\n"
              + "\t\tsh.mu.RLock()\n"
              + "\t\trow := &sh.sets[setIdx]\n"
              + "\t\tfor i := 0; i < jWeakInternWays; i++ {\n"
              + "\t\t\tif v := row[i]; len(v) == len(str) && v == str {\n"
              + "\t\t\t\tsh.mu.RUnlock()\n"
              + "\t\t\t\treturn v\n"
              + "\t\t\t}\n"
              + "\t\t}\n"
              + "\t\tsh.mu.RUnlock()\n"
              + "\t\tsh.mu.Lock()\n"
              + "\t\tfor i := 0; i < jWeakInternWays; i++ {\n"
              + "\t\t\tif v := row[i]; len(v) == len(str) && v == str {\n"
              + "\t\t\t\tsh.mu.Unlock()\n"
              + "\t\t\t\treturn v\n"
              + "\t\t\t}\n"
              + "\t\t}\n"
              + "\t\trow[3] = row[2]\n"
              + "\t\trow[2] = row[1]\n"
              + "\t\trow[1] = row[0]\n"
              + "\t\trow[0] = str\n"
              + "\t\tsh.mu.Unlock()\n"
              + "\t\treturn str\n"
              + "\t}\n"
              + "\tif shs := s.strShards.Load(); shs != nil {\n"
              + "\t\treturn internIntoShard(&shs[internStrHash(str)&63], str)\n"
              + "\t}\n"
              + "\ts.strMu.RLock()\n"
              + "\tif s.strs != nil {\n"
              + "\t\tif v, ok := s.strs[str]; ok {\n"
              + "\t\t\ts.strMu.RUnlock()\n"
              + "\t\t\treturn v\n"
              + "\t\t}\n"
              + "\t}\n"
              + "\ts.strMu.RUnlock()\n"
              + "\ts.strMu.Lock()\n"
              + "\tif shs := s.strShards.Load(); shs != nil {\n"
              + "\t\ts.strMu.Unlock()\n"
              + "\t\treturn internIntoShard(&shs[internStrHash(str)&63], str)\n"
              + "\t}\n"
              + "\tif s.strs == nil {\n"
              + "\t\ts.strs = make(map[string]string, 16)\n"
              + "\t} else if v, ok := s.strs[str]; ok {\n"
              + "\t\ts.strMu.Unlock()\n"
              + "\t\treturn v\n"
              + "\t}\n"
              + "\tif len(s.strs) >= jInternShardThreshold {\n"
              + "\t\tshs := new([64]jInternStrShard)\n"
              + "\t\t// order-independent: partitioning unique interned strings into hash shards\n"
              + "\t\tfor k, v := range s.strs {\n"
              + "\t\t\tsh := &shs[internStrHash(k)&63]\n"
              + "\t\t\tif sh.m == nil { sh.m = make(map[string]string, 16) }\n"
              + "\t\t\tsh.m[k] = v\n"
              + "\t\t}\n"
              + "\t\tsh := &shs[internStrHash(str)&63]\n"
              + "\t\tif sh.m == nil { sh.m = make(map[string]string, 16) }\n"
              + "\t\tsh.m[str] = str\n"
              + "\t\ts.strShards.Store(shs)\n"
              + "\t\ts.strs = nil\n"
              + "\t\ts.strMu.Unlock()\n"
              + "\t\treturn str\n"
              + "\t}\n"
              + "\ts.strs[str] = str\n"
              + "\ts.strMu.Unlock()\n"
              + "\treturn str\n"
              + "}\n");
      for (Element m : intTe.getEnclosedElements()) {
        if (m instanceof ExecutableElement ee
            && ee.getKind() == ElementKind.METHOD
            && !ee.getModifiers().contains(Modifier.STATIC)) {
          TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
          if (sig == null || !resolver.isNeededExternalSig(sig)) continue;
          out.append("func (s *simpleInterner) ")
              .append(sig.goName)
              .append(formatExtSigWithParamNames(sig, imports))
              .append(" {\n");
          if (sig.simpleName.equals("intern") && sig.canonParamTypes.size() == 1) {
            out.append(
                "\tif sStr, ok := arg0.(string); ok {\n\t\treturn s.InternStr(sStr)\n\t}\n"
                    + "\tCheckNotNull(arg0)\n\ts.mu.Lock()\n\tdefer s.mu.Unlock()\n"
                    + "\tif s.mc == nil {\n\t\ts.mc = &JMapCore{}\n\t} else if ex := s.mc.Get(arg0); ex != nil {\n\t\treturn ex\n\t} else if s.weak && s.mc.Live >= 4096 {\n\t\t*s.mc = JMapCore{}\n\t}\n"
                    + "\ts.mc.Put(arg0, arg0)\n\treturn arg0\n}\n");
          } else {
            out.append("\tpanic(")
                .append(
                    GoNames.goStringLiteral("UNIMPLEMENTED_EXT: simpleInterner." + sig.goName))
                .append(")\n}\n");
          }
        }
      }
      out.append("func InternString(interner ")
          .append(internerGo)
          .append(", s string) string {\n")
          .append("\tif interner == nil { ThrowNPE() }\n")
          .append("\treturn interner.InternStr(s)\n}\n\n");
    }

    String mapGetGo = findExtMethodGoName("java.util.Map", "get", 1);
    String mapGetOrDefGo = findExtMethodGoName("java.util.Map", "getOrDefault", 2);
    String mapContainsKeyGo = findExtMethodGoName("java.util.Map", "containsKey", 1);
    String mapPutGo = findExtMethodGoName("java.util.Map", "put", 2);
    String mapRemoveGo = findExtMethodGoName("java.util.Map", "remove", 1);
    String mapComputeIfAbsentGo = findExtMethodGoName("java.util.Map", "computeIfAbsent", 2);
    String colContainsGo = findExtMethodGoName("java.util.Collection", "contains", 1);
    String colAddGo = findExtMethodGoName("java.util.Collection", "add", 1);
    String colRemoveGo = findExtMethodGoName("java.util.Collection", "remove", 1);

    out.append("func MapGetStr(m any, k string) any {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastStrMap); ok { return fm.MapGetStr_(k) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(Ext_java_util_Map).")
        .append(mapGetGo)
        .append("(BoxString(k))\n}\n")
        .append("func MapGetOrDefaultStr(m any, k string, def any) any {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastStrMap); ok { return fm.MapGetOrDefaultStr_(k, def) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(Ext_java_util_Map).")
        .append(mapGetOrDefGo)
        .append("(BoxString(k), ToAny(def))\n}\n")
        .append("func MapContainsKeyStr(m any, k string) bool {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastStrMap); ok { return fm.MapContainsStr_(k) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(Ext_java_util_Map).")
        .append(mapContainsKeyGo)
        .append("(BoxString(k))\n}\n")
        .append("func MapPutStr(m any, k string, v any) any {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastStrMap); ok { return fm.MapPutStr_(k, v) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(Ext_java_util_Map).")
        .append(mapPutGo)
        .append("(BoxString(k), ToAny(v))\n}\n")
        .append("func MapRemoveStr(m any, k string) any {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastStrMap); ok { return fm.MapRemoveStr_(k) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(Ext_java_util_Map).")
        .append(mapRemoveGo)
        .append("(BoxString(k))\n}\n")
        .append("func SetContainsStr(s any, k string) bool {\n")
        .append("\tif s == nil { ThrowNPE() }\n")
        .append("\tif fs, ok := s.(FastStrSet); ok { return fs.SetContainsStr_(k) }\n")
        .append("\tif isJavaNilAny(s) { ThrowNPE() }\n")
        .append("\treturn s.(Ext_java_util_Collection).")
        .append(colContainsGo)
        .append("(BoxString(k))\n}\n")
        .append("func SetAddStr(s any, k string) bool {\n")
        .append("\tif s == nil { ThrowNPE() }\n")
        .append("\tif fs, ok := s.(FastStrSet); ok { return fs.SetAddStr_(k) }\n")
        .append("\tif isJavaNilAny(s) { ThrowNPE() }\n")
        .append("\treturn s.(Ext_java_util_Collection).")
        .append(colAddGo)
        .append("(BoxString(k))\n}\n")
        .append("func SetRemoveStr(s any, k string) bool {\n")
        .append("\tif s == nil { ThrowNPE() }\n")
        .append("\tif fs, ok := s.(FastStrSet); ok { return fs.SetRemoveStr_(k) }\n")
        .append("\tif isJavaNilAny(s) { ThrowNPE() }\n")
        .append("\treturn s.(Ext_java_util_Collection).")
        .append(colRemoveGo)
        .append("(BoxString(k))\n}\n")
        .append("func MapGetNode[N ~uint32](m any, n N) any {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastNodeMap); ok { return fm.MapGetNode_(uint32(n)) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(Ext_java_util_Map).")
        .append(mapGetGo)
        .append("(BoxNode(n))\n}\n")
        .append("func MapGetOrDefaultNode[N ~uint32](m any, n N, def any) any {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastNodeMap); ok { return fm.MapGetOrDefaultNode_(uint32(n), def) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(Ext_java_util_Map).")
        .append(mapGetOrDefGo)
        .append("(BoxNode(n), ToAny(def))\n}\n")
        .append("func MapContainsKeyNode[N ~uint32](m any, n N) bool {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastNodeMap); ok { return fm.MapContainsNode_(uint32(n)) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(Ext_java_util_Map).")
        .append(mapContainsKeyGo)
        .append("(BoxNode(n))\n}\n")
        .append("func MapPutNode[N ~uint32](m any, n N, v any) any {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastNodeMap); ok { return fm.MapPutNode_(uint32(n), v) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(Ext_java_util_Map).")
        .append(mapPutGo)
        .append("(BoxNode(n), ToAny(v))\n}\n")
        .append("func MapRemoveNode[N ~uint32](m any, n N) any {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastNodeMap); ok { return fm.MapRemoveNode_(uint32(n)) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(Ext_java_util_Map).")
        .append(mapRemoveGo)
        .append("(BoxNode(n))\n}\n")
        .append("func SetContainsNode[N ~uint32](s any, n N) bool {\n")
        .append("\tif s == nil { ThrowNPE() }\n")
        .append("\tif fs, ok := s.(FastNodeSet); ok { return fs.SetContainsNode_(uint32(n)) }\n")
        .append("\tif isJavaNilAny(s) { ThrowNPE() }\n")
        .append("\treturn s.(Ext_java_util_Collection).")
        .append(colContainsGo)
        .append("(BoxNode(n))\n}\n")
        .append("func SetAddNode[N ~uint32](s any, n N) bool {\n")
        .append("\tif s == nil { ThrowNPE() }\n")
        .append("\tif fs, ok := s.(FastNodeSet); ok { return fs.SetAddNode_(uint32(n)) }\n")
        .append("\tif isJavaNilAny(s) { ThrowNPE() }\n")
        .append("\treturn s.(Ext_java_util_Collection).")
        .append(colAddGo)
        .append("(BoxNode(n))\n}\n")
        .append("func SetRemoveNode[N ~uint32](s any, n N) bool {\n")
        .append("\tif s == nil { ThrowNPE() }\n")
        .append("\tif fs, ok := s.(FastNodeSet); ok { return fs.SetRemoveNode_(uint32(n)) }\n")
        .append("\tif isJavaNilAny(s) { ThrowNPE() }\n")
        .append("\treturn s.(Ext_java_util_Collection).")
        .append(colRemoveGo)
        .append("(BoxNode(n))\n}\n")
        .append("func MapComputeIfAbsentStr(m any, k string, fn Ext_java_util_function_Function) any {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastComputeMap); ok { return fm.MapComputeIfAbsentStr_(k, fn) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(Ext_java_util_Map).")
        .append(mapComputeIfAbsentGo)
        .append("(BoxString(k), fn)\n}\n")
        .append("func MapComputeIfAbsentNode[N ~uint32](m any, n N, fn Ext_java_util_function_Function) any {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastComputeMap); ok { return fm.MapComputeIfAbsentNode_(uint32(n), fn) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(Ext_java_util_Map).")
        .append(mapComputeIfAbsentGo)
        .append("(BoxNode(n), fn)\n}\n")
        .append("func MapTryGetComputeStr(m any, k string) (any, bool) {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastComputeMap); ok { return fm.MapTryGetComputeStr_(k) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn nil, false\n}\n")
        .append("func MapTryGetComputeNode[N ~uint32](m any, n N) (any, bool) {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastComputeMap); ok { return fm.MapTryGetComputeNode_(uint32(n)) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn nil, false\n}\n")
        .append("func MapTryGetComputeAny(m any, k any) (any, bool) {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastComputeMap); ok { return fm.MapTryGetComputeAny_(k) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn nil, false\n}\n");

    String dequePushGo = findExtMethodGoName("java.util.Deque", "push", 1);
    String dequePopGo = findExtMethodGoName("java.util.Deque", "pop", 0);
    String dequeOfferLastGo = findExtMethodGoName("java.util.Deque", "offerLast", 1);
    String dequePollFirstGo = findExtMethodGoName("java.util.Deque", "pollFirst", 0);
    String dequePeekFirstGo = findExtMethodGoName("java.util.Deque", "peekFirst", 0);
    String dequePeekLastGo = findExtMethodGoName("java.util.Deque", "peekLast", 0);
    String seqAddFirstGo = findExtMethodGoName("java.util.SequencedCollection", "addFirst", 1);
    String seqAddLastGo = findExtMethodGoName("java.util.SequencedCollection", "addLast", 1);
    String seqRemoveFirstGo = findExtMethodGoName("java.util.SequencedCollection", "removeFirst", 0);
    String seqRemoveLastGo = findExtMethodGoName("java.util.SequencedCollection", "removeLast", 0);
    String queuePollGo = findExtMethodGoName("java.util.Queue", "poll", 0);
    String queuePeekGo = findExtMethodGoName("java.util.Queue", "peek", 0);

    out.append("func DequePushNode(d any, n uint32) {\n")
        .append("\tif d == nil { ThrowNPE() }\n")
        .append("\tif fd, ok := d.(jFastDequeCoreProvider); ok {\n")
        .append("\t\tif dc := fd.fastDequeCore_(); dc != nil {\n\t\t\tdc.AddFirstNode(n)\n\t\t\treturn\n\t\t}\n\t}\n")
        .append("\tif ld, ok := d.(jLinkedDequeCoreProvider); ok {\n")
        .append("\t\tif lc := ld.fastLinkedDequeCore_(); lc != nil {\n\t\t\tlc.AddAt(0, boxNodeById(n))\n\t\t\treturn\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(d) { ThrowNPE() }\n")
        .append("\tif dq, ok := d.(interface{ ").append(dequePushGo).append("(any) }); ok { dq.").append(dequePushGo).append("(boxNodeById(n)); return }\n")
        .append("\tif dq, ok := d.(interface{ ").append(seqAddFirstGo).append("(any) }); ok { dq.").append(seqAddFirstGo).append("(boxNodeById(n)); return }\n")
        .append("\tpanic(javalang.NewUnsupportedOperationException(fmt.Sprintf(\"DequePushNode unsupported on %T\", d)))\n}\n")
        .append("func DequeAddLastNode(d any, n uint32) bool {\n")
        .append("\tif d == nil { ThrowNPE() }\n")
        .append("\tif fd, ok := d.(jFastDequeCoreProvider); ok {\n")
        .append("\t\tif dc := fd.fastDequeCore_(); dc != nil {\n\t\t\tdc.AddLastNode(n)\n\t\t\treturn true\n\t\t}\n\t}\n")
        .append("\tif ld, ok := d.(jLinkedDequeCoreProvider); ok {\n")
        .append("\t\tif lc := ld.fastLinkedDequeCore_(); lc != nil {\n\t\t\tlc.Add(boxNodeById(n))\n\t\t\treturn true\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(d) { ThrowNPE() }\n")
        .append("\tif dq, ok := d.(interface{ ").append(dequeOfferLastGo).append("(any) bool }); ok { return dq.").append(dequeOfferLastGo).append("(boxNodeById(n)) }\n")
        .append("\tif dq, ok := d.(interface{ ").append(seqAddLastGo).append("(any) }); ok { dq.").append(seqAddLastGo).append("(boxNodeById(n)); return true }\n")
        .append("\treturn d.(Ext_java_util_Collection).").append(colAddGo).append("(boxNodeById(n))\n}\n")
        .append("func DequePopNode(d any) uint32 {\n")
        .append("\tif d == nil { ThrowNPE() }\n")
        .append("\tif fd, ok := d.(jFastDequeCoreProvider); ok {\n")
        .append("\t\tif dc := fd.fastDequeCore_(); dc != nil { return dc.PopNode() }\n\t}\n")
        .append("\tif ld, ok := d.(jLinkedDequeCoreProvider); ok {\n")
        .append("\t\tif lc := ld.fastLinkedDequeCore_(); lc != nil {\n\t\t\tif lc.IsEmpty() { panic(javalang.NewNoSuchElementException(\"empty\")) }\n\t\t\treturn CastNode[uint32](lc.RemoveAt(0))\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(d) { ThrowNPE() }\n")
        .append("\tif dq, ok := d.(interface{ ").append(dequePopGo).append("() any }); ok { return CastNode[uint32](dq.").append(dequePopGo).append("()) }\n")
        .append("\tif dq, ok := d.(interface{ ").append(seqRemoveFirstGo).append("() any }); ok { return CastNode[uint32](dq.").append(seqRemoveFirstGo).append("()) }\n")
        .append("\tpanic(javalang.NewUnsupportedOperationException(fmt.Sprintf(\"DequePopNode unsupported on %T\", d)))\n}\n")
        .append("func DequeRemoveLastNode(d any) uint32 {\n")
        .append("\tif d == nil { ThrowNPE() }\n")
        .append("\tif fd, ok := d.(jFastDequeCoreProvider); ok {\n")
        .append("\t\tif dc := fd.fastDequeCore_(); dc != nil { return dc.RemoveLastNode() }\n\t}\n")
        .append("\tif ld, ok := d.(jLinkedDequeCoreProvider); ok {\n")
        .append("\t\tif lc := ld.fastLinkedDequeCore_(); lc != nil {\n\t\t\tif lc.IsEmpty() { panic(javalang.NewNoSuchElementException(\"empty\")) }\n\t\t\treturn CastNode[uint32](lc.RemoveAt(lc.Size() - 1))\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(d) { ThrowNPE() }\n")
        .append("\tif dq, ok := d.(interface{ ").append(seqRemoveLastGo).append("() any }); ok { return CastNode[uint32](dq.").append(seqRemoveLastGo).append("()) }\n")
        .append("\tpanic(javalang.NewUnsupportedOperationException(fmt.Sprintf(\"DequeRemoveLastNode unsupported on %T\", d)))\n}\n")
        .append("func DequePollFirstNode(d any) uint32 {\n")
        .append("\tif d == nil { ThrowNPE() }\n")
        .append("\tif fd, ok := d.(jFastDequeCoreProvider); ok {\n")
        .append("\t\tif dc := fd.fastDequeCore_(); dc != nil { return dc.PollFirstNode() }\n\t}\n")
        .append("\tif ld, ok := d.(jLinkedDequeCoreProvider); ok {\n")
        .append("\t\tif lc := ld.fastLinkedDequeCore_(); lc != nil {\n\t\t\tif lc.IsEmpty() { return 0 }\n\t\t\treturn CastNode[uint32](lc.RemoveAt(0))\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(d) { ThrowNPE() }\n")
        .append("\tif dq, ok := d.(interface{ ").append(dequePollFirstGo).append("() any }); ok { return CastNode[uint32](dq.").append(dequePollFirstGo).append("()) }\n")
        .append("\tif dq, ok := d.(interface{ ").append(queuePollGo).append("() any }); ok { return CastNode[uint32](dq.").append(queuePollGo).append("()) }\n")
        .append("\tpanic(javalang.NewUnsupportedOperationException(fmt.Sprintf(\"DequePollFirstNode unsupported on %T\", d)))\n}\n")
        .append("func DequePeekFirstNode(d any) uint32 {\n")
        .append("\tif d == nil { ThrowNPE() }\n")
        .append("\tif fd, ok := d.(jFastDequeCoreProvider); ok {\n")
        .append("\t\tif dc := fd.fastDequeCore_(); dc != nil { return dc.PeekFirstNode() }\n\t}\n")
        .append("\tif ld, ok := d.(jLinkedDequeCoreProvider); ok {\n")
        .append("\t\tif lc := ld.fastLinkedDequeCore_(); lc != nil {\n\t\t\tif lc.IsEmpty() { return 0 }\n\t\t\treturn CastNode[uint32](lc.Get(0))\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(d) { ThrowNPE() }\n")
        .append("\tif dq, ok := d.(interface{ ").append(dequePeekFirstGo).append("() any }); ok { return CastNode[uint32](dq.").append(dequePeekFirstGo).append("()) }\n")
        .append("\tif dq, ok := d.(interface{ ").append(queuePeekGo).append("() any }); ok { return CastNode[uint32](dq.").append(queuePeekGo).append("()) }\n")
        .append("\tpanic(javalang.NewUnsupportedOperationException(fmt.Sprintf(\"DequePeekFirstNode unsupported on %T\", d)))\n}\n")
        .append("func DequePeekLastNode(d any) uint32 {\n")
        .append("\tif d == nil { ThrowNPE() }\n")
        .append("\tif fd, ok := d.(jFastDequeCoreProvider); ok {\n")
        .append("\t\tif dc := fd.fastDequeCore_(); dc != nil { return dc.PeekLastNode() }\n\t}\n")
        .append("\tif ld, ok := d.(jLinkedDequeCoreProvider); ok {\n")
        .append("\t\tif lc := ld.fastLinkedDequeCore_(); lc != nil {\n\t\t\tif lc.IsEmpty() { return 0 }\n\t\t\treturn CastNode[uint32](lc.Get(lc.Size() - 1))\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(d) { ThrowNPE() }\n")
        .append("\tif dq, ok := d.(interface{ ").append(dequePeekLastGo).append("() any }); ok { return CastNode[uint32](dq.").append(dequePeekLastGo).append("()) }\n")
        .append("\tpanic(javalang.NewUnsupportedOperationException(fmt.Sprintf(\"DequePeekLastNode unsupported on %T\", d)))\n}\n")
        .append("func DequePushColorId[R RawColorIdKey](d any, r R) {\n")
        .append("\tif d == nil { ThrowNPE() }\n")
        .append("\traw, valid := r.RawColorId_()\n")
        .append("\tif !valid { ThrowNPE() }\n")
        .append("\tif fd, ok := d.(jFastDequeCoreProvider); ok {\n")
        .append("\t\tif dc := fd.fastDequeCore_(); dc != nil {\n\t\t\tdc.AddFirstColorId(raw)\n\t\t\treturn\n\t\t}\n\t}\n")
        .append("\tif ld, ok := d.(jLinkedDequeCoreProvider); ok {\n")
        .append("\t\tif lc := ld.fastLinkedDequeCore_(); lc != nil {\n\t\t\tlc.AddAt(0, r)\n\t\t\treturn\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(d) { ThrowNPE() }\n")
        .append("\tif dq, ok := d.(interface{ ").append(dequePushGo).append("(any) }); ok { dq.").append(dequePushGo).append("(r); return }\n")
        .append("\tif dq, ok := d.(interface{ ").append(seqAddFirstGo).append("(any) }); ok { dq.").append(seqAddFirstGo).append("(r); return }\n")
        .append("\tpanic(javalang.NewUnsupportedOperationException(fmt.Sprintf(\"DequePushColorId unsupported on %T\", d)))\n}\n")
        .append("func DequeAddLastColorIdBool[R RawColorIdKey](d any, r R) bool {\n")
        .append("\tif d == nil { ThrowNPE() }\n")
        .append("\traw, valid := r.RawColorId_()\n")
        .append("\tif !valid { ThrowNPE() }\n")
        .append("\tif fd, ok := d.(jFastDequeCoreProvider); ok {\n")
        .append("\t\tif dc := fd.fastDequeCore_(); dc != nil {\n\t\t\tdc.AddLastColorId(raw)\n\t\t\treturn true\n\t\t}\n\t}\n")
        .append("\tif ld, ok := d.(jLinkedDequeCoreProvider); ok {\n")
        .append("\t\tif lc := ld.fastLinkedDequeCore_(); lc != nil {\n\t\t\tlc.Add(r)\n\t\t\treturn true\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(d) { ThrowNPE() }\n")
        .append("\tif dq, ok := d.(interface{ ").append(dequeOfferLastGo).append("(any) bool }); ok { return dq.").append(dequeOfferLastGo).append("(r) }\n")
        .append("\tif dq, ok := d.(interface{ ").append(seqAddLastGo).append("(any) }); ok { dq.").append(seqAddLastGo).append("(r); return true }\n")
        .append("\treturn d.(Ext_java_util_Collection).").append(colAddGo).append("(r)\n}\n")
        .append("func DequeAddLastColorIdVoid[R RawColorIdKey](d any, r R) {\n")
        .append("\t_ = DequeAddLastColorIdBool(d, r)\n}\n")
        .append("func DequePopVoid(d any) {\n")
        .append("\tif d == nil { ThrowNPE() }\n")
        .append("\tif fd, ok := d.(jFastDequeCoreProvider); ok {\n")
        .append("\t\tif dc := fd.fastDequeCore_(); dc != nil {\n\t\t\tdc.RemoveFirstVoid()\n\t\t\treturn\n\t\t}\n\t}\n")
        .append("\tif ld, ok := d.(jLinkedDequeCoreProvider); ok {\n")
        .append("\t\tif lc := ld.fastLinkedDequeCore_(); lc != nil {\n\t\t\tif lc.IsEmpty() { panic(javalang.NewNoSuchElementException(\"empty\")) }\n\t\t\t_ = lc.RemoveAt(0)\n\t\t\treturn\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(d) { ThrowNPE() }\n")
        .append("\tif dq, ok := d.(interface{ ").append(dequePopGo).append("() any }); ok { _ = dq.").append(dequePopGo).append("(); return }\n")
        .append("\tif dq, ok := d.(interface{ ").append(seqRemoveFirstGo).append("() any }); ok { _ = dq.").append(seqRemoveFirstGo).append("(); return }\n")
        .append("\tpanic(javalang.NewUnsupportedOperationException(fmt.Sprintf(\"DequePopVoid unsupported on %T\", d)))\n}\n")
        .append("func DequeRemoveLastVoid(d any) {\n")
        .append("\tif d == nil { ThrowNPE() }\n")
        .append("\tif fd, ok := d.(jFastDequeCoreProvider); ok {\n")
        .append("\t\tif dc := fd.fastDequeCore_(); dc != nil {\n\t\t\tdc.RemoveLastVoid()\n\t\t\treturn\n\t\t}\n\t}\n")
        .append("\tif ld, ok := d.(jLinkedDequeCoreProvider); ok {\n")
        .append("\t\tif lc := ld.fastLinkedDequeCore_(); lc != nil {\n\t\t\tif lc.IsEmpty() { panic(javalang.NewNoSuchElementException(\"empty\")) }\n\t\t\t_ = lc.RemoveAt(lc.Size() - 1)\n\t\t\treturn\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(d) { ThrowNPE() }\n")
        .append("\tif dq, ok := d.(interface{ ").append(seqRemoveLastGo).append("() any }); ok { _ = dq.").append(seqRemoveLastGo).append("(); return }\n")
        .append("\tpanic(javalang.NewUnsupportedOperationException(fmt.Sprintf(\"DequeRemoveLastVoid unsupported on %T\", d)))\n}\n")
        .append("func MapGetEqVal[K FastEqValKey](m any, k K) any {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(jFastMapCoreProvider); ok {\n")
        .append("\t\tif mc := fm.fastMapCore_(); mc != nil {\n\t\t\tv, _ := findMapCoreEqVal(mc, k)\n\t\t\treturn v\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(Ext_java_util_Map).").append(mapGetGo).append("(BoxNullableVal(k))\n}\n")
        .append("func MapGetOrDefaultEqVal[K FastEqValKey](m any, k K, def any) any {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(jFastMapCoreProvider); ok {\n")
        .append("\t\tif mc := fm.fastMapCore_(); mc != nil {\n\t\t\tif v, found := findMapCoreEqVal(mc, k); found {\n\t\t\t\treturn v\n\t\t\t}\n\t\t\treturn ToAny(def)\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(Ext_java_util_Map).").append(mapGetOrDefGo).append("(BoxNullableVal(k), ToAny(def))\n}\n")
        .append("func MapContainsEqVal[K FastEqValKey](m any, k K) bool {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(jFastMapCoreProvider); ok {\n")
        .append("\t\tif mc := fm.fastMapCore_(); mc != nil {\n\t\t\t_, found := findMapCoreEqVal(mc, k)\n\t\t\treturn found\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(Ext_java_util_Map).").append(mapContainsKeyGo).append("(BoxNullableVal(k))\n}\n")
        .append("func SetContainsEqVal[K FastEqValKey](s any, k K) bool {\n")
        .append("\tif s == nil { ThrowNPE() }\n")
        .append("\tif hit, handled := fastSetContainsEqVal(s, k); handled { return hit }\n")
        .append("\tif isJavaNilAny(s) { ThrowNPE() }\n")
        .append("\treturn s.(Ext_java_util_Collection).").append(colContainsGo).append("(BoxNullableVal(k))\n}\n\n");

    String tblContainsRowGo = findExtMethodGoName("com.google.common.collect.Table", "containsRow", 1);
    String tblContainsGo = findExtMethodGoName("com.google.common.collect.Table", "contains", 2);
    String tblGetGo = findExtMethodGoName("com.google.common.collect.Table", "get", 2);
    String tblPutGo = findExtMethodGoName("com.google.common.collect.Table", "put", 3);
    String tblRemoveGo = findExtMethodGoName("com.google.common.collect.Table", "remove", 2);
    String tblRowGo = findExtMethodGoName("com.google.common.collect.Table", "row", 1);

    out.append("func TableContainsRowColorId[R RawColorIdKey](tbl any, r R) bool {\n")
        .append("\tif tbl == nil { ThrowNPE() }\n")
        .append("\traw, valid := r.RawColorId_()\n")
        .append("\tif ft, ok := tbl.(jFastTableCoreProvider); ok {\n")
        .append("\t\tif tc, _ := ft.fastTableCore_(); tc != nil {\n")
        .append("\t\t\tif !valid { return false }\n")
        .append("\t\t\treturn tc.containsRowColorId(raw)\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(tbl) { ThrowNPE() }\n")
        .append("\treturn tbl.(interface{ ").append(tblContainsRowGo).append("(any) bool }).").append(tblContainsRowGo).append("(BoxNullableVal(r))\n}\n")
        .append("func TableContainsColorId[R RawColorIdKey](tbl any, r R, c any) bool {\n")
        .append("\tif tbl == nil { ThrowNPE() }\n")
        .append("\traw, valid := r.RawColorId_()\n")
        .append("\tif ft, ok := tbl.(jFastTableCoreProvider); ok {\n")
        .append("\t\tif tc, _ := ft.fastTableCore_(); tc != nil {\n")
        .append("\t\t\tif !valid { return false }\n")
        .append("\t\t\treturn tc.containsColorId(raw, c)\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(tbl) { ThrowNPE() }\n")
        .append("\treturn tbl.(interface{ ").append(tblContainsGo).append("(any, any) bool }).").append(tblContainsGo).append("(BoxNullableVal(r), c)\n}\n")
        .append("func TableGetColorId[R RawColorIdKey](tbl any, r R, c any) any {\n")
        .append("\tif tbl == nil { ThrowNPE() }\n")
        .append("\traw, valid := r.RawColorId_()\n")
        .append("\tif ft, ok := tbl.(jFastTableCoreProvider); ok {\n")
        .append("\t\tif tc, _ := ft.fastTableCore_(); tc != nil {\n")
        .append("\t\t\tif !valid { return nil }\n")
        .append("\t\t\treturn tc.getColorId(raw, c)\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(tbl) { ThrowNPE() }\n")
        .append("\treturn tbl.(interface{ ").append(tblGetGo).append("(any, any) any }).").append(tblGetGo).append("(BoxNullableVal(r), c)\n}\n")
        .append("func TablePutColorId[R RawColorIdKey](tbl any, r R, c, v any) any {\n")
        .append("\tif tbl == nil { ThrowNPE() }\n")
        .append("\traw, valid := r.RawColorId_()\n")
        .append("\tif ft, ok := tbl.(jFastTableCoreProvider); ok {\n")
        .append("\t\ttc, imm := ft.fastTableCore_()\n")
        .append("\t\tif imm { panic(javalang.NewUnsupportedOperationException(\"ImmutableTable\")) }\n")
        .append("\t\tif tc != nil {\n")
        .append("\t\t\tif !valid { ThrowNPE() }\n")
        .append("\t\t\treturn tc.putColorId(raw, c, v)\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(tbl) { ThrowNPE() }\n")
        .append("\treturn tbl.(interface{ ").append(tblPutGo).append("(any, any, any) any }).").append(tblPutGo).append("(BoxNullableVal(r), c, v)\n}\n")
        .append("func TableRemoveColorId[R RawColorIdKey](tbl any, r R, c any) any {\n")
        .append("\tif tbl == nil { ThrowNPE() }\n")
        .append("\traw, valid := r.RawColorId_()\n")
        .append("\tif ft, ok := tbl.(jFastTableCoreProvider); ok {\n")
        .append("\t\ttc, imm := ft.fastTableCore_()\n")
        .append("\t\tif imm { panic(javalang.NewUnsupportedOperationException(\"ImmutableTable\")) }\n")
        .append("\t\tif tc != nil {\n")
        .append("\t\t\tif !valid { return nil }\n")
        .append("\t\t\treturn tc.removeColorId(raw, c)\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(tbl) { ThrowNPE() }\n")
        .append("\treturn tbl.(interface{ ").append(tblRemoveGo).append("(any, any) any }).").append(tblRemoveGo).append("(BoxNullableVal(r), c)\n}\n")
        .append("func TableRowColorId[R RawColorIdKey](tbl any, r R) Ext_java_util_Map {\n")
        .append("\tif tbl == nil { ThrowNPE() }\n")
        .append("\traw, valid := r.RawColorId_()\n")
        .append("\tif ft, ok := tbl.(jFastTableCoreProvider); ok {\n")
        .append("\t\tif tc, _ := ft.fastTableCore_(); tc != nil {\n")
        .append("\t\t\tif !valid { ThrowNPE() }\n")
        .append("\t\t\treturn tc.rowColorId(raw)\n\t\t}\n\t}\n")
        .append("\tif isJavaNilAny(tbl) { ThrowNPE() }\n")
        .append("\treturn tbl.(interface{ ").append(tblRowGo).append("(any) Ext_java_util_Map }).").append(tblRowGo).append("(BoxNullableVal(r))\n}\n");

    String listGetGo = findExtMethodGoName("java.util.List", "get", 1);
    String mmContainsKeyGo = findExtMethodGoName("com.google.common.collect.Multimap", "containsKey", 1);
    String mmContainsEntryGo = findExtMethodGoName("com.google.common.collect.Multimap", "containsEntry", 2);
    String mmGetGo = findExtMethodGoName("com.google.common.collect.Multimap", "get", 1);
    String mmPutGo = findExtMethodGoName("com.google.common.collect.Multimap", "put", 2);
    String mmRemoveGo = findExtMethodGoName("com.google.common.collect.Multimap", "remove", 2);
    String mmRemoveAllGo = findExtMethodGoName("com.google.common.collect.Multimap", "removeAll", 1);
    String mmPutAllGo = findExtMethodGoName("com.google.common.collect.Multimap", "putAll", 2);
    String msCountGo = findExtMethodGoName("com.google.common.collect.Multiset", "count", 1);
    String msContainsGo = findExtMethodGoName("com.google.common.collect.Multiset", "contains", 1);
    String msAdd1Go = findExtMethodGoName("com.google.common.collect.Multiset", "add", 1);
    String msAdd2Go = findExtMethodGoName("com.google.common.collect.Multiset", "add", 2);
    String msRemove1Go = findExtMethodGoName("com.google.common.collect.Multiset", "remove", 1);
    String msSetCount3Go = findExtMethodGoName("com.google.common.collect.Multiset", "setCount", 3);

    out.append("func ListGetStr(l any, i int32) string {\n")
        .append("\tif l == nil { ThrowNPE() }\n")
        .append("\tif fl, ok := l.(FastStrList); ok { return fl.ListGetStr_(i) }\n")
        .append("\tif isJavaNilAny(l) { ThrowNPE() }\n")
        .append("\treturn CastString(l.(Ext_java_util_List).").append(listGetGo).append("(i))\n}\n")
        .append("func MultimapContainsKeyStr(m any, k string) bool {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultimap); ok { return fm.MultimapContainsKeyStr_(k) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(interface{ ").append(mmContainsKeyGo).append("(any) bool }).").append(mmContainsKeyGo).append("(BoxString(k))\n}\n")
        .append("func MultimapContainsKeyNode[N ~uint32](m any, n N) bool {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultimap); ok { return fm.MultimapContainsKeyNode_(uint32(n)) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(interface{ ").append(mmContainsKeyGo).append("(any) bool }).").append(mmContainsKeyGo).append("(BoxNode(n))\n}\n")
        .append("func MultimapContainsEntryStr(m any, k string, v any) bool {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultimap); ok { return fm.MultimapContainsEntryStr_(k, v) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(interface{ ").append(mmContainsEntryGo).append("(any, any) bool }).").append(mmContainsEntryGo).append("(BoxString(k), ToAny(v))\n}\n")
        .append("func MultimapContainsEntryNode[N ~uint32](m any, n N, v any) bool {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultimap); ok { return fm.MultimapContainsEntryNode_(uint32(n), v) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(interface{ ").append(mmContainsEntryGo).append("(any, any) bool }).").append(mmContainsEntryGo).append("(BoxNode(n), ToAny(v))\n}\n")
        .append("func MultimapGetStr(m any, k string) any {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultimap); ok { return fm.MultimapGetStr_(k) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\tif gm, ok := m.(interface{ ").append(mmGetGo).append("(any) any }); ok { return gm.").append(mmGetGo).append("(BoxString(k)) }\n")
        .append("\treturn m.(interface{ ").append(mmGetGo).append("(any) Ext_java_util_Collection }).").append(mmGetGo).append("(BoxString(k))\n}\n")
        .append("func MultimapGetNode[N ~uint32](m any, n N) any {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultimap); ok { return fm.MultimapGetNode_(uint32(n)) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\tif gm, ok := m.(interface{ ").append(mmGetGo).append("(any) any }); ok { return gm.").append(mmGetGo).append("(BoxNode(n)) }\n")
        .append("\treturn m.(interface{ ").append(mmGetGo).append("(any) Ext_java_util_Collection }).").append(mmGetGo).append("(BoxNode(n))\n}\n")
        .append("func MultimapPutStr(m any, k string, v any) bool {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultimap); ok { return fm.MultimapPutStr_(k, v) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(interface{ ").append(mmPutGo).append("(any, any) bool }).").append(mmPutGo).append("(BoxString(k), ToAny(v))\n}\n")
        .append("func MultimapPutNode[N ~uint32](m any, n N, v any) bool {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultimap); ok { return fm.MultimapPutNode_(uint32(n), v) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(interface{ ").append(mmPutGo).append("(any, any) bool }).").append(mmPutGo).append("(BoxNode(n), ToAny(v))\n}\n")
        .append("func MultimapRemoveStr(m any, k string, v any) bool {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultimap); ok { return fm.MultimapRemoveStr_(k, v) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(interface{ ").append(mmRemoveGo).append("(any, any) bool }).").append(mmRemoveGo).append("(BoxString(k), ToAny(v))\n}\n")
        .append("func MultimapRemoveNode[N ~uint32](m any, n N, v any) bool {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultimap); ok { return fm.MultimapRemoveNode_(uint32(n), v) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(interface{ ").append(mmRemoveGo).append("(any, any) bool }).").append(mmRemoveGo).append("(BoxNode(n), ToAny(v))\n}\n")
        .append("func MultimapRemoveAllStr(m any, k string) any {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultimap); ok { return fm.MultimapRemoveAllStr_(k) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\tif gm, ok := m.(interface{ ").append(mmRemoveAllGo).append("(any) any }); ok { return gm.").append(mmRemoveAllGo).append("(BoxString(k)) }\n")
        .append("\treturn m.(interface{ ").append(mmRemoveAllGo).append("(any) Ext_java_util_Collection }).").append(mmRemoveAllGo).append("(BoxString(k))\n}\n")
        .append("func MultimapRemoveAllNode[N ~uint32](m any, n N) any {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultimap); ok { return fm.MultimapRemoveAllNode_(uint32(n)) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\tif gm, ok := m.(interface{ ").append(mmRemoveAllGo).append("(any) any }); ok { return gm.").append(mmRemoveAllGo).append("(BoxNode(n)) }\n")
        .append("\treturn m.(interface{ ").append(mmRemoveAllGo).append("(any) Ext_java_util_Collection }).").append(mmRemoveAllGo).append("(BoxNode(n))\n}\n")
        .append("func MultimapPutAllStr(m any, k string, vals Ext_java_lang_Iterable) bool {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultimap); ok { return fm.MultimapPutAllStr_(k, vals) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\tif pm, ok := m.(interface{ ").append(mmPutAllGo).append("(any, Ext_java_lang_Iterable) bool }); ok { return pm.").append(mmPutAllGo).append("(BoxString(k), vals) }\n")
        .append("\treturn m.(interface{ ").append(mmPutAllGo).append("(any, any) bool }).").append(mmPutAllGo).append("(BoxString(k), vals)\n}\n")
        .append("func MultimapPutAllNode[N ~uint32](m any, n N, vals Ext_java_lang_Iterable) bool {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultimap); ok { return fm.MultimapPutAllNode_(uint32(n), vals) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\tif pm, ok := m.(interface{ ").append(mmPutAllGo).append("(any, Ext_java_lang_Iterable) bool }); ok { return pm.").append(mmPutAllGo).append("(BoxNode(n), vals) }\n")
        .append("\treturn m.(interface{ ").append(mmPutAllGo).append("(any, any) bool }).").append(mmPutAllGo).append("(BoxNode(n), vals)\n}\n")
        .append("func MultisetCountStr(m any, s string) int32 {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultiset); ok { return fm.MultisetCountStr_(s) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(interface{ ").append(msCountGo).append("(any) int32 }).").append(msCountGo).append("(BoxString(s))\n}\n")
        .append("func MultisetContainsStr(m any, s string) bool {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultiset); ok { return fm.MultisetContainsStr_(s) }\n")
        .append("\tif fs, ok := m.(FastStrSet); ok { return fs.SetContainsStr_(s) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(interface{ ").append(msContainsGo).append("(any) bool }).").append(msContainsGo).append("(BoxString(s))\n}\n")
        .append("func MultisetAddStr(m any, s string) bool {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultiset); ok { return fm.MultisetAddStr1_(s) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(interface{ ").append(msAdd1Go).append("(any) bool }).").append(msAdd1Go).append("(BoxString(s))\n}\n")
        .append("func MultisetAddCountStr(m any, s string, n int32) int32 {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultiset); ok { return fm.MultisetAddStr2_(s, n) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(interface{ ").append(msAdd2Go).append("(any, int32) int32 }).").append(msAdd2Go).append("(BoxString(s), n)\n}\n")
        .append("func MultisetRemoveStr(m any, s string) bool {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultiset); ok { return fm.MultisetRemoveStr1_(s) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(interface{ ").append(msRemove1Go).append("(any) bool }).").append(msRemove1Go).append("(BoxString(s))\n}\n")
        .append("func MultisetSetCount3Str(m any, s string, oldCount, newCount int32) bool {\n")
        .append("\tif m == nil { ThrowNPE() }\n")
        .append("\tif fm, ok := m.(FastMultiset); ok { return fm.MultisetSetCountStr3_(s, oldCount, newCount) }\n")
        .append("\tif isJavaNilAny(m) { ThrowNPE() }\n")
        .append("\treturn m.(interface{ ").append(msSetCount3Go).append("(any, int32, int32) bool }).").append(msSetCount3Go).append("(BoxString(s), oldCount, newCount)\n}\n\n");

    TypeElement biTe = resolver.elements.getTypeElement("java.math.BigInteger");
    if (biTe != null && resolver.getReferencedExternalTypes().contains(biTe)) {
      String st = goStructType(biTe, imports);
      out.append("func wrapBigInteger(bi *javamath.BigInteger) *").append(st).append(" {\n");
      out.append("\tif bi == nil { return nil }\n");
      out.append("\tobj := &").append(st).append("{}\n");
      out.append("\tobj.BindExtSelf_(obj)\n");
      out.append("\tobj.Core_ = bi\n\treturn obj\n}\n");
      out.append("func (b *").append(st).append(") bi() *javamath.BigInteger {\n");
      out.append("\tCheckNotNull(b)\n");
      out.append("\tif v, ok := b.Core_.(*javamath.BigInteger); ok && v != nil { return v }\n");
      out.append("\treturn javamath.BigInteger_ZERO\n}\n\n");
    }
  }

  private void emitImmListVariantsAndBuilder(
      TypeElement ilTe, PackagePlanner.ImportTracker imports, StringBuilder out) {
    String st = goStructType(ilTe, imports);
    String ilGo = goTypeInJ2gort(ilTe.asType(), imports);
    out.append("var emptyImmutableList_ ").append(ilGo).append("\n");
    out.append("func init() {\n\temptyImmutableList_ = &")
        .append(st)
        .append("{items_: emptyAnySlice_}\n}\n\n");
    out.append("func newImmutableListShell() *").append(st).append(" {\n");
    out.append("\treturn &").append(st).append("{}\n}\n\n");
    out.append("func newImmListBaseOwn(items []any) *").append(st).append(" {\n");
    out.append("\tfor i, v := range items { items[i] = ToAny(v) }\n");
    out.append("\treturn &").append(st).append("{items_: items[:len(items):len(items)]}\n}\n\n");

    // Collect markers and nil-Base getters required by Ext_com_google_common_collect_ImmutableList
    Set<String> immListMarkers = new LinkedHashSet<>();
    immListMarkers.add(resolver.getInterfaceMarkerMethod(ilTe));
    List<String> nilBaseGetters = new ArrayList<>();
    if (isSubclassedExternal(ilTe)) {
      nilBaseGetters.add(
          "Get" + extGoTypeName(ilTe) + "Base() *" + st + " { return nil }");
    }
    for (TypeElement anc : resolver.getAllAncestorTypes(ilTe)) {
      if (!mapsToAny(anc)
          && !isShimException(anc)
          && !anc.getQualifiedName().contentEquals("java.lang.Number")) {
        if (anc.getKind() == ElementKind.INTERFACE || isSubclassedExternal(anc)) {
          immListMarkers.add(resolver.getInterfaceMarkerMethod(anc));
        }
        if (anc.getKind() == ElementKind.CLASS
            && isSubclassedExternal(anc)
            && resolver.getReferencedExternalTypes().contains(anc)) {
          nilBaseGetters.add(
              "Get"
                  + extGoTypeName(anc)
                  + "Base() *"
                  + goStructType(anc, imports)
                  + " { return nil }");
        }
      }
    }
    immListMarkers.add("Implements_j2gort_Ext_java_util_RandomAccess");
    immListMarkers.add("Implements_j2gort_Ext_java_io_Serializable");

    List<TypeElement> ilChain = new ArrayList<>();
    ilChain.add(ilTe);
    ilChain.addAll(resolver.getAllAncestorTypes(ilTe));
    String intFnApply = findExtMethodGoName("java.util.function.IntFunction", "apply", 1);
    String consumerAccept = findExtMethodGoName("java.util.function.Consumer", "accept", 1);

    // 0. Zero-alloc direct-iface 1-pointer ImmutableList: jImmListPtr[T]
    out.append("type jImmListPtr[T any] struct {\n\tp *T\n}\n");
    out.append("func ImmListOfPtr[T any](p *T) ")
        .append(ilGo)
        .append(" {\n\tif p == nil { ThrowNPE() }\n\treturn jImmListPtr[T]{p: p}\n}\n");
    out.append("func (l jImmListPtr[T]) IsJavaNil_() bool { return l.p == nil }\n");
    out.append(
        "func (jImmListPtr[T]) JavaClassInfo_() (string, string) { return \"com.google.common.collect.ImmutableList\", \"ImmutableList\" }\n");
    out.append(
        "func (l jImmListPtr[T]) singleNonNullElem_() any {\n\tif l.p == nil { ThrowNPE() }\n\treturn l.p\n}\n");
    out.append(
        "func (l jImmListPtr[T]) typedViewLen_() (int, bool) {\n\tif l.p == nil { ThrowNPE() }\n\treturn 1, true\n}\n");
    out.append("func (l jImmListPtr[T]) typedViewAt_(int) any { return l.p }\n");
    out.append("func (jImmListPtr[T]) typedViewNodeAt_(int) (uint32, bool) { return 0, false }\n");
    out.append(
        "func (l jImmListPtr[T]) IterSlice_() []any {\n\tif l.p == nil { ThrowNPE() }\n\treturn []any{l.p}\n}\n");
    out.append(
        "func (l jImmListPtr[T]) ListGetStr_(i int32) string {\n\tif l.p == nil { ThrowNPE() }\n\tCheckElementIndex(i, 1)\n\treturn CastString(l.p)\n}\n");
    out.append(
        "func (l jImmListPtr[T]) SetContainsStr_(k string) bool {\n\tif l.p == nil { ThrowNPE() }\n\t_ = k\n\treturn false\n}\n");
    out.append(
        "func (l jImmListPtr[T]) SetAddStr_(string) bool {\n\tif l.p == nil { ThrowNPE() }\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
    out.append(
        "func (l jImmListPtr[T]) SetRemoveStr_(string) bool {\n\tif l.p == nil { ThrowNPE() }\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
    out.append(
        "func (l jImmListPtr[T]) SetContainsNode_(u uint32) bool {\n\tif l.p == nil { ThrowNPE() }\n\treturn singletonElemMatchesNode(l.p, u)\n}\n");
    out.append(
        "func (l jImmListPtr[T]) SetAddNode_(uint32) bool {\n\tif l.p == nil { ThrowNPE() }\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
    out.append(
        "func (l jImmListPtr[T]) SetRemoveNode_(uint32) bool {\n\tif l.p == nil { ThrowNPE() }\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
    out.append(
        "func (l jImmListPtr[T]) Equals(arg0 any) bool { var a [1]any = [1]any{l.p}; return listEqualsSlice(a[:], arg0) }\n");
    out.append(
        "func (l jImmListPtr[T]) HashCode() int32 { return 31 + elemJavaHashCode(l.p, false) }\n");
    out.append("func (l jImmListPtr[T]) String() string { return \"[\" + StrVal(l.p) + \"]\" }\n");
    for (String mk : immListMarkers) {
      out.append("func (jImmListPtr[T]) ").append(mk).append("() {}\n");
    }
    for (String bg : nilBaseGetters) {
      out.append("func (jImmListPtr[T]) ").append(bg).append('\n');
    }
    {
      Set<String> emittedPtr = new HashSet<>(Set.of("Equals", "HashCode", "String", "IterSlice_"));
      for (TypeElement cur : ilChain) {
        if (mapsToAny(cur) || isShimException(cur)) continue;
        for (Element m : cur.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && !ee.getModifiers().contains(Modifier.STATIC)
              && !ee.getModifiers().contains(Modifier.PRIVATE)) {
            TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
            if (sig == null
                || !resolver.isNeededExternalSig(sig)
                || !emittedPtr.add(sig.goName)) {
              continue;
            }
            int arity = sig.canonParamTypes.size();
            String p0 = arity >= 1 ? goTypeInJ2gort(sig.canonParamTypes.get(0), imports) : "";
            out.append("func (l jImmListPtr[T]) ")
                .append(sig.goName)
                .append(formatExtSigWithParamNames(sig, imports))
                .append(" {\n");
            switch (sig.simpleName) {
              case "size" -> out.append("\treturn 1\n");
              case "isEmpty" -> out.append("\treturn false\n");
              case "get" -> out.append("\tCheckElementIndex(arg0, 1)\n\treturn l.p\n");
              case "getFirst", "element", "getLast" -> out.append("\treturn l.p\n");
              case "contains" -> out.append("\treturn jSetElemEquals(l.p, arg0)\n");
              case "indexOf", "lastIndexOf" ->
                  out.append("\tif jSetElemEquals(l.p, arg0) { return 0 }\n\treturn -1\n");
              case "containsAll" ->
                  out.append(
                      "\treturn colContainsAll(func(x any) bool { return jSetElemEquals(l.p, x) }, arg0)\n");
              case "iterator" -> out.append("\treturn newSingleObjIter(l.p)\n");
              case "listIterator" ->
                  out.append(
                      arity == 0
                          ? "\treturn newImmListIter([]any{l.p}, 0)\n"
                          : "\treturn newImmListIter([]any{l.p}, arg0)\n");
              case "subList" ->
                  out.append(
                      "\tCheckPositionIndexes(arg0, arg1, 1)\n\tif arg0 == arg1 { return emptyImmutableList_ }\n\treturn l\n");
              case "asList", "reverse", "reversed" -> out.append("\treturn l\n");
              case "toArray" -> {
                if (arity == 0) {
                  out.append("\treturn []any{l.p}\n");
                } else if (p0.startsWith("[]")) {
                  out.append("\treturn sliceToArray1([]any{l.p}, arg0)\n");
                } else {
                  out.append("\tCheckNotNull(arg0)\n\treturn sliceToArray1([]any{l.p}, Cast[[]any](arg0.")
                      .append(intFnApply)
                      .append("(1)))\n");
                }
              }
              case "forEach" ->
                  out.append("\tCheckNotNull(arg0)\n\targ0.")
                      .append(consumerAccept)
                      .append("(l.p)\n");
              case "stream", "parallelStream" ->
                  out.append("\treturn streamFromSlice([]any{l.p})\n");
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
                  "replaceAll" ->
                  out.append("\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n");
              default ->
                  out.append("\tpanic(")
                      .append(GoNames.goStringLiteral("UNIMPLEMENTED_EXT: jImmListPtr." + sig.goName))
                      .append(")\n");
            }
            out.append("}\n");
          }
        }
      }
      out.append('\n');
    }

    // 1. Fixed-size single-allocation any lists: *jImmList1..*jImmList8
    for (int n = 1; n <= 8; n++) {
      String tn = "jImmList" + n;
      out.append("type ").append(tn).append(" struct {\n\telem [").append(n).append("]any\n}\n");
      out.append("func (l *").append(tn).append(") IsJavaNil_() bool { return l == nil }\n");
      out.append("func (*")
          .append(tn)
          .append(") JavaClassInfo_() (string, string) { return \"com.google.common.collect.ImmutableList\", \"ImmutableList\" }\n");
      out.append("func (l *")
          .append(tn)
          .append(") IterSlice_() []any {\n\tif l == nil { ThrowNPE() }\n\treturn l.elem[:]\n}\n");
      out.append("func (l *")
          .append(tn)
          .append(") ListGetStr_(i int32) string {\n\tif l == nil { ThrowNPE() }\n\tCheckElementIndex(i, ")
          .append(n)
          .append(")\n\treturn CastString(l.elem[i])\n}\n");
      out.append("func (l *")
          .append(tn)
          .append(") SetContainsStr_(k string) bool {\n\tif l == nil { ThrowNPE() }\n\tif k == NullString { return false }\n\tfor i := 0; i < ")
          .append(n)
          .append("; i++ {\n\t\tif es, ok := l.elem[i].(string); ok && es == k { return true }\n\t}\n\treturn false\n}\n");
      out.append("func (l *").append(tn).append(") SetAddStr_(string) bool {\n\tif l == nil { ThrowNPE() }\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      out.append("func (l *").append(tn).append(") SetRemoveStr_(string) bool {\n\tif l == nil { ThrowNPE() }\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      out.append("func (l *")
          .append(tn)
          .append(") SetContainsNode_(u uint32) bool {\n\tif l == nil { ThrowNPE() }\n\tif u == 0 { return false }\n\tfor i := 0; i < ")
          .append(n)
          .append("; i++ {\n\t\tif singletonElemMatchesNode(l.elem[i], u) { return true }\n\t}\n\treturn false\n}\n");
      out.append("func (l *").append(tn).append(") SetAddNode_(uint32) bool {\n\tif l == nil { ThrowNPE() }\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      out.append("func (l *").append(tn).append(") SetRemoveNode_(uint32) bool {\n\tif l == nil { ThrowNPE() }\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      out.append("func (l *").append(tn).append(") Equals(arg0 any) bool { return listEqualsSlice(l.elem[:], arg0) }\n");
      out.append("func (l *").append(tn).append(") HashCode() int32 { return listHashSlice(l.elem[:]) }\n");
      out.append("func (l *").append(tn).append(") String() string { return sliceJavaString(l.elem[:]) }\n");
      for (String mk : immListMarkers) {
        out.append("func (*").append(tn).append(") ").append(mk).append("() {}\n");
      }
      for (String bg : nilBaseGetters) {
        out.append("func (*").append(tn).append(") ").append(bg).append('\n');
      }
      Set<String> emitted = new HashSet<>(Set.of("Equals", "HashCode", "String", "IterSlice_"));
      for (TypeElement cur : ilChain) {
        if (mapsToAny(cur) || isShimException(cur)) continue;
        for (Element m : cur.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && !ee.getModifiers().contains(Modifier.STATIC)
              && !ee.getModifiers().contains(Modifier.PRIVATE)) {
            TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
            if (sig == null || !resolver.isNeededExternalSig(sig) || !emitted.add(sig.goName)) {
              continue;
            }
            int arity = sig.canonParamTypes.size();
            String p0 = arity >= 1 ? goTypeInJ2gort(sig.canonParamTypes.get(0), imports) : "";
            out.append("func (l *")
                .append(tn)
                .append(") ")
                .append(sig.goName)
                .append(formatExtSigWithParamNames(sig, imports))
                .append(" {\n");
            switch (sig.simpleName) {
              case "size" -> out.append("\treturn ").append(n).append('\n');
              case "isEmpty" -> out.append("\treturn false\n");
              case "get" -> out.append("\tCheckElementIndex(arg0, ").append(n).append(")\n\treturn l.elem[arg0]\n");
              case "getFirst", "element" -> out.append("\treturn l.elem[0]\n");
              case "getLast" -> out.append("\treturn l.elem[").append(n - 1).append("]\n");
              case "contains" -> out.append("\treturn immSliceContains(l.elem[:], arg0)\n");
              case "indexOf" -> out.append("\treturn immSliceIndexOf(l.elem[:], arg0)\n");
              case "lastIndexOf" -> out.append("\treturn immSliceLastIndexOf(l.elem[:], arg0)\n");
              case "containsAll" ->
                  out.append("\treturn colContainsAll(func(x any) bool { return immSliceContains(l.elem[:], x) }, arg0)\n");
              case "iterator" ->
                  out.append(
                      n == 1
                          ? "\treturn newSingleObjIter(l.elem[0])\n"
                          : "\treturn newSliceObjIter(l.elem[:])\n");
              case "listIterator" ->
                  out.append(arity == 0 ? "\treturn newImmListIter(l.elem[:], 0)\n" : "\treturn newImmListIter(l.elem[:], arg0)\n");
              case "subList" ->
                  out.append("\tCheckPositionIndexes(arg0, arg1, ")
                      .append(n)
                      .append(")\n\tif arg0 == 0 && arg1 == ")
                      .append(n)
                      .append(" { return l }\n\treturn wrapImmutableList(l.elem[arg0:arg1])\n");
              case "asList" -> out.append("\treturn l\n");
              case "reverse", "reversed" -> {
                if (n == 1) {
                  out.append("\treturn l\n");
                } else {
                  out.append("\treturn &").append(tn).append("{elem: [").append(n).append("]any{");
                  for (int j = n - 1; j >= 0; j--) {
                    if (j < n - 1) out.append(", ");
                    out.append("l.elem[").append(j).append(']');
                  }
                  out.append("}}\n");
                }
              }
              case "toArray" -> {
                if (arity == 0) {
                  out.append("\tout := make([]any, ").append(n).append(")\n\tcopy(out, l.elem[:])\n\treturn out\n");
                } else if (p0.startsWith("[]")) {
                  out.append("\treturn sliceToArray1(l.elem[:], arg0)\n");
                } else {
                  out.append("\tCheckNotNull(arg0)\n\treturn sliceToArray1(l.elem[:], Cast[[]any](arg0.")
                      .append(intFnApply)
                      .append('(')
                      .append(n)
                      .append(")))\n");
                }
              }
              case "forEach" ->
                  out.append("\tCheckNotNull(arg0)\n\tfor i := 0; i < ")
                      .append(n)
                      .append("; i++ { arg0.")
                      .append(consumerAccept)
                      .append("(l.elem[i]) }\n");
              case "stream", "parallelStream" -> out.append("\treturn streamFromSlice(l.elem[:])\n");
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
                  "replaceAll" ->
                  out.append("\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n");
              default ->
                  out.append("\tpanic(")
                      .append(GoNames.goStringLiteral("UNIMPLEMENTED_EXT: " + tn + "." + sig.goName))
                      .append(")\n");
            }
            out.append("}\n");
          }
        }
      }
      out.append('\n');
    }

    // 2. Unboxed string lists: *jImmStrList1..*jImmStrList8 (fixed, 1 alloc) and *jImmStrList (slice, >=9)
    for (int n = 0; n <= 8; n++) {
      String tn = n == 0 ? "jImmStrList" : ("jImmStrList" + n);
      String fieldDecl = n == 0 ? "strs []string" : ("strs [" + n + "]string");
      String sliceExpr = n == 0 ? "l.strs" : "l.strs[:]";
      out.append("type ").append(tn).append(" struct {\n\t").append(fieldDecl).append("\n}\n");
      out.append("func (l *").append(tn).append(") IsJavaNil_() bool { return l == nil }\n");
      out.append("func (*")
          .append(tn)
          .append(") JavaClassInfo_() (string, string) { return \"com.google.common.collect.ImmutableList\", \"ImmutableList\" }\n");
      out.append("func (l *")
          .append(tn)
          .append(") IterStrSlice_() ([]string, bool) {\n\tif l == nil { ThrowNPE() }\n\treturn ")
          .append(sliceExpr)
          .append(", true\n}\n");
      out.append("func (l *")
          .append(tn)
          .append(") IterSlice_() []any {\n\tif l == nil { ThrowNPE() }\n\treturn strsToAnySlice(")
          .append(sliceExpr)
          .append(")\n}\n");
      out.append("func (l *")
          .append(tn)
          .append(") ListGetStr_(i int32) string {\n\tif l == nil { ThrowNPE() }\n\tCheckElementIndex(i, int32(len(l.strs)))\n\treturn l.strs[i]\n}\n");
      out.append("func (l *")
          .append(tn)
          .append(") SetContainsStr_(k string) bool {\n\tif l == nil { ThrowNPE() }\n\tif k == NullString { return false }\n\tfor _, s := range l.strs {\n\t\tif s == k { return true }\n\t}\n\treturn false\n}\n");
      out.append("func (l *").append(tn).append(") SetAddStr_(string) bool {\n\tif l == nil { ThrowNPE() }\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      out.append("func (l *").append(tn).append(") SetRemoveStr_(string) bool {\n\tif l == nil { ThrowNPE() }\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      out.append("func (l *").append(tn).append(") SetContainsNode_(uint32) bool {\n\tif l == nil { ThrowNPE() }\n\treturn false\n}\n");
      out.append("func (l *").append(tn).append(") SetAddNode_(uint32) bool {\n\tif l == nil { ThrowNPE() }\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      out.append("func (l *").append(tn).append(") SetRemoveNode_(uint32) bool {\n\tif l == nil { ThrowNPE() }\n\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n}\n");
      out.append("func (l *").append(tn).append(") Equals(arg0 any) bool { return listEqualsStrSlice(").append(sliceExpr).append(", arg0) }\n");
      out.append("func (l *").append(tn).append(") HashCode() int32 { return listHashStrSlice(").append(sliceExpr).append(") }\n");
      out.append("func (l *").append(tn).append(") String() string { return strSliceJavaString(").append(sliceExpr).append(") }\n");
      for (String mk : immListMarkers) {
        out.append("func (*").append(tn).append(") ").append(mk).append("() {}\n");
      }
      for (String bg : nilBaseGetters) {
        out.append("func (*").append(tn).append(") ").append(bg).append('\n');
      }
      Set<String> emitted = new HashSet<>(Set.of("Equals", "HashCode", "String", "IterSlice_"));
      for (TypeElement cur : ilChain) {
        if (mapsToAny(cur) || isShimException(cur)) continue;
        for (Element m : cur.getEnclosedElements()) {
          if (m instanceof ExecutableElement ee
              && ee.getKind() == ElementKind.METHOD
              && !ee.getModifiers().contains(Modifier.STATIC)
              && !ee.getModifiers().contains(Modifier.PRIVATE)) {
            TypeResolver.MethodSig sig = resolver.getInstanceMethodSig(ee);
            if (sig == null || !resolver.isNeededExternalSig(sig) || !emitted.add(sig.goName)) {
              continue;
            }
            int arity = sig.canonParamTypes.size();
            String p0 = arity >= 1 ? goTypeInJ2gort(sig.canonParamTypes.get(0), imports) : "";
            out.append("func (l *")
                .append(tn)
                .append(") ")
                .append(sig.goName)
                .append(formatExtSigWithParamNames(sig, imports))
                .append(" {\n");
            switch (sig.simpleName) {
              case "size" -> out.append(n == 0 ? "\treturn int32(len(l.strs))\n" : ("\treturn " + n + "\n"));
              case "isEmpty" -> out.append(n == 0 ? "\treturn len(l.strs) == 0\n" : "\treturn false\n");
              case "get" -> out.append("\tCheckElementIndex(arg0, int32(len(l.strs)))\n\treturn l.strs[arg0]\n");
              case "getFirst", "element" ->
                  out.append(
                      n == 0
                          ? "\tif len(l.strs) == 0 { panic(javalang.NewNoSuchElementException(\"empty\")) }\n\treturn l.strs[0]\n"
                          : "\treturn l.strs[0]\n");
              case "getLast" ->
                  out.append(
                      n == 0
                          ? "\tif len(l.strs) == 0 { panic(javalang.NewNoSuchElementException(\"empty\")) }\n\treturn l.strs[len(l.strs)-1]\n"
                          : ("\treturn l.strs[" + (n - 1) + "]\n"));
              case "contains" -> out.append("\treturn immStrSliceContains(").append(sliceExpr).append(", arg0)\n");
              case "indexOf" -> out.append("\treturn immStrSliceIndexOf(").append(sliceExpr).append(", arg0)\n");
              case "lastIndexOf" -> out.append("\treturn immStrSliceLastIndexOf(").append(sliceExpr).append(", arg0)\n");
              case "containsAll" ->
                  out.append("\treturn colContainsAll(func(x any) bool { return immStrSliceContains(")
                      .append(sliceExpr)
                      .append(", x) }, arg0)\n");
              case "iterator" ->
                  out.append(
                      n == 1
                          ? "\treturn newSingleStrIter(l.strs[0])\n"
                          : ("\treturn newStrSliceIter(" + sliceExpr + ")\n"));
              case "listIterator" ->
                  out.append(
                      arity == 0
                          ? "\treturn newImmListIter(strsToAnySlice(" + sliceExpr + "), 0)\n"
                          : "\treturn newImmListIter(strsToAnySlice(" + sliceExpr + "), arg0)\n");
              case "subList" ->
                  out.append(
                      "\tsz := int32(len(l.strs))\n\tCheckPositionIndexes(arg0, arg1, sz)\n\tif arg0 == 0 && arg1 == sz { return l }\n\treturn wrapImmStrSliceCopy(l.strs[arg0:arg1])\n");
              case "asList" -> out.append("\treturn l\n");
              case "reverse", "reversed" -> {
                if (n == 1) {
                  out.append("\treturn l\n");
                } else if (n >= 2) {
                  out.append("\treturn &").append(tn).append("{strs: [").append(n).append("]string{");
                  for (int j = n - 1; j >= 0; j--) {
                    if (j < n - 1) out.append(", ");
                    out.append("l.strs[").append(j).append(']');
                  }
                  out.append("}}\n");
                } else {
                  out.append(
                      "\tif len(l.strs) <= 1 { return l }\n\trev := make([]string, len(l.strs))\n\tfor i, v := range l.strs { rev[len(l.strs)-1-i] = v }\n\treturn &jImmStrList{strs: rev}\n");
                }
              }
              case "toArray" -> {
                if (arity == 0) {
                  out.append("\treturn strsToAnySlice(").append(sliceExpr).append(")\n");
                } else if (p0.startsWith("[]")) {
                  out.append("\treturn sliceToArray1(strsToAnySlice(").append(sliceExpr).append("), arg0)\n");
                } else {
                  out.append("\tCheckNotNull(arg0)\n\treturn sliceToArray1(strsToAnySlice(")
                      .append(sliceExpr)
                      .append("), Cast[[]any](arg0.")
                      .append(intFnApply)
                      .append("(int32(len(l.strs)))))\n");
                }
              }
              case "forEach" ->
                  out.append("\tCheckNotNull(arg0)\n\tfor _, x := range l.strs { arg0.")
                      .append(consumerAccept)
                      .append("(x) }\n");
              case "stream", "parallelStream" ->
                  out.append("\treturn streamFromSlice(strsToAnySlice(").append(sliceExpr).append("))\n");
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
                  "replaceAll" ->
                  out.append("\tpanic(javalang.NewUnsupportedOperationException(\"immutable\"))\n");
              default ->
                  out.append("\tpanic(")
                      .append(GoNames.goStringLiteral("UNIMPLEMENTED_EXT: " + tn + "." + sig.goName))
                      .append(")\n");
            }
            out.append("}\n");
          }
        }
      }
      out.append('\n');
    }

    // 3. Factory functions for fixed 1..8 any & string ImmutableLists + slice wrappers
    out.append("func newImmList1(e0 any) ").append(ilGo).append(" {\n")
        .append("\tif s0, ok := e0.(string); ok { return ImmListOfStr1(s0) }\n")
        .append("\te0 = CheckNotNull(ToAny(e0))\n\treturn &jImmList1{elem: [1]any{e0}}\n}\n")
        .append("func newImmList2(e0, e1 any) ").append(ilGo).append(" {\n")
        .append("\te0 = CheckNotNull(ToAny(e0))\n\te1 = CheckNotNull(ToAny(e1))\n\treturn &jImmList2{elem: [2]any{e0, e1}}\n}\n")
        .append("func newImmList3(e0, e1, e2 any) ").append(ilGo).append(" {\n")
        .append("\te0 = CheckNotNull(ToAny(e0))\n\te1 = CheckNotNull(ToAny(e1))\n\te2 = CheckNotNull(ToAny(e2))\n\treturn &jImmList3{elem: [3]any{e0, e1, e2}}\n}\n")
        .append("func newImmList4(e0, e1, e2, e3 any) ").append(ilGo).append(" {\n")
        .append("\te0 = CheckNotNull(ToAny(e0))\n\te1 = CheckNotNull(ToAny(e1))\n\te2 = CheckNotNull(ToAny(e2))\n\te3 = CheckNotNull(ToAny(e3))\n\treturn &jImmList4{elem: [4]any{e0, e1, e2, e3}}\n}\n")
        .append("func newImmList5(e0, e1, e2, e3, e4 any) ").append(ilGo).append(" {\n")
        .append("\te0 = CheckNotNull(ToAny(e0))\n\te1 = CheckNotNull(ToAny(e1))\n\te2 = CheckNotNull(ToAny(e2))\n\te3 = CheckNotNull(ToAny(e3))\n\te4 = CheckNotNull(ToAny(e4))\n\treturn &jImmList5{elem: [5]any{e0, e1, e2, e3, e4}}\n}\n")
        .append("func newImmList6(e0, e1, e2, e3, e4, e5 any) ").append(ilGo).append(" {\n")
        .append("\te0 = CheckNotNull(ToAny(e0))\n\te1 = CheckNotNull(ToAny(e1))\n\te2 = CheckNotNull(ToAny(e2))\n\te3 = CheckNotNull(ToAny(e3))\n\te4 = CheckNotNull(ToAny(e4))\n\te5 = CheckNotNull(ToAny(e5))\n\treturn &jImmList6{elem: [6]any{e0, e1, e2, e3, e4, e5}}\n}\n")
        .append("func newImmList7(e0, e1, e2, e3, e4, e5, e6 any) ").append(ilGo).append(" {\n")
        .append("\te0 = CheckNotNull(ToAny(e0))\n\te1 = CheckNotNull(ToAny(e1))\n\te2 = CheckNotNull(ToAny(e2))\n\te3 = CheckNotNull(ToAny(e3))\n\te4 = CheckNotNull(ToAny(e4))\n\te5 = CheckNotNull(ToAny(e5))\n\te6 = CheckNotNull(ToAny(e6))\n\treturn &jImmList7{elem: [7]any{e0, e1, e2, e3, e4, e5, e6}}\n}\n")
        .append("func newImmList8(e0, e1, e2, e3, e4, e5, e6, e7 any) ").append(ilGo).append(" {\n")
        .append("\te0 = CheckNotNull(ToAny(e0))\n\te1 = CheckNotNull(ToAny(e1))\n\te2 = CheckNotNull(ToAny(e2))\n\te3 = CheckNotNull(ToAny(e3))\n\te4 = CheckNotNull(ToAny(e4))\n\te5 = CheckNotNull(ToAny(e5))\n\te6 = CheckNotNull(ToAny(e6))\n\te7 = CheckNotNull(ToAny(e7))\n\treturn &jImmList8{elem: [8]any{e0, e1, e2, e3, e4, e5, e6, e7}}\n}\n")
        .append("func ImmListOfStr1(s0 string) ").append(ilGo).append(" {\n")
        .append("\tif s0 == NullString { ThrowNPE() }\n\treturn &jImmStrList1{strs: [1]string{s0}}\n}\n")
        .append("func ImmListOfStr2(s0, s1 string) ").append(ilGo).append(" {\n")
        .append("\tif s0 == NullString || s1 == NullString { ThrowNPE() }\n\treturn &jImmStrList2{strs: [2]string{s0, s1}}\n}\n")
        .append("func ImmListOfStr3(s0, s1, s2 string) ").append(ilGo).append(" {\n")
        .append("\tif s0 == NullString || s1 == NullString || s2 == NullString { ThrowNPE() }\n\treturn &jImmStrList3{strs: [3]string{s0, s1, s2}}\n}\n")
        .append("func ImmListOfStr4(s0, s1, s2, s3 string) ").append(ilGo).append(" {\n")
        .append("\tif s0 == NullString || s1 == NullString || s2 == NullString || s3 == NullString { ThrowNPE() }\n\treturn &jImmStrList4{strs: [4]string{s0, s1, s2, s3}}\n}\n")
        .append("func ImmListOfStr5(s0, s1, s2, s3, s4 string) ").append(ilGo).append(" {\n")
        .append("\tif s0 == NullString || s1 == NullString || s2 == NullString || s3 == NullString || s4 == NullString { ThrowNPE() }\n\treturn &jImmStrList5{strs: [5]string{s0, s1, s2, s3, s4}}\n}\n")
        .append("func ImmListOfStr6(s0, s1, s2, s3, s4, s5 string) ").append(ilGo).append(" {\n")
        .append("\tif s0 == NullString || s1 == NullString || s2 == NullString || s3 == NullString || s4 == NullString || s5 == NullString { ThrowNPE() }\n\treturn &jImmStrList6{strs: [6]string{s0, s1, s2, s3, s4, s5}}\n}\n")
        .append("func ImmListOfStr7(s0, s1, s2, s3, s4, s5, s6 string) ").append(ilGo).append(" {\n")
        .append("\tif s0 == NullString || s1 == NullString || s2 == NullString || s3 == NullString || s4 == NullString || s5 == NullString || s6 == NullString { ThrowNPE() }\n\treturn &jImmStrList7{strs: [7]string{s0, s1, s2, s3, s4, s5, s6}}\n}\n")
        .append("func ImmListOfStr8(s0, s1, s2, s3, s4, s5, s6, s7 string) ").append(ilGo).append(" {\n")
        .append("\tif s0 == NullString || s1 == NullString || s2 == NullString || s3 == NullString || s4 == NullString || s5 == NullString || s6 == NullString || s7 == NullString { ThrowNPE() }\n\treturn &jImmStrList8{strs: [8]string{s0, s1, s2, s3, s4, s5, s6, s7}}\n}\n")
        .append("func wrapImmStrSliceCopy(src []string) ").append(ilGo).append(" {\n")
        .append("\tswitch len(src) {\n\tcase 0:\n\t\treturn emptyImmutableList_\n")
        .append("\tcase 1:\n\t\treturn ImmListOfStr1(src[0])\n")
        .append("\tcase 2:\n\t\treturn ImmListOfStr2(src[0], src[1])\n")
        .append("\tcase 3:\n\t\treturn ImmListOfStr3(src[0], src[1], src[2])\n")
        .append("\tcase 4:\n\t\treturn ImmListOfStr4(src[0], src[1], src[2], src[3])\n")
        .append("\tcase 5:\n\t\treturn ImmListOfStr5(src[0], src[1], src[2], src[3], src[4])\n")
        .append("\tcase 6:\n\t\treturn ImmListOfStr6(src[0], src[1], src[2], src[3], src[4], src[5])\n")
        .append("\tcase 7:\n\t\treturn ImmListOfStr7(src[0], src[1], src[2], src[3], src[4], src[5], src[6])\n")
        .append("\tcase 8:\n\t\treturn ImmListOfStr8(src[0], src[1], src[2], src[3], src[4], src[5], src[6], src[7])\n")
        .append("\tdefault:\n\t\tcp := make([]string, len(src))\n\t\tfor i, s := range src {\n\t\t\tif s == NullString { ThrowNPE() }\n\t\t\tcp[i] = s\n\t\t}\n\t\treturn &jImmStrList{strs: cp}\n\t}\n}\n")
        .append("func wrapImmutableList(items []any) ").append(ilGo).append(" {\n")
        .append("\tswitch len(items) {\n\tcase 0:\n\t\treturn emptyImmutableList_\n")
        .append("\tcase 1:\n\t\tif s0, ok := items[0].(string); ok && s0 != NullString { return &jImmStrList1{strs: [1]string{s0}} }\n\t\treturn &jImmList1{elem: [1]any{ToAny(items[0])}}\n")
        .append("\tcase 2:\n\t\treturn &jImmList2{elem: [2]any{ToAny(items[0]), ToAny(items[1])}}\n")
        .append("\tcase 3:\n\t\treturn &jImmList3{elem: [3]any{ToAny(items[0]), ToAny(items[1]), ToAny(items[2])}}\n")
        .append("\tcase 4:\n\t\treturn &jImmList4{elem: [4]any{ToAny(items[0]), ToAny(items[1]), ToAny(items[2]), ToAny(items[3])}}\n")
        .append("\tcase 5:\n\t\treturn &jImmList5{elem: [5]any{ToAny(items[0]), ToAny(items[1]), ToAny(items[2]), ToAny(items[3]), ToAny(items[4])}}\n")
        .append("\tcase 6:\n\t\treturn &jImmList6{elem: [6]any{ToAny(items[0]), ToAny(items[1]), ToAny(items[2]), ToAny(items[3]), ToAny(items[4]), ToAny(items[5])}}\n")
        .append("\tcase 7:\n\t\treturn &jImmList7{elem: [7]any{ToAny(items[0]), ToAny(items[1]), ToAny(items[2]), ToAny(items[3]), ToAny(items[4]), ToAny(items[5]), ToAny(items[6])}}\n")
        .append("\tcase 8:\n\t\treturn &jImmList8{elem: [8]any{ToAny(items[0]), ToAny(items[1]), ToAny(items[2]), ToAny(items[3]), ToAny(items[4]), ToAny(items[5]), ToAny(items[6]), ToAny(items[7])}}\n")
        .append("\tdefault:\n\t\tcp := make([]any, len(items))\n\t\tfor i, v := range items { cp[i] = ToAny(v) }\n\t\treturn &").append(st).append("{items_: cp[:len(cp):len(cp)]}\n\t}\n}\n")
        .append("func wrapImmutableListOwn(items []any) ").append(ilGo).append(" {\n")
        .append("\tswitch len(items) {\n\tcase 0:\n\t\treturn emptyImmutableList_\n")
        .append("\tcase 1:\n\t\tif s0, ok := items[0].(string); ok && s0 != NullString { return &jImmStrList1{strs: [1]string{s0}} }\n\t\treturn &jImmList1{elem: [1]any{ToAny(items[0])}}\n")
        .append("\tcase 2:\n\t\treturn &jImmList2{elem: [2]any{ToAny(items[0]), ToAny(items[1])}}\n")
        .append("\tcase 3:\n\t\treturn &jImmList3{elem: [3]any{ToAny(items[0]), ToAny(items[1]), ToAny(items[2])}}\n")
        .append("\tcase 4:\n\t\treturn &jImmList4{elem: [4]any{ToAny(items[0]), ToAny(items[1]), ToAny(items[2]), ToAny(items[3])}}\n")
        .append("\tdefault:\n\t\tfor i, v := range items { items[i] = ToAny(v) }\n\t\treturn &").append(st).append("{items_: items[:len(items):len(items)]}\n\t}\n}\n\n");

    // 4. Zero-escape inline-buffer methods on *Ext_com_google_common_collect_ImmutableList_Builder
    TypeElement ilbTe = resolver.elements.getTypeElement("com.google.common.collect.ImmutableList.Builder");
    if (ilbTe != null && resolver.getReferencedExternalTypes().contains(ilbTe)) {
      String bst = goStructType(ilbTe, imports);
      out.append("func InitImmListBuilderCap(b *").append(bst).append(", expected int32) *").append(bst).append(" {\n")
          .append("\tif expected < 0 {\n\t\tpanic(javalang.NewIllegalArgumentException(\"expectedSize cannot be negative but was: \" + fmt.Sprint(expected)))\n\t}\n")
          .append("\tif expected > 8 && expected > b.capHint {\n\t\tb.capHint = expected\n\t}\n\treturn b\n}\n")
          .append("func (b *").append(bst).append(") reserveTotal(need int) {\n")
          .append("\tif need <= 8 {\n\t\treturn\n\t}\n")
          .append("\tif int32(need) > b.capHint {\n\t\tb.capHint = int32(need)\n\t}\n")
          .append("\tif b.mode == 1 && b.n > 8 && (cap(b.strOver) < need || b.aliased) {\n")
          .append("\t\tcp := make([]string, len(b.strOver), max(need, len(b.strOver)*2))\n\t\tcopy(cp, b.strOver)\n\t\tb.strOver = cp\n\t\tb.aliased = false\n\t}\n")
          .append("\tif b.mode == 2 && b.n > 8 && (cap(b.overflow) < need || b.aliased) {\n")
          .append("\t\tcp := make([]any, len(b.overflow), max(need, len(b.overflow)*2))\n\t\tcopy(cp, b.overflow)\n\t\tb.overflow = cp\n\t\tb.aliased = false\n\t}\n}\n")
          .append("func (b *").append(bst).append(") promoteStrToAny() {\n")
          .append("\tif b.n <= 8 {\n\t\tfor i := int32(0); i < b.n; i++ { b.anyBuf[i] = b.strBuf[i] }\n\t\tb.strBuf = [8]string{}\n\t\tb.mode = 2\n\t\treturn\n\t}\n")
          .append("\tcap0 := max(len(b.strOver)+4, int(b.capHint))\n\tb.overflow = make([]any, len(b.strOver), cap0)\n\tfor i, s := range b.strOver { b.overflow[i] = s }\n\tb.strOver = nil\n\tb.aliased = false\n\tb.mode = 2\n}\n")
          .append("func (b *").append(bst).append(") AddStrVoid_(s string) {\n")
          .append("\tif b == nil || s == NullString { ThrowNPE() }\n")
          .append("\tswitch b.mode {\n\tcase 0:\n\t\tb.mode = 1\n\t\tb.strBuf[0] = s\n\t\tb.n = 1\n")
          .append("\tcase 1:\n\t\tif b.n < 8 {\n\t\t\tb.strBuf[b.n] = s\n\t\t\tb.n++\n\t\t\treturn\n\t\t}\n")
          .append("\t\tif b.n == 8 {\n\t\t\tcap0 := max(16, int(b.capHint))\n\t\t\tb.strOver = make([]string, 9, cap0)\n\t\t\tcopy(b.strOver[:8], b.strBuf[:])\n\t\t\tb.strBuf = [8]string{}\n\t\t\tb.strOver[8] = s\n\t\t\tb.n = 9\n\t\t\tb.aliased = false\n\t\t\treturn\n\t\t}\n")
          .append("\t\tif b.aliased {\n\t\t\tcp := make([]string, len(b.strOver), max(len(b.strOver)*2, 16))\n\t\t\tcopy(cp, b.strOver)\n\t\t\tb.strOver = cp\n\t\t\tb.aliased = false\n\t\t}\n")
          .append("\t\tb.strOver = append(b.strOver, s)\n\t\tb.n++\n")
          .append("\tdefault:\n\t\tb.addAnyNonNull(s)\n\t}\n}\n")
          .append("func (b *").append(bst).append(") AddStrSelf_(s string) *").append(bst).append(" {\n\tb.AddStrVoid_(s)\n\treturn b\n}\n")
          .append("func (b *").append(bst).append(") addAnyNonNull(v any) {\n")
          .append("\tif b.mode == 0 {\n\t\tif s, ok := v.(string); ok {\n\t\t\tb.mode = 1\n\t\t\tb.strBuf[0] = s\n\t\t\tb.n = 1\n\t\t\treturn\n\t\t}\n\t\tb.mode = 2\n\t} else if b.mode == 1 {\n")
          .append("\t\tif s, ok := v.(string); ok {\n\t\t\tb.AddStrVoid_(s)\n\t\t\treturn\n\t\t}\n\t\tb.promoteStrToAny()\n\t}\n")
          .append("\tif b.n < 8 {\n\t\tb.anyBuf[b.n] = v\n\t\tb.n++\n\t\treturn\n\t}\n")
          .append("\tif b.n == 8 {\n\t\tcap0 := max(16, int(b.capHint))\n\t\tb.overflow = make([]any, 9, cap0)\n\t\tcopy(b.overflow[:8], b.anyBuf[:])\n\t\tb.anyBuf = [8]any{}\n\t\tb.overflow[8] = v\n\t\tb.n = 9\n\t\tb.aliased = false\n\t\treturn\n\t}\n")
          .append("\tif b.aliased {\n\t\tcp := make([]any, len(b.overflow), max(len(b.overflow)*2, 16))\n\t\tcopy(cp, b.overflow)\n\t\tb.overflow = cp\n\t\tb.aliased = false\n\t}\n")
          .append("\tb.overflow = append(b.overflow, v)\n\tb.n++\n}\n")
          .append("func (b *").append(bst).append(") AddVoid_(v any) {\n\tif b == nil { ThrowNPE() }\n\tb.addAnyNonNull(CheckNotNull(ToAny(v)))\n}\n")
          .append("func (b *").append(bst).append(") AddSelf_(v any) *").append(bst).append(" {\n\tb.AddVoid_(v)\n\treturn b\n}\n")
          .append("func (b *").append(bst).append(") AddAllVoid_(src any) {\n")
          .append("\tif b == nil { ThrowNPE() }\n\tCheckNotNull(src)\n")
          .append("\tif ssi, ok := src.(StrSliceIterable); ok {\n\t\tif ss, ok2 := ssi.IterStrSlice_(); ok2 {\n\t\t\tb.reserveTotal(int(b.n) + len(ss))\n\t\t\tfor _, s := range ss { b.AddStrVoid_(s) }\n\t\t\treturn\n\t\t}\n\t}\n")
          .append("\tif tv, ok := src.(jTypedViewCursorSource); ok {\n\t\tif n, ok2 := tv.typedViewLen_(); ok2 {\n\t\t\tb.reserveTotal(int(b.n) + n)\n\t\t\tfor i := 0; i < n; i++ { b.addAnyNonNull(CheckNotNull(ToAny(tv.typedViewAt_(i)))) }\n\t\t\treturn\n\t\t}\n\t}\n")
          .append("\tif si, ok := src.(SliceIterable); ok {\n\t\titems := si.IterSlice_()\n\t\tb.reserveTotal(int(b.n) + len(items))\n\t\tfor _, x := range items { b.addAnyNonNull(CheckNotNull(ToAny(x))) }\n\t\treturn\n\t}\n")
          .append("\tif col, ok := src.(Ext_java_util_Collection); ok {\n\t\tif sz := int(col.Size()); sz > 0 {\n\t\t\tb.reserveTotal(int(b.n) + sz)\n\t\t}\n\t}\n")
          .append("\tfor it := NewIterCursor(src); it.HasNext(); {\n\t\tb.addAnyNonNull(CheckNotNull(ToAny(it.Next())))\n\t}\n}\n")
          .append("func (b *").append(bst).append(") AddAllSelf_(src any) *").append(bst).append(" {\n\tb.AddAllVoid_(src)\n\treturn b\n}\n")
          .append("func (b *").append(bst).append(") BuildImmList_() ").append(ilGo).append(" {\n")
          .append("\tif b == nil { ThrowNPE() }\n\tif b.n == 0 { return emptyImmutableList_ }\n")
          .append("\tif b.mode == 1 {\n\t\tswitch b.n {\n")
          .append("\t\tcase 1: return &jImmStrList1{strs: [1]string{b.strBuf[0]}}\n")
          .append("\t\tcase 2: return &jImmStrList2{strs: [2]string{b.strBuf[0], b.strBuf[1]}}\n")
          .append("\t\tcase 3: return &jImmStrList3{strs: [3]string{b.strBuf[0], b.strBuf[1], b.strBuf[2]}}\n")
          .append("\t\tcase 4: return &jImmStrList4{strs: [4]string{b.strBuf[0], b.strBuf[1], b.strBuf[2], b.strBuf[3]}}\n")
          .append("\t\tcase 5: return &jImmStrList5{strs: [5]string{b.strBuf[0], b.strBuf[1], b.strBuf[2], b.strBuf[3], b.strBuf[4]}}\n")
          .append("\t\tcase 6: return &jImmStrList6{strs: [6]string{b.strBuf[0], b.strBuf[1], b.strBuf[2], b.strBuf[3], b.strBuf[4], b.strBuf[5]}}\n")
          .append("\t\tcase 7: return &jImmStrList7{strs: [7]string{b.strBuf[0], b.strBuf[1], b.strBuf[2], b.strBuf[3], b.strBuf[4], b.strBuf[5], b.strBuf[6]}}\n")
          .append("\t\tcase 8: return &jImmStrList8{strs: b.strBuf}\n")
          .append("\t\tdefault:\n\t\t\tb.aliased = true\n\t\t\ts := b.strOver[:b.n:b.n]\n\t\t\treturn &jImmStrList{strs: s}\n\t\t}\n\t}\n")
          .append("\tswitch b.n {\n")
          .append("\tcase 1: return &jImmList1{elem: [1]any{b.anyBuf[0]}}\n")
          .append("\tcase 2: return &jImmList2{elem: [2]any{b.anyBuf[0], b.anyBuf[1]}}\n")
          .append("\tcase 3: return &jImmList3{elem: [3]any{b.anyBuf[0], b.anyBuf[1], b.anyBuf[2]}}\n")
          .append("\tcase 4: return &jImmList4{elem: [4]any{b.anyBuf[0], b.anyBuf[1], b.anyBuf[2], b.anyBuf[3]}}\n")
          .append("\tcase 5: return &jImmList5{elem: [5]any{b.anyBuf[0], b.anyBuf[1], b.anyBuf[2], b.anyBuf[3], b.anyBuf[4]}}\n")
          .append("\tcase 6: return &jImmList6{elem: [6]any{b.anyBuf[0], b.anyBuf[1], b.anyBuf[2], b.anyBuf[3], b.anyBuf[4], b.anyBuf[5]}}\n")
          .append("\tcase 7: return &jImmList7{elem: [7]any{b.anyBuf[0], b.anyBuf[1], b.anyBuf[2], b.anyBuf[3], b.anyBuf[4], b.anyBuf[5], b.anyBuf[6]}}\n")
          .append("\tcase 8: return &jImmList8{elem: b.anyBuf}\n")
          .append("\tdefault:\n\t\tb.aliased = true\n\t\ts := b.overflow[:b.n:b.n]\n\t\treturn &").append(st).append("{items_: s}\n\t}\n}\n\n");
    }
  }
}
