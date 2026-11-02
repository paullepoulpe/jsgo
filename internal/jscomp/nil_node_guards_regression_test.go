package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestNilNodeGuardsRegression(t *testing.T) {
	t.Run("PolymerDescriptorAndPropertiesWithSpread", func(t *testing.T) {
		compiler := NewCompiler()
		opts := NewCompilerOptions()
		pv := int32(2)
		opts.SetPolymerVersion(&pv)
		compiler.Init([]*SourceFile{
			SourceFile_FromCode("externs.js", polymerPassTest_EXTERNS+"\nPolymerElement;\n"),
		}, []*SourceFile{
			SourceFile_FromCode("app.js", `
const shared = {extra() {}};
/** @polymerBehavior */
const MyBehavior = {
  ...shared,
  properties: {
    propA: String,
  },
  ready() {},
};
Polymer({
  is: 'x-foo',
  behaviors: [MyBehavior],
  ...shared,
  properties: {
    count: Number,
  },
  attached() {},
});
`),
		}, opts)
		compiler.ParseInputs()
		NewGatherModuleMetadata(compiler, false, ModuleLoader_ResolutionMode_BROWSER).Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
		pass := NewPolymerPass(compiler)
		pass.Process(compiler.GetExternsRoot(), compiler.GetJsRoot())

		// Also exercise addBehaviorSuppressions and suppressDefaultValues on an
		// OBJECTLIT containing childless STRING_KEY nodes directly.
		propsObj := rhino.IR_Objectlit(rhino.IR_StringKey("childlessProp"))
		behaviorObj := rhino.IR_Objectlit(
			rhino.IR_StringKey("childlessKey"),
			rhino.IR_Propdef(rhino.IR_StringKey("properties"), propsObj),
		)
		(polymerPassSuppressBehaviorsAndProtectKeys{compiler: compiler}).addBehaviorSuppressions(behaviorObj)
	})

	t.Run("PolymerElementPropExprBareChildlessExprStmt", func(t *testing.T) {
		// Construct an EXPR_RESULT wrapping a childless NAME ("PolymerElement;")
		// and verify polymerPassFindExterns_isPolymerElementPropExpr returns false cleanly.
		stmt := rhino.IR_ExprResult(rhino.IR_Name("PolymerElement"))
		if polymerPassFindExterns_isPolymerElementPropExpr(stmt) {
			t.Fatalf("expected false for bare NAME expression statement")
		}
	})

	t.Run("CheckEs6ModuleFileStructureNilParentAtRoot", func(t *testing.T) {
		compiler := NewCompiler()
		compiler.InitCompilerOptionsIfTesting()
		pass := NewCheckEs6ModuleFileStructure(compiler)
		nt := &NodeTraversal{compiler: compiler, callback: pass}

		block := rhino.IR_Block()
		if pass.ShouldTraverse(nt, block, rhino.NilNode) {
			t.Fatalf("expected false for detached BLOCK at root")
		}
		exprStmt := rhino.IR_ExprResult(rhino.IR_Name("x"))
		if !pass.VisitExprResult(nt, ast.MustExprStmt(exprStmt), rhino.NilNode) {
			t.Fatalf("expected true for detached EXPR_RESULT")
		}
		varDecl := rhino.IR_Var(rhino.IR_Name("x"))
		if !pass.VisitDeclaration(nt, ast.MustNameDeclaration(varDecl), rhino.NilNode) {
			t.Fatalf("expected true for detached VAR")
		}
	})

	t.Run("CollapsePropertiesDynamicImportNamespaceStubWithoutRValue", func(t *testing.T) {
		compiler := NewCompiler()
		opts := NewCompilerOptions()
		compiler.Init([]*SourceFile{}, []*SourceFile{
			SourceFile_FromCode("mod.js", `
/** @const */ var module$foo = {};
/** @type {?} */ module$foo.bar;
`),
		}, opts)
		compiler.ParseInputs()
		outer := NewInlineAndCollapseProperties(compiler, InlineAndCollapsePropertiesOptions{
			PropertyCollapseLevel:    CompilerOptions_PropertyCollapseLevel_ALL,
			ChunkOutputType:          CompilerOptions_ChunkOutputType_ES_MODULES,
			HaveModulesBeenRewritten: true,
			ModuleResolutionMode:     ModuleLoader_ResolutionMode_BROWSER,
		})
		outer.namespace = NewGlobalNamespace(compiler, compiler.GetJsRoot())
		cp := NewInlineAndCollapseProperties_CollapseProperties(outer)
		cp.addDynamicallyImportedModule("module$foo")
		cp.Process(compiler.GetExternsRoot(), compiler.GetJsRoot())
	})

	t.Run("AnalyzePrototypePropertiesMaybeGetVarNilNode", func(t *testing.T) {
		compiler := NewCompiler()
		compiler.InitCompilerOptionsIfTesting()
		app := NewAnalyzePrototypeProperties(compiler, nil, false, false, false)
		pp := &AnalyzePrototypeProperties_ProcessProperties{outer: app}
		nt := &NodeTraversal{compiler: compiler, callback: pp}
		if got := pp.MaybeGetVar(nt, rhino.NilNode); got != nil {
			t.Fatalf("MaybeGetVar(NilNode) = %v, want nil", got)
		}
	})

	t.Run("ConformanceRegexWordBoundariesAndPostDollarSentinel", func(t *testing.T) {
		tc := NewCheckConformanceTest()
		tc.SetUp()
		tc.configuration = `requirement: {
  type: BANNED_NAME
  value: 'eval'
  allowlist_regexp: '\\bSRC1\\b'
  allowlist_regexp: 'never$match'
  error_message: 'eval is banned'
}`
		// SRC1 matches \bSRC1\b (allowlisted); SRC2 matches neither \bSRC1\b nor never$match, so 1 violation is reported.
		tc.TestConformance("eval();", "eval();", CheckConformance_CONFORMANCE_VIOLATION)
	})
}
