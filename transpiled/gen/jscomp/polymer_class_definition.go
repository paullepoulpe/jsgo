/*
 * Copyright 2016 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/PolymerClassDefinition.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.PolymerClassDefinition
type PolymerClassDefinition struct {
	input             *CompilerInput                                 `gson:"input"`
	constructor       *PolymerPass_MemberDefinition                  `gson:"constructor"`
	nativeBaseElement string                                         `gson:"nativeBaseElement"`
	props             j2gort.Ext_java_util_List                      `gson:"props"`
	behaviorProps     j2gort.Ext_java_util_Map                       `gson:"behaviorProps"`
	methods           j2gort.Ext_java_util_List                      `gson:"methods"`
	behaviors         []*PolymerBehaviorExtractor_BehaviorDefinition `gson:"behaviors"`
	features          *pkg_parser.FeatureSet                         `gson:"features"`
	interfaceName     string                                         `gson:"interfaceName"`
	definition        pkg_rhino.Node                                 `gson:"definition"`
	target            pkg_rhino.Node                                 `gson:"target"`
	descriptor        pkg_rhino.Node                                 `gson:"descriptor"`
	idHash_           atomic.Int32
	defType           PolymerClassDefinition_DefinitionType `gson:"defType"`
	hasGeneratedLhs   bool                                  `gson:"hasGeneratedLhs"`
}

func (this_ *PolymerClassDefinition) BindSelf_(s any) {
	_ = s
}
func (this_ *PolymerClassDefinition) IsJavaNil_() bool { return this_ == nil }
func (*PolymerClassDefinition) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolymerClassDefinition", "PolymerClassDefinition"
}
func (this_ *PolymerClassDefinition) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PolymerClassDefinition) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolymerClassDefinition", "com.google.javascript.jscomp.PolymerClassDefinition", func(v_ any) bool {
		return j2gort.InstanceOf[*PolymerClassDefinition](v_)
	})
}

func (this_ *PolymerClassDefinition) initFields_PolymerClassDefinition() {
	_ = this_
	this_.nativeBaseElement = j2gort.NullString
	this_.interfaceName = j2gort.NullString
}

// Java constructor: PolymerClassDefinition(com.google.javascript.jscomp.PolymerClassDefinition.DefinitionType defType, com.google.javascript.rhino.Node definition, com.google.javascript.rhino.Node target, boolean hasGeneratedLhs, com.google.javascript.rhino.Node descriptor, com.google.javascript.rhino.JSDocInfo classInfo, com.google.javascript.jscomp.PolymerPass.MemberDefinition constructor, java.lang.@org.jspecify.annotations.Nullable String nativeBaseElement, java.util.List<com.google.javascript.jscomp.PolymerPass.MemberDefinition> props, java.util.@org.jspecify.annotations.Nullable Map<com.google.javascript.jscomp.PolymerPass.MemberDefinition,com.google.javascript.jscomp.PolymerBehaviorExtractor.BehaviorDefinition> behaviorProps, java.util.List<com.google.javascript.jscomp.PolymerPass.MemberDefinition> methods, com.google.common.collect.@org.jspecify.annotations.Nullable ImmutableList<com.google.javascript.jscomp.PolymerBehaviorExtractor.BehaviorDefinition> behaviors, com.google.javascript.jscomp.parsing.parser.@org.jspecify.annotations.Nullable FeatureSet features, com.google.javascript.jscomp.CompilerInput input)
func NewPolymerClassDefinition(defType PolymerClassDefinition_DefinitionType, definition pkg_rhino.Node, target pkg_rhino.Node, hasGeneratedLhs bool, descriptor pkg_rhino.Node, classInfo *pkg_rhino.JSDocInfo, constructor *PolymerPass_MemberDefinition, nativeBaseElement string, props j2gort.Ext_java_util_List, behaviorProps j2gort.Ext_java_util_Map, methods j2gort.Ext_java_util_List, behaviors j2gort.Ext_com_google_common_collect_ImmutableList, features *pkg_parser.FeatureSet, input *CompilerInput) *PolymerClassDefinition {
	this_ := &PolymerClassDefinition{}
	InitPolymerClassDefinition(this_, defType, definition, target, hasGeneratedLhs, descriptor, classInfo, constructor, nativeBaseElement, props, behaviorProps, methods, behaviors, features, input)
	return this_
}

func InitPolymerClassDefinition(this_ *PolymerClassDefinition, defType PolymerClassDefinition_DefinitionType, definition pkg_rhino.Node, target pkg_rhino.Node, hasGeneratedLhs bool, descriptor pkg_rhino.Node, classInfo *pkg_rhino.JSDocInfo, constructor *PolymerPass_MemberDefinition, nativeBaseElement string, props j2gort.Ext_java_util_List, behaviorProps j2gort.Ext_java_util_Map, methods j2gort.Ext_java_util_List, behaviors j2gort.Ext_com_google_common_collect_ImmutableList, features *pkg_parser.FeatureSet, input *CompilerInput) {
	_ = defType
	_ = definition
	_ = target
	_ = hasGeneratedLhs
	_ = descriptor
	_ = classInfo
	_ = constructor
	_ = nativeBaseElement
	_ = props
	_ = behaviorProps
	_ = methods
	_ = behaviors
	_ = features
	_ = input
	this_.initFields_PolymerClassDefinition()
	this_.defType = defType
	this_.definition = definition
	this_.target = target
	this_.hasGeneratedLhs = hasGeneratedLhs
	if !((descriptor == 0) || descriptor.IsObjectLit()) {
		j2gort.CheckState(false)
	}
	this_.descriptor = descriptor
	this_.constructor = constructor
	this_.nativeBaseElement = nativeBaseElement
	this_.props = props
	this_.behaviorProps = behaviorProps
	this_.methods = methods
	this_.behaviors = j2gort.ImmListToSliceNullable[*PolymerBehaviorExtractor_BehaviorDefinition](behaviors)
	this_.features = features
	this_.input = input
}

// Java static method: PolymerClassDefinition.extractFromCallNode(com.google.javascript.rhino.Node callNode, com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata moduleMetadata, com.google.javascript.jscomp.PolymerBehaviorExtractor behaviorExtractor)
func PolymerClassDefinition_ExtractFromCallNode(callNode pkg_rhino.Node, compiler AbstractCompiler, moduleMetadata ModuleMetadataMap_ModuleMetadata, behaviorExtractor *PolymerBehaviorExtractor) *PolymerClassDefinition {
	_ = callNode
	_ = compiler
	_ = moduleMetadata
	_ = behaviorExtractor
	var descriptor pkg_rhino.Node = NodeUtil_GetArgumentForCallOrNew(callNode, int32(0))
	_ = descriptor
	if (descriptor == 0) || !(descriptor.IsObjectLit()) {
		compiler.Report_1(JSError_Make_3(callNode, StaticGet_PolymerPassErrors_POLYMER_DESCRIPTOR_NOT_VALID(), []string{}))
		return nil
	}
	var paramCount int32 = (callNode.GetChildCount() - int32(1))
	_ = paramCount
	if paramCount != int32(1) {
		compiler.Report_1(JSError_Make_3(callNode, StaticGet_PolymerPassErrors_POLYMER_UNEXPECTED_PARAMS(), []string{}))
		return nil
	}
	var elName pkg_rhino.Node = NodeUtil_GetFirstPropMatchingKey(descriptor, "is")
	_ = elName
	if elName == 0 {
		compiler.Report_1(JSError_Make_3(callNode, StaticGet_PolymerPassErrors_POLYMER_MISSING_IS(), []string{}))
		return nil
	}
	var enclosingModule pkg_rhino.Node = NodeUtil_GetEnclosingModuleIfPresent(callNode)
	_ = enclosingModule
	var hasGeneratedLhs bool = false
	_ = hasGeneratedLhs
	var target pkg_rhino.Node = pkg_rhino.Node(0)
	_ = target
	if NodeUtil_IsNameDeclaration(callNode.GetGrandparent()) {
		target = pkg_rhino.IR_Name(callNode.GetParent().GetString())
	} else if callNode.GetParent().IsAssign() {
		if PolymerClassDefinition_IsGoogModuleExports(callNode.GetParent()) {
			target = PolymerClassDefinition_CreateDummyGoogModuleExportsTarget(compiler, callNode)
		} else {
			target = callNode.GetParent().GetFirstChild().CloneTree_0()
		}
	} else {
		var elNameStringBase string
		_ = elNameStringBase
		if elName.IsQualifiedName() {
			elNameStringBase = j2gort.JStrReplaceChar(elName.GetQualifiedName(), uint16(46), uint16(36))
		} else {
			elNameStringBase = elName.GetString()
		}
		var elNameString string = j2gort.Ext_com_google_common_base_CaseFormat_F_LOWER_HYPHEN.To(j2gort.Ext_com_google_common_base_CaseFormat_F_UPPER_CAMEL, elNameStringBase)
		_ = elNameString
		elNameString = (j2gort.StrValString(elNameString) + "Element")
		target = pkg_rhino.IR_Name(elNameString)
		hasGeneratedLhs = true
	}
	var classInfo *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(target)
	_ = classInfo
	var ctorInfo *pkg_rhino.JSDocInfo = nil
	_ = ctorInfo
	var constructor pkg_rhino.Node = NodeUtil_GetFirstPropMatchingKey(descriptor, "factoryImpl")
	_ = constructor
	if constructor == 0 {
		constructor = NodeUtil_EmptyFunction()
		compiler.ReportChangeToChangeScope(constructor)
		constructor.SrcrefTree(callNode)
	} else {
		ctorInfo = NodeUtil_GetBestJSDocInfo(constructor)
	}
	var baseClass pkg_rhino.Node = NodeUtil_GetFirstPropMatchingKey(descriptor, "extends")
	_ = baseClass
	var nativeBaseElement string
	_ = nativeBaseElement
	if baseClass == 0 {
		nativeBaseElement = j2gort.NullString
	} else {
		nativeBaseElement = baseClass.GetString()
	}
	var behaviorArray pkg_rhino.Node = NodeUtil_GetFirstPropMatchingKey(descriptor, "behaviors")
	_ = behaviorArray
	var behaviors []*PolymerBehaviorExtractor_BehaviorDefinition = j2gort.ImmListToSliceNullable[*PolymerBehaviorExtractor_BehaviorDefinition](behaviorExtractor.ExtractBehaviors(behaviorArray, moduleMetadata))
	_ = behaviors
	var properties j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = properties
	var behaviorProps j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = behaviorProps
	{

		for _, behavior := range behaviors {

			_ = behavior
			{
				cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(behavior.props))
				for cur_3.HasNext() {
					rawElem_4 := cur_3.Next()
					prop := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_4)
					_ = prop
					behaviorProps.Put_2(j2gort.PtrToAny(prop), j2gort.PtrToAny(behavior))
				}
			}
		}
	}
	PolymerClassDefinition_OverwriteMembersIfPresent(properties, j2gort.AsIface[j2gort.Ext_java_util_List](PolymerPassStaticUtils_ExtractProperties(descriptor, StaticGet_PolymerClassDefinition_DefinitionType_ObjectLiteral(), compiler, pkg_rhino.Node(0))))
	PolymerClassDefinition_RemoveDuplicateBehaviorProps(behaviorProps)
	PolymerClassDefinition_RemoveBehaviorPropsOverlappingWithElementProps(behaviorProps, properties)
	var newFeatures *pkg_parser.FeatureSet = nil
	_ = newFeatures
	if len(behaviors) != 0 {
		newFeatures = behaviors[int32(0)].features
		{
			var i int32 = int32(1)
			_ = i
			for ; i < int32(len(behaviors)); i = i + 1 {
				newFeatures = newFeatures.Union(behaviors[i].features)
			}
		}
	}
	var methods j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = methods
	{
		var keyNode pkg_rhino.Node = descriptor.GetFirstChild()
		_ = keyNode
		for ; keyNode != 0; keyNode = keyNode.GetNext() {
			var isFunctionDefinition bool = (keyNode.IsMemberFunctionDef() || (keyNode.IsStringKey() && keyNode.GetFirstChild().IsFunction()))
			_ = isFunctionDefinition
			if isFunctionDefinition {
				methods.Add_1(NewPolymerPass_MemberDefinition(NodeUtil_GetBestJSDocInfo(keyNode), keyNode, keyNode.GetFirstChild(), enclosingModule))
			}
		}
	}
	var input *CompilerInput = compiler.GetInput(NodeUtil_GetEnclosingScript(callNode).GetInputId())
	_ = input
	return NewPolymerClassDefinition(StaticGet_PolymerClassDefinition_DefinitionType_ObjectLiteral(), callNode, target, hasGeneratedLhs, descriptor, classInfo, NewPolymerPass_MemberDefinition(ctorInfo, pkg_rhino.Node(0), constructor, enclosingModule), nativeBaseElement, properties, behaviorProps, methods, j2gort.ImmListFromTypedSlice(behaviors), newFeatures, input)
}

// Java static method: PolymerClassDefinition.isGoogModuleExports(com.google.javascript.rhino.Node assign)
func PolymerClassDefinition_IsGoogModuleExports(assign pkg_rhino.Node) bool {
	_ = assign
	if !(assign.GetParent().IsExprResult()) || !(assign.GetFirstChild().MatchesName_1_string("exports")) {
		return false
	}
	var containingBlock pkg_rhino.Node = assign.GetGrandparent()
	_ = containingBlock
	return ((containingBlock.IsModuleBody() && containingBlock.GetParent().GetBooleanProp(pkg_rhino.StaticGet_Node_GOOG_MODULE())) || NodeUtil_IsBundledGoogModuleScopeRoot(containingBlock))
}

// Java static method: PolymerClassDefinition.createDummyGoogModuleExportsTarget(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node callNode)
func PolymerClassDefinition_CreateDummyGoogModuleExportsTarget(compiler AbstractCompiler, callNode pkg_rhino.Node) pkg_rhino.Node {
	_ = compiler
	_ = callNode
	var madeUpName string = ("exportsForPolymer$jscomp" + j2gort.StrValString(j2gort.CastString(j2gort.AsIface[j2gort.Ext_java_util_function_Supplier](compiler.GetUniqueNameIdSupplier()).Get())))
	_ = madeUpName
	var assignExpr pkg_rhino.Node = callNode.GetGrandparent()
	_ = assignExpr
	var exportName pkg_rhino.Node = callNode.GetParent().GetFirstChild()
	_ = exportName
	var target pkg_rhino.Node = pkg_rhino.IR_Name(madeUpName).ClonePropsFrom(exportName).Srcref(exportName)
	_ = target
	callNode.ReplaceWith(target)
	var newDecl pkg_rhino.Node = pkg_rhino.IR_Var_2(target.CloneNode_0(), callNode).Srcref(assignExpr)
	_ = newDecl
	newDecl.InsertBefore(assignExpr)
	newDecl.SetJSDocInfo(assignExpr.GetJSDocInfo())
	assignExpr.SetJSDocInfo(nil)
	return target
}

// Java static method: PolymerClassDefinition.extractFromClassNode(com.google.javascript.rhino.Node classNode, com.google.javascript.jscomp.AbstractCompiler compiler)
func PolymerClassDefinition_ExtractFromClassNode(classNode pkg_rhino.Node, compiler AbstractCompiler) *PolymerClassDefinition {
	_ = classNode
	_ = compiler
	if !((classNode != 0) && classNode.IsClass()) {
		j2gort.CheckState(false)
	}
	var propertiesDescriptor pkg_rhino.Node = pkg_rhino.Node(0)
	_ = propertiesDescriptor
	var propertiesGetter pkg_rhino.Node = NodeUtil_GetFirstGetterMatchingKey(NodeUtil_GetClassMembers(classNode), "properties")
	_ = propertiesGetter
	if propertiesGetter != 0 {
		if !(propertiesGetter.IsStaticMember()) {
			compiler.Report_1(JSError_Make_3(classNode, StaticGet_PolymerPassErrors_POLYMER_CLASS_PROPERTIES_NOT_STATIC(), []string{}))
		} else {
			{
				var child pkg_rhino.Node = NodeUtil_GetFunctionBody(propertiesGetter.GetFirstChild()).GetFirstChild()
				_ = child
				for ; child != 0; child = child.GetNext() {
					if child.IsReturn() {
						if child.HasChildren() && child.GetFirstChild().IsObjectLit() {
							propertiesDescriptor = child.GetFirstChild()
							break
						} else {
							compiler.Report_1(JSError_Make_3(propertiesGetter, StaticGet_PolymerPassErrors_POLYMER_CLASS_PROPERTIES_INVALID(), []string{}))
						}
					}
				}
			}
		}
	}
	var target pkg_rhino.Node = pkg_rhino.Node(0)
	_ = target
	if NodeUtil_IsNameDeclaration(classNode.GetGrandparent()) {
		target = pkg_rhino.IR_Name(classNode.GetParent().GetString())
	} else if classNode.GetParent().IsAssign() && classNode.GetParent().GetFirstChild().IsQualifiedName() {
		target = classNode.GetParent().GetFirstChild()
	} else if !(classNode.GetFirstChild().IsEmpty()) {
		target = classNode.GetFirstChild()
	} else {
		compiler.Report_1(JSError_Make_3(classNode, StaticGet_PolymerPassErrors_POLYMER_CLASS_UNNAMED(), []string{}))
		return nil
	}
	var enclosingModule pkg_rhino.Node = NodeUtil_GetEnclosingModuleIfPresent(classNode)
	_ = enclosingModule
	var classInfo *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(classNode)
	_ = classInfo
	var constructor pkg_rhino.Node = NodeUtil_GetEs6ClassConstructorMemberFunctionDef(classNode)
	_ = constructor
	var properties j2gort.Ext_com_google_common_collect_ImmutableList = PolymerPassStaticUtils_ExtractProperties(propertiesDescriptor, StaticGet_PolymerClassDefinition_DefinitionType_ES6Class(), compiler, constructor)
	_ = properties
	var methods j2gort.Ext_java_util_List = PolymerClassDefinition_ExtractClassMethods(classNode, enclosingModule)
	_ = methods
	var ctorInfo *pkg_rhino.JSDocInfo = nil
	_ = ctorInfo
	if constructor != 0 {
		ctorInfo = NodeUtil_GetBestJSDocInfo(constructor)
		PolymerClassDefinition_AddFunctionAssignmentsFromConstructor(constructor, enclosingModule, methods)
	}
	var input *CompilerInput = compiler.GetInput(NodeUtil_GetEnclosingScript(classNode).GetInputId())
	_ = input
	return NewPolymerClassDefinition(StaticGet_PolymerClassDefinition_DefinitionType_ES6Class(), classNode, target, false, propertiesDescriptor, classInfo, NewPolymerPass_MemberDefinition(ctorInfo, pkg_rhino.Node(0), constructor, enclosingModule), j2gort.NullString, j2gort.AsIface[j2gort.Ext_java_util_List](properties), nil, methods, nil, nil, input)
}

// Java static method: PolymerClassDefinition.extractClassMethods(com.google.javascript.rhino.Node classNode, com.google.javascript.rhino.Node enclosingModule)
func PolymerClassDefinition_ExtractClassMethods(classNode pkg_rhino.Node, enclosingModule pkg_rhino.Node) j2gort.Ext_java_util_List {
	_ = classNode
	_ = enclosingModule
	var methods []*PolymerPass_MemberDefinition = []*PolymerPass_MemberDefinition{}
	_ = methods
	{
		var keyNode pkg_rhino.Node = NodeUtil_GetClassMembers(classNode).GetFirstChild()
		_ = keyNode
		for ; keyNode != 0; keyNode = keyNode.GetNext() {
			if keyNode.IsMemberFunctionDef() || PolymerClassDefinition_IsFunctionField(keyNode) {
				methods = append(methods, NewPolymerPass_MemberDefinition(NodeUtil_GetBestJSDocInfo(keyNode), keyNode, keyNode.GetFirstChild(), enclosingModule))
			}
		}
	}
	return j2gort.WrapTypedArrayList(methods)
}

// Java static method: PolymerClassDefinition.isFunctionField(com.google.javascript.rhino.Node field)
func PolymerClassDefinition_IsFunctionField(field pkg_rhino.Node) bool {
	_ = field
	if !(field.IsMemberFieldDef()) {
		return false
	}
	var value pkg_rhino.Node = field.GetFirstChild()
	_ = value
	return ((value != 0) && value.IsFunction())
}

// Java static method: PolymerClassDefinition.addFunctionAssignmentsFromConstructor(com.google.javascript.rhino.Node constructor, com.google.javascript.rhino.Node enclosingModule, java.util.List<com.google.javascript.jscomp.PolymerPass.MemberDefinition> methods)
func PolymerClassDefinition_AddFunctionAssignmentsFromConstructor(constructor pkg_rhino.Node, enclosingModule pkg_rhino.Node, methods j2gort.Ext_java_util_List) {
	_ = constructor
	_ = enclosingModule
	_ = methods
	var body pkg_rhino.Node = NodeUtil_GetFunctionBody(constructor.GetFirstChild())
	_ = body
	{
		var stmt pkg_rhino.Node = body.GetFirstChild()
		_ = stmt
		for ; stmt != 0; stmt = stmt.GetNext() {
			if !(stmt.IsExprResult()) || !(stmt.GetFirstChild().IsAssign()) {
				continue
			}
			var assignment pkg_rhino.Node = stmt.GetFirstChild()
			_ = assignment
			var propName pkg_rhino.Node = assignment.GetFirstChild()
			_ = propName
			var propValue pkg_rhino.Node = assignment.GetLastChild()
			_ = propValue
			if (propName.IsGetProp() && propName.GetFirstChild().IsThis()) && propValue.IsFunction() {
				methods.Add_1(NewPolymerPass_MemberDefinition(NodeUtil_GetBestJSDocInfo(propName), propName, propValue, enclosingModule))
			}
		}
	}
}

// Java static method: PolymerClassDefinition.overwriteMembersIfPresent(java.util.List<com.google.javascript.jscomp.PolymerPass.MemberDefinition> list, java.util.List<com.google.javascript.jscomp.PolymerPass.MemberDefinition> newMembers)
func PolymerClassDefinition_OverwriteMembersIfPresent(list j2gort.Ext_java_util_List, newMembers j2gort.Ext_java_util_List) {
	_ = list
	_ = newMembers
	{
		cur_5 := j2gort.NewIterCursor(j2gort.CheckNotNull(newMembers))
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			newMember := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_6)
			_ = newMember
			{
				cur_7 := j2gort.NewIterCursor(j2gort.CheckNotNull(list))
				for cur_7.HasNext() {
					rawElem_8 := cur_7.Next()
					member := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_8)
					_ = member
					{
						var iifeCond_1 bool
						{

							sr_9 := member.name.GetString()
							se_10 := newMember.name.GetString()
							iifeCond_1 = j2gort.CheckNotNullString(sr_9) == se_10
						}
						if iifeCond_1 {
							list.Remove_1_any(j2gort.PtrToAny(member))
							break
						}
					}
				}
			}
			list.Add_1(j2gort.PtrToAny(newMember))
		}
	}
}

// Java static method: PolymerClassDefinition.removeDuplicateBehaviorProps(java.util.Map<com.google.javascript.jscomp.PolymerPass.MemberDefinition,com.google.javascript.jscomp.PolymerBehaviorExtractor.BehaviorDefinition> behaviorProps)
func PolymerClassDefinition_RemoveDuplicateBehaviorProps(behaviorProps j2gort.Ext_java_util_Map) {
	_ = behaviorProps
	if behaviorProps == nil {
		return
	}
	var behaviorsItr j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](j2gort.AsIface[j2gort.Ext_java_util_Set](behaviorProps.EntrySet()).Iterator())
	_ = behaviorsItr
	var seen j2gort.LocalOrderedSet[string] = j2gort.NewLocalOrderedSet[string]()
	_ = seen
	for behaviorsItr.HasNext() {
		var memberDefinition *PolymerPass_MemberDefinition = j2gort.CastPtr[PolymerPass_MemberDefinition](j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](behaviorsItr.Next()).GetKey())
		_ = memberDefinition
		var propertyName string = memberDefinition.name.GetString()
		_ = propertyName
		if !(seen.Add(propertyName)) {
			behaviorsItr.Remove()
		}
	}
}

// Java static method: PolymerClassDefinition.removeBehaviorPropsOverlappingWithElementProps(java.util.Map<com.google.javascript.jscomp.PolymerPass.MemberDefinition,com.google.javascript.jscomp.PolymerBehaviorExtractor.BehaviorDefinition> behaviorProps, java.util.List<com.google.javascript.jscomp.PolymerPass.MemberDefinition> polymerElementProps)
func PolymerClassDefinition_RemoveBehaviorPropsOverlappingWithElementProps(behaviorProps j2gort.Ext_java_util_Map, polymerElementProps j2gort.Ext_java_util_List) {
	_ = behaviorProps
	_ = polymerElementProps
	if behaviorProps == nil {
		return
	}
	var elementPropNames j2gort.LocalOrderedSet[string] = j2gort.NewLocalOrderedSet[string]()
	_ = elementPropNames
	{
		cur_11 := j2gort.NewIterCursor(j2gort.CheckNotNull(polymerElementProps))
		for cur_11.HasNext() {
			rawElem_12 := cur_11.Next()
			x := j2gort.CastPtr[PolymerPass_MemberDefinition](rawElem_12)
			_ = x
			elementPropNames.Add(x.name.GetString())
		}
	}
	var behaviorsItr j2gort.Ext_java_util_Iterator = j2gort.AsIface[j2gort.Ext_java_util_Iterator](j2gort.AsIface[j2gort.Ext_java_util_Set](behaviorProps.EntrySet()).Iterator())
	_ = behaviorsItr
	for behaviorsItr.HasNext() {
		var memberDefinition *PolymerPass_MemberDefinition = j2gort.CastPtr[PolymerPass_MemberDefinition](j2gort.AsIface[j2gort.Ext_java_util_Map_Entry](behaviorsItr.Next()).GetKey())
		_ = memberDefinition
		if elementPropNames.Contains(memberDefinition.name.GetString()) {
			behaviorsItr.Remove()
		}
	}
}

// Java method: PolymerClassDefinition.toString()
func (this_ *PolymerClassDefinition) String() string {
	return j2gort.Ext_com_google_common_base_MoreObjects_S_ToStringHelper_1_any(this_).Add_2_string_any("defType", j2gort.BoxNullableVal(this_.defType)).Add_2_string_any("definition", j2gort.BoxNode(this_.definition)).Add_2_string_any("target", j2gort.BoxNode(this_.target)).Add_2_string_any("nativeBaseElement", j2gort.BoxString(this_.nativeBaseElement)).OmitNullValues().String()
}

// Java method: PolymerClassDefinition.getInterfaceName(com.google.javascript.jscomp.UniqueIdSupplier uniqueIdSupplier)
func (this_ *PolymerClassDefinition) GetInterfaceName(uniqueIdSupplier *UniqueIdSupplier) string {
	_ = uniqueIdSupplier
	if this_.interfaceName == j2gort.NullString {
		this_.interfaceName = ("Polymer" + j2gort.StrValString(j2gort.JStrReplaceChar(this_.target.GetQualifiedName(), uint16(46), uint16(95))) + "Interface$" + j2gort.StrValString(uniqueIdSupplier.GetUniqueId(this_.input)))
	}
	return this_.interfaceName
}

// Java enum: com.google.javascript.jscomp.PolymerClassDefinition.DefinitionType
type PolymerClassDefinition_DefinitionType uint8

const (
	enum_PolymerClassDefinition_DefinitionType_ObjectLiteral PolymerClassDefinition_DefinitionType = 1
	enum_PolymerClassDefinition_DefinitionType_ES6Class      PolymerClassDefinition_DefinitionType = 2
)

var enumBaseTable_PolymerClassDefinition_DefinitionType = [...]j2gort.EnumBase{
	{Name_: "ObjectLiteral", Ordinal_: 0},
	{Name_: "ES6Class", Ordinal_: 1},
}
var enumValues_PolymerClassDefinition_DefinitionType = []PolymerClassDefinition_DefinitionType{1, 2}

func StaticGet_PolymerClassDefinition_DefinitionType_ObjectLiteral() PolymerClassDefinition_DefinitionType {
	return enum_PolymerClassDefinition_DefinitionType_ObjectLiteral
}
func StaticGet_PolymerClassDefinition_DefinitionType_ES6Class() PolymerClassDefinition_DefinitionType {
	return enum_PolymerClassDefinition_DefinitionType_ES6Class
}

func PolymerClassDefinition_DefinitionType_Values() []PolymerClassDefinition_DefinitionType {
	return j2gort.CloneSlice(enumValues_PolymerClassDefinition_DefinitionType)
}

func PolymerClassDefinition_DefinitionType_ValueOf(name string) PolymerClassDefinition_DefinitionType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_PolymerClassDefinition_DefinitionType {
		if enumBaseTable_PolymerClassDefinition_DefinitionType[i_].Name_ == name {
			return PolymerClassDefinition_DefinitionType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.PolymerClassDefinition.DefinitionType." + name))
}

func (this_ PolymerClassDefinition_DefinitionType) IsJavaNil_() bool                              { return this_ == 0 }
func (PolymerClassDefinition_DefinitionType) HasIdentityEquals_()                                 {}
func (PolymerClassDefinition_DefinitionType) Implements_j2gort_Ext_java_lang_Enum()               {}
func (PolymerClassDefinition_DefinitionType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ PolymerClassDefinition_DefinitionType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (PolymerClassDefinition_DefinitionType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolymerClassDefinition$DefinitionType", "DefinitionType"
}
func (this_ PolymerClassDefinition_DefinitionType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_PolymerClassDefinition_DefinitionType[this_-1]
}
func (this_ PolymerClassDefinition_DefinitionType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ PolymerClassDefinition_DefinitionType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ PolymerClassDefinition_DefinitionType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_PolymerClassDefinition_DefinitionType[this_-1].Name_
}
func (this_ PolymerClassDefinition_DefinitionType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_PolymerClassDefinition_DefinitionType[this_-1].Name_
}
func (this_ PolymerClassDefinition_DefinitionType) Equals(other_ any) bool {
	o_, ok_ := other_.(PolymerClassDefinition_DefinitionType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ PolymerClassDefinition_DefinitionType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_PolymerClassDefinition_DefinitionType[this_-1].HashCode()
}
func (this_ PolymerClassDefinition_DefinitionType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_PolymerClassDefinition_DefinitionType[this_-1].CompareTo(other_)
}
func (this_ PolymerClassDefinition_DefinitionType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_PolymerClassDefinition_DefinitionType[this_-1].Name_)
}
func init() {
	vals_ := []any{PolymerClassDefinition_DefinitionType(1), PolymerClassDefinition_DefinitionType(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.PolymerClassDefinition.DefinitionType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.PolymerClassDefinition$DefinitionType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolymerClassDefinition.DefinitionType", "com.google.javascript.jscomp.PolymerClassDefinition$DefinitionType", func(v_ any) bool {
		return j2gort.InstanceOf[PolymerClassDefinition_DefinitionType](v_)
	})
}
