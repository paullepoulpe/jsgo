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
// Port of com/google/javascript/jscomp/PolymerPassStaticUtils.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.PolymerPassStaticUtils
type PolymerPassStaticUtils struct {
	idHash_ atomic.Int32
}

func (this_ *PolymerPassStaticUtils) BindSelf_(s any) {
	_ = s
}
func (this_ *PolymerPassStaticUtils) IsJavaNil_() bool { return this_ == nil }
func (*PolymerPassStaticUtils) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolymerPassStaticUtils", "PolymerPassStaticUtils"
}
func (this_ *PolymerPassStaticUtils) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PolymerPassStaticUtils) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.PolymerPassStaticUtils", "com.google.javascript.jscomp.PolymerPassStaticUtils", func(v_ any) bool {
		return j2gort.InstanceOf[*PolymerPassStaticUtils](v_)
	})
}

func StaticGet_PolymerPassStaticUtils_VIRTUAL_FILE() string { return "<PolymerPassStaticUtils.java>" }

var static_PolymerPassStaticUtils_POLYMER_DOT_ELEMENT pkg_rhino.QualifiedName

func StaticGet_PolymerPassStaticUtils_POLYMER_DOT_ELEMENT() pkg_rhino.QualifiedName {
	clinit_PolymerPassStaticUtils()
	return static_PolymerPassStaticUtils_POLYMER_DOT_ELEMENT
}
func StaticSet_PolymerPassStaticUtils_POLYMER_DOT_ELEMENT(v pkg_rhino.QualifiedName) {
	clinit_PolymerPassStaticUtils()
	static_PolymerPassStaticUtils_POLYMER_DOT_ELEMENT = v
}

var clinitDone_PolymerPassStaticUtils uint32
var clinitState_PolymerPassStaticUtils j2gort.ClinitState

func clinit_PolymerPassStaticUtils() {
	if atomic.LoadUint32(&clinitDone_PolymerPassStaticUtils) != 0 {
		return
	}
	clinitSlow_PolymerPassStaticUtils()
}

func clinitSlow_PolymerPassStaticUtils() {
	if !j2gort.EnterClinit(&clinitState_PolymerPassStaticUtils, &clinitDone_PolymerPassStaticUtils) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_PolymerPassStaticUtils, &clinitDone_PolymerPassStaticUtils)
	static_PolymerPassStaticUtils_POLYMER_DOT_ELEMENT = pkg_rhino.QualifiedName_Of("Polymer.Element")
}

func (this_ *PolymerPassStaticUtils) initFields_PolymerPassStaticUtils() {
	_ = this_
}

// Java constructor: PolymerPassStaticUtils()
func NewPolymerPassStaticUtils() *PolymerPassStaticUtils {
	this_ := &PolymerPassStaticUtils{}
	InitPolymerPassStaticUtils(this_)
	return this_
}

func InitPolymerPassStaticUtils(this_ *PolymerPassStaticUtils) {
	this_.initFields_PolymerPassStaticUtils()
}

// Java static method: PolymerPassStaticUtils.isPolymerCall(com.google.javascript.rhino.Node call)
func PolymerPassStaticUtils_IsPolymerCall(call pkg_rhino.Node) bool {
	_ = call
	if (call == 0) || !(call.IsCall()) {
		return false
	}
	var name pkg_rhino.Node = call.GetFirstChild()
	_ = name
	return ((name.MatchesName_1_string("Polymer") || func() bool {
		sr_1 := "Polymer"
		se_2 := name.GetOriginalQualifiedName()
		return j2gort.CheckNotNullString(sr_1) == se_2
	}()) || (name.IsGetProp() && (j2gort.CheckNotNullString(name.GetString()) == "Polymer")))
}

// Java static method: PolymerPassStaticUtils.isPolymerClass(com.google.javascript.rhino.Node cls)
func PolymerPassStaticUtils_IsPolymerClass(cls pkg_rhino.Node) bool {
	_ = cls
	if (cls == 0) || !(cls.IsClass()) {
		return false
	}
	var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(cls)
	_ = info
	if (info != nil) && info.IsPolymer() {
		return true
	}
	var heritage pkg_rhino.Node = cls.GetSecondChild()
	_ = heritage
	return (!(heritage.IsEmpty()) && (((((StaticGet_PolymerPassStaticUtils_POLYMER_DOT_ELEMENT().Matches(heritage) || PolymerPassStaticUtils_Matches(heritage, "PolymerElement")) || PolymerPassStaticUtils_Matches(heritage, "PolymerElementWithoutHtml")) || PolymerPassStaticUtils_Matches(heritage, "PolymerLiteControllerBase")) || PolymerPassStaticUtils_Matches(heritage, "PolymerBehavior")) || PolymerPassStaticUtils_Matches(heritage, "PolymerBehaviorWithoutHtml")))
}

