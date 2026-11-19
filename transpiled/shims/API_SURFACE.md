# `transpiled/shims/` API Surface Reference

Module import prefix: `github.com/paullepoulpe/jsgo/transpiled/shims/<pkg>`

All 12 shim packages are clean-room (`PORTING.md` compliant), pure Go (zero `"unsafe"` in production `.go` files), race-free, and deterministic (zero unannotated Go `map` `range` iteration; verified by `go run ./tools/maprangecheck ./transpiled/shims/...`).

---

## 1. Java -> Go Package Mapping

| Java Package(s) | Go Shim Package | Purpose |
|---|---|---|
| `com.google.javascript.rhino.Node` arena storage | `shims/nodestore` | Handle-based 40-byte SoA AST `Node` arena (`NodeStore`, `Node`, `StringID`, `HandleTable[P]`) |
| `java.lang.*`, `java.util.Objects` | `shims/javalang` | WTF-8/UTF-16 `String`, `Character`, `Double`/`Integer`/`Long` math & formatting, `Object`, `Class`, `System`, `Throwable` hierarchy |
| `java.util.*` | `shims/javautil` | Deterministic `LinkedHashMap`/`HashMap`/`IdentityHashMap`, `LinkedHashSet`/`HashSet`, `TreeMap`, `TreeSet`, `ArrayList`, `ArrayDeque`, `PriorityQueue`, `BitSet`, `Optional`, `Random`, `Comparator`, `Collections` |
| `java.io.*` | `shims/javaio` | `InputStream`/`OutputStream`, `Reader`/`Writer` (`uint16` UTF-16 buffers), `File`, `PrintStream`, `IOException` hierarchy |
| `java.nio.charset.*`, `java.nio.file.*`, `java.net.URI` | `shims/javanio` | `Charset`, `StandardCharsets`, `CharsetEncoder`, `Path`, `Files`, `FileSystems`, `URI` |
| `java.util.zip.*` | `shims/javazip` | `ZipInputStream`, `ZipFile`, `ZipEntry`, `GZIPInputStream`, `GZIPOutputStream`, `Deflater` constants |
| `java.math.BigInteger` | `shims/javamath` | Immutable arbitrary-precision `BigInteger` |
| `java.util.regex.*`, `com.google.re2j.*` | `shims/javaregex` | Java-semantics UTF-16 `Pattern`, `Matcher`, `PatternSyntaxException`, regex `String` helpers |
| `com.google.common.{base,collect,primitives,io,hash,escape,net}.*` | `shims/guava` | `Preconditions`, `ImmutableList`/`Set`/`Map` (+ builders), `Lists`/`Maps`/`Sets`/`Iterables`/`Chars`, `Strings`, `Ascii`, `Joiner`, `Splitter`, `CharMatcher`, `MoreObjects`, `PercentEscaper`, `BaseEncoding`, `ByteStreams`/`CharStreams`/`Files`, `Hashing` (`farmHashFingerprint64`, `murmur3_128`) |
| `com.google.gson.*` | `shims/gson` | Streaming `JsonReader`/`JsonWriter`, DOM `JsonElement`/`JsonObject`/`JsonArray`/`JsonPrimitive`/`JsonNull`, `Gson`/`GsonBuilder` |
| `com.google.protobuf.*`, Closure `.proto` messages | `shims/protobuf` | `ByteString`, `CodedInputStream`, `WireFormat`, `ExtensionRegistry`, `TextFormat`, and all Closure v20261104 `.proto` messages + `*_NewBuilder()` builders |
| `Class.getResource*`, `ResourceLoader`, `java.util.ResourceBundle` | `shims/resources` | Embedded `externs.zip`, `runtime_libs.typedast`, `com/google/javascript/jscomp/js/**`, `ParserConfig.properties`, `lib/*.js`, `GetResourceAsStream`, `LoadTextResource`, `ResourceBundle_GetBundle` |

---

## 2. Exception & Error Conventions

