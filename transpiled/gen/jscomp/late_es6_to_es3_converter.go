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
// Port of com/google/javascript/jscomp/LateEs6ToEs3Converter.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.LateEs6ToEs3Converter
type LateEs6ToEs3Converter struct {
	compiler                  AbstractCompiler     `gson:"compiler"`
	astFactory                *AstFactory          `gson:"astFactory"`
	templateLiteralConverter  *Es6TemplateLiterals `gson:"templateLiteralConverter"`
	templateLitInsertionPoint pkg_rhino.Node       `gson:"templateLitInsertionPoint"`
	idHash_                   atomic.Int32
}

func (this_ *LateEs6ToEs3Converter) BindSelf_(s any) {
	_ = s
}
func (this_ *LateEs6ToEs3Converter) IsJavaNil_() bool { return this_ == nil }
func (*LateEs6ToEs3Converter) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.LateEs6ToEs3Converter", "LateEs6ToEs3Converter"
}
func (this_ *LateEs6ToEs3Converter) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*LateEs6ToEs3Converter) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.LateEs6ToEs3Converter", "com.google.javascript.jscomp.LateEs6ToEs3Converter", func(v_ any) bool {
		return j2gort.InstanceOf[*LateEs6ToEs3Converter](v_)
	})
}
func (*LateEs6ToEs3Converter) Implements_jscomp_NodeTraversal_Callback() {}
func (*LateEs6ToEs3Converter) Implements_jscomp_CompilerPass()           {}

var static_LateEs6ToEs3Converter_transpiledFeatures *pkg_parser.FeatureSet

func StaticGet_LateEs6ToEs3Converter_TranspiledFeatures() *pkg_parser.FeatureSet {
	clinit_LateEs6ToEs3Converter()
	return static_LateEs6ToEs3Converter_transpiledFeatures
}
func StaticSet_LateEs6ToEs3Converter_TranspiledFeatures(v *pkg_parser.FeatureSet) {
	clinit_LateEs6ToEs3Converter()
	static_LateEs6ToEs3Converter_transpiledFeatures = v
}
func StaticGet_LateEs6ToEs3Converter_FRESH_COMP_PROP_VAR() string { return "$jscomp$compprop" }

var clinitDone_LateEs6ToEs3Converter uint32
var clinitState_LateEs6ToEs3Converter j2gort.ClinitState

func clinit_LateEs6ToEs3Converter() {
	if atomic.LoadUint32(&clinitDone_LateEs6ToEs3Converter) != 0 {
		return
	}
	clinitSlow_LateEs6ToEs3Converter()
}

func clinitSlow_LateEs6ToEs3Converter() {
	if !j2gort.EnterClinit(&clinitState_LateEs6ToEs3Converter, &clinitDone_LateEs6ToEs3Converter) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_LateEs6ToEs3Converter, &clinitDone_LateEs6ToEs3Converter)
	static_LateEs6ToEs3Converter_transpiledFeatures = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature_7251([]*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_COMPUTED_PROPERTIES(), pkg_parser.StaticGet_FeatureSet_Feature_MEMBER_DECLARATIONS(), pkg_parser.StaticGet_FeatureSet_Feature_TEMPLATE_LITERALS()})
}

func (this_ *LateEs6ToEs3Converter) initFields_LateEs6ToEs3Converter() {
	_ = this_

}

// Java constructor: LateEs6ToEs3Converter(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewLateEs6ToEs3Converter(compiler AbstractCompiler) *LateEs6ToEs3Converter {
	clinit_LateEs6ToEs3Converter()
	this_ := &LateEs6ToEs3Converter{}
	InitLateEs6ToEs3Converter(this_, compiler)
	return this_
}

func InitLateEs6ToEs3Converter(this_ *LateEs6ToEs3Converter, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_LateEs6ToEs3Converter()
	this_.compiler = compiler
	this_.astFactory = compiler.CreateAstFactory()
	this_.templateLiteralConverter = NewEs6TemplateLiterals(compiler)
}

// Java method: LateEs6ToEs3Converter.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *LateEs6ToEs3Converter) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	TranspilationPasses_ProcessTranspile(this_.compiler, root, StaticGet_LateEs6ToEs3Converter_TranspiledFeatures(), []NodeTraversal_Callback{this_})
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(this_.compiler, root, StaticGet_LateEs6ToEs3Converter_TranspiledFeatures())
}

