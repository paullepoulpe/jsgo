// Port of com/google/javascript/jscomp/parsing/JsDocTokenStreamTest.java @ closure-compiler v20260415

package parsing

import "testing"

// Tests for JsDocTokenStream.

func TestJsDocTokenization1(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_STAR, JsDocToken_ANNOTATION, JsDocToken_LEFT_CURLY, JsDocToken_STRING, JsDocToken_RIGHT_CURLY, JsDocToken_EOL, JsDocToken_STAR, JsDocToken_ANNOTATION}
	strings := []string{"type", "string", "private"}
	testJSDocTokenStream(t, " * @type {string}\n * @private", tokens, strings)
	testJSDocTokenStream(t, " *    @type { string } \n * @private", tokens, strings)
	testJSDocTokenStream(t, " * @type   {  string}\n * @private", tokens, strings)
	testJSDocTokenStream(t, " * @type {string  }\n * @private", tokens, strings)
	testJSDocTokenStream(t, " * @type {string}\n *   @private", tokens, strings)
	testJSDocTokenStream(t, " * @type {string}   \n * @private", tokens, strings)
}

func TestJsDocTokenization2(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_ANNOTATION, JsDocToken_LEFT_CURLY, JsDocToken_STRING, JsDocToken_LEFT_ANGLE, JsDocToken_STRING, JsDocToken_PIPE, JsDocToken_STRING, JsDocToken_RIGHT_ANGLE, JsDocToken_RIGHT_CURLY}
	strings := []string{"param", "Array", "string", "null"}
	testJSDocTokenStream(t, "@param {Array.<string|null>}", tokens, strings)
	testJSDocTokenStream(t, "@param {Array.<string|null>}", tokens, strings)
	testJSDocTokenStream(t, "@param {Array.<string |null>}", tokens, strings)
	testJSDocTokenStream(t, " @param {Array.<string |  null>}", tokens, strings)
	testJSDocTokenStream(t, " @param {Array.<string|null  >}", tokens, strings)
	testJSDocTokenStream(t, "@param {Array  .<string|null>}", tokens, strings)
	testJSDocTokenStream(t, "@param   {Array.<string|null>}", tokens, strings)
	testJSDocTokenStream(t, "@param {  Array.<string|null>}", tokens, strings)
	testJSDocTokenStream(t, "@param {Array.<string|   null>}  ", tokens, strings)
	testJSDocTokenStream(t, "@param {Array.<string|null>}", tokens, strings)
	testJSDocTokenStream(t, "     @param { Array .< string |null > } ", tokens, strings)
}

func TestJsDocTokenization4(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_ANNOTATION, JsDocToken_LEFT_CURLY, JsDocToken_STRING, JsDocToken_LEFT_ANGLE, JsDocToken_LEFT_PAREN, JsDocToken_STRING, JsDocToken_COMMA, JsDocToken_STRING, JsDocToken_RIGHT_PAREN, JsDocToken_RIGHT_ANGLE, JsDocToken_RIGHT_CURLY, JsDocToken_EOF}
	strings := []string{"param", "Array", "string", "null"}
	testJSDocTokenStream(t, "@param {Array.<(string,null)>}", tokens, strings)
	testJSDocTokenStream(t, "@param {Array  .<(string,null)> } ", tokens, strings)
	testJSDocTokenStream(t, " @param {Array.<  (  string,null)>}", tokens, strings)
	testJSDocTokenStream(t, "@param {Array.<(string  , null)>}", tokens, strings)
	testJSDocTokenStream(t, "@param {Array.<(string,   null)  > }  ", tokens, strings)
	testJSDocTokenStream(t, "@param {  Array  .<  (string,null)>}   ", tokens, strings)
}

func TestJsDocTokenization5(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_ANNOTATION, JsDocToken_STRING, JsDocToken_EOC, JsDocToken_EOF}
	strings := []string{"param", "foo.Bar"}
	testJSDocTokenStream(t, "@param foo.Bar*/", tokens, strings)
	testJSDocTokenStream(t, " @param   foo.Bar*/", tokens, strings)
	testJSDocTokenStream(t, " @param foo.Bar   */", tokens, strings)
}

