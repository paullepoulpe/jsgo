package javaast_test

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/tools/javaport/javaast"
)

// testCtx is a minimal MethodContext used across all test cases.
var testCtx = javaast.MethodContext{
	RecvName:  "c",
	ClassName: "MyClass",
	TypeMap: func(jt string) string {
		switch jt {
		case "boolean":
			return "bool"
		case "int":
			return "int"
		case "String":
			return "string"
		case "void":
			return ""
		}
		return jt
	},
}

// translate is a helper that runs TranslateBody and strips leading/trailing whitespace.
func translate(t *testing.T, javaBody string) string {
	t.Helper()
	got, ok := javaast.TranslateBody(javaBody, testCtx)
	if !ok {
		t.Fatalf("TranslateBody returned ok=false for:\n%s", javaBody)
	}
	return strings.TrimSpace(got)
}

// trivial is a helper that runs TryTrivial.
func trivial(t *testing.T, javaBody string) (string, bool) {
	t.Helper()
	return javaast.TryTrivial(javaBody, "c", "MyClass", func(jt string) string { return jt })
}

// ---- #161 Trivial body patterns ----------------------------------------

func TestTrivial_Getter(t *testing.T) {
	cases := []struct {
		body string
		want string
	}{
		{
			body: "return foo;",
			want: "\treturn c.foo\n",
		},
		{
			body: "return this.foo;",
			want: "\treturn c.foo\n",
		},
	}
	for _, tc := range cases {
		got, ok := trivial(t, tc.body)
		if !ok {
			t.Errorf("TryTrivial(%q) returned false", tc.body)
			continue
		}
		if got != tc.want {
			t.Errorf("TryTrivial(%q)\n  got:  %q\n  want: %q", tc.body, got, tc.want)
		}
	}
}

func TestTrivial_Setter(t *testing.T) {
	body := "this.foo = foo;"
	got, ok := trivial(t, body)
	if !ok {
		t.Fatalf("TryTrivial(%q) returned false", body)
	}
	want := "\tc.foo = foo\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTrivial_Delegation_Return(t *testing.T) {
	body := "return other.bar(y);"
	got, ok := trivial(t, body)
	if !ok {
		t.Fatalf("TryTrivial(%q) returned false", body)
	}
	want := "\treturn c.other.Bar(y)\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTrivial_Delegation_Void(t *testing.T) {
	body := "other.bar(y);"
	got, ok := trivial(t, body)
	if !ok {
		t.Fatalf("TryTrivial(%q) returned false", body)
	}
	want := "\tc.other.Bar(y)\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestTrivial_NullCheck(t *testing.T) {
	body := "if (x == null) return y;"
	got, ok := trivial(t, body)
	if !ok {
		t.Fatalf("TryTrivial(%q) returned false", body)
	}
	if !strings.Contains(got, "if x == nil") {
		t.Errorf("expected 'if x == nil' in output, got: %q", got)
	}
	if !strings.Contains(got, "return") {
		t.Errorf("expected return in output, got: %q", got)
	}
}

func TestTrivial_NullCheck_NotEq(t *testing.T) {
	body := "if (x != null) return y;"
	got, ok := trivial(t, body)
	if !ok {
		t.Fatalf("TryTrivial(%q) returned false", body)
	}
	if !strings.Contains(got, "if x != nil") {
		t.Errorf("expected 'if x != nil' in output, got: %q", got)
	}
}

func TestTrivial_FieldInit(t *testing.T) {
	body := "this.foo = bar;\nthis.baz = qux;"
	got, ok := trivial(t, body)
	if !ok {
		t.Fatalf("TryTrivial(%q) returned false", body)
	}
	if !strings.Contains(got, "c.foo = bar") {
		t.Errorf("expected 'c.foo = bar' in output, got: %q", got)
	}
	if !strings.Contains(got, "c.baz = qux") {
		t.Errorf("expected 'c.baz = qux' in output, got: %q", got)
	}
}

