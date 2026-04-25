package main

import (
	"strings"
	"testing"
)

// ── #166 Generics ────────────────────────────────────────────────────────────

func TestParseTypeParams_simple(t *testing.T) {
	tps := parseTypeParams("T")
	if len(tps) != 1 || tps[0].name != "T" || tps[0].bound != "any" {
		t.Fatalf("got %+v", tps)
	}
}

func TestParseTypeParams_bounded(t *testing.T) {
	tps := parseTypeParams("T extends Foo")
	if len(tps) != 1 || tps[0].name != "T" || tps[0].bound != "Foo" {
		t.Fatalf("got %+v", tps)
	}
}

func TestParseTypeParams_wildcard_extends(t *testing.T) {
	tps := parseTypeParams("? extends Foo")
	if len(tps) != 1 || tps[0].name != "T" || tps[0].bound != "Foo" {
		t.Fatalf("got %+v", tps)
	}
}

func TestParseTypeParams_wildcard_super(t *testing.T) {
	tps := parseTypeParams("? super Foo")
	if len(tps) != 1 || tps[0].name != "T" || tps[0].bound != "any" {
		t.Fatalf("got %+v", tps)
	}
}

func TestParseTypeParams_multiple(t *testing.T) {
	tps := parseTypeParams("K, V")
	if len(tps) != 2 {
		t.Fatalf("want 2 params, got %d: %+v", len(tps), tps)
	}
	if tps[0].name != "K" || tps[1].name != "V" {
		t.Fatalf("got %+v", tps)
	}
}

func TestTypeParamsGoString(t *testing.T) {
	tps := []typeParam{{name: "K", bound: "any"}, {name: "V", bound: "any"}}
	got := typeParamsGoString(tps)
	if got != "[K any, V any]" {
		t.Fatalf("got %q", got)
	}
}

func TestTypeParamsGoString_bounded(t *testing.T) {
	tps := []typeParam{{name: "T", bound: "Foo"}}
	got := typeParamsGoString(tps)
	if got != "[T Foo]" {
		t.Fatalf("got %q", got)
	}
}

func TestGenerics_emittedInStruct(t *testing.T) {
	src := `package com.example;
public class Pair<K, V> {
  public K key;
  public V value;
}`
	// parseClass should pick up typeParams
	cls := parseClass(src)
	if cls == nil {
		t.Fatal("parseClass returned nil")
	}
	if len(cls.typeParams) != 2 {
		t.Fatalf("want 2 type params, got %d: %+v", len(cls.typeParams), cls.typeParams)
	}
	out := emitStruct(cls)
	if !strings.Contains(out, "[K any, V any]") {
		t.Errorf("expected generic params in struct, got:\n%s", out)
	}
}

// ── #167 Collections ─────────────────────────────────────────────────────────

func TestJavaToGoType_ImmutableList(t *testing.T) {
	got := javaToGoType("ImmutableList<String>")
	if got != "[]string" {
		t.Fatalf("got %q", got)
	}
}

func TestJavaToGoType_ImmutableMap(t *testing.T) {
	got := javaToGoType("ImmutableMap<String, int>")
	if got != "map[string]int" {
		t.Fatalf("got %q", got)
	}
}

func TestJavaToGoType_ImmutableSet(t *testing.T) {
	got := javaToGoType("ImmutableSet<String>")
	if got != "map[string]struct{}" {
		t.Fatalf("got %q", got)
	}
}

func TestJavaToGoType_Optional(t *testing.T) {
	got := javaToGoType("Optional<String>")
	if got != "*string" {
		t.Fatalf("got %q", got)
	}
}

func TestJavaToGoType_OptionalPointer(t *testing.T) {
	// Optional<Node> → *rhino.Node (already a pointer, no double pointer)
	got := javaToGoType("Optional<Node>")
	if got != "*rhino.Node" {
		t.Fatalf("got %q", got)
	}
}

func TestJavaToGoType_ArrayList(t *testing.T) {
	got := javaToGoType("ArrayList<String>")
	if got != "[]string" {
		t.Fatalf("got %q", got)
	}
}

