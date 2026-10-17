// Vector generator for jutil.FileSystem_GetPathMatcher (black-box JDK output):
// glob patterns, the paths each matches, or the exception.
// Output: internal/jutil/testdata/glob_vectors.tsv
import java.nio.file.*;
public final class JutilGlobVectors {
  public static void main(String[] a) {
    String[] globs={"*.js","**.js","./**.js","./*/a.js","a?.js","[ab].js","[!ab].js","[a-c].js","[-a].js","[a-].js","{a,b}.js","a\\*.js","*","**","./src/**_test.js","[a/b]","{a,{b}}","[abc","{a,b","a\\","[]].js","[!].js","[z-a]","x[*?\\\\]y","a.b"};
    String[] paths={"a.js","b.js","c.js","d.js","-.js","./a.js","./x/a.js","x/y/a.js","ab.js","a*.js","./src/x_test.js","./src/d/y_test.js","a/b","]].js","].js","x*y","x?y","x\\y","axb",".js"};
    for (String g : globs) {
      PathMatcher m;
      try { m = FileSystems.getDefault().getPathMatcher("glob:"+g); }
      catch (Exception e) { System.out.println(g+"\tERR\t"+e.toString().replace("\n","\\n")); continue; }
      StringBuilder sb=new StringBuilder();
      for (String p : paths) if (m.matches(Paths.get(p))) sb.append(p).append(' ');
      System.out.println(g+"\tOK\t"+sb.toString().trim());
    }
  }
}