// Java static method: PolymerPassStaticUtils.matches(com.google.javascript.rhino.Node n, java.lang.String name)
func PolymerPassStaticUtils_Matches(n pkg_rhino.Node, name string) bool {
	_ = n
	_ = name
	return ((n.MatchesName_1_string(name) || func() bool {
		sr_3 := name
		se_4 := n.GetOriginalQualifiedName()
		return j2gort.CheckNotNullString(sr_3) == se_4
	}()) || (n.IsGetProp() && (j2gort.CheckNotNullString(n.GetString()) == name)))
}

// Java static method: PolymerPassStaticUtils.switchDollarSignPropsToBrackets(com.google.javascript.rhino.Node def, com.google.javascript.jscomp.AbstractCompiler compiler)
func PolymerPassStaticUtils_SwitchDollarSignPropsToBrackets(def pkg_rhino.Node, compiler AbstractCompiler) {
	_ = def
	_ = compiler
	if !(def.IsObjectLit() || def.IsClassMembers()) {
		j2gort.CheckState(false)
	}
	{
		var keyNode pkg_rhino.Node = def.GetFirstChild()
		_ = keyNode
		for ; keyNode != 0; keyNode = keyNode.GetNext() {
			var value pkg_rhino.Node = keyNode.GetFirstChild()
			_ = value
			if (value != 0) && value.IsFunction() {
				NodeUtil_VisitPostOrder_2(value.GetLastChild(), NewPolymerPassStaticUtils__Anon_70(compiler))
			}
		}
	}
}

// Java static method: PolymerPassStaticUtils.quoteListenerAndHostAttributeKeys(com.google.javascript.rhino.Node objLit, com.google.javascript.jscomp.AbstractCompiler compiler)
func PolymerPassStaticUtils_QuoteListenerAndHostAttributeKeys(objLit pkg_rhino.Node, compiler AbstractCompiler) {
	_ = objLit
	_ = compiler
	if !(objLit.IsObjectLit()) {
		j2gort.CheckState(false)
	}
	{
		var keyNode pkg_rhino.Node = objLit.GetFirstChild()
		_ = keyNode
		for ; keyNode != 0; keyNode = keyNode.GetNext() {
			if !(keyNode.IsStringKey()) {
				continue
			}
			if !(j2gort.CheckNotNullString(keyNode.GetString()) == "listeners") && !(j2gort.CheckNotNullString(keyNode.GetString()) == "hostAttributes") {
				continue
			}
			{
				var keyToQuote pkg_rhino.Node = keyNode.GetFirstFirstChild()
				_ = keyToQuote
				for ; keyToQuote != 0; keyToQuote = keyToQuote.GetNext() {
					if !(keyToQuote.IsQuotedStringKey()) {
						keyToQuote.SetQuotedStringKey()
						compiler.ReportChangeToEnclosingScope(keyToQuote)
					}
				}
			}
		}
	}
}

// Java static method: PolymerPassStaticUtils.protectObserverAndPropertyFunctionKeys(com.google.javascript.rhino.Node behaviorObjLit)
func PolymerPassStaticUtils_ProtectObserverAndPropertyFunctionKeys(behaviorObjLit pkg_rhino.Node) {
	_ = behaviorObjLit
	if !(behaviorObjLit.IsObjectLit()) {
		j2gort.CheckState(false)
	}
	var methodsToProtect j2gort.OrderedSet[string] = j2gort.NewOrderedSet[string]()
	_ = methodsToProtect
	methodsToProtect.AddAllSet(ordSet_PolymerPassStaticUtils_AddObserverMethodsToProtect(behaviorObjLit))
	methodsToProtect.AddAllSet(ordSet_PolymerPassStaticUtils_AddPropertyMethodsToProtect(behaviorObjLit))
	{
		var keyNode pkg_rhino.Node = behaviorObjLit.GetFirstChild()
		_ = keyNode
		for ; keyNode != 0; keyNode = keyNode.GetNext() {
			if !(NodeUtil_IsObjectLitKey(keyNode)) {
				continue
			}
			var keyName string = NodeUtil_GetObjectOrClassLitKeyName(keyNode)
			_ = keyName
			if !(methodsToProtect.Contains(keyName)) {
				continue
			}
			PolymerPassStaticUtils_AddExportAndNoCollapseToKey(keyNode)
		}
	}
}