// Java method: LateEs6ToEs3Converter.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *LateEs6ToEs3Converter) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() {
		this_.templateLitInsertionPoint = LateEs6ToEs3Converter_FindTemplateLitInsertionPoint(n)
	}
	return true
}

// Java static method: LateEs6ToEs3Converter.findTemplateLitInsertionPoint(com.google.javascript.rhino.Node script)
func LateEs6ToEs3Converter_FindTemplateLitInsertionPoint(script pkg_rhino.Node) pkg_rhino.Node {
	clinit_LateEs6ToEs3Converter()
	_ = script
	if script.GetIsInClosureUnawareSubtree() {
		var closureUnawareBlock pkg_rhino.Node = NodeUtil_FindClosureUnawareScriptRoot(script)
		_ = closureUnawareBlock
		return NodeUtil_GetInsertionPointAfterAllInnerFunctionDeclarations(closureUnawareBlock)
	}
	return script.GetFirstChild()
}

// Java method: LateEs6ToEs3Converter.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *LateEs6ToEs3Converter) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(51) /* Token.ASSIGN */ :
		var lhs pkg_rhino.Node = n.GetFirstChild()
		_ = lhs
		var rhs pkg_rhino.Node = n.GetSecondChild()
		_ = rhs
		if (lhs.IsGetProp() && rhs.IsFunction()) && lhs.GetFirstChild().IsName() {
			var qName pkg_rhino.QualifiedName = pkg_rhino.QualifiedName_Of("$jscomp.createTemplateTagFirstArgWithRaw")
			_ = qName
			if qName.Matches(lhs) {
				j2gort.CheckNotNullNode(n.GetParent())
				if !(n.GetParent().IsExprResult()) {
					j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
				}
				this_.templateLitInsertionPoint = n.GetParent().GetNext()
			}
		}
	case int32(47) /* Token.OBJECTLIT */ :
		this_.priv_LateEs6ToEs3Converter_visitObject_3739(t, n)
	case int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
		if parent.IsObjectLit() {
			this_.priv_LateEs6ToEs3Converter_visitMemberFunctionDefInObjectLit_3738(n)
		}
	case int32(130) /* Token.TAGGED_TEMPLATELIT */ :
		this_.templateLiteralConverter.VisitTaggedTemplateLiteral(t, n, this_.templateLitInsertionPoint)
	case int32(131) /* Token.TEMPLATELIT */ :
		if !(parent.IsTaggedTemplateLit()) {
			this_.templateLiteralConverter.VisitTemplateLiteral(t, n)
		}
	default:
	}
}

// Java method: LateEs6ToEs3Converter.visitMemberFunctionDefInObjectLit(com.google.javascript.rhino.Node n)
func (this_ *LateEs6ToEs3Converter) priv_LateEs6ToEs3Converter_visitMemberFunctionDefInObjectLit_3738(n pkg_rhino.Node) {
	_ = n
	var name string = n.GetString()
	_ = name
	var nameNode pkg_rhino.Node = n.GetFirstFirstChild()
	_ = nameNode
	var stringKey pkg_rhino.Node = this_.astFactory.CreateStringKey(name, n.RemoveFirstChild())
	_ = stringKey
	stringKey.SetJSDocInfo(n.GetJSDocInfo())
	n.ReplaceWith(stringKey)
	stringKey.Srcref(nameNode)
	this_.compiler.ReportChangeToEnclosingScope(stringKey)
}

// Java method: LateEs6ToEs3Converter.visitObject(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node obj)
func (this_ *LateEs6ToEs3Converter) priv_LateEs6ToEs3Converter_visitObject_3739(t *NodeTraversal, obj pkg_rhino.Node) {
	_ = t
	_ = obj
	{
		var child pkg_rhino.Node = obj.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if child.IsComputedProp() {
				this_.priv_LateEs6ToEs3Converter_visitObjectWithComputedProperty_3740(t, obj)
				return
			}
		}
	}
}

