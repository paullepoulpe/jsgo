// Port of com/google/javascript/rhino/Node.java @ closure-compiler v20261028

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

	"github.com/paullepoulpe/jsgo/internal/base"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/strpool"
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

func (s *NodeStore) internBigInt(v *big.Int) uint32 {
	id := s.bigInts.add(&s.mu, v)
	if s.static {
		return staticBit | id
	}
	return id
}

// Java: Node.getBigInt()
func (n Node) GetBigInt() *big.Int {
	n.node_castKind(node_kindBigInt)
	id := uint32(n.rec().payload)
	if id == 0 {
		return nil
	}
	if id&staticBit != 0 {
		return *staticStore.bigInts.at(id &^ staticBit)
	}
	return *n.store().bigInts.at(id)
}

// Java: Node.setBigInt(BigInteger)
func (n Node) SetBigInt(number *big.Int) {
	jutil.CheckNotNull(number)
	jutil.CheckStateMsg1(number.Sign() >= 0, "%s", number)
	n.node_castKind(node_kindBigInt)
	n.rec().payload = n.store().internBigInt(number)
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
	n.setStringPayload(strpool.Intern(str))
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

// CookedString returns the cooked string of a TemplateLiteralSubstringNode,
// or ("", false) if unset.
//
// Java: Node.getCookedString()
func (n Node) CookedString() (string, bool) {
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
			sb.WriteString(strpool.String(origID))
			sb.WriteByte(']')
		}
	}

	if printType && n.tc() != 0 {
		typeString := jutil.ToJavaString(n.typeOrColorValue())
		sb.WriteString(" : ")
		sb.WriteString(typeString)
	}
}

func writeJsonPair(sb *strings.Builder, name, value string) {
	sb.WriteByte('"')
	sb.WriteString(name)
	sb.WriteString(`":"`)
	sb.WriteString(value)
	sb.WriteByte('"')
}

func writeJsonPairRawValue(sb *strings.Builder, name, value string) {
	sb.WriteByte('"')
	sb.WriteString(name)
	sb.WriteString(`":`)
	sb.WriteString(value)
}

// Java: Node.createJsonPair(String, String)
func Node_CreateJsonPair(name string, value string) string {
	return `"` + name + `":"` + value + `"`
}

// Java: Node.createJsonPairRawValue(String, String)
func Node_CreateJsonPairRawValue(name string, value string) string {
	return `"` + name + `":` + value
}

// node_appender is Java's Appendable over an io.StringWriter: write errors
// become the IOException Appendable.append throws.
type node_appender struct{ w io.StringWriter }

func (a node_appender) append(s string) {
	if _, err := a.w.WriteString(s); err != nil {
		panic(jutil.NewIOException(err.Error()))
	}
}

func (n Node) toJsonInto(sb *strings.Builder) {
	sb.WriteByte('{')
	writeJsonPair(sb, "token", n.rec().token.String())
	if n.nodeKind() == node_kindString {
		sb.WriteByte(',')
		writeJsonPair(sb, "string", n.GetString())
	} else if n.rec().token == Token_FUNCTION {
		sb.WriteByte(',')
		// In the case of JsDoc trees, the first child is often not a string
		// which causes exceptions to be thrown when calling toString or
		// toStringTree.
		if n.rec().first == 0 || n.firstNode().rec().token != Token_NAME {
			writeJsonPair(sb, "functionName", "<invalid>")
		} else {
			writeJsonPair(sb, "functionName", n.firstNode().GetString())
		}
	} else if n.rec().token == Token_NUMBER {
		sb.WriteByte(',')
		writeJsonPair(sb, "number", jutil.DoubleToString(n.GetDouble()))
	}
	lineno := n.GetLineno()
	if lineno != -1 {
		sb.WriteString(`,"sourceLocation":{`)
		writeJsonPairRawValue(sb, "line", strconv.Itoa(lineno))
		sb.WriteByte(',')
		writeJsonPairRawValue(sb, "col", strconv.Itoa(n.GetCharno()))
		if length := n.GetLength(); length != 0 {
			sb.WriteByte(',')
			writeJsonPairRawValue(sb, "length", strconv.Itoa(length))
		}
		sb.WriteByte('}')
	}

	if origID := n.node_originalNameID(); origID != 0 {
		sb.WriteByte(',')
		writeJsonPair(sb, "original_name", strpool.String(origID))
	}

	keys := n.GetSortedPropTypes()
	if len(keys) != 0 {
		sb.WriteString(`,"props":{`)
		for i := 0; i < len(keys); i++ {
			type_ := Node_propValues[keys[i]]
			x := n.LookupProperty(type_)
			escaped := strings.ReplaceAll(strings.ReplaceAll(x.String(), "\n", `\n`), `"`, `\"`)
			writeJsonPair(sb, jutil.StringToLowerCaseASCII(type_.String()), escaped)
			if i+1 < len(keys) {
				sb.WriteByte(',')
			}
		}
		sb.WriteByte('}')
	}

	if n.tc() != 0 {
		typeString := jutil.ToJavaString(n.typeOrColorValue())
		sb.WriteByte(',')
		writeJsonPair(sb, "typeString", strings.ReplaceAll(typeString, `"`, `\"`))
	}

	if n.rec().first != 0 {
		sb.WriteString(`,"children":[`)
		for child := n.firstNode(); !child.IsNil(); child = child.nextNode() {
			child.toJsonInto(sb)
			if child.rec().next != 0 {
				sb.WriteByte(',')
			}
		}
		sb.WriteByte(']')
	}
	sb.WriteByte('}')
}