// Java static method: PolymerPassStaticUtils.addObserverMethodsToProtect(com.google.javascript.rhino.Node behaviorObjLit)
func ordSet_PolymerPassStaticUtils_AddObserverMethodsToProtect(behaviorObjLit pkg_rhino.Node) j2gort.OrderedSet[string] {
	_ = behaviorObjLit
	var methodsToProtect j2gort.OrderedSet[string] = j2gort.NewOrderedSet[string]()
	_ = methodsToProtect
	var observers pkg_rhino.Node = NodeUtil_GetFirstPropMatchingKey(behaviorObjLit, "observers")
	_ = observers
	if (observers != 0) && observers.IsArrayLit() {
		{
			var child pkg_rhino.Node = observers.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				if !(child.IsStringLit()) {
					continue
				}
				var name string = PolymerPassStaticUtils_ExtractMethodName(child.GetString())
				_ = name
				if name != j2gort.NullString {
					methodsToProtect.Add(name)
				}
			}
		}
	}
	return methodsToProtect
}

// Java static method: PolymerPassStaticUtils.addPropertyMethodsToProtect(com.google.javascript.rhino.Node behaviorObjLit)
func ordSet_PolymerPassStaticUtils_AddPropertyMethodsToProtect(behaviorObjLit pkg_rhino.Node) j2gort.OrderedSet[string] {
	_ = behaviorObjLit
	var methodsToProtect j2gort.OrderedSet[string] = j2gort.NewOrderedSet[string]()
	_ = methodsToProtect
	var properties pkg_rhino.Node = NodeUtil_GetFirstPropMatchingKey(behaviorObjLit, "properties")
	_ = properties
	if (properties != 0) && properties.IsObjectLit() {
		{
			var prop pkg_rhino.Node = properties.GetFirstChild()
			_ = prop
			for ; prop != 0; prop = prop.GetNext() {
				if !(prop.GetLastChild().IsObjectLit()) {
					continue
				}
				var observer pkg_rhino.Node = NodeUtil_GetFirstPropMatchingKey(prop.GetLastChild(), "observer")
				_ = observer
				if (observer != 0) && observer.IsStringLit() {
					methodsToProtect.Add(observer.GetString())
				}
				var computed pkg_rhino.Node = NodeUtil_GetFirstPropMatchingKey(prop.GetLastChild(), "computed")
				_ = computed
				if (computed != 0) && computed.IsStringLit() {
					var name string = PolymerPassStaticUtils_ExtractMethodName(computed.GetString())
					_ = name
					if name != j2gort.NullString {
						methodsToProtect.Add(name)
					}
				}
			}
		}
	}
	return methodsToProtect
}

// Java static method: PolymerPassStaticUtils.extractMethodName(java.lang.String signature)
func PolymerPassStaticUtils_ExtractMethodName(signature string) string {
	_ = signature
	var openParenIndex int32 = j2gort.JStrIndexOfCp(signature, int32(uint16(40)))
	_ = openParenIndex
	if openParenIndex > int32(0) {
		return javalang.StringTrim(j2gort.CheckNotNullString(j2gort.JStrSub(signature, int32(0), openParenIndex)))
	} else {
		return j2gort.NullString
	}
}

// Java static method: PolymerPassStaticUtils.addExportAndNoCollapseToKey(com.google.javascript.rhino.Node keyNode)
func PolymerPassStaticUtils_AddExportAndNoCollapseToKey(keyNode pkg_rhino.Node) {
	_ = keyNode
	var newDocs *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(keyNode.GetJSDocInfo())
	_ = newDocs
	newDocs.RecordNoCollapse()
	newDocs.RecordExport()
	keyNode.SetJSDocInfo(newDocs.Build_0())
}

