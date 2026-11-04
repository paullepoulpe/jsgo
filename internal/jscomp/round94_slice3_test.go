/*
 * Copyright 2026 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parser"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type round94ConvertChunksCase struct {
	CompilerTestCaseBase
}

func newRound94ConvertChunksCase() *round94ConvertChunksCase {
	c := &round94ConvertChunksCase{}
	c.InitCompilerTestCaseBase(c)
	return c
}

func (c *round94ConvertChunksCase) GetProcessor(compiler *Compiler) CompilerPass {
	return NewConvertChunksToESModules(compiler)
}

func TestRound94Slice3_ConvertChunksToESModulesMergedScriptFeatures(t *testing.T) {
	RunJUnit4(t, func() JUnitTestCase { return newRound94ConvertChunksCase() })
}

func (c *round94ConvertChunksCase) TestMergedInputsAccumulateFeaturesOnFirstInputScript() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)

	original := JSChunkGraphBuilder_ForStar().
		AddChunk("var a = 1;").
		AddChunk("c;").
		Build()
	original[0].AddSourceFile(SourceFile_FromCode("m0-1", "var b = a ?? 2;"))
	original[0].AddSourceFile(SourceFile_FromCode("m0-2", "var c = b?.x ** 2;"))

	expected := JSChunkGraphBuilder_ForStar().
		AddChunk("var a = 1; var b = a ?? 2; var c = b?.x ** 2; export {c};").
		AddChunk("import {c} from './m0.js'; c;").
		Build()
	expected[0].AddSourceFile(SourceFile_FromCode("m0-1", ""))
	expected[0].AddSourceFile(SourceFile_FromCode("m0-2", ""))

	c.Test(c.Srcs(original), c.Expected(expected))

	lastCompiler := c.GetLastCompiler()
	firstScript := lastCompiler.GetChunkGraph().GetAllChunks()[0].GetFirst().GetAstRoot(lastCompiler)
	fs := NodeUtil_GetFeatureSetOfScript(firstScript)
	for _, want := range []parser.FeatureSet_Feature{
		parser.FeatureSet_Feature_MODULES,
		parser.FeatureSet_Feature_NULL_COALESCE_OP,
		parser.FeatureSet_Feature_OPTIONAL_CHAINING,
		parser.FeatureSet_Feature_EXPONENT_OP,
	} {
		if !fs.ContainsFeature(want) {
			c.T().Fatalf("expected merged chunk first script FeatureSet %v to contain %v", fs, want)
		}
	}
}

func (c *round94ConvertChunksCase) TestZeroArgDynamicImportCallbackReportsDiagnostic() {
	c.IgnoreWarnings(ModuleLoader_LOAD_WARNING)
	c.TestError(
		c.Srcs(
			JSChunkGraphBuilder_ForStar().
				AddChunk("var a = 1; import('./m1.js').then("+ConvertChunksToESModules_DYNAMIC_IMPORT_CALLBACK_FN+"());").
				AddChunk("var b = 2;").
				Build()),
		ConvertChunksToESModules_UNRECOGNIZED_DYNAMIC_IMPORT_CALLBACK,
	)

	compiler := c.GetLastCompiler()
	emptyCall := rhino.NewNode(rhino.Token_CALL)
	if got := ConvertChunksToESModules_GetDynamicImportCallbackModuleNamespace(compiler, emptyCall); !got.IsNil() {
		c.T().Fatalf("expected NilNode for malformed 0-child CALL, got %v", got)
	}
}

func TestRound94Slice3_ReplaceCssNamesMalformedCalls(t *testing.T) {
	if isGetCssNameCall(rhino.NilNode) {
		t.Fatal("expected NilNode not to match goog.getCssName call")
	}
	emptyCall := rhino.NewNode(rhino.Token_CALL)
	if isGetCssNameCall(emptyCall) {
		t.Fatal("expected 0-child CALL not to match goog.getCssName call")
	}
	if err := validateGetCssNameCall(emptyCall); err == nil || err.Type() != ReplaceCssNames_INVALID_NUM_ARGUMENTS_ERROR {
		t.Fatalf("expected INVALID_NUM_ARGUMENTS_ERROR for 0-child CALL, got %v", err)
	}

	calleeOnly := rhino.IR_Call(rhino.IR_Getprop(rhino.IR_Name("goog"), "getCssName"))
	if !isGetCssNameCall(calleeOnly) {
		t.Fatal("expected 0-arg goog.getCssName() to match isGetCssNameCall")
	}
	if err := validateGetCssNameCall(calleeOnly); err == nil || err.Type() != ReplaceCssNames_INVALID_NUM_ARGUMENTS_ERROR {
		t.Fatalf("expected INVALID_NUM_ARGUMENTS_ERROR for 0-arg goog.getCssName(), got %v", err)
	}

	threeArg := rhino.IR_Call(
		rhino.IR_Getprop(rhino.IR_Name("goog"), "getCssName"),
		rhino.IR_Name("base"),
		rhino.IR_String("a"),
		rhino.IR_String("b"),
	)
	if err := validateGetCssNameCall(threeArg); err == nil || err.Type() != ReplaceCssNames_INVALID_NUM_ARGUMENTS_ERROR {
		t.Fatalf("expected INVALID_NUM_ARGUMENTS_ERROR for 3-arg goog.getCssName(), got %v", err)
	}
}

func TestRound94Slice3_ReplaceMessagesMalformedAndProtectedCalls(t *testing.T) {
	if ReplaceMessagesConstants_IsProtectedMessage(rhino.NilNode) {
		t.Fatal("expected NilNode not to be a protected message")
	}
	emptyCall := rhino.NewNode(rhino.Token_CALL)
	if ReplaceMessagesConstants_IsProtectedMessage(emptyCall) {
		t.Fatal("expected 0-child CALL not to be a protected message")
	}
	if _, ok := ReplaceMessages_ProtectedMsgFallback_FromAstNode(rhino.NilNode); ok {
		t.Fatal("expected NilNode not to parse as ProtectedMsgFallback")
	}
	if _, ok := ReplaceMessages_ProtectedMsgFallback_FromAstNode(emptyCall); ok {
		t.Fatal("expected 0-child CALL not to parse as ProtectedMsgFallback")
	}
	if _, ok := ReplaceMessages_ProtectedJsMessage_FromAstNode(rhino.NilNode, nil); ok {
		t.Fatal("expected NilNode not to parse as ProtectedJsMessage")
	}
	if _, ok := ReplaceMessages_ProtectedJsMessage_FromAstNode(emptyCall, nil); ok {
		t.Fatal("expected 0-child CALL not to parse as ProtectedJsMessage")
	}

	// 0-arg __jscomp_define_msg__() throws NullPointerException just like Java checkNotNull(calleeNode.getNext()).
	defineNoArgs := rhino.IR_Call(rhino.IR_Name(ReplaceMessagesConstants_DEFINE_MSG_CALLEE))
	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected panic on 0-arg __jscomp_define_msg__()")
			}
			if _, ok := r.(*jutil.NullPointerException); !ok {
				t.Fatalf("expected *jutil.NullPointerException, got %T (%v)", r, r)
			}
		}()
		ReplaceMessages_ProtectedJsMessage_FromAstNode(defineNoArgs, nil)
	}()

	// Malformed 0-arg and 2-arg __jscomp_msg_fallback__(...) throw IllegalStateException (never NilNode SIGSEGV).
	for _, badFallback := range []rhino.Node{
		rhino.IR_Call(rhino.IR_Name(ReplaceMessagesConstants_FALLBACK_MSG_CALLEE)),
		rhino.IR_Call(
			rhino.IR_Name(ReplaceMessagesConstants_FALLBACK_MSG_CALLEE),
			rhino.IR_String("MSG_A"),
			rhino.IR_Name("a"),
		),
	} {
		func() {
			defer func() {
				r := recover()
				if r == nil {
					t.Fatal("expected panic on malformed __jscomp_msg_fallback__()")
				}
				if _, ok := r.(*jutil.IllegalStateException); !ok {
					t.Fatalf("expected *jutil.IllegalStateException, got %T (%v)", r, r)
				}
			}()
			ReplaceMessages_ProtectedMsgFallback_FromAstNode(badFallback)
		}()
	}

	// Valid 4-arg __jscomp_msg_fallback__ parses cleanly.
	valA := rhino.IR_String("hello")
	valB := rhino.IR_String("bonjour")
	goodFallback := rhino.IR_Call(
		rhino.IR_Name(ReplaceMessagesConstants_FALLBACK_MSG_CALLEE),
		rhino.IR_String("MSG_NEW"),
		valA,
		rhino.IR_String("MSG_OLD"),
		valB,
	)
	fb, ok := ReplaceMessages_ProtectedMsgFallback_FromAstNode(goodFallback)
	if !ok || fb.firstMsgKey != "MSG_NEW" || fb.secondMsgKey != "MSG_OLD" || fb.firstMsgValue != valA || fb.secondMsgValue != valB {
		t.Fatalf("unexpected ProtectedMsgFallback result: ok=%v fb=%+v", ok, fb)
	}
}

type round94StripCodeCase struct {
	CompilerTestCaseBase
}

func newRound94StripCodeCase() *round94StripCodeCase {
	c := &round94StripCodeCase{}
	c.InitCompilerTestCaseBaseWithExterns(c, "function keep(x) {}\nvar obj = {};\n")
	return c
}

func (c *round94StripCodeCase) SetUp() {
	c.CompilerTestCaseBase.SetUp()
	c.EnableNormalize()
}

func (c *round94StripCodeCase) GetProcessor(compiler *Compiler) CompilerPass {
	stripTypes := NewOrderedSet("goog.debug.Logger")
	stripSuffixes := NewOrderedSet("logger_")
	stripPrefixes := NewOrderedSet[string]()
	return NewStripCode(compiler, stripTypes, stripSuffixes, stripPrefixes, false)
}

func TestRound94Slice3_StripCodeMalformedAndNestedCalls(t *testing.T) {
	compiler := NewCompiler()
	compiler.InitCompilerOptionsIfTesting()
	stripTypes := NewOrderedSet("goog.debug.Logger")
	stripSuffixes := NewOrderedSet("logger_")
	strip := NewStripCode(compiler, stripTypes, stripSuffixes, nil, false)
	traversal := &NodeTraversal{compiler: compiler}

	// Malformed / NilNode inputs must not panic.
	if strip.IsCallWhoseReturnValueShouldBeStripped(rhino.NilNode) {
		t.Fatal("expected false for NilNode")
	}
	emptyCall := rhino.NewNode(rhino.Token_CALL)
	if strip.IsCallWhoseReturnValueShouldBeStripped(emptyCall) {
		t.Fatal("expected false for 0-child CALL")
	}
	emptyNew := rhino.NewNode(rhino.Token_NEW)
	if strip.IsCallWhoseReturnValueShouldBeStripped(emptyNew) {
		t.Fatal("expected false for 0-child NEW")
	}
	if strip.IsReferenceToRemovedVar(rhino.NilNode) || strip.IsReferenceToRemovedVar(emptyCall) {
		t.Fatal("expected false for non-NAME nodes in IsReferenceToRemovedVar")
	}
	if strip.IsMethodOrCtorCallThatTriggersRemoval(traversal, rhino.NilNode, rhino.NilNode) {
		t.Fatal("expected false for NilNode in IsMethodOrCtorCallThatTriggersRemoval")
	}
	if strip.IsMethodOrCtorCallThatTriggersRemoval(traversal, emptyCall, rhino.NilNode) {
		t.Fatal("expected false for 0-child CALL in IsMethodOrCtorCallThatTriggersRemoval")
	}

	// Malformed 0-child ASSIGN / INC / EXPR_RESULT attached to a block must not panic.
	block := rhino.IR_Block()
	emptyAssign := rhino.NewNode(rhino.Token_ASSIGN)
	emptyInc := rhino.NewNode(rhino.Token_INC)
	emptyExprStmt := rhino.NewNode(rhino.Token_EXPR_RESULT)
	block.AddChildToBack(emptyAssign)
	block.AddChildToBack(emptyInc)
	block.AddChildToBack(emptyExprStmt)
	strip.MaybeEliminateAssignmentByLvalueName(traversal, emptyAssign, block)
	strip.MaybeEliminateAssignmentByLvalueName(traversal, emptyInc, block)
	strip.MaybeEliminateExpressionByName(emptyExprStmt)
	strip.ReplaceHighestNestedCallWithNull(traversal, emptyCall, rhino.NilNode)

	RunJUnit4(t, func() JUnitTestCase { return newRound94StripCodeCase() })
}

func (c *round94StripCodeCase) TestNestedCallAndCollapsedDeclarationStripping() {
	c.Test(
		"var goog = {debug: {Logger: {}, logger_: 1}};\n"+
			"var my$logger_;\n"+
			"var log = goog.debug.Logger.getLogger('a');\n"+
			"log.info('x').chain();\n"+
			"var y; y = log.info('x');\n"+
			"log.sub.prop = keep(1);\n"+
			"log.arr[0]++;\n"+
			"keep(obj[log.key()], log.arg());\n",
		"var goog = {debug: {Logger: {}}};\n"+
			"var my$logger_ = null;\n"+
			"var y; y = null;\n"+
			"keep(1);\n"+
			"keep(obj[null], null);\n",
	)
}
