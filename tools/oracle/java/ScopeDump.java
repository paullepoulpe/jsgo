// Oracle tool: sweeps the API of NodeTraversal, SyntacticScopeCreator / Scope / Var,
// MemoizedScopeCreator and ReferenceCollector / ReferenceCollection / Reference over real files,
// and prints what they report. The Go side is internal/jscomp/scope_dump_test.go.
//
// Build/run (see tools/oracle/README.md):
//   javac -cp closure-compiler.jar -d build ScopeDump.java
//   java -Xss1g -cp closure-compiler.jar:build com.google.javascript.jscomp.ScopeDump
//        <out.gz> <root> <rel-file>... | @<listfile>
//
// Writes one gzipped dump: for each <rel-file>, a "# file <rel-file>" line and the sections
// below. Each file is parsed with
// ParserRunner.parse and the parsing.Config the compiler uses for --language_in=ECMASCRIPT_NEXT
// (as ParsingDump does), the SCRIPT gets the InputId <rel-file> and is put under a ROOT, and a
// fresh Compiler (new Compiler() + initOptions) gets a CompilerInput for that id (Compiler's
// private putCompilerInput, called reflectively) so that NodeTraversal.getInput() finds it.
//
// Nodes print as TOKEN@line:charno (or "null"). Exceptions print as "!" + the class names of the
// exception and its causes, joined by "<" (messages are JVM-specific for NPEs, so left out).
//
// Sections:
//   # traversal     NodeTraversal.traverse(ROOT) with a ScopedCallback recording, in order,
//                   S <n> p=<parent> d=<scopeDepth> r=<scopeRoot> h=<closestHoistScopeRoot>
//                     f=<enclosingFunction> b=<closestScopeRootNodeBindingThisOrSuper> m=<flags>
//                   (flags: G inGlobalHoistScope, M inModuleHoistScope, H isHoistScope,
//                    K inFunctionBlockScope, O inModuleScope, g inGlobalScope),
//                   V <n> p=<parent>, E <scopeRoot> d=<depth>, X <scopeRoot> d=<depth>
//   # obey          the same traversal with setObeyDestructuringAndDefaultValueExecutionOrder(true),
//                   recording only "S <n>" and "V <n>"
//   # scoperoots    for every FUNCTION (pre-order): "R <fn>", then
//                   NodeTraversal.traverseScopeRoots(compiler, [fn], recorder, false) recording
//                   "S <n>", "V <n>", "E <root>", "X <root>"
//   # atscope       for every FUNCTION scope (pre-order) of at most 64 nodes: "A <fn>", then
//                   NodeTraversal.builder()...setScopeCreator(memoized).traverseAtScope(scope),
//                   recording as in # scoperoots
//   # scopes        NodeTraversal.traverse(ROOT) with scope creator
//                   new SyntacticScopeCreator(compiler, <recording handler>, true) and a callback
//                   that calls t.getScope() on enterScope and prints the scope:
//                     scope <root> d=<depth> p=<parentRoot> f=<flags> h=<closestHoistScope>
//                       c=<closestContainerScope> t=<scopeOfThis> cfg=<closestCfgRootScope> n=<count>
//                   (flags: B block, S static block, K function block, F function, M module,
//                    D member field def, C computed field def rhs, T catch, R cfg root, H hoist,
//                    G global, L local), its implicit slots
//                     implicit arguments=<..> this=<..> super=<..> exports=<..>
//                   (Y/N hasOwnSlot, then V/- getOwnSlot non-null), getArgumentsVar
//                     args=<scope root of the var or null>
//                   and each var of getVarIterable():
//                     var <name> n=<nameNode> pn=<parentNode> i=<index> dt=<declarationType>
//                       f=<flags> init=<initialValue> doc=<jsdoc 0/1> sym=<getSymbol()==var>
//                       decl=<getDeclaration()==var> str=<implicit goog namespace strength>
//                   (flags: v var, l let, c const, k class, p param, h catch, i import,
//                    b bleeding function, a arguments, t this, e goog module exports,
//                    n implicit goog namespace, C declared or inferred const, D define, x extern,
//                    g global, d default param; each failing one is "!")
//                   Redeclarations reported to the handler print as
//                     redeclare <scopeRoot> <name> <n> input=<input name or null>
//   # references    new ReferenceCollector(compiler, <behavior>, new SyntacticScopeCreator(compiler))
//                   .process(ROOT); the behavior prints on each afterExitScope
//                     exit <scopeRoot> <var name>=<number of references or -> ...
//                   for the vars of t.getScope(); then for each var of getAllSymbols():
//                     sym <name> s=<scopeRoot> n=<nameNode> g=<getScope(var) root>
//                     coll wd=.. esc=.. once=.. never=.. first=.. weak=.. init=<index>
//                       initc=<index> one=<index>
//                     ref <n> p=<parent> s=<scope root> bb=<block root> bbp=<parent block root>
//                       bbf=<isFunction/isLoop> f=<flags> av=<assignedValue> sym=<symbol name node>
//                   (ref flags: D declaration, L lvalue, I initializing declaration, v var decl,
//                    l let decl, c const decl, h hoisted function, s simple assignment to name,
//                    o dot property access, O assigned to object destructuring pattern)
//   # memoized      MemoizedScopeCreator over a SyntacticScopeCreator: createScope(ROOT, null)
//                   twice (same instance) and createScope(ROOT, global) (the precondition)
//   # exception ... if a section threw
// Every line is escaped like ParsingDump: '\' -> "\\", every char outside printable ASCII ->
// \\uXXXX.
package com.google.javascript.jscomp;

