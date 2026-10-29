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

	"github.com/paullepoulpe/jsgo/internal/ast"
	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

var (
	rkmIDArray         = uint32(strpool.Intern("Array"))
	rkmIDMath          = uint32(strpool.Intern("Math"))
	rkmIDNumber        = uint32(strpool.Intern("Number"))
	rkmIDOf            = uint32(strpool.Intern("of"))
	rkmIDIsArray       = uint32(strpool.Intern("isArray"))
	rkmIDParseInt      = uint32(strpool.Intern("parseInt"))
	rkmIDParseFloat    = uint32(strpool.Intern("parseFloat"))
	rkmIDIsFinite      = uint32(strpool.Intern("isFinite"))
	rkmIDIsNaN         = uint32(strpool.Intern("isNaN"))
	rkmIDIsSafeInteger = uint32(strpool.Intern("isSafeInteger"))
	rkmIDAbs           = uint32(strpool.Intern("abs"))
	rkmIDCeil          = uint32(strpool.Intern("ceil"))
	rkmIDFloor         = uint32(strpool.Intern("floor"))
	rkmIDFround        = uint32(strpool.Intern("fround"))
	rkmIDRound         = uint32(strpool.Intern("round"))
	rkmIDSign          = uint32(strpool.Intern("sign"))
	rkmIDTrunc         = uint32(strpool.Intern("trunc"))
	rkmIDClz32         = uint32(strpool.Intern("clz32"))
	rkmIDMax           = uint32(strpool.Intern("max"))
	rkmIDMin           = uint32(strpool.Intern("min"))
	rkmIDImul          = uint32(strpool.Intern("imul"))
	rkmIDPow           = uint32(strpool.Intern("pow"))
	rkmIDJoin          = uint32(strpool.Intern("join"))
	rkmIDConcat        = uint32(strpool.Intern("concat"))
	rkmIDComma         = uint32(strpool.Intern(","))
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
	return &PeepholeReplaceKnownMethods{
		late:     late,
		useTypes: useTypes,
	}
}

func (p *PeepholeReplaceKnownMethods) peepholeConfig() peepholeOptConfig {
	return peepholeOptConfig{
		tokens:     peepholeCallOnlyTokens,
		isStandard: true,
	}
}

