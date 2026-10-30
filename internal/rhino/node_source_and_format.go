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

// GetOriginalNameString returns the node's originalName without allocating a *string.
func (n Node) GetOriginalNameString() (string, bool) {
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
