// Port of com/google/javascript/rhino/Node.java @ closure-compiler v20261021

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
 *   Norris Boyd
 *   Roger Lawrence
 *   Mike McCabe
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"io"
	"math"
	"math/big"
	"strconv"
	"strings"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// Java: Node.getDouble()
func (n Node) GetDouble() float64 {
	n.node_castKind(node_kindNumber)
	return n.node_number()
}

// Java: Node.setDouble(double)
func (n Node) SetDouble(x float64) {
	jutil.CheckStateMsg1(!math.IsNaN(x), "%s", x)
	jutil.CheckStateMsg1(base.JSCompDoubles_IsPositive(x), "%s", x)
	n.node_castKind(node_kindNumber)
	n.rec().payload = numberPayload(x)
}

var (
	rhinoBigIntMu   sync.RWMutex
	rhinoBigIntPool []*big.Int
)

func rhinoInternBigInt(v *big.Int) uint32 {
	rhinoBigIntMu.Lock()
	rhinoBigIntPool = append(rhinoBigIntPool, v)
	id := uint32(len(rhinoBigIntPool))
	rhinoBigIntMu.Unlock()
	return id
}

func rhinoBigIntFromID(id uint32) *big.Int {
	if id == 0 {
		return nil
	}
	rhinoBigIntMu.RLock()
	v := rhinoBigIntPool[id-1]
	rhinoBigIntMu.RUnlock()
	return v
}

// Java: Node.getBigInt()
func (n Node) GetBigInt() *big.Int {
	n.node_castKind(node_kindBigInt)
	return rhinoBigIntFromID(uint32(n.rec().payload))
}

// Java: Node.setBigInt(BigInteger)
func (n Node) SetBigInt(number *big.Int) {
	jutil.CheckNotNull(number)
	jutil.CheckStateMsg1(number.Sign() >= 0, "%s", number)
	n.node_castKind(node_kindBigInt)
	n.rec().payload = rhinoInternBigInt(number)
}

// Java: Node.getString()
func (n Node) GetString() string {
	n.node_castKind(node_kindString)
	return n.nodeStringValue()
}

//go:noinline
func (n Node) failGetStringID() {
	if got := n.nodeKind(); got != node_kindString {
		failNodeCastKind(got, node_kindString)
	}
	panic(jutil.NewNullPointerException())
}

//go:noinline
func (n Node) failGetRawStringID() {
	if got := n.nodeKind(); got != node_kindTemplateLiteralSubstring {
		failNodeCastKind(got, node_kindTemplateLiteralSubstring)
	}
	panic(jutil.NewNullPointerException())
}

// GetStringID returns the 32-bit interned string ID of a StringNode.
func (n Node) GetStringID() uint32 {
	id := uint32(n.rec().payload)
	if n.rec().kindFlags&kindMask != uint8(node_kindString) || id == 0 {
		n.failGetStringID()
	}
	return id
}

// Java: Node.setString(String)
func (n Node) SetString(str string) {
	n.node_castKind(node_kindString)
	n.setStringPayload(rhinoStringPoolIntern(str))
}

// Java: Node.setStringFromStringPool(RhinoStringPool.LazyInternedStringList, int)
func (n Node) SetStringFromStringPool(stringPool *RhinoStringPool_LazyInternedStringList, offset int) {
	n.node_castKind(node_kindString)
	n.setStringPayload(stringPool.GetID(offset))
}

// Java: Node.getRawString()
func (n Node) GetRawString() string {
	n.node_castKind(node_kindTemplateLiteralSubstring)
	return n.nodeStringValue() // raw
}

// GetRawStringID returns the 32-bit interned raw string ID of a TemplateLiteralSubstringNode.
func (n Node) GetRawStringID() uint32 {
	id := uint32(n.rec().payload)
	if n.rec().kindFlags&kindMask != uint8(node_kindTemplateLiteralSubstring) || id == 0 {
		n.failGetRawStringID()
	}
	return id
}

// Java: Node.getCookedString()
func (n Node) GetCookedString() *string {
	n.node_castKind(node_kindTemplateLiteralSubstring)
	return n.node_cooked()
}

// GetCookedStringID returns the 32-bit interned cooked string ID of a TemplateLiteralSubstringNode, or 0 if unset.
func (n Node) GetCookedStringID() uint32 {
	n.node_castKind(node_kindTemplateLiteralSubstring)
	return uint32(n.node_cookedID())
}

// Java: Node.toString()
func (n Node) String() string {
	return n.ToStringWithFlags(true, true, true)
}

// Java: Node.toString(boolean, boolean, boolean)
func (n Node) ToStringWithFlags(printSource bool, printAnnotations bool, printType bool) string {
	var sb strings.Builder
	n.ToStringInto(&sb, printSource, printAnnotations, printType)
	return sb.String()
}

// Java: Node.toString(StringBuilder, boolean, boolean, boolean)
func (n Node) ToStringInto(sb *strings.Builder, printSource bool, printAnnotations bool, printType bool) {
	sb.WriteString(n.rec().token.String())
	if n.nodeKind() == node_kindString {
		sb.WriteByte(' ')
		sb.WriteString(n.GetString())
	} else if n.rec().token == Token_FUNCTION {
		sb.WriteByte(' ')
		// In the case of JsDoc trees, the first child is often not a string
		// which causes exceptions to be thrown when calling toString or
		// toStringTree.
		if n.rec().first == 0 || n.firstNode().rec().token != Token_NAME {
			sb.WriteString("<invalid>")
		} else {
			sb.WriteString(n.firstNode().GetString())
		}
	} else if n.rec().token == Token_NUMBER {
		sb.WriteByte(' ')
		sb.WriteString(jutil.DoubleToString(n.GetDouble()))
	}
	if printSource {
		lineno := n.GetLineno()
		if lineno != -1 {
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(lineno))
			sb.WriteByte(':')
			sb.WriteString(strconv.Itoa(n.GetCharno()))
			sb.WriteByte(' ')
		}
		if length := n.GetLength(); length != 0 {
			sb.WriteString(" [length: ")
			sb.WriteString(strconv.Itoa(length))
			sb.WriteByte(']')
		}
	}

	if printAnnotations {
		keys := n.GetSortedPropTypes()
		for i := 0; i < len(keys); i++ {
			type_ := Node_propValues[keys[i]]
			if type_ == Node_Prop_CLOSURE_UNAWARE_SHADOW {
				sb.WriteString(" [is_shadow_host]")
				continue
			}
			x := n.LookupProperty(type_)
			sb.WriteString(" [")
			sb.WriteString(jutil.StringToLowerCaseASCII(type_.String()))
			sb.WriteString(": ")
			sb.WriteString(jutil.ToJavaString(x))
			sb.WriteByte(']')
		}
		if origID := n.node_originalNameID(); origID != 0 {
			sb.WriteString(" [original_name: ")
			sb.WriteString(rhinoStringFromID(origID))
			sb.WriteByte(']')
		}
	}

	if printType && n.tc() != 0 {
		typeString := jutil.ToJavaString(n.typeOrColorValue())
		sb.WriteString(" : ")
		sb.WriteString(typeString)
	}
}

