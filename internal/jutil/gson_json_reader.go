package jutil

// Port of com/google/gson/stream/JsonReader.java (and JsonToken, JsonScope)
// @ gson 2.9.1, the version shaded into the pinned closure-compiler.jar.
//
// Go-only: language-adapter — Gson is a library dependency of Closure; it is
// ported here as Java-semantics runtime. The reader keeps Gson's 1024-char
// buffer and refill logic because they are observable (a number longer than
// the buffer is read as an unquoted literal, and positions in error messages
// are computed from the buffer). nextDouble/nextLong/nextInt (which need
// Java's Double.parseDouble) and promoteNameToValue are not ported.

// JsonToken is com.google.gson.stream.JsonToken.
type JsonToken int32

const (
	JsonToken_BEGIN_ARRAY JsonToken = iota
	JsonToken_END_ARRAY
	JsonToken_BEGIN_OBJECT
	JsonToken_END_OBJECT
	JsonToken_NAME
	JsonToken_STRING
	JsonToken_NUMBER
	JsonToken_BOOLEAN
	JsonToken_NULL
	JsonToken_END_DOCUMENT
)

var jsonTokenNames = [...]string{
	"BEGIN_ARRAY", "END_ARRAY", "BEGIN_OBJECT", "END_OBJECT", "NAME",
	"STRING", "NUMBER", "BOOLEAN", "NULL", "END_DOCUMENT",
}

// String is the enum constant's name().
func (t JsonToken) String() string { return jsonTokenNames[t] }

// com.google.gson.stream.JsonScope
const (
	jsonScope_EMPTY_ARRAY       = 1
	jsonScope_NONEMPTY_ARRAY    = 2
	jsonScope_EMPTY_OBJECT      = 3
	jsonScope_DANGLING_NAME     = 4
	jsonScope_NONEMPTY_OBJECT   = 5
	jsonScope_EMPTY_DOCUMENT    = 6
	jsonScope_NONEMPTY_DOCUMENT = 7
	jsonScope_CLOSED            = 8
)

// JavaReader is java.io.Reader's read(char[], int, int) / close().
type JavaReader interface {
	// Read is read(char[] cbuf, int off, int len): the number of chars read,
	// or -1 at the end of the stream.
	Read(cbuf []uint16, off, length int) int
	// Close is close().
	Close()
}

// StringReader is java.io.StringReader.
type StringReader struct {
	str  []uint16
	next int
}

// NewStringReader is new StringReader(String).
func NewStringReader(s string) *StringReader {
	return &StringReader{str: StringToUTF16(s)}
}

// Read is StringReader.read(char[], int, int).
func (r *StringReader) Read(cbuf []uint16, off, length int) int {
	if off < 0 || length < 0 || off+length > len(cbuf) {
		panic(NewIndexOutOfBoundsException("Range [" + itoa(off) + ", " + itoa(off) + " + " + itoa(length) + ") out of bounds for length " + itoa(len(cbuf))))
	}
	if length == 0 {
		return 0
	}
	if r.next >= len(r.str) {
		return -1
	}
	n := min(len(r.str)-r.next, length)
	copy(cbuf[off:off+n], r.str[r.next:r.next+n])
	r.next += n
	return n
}

// Close is StringReader.close().
func (r *StringReader) Close() {}

const jsonReader_MIN_INCOMPLETE_INTEGER = -0x8000000000000000 / 10

const (
	jsonReader_PEEKED_NONE               = 0
	jsonReader_PEEKED_BEGIN_OBJECT       = 1
	jsonReader_PEEKED_END_OBJECT         = 2
	jsonReader_PEEKED_BEGIN_ARRAY        = 3
	jsonReader_PEEKED_END_ARRAY          = 4
	jsonReader_PEEKED_TRUE               = 5
	jsonReader_PEEKED_FALSE              = 6
	jsonReader_PEEKED_NULL               = 7
	jsonReader_PEEKED_SINGLE_QUOTED      = 8
	jsonReader_PEEKED_DOUBLE_QUOTED      = 9
	jsonReader_PEEKED_UNQUOTED           = 10
	jsonReader_PEEKED_BUFFERED           = 11
	jsonReader_PEEKED_SINGLE_QUOTED_NAME = 12
	jsonReader_PEEKED_DOUBLE_QUOTED_NAME = 13
	jsonReader_PEEKED_UNQUOTED_NAME      = 14
	jsonReader_PEEKED_LONG               = 15
	jsonReader_PEEKED_NUMBER             = 16
	jsonReader_PEEKED_EOF                = 17
)

// State machine when parsing numbers
const (
	jsonReader_NUMBER_CHAR_NONE           = 0
	jsonReader_NUMBER_CHAR_SIGN           = 1
	jsonReader_NUMBER_CHAR_DIGIT          = 2
	jsonReader_NUMBER_CHAR_DECIMAL        = 3
	jsonReader_NUMBER_CHAR_FRACTION_DIGIT = 4
	jsonReader_NUMBER_CHAR_EXP_E          = 5
	jsonReader_NUMBER_CHAR_EXP_SIGN       = 6
	jsonReader_NUMBER_CHAR_EXP_DIGIT      = 7
)

const jsonReader_BUFFER_SIZE = 1024

