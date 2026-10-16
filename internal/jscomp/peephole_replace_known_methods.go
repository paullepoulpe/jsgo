// Port of com/google/javascript/jscomp/PeepholeReplaceKnownMethods.java @ closure-compiler v20260415

package jscomp

import (
	"math"
	"math/bits"
	"strings"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// PeepholeReplaceKnownMethods is just to fold known methods when they are
// called with constants.
type PeepholeReplaceKnownMethods struct {
	AbstractPeepholeOptimizationBase
	late     bool
	useTypes bool
}

// NewPeepholeReplaceKnownMethods creates the optimization.
//
// late: When late is true, this mean we are currently running after most of
// the other optimizations. In this case we avoid changes that make the code
// larger (but otherwise easier to analyze - such as using string splitting).
//
// Java: PeepholeReplaceKnownMethods(boolean, boolean)
func NewPeepholeReplaceKnownMethods(late bool, useTypes bool) *PeepholeReplaceKnownMethods {
	p := &PeepholeReplaceKnownMethods{}
	p.InitAbstractPeepholeOptimizationBase(p) // the implicit super()
	p.late = late
	p.useTypes = useTypes
	return p
}

// Java: PeepholeReplaceKnownMethods.optimizeSubtree(Node)
func (p *PeepholeReplaceKnownMethods) OptimizeSubtree(subtree *rhino.Node) *rhino.Node {
	if NodeUtil_IsGoogWeakUsageCall(subtree) && p.late {
		// goog.weakUsage(x) -> x.
		name := subtree.GetSecondChild().Detach()
		subtree.ReplaceWith(name)
		p.ReportChangeToEnclosingScope(name)
		return name
	} else if subtree.IsCall() {
		return p.TryFoldKnownMethods(subtree)
	}
	return subtree
}

// Java: PeepholeReplaceKnownMethods.tryFoldKnownMethods(Node)
func (p *PeepholeReplaceKnownMethods) TryFoldKnownMethods(subtree *rhino.Node) *rhino.Node {
	// For now we only support string methods .join(),
	// .indexOf(), .substring() and .substr()
	// array method concat()
	// and numeric methods parseInt() and parseFloat().

	jutil.CheckArgumentMsg1(subtree.IsCall(), "%s", subtree)
	subtree = p.TryFoldArrayJoin(subtree)
	// tryFoldArrayJoin may return a string literal instead of a CALL node
	if subtree.IsCall() {
		subtree = p.TryToFoldArrayConcat(subtree)
		jutil.CheckStateMsg1(subtree.IsCall(), "%s", subtree)
		callTarget := jutil.CheckNotNull(subtree.GetFirstChild())

		if callTarget.IsGetProp() {
			if p.IsASTNormalized() && callTarget.GetFirstChild().IsQualifiedName() {
				switch *callTarget.GetFirstChild().GetQualifiedName() {
				case "Array":
					return p.TryFoldKnownArrayMethods(subtree, callTarget)
				case "Math":
					return p.TryFoldKnownMathMethods(subtree, callTarget)
				case "Number":
					return p.TryFoldKnownNumberMethods(subtree, callTarget)
				default:
				}
			}
			subtree = p.TryFoldKnownStringMethods(subtree, callTarget)
		} else if callTarget.IsName() {
			subtree = p.TryFoldKnownNumericMethods(subtree, callTarget)
		}
	}

	return subtree
}

// TryFoldKnownArrayMethods tries to evaluate a method on the Array object.
//
// Java: PeepholeReplaceKnownMethods.tryFoldKnownArrayMethods(Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldKnownArrayMethods(subtree *rhino.Node, callTarget *rhino.Node) *rhino.Node {
	jutil.CheckArgument(subtree.IsCall() && callTarget.IsGetProp())

	// Method node might not be a string if callTarget is a GETELEM.
	// e.g. Array[something]()
	if callTarget.GetString() != "of" {
		return subtree
	}

	subtree.RemoveFirstChild()

	arraylit := rhino.NewNode(rhino.Token_ARRAYLIT)
	arraylit.AddChildrenToBack(subtree.RemoveChildren())
	subtree.ReplaceWith(arraylit)
	p.ReportChangeToEnclosingScope(arraylit)
	return arraylit
}

// Java: PeepholeReplaceKnownMethods.tryFoldKnownNumberMethods(Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldKnownNumberMethods(subtree *rhino.Node, callTarget *rhino.Node) *rhino.Node {
	jutil.CheckArgument(subtree.IsCall() && callTarget.IsGetProp())
	methodName := callTarget.GetString()
	switch methodName {
	case "parseInt", "parseFloat":
		firstArg := callTarget.GetNext()
		if (firstArg != nil && firstArg.IsStringLit()) || p.IsNumericLiteral(firstArg) {
			return p.TryFoldParseNumber(subtree, methodName, firstArg)
		}
		return subtree
	default:
	}

	// If the value is a constant but known not to be a number we could still constant fold many of
	// the methods below but it doesn't seem worth it.
	// Get the only number arg as a double as long as it is the only argument
	var onlyArg *float64
	valueNode := callTarget.GetNext()
	if valueNode == nil || valueNode.GetNext() != nil {
		return subtree
	}
	onlyArg = p.GetSideEffectFreeNumberValueNoConversion(valueNode)
	if onlyArg == nil {
		return subtree
	}
	var replacement *rhino.Node
	switch methodName {
	case "isFinite":
		replacement = NodeUtil_BooleanNode(!math.IsInf(*onlyArg, 0) && !math.IsNaN(*onlyArg))
	case "isNaN":
		replacement = NodeUtil_BooleanNode(math.IsNaN(*onlyArg))
	case "isSafeInteger":
		replacement =
			NodeUtil_BooleanNode(
				*onlyArg < jutil.MathPow(2, 53) &&
					*onlyArg >= (-1*(jutil.MathPow(2, 53)-1)) &&
					float64(jutil.D2L(*onlyArg)) == *onlyArg)
	default:
	}
	if replacement != nil {
		subtree.ReplaceWith(replacement)
		p.ReportChangeToEnclosingScope(replacement)
		return replacement
	}
	return subtree
}

