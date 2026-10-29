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
	"slices"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
	"github.com/paullepoulpe/jsgo/internal/strpool"
)

// TypedAstSerializer transforms a compiler AST into a serialized TypedAst
// object.
//
// Java: final class TypedAstSerializer
type TypedAstSerializer struct {
	compiler          AbstractCompiler
	serializationMode *SerializationOptions
	stringPool        *StringPool_Builder
	previousLine      int
	previousColumn    int
	// Go-only: language-adapter — the ArrayDeque<SourceFile> is used only as
	// a stack (addLast, removeLast, peekLast, clear): a slice, top last.
	subtreeSourceFiles []*SourceFile
	sourceFilePointers map[*SourceFile]int
	sourceFileList     []*SourceFile

	colorTypeSerializer *TypedAstSerializer_ColorTypeSerializer
	jsTypeSerializer    *TypedAstSerializer_JSTypeSerializer
	usedPropBits        []uint64
	recordUsedProps     bool

	scriptEnc      scriptWireEncoder
	lazyAstArena   slabArena[LazyAst]
	scriptByteSlab []byte
}

// slabArena is a two-level bump allocator that retains backing slabs across resets.
type slabArena[T any] struct {
	slabs   [][]T
	slabIdx int
	slabPos int
}

func (a *slabArena[T]) alloc(chunkSize int) *T {
	if a.slabIdx >= len(a.slabs) {
		a.slabs = append(a.slabs, make([]T, chunkSize))
	}
	slab := a.slabs[a.slabIdx]
	if a.slabPos >= len(slab) {
		a.slabIdx++
		a.slabPos = 0
		if a.slabIdx >= len(a.slabs) {
			a.slabs = append(a.slabs, make([]T, chunkSize))
		}
		slab = a.slabs[a.slabIdx]
	}
	p := &slab[a.slabPos]
	a.slabPos++
	return p
}

// Java: TypedAstSerializer(AbstractCompiler, SerializationOptions)
func NewTypedAstSerializer(compiler AbstractCompiler, serializationMode *SerializationOptions) *TypedAstSerializer {
	return &TypedAstSerializer{
		compiler:           compiler,
		serializationMode:  serializationMode,
		stringPool:         StringPool_Builder_(),
		sourceFilePointers: make(map[*SourceFile]int, 64),
	}
}

