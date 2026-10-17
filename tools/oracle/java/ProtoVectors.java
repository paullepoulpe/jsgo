// Oracle tool: writes random protobuf messages of the types reconstructed by
// ProtoGoGen (serialized bytes, toString()) to
// internal/jscomp/testdata/proto_vectors.tsv, for the Go round-trip test.
//   java -cp .cache/closure-compiler.jar:.cache/oracle-build ProtoVectors internal/jscomp/testdata
import com.google.javascript.jscomp.ConformanceConfig;
import com.google.javascript.jscomp.Requirement;
import com.google.javascript.jscomp.RequirementScopeEntry;
import com.google.javascript.jscomp.jarjar.com.google.protobuf.ByteString;
import com.google.javascript.jscomp.jarjar.com.google.protobuf.Message;
import com.google.javascript.jscomp.serialization.*;
import java.io.PrintStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Random;

public class ProtoVectors {
  static Random r = new Random(777);

  static String hex(byte[] b) {
    StringBuilder sb = new StringBuilder();
    for (byte x : b) sb.append(String.format("%02x", x & 0xff));
    return sb.toString();
  }

  static String esc(String s) {
    return s.replace("\\", "\\\\").replace("\n", "\\n").replace("\t", "\\t");
  }

  static String str() {
    String[] pool = {"", "a", "foo.js", "\u00e9\u4e2d", "x\ny", "q\"uote", "tab\tbed", "\ud83d\ude00", "back\\slash", "'"};
    return pool[r.nextInt(pool.length)];
  }

  static ByteString bytes() {
    byte[] b = new byte[r.nextInt(6)];
    r.nextBytes(b);
    return ByteString.copyFrom(b);
  }

  static int smallInt() {
    switch (r.nextInt(5)) {
      case 0: return 0;
      case 1: return r.nextInt(10);
      case 2: return -r.nextInt(1000);
      case 3: return r.nextInt();
      default: return r.nextInt(300);
    }
  }

  static OptimizationJsdoc jsdoc() {
    OptimizationJsdoc.Builder b = OptimizationJsdoc.newBuilder();
    int n = r.nextInt(4);
    for (int i = 0; i < n; i++) {
      if (r.nextInt(6) == 0) b.addKindValue(100 + r.nextInt(5));
      else b.addKind(JsdocTag.values()[r.nextInt(JsdocTag.values().length - 1)]);
    }
    if (r.nextBoolean()) b.setLicenseTextPointer(smallInt());
    if (r.nextBoolean()) b.setDescriptionPointer(smallInt());
    if (r.nextInt(4) == 0) b.setClosureUnawareModePointer(smallInt());
    return b.build();
  }

  static AstNode ast(int depth) {
    AstNode.Builder b = AstNode.newBuilder();
    if (r.nextInt(10) == 0) b.setKindValue(500 + r.nextInt(10));
    else b.setKind(NodeKind.values()[r.nextInt(NodeKind.values().length - 1)]);
    switch (r.nextInt(5)) {
      case 0: b.setStringValuePointer(r.nextInt(1000)); break;
      case 1: b.setDoubleValue(r.nextBoolean() ? r.nextDouble() * 1000 : (r.nextBoolean() ? -0.0 : Double.NaN)); break;
      case 2: b.setTemplateStringValue(TemplateStringValue.newBuilder().setRawStringPointer(r.nextInt(50)).setCookedStringPointer(r.nextInt(3) - 1)); break;
      default: break;
    }
    if (r.nextBoolean()) b.setRelativeLine(r.nextInt(20) - 10);
    if (r.nextBoolean()) b.setRelativeColumn(r.nextInt(200) - 100);
    if (r.nextInt(3) == 0) b.setJsdoc(jsdoc());
    if (r.nextInt(3) == 0) b.setOriginalNamePointer(r.nextInt(100));
    if (r.nextInt(3) == 0) b.setBooleanProperties(r.nextLong() >>> r.nextInt(64));
    if (r.nextInt(3) == 0) b.setType(r.nextInt(20));
    if (r.nextInt(4) == 0) b.setSourceFile(r.nextInt(5));
    if (depth < 4) {
      int n = r.nextInt(4);
      for (int i = 0; i < n; i++) b.addChild(ast(depth + 1));
    }
    return b.build();
  }

  static TypePool typePool() {
    TypePool.Builder b = TypePool.newBuilder();
    int n = r.nextInt(5);
    for (int i = 0; i < n; i++) {
      TypeProto.Builder t = b.addTypeBuilder();
      if (r.nextBoolean()) {
        ObjectTypeProto.Builder o = t.getObjectBuilder();
        o.setUuid(bytes()).setIsInvalidating(r.nextBoolean()).setMarkedConstructor(r.nextBoolean())
            .setClosureAssert(r.nextBoolean()).setPropertiesKeepOriginalName(r.nextBoolean());
        for (int k = r.nextInt(3); k > 0; k--) o.addInstanceType(smallInt());
        for (int k = r.nextInt(3); k > 0; k--) o.addPrototype(smallInt());
        for (int k = r.nextInt(3); k > 0; k--) o.addOwnProperty(smallInt());
        if (r.nextInt(4) == 0) o.getDebugInfoBuilder().addTypenamePointer(smallInt());
      } else if (r.nextBoolean()) {
        UnionTypeProto.Builder u = t.getUnionBuilder();
        for (int k = r.nextInt(4); k > 0; k--) u.addUnionMember(smallInt());
      }
    }
    for (int k = r.nextInt(3); k > 0; k--) b.addDisambiguationEdgesBuilder().setSubtype(smallInt()).setSupertype(smallInt());
    if (r.nextBoolean()) {
      TypePool.DebugInfo.Builder d = b.getDebugInfoBuilder();
      for (int k = r.nextInt(3); k > 0; k--) d.addMismatchBuilder().setSourceRef(str()).addInvolvedColor(smallInt()).addInvolvedColor(smallInt());
    }
    return b.build();
  }