// JsonReader is com.google.gson.stream.JsonReader: reads a JSON encoded
// value as a stream of tokens.
type JsonReader struct {
	in JavaReader

	// True to accept non-spec compliant JSON.
	lenient bool

	// Use a manual buffer to easily read and unread upcoming characters, and
	// also so we can create strings without an intermediate StringBuilder.
	buffer []uint16
	pos    int
	limit  int

	lineNumber int
	lineStart  int

	peeked int

	// A peeked value that was composed entirely of digits with an optional
	// leading dash. Positive values may not have a leading 0.
	peekedLong int64

	// The number of characters in a peeked number literal.
	peekedNumberLength int

	// A peeked string that should be parsed on the next double, long or
	// string. This is populated before a numeric value is parsed and used if
	// that parsing fails.
	peekedString string

	// The nesting stack.
	stack     []int
	stackSize int

	// The path members. It corresponds directly to stack: At indices where
	// the stack contains an object (EMPTY_OBJECT, DANGLING_NAME or
	// NONEMPTY_OBJECT), pathNames contains the name at this scope. Where it
	// contains an array (EMPTY_ARRAY, NONEMPTY_ARRAY) pathIndices contains
	// the current index in that array. Otherwise the value is undefined, and
	// we take advantage of that by incrementing pathIndices when doing so
	// isn't useful.
	pathNames   []*string
	pathIndices []int
}

// NewJsonReader creates a new instance that reads a JSON-encoded stream
// from in.
//
// Java: JsonReader(Reader)
func NewJsonReader(in JavaReader) *JsonReader {
	r := &JsonReader{
		buffer:      make([]uint16, jsonReader_BUFFER_SIZE),
		stack:       make([]int, 32),
		pathNames:   make([]*string, 32),
		pathIndices: make([]int, 32),
	}
	r.stack[r.stackSize] = jsonScope_EMPTY_DOCUMENT
	r.stackSize++
	if in == nil {
		panic(NewNullPointerExceptionMsg("in == null"))
	}
	r.in = in
	return r
}

// SetLenient configures this parser to be liberal in what it accepts.
//
// Java: JsonReader.setLenient(boolean)
func (r *JsonReader) SetLenient(lenient bool) {
	r.lenient = lenient
}

// IsLenient returns true if this parser is liberal in what it accepts.
//
// Java: JsonReader.isLenient()
func (r *JsonReader) IsLenient() bool {
	return r.lenient
}

// BeginArray consumes the next token from the JSON stream and asserts that
// it is the beginning of a new array.
//
// Java: JsonReader.beginArray()
func (r *JsonReader) BeginArray() {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}
	if p == jsonReader_PEEKED_BEGIN_ARRAY {
		r.push(jsonScope_EMPTY_ARRAY)
		r.pathIndices[r.stackSize-1] = 0
		r.peeked = jsonReader_PEEKED_NONE
	} else {
		panic(NewIllegalStateExceptionMsg("Expected BEGIN_ARRAY but was " + r.Peek().String() + r.locationString()))
	}
}

// EndArray consumes the next token from the JSON stream and asserts that it
// is the end of the current array.
//
// Java: JsonReader.endArray()
func (r *JsonReader) EndArray() {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}
	if p == jsonReader_PEEKED_END_ARRAY {
		r.stackSize--
		r.pathIndices[r.stackSize-1]++
		r.peeked = jsonReader_PEEKED_NONE
	} else {
		panic(NewIllegalStateExceptionMsg("Expected END_ARRAY but was " + r.Peek().String() + r.locationString()))
	}
}

// BeginObject consumes the next token from the JSON stream and asserts that
// it is the beginning of a new object.
//
// Java: JsonReader.beginObject()
func (r *JsonReader) BeginObject() {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}
	if p == jsonReader_PEEKED_BEGIN_OBJECT {
		r.push(jsonScope_EMPTY_OBJECT)
		r.peeked = jsonReader_PEEKED_NONE
	} else {
		panic(NewIllegalStateExceptionMsg("Expected BEGIN_OBJECT but was " + r.Peek().String() + r.locationString()))
	}
}

// EndObject consumes the next token from the JSON stream and asserts that
// it is the end of the current object.
//
// Java: JsonReader.endObject()
func (r *JsonReader) EndObject() {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}
	if p == jsonReader_PEEKED_END_OBJECT {
		r.stackSize--
		r.pathNames[r.stackSize] = nil // Free the last path name so that it can be garbage collected!
		r.pathIndices[r.stackSize-1]++
		r.peeked = jsonReader_PEEKED_NONE
	} else {
		panic(NewIllegalStateExceptionMsg("Expected END_OBJECT but was " + r.Peek().String() + r.locationString()))
	}
}

// HasNext returns true if the current array or object has another element.
//
// Java: JsonReader.hasNext()
func (r *JsonReader) HasNext() bool {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}
	return p != jsonReader_PEEKED_END_OBJECT && p != jsonReader_PEEKED_END_ARRAY && p != jsonReader_PEEKED_EOF
}

