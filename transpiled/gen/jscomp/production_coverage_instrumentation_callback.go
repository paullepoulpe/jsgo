/*
 * Copyright 2020 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/instrumentation/ProductionCoverageInstrumentationCallback.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	pkg_sourcemap "github.com/paullepoulpe/jsgo/transpiled/gen/sourcemap"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback
type ProductionCoverageInstrumentationCallback struct {
	instrumentationArrayName string                                                      `gson:"instrumentationArrayName"`
	compiler                 AbstractCompiler                                            `gson:"compiler"`
	parameterMapping         *ProductionCoverageInstrumentationCallback_ParameterMapping `gson:"parameterMapping"`
	functionNameStack        j2gort.Deque[string]                                        `gson:"functionNameStack"`
	idHash_                  atomic.Int32
}

func (this_ *ProductionCoverageInstrumentationCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *ProductionCoverageInstrumentationCallback) IsJavaNil_() bool { return this_ == nil }
func (*ProductionCoverageInstrumentationCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback", "ProductionCoverageInstrumentationCallback"
}
func (this_ *ProductionCoverageInstrumentationCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProductionCoverageInstrumentationCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback", "com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*ProductionCoverageInstrumentationCallback](v_)
	})
}
func (*ProductionCoverageInstrumentationCallback) Implements_jscomp_NodeTraversal_Callback() {}

func StaticGet_ProductionCoverageInstrumentationCallback_ANONYMOUS_FUNCTION_NAME() string {
	return "<Anonymous>"
}
func (this_ *ProductionCoverageInstrumentationCallback) initFields_ProductionCoverageInstrumentationCallback() {
	_ = this_
	this_.instrumentationArrayName = j2gort.NullString
	this_.functionNameStack = j2gort.Deque[string]{}
}

// Java constructor: ProductionCoverageInstrumentationCallback(com.google.javascript.jscomp.AbstractCompiler compiler, java.lang.String instrumentationArrayName)
func NewProductionCoverageInstrumentationCallback(compiler AbstractCompiler, instrumentationArrayName string) *ProductionCoverageInstrumentationCallback {
	this_ := &ProductionCoverageInstrumentationCallback{}
	InitProductionCoverageInstrumentationCallback(this_, compiler, instrumentationArrayName)
	return this_
}

func InitProductionCoverageInstrumentationCallback(this_ *ProductionCoverageInstrumentationCallback, compiler AbstractCompiler, instrumentationArrayName string) {
	_ = compiler
	_ = instrumentationArrayName
	this_.initFields_ProductionCoverageInstrumentationCallback()
	this_.compiler = compiler
	this_.parameterMapping = NewProductionCoverageInstrumentationCallback_ParameterMapping()
	this_.instrumentationArrayName = instrumentationArrayName
}

// Java method: ProductionCoverageInstrumentationCallback.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ProductionCoverageInstrumentationCallback) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if !(n.IsRoot()) && !(t.GetSourceName() == n.GetSourceFileName()) {
		return false
	}
	if n.IsFunction() {
		var fnName string = NodeUtil_GetBestLValueName(NodeUtil_GetBestLValue(n))
		_ = fnName
		if fnName == j2gort.NullString {
			fnName = "<Anonymous>"
		} else {
			fnName = fnName
		}
		this_.functionNameStack.Push(fnName)
	}
	if this_.functionNameStack.IsEmpty() {
		this_.functionNameStack.Push("<Anonymous>")
	}
	return true
}

// Java method: ProductionCoverageInstrumentationCallback.visit(com.google.javascript.jscomp.NodeTraversal traversal, com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node parent)
func (this_ *ProductionCoverageInstrumentationCallback) Visit_3(traversal *NodeTraversal, node pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = traversal
	_ = node
	_ = parent
	var fileName string = traversal.GetSourceName()
	_ = fileName
	var sourceFileName string = node.GetSourceFileName()
	_ = sourceFileName
	var functionName string = this_.functionNameStack.Peek()
	_ = functionName
	switch node.TokenSwitchOrd_() {
	case int32(73) /* Token.FUNCTION */ :
		this_.functionNameStack.Pop()
		this_.priv_ProductionCoverageInstrumentationCallback_instrumentBlockNode_6511(node.GetLastChild(), fileName, functionName, StaticGet_ProductionCoverageInstrumentationCallback_Type_FUNCTION())
	case int32(74) /* Token.IF */ :
		var ifTrueNode pkg_rhino.Node = node.GetSecondChild()
		_ = ifTrueNode
		this_.priv_ProductionCoverageInstrumentationCallback_instrumentBlockNode_6511(ifTrueNode, sourceFileName, functionName, StaticGet_ProductionCoverageInstrumentationCallback_Type_BRANCH())
		if node.GetChildCount() == int32(2) {
			this_.priv_ProductionCoverageInstrumentationCallback_addElseBlock_6513(node)
		}
		var ifFalseNode pkg_rhino.Node = node.GetLastChild()
		_ = ifFalseNode
		if NodeUtil_IsEmptyBlock(ifFalseNode) || (ifFalseNode.HasChildren() && !(ifFalseNode.GetFirstChild().IsIf())) {
			this_.priv_ProductionCoverageInstrumentationCallback_instrumentBlockNode_6511(ifFalseNode, sourceFileName, functionName, StaticGet_ProductionCoverageInstrumentationCallback_Type_BRANCH_DEFAULT())
		}
	case int32(75) /* Token.SWITCH */ :
		var hasDefaultCase bool = false
		_ = hasDefaultCase
		var switchBody pkg_rhino.Node = node.GetSecondChild()
		_ = switchBody
		{
			var c pkg_rhino.Node = switchBody.GetFirstChild()
			_ = c
			for ; c != 0; c = c.GetNext() {
				if c.IsDefaultCase() {
					this_.priv_ProductionCoverageInstrumentationCallback_instrumentBlockNode_6511(c.GetLastChild(), sourceFileName, functionName, StaticGet_ProductionCoverageInstrumentationCallback_Type_BRANCH_DEFAULT())
					hasDefaultCase = true
				} else {
					this_.priv_ProductionCoverageInstrumentationCallback_instrumentBlockNode_6511(c.GetLastChild(), sourceFileName, functionName, StaticGet_ProductionCoverageInstrumentationCallback_Type_BRANCH())
				}
			}
		}
		if !(hasDefaultCase) {
			var defaultBlock pkg_rhino.Node = pkg_rhino.IR_Block_0()
			_ = defaultBlock
			defaultBlock.SrcrefTreeIfMissing(node)
			var defaultCase pkg_rhino.Node = pkg_rhino.IR_DefaultCase(defaultBlock).SrcrefTreeIfMissing(node)
			_ = defaultCase
			switchBody.AddChildToBack(defaultCase)
			this_.priv_ProductionCoverageInstrumentationCallback_instrumentBlockNode_6511(defaultBlock, sourceFileName, functionName, StaticGet_ProductionCoverageInstrumentationCallback_Type_BRANCH_DEFAULT())
		}
	case int32(67) /* Token.HOOK */ :
		var ifTernaryIsTrueExpression pkg_rhino.Node = node.GetSecondChild()
		_ = ifTernaryIsTrueExpression
		var ifTernaryIsFalseExpression pkg_rhino.Node = node.GetLastChild()
		_ = ifTernaryIsFalseExpression
		this_.priv_ProductionCoverageInstrumentationCallback_addInstrumentationNodeWithComma_6510(ifTernaryIsTrueExpression, sourceFileName, functionName, StaticGet_ProductionCoverageInstrumentationCallback_Type_BRANCH())
		this_.priv_ProductionCoverageInstrumentationCallback_addInstrumentationNodeWithComma_6510(ifTernaryIsFalseExpression, sourceFileName, functionName, StaticGet_ProductionCoverageInstrumentationCallback_Type_BRANCH())
		this_.compiler.ReportChangeToEnclosingScope(node)
	case int32(68) /* Token.OR */, int32(69) /* Token.AND */, int32(70) /* Token.COALESCE */ :
		var secondExpression pkg_rhino.Node = node.GetLastChild()
		_ = secondExpression
		this_.priv_ProductionCoverageInstrumentationCallback_addInstrumentationNodeWithComma_6510(secondExpression, sourceFileName, functionName, StaticGet_ProductionCoverageInstrumentationCallback_Type_BRANCH())
		this_.compiler.ReportChangeToEnclosingScope(node)
	default:
		if NodeUtil_IsLoopStructure(node) {
			var blockNode pkg_rhino.Node = NodeUtil_GetLoopCodeBlock(node)
			_ = blockNode
			j2gort.CheckNotNullNode(blockNode)
			this_.priv_ProductionCoverageInstrumentationCallback_instrumentBlockNode_6511(blockNode, sourceFileName, functionName, StaticGet_ProductionCoverageInstrumentationCallback_Type_BRANCH())
		}
	}
}

