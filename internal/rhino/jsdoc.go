// Port of com/google/javascript/rhino/JSDocInfo.java @ closure-compiler v20260415
//
// The Java class is substantial; this port keeps only the @type annotation
// (as an unresolved JSTypeExpression). Further annotations (@param,
// @return, @template, …) are added as later tiers require them.

package rhino

type JSDocInfo struct {
	typeExpr *JSTypeExpression
}

func (i *JSDocInfo) Type() *JSTypeExpression { return i.typeExpr }
func (i *JSDocInfo) SetType(t *JSTypeExpression) { i.typeExpr = t }
func (i *JSDocInfo) HasType() bool { return i != nil && i.typeExpr != nil }