// TryFoldKnownMathMethods tries to evaluate a method on the Math object.
//
// Java: PeepholeReplaceKnownMethods.tryFoldKnownMathMethods(Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldKnownMathMethods(subtree *rhino.Node, callTarget *rhino.Node) *rhino.Node {
	jutil.CheckArgument(subtree.IsCall() && callTarget.IsGetProp())

	// first collect the arguments, if they are all numbers then we proceed
	var args []float64
	for arg := callTarget.GetNext(); arg != nil; arg = arg.GetNext() {
		d := p.GetSideEffectFreeNumberValue(arg)
		if d != nil {
			// lazily allocate, most calls will not be optimizable
			args = append(args, *d)
		} else {
			return subtree
		}
	}
	var replacement *float64
	methodName := callTarget.GetString()
	// NOTE: the standard does not define precision for these methods, but we are conservative, so
	// for now we only implement the methods that are guaranteed to not increase the size of the
	// numeric constants.
	if len(args) == 1 {
		arg := args[0]
		switch methodName {
		case "abs":
			r := math.Abs(arg)
			replacement = &r
		case "ceil":
			r := math.Ceil(arg)
			replacement = &r
		case "floor":
			r := math.Floor(arg)
			replacement = &r
		case "fround":
			if math.IsNaN(arg) || math.IsInf(arg, 0) || arg == 0 {
				r := arg
				replacement = &r
				// if the double is exactly representable as a float, then just cast since no rounding
				// is involved
			} else if float64(float32(arg)) == arg {
				// TODO(b/155511629): This condition is always true after J2CL transpilation.
				r := float64(float32(arg))
				replacement = &r
			} else {
				// (float) arg does not necessarily use the correct rounding mode, so don't do anything
				replacement = nil
			}
		case "round":
			if math.IsNaN(arg) || math.IsInf(arg, 0) {
				r := arg
				replacement = &r
			} else {
				r := float64(jutil.MathRound(arg))
				replacement = &r
			}
		case "sign":
			r := jutil.MathSignum(arg)
			replacement = &r
		case "trunc":
			if math.IsNaN(arg) || math.IsInf(arg, 0) {
				r := arg
				replacement = &r
			} else {
				r := jutil.MathSignum(arg) * math.Floor(math.Abs(arg))
				replacement = &r
			}
		case "clz32":
			r := float64(bits.LeadingZeros32(uint32(base.JSCompDoubles_EcmascriptToUint32(arg))))
			replacement = &r
		default:
		}
	}
	// handle the variadic functions now if we haven't already
	// For each of these we could allow for some of the values to be unknown and either reduce to
	// NaN or simplify the existing args. e.g. Math.max(3, x, 2) -> Math.max(3, x)
	if replacement == nil {
		switch methodName {
		case "max":
			result := math.Inf(-1)
			for _, d := range args {
				result = jutil.MathMax(result, d)
			}
			replacement = &result
		case "min":
			result := math.Inf(1)
			for _, d := range args {
				result = jutil.MathMin(result, d)
			}
			replacement = &result
		case "imul":
			if len(args) < 2 {
				r := 0.0
				replacement = &r
			} else {
				// Ignore args3+
				r := float64(base.JSCompDoubles_EcmascriptToInt32(args[0]) * base.JSCompDoubles_EcmascriptToInt32(args[1]))
				replacement = &r
			}
		case "pow":
			if len(args) == 2 {
				r := jutil.MathPow(args[0], args[1])
				replacement = &r
			}
		default:
		}
	}

	if replacement != nil {
		numberNode := NodeUtil_NumberNode(*replacement, subtree)
		subtree.ReplaceWith(numberNode)
		p.ReportChangeToEnclosingScope(numberNode)
		return numberNode
	}
	return subtree
}

// TryFoldKnownStringMethods tries to evaluate known String methods
// .indexOf(), .substr(), .substring()
//
// Java: PeepholeReplaceKnownMethods.tryFoldKnownStringMethods(Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldKnownStringMethods(subtree *rhino.Node, callTarget *rhino.Node) *rhino.Node {
	jutil.CheckArgument(subtree.IsCall() && callTarget.IsGetProp())

	// check if this is a call on a string method
	// then dispatch to specific folding method.
	stringNode := callTarget.GetFirstChild()

	isStringLiteral := stringNode.IsStringLit()
	functionNameString := callTarget.GetString()
	firstArg := callTarget.GetNext()
	if isStringLiteral {
		if functionNameString == "split" {
			return p.TryFoldStringSplit(subtree, stringNode, firstArg)
		} else if firstArg == nil {
			switch functionNameString {
			case "toLowerCase":
				return p.TryFoldStringToLowerCase(subtree, stringNode)
			case "toUpperCase":
				return p.TryFoldStringToUpperCase(subtree, stringNode)
			case "trim":
				return p.TryFoldStringTrim(subtree, stringNode)
			default:
			}
		} else {
			if NodeUtil_IsImmutableValue(firstArg) {
				switch functionNameString {
				case "indexOf", "lastIndexOf":
					return p.TryFoldStringIndexOf(subtree, functionNameString, stringNode, firstArg)
				case "substr":
					return p.TryFoldStringSubstr(subtree, stringNode, firstArg)
				case "substring", "slice":
					return p.TryFoldStringSubstringOrSlice(subtree, stringNode, firstArg)
				case "charAt":
					return p.TryFoldStringCharAt(subtree, stringNode, firstArg)
				case "charCodeAt":
					return p.TryFoldStringCharCodeAt(subtree, stringNode, firstArg)
				case "replace":
					return p.TryFoldStringReplace(subtree, stringNode, firstArg)
				case "replaceAll":
					return p.TryFoldStringReplaceAll(subtree, stringNode, firstArg)
				default:
				}
			}
		}
	}
	if p.useTypes &&
		firstArg != nil &&
		(isStringLiteral || colors.StandardColors_STRING.Equals(stringNode.GetColor())) {
		if subtree.HasXChildren(3) {
			maybeStart := p.GetSideEffectFreeNumberValue(firstArg)
			if maybeStart != nil {
				start := int(jutil.D2I(*maybeStart))
				maybeLengthOrEnd := p.GetSideEffectFreeNumberValue(firstArg.GetNext())
				if maybeLengthOrEnd != nil {
					switch functionNameString {
					case "substr":
						length := int(jutil.D2I(*maybeLengthOrEnd))
						if start >= 0 && length == 1 {
							return p.ReplaceWithCharAt(subtree, callTarget, firstArg)
						}
					case "substring", "slice":
						end := int(jutil.D2I(*maybeLengthOrEnd))
						// unlike slice and substring, chatAt can not be used with negative indexes
						if start >= 0 && jutil.IntSub(end, start) == 1 {
							return p.ReplaceWithCharAt(subtree, callTarget, firstArg)
						}
					default:
					}
				}
			}
		}
	}
	return subtree
}

