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
	"github.com/paullepoulpe/jsgo/internal/colors"
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
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

	typeSerializer      TypedAstSerializer_TypeSerializer
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
	s := &TypedAstSerializer{
		stringPool:         StringPool_Builder_(),
		sourceFilePointers: make(map[*SourceFile]int, 64),
		typeSerializer:     nil,
	}
	s.compiler = compiler
	s.serializationMode = serializationMode
	return s
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
		s.typeSerializer = s.colorTypeSerializer
	} else if s.compiler.HasTypeCheckingRun() {
		// The AST has JSTypes, but we want to serialize colors instead.
		// TODO(bradfordcsmith): Change this branch to throw an error and delete the related logic.
		//     Nothing should be using it anymore.
		s.jsTypeSerializer = TypedAstSerializer_CreateJSTypeSerializer(s.compiler, s.serializationMode, s.stringPool)
		s.typeSerializer = s.jsTypeSerializer
	} else {
		s.typeSerializer = NewTypedAstSerializer_NoOpTypeSerializer()
	}

	builder := TypedAst_NewBuilder()
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
		SetTypePool(s.typeSerializer.GenerateTypePool()).
		SetStringPool(s.stringPool.Build().ToProto()).
		SetSourceFilePool(sourceFiles).
		Build()
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

// Java: TypedAstSerializer.createWithPositionInfo(Node)
func (s *TypedAstSerializer) CreateWithPositionInfo(n rhino.Node) *AstNode_Builder {
	jutil.CheckState(n.GetLength() >= 0)
	currentLine := n.GetLineno()
	currentColumn := n.GetCharno()
	builder := AstNode_NewBuilder().
		SetRelativeLine(int32(currentLine - s.previousLine)).
		SetRelativeColumn(int32(currentColumn - s.previousColumn))
	s.previousLine = currentLine
	s.previousColumn = currentColumn
	return builder
}

// Java: TypedAstSerializer.visit(Node)
func (s *TypedAstSerializer) Visit(n rhino.Node) *AstNode {
	return s.visitArena(n)
}

