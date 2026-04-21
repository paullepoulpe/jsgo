// Port of com/google/javascript/jscomp/parsing/parser/Scanner.java
// @ closure-compiler v20260415 (subset).
//
// Java's Scanner is a full ES tokenizer that emits TokenType-bearing
// Token records (IDENTIFIER, NUMBER_LITERAL, STRING_LITERAL, OPEN_CURLY,
// etc.) from a char stream. This file carries the character-level
// scanning primitives the narrow jsgo grammar uses — SkipTrivia,
// ConsumeIdent, ConsumeKeyword, ConsumeBlockComment, position
// tracking — plus the full TokenType enum (mirroring TokenType.java
// verbatim) and literal tokenizers for numbers, strings, templates,
// and regular expressions. The Parser still peeks at runes directly
// in most of its recursive-descent branches; the TokenType enum below
// primarily exists so shared literal helpers can return a matching
// Java-enum value without leaking char-level details.

package parser

import "unicode"

// TokenType ports com/google/javascript/jscomp/parsing/parser/TokenType.java
// @ closure-compiler v20260415 — the full enum, in declaration order.
// Every value in the Java enum has a corresponding constant here so
// downstream additions (lookahead helpers, lexer state machines) can
// return a TokenType that matches what Java's Scanner would emit.
type TokenType int

const (
	TTEndOfFile TokenType = iota
	TTError
	TTIdentifier

	// Keywords.
	TTBreak
	TTCase
	TTCatch
	TTContinue
	TTDebugger
	TTDefault
	TTDelete
	TTDo
	TTElse
	TTFinally
	TTFor
	TTFunction
	TTIf
	TTIn
	TTInstanceof
	TTNew
	TTReturn
	TTSwitch
	TTThis
	TTThrow
	TTTry
	TTTypeof
	TTVar
	TTVoid
	TTWhile
	TTWith
	TTClass
	TTConst
	TTEnum
	TTExport
	TTExtends
	TTImport
	TTSuper
	TTImplements
	TTInterface
	TTLet
	TTPackage
	TTPrivate
	TTProtected
	TTPublic
	TTStatic
	TTYield

	// Punctuators.
	TTOpenCurly
	TTCloseCurly
	TTOpenParen
	TTCloseParen
	TTOpenSquare
	TTCloseSquare
	TTPeriod
	TTSemiColon
	TTComma
	TTOpenAngle               // <
	TTCloseAngle              // >
	TTLessEqual               // <=
	TTGreaterEqual            // >=
	TTArrow                   // =>
	TTEqualEqual              // ==
	TTNotEqual                // !=
	TTEqualEqualEqual         // ===
	TTNotEqualEqual           // !==
	TTPlus                    // +
	TTMinus                   // -
	TTStar                    // *
	TTStarStar                // **
	TTPercent                 // %
	TTPlusPlus                // ++
	TTMinusMinus              // --
	TTLeftShift               // <<
	TTRightShift              // >>
	TTUnsignedRightShift      // >>>
	TTAmpersand               // &
	TTBar                     // |
	TTCaret                   // ^
	TTBang                    // !
	TTTilde                   // ~
	TTAnd                     // &&
	TTOr                      // ||
	TTQuestion                // ?
	TTQuestionQuestion        // ??
	TTQuestionDot             // ?.
	TTColon                   // :
	TTEqual                   // =
	TTPlusEqual               // +=
	TTMinusEqual              // -=
	TTStarEqual               // *=
	TTStarStarEqual           // **=
	TTPercentEqual            // %=
	TTLeftShiftEqual          // <<=
	TTRightShiftEqual         // >>=
	TTUnsignedRightShiftEqual // >>>=
	TTAmpersandEqual          // &=
	TTBarEqual                // |=
	TTCaretEqual              // ^=
	TTSlash                   // /
	TTSlashEqual              // /=
	TTAndEqual                // &&=
	TTOrEqual                 // ||=
	TTQuestionQuestionEqual   // ??=

	// Literals.
	TTNull
	TTTrue
	TTFalse
	TTNumber
	TTString
	TTBigint
	TTRegularExpression
	TTEllipsis // ...

	// Template literal pieces.
	TTTemplateHead           // `foo${   — open + head segment + ${
	TTTemplateMiddle         // }foo${   — middle segment
	TTTemplateTail           // }foo`    — tail segment + close
	TTNoSubstitutionTemplate // `foo`    — entire template, no ${ }

	// TypeScript-ish contextual keywords. Present so round-trip lexers
	// over TS-flavoured sources don't fall off the end of the enum.
	TTType
	TTDeclare
	TTModule
	TTNamespace
)

