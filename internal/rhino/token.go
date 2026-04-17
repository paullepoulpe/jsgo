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
	LET
	CONST
	NAME
	STRINGLIT
	NUMBER
	TRUE
	FALSE
	NULL
	UNDEFINED
	ASSIGN
	EXPR_RESULT
	FUNCTION
	PARAM_LIST
	BLOCK
	CALL
	RETURN
)

func (t Token) String() string {
	switch t {
	case SCRIPT:
		return "SCRIPT"
	case VAR:
		return "VAR"
	case LET:
		return "LET"
	case CONST:
		return "CONST"
	case NAME:
		return "NAME"
	case STRINGLIT:
		return "STRINGLIT"
	case NUMBER:
		return "NUMBER"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case NULL:
		return "NULL"
	case UNDEFINED:
		return "UNDEFINED"
	case ASSIGN:
		return "ASSIGN"
	case EXPR_RESULT:
		return "EXPR_RESULT"
	case FUNCTION:
		return "FUNCTION"
	case PARAM_LIST:
		return "PARAM_LIST"
	case BLOCK:
		return "BLOCK"
	case CALL:
		return "CALL"
	case RETURN:
		return "RETURN"
	}
	return "ERROR"
}