// Java method: LateEs6ToEs3Converter.visitObjectWithComputedProperty(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node obj)
func (this_ *LateEs6ToEs3Converter) priv_LateEs6ToEs3Converter_visitObjectWithComputedProperty_3740(t *NodeTraversal, obj pkg_rhino.Node) {
	_ = t
	_ = obj
	if !(obj.IsObjectLit()) {
		j2gort.CheckArgument(false)
	}
	var props []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = props
	var currElement pkg_rhino.Node = obj.GetFirstChild()
	_ = currElement
	var objectType AstFactory_Type = AstFactory_S_Type_1_Node(obj)
	_ = objectType
	for currElement != 0 {
		if currElement.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_GETTER()) || currElement.GetBooleanProp(pkg_rhino.StaticGet_Node_COMPUTED_PROP_SETTER()) {
			this_.compiler.Report_1(JSError_Make_3(currElement, StaticGet_ReportUntranspilableFeatures_UNTRANSPILABLE_FEATURE_PRESENT(), []string{"computed getter/setter in an object literal", "ES2015", ""}))
			return
		} else if currElement.IsGetterDef() || currElement.IsSetterDef() {
			currElement = currElement.GetNext()
		} else {
			var nextNode pkg_rhino.Node = currElement.GetNext()
			_ = nextNode
			currElement.Detach()
			props = append(props, currElement)
			currElement = nextNode
		}
	}
	var objName string = ("$jscomp$compprop" + j2gort.StrValString(this_.compiler.GetUniqueIdSupplier().GetUniqueId(t.GetInput())))
	_ = objName
	props = j2gort.SliceReverseInPlace(props)
	var result pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(objName, objectType)
	_ = result
	for _, propdef := range props {
		_ = propdef
		if propdef.IsComputedProp() {
			var propertyExpression pkg_rhino.Node = propdef.RemoveFirstChild()
			_ = propertyExpression
			var value pkg_rhino.Node = propdef.RemoveFirstChild()
			_ = value
			result = this_.astFactory.CreateComma(this_.astFactory.CreateAssign_2_Node_Node(this_.astFactory.CreateGetElem(this_.astFactory.CreateName_2_string_AstFactory_Type(objName, objectType), propertyExpression), value), result)
		} else {
			var val pkg_rhino.Node = propdef.RemoveFirstChild()
			_ = val
			var isQuotedAccess bool = propdef.IsQuotedStringKey()
			_ = isQuotedAccess
			propdef.SetToken(pkg_rhino.StaticGet_Token_STRINGLIT())
			propdef.SetColor(pkg_colors.StaticGet_StandardColors_STRING())
			propdef.PutBooleanProp(pkg_rhino.StaticGet_Node_QUOTED_PROP(), false)
			var objNameNode pkg_rhino.Node = this_.astFactory.CreateName_2_string_AstFactory_Type(objName, objectType)
			_ = objNameNode
			var access pkg_rhino.Node
			_ = access
			if isQuotedAccess {
				access = this_.astFactory.CreateGetElem(objNameNode, propdef)
			} else {
				access = this_.astFactory.CreateGetProp(objNameNode, propdef.GetString(), AstFactory_S_Type_1_Node(propdef))
			}
			result = this_.astFactory.CreateComma(this_.astFactory.CreateAssign_2_Node_Node(access, val), result)
		}
	}

	var statement pkg_rhino.Node = obj
	_ = statement
	for !(NodeUtil_IsStatement(statement)) {
		statement = statement.GetParent()
	}
	var isInsideLoopHeader bool = NodeUtil_IsLoopStructure(statement)
	_ = isInsideLoopHeader
	if isInsideLoopHeader {
		var placeholder pkg_rhino.Node = pkg_rhino.IR_Empty()
		_ = placeholder
		obj.ReplaceWith(placeholder)
		var assignEmptyObj pkg_rhino.Node = this_.astFactory.CreateAssign_2_Node_Node(this_.astFactory.CreateName_2_string_AstFactory_Type(objName, objectType), obj)
		_ = assignEmptyObj
		result = this_.astFactory.CreateComma(assignEmptyObj, result)
		result.SrcrefTreeIfMissing(obj)
		placeholder.ReplaceWith(result)
		var var_ pkg_rhino.Node = pkg_rhino.IR_Var_1(this_.astFactory.CreateName_2_string_AstFactory_Type(objName, objectType))
		_ = var_
		var_.SrcrefTreeIfMissing(statement)
		var_.InsertBefore(statement)
		this_.compiler.ReportChangeToEnclosingScope(var_)
	} else {
		result.SrcrefTreeIfMissing(obj)
		obj.ReplaceWith(result)
		var var__2 pkg_rhino.Node = pkg_rhino.IR_Var_2(this_.astFactory.CreateName_2_string_AstFactory_Type(objName, objectType), obj)
		_ = var__2
		var__2.SrcrefTreeIfMissing(statement)
		var__2.InsertBefore(statement)
		this_.compiler.ReportChangeToEnclosingScope(var__2)
	}
}
func (*LateEs6ToEs3Converter) CanReleaseSyntacticScopes_() {}
