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
// Port of com/google/javascript/jscomp/ConcretizeStaticInheritanceForInlining.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.ConcretizeStaticInheritanceForInlining
type ConcretizeStaticInheritanceForInlining struct {
	duplicateClassNames j2gort.OrderedSet[string]                                                          `gson:"duplicateClassNames"`
	inherits            RuntimeJsLibManager_JsLibField                                                     `gson:"inherits"`
	compiler            AbstractCompiler                                                                   `gson:"compiler"`
	classByAlias        j2gort.OrderedMap[string, *ConcretizeStaticInheritanceForInlining_JavascriptClass] `gson:"classByAlias"`
	classLeafNames      j2gort.OrderedSet[string]                                                          `gson:"classLeafNames"`
	idHash_             atomic.Int32
}

func (this_ *ConcretizeStaticInheritanceForInlining) BindSelf_(s any) {
	_ = s
}
func (this_ *ConcretizeStaticInheritanceForInlining) IsJavaNil_() bool { return this_ == nil }
func (*ConcretizeStaticInheritanceForInlining) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConcretizeStaticInheritanceForInlining", "ConcretizeStaticInheritanceForInlining"
}
func (this_ *ConcretizeStaticInheritanceForInlining) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ConcretizeStaticInheritanceForInlining) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConcretizeStaticInheritanceForInlining", "com.google.javascript.jscomp.ConcretizeStaticInheritanceForInlining", func(v_ any) bool {
		return j2gort.InstanceOf[*ConcretizeStaticInheritanceForInlining](v_)
	})
}
func (*ConcretizeStaticInheritanceForInlining) Implements_jscomp_CompilerPass() {}

var static_ConcretizeStaticInheritanceForInlining_DUPLICATE_CLASS *DiagnosticType

func StaticGet_ConcretizeStaticInheritanceForInlining_DUPLICATE_CLASS() *DiagnosticType {
	clinit_ConcretizeStaticInheritanceForInlining()
	return static_ConcretizeStaticInheritanceForInlining_DUPLICATE_CLASS
}
func StaticSet_ConcretizeStaticInheritanceForInlining_DUPLICATE_CLASS(v *DiagnosticType) {
	clinit_ConcretizeStaticInheritanceForInlining()
	static_ConcretizeStaticInheritanceForInlining_DUPLICATE_CLASS = v
}

var static_ConcretizeStaticInheritanceForInlining_BANNED_PROP_NAMES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_ConcretizeStaticInheritanceForInlining_BANNED_PROP_NAMES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ConcretizeStaticInheritanceForInlining()
	return static_ConcretizeStaticInheritanceForInlining_BANNED_PROP_NAMES
}
func StaticSet_ConcretizeStaticInheritanceForInlining_BANNED_PROP_NAMES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_ConcretizeStaticInheritanceForInlining()
	static_ConcretizeStaticInheritanceForInlining_BANNED_PROP_NAMES = v
}

var clinitDone_ConcretizeStaticInheritanceForInlining uint32
var clinitState_ConcretizeStaticInheritanceForInlining j2gort.ClinitState

func clinit_ConcretizeStaticInheritanceForInlining() {
	if atomic.LoadUint32(&clinitDone_ConcretizeStaticInheritanceForInlining) != 0 {
		return
	}
	clinitSlow_ConcretizeStaticInheritanceForInlining()
}

func clinitSlow_ConcretizeStaticInheritanceForInlining() {
	if !j2gort.EnterClinit(&clinitState_ConcretizeStaticInheritanceForInlining, &clinitDone_ConcretizeStaticInheritanceForInlining) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ConcretizeStaticInheritanceForInlining, &clinitDone_ConcretizeStaticInheritanceForInlining)
	static_ConcretizeStaticInheritanceForInlining_DUPLICATE_CLASS = DiagnosticType_Error("DUPLICATE_CLASS", "Multiple classes cannot share the same name: {0}")
	static_ConcretizeStaticInheritanceForInlining_BANNED_PROP_NAMES = j2gort.ImmSetOfStr2("prototype", "getInstance")
}