// Java method: ProductionCoverageInstrumentationCallback.addInstrumentationNodeWithComma(com.google.javascript.rhino.Node originalNode, java.lang.String fileName, java.lang.String functionName, com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback.Type type)
func (this_ *ProductionCoverageInstrumentationCallback) priv_ProductionCoverageInstrumentationCallback_addInstrumentationNodeWithComma_6510(originalNode pkg_rhino.Node, fileName string, functionName string, type_ ProductionCoverageInstrumentationCallback_Type) {
	_ = originalNode
	_ = fileName
	_ = functionName
	_ = type_
	var cloneOfOriginal pkg_rhino.Node = originalNode.CloneTree_0()
	_ = cloneOfOriginal
	var newInstrumentationNode pkg_rhino.Node = this_.priv_ProductionCoverageInstrumentationCallback_newInstrumentationNode_6512(cloneOfOriginal, fileName, functionName, type_)
	_ = newInstrumentationNode
	var childOfInstrumentationNode pkg_rhino.Node = newInstrumentationNode.RemoveFirstChild()
	_ = childOfInstrumentationNode
	var infusedExp pkg_rhino.Node = AstManipulations_FuseExpressions(childOfInstrumentationNode, cloneOfOriginal)
	_ = infusedExp
	originalNode.ReplaceWith(infusedExp)
}