// Peek returns the type of the next token without consuming it.
//
// Java: JsonReader.peek()
func (r *JsonReader) Peek() JsonToken {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}

	switch p {
	case jsonReader_PEEKED_BEGIN_OBJECT:
		return JsonToken_BEGIN_OBJECT
	case jsonReader_PEEKED_END_OBJECT:
		return JsonToken_END_OBJECT
	case jsonReader_PEEKED_BEGIN_ARRAY:
		return JsonToken_BEGIN_ARRAY
	case jsonReader_PEEKED_END_ARRAY:
		return JsonToken_END_ARRAY
	case jsonReader_PEEKED_SINGLE_QUOTED_NAME, jsonReader_PEEKED_DOUBLE_QUOTED_NAME, jsonReader_PEEKED_UNQUOTED_NAME:
		return JsonToken_NAME
	case jsonReader_PEEKED_TRUE, jsonReader_PEEKED_FALSE:
		return JsonToken_BOOLEAN
	case jsonReader_PEEKED_NULL:
		return JsonToken_NULL
	case jsonReader_PEEKED_SINGLE_QUOTED, jsonReader_PEEKED_DOUBLE_QUOTED, jsonReader_PEEKED_UNQUOTED, jsonReader_PEEKED_BUFFERED:
		return JsonToken_STRING
	case jsonReader_PEEKED_LONG, jsonReader_PEEKED_NUMBER:
		return JsonToken_NUMBER
	case jsonReader_PEEKED_EOF:
		return JsonToken_END_DOCUMENT
	default:
		panic(NewAssertionErrorNoMessage())
	}
}

func (r *JsonReader) doPeek() int {
	peekStack := r.stack[r.stackSize-1]
	if peekStack == jsonScope_EMPTY_ARRAY {
		r.stack[r.stackSize-1] = jsonScope_NONEMPTY_ARRAY
	} else if peekStack == jsonScope_NONEMPTY_ARRAY {
		// Look for a comma before the next element.
		c := r.nextNonWhitespace(true)
		switch c {
		case ']':
			r.peeked = jsonReader_PEEKED_END_ARRAY
			return r.peeked
		case ';':
			r.checkLenient() // fall-through
		case ',':
		default:
			panic(r.syntaxError("Unterminated array"))
		}
	} else if peekStack == jsonScope_EMPTY_OBJECT || peekStack == jsonScope_NONEMPTY_OBJECT {
		r.stack[r.stackSize-1] = jsonScope_DANGLING_NAME
		// Look for a comma before the next element.
		if peekStack == jsonScope_NONEMPTY_OBJECT {
			c := r.nextNonWhitespace(true)
			switch c {
			case '}':
				r.peeked = jsonReader_PEEKED_END_OBJECT
				return r.peeked
			case ';':
				r.checkLenient() // fall-through
			case ',':
			default:
				panic(r.syntaxError("Unterminated object"))
			}
		}
		c := r.nextNonWhitespace(true)
		switch c {
		case '"':
			r.peeked = jsonReader_PEEKED_DOUBLE_QUOTED_NAME
			return r.peeked
		case '\'':
			r.checkLenient()
			r.peeked = jsonReader_PEEKED_SINGLE_QUOTED_NAME
			return r.peeked
		case '}':
			if peekStack != jsonScope_NONEMPTY_OBJECT {
				r.peeked = jsonReader_PEEKED_END_OBJECT
				return r.peeked
			} else {
				panic(r.syntaxError("Expected name"))
			}
		default:
			r.checkLenient()
			r.pos-- // Don't consume the first character in an unquoted string.
			if r.isLiteral(uint16(c)) {
				r.peeked = jsonReader_PEEKED_UNQUOTED_NAME
				return r.peeked
			} else {
				panic(r.syntaxError("Expected name"))
			}
		}
	} else if peekStack == jsonScope_DANGLING_NAME {
		r.stack[r.stackSize-1] = jsonScope_NONEMPTY_OBJECT
		// Look for a colon before the value.
		c := r.nextNonWhitespace(true)
		switch c {
		case ':':
		case '=':
			r.checkLenient()
			if (r.pos < r.limit || r.fillBuffer(1)) && r.buffer[r.pos] == '>' {
				r.pos++
			}
		default:
			panic(r.syntaxError("Expected ':'"))
		}
	} else if peekStack == jsonScope_EMPTY_DOCUMENT {
		if r.lenient {
			r.consumeNonExecutePrefix()
		}
		r.stack[r.stackSize-1] = jsonScope_NONEMPTY_DOCUMENT
	} else if peekStack == jsonScope_NONEMPTY_DOCUMENT {
		c := r.nextNonWhitespace(false)
		if c == -1 {
			r.peeked = jsonReader_PEEKED_EOF
			return r.peeked
		} else {
			r.checkLenient()
			r.pos--
		}
	} else if peekStack == jsonScope_CLOSED {
		panic(NewIllegalStateExceptionMsg("JsonReader is closed"))
	}

	c := r.nextNonWhitespace(true)
	switch c {
	case ']':
		if peekStack == jsonScope_EMPTY_ARRAY {
			r.peeked = jsonReader_PEEKED_END_ARRAY
			return r.peeked
		}
		// fall-through to handle ",]"
		fallthrough
	case ';', ',':
		// In lenient mode, a 0-length literal in an array means 'null'.
		if peekStack == jsonScope_EMPTY_ARRAY || peekStack == jsonScope_NONEMPTY_ARRAY {
			r.checkLenient()
			r.pos--
			r.peeked = jsonReader_PEEKED_NULL
			return r.peeked
		} else {
			panic(r.syntaxError("Unexpected value"))
		}
	case '\'':
		r.checkLenient()
		r.peeked = jsonReader_PEEKED_SINGLE_QUOTED
		return r.peeked
	case '"':
		r.peeked = jsonReader_PEEKED_DOUBLE_QUOTED
		return r.peeked
	case '[':
		r.peeked = jsonReader_PEEKED_BEGIN_ARRAY
		return r.peeked
	case '{':
		r.peeked = jsonReader_PEEKED_BEGIN_OBJECT
		return r.peeked
	default:
		r.pos-- // Don't consume the first character in a literal value.
	}

	result := r.peekKeyword()
	if result != jsonReader_PEEKED_NONE {
		return result
	}

	result = r.peekNumber()
	if result != jsonReader_PEEKED_NONE {
		return result
	}

	if !r.isLiteral(r.buffer[r.pos]) {
		panic(r.syntaxError("Expected value"))
	}

	r.checkLenient()
	r.peeked = jsonReader_PEEKED_UNQUOTED
	return r.peeked
}

