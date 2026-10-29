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

type serializerFixupKind uint8

const (
	fixupKindStringID serializerFixupKind = iota
	fixupKindStringIDUsedProp
	fixupKindColor
	fixupKindSourceFile
)

type serializerFixup struct {
	dst   *int32
	color *colors.Color
	sf    *SourceFile
	strID uint32
	kind  serializerFixupKind
}

type typedAstScriptArena struct {
	TypedAstSerializer
	rootProto      *AstNode
	lazyAst        LazyAst
	fixups         []serializerFixup
	jsdocFixups    []*OptimizationJsdoc
	templateFixups []*TemplateStringValue
	postOrderNodes []*AstNode
}

func (a *typedAstScriptArena) resetForScript() {
	a.previousColumn = 0
	a.previousLine = 0
	a.resetScriptSlabs()
	a.subtreeSourceFiles = a.subtreeSourceFiles[:0]
	a.rootProto = nil
	a.lazyAst = LazyAst{frozen: true, memoizedSize: -1}
	a.fixups = a.fixups[:0]
	a.jsdocFixups = a.jsdocFixups[:0]
	a.templateFixups = a.templateFixups[:0]
	a.postOrderNodes = a.postOrderNodes[:0]
}

func (a *typedAstScriptArena) buildScriptPassA(script rhino.Node, parent *TypedAstSerializer) {
	jutil.CheckState(script.IsScript())
	a.resetForScript()

	var scriptSF *SourceFile
	if ssf := script.GetStaticSourceFile(); ssf != nil {
		scriptSF = ssf.(*SourceFile)
	}
	if scriptSF == nil {
		if !a.peekLastSubtreeSourceFile().IsExtern() {
			jutil.CheckStateMsg2(
				false,
				"Unexpected null SourceFile for node %s with parent %s",
				script.ToStringTree(),
				script.GetParent())
		}
	} else {
		a.subtreeSourceFiles = append(a.subtreeSourceFiles, scriptSF)
		a.fixups = append(a.fixups, serializerFixup{
			kind: fixupKindSourceFile,
			sf:   scriptSF,
			dst:  &a.lazyAst.sourceFile_,
		})
	}

	a.rootProto = a.visitArenaParallel(script, parent)
	a.subtreeSourceFiles = a.subtreeSourceFiles[:0]

	if encodedSourceMap := parent.compiler.GetBase64SourceMapContents(*script.GetSourceFileName()); encodedSourceMap != nil {
		a.lazyAst.sourceMappingUrl_ = *encodedSourceMap
	}
}

