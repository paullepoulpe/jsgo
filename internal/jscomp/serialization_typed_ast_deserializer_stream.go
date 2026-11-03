// Port of com/google/javascript/jscomp/serialization/TypedAstDeserializer.java @ closure-compiler v20261028

/*
 * Copyright 2021 The Closure Compiler Authors.
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
	"bufio"
	"bytes"
	"io"
	"math"
	"strconv"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// typedAstDeserializer_WIRETYPE_LENGTH_DELIMITED is
// WireFormat.WIRETYPE_LENGTH_DELIMITED.
const (
	typedAstDeserializer_WIRETYPE_LENGTH_DELIMITED = 2
	typedAstList_TYPED_ASTS_FIELD_NUMBER           = 1
)

func typedAstValidateShardTag(rawTag uint64) jutil.Throwable {
	tag := int32(rawTag)
	if uint32(tag)>>3 == 0 {
		return protoInvalidTag()
	}
	if jutil.Ushr32(tag, 3) != typedAstList_TYPED_ASTS_FIELD_NUMBER ||
		int(tag)&7 != typedAstDeserializer_WIRETYPE_LENGTH_DELIMITED {
		return NewInvalidProtocolBufferException(
			"Unexpected field number " +
				strconv.Itoa(int(jutil.Ushr32(tag, 3))) +
				" or wire type " +
				strconv.Itoa(int(tag)&7))
	}
	return nil
}

func typedAstDeserializer_decodeShardBytes(messageBytes []byte, buildStringPool func(*TypedAst) bool) (*TypedAst, *StringPool) {
	rd := &protoReader{b: messageBytes, depth: 1, limit: math.MaxInt32}
	typedAst := &TypedAst{}
	typedAst.mergeFromReader(rd)
	var stringShard *StringPool
	if buildStringPool == nil || buildStringPool(typedAst) {
		stringShard = StringPool_FromProto(typedAst.GetStringPool())
	}
	return typedAst, stringShard
}

type byteAndBlockReader interface {
	io.Reader
	io.ByteReader
}

func typedAstDeserializer_scanTypedAstStream(
	typedAstsStream io.Reader,
	processShardBatch func([][]byte),
	processLargeShard func(io.Reader, int),
) {
	const batchSize = 32
	if buf, ok := typedAstsStream.(*bytes.Buffer); ok {
		data := buf.Next(buf.Len())
		pos := 0
		readNextFromSlice := func() ([]byte, bool, jutil.Throwable) {
			rawTag, ok, err := protoReadSliceVarint64(data, &pos, true)
			if err != nil || !ok {
				return nil, false, err
			}
			if err := typedAstValidateShardTag(rawTag); err != nil {
				return nil, false, err
			}
			rawSize, _, err := protoReadSliceVarint64(data, &pos, false)
			if err != nil {
				return nil, false, err
			}
			size := int32(rawSize)
			if size < 0 {
				return nil, false, protoNegativeSize()
			}
			if int(size) > len(data)-pos {
				return nil, false, protoTruncatedMessage()
			}
			v := data[pos : pos+int(size) : pos+int(size)]
			pos += int(size)
			return v, true, nil
		}
		firstMsg, ok, firstErr := readNextFromSlice()
		if firstErr != nil {
			panic(firstErr)
		}
		if !ok {
			return
		}
		if pos >= len(data) {
			processShardBatch([][]byte{firstMsg})
			return
		}
		batch := make([][]byte, 0, batchSize)
		batch = append(batch, firstMsg)
		for {
			var readErr jutil.Throwable
			for len(batch) < batchSize {
				msg, hasMore, err := readNextFromSlice()
				if err != nil {
					readErr = err
					break
				}
				if !hasMore {
					break
				}
				batch = append(batch, msg)
			}
			if len(batch) > 0 {
				processShardBatch(batch)
				batch = batch[:0]
			}
			if readErr != nil {
				panic(readErr)
			}
			if pos >= len(data) {
				break
			}
		}
		return
	}

	const (
		probeSize          = 32 * 1024
		streamLargeShardAt = 1 << 20
	)
	bbr, ok := typedAstsStream.(byteAndBlockReader)
	if !ok {
		bbr = bufio.NewReaderSize(typedAstsStream, probeSize)
	}

	var probe [probeSize]byte
	readNextTagAndLen := func() (int, bool, jutil.Throwable) {
		uTag, ok, err := protoReadStreamVarint64(bbr, true)
		if err != nil || !ok {
			return 0, false, err
		}
		if uTag > math.MaxUint32 {
			return 0, false, protoInvalidTag()
		}
		if err := typedAstValidateShardTag(uTag); err != nil {
			return 0, false, err
		}
		rawLen, _, err := protoReadStreamVarint64(bbr, false)
		if err != nil {
			return 0, false, err
		}
		msgLen := int(int32(rawLen))
		if msgLen < 0 {
			return 0, false, protoNegativeSize()
		}
		return msgLen, true, nil
	}
	readBodyBytes := func(msgLen int) ([]byte, jutil.Throwable) {
		// Do not pool messageBytes buffers because LazyAst.script uses zero-copy readBytesView().
		if msgLen <= probeSize {
			messageBytes := make([]byte, msgLen)
			if err := protoReadFull(bbr, messageBytes); err != nil {
				return nil, err
			}
			return messageBytes, nil
		}
		if err := protoReadFull(bbr, probe[:]); err != nil {
			return nil, err
		}
		messageBytes := make([]byte, msgLen)
		copy(messageBytes, probe[:])
		if err := protoReadFull(bbr, messageBytes[probeSize:]); err != nil {
			return nil, err
		}
		return messageBytes, nil
	}

	batch := make([][]byte, 0, batchSize)
	for {
		var readErr jutil.Throwable
		var reachedEOF bool
		var largeMsgLen int
		hasLarge := false
		for len(batch) < batchSize {
			msgLen, hasMore, err := readNextTagAndLen()
			if err != nil {
				readErr = err
				break
			}
			if !hasMore {
				reachedEOF = true
				break
			}
			if msgLen > streamLargeShardAt {
				largeMsgLen = msgLen
				hasLarge = true
				break
			}
			msg, err := readBodyBytes(msgLen)
			if err != nil {
				readErr = err
				break
			}
			batch = append(batch, msg)
		}
		if len(batch) > 0 {
			processShardBatch(batch)
			batch = batch[:0]
		}
		if readErr != nil {
			panic(readErr)
		}
		if hasLarge {
			processLargeShard(bbr, largeMsgLen)
		}
		if reachedEOF {
			break
		}
	}
}

type typedAstStreamReader struct {
	r         io.Reader
	win       []byte
	rpos      int
	wpos      int
	remaining int
	scratch   []byte
}

func (s *typedAstStreamReader) totalLeft() int {
	return (s.wpos - s.rpos) + s.remaining
}

func (s *typedAstStreamReader) fill(need int) {
	if s.rpos > 0 {
		copy(s.win, s.win[s.rpos:s.wpos])
		s.wpos -= s.rpos
		s.rpos = 0
	}
	for s.wpos < need {
		toRead := min(len(s.win)-s.wpos, s.remaining)
		if toRead <= 0 {
			panic(protoTruncatedMessage())
		}
		n, err := s.r.Read(s.win[s.wpos : s.wpos+toRead])
		s.wpos += n
		s.remaining -= n
		if n == 0 && err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				panic(protoTruncatedMessage())
			}
			panic(jutil.NewIOException(err.Error()))
		}
	}
}

func (s *typedAstStreamReader) readByte(minLeft int) byte {
	if s.totalLeft() <= minLeft {
		panic(protoTruncatedMessage())
	}
	if s.rpos == s.wpos {
		s.fill(1)
	}
	b := s.win[s.rpos]
	s.rpos++
	return b
}

func (s *typedAstStreamReader) readVarint64(minLeft int) uint64 {
	b := s.readByte(minLeft)
	x := uint64(b & 0x7f)
	if b < 0x80 {
		return x
	}
	for shift := uint(7); shift < 64; shift += 7 {
		b = s.readByte(minLeft)
		x |= uint64(b&0x7f) << shift
		if b < 0x80 {
			return x
		}
	}
	panic(protoMalformedVarint())
}

func (s *typedAstStreamReader) readTag(minLeft int) int {
	if s.totalLeft() <= minLeft {
		return 0
	}
	tag := int32(s.readVarint64(minLeft))
	if uint32(tag)>>3 == 0 {
		panic(protoInvalidTag())
	}
	return int(tag)
}

func (s *typedAstStreamReader) readLength(minLeft int) int {
	n := int(int32(s.readVarint64(minLeft)))
	if n < 0 {
		panic(protoNegativeSize())
	}
	if s.totalLeft()-n < minLeft {
		panic(protoTruncatedMessage())
	}
	return n
}

func (s *typedAstStreamReader) readSliceView(n int, minLeft int) []byte {
	if s.totalLeft()-n < minLeft {
		panic(protoTruncatedMessage())
	}
	if n == 0 {
		return nil
	}
	if n <= s.wpos-s.rpos {
		v := s.win[s.rpos : s.rpos+n : s.rpos+n]
		s.rpos += n
		return v
	}
	if n <= len(s.win) {
		s.fill(n)
		v := s.win[0:n:n]
		s.rpos = n
		return v
	}
	if cap(s.scratch) < n {
		s.scratch = make([]byte, n)
	} else {
		s.scratch = s.scratch[:n]
	}
	avail := s.wpos - s.rpos
	if avail > 0 {
		copy(s.scratch[:avail], s.win[s.rpos:s.wpos])
		s.rpos = 0
		s.wpos = 0
	}
	need := n - avail
	if err := protoReadFull(s.r, s.scratch[avail:n]); err != nil {
		panic(err)
	}
	s.remaining -= need
	return s.scratch[:n:n]
}

func streamMergeSubmessage[T interface{ mergeFromReader(*protoReader) }](sr *typedAstStreamReader, dst T, minLeft int) {
	pr := protoReader{b: sr.readSliceView(sr.readLength(minLeft), minLeft), depth: 2, limit: math.MaxInt32}
	dst.mergeFromReader(&pr)
}

func (s *typedAstStreamReader) skipField(tag int, minLeft int) {
	switch tag & 7 {
	case protoWireVarint:
		_ = s.readVarint64(minLeft)
	case protoWireFixed64:
		_ = s.readSliceView(8, minLeft)
	case protoWireBytes:
		_ = s.readSliceView(s.readLength(minLeft), minLeft)
	case protoWireStartGroup:
		groupNum := tag >> 3
		for {
			gt := s.readTag(minLeft)
			if gt == 0 {
				panic(protoTruncatedMessage())
			}
			if gt&7 == protoWireEndGroup {
				if gt>>3 != groupNum {
					panic(protoInvalidEndTag())
				}
				return
			}
			s.skipField(gt, minLeft)
		}
	case protoWireEndGroup:
		panic(protoInvalidEndTag())
	case protoWireFixed32:
		_ = s.readSliceView(4, minLeft)
	default:
		panic(protoInvalidWireType())
	}
}

func typedAstDeserializer_streamLargeShard(
	r io.Reader,
	msgLen int,
	buildStringPool func(*TypedAst) bool,
) (*TypedAst, *StringPool) {
	sr := typedAstStreamReader{
		r:         r,
		win:       make([]byte, 32768),
		remaining: msgLen,
	}
	m := &TypedAst{}
	typeProtoSlab := jutil.NewFixedSlab[TypeProto](512)
	edgeSlab := jutil.NewFixedSlab[SubtypingEdge](512)
	lazyAstSlab := jutil.NewFixedSlab[LazyAst](256)
	sourceFileSlab := jutil.NewFixedSlab[SourceFileProto](256)

	var scriptSlab []byte
	allocScriptBytes := func(src []byte) []byte {
		n := len(src)
		if n == 0 {
			return nil
		}
		if n > 32760 {
			return bytes.Clone(src)
		}
		if len(scriptSlab) < n {
			scriptSlab = make([]byte, 32760)
		}
		dst := scriptSlab[:n:n]
		scriptSlab = scriptSlab[n:]
		copy(dst, src)
		return dst
	}

	type deferredWtf8 struct {
		idx  int
		data []byte
	}
	var (
		hasStringPool bool
		stringIDs     []strpool.ID
		maxAsciiLen   int
		maxLength     int32
		deferred      []deferredWtf8
	)

	for sr.totalLeft() > 0 {
		tag := sr.readTag(0)
		switch tag {
		case 0:
		case 10:
			subLen := sr.readLength(0)
			endLeft := sr.totalLeft() - subLen
			var x *TypePool
			if m.typePool_ != nil {
				x = m.typePool_
				x.cachedSize = 0
			} else {
				x = &TypePool{}
			}
			for sr.totalLeft() > endLeft {
				subTag := sr.readTag(endLeft)
				switch subTag {
				case 10:
					tp := typeProtoSlab.New()
					streamMergeSubmessage(&sr, tp, endLeft)
					x.type_ = append(x.type_, tp)
				case 18:
					ed := edgeSlab.New()
					streamMergeSubmessage(&sr, ed, endLeft)
					x.disambiguationEdges_ = append(x.disambiguationEdges_, ed)
				case 26:
					var dbg *TypePool_DebugInfo
					if x.debugInfo_ != nil {
						dbg = x.debugInfo_
						dbg.cachedSize = 0
					} else {
						dbg = &TypePool_DebugInfo{}
					}
					streamMergeSubmessage(&sr, dbg, endLeft)
					x.debugInfo_ = dbg
				default:
					sr.skipField(subTag, endLeft)
				}
			}
			m.typePool_ = x
		case 18:
			hasStringPool = true
			subLen := sr.readLength(0)
			endLeft := sr.totalLeft() - subLen
			if stringIDs == nil {
				stringIDs = make([]strpool.ID, 1, min(subLen/12+2, 4_000_000))
				stringIDs[0] = strpool.Intern("")
			}
			for sr.totalLeft() > endLeft {
				subTag := sr.readTag(endLeft)
				switch subTag {
				case 10:
					raw := sr.readSliceView(sr.readLength(endLeft), endLeft)
					if wtf8_isASCIIBytes(raw) {
						if len(raw) > maxAsciiLen {
							maxAsciiLen = len(raw)
						}
						stringIDs = append(stringIDs, strpool.InternBytes(raw))
					} else {
						stringIDs = append(stringIDs, 0)
						deferred = append(deferred, deferredWtf8{idx: len(stringIDs) - 1, data: bytes.Clone(raw)})
					}
				case 16:
					maxLength = int32(sr.readVarint64(endLeft))
				default:
					sr.skipField(subTag, endLeft)
				}
			}
		case 26, 34:
			la := lazyAstSlab.New()
			streamMergeSubmessage(&sr, la, 0)
			la.script_ = allocScriptBytes(la.script_)
			if tag == 26 {
				m.externAst_ = append(m.externAst_, la)
			} else {
				m.codeAst_ = append(m.codeAst_, la)
			}
		case 42:
			subLen := sr.readLength(0)
			endLeft := sr.totalLeft() - subLen
			var x *SourceFilePool
			if m.sourceFilePool_ != nil {
				x = m.sourceFilePool_
				x.cachedSize = 0
			} else {
				x = &SourceFilePool{}
			}
			for sr.totalLeft() > endLeft {
				subTag := sr.readTag(endLeft)
				switch subTag {
				case 10:
					sf := sourceFileSlab.New()
					streamMergeSubmessage(&sr, sf, endLeft)
					x.sourceFile_ = append(x.sourceFile_, sf)
				default:
					sr.skipField(subTag, endLeft)
				}
			}
			m.sourceFilePool_ = x
		case 50:
			var x *ExternsSummary
			if m.externsSummary_ != nil {
				x = m.externsSummary_
				x.cachedSize = 0
			} else {
				x = &ExternsSummary{}
			}
			streamMergeSubmessage(&sr, x, 0)
			m.externsSummary_ = x
		case 58:
			raw := sr.readSliceView(sr.readLength(0), 0)
			if !utf8Valid(raw) {
				panic(protoInvalidUtf8())
			}
			m.runtimeLibraryToInject_ = append(m.runtimeLibraryToInject_, string(raw))
		default:
			sr.skipField(tag, 0)
		}
	}

	var stringShard *StringPool
	if buildStringPool == nil || buildStringPool(m) {
		if !hasStringPool {
			stringShard = StringPool_FromProto(m.GetStringPool())
		} else {
			if maxLength < 0 {
				_ = Wtf8_Decoder_(int(maxLength))
			}
			if maxAsciiLen > int(maxLength) {
				panic(jutil.NewArrayIndexOutOfBoundsExceptionForLength(int(maxLength), int(maxLength)))
			}
			if len(deferred) > 0 {
				decoder := Wtf8_Decoder_(int(maxLength))
				for _, d := range deferred {
					stringIDs[d.idx] = decoder.DecodeInternID(d.data)
				}
			}
			stringShard = &StringPool{
				maxLength: int(maxLength),
				pool:      rhino.NewRhinoStringPool_LazyInternedStringListFromIDs(stringIDs),
			}
		}
	}
	return m, stringShard
}
