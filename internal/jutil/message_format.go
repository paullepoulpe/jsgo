package jutil

import (
	"errors"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// Go-only: language-adapter — java.text.MessageFormat for the default
// locale en_US (see the Javadoc of java.text.MessageFormat, ChoiceFormat,
// NumberFormat and SimpleDateFormat).
//
// A pattern is literal text with format elements { index [, type [, style]] }.
// Outside elements a quote starts or ends quoted text and two quotes stand
// for one; inside elements quotes are kept (the style is handed on as it is
// written) but braces and commas within them do not count. Arguments are
// rendered as the Javadoc table of MessageFormat.format says: "{n}" when
// there is no argument n, "null" for null, a subformat when the element has
// a type, and otherwise numbers through NumberFormat.getInstance() and
// other objects through toString().
//
// Observed behaviour (black-box, JDK 21, en_US):
//   - errors are IllegalArgumentExceptions, thrown while the pattern is read
//     (so the first bad element wins): "can't parse argument number: <s>"
//     (Integer.parseInt rules, no trimming: Unicode decimal digits and a
//     sign are accepted), "negative argument number: <n>", "<n> exceeds the
//     ArgumentIndex implementation limit" (n >= 10000), "unknown format
//     type: <type>", "Choice Pattern incorrect: <style>" and, at the end,
//     "Unmatched braces in the pattern." (but a pattern that ends inside
//     nested braces of an element just drops that element);
//   - inside an element, a ',' starts the next part (index, type, style)
//     even inside nested braces; after the style it is text. Spaces at the
//     start of the type are dropped; the type and the named styles are then
//     compared after trimming and ignoring ASCII case. Any other style is
//     passed on untrimmed;
//   - number elements need a Number ("Cannot format given Object as a
//     Number"), date and time elements a Date or a Number ("Cannot format
//     given Object as a Date"); null is "null" whatever the type;
//   - a ChoiceFormat result that contains '{' is formatted again as a
//     pattern with the same arguments;
//   - NumberFormat: "#,##0.###" (no type or style), "#,##0" (integer),
//     "#,##0%" (percent) and "$#,##0.00" (currency), rounding half-even. A
//     double is written from its shortest decimal digits (padded with
//     zeros), and a tie in those digits is settled by the exact binary
//     value, except that a lone digit 5 at the third decimal place or
//     further right is rounded down to 0 (0.0005 gives "0"). Integral
//     doubles from 2^53 to 2^63 are written differently (see mfLongDigits).
//     Percent multiplies a double by 100 in double arithmetic. NaN is
//     "NaN", infinity "∞", and -0.0 keeps its sign ("-0");
//   - SimpleDateFormat patterns may use only the ASCII letters
//     GyYMLwWDdFEuaHkKhmsSzZX outside quotes ("Illegal pattern character
//     'c'").
//
// Not reproduced:
//   - there is no java.util.Date in Go, and formatting a number through a
//     date or time element, or through a DecimalFormat pattern given as the
//     style, throws UnsupportedOperationException. DecimalFormat patterns
//     are not checked when the format is built (the JDK rejects malformed
//     ones with IllegalArgumentException);
//   - for a few doubles of 2^63 or more (about 1 in 25000 random ones) the
//     JDK writes digits other than the shortest ones, e.g.
//     1.9546806789078286E25 as 19,546,806,789,078,285,000,000,000.

// MessageFormat is java.text.MessageFormat.
type MessageFormat struct {
	pattern  []uint16
	elements []mfElement
	tail     string
}

// mfElement is the literal text before a format element and the element.
type mfElement struct {
	prefix string
	index  int
	kind   mfKind
	style  mfNumberStyle // kind == mfNumber
	custom string        // a style that is not a keyword
	choice *mfChoice     // kind == mfChoiceKind
}

type mfKind int

const (
	mfNone mfKind = iota
	mfNumber
	mfDate
	mfTime
	mfChoiceKind
)

type mfNumberStyle int

const (
	mfNumberDefault mfNumberStyle = iota
	mfNumberInteger
	mfNumberPercent
	mfNumberCurrency
	mfNumberPattern
)

// mfArgumentIndexLimit is the bound the JDK puts on argument numbers.
const mfArgumentIndexLimit = 10000

// NewMessageFormat is new MessageFormat(pattern).
func NewMessageFormat(pattern string) *MessageFormat {
	m := &MessageFormat{}
	m.ApplyPattern(pattern)
	return m
}

// MessageFormat_Format is MessageFormat.format(String, Object...).
func MessageFormat_Format(pattern string, arguments []any) string {
	return NewMessageFormat(pattern).Format(arguments)
}

// ApplyPattern is MessageFormat.applyPattern(String). The pattern is read
// in two alternating stages: mfLiteralText reads text up to the next
// element, mfElementText reads that element; each element is checked and
// built as soon as it is read, so the first bad one decides the exception.
func (m *MessageFormat) ApplyPattern(pattern string) {
	src := StringToUTF16(pattern)
	var elements []mfElement
	var tail []uint16
	for at := 0; ; {
		text, open := mfLiteralText(src, at)
		if open < 0 {
			tail = text
			break
		}
		body, next, nested := mfElementText(src, open)
		if next < 0 {
			if !nested {
				panic(NewIllegalArgumentExceptionMsg("Unmatched braces in the pattern."))
			}
			// Observed (black-box): when the pattern ends inside nested
			// braces of an element, that element is silently dropped and
			// the text before it becomes the tail.
			tail = text
			break
		}
		index, kind, style := mfElementParts(body)
		elements = append(elements, mfMakeElement([4][]uint16{text, index, kind, style}))
		at = next
	}
	m.pattern = src
	m.elements = elements
	m.tail = StringFromUTF16(tail)
}

// mfLiteralText reads the literal text of src from offset at: "”" stands
// for one quote, a lone quote opens or closes a quoted run (and is not part
// of the text), and outside quoted runs '{' starts an element. It returns
// the text and the offset just after that '{', or -1 when src ended first.
func mfLiteralText(src []uint16, at int) ([]uint16, int) {
	var text []uint16
	quoted := false
	for i := at; i < len(src); i++ {
		switch c := src[i]; {
		case c == '\'' && i+1 < len(src) && src[i+1] == '\'':
			text = append(text, c)
			i++
		case c == '\'':
			quoted = !quoted
		case c == '{' && !quoted:
			return text, i + 1
		default:
			text = append(text, c)
		}
	}
	return text, -1
}

// mfElementText reads an element whose '{' is just before offset at: its
// raw text, quotes included, up to the '}' that closes it. Braces in quoted
// runs do not count, other braces nest. It returns the text and the offset
// after the closing '}', or next = -1 when src ends first; nested then says
// whether it ended inside nested braces.
func mfElementText(src []uint16, at int) (body []uint16, next int, nested bool) {
	level := 0
	quoted := false
	for i := at; i < len(src); i++ {
		c := src[i]
		switch {
		case quoted:
			quoted = c != '\''
		case c == '\'':
			quoted = true
		case c == '{':
			level++
		case c == '}' && level == 0:
			return src[at:i], i + 1, false
		case c == '}':
			level--
		}
	}
	return nil, -1, level > 0
}

// mfElementParts splits the raw text of an element into its index, type
// and style at its first two commas outside quoted runs (header: they count
// even inside nested braces; later commas belong to the style), and drops
// the spaces that start the type.
func mfElementParts(body []uint16) (index, kind, style []uint16) {
	var cuts [2]int
	found := 0
	quoted := false
	for i, c := range body {
		switch {
		case quoted:
			quoted = c != '\''
		case c == '\'':
			quoted = true
		case c == ',' && found < 2:
			cuts[found] = i
			found++
		}
	}
	switch found {
	case 0:
		return body, nil, nil
	case 1:
		index, kind = body[:cuts[0]], body[cuts[0]+1:]
	default:
		index, kind, style = body[:cuts[0]], body[cuts[0]+1:cuts[1]], body[cuts[1]+1:]
	}
	for len(kind) > 0 && kind[0] == ' ' {
		kind = kind[1:]
	}
	return index, kind, style
}

// mfMakeElement builds the element read into parts (prefix, index, type,
// style).
func mfMakeElement(parts [4][]uint16) mfElement {
	e := mfElement{prefix: StringFromUTF16(parts[0])}
	indexText := StringFromUTF16(parts[1])
	index, ok := mfParseInt(parts[1])
	if !ok {
		panic(NewIllegalArgumentExceptionMsg("can't parse argument number: " + indexText))
	}
	if index < 0 {
		panic(NewIllegalArgumentExceptionMsg("negative argument number: " + strconv.Itoa(index)))
	}
	if index >= mfArgumentIndexLimit {
		panic(NewIllegalArgumentExceptionMsg(strconv.Itoa(index) + " exceeds the ArgumentIndex implementation limit"))
	}
	e.index = index
	style := StringFromUTF16(parts[3])
	switch mfKeyword(parts[2]) {
	case "":
	case "number":
		e.kind = mfNumber
		switch mfKeyword(parts[3]) {
		case "":
			e.style = mfNumberDefault
		case "integer":
			e.style = mfNumberInteger
		case "percent":
			e.style = mfNumberPercent
		case "currency":
			e.style = mfNumberCurrency
		default:
			e.style = mfNumberPattern
			e.custom = style
		}
	case "date", "time":
		e.kind = mfDate
		if mfKeyword(parts[2]) == "time" {
			e.kind = mfTime
		}
		switch mfKeyword(parts[3]) {
		case "", "short", "medium", "long", "full":
		default:
			mfCheckDatePattern(parts[3])
			e.custom = style
		}
	case "choice":
		e.kind = mfChoiceKind
		func() {
			defer func() {
				if r := recover(); r != nil {
					panic(NewIllegalArgumentExceptionMsg("Choice Pattern incorrect: " + style))
				}
			}()
			e.choice = mfParseChoice(parts[3])
		}()
	default:
		panic(NewIllegalArgumentExceptionMsg("unknown format type: " + StringFromUTF16(parts[2])))
	}
	return e
}

// mfKeyword trims s like String.trim() and lower-cases ASCII letters; the
// keywords are ASCII, and no other character's lower case is one of their
// letters.
func mfKeyword(s []uint16) string {
	start, end := 0, len(s)
	for start < end && s[start] <= ' ' {
		start++
	}
	for end > start && s[end-1] <= ' ' {
		end--
	}
	var sb strings.Builder
	for _, c := range s[start:end] {
		if c >= 0x80 {
			return "\x00" // not a keyword
		}
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		sb.WriteByte(byte(c))
	}
	return sb.String()
}

// mfParseInt is Integer.parseInt(s) (radix 10), reporting failure instead
// of throwing.
func mfParseInt(s []uint16) (int, bool) {
	if len(s) == 0 {
		return 0, false
	}
	negative := false
	i := 0
	if s[0] == '-' || s[0] == '+' {
		negative = s[0] == '-'
		i++
		if len(s) == 1 {
			return 0, false
		}
	}
	var v int64
	for ; i < len(s); i++ {
		d := CharacterDigit(s[i], 10)
		if d < 0 {
			return 0, false
		}
		v = v*10 + int64(d)
		if v > math.MaxInt32+1 {
			return 0, false
		}
	}
	if negative {
		v = -v
	}
	if v > math.MaxInt32 || v < math.MinInt32 {
		return 0, false
	}
	return int(v), true
}

// mfCheckDatePattern checks the letters of a SimpleDateFormat pattern.
func mfCheckDatePattern(p []uint16) {
	inQuote := false
	for i := 0; i < len(p); i++ {
		c := p[i]
		if c == '\'' {
			if i+1 < len(p) && p[i+1] == '\'' {
				i++
				continue
			}
			inQuote = !inQuote
			continue
		}
		if inQuote {
			continue
		}
		if (c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z') && !strings.ContainsRune("GyYMLwWDdFEuaHkKhmsSzZX", rune(c)) {
			panic(NewIllegalArgumentExceptionMsg("Illegal pattern character '" + string(rune(c)) + "'"))
		}
	}
}

// Format is MessageFormat.format(Object[]) (through
// Format.format(Object)).
func (m *MessageFormat) Format(arguments []any) string {
	var sb strings.Builder
	for _, e := range m.elements {
		sb.WriteString(e.prefix)
		if arguments == nil || e.index >= len(arguments) {
			sb.WriteString("{" + strconv.Itoa(e.index) + "}")
			continue
		}
		arg := arguments[e.index]
		if IsNil(arg) {
			sb.WriteString("null")
			continue
		}
		num, isNum := mfNumberOf(arg)
		switch e.kind {
		case mfNumber:
			if !isNum {
				panic(NewIllegalArgumentExceptionMsg("Cannot format given Object as a Number"))
			}
			sb.WriteString(e.formatNumber(num))
		case mfDate, mfTime:
			if !isNum {
				panic(NewIllegalArgumentExceptionMsg("Cannot format given Object as a Date"))
			}
			panic(NewUnsupportedOperationExceptionMsg("MessageFormat: date and time formats are not supported"))
		case mfChoiceKind:
			if !isNum {
				panic(NewIllegalArgumentExceptionMsg("Cannot format given Object as a Number"))
			}
			s := e.choice.format(num.float())
			if strings.ContainsRune(s, '{') {
				s = MessageFormat_Format(s, arguments)
			}
			sb.WriteString(s)
		default:
			switch {
			case isNum:
				sb.WriteString(mfFormatDecimal(num, mfDefaultNumber))
			default:
				if s, ok := arg.(string); ok {
					sb.WriteString(s)
				} else {
					sb.WriteString(ToJavaString(arg))
				}
			}
		}
	}
	sb.WriteString(m.tail)
	return sb.String()
}

// FormatStrings formats string arguments without boxing []any or allocating a
// strings.Builder when the pattern has no placeholders.
func (m *MessageFormat) FormatStrings(arguments ...string) string {
	if len(m.elements) == 0 {
		return m.tail
	}
	n := len(m.tail)
	for i := range m.elements {
		e := &m.elements[i]
		if e.kind != mfNone {
			args := make([]any, len(arguments))
			for j, a := range arguments {
				args[j] = a
			}
			return m.Format(args)
		}
		n += len(e.prefix)
		if e.index < len(arguments) {
			n += len(arguments[e.index])
		} else {
			n += 3
		}
	}
	var sb strings.Builder
	sb.Grow(n)
	for i := range m.elements {
		e := &m.elements[i]
		sb.WriteString(e.prefix)
		if e.index >= len(arguments) {
			sb.WriteByte('{')
			sb.WriteString(strconv.Itoa(e.index))
			sb.WriteByte('}')
		} else {
			sb.WriteString(arguments[e.index])
		}
	}
	sb.WriteString(m.tail)
	return sb.String()
}

func (e *mfElement) formatNumber(n mfNum) string {
	switch e.style {
	case mfNumberInteger:
		return mfFormatDecimal(n, mfIntegerNumber)
	case mfNumberPercent:
		return mfFormatDecimal(n, mfPercentNumber)
	case mfNumberCurrency:
		return mfFormatDecimal(n, mfCurrencyNumber)
	case mfNumberPattern:
		panic(NewUnsupportedOperationExceptionMsg("MessageFormat: DecimalFormat patterns are not supported: " + e.custom))
	}
	return mfFormatDecimal(n, mfDefaultNumber)
}

// ---------------------------------------------------------------------------
// ChoiceFormat (see its Javadoc).
//
// Observed behaviour (black-box): a choice is a limit, a relation and a
// format, and choices are separated by '|'. The relation is '#' (the limit
// itself), U+2264 LESS-THAN OR EQUAL TO (the same) or '<' (the least double
// above the limit). Quoting is MessageFormat's. The limit text is "∞",
// "-∞", or anything Double.parseDouble accepts (it trims). Errors: a
// relation with no limit text before it (which includes a relation inside
// a format), and a limit that does not exceed the one before it. The
// separator commits a choice from what has been read, even without a
// relation: the last limit (0 if none yet) with the format text so far
// ("|1#a" makes 0 → "" and 1 → "a"). Limit text with no relation after it
// carries over to the next limit ("1#a|2|3#c" makes 23 → "c"), or is
// ignored at the end. Formatting picks the last choice whose limit is <=
// the number, else the first one (NaN too); with no choice at all it
// throws ArrayIndexOutOfBoundsException.

type mfChoice struct {
	limits  []float64
	formats []string
}

// mfParseChoice reads a ChoiceFormat pattern. It alternates between
// reading limit text and reading a format: a relation character closes the
// limit (mfCloseLimit) and opens a format, a separator commits the choice
// read so far, and an open format is committed at the end.
func mfParseChoice(p []uint16) *mfChoice {
	ch := &mfChoice{}
	var limitText, format []uint16
	limit := 0.0 // the limit of the choice being read
	formatOpen := false
	commit := func() {
		ch.limits = append(ch.limits, limit)
		ch.formats = append(ch.formats, StringFromUTF16(format))
		format = nil
	}
	quoted := false
	for i := 0; i < len(p); i++ {
		c := p[i]
		switch {
		case c == '\'' && i+1 < len(p) && p[i+1] == '\'':
			i++ // one literal quote, kept below
		case c == '\'':
			quoted = !quoted
			continue
		case quoted:
		case c == '|':
			commit()
			formatOpen = false
			continue
		case c == '#' || c == '<' || c == 0x2264:
			limit = mfCloseLimit(limitText, c == '<', limit, len(ch.limits) > 0)
			limitText = nil
			formatOpen = true
			continue
		}
		if formatOpen {
			format = append(format, c)
		} else {
			limitText = append(limitText, c)
		}
	}
	if formatOpen {
		commit()
	}
	return ch
}

// mfCloseLimit turns limit text into the limit of a choice: the parsed
// value, or the next double above it for a strict relation. With an
// earlier choice, the limit must exceed prev, the limit read before.
func mfCloseLimit(text []uint16, strict bool, prev float64, earlier bool) float64 {
	if len(text) == 0 {
		panic(NewIllegalArgumentExceptionMsg("Each interval must contain a number before a format"))
	}
	v := mfChoiceLimit(StringFromUTF16(text))
	if strict {
		v = mfNextDouble(v)
	}
	if earlier && v <= prev {
		panic(NewIllegalArgumentExceptionMsg("Incorrect order of intervals, must be in ascending order"))
	}
	return v
}

func mfChoiceLimit(s string) float64 {
	switch s {
	case "∞":
		return math.Inf(1)
	case "-∞":
		return math.Inf(-1)
	}
	return mfParseJavaDouble(s)
}

// mfParseJavaDouble is Double.parseDouble(s) (see the Javadoc of
// Double.valueOf(String)): surrounding characters <= ' ' are ignored, then
// an optional sign and "NaN", "Infinity", a decimal number with an optional
// exponent or a hexadecimal one with a binary exponent, optionally followed
// by f, F, d or D. It throws NumberFormatException otherwise.
func mfParseJavaDouble(s string) float64 {
	fail := func() { panic(NewNumberFormatException("For input string: \"" + s + "\"")) }
	t := []byte(strings.TrimFunc(s, func(r rune) bool { return r <= ' ' }))
	if len(t) == 0 {
		panic(NewNumberFormatException("empty String"))
	}
	sign := ""
	if t[0] == '+' || t[0] == '-' {
		sign, t = string(t[0]), t[1:]
	}
	switch string(t) {
	case "NaN":
		return math.NaN()
	case "Infinity":
		if sign == "-" {
			return math.Inf(-1)
		}
		return math.Inf(1)
	}
	if n := len(t); n > 0 && strings.IndexByte("fFdD", t[n-1]) >= 0 {
		t = t[:n-1]
	}
	isDigit := func(c byte, hex bool) bool {
		return c >= '0' && c <= '9' || hex && (c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F')
	}
	hex := len(t) > 2 && t[0] == '0' && (t[1] == 'x' || t[1] == 'X')
	i := 0
	if hex {
		i = 2
	}
	mantissaDigits := 0
	for ; i < len(t) && isDigit(t[i], hex); i++ {
		mantissaDigits++
	}
	if i < len(t) && t[i] == '.' {
		for i++; i < len(t) && isDigit(t[i], hex); i++ {
			mantissaDigits++
		}
	}
	if mantissaDigits == 0 {
		fail()
	}
	if i < len(t) && (!hex && (t[i] == 'e' || t[i] == 'E') || hex && (t[i] == 'p' || t[i] == 'P')) {
		i++
		if i < len(t) && (t[i] == '+' || t[i] == '-') {
			i++
		}
		start := i
		for i < len(t) && t[i] >= '0' && t[i] <= '9' {
			i++
		}
		if i == start {
			fail()
		}
	} else if hex {
		fail() // a hexadecimal number needs its binary exponent
	}
	if i != len(t) {
		fail()
	}
	v, err := strconv.ParseFloat(sign+string(t), 64)
	if err != nil && !errors.Is(err, strconv.ErrRange) {
		fail()
	}
	return v
}

// mfNextDouble is ChoiceFormat.nextDouble(double): the least double above
// d (infinities and NaN unchanged).
func mfNextDouble(d float64) float64 {
	if math.IsNaN(d) || math.IsInf(d, 0) {
		return d
	}
	return math.Nextafter(d, math.Inf(1))
}

func (ch *mfChoice) format(x float64) string {
	if len(ch.limits) == 0 {
		panic(NewArrayIndexOutOfBoundsExceptionForLength(0, 0))
	}
	i := 0
	for i < len(ch.limits) && x >= ch.limits[i] {
		i++
	}
	if i > 0 {
		i--
	}
	return ch.formats[i]
}

// ---------------------------------------------------------------------------
// Numbers

// mfNum is a java.lang.Number argument: an integer (exact) or a double.
type mfNum struct {
	isInt bool
	i     *big.Int
	f     float64
}

func (n mfNum) float() float64 {
	if n.isInt {
		f, _ := new(big.Float).SetInt(n.i).Float64()
		return f
	}
	return n.f
}

// mfNumberOf reports whether arg is a Java Number (Go's integer and
// floating-point types other than uint16, which is char, and *big.Int,
// which is BigInteger) and converts it.
func mfNumberOf(arg any) (mfNum, bool) {
	switch x := arg.(type) {
	case int:
		return mfNum{isInt: true, i: big.NewInt(int64(x))}, true
	case int8:
		return mfNum{isInt: true, i: big.NewInt(int64(x))}, true
	case int16:
		return mfNum{isInt: true, i: big.NewInt(int64(x))}, true
	case int32:
		return mfNum{isInt: true, i: big.NewInt(int64(x))}, true
	case int64:
		return mfNum{isInt: true, i: big.NewInt(x)}, true
	case uint8:
		return mfNum{isInt: true, i: big.NewInt(int64(x))}, true
	case uint32:
		return mfNum{isInt: true, i: new(big.Int).SetUint64(uint64(x))}, true
	case uint64:
		return mfNum{isInt: true, i: new(big.Int).SetUint64(x)}, true
	case uint:
		return mfNum{isInt: true, i: new(big.Int).SetUint64(uint64(x))}, true
	case float32:
		return mfNum{f: float64(x)}, true
	case float64:
		return mfNum{f: x}, true
	case *big.Int:
		return mfNum{isInt: true, i: x}, true
	case *int:
		return mfNumberOf(*x)
	case *int8:
		return mfNumberOf(*x)
	case *int16:
		return mfNumberOf(*x)
	case *int32:
		return mfNumberOf(*x)
	case *int64:
		return mfNumberOf(*x)
	case *float32:
		return mfNumberOf(*x)
	case *float64:
		return mfNumberOf(*x)
	}
	return mfNum{}, false
}

// mfDecimalFormat holds the settings of the NumberFormat instances used.
type mfDecimalFormat struct {
	minFrac, maxFrac int
	multiplier       int64
	prefix, suffix   string
}

var (
	mfDefaultNumber  = &mfDecimalFormat{maxFrac: 3, multiplier: 1}
	mfIntegerNumber  = &mfDecimalFormat{maxFrac: 0, multiplier: 1}
	mfPercentNumber  = &mfDecimalFormat{maxFrac: 0, multiplier: 100, suffix: "%"}
	mfCurrencyNumber = &mfDecimalFormat{minFrac: 2, maxFrac: 2, multiplier: 1, prefix: "$"}
)

// mfFormatDecimal formats n with grouping by three, rounding half-even
// (see the observed behaviour at the top of the file).
func mfFormatDecimal(n mfNum, f *mfDecimalFormat) string {
	var digits string // significant digits, no leading zeros ("" for zero)
	var exp int       // value = 0.digits × 10^exp
	negative := false
	if n.isInt {
		v := new(big.Int).Mul(n.i, big.NewInt(f.multiplier))
		negative = v.Sign() < 0
		s := new(big.Int).Abs(v).String()
		if s != "0" {
			digits, exp = strings.TrimRight(s, "0"), len(s)
		}
	} else {
		v := n.f
		if math.IsNaN(v) {
			return "NaN"
		}
		negative = math.Signbit(v)
		v *= float64(f.multiplier)
		if math.IsInf(v, 0) {
			return mfSign(negative) + f.prefix + "∞" + f.suffix
		}
		if a := math.Abs(v); a >= 1<<53 && a < 1<<63 {
			digits, exp = mfLongDigits(a)
		} else {
			digits, exp = mfShortestDigits(a)
			digits, exp = mfRoundHalfEven(digits, exp, f.maxFrac, a)
		}
	}
	var intPart, fracPart string
	switch {
	case exp <= 0:
		intPart = "0"
		fracPart = strings.Repeat("0", -exp) + digits
	case exp >= len(digits):
		intPart = digits + strings.Repeat("0", exp-len(digits))
	default:
		intPart, fracPart = digits[:exp], digits[exp:]
	}
	if intPart == "" {
		intPart = "0"
	}
	fracPart = strings.TrimRight(fracPart, "0")
	for len(fracPart) < f.minFrac {
		fracPart += "0"
	}
	var sb strings.Builder
	sb.WriteString(mfSign(negative))
	sb.WriteString(f.prefix)
	for i := 0; i < len(intPart); i++ {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			sb.WriteByte(',')
		}
		sb.WriteByte(intPart[i])
	}
	if fracPart != "" {
		sb.WriteByte('.')
		sb.WriteString(fracPart)
	}
	sb.WriteString(f.suffix)
	return sb.String()
}

func mfSign(negative bool) string {
	if negative {
		return "-"
	}
	return ""
}

// mfShortestDigits returns the shortest decimal digits of v > 0 that read
// back as v, as 0.digits × 10^exp ("" for zero).
func mfShortestDigits(v float64) (string, int) {
	if v == 0 {
		return "", 0
	}
	s := strconv.FormatFloat(v, 'e', -1, 64) // d.ddde±xx
	mant, expText, _ := strings.Cut(s, "e")
	e, _ := strconv.Atoi(expText)
	digits := strings.Replace(mant, ".", "", 1)
	return strings.TrimRight(digits, "0"), e + 1
}

// mfLongDigits returns the digits the JDK writes for an integral double v,
// 2^53 <= v < 2^63.
//
// Observed behaviour (black-box): these are not the shortest digits but
// those of the exact integer, rounded to a power of ten that grows with the
// binary exponent e: 10^floor((e-54)·log10 2), that is, exact up to 2^58,
// to tens up to 2^61 and to hundreds up to 2^63 (no ties can occur).
func mfLongDigits(v float64) (string, int) {
	e := math.Ilogb(v)
	k := 0
	if e > 54 {
		k = int(math.Floor(float64(e-54) * math.Log10(2)))
	}
	n := new(big.Int)
	new(big.Float).SetFloat64(v).Int(n)
	if k > 0 {
		p := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(k)), nil)
		half := new(big.Int).Rsh(p, 1)
		n.Add(n, half).Quo(n, p).Mul(n, p)
	}
	s := n.String()
	return strings.TrimRight(s, "0"), len(s)
}