// String returns the Java-style name of the TokenType. Matches the
// output of TokenType.toString() in Java so diagnostics with
// "expected FOO, got BAR" look the same.
func (t TokenType) String() string {
	switch t {
	case TTEndOfFile:
		return "END_OF_FILE"
	case TTError:
		return "ERROR"
	case TTIdentifier:
		return "IDENTIFIER"
	case TTBreak:
		return "BREAK"
	case TTCase:
		return "CASE"
	case TTCatch:
		return "CATCH"
	case TTContinue:
		return "CONTINUE"
	case TTDebugger:
		return "DEBUGGER"
	case TTDefault:
		return "DEFAULT"
	case TTDelete:
		return "DELETE"
	case TTDo:
		return "DO"
	case TTElse:
		return "ELSE"
	case TTFinally:
		return "FINALLY"
	case TTFor:
		return "FOR"
	case TTFunction:
		return "FUNCTION"
	case TTIf:
		return "IF"
	case TTIn:
		return "IN"
	case TTInstanceof:
		return "INSTANCEOF"
	case TTNew:
		return "NEW"
	case TTReturn:
		return "RETURN"
	case TTSwitch:
		return "SWITCH"
	case TTThis:
		return "THIS"
	case TTThrow:
		return "THROW"
	case TTTry:
		return "TRY"
	case TTTypeof:
		return "TYPEOF"
	case TTVar:
		return "VAR"
	case TTVoid:
		return "VOID"
	case TTWhile:
		return "WHILE"
	case TTWith:
		return "WITH"
	case TTClass:
		return "CLASS"
	case TTConst:
		return "CONST"
	case TTEnum:
		return "ENUM"
	case TTExport:
		return "EXPORT"
	case TTExtends:
		return "EXTENDS"
	case TTImport:
		return "IMPORT"
	case TTSuper:
		return "SUPER"
	case TTImplements:
		return "IMPLEMENTS"
	case TTInterface:
		return "INTERFACE"
	case TTLet:
		return "LET"
	case TTPackage:
		return "PACKAGE"
	case TTPrivate:
		return "PRIVATE"
	case TTProtected:
		return "PROTECTED"
	case TTPublic:
		return "PUBLIC"
	case TTStatic:
		return "STATIC"
	case TTYield:
		return "YIELD"
	case TTOpenCurly:
		return "OPEN_CURLY"
	case TTCloseCurly:
		return "CLOSE_CURLY"
	case TTOpenParen:
		return "OPEN_PAREN"
	case TTCloseParen:
		return "CLOSE_PAREN"
	case TTOpenSquare:
		return "OPEN_SQUARE"
	case TTCloseSquare:
		return "CLOSE_SQUARE"
	case TTPeriod:
		return "PERIOD"
	case TTSemiColon:
		return "SEMI_COLON"
	case TTComma:
		return "COMMA"
	case TTOpenAngle:
		return "OPEN_ANGLE"
	case TTCloseAngle:
		return "CLOSE_ANGLE"
	case TTLessEqual:
		return "LESS_EQUAL"
	case TTGreaterEqual:
		return "GREATER_EQUAL"
	case TTArrow:
		return "ARROW"
	case TTEqualEqual:
		return "EQUAL_EQUAL"
	case TTNotEqual:
		return "NOT_EQUAL"
	case TTEqualEqualEqual:
		return "EQUAL_EQUAL_EQUAL"
	case TTNotEqualEqual:
		return "NOT_EQUAL_EQUAL"
	case TTPlus:
		return "PLUS"
	case TTMinus:
		return "MINUS"
	case TTStar:
		return "STAR"
	case TTStarStar:
		return "STAR_STAR"
	case TTPercent:
		return "PERCENT"
	case TTPlusPlus:
		return "PLUS_PLUS"
	case TTMinusMinus:
		return "MINUS_MINUS"
	case TTLeftShift:
		return "LEFT_SHIFT"
	case TTRightShift:
		return "RIGHT_SHIFT"
	case TTUnsignedRightShift:
		return "UNSIGNED_RIGHT_SHIFT"
	case TTAmpersand:
		return "AMPERSAND"
	case TTBar:
		return "BAR"
	case TTCaret:
		return "CARET"
	case TTBang:
		return "BANG"
	case TTTilde:
		return "TILDE"
	case TTAnd:
		return "AND"
	case TTOr:
		return "OR"
	case TTQuestion:
		return "QUESTION"
	case TTQuestionQuestion:
		return "QUESTION_QUESTION"
	case TTQuestionDot:
		return "QUESTION_DOT"
	case TTColon:
		return "COLON"
	case TTEqual:
		return "EQUAL"
	case TTPlusEqual:
		return "PLUS_EQUAL"
	case TTMinusEqual:
		return "MINUS_EQUAL"
	case TTStarEqual:
		return "STAR_EQUAL"
	case TTStarStarEqual:
		return "STAR_STAR_EQUAL"
	case TTPercentEqual:
		return "PERCENT_EQUAL"
	case TTLeftShiftEqual:
		return "LEFT_SHIFT_EQUAL"
	case TTRightShiftEqual:
		return "RIGHT_SHIFT_EQUAL"
	case TTUnsignedRightShiftEqual:
		return "UNSIGNED_RIGHT_SHIFT_EQUAL"
	case TTAmpersandEqual:
		return "AMPERSAND_EQUAL"
	case TTBarEqual:
		return "BAR_EQUAL"
	case TTCaretEqual:
		return "CARET_EQUAL"
	case TTSlash:
		return "SLASH"
	case TTSlashEqual:
		return "SLASH_EQUAL"
	case TTAndEqual:
		return "AND_EQUAL"
	case TTOrEqual:
		return "OR_EQUAL"
	case TTQuestionQuestionEqual:
		return "QUESTION_QUESTION_EQUAL"
	case TTNull:
		return "NULL"
	case TTTrue:
		return "TRUE"
	case TTFalse:
		return "FALSE"
	case TTNumber:
		return "NUMBER"
	case TTString:
		return "STRING"
	case TTBigint:
		return "BIGINT"
	case TTRegularExpression:
		return "REGULAR_EXPRESSION"
	case TTEllipsis:
		return "ELLIPSIS"
	case TTTemplateHead:
		return "TEMPLATE_HEAD"
	case TTTemplateMiddle:
		return "TEMPLATE_MIDDLE"
	case TTTemplateTail:
		return "TEMPLATE_TAIL"
	case TTNoSubstitutionTemplate:
		return "NO_SUBSTITUTION_TEMPLATE"
	case TTType:
		return "TYPE"
	case TTDeclare:
		return "DECLARE"
	case TTModule:
		return "MODULE"
	case TTNamespace:
		return "NAMESPACE"
	}
	return "UNKNOWN"
}