- **Unchecked Java exceptions (`RuntimeException`, `Error`)**: Thrown via `panic(ex)` where `ex` implements `javalang.Throwable` (`IRuntimeException`, `INullPointerException`, `IIllegalArgumentException`, `IIllegalStateException`, `IIndexOutOfBoundsException`, `IClassCastException`, `INoSuchElementException`, `IUnsupportedOperationException`, `IAssertionError`, etc.). Recover in translated `try/catch` blocks via `ex := javalang.RecoverException(recover())` and type-assert against the marker interface (e.g. `if e, ok := ex.(javalang.IIllegalArgumentException); ok { ... }`).
- **Checked Java I/O & parse exceptions (`IOException`, `InvalidProtocolBufferException`, `URISyntaxException`, `ZipException`)**: Returned as `(T, error)` where the concrete error (`*javaio.IOException`, `*javaio.FileNotFoundException`, `*javaio.EOFException`, `*javazip.ZipException`, `*protobuf.InvalidProtocolBufferException`) ALSO implements `javalang.Throwable` so it can either be checked via `err != nil` or panicked/caught uniformly.

---

## 3. Package-by-Package Exported Symbols

### 3.1 `shims/nodestore`
- **Core types**: `NodeStore`, `Node` (`uint32` handle, `0` = null), `NodeData`, `NodeAux`, `StringID`, `HandleTable[P]`, `NodeAllocator`, `BaselineSnapshot`.
- **Constructors & lifecycle**: `NewNodeStore() *NodeStore`, `ResetNodeStore()`, `LockStandaloneStore()`, `UnlockStandaloneStore()`, `Epoch() uint64`, `RewindBaseline(BaselineSnapshot)`.
- **Source position sentinel**: Unknown line/column returns `-1` (`GetLineno()`, `GetCharno()`, `GetSourceOffset()`), matching `com.google.javascript.rhino.Node`.

### 3.2 `shims/javalang`
- **UTF-16 / WTF-8 `String` (`jstring.go`)**:
  - `StringLength(s string) int`, `StringCharAt(s string, index int) uint16`
  - `StringSubstring(s string, begin, end int) string`, `StringSubstringFrom(s string, begin int) string` (preserves split surrogate halves in WTF-8)
  - `StringIndexOf(s, sub string) int`, `StringIndexOfFrom(s, sub string, fromIndex int) int`, `StringIndexOfChar(s string, ch uint16) int`
  - `StringLastIndexOf(s, sub string) int`, `StringLastIndexOfFrom(s, sub string, fromIndex int) int`, `StringLastIndexOfCharFrom(s string, ch uint16, fromIndex int) int`
  - `StringStartsWith(s, prefix string) bool`, `StringEndsWith(s, suffix string) bool`, `StringRegionMatches(s string, toffset int, other string, ooffset, length int) bool`
  - `StringReplace(s, target, replacement string) string`, `StringTrim(s string) string`
  - `StringSplitLiteral(s, sep string) []string`, `StringSplitLiteralLimit(s, sep string, limit int) []string`
  - `StringCompareTo(a, b string) int`, `StringHashCode(s string) int32`
  - `StringToLowerCaseASCII(s string) string`, `StringToUpperCaseASCII(s string) string`
  - `StringToUTF16(s string) []uint16`, `StringFromUTF16(units []uint16) string`, `StringFromChar(c uint16) string`
  - `IsASCII(s string) bool`, `ASCIIPrefixLen(s string) int`, `IsCanonicalWTF8(s string) bool`
- **`Character` (`character.go`)**:
  - Constants: `CharacterMinRadix` (`2`), `CharacterMaxRadix` (`36`)
  - `CharacterIsDigit`, `CharacterIsLetter`, `CharacterIsLetterOrDigit`, `CharacterIsUpperCase`, `CharacterIsLowerCase`, `CharacterIsWhitespace` (`uint16 -> bool`), `CharacterIsWhitespaceCodePoint(int) bool`
  - `CharacterIsJavaIdentifierStart`, `CharacterIsJavaIdentifierPart`, `CharacterIsIdentifierIgnorable` (`uint16 -> bool`)
  - `CharacterIsSupplementaryCodePoint(int) bool`, `CharacterToLowerCase(uint16) uint16`, `CharacterToUpperCase(uint16) uint16`, `CharacterDigit(ch uint16, radix int) int`
