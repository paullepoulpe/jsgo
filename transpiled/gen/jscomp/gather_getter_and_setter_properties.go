/*
 * Copyright 2018 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/GatherGetterAndSetterProperties.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.GatherGetterAndSetterProperties
type GatherGetterAndSetterProperties struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_  atomic.Int32
}

func (this_ *GatherGetterAndSetterProperties) BindSelf_(s any) {
	_ = s
}
func (this_ *GatherGetterAndSetterProperties) IsJavaNil_() bool { return this_ == nil }
func (*GatherGetterAndSetterProperties) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GatherGetterAndSetterProperties", "GatherGetterAndSetterProperties"
}
func (this_ *GatherGetterAndSetterProperties) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GatherGetterAndSetterProperties) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GatherGetterAndSetterProperties", "com.google.javascript.jscomp.GatherGetterAndSetterProperties", func(v_ any) bool {
		return j2gort.InstanceOf[*GatherGetterAndSetterProperties](v_)
	})
}
func (*GatherGetterAndSetterProperties) Implements_jscomp_CompilerPass() {}

func (this_ *GatherGetterAndSetterProperties) initFields_GatherGetterAndSetterProperties() {
	_ = this_
}

// Java constructor: GatherGetterAndSetterProperties(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewGatherGetterAndSetterProperties(compiler AbstractCompiler) *GatherGetterAndSetterProperties {
	this_ := &GatherGetterAndSetterProperties{}
	InitGatherGetterAndSetterProperties(this_, compiler)
	return this_
}

func InitGatherGetterAndSetterProperties(this_ *GatherGetterAndSetterProperties, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_GatherGetterAndSetterProperties()
	this_.compiler = compiler
}

// Java method: GatherGetterAndSetterProperties.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *GatherGetterAndSetterProperties) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	GatherGetterAndSetterProperties_Update(this_.compiler, externs, root)
}

// Java static method: GatherGetterAndSetterProperties.update(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func GatherGetterAndSetterProperties_Update(compiler AbstractCompiler, externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = compiler
	_ = externs
	_ = root
	if !(externs.GetParent() == root.GetParent()) {
		j2gort.CheckState(false)
	}
	var accessorSummary *AccessorSummary = nil
	_ = accessorSummary
	if compiler.GetOptions().GetAssumePropertiesAreStaticallyAnalyzable() {
		accessorSummary = AccessorSummary_Create(j2gort.AsIface[j2gort.Ext_java_util_Map](GatherGetterAndSetterProperties_Gather(compiler, externs.GetParent())))
	} else {
		accessorSummary = AccessorSummary_CreateAssumingAlwaysGetterAndSetter()
	}
	compiler.SetAccessorSummary(accessorSummary)
}

// Java static method: GatherGetterAndSetterProperties.gather(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node root)
func GatherGetterAndSetterProperties_Gather(compiler AbstractCompiler, root pkg_rhino.Node) j2gort.Ext_com_google_common_collect_ImmutableMap {
	_ = compiler
	_ = root
	var gatherCallback *GatherGetterAndSetterProperties_GatherCallback = NewGatherGetterAndSetterProperties_GatherCallback()
	_ = gatherCallback
	NodeTraversal_Traverse(compiler, root, gatherCallback)
	return gatherCallback.properties.ToImmutableMap()
}

// Java class: com.google.javascript.jscomp.GatherGetterAndSetterProperties.GatherCallback
type GatherGetterAndSetterProperties_GatherCallback struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	properties j2gort.OrderedMap[string, *AccessorSummary_PropertyAccessKind] `gson:"properties"`
}

func (this_ *GatherGetterAndSetterProperties_GatherCallback) BindSelf_(s any) {
	_ = s
}
func (this_ *GatherGetterAndSetterProperties_GatherCallback) IsJavaNil_() bool { return this_ == nil }
func (*GatherGetterAndSetterProperties_GatherCallback) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GatherGetterAndSetterProperties$GatherCallback", "GatherCallback"
}
func (this_ *GatherGetterAndSetterProperties_GatherCallback) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*GatherGetterAndSetterProperties_GatherCallback) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GatherGetterAndSetterProperties.GatherCallback", "com.google.javascript.jscomp.GatherGetterAndSetterProperties$GatherCallback", func(v_ any) bool {
		return j2gort.InstanceOf[*GatherGetterAndSetterProperties_GatherCallback](v_)
	})
}
func (*GatherGetterAndSetterProperties_GatherCallback) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*GatherGetterAndSetterProperties_GatherCallback) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *GatherGetterAndSetterProperties_GatherCallback) initFields_GatherGetterAndSetterProperties_GatherCallback() {
	_ = this_
	this_.properties = j2gort.NewOrderedMap[string, *AccessorSummary_PropertyAccessKind](

	// Java constructor: GatherCallback()
	)
}

func NewGatherGetterAndSetterProperties_GatherCallback() *GatherGetterAndSetterProperties_GatherCallback {
	this_ := &GatherGetterAndSetterProperties_GatherCallback{}
	InitGatherGetterAndSetterProperties_GatherCallback(this_)
	return this_
}

func InitGatherGetterAndSetterProperties_GatherCallback(this_ *GatherGetterAndSetterProperties_GatherCallback) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_GatherGetterAndSetterProperties_GatherCallback()
}

// Java method: GatherCallback.record(java.lang.String property, com.google.javascript.jscomp.AccessorSummary.PropertyAccessKind kind)
func (this_ *GatherGetterAndSetterProperties_GatherCallback) priv_GatherCallback_record_3009(property string, kind *AccessorSummary_PropertyAccessKind) {
	_ = property
	_ = kind
	var existing *AccessorSummary_PropertyAccessKind = this_.properties.Get(property)
	_ = existing
	this_.properties.Put(property, func() *AccessorSummary_PropertyAccessKind {
		if existing == nil {
			return kind
		}
		return existing.UnionWith(kind)
	}())
}

// Java method: GatherCallback.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *GatherGetterAndSetterProperties_GatherCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(95) /* Token.GETTER_DEF */ :
		this_.priv_GatherCallback_recordGetterDef_3010(n)
	case int32(96) /* Token.SETTER_DEF */ :
		this_.priv_GatherCallback_recordSetterDef_3011(n)
	case int32(28) /* Token.CALL */ :
		if NodeUtil_IsObjectDefinePropertyDefinition(n) {
			this_.priv_GatherCallback_visitDefineProperty_3013(n)
		} else if NodeUtil_IsObjectDefinePropertiesDefinition(n) || NodeUtil_IsObjectCreateWithPropertiesDefinition(n) {
			this_.priv_GatherCallback_visitDefineProperties_3014(n)
		}
	default:
	}
}

