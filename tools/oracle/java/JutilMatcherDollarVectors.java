// Oracle tool: writes Java reference vectors for jutil.Matcher with the
// trailing-`$` support of jutil.Pattern (java.util.regex) into
// internal/jutil/testdata/matcher_dollar.tsv (gzip it to matcher_dollar.tsv.gz).
//   java -cp .cache/oracle-build JutilMatcherDollarVectors internal/jutil/testdata
//
// Each line is tab-separated; strings are written as "u" + 4-hex-digit UTF-16
// units:
//   find      <pattern> <input> <match>;<match>;...   (every find() result)
//   matches   <pattern> <input> <match> | -
//   first     <pattern> <input> <replacement> <result>   (replaceFirst)
//   all       <pattern> <input> <replacement> <result>   (replaceAll)
//   charup    <hex char list of Character.isUpperCase(char) == true>
// where a match is start,end,g1start:g1end,... (-1:-1 for a group that did
// not participate).
import java.io.PrintStream;
import java.nio.file.*;
import java.util.*;
import java.util.regex.*;

public class JutilMatcherDollarVectors {
  static String units(String s) {
    StringBuilder sb = new StringBuilder("u");
    for (int i = 0; i < s.length(); i++) sb.append(String.format("%04x", (int) s.charAt(i)));
    return sb.toString();
  }

  static String match(Matcher m) {
    StringBuilder sb = new StringBuilder();
    sb.append(m.start()).append(',').append(m.end());
    for (int g = 1; g <= m.groupCount(); g++) {
      sb.append(',').append(m.start(g)).append(':').append(m.end(g));
    }
    return sb.toString();
  }

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    String[] patterns = {
      ":-?\\d+",
      "blaze-out/[^/]*/(bin|genfiles)/",
      "^(.*)/(?:test|tests|testing)/(?:[^/]+)$",
      "-out/.*/(bin|genfiles)/(.*)$",
      "[A-Z0-9][A-Z0-9_]*",
      "a*",
      "(a)|(b)",
      "x$",
      "a|b$",
      "(\\d+)-(\\d+)",
      "[^/]+$",
      ".*$",
      "(x)?y",
      "^a",
      "^x|y",
      "^$",
    };
    String[] inputs = {
      "",
      "a",
      "x",
      "x\n",
      "x\r\n",
      "x\r",
      "x ",
      "x\n\n",
      "foo/test/bar.js",
      "foo/test/bar.js\n",
      "foo/test/bar.js\r\n",
      "a/b/tests/c",
      "a/testing/b/c",
      "blaze-out/k8-fastbuild/genfiles/some/path/foo.js",
      "blaze-out/k8-opt/bin/x/y.js",
      "x/blaze-out/k/bin/blaze-out/k/genfiles/z",
      "file.js:12  warning",
      "file.js:-1 a:2 b:345",
      "FOO_BAR",
      "foo",
      "aaa",
      "baab",
      "12-34 56-78",
      "y xy",
      "éé/test/😀.js",
      "b\n",
      "b\r\n",
      "\r\n",
      "aa",
      "xyxy",
      "yx",
    };
    String[] replacements = {":", "", "[$1]", "\\$x", "$0$0"};
    try (PrintStream ps = new PrintStream(Files.newOutputStream(out.resolve("matcher_dollar.tsv")), true, "UTF-8")) {
      for (String p : patterns) {
        Pattern pat = Pattern.compile(p);
        for (String in : inputs) {
          Matcher m = pat.matcher(in);
          StringBuilder sb = new StringBuilder();
          int n = 0;
          while (m.find() && n++ < 50) {
            if (sb.length() > 0) sb.append(';');
            sb.append(match(m));
          }
          ps.println("find\t" + units(p) + "\t" + units(in) + "\t" + sb);
          m = pat.matcher(in);
          ps.println("matches\t" + units(p) + "\t" + units(in) + "\t" + (m.matches() ? match(m) : "-"));
          for (String r : replacements) {
            String first, all;
            try {
              first = units(pat.matcher(in).replaceFirst(r));
            } catch (RuntimeException e) {
              first = "!" + e.getClass().getName();
            }
            try {
              all = units(pat.matcher(in).replaceAll(r));
            } catch (RuntimeException e) {
              all = "!" + e.getClass().getName();
            }
            ps.println("first\t" + units(p) + "\t" + units(in) + "\t" + units(r) + "\t" + first);
            ps.println("all\t" + units(p) + "\t" + units(in) + "\t" + units(r) + "\t" + all);
          }
        }
      }
      StringBuilder up = new StringBuilder();
      for (int c = 0; c < 0x10000; c++) {
        if (Character.isUpperCase((char) c)) {
          if (up.length() > 0) up.append(',');
          up.append(String.format("%04x", c));
        }
      }
      ps.println("charup\t" + up);
    }
  }
}