// Java: Node.createJsonPair(String, String)
func Node_CreateJsonPair(name string, value string) string {
	return "\"" + name + "\":\"" + value + "\""
}

// Java: Node.createJsonPairRawValue(String, String)
func Node_CreateJsonPairRawValue(name string, value string) string {
	return "\"" + name + "\":" + value
}

// node_appender is Java's Appendable over an io.StringWriter: write errors
// become the IOException Appendable.append throws.
type node_appender struct{ w io.StringWriter }

func (a node_appender) append(s string) {
	if _, err := a.w.WriteString(s); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
}

// ToJson writes the JSON form of this node and its subtree.
//
// Java: Node.toJson(Appendable)
func (n Node) ToJson(appendable io.StringWriter) {
	sb := node_appender{appendable}
	sb.append("{")
	sb.append(Node_CreateJsonPair("token", n.rec().token.String()))
	if n.nodeKind() == node_kindString {
		sb.append(",")
		sb.append(Node_CreateJsonPair("string", n.GetString()))
	} else if n.rec().token == Token_FUNCTION {
		sb.append(",")
		// In the case of JsDoc trees, the first child is often not a string
		// which causes exceptions to be thrown when calling toString or
		// toStringTree.
		if n.rec().first == 0 || n.firstNode().rec().token != Token_NAME {
			sb.append(Node_CreateJsonPair("functionName", "<invalid>"))
		} else {
			sb.append(Node_CreateJsonPair("functionName", n.firstNode().GetString()))
		}
	} else if n.rec().token == Token_NUMBER {
		sb.append(",")
		sb.append(Node_CreateJsonPair("number", jutil.DoubleToString(n.GetDouble())))
	}
	lineno := n.GetLineno()
	if lineno != -1 {
		sb.append(",")
		sb.append("\"sourceLocation\":{")
		sb.append(Node_CreateJsonPairRawValue("line", strconv.Itoa(lineno)))
		sb.append(",")
		sb.append(Node_CreateJsonPairRawValue("col", strconv.Itoa(n.GetCharno())))
		if length := n.GetLength(); length != 0 {
			sb.append(",")
			sb.append(Node_CreateJsonPairRawValue("length", strconv.Itoa(length)))
		}
		sb.append("}")
	}

	if origID := n.node_originalNameID(); origID != 0 {
		sb.append(",")
		sb.append(Node_CreateJsonPair("original_name", rhinoStringFromID(origID)))
	}

	keys := n.GetSortedPropTypes()
	if len(keys) != 0 {
		sb.append(",")
		sb.append("\"props\":{")
		for i := 0; i < len(keys); i++ {
			type_ := Node_propValues[keys[i]]
			x := n.LookupProperty(type_)
			sb.append(
				Node_CreateJsonPair(
					jutil.StringToLowerCaseASCII(type_.String()),
					jutil.StringReplace(jutil.StringReplace(x.String(), "\n", "\\n"), "\"", "\\\"")))
			if i+1 < len(keys) {
				sb.append(",")
			}
		}
		sb.append("}")
	}

	if n.tc() != 0 {
		typeString := jutil.ToJavaString(n.typeOrColorValue())
		sb.append(",")
		sb.append(Node_CreateJsonPair("typeString", jutil.StringReplace(typeString, "\"", "\\\"")))
	}

	if n.rec().first != 0 {
		sb.append(",")
		sb.append("\"children\":[")
		for child := n.firstNode(); !child.IsNil(); child = child.nextNode() {
			child.ToJson(appendable)
			if child.rec().next != 0 {
				sb.append(",")
			}
		}
		sb.append("]")
	}
	sb.append("}")
}

// Java: Node.toStringTree()
func (n Node) ToStringTree() string {
	return n.ToStringTreeImpl()
}

// Java: Node.toStringTreeImpl()
func (n Node) ToStringTreeImpl() string {
	var s strings.Builder
	n.AppendStringTree(&s)
	return s.String()
}

// Java: Node.appendStringTree(Appendable)
func (n Node) AppendStringTree(appendable io.StringWriter) {
	Node_ToStringTreeHelper(n, 0, appendable)
}

// Java: Node.toStringTreeHelper(Node, int, Appendable)
func Node_ToStringTreeHelper(n Node, level int, appendable io.StringWriter) {
	sb := node_appender{appendable}
	for i := 0; i != level; i++ {
		sb.append("    ")
	}
	sb.append(n.String())
	sb.append("\n")
	for cursor := n.firstNode(); !cursor.IsNil(); cursor = cursor.nextNode() {
		Node_ToStringTreeHelper(cursor, level+1, appendable)
	}
	shadow := n.GetClosureUnawareShadow()
	if !shadow.IsNil() {
		Node_ToStringTreeHelper(shadow, level+1, appendable)
	}
}

// Java: Node.appendJsonTree(Appendable)
func (n Node) AppendJsonTree(appendable io.StringWriter) {
	Node_ToJsonTreeHelper(n, appendable)
}

// Java: Node.toJsonTreeHelper(Node, Appendable)
func Node_ToJsonTreeHelper(n Node, appendable io.StringWriter) {
	n.ToJson(appendable)
}

