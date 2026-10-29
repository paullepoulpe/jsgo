// Port of com/google/protobuf/TextFormat.java @ protobuf-java 4.30.2 (the text-format parser)

// Protocol Buffers - Google's data interchange format
// Copyright 2008 Google Inc.  All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd
//
// Here that LICENSE file is LICENSES/BSD-3-Clause-protobuf.txt.

// Go-only: language-adapter — protobuf is a leaf package (PORTING.md §2) and
// its messages are reconstructed (proto_runtime.go, which already prints the
// text format). This file translates the parsing half of
// com.google.protobuf.TextFormat (the version shaded into the pinned
// closure-compiler.jar): Parser (merge, the skipping of unknown fields) and
// parseInteger / parseJavaDouble, over the generated builders'
// MessageReflection.MergeTarget (proto_descriptors.go). The reconstructed
// messages have no extensions, no google.protobuf.Any and no groups, and no
// caller asks for a parse-info tree, so those branches are not translated;
// an extension name always resolves to "unknown", as with an empty
// ExtensionRegistry.

package jscomp

import (
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

const textFormat_DEFAULT_RECURSION_LIMIT = 100

// Java: TextFormat.PARSER
var textFormat_PARSER = &TextFormat_Parser{}

// TextFormat_Merge parses a text-format message from input and merges the
// contents into builder. It panics with TextFormat_ParseException (a checked
// exception).
//
// Java: TextFormat.merge(CharSequence, Message.Builder)
func TextFormat_Merge(input string, builder TextFormat_MergeTarget) {
	textFormat_PARSER.Merge(input, builder)
}

// TextFormat_Parser is the parser for text-format proto2 instances. This
// class is thread-safe. The implementation largely follows
// google/protobuf/text_format.cc.
type TextFormat_Parser struct{}

// CheckUnknownFields checks both unknown fields and unknown extensions and
// throws an exception if any were encountered.
//
// Java: TextFormat.Parser.checkUnknownFields(List<UnknownField>)
func (p *TextFormat_Parser) CheckUnknownFields(unknownFields []string) {
	if len(unknownFields) == 0 {
		return
	}

	var msg strings.Builder
	msg.WriteString("Input contains unknown fields and/or extensions:")
	for _, field := range unknownFields {
		msg.WriteByte('\n')
		msg.WriteString(field)
	}

	lineColumn := strings.Split(unknownFields[0], ":")
	line, _ := strconv.Atoi(lineColumn[0])
	column, _ := strconv.Atoi(lineColumn[1])
	panic(NewTextFormat_ParseExceptionAt(line, column, msg.String()))
}

// Merge parses a text-format message from input and merges the contents
// into builder.
//
// Java: TextFormat.Parser.merge(CharSequence, ExtensionRegistry, Message.Builder)
func (p *TextFormat_Parser) Merge(input string, target TextFormat_MergeTarget) {
	tokenizer := NewTextFormat_Tokenizer(input)
	var unknownFields []string

	for !tokenizer.AtEnd() {
		p.MergeField(tokenizer, target, &unknownFields, textFormat_DEFAULT_RECURSION_LIMIT)
	}
	p.CheckUnknownFields(unknownFields)
}

// MergeField parses a single field from tokenizer and merges it into target.
//
// Java: TextFormat.Parser.mergeField(Tokenizer, ExtensionRegistry, MergeTarget, TextFormatParseInfoTree.Builder, List<UnknownField>, int)
func (p *TextFormat_Parser) MergeField(
	tokenizer *TextFormat_Tokenizer,
	target TextFormat_MergeTarget,
	unknownFields *[]string,
	recursionLimit int) {
	var field *Descriptors_FieldDescriptor
	var name string
	type_ := target.GetDescriptorForType()

	// (A google.protobuf.Any type URL "[...]" needs the Any descriptor; no
	// reconstructed message is an Any.)

	if tokenizer.TryConsume("[") {
		// An extension.
		var nameBuilder strings.Builder
		nameBuilder.WriteString(tokenizer.ConsumeIdentifier())
		for tokenizer.TryConsume(".") {
			nameBuilder.WriteByte('.')
			nameBuilder.WriteString(tokenizer.ConsumeIdentifier())
		}
		name = nameBuilder.String()

		// target.findExtensionByName: the reconstructed messages declare no
		// extensions (an empty ExtensionRegistry).
		message := strconv.Itoa(tokenizer.GetPreviousLine()+1) +
			":" +
			strconv.Itoa(tokenizer.GetPreviousColumn()+1) +
			":\t" +
			type_.GetFullName() +
			".[" +
			name +
			"]"
		*unknownFields = append(*unknownFields, message)

		tokenizer.Consume("]")
	} else {
		name = tokenizer.ConsumeIdentifier()
		field = type_.FindFieldByName(name)

		// Group names are expected to be capitalized as they appear in the
		// .proto file, which actually matches their type names, not their field
		// names.
		if field == nil {
			// Explicitly specify US locale so that this code does not break when
			// executing in Turkey.
			lowerName := jutil.StringToLowerCaseRoot(name)
			field = type_.FindFieldByName(lowerName)
			// If the case-insensitive match worked but the field is NOT a group,
			if field != nil && !field.IsGroupLike() {
				field = nil
			}
			if field != nil && field.GetMessageType().GetName() != name {
				field = nil
			}
		}

		if field == nil {
			message := strconv.Itoa(tokenizer.GetPreviousLine()+1) +
				":" +
				strconv.Itoa(tokenizer.GetPreviousColumn()+1) +
				":\t" +
				type_.GetFullName() +
				"." +
				name
			*unknownFields = append(*unknownFields, message)
		}
	}

	// Skips unknown fields.
	if field == nil {
		p.GuessFieldTypeAndSkip(tokenizer, type_, recursionLimit)
		return
	}

	// Handle potential ':'.
	if field.GetJavaType() == Descriptors_FieldDescriptor_JavaType_MESSAGE {
		tokenizer.TryConsume(":") // optional
		p.ConsumeFieldValues(tokenizer, target, field, unknownFields, recursionLimit)
	} else {
		tokenizer.Consume(":") // required
		p.ConsumeFieldValues(tokenizer, target, field, unknownFields, recursionLimit)
	}

	// For historical reasons, fields may optionally be separated by commas or
	// semicolons.
	if !tokenizer.TryConsume(";") {
		tokenizer.TryConsume(",")
	}
}

// Java: TextFormat.Parser.consumeFullTypeName(Tokenizer)
func (p *TextFormat_Parser) ConsumeFullTypeName(tokenizer *TextFormat_Tokenizer) {
	// If there is not a leading `[`, this is just a type name.
	if !tokenizer.TryConsume("[") {
		tokenizer.ConsumeIdentifier()
		return
	}

	// Otherwise, this is an extension or google.protobuf.Any type URL: we consume proto path
	// elements until we've addressed the type.
	tokenizer.ConsumeIdentifier()
	for tokenizer.TryConsume(".") {
		tokenizer.ConsumeIdentifier()
	}
	if tokenizer.TryConsume("/") {
		tokenizer.ConsumeIdentifier()
		for tokenizer.TryConsume(".") {
			tokenizer.ConsumeIdentifier()
		}
	}
	tokenizer.Consume("]")
}

// ConsumeFieldValues parses one or more field values from tokenizer and
// merges them into target.
//
// Java: TextFormat.Parser.consumeFieldValues(Tokenizer, ExtensionRegistry, MergeTarget, FieldDescriptor, ExtensionInfo, TextFormatParseInfoTree.Builder, List<UnknownField>, int)
func (p *TextFormat_Parser) ConsumeFieldValues(
	tokenizer *TextFormat_Tokenizer,
	target TextFormat_MergeTarget,
	field *Descriptors_FieldDescriptor,
	unknownFields *[]string,
	recursionLimit int) {
	// Support specifying repeated field values as a comma-separated list.
	// Ex."foo: [1, 2, 3]"
	if field.IsRepeated() && tokenizer.TryConsume("[") {
		if !tokenizer.TryConsume("]") { // Allow "foo: []" to be treated as empty.
			for {
				p.ConsumeFieldValue(tokenizer, target, field, unknownFields, recursionLimit)
				if tokenizer.TryConsume("]") {
					// End of list.
					break
				}
				tokenizer.Consume(",")
			}
		}
	} else {
		p.ConsumeFieldValue(tokenizer, target, field, unknownFields, recursionLimit)
	}
}

// ConsumeFieldValue parses a single field value from tokenizer and merges it
// into target.
//
// Java: TextFormat.Parser.consumeFieldValue(Tokenizer, ExtensionRegistry, MergeTarget, FieldDescriptor, ExtensionInfo, TextFormatParseInfoTree.Builder, List<UnknownField>, int)
func (p *TextFormat_Parser) ConsumeFieldValue(
	tokenizer *TextFormat_Tokenizer,
	target TextFormat_MergeTarget,
	field *Descriptors_FieldDescriptor,
	unknownFields *[]string,
	recursionLimit int) {
	var value any

	if field.GetJavaType() == Descriptors_FieldDescriptor_JavaType_MESSAGE {
		if recursionLimit < 1 {
			panic(tokenizer.ParseException("Message is nested too deep"))
		}

		var endToken string
		if tokenizer.TryConsume("<") {
			endToken = ">"
		} else {
			tokenizer.Consume("{")
			endToken = "}"
		}

		subField := target.protoNewMergeTargetForField(field)

		for !tokenizer.TryConsume(endToken) {
			if tokenizer.AtEnd() {
				panic(tokenizer.ParseException("Expected \"" + endToken + "\"."))
			}
			p.MergeField(tokenizer, subField, unknownFields, recursionLimit-1)
		}

		value = subField.protoFinish()
	} else {
		switch field.GetType() {
		case Descriptors_FieldDescriptor_Type_INT32,
			Descriptors_FieldDescriptor_Type_SINT32,
			Descriptors_FieldDescriptor_Type_SFIXED32:
			value = tokenizer.ConsumeInt32()

		case Descriptors_FieldDescriptor_Type_INT64,
			Descriptors_FieldDescriptor_Type_SINT64,
			Descriptors_FieldDescriptor_Type_SFIXED64:
			value = tokenizer.ConsumeInt64()

		case Descriptors_FieldDescriptor_Type_UINT32,
			Descriptors_FieldDescriptor_Type_FIXED32:
			value = tokenizer.ConsumeUInt32()

		case Descriptors_FieldDescriptor_Type_UINT64,
			Descriptors_FieldDescriptor_Type_FIXED64:
			value = tokenizer.ConsumeUInt64()

		case Descriptors_FieldDescriptor_Type_FLOAT:
			value = tokenizer.ConsumeFloat()

		case Descriptors_FieldDescriptor_Type_DOUBLE:
			value = tokenizer.ConsumeDouble()

		case Descriptors_FieldDescriptor_Type_BOOL:
			value = tokenizer.ConsumeBoolean()

		case Descriptors_FieldDescriptor_Type_STRING:
			value = tokenizer.ConsumeString()

		case Descriptors_FieldDescriptor_Type_BYTES:
			value = tokenizer.ConsumeByteString()

		case Descriptors_FieldDescriptor_Type_ENUM:
			enumType := field.GetEnumType()

			if tokenizer.LookingAtInteger() {
				number := tokenizer.ConsumeInt32()
				// (Every reconstructed enum that a parsed message uses is closed.)
				var enumValue *Descriptors_EnumValueDescriptor
				if enumType.IsClosed() {
					enumValue = enumType.FindValueByNumber(number)
				} else {
					panic(jutil.NewUnsupportedOperationExceptionMsg("open enum " + enumType.GetFullName()))
				}
				if enumValue == nil {
					panic(tokenizer.ParseExceptionPreviousToken("Enum type \"" +
						enumType.GetFullName() +
						"\" has no value with number " +
						strconv.Itoa(int(number)) +
						"."))
				}
				value = enumValue
			} else {
				id := tokenizer.ConsumeIdentifier()
				enumValue := enumType.FindValueByName(id)
				if enumValue == nil {
					panic(tokenizer.ParseExceptionPreviousToken("Enum type \"" +
						enumType.GetFullName() +
						"\" has no value named \"" +
						id +
						"\"."))
				}
				value = enumValue
			}

		case Descriptors_FieldDescriptor_Type_MESSAGE, Descriptors_FieldDescriptor_Type_GROUP:
			panic(jutil.NewRuntimeException("Can't get here."))
		}
	}

	if field.IsRepeated() {
		target.protoAddRepeatedField(field, value)
	} else {
		target.protoSetField(field, value)
	}
}

// SkipField skips the next field including the field's name and value.
//
// Java: TextFormat.Parser.skipField(Tokenizer, Descriptor, int)
func (p *TextFormat_Parser) SkipField(tokenizer *TextFormat_Tokenizer, type_ *Descriptors_Descriptor, recursionLimit int) {
	p.ConsumeFullTypeName(tokenizer)
	p.GuessFieldTypeAndSkip(tokenizer, type_, recursionLimit)

	if !tokenizer.TryConsume(";") {
		tokenizer.TryConsume(",")
	}
}

// SkipFieldMessage skips the whole body of a message including the beginning
// delimiter and the ending delimiter.
//
// Java: TextFormat.Parser.skipFieldMessage(Tokenizer, Descriptor, int)
func (p *TextFormat_Parser) SkipFieldMessage(tokenizer *TextFormat_Tokenizer, type_ *Descriptors_Descriptor, recursionLimit int) {
	var delimiter string
	if tokenizer.TryConsume("<") {
		delimiter = ">"
	} else {
		tokenizer.Consume("{")
		delimiter = "}"
	}
	for !tokenizer.LookingAt(">") && !tokenizer.LookingAt("}") {
		p.SkipField(tokenizer, type_, recursionLimit)
	}
	tokenizer.Consume(delimiter)
}

// SkipFieldValue skips a field value.
//
// Java: TextFormat.Parser.skipFieldValue(Tokenizer)
func (p *TextFormat_Parser) SkipFieldValue(tokenizer *TextFormat_Tokenizer) {
	if !tokenizer.TryConsumeByteString() &&
		!tokenizer.TryConsumeIdentifier() && // includes enum & boolean
		!tokenizer.TryConsumeInt64() && // includes int32
		!tokenizer.TryConsumeUInt64() && // includes uint32
		!tokenizer.TryConsumeDouble() &&
		!tokenizer.TryConsumeFloat() {
		panic(tokenizer.ParseException("Invalid field value: " + tokenizer.currentToken))
	}
}

// GuessFieldTypeAndSkip tries to guess the type of this field and skips it.
//
// Java: TextFormat.Parser.guessFieldTypeAndSkip(Tokenizer, Descriptor, int)
func (p *TextFormat_Parser) GuessFieldTypeAndSkip(tokenizer *TextFormat_Tokenizer, type_ *Descriptors_Descriptor, recursionLimit int) {
	semicolonConsumed := tokenizer.TryConsume(":")
	if tokenizer.LookingAt("[") {
		p.SkipFieldShortFormedRepeated(tokenizer, semicolonConsumed, type_, recursionLimit)
	} else if semicolonConsumed && !tokenizer.LookingAt("{") && !tokenizer.LookingAt("<") {
		p.SkipFieldValue(tokenizer)
	} else {
		if recursionLimit < 1 {
			panic(tokenizer.ParseException("Message is nested too deep"))
		}
		p.SkipFieldMessage(tokenizer, type_, recursionLimit-1)
	}
}

// SkipFieldShortFormedRepeated skips a short-formed repeated field value.
//
// Reports an error if scalar type is not allowed but showing up inside "[]".
//
// Java: TextFormat.Parser.skipFieldShortFormedRepeated(Tokenizer, boolean, Descriptor, int)
func (p *TextFormat_Parser) SkipFieldShortFormedRepeated(
	tokenizer *TextFormat_Tokenizer, scalarAllowed bool, type_ *Descriptors_Descriptor, recursionLimit int) {
	if !tokenizer.TryConsume("[") || tokenizer.TryConsume("]") {
		return
	}

	for {
		if tokenizer.LookingAt("{") || tokenizer.LookingAt("<") {
			if recursionLimit < 1 {
				panic(tokenizer.ParseException("Message is nested too deep"))
			}
			p.SkipFieldMessage(tokenizer, type_, recursionLimit-1)
		} else if scalarAllowed {
			p.SkipFieldValue(tokenizer)
		} else {
			panic(tokenizer.ParseException(
				"Invalid repeated scalar field: missing \":\" before \"[\"."))
		}
		if tokenizer.TryConsume("]") {
			break
		}
		tokenizer.Consume(",")
	}
}

func textFormat_parseInteger(text string, isSigned bool, isLong bool) (int64, string, bool) {
	pos := 0

	negative := false
	if strings.HasPrefix(text[pos:], "-") {
		if !isSigned {
			return 0, "Number must be positive: " + text, false
		}
		pos++
		negative = true
	}

	radix := 10
	if strings.HasPrefix(text[pos:], "0x") {
		pos += 2
		radix = 16
	} else if strings.HasPrefix(text[pos:], "0") {
		radix = 8
	}

	numberText := jutil.StringSubstringFrom(text, pos)

	var result int64
	if jutil.StringLength(numberText) < 16 {
		var errMsg string
		var ok bool
		result, errMsg, ok = textFormat_parseLong(numberText, radix)
		if !ok {
			return 0, errMsg, false
		}
		if negative {
			result = -result
		}

		if !isLong {
			if isSigned {
				if result > 2147483647 || result < -2147483648 {
					return 0, "Number out of range for 32-bit signed integer: " + text, false
				}
			} else {
				if result >= (int64(1)<<32) || result < 0 {
					return 0, "Number out of range for 32-bit unsigned integer: " + text, false
				}
			}
		}
	} else {
		bigValue, errMsg, ok := textFormat_parseBigInteger(numberText, radix)
		if !ok {
			return 0, errMsg, false
		}
		if negative {
			bigValue = new(big.Int).Neg(bigValue)
		}

		// BigInteger.bitLength(): the minimal two's-complement width
		// excluding the sign bit.
		bitLength := bigValue.BitLen()
		if bigValue.Sign() < 0 {
			bitLength = new(big.Int).Not(bigValue).BitLen()
		}
		if !isLong {
			if isSigned {
				if bitLength > 31 {
					return 0, "Number out of range for 32-bit signed integer: " + text, false
				}
			} else {
				if bitLength > 32 {
					return 0, "Number out of range for 32-bit unsigned integer: " + text, false
				}
			}
		} else {
			if isSigned {
				if bitLength > 63 {
					return 0, "Number out of range for 64-bit signed integer: " + text, false
				}
			} else {
				if bitLength > 64 {
					return 0, "Number out of range for 64-bit unsigned integer: " + text, false
				}
			}
		}

		// BigInteger.longValue(): the low-order 64 bits.
		result = int64(new(big.Int).And(bigValue, new(big.Int).SetUint64(^uint64(0))).Uint64())
	}

	return result, "", true
}

func textFormat_badNumberInput(s string, radix int) string {
	msg := "For input string: \"" + s + "\""
	if radix != 10 {
		msg += " under radix " + strconv.Itoa(radix)
	}
	return msg
}

func textFormat_parseLong(s string, radix int) (int64, string, bool) {
	units := jutil.StringToUTF16(s)
	i := 0
	negative := false
	if len(units) > 0 && (units[0] == '-' || units[0] == '+') {
		negative = units[0] == '-'
		i = 1
	}
	if i == len(units) {
		return 0, textFormat_badNumberInput(s, radix), false
	}
	limit := uint64(math.MaxInt64)
	if negative {
		limit++
	}
	base := uint64(radix)
	var magnitude uint64
	for ; i < len(units); i++ {
		digit := jutil.CharacterDigit(units[i], radix)
		if digit < 0 || magnitude > (limit-uint64(digit))/base {
			return 0, textFormat_badNumberInput(s, radix), false
		}
		magnitude = magnitude*base + uint64(digit)
	}
	if negative {
		return -int64(magnitude), "", true
	}
	return int64(magnitude), "", true
}

func textFormat_parseBigInteger(s string, radix int) (*big.Int, string, bool) {
	units := jutil.StringToUTF16(s)
	start := 0
	negative := false
	if len(units) > 0 && (units[0] == '-' || units[0] == '+') {
		negative = units[0] == '-'
		start = 1
	}
	if start == len(units) {
		return nil, "Zero length BigInteger", false
	}
	for _, c := range units[start:] {
		if c == '-' || c == '+' {
			return nil, "Illegal embedded sign character", false
		}
	}
	for start < len(units) && jutil.CharacterDigit(units[start], radix) == 0 {
		start++
	}
	if start == len(units) {
		return big.NewInt(0), "", true
	}
	digits := units[start:]
	group := 0
	for p := int64(radix); p <= math.MaxInt32; p *= int64(radix) {
		group++
	}
	first := len(digits) % group
	if first == 0 {
		first = group
	}
	text := make([]byte, len(digits))
	for from, to := 0, first; from < len(digits); from, to = to, to+group {
		for i := from; i < to; i++ {
			d := jutil.CharacterDigit(digits[i], radix)
			if d < 0 {
				return nil, textFormat_badNumberInput(jutil.StringFromUTF16(digits[from:to]), radix), false
			}
			text[i] = "0123456789abcdefghijklmnopqrstuvwxyz"[d]
		}
	}
	x, ok := new(big.Int).SetString(string(text), radix)
	if !ok {
		return nil, textFormat_badNumberInput(s, radix), false
	}
	if negative {
		x.Neg(x)
	}
	return x, "", true
}

func textFormat_parseJavaDouble(s string) (float64, string, bool) {
	t := jutil.StringTrim(s)
	if t == "" {
		return 0, "empty String", false
	}
	i := 0
	negative := false
	if t[0] == '-' || t[0] == '+' {
		negative = t[0] == '-'
		i = 1
	}
	body := t[i:]
	if body == "" {
		return 0, "For input string: \"" + t + "\"", false
	}
	var v float64
	var errKind int // 0=ok, 1=multiplePoints, 2=invalid
	switch {
	case body[0] == 'N':
		if body != "NaN" {
			return 0, "For input string: \"" + t + "\"", false
		}
		return jutil.Double_NaN, "", true
	case body[0] == 'I':
		if body != "Infinity" {
			return 0, "For input string: \"" + t + "\"", false
		}
		v = math.Inf(1)
	case len(body) >= 2 && body[0] == '0' && (body[1] == 'x' || body[1] == 'X'):
		v, errKind = textFormat_parseHexDouble(body[2:])
	default:
		v, errKind = textFormat_parseDecimalDouble(body)
	}
	switch errKind {
	case 0:
		if negative {
			v = -v
		}
		return v, "", true
	case 1:
		return 0, "multiple points", false
	default:
		return 0, "For input string: \"" + t + "\"", false
	}
}

func textFormat_skipFloatSuffix(s string, i int) bool {
	if i < len(s) {
		switch s[i] {
		case 'f', 'F', 'd', 'D':
			i++
		}
	}
	return i == len(s)
}

func textFormat_scanSignedExponent(s string, i int) (int64, int, bool) {
	var value int64
	negative := false
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		negative = s[i] == '-'
		i++
	}
	start := i
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		if value < 1e15 {
			value = value*10 + int64(s[i]-'0')
		}
		i++
	}
	if i == start {
		return 0, i, false
	}
	if negative {
		value = -value
	}
	return value, i, true
}

