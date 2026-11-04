package rhino_test

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

func TestRound96NodeChildrenNilSafety(t *testing.T) {
	prev := rhino.SetCurrentNodeStore(rhino.NewNodeStore())
	defer rhino.SetCurrentNodeStore(prev)

	leaf := rhino.IR_Empty()
	if !leaf.GetFirstFirstChild().IsNil() {
		t.Fatalf("leaf.GetFirstFirstChild() = %v, want NilNode", leaf.GetFirstFirstChild())
	}
	if !leaf.GetSecondChild().IsNil() {
		t.Fatalf("leaf.GetSecondChild() = %v, want NilNode", leaf.GetSecondChild())
	}
	if !leaf.GetPrevious().IsNil() {
		t.Fatalf("detached leaf.GetPrevious() = %v, want NilNode", leaf.GetPrevious())
	}
	for _, idx := range []int{0, 1, 2, 10} {
		if !leaf.GetChildAtIndex(idx).IsNil() {
			t.Fatalf("leaf.GetChildAtIndex(%d) = %v, want NilNode", idx, leaf.GetChildAtIndex(idx))
		}
	}

	// NilNode and detached checks on IsFirstChildOf / IsSecondChildOf / IsOnlyChildOf.
	if rhino.NilNode.IsFirstChildOf(leaf) || rhino.NilNode.IsFirstChildOf(rhino.NilNode) || leaf.IsFirstChildOf(rhino.NilNode) {
		t.Fatal("IsFirstChildOf returned true for NilNode or detached node")
	}
	if rhino.NilNode.IsSecondChildOf(leaf) || rhino.NilNode.IsSecondChildOf(rhino.NilNode) || leaf.IsSecondChildOf(rhino.NilNode) {
		t.Fatal("IsSecondChildOf returned true for NilNode or detached node")
	}
	if rhino.NilNode.IsOnlyChildOf(leaf) || rhino.NilNode.IsOnlyChildOf(rhino.NilNode) || leaf.IsOnlyChildOf(rhino.NilNode) {
		t.Fatal("IsOnlyChildOf returned true for NilNode or detached node")
	}

	// 1-child parent: c0 is first & only child, NOT second child (even though c0.previous == c0 in circular list!).
	p1 := rhino.IR_Block()
	c0 := rhino.IR_Empty()
	p1.AddChildToBack(c0)
	if !p1.GetFirstFirstChild().IsNil() {
		t.Fatal("expected p1.GetFirstFirstChild() to be NilNode when c0 has no children")
	}
	if !p1.GetSecondChild().IsNil() {
		t.Fatal("expected p1.GetSecondChild() to be NilNode")
	}
	if p1.GetChildAtIndex(0) != c0 || !p1.GetChildAtIndex(1).IsNil() || !p1.GetChildAtIndex(5).IsNil() {
		t.Fatal("unexpected GetChildAtIndex results on 1-child parent")
	}
	if !c0.IsFirstChildOf(p1) || !c0.IsOnlyChildOf(p1) || c0.IsSecondChildOf(p1) {
		t.Fatalf("1-child predicates wrong: first=%v only=%v second=%v",
			c0.IsFirstChildOf(p1), c0.IsOnlyChildOf(p1), c0.IsSecondChildOf(p1))
	}

	// 3-child parent.
	c1 := rhino.IR_Empty()
	c2 := rhino.IR_Empty()
	p1.AddChildToBack(c1)
	p1.AddChildToBack(c2)
	if !c0.IsFirstChildOf(p1) || c0.IsOnlyChildOf(p1) || c0.IsSecondChildOf(p1) {
		t.Fatal("c0 predicates wrong on 3-child parent")
	}
	if c1.IsFirstChildOf(p1) || c1.IsOnlyChildOf(p1) || !c1.IsSecondChildOf(p1) {
		t.Fatal("c1 predicates wrong on 3-child parent")
	}
	if c2.IsFirstChildOf(p1) || c2.IsOnlyChildOf(p1) || c2.IsSecondChildOf(p1) {
		t.Fatal("c2 predicates wrong on 3-child parent")
	}
	if p1.GetSecondChild() != c1 || p1.GetChildAtIndex(2) != c2 || !p1.GetChildAtIndex(3).IsNil() {
		t.Fatal("sibling lookup wrong on 3-child parent")
	}
}

