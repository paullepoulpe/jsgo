// Port of com/google/javascript/jscomp/serialization/TypedAstSerializer.java @ closure-compiler v20261021

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
	"math"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// flatAstNode is a pointer-free (noscan) pre-order record of a single AstNode
// used to compute exact bottom-up wire sizes and emit protobuf bytes directly
// without materializing an intermediate *AstNode proto tree.
type flatAstNode struct {
	boolProps  uint64
	valBits    uint64
	endIdx     uint32
	wireSize   uint32
	relLine    int32
	relCol     int32
	origName   int32
	typePtr    int32
	sourceFile int32
	jsdocIdx   int32
	kind       uint16
	valCase    uint8
	hasType    bool
}

// flatJsdoc is a pointer-free (noscan) record of an OptimizationJsdoc message.
type flatJsdoc struct {
	kindOff    uint32
	kindLen    uint32
	licensePtr int32
	descPtr    int32
	meaningPtr int32
	altMsgPtr  int32
	cuModePtr  int32
	wireSize   uint32
}

type fixupTarget uint8

const (
	fixupNodeStrVal fixupTarget = iota
	fixupNodeStrValUsedProp
	fixupNodeRawTmpl
	fixupNodeCookedTmpl
	fixupNodeOrigName
	fixupNodeColor
	fixupNodeSourceFile
	fixupLazySourceFile
	fixupJsdocLicense
	fixupJsdocDesc
	fixupJsdocAltMsg
	fixupJsdocMeaning
	fixupJsdocCuMode
)

// poolFixup is a 12-byte pointer-free (noscan) deferred pool registration.
type poolFixup struct {
	targetIdx uint32
	refIdx    uint32
	target    fixupTarget
}

type scriptWireEncoder struct {
	nodes              []flatAstNode
	jsdocs             []flatJsdoc
	jsdocKinds         []int32
	fixups             []poolFixup
	colorRefs          []*colors.Color
	sfRefs             []*SourceFile
	subtreeSourceFiles []*SourceFile
	previousLine       int
	previousColumn     int
	lazyAst            LazyAst
}

func (enc *scriptWireEncoder) reset() {
	enc.nodes = enc.nodes[:0]
	enc.jsdocs = enc.jsdocs[:0]
	enc.jsdocKinds = enc.jsdocKinds[:0]
	enc.fixups = enc.fixups[:0]
	clear(enc.colorRefs)
	enc.colorRefs = enc.colorRefs[:0]
	clear(enc.sfRefs)
	enc.sfRefs = enc.sfRefs[:0]
	clear(enc.subtreeSourceFiles)
	enc.subtreeSourceFiles = enc.subtreeSourceFiles[:0]
	enc.previousLine = 0
	enc.previousColumn = 0
	enc.lazyAst = LazyAst{}
}

func (enc *scriptWireEncoder) peekLastSubtreeSourceFile() *SourceFile {
	if len(enc.subtreeSourceFiles) == 0 {
		return nil
	}
	return enc.subtreeSourceFiles[len(enc.subtreeSourceFiles)-1]
}

func (s *TypedAstSerializer) visitDirect(n rhino.Node) {
	jutil.CheckState(n.GetLength() >= 0)
	currentLine := n.GetLineno()
	currentColumn := n.GetCharno()
	fn := flatAstNode{
		relLine:  int32(currentLine - s.previousLine),
		relCol:   int32(currentColumn - s.previousColumn),
		jsdocIdx: -1,
	}
	s.previousLine = currentLine
	s.previousColumn = currentColumn

	if s.colorTypeSerializer != nil {
		if color := n.GetColor(); color != nil {
			fn.typePtr = int32(s.colorTypeSerializer.colorSerializer.AddColor(color))
			fn.hasType = true
		}
	} else if s.jsTypeSerializer != nil {
		s.jsTypeSerializer.addTypeToFlatNode(n, &fn)
	}

	if jsdoc := n.GetJSDocInfo(); jsdoc != nil {
		kindOff := len(s.scriptEnc.jsdocKinds)
		s.scriptEnc.jsdocKinds = jsDocSerializer_fillJsdocTagsOnly(jsdoc, s.scriptEnc.jsdocKinds)
		kindLen := len(s.scriptEnc.jsdocKinds) - kindOff
		lic, desc, alt, mean, mode, hasMode := jsDocSerializer_stringFields(jsdoc)
		var doc flatJsdoc
		if lic != "" {
			doc.licensePtr = int32(s.stringPool.Put(lic))
		}
		if desc != "" {
			doc.descPtr = int32(s.stringPool.Put(desc))
		}
		if alt != "" {
			doc.altMsgPtr = int32(s.stringPool.Put(alt))
		}
		if mean != "" {
			doc.meaningPtr = int32(s.stringPool.Put(mean))
		}
		if hasMode {
			doc.cuModePtr = int32(s.stringPool.Put(mode.String()))
		}
		if kindLen != 0 || doc.licensePtr != 0 || doc.descPtr != 0 || doc.meaningPtr != 0 || doc.altMsgPtr != 0 || doc.cuModePtr != 0 {
			doc.kindOff = uint32(kindOff)
			doc.kindLen = uint32(kindLen)
			fn.jsdocIdx = int32(len(s.scriptEnc.jsdocs))
			s.scriptEnc.jsdocs = append(s.scriptEnc.jsdocs, doc)
		}
	}

	fn.kind = uint16(s.KindTranslator(n).GetNumber())
	s.translateValueDirect(&fn, n)
	fn.boolProps = uint64(n.SerializeProperties())
	sourceFile := s.GetSourceFilePointer(n)
	fn.sourceFile = int32(sourceFile)

	idx := len(s.scriptEnc.nodes)
	s.scriptEnc.nodes = append(s.scriptEnc.nodes, fn)

	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		s.visitDirect(child)
	}
	if shadowedCode := n.GetClosureUnawareShadow(); !shadowedCode.IsNil() {
		// For Closure shadow hosts, the ASTNode will get a boolean property bit set
		// (CLOSURE_UNAWARE_SHADOW) that indicates this child ASTNode is not a normal child, but is
		// the contents of the shadow.
		// Shadow roots are structured as
		// ROOT -> SCRIPT -> EXPR_RESULT -> FUNCTION
		// We avoid including the ROOT -> SCRIPT -> EXPR_RESULT structure in the serialization as it
		// is all synthetic code that would unnecessarily bloat the TypedAST and is instead recreated
		// upon deserialization.
		// The child ASTNode is just the FUNCTION.
		s.visitDirect(shadowedCode.GetFirstFirstChild().GetFirstChild())
	}
	s.scriptEnc.nodes[idx].endIdx = uint32(len(s.scriptEnc.nodes))

	if sourceFile != 0 {
		// ArrayDeque.removeLast() throws NoSuchElementException on an empty deque.
		if len(s.subtreeSourceFiles) == 0 {
			panic(jutil.NewNoSuchElementException())
		}
		s.subtreeSourceFiles = s.subtreeSourceFiles[:len(s.subtreeSourceFiles)-1]
	}
	if origID := n.GetOriginalNameID(); origID != 0 {
		s.scriptEnc.nodes[idx].origName = int32(s.stringPool.PutID(origID))
	}
}