// ToJson writes the JSON form of this node and its subtree.
//
// Java: Node.toJson(Appendable)
func (n Node) ToJson(appendable io.StringWriter) {
	if sb, ok := appendable.(*strings.Builder); ok {
		n.toJsonInto(sb)
		return
	}
	var sb strings.Builder
	n.toJsonInto(&sb)
	node_appender{appendable}.append(sb.String())
}

// Java: Node.toStringTree()
func (n Node) ToStringTree() string {
	return n.ToStringTreeImpl()
}

// Java: Node.toStringTreeImpl()
func (n Node) ToStringTreeImpl() string {
	var s strings.Builder
	appendStringTreeInto(n, 0, &s)
	return s.String()
}

// Java: Node.appendStringTree(Appendable)
func (n Node) AppendStringTree(appendable io.StringWriter) {
	Node_ToStringTreeHelper(n, 0, appendable)
}

func appendStringTreeInto(n Node, level int, sb *strings.Builder) {
	for range level {
		sb.WriteString("    ")
	}
	n.ToStringInto(sb, true, true, true)
	sb.WriteByte('\n')
	for cursor := n.firstNode(); !cursor.IsNil(); cursor = cursor.nextNode() {
		appendStringTreeInto(cursor, level+1, sb)
	}
	if shadow := n.GetClosureUnawareShadow(); !shadow.IsNil() {
		appendStringTreeInto(shadow, level+1, sb)
	}
}

