// Oracle tool: writes black-box reference vectors for the charsets of
// internal/jutil (charset_registry.go) into <dir>/charset_vectors.tsv.gz.
// Only the public API of the JDK is called: Charset.forName, isSupported,
// availableCharsets, name, aliases, new String(byte[], Charset),
// InputStreamReader, CharsetDecoder.decode (REPORT), CharsetEncoder.canEncode
// and OutputStreamWriter.
//   javac -d .cache/oracle-build tools/oracle/java/JutilCharsetVectors.java
//   java -cp .cache/oracle-build JutilCharsetVectors internal/jutil/testdata 20261017
// The Go side is internal/jutil/charset_vectors_test.go;
// JUTIL_CHARSET_VECTORS=<dir> points it at vectors written with another seed.
//
// Rows (tab-separated; strings are written as their UTF-16 units, four hex
// digits each, bytes as hex):
//   cs      <canonical name> <aliases, sorted, comma-separated>
//           one row for every charset the JDK has;
//   name    <name> <Charset.forName(name).name() or exception> <isSupported or exception>
//           every canonical name and alias in random case, and random names;
//   dec     <charset> <bytes> <new String(bytes, cs)> <InputStreamReader result>
//           <REPORT result: the units, or M<length> / U<length>>
//           random byte strings, valid encodings of random text with random
//           damage, and charset-specific fragments (escape sequences, marks);
//   enc     <charset> <string, its writes separated by '|'> <OutputStreamWriter bytes>
//   can     <charset> <string> <canEncode of each char, '0'/'1', on one encoder>
//   digest  <charset> <plan> <SHA-256 of the exhaustive lines of the plan>
//           the plan names input sets, each giving one line per input:
//             N     all N-byte inputs (N = 1 or 2);
//             P+N   the prefix P (hex) followed by all N-byte suffixes;
//             N+S   all N-byte inputs followed by the suffix S (hex);
//             gb4   all four-byte GB18030 sequences 81-FE 30-39 81-FE 30-39;
//           with the line  <bytes hex> \t <new String units> \t <REPORT result> \n;
//             bmp   every char c: e<c, 4 hex> \t <canEncode on a fresh encoder: 0/1>
//                   \t <OutputStreamWriter bytes of the one-char string> \n;
//             sup   every 257th code point from U+10000: s<cp, 6 hex> \t
//                   <OutputStreamWriter bytes of its surrogate pair> \n.
//           The lines are UTF-8 (ASCII) text; the digest covers them in plan order.
import java.io.*;
import java.nio.ByteBuffer;
import java.nio.charset.*;
import java.security.MessageDigest;
import java.util.*;
import java.util.zip.GZIPOutputStream;

public class JutilCharsetVectors {
  // The charsets jutil implements, with their exhaustive plans.
  static final String SB = "1 bmp sup";
  static final String DB = "1 2 bmp sup";
  static final String[][] CHARSETS = {
    {"US-ASCII", SB}, {"ISO-8859-1", SB}, {"UTF-8", "1 2 bmp sup"},
    {"UTF-16", DB}, {"UTF-16BE", DB}, {"UTF-16LE", DB},
    {"UTF-32", "1 2 0000+2 0010+2 2+0000 2+1000 bmp sup"}, {"UTF-32BE", "1 2 0000+2 0010+2 bmp sup"},
    {"UTF-32LE", "1 2 2+0000 2+1000 bmp sup"}, {"X-UTF-32BE-BOM", "1 2 0000+2 bmp sup"},
    {"X-UTF-32LE-BOM", "1 2 2+0000 bmp sup"},
    {"windows-1250", SB}, {"windows-1251", SB}, {"windows-1252", SB}, {"windows-1253", SB},
    {"windows-1254", SB}, {"windows-1255", SB}, {"windows-1256", SB}, {"windows-1257", SB},
    {"windows-1258", SB}, {"x-windows-874", SB},
    {"ISO-8859-2", SB}, {"ISO-8859-3", SB}, {"ISO-8859-4", SB}, {"ISO-8859-5", SB},
    {"ISO-8859-6", SB}, {"ISO-8859-7", SB}, {"ISO-8859-8", SB}, {"ISO-8859-9", SB},
    {"x-iso-8859-11", SB}, {"ISO-8859-13", SB}, {"ISO-8859-15", SB}, {"ISO-8859-16", SB},
    {"TIS-620", SB}, {"KOI8-R", SB}, {"KOI8-U", SB},
    {"IBM437", SB}, {"IBM850", SB}, {"IBM852", SB}, {"IBM855", SB}, {"IBM00858", SB},
    {"IBM860", SB}, {"IBM862", SB}, {"IBM863", SB}, {"IBM865", SB}, {"IBM866", SB},
    {"x-MacRoman", SB}, {"x-MacCyrillic", SB},
    {"Shift_JIS", DB}, {"windows-31j", DB}, {"EUC-JP", "1 2 8f+2 bmp sup"},
    {"ISO-2022-JP", "1 2 1b2442+1 1b2442+2 1b2849+1 1b284a+1 0e+1 1b+2 1b24+2 1b28+2 1b24420e+1 bmp sup"},
    {"GBK", DB}, {"GB2312", DB}, {"GB18030", "1 2 8130+2 gb4 bmp sup"},
    {"Big5", DB}, {"EUC-KR", DB}, {"x-windows-949", DB},
  };

