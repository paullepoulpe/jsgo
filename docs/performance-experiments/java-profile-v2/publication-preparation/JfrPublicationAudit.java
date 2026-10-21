// Source only: bounded standard-consumer audit. Never emits observed string values/hashes.
// Usage: JfrPublicationAudit recording.jfr [private-exact-allowlist.tsv]
import java.io.BufferedReader;
import java.io.IOException;
import java.io.PrintWriter;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.time.Duration;
import java.util.HashMap;
import java.util.HashSet;
import java.util.IdentityHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.TreeMap;
import jdk.jfr.EventType;
import jdk.jfr.ValueDescriptor;
import jdk.jfr.consumer.RecordedEvent;
import jdk.jfr.consumer.RecordedObject;
import jdk.jfr.consumer.RecordingFile;

public final class JfrPublicationAudit {
  private static final long MAX_BYTES = 512L << 20, MAX_EVENTS = 5_000_000;
  private static final long MAX_NODES = 100_000_000, MAX_SECONDS = 120;
  private static final int MAX_DEPTH = 64, MAX_SCHEMAS = 20_000;
  private static final int MAX_STRING = 1 << 20, MAX_POLICY_ENTRIES = 100_000;
  private static final String ENV = "jdk.InitialEnvironmentVariable";
  private static final String PROP = "jdk.InitialSystemProperty";
  private static final Set<String> BENIGN_ENV_KEYS = Set.of(
      "LC_ALL", "LANG", "LANGUAGE", "PATH", "HOME", "PWD", "HOSTNAME", "USER",
      "LOGNAME", "TZ", "PYTHON_VERSION", "PYTHON_SHA256", "GPG_KEY", "PYTHON_PIP_VERSION",
      "PYTHON_SETUPTOOLS_VERSION", "PYTHON_GET_PIP_URL", "PYTHON_GET_PIP_SHA256");
  private static final Set<String> BENIGN_PROPERTY_KEYS = Set.of(
      "java.version", "java.version.date", "java.vendor", "java.vendor.url", "java.vendor.url.bug",
      "java.vendor.version", "java.home", "java.class.path", "java.library.path", "java.class.version",
      "java.vm.name", "java.vm.vendor", "java.vm.version", "java.vm.info", "java.vm.compressedOopsMode",
      "java.vm.specification.name", "java.vm.specification.vendor", "java.vm.specification.version",
      "java.runtime.name", "java.runtime.version", "java.specification.name", "java.specification.vendor",
      "java.specification.version", "java.specification.maintenance.version", "java.io.tmpdir",
      "java.net.preferIPv4Stack", "java.net.preferIPv6Addresses", "os.name", "os.arch", "os.version",
      "user.name", "user.home", "user.dir", "user.language", "user.country", "user.timezone",
      "file.encoding", "native.encoding", "file.separator", "path.separator", "line.separator",
      "stdout.encoding", "stderr.encoding", "sun.stdout.encoding", "sun.stderr.encoding",
      "sun.arch.data.model", "sun.boot.library.path", "sun.cpu.endian", "sun.io.unicode.encoding",
      "sun.java.command", "sun.java.launcher", "sun.jnu.encoding", "jdk.debug", "jdk.module.main",
      "jdk.module.main.class", "jdk.module.path", "jdk.module.addmods.0");
  // Names/descriptors originate from the pinned JDK/compiler classes. They are
  // inventoried as diagnostic metadata, not treated as environment/path payloads.
  private static final Set<String> DIAGNOSTIC_METADATA = Set.of(
      "java.lang.Class.name", "jdk.types.Method.name", "jdk.types.Method.descriptor",
      "jdk.types.Package.name", "jdk.types.Module.name", "jdk.types.Module.version",
      "jdk.types.FrameType.description",
      "jdk.types.ThreadState.name", "jdk.types.GCName.name", "jdk.types.GCCause.cause",
      "jdk.types.CompilerType.compiler", "jdk.types.CodeBlobType.type");

  private static final class Counts {
    long seen, exactAllowed, diagnostic, unresolved;
  }
  private static final class Stop extends RuntimeException {
    final String reason;
    Stop(String reason) { this.reason = reason; }
  }
  private final Map<String, Set<String>> exact = new HashMap<>();
  private final Map<String, Counts> strings = new TreeMap<>();
  private final Map<String, Long> events = new TreeMap<>(), keys = new TreeMap<>();
  private final Map<String, String> schema = new TreeMap<>();
  private final Set<Long> schemaTypes = new HashSet<>();
  private long started, eventCount, nodes, unresolved, omittedNames, unknownKeys;
  private String outcome = "COMPLETE";