// node_number is the NumberNode.number field.
func (n Node) node_number() float64 {
	return payloadNumber(n.rec().payload)
}

// node_cooked is the TemplateLiteralSubstringNode.cooked field.
func (n Node) node_cooked() *string {
	return rhinoStringPtrFromID(n.node_cookedID())
}

// ==========================================================================
// Source position management

// Java: Node.setStaticSourceFileFrom(Node)
func (n Node) SetStaticSourceFileFrom(other Node) {
	n.rec().sourceFile = other.rec().sourceFile
}

// Java: Node.setStaticSourceFile(StaticSourceFile)
func (n Node) SetStaticSourceFile(file StaticSourceFile) Node {
	n.PutProp(Node_Prop_SOURCE_FILE, file)
	return n
}

// SetSourceFileForTesting sets the source file to a non-extern file of the
// given name.
//
// Java: Node.setSourceFileForTesting(String)
func (n Node) SetSourceFileForTesting(name string) {
	n.PutProp(Node_Prop_SOURCE_FILE, NewSimpleSourceFile(name, StaticSourceFile_SourceKind_STRONG))
}

// TODO(johnlenz): make this final
//
// Java: Node.getSourceFileName()
func (n Node) GetSourceFileName() *string {
	sfID := n.rec().sourceFile
	if sfID == 0 {
		return nil
	}
	idx := sfID - 1
	if p := sourceFileNamePtrChunks[idx>>10].Load()[idx&1023]; p != nil {
		return p
	}
	_ = node_cast[StaticSourceFile](sourceFilePropItem(sfID).objectValue)
	return nil
}

// GetStaticSourceFile returns the source file associated with this input.
//
// Java: Node.getStaticSourceFile()
func (n Node) GetStaticSourceFile() StaticSourceFile {
	sfID := n.rec().sourceFile
	if sfID == 0 {
		return nil
	}
	return node_cast[StaticSourceFile](sourceFilePropItem(sfID).objectValue)
}

// SetInputId sets the ID of the input this Node came from.
//
// Java: Node.setInputId(InputId)
func (n Node) SetInputId(inputId *InputId) {
	n.PutProp(Node_Prop_INPUT_ID, inputId)
}

// GetInputId returns the Id of the CompilerInput associated with this Node.
//
// Java: Node.getInputId()
func (n Node) GetInputId() *InputId {
	return node_cast[*InputId](n.GetProp(Node_Prop_INPUT_ID))
}

// GetOriginalName returns the original name of this node, if the node has
// been renamed.
//
// Deprecated: "original name" is poorly defined.
//
// Java: Node.getOriginalName()
func (n Node) GetOriginalName() *string {
	return rhinoStringPtrFromID(n.node_originalNameID())
}

// GetOriginalNameID returns the 32-bit interned originalName ID of n, or 0 if unset.
func (n Node) GetOriginalNameID() uint32 {
	return uint32(n.node_originalNameID())
}

func (n Node) getOriginalNameString() (string, bool) {
	id := n.node_originalNameID()
	if id == 0 {
		return "", false
	}
	return rhinoStringFromID(id), true
}

// Java: Node.setOriginalName(String)
func (n Node) SetOriginalName(s *string) {
	if s == nil {
		n.node_setOriginalNameID(0)
	} else {
		n.node_setOriginalNameID(rhinoStringPoolIntern(*s))
	}
}

// SetOriginalNameString interns s in RhinoStringPool and sets it as the node's
// original name without boxing s on the heap.
func (n Node) SetOriginalNameString(s string) {
	n.node_setOriginalNameID(rhinoStringPoolIntern(s))
}

// Java: Node.setOriginalNameFromStringPool(RhinoStringPool.LazyInternedStringList, int)
func (n Node) SetOriginalNameFromStringPool(stringPool *RhinoStringPool_LazyInternedStringList, offset int) {
	n.node_setOriginalNameID(stringPool.GetID(offset))
}

// SetOriginalNameFromName copies the original name.
//
// Java: Node.setOriginalNameFromName(Node)
func (n Node) SetOriginalNameFromName(name Node) {
	name.node_castKind(node_kindString)
	n.node_setOriginalNameID(rhinoStringID(name.rec().payload))
}

// IsIndexable reports whether this node should be indexed by static
// analysis / code indexing tools.
//
// Java: Node.isIndexable()
func (n Node) IsIndexable() bool {
	return !n.GetBooleanProp(Node_Prop_NON_INDEXABLE)
}

// Java: Node.makeNonIndexable()
func (n Node) MakeNonIndexable() {
	n.PutBooleanProp(Node_Prop_NON_INDEXABLE, true)
}

// Java: Node.makeNonIndexableRecursive()
func (n Node) MakeNonIndexableRecursive() {
	n.MakeNonIndexable()
	for child := n.firstNode(); !child.IsNil(); child = child.GetNext() {
		child.MakeNonIndexableRecursive()
	}
}

// Java: Node.isFromExterns()
func (n Node) IsFromExterns() bool {
	file := n.GetStaticSourceFile()
	if file == nil {
		return false
	}
	return file.IsExtern()
}

// IsClosureUnawareCode indicates that this node is for source that was
// written without any understanding of the compiler's various type-checking
// or optimization limitations, and should be optimized very cautiously by
// the compiler.
//
// This attribute is per source file, not per node / subset of an source
// file's AST.
//
// Java: Node.isClosureUnawareCode()
func (n Node) IsClosureUnawareCode() bool {
	file := n.GetStaticSourceFile()
	return file != nil && file.IsClosureUnawareCode()
}

// Java: Node.getLength()
func (n Node) GetLength() int {
	if n.rec().length16 != 0xFFFF {
		return int(n.rec().length16)
	}
	if box := n.extBoxOrNil(); box != nil {
		return int(box.length)
	}
	return 0
}

// Java: Node.setLength(int)
func (n Node) SetLength(length int) {
	if length >= 0 && length < 0xFFFF {
		wasInBox := n.rec().length16 == 0xFFFF
		n.rec().length16 = uint16(length)
		if wasInBox {
			n.maybeCollapseExtBox(n.extBox())
		}
		return
	}
	n.ensureExtBox().length = int32(length)
	n.rec().length16 = 0xFFFF
}