// Java method: ProductionCoverageInstrumentationCallback.instrumentBlockNode(com.google.javascript.rhino.Node block, java.lang.String fileName, java.lang.String fnName, com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback.Type type)
func (this_ *ProductionCoverageInstrumentationCallback) priv_ProductionCoverageInstrumentationCallback_instrumentBlockNode_6511(block pkg_rhino.Node, fileName string, fnName string, type_ ProductionCoverageInstrumentationCallback_Type) {
	_ = block
	_ = fileName
	_ = fnName
	_ = type_
	block.AddChildToFront(this_.priv_ProductionCoverageInstrumentationCallback_newInstrumentationNode_6512(block, fileName, fnName, type_))
	this_.compiler.ReportChangeToEnclosingScope(block)
}

// Java method: ProductionCoverageInstrumentationCallback.newInstrumentationNode(com.google.javascript.rhino.Node node, java.lang.String fileName, java.lang.String fnName, com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback.Type type)
func (this_ *ProductionCoverageInstrumentationCallback) priv_ProductionCoverageInstrumentationCallback_newInstrumentationNode_6512(node pkg_rhino.Node, fileName string, fnName string, type_ ProductionCoverageInstrumentationCallback_Type) pkg_rhino.Node {
	_ = node
	_ = fileName
	_ = fnName
	_ = type_
	var lineNo int32 = node.GetLineno()
	_ = lineNo
	var columnNo int32 = node.GetCharno()
	_ = columnNo
	if node.IsBlock() {
		lineNo = node.GetParent().GetLineno()
		columnNo = node.GetParent().GetCharno()
	}
	var encodedParam string = this_.parameterMapping.priv_ParameterMapping_getEncodedParam_6515(fileName, fnName, type_, lineNo, columnNo)
	_ = encodedParam
	var prop pkg_rhino.Node = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name(this_.instrumentationArrayName), "push")
	_ = prop
	var functionCall pkg_rhino.Node = pkg_rhino.IR_Call(prop, []pkg_rhino.Node{pkg_rhino.IR_String(encodedParam)})
	_ = functionCall
	var exprNode pkg_rhino.Node = pkg_rhino.IR_ExprResult(functionCall)
	_ = exprNode
	return exprNode.SrcrefTreeIfMissing(node)
}

// Java method: ProductionCoverageInstrumentationCallback.addElseBlock(com.google.javascript.rhino.Node node)
func (this_ *ProductionCoverageInstrumentationCallback) priv_ProductionCoverageInstrumentationCallback_addElseBlock_6513(node pkg_rhino.Node) pkg_rhino.Node {
	_ = node
	var defaultBlock pkg_rhino.Node = pkg_rhino.IR_Block_0()
	_ = defaultBlock
	node.AddChildToBack(defaultBlock)
	return defaultBlock.SrcrefTreeIfMissing(node)
}

// Java method: ProductionCoverageInstrumentationCallback.getInstrumentationMapping()
func (this_ *ProductionCoverageInstrumentationCallback) GetInstrumentationMapping() *VariableMap {
	return this_.parameterMapping.priv_ParameterMapping_getParamMappingAsVariableMap_6517()
}