func (s *TypedAstSerializer) visitArena(n rhino.Node) *AstNode {
	jutil.CheckState(n.GetLength() >= 0)
	currentLine := n.GetLineno()
	currentColumn := n.GetCharno()
	m := s.allocAstNode()
	m.relativeLine_ = int32(currentLine - s.previousLine)
	m.relativeColumn_ = int32(currentColumn - s.previousColumn)
	s.previousLine = currentLine
	s.previousColumn = currentColumn

	if s.colorTypeSerializer != nil {
		if color := n.GetColor(); color != nil {
			m.type_ = int32(s.colorTypeSerializer.colorSerializer.AddColor(color))
			m.bitField0_ |= 1 << 0
		}
	} else if s.jsTypeSerializer != nil {
		s.jsTypeSerializer.addTypeToNode(n, m)
	} else if s.typeSerializer != nil {
		b := AstNode_Builder{m: m}
		s.typeSerializer.AddTypeForNode(n, &b)
	}

	if jsdoc := n.GetJSDocInfo(); jsdoc != nil {
		var doc OptimizationJsdoc
		if s.jsdocKindScratch == nil {
			s.jsdocKindScratch = make([]int32, 0, 8)
		}
		hasOptJsdoc := jsDocSerializer_fillJsdocProto(jsdoc, s.stringPool, &doc, s.jsdocKindScratch[:0])
		s.jsdocKindScratch = doc.kind_[:0]
		if hasOptJsdoc {
			m.jsdoc_ = s.allocJsdoc(doc)
		}
	}
	m.kind_ = s.KindTranslator(n).GetNumber()
	s.valueTranslatorToNode(m, n)
	m.booleanProperties_ = n.SerializeProperties()
	sourceFile := s.GetSourceFilePointer(n)
	m.sourceFile_ = int32(sourceFile)

	firstChild := n.GetFirstChild()
	shadowedCode := n.GetClosureUnawareShadow()
	if !firstChild.IsNil() || !shadowedCode.IsNil() {
		childrenStart := len(s.childStack)
		for child := firstChild; !child.IsNil(); child = child.GetNext() {
			s.childStack = append(s.childStack, s.visitArena(child))
		}
		if !shadowedCode.IsNil() {
			// For Closure shadow hosts, the ASTNode will get a boolean property bit set
			// (CLOSURE_UNAWARE_SHADOW) that indicates this child ASTNode is not a normal child, but is
			// the contents of the shadow.
			// Shadow roots are structured as
			// ROOT -> SCRIPT -> EXPR_RESULT -> FUNCTION
			// We avoid including the ROOT -> SCRIPT -> EXPR_RESULT structure in the serialization as it
			// is all synthetic code that would unnecessarily bloat the TypedAST and is instead recreated
			// upon deserialization.
			// The child ASTNode is just the FUNCTION.
			s.childStack = append(s.childStack, s.visitArena(shadowedCode.GetFirstFirstChild().GetFirstChild()))
		}
		m.child_ = s.allocChildren(s.childStack[childrenStart:])
		s.childStack = s.childStack[:childrenStart]
	}

	if sourceFile != 0 {
		// ArrayDeque.removeLast() throws NoSuchElementException on an empty deque.
		if len(s.subtreeSourceFiles) == 0 {
			panic(jutil.NewNoSuchElementException())
		}
		s.subtreeSourceFiles = s.subtreeSourceFiles[:len(s.subtreeSourceFiles)-1]
	}
	if origID := n.GetOriginalNameID(); origID != 0 {
		m.originalNamePointer_ = int32(s.stringPool.PutID(origID))
	}

	size := m.scalarsSerializedSize()
	for _, c := range m.child_ {
		cs := c.memoizedSize
		size += 1 + protoVarintLen(uint64(cs)) + cs
	}
	if m.jsdoc_ != nil {
		js := m.jsdoc_.memoizedSize
		size += 1 + protoVarintLen(uint64(js)) + js
	}
	if m.valueCase_ == 8 {
		ts := m.templateStringValue_.memoizedSize
		size += 1 + protoVarintLen(uint64(ts)) + ts
	}
	m.memoizedSize = size

	return m
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
	if idx, ok := s.sourceFilePointers[sourceFile]; ok {
		return idx
	}
	idx := 1 + len(s.sourceFileList)
	s.sourceFileList = append(s.sourceFileList, sourceFile)
	s.sourceFilePointers[sourceFile] = idx
	return idx
}

// Java: TypedAstSerializer.addType(Node, AstNode.Builder)
func (s *TypedAstSerializer) AddType(n rhino.Node, builder *AstNode_Builder) {
	s.typeSerializer.AddTypeForNode(n, builder)
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

func (s *TypedAstSerializer) valueTranslatorToNode(m *AstNode, n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_GETPROP,
		rhino.Token_OPTCHAIN_GETPROP:
		ptr := s.stringPool.PutID(n.GetStringID())
		if s.recordUsedProps {
			s.markUsedProp(ptr)
		}
		m.stringValuePointer_ = int32(ptr)
		m.valueCase_ = 3
		return
	case rhino.Token_STRING_KEY,
		rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_MEMBER_FIELD_DEF,
		rhino.Token_GETTER_DEF,
		rhino.Token_SETTER_DEF:
		ptr := s.stringPool.PutID(n.GetStringID())
		if s.recordUsedProps && !n.IsQuotedStringKey() {
			s.markUsedProp(ptr)
		}
		m.stringValuePointer_ = int32(ptr)
		m.valueCase_ = 3
		return
	case rhino.Token_NAME,
		rhino.Token_STRINGLIT,
		rhino.Token_LABEL_NAME,
		rhino.Token_IMPORT_STAR:
		m.stringValuePointer_ = int32(s.stringPool.PutID(n.GetStringID()))
		m.valueCase_ = 3
		return
	case rhino.Token_TEMPLATELIT_STRING:
		// The raw string is put in the pool before the cooked one (Java's
		// evaluation order).
		rawPtr := int32(s.stringPool.PutID(n.GetRawStringID()))
		cookedPtr := int32(-1)
		if cookedID := n.GetCookedStringID(); cookedID != 0 {
			cookedPtr = int32(s.stringPool.PutID(cookedID))
		}
		m.templateStringValue_ = s.allocTemplateStringValue(rawPtr, cookedPtr)
		m.valueCase_ = 8
		return
	case rhino.Token_NUMBER:
		m.doubleValue_ = n.GetDouble()
		m.valueCase_ = 4
		return
	case rhino.Token_BIGINT:
		m.stringValuePointer_ = int32(s.stringPool.Put(n.GetBigInt().String()))
		m.valueCase_ = 3
		return
	default:
		// No value
		return
	}
}