// Scanner walks a source-string by rune. All positions are reported
// as 1-indexed line, 0-indexed col to match Java's getLineno /
// getCharno convention.
type Scanner struct {
	src     []rune
	pos     int
	line    int
	lineBOL int
	// lastTokenLine is the line of the most recently consumed
	// non-trivia character. Used by the parser's ASI path to decide
	// whether the next token sits on a later line. Trivia-only
	// consumption (SkipTrivia, ConsumeLineComment, non-JSDoc block
	// comments) does not touch this field.
	lastTokenLine int
}

// NewScanner mirrors Java's `new Scanner(source)` constructor.
func NewScanner(src string) *Scanner {
	return &Scanner{src: []rune(src), line: 1, lastTokenLine: 1}
}

// LastTokenLine returns the line of the most recently consumed
// non-trivia character. The parser compares this to Line() to
// decide whether ES5 §7.9 ASI should fire after a statement that's
// missing its terminating `;`.
func (s *Scanner) LastTokenLine() int { return s.lastTokenLine }

// MarkTokenConsumed updates lastTokenLine to the current scanner
// line. The Parser calls this after any Advance-based token
// consumption that doesn't go through one of Scanner's Consume*
// helpers (e.g. parseNumberLiteral / parseStringLiteral loops).
func (s *Scanner) MarkTokenConsumed() { s.lastTokenLine = s.line }

// ---- basic position / peek / advance ----------------------------- //

// Text returns the substring between two byte-positions previously
// observed via Pos(). Useful for literal-value extraction
// (parseStringLiteral, parseNumberLiteral) without callers reaching
// into the underlying rune slice.
func (s *Scanner) Text(start, end int) string {
	if start < 0 {
		start = 0
	}
	if end > len(s.src) {
		end = len(s.src)
	}
	if start >= end {
		return ""
	}
	return string(s.src[start:end])
}

func (s *Scanner) Pos() int    { return s.pos }
func (s *Scanner) Line() int   { return s.line }
func (s *Scanner) Col() int    { return s.pos - s.lineBOL }
func (s *Scanner) AtEnd() bool { return s.pos >= len(s.src) }