// Java enum: com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback.Type
type ProductionCoverageInstrumentationCallback_Type uint8

const (
	enum_ProductionCoverageInstrumentationCallback_Type_FUNCTION       ProductionCoverageInstrumentationCallback_Type = 1
	enum_ProductionCoverageInstrumentationCallback_Type_BRANCH         ProductionCoverageInstrumentationCallback_Type = 2
	enum_ProductionCoverageInstrumentationCallback_Type_BRANCH_DEFAULT ProductionCoverageInstrumentationCallback_Type = 3
)

var enumBaseTable_ProductionCoverageInstrumentationCallback_Type = [...]j2gort.EnumBase{
	{Name_: "FUNCTION", Ordinal_: 0},
	{Name_: "BRANCH", Ordinal_: 1},
	{Name_: "BRANCH_DEFAULT", Ordinal_: 2},
}
var enumValues_ProductionCoverageInstrumentationCallback_Type = []ProductionCoverageInstrumentationCallback_Type{1, 2, 3}

func StaticGet_ProductionCoverageInstrumentationCallback_Type_FUNCTION() ProductionCoverageInstrumentationCallback_Type {
	return enum_ProductionCoverageInstrumentationCallback_Type_FUNCTION
}
func StaticGet_ProductionCoverageInstrumentationCallback_Type_BRANCH() ProductionCoverageInstrumentationCallback_Type {
	return enum_ProductionCoverageInstrumentationCallback_Type_BRANCH
}
func StaticGet_ProductionCoverageInstrumentationCallback_Type_BRANCH_DEFAULT() ProductionCoverageInstrumentationCallback_Type {
	return enum_ProductionCoverageInstrumentationCallback_Type_BRANCH_DEFAULT
}

func ProductionCoverageInstrumentationCallback_Type_Values() []ProductionCoverageInstrumentationCallback_Type {
	return j2gort.CloneSlice(enumValues_ProductionCoverageInstrumentationCallback_Type)
}

func ProductionCoverageInstrumentationCallback_Type_ValueOf(name string) ProductionCoverageInstrumentationCallback_Type {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ProductionCoverageInstrumentationCallback_Type {
		if enumBaseTable_ProductionCoverageInstrumentationCallback_Type[i_].Name_ == name {
			return ProductionCoverageInstrumentationCallback_Type(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback.Type." + name))
}

func (this_ ProductionCoverageInstrumentationCallback_Type) IsJavaNil_() bool                { return this_ == 0 }
func (ProductionCoverageInstrumentationCallback_Type) HasIdentityEquals_()                   {}
func (ProductionCoverageInstrumentationCallback_Type) Implements_j2gort_Ext_java_lang_Enum() {}
func (ProductionCoverageInstrumentationCallback_Type) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ ProductionCoverageInstrumentationCallback_Type) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ProductionCoverageInstrumentationCallback_Type) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback$Type", "Type"
}
func (this_ ProductionCoverageInstrumentationCallback_Type) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ProductionCoverageInstrumentationCallback_Type[this_-1]
}
func (this_ ProductionCoverageInstrumentationCallback_Type) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ProductionCoverageInstrumentationCallback_Type) EnumOrdinal_() int32 {
	return int32(this_) - 1
}
func (this_ ProductionCoverageInstrumentationCallback_Type) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ProductionCoverageInstrumentationCallback_Type[this_-1].Name_
}
func (this_ ProductionCoverageInstrumentationCallback_Type) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ProductionCoverageInstrumentationCallback_Type[this_-1].Name_
}
func (this_ ProductionCoverageInstrumentationCallback_Type) Equals(other_ any) bool {
	o_, ok_ := other_.(ProductionCoverageInstrumentationCallback_Type)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ProductionCoverageInstrumentationCallback_Type) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ProductionCoverageInstrumentationCallback_Type[this_-1].HashCode()
}
func (this_ ProductionCoverageInstrumentationCallback_Type) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ProductionCoverageInstrumentationCallback_Type[this_-1].CompareTo(other_)
}
func (this_ ProductionCoverageInstrumentationCallback_Type) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ProductionCoverageInstrumentationCallback_Type[this_-1].Name_)
}
func init() {
	vals_ := []any{ProductionCoverageInstrumentationCallback_Type(1), ProductionCoverageInstrumentationCallback_Type(2), ProductionCoverageInstrumentationCallback_Type(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback.Type", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback$Type", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback.Type", "com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback$Type", func(v_ any) bool {
		return j2gort.InstanceOf[ProductionCoverageInstrumentationCallback_Type](v_)
	})
}