  static Random r;

  static String hex(byte[] b) {
    StringBuilder sb = new StringBuilder();
    for (byte x : b) sb.append(String.format("%02x", x & 0xff));
    return sb.toString();
  }

  static byte[] unhex(String h) {
    byte[] b = new byte[h.length() / 2];
    for (int i = 0; i < b.length; i++) b[i] = (byte) Integer.parseInt(h.substring(2 * i, 2 * i + 2), 16);
    return b;
  }

  static String units(String s) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < s.length(); i++) sb.append(String.format("%04x", (int) s.charAt(i)));
    return sb.toString();
  }

  static String report(Charset cs, byte[] b) {
    try {
      return units(cs.newDecoder().decode(ByteBuffer.wrap(b)).toString());
    } catch (MalformedInputException e) {
      return "M" + e.getInputLength();
    } catch (UnmappableCharacterException e) {
      return "U" + e.getInputLength();
    } catch (CharacterCodingException e) {
      throw new RuntimeException(e);
    }
  }

  static String reader(Charset cs, byte[] b) throws IOException {
    Reader in = new InputStreamReader(new ByteArrayInputStream(b), cs);
    StringBuilder sb = new StringBuilder();
    char[] buf = new char[3];
    for (int n; (n = in.read(buf)) >= 0; ) sb.append(buf, 0, n);
    return sb.toString();
  }

  static byte[] write(Charset cs, String... writes) throws IOException {
    ByteArrayOutputStream bo = new ByteArrayOutputStream();
    Writer w = new OutputStreamWriter(bo, cs);
    for (String s : writes) w.write(s);
    w.close();
    return bo.toByteArray();
  }

  static String exc(Throwable e) {
    return e.toString();
  }

  // Random case for a name.
  static String mangle(String n) {
    StringBuilder sb = new StringBuilder();
    for (char c : n.toCharArray()) {
      int k = r.nextInt(3);
      sb.append(k == 0 ? Character.toUpperCase(c) : k == 1 ? Character.toLowerCase(c) : c);
    }
    return sb.toString();
  }

  static String nameRow(String n) {
    String a, b;
    try { a = Charset.forName(n).name(); } catch (Exception e) { a = exc(e); }
    try { b = String.valueOf(Charset.isSupported(n)); } catch (Exception e) { b = exc(e); }
    return "name\t" + n + "\t" + a + "\t" + b;
  }

  // The encodable BMP chars of a charset (fresh encoders), for random text.
  static char[] repertoire(Charset cs) {
    StringBuilder sb = new StringBuilder();
    for (int c = 0x80; c < 0x10000; c++) {
      if (c >= 0xD800 && c <= 0xDFFF) continue;
      if (cs.newEncoder().canEncode((char) c)) sb.append((char) c);
    }
    return sb.toString().toCharArray();
  }

  static String randomText(char[] rep, int len) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < len; i++) {
      int k = r.nextInt(20);
      if (k < 11 && rep.length > 0) sb.append(rep[r.nextInt(rep.length)]);
      else if (k < 15) sb.append((char) (0x20 + r.nextInt(0x5f)));
      else if (k < 16) sb.append((char) r.nextInt(0x20));
      else if (k < 17) sb.append((char) (0x80 + r.nextInt(0xff80)));
      else if (k < 18) sb.appendCodePoint(0x10000 + r.nextInt(0x100000));
      else if (k < 19) sb.append((char) (0xD800 + r.nextInt(0x800)));
      else sb.append("¥‾ｱ€﻿�".charAt(r.nextInt(6)));
    }
    return sb.toString();
  }

  static final String[] FRAGMENTS = {
    "1b2442", "1b2440", "1b2842", "1b284a", "1b2849", "1b2441", "1b2428", "1b28", "1b24", "1b", "0e", "0f",
    "0000feff", "fffe0000", "feff", "fffe", "0000d800", "0010ffff", "00110000", "efbbbf",
    "8130", "81308130", "8431a439", "90308130", "e3329a35", "e3329a36", "8e", "8f", "a1", "81", "80", "ff",
  };

  static byte[] damage(byte[] b) {
    ByteArrayOutputStream o = new ByteArrayOutputStream();
    for (int i = 0; i < b.length; i++) {
      int k = r.nextInt(40);
      if (k == 0) continue;
      if (k == 1) o.write(r.nextInt(256));
      if (k == 2) { o.write(r.nextInt(256)); continue; }
      if (k == 3) { byte[] f = unhex(FRAGMENTS[r.nextInt(FRAGMENTS.length)]); o.write(f, 0, f.length); }
      o.write(b[i]);
    }
    byte[] out = o.toByteArray();
    if (r.nextInt(4) == 0 && out.length > 0) out = Arrays.copyOf(out, r.nextInt(out.length));
    return out;
  }

  static byte[] randomBytes(String name, char[] rep) throws IOException {
    switch (r.nextInt(4)) {
      case 0: {
        byte[] b = new byte[r.nextInt(10)];
        r.nextBytes(b);
        return b;
      }
      case 1: {
        ByteArrayOutputStream o = new ByteArrayOutputStream();
        for (int n = r.nextInt(6); n > 0; n--) {
          if (r.nextBoolean()) {
            byte[] f = unhex(FRAGMENTS[r.nextInt(FRAGMENTS.length)]);
            o.write(f, 0, f.length);
          } else {
            o.write(r.nextInt(256));
          }
        }
        return o.toByteArray();
      }
      default:
        return damage(write(Charset.forName(name), randomText(rep, r.nextInt(8))));
    }
  }

  static void plan(Charset cs, String plan, MessageDigest md) throws IOException {
    for (String item : plan.split(" ")) {
      if (item.equals("bmp")) {
        for (int c = 0; c < 0x10000; c++) {
          String s = String.valueOf((char) c);
          line(md, String.format("e%04x\t%s\t%s\n", c, cs.newEncoder().canEncode((char) c) ? "1" : "0", hex(write(cs, s))));
        }
      } else if (item.equals("sup")) {
        for (int cp = 0x10000; cp <= 0x10FFFF; cp += 257) {
          line(md, String.format("s%06x\t%s\n", cp, hex(write(cs, new String(Character.toChars(cp))))));
        }
      } else if (item.equals("gb4")) {
        byte[] b = new byte[4];
        for (int b1 = 0x81; b1 <= 0xFE; b1++)
          for (int b2 = 0x30; b2 <= 0x39; b2++)
            for (int b3 = 0x81; b3 <= 0xFE; b3++)
              for (int b4 = 0x30; b4 <= 0x39; b4++) {
                b[0] = (byte) b1; b[1] = (byte) b2; b[2] = (byte) b3; b[3] = (byte) b4;
                decLine(cs, b, md);
              }
      } else {
        String pre = "", suf = "";
        int n;
        int plus = item.indexOf('+');
        if (plus < 0) {
          n = Integer.parseInt(item);
        } else if (plus == 1) {
          n = Integer.parseInt(item.substring(0, 1));
          suf = item.substring(2);
        } else {
          pre = item.substring(0, plus);
          n = Integer.parseInt(item.substring(plus + 1));
        }
        byte[] p = unhex(pre), s = unhex(suf);
        byte[] b = new byte[p.length + n + s.length];
        System.arraycopy(p, 0, b, 0, p.length);
        System.arraycopy(s, 0, b, p.length + n, s.length);
        for (int v = 0; v < 1 << (8 * n); v++) {
          for (int k = 0; k < n; k++) b[p.length + k] = (byte) (v >> (8 * (n - 1 - k)));
          decLine(cs, b, md);
        }
      }
    }
  }

  static void decLine(Charset cs, byte[] b, MessageDigest md) {
    line(md, hex(b) + "\t" + units(new String(b, cs)) + "\t" + report(cs, b) + "\n");
  }

  static PrintStream debug; // the plan lines of one charset, when asked

  static void line(MessageDigest md, String s) {
    md.update(s.getBytes(StandardCharsets.US_ASCII));
    if (debug != null) debug.print(s);
  }

  public static void main(String[] args) throws Exception {
    String dir = args[0];
    long seed = Long.parseLong(args[1]);
    String debugCharset = args.length > 3 && args[2].equals("--lines") ? args[3] : null;
    if (debugCharset != null) {
      debug = new PrintStream(new FileOutputStream(dir + "/" + debugCharset + ".lines"), false, "US-ASCII");
    }
    r = new Random(seed);
    PrintStream out = new PrintStream(new GZIPOutputStream(new FileOutputStream(dir + "/charset_vectors.tsv.gz")), false, "US-ASCII");
    out.println("seed\t" + seed);
    // Every JDK charset with its aliases; names in random case.
    List<String> names = new ArrayList<>();
    for (Charset c : Charset.availableCharsets().values()) {
      TreeSet<String> as = new TreeSet<>(c.aliases());
      out.println("cs\t" + c.name() + "\t" + String.join(",", as));
      names.add(c.name());
      names.addAll(as);
    }
    for (String n : names) {
      out.println(nameRow(n));
      out.println(nameRow(mangle(n)));
    }
    String alphabet = "abcdefgijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-+.:_ /@*";
    for (int i = 0; i < 300; i++) {
      StringBuilder sb = new StringBuilder();
      for (int n = r.nextInt(14); n >= 0; n--) sb.append(alphabet.charAt(r.nextInt(alphabet.length())));
      if (r.nextInt(3) == 0) sb.insert(0, names.get(r.nextInt(names.size())).substring(0, 1));
      out.println(nameRow(sb.toString()));
    }
    for (String[] e : CHARSETS) {
      String name = e[0];
      Charset cs = Charset.forName(name);
      char[] rep = repertoire(cs);
      for (int i = 0; i < 150; i++) {
        byte[] b = randomBytes(name, rep);
        out.println("dec\t" + name + "\t" + hex(b) + "\t" + units(new String(b, cs)) + "\t" + units(reader(cs, b)) + "\t" + report(cs, b));
      }
      for (int i = 0; i < 100; i++) {
        String s = randomText(rep, r.nextInt(10));
        List<String> writes = new ArrayList<>();
        int at = 0;
        while (at < s.length()) {
          int end = r.nextInt(3) == 0 ? s.length() : at + 1 + r.nextInt(s.length() - at);
          writes.add(s.substring(at, end));
          at = end;
        }
        if (r.nextInt(5) == 0) writes.add("");
        List<String> us = new ArrayList<>();
        for (String w : writes) us.add(units(w));
        out.println("enc\t" + name + "\t" + String.join("|", us) + "\t" + hex(write(cs, writes.toArray(new String[0]))));
      }
      for (int i = 0; i < 40; i++) {
        String s = randomText(rep, 1 + r.nextInt(12));
        CharsetEncoder enc = cs.newEncoder();
        StringBuilder bits = new StringBuilder();
        for (char c : s.toCharArray()) bits.append(enc.canEncode(c) ? '1' : '0');
        out.println("can\t" + name + "\t" + units(s) + "\t" + bits);
      }
      MessageDigest md = MessageDigest.getInstance("SHA-256");
      PrintStream saved = debug;
      if (!name.equals(debugCharset)) debug = null;
      plan(cs, e[1], md);
      debug = saved;
      out.println("digest\t" + name + "\t" + e[1] + "\t" + hex(md.digest()));
    }
    out.close();
    if (debug != null) debug.close();
  }
}
