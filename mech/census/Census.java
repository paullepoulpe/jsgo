// Census: type-attributes Closure Compiler's production sources with javac's
// public com.sun.source API and reports every member of a class outside the
// source set that the sources use. This is the external surface that the Go
// side must supply (handwritten shim, or a direct lowering in the transpiler).
//
//   java Census.java <classpath> <out.tsv> <file.java>...
import com.sun.source.tree.*;
import com.sun.source.util.*;
import java.io.*;
import java.nio.file.*;
import java.util.*;
import javax.lang.model.element.*;
import javax.lang.model.type.*;
import javax.tools.*;

public class Census {
  public static void main(String[] a) throws Exception {
    JavaCompiler jc = ToolProvider.getSystemJavaCompiler();
    StandardJavaFileManager fm = jc.getStandardFileManager(null, null, null);
    List<File> files = new ArrayList<>();
    for (int i = 2; i < a.length; i++) {
      if (a[i].startsWith("@")) for (String l : Files.readAllLines(Path.of(a[i].substring(1)))) files.add(new File(l));
      else files.add(new File(a[i]));
    }
    DiagnosticCollector<JavaFileObject> diags = new DiagnosticCollector<>();
    JavacTask task = (JavacTask) jc.getTask(null, fm, diags,
        List.of("-proc:none", "-cp", a[0], "--release", "21", "-Xmaxerrs", "100000"),
        null, fm.getJavaFileObjectsFromFiles(files));
    Iterable<? extends CompilationUnitTree> units = task.parse();
    task.analyze();
    long errors = diags.getDiagnostics().stream().filter(d -> d.getKind() == Diagnostic.Kind.ERROR).count();
    diags.getDiagnostics().stream().filter(d -> d.getKind() == Diagnostic.Kind.ERROR).limit(15)
        .forEach(d -> System.err.println(d));
    Trees trees = Trees.instance(task);
    javax.lang.model.util.Types types = task.getTypes();
    // kind \t owner \t member  ->  use count
    TreeMap<String, Integer> uses = new TreeMap<>();
    int nUnits = 0;
    for (CompilationUnitTree cu : units) {
      nUnits++;
      new TreePathScanner<Void, Void>() {
        void rec(Element e) {
          if (e == null) return;
          Element owner = e;
          String kind;
          switch (e.getKind()) {
            case METHOD: kind = e.getModifiers().contains(Modifier.STATIC) ? "static-method" : "method"; owner = e.getEnclosingElement(); break;
            case CONSTRUCTOR: kind = "constructor"; owner = e.getEnclosingElement(); break;
            case FIELD: case ENUM_CONSTANT: kind = "field"; owner = e.getEnclosingElement(); break;
            case CLASS: case INTERFACE: case ENUM: case RECORD: case ANNOTATION_TYPE: kind = "type"; break;
            default: return;
          }
          if (!(owner instanceof TypeElement)) return;
          TypeElement te = (TypeElement) owner;
          if (trees.getTree(te) != null) return; // declared in the source set
          String on = te.getQualifiedName().toString();
          if (on.startsWith("com.google.javascript.") || on.startsWith("com.google.debugging.")) kind = "proto-" + kind; // generated protobuf classes
          String member = "";
          if (e instanceof ExecutableElement) {
            ExecutableElement x = (ExecutableElement) e;
            StringBuilder sb = new StringBuilder(x.getSimpleName()).append('(');
            String sep = "";
            for (VariableElement p : x.getParameters()) { sb.append(sep).append(types.erasure(p.asType())); sep = ","; }
            member = sb.append(')').toString();
          } else if (!kind.endsWith("type")) member = e.getSimpleName().toString();
          uses.merge(kind + "\t" + on + "\t" + member, 1, Integer::sum);
        }
        Element el() { return trees.getElement(getCurrentPath()); }
        @Override public Void visitMethodInvocation(MethodInvocationTree t, Void v) {
          rec(trees.getElement(new TreePath(getCurrentPath(), t.getMethodSelect()))); return super.visitMethodInvocation(t, v); }
        @Override public Void visitNewClass(NewClassTree t, Void v) { rec(el()); return super.visitNewClass(t, v); }
        @Override public Void visitMemberReference(MemberReferenceTree t, Void v) { rec(el()); return super.visitMemberReference(t, v); }
        @Override public Void visitMemberSelect(MemberSelectTree t, Void v) {
          Element e = el(); if (e != null && (e.getKind().isField() || e.getKind().isClass() || e.getKind().isInterface())) rec(e);
          return super.visitMemberSelect(t, v); }
        @Override public Void visitIdentifier(IdentifierTree t, Void v) {
          Element e = el(); if (e != null && (e.getKind().isField() || e.getKind().isClass() || e.getKind().isInterface())) rec(e);
          return super.visitIdentifier(t, v); }
        @Override public Void visitImport(ImportTree t, Void v) { return null; }
        @Override public Void visitAnnotation(AnnotationTree t, Void v) { return null; }
      }.scan(cu, null);
    }
    try (PrintWriter w = new PrintWriter(Files.newBufferedWriter(Path.of(a[1])))) {
      uses.forEach((k, n) -> w.println(k + "\t" + n));
    }
    System.err.println("units=" + nUnits + " errors=" + errors + " distinct=" + uses.size());
  }
}
