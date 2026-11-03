package args4j

import "testing"

type args4jMutationBean struct {
	out        string
	withDepend string
	withForbid string
}

func (b *args4jMutationBean) Args4jFields() []AnnotatedField {
	return []AnnotatedField{
		{Option: &Option{Name: "--out"}, Field: String(&b.out)},
		{Option: &Option{Name: "--with-missing-dep", Depends: []string{"--unregistered"}}, Field: String(&b.withDepend)},
		{Option: &Option{Name: "--with-missing-forbid", Forbids: []string{"--unregistered"}}, Field: String(&b.withForbid)},
	}
}

func TestTryParseArgumentDoesNotMutateCallerArgs(t *testing.T) {
	bean := &args4jMutationBean{}
	p := NewCmdLineParser(bean)
	callerArgs := []string{"--out=bundle.js"}
	if err := p.TryParseArgument(callerArgs...); err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	if bean.out != "bundle.js" {
		t.Fatalf("bean.out = %q, want bundle.js", bean.out)
	}
	if callerArgs[0] != "--out=bundle.js" {
		t.Fatalf("TryParseArgument mutated caller slice: got %q, want --out=bundle.js", callerArgs[0])
	}
}

func TestDependsAndForbidsWithUnregisteredOption(t *testing.T) {
	bean := &args4jMutationBean{}
	p := NewCmdLineParser(bean)
	if err := p.TryParseArgument("--with-missing-dep", "x"); err == nil {
		t.Fatal("expected error when dependent option is unregistered/absent")
	}

	bean2 := &args4jMutationBean{}
	p2 := NewCmdLineParser(bean2)
	if err := p2.TryParseArgument("--with-missing-forbid", "y"); err != nil {
		t.Fatalf("unexpected error for unregistered forbidden option: %v", err)
	}
}
