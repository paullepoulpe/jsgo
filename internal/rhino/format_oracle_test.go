package rhino

// Differential check of the formatting code against Java: this test builds
// the values that tools/oracle/java/RhinoFormatOracle.java builds, with the
// ported API, and must reproduce its output (testdata/format_oracle.txt)
// byte for byte. Keep the two files in step.

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/dtoa"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
)

type formatOracle struct {
	sb strings.Builder
}

func oracleEsc(s string) string {
	var sb strings.Builder
	for _, c := range jutil.StringToUTF16(s) {
		switch {
		case c == '\\':
			sb.WriteString("\\\\")
		case c == '\n':
			sb.WriteString("\\n")
		case c < 0x20 || c >= 0x7f:
			fmt.Fprintf(&sb, "\\u%04x", c)
		default:
			sb.WriteByte(byte(c))
		}
	}
	return sb.String()
}

func (o *formatOracle) line(s string) { o.sb.WriteString(s + "\n") }

func (o *formatOracle) rec(label, value string) { o.line(label + "\t" + oracleEsc(value)) }

func (o *formatOracle) dumpNodes(n Node, depth int) {
	o.rec("node"+strconv.Itoa(depth), n.ToStringWithFlags(true, true, false))
	if info := n.GetJSDocInfo(); info != nil {
		o.rec("jsdoc"+strconv.Itoa(depth), info.ToStringVerbose())
	}
	for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
		o.dumpNodes(c, depth+1)
	}
}

func (o *formatOracle) tree(name string, root Node) {
	o.line("=== " + name)
	o.rec("tree", root.ToStringTree())
	o.rec("str", root.String())
	o.dumpNodes(root, 0)
	var sb strings.Builder
	root.AppendJsonTree(&sb)
	o.rec("json", sb.String())
}

func (o *formatOracle) info(name string, info *JSDocInfo) {
	o.line("=== " + name)
	if info == nil {
		o.rec("verbose", "null")
	} else {
		o.rec("verbose", info.ToStringVerbose())
	}
}

func oracleType(name string) *JSTypeExpression {
	return NewJSTypeExpression(Node_NewString(name), "src.js")
}

func oracleBangType(name string) *JSTypeExpression {
	return NewJSTypeExpression(NewNodeWithChild(Token_BANG, Node_NewString(name)), "src.js")
}

func set(elems ...string) *jutil.LinkedHashSet[string] { return jutil.NewLinkedHashSet(elems...) }