// SerializeRoots transforms the given compiler AST root nodes into into a
// serialized TypedAst object.
//
// Java: TypedAstSerializer.serializeRoots(Node, Node)
func (s *TypedAstSerializer) SerializeRoots(externsRoot rhino.Node, jsRoot rhino.Node) *TypedAst {
	jutil.CheckArgument(externsRoot.IsRoot())
	jutil.CheckArgument(jsRoot.IsRoot())

	if s.compiler.HasOptimizationColors() {
		if root := s.compiler.GetRoot(); !root.IsNil() && externsRoot.GetParent() == root {
			s.usedPropBits = make([]uint64, 64)
			s.recordUsedProps = true
			colorSerializer := NewColorSerializer(
				s.serializationMode,
				s.stringPool.Put,
				func(prop string) bool {
					_, ok := s.lookupUsedPropIndex(prop)
					return ok
				})
			colorSerializer.lookupUsedPropIndexFn = s.lookupUsedPropIndex
			s.colorTypeSerializer = NewTypedAstSerializer_ColorTypeSerializer(colorSerializer, s.compiler.GetColorRegistry())
		} else {
			s.colorTypeSerializer = s.CreateColorTypeSerializer(s.compiler, s.serializationMode, s.stringPool)
		}
	} else if s.compiler.HasTypeCheckingRun() {
		// The AST has JSTypes, but we want to serialize colors instead.
		// TODO(bradfordcsmith): Change this branch to throw an error and delete the related logic.
		//     Nothing should be using it anymore.
		s.jsTypeSerializer = TypedAstSerializer_CreateJSTypeSerializer(s.compiler, s.serializationMode, s.stringPool)
	}

	numThreads := 1
	if opts := s.compiler.GetOptions(); opts != nil && opts.GetNumParallelThreads() > 1 {
		numThreads = opts.GetNumParallelThreads()
	}

	var externAsts []*LazyAst
	var codeAsts []*LazyAst
	if numThreads > 1 && s.jsTypeSerializer == nil {
		externAsts, codeAsts = s.serializeRootsParallel(externsRoot, jsRoot, numThreads)
	} else {
		for script := externsRoot.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
			if NodeUtil_IsFromTypeSummary(script) {
				continue
			}
			externAsts = append(externAsts, s.SerializeScriptNode(script))
		}
		for script := jsRoot.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
			if NodeUtil_IsFromTypeSummary(script) {
				continue
			}
			codeAsts = append(codeAsts, s.SerializeScriptNode(script))
		}
	}

	sourceFileProtos := make([]*SourceFileProto, 0, len(s.sourceFileList))
	for _, f := range s.sourceFileList {
		sourceFileProtos = append(sourceFileProtos, f.GetProto())
	}
	sourceFiles := &SourceFilePool{
		sourceFile_: sourceFileProtos,
	}

	var externsSummary *ExternsSummary
	if s.compiler.GetExternProperties() != nil {
		var propNamePtrs []int32
		for prop := range s.compiler.GetExternProperties().All() {
			propNamePtrs = append(propNamePtrs, int32(s.stringPool.Put(prop)))
		}
		externsSummary = &ExternsSummary{
			propNamePtr_: propNamePtrs,
		}
	}

	var runtimeLibs []string
	if libs := s.serializationMode.RuntimeLibraries(); len(libs) > 0 {
		runtimeLibs = slices.Clone(libs)
	}

	return &TypedAst{
		typePool_:               s.generateTypePool(),
		stringPool_:             s.stringPool.Build().ToProtoWithThreads(numThreads),
		externAst_:              externAsts,
		codeAst_:                codeAsts,
		sourceFilePool_:         sourceFiles,
		externsSummary_:         externsSummary,
		runtimeLibraryToInject_: runtimeLibs,
	}
}

func (s *TypedAstSerializer) generateTypePool() *TypePool {
	if s.colorTypeSerializer != nil {
		return s.colorTypeSerializer.GenerateTypePool()
	}
	if s.jsTypeSerializer != nil {
		return s.jsTypeSerializer.GenerateTypePool()
	}
	return TypePool_GetDefaultInstance()
}

// Java: TypedAstSerializer.serializeScriptNode(Node)
func (s *TypedAstSerializer) SerializeScriptNode(script rhino.Node) *LazyAst {
	jutil.CheckState(script.IsScript())
	s.previousColumn = 0
	s.previousLine = s.previousColumn
	s.scriptEnc.reset()

	sourceFile := s.GetSourceFilePointer(script)
	s.visitDirect(script)
	s.subtreeSourceFiles = s.subtreeSourceFiles[:0]

	encodedSourceMap := s.compiler.GetBase64SourceMapContents(*script.GetSourceFileName())
	scriptBytes := s.scriptEnc.encodeToByteSlab(&s.scriptByteSlab)

	lazyAst := s.lazyAstArena.alloc(256)
	*lazyAst = LazyAst{
		script_:     scriptBytes,
		sourceFile_: int32(sourceFile),
	}

	if encodedSourceMap != nil {
		// This is the encoded source map taken from the inline sourcemap comment. It does not include
		// the base64 prefix.
		// E.g. We serialize "eyJ2ZXJzaW9uI..." from the "//# sourceMappingURL=
		// data:application/json;base64,eyJ2ZXJzaW9uI..." comment.
		lazyAst.sourceMappingUrl_ = *encodedSourceMap
		lazyAst.serializedSize()
	} else {
		lazyAst.cachedSize = lazyAstFastSize(scriptBytes, lazyAst.sourceFile_) + 1
	}

	return lazyAst
}

