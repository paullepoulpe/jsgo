// Generates internal/parser/testdata/enums.txt:
//   java -cp closure-compiler.jar:build ParserEnumDump > internal/parser/testdata/enums.txt
import com.google.javascript.jscomp.parsing.parser.*;
import com.google.javascript.jscomp.parsing.parser.trees.*;
import java.lang.reflect.*;

// Oracle tool: dumps every enum of the parser packages (ordinal, name, toString) and every static FeatureSet
// constant (toString and version()) for internal/parser/testdata/enums.txt.
public class ParserEnumDump {
  public static void main(String[] a) throws Exception {
    String[] enums = {
      "parser.TokenType", "parser.Keywords", "parser.FeatureSet$Feature", "parser.FeatureSet$LangVersion",
      "parser.TemplateLiteralToken$ErrorLevel", "parser.Parser$FunctionFlavor", "parser.Parser$Config$Mode",
      "parser.Parser$Expression", "parser.Parser$PatternKind", "parser.trees.ParseTreeType",
      "parser.trees.Comment$Type", "parser.trees.ComprehensionTree$ComprehensionType",
      "parser.trees.FunctionDeclarationTree$Kind", "parser.trees.UpdateExpressionTree$OperatorPosition"};
    for (String e : enums) {
      Class<?> c = Class.forName("com.google.javascript.jscomp.parsing." + e);
      for (Object o : c.getEnumConstants()) {
        String s;
        try { s = o.toString(); } catch (RuntimeException ex) { s = "<" + ex.getClass().getSimpleName() + ">"; }
        System.out.println("enum " + c.getName().substring(c.getName().lastIndexOf('.') + 1) + " " + ((Enum<?>) o).ordinal() + " " + ((Enum<?>) o).name() + " " + s);
      }
    }
    for (Field f : FeatureSet.class.getFields()) {
      if (Modifier.isStatic(f.getModifiers()) && f.getType() == FeatureSet.class) {
        FeatureSet fs = (FeatureSet) f.get(null);
        System.out.println("featureset " + f.getName() + " " + fs.version() + " " + fs);
      }
    }
    for (Keywords k : Keywords.values()) System.out.println("keyword " + k.name() + " " + k.type.name() + " strict=" + Keywords.isStrictKeyword(k.type));
  }
}