func (p *PeepholeReplaceKnownMethods) cloneForWorker() AbstractPeepholeOptimization {
	return NewPeepholeReplaceKnownMethods(p.late, p.useTypes)
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
			switch receiver.Node().GetStringID() {
			case rkmIDArray:
				if p.IsASTNormalized() {
					return p.TryFoldKnownArrayMethods(call, method)
				}
			case rkmIDMath:
				if p.IsASTNormalized() {
					return p.TryFoldKnownMathMethods(call, method)
				}
			case rkmIDNumber:
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
	switch callTarget.Node().GetStringID() {
	case rkmIDOf:
		// Array.of(args) -> [args]: the arguments move to the literal at once.
		callTarget.Detach()

		arraylit := rhino.NewNode(rhino.Token_ARRAYLIT)
		arraylit.AddChildrenToBack(subtree.RemoveChildren())
		subtree.ReplaceWith(arraylit)
		p.ReportChangeToEnclosingScope(arraylit)
		return arraylit
	case rkmIDIsArray:
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
	methodID := callTarget.Node().GetStringID()
	switch methodID {
	case rkmIDParseInt, rkmIDParseFloat:
		if firstArg, ok := call.FirstArg(); ok &&
			(firstArg.Token() == rhino.Token_STRINGLIT || NodeUtil_IsNumericLiteral(firstArg.Node())) {
			methodName := "parseFloat"
			if methodID == rkmIDParseInt {
				methodName = "parseInt"
			}
			return p.TryFoldParseNumber(call, methodName, firstArg)
		}
		return subtree
	case rkmIDIsFinite, rkmIDIsNaN, rkmIDIsSafeInteger:
	default:
		return subtree
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
		switch methodID {
		case rkmIDIsFinite:
			replacement = NodeUtil_BooleanNode(!math.IsInf(onlyArg, 0) && !math.IsNaN(onlyArg))
		case rkmIDIsNaN:
			replacement = NodeUtil_BooleanNode(math.IsNaN(onlyArg))
		case rkmIDIsSafeInteger:
			replacement =
				NodeUtil_BooleanNode(
					onlyArg < base.MathPow(2, 53) &&
						onlyArg >= (-1*(base.MathPow(2, 53)-1)) &&
						float64(jutil.D2L(onlyArg)) == onlyArg)
		}
	} else if !p.MayHaveSideEffects(valueNode) &&
		(NodeUtil_IsImmutableValue(valueNode) || NodeUtil_IsLiteralValue(valueNode, false)) {
		replacement = NodeUtil_BooleanNode(false)
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
	var argsBuf [4]float64
	args := argsBuf[:0]
	for arg := range call.Args() {
		d, ok := p.GetSideEffectFreeNumberValue(arg.Node())
		if ok {
			args = append(args, d)
		} else {
			return subtree
		}
	}
	var replacement float64
	hasReplacement := false
	methodID := callTarget.Node().GetStringID()
	// NOTE: the standard does not define precision for these methods, but we are conservative, so
	// for now we only implement the methods that are guaranteed to not increase the size of the
	// numeric constants.
	if len(args) == 1 {
		arg := args[0]
		switch methodID {
		case rkmIDAbs:
			replacement, hasReplacement = math.Abs(arg), true
		case rkmIDCeil:
			replacement, hasReplacement = math.Ceil(arg), true
		case rkmIDFloor:
			replacement, hasReplacement = math.Floor(arg), true
		case rkmIDFround:
			if math.IsNaN(arg) || math.IsInf(arg, 0) || arg == 0 {
				replacement, hasReplacement = arg, true
				// if the double is exactly representable as a float, then just cast since no rounding
				// is involved
			} else if float64(float32(arg)) == arg {
				// TODO(b/155511629): This condition is always true after J2CL transpilation.
				replacement, hasReplacement = float64(float32(arg)), true
			}
			// else: (float) arg does not necessarily use the correct rounding mode, so don't do anything
		case rkmIDRound:
			if math.IsNaN(arg) || math.IsInf(arg, 0) {
				replacement, hasReplacement = arg, true
			} else if arg >= -0.5 &&
				(arg < 0.0 || (arg == 0.0 && int64(math.Float64bits(arg)) < 0)) {
				// Go-only: language-adapter — the Go constant -0.0 is +0.
				replacement, hasReplacement = math.Copysign(0, -1), true
			} else {
				replacement, hasReplacement = float64(base.MathRound(arg)), true
			}
		case rkmIDSign:
			replacement, hasReplacement = base.MathSignum(arg), true
		case rkmIDTrunc:
			if math.IsNaN(arg) || math.IsInf(arg, 0) {
				replacement, hasReplacement = arg, true
			} else {
				replacement, hasReplacement = base.MathSignum(arg)*math.Floor(math.Abs(arg)), true
			}
		case rkmIDClz32:
			replacement, hasReplacement = float64(bits.LeadingZeros32(uint32(base.JSCompDoubles_EcmascriptToUint32(arg)))), true
		default:
		}
	}
	// handle the variadic functions now if we haven't already
	// For each of these we could allow for some of the values to be unknown and either reduce to
	// NaN or simplify the existing args. e.g. Math.max(3, x, 2) -> Math.max(3, x)
	if !hasReplacement {
		switch methodID {
		case rkmIDMax:
			result := math.Inf(-1)
			for _, d := range args {
				result = base.MathMax(result, d)
			}
			replacement, hasReplacement = result, true
		case rkmIDMin:
			result := math.Inf(1)
			for _, d := range args {
				result = base.MathMin(result, d)
			}
			replacement, hasReplacement = result, true
		case rkmIDImul:
			if len(args) < 2 {
				replacement, hasReplacement = 0.0, true
			} else {
				// Ignore args3+
				replacement = float64(base.JSCompDoubles_EcmascriptToInt32(args[0]) * base.JSCompDoubles_EcmascriptToInt32(args[1]))
				hasReplacement = true
			}
		case rkmIDPow:
			if len(args) == 2 {
				replacement, hasReplacement = base.MathPow(args[0], args[1]), true
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

// TryFoldKnownNumericMethods tries to evaluate known Numeric methods
// parseInt(), parseFloat()
//
// Java: PeepholeReplaceKnownMethods.tryFoldKnownNumericMethods(Node, Node)
func (p *PeepholeReplaceKnownMethods) TryFoldKnownNumericMethods(call ast.Call, callTarget ast.Name) rhino.Node {
	id := callTarget.Node().GetStringID()
	if id == rkmIDParseInt || id == rkmIDParseFloat {
		if firstArgument, ok := call.FirstArg(); ok &&
			(firstArgument.Token() == rhino.Token_STRINGLIT || NodeUtil_IsNumericLiteral(firstArgument.Node())) &&
			p.IsASTNormalized() {
			name := "parseFloat"
			if id == rkmIDParseInt {
				name = "parseInt"
			}
			return p.TryFoldParseNumber(call, name, firstArgument)
		}
	}
	return call.Node()
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
	endIndex := len(input) - 1
	for startIndex < len(input) && input[startIndex] == '0' {
		startIndex++
	}
	if strings.IndexByte(input, '.') >= 0 {
		for endIndex >= 0 && input[endIndex] == '0' {
			endIndex--
		}
		if endIndex >= 0 && input[endIndex] == '.' {
			endIndex--
		}
	}
	if startIndex >= endIndex {
		return input
	}
	return input[startIndex : endIndex+1]
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
		}
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
		newVal, ok := jutil.TryParseDouble(stringVal)
		if !ok {
			return n
		}
		newNode := NodeUtil_NumberNode(newVal, n)
		normalizedNewVal := PeepholeReplaceKnownMethods_NormalizeNumericString(jutil.DoubleToString(newVal))
		if PeepholeReplaceKnownMethods_NormalizeNumericString(stringVal) != normalizedNewVal {
			return n
		}
		n.ReplaceWith(newNode)
		p.ReportChangeToEnclosingScope(newNode)
		return newNode
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
	if stringVal[0] == '-' {
		isNegative = true
		stringVal = stringVal[1:]
	} else if stringVal[0] == '+' {
		stringVal = stringVal[1:]
	}

	if radix == 0 || radix == 16 {
		if len(stringVal) >= 2 && stringVal[0] == '0' && (stringVal[1] == 'x' || stringVal[1] == 'X') {
			radix = 16
			stringVal = stringVal[2:]
		} else if radix == 0 {
			if !p.IsEcmaScript5OrGreater() && len(stringVal) > 0 && stringVal[0] == '0' {
				return n
			}
			radix = 10
		}
	}

	endDigitIndex := 0
	for endDigitIndex < len(stringVal) {
		c := stringVal[endDigitIndex]
		if c >= 0x80 || PeepholeReplaceKnownMethods_GetRadixDigit(uint16(c), radix) < 0 {
			break
		}
		endDigitIndex++
	}

	if endDigitIndex == 0 {
		return n
	}

	digits := stringVal[:endDigitIndex]
	var newVal float64
	if val, ok := jutil.TryParseLong(digits, radix); ok && val >= -(1<<53) && val <= 1<<53 {
		newVal = float64(val)
	} else {
		bi, ok := jutil.TryParseBigInt(digits, radix)
		if !ok {
			return n
		}
		// BigInteger.doubleValue(): the nearest double, ties to even, or an
		// infinity when out of range.
		newVal, _ = new(big.Float).SetInt(bi).Float64()
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
	if !ok || callTarget.Node().GetStringID() != rkmIDJoin {
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
		if NodeUtil_IsUndefined(right.Node()) || (isStringLit && str.Node().GetStringID() == rkmIDComma) {
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
	var arrayFoldedBuf [8]rhino.Node
	arrayFoldedChildren := arrayFoldedBuf[:0]
	var sb strings.Builder
	hasSB := false
	foldedSize := 0
	var prev rhino.Node
	// Merges adjacent String nodes.
	for e := range array.Elements() {
		elem := e.Node()
		if NodeUtil_IsImmutableValue(elem) || elem.IsEmpty() {
			if !hasSB {
				hasSB = true
			} else {
				sb.WriteString(joinString)
			}
			elementStr, ok := nodeUtil_getArrayElementStringValue(elem)
			if !ok {
				return n // TODO(nickreid): Is this ever null?
			}
			sb.WriteString(elementStr)
		} else {
			if hasSB {
				rhino.CheckNotNull(prev)
				merged := sb.String()
				sb.Reset()
				hasSB = false
				// + 2 for the quotes.
				foldedSize += jutil.StringLength(merged) + 2
				arrayFoldedChildren = append(arrayFoldedChildren, rhino.IR_String(merged).SrcrefIfMissing(prev))
			}
			foldedSize += InlineCostEstimator_GetCost(elem)
			arrayFoldedChildren = append(arrayFoldedChildren, elem)
		}
		prev = elem
	}

	if hasSB {
		rhino.CheckNotNull(prev)
		merged := sb.String()
		// + 2 for the quotes.
		foldedSize += jutil.StringLength(merged) + 2
		arrayFoldedChildren = append(arrayFoldedChildren, rhino.IR_String(merged).SrcrefIfMissing(prev))
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
		}
		// Because of special case behavior for `null` and `undefined` values, there's no safe way
		// to convert `[someNonStringValue].join()` to something shorter.
		// e.g. String(someNonStringValue) would turn `null` into `"null"`, which isn't right.
		return n
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

// Java: PeepholeReplaceKnownMethods.tryToFoldArrayConcat(Node)
func (p *PeepholeReplaceKnownMethods) TryToFoldArrayConcat(n ast.Call) ast.Call {
	if !p.IsASTNormalized() || !p.useTypes {
		return n
	}
	concatFunctionCall, ok := PeepholeReplaceKnownMethods_CreateConcatFunctionCallForNode(n)
	if !ok {
		return n
	}
	concatFunctionCall, ok = p.TryToRemoveArrayLiteralFromFrontOfConcat(concatFunctionCall)
	jutil.CheckState(ok)
	return p.TryToFoldConcatChaining(concatFunctionCall)
}

// TryToRemoveArrayLiteralFromFrontOfConcat checks if we have this code
// pattern `[].concat(exactlyArrayArgument,...*)` and if yes replace empty
// array literal from the front of concatenation by the first argument of
// concat function call `[].concat(arr,1)` -> `arr.concat(1)`.
//
// Java: PeepholeReplaceKnownMethods.tryToRemoveArrayLiteralFromFrontOfConcat(ConcatFunctionCall)
func (p *PeepholeReplaceKnownMethods) TryToRemoveArrayLiteralFromFrontOfConcat(
	concatFunctionCall peepholeReplaceKnownMethods_ConcatFunctionCall) (peepholeReplaceKnownMethods_ConcatFunctionCall, bool) {
	callNode := concatFunctionCall.callNode
	arrayLiteralToRemove := concatFunctionCall.calleeNode
	if array, ok := arrayLiteralToRemove.AsArrayLit(); !ok || array.HasElements() {
		return concatFunctionCall, true
	}
	firstArgument, ok := callNode.FirstArg()
	if !ok {
		return concatFunctionCall, true
	}
	firstArg, ok := firstArgument.AsExpr() // not a spread
	if !ok || !PeepholeReplaceKnownMethods_ContainsExactlyArray(firstArg) {
		return concatFunctionCall, true
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
func (p *PeepholeReplaceKnownMethods) TryToFoldConcatChaining(concatFunctionCall peepholeReplaceKnownMethods_ConcatFunctionCall) ast.Call {
	concatCallNode := concatFunctionCall.callNode

	maybeFunctionCall, ok := concatFunctionCall.calleeNode.AsCall()
	if !ok {
		return concatCallNode
	}
	previousConcatFunctionCall, ok :=
		PeepholeReplaceKnownMethods_CreateConcatFunctionCallForNode(maybeFunctionCall)
	if !ok {
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

// PeepholeReplaceKnownMethods_CreateConcatFunctionCallForNode returns, if the
// argument node is a call to `Array.prototype.concat`, a
// `ConcatFunctionCall` value for it and true, otherwise false.
//
// Java: PeepholeReplaceKnownMethods.createConcatFunctionCallForNode(Node)
func PeepholeReplaceKnownMethods_CreateConcatFunctionCallForNode(n ast.Call) (peepholeReplaceKnownMethods_ConcatFunctionCall, bool) {
	callTarget, ok := n.Callee().AsGetProp()
	if !ok || callTarget.Node().GetStringID() != rkmIDConcat {
		return peepholeReplaceKnownMethods_ConcatFunctionCall{}, false
	}
	calleeNode := callTarget.Object()
	if !PeepholeReplaceKnownMethods_ContainsExactlyArray(calleeNode) {
		return peepholeReplaceKnownMethods_ConcatFunctionCall{}, false
	}
	return peepholeReplaceKnownMethods_ConcatFunctionCall{callNode: n, calleeNode: calleeNode}, true
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
		callee.Node().GetStringID() == rkmIDConcat &&
		PeepholeReplaceKnownMethods_ContainsExactlyArray(callee.Object())
}
