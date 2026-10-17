// Port of com/google/javascript/jscomp/NodeNameExtractor.java @ closure-compiler v20260415

/*
 * Copyright 2004 The Closure Compiler Authors.
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
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// NodeNameExtractor is a utility class that extracts the qualified name out
// of a node. Useful when trying to get a human-friendly string
// representation of a property node that can be used to describe the node
// or name related nodes based on it (as done by the NameAnonymousFunctions
// compiler pass).
type NodeNameExtractor struct {
	delimiter     uint16
	nextUniqueInt int32
}

// Java: NodeNameExtractor(char)
func NewNodeNameExtractor(delimiter uint16) *NodeNameExtractor {
	return &NodeNameExtractor{delimiter: delimiter, nextUniqueInt: 0}
}

// GetName returns a qualified name of the specified node. Dots and brackets
// are changed to the delimiter passed in when constructing the
// NodeNameExtractor object. We also replace ".prototype" with the delimiter
// to keep names short, while still differentiating them from static
// properties. (Prototype properties will end up looking like "a$b$$c" if
// this.delimiter = '$'.)
//
// Go-only: language-adapter — Java's String result is never null except
// through NodeUtil.getName / getStringValue, whose null is concatenated as
// "null" by the callers' string concatenation; GetName returns *string so
// that the null of those cases is returned as Java returns it.
//
// Java: NodeNameExtractor.getName(Node)
func (e *NodeNameExtractor) GetName(node *rhino.Node) *string {
	switch node.GetToken() {
	case rhino.Token_CLASS, rhino.Token_FUNCTION:
		return NodeUtil_GetName(node)
	case rhino.Token_GETPROP:
		lhsOfDot := node.GetFirstChild()
		rhsOfDot := lhsOfDot.GetNext()
		lhsOfDotName := e.GetName(lhsOfDot)
		rhsOfDotName := e.GetName(rhsOfDot)
		if rhsOfDotName != nil && "prototype" == *rhsOfDotName {
			return jutil.StringPtr(jutil.ToJavaString(lhsOfDotName) + jutil.StringFromChar(e.delimiter))
		} else {
			return jutil.StringPtr(jutil.ToJavaString(lhsOfDotName) + jutil.StringFromChar(e.delimiter) + jutil.ToJavaString(rhsOfDotName))
		}
	case rhino.Token_GETELEM:
		outsideBrackets := node.GetFirstChild()
		insideBrackets := outsideBrackets.GetNext()
		nameOutsideBrackets := e.GetName(outsideBrackets)
		nameInsideBrackets := e.GetName(insideBrackets)
		if nameInsideBrackets != nil && "prototype" == *nameInsideBrackets {
			return jutil.StringPtr(jutil.ToJavaString(nameOutsideBrackets) + jutil.StringFromChar(e.delimiter))
		} else {
			return jutil.StringPtr(jutil.ToJavaString(nameOutsideBrackets) + jutil.StringFromChar(e.delimiter) + jutil.ToJavaString(nameInsideBrackets))
		}
	case rhino.Token_NAME:
		return jutil.StringPtr(node.GetString())
	case rhino.Token_STRINGLIT, rhino.Token_STRING_KEY, rhino.Token_MEMBER_FUNCTION_DEF:
		if rhino.TokenStream_IsJSIdentifier(node.GetString()) {
			return jutil.StringPtr(node.GetString())
		}
		s := "__" + strconv.FormatInt(int64(e.nextUniqueInt), 10)
		e.nextUniqueInt++
		return &s
	case rhino.Token_NUMBER:
		return NodeUtil_GetStringValue(node)
	case rhino.Token_THIS:
		return jutil.StringPtr("this")
	case rhino.Token_CALL, rhino.Token_COMPUTED_PROP:
		return e.GetName(node.GetFirstChild())
	default:
		sb := jutil.NewStringBuilder()
		for child := node.GetFirstChild(); child != nil; child = child.GetNext() {
			if sb.Length() > 0 {
				sb.AppendChar(e.delimiter)
			}
			sb.Append(jutil.ToJavaString(e.GetName(child)))
		}
		return jutil.StringPtr(sb.String())
	}
}