func (s *TypedAstSerializer) translateValueDirect(fn *flatAstNode, n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_GETPROP,
		rhino.Token_OPTCHAIN_GETPROP:
		ptr := s.stringPool.PutID(n.GetStringID())
		if s.recordUsedProps {
			s.markUsedProp(ptr)
		}
		fn.valBits = uint64(uint32(int32(ptr)))
		fn.valCase = 3
	case rhino.Token_STRING_KEY,
		rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_MEMBER_FIELD_DEF,
		rhino.Token_GETTER_DEF,
		rhino.Token_SETTER_DEF:
		ptr := s.stringPool.PutID(n.GetStringID())
		if s.recordUsedProps && !n.IsQuotedStringKey() {
			s.markUsedProp(ptr)
		}
		fn.valBits = uint64(uint32(int32(ptr)))
		fn.valCase = 3
	case rhino.Token_NAME,
		rhino.Token_STRINGLIT,
		rhino.Token_LABEL_NAME,
		rhino.Token_IMPORT_STAR:
		fn.valBits = uint64(uint32(int32(s.stringPool.PutID(n.GetStringID()))))
		fn.valCase = 3
	case rhino.Token_TEMPLATELIT_STRING:
		// The raw string is put in the pool before the cooked one (Java's
		// evaluation order).
		rawPtr := uint32(int32(s.stringPool.PutID(n.GetRawStringID())))
		cookedPtr := uint32(0xFFFFFFFF) // int32(-1)
		if cookedID := n.GetCookedStringID(); cookedID != 0 {
			cookedPtr = uint32(int32(s.stringPool.PutID(cookedID)))
		}
		fn.valBits = (uint64(rawPtr) << 32) | uint64(cookedPtr)
		fn.valCase = 8
	case rhino.Token_NUMBER:
		fn.valBits = math.Float64bits(n.GetDouble())
		fn.valCase = 4
	case rhino.Token_BIGINT:
		fn.valBits = uint64(uint32(int32(s.stringPool.Put(n.GetBigInt().String()))))
		fn.valCase = 3
	}
}

func (enc *scriptWireEncoder) buildScriptPassA(script rhino.Node, parent *TypedAstSerializer) {
	jutil.CheckState(script.IsScript())
	enc.reset()

	var scriptSF *SourceFile
	if ssf := script.GetStaticSourceFile(); ssf != nil {
		scriptSF = ssf.(*SourceFile)
	}
	if scriptSF == nil {
		if !enc.peekLastSubtreeSourceFile().IsExtern() {
			jutil.CheckStateMsg2(
				false,
				"Unexpected null SourceFile for node %s with parent %s",
				script.ToStringTree(),
				script.GetParent())
		}
	} else {
		enc.subtreeSourceFiles = append(enc.subtreeSourceFiles, scriptSF)
		sfRef := uint32(len(enc.sfRefs))
		enc.sfRefs = append(enc.sfRefs, scriptSF)
		enc.fixups = append(enc.fixups, poolFixup{
			refIdx: sfRef,
			target: fixupLazySourceFile,
		})
	}

	enc.visitParallel(script, parent)
	clear(enc.subtreeSourceFiles)
	enc.subtreeSourceFiles = enc.subtreeSourceFiles[:0]

	if encodedSourceMap := parent.compiler.GetBase64SourceMapContents(*script.GetSourceFileName()); encodedSourceMap != nil {
		enc.lazyAst.sourceMappingUrl_ = *encodedSourceMap
	}
}