func (n Node) copyLengthFrom(other Node) {
	if other.rec().length16 != 0xFFFF && n.rec().length16 != 0xFFFF {
		n.rec().length16 = other.rec().length16
		return
	}
	n.SetLength(other.GetLength())
}

// Java: Node.getLineno()
func (n Node) GetLineno() int {
	if n.rec().linenoCharno == -1 {
		return -1
	}
	return int(jutil.Ushr32(n.rec().linenoCharno, Node_CHARNO_BITS))
}

// GetCharno returns the 0-based column number.
//
// Java: Node.getCharno()
func (n Node) GetCharno() int {
	if n.rec().linenoCharno == -1 {
		return -1
	}
	return int(n.rec().linenoCharno & Node_MAX_COLUMN_NUMBER)
}

// Java: Node.getLocation()
func (n Node) GetLocation() string {
	return jutil.ToJavaString(n.GetSourceFileName()) + ":" + strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno())
}

// TODO(johnlenz): make this final
//
// Java: Node.getSourceOffset()
func (n Node) GetSourceOffset() int {
	file := n.GetStaticSourceFile()
	if file == nil {
		return -1
	}
	lineno := n.GetLineno()
	if lineno == -1 {
		return -1
	}
	return int(int32(file.GetLineOffset(lineno)) + int32(n.GetCharno()))
}

// Java: Node.getSourcePosition()
func (n Node) GetSourcePosition() int {
	return int(n.rec().linenoCharno)
}

// Node_CHARNO_BITS represents how many of the lower-order bits of
// linenoCharno are reserved for storing the column number. Bits above these
// store the line number. This gives us decent position information for
// everything except files already passed through a minimizer, where lines
// might be longer than 4096 characters.
const Node_CHARNO_BITS = 12

// Node_MAX_COLUMN_NUMBER represents the maximum column number that can be
// represented. JSCompiler's modifications to Rhino cause all tokens located
// beyond the maximum column to MAX_COLUMN_NUMBER.
const Node_MAX_COLUMN_NUMBER = (1 << Node_CHARNO_BITS) - 1

// SetLinenoCharno merges the line number and character number in one
// integer.
//
// The charno takes the first 12 bits and the line number takes the rest. If
// the charno is greater than (2^12)-1 it is adjusted to (2^12)-1
//
// Java: Node.setLinenoCharno(int, int)
func (n Node) SetLinenoCharno(lineno int, charno int) Node {
	if lineno < 0 || charno < 0 {
		n.rec().linenoCharno = -1
		return n
	}

	if charno > Node_MAX_COLUMN_NUMBER {
		charno = Node_MAX_COLUMN_NUMBER
	}
	n.rec().linenoCharno = jutil.Shl32(int32(lineno), Node_CHARNO_BITS) | int32(charno)

	return n
}

// Node_RecursionMode is the recursion mode for isEquivalentTo.
//
// Java: Node.RecursionMode
type Node_RecursionMode int32

const (
	// Do not compare children or shadow AST.
	Node_RecursionMode_SHALLOW Node_RecursionMode = iota
	// Recursively compare children but not shadow AST.
	Node_RecursionMode_DEEP_NO_SHADOW
	// Recursively compare children and shadow AST.
	Node_RecursionMode_DEEP
)

var node_RecursionMode_names = [...]string{"SHALLOW", "DEEP_NO_SHADOW", "DEEP"}

// String is Enum.name().
func (m Node_RecursionMode) String() string { return node_RecursionMode_names[m] }

// Java: Node.RecursionMode.recurseShadowAst()
func (m Node_RecursionMode) RecurseShadowAst() bool {
	return m == Node_RecursionMode_DEEP
}

// Java: Node.RecursionMode.recurseChildren()
func (m Node_RecursionMode) RecurseChildren() bool {
	return m == Node_RecursionMode_DEEP || m == Node_RecursionMode_DEEP_NO_SHADOW
}

// Node_TypeComparison is the type comparison mode for isEquivalentTo.
//
// Java: Node.TypeComparison
type Node_TypeComparison int32

const (
	// Ignore JSTypes/Colors.
	Node_TypeComparison_IGNORE Node_TypeComparison = iota
	// Compare JSTypes/Colors.
	Node_TypeComparison_COMPARE
)

// String is Enum.name().
func (c Node_TypeComparison) String() string { return [...]string{"IGNORE", "COMPARE"}[c] }

// Node_JsDocComparison is the JsDoc comparison mode for isEquivalentTo.
//
// Java: Node.JsDocComparison
type Node_JsDocComparison int32

const (
	// Ignore JSDoc.
	Node_JsDocComparison_IGNORE Node_JsDocComparison = iota
	// Compare JSDoc.
	Node_JsDocComparison_COMPARE
)

// String is Enum.name().
func (c Node_JsDocComparison) String() string { return [...]string{"IGNORE", "COMPARE"}[c] }

// Node_SideEffectComparison is the side effect comparison mode for
// isEquivalentTo.
//
// Java: Node.SideEffectComparison
type Node_SideEffectComparison int32

const (
	// Ignore side effect flags.
	Node_SideEffectComparison_IGNORE Node_SideEffectComparison = iota
	// Compare side effect flags.
	Node_SideEffectComparison_COMPARE
)

// String is Enum.name().
func (c Node_SideEffectComparison) String() string { return [...]string{"IGNORE", "COMPARE"}[c] }

// IsEquivalentToShallow checks equivalence without going into child nodes.
//
// Java: Node.isEquivalentToShallow(Node)
func (n Node) IsEquivalentToShallow(node Node) bool {
	return n.IsEquivalentToWithModes(
		node,
		Node_RecursionMode_SHALLOW,
		Node_TypeComparison_IGNORE,
		Node_JsDocComparison_IGNORE,
		Node_SideEffectComparison_IGNORE)
}

