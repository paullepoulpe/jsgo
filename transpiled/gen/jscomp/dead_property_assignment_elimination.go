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
// Port of com/google/javascript/jscomp/DeadPropertyAssignmentElimination.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

var lambdaSingleton_dead_property_assignment_elimination_1 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_dead_property_assignment_elimination_1 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool {
		return NodeUtil_IsFunctionOrNonStaticClassFieldInitializer(pkg_rhino.Node(rawNode_))
	})
}

// Java class: com.google.javascript.jscomp.DeadPropertyAssignmentElimination
type DeadPropertyAssignmentElimination struct {
	compiler                                        AbstractCompiler `gson:"compiler"`
	idHash_                                         atomic.Int32
	slab_DeadPropertyAssignmentElimination_Property j2gort.PooledSlab[DeadPropertyAssignmentElimination_Property]
}

func (this_ *DeadPropertyAssignmentElimination) BindSelf_(s any) {
	_ = s
}
func (this_ *DeadPropertyAssignmentElimination) IsJavaNil_() bool { return this_ == nil }
func (*DeadPropertyAssignmentElimination) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DeadPropertyAssignmentElimination", "DeadPropertyAssignmentElimination"
}
func (this_ *DeadPropertyAssignmentElimination) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DeadPropertyAssignmentElimination) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DeadPropertyAssignmentElimination", "com.google.javascript.jscomp.DeadPropertyAssignmentElimination", func(v_ any) bool {
		return j2gort.InstanceOf[*DeadPropertyAssignmentElimination](v_)
	})
}
func (*DeadPropertyAssignmentElimination) Implements_jscomp_CompilerPass() {}

func (this_ *DeadPropertyAssignmentElimination) initFields_DeadPropertyAssignmentElimination() {
	_ = this_
}

// Java constructor: DeadPropertyAssignmentElimination(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewDeadPropertyAssignmentElimination(compiler AbstractCompiler) *DeadPropertyAssignmentElimination {
	this_ := &DeadPropertyAssignmentElimination{}
	InitDeadPropertyAssignmentElimination(this_, compiler)
	return this_
}

func InitDeadPropertyAssignmentElimination(this_ *DeadPropertyAssignmentElimination, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_DeadPropertyAssignmentElimination()
	this_.compiler = compiler
}

// Java method: DeadPropertyAssignmentElimination.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *DeadPropertyAssignmentElimination) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	defer this_.slab_DeadPropertyAssignmentElimination_Property.Release()

	if (this_.compiler.GetExternProperties() == nil) || (this_.compiler.GetAccessorSummary() == nil) {
		return
	}
	var skiplistedPropNames j2gort.Ext_java_util_Set = j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_Sets_Union(j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableSet](this_.compiler.GetAccessorSummary().GetAccessors().KeySet_0()), j2gort.AsIface[j2gort.Ext_java_util_Set](this_.compiler.GetExternProperties())))
	_ = skiplistedPropNames
	NodeTraversal_Traverse(this_.compiler, this_.compiler.GetJsRoot(), allocFunctionVisitorWithSlab_(&this_.slab_DeadPropertyAssignmentElimination_Property, skiplistedPropNames))
}

// Java class: com.google.javascript.jscomp.DeadPropertyAssignmentElimination.FunctionVisitor
type DeadPropertyAssignmentElimination_FunctionVisitor struct {
	NodeTraversal_AbstractChangedScopeCallbackBase
	skiplistedPropNames j2gort.Ext_java_util_Set `gson:"skiplistedPropNames"`
	propSlab_           *j2gort.PooledSlab[DeadPropertyAssignmentElimination_Property]
}

func (this_ *DeadPropertyAssignmentElimination_FunctionVisitor) BindSelf_(s any) {
	_ = s
	this_.NodeTraversal_AbstractChangedScopeCallbackBase.BindSelf_(s)
}
func (this_ *DeadPropertyAssignmentElimination_FunctionVisitor) IsJavaNil_() bool {
	return this_ == nil
}
func (*DeadPropertyAssignmentElimination_FunctionVisitor) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DeadPropertyAssignmentElimination$FunctionVisitor", "FunctionVisitor"
}
func (this_ *DeadPropertyAssignmentElimination_FunctionVisitor) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractChangedScopeCallbackBase.IdentityHash_()
}
func (*DeadPropertyAssignmentElimination_FunctionVisitor) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DeadPropertyAssignmentElimination.FunctionVisitor", "com.google.javascript.jscomp.DeadPropertyAssignmentElimination$FunctionVisitor", func(v_ any) bool {
		return j2gort.InstanceOf[*DeadPropertyAssignmentElimination_FunctionVisitor](v_)
	})
}
func (*DeadPropertyAssignmentElimination_FunctionVisitor) Implements_jscomp_NodeTraversal_AbstractChangedScopeCallback() {
}
func (*DeadPropertyAssignmentElimination_FunctionVisitor) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {
}
func (*DeadPropertyAssignmentElimination_FunctionVisitor) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *DeadPropertyAssignmentElimination_FunctionVisitor) initFields_DeadPropertyAssignmentElimination_FunctionVisitor() {
	_ = this_
}