- **Primitives, Narrowing, Shifts & `Double.toString` (`primitives.go`)**:
  - `D2I(float64) int32`, `D2L(float64) int64`, `DRem(a, b float64) float64`
  - `Shl32`, `Shr32`, `Ushr32` (`(int32, int) int32`); `Shl64`, `Shr64`, `Ushr64` (`(int64, int) int64`)
  - `DoubleToString(float64) string` (Java 11+ Schubfach exact formatting), `DoubleParseDouble(string) float64`, `DoubleCompare(a, b float64) int`, `DoubleIsNaN`, `DoubleIsInfinite`, `DoubleIsFinite`
  - `IntegerParseInt(string) int32`, `IntegerParseIntRadix(string, int) int32`, `IntegerToHexString(int32) string`
  - `LongParseLong(string) int64`, `LongParseLongRadix(string, int) int64`, `LongToHexString(int64) string`
- **`Object`, `Class`, `System` (`object.go`, `class.go`)**:
  - `Equaler` (`Equals(any) bool`), `Hashable` (`HashCode() int32`), `ObjectsEquals(a, b any) bool`, `ObjectHashCode(v any) int32`, `IdentityHashCode(v any) int32`, `ToJavaString(v any) string`, `IsNil(v any) bool`
  - `Class` (`ForName`, `NewClass`, `GetName`, `GetSimpleName`)
  - `System_Arraycopy[T any](src []T, srcPos int, dest []T, destPos, length int)`, `System_CurrentTimeMillis() int64`, `System_NanoTime() int64`
- **`Throwable` Hierarchy (`throwable.go`)**:
  - `Throwable`, `ThrowableBase`, `InitThrowable`, `RecoverException(any) Throwable`
  - Constructors: `NewException`, `NewExceptionMsg`, `NewExceptionCause`, `NewRuntimeException`, `NewRuntimeExceptionMsg`, `NewRuntimeExceptionCause`, `NewIllegalArgumentException`, `NewIllegalStateException`, `NewNullPointerException`, `NewNullPointerExceptionMsg`, `NewIndexOutOfBoundsException`, `NewIndexOutOfBoundsExceptionForIndex`, `NewArrayIndexOutOfBoundsException`, `NewArrayIndexOutOfBoundsExceptionForIndex`, `NewStringIndexOutOfBoundsException`, `NewStringIndexOutOfBoundsExceptionForIndex`, `NewStringIndexOutOfBoundsExceptionForRange`, `NewUnsupportedOperationException`, `NewClassCastException`, `NewNoSuchElementException`, `NewArithmeticException`, `NewNumberFormatException`, `NewError`, `NewAssertionError`.

### 3.3 `shims/javautil`
- **Maps & Sets**:
  - `LinkedHashMap[K comparable, V any]` (and aliases `HashMap`, `IdentityHashMap`): `NewLinkedHashMap`, `NewLinkedHashMapWithCapacity`, `NewHashMap`, `NewIdentityHashMap`, `Put`, `PutVal`, `PutIfAbsent`, `Get`, `GetVal`, `GetOrDefault`, `Remove`, `ContainsKey`, `ContainsValue`, `Size`, `IsEmpty`, `Clear`, `PutAll`, `All() iter.Seq2[K, V]`, `Keys() iter.Seq[K]`, `Values() iter.Seq[V]`, `EntrySet() []MapEntry[K, V]`.
  - `LinkedHashSet[T comparable]` (and alias `HashSet`): `NewLinkedHashSet`, `NewLinkedHashSetWithCapacity`, `NewHashSet`, `Add`, `AddAll`, `Remove`, `Contains`, `Size`, `IsEmpty`, `Clear`, `ToSlice`, `All() iter.Seq[T]`.
  - `TreeMap[K, V]`: `NewTreeMap[K cmp.Ordered, V]`, `NewTreeMapComparator[K, V](Comparator[K])`, `Put`, `Get`, `GetVal`, `Remove`, `ContainsKey`, `FirstKey`, `LastKey`, `FloorKey`, `CeilingKey`, `Size`, `IsEmpty`, `Clear`, `All`, `Keys`, `Values`.
  - `TreeSet[T]`: `NewTreeSet[T cmp.Ordered]`, `NewTreeSetComparator[T](Comparator[T])`, `Add`, `Remove`, `Contains`, `First`, `Last`, `Floor`, `Ceiling`, `PollFirst`, `Size`, `IsEmpty`, `Clear`, `ToSlice`, `All`.