import com.google.javascript.jscomp.parsing.Config;
import com.google.javascript.jscomp.parsing.ParserRunner;
import com.google.javascript.rhino.ErrorReporter;
import com.google.javascript.rhino.IR;
import com.google.javascript.rhino.InputId;
import com.google.javascript.rhino.Node;
import java.io.OutputStream;
import java.io.PrintStream;
import java.lang.reflect.Method;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.function.Supplier;
import java.util.zip.GZIPOutputStream;

public final class ScopeDump {
  public static void main(String[] args) throws Exception {
    Path out = Paths.get(args[0]);
    Path root = Paths.get(args[1]);
    List<String> files = new ArrayList<>();
    for (int i = 2; i < args.length; i++) {
      String f = args[i];
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
    try (OutputStream os = new GZIPOutputStream(Files.newOutputStream(out))) {
      for (String rel : files) {
        StringBuilder sb = new StringBuilder();
        dumpOne(sb, rel, Files.readString(root.resolve(rel), StandardCharsets.UTF_8));
        os.write(sb.toString().getBytes(StandardCharsets.UTF_8));
      }
    }
  }

  static Compiler newCompiler() {
    Compiler compiler = new Compiler(new PrintStream(OutputStream.nullOutputStream()));
    CompilerOptions options = new CompilerOptions();
    options.setLanguageIn(CompilerOptions.LanguageMode.ECMASCRIPT_NEXT);
    compiler.initOptions(options);
    return compiler;
  }

  static void putCompilerInput(Compiler compiler, CompilerInput input) throws Exception {
    Method m = Compiler.class.getDeclaredMethod("putCompilerInput", CompilerInput.class);
    m.setAccessible(true);
    m.invoke(compiler, input);
  }

  static void dumpOne(StringBuilder sb, String name, String src) throws Exception {
    Compiler compiler = newCompiler();
    Config config = compiler.getParserConfig(AbstractCompiler.ConfigContext.DEFAULT);
    ErrorReporter reporter =
        new ErrorReporter() {
          @Override
          public void warning(String message, String sourceName, int line, int lineOffset) {}

          @Override
          public void error(String message, String sourceName, int line, int lineOffset) {}
        };
    SourceFile file = SourceFile.fromCode(name, src);
    line(sb, "# file " + name);
    ParserRunner.ParseResult result;
    try {
      result = ParserRunner.parse(file, src, config, reporter);
    } catch (RuntimeException e) {
      line(sb, "# exception " + exc(e));
      return;
    }
    if (result.ast == null) {
      line(sb, "# no ast");
      return;
    }
    Node script = result.ast;
    script.setInputId(new InputId(name));
    Node root = IR.root(script);
    putCompilerInput(compiler, new CompilerInput(file));

    section(sb, "# traversal", () -> traversal(sb, compiler, root, false, true));
    section(sb, "# obey", () -> traversal(sb, compiler, root, true, false));
    section(sb, "# scoperoots", () -> scopeRoots(sb, compiler, root));
    section(sb, "# atscope", () -> atScope(sb, compiler, root));
    section(sb, "# scopes", () -> scopes(sb, compiler, root));
    section(sb, "# references", () -> references(sb, compiler, root));
    section(sb, "# memoized", () -> memoized(sb, compiler, root));
  }

  static void section(StringBuilder sb, String title, Runnable r) {
    line(sb, title);
    try {
      r.run();
    } catch (RuntimeException | Error e) {
      line(sb, "# exception " + exc(e));
    }
  }

  // ---------------------------------------------------------------------------------------------
  // NodeTraversal

  /** Records the callbacks of a traversal. */
  static final class Recorder implements NodeTraversal.ScopedCallback {
    final StringBuilder sb;
    final boolean verbose;
    final boolean scoped;

    Recorder(StringBuilder sb, boolean verbose, boolean scoped) {
      this.sb = sb;
      this.verbose = verbose;
      this.scoped = scoped;
    }

    @Override
    public boolean shouldTraverse(NodeTraversal t, Node n, Node parent) {
      if (verbose) {
        line(
            sb,
            "S "
                + node(n)
                + " p="
                + node(parent)
                + " d="
                + str(() -> String.valueOf(t.getScopeDepth()))
                + " r="
                + node(t.getScopeRoot())
                + " h="
                + node(t.getClosestHoistScopeRoot())
                + " f="
                + node(t.getEnclosingFunction())
                + " b="
                + node(t.getClosestScopeRootNodeBindingThisOrSuper())
                + " m="
                + str(() -> traversalFlags(t)));
      } else {
        line(sb, "S " + node(n));
      }
      return true;
    }

    @Override
    public void visit(NodeTraversal t, Node n, Node parent) {
      if (verbose) {
        line(sb, "V " + node(n) + " p=" + node(parent));
      } else {
        line(sb, "V " + node(n));
      }
    }

    @Override
    public void enterScope(NodeTraversal t) {
      if (scoped) {
        line(sb, "E " + node(t.getScopeRoot()) + " d=" + str(() -> String.valueOf(t.getScopeDepth())));
      }
    }

    @Override
    public void exitScope(NodeTraversal t) {
      if (scoped) {
        line(sb, "X " + node(t.getScopeRoot()) + " d=" + str(() -> String.valueOf(t.getScopeDepth())));
      }
    }
  }

  static String traversalFlags(NodeTraversal t) {
    StringBuilder f = new StringBuilder();
    if (t.inGlobalHoistScope()) {
      f.append('G');
    }
    if (t.inModuleHoistScope()) {
      f.append('M');
    }
    if (t.isHoistScope()) {
      f.append('H');
    }
    if (t.inFunctionBlockScope()) {
      f.append('K');
    }
    if (t.inModuleScope()) {
      f.append('O');
    }
    if (t.inGlobalScope()) {
      f.append('g');
    }
    return f.toString();
  }

  static void traversal(StringBuilder sb, Compiler compiler, Node root, boolean obey, boolean verbose) {
    NodeTraversal.builder()
        .setCompiler(compiler)
        .setCallback(new Recorder(sb, verbose, verbose))
        .setObeyDestructuringAndDefaultValueExecutionOrder(obey)
        .traverse(root);
  }

  static List<Node> functions(Node root) {
    List<Node> fns = new ArrayList<>();
    collectFunctions(root, fns);
    return fns;
  }

  static void collectFunctions(Node n, List<Node> out) {
    if (n.isFunction()) {
      out.add(n);
    }
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      collectFunctions(c, out);
    }
  }