// Java class: com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback.ParameterMapping
type ProductionCoverageInstrumentationCallback_ParameterMapping struct {
	paramValueEncodings  j2gort.Ext_java_util_Map         `gson:"paramValueEncodings"`
	fileNameToIndex      j2gort.OrderedMap[string, int32] `gson:"fileNameToIndex"`
	functionNameToIndex  j2gort.OrderedMap[string, int32] `gson:"functionNameToIndex"`
	typeToIndex          j2gort.OrderedMap[string, int32] `gson:"typeToIndex"`
	nextUniqueIdentifier int64                            `gson:"nextUniqueIdentifier"`
	idHash_              atomic.Int32
}

func (this_ *ProductionCoverageInstrumentationCallback_ParameterMapping) BindSelf_(s any) {
	_ = s
}
func (this_ *ProductionCoverageInstrumentationCallback_ParameterMapping) IsJavaNil_() bool {
	return this_ == nil
}
func (*ProductionCoverageInstrumentationCallback_ParameterMapping) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback$ParameterMapping", "ParameterMapping"
}
func (this_ *ProductionCoverageInstrumentationCallback_ParameterMapping) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProductionCoverageInstrumentationCallback_ParameterMapping) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback.ParameterMapping", "com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback$ParameterMapping", func(v_ any) bool {
		return j2gort.InstanceOf[*ProductionCoverageInstrumentationCallback_ParameterMapping](v_)
	})
}

func (this_ *ProductionCoverageInstrumentationCallback_ParameterMapping) initFields_ProductionCoverageInstrumentationCallback_ParameterMapping() {
	_ = this_
}

// Java constructor: ParameterMapping()
func NewProductionCoverageInstrumentationCallback_ParameterMapping() *ProductionCoverageInstrumentationCallback_ParameterMapping {
	this_ := &ProductionCoverageInstrumentationCallback_ParameterMapping{}
	InitProductionCoverageInstrumentationCallback_ParameterMapping(this_)
	return this_
}

func InitProductionCoverageInstrumentationCallback_ParameterMapping(this_ *ProductionCoverageInstrumentationCallback_ParameterMapping) {
	this_.initFields_ProductionCoverageInstrumentationCallback_ParameterMapping()
	this_.nextUniqueIdentifier = int64(int32(0))
	this_.paramValueEncodings = j2gort.NewExt_java_util_LinkedHashMap()
	this_.fileNameToIndex = j2gort.NewOrderedMap[string, int32]()
	this_.functionNameToIndex = j2gort.NewOrderedMap[string, int32]()
	this_.typeToIndex = j2gort.NewOrderedMap[string, int32]()
}