// Java: TypedAstSerializer.valueTranslator(AstNode.Builder, Node)
func (s *TypedAstSerializer) ValueTranslator(builder *AstNode_Builder, n rhino.Node) {
	switch n.GetToken() {
	case rhino.Token_GETPROP,
		rhino.Token_OPTCHAIN_GETPROP,
		rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_MEMBER_FIELD_DEF,
		rhino.Token_NAME,
		rhino.Token_STRINGLIT,
		rhino.Token_STRING_KEY,
		rhino.Token_GETTER_DEF,
		rhino.Token_SETTER_DEF,
		rhino.Token_LABEL_NAME,
		rhino.Token_IMPORT_STAR:
		builder.SetStringValuePointer(int32(s.stringPool.PutID(n.GetStringID())))
		return
	case rhino.Token_TEMPLATELIT_STRING:
		// The raw string is put in the pool before the cooked one (Java's
		// evaluation order).
		templateStringValue := TemplateStringValue_NewBuilder().
			SetRawStringPointer(int32(s.stringPool.PutID(n.GetRawStringID())))
		cookedStringPointer := -1
		if cookedID := n.GetCookedStringID(); cookedID != 0 {
			cookedStringPointer = s.stringPool.PutID(cookedID)
		}
		builder.SetTemplateStringValue(
			templateStringValue.
				SetCookedStringPointer(int32(cookedStringPointer)).
				Build())
		return
	case rhino.Token_NUMBER:
		builder.SetDoubleValue(n.GetDouble())
		return
	case rhino.Token_BIGINT:
		builder.SetStringValuePointer(int32(s.stringPool.Put(n.GetBigInt().String())))
		return
	default:
		// No value
		return
	}
}

// Java: TypedAstSerializer.setOriginalName(AstNode.Builder, Node)
func (s *TypedAstSerializer) SetOriginalName(builder *AstNode_Builder, n rhino.Node) {
	origID := n.GetOriginalNameID()
	if origID == 0 {
		builder.SetOriginalNamePointer(0) // equivalent to 'not set'
		return
	}

	builder.SetOriginalNamePointer(int32(s.stringPool.PutID(origID)))
}