// Java static method: PolymerPassStaticUtils.extractProperties(com.google.javascript.rhino.Node descriptor, com.google.javascript.jscomp.PolymerClassDefinition.DefinitionType defType, com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node constructor)
func PolymerPassStaticUtils_ExtractProperties(descriptor pkg_rhino.Node, defType PolymerClassDefinition_DefinitionType, compiler AbstractCompiler, constructor pkg_rhino.Node) j2gort.Ext_com_google_common_collect_ImmutableList {
	_ = descriptor
	_ = defType
	_ = compiler
	_ = constructor
	var properties pkg_rhino.Node = descriptor
	_ = properties
	if defType == StaticGet_PolymerClassDefinition_DefinitionType_ObjectLiteral() {
		properties = NodeUtil_GetFirstPropMatchingKey(descriptor, "properties")
	}
	if properties == 0 {
		return j2gort.Ext_com_google_common_collect_ImmutableList_Of_0()
	}
	var constructorPropertyJsDoc j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = constructorPropertyJsDoc
	if constructor != 0 {
		PolymerPassStaticUtils_CollectConstructorPropertyJsDoc(constructor, constructorPropertyJsDoc)
	}
	var enclosingModule pkg_rhino.Node = NodeUtil_GetEnclosingModuleIfPresent(descriptor)
	_ = enclosingModule
	var members []*PolymerPass_MemberDefinition = []*PolymerPass_MemberDefinition{}
	_ = members
	{
		var keyNode pkg_rhino.Node = properties.GetFirstChild()
		_ = keyNode
		for ; keyNode != 0; keyNode = keyNode.GetNext() {
			var constructorJsDoc *pkg_rhino.JSDocInfo = j2gort.CastPtr[pkg_rhino.JSDocInfo](j2gort.MapGetStr(constructorPropertyJsDoc, keyNode.GetString()))
			_ = constructorJsDoc
			var propertiesConfigJsDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(keyNode)
			_ = propertiesConfigJsDoc
			var bestJsDoc *pkg_rhino.JSDocInfo = nil
			_ = bestJsDoc
			if constructorJsDoc != nil {
				bestJsDoc = constructorJsDoc
				if propertiesConfigJsDoc != nil {
					compiler.Report_1(JSError_Make_3(keyNode, StaticGet_PolymerPassErrors_POLYMER_MISPLACED_PROPERTY_JSDOC(), []string{}))
				}
			} else {
				bestJsDoc = propertiesConfigJsDoc
			}
			members = append(members, NewPolymerPass_MemberDefinition(bestJsDoc, keyNode, keyNode.GetFirstChild(), enclosingModule))
		}
	}
	return j2gort.ImmListFromTypedSlice(members)
}

// Java static method: PolymerPassStaticUtils.collectConstructorPropertyJsDoc(com.google.javascript.rhino.Node node, java.util.Map<java.lang.String,com.google.javascript.rhino.JSDocInfo> map)
func PolymerPassStaticUtils_CollectConstructorPropertyJsDoc(node pkg_rhino.Node, map_ j2gort.Ext_java_util_Map) {
	_ = node
	_ = map_
	j2gort.CheckNotNullNode(node)
	{
		var child pkg_rhino.Node = node.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if child.IsGetProp() && child.GetFirstChild().IsThis() {
				j2gort.MapPutStr(map_, child.GetString(), j2gort.PtrToAny(NodeUtil_GetBestJSDocInfo(child)))
			} else {
				PolymerPassStaticUtils_CollectConstructorPropertyJsDoc(child, map_)
			}
		}
	}
}

// Java static method: PolymerPassStaticUtils.getTypeFromProperty(com.google.javascript.jscomp.PolymerPass.MemberDefinition property, com.google.javascript.jscomp.AbstractCompiler compiler)
func PolymerPassStaticUtils_GetTypeFromProperty(property *PolymerPass_MemberDefinition, compiler AbstractCompiler) pkg_rhino.JSTypeExpression {
	_ = property
	_ = compiler
	if (property.info != nil) && property.info.HasType() {
		return property.info.GetType().Copy()
	}
	var typeValue pkg_rhino.Node = pkg_rhino.Node(0)
	_ = typeValue
	if property.value.IsObjectLit() {
		typeValue = NodeUtil_GetFirstPropMatchingKey(property.value, "type")
		if (typeValue == 0) || !(typeValue.IsName()) {
			compiler.Report_1(JSError_Make_3(property.name, StaticGet_PolymerPassErrors_POLYMER_INVALID_PROPERTY(), []string{}))
			return (pkg_rhino.JSTypeExpression{})
		}
	} else if property.value.IsName() {
		typeValue = property.value
	} else {
		typeValue = pkg_rhino.Node(0)
	}
	if typeValue == 0 {
		compiler.Report_1(JSError_Make_3(property.value, StaticGet_PolymerPassErrors_POLYMER_INVALID_PROPERTY(), []string{}))
		return (pkg_rhino.JSTypeExpression{})
	}
	var typeString string = typeValue.GetString()
	_ = typeString
	var typeNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = typeNode
	switch j2gort.CheckNotNullString(typeString) {
	case "Boolean", "String", "Number":
		typeNode = pkg_rhino.IR_String(j2gort.Ext_com_google_common_base_Ascii_ToLowerCase_1_string(typeString))
	case "Array", "Function", "Object", "Date":
		typeNode = pkg_rhino.NewNode_2(pkg_rhino.StaticGet_Token_BANG(), pkg_rhino.IR_String(typeString))
	default:
		compiler.Report_1(JSError_Make_3(property.value, StaticGet_PolymerPassErrors_POLYMER_INVALID_PROPERTY(), []string{}))
		typeNode = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_QMARK())
	}
	typeNode.SrcrefTree(typeValue)
	return pkg_rhino.NewJSTypeExpression(typeNode, "<PolymerPassStaticUtils.java>")
}