func TestTrivial_SingleReturn(t *testing.T) {
	body := "return x + y;"
	got, ok := trivial(t, body)
	if !ok {
		t.Fatalf("TryTrivial(%q) returned false", body)
	}
	if !strings.Contains(got, "return x + y") {
		t.Errorf("expected 'return x + y' in output, got: %q", got)
	}
}

// ---- #165 Control-flow patterns ----------------------------------------

func TestCF_IfElse(t *testing.T) {
	body := `
if (x == null) {
    return 0;
} else {
    return 1;
}
`
	got := translate(t, body)
	if !strings.Contains(got, "if x == nil") {
		t.Errorf("expected 'if x == nil', got:\n%s", got)
	}
	if !strings.Contains(got, "} else {") {
		t.Errorf("expected '} else {', got:\n%s", got)
	}
}

func TestCF_IfElseIf(t *testing.T) {
	body := `
if (a > 0) {
    return 1;
} else if (a < 0) {
    return -1;
} else {
    return 0;
}
`
	got := translate(t, body)
	if !strings.Contains(got, "else if") {
		t.Errorf("expected 'else if', got:\n%s", got)
	}
}

func TestCF_ForClassic(t *testing.T) {
	body := `
for (int i = 0; i < n; i++) {
    process(i);
}
`
	got := translate(t, body)
	if !strings.Contains(got, "for") {
		t.Errorf("expected 'for', got:\n%s", got)
	}
	if !strings.Contains(got, "i < n") {
		t.Errorf("expected 'i < n', got:\n%s", got)
	}
}

func TestCF_ForEach(t *testing.T) {
	body := `
for (String s : list) {
    process(s);
}
`
	got := translate(t, body)
	if !strings.Contains(got, "range") {
		t.Errorf("expected 'range', got:\n%s", got)
	}
	if !strings.Contains(got, "_, s") {
		t.Errorf("expected '_, s', got:\n%s", got)
	}
}

func TestCF_While(t *testing.T) {
	body := `
while (x > 0) {
    x--;
}
`
	got := translate(t, body)
	if !strings.Contains(got, "for x > 0") {
		t.Errorf("expected 'for x > 0', got:\n%s", got)
	}
}

func TestCF_DoWhile(t *testing.T) {
	body := `
do {
    x++;
} while (x < 10);
`
	got := translate(t, body)
	if !strings.Contains(got, "for {") {
		t.Errorf("expected 'for {' for do-while, got:\n%s", got)
	}
	if !strings.Contains(got, "break") {
		t.Errorf("expected break in do-while, got:\n%s", got)
	}
}

func TestCF_Switch(t *testing.T) {
	body := `
switch (token) {
    case 1:
        return "one";
    case 2:
        return "two";
    default:
        return "other";
}
`
	got := translate(t, body)
	if !strings.Contains(got, "switch token") {
		t.Errorf("expected 'switch token', got:\n%s", got)
	}
	if !strings.Contains(got, "case 1:") {
		t.Errorf("expected 'case 1:', got:\n%s", got)
	}
	if !strings.Contains(got, "default:") {
		t.Errorf("expected 'default:', got:\n%s", got)
	}
}

func TestCF_SwitchFallthrough(t *testing.T) {
	body := `
switch (x) {
    case 1:
        doA();
    case 2:
        doB();
        break;
}
`
	got := translate(t, body)
	if !strings.Contains(got, "fallthrough") {
		t.Errorf("expected 'fallthrough', got:\n%s", got)
	}
}

func TestCF_TryCatch(t *testing.T) {
	body := `
try {
    doSomething();
} catch (Exception e) {
    handleError(e);
}
`
	got := translate(t, body)
	if !strings.Contains(got, "defer func()") {
		t.Errorf("expected 'defer func()', got:\n%s", got)
	}
	if !strings.Contains(got, "recover()") {
		t.Errorf("expected 'recover()', got:\n%s", got)
	}
}