func (this_ *ConcretizeStaticInheritanceForInlining) initFields_ConcretizeStaticInheritanceForInlining() {
	_ = this_
	this_.duplicateClassNames = j2gort.NewOrderedSet[string]()
	this_.classByAlias = j2gort.NewOrderedMap[string, *ConcretizeStaticInheritanceForInlining_JavascriptClass]()
	this_.classLeafNames = j2gort.NewOrderedSet[string]()
}

// Java constructor: ConcretizeStaticInheritanceForInlining(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewConcretizeStaticInheritanceForInlining(compiler AbstractCompiler) *ConcretizeStaticInheritanceForInlining {
	clinit_ConcretizeStaticInheritanceForInlining()
	this_ := &ConcretizeStaticInheritanceForInlining{}
	InitConcretizeStaticInheritanceForInlining(this_, compiler)
	return this_
}

func InitConcretizeStaticInheritanceForInlining(this_ *ConcretizeStaticInheritanceForInlining, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_ConcretizeStaticInheritanceForInlining()
	this_.compiler = compiler
	this_.inherits = compiler.GetRuntimeJsLibManager().GetJsLibField("$jscomp.inherits")
}

// Java method: ConcretizeStaticInheritanceForInlining.putClassAlias(java.lang.String name, com.google.javascript.jscomp.ConcretizeStaticInheritanceForInlining.JavascriptClass cls)
func (this_ *ConcretizeStaticInheritanceForInlining) priv_ConcretizeStaticInheritanceForInlining_putClassAlias_2123(name string, cls *ConcretizeStaticInheritanceForInlining_JavascriptClass) {
	_ = name
	_ = cls
	this_.classByAlias.Put(name, cls)
	var lastDot int32 = j2gort.JStrLastIndexOfCp(name, int32(uint16(46)))
	_ = lastDot
	this_.classLeafNames.Add(func() string {
		if lastDot < int32(0) {
			return name
		}
		return j2gort.JStrSubFrom(name, (lastDot + int32(1)))
	}())
}

// Java method: ConcretizeStaticInheritanceForInlining.getClassKeyIfPresent(com.google.javascript.rhino.Node n)
func (this_ *ConcretizeStaticInheritanceForInlining) priv_ConcretizeStaticInheritanceForInlining_getClassKeyIfPresent_2124(n pkg_rhino.Node) string {
	_ = n
	if this_.classByAlias.IsEmpty() {
		return j2gort.NullString
	}
	if n.IsName() {
		var name string = n.GetString()
		_ = name
		if !(j2gort.CheckNotNullString(name) == "") && this_.classByAlias.ContainsKey(name) {
			return name
		} else {
			return j2gort.NullString
		}
	}
	if n.IsGetProp() {
		var prop string = n.GetString()
		_ = prop
		var dot int32 = j2gort.JStrLastIndexOfCp(prop, int32(uint16(46)))
		_ = dot
		var leaf string
		_ = leaf
		if dot < int32(0) {
			leaf = prop
		} else {
			leaf = j2gort.JStrSubFrom(prop, (dot + int32(1)))
		}
		if !(this_.classLeafNames.Contains(leaf)) {
			return j2gort.NullString
		}
		var qname string = n.GetQualifiedName()
		_ = qname
		if (qname != j2gort.NullString) && this_.classByAlias.ContainsKey(qname) {
			return qname
		} else {
			return j2gort.NullString
		}
	}
	if n.IsThis() {
		if this_.classByAlias.ContainsKey("this") {
			return "this"
		} else {
			return j2gort.NullString
		}
	}
	if n.IsSuper() {
		if this_.classByAlias.ContainsKey("super") {
			return "super"
		} else {
			return j2gort.NullString
		}
	}
	return j2gort.NullString
}

