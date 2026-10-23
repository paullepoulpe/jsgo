// Derived from com/google/protobuf/TextFormatEscaper.java
// (protobuf-java 4.30.2): protoTextBytes is its escapeBytes. The notice of
// that file:
//
// Protocol Buffers - Google's data interchange format
// Copyright 2008 Google Inc.  All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd
//
// Here that LICENSE file is LICENSES/BSD-3-Clause-protobuf.txt.

// Runtime support for the reconstructed protobuf messages (proto_*.go).
//
// Go-only: language-adapter — generated protobuf classes are a leaf package
// in PORTING.md §2. The messages Closure uses are reconstructed from the
// descriptors in the pinned closure-compiler.jar by
// tools/oracle/java/ProtoGoGen.java; this file is the hand-written part
// they share: the wire format (CodedInputStream / CodedOutputStream
// semantics for the field types Closure's protos use), the text format of
// Message.toString() (TextFormat.printer() with escapeNonAscii, including
// unknown fields), and InvalidProtocolBufferException. Extensions,
// descriptors/reflection and text-format parsing are not ported.

package jscomp

import (
	"encoding/binary"
	"math"
	"math/bits"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// InvalidProtocolBufferException is
// com.google.protobuf.InvalidProtocolBufferException (an IOException).
type InvalidProtocolBufferException struct{ jutil.IOException }

func (*InvalidProtocolBufferException) invalidProtocolBufferException() {}

// IInvalidProtocolBufferException matches InvalidProtocolBufferException.
type IInvalidProtocolBufferException interface {
	jutil.IIOException
	invalidProtocolBufferException()
}

// NewInvalidProtocolBufferException is new InvalidProtocolBufferException(String).
func NewInvalidProtocolBufferException(description string) *InvalidProtocolBufferException {
	e := &InvalidProtocolBufferException{}
	jutil.InitThrowable(e, "com.google.javascript.jscomp.jarjar.com.google.protobuf.InvalidProtocolBufferException", &description, nil)
	return e
}

func protoTruncatedMessage() *InvalidProtocolBufferException {
	return NewInvalidProtocolBufferException("While parsing a protocol message, the input ended unexpectedly " +
		"in the middle of a field.  This could mean either that the " +
		"input has been truncated or that an embedded message " +
		"misreported its own length.")
}

func protoNegativeSize() *InvalidProtocolBufferException {
	return NewInvalidProtocolBufferException("CodedInputStream encountered an embedded string or message " +
		"which claimed to have negative size.")
}

func protoMalformedVarint() *InvalidProtocolBufferException {
	return NewInvalidProtocolBufferException("CodedInputStream encountered a malformed varint.")
}

func protoInvalidTag() *InvalidProtocolBufferException {
	return NewInvalidProtocolBufferException("Protocol message contained an invalid tag (zero).")
}

func protoInvalidEndTag() *InvalidProtocolBufferException {
	return NewInvalidProtocolBufferException("Protocol message end-group tag did not match expected tag.")
}

func protoInvalidWireType() *InvalidProtocolBufferException {
	return NewInvalidProtocolBufferException("Protocol message tag had invalid wire type.")
}

func protoInvalidUtf8() *InvalidProtocolBufferException {
	return NewInvalidProtocolBufferException("Protocol message had invalid UTF-8.")
}

// Wire types.
const (
	protoWireVarint     = 0
	protoWireFixed64    = 1
	protoWireBytes      = 2
	protoWireStartGroup = 3
	protoWireEndGroup   = 4
	protoWireFixed32    = 5
)

// ---------------------------------------------------------------------------
// Encoding

func protoAppendVarint(b []byte, v uint64) []byte {
	if v < 0x80 {
		return append(b, byte(v))
	}
	if v < 0x4000 {
		return append(b, byte(v)|0x80, byte(v>>7))
	}
	return protoAppendVarintSlow(b, v)
}

func protoAppendVarintSlow(b []byte, v uint64) []byte {
	for v >= 0x80 {
		b = append(b, byte(v)|0x80)
		v >>= 7
	}
	return append(b, byte(v))
}

func protoPutVarint(dst []byte, v uint64) int {
	i := 0
	for v >= 0x80 {
		dst[i] = byte(v) | 0x80
		v >>= 7
		i++
	}
	dst[i] = byte(v)
	return i + 1
}

func protoAppendTag(b []byte, number int, wireType int) []byte {
	return protoAppendVarint(b, uint64(uint32(number)<<3|uint32(wireType)))
}

// protoAppendInt32 writes an int32/enum value (sign-extended to 64 bits).
func protoAppendInt32(b []byte, v int32) []byte { return protoAppendVarint(b, uint64(int64(v))) }

func protoAppendUint32(b []byte, v int32) []byte { return protoAppendVarint(b, uint64(uint32(v))) }

func protoAppendSint32(b []byte, v int32) []byte {
	return protoAppendVarint(b, uint64(uint32((v<<1)^(v>>31))))
}

func protoAppendSint64(b []byte, v int64) []byte {
	return protoAppendVarint(b, uint64((v<<1)^(v>>63)))
}

func protoAppendBool(b []byte, v bool) []byte {
	if v {
		return append(b, 1)
	}
	return append(b, 0)
}

func protoAppendFixed32(b []byte, v uint32) []byte { return binary.LittleEndian.AppendUint32(b, v) }

func protoAppendFixed64(b []byte, v uint64) []byte { return binary.LittleEndian.AppendUint64(b, v) }

func protoAppendDouble(b []byte, v float64) []byte {
	return protoAppendFixed64(b, math.Float64bits(v))
}

func protoAppendBytes(b []byte, v []byte) []byte {
	b = protoAppendVarint(b, uint64(len(v)))
	return append(b, v...)
}

// protoAppendString writes a string as UTF-8 (unpaired surrogates become
// '?', as CodedOutputStream's fallback encoder does).
func protoAppendString(b []byte, v string) []byte {
	return protoAppendBytes(b, jutil.StringGetBytesUTF8(v))
}

// ---------------------------------------------------------------------------
// Decoding

// protoReader is a CodedInputStream over a byte slice. depth is the
// CodedInputStream recursionDepth of the message being read; limit is its
// recursion limit (0 means the default of 100).
type protoReader struct {
	b     []byte
	pos   int
	depth int
	limit int
}

const protoDefaultRecursionLimit = 100

func protoRecursionLimitExceeded() *InvalidProtocolBufferException {
	return NewInvalidProtocolBufferException("Protocol message had too many levels of nesting.  May be malicious.  " +
		"Use setRecursionLimit() to increase the recursion depth limit.")
}

// checkRecursionLimit is CodedInputStream.checkRecursionLimit().
func (r *protoReader) checkRecursionLimit() {
	limit := r.limit
	if limit == 0 {
		limit = protoDefaultRecursionLimit
	}
	if r.depth >= limit {
		panic(protoRecursionLimitExceeded())
	}
}

// sub returns a reader for an embedded message (CodedInputStream.readMessage
// with recursionDepth incremented).
func (r *protoReader) sub(data []byte) *protoReader {
	r.checkRecursionLimit()
	return &protoReader{b: data, depth: r.depth + 1, limit: r.limit}
}

func (r *protoReader) atEnd() bool { return r.pos >= len(r.b) }

func (r *protoReader) readRawVarint64() uint64 {
	var x uint64
	for shift := 0; shift < 64; shift += 7 {
		if r.pos >= len(r.b) {
			panic(protoTruncatedMessage())
		}
		c := r.b[r.pos]
		r.pos++
		x |= uint64(c&0x7F) << shift
		if c < 0x80 {
			return x
		}
	}
	panic(protoMalformedVarint())
}

// readTag is CodedInputStream.readTag(): 0 at the end of input.
func (r *protoReader) readTag() int {
	if r.atEnd() {
		return 0
	}
	tag := int32(r.readRawVarint64())
	if uint32(tag)>>3 == 0 {
		panic(protoInvalidTag())
	}
	return int(tag)
}

func (r *protoReader) readInt32() int32  { return int32(r.readRawVarint64()) }
func (r *protoReader) readInt64() int64  { return int64(r.readRawVarint64()) }
func (r *protoReader) readUint32() int32 { return int32(r.readRawVarint64()) }
func (r *protoReader) readBool() bool    { return r.readRawVarint64() != 0 }

func (r *protoReader) readSint32() int32 {
	n := uint32(r.readRawVarint64())
	return int32(n>>1) ^ -int32(n&1)
}

func (r *protoReader) readSint64() int64 {
	n := r.readRawVarint64()
	return int64(n>>1) ^ -int64(n&1)
}

func (r *protoReader) readFixed32() uint32 {
	if len(r.b)-r.pos < 4 {
		panic(protoTruncatedMessage())
	}
	v := binary.LittleEndian.Uint32(r.b[r.pos:])
	r.pos += 4
	return v
}

func (r *protoReader) readFixed64() uint64 {
	if len(r.b)-r.pos < 8 {
		panic(protoTruncatedMessage())
	}
	v := binary.LittleEndian.Uint64(r.b[r.pos:])
	r.pos += 8
	return v
}

func (r *protoReader) readDouble() float64 { return math.Float64frombits(r.readFixed64()) }

func (r *protoReader) readBytes() []byte {
	return append([]byte(nil), r.readBytesView()...)
}

// readBytesView is readBytes without the copy: the result aliases the
// reader's buffer, which is never written, so it is only for callers that
// do not keep it (decoding a string or an embedded message, skipping).
//
// Go-only: performance — copying every embedded message before decoding it
// copied each byte once per level of nesting.
func (r *protoReader) readBytesView() []byte {
	size := int32(r.readRawVarint64())
	if size < 0 {
		panic(protoNegativeSize())
	}
	if int(size) > len(r.b)-r.pos {
		panic(protoTruncatedMessage())
	}
	v := r.b[r.pos : r.pos+int(size) : r.pos+int(size)]
	r.pos += int(size)
	return v
}

// readStringRequireUtf8 is CodedInputStream.readStringRequireUtf8 (proto3).
func (r *protoReader) readStringRequireUtf8() string {
	b := r.readBytesView()
	if !utf8Valid(b) {
		panic(protoInvalidUtf8())
	}
	return string(b)
}

// readStringProto2 decodes a proto2 string like ByteString.toStringUtf8
// (malformed input is replaced).
func (r *protoReader) readStringProto2() string {
	return jutil.NewStringFromUTF8Bytes(r.readBytesView())
}

// readMessageBytes returns the bytes of a length-delimited embedded message.
// It aliases the reader's buffer (see readBytesView): every caller only
// decodes it with a sub-reader.
func (r *protoReader) readMessageBytes() []byte { return r.readBytesView() }

// skipField skips the field with the given tag (CodedInputStream.skipField).
func (r *protoReader) skipField(tag int) {
	switch tag & 7 {
	case protoWireVarint:
		r.readRawVarint64()
	case protoWireFixed64:
		r.readFixed64()
	case protoWireBytes:
		r.readBytesView()
	case protoWireStartGroup:
		r.skipGroup(tag >> 3)
	case protoWireEndGroup:
		panic(protoInvalidEndTag())
	case protoWireFixed32:
		r.readFixed32()
	default:
		panic(protoInvalidWireType())
	}
}

// skipGroup skips an unknown group; like Java's parse of the group into the
// UnknownFieldSet (CodedInputStream.readGroup), it counts towards the
// recursion limit.
func (r *protoReader) skipGroup(number int) {
	r.checkRecursionLimit()
	r.depth++
	defer func() { r.depth-- }()
	for {
		tag := r.readTag()
		if tag == 0 {
			panic(protoTruncatedMessage())
		}
		if tag&7 == protoWireEndGroup {
			if tag>>3 != number {
				panic(protoInvalidEndTag())
			}
			return
		}
		r.skipField(tag)
	}
}

// readPacked calls read for each element of a packed repeated field, or once
// for an unpacked element, depending on the wire type of tag.
func (r *protoReader) readRepeated(tag int, scalarWireType int, read func(*protoReader)) {
	if tag&7 == protoWireBytes && scalarWireType != protoWireBytes {
		data := r.readBytesView()
		origB, origPos := r.b, r.pos
		r.b, r.pos = data, 0
		for !r.atEnd() {
			read(r)
		}
		r.b, r.pos = origB, origPos
		return
	}
	read(r)
}

// ---------------------------------------------------------------------------
// Text format (Message.toString())

type protoTextPrinter struct {
	sb     strings.Builder
	indent int
}

func (p *protoTextPrinter) startLine() {
	for i := 0; i < p.indent; i++ {
		p.sb.WriteString("  ")
	}
}

// scalar prints "name: value\n".
func (p *protoTextPrinter) scalar(name, value string) {
	p.startLine()
	p.sb.WriteString(name)
	p.sb.WriteString(": ")
	p.sb.WriteString(value)
	p.sb.WriteByte('\n')
}

// message prints "name {\n" ... "}\n".
func (p *protoTextPrinter) message(name string, print func(*protoTextPrinter)) {
	p.messageWithSeparator(name, " {\n", print)
}

func (p *protoTextPrinter) messageWithSeparator(name, open string, print func(*protoTextPrinter)) {
	p.startLine()
	p.sb.WriteString(name)
	p.sb.WriteString(open)
	p.indent++
	print(p)
	p.indent--
	p.startLine()
	p.sb.WriteString("}\n")
}

func protoTextInt32(v int32) string    { return strconv.FormatInt(int64(v), 10) }
func protoTextInt64(v int64) string    { return strconv.FormatInt(v, 10) }
func protoTextUint32(v int32) string   { return strconv.FormatUint(uint64(uint32(v)), 10) }
func protoTextUint64(v int64) string   { return strconv.FormatUint(uint64(v), 10) }
func protoTextBool(v bool) string      { return strconv.FormatBool(v) }
func protoTextDouble(v float64) string { return jutil.DoubleToString(v) }

// protoTextBytes is "\"" + TextFormat.escapeBytes(b) + "\"".
// The escaping is protobuf's TextFormatEscaper.escapeBytes (BSD-3-Clause;
// see the notice at the top of this file).
func protoTextBytes(b []byte) string {
	var sb strings.Builder
	sb.WriteByte('"')
	for _, c := range b {
		switch c {
		case 0x07:
			sb.WriteString(`\a`)
		case '\b':
			sb.WriteString(`\b`)
		case '\f':
			sb.WriteString(`\f`)
		case '\n':
			sb.WriteString(`\n`)
		case '\r':
			sb.WriteString(`\r`)
		case '\t':
			sb.WriteString(`\t`)
		case 0x0b:
			sb.WriteString(`\v`)
		case '\\':
			sb.WriteString(`\\`)
		case '\'':
			sb.WriteString(`\'`)
		case '"':
			sb.WriteString(`\"`)
		default:
			if c >= 0x20 && c <= 0x7e {
				sb.WriteByte(c)
			} else {
				sb.WriteByte('\\')
				sb.WriteByte('0' + (c>>6)&3)
				sb.WriteByte('0' + (c>>3)&7)
				sb.WriteByte('0' + c&7)
			}
		}
	}
	sb.WriteByte('"')
	return sb.String()
}

// protoTextString is "\"" + TextFormat.escapeText(s) + "\"".
func protoTextString(s string) string { return protoTextBytes(jutil.StringGetBytesUTF8(s)) }

// protoUnknownField is UnknownFieldSet.Field.
type protoUnknownField struct {
	varint          []uint64
	fixed32         []uint32
	fixed64         []uint64
	lengthDelimited [][]byte
	group           []*protoUnknownFieldSet
}

// protoUnknownFieldSet is UnknownFieldSet: fields sorted by number.
type protoUnknownFieldSet struct {
	numbers []int
	fields  map[int]*protoUnknownField
}

func (s *protoUnknownFieldSet) field(number int) *protoUnknownField {
	if f, ok := s.fields[number]; ok {
		return f
	}
	f := &protoUnknownField{}
	s.fields[number] = f
	i := 0
	for i < len(s.numbers) && s.numbers[i] < number {
		i++
	}
	s.numbers = append(s.numbers, 0)
	copy(s.numbers[i+1:], s.numbers[i:])
	s.numbers[i] = number
	return f
}

// protoParseUnknownFieldSet is UnknownFieldSet.parseFrom; it panics with
// InvalidProtocolBufferException.
func protoParseUnknownFieldSet(data []byte) *protoUnknownFieldSet {
	r := &protoReader{b: data}
	s := r.readUnknownFieldSet(-1)
	return s
}

func (r *protoReader) readUnknownFieldSet(endGroup int) *protoUnknownFieldSet {
	s := &protoUnknownFieldSet{fields: map[int]*protoUnknownField{}}
	for {
		tag := r.readTag()
		if tag == 0 {
			if endGroup >= 0 {
				panic(protoTruncatedMessage())
			}
			return s
		}
		number := tag >> 3
		switch tag & 7 {
		case protoWireVarint:
			f := s.field(number)
			f.varint = append(f.varint, r.readRawVarint64())
		case protoWireFixed64:
			f := s.field(number)
			f.fixed64 = append(f.fixed64, r.readFixed64())
		case protoWireBytes:
			f := s.field(number)
			f.lengthDelimited = append(f.lengthDelimited, r.readBytes())
		case protoWireStartGroup:
			r.checkRecursionLimit()
			r.depth++
			g := r.readUnknownFieldSet(number)
			r.depth--
			f := s.field(number)
			f.group = append(f.group, g)
		case protoWireEndGroup:
			if number != endGroup {
				panic(protoInvalidEndTag())
			}
			return s
		case protoWireFixed32:
			f := s.field(number)
			f.fixed32 = append(f.fixed32, r.readFixed32())
		default:
			panic(protoInvalidWireType())
		}
	}
}

// writeTo is UnknownFieldSet.writeTo: fields in number order, each as
// varints, fixed32s, fixed64s, length-delimited values, then groups.
func (s *protoUnknownFieldSet) writeTo(b []byte) []byte {
	for _, number := range s.numbers {
		f := s.fields[number]
		for _, v := range f.varint {
			b = protoAppendTag(b, number, protoWireVarint)
			b = protoAppendVarint(b, v)
		}
		for _, v := range f.fixed32 {
			b = protoAppendTag(b, number, protoWireFixed32)
			b = protoAppendFixed32(b, v)
		}
		for _, v := range f.fixed64 {
			b = protoAppendTag(b, number, protoWireFixed64)
			b = protoAppendFixed64(b, v)
		}
		for _, v := range f.lengthDelimited {
			b = protoAppendTag(b, number, protoWireBytes)
			b = protoAppendBytes(b, v)
		}
		for _, g := range f.group {
			b = protoAppendTag(b, number, protoWireStartGroup)
			b = g.writeTo(b)
			b = protoAppendTag(b, number, protoWireEndGroup)
		}
	}
	return b
}

// protoCanonicalUnknownFields re-encodes raw unknown fields the way Java
// serializes its UnknownFieldSet.
func protoCanonicalUnknownFields(raw []byte) []byte {
	if len(raw) == 0 {
		return nil
	}
	return protoParseUnknownFieldSet(raw).writeTo(nil)
}

// unknownFields prints raw unknown fields like TextFormat.printUnknownFields.
func (p *protoTextPrinter) unknownFields(raw []byte) {
	if len(raw) == 0 {
		return
	}
	p.unknownFieldSet(protoParseUnknownFieldSet(raw))
}

func (p *protoTextPrinter) unknownFieldSet(s *protoUnknownFieldSet) {
	for _, number := range s.numbers {
		f := s.fields[number]
		n := strconv.Itoa(number)
		for _, v := range f.varint {
			p.scalar(n, strconv.FormatUint(v, 10))
		}
		for _, v := range f.fixed32 {
			h := strconv.FormatUint(uint64(v), 16)
			p.scalar(n, "0x"+strings.Repeat("0", 8-len(h))+h)
		}
		for _, v := range f.fixed64 {
			h := strconv.FormatUint(v, 16)
			p.scalar(n, "0x"+strings.Repeat("0", 16-len(h))+h)
		}
		for _, v := range f.lengthDelimited {
			var nested *protoUnknownFieldSet
			jutil.Try(func() {
				nested = protoParseUnknownFieldSet(v)
			}, func(e IInvalidProtocolBufferException) {
				nested = nil
			})
			if nested != nil {
				p.messageWithSeparator(n, ": {\n", func(p *protoTextPrinter) { p.unknownFieldSet(nested) })
			} else {
				p.scalar(n, protoTextBytes(v))
			}
		}
		for _, g := range f.group {
			p.message(n, func(p *protoTextPrinter) { p.unknownFieldSet(g) })
		}
	}
}

// ---------------------------------------------------------------------------
// Misc helpers for the generated code

func utf8Valid(b []byte) bool { return utf8.Valid(b) }

func protoBytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func protoHashBytes(h int32, b []byte) int32 {
	for _, c := range b {
		h = 31*h + int32(int8(c))
	}
	return h
}

func protoHashString(h int32, s string) int32 { return 31*h + jutil.StringHashCode(s) }

func protoHashBool(h int32, v bool) int32 {
	if v {
		return 31*h + 1231
	}
	return 31*h + 1237
}

func protoHashInt64(h int32, v int64) int32 { return 31*h + int32(v^int64(uint64(v)>>32)) }

func protoDoubleBits(v float64) int64 { return jutil.DoubleToLongBits(v) }

func protoVarintLen(v uint64) int {
	return int(9*uint32(bits.Len64(v))+64) >> 6
}

// protoCheckIndex is the bounds check of the primitive protobuf lists
// (IntArrayList.ensureIndexInRange).
func protoCheckIndex(index, size int) {
	if index < 0 || index >= size {
		panic(jutil.NewIndexOutOfBoundsException("Index:" + strconv.Itoa(index) + ", Size:" + strconv.Itoa(size)))
	}
}

// protoCheckIndexJDK is the bounds check of the object protobuf lists
// (ArrayList.get: Objects.checkIndex).
func protoCheckIndexJDK(index, size int) {
	if index < 0 || index >= size {
		panic(jutil.NewIndexOutOfBoundsException("Index " + strconv.Itoa(index) + " out of bounds for length " + strconv.Itoa(size)))
	}
}