// Java static method: PolymerPassStaticUtils.getPolymerElementType(com.google.javascript.jscomp.PolymerClassDefinition cls)
func PolymerPassStaticUtils_GetPolymerElementType(cls *PolymerClassDefinition) string {
	_ = cls
	var nativeElementName string
	_ = nativeElementName
	if cls.nativeBaseElement == j2gort.NullString {
		nativeElementName = ""
	} else {
		nativeElementName = j2gort.Ext_com_google_common_base_CaseFormat_F_LOWER_HYPHEN.To(j2gort.Ext_com_google_common_base_CaseFormat_F_UPPER_CAMEL, cls.nativeBaseElement)
	}
	return j2gort.Ext_java_lang_String_Format_2("Polymer%sElement", []any{j2gort.BoxString(nativeElementName)})
}

// Java class:
type PolymerPassStaticUtils__Anon_70 struct {
	cap_0_compiler AbstractCompiler
	idHash_        atomic.Int32
}

func (this_ *PolymerPassStaticUtils__Anon_70) BindSelf_(s any) {
	_ = s
}
func (this_ *PolymerPassStaticUtils__Anon_70) IsJavaNil_() bool { return this_ == nil }
func (*PolymerPassStaticUtils__Anon_70) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.PolymerPassStaticUtils$1", ""
}
func (this_ *PolymerPassStaticUtils__Anon_70) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*PolymerPassStaticUtils__Anon_70) HasIdentityEquals_()                 {}
func (*PolymerPassStaticUtils__Anon_70) Implements_jscomp_NodeUtil_Visitor() {}

func (this_ *PolymerPassStaticUtils__Anon_70) initFields_PolymerPassStaticUtils__Anon_70() {
	_ = this_
}

func NewPolymerPassStaticUtils__Anon_70(cap_p_compiler AbstractCompiler) *PolymerPassStaticUtils__Anon_70 {
	this_ := &PolymerPassStaticUtils__Anon_70{}
	this_.cap_0_compiler = cap_p_compiler
	this_.initFields_PolymerPassStaticUtils__Anon_70()
	return this_
}

// Java method: .visit(com.google.javascript.rhino.Node n)
func (this_ *PolymerPassStaticUtils__Anon_70) Visit(n pkg_rhino.Node) {
	_ = n
	if !(n.IsGetProp()) {
		return
	}
	var dollarSign pkg_rhino.Node = n.GetFirstChild()
	_ = dollarSign
	if !(dollarSign.IsGetProp()) || !(j2gort.CheckNotNullString(dollarSign.GetString()) == "$") {
		return
	}
	if n.GetParent().IsCall() && n.GetParent().HasOneChild() {
		return
	}
	var key pkg_rhino.Node = pkg_rhino.IR_String(n.GetString()).ClonePropsFrom(n).Srcref(n)
	_ = key
	var getelem pkg_rhino.Node = pkg_rhino.IR_Getelem(dollarSign.Detach(), key).ClonePropsFrom(n).Srcref(n)
	_ = getelem
	n.ReplaceWith(getelem)
	this_.cap_0_compiler.ReportChangeToEnclosingScope(getelem)
}
func PolymerPassStaticUtils_AddObserverMethodsToProtect(behaviorObjLit pkg_rhino.Node) j2gort.Ext_java_util_Set {
	return ordSet_PolymerPassStaticUtils_AddObserverMethodsToProtect(behaviorObjLit).ToJavaSet()
}
func PolymerPassStaticUtils_AddPropertyMethodsToProtect(behaviorObjLit pkg_rhino.Node) j2gort.Ext_java_util_Set {
	return ordSet_PolymerPassStaticUtils_AddPropertyMethodsToProtect(behaviorObjLit).ToJavaSet()
}