// Java method: ConcretizeStaticInheritanceForInlining.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ConcretizeStaticInheritanceForInlining) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var findStaticMembers *ConcretizeStaticInheritanceForInlining_FindStaticMembers = NewConcretizeStaticInheritanceForInlining_FindStaticMembers(this_)
	_ = findStaticMembers
	NodeTraversal_Traverse(this_.compiler, root, findStaticMembers)
	this_.priv_ConcretizeStaticInheritanceForInlining_processInherits_2125(findStaticMembers)
}

// Java method: ConcretizeStaticInheritanceForInlining.processInherits(com.google.javascript.jscomp.ConcretizeStaticInheritanceForInlining.FindStaticMembers findStaticMembers)
func (this_ *ConcretizeStaticInheritanceForInlining) priv_ConcretizeStaticInheritanceForInlining_processInherits_2125(findStaticMembers *ConcretizeStaticInheritanceForInlining_FindStaticMembers) {
	_ = findStaticMembers
	for _, inheritsCall := range findStaticMembers.inheritsCalls {
		_ = inheritsCall
		var superclassNameNode pkg_rhino.Node = inheritsCall.GetLastChild()
		_ = superclassNameNode
		var superclassQname string = superclassNameNode.GetQualifiedName()
		_ = superclassQname
		var subclassNameNode pkg_rhino.Node = superclassNameNode.GetPrevious()
		_ = subclassNameNode
		var subclassQname string = subclassNameNode.GetQualifiedName()
		_ = subclassQname
		var superClass *ConcretizeStaticInheritanceForInlining_JavascriptClass = this_.classByAlias.Get(superclassQname)
		_ = superClass
		var subClass *ConcretizeStaticInheritanceForInlining_JavascriptClass = this_.classByAlias.Get(subclassQname)
		_ = subClass
		if this_.duplicateClassNames.Contains(superclassQname) {
			this_.compiler.Report_1(JSError_Make_3(inheritsCall, StaticGet_ConcretizeStaticInheritanceForInlining_DUPLICATE_CLASS(), []string{superclassQname}))
			return
		}
		if (superClass == nil) || (subClass == nil) {
			continue
		}
		this_.priv_ConcretizeStaticInheritanceForInlining_copyStaticMembers_2126(superClass, subClass, inheritsCall, findStaticMembers)
	}

}

// Java method: ConcretizeStaticInheritanceForInlining.copyStaticMembers(com.google.javascript.jscomp.ConcretizeStaticInheritanceForInlining.JavascriptClass superClass, com.google.javascript.jscomp.ConcretizeStaticInheritanceForInlining.JavascriptClass subClass, com.google.javascript.rhino.Node inheritsCall, com.google.javascript.jscomp.ConcretizeStaticInheritanceForInlining.FindStaticMembers findStaticMembers)
func (this_ *ConcretizeStaticInheritanceForInlining) priv_ConcretizeStaticInheritanceForInlining_copyStaticMembers_2126(superClass *ConcretizeStaticInheritanceForInlining_JavascriptClass, subClass *ConcretizeStaticInheritanceForInlining_JavascriptClass, inheritsCall pkg_rhino.Node, findStaticMembers *ConcretizeStaticInheritanceForInlining_FindStaticMembers) {
	_ = superClass
	_ = subClass
	_ = inheritsCall
	_ = findStaticMembers
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(superClass.staticMembers))
		for cur_2.HasNext() {
			staticMember := pkg_rhino.Node(cur_2.NextNode())
			_ = staticMember
			if !(staticMember.IsAssign()) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(staticMember))
			}
			var memberName string = staticMember.GetFirstChild().GetString()
			_ = memberName
			if superClass.definedProperties.Contains(memberName) {
				continue
			}
			if this_.priv_ConcretizeStaticInheritanceForInlining_isOverriden_2127(subClass, memberName) {
				continue
			}
			if findStaticMembers.IsBefore(inheritsCall, staticMember) {
				continue
			}
			var info *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(staticMember.GetJSDocInfo())
			_ = info
			var function pkg_rhino.Node = staticMember.GetLastChild()
			_ = function
			var sourceInfoNode pkg_rhino.Node = function
			_ = sourceInfoNode
			if function.IsFunction() {
				sourceInfoNode = function.GetFirstChild()
			}
			var subclassNameNode pkg_rhino.Node = inheritsCall.GetSecondChild()
			_ = subclassNameNode
			var superclassNameNode pkg_rhino.Node = subclassNameNode.GetNext()
			_ = superclassNameNode
			var assign pkg_rhino.Node = pkg_rhino.IR_Assign(pkg_rhino.IR_Getprop_2(subclassNameNode.CloneTree_0(), memberName), pkg_rhino.IR_Getprop_2(superclassNameNode.CloneTree_0(), memberName))
			_ = assign
			assign.SetJSDocInfo(info.Build_0())
			var exprResult pkg_rhino.Node = pkg_rhino.IR_ExprResult(assign)
			_ = exprResult
			exprResult.SrcrefTreeIfMissing(sourceInfoNode)
			var inheritsExpressionResult pkg_rhino.Node = inheritsCall.GetParent()
			_ = inheritsExpressionResult
			exprResult.InsertAfter(inheritsExpressionResult)
			this_.compiler.ReportChangeToEnclosingScope(inheritsExpressionResult)
			subClass.AddStaticMember(assign)
		}
	}
}

