/*
 * Copyright 2014 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/TemplateAstMatcher.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.TemplateAstMatcher
type TemplateAstMatcher struct {
	topScope             *TypedScope                         `gson:"topScope"`
	typeRegistry         *pkg_rhino.JSTypeRegistry           `gson:"typeRegistry"`
	templateParams       j2gort.Ext_java_util_List           `gson:"templateParams"`
	paramNodeMatches     j2gort.Ext_java_util_ArrayList      `gson:"paramNodeMatches"`
	templateLocals       j2gort.Ext_java_util_List           `gson:"templateLocals"`
	localVarMatches      j2gort.Ext_java_util_ArrayList      `gson:"localVarMatches"`
	stringLiteralMatches *j2gort.Ext_java_util_LinkedHashMap `gson:"stringLiteralMatches"`
	templateStart        pkg_rhino.Node                      `gson:"templateStart"`
	idHash_              atomic.Int32
	isLooseMatch         bool                 `gson:"isLooseMatch"`
	typeMatchingStrategy TypeMatchingStrategy `gson:"typeMatchingStrategy"`
}

func (this_ *TemplateAstMatcher) BindSelf_(s any) {
	_ = s
}
func (this_ *TemplateAstMatcher) IsJavaNil_() bool { return this_ == nil }
func (*TemplateAstMatcher) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TemplateAstMatcher", "TemplateAstMatcher"
}
func (this_ *TemplateAstMatcher) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TemplateAstMatcher) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TemplateAstMatcher", "com.google.javascript.jscomp.TemplateAstMatcher", func(v_ any) bool {
		return j2gort.InstanceOf[*TemplateAstMatcher](v_)
	})
}

var static_TemplateAstMatcher_TEMPLATE_TYPE_PARAM pkg_rhino.Token

func StaticGet_TemplateAstMatcher_TEMPLATE_TYPE_PARAM() pkg_rhino.Token {
	clinit_TemplateAstMatcher()
	return static_TemplateAstMatcher_TEMPLATE_TYPE_PARAM
}
func StaticSet_TemplateAstMatcher_TEMPLATE_TYPE_PARAM(v pkg_rhino.Token) {
	clinit_TemplateAstMatcher()
	static_TemplateAstMatcher_TEMPLATE_TYPE_PARAM = v
}

var static_TemplateAstMatcher_TEMPLATE_LOCAL_NAME pkg_rhino.Token

func StaticGet_TemplateAstMatcher_TEMPLATE_LOCAL_NAME() pkg_rhino.Token {
	clinit_TemplateAstMatcher()
	return static_TemplateAstMatcher_TEMPLATE_LOCAL_NAME
}
func StaticSet_TemplateAstMatcher_TEMPLATE_LOCAL_NAME(v pkg_rhino.Token) {
	clinit_TemplateAstMatcher()
	static_TemplateAstMatcher_TEMPLATE_LOCAL_NAME = v
}

var static_TemplateAstMatcher_TEMPLATE_STRING_LITERAL pkg_rhino.Token

func StaticGet_TemplateAstMatcher_TEMPLATE_STRING_LITERAL() pkg_rhino.Token {
	clinit_TemplateAstMatcher()
	return static_TemplateAstMatcher_TEMPLATE_STRING_LITERAL
}
func StaticSet_TemplateAstMatcher_TEMPLATE_STRING_LITERAL(v pkg_rhino.Token) {
	clinit_TemplateAstMatcher()
	static_TemplateAstMatcher_TEMPLATE_STRING_LITERAL = v
}

var clinitDone_TemplateAstMatcher uint32
var clinitState_TemplateAstMatcher j2gort.ClinitState

func clinit_TemplateAstMatcher() {
	if atomic.LoadUint32(&clinitDone_TemplateAstMatcher) != 0 {
		return
	}
	clinitSlow_TemplateAstMatcher()
}

func clinitSlow_TemplateAstMatcher() {
	if !j2gort.EnterClinit(&clinitState_TemplateAstMatcher, &clinitDone_TemplateAstMatcher) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_TemplateAstMatcher, &clinitDone_TemplateAstMatcher)
	static_TemplateAstMatcher_TEMPLATE_TYPE_PARAM = pkg_rhino.StaticGet_Token_PLACEHOLDER1()
	static_TemplateAstMatcher_TEMPLATE_LOCAL_NAME = pkg_rhino.StaticGet_Token_PLACEHOLDER2()
	static_TemplateAstMatcher_TEMPLATE_STRING_LITERAL = pkg_rhino.StaticGet_Token_PLACEHOLDER3()
}

func (this_ *TemplateAstMatcher) initFields_TemplateAstMatcher() {
	_ = this_
	this_.templateParams = j2gort.NewExt_java_util_ArrayList()
	this_.paramNodeMatches = j2gort.NewExt_java_util_ArrayList()
	this_.templateLocals = j2gort.NewExt_java_util_ArrayList()
	this_.localVarMatches = j2gort.NewExt_java_util_ArrayList()
	this_.stringLiteralMatches = j2gort.NewExt_java_util_LinkedHashMap()

}

// Java constructor: TemplateAstMatcher(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node templateFunctionNode, com.google.javascript.jscomp.TypeMatchingStrategy typeMatchingStrategy)
func NewTemplateAstMatcher(compiler AbstractCompiler, templateFunctionNode pkg_rhino.Node, typeMatchingStrategy TypeMatchingStrategy) *TemplateAstMatcher {
	this_ := &TemplateAstMatcher{}
	InitTemplateAstMatcher(this_, compiler, templateFunctionNode, typeMatchingStrategy)
	return this_
}

func InitTemplateAstMatcher(this_ *TemplateAstMatcher, compiler AbstractCompiler, templateFunctionNode pkg_rhino.Node, typeMatchingStrategy TypeMatchingStrategy) {
	_ = compiler
	_ = templateFunctionNode
	_ = typeMatchingStrategy
	this_.initFields_TemplateAstMatcher()
	if !(templateFunctionNode.IsFunction()) {
		j2gort.CheckStateFmt(false, "Template node must be a function node. Received: %s", j2gort.BoxNode(templateFunctionNode))
	}
	this_.topScope = compiler.GetTopScope()
	this_.typeRegistry = j2gort.CheckNotNullPtr(compiler.GetTypeRegistry())
	this_.templateStart = this_.priv_TemplateAstMatcher_initTemplate_5317(templateFunctionNode)
	this_.typeMatchingStrategy = j2gort.CheckNotNullEnum(typeMatchingStrategy)
}

// Java method: TemplateAstMatcher.matches(com.google.javascript.rhino.Node n)
func (this_ *TemplateAstMatcher) Matches(n pkg_rhino.Node) bool {
	_ = n
	if this_.priv_TemplateAstMatcher_matchesTemplateShape_5327(this_.templateStart, n) {
		if this_.paramNodeMatches.IsEmpty() && this_.localVarMatches.IsEmpty() {
			return true
		}
		this_.priv_TemplateAstMatcher_reset_5321()
		return this_.priv_TemplateAstMatcher_matchesTemplate_5329(this_.templateStart, n)
	}
	return false
}

// Java method: TemplateAstMatcher.isLooseMatch()
func (this_ *TemplateAstMatcher) IsLooseMatch() bool {
	return this_.isLooseMatch
}

// Java method: TemplateAstMatcher.getTemplateNodeToMatchMap()
func (this_ *TemplateAstMatcher) GetTemplateNodeToMatchMap() j2gort.Ext_java_util_Map {
	var map_ j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap_1_Ext_java_util_Map(j2gort.AsIface[j2gort.Ext_java_util_Map](this_.stringLiteralMatches))
	_ = map_
	{
		var i int32 = int32(0)
		_ = i
		for ; i < this_.templateParams.Size(); i = i + 1 {
			var name string = j2gort.ListGetStr(this_.templateParams, i)
			_ = name
			j2gort.MapPutStr(map_, name, j2gort.BoxNode(j2gort.CastNode[pkg_rhino.Node](this_.paramNodeMatches.Get(i))))
		}
	}
	{
		var i_2 int32 = int32(0)
		_ = i_2
		for ; i_2 < this_.templateLocals.Size(); i_2 = i_2 + 1 {
			var name_2 string = j2gort.ListGetStr(this_.templateLocals, i_2)
			_ = name_2
			j2gort.MapPutStr(map_, name_2, j2gort.BoxNode(pkg_rhino.IR_Name(j2gort.ListGetStr(this_.localVarMatches, i_2))))
		}
	}
	return map_
}

// Java method: TemplateAstMatcher.initTemplate(com.google.javascript.rhino.Node templateFunctionNode)
func (this_ *TemplateAstMatcher) priv_TemplateAstMatcher_initTemplate_5317(templateFunctionNode pkg_rhino.Node) pkg_rhino.Node {
	_ = templateFunctionNode
	var prepped pkg_rhino.Node = templateFunctionNode.CloneTree_0()
	_ = prepped
	this_.priv_TemplateAstMatcher_prepTemplatePlaceholders_5318(prepped)
	var body pkg_rhino.Node = prepped.GetLastChild()
	_ = body
	var startNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = startNode
	if body.HasOneChild() && body.GetFirstChild().IsExprResult() {
		startNode = body.GetFirstFirstChild()
	} else {
		startNode = body.GetFirstChild()
	}
	{
		var i int32 = int32(0)
		_ = i
		for ; i < this_.templateLocals.Size(); i = i + 1 {
			this_.localVarMatches.Add_1(nil)
		}
	}
	{
		var i_2 int32 = int32(0)
		_ = i_2
		for ; i_2 < this_.templateParams.Size(); i_2 = i_2 + 1 {
			this_.paramNodeMatches.Add_1(nil)
		}
	}
	return startNode
}

// Java method: TemplateAstMatcher.prepTemplatePlaceholders(com.google.javascript.rhino.Node fn)
func (this_ *TemplateAstMatcher) priv_TemplateAstMatcher_prepTemplatePlaceholders_5318(fn pkg_rhino.Node) {
	_ = fn
	var locals j2gort.Ext_java_util_List = this_.templateLocals
	_ = locals
	var params j2gort.Ext_java_util_List = this_.templateParams
	_ = params
	var paramTypes j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = paramTypes
	var fnName string = fn.GetFirstChild().GetString()
	_ = fnName
	fn.GetFirstChild().SetString("")
	var templateParametersNode pkg_rhino.Node = fn.GetSecondChild()
	_ = templateParametersNode
	var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(fn)
	_ = info
	if templateParametersNode.HasChildren() {
		j2gort.CheckNotNullPtr(info)
	}
	{
		var paramNode pkg_rhino.Node = templateParametersNode.GetFirstChild()
		_ = paramNode
		for ; paramNode != 0; paramNode = paramNode.GetNext() {
			var name string = paramNode.GetString()
			_ = name
			var expression pkg_rhino.JSTypeExpression = info.GetParameterType(name)
			_ = expression
			j2gort.CheckNotNullVal(expression)
			var type_ pkg_rhino.JSType = this_.typeRegistry.EvaluateTypeExpression(expression, pkg_rhino.ToIface_StaticTypedScope(this_.topScope))
			_ = type_
			j2gort.CheckNotNullIface(type_)
			params.Add_1(j2gort.BoxString(name))
			j2gort.MapPutStr(paramTypes, name, type_)
		}
	}
	this_.priv_TemplateAstMatcher_traverse_5320(fn, NewTemplateAstMatcher__Anon_97(this_, locals, params, paramTypes))
}

// Java method: TemplateAstMatcher.replaceNodeInPlace(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node replacement)
func (this_ *TemplateAstMatcher) ReplaceNodeInPlace(n pkg_rhino.Node, replacement pkg_rhino.Node) {
	_ = n
	_ = replacement
	if n.HasChildren() {
		var children pkg_rhino.Node = n.RemoveChildren()
		_ = children
		replacement.AddChildrenToFront(children)
	}
	n.ReplaceWith(replacement)
}

// Java method: TemplateAstMatcher.traverse(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.TemplateAstMatcher.Visitor callback)
func (this_ *TemplateAstMatcher) priv_TemplateAstMatcher_traverse_5320(n pkg_rhino.Node, callback TemplateAstMatcher_Visitor) {
	_ = n
	_ = callback
	var next pkg_rhino.Node = pkg_rhino.Node(0)
	_ = next
	{
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = next {
			next = c.GetNext()
			this_.priv_TemplateAstMatcher_traverse_5320(c, callback)
		}
	}
	callback.Visit(n)
}

// Java method: TemplateAstMatcher.reset()
func (this_ *TemplateAstMatcher) priv_TemplateAstMatcher_reset_5321() {
	this_.isLooseMatch = false
	j2gort.Ext_java_util_Collections_Fill(j2gort.AsIface[j2gort.Ext_java_util_List](this_.localVarMatches), nil)
	{
		var i int32 = int32(0)
		_ = i
		for ; i < this_.paramNodeMatches.Size(); i = i + 1 {
			this_.paramNodeMatches.Set(i, nil)
		}
	}
}

// Java method: TemplateAstMatcher.isTemplateParameterNode(com.google.javascript.rhino.Node n)
func (this_ *TemplateAstMatcher) priv_TemplateAstMatcher_isTemplateParameterNode_5322(n pkg_rhino.Node) bool {
	_ = n
	return (n.GetToken() == StaticGet_TemplateAstMatcher_TEMPLATE_TYPE_PARAM())
}

// Java method: TemplateAstMatcher.isTemplateParameterStringLiteralNode(com.google.javascript.rhino.Node n)
func (this_ *TemplateAstMatcher) priv_TemplateAstMatcher_isTemplateParameterStringLiteralNode_5323(n pkg_rhino.Node) bool {
	_ = n
	return (n.GetToken() == StaticGet_TemplateAstMatcher_TEMPLATE_STRING_LITERAL())
}

// Java method: TemplateAstMatcher.createTemplateParameterNode(int index, com.google.javascript.rhino.jstype.JSType type, boolean isStringLiteral)
func (this_ *TemplateAstMatcher) priv_TemplateAstMatcher_createTemplateParameterNode_5324(index int32, type_ pkg_rhino.JSType, isStringLiteral bool) pkg_rhino.Node {
	_ = index
	_ = type_
	_ = isStringLiteral
	if !(index >= int32(0)) {
		j2gort.CheckState(false)
	}
	j2gort.CheckNotNullIface(type_)
	var n pkg_rhino.Node = pkg_rhino.Node_NewNumber(float64(index))
	_ = n
	if isStringLiteral {
		n.SetToken(StaticGet_TemplateAstMatcher_TEMPLATE_STRING_LITERAL())
	} else {
		n.SetToken(StaticGet_TemplateAstMatcher_TEMPLATE_TYPE_PARAM())
	}
	n.SetJSType(type_)
	return n
}

// Java method: TemplateAstMatcher.isTemplateLocalNameNode(com.google.javascript.rhino.Node n)
func (this_ *TemplateAstMatcher) priv_TemplateAstMatcher_isTemplateLocalNameNode_5325(n pkg_rhino.Node) bool {
	_ = n
	return (n.GetToken() == StaticGet_TemplateAstMatcher_TEMPLATE_LOCAL_NAME())
}

// Java method: TemplateAstMatcher.createTemplateLocalNameNode(int index)
func (this_ *TemplateAstMatcher) priv_TemplateAstMatcher_createTemplateLocalNameNode_5326(index int32) pkg_rhino.Node {
	_ = index
	if !(index >= int32(0)) {
		j2gort.CheckState(false)
	}
	var n pkg_rhino.Node = pkg_rhino.Node_NewNumber(float64(index))
	_ = n
	n.SetToken(StaticGet_TemplateAstMatcher_TEMPLATE_LOCAL_NAME())
	return n
}

// Java method: TemplateAstMatcher.matchesTemplateShape(com.google.javascript.rhino.Node template, com.google.javascript.rhino.Node ast)
func (this_ *TemplateAstMatcher) priv_TemplateAstMatcher_matchesTemplateShape_5327(template pkg_rhino.Node, ast pkg_rhino.Node) bool {
	_ = template
	_ = ast
	for template != 0 {
		if (ast == 0) || !(this_.priv_TemplateAstMatcher_matchesNodeShape_5328(template, ast)) {
			return false
		}
		template = template.GetNext()
		ast = ast.GetNext()
	}
	return true
}

// Java method: TemplateAstMatcher.matchesNodeShape(com.google.javascript.rhino.Node template, com.google.javascript.rhino.Node ast)
func (this_ *TemplateAstMatcher) priv_TemplateAstMatcher_matchesNodeShape_5328(template pkg_rhino.Node, ast pkg_rhino.Node) bool {
	_ = template
	_ = ast
	if this_.priv_TemplateAstMatcher_isTemplateParameterNode_5322(template) {
		return !(NodeUtil_IsStatement(ast))
	} else if this_.priv_TemplateAstMatcher_isTemplateLocalNameNode_5325(template) {
		if !(ast.IsName()) {
			return false
		}
	} else if this_.priv_TemplateAstMatcher_isTemplateParameterStringLiteralNode_5323(template) {
		return NodeUtil_IsSomeCompileTimeConstStringValue(ast)
	} else if template.IsCall() {
		if ((ast == 0) || !(ast.IsCall())) || !(ast.HasXChildren(template.GetChildCount())) {
			return false
		}
	} else if !(template.IsEquivalentToShallow(ast)) {
		return false
	}
	var templateChild pkg_rhino.Node = template.GetFirstChild()
	_ = templateChild
	var astChild pkg_rhino.Node = ast.GetFirstChild()
	_ = astChild
	for templateChild != 0 {
		if !(this_.priv_TemplateAstMatcher_matchesNodeShape_5328(templateChild, astChild)) {
			return false
		}
		templateChild = templateChild.GetNext()
		astChild = astChild.GetNext()
	}
	return true
}

// Java method: TemplateAstMatcher.matchesTemplate(com.google.javascript.rhino.Node template, com.google.javascript.rhino.Node ast)
func (this_ *TemplateAstMatcher) priv_TemplateAstMatcher_matchesTemplate_5329(template pkg_rhino.Node, ast pkg_rhino.Node) bool {
	_ = template
	_ = ast
	for template != 0 {
		if (ast == 0) || !(this_.priv_TemplateAstMatcher_matchesNode_5330(template, ast)) {
			return false
		}
		template = template.GetNext()
		ast = ast.GetNext()
	}
	return true
}

// Java method: TemplateAstMatcher.matchesNode(com.google.javascript.rhino.Node template, com.google.javascript.rhino.Node ast)
func (this_ *TemplateAstMatcher) priv_TemplateAstMatcher_matchesNode_5330(template pkg_rhino.Node, ast pkg_rhino.Node) bool {
	_ = template
	_ = ast
	if this_.priv_TemplateAstMatcher_isTemplateParameterNode_5322(template) {
		var paramIndex int32 = javalang.D2I((template.GetDouble()))
		_ = paramIndex
		var previousMatch pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](this_.paramNodeMatches.Get(paramIndex))
		_ = previousMatch
		if previousMatch != 0 {
			return ast.IsEquivalentTo_1(previousMatch)
		}
		var templateType pkg_rhino.JSType = template.GetJSType()
		_ = templateType
		j2gort.CheckNotNullIface(templateType)
		if this_.priv_TemplateAstMatcher_isUnresolvedType_5331(templateType) {
			return false
		}
		var matchResult TypeMatchingStrategy_MatchResult = this_.typeMatchingStrategy.Match(templateType, ast.GetJSType())
		_ = matchResult
		this_.isLooseMatch = matchResult.IsLooseMatch()
		var isMatch bool = matchResult.IsMatch()
		_ = isMatch
		if isMatch && (previousMatch == 0) {
			this_.paramNodeMatches.Set(paramIndex, j2gort.BoxNode(ast))
		}
		return isMatch
	} else if this_.priv_TemplateAstMatcher_isTemplateLocalNameNode_5325(template) {
		var paramIndex_2 int32 = javalang.D2I((template.GetDouble()))
		_ = paramIndex_2
		var previouslyMatched bool = (j2gort.ListGetStr(this_.localVarMatches, paramIndex_2) != j2gort.NullString)
		_ = previouslyMatched
		if previouslyMatched {
			return func() bool {
				sr_1 := ast.GetString()
				se_2 := j2gort.ListGetStr(this_.localVarMatches, paramIndex_2)
				return j2gort.CheckNotNullString(sr_1) == se_2
			}()
		} else {
			var originalName string = ast.GetOriginalName()
			_ = originalName
			var name string
			_ = name
			if originalName != j2gort.NullString {
				name = originalName
			} else {
				name = ast.GetString()
			}
			this_.localVarMatches.Set(paramIndex_2, j2gort.BoxString(name))
		}
	} else if this_.priv_TemplateAstMatcher_isTemplateParameterStringLiteralNode_5323(template) {
		var paramIndex_3 int32 = javalang.D2I((template.GetDouble()))
		_ = paramIndex_3
		var previousMatch_2 pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](this_.paramNodeMatches.Get(paramIndex_3))
		_ = previousMatch_2
		if previousMatch_2 != 0 {
			return ast.IsEquivalentTo_1(previousMatch_2)
		}
		if NodeUtil_IsSomeCompileTimeConstStringValue(ast) {
			this_.paramNodeMatches.Set(paramIndex_3, j2gort.BoxNode(ast))
			return true
		}
		return false
	}
	var templateChild pkg_rhino.Node = template.GetFirstChild()
	_ = templateChild
	var astChild pkg_rhino.Node = ast.GetFirstChild()
	_ = astChild
	for templateChild != 0 {
		if !(this_.priv_TemplateAstMatcher_matchesNode_5330(templateChild, astChild)) {
			return false
		}
		templateChild = templateChild.GetNext()
		astChild = astChild.GetNext()
	}
	return true
}

// Java method: TemplateAstMatcher.isUnresolvedType(com.google.javascript.rhino.jstype.JSType type)
func (this_ *TemplateAstMatcher) priv_TemplateAstMatcher_isUnresolvedType_5331(type_ pkg_rhino.JSType) bool {
	_ = type_
	base_type__ := type_.GetJSTypeBase()
	if base_type__.IsNoResolvedType() || (base_type__.IsNamedType() && base_type__.IsUnknownType()) {
		return true
	}
	if base_type__.IsUnionType() {
		{
			cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(base_type__.GetUnionMembers()))
			for cur_3.HasNext() {
				rawElem_4 := cur_3.Next()
				alternate := j2gort.AsIface[pkg_rhino.JSType](rawElem_4)
				_ = alternate
				if this_.priv_TemplateAstMatcher_isUnresolvedType_5331(alternate) {
					return true
				}
			}
		}
	}
	return false
}

// Java interface: com.google.javascript.jscomp.TemplateAstMatcher.Visitor
type TemplateAstMatcher_Visitor interface {
	Implements_jscomp_TemplateAstMatcher_Visitor()
	IsJavaNil_() bool
	Visit(pkg_rhino.Node)
}

func ToIface_TemplateAstMatcher_Visitor[T interface {
	TemplateAstMatcher_Visitor
	comparable
}](v T) TemplateAstMatcher_Visitor {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_TemplateAstMatcher_Visitor[T TemplateAstMatcher_Visitor](v T) TemplateAstMatcher_Visitor {
	var vi TemplateAstMatcher_Visitor = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.TemplateAstMatcher.Visitor", "com.google.javascript.jscomp.TemplateAstMatcher$Visitor", func(v_ any) bool {
		return j2gort.InstanceOf[TemplateAstMatcher_Visitor](v_)
	})
}

// Java class:
type TemplateAstMatcher__Anon_97 struct {
	outer_           *TemplateAstMatcher
	cap_0_locals     j2gort.Ext_java_util_List
	cap_1_params     j2gort.Ext_java_util_List
	cap_2_paramTypes j2gort.Ext_java_util_Map
	idHash_          atomic.Int32
}

func (this_ *TemplateAstMatcher__Anon_97) BindSelf_(s any) {
	_ = s
}
func (this_ *TemplateAstMatcher__Anon_97) IsJavaNil_() bool { return this_ == nil }
func (*TemplateAstMatcher__Anon_97) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.TemplateAstMatcher$1", ""
}
func (this_ *TemplateAstMatcher__Anon_97) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*TemplateAstMatcher__Anon_97) HasIdentityEquals_()                           {}
func (this_ *TemplateAstMatcher__Anon_97) SetOuter_(o *TemplateAstMatcher)         { this_.outer_ = o }
func (this_ *TemplateAstMatcher__Anon_97) GetOuter_() *TemplateAstMatcher          { return this_.outer_ }
func (*TemplateAstMatcher__Anon_97) Implements_jscomp_TemplateAstMatcher_Visitor() {}

func (this_ *TemplateAstMatcher__Anon_97) initFields_TemplateAstMatcher__Anon_97() {
	_ = this_
}

func NewTemplateAstMatcher__Anon_97(outer_ *TemplateAstMatcher, cap_p_locals j2gort.Ext_java_util_List, cap_p_params j2gort.Ext_java_util_List, cap_p_paramTypes j2gort.Ext_java_util_Map) *TemplateAstMatcher__Anon_97 {
	this_ := &TemplateAstMatcher__Anon_97{}
	this_.outer_ = outer_
	this_.cap_0_locals = cap_p_locals
	this_.cap_1_params = cap_p_params
	this_.cap_2_paramTypes = cap_p_paramTypes
	this_.initFields_TemplateAstMatcher__Anon_97()
	return this_
}

// Java method: .visit(com.google.javascript.rhino.Node n)
func (this_ *TemplateAstMatcher__Anon_97) Visit(n pkg_rhino.Node) {
	_ = n
	if n.IsName() {
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		var name string = n.GetString()
		_ = name
		if (!(j2gort.CheckNotNullString(name) == "") && parent.IsVar()) && !(this_.cap_0_locals.Contains(j2gort.BoxString(name))) {
			this_.cap_0_locals.Add_1(j2gort.BoxString(n.GetString()))
		}
		if this_.cap_1_params.Contains(j2gort.BoxString(name)) {
			var type_ pkg_rhino.JSType = j2gort.AsIface[pkg_rhino.JSType](j2gort.MapGetStr(this_.cap_2_paramTypes, name))
			_ = type_
			var isStringLiteral bool = (type_.GetJSTypeBase().IsStringValueType() && javalang.StringStartsWith(j2gort.CheckNotNullString(name), "string_literal"))
			_ = isStringLiteral
			this_.outer_.ReplaceNodeInPlace(n, this_.outer_.priv_TemplateAstMatcher_createTemplateParameterNode_5324(this_.cap_1_params.IndexOf_1(j2gort.BoxString(name)), type_, isStringLiteral))
		} else if this_.cap_0_locals.Contains(j2gort.BoxString(name)) {
			this_.outer_.ReplaceNodeInPlace(n, this_.outer_.priv_TemplateAstMatcher_createTemplateLocalNameNode_5326(this_.cap_0_locals.IndexOf_1(j2gort.BoxString(name))))
		}
	}
}
