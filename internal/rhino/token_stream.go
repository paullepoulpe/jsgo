// Port of com/google/javascript/rhino/TokenStream.java @ closure-compiler v20260415

package rhino

// TokenStream implements the JavaScript scanner. It is based on the C source
// files jsscan.c and jsscan.h in the jsref package.
// Mirrors com.google.javascript.rhino.TokenStream.
type TokenStream struct{}

// IsKeyword reports whether name is an ES3 keyword. Keywords for versions of
// JavaScript after ES3 are not included because the parser would reject valid
// ES3 code that happened to use newer keywords.
// Since IsKeyword is used repeatedly in scanning, it is implemented to do the
// minimum number of character comparisons to ascertain whether the given name
// is a keyword, rather than doing a string search.
// Mirrors TokenStream.isKeyword(String).
func IsKeyword(name string) bool {
	id := false
	s := name
	var x string
outer:
	switch len(s) {
	case 2:
		switch s[1] {
		case 'f':
			if s[0] == 'i' {
				id = true
				break outer
			}
		case 'n':
			if s[0] == 'i' {
				id = true
				break outer
			}
		case 'o':
			if s[0] == 'd' {
				id = true
				break outer
			}
		}
	case 3:
		switch s[0] {
		case 'f':
			if s[2] == 'r' && s[1] == 'o' {
				id = true
				break outer
			}
		case 'i':
			if s[2] == 't' && s[1] == 'n' {
				id = true
				break outer
			}
		case 'n':
			if s[2] == 'w' && s[1] == 'e' {
				id = true
				break outer
			}
		case 't':
			if s[2] == 'y' && s[1] == 'r' {
				id = true
				break outer
			}
		case 'v':
			if s[2] == 'r' && s[1] == 'a' {
				id = true
				break outer
			}
		}
	case 4:
		switch s[0] {
		case 'b':
			x = "byte"
			id = true
		case 'c':
			switch s[3] {
			case 'e':
				if s[2] == 's' && s[1] == 'a' {
					id = true
					break outer
				}
			case 'r':
				if s[2] == 'a' && s[1] == 'h' {
					id = true
					break outer
				}
			}
		case 'e':
			switch s[3] {
			case 'e':
				if s[2] == 's' && s[1] == 'l' {
					id = true
					break outer
				}
			case 'm':
				if s[2] == 'u' && s[1] == 'n' {
					id = true
					break outer
				}
			}
		case 'g':
			x = "goto"
			id = true
		case 'l':
			x = "long"
			id = true
		case 'n':
			x = "null"
			id = true
		case 't':
			switch s[3] {
			case 'e':
				if s[2] == 'u' && s[1] == 'r' {
					id = true
					break outer
				}
			case 's':
				if s[2] == 'i' && s[1] == 'h' {
					id = true
					break outer
				}
			}
		case 'v':
			x = "void"
			id = true
		case 'w':
			x = "with"
			id = true
		}
	case 5:
		switch s[2] {
		case 'a':
			x = "class"
			id = true
		case 'e':
			x = "break"
			id = true
		case 'i':
			x = "while"
			id = true
		case 'l':
			x = "false"
			id = true
		case 'n':
			switch s[0] {
			case 'c':
				x = "const"
				id = true
			case 'f':
				x = "final"
				id = true
			}
		case 'o':
			switch s[0] {
			case 'f':
				x = "float"
				id = true
			case 's':
				x = "short"
				id = true
			}
		case 'p':
			x = "super"
			id = true
		case 'r':
			x = "throw"
			id = true
		case 't':
			x = "catch"
			id = true
		}
	case 6:
		switch s[1] {
		case 'a':
			x = "native"
			id = true
		case 'e':
			switch s[0] {
			case 'd':
				x = "delete"
				id = true
			case 'r':
				x = "return"
				id = true
			}
		case 'h':
			x = "throws"
			id = true
		case 'm':
			x = "import"
			id = true
		case 'o':
			x = "double"
			id = true
		case 't':
			x = "static"
			id = true
		case 'u':
			x = "public"
			id = true
		case 'w':
			x = "switch"
			id = true
		case 'x':
			x = "export"
			id = true
		case 'y':
			x = "typeof"
			id = true
		}
	case 7:
		switch s[1] {
		case 'a':
			x = "package"
			id = true
		case 'e':
			x = "default"
			id = true
		case 'i':
			x = "finally"
			id = true
		case 'o':
			x = "boolean"
			id = true
		case 'r':
			x = "private"
			id = true
		case 'x':
			x = "extends"
			id = true
		}
	case 8:
		switch s[0] {
		case 'a':
			x = "abstract"
			id = true
		case 'c':
			x = "continue"
			id = true
		case 'd':
			x = "debugger"
			id = true
		case 'f':
			x = "function"
			id = true
		case 'v':
			x = "volatile"
			id = true
		}
	case 9:
		switch s[0] {
		case 'i':
			x = "interface"
			id = true
		case 'p':
			x = "protected"
			id = true
		case 't':
			x = "transient"
			id = true
		}
	case 10:
		switch s[1] {
		case 'm':
			x = "implements"
			id = true
		case 'n':
			x = "instanceof"
			id = true
		}
	case 12:
		x = "synchronized"
		id = true
	}
	// Partial match — validate the entire string against the one possibility.
	if x != "" && x != s {
		return false
	}
	return id
}

// IsJSIdentifier reports whether s is a valid JavaScript identifier.
// Mirrors TokenStream.isJSIdentifier(String).
func IsJSIdentifier(s string) bool {
	return JSIdentifierIsValid(s)
}

// JSIdentifierIsValid reports whether s is a valid JavaScript identifier name.
func JSIdentifierIsValid(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !isJSIdentStart(r) {
				return false
			}
		} else {
			if !isJSIdentPart(r) {
				return false
			}
		}
	}
	return true
}

func isJSIdentStart(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || r == '$'
}

func isJSIdentPart(r rune) bool {
	return isJSIdentStart(r) || (r >= '0' && r <= '9')
}
