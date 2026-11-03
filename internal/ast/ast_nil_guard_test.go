package ast

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func TestIsEs6ConstructorDetachedAndNilGuards(t *testing.T) {
	store := rhino.NewNodeStore()
	prev := rhino.SetCurrentNodeStore(store)
	defer rhino.SetCurrentNodeStore(prev)

	// Zero-value AST wrappers (wrapping NilNode) must return false without panicking.
	var zeroMember MemberFunctionDef
	var zeroFn Function
	if zeroMember.IsEs6Constructor() {
		t.Fatal("zero MemberFunctionDef.IsEs6Constructor() = true, want false")
	}
	if zeroFn.IsEs6Constructor() {
		t.Fatal("zero Function.IsEs6Constructor() = true, want false")
	}

	// Detached MemberFunctionDef named "constructor" (parent is NilNode) must return false.
	rawFn := rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(), rhino.IR_Block())
	rawMember := rhino.IR_MemberFunctionDef("constructor", rawFn)
	detachedMember := MustMemberFunctionDef(rawMember)
	detachedFn := MustFunction(rawFn)
	if detachedMember.IsEs6Constructor() {
		t.Fatal("detached MemberFunctionDef.IsEs6Constructor() = true, want false")
	}
	if detachedFn.IsEs6Constructor() {
		t.Fatal("Function inside detached MemberFunctionDef: IsEs6Constructor() = true, want false")
	}

	// Attached to an OBJECTLIT rather than CLASS_MEMBERS must return false.
	objLit := rhino.IR_Objectlit(rawMember)
	if detachedMember.IsEs6Constructor() || detachedFn.IsEs6Constructor() {
		t.Fatal("constructor() inside OBJECTLIT should not be an ES6 class constructor")
	}
	rawMember.Detach()
	_ = objLit

	// Attached to CLASS_MEMBERS: non-static "constructor" is true; static "constructor" is false.
	classMembers := rhino.NewNode(rhino.Token_CLASS_MEMBERS)
	classMembers.AddChildToBack(rawMember)
	if !detachedMember.IsEs6Constructor() {
		t.Fatal("non-static constructor() inside CLASS_MEMBERS: MemberFunctionDef.IsEs6Constructor() = false, want true")
	}
	if !detachedFn.IsEs6Constructor() {
		t.Fatal("non-static constructor() inside CLASS_MEMBERS: Function.IsEs6Constructor() = false, want true")
	}

	detachedMember.SetStatic(true)
	if detachedMember.IsEs6Constructor() || detachedFn.IsEs6Constructor() {
		t.Fatal("static constructor() inside CLASS_MEMBERS must return false")
	}
}

func TestRhinoNilNodeMatchingAndEquivalenceGuards(t *testing.T) {
	store := rhino.NewNodeStore()
	prev := rhino.SetCurrentNodeStore(store)
	defer rhino.SetCurrentNodeStore(prev)

	nameFoo := rhino.IR_Name("foo")
	qnameFooBar := rhino.IR_Getprop(rhino.IR_Name("foo"), "bar")
	malformedGetprop := rhino.IR_Getprop(rhino.IR_Name("tmp"), "bar")
	malformedGetprop.RemoveChildren()

	if rhino.NilNode.MatchesName("foo") || rhino.NilNode.MatchesName("") {
		t.Fatal("NilNode.MatchesName must return false")
	}
	if rhino.NilNode.MatchesNameNode(rhino.NilNode) ||
		rhino.NilNode.MatchesNameNode(nameFoo) ||
		nameFoo.MatchesNameNode(rhino.NilNode) {
		t.Fatal("MatchesNameNode with NilNode must return false")
	}

	if rhino.NilNode.MatchesQualifiedName("foo") ||
		rhino.NilNode.MatchesQualifiedName("foo.bar") ||
		malformedGetprop.MatchesQualifiedName("foo.bar") {
		t.Fatal("MatchesQualifiedName on NilNode or childless GETPROP must return false")
	}
	if rhino.NilNode.MatchesQualifiedNameWithEnd("foo.bar", 7) ||
		malformedGetprop.MatchesQualifiedNameWithEnd("foo.bar", 7) {
		t.Fatal("MatchesQualifiedNameWithEnd on NilNode or childless GETPROP must return false")
	}
	if rhino.NilNode.MatchesQualifiedNameNode(rhino.NilNode) ||
		rhino.NilNode.MatchesQualifiedNameNode(qnameFooBar) ||
		qnameFooBar.MatchesQualifiedNameNode(rhino.NilNode) ||
		malformedGetprop.MatchesQualifiedNameNode(qnameFooBar) ||
		qnameFooBar.MatchesQualifiedNameNode(malformedGetprop) {
		t.Fatal("MatchesQualifiedNameNode with NilNode or childless GETPROP must return false")
	}

	if rhino.NilNode.IsQualifiedName() || malformedGetprop.IsQualifiedName() {
		t.Fatal("IsQualifiedName on NilNode or childless GETPROP must return false")
	}
	if rhino.NilNode.IsUnscopedQualifiedName() || malformedGetprop.IsUnscopedQualifiedName() {
		t.Fatal("IsUnscopedQualifiedName on NilNode or childless GETPROP must return false")
	}
	if s, ok := rhino.NilNode.QualifiedName(); ok || s != "" {
		t.Fatalf("NilNode.QualifiedName() = (%q, %v), want (\"\", false)", s, ok)
	}
	if rhino.NilNode.IsValidAssignmentTarget() {
		t.Fatal("NilNode.IsValidAssignmentTarget() = true, want false")
	}
	if rhino.NilNode.IsEquivalentTo(rhino.NilNode) ||
		rhino.NilNode.IsEquivalentTo(nameFoo) ||
		nameFoo.IsEquivalentTo(rhino.NilNode) {
		t.Fatal("IsEquivalentTo with NilNode must return false")
	}

	// QualifiedName.Matches(NilNode) across simple, dotted, GetpropQname, and NodeQname implementations.
	qnames := []rhino.QualifiedName{
		rhino.QualifiedName_Of("foo"),
		rhino.QualifiedName_Of("foo.bar.baz"),
		rhino.QualifiedName_Of("foo").Getprop("bar"),
		nameFoo.GetQualifiedNameObject(),
		qnameFooBar.GetQualifiedNameObject(),
	}
	for i, q := range qnames {
		if q.Matches(rhino.NilNode) {
			t.Fatalf("qnames[%d] (%s).Matches(NilNode) = true, want false", i, q.Join())
		}
		if q.Matches(malformedGetprop) {
			t.Fatalf("qnames[%d] (%s).Matches(malformedGetprop) = true, want false", i, q.Join())
		}
	}
}
