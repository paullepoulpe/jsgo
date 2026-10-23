package parser

import (
	"strings"
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Expected strings were printed by the pinned Java classes (toString() of
// SourcePosition, SourceRange, ParseTree, the tokens and FeatureSet, and the
// FormalParameterListTree checkArgument message). Go positions carry no
// source file, so they print like Java positions whose file name is empty.
func TestToStringMatchesJava(t *testing.T) {
	p0 := NewSourcePosition(0, 0, 0)
	p5 := NewSourcePosition(5, 1, 2)
	r := NewSourceRange(p0, p5)
	cases := []struct{ got, want string }{
		{p5.String(), "(2, 3)"},
		{r.String(), "<(1, 1) - (2, 3)>"},
		{NewBlockTree(r, []ParseTree{}).String(), "BLOCK@<(1, 1) - (2, 3)>"},
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
	if want := "Unexpected # of comma and formal params.|params: []|comma positions: [(1, 1), (2, 3)]"; strings.ReplaceAll(msg, "\n", "|") != want {
		t.Errorf("checkArgument message %q, Java %q", msg, want)
	}
	// Token casts are type assertions, like Java's (T) this.
	var tok Token = NewIdentifierToken(r, "x")
	if tok.(*IdentifierToken).GetValue() != "x" {
		t.Errorf("AsIdentifier")
	}
	var lit Token = NewStringLiteralToken("'s'", r, false)
	if lit.(LiteralToken).GetLiteralTokenBase().Value != "'s'" {
		t.Errorf("AsLiteral")
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