func (enc *scriptWireEncoder) visitParallel(n rhino.Node, parent *TypedAstSerializer) {
	jutil.CheckState(n.GetLength() >= 0)
	nodeIdx := uint32(len(enc.nodes))
	currentLine := n.GetLineno()
	currentColumn := n.GetCharno()
	fn := flatAstNode{
		relLine:  int32(currentLine - enc.previousLine),
		relCol:   int32(currentColumn - enc.previousColumn),
		jsdocIdx: -1,
	}
	enc.previousLine = currentLine
	enc.previousColumn = currentColumn

	if parent.colorTypeSerializer != nil {
		if color := n.GetColor(); color != nil {
			fn.hasType = true
			cRef := uint32(len(enc.colorRefs))
			enc.colorRefs = append(enc.colorRefs, color)
			enc.fixups = append(enc.fixups, poolFixup{
				targetIdx: nodeIdx,
				refIdx:    cRef,
				target:    fixupNodeColor,
			})
		}
	} else if parent.jsTypeSerializer != nil {
		parent.jsTypeSerializer.addTypeToFlatNode(n, &fn)
	}

	if jsdoc := n.GetJSDocInfo(); jsdoc != nil {
		kindOff := len(enc.jsdocKinds)
		enc.jsdocKinds = jsDocSerializer_fillJsdocTagsOnly(jsdoc, enc.jsdocKinds)
		kindLen := len(enc.jsdocKinds) - kindOff

		lic, desc, alt, mean, mode, hasMode := jsDocSerializer_stringFields(jsdoc)
		hasStr := lic != "" || desc != "" || alt != "" || mean != "" || hasMode
		if kindLen != 0 || hasStr {
			jIdx := uint32(len(enc.jsdocs))
			enc.jsdocs = append(enc.jsdocs, flatJsdoc{
				kindOff: uint32(kindOff),
				kindLen: uint32(kindLen),
			})
			fn.jsdocIdx = int32(jIdx)
			if hasStr {
				if lic != "" {
					enc.fixups = append(enc.fixups, poolFixup{
						targetIdx: jIdx,
						refIdx:    rhino.RhinoStringPool_InternID(lic),
						target:    fixupJsdocLicense,
					})
				}
				if desc != "" {
					enc.fixups = append(enc.fixups, poolFixup{
						targetIdx: jIdx,
						refIdx:    rhino.RhinoStringPool_InternID(desc),
						target:    fixupJsdocDesc,
					})
				}
				if alt != "" {
					enc.fixups = append(enc.fixups, poolFixup{
						targetIdx: jIdx,
						refIdx:    rhino.RhinoStringPool_InternID(alt),
						target:    fixupJsdocAltMsg,
					})
				}
				if mean != "" {
					enc.fixups = append(enc.fixups, poolFixup{
						targetIdx: jIdx,
						refIdx:    rhino.RhinoStringPool_InternID(mean),
						target:    fixupJsdocMeaning,
					})
				}
				if hasMode {
					enc.fixups = append(enc.fixups, poolFixup{
						targetIdx: jIdx,
						refIdx:    rhino.RhinoStringPool_InternID(mode.String()),
						target:    fixupJsdocCuMode,
					})
				}
			}
		}
	}

	fn.kind = uint16(parent.KindTranslator(n).GetNumber())
	enc.translateValueParallel(&fn, nodeIdx, n, parent.recordUsedProps)
	fn.boolProps = uint64(n.SerializeProperties())

	pushedSF := false
	var sf *SourceFile
	if ssf := n.GetStaticSourceFile(); ssf != nil {
		sf = ssf.(*SourceFile)
	}
	if sf == nil {
		if !enc.peekLastSubtreeSourceFile().IsExtern() {
			jutil.CheckStateMsg2(
				false,
				"Unexpected null SourceFile for node %s with parent %s",
				n.ToStringTree(),
				n.GetParent())
		}
	} else if sf != enc.peekLastSubtreeSourceFile() {
		enc.subtreeSourceFiles = append(enc.subtreeSourceFiles, sf)
		pushedSF = true
		sfRef := uint32(len(enc.sfRefs))
		enc.sfRefs = append(enc.sfRefs, sf)
		enc.fixups = append(enc.fixups, poolFixup{
			targetIdx: nodeIdx,
			refIdx:    sfRef,
			target:    fixupNodeSourceFile,
		})
	}

	enc.nodes = append(enc.nodes, fn)

	for child := n.GetFirstChild(); !child.IsNil(); child = child.GetNext() {
		enc.visitParallel(child, parent)
	}
	if shadowedCode := n.GetClosureUnawareShadow(); !shadowedCode.IsNil() {
		enc.visitParallel(shadowedCode.GetFirstFirstChild().GetFirstChild(), parent)
	}
	enc.nodes[nodeIdx].endIdx = uint32(len(enc.nodes))

	if pushedSF {
		if len(enc.subtreeSourceFiles) == 0 {
			panic(jutil.NewNoSuchElementException())
		}
		enc.subtreeSourceFiles = enc.subtreeSourceFiles[:len(enc.subtreeSourceFiles)-1]
	}
	if origID := n.GetOriginalNameID(); origID != 0 {
		enc.fixups = append(enc.fixups, poolFixup{
			targetIdx: nodeIdx,
			refIdx:    origID,
			target:    fixupNodeOrigName,
		})
	}
}