- **Lists, Queues, BitSet, Optional, Random**:
  - `ArrayList[T]`: `NewArrayList`, `NewArrayListWithCapacity`, `NewArrayListFrom`, `Add`, `AddAt`, `AddAll`, `Get`, `Set`, `RemoveAt`, `RemoveElem`, `IndexOf`, `Contains`, `Size`, `IsEmpty`, `Clear`, `Sort`, `SubList`, `ToSlice`, `All`.
  - `ArrayDeque[T]` (O(1) ring buffer): `NewArrayDeque`, `NewArrayDequeWithCapacity`, `AddFirst`, `AddLast`, `Add`, `OfferFirst`, `OfferLast`, `RemoveFirst`, `RemoveLast`, `PollFirst`, `PollLast`, `GetFirst`, `GetLast`, `PeekFirst`, `PeekLast`, `Push`, `Pop`, `RemoveFirstOccurrence`, `Size`, `IsEmpty`, `Clear`, `ToSlice`, `All`.
  - `PriorityQueue[T]`: `NewPriorityQueue`, `NewPriorityQueueComparator`, `Offer`, `Add`, `Peek`, `Poll`, `Remove`, `Size`, `IsEmpty`, `Clear`, `All`.
  - `BitSet`: `NewBitSet`, `NewBitSetWithSize`, `BitSet_ValueOf([]int64)`, `Set`, `SetTo`, `SetRange`, `ClearBit`, `ClearRange`, `ClearAll`, `Flip`, `FlipRange`, `Get`, `NextSetBit`, `NextClearBit`, `PreviousSetBit`, `Length`, `Cardinality`, `IsEmpty`, `Intersects`, `And`, `Or`, `Xor`, `AndNot`, `ToLongArray`, `Clone`, `Equals`, `HashCode`, `String`.
  - `Optional[T]`: `Optional_Empty`, `Optional_Of`, `Optional_OfNullable`, `IsPresent`, `IsEmpty`, `Get`, `OrElse`, `OrElseGet`, `IfPresent`, `Map`.
  - `Random`: `NewRandom`, `NewRandomSeed(int64)`, `SetSeed`, `NextInt`, `NextIntBound(int32)`, `NextLong`, `NextBoolean`.
  - `Comparator[T]` & `Collections_*`: `NaturalOrder`, `ReverseOrder`, `Comparing`, `NullsFirst`, `NullsLast`; `Collections_Sort`, `Collections_Reverse`, `Collections_Shuffle`, `Collections_BinarySearch`, `Collections_Disjoint`, `Collections_Min`, `Collections_Max`, `Collections_EmptyList/Set/Map`, `Collections_Singleton/List/Map`, `Collections_UnmodifiableList/Set/Map`.

