// Generates internal/jutil/testdata/deflater_vectors.txt: vectors for java.util.zip's Deflater,
// DeflaterOutputStream and GZIPOutputStream (the JDK's public API only), replayed by
// internal/jutil/deflater_vectors_test.go.
//   java -cp .cache/oracle-build JutilDeflaterVectors > internal/jutil/testdata/deflater_vectors.txt
//
// jutil compresses with Go's compress/flate and the JDK with zlib, so the compressed bytes may
// differ; what the vectors pin down is the stream contract: each case's output must decompress to
// exactly its input (the Go side checks its own output; this side checks the JDK's while
// generating), plus the GZIP header, the exceptions and the API behaviour after end().
//
// Each line is a case spec, " => ", and the result: "in=<n> sha=<SHA-256 of the input, hex>" (the
// Go side regenerates the input and checks it against these), or "exc=<class>: <message>", or the
// GZIP header as "hex=...". The spec is space-separated key=value pairs that both sides interpret
// the same way:
//   kind=deflater|gzip|gzipsync
//   input=<shape>:<size>:<seed>   the uncompressed data (see input() below; shapes rand, text,
//                                 words, runs, mix, zeros)
//   level=<n> strategy=<n> nowrap=0|1   (deflater) new Deflater(level, nowrap), then
//                                 setStrategy(strategy) before any input
//   prelevel=<n>                  (gzip) def.setLevel(n) right after construction, as Closure's
//                                 FastGzipOutputStream does
//   chunks=all|fixed:<n>|rand:<seed>:<max>   how the data is split into setInput/write calls
//   out=<n>                       (deflater) the output buffer of each deflate call; (gzip) the
//                                 stream's buffer size
//   mid=<offset>:<level>:<strategy>   (deflater) setLevel/setStrategy once <offset> bytes have
//                                 been fed (at the end of the chunk that reaches it)
// The deflater driver is: for each chunk, setInput(chunk) then deflate(buf) until needsInput();
// then finish() and deflate(buf) until finished(). The gzip driver writes each chunk (gzipsync
// also calls flush() after each, after which everything written so far must be decodable) and
// closes the stream.
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.zip.Deflater;
import java.util.zip.GZIPInputStream;
import java.util.zip.GZIPOutputStream;
import java.util.zip.Inflater;
import java.util.zip.InflaterInputStream;

public class JutilDeflaterVectors {

  /** SplitMix64 (Steele, Lea and Flood, "Fast splittable pseudorandom number generators"). */
  static final class Rng {
    long state;

    Rng(long seed) {
      state = seed;
    }

    long next() {
      long z = (state += 0x9e3779b97f4a7c15L);
      z = (z ^ (z >>> 30)) * 0xbf58476d1ce4e5b9L;
      z = (z ^ (z >>> 27)) * 0x94d049bb133111ebL;
      return z ^ (z >>> 31);
    }

    /** A value in [0, n). */
    int below(int n) {
      return (int) Long.remainderUnsigned(next(), n);
    }
  }

  static byte[] input(String spec) {
    String[] p = spec.split(":");
    String shape = p[0];
    int size = Integer.parseInt(p[1]);
    Rng r = new Rng(Long.parseLong(p[2]));
    byte[] b = new byte[size];
    switch (shape) {
      case "rand" -> {
        for (int i = 0; i < size; i++) {
          b[i] = (byte) r.next();
        }
      }
      case "text" -> {
        byte[] alphabet = "abcdefgh \n".getBytes(StandardCharsets.US_ASCII);
        for (int i = 0; i < size; i++) {
          b[i] = alphabet[r.below(alphabet.length)];
        }
      }
      case "words", "mix" -> {
        // A vocabulary of 300 words of 1 to 10 letters, then words joined by separators; for mix,
        // segments of random bytes are interleaved.
        Rng vr = new Rng(12345);
        byte[][] vocab = new byte[300][];
        for (int w = 0; w < vocab.length; w++) {
          vocab[w] = new byte[1 + vr.below(10)];
          for (int k = 0; k < vocab[w].length; k++) {
            vocab[w][k] = (byte) ('a' + vr.below(26));
          }
        }
        byte[] seps = " , . ;\n(){}=".getBytes(StandardCharsets.US_ASCII);
        int i = 0;
        while (i < size) {
          if (shape.equals("mix") && r.below(8) == 0) {
            int n = 1 + r.below(40000);
            for (int k = 0; k < n && i < size; k++) {
              b[i++] = (byte) r.next();
            }
            continue;
          }
          // Zipf-like: favour the first words.
          int w = r.below(1 + r.below(vocab.length));
          for (int k = 0; k < vocab[w].length && i < size; k++) {
            b[i++] = vocab[w][k];
          }
          if (i < size) {
            b[i++] = seps[r.below(seps.length)];
          }
        }
      }
      case "runs" -> {
        int i = 0;
        while (i < size) {
          byte v = (byte) r.below(4);
          int n = 1 + r.below(300);
          for (int k = 0; k < n && i < size; k++) {
            b[i++] = v;
          }
        }
      }
      case "zeros" -> {}
      default -> throw new IllegalArgumentException(shape);
    }
    return b;
  }

