package parser

import (
	"reflect"
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Expected strings were printed by the pinned Java classes (toString() of
// SourcePosition, SourceRange, ParseTree, the tokens and FeatureSet, and the
// FormalParameterListTree checkArgument message). Go positions carry no
// source file, so they print like Java positions whose file name is empty.
//
// Go-only: SourcePosition is an offset and prints as "@offset"; a
// LinePosition (with line and column) prints like the Java position.
func TestToStringMatchesJava(t *testing.T) {
	p0 := NewSourcePosition(0)
	p5 := NewSourcePosition(5)
	r := NewSourceRange(p0, p5)
	l0 := LinePosition{Offset: 0, Line: 0, Column: 0}
	l5 := LinePosition{Offset: 5, Line: 1, Column: 2}
	cases := []struct{ got, want string }{
		{l5.String(), "(2, 3)"},
		{LineRange{l0, l5}.String(), "<(1, 1) - (2, 3)>"},
		{p5.String(), "@5"},
		{r.String(), "<@0 - @5>"},
		{NewBlockTree(r, []ParseTree{}).String(), "BLOCK@<@0 - @5>"},
		{NewToken(TokenType_OPEN_CURLY, r).String(), "{"},
		{NewToken(TokenType_VAR, r).String(), "var"},
		{NewToken(TokenType_END_OF_FILE, r).String(), "End of File"},
		{NewIdentifierToken(r, "#priv").String(), "#priv"},
		{NewLiteralToken(TokenType_NUMBER, "1_0", r).String(), "1_0"},
		{FeatureSet_BARE_MINIMUM.With(FeatureSet_Feature_MODULES).With(FeatureSet_Feature_GETTER).String(), "[getters, modules]"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("got %q, Java %q", c.got, c.want)
		}
	}
	var msg string
	jutil.Try(func() {
		NewFormalParameterListTree(r, []ParseTree{}, false, []SourcePosition{p0, p5})
	}, func(e jutil.IIllegalArgumentException) { msg = e.GetMessage() })
	if want := "Unexpected # of comma and formal params.|params: []|comma positions: [@0, @5]"; strings.ReplaceAll(msg, "\n", "|") != want {
		t.Errorf("checkArgument message %q, Java %q", msg, want)
	}
	// Token is a 16-byte pointer-free value type.
	var tok Token = NewIdentifierToken(r, "x")
	if tok.GetValue() != "x" {
		t.Errorf("GetValue")
	}
	var lit Token = NewStringLiteralToken("'s'", r, false)
	if lit.String() != "'s'" {
		t.Errorf("StringLiteralToken.String")
	}
	var tree ParseTree = NewParenExpressionTree(r, NewArrayPatternTree(r, []ParseTree{}))
	if !tree.IsPattern() || !tree.IsValidAssignmentTarget() || tree.IsValidNonVanillaAssignmentTarget() {
		t.Errorf("paren-wrapped pattern predicates")
	}
	threw := false
	jutil.Try(func() { NewIdentifierToken(r, "#p").GetValue() }, func(e jutil.IIllegalStateException) { threw = true })
	if !threw {
		t.Errorf("GetValue on a private identifier must throw IllegalStateException")
	}
}

func typeContainsPointers(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Pointer, reflect.String, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface, reflect.UnsafePointer:
		return true
	case reflect.Array:
		return typeContainsPointers(t.Elem())
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			if typeContainsPointers(t.Field(i).Type) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func TestTokenValueLayout(t *testing.T) {
	tt := reflect.TypeOf(Token{})
	if tt.Size() != 16 {
		t.Fatalf("Token size = %d, want 16", tt.Size())
	}
	if typeContainsPointers(tt) {
		t.Fatalf("Token contains pointers; expected pointer-free (noscan) value struct")
	}
	var zero Token
	if !zero.IsNull() || zero.IsPresent() {
		t.Fatalf("zero Token must report IsNull=true, IsPresent=false")
	}
}
