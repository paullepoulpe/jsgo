// Port of com/google/javascript/jscomp/AstManipulationsTest.java @ closure-compiler v20260415

package jscomp

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

func assertNodeEquivalent(t *testing.T, got, want *rhino.Node) {
	t.Helper()
	if !got.IsEquivalentTo(want) {
		t.Errorf("got\n%s\nwant\n%s", got.ToStringTree(), want.ToStringTree())
	}
}

func TestAstManipulations_FuseExpressionsFusesTwoNumbers(t *testing.T) {
	one := rhino.IR_Number(1)
	two := rhino.IR_Number(2)

	assertNodeEquivalent(t, AstManipulations_FuseExpressions(one, two),
		rhino.IR_Comma(one.CloneNode(), two.CloneNode()))
}

func TestAstManipulations_FuseExpressionsDropsEmptySecondExpression(t *testing.T) {
	one := rhino.IR_Number(1)
	empty := rhino.IR_Empty()

	assertNodeEquivalent(t, AstManipulations_FuseExpressions(one, empty), one.CloneNode())
}

func TestAstManipulations_FuseExpressionsKeepsEmptyFirstExpression(t *testing.T) {
	empty := rhino.IR_Empty()
	one := rhino.IR_Number(1)

	// Note: it should be fine to return one instead of this comma
	assertNodeEquivalent(t, AstManipulations_FuseExpressions(empty, one),
		rhino.NewNodeWith2Children(rhino.Token_COMMA, empty.CloneNode(), one.CloneNode()))
}

func TestAstManipulations_FuseExpressionsFusesNumberIntoComma(t *testing.T) {
	zero := rhino.IR_Number(0)
	one := rhino.IR_Number(1)
	zeroCommaOne := rhino.IR_Comma(zero, one)
	two := rhino.IR_Number(2)

	assertNodeEquivalent(t, AstManipulations_FuseExpressions(zeroCommaOne, two),
		rhino.IR_Comma(zeroCommaOne.CloneTree(), two.CloneNode()))
}

func TestAstManipulations_FuseExpressionsFusesCommaIntoNumber(t *testing.T) {
	zero := rhino.IR_Number(0)
	one := rhino.IR_Number(1)
	two := rhino.IR_Number(2)
	oneCommaTwo := rhino.IR_Comma(one, two)

	// Utility keeps commas in the left branch of the comma tree, i.e.
	//  COMMA
	//    COMMA
	//      NUMBER 0
	//      NUMBER 1
	//    NUMBER 2
	// instead of this:
	//  COMMA
	//    NUMBER 0
	//    COMMA
	//      NUMBER 1
	//      NUMBER 2
	// to make unit-testing passes that use this utility easier. Writing
	// 'expected' JS code for the latter tree requires more parentheses.
	assertNodeEquivalent(t, AstManipulations_FuseExpressions(zero, oneCommaTwo),
		rhino.IR_Comma(rhino.IR_Comma(zero.CloneNode(), one.CloneNode()), two.CloneNode()))
}