func TestRound96MalformedZeroChildGetPropQualifiedName(t *testing.T) {
	prev := rhino.SetCurrentNodeStore(rhino.NewNodeStore())
	defer rhino.SetCurrentNodeStore(prev)

	gp0 := rhino.Node_NewStringWithToken(rhino.Token_GETPROP, "a")
	if q, ok := gp0.GetQualifiedNameString(); ok || q != "" {
		t.Fatalf("gp0.GetQualifiedNameString() = (%q, %v), want (\"\", false)", q, ok)
	}
	if q, ok := gp0.OriginalQualifiedName(); ok || q != "" {
		t.Fatalf("gp0.OriginalQualifiedName() = (%q, %v), want (\"\", false)", q, ok)
	}

	outer := rhino.Node_NewStringWithToken(rhino.Token_GETPROP, "b")
	outer.AddChildToBack(gp0)
	if q, ok := outer.GetQualifiedNameString(); ok || q != "" {
		t.Fatalf("outer.GetQualifiedNameString() = (%q, %v), want (\"\", false)", q, ok)
	}
	if q, ok := outer.OriginalQualifiedName(); ok || q != "" {
		t.Fatalf("outer.OriginalQualifiedName() = (%q, %v), want (\"\", false)", q, ok)
	}
}

func TestRound96ParseTemplateArgsObjectSingleTemplateParam(t *testing.T) {
	base := rhinotesting.NewBaseJSTypeTestCase(t)

	// Define a nominal `Object` type that only accepts 1 template parameter (`T`).
	closer := base.Registry.GetResolver().OpenForDefinition()
	ctor := base.Registry.CreateConstructorType(
		"Object",
		rhino.NilNode,
		nil,
		nil,
		[]*rhino.TemplateType{base.Registry.CreateTemplateType("T")},
		false)
	objInstance := ctor.GetInstanceType()
	closer.Close()

	scope := rhinotesting.MapBasedScope_EmptyScope()
	base.Registry.DeclareType(scope, "window.Object", objInstance)

	// Direct ParseTemplateArgs on `Object<number>` against a 1-template-param nominal type.
	objNode := rhino.Node_NewString("Object")
	objNode.SetLinenoCharno(12, 4)
	objBlock := rhino.IR_Block()
	objArg := rhino.Node_NewString("number")
	objArg.SetLinenoCharno(12, 11)
	objBlock.AddChildToBack(objArg)
	objNode.AddChildToBack(objBlock)

	base.ErrorReporter.ExpectAllWarnings(
		"Too many template parameters\nFound 2, required at most 1",
		"Too many template parameters\nFound 2, required at most 1")
	args := base.Registry.ParseTemplateArgs(objInstance, objNode, "test.js", scope)
	if len(args) != 1 {
		t.Fatalf("ParseTemplateArgs(Object<number>) len = %d, want 1", len(args))
	}

	// End-to-end CreateTypeFromCommentNodeWithScope on `window.Object<number>`.
	winObjNode := rhino.Node_NewString("window.Object")
	winObjNode.SetLinenoCharno(15, 4)
	winBlock := rhino.IR_Block()
	winArg := rhino.Node_NewString("number")
	winArg.SetLinenoCharno(15, 18)
	winBlock.AddChildToBack(winArg)
	winObjNode.AddChildToBack(winBlock)

	got := base.Registry.CreateTypeFromCommentNodeWithScope(winObjNode, "test.js", scope)
	if got == nil {
		t.Fatal("CreateTypeFromCommentNodeWithScope(window.Object<number>) returned nil")
	}
	base.ErrorReporter.VerifyHasEncounteredAllWarningsAndErrors(t)
}
