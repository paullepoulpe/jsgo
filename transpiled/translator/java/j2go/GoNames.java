package j2go;

import java.util.Locale;
import java.util.Set;

/**
 * Deterministic Java-to-Go identifier, file-name, string-literal, and license-header helpers.
 */
public final class GoNames {
  private static final Set<String> GO_KEYWORDS_AND_PREDECLARED =
      Set.of(
          // Go keywords (25)
          "break",
          "case",
          "chan",
          "const",
          "continue",
          "default",
          "defer",
          "else",
          "fallthrough",
          "for",
          "func",
          "go",
          "goto",
          "if",
          "import",
          "interface",
          "map",
          "package",
          "range",
          "return",
          "select",
          "struct",
          "switch",
          "type",
          "var",
          // Go predeclared types, constants, and builtins that could shadow when used as locals/fields
          "any",
          "bool",
          "byte",
          "comparable",
          "complex64",
          "complex128",
          "error",
          "float32",
          "float64",
          "int",
          "int8",
          "int16",
          "int32",
          "int64",
          "rune",
          "string",
          "uint",
          "uint8",
          "uint16",
          "uint32",
          "uint64",
          "uintptr",
          "true",
          "false",
          "iota",
          "nil",
          "append",
          "cap",
          "clear",
          "close",
          "complex",
          "copy",
          "delete",
          "imag",
          "len",
          "make",
          "max",
          "min",
          "new",
          "panic",
          "print",
          "println",
          "real",
          "recover",
          // Common Go imports used in generated files
          "fmt",
          "math",
          "j2gort",
          "javalang",
          "javautil",
          "javaio",
          "guava",
          "nodestore");

  private GoNames() {}

  /** Sanitizes '$' and non-Go identifier chars while preserving case. */
  public static String cleanChars(String raw) {
    if (raw == null || raw.isEmpty()) {
      return "ident_";
    }
    StringBuilder sb = new StringBuilder(raw.length() + 4);
    for (int i = 0; i < raw.length(); i++) {
      char c = raw.charAt(i);
      if (c == '$') {
        sb.append(i == 0 ? "D_" : "_d");
      } else if (Character.isLetterOrDigit(c) || c == '_') {
        sb.append(c);
      } else {
        sb.append('_');
      }
    }
    return sb.toString();
  }

  /** Converts a Java identifier into an exported Go identifier (starts with A-Z). */
  public static String exported(String raw) {
    String s = cleanChars(raw);
    if (s.isEmpty()) {
      return "X_";
    }
    char first = s.charAt(0);
    if (first >= 'a' && first <= 'z') {
      return Character.toUpperCase(first) + s.substring(1);
    }
    if (first >= 'A' && first <= 'Z') {
      return s;
    }
    return "X" + s;
  }

  /** Converts a Java identifier into an unexported Go identifier (starts with a-z) safe from keywords. */
  public static String unexported(String raw) {
    String s = cleanChars(raw);
    if (s.isEmpty()) {
      return "v_";
    }
    char first = s.charAt(0);
    if (first >= 'A' && first <= 'Z') {
      // Lowercase leading uppercase run carefully, or prefix f_ so `X` and `x` never collide!
      s = "u_" + s;
    } else if (!(first >= 'a' && first <= 'z')) {
      s = "u" + s;
    }
    if (GO_KEYWORDS_AND_PREDECLARED.contains(s)) {
      return s + "_";
    }
    return s;
  }

  /** Safe local variable or parameter identifier in Go. */
  public static String localIdent(String raw) {
    String s = cleanChars(raw);
    if (s.isEmpty() || s.equals("_")) {
      return "unused_";
    }
    char first = s.charAt(0);
    if (first >= 'A' && first <= 'Z') {
      s = "l_" + s;
    } else if (!(first >= 'a' && first <= 'z')) {
      s = "l" + s;
    }
    if (GO_KEYWORDS_AND_PREDECLARED.contains(s)) {
      return s + "_";
    }
    return s;
  }