// Java: TypedAstSerializer.kindTranslator(Node)
func (s *TypedAstSerializer) KindTranslator(n rhino.Node) NodeKind {
	switch n.GetToken() {
	case rhino.Token_SCRIPT:
		return NodeKind_SOURCE_FILE
	case rhino.Token_NUMBER:
		return NodeKind_NUMBER_LITERAL
	case rhino.Token_STRINGLIT:
		return NodeKind_STRING_LITERAL
	case rhino.Token_BIGINT:
		return NodeKind_BIGINT_LITERAL
	case rhino.Token_REGEXP:
		return NodeKind_REGEX_LITERAL
	case rhino.Token_FALSE:
		return NodeKind_FALSE
	case rhino.Token_TRUE:
		return NodeKind_TRUE
	case rhino.Token_NULL:
		return NodeKind_NULL
	case rhino.Token_THIS:
		return NodeKind_THIS
	case rhino.Token_VOID:
		return NodeKind_VOID
	case rhino.Token_ARRAYLIT:
		return NodeKind_ARRAY_LITERAL
	case rhino.Token_OBJECTLIT:
		return NodeKind_OBJECT_LITERAL
	case rhino.Token_NAME:
		return NodeKind_IDENTIFIER
	case rhino.Token_ASSIGN:
		return NodeKind_ASSIGNMENT
	case rhino.Token_CALL:
		return NodeKind_CALL
	case rhino.Token_NEW:
		return NodeKind_NEW
	case rhino.Token_GETPROP:
		return NodeKind_PROPERTY_ACCESS
	case rhino.Token_GETELEM:
		return NodeKind_ELEMENT_ACCESS
	case rhino.Token_COMMA:
		return NodeKind_COMMA
	case rhino.Token_OR:
		return NodeKind_BOOLEAN_OR
	case rhino.Token_AND:
		return NodeKind_BOOLEAN_AND
	case rhino.Token_HOOK:
		return NodeKind_HOOK
	case rhino.Token_EQ:
		return NodeKind_EQUAL
	case rhino.Token_NE:
		return NodeKind_NOT_EQUAL
	case rhino.Token_LT:
		return NodeKind_LESS_THAN
	case rhino.Token_LE:
		return NodeKind_LESS_THAN_EQUAL
	case rhino.Token_GT:
		return NodeKind_GREATER_THAN
	case rhino.Token_GE:
		return NodeKind_GREATER_THAN_EQUAL
	case rhino.Token_SHEQ:
		return NodeKind_TRIPLE_EQUAL
	case rhino.Token_SHNE:
		return NodeKind_NOT_TRIPLE_EQUAL
	case rhino.Token_NOT:
		return NodeKind_NOT
	case rhino.Token_POS:
		return NodeKind_POSITIVE
	case rhino.Token_NEG:
		return NodeKind_NEGATIVE
	case rhino.Token_TYPEOF:
		return NodeKind_TYPEOF
	case rhino.Token_INSTANCEOF:
		return NodeKind_INSTANCEOF
	case rhino.Token_IN:
		return NodeKind_IN
	case rhino.Token_ADD:
		return NodeKind_ADD
	case rhino.Token_SUB:
		return NodeKind_SUBTRACT
	case rhino.Token_MUL:
		return NodeKind_MULTIPLY
	case rhino.Token_DIV:
		return NodeKind_DIVIDE
	case rhino.Token_MOD:
		return NodeKind_MODULO
	case rhino.Token_EXPONENT:
		return NodeKind_EXPONENT
	case rhino.Token_BITNOT:
		return NodeKind_BITWISE_NOT
	case rhino.Token_BITOR:
		return NodeKind_BITWISE_OR
	case rhino.Token_BITAND:
		return NodeKind_BITWISE_AND
	case rhino.Token_BITXOR:
		return NodeKind_BITWISE_XOR
	case rhino.Token_LSH:
		return NodeKind_LEFT_SHIFT
	case rhino.Token_RSH:
		return NodeKind_RIGHT_SHIFT
	case rhino.Token_URSH:
		return NodeKind_UNSIGNED_RIGHT_SHIFT
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
	case rhino.Token_ASSIGN_ADD:
		return NodeKind_ASSIGN_ADD
	case rhino.Token_ASSIGN_SUB:
		return NodeKind_ASSIGN_SUBTRACT
	case rhino.Token_ASSIGN_MUL:
		return NodeKind_ASSIGN_MULTIPLY
	case rhino.Token_ASSIGN_DIV:
		return NodeKind_ASSIGN_DIVIDE
	case rhino.Token_ASSIGN_MOD:
		return NodeKind_ASSIGN_MODULO
	case rhino.Token_ASSIGN_EXPONENT:
		return NodeKind_ASSIGN_EXPONENT
	case rhino.Token_ASSIGN_BITOR:
		return NodeKind_ASSIGN_BITWISE_OR
	case rhino.Token_ASSIGN_BITAND:
		return NodeKind_ASSIGN_BITWISE_AND
	case rhino.Token_ASSIGN_BITXOR:
		return NodeKind_ASSIGN_BITWISE_XOR
	case rhino.Token_ASSIGN_LSH:
		return NodeKind_ASSIGN_LEFT_SHIFT
	case rhino.Token_ASSIGN_RSH:
		return NodeKind_ASSIGN_RIGHT_SHIFT
	case rhino.Token_ASSIGN_URSH:
		return NodeKind_ASSIGN_UNSIGNED_RIGHT_SHIFT
	case rhino.Token_YIELD:
		return NodeKind_YIELD
	case rhino.Token_AWAIT:
		return NodeKind_AWAIT
	case rhino.Token_DELPROP:
		return NodeKind_DELETE
	case rhino.Token_TAGGED_TEMPLATELIT:
		return NodeKind_TAGGED_TEMPLATELIT
	case rhino.Token_TEMPLATELIT:
		return NodeKind_TEMPLATELIT
	case rhino.Token_TEMPLATELIT_SUB:
		return NodeKind_TEMPLATELIT_SUB
	case rhino.Token_TEMPLATELIT_STRING:
		return NodeKind_TEMPLATELIT_STRING
	case rhino.Token_NEW_TARGET:
		return NodeKind_NEW_TARGET
	case rhino.Token_COMPUTED_PROP:
		return NodeKind_COMPUTED_PROP
	case rhino.Token_IMPORT_META:
		return NodeKind_IMPORT_META
	case rhino.Token_OPTCHAIN_GETPROP:
		return NodeKind_OPTCHAIN_PROPERTY_ACCESS
	case rhino.Token_OPTCHAIN_CALL:
		return NodeKind_OPTCHAIN_CALL
	case rhino.Token_OPTCHAIN_GETELEM:
		return NodeKind_OPTCHAIN_ELEMENT_ACCESS
	case rhino.Token_COALESCE:
		return NodeKind_COALESCE
	case rhino.Token_DYNAMIC_IMPORT:
		return NodeKind_DYNAMIC_IMPORT
	case rhino.Token_ASSIGN_OR:
		return NodeKind_ASSIGN_OR
	case rhino.Token_ASSIGN_AND:
		return NodeKind_ASSIGN_AND
	case rhino.Token_ASSIGN_COALESCE:
		return NodeKind_ASSIGN_COALESCE
	case rhino.Token_EXPR_RESULT:
		return NodeKind_EXPRESSION_STATEMENT
	case rhino.Token_BREAK:
		return NodeKind_BREAK_STATEMENT
	case rhino.Token_CONTINUE:
		return NodeKind_CONTINUE_STATEMENT
	case rhino.Token_DEBUGGER:
		return NodeKind_DEBUGGER_STATEMENT
	case rhino.Token_DO:
		return NodeKind_DO_STATEMENT
	case rhino.Token_FOR:
		return NodeKind_FOR_STATEMENT
	case rhino.Token_FOR_IN:
		return NodeKind_FOR_IN_STATEMENT
	case rhino.Token_FOR_OF:
		return NodeKind_FOR_OF_STATEMENT
	case rhino.Token_FOR_AWAIT_OF:
		return NodeKind_FOR_AWAIT_OF_STATEMENT
	case rhino.Token_IF:
		return NodeKind_IF_STATEMENT
	case rhino.Token_RETURN:
		return NodeKind_RETURN_STATEMENT
	case rhino.Token_SWITCH:
		return NodeKind_SWITCH_STATEMENT
	case rhino.Token_THROW:
		return NodeKind_THROW_STATEMENT
	case rhino.Token_TRY:
		return NodeKind_TRY_STATEMENT
	case rhino.Token_WHILE:
		return NodeKind_WHILE_STATEMENT
	case rhino.Token_EMPTY:
		return NodeKind_EMPTY
	case rhino.Token_WITH:
		return NodeKind_WITH
	case rhino.Token_IMPORT:
		return NodeKind_IMPORT
	case rhino.Token_EXPORT:
		return NodeKind_EXPORT
	case rhino.Token_VAR:
		return NodeKind_VAR_DECLARATION
	case rhino.Token_CONST:
		return NodeKind_CONST_DECLARATION
	case rhino.Token_LET:
		return NodeKind_LET_DECLARATION
	case rhino.Token_FUNCTION:
		return NodeKind_FUNCTION_LITERAL
	case rhino.Token_CLASS:
		return NodeKind_CLASS_LITERAL
	case rhino.Token_BLOCK:
		return NodeKind_BLOCK
	case rhino.Token_LABEL:
		return NodeKind_LABELED_STATEMENT
	case rhino.Token_LABEL_NAME:
		return NodeKind_LABELED_NAME
	case rhino.Token_CLASS_MEMBERS:
		return NodeKind_CLASS_MEMBERS
	case rhino.Token_MEMBER_FUNCTION_DEF:
		return NodeKind_METHOD_DECLARATION
	case rhino.Token_MEMBER_FIELD_DEF:
		return NodeKind_FIELD_DECLARATION
	case rhino.Token_COMPUTED_FIELD_DEF:
		return NodeKind_COMPUTED_PROP_FIELD
	case rhino.Token_PARAM_LIST:
		return NodeKind_PARAMETER_LIST
	case rhino.Token_STRING_KEY:
		if n.IsQuotedStringKey() {
			return NodeKind_QUOTED_STRING_KEY
		}
		return NodeKind_RENAMABLE_STRING_KEY
	case rhino.Token_CASE:
		return NodeKind_CASE
	case rhino.Token_DEFAULT_CASE:
		return NodeKind_DEFAULT_CASE
	case rhino.Token_CATCH:
		return NodeKind_CATCH
	case rhino.Token_SUPER:
		return NodeKind_SUPER
	case rhino.Token_ARRAY_PATTERN:
		return NodeKind_ARRAY_PATTERN
	case rhino.Token_OBJECT_PATTERN:
		return NodeKind_OBJECT_PATTERN
	case rhino.Token_DESTRUCTURING_LHS:
		return NodeKind_DESTRUCTURING_LHS
	case rhino.Token_DEFAULT_VALUE:
		return NodeKind_DEFAULT_VALUE
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
	case rhino.Token_IMPORT_SPECS:
		return NodeKind_IMPORT_SPECS
	case rhino.Token_IMPORT_SPEC:
		return NodeKind_IMPORT_SPEC
	case rhino.Token_IMPORT_STAR:
		return NodeKind_IMPORT_STAR
	case rhino.Token_EXPORT_SPECS:
		return NodeKind_EXPORT_SPECS
	case rhino.Token_EXPORT_SPEC:
		return NodeKind_EXPORT_SPEC
	case rhino.Token_MODULE_BODY:
		return NodeKind_MODULE_BODY
	case rhino.Token_ITER_REST:
		return NodeKind_ITER_REST
	case rhino.Token_ITER_SPREAD:
		return NodeKind_ITER_SPREAD
	case rhino.Token_OBJECT_REST:
		return NodeKind_OBJECT_REST
	case rhino.Token_OBJECT_SPREAD:
		return NodeKind_OBJECT_SPREAD
	case rhino.Token_SWITCH_BODY:
		return NodeKind_SWITCH_BODY
	// Explicitly unsupported token types. Not serialized.
	case rhino.Token_ROOT,
		// TS type tokens
		rhino.Token_STRING_TYPE,
		rhino.Token_BOOLEAN_TYPE,
		rhino.Token_NUMBER_TYPE,
		rhino.Token_FUNCTION_TYPE,
		rhino.Token_PARAMETERIZED_TYPE,
		rhino.Token_UNION_TYPE,
		rhino.Token_ANY_TYPE,
		rhino.Token_NULLABLE_TYPE,
		rhino.Token_VOID_TYPE,
		rhino.Token_REST_PARAMETER_TYPE,
		rhino.Token_NAMED_TYPE,
		rhino.Token_RECORD_TYPE,
		rhino.Token_UNDEFINED_TYPE,
		rhino.Token_OPTIONAL_PARAMETER,
		rhino.Token_ARRAY_TYPE,
		rhino.Token_GENERIC_TYPE,
		rhino.Token_GENERIC_TYPE_LIST,
		rhino.Token_INTERFACE,
		rhino.Token_INTERFACE_EXTENDS,
		rhino.Token_INTERFACE_MEMBERS,
		rhino.Token_MEMBER_VARIABLE_DEF,
		rhino.Token_ENUM,
		rhino.Token_ENUM_MEMBERS,
		rhino.Token_IMPLEMENTS,
		rhino.Token_TYPE_ALIAS,
		rhino.Token_DECLARE,
		rhino.Token_INDEX_SIGNATURE,
		rhino.Token_CALL_SIGNATURE,
		rhino.Token_NAMESPACE,
		rhino.Token_NAMESPACE_ELEMENTS,
		// JSDoc tokens
		rhino.Token_ANNOTATION,
		rhino.Token_PIPE,
		rhino.Token_STAR,
		rhino.Token_EOC,
		rhino.Token_QMARK,
		rhino.Token_BANG,
		rhino.Token_EQUALS,
		rhino.Token_LB,
		rhino.Token_LC,
		rhino.Token_COLON,
		// Closure-type-system-specific token
		rhino.Token_CAST,
		// Unused tokens
		rhino.Token_PLACEHOLDER1,
		rhino.Token_PLACEHOLDER2,
		rhino.Token_PLACEHOLDER3:
	}
	panic(jutil.NewIllegalStateExceptionMsg("Unserializable token for node: " + n.String()))
}