// Java constructor: FunctionVisitor(java.util.Set<java.lang.String> skiplistedPropNames)
func NewDeadPropertyAssignmentElimination_FunctionVisitor(skiplistedPropNames j2gort.Ext_java_util_Set) *DeadPropertyAssignmentElimination_FunctionVisitor {
	this_ := &DeadPropertyAssignmentElimination_FunctionVisitor{}
	this_.BindSelf_(this_)
	InitDeadPropertyAssignmentElimination_FunctionVisitor(this_, skiplistedPropNames)
	return this_
}
func allocFunctionVisitorWithSlab_(slab_ *j2gort.PooledSlab[DeadPropertyAssignmentElimination_Property], skiplistedPropNames j2gort.Ext_java_util_Set) *DeadPropertyAssignmentElimination_FunctionVisitor {
	v_ := NewDeadPropertyAssignmentElimination_FunctionVisitor(skiplistedPropNames)
	v_.propSlab_ = slab_
	return v_
}

func InitDeadPropertyAssignmentElimination_FunctionVisitor(this_ *DeadPropertyAssignmentElimination_FunctionVisitor, skiplistedPropNames j2gort.Ext_java_util_Set) {
	_ = skiplistedPropNames
	InitNodeTraversal_AbstractChangedScopeCallback(&this_.NodeTraversal_AbstractChangedScopeCallbackBase)
	this_.initFields_DeadPropertyAssignmentElimination_FunctionVisitor()
	this_.skiplistedPropNames = skiplistedPropNames
}

// Java method: FunctionVisitor.enterChangedScopeRoot(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node root)
func (this_ *DeadPropertyAssignmentElimination_FunctionVisitor) EnterChangedScopeRoot(compiler AbstractCompiler, root pkg_rhino.Node) {
	_ = compiler
	_ = root
	if !(root.IsFunction()) {
		return
	}
	var body pkg_rhino.Node = NodeUtil_GetFunctionBody(root)
	_ = body
	if !(body.HasChildren()) || NodeUtil_Has(body, lambdaSingleton_dead_property_assignment_elimination_1, j2gort.Ext_com_google_common_base_Predicates_AlwaysTrue()) {
		return
	}
	var val_1 DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal
	InitDeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal(&val_1, this_.skiplistedPropNames)
	val_1.propSlab_ = this_.propSlab_
	if this_.propSlab_ != nil {
		defer this_.propSlab_.Reset()
	}

	var traversal *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal = &val_1
	_ = traversal
	NodeTraversal_Traverse(compiler, body, traversal)
	{
		cur_2 := traversal.propertyMap.IterValues()
		for cur_2.HasNext() {

			property := cur_2.
				Next()
			_ = property
			if property.writes.Size() <= int32(1) {
				continue
			}
			var iter j2gort.Ext_com_google_common_collect_PeekingIterator = j2gort.Ext_com_google_common_collect_Iterators_PeekingIterator_1_Ext_java_util_Iterator(j2gort.DequeIterator(&property.writes))
			_ = iter
			for iter.HasNext() {
				var propertyWrite *DeadPropertyAssignmentElimination_PropertyWrite = j2gort.CastPtr[DeadPropertyAssignmentElimination_PropertyWrite](iter.Next())
				_ = propertyWrite
				if iter.HasNext() && propertyWrite.IsSafeToRemove(j2gort.CastPtr[DeadPropertyAssignmentElimination_PropertyWrite](iter.Peek())) {
					var lhs pkg_rhino.Node = propertyWrite.assignedAt
					_ = lhs
					var rhs pkg_rhino.Node = lhs.GetNext()
					_ = rhs
					var assignNode pkg_rhino.Node = lhs.GetParent()
					_ = assignNode
					if assignNode.IsAssign() {
						rhs.Detach()
						assignNode.ReplaceWith(rhs)
						compiler.ReportChangeToEnclosingScope(rhs)
					} else {
						if !(NodeUtil_IsAssignmentOp(assignNode)) {
							j2gort.CheckState(false)
						}
						var opType pkg_rhino.Token = NodeUtil_GetOpFromAssignmentOp(assignNode)
						_ = opType
						assignNode.SetToken(opType)
						compiler.ReportChangeToEnclosingScope(assignNode)
					}
				}
			}
		}
	}
}

