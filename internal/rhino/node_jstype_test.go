// Port of com/google/javascript/rhino/NodeTest.java @ closure-compiler v20261021

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Nick Santos
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

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