func lazyAstFastSize(scriptBytes []byte, sourceFile int32) int {
	ls := 0
	if len(scriptBytes) != 0 {
		ls += 1 + protoVarintLen(uint64(len(scriptBytes))) + len(scriptBytes)
	}
	if sourceFile != 0 {
		ls += 1 + protoVarintLen(uint64(uint32(sourceFile)))
	}
	return ls
}

// peekLastSubtreeSourceFile is subtreeSourceFiles.peekLast(): the top of the
// stack, or null when it is empty.
//
// Go-only: language-adapter — ArrayDeque.peekLast on the slice stack.
func (s *TypedAstSerializer) peekLastSubtreeSourceFile() *SourceFile {
	if len(s.subtreeSourceFiles) == 0 {
		return nil
	}
	return s.subtreeSourceFiles[len(s.subtreeSourceFiles)-1]
}

// Java: TypedAstSerializer.getSourceFilePointer(Node)
func (s *TypedAstSerializer) GetSourceFilePointer(n rhino.Node) int {
	var sourceFile *SourceFile
	if ssf := n.GetStaticSourceFile(); ssf != nil {
		sourceFile = ssf.(*SourceFile)
	}
	if sourceFile == nil {
		// TODO(b/186056977): enforce that SourceFile is not null in externs as well as code.
		if !s.peekLastSubtreeSourceFile().IsExtern() {
			jutil.CheckStateMsg2(
				false,
				"Unexpected null SourceFile for node %s with parent %s",
				n.ToStringTree(),
				n.GetParent())
		}
		return 0 // not set
	}

	if sourceFile == s.peekLastSubtreeSourceFile() {
		// To save space, only serialize a SourceFile for a node if it is different than the parent's
		// source.
		return 0 // not set
	}

	s.subtreeSourceFiles = append(s.subtreeSourceFiles, sourceFile)
	return s.internSourceFile(sourceFile)
}

func (s *TypedAstSerializer) internSourceFile(sourceFile *SourceFile) int {
	if idx, ok := s.sourceFilePointers[sourceFile]; ok {
		return idx
	}
	idx := 1 + len(s.sourceFileList)
	s.sourceFileList = append(s.sourceFileList, sourceFile)
	s.sourceFilePointers[sourceFile] = idx
	return idx
}

func (s *TypedAstSerializer) applyFixupsSequential(enc *scriptWireEncoder) {
	var cs *ColorSerializer
	if s.colorTypeSerializer != nil {
		cs = s.colorTypeSerializer.colorSerializer
	}
	sp := s.stringPool
	nodes := enc.nodes
	jsdocs := enc.jsdocs
	fixups := enc.fixups
	for i := range fixups {
		op := fixups[i]
		switch op.target {
		case fixupNodeStrVal:
			nodes[op.targetIdx].valBits = uint64(uint32(int32(sp.PutID(op.refIdx))))
		case fixupNodeStrValUsedProp:
			ptr := sp.PutID(op.refIdx)
			s.markUsedProp(ptr)
			nodes[op.targetIdx].valBits = uint64(uint32(int32(ptr)))
		case fixupNodeOrigName:
			nodes[op.targetIdx].origName = int32(sp.PutID(op.refIdx))
		case fixupNodeColor:
			nodes[op.targetIdx].typePtr = int32(cs.AddColor(enc.colorRefs[op.refIdx]))
		case fixupNodeRawTmpl:
			ptr := uint32(int32(sp.PutID(op.refIdx)))
			n := &nodes[op.targetIdx]
			n.valBits = (n.valBits & 0xFFFFFFFF) | (uint64(ptr) << 32)
		case fixupNodeCookedTmpl:
			ptr := uint32(int32(sp.PutID(op.refIdx)))
			n := &nodes[op.targetIdx]
			n.valBits = (n.valBits &^ 0xFFFFFFFF) | uint64(ptr)
		case fixupNodeSourceFile:
			nodes[op.targetIdx].sourceFile = int32(s.internSourceFile(enc.sfRefs[op.refIdx]))
		case fixupLazySourceFile:
			enc.lazyAst.sourceFile_ = int32(s.internSourceFile(enc.sfRefs[op.refIdx]))
		case fixupJsdocLicense:
			jsdocs[op.targetIdx].licensePtr = int32(sp.PutID(op.refIdx))
		case fixupJsdocDesc:
			jsdocs[op.targetIdx].descPtr = int32(sp.PutID(op.refIdx))
		case fixupJsdocAltMsg:
			jsdocs[op.targetIdx].altMsgPtr = int32(sp.PutID(op.refIdx))
		case fixupJsdocMeaning:
			jsdocs[op.targetIdx].meaningPtr = int32(sp.PutID(op.refIdx))
		case fixupJsdocCuMode:
			jsdocs[op.targetIdx].cuModePtr = int32(sp.PutID(op.refIdx))
		}
	}
}