  static List<int[]> chunks(String spec, int size) {
    List<int[]> out = new ArrayList<>();
    String[] p = spec.split(":");
    switch (p[0]) {
      case "all" -> out.add(new int[] {0, size});
      case "fixed" -> {
        int n = Integer.parseInt(p[1]);
        for (int i = 0; i < size; i += n) {
          out.add(new int[] {i, Math.min(n, size - i)});
        }
      }
      case "rand" -> {
        Rng r = new Rng(Long.parseLong(p[1]));
        int max = Integer.parseInt(p[2]);
        for (int i = 0; i < size; ) {
          int n = Math.min(1 + r.below(max), size - i);
          out.add(new int[] {i, n});
          i += n;
        }
      }
      default -> throw new IllegalArgumentException(spec);
    }
    if (out.isEmpty()) {
      out.add(new int[] {0, 0});
    }
    return out;
  }

  static Map<String, String> parse(String spec) {
    Map<String, String> m = new LinkedHashMap<>();
    for (String kv : spec.split(" ")) {
      int eq = kv.indexOf('=');
      m.put(kv.substring(0, eq), kv.substring(eq + 1));
    }
    return m;
  }

  static byte[] run(String spec) throws IOException {
    Map<String, String> m = parse(spec);
    byte[] data = input(m.get("input"));
    List<int[]> chunks = chunks(m.getOrDefault("chunks", "all"), data.length);
    int outSize = Integer.parseInt(m.getOrDefault("out", "512"));
    ByteArrayOutputStream bos = new ByteArrayOutputStream();
    switch (m.get("kind")) {
      case "deflater" -> {
        Deflater d =
            new Deflater(
                Integer.parseInt(m.getOrDefault("level", "-1")), m.getOrDefault("nowrap", "0").equals("1"));
        int strategy = Integer.parseInt(m.getOrDefault("strategy", "0"));
        if (strategy != 0) {
          d.setStrategy(strategy);
        }
        long midAt = -1;
        int midLevel = 0;
        int midStrategy = 0;
        if (m.containsKey("mid")) {
          String[] mp = m.get("mid").split(":");
          midAt = Long.parseLong(mp[0]);
          midLevel = Integer.parseInt(mp[1]);
          midStrategy = Integer.parseInt(mp[2]);
        }
        byte[] buf = new byte[outSize];
        long fed = 0;
        for (int[] c : chunks) {
          d.setInput(data, c[0], c[1]);
          while (!d.needsInput()) {
            int n = d.deflate(buf);
            bos.write(buf, 0, n);
          }
          fed += c[1];
          if (midAt >= 0 && fed >= midAt) {
            d.setLevel(midLevel);
            d.setStrategy(midStrategy);
            midAt = -1;
          }
        }
        d.finish();
        while (!d.finished()) {
          int n = d.deflate(buf);
          bos.write(buf, 0, n);
        }
        d.end();
      }
      case "gzip", "gzipsync" -> {
        boolean sync = m.get("kind").equals("gzipsync");
        int prelevel = m.containsKey("prelevel") ? Integer.parseInt(m.get("prelevel")) : -2;
        GZIPOutputStream g =
            new GZIPOutputStream(bos, outSize, sync) {
              {
                if (prelevel != -2) {
                  def.setLevel(prelevel);
                }
              }
            };
        for (int[] c : chunks) {
          g.write(data, c[0], c[1]);
          if (sync) {
            g.flush();
          }
        }
        g.close();
      }
      default -> throw new IllegalArgumentException(spec);
    }
    return bos.toByteArray();
  }

