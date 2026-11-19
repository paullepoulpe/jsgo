package j2go;

import java.io.BufferedOutputStream;
import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.io.DataInputStream;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Enumeration;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import java.util.jar.JarEntry;
import java.util.jar.JarOutputStream;
import java.util.zip.ZipEntry;
import java.util.zip.ZipFile;

/**
 * Prepares an unshaded dependency JAR from Closure Compiler's shaded Maven uber-jar
 * so {@code JavacTask} resolves all third-party libraries (Guava, Gson, Protobuf,
 * RE2J, args4j, AutoValue, ErrorProne, JSpecify) and compiled proto classes with
 * zero javac error diagnostics.
 *
 * <p>Operates strictly on Closure Compiler's own Apache-2.0/BSD/MIT jar using the
 * public JVMS §4.4 classfile constant-pool specification (100% clean-room).
 */
public final class ClasspathUnshader {
  private static final String JARJAR_PREFIX = "com/google/javascript/jscomp/jarjar/";

  private ClasspathUnshader() {}

  public static Path ensureUnshadedJar(Path shadedJar, List<Path> sourceFiles, Path outJar)
      throws IOException {
    Set<String> srcTopInternalNames = collectSourceInternalNames(sourceFiles);
    Files.createDirectories(outJar.getParent());

    // Rebuild if missing or older than shadedJar
    if (Files.exists(outJar)
        && Files.size(outJar) > 0
        && Files.getLastModifiedTime(outJar).compareTo(Files.getLastModifiedTime(shadedJar)) >= 0) {
      return outJar;
    }

    Path tmpJar = outJar.resolveSibling(outJar.getFileName().toString() + ".tmp");
    Set<String> seenEntries = new HashSet<>();
    try (ZipFile zf = new ZipFile(shadedJar.toFile());
        JarOutputStream jos =
            new JarOutputStream(new BufferedOutputStream(Files.newOutputStream(tmpJar)))) {
      Enumeration<? extends ZipEntry> entries = zf.entries();
      while (entries.hasMoreElements()) {
        ZipEntry entry = entries.nextElement();
        if (entry.isDirectory()) {
          continue;
        }
        String name = entry.getName();
        if (name.startsWith(JARJAR_PREFIX)) {
          name = name.substring(JARJAR_PREFIX.length());
        }
        if (!name.endsWith(".class")) {
          continue;
        }
        // Skip XML APIs provided by OpenJDK's standard java.xml module
        if (name.startsWith("javax/xml/")
            || name.startsWith("org/xml/")
            || name.startsWith("org/w3c/")) {
          continue;
        }
        String internalName = name.substring(0, name.length() - 6);
        String topInternal = topLevelInternalName(internalName);
        if (srcTopInternalNames.contains(topInternal)) {
          continue;
        }
        if (!seenEntries.add(name)) {
          continue;
        }

        byte[] raw;
        try (InputStream is = zf.getInputStream(entry)) {
          raw = is.readAllBytes();
        }
        byte[] unshaded = unshadeConstantPool(raw);
        JarEntry je = new JarEntry(name);
        je.setTime(0L);
        jos.putNextEntry(je);
        jos.write(unshaded);
        jos.closeEntry();
      }
    }
    Files.move(tmpJar, outJar, java.nio.file.StandardCopyOption.REPLACE_EXISTING);
    return outJar;
  }

  private static String topLevelInternalName(String internalName) {
    int autoIdx = internalName.indexOf("/$AutoValue_");
    if (autoIdx >= 0) {
      int nestedDollar = internalName.indexOf('$', autoIdx + 2);
      return nestedDollar >= 0 ? internalName.substring(0, nestedDollar) : internalName;
    }
    int dollar = internalName.indexOf('$');
    return dollar >= 0 ? internalName.substring(0, dollar) : internalName;
  }

  private static Set<String> collectSourceInternalNames(List<Path> sourceFiles) throws IOException {
    Set<String> names = new HashSet<>();
    for (Path p : sourceFiles) {
      String content = Files.readString(p, StandardCharsets.UTF_8);
      String pkg = "";
      int start = 0;
      while (start < content.length()) {
        int nl = content.indexOf('\n', start);
        if (nl < 0) nl = content.length();
        String line = content.substring(start, nl).trim();
        if (line.startsWith("package ") && line.endsWith(";")) {
          pkg = line.substring(8, line.length() - 1).trim().replace('.', '/');
          break;
        }
        start = nl + 1;
      }
      String base = p.getFileName().toString();
      if (base.endsWith(".java")) {
        base = base.substring(0, base.length() - 5);
      }
      names.add(pkg.isEmpty() ? base : pkg + "/" + base);
    }
    return names;
  }

  private static byte[] unshadeConstantPool(byte[] raw) throws IOException {
    DataInputStream in = new DataInputStream(new ByteArrayInputStream(raw));
    ByteArrayOutputStream bos = new ByteArrayOutputStream(raw.length);
    DataOutputStream out = new DataOutputStream(bos);

    int magic = in.readInt();
    if (magic != 0xCAFEBABE) {
      return raw;
    }
    out.writeInt(magic);
    out.writeShort(in.readUnsignedShort()); // minor_version
    out.writeShort(in.readUnsignedShort()); // major_version
    int cpCount = in.readUnsignedShort();
    out.writeShort(cpCount);

    for (int i = 1; i < cpCount; i++) {
      int tag = in.readUnsignedByte();
      out.writeByte(tag);
      switch (tag) {
        case 1 -> { // CONSTANT_Utf8
          int len = in.readUnsignedShort();
          byte[] bytes = new byte[len];
          in.readFully(bytes);
          String s = new String(bytes, StandardCharsets.UTF_8);
          if (s.contains(JARJAR_PREFIX)) {
            s = s.replace(JARJAR_PREFIX, "");
            bytes = s.getBytes(StandardCharsets.UTF_8);
          }
          out.writeShort(bytes.length);
          out.write(bytes);
        }
        case 7, 8, 16, 19, 20 -> out.writeShort(in.readUnsignedShort());
        case 15 -> {
          out.writeByte(in.readUnsignedByte());
          out.writeShort(in.readUnsignedShort());
        }
        case 3, 4, 9, 10, 11, 12, 17, 18 -> out.writeInt(in.readInt());
        case 5, 6 -> {
          out.writeLong(in.readLong());
          i++; // 8-byte constants occupy two constant-pool entries (JVMS §4.4.5)
        }
        default -> throw new IOException("Unsupported constant-pool tag: " + tag);
      }
    }
    out.write(in.readAllBytes());
    out.flush();
    return bos.toByteArray();
  }
}