  public static void main(String[] args) {
    JfrPublicationAudit audit = new JfrPublicationAudit();
    try {
      if (args.length < 1 || args.length > 2) throw new Stop("USAGE_RECORDING_OPTIONAL_PRIVATE_POLICY");
      if (args.length == 2) audit.policy(Path.of(args[1]));
      audit.exact.computeIfAbsent("ENV:LC_ALL", k -> new HashSet<>()).add("C.UTF-8");
      audit.scan(Path.of(args[0]));
    } catch (Stop stop) {
      audit.outcome = stop.reason;
    } catch (IOException | RuntimeException error) {
      // Exception messages can contain paths/payloads. Do not serialize them.
      audit.outcome = "READ_OR_SCHEMA_ERROR";
    }
    audit.write();
    if (!audit.outcome.equals("COMPLETE") || audit.unresolved > 0 || audit.omittedNames > 0 || audit.unknownKeys > 0)
      System.exit(2);
  }

  private void policy(Path path) throws IOException {
    if (Files.size(path) > (16L << 20)) throw new Stop("POLICY_SIZE_LIMIT");
    try (BufferedReader reader = Files.newBufferedReader(path, StandardCharsets.UTF_8)) {
      String line;
      int count = 0;
      while ((line = reader.readLine()) != null) {
        if (line.isEmpty() || line.startsWith("#")) continue;
        if (++count > MAX_POLICY_ENTRIES || line.length() > MAX_STRING) throw new Stop("POLICY_LIMIT");
        int tab = line.indexOf('\t');
        if (tab < 1 || line.indexOf('\t', tab + 1) >= 0) throw new Stop("POLICY_FORMAT");
        String selector = line.substring(0, tab);
        if (!(selector.startsWith("ENV:") || selector.startsWith("PROP:") || selector.startsWith("FIELD:")))
          throw new Stop("POLICY_SELECTOR");
        if (!selector.matches("[A-Za-z0-9_.$:@\\[\\]-]+")) throw new Stop("POLICY_SELECTOR");
        exact.computeIfAbsent(selector, k -> new HashSet<>()).add(unescape(line.substring(tab + 1)));
      }
    }
  }

  private static String unescape(String text) {
    StringBuilder out = new StringBuilder();
    for (int i = 0; i < text.length(); i++) {
      char c = text.charAt(i);
      if (c != '\\') { out.append(c); continue; }
      if (++i == text.length()) throw new Stop("POLICY_ESCAPE");
      out.append(switch (text.charAt(i)) {
        case '\\' -> '\\'; case 'n' -> '\n'; case 'r' -> '\r'; case 't' -> '\t';
        default -> throw new Stop("POLICY_ESCAPE");
      });
    }
    return out.toString();
  }

  private void scan(Path path) throws IOException {
    if (Files.size(path) > MAX_BYTES) throw new Stop("RECORDING_SIZE_LIMIT");
    started = System.nanoTime();
    try (RecordingFile file = new RecordingFile(path)) {
      for (EventType type : file.readEventTypes()) {
        String owner = schemaName(type.getName());
        events.putIfAbsent(owner, 0L);
        if (events.size() > MAX_SCHEMAS) throw new Stop("EVENT_TYPE_LIMIT");
        fields(owner, type.getFields(), 0);
      }
      while (file.hasMoreEvents()) {
        if (eventCount >= MAX_EVENTS) throw new Stop("EVENT_LIMIT");
        checkTime();
        RecordedEvent event = file.readEvent();
        eventCount++;
        String owner = event.getEventType().getName();
        events.merge(schemaName(owner), 1L, Long::sum);
        walk(event, owner, 0, new IdentityHashMap<>());
      }
    }
  }

  private void fields(String owner, List<ValueDescriptor> descriptors, int depth) {
    if (depth > MAX_DEPTH) throw new Stop("SCHEMA_DEPTH_LIMIT");
    for (ValueDescriptor field : descriptors) {
      String selector = owner + "." + schemaName(field.getName());
      String type = schemaName(field.getTypeName()) + (field.isArray() ? "[]" : "");
      if (field.getTypeName().equals("java.lang.String") || !field.getFields().isEmpty()) {
        schema.put(selector, type);
        if (schema.size() > MAX_SCHEMAS) throw new Stop("SCHEMA_COUNT_LIMIT");
      }
      if (field.getTypeName().equals("java.lang.String")) strings.computeIfAbsent(selector, k -> new Counts());
      if (!field.getFields().isEmpty() && schemaTypes.add(field.getTypeId()))
        fields(schemaName(field.getTypeName()), field.getFields(), depth + 1);
    }
  }

  private void walk(RecordedObject object, String owner, int depth, IdentityHashMap<Object, Boolean> visited) {
    if (depth > MAX_DEPTH) throw new Stop("VALUE_DEPTH_LIMIT");
    if (visited.put(object, Boolean.TRUE) != null) return;
    String specialKey = null;
    boolean initial = owner.equals(ENV) || owner.equals(PROP);
    if (initial && object.hasField("key")) {
      specialKey = object.getString("key");
      Set<String> benign = owner.equals(ENV) ? BENIGN_ENV_KEYS : BENIGN_PROPERTY_KEYS;
      if (specialKey != null && benign.contains(specialKey))
        keys.merge((owner.equals(ENV) ? "ENV:" : "PROP:") + specialKey, 1L, Long::sum);
      else unknownKeys++;
    }
    for (ValueDescriptor field : object.getFields()) {
      Object value = object.getValue(field.getName());
      inspect(value, owner, field, specialKey, depth + 1, visited);
    }
  }

