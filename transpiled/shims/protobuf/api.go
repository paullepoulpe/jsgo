package protobuf

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"sync"
	"unicode/utf8"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javanio"
)

// ---------------------------------------------------------------------------
// WireFormat & ExtensionRegistry
// ---------------------------------------------------------------------------

const (
	WireFormat_WIRETYPE_VARINT           int32 = 0
	WireFormat_WIRETYPE_FIXED64          int32 = 1
	WireFormat_WIRETYPE_LENGTH_DELIMITED int32 = 2
	WireFormat_WIRETYPE_START_GROUP      int32 = 3
	WireFormat_WIRETYPE_END_GROUP        int32 = 4
	WireFormat_WIRETYPE_FIXED32          int32 = 5
)

func WireFormat_GetTagWireType(tag int32) int32    { return tag & 7 }
func WireFormat_GetTagFieldNumber(tag int32) int32 { return int32(uint32(tag) >> 3) }

type ExtensionRegistry struct{}

var emptyExtensionRegistry = &ExtensionRegistry{}

func ExtensionRegistry_GetEmptyRegistry() *ExtensionRegistry {
	return emptyExtensionRegistry
}

// ---------------------------------------------------------------------------
// ByteString
// ---------------------------------------------------------------------------

type ByteString struct {
	data []byte
	str  string
}

type byteStringViewSlot struct {
	mu sync.Mutex
	s  string
	b  []byte
}

var (
	ByteString_EMPTY    = ByteString{}
	emptyByteSlice      = []byte{}
	byteStringViewCache [64]byteStringViewSlot
	byteString1Bytes    = func() [256]byte {
		var b [256]byte
		for i := range b {
			b[i] = byte(i)
		}
		return b
	}()
)

func ByteString_CopyFrom(b []byte) ByteString {
	n := len(b)
	if n == 0 {
		return ByteString_EMPTY
	}
	if n == 1 {
		c := b[0]
		return ByteString{data: byteString1Bytes[c : int(c)+1 : int(c)+1]}
	}
	if n <= 256 {
		dst := allocWireBytes(n)
		dst = dst[:n:n]
		copy(dst, b)
		return ByteString{data: dst}
	}
	cp := make([]byte, n)
	copy(cp, b)
	return ByteString{data: cp}
}

func ByteString_FromBytesUnsafeWrap(b []byte) ByteString {
	if len(b) == 0 {
		return ByteString_EMPTY
	}
	return ByteString{data: b}
}

func ByteString_CopyFromUtf8(s string) ByteString {
	n := len(s)
	if n == 0 {
		return ByteString_EMPTY
	}
	if n == 1 && s[0] < 0x80 {
		c := s[0]
		return ByteString{data: byteString1Bytes[c : int(c)+1 : int(c)+1]}
	}
	if javalang.IsASCII(s) || utf8.ValidString(s) {
		return ByteString{str: s}
	}
	return ByteString{data: javanio.StandardCharsets_UTF_8.Encode(s)}
}

func (bs ByteString) Size() int {
	if bs.data != nil {
		return len(bs.data)
	}
	return len(bs.str)
}

func (bs ByteString) IsEmpty() bool { return len(bs.data) == 0 && len(bs.str) == 0 }

func (bs ByteString) ByteAt(index int) byte {
	if bs.data != nil || bs.str == "" {
		if index < 0 || index >= len(bs.data) {
			panic(javalang.NewIndexOutOfBoundsExceptionForIndex(index, len(bs.data)))
		}
		return bs.data[index]
	}
	if index < 0 || index >= len(bs.str) {
		panic(javalang.NewIndexOutOfBoundsExceptionForIndex(index, len(bs.str)))
	}
	return bs.str[index]
}

func (bs ByteString) ToByteArray() []byte {
	if bs.data != nil || bs.str == "" {
		if len(bs.data) == 0 {
			return []byte{}
		}
		cp := make([]byte, len(bs.data))
		copy(cp, bs.data)
		return cp
	}
	cp := make([]byte, len(bs.str))
	copy(cp, bs.str)
	return cp
}

func (bs ByteString) BytesView() []byte {
	if bs.data != nil || bs.str == "" {
		return bs.data
	}
	s := bs.str
	idx := (uint(len(s))*31 ^ uint(s[0])<<3 ^ uint(s[len(s)-1])) & 63
	sh := &byteStringViewCache[idx]
	sh.mu.Lock()
	if sh.s == s && sh.b != nil {
		out := sh.b
		sh.mu.Unlock()
		return out
	}
	out := []byte(s)
	out = out[:len(out):len(out)]
	sh.s = s
	sh.b = out
	sh.mu.Unlock()
	return out
}

// bytesImmutable returns bs's immutable backing slice without cloning when []byte-backed,
// capped to its length so any accidental append allocates rather than mutating in place.
func (bs ByteString) bytesImmutable() []byte {
	if bs.data != nil || bs.str == "" {
		n := len(bs.data)
		if n == 0 {
			return emptyByteSlice[:0:0]
		}
		return bs.data[:n:n]
	}
	n := len(bs.str)
	if n <= 256 {
		dst := allocWireBytes(n)[:n:n]
		copy(dst, bs.str)
		return dst
	}
	cp := make([]byte, n)
	copy(cp, bs.str)
	return cp
}

func (bs ByteString) ToStringUtf8() string {
	if bs.data == nil && bs.str != "" {
		return bs.str
	}
	return javanio.StandardCharsets_UTF_8.Decode(bs.data)
}

func (bs ByteString) NewCodedInput() *CodedInputStream {
	if bs.data != nil || bs.str == "" {
		return CodedInputStream_NewInstanceBytes(bs.data)
	}
	return CodedInputStream_NewInstanceBytes(bs.bytesImmutable())
}