// IsEquivalentWithSideEffectsTo returns true if this node is equivalent
// semantically to another including side effects.
//
// Java: Node.isEquivalentWithSideEffectsTo(Node)
func (n Node) IsEquivalentWithSideEffectsTo(node Node) bool {
	return n.IsEquivalentToWithModes(
		node,
		Node_RecursionMode_DEEP,
		Node_TypeComparison_IGNORE,
		Node_JsDocComparison_IGNORE,
		Node_SideEffectComparison_COMPARE)
}

// IsEquivalentWithSideEffectsToShallow returns true if this node is
// equivalent semantically to another including side effects.
//
// Java: Node.isEquivalentWithSideEffectsToShallow(Node)
func (n Node) IsEquivalentWithSideEffectsToShallow(node Node) bool {
	return n.IsEquivalentToWithModes(
		node,
		Node_RecursionMode_SHALLOW,
		Node_TypeComparison_IGNORE,
		Node_JsDocComparison_IGNORE,
		Node_SideEffectComparison_COMPARE)
}

// IsEquivalentToTyped returns true if this node is equivalent semantically
// to another and the types are equivalent.
//
// Java: Node.isEquivalentToTyped(Node)
func (n Node) IsEquivalentToTyped(node Node) bool {
	return n.IsEquivalentToWithModes(
		node,
		Node_RecursionMode_DEEP,
		Node_TypeComparison_COMPARE,
		Node_JsDocComparison_COMPARE,
		Node_SideEffectComparison_IGNORE)
}

// IsEquivalentTo returns true if this node is equivalent semantically to
// another.
//
// Java: Node.isEquivalentTo(Node)
func (n Node) IsEquivalentTo(node Node) bool {
	return n.IsEquivalentToWithModes(
		node,
		Node_RecursionMode_DEEP,
		Node_TypeComparison_IGNORE,
		Node_JsDocComparison_IGNORE,
		Node_SideEffectComparison_IGNORE)
}

// IsEquivalentToWithFlags returns whether this node is equivalent
// semantically to the provided node.
//
// compareType: whether to compare the JSTypes of the nodes. recurse: whether
// to compare the children of the current node. If not, only the count of
// the children are compared. jsDoc: whether to check that the JsDoc of the
// nodes are equivalent.
//
// Deprecated: Use IsEquivalentToWithModes instead.
//
// Java: Node.isEquivalentTo(Node, boolean, boolean, boolean)
func (n Node) IsEquivalentToWithFlags(node Node, compareType bool, recurse bool, jsDoc bool) bool {
	return n.IsEquivalentToWithFlagsAndSideEffect(node, compareType, recurse, jsDoc, false)
}

// IsEquivalentToWithFlagsAndSideEffect returns whether this node is
// equivalent semantically to the provided node.
//
// sideEffect: whether to check that the side-effect flags of the nodes are
// equivalent.
//
// Deprecated: Use IsEquivalentToWithModes instead.
//
// Java: Node.isEquivalentTo(Node, boolean, boolean, boolean, boolean)
func (n Node) IsEquivalentToWithFlagsAndSideEffect(node Node, compareType bool, recurse bool, jsDoc bool, sideEffect bool) bool {
	recursion := Node_RecursionMode_SHALLOW
	if recurse {
		recursion = Node_RecursionMode_DEEP
	}
	typeConfig := Node_TypeComparison_IGNORE
	if compareType {
		typeConfig = Node_TypeComparison_COMPARE
	}
	jsDocConfig := Node_JsDocComparison_IGNORE
	if jsDoc {
		jsDocConfig = Node_JsDocComparison_COMPARE
	}
	sideEffectConfig := Node_SideEffectComparison_IGNORE
	if sideEffect {
		sideEffectConfig = Node_SideEffectComparison_COMPARE
	}
	return n.IsEquivalentToWithModes(node, recursion, typeConfig, jsDocConfig, sideEffectConfig)
}

// IsEquivalentToWithModes returns whether this node is equivalent
// semantically to the provided node.
//
// Go-only: language-adapter — dispatches to the overrides of the folded
// subclasses (NumberNode, BigIntNode, StringNode,
// TemplateLiteralSubstringNode), each of which calls the base
// implementation first.
//
// Java: Node.isEquivalentTo(Node, RecursionMode, TypeComparison, JsDocComparison, SideEffectComparison)
func (n Node) IsEquivalentToWithModes(
	node Node,
	recursion Node_RecursionMode,
	typeConfig Node_TypeComparison,
	jsDocConfig Node_JsDocComparison,
	sideEffectConfig Node_SideEffectComparison,
) bool {
	switch n.nodeKind() {
	case node_kindNumber:
		// Java: NumberNode.isEquivalentTo
		return n.isEquivalentToBase(node, recursion, typeConfig, jsDocConfig, sideEffectConfig) &&
			(n.node_number() == node.node_number()) // -0.0 and NaN are forbidden.
	case node_kindBigInt:
		// Java: BigIntNode.isEquivalentTo
		return n.isEquivalentToBase(node, recursion, typeConfig, jsDocConfig, sideEffectConfig) &&
			jutil.BigIntegerEquals(n.GetBigInt(), node.GetBigInt())
	case node_kindString:
		// Java: StringNode.isEquivalentTo
		return n.isEquivalentToBase(node, recursion, typeConfig, jsDocConfig, sideEffectConfig) &&
			uint32(n.rec().payload) == uint32(node.rec().payload)
	case node_kindTemplateLiteralSubstring:
		// Java: TemplateLiteralSubstringNode.isEquivalentTo
		if !n.isEquivalentToBase(node, recursion, typeConfig, jsDocConfig, sideEffectConfig) {
			return false
		}
		return n.rec().payload == node.rec().payload && n.node_cookedID() == node.node_cookedID()
	}
	return n.isEquivalentToBase(node, recursion, typeConfig, jsDocConfig, sideEffectConfig)
}