// Java class: com.google.javascript.jscomp.DeadPropertyAssignmentElimination.Property
type DeadPropertyAssignmentElimination_Property struct {
	name     string                                                         `gson:"name"`
	writes   j2gort.Deque[*DeadPropertyAssignmentElimination_PropertyWrite] `gson:"writes"`
	children j2gort.OrderedSet[*DeadPropertyAssignmentElimination_Property] `gson:"children"`
	idHash_  atomic.Int32
}

func (this_ *DeadPropertyAssignmentElimination_Property) BindSelf_(s any) {
	_ = s
}
func (this_ *DeadPropertyAssignmentElimination_Property) IsJavaNil_() bool { return this_ == nil }
func (*DeadPropertyAssignmentElimination_Property) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DeadPropertyAssignmentElimination$Property", "Property"
}
func (this_ *DeadPropertyAssignmentElimination_Property) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DeadPropertyAssignmentElimination_Property) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DeadPropertyAssignmentElimination.Property", "com.google.javascript.jscomp.DeadPropertyAssignmentElimination$Property", func(v_ any) bool {
		return j2gort.InstanceOf[*DeadPropertyAssignmentElimination_Property](v_)
	})
}

func (this_ *DeadPropertyAssignmentElimination_Property) initFields_DeadPropertyAssignmentElimination_Property() {
	_ = this_
	this_.name = j2gort.NullString
	this_.writes = j2gort.Deque[*DeadPropertyAssignmentElimination_PropertyWrite]{}
	this_.children = j2gort.NewOrderedSet[*DeadPropertyAssignmentElimination_Property](

	// Java constructor: Property(java.lang.String name)
	)
}

func NewDeadPropertyAssignmentElimination_Property(name string) *DeadPropertyAssignmentElimination_Property {
	this_ := &DeadPropertyAssignmentElimination_Property{}
	InitDeadPropertyAssignmentElimination_Property(this_, name)
	return this_
}
func allocOptPooled_DeadPropertyAssignmentElimination_Property(slab_ *j2gort.PooledSlab[DeadPropertyAssignmentElimination_Property], name string) *DeadPropertyAssignmentElimination_Property {
	var this_ *DeadPropertyAssignmentElimination_Property
	if slab_ != nil {
		this_ = slab_.Alloc()
	} else {
		this_ = &DeadPropertyAssignmentElimination_Property{}
	}
	InitDeadPropertyAssignmentElimination_Property(this_, name)
	return this_
}

func InitDeadPropertyAssignmentElimination_Property(this_ *DeadPropertyAssignmentElimination_Property, name string) {
	_ = name
	this_.initFields_DeadPropertyAssignmentElimination_Property()
	this_.name = name
}

// Java method: Property.markLastWriteRead()
func (this_ *DeadPropertyAssignmentElimination_Property) MarkLastWriteRead() {
	if !(this_.writes.IsEmpty()) {
		this_.writes.GetLast().MarkRead()
	}
}

// Java method: Property.markChildrenRead()
func (this_ *DeadPropertyAssignmentElimination_Property) MarkChildrenRead() {
	if this_.children.IsEmpty() {
		return
	}
	var propertiesSet j2gort.OrderedSet[*DeadPropertyAssignmentElimination_Property] = j2gort.OrderedSetFromSet(this_.children)
	_ = propertiesSet
	var propertyQueue j2gort.Deque[*DeadPropertyAssignmentElimination_Property] = j2gort.DequeFromOrderedSet(propertiesSet)
	_ = propertyQueue
	propertiesSet.Add(this_)
	for !(propertyQueue.IsEmpty()) {
		var childProperty *DeadPropertyAssignmentElimination_Property = propertyQueue.RemoveFirst()
		_ = childProperty
		childProperty.MarkLastWriteRead()
		{
			cur_4 := childProperty.children.Iter()
			for cur_4.HasNext() {

				grandchildProperty := cur_4.
					Next()
				_ = grandchildProperty
				if propertiesSet.Add(grandchildProperty) {
					propertyQueue.AddLast(grandchildProperty)
				}
			}
		}
	}
}