// Java method: ConcretizeStaticInheritanceForInlining.isOverriden(com.google.javascript.jscomp.ConcretizeStaticInheritanceForInlining.JavascriptClass subClass, java.lang.String memberName)
func (this_ *ConcretizeStaticInheritanceForInlining) priv_ConcretizeStaticInheritanceForInlining_isOverriden_2127(subClass *ConcretizeStaticInheritanceForInlining_JavascriptClass, memberName string) bool {
	_ = subClass
	_ = memberName
	if subClass.staticMemberNames.Contains(memberName) {
		return true
	}
	if subClass.definedProperties.Contains(memberName) {
		return true
	}
	return false
}

// Java method: ConcretizeStaticInheritanceForInlining.getReferencedClass(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConcretizeStaticInheritanceForInlining) priv_ConcretizeStaticInheritanceForInlining_getReferencedClass_2128(t *NodeTraversal, n pkg_rhino.Node) *ConcretizeStaticInheritanceForInlining_JavascriptClass {
	_ = t
	_ = n
	var className string = this_.priv_ConcretizeStaticInheritanceForInlining_getClassKeyIfPresent_2124(n)
	_ = className
	if className == j2gort.NullString {
		return nil
	}
	if n.IsName() {
		var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(className))
		_ = var_
		if (var_ != nil) && var_.IsLocal() {
			return nil
		}
	}
	return this_.classByAlias.Get(className)
}

// Java class: com.google.javascript.jscomp.ConcretizeStaticInheritanceForInlining.JavascriptClass
type ConcretizeStaticInheritanceForInlining_JavascriptClass struct {
	staticMembers     j2gort.Ext_java_util_Set  `gson:"staticMembers"`
	staticMemberNames j2gort.OrderedSet[string] `gson:"staticMemberNames"`
	definedProperties j2gort.OrderedSet[string] `gson:"definedProperties"`
	idHash_           atomic.Int32
}

func (this_ *ConcretizeStaticInheritanceForInlining_JavascriptClass) BindSelf_(s any) {
	_ = s
}
func (this_ *ConcretizeStaticInheritanceForInlining_JavascriptClass) IsJavaNil_() bool {
	return this_ == nil
}
func (*ConcretizeStaticInheritanceForInlining_JavascriptClass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConcretizeStaticInheritanceForInlining$JavascriptClass", "JavascriptClass"
}
func (this_ *ConcretizeStaticInheritanceForInlining_JavascriptClass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ConcretizeStaticInheritanceForInlining_JavascriptClass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConcretizeStaticInheritanceForInlining.JavascriptClass", "com.google.javascript.jscomp.ConcretizeStaticInheritanceForInlining$JavascriptClass", func(v_ any) bool {
		return j2gort.InstanceOf[*ConcretizeStaticInheritanceForInlining_JavascriptClass](v_)
	})
}