func (r *JsonReader) peekKeyword() int {
	// Figure out which keyword we're matching against by its first character.
	c := r.buffer[r.pos]
	var keyword, keywordUpper string
	var peeking int
	if c == 't' || c == 'T' {
		keyword = "true"
		keywordUpper = "TRUE"
		peeking = jsonReader_PEEKED_TRUE
	} else if c == 'f' || c == 'F' {
		keyword = "false"
		keywordUpper = "FALSE"
		peeking = jsonReader_PEEKED_FALSE
	} else if c == 'n' || c == 'N' {
		keyword = "null"
		keywordUpper = "NULL"
		peeking = jsonReader_PEEKED_NULL
	} else {
		return jsonReader_PEEKED_NONE
	}

	// Confirm that chars [1..length) match the keyword.
	length := len(keyword)
	for i := 1; i < length; i++ {
		if r.pos+i >= r.limit && !r.fillBuffer(i+1) {
			return jsonReader_PEEKED_NONE
		}
		c = r.buffer[r.pos+i]
		if c != uint16(keyword[i]) && c != uint16(keywordUpper[i]) {
			return jsonReader_PEEKED_NONE
		}
	}

	if (r.pos+length < r.limit || r.fillBuffer(length+1)) && r.isLiteral(r.buffer[r.pos+length]) {
		return jsonReader_PEEKED_NONE // Don't match trues, falsey or nullsoft!
	}

	// We've found the keyword followed either by EOF or by a non-literal
	// character.
	r.pos += length
	r.peeked = peeking
	return r.peeked
}

func (r *JsonReader) peekNumber() int {
	// Like nextNonWhitespace, this uses locals 'p' and 'l' to save
	// inner-loop field access.
	buffer := r.buffer
	p := r.pos
	l := r.limit

	var value int64 // Negative to accommodate Long.MIN_VALUE more easily.
	negative := false
	fitsInLong := true
	last := jsonReader_NUMBER_CHAR_NONE

	i := 0

charactersOfNumber:
	for ; true; i++ {
		if p+i == l {
			if i == len(buffer) {
				// Though this looks like a well-formed number, it's too long
				// to continue reading. Give up and let the application handle
				// this as an unquoted literal.
				return jsonReader_PEEKED_NONE
			}
			if !r.fillBuffer(i + 1) {
				break
			}
			p = r.pos
			l = r.limit
		}

		c := buffer[p+i]
		switch c {
		case '-':
			if last == jsonReader_NUMBER_CHAR_NONE {
				negative = true
				last = jsonReader_NUMBER_CHAR_SIGN
				continue
			} else if last == jsonReader_NUMBER_CHAR_EXP_E {
				last = jsonReader_NUMBER_CHAR_EXP_SIGN
				continue
			}
			return jsonReader_PEEKED_NONE

		case '+':
			if last == jsonReader_NUMBER_CHAR_EXP_E {
				last = jsonReader_NUMBER_CHAR_EXP_SIGN
				continue
			}
			return jsonReader_PEEKED_NONE

		case 'e', 'E':
			if last == jsonReader_NUMBER_CHAR_DIGIT || last == jsonReader_NUMBER_CHAR_FRACTION_DIGIT {
				last = jsonReader_NUMBER_CHAR_EXP_E
				continue
			}
			return jsonReader_PEEKED_NONE

		case '.':
			if last == jsonReader_NUMBER_CHAR_DIGIT {
				last = jsonReader_NUMBER_CHAR_DECIMAL
				continue
			}
			return jsonReader_PEEKED_NONE

		default:
			if c < '0' || c > '9' {
				if !r.isLiteral(c) {
					break charactersOfNumber
				}
				return jsonReader_PEEKED_NONE
			}
			if last == jsonReader_NUMBER_CHAR_SIGN || last == jsonReader_NUMBER_CHAR_NONE {
				value = -int64(c - '0')
				last = jsonReader_NUMBER_CHAR_DIGIT
			} else if last == jsonReader_NUMBER_CHAR_DIGIT {
				if value == 0 {
					return jsonReader_PEEKED_NONE // Leading '0' prefix is not allowed (since it could be octal).
				}
				newValue := value*10 - int64(c-'0')
				fitsInLong = fitsInLong && (value > jsonReader_MIN_INCOMPLETE_INTEGER ||
					(value == jsonReader_MIN_INCOMPLETE_INTEGER && newValue < value))
				value = newValue
			} else if last == jsonReader_NUMBER_CHAR_DECIMAL {
				last = jsonReader_NUMBER_CHAR_FRACTION_DIGIT
			} else if last == jsonReader_NUMBER_CHAR_EXP_E || last == jsonReader_NUMBER_CHAR_EXP_SIGN {
				last = jsonReader_NUMBER_CHAR_EXP_DIGIT
			}
		}
	}

	// We've read a complete number. Decide if it's a PEEKED_LONG or a
	// PEEKED_NUMBER.
	if last == jsonReader_NUMBER_CHAR_DIGIT && fitsInLong && (value != -0x8000000000000000 || negative) && (value != 0 || !negative) {
		if negative {
			r.peekedLong = value
		} else {
			r.peekedLong = -value
		}
		r.pos += i
		r.peeked = jsonReader_PEEKED_LONG
		return r.peeked
	} else if last == jsonReader_NUMBER_CHAR_DIGIT || last == jsonReader_NUMBER_CHAR_FRACTION_DIGIT ||
		last == jsonReader_NUMBER_CHAR_EXP_DIGIT {
		r.peekedNumberLength = i
		r.peeked = jsonReader_PEEKED_NUMBER
		return r.peeked
	} else {
		return jsonReader_PEEKED_NONE
	}
}