// Java: Node.toStringTreeHelper(Node, int, Appendable)
func Node_ToStringTreeHelper(n Node, level int, appendable io.StringWriter) {
	if sb, ok := appendable.(*strings.Builder); ok {
		appendStringTreeInto(n, level, sb)
		return
	}
	var sb strings.Builder
	appendStringTreeInto(n, level, &sb)
	node_appender{appendable}.append(sb.String())
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
func (n Node) node_cooked() (string, bool) {
	id := n.node_cookedID()
	if id == 0 {
		return "", false
	}
	return strpool.String(id), true
}

// ==========================================================================
// Source position management

// Java: Node.setStaticSourceFileFrom(Node)
func (n Node) SetStaticSourceFileFrom(other Node) {
	sfID := other.rec().sourceFile
	if sfID != 0 && sfID&staticBit == 0 && uint32(n)&staticBit != 0 {
		n.rec().sourceFile = staticStore.allocSourceFileID(lookupSourceFileSlot(sfID).item.objectValue)
		return
	}
	n.rec().sourceFile = sfID
}

// Java: Node.setStaticSourceFile(StaticSourceFile)
func (n Node) SetStaticSourceFile(file StaticSourceFile) Node {
	if jutil.IsNil(file) {
		n.RemoveProp(Node_Prop_SOURCE_FILE)
		return n
	}
	r := n.rec()
	r.sourceFile = n.store().allocSourceFileID(file)
	return n
}

// SetSourceFileForTesting sets the source file to a non-extern file of the
// given name.
//
// Java: Node.setSourceFileForTesting(String)
func (n Node) SetSourceFileForTesting(name string) {
	r := n.rec()
	r.sourceFile = n.store().allocSourceFileID(NewSimpleSourceFile(name, StaticSourceFile_SourceKind_STRONG))
}

// TODO(johnlenz): make this final
//
// SourceFileName returns the source file name of n, or ("", false) if n has no
// source file.
//
// Java: Node.getSourceFileName()
func (n Node) SourceFileName() (string, bool) {
	sf := n.GetStaticSourceFile()
	if sf == nil {
		return "", false
	}
	return sf.GetName(), true
}

// SourceFileID returns the internal source-file table ID of n, or 0 if unset.
func (n Node) SourceFileID() uint32 {
	if n.IsNil() {
		return 0
	}
	return n.rec().sourceFile
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

// GetOriginalNameID returns the 32-bit interned originalName ID of n, or 0 if unset.
func (n Node) GetOriginalNameID() uint32 {
	return uint32(n.node_originalNameID())
}

// OriginalName returns the original name of this node if the node has been
// renamed, or ("", false) if unset.
//
// Deprecated: "original name" is poorly defined.
//
// Java: Node.getOriginalName()
func (n Node) OriginalName() (string, bool) {
	id := n.node_originalNameID()
	if id == 0 {
		return "", false
	}
	return strpool.String(id), true
}

// GetOriginalNameString returns the node's originalName, or ("", false) if unset.
func (n Node) GetOriginalNameString() (string, bool) {
	return n.OriginalName()
}

// SetOriginalName interns s in RhinoStringPool and sets it as the node's
// original name.
//
// Java: Node.setOriginalName(String)
func (n Node) SetOriginalName(s string) {
	n.node_setOriginalNameID(strpool.Intern(s))
}

// ClearOriginalName clears the node's original name.
func (n Node) ClearOriginalName() {
	n.node_setOriginalNameID(0)
}

// SetOriginalNameString interns s in RhinoStringPool and sets it as the node's
// original name.
func (n Node) SetOriginalNameString(s string) {
	n.SetOriginalName(s)
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

// ExtractLineno extracts the line number from a packed linenoCharno value.
func ExtractLineno(linenoCharno int32) int {
	if linenoCharno == -1 {
		return -1
	}
	return int(uint32(linenoCharno) >> Node_CHARNO_BITS)
}

// ExtractCharno extracts the 0-based column number from a packed linenoCharno value.
func ExtractCharno(linenoCharno int32) int {
	if linenoCharno == -1 {
		return -1
	}
	return int(linenoCharno & Node_MAX_COLUMN_NUMBER)
}

// Java: Node.getLineno()
func (n Node) GetLineno() int {
	return ExtractLineno(n.rec().linenoCharno)
}

// GetCharno returns the 0-based column number.
//
// Java: Node.getCharno()
func (n Node) GetCharno() int {
	return ExtractCharno(n.rec().linenoCharno)
}

// Java: Node.getLocation()
func (n Node) GetLocation() string {
	sourceName := "null"
	if name, ok := n.SourceFileName(); ok {
		sourceName = name
	}
	return sourceName + ":" + strconv.Itoa(n.GetLineno()) + ":" + strconv.Itoa(n.GetCharno())
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
	n.rec().linenoCharno = int32(uint32(lineno)<<Node_CHARNO_BITS) | int32(charno)

	return n
}

// SetSourceInfoIfUnset sets lineno, charno, and length on n in a single record
// lookup when n's lineno is unset (-1).
func (n Node) SetSourceInfoIfUnset(lineno, charno, length int) {
	nr := n.rec()
	if nr.linenoCharno != -1 {
		return
	}
	if uint(lineno) < 0xFFFFF && uint(charno) <= Node_MAX_COLUMN_NUMBER && uint(length) < 0xFFFF && nr.length16 != 0xFFFF {
		nr.linenoCharno = int32(uint32(lineno)<<Node_CHARNO_BITS) | int32(charno)
		nr.length16 = uint16(length)
		return
	}
	n.setSourceInfoSlow(lineno, charno, length)
}

//go:noinline
func (n Node) setSourceInfoSlow(lineno, charno, length int) {
	n.SetLinenoCharno(lineno, charno)
	n.SetLength(length)
}