  private void inspect(Object value, String owner, ValueDescriptor field, String specialKey,
                       int depth, IdentityHashMap<Object, Boolean> visited) {
    if (++nodes > MAX_NODES) throw new Stop("NODE_LIMIT");
    if ((nodes & 65535) == 0) checkTime();
    if (depth > MAX_DEPTH) throw new Stop("VALUE_DEPTH_LIMIT");
    if (value == null) return;
    if (value instanceof String text) {
      String rawSelector = owner + "." + field.getName();
      String outputSelector = schemaName(owner) + "." + schemaName(field.getName());
      Counts count = strings.computeIfAbsent(outputSelector, k -> new Counts());
      if (strings.size() > MAX_SCHEMAS) throw new Stop("STRING_SCHEMA_LIMIT");
      count.seen++;
      boolean allowed = false, diagnostic = false;
      if (text.length() <= MAX_STRING) {
        String special = owner.equals(ENV) ? "ENV:" : owner.equals(PROP) ? "PROP:" : null;
        if (special != null && field.getName().equals("value") && specialKey != null)
          allowed = exact.getOrDefault(special + specialKey, Set.of()).contains(text);
        else if (special != null && field.getName().equals("key"))
          allowed = (owner.equals(ENV) ? BENIGN_ENV_KEYS : BENIGN_PROPERTY_KEYS).contains(text);
        allowed |= text.isEmpty() || exact.getOrDefault("FIELD:" + rawSelector, Set.of()).contains(text);
        diagnostic = DIAGNOSTIC_METADATA.contains(rawSelector)
            && text.length() <= 4096 && text.chars().allMatch(c -> c >= 32 && c < 127);
      }
      if (allowed) count.exactAllowed++;
      else if (diagnostic) count.diagnostic++;
      else { count.unresolved++; unresolved++; }
    } else if (value instanceof RecordedObject nested) {
      walk(nested, field.getTypeName(), depth, visited);
    } else if (value instanceof Object[] array) {
      if (visited.put(array, Boolean.TRUE) != null) return;
      for (Object item : array) inspect(item, owner, field, specialKey, depth + 1, visited);
    } else if (!(value instanceof Number || value instanceof Boolean || value instanceof Character)) {
      // Unexpected structured/array representation is not silently skipped.
      throw new Stop("UNSUPPORTED_VALUE_REPRESENTATION");
    }
  }

  private String schemaName(String name) {
    if (name != null && name.length() <= 256 && name.matches("[A-Za-z_$][A-Za-z0-9_.$]*")) return name;
    omittedNames++;
    return "REDACTED_SCHEMA_NAME";
  }

  private void checkTime() {
    if (Duration.ofNanos(System.nanoTime() - started).getSeconds() > MAX_SECONDS)
      throw new Stop("TIME_LIMIT");
  }

  private void write() {
    PrintWriter out = new PrintWriter(System.out, false, StandardCharsets.UTF_8);
    String status = !outcome.equals("COMPLETE") ? "INCOMPLETE" :
        unresolved > 0 || omittedNames > 0 || unknownKeys > 0 ? "UNRESOLVED" : "PASS_STANDARD_CONSUMER_STRING_AUDIT";
    out.println("STATUS\t" + status);
    out.println("SCAN_OUTCOME\t" + outcome);
    out.println("TOTALS\tevents\t" + eventCount + "\tnodes\t" + nodes
        + "\tunresolved_strings\t" + unresolved + "\tunknown_initial_keys\t" + unknownKeys
        + "\tredacted_schema_names\t" + omittedNames);
    out.println("SCOPE\tInitial environment/property values and other reachable strings; standard nested JDK/compiler diagnostic metadata is classified separately; no observed values or value hashes emitted.");
    out.println("POLICY\tOnly LC_ALL=C.UTF-8 is built in; private exact allowlist values remain in memory. Unknown strings prevent PASS. Source classes and raw recordings are unchanged.");
    for (Map.Entry<String, Long> entry : events.entrySet())
      out.println("EVENT\t" + entry.getKey() + "\t" + entry.getValue());
    for (Map.Entry<String, String> entry : schema.entrySet())
      out.println("SCHEMA\t" + entry.getKey() + "\t" + entry.getValue());
    for (Map.Entry<String, Long> entry : keys.entrySet())
      out.println("BENIGN_KEY\t" + entry.getKey() + "\t" + entry.getValue());
    out.println("STRING_FIELDS\tfield\toccurrences\texact_allowed_or_empty\tdiagnostic_metadata\tunresolved");
    for (Map.Entry<String, Counts> entry : strings.entrySet()) {
      Counts c = entry.getValue();
      out.println("STRING\t" + entry.getKey() + "\t" + c.seen + "\t" + c.exactAllowed + "\t" + c.diagnostic + "\t" + c.unresolved);
    }
    out.flush();
    if (out.checkError()) System.exit(3);
  }
}
