// Oracle tool: builds rhino Node / JSDocInfo / Color values through the
// public API and prints their Java renderings (Node.toStringTree(),
// Node.toString(true, true, false), JSDocInfo.toStringVerbose(),
// Node.appendJsonTree, Color.toString), one escaped record per line. The Go
// test internal/rhino/format_oracle_test.go builds the same values with the
// ported API and must reproduce internal/rhino/testdata/format_oracle.txt
// byte for byte.
//
//   java -cp .cache/$CLOSURE_JAR:.cache/oracle-build \
//       com.google.javascript.rhino.RhinoFormatOracle > internal/rhino/testdata/format_oracle.txt
//
// It lives in package com.google.javascript.rhino to reach package-private
// API (Node.Prop, Node.putProp) like the Java unit tests do.
package com.google.javascript.rhino;

import com.google.javascript.jscomp.colors.Color;
import com.google.javascript.jscomp.colors.ColorId;
import com.google.javascript.jscomp.colors.StandardColors;
import com.google.javascript.jscomp.jarjar.com.google.common.collect.ImmutableSet;
import com.google.javascript.jscomp.parsing.parser.FeatureSet;
import com.google.javascript.jscomp.parsing.parser.util.SourcePosition;
import com.google.javascript.rhino.dtoa.DToA;
import java.io.IOException;
import java.io.PrintStream;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.util.LinkedHashSet;

public final class RhinoFormatOracle {
  private static PrintStream out;

  /** Escapes a Java string to printable ASCII: \n, \\ and \\uXXXX for everything else. */
  static String esc(String s) {
    StringBuilder sb = new StringBuilder();
    for (int i = 0; i < s.length(); i++) {
      char c = s.charAt(i);
      if (c == '\\') {
        sb.append("\\\\");
      } else if (c == '\n') {
        sb.append("\\n");
      } else if (c < 0x20 || c >= 0x7f) {
        sb.append(String.format("\\u%04x", (int) c));
      } else {
        sb.append(c);
      }
    }
    return sb.toString();
  }

  static void rec(String label, String value) {
    out.println(label + "\t" + esc(value));
  }

  static void dumpNodes(Node n, int depth) {
    rec("node" + depth, n.toString(true, true, false));
    JSDocInfo info = n.getJSDocInfo();
    if (info != null) {
      rec("jsdoc" + depth, info.toStringVerbose());
    }
    for (Node c = n.getFirstChild(); c != null; c = c.getNext()) {
      dumpNodes(c, depth + 1);
    }
  }

  static void tree(String name, Node root) {
    out.println("=== " + name);
    rec("tree", root.toStringTree());
    rec("str", root.toString());
    dumpNodes(root, 0);
    StringBuilder sb = new StringBuilder();
    try {
      root.appendJsonTree(sb);
    } catch (IOException e) {
      throw new IllegalStateException(e);
    }
    rec("json", sb.toString());
  }

  static void info(String name, JSDocInfo info) {
    out.println("=== " + name);
    rec("verbose", info == null ? "null" : info.toStringVerbose());
  }

  static JSTypeExpression type(String name) {
    return new JSTypeExpression(Node.newString(name), "src.js");
  }

  static JSTypeExpression bangType(String name) {
    return new JSTypeExpression(new Node(Token.BANG, Node.newString(name)), "src.js");
  }