// Peek returns the current rune without consuming. Returns 0 at EOF.
func (s *Scanner) Peek() rune {
	if s.AtEnd() {
		return 0
	}
	return s.src[s.pos]
}

// PeekAt returns the rune at offset ahead of the cursor, or 0 past end.
func (s *Scanner) PeekAt(offset int) rune {
	i := s.pos + offset
	if i < 0 || i >= len(s.src) {
		return 0
	}
	return s.src[i]
}

// Advance consumes one rune and updates line/col tracking.
func (s *Scanner) Advance() rune {
	if s.AtEnd() {
		return 0
	}
	c := s.src[s.pos]
	s.pos++
	if c == '\n' {
		s.line++
		s.lineBOL = s.pos
	}
	return c
}

// Save / Restore give the parser a checkpoint it can rewind to for
// speculative parse paths. Mirrors how Java's Scanner is coupled with
// the parser via Token-level markSavePoint / resetTo methods.
type ScannerMark struct {
	pos, line, lineBOL, lastTokenLine int
}

func (s *Scanner) Save() ScannerMark {
	return ScannerMark{s.pos, s.line, s.lineBOL, s.lastTokenLine}
}
func (s *Scanner) Restore(m ScannerMark) {
	s.pos, s.line, s.lineBOL, s.lastTokenLine = m.pos, m.line, m.lineBOL, m.lastTokenLine
}

// ---- byte-level conditionals ------------------------------------- //

// ConsumeByte advances past `b` if the current rune matches. Returns
// whether it consumed.
func (s *Scanner) ConsumeByte(b byte) bool {
	if s.AtEnd() || byte(s.src[s.pos]) != b {
		return false
	}
	s.Advance()
	s.lastTokenLine = s.line
	return true
}

// ---- identifier / keyword --------------------------------------- //

// ConsumeIdent reads a JS identifier starting at the cursor.
// Returns ("", false) if the cursor isn't on an ident start or if the
// identifier matches a reserved word (Java's Scanner keeps keyword +
// identifier disjoint via its TokenType enum; we emulate that by
// refusing to return reserved names, which prevents `if`, `return`,
// etc. from being wrongly consumed as identifiers in speculative
// statement parses).
func (s *Scanner) ConsumeIdent() (string, bool) {
	if s.AtEnd() {
		return "", false
	}
	c := s.src[s.pos]
	if !isIdentStart(c) {
		return "", false
	}
	start := s.pos
	end := s.pos
	for end < len(s.src) && isIdentCont(s.src[end]) {
		end++
	}
	ident := string(s.src[start:end])
	if isReservedWord(ident) {
		return "", false
	}
	for s.pos < end {
		s.Advance()
	}
	s.lastTokenLine = s.line
	return ident, true
}

// isReservedWord returns true for JS reserved words + the subset of
// contextual keywords we must not misread as identifiers. Keeping the
// list in sync with Parser.java's ReservedWord handling is important
// for parser correctness; anything missing here risks a bad AST (see
// `if (x) { … }` being parsed as `if(x) { … };`).
func isReservedWord(s string) bool {
	switch s {
	// Control flow / declarations the parser might speculatively
	// treat as identifiers in the LHS slot.
	case "break", "case", "catch", "class", "const", "continue",
		"debugger", "default", "do", "else", "export",
		"extends", "finally", "for", "function", "if", "import",
		"return", "switch", "throw", "try", "var", "while",
		"let":
		return true
	// Literal-shaped keywords that are parsed in specific
	// expression branches (tryPrimary consumes them before the
	// identifier fallback runs).
	case "true", "false", "null", "undefined", "this":
		return true
	// Operator-shaped keywords handled by tryUnary / peekBinOp.
	case "new", "typeof", "void", "delete", "in", "instanceof":
		return true
	}
	// Contextual keywords (`super`, `yield`, `async`, `await`,
	// `static`, `enum`, `with`) intentionally pass through as
	// identifiers — the scope / checker layers give them meaning,
	// and some corpus fixtures use them unadorned.
	return false
}

// AtKeyword reports whether `kw` is the next keyword (followed by a
// non-ident boundary) at the cursor — non-consuming. Used by the
// binary-op parser to decide whether to consume `instanceof` / `in`
// without advancing the scanner.
func (s *Scanner) AtKeyword(kw string) bool {
	rs := []rune(kw)
	if s.pos+len(rs) > len(s.src) {
		return false
	}
	for i, r := range rs {
		if s.src[s.pos+i] != r {
			return false
		}
	}
	if s.pos+len(rs) < len(s.src) && isIdentCont(s.src[s.pos+len(rs)]) {
		return false
	}
	return true
}