func TestCF_TryFinally(t *testing.T) {
	body := `
try {
    doSomething();
} finally {
    cleanup();
}
`
	got := translate(t, body)
	if !strings.Contains(got, "defer func()") {
		t.Errorf("expected 'defer func()', got:\n%s", got)
	}
	if !strings.Contains(got, "Cleanup()") && !strings.Contains(got, "cleanup()") {
		t.Errorf("expected cleanup call in finally block, got:\n%s", got)
	}
}

func TestCF_Throw(t *testing.T) {
	body := `throw new IllegalArgumentException("bad value");`
	got := translate(t, body)
	if !strings.Contains(got, "panic(") {
		t.Errorf("expected 'panic(', got:\n%s", got)
	}
}

func TestCF_Break(t *testing.T) {
	body := `
for (int i = 0; i < 10; i++) {
    if (i == 5) {
        break;
    }
}
`
	got := translate(t, body)
	if !strings.Contains(got, "break") {
		t.Errorf("expected 'break', got:\n%s", got)
	}
}

func TestCF_Continue(t *testing.T) {
	body := `
for (int i = 0; i < 10; i++) {
    if (i % 2 == 0) {
        continue;
    }
    process(i);
}
`
	got := translate(t, body)
	if !strings.Contains(got, "continue") {
		t.Errorf("expected 'continue', got:\n%s", got)
	}
}

func TestCF_LocalVarDecl(t *testing.T) {
	body := `
int x = 42;
String name = "hello";
`
	got := translate(t, body)
	if !strings.Contains(got, "x := 42") {
		t.Errorf("expected 'x := 42', got:\n%s", got)
	}
	if !strings.Contains(got, `name := "hello"`) {
		t.Errorf(`expected 'name := "hello"', got:\n%s`, got)
	}
}

func TestCF_LocalVarNoInit(t *testing.T) {
	body := `int count;`
	got := translate(t, body)
	if !strings.Contains(got, "var count") {
		t.Errorf("expected 'var count', got:\n%s", got)
	}
}

func TestCF_CompoundAssign(t *testing.T) {
	body := `x += 5;`
	got := translate(t, body)
	if !strings.Contains(got, "x += 5") {
		t.Errorf("expected 'x += 5', got:\n%s", got)
	}
}

func TestCF_NewObject(t *testing.T) {
	body := `return new Foo(a, b);`
	got := translate(t, body)
	if !strings.Contains(got, "NewFoo(") {
		t.Errorf("expected 'NewFoo(', got:\n%s", got)
	}
}

func TestCF_ArrayAccess(t *testing.T) {
	body := `return arr[i];`
	got := translate(t, body)
	if !strings.Contains(got, "arr[i]") {
		t.Errorf("expected 'arr[i]', got:\n%s", got)
	}
}

func TestCF_NullComparison_Binary(t *testing.T) {
	body := `
if (x == null) {
    return 0;
}
`
	got := translate(t, body)
	if !strings.Contains(got, "x == nil") {
		t.Errorf("expected 'x == nil', got:\n%s", got)
	}
}

func TestCF_BreakLabeled(t *testing.T) {
	body := `
outer:
for (int i = 0; i < 10; i++) {
    for (int j = 0; j < 10; j++) {
        if (i + j > 5) {
            break outer;
        }
    }
}
`
	got := translate(t, body)
	if !strings.Contains(got, "break outer") {
		t.Errorf("expected 'break outer', got:\n%s", got)
	}
}

func TestCF_ChainedMethodCall(t *testing.T) {
	body := `return foo.bar().baz();`
	got := translate(t, body)
	// Should contain chained calls
	if !strings.Contains(got, "Bar()") && !strings.Contains(got, "bar()") {
		t.Errorf("expected method call chain, got:\n%s", got)
	}
}

func TestCF_NewArray(t *testing.T) {
	body := `int[] arr = new int[10];`
	got := translate(t, body)
	if !strings.Contains(got, "make(") || !strings.Contains(got, "10") {
		t.Errorf("expected make(..., 10), got:\n%s", got)
	}
}
