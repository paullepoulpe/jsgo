import com.google.javascript.jscomp.diagnostic.LogsGson;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableMultimap;
import com.google.javascript.jscomp.jarjar.com.google.gson.Gson;
import com.google.javascript.jscomp.jarjar.com.google.gson.GsonBuilder;
import com.google.javascript.jscomp.jarjar.com.google.gson.JsonArray;
import com.google.javascript.jscomp.jarjar.com.google.gson.JsonNull;
import com.google.javascript.jscomp.jarjar.com.google.gson.JsonObject;
import com.google.javascript.jscomp.jarjar.com.google.gson.JsonPrimitive;
import java.lang.reflect.Method;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.EnumSet;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.TreeMap;

/**
 * Vectors for jutil's Gson.toJson (internal/jutil/gson_to_json_test.go) and LogsGson.toJson
 * (internal/diagnostic/logs_gson_test.go), from the Gson 2.9.1 shaded into the pinned jar.
 *
 * <pre>
 * java -cp $CP GsonToJsonVectors > gson_tojson.tsv
 * </pre>
 *
 * One line per (configuration, case): configuration TAB case name TAB the JSON or the exception
 * class and message, escaped like RegExpDump. The Go tests build the same values by case name.
 */
public final class GsonToJsonVectors {
  enum Color {
    RED,
    GREEN
  }

  static class Point {
    int x = 1;
    String name = "p<q>";
    String missing = null;
    transient int t = 5;
    static int s = 3;
    List<Integer> list = List.of(1, 2);
    Object me = this;
    double d = 0.5;
    char c = 'c';
    boolean b = true;
    long l = 1L << 40;
  }

  static class Point3 extends Point {
    int z = 9;
    Map<String, Integer> m = new LinkedHashMap<>(Map.of("k", 1));
  }

  static class Holder {
    Object value;

    Holder(Object value) {
      this.value = value;
    }
  }

  static class AbleThing implements LogsGson.Able {
    Object v;

    AbleThing(Object v) {
      this.v = v;
    }

    @Override
    public Object toLogsGson() {
      return v;
    }
  }

  static Map<String, Object> cases() {
    Map<String, Object> c = new LinkedHashMap<>();
    c.put("null", null);
    c.put("string", "a<b>&'=\"\\  \u0001\t\n\b\f\r/é😀");
    c.put("loneSurrogate", "\ud800x\udc00");
    c.put("int", 42);
    c.put("long", -7L);
    c.put("short", (short) 3);
    c.put("byte", (byte) -1);
    c.put("intMin", Integer.MIN_VALUE);
    c.put("longMax", Long.MAX_VALUE);
    c.put("double1", 1.0);
    c.put("double0.1", 0.1);
    c.put("double1e21", 1e21);
    c.put("doubleNegZero", -0.0);
    c.put("double1e-7", 1e-7);
    c.put("double123456789", 123456789.0);
    c.put("double1e7", 1e7);
    c.put("doubleMin", Double.MIN_VALUE);
    c.put("doubleNaN", Double.NaN);
    c.put("doubleInf", Double.NEGATIVE_INFINITY);
    c.put("char", 'x');
    c.put("charLS", ' ');
    c.put("charQuote", '"');
    c.put("true", true);
    c.put("listMixed", Arrays.asList(1, "a", null, true));
    c.put("listEmpty", List.of());
    c.put("listNested", List.of(List.of(1), List.of()));
    c.put("intArray", new int[] {1, 2, 3});
    c.put("stringArray", new String[] {"x", null});
    Map<String, Object> m = new LinkedHashMap<>();
    m.put("a", 1);
    m.put("b", null);
    m.put("c", List.of(1, 2));
    m.put("d", "<");
    c.put("map", m);
    Map<Integer, String> mi = new LinkedHashMap<>();
    mi.put(3, "x");
    mi.put(1, "y");
    c.put("mapIntKeys", mi);
    Map<String, Integer> mn = new LinkedHashMap<>();
    mn.put(null, 1);
    c.put("mapNullKey", mn);
    c.put("enum", Color.GREEN);
    c.put("enumList", List.of(Color.RED, Color.GREEN));
    Map<Color, Integer> me = new LinkedHashMap<>();
    me.put(Color.RED, 1);
    c.put("enumKeyMap", me);
    c.put("enumSet", EnumSet.of(Color.GREEN, Color.RED));
    c.put("point", new Point());
    c.put("point3", new Point3());
    c.put("holderNull", new Holder(null));
    c.put("holderPoint", new Holder(new Point()));
    c.put("holderList", new Holder(List.of(1)));
    c.put("ableList", new AbleThing(List.of(0)));
    c.put("ableNull", new AbleThing(null));
    c.put("ableString", new AbleThing("<"));
    c.put("ableInList", List.of(new AbleThing(Map.of("a", 1))));
    c.put("holderAble", new Holder(new AbleThing(1)));
    JsonObject o = new JsonObject();
    o.add("a", new JsonPrimitive(1));
    JsonArray arr = new JsonArray();
    arr.add(new JsonPrimitive(true));
    arr.add(JsonNull.INSTANCE);
    o.add("b", arr);
    o.add("c", new JsonPrimitive("<x>"));
    c.put("jsonTree", o);
    c.put("nested", List.of(Map.of("x", List.of(Map.of()))));
    TreeMap<String, Integer> tm = new TreeMap<>();
    tm.put("b", 1);
    tm.put("a", 2);
    c.put("treeMap", tm);
    c.put("linkedHashSet", new LinkedHashSet<>(List.of("b", "a")));
    c.put("arrayList", new ArrayList<>(List.of(3, 1)));
    return c;
  }

  static Map<String, Object> logsGsonOnlyCases() {
    Map<String, Object> c = new LinkedHashMap<>();
    c.put(
        "multimap",
        ImmutableMultimap.builder().put("foo", 0).put("foo", 1).put("bar", 2).build());
    c.put(
        "holderMultimap",
        new Holder(ImmutableMultimap.builder().put("<", "&").build()));
    return c;
  }

  public static void main(String[] args) throws Exception {
    Map<String, Gson> configs = new LinkedHashMap<>();
    configs.put("default", new Gson());
    configs.put("builder", new GsonBuilder().create());
    configs.put("noHtmlPretty", new GsonBuilder().disableHtmlEscaping().setPrettyPrinting().create());
    configs.put("serializeNulls", new GsonBuilder().serializeNulls().create());
    StringBuilder sb = new StringBuilder();
    for (Map.Entry<String, Gson> config : configs.entrySet()) {
      for (Map.Entry<String, Object> c : cases().entrySet()) {
        String result;
        try {
          result = config.getValue().toJson(c.getValue());
        } catch (Throwable e) {
          result = e.getClass().getName() + ": " + e.getMessage();
        }
        sb.append(config.getKey()).append('\t').append(c.getKey()).append('\t');
        sb.append(RegExpDump.escape(result)).append('\n');
      }
    }
    Method toJson = LogsGson.class.getDeclaredMethod("toJson", Object.class);
    toJson.setAccessible(true);
    Map<String, Object> all = cases();
    all.putAll(logsGsonOnlyCases());
    for (Map.Entry<String, Object> c : all.entrySet()) {
      String result;
      try {
        result = (String) toJson.invoke(null, c.getValue());
      } catch (java.lang.reflect.InvocationTargetException ite) {
        Throwable e = ite.getCause();
        result = e.getClass().getName() + ": " + e.getMessage();
      }
      sb.append("LogsGson\t").append(c.getKey()).append('\t');
      sb.append(RegExpDump.escape(result)).append('\n');
    }
    System.out.print(sb);
  }
}