func (r *JsonReader) isLiteral(c uint16) bool {
	switch c {
	case '/', '\\', ';', '#', '=':
		r.checkLenient() // fall-through
		return false
	case '{', '}', '[', ']', ':', ',', ' ', '\t', '\f', '\r', '\n':
		return false
	default:
		return true
	}
}

// NextName returns the next token, a property name, and consumes it.
//
// Java: JsonReader.nextName()
func (r *JsonReader) NextName() string {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}
	var result string
	if p == jsonReader_PEEKED_UNQUOTED_NAME {
		result = r.nextUnquotedValue()
	} else if p == jsonReader_PEEKED_SINGLE_QUOTED_NAME {
		result = r.nextQuotedValue('\'')
	} else if p == jsonReader_PEEKED_DOUBLE_QUOTED_NAME {
		result = r.nextQuotedValue('"')
	} else {
		panic(NewIllegalStateExceptionMsg("Expected a name but was " + r.Peek().String() + r.locationString()))
	}
	r.peeked = jsonReader_PEEKED_NONE
	r.pathNames[r.stackSize-1] = &result
	return result
}

// NextString returns the string value of the next token, consuming it. If
// the next token is a number, this method will return its string form.
//
// Java: JsonReader.nextString()
func (r *JsonReader) NextString() string {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}
	var result string
	if p == jsonReader_PEEKED_UNQUOTED {
		result = r.nextUnquotedValue()
	} else if p == jsonReader_PEEKED_SINGLE_QUOTED {
		result = r.nextQuotedValue('\'')
	} else if p == jsonReader_PEEKED_DOUBLE_QUOTED {
		result = r.nextQuotedValue('"')
	} else if p == jsonReader_PEEKED_BUFFERED {
		result = r.peekedString
		r.peekedString = ""
	} else if p == jsonReader_PEEKED_LONG {
		result = LongToString(r.peekedLong)
	} else if p == jsonReader_PEEKED_NUMBER {
		result = StringFromUTF16(r.buffer[r.pos : r.pos+r.peekedNumberLength])
		r.pos += r.peekedNumberLength
	} else {
		panic(NewIllegalStateExceptionMsg("Expected a string but was " + r.Peek().String() + r.locationString()))
	}
	r.peeked = jsonReader_PEEKED_NONE
	r.pathIndices[r.stackSize-1]++
	return result
}

// NextBoolean returns the boolean value of the next token, consuming it.
//
// Java: JsonReader.nextBoolean()
func (r *JsonReader) NextBoolean() bool {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}
	if p == jsonReader_PEEKED_TRUE {
		r.peeked = jsonReader_PEEKED_NONE
		r.pathIndices[r.stackSize-1]++
		return true
	} else if p == jsonReader_PEEKED_FALSE {
		r.peeked = jsonReader_PEEKED_NONE
		r.pathIndices[r.stackSize-1]++
		return false
	}
	panic(NewIllegalStateExceptionMsg("Expected a boolean but was " + r.Peek().String() + r.locationString()))
}

// NextNull consumes the next token from the JSON stream and asserts that it
// is a literal null.
//
// Java: JsonReader.nextNull()
func (r *JsonReader) NextNull() {
	p := r.peeked
	if p == jsonReader_PEEKED_NONE {
		p = r.doPeek()
	}
	if p == jsonReader_PEEKED_NULL {
		r.peeked = jsonReader_PEEKED_NONE
		r.pathIndices[r.stackSize-1]++
	} else {
		panic(NewIllegalStateExceptionMsg("Expected null but was " + r.Peek().String() + r.locationString()))
	}
}

