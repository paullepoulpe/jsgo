import com.google.javascript.jscomp.jarjar.com.google.gson.Gson;
import com.google.javascript.jscomp.jarjar.com.google.gson.JsonArray;
import com.google.javascript.jscomp.jarjar.com.google.gson.JsonElement;
import com.google.javascript.jscomp.jarjar.com.google.gson.JsonObject;
import com.google.javascript.jscomp.jarjar.com.google.gson.internal.LazilyParsedNumber;
import java.io.FileOutputStream;
import java.io.OutputStreamWriter;
import java.io.PrintWriter;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Random;
import java.util.function.Supplier;
import java.util.zip.GZIPOutputStream;

/**
 * Vectors for the Gson 2.9.1 runtime in internal/jutil (gson*.go): the Gson
 * that closure-compiler.jar shades. For each JSON input, the result of
 * Gson.fromJson(input, JsonElement/JsonObject/JsonArray.class) (toString of
 * the tree, or the exception's toString), plus the getAs* accessors on the
 * values; and LazilyParsedNumber.intValue/longValue on number strings.
 *
 * <p>Output (gzipped): internal/jutil/testdata/gson.tsv.gz. Every field is escaped to
 * printable ASCII: '\\' as "\\\\", tab as "\\t", newline as "\\n", any other
 * char outside 0x20..0x7e as "\\uXXXX" (UTF-16 code units, so lone surrogates
 * survive).
 *
 * <p>Usage: java -cp closure-compiler.jar:. GsonVectors out.tsv.gz
 */
public final class GsonVectors {
  private static final Gson GSON = new Gson();