func byteStringEqual(a, b ByteString) bool {
	if a.data != nil || a.str == "" {
		if b.data != nil || b.str == "" {
			return bytes.Equal(a.data, b.data)
		}
		return string(a.data) == b.str
	}
	if b.data != nil || b.str == "" {
		return a.str == string(b.data)
	}
	return a.str == b.str
}

func (bs ByteString) Equals(other any) bool {
	switch o := other.(type) {
	case ByteString:
		return byteStringEqual(bs, o)
	case *ByteString:
		return o != nil && byteStringEqual(bs, *o)
	default:
		return false
	}
}

func (bs ByteString) HashCode() int32 {
	if bs.data == nil && bs.str != "" {
		h := int32(len(bs.str))
		for i := 0; i < len(bs.str); i++ {
			h = h*31 + int32(int8(bs.str[i]))
		}
		if h == 0 {
			return 1
		}
		return h
	}
	h := int32(len(bs.data))
	for _, b := range bs.data {
		h = h*31 + int32(int8(b))
	}
	if h == 0 {
		return 1
	}
	return h
}

func (bs ByteString) String() string {
	return fmt.Sprintf("<ByteString@%x size=%d>", bs.HashCode(), bs.Size())
}

// ByteString_Output models ByteString.Output in 48 bytes:
//   - <= 15 bytes stay inline with zero pool operations;
//   - 16..32768 bytes use a sync.Pool 32 KiB noscan scratch buffer (zero growth reallocs)
//     returned to the pool in ToByteString() while copying into the 2 MiB noscan wire slab;
//   - frozen outputs clamp cap(buf) == len(buf) so post-ToByteString writes copy-on-write safely.
const (
	byteStringOutputInlineCap = 15
	byteStringScratchSize     = 32 * 1024
)

type byteStringScratch [byteStringScratchSize]byte

var byteStringScratchPool = sync.Pool{
	New: func() any {
		return new(byteStringScratch)
	},
}

type ByteString_Output struct {
	buf     []byte
	scratch *byteStringScratch
	n       uint8
	inline  [byteStringOutputInlineCap]byte
}

//go:noinline
func panicByteStringOutputNegSize() {
	panic(javalang.NewIllegalArgumentException("Buffer size < 0"))
}

func ByteString_NewOutput(initialCapacity int) *ByteString_Output {
	if initialCapacity < 0 {
		panicByteStringOutputNegSize()
	}
	return &ByteString_Output{}
}

//go:noinline
func writeByteOutputSlow(o *ByteString_Output, b byte) {
	if o.buf == nil {
		cur := int(o.n)
		sc := byteStringScratchPool.Get().(*byteStringScratch)
		copy(sc[:cur], o.inline[:cur])
		sc[cur] = b
		o.scratch = sc
		o.buf = sc[: cur+1 : len(sc)]
		o.n = 0
		return
	}
	cur := len(o.buf)
	if o.scratch != nil {
		cp := make([]byte, cur+1, cur*2)
		copy(cp, o.buf)
		cp[cur] = b
		byteStringScratchPool.Put(o.scratch)
		o.scratch = nil
		o.buf = cp
		return
	}
	if cur+1 <= byteStringScratchSize {
		sc := byteStringScratchPool.Get().(*byteStringScratch)
		copy(sc[:cur], o.buf)
		sc[cur] = b
		o.scratch = sc
		o.buf = sc[: cur+1 : len(sc)]
		return
	}
	cp := make([]byte, cur+1, max(cur*2, cur+1))
	copy(cp, o.buf)
	cp[cur] = b
	o.buf = cp
}

//go:noinline
func writeSliceOutputSlow(o *ByteString_Output, p []byte) {
	if o.buf == nil {
		cur := int(o.n)
		need := cur + len(p)
		if need <= byteStringOutputInlineCap {
			copy(o.inline[cur:], p)
			o.n = uint8(need)
			return
		}
		if need <= byteStringScratchSize {
			sc := byteStringScratchPool.Get().(*byteStringScratch)
			copy(sc[:cur], o.inline[:cur])
			copy(sc[cur:need], p)
			o.scratch = sc
			o.buf = sc[:need:len(sc)]
			o.n = 0
			return
		}
		cp := make([]byte, need, max(need, 512))
		copy(cp, o.inline[:cur])
		copy(cp[cur:], p)
		o.buf = cp
		o.n = 0
		return
	}
	cur := len(o.buf)
	need := cur + len(p)
	if o.scratch != nil {
		cp := make([]byte, need, max(len(o.scratch)*2, need))
		copy(cp, o.buf)
		copy(cp[cur:], p)
		byteStringScratchPool.Put(o.scratch)
		o.scratch = nil
		o.buf = cp
		return
	}
	if need <= byteStringScratchSize {
		sc := byteStringScratchPool.Get().(*byteStringScratch)
		copy(sc[:cur], o.buf)
		copy(sc[cur:need], p)
		o.scratch = sc
		o.buf = sc[:need:len(sc)]
		return
	}
	cp := make([]byte, need, max(cur*2, need))
	copy(cp, o.buf)
	copy(cp[cur:], p)
	o.buf = cp
}

func byteStringWriteByte(o *ByteString_Output, b byte, slow func(*ByteString_Output, byte)) {
	if len(o.buf) < cap(o.buf) {
		o.buf = append(o.buf, b)
		return
	}
	if o.buf == nil && o.n < byteStringOutputInlineCap {
		o.inline[o.n] = b
		o.n++
		return
	}
	slow(o, b)
}

func byteStringWriteSlice(o *ByteString_Output, p []byte, slow func(*ByteString_Output, []byte)) (int, error) {
	if len(p) > 0 {
		if len(o.buf)+len(p) <= cap(o.buf) {
			o.buf = append(o.buf, p...)
		} else {
			slow(o, p)
		}
	}
	return len(p), nil
}