// TryFoldKnownNumericMethods tries to evaluate known Numeric methods
// parseInt(), parseFloat()
//
// Java: PeepholeReplaceKnownMethods.tryFoldKnownNumericMethods(Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldKnownNumericMethods(subtree *rhino.Node, callTarget *rhino.Node) *rhino.Node {
	jutil.CheckArgument(subtree.IsCall())

	if p.IsASTNormalized() {
		// check if this is a call on a string method
		// then dispatch to specific folding method.
		functionNameString := callTarget.GetString()
		firstArgument := callTarget.GetNext()
		if (firstArgument != nil) &&
			(firstArgument.IsStringLit() || p.IsNumericLiteral(firstArgument)) &&
			(functionNameString == "parseInt" || functionNameString == "parseFloat") {
			subtree = p.TryFoldParseNumber(subtree, functionNameString, firstArgument)
		}
	}
	return subtree
}

// IsNumericLiteral returns true both for number literals and their negations
// (e.g. `-12.3`).
//
// Java: PeepholeReplaceKnownMethods.isNumericLiteral(Node)
func (p *PeepholeReplaceKnownMethods) IsNumericLiteral(n *rhino.Node) bool {
	return n.IsNumber() || (n.IsNeg() && n.GetOnlyChild().IsNumber())
}

