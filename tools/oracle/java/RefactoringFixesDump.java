// Oracle tool for com.google.javascript.refactoring (SuggestedFix, ErrorToFixMapper,
// ApplySuggestedFixes, FixingErrorManager): runs the linter's fixes over each input file, the way
// Linter.fixRepeatedly(String) does, in memory, and dumps every suggested fix of every round and
// the code each round produces.
//
// Each file is linted on its own. A round is Linter.fix(filename, ImmutableSet.of()) with the
// file's current code in place of the file on disk: a fresh Compiler (output to a null stream)
// with a fresh FixingErrorManager (no unfixable errors) as its error manager, setCompiler, then
// Linter.lint(path, compiler): LintPassConfig over Linter.Builder()'s options (read from a fresh
// Linter.Builder by reflection, so that they are exactly the linter's), disableThreads, and
// compile([SourceFile.fromCode("<Linter externs>", "")], [SourceFile.fromCode(<rel-file>, code)],
// options). The round then takes FixingErrorManager.getSureFixes(); if there are none, the file
// is done; otherwise ApplySuggestedFixes.applySuggestedFixesToCode(sureFixes, {<rel-file>: code})
// gives the code of the next round (Linter writes it back to the file). At most
// Linter.MAX_FIXES (5) rounds run.
//
// It lives in package com.google.javascript.jscomp to reach the package-private LintPassConfig,
// Linter.builder() and Linter.MAX_FIXES.
//
// Build/run (see tools/oracle/README.md):
//   javac -cp closure-compiler.jar -d build RefactoringFixesDump.java
//   java -Xss1g -cp closure-compiler.jar:build com.google.javascript.jscomp.RefactoringFixesDump
//        [--gzip] <outfile> <root> <rel-file>... | @<listfile>
//
// Format (one block per file, in input order):
//   # file <rel-file>
//   ## round <n>                            n = 0, 1, ...
//   fix <i> alts=<k> desc=<description|null> at=<sourceFilename>:<lineno>:<charno>
//       closurized=<true|false>               one line per FixingErrorManager.getAllFixes()
//                                           entry, in its order (errors, then warnings); "at=null"
//                                           when no MatchedNodeInfo is attached
//     <filename> <CodeReplacement.toString()>
//                                           one line per replacement, in
//                                           getReplacements().entries() order
//   ## sure <count>                         the number of getSureFixes()
//   ## code len=<length> sha256=<hex>       the code after applying them (length in UTF-16
//                                           units; SHA-256 of its UTF-8 bytes)
//   ## exception <class> [<- <cause class>]...   if the round threw (the file is then done)
//   ## final len=<length> sha256=<hex>      the code after the last round
// Every line is escaped: '\' -> "\\", and every char outside printable ASCII -> \\uXXXX
// (newlines included), so the output is ASCII.
package com.google.javascript.jscomp;

import com.google.javascript.refactoring.ApplySuggestedFixes;
import com.google.javascript.refactoring.CodeReplacement;
import com.google.javascript.refactoring.FixingErrorManager;
import com.google.javascript.refactoring.SuggestedFix;
import java.io.OutputStream;
import java.io.PrintStream;
import java.lang.reflect.Field;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.Collection;
import java.util.List;
import java.util.Map;
import java.util.zip.GZIPOutputStream;
import jsgo.oracle.CorpusExclusions;

public final class RefactoringFixesDump {
  public static void main(String[] args) throws Exception {
    boolean gzip = false;
    List<String> rest = new ArrayList<>();
    for (String a : args) {
      if (a.equals("--gzip")) {
        gzip = true;
      } else {
        rest.add(a);
      }
    }
    Path out = Paths.get(rest.get(0));
    Path root = Paths.get(rest.get(1));
    List<String> files = new ArrayList<>();
    for (String f : rest.subList(2, rest.size())) {
      if (f.startsWith("@")) {
        for (String l : Files.readAllLines(Paths.get(f.substring(1)), StandardCharsets.UTF_8)) {
          if (!l.isEmpty()) {
            files.add(l);
          }
        }
      } else {
        files.add(f);
      }
    }
    CorpusExclusions.removeFrom(files); // tools/oracle/gpl-excluded-corpus.txt
    if (out.getParent() != null) {
      Files.createDirectories(out.getParent());
    }
    try (OutputStream os =
        gzip ? new GZIPOutputStream(Files.newOutputStream(out)) : Files.newOutputStream(out)) {
      for (String rel : files) {
        StringBuilder sb = new StringBuilder();
        dumpOne(sb, rel, Files.readString(root.resolve(rel), StandardCharsets.UTF_8));
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
      }
    }
  }