// ConsumeKeywordRaw advances past a keyword that the caller has
// already verified via AtKeyword. Skips the reserved-word filter that
// ConsumeIdent applies, so keywords known to the binary-op parser
// (instanceof, in) can still be consumed.
func (s *Scanner) ConsumeKeywordRaw(kw string) {
	for i := 0; i < len([]rune(kw)); i++ {
		s.Advance()
	}
	s.lastTokenLine = s.line
}

// ConsumeKeyword matches an exact keyword followed by a non-ident
// boundary. Returns whether the keyword was consumed.
func (s *Scanner) ConsumeKeyword(kw string) bool {
	rs := []rune(kw)
	if s.pos+len(rs) > len(s.src) {
		return false
	}
	for i, r := range rs {
		if s.src[s.pos+i] != r {
			return false
		}
	}
	if s.pos+len(rs) < len(s.src) {
		nxt := s.src[s.pos+len(rs)]
		if isIdentCont(nxt) {
			return false
		}
	}
	for i := 0; i < len(rs); i++ {
		s.Advance()
	}
	s.lastTokenLine = s.line
	return true
}

// ---- trivia (whitespace + comments) ----------------------------- //

// ConsumeBlockComment reads a `/* … */` block comment (including the
// `/**` JSDoc variant). Returns the comment body (the text between
// the delimiters) and whether a comment was consumed.
func (s *Scanner) ConsumeBlockComment() (string, bool) {
	if s.pos+1 >= len(s.src) || s.src[s.pos] != '/' || s.src[s.pos+1] != '*' {
		return "", false
	}
	s.Advance()
	s.Advance()
	start := s.pos
	for s.pos+1 < len(s.src) && !(s.src[s.pos] == '*' && s.src[s.pos+1] == '/') {
		s.Advance()
	}
	if s.pos+1 >= len(s.src) {
		return "", false
	}
	body := string(s.src[start:s.pos])
	s.Advance()
	s.Advance()
	s.lastTokenLine = s.line
	return body, true
}

// ConsumeLineComment reads a `// …\n` comment. Returns whether it
// consumed one.
func (s *Scanner) ConsumeLineComment() bool {
	if s.pos+1 < len(s.src) && s.src[s.pos] == '/' && s.src[s.pos+1] == '/' {
		for !s.AtEnd() && s.src[s.pos] != '\n' {
			s.Advance()
		}
		return true
	}
	return false
}

// SkipTrivia consumes whitespace plus plain (non-JSDoc) block + line
// comments. JSDoc block comments (`/**`) are left for the parser to
// attach to the following declaration.
func (s *Scanner) SkipTrivia() {
	s.skipTriviaTracked()
}

// SkipTriviaSeenNewline is SkipTrivia but reports whether the skipped
// region crossed a line boundary. Used by the parser's ASI path —
// when a statement isn't terminated by `;` but the next token sits on
// a later line, ES5 §7.9 says to insert a semicolon.
func (s *Scanner) SkipTriviaSeenNewline() bool {
	return s.skipTriviaTracked()
}

func (s *Scanner) skipTriviaTracked() bool {
	startLine := s.line
	for !s.AtEnd() {
		c := s.src[s.pos]
		if unicode.IsSpace(c) {
			s.Advance()
			continue
		}
		if s.ConsumeLineComment() {
			continue
		}
		if s.pos+2 < len(s.src) && c == '/' && s.src[s.pos+1] == '*' && s.src[s.pos+2] != '*' {
			s.Advance()
			s.Advance()
			for s.pos+1 < len(s.src) && !(s.src[s.pos] == '*' && s.src[s.pos+1] == '/') {
				s.Advance()
			}
			if s.pos+1 < len(s.src) {
				s.Advance()
				s.Advance()
			}
			continue
		}
		break
	}
	return s.line != startLine
}

