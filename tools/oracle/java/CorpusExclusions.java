// Oracle helper: keeps the corpus files listed in tools/oracle/gpl-excluded-corpus.txt out of
// every dump made over a corpus. Their licence notices are not allowed in this repository
// (PORTING.md), and goldens made from them would hold their code. The Go side is
// tools/oracle/corpus; "go run ./tools/licenses -check" fails if a tracked file mentions one of
// them (see tools/oracle/README.md).
//
// Each tool that reads corpus files passes its file list, from the command line or from an
// @<listfile> / --list=<listfile>, through removeFrom() before dumping anything.
//
// The list is read from -Djsgo.excludedCorpus=<path> if set, else from
// tools/oracle/gpl-excluded-corpus.txt in the current directory or the nearest parent that has
// it. Without the list the tools stop: an excluded file must never get through.
//
// A path names an excluded file when, with forward slashes, it equals the file's key or ends with
// "/" and the key. The key is the listed path inside its npm package: the part after
// "<name>@<version>/package/" (the npm tarball root), or the whole path when there is no such
// part. So the match holds whatever directory the corpus root is.

package jsgo.oracle;

import java.io.IOException;
import java.io.UncheckedIOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Iterator;
import java.util.List;

public final class CorpusExclusions {
  /** The list, relative to the repository root. */
  public static final String LIST = "tools/oracle/gpl-excluded-corpus.txt";

  private static List<String> keys;

  private CorpusExclusions() {}

  /** The keys of the listed files, read once. */
  private static synchronized List<String> keys() {
    if (keys == null) {
      List<String> k = new ArrayList<>();
      for (String l : readList()) {
        l = l.trim();
        if (!l.isEmpty() && !l.startsWith("#")) {
          k.add(key(l));
        }
      }
      keys = k;
    }
    return keys;
  }

  private static List<String> readList() {
    String override = System.getProperty("jsgo.excludedCorpus");
    try {
      if (override != null) {
        return Files.readAllLines(Paths.get(override), StandardCharsets.UTF_8);
      }
      for (Path dir = Paths.get("").toAbsolutePath(); dir != null; dir = dir.getParent()) {
        Path p = dir.resolve(LIST);
        if (Files.isRegularFile(p)) {
          return Files.readAllLines(p, StandardCharsets.UTF_8);
        }
      }
    } catch (IOException e) {
      throw new UncheckedIOException(e);
    }
    throw new IllegalStateException(
        LIST + " not found above the current directory; run from the repository or set"
            + " -Djsgo.excludedCorpus=<path>");
  }

  /** The part of a listed path that names the file whatever the corpus root. */
  public static String key(String listed) {
    String[] parts = listed.split("/", -1);
    for (int i = 0; i + 2 < parts.length; i++) {
      if (parts[i].contains("@") && parts[i + 1].equals("package")) {
        return String.join("/", Arrays.asList(parts).subList(i + 2, parts.length));
      }
    }
    return listed;
  }

  /** Whether path (a corpus-relative or absolute file name) names an excluded file. */
  public static boolean isExcluded(String path) {
    String p = path.replace(java.io.File.separatorChar, '/');
    for (String k : keys()) {
      if (p.equals(k) || p.endsWith("/" + k)) {
        return true;
      }
    }
    return false;
  }

  /**
   * Removes the excluded files from files (a mutable list, or a sublist view of one), keeping the
   * order of the others, and says on stderr how many it removed.
   */
  public static void removeFrom(List<String> files) {
    int removed = 0;
    for (Iterator<String> it = files.iterator(); it.hasNext(); ) {
      if (isExcluded(it.next())) {
        it.remove();
        removed++;
      }
    }
    if (removed > 0) {
      System.err.println("CorpusExclusions: skipped " + removed + " file(s) listed in " + LIST);
    }
  }
}