### 3.4 `shims/javaio`
- **Byte streams**: `InputStream`, `OutputStream`, `ReaderInputStream` (`NewReaderInputStream`), `ByteArrayInputStream` (`NewByteArrayInputStream`), `ByteArrayOutputStream` (`NewByteArrayOutputStream`, `ToByteArray`, `ToString`, `Size`, `Reset`), `FileInputStream` (`NewFileInputStream`, `NewFileInputStreamFromFile`), `FileOutputStream` (`NewFileOutputStream`, `NewFileOutputStreamFromFile`, `NewFileOutputStreamFromOSFile`), `BufferedInputStream` (`NewBufferedInputStream`), `BufferedOutputStream` (`NewBufferedOutputStream`), `ReadAllBytes(io.Reader) ([]byte, error)`.
- **Character streams (`uint16` UTF-16)**: `Reader` (`Read(cbuf []uint16, off, length int) int`, `Close() error`), `Writer` (`Write(string)`, `AppendCharSequence(string)`, `AppendCharacter(uint16)`, `Flush()`, `Close()`), `StringReader` (`NewStringReader`), `StringWriter` (`NewStringWriter`, `ToString`), `InputStreamReader` (`NewInputStreamReader`, `NewInputStreamReaderCharset`, `ReadAllString`), `OutputStreamWriter` (`NewOutputStreamWriter`, `NewOutputStreamWriterCharset`), `BufferedReader` (`NewBufferedReader`, `ReadLine() (string, error)`), `BufferedWriter` (`NewBufferedWriter`, `NewLine()`).
- **`File` & `PrintStream` & Exceptions**:
  - `File`: `NewFile(string)`, `NewFileParentChild(parent, child string)`, `CreateTempFile(prefix, suffix string)`, `GetPath`, `GetName`, `GetParent() string`, `Parent() (string, bool)`, `GetParentFile`, `GetAbsolutePath`, `GetCanonicalPath`, `Exists`, `IsDirectory`, `IsFile`, `Mkdir`, `Mkdirs`, `Delete`, `ListFiles`.
  - `PrintStream`: `System_Out`, `System_Err`, `NewPrintStream`, `Print`, `Println`, `Printf`, `Flush`, `Close`.
  - Exceptions: `IOException` (`NewIOException`, `NewIOExceptionCause`), `FileNotFoundException` (`NewFileNotFoundException`), `EOFException` (`NewEOFException`).

### 3.5 `shims/javanio`
- **`Charset` & `CharsetEncoder`**: `StandardCharsets_UTF_8`, `StandardCharsets_US_ASCII`, `StandardCharsets_ISO_8859_1`, `StandardCharsets_UTF_16`, `StandardCharsets_UTF_16BE`, `StandardCharsets_UTF_16LE`; `Charset_ForName`, `Charset_IsSupported`, `Charset_DefaultCharset`; `(*Charset).Encode`, `Decode`, `CanEncode`, `CanEncodeChar`, `NewEncoder`; `NewInputStreamReaderWithCharset`, `NewOutputStreamWriterWithCharset`, `NewOutputStreamWriterWithEncoder`, `StringGetBytes`, `NewStringFromBytes`.
- **`Path`, `Files`, `FileSystems`**: `Path_Of(first string, more ...string) *Path` (alias `Paths_Get`), `(*Path).Resolve`, `ResolveSibling`, `Relativize`, `Normalize`, `ToAbsolutePath`, `ToRealPath`, `GetParent`, `GetFileName`, `GetRoot`, `StartsWith`, `EndsWith`, `ToFile`, `ToUri`; `FileSystems_GetDefault() *FileSystem` (`GetPath`, `GetPathMatcher("glob:..." | "regex:...")`); `Files_ReadString`, `Files_ReadStringCharset`, `Files_ReadAllBytes`, `Files_WriteBytes`, `Files_WriteString`, `Files_NewInputStream`, `Files_NewOutputStream`, `Files_NewBufferedReader`, `Files_NewBufferedWriter`, `Files_CreateDirectories`, `Files_CreateTempDirectory`, `Files_CreateTempFile`, `Files_WalkFileTree`, `Files_Exists`, `Files_IsDirectory`, `Files_IsRegularFile`, `Files_Delete`, `Files_DeleteIfExists`.
- **`URI`**: `NewURI(string) (*URI, error)`, `URI_Create(string) *URI`, `Resolve`, `Relativize`, `Normalize`, `GetScheme`, `GetAuthority`, `GetPath`, `GetQuery`, `GetFragment`, `IsAbsolute`, `ToPath`.

### 3.6 `shims/javazip`
- `Deflater_NO_COMPRESSION`, `Deflater_BEST_SPEED`, `Deflater_BEST_COMPRESSION`, `Deflater_DEFAULT_COMPRESSION`.
- `ZipEntry`: `NewZipEntry(string)`, `GetName`, `IsDirectory`, `GetSize`.
- `ZipInputStream`: `NewZipInputStream(io.Reader)`, `NewZipInputStreamCharset(io.Reader, *javanio.Charset)`, `GetNextEntry() (*ZipEntry, error)`, `CloseEntry() error`, `Read`, `Close`.
- `ZipFile`: `NewZipFile(string) (*ZipFile, error)`, `NewZipFileFromFile(*javaio.File)`, `GetEntry(string) *ZipEntry`, `GetInputStream(*ZipEntry) (javaio.InputStream, error)`, `Entries() []*ZipEntry`, `Size() int`, `Close() error`.
- `GZIPInputStream` / `GZIPOutputStream`: `NewGZIPInputStream`, `NewGZIPInputStreamWithSize`, `NewGZIPOutputStream`, `NewGZIPOutputStreamLevel`, `NewFastGzipOutputStream`.

