// Port of com/google/javascript/jscomp/serialization/ScriptNodeDeserializer.java @ closure-compiler v20261021

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
	"bytes"
	"math"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/parsing"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ScriptNodeDeserializer deserializes an AstNode-tree representing a SCRIPT
// into a Node-tree.
//
// This process depends on other information from the TypedAST format, but
// the output it limited to only a single SCRIPT. The other deserialized
// content must be provided beforehand.
type ScriptNodeDeserializer struct {
	sourceFile       *SourceFile
	scriptBytes      []byte
	sourceMappingURL string
	colorPoolShard   *ColorPool_ShardView // Optional: nil when absent
	stringPool       *StringPool
	filePool         []*SourceFile
}

// Java: ScriptNodeDeserializer(LazyAst, StringPool, Optional<ColorPool.ShardView>, ImmutableList<SourceFile>)
func NewScriptNodeDeserializer(
	ast *LazyAst,
	stringPool *StringPool,
	colorPoolShard *ColorPool_ShardView,
	filePool []*SourceFile) *ScriptNodeDeserializer {
	scriptBytes := ast.GetScript()
	if len(scriptBytes) <= 256 {
		scriptBytes = bytes.Clone(scriptBytes)
	}
	return &ScriptNodeDeserializer{
		scriptBytes:      scriptBytes,
		sourceFile:       filePool[ast.GetSourceFile()-1],
		sourceMappingURL: ast.GetSourceMappingUrl(),
		colorPoolShard:   colorPoolShard,
		stringPool:       stringPool,
		filePool:         filePool,
	}
}

// Java: ScriptNodeDeserializer.deserializeNew()
func (d *ScriptNodeDeserializer) DeserializeNew() rhino.Node {
	return newScriptNodeDeserializer_Runner(d).Run()
}

// Java: ScriptNodeDeserializer.Runner
type scriptNodeDeserializer_Runner struct {
	owner_           *ScriptNodeDeserializer
	featureCollector *parsing.FeatureCollector
	previousLine     int
	previousColumn   int
	childBytes       [][]byte
	contextStack     []parsing.FeatureCollector_FeatureContext
	jsdoc            OptimizationJsdoc
	jsdocKindBuf     [16]int32
	jsdocBuilder     rhino.JSDocInfo_Builder
}

func newScriptNodeDeserializer_Runner(owner *ScriptNodeDeserializer) *scriptNodeDeserializer_Runner {
	return &scriptNodeDeserializer_Runner{owner_: owner, featureCollector: parsing.NewFeatureCollector()}
}

// Java: ScriptNodeDeserializer.Runner.run()
func (r *scriptNodeDeserializer_Runner) Run() rhino.Node {
	defer func() {
		if rec := recover(); rec != nil {
			t, ok := rec.(jutil.Throwable)
			if !ok {
				t = jutil.AsThrowable(rec)
			}
			if ex, ok := t.(jutil.IIOException); ok {
				panic(NewMalformedTypedAstExceptionWithCause(r.Owner().sourceFile, ex))
			}
			panic(rec)
		}
	}()
	r.contextStack = append(r.contextStack[:0], parsing.FeatureCollector_FeatureContext_NONE)

	tmpl := r.Owner().sourceFile.getOrCreateSourceInfoTemplate()

	scriptNode, _ := r.visitBytes(
		r.Owner().scriptBytes,
		0,
		tmpl)
	scriptNode.PutProp(rhino.Node_FEATURE_SET, r.featureCollector.AllFeatures())
	return scriptNode
}

// Java: ScriptNodeDeserializer.Runner.owner()
func (r *scriptNodeDeserializer_Runner) Owner() *ScriptNodeDeserializer {
	return r.owner_
}