// SkipToStatementBoundary walks forward until just past the next `;`
// or `}` at top-level, or to EOF. Used by the Parser to recover from
// unrecognised input without losing subsequent valid statements.
// Strings/comments are coarsely handled so the scan doesn't mistake
// content for a boundary.
func (s *Scanner) SkipToStatementBoundary() {
	depth := 0
	for !s.AtEnd() {
		c := s.src[s.pos]
		switch c {
		case ';':
			s.Advance()
			if depth == 0 {
				return
			}
		case '{':
			depth++
			s.Advance()
		case '}':
			if depth == 0 {
				return
			}
			depth--
			s.Advance()
			if depth == 0 {
				return
			}
		case '/':
			if s.pos+1 < len(s.src) && s.src[s.pos+1] == '/' {
				for !s.AtEnd() && s.src[s.pos] != '\n' {
					s.Advance()
				}
			} else if s.pos+1 < len(s.src) && s.src[s.pos+1] == '*' {
				s.Advance()
				s.Advance()
				for s.pos+1 < len(s.src) && !(s.src[s.pos] == '*' && s.src[s.pos+1] == '/') {
					s.Advance()
				}
				if s.pos+1 < len(s.src) {
					s.Advance()
					s.Advance()
				}
			} else {
				s.Advance()
			}
		case '"', '\'':
			s.Advance()
			for !s.AtEnd() && s.src[s.pos] != c {
				if s.src[s.pos] == '\\' && s.pos+1 < len(s.src) {
					s.Advance()
					s.Advance()
					continue
				}
				if s.src[s.pos] == '\n' {
					break
				}
				s.Advance()
			}
			if !s.AtEnd() {
				s.Advance()
			}
		default:
			s.Advance()
		}
	}
}

// ---- character-class predicates --------------------------------- //

func isIdentStart(c rune) bool {
	return unicode.IsLetter(c) || c == '_' || c == '$'
}

func isIdentCont(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' || c == '$'
}

// ---- full literal tokenizers ------------------------------------ //
// These helpers mirror Scanner.java's scanNumericLiteral /
// scanStringLiteral / scanRegularExpressionLiteral / scanTemplate*
// entry points. They consume the entire literal from the cursor and
// return a text + structured-value tuple. Callers build the matching
// rhino.Node. Each helper leaves the scanner just past the literal
// on success, and rewinds to the cursor start on failure.

// ScanNumber reads a JS numeric literal starting at the cursor. Supports
// decimal (incl. `1_000` separators), hex (0x), octal (0o), binary
// (0b), exponent (`1e10`), fraction (`.5`), and a trailing bigint `n`.
// Returns (text, isBigInt, ok); `ok=false` means the cursor wasn't on a
// numeric literal and nothing was consumed.
func (s *Scanner) ScanNumber() (string, bool, bool) {
	if s.AtEnd() {
		return "", false, false
	}
	start := s.pos
	c := s.Peek()
	// Leading `.5` form.
	if c == '.' {
		if !unicode.IsDigit(s.PeekAt(1)) {
			return "", false, false
		}
		s.Advance()
		for !s.AtEnd() && (unicode.IsDigit(s.Peek()) || s.Peek() == '_') {
			s.Advance()
		}
		s.scanExponent()
		s.lastTokenLine = s.line
		return string(s.src[start:s.pos]), false, true
	}
	if !unicode.IsDigit(c) {
		return "", false, false
	}
	// Hex / oct / bin prefixes.
	if c == '0' && s.pos+1 < len(s.src) {
		nxt := s.src[s.pos+1]
		switch nxt {
		case 'x', 'X':
			s.Advance()
			s.Advance()
			for !s.AtEnd() && (isHexDigit(s.Peek()) || s.Peek() == '_') {
				s.Advance()
			}
			isBig := false
			if !s.AtEnd() && s.Peek() == 'n' {
				isBig = true
				s.Advance()
			}
			s.lastTokenLine = s.line
			return string(s.src[start:s.pos]), isBig, true
		case 'o', 'O':
			s.Advance()
			s.Advance()
			for !s.AtEnd() && (isOctDigit(s.Peek()) || s.Peek() == '_') {
				s.Advance()
			}
			isBig := false
			if !s.AtEnd() && s.Peek() == 'n' {
				isBig = true
				s.Advance()
			}
			s.lastTokenLine = s.line
			return string(s.src[start:s.pos]), isBig, true
		case 'b', 'B':
			s.Advance()
			s.Advance()
			for !s.AtEnd() && (s.Peek() == '0' || s.Peek() == '1' || s.Peek() == '_') {
				s.Advance()
			}
			isBig := false
			if !s.AtEnd() && s.Peek() == 'n' {
				isBig = true
				s.Advance()
			}
			s.lastTokenLine = s.line
			return string(s.src[start:s.pos]), isBig, true
		}
	}
	// Decimal + fraction + exponent.
	for !s.AtEnd() && (unicode.IsDigit(s.Peek()) || s.Peek() == '_') {
		s.Advance()
	}
	if !s.AtEnd() && s.Peek() == '.' {
		s.Advance()
		for !s.AtEnd() && (unicode.IsDigit(s.Peek()) || s.Peek() == '_') {
			s.Advance()
		}
	}
	s.scanExponent()
	isBig := false
	if !s.AtEnd() && s.Peek() == 'n' {
		// BigInt can't have a fraction or exponent — but we leave that
		// validation to downstream parser checks. Consume the suffix so
		// the returned text is `123n`.
		isBig = true
		s.Advance()
	}
	s.lastTokenLine = s.line
	return string(s.src[start:s.pos]), isBig, true
}