func runSerializerBatchParallel(batchLen, numThreads int, fn func(workerIdx, slot int)) {
	workers := min(numThreads, batchLen)
	if workers <= 1 {
		for slot := 0; slot < batchLen; slot++ {
			fn(0, slot)
		}
		return
	}
	var next atomic.Int32
	var panicked atomic.Bool
	panics := make([]any, batchLen)
	hasPanic := make([]bool, batchLen)
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 0; w < workers; w++ {
		workerIdx := w
		go func() {
			defer wg.Done()
			for !panicked.Load() {
				slot := int(next.Add(1) - 1)
				if slot >= batchLen {
					return
				}
				func() {
					defer func() {
						if r := recover(); r != nil {
							panics[slot] = r
							hasPanic[slot] = true
							panicked.Store(true)
						}
					}()
					fn(workerIdx, slot)
				}()
			}
		}()
	}
	wg.Wait()
	if panicked.Load() {
		for slot := 0; slot < batchLen; slot++ {
			if hasPanic[slot] {
				panic(panics[slot])
			}
		}
	}
}

func (s *TypedAstSerializer) serializeRootsParallel(
	externsRoot rhino.Node,
	jsRoot rhino.Node,
	numThreads int,
) (externAsts, codeAsts []*LazyAst) {
	var allScripts []rhino.Node
	for script := externsRoot.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
		if !NodeUtil_IsFromTypeSummary(script) {
			allScripts = append(allScripts, script)
		}
	}
	numExtern := len(allScripts)
	for script := jsRoot.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
		if !NodeUtil_IsFromTypeSummary(script) {
			allScripts = append(allScripts, script)
		}
	}
	total := len(allScripts)
	if total == 0 {
		return nil, nil
	}
	if total < 4 {
		if numExtern > 0 {
			externAsts = make([]*LazyAst, 0, numExtern)
		}
		if total > numExtern {
			codeAsts = make([]*LazyAst, 0, total-numExtern)
		}
		for i, script := range allScripts {
			lazyAst := s.SerializeScriptNode(script)
			if i < numExtern {
				externAsts = append(externAsts, lazyAst)
			} else {
				codeAsts = append(codeAsts, lazyAst)
			}
		}
		return externAsts, codeAsts
	}

	batchSize := min(total, max(numThreads*8, 16))
	encoders := make([]scriptWireEncoder, batchSize)
	workerByteSlabs := make([][]byte, numThreads)
	lazyAstOut := make([]LazyAst, total)

	for base := 0; base < total; base += batchSize {
		batchEnd := min(base+batchSize, total)
		batch := allScripts[base:batchEnd]
		batchLen := len(batch)

		// Pass A: Build per-script flat AST tables and pointer-free fixup streams in parallel.
		runSerializerBatchParallel(batchLen, numThreads, func(_ int, slot int) {
			encoders[slot].buildScriptPassA(batch[slot], s)
		})

		// Pass B: Replay pool registrations in strict script & AST pre/post order.
		for slot := 0; slot < batchLen; slot++ {
			s.applyFixupsSequential(&encoders[slot])
		}

		// Pass C: Compute bottom-up wire sizes and emit protobuf wire bytes directly in parallel.
		runSerializerBatchParallel(batchLen, numThreads, func(workerIdx int, slot int) {
			encoders[slot].marshalScriptPassC(&workerByteSlabs[workerIdx])
			lazyAstOut[base+slot] = encoders[slot].lazyAst
		})
	}

	lazyPtrs := make([]*LazyAst, total)
	for i := 0; i < total; i++ {
		lazyPtrs[i] = &lazyAstOut[i]
	}
	if numExtern > 0 {
		externAsts = lazyPtrs[:numExtern:numExtern]
	}
	if total > numExtern {
		codeAsts = lazyPtrs[numExtern:total:total]
	}
	return externAsts, codeAsts
}