  static int size(Node n) {
    int s = 1;
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      s += size(c);
    }
    return s;
  }

  static void scopeRoots(StringBuilder sb, Compiler compiler, Node root) {
    for (Node fn : functions(root)) {
      line(sb, "R " + node(fn));
      try {
        NodeTraversal.traverseScopeRoots(
            compiler, List.of(fn), new Recorder(sb, false, true), false);
      } catch (RuntimeException | Error e) {
        line(sb, "# exception " + exc(e));
      }
    }
  }

  static void atScope(StringBuilder sb, Compiler compiler, Node root) {
    MemoizedScopeCreator creator = new MemoizedScopeCreator(new SyntacticScopeCreator(compiler));
    List<Scope> fnScopes = new ArrayList<>();
    NodeTraversal.builder()
        .setCompiler(compiler)
        .setScopeCreator(creator)
        .setCallback(
            new NodeTraversal.AbstractScopedCallback() {
              @Override
              public void enterScope(NodeTraversal t) {
                if (t.getScopeRoot().isFunction()) {
                  fnScopes.add(t.getScope());
                }
              }

              @Override
              public void visit(NodeTraversal t, Node n, Node parent) {}
            })
        .traverse(root);
    for (Scope s : fnScopes) {
      if (size(s.getRootNode()) > 64) {
        continue;
      }
      line(sb, "A " + node(s.getRootNode()));
      try {
        NodeTraversal.builder()
            .setCompiler(compiler)
            .setCallback(new Recorder(sb, false, true))
            .setScopeCreator(creator)
            .traverseAtScope(s);
      } catch (RuntimeException | Error e) {
        line(sb, "# exception " + exc(e));
      }
    }
  }

  // ---------------------------------------------------------------------------------------------
  // Scopes

  static void scopes(StringBuilder sb, Compiler compiler, Node root) {
    SyntacticScopeCreator.RedeclarationHandler handler =
        (s, name, n, input) ->
            line(
                sb,
                "redeclare "
                    + node(s.getRootNode())
                    + " "
                    + name
                    + " "
                    + node(n)
                    + " input="
                    + (input == null ? "null" : input.getName()));
    SyntacticScopeCreator creator = new SyntacticScopeCreator(compiler, handler, true);
    NodeTraversal.builder()
        .setCompiler(compiler)
        .setScopeCreator(creator)
        .setCallback(
            new NodeTraversal.AbstractScopedCallback() {
              @Override
              public void enterScope(NodeTraversal t) {
                printScope(sb, t.getScope());
              }

              @Override
              public void visit(NodeTraversal t, Node n, Node parent) {}
            })
        .traverse(root);
  }

  static String scopeRoot(AbstractScope<?, ?> s) {
    return s == null ? "null" : node(s.getRootNode());
  }

  static void printScope(StringBuilder sb, Scope s) {
    line(
        sb,
        "scope "
            + node(s.getRootNode())
            + " d="
            + s.getDepth()
            + " p="
            + scopeRoot(s.getParent())
            + " f="
            + scopeFlags(s)
            + " h="
            + str(() -> scopeRoot(s.getClosestHoistScope()))
            + " c="
            + str(() -> scopeRoot(s.getClosestContainerScope()))
            + " t="
            + str(() -> scopeRoot(s.getScopeOfThis()))
            + " cfg="
            + str(() -> scopeRoot(s.getClosestCfgRootScope()))
            + " n="
            + s.getVarCount());
    StringBuilder implicit = new StringBuilder("implicit");
    for (String name : new String[] {"arguments", "this", "super", "exports"}) {
      implicit
          .append(' ')
          .append(name)
          .append('=')
          .append(str(() -> s.hasOwnSlot(name) ? "Y" : "N"))
          .append(str(() -> s.getOwnSlot(name) != null ? "V" : "-"));
    }
    line(sb, implicit.toString());
    line(
        sb,
        "args="
            + str(
                () -> {
                  Var v = s.getArgumentsVar();
                  return v == null ? "null" : node(v.getScopeRoot()) + "/" + v.getIndex();
                }));
    for (Var v : s.getVarIterable()) {
      printVar(sb, v);
    }
  }

  static String scopeFlags(Scope s) {
    StringBuilder f = new StringBuilder();
    flag(f, 'B', s::isBlockScope);
    flag(f, 'S', s::isStaticBlockScope);
    flag(f, 'K', s::isFunctionBlockScope);
    flag(f, 'F', s::isFunctionScope);
    flag(f, 'M', s::isModuleScope);
    flag(f, 'D', s::isMemberFieldDefScope);
    flag(f, 'C', s::isComputedFieldDefRhsScope);
    flag(f, 'T', s::isCatchScope);
    flag(f, 'R', s::isCfgRootScope);
    flag(f, 'H', s::isHoistScope);
    flag(f, 'G', s::isGlobal);
    flag(f, 'L', s::isLocal);
    return f.toString();
  }

  static void printVar(StringBuilder sb, Var v) {
    StringBuilder f = new StringBuilder();
    flag(f, 'v', v::isVar);
    flag(f, 'l', v::isLet);
    flag(f, 'c', v::isConst);
    flag(f, 'k', v::isClass);
    flag(f, 'p', v::isParam);
    flag(f, 'h', v::isCatch);
    flag(f, 'i', v::isImport);
    flag(f, 'b', v::isBleedingFunction);
    flag(f, 'a', v::isArguments);
    flag(f, 't', v::isThis);
    flag(f, 'e', v::isGoogModuleExports);
    flag(f, 'n', v::isImplicitGoogNamespace);
    flag(f, 'C', v::isDeclaredOrInferredConst);
    flag(f, 'D', v::isDefine);
    flag(f, 'x', v::isExtern);
    flag(f, 'g', v::isGlobal);
    flag(f, 'd', v::isDefaultParam);
    line(
        sb,
        "var "
            + v.getName()
            + " n="
            + node(v.getNameNode())
            + " pn="
            + node(v.getParentNode())
            + " i="
            + v.getIndex()
            + " dt="
            + str(() -> String.valueOf(v.declarationType()))
            + " f="
            + f
            + " init="
            + str(() -> node(v.getInitialValue()))
            + " doc="
            + (v.getJSDocInfo() != null ? 1 : 0)
            + " sym="
            + (v.getSymbol() == v ? 1 : 0)
            + " decl="
            + (v.getDeclaration() == v ? 1 : 0)
            + " str="
            + str(() -> String.valueOf(v.getImplicitGoogNamespaceStrength())));
  }

  interface BoolSupplier {
    boolean get();
  }

  static void flag(StringBuilder f, char c, BoolSupplier b) {
    try {
      if (b.get()) {
        f.append(c);
      }
    } catch (RuntimeException | Error e) {
      f.append('!');
    }
  }

  // ---------------------------------------------------------------------------------------------
  // References

  static void references(StringBuilder sb, Compiler compiler, Node root) {
    ReferenceCollector.Behavior behavior =
        (t, referenceMap) -> {
          StringBuilder l = new StringBuilder("exit ").append(node(t.getScopeRoot()));
          for (Var v : t.getScope().getVarIterable()) {
            ReferenceCollection refs = referenceMap.getReferences(v);
            l.append(' ')
                .append(v.getName())
                .append('=')
                .append(refs == null ? "-" : String.valueOf(refs.references.size()));
          }
          line(sb, l.toString());
        };
    ReferenceCollector collector =
        new ReferenceCollector(compiler, behavior, new SyntacticScopeCreator(compiler));
    collector.process(root);
    for (Var v : collector.getAllSymbols()) {
      ReferenceCollection refs = collector.getReferences(v);
      line(
          sb,
          "sym "
              + v.getName()
              + " s="
              + node(v.getScopeRoot())
              + " n="
              + node(v.getNameNode())
              + " g="
              + node(collector.getScope(v).getRootNode()));
      line(
          sb,
          "coll wd="
              + str(() -> String.valueOf(refs.isWellDefined()))
              + " esc="
              + str(() -> String.valueOf(refs.isEscaped()))
              + " once="
              + str(() -> String.valueOf(refs.isAssignedOnceInLifetime()))
              + " never="
              + str(() -> String.valueOf(refs.isNeverAssigned()))
              + " first="
              + str(() -> String.valueOf(refs.firstReferenceIsAssigningDeclaration()))
              + " weak="
              + str(() -> String.valueOf(refs.isReferencedWeakly()))
              + " init="
              + str(() -> String.valueOf(refs.references.indexOf(refs.getInitializingReference())))
              + " initc="
              + str(
                  () ->
                      String.valueOf(
                          refs.references.indexOf(refs.getInitializingReferenceForConstants())))
              + " one="
              + str(() -> String.valueOf(refs.references.indexOf(refs.getOneAndOnlyAssignment()))));
      for (Reference ref : refs) {
        printReference(sb, ref);
      }
    }
  }

  static void printReference(StringBuilder sb, Reference ref) {
    StringBuilder f = new StringBuilder();
    flag(f, 'D', ref::isDeclaration);
    flag(f, 'L', ref::isLvalue);
    flag(f, 'I', ref::isInitializingDeclaration);
    flag(f, 'v', ref::isVarDeclaration);
    flag(f, 'l', ref::isLetDeclaration);
    flag(f, 'c', ref::isConstDeclaration);
    flag(f, 'h', ref::isHoistedFunction);
    flag(f, 's', ref::isSimpleAssignmentToName);
    flag(f, 'o', ref::isDotPropertyAccess);
    flag(f, 'O', ref::isAssignedToObjectDestructuringPattern);
    BasicBlock bb = ref.getBasicBlock();
    line(
        sb,
        "ref "
            + node(ref.getNode())
            + " p="
            + node(ref.getParent())
            + " s="
            + node(ref.getScope().getRootNode())
            + " bb="
            + node(bb.getRoot())
            + " bbp="
            + (bb.getParent() == null ? "null" : node(bb.getParent().getRoot()))
            + " bbf="
            + (bb.isFunction() ? "F" : "")
            + (bb.isLoop() ? "L" : "")
            + " f="
            + f
            + " av="
            + str(() -> node(ref.getAssignedValue()))
            + " sym="
            + str(
                () -> {
                  Var s = ref.getSymbol();
                  return s == null ? "null" : node(s.getNameNode()) + "/" + node(s.getScopeRoot());
                }));
  }

  // ---------------------------------------------------------------------------------------------
  // MemoizedScopeCreator

  static void memoized(StringBuilder sb, Compiler compiler, Node root) {
    ScopeCreator creator = new MemoizedScopeCreator(new SyntacticScopeCreator(compiler));
    AbstractScope<?, ?> a = creator.createScope(root, null);
    AbstractScope<?, ?> b = creator.createScope(root, null);
    line(sb, "same=" + (a == b) + " vars=" + ((Scope) a).getVarCount());
    line(
        sb,
        "reparent="
            + str(
                () -> {
                  creator.createScope(root, a);
                  return "ok";
                }));
  }

  // ---------------------------------------------------------------------------------------------
  // Formatting

  static String node(Node n) {
    if (n == null) {
      return "null";
    }
    return n.getToken() + "@" + n.getLineno() + ":" + n.getCharno();
  }

  static String str(Supplier<String> s) {
    try {
      return s.get();
    } catch (RuntimeException | Error e) {
      return "!" + exc(e);
    }
  }

  static String exc(Throwable e) {
    StringBuilder b = new StringBuilder(e.getClass().getName());
    for (Throwable c = e.getCause(); c != null && c != e; c = c.getCause()) {
      b.append('<').append(c.getClass().getName());
      if (c.getCause() == c) {
        break;
      }
    }
    return b.toString();
  }

  /** Backslash and every char outside printable ASCII become escapes. */
  static String escape(String s) {
    StringBuilder b = new StringBuilder(s.length());
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '\\') {
        b.append("\\\\");
      } else if (c >= 0x20 && c < 0x7f) {
        b.append(c);
      } else {
        b.append(String.format("\\u%04x", (int) c));
      }
    }
    return b.toString();
  }

  static void line(StringBuilder sb, String s) {
    sb.append(escape(s)).append('\n');
  }
}