// Java method: Property.addWrite(com.google.javascript.rhino.Node lhs)
func (this_ *DeadPropertyAssignmentElimination_Property) AddWrite(lhs pkg_rhino.Node) {
	_ = lhs
	if !(lhs.IsQualifiedName()) {
		j2gort.CheckArgument(false)
	}
	this_.writes.AddLast(NewDeadPropertyAssignmentElimination_PropertyWrite(lhs))
}

// Java method: Property.toString()
func (this_ *DeadPropertyAssignmentElimination_Property) String() string {
	return ("Property " + j2gort.StrValString(this_.name))
}

// Java class: com.google.javascript.jscomp.DeadPropertyAssignmentElimination.PropertyWrite
type DeadPropertyAssignmentElimination_PropertyWrite struct {
	qualifiedName string         `gson:"qualifiedName"`
	assignedAt    pkg_rhino.Node `gson:"assignedAt"`
	idHash_       atomic.Int32
	isRead        bool `gson:"isRead"`
}

func (this_ *DeadPropertyAssignmentElimination_PropertyWrite) BindSelf_(s any) {
	_ = s
}
func (this_ *DeadPropertyAssignmentElimination_PropertyWrite) IsJavaNil_() bool { return this_ == nil }
func (*DeadPropertyAssignmentElimination_PropertyWrite) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DeadPropertyAssignmentElimination$PropertyWrite", "PropertyWrite"
}
func (this_ *DeadPropertyAssignmentElimination_PropertyWrite) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DeadPropertyAssignmentElimination_PropertyWrite) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DeadPropertyAssignmentElimination.PropertyWrite", "com.google.javascript.jscomp.DeadPropertyAssignmentElimination$PropertyWrite", func(v_ any) bool {
		return j2gort.InstanceOf[*DeadPropertyAssignmentElimination_PropertyWrite](v_)
	})
}

func (this_ *DeadPropertyAssignmentElimination_PropertyWrite) initFields_DeadPropertyAssignmentElimination_PropertyWrite() {
	_ = this_

	this_.qualifiedName = j2gort.NullString
}

// Java constructor: PropertyWrite(com.google.javascript.rhino.Node assignedAt)
func NewDeadPropertyAssignmentElimination_PropertyWrite(assignedAt pkg_rhino.Node) *DeadPropertyAssignmentElimination_PropertyWrite {
	this_ := &DeadPropertyAssignmentElimination_PropertyWrite{}
	InitDeadPropertyAssignmentElimination_PropertyWrite(this_, assignedAt)
	return this_
}

func InitDeadPropertyAssignmentElimination_PropertyWrite(this_ *DeadPropertyAssignmentElimination_PropertyWrite, assignedAt pkg_rhino.Node) {
	_ = assignedAt
	this_.initFields_DeadPropertyAssignmentElimination_PropertyWrite()
	if !(assignedAt.IsQualifiedName()) {
		j2gort.CheckArgument(false)
	}
	this_.assignedAt = assignedAt
	this_.qualifiedName = assignedAt.GetQualifiedName()
}

// Java method: PropertyWrite.isSafeToRemove(com.google.javascript.jscomp.DeadPropertyAssignmentElimination.@org.jspecify.annotations.Nullable PropertyWrite nextWrite)
func (this_ *DeadPropertyAssignmentElimination_PropertyWrite) IsSafeToRemove(nextWrite *DeadPropertyAssignmentElimination_PropertyWrite) bool {
	_ = nextWrite
	return ((!(this_.isRead) && (nextWrite != nil)) && (this_.qualifiedName == nextWrite.qualifiedName))
}

// Java method: PropertyWrite.markRead()
func (this_ *DeadPropertyAssignmentElimination_PropertyWrite) MarkRead() {
	this_.isRead = true
}

// Java class: com.google.javascript.jscomp.DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal
type DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal struct {
	propertyMap         j2gort.OrderedMap[string, *DeadPropertyAssignmentElimination_Property] `gson:"propertyMap"`
	skiplistedPropNames j2gort.Ext_java_util_Set                                               `gson:"skiplistedPropNames"`
	idHash_             atomic.Int32
	propSlab_           *j2gort.PooledSlab[DeadPropertyAssignmentElimination_Property]
}