// isEquivalentToBase is Node's own isEquivalentTo (the super call of the
// folded subclasses).
func (n Node) isEquivalentToBase(
	node Node,
	recursion Node_RecursionMode,
	typeConfig Node_TypeComparison,
	jsDocConfig Node_JsDocComparison,
	sideEffectConfig Node_SideEffectComparison,
) bool {
	if n.rec().token != node.rec().token ||
		n.GetChildCount() != node.GetChildCount() ||
		n.nodeKind() != node.nodeKind() {
		return false
	}

	if typeConfig == Node_TypeComparison_COMPARE &&
		!jutil.ObjectsEquals(n.typeOrColorValue(), node.typeOrColorValue()) {
		return false
	}

	if jsDocConfig == Node_JsDocComparison_COMPARE &&
		!JSDocInfo_AreEquivalent(n.GetJSDocInfo(), node.GetJSDocInfo()) {
		return false
	}

	thisDte := n.GetDeclaredTypeExpression()
	thatDte := node.GetDeclaredTypeExpression()
	if thisDte == thatDte {
		// Do nothing
	} else if thisDte.IsNil() || thatDte.IsNil() {
		return false
	} else if !thisDte.IsEquivalentToWithModes(
		thatDte, recursion, typeConfig, jsDocConfig, sideEffectConfig) {
		return false
	}

	if (n.props() & node_propEqualityMask) != (node.props() & node_propEqualityMask) {
		return false
	}
	if ((n.rec().kindFlags | node.rec().kindFlags) & kindFlagHasExtra) != 0 {
		if n.GetIntProp(Node_Prop_EXPORT_DEFAULT) != node.GetIntProp(Node_Prop_EXPORT_DEFAULT) ||
			n.GetIntProp(Node_Prop_EXPORT_ALL_FROM) != node.GetIntProp(Node_Prop_EXPORT_ALL_FROM) ||
			n.GetIntProp(Node_Prop_INCRDECR) != node.GetIntProp(Node_Prop_INCRDECR) ||
			n.GetIntProp(Node_Prop_QUOTED) != node.GetIntProp(Node_Prop_QUOTED) {
			return false
		}
	}

	if sideEffectConfig == Node_SideEffectComparison_COMPARE {
		if n.GetSideEffectFlags() != node.GetSideEffectFlags() {
			return false
		}

		if n.IsUnusedParameter() != node.IsUnusedParameter() {
			return false
		}
	}

	thisIsShadowHost := !n.GetClosureUnawareShadow().IsNil()
	nodeIsShadowHost := !node.GetClosureUnawareShadow().IsNil()
	if thisIsShadowHost != nodeIsShadowHost {
		return false
	}

	if recursion.RecurseShadowAst() && thisIsShadowHost && nodeIsShadowHost {
		thisShadow := n.GetClosureUnawareShadow()
		nodeShadow := node.GetClosureUnawareShadow()
		if !thisShadow.IsEquivalentToWithModes(
			nodeShadow, recursion, typeConfig, jsDocConfig, sideEffectConfig) {
			return false
		}
	}

	if recursion.RecurseChildren() {
		for c, c2 := n.firstNode(), node.firstNode(); !c.IsNil(); c, c2 = c.nextNode(), c2.nextNode() {
			if !c.IsEquivalentToWithModes(c2, recursion, typeConfig, jsDocConfig, sideEffectConfig) {
				return false
			}
		}
	}

	return true
}

const node_propEqualityMask uint64 = propFlag_ARROW_FN |
	propFlag_ASYNC_FN |
	propFlag_GENERATOR_FN |
	propFlag_START_OF_OPT_CHAIN |
	propFlag_STATIC_MEMBER |
	propFlag_YIELD_ALL |
	propFlag_EXPORT_DEFAULT |
	propFlag_EXPORT_ALL_FROM |
	propFlag_INCRDECR |
	propFlag_QUOTED |
	propFlag_FREE_CALL |
	propFlag_DIRECT_EVAL |
	propFlag_COMPUTED_PROP_METHOD |
	propFlag_COMPUTED_PROP_GETTER |
	propFlag_COMPUTED_PROP_SETTER

// GetQualifiedName takes a set of GETPROP nodes and produces a string that
// is each property separated by dots. If the node ultimately under the left
// sub-tree is not a simple name, this is not a valid qualified name.
//
// Returns null if this is not a qualified name, or a dot-separated string of
// the name and properties.
//
// Java: Node.getQualifiedName()
//
// Go-only: performance — the body is in getQualifiedNameString, so that this
// wrapper inlines: a caller that only reads the result keeps the string
// header on its stack instead of the heap (these headers were 6% of
// checks-only's allocations).
func (n Node) GetQualifiedName() *string {
	s, ok := n.getQualifiedNameString()
	if !ok {
		return nil
	}
	return &s
}

// GetQualifiedNameString is the body of GetQualifiedName; ok is false
// where Java returns null.
func (n Node) GetQualifiedNameString() (string, bool) {
	if n.IsNil() {
		return "", false
	}
	return n.getQualifiedNameString()
}

// getQualifiedNameString is the body of getQualifiedName; ok is false
// where Java returns null.
func (n Node) getQualifiedNameString() (string, bool) {
	switch n.rec().token {
	case Token_NAME:
		name := n.GetString()
		if name == "" {
			return "", false
		}
		return name, true
	case Token_GETPROP:
		first := n.firstNode()
		if !first.IsGetProp() {
			left, ok := first.getQualifiedNameString()
			if !ok {
				return "", false
			}
			return left + "." + n.GetString(), true
		}
		var partsBuf [8]string
		parts := partsBuf[:0]
		totalLen := 0
		cur := n
		for cur.IsGetProp() {
			p := cur.GetString()
			parts = append(parts, p)
			totalLen += 1 + len(p)
			cur = cur.firstNode()
		}
		left, ok := cur.getQualifiedNameString()
		if !ok {
			return "", false
		}
		var builder strings.Builder
		builder.Grow(len(left) + totalLen)
		builder.WriteString(left)
		for i := len(parts) - 1; i >= 0; i-- {
			builder.WriteByte('.')
			builder.WriteString(parts[i])
		}
		return builder.String(), true
	case Token_THIS:
		return "this", true
	case Token_SUPER:
		return "super", true
	default:
		return "", false
	}
}

// Java: Node.getQualifiedNameObject()
func (n Node) GetQualifiedNameObject() QualifiedName {
	if n.IsQualifiedName() {
		return NewQualifiedName_NodeQname(n)
	}
	return nil
}