func (enc *scriptWireEncoder) translateValueParallel(fn *flatAstNode, nodeIdx uint32, n rhino.Node, recordUsedProps bool) {
	switch n.GetToken() {
	case rhino.Token_GETPROP,
		rhino.Token_OPTCHAIN_GETPROP:
		target := fixupNodeStrVal
		if recordUsedProps {
			target = fixupNodeStrValUsedProp
		}
		enc.fixups = append(enc.fixups, poolFixup{
			targetIdx: nodeIdx,
			refIdx:    n.GetStringID(),
			target:    target,
		})
		fn.valCase = 3
	case rhino.Token_STRING_KEY,
		rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_MEMBER_FIELD_DEF,
		rhino.Token_GETTER_DEF,
		rhino.Token_SETTER_DEF:
		target := fixupNodeStrVal
		if recordUsedProps && !n.IsQuotedStringKey() {
			target = fixupNodeStrValUsedProp
		}
		enc.fixups = append(enc.fixups, poolFixup{
			targetIdx: nodeIdx,
			refIdx:    n.GetStringID(),
			target:    target,
		})
		fn.valCase = 3
	case rhino.Token_NAME,
		rhino.Token_STRINGLIT,
		rhino.Token_LABEL_NAME,
		rhino.Token_IMPORT_STAR:
		enc.fixups = append(enc.fixups, poolFixup{
			targetIdx: nodeIdx,
			refIdx:    n.GetStringID(),
			target:    fixupNodeStrVal,
		})
		fn.valCase = 3
	case rhino.Token_TEMPLATELIT_STRING:
		fn.valBits = 0xFFFFFFFF // cooked = int32(-1) by default
		fn.valCase = 8
		enc.fixups = append(enc.fixups, poolFixup{
			targetIdx: nodeIdx,
			refIdx:    n.GetRawStringID(),
			target:    fixupNodeRawTmpl,
		})
		if cookedID := n.GetCookedStringID(); cookedID != 0 {
			enc.fixups = append(enc.fixups, poolFixup{
				targetIdx: nodeIdx,
				refIdx:    cookedID,
				target:    fixupNodeCookedTmpl,
			})
		}
	case rhino.Token_NUMBER:
		fn.valBits = math.Float64bits(n.GetDouble())
		fn.valCase = 4
	case rhino.Token_BIGINT:
		enc.fixups = append(enc.fixups, poolFixup{
			targetIdx: nodeIdx,
			refIdx:    rhino.RhinoStringPool_InternID(n.GetBigInt().String()),
			target:    fixupNodeStrVal,
		})
		fn.valCase = 3
	}
}

func (enc *scriptWireEncoder) marshalScriptPassC(byteSlab *[]byte) {
	scriptBytes := enc.encodeToByteSlab(byteSlab)
	enc.lazyAst.script_ = scriptBytes
	if len(enc.lazyAst.sourceMappingUrl_) == 0 {
		enc.lazyAst.cachedSize = lazyAstFastSize(scriptBytes, enc.lazyAst.sourceFile_) + 1
	} else {
		enc.lazyAst.serializedSize()
	}
}

func (enc *scriptWireEncoder) encodeToByteSlab(byteSlab *[]byte) []byte {
	jsdocs := enc.jsdocs
	for i := range jsdocs {
		jsdocs[i].wireSize = enc.computeJsdocWireSize(&jsdocs[i])
	}
	nodes := enc.nodes
	for i := len(nodes) - 1; i >= 0; i-- {
		nodes[i].wireSize = enc.computeNodeWireSize(&nodes[i], i)
	}
	sz := int(nodes[0].wireSize)
	if sz > 262144 {
		out := enc.appendNode(make([]byte, 0, sz), 0)
		if len(out) != sz {
			panic("AstNode wire size mismatch")
		}
		return out
	}
	slab := *byteSlab
	if len(slab) < sz {
		slab = make([]byte, max(4*1024*1024, sz))
	}
	buf := slab[:0:sz]
	*byteSlab = slab[sz:]
	out := enc.appendNode(buf, 0)
	if len(out) != sz {
		panic("AstNode wire size mismatch")
	}
	return out
}

func (enc *scriptWireEncoder) computeJsdocWireSize(d *flatJsdoc) uint32 {
	size := 0
	if d.kindLen > 0 {
		kinds := enc.jsdocKinds[d.kindOff : d.kindOff+d.kindLen]
		packedLen := 0
		for _, k := range kinds {
			packedLen += protoVarintLen(uint64(int64(k)))
		}
		size += 1 + protoVarintLen(uint64(packedLen)) + packedLen
	}
	if d.licensePtr != 0 {
		size += 1 + protoVarintLen(uint64(int64(d.licensePtr)))
	}
	if d.descPtr != 0 {
		size += 1 + protoVarintLen(uint64(int64(d.descPtr)))
	}
	if d.meaningPtr != 0 {
		size += 1 + protoVarintLen(uint64(int64(d.meaningPtr)))
	}
	if d.altMsgPtr != 0 {
		size += 1 + protoVarintLen(uint64(int64(d.altMsgPtr)))
	}
	if d.cuModePtr != 0 {
		size += 1 + protoVarintLen(uint64(int64(d.cuModePtr)))
	}
	return uint32(size)
}

func templateStringWireSize(valBits uint64) int {
	rawPtr := int32(valBits >> 32)
	cookedPtr := int32(valBits)
	size := 0
	if rawPtr != 0 {
		size += 1 + protoVarintLen(uint64(uint32(rawPtr)))
	}
	if cookedPtr != 0 {
		size += 1 + protoVarintLen(uint64(int64(cookedPtr)))
	}
	return size
}

func (enc *scriptWireEncoder) computeNodeWireSize(fn *flatAstNode, idx int) uint32 {
	size := 0
	if fn.kind != 0 {
		size += 1 + protoVarintLen(uint64(fn.kind))
	}
	nodes := enc.nodes
	for c := idx + 1; c < int(fn.endIdx); c = int(nodes[c].endIdx) {
		cs := int(nodes[c].wireSize)
		size += 1 + protoVarintLen(uint64(cs)) + cs
	}
	switch fn.valCase {
	case 3:
		size += 1 + protoVarintLen(fn.valBits)
	case 4:
		size += 9
	case 8:
		ts := templateStringWireSize(fn.valBits)
		size += 1 + protoVarintLen(uint64(ts)) + ts
	}
	if fn.relLine != 0 {
		size += 1 + protoVarintLen(uint64(uint32((fn.relLine<<1)^(fn.relLine>>31))))
	}
	if fn.relCol != 0 {
		size += 1 + protoVarintLen(uint64(uint32((fn.relCol<<1)^(fn.relCol>>31))))
	}
	if fn.jsdocIdx >= 0 {
		js := int(enc.jsdocs[fn.jsdocIdx].wireSize)
		size += 1 + protoVarintLen(uint64(js)) + js
	}
	if fn.origName != 0 {
		size += 1 + protoVarintLen(uint64(uint32(fn.origName)))
	}
	if fn.boolProps != 0 {
		size += 1 + protoVarintLen(fn.boolProps)
	}
	if fn.hasType {
		size += 1 + protoVarintLen(uint64(int64(fn.typePtr)))
	}
	if fn.sourceFile != 0 {
		size += 2 + protoVarintLen(uint64(uint32(fn.sourceFile)))
	}
	return uint32(size)
}