func (a *typedAstScriptArena) visitArenaParallel(n rhino.Node, parent *TypedAstSerializer) *AstNode {
	jutil.CheckState(n.GetLength() >= 0)
	currentLine := n.GetLineno()
	currentColumn := n.GetCharno()
	m := a.allocAstNode()
	m.relativeLine_ = int32(currentLine - a.previousLine)
	m.relativeColumn_ = int32(currentColumn - a.previousColumn)
	a.previousLine = currentLine
	a.previousColumn = currentColumn

	if parent.colorTypeSerializer != nil {
		if color := n.GetColor(); color != nil {
			m.bitField0_ |= 1 << 0
			a.fixups = append(a.fixups, serializerFixup{
				kind:  fixupKindColor,
				color: color,
				dst:   &m.type_,
			})
		}
	} else if parent.typeSerializer != nil {
		b := AstNode_Builder{m: m}
		parent.typeSerializer.AddTypeForNode(n, &b)
	}

	if jsdoc := n.GetJSDocInfo(); jsdoc != nil {
		var doc OptimizationJsdoc
		if a.jsdocKindScratch == nil {
			a.jsdocKindScratch = make([]int32, 0, 8)
		}
		doc.kind_ = jsDocSerializer_fillJsdocTagsOnly(jsdoc, a.jsdocKindScratch[:0])
		a.jsdocKindScratch = doc.kind_[:0]

		lic, desc, alt, mean, mode, hasMode := jsDocSerializer_stringFields(jsdoc)
		hasStr := lic != "" || desc != "" || alt != "" || mean != "" || hasMode
		if len(doc.kind_) != 0 || hasStr {
			allocatedDoc := a.allocJsdoc(doc)
			m.jsdoc_ = allocatedDoc
			if hasStr {
				a.jsdocFixups = append(a.jsdocFixups, allocatedDoc)
				if lic != "" {
					a.fixups = append(a.fixups, serializerFixup{
						kind:  fixupKindStringID,
						strID: rhino.RhinoStringPool_InternID(lic),
						dst:   &allocatedDoc.licenseTextPointer_,
					})
				}
				if desc != "" {
					a.fixups = append(a.fixups, serializerFixup{
						kind:  fixupKindStringID,
						strID: rhino.RhinoStringPool_InternID(desc),
						dst:   &allocatedDoc.descriptionPointer_,
					})
				}
				if alt != "" {
					a.fixups = append(a.fixups, serializerFixup{
						kind:  fixupKindStringID,
						strID: rhino.RhinoStringPool_InternID(alt),
						dst:   &allocatedDoc.alternateMessageIdPointer_,
					})
				}
				if mean != "" {
					a.fixups = append(a.fixups, serializerFixup{
						kind:  fixupKindStringID,
						strID: rhino.RhinoStringPool_InternID(mean),
						dst:   &allocatedDoc.meaningPointer_,
					})
				}
				if hasMode {
					a.fixups = append(a.fixups, serializerFixup{
						kind:  fixupKindStringID,
						strID: rhino.RhinoStringPool_InternID(mode.String()),
						dst:   &allocatedDoc.closureUnawareModePointer_,
					})
				}
			}
		}
	}

	m.kind_ = a.KindTranslator(n).GetNumber()
	a.valueTranslatorParallel(m, n, parent.recordUsedProps)
	m.booleanProperties_ = n.SerializeProperties()

	pushedSF := false
	var sf *SourceFile
	if ssf := n.GetStaticSourceFile(); ssf != nil {
		sf = ssf.(*SourceFile)
	}
	if sf == nil {
		if !a.peekLastSubtreeSourceFile().IsExtern() {
			jutil.CheckStateMsg2(
				false,
				"Unexpected null SourceFile for node %s with parent %s",
				n.ToStringTree(),
				n.GetParent())
		}
	} else if sf != a.peekLastSubtreeSourceFile() {
		a.subtreeSourceFiles = append(a.subtreeSourceFiles, sf)
		pushedSF = true
		a.fixups = append(a.fixups, serializerFixup{
			kind: fixupKindSourceFile,
			sf:   sf,
			dst:  &m.sourceFile_,
		})
	}

	firstChild := n.GetFirstChild()
	shadowedCode := n.GetClosureUnawareShadow()
	if !firstChild.IsNil() || !shadowedCode.IsNil() {
		childrenStart := len(a.childStack)
		for child := firstChild; !child.IsNil(); child = child.GetNext() {
			a.childStack = append(a.childStack, a.visitArenaParallel(child, parent))
		}
		if !shadowedCode.IsNil() {
			a.childStack = append(a.childStack, a.visitArenaParallel(shadowedCode.GetFirstFirstChild().GetFirstChild(), parent))
		}
		m.child_ = a.allocChildren(a.childStack[childrenStart:])
		a.childStack = a.childStack[:childrenStart]
	}

	if pushedSF {
		if len(a.subtreeSourceFiles) == 0 {
			panic(jutil.NewNoSuchElementException())
		}
		a.subtreeSourceFiles = a.subtreeSourceFiles[:len(a.subtreeSourceFiles)-1]
	}
	if origID := n.GetOriginalNameID(); origID != 0 {
		a.fixups = append(a.fixups, serializerFixup{
			kind:  fixupKindStringID,
			strID: origID,
			dst:   &m.originalNamePointer_,
		})
	}

	a.postOrderNodes = append(a.postOrderNodes, m)
	return m
}