func (s *Scanner) scanExponent() {
	if s.AtEnd() {
		return
	}
	c := s.Peek()
	if c != 'e' && c != 'E' {
		return
	}
	s.Advance()
	if !s.AtEnd() && (s.Peek() == '+' || s.Peek() == '-') {
		s.Advance()
	}
	for !s.AtEnd() && (unicode.IsDigit(s.Peek()) || s.Peek() == '_') {
		s.Advance()
	}
}

// ScanString reads a `'…'` or `"…"` string literal (including escape
// sequences). The opening quote must be at the cursor. Returns the
// decoded value (escapes processed) and whether scanning succeeded
// (e.g. on a lone newline inside a single-line string we return ok=false).
func (s *Scanner) ScanString() (string, bool) {
	if s.AtEnd() {
		return "", false
	}
	quote := s.Peek()
	if quote != '\'' && quote != '"' {
		return "", false
	}
	s.Advance()
	var b []rune
	for !s.AtEnd() {
		c := s.Peek()
		if c == quote {
			s.Advance()
			s.lastTokenLine = s.line
			return string(b), true
		}
		if c == '\n' || c == '\r' {
			return "", false
		}
		if c == '\\' {
			s.Advance()
			dec, ok := s.scanEscapeChar()
			if !ok {
				return "", false
			}
			b = append(b, dec...)
			continue
		}
		b = append(b, c)
		s.Advance()
	}
	return "", false
}

// scanEscapeChar decodes a single `\X` escape after the backslash is
// consumed. Supports standard JS escapes (`\n`, `\t`, …), `\xHH`,
// `\uHHHH`, `\u{HHHH}`, octal escapes, and `\\r?\n` line continuations
// (which produce zero output). Returns the decoded rune(s) and ok.
func (s *Scanner) scanEscapeChar() ([]rune, bool) {
	if s.AtEnd() {
		return nil, false
	}
	c := s.Advance()
	switch c {
	case 'n':
		return []rune{'\n'}, true
	case 't':
		return []rune{'\t'}, true
	case 'r':
		return []rune{'\r'}, true
	case 'b':
		return []rune{'\b'}, true
	case 'f':
		return []rune{'\f'}, true
	case 'v':
		return []rune{'\v'}, true
	case '0':
		// `\0` is the NUL escape — but only when not followed by a digit
		// (otherwise it's a legacy octal, which the parser accepts in
		// non-strict mode).
		if !s.AtEnd() && isOctDigit(s.Peek()) {
			return s.scanOctalTail('0')
		}
		return []rune{0}, true
	case '\'', '"', '\\', '`':
		return []rune{c}, true
	case '\n':
		// Line continuation — produces nothing in the decoded value.
		return nil, true
	case '\r':
		if !s.AtEnd() && s.Peek() == '\n' {
			s.Advance()
		}
		return nil, true
	case 'x':
		if s.pos+1 >= len(s.src) || !isHexDigit(s.src[s.pos]) || !isHexDigit(s.src[s.pos+1]) {
			return nil, false
		}
		hi := hexVal(s.Advance())
		lo := hexVal(s.Advance())
		return []rune{rune(hi*16 + lo)}, true
	case 'u':
		if !s.AtEnd() && s.Peek() == '{' {
			s.Advance()
			v := 0
			for !s.AtEnd() && s.Peek() != '}' {
				if !isHexDigit(s.Peek()) {
					return nil, false
				}
				v = v*16 + hexVal(s.Advance())
				if v > 0x10FFFF {
					return nil, false
				}
			}
			if s.AtEnd() {
				return nil, false
			}
			s.Advance() // '}'
			return []rune{rune(v)}, true
		}
		if s.pos+3 >= len(s.src) {
			return nil, false
		}
		v := 0
		for i := 0; i < 4; i++ {
			if !isHexDigit(s.Peek()) {
				return nil, false
			}
			v = v*16 + hexVal(s.Advance())
		}
		return []rune{rune(v)}, true
	}
	if isOctDigit(c) {
		return s.scanOctalTail(c)
	}
	return []rune{c}, true
}

