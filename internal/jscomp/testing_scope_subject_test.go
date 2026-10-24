// Port of com/google/javascript/jscomp/testing/ScopeSubject.java @ closure-compiler v20261021

/*
 * Copyright 2015 The Closure Compiler Authors.
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
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

// ScopeSubject is a Truth Subject for the AbstractScope class. Usage:
//
//	ScopeSubject_AssertScope(scope).Declares("somevar")
//	ScopeSubject_AssertScope(scope).Declares("otherVar").Directly()
//	ScopeSubject_AssertScope(scope).Declares("yetAnotherVar").OnClosestContainerScope()
//
// The scope is an AbstractScope<?, ?> (a *Scope or a *TypedScope); the
// AbstractScope methods whose Java signatures mention S or V are reached
// through scopeSubject_getVar and its siblings below.
type ScopeSubject struct {
	*Subject
	actual AbstractScope
}

// Java: ScopeSubject.assertScope(AbstractScope<?, ?>)
func ScopeSubject_AssertScope(scope AbstractScope) *ScopeSubject {
	// NB: Eclipse's Java compiler bails on just passing ScopeSubject::new below, so wrap it in a
	// Closure.
	return newScopeSubject(&Subject{actual: scope}, scope)
}

// AssertScope is ScopeSubject_AssertScope for Java's
// `import static ...ScopeSubject.assertScope` (Go-only: language-adapter).
func AssertScope(scope AbstractScope) *ScopeSubject {
	return ScopeSubject_AssertScope(scope)
}

// Java: ScopeSubject(FailureMetadata, AbstractScope<?, ?>)
func newScopeSubject(failureMetadata *Subject, scope AbstractScope) *ScopeSubject {
	return &ScopeSubject{Subject: failureMetadata, actual: scope}
}

// Java: ScopeSubject.doesNotDeclare(String)
func (s *ScopeSubject) DoesNotDeclare(name string) {
	var_ := s.getVar(name)
	if var_ != nil {
		s.FailWithoutActual(
			Fact_Fact("expected not to declare", name),
			Fact_Fact("but declared it with value", var_),
			Fact_Fact("scope was", s.actual))
	}
}

// Java: ScopeSubject.declares(String)
func (s *ScopeSubject) Declares(name string) *ScopeSubject_DeclarationSubject {
	var_ := s.getVar(name)
	if var_ == nil {
		declared := scopeSubject_getAllAccessibleVariables(s.actual)
		names := make([]string, 0, len(declared))
		for _, v := range declared {
			names = append(names, v.GetName())
		}
		if len(names) > 10 {
			names = append(append([]string{}, names[0:9]...), "and "+strconv.Itoa(len(names)-9)+" others")
		}
		s.FailWithoutActual(
			Fact_Fact("expected to declare", name),
			Fact_SimpleFact("but did not"),
			Fact_Fact("did declare", strings.Join(names, ", ")),
			Fact_Fact("scope was", s.actual))
	}
	AssertThat(s.actual.(interface {
		GetTopmostScopeOfEventualDeclaration(name string) rhino.StaticScope
	}).GetTopmostScopeOfEventualDeclaration(name)).IsEqualTo(var_.GetScope())
	return newScopeSubject_DeclarationSubject(s, var_)
}

// Java: ScopeSubject.getVar(String)
func (s *ScopeSubject) getVar(name string) AbstractVar {
	if s.actual.HasSlot(name) {
		return jutil.CheckNotNull(scopeSubject_getVar(s.actual, name))
	}
	return nil
}

// ScopeSubject_DeclarationSubject is a subject for an AbstractVar declared
// by this particular scope.
type ScopeSubject_DeclarationSubject struct {
	// Go-only: language-adapter — the enclosing ScopeSubject.
	outer *ScopeSubject
	var_  AbstractVar
}

// Java: ScopeSubject.DeclarationSubject(AbstractVar<?, ?>)
func newScopeSubject_DeclarationSubject(outer *ScopeSubject, var_ AbstractVar) *ScopeSubject_DeclarationSubject {
	return &ScopeSubject_DeclarationSubject{outer: outer, var_: jutil.CheckNotNull(var_)}
}

// expectScope expects the variable to be defined on the given scope. The
// preposition is either "on" or "", depending on whether it is needed for
// grammatical correctness. The expected object is displayed in brackets, in
// parallel to the actual scope that the variable is defined on.
//
// Java: ScopeSubject.DeclarationSubject.expectScope(String, Object, AbstractScope<?, ?>)
func (d *ScopeSubject_DeclarationSubject) expectScope(preposition string, expected any, scope AbstractScope) {
	if d.var_.GetScope() != rhino.StaticScope(scope) {
		space := ""
		if preposition != "" {
			space = " "
		}
		d.outer.FailWithoutActual(
			Fact_Fact("for var", d.var_.GetName()),
			Fact_Fact(
				"expected to be declared"+space+preposition,
				expected),
			Fact_Fact("but is declared on", d.var_.GetScope()))
	}
}

// Directly expects the declared variable to be declared on the subject
// scope.
//
// Java: ScopeSubject.DeclarationSubject.directly()
func (d *ScopeSubject_DeclarationSubject) Directly() *ScopeSubject_DeclarationSubject {
	d.expectScope("", "directly", d.outer.actual)
	return d
}

// On expects the declared variable to be declared on the given scope.
//
// Java: ScopeSubject.DeclarationSubject.on(AbstractScope<?, ?>)
func (d *ScopeSubject_DeclarationSubject) On(scope AbstractScope) *ScopeSubject_DeclarationSubject {
	jutil.CheckStateMsg(
		scope != d.outer.actual,
		"It doesn't make sense to pass the scope already being asserted about. Use .directly()")
	d.expectScope("on", scope, scope)
	return d
}

// OnClosestContainerScope expects the declared variable to be declared on
// the closest container scope.
//
// Java: ScopeSubject.DeclarationSubject.onClosestContainerScope()
func (d *ScopeSubject_DeclarationSubject) OnClosestContainerScope() *ScopeSubject_DeclarationSubject {
	d.expectScope("on", "the closest container scope", scopeSubject_getClosestContainerScope(d.outer.actual))
	return d
}

// OnClosestHoistScope expects the declared variable to be declared on the
// closest hoist scope.
//
// Java: ScopeSubject.DeclarationSubject.onClosestHoistScope()
func (d *ScopeSubject_DeclarationSubject) OnClosestHoistScope() *ScopeSubject_DeclarationSubject {
	d.expectScope("on", "the closest hoist scope", scopeSubject_getClosestHoistScope(d.outer.actual))
	return d
}

// Globally expects the declared variable to be declared on the global
// scope.
//
// Java: ScopeSubject.DeclarationSubject.globally()
func (d *ScopeSubject_DeclarationSubject) Globally() *ScopeSubject_DeclarationSubject {
	d.expectScope("", "globally", scopeSubject_getGlobalScope(d.outer.actual))
	return d
}

// OnSomeParent expects the declared variable to be declared on any scope
// other than the subject.
//
// Java: ScopeSubject.DeclarationSubject.onSomeParent()
func (d *ScopeSubject_DeclarationSubject) OnSomeParent() *ScopeSubject_DeclarationSubject {
	if d.var_ != nil && d.var_.GetScope() == rhino.StaticScope(d.outer.actual) {
		d.outer.FailWithoutActual(
			Fact_Fact("for var", d.var_.GetName()),
			Fact_SimpleFact("expected a declaration on a parent scope"),
			Fact_SimpleFact("but found it declared directly"),
			Fact_Fact("scope was", d.outer.actual))
	}
	return d
}

// OnScopeLabeled expects the declared variable to be declared on some
// scope with the given label.
//
// Java: ScopeSubject.DeclarationSubject.onScopeLabeled(String)
func (d *ScopeSubject_DeclarationSubject) OnScopeLabeled(expectedLabel string) *ScopeSubject_DeclarationSubject {
	actualLabel := d.outer.getLabel(d.var_.GetScopeRoot())
	if actualLabel == nil {
		d.outer.FailWithoutActual(
			Fact_Fact("expected to declare", d.var_.GetName()),
			Fact_Fact("on a scope labeled", expectedLabel),
			Fact_SimpleFact("but declared it on an unlabeled scope"),
			Fact_Fact("scope under test was", d.outer.actual))
	} else if *actualLabel != expectedLabel {
		d.outer.FailWithoutActual(
			Fact_Fact("expected to declare", d.var_.GetName()),
			Fact_Fact("on a scope labeled", expectedLabel),
			Fact_Fact("but declared it on an scope labeled", *actualLabel),
			Fact_Fact("scope under test was", d.outer.actual))
	}
	return d
}

// Java: ScopeSubject.DeclarationSubject.withTypeThat()
func (d *ScopeSubject_DeclarationSubject) WithTypeThat() *rhinotesting.TypeSubject {
	typedVar := d.var_.GetSymbol().(*TypedVar)
	return TypeSubject_AssertType(typedVar.GetType())
}

// getLabel returns the name of the label applied to n, or nil if none
// exists.
//
// Java: ScopeSubject.getLabel(Node)
func (s *ScopeSubject) getLabel(n *rhino.Node) *string {
	// If the node is labeled it will be the second child of a LABEL and the first child
	// will be a LABEL_NAME.
	parent := n.GetParent()
	if parent != nil && parent.IsLabel() {
		labelNameNode := parent.GetFirstChild()
		jutil.CheckStateMsg1(labelNameNode.IsLabelName(), "%s", labelNameNode)
		jutil.CheckStateMsg1(labelNameNode.GetNext() == n, "%s", n)
		label := labelNameNode.GetString()
		return &label
	} else {
		return nil
	}
}

// Go-only: language-adapter — the methods of AbstractScope<?, ?> whose
// Java signatures mention S or V (see the note at the top of
// abstract_scope.go), on the two concrete scopes. They return nil for Java
// null, never a typed nil.

func scopeSubject_getVar(scope AbstractScope, name string) AbstractVar {
	switch s := scope.(type) {
	case *Scope:
		if v := s.GetVar(name); v != nil {
			return v
		}
	case *TypedScope:
		if v := s.GetVar(name); v != nil {
			return v
		}
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg("not a Scope or TypedScope: " + show(scope)))
	}
	return nil
}

func scopeSubject_getAllAccessibleVariables(scope AbstractScope) []AbstractVar {
	var out []AbstractVar
	switch s := scope.(type) {
	case *Scope:
		for v := range s.GetAllAccessibleVariables() {
			out = append(out, v)
		}
	case *TypedScope:
		for v := range s.GetAllAccessibleVariables() {
			out = append(out, v)
		}
	default:
		panic(jutil.NewIllegalArgumentExceptionMsg("not a Scope or TypedScope: " + show(scope)))
	}
	return out
}

func scopeSubject_scopeOf[S interface {
	comparable
	AbstractScope
}](s S) AbstractScope {
	var null S
	if s == null {
		return nil
	}
	return s
}

func scopeSubject_getClosestContainerScope(scope AbstractScope) AbstractScope {
	switch s := scope.(type) {
	case *Scope:
		return scopeSubject_scopeOf(s.GetClosestContainerScope())
	case *TypedScope:
		return scopeSubject_scopeOf(s.GetClosestContainerScope())
	}
	panic(jutil.NewIllegalArgumentExceptionMsg("not a Scope or TypedScope: " + show(scope)))
}

func scopeSubject_getClosestHoistScope(scope AbstractScope) AbstractScope {
	switch s := scope.(type) {
	case *Scope:
		return scopeSubject_scopeOf(s.GetClosestHoistScope())
	case *TypedScope:
		return scopeSubject_scopeOf(s.GetClosestHoistScope())
	}
	panic(jutil.NewIllegalArgumentExceptionMsg("not a Scope or TypedScope: " + show(scope)))
}

func scopeSubject_getGlobalScope(scope AbstractScope) AbstractScope {
	switch s := scope.(type) {
	case *Scope:
		return scopeSubject_scopeOf(s.GetGlobalScope())
	case *TypedScope:
		return scopeSubject_scopeOf(s.GetGlobalScope())
	}
	panic(jutil.NewIllegalArgumentExceptionMsg("not a Scope or TypedScope: " + show(scope)))
}