// GetOriginalQualifiedName takes a set of GETPROP nodes and produces a
// string that is each property separated by dots. If the node ultimately
// under the left sub-tree is not a simple name, this is not a valid
// qualified name. This method returns the original name of each segment
// rather than the renamed version.
//
// Deprecated: "original name" is poorly defined. See #getOriginalName
//
// Java: Node.getOriginalQualifiedName()
func (n Node) GetOriginalQualifiedName() *string {
	if n.rec().token == Token_NAME {
		var name string
		if n.GetOriginalName() == nil {
			name = n.GetString()
		} else {
			name = *n.GetOriginalName()
		}
		if name == "" {
			return nil
		}
		return &name
	} else if n.rec().token == Token_GETPROP {
		left := n.GetFirstChild().GetOriginalQualifiedName()
		if left == nil {
			return nil
		}

		var right string
		if n.GetOriginalName() == nil {
			right = n.GetString()
		} else {
			right = *n.GetOriginalName()
		}

		s := *left + "." + right
		return &s
	} else if n.rec().token == Token_THIS {
		s := "this"
		return &s
	} else if n.rec().token == Token_SUPER {
		s := "super"
		return &s
	} else {
		return nil
	}
}

// IsQualifiedName returns whether a node corresponds to a simple or a
// qualified name, such as `x` or `a.b.c` or `this.a`.
//
// Java: Node.isQualifiedName()
func (n Node) IsQualifiedName() bool {
	for {
		switch n.rec().token {
		case Token_NAME:
			return n.GetStringID() > 1
		case Token_THIS, Token_SUPER:
			return true
		case Token_GETPROP:
			n = n.firstNode()
		default:
			return false
		}
	}
}

// MatchesName returns whether a node matches a simple name, such as `x`,
// returns false if this is not a NAME node.
//
// Java: Node.matchesName(String)
func (n Node) MatchesName(name string) bool {
	if n.rec().token != Token_NAME {
		return false
	}
	internalString := n.GetString()
	return internalString != "" && name == internalString
}

// MatchesNameNode checks that if two NAME node match, returns false if
// either node is not a NAME node. As a empty string is not considered a
// valid Name (it is an AST placeholder), empty strings are never considered
// to be matches.
//
// Java: Node.matchesName(Node)
func (n Node) MatchesNameNode(other Node) bool {
	if n.rec().token != Token_NAME || other.rec().token != Token_NAME {
		return false
	}

	// Compare interned string IDs directly (id > 1 means non-empty!).
	id := rhinoStringID(n.rec().payload)
	return id > 1 && id == rhinoStringID(other.rec().payload)
}

// MatchesQualifiedName returns whether a node matches a simple or a
// qualified name, such as `x` or `a.b.c` or `this.a`.
//
// Java: Node.matchesQualifiedName(String)
func (n Node) MatchesQualifiedName(qname string) bool {
	for {
		switch n.rec().token {
		case Token_NAME, Token_IMPORT_STAR:
			name := n.GetString()
			return name != "" && name == qname && strings.IndexByte(name, '.') < 0
		case Token_THIS:
			return qname == "this"
		case Token_SUPER:
			return qname == "super"
		case Token_GETPROP:
			prop := n.GetString()
			dot := len(qname) - len(prop) - 1
			if dot < 1 || qname[dot] != '.' || qname[dot+1:] != prop || strings.IndexByte(prop, '.') >= 0 {
				return false
			}
			qname = qname[:dot]
			n = n.firstNode()
		default:
			return false
		}
	}
}

// MatchesQualifiedNameWithEnd returns whether a node matches a simple or a
// qualified name, such as `x` or `a.b.c` or `this.a`.
//
// Java: Node.matchesQualifiedName(String, int)
func (n Node) MatchesQualifiedNameWithEnd(qname string, endIndex int) bool {
	start := jutil.StringLastIndexOfCharFrom(qname, '.', endIndex-1) + 1

	switch n.GetToken() {
	case Token_NAME, Token_IMPORT_STAR:
		name := n.GetString()
		return start == 0 && name != "" && jutil.StringLength(name) == endIndex && jutil.StringStartsWith(qname, name)
	case Token_THIS:
		return start == 0 && endIndex == 4 && jutil.StringStartsWith(qname, "this")
	case Token_SUPER:
		return start == 0 && endIndex == 5 && jutil.StringStartsWith(qname, "super")
	case Token_GETPROP:
		prop := n.GetString()
		return start > 1 &&
			jutil.StringLength(prop) == endIndex-start &&
			jutil.StringRegionMatches(prop, 0, qname, start, endIndex-start) &&
			n.GetFirstChild().MatchesQualifiedNameWithEnd(qname, start-1)

	case Token_MEMBER_FUNCTION_DEF:
		// These are explicitly *not* qualified name components.
		return false
	default:
		return false
	}
}

// MatchesQualifiedNameNode returns whether a node matches a simple or a
// qualified name, such as `x` or `a.b.c` or `this.a`.
//
// Java: Node.matchesQualifiedName(Node)
func (n Node) MatchesQualifiedNameNode(other Node) bool {
	if other.rec().token != n.rec().token {
		return false
	}

	switch n.rec().token {
	case Token_NAME:
		return n.MatchesNameNode(other)
	case Token_THIS, Token_SUPER:
		return true
	case Token_GETPROP:
		return uint32(n.rec().payload) == uint32(other.rec().payload) &&
			n.GetFirstChild().MatchesQualifiedNameNode(other.GetFirstChild())

	case Token_MEMBER_FUNCTION_DEF:
		// These are explicitly *not* qualified name components.
		return false
	default:
		return false
	}
}

// IsUnscopedQualifiedName returns whether a node corresponds to a simple or
// a qualified name without a "this" reference, such as `a.b.c`, but not
// `this.a`.
//
// Java: Node.isUnscopedQualifiedName()
func (n Node) IsUnscopedQualifiedName() bool {
	for {
		switch n.rec().token {
		case Token_NAME:
			return n.GetStringID() > 1
		case Token_GETPROP:
			n = n.firstNode()
		default:
			return false
		}
	}
}