func buildFormatOracle() string {
	o := &formatOracle{}

	// ---- numbers
	numbers := []float64{
		0.0, 1.0, 0.5, 100.0, 1e21, 1e-7, 0.001, 1e7, 123456789.0, 1234567.125,
		1.7976931348623157e308, 4.9e-324, 2.5e-5, 9007199254740993.0, 0.1, 1e300,
	}
	numbers0 := NewNode(Token_ARRAYLIT)
	for _, d := range numbers {
		numbers0.AddChildToBack(IR_Number(d))
	}
	o.tree("numbers", numbers0)
	for _, d := range numbers {
		o.rec("dtoa", dtoa.DToA_NumberToString(d))
	}

	// ---- strings
	strs := []string{"", "plain", "h\u00e9llo", "\U0001F600",
		"lone" + jutil.StringFromChar(0xD800) + "x", "tab\tnl\nq\"b\\", "\u2028"}
	stringsNode := NewNode(Token_ARRAYLIT)
	for _, s := range strs {
		stringsNode.AddChildToBack(IR_String(s))
	}
	stringsNode.AddChildToBack(Node_NewTemplateLitString(jutil.StringPtr("cooked"), "raw"))
	stringsNode.AddChildToBack(Node_NewTemplateLitString(nil, "\\u{bad"))
	stringsNode.AddChildToBack(Node_NewBigInt(jutil.NewBigIntegerFromString("12345678901234567890")))
	o.tree("strings", stringsNode)

	// ---- a function with source positions, lengths and original names
	fnName := IR_Name("f")
	param := IR_Name("p")
	fn := IR_Function(fnName, IR_ParamList(param), IR_BlockOfNode(IR_ReturnNodeWithExpr(IR_Name("p"))))
	fn.SetLinenoCharno(3, 4)
	fn.SetLength(27)
	fnName.SetLinenoCharno(3, 13)
	fnName.SetLength(1)
	param.SetLinenoCharno(3, 15)
	param.SetOriginalName(jutil.StringPtr("origP"))
	fn.SetSourceFileForTesting("src.js")
	fnName.SetStaticSourceFileFrom(fn)
	param.SetInputId(NewInputId("src.js"))
	script := IR_ScriptOfNodes(fn)
	script.SetLinenoCharno(1, 0)
	script.SetLength(100)
	script.SetStaticSourceFile(NewSimpleSourceFile("ext.js", StaticSourceFile_SourceKind_EXTERN))
	script.PutProp(Node_FEATURE_SET, parser.FeatureSet_ES2015)
	script.PutBooleanProp(Node_GOOG_MODULE, true)
	script.SetUseStrict(true)
	o.tree("function", IR_Root(script))
	o.rec("location", fnName.GetLocation())
	o.rec("sourceOffset", strconv.Itoa(fn.GetSourceOffset()))

	// ---- invalid FUNCTION, charno clamping, positions off
	badFn := NewNode(Token_FUNCTION)
	badFn.SetLinenoCharno(7, 5000)
	unpositioned := IR_Name("u")
	unpositioned.SetLinenoCharno(-1, 3)
	unpositioned.SetLength(4)
	o.tree("misc", NewNodeWith2Children(Token_BLOCK, badFn, unpositioned))

	// ---- every boolean / int prop reachable through the public API
	props := NewNode(Token_BLOCK)
	paren := IR_Name("paren")
	paren.SetIsParenthesized(true)
	props.AddChildToBack(IR_ExprResult(paren))
	props.AddChildToBack(IR_ExprResult(IR_Inc(IR_Name("i"), true)))
	props.AddChildToBack(IR_ExprResult(IR_Dec(IR_Name("j"), false)))
	props.AddChildToBack(
		IR_ExprResult(IR_Objectlit(IR_QuotedStringKey("q", IR_Number(1)), IR_StringKeyWithValue("s", IR_Number(2)))))
	synthetic := IR_Block()
	synthetic.SetIsSyntheticBlock(true)
	synthetic.SetIsAddedBlock(true)
	props.AddChildToBack(synthetic)
	call := IR_Call(IR_Name("g"), IR_Name("a"))
	call.SetSideEffectFlags(5)
	call.PutBooleanProp(Node_FREE_CALL, true)
	call.PutBooleanProp(Node_DIRECT_EVAL, true)
	props.AddChildToBack(IR_ExprResult(call))
	call2 := IR_Call(IR_Name("h"))
	call2.SetSideEffectFlagsObject(NewNode_SideEffectFlags().ClearAllFlags().SetMutatesThis())
	props.AddChildToBack(IR_ExprResult(call2))
	constName := IR_Name("C")
	constName.PutBooleanProp(Node_IS_CONSTANT_NAME, true)
	constName.PutBooleanProp(Node_IS_NAMESPACE, true)
	constName.SetDeclaredConstantVar(true)
	constName.SetInferredConstantVar(true)
	constName.PutBooleanProp(Node_MODULE_ALIAS, true)
	constName.PutBooleanProp(Node_MODULE_EXPORT, true)
	constName.SetUnusedParameter(true)
	constName.MakeNonIndexable()
	constName.SetDeleted(true)
	constName.SetChangeTime(7)
	props.AddChildToBack(IR_VarWithValue(constName, IR_Number(3)))
	gen := IR_Function(IR_Name("gen"), IR_ParamList(), IR_Block())
	gen.SetIsGeneratorFunction(true)
	gen.SetIsAsyncFunction(true)
	gen.SetGeneratorMarker(true)
	gen.SetGeneratorSafe(true)
	props.AddChildToBack(gen)
	arrow := IR_ArrowFunction(IR_Name(""), IR_ParamList(), IR_Number(1))
	props.AddChildToBack(IR_ExprResult(arrow))
	yieldAll := IR_YieldNode(IR_Name("it"))
	yieldAll.SetYieldAll(true)
	props.AddChildToBack(IR_ExprResult(yieldAll))
	member := IR_MemberFunctionDef("m", IR_Function(IR_Name(""), IR_ParamList(), IR_Block()))
	member.SetStaticMember(true)
	computed := IR_ComputedProp(IR_String("k"), IR_Number(1))
	computed.PutBooleanProp(Node_COMPUTED_PROP_METHOD, true)
	computed.PutBooleanProp(Node_COMPUTED_PROP_GETTER, true)
	computed.PutBooleanProp(Node_COMPUTED_PROP_SETTER, true)
	computed.PutBooleanProp(Node_COMPUTED_PROP_VARIABLE, true)
	classMembers := NewNodeWithChild(Token_CLASS_MEMBERS, member)
	classMembers.AddChildToBack(computed)
	clazz := NewNodeWith3Children(Token_CLASS, IR_Name("K"), IR_Empty(), classMembers)
	clazz.PutBooleanProp(Node_IS_TYPESCRIPT_ABSTRACT, true)
	clazz.PutBooleanProp(Node_REFLECTED_OBJECT, true)
	props.AddChildToBack(clazz)
	exportDefault := NewNode(Token_EXPORT)
	exportDefault.PutBooleanProp(Node_EXPORT_DEFAULT, true)
	exportDefault.PutBooleanProp(Node_EXPORT_ALL_FROM, true)
	exportDefault.PutBooleanProp(Node_ES6_MODULE, true)
	props.AddChildToBack(exportDefault)
	shorthand := IR_StringKeyWithValue("sh", IR_Name("sh"))
	shorthand.SetShorthandProperty(true)
	props.AddChildToBack(IR_ExprResult(IR_Objectlit(shorthand)))
	optChain := IR_StartOptChainGetprop(IR_Name("o"), "p")
	props.AddChildToBack(IR_ExprResult(IR_ContinueOptChainCall(optChain)))
	arr := IR_Arraylit(IR_Number(1))
	arr.SetTrailingComma(true)
	props.AddChildToBack(IR_ExprResult(arr))
	synthVar := IR_Var(IR_Name("synth"))
	synthVar.SetIsSynthesizedUnfulfilledNameDeclaration(true)
	props.AddChildToBack(synthVar)
	privateId := Node_NewStringWithToken(Token_STRING_KEY, "#priv")
	privateId.SetPrivateIdentifier()
	props.AddChildToBack(IR_ExprResult(IR_Objectlit(privateId)))
	typed := IR_Name("typed")
	typed.PutBooleanProp(Node_OPT_ES6_TYPED, true)
	typed.PutBooleanProp(Node_CONSTRUCT_SIGNATURE, true)
	typed.PutProp(Node_ACCESS_MODIFIER, "public")
	typed.PutProp(Node_GENERIC_TYPE_LIST, NewNodeWithChild(Token_GENERIC_TYPE_LIST, IR_Name("T")))
	typed.PutProp(Node_IMPLEMENTS, NewNode(Token_IMPLEMENTS))
	typed.SetDeclaredTypeExpression(NewNode(Token_NUMBER_TYPE))
	props.AddChildToBack(IR_ExprResult(typed))
	colored := IR_Name("colored")
	colored.SetColor(colors.StandardColors_NUMBER)
	colored.SetColorFromTypeCast()
	props.AddChildToBack(IR_ExprResult(colored))
	shadowHost := IR_Name("host")
	shadowHost.SetClosureUnawareShadow(IR_Root(IR_ScriptOfNodes(IR_ExprResult(IR_Name("shadowed")))))
	props.AddChildToBack(IR_ExprResult(shadowHost))
	unaware := IR_Name("unaware")
	unaware.SetSourceFileForTesting("unaware.js")
	unaware.SetIsInClosureUnawareSubtree(true)
	props.AddChildToBack(IR_ExprResult(unaware))
	comment := IR_Name("commented")
	comment.SetNonJSDocComment(
		NewNonJSDocComment(parser.LinePosition{Offset: 0, Line: 1, Column: 2}, parser.LinePosition{Offset: 5, Line: 1, Column: 7}, jutil.StringPtr("// c")))
	comment.SetTrailingNonJSDocComment(
		NewNonJSDocComment(parser.LinePosition{Offset: 8, Line: 1, Column: 9}, parser.LinePosition{Offset: 12, Line: 1, Column: 13}, jutil.StringPtr("/* t */")))
	props.AddChildToBack(IR_ExprResult(comment))
	odd := IR_Name("odd")
	odd.PutProp(Node_FEATURE_SET, math.NaN())
	odd.PutProp(Node_ACCESS_MODIFIER, math.Copysign(0, -1))
	odd.PutIntProp(Node_QUOTED_PROP, -42)
	props.AddChildToBack(IR_ExprResult(odd))
	o.tree("props", props)

	// ---- JSDoc on nodes and standalone
	b := JSDocInfo_Builder_()
	b.RecordType(oracleBangType("Foo"))
	b.RecordVisibility(JSDocInfo_Visibility_PRIVATE)
	b.RecordConstancy()
	b.RecordDeprecated()
	b.RecordDeprecationReason("old")
	b.RecordSuppressionsWithDescription(set("visibility", "checkTypes"), "because")
	b.RecordSuppression("missingRequire")
	b.RecordNoInline()
	b.RecordExport()
	b.RecordClosurePrimitiveId("asserts.fail")
	b.RecordIdGenerator()
	b.RecordMods("ns")
	b.RecordLicense("MIT")
	b.AddLicense(" extra")
	b.RecordTsType("number")
	b.RecordClosureUnawareCodeWithMode(JSDocInfo_PerFileClosureUnawareMode_WHITESPACE)
	typeInfo := b.Build()
	withDoc := IR_VarWithValue(IR_Name("x"), IR_Number(1))
	withDoc.SetJSDocInfo(typeInfo)
	castNode := IR_Cast(IR_Name("y"), typeInfo)
	o.tree("jsdocNodes", IR_BlockOfNodes(withDoc, IR_ExprResult(castNode)))

	f := JSDocInfo_Builder_().ParseDocumentation()
	f.RecordParameter("a", oracleType("string"))
	f.RecordParameter("b", NewJSTypeExpression(NewNodeWithChild(Token_EQUALS, Node_NewString("number")), "src.js"))
	f.RecordParameterDescription("a", "the a")
	f.RecordReturnType(oracleType("boolean"))
	f.RecordReturnDescription("whether")
	f.RecordThisType(oracleBangType("Bar"))
	f.RecordTemplateTypeName("T")
	f.RecordTemplateTypeNameWithBound("U", oracleType("Object"))
	f.RecordTypeTransformation("R", IR_Name("ttl"))
	f.RecordThrowsAnnotation("{Error} when bad")
	f.RecordAuthor("me")
	f.RecordReference("http://example.com")
	f.RecordDescription("desc")
	f.RecordMeaning("meaning")
	f.RecordAlternateMessageId("alt")
	f.RecordModifies(set("this", "arguments"))
	f.RecordBlockDescription("block")
	f.RecordOriginalCommentString("/** original */")
	f.RecordOriginalCommentPosition(42)
	f.RecordNoSideEffects()
	f.RecordImplicitCast()
	f.RecordStableIdGenerator()
	f.RecordInlineType()
	fnInfo := f.Build()
	o.info("functionDoc", fnInfo)
	o.rec("paramNames", fnInfo.GetParameterNames().String())
	o.rec("paramB", fnInfo.GetParameterType("b").String())
	o.rec("templates", jutil.CollectionToString(fnInfo.GetTemplateTypeNames()))
	o.rec("typeNodes", strconv.Itoa(len(fnInfo.GetTypeNodes())))
	o.rec("modifiesThis", jutil.ToJavaString(fnInfo.ModifiesThis()))
	o.rec("suppressions", typeInfo.GetSuppressions().String())
	o.rec("suppressionsMap", FormatJSDocSuppressions(typeInfo.GetSuppressionsAndTheirDescription()))
	o.rec("visibility", typeInfo.GetVisibility().String())
	o.rec("license", *typeInfo.GetLicense())
	o.rec("clone", fnInfo.CloneWithTypeNodes(true).ToStringVerbose())
	o.rec("equivalent", jutil.ToJavaString(JSDocInfo_AreEquivalent(fnInfo, fnInfo.Clone())))

	c := JSDocInfo_Builder_()
	c.RecordConstructor()
	c.RecordBaseType(oracleBangType("Base"))
	c.RecordImplementedInterface(oracleBangType("I1"))
	c.RecordImplementedInterface(oracleBangType("I2"))
	c.RecordImplementedInterface(oracleBangType("I1"))
	c.RecordStruct()
	c.RecordFinality()
	c.RecordFileOverview("ignored without documentation")
	c.RecordEnhance("enh")
	c.RecordPolymer()
	c.RecordCustomElement()
	c.RecordNgInject(true)
	c.RecordPureOrBreakMyCode()
	c.RecordLogTypeInCompiler()
	o.info("ctorDoc", c.Build())

	i := JSDocInfo_Builder_()
	i.RecordInterface()
	i.RecordExtendedInterface(oracleBangType("Parent"))
	i.RecordLends(oracleType("Lent"))
	i.RecordVisibility(JSDocInfo_Visibility_PACKAGE)
	o.info("interfaceDoc", i.Build())

	e := JSDocInfo_Builder_()
	e.RecordEnumParameterType(oracleType("number"))
	o.info("enumDoc", e.Build())
	t := JSDocInfo_Builder_()
	t.RecordTypedef(oracleType("Object"))
	t.RecordXidGenerator()
	o.info("typedefDoc", t.Build())
	d := JSDocInfo_Builder_()
	d.RecordDefineType(oracleType("boolean"))
	d.RecordDescription("define")
	o.info("defineDoc", d.Build())
	o.info("emptyAlways", JSDocInfo_Builder_().BuildWithAlways(true))
	o.info("emptyNull", JSDocInfo_Builder_().Build())
	r := JSDocInfo_Builder_()
	r.RecordImplicitMatch()
	r.RecordUnrestricted()
	r.RecordDict()
	r.RecordNoCollapse()
	r.RecordNoCompile()
	r.RecordTypeSummary()
	r.RecordExterns()
	r.RecordNoCoverage()
	r.RecordOverride()
	r.RecordProvideGoog()
	r.RecordWizaction()
	r.RecordWizcallback()
	r.RecordMixinClass()
	r.RecordMixinFunction()
	r.RecordPolymerBehavior()
	r.RecordSassGeneratedCssTs()
	r.RecordUsedViaDotConstructor()
	r.RecordNoDts()
	r.RecordProvideAlreadyProvided()
	r.RecordCollapsibleOrBreakMyCode()
	r.RecordEncourageInlining()
	r.RecordConsistentIdGenerator()
	r.RecordMappedIdGenerator()
	o.info("bitsDoc", r.Build())
	fo := JSDocInfo_Builder_().ParseDocumentation()
	fo.RecordFileOverview("the file")
	fo.RecordClosureUnawareCode()
	o.info("fileoverviewDoc", fo.Build())

	// ---- colors
	o.line("=== colors")
	o.rec("number", colors.StandardColors_NUMBER.String())
	o.rec("unknown", colors.StandardColors_UNKNOWN.String())
	elems := []*colors.Color{colors.StandardColors_NUMBER, colors.StandardColors_STRING, colors.StandardColors_NULL_OR_VOID}
	union := colors.Color_CreateUnionSlice(elems)
	o.rec("union", union.String())
	o.rec("unionMinusNull", union.SubtractNullOrVoid().String())
	custom := colors.Color_SingleBuilder().
		SetId(colors.ColorId_FromAscii("Foo")).
		SetOwnProperties([]string{"a", "b"}).
		SetInstanceColor(colors.StandardColors_NUMBER).
		SetConstructor(true).
		Build()
	o.rec("custom", custom.String())
	o.rec("customHash", jutil.ToJavaString(custom.HashCode()))
	o.rec("unionHash", jutil.ToJavaString(union.HashCode()))
	colorTree := IR_Name("c")
	colorTree.SetColor(union)
	o.tree("colorNode", colorTree)

	// ---- qualified names and misc accessors
	o.line("=== qnames")
	q := IR_Getprop(IR_Getprop(IR_ThisNode(), "a"), "b")
	q.GetFirstChild().SetOriginalName(jutil.StringPtr("origA"))
	o.rec("qname", *q.GetQualifiedName())
	o.rec("origQname", *q.GetOriginalQualifiedName())
	o.rec("qobj", q.GetQualifiedNameObject().Join())
	o.rec("of", strings.Join(QualifiedName_Of("x.y.z").Components(), "|"))
	o.rec("propList", constName.GetPropListDebugString())
	o.rec("serialized", jutil.LongToHexString(constName.SerializeProperties()))
	o.rec("sideEffects", NewNode_SideEffectFlagsWithValue(5).String())
	o.rec("tokens", strconv.Itoa(len(Token_values())))
	o.rec("msg", Msg_JSDOC_TYPETRANSFORMATION_INVALID.FormatArgs("x"))

	return o.sb.String()
}

func TestFormatMatchesJava(t *testing.T) {
	want, err := os.ReadFile("testdata/format_oracle.txt")
	if err != nil {
		t.Fatal(err)
	}
	got := buildFormatOracle()
	if got == string(want) {
		return
	}
	gl, wl := strings.Split(got, "\n"), strings.Split(string(want), "\n")
	diffs := 0
	for i := 0; i < len(gl) || i < len(wl); i++ {
		var g, w string
		if i < len(gl) {
			g = gl[i]
		}
		if i < len(wl) {
			w = wl[i]
		}
		if g != w {
			diffs++
			if diffs <= 15 {
				t.Errorf("line %d:\n go:   %s\n java: %s", i+1, g, w)
			}
		}
	}
	t.Fatalf("%d differing lines", diffs)
}