// TryFoldStringToLowerCase returns The lowered string Node.
//
// This method is believed to be correct independent of the locale of the
// compiler and the JSVM executing the compiled code, assuming both are
// implementations of Unicode are correct.
//
// See https://tc39.es/ecma262/#sec-string.prototype.tolowercase,
// https://unicode.org/faq/casemap_charprop.html#5 and
// https://docs.oracle.com/javase/8/docs/api/java/lang/String.html#toLowerCase-java.util.Locale-
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringToLowerCase(Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringToLowerCase(subtree *rhino.Node, stringNode *rhino.Node) *rhino.Node {
	lowered := jutil.StringToLowerCaseRoot(stringNode.GetString())
	replacement := rhino.IR_String(lowered)
	subtree.ReplaceWith(replacement)
	p.ReportChangeToEnclosingScope(replacement)
	return replacement
}

// TryFoldStringToUpperCase returns The upped string Node.
//
// This method is believed to be correct independent of the locale of the
// compiler and the JSVM executing the compiled code, assuming both are
// implementations of Unicode are correct.
//
// See https://tc39.es/ecma262/#sec-string.prototype.touppercase,
// https://unicode.org/faq/casemap_charprop.html#5 and
// https://docs.oracle.com/javase/8/docs/api/java/lang/String.html#toUpperCase-java.util.Locale-
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringToUpperCase(Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringToUpperCase(subtree *rhino.Node, stringNode *rhino.Node) *rhino.Node {
	upped := jutil.StringToUpperCaseRoot(stringNode.GetString())
	replacement := rhino.IR_String(upped)
	subtree.ReplaceWith(replacement)
	p.ReportChangeToEnclosingScope(replacement)
	return replacement
}

// TryFoldStringTrim returns the trimmed string Node.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringTrim(Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringTrim(subtree *rhino.Node, stringNode *rhino.Node) *rhino.Node {
	// See ECMA 15.5.4.20, 7.2, and 7.3
	// All Unicode 10.0 whitespace + BOM
	// whitespace :=
	//     "[ \t\n-\r\\u0085\\u00A0\\u1680\\u2000-\\u200A\\u2028\\u2029\\u202F\\u205F\\u3000\\uFEFF]+"
	// trimmed := s.replaceAll("^" + whitespace + "|" + whitespace + "$", "")
	trimmed := peepholeReplaceKnownMethods_trimWhitespaceRegex(stringNode.GetString())
	replacement := rhino.IR_String(trimmed)
	subtree.ReplaceWith(replacement)
	p.ReportChangeToEnclosingScope(replacement)
	return replacement
}

// peepholeReplaceKnownMethods_isTrimWhitespace is the character class of
// tryFoldStringTrim's regex: [ \t\n-\r\u0085\u00A0\u1680\u2000-\u200A
// \u2028\u2029\u202F\u205F\u3000\uFEFF].
func peepholeReplaceKnownMethods_isTrimWhitespace(r rune) bool {
	switch {
	case r == ' ' || r == '\t' || (r >= '\n' && r <= '\r'):
		return true
	case r == 0x0085 || r == 0x00A0 || r == 0x1680 || (r >= 0x2000 && r <= 0x200A):
		return true
	case r == 0x2028 || r == 0x2029 || r == 0x202F || r == 0x205F || r == 0x3000 || r == 0xFEFF:
		return true
	}
	return false
}

// peepholeReplaceKnownMethods_trimWhitespaceRegex is
// s.replaceAll("^" + whitespace + "|" + whitespace + "$", "") for
// tryFoldStringTrim's whitespace class W.
//
// Go-only: language-adapter — Java `$` and String.replaceAll (jutil.Pattern
// supports neither: Java's `$` also matches before a final line terminator),
// so this is an equivalent hand matcher for this one java.util.regex use
// (PORTING.md, Regex). The first alternative removes the maximal leading run of
// W. The second only matches a run of W that reaches the end of the input:
// Java's `$` could also match before a final line terminator, but every line
// terminator is in W, so the greedy run always extends to the end first; and
// a run followed by a non-W character can never be followed by `$`. Every
// character of W is in the BMP and is not a surrogate, so matching code
// points of the WTF-8 string is matching Java's.
func peepholeReplaceKnownMethods_trimWhitespaceRegex(s string) string {
	start := 0
	for start < len(s) {
		r, size := utf8.DecodeRuneInString(s[start:])
		if !peepholeReplaceKnownMethods_isTrimWhitespace(r) {
			break
		}
		start += size
	}
	end := len(s)
	for end > start {
		r, size := utf8.DecodeLastRuneInString(s[start:end])
		if !peepholeReplaceKnownMethods_isTrimWhitespace(r) {
			break
		}
		end -= size
	}
	return s[start:end]
}

// PeepholeReplaceKnownMethods_NormalizeNumericString takes the string
// representation of a number (input) and returns the string with leading and
// trailing zeros removed.
//
// Java: PeepholeReplaceKnownMethods.normalizeNumericString(String)
func PeepholeReplaceKnownMethods_NormalizeNumericString(input string) string {
	if input == "" {
		return input
	}

	startIndex := 0
	endIndex := jutil.StringLength(input) - 1

	// Remove leading zeros
	for startIndex < jutil.StringLength(input) &&
		jutil.StringCharAt(input, startIndex) == '0' &&
		jutil.StringCharAt(input, startIndex) != '.' {
		startIndex++
	}

	// Remove trailing zeros only after the decimal
	if jutil.StringIndexOfChar(input, '.') >= 0 {
		for endIndex >= 0 && jutil.StringCharAt(input, endIndex) == '0' {
			endIndex--
		}
		if jutil.StringCharAt(input, endIndex) == '.' {
			endIndex--
		}
	}
	if startIndex >= endIndex {
		return input
	}

	return jutil.StringSubstring(input, startIndex, endIndex+1)
}

// TryFoldParseNumber tries to evaluate parseInt, parseFloat:
//
//	parseInt("1") -> 1
//	parseInt("1", 10) -> 1
//	parseFloat("1.11") -> 1.11
//
// Java: PeepholeReplaceKnownMethods.tryFoldParseNumber(Node, String, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldParseNumber(n *rhino.Node, functionName string, firstArg *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsCall())

	isParseInt := functionName == "parseInt"
	secondArg := firstArg.GetNext()

	// Second argument is only used as the radix for parseInt
	radix := 0
	if secondArg != nil {
		if !isParseInt {
			return n
		}

		// Third-argument and non-numeric second arg are problematic. Discard.
		if secondArg.GetNext() != nil || !secondArg.IsNumber() {
			return n
		} else {
			tmpRadix := secondArg.GetDouble()
			if tmpRadix != float64(jutil.D2I(tmpRadix)) {
				return n
			}
			radix = int(jutil.D2I(tmpRadix))
			if radix < 0 || radix == 1 || radix > 36 {
				return n
			}
		}
	}

	// stringVal must be a valid string.
	var stringVal string
	var checkVal *float64
	if p.IsNumericLiteral(firstArg) {
		checkVal = p.GetSideEffectFreeNumberValue(firstArg)
		if !(radix == 0 || radix == 10) && isParseInt {
			// Convert a numeric first argument to a different base
			stringVal = jutil.LongToString(int64(jutil.D2I(*checkVal)))
		} else {
			// If parseFloat is called with a numeric argument,
			// replace it with just the number.
			// If parseInt is called with a numeric first argument and the radix
			// is 10 or omitted, just replace it with the number
			var numericNode *rhino.Node
			if isParseInt {
				numericNode = NodeUtil_NumberNode(float64(jutil.D2I(*checkVal)), n)
			} else {
				numericNode = NodeUtil_NumberNode(*checkVal, n)
			}
			n.ReplaceWith(numericNode)
			p.ReportChangeToEnclosingScope(numericNode)
			return numericNode
		}
	} else {
		maybeStringVal := p.GetSideEffectFreeStringValue(firstArg)
		if maybeStringVal == nil {
			return n
		}
		stringVal = *maybeStringVal

		// Check that the string is in a format we can recognize
		checkVal = NodeUtil_GetStringNumberValue(stringVal)
		if checkVal == nil {
			return n
		}

		stringVal = NodeUtil_TrimJsWhiteSpace(stringVal)
		if stringVal == "" {
			return n
		}
	}

	var newNode *rhino.Node
	if stringVal == "0" {
		// Special case for parseInt("0") or parseFloat("0")
		newNode = rhino.IR_Number(0)
	} else if isParseInt {
		if radix == 0 || radix == 16 {
			if jutil.StringLength(stringVal) > 1 && jutil.AsciiEqualsIgnoreCase(jutil.StringSubstring(stringVal, 0, 2), "0x") {
				radix = 16
				stringVal = jutil.StringSubstringFrom(stringVal, 2)
			} else if radix == 0 {
				// if a radix is not specified or is 0 and the most
				// significant digit is "0", the string will parse
				// with a radix of 8 on some browsers, so leave
				// this case alone. This check does not apply in
				// script mode ECMA5 or greater
				if !p.IsEcmaScript5OrGreater() && jutil.StringSubstring(stringVal, 0, 1) == "0" {
					return n
				}

				radix = 10
			}
		}
		var newVal int32 = 0
		failed := false
		jutil.Try(func() {
			newVal = jutil.IntegerParseInt(stringVal, radix)
		}, func(e jutil.INumberFormatException) {
			failed = true
		})
		if failed {
			return n
		}
		if newVal == 0 && jutil.StringStartsWith(stringVal, "-") {
			newNode = NodeUtil_NumberNode(math.Copysign(0, -1), n)
		} else {
			newNode = NodeUtil_NumberNode(float64(newVal), n)
		}
	} else {
		normalizedNewVal := "0"
		failed := false
		jutil.Try(func() {
			newVal := jutil.DoubleParseDouble(stringVal)
			newNode = NodeUtil_NumberNode(newVal, n)
			normalizedNewVal = PeepholeReplaceKnownMethods_NormalizeNumericString(jutil.DoubleToString(newVal))
		}, func(e jutil.INumberFormatException) {
			failed = true
		})
		if failed {
			return n
		}
		// Make sure that the parsed number matches the original string
		// This prevents rounding differences between the Java implementation
		// and native script.
		if PeepholeReplaceKnownMethods_NormalizeNumericString(stringVal) != normalizedNewVal {
			return n
		}
	}

	n.ReplaceWith(newNode)
	p.ReportChangeToEnclosingScope(newNode)

	return newNode
}