// TypedAstSerializer_TypeSerializer is used to provide TypePointers for
// serializing Nodes and to generate the TypePool.
//
// Java: interface TypedAstSerializer.TypeSerializer
type TypedAstSerializer_TypeSerializer interface {
	// AddTypeForNode adds, if appropriate for `node`, an `Integer` to
	// `astNodeBuilder`.
	//
	// Java: TypedAstSerializer.TypeSerializer.addTypeForNode(Node, AstNode.Builder)
	AddTypeForNode(node rhino.Node, astNodeBuilder *AstNode_Builder)

	// GenerateTypePool returns a `TypePool` containing the types used by
	// `addTypeForNode()`.
	//
	// Java: TypedAstSerializer.TypeSerializer.generateTypePool()
	GenerateTypePool() *TypePool
}

// TypedAstSerializer_NoOpTypeSerializer is used when type checking has not
// been done.
//
// Java: private static class TypedAstSerializer.NoOpTypeSerializer
type TypedAstSerializer_NoOpTypeSerializer struct {
	_ byte // Go-only: language-adapter — distinct identity
}

// Java: TypedAstSerializer.NoOpTypeSerializer()
func NewTypedAstSerializer_NoOpTypeSerializer() *TypedAstSerializer_NoOpTypeSerializer {
	return &TypedAstSerializer_NoOpTypeSerializer{}
}

