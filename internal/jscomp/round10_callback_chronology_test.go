package jscomp

import (
	"slices"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type round10PropertyTestConvention struct {
	CodingConvention
	result bool
	calls  int
}

func (c *round10PropertyTestConvention) IsPropertyTestFunction(n *rhino.Node) bool {
	c.calls++
	return c.result
}

type round10ConventionSwitchCallback struct {
	GuardedCallbackBase[string]
	compiler *Compiler
	next     CodingConvention
	guarded  map[string]bool
}

func (c *round10ConventionSwitchCallback) VisitGuarded(_ *NodeTraversal, n *rhino.Node, _ *rhino.Node) {
	if !n.IsName() {
		return
	}
	name := n.GetString()
	if name == "switchConvention" {
		c.compiler.GetOptions().SetCodingConvention(c.next)
	}
	if name == "before" || name == "after" {
		c.guarded[name] = c.IsGuarded(name)
	}
}

func TestGuardedCallbackUsesCurrentCodingConvention(t *testing.T) {
	for _, mode := range []string{"before_traversal", "between_sibling_calls"} {
		t.Run(mode, func(t *testing.T) {
			compiler := NewCompiler()
			compiler.InitOptions(NewCompilerOptions())
			first := &round10PropertyTestConvention{CodingConvention: compiler.GetCodingConvention()}
			second := &round10PropertyTestConvention{CodingConvention: compiler.GetCodingConvention(), result: true}
			compiler.GetOptions().SetCodingConvention(first)
			cb := &round10ConventionSwitchCallback{compiler: compiler, next: second, guarded: make(map[string]bool)}
			cb.InitGuardedCallbackBase(cb, compiler)
			var script *rhino.Node
			if mode == "before_traversal" {
				compiler.GetOptions().SetCodingConvention(second)
				script = rhino.IR_ScriptOfNodes(rhino.IR_ExprResult(rhino.IR_Call(rhino.IR_Name("probe"), rhino.IR_Name("after"))))
			} else {
				script = rhino.IR_ScriptOfNodes(
					rhino.IR_ExprResult(rhino.IR_Call(rhino.IR_Name("probe"), rhino.IR_Name("before"))),
					rhino.IR_ExprResult(rhino.IR_Name("switchConvention")),
					rhino.IR_ExprResult(rhino.IR_Call(rhino.IR_Name("probe"), rhino.IR_Name("after"))))
			}
			NodeTraversal_Traverse(compiler, rhino.IR_Root(script), cb)
			if guarded, seen := cb.guarded["after"]; !seen || !guarded {
				t.Fatal("property-test argument ignored the current coding convention")
			}
			wantFirstCalls := 0
			if mode == "between_sibling_calls" {
				wantFirstCalls = 1
				if guarded, seen := cb.guarded["before"]; !seen || guarded {
					t.Fatal("first property-test argument used the later convention")
				}
			}
			if first.calls != wantFirstCalls || second.calls != 1 {
				t.Fatalf("property-test callback counts = %d/%d, want %d/1", first.calls, second.calls, wantFirstCalls)
			}
		})
	}
}

type round10ExportTraceConvention struct {
	CodingConvention
	trace      []string
	panicValue any
}

func (c *round10ExportTraceConvention) IsExportedWithLocal(name string, local bool) bool {
	if !local {
		panic("FSIV must classify local names as local")
	}
	c.trace = append(c.trace, name)
	if c.panicValue != nil {
		panic(c.panicValue)
	}
	// The BLOCK definition is deliberately rejected by the convention before
	// candidate analysis, as it was on the validated pre-optimization path.
	return name == "block"
}

func TestFSIVExportConventionBeforeDefinitionFiltering(t *testing.T) {
	for _, mode := range []string{"missing_function_block", "panic_before_missing_definition"} {
		t.Run(mode, func(t *testing.T) {
			compiler := NewCompiler()
			compiler.InitOptions(NewCompilerOptions())
			convention := &round10ExportTraceConvention{CodingConvention: compiler.GetCodingConvention()}
			compiler.GetOptions().SetCodingConvention(convention)
			pass := NewFlowSensitiveInlineVariables(compiler)
			pass.reachingDef = &MustBeReachingVariableDef{allVarsInFn: jutil.NewLinkedHashMap[string, *Var]()}
			state := NewMustBeReachingVariableDef_MustDef()
			global := Scope_CreateGlobalScope(rhino.IR_Root())
			for i, name := range []string{"missing", "fn", "block"} {
				v := NewVar(name, rhino.IR_Name(name), global, i, nil, nil)
				pass.reachingDef.allVarsInFn.Put(name, v)
				if name == "fn" {
					fn := rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(), rhino.IR_Block())
					state.reachingDef.Put(v, NewMustBeReachingVariableDef_Definition(fn))
				} else if name == "block" {
					state.reachingDef.Put(v, NewMustBeReachingVariableDef_Definition(rhino.IR_Block()))
				}
			}
			cb := NewFlowSensitiveInlineVariables_GatherCandidatesCfgNodeCallback(pass)
			cb.inState = state
			visit := func(name string) {
				n := rhino.IR_Name(name)
				parent := rhino.IR_ExprResult(n)
				cb.Visit(nil, n, parent)
			}
			if mode == "panic_before_missing_definition" {
				sentinel := &struct{ label string }{"custom convention"}
				convention.panicValue = sentinel
				var got any
				func() {
					defer func() { got = recover() }()
					visit("missing")
				}()
				if got != sentinel || !slices.Equal(convention.trace, []string{"missing"}) {
					t.Fatalf("custom convention panic suppressed or replaced: panic=%v trace=%v", got, convention.trace)
				}
				return
			}
			visit("missing")
			visit("fn")
			visit("block")
			if !slices.Equal(convention.trace, []string{"missing", "fn", "block"}) {
				t.Fatalf("export callback trace = %v, want missing/fn/block before definition filtering", convention.trace)
			}
		})
	}
}

