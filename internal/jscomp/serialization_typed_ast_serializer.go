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

	nodeSlabs        [][]AstNode
	nodeSlabIdx      int
	nodeSlabPos      int
	childPtrSlabs    [][]*AstNode
	childSlabIdx     int
	childSlabPos     int
	childStack       []*AstNode
	jsdocSlabs       [][]OptimizationJsdoc
	jsdocSlabIdx     int
	jsdocSlabPos     int
	jsdocKindSlabs   [][]int32
	jsdocKindSlabIdx int
	jsdocKindSlabPos int
	jsdocKindScratch []int32
	templateSlabs    [][]TemplateStringValue
	templateSlabIdx  int
	templateSlabPos  int
	scriptByteSlab   []byte
	lazyAstSlab      []LazyAst
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

func (s *TypedAstSerializer) resetScriptSlabs() {
	s.nodeSlabIdx = 0
	s.nodeSlabPos = 0
	s.childSlabIdx = 0
	s.childSlabPos = 0
	s.childStack = s.childStack[:0]
	s.jsdocSlabIdx = 0
	s.jsdocSlabPos = 0
	s.jsdocKindSlabIdx = 0
	s.jsdocKindSlabPos = 0
	s.templateSlabIdx = 0
	s.templateSlabPos = 0
}

func (s *TypedAstSerializer) allocAstNode() *AstNode {
	const chunkSize = 256
	if s.nodeSlabIdx >= len(s.nodeSlabs) {
		s.nodeSlabs = append(s.nodeSlabs, make([]AstNode, chunkSize))
	}
	slab := s.nodeSlabs[s.nodeSlabIdx]
	if s.nodeSlabPos >= len(slab) {
		s.nodeSlabIdx++
		s.nodeSlabPos = 0
		if s.nodeSlabIdx >= len(s.nodeSlabs) {
			s.nodeSlabs = append(s.nodeSlabs, make([]AstNode, chunkSize))
		}
		slab = s.nodeSlabs[s.nodeSlabIdx]
	}
	m := &slab[s.nodeSlabPos]
	s.nodeSlabPos++
	*m = AstNode{frozen: true, memoizedSize: -1}
	return m
}

func (s *TypedAstSerializer) allocChildren(children []*AstNode) []*AstNode {
	const chunkSize = 4096
	n := len(children)
	if n == 0 {
		return nil
	}
	if s.childSlabIdx >= len(s.childPtrSlabs) {
		s.childPtrSlabs = append(s.childPtrSlabs, make([]*AstNode, max(chunkSize, n)))
	}
	slab := s.childPtrSlabs[s.childSlabIdx]
	if s.childSlabPos+n > len(slab) {
		s.childSlabIdx++
		s.childSlabPos = 0
		if s.childSlabIdx >= len(s.childPtrSlabs) || len(s.childPtrSlabs[s.childSlabIdx]) < n {
			s.childPtrSlabs = append(s.childPtrSlabs, make([]*AstNode, max(chunkSize, n)))
			s.childSlabIdx = len(s.childPtrSlabs) - 1
		}
		slab = s.childPtrSlabs[s.childSlabIdx]
	}
	dst := slab[s.childSlabPos : s.childSlabPos+n : s.childSlabPos+n]
	s.childSlabPos += n
	copy(dst, children)
	return dst
}

