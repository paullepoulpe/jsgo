package jscomp

// Shared helpers of the ports of NodeTraversalTest, SyntacticScopeCreatorTest,
// MemoizedScopeCreatorTest and ReferenceCollectorTest.

import (
	"slices"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// scopeTestParse is the tests' parse(Compiler, String):
// compiler.parseTestCode(js), checked to report no errors, put under a
// ROOT.
func scopeTestParse(t *testing.T, compiler *testCompiler, js string) *rhino.Node {
	t.Helper()
	n := compiler.ParseTestCode(js)
	if errors := compiler.GetErrors(); len(errors) != 0 {
		t.Fatalf("parse errors: %v", errors)
	}
	rhino.IR_Root(n)
	return n
}

// scopeTestParseRoots is NodeTraversalTest.parseRoots.
func scopeTestParseRoots(t *testing.T, compiler *testCompiler, externs, js string) *rhino.Node {
	t.Helper()
	extern := scopeTestParse(t, compiler, externs).Detach()
	main := scopeTestParse(t, compiler, js).Detach()
	return rhino.IR_Root(rhino.IR_Root(extern), rhino.IR_Root(main))
}

// scopeTestNewCompiler is `new Compiler()` followed by
// `initOptions(options)`.
func scopeTestNewCompiler(options *CompilerOptions) *testCompiler {
	return newTestCompiler(options).(*testCompiler)
}

// scopeTestCatch runs body and returns what it throws (nil if nothing).
func scopeTestCatch(body func()) (thrown jutil.Throwable) {
	jutil.Try(body, func(e jutil.Throwable) { thrown = e })
	return thrown
}

// scopeTestVarNames lists the names of vars.
func scopeTestVarNames[V interface{ GetName() string }](vars func(func(V) bool)) []string {
	var names []string
	for v := range vars {
		names = append(names, v.GetName())
	}
	return names
}

// scopeSubject is the part of com.google.javascript.jscomp.testing.ScopeSubject
// these tests use: assertScope(scope).declares(name)... and
// .doesNotDeclare(name), for untyped scopes.
type scopeSubject struct {
	t      *testing.T
	actual *Scope
}

func assertScope(t *testing.T, scope *Scope) *scopeSubject {
	return &scopeSubject{t: t, actual: scope}
}

// getVar is ScopeSubject.getVar.
func (s *scopeSubject) getVar(name string) *Var {
	if s.actual.HasSlot(name) {
		return jutil.CheckNotNull(s.actual.GetVar(name))
	}
	return nil
}

// Java: ScopeSubject.doesNotDeclare(String)
func (s *scopeSubject) doesNotDeclare(name string) {
	s.t.Helper()
	if v := s.getVar(name); v != nil {
		s.t.Errorf("expected not to declare %s, but declared it with value %v (scope was %v)", name, v, s.actual)
	}
}

// Java: ScopeSubject.declares(String)
func (s *scopeSubject) declares(name string) *scopeDeclarationSubject {
	s.t.Helper()
	v := s.getVar(name)
	if v == nil {
		names := scopeTestVarNames(s.actual.GetAllAccessibleVariables())
		if len(names) > 10 {
			names = append(slices.Clone(names[:9]), "and "+jutil.ToJavaString(len(names)-9)+" others")
		}
		s.t.Errorf("expected to declare %s, but did not; did declare %s (scope was %v)",
			name, strings.Join(names, ", "), s.actual)
		return &scopeDeclarationSubject{s: s}
	}
	if got := s.actual.GetTopmostScopeOfEventualDeclaration(name); got != v.GetScope() {
		s.t.Errorf("%s: getTopmostScopeOfEventualDeclaration is %v, the var is declared on %v", name, got, v.GetScope())
	}
	return &scopeDeclarationSubject{s: s, var_: v}
}

// scopeDeclarationSubject is ScopeSubject.DeclarationSubject.
type scopeDeclarationSubject struct {
	s    *scopeSubject
	var_ *Var
}

// Java: ScopeSubject.DeclarationSubject.expectScope(String, Object, AbstractScope<?, ?>)
func (d *scopeDeclarationSubject) expectScope(preposition string, expected any, scope *Scope) {
	d.s.t.Helper()
	if d.var_ == nil {
		return // declares already failed
	}
	if d.var_.GetScope() != scope {
		if preposition != "" {
			preposition = " " + preposition
		}
		d.s.t.Errorf("for var %s: expected to be declared%s %v, but is declared on %v",
			d.var_.GetName(), preposition, expected, d.var_.GetScope())
	}
}

// Java: ScopeSubject.DeclarationSubject.directly()
func (d *scopeDeclarationSubject) directly() *scopeDeclarationSubject {
	d.s.t.Helper()
	d.expectScope("", "directly", d.s.actual)
	return d
}

// Java: ScopeSubject.DeclarationSubject.on(AbstractScope<?, ?>)
func (d *scopeDeclarationSubject) on(scope *Scope) *scopeDeclarationSubject {
	d.s.t.Helper()
	jutil.CheckStateMsg1(scope != d.s.actual,
		"It doesn't make sense to pass the scope already being asserted about. Use .directly()%s", "")
	d.expectScope("on", scope, scope)
	return d
}

// Java: ScopeSubject.DeclarationSubject.onClosestContainerScope()
func (d *scopeDeclarationSubject) onClosestContainerScope() *scopeDeclarationSubject {
	d.s.t.Helper()
	d.expectScope("on", "the closest container scope", d.s.actual.GetClosestContainerScope())
	return d
}

// Java: ScopeSubject.DeclarationSubject.onClosestHoistScope()
func (d *scopeDeclarationSubject) onClosestHoistScope() *scopeDeclarationSubject {
	d.s.t.Helper()
	d.expectScope("on", "the closest hoist scope", d.s.actual.GetClosestHoistScope())
	return d
}

// Java: ScopeSubject.DeclarationSubject.globally()
func (d *scopeDeclarationSubject) globally() *scopeDeclarationSubject {
	d.s.t.Helper()
	d.expectScope("", "globally", d.s.actual.GetGlobalScope())
	return d
}

// Java: ScopeSubject.DeclarationSubject.onSomeParent()
func (d *scopeDeclarationSubject) onSomeParent() *scopeDeclarationSubject {
	d.s.t.Helper()
	if d.var_ != nil && d.var_.GetScope() == d.s.actual {
		d.s.t.Errorf("for var %s: expected a declaration on a parent scope, but found it declared directly (scope was %v)",
			d.var_.GetName(), d.s.actual)
	}
	return d
}