// Java method: GatherCallback.recordGetterDef(com.google.javascript.rhino.Node getterDef)
func (this_ *GatherGetterAndSetterProperties_GatherCallback) priv_GatherCallback_recordGetterDef_3010(getterDef pkg_rhino.Node) {
	_ = getterDef
	if !(getterDef.IsGetterDef()) {
		j2gort.CheckState(false)
	}
	var name string = getterDef.GetString()
	_ = name
	this_.priv_GatherCallback_record_3009(name, StaticGet_AccessorSummary_PropertyAccessKind_GETTER_ONLY())
}

// Java method: GatherCallback.recordSetterDef(com.google.javascript.rhino.Node setterDef)
func (this_ *GatherGetterAndSetterProperties_GatherCallback) priv_GatherCallback_recordSetterDef_3011(setterDef pkg_rhino.Node) {
	_ = setterDef
	if !(setterDef.IsSetterDef()) {
		j2gort.CheckState(false)
	}
	var name string = setterDef.GetString()
	_ = name
	this_.priv_GatherCallback_record_3009(name, StaticGet_AccessorSummary_PropertyAccessKind_SETTER_ONLY())
}

// Java method: GatherCallback.visitDescriptor(java.lang.String propertyName, com.google.javascript.rhino.Node descriptor)
func (this_ *GatherGetterAndSetterProperties_GatherCallback) priv_GatherCallback_visitDescriptor_3012(propertyName string, descriptor pkg_rhino.Node) {
	_ = propertyName
	_ = descriptor
	{
		var key pkg_rhino.Node = descriptor.GetFirstChild()
		_ = key
		for ; key != 0; key = key.GetNext() {
			if key.IsStringKey() || key.IsMemberFunctionDef() {
				{
					var iifeCond_2 bool
					{

						sr_1 := "get"
						se_2 := key.GetString()
						iifeCond_2 = j2gort.CheckNotNullString(sr_1) == se_2
					}
					if iifeCond_2 {
						this_.priv_GatherCallback_record_3009(propertyName, StaticGet_AccessorSummary_PropertyAccessKind_GETTER_ONLY())
					} else {
						var iifeCond_1 bool
						{

							sr_3 := "set"
							se_4 := key.GetString()
							iifeCond_1 = j2gort.CheckNotNullString(sr_3) == se_4
						}
						if iifeCond_1 {
							this_.priv_GatherCallback_record_3009(propertyName, StaticGet_AccessorSummary_PropertyAccessKind_SETTER_ONLY())
						}
					}
				}
			}
		}
	}
}

// Java method: GatherCallback.visitDefineProperty(com.google.javascript.rhino.Node definePropertyCall)
func (this_ *GatherGetterAndSetterProperties_GatherCallback) priv_GatherCallback_visitDefineProperty_3013(definePropertyCall pkg_rhino.Node) {
	_ = definePropertyCall
	var propertyNameNode pkg_rhino.Node = definePropertyCall.GetChildAtIndex(int32(2))
	_ = propertyNameNode
	var descriptor pkg_rhino.Node = definePropertyCall.GetChildAtIndex(int32(3))
	_ = descriptor
	if !(propertyNameNode.IsStringLit()) || !(descriptor.IsObjectLit()) {
		return
	}
	var propertyName string = propertyNameNode.GetString()
	_ = propertyName
	this_.priv_GatherCallback_visitDescriptor_3012(propertyName, descriptor)
}

// Java method: GatherCallback.visitDefineProperties(com.google.javascript.rhino.Node definePropertiesCall)
func (this_ *GatherGetterAndSetterProperties_GatherCallback) priv_GatherCallback_visitDefineProperties_3014(definePropertiesCall pkg_rhino.Node) {
	_ = definePropertiesCall
	var props pkg_rhino.Node = definePropertiesCall.GetChildAtIndex(int32(2))
	_ = props
	if !(props.IsObjectLit()) {
		return
	}
	{
		var prop pkg_rhino.Node = props.GetFirstChild()
		_ = prop
		for ; prop != 0; prop = prop.GetNext() {
			if (prop.IsStringKey() && prop.HasOneChild()) && prop.GetFirstChild().IsObjectLit() {
				var propertyName string = prop.GetString()
				_ = propertyName
				var descriptor pkg_rhino.Node = prop.GetFirstChild()
				_ = descriptor
				this_.priv_GatherCallback_visitDescriptor_3012(propertyName, descriptor)
			}
		}
	}
}
func (*GatherGetterAndSetterProperties_GatherCallback) AlwaysShouldTraverse_()      {}
func (*GatherGetterAndSetterProperties_GatherCallback) CanReleaseSyntacticScopes_() {}