func (enc *scriptWireEncoder) appendJsdoc(b []byte, d *flatJsdoc) []byte {
	if d.kindLen > 0 {
		kinds := enc.jsdocKinds[d.kindOff : d.kindOff+d.kindLen]
		packedLen := 0
		for _, k := range kinds {
			packedLen += protoVarintLen(uint64(int64(k)))
		}
		b = append(b, 10)
		b = protoAppendVarint(b, uint64(packedLen))
		for _, k := range kinds {
			b = protoAppendInt32(b, k)
		}
	}
	if d.licensePtr != 0 {
		b = append(b, 16)
		b = protoAppendInt32(b, d.licensePtr)
	}
	if d.descPtr != 0 {
		b = append(b, 24)
		b = protoAppendInt32(b, d.descPtr)
	}
	if d.meaningPtr != 0 {
		b = append(b, 32)
		b = protoAppendInt32(b, d.meaningPtr)
	}
	if d.altMsgPtr != 0 {
		b = append(b, 40)
		b = protoAppendInt32(b, d.altMsgPtr)
	}
	if d.cuModePtr != 0 {
		b = append(b, 48)
		b = protoAppendInt32(b, d.cuModePtr)
	}
	return b
}

func (enc *scriptWireEncoder) appendNode(b []byte, idx int) []byte {
	fn := enc.nodes[idx]
	if fn.kind != 0 {
		b = append(b, 8)
		b = protoAppendVarint(b, uint64(fn.kind))
	}
	for c := idx + 1; c < int(fn.endIdx); c = int(enc.nodes[c].endIdx) {
		b = append(b, 18)
		b = protoAppendVarint(b, uint64(enc.nodes[c].wireSize))
		b = enc.appendNode(b, c)
	}
	switch fn.valCase {
	case 3:
		b = append(b, 24)
		b = protoAppendVarint(b, fn.valBits)
	case 4:
		b = append(b, 33)
		b = protoAppendFixed64(b, fn.valBits)
	}
	if fn.relLine != 0 {
		b = append(b, 40)
		b = protoAppendSint32(b, fn.relLine)
	}
	if fn.relCol != 0 {
		b = append(b, 48)
		b = protoAppendSint32(b, fn.relCol)
	}
	if fn.jsdocIdx >= 0 {
		d := &enc.jsdocs[fn.jsdocIdx]
		b = append(b, 58)
		b = protoAppendVarint(b, uint64(d.wireSize))
		b = enc.appendJsdoc(b, d)
	}
	if fn.valCase == 8 {
		rawPtr := int32(fn.valBits >> 32)
		cookedPtr := int32(fn.valBits)
		ts := templateStringWireSize(fn.valBits)
		b = append(b, 66)
		b = protoAppendVarint(b, uint64(ts))
		if rawPtr != 0 {
			b = append(b, 8)
			b = protoAppendUint32(b, rawPtr)
		}
		if cookedPtr != 0 {
			b = append(b, 16)
			b = protoAppendInt32(b, cookedPtr)
		}
	}
	if fn.origName != 0 {
		b = append(b, 72)
		b = protoAppendUint32(b, fn.origName)
	}
	if fn.boolProps != 0 {
		b = append(b, 112)
		b = protoAppendVarint(b, fn.boolProps)
	}
	if fn.hasType {
		b = append(b, 120)
		b = protoAppendInt32(b, fn.typePtr)
	}
	if fn.sourceFile != 0 {
		b = append(b, 0x80, 0x01)
		b = protoAppendUint32(b, fn.sourceFile)
	}
	return b
}

type nodeDeserMode uint8

const (
	deserModeInvalid nodeDeserMode = iota
	deserModePlain
	deserModeString
	deserModeQuotedString
	deserModePostIncrDecr
	deserModeNumber
	deserModeBigInt
	deserModeTemplateLitString
)

type nodeKindDeserSpec struct {
	token rhino.Token
	mode  nodeDeserMode
}

type tokenNodeKindEntry struct {
	token rhino.Token
	kind  NodeKind
	mode  nodeDeserMode
}

const nodeKindTokenConditional NodeKind = -1

var (
	tokenToNodeKind     [256]NodeKind
	nodeKindToDeserSpec [int(NodeKind_UNRECOGNIZED) + 1]nodeKindDeserSpec
)