### 3.7 `shims/javamath`
- `BigInteger`: `BigInteger_ZERO`, `BigInteger_ONE`, `BigInteger_TWO`, `BigInteger_TEN`, `NewBigInteger(string)`, `NewBigIntegerRadix(string, int)`, `NewBigIntegerBytes([]byte)`, `NewBigIntegerSignMagnitude(int, []byte)`, `BigInteger_ValueOf(int64)`, `Add`, `Subtract`, `Multiply`, `Divide`, `Remainder`, `DivideAndRemainder`, `Pow`, `Abs`, `Negate`, `Signum`, `Gcd`, `Mod`, `ShiftLeft`, `ShiftRight`, `And`, `Or`, `Xor`, `Not`, `AndNot`, `TestBit`, `SetBit`, `ClearBit`, `BitLength`, `BitCount`, `GetLowestSetBit`, `CompareTo`, `Equals`, `HashCode`, `Min`, `Max`, `IntValue`, `IntValueExact`, `LongValue`, `LongValueExact`, `FloatValue`, `DoubleValue`, `ToByteArray`, `String`, `ToStringRadix`.

### 3.8 `shims/javaregex`
- Constants: `Pattern_UNIX_LINES`, `Pattern_CASE_INSENSITIVE`, `Pattern_COMMENTS`, `Pattern_MULTILINE`, `Pattern_LITERAL`, `Pattern_DOTALL`, `Pattern_UNICODE_CASE`.
- `Pattern`: `Pattern_Compile(string) *Pattern`, `Pattern_CompileFlags(string, int) *Pattern`, `Pattern_Quote(string) string`, `Pattern_Matches(regex, input string) bool`, `IsValidPattern(string) bool`, `(*Pattern).Matcher(string) *Matcher`, `Split(string) []string`, `SplitLimit(string, int) []string`, `Pattern() string`, `Flags() int`.
- `Matcher`: `Matches() bool`, `LookingAt() bool`, `Find() bool`, `FindFrom(int) bool`, `Reset() *Matcher`, `ResetInput(string) *Matcher`, `Region(start, end int) *Matcher`, `Start() int`, `StartGroup(int) int`, `End() int`, `EndGroup(int) int`, `Group() string`, `GroupIndex(int) (string, bool)`, `GroupVal(int) string`, `GroupByName(string) (string, bool)`, `GroupCount() int`, `ReplaceAll(string) string`, `ReplaceFirst(string) string`, `AppendReplacement(*strings.Builder, string) *Matcher`, `AppendTail(*strings.Builder) *strings.Builder`, `Matcher_QuoteReplacement(string) string`.
- `String` regex helpers: `StringMatches`, `StringReplaceAll`, `StringReplaceFirst`, `StringSplit`, `StringSplitLimit`.