func (o *ByteString_Output) WriteByteInt(b int) {
	byteStringWriteByte(o, byte(b), writeByteOutputSlow)
}

func (o *ByteString_Output) WriteByte(b byte) error {
	byteStringWriteByte(o, b, writeByteOutputSlow)
	return nil
}

func (o *ByteString_Output) Write(p []byte) (int, error) {
	return byteStringWriteSlice(o, p, writeSliceOutputSlow)
}

func (o *ByteString_Output) Size() int {
	if o.buf != nil {
		return len(o.buf)
	}
	return int(o.n)
}

//go:noinline
func toByteStringOutputSlow(o *ByteString_Output) ByteString {
	if o.buf == nil {
		n := int(o.n)
		if n == 1 {
			b := o.inline[0]
			dst := byteString1Bytes[b : int(b)+1 : int(b)+1]
			o.buf = dst
			o.n = 0
			return ByteString{data: dst}
		}
		dst := allocWireBytes(n)
		dst = dst[:n:n]
		copy(dst, o.inline[:n])
		o.buf = dst
		o.n = 0
		return ByteString{data: dst}
	}
	n := len(o.buf)
	if o.scratch != nil {
		var dst []byte
		if n == 1 {
			b := o.buf[0]
			dst = byteString1Bytes[b : int(b)+1 : int(b)+1]
		} else {
			dst = allocWireBytes(n)
			dst = dst[:n:n]
			copy(dst, o.buf)
		}
		byteStringScratchPool.Put(o.scratch)
		o.scratch = nil
		o.buf = dst
		return ByteString{data: dst}
	}
	if n == 0 {
		return ByteString_EMPTY
	}
	if cap(o.buf) > 64 && cap(o.buf) > n*2 {
		dst := allocWireBytes(n)
		dst = dst[:n:n]
		copy(dst, o.buf)
		o.buf = dst
		return ByteString{data: dst}
	}
	dst := o.buf[:n:n]
	o.buf = dst
	return ByteString{data: dst}
}

func byteStringFreeze(o *ByteString_Output, slow func(*ByteString_Output) ByteString) ByteString {
	if o.buf == nil {
		if o.n == 0 {
			return ByteString_EMPTY
		}
		return slow(o)
	}
	if o.scratch == nil && len(o.buf) == cap(o.buf) {
		return ByteString{data: o.buf}
	}
	return slow(o)
}

func (o *ByteString_Output) ToByteString() ByteString {
	return byteStringFreeze(o, toByteStringOutputSlow)
}

// ---------------------------------------------------------------------------
// CodedInputStream
// ---------------------------------------------------------------------------

type MessageBuilder interface {
	mergeFromReader(r *protoReader)
}

type CodedInputStream struct {
	rd         protoReader
	br         *bufio.Reader
	scratch    []byte
	scriptSlab []byte
	scriptTail []byte
	scriptStep uint8
}

const (
	scriptSlabMaxItem = 1024 * 1024
	scriptSlabMaxCap  = 4 * 1024 * 1024
)

func (c *CodedInputStream) allocScriptBytes(sz, streamRem int) []byte {
	if sz <= 0 {
		return emptyByteSlice[:0:0]
	}
	if sz > scriptSlabMaxItem {
		if sz > 1<<20 && c.br != nil && c.br.Buffered() == 0 {
			if _, err := c.br.Peek(1); err != nil {
				if err == io.EOF {
					panic(protoTruncatedMessage())
				}
				panic(wrapInputStreamIOErr(err))
			}
		}
		return make([]byte, sz)
	}
	if tailRem := cap(c.scriptTail) - len(c.scriptTail); tailRem >= sz {
		off := len(c.scriptTail)
		end := off + sz
		s := c.scriptTail[:end]
		out := s[off:end:end]
		if cap(s)-end < 512 {
			c.scriptTail = nil
		} else {
			c.scriptTail = s
		}
		return out
	}
	rem := cap(c.scriptSlab) - len(c.scriptSlab)
	if rem < sz {
		if rem >= 512 && rem > cap(c.scriptTail)-len(c.scriptTail) {
			c.scriptTail = c.scriptSlab
		}
		if streamRem >= 2<<20 && c.scriptStep < 8 {
			c.scriptStep = 8
		}
		slabCap := min(max(16384<<min(int(c.scriptStep), 8), sz), scriptSlabMaxCap)
		if streamRem > 0 && slabCap > streamRem {
			slabCap = max(streamRem, sz)
		}
		if c.scriptStep < 8 {
			c.scriptStep++
		}
		c.scriptSlab = make([]byte, 0, slabCap)
	}
	off := len(c.scriptSlab)
	end := off + sz
	c.scriptSlab = c.scriptSlab[:end]
	return c.scriptSlab[off:end:end]
}

type limProtoStream struct {
	cis      *CodedInputStream
	br       *bufio.Reader
	scratch  []byte
	rem      int
	outerRem int
	depth    int
	limit    int
}

func (s *limProtoStream) checkRecursionLimit() {
	limit := s.limit
	if limit == 0 {
		limit = protoDefaultRecursionLimit
	}
	if s.depth >= limit {
		panic(protoRecursionLimitExceeded())
	}
}

func (s *limProtoStream) pushSub(subLen int) int {
	s.checkRecursionLimit()
	saved := s.rem - subLen
	s.outerRem += saved
	s.rem = subLen
	s.depth++
	return saved
}

func (s *limProtoStream) popSub(saved int) {
	if s.rem != 0 {
		panic(protoTruncatedMessage())
	}
	s.depth--
	s.outerRem -= saved
	s.rem = saved
}