var bijectiveTokenNodeKinds = [...]tokenNodeKindEntry{
	{rhino.Token_SCRIPT, NodeKind_SOURCE_FILE, deserModePlain},
	{rhino.Token_NUMBER, NodeKind_NUMBER_LITERAL, deserModeNumber},
	{rhino.Token_STRINGLIT, NodeKind_STRING_LITERAL, deserModeString},
	{rhino.Token_BIGINT, NodeKind_BIGINT_LITERAL, deserModeBigInt},
	{rhino.Token_REGEXP, NodeKind_REGEX_LITERAL, deserModePlain},
	{rhino.Token_FALSE, NodeKind_FALSE, deserModePlain},
	{rhino.Token_TRUE, NodeKind_TRUE, deserModePlain},
	{rhino.Token_NULL, NodeKind_NULL, deserModePlain},
	{rhino.Token_THIS, NodeKind_THIS, deserModePlain},
	{rhino.Token_VOID, NodeKind_VOID, deserModePlain},
	{rhino.Token_ARRAYLIT, NodeKind_ARRAY_LITERAL, deserModePlain},
	{rhino.Token_OBJECTLIT, NodeKind_OBJECT_LITERAL, deserModePlain},
	{rhino.Token_NAME, NodeKind_IDENTIFIER, deserModeString},
	{rhino.Token_ASSIGN, NodeKind_ASSIGNMENT, deserModePlain},
	{rhino.Token_CALL, NodeKind_CALL, deserModePlain},
	{rhino.Token_NEW, NodeKind_NEW, deserModePlain},
	{rhino.Token_GETPROP, NodeKind_PROPERTY_ACCESS, deserModeString},
	{rhino.Token_GETELEM, NodeKind_ELEMENT_ACCESS, deserModePlain},
	{rhino.Token_COMMA, NodeKind_COMMA, deserModePlain},
	{rhino.Token_OR, NodeKind_BOOLEAN_OR, deserModePlain},
	{rhino.Token_AND, NodeKind_BOOLEAN_AND, deserModePlain},
	{rhino.Token_HOOK, NodeKind_HOOK, deserModePlain},
	{rhino.Token_EQ, NodeKind_EQUAL, deserModePlain},
	{rhino.Token_NE, NodeKind_NOT_EQUAL, deserModePlain},
	{rhino.Token_LT, NodeKind_LESS_THAN, deserModePlain},
	{rhino.Token_LE, NodeKind_LESS_THAN_EQUAL, deserModePlain},
	{rhino.Token_GT, NodeKind_GREATER_THAN, deserModePlain},
	{rhino.Token_GE, NodeKind_GREATER_THAN_EQUAL, deserModePlain},
	{rhino.Token_SHEQ, NodeKind_TRIPLE_EQUAL, deserModePlain},
	{rhino.Token_SHNE, NodeKind_NOT_TRIPLE_EQUAL, deserModePlain},
	{rhino.Token_NOT, NodeKind_NOT, deserModePlain},
	{rhino.Token_POS, NodeKind_POSITIVE, deserModePlain},
	{rhino.Token_NEG, NodeKind_NEGATIVE, deserModePlain},
	{rhino.Token_TYPEOF, NodeKind_TYPEOF, deserModePlain},
	{rhino.Token_INSTANCEOF, NodeKind_INSTANCEOF, deserModePlain},
	{rhino.Token_IN, NodeKind_IN, deserModePlain},
	{rhino.Token_ADD, NodeKind_ADD, deserModePlain},
	{rhino.Token_SUB, NodeKind_SUBTRACT, deserModePlain},
	{rhino.Token_MUL, NodeKind_MULTIPLY, deserModePlain},
	{rhino.Token_DIV, NodeKind_DIVIDE, deserModePlain},
	{rhino.Token_MOD, NodeKind_MODULO, deserModePlain},
	{rhino.Token_EXPONENT, NodeKind_EXPONENT, deserModePlain},
	{rhino.Token_BITNOT, NodeKind_BITWISE_NOT, deserModePlain},
	{rhino.Token_BITOR, NodeKind_BITWISE_OR, deserModePlain},
	{rhino.Token_BITAND, NodeKind_BITWISE_AND, deserModePlain},
	{rhino.Token_BITXOR, NodeKind_BITWISE_XOR, deserModePlain},
	{rhino.Token_LSH, NodeKind_LEFT_SHIFT, deserModePlain},
	{rhino.Token_RSH, NodeKind_RIGHT_SHIFT, deserModePlain},
	{rhino.Token_URSH, NodeKind_UNSIGNED_RIGHT_SHIFT, deserModePlain},
	{rhino.Token_ASSIGN_ADD, NodeKind_ASSIGN_ADD, deserModePlain},
	{rhino.Token_ASSIGN_SUB, NodeKind_ASSIGN_SUBTRACT, deserModePlain},
	{rhino.Token_ASSIGN_MUL, NodeKind_ASSIGN_MULTIPLY, deserModePlain},
	{rhino.Token_ASSIGN_DIV, NodeKind_ASSIGN_DIVIDE, deserModePlain},
	{rhino.Token_ASSIGN_MOD, NodeKind_ASSIGN_MODULO, deserModePlain},
	{rhino.Token_ASSIGN_EXPONENT, NodeKind_ASSIGN_EXPONENT, deserModePlain},
	{rhino.Token_ASSIGN_BITOR, NodeKind_ASSIGN_BITWISE_OR, deserModePlain},
	{rhino.Token_ASSIGN_BITAND, NodeKind_ASSIGN_BITWISE_AND, deserModePlain},
	{rhino.Token_ASSIGN_BITXOR, NodeKind_ASSIGN_BITWISE_XOR, deserModePlain},
	{rhino.Token_ASSIGN_LSH, NodeKind_ASSIGN_LEFT_SHIFT, deserModePlain},
	{rhino.Token_ASSIGN_RSH, NodeKind_ASSIGN_RIGHT_SHIFT, deserModePlain},
	{rhino.Token_ASSIGN_URSH, NodeKind_ASSIGN_UNSIGNED_RIGHT_SHIFT, deserModePlain},
	{rhino.Token_YIELD, NodeKind_YIELD, deserModePlain},
	{rhino.Token_AWAIT, NodeKind_AWAIT, deserModePlain},
	{rhino.Token_DELPROP, NodeKind_DELETE, deserModePlain},
	{rhino.Token_TAGGED_TEMPLATELIT, NodeKind_TAGGED_TEMPLATELIT, deserModePlain},
	{rhino.Token_TEMPLATELIT, NodeKind_TEMPLATELIT, deserModePlain},
	{rhino.Token_TEMPLATELIT_SUB, NodeKind_TEMPLATELIT_SUB, deserModePlain},
	{rhino.Token_TEMPLATELIT_STRING, NodeKind_TEMPLATELIT_STRING, deserModeTemplateLitString},
	{rhino.Token_NEW_TARGET, NodeKind_NEW_TARGET, deserModePlain},
	{rhino.Token_COMPUTED_PROP, NodeKind_COMPUTED_PROP, deserModePlain},
	{rhino.Token_IMPORT_META, NodeKind_IMPORT_META, deserModePlain},
	{rhino.Token_OPTCHAIN_GETPROP, NodeKind_OPTCHAIN_PROPERTY_ACCESS, deserModeString},
	{rhino.Token_OPTCHAIN_CALL, NodeKind_OPTCHAIN_CALL, deserModePlain},
	{rhino.Token_OPTCHAIN_GETELEM, NodeKind_OPTCHAIN_ELEMENT_ACCESS, deserModePlain},
	{rhino.Token_COALESCE, NodeKind_COALESCE, deserModePlain},
	{rhino.Token_DYNAMIC_IMPORT, NodeKind_DYNAMIC_IMPORT, deserModePlain},
	{rhino.Token_ASSIGN_OR, NodeKind_ASSIGN_OR, deserModePlain},
	{rhino.Token_ASSIGN_AND, NodeKind_ASSIGN_AND, deserModePlain},
	{rhino.Token_ASSIGN_COALESCE, NodeKind_ASSIGN_COALESCE, deserModePlain},
	{rhino.Token_EXPR_RESULT, NodeKind_EXPRESSION_STATEMENT, deserModePlain},
	{rhino.Token_BREAK, NodeKind_BREAK_STATEMENT, deserModePlain},
	{rhino.Token_CONTINUE, NodeKind_CONTINUE_STATEMENT, deserModePlain},
	{rhino.Token_DEBUGGER, NodeKind_DEBUGGER_STATEMENT, deserModePlain},
	{rhino.Token_DO, NodeKind_DO_STATEMENT, deserModePlain},
	{rhino.Token_FOR, NodeKind_FOR_STATEMENT, deserModePlain},
	{rhino.Token_FOR_IN, NodeKind_FOR_IN_STATEMENT, deserModePlain},
	{rhino.Token_FOR_OF, NodeKind_FOR_OF_STATEMENT, deserModePlain},
	{rhino.Token_FOR_AWAIT_OF, NodeKind_FOR_AWAIT_OF_STATEMENT, deserModePlain},
	{rhino.Token_IF, NodeKind_IF_STATEMENT, deserModePlain},
	{rhino.Token_RETURN, NodeKind_RETURN_STATEMENT, deserModePlain},
	{rhino.Token_SWITCH, NodeKind_SWITCH_STATEMENT, deserModePlain},
	{rhino.Token_THROW, NodeKind_THROW_STATEMENT, deserModePlain},
	{rhino.Token_TRY, NodeKind_TRY_STATEMENT, deserModePlain},
	{rhino.Token_WHILE, NodeKind_WHILE_STATEMENT, deserModePlain},
	{rhino.Token_EMPTY, NodeKind_EMPTY, deserModePlain},
	{rhino.Token_WITH, NodeKind_WITH, deserModePlain},
	{rhino.Token_IMPORT, NodeKind_IMPORT, deserModePlain},
	{rhino.Token_EXPORT, NodeKind_EXPORT, deserModePlain},
	{rhino.Token_VAR, NodeKind_VAR_DECLARATION, deserModePlain},
	{rhino.Token_CONST, NodeKind_CONST_DECLARATION, deserModePlain},
	{rhino.Token_LET, NodeKind_LET_DECLARATION, deserModePlain},
	{rhino.Token_FUNCTION, NodeKind_FUNCTION_LITERAL, deserModePlain},
	{rhino.Token_CLASS, NodeKind_CLASS_LITERAL, deserModePlain},
	{rhino.Token_BLOCK, NodeKind_BLOCK, deserModePlain},
	{rhino.Token_LABEL, NodeKind_LABELED_STATEMENT, deserModePlain},
	{rhino.Token_LABEL_NAME, NodeKind_LABELED_NAME, deserModeString},
	{rhino.Token_CLASS_MEMBERS, NodeKind_CLASS_MEMBERS, deserModePlain},
	{rhino.Token_MEMBER_FUNCTION_DEF, NodeKind_METHOD_DECLARATION, deserModeString},
	{rhino.Token_MEMBER_FIELD_DEF, NodeKind_FIELD_DECLARATION, deserModeString},
	{rhino.Token_COMPUTED_FIELD_DEF, NodeKind_COMPUTED_PROP_FIELD, deserModePlain},
	{rhino.Token_PARAM_LIST, NodeKind_PARAMETER_LIST, deserModePlain},
	{rhino.Token_CASE, NodeKind_CASE, deserModePlain},
	{rhino.Token_DEFAULT_CASE, NodeKind_DEFAULT_CASE, deserModePlain},
	{rhino.Token_CATCH, NodeKind_CATCH, deserModePlain},
	{rhino.Token_SUPER, NodeKind_SUPER, deserModePlain},
	{rhino.Token_ARRAY_PATTERN, NodeKind_ARRAY_PATTERN, deserModePlain},
	{rhino.Token_OBJECT_PATTERN, NodeKind_OBJECT_PATTERN, deserModePlain},
	{rhino.Token_DESTRUCTURING_LHS, NodeKind_DESTRUCTURING_LHS, deserModePlain},
	{rhino.Token_DEFAULT_VALUE, NodeKind_DEFAULT_VALUE, deserModePlain},
	{rhino.Token_IMPORT_SPECS, NodeKind_IMPORT_SPECS, deserModePlain},
	{rhino.Token_IMPORT_SPEC, NodeKind_IMPORT_SPEC, deserModePlain},
	{rhino.Token_IMPORT_STAR, NodeKind_IMPORT_STAR, deserModeString},
	{rhino.Token_EXPORT_SPECS, NodeKind_EXPORT_SPECS, deserModePlain},
	{rhino.Token_EXPORT_SPEC, NodeKind_EXPORT_SPEC, deserModePlain},
	{rhino.Token_MODULE_BODY, NodeKind_MODULE_BODY, deserModePlain},
	{rhino.Token_ITER_REST, NodeKind_ITER_REST, deserModePlain},
	{rhino.Token_ITER_SPREAD, NodeKind_ITER_SPREAD, deserModePlain},
	{rhino.Token_OBJECT_REST, NodeKind_OBJECT_REST, deserModePlain},
	{rhino.Token_OBJECT_SPREAD, NodeKind_OBJECT_SPREAD, deserModePlain},
	{rhino.Token_SWITCH_BODY, NodeKind_SWITCH_BODY, deserModePlain},
}

