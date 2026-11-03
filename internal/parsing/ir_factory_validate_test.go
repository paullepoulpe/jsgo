package parsing

import (
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type recordingErrorReporter struct {
	errors   []string
	warnings []string
}

func (r *recordingErrorReporter) Error(message string, sourceName string, line int, lineOffset int) {
	r.errors = append(r.errors, message)
}

func (r *recordingErrorReporter) Warning(message string, sourceName string, line int, lineOffset int) {
	r.warnings = append(r.warnings, message)
}

func TestIRFactoryValidationDetachedAndShallowAncestors(t *testing.T) {
	prev := rhino.SetCurrentNodeStore(rhino.NewNodeStore())
	defer rhino.SetCurrentNodeStore(prev)

	reporter := &recordingErrorReporter{}
	sf := rhino.NewSimpleSourceFile("test.js", rhino.StaticSourceFile_SourceKind_STRONG)
	cfg := ParserRunner_CreateConfigWithModes(
		Config_LanguageMode_ECMASCRIPT5,
		Config_JsDocParsing_TYPES_ONLY,
		Config_RunMode_STOP_AFTER_ERROR,
		nil,
		true,
		Config_StrictMode_STRICT,
	)

	docBuilder := rhino.NewJSDocInfo_Builder()
	docBuilder.RecordClosureUnawareCode()
	closureUnawareDoc := docBuilder.Build()

	factory := &IRFactory{
		sourceName:        "test.js",
		sourceFile:        sf,
		config:            cfg,
		errorReporter:     reporter,
		firstFileoverview: closureUnawareDoc,
	}

	// 1. ValidateShadowHost on detached shadow host (parent == NilNode) and
	// parent GETPROP with no grandparent must report an error cleanly instead of SIGSEGV.
	shadowHost := rhino.IR_Name("host")
	shadowHost.SetClosureUnawareShadow(rhino.IR_Script())
	factory.ValidateShadowHost(shadowHost)
	if len(reporter.errors) != 1 || reporter.errors[0] != "invalid @closureUnaware node: must be directly called" {
		t.Fatalf("unexpected ValidateShadowHost(detached) errors: %v", reporter.errors)
	}

	reporter.errors = nil
	getpropNoGrandparent := rhino.IR_Getprop(shadowHost, "call")
	factory.ValidateShadowHost(shadowHost)
	_ = getpropNoGrandparent
	if len(reporter.errors) != 1 || reporter.errors[0] != "invalid @closureUnaware node: must be directly called" {
		t.Fatalf("unexpected ValidateShadowHost(no grandparent) errors: %v", reporter.errors)
	}

	// 2. ValidateAwait when parent is a detached BLOCK or detached MEMBER_FIELD_DEF (no grandparent).
	reporter.errors = nil
	awaitInDetachedBlock := rhino.NewNodeWithChild(rhino.Token_AWAIT, rhino.IR_Number(1))
	detachedBlock := rhino.IR_BlockOfNode(rhino.IR_ExprResult(awaitInDetachedBlock))
	_ = detachedBlock
	factory.ValidateAwait(awaitInDetachedBlock)
	if len(reporter.errors) != 1 || reporter.errors[0] != IRFactory_UNEXPECTED_AWAIT {
		t.Fatalf("unexpected ValidateAwait(detached BLOCK) errors: %v", reporter.errors)
	}

	reporter.errors = nil
	awaitInDetachedField := rhino.NewNodeWithChild(rhino.Token_AWAIT, rhino.IR_Number(2))
	detachedField := rhino.Node_NewStringWithToken(rhino.Token_MEMBER_FIELD_DEF, "x")
	detachedField.AddChildToBack(awaitInDetachedField)
	factory.ValidateAwait(awaitInDetachedField)
	if len(reporter.errors) != 1 || reporter.errors[0] != IRFactory_UNEXPECTED_AWAIT {
		t.Fatalf("unexpected ValidateAwait(detached MEMBER_FIELD_DEF) errors: %v", reporter.errors)
	}

	// 3. ValidateBlockScopedFunctions on detached FUNCTION and FUNCTION in detached BLOCK (no grandparent).
	reporter.warnings = nil
	detachedFn := rhino.IR_Function(rhino.IR_Name("f"), rhino.IR_ParamList(), rhino.IR_Block())
	factory.ValidateBlockScopedFunctions(detachedFn)
	parentBlockNoGrandparent := rhino.IR_BlockOfNode(detachedFn)
	_ = parentBlockNoGrandparent
	factory.ValidateBlockScopedFunctions(detachedFn)
	if len(reporter.warnings) != 0 {
		t.Fatalf("expected 0 warnings on FUNCTION in parentless BLOCK, got %v", reporter.warnings)
	}
	outerScript := rhino.IR_ScriptOfNodes(parentBlockNoGrandparent)
	_ = outerScript
	factory.ValidateBlockScopedFunctions(detachedFn)
	if len(reporter.warnings) != 1 {
		t.Fatalf("expected 1 warning when BLOCK has non-FUNCTION grandparent, got %v", reporter.warnings)
	}

	// 4. validatePrivatePropertyUsage when classMembers directly holds a private GETPROP
	// and via full parse of undeclared private field + private method call.
	reporter.errors = nil
	privGetProp := rhino.IR_Getprop(rhino.NewNode(rhino.Token_THIS), "secret")
	privGetProp.PutBooleanProp(rhino.Node_Prop_PRIVATE_IDENTIFIER, true)
	classMembers := rhino.NewNodeWithChild(rhino.Token_CLASS_MEMBERS, privGetProp)
	factory.validatePrivatePropertyUsage(&privatePropScope{}, classMembers)
	if len(reporter.errors) != 1 || reporter.errors[0] != "Private fields must be declared in an enclosing class" {
		t.Fatalf("unexpected validatePrivatePropertyUsage errors: %v", reporter.errors)
	}

	parseReporter := &recordingErrorReporter{}
	esNextCfg := ParserRunner_CreateConfigWithModes(
		Config_LanguageMode_ES_NEXT,
		Config_JsDocParsing_TYPES_ONLY,
		Config_RunMode_KEEP_GOING,
		nil,
		true,
		Config_StrictMode_STRICT,
	)
	_ = ParserRunner_Parse(
		rhino.NewSimpleSourceFile("priv.js", rhino.StaticSourceFile_SourceKind_STRONG),
		"class C { m() { return this.#missingField + this.#missingMethod(); } }",
		esNextCfg,
		parseReporter,
	)
	wantParseErrors := []string{
		"Private methods must be declared in an enclosing class",
		"Private fields must be declared in an enclosing class",
	}
	for _, wantErr := range wantParseErrors {
		if !slices.Contains(parseReporter.errors, wantErr) {
			t.Fatalf("missing expected error %q in %v", wantErr, parseReporter.errors)
		}
	}
}

func TestParsingDefensiveClones(t *testing.T) {
	reserved1 := ParserRunner_GetReservedVars()
	origReservedLen := reserved1.Len()
	if origReservedLen == 0 {
		t.Fatal("expected non-empty reserved vars")
	}
	reserved1.Add("__mutated_reserved_var__")
	reserved2 := ParserRunner_GetReservedVars()
	if reserved2.Len() != origReservedLen || reserved2.Contains("__mutated_reserved_var__") {
		t.Fatalf("ParserRunner_GetReservedVars() returned shared mutable OrderedSet")
	}

	suppressions1 := ParserRunner_GetSuppressionNames()
	origSuppLen := suppressions1.Len()
	if origSuppLen == 0 {
		t.Fatal("expected non-empty suppression names")
	}
	suppressions1.Add("__mutated_suppression__")
	suppressions2 := ParserRunner_GetSuppressionNames()
	if suppressions2.Len() != origSuppLen || suppressions2.Contains("__mutated_suppression__") {
		t.Fatalf("ParserRunner_GetSuppressionNames() returned shared mutable OrderedSet")
	}

	modes1 := Config_JsDocParsing_values()
	if len(modes1) == 0 {
		t.Fatal("expected non-empty Config_JsDocParsing_values()")
	}
	origMode0 := modes1[0]
	modes1[0] = Config_JsDocParsing(-123)
	modes2 := Config_JsDocParsing_values()
	if modes2[0] != origMode0 {
		t.Fatalf("Config_JsDocParsing_values() returned shared backing slice: got %v, want %v", modes2[0], origMode0)
	}
}