func (s *limProtoStream) readByte(allowEOF bool) (byte, bool) {
	if s.rem <= 0 {
		if allowEOF {
			return 0, false
		}
		panic(protoTruncatedMessage())
	}
	b, err := s.br.ReadByte()
	if err != nil {
		if err == io.EOF {
			panic(protoTruncatedMessage())
		}
		panic(wrapInputStreamIOErr(err))
	}
	s.rem--
	return b, true
}

func (s *limProtoStream) readVarint64(allowEOF bool) (uint64, bool) {
	first, ok := s.readByte(allowEOF)
	if !ok {
		return 0, false
	}
	if first < 0x80 {
		return uint64(first), true
	}
	x := uint64(first & 0x7f)
	for shift := uint(7); shift < 64; shift += 7 {
		b, _ := s.readByte(false)
		x |= uint64(b&0x7f) << shift
		if b < 0x80 {
			return x, true
		}
	}
	panic(protoMalformedVarint())
}

func (s *limProtoStream) readRawTag() (int, uint64) {
	raw, ok := s.readVarint64(true)
	if !ok {
		return 0, 0
	}
	tag := int32(raw)
	if uint32(tag)>>3 == 0 {
		panic(protoInvalidTag())
	}
	return int(tag), raw
}

func (s *limProtoStream) readTag() int {
	tag, _ := s.readRawTag()
	return tag
}

func (s *limProtoStream) readSubLen() int {
	raw, _ := s.readVarint64(false)
	sz := int32(raw)
	if sz < 0 {
		panic(protoNegativeSize())
	}
	if int(sz) > s.rem {
		panic(protoTruncatedMessage())
	}
	return int(sz)
}

func (s *limProtoStream) readInto(dst []byte) {
	if len(dst) == 0 {
		return
	}
	if len(dst) > s.rem {
		panic(protoTruncatedMessage())
	}
	n, err := io.ReadFull(s.br, dst)
	s.rem -= n
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			panic(protoTruncatedMessage())
		}
		panic(wrapInputStreamIOErr(err))
	}
}

func (s *limProtoStream) readScratch(n int) []byte {
	if n == 0 {
		return emptyByteSlice[:0:0]
	}
	if n > s.rem {
		panic(protoTruncatedMessage())
	}
	if s.br.Buffered() >= n {
		b, _ := s.br.Peek(n)
		_, _ = s.br.Discard(n)
		s.rem -= n
		return b[:n:n]
	}
	if cap(s.scratch) < n {
		if n > 1<<20 && s.br.Buffered() == 0 {
			if _, err := s.br.Peek(1); err != nil {
				if err == io.EOF {
					panic(protoTruncatedMessage())
				}
				panic(wrapInputStreamIOErr(err))
			}
		}
		s.scratch = make([]byte, n, max(n, max(cap(s.scratch)*2, 4096)))
	} else {
		s.scratch = s.scratch[:n]
	}
	s.readInto(s.scratch)
	return s.scratch[:n:n]
}

func (s *limProtoStream) skipBytes(n int) {
	if n == 0 {
		return
	}
	if n > s.rem {
		panic(protoTruncatedMessage())
	}
	discarded, err := s.br.Discard(n)
	s.rem -= discarded
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			panic(protoTruncatedMessage())
		}
		panic(wrapInputStreamIOErr(err))
	}
}

func (s *limProtoStream) skipField(tag int) {
	switch tag & 7 {
	case protoWireVarint:
		s.readVarint64(false)
	case protoWireFixed64:
		s.skipBytes(8)
	case protoWireBytes:
		n := s.readSubLen()
		s.skipBytes(n)
	case protoWireStartGroup:
		s.skipGroup(tag >> 3)
	case protoWireEndGroup:
		panic(protoInvalidEndTag())
	case protoWireFixed32:
		s.skipBytes(4)
	default:
		panic(protoInvalidWireType())
	}
}

func (s *limProtoStream) skipGroup(number int) {
	s.checkRecursionLimit()
	s.depth++
	defer func() { s.depth-- }()
	for {
		tag := s.readTag()
		if tag == 0 {
			panic(protoTruncatedMessage())
		}
		if tag&7 == protoWireEndGroup {
			if tag>>3 != number {
				panic(protoInvalidEndTag())
			}
			return
		}
		s.skipField(tag)
	}
}

func (s *limProtoStream) appendUnknownField(dst []byte, tag int, rawTag uint64) []byte {
	dst = protoAppendVarint(dst, rawTag)
	switch tag & 7 {
	case protoWireVarint:
		v, _ := s.readVarint64(false)
		return protoAppendVarint(dst, v)
	case protoWireFixed64:
		b := s.readScratch(8)
		return append(dst, b...)
	case protoWireBytes:
		n := s.readSubLen()
		dst = protoAppendVarint(dst, uint64(n))
		b := s.readScratch(n)
		return append(dst, b...)
	case protoWireStartGroup:
		s.checkRecursionLimit()
		s.depth++
		defer func() { s.depth-- }()
		endNum := tag >> 3
		for {
			gt, gRaw := s.readRawTag()
			if gt == 0 {
				panic(protoTruncatedMessage())
			}
			if gt&7 == protoWireEndGroup {
				if gt>>3 != endNum {
					panic(protoInvalidEndTag())
				}
				return protoAppendVarint(dst, gRaw)
			}
			dst = s.appendUnknownField(dst, gt, gRaw)
		}
	case protoWireEndGroup:
		panic(protoInvalidEndTag())
	case protoWireFixed32:
		b := s.readScratch(4)
		return append(dst, b...)
	default:
		panic(protoInvalidWireType())
	}
}

