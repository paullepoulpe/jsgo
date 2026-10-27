// Port of com/google/javascript/jscomp/PeepholeReplaceKnownMethods.java @ closure-compiler v20261021

/*
 * Copyright 2011 The Closure Compiler Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package jscomp

import (
	"math"
	"math/big"
	"math/bits"
	"strings"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// Go-only: this pass uses the typed AST view (internal/ast) to read the
// calls it folds: the callee (call.Callee(), a GetProp or a Name), the
// receiver (getProp.Object(): a string or array literal, or the Array, Math
// and Number names), the method name (getProp.Property()) and the arguments
// (call.FirstArg(), call.OnlyArg(), arg.Next(), call.Args()) instead of
// child positions. Each tryFoldX method takes the kinds it folds (ast.Call,
// ast.GetProp, ast.StringLit, ast.Arg), so the Java checks of node tokens
// and of missing arguments are gone, and an argument that may be missing is
// read by the method that handles its absence (tryFoldStringAt,
// tryFoldStringSplit) instead of being passed as null. The string and
// number value logic, the elements folded by Array.join and the edits that
// move children in bulk or in place (Array.of, the join's new elements, the
// charAt rewrite) use the rhino API, in the Java's order.

// Java: PeepholeReplaceKnownMethods.TRIM_BOTH_PATTERN
var peepholeReplaceKnownMethods_TRIM_BOTH_PATTERN = peepholeReplaceKnownMethods_trimPattern{leading: true, trailing: true}

// Java: PeepholeReplaceKnownMethods.TRIM_START_PATTERN
var peepholeReplaceKnownMethods_TRIM_START_PATTERN = peepholeReplaceKnownMethods_trimPattern{leading: true}

// Java: PeepholeReplaceKnownMethods.TRIM_END_PATTERN
var peepholeReplaceKnownMethods_TRIM_END_PATTERN = peepholeReplaceKnownMethods_trimPattern{trailing: true}

// peepholeReplaceKnownMethods_trimPattern is one of the compiled patterns
// TRIM_BOTH_PATTERN (`^W|W\z`), TRIM_START_PATTERN (`^W`) and
// TRIM_END_PATTERN (`W\z`), where W is ECMA_WHITESPACE.
//
// Go-only: language-adapter — `\z` (jutil.Pattern rejects it), so this is an
// equivalent hand matcher for these java.util.regex uses (PORTING.md,
// Regex). With W a greedy run of whitespace characters, `^W` matches the
// maximal leading run, and `W\z` matches only a run that reaches the end of
// the input: the first position from which one does is the start of the
// maximal trailing run. Every character of W is in the BMP and is not a
// surrogate, so matching code points of the WTF-8 string is matching Java's.
type peepholeReplaceKnownMethods_trimPattern struct {
	leading, trailing bool
}

// replaceAll is Pattern.matcher(s).replaceAll("").
func (t peepholeReplaceKnownMethods_trimPattern) replaceAll(s string) string {
	start := 0
	if t.leading {
		for start < len(s) {
			r, size := utf8.DecodeRuneInString(s[start:])
			if !peepholeReplaceKnownMethods_isEcmaWhitespace(r) {
				break
			}
			start += size
		}
	}
	end := len(s)
	if t.trailing {
		for end > start {
			r, size := utf8.DecodeLastRuneInString(s[start:end])
			if !peepholeReplaceKnownMethods_isEcmaWhitespace(r) {
				break
			}
			end -= size
		}
	}
	return s[start:end]
}

// peepholeReplaceKnownMethods_isEcmaWhitespace is the character class of
// ECMA_WHITESPACE: [ \t\n-\r\u00A0\u1680\u2000-\u200A\u2028\u2029\u202F
// \u205F\u3000\uFEFF].
func peepholeReplaceKnownMethods_isEcmaWhitespace(r rune) bool {
	switch {
	case r == ' ' || r == '\t' || (r >= '\n' && r <= '\r'):
		return true
	case r == 0x00A0 || r == 0x1680 || (r >= 0x2000 && r <= 0x200A):
		return true
	case r == 0x2028 || r == 0x2029 || r == 0x202F || r == 0x205F || r == 0x3000 || r == 0xFEFF:
		return true
	}
	return false
}

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
func (p *PeepholeReplaceKnownMethods) OptimizeSubtree(subtree rhino.Node) rhino.Node {
	call, ok := ast.AsCall(subtree)
	if !ok {
		return subtree
	}
	if p.late && NodeUtil_IsGoogWeakUsageCall(subtree) {
		// goog.weakUsage(x) -> x.
		// (ProcessClosurePrimitives reports a goog.weakUsage call that does
		// not have exactly one argument.)
		x, ok := call.FirstArg()
		jutil.CheckState(ok)
		name := x.Detach().Node()
		subtree.ReplaceWith(name)
		p.ReportChangeToEnclosingScope(name)
		return name
	}
	return p.TryFoldKnownMethods(call)
}

// Java: PeepholeReplaceKnownMethods.tryFoldKnownMethods(Node)
func (p *PeepholeReplaceKnownMethods) TryFoldKnownMethods(call ast.Call) rhino.Node {
	// For now we only support string methods .join(),
	// .indexOf(), .substring() and .substr()
	// array method concat()
	// and numeric methods parseInt() and parseFloat().

	folded := p.TryFoldArrayJoin(call)
	// tryFoldArrayJoin may return a string literal instead of a CALL node
	call, ok := ast.AsCall(folded)
	if !ok {
		return folded
	}
	call = p.TryToFoldArrayConcat(call)

	callTarget := call.Callee()
	if method, ok := callTarget.AsGetProp(); ok {
		if receiver, ok := method.Object().AsName(); ok {
			switch receiver.Ident() {
			case "Array":
				if p.IsASTNormalized() {
					return p.TryFoldKnownArrayMethods(call, method)
				}
			case "Math":
				if p.IsASTNormalized() {
					return p.TryFoldKnownMathMethods(call, method)
				}
			case "Number":
				if p.IsASTNormalized() {
					return p.TryFoldKnownNumberMethods(call, method)
				}
			default:
			}
		}
		return p.TryFoldKnownStringMethods(call, method)
	} else if function, ok := callTarget.AsName(); ok {
		return p.TryFoldKnownNumericMethods(call, function)
	}
	return call.Node()
}

// TryFoldKnownArrayMethods tries to evaluate a method on the Array object.
//
// Java: PeepholeReplaceKnownMethods.tryFoldKnownArrayMethods(Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldKnownArrayMethods(call ast.Call, callTarget ast.GetProp) rhino.Node {
	subtree := call.Node()
	switch callTarget.Property() {
	case "of":
		// Array.of(args) -> [args]: the arguments move to the literal at once.
		callTarget.Detach()

		arraylit := rhino.NewNode(rhino.Token_ARRAYLIT)
		arraylit.AddChildrenToBack(subtree.RemoveChildren())
		subtree.ReplaceWith(arraylit)
		p.ReportChangeToEnclosingScope(arraylit)
		return arraylit
	case "isArray":
		return p.TryFoldArrayIsArray(call)
	default:
		return subtree
	}
}

// TryFoldArrayIsArray tries to evaluate Array.isArray:
//
//	Array.isArray([]) -> true
//	Array.isArray([1, 2, 3]) -> true
//	Array.isArray({}) -> false
//	Array.isArray(123) -> false
//	Array.isArray('hello') -> false
//	Array.isArray(null) -> false
//	Array.isArray() -> false
//
// Java: PeepholeReplaceKnownMethods.tryFoldArrayIsArray(Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldArrayIsArray(call ast.Call) rhino.Node {
	subtree := call.Node()
	if !call.HasArgs() {
		// Array.isArray() with 0 arguments evaluates to Array.isArray(undefined) -> false.
		newNode := NodeUtil_BooleanNode(false).Srcref(subtree)
		subtree.ReplaceWith(newNode)
		p.ReportChangeToEnclosingScope(newNode)
		return newNode
	}

	arg, ok := call.OnlyArg()
	if !ok {
		// 2 or more arguments: preserve call to avoid dropping potential side effects or arguments.
		return subtree
	}

	firstArg := arg.Node()
	if p.MayHaveSideEffects(firstArg) {
		return subtree
	}

	var isArray bool
	if firstArg.IsArrayLit() {
		isArray = true
	} else if NodeUtil_IsImmutableValue(firstArg) ||
		firstArg.IsObjectLit() ||
		firstArg.IsRegExp() ||
		firstArg.IsFunction() ||
		firstArg.IsClass() {
		isArray = false
	} else {
		return subtree
	}

	p.MarkFunctionsDeleted(subtree)
	newNode := NodeUtil_BooleanNode(isArray).Srcref(subtree)
	subtree.ReplaceWith(newNode)
	p.ReportChangeToEnclosingScope(newNode)
	return newNode
}

// Java: PeepholeReplaceKnownMethods.tryFoldKnownNumberMethods(Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldKnownNumberMethods(call ast.Call, callTarget ast.GetProp) rhino.Node {
	subtree := call.Node()
	methodName := callTarget.Property()
	switch methodName {
	case "parseInt", "parseFloat":
		if firstArg, ok := call.FirstArg(); ok &&
			(firstArg.Token() == rhino.Token_STRINGLIT || NodeUtil_IsNumericLiteral(firstArg.Node())) {
			return p.TryFoldParseNumber(call, methodName, firstArg)
		}
		return subtree
	default:
	}

	// Get the only number arg as a double as long as it is the only argument
	arg, ok := call.OnlyArg()
	if !ok {
		return subtree
	}
	valueNode := arg.Node()
	onlyArg, ok := p.GetSideEffectFreeNumberValueNoConversion(valueNode)
	var replacement rhino.Node
	if ok {
		switch methodName {
		case "isFinite":
			replacement = NodeUtil_BooleanNode(!math.IsInf(onlyArg, 0) && !math.IsNaN(onlyArg))
		case "isNaN":
			replacement = NodeUtil_BooleanNode(math.IsNaN(onlyArg))
		case "isSafeInteger":
			replacement =
				NodeUtil_BooleanNode(
					onlyArg < jutil.MathPow(2, 53) &&
						onlyArg >= (-1*(jutil.MathPow(2, 53)-1)) &&
						float64(jutil.D2L(onlyArg)) == onlyArg)
		default:
		}
	} else if !p.MayHaveSideEffects(valueNode) &&
		(NodeUtil_IsImmutableValue(valueNode) || NodeUtil_IsLiteralValue(valueNode, false)) {
		switch methodName {
		case "isFinite", "isNaN", "isSafeInteger":
			replacement = NodeUtil_BooleanNode(false)
		default:
		}
	}
	if !replacement.IsNil() {
		subtree.ReplaceWith(replacement)
		p.ReportChangeToEnclosingScope(replacement)
		return replacement
	}
	return subtree
}

// TryFoldKnownMathMethods tries to evaluate a method on the Math object.
//
// Java: PeepholeReplaceKnownMethods.tryFoldKnownMathMethods(Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldKnownMathMethods(call ast.Call, callTarget ast.GetProp) rhino.Node {
	subtree := call.Node()

	// first collect the arguments, if they are all numbers then we proceed
	var args []float64
	for arg := range call.Args() {
		d, ok := p.GetSideEffectFreeNumberValue(arg.Node())
		if ok {
			// lazily allocate, most calls will not be optimizable
			args = append(args, d)
		} else {
			return subtree
		}
	}
	var replacement float64
	hasReplacement := false
	methodName := callTarget.Property()
	// NOTE: the standard does not define precision for these methods, but we are conservative, so
	// for now we only implement the methods that are guaranteed to not increase the size of the
	// numeric constants.
	if len(args) == 1 {
		arg := args[0]
		switch methodName {
		case "abs":
			replacement, hasReplacement = math.Abs(arg), true
		case "ceil":
			replacement, hasReplacement = math.Ceil(arg), true
		case "floor":
			replacement, hasReplacement = math.Floor(arg), true
		case "fround":
			if math.IsNaN(arg) || math.IsInf(arg, 0) || arg == 0 {
				replacement, hasReplacement = arg, true
				// if the double is exactly representable as a float, then just cast since no rounding
				// is involved
			} else if float64(float32(arg)) == arg {
				// TODO(b/155511629): This condition is always true after J2CL transpilation.
				replacement, hasReplacement = float64(float32(arg)), true
			}
			// else: (float) arg does not necessarily use the correct rounding mode, so don't do anything
		case "round":
			if math.IsNaN(arg) || math.IsInf(arg, 0) {
				replacement, hasReplacement = arg, true
			} else if arg >= -0.5 &&
				(arg < 0.0 || (arg == 0.0 && int64(math.Float64bits(arg)) < 0)) {
				// Go-only: language-adapter — the Go constant -0.0 is +0.
				replacement, hasReplacement = math.Copysign(0, -1), true
			} else {
				replacement, hasReplacement = float64(jutil.MathRound(arg)), true
			}
		case "sign":
			replacement, hasReplacement = jutil.MathSignum(arg), true
		case "trunc":
			if math.IsNaN(arg) || math.IsInf(arg, 0) {
				replacement, hasReplacement = arg, true
			} else {
				replacement, hasReplacement = jutil.MathSignum(arg)*math.Floor(math.Abs(arg)), true
			}
		case "clz32":
			replacement, hasReplacement = float64(bits.LeadingZeros32(uint32(base.JSCompDoubles_EcmascriptToUint32(arg)))), true
		default:
		}
	}
	// handle the variadic functions now if we haven't already
	// For each of these we could allow for some of the values to be unknown and either reduce to
	// NaN or simplify the existing args. e.g. Math.max(3, x, 2) -> Math.max(3, x)
	if !hasReplacement {
		switch methodName {
		case "max":
			result := math.Inf(-1)
			for _, d := range args {
				result = jutil.MathMax(result, d)
			}
			replacement, hasReplacement = result, true
		case "min":
			result := math.Inf(1)
			for _, d := range args {
				result = jutil.MathMin(result, d)
			}
			replacement, hasReplacement = result, true
		case "imul":
			if len(args) < 2 {
				replacement, hasReplacement = 0.0, true
			} else {
				// Ignore args3+
				replacement = float64(base.JSCompDoubles_EcmascriptToInt32(args[0]) * base.JSCompDoubles_EcmascriptToInt32(args[1]))
				hasReplacement = true
			}
		case "pow":
			if len(args) == 2 {
				replacement, hasReplacement = jutil.MathPow(args[0], args[1]), true
			}
		default:
		}
	}

	if hasReplacement {
		numberNode := NodeUtil_NumberNode(replacement, subtree)
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
func (p *PeepholeReplaceKnownMethods) TryFoldKnownStringMethods(call ast.Call, callTarget ast.GetProp) rhino.Node {
	subtree := call.Node()

	// check if this is a call on a string method
	// then dispatch to specific folding method.
	stringNode := callTarget.Object()

	str, isStringLiteral := stringNode.AsStringLit()
	if !isStringLiteral && !p.useTypes {
		return subtree
	}
	functionNameString := callTarget.Property()
	firstArg, hasArgs := call.FirstArg()
	if isStringLiteral {
		if functionNameString == "split" {
			return p.TryFoldStringSplit(call, str)
		} else if !hasArgs {
			switch functionNameString {
			case "toLowerCase":
				return p.TryFoldStringToLowerCase(call, str)
			case "toUpperCase":
				return p.TryFoldStringToUpperCase(call, str)
			case "trim", "trimStart", "trimEnd", "trimLeft", "trimRight":
				return p.TryFoldStringTrim(call, str, functionNameString)
			case "at":
				return p.TryFoldStringAt(call, str)
			default:
			}
		} else {
			if NodeUtil_IsImmutableValue(firstArg.Node()) {
				switch functionNameString {
				case "indexOf", "lastIndexOf":
					return p.TryFoldStringIndexOf(call, functionNameString, str, firstArg)
				case "includes":
					return p.TryFoldStringIncludes(call, str, firstArg)
				case "startsWith":
					return p.TryFoldStringStartsWith(call, str, firstArg)
				case "endsWith":
					return p.TryFoldStringEndsWith(call, str, firstArg)
				case "at":
					return p.TryFoldStringAt(call, str)
				case "substr":
					return p.TryFoldStringSubstr(call, str, firstArg)
				case "substring", "slice":
					return p.TryFoldStringSubstringOrSlice(call, str, firstArg)
				case "charAt":
					return p.TryFoldStringCharAt(call, str, firstArg)
				case "charCodeAt":
					return p.TryFoldStringCharCodeAt(call, str, firstArg)
				case "replace":
					return p.TryFoldStringReplace(call, str, firstArg)
				case "replaceAll":
					return p.TryFoldStringReplaceAll(call, str, firstArg)
				default:
				}
			}
		}
	}
	if !p.useTypes ||
		!hasArgs ||
		!(isStringLiteral || colors.StandardColors_STRING.Equals(stringNode.Node().GetColor())) {
		return subtree
	}
	// A call with exactly two arguments: the start, and the length or the end.
	secondArg, ok := firstArg.Next()
	if !ok {
		return subtree
	}
	if _, ok := secondArg.Next(); ok {
		return subtree
	}
	maybeStart, ok := p.GetSideEffectFreeNumberValue(firstArg.Node())
	if !ok {
		return subtree
	}
	start := int(jutil.D2I(maybeStart))
	maybeLengthOrEnd, ok := p.GetSideEffectFreeNumberValue(secondArg.Node())
	if !ok {
		return subtree
	}
	switch functionNameString {
	case "substr":
		length := int(jutil.D2I(maybeLengthOrEnd))
		if start >= 0 && length == 1 {
			return p.ReplaceWithCharAt(call, callTarget, firstArg, secondArg)
		}
	case "substring", "slice":
		end := int(jutil.D2I(maybeLengthOrEnd))
		// unlike slice and substring, chatAt can not be used with negative indexes
		if start >= 0 && jutil.IntSub(end, start) == 1 {
			return p.ReplaceWithCharAt(call, callTarget, firstArg, secondArg)
		}
	default:
	}
	return subtree
}

// TryFoldKnownNumericMethods tries to evaluate known Numeric methods
// parseInt(), parseFloat()
//
// Java: PeepholeReplaceKnownMethods.tryFoldKnownNumericMethods(Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldKnownNumericMethods(call ast.Call, callTarget ast.Name) rhino.Node {
	functionNameString := callTarget.Ident()
	if functionNameString == "parseInt" || functionNameString == "parseFloat" {
		if firstArgument, ok := call.FirstArg(); ok &&
			(firstArgument.Token() == rhino.Token_STRINGLIT || NodeUtil_IsNumericLiteral(firstArgument.Node())) &&
			p.IsASTNormalized() {
			return p.TryFoldParseNumber(call, functionNameString, firstArgument)
		}
	}
	return call.Node()
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
func (p *PeepholeReplaceKnownMethods) TryFoldStringToLowerCase(call ast.Call, stringNode ast.StringLit) rhino.Node {
	lowered := jutil.StringToLowerCaseRoot(stringNode.Value())
	replacement := rhino.IR_String(lowered)
	call.Node().ReplaceWith(replacement)
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
func (p *PeepholeReplaceKnownMethods) TryFoldStringToUpperCase(call ast.Call, stringNode ast.StringLit) rhino.Node {
	upped := jutil.StringToUpperCaseRoot(stringNode.Value())
	replacement := rhino.IR_String(upped)
	call.Node().ReplaceWith(replacement)
	p.ReportChangeToEnclosingScope(replacement)
	return replacement
}

// TryFoldStringTrim returns The trimmed string Node.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringTrim(Node, Node, String)
func (p *PeepholeReplaceKnownMethods) TryFoldStringTrim(call ast.Call, stringNode ast.StringLit, functionName string) rhino.Node {
	subtree := call.Node()
	// See ECMA-262 § 22.1.3.32 (trim), § 22.1.3.34 (trimStart), § 22.1.3.33 (trimEnd),
	// Annex B.2.2.15 (trimLeft), Annex B.2.2.16 (trimRight), § 12.2, and § 12.3
	// All Unicode whitespace + LineTerminator + BOM
	var result string
	switch functionName {
	case "trim":
		result = peepholeReplaceKnownMethods_TRIM_BOTH_PATTERN.replaceAll(stringNode.Value())
	case "trimStart", "trimLeft":
		result = peepholeReplaceKnownMethods_TRIM_START_PATTERN.replaceAll(stringNode.Value())
	case "trimEnd", "trimRight":
		result = peepholeReplaceKnownMethods_TRIM_END_PATTERN.replaceAll(stringNode.Value())
	default:
		return subtree
	}
	replacement := rhino.IR_String(result).Srcref(subtree)
	subtree.ReplaceWith(replacement)
	p.ReportChangeToEnclosingScope(replacement)
	return replacement
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
func (p *PeepholeReplaceKnownMethods) TryFoldParseNumber(call ast.Call, functionName string, firstArg ast.Arg) rhino.Node {
	n := call.Node()
	isParseInt := functionName == "parseInt"

	// Second argument is only used as the radix for parseInt
	radix := 0
	if secondArg, ok := firstArg.Next(); ok {
		if !isParseInt {
			return n
		}

		// Third-argument and non-numeric second arg are problematic. Discard.
		radixNode, isNumber := secondArg.AsNumber()
		if _, ok := secondArg.Next(); ok || !isNumber {
			return n
		}
		tmpRadix := radixNode.Value()
		if tmpRadix != float64(jutil.D2I(tmpRadix)) {
			return n
		}
		radix = int(jutil.D2I(tmpRadix))
		if radix < 0 || radix == 1 || radix > 36 {
			return n
		}
	}

	if !isParseInt {
		// parseFloat logic
		if NodeUtil_IsNumericLiteral(firstArg.Node()) {
			checkVal, _ := p.GetSideEffectFreeNumberValue(firstArg.Node())
			numericNode := NodeUtil_NumberNode(checkVal, n)
			n.ReplaceWith(numericNode)
			p.ReportChangeToEnclosingScope(numericNode)
			return numericNode
		} else {
			stringVal, ok := p.GetSideEffectFreeStringValue(firstArg.Node())
			if !ok {
				return n
			}
			if _, ok := nodeUtil_getStringNumberValue(stringVal); !ok {
				return n
			}
			stringVal = NodeUtil_TrimJsWhiteSpace(stringVal)
			if stringVal == "" {
				return n
			}
			if stringVal == "0" {
				newNode := rhino.IR_Number(0)
				n.ReplaceWith(newNode)
				p.ReportChangeToEnclosingScope(newNode)
				return newNode
			}
			normalizedNewVal := "0"
			var result rhino.Node
			jutil.Try(func() {
				newVal := jutil.DoubleParseDouble(stringVal)
				newNode := NodeUtil_NumberNode(newVal, n)
				normalizedNewVal = PeepholeReplaceKnownMethods_NormalizeNumericString(jutil.DoubleToString(newVal))
				if PeepholeReplaceKnownMethods_NormalizeNumericString(stringVal) != normalizedNewVal {
					result = n
					return
				}
				n.ReplaceWith(newNode)
				p.ReportChangeToEnclosingScope(newNode)
				result = newNode
			}, func(e jutil.INumberFormatException) {
				result = n
			})
			return result
		}
	}

	// parseInt logic following ECMA-262 semantics
	stringVal, ok := p.GetSideEffectFreeStringValue(firstArg.Node())
	if !ok {
		return n
	}

	stringVal = NodeUtil_TrimJsWhiteSpace(stringVal)
	if stringVal == "" {
		return n
	}

	isNegative := false
	if jutil.StringStartsWith(stringVal, "-") {
		isNegative = true
		stringVal = jutil.StringSubstringFrom(stringVal, 1)
	} else if jutil.StringStartsWith(stringVal, "+") {
		stringVal = jutil.StringSubstringFrom(stringVal, 1)
	}

	if radix == 0 || radix == 16 {
		if jutil.StringLength(stringVal) > 1 && jutil.AsciiEqualsIgnoreCase(jutil.StringSubstring(stringVal, 0, 2), "0x") {
			radix = 16
			stringVal = jutil.StringSubstringFrom(stringVal, 2)
		} else if radix == 0 {
			if !p.IsEcmaScript5OrGreater() && jutil.StringStartsWith(stringVal, "0") {
				return n
			}
			radix = 10
		}
	}

	chars := jutil.StringToUTF16(stringVal)
	endDigitIndex := 0
	for endDigitIndex < len(chars) {
		c := chars[endDigitIndex]
		if PeepholeReplaceKnownMethods_GetRadixDigit(c, radix) < 0 {
			break
		}
		endDigitIndex++
	}

	if endDigitIndex == 0 {
		return n
	}

	digits := jutil.StringSubstring(stringVal, 0, endDigitIndex)
	var newVal float64
	failed := false
	jutil.Try(func() {
		bi := jutil.NewBigIntegerFromStringRadix(digits, radix)
		// BigInteger.doubleValue(): the nearest double, ties to even, or an
		// infinity when out of range.
		newVal, _ = new(big.Float).SetInt(bi).Float64()
	}, func(e jutil.INumberFormatException) {
		failed = true
	})
	if failed {
		return n
	}

	if isNegative {
		newVal = -newVal
	}

	newNode := NodeUtil_NumberNode(newVal, n)
	n.ReplaceWith(newNode)
	p.ReportChangeToEnclosingScope(newNode)
	return newNode
}

// Java: PeepholeReplaceKnownMethods.getRadixDigit(char, int)
func PeepholeReplaceKnownMethods_GetRadixDigit(c uint16, radix int) int {
	if c >= '0' && c <= '9' {
		d := int(c - '0')
		if d < radix {
			return d
		}
		return -1
	}
	if c >= 'a' && c <= 'z' {
		d := int(c-'a') + 10
		if d < radix {
			return d
		}
		return -1
	}
	if c >= 'A' && c <= 'Z' {
		d := int(c-'A') + 10
		if d < radix {
			return d
		}
		return -1
	}
	return -1
}

// TryFoldStringIndexOf tries to evaluate String.indexOf/lastIndexOf:
//
//	"abcdef".indexOf("bc") -> 1
//	"abcdefbc".indexOf("bc", 3) -> 6
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringIndexOf(Node, String, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringIndexOf(call ast.Call, functionName string, lstringNode ast.StringLit, firstArg ast.Arg) rhino.Node {
	n := call.Node()
	lstring := lstringNode.Value()
	isIndexOf := functionName == "indexOf"
	searchValue, ok := p.GetSideEffectFreeStringValue(firstArg.Node())
	// searchValue must be a valid string.
	if !ok {
		return n
	}
	fromIndex := 0
	if !isIndexOf {
		fromIndex = jutil.StringLength(lstring)
	}
	if secondArg, ok := firstArg.Next(); ok {
		// Third-argument and non-numeric second arg are problematic. Discard.
		from, isNumber := secondArg.AsNumber()
		if _, ok := secondArg.Next(); ok || !isNumber {
			return n
		}
		fromIndex = int(jutil.D2I(from.Value()))
	}
	var indexVal int
	if isIndexOf {
		indexVal = jutil.StringIndexOfFrom(lstring, searchValue, fromIndex)
	} else {
		indexVal = jutil.StringLastIndexOfFrom(lstring, searchValue, fromIndex)
	}
	newNode := NodeUtil_NumberNode(float64(indexVal), n)
	n.ReplaceWith(newNode)
	p.ReportChangeToEnclosingScope(newNode)

	return newNode
}

// TryFoldStringIncludes tries to evaluate String.prototype.includes:
//
//	"abcdef".includes("bc") -> true
//	"abcdef".includes("bc", 2) -> false
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringIncludes(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringIncludes(call ast.Call, lstringNode ast.StringLit, firstArg ast.Arg) rhino.Node {
	n := call.Node()
	if _, ok := firstArg.AsRegExp(); ok {
		return n
	}

	lstring := lstringNode.Value()
	searchValue, ok := p.GetSideEffectFreeStringValue(firstArg.Node())
	if !ok {
		return n
	}

	fromIndex := 0
	if arg, ok := firstArg.Next(); ok {
		secondArg := arg.Node()
		if _, ok := arg.Next(); ok || p.MayHaveSideEffects(secondArg) {
			return n
		}
		if NodeUtil_IsNumericLiteral(secondArg) {
			posVal := NodeUtil_GetNumericLiteralValue(secondArg)
			fromIndex = int(jutil.D2I(jutil.MathMin(jutil.MathMax(posVal, 0), float64(jutil.StringLength(lstring)))))
		} else if !NodeUtil_IsNullOrUndefined(secondArg) {
			return n
		}
	}

	contains := jutil.StringIndexOfFrom(lstring, searchValue, fromIndex) != -1
	newNode := NodeUtil_BooleanNode(contains).Srcref(n)
	n.ReplaceWith(newNode)
	p.ReportChangeToEnclosingScope(newNode)
	return newNode
}

// TryFoldStringStartsWith tries to evaluate String.prototype.startsWith:
//
//	"abcdef".startsWith("abc") -> true
//	"abcdef".startsWith("bc", 1) -> true
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringStartsWith(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringStartsWith(call ast.Call, lstringNode ast.StringLit, firstArg ast.Arg) rhino.Node {
	n := call.Node()
	if _, ok := firstArg.AsRegExp(); ok {
		return n
	}

	lstring := lstringNode.Value()
	searchValue, ok := p.GetSideEffectFreeStringValue(firstArg.Node())
	if !ok {
		return n
	}

	start := 0
	if arg, ok := firstArg.Next(); ok {
		secondArg := arg.Node()
		if _, ok := arg.Next(); ok || p.MayHaveSideEffects(secondArg) {
			return n
		}
		if NodeUtil_IsNumericLiteral(secondArg) {
			posVal := NodeUtil_GetNumericLiteralValue(secondArg)
			start = int(jutil.D2I(jutil.MathMin(jutil.MathMax(posVal, 0), float64(jutil.StringLength(lstring)))))
		} else if !NodeUtil_IsNullOrUndefined(secondArg) {
			return n
		}
	}

	// String.startsWith(String, int)
	result := jutil.StringRegionMatches(lstring, start, searchValue, 0, jutil.StringLength(searchValue))
	newNode := NodeUtil_BooleanNode(result).Srcref(n)
	n.ReplaceWith(newNode)
	p.ReportChangeToEnclosingScope(newNode)
	return newNode
}

// TryFoldStringEndsWith tries to evaluate String.prototype.endsWith:
//
//	"abcdef".endsWith("def") -> true
//	"abcdef".endsWith("de", 5) -> true
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringEndsWith(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringEndsWith(call ast.Call, lstringNode ast.StringLit, firstArg ast.Arg) rhino.Node {
	n := call.Node()
	if _, ok := firstArg.AsRegExp(); ok {
		return n
	}

	lstring := lstringNode.Value()
	searchValue, ok := p.GetSideEffectFreeStringValue(firstArg.Node())
	if !ok {
		return n
	}

	end := jutil.StringLength(lstring)
	if arg, ok := firstArg.Next(); ok {
		secondArg := arg.Node()
		if _, ok := arg.Next(); ok || p.MayHaveSideEffects(secondArg) {
			return n
		}
		if NodeUtil_IsNumericLiteral(secondArg) {
			posVal := NodeUtil_GetNumericLiteralValue(secondArg)
			end = int(jutil.D2I(jutil.MathMin(jutil.MathMax(posVal, 0), float64(jutil.StringLength(lstring)))))
		} else if secondArg.IsNull() {
			end = 0
		} else if NodeUtil_IsUndefined(secondArg) {
			end = jutil.StringLength(lstring)
		} else {
			return n
		}
	}

	searchLength := jutil.StringLength(searchValue)
	start := end - searchLength
	// String.startsWith(String, int)
	result := start >= 0 && jutil.StringRegionMatches(lstring, start, searchValue, 0, searchLength)
	newNode := NodeUtil_BooleanNode(result).Srcref(n)
	n.ReplaceWith(newNode)
	p.ReportChangeToEnclosingScope(newNode)
	return newNode
}

// TryFoldStringAt tries to evaluate String.prototype.at:
//
//	"abcdef".at(0) -> "a"
//	"abcdef".at(-1) -> "f"
//	"abcdef".at(10) -> void 0
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringAt(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringAt(call ast.Call, stringNode ast.StringLit) rhino.Node {
	n := call.Node()
	var posVal float64
	firstArg, ok := call.FirstArg()
	if !ok {
		posVal = 0.0
	} else {
		if _, ok := firstArg.Next(); ok {
			return n
		}
		maybePos, ok := p.GetSideEffectFreeNumberValue(firstArg.Node())
		if !ok {
			return n
		}
		if math.IsNaN(maybePos) {
			posVal = 0.0
		} else {
			posVal = maybePos
		}
	}

	stringVal := stringNode.Value()
	len_ := jutil.StringLength(stringVal)

	var replacement rhino.Node
	if math.IsInf(posVal, 0) || posVal >= float64(len_) || posVal <= float64(-(len_+1)) {
		replacement = NodeUtil_NewUndefinedNode(n)
	} else {
		intPos := jutil.D2L(posVal)
		var k int64
		if intPos >= 0 {
			k = intPos
		} else {
			k = int64(len_) + intPos
		}
		if k < 0 || k >= int64(len_) {
			replacement = NodeUtil_NewUndefinedNode(n)
		} else {
			replacement = rhino.IR_String(jutil.StringFromChar(jutil.StringCharAt(stringVal, int(k)))).Srcref(n)
		}
	}

	n.ReplaceWith(replacement)
	p.ReportChangeToEnclosingScope(replacement)
	return replacement
}

// TryFoldArrayJoin tries to fold an array join: ['a', 'b', 'c'].join(”) ->
// 'abc';
//
// Java: PeepholeReplaceKnownMethods.tryFoldArrayJoin(Node)
func (p *PeepholeReplaceKnownMethods) TryFoldArrayJoin(call ast.Call) rhino.Node {
	n := call.Node()
	callTarget, ok := call.Callee().AsGetProp()
	if !ok {
		return n
	}

	array, ok := callTarget.Object().AsArrayLit()
	if !ok || callTarget.Property() != "join" {
		return n
	}
	arrayNode := array.Node()

	// The separator: at most one argument, immutable.
	right, hasRight := call.FirstArg()
	if hasRight {
		if _, ok := right.Next(); ok || !NodeUtil_IsImmutableValue(right.Node()) {
			return n
		}
	}

	if hasRight {
		str, isStringLit := right.AsStringLit()
		if NodeUtil_IsUndefined(right.Node()) || (isStringLit && "," == str.Value()) {
			// "," is the default, it doesn't need to be explicit
			right.Detach()
			p.ReportChangeToEnclosingScope(n)
		}
	}

	// logic above ensures that `right` is immutable, so no need to check for
	// side effects with getSideEffectFreeStringValue(right)
	var joinString string
	if !hasRight || NodeUtil_IsUndefined(right.Node()) {
		joinString = ","
	} else if s, ok := nodeUtil_getStringValue(right.Node()); ok {
		joinString = s
	} else {
		joinString = "null"
	}
	var arrayFoldedChildren []rhino.Node
	var sb *jutil.StringBuilder
	foldedSize := 0
	var prev rhino.Node
	// Merges adjacent String nodes.
	for e := range array.Elements() {
		elem := e.Node()
		if NodeUtil_IsImmutableValue(elem) || elem.IsEmpty() {
			if sb == nil {
				sb = jutil.NewStringBuilder()
			} else {
				sb.Append(joinString)
			}
			elementStr, ok := nodeUtil_getArrayElementStringValue(elem)
			if !ok {
				return n // TODO(nickreid): Is this ever null?
			}
			sb.Append(elementStr)
		} else {
			if sb != nil {
				rhino.CheckNotNull(prev)
				// + 2 for the quotes.
				foldedSize += sb.Length() + 2
				arrayFoldedChildren = append(arrayFoldedChildren, rhino.IR_String(sb.String()).SrcrefIfMissing(prev))
				sb = nil
			}
			foldedSize += InlineCostEstimator_GetCost(elem)
			arrayFoldedChildren = append(arrayFoldedChildren, elem)
		}
		prev = elem
	}

	if sb != nil {
		rhino.CheckNotNull(prev)
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
		// (Like the Java, this counts the separator even when it was removed
		// above as the default one.)
		if hasRight {
			foldedSize += InlineCostEstimator_GetCost(right.Node())
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
func (p *PeepholeReplaceKnownMethods) TryFoldStringSubstr(call ast.Call, stringNode ast.StringLit, arg1 ast.Arg) rhino.Node {
	n := call.Node()

	var start int
	var length int
	stringAsString := stringNode.Value()

	maybeStart, ok := p.GetSideEffectFreeNumberValue(arg1.Node())
	if ok {
		start = int(jutil.D2I(maybeStart))
	} else {
		return n
	}

	if arg2, ok := arg1.Next(); ok {
		maybeLength, ok := p.GetSideEffectFreeNumberValue(arg2.Node())
		if ok {
			length = int(jutil.D2I(maybeLength))
		} else {
			return n
		}

		if _, ok := arg2.Next(); ok {
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
func (p *PeepholeReplaceKnownMethods) TryFoldStringSubstringOrSlice(call ast.Call, stringNode ast.StringLit, arg1 ast.Arg) rhino.Node {
	n := call.Node()

	var start int
	var end int
	stringAsString := stringNode.Value()

	maybeStart, ok := p.GetSideEffectFreeNumberValue(arg1.Node())
	if ok {
		start = int(jutil.D2I(maybeStart))
	} else {
		return n
	}

	if arg2, ok := arg1.Next(); ok {
		maybeEnd, ok := p.GetSideEffectFreeNumberValue(arg2.Node())
		if ok {
			end = int(jutil.D2I(maybeEnd))
		} else {
			return n
		}

		if _, ok := arg2.Next(); ok {
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

// ReplaceWithCharAt rewrites s.substr(start, 1), s.substring(start, start
// + 1) or s.slice(start, start + 1) in place to s.charAt(start): it renames
// the method and removes secondArg.
//
// Go-only: it takes the second argument (Java: firstArg.getNext()).
//
// Java: PeepholeReplaceKnownMethods.replaceWithCharAt(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) ReplaceWithCharAt(call ast.Call, callTarget ast.GetProp, firstArg ast.Arg, secondArg ast.Arg) rhino.Node {
	// TODO(moz): Maybe correct the arity of the function type here.
	callTarget.SetProperty("charAt")
	secondArg.Detach()
	p.ReportChangeToEnclosingScope(firstArg.Node())
	return call.Node()
}

// TryFoldStringCharAt tries to fold .charAt() calls on strings.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringCharAt(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringCharAt(call ast.Call, stringNode ast.StringLit, arg1 ast.Arg) rhino.Node {
	n := call.Node()

	var index int
	stringAsString := stringNode.Value()

	// The only argument, a number.
	num, isNumber := arg1.AsNumber()
	if _, more := arg1.Next(); isNumber && !more {
		index = int(jutil.D2I(num.Value()))
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
func (p *PeepholeReplaceKnownMethods) TryFoldStringReplace(call ast.Call, stringNode ast.StringLit, arg1 ast.Arg) rhino.Node {
	n := call.Node()
	arg2, ok := arg1.Next()
	if !ok {
		// too few parameters
		return n
	}
	if _, ok := arg2.Next(); ok {
		// too many parameters
		return n
	}

	pattern, isString1 := arg1.AsStringLit()
	replacement, isString2 := arg2.AsStringLit()
	if !isString1 || !isString2 {
		// only string literals are supported for folding.
		return n
	}

	lookForPattern := pattern.Value()
	replacementPattern := replacement.Value()
	if strings.Contains(replacementPattern, "$") {
		// 'special' replacements aren't supported yet.
		return n
	}

	original := stringNode.Value()

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

	resultNode := rhino.IR_String(newString).Srcref(stringNode.Node())
	parent := n.GetParent()
	n.ReplaceWith(resultNode)
	p.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// TryFoldStringReplaceAll tries to fold .replaceAll() calls on strings.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringReplaceAll(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringReplaceAll(call ast.Call, stringNode ast.StringLit, arg1 ast.Arg) rhino.Node {
	n := call.Node()
	arg2, ok := arg1.Next()
	if !ok {
		// too few parameters
		return n
	}
	if _, ok := arg2.Next(); ok {
		// too many parameters
		return n
	}

	pattern, isString1 := arg1.AsStringLit()
	replacement, isString2 := arg2.AsStringLit()
	if !isString1 || !isString2 {
		// only string literals are supported for folding.
		return n
	}

	replacementPattern := replacement.Value()
	if strings.Contains(replacementPattern, "$") {
		// 'special' replacements aren't supported yet.
		return n
	}

	// Java "replace" acts like JavaScript's "replaceAll" and replaces all occurrences.
	original := stringNode.Value()
	newString := jutil.StringReplace(original, pattern.Value(), replacementPattern)

	resultNode := rhino.IR_String(newString).Srcref(stringNode.Node())
	parent := n.GetParent()
	n.ReplaceWith(resultNode)
	p.ReportChangeToEnclosingScope(parent)
	return resultNode
}

// TryFoldStringCharCodeAt tries to fold .charCodeAt() calls on strings.
//
// Java: PeepholeReplaceKnownMethods.tryFoldStringCharCodeAt(Node, Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldStringCharCodeAt(call ast.Call, stringNode ast.StringLit, arg1 ast.Arg) rhino.Node {
	n := call.Node()

	var index int
	stringAsString := stringNode.Value()

	// The only argument, a number.
	num, isNumber := arg1.AsNumber()
	if _, more := arg1.Next(); isNumber && !more {
		index = int(jutil.D2I(num.Value()))
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
func (p *PeepholeReplaceKnownMethods) JsSplit(stringValue string, separator string, hasSeparator bool, limit int) []string {
	jutil.CheckArgument(limit >= 0)

	// For limits of 0, return an empty array
	if limit == 0 {
		return []string{}
	}

	// If a separator is not specified, return the entire string as
	// the only element of an array.
	if !hasSeparator {
		return []string{stringValue}
	}

	var splitStrings []string

	// If an empty string is specified for the separator, split apart each
	// character of the string.
	if separator == "" {
		for i := 0; i < jutil.StringLength(stringValue) && i < limit; i++ {
			splitStrings = append(splitStrings, jutil.StringSubstring(stringValue, i, i+1))
		}
	} else {
		startIndex := 0
		var matchIndex int
		for {
			matchIndex = PeepholeReplaceKnownMethods_JsSplitMatch(stringValue, startIndex, separator)
			if !(matchIndex >= 0 && len(splitStrings) < limit) {
				break
			}
			splitStrings = append(splitStrings, jutil.StringSubstring(stringValue, startIndex, matchIndex))

			startIndex = matchIndex + jutil.StringLength(separator)
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
func (p *PeepholeReplaceKnownMethods) TryFoldStringSplit(call ast.Call, stringNode ast.StringLit) rhino.Node {
	n := call.Node()
	if p.late {
		return n
	}

	var separator string
	hasSeparator := false
	stringValue := stringNode.Value()

	// Maximum number of possible splits
	limit := jutil.StringLength(stringValue) + 1

	if arg1, ok := call.FirstArg(); ok {
		if sep, ok := arg1.AsStringLit(); ok {
			separator, hasSeparator = sep.Value(), true
		} else if arg1.Token() != rhino.Token_NULL {
			return n
		}

		// Go-only note: like the Java, this does not look past the limit, so
		// "a,b".split(",", 1, f()) folds to ["a"] and drops the call to f.
		// Kept for identical output.
		if arg2, ok := arg1.Next(); ok {
			if NodeUtil_IsNumericLiteral(arg2.Node()) {
				// Integer.toUnsignedLong(int)
				uint32Limit :=
					int64(uint32(base.JSCompDoubles_EcmascriptToUint32(NodeUtil_GetNumericLiteralValue(arg2.Node()))))
				limit = int(min(uint32Limit, int64(limit)))
			} else {
				return n
			}
		}
	}

	// Split the string and convert the returned array into JS nodes
	stringArray := p.JsSplit(stringValue, separator, hasSeparator, limit)
	arrayOfStrings := rhino.IR_Arraylit()
	for _, element := range stringArray {
		arrayOfStrings.AddChildToBack(rhino.IR_String(element).Srcref(stringNode.Node()))
	}

	parent := n.GetParent()
	n.ReplaceWith(arrayOfStrings)
	p.ReportChangeToEnclosingScope(parent)
	return arrayOfStrings
}

// Java: PeepholeReplaceKnownMethods.tryToFoldArrayConcat(Node)
func (p *PeepholeReplaceKnownMethods) TryToFoldArrayConcat(n ast.Call) ast.Call {
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
	if array, ok := arrayLiteralToRemove.AsArrayLit(); !ok || array.HasElements() {
		return concatFunctionCall
	}
	firstArgument, ok := callNode.FirstArg()
	if !ok {
		return concatFunctionCall
	}
	firstArg, ok := firstArgument.AsExpr() // not a spread
	if !ok || !PeepholeReplaceKnownMethods_ContainsExactlyArray(firstArg) {
		return concatFunctionCall
	}

	firstArg.Detach()
	arrayLiteralToRemove.ReplaceWith(firstArg)

	p.ReportChangeToEnclosingScope(callNode.Node())
	return PeepholeReplaceKnownMethods_CreateConcatFunctionCallForNode(callNode)
}

// TryToFoldConcatChaining checks if we have this code pattern
// `array.concat(...*).concat(sideEffectFreeArguments)` and if yes fold
// chained concat functions, so `arr.concat(a).concat(b)` will be fold into
// `arr.concat(a,b)`.
//
// Java: PeepholeReplaceKnownMethods.tryToFoldConcatChaining(ConcatFunctionCall)
func (p *PeepholeReplaceKnownMethods) TryToFoldConcatChaining(concatFunctionCall *peepholeReplaceKnownMethods_ConcatFunctionCall) ast.Call {
	jutil.CheckNotNull(concatFunctionCall)

	concatCallNode := concatFunctionCall.callNode

	maybeFunctionCall, ok := concatFunctionCall.calleeNode.AsCall()
	if !ok {
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
	for arg := range concatCallNode.Args() {
		if p.MayHaveSideEffects(arg.Node()) {
			return concatCallNode
		}
	}

	// perform folding
	previousConcatCallNode := previousConcatFunctionCall.callNode
	arg, ok := concatCallNode.FirstArg()
	for ok {
		currentArg := arg
		arg, ok = arg.Next()
		previousConcatCallNode.AppendArg(currentArg.Detach())
	}
	concatCallNode.Node().ReplaceWith(previousConcatCallNode.Detach().Node())
	p.ReportChangeToEnclosingScope(previousConcatCallNode.Node())
	return previousConcatCallNode
}

// peepholeReplaceKnownMethods_ConcatFunctionCall is the ConcatFunctionCall
// class.
//
// Go-only: language-adapter — the Java class is abstract only so that it is
// instantiated as an empty anonymous subclass (`new ConcatFunctionCall(...)
// {}`), which adds nothing; one struct stands for both. The Java field
// firstArgumentNode (the first argument, or null) is read from callNode
// (callNode.FirstArg(), callNode.Args()) where it is used.
type peepholeReplaceKnownMethods_ConcatFunctionCall struct {
	callNode ast.Call
	// calleeNode is the array whose concat method is called.
	calleeNode ast.Expr
}

// Java: PeepholeReplaceKnownMethods.ConcatFunctionCall(Node, Node, Node)
func newPeepholeReplaceKnownMethods_ConcatFunctionCall(callNode ast.Call, calleeNode ast.Expr) *peepholeReplaceKnownMethods_ConcatFunctionCall {
	c := &peepholeReplaceKnownMethods_ConcatFunctionCall{}
	c.callNode = callNode
	c.calleeNode = calleeNode
	return c
}

// PeepholeReplaceKnownMethods_CreateConcatFunctionCallForNode returns, if the
// argument node is a call to `Array.prototype.concat`, a
// `ConcatFunctionCall` object for it, otherwise `null`.
//
// Java: PeepholeReplaceKnownMethods.createConcatFunctionCallForNode(Node)
func PeepholeReplaceKnownMethods_CreateConcatFunctionCallForNode(n ast.Call) *peepholeReplaceKnownMethods_ConcatFunctionCall {
	callTarget, ok := n.Callee().AsGetProp()
	if !ok || callTarget.Property() != "concat" {
		return nil
	}
	calleeNode := callTarget.Object()
	if !PeepholeReplaceKnownMethods_ContainsExactlyArray(calleeNode) {
		return nil
	}
	return newPeepholeReplaceKnownMethods_ConcatFunctionCall(n, calleeNode)
}

// PeepholeReplaceKnownMethods_ContainsExactlyArray checks if a node is a
// known array. Checks for array literals and nested .concat calls.
//
// Java: PeepholeReplaceKnownMethods.containsExactlyArray(Node)
func PeepholeReplaceKnownMethods_ContainsExactlyArray(n ast.Expr) bool {
	if _, ok := n.AsArrayLit(); ok {
		return true
	}

	// Check for "[].concat(1)"
	call, ok := n.AsCall()
	if !ok {
		return false
	}
	callee, ok := call.Callee().AsGetProp()
	return ok &&
		callee.Property() == "concat" &&
		PeepholeReplaceKnownMethods_ContainsExactlyArray(callee.Object())
}