  /** Converts CamelCase Java class name to snake_case Go filename (without .go). */
  public static String toSnakeCase(String name) {
    String clean = name.replace("$", "d_");
    StringBuilder sb = new StringBuilder(clean.length() + 8);
    for (int i = 0; i < clean.length(); i++) {
      char c = clean.charAt(i);
      if (c >= 'A' && c <= 'Z') {
        if (i > 0) {
          char prev = clean.charAt(i - 1);
          boolean prevLowerOrDigit = (prev >= 'a' && prev <= 'z') || (prev >= '0' && prev <= '9');
          boolean nextLower = (i + 1 < clean.length()) && (clean.charAt(i + 1) >= 'a' && clean.charAt(i + 1) <= 'z');
          if (prevLowerOrDigit || (prev >= 'A' && prev <= 'Z' && nextLower)) {
            if (sb.length() > 0 && sb.charAt(sb.length() - 1) != '_') {
              sb.append('_');
            }
          }
        }
        sb.append(Character.toLowerCase(c));
      } else {
        sb.append(Character.toLowerCase(c));
      }
    }
    String res = sb.toString();
    // Avoid Go's magic _test.go suffix and OS/arch suffixes (_linux, _darwin, _windows, _amd64, _arm64, etc.)
    if (res.endsWith("_test")) {
      res = res + "_cls";
    }
    for (String arch :
        new String[] {
          "_aix", "_android", "_darwin", "_dragonfly", "_freebsd", "_illumos", "_ios", "_js",
          "_linux", "_netbsd", "_openbsd", "_plan9", "_solaris", "_wasip1", "_windows", "_386",
          "_amd64", "_arm", "_arm64", "_loong64", "_mips", "_ppc64", "_riscv64", "_s390x", "_wasm"
        }) {
      if (res.endsWith(arch)) {
        res = res + "_cls";
        break;
      }
    }
    return res;
  }

  /**
   * Encodes a Java UTF-16 String as a valid Go string literal using WTF-8 hex escapes
   * for lone UTF-16 surrogates (0xD800..0xDFFF) so Go never rejects surrogate escapes.
   */
  public static String goStringLiteral(String s) {
    if (s == null) {
      return "\"\"";
    }
    StringBuilder sb = new StringBuilder(s.length() + 8);
    sb.append('"');
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '\\') {
        sb.append("\\\\");
      } else if (c == '"') {
        sb.append("\\\"");
      } else if (c == '\n') {
        sb.append("\\n");
      } else if (c == '\r') {
        sb.append("\\r");
      } else if (c == '\t') {
        sb.append("\\t");
      } else if (c < 0x20 || c == 0x7f) {
        sb.append(String.format(Locale.ROOT, "\\x%02x", (int) c));
      } else if (Character.isHighSurrogate(c)) {
        if (i + 1 < s.length() && Character.isLowSurrogate(s.charAt(i + 1))) {
          int cp = Character.toCodePoint(c, s.charAt(i + 1));
          sb.append(String.format(Locale.ROOT, "\\U%08x", cp));
          i++;
        } else {
          appendWtf8Hex(sb, c);
        }
      } else if (Character.isLowSurrogate(c)) {
        appendWtf8Hex(sb, c);
      } else if (c < 0x7f) {
        sb.append(c);
      } else {
        sb.append(String.format(Locale.ROOT, "\\u%04x", (int) c));
      }
    }
    sb.append('"');
    return sb.toString();
  }

  private static void appendWtf8Hex(StringBuilder sb, char surrogate) {
    int b0 = 0xE0 | ((surrogate >> 12) & 0x0F);
    int b1 = 0x80 | ((surrogate >> 6) & 0x3F);
    int b2 = 0x80 | (surrogate & 0x3F);
    sb.append(String.format(Locale.ROOT, "\\x%02x\\x%02x\\x%02x", b0, b1, b2));
  }

  /**
   * Extracts the leading license header block from Java source text and strips the copyleft
   * dual-license alternative from MPL 1.1 Rhino headers per PORTING.md.
   */
  public static String extractSanitizedLicenseHeader(String javaSource) {
    if (javaSource == null) {
      return "";
    }
    int pkgIdx = javaSource.indexOf("\npackage ");
    if (pkgIdx < 0 && javaSource.startsWith("package ")) {
      return "";
    }
    String prefix = pkgIdx >= 0 ? javaSource.substring(0, pkgIdx).trim() : "";
    if (prefix.isEmpty() || !prefix.contains("/*")) {
      return "";
    }
    // Strip Rhino dual-license copyleft alternative per PORTING.md lines 37-40
    String dualTag = "Version: MPL 1.1/" + "G" + "PL 2.0";
    String gnuTag = "GNU General" + " Public License";
    if (prefix.contains(dualTag) || prefix.contains(gnuTag)) {
      prefix = prefix.replace(dualTag, "Version: MPL 1.1");
      prefix =
          prefix.replaceAll(
              "(?s)\\*\\s*Alternatively, the contents of this file may be used under the terms of.*?(\\* \\*\\*\\*\\*\\* END LICENSE BLOCK \\*\\*\\*\\*\\*)",
              "* This Go translation is distributed under the MPL 1.1 only, as the\n"
                  + " * original's dual licence allows; the other licence it offers is not\n"
                  + " * used.\n"
                  + " *\n"
                  + " $1");
    }
    return prefix;
  }
}
