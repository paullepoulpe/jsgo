// Port of com/google/javascript/jscomp/parsing/JsDocToken.java @ closure-compiler v20260415

package parsing

// JsDocToken is the lexical category returned by JsDocTokenStream.
// The enum preserves Java's names and ordering so the parser reads
// 1:1 against the Java side.
type JsDocToken int

const (
	JsDocTokenEOF JsDocToken = iota // end of file token (not EOF_CHAR)
	JsDocTokenEOL                   // end of line
	JsDocTokenLeftAngle
	JsDocTokenRightAngle
	JsDocTokenString
	JsDocTokenLeftSquare
	JsDocTokenRightSquare
	JsDocTokenLeftCurly
	JsDocTokenRightCurly
	JsDocTokenLeftParen
	JsDocTokenRightParen
	JsDocTokenComma
	JsDocTokenColon
	JsDocTokenIterRest
	JsDocTokenAnnotation
	JsDocTokenPipe
	JsDocTokenStar
	JsDocTokenEOC // end-of-comment (`*/`)
	JsDocTokenQMark
	JsDocTokenBang
	JsDocTokenEquals
)

func (t JsDocToken) String() string {
	switch t {
	case JsDocTokenEOF:
		return "EOF"
	case JsDocTokenEOL:
		return "EOL"
	case JsDocTokenLeftAngle:
		return "LEFT_ANGLE"
	case JsDocTokenRightAngle:
		return "RIGHT_ANGLE"
	case JsDocTokenString:
		return "STRING"
	case JsDocTokenLeftSquare:
		return "LEFT_SQUARE"
	case JsDocTokenRightSquare:
		return "RIGHT_SQUARE"
	case JsDocTokenLeftCurly:
		return "LEFT_CURLY"
	case JsDocTokenRightCurly:
		return "RIGHT_CURLY"
	case JsDocTokenLeftParen:
		return "LEFT_PAREN"
	case JsDocTokenRightParen:
		return "RIGHT_PAREN"
	case JsDocTokenComma:
		return "COMMA"
	case JsDocTokenColon:
		return "COLON"
	case JsDocTokenIterRest:
		return "ITER_REST"
	case JsDocTokenAnnotation:
		return "ANNOTATION"
	case JsDocTokenPipe:
		return "PIPE"
	case JsDocTokenStar:
		return "STAR"
	case JsDocTokenEOC:
		return "EOC"
	case JsDocTokenQMark:
		return "QMARK"
	case JsDocTokenBang:
		return "BANG"
	case JsDocTokenEquals:
		return "EQUALS"
	}
	return "UNKNOWN"
}