func wrapInputStreamIOErr(err error) error {
	if err == nil {
		return nil
	}
	if th, ok := err.(javalang.Throwable); ok {
		return th
	}
	return javaio.NewIOException(err.Error())
}

func inputStreamRemainingHint(in io.Reader) (int64, bool) {
	if in == nil {
		return 0, false
	}
	if h, ok := in.(javaio.RemainingBytesHinter); ok {
		return h.RemainingBytesHint()
	}
	if l, ok := in.(interface{ Len() int }); ok {
		n := l.Len()
		if n >= 0 {
			return int64(n), true
		}
	}
	return 0, false
}

func readAllFromInputStream(in io.Reader) ([]byte, error) {
	if rem, ok := inputStreamRemainingHint(in); ok && rem >= 0 && rem < 1<<31-1 {
		n := int(rem)
		buf := make([]byte, 0, n+1)
		for {
			if len(buf) == cap(buf) {
				buf = append(buf, 0)[:len(buf)]
			}
			m, err := in.Read(buf[len(buf):cap(buf)])
			buf = buf[:len(buf)+m]
			if err != nil {
				if err == io.EOF {
					return buf, nil
				}
				return buf, err
			}
		}
	}
	return io.ReadAll(in)
}

func CodedInputStream_NewInstanceBytes(b []byte) *CodedInputStream {
	return &CodedInputStream{rd: protoReader{b: b, limit: math.MaxInt32}}
}

func CodedInputStream_NewInstance(in io.Reader) (*CodedInputStream, error) {
	if in == nil {
		return CodedInputStream_NewInstanceBytes(nil), nil
	}
	if rem, ok := inputStreamRemainingHint(in); ok && rem >= 0 && rem <= 256*1024 {
		data, err := readAllFromInputStream(in)
		if err != nil {
			return nil, wrapInputStreamIOErr(err)
		}
		return CodedInputStream_NewInstanceBytes(data), nil
	}
	br := bufio.NewReaderSize(in, 64*1024)
	if _, err := br.Peek(1); err != nil {
		if err == io.EOF {
			return CodedInputStream_NewInstanceBytes(nil), nil
		}
		return nil, wrapInputStreamIOErr(err)
	}
	return &CodedInputStream{
		rd: protoReader{limit: math.MaxInt32},
		br: br,
	}, nil
}

func (c *CodedInputStream) materializeBytes() error {
	if c.br == nil {
		return nil
	}
	data, err := readAllFromInputStream(c.br)
	c.br = nil
	if err != nil {
		return wrapInputStreamIOErr(err)
	}
	c.rd.b = data
	c.rd.pos = 0
	return nil
}

func (c *CodedInputStream) IsAtEnd() bool {
	if c.br != nil {
		_, err := c.br.Peek(1)
		if err == nil {
			return false
		}
		if err == io.EOF {
			c.scratch = nil
			c.scriptSlab = nil
			c.scriptTail = nil
			return true
		}
		panic(wrapInputStreamIOErr(err))
	}
	return c.rd.atEnd()
}

func (c *CodedInputStream) ReadTag() (tag int32, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if ex, ok := rec.(error); ok {
				err = ex
				return
			}
			panic(rec)
		}
	}()
	if c.br != nil {
		raw, ok, th := protoReadStreamVarint64(c.br, true)
		if th != nil {
			return 0, th
		}
		if !ok {
			return 0, nil
		}
		t := int32(raw)
		if uint32(t)>>3 == 0 {
			return 0, protoInvalidTag()
		}
		return t, nil
	}
	return int32(c.rd.readTag()), nil
}

func (c *CodedInputStream) ReadMessage(builder MessageBuilder, _ *ExtensionRegistry) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			if ex, ok := rec.(error); ok {
				err = ex
				return
			}
			panic(rec)
		}
	}()
	if c.br != nil {
		c.rd.checkRecursionLimit()
		rawSize, _, th := protoReadStreamVarint64(c.br, false)
		if th != nil {
			return th
		}
		size := int32(rawSize)
		if size < 0 {
			return protoNegativeSize()
		}
		if tb, ok := builder.(*TypedAst); ok {
			ls := limProtoStream{
				cis:     c,
				br:      c.br,
				scratch: c.scratch,
				rem:     int(size),
				depth:   c.rd.depth + 1,
				limit:   c.rd.limit,
			}
			tb.mergeFromStream(&ls)
			if ls.rem != 0 {
				return protoTruncatedMessage()
			}
			c.scratch = ls.scratch
			return nil
		}
		var msgBytes []byte
		if size > 0 {
			if int(size) > 16<<20 && c.br.Buffered() == 0 {
				if _, peekErr := c.br.Peek(1); peekErr != nil {
					if peekErr == io.EOF {
						return protoTruncatedMessage()
					}
					return wrapInputStreamIOErr(peekErr)
				}
			}
			msgBytes = make([]byte, int(size))
			if th := protoReadFull(c.br, msgBytes); th != nil {
				return th
			}
		}
		sub := protoReader{b: msgBytes, depth: c.rd.depth + 1, limit: c.rd.limit}
		builder.mergeFromReader(&sub)
		return nil
	}
	msgBytes := c.rd.readMessageBytes()
	sub := c.rd.sub(msgBytes)
	builder.mergeFromReader(&sub)
	return nil
}

func (c *CodedInputStream) ResetSizeCounter() {}

func (c *CodedInputStream) SetRecursionLimit(limit int) int {
	if limit < 0 {
		panic(javalang.NewIllegalArgumentException("Recursion limit cannot be negative"))
	}
	old := c.rd.limit
	c.rd.limit = limit
	return old
}

// ---------------------------------------------------------------------------
// TextFormat.Printer
// ---------------------------------------------------------------------------

type TextFormat_Printer struct{}