func (s *TypedAstSerializer) markUsedProp(ptr int) {
	word := ptr >> 6
	if word >= len(s.usedPropBits) {
		newLen := max(len(s.usedPropBits)*2, word+1)
		grown := make([]uint64, newLen)
		copy(grown, s.usedPropBits)
		s.usedPropBits = grown
	}
	s.usedPropBits[word] |= uint64(1) << (uint(ptr) & 63)
}

func (s *TypedAstSerializer) lookupUsedPropIndex(prop string) (int, bool) {
	idx, ok := s.stringPool.Lookup(prop)
	if !ok {
		return 0, false
	}
	word := idx >> 6
	if word >= len(s.usedPropBits) || (s.usedPropBits[word]&(uint64(1)<<(uint(idx)&63))) == 0 {
		return 0, false
	}
	return idx, true
}

// TypedAstSerializer_CreateJSTypeSerializer creates the `TypeSerializer`
// appropriate for an AST that contains JSTypes.
//
// Java: TypedAstSerializer.createJSTypeSerializer(AbstractCompiler, SerializationOptions, StringPool.Builder)
func TypedAstSerializer_CreateJSTypeSerializer(
	compiler AbstractCompiler,
	serializationMode *SerializationOptions,
	stringPoolBuilder *StringPool_Builder) *TypedAstSerializer_JSTypeSerializer {
	serializeTypesToPointers := SerializeTypesToPointers_Create(compiler, stringPoolBuilder, serializationMode)
	// Gather and serialize all the types now.
	serializeTypesToPointers.GatherTypesOnAst(compiler.GetRoot())
	return &TypedAstSerializer_JSTypeSerializer{
		reconserializer: serializeTypesToPointers.jstypeReconserializer,
		typePool:        serializeTypesToPointers.GetTypePool(),
	}
}

// TypedAstSerializer_JSTypeSerializer is used when the AST's JSTypes have not
// been converted to Colors.
//
// Java: private static class TypedAstSerializer.JSTypeSerializer
type TypedAstSerializer_JSTypeSerializer struct {
	// Everything is pre-calculated with this form of serialization.
	reconserializer *JSTypeReconserializer
	typePool        *TypePool
}

func (s *TypedAstSerializer_JSTypeSerializer) addTypeToFlatNode(node rhino.Node, fn *flatAstNode) {
	tb := node.GetJSTypeBase()
	if tb == nil {
		return
	}
	pointer, ok := s.reconserializer.typePointers[tb]
	if !ok {
		jutil.CheckNotNullMsg1[any](nil, "cannot find pointer for %s", tb.Self())
	}
	fn.typePtr = pointer
	fn.hasType = true
}

// Java: TypedAstSerializer.JSTypeSerializer.generateTypePool()
func (s *TypedAstSerializer_JSTypeSerializer) GenerateTypePool() *TypePool {
	return s.typePool
}