type round10CoalesceExportConvention struct {
	CodingConvention
	trace      []string
	panicValue *jutil.RuntimeException
}

func (c *round10CoalesceExportConvention) IsExported(name string) bool {
	c.trace = append(c.trace, name)
	if c.panicValue != nil {
		panic(c.panicValue)
	}
	return false
}

func TestCoalesceSingleLocalPreservesCustomExportCallback(t *testing.T) {
	for _, mode := range []string{"return", "panic"} {
		t.Run(mode, func(t *testing.T) {
			fixture := newSyntacticScopeCreatorTest(t)
			script := fixture.getRoot("function f() { var only; return only; }")
			externs := rhino.IR_Root()
			root := rhino.IR_Root(script)
			_ = rhino.IR_Root(externs, root)
			fixture.compiler.SetLifeCycleStage(AbstractCompiler_LifeCycleStage_NORMALIZED)
			fixture.compiler.SetAccessorSummary(newAccessorSummaryAssuming(false))
			convention := &round10CoalesceExportConvention{CodingConvention: fixture.compiler.GetCodingConvention()}
			if mode == "panic" {
				convention.panicValue = jutil.NewRuntimeException("round10 custom export callback")
			}
			fixture.compiler.GetOptions().SetCodingConvention(convention)
			pass := NewCoalesceVariableNames(fixture.compiler, false)
			var caught any
			func() {
				defer func() { caught = recover() }()
				pass.Process(externs, root)
			}()
			if !slices.Equal(convention.trace, []string{"only"}) {
				t.Fatalf("single-local export callback trace = %v, want [only]", convention.trace)
			}
			if mode == "return" {
				if caught != nil {
					t.Fatalf("single-local analysis unexpectedly failed: %v", caught)
				}
			} else if caught != convention.panicValue {
				wrapped, ok := caught.(jutil.Throwable)
				if !ok || wrapped.GetCause() != convention.panicValue {
					t.Fatalf("custom export callback exception was suppressed or replaced: %v", caught)
				}
			}
		})
	}
}