  public static void main(String[] args) throws Exception {
    out = new PrintStream(System.out, true, StandardCharsets.UTF_8);

    // ---- numbers
    double[] numbers = {
      0.0, 1.0, 0.5, 100.0, 1e21, 1e-7, 0.001, 1e7, 123456789.0, 1234567.125,
      1.7976931348623157e308, 4.9e-324, 2.5e-5, 9007199254740993.0, 0.1, 1e300
    };
    Node numbers0 = new Node(Token.ARRAYLIT);
    for (double d : numbers) {
      numbers0.addChildToBack(IR.number(d));
    }
    tree("numbers", numbers0);
    for (double d : numbers) {
      rec("dtoa", DToA.numberToString(d));
    }

    // ---- strings
    Node strings = new Node(Token.ARRAYLIT);
    String[] strs = {"", "plain", "héllo", "😀", "lone\ud800x", "tab\tnl\nq\"b\\", " "};
    for (String s : strs) {
      strings.addChildToBack(IR.string(s));
    }
    strings.addChildToBack(Node.newTemplateLitString("cooked", "raw"));
    strings.addChildToBack(Node.newTemplateLitString(null, "\\u{bad"));
    strings.addChildToBack(Node.newBigInt(new BigInteger("12345678901234567890")));
    tree("strings", strings);

    // ---- a function with source positions, lengths and original names
    Node fnName = IR.name("f");
    Node param = IR.name("p");
    Node fn = IR.function(fnName, IR.paramList(param), IR.block(IR.returnNode(IR.name("p"))));
    fn.setLinenoCharno(3, 4);
    fn.setLength(27);
    fnName.setLinenoCharno(3, 13);
    fnName.setLength(1);
    param.setLinenoCharno(3, 15);
    param.setOriginalName("origP");
    fn.setSourceFileForTesting("src.js");
    fnName.setStaticSourceFileFrom(fn);
    param.setInputId(new InputId("src.js"));
    Node script = IR.script(fn);
    script.setLinenoCharno(1, 0);
    script.setLength(100);
    script.setStaticSourceFile(new SimpleSourceFile("ext.js", StaticSourceFile.SourceKind.EXTERN));
    script.putProp(Node.FEATURE_SET, FeatureSet.ES2015);
    script.putBooleanProp(Node.GOOG_MODULE, true);
    script.setUseStrict(true);
    tree("function", IR.root(script));
    rec("location", fnName.getLocation());
    rec("sourceOffset", String.valueOf(fn.getSourceOffset()));

    // ---- invalid FUNCTION, charno clamping, positions off
    Node badFn = new Node(Token.FUNCTION);
    badFn.setLinenoCharno(7, 5000);
    Node unpositioned = IR.name("u");
    unpositioned.setLinenoCharno(-1, 3);
    unpositioned.setLength(4);
    tree("misc", new Node(Token.BLOCK, badFn, unpositioned));

    // ---- every boolean / int prop reachable through the public API
    Node props = new Node(Token.BLOCK);
    Node paren = IR.name("paren");
    paren.setIsParenthesized(true);
    props.addChildToBack(IR.exprResult(paren));
    props.addChildToBack(IR.exprResult(IR.inc(IR.name("i"), true)));
    props.addChildToBack(IR.exprResult(IR.dec(IR.name("j"), false)));
    props.addChildToBack(
        IR.exprResult(IR.objectlit(IR.quotedStringKey("q", IR.number(1)), IR.stringKey("s", IR.number(2)))));
    Node synthetic = IR.block();
    synthetic.setIsSyntheticBlock(true);
    synthetic.setIsAddedBlock(true);
    props.addChildToBack(synthetic);
    Node call = IR.call(IR.name("g"), IR.name("a"));
    call.setSideEffectFlags(5);
    call.putBooleanProp(Node.FREE_CALL, true);
    call.putBooleanProp(Node.DIRECT_EVAL, true);
    props.addChildToBack(IR.exprResult(call));
    Node call2 = IR.call(IR.name("h"));
    call2.setSideEffectFlags(new Node.SideEffectFlags().clearAllFlags().setMutatesThis());
    props.addChildToBack(IR.exprResult(call2));
    Node constName = IR.name("C");
    constName.putBooleanProp(Node.IS_CONSTANT_NAME, true);
    constName.putBooleanProp(Node.IS_NAMESPACE, true);
    constName.setDeclaredConstantVar(true);
    constName.setInferredConstantVar(true);
    constName.putBooleanProp(Node.MODULE_ALIAS, true);
    constName.putBooleanProp(Node.MODULE_EXPORT, true);
    constName.setUnusedParameter(true);
    constName.makeNonIndexable();
    constName.setDeleted(true);
    constName.setChangeTime(7);
    props.addChildToBack(IR.var(constName, IR.number(3)));
    Node gen = IR.function(IR.name("gen"), IR.paramList(), IR.block());
    gen.setIsGeneratorFunction(true);
    gen.setIsAsyncFunction(true);
    gen.setGeneratorMarker(true);
    gen.setGeneratorSafe(true);
    props.addChildToBack(gen);
    Node arrow = IR.arrowFunction(IR.name(""), IR.paramList(), IR.number(1));
    props.addChildToBack(IR.exprResult(arrow));
    Node yieldAll = IR.yieldNode(IR.name("it"));
    yieldAll.setYieldAll(true);
    props.addChildToBack(IR.exprResult(yieldAll));
    Node member = IR.memberFunctionDef("m", IR.function(IR.name(""), IR.paramList(), IR.block()));
    member.setStaticMember(true);
    Node computed = IR.computedProp(IR.string("k"), IR.number(1));
    computed.putBooleanProp(Node.COMPUTED_PROP_METHOD, true);
    computed.putBooleanProp(Node.COMPUTED_PROP_GETTER, true);
    computed.putBooleanProp(Node.COMPUTED_PROP_SETTER, true);
    computed.putBooleanProp(Node.COMPUTED_PROP_VARIABLE, true);
    Node classMembers = new Node(Token.CLASS_MEMBERS, member);
    classMembers.addChildToBack(computed);
    Node clazz = new Node(Token.CLASS, IR.name("K"), IR.empty(), classMembers);
    clazz.putBooleanProp(Node.IS_TYPESCRIPT_ABSTRACT, true);
    clazz.putBooleanProp(Node.REFLECTED_OBJECT, true);
    props.addChildToBack(clazz);
    Node exportDefault = new Node(Token.EXPORT);
    exportDefault.putBooleanProp(Node.EXPORT_DEFAULT, true);
    exportDefault.putBooleanProp(Node.EXPORT_ALL_FROM, true);
    exportDefault.putBooleanProp(Node.ES6_MODULE, true);
    props.addChildToBack(exportDefault);
    Node shorthand = IR.stringKey("sh", IR.name("sh"));
    shorthand.setShorthandProperty(true);
    props.addChildToBack(IR.exprResult(IR.objectlit(shorthand)));
    Node optChain = IR.startOptChainGetprop(IR.name("o"), "p");
    props.addChildToBack(IR.exprResult(IR.continueOptChainCall(optChain)));
    Node arr = IR.arraylit(IR.number(1));
    arr.setTrailingComma(true);
    props.addChildToBack(IR.exprResult(arr));
    Node synthVar = IR.var(IR.name("synth"));
    synthVar.setIsSynthesizedUnfulfilledNameDeclaration(true);
    props.addChildToBack(synthVar);
    Node privateId = Node.newString(Token.STRING_KEY, "#priv");
    privateId.setPrivateIdentifier();
    props.addChildToBack(IR.exprResult(IR.objectlit(privateId)));
    Node typed = IR.name("typed");
    typed.putBooleanProp(Node.OPT_ES6_TYPED, true);
    typed.putBooleanProp(Node.CONSTRUCT_SIGNATURE, true);
    typed.putProp(Node.ACCESS_MODIFIER, "public");
    typed.putProp(Node.GENERIC_TYPE_LIST, new Node(Token.GENERIC_TYPE_LIST, IR.name("T")));
    typed.putProp(Node.IMPLEMENTS, new Node(Token.IMPLEMENTS));
    typed.setDeclaredTypeExpression(new Node(Token.NUMBER_TYPE));
    props.addChildToBack(IR.exprResult(typed));
    Node colored = IR.name("colored");
    colored.setColor(StandardColors.NUMBER);
    colored.setColorFromTypeCast();
    props.addChildToBack(IR.exprResult(colored));
    Node shadowHost = IR.name("host");
    shadowHost.setClosureUnawareShadow(IR.root(IR.script(IR.exprResult(IR.name("shadowed")))));
    props.addChildToBack(IR.exprResult(shadowHost));
    Node unaware = IR.name("unaware");
    unaware.setSourceFileForTesting("unaware.js");
    unaware.setIsInClosureUnawareSubtree(true);
    props.addChildToBack(IR.exprResult(unaware));
    Node comment = IR.name("commented");
    comment.setNonJSDocComment(
        new NonJSDocComment(new SourcePosition(null, 0, 1, 2), new SourcePosition(null, 5, 1, 7), "// c"));
    comment.setTrailingNonJSDocComment(
        new NonJSDocComment(new SourcePosition(null, 8, 1, 9), new SourcePosition(null, 12, 1, 13), "/* t */"));
    props.addChildToBack(IR.exprResult(comment));
    Node odd = IR.name("odd");
    odd.putProp(Node.FEATURE_SET, Double.NaN);
    odd.putProp(Node.ACCESS_MODIFIER, -0.0);
    odd.putIntProp(Node.QUOTED_PROP, -42);
    props.addChildToBack(IR.exprResult(odd));
    tree("props", props);

    // ---- JSDoc on nodes and standalone
    JSDocInfo.Builder b = JSDocInfo.builder();
    b.recordType(bangType("Foo"));
    b.recordVisibility(JSDocInfo.Visibility.PRIVATE);
    b.recordConstancy();
    b.recordDeprecated();
    b.recordDeprecationReason("old");
    b.recordSuppressions(ImmutableSet.of("visibility", "checkTypes"), "because");
    b.recordSuppression("missingRequire");
    b.recordNoInline();
    b.recordExport();
    b.recordClosurePrimitiveId("asserts.fail");
    b.recordIdGenerator();
    b.recordMods("ns");
    b.recordLicense("MIT");
    b.addLicense(" extra");
    b.recordTsType("number");
    b.recordClosureUnawareCode(JSDocInfo.PerFileClosureUnawareMode.WHITESPACE);
    JSDocInfo typeInfo = b.build();
    Node withDoc = IR.var(IR.name("x"), IR.number(1));
    withDoc.setJSDocInfo(typeInfo);
    Node castNode = IR.cast(IR.name("y"), typeInfo);
    tree("jsdocNodes", IR.block(withDoc, IR.exprResult(castNode)));

    JSDocInfo.Builder f = JSDocInfo.builder().parseDocumentation();
    f.recordParameter("a", type("string"));
    f.recordParameter("b", new JSTypeExpression(new Node(Token.EQUALS, Node.newString("number")), "src.js"));
    f.recordParameterDescription("a", "the a");
    f.recordReturnType(type("boolean"));
    f.recordReturnDescription("whether");
    f.recordThisType(bangType("Bar"));
    f.recordTemplateTypeName("T");
    f.recordTemplateTypeName("U", type("Object"));
    f.recordTypeTransformation("R", IR.name("ttl"));
    f.recordThrowsAnnotation("{Error} when bad");
    f.recordAuthor("me");
    f.recordReference("http://example.com");
    f.recordDescription("desc");
    f.recordMeaning("meaning");
    f.recordAlternateMessageId("alt");
    f.recordModifies(ImmutableSet.of("this", "arguments"));
    f.recordBlockDescription("block");
    f.recordOriginalCommentString("/** original */");
    f.recordOriginalCommentPosition(42);
    f.recordNoSideEffects();
    f.recordImplicitCast();
    f.recordStableIdGenerator();
    f.recordInlineType();
    JSDocInfo fnInfo = f.build();
    info("functionDoc", fnInfo);
    rec("paramNames", fnInfo.getParameterNames().toString());
    rec("paramB", fnInfo.getParameterType("b").toString());
    rec("templates", fnInfo.getTemplateTypeNames().toString());
    rec("typeNodes", String.valueOf(fnInfo.getTypeNodes().size()));
    rec("modifiesThis", String.valueOf(fnInfo.modifiesThis()));
    rec("suppressions", typeInfo.getSuppressions().toString());
    rec("suppressionsMap", typeInfo.getSuppressionsAndTheirDescription().toString());
    rec("visibility", String.valueOf(typeInfo.getVisibility()));
    rec("license", typeInfo.getLicense());
    rec("clone", fnInfo.clone(true).toStringVerbose());
    rec("equivalent", String.valueOf(JSDocInfo.areEquivalent(fnInfo, fnInfo.clone())));

    JSDocInfo.Builder c = JSDocInfo.builder();
    c.recordConstructor();
    c.recordBaseType(bangType("Base"));
    c.recordImplementedInterface(bangType("I1"));
    c.recordImplementedInterface(bangType("I2"));
    c.recordImplementedInterface(bangType("I1"));
    c.recordStruct();
    c.recordFinality();
    c.recordFileOverview("ignored without documentation");
    c.recordEnhance("enh");
    c.recordPolymer();
    c.recordCustomElement();
    c.recordNgInject(true);
    c.recordPureOrBreakMyCode();
    c.recordLogTypeInCompiler();
    info("ctorDoc", c.build());

    JSDocInfo.Builder i = JSDocInfo.builder();
    i.recordInterface();
    i.recordExtendedInterface(bangType("Parent"));
    i.recordLends(type("Lent"));
    i.recordVisibility(JSDocInfo.Visibility.PACKAGE);
    info("interfaceDoc", i.build());

    JSDocInfo.Builder e = JSDocInfo.builder();
    e.recordEnumParameterType(type("number"));
    info("enumDoc", e.build());
    JSDocInfo.Builder t = JSDocInfo.builder();
    t.recordTypedef(type("Object"));
    t.recordXidGenerator();
    info("typedefDoc", t.build());
    JSDocInfo.Builder d = JSDocInfo.builder();
    d.recordDefineType(type("boolean"));
    d.recordDescription("define");
    info("defineDoc", d.build());
    info("emptyAlways", JSDocInfo.builder().build(true));
    info("emptyNull", JSDocInfo.builder().build());
    JSDocInfo.Builder r = JSDocInfo.builder();
    r.recordImplicitMatch();
    r.recordUnrestricted();
    r.recordDict();
    r.recordNoCollapse();
    r.recordNoCompile();
    r.recordTypeSummary();
    r.recordExterns();
    r.recordNoCoverage();
    r.recordOverride();
    r.recordProvideGoog();
    r.recordWizaction();
    r.recordWizcallback();
    r.recordMixinClass();
    r.recordMixinFunction();
    r.recordPolymerBehavior();
    r.recordSassGeneratedCssTs();
    r.recordUsedViaDotConstructor();
    r.recordNoDts();
    r.recordProvideAlreadyProvided();
    r.recordCollapsibleOrBreakMyCode();
    r.recordEncourageInlining();
    r.recordConsistentIdGenerator();
    r.recordMappedIdGenerator();
    info("bitsDoc", r.build());
    JSDocInfo.Builder fo = JSDocInfo.builder().parseDocumentation();
    fo.recordFileOverview("the file");
    fo.recordClosureUnawareCode();
    info("fileoverviewDoc", fo.build());

    // ---- colors
    out.println("=== colors");
    rec("number", StandardColors.NUMBER.toString());
    rec("unknown", StandardColors.UNKNOWN.toString());
    LinkedHashSet<Color> elems = new LinkedHashSet<>();
    elems.add(StandardColors.NUMBER);
    elems.add(StandardColors.STRING);
    elems.add(StandardColors.NULL_OR_VOID);
    Color union = Color.createUnion(elems);
    rec("union", union.toString());
    rec("unionMinusNull", union.subtractNullOrVoid().toString());
    Color custom =
        Color.singleBuilder()
            .setId(ColorId.fromAscii("Foo"))
            .setOwnProperties(ImmutableSet.of("a", "b"))
            .setInstanceColor(StandardColors.NUMBER)
            .setConstructor(true)
            .build();
    rec("custom", custom.toString());
    rec("customHash", String.valueOf(custom.hashCode()));
    rec("unionHash", String.valueOf(union.hashCode()));
    Node colorTree = IR.name("c");
    colorTree.setColor(union);
    tree("colorNode", colorTree);

    // ---- qualified names and misc accessors
    out.println("=== qnames");
    Node q = IR.getprop(IR.getprop(IR.thisNode(), "a"), "b");
    q.getFirstChild().setOriginalName("origA");
    rec("qname", q.getQualifiedName());
    rec("origQname", q.getOriginalQualifiedName());
    rec("qobj", q.getQualifiedNameObject().join());
    rec("of", String.join("|", QualifiedName.of("x.y.z").components()));
    rec("propList", constName.getPropListDebugString());
    rec("serialized", Long.toHexString(constName.serializeProperties()));
    rec("sideEffects", new Node.SideEffectFlags(5).toString());
    rec("tokens", String.valueOf(Token.values().length));
    rec("msg", Msg.JSDOC_TYPETRANSFORMATION_INVALID.format("x"));
  }
}