  static String esc(String s) {
    if (s == null) {
      return "\\0";
    }
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '\\') {
        sb.append("\\\\");
      } else if (c == '\t') {
        sb.append("\\t");
      } else if (c == '\n') {
        sb.append("\\n");
      } else if (c >= 0x20 && c <= 0x7e) {
        sb.append(c);
      } else {
        sb.append(String.format("\\u%04x", (int) c));
      }
    }
    return sb.toString();
  }

  static String run(Supplier<Object> s) {
    try {
      Object o = s.get();
      return "=" + o;
    } catch (Throwable t) {
      return "!" + t;
    }
  }

  static String accessors(JsonElement e) {
    StringBuilder sb = new StringBuilder();
    sb.append(run(() -> e.getAsInt())).append(" | ");
    sb.append(run(() -> e.getAsLong())).append(" | ");
    sb.append(run(() -> e.getAsString())).append(" | ");
    sb.append(run(() -> e.getAsBoolean())).append(" | ");
    sb.append(run(() -> e.getAsJsonObject())).append(" | ");
    sb.append(run(() -> e.getAsJsonArray())).append(" | ");
    sb.append(run(() -> e.isJsonNull())).append(e.isJsonPrimitive()).append(e.isJsonArray()).append(e.isJsonObject());
    return sb.toString();
  }

  static List<String> inputs() {
    List<String> in = new ArrayList<>();
    String[] fixed = {
      "", " ", "\n", "\t\r\n ", "{}", "[]", "null", "true", "false", "1", "-1", "\"s\"", "'s'", "s",
      "{\"a\":1}", "{'a':'b'}", "{a:b}", "{a=b}", "{a=>b}", "{a=b>c}", "{\"a\":1;\"b\":2}", "[1;2]",
      "[1,,2]", "[,]", "[,1]", "[1,]", "[;]", "{\"a\":1,}", "{,}", "{;}", "{\"a\"}", "{\"a\":}",
      "/* c */ {}", "// c\n{}", "# c\n{}", "/{}", "/*{}", "/* unterminated", ")]}'\n{}", ")]}'{}",
      ")]}", "\ufeff{}", "\ufeff", "{} x", "{}{}", "{} //c", "{} /* c */", "{} #c", "{}\n", "{} ",
      "[] ]", "{\"a\":01}", "{\"a\":-0}", "{\"a\":0}", "{\"a\":-}", "{\"a\":1.}", "{\"a\":.5}",
      "{\"a\":1.5e3}", "{\"a\":1e400}", "{\"a\":1E+2}", "{\"a\":1e}", "{\"a\":1e+}", "{\"a\":+1}",
      "9223372036854775807", "9223372036854775808", "-9223372036854775808", "-9223372036854775809",
      "2147483647", "2147483648", "-2147483649", "{\"a\":\"\\u00e9\\ud83d\\ude00\\ud800\"}",
      "\"\\x\"", "\"\\u12\"", "\"\\u12G4\"", "\"\\u", "\"\\", "\"abc", "'abc", "{\"a\"", "{\"a\":",
      "[1", "[1,", "{\"a\" 1}", "NaN", "{\"a\":NaN}", "{\"a\":Infinity}", "{\"a\":-Infinity}",
      "{\"a\":tru}", "{\"a\":truex}", "{\"a\":nullx}", "{\"a\":TRUE}", "{\"a\":True}", "{\"a\":nULL}",
      "{\"a\":fAlSe}", "{\"a\":t}", "{\"a\":nul}", "{\"a\":{\"b\":[1,{\"c\":null}]}}",
      "{\"a\":1,\"b\":2,\"a\":3}", "{\"\":1}", "{\n\"a\":1,\n\"b\":\n  x y}", "\n\n  {\"a\":[1,\n2 3]}",
      "{\"a\":\"b\"\"c\":1}", "[\"a\" \"b\"]", "{\"a\":\"\\t\\b\\n\\r\\f\\/\\\\\\\"\\'\"}",
      "{\"a\":\"\\\nb\"}", "{\"a\":\"x\u2028y\u2029z\u0001\u001f\u007f<>&='\"}", "\"\u0000\"",
      "{\"version\":3,\"file\":\"x\",\"mappings\":\"AAAA\",\"sources\":[\"a\"],\"names\":[]}",
      "{\"a\":[1,2]}", "{\"a\":[\"7\"]}", "{\"a\":[]}", "{\"a\":[true]}", "{\"a\":\"12\"}",
      "{\"a\":\"+12\"}", "{\"a\":\" 12\"}", "{\"a\":\"1.5\"}", "{\"a\":\"TRUE\"}", "{\"a\":\"\"}",
      "{\"a\":\"\\u0661\\u0662\"}", "{\"a\":1.9}", "{\"a\":-1.9}", "{\"a\":1e2}", "{\"a\":4294967296}",
      "{\"a\":-4294967297.5}", "{\"a\":12345678901234567890123}", "{\"a\":1e-5}",
      "{a:1, 'b':2, \"c\":3, d:'x', e:\"y\", f:z}", "{a:1 b:2}", "[a b]", "[a,b]", "{\"a\\u0062\":1}",
      "{\"a\":1}}", "{\"a\":1}]", "[1]]", "\"a\"\"b\"", "1 2", "true false", "{\"a\":[[[[[]]]]]}",
      "[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[[1]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]",
      "{\"a\":/*x*/1}", "{\"a\"/*x*/:1}", "{/*x*/\"a\":1}", "{\"a\":1/*x*/}", "{\"a\":1//x\n}",
      "{\"a\":1#x\n,\"b\":2}", "{\"a\":1\r\n,\"b\":x\r\n}", "{\"a\":'it''s'}", "{\"a\":'\\''}",
      "{'a\"':\"b'\"}", "{\"a\":b/c}", "{\"a\":b\\c}", "{\"a\":b;c}", "{\"a\":b#c}", "{\"a\":b=c}",
      "{\"a\":\f1}", "{\"a\":1\f}", "{\f}", "\f{}",
    };
    for (String s : fixed) {
      in.add(s);
    }
    // Buffer boundary cases (JsonReader's buffer is 1024 chars).
    for (int pad : new int[] {1017, 1019, 1020, 1021, 1022, 1023, 1024, 2047}) {
      StringBuilder p = new StringBuilder();
      for (int i = 0; i < pad; i++) {
        p.append(' ');
      }
      in.add(p + "{\"a\":\"x\\u00e9\\ny\",\"b\":true,\"c\":null,\"d\":123456,\"e\":-1.5e3,f:g,'h':'i'}");
      in.add("{\"k\":\"" + p.toString().replace(' ', 'x') + "\\u0041\\n\"}");
      in.add("[" + p.toString().replace(' ', '1') + "]");
      in.add("[" + p.toString().replace(' ', 'z') + "]");
      in.add("{\"k\":" + p.toString().replace(' ', '7') + ".5}");
      in.add(p + "/* c */ {\"a\":1}");
      in.add(p + "{\"a\":1} x");
      in.add("{\"k\":\"" + p.toString().replace(' ', 'x') + "\\u12");
      in.add("{\"k\":\"" + p.toString().replace(' ', 'x') + "\\");
      in.add(p + ")]}'\n{}");
      in.add("\ufeff" + p + "{\"a\":true}");
      in.add("[" + p.toString().replace(' ', 't') + "rue]");
      in.add(p + "{\"a\":nul");
      in.add(p + "{\"a\":tru");
      in.add(p + "{\"a\":1}// c");
    }
    // Random token soups.
    String[] tokens = {
      "{", "}", "[", "]", ":", ",", ";", "=", "=>", "\"a\"", "'b'", "c", "1", "-1.5e3", "0", "01",
      "true", "null", "NaN", " ", "\n", "/*x*/", "//x\n", "#x\n", "\"\\u00e9\"", "\"\\\\\"", "\"x",
      "'", "\"", "\\", "/", "\f", "\r", "TRUE", "-", ".", "e", "\ud800", "\u2028",
    };
    Random r = new Random(20260415);
    for (int n = 0; n < 700; n++) {
      StringBuilder sb = new StringBuilder();
      int len = 1 + r.nextInt(14);
      if (r.nextInt(3) > 0) {
        sb.append(r.nextBoolean() ? "{" : "[");
      }
      for (int i = 0; i < len; i++) {
        sb.append(tokens[r.nextInt(tokens.length)]);
      }
      in.add(sb.toString());
    }
    // Random well-formed documents, then truncated or with one char replaced.
    for (int n = 0; n < 200; n++) {
      String doc = randomValue(r, 0);
      in.add(doc);
      if (doc.length() > 1) {
        in.add(doc.substring(0, r.nextInt(doc.length())));
        int at = r.nextInt(doc.length());
        in.add(doc.substring(0, at) + tokens[r.nextInt(tokens.length)] + doc.substring(at + 1));
      }
    }
    return in;
  }

  static String randomValue(Random r, int depth) {
    int k = depth > 3 ? r.nextInt(5) : r.nextInt(8);
    switch (k) {
      case 0:
        return String.valueOf(r.nextInt(2000) - 1000);
      case 1:
        return "\"" + (char) ('a' + r.nextInt(26)) + (r.nextBoolean() ? "\\n\\u00e9" : "") + "\"";
      case 2:
        return r.nextBoolean() ? "true" : "false";
      case 3:
        return "null";
      case 4:
        return (r.nextInt(100) - 50) + "." + r.nextInt(100) + (r.nextBoolean() ? "e" + (r.nextInt(20) - 10) : "");
      case 5:
      case 6:
        {
          StringBuilder sb = new StringBuilder("{");
          int n = r.nextInt(4);
          for (int i = 0; i < n; i++) {
            if (i > 0) {
              sb.append(r.nextInt(5) == 0 ? " , " : ",");
            }
            sb.append("\"k").append(r.nextInt(5)).append("\":").append(randomValue(r, depth + 1));
          }
          return sb.append("}").toString();
        }
      default:
        {
          StringBuilder sb = new StringBuilder("[");
          int n = r.nextInt(4);
          for (int i = 0; i < n; i++) {
            if (i > 0) {
              sb.append(",");
            }
            sb.append(randomValue(r, depth + 1));
          }
          return sb.append("]").toString();
        }
    }
  }

  public static void main(String[] args) throws Exception {
    try (PrintWriter out =
        new PrintWriter(
            new OutputStreamWriter(
                new GZIPOutputStream(new FileOutputStream(args[0])), StandardCharsets.UTF_8))) {
      for (String input : inputs()) {
        out.print("J\t" + esc(input));
        JsonElement[] holder = new JsonElement[1];
        out.print("\t" + esc(run(() -> holder[0] = GSON.fromJson(input, JsonElement.class))));
        out.print("\t" + esc(run(() -> GSON.fromJson(input, JsonObject.class))));
        out.print("\t" + esc(run(() -> GSON.fromJson(input, JsonArray.class))));
        out.println();
        JsonElement e = holder[0];
        if (e != null) {
          out.println("V\t\t" + esc(accessors(e)));
          if (e.isJsonObject()) {
            for (Map.Entry<String, JsonElement> m : e.getAsJsonObject().entrySet()) {
              out.println("V\t" + esc(m.getKey()) + "\t" + esc(accessors(m.getValue())));
            }
          } else if (e.isJsonArray()) {
            for (JsonElement m : e.getAsJsonArray()) {
              out.println("V\t[]\t" + esc(accessors(m)));
            }
          }
        }
      }
      String[] numbers = {
        "0", "1", "-1", "1.5", "-2.7", "1e2", "1E+2", "1e-2", "2147483648", "-2147483649",
        "9223372036854775807", "9223372036854775808", "-9223372036854775809", "1e19", "1e20",
        "12345678901234567890123", "-12345678901234567890123", "1.5e10", "-0.0", "0e5", "-0",
        "1e2147483647", "1e2147483648", "1e-2147483648", "1e-2147483649", "1e99999999999",
        "1e00000000000001", "0.000001e6", "123.456e-1", "1e-64", "1e63", "1e64", "-1e64", "5e63",
        "3e9", "-3e9", "4294967295.9", "4294967296", "18446744073709551616", "18446744073709551617.5",
        "0.0000", "100e-2", "7e-1", "1.e1", ".5", "1..2", "e5", "1e", "1e+", "", "-", "+1", "+",
        "1e5x", "abc", "\u0661\u0662", "1\u0661", "1e\u0661",
        "1e2147483649", "1.5e2147483648", "1e9999999999", "1e-9999999999", "1.5e-2147483647",
        "1e+2147483648", "0e-2147483649", "1.5e-2147483646", "10e-2147483647",
      };
      for (String s : numbers) {
        LazilyParsedNumber n = new LazilyParsedNumber(s);
        out.println("N\t" + esc(s) + "\t" + esc(run(() -> n.intValue())) + "\t" + esc(run(() -> n.longValue())));
      }
    }
  }
}