  /** The options of a fresh Linter.Builder (Linter.Builder.options). */
  private static CompilerOptions linterOptions() throws ReflectiveOperationException {
    Linter.Builder builder = Linter.builder();
    Field f = Linter.Builder.class.getDeclaredField("options");
    f.setAccessible(true);
    return (CompilerOptions) f.get(builder);
  }

  private static void dumpOne(StringBuilder sb, String name, String src) throws Exception {
    line(sb, "# file " + name);
    String code = src;
    for (int round = 0; round < Linter.MAX_FIXES; round++) {
      line(sb, "## round " + round);
      try {
        CompilerOptions options = linterOptions();
        Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
        FixingErrorManager errorManager = new FixingErrorManager();
        compiler.setErrorManager(errorManager);
        errorManager.setCompiler(compiler);

        // Linter.lint(Path, Compiler), with the code in memory.
        compiler.setPassConfig(new LintPassConfig(options));
        compiler.disableThreads();
        SourceFile externs = SourceFile.fromCode("<Linter externs>", "");
        compiler.compile(List.of(externs), List.of(SourceFile.fromCode(name, code)), options);

        int i = 0;
        for (SuggestedFix fix : errorManager.getAllFixes()) {
          dumpFix(sb, i++, fix);
        }
        Collection<SuggestedFix> fixes = errorManager.getSureFixes();
        line(sb, "## sure " + fixes.size());
        if (fixes.isEmpty()) {
          break;
        }
        // ApplySuggestedFixes.applySuggestedFixesToFiles writes only the files that have
        // replacements.
        String newCode =
            ApplySuggestedFixes.applySuggestedFixesToCode(fixes, Map.of(name, code)).get(name);
        if (newCode != null) {
          code = newCode;
        }
        line(sb, "## code " + digest(code));
      } catch (RuntimeException | StackOverflowError e) {
        StringBuilder ex = new StringBuilder("## exception " + e.getClass().getName());
        for (Throwable c = e.getCause(); c != null; c = c.getCause()) {
          ex.append(" <- ").append(c.getClass().getName());
        }
        line(sb, ex.toString());
        break;
      }
    }
    line(sb, "## final " + digest(code));
  }

  private static void dumpFix(StringBuilder sb, int i, SuggestedFix fix) {
    SuggestedFix.MatchedNodeInfo info = fix.getMatchedNodeInfo();
    String at =
        info == null
            ? "null"
            : info.sourceFilename() + ":" + info.lineno() + ":" + info.charno()
                + " closurized=" + info.inClosurizedFile();
    line(
        sb,
        "fix " + i + " alts=" + fix.getAlternatives().size() + " desc=" + fix.getDescription()
            + " at=" + at);
    for (Map.Entry<String, CodeReplacement> e : fix.getReplacements().entries()) {
      line(sb, "  " + e.getKey() + " " + e.getValue());
    }
  }

  private static String digest(String code) throws Exception {
    MessageDigest md = MessageDigest.getInstance("SHA-256");
    byte[] h = md.digest(code.getBytes(StandardCharsets.UTF_8));
    StringBuilder sb = new StringBuilder("len=" + code.length() + " sha256=");
    for (byte b : h) {
      sb.append(String.format("%02x", b & 0xff));
    }
    return sb.toString();
  }

  private static void line(StringBuilder sb, String s) {
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '\\') {
        sb.append("\\\\");
      } else if (c >= 0x20 && c < 0x7f) {
        sb.append(c);
      } else {
        sb.append(String.format("\\u%04x", (int) c));
      }
    }
    sb.append('\n');
  }

  private RefactoringFixesDump() {}
}