func (r *scriptNodeDeserializer_Runner) visitBytes(
	data []byte,
	depth int,
	sourceFileTemplate rhino.Node,
) (rhino.Node, int32) {
	var (
		kindVal             int32
		stringValuePointer  int32
		doubleValue         float64
		templateRawPtr      int32
		templateCookedPtr   int32
		relativeLine        int32
		relativeColumn      int32
		originalNamePointer int32
		booleanProperties   int64
		typeVal             int32
		hasType             bool
		sourceFile          int32
		hasJsdoc            bool
	)
	childrenStart := len(r.childBytes)

	rd := protoReader{b: data, depth: depth, limit: math.MaxInt32}
	for !rd.atEnd() {
		tag := rd.readTag()
		switch tag {
		case 0:
		case 8:
			kindVal = rd.readInt32()
		case 18:
			rd.checkRecursionLimit()
			r.childBytes = append(r.childBytes, rd.readMessageBytes())
		case 24:
			stringValuePointer = rd.readUint32()
			doubleValue = 0
			templateRawPtr = 0
			templateCookedPtr = 0
		case 33:
			doubleValue = rd.readDouble()
			stringValuePointer = 0
			templateRawPtr = 0
			templateCookedPtr = 0
		case 40:
			relativeLine = rd.readSint32()
		case 48:
			relativeColumn = rd.readSint32()
		case 58:
			rd.checkRecursionLimit()
			if !hasJsdoc {
				r.jsdoc = OptimizationJsdoc{kind_: r.jsdocKindBuf[:0]}
				hasJsdoc = true
			}
			subRd := protoReader{b: rd.readMessageBytes(), depth: depth + 1, limit: math.MaxInt32}
			r.jsdoc.mergeFromReader(&subRd)
		case 66:
			rd.checkRecursionLimit()
			stringValuePointer = 0
			doubleValue = 0
			subRd := protoReader{b: rd.readMessageBytes(), depth: depth + 1, limit: math.MaxInt32}
			for !subRd.atEnd() {
				subTag := subRd.readTag()
				switch subTag {
				case 0:
				case 8:
					templateRawPtr = subRd.readUint32()
				case 16:
					templateCookedPtr = subRd.readInt32()
				default:
					subRd.skipField(subTag)
				}
			}
		case 72:
			originalNamePointer = rd.readUint32()
		case 112:
			booleanProperties = rd.readInt64()
		case 120:
			typeVal = rd.readInt32()
			hasType = true
		case 128:
			sourceFile = rd.readUint32()
		default:
			rd.skipField(tag)
		}
	}

	owner := r.owner_
	if sourceFileTemplate.IsNil() || sourceFile != 0 {
		// 0 == 'not set'
		sourceFileTemplate = owner.CreateSourceInfoTemplate(owner.filePool[sourceFile-1])
	}

	currentLine := r.previousLine + int(relativeLine)
	currentColumn := r.previousColumn + int(relativeColumn)

	n := r.deserializeSingleNode(
		nodeKind_fromNumber(kindVal),
		stringValuePointer,
		doubleValue,
		templateRawPtr,
		templateCookedPtr,
		sourceFileTemplate)
	if hasType && owner.colorPoolShard != nil {
		n.SetColor(owner.colorPoolShard.GetColor(int(typeVal)))
	}
	if booleanProperties > 0 {
		n.DeserializeProperties(
			owner.FilterOutCastProp(booleanProperties), sourceFileTemplate.GetIsInClosureUnawareSubtree())
	}
	if hasJsdoc {
		n.SetJSDocInfo(jsDocSerializer_deserializeJsdocWithBuilder(&r.jsdoc, owner.stringPool, &r.jsdocBuilder))
	}
	n.SetLinenoCharno(currentLine, currentColumn)
	r.previousLine = currentLine
	r.previousColumn = currentColumn
	r.contextStack = append(r.contextStack, r.featureCollector.VisitSingleNode(r.contextStack, n))

	if rhino.Node_HasBitSet(booleanProperties, int(rhino.NodeProperty_CLOSURE_UNAWARE_SHADOW.GetNumber())) {
		if childrenStart >= len(r.childBytes) {
			panic(jutil.NewIndexOutOfBoundsException("Index 0 out of bounds for length 0"))
		}
		serializedShadowChild := r.childBytes[childrenStart]
		r.childBytes = r.childBytes[:childrenStart]
		closureUnawareSubtreeTemplateNode := sourceFileTemplate.CloneTree()
		closureUnawareSubtreeTemplateNode.SetIsInClosureUnawareSubtree(true)
		shadowedCode, shadowOrigNamePtr := r.visitBytes(
			serializedShadowChild, depth+1, closureUnawareSubtreeTemplateNode)
		if shadowOrigNamePtr != 0 {
			shadowedCode.SetOriginalNameFromStringPool(
				owner.stringPool.GetInternedStrings(), int(shadowOrigNamePtr))
		}
		shadowRoot := rhino.IR_Root(rhino.IR_ScriptOfNodes(rhino.IR_ExprResult(shadowedCode)))
		shadowRoot.ClonePropsFrom(closureUnawareSubtreeTemplateNode)
		shadowRoot.GetFirstChild().ClonePropsFrom(closureUnawareSubtreeTemplateNode)
		shadowRoot.GetFirstFirstChild().ClonePropsFrom(closureUnawareSubtreeTemplateNode)

		n.SetClosureUnawareShadow(shadowRoot)
		return n, originalNamePointer
	}

	childrenEnd := len(r.childBytes)
	internedStrings := owner.stringPool.GetInternedStrings()
	for i := childrenStart; i < childrenEnd; i++ {
		deserializedChild, childOrigNamePtr := r.visitBytes(
			r.childBytes[i], depth+1, sourceFileTemplate)
		n.AddChildToBack(deserializedChild)
		if childOrigNamePtr != 0 {
			deserializedChild.SetOriginalNameFromStringPool(internedStrings, int(childOrigNamePtr))
		}
	}
	r.childBytes = r.childBytes[:childrenStart]

	r.contextStack = r.contextStack[:len(r.contextStack)-1]
	return n, originalNamePointer
}