// TryFoldStringIndexOf tries to evaluate String.indexOf/lastIndexOf:
//
//	"abcdef".indexOf("bc") -> 1
//	"abcdefbc".indexOf("bc", 3) -> 6
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringIndexOf(Node, String, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringIndexOf(n *rhino.Node, functionName string, lstringNode *rhino.Node, firstArg *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsCall())
	jutil.CheckArgument(lstringNode.IsStringLit())

	lstring := lstringNode.GetString()
	isIndexOf := functionName == "indexOf"
	secondArg := firstArg.GetNext()
	searchValue := p.GetSideEffectFreeStringValue(firstArg)
	// searchValue must be a valid string.
	if searchValue == nil {
		return n
	}
	fromIndex := 0
	if !isIndexOf {
		fromIndex = jutil.StringLength(lstring)
	}
	if secondArg != nil {
		// Third-argument and non-numeric second arg are problematic. Discard.
		if secondArg.GetNext() != nil || !secondArg.IsNumber() {
			return n
		} else {
			fromIndex = int(jutil.D2I(secondArg.GetDouble()))
		}
	}
	var indexVal int
	if isIndexOf {
		indexVal = jutil.StringIndexOfFrom(lstring, *searchValue, fromIndex)
	} else {
		indexVal = jutil.StringLastIndexOfFrom(lstring, *searchValue, fromIndex)
	}
	newNode := NodeUtil_NumberNode(float64(indexVal), n)
	n.ReplaceWith(newNode)
	p.ReportChangeToEnclosingScope(newNode)

	return newNode
}

// TryFoldArrayJoin tries to fold an array join: ['a', 'b', 'c'].join(”) ->
// 'abc';
//
// Java: PeepholeReplaceKnownMethods.tryFoldArrayJoin(Node)
func (p *PeepholeReplaceKnownMethods) TryFoldArrayJoin(n *rhino.Node) *rhino.Node {
	jutil.CheckStateMsg1(n.IsCall(), "%s", n)
	callTarget := n.GetFirstChild()

	if callTarget == nil || !callTarget.IsGetProp() {
		return n
	}

	right := callTarget.GetNext()
	if right != nil {
		if right.GetNext() != nil || !NodeUtil_IsImmutableValue(right) {
			return n
		}
	}

	arrayNode := callTarget.GetFirstChild()

	if !arrayNode.IsArrayLit() || callTarget.GetString() != "join" {
		return n
	}

	if right != nil && right.IsStringLit() && "," == right.GetString() {
		// "," is the default, it doesn't need to be explicit
		right.Detach()
		p.ReportChangeToEnclosingScope(n)
	}

	// logic above ensures that `right` is immutable, so no need to check for
	// side effects with getSideEffectFreeStringValue(right)
	var joinString *string
	if right == nil {
		comma := ","
		joinString = &comma
	} else {
		joinString = NodeUtil_GetStringValue(right)
	}
	var arrayFoldedChildren []*rhino.Node
	var sb *jutil.StringBuilder
	foldedSize := 0
	var prev *rhino.Node
	elem := arrayNode.GetFirstChild()
	// Merges adjacent String nodes.
	for elem != nil {
		if NodeUtil_IsImmutableValue(elem) || elem.IsEmpty() {
			if sb == nil {
				sb = jutil.NewStringBuilder()
			} else {
				sb.Append(jutil.ToJavaString(joinString))
			}
			elementStr := NodeUtil_GetArrayElementStringValue(elem)
			if elementStr == nil {
				return n // TODO(nickreid): Is this ever null?
			}
			sb.Append(*elementStr)
		} else {
			if sb != nil {
				jutil.CheckNotNull(prev)
				// + 2 for the quotes.
				foldedSize += sb.Length() + 2
				arrayFoldedChildren = append(arrayFoldedChildren, rhino.IR_String(sb.String()).SrcrefIfMissing(prev))
				sb = nil
			}
			foldedSize += InlineCostEstimator_GetCost(elem)
			arrayFoldedChildren = append(arrayFoldedChildren, elem)
		}
		prev = elem
		elem = elem.GetNext()
	}

	if sb != nil {
		jutil.CheckNotNull(prev)
		// + 2 for the quotes.
		foldedSize += sb.Length() + 2
		arrayFoldedChildren = append(arrayFoldedChildren, rhino.IR_String(sb.String()).SrcrefIfMissing(prev))
	}
	// one for each comma.
	foldedSize += len(arrayFoldedChildren) - 1

	originalSize := InlineCostEstimator_GetCost(n)
	switch len(arrayFoldedChildren) {
	case 0:
		emptyStringNode := rhino.IR_String("")
		n.ReplaceWith(emptyStringNode)
		p.ReportChangeToEnclosingScope(emptyStringNode)
		return emptyStringNode
	case 1:
		foldedStringNode := arrayFoldedChildren[0]
		arrayFoldedChildren = arrayFoldedChildren[1:]
		// The spread isn't valid outside any array literal (or would change meaning)
		// so don't try to fold it.
		if foldedStringNode.IsSpread() || foldedSize > originalSize {
			return n
		}
		if foldedStringNode.IsStringLit() {
			arrayNode.DetachChildren()
			n.ReplaceWith(foldedStringNode)
			p.ReportChangeToEnclosingScope(foldedStringNode)
			return foldedStringNode
		} else {
			// Because of special case behavior for `null` and `undefined` values, there's no safe way
			// to convert `[someNonStringValue].join()` to something shorter.
			// e.g. String(someNonStringValue) would turn `null` into `"null"`, which isn't right.
			return n
		}
	default:
		if arrayNode.HasXChildren(len(arrayFoldedChildren)) {
			// No folding could actually be performed.
			return n
		}
		kJoinOverhead := len("[].join()")
		foldedSize += kJoinOverhead
		if right != nil {
			foldedSize += InlineCostEstimator_GetCost(right)
		}
		if foldedSize > originalSize {
			return n
		}
		arrayNode.DetachChildren()
		for _, node := range arrayFoldedChildren {
			arrayNode.AddChildToBack(node)
		}
		p.ReportChangeToEnclosingScope(arrayNode)
	}

	return n
}