// CreateColorTypeSerializer creates a `TypeSerializer` that knows how to
// serialize `Color`s from the AST.
//
// Java: TypedAstSerializer.createColorTypeSerializer(AbstractCompiler, SerializationOptions, StringPool.Builder)
func (s *TypedAstSerializer) CreateColorTypeSerializer(
	compiler AbstractCompiler,
	serializationMode *SerializationOptions,
	stringPoolBuilder *StringPool_Builder) *TypedAstSerializer_ColorTypeSerializer {
	// Gather all the property names that are actually used in the AST.
	var usedPropertyNames jutil.BitSet
	var lastPropID uint32
	var walk func(n rhino.Node)
	walk = func(n rhino.Node) {
		if n.IsScript() && NodeUtil_IsFromTypeSummary(n) {
			return
		}
		for c := n.GetFirstChild(); !c.IsNil(); c = c.GetNext() {
			walk(c)
		}
		switch n.GetToken() {
		case rhino.Token_GETPROP,
			rhino.Token_OPTCHAIN_GETPROP:
			propID := n.GetStringID()
			if propID != lastPropID {
				lastPropID = propID
				usedPropertyNames.Set(int(propID))
			}
		case rhino.Token_STRING_KEY,
			rhino.Token_MEMBER_FUNCTION_DEF,
			rhino.Token_MEMBER_FIELD_DEF,
			rhino.Token_GETTER_DEF,
			rhino.Token_SETTER_DEF:
			if !n.IsQuotedStringKey() {
				propID := n.GetStringID()
				if propID != lastPropID {
					lastPropID = propID
					usedPropertyNames.Set(int(propID))
				}
			}
		case rhino.Token_NAME:
			if shadow := n.GetClosureUnawareShadow(); !shadow.IsNil() {
				walk(shadow.GetFirstFirstChild().GetFirstChild())
			}
		}
	}
	walk(compiler.GetRoot())
	colorSerializer := NewColorSerializer(
		serializationMode,
		// lookup / allocate strings from the shared string pool
		stringPoolBuilder.Put,
		// only include property names known to be used in the AST
		func(prop string) bool {
			id := strpool.Lookup(prop)
			return id != 0 && usedPropertyNames.Get(int(id))
		})
	return NewTypedAstSerializer_ColorTypeSerializer(colorSerializer, compiler.GetColorRegistry())
}

// TypedAstSerializer_ColorTypeSerializer is used when the AST has `Color`s
// rather than `JSType`s.
//
// Java: private static class TypedAstSerializer.ColorTypeSerializer
type TypedAstSerializer_ColorTypeSerializer struct {
	colorSerializer *ColorSerializer
	colorRegistry   *colors.ColorRegistry
}

// Java: TypedAstSerializer.ColorTypeSerializer(ColorSerializer, ColorRegistry)
func NewTypedAstSerializer_ColorTypeSerializer(colorSerializer *ColorSerializer, colorRegistry *colors.ColorRegistry) *TypedAstSerializer_ColorTypeSerializer {
	return &TypedAstSerializer_ColorTypeSerializer{colorSerializer: colorSerializer, colorRegistry: colorRegistry}
}

// Java: TypedAstSerializer.ColorTypeSerializer.generateTypePool()
func (s *TypedAstSerializer_ColorTypeSerializer) GenerateTypePool() *TypePool {
	mismatchLocationsForDebugging := s.colorRegistry.GetMismatchLocationsForDebugging()
	return s.colorSerializer.GenerateTypePool(
		s.colorRegistry.GetDisambiguationSupertypesSlice,
		func(color *colors.Color) *jutil.LinkedHashSet[string] {
			// Go-only: language-adapter — the ColorSerializer port takes the
			// ImmutableSet<String> as a LinkedHashSet, the SetMultimap port
			// returns an ObjectLinkedHashSet: copied in iteration order.
			locations := jutil.NewLinkedHashSet[string]()
			for location := range mismatchLocationsForDebugging.Get(color.GetId()).All() {
				locations.Add(location)
			}
			return locations
		})
}