func (this_ *ConcretizeStaticInheritanceForInlining_JavascriptClass) initFields_ConcretizeStaticInheritanceForInlining_JavascriptClass() {
	_ = this_
	this_.staticMembers = j2gort.NewExt_java_util_LinkedHashSet()
	this_.staticMemberNames = j2gort.NewOrderedSet[string]()
	this_.definedProperties = j2gort.NewOrderedSet[string]()
}

// Java constructor: JavascriptClass()
func NewConcretizeStaticInheritanceForInlining_JavascriptClass() *ConcretizeStaticInheritanceForInlining_JavascriptClass {
	this_ := &ConcretizeStaticInheritanceForInlining_JavascriptClass{}
	InitConcretizeStaticInheritanceForInlining_JavascriptClass(this_)
	return this_
}

func InitConcretizeStaticInheritanceForInlining_JavascriptClass(this_ *ConcretizeStaticInheritanceForInlining_JavascriptClass) {
	this_.initFields_ConcretizeStaticInheritanceForInlining_JavascriptClass()
}

// Java method: JavascriptClass.addStaticMember(com.google.javascript.rhino.Node node)
func (this_ *ConcretizeStaticInheritanceForInlining_JavascriptClass) AddStaticMember(node pkg_rhino.Node) {
	_ = node
	j2gort.SetAddNode(this_.staticMembers, node)
	this_.staticMemberNames.Add(node.GetFirstChild().GetString())
}

// Java class: com.google.javascript.jscomp.ConcretizeStaticInheritanceForInlining.FindStaticMembers
type ConcretizeStaticInheritanceForInlining_FindStaticMembers struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_        *ConcretizeStaticInheritanceForInlining
	inheritsCalls []pkg_rhino.Node                         `gson:"inheritsCalls"`
	nodeOrder     j2gort.OrderedMap[pkg_rhino.Node, int32] `gson:"nodeOrder"`
}

func (this_ *ConcretizeStaticInheritanceForInlining_FindStaticMembers) BindSelf_(s any) {
	_ = s
}
func (this_ *ConcretizeStaticInheritanceForInlining_FindStaticMembers) IsJavaNil_() bool {
	return this_ == nil
}
func (*ConcretizeStaticInheritanceForInlining_FindStaticMembers) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ConcretizeStaticInheritanceForInlining$FindStaticMembers", "FindStaticMembers"
}
func (this_ *ConcretizeStaticInheritanceForInlining_FindStaticMembers) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ConcretizeStaticInheritanceForInlining_FindStaticMembers) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ConcretizeStaticInheritanceForInlining.FindStaticMembers", "com.google.javascript.jscomp.ConcretizeStaticInheritanceForInlining$FindStaticMembers", func(v_ any) bool {
		return j2gort.InstanceOf[*ConcretizeStaticInheritanceForInlining_FindStaticMembers](v_)
	})
}
func (this_ *ConcretizeStaticInheritanceForInlining_FindStaticMembers) SetOuter_(o *ConcretizeStaticInheritanceForInlining) {
	this_.outer_ = o
}
func (this_ *ConcretizeStaticInheritanceForInlining_FindStaticMembers) GetOuter_() *ConcretizeStaticInheritanceForInlining {
	return this_.outer_
}
func (*ConcretizeStaticInheritanceForInlining_FindStaticMembers) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*ConcretizeStaticInheritanceForInlining_FindStaticMembers) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *ConcretizeStaticInheritanceForInlining_FindStaticMembers) initFields_ConcretizeStaticInheritanceForInlining_FindStaticMembers() {
	_ = this_
	this_.inheritsCalls = []pkg_rhino.Node{}
	this_.nodeOrder = j2gort.NewOrderedMap[pkg_rhino.Node, int32](

	// Java constructor: FindStaticMembers()
	)
}