func TestJavaToGoType_HashMap(t *testing.T) {
	got := javaToGoType("HashMap<String, boolean>")
	if got != "map[string]bool" {
		t.Fatalf("got %q", got)
	}
}

// ── #163 Enum-with-methods ───────────────────────────────────────────────────

func TestEnumSimple(t *testing.T) {
	src := `package com.example;
public enum Color {
  RED, GREEN, BLUE
}`
	cls := parseClass(src)
	if cls == nil || !cls.isEnum {
		t.Fatal("expected enum class")
	}
	if len(cls.enumConsts) != 3 {
		t.Fatalf("want 3 consts, got %d: %+v", len(cls.enumConsts), cls.enumConsts)
	}
	out := emitEnum(cls, "com/example/Color.java")
	if !strings.Contains(out, "type Color int") {
		t.Errorf("expected iota enum, got:\n%s", out)
	}
	if !strings.Contains(out, "RED Color = iota") {
		t.Errorf("expected iota, got:\n%s", out)
	}
}

func TestEnumWithMethods(t *testing.T) {
	src := `package com.example;
public enum Kind {
  A {
    @Override
    public boolean foo() { return true; }
  },
  B {
    @Override
    public boolean foo() { return false; }
  };
  public abstract boolean foo();
}`
	cls := parseClass(src)
	if cls == nil || !cls.isEnum {
		t.Fatal("expected enum")
	}
	out := emitEnum(cls, "com/example/Kind.java")
	if !strings.Contains(out, "type Kind interface") {
		t.Errorf("expected interface, got:\n%s", out)
	}
	if !strings.Contains(out, "type Kind_A struct") {
		t.Errorf("expected Kind_A struct, got:\n%s", out)
	}
	if !strings.Contains(out, "type Kind_B struct") {
		t.Errorf("expected Kind_B struct, got:\n%s", out)
	}
	if !strings.Contains(out, "func NewKind(name string) Kind") {
		t.Errorf("expected factory, got:\n%s", out)
	}
	if !strings.Contains(out, "KindA Kind") {
		t.Errorf("expected sentinel var KindA, got:\n%s", out)
	}
}

// ── #168 AutoValue ───────────────────────────────────────────────────────────

func TestAutoValue(t *testing.T) {
	src := `package com.example;
/** A simple value class. */
@AutoValue
abstract class Point {
  abstract int x();
  abstract int y();
  static Point create(int x, int y) { return new AutoValue_Point(x, y); }
}`
	cls := parseClass(src)
	if cls == nil {
		t.Fatal("parseClass returned nil")
	}
	if !cls.isAutoValue {
		t.Fatal("expected isAutoValue=true")
	}
	out := emitAutoValue(cls, "com/example/Point.java")
	if !strings.Contains(out, "type Point struct") {
		t.Errorf("expected struct, got:\n%s", out)
	}
	if !strings.Contains(out, "func NewPoint(") {
		t.Errorf("expected factory, got:\n%s", out)
	}
	if !strings.Contains(out, "x int") && !strings.Contains(out, "X int") {
		t.Errorf("expected x field, got:\n%s", out)
	}
}

// ── #162 Inner classes ───────────────────────────────────────────────────────

func TestInnerClassParsing(t *testing.T) {
	src := `package com.example;
public class Outer {
  private int x;
  static class Inner {
    private String name;
  }
}`
	cls := parseClass(src)
	if cls == nil {
		t.Fatal("parseClass returned nil")
	}
	if len(cls.innerClasses) != 1 {
		t.Fatalf("want 1 inner class, got %d", len(cls.innerClasses))
	}
	if cls.innerClasses[0].name != "Inner" {
		t.Fatalf("got inner name %q", cls.innerClasses[0].name)
	}
	out := emitInnerClass("Outer", cls.innerClasses[0], "com/example/Outer.java")
	if !strings.Contains(out, "type OuterInner struct") {
		t.Errorf("expected OuterInner struct, got:\n%s", out)
	}
}

// ── #169 Panic stubs ─────────────────────────────────────────────────────────

