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
// its messages are reconstructed (proto_runtime.go). Closure Compiler only
// parses text-format ConformanceConfig protos, so field dispatch is wired
// directly onto ConformanceConfig, Requirement, and RequirementScopeEntry
// without runtime descriptor tables or reflection.

package jscomp

import (
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

const textFormat_DEFAULT_RECURSION_LIMIT = 100

// ConformanceTextFormatTarget is implemented by the reconstructed conformance
// message builders that support textproto parsing.
type ConformanceTextFormatTarget interface {
	conformanceFullName() string
	conformanceMergeField(name string, t *TextFormat_Tokenizer, unknownFields *[]string, recursionLimit int) (bool, *TextFormat_ParseException)
}

// TextFormat_Merge parses a text-format message from input and merges the
// contents into builder, panicking with *TextFormat_ParseException on error.
//
// Java: TextFormat.merge(CharSequence, Message.Builder)
func TextFormat_Merge(input string, builder ConformanceTextFormatTarget) {
	if err := TextFormat_MergeErr(input, builder); err != nil {
		panic(err)
	}
}

// TextFormat_MergeErr parses a text-format message from input and merges the
// contents into builder, returning a *TextFormat_ParseException on malformed
// input instead of panicking.
func TextFormat_MergeErr(input string, builder ConformanceTextFormatTarget) *TextFormat_ParseException {
	tokenizer := NewTextFormat_Tokenizer(input)
	var unknownFields []string
	for !tokenizer.AtEnd() {
		if err := textFormat_mergeField(tokenizer, builder, &unknownFields, textFormat_DEFAULT_RECURSION_LIMIT); err != nil {
			return err
		}
	}
	return textFormat_checkUnknownFields(unknownFields)
}

func (*ConformanceConfig) conformanceFullName() string     { return "jscomp.ConformanceConfig" }
func (*Requirement) conformanceFullName() string           { return "jscomp.Requirement" }
func (*RequirementScopeEntry) conformanceFullName() string { return "jscomp.RequirementScopeEntry" }

func (b *ConformanceConfig) conformanceMergeField(
	name string, t *TextFormat_Tokenizer, unknownFields *[]string, recursionLimit int,
) (bool, *TextFormat_ParseException) {
	switch name {
	case "requirement":
		return consumeProtoMessageField(t, true, unknownFields, recursionLimit, Requirement_NewBuilder, b.AddRequirement)
	case "library_level_non_allowlisted_conformance_violations_behavior":
		return consumeProtoEnumField(t,
			"jscomp.ConformanceConfig.LibraryLevelNonAllowlistedConformanceViolationsBehavior",
			conformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_names[:],
			ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_ForNumber,
			b.SetLibraryLevelNonAllowlistedConformanceViolationsBehavior)
	}
	return false, nil
}

func (b *Requirement) conformanceMergeField(
	name string, t *TextFormat_Tokenizer, unknownFields *[]string, recursionLimit int,
) (bool, *TextFormat_ParseException) {
	switch name {
	case "error_message":
		return consumeProtoStringField(t, false, b.SetErrorMessage)
	case "whitelist":
		return consumeProtoStringField(t, true, b.AddWhitelist)
	case "whitelist_regexp":
		return consumeProtoStringField(t, true, b.AddWhitelistRegexp)
	case "allowlist":
		return consumeProtoStringField(t, true, b.AddAllowlist)
	case "allowlist_regexp":
		return consumeProtoStringField(t, true, b.AddAllowlistRegexp)
	case "whitelist_entry":
		return consumeProtoMessageField(t, true, unknownFields, recursionLimit, RequirementScopeEntry_NewBuilder, b.AddWhitelistEntry)
	case "allowlist_entry":
		return consumeProtoMessageField(t, true, unknownFields, recursionLimit, RequirementScopeEntry_NewBuilder, b.AddAllowlistEntry)
	case "only_apply_to":
		return consumeProtoStringField(t, true, b.AddOnlyApplyTo)
	case "only_apply_to_regexp":
		return consumeProtoStringField(t, true, b.AddOnlyApplyToRegexp)
	case "type":
		return consumeProtoEnumField(t, "jscomp.Requirement.Type", requirement_Type_names[:], Requirement_Type_ForNumber, b.SetType)
	case "value":
		return consumeProtoStringField(t, true, b.AddValue)
	case "allow_extending_value":
		return consumeProtoBoolField(t, b.SetAllowExtendingValue)
	case "type_matching_strategy":
		return consumeProtoEnumField(t, "jscomp.Requirement.TypeMatchingStrategy", requirement_TypeMatchingStrategy_names[:], Requirement_TypeMatchingStrategy_ForNumber, b.SetTypeMatchingStrategy)
	case "java_class":
		return consumeProtoStringField(t, false, b.SetJavaClass)
	case "rule_id":
		return consumeProtoStringField(t, false, b.SetRuleId)
	case "extends":
		return consumeProtoStringField(t, false, b.SetExtends)
	case "report_loose_type_violations":
		return consumeProtoBoolField(t, b.SetReportLooseTypeViolations)
	case "severity":
		return consumeProtoEnumField(t, "jscomp.Requirement.Severity", requirement_Severity_names[:], Requirement_Severity_ForNumber, b.SetSeverity)
	case "config_file":
		return consumeProtoStringField(t, true, b.AddConfigFile)
	}
	return false, nil
}

func (b *RequirementScopeEntry) conformanceMergeField(
	name string, t *TextFormat_Tokenizer, _ *[]string, _ int,
) (bool, *TextFormat_ParseException) {
	switch name {
	case "reason":
		return consumeProtoEnumField(t, "jscomp.RequirementScopeEntry.Reason", requirementScopeEntry_Reason_names[:], RequirementScopeEntry_Reason_ForNumber, b.SetReason)
	case "prefix":
		return consumeProtoStringField(t, true, b.AddPrefix)
	case "regexp":
		return consumeProtoStringField(t, true, b.AddRegexp)
	case "explanation":
		return consumeProtoStringField(t, false, b.SetExplanation)
	case "comment":
		return consumeProtoStringField(t, true, b.AddComment)
	case "automatically_prune":
		return consumeProtoBoolField(t, b.SetAutomaticallyPrune)
	case "do_not_record_violations_in_summary_because_this_is_gencode_and_always_safe":
		return consumeProtoBoolField(t, b.SetDoNotRecordViolationsInSummaryBecauseThisIsGencodeAndAlwaysSafe)
	}
	return false, nil
}

func consumeProtoScalarValues(
	t *TextFormat_Tokenizer, repeated bool, consumeOne func() *TextFormat_ParseException,
) (bool, *TextFormat_ParseException) {
	if err := t.Consume(":"); err != nil {
		return true, err
	}
	if repeated && t.TryConsume("[") {
		if !t.TryConsume("]") {
			for {
				if err := consumeOne(); err != nil {
					return true, err
				}
				if t.TryConsume("]") {
					break
				}
				if err := t.Consume(","); err != nil {
					return true, err
				}
			}
		}
		return true, nil
	}
	return true, consumeOne()
}

func consumeProtoStringField[R any](
	t *TextFormat_Tokenizer, repeated bool, set func(string) R,
) (bool, *TextFormat_ParseException) {
	return consumeProtoScalarValues(t, repeated, func() *TextFormat_ParseException {
		v, err := t.ConsumeString()
		if err != nil {
			return err
		}
		set(v)
		return nil
	})
}

func consumeProtoBoolField[R any](
	t *TextFormat_Tokenizer, set func(bool) R,
) (bool, *TextFormat_ParseException) {
	return consumeProtoScalarValues(t, false, func() *TextFormat_ParseException {
		v, err := t.ConsumeBoolean()
		if err != nil {
			return err
		}
		set(v)
		return nil
	})
}

func consumeProtoEnumField[E ~int32, R any](
	t *TextFormat_Tokenizer, enumFullName string, names []string, forNumber func(int32) (E, bool), set func(E) R,
) (bool, *TextFormat_ParseException) {
	return consumeProtoScalarValues(t, false, func() *TextFormat_ParseException {
		if t.LookingAtInteger() {
			number, err := t.ConsumeInt32()
			if err != nil {
				return err
			}
			enumValue, ok := forNumber(number)
			if !ok {
				return t.ParseExceptionPreviousToken("Enum type \"" +
					enumFullName +
					"\" has no value with number " +
					strconv.Itoa(int(number)) +
					".")
			}
			set(enumValue)
			return nil
		}
		id, err := t.ConsumeIdentifier()
		if err != nil {
			return err
		}
		for i, name := range names {
			if name == id {
				set(E(i))
				return nil
			}
		}
		return t.ParseExceptionPreviousToken("Enum type \"" +
			enumFullName +
			"\" has no value named \"" +
			id +
			"\".")
	})
}

func consumeProtoMessageField[M ConformanceTextFormatTarget, R any](
	t *TextFormat_Tokenizer,
	repeated bool,
	unknownFields *[]string,
	recursionLimit int,
	newMsg func() M,
	add func(M) R,
) (bool, *TextFormat_ParseException) {
	t.TryConsume(":") // optional before message value
	consumeOne := func() *TextFormat_ParseException {
		if recursionLimit < 1 {
			return t.ParseException("Message is nested too deep")
		}
		var endToken string
		if t.TryConsume("<") {
			endToken = ">"
		} else {
			if err := t.Consume("{"); err != nil {
				return err
			}
			endToken = "}"
		}
		sub := newMsg()
		for !t.TryConsume(endToken) {
			if t.AtEnd() {
				return t.ParseException("Expected \"" + endToken + "\".")
			}
			if err := textFormat_mergeField(t, sub, unknownFields, recursionLimit-1); err != nil {
				return err
			}
		}
		add(sub)
		return nil
	}
	if repeated && t.TryConsume("[") {
		if !t.TryConsume("]") {
			for {
				if err := consumeOne(); err != nil {
					return true, err
				}
				if t.TryConsume("]") {
					break
				}
				if err := t.Consume(","); err != nil {
					return true, err
				}
			}
		}
		return true, nil
	}
	return true, consumeOne()
}

// Java: TextFormat.Parser.checkUnknownFields(List<UnknownField>)
func textFormat_checkUnknownFields(unknownFields []string) *TextFormat_ParseException {
	if len(unknownFields) == 0 {
		return nil
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
	return NewTextFormat_ParseExceptionAt(line, column, msg.String())
}

// Java: TextFormat.Parser.mergeField(...)
func textFormat_mergeField(
	tokenizer *TextFormat_Tokenizer,
	target ConformanceTextFormatTarget,
	unknownFields *[]string,
	recursionLimit int,
) *TextFormat_ParseException {
	fullName := target.conformanceFullName()

	if tokenizer.TryConsume("[") {
		var nameBuilder strings.Builder
		id, err := tokenizer.ConsumeIdentifier()
		if err != nil {
			return err
		}
		nameBuilder.WriteString(id)
		for tokenizer.TryConsume(".") {
			nameBuilder.WriteByte('.')
			id, err = tokenizer.ConsumeIdentifier()
			if err != nil {
				return err
			}
			nameBuilder.WriteString(id)
		}
		*unknownFields = append(*unknownFields,
			strconv.Itoa(tokenizer.GetPreviousLine()+1)+":"+
				strconv.Itoa(tokenizer.GetPreviousColumn()+1)+":\t"+
				fullName+".["+nameBuilder.String()+"]")
		if err := tokenizer.Consume("]"); err != nil {
			return err
		}
		return textFormat_guessFieldTypeAndSkip(tokenizer, recursionLimit)
	}

	name, err := tokenizer.ConsumeIdentifier()
	if err != nil {
		return err
	}
	matched, err := target.conformanceMergeField(name, tokenizer, unknownFields, recursionLimit)
	if err != nil {
		return err
	}
	if !matched {
		*unknownFields = append(*unknownFields,
			strconv.Itoa(tokenizer.GetPreviousLine()+1)+":"+
				strconv.Itoa(tokenizer.GetPreviousColumn()+1)+":\t"+
				fullName+"."+name)
		return textFormat_guessFieldTypeAndSkip(tokenizer, recursionLimit)
	}

	// For historical reasons, fields may optionally be separated by commas or semicolons.
	if !tokenizer.TryConsume(";") {
		tokenizer.TryConsume(",")
	}
	return nil
}

// Java: TextFormat.Parser.consumeFullTypeName(Tokenizer)
func textFormat_consumeFullTypeName(tokenizer *TextFormat_Tokenizer) *TextFormat_ParseException {
	if !tokenizer.TryConsume("[") {
		_, err := tokenizer.ConsumeIdentifier()
		return err
	}
	if _, err := tokenizer.ConsumeIdentifier(); err != nil {
		return err
	}
	for tokenizer.TryConsume(".") {
		if _, err := tokenizer.ConsumeIdentifier(); err != nil {
			return err
		}
	}
	if tokenizer.TryConsume("/") {
		if _, err := tokenizer.ConsumeIdentifier(); err != nil {
			return err
		}
		for tokenizer.TryConsume(".") {
			if _, err := tokenizer.ConsumeIdentifier(); err != nil {
				return err
			}
		}
	}
	return tokenizer.Consume("]")
}

// Java: TextFormat.Parser.skipField(Tokenizer, Descriptor, int)
func textFormat_skipField(tokenizer *TextFormat_Tokenizer, recursionLimit int) *TextFormat_ParseException {
	if err := textFormat_consumeFullTypeName(tokenizer); err != nil {
		return err
	}
	if err := textFormat_guessFieldTypeAndSkip(tokenizer, recursionLimit); err != nil {
		return err
	}
	if !tokenizer.TryConsume(";") {
		tokenizer.TryConsume(",")
	}
	return nil
}

// Java: TextFormat.Parser.skipFieldMessage(Tokenizer, Descriptor, int)
func textFormat_skipFieldMessage(tokenizer *TextFormat_Tokenizer, recursionLimit int) *TextFormat_ParseException {
	var delimiter string
	if tokenizer.TryConsume("<") {
		delimiter = ">"
	} else {
		if err := tokenizer.Consume("{"); err != nil {
			return err
		}
		delimiter = "}"
	}
	for !tokenizer.LookingAt(">") && !tokenizer.LookingAt("}") {
		if err := textFormat_skipField(tokenizer, recursionLimit); err != nil {
			return err
		}
	}
	return tokenizer.Consume(delimiter)
}

// Java: TextFormat.Parser.skipFieldValue(Tokenizer)
func textFormat_skipFieldValue(tokenizer *TextFormat_Tokenizer) *TextFormat_ParseException {
	if !tokenizer.TryConsumeByteString() &&
		!tokenizer.TryConsumeIdentifier() && // includes enum & boolean
		!tokenizer.TryConsumeInt64() && // includes int32
		!tokenizer.TryConsumeUInt64() && // includes uint32
		!tokenizer.TryConsumeDouble() &&
		!tokenizer.TryConsumeFloat() {
		return tokenizer.ParseException("Invalid field value: " + tokenizer.currentToken)
	}
	return nil
}

// Java: TextFormat.Parser.guessFieldTypeAndSkip(Tokenizer, Descriptor, int)
func textFormat_guessFieldTypeAndSkip(tokenizer *TextFormat_Tokenizer, recursionLimit int) *TextFormat_ParseException {
	colonConsumed := tokenizer.TryConsume(":")
	if tokenizer.LookingAt("[") {
		return textFormat_skipFieldShortFormedRepeated(tokenizer, colonConsumed, recursionLimit)
	} else if colonConsumed && !tokenizer.LookingAt("{") && !tokenizer.LookingAt("<") {
		return textFormat_skipFieldValue(tokenizer)
	}
	if recursionLimit < 1 {
		return tokenizer.ParseException("Message is nested too deep")
	}
	return textFormat_skipFieldMessage(tokenizer, recursionLimit-1)
}

// Java: TextFormat.Parser.skipFieldShortFormedRepeated(Tokenizer, boolean, Descriptor, int)
func textFormat_skipFieldShortFormedRepeated(
	tokenizer *TextFormat_Tokenizer, scalarAllowed bool, recursionLimit int,
) *TextFormat_ParseException {
	if !tokenizer.TryConsume("[") || tokenizer.TryConsume("]") {
		return nil
	}
	for {
		if tokenizer.LookingAt("{") || tokenizer.LookingAt("<") {
			if recursionLimit < 1 {
				return tokenizer.ParseException("Message is nested too deep")
			}
			if err := textFormat_skipFieldMessage(tokenizer, recursionLimit-1); err != nil {
				return err
			}
		} else if scalarAllowed {
			if err := textFormat_skipFieldValue(tokenizer); err != nil {
				return err
			}
		} else {
			return tokenizer.ParseException(
				"Invalid repeated scalar field: missing \":\" before \"[\".")
		}
		if tokenizer.TryConsume("]") {
			break
		}
		if err := tokenizer.Consume(","); err != nil {
			return err
		}
	}
	return nil
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
