// Original setup helper against pinned stock Closure. No compiler replacement.
package com.google.javascript.jscomp;

import com.google.gson.GsonBuilder;
import com.google.javascript.jscomp.deps.DependencyInfo;
import java.nio.charset.StandardCharsets;
import java.nio.file.*;
import java.security.MessageDigest;
import java.util.*;

/** Extract original AST dependencies before any module rewriting/optimization. */
public final class CliDependencyInventory {
  private static String sha(Path file) throws Exception {
    return HexFormat.of().formatHex(MessageDigest.getInstance("SHA-256").digest(Files.readAllBytes(file)));
  }
  public static void main(String[] args) throws Exception {
    Path cg = Path.of("/sys/fs/cgroup");
    String memory = Files.readString(cg.resolve("memory.max")).trim();
    String swap = Files.readString(cg.resolve("memory.swap.max")).trim();
    String cpu = Files.readString(cg.resolve("cpu.max")).trim();
    String cpuset = Files.readString(cg.resolve("cpuset.cpus.effective")).trim();
    String[] quota = cpu.split("\\s+");
    if (!memory.equals("4294967296") || !swap.equals("0") || quota[0].equals("max")
        || Long.parseLong(quota[0]) != 2*Long.parseLong(quota[1]) || !cpuset.equals("0-1")) {
      throw new IllegalStateException("requires real4GiB,swap0,CPU2,affinity0-1 setup container");
    }
    List<String> names = Files.readAllLines(Path.of(args[0]), StandardCharsets.UTF_8);
    if (names.isEmpty() || new HashSet<>(names).size() != names.size()) {
      throw new IllegalArgumentException("need distinct nonempty original input paths");
    }
    List<SourceFile> sources = new ArrayList<>();
    for (String name : names) sources.add(SourceFile.fromFile(name, StandardCharsets.UTF_8));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(CompilerOptions.LanguageMode.ECMASCRIPT_NEXT);
    options.setDependencyOptions(DependencyOptions.sortOnly());
    Compiler compiler = new Compiler();
    compiler.init(List.of(), sources, options);
    compiler.initializeModuleLoader();
    // This setup helper asks the stock AST DepsFinder used by SORT_ONLY.
    // It is not a regex heuristic and does not run optimization passes.
    compiler.setPreferRegexParser(false);
    Map<String, CompilerInput> byName = new HashMap<>();
    for (CompilerInput input : compiler.getInputsById().values()) byName.put(input.getName(), input);
    List<Map<String, Object>> files = new ArrayList<>();
    for (String name : names) {
      CompilerInput input = Objects.requireNonNull(byName.get(name), name);
      DependencyInfo info = input.getDependencyInfo();
      Map<String, Object> row = new LinkedHashMap<>();
      row.put("path", name);
      row.put("bytes", Files.size(Path.of(name)));
      row.put("sha256", sha(Path.of(name)));
      row.put("provides", info.getProvides());
      List<String> requires = new ArrayList<>();
      for (DependencyInfo.Require dependency : info.getRequires()) requires.add(dependency.getSymbol());
      row.put("requires", requires);
      row.put("type_requires", info.getTypeRequires());
      row.put("load_flags", info.getLoadFlags());
      row.put("externs_annotation", info.getHasExternsAnnotation());
      row.put("nocompile_annotation", info.getHasNoCompileAnnotation());
      files.add(row);
    }
    Path jar = Path.of(Compiler.class.getProtectionDomain().getCodeSource().getLocation().toURI());
    Map<String, Object> result = new LinkedHashMap<>();
    result.put("schema", "jsgo-stock-ast-dependencies-v1");
    result.put("stock_jar_sha256", sha(jar));
    result.put("method", "stock CompilerInput AST DepsFinder, SORT_ONLY, before rewrites; preparation only");
    result.put("kernel_budget", Map.of("memory.max", memory, "memory.swap.max", swap, "cpu.max", cpu, "cpuset.cpus.effective", cpuset));
    result.put("files", files);
    result.put("errors", compiler.getErrors().stream().map(Object::toString).toList());
    result.put("warnings", compiler.getWarnings().stream().map(Object::toString).toList());
    System.out.println(new GsonBuilder().setPrettyPrinting().create().toJson(result));
    if (!compiler.getErrors().isEmpty()) System.exit(1);
  }
}