func init() {
	for _, e := range bijectiveTokenNodeKinds {
		if tokenToNodeKind[e.token] != 0 || nodeKindToDeserSpec[e.kind].mode != deserModeInvalid {
			panic("duplicate token<->NodeKind entry")
		}
		tokenToNodeKind[e.token] = e.kind
		nodeKindToDeserSpec[e.kind] = nodeKindDeserSpec{token: e.token, mode: e.mode}
	}

	tokenToNodeKind[rhino.Token_INC] = nodeKindTokenConditional
	nodeKindToDeserSpec[NodeKind_PRE_INCREMENT] = nodeKindDeserSpec{token: rhino.Token_INC, mode: deserModePlain}
	nodeKindToDeserSpec[NodeKind_POST_INCREMENT] = nodeKindDeserSpec{token: rhino.Token_INC, mode: deserModePostIncrDecr}

	tokenToNodeKind[rhino.Token_DEC] = nodeKindTokenConditional
	nodeKindToDeserSpec[NodeKind_PRE_DECREMENT] = nodeKindDeserSpec{token: rhino.Token_DEC, mode: deserModePlain}
	nodeKindToDeserSpec[NodeKind_POST_DECREMENT] = nodeKindDeserSpec{token: rhino.Token_DEC, mode: deserModePostIncrDecr}

	tokenToNodeKind[rhino.Token_STRING_KEY] = nodeKindTokenConditional
	nodeKindToDeserSpec[NodeKind_RENAMABLE_STRING_KEY] = nodeKindDeserSpec{token: rhino.Token_STRING_KEY, mode: deserModeString}
	nodeKindToDeserSpec[NodeKind_QUOTED_STRING_KEY] = nodeKindDeserSpec{token: rhino.Token_STRING_KEY, mode: deserModeQuotedString}

	tokenToNodeKind[rhino.Token_GETTER_DEF] = nodeKindTokenConditional
	nodeKindToDeserSpec[NodeKind_RENAMABLE_GETTER_DEF] = nodeKindDeserSpec{token: rhino.Token_GETTER_DEF, mode: deserModeString}
	nodeKindToDeserSpec[NodeKind_QUOTED_GETTER_DEF] = nodeKindDeserSpec{token: rhino.Token_GETTER_DEF, mode: deserModeQuotedString}

	tokenToNodeKind[rhino.Token_SETTER_DEF] = nodeKindTokenConditional
	nodeKindToDeserSpec[NodeKind_RENAMABLE_SETTER_DEF] = nodeKindDeserSpec{token: rhino.Token_SETTER_DEF, mode: deserModeString}
	nodeKindToDeserSpec[NodeKind_QUOTED_SETTER_DEF] = nodeKindDeserSpec{token: rhino.Token_SETTER_DEF, mode: deserModeQuotedString}

	for k := NodeKind(1); k < NodeKind_UNRECOGNIZED; k++ {
		if nodeKindToDeserSpec[k].mode == deserModeInvalid {
			panic("missing NodeKind deserializer spec: " + k.String())
		}
	}
}