func (r *scriptNodeDeserializer_Runner) stringNode(token rhino.Token, stringValPtr int32, templateNode rhino.Node) rhino.Node {
	id := r.owner_.stringPool.GetInternedStrings().GetID(int(stringValPtr))
	return rhino.AllocNode().InitFromTemplate(token, templateNode).SetStringID(id)
}

func (r *scriptNodeDeserializer_Runner) deserializeSingleNode(
	kind NodeKind,
	stringValPtr int32,
	doubleVal float64,
	templateRawPtr int32,
	templateCookedPtr int32,
	templateNode rhino.Node,
) rhino.Node {
	if uint(kind) < uint(len(nodeKindToDeserSpec)) {
		spec := nodeKindToDeserSpec[kind]
		switch spec.mode {
		case deserModePlain:
			return rhino.AllocNode().InitFromTemplate(spec.token, templateNode)
		case deserModeString:
			return r.stringNode(spec.token, stringValPtr, templateNode)
		case deserModeQuotedString:
			n := r.stringNode(spec.token, stringValPtr, templateNode)
			n.SetQuotedStringKey()
			return n
		case deserModePostIncrDecr:
			n := rhino.AllocNode().InitFromTemplate(spec.token, templateNode)
			n.PutBooleanProp(rhino.Node_INCRDECR_PROP, true)
			return n
		case deserModeNumber:
			return rhino.AllocNode().InitFromTemplate(rhino.Token_NUMBER, templateNode).SetNumberRaw(doubleVal)
		case deserModeBigInt:
			bigintString := r.owner_.stringPool.Get(int(stringValPtr))
			return rhino.AllocNode().InitBigIntFromTemplate(jutil.NewBigIntegerFromString(bigintString), templateNode)
		case deserModeTemplateLitString:
			return rhino.AllocNode().InitTemplateLitStringFromStringPool(
				r.owner_.stringPool.GetInternedStrings(),
				int(templateCookedPtr),
				int(templateRawPtr),
				templateNode)
		}
	}
	panic(jutil.NewIllegalStateExceptionMsg("Unexpected serialized kind for AstNode: " + kind.String()))
}

// Java: ScriptNodeDeserializer.getSourceMappingURL()
func (d *ScriptNodeDeserializer) GetSourceMappingURL() string {
	return d.sourceMappingURL
}

// Java: ScriptNodeDeserializer.getSourceFile()
func (d *ScriptNodeDeserializer) GetSourceFile() *SourceFile {
	return d.sourceFile
}

// CreateSourceInfoTemplate creates a template node to use as a source of
// common attributes.
//
// This allows the prop structure to be shared among all the node from this
// source file. This reduces the cost of these properties to O(nodes) to
// O(files).
//
// Java: ScriptNodeDeserializer.createSourceInfoTemplate(SourceFile)
func (d *ScriptNodeDeserializer) CreateSourceInfoTemplate(file *SourceFile) rhino.Node {
	return file.getOrCreateSourceInfoTemplate()
}

// FilterOutCastProp filters out any NodeProperty.COLOR_BEFORE_CASTs if no
// colors are being deserialized.
//
// This is because it doesn't make sense to have that property present on
// nodes that don't have colors.
//
// Java: ScriptNodeDeserializer.filterOutCastProp(long)
func (d *ScriptNodeDeserializer) FilterOutCastProp(nodeProperties int64) int64 {
	if d.colorPoolShard != nil {
		return nodeProperties // we are deserializing colors, so this is fine.
	}
	return nodeProperties &^ (int64(1) << rhino.NodeProperty_COLOR_FROM_CAST.GetNumber())
}
