package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

func TestControlFlowClassEvalOrder_ForwardAndBackward(t *testing.T) {
	compiler := NewCompiler()
	compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	options := NewCompilerOptions()
	options.SetLanguage(CompilerOptions_LanguageMode_UNSUPPORTED)
	compiler.InitOptions(options)

	src := `function f() {
  var C = class Inner extends sup() {
    instField = deferredInst();
    static s1 = stat1();
    [key1()]() { body1(); }
    get [key2()]() { return body2(); }
    [key3()] = deferredCompInst();
    static [key4()] = stat2();
    static { staticBlockBody(); }
    method() { body3(); }
  };
}`
	script := compiler.ParseTestCode(src)
	if len(compiler.GetErrors()) > 0 {
		t.Fatalf("unexpected parse errors: %v", compiler.GetErrors())
	}
	fn := script.GetFirstChild()
	cls := fn.GetLastChild().GetFirstChild().GetFirstChild().GetFirstChild()

	var forwardCalls []string
	walkCfgCondEvalForward(cls, false, &forwardCalls, func(out *[]string, n rhino.Node, _ bool) (cfgCondEvalState, bool) {
		if n.IsName() {
			*out = append(*out, n.GetString())
		}
		return cfgCondEvalContinue, false
	})

	wantForward := []string{"Inner", "sup", "key1", "key2", "key3", "key4", "stat1", "stat2"}
	if len(forwardCalls) != len(wantForward) {
		t.Fatalf("forward order = %v, want %v", forwardCalls, wantForward)
	}
	for i := range wantForward {
		if forwardCalls[i] != wantForward[i] {
			t.Fatalf("forward order[%d] = %q, want %q (full: %v)", i, forwardCalls[i], wantForward[i], forwardCalls)
		}
	}

	var backwardCalls []string
	walkCfgCondEvalBackward(cls, false, &backwardCalls, func(out *[]string, n rhino.Node, _ bool) bool {
		if n.IsName() {
			*out = append(*out, n.GetString())
		}
		return false
	})

	wantBackward := []string{"stat2", "stat1", "key4", "key3", "key2", "key1", "sup", "Inner"}
	if len(backwardCalls) != len(wantBackward) {
		t.Fatalf("backward order = %v, want %v", backwardCalls, wantBackward)
	}
	for i := range wantBackward {
		if backwardCalls[i] != wantBackward[i] {
			t.Fatalf("backward order[%d] = %q, want %q (full: %v)", i, backwardCalls[i], wantBackward[i], backwardCalls)
		}
	}
}

func TestMaybeReachingVariableUse_IgnoresDeferredNonStaticClassFields(t *testing.T) {
	compiler := NewCompiler()
	compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	options := NewCompilerOptions()
	options.SetLanguage(CompilerOptions_LanguageMode_UNSUPPORTED)
	compiler.InitOptions(options)

	src := `function f() {
  var x = 1;
  var C = class {
    a = use(x);
    [key()] = use(x);
    m() { return x; }
  };
  return C;
}`
	script := compiler.ParseTestCode(src)
	if len(compiler.GetErrors()) > 0 {
		t.Fatalf("unexpected parse errors: %v", compiler.GetErrors())
	}
	fn := script.GetFirstChild()
	cfg := ComputeControlFlowGraph(compiler, fn, false, true)
	xID := uint32(strpool.Intern("x"))

	var ru MaybeReachingVariableUse
	ru.resetWithVarIDs(cfg, []uint32{xID})
	ru.Analyze()

	defStmt := fn.GetLastChild().GetFirstChild() // var x = 1;
	for u := range ru.GetUses("x", defStmt) {
		t.Fatalf("expected no synchronous reaching use from deferred class members, got %v", u)
	}
}

func TestControlFlowClassEvalOrder_ZeroAllocs(t *testing.T) {
	compiler := NewCompiler()
	compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
	options := NewCompilerOptions()
	options.SetLanguage(CompilerOptions_LanguageMode_UNSUPPORTED)
	compiler.InitOptions(options)

	src := `function f() {
  var x = 1;
  var C = class extends sup(x) {
    a = deferred(x);
    static s1 = stat1(x);
    [key1(x)]() {}
    static [key2(x)] = stat2(x);
  };
}`
	script := compiler.ParseTestCode(src)
	if len(compiler.GetErrors()) > 0 {
		t.Fatalf("unexpected parse errors: %v", compiler.GetErrors())
	}
	varStmt := script.GetFirstChild().GetLastChild().GetFirstChild().GetNext()
	cls := varStmt.GetFirstChild().GetFirstChild()
	classMembers := cls.GetLastChild()
	staticField := classMembers.GetFirstChild().GetNext()
	compMethod := staticField.GetNext()
	xID := uint32(strpool.Intern("x"))

	fwdAllocs := testing.AllocsPerRun(100, func() {
		_ = walkCfgCondEvalForward(varStmt, false, xID, daeVisitReadBeforeKill)
	})
	if fwdAllocs != 0 {
		t.Fatalf("walkCfgCondEvalForward allocated %v times per run, want 0", fwdAllocs)
	}

	bwdAllocs := testing.AllocsPerRun(100, func() {
		walkCfgCondEvalBackward(varStmt, false, xID, func(id uint32, n rhino.Node, _ bool) bool {
			return n.IsName() && n.GetStringID() == id
		})
	})
	if bwdAllocs != 0 {
		t.Fatalf("walkCfgCondEvalBackward allocated %v times per run, want 0", bwdAllocs)
	}

	beforeAllocs := testing.AllocsPerRun(100, func() {
		_ = walkClassMembersBefore(classMembers, staticField, true, xID, func(id uint32, expr rhino.Node) bool {
			return expr.IsName() && expr.GetStringID() == id
		})
	})
	if beforeAllocs != 0 {
		t.Fatalf("walkClassMembersBefore allocated %v times per run, want 0", beforeAllocs)
	}

	afterAllocs := testing.AllocsPerRun(100, func() {
		_ = walkClassMembersAfterState(classMembers, compMethod, false, xID, func(id uint32, expr rhino.Node) cfgCondEvalState {
			return walkCfgCondEvalForward(expr, false, id, daeVisitReadBeforeKill)
		})
	})
	if afterAllocs != 0 {
		t.Fatalf("walkClassMembersAfterState allocated %v times per run, want 0", afterAllocs)
	}
}