// Java: TypedAstSerializer.kindTranslator(Node)
func (s *TypedAstSerializer) KindTranslator(n rhino.Node) NodeKind {
	kind := tokenToNodeKind[n.GetToken()]
	if kind > 0 {
		return kind
	}
	return translateConditionalOrInvalidToken(n, kind)
}

func translateConditionalOrInvalidToken(n rhino.Node, kind NodeKind) NodeKind {
	if kind == nodeKindTokenConditional {
		switch n.GetToken() {
		case rhino.Token_INC:
			if n.GetBooleanProp(rhino.Node_INCRDECR_PROP) {
				return NodeKind_POST_INCREMENT
			}
			return NodeKind_PRE_INCREMENT
		case rhino.Token_DEC:
			if n.GetBooleanProp(rhino.Node_INCRDECR_PROP) {
				return NodeKind_POST_DECREMENT
			}
			return NodeKind_PRE_DECREMENT
		case rhino.Token_STRING_KEY:
			if n.IsQuotedStringKey() {
				return NodeKind_QUOTED_STRING_KEY
			}
			return NodeKind_RENAMABLE_STRING_KEY
		case rhino.Token_GETTER_DEF:
			if n.IsQuotedStringKey() {
				return NodeKind_QUOTED_GETTER_DEF
			}
			return NodeKind_RENAMABLE_GETTER_DEF
		case rhino.Token_SETTER_DEF:
			if n.IsQuotedStringKey() {
				return NodeKind_QUOTED_SETTER_DEF
			}
			return NodeKind_RENAMABLE_SETTER_DEF
		}
	}
	panic(jutil.NewIllegalStateExceptionMsg("Unserializable token for node: " + n.String()))
}