func TestPanicStub(t *testing.T) {
	// Bodies that can be translated by the javaast package emit Go code
	// instead of a panic stub. `return 42;` translates to `return 42`.
	// Bodies that cannot be parsed at all still get the panic stub.
	src := `package com.example;
public class Foo {
  public int bar(String s) { return 42; }
}`
	cls := parseClass(src)
	if cls == nil {
		t.Fatal("parseClass returned nil")
	}
	out := emitMethods(cls, "com/example/Foo.java")
	// Translator successfully handles `return 42;` — no panic stub needed.
	if !strings.Contains(out, "return 42") {
		t.Errorf("expected translated 'return 42', got:\n%s", out)
	}
}

// ── #170 Source-line pins ─────────────────────────────────────────────────────

func TestSourceLinePins(t *testing.T) {
	src := `package com.example;
public class Foo {
  public void doIt() {}
  public int calc(int x) { return x; }
}`
	cls := parseClass(src)
	if cls == nil {
		t.Fatal("parseClass returned nil")
	}
	out := emitMethods(cls, "com/example/Foo.java")
	if !strings.Contains(out, "// @ com/example/Foo.java:") {
		t.Errorf("expected source-line pin comment, got:\n%s", out)
	}
}

// ── #172 Javadoc preservation ────────────────────────────────────────────────

func TestJavadocPreservation(t *testing.T) {
	src := `package com.example;
/**
 * Checks things.
 * @param x the value
 * @return the result
 */
public class Checker {
  /**
   * Processes a node.
   * @param n the node
   * @see OtherClass
   */
  public void process(Object n) {}
}`
	cls := parseClass(src)
	if cls == nil {
		t.Fatal("parseClass returned nil")
	}
	if !strings.Contains(cls.javadoc, "Checks things") {
		t.Errorf("class javadoc not preserved: %q", cls.javadoc)
	}
	if !strings.Contains(cls.javadoc, "Param: x the value") {
		t.Errorf("@param not converted, javadoc: %q", cls.javadoc)
	}
	if !strings.Contains(cls.javadoc, "Returns: the result") {
		t.Errorf("@return not converted, javadoc: %q", cls.javadoc)
	}
	out := emitMethods(cls, "com/example/Checker.java")
	if !strings.Contains(out, "// Processes a node") {
		t.Errorf("method javadoc missing, got:\n%s", out)
	}
	if !strings.Contains(out, "// Param: n the node") {
		t.Errorf("method @param not converted, got:\n%s", out)
	}
	if !strings.Contains(out, "// See: OtherClass") {
		t.Errorf("method @see not converted, got:\n%s", out)
	}
}

// ── #164 Import graph ─────────────────────────────────────────────────────────

func TestJavaImportParsing(t *testing.T) {
	src := `import com.google.javascript.jscomp.Compiler;
import com.google.javascript.rhino.Node;
import java.util.List;`
	imports := parseJavaImports(src)
	if len(imports) != 3 {
		t.Fatalf("want 3 imports, got %d: %v", len(imports), imports)
	}
}

func TestJavaImportsToGoPkgs(t *testing.T) {
	javaImports := []string{
		"com.google.javascript.jscomp.Compiler",
		"com.google.javascript.rhino.Node",
		"java.util.List",
	}
	pkgs := javaImportsToGoPkgs(javaImports)
	foundJscomp := false
	foundRhino := false
	for _, p := range pkgs {
		if strings.Contains(p, "internal/jscomp") {
			foundJscomp = true
		}
		if strings.Contains(p, "internal/rhino") {
			foundRhino = true
		}
	}
	if !foundJscomp {
		t.Error("expected jscomp package in imports")
	}
	if !foundRhino {
		t.Error("expected rhino package in imports")
	}
}

// ── #171 Cache ───────────────────────────────────────────────────────────────

func TestCacheRoundtrip(t *testing.T) {
	dir := t.TempDir()
	key := computeCacheKey([]byte("hello world"))
	if got := loadCache(dir, key); got != "" {
		t.Fatal("expected empty cache miss")
	}
	saveCache(dir, key, "generated output")
	got := loadCache(dir, key)
	if got != "generated output" {
		t.Fatalf("cache roundtrip failed, got %q", got)
	}
}