// nextQuotedValue returns the string up to but not including quote,
// unescaping any character escape sequences encountered along the way. The
// opening quote should have already been read. This consumes the closing
// quote, but does not include it in the returned string.
func (r *JsonReader) nextQuotedValue(quote uint16) string {
	// Like nextNonWhitespace, this uses locals 'p' and 'l' to save
	// inner-loop field access.
	buffer := r.buffer
	var builder []uint16
	haveBuilder := false
	for {
		p := r.pos
		l := r.limit
		// the index of the first character not yet appended to the builder.
		start := p
		for p < l {
			c := buffer[p]
			p++

			if c == quote {
				r.pos = p
				length := p - start - 1
				if !haveBuilder {
					return StringFromUTF16(buffer[start : start+length])
				} else {
					builder = append(builder, buffer[start:start+length]...)
					return StringFromUTF16(builder)
				}
			} else if c == '\\' {
				r.pos = p
				length := p - start - 1
				if !haveBuilder {
					estimatedLength := (length + 1) * 2
					builder = make([]uint16, 0, max(estimatedLength, 16))
					haveBuilder = true
				}
				builder = append(builder, buffer[start:start+length]...)
				builder = append(builder, r.readEscapeCharacter())
				p = r.pos
				l = r.limit
				start = p
			} else if c == '\n' {
				r.lineNumber++
				r.lineStart = p
			}
		}

		if !haveBuilder {
			estimatedLength := (p - start) * 2
			builder = make([]uint16, 0, max(estimatedLength, 16))
			haveBuilder = true
		}
		builder = append(builder, buffer[start:p]...)
		r.pos = p
		if !r.fillBuffer(1) {
			panic(r.syntaxError("Unterminated string"))
		}
	}
}

// nextUnquotedValue returns an unquoted value as a string.
func (r *JsonReader) nextUnquotedValue() string {
	var builder []uint16
	haveBuilder := false
	i := 0

findNonLiteralCharacter:
	for {
		for ; r.pos+i < r.limit; i++ {
			switch r.buffer[r.pos+i] {
			case '/', '\\', ';', '#', '=':
				r.checkLenient() // fall-through
				break findNonLiteralCharacter
			case '{', '}', '[', ']', ':', ',', ' ', '\t', '\f', '\r', '\n':
				break findNonLiteralCharacter
			}
		}

		// Attempt to load the entire literal into the buffer at once.
		if i < len(r.buffer) {
			if r.fillBuffer(i + 1) {
				continue
			} else {
				break
			}
		}

		// use a StringBuilder when the value is too long. This is too long to
		// be a number!
		if !haveBuilder {
			builder = make([]uint16, 0, max(i, 16))
			haveBuilder = true
		}
		builder = append(builder, r.buffer[r.pos:r.pos+i]...)
		r.pos += i
		i = 0
		if !r.fillBuffer(1) {
			break
		}
	}

	var result string
	if !haveBuilder {
		result = StringFromUTF16(r.buffer[r.pos : r.pos+i])
	} else {
		builder = append(builder, r.buffer[r.pos:r.pos+i]...)
		result = StringFromUTF16(builder)
	}
	r.pos += i
	return result
}

func (r *JsonReader) skipQuotedValue(quote uint16) {
	// Like nextNonWhitespace, this uses locals 'p' and 'l' to save
	// inner-loop field access.
	buffer := r.buffer
	for {
		p := r.pos
		l := r.limit
		// the index of the first character not yet appended to the builder.
		for p < l {
			c := buffer[p]
			p++
			if c == quote {
				r.pos = p
				return
			} else if c == '\\' {
				r.pos = p
				r.readEscapeCharacter()
				p = r.pos
				l = r.limit
			} else if c == '\n' {
				r.lineNumber++
				r.lineStart = p
			}
		}
		r.pos = p
		if !r.fillBuffer(1) {
			break
		}
	}
	panic(r.syntaxError("Unterminated string"))
}

func (r *JsonReader) skipUnquotedValue() {
	for {
		i := 0
		for ; r.pos+i < r.limit; i++ {
			switch r.buffer[r.pos+i] {
			case '/', '\\', ';', '#', '=':
				r.checkLenient() // fall-through
				r.pos += i
				return
			case '{', '}', '[', ']', ':', ',', ' ', '\t', '\f', '\r', '\n':
				r.pos += i
				return
			}
		}
		r.pos += i
		if !r.fillBuffer(1) {
			break
		}
	}
}

// Close closes this JSON reader and the underlying Reader.
//
// Java: JsonReader.close()
func (r *JsonReader) Close() {
	r.peeked = jsonReader_PEEKED_NONE
	r.stack[0] = jsonScope_CLOSED
	r.stackSize = 1
	r.in.Close()
}