// TryFoldStringSubstr tries to fold .substr() calls on strings.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringSubstr(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringSubstr(n *rhino.Node, stringNode *rhino.Node, arg1 *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsCall())
	jutil.CheckArgument(stringNode.IsStringLit())
	jutil.CheckArgument(arg1 != nil)

	var start int
	var length int
	stringAsString := stringNode.GetString()

	maybeStart := p.GetSideEffectFreeNumberValue(arg1)
	if maybeStart != nil {
		start = int(jutil.D2I(*maybeStart))
	} else {
		return n
	}

	arg2 := arg1.GetNext()
	if arg2 != nil {
		maybeLength := p.GetSideEffectFreeNumberValue(arg2)
		if maybeLength != nil {
			length = int(jutil.D2I(*maybeLength))
		} else {
			return n
		}

		if arg2.GetNext() != nil {
			// If we got more args than we expected, bail out.
			return n
		}
	} else {
		// parameter 2 not passed
		length = jutil.IntSub(jutil.StringLength(stringAsString), start)
	}

	// Don't handle these cases. The specification actually does
	// specify the behavior in some of these cases, but we haven't
	// done a thorough investigation that it is correctly implemented
	// in all browsers.
	if jutil.IntAdd(start, length) > jutil.StringLength(stringAsString) || (length < 0) || (start < 0) {
		return n
	}

	result := jutil.StringSubstring(stringAsString, start, jutil.IntAdd(start, length))
	resultNode := rhino.IR_String(result)

	parent := n.GetParent()
	n.ReplaceWith(resultNode)
	p.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// TryFoldStringSubstringOrSlice tries to fold .substring() or .slice() calls
// on strings.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringSubstringOrSlice(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringSubstringOrSlice(n *rhino.Node, stringNode *rhino.Node, arg1 *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsCall())
	jutil.CheckArgument(stringNode.IsStringLit())
	jutil.CheckArgument(arg1 != nil)

	var start int
	var end int
	stringAsString := stringNode.GetString()

	maybeStart := p.GetSideEffectFreeNumberValue(arg1)
	if maybeStart != nil {
		start = int(jutil.D2I(*maybeStart))
	} else {
		return n
	}

	arg2 := arg1.GetNext()
	if arg2 != nil {
		maybeEnd := p.GetSideEffectFreeNumberValue(arg2)
		if maybeEnd != nil {
			end = int(jutil.D2I(*maybeEnd))
		} else {
			return n
		}

		if arg2.GetNext() != nil {
			// If we got more args than we expected, bail out.
			return n
		}
	} else {
		// parameter 2 not passed
		end = jutil.StringLength(stringAsString)
	}

	// Don't handle these cases. The specification actually does
	// specify the behavior in some of these cases, but we haven't
	// done a thorough investigation that it is correctly implemented
	// in all browsers.
	if (end > jutil.StringLength(stringAsString)) ||
		(start > jutil.StringLength(stringAsString)) ||
		(start < 0) ||
		(end < 0) ||
		(start > end) {
		return n
	}

	result := jutil.StringSubstring(stringAsString, start, end)
	resultNode := rhino.IR_String(result)

	parent := n.GetParent()
	n.ReplaceWith(resultNode)
	p.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// Java: PeepholeReplaceKnownMethods.replaceWithCharAt(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) ReplaceWithCharAt(n *rhino.Node, callTarget *rhino.Node, firstArg *rhino.Node) *rhino.Node {
	// TODO(moz): Maybe correct the arity of the function type here.
	callTarget.SetString("charAt")
	firstArg.GetNext().Detach()
	p.ReportChangeToEnclosingScope(firstArg)
	return n
}

// TryFoldStringCharAt tries to fold .charAt() calls on strings.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringCharAt(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringCharAt(n *rhino.Node, stringNode *rhino.Node, arg1 *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsCall())
	jutil.CheckArgument(stringNode.IsStringLit())

	var index int
	stringAsString := stringNode.GetString()

	if arg1 != nil && arg1.IsNumber() && arg1.GetNext() == nil {
		index = int(jutil.D2I(arg1.GetDouble()))
	} else {
		return n
	}

	if index < 0 || jutil.StringLength(stringAsString) <= index {
		// http://es5.github.com/#x15.5.4.4 says "" is returned when index is
		// out of bounds but we bail.
		return n
	}

	resultNode := rhino.IR_String(jutil.StringSubstring(stringAsString, index, index+1))
	parent := n.GetParent()
	n.ReplaceWith(resultNode)
	p.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// TryFoldStringReplace tries to fold .replace() calls on strings.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringReplace(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringReplace(n *rhino.Node, stringNode *rhino.Node, arg1 *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsCall())
	jutil.CheckArgument(stringNode.IsStringLit())

	arg2 := arg1.GetNext()
	if arg2 == nil || arg2.GetNext() != nil {
		// too few or too many parameters
		return n
	}

	if !arg1.IsStringLit() || !arg2.IsStringLit() {
		// only string literals are supported for folding.
		return n
	}

	lookForPattern := arg1.GetString()
	replacementPattern := arg2.GetString()
	if strings.Contains(replacementPattern, "$") {
		// 'special' replacements aren't supported yet.
		return n
	}

	original := stringNode.GetString()

	index := jutil.StringIndexOf(original, lookForPattern)
	if index == -1 {
		return n
	}

	// Java "replace" acts like JavaScript's "replaceAll" here we only want to replace the first
	// instance of the string
	newString :=
		jutil.StringSubstring(original, 0, index) +
			replacementPattern +
			jutil.StringSubstringFrom(original, index+jutil.StringLength(lookForPattern))

	resultNode := rhino.IR_String(newString).Srcref(stringNode)
	parent := n.GetParent()
	n.ReplaceWith(resultNode)
	p.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// TryFoldStringReplaceAll tries to fold .replaceAll() calls on strings.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringReplaceAll(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringReplaceAll(n *rhino.Node, stringNode *rhino.Node, arg1 *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsCall())
	jutil.CheckArgument(stringNode.IsStringLit())

	arg2 := arg1.GetNext()
	if arg2 == nil || arg2.GetNext() != nil {
		// too few or too many parameters
		return n
	}

	if !arg1.IsStringLit() || !arg2.IsStringLit() {
		// only string literals are supported for folding.
		return n
	}

	replacementPattern := arg2.GetString()
	if strings.Contains(replacementPattern, "$") {
		// 'special' replacements aren't supported yet.
		return n
	}

	// Java "replace" acts like JavaScript's "replaceAll" and replaces all occurrences.
	original := stringNode.GetString()
	newString := jutil.StringReplace(original, arg1.GetString(), replacementPattern)

	resultNode := rhino.IR_String(newString).Srcref(stringNode)
	parent := n.GetParent()
	n.ReplaceWith(resultNode)
	p.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// TryFoldStringCharCodeAt tries to fold .charCodeAt() calls on strings.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringCharCodeAt(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringCharCodeAt(n *rhino.Node, stringNode *rhino.Node, arg1 *rhino.Node) *rhino.Node {
	jutil.CheckArgument(n.IsCall())
	jutil.CheckArgument(stringNode.IsStringLit())

	var index int
	stringAsString := stringNode.GetString()

	if arg1 != nil && arg1.IsNumber() && arg1.GetNext() == nil {
		index = int(jutil.D2I(arg1.GetDouble()))
	} else {
		return n
	}

	if index < 0 || jutil.StringLength(stringAsString) <= index {
		// http://es5.github.com/#x15.5.4.5 says NaN is returned when index is
		// out of bounds but we bail.
		return n
	}

	resultNode := rhino.IR_Number(float64(jutil.StringCharAt(stringAsString, index)))
	parent := n.GetParent()
	n.ReplaceWith(resultNode)
	p.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// PeepholeReplaceKnownMethods_JsSplitMatch is a support function for jsSplit:
// it finds the first occurrence of separator within stringValue starting at
// startIndex.
//
// Java: PeepholeReplaceKnownMethods.jsSplitMatch(String, int, String)
func PeepholeReplaceKnownMethods_JsSplitMatch(stringValue string, startIndex int, separator string) int {

	if startIndex+jutil.StringLength(separator) > jutil.StringLength(stringValue) {
		return -1
	}

	matchIndex := jutil.StringIndexOfFrom(stringValue, separator, startIndex)

	if matchIndex < 0 {
		return -1
	}

	return matchIndex
}

// JsSplit implements the JS String.split method using a string separator.
//
// Java: PeepholeReplaceKnownMethods.jsSplit(String, String, int)
func (p *PeepholeReplaceKnownMethods) JsSplit(stringValue string, separator *string, limit int) []string {
	jutil.CheckArgument(limit >= 0)

	// For limits of 0, return an empty array
	if limit == 0 {
		return []string{}
	}

	// If a separator is not specified, return the entire string as
	// the only element of an array.
	if separator == nil {
		return []string{stringValue}
	}

	var splitStrings []string

	// If an empty string is specified for the separator, split apart each
	// character of the string.
	if *separator == "" {
		for i := 0; i < jutil.StringLength(stringValue) && i < limit; i++ {
			splitStrings = append(splitStrings, jutil.StringSubstring(stringValue, i, i+1))
		}
	} else {
		startIndex := 0
		var matchIndex int
		for {
			matchIndex = PeepholeReplaceKnownMethods_JsSplitMatch(stringValue, startIndex, *separator)
			if !(matchIndex >= 0 && len(splitStrings) < limit) {
				break
			}
			splitStrings = append(splitStrings, jutil.StringSubstring(stringValue, startIndex, matchIndex))

			startIndex = matchIndex + jutil.StringLength(*separator)
		}

		if len(splitStrings) < limit {
			if startIndex < jutil.StringLength(stringValue) {
				splitStrings = append(splitStrings, jutil.StringSubstringFrom(stringValue, startIndex))
			} else {
				splitStrings = append(splitStrings, "")
			}
		}
	}

	return splitStrings
}

// TryFoldStringSplit tries to fold .split() calls on strings.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringSplit(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringSplit(n *rhino.Node, stringNode *rhino.Node, arg1 *rhino.Node) *rhino.Node {
	if p.late {
		return n
	}

	jutil.CheckArgument(n.IsCall())
	jutil.CheckArgument(stringNode.IsStringLit())

	var separator *string
	stringValue := stringNode.GetString()

	// Maximum number of possible splits
	limit := jutil.StringLength(stringValue) + 1

	if arg1 != nil {
		if arg1.IsStringLit() {
			s := arg1.GetString()
			separator = &s
		} else if !arg1.IsNull() {
			return n
		}

		arg2 := arg1.GetNext()
		if arg2 != nil {
			if arg2.IsNumber() {
				limit = min(int(jutil.D2I(arg2.GetDouble())), limit)
				if limit < 0 {
					return n
				}
			} else {
				return n
			}
		}
	}

	// Split the string and convert the returned array into JS nodes
	stringArray := p.JsSplit(stringValue, separator, limit)
	arrayOfStrings := rhino.IR_Arraylit()
	for _, element := range stringArray {
		arrayOfStrings.AddChildToBack(rhino.IR_String(element).Srcref(stringNode))
	}

	parent := n.GetParent()
	n.ReplaceWith(arrayOfStrings)
	p.ReportChangeToEnclosingScope(parent)
	return arrayOfStrings
}

// Java: PeepholeReplaceKnownMethods.tryToFoldArrayConcat(Node)
func (p *PeepholeReplaceKnownMethods) TryToFoldArrayConcat(n *rhino.Node) *rhino.Node {
	jutil.CheckArgumentMsg1(n.IsCall(), "%s", n)

	if !p.IsASTNormalized() || !p.useTypes {
		return n
	}
	concatFunctionCall := PeepholeReplaceKnownMethods_CreateConcatFunctionCallForNode(n)
	if concatFunctionCall == nil {
		return n
	}
	concatFunctionCall = p.TryToRemoveArrayLiteralFromFrontOfConcat(concatFunctionCall)
	jutil.CheckNotNull(concatFunctionCall)
	return p.TryToFoldConcatChaining(concatFunctionCall)
}