func TestCacheKeyDiffers(t *testing.T) {
	k1 := computeCacheKey([]byte("src1"))
	k2 := computeCacheKey([]byte("src2"))
	if k1 == k2 {
		t.Fatal("different inputs should yield different keys")
	}
}

// ── Full pipeline smoke test ─────────────────────────────────────────────────

func TestGenerateSkeleton_smoke(t *testing.T) {
	src := `package com.google.javascript.jscomp;
import com.google.javascript.rhino.Node;
/**
 * Does compilation.
 * @param n the node
 */
public final class MyPass implements CompilerPass {
  private final AbstractCompiler compiler;
  private int count;

  public MyPass(AbstractCompiler compiler) {
    this.compiler = compiler;
  }

  /**
   * Main entry point.
   * @param root the root node
   * @return whether succeeded
   */
  @Override
  public boolean process(Node root, Node externs) {
    return false;
  }
}`
	out := generateSkeleton(src, "com/google/javascript/jscomp/MyPass.java", "jscomp", "")
	checks := []string{
		"// Port of com/google/javascript/jscomp/MyPass.java",
		"package jscomp",
		"type MyPass struct",
		"func NewMyPass(",
		"func (c *MyPass) Process(",
		// Note: `return false;` is translated by javaast — no panic stub.
		"// @ com/google/javascript/jscomp/MyPass.java:",
		"// Main entry point",
		"// Param: root the root node",
		"// Returns: whether succeeded",
	}
	for _, c := range checks {
		if !strings.Contains(out, c) {
			t.Errorf("expected %q in output", c)
		}
	}
}

func TestGenerateSkeleton_innerClass(t *testing.T) {
	src := `package com.google.javascript.jscomp;
public class Outer {
  private int x;
  static class State {
    boolean active;
  }
  public void reset() {}
}`
	out := generateSkeleton(src, "com/google/javascript/jscomp/Outer.java", "jscomp", "")
	if !strings.Contains(out, "type OuterState struct") {
		t.Errorf("expected OuterState inner class, got:\n%s", out)
	}
}

func TestGenerateSkeleton_genericClass(t *testing.T) {
	src := `package com.google.javascript.jscomp;
public class Container<T extends Object> {
  T value;
  public T get() { return value; }
}`
	out := generateSkeleton(src, "com/google/javascript/jscomp/Container.java", "jscomp", "")
	if !strings.Contains(out, "[T") {
		t.Errorf("expected generic type param, got:\n%s", out)
	}
}

// ── Bug fix regressions ──────────────────────────────────────────────────────

// Bug #1 — Field with generic constructor initialiser (new ArrayList<>())
// must NOT be emitted as a method stub like `arrayList<>()`.
func TestBug1_GenericConstructorInitFieldNotMethod(t *testing.T) {
	src := `package com.google.javascript.jscomp;
import java.util.ArrayList;
import java.util.LinkedHashSet;
import com.google.javascript.rhino.Node;
import java.util.List;
import java.util.Set;
final class CheckSideEffects {
  private final List<Node> problemNodes = new ArrayList<>();
  private final Set<String> noSideEffectExterns = new LinkedHashSet<>();
  CheckSideEffects(boolean report) { this.report = report; }
  private final boolean report;
}`
	out := generateSkeleton(src, "com/google/javascript/jscomp/CheckSideEffects.java", "jscomp", "")
	// The generic constructor initialiser must NOT produce a method stub.
	if strings.Contains(out, "arrayList<>") || strings.Contains(out, "ArrayList<>") ||
		strings.Contains(out, "linkedHashSet<>") || strings.Contains(out, "LinkedHashSet<>") {
		t.Errorf("Bug #1: field with new ArrayList<>() emitted as method stub:\n%s", out)
	}
	// The struct must still contain the fields.
	if !strings.Contains(out, "ProblemNodes") && !strings.Contains(out, "problemNodes") {
		t.Errorf("Bug #1: struct missing ProblemNodes field:\n%s", out)
	}
}

