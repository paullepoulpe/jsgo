import com.google.javascript.jscomp.jarjar.com.google.gson.Gson;
import com.google.javascript.jscomp.jarjar.com.google.gson.JsonArray;
import com.google.javascript.jscomp.jarjar.com.google.gson.JsonElement;
import com.google.javascript.jscomp.jarjar.com.google.gson.JsonNull;
import com.google.javascript.jscomp.jarjar.com.google.gson.JsonObject;
import com.google.javascript.jscomp.jarjar.com.google.gson.JsonPrimitive;
import com.google.javascript.jscomp.jarjar.com.google.gson.internal.LazilyParsedNumber;
import java.io.FileOutputStream;
import java.io.OutputStreamWriter;
import java.io.PrintWriter;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;

/**
 * Vectors for JsonElement equals/hashCode in internal/jutil (gson_json_element.go): the Gson 2.9.1
 * that closure-compiler.jar shades. A list of elements (parsed JSON documents, primitives built
 * from Integer, Long, Double and String values, LazilyParsedNumbers), then equals for every
 * ordered pair and hashCode for every element whose tree holds no JsonNull (JsonNull.hashCode()
 * is the identity hash of the JsonNull class, which is arbitrary).
 *
 * <p>Output: internal/jutil/testdata/gson_equals.tsv. Lines:
 *
 * <ul>
 *   <li>{@code D <index> <kind> <payload>}: an element; kind is json (payload parsed by
 *       Gson.fromJson(payload, JsonElement.class)), int, long, double (the payload is the
 *       Double.toString of the value), string, bool, lazy (a LazilyParsedNumber), null,
 *       objectOfInt ({"a": Integer 1}) or arrayOfLongAndDouble ([Long 1, Double 2.0]);
 *   <li>{@code H <index> <hashCode or ->};
 *   <li>{@code E <index> <index> <0|1 per element>}: equals of element i with every element j.
 * </ul>
 *
 * <p>Payloads are printable ASCII. Usage: java -cp closure-compiler.jar:. GsonEqualsVectors out.tsv
 */
public final class GsonEqualsVectors {
  private static final Gson GSON = new Gson();

  private static final List<String[]> defs = new ArrayList<>();
  private static final List<JsonElement> elements = new ArrayList<>();

  private static void add(String kind, String payload, JsonElement e) {
    defs.add(new String[] {kind, payload});
    elements.add(e);
  }

  private static void json(String s) {
    add("json", s, GSON.fromJson(s, JsonElement.class));
  }

  private static boolean hasNull(JsonElement e) {
    if (e.isJsonNull()) {
      return true;
    }
    if (e.isJsonObject()) {
      for (var entry : e.getAsJsonObject().entrySet()) {
        if (hasNull(entry.getValue())) {
          return true;
        }
      }
    }
    if (e.isJsonArray()) {
      for (JsonElement m : e.getAsJsonArray()) {
        if (hasNull(m)) {
          return true;
        }
      }
    }
    return false;
  }

  public static void main(String[] args) throws Exception {
    String[] docs = {
      "1", "1.0", "1e0", "01", "-0", "0", "0.0", "-0.0", "3", "3.0", "3.5", "1e400", "-1e400",
      "9223372036854775807", "9223372036854775808", "-9223372036854775808", "4294967296",
      "4294967297", "2147483648", "12345678901234567890123", "0.1", "1e-5", "NaN", "Infinity",
      "\"1\"", "\"3\"", "\"a\"", "\"\"", "\"\\u00e9\\ud83d\\ude00\"", "true", "false", "\"true\"",
      "null", "{}", "[]", "[null]", "{\"a\":null}", "{\"a\":1}", "{\"a\":1.0}", "{\"a\":\"1\"}",
      "{\"b\":1}", "{\"a\":1,\"b\":2}", "{\"b\":2,\"a\":1}", "{\"a\":1,\"b\":3}", "[1,2]", "[2,1]",
      "[1.0,2.0]", "[1,2,3]", "[[1],[2]]", "[[1],[2.0]]", "{\"a\":[1,{\"b\":true}]}",
      "{\"a\":[1,{\"b\":true}],\"c\":\"x\"}", "{\"a\":[1.00,{\"b\":true}]}", "[\"a\",true,1]",
      "{\"version\":3,\"file\":\"x\",\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[]}",
      "{\"version\":3.0,\"file\":\"x\",\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[]}",
    };
    for (String d : docs) {
      json(d);
    }
    for (int v : new int[] {0, 1, 3, -1, Integer.MAX_VALUE, Integer.MIN_VALUE}) {
      add("int", String.valueOf(v), new JsonPrimitive(Integer.valueOf(v)));
    }
    for (long v :
        new long[] {0L, 3L, 4294967296L, -4294967297L, Long.MAX_VALUE, Long.MIN_VALUE, 1L << 53,
          (1L << 53) + 1}) {
      add("long", String.valueOf(v), new JsonPrimitive(Long.valueOf(v)));
    }
    for (double v :
        new double[] {0.0, -0.0, 1.0, 3.0, 3.5, 0.1, 1e300, Double.NaN,
          Double.POSITIVE_INFINITY, 9.223372036854776E18, 9.007199254740992E15}) {
      add("double", Double.toString(v), new JsonPrimitive(Double.valueOf(v)));
    }
    for (String v : new String[] {"3", "3.0", "x", "", "true"}) {
      add("string", v, new JsonPrimitive(v));
    }
    add("bool", "true", new JsonPrimitive(Boolean.TRUE));
    add("bool", "false", new JsonPrimitive(Boolean.FALSE));
    for (String v : new String[] {"3", "3.0", "-0", "9007199254740993", "1e400", "NaN"}) {
      add("lazy", v, new JsonPrimitive(new LazilyParsedNumber(v)));
    }
    add("null", "", JsonNull.INSTANCE);
    JsonObject built = new JsonObject();
    built.add("a", new JsonPrimitive(Integer.valueOf(1)));
    add("objectOfInt", "a=1", built);
    JsonArray builtArray = new JsonArray();
    builtArray.add(new JsonPrimitive(Long.valueOf(1)));
    builtArray.add(new JsonPrimitive(Double.valueOf(2)));
    add("arrayOfLongAndDouble", "1,2.0", builtArray);

    try (PrintWriter out =
        new PrintWriter(
            new OutputStreamWriter(new FileOutputStream(args[0]), StandardCharsets.UTF_8))) {
      for (int i = 0; i < defs.size(); i++) {
        out.println("D\t" + i + "\t" + defs.get(i)[0] + "\t" + esc(defs.get(i)[1]));
      }
      for (int i = 0; i < elements.size(); i++) {
        JsonElement e = elements.get(i);
        out.println("H\t" + i + "\t" + (hasNull(e) ? "-" : String.valueOf(e.hashCode())));
      }
      for (int i = 0; i < elements.size(); i++) {
        StringBuilder sb = new StringBuilder();
        for (int j = 0; j < elements.size(); j++) {
          sb.append(elements.get(i).equals(elements.get(j)) ? '1' : '0');
        }
        out.println("E\t" + i + "\t" + sb);
      }
    }
  }

  static String esc(String s) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '\\') {
        sb.append("\\\\");
      } else if (c >= 0x20 && c <= 0x7e && c != '\t') {
        sb.append(c);
      } else {
        sb.append(String.format("\\u%04x", (int) c));
      }
    }
    return sb.toString();
  }
}