// SkipValue skips the next value recursively. If it is an object or array,
// all nested elements are skipped.
//
// Java: JsonReader.skipValue()
func (r *JsonReader) SkipValue() {
	count := 0
	for {
		p := r.peeked
		if p == jsonReader_PEEKED_NONE {
			p = r.doPeek()
		}

		if p == jsonReader_PEEKED_BEGIN_ARRAY {
			r.push(jsonScope_EMPTY_ARRAY)
			count++
		} else if p == jsonReader_PEEKED_BEGIN_OBJECT {
			r.push(jsonScope_EMPTY_OBJECT)
			count++
		} else if p == jsonReader_PEEKED_END_ARRAY {
			r.stackSize--
			count--
		} else if p == jsonReader_PEEKED_END_OBJECT {
			r.stackSize--
			count--
		} else if p == jsonReader_PEEKED_UNQUOTED_NAME || p == jsonReader_PEEKED_UNQUOTED {
			r.skipUnquotedValue()
		} else if p == jsonReader_PEEKED_SINGLE_QUOTED || p == jsonReader_PEEKED_SINGLE_QUOTED_NAME {
			r.skipQuotedValue('\'')
		} else if p == jsonReader_PEEKED_DOUBLE_QUOTED || p == jsonReader_PEEKED_DOUBLE_QUOTED_NAME {
			r.skipQuotedValue('"')
		} else if p == jsonReader_PEEKED_NUMBER {
			r.pos += r.peekedNumberLength
		}
		r.peeked = jsonReader_PEEKED_NONE
		if count == 0 {
			break
		}
	}

	r.pathIndices[r.stackSize-1]++
	null := "null"
	r.pathNames[r.stackSize-1] = &null
}

func (r *JsonReader) push(newTop int) {
	if r.stackSize == len(r.stack) {
		newLength := r.stackSize * 2
		r.stack = append(r.stack, make([]int, newLength-len(r.stack))...)
		r.pathIndices = append(r.pathIndices, make([]int, newLength-len(r.pathIndices))...)
		r.pathNames = append(r.pathNames, make([]*string, newLength-len(r.pathNames))...)
	}
	r.stack[r.stackSize] = newTop
	r.stackSize++
}

// fillBuffer returns true once limit - pos >= minimum. If the data is
// exhausted before that many characters are available, this returns false.
func (r *JsonReader) fillBuffer(minimum int) bool {
	buffer := r.buffer
	r.lineStart -= r.pos
	if r.limit != r.pos {
		r.limit -= r.pos
		copy(buffer[0:r.limit], buffer[r.pos:r.pos+r.limit])
	} else {
		r.limit = 0
	}

	r.pos = 0
	var total int
	for {
		total = r.in.Read(buffer, r.limit, len(buffer)-r.limit)
		if total == -1 {
			break
		}
		r.limit += total

		// if this is the first read, consume an optional byte order mark
		// (BOM) if it exists
		if r.lineNumber == 0 && r.lineStart == 0 && r.limit > 0 && buffer[0] == 0xfeff {
			r.pos++
			r.lineStart++
			minimum++
		}

		if r.limit >= minimum {
			return true
		}
	}
	return false
}

// nextNonWhitespace returns the next character in the stream that is
// neither whitespace nor a part of a comment. When this returns, the
// returned character is always at buffer[pos-1]; this means the caller can
// always push back the returned character by decrementing pos.
func (r *JsonReader) nextNonWhitespace(throwOnEof bool) int {
	// This code uses ugly local variables 'p' and 'l' representing the 'pos'
	// and 'limit' fields respectively. Using locals rather than fields saves
	// a few field reads for each whitespace character in a pretty-printed
	// document, resulting in a 5% speedup. We need to flush 'p' to its field
	// before any (potentially indirect) call to fillBuffer() and reread both
	// 'p' and 'l' after any (potentially indirect) call to the same method.
	buffer := r.buffer
	p := r.pos
	l := r.limit
	for {
		if p == l {
			r.pos = p
			if !r.fillBuffer(1) {
				break
			}
			p = r.pos
			l = r.limit
		}

		c := int(buffer[p])
		p++
		if c == '\n' {
			r.lineNumber++
			r.lineStart = p
			continue
		} else if c == ' ' || c == '\r' || c == '\t' {
			continue
		}

		if c == '/' {
			r.pos = p
			if p == l {
				r.pos-- // push back '/' so it's still in the buffer when this method returns
				charsLoaded := r.fillBuffer(2)
				r.pos++ // consume the '/' again
				if !charsLoaded {
					return c
				}
			}

			r.checkLenient()
			peek := buffer[r.pos]
			switch peek {
			case '*':
				// skip a /* c-style comment */
				r.pos++
				if !r.skipTo("*/") {
					panic(r.syntaxError("Unterminated comment"))
				}
				p = r.pos + 2
				l = r.limit
				continue

			case '/':
				// skip a // end-of-line comment
				r.pos++
				r.skipToEndOfLine()
				p = r.pos
				l = r.limit
				continue

			default:
				return c
			}
		} else if c == '#' {
			r.pos = p
			// Skip a # hash end-of-line comment. The JSON RFC doesn't specify
			// this behaviour, but it's required to parse existing documents.
			// See http://b/2571423.
			r.checkLenient()
			r.skipToEndOfLine()
			p = r.pos
			l = r.limit
		} else {
			r.pos = p
			return c
		}
	}
	if throwOnEof {
		panic(NewEOFException("End of input" + r.locationString()))
	} else {
		return -1
	}
}

func (r *JsonReader) checkLenient() {
	if !r.lenient {
		panic(r.syntaxError("Use JsonReader.setLenient(true) to accept malformed JSON"))
	}
}