func (this_ *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) BindSelf_(s any) {
	_ = s
}
func (this_ *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) IsJavaNil_() bool {
	return this_ == nil
}
func (*DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.DeadPropertyAssignmentElimination$FindCandidateAssignmentTraversal", "FindCandidateAssignmentTraversal"
}
func (this_ *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.DeadPropertyAssignmentElimination.FindCandidateAssignmentTraversal", "com.google.javascript.jscomp.DeadPropertyAssignmentElimination$FindCandidateAssignmentTraversal", func(v_ any) bool {
		return j2gort.InstanceOf[*DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal](v_)
	})
}
func (*DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) initFields_DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal() {
	_ = this_
	this_.propertyMap = j2gort.NewOrderedMap[string, *DeadPropertyAssignmentElimination_Property](

	// Java constructor: FindCandidateAssignmentTraversal(java.util.Set<java.lang.String> skiplistedPropNames)
	)
}

func NewDeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal(skiplistedPropNames j2gort.Ext_java_util_Set) *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal {
	this_ := &DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal{}
	InitDeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal(this_, skiplistedPropNames)
	return this_
}

func InitDeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal(this_ *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal, skiplistedPropNames j2gort.Ext_java_util_Set) {
	_ = skiplistedPropNames
	this_.initFields_DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal()
	this_.skiplistedPropNames = skiplistedPropNames
}

// Java method: FindCandidateAssignmentTraversal.getOrCreateProperty(com.google.javascript.rhino.Node propNode)
func (this_ *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) priv_FindCandidateAssignmentTraversal_getOrCreateProperty_2349(propNode pkg_rhino.Node) *DeadPropertyAssignmentElimination_Property {
	_ = propNode
	if !(propNode.IsQualifiedName()) {
		return nil
	}
	var propName string
	_ = propName
	if propNode.IsGetProp() {
		propName = propNode.GetString()
	} else {
		propName = propNode.GetQualifiedName()
	}
	var property *DeadPropertyAssignmentElimination_Property = this_.propertyMap.Get(propName)
	_ = property
	if property == nil {
		property = allocOptPooled_DeadPropertyAssignmentElimination_Property(this_.propSlab_, propName)
		this_.propertyMap.Put(propName, property)
	}
	if propNode.IsGetProp() {
		var parentProperty *DeadPropertyAssignmentElimination_Property = this_.priv_FindCandidateAssignmentTraversal_getOrCreateProperty_2349(propNode.GetFirstChild())
		_ = parentProperty
		if parentProperty != nil {
			parentProperty.children.Add(property)
		}
	}
	return property
}

// Java method: FindCandidateAssignmentTraversal.shouldTraverse(com.google.javascript.jscomp.NodeTraversal nodeTraversal, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) ShouldTraverse_3(nodeTraversal *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = nodeTraversal
	_ = n
	_ = parent
	return this_.priv_FindCandidateAssignmentTraversal_visitNode_2352(n, parent)
}

// Java method: FindCandidateAssignmentTraversal.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if NodeUtil_IsAssignmentOp(n) {
		this_.priv_FindCandidateAssignmentTraversal_visitAssignmentLhs_2351(n.GetFirstChild())
	}
	if (NodeUtil_IsInvocation(n) || n.IsYield()) || n.IsAwait() {
		this_.priv_FindCandidateAssignmentTraversal_markAllPropsRead_2353()
	}
	if n.IsBlock() {
		this_.priv_FindCandidateAssignmentTraversal_visitBlock_2350(n)
	}
}

// Java method: FindCandidateAssignmentTraversal.visitBlock(com.google.javascript.rhino.Node blockNode)
func (this_ *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) priv_FindCandidateAssignmentTraversal_visitBlock_2350(blockNode pkg_rhino.Node) {
	_ = blockNode
	if !(blockNode.IsBlock()) {
		j2gort.CheckArgument(false)
	}
	if blockNode.HasChildren() {
		this_.priv_FindCandidateAssignmentTraversal_markAllPropsRead_2353()
	}
}