  static SourceFileProto sourceFile() {
    SourceFileProto.Builder b = SourceFileProto.newBuilder().setFilename(str());
    switch (r.nextInt(5)) {
      case 0: b.setPreloadedContents(str()); break;
      case 1: b.setFileOnDisk(SourceFileProto.FileOnDisk.newBuilder().setActualPath(str()).setCharset(r.nextBoolean() ? "" : "US-ASCII")); break;
      case 2: b.setZipEntry(SourceFileProto.ZipEntryOnDisk.newBuilder().setZipPath(str()).setEntryName(str())); break;
      case 3: b.setStubFile(r.nextBoolean()); break;
      default: break;
    }
    if (r.nextBoolean()) b.setSourceKindValue(r.nextInt(4));
    if (r.nextBoolean()) b.setIsClosureUnawareCode(true);
    if (r.nextBoolean()) b.setNumBytesPlusOne(r.nextInt());
    if (r.nextBoolean()) b.setNumLinesPlusOne(r.nextInt(100));
    return b.build();
  }

  static Requirement requirement() {
    Requirement.Builder b = Requirement.newBuilder();
    if (r.nextBoolean()) b.setErrorMessage(str());
    for (int k = r.nextInt(3); k > 0; k--) b.addWhitelist(str());
    for (int k = r.nextInt(2); k > 0; k--) b.addAllowlistRegexp(str());
    for (int k = r.nextInt(2); k > 0; k--) b.addAllowlistEntry(RequirementScopeEntry.newBuilder().addPrefix(str()).addRegexp(str()).setReason(RequirementScopeEntry.Reason.values()[r.nextInt(4)]).setAutomaticallyPrune(r.nextBoolean()));
    if (r.nextBoolean()) b.setType(Requirement.Type.values()[r.nextInt(Requirement.Type.values().length)]);
    for (int k = r.nextInt(3); k > 0; k--) b.addValue(str());
    if (r.nextBoolean()) b.setTypeMatchingStrategy(Requirement.TypeMatchingStrategy.values()[r.nextInt(5)]);
    if (r.nextBoolean()) b.setRuleId(str());
    if (r.nextBoolean()) b.setReportLooseTypeViolations(r.nextBoolean());
    if (r.nextBoolean()) b.setSeverity(Requirement.Severity.values()[r.nextInt(3)]);
    for (int k = r.nextInt(2); k > 0; k--) b.addConfigFile(str());
    if (r.nextBoolean()) b.setAllowExtendingValue(r.nextBoolean());
    return b.build();
  }

  static byte[] withUnknown(byte[] raw) {
    if (r.nextInt(6) != 0) return raw;
    java.io.ByteArrayOutputStream out = new java.io.ByteArrayOutputStream();
    out.writeBytes(raw);
    int field = 900 + r.nextInt(5);
    switch (r.nextInt(4)) {
      case 0: out.write((field << 3) & 0x7f | 0x80); out.write(field >> 4); out.write(r.nextInt(128)); break;
      case 1: out.write((field << 3 | 2) & 0x7f | 0x80); out.write(field >> 4); out.write(2); out.write(0x08); out.write(0x05); break;
      case 2: out.write((field << 3 | 2) & 0x7f | 0x80); out.write(field >> 4); out.write(1); out.write(0xff); break;
      default: out.write((field << 3 | 5) & 0x7f | 0x80); out.write(field >> 4); out.write(1); out.write(2); out.write(3); out.write(4); break;
    }
    return out.toByteArray();
  }

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    Files.createDirectories(out);
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("proto_vectors.tsv")), true, "UTF-8")) {
      for (int i = 0; i < 400; i++) {
        String type;
        Message m;
        switch (i % 6) {
          case 0: type = "AstNode"; m = ast(0); break;
          case 1: type = "TypePool"; m = typePool(); break;
          case 2: type = "SourceFileProto"; m = sourceFile(); break;
          case 3: type = "Requirement"; m = requirement(); break;
          case 4: type = "OptimizationJsdoc"; m = jsdoc(); break;
          default:
            type = "ConformanceConfig";
            ConformanceConfig.Builder c = ConformanceConfig.newBuilder();
            for (int k = r.nextInt(3); k > 0; k--) c.addRequirement(requirement());
            if (r.nextBoolean()) c.setLibraryLevelNonAllowlistedConformanceViolationsBehavior(ConformanceConfig.LibraryLevelNonAllowlistedConformanceViolationsBehavior.values()[r.nextInt(3)]);
            m = c.build();
            break;
        }
        byte[] raw = withUnknown(m.toByteArray());
        Message parsed = m.getParserForType().parseFrom(raw);
        ps.println(type + "\t" + hex(raw) + "\t" + hex(parsed.toByteArray()) + "\t" + esc(parsed.toString()));
      }
    }
  }
}
