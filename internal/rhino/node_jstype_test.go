// Port of com/google/javascript/rhino/NodeTest.java @ closure-compiler v20260415

package rhino_test

// The NodeTest cases that attach JSTypes from a JSTypeRegistry. Go-only:
// language-adapter — they use rhinotesting, which imports rhino, so they
// live in the external test package; the other NodeTest cases are in
// node_test.go.

import (
	"testing"

	. "github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/rhino/rhinotesting"
)

func TestNode_CheckTreeTypeAwareEqualsSame(t *testing.T) {
	tr := rhinotesting.NewTruth(t)
	testErrorReporter := rhinotesting.NewTestErrorReporter()
	registry := NewJSTypeRegistry(testErrorReporter)
	node1 := Node_NewStringWithToken(Token_NAME, "f")
	node1.SetJSType(registry.GetNativeType(JSTypeNative_NUMBER_TYPE))
	node2 := Node_NewStringWithToken(Token_NAME, "f")
	node2.SetJSType(registry.GetNativeType(JSTypeNative_NUMBER_TYPE))
	tr.AssertThat(node1.IsEquivalentToTyped(node2)).IsTrue()
	testErrorReporter.VerifyHasEncounteredAllWarningsAndErrors(t)
}

func TestNode_CheckTreeTypeAwareEqualsDifferent(t *testing.T) {
	tr := rhinotesting.NewTruth(t)
	testErrorReporter := rhinotesting.NewTestErrorReporter()
	registry := NewJSTypeRegistry(testErrorReporter)
	node1 := Node_NewStringWithToken(Token_NAME, "f")
	node1.SetJSType(registry.GetNativeType(JSTypeNative_NUMBER_TYPE))
	node2 := Node_NewStringWithToken(Token_NAME, "f")
	node2.SetJSType(registry.GetNativeType(JSTypeNative_STRING_TYPE))
	tr.AssertThat(node1.IsEquivalentToTyped(node2)).IsFalse()
	testErrorReporter.VerifyHasEncounteredAllWarningsAndErrors(t)
}

func TestNode_CheckTreeTypeAwareEqualsDifferentNull(t *testing.T) {
	tr := rhinotesting.NewTruth(t)
	testErrorReporter := rhinotesting.NewTestErrorReporter()
	registry := NewJSTypeRegistry(testErrorReporter)
	node1 := Node_NewStringWithToken(Token_NAME, "f")
	node1.SetJSType(registry.GetNativeType(JSTypeNative_NUMBER_TYPE))
	node2 := Node_NewStringWithToken(Token_NAME, "f")
	tr.AssertThat(node1.IsEquivalentToTyped(node2)).IsFalse()
	testErrorReporter.VerifyHasEncounteredAllWarningsAndErrors(t)
}

// Java: NodeTest.bitsetFromNodeProperties(ImmutableSet<NodeProperty>)
func nodeTest_bitsetFromNodeProperties(props ...NodeProperty) int64 {
	var bitset int64
	for _, prop := range props {
		bitset = Node_SetNodePropertyBit(bitset, prop)
	}
	return bitset
}

func TestNode_SerializeProperties_TypeBeforeCast(t *testing.T) {
	tr := rhinotesting.NewTruth(t)
	testErrorReporter := rhinotesting.NewTestErrorReporter()
	registry := NewJSTypeRegistry(testErrorReporter)
	node := Node_NewStringWithToken(Token_NAME, "f")
	node.SetJSTypeBeforeCast(registry.GetNativeType(JSTypeNative_NUMBER_TYPE))
	result := node.SerializeProperties()
	// Special case: Rhino node prop TYPE_BEFORE_CAST is converted to NodeProperty.COLOR_FROM_CAST
	tr.AssertThat(result).
		IsEqualTo(nodeTest_bitsetFromNodeProperties(NodeProperty_COLOR_FROM_CAST))
}