// scanOctalTail reads up to two more octal digits after an initial
// digit. Matches the legacy JS octal escape form (e.g. `\077`).
func (s *Scanner) scanOctalTail(first rune) ([]rune, bool) {
	v := int(first - '0')
	for i := 0; i < 2; i++ {
		if s.AtEnd() || !isOctDigit(s.Peek()) {
			break
		}
		v = v*8 + int(s.Advance()-'0')
	}
	return []rune{rune(v)}, true
}

func isHexDigit(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isOctDigit(c rune) bool {
	return c >= '0' && c <= '7'
}

func hexVal(c rune) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	case c >= 'A' && c <= 'F':
		return int(c-'A') + 10
	}
	return 0
}

// ScanRegex reads a `/pattern/flags` regular-expression literal. The
// cursor must be at the leading `/`. Returns (pattern, flags, ok).
// Mirrors Scanner.java's scanRegularExpressionLiteralPortion +
// scanRegularExpressionFlags. `ok=false` means nothing was consumed.
//
// Regex vs division disambiguation is left to the caller — by the time
// it reaches this helper, the parser has already decided the leading
// `/` starts a regex (via expression-context heuristics).
func (s *Scanner) ScanRegex() (string, string, bool) {
	if s.AtEnd() || s.Peek() != '/' {
		return "", "", false
	}
	save := s.Save()
	s.Advance() // '/'
	start := s.pos
	inClass := false
	for !s.AtEnd() {
		c := s.Peek()
		if c == '\n' || c == '\r' {
			s.Restore(save)
			return "", "", false
		}
		if c == '\\' {
			s.Advance()
			if s.AtEnd() {
				s.Restore(save)
				return "", "", false
			}
			s.Advance()
			continue
		}
		if c == '[' {
			inClass = true
			s.Advance()
			continue
		}
		if c == ']' && inClass {
			inClass = false
			s.Advance()
			continue
		}
		if c == '/' && !inClass {
			break
		}
		s.Advance()
	}
	if s.AtEnd() {
		s.Restore(save)
		return "", "", false
	}
	pat := string(s.src[start:s.pos])
	s.Advance() // closing '/'
	flagStart := s.pos
	for !s.AtEnd() && isIdentCont(s.Peek()) {
		s.Advance()
	}
	flags := string(s.src[flagStart:s.pos])
	s.lastTokenLine = s.line
	return pat, flags, true
}

// TemplateSegment represents one piece of a template literal: the
// contiguous text between backticks / `${` / `}`. `Tail` is true when
// the segment is the terminating piece (ends at a backtick rather
// than a `${`). Matches TemplateLiteralToken's TEMPLATE_HEAD /
// TEMPLATE_MIDDLE / TEMPLATE_TAIL / NO_SUBSTITUTION_TEMPLATE
// distinction in Java's Scanner.
type TemplateSegment struct {
	Text string
	Tail bool
}

// ScanTemplateStart reads from a leading backtick through the first
// segment boundary (either `${` or closing backtick). Returns the
// segment and whether the scan succeeded. Mirrors Scanner.java's
// scanTemplateLiteral entry point.
func (s *Scanner) ScanTemplateStart() (TemplateSegment, bool) {
	if s.AtEnd() || s.Peek() != '`' {
		return TemplateSegment{}, false
	}
	s.Advance() // opening backtick
	return s.scanTemplateSegment()
}

// ScanTemplateContinue reads from a `}` that terminates a substitution
// expression through the next segment boundary. The caller must have
// already consumed the `}`. Mirrors Scanner.java's
// nextTemplateLiteralTokenImpl.
func (s *Scanner) ScanTemplateContinue() (TemplateSegment, bool) {
	return s.scanTemplateSegment()
}

func (s *Scanner) scanTemplateSegment() (TemplateSegment, bool) {
	var b []rune
	for !s.AtEnd() {
		c := s.Peek()
		if c == '`' {
			s.Advance()
			s.lastTokenLine = s.line
			return TemplateSegment{Text: string(b), Tail: true}, true
		}
		if c == '$' && s.PeekAt(1) == '{' {
			s.Advance()
			s.Advance()
			s.lastTokenLine = s.line
			return TemplateSegment{Text: string(b), Tail: false}, true
		}
		if c == '\\' {
			s.Advance()
			dec, ok := s.scanEscapeChar()
			if !ok {
				return TemplateSegment{}, false
			}
			b = append(b, dec...)
			continue
		}
		b = append(b, c)
		s.Advance()
	}
	return TemplateSegment{}, false
}