// Java: Node.isValidAssignmentTarget()
func (n Node) IsValidAssignmentTarget() bool {
	switch n.GetToken() {
	case Token_NAME, Token_GETPROP, Token_GETELEM, Token_ARRAY_PATTERN, Token_OBJECT_PATTERN:
		return true
	default:
		return false
	}
}

// Clone always throws: use CloneNode.
//
// Java: Node.clone()
func (n Node) Clone() any {
	panic(jutil.NewUnsupportedOperationExceptionMsg("Did you mean cloneNode?"))
}

// CloneNode returns a detached clone of the Node, specifically excluding its
// children.
//
// Java: Node.cloneNode()
func (n Node) CloneNode() Node {
	return n.CloneNodeWithTypeExprs(false)
}

// CloneNodeWithTypeExprs returns a detached clone of the Node, specifically
// excluding its children.
//
// Go-only: language-adapter — dispatches to the cloneNode overrides of the
// folded subclasses.
//
// Java: Node.cloneNode(boolean)
func (n Node) CloneNodeWithTypeExprs(cloneTypeExprs bool) Node {
	clone := allocNode()
	clone.rec().token = n.rec().token
	clone.rec().kindFlags = n.rec().kindFlags & kindMask
	clone.rec().payload = n.rec().payload
	Node_CopyBaseNodeFields(n, clone, cloneTypeExprs)
	return clone
}

// Java: Node.copyBaseNodeFields(Node, Node, boolean)
func Node_CopyBaseNodeFields(source Node, dest Node, cloneTypeExprs bool) {
	dest.rec().linenoCharno = source.rec().linenoCharno
	dest.setTC(dest.store().translateTC(source.store(), source.tc(), source.rec().kindFlags&kindFlagTypeIsColor != 0))
	dest.rec().kindFlags = (dest.rec().kindFlags &^ (kindFlagTypeIsColor | kindFlagHasExtBox | kindFlagHasExtra)) |
		(source.rec().kindFlags & (kindFlagTypeIsColor | kindFlagHasExtBox))
	if srcBox := source.extBoxOrNil(); srcBox != nil {
		srcBox.propListHead.markShared()
		boxCopy := *srcBox
		boxCopy.typeBeforeCast = dest.store().translateTypeHandle(source.store(), boxCopy.typeBeforeCast)
		dest.setAux(nodeAux{ext: &boxCopy})
	} else {
		dest.setAux(source.auxValue())
	}
	if source.rec().kindFlags&kindFlagHasExtra != 0 && source.coldPropListHead() != nil {
		dest.rec().kindFlags |= kindFlagHasExtra
	}
	dest.copyLengthFrom(source)
	dest.node_setOriginalNameID(source.node_originalNameID())
	dest.setPropsAndSourceFile(source)

	// TODO(johnlenz): Remove this once JSTypeExpression are immutable
	if cloneTypeExprs {
		info := source.GetJSDocInfo()
		if info != nil {
			dest.SetJSDocInfo(info.CloneWithTypeNodes(true))
		}
	}
}

// CloneTree returns a detached clone of the Node and all its children.
//
// Java: Node.cloneTree()
func (n Node) CloneTree() Node {
	return n.CloneTreeWithTypeExprs(false)
}

// Java: Node.cloneTree(boolean)
func (n Node) CloneTreeWithTypeExprs(cloneTypeExprs bool) Node {
	result := n.CloneNodeWithTypeExprs(cloneTypeExprs)
	if n.rec().first != 0 {
		var firstChild Node
		var lastChild Node
		for n2 := n.firstNode(); !n2.IsNil(); n2 = n2.nextNode() {
			n2clone := n2.CloneTreeWithTypeExprs(cloneTypeExprs)
			n2clone.rec().parent = n2clone.ref(result)
			if firstChild.IsNil() {
				firstChild = n2clone
				lastChild = firstChild
			} else {
				lastChild.rec().next = lastChild.ref(n2clone)
				n2clone.rec().previous = n2clone.ref(lastChild)
				lastChild = n2clone
			}
		}
		firstChild.rec().previous = firstChild.ref(lastChild)
		lastChild.rec().next = 0
		result.rec().first = result.ref(firstChild)
	}

	if shadow := n.GetClosureUnawareShadow(); !shadow.IsNil() {
		result.SetClosureUnawareShadow(shadow.CloneTreeWithTypeExprs(cloneTypeExprs))
	}

	return result
}

// Srcref copies the source info from `other` onto `this`.
//
// Java: Node.srcref(Node)
func (n Node) Srcref(other Node) Node {
	n.SetStaticSourceFileFrom(other)
	n.node_setOriginalNameID(other.node_originalNameID())
	n.rec().linenoCharno = other.rec().linenoCharno
	n.copyLengthFrom(other)
	return n
}

// SrcrefTree copies, for all Nodes in the subtree of `this`, the source info
// from `other`.
//
// Java: Node.srcrefTree(Node)
func (n Node) SrcrefTree(other Node) Node {
	n.Srcref(other)
	for child := n.firstNode(); !child.IsNil(); child = child.nextNode() {
		child.SrcrefTree(other)
	}
	return n
}

// SrcrefIfMissing copies the source info from `other` iff source info is not
// set on `this`.
//
// Java: Node.srcrefIfMissing(Node)
func (n Node) SrcrefIfMissing(other Node) Node {
	if n.rec().sourceFile == 0 {
		n.SetStaticSourceFileFrom(other)
		n.rec().linenoCharno = other.rec().linenoCharno
		n.copyLengthFrom(other)
	}

	// TODO(lharker): should this be inside the above if condition?
	// If the node already has a source file, it seems strange to
	// go ahead and set the original name anyway.
	if n.node_originalNameID() == 0 {
		n.node_setOriginalNameID(other.node_originalNameID())
	}

	return n
}

// SrcrefTreeIfMissing copies, for all Nodes in the subtree of `this`, iff
// source info is not set, the source info from `other`.
//
// Java: Node.srcrefTreeIfMissing(Node)
func (n Node) SrcrefTreeIfMissing(other Node) Node {
	n.SrcrefIfMissing(other)
	for child := n.firstNode(); !child.IsNil(); child = child.nextNode() {
		child.SrcrefTreeIfMissing(other)
	}
	return n
}