func textFormat_parseDecimalDouble(s string) (float64, int) {
	i, points, digitCount := 0, 0, 0
	for i < len(s) && ((s[i] >= '0' && s[i] <= '9') || s[i] == '.') {
		if s[i] == '.' {
			points++
		} else {
			digitCount++
		}
		i++
	}
	if points > 1 {
		return 0, 1
	}
	if digitCount == 0 {
		return 0, 2
	}
	significandEnd := i
	var exponent int64
	if i < len(s) && (s[i] == 'e' || s[i] == 'E') {
		var ok bool
		exponent, i, ok = textFormat_scanSignedExponent(s, i+1)
		if !ok {
			return 0, 2
		}
	}
	literalEnd := i
	if !textFormat_skipFloatSuffix(s, i) {
		return 0, 2
	}
	if exponent > -10000 && exponent < 10000 {
		if v, err := strconv.ParseFloat(s[:literalEnd], 64); err == nil {
			return v, 0
		} else if ne, ok := err.(*strconv.NumError); ok && ne.Err == strconv.ErrRange {
			return v, 0
		}
	}
	sig := make([]byte, 0, digitCount)
	var intDigits int64
	seenPoint := false
	for j := 0; j < significandEnd; j++ {
		switch c := s[j]; {
		case c == '.':
			seenPoint = true
		case c == '0' && len(sig) == 0:
			if seenPoint {
				intDigits--
			}
		default:
			sig = append(sig, c)
			if !seenPoint {
				intDigits++
			}
		}
	}
	if len(sig) == 0 {
		return 0, 0
	}
	scale10 := intDigits + exponent
	switch {
	case scale10 > 310:
		return math.Inf(1), 0
	case scale10 < -400:
		return 0, 0
	}
	v, _ := strconv.ParseFloat("0."+string(sig)+"e"+strconv.FormatInt(scale10, 10), 64)
	return v, 0
}