### 3.9 `shims/guava`
- **Preconditions**: `CheckState`, `CheckStateMsg`, `CheckStateMsg1..3`, `CheckArgument`, `CheckArgumentMsg`, `CheckArgumentMsg1..3`, `CheckNotNull`, `CheckNotNullMsg`, `CheckNotNullMsg1..2`, `CheckElementIndex`, `CheckPositionIndex`, `CheckPositionIndexes`, `LenientFormat`.
- **Immutable & mutable collection helpers**: `ImmutableList[T]` (`ImmutableList_Of`, `ImmutableList_CopyOf`, `ImmutableList_NewBuilder`), `ImmutableSet[T]` (`ImmutableSet_Of`, `ImmutableSet_CopyOf`, `ImmutableSet_NewBuilder`), `ImmutableMap[K, V]` (`ImmutableMap_Of`, `ImmutableMap_Of1`, `ImmutableMap_Of2`, `ImmutableMap_CopyOf`, `ImmutableMap_NewBuilder` with `Put`, `PutAll`, `BuildOrThrow`, `BuildKeepingLast`), `Lists_*`, `Maps_*`, `Sets_*`, `Iterables_GetOnlyElement`, `Iterables_GetOnlyElementSeq`, `Iterables_Concat`, `Chars_AsList`, `Chars_ToArray`, `Chars_Contains`.
- **Strings, Ascii, Joiner, Splitter, CharMatcher, MoreObjects, Escaper**: `Strings_IsNullOrEmpty`, `Strings_NullToEmpty`, `Strings_EmptyToNull`, `Strings_Repeat`, `Strings_PadStart`, `Strings_PadEnd`, `Strings_CommonPrefix`, `Strings_CommonSuffix`, `Strings_LenientFormat`; `Ascii_ToLowerCase`, `Ascii_ToUpperCase`, `Ascii_ToLowerCaseChar`, `Ascii_ToUpperCaseChar`, `Ascii_IsLowerCase`, `Ascii_IsUpperCase`, `Ascii_EqualsIgnoreCase`; `Joiner_On`, `Joiner_OnChar`; `Splitter_On`, `Splitter_OnChar`, `Splitter_OnMatcher`, `Splitter_OnPattern`; `CharMatcher_Whitespace`, `CharMatcher_Ascii`, `CharMatcher_Any`, `CharMatcher_None`, `CharMatcher_Is`, `CharMatcher_IsNot`, `CharMatcher_InRange`, `CharMatcher_AnyOf`; `MoreObjects_FirstNonNull`, `MoreObjects_ToStringHelper`; `Escaper`, `NewPercentEscaper`.
- **I/O & Hashing**: `BaseEncoding_Base64()`, `BaseEncoding_Base64Url()`; `ByteStreams_ToByteArray`, `ByteStreams_Copy`; `CharStreams_ToString`, `CharStreams_ReadLines`; `CharSource_Wrap`, `Files_AsCharSource`, `Files_ToByteArray`, `Files_WriteBytes`, `Files_AsCharSinkWrite`; `Hashing_FarmHashFingerprint64()`, `Hashing_Murmur3_128()`, `Hashing_Murmur3_128WithSeed(int32)`, `FarmHashFingerprint64([]byte) uint64`, `Murmur3_128([]byte, uint32) (uint64, uint64)`.

### 3.10 `shims/gson`
- **Streaming**: `JsonReader` (`NewJsonReader`, `SetLenient`, `BeginArray`, `EndArray`, `BeginObject`, `EndObject`, `HasNext`, `Peek`, `NextName`, `NextString`, `NextBoolean`, `NextNull`, `NextDouble`, `NextLong`, `NextInt`, `SkipValue`, `Close`), `JsonWriter` (`NewJsonWriter`, `SetIndent`, `SetHtmlSafe`, `SetSerializeNulls`, `SetLenient`, `BeginArray`, `EndArray`, `BeginObject`, `EndObject`, `Name`, `ValueString`, `ValueBool`, `ValueInt64`, `ValueFloat64`, `ValueNumber`, `NullValue`, `JsonValue`, `Flush`, `Close`), `JsonToken`.
- **DOM & `Gson`**: `JsonElement`, `JsonObject` (`NewJsonObject`, `Add`, `AddPropertyString`, `AddPropertyNumber`, `AddPropertyBool`, `Get`, `GetAsJsonObject`, `GetAsJsonArray`, `GetAsJsonPrimitive`, `Has`, `EntrySet`), `JsonArray` (`NewJsonArray`, `Add`, `AddString`, `Size`, `Get`, `All`), `JsonPrimitive`, `JsonNull_INSTANCE`; `Gson`, `NewGson()`, `NewGsonBuilder()` (`SetPrettyPrinting`, `DisableHtmlEscaping`, `SerializeNulls`, `RegisterTypeAdapter`, `Create`), `JsonParser_ParseString`, `JsonParser_ParseReader`, `Streams_Write`.