// mfRoundHalfEven rounds 0.digits × 10^exp to maxFrac decimal places; exact
// is the value, which settles ties of the digits.
func mfRoundHalfEven(digits string, exp, maxFrac int, exact float64) (string, int) {
	keep := exp + maxFrac
	if keep >= len(digits) {
		return digits, exp
	}
	if keep < 0 {
		return "", 0
	}
	up := false
	switch dropped := digits[keep:]; {
	case dropped[0] > '5':
		up = true
	case dropped[0] < '5':
	case strings.TrimRight(dropped[1:], "0") != "":
		up = true
	case keep == 0 && digits == "5" && exp <= -3:
		// Observed behaviour (black-box): a lone 5 at the third decimal
		// place or further right counts as an exact tie.
		up = false
	default:
		// A tie in the shortest digits: compare the binary value with it.
		tie, _ := new(big.Rat).SetString("0." + digits[:keep] + "5e" + strconv.Itoa(exp))
		switch new(big.Rat).SetFloat64(exact).Cmp(tie) {
		case 1:
			up = true
		case 0:
			up = keep > 0 && (digits[keep-1]-'0')%2 == 1
		}
	}
	kept := []byte(digits[:keep])
	if up {
		i := len(kept) - 1
		for i >= 0 && kept[i] == '9' {
			kept[i] = '0'
			i--
		}
		if i < 0 {
			kept = append([]byte{'1'}, kept...)
			exp++
		} else {
			kept[i]++
		}
	}
	return strings.TrimRight(string(kept), "0"), exp
}
