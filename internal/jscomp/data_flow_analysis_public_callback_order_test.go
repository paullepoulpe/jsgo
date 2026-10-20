package jscomp

import (
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type dfaPublicOrderConvention struct {
	CodingConvention
	events *[]string
}

func (c *dfaPublicOrderConvention) IsExported(name string) bool {
	*c.events = append(*c.events, "export:"+name)
	return name == "exported"
}

type dfaPublicOrderCreator struct {
	delegate ScopeCreator
	events *[]string
	onCreate func()
}

func (c *dfaPublicOrderCreator) CreateScope(n *rhino.Node, parent AbstractScope) AbstractScope {
	*c.events = append(*c.events, "scope:"+n.GetToken().String())
	if c.onCreate != nil {
		run := c.onCreate
		c.onCreate = nil
		run()
	}
	return c.delegate.CreateScope(n, parent)
}

func dfaPublicOrderFixture(t *testing.T) (*Compiler, *Scope, *SyntacticScopeCreator, *jutil.LinkedHashMap[string, *Var]) {
	t.Helper()
	fixture := newSyntacticScopeCreatorTest(t)
	script := fixture.getRoot("function outer() { var local; var exported; function inner() { return local; } }")
	fn := script.GetFirstChild()
	if fn == nil || !fn.IsFunction() {
		t.Fatal("fixture did not parse an attached function declaration")
	}
	creator := NewSyntacticScopeCreator(fixture.compiler)
	global := creator.CreateScopeOfScope(script, nil)
	function := creator.CreateScopeOfScope(fn, global)
	all := NodeUtil_GetAllVarsDeclaredInFunction(fixture.compiler, creator, function).GetAllVariables()
	if all.Get("local") == nil || all.Get("exported") == nil {
		t.Fatal("fixture did not collect the outer bindings")
	}
	return fixture.compiler, function, creator, all
}

func TestDataFlowComputeEscapedPreservesPublicCallbackOrder(t *testing.T) {
	compiler, scope, delegate, all := dfaPublicOrderFixture(t)
	var events []string
	compiler.GetOptions().SetCodingConvention(&dfaPublicOrderConvention{CodingConvention: NewGoogleCodingConvention(), events: &events})
	creator := &dfaPublicOrderCreator{delegate: delegate, events: &events}
	escaped := jutil.NewObjectLinkedHashSet[*Var]()
	DataFlowAnalysis_ComputeEscaped(scope, escaped, compiler, creator, all)
	seenScope, seenExport := false, false
	for _, event := range events {
		if strings.HasPrefix(event, "scope:") {
			if seenExport {
				t.Fatalf("scope callback moved after export classification: %v", events)
			}
			seenScope = true
		} else if strings.HasPrefix(event, "export:") {
			if !seenScope {
				t.Fatalf("export classification moved before scope callbacks: %v", events)
			}
			seenExport = true
		}
	}
	if !seenScope || !seenExport {
		t.Fatalf("fixture did not exercise both public callback phases: %v", events)
	}
	var names []string
	for variable := range escaped.All() {
		names = append(names, variable.GetName())
	}
	if !slices.Equal(names, []string{"local", "exported"}) {
		t.Fatalf("escaped insertion order = %v, want [local exported]", names)
	}
}

func TestDataFlowComputeEscapedObservesCreatorMutationBeforeExportClassification(t *testing.T) {
	compiler, scope, delegate, all := dfaPublicOrderFixture(t)
	var events []string
	compiler.GetOptions().SetCodingConvention(&dfaPublicOrderConvention{CodingConvention: NewGoogleCodingConvention(), events: &events})
	creator := &dfaPublicOrderCreator{delegate: delegate, events: &events, onCreate: func() {
		all.Remove("exported")
	}}
	escaped := jutil.NewObjectLinkedHashSet[*Var]()
	DataFlowAnalysis_ComputeEscaped(scope, escaped, compiler, creator, all)
	var names []string
	for variable := range escaped.All() {
		names = append(names, variable.GetName())
	}
	if all.Get("exported") != nil || !slices.Equal(names, []string{"local"}) || slices.Contains(events, "export:exported") {
		t.Fatalf("creator mutation was observed too late: escaped=%v events=%v", names, events)
	}
}
