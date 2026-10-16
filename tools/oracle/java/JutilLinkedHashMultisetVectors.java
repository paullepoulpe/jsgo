// Generates internal/jutil/testdata/linked_hash_multiset.tsv: seeded scripts
// of add/remove operations on Guava's LinkedHashMultiset<String> (the shaded
// copy in the pinned jar), with every result.
//   java -cp .cache/closure-compiler.jar:.cache/oracle-build JutilLinkedHashMultisetVectors internal/jutil/testdata
//
// Rows: "<ops>\t<results>", where <ops> is ";"-separated "a<word>" (add) or
// "r<word>" (remove) and <results> is "|"-separated, one per op:
// "<op result>,<size>,<count of word>,<contains word>,<toString>,<iteration>"
// with the iteration as the ":"-joined elements. Words are single ASCII
// letters, so no escaping is needed.
import com.google.javascript.jscomp.jarjar.com.google.common.collect.LinkedHashMultiset;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.*;
import java.util.ArrayList;
import java.util.List;
import java.util.Random;

public class JutilLinkedHashMultisetVectors {
  static final String[] WORDS = {"a", "b", "c", "d", "e", "f"};

  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    Random r = new Random(7);
    try (PrintStream ps =
        new PrintStream(
            Files.newOutputStream(out.resolve("linked_hash_multiset.tsv")),
            false,
            StandardCharsets.UTF_8)) {
      for (int k = 0; k < 400; k++) {
        LinkedHashMultiset<String> m = LinkedHashMultiset.create();
        int n = r.nextInt(16);
        List<String> ops = new ArrayList<>();
        List<String> results = new ArrayList<>();
        for (int i = 0; i < n; i++) {
          String w = WORDS[r.nextInt(WORDS.length)];
          boolean add = r.nextInt(3) != 0;
          boolean res = add ? m.add(w) : m.remove(w);
          ops.add((add ? "a" : "r") + w);
          results.add(
              res + "," + m.size() + "," + m.count(w) + "," + m.contains(w) + "," + m + ","
                  + String.join(":", m));
        }
        ps.println(String.join(";", ops) + "\t" + String.join("|", results));
      }
    }
  }
}