  static String hex(byte[] b) {
    StringBuilder sb = new StringBuilder();
    for (byte x : b) {
      sb.append(String.format("%02x", x & 0xff));
    }
    return sb.toString();
  }

  static List<String> specs() {
    List<String> s = new ArrayList<>();
    // Level 1 (Closure's FastGzipOutputStream), raw deflate: every shape at sizes around block and
    // window boundaries, fed in one chunk.
    int[] sizes = {
      0, 1, 2, 3, 4, 10, 100, 257, 258, 259, 262, 4096, 16383, 16384, 32767, 32768, 32769, 65273,
      65274, 65275, 65276, 65535, 65536, 65537, 98304, 100000, 300000
    };
    String[] shapes = {"rand", "text", "words", "runs", "mix", "zeros"};
    for (String shape : shapes) {
      for (int size : sizes) {
        s.add("kind=deflater level=1 nowrap=1 input=" + shape + ":" + size + ":1 out=65536");
      }
    }
    // Large inputs.
    for (String shape : new String[] {"words", "mix", "text", "rand"}) {
      s.add("kind=deflater level=1 nowrap=1 input=" + shape + ":1000000:2 out=65536");
    }
    // Chunked input and small output buffers.
    for (String shape : new String[] {"words", "mix", "runs"}) {
      for (String chunks : new String[] {"fixed:1", "fixed:7", "fixed:4096", "fixed:65535", "rand:5:9000", "rand:6:70000"}) {
        s.add(
            "kind=deflater level=1 nowrap=1 input="
                + shape
                + ":200000:3 chunks="
                + chunks
                + " out=65536");
      }
      for (int out : new int[] {1, 2, 7, 512, 4096}) {
        s.add("kind=deflater level=1 nowrap=1 input=" + shape + ":200000:4 out=" + out);
      }
    }
    // Every level, both wrappers.
    for (int level = 1; level <= 9; level++) {
      for (String shape : new String[] {"words", "mix", "text", "runs"}) {
        s.add("kind=deflater level=" + level + " nowrap=1 input=" + shape + ":150000:8 out=65536");
      }
      s.add("kind=deflater level=" + level + " nowrap=0 input=words:150000:9 out=512");
      s.add("kind=deflater level=" + level + " nowrap=0 input=rand:300:9 out=512");
      s.add(
          "kind=deflater level="
              + level
              + " nowrap=1 input=words:300000:10 chunks=rand:11:20000 out=1000");
    }
    s.add("kind=deflater level=-1 nowrap=0 input=words:0:1 out=512");
    s.add("kind=deflater level=-1 nowrap=0 input=words:100000:1 out=512");
    s.add("kind=deflater level=6 nowrap=0 input=zeros:1000000:1 out=512");
    // Strategies.
    for (int level : new int[] {1, 4, 6, 9}) {
      for (int strategy : new int[] {1, 2}) {
        for (String shape : new String[] {"words", "mix"}) {
          s.add(
              "kind=deflater level="
                  + level
                  + " strategy="
                  + strategy
                  + " nowrap=1 input="
                  + shape
                  + ":150000:12 out=65536");
        }
      }
    }
    // setLevel/setStrategy in the middle of the stream.
    for (String mid :
        new String[] {"50000:6:0", "50000:1:0", "50000:2:0", "50000:9:1", "50000:1:2", "0:6:0"}) {
      for (int level : new int[] {1, 6}) {
        for (String chunks : new String[] {"all", "fixed:20000"}) {
          s.add(
              "kind=deflater level="
                  + level
                  + " nowrap=1 input=words:150000:13 chunks="
                  + chunks
                  + " mid="
                  + mid
                  + " out=65536");
        }
      }
    }
    s.add("kind=deflater level=1 nowrap=1 input=words:150000:14 chunks=fixed:20000 mid=60000:6:0 out=100");
    // GZIPOutputStream, default level and as FastGzipOutputStream (BEST_SPEED, 64K buffer).
    for (String shape : new String[] {"words", "mix", "rand", "text"}) {
      for (int size : new int[] {0, 1, 100, 70000, 400000}) {
        s.add("kind=gzip input=" + shape + ":" + size + ":15 out=512");
        s.add("kind=gzip prelevel=1 input=" + shape + ":" + size + ":15 out=65536");
      }
      s.add("kind=gzip prelevel=1 input=" + shape + ":400000:16 chunks=fixed:4096 out=65536");
      s.add("kind=gzip prelevel=1 input=" + shape + ":400000:16 chunks=rand:17:10000 out=65536");
      s.add("kind=gzip prelevel=9 input=" + shape + ":100000:16 out=10");
    }
    s.add("kind=gzipsync input=words:100000:18 chunks=fixed:30000 out=512");
    s.add("kind=gzipsync prelevel=1 input=mix:200000:18 chunks=rand:19:50000 out=4096");
    return s;
  }