// Java: TypedAstSerializer.NoOpTypeSerializer.addTypeForNode(Node, AstNode.Builder)
func (*TypedAstSerializer_NoOpTypeSerializer) AddTypeForNode(node rhino.Node, astNodeBuilder *AstNode_Builder) {
	// Do nothing.
}

// Java: TypedAstSerializer.NoOpTypeSerializer.generateTypePool()
func (*TypedAstSerializer_NoOpTypeSerializer) GenerateTypePool() *TypePool {
	return TypePool_GetDefaultInstance() // empty TypePool
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

// Java: TypedAstSerializer.JSTypeSerializer.addTypeForNode(Node, AstNode.Builder)
func (s *TypedAstSerializer_JSTypeSerializer) AddTypeForNode(node rhino.Node, astNodeBuilder *AstNode_Builder) {
	s.addTypeToNode(node, astNodeBuilder.m)
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
	usedPropertyNames := make(map[string]struct{}, 256)
	var lastProp string
	hasLastProp := false
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
			prop := n.GetString()
			if !hasLastProp || prop != lastProp {
				lastProp = prop
				hasLastProp = true
				usedPropertyNames[prop] = struct{}{}
			}
		case rhino.Token_STRING_KEY,
			rhino.Token_MEMBER_FUNCTION_DEF,
			rhino.Token_MEMBER_FIELD_DEF,
			rhino.Token_GETTER_DEF,
			rhino.Token_SETTER_DEF:
			if !n.IsQuotedStringKey() {
				prop := n.GetString()
				if !hasLastProp || prop != lastProp {
					lastProp = prop
					hasLastProp = true
					usedPropertyNames[prop] = struct{}{}
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
			_, ok := usedPropertyNames[prop]
			return ok
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

// Java: TypedAstSerializer.ColorTypeSerializer.addTypeForNode(Node, AstNode.Builder)
func (s *TypedAstSerializer_ColorTypeSerializer) AddTypeForNode(node rhino.Node, astNodeBuilder *AstNode_Builder) {
	color := node.GetColor()
	if color != nil {
		astNodeBuilder.SetType(int32(s.colorSerializer.AddColor(color)))
	}
}

// Java: TypedAstSerializer.ColorTypeSerializer.generateTypePool()
func (s *TypedAstSerializer_ColorTypeSerializer) GenerateTypePool() *TypePool {
	mismatchLocationsForDebugging := s.colorRegistry.GetMismatchLocationsForDebugging()
	return s.colorSerializer.GenerateTypePool(
		s.colorRegistry.GetDisambiguationSupertypes,
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