func NewConcretizeStaticInheritanceForInlining_FindStaticMembers(outer_ *ConcretizeStaticInheritanceForInlining) *ConcretizeStaticInheritanceForInlining_FindStaticMembers {
	this_ := &ConcretizeStaticInheritanceForInlining_FindStaticMembers{}
	this_.outer_ = outer_
	InitConcretizeStaticInheritanceForInlining_FindStaticMembers(this_)
	return this_
}

func InitConcretizeStaticInheritanceForInlining_FindStaticMembers(this_ *ConcretizeStaticInheritanceForInlining_FindStaticMembers) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ConcretizeStaticInheritanceForInlining_FindStaticMembers()
}

// Java method: FindStaticMembers.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ConcretizeStaticInheritanceForInlining_FindStaticMembers) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(28) /* Token.CALL */ :
		if this_.outer_.inherits.Matches(n.GetFirstChild()) {
			this_.inheritsCalls = append(this_.inheritsCalls, n)
			this_.nodeOrder.Put(n, this_.nodeOrder.Size())
		}
		if NodeUtil_IsObjectDefinePropertiesDefinition(n) {
			this_.priv_FindStaticMembers_visitDefinedPropertiesCall_2130(t, n)
		}
	case int32(97) /* Token.CONST */, int32(111) /* Token.LET */, int32(84) /* Token.VAR */ :
		this_.priv_FindStaticMembers_visitVariableDeclaration_2134(n)
	case int32(51) /* Token.ASSIGN */ :
		this_.priv_FindStaticMembers_visitAssign_2133(t, n)
	case int32(73) /* Token.FUNCTION */ :
		this_.priv_FindStaticMembers_visitFunctionClassDef_2131(n)
	default:
	}
}

// Java method: FindStaticMembers.visitDefinedPropertiesCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node definePropertiesCall)
func (this_ *ConcretizeStaticInheritanceForInlining_FindStaticMembers) priv_FindStaticMembers_visitDefinedPropertiesCall_2130(t *NodeTraversal, definePropertiesCall pkg_rhino.Node) {
	_ = t
	_ = definePropertiesCall
	var object pkg_rhino.Node = definePropertiesCall.GetSecondChild()
	_ = object
	var c *ConcretizeStaticInheritanceForInlining_JavascriptClass = this_.outer_.priv_ConcretizeStaticInheritanceForInlining_getReferencedClass_2128(t, object)
	_ = c
	if c != nil {
		{
			cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(NodeUtil_GetObjectDefinedPropertiesKeys(definePropertiesCall)))
			for cur_3.HasNext() {
				prop := pkg_rhino.Node(cur_3.NextNode())
				_ = prop
				c.definedProperties.Add(prop.GetString())
			}
		}
	}
}

// Java method: FindStaticMembers.visitFunctionClassDef(com.google.javascript.rhino.Node n)
func (this_ *ConcretizeStaticInheritanceForInlining_FindStaticMembers) priv_FindStaticMembers_visitFunctionClassDef_2131(n pkg_rhino.Node) {
	_ = n
	var classInfo *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
	_ = classInfo
	if (classInfo != nil) && classInfo.IsConstructor() {
		var name string = NodeUtil_GetName(n)
		_ = name
		if name == j2gort.NullString {
			return
		}
		if this_.outer_.classByAlias.ContainsKey(name) {
			this_.outer_.duplicateClassNames.Add(name)
		} else {
			this_.outer_.priv_ConcretizeStaticInheritanceForInlining_putClassAlias_2123(name, NewConcretizeStaticInheritanceForInlining_JavascriptClass())
		}
	}
}

// Java method: FindStaticMembers.setAlias(java.lang.String original, java.lang.String alias)
func (this_ *ConcretizeStaticInheritanceForInlining_FindStaticMembers) priv_FindStaticMembers_setAlias_2132(original string, alias string) {
	_ = original
	_ = alias
	j2gort.CheckNotNullString(original)
	j2gort.CheckNotNullString(alias)
	if !(this_.outer_.classByAlias.ContainsKey(original)) {
		j2gort.CheckArgument(false)
	}
	this_.outer_.priv_ConcretizeStaticInheritanceForInlining_putClassAlias_2123(alias, this_.outer_.classByAlias.Get(original))
}