  /** Decompresses a case's output. */
  static byte[] inflate(String spec, byte[] b) throws IOException {
    Map<String, String> m = parse(spec);
    if (m.get("kind").equals("deflater")) {
      boolean nowrap = m.getOrDefault("nowrap", "0").equals("1");
      try (InflaterInputStream in =
          new InflaterInputStream(new ByteArrayInputStream(b), new Inflater(nowrap))) {
        return in.readAllBytes();
      }
    }
    try (GZIPInputStream in = new GZIPInputStream(new ByteArrayInputStream(b))) {
      return in.readAllBytes();
    }
  }

  static String exc(Throwable e) {
    return e.getClass().getName() + ": " + e.getMessage();
  }

  public static void main(String[] args) throws Exception {
    PrintStream out = new PrintStream(System.out, true, StandardCharsets.UTF_8);
    MessageDigest sha = MessageDigest.getInstance("SHA-256");
    for (String spec : specs()) {
      byte[] b = run(spec);
      byte[] data = input(parse(spec).get("input"));
      if (!java.util.Arrays.equals(inflate(spec, b), data)) {
        throw new AssertionError("does not decompress to its input: " + spec);
      }
      out.println(spec + " => in=" + data.length + " sha=" + hex(sha.digest(data)));
    }
    // Exceptions.
    try {
      new Deflater(10, true);
    } catch (RuntimeException e) {
      out.println("exc=newlevel10 => exc=" + exc(e));
    }
    try {
      new Deflater(1, true).setLevel(-2);
    } catch (RuntimeException e) {
      out.println("exc=setlevel-2 => exc=" + exc(e));
    }
    try {
      new Deflater(1, true).setStrategy(3);
    } catch (RuntimeException e) {
      out.println("exc=setstrategy3 => exc=" + exc(e));
    }
    try {
      new Deflater(1, true).deflate(new byte[10], 0, 10, 1);
    } catch (RuntimeException e) {
      out.println("exc=flush1 => exc=" + exc(e));
    }
    try {
      Deflater d = new Deflater(1, true);
      d.end();
      // setInput, setLevel, setStrategy, needsInput, finish and finished still work.
      d.setInput(new byte[3]);
      d.setLevel(5);
      d.needsInput();
      d.finished();
      d.deflate(new byte[10]);
    } catch (RuntimeException e) {
      out.println("exc=ended => exc=" + exc(e));
    }
    try {
      GZIPOutputStream g = new GZIPOutputStream(new ByteArrayOutputStream(), 0);
    } catch (RuntimeException e) {
      out.println("exc=gzipsize0 => exc=" + exc(e));
    }
    try {
      GZIPOutputStream g = new GZIPOutputStream(new ByteArrayOutputStream());
      g.finish();
      g.write(1);
    } catch (IOException e) {
      out.println("exc=writeafterfinish => exc=" + exc(e));
    }
    // The header is written by the constructor.
    ByteArrayOutputStream hb = new ByteArrayOutputStream();
    new GZIPOutputStream(hb);
    out.println("header => hex=" + hex(hb.toByteArray()));
  }
}