func textFormat_parseHexDouble(s string) (float64, int) {
	i, digitCount := 0, 0
	sig := make([]byte, 0, len(s))
	var intDigits int64
	seenPoint := false
	for i < len(s) {
		c := s[i]
		switch {
		case textFormat_isHex(c):
			digitCount++
			if c != '0' || len(sig) > 0 {
				sig = append(sig, c)
				if !seenPoint {
					intDigits++
				}
			} else if seenPoint {
				intDigits--
			}
		case c == '.' && !seenPoint:
			seenPoint = true
		default:
			goto digitsDone
		}
		i++
	}
digitsDone:
	if digitCount == 0 || i == len(s) || (s[i] != 'p' && s[i] != 'P') {
		return 0, 2
	}
	exponent, end, ok := textFormat_scanSignedExponent(s, i+1)
	if !ok || !textFormat_skipFloatSuffix(s, end) {
		return 0, 2
	}
	if len(sig) == 0 {
		return 0, 0
	}
	scale2 := 4*intDigits + exponent
	switch {
	case scale2 > 1100:
		return math.Inf(1), 0
	case scale2 < -1200:
		return 0, 0
	}
	v, _ := strconv.ParseFloat("0x0."+string(sig)+"p"+strconv.FormatInt(scale2, 10), 64)
	return v, 0
}