func TestJsDocTokenization6(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_ANNOTATION, JsDocToken_EOL, JsDocToken_ANNOTATION, JsDocToken_EOL, JsDocToken_ANNOTATION, JsDocToken_EOC}
	strings := []string{"private", "static", "desc"}
	testJSDocTokenStream(t, "@private\n@static\n@desc*/", tokens, strings)
	testJSDocTokenStream(t, "@private\n @static\n@desc*/", tokens, strings)
	testJSDocTokenStream(t, "@private\n@static\n @desc*/", tokens, strings)
	testJSDocTokenStream(t, "@private\n@static\n@desc */", tokens, strings)
	testJSDocTokenStream(t, " @private \n@static\n @desc*/", tokens, strings)
	testJSDocTokenStream(t, "@private\n@static    \n @desc  */", tokens, strings)
	testJSDocTokenStream(t, "@private\n@static\n@desc*/", tokens, strings)
	testJSDocTokenStream(t, "@private   \n@static   \n @desc*/", tokens, strings)
}

func TestJsDocTokenization7(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_ITER_REST, JsDocToken_ITER_REST, JsDocToken_ITER_REST, JsDocToken_ITER_REST, JsDocToken_ITER_REST, JsDocToken_LEFT_ANGLE, JsDocToken_EOC}
	strings := []string{}

	testJSDocTokenStream(t, "................<*/", tokens, strings)
	testJSDocTokenStream(t, "............... .<*/", tokens, strings)
	testJSDocTokenStream(t, "................< */", tokens, strings)
	testJSDocTokenStream(t, "............... .< */", tokens, strings)
	testJSDocTokenStream(t, "............... .< */ ", tokens, strings)
	testJSDocTokenStream(t, " ............... .< */ ", tokens, strings)
}

func TestJsDocTokenization8(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_STAR, JsDocToken_ANNOTATION, JsDocToken_STRING, JsDocToken_STRING, JsDocToken_STRING, JsDocToken_STRING, JsDocToken_STRING, JsDocToken_STRING, JsDocToken_STRING, JsDocToken_EOL, JsDocToken_EOC}
	strings := []string{"param", "foo.Bar", "opt_name", "this", "parameter", "is", "a", "name"}
	testJSDocTokenStream(t, "* @param foo.Bar opt_name this parameter is a name\n*/\n", tokens, strings)
	testJSDocTokenStream(t, " *  @param foo.Bar opt_name this parameter is a name \n*/ \n", tokens, strings)
}

func TestJsDocTokenization9(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_STAR, JsDocToken_ANNOTATION, JsDocToken_STRING, JsDocToken_STRING, JsDocToken_STRING, JsDocToken_STRING, JsDocToken_STRING, JsDocToken_ANNOTATION, JsDocToken_STRING, JsDocToken_EOL, JsDocToken_EOC}
	strings := []string{"param", "foo.Bar", "opt_name", "this", "parameter", "does", "media", "blah"}
	testJSDocTokenStream(t, "* @param foo.Bar opt_name this parameter does @media blah\n*/\n", tokens, strings)
}

func TestJsDocTokenization10(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_STRING, JsDocToken_LEFT_ANGLE, JsDocToken_STRING, JsDocToken_RIGHT_ANGLE, JsDocToken_EOC}
	strings := []string{"Array", "String"}
	testJSDocTokenStream(t, "Array<String>*/", tokens, strings)
}

func TestJsDocTokenization11(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_ANNOTATION, JsDocToken_LEFT_CURLY, JsDocToken_STRING, JsDocToken_QMARK, JsDocToken_RIGHT_CURLY, JsDocToken_EOC, JsDocToken_EOF}
	strings := []string{"param", "string"}
	testJSDocTokenStream(t, "@param {string?}*/", tokens, strings)
	testJSDocTokenStream(t, " @param {string?}*/", tokens, strings)
	testJSDocTokenStream(t, "@param { string?}*/", tokens, strings)
	testJSDocTokenStream(t, "@param {string ?}*/", tokens, strings)
	testJSDocTokenStream(t, "@param  {string ?  } */", tokens, strings)
	testJSDocTokenStream(t, "@param { string  ?  }*/", tokens, strings)
	testJSDocTokenStream(t, "@param {string?  }*/", tokens, strings)
}

func TestJsDocTokenization12(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_STRING, JsDocToken_ITER_REST, JsDocToken_EOC}
	strings := []string{"function"}

	testJSDocTokenStream(t, "function ...*/", tokens, strings)
}

func TestJsDocTokenization13(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_ITER_REST, JsDocToken_LEFT_SQUARE, JsDocToken_STRING, JsDocToken_RIGHT_SQUARE, JsDocToken_EOC}
	strings := []string{"number"}

	testJSDocTokenStream(t, "...[number]*/", tokens, strings)
}

func TestJsDocTokenization14(t *testing.T) {
	// Since ES4 type parsing only requires to parse an ITER_REST when it is
	// followed by a comma (,) we are allowing this case to parse this way.
	// This is a simplification of the tokenizer, but the extra complexity is
	// never used.
	tokens := []JsDocToken{JsDocToken_STRING, JsDocToken_LEFT_SQUARE, JsDocToken_STRING, JsDocToken_EOC}
	strings := []string{"foo", "bar..."}

	testJSDocTokenStream(t, "foo[ bar...*/", tokens, strings)
}

