/*
 * Copyright 2026 The Closure Compiler Authors.
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
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"runtime"
	"slices"
	"sync"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

const (
	binarySaveStateMagic           = "JSGOSAV1"
	binarySaveStateHeaderSize      = 16
	binarySaveStateVersion         = uint32(1)
	binarySaveStateFormatTag       = uint32(0x4a53_4202) // "JSB\x02": flat binary AST arena v2
	maxBinarySaveStateSectionBytes = uint64(16) << 30    // 16 GiB sanity cap per binary section
	maxBinarySaveStateCount        = 64 << 20            // 64M entries sanity cap per table
)

// Per-node compact binary flag bits following the 1-byte NodeKind ordinal.
const (
	binNodeValMask     uint16 = 0x3 // 0=none, 1=string, 2=double, 3=template
	binNodeValNone     uint16 = 0
	binNodeValString   uint16 = 1
	binNodeValDouble   uint16 = 2
	binNodeValTemplate uint16 = 3
	binNodeHasChildren uint16 = 1 << 2
	binNodeLineShift          = 3
	binNodeLineMask    uint16 = 0x3 << binNodeLineShift // 0=zero, 1=int8, 2=int32
	binNodeColShift           = 5
	binNodeColMask     uint16 = 0x3 << binNodeColShift // 0=zero, 1=int8, 2=int32
	binNodePropsShift         = 7
	binNodePropsMask   uint16 = 0x3 << binNodePropsShift // 0=zero, 1=uint16, 2=uint32, 3=uint64
	binNodeHasType     uint16 = 1 << 9
	binNodeHasOrigName uint16 = 1 << 10
	binNodeHasJsdoc    uint16 = 1 << 11
	binNodeHasSource   uint16 = 1 << 12
	binNodeKnownFlags  uint16 = 0x1FFF
	binNodeMaxBytes           = 64 // 3-byte header + max(10 val, 5 endIdx, 4 line, 4 col, 8 props, 5 type, 5 orig, 5 jsdoc, 5 sf)
)

type replayErrorReader struct {
	prefix []byte
	err    error
}

func newReplayErrorReader(prefix []byte, err error) io.Reader {
	return &replayErrorReader{prefix: bytes.Clone(prefix), err: err}
}

func (r *replayErrorReader) Read(p []byte) (int, error) {
	if len(r.prefix) > 0 {
		n := copy(p, r.prefix)
		r.prefix = r.prefix[n:]
		return n, nil
	}
	return 0, r.err
}

func (c *Compiler) readAndValidateBinaryStateHeader(r io.Reader) {
	var hdr [binarySaveStateHeaderSize]byte
	if _, err := io.ReadFull(r, hdr[:]); err != nil {
		binaryFail("Truncated binary compiler state header: " + err.Error())
	}
	if string(hdr[:8]) != binarySaveStateMagic {
		panic(jutil.NewIllegalStateExceptionMsg("Invalid binary compiler state magic"))
	}
	if v := binary.LittleEndian.Uint32(hdr[8:12]); v != binarySaveStateVersion {
		panic(jutil.NewIllegalStateExceptionMsg(fmt.Sprintf("Unsupported binary compiler state version %d (expected %d)", v, binarySaveStateVersion)))
	}
	if tag := binary.LittleEndian.Uint32(hdr[12:16]); tag != binarySaveStateFormatTag {
		panic(jutil.NewIllegalStateExceptionMsg(fmt.Sprintf("Incompatible binary compiler state format tag 0x%08x (expected 0x%08x)", tag, binarySaveStateFormatTag)))
	}
}

type binaryEncodedScript struct {
	data             []byte
	sourceFile       int32
	sourceMappingURL string
}

func (c *Compiler) saveBinaryState(outputStream io.Writer) {
	c.runInCompilerThread(func() {
		tracer := c.NewTracer("serializeCompilerState")
		nonClosingOut := NewCompiler_NonClosingOutputStream(outputStream)
		var outerHdr [binarySaveStateHeaderSize]byte
		copy(outerHdr[:8], binarySaveStateMagic)
		binary.LittleEndian.PutUint32(outerHdr[8:12], binarySaveStateVersion)
		binary.LittleEndian.PutUint32(outerHdr[12:16], binarySaveStateFormatTag)
		if _, err := nonClosingOut.Write(outerHdr[:]); err != nil {
			binaryFail(err.Error())
		}
		var compressionStream io.Writer
		if wrapper := c.options.GetStateCompressionWrapper(); wrapper != nil {
			compressionStream = wrapper(nonClosingOut)
		} else {
			compressionStream = newSaveStateGzipOutputStream(nonClosingOut, c.options.GetNumParallelThreads())
		}
		bufferedOut := compiler_newBufferedOutputStream(compressionStream, 64*1024)
		stateBytes := c.GetCompilerStateProtoBuilder().Build().ToByteArray()
		c.StopTracer(tracer, "serializeCompilerState")

		tracer = c.NewTracer("serializeTypedAst")
		NewRemoveCastNodes(c).Process(c.externsRoot, c.jsRoot)
		serOpts := &SerializationOptions{
			includeDebugInfo: c.GetOptions().ShouldSerializeExtraDebugInfo(),
			runtimeLibraries: emptyRuntimeLibraries,
		}
		c.writeBinaryTypedAstPayload(bufferedOut.w, stateBytes, serOpts)
		c.StopTracer(tracer, "serializeTypedAst")
		bufferedOut.Close()
	})
}

func initBinaryTypedAstSerializer(s *TypedAstSerializer, externsRoot rhino.Node, jsRoot rhino.Node) {
	jutil.CheckArgument(externsRoot.IsRoot())
	jutil.CheckArgument(jsRoot.IsRoot())
	if s.compiler.HasOptimizationColors() {
		if root := s.compiler.GetRoot(); !root.IsNil() && externsRoot.GetParent() == root {
			s.recordUsedProps = true
			cs := NewColorSerializer(s.serializationMode, s.stringPool.Put, func(prop string) bool {
				_, ok := s.lookupUsedPropIndex(prop)
				return ok
			})
			cs.lookupUsedPropIndexFn = s.lookupUsedPropIndex
			s.colorTypeSerializer = NewTypedAstSerializer_ColorTypeSerializer(cs, s.compiler.GetColorRegistry())
		} else {
			s.colorTypeSerializer = s.CreateColorTypeSerializer(s.compiler, s.serializationMode, s.stringPool)
		}
	} else if s.compiler.HasTypeCheckingRun() {
		s.jsTypeSerializer = TypedAstSerializer_CreateJSTypeSerializer(s.compiler, s.serializationMode, s.stringPool)
	}
}

func (c *Compiler) writeBinaryTypedAstPayload(bw *bufio.Writer, stateBytes []byte, serOpts *SerializationOptions) {
	s := NewTypedAstSerializer(c, serOpts)
	initBinaryTypedAstSerializer(s, c.externsRoot, c.jsRoot)

	numThreads := 1
	if opts := c.GetOptions(); opts != nil && opts.GetNumParallelThreads() > 1 {
		numThreads = opts.GetNumParallelThreads()
	}

	var allScripts []rhino.Node
	appendScripts := func(root rhino.Node) {
		for sc := root.GetFirstChild(); !sc.IsNil(); sc = sc.GetNext() {
			if !NodeUtil_IsFromTypeSummary(sc) {
				allScripts = append(allScripts, sc)
			}
		}
	}
	appendScripts(c.externsRoot)
	numExtern := len(allScripts)
	appendScripts(c.jsRoot)
	totalScripts := len(allScripts)

	encodedScripts := make([]binaryEncodedScript, totalScripts)
	if numThreads > 1 && s.jsTypeSerializer == nil && totalScripts >= 4 {
		batchSize := min(totalScripts, max(numThreads*8, 16))
		encoders := make([]scriptWireEncoder, batchSize)
		workerScratch := make([][]byte, numThreads)
		workerSlabs := make([][]byte, numThreads)
		for base := 0; base < totalScripts; base += batchSize {
			batch := allScripts[base:min(base+batchSize, totalScripts)]
			batchLen := len(batch)
			runSerializerBatchParallel(batchLen, numThreads, func(_ int, slot int) {
				encoders[slot].buildScript(batch[slot], s)
			})
			for slot := 0; slot < batchLen; slot++ {
				s.applyFixupsSequential(&encoders[slot])
			}
			runSerializerBatchParallel(batchLen, numThreads, func(workerIdx int, slot int) {
				encodedScripts[base+slot] = encodeBinaryScript(&encoders[slot], &workerScratch[workerIdx], &workerSlabs[workerIdx])
			})
		}
	} else {
		var scratch []byte
		for i, script := range allScripts {
			s.scriptEnc.buildScript(script, s)
			s.applyFixupsSequential(&s.scriptEnc)
			encodedScripts[i] = encodeBinaryScript(&s.scriptEnc, &scratch, &s.scriptByteSlab)
		}
	}

	var externPropPtrs []uint32
	if externProps := c.GetExternPropSet(); externProps != nil {
		externPropPtrs = make([]uint32, 0, externProps.Len())
		for prop := range externProps.AllStrings() {
			externPropPtrs = append(externPropPtrs, uint32(s.stringPool.Put(prop)))
		}
	}
	var runtimeLibs []string
	if libs := serOpts.RuntimeLibraries(); len(libs) > 0 {
		runtimeLibs = slices.Clone(libs)
	}
	typePoolBytes := marshalBinaryTypePool(s.generateTypePool())

	w := binarySectionWriter{bw: bw}
	w.u32(uint32(len(stateBytes)))
	w.bytes(stateBytes)

	ids := s.stringPool.ids
	numStrings := max(len(ids)-1, 0)
	var totalStringBytes uint64
	for i := 1; i <= numStrings; i++ {
		totalStringBytes += uint64(len(strpool.String(ids[i])))
	}
	if totalStringBytes > math.MaxUint32 {
		panic(jutil.NewIllegalStateExceptionMsg("String table exceeds maximum binary SaveState size"))
	}
	w.u32(uint32(s.stringPool.maxLength))
	w.u32(uint32(numStrings))
	w.u32(uint32(totalStringBytes))
	for i := 1; i <= numStrings; i++ {
		w.u32(uint32(len(strpool.String(ids[i]))))
	}
	for i := 1; i <= numStrings; i++ {
		w.str(strpool.String(ids[i]))
	}

	w.u32(uint32(len(typePoolBytes)))
	w.bytes(typePoolBytes)

	w.u32(uint32(len(s.sourceFileList)))
	var sfScratch []byte
	for _, sf := range s.sourceFileList {
		p := sf.GetProto()
		_ = p.serializedSize()
		sfScratch = p.appendTo(sfScratch[:0])
		w.u32(uint32(len(sfScratch)))
		w.bytes(sfScratch)
	}

	w.u32(uint32(len(externPropPtrs)))
	for _, ptr := range externPropPtrs {
		w.u32(ptr)
	}
	w.u32(uint32(len(runtimeLibs)))
	for _, lib := range runtimeLibs {
		w.u32(uint32(len(lib)))
		w.str(lib)
	}

	w.u32(uint32(numExtern))
	w.u32(uint32(totalScripts - numExtern))
	var totalSourceMapURLBytes, totalScriptBytes uint64
	for i := range encodedScripts {
		sc := &encodedScripts[i]
		w.u32(uint32(sc.sourceFile))
		w.u32(uint32(len(sc.data)))
		w.u32(uint32(len(sc.sourceMappingURL)))
		totalSourceMapURLBytes += uint64(len(sc.sourceMappingURL))
		totalScriptBytes += uint64(len(sc.data))
	}
	if totalSourceMapURLBytes > math.MaxUint32 || totalScriptBytes > maxBinarySaveStateSectionBytes {
		panic(jutil.NewIllegalStateExceptionMsg("Script arena exceeds maximum binary SaveState size"))
	}
	w.u32(uint32(totalSourceMapURLBytes))
	for i := range encodedScripts {
		w.str(encodedScripts[i].sourceMappingURL)
	}
	w.u64(totalScriptBytes)
	for i := range encodedScripts {
		w.bytes(encodedScripts[i].data)
	}
}

type binarySectionWriter struct {
	bw  *bufio.Writer
	tmp [8]byte
}

func (w *binarySectionWriter) u32(v uint32) {
	binary.LittleEndian.PutUint32(w.tmp[:4], v)
	if _, err := w.bw.Write(w.tmp[:4]); err != nil {
		binaryFail(err.Error())
	}
}

func (w *binarySectionWriter) u64(v uint64) {
	binary.LittleEndian.PutUint64(w.tmp[:8], v)
	if _, err := w.bw.Write(w.tmp[:8]); err != nil {
		binaryFail(err.Error())
	}
}

func (w *binarySectionWriter) bytes(p []byte) {
	if len(p) > 0 {
		if _, err := w.bw.Write(p); err != nil {
			binaryFail(err.Error())
		}
	}
}

func (w *binarySectionWriter) str(s string) {
	if len(s) > 0 {
		if _, err := w.bw.WriteString(s); err != nil {
			binaryFail(err.Error())
		}
	}
}

func marshalBinaryTypePool(tp *TypePool) []byte {
	if tp == nil {
		return nil
	}
	if sz := ensureTypePoolMemoizedSizes(tp); sz > 0 {
		return tp.appendTo(make([]byte, 0, sz))
	}
	return nil
}

func unmarshalBinaryTypePool(data []byte) *TypePool {
	if len(data) == 0 {
		return TypePool_GetDefaultInstance()
	}
	tp := &TypePool{}
	typeProtoSlab := jutil.NewFixedSlab[TypeProto](min(max(len(data)/16, 16), 512))
	edgeSlab := jutil.NewFixedSlab[SubtypingEdge](min(max(len(data)/32, 8), 512))
	rd := protoReader{b: data, depth: 1, limit: math.MaxInt32}
	for !rd.atEnd() {
		switch tag := rd.readTag(); tag {
		case 0:
		case 10:
			t := typeProtoSlab.New()
			sub := protoReader{b: rd.readMessageBytes(), depth: 2, limit: math.MaxInt32}
			t.mergeFromReader(&sub)
			tp.type_ = append(tp.type_, t)
		case 18:
			ed := edgeSlab.New()
			sub := protoReader{b: rd.readMessageBytes(), depth: 2, limit: math.MaxInt32}
			ed.mergeFromReader(&sub)
			tp.disambiguationEdges_ = append(tp.disambiguationEdges_, ed)
		case 26:
			if tp.debugInfo_ == nil {
				tp.debugInfo_ = &TypePool_DebugInfo{}
			} else {
				tp.debugInfo_.cachedSize = 0
			}
			sub := protoReader{b: rd.readMessageBytes(), depth: 2, limit: math.MaxInt32}
			tp.debugInfo_.mergeFromReader(&sub)
		default:
			rd.skipField(tag)
		}
	}
	return tp
}

func encodeBinaryScript(enc *scriptWireEncoder, scratch *[]byte, byteSlab *[]byte) binaryEncodedScript {
	nodes, jsdocs, jsdocKinds := enc.nodes, enc.jsdocs, enc.jsdocKinds
	maxBytes := 12 + len(jsdocs)*35 + len(jsdocKinds)*5 + len(nodes)*binNodeMaxBytes
	if cap(*scratch) < maxBytes {
		*scratch = make([]byte, max(maxBytes, 64*1024))
	}
	buf := (*scratch)[:maxBytes]
	binary.LittleEndian.PutUint32(buf[0:4], uint32(len(nodes)))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(len(jsdocs)))
	binary.LittleEndian.PutUint32(buf[8:12], uint32(len(jsdocKinds)))
	pos := 12

	putUvar32 := func(v uint32) {
		pos += binary.PutUvarint(buf[pos:], uint64(v))
	}

	for i := range jsdocs {
		d := &jsdocs[i]
		putUvar32(d.kindOff)
		putUvar32(d.kindLen)
		putUvar32(uint32(d.licensePtr))
		putUvar32(uint32(d.descPtr))
		putUvar32(uint32(d.meaningPtr))
		putUvar32(uint32(d.altMsgPtr))
		putUvar32(uint32(d.cuModePtr))
	}
	for _, k := range jsdocKinds {
		putUvar32(uint32(k))
	}

	putRelDelta := func(v int32, shift uint) uint16 {
		if v == 0 {
			return 0
		}
		if v >= -128 && v <= 127 {
			buf[pos] = byte(int8(v))
			pos++
			return 1 << shift
		}
		binary.LittleEndian.PutUint32(buf[pos:], uint32(v))
		pos += 4
		return 2 << shift
	}
	putOptUvar32 := func(cond bool, bit uint16, v uint32) uint16 {
		if !cond {
			return 0
		}
		pos += binary.PutUvarint(buf[pos:], uint64(v))
		return bit
	}

	for idx := range nodes {
		fn := &nodes[idx]
		kindOrd := nodeKind_fromNumber(int32(fn.kind))
		hdrPos := pos
		pos += 3
		var flags uint16

		switch fn.valCase {
		case astValCase_String:
			flags |= binNodeValString
			putUvar32(uint32(fn.valBits))
		case astValCase_Double:
			flags |= binNodeValDouble
			binary.LittleEndian.PutUint64(buf[pos:], fn.valBits)
			pos += 8
		case astValCase_Template:
			flags |= binNodeValTemplate
			rawPtr, cookedPtr := fn.templatePointers()
			putUvar32(uint32(rawPtr))
			putUvar32(uint32(cookedPtr + 1))
		}
		flags |= putOptUvar32(fn.endIdx != uint32(idx+1), binNodeHasChildren, fn.endIdx)
		flags |= putRelDelta(fn.relLine, binNodeLineShift)
		flags |= putRelDelta(fn.relCol, binNodeColShift)
		if fn.boolProps != 0 {
			if fn.boolProps <= 0xFFFF {
				flags |= 1 << binNodePropsShift
				binary.LittleEndian.PutUint16(buf[pos:], uint16(fn.boolProps))
				pos += 2
			} else if fn.boolProps <= 0xFFFFFFFF {
				flags |= 2 << binNodePropsShift
				binary.LittleEndian.PutUint32(buf[pos:], uint32(fn.boolProps))
				pos += 4
			} else {
				flags |= 3 << binNodePropsShift
				binary.LittleEndian.PutUint64(buf[pos:], fn.boolProps)
				pos += 8
			}
		}
		flags |= putOptUvar32(fn.hasType, binNodeHasType, uint32(fn.typePtr))
		flags |= putOptUvar32(fn.origName != 0, binNodeHasOrigName, uint32(fn.origName))
		flags |= putOptUvar32(fn.jsdocIdx >= 0, binNodeHasJsdoc, uint32(fn.jsdocIdx))
		flags |= putOptUvar32(fn.sourceFile != 0, binNodeHasSource, uint32(fn.sourceFile))

		buf[hdrPos] = uint8(kindOrd)
		binary.LittleEndian.PutUint16(buf[hdrPos+1:], flags)
	}

	var out []byte
	if pos > 262144 {
		out = bytes.Clone(buf[:pos])
	} else {
		slab := *byteSlab
		if len(slab) < pos {
			slab = make([]byte, min(max(64*1024, pos*4), 1024*1024))
		}
		out = slab[:pos:pos]
		copy(out, buf[:pos])
		*byteSlab = slab[pos:]
	}
	return binaryEncodedScript{data: out, sourceFile: enc.lazyAst.sourceFile_, sourceMappingURL: enc.lazyAst.sourceMappingUrl_}
}

func binaryFail(msg string) { panic(jutil.NewIOException(msg)) }

type binaryStreamReader struct {
	r   io.Reader
	tmp [8]byte
}

func (sr *binaryStreamReader) readFull(dst []byte) {
	if _, err := io.ReadFull(sr.r, dst); err != nil {
		binaryFail("Truncated binary compiler state: " + err.Error())
	}
}

func (sr *binaryStreamReader) u32() uint32 {
	sr.readFull(sr.tmp[:4])
	return binary.LittleEndian.Uint32(sr.tmp[:4])
}

func (sr *binaryStreamReader) u64() uint64 {
	sr.readFull(sr.tmp[:8])
	return binary.LittleEndian.Uint64(sr.tmp[:8])
}

func (sr *binaryStreamReader) count() uint32 {
	n := sr.u32()
	if n > maxBinarySaveStateCount {
		binaryFail("Corrupt binary compiler state table count")
	}
	return n
}

const (
	binaryStreamDirectAllocMaxBytes = 4 << 20
	binaryStreamReadChunkBytes      = 16 << 20
)

func (sr *binaryStreamReader) bytes(n uint64) []byte {
	if n == 0 {
		return nil
	}
	if n > maxBinarySaveStateSectionBytes || n > uint64(math.MaxInt) {
		binaryFail("Binary compiler state section exceeds maximum size")
	}
	if n <= binaryStreamDirectAllocMaxBytes {
		buf := make([]byte, int(n))
		sr.readFull(buf)
		return buf
	}
	first := make([]byte, binaryStreamDirectAllocMaxBytes)
	sr.readFull(first)
	buf := make([]byte, int(n))
	copy(buf, first)
	for off := len(first); off < len(buf); {
		end := min(off+binaryStreamReadChunkBytes, len(buf))
		sr.readFull(buf[off:end])
		off = end
	}
	return buf
}

func (sr *binaryStreamReader) readScratch(scratch *[]byte, n uint64) []byte {
	if n == 0 || n > 4<<20 {
		return sr.bytes(n)
	}
	if cap(*scratch) < int(n) {
		*scratch = make([]byte, max(int(n), 256))
	}
	sr.readFull((*scratch)[:int(n)])
	return (*scratch)[:int(n)]
}

func (c *Compiler) deserializeBinaryCompilerState(in io.Reader) {
	sr := &binaryStreamReader{r: in}
	stateProto := JSCompilerStateProto_ParseFrom(sr.bytes(uint64(sr.u32())))
	externFiles, codeFiles, allInputFiles := c.collectRestoreInputFiles()

	numWorkers := 1
	if c.options != nil && c.options.GetNumParallelThreads() > 1 {
		numWorkers = c.options.GetNumParallelThreads()
	} else if n := runtime.GOMAXPROCS(0); n > 1 {
		numWorkers = n
	}

	stringPool := readBinaryStringPool(sr, numWorkers)

	typePoolLen := sr.u32()
	var typePoolBytes []byte
	includeTypes := stateProto.GetTypeCheckingHasRun()
	if includeTypes {
		typePoolBytes = sr.bytes(uint64(typePoolLen))
	} else if typePoolLen > 0 {
		if uint64(typePoolLen) > maxBinarySaveStateSectionBytes {
			binaryFail("Binary TypePool section exceeds maximum size")
		}
		if _, err := io.CopyN(io.Discard, sr.r, int64(typePoolLen)); err != nil {
			binaryFail("Truncated binary compiler state: " + err.Error())
		}
	}

	numFiles := int(sr.count())
	sourceFileProtos := make([]*SourceFileProto, numFiles)
	sfBacking := make([]SourceFileProto, numFiles)
	var sfScratch []byte
	for i := 0; i < numFiles; i++ {
		sfBytes := sr.readScratch(&sfScratch, uint64(sr.u32()))
		sfProto := &sfBacking[i]
		rd := protoReader{b: sfBytes, depth: 1, limit: math.MaxInt32}
		sfProto.mergeFromReader(&rd)
		sourceFileProtos[i] = sfProto
	}

	numExternProps := int(sr.count())
	externPropRaw := sr.bytes(uint64(numExternProps) * 4)
	stringPoolCount := uint32(stringPool.GetInternedStrings().Len())

	numRuntimeLibs := int(sr.count())
	runtimeLibs := make([]string, numRuntimeLibs)
	for i := 0; i < numRuntimeLibs; i++ {
		runtimeLibs[i] = string(sr.readScratch(&sfScratch, uint64(sr.u32())))
	}

	numExternScripts := int(sr.count())
	numCodeScripts := int(sr.count())
	totalScripts := numExternScripts + numCodeScripts
	if uint64(totalScripts) > maxBinarySaveStateCount {
		binaryFail("Corrupt binary script count")
	}
	scriptMetaRaw := sr.bytes(uint64(totalScripts) * 12)
	sourceMapBlob := sr.bytes(uint64(sr.u32()))
	scriptBlob := sr.bytes(sr.u64())

	deserializedAst := c.buildBinaryDeserializedAst(
		allInputFiles, includeTypes, stringPool, stringPoolCount, typePoolBytes,
		sourceFileProtos, numExternProps, externPropRaw, runtimeLibs,
		totalScripts, scriptMetaRaw, sourceMapBlob, scriptBlob,
	)
	c.applyDeserializedCompilerState(stateProto, externFiles, codeFiles, deserializedAst)
}

func readBinaryStringPool(sr *binaryStreamReader, numWorkers int) *StringPool {
	maxLength := int(int32(sr.u32()))
	if maxLength < 0 {
		binaryFail("Corrupt binary StringPool maxLength")
	}
	numStrings := int(sr.count())
	totalStringBytes := sr.u32()
	lenBytes := sr.bytes(uint64(numStrings) * 4)
	blob := sr.bytes(uint64(totalStringBytes))

	ids := make([]strpool.ID, numStrings+1)
	ids[0] = strpool.EmptyID
	useParallel := numWorkers > 1 && numStrings >= 4096
	var offsets []uint32
	if useParallel {
		offsets = make([]uint32, numStrings+1)
	}
	var off uint64
	for i := 0; i < numStrings; i++ {
		n := uint64(binary.LittleEndian.Uint32(lenBytes[i*4:]))
		if off+n > uint64(len(blob)) {
			binaryFail("Corrupt binary StringPool string length")
		}
		if useParallel {
			offsets[i+1] = uint32(off + n)
		} else {
			ids[i+1] = strpool.InternBytes(blob[off : off+n])
		}
		off += n
	}
	if off != uint64(len(blob)) {
		binaryFail("Corrupt binary StringPool total byte length")
	}
	if useParallel {
		workers := min(numWorkers, (numStrings+2047)/2048)
		var wg sync.WaitGroup
		wg.Add(workers)
		for w := 0; w < workers; w++ {
			start, end := (w*numStrings)/workers, ((w+1)*numStrings)/workers
			go func() {
				defer wg.Done()
				for i := start; i < end; i++ {
					ids[i+1] = strpool.InternBytes(blob[offsets[i]:offsets[i+1]])
				}
			}()
		}
		wg.Wait()
	}
	return &StringPool{maxLength: maxLength, pool: rhino.NewRhinoStringPool_LazyInternedStringListFromIDs(ids)}
}

func (c *Compiler) buildBinaryDeserializedAst(
	allInputFiles []*SourceFile, includeTypes bool, stringPool *StringPool, stringPoolCount uint32,
	typePoolBytes []byte, sourceFileProtos []*SourceFileProto, numExternProps int, externPropRaw []byte,
	runtimeLibs []string, totalScripts int, scriptMetaRaw, sourceMapBlob, scriptBlob []byte,
) *TypedAstDeserializer_DeserializedAst {
	requiredOrder := make([]*SourceFile, 0, len(allInputFiles))
	requiredSet := make(map[*SourceFile]struct{}, len(allInputFiles))
	seenNames := make(map[string]*SourceFile, len(allInputFiles))
	for _, f := range allInputFiles {
		if _, dup := requiredSet[f]; dup {
			continue
		}
		requiredSet[f] = struct{}{}
		requiredOrder = append(requiredOrder, f)
		if prev, dup := seenNames[f.GetName()]; dup {
			panic(jutil.NewIllegalArgumentExceptionMsg(
				"Multiple entries with same key: " + f.GetName() + "=" + f.String() + " and " + f.GetName() + "=" + prev.String()))
		}
		seenNames[f.GetName()] = f
	}
	syntheticExterns := c.SYNTHETIC_EXTERNS_FILE
	deserializer := newTypedAstDeserializer(
		syntheticExterns, nil, requiredOrder, requiredSet, typedAstDeserializer_Mode_FULL_AST, includeTypes)
	deserializer.addToFilePool(syntheticExterns.GetName(), syntheticExterns)
	for _, file := range requiredOrder {
		deserializer.addToFilePool(file.GetName(), file)
	}
	c.InitRuntimeLibraryTypedAsts(deserializer.colorPoolBuilder)

	fileShard := make([]*SourceFile, 0, len(sourceFileProtos))
	for _, p := range sourceFileProtos {
		if existingFile := deserializer.findExistingFile(p.GetFilename()); existingFile != nil {
			existingFile.RestoreCachedStateFrom(p)
			fileShard = append(fileShard, existingFile)
		} else {
			protoFile := SourceFile_FromProto(p)
			deserializer.addToFilePool(p.GetFilename(), protoFile)
			fileShard = append(fileShard, protoFile)
		}
	}

	readScriptMeta := func(i int) (sfIdx int, dataLen, smLen uint64) {
		off := i * 12
		return int(binary.LittleEndian.Uint32(scriptMetaRaw[off:])),
			uint64(binary.LittleEndian.Uint32(scriptMetaRaw[off+4:])),
			uint64(binary.LittleEndian.Uint32(scriptMetaRaw[off+8:]))
	}
	containsRequiredInput := false
	var smPos, scPos uint64
	for i := 0; i < totalScripts; i++ {
		sfIdx, dataLen, smLen := readScriptMeta(i)
		if sfIdx <= 0 || sfIdx > len(fileShard) || scPos+dataLen > uint64(len(scriptBlob)) || smPos+smLen > uint64(len(sourceMapBlob)) {
			binaryFail("Corrupt binary script metadata")
		}
		scPos, smPos = scPos+dataLen, smPos+smLen
		if !containsRequiredInput {
			_, containsRequiredInput = requiredSet[fileShard[sfIdx-1]]
		}
	}
	if scPos != uint64(len(scriptBlob)) || smPos != uint64(len(sourceMapBlob)) {
		binaryFail("Corrupt binary script blob length")
	}

	if containsRequiredInput {
		for _, lib := range runtimeLibs {
			deserializer.runtimeLibraries.Add(lib)
		}
		var colorShard *ColorPool_ShardView
		if deserializer.colorPoolBuilder != nil {
			colorShard = deserializer.colorPoolBuilder.AddShard(unmarshalBinaryTypePool(typePoolBytes), stringPool)
		}
		if deserializer.externProperties.IsEmpty() && numExternProps > orderedMapIndexThreshold {
			deserializer.externProperties = NewOrderedSetWithCapacity[string](numExternProps)
		}
		for i := 0; i < numExternProps; i++ {
			ptr := binary.LittleEndian.Uint32(externPropRaw[i*4:])
			if ptr >= stringPoolCount {
				binaryFail("Corrupt binary extern property string pointer")
			}
			deserializer.externProperties.Add(stringPool.Get(int(ptr)))
		}

		resolveSourceMaps := c.GetOptions().GetResolveSourceMapAnnotations()
		parseInlineSourceMaps := c.GetOptions().GetParseInlineSourceMaps()
		deserBacking := make([]ScriptNodeDeserializer, totalScripts)
		var smOff, scOff int
		for i := 0; i < totalScripts; i++ {
			sfIdx, dataLenU, smLenU := readScriptMeta(i)
			dataLen, smLen := int(dataLenU), int(smLenU)
			scriptData := scriptBlob[scOff : scOff+dataLen : scOff+dataLen]
			scOff += dataLen
			var smURL string
			if smLen > 0 {
				smURL = string(sourceMapBlob[smOff : smOff+smLen])
				smOff += smLen
			}
			file := fileShard[sfIdx-1]
			if _, ok := requiredSet[file]; !ok && file != syntheticExterns {
				continue
			}
			if file.IsWeak() {
				if _, exists := deserializer.typedAstFilesystem[file]; !exists {
					deserializer.typedAstFilesystem[file] = typedAstFileEntry{weakFile: file}
				}
				continue
			}
			sd := &deserBacking[i]
			*sd = ScriptNodeDeserializer{
				sourceFile: file, binaryScript: scriptData, sourceMappingURL: smURL,
				colorPoolShard: colorShard, stringPool: stringPool, filePool: fileShard,
			}
			if file == syntheticExterns {
				deserializer.syntheticExternsDeserializers = append(deserializer.syntheticExternsDeserializers, sd)
			} else if _, exists := deserializer.typedAstFilesystem[file]; !exists {
				deserializer.typedAstFilesystem[file] = typedAstFileEntry{sd: sd}
			}
			typedAstDeserializer_addInputSourceMap(sd, c, resolveSourceMaps, parseInlineSourceMaps)
		}
	}

	deserializer.typedAstFilesystem[syntheticExterns] = typedAstFileEntry{
		synthetic: newTypedAstDeserializer_SyntheticExternsDeserializer(syntheticExterns, deserializer.syntheticExternsDeserializers),
	}
	var registry *colors.ColorRegistry
	if deserializer.colorPoolBuilder != nil {
		registry = deserializer.colorPoolBuilder.Build().GetRegistry()
	}
	return typedAstDeserializer_DeserializedAst_create(
		deserializer.typedAstFilesystem, registry, deserializer.externProperties, deserializer.runtimeLibraries)
}

type binSliceReader struct {
	b   []byte
	pos int
}

func (r *binSliceReader) fail(msg string) { binaryFail(msg) }

func (r *binSliceReader) need(n int, msg string) int {
	p := r.pos
	if p+n > len(r.b) {
		r.fail(msg)
	}
	r.pos = p + n
	return p
}

func (r *binSliceReader) u16(msg string) uint16 {
	return binary.LittleEndian.Uint16(r.b[r.need(2, msg):])
}
func (r *binSliceReader) u32(msg string) uint32 {
	return binary.LittleEndian.Uint32(r.b[r.need(4, msg):])
}
func (r *binSliceReader) u64(msg string) uint64 {
	return binary.LittleEndian.Uint64(r.b[r.need(8, msg):])
}

func (r *binSliceReader) uvarint32(truncMsg string) uint32 {
	p := r.pos
	b := r.b
	if p < len(b) && b[p] < 0x80 {
		r.pos = p + 1
		return uint32(b[p])
	}
	v, n := binary.Uvarint(b[p:])
	if n <= 0 {
		if n == 0 {
			r.fail(truncMsg)
		}
		r.fail("Overflow in binary script varint")
	}
	if v > math.MaxUint32 {
		r.fail("Overflow in binary script uint32 varint")
	}
	r.pos = p + n
	return uint32(v)
}

func (r *binSliceReader) strPtr(poolLen uint32, truncMsg, oobMsg string) int32 {
	u := r.uvarint32(truncMsg)
	if u >= poolLen {
		r.fail(oobMsg)
	}
	return int32(u)
}

func (r *binSliceReader) relDelta(mode uint16, truncMsg, modeMsg string) int32 {
	switch mode {
	case 0:
		return 0
	case 1:
		return int32(int8(r.b[r.need(1, truncMsg)]))
	case 2:
		return int32(r.u32(truncMsg))
	default:
		r.fail(modeMsg)
		return 0
	}
}

type binaryParentFrame struct {
	node, template, shadowTemplate, shadowChild rhino.Node
	endIdx                                      uint32
	origNamePtr                                 int32
	isShadowHost                                bool
}

func (d *ScriptNodeDeserializer) deserializeBinaryWithAllocator(alloc *rhino.NodeAllocator) rhino.Node {
	defer func() {
		if rec := recover(); rec != nil {
			t, ok := rec.(jutil.Throwable)
			if !ok {
				t = jutil.AsThrowable(rec)
			}
			if ex, ok := t.(jutil.IIOException); ok {
				panic(NewMalformedTypedAstExceptionWithCause(d.sourceFile, ex))
			}
			panic(rec)
		}
	}()

	b := d.binaryScript
	if len(b) < 12 {
		binaryFail("Truncated binary script header")
	}
	nodeCount, jsdocCount, jsdocKindCount := binary.LittleEndian.Uint32(b[0:4]), binary.LittleEndian.Uint32(b[4:8]), binary.LittleEndian.Uint32(b[8:12])
	if nodeCount == 0 || nodeCount > maxBinarySaveStateCount || jsdocCount > maxBinarySaveStateCount || jsdocKindCount > maxBinarySaveStateCount {
		binaryFail("Invalid binary script header counts")
	}
	if uint64(12)+uint64(jsdocCount)*7+uint64(jsdocKindCount) > uint64(len(b)) {
		binaryFail("Truncated binary script jsdoc table")
	}

	rd := binSliceReader{b: b, pos: 12}
	internedStrings := d.stringPool.GetInternedStrings()
	stringPoolLen := uint32(internedStrings.Len())

	var jsdocStack [16]flatJsdoc
	var jsdocs []flatJsdoc
	if jsdocCount > 0 {
		if jsdocCount <= uint32(len(jsdocStack)) {
			jsdocs = jsdocStack[:jsdocCount]
		} else {
			jsdocs = make([]flatJsdoc, jsdocCount)
		}
		for ji := uint32(0); ji < jsdocCount; ji++ {
			kindOff := rd.uvarint32("Truncated binary script jsdoc table")
			kindLen := rd.uvarint32("Truncated binary script jsdoc table")
			licPtr := rd.uvarint32("Truncated binary script jsdoc table")
			descPtr := rd.uvarint32("Truncated binary script jsdoc table")
			meanPtr := rd.uvarint32("Truncated binary script jsdoc table")
			altPtr := rd.uvarint32("Truncated binary script jsdoc table")
			cuPtr := rd.uvarint32("Truncated binary script jsdoc table")
			if uint64(kindOff)+uint64(kindLen) > uint64(jsdocKindCount) ||
				licPtr >= stringPoolLen || descPtr >= stringPoolLen ||
				meanPtr >= stringPoolLen || altPtr >= stringPoolLen || cuPtr >= stringPoolLen {
				rd.fail("Binary OptimizationJsdoc entry out of bounds")
			}
			jsdocs[ji] = flatJsdoc{
				kindOff:    kindOff,
				kindLen:    kindLen,
				licensePtr: int32(licPtr),
				descPtr:    int32(descPtr),
				meaningPtr: int32(meanPtr),
				altMsgPtr:  int32(altPtr),
				cuModePtr:  int32(cuPtr),
			}
		}
	}
	var jsdocKindStack [32]int32
	var jsdocKinds []int32
	if jsdocKindCount > 0 {
		if jsdocKindCount <= uint32(len(jsdocKindStack)) {
			jsdocKinds = jsdocKindStack[:jsdocKindCount]
		} else {
			jsdocKinds = make([]int32, jsdocKindCount)
		}
		for ki := uint32(0); ki < jsdocKindCount; ki++ {
			jsdocKinds[ki] = int32(rd.uvarint32("Truncated binary script jsdoc kinds"))
		}
	}

	var featureCollector parsing.FeatureCollector
	featureCollector.Reset()
	var ctxBuf [64]parsing.FeatureCollector_FeatureContext
	contextStack := append(ctxBuf[:0], parsing.FeatureCollector_FeatureContext_NONE)
	var frameBuf [64]binaryParentFrame
	stack := frameBuf[:0]
	var jsdocBuilder rhino.JSDocInfo_Builder
	rootTmpl := d.sourceFile.getOrCreateSourceInfoTemplate()
	var previousLine, previousColumn int
	var scriptNode rhino.Node

	popCompleted := func(curIdx uint32) {
		for len(stack) > 0 && stack[len(stack)-1].endIdx <= curIdx {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if top.isShadowHost {
				jutil.CheckStateMsg1(rhino.IR_MayBeExpression(top.shadowChild), "%s", top.shadowChild)
				exprResult := alloc.Alloc().InitWithChildFromTemplate(rhino.Token_EXPR_RESULT, top.shadowChild, top.shadowTemplate)
				script := alloc.Alloc().InitWithChildFromTemplate(rhino.Token_SCRIPT, exprResult, top.shadowTemplate)
				shadowRoot := alloc.Alloc().InitWithChildFromTemplate(rhino.Token_ROOT, script, top.shadowTemplate)
				top.node.SetClosureUnawareShadow(shadowRoot)
				setOriginalNameIfPresent(top.node, internedStrings, top.origNamePtr)
			} else {
				contextStack = contextStack[:len(contextStack)-1]
			}
		}
	}

	for idx := uint32(0); idx < nodeCount; idx++ {
		popCompleted(idx)
		hdrOff := rd.need(3, "Truncated binary AstNode header")
		kindByte := b[hdrOff]
		flags := binary.LittleEndian.Uint16(b[hdrOff+1:])
		if flags&^binNodeKnownFlags != 0 {
			rd.fail("Invalid binary AstNode flags")
		}

		var stringValPtr, templateRawPtr, templateCookedPtr int32
		var doubleVal float64
		switch flags & binNodeValMask {
		case binNodeValNone:
		case binNodeValString:
			stringValPtr = rd.strPtr(stringPoolLen, "Truncated binary AstNode string value", "Binary AstNode string pointer out of bounds")
		case binNodeValDouble:
			doubleVal = math.Float64frombits(rd.u64("Truncated binary AstNode double value"))
		case binNodeValTemplate:
			rawU := rd.uvarint32("Truncated binary AstNode template value")
			cookedPlus1 := rd.uvarint32("Truncated binary AstNode template value")
			templateRawPtr = int32(rawU)
			templateCookedPtr = int32(cookedPlus1) - 1
			if rawU >= stringPoolLen || (cookedPlus1 != 0 && cookedPlus1-1 >= stringPoolLen) {
				rd.fail("Binary AstNode template pointer out of bounds")
			}
		}

		endIdx := idx + 1
		if flags&binNodeHasChildren != 0 {
			endIdx = rd.uvarint32("Truncated binary AstNode endIdx")
		}
		switch {
		case idx == 0 && endIdx != nodeCount:
			rd.fail("Binary script root endIdx mismatch")
		case idx > 0 && len(stack) == 0:
			rd.fail("Binary AstNode missing parent")
		case idx > 0 && stack[len(stack)-1].isShadowHost && endIdx != stack[len(stack)-1].endIdx:
			rd.fail("Binary shadow AstNode endIdx mismatch")
		case idx > 0 && !stack[len(stack)-1].isShadowHost && (endIdx <= idx || endIdx > stack[len(stack)-1].endIdx):
			rd.fail("Binary AstNode endIdx out of bounds")
		}

		relLine := rd.relDelta((flags&binNodeLineMask)>>binNodeLineShift, "Truncated binary AstNode relLine", "Invalid binary AstNode relLine mode")
		relCol := rd.relDelta((flags&binNodeColMask)>>binNodeColShift, "Truncated binary AstNode relCol", "Invalid binary AstNode relCol mode")

		var boolProps int64
		var typeVal, origNamePtr int32
		jsdocIdx := int32(-1)
		var sourceFile uint32
		hasType := (flags & binNodeHasType) != 0
		if flags&(binNodePropsMask|binNodeHasType|binNodeHasOrigName|binNodeHasJsdoc|binNodeHasSource) != 0 {
			switch (flags & binNodePropsMask) >> binNodePropsShift {
			case 0:
			case 1:
				boolProps = int64(rd.u16("Truncated binary AstNode boolProps"))
			case 2:
				boolProps = int64(rd.u32("Truncated binary AstNode boolProps"))
			case 3:
				boolProps = int64(rd.u64("Truncated binary AstNode boolProps"))
			}
			if boolProps < 0 {
				rd.fail("Invalid negative binary AstNode boolProps")
			}
			if hasType {
				typeVal = int32(rd.uvarint32("Truncated binary AstNode typePtr"))
			}
			if flags&binNodeHasOrigName != 0 {
				origNamePtr = rd.strPtr(stringPoolLen, "Truncated binary AstNode origName", "Binary AstNode origName pointer out of bounds")
			}
			if flags&binNodeHasJsdoc != 0 {
				u := rd.uvarint32("Truncated binary AstNode jsdocIdx")
				if u >= jsdocCount {
					rd.fail("Binary AstNode jsdocIdx out of bounds")
				}
				jsdocIdx = int32(u)
			}
			if flags&binNodeHasSource != 0 {
				sourceFile = rd.uvarint32("Truncated binary AstNode sourceFile")
				if sourceFile == 0 || int(sourceFile) > len(d.filePool) {
					rd.fail("Binary AstNode sourceFile out of bounds")
				}
			}
		}

		var curTemplate rhino.Node
		switch {
		case idx == 0:
			curTemplate = rootTmpl
		case stack[len(stack)-1].isShadowHost:
			curTemplate = stack[len(stack)-1].shadowTemplate
		default:
			curTemplate = stack[len(stack)-1].template
		}
		if curTemplate.IsNil() || sourceFile != 0 {
			if sourceFile == 0 {
				rd.fail("Missing SourceFile on binary AstNode")
			}
			curTemplate = d.CreateSourceInfoTemplate(d.filePool[sourceFile-1])
		}

		currentLine, currentColumn := previousLine+int(relLine), previousColumn+int(relCol)
		n := deserializeBinarySingleNode(
			alloc, NodeKind(kindByte), stringValPtr, doubleVal,
			templateRawPtr, templateCookedPtr, curTemplate, d.stringPool, internedStrings)
		if hasType && d.colorPoolShard != nil {
			if uint(int(typeVal)) >= uint(len(d.colorPoolShard.offsetToColor)) {
				rd.fail("Binary AstNode color pointer out of bounds")
			}
			n.SetColor(d.colorPoolShard.GetColor(int(typeVal)))
		}
		if boolProps > 0 {
			n.DeserializeProperties(d.FilterOutCastProp(boolProps), curTemplate.GetIsInClosureUnawareSubtree())
		}
		if jsdocIdx >= 0 {
			jd := &jsdocs[jsdocIdx]
			jsdoc := OptimizationJsdoc{
				kind_:                      jsdocKinds[jd.kindOff : jd.kindOff+jd.kindLen],
				licenseTextPointer_:        jd.licensePtr,
				descriptionPointer_:        jd.descPtr,
				meaningPointer_:            jd.meaningPtr,
				alternateMessageIdPointer_: jd.altMsgPtr,
				closureUnawareModePointer_: jd.cuModePtr,
			}
			n.SetJSDocInfo(jsDocSerializer_deserializeJsdocWithBuilder(&jsdoc, d.stringPool, &jsdocBuilder))
		}

		n.SetLinenoCharno(currentLine, currentColumn)
		previousLine, previousColumn = currentLine, currentColumn

		isShadowHost := rhino.Node_HasBitSet(boolProps, int(rhino.NodeProperty_CLOSURE_UNAWARE_SHADOW.GetNumber()))
		if idx == 0 {
			scriptNode = n
		} else {
			parent := &stack[len(stack)-1]
			if parent.isShadowHost {
				parent.shadowChild = n
			} else {
				parent.node.AddChildToBack(n)
			}
			if !isShadowHost {
				setOriginalNameIfPresent(n, internedStrings, origNamePtr)
			}
		}

		if isShadowHost {
			contextStack = append(contextStack, featureCollector.VisitSingleNode(contextStack, n))
			if endIdx <= idx+1 {
				panic(jutil.NewIndexOutOfBoundsException("Index 0 out of bounds for length 0"))
			}
			shadowTmpl := alloc.Alloc().InitFromTemplate(curTemplate.GetToken(), curTemplate)
			shadowTmpl.SetIsInClosureUnawareSubtree(true)
			var shadowOrigName int32
			if idx > 0 {
				shadowOrigName = origNamePtr
			}
			stack = append(stack, binaryParentFrame{
				node: n, template: curTemplate, shadowTemplate: shadowTmpl,
				endIdx: endIdx, origNamePtr: shadowOrigName, isShadowHost: true,
			})
		} else if endIdx > idx+1 {
			contextStack = append(contextStack, featureCollector.VisitSingleNode(contextStack, n))
			stack = append(stack, binaryParentFrame{node: n, template: curTemplate, endIdx: endIdx})
		} else {
			_ = featureCollector.VisitSingleNode(contextStack, n)
		}
	}

	if rd.pos != len(b) {
		rd.fail("Trailing bytes in binary script")
	}
	popCompleted(nodeCount)
	scriptNode.PutProp(rhino.Node_FEATURE_SET, featureCollector.AllFeatures())
	return scriptNode
}

func deserializeBinarySingleNode(
	alloc *rhino.NodeAllocator, kind NodeKind, stringValPtr int32, doubleVal float64,
	templateRawPtr, templateCookedPtr int32, templateNode rhino.Node,
	stringPool *StringPool, internedStrings *rhino.RhinoStringPool_LazyInternedStringList,
) rhino.Node {
	if uint(kind) < uint(len(nodeKindToDeserSpec)) {
		spec := nodeKindToDeserSpec[kind]
		switch spec.mode {
		case deserModePlain:
			return alloc.Alloc().InitFromTemplate(spec.token, templateNode)
		case deserModeString:
			return alloc.Alloc().InitFromTemplate(spec.token, templateNode).SetStringID(internedStrings.GetID(int(stringValPtr)))
		case deserModeQuotedString:
			n := alloc.Alloc().InitFromTemplate(spec.token, templateNode).SetStringID(internedStrings.GetID(int(stringValPtr)))
			n.SetQuotedStringKey()
			return n
		case deserModePostIncrDecr:
			n := alloc.Alloc().InitFromTemplate(spec.token, templateNode)
			n.PutBooleanProp(rhino.Node_INCRDECR_PROP, true)
			return n
		case deserModeNumber:
			return alloc.Alloc().InitFromTemplate(rhino.Token_NUMBER, templateNode).SetNumberRaw(doubleVal)
		case deserModeBigInt:
			return alloc.Alloc().InitBigIntFromTemplate(jutil.NewBigIntegerFromString(stringPool.Get(int(stringValPtr))), templateNode)
		case deserModeTemplateLitString:
			return alloc.Alloc().InitTemplateLitStringFromStringPool(
				internedStrings, int(templateCookedPtr), int(templateRawPtr), templateNode)
		}
	}
	panic(jutil.NewIOException("Invalid binary AstNode kind"))
}