// Java method: ParameterMapping.getEncodedParam(java.lang.String fileName, java.lang.String functionName, com.google.javascript.jscomp.instrumentation.ProductionCoverageInstrumentationCallback.Type type, int lineNo, int colNo)
func (this_ *ProductionCoverageInstrumentationCallback_ParameterMapping) priv_ParameterMapping_getEncodedParam_6515(fileName string, functionName string, type_ ProductionCoverageInstrumentationCallback_Type, lineNo int32, colNo int32) string {
	_ = fileName
	_ = functionName
	_ = type_
	_ = lineNo
	_ = colNo
	this_.fileNameToIndex.PutIfAbsent(fileName, this_.fileNameToIndex.Size())
	this_.functionNameToIndex.PutIfAbsent(functionName, this_.functionNameToIndex.Size())
	this_.typeToIndex.PutIfAbsent((type_).Name(), this_.typeToIndex.Size())
	var sb *j2gort.Ext_java_lang_StringBuilder = j2gort.NewExt_java_lang_StringBuilder()
	_ = sb
	func() {
		defer func() {
			if rec_1 := recover(); rec_1 != nil {
				jex_2 := javalang.RecoverException(rec_1)
				_ = jex_2
				if j2gort.InstanceOf[javaio.IIOException](jex_2) {
					e := j2gort.AsIface[javaio.IIOException](jex_2)
					_ = e
					panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", e)))
				} else {
					panic(rec_1)
				}
			}
		}()
		pkg_sourcemap.Base64VLQ_Encode(j2gort.AsIface[j2gort.Ext_java_lang_Appendable](sb), func() int32 {
			if v_, ok_ := this_.fileNameToIndex.GetOk(fileName); ok_ {
				return v_
			}
			j2gort.ThrowNPE()
			return 0
		}())
		pkg_sourcemap.Base64VLQ_Encode(j2gort.AsIface[j2gort.Ext_java_lang_Appendable](sb), func() int32 {
			if v_, ok_ := this_.functionNameToIndex.GetOk(functionName); ok_ {
				return v_
			}
			j2gort.ThrowNPE()
			return 0
		}())
		pkg_sourcemap.Base64VLQ_Encode(j2gort.AsIface[j2gort.Ext_java_lang_Appendable](sb), func() int32 {
			if v_, ok_ := this_.typeToIndex.GetOk((type_).Name()); ok_ {
				return v_
			}
			j2gort.ThrowNPE()
			return 0
		}())
		pkg_sourcemap.Base64VLQ_Encode(j2gort.AsIface[j2gort.Ext_java_lang_Appendable](sb), lineNo)
		pkg_sourcemap.Base64VLQ_Encode(j2gort.AsIface[j2gort.Ext_java_lang_Appendable](sb), colNo)
	}()
	var encodedParam string = sb.String()
	_ = encodedParam
	if !(j2gort.MapContainsKeyStr(this_.paramValueEncodings, encodedParam)) {
		var uniqueIdentifier int64 = this_.priv_ParameterMapping_generateUniqueIdentifier_6516()
		_ = uniqueIdentifier
		if uniqueIdentifier > int64(int32(2147483647)) {
			panic(j2gort.AsIface[javalang.IArithmeticException](j2gort.ConstructShimThrowable("java.lang.ArithmeticException", "Unique Identifier exceeds value of Integer.MAX_VALUE, could not encode with Base 64 VLQ")))
		}
		sb = j2gort.NewExt_java_lang_StringBuilder()
		func() {
			defer func() {
				if rec_3 := recover(); rec_3 != nil {
					jex_4 := javalang.RecoverException(rec_3)
					_ = jex_4
					if j2gort.InstanceOf[javaio.IIOException](jex_4) {
						e_2 := j2gort.AsIface[javaio.IIOException](jex_4)
						_ = e_2
						panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", e_2)))
					} else {
						panic(rec_3)
					}
				}
			}()
			pkg_sourcemap.Base64VLQ_Encode(j2gort.AsIface[j2gort.Ext_java_lang_Appendable](sb), j2gort.Ext_java_lang_Math_ToIntExact(uniqueIdentifier))
		}()
		j2gort.MapPutStr(this_.paramValueEncodings, encodedParam, sb.String())
	}
	return j2gort.CastString(j2gort.MapGetStr(this_.paramValueEncodings, encodedParam))
}

// Java method: ParameterMapping.generateUniqueIdentifier()
func (this_ *ProductionCoverageInstrumentationCallback_ParameterMapping) priv_ParameterMapping_generateUniqueIdentifier_6516() int64 {
	this_.nextUniqueIdentifier = this_.nextUniqueIdentifier + 1
	return this_.nextUniqueIdentifier
}

// Java method: ParameterMapping.getParamMappingAsVariableMap()
func (this_ *ProductionCoverageInstrumentationCallback_ParameterMapping) priv_ParameterMapping_getParamMappingAsVariableMap_6517() *VariableMap {
	var gson *j2gort.Ext_com_google_gson_Gson = j2gort.NewExt_com_google_gson_GsonBuilder().DisableHtmlEscaping().Create()
	_ = gson
	j2gort.MapPutStr(this_.paramValueEncodings, gson.ToJson_1_any(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.fileNameToIndex.LiveKeySet())), " FileNames")
	j2gort.MapPutStr(this_.paramValueEncodings, gson.ToJson_1_any(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.functionNameToIndex.LiveKeySet())), " FunctionNames")
	j2gort.MapPutStr(this_.paramValueEncodings, gson.ToJson_1_any(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.typeToIndex.LiveKeySet())), " Types")
	var val_5 VariableMap
	InitVariableMap(&val_5, this_.paramValueEncodings)
	var preInversedMap *VariableMap = &val_5
	_ = preInversedMap
	var inversedMap j2gort.Ext_com_google_common_collect_ImmutableMap = preInversedMap.GetNewNameToOriginalNameMap()
	_ = inversedMap
	return NewVariableMap(j2gort.AsIface[j2gort.Ext_java_util_Map](inversedMap))
}
func (*ProductionCoverageInstrumentationCallback) CanReleaseSyntacticScopes_() {}
