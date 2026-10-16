// Oracle tool: a fixed scenario on TypedScope / TypedVar (which have no unit tests of their own),
// printing the result of every query. The Go side is internal/jscomp/typed_scope_test.go, which
// runs the same scenario and compares with internal/jscomp/testdata/typed_scope.txt.
//
//   javac -cp closure-compiler.jar -d build TypedScopeVectors.java
//   java -cp closure-compiler.jar:build com.google.javascript.jscomp.TypedScopeVectors <out.txt>
//
// The code below is parsed like ScopeDump does (ParserRunner.parse, ECMASCRIPT_NEXT compiler
// config); TypedScopes are built by hand over its nodes with types from the Compiler's
// JSTypeRegistry, then every scope is queried for a list of names. Exceptions print as "!" + the
// exception class name.
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.parsing.Config;
import com.google.javascript.jscomp.parsing.ParserRunner;
import com.google.javascript.rhino.ErrorReporter;
import com.google.javascript.rhino.IR;
import com.google.javascript.rhino.InputId;
import com.google.javascript.rhino.Node;
import com.google.javascript.rhino.QualifiedName;
import com.google.javascript.rhino.StaticScope;
import com.google.javascript.rhino.Token;
import com.google.javascript.rhino.jstype.JSType;
import com.google.javascript.rhino.jstype.JSTypeNative;
import com.google.javascript.rhino.jstype.JSTypeRegistry;
import java.io.OutputStream;
import java.io.PrintStream;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.function.Supplier;

public final class TypedScopeVectors {
  static final String SRC =
      String.join(
          "\n",
          "var a = 1;",
          "a.b = 2;",
          "function f(x) { var a; var u; let w; this.q; }",
          "class C { m() {} static { } fld = 1; static sfld = 2; }",
          "var u;",
          "var arguments;",
          "");

  static final String[] NAMES = {
    "a", "a.b", "a.c", "b", "x", "u", "z", "w", "this", "super", "arguments", "exports", "f", "C"
  };

  static final Map<Object, String> labels = new LinkedHashMap<>();

  static String label(Object scope) {
    if (scope == null) {
      return "null";
    }
    String l = labels.get(scope);
    return l == null ? "?" : l;
  }

  static String str(Supplier<String> s) {
    try {
      return s.get();
    } catch (RuntimeException | Error e) {
      return "!" + e.getClass().getName();
    }
  }

  static String describe(TypedVar v) {
    if (v == null) {
      return "null";
    }
    return v
        + "|scope="
        + label(v.getScope())
        + "|type="
        + v.getType()
        + "|inferred="
        + v.isTypeInferred()
        + "|index="
        + v.getIndex()
        + "|node="
        + (v.getNode() == null ? "null" : v.getNode().getToken())
        + "|input="
        + v.getInputName()
        + "|extern="
        + v.isExtern();
  }