func (a *typedAstScriptArena) valueTranslatorParallel(m *AstNode, n rhino.Node, recordUsedProps bool) {
	switch n.GetToken() {
	case rhino.Token_GETPROP,
		rhino.Token_OPTCHAIN_GETPROP:
		kind := fixupKindStringID
		if recordUsedProps {
			kind = fixupKindStringIDUsedProp
		}
		a.fixups = append(a.fixups, serializerFixup{
			kind:  kind,
			strID: n.GetStringID(),
			dst:   &m.stringValuePointer_,
		})
		m.valueCase_ = 3
	case rhino.Token_STRING_KEY,
		rhino.Token_MEMBER_FUNCTION_DEF,
		rhino.Token_MEMBER_FIELD_DEF,
		rhino.Token_GETTER_DEF,
		rhino.Token_SETTER_DEF:
		kind := fixupKindStringID
		if recordUsedProps && !n.IsQuotedStringKey() {
			kind = fixupKindStringIDUsedProp
		}
		a.fixups = append(a.fixups, serializerFixup{
			kind:  kind,
			strID: n.GetStringID(),
			dst:   &m.stringValuePointer_,
		})
		m.valueCase_ = 3
	case rhino.Token_NAME,
		rhino.Token_STRINGLIT,
		rhino.Token_LABEL_NAME,
		rhino.Token_IMPORT_STAR:
		a.fixups = append(a.fixups, serializerFixup{
			kind:  fixupKindStringID,
			strID: n.GetStringID(),
			dst:   &m.stringValuePointer_,
		})
		m.valueCase_ = 3
	case rhino.Token_TEMPLATELIT_STRING:
		tmpl := a.allocTemplateStringValue(0, -1)
		m.templateStringValue_ = tmpl
		m.valueCase_ = 8
		a.templateFixups = append(a.templateFixups, tmpl)
		a.fixups = append(a.fixups, serializerFixup{
			kind:  fixupKindStringID,
			strID: n.GetRawStringID(),
			dst:   &tmpl.rawStringPointer_,
		})
		if cookedID := n.GetCookedStringID(); cookedID != 0 {
			a.fixups = append(a.fixups, serializerFixup{
				kind:  fixupKindStringID,
				strID: cookedID,
				dst:   &tmpl.cookedStringPointer_,
			})
		}
	case rhino.Token_NUMBER:
		m.doubleValue_ = n.GetDouble()
		m.valueCase_ = 4
	case rhino.Token_BIGINT:
		a.fixups = append(a.fixups, serializerFixup{
			kind:  fixupKindStringID,
			strID: rhino.RhinoStringPool_InternID(n.GetBigInt().String()),
			dst:   &m.stringValuePointer_,
		})
		m.valueCase_ = 3
	}
}

func (a *typedAstScriptArena) marshalScriptPassC(byteSlab *[]byte) {
	for _, doc := range a.jsdocFixups {
		doc.memoizedSize = doc.scalarsSerializedSize()
	}
	for _, tmpl := range a.templateFixups {
		tmpl.memoizedSize = tmpl.scalarsSerializedSize()
	}
	for _, m := range a.postOrderNodes {
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
	}

	sz := a.rootProto.serializedSize()
	var scriptBytes []byte
	if sz <= 262144 {
		slab := *byteSlab
		if len(slab) < sz {
			slab = make([]byte, max(4*1024*1024, sz))
		}
		buf := slab[:0:sz]
		*byteSlab = slab[sz:]
		scriptBytes = a.rootProto.appendTo(buf)
	} else {
		scriptBytes = a.rootProto.ToByteString()
	}
	a.lazyAst.script_ = scriptBytes
	if len(a.lazyAst.sourceMappingUrl_) == 0 {
		ls := 0
		if len(scriptBytes) != 0 {
			ls += 1 + protoVarintLen(uint64(len(scriptBytes))) + len(scriptBytes)
		}
		if a.lazyAst.sourceFile_ != 0 {
			ls += 1 + protoVarintLen(uint64(uint32(a.lazyAst.sourceFile_)))
		}
		a.lazyAst.memoizedSize = ls
	} else {
		a.lazyAst.serializedSize()
	}
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
