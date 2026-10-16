// Port of com/google/javascript/jscomp/NodeNameExtractor.java
// @ closure-compiler v20260415.
//
// Utility that extracts the qualified name out of a node. Useful
// when producing a human-friendly string representation of a
// property node, e.g. to describe the node or to name related
// nodes (as done by the NameAnonymousFunctions compiler pass).
//
// NodeUtil.getName and NodeUtil.getStringValue are inlined here
// as nodeUtilGetName / nodeUtilGetStringValue since NodeUtil.java
// is not yet ported. TokenStream.isJSIdentifier is inlined as
// isJSIdentifier. TODO: replace with NodeUtil / TokenStream ports.

package jscomp

import (
	"fmt"
	"unicode"

	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
)

// NodeNameExtractor extracts a qualified name from an AST node.
// The delimiter character replaces dots and brackets, giving a
// flat identifier-like string. ".prototype" is replaced with a
// bare delimiter to keep names short while still distinguishing
// prototype properties from static ones.
//
// Mirrors Java's NodeNameExtractor class.
type NodeNameExtractor struct {
	delimiter     rune
	nextUniqueInt int
}

// NewNodeNameExtractor constructs a NodeNameExtractor with the
// given delimiter. Mirrors Java's NodeNameExtractor(char delimiter).
func NewNodeNameExtractor(delimiter rune) *NodeNameExtractor {
	return &NodeNameExtractor{delimiter: delimiter}
}

// GetName returns a qualified name for the specified node. Dots
// and brackets are replaced by the delimiter. ".prototype" is
// collapsed to a bare delimiter. Mirrors Java's getName(Node).
func (e *NodeNameExtractor) GetName(node *rhino.Node) string {
	switch node.Token() {
	case rhino.CLASS, rhino.FUNCTION:
		return nodeUtilGetName(node)

	case rhino.GETPROP:
		lhsOfDot := node.FirstChild()
		rhsOfDot := lhsOfDot.Next()
		lhsOfDotName := e.GetName(lhsOfDot)
		rhsOfDotName := e.GetName(rhsOfDot)
		if rhsOfDotName == "prototype" {
			return lhsOfDotName + string(e.delimiter)
		}
		return lhsOfDotName + string(e.delimiter) + rhsOfDotName

	case rhino.GETELEM:
		outsideBrackets := node.FirstChild()
		insideBrackets := outsideBrackets.Next()
		nameOutsideBrackets := e.GetName(outsideBrackets)
		nameInsideBrackets := e.GetName(insideBrackets)
		if nameInsideBrackets == "prototype" {
			return nameOutsideBrackets + string(e.delimiter)
		}
		return nameOutsideBrackets + string(e.delimiter) + nameInsideBrackets

	case rhino.NAME:
		return node.String_()

	case rhino.STRINGLIT, rhino.STRING_KEY, rhino.MEMBER_FUNCTION_DEF:
		s := node.String_()
		if isJSIdentifier(s) {
			return s
		}
		name := fmt.Sprintf("__%d", e.nextUniqueInt)
		e.nextUniqueInt++
		return name

	case rhino.NUMBER:
		return nodeUtilGetStringValue(node)

	case rhino.THIS:
		return "this"

	case rhino.CALL, rhino.COMPUTED_PROP:
		return e.GetName(node.FirstChild())

	default:
		sb := ""
		for child := node.FirstChild(); child != nil; child = child.Next() {
			if sb != "" {
				sb += string(e.delimiter)
			}
			sb += e.GetName(child)
		}
		return sb
	}
}

// nodeUtilGetName returns the name of a FUNCTION or CLASS node, or ""
// when the node is anonymous. Inlines NodeUtil.getName(Node) from Java,
// which delegates to NodeUtil.getNameNode(Node).getQualifiedName().
// TODO: replace with NodeUtil port.
func nodeUtilGetName(n *rhino.Node) string {
	nameNode := nodeUtilGetNameNode(n)
	if nameNode == nil {
		return ""
	}
	return qualifiedName(nameNode)
}

// nodeUtilGetNameNode returns the name-bearing node of a FUNCTION or
// CLASS, mirroring NodeUtil.getNameNode in Java.
// TODO: replace with NodeUtil port.
func nodeUtilGetNameNode(n *rhino.Node) *rhino.Node {
	parent := n.Parent()
	if parent == nil {
		// Standalone function/class declaration — name is first child.
		first := n.FirstChild()
		if first == nil || first.String_() == "" {
			return nil
		}
		return first
	}
	switch parent.Token() {
	case rhino.NAME:
		// var name = function() ...
		return parent
	case rhino.ASSIGN:
		// qualified.name = function() ...
		first := parent.FirstChild()
		if first != nil && isQualifiedName(first) {
			return first
		}
		return nil
	default:
		// function name() ... or class Name ...
		// Anonymous functions have a NAME node with empty string;
		// anonymous classes may have no first child or an empty node.
		// Mirrors: funNameNode.isEmpty() || funNameNode.getString().isEmpty()
		first := n.FirstChild()
		if first == nil || first.String_() == "" {
			return nil
		}
		return first
	}
}

// isQualifiedName reports whether n is a NAME or a GETPROP chain of
// NAMEs (e.g. a.b.c). Mirrors Node.isQualifiedName in Java.
func isQualifiedName(n *rhino.Node) bool {
	switch n.Token() {
	case rhino.NAME:
		return true
	case rhino.GETPROP:
		kids := n.Children()
		return len(kids) == 2 && isQualifiedName(kids[0])
	}
	return false
}

// nodeUtilGetStringValue returns the string representation of a NUMBER
// node. Inlines the NUMBER branch of NodeUtil.getStringValue(Node)
// which uses DToA.numberToString. For integers we omit the decimal
// point; for other values we use Go's default float formatting.
// TODO: replace with NodeUtil / DToA port.
func nodeUtilGetStringValue(n *rhino.Node) string {
	v := n.Number()
	if v == float64(int64(v)) {
		return fmt.Sprintf("%d", int64(v))
	}
	return fmt.Sprintf("%g", v)
}

// isJSIdentifier reports whether s is a valid JS identifier. Inlines
// Java's TokenStream.isJSIdentifier(String) → JSIdentifier.isJSIdentifier.
// A JS identifier has length > 0, starts with a character where
// Character.isJavaIdentifierStart is true, and all subsequent characters
// satisfy Character.isJavaIdentifierPart — and none are
// Character.isIdentifierIgnorable. In practice: first rune must be a
// letter, underscore, or dollar sign; subsequent runes may also be
// digits; zero-width joiners and other ignorable code points are rejected.
// TODO: replace with TokenStream port.
func isJSIdentifier(s string) bool {
	if s == "" {
		return false
	}
	runes := []rune(s)
	first := runes[0]
	if unicode.Is(unicode.Mn, first) || // ignorable (approximation)
		!(unicode.IsLetter(first) || first == '_' || first == '$') {
		return false
	}
	for _, r := range runes[1:] {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '$') {
			return false
		}
	}
	return true
}