func (s *TypedAstSerializer) allocJsdoc(src OptimizationJsdoc) *OptimizationJsdoc {
	const chunkSize = 512
	if s.jsdocSlabIdx >= len(s.jsdocSlabs) {
		s.jsdocSlabs = append(s.jsdocSlabs, make([]OptimizationJsdoc, chunkSize))
	}
	slab := s.jsdocSlabs[s.jsdocSlabIdx]
	if s.jsdocSlabPos >= len(slab) {
		s.jsdocSlabIdx++
		s.jsdocSlabPos = 0
		if s.jsdocSlabIdx >= len(s.jsdocSlabs) {
			s.jsdocSlabs = append(s.jsdocSlabs, make([]OptimizationJsdoc, chunkSize))
		}
		slab = s.jsdocSlabs[s.jsdocSlabIdx]
	}
	dst := &slab[s.jsdocSlabPos]
	s.jsdocSlabPos++
	*dst = src
	if n := len(src.kind_); n > 0 {
		if s.jsdocKindSlabIdx >= len(s.jsdocKindSlabs) {
			s.jsdocKindSlabs = append(s.jsdocKindSlabs, make([]int32, max(1024, n)))
		}
		kslab := s.jsdocKindSlabs[s.jsdocKindSlabIdx]
		if s.jsdocKindSlabPos+n > len(kslab) {
			s.jsdocKindSlabIdx++
			s.jsdocKindSlabPos = 0
			if s.jsdocKindSlabIdx >= len(s.jsdocKindSlabs) || len(s.jsdocKindSlabs[s.jsdocKindSlabIdx]) < n {
				s.jsdocKindSlabs = append(s.jsdocKindSlabs, make([]int32, max(1024, n)))
				s.jsdocKindSlabIdx = len(s.jsdocKindSlabs) - 1
			}
			kslab = s.jsdocKindSlabs[s.jsdocKindSlabIdx]
		}
		kdst := kslab[s.jsdocKindSlabPos : s.jsdocKindSlabPos+n : s.jsdocKindSlabPos+n]
		s.jsdocKindSlabPos += n
		copy(kdst, src.kind_)
		dst.kind_ = kdst
	} else {
		dst.kind_ = nil
	}
	dst.frozen = true
	dst.memoizedSize = dst.scalarsSerializedSize()
	return dst
}

func (s *TypedAstSerializer) allocTemplateStringValue(rawPtr, cookedPtr int32) *TemplateStringValue {
	if s.templateSlabIdx >= len(s.templateSlabs) {
		s.templateSlabs = append(s.templateSlabs, make([]TemplateStringValue, 64))
	}
	slab := s.templateSlabs[s.templateSlabIdx]
	if s.templateSlabPos >= len(slab) {
		s.templateSlabIdx++
		s.templateSlabPos = 0
		if s.templateSlabIdx >= len(s.templateSlabs) {
			s.templateSlabs = append(s.templateSlabs, make([]TemplateStringValue, 64))
		}
		slab = s.templateSlabs[s.templateSlabIdx]
	}
	dst := &slab[s.templateSlabPos]
	s.templateSlabPos++
	*dst = TemplateStringValue{
		rawStringPointer_:    rawPtr,
		cookedStringPointer_: cookedPtr,
		frozen:               true,
	}
	dst.memoizedSize = dst.scalarsSerializedSize()
	return dst
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

	builder := TypedAst_NewBuilder()
	if numThreads > 1 && s.jsTypeSerializer == nil {
		s.serializeRootsParallel(externsRoot, jsRoot, numThreads, builder)
	} else {
		for script := externsRoot.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
			if NodeUtil_IsFromTypeSummary(script) {
				continue
			}
			builder.AddExternAst(s.SerializeScriptNode(script))
		}
		for script := jsRoot.GetFirstChild(); !script.IsNil(); script = script.GetNext() {
			if NodeUtil_IsFromTypeSummary(script) {
				continue
			}
			builder.AddCodeAst(s.SerializeScriptNode(script))
		}
	}

	sourceFileProtos := make([]*SourceFileProto, 0, len(s.sourceFileList))
	for _, f := range s.sourceFileList {
		sourceFileProtos = append(sourceFileProtos, f.GetProto())
	}
	sourceFiles := SourceFilePool_NewBuilder().
		AddAllSourceFile(sourceFileProtos).
		Build()

	if s.compiler.GetExternProperties() != nil {
		externsSummary := ExternsSummary_NewBuilder()
		for prop := range s.compiler.GetExternProperties().All() {
			externsSummary.AddPropNamePtr(int32(s.stringPool.Put(prop)))
		}
		builder.SetExternsSummaryFromBuilder(externsSummary)
	}

	return builder.
		AddAllRuntimeLibraryToInject(s.serializationMode.RuntimeLibraries()).
		SetTypePool(s.generateTypePool()).
		SetStringPool(s.stringPool.Build().ToProtoWithThreads(numThreads)).
		SetSourceFilePool(sourceFiles).
		Build()
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
	s.resetScriptSlabs()

	sourceFile := s.GetSourceFilePointer(script)
	scriptProto := s.visitArena(script)
	s.subtreeSourceFiles = s.subtreeSourceFiles[:0]

	encodedSourceMap := s.compiler.GetBase64SourceMapContents(*script.GetSourceFileName())

	var scriptBytes []byte
	sz := scriptProto.serializedSize()
	if sz <= 262144 {
		if len(s.scriptByteSlab) < sz {
			s.scriptByteSlab = make([]byte, max(4*1024*1024, sz))
		}
		buf := s.scriptByteSlab[:0:sz]
		s.scriptByteSlab = s.scriptByteSlab[sz:]
		scriptBytes = scriptProto.appendTo(buf)
	} else {
		scriptBytes = scriptProto.ToByteString()
	}

	if len(s.lazyAstSlab) == 0 {
		s.lazyAstSlab = make([]LazyAst, 256)
	}
	lazyAst := &s.lazyAstSlab[0]
	s.lazyAstSlab = s.lazyAstSlab[1:]
	*lazyAst = LazyAst{
		script_:      scriptBytes,
		sourceFile_:  int32(sourceFile),
		frozen:       true,
		memoizedSize: -1,
	}

	if encodedSourceMap != nil {
		// This is the encoded source map taken from the inline sourcemap comment. It does not include
		// the base64 prefix.
		// E.g. We serialize "eyJ2ZXJzaW9uI..." from the "//# sourceMappingURL=
		// data:application/json;base64,eyJ2ZXJzaW9uI..." comment.
		lazyAst.sourceMappingUrl_ = *encodedSourceMap
	}

	return lazyAst
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