// TryToRemoveArrayLiteralFromFrontOfConcat checks if we have this code
// pattern `[].concat(exactlyArrayArgument,...*)` and if yes replace empty
// array literal from the front of concatenation by the first argument of
// concat function call `[].concat(arr,1)` -> `arr.concat(1)`.
//
// Java: PeepholeReplaceKnownMethods.tryToRemoveArrayLiteralFromFrontOfConcat(ConcatFunctionCall)
func (p *PeepholeReplaceKnownMethods) TryToRemoveArrayLiteralFromFrontOfConcat(
	concatFunctionCall *peepholeReplaceKnownMethods_ConcatFunctionCall) *peepholeReplaceKnownMethods_ConcatFunctionCall {
	jutil.CheckNotNull(concatFunctionCall)

	callNode := concatFunctionCall.callNode
	arrayLiteralToRemove := concatFunctionCall.calleeNode
	if !arrayLiteralToRemove.IsArrayLit() || arrayLiteralToRemove.HasChildren() {
		return concatFunctionCall
	}
	firstArg := concatFunctionCall.firstArgumentNode
	if !PeepholeReplaceKnownMethods_ContainsExactlyArray(firstArg) {
		return concatFunctionCall
	}

	firstArg.Detach()
	arrayLiteralToRemove.ReplaceWith(firstArg)

	p.ReportChangeToEnclosingScope(callNode)
	return PeepholeReplaceKnownMethods_CreateConcatFunctionCallForNode(callNode)
}

// TryToFoldConcatChaining checks if we have this code pattern
// `array.concat(...*).concat(sideEffectFreeArguments)` and if yes fold
// chained concat functions, so `arr.concat(a).concat(b)` will be fold into
// `arr.concat(a,b)`.
//
// Java: PeepholeReplaceKnownMethods.tryToFoldConcatChaining(ConcatFunctionCall)
func (p *PeepholeReplaceKnownMethods) TryToFoldConcatChaining(concatFunctionCall *peepholeReplaceKnownMethods_ConcatFunctionCall) *rhino.Node {
	jutil.CheckNotNull(concatFunctionCall)

	concatCallNode := concatFunctionCall.callNode

	maybeFunctionCall := concatFunctionCall.calleeNode
	if !maybeFunctionCall.IsCall() {
		return concatCallNode
	}
	previousConcatFunctionCall :=
		PeepholeReplaceKnownMethods_CreateConcatFunctionCallForNode(maybeFunctionCall)
	if previousConcatFunctionCall == nil {
		return concatCallNode
	}
	// make sure that arguments in second concat function call can't change the array
	// so we can fold chained concat functions
	// to clarify, consider this code
	// here we can't fold concatenation
	// var a = [];
	// a.concat(1).concat(a.push(1)); -> [1,1]
	// a.concat(1,a.push(1)); -> [1,1,1]
	for arg := concatFunctionCall.firstArgumentNode; arg != nil; arg = arg.GetNext() {
		if p.MayHaveSideEffects(arg) {
			return concatCallNode
		}
	}

	// perform folding
	previousConcatCallNode := previousConcatFunctionCall.callNode
	arg := concatFunctionCall.firstArgumentNode
	for arg != nil {
		currentArg := arg
		arg = arg.GetNext()
		previousConcatCallNode.AddChildToBack(currentArg.Detach())
	}
	concatCallNode.ReplaceWith(previousConcatCallNode.Detach())
	p.ReportChangeToEnclosingScope(previousConcatCallNode)
	return previousConcatCallNode
}

// peepholeReplaceKnownMethods_ConcatFunctionCall is the ConcatFunctionCall
// class.
//
// Go-only: language-adapter — the Java class is abstract only so that it is
// instantiated as an empty anonymous subclass (`new ConcatFunctionCall(...)
// {}`), which adds nothing; one struct stands for both.
type peepholeReplaceKnownMethods_ConcatFunctionCall struct {
	callNode          *rhino.Node
	calleeNode        *rhino.Node
	firstArgumentNode *rhino.Node
}

// Java: PeepholeReplaceKnownMethods.ConcatFunctionCall(Node, Node, Node)
func newPeepholeReplaceKnownMethods_ConcatFunctionCall(callNode *rhino.Node, calleeNode *rhino.Node, firstArgumentNode *rhino.Node) *peepholeReplaceKnownMethods_ConcatFunctionCall {
	c := &peepholeReplaceKnownMethods_ConcatFunctionCall{}
	c.callNode = jutil.CheckNotNull(callNode)
	c.calleeNode = jutil.CheckNotNull(calleeNode)
	c.firstArgumentNode = firstArgumentNode
	return c
}

// PeepholeReplaceKnownMethods_CreateConcatFunctionCallForNode returns, if the
// argument node is a call to `Array.prototype.concat`, a
// `ConcatFunctionCall` object for it, otherwise `null`.
//
// Java: PeepholeReplaceKnownMethods.createConcatFunctionCallForNode(Node)
func PeepholeReplaceKnownMethods_CreateConcatFunctionCallForNode(n *rhino.Node) *peepholeReplaceKnownMethods_ConcatFunctionCall {
	jutil.CheckArgumentMsg1(n.IsCall(), "%s", n)
	callTarget := jutil.CheckNotNull(n.GetFirstChild())
	if !callTarget.IsGetProp() || callTarget.GetString() != "concat" {
		return nil
	}
	calleeNode := callTarget.GetFirstChild()
	if !PeepholeReplaceKnownMethods_ContainsExactlyArray(calleeNode) {
		return nil
	}
	firstArgumentNode := n.GetSecondChild()
	return newPeepholeReplaceKnownMethods_ConcatFunctionCall(n, calleeNode, firstArgumentNode)
}

// PeepholeReplaceKnownMethods_ContainsExactlyArray checks if a node is a
// known array. Checks for array literals and nested .concat calls.
//
// Java: PeepholeReplaceKnownMethods.containsExactlyArray(Node)
func PeepholeReplaceKnownMethods_ContainsExactlyArray(n *rhino.Node) bool {
	if n == nil {
		return false
	}

	if n.IsArrayLit() {
		return true
	}

	// Check for "[].concat(1)"
	if !n.IsCall() {
		return false
	}
	callee := n.GetFirstChild()
	return callee.IsGetProp() &&
		callee.GetString() == "concat" &&
		PeepholeReplaceKnownMethods_ContainsExactlyArray(callee.GetFirstChild())
}