// Bug #2 — Constructor body must declare the receiver before assigning fields.
// `this.x = x; this.y = y;` should NOT reference an undeclared `c` variable.
func TestBug2_ConstructorBodyDeclaredReceiver(t *testing.T) {
	src := `package com.google.javascript.jscomp;
final class CheckSideEffects {
  private final boolean report;
  private final boolean protect;
  CheckSideEffects(boolean report, boolean protect) {
    this.report = report;
    this.protect = protect;
  }
}`
	out := generateSkeleton(src, "com/google/javascript/jscomp/CheckSideEffects.java", "jscomp", "")
	// The constructor must declare the receiver (e.g. `c := &CheckSideEffects{}`).
	if !strings.Contains(out, ":= &CheckSideEffects{") {
		t.Errorf("Bug #2: constructor does not declare receiver with :=:\n%s", out)
	}
	// The constructor must return the receiver.
	if !strings.Contains(out, "return c") {
		t.Errorf("Bug #2: constructor does not return receiver variable:\n%s", out)
	}
}

// Bug #3a — Arithmetic shift `<<` in a field initialiser must not confuse the
// depth counter in findDeclTerminator, causing it to miss the terminating `;`
// and thus fail to reach the method that follows.
func TestBug3a_ShiftOperatorInFieldInitialiser(t *testing.T) {
	src := `package com.google.javascript.jscomp;
public final class NodeUtil {
  static final long MAX_POSITIVE_INTEGER_NUMBER = (1L << 53) - 1;
  static final String JSC_PROPERTY_NAME_FN = "JSCompiler_renameProperty";
  static boolean isExpression(Object n) { return true; }
  public static boolean isSomething(Object n) { return false; }
}`
	out := generateSkeleton(src, "com/google/javascript/jscomp/NodeUtil.java", "jscomp", "")
	// The static methods AFTER the tricky field must be emitted.
	// Before the fix, the depth counter would desync on `<<` and the parser
	// would never find the closing `;`, stopping all further processing.
	if !strings.Contains(out, "IsExpression") {
		t.Errorf("Bug #3a: static method after `<<` field missing from output:\n%s", out)
	}
	if !strings.Contains(out, "IsSomething") {
		t.Errorf("Bug #3a: second static method missing from output:\n%s", out)
	}
}

// Bug #3b — A Java file whose top-level declaration is a class containing an
// inner enum must be parsed as a class, not as that inner enum.
func TestBug3b_InnerEnumNotTopLevelClass(t *testing.T) {
	src := `package com.google.javascript.jscomp;
public final class NodeUtil {
  private NodeUtil() {}
  public enum ValueType {
    UNDETERMINED,
    NULL,
    VOID
  }
  public static boolean isExpression(Object n) { return true; }
}`
	out := generateSkeleton(src, "com/google/javascript/jscomp/NodeUtil.java", "jscomp", "")
	// The output must contain a NodeUtil struct, NOT a NodeUtil enum.
	if !strings.Contains(out, "type NodeUtil struct") {
		t.Errorf("Bug #3b: top-level should be class NodeUtil struct, got:\n%s", out)
	}
	// The inner enum must appear as an inner declaration, not as the file's type.
	if strings.Contains(out, "type NodeUtil int") {
		t.Errorf("Bug #3b: top-level was wrongly treated as enum:\n%s", out)
	}
	// Static method must be emitted as a package-level function.
	if !strings.Contains(out, "func IsExpression(") {
		t.Errorf("Bug #3b: static method not emitted as package-level func:\n%s", out)
	}
}

// Static methods must be emitted as package-level functions, not skipped.
func TestStaticMethodsEmittedAsPackageLevelFunctions(t *testing.T) {
	src := `package com.google.javascript.jscomp;
public final class NodeUtil {
  private NodeUtil() {}
  public static boolean isExpression(Object n) { return true; }
  private static String getNameHelper(Object n) { return ""; }
}`
	out := generateSkeleton(src, "com/google/javascript/jscomp/NodeUtil.java", "jscomp", "")
	if !strings.Contains(out, "func IsExpression(") {
		t.Errorf("public static method should be package-level func, got:\n%s", out)
	}
	if !strings.Contains(out, "func getNameHelper(") {
		t.Errorf("private static method should be package-level func, got:\n%s", out)
	}
}