  public static void main(String[] args) throws Exception {
    PrintStream out =
        new PrintStream(Files.newOutputStream(Paths.get(args[0])), true, "UTF-8");
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(CompilerOptions.LanguageMode.ECMASCRIPT_NEXT);
    compiler.initOptions(options);
    Config config = compiler.getParserConfig(AbstractCompiler.ConfigContext.DEFAULT);
    ErrorReporter reporter =
        new ErrorReporter() {
          @Override
          public void warning(String message, String sourceName, int line, int lineOffset) {}

          @Override
          public void error(String message, String sourceName, int line, int lineOffset) {}
        };
    SourceFile file = SourceFile.fromCode("typed.js", SRC);
    Node script = ParserRunner.parse(file, SRC, config, reporter).ast;
    script.setInputId(new InputId("typed.js"));
    Node root = IR.root(script);
    CompilerInput input = new CompilerInput(file);
    JSTypeRegistry registry = compiler.getTypeRegistry();

    Node aName = script.getFirstChild().getFirstChild();
    Node abGetprop = script.getSecondChild().getFirstChild().getFirstChild();
    Node fn = script.getChildAtIndex(2);
    Node xName = fn.getSecondChild().getFirstChild();
    Node fBody = fn.getLastChild();
    Node innerA = fBody.getFirstChild().getFirstChild();
    Node innerU = fBody.getSecondChild().getFirstChild();
    Node cls = script.getChildAtIndex(3);
    Node members = cls.getLastChild();
    Node mFn = members.getFirstChild().getFirstChild();
    Node staticBlock = members.getSecondChild();
    Node fld = members.getChildAtIndex(2);
    Node sfld = members.getChildAtIndex(3);
    Node uName = script.getChildAtIndex(4).getFirstChild();
    Node argsName = script.getChildAtIndex(5).getFirstChild();

    JSType num = registry.getNativeType(JSTypeNative.NUMBER_TYPE);
    JSType string = registry.getNativeType(JSTypeNative.STRING_TYPE);
    JSType bool = registry.getNativeType(JSTypeNative.BOOLEAN_TYPE);
    JSType obj = registry.getNativeType(JSTypeNative.OBJECT_TYPE);
    JSType ctor = registry.getNativeType(JSTypeNative.DATE_FUNCTION_TYPE);
    root.setJSType(obj);
    fn.setJSType(ctor);
    cls.setJSType(ctor);

    TypedScope global = TypedScope.createGlobalScope(root);
    labels.put(global, "global");
    global.declare("a", aName, num, input, false);
    global.declare("a.b", abGetprop, string, input, true);
    global.declare("u", uName, null, input, false);
    global.declare("arguments", argsName, obj, null, false);
    global.declare("f", fn.getFirstChild(), ctor, input, false);
    global.declare("C", cls.getFirstChild(), ctor, input, false);
    out.println("declare-again " + str(() -> describe(global.declare("a", aName, num, input, false))));

    TypedScope fScope = new TypedScope(global, fn);
    labels.put(fScope, "fScope");
    fScope.declare("x", xName, null, input, true);
    TypedScope fBlock =
        new TypedScope(fScope, fBody, new LinkedHashSet<>(List.of("a", "z", "u")), null);
    labels.put(fBlock, "fBlock");
    fBlock.declare("a", innerA, bool, input, false);
    fBlock.declare("u", innerU, null, input, false);
    TypedScope clsScope = new TypedScope(global, cls);
    labels.put(clsScope, "clsScope");
    TypedScope mScope = new TypedScope(clsScope, mFn);
    labels.put(mScope, "mScope");
    TypedScope sbScope = new TypedScope(clsScope, staticBlock);
    labels.put(sbScope, "sbScope");
    TypedScope fldScope = new TypedScope(clsScope, fld);
    labels.put(fldScope, "fldScope");
    TypedScope sfldScope = new TypedScope(clsScope, sfld);
    labels.put(sfldScope, "sfldScope");
    TypedScope bottom = TypedScope.createLatticeBottom(new Node(Token.ROOT));
    labels.put(bottom, "bottom");
    out.println(
        "bad-child " + str(() -> label(new TypedScope(global, script.getChildAtIndex(4)))));

    List<TypedScope> scopes =
        List.of(global, fScope, fBlock, clsScope, mScope, sbScope, fldScope, sfldScope, bottom);
    for (TypedScope s : scopes) {
      out.println(
          "scope "
              + label(s)
              + " depth="
              + s.getDepth()
              + " parent="
              + label(s.getParent())
              + " bottom="
              + s.isBottom()
              + " module="
              + s.getModule()
              + " this="
              + str(() -> String.valueOf(s.getTypeOfThis()))
              + " count="
              + s.getVarCount()
              + " global="
              + label(s.getGlobalScope()));
      List<String> names = new ArrayList<>();
      for (TypedVar v : s.getVarIterable()) {
        names.add(v.getName());
      }
      out.println("  vars " + names);
      List<String> unbound = new ArrayList<>();
      for (TypedVar v : s.getDeclarativelyUnboundVarsWithoutTypes()) {
        unbound.add(v.getName());
      }
      out.println("  unbound " + unbound);
      out.println("  args " + str(() -> describe(s.getArgumentsVar())));
      for (String name : NAMES) {
        out.println(
            "  "
                + name
                + " get="
                + str(() -> describe(s.getVar(name)))
                + " slot="
                + str(() -> describe((TypedVar) s.getSlot(name)))
                + " own="
                + str(() -> describe(s.getOwnSlot(name)))
                + " hasOwn="
                + str(() -> String.valueOf(s.hasOwnSlot(name)))
                + " has="
                + str(() -> String.valueOf(s.hasSlot(name)))
                + " topmost="
                + str(() -> label((StaticScope) s.getTopmostScopeOfEventualDeclaration(name))));
      }
      for (String ns : new String[] {"a", "a.b", "a.b.c", "f", "nope"}) {
        out.println(
            "  ns "
                + ns
                + " "
                + str(() -> String.valueOf(s.getTypeThroughNamespace(ns)))
                + " lookup="
                + str(() -> String.valueOf(s.lookupQualifiedName(QualifiedName.of(ns)))));
      }
    }

    out.println("validate fBlock " + str(() -> { fBlock.validateCompletelyBuilt(); return "ok"; }));
    out.println("validate fBlock " + str(() -> { fBlock.validateCompletelyBuilt(); return "ok"; }));
    out.println("validate global " + str(() -> { global.validateCompletelyBuilt(); return "ok"; }));

    TypedVar u = global.getVar("u");
    u.setType(num);
    u.resolveType(reporter);
    u.markEscaped();
    out.println("u " + describe(u) + " escaped=" + u.isMarkedEscaped() + " once=" + u.isMarkedAssignedExactlyOnce());
    u.markAssignedExactlyOnce();
    out.println("u once=" + u.isMarkedAssignedExactlyOnce() + " symbol=" + (u.getSymbol() == u) + " decl=" + (u.getDeclaration() == u));
    out.println("undeclare " + str(() -> { fBlock.undeclare(fBlock.getVar("u")); return "ok"; }));
    out.println("after-undeclare " + str(() -> describe(fBlock.getVar("u"))));
    out.println("undeclare-foreign " + str(() -> { fBlock.undeclare(global.getVar("a")); return "ok"; }));
    out.println(
        "bad-name-node "
            + str(() -> describe(new TypedVar(false, "q", fn.getSecondChild(), null, global, 0, input))));
    out.println("this-var " + str(() -> describe(fScope.getVar("this"))) + " type=" + str(() -> String.valueOf(fScope.getVar("this").getType())));
    out.println("super-var " + str(() -> describe(mScope.getVar("super"))));
    out.println("contains " + global.contains(fBlock) + " " + fBlock.contains(global) + " common=" + label(fBlock.getCommonParent(mScope)));
    out.println("hoist " + label(fBlock.getClosestHoistScope()) + " container=" + label(fBlock.getClosestContainerScope()) + " sameContainer=" + fBlock.hasSameContainerScope(fScope));
  }
}