func TestJsDocTokenization15(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_STRING, JsDocToken_LEFT_SQUARE, JsDocToken_STRING, JsDocToken_COMMA, JsDocToken_ITER_REST, JsDocToken_EOC}
	strings := []string{"foo", "bar"}

	testJSDocTokenStream(t, "foo[ bar,...*/", tokens, strings)
	testJSDocTokenStream(t, "foo[ bar ,...*/", tokens, strings)
	testJSDocTokenStream(t, "foo[bar, ...*/", tokens, strings)
	testJSDocTokenStream(t, "foo[ bar  ,   ...  */", tokens, strings)
	testJSDocTokenStream(t, "foo [bar,... */", tokens, strings)
}

func TestJsDocTokenization16(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_STRING, JsDocToken_COLON, JsDocToken_COLON, JsDocToken_COLON, JsDocToken_ITER_REST, JsDocToken_STRING, JsDocToken_COLON, JsDocToken_STRING, JsDocToken_EOC}
	strings := []string{"foo", "bar", "bar2"}

	testJSDocTokenStream(t, "foo:::...bar:bar2*/", tokens, strings)
}

func TestJsDocTokenization17(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_STRING, JsDocToken_EOL, JsDocToken_EOC}
	strings := []string{".."}

	testJSDocTokenStream(t, "..\n*/", tokens, strings)
}

func TestJsDocTokenization18(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_STRING, JsDocToken_EOL, JsDocToken_EOC}
	strings := []string{"."}

	testJSDocTokenStream(t, ".\n*/", tokens, strings)
}

func TestJsDocTokenization19(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_ANNOTATION, JsDocToken_LEFT_CURLY, JsDocToken_STAR, JsDocToken_RIGHT_CURLY, JsDocToken_EOC}
	strings := []string{"type", "*"}

	testJSDocTokenStream(t, "@type {*}*/", tokens, strings)
}

func TestJsDocTokenization20(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_ANNOTATION, JsDocToken_LEFT_CURLY, JsDocToken_BANG, JsDocToken_STRING, JsDocToken_RIGHT_CURLY, JsDocToken_EOC, JsDocToken_EOF}
	strings := []string{"param", "Object"}
	testJSDocTokenStream(t, "@param {!Object}*/", tokens, strings)
	testJSDocTokenStream(t, " @param {!Object}*/", tokens, strings)
	testJSDocTokenStream(t, "@param {! Object}*/", tokens, strings)
	testJSDocTokenStream(t, "@param { !Object}*/", tokens, strings)
	testJSDocTokenStream(t, "@param  {!Object  } */", tokens, strings)
	testJSDocTokenStream(t, "@param {  ! Object  }*/", tokens, strings)
	testJSDocTokenStream(t, "@param {!Object  }*/", tokens, strings)
}

func TestJsDocTokenization21(t *testing.T) {
	tokens := []JsDocToken{JsDocToken_ANNOTATION, JsDocToken_LEFT_CURLY, JsDocToken_STRING, JsDocToken_EQUALS, JsDocToken_RIGHT_CURLY, JsDocToken_EOC, JsDocToken_EOF}
	strings := []string{"param", "Object"}
	testJSDocTokenStream(t, "@param {Object=}*/", tokens, strings)
	testJSDocTokenStream(t, " @param {Object=}*/", tokens, strings)
	testJSDocTokenStream(t, "@param { Object =}*/", tokens, strings)
	testJSDocTokenStream(t, "@param { Object=}*/", tokens, strings)
	testJSDocTokenStream(t, "@param  {Object=  } */", tokens, strings)
	testJSDocTokenStream(t, "@param { Object = }*/", tokens, strings)
	testJSDocTokenStream(t, "@param {Object=  }*/", tokens, strings)
}

func testJSDocTokenStream(t *testing.T, comment string, tokens []JsDocToken, strings []string) {
	t.Helper()
	stream := NewJsDocTokenStreamWithLineno(comment, 0)
	stringsIndex := 0
	for _, token := range tokens {
		readToken := stream.GetJsDocToken()

		// token equality
		if token != readToken {
			t.Fatalf("%q: got token %v, want %v", comment, readToken, token)
		}

		// string equality
		if token == JsDocToken_ANNOTATION || token == JsDocToken_STRING {
			if got, want := stream.GetString(), strings[stringsIndex]; got != want {
				t.Fatalf("%q: got string %q, want %q", comment, got, want)
			}
			stringsIndex++
		}
	}
}