var defaultTextFormatPrinter = &TextFormat_Printer{}

func TextFormat_PrinterInstance() *TextFormat_Printer {
	return defaultTextFormatPrinter
}

func (p *TextFormat_Printer) PrintToString(msg any) string {
	if msg == nil {
		return ""
	}
	if s, ok := msg.(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprint(msg)
}

// ---------------------------------------------------------------------------
// NodeProperty (com.google.javascript.jscomp.serialization.NodeProperty)
// ---------------------------------------------------------------------------

type NodeProperty int32

const (
	NodeProperty_NODE_PROPERTY_UNSPECIFIED NodeProperty = iota
	NodeProperty_IS_PARENTHESIZED
	NodeProperty_SYNTHETIC
	NodeProperty_ADDED_BLOCK
	NodeProperty_IS_CONSTANT_NAME
	NodeProperty_IS_NAMESPACE
	NodeProperty_DIRECT_EVAL
	NodeProperty_FREE_CALL
	NodeProperty_UNUSED_11
	NodeProperty_REFLECTED_OBJECT
	NodeProperty_STATIC_MEMBER
	NodeProperty_GENERATOR_FN
	NodeProperty_ARROW_FN
	NodeProperty_ASYNC_FN
	NodeProperty_YIELD_ALL
	NodeProperty_EXPORT_DEFAULT
	NodeProperty_EXPORT_ALL_FROM
	NodeProperty_CONSTANT_VAR_FLAGS
	NodeProperty_IS_GENERATOR_MARKER
	NodeProperty_IS_GENERATOR_SAFE
	NodeProperty_COMPUTED_PROP_METHOD
	NodeProperty_COMPUTED_PROP_GETTER
	NodeProperty_COMPUTED_PROP_SETTER
	NodeProperty_COMPUTED_PROP_VARIABLE
	NodeProperty_COLOR_FROM_CAST
	NodeProperty_NON_INDEXABLE
	NodeProperty_GOOG_MODULE
	NodeProperty_DELETED
	NodeProperty_MODULE_ALIAS
	NodeProperty_IS_UNUSED_PARAMETER
	NodeProperty_MODULE_EXPORT
	NodeProperty_IS_SHORTHAND_PROPERTY
	NodeProperty_ES6_MODULE
	NodeProperty_START_OF_OPT_CHAIN
	NodeProperty_TRAILING_COMMA
	NodeProperty_IS_INFERRED_CONSTANT
	NodeProperty_IS_DECLARED_CONSTANT
	NodeProperty_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION
	NodeProperty_MUTATES_GLOBAL_STATE
	NodeProperty_MUTATES_THIS
	NodeProperty_MUTATES_ARGUMENTS
	NodeProperty_THROWS
	NodeProperty_CLOSURE_UNAWARE_SHADOW
	NodeProperty_PRIVATE_IDENTIFIER
	NodeProperty_UNRECOGNIZED
)

var nodeProperty_names = [...]string{
	"NODE_PROPERTY_UNSPECIFIED",
	"IS_PARENTHESIZED",
	"SYNTHETIC",
	"ADDED_BLOCK",
	"IS_CONSTANT_NAME",
	"IS_NAMESPACE",
	"DIRECT_EVAL",
	"FREE_CALL",
	"UNUSED_11",
	"REFLECTED_OBJECT",
	"STATIC_MEMBER",
	"GENERATOR_FN",
	"ARROW_FN",
	"ASYNC_FN",
	"YIELD_ALL",
	"EXPORT_DEFAULT",
	"EXPORT_ALL_FROM",
	"CONSTANT_VAR_FLAGS",
	"IS_GENERATOR_MARKER",
	"IS_GENERATOR_SAFE",
	"COMPUTED_PROP_METHOD",
	"COMPUTED_PROP_GETTER",
	"COMPUTED_PROP_SETTER",
	"COMPUTED_PROP_VARIABLE",
	"COLOR_FROM_CAST",
	"NON_INDEXABLE",
	"GOOG_MODULE",
	"DELETED",
	"MODULE_ALIAS",
	"IS_UNUSED_PARAMETER",
	"MODULE_EXPORT",
	"IS_SHORTHAND_PROPERTY",
	"ES6_MODULE",
	"START_OF_OPT_CHAIN",
	"TRAILING_COMMA",
	"IS_INFERRED_CONSTANT",
	"IS_DECLARED_CONSTANT",
	"SYNTHESIZED_UNFULFILLED_NAME_DECLARATION",
	"MUTATES_GLOBAL_STATE",
	"MUTATES_THIS",
	"MUTATES_ARGUMENTS",
	"THROWS",
	"CLOSURE_UNAWARE_SHADOW",
	"PRIVATE_IDENTIFIER",
	"UNRECOGNIZED",
}

var nodeProperty_numbers = [...]int32{
	0, 1, 3, 4, 6, 7, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	21, 22, 23, 24, 25, 26, 28, 29, 30, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, 48, 49, 50, 51, -1,
}

func (p NodeProperty) String() string { return nodeProperty_names[p] }
func (p NodeProperty) Name() string   { return nodeProperty_names[p] }

func (p NodeProperty) GetNumber() int32 {
	if p == NodeProperty_UNRECOGNIZED {
		panic(javalang.NewIllegalArgumentException("Can't get the number of an unknown enum value."))
	}
	return nodeProperty_numbers[p]
}

var nodePropertyByNumber = func() [64]*NodeProperty {
	var out [64]*NodeProperty
	var vals [len(nodeProperty_numbers)]NodeProperty
	for i, n := range nodeProperty_numbers {
		if n >= 0 && int(n) < len(out) && NodeProperty(i) != NodeProperty_UNRECOGNIZED {
			vals[i] = NodeProperty(i)
			out[n] = &vals[i]
		}
	}
	return out
}()

func NodeProperty_ForNumber(value int32) *NodeProperty {
	if uint32(value) < uint32(len(nodePropertyByNumber)) {
		return nodePropertyByNumber[value]
	}
	return nil
}

// ---------------------------------------------------------------------------
// Mapping.OriginalMapping (com.google.debugging.sourcemap.proto.Mapping)
// ---------------------------------------------------------------------------

type Mapping_OriginalMapping_Precision int32

const (
	Mapping_OriginalMapping_Precision_UNKNOWN_PRECISION Mapping_OriginalMapping_Precision = 0
	Mapping_OriginalMapping_Precision_EXACT             Mapping_OriginalMapping_Precision = 1
	Mapping_OriginalMapping_Precision_APPROXIMATE_LINE  Mapping_OriginalMapping_Precision = 2
)

var mappingPrecisionNames = [...]string{"UNKNOWN_PRECISION", "EXACT", "APPROXIMATE_LINE"}

func (p Mapping_OriginalMapping_Precision) String() string {
	if int(p) >= 0 && int(p) < len(mappingPrecisionNames) {
		return mappingPrecisionNames[p]
	}
	return "UNKNOWN_PRECISION"
}

func (p Mapping_OriginalMapping_Precision) Name() string     { return p.String() }
func (p Mapping_OriginalMapping_Precision) GetNumber() int32 { return int32(p) }

type Mapping_OriginalMapping struct {
	OriginalFile   string
	LineNumber     int32
	ColumnPosition int32
	Identifier     string
	Precision      Mapping_OriginalMapping_Precision
}

type Mapping_OriginalMapping_Builder = Mapping_OriginalMapping

var mappingOriginalMappingDefault = &Mapping_OriginalMapping{}

func Mapping_OriginalMapping_GetDefaultInstance() *Mapping_OriginalMapping {
	return mappingOriginalMappingDefault
}

func Mapping_OriginalMapping_NewBuilder() *Mapping_OriginalMapping_Builder {
	return &Mapping_OriginalMapping{}
}

func (m *Mapping_OriginalMapping) GetOriginalFile() string   { return m.OriginalFile }
func (m *Mapping_OriginalMapping) GetLineNumber() int32      { return m.LineNumber }
func (m *Mapping_OriginalMapping) GetColumnPosition() int32  { return m.ColumnPosition }
func (m *Mapping_OriginalMapping) GetIdentifier() string     { return m.Identifier }
func (m *Mapping_OriginalMapping) GetPrecision() Mapping_OriginalMapping_Precision {
	if uint32(m.Precision) <= 2 {
		return m.Precision
	}
	return Mapping_OriginalMapping_Precision_UNKNOWN_PRECISION
}

func (m *Mapping_OriginalMapping) SetOriginalFile(v string) *Mapping_OriginalMapping {
	m.OriginalFile = v
	return m
}
func (m *Mapping_OriginalMapping) SetLineNumber(v int32) *Mapping_OriginalMapping {
	m.LineNumber = v
	return m
}
func (m *Mapping_OriginalMapping) SetColumnPosition(v int32) *Mapping_OriginalMapping {
	m.ColumnPosition = v
	return m
}
func (m *Mapping_OriginalMapping) SetIdentifier(v string) *Mapping_OriginalMapping {
	m.Identifier = v
	return m
}
func (m *Mapping_OriginalMapping) SetPrecision(v Mapping_OriginalMapping_Precision) *Mapping_OriginalMapping {
	m.Precision = v
	return m
}
func (m *Mapping_OriginalMapping) Build() *Mapping_OriginalMapping { return m }
func (m *Mapping_OriginalMapping) ToBuilder() *Mapping_OriginalMapping {
	cp := *m
	return &cp
}

func (m *Mapping_OriginalMapping) Equals(other any) bool {
	o, ok := other.(*Mapping_OriginalMapping)
	return ok && o != nil && *m == *o
}

func (m *Mapping_OriginalMapping) HashCode() int32 {
	h := javalang.StringHashCode(m.OriginalFile)
	h = 31*h + m.LineNumber
	h = 31*h + m.ColumnPosition
	h = 31*h + javalang.StringHashCode(m.Identifier)
	h = 31*h + int32(m.Precision)
	return h
}

// ---------------------------------------------------------------------------
// Instrumentation Reporter Protos
// (com.google.javascript.jscomp.instrumentation.reporter.proto.*)
// ---------------------------------------------------------------------------

type InstrumentationPoint_Type int32

const (
	InstrumentationPoint_Type_TYPE_UNKNOWN   InstrumentationPoint_Type = 0
	InstrumentationPoint_Type_FUNCTION       InstrumentationPoint_Type = 1
	InstrumentationPoint_Type_BRANCH         InstrumentationPoint_Type = 2
	InstrumentationPoint_Type_BRANCH_DEFAULT InstrumentationPoint_Type = 3
)

var instrumentationPointTypeNames = [...]string{
	"TYPE_UNKNOWN",
	"FUNCTION",
	"BRANCH",
	"BRANCH_DEFAULT",
}

func (t InstrumentationPoint_Type) Name() string {
	if int(t) >= 0 && int(t) < len(instrumentationPointTypeNames) {
		return instrumentationPointTypeNames[t]
	}
	return "UNRECOGNIZED"
}

func (t InstrumentationPoint_Type) String() string   { return t.Name() }
func (t InstrumentationPoint_Type) GetNumber() int32 { return int32(t) }

func InstrumentationPoint_Type_ValueOf(name string) InstrumentationPoint_Type {
	for i, n := range instrumentationPointTypeNames {
		if n == name {
			return InstrumentationPoint_Type(i)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant InstrumentationPoint.Type." + name))
}

type InstrumentationPoint struct {
	fileName_     string
	functionName_ string
	type_         InstrumentationPoint_Type
	lineNumber_   int32
	columnNumber_ int32
}

type InstrumentationPoint_Builder = InstrumentationPoint

func InstrumentationPoint_NewBuilder() *InstrumentationPoint_Builder {
	return &InstrumentationPoint{}
}

func (p *InstrumentationPoint) GetFileName() string               { return p.fileName_ }
func (p *InstrumentationPoint) GetFunctionName() string           { return p.functionName_ }
func (p *InstrumentationPoint) GetType() InstrumentationPoint_Type { return p.type_ }
func (p *InstrumentationPoint) GetLineNumber() int32              { return p.lineNumber_ }
func (p *InstrumentationPoint) GetColumnNumber() int32            { return p.columnNumber_ }

func (p *InstrumentationPoint) SetFileName(v string) *InstrumentationPoint {
	p.fileName_ = v
	return p
}
func (p *InstrumentationPoint) SetFunctionName(v string) *InstrumentationPoint {
	p.functionName_ = v
	return p
}
func (p *InstrumentationPoint) SetType(v InstrumentationPoint_Type) *InstrumentationPoint {
	p.type_ = v
	return p
}
func (p *InstrumentationPoint) SetLineNumber(v int32) *InstrumentationPoint {
	p.lineNumber_ = v
	return p
}
func (p *InstrumentationPoint) SetColumnNumber(v int32) *InstrumentationPoint {
	p.columnNumber_ = v
	return p
}
func (p *InstrumentationPoint) Build() *InstrumentationPoint { return p }

func (p *InstrumentationPoint) Equals(other any) bool {
	o, ok := other.(*InstrumentationPoint)
	return ok && o != nil && *p == *o
}

func (p *InstrumentationPoint) HashCode() int32 {
	h := javalang.StringHashCode(p.fileName_)
	h = 31*h + javalang.StringHashCode(p.functionName_)
	h = 31*h + int32(p.type_)
	h = 31*h + p.lineNumber_
	h = 31*h + p.columnNumber_
	return h
}

type InstrumentationPointStats_Presence int32

const (
	InstrumentationPointStats_Presence_PRESENCE_UNKNOWN   InstrumentationPointStats_Presence = 0
	InstrumentationPointStats_Presence_PRESENT            InstrumentationPointStats_Presence = 1
	InstrumentationPointStats_Presence_STATICALLY_REMOVED InstrumentationPointStats_Presence = 2
)

func (p InstrumentationPointStats_Presence) GetNumber() int32 { return int32(p) }
func (p InstrumentationPointStats_Presence) String() string {
	switch p {
	case InstrumentationPointStats_Presence_PRESENT:
		return "PRESENT"
	case InstrumentationPointStats_Presence_STATICALLY_REMOVED:
		return "STATICALLY_REMOVED"
	default:
		return "PRESENCE_UNKNOWN"
	}
}

type InstrumentationPointStats struct {
	point_         *InstrumentationPoint
	timesExecuted_ int64
	pointPresence_ InstrumentationPointStats_Presence
}

type InstrumentationPointStats_Builder = InstrumentationPointStats

func InstrumentationPointStats_NewBuilder() *InstrumentationPointStats_Builder {
	return &InstrumentationPointStats{}
}

func (s *InstrumentationPointStats) GetPoint() *InstrumentationPoint {
	if s.point_ != nil {
		return s.point_
	}
	return &InstrumentationPoint{}
}
func (s *InstrumentationPointStats) GetTimesExecuted() int64 { return s.timesExecuted_ }
func (s *InstrumentationPointStats) GetPointPresence() InstrumentationPointStats_Presence {
	return s.pointPresence_
}
func (s *InstrumentationPointStats) SetPoint(v *InstrumentationPoint) *InstrumentationPointStats {
	s.point_ = v
	return s
}
func (s *InstrumentationPointStats) SetTimesExecuted(v int64) *InstrumentationPointStats {
	s.timesExecuted_ = v
	return s
}
func (s *InstrumentationPointStats) SetPointPresence(v InstrumentationPointStats_Presence) *InstrumentationPointStats {
	s.pointPresence_ = v
	return s
}
func (s *InstrumentationPointStats) Build() *InstrumentationPointStats { return s }

type FileProfile struct {
	fileName_                   string
	instrumentationPointsStats_ []*InstrumentationPointStats
}

type FileProfile_Builder = FileProfile

func FileProfile_NewBuilder() *FileProfile_Builder { return &FileProfile{} }

func (f *FileProfile) GetFileName() string { return f.fileName_ }
func (f *FileProfile) GetInstrumentationPointsStatsList() []*InstrumentationPointStats {
	return f.instrumentationPointsStats_
}
func (f *FileProfile) SetFileName(v string) *FileProfile {
	f.fileName_ = v
	return f
}
func (f *FileProfile) AddAllInstrumentationPointsStats(v []*InstrumentationPointStats) *FileProfile {
	f.instrumentationPointsStats_ = append(f.instrumentationPointsStats_, v...)
	return f
}
func (f *FileProfile) Build() *FileProfile { return f }

type ReportProfile struct {
	fileProfile_ []*FileProfile
}

type ReportProfile_Builder = ReportProfile

func ReportProfile_NewBuilder() *ReportProfile_Builder { return &ReportProfile{} }

func (r *ReportProfile) GetFileProfileList() []*FileProfile { return r.fileProfile_ }
func (r *ReportProfile) AddAllFileProfile(v []*FileProfile) *ReportProfile {
	r.fileProfile_ = append(r.fileProfile_, v...)
	return r
}
func (r *ReportProfile) Build() *ReportProfile { return r }