func (s *TypedAstSerializer) applyFixupsSequential(fixups []serializerFixup) {
	var cs *ColorSerializer
	if s.colorTypeSerializer != nil {
		cs = s.colorTypeSerializer.colorSerializer
	}
	for i := range fixups {
		op := &fixups[i]
		switch op.kind {
		case fixupKindStringID:
			*op.dst = int32(s.stringPool.PutID(op.strID))
		case fixupKindStringIDUsedProp:
			ptr := s.stringPool.PutID(op.strID)
			s.markUsedProp(ptr)
			*op.dst = int32(ptr)
		case fixupKindColor:
			*op.dst = int32(cs.AddColor(op.color))
		case fixupKindSourceFile:
			*op.dst = int32(s.internSourceFile(op.sf))
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
	builder *TypedAst_Builder,
) {
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
		return
	}
	if total < 4 {
		for i, script := range allScripts {
			lazyAst := s.SerializeScriptNode(script)
			if i < numExtern {
				builder.AddExternAst(lazyAst)
			} else {
				builder.AddCodeAst(lazyAst)
			}
		}
		return
	}

	batchSize := min(total, max(numThreads*8, 16))
	arenas := make([]typedAstScriptArena, batchSize)
	workerByteSlabs := make([][]byte, numThreads)
	lazyAstOut := make([]LazyAst, total)

	for base := 0; base < total; base += batchSize {
		batchEnd := min(base+batchSize, total)
		batch := allScripts[base:batchEnd]
		batchLen := len(batch)

		// Pass A: Build per-script AstNode arenas and ordered fixup streams in parallel.
		runSerializerBatchParallel(batchLen, numThreads, func(_ int, slot int) {
			arenas[slot].buildScriptPassA(batch[slot], s)
		})

		// Pass B: Replay pool registrations in strict script & AST pre/post order.
		for slot := 0; slot < batchLen; slot++ {
			s.applyFixupsSequential(arenas[slot].fixups)
		}

		// Pass C: Finalize bottom-up memoized sizes and marshal LazyAst bytes in parallel.
		runSerializerBatchParallel(batchLen, numThreads, func(workerIdx int, slot int) {
			arenas[slot].marshalScriptPassC(&workerByteSlabs[workerIdx])
			lazyAstOut[base+slot] = arenas[slot].lazyAst
		})
	}

	for i := 0; i < numExtern; i++ {
		builder.AddExternAst(&lazyAstOut[i])
	}
	for i := numExtern; i < total; i++ {
		builder.AddCodeAst(&lazyAstOut[i])
	}
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

func (s *TypedAstSerializer_JSTypeSerializer) addTypeToNode(node rhino.Node, m *AstNode) {
	tb := node.GetJSTypeBase()
	if tb == nil {
		return
	}
	pointer, ok := s.reconserializer.typePointers[tb]
	if !ok {
		jutil.CheckNotNullMsg1[any](nil, "cannot find pointer for %s", tb.Self())
	}
	m.type_ = pointer
	m.bitField0_ |= 1 << 0
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
