// Port of com/google/javascript/jscomp/parsing/JsDocToken.java @ closure-compiler v20260415

package parsing

// JsDocToken holds the JSDoc-specific tokens.
//
// This class is based on Rhino's Token.
type JsDocToken int32

const (
	// Tokens recycled from Rhino

	JsDocToken_EOF JsDocToken = iota // end of file token - (not EOF_CHAR)
	JsDocToken_EOL                   // end of line
	JsDocToken_LEFT_ANGLE
	JsDocToken_RIGHT_ANGLE
	JsDocToken_STRING
	JsDocToken_LEFT_SQUARE
	JsDocToken_RIGHT_SQUARE
	JsDocToken_LEFT_CURLY
	JsDocToken_RIGHT_CURLY
	JsDocToken_LEFT_PAREN
	JsDocToken_RIGHT_PAREN
	JsDocToken_COMMA // comma operator
	JsDocToken_COLON
	JsDocToken_ITER_REST

	// JsDoc-only tokens

	JsDocToken_ANNOTATION
	JsDocToken_PIPE
	JsDocToken_STAR
	JsDocToken_EOC
	JsDocToken_QMARK
	JsDocToken_BANG
	JsDocToken_EQUALS
)

var jsDocToken_names = [...]string{
	"EOF",
	"EOL",
	"LEFT_ANGLE",
	"RIGHT_ANGLE",
	"STRING",
	"LEFT_SQUARE",
	"RIGHT_SQUARE",
	"LEFT_CURLY",
	"RIGHT_CURLY",
	"LEFT_PAREN",
	"RIGHT_PAREN",
	"COMMA",
	"COLON",
	"ITER_REST",
	"ANNOTATION",
	"PIPE",
	"STAR",
	"EOC",
	"QMARK",
	"BANG",
	"EQUALS",
}

// JsDocToken_values is JsDocToken.values().
//
// Java: JsDocToken.values()
func JsDocToken_values() []JsDocToken {
	values := make([]JsDocToken, len(jsDocToken_names))
	for i := range values {
		values[i] = JsDocToken(i)
	}
	return values
}

// Name is Enum.name().
//
// Java: JsDocToken.name()
func (t JsDocToken) Name() string { return jsDocToken_names[t] }

// String is Enum.toString().
//
// Java: JsDocToken.toString()
func (t JsDocToken) String() string { return t.Name() }