### 3.11 `shims/protobuf`
- **Runtime**: `ByteString` (`ByteString_EMPTY`, `ByteString_CopyFrom`, `ByteString_CopyFromUtf8`, `ByteString_NewOutput`, `Size`, `ByteAt`, `ToByteArray`, `ToStringUtf8`, `NewCodedInput`), `CodedInputStream` (`CodedInputStream_NewInstance`, `CodedInputStream_NewInstanceBytes`, `IsAtEnd`, `ReadTag`, `ReadMessage`, `ResetSizeCounter`), `ExtensionRegistry_GetEmptyRegistry()`, `WireFormat_GetTagFieldNumber`, `WireFormat_GetTagWireType`, `InvalidProtocolBufferException`, `TextFormat_Merge`, `TextFormat_PrinterInstance()`.
- **Closure `.proto` messages & builders**:
  - TypedAST & Serialization: `TypedAst` (`TypedAst_NewBuilder`, `Clear`, `Build`, `ToByteArray`), `TypedAst_List` (`TypedAst_List_NewBuilder`, `WriteTo`, `WriteToStream`), `LazyAst` (`LazyAst_NewBuilder`), `AstNode` (`AstNode_NewBuilder`, `AstNode_ParseFrom`, `AstNode_ParseFromCodedInput`, `NodeKind`, `NodeProperty`), `TemplateStringValue` (`TemplateStringValue_NewBuilder`), `OptimizationJsdoc` (`OptimizationJsdoc_NewBuilder`, `JsdocTag`), `StringPoolProto` (`StringPoolProto_NewBuilder`), `ExternsSummary` (`ExternsSummary_NewBuilder`), `SourceFilePool` (`SourceFilePool_NewBuilder`), `SourceFileProto` (`SourceFileProto_NewBuilder`, `SourceFileProto_FileOnDisk_NewBuilder`, `SourceFileProto_ZipEntryOnDisk_NewBuilder`), `TypePool` (`TypePool_NewBuilder`), `TypeProto` (`TypeProto_NewBuilder`), `ObjectTypeProto` (`ObjectTypeProto_NewBuilder`), `UnionTypeProto` (`UnionTypeProto_NewBuilder`), `SubtypingEdge` (`SubtypingEdge_NewBuilder`), `PrimitiveType`.
  - Conformance: `ConformanceConfig` (`ConformanceConfig_NewBuilder`), `Requirement` (`Requirement_NewBuilder`, `ToBuilder`), `RequirementScopeEntry`, `ConformanceViolationsProto`, `ViolationProto`.
  - SourceMap: `Mapping_OriginalMapping` (`Mapping_OriginalMapping_NewBuilder`, `Mapping_OriginalMapping_Precision`).
  - Instrumentation: `ReportProfile` (`ReportProfile_NewBuilder`, `ReportProfile_ParseFrom`), `FileProfile` (`FileProfile_NewBuilder`), `InstrumentationPoint` (`InstrumentationPoint_NewBuilder`, `InstrumentationPoint_ParseFrom`), `InstrumentationPointStats` (`InstrumentationPointStats_NewBuilder`).
  - Compiler state: `JSCompilerStateProto` (`JSCompilerStateProto_NewBuilder`, `JSCompilerStateProto_ParseFrom`).

### 3.12 `shims/resources`
- `GetResourceAsStream(basePackage, resourcePath string) javaio.InputStream` (`nil` if missing; supports leading `'/'` root-relative or package-relative paths)
- `GetResourceBytes(basePackage, resourcePath string) ([]byte, bool)`
- `ResourceExists(basePackage, resourcePath string) bool`
- `LoadTextResource(basePackage, resourcePath string) string` / `LoadGlobalResource(basePackage, resourcePath string) string`
- `ResourceBundle_GetBundle(baseName string) *ResourceBundle` (`(*ResourceBundle).GetString(key string) string`, `ContainsKey(key string) bool`), `ParseJavaProperties(string) map[string]string`
- `ExternsZipBytes() []byte`, `RuntimeLibsTypedAstBytes() []byte`, `RegisterOverride`, `UnregisterOverride`, `ClearOverrides`.