// skipToEndOfLine advances the position until after the next newline
// character. If the line is terminated by "\r\n", the '\n' must be consumed
// as whitespace by the caller.
func (r *JsonReader) skipToEndOfLine() {
	for r.pos < r.limit || r.fillBuffer(1) {
		c := r.buffer[r.pos]
		r.pos++
		if c == '\n' {
			r.lineNumber++
			r.lineStart = r.pos
			break
		} else if c == '\r' {
			break
		}
	}
}

// skipTo advances to the beginning of toFind (ASCII).
func (r *JsonReader) skipTo(toFind string) bool {
	length := len(toFind)
outer:
	for ; r.pos+length <= r.limit || r.fillBuffer(length); r.pos++ {
		if r.buffer[r.pos] == '\n' {
			r.lineNumber++
			r.lineStart = r.pos + 1
			continue
		}
		for c := 0; c < length; c++ {
			if r.buffer[r.pos+c] != uint16(toFind[c]) {
				continue outer
			}
		}
		return true
	}
	return false
}

// Java: JsonReader.toString()
func (r *JsonReader) String() string {
	return "JsonReader" + r.locationString()
}

func (r *JsonReader) locationString() string {
	line := r.lineNumber + 1
	column := r.pos - r.lineStart + 1
	return " at line " + itoa(line) + " column " + itoa(column) + " path " + r.GetPath()
}

func (r *JsonReader) getPath(usePreviousPath bool) string {
	result := NewStringBuilder().AppendChar('$')
	for i := 0; i < r.stackSize; i++ {
		switch r.stack[i] {
		case jsonScope_EMPTY_ARRAY, jsonScope_NONEMPTY_ARRAY:
			pathIndex := r.pathIndices[i]
			// If index is last path element it points to next array element;
			// have to decrement
			if usePreviousPath && pathIndex > 0 && i == r.stackSize-1 {
				pathIndex--
			}
			result.AppendChar('[').AppendInt(int64(pathIndex)).AppendChar(']')
		case jsonScope_EMPTY_OBJECT, jsonScope_DANGLING_NAME, jsonScope_NONEMPTY_OBJECT:
			result.AppendChar('.')
			if r.pathNames[i] != nil {
				result.Append(*r.pathNames[i])
			}
		case jsonScope_NONEMPTY_DOCUMENT, jsonScope_EMPTY_DOCUMENT, jsonScope_CLOSED:
		}
	}
	return result.String()
}

// GetPreviousPath returns a JsonPath to the element just read.
//
// Java: JsonReader.getPreviousPath()
func (r *JsonReader) GetPreviousPath() string {
	return r.getPath(true)
}

// GetPath returns a JsonPath to the next (to be read) element.
//
// Java: JsonReader.getPath()
func (r *JsonReader) GetPath() string {
	return r.getPath(false)
}

// readEscapeCharacter unescapes the character identified by the
// character or characters that immediately follow a backslash. The backslash
// '\' should have already been read. This supports both unicode escapes "u000A"
// and two-character escapes "\n".
func (r *JsonReader) readEscapeCharacter() uint16 {
	if r.pos == r.limit && !r.fillBuffer(1) {
		panic(r.syntaxError("Unterminated escape sequence"))
	}

	escaped := r.buffer[r.pos]
	r.pos++
	switch escaped {
	case 'u':
		if r.pos+4 > r.limit && !r.fillBuffer(4) {
			panic(r.syntaxError("Unterminated escape sequence"))
		}
		// Equivalent to Integer.parseInt(stringPool.get(buffer, pos, 4), 16);
		var result uint16
		for i, end := r.pos, r.pos+4; i < end; i++ {
			c := r.buffer[i]
			result <<= 4
			if c >= '0' && c <= '9' {
				result += c - '0'
			} else if c >= 'a' && c <= 'f' {
				result += c - 'a' + 10
			} else if c >= 'A' && c <= 'F' {
				result += c - 'A' + 10
			} else {
				panic(NewNumberFormatException("\\u" + StringFromUTF16(r.buffer[r.pos:r.pos+4])))
			}
		}
		r.pos += 4
		return result

	case 't':
		return '\t'

	case 'b':
		return '\b'

	case 'n':
		return '\n'

	case 'r':
		return '\r'

	case 'f':
		return '\f'

	case '\n':
		r.lineNumber++
		r.lineStart = r.pos
		// fall-through
		return escaped

	case '\'', '"', '\\', '/':
		return escaped
	default:
		// throw error when none of the above cases are matched
		panic(r.syntaxError("Invalid escape sequence"))
	}
}

// syntaxError builds the MalformedJsonException with the given message and
// a context snippet with this reader's content; the caller throws it.
func (r *JsonReader) syntaxError(message string) *MalformedJsonException {
	return NewMalformedJsonException(message + r.locationString())
}

// consumeNonExecutePrefix consumes the non-execute prefix if it exists.
func (r *JsonReader) consumeNonExecutePrefix() {
	// fast forward through the leading whitespace
	r.nextNonWhitespace(true)
	r.pos--

	if r.pos+5 > r.limit && !r.fillBuffer(5) {
		return
	}

	p := r.pos
	buf := r.buffer
	if buf[p] != ')' || buf[p+1] != ']' || buf[p+2] != '}' || buf[p+3] != '\'' || buf[p+4] != '\n' {
		return // not a security token!
	}

	// we consumed a security token!
	r.pos += 5
}
