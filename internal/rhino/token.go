// Port of com/google/javascript/rhino/Token.java @ closure-compiler v20260415
//
// Only the subset of Token enum values needed for Tier-0 fixtures is
// defined here. More are added as corpus tiers grow.

package rhino

type Token int

const (
	ERROR Token = iota
	SCRIPT
	VAR
	NAME
	STRINGLIT
	NUMBER
)

func (t Token) String() string {
	switch t {
	case SCRIPT:
		return "SCRIPT"
	case VAR:
		return "VAR"
	case NAME:
		return "NAME"
	case STRINGLIT:
		return "STRINGLIT"
	case NUMBER:
		return "NUMBER"
	}
	return "ERROR"
}