// Java method: FindStaticMembers.visitAssign(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ConcretizeStaticInheritanceForInlining_FindStaticMembers) priv_FindStaticMembers_visitAssign_2133(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if this_.outer_.classByAlias.IsEmpty() {
		return
	}
	var existingClassQname string = this_.outer_.priv_ConcretizeStaticInheritanceForInlining_getClassKeyIfPresent_2124(n.GetLastChild())
	_ = existingClassQname
	if existingClassQname != j2gort.NullString {
		var alias string = n.GetFirstChild().GetQualifiedName()
		_ = alias
		if alias != j2gort.NullString {
			this_.priv_FindStaticMembers_setAlias_2132(existingClassQname, alias)
		}
	} else if n.GetFirstChild().IsGetProp() {
		var getProp pkg_rhino.Node = n.GetFirstChild()
		_ = getProp
		var classNode pkg_rhino.Node = getProp.GetFirstChild()
		_ = classNode
		var refClass *ConcretizeStaticInheritanceForInlining_JavascriptClass = this_.outer_.priv_ConcretizeStaticInheritanceForInlining_getReferencedClass_2128(t, classNode)
		_ = refClass
		if (refClass != nil) && !(j2gort.SetContainsStr(StaticGet_ConcretizeStaticInheritanceForInlining_BANNED_PROP_NAMES(), getProp.GetString())) {
			refClass.AddStaticMember(n)
			this_.nodeOrder.Put(n, this_.nodeOrder.Size())
		}
	}
}

// Java method: FindStaticMembers.visitVariableDeclaration(com.google.javascript.rhino.Node n)
func (this_ *ConcretizeStaticInheritanceForInlining_FindStaticMembers) priv_FindStaticMembers_visitVariableDeclaration_2134(n pkg_rhino.Node) {
	_ = n
	if this_.outer_.classByAlias.IsEmpty() {
		return
	}
	var child pkg_rhino.Node = n.GetFirstChild()
	_ = child
	if !(child.HasChildren()) {
		return
	}
	var maybeOriginalName string = this_.outer_.priv_ConcretizeStaticInheritanceForInlining_getClassKeyIfPresent_2124(child.GetFirstChild())
	_ = maybeOriginalName
	if maybeOriginalName != j2gort.NullString {
		var maybeAlias string = child.GetQualifiedName()
		_ = maybeAlias
		if maybeAlias != j2gort.NullString {
			this_.priv_FindStaticMembers_setAlias_2132(maybeOriginalName, maybeAlias)
		}
	}
}

// Java method: FindStaticMembers.isBefore(com.google.javascript.rhino.Node earlier, com.google.javascript.rhino.Node later)
func (this_ *ConcretizeStaticInheritanceForInlining_FindStaticMembers) IsBefore(earlier pkg_rhino.Node, later pkg_rhino.Node) bool {
	_ = earlier
	_ = later
	earlierPosition, earlierPositionOk_ := this_.nodeOrder.GetOk(earlier)
	_ = earlierPositionOk_
	_ = earlierPosition
	laterPosition, laterPositionOk_ := this_.nodeOrder.GetOk(later)
	_ = laterPositionOk_
	_ = laterPosition
	return (((earlierPositionOk_) && (laterPositionOk_)) && (func() int32 {
		if !earlierPositionOk_ {
			j2gort.ThrowNPE()
		}
		return earlierPosition
	}() < func() int32 {
		if !laterPositionOk_ {
			j2gort.ThrowNPE()
		}
		return laterPosition
	}()))
}
func (*ConcretizeStaticInheritanceForInlining_FindStaticMembers) AlwaysShouldTraverse_()      {}
func (*ConcretizeStaticInheritanceForInlining_FindStaticMembers) CanReleaseSyntacticScopes_() {}