// Java static method: FindCandidateAssignmentTraversal.isConditionalExpression(com.google.javascript.rhino.Node n)
func DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal_IsConditionalExpression(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(69) /* Token.AND */, int32(68) /* Token.OR */, int32(67) /* Token.HOOK */, int32(70) /* Token.COALESCE */, int32(31) /* Token.OPTCHAIN_CALL */, int32(30) /* Token.OPTCHAIN_GETELEM */, int32(29) /* Token.OPTCHAIN_GETPROP */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: FindCandidateAssignmentTraversal.visitAssignmentLhs(com.google.javascript.rhino.Node lhs)
func (this_ *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) priv_FindCandidateAssignmentTraversal_visitAssignmentLhs_2351(lhs pkg_rhino.Node) {
	_ = lhs
	var property *DeadPropertyAssignmentElimination_Property = this_.priv_FindCandidateAssignmentTraversal_getOrCreateProperty_2349(lhs)
	_ = property
	if property == nil {
		return
	}
	if !(lhs.IsGetProp()) {
		property.MarkLastWriteRead()
		property.MarkChildrenRead()
		return
	}
	var assignNode pkg_rhino.Node = lhs.GetParent()
	_ = assignNode
	if !(assignNode.IsAssign()) {
		property.MarkLastWriteRead()
	}
	property.MarkChildrenRead()
	property.AddWrite(lhs)
	var child pkg_rhino.Node = lhs.GetFirstChild()
	_ = child
	for child != 0 {
		var childProperty *DeadPropertyAssignmentElimination_Property = this_.priv_FindCandidateAssignmentTraversal_getOrCreateProperty_2349(child)
		_ = childProperty
		if childProperty == nil {
			break
		}
		childProperty.MarkLastWriteRead()
		child = child.GetFirstChild()
	}
}

// Java method: FindCandidateAssignmentTraversal.visitNode(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) priv_FindCandidateAssignmentTraversal_visitNode_2352(n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */ :
		if n.IsGetProp() && j2gort.SetContainsStr(this_.skiplistedPropNames, n.GetString()) {
			this_.priv_FindCandidateAssignmentTraversal_markAllPropsRead_2353()
			return true
		}
		if NodeUtil_IsAssignmentOp(parent) && (parent.GetFirstChild() == n) {
			return true
		}
		var property *DeadPropertyAssignmentElimination_Property = this_.priv_FindCandidateAssignmentTraversal_getOrCreateProperty_2349(n)
		_ = property
		if property != nil {
			property.MarkLastWriteRead()
			if !(parent.IsGetProp()) {
				property.MarkChildrenRead()
			}
		}
		return true
	case int32(37) /* Token.THIS */, int32(32) /* Token.NAME */ :
		var nameProp *DeadPropertyAssignmentElimination_Property = j2gort.CheckNotNullPtr(this_.priv_FindCandidateAssignmentTraversal_getOrCreateProperty_2349(n))
		_ = nameProp
		nameProp.MarkLastWriteRead()
		if !(parent.IsGetProp()) {
			nameProp.MarkChildrenRead()
		}
		return true
	case int32(43) /* Token.THROW */, int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(127) /* Token.ITER_SPREAD */, int32(102) /* Token.ARRAY_PATTERN */, int32(75) /* Token.SWITCH */ :
		this_.priv_FindCandidateAssignmentTraversal_markAllPropsRead_2353()
		return false
	case int32(103) /* Token.OBJECT_PATTERN */ :
		{
			var child pkg_rhino.Node = n.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				if child.IsStringKey() {
					var propName string = child.GetString()
					_ = propName
					var property_2 *DeadPropertyAssignmentElimination_Property = this_.propertyMap.Get(propName)
					_ = property_2
					if property_2 != nil {
						property_2.MarkLastWriteRead()
						property_2.MarkChildrenRead()
					}
				} else {
					this_.priv_FindCandidateAssignmentTraversal_markAllPropsRead_2353()
					break
				}
			}
		}
		return true
	case int32(90) /* Token.BLOCK */ :
		this_.priv_FindCandidateAssignmentTraversal_visitBlock_2350(n)
		return true
	default:
		if DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal_IsConditionalExpression(n) {
			this_.priv_FindCandidateAssignmentTraversal_markAllPropsRead_2353()
			return false
		}
		return true
	}
	panic("unreachable switch expression")
}

// Java method: FindCandidateAssignmentTraversal.markAllPropsRead()
func (this_ *DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) priv_FindCandidateAssignmentTraversal_markAllPropsRead_2353() {
	{
		cur_6 := this_.propertyMap.IterValues()
		for cur_6.HasNext() {

			property := cur_6.
				Next()
			_ = property
			if property.writes.IsEmpty() {
				continue
			}
			property.MarkLastWriteRead()
		}
	}
}
func (*DeadPropertyAssignmentElimination_FunctionVisitor) NoOpVisit_()                 {}
func (*DeadPropertyAssignmentElimination_FunctionVisitor) CanReleaseSyntacticScopes_() {}
func (*DeadPropertyAssignmentElimination_FindCandidateAssignmentTraversal) CanReleaseSyntacticScopes_() {
}
