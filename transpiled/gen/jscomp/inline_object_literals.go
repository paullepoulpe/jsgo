/*
 * Copyright 2011 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/InlineObjectLiterals.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.InlineObjectLiterals
type InlineObjectLiterals struct {
	compiler                                   AbstractCompiler                       `gson:"compiler"`
	safeNameIdSupplier                         j2gort.Ext_java_util_function_Supplier `gson:"safeNameIdSupplier"`
	idHash_                                    atomic.Int32
	slab_InlineObjectLiterals_InliningBehavior j2gort.PooledSlab[InlineObjectLiterals_InliningBehavior]
}

func (this_ *InlineObjectLiterals) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineObjectLiterals) IsJavaNil_() bool { return this_ == nil }
func (*InlineObjectLiterals) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineObjectLiterals", "InlineObjectLiterals"
}
func (this_ *InlineObjectLiterals) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineObjectLiterals) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineObjectLiterals", "com.google.javascript.jscomp.InlineObjectLiterals", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineObjectLiterals](v_)
	})
}
func (*InlineObjectLiterals) Implements_jscomp_CompilerPass() {}

func StaticGet_InlineObjectLiterals_VAR_PREFIX() string            { return "JSCompiler_object_inline_" }
func StaticGet_InlineObjectLiterals_STRING_KEY_IDENTIFIER() string { return "string_key" }
func (this_ *InlineObjectLiterals) initFields_InlineObjectLiterals() {
	_ = this_
}

// Java constructor: InlineObjectLiterals(com.google.javascript.jscomp.AbstractCompiler compiler, java.util.function.Supplier<java.lang.String> safeNameIdSupplier)
func NewInlineObjectLiterals(compiler AbstractCompiler, safeNameIdSupplier j2gort.Ext_java_util_function_Supplier) *InlineObjectLiterals {
	this_ := &InlineObjectLiterals{}
	InitInlineObjectLiterals(this_, compiler, safeNameIdSupplier)
	return this_
}

func InitInlineObjectLiterals(this_ *InlineObjectLiterals, compiler AbstractCompiler, safeNameIdSupplier j2gort.Ext_java_util_function_Supplier) {
	_ = compiler
	_ = safeNameIdSupplier
	this_.initFields_InlineObjectLiterals()
	this_.compiler = compiler
	this_.safeNameIdSupplier = safeNameIdSupplier
}

// Java method: InlineObjectLiterals.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *InlineObjectLiterals) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	defer this_.ReleaseSlabs_()
	clinit_ReferenceCollector()
	var val_1 ReferenceCollector
	InitReferenceCollector_3(&val_1, this_.compiler, ToIface_ReferenceCollector_Behavior(allocIolPooled_InlineObjectLiterals_InliningBehavior(this_)), NewSyntacticScopeCreator_1(this_.compiler))
	var callback *ReferenceCollector = &val_1
	_ = callback
	callback.ownsScopeCreator_ = true
	defer callback.ReleaseSlabs_()

	callback.Process_2(externs, root)
}

// Java class: com.google.javascript.jscomp.InlineObjectLiterals.InliningBehavior
type InlineObjectLiterals_InliningBehavior struct {
	outer_          *InlineObjectLiterals
	staleVars       j2gort.CustomOrderedSet[*Var] `gson:"staleVars"`
	validProperties j2gort.Ext_java_util_HashSet  `gson:"validProperties"`
	idHash_         atomic.Int32
}

func (this_ *InlineObjectLiterals_InliningBehavior) BindSelf_(s any) {
	_ = s
}
func (this_ *InlineObjectLiterals_InliningBehavior) IsJavaNil_() bool { return this_ == nil }
func (*InlineObjectLiterals_InliningBehavior) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.InlineObjectLiterals$InliningBehavior", "InliningBehavior"
}
func (this_ *InlineObjectLiterals_InliningBehavior) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*InlineObjectLiterals_InliningBehavior) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.InlineObjectLiterals.InliningBehavior", "com.google.javascript.jscomp.InlineObjectLiterals$InliningBehavior", func(v_ any) bool {
		return j2gort.InstanceOf[*InlineObjectLiterals_InliningBehavior](v_)
	})
}
func (this_ *InlineObjectLiterals_InliningBehavior) SetOuter_(o *InlineObjectLiterals) {
	this_.outer_ = o
}
func (this_ *InlineObjectLiterals_InliningBehavior) GetOuter_() *InlineObjectLiterals {
	return this_.outer_
}
func (*InlineObjectLiterals_InliningBehavior) Implements_jscomp_ReferenceCollector_Behavior() {}

func (this_ *InlineObjectLiterals_InliningBehavior) initFields_InlineObjectLiterals_InliningBehavior() {
	_ = this_
	this_.staleVars = j2gort.NewCustomOrderedSet[*Var]()
	this_.validProperties = j2gort.NewExt_java_util_HashSet()
}

// Java constructor: InliningBehavior()
func NewInlineObjectLiterals_InliningBehavior(outer_ *InlineObjectLiterals) *InlineObjectLiterals_InliningBehavior {
	this_ := &InlineObjectLiterals_InliningBehavior{}
	this_.outer_ = outer_
	InitInlineObjectLiterals_InliningBehavior(this_)
	return this_
}
func allocIolPooled_InlineObjectLiterals_InliningBehavior(outer_ *InlineObjectLiterals) *InlineObjectLiterals_InliningBehavior {
	if outer_ == nil {
		return NewInlineObjectLiterals_InliningBehavior(outer_)
	}
	this_ := outer_.slab_InlineObjectLiterals_InliningBehavior.Alloc()
	this_.outer_ = outer_
	InitInlineObjectLiterals_InliningBehavior(this_)
	return this_
}

func InitInlineObjectLiterals_InliningBehavior(this_ *InlineObjectLiterals_InliningBehavior) {
	this_.initFields_InlineObjectLiterals_InliningBehavior()
}

// Java method: InliningBehavior.afterExitScope(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.ReferenceMap referenceMap)
func (this_ *InlineObjectLiterals_InliningBehavior) AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap) {
	_ = t
	_ = referenceMap
	var scope *Scope = t.GetScope()
	_ = scope
	if scope.IsGlobal() {
		return
	}
	{
		cur_2 := (scope).iterVars_()
		for cur_2.HasNext() {

			v := cur_2.Next()

			_ = v
			if v.IsGlobal() || v.IsExtern() {
				continue
			}
			var referenceInfo *ReferenceCollection = j2gort.CastPtr[ReferenceCollection](referenceMap.GetReferences(j2gort.PtrToAny(v)))
			_ = referenceInfo
			if (referenceInfo == nil) || this_.priv_InliningBehavior_isVarInlineForbidden_3334(v) {
				continue
			}
			if this_.priv_InliningBehavior_isInlinableObject_3335(referenceInfo) {
				this_.staleVars.Add(v)
				var init *Reference = referenceInfo.GetInitializingReference()
				_ = init
				this_.priv_InliningBehavior_splitObject_3340(v, init, referenceInfo)
			}
		}
	}
}

// Java method: InliningBehavior.recordStaleVarReferencesInTree(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.Scope scope)
func (this_ *InlineObjectLiterals_InliningBehavior) priv_InliningBehavior_recordStaleVarReferencesInTree_3333(node pkg_rhino.Node, scope *Scope) {
	_ = node
	_ = scope
	if node.IsName() {
		var var_ *Var = scope.LookupVar(node.GetString())
		_ = var_
		if var_ != nil {
			this_.staleVars.Add(var_)
		}
	}
	if !(node.IsFunction()) {
		{
			var c pkg_rhino.Node = node.GetFirstChild()
			_ = c
			for ; c != 0; c = c.GetNext() {
				this_.priv_InliningBehavior_recordStaleVarReferencesInTree_3333(c, scope)
			}
		}
	}
}

// Java method: InliningBehavior.isVarInlineForbidden(com.google.javascript.jscomp.Var var)
func (this_ *InlineObjectLiterals_InliningBehavior) priv_InliningBehavior_isVarInlineForbidden_3334(var_ *Var) bool {
	_ = var_
	return ((((var_.IsGlobal() || var_.IsExtern()) || (!(this_.staleVars.IsEmpty()) && this_.staleVars.Contains(var_))) || this_.outer_.compiler.GetCodingConvention().IsExported_2(var_.GetName(), true)) || this_.outer_.compiler.GetCodingConvention().IsPropertyRenameFunction(var_.GetNameNode()))
}

// Java method: InliningBehavior.isInlinableObject(com.google.javascript.jscomp.ReferenceCollection refs)
func (this_ *InlineObjectLiterals_InliningBehavior) priv_InliningBehavior_isInlinableObject_3335(refs *ReferenceCollection) bool {
	_ = refs
	var ret bool = false
	_ = ret
	var hasAnyValidProps bool = false
	_ = hasAnyValidProps
	var validProps j2gort.Ext_java_util_HashSet = this_.validProperties
	_ = validProps
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = refs.Size()
		_ = size
		for ; i < size; i = i + 1 {
			var ref *Reference = refs.Get(i)
			_ = ref
			var name pkg_rhino.Node = ref.GetNode()
			_ = name
			var parent pkg_rhino.Node = ref.GetParent()
			_ = parent
			var grandparent pkg_rhino.Node = ref.GetGrandparent()
			_ = grandparent
			if NodeUtil_IsNormalOrOptChainGetProp(parent) {
				if !(parent.GetFirstChild() == name) {
					j2gort.CheckState(false)
				}
				if NodeUtil_IsInvocationTarget(parent) {
					return false
				}
				if grandparent.IsDelProp() {
					return false
				}
				var propName string = parent.GetString()
				_ = propName
				if !(hasAnyValidProps) || !(j2gort.SetContainsStr(validProps, propName)) {
					if NodeUtil_IsNameDeclOrSimpleAssignLhs(parent, grandparent) {
						if !(hasAnyValidProps) {
							validProps.Clear()
							hasAnyValidProps = true
						}
						j2gort.SetAddStr(validProps, propName)
					} else {
						return false
					}
				}
				continue
			}
			if !(this_.priv_InliningBehavior_isVarOrAssignExprLhs_3336(name)) {
				return false
			}
			if NodeUtil_IsNameDeclaration(parent) && NodeUtil_IsAnyFor(grandparent) {
				return false
			}
			var val pkg_rhino.Node = ref.GetAssignedValue()
			_ = val
			if val == 0 {
				continue
			}
			if !(val.IsObjectLit()) {
				return false
			}
			{
				var child pkg_rhino.Node = val.GetFirstChild()
				_ = child
				for ; child != 0; child = child.GetNext() {
					switch child.TokenSwitchOrd_() {
					case int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(129) /* Token.COMPUTED_PROP */, int32(128) /* Token.OBJECT_SPREAD */ :
						return false
					case int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(100) /* Token.STRING_KEY */ :
					default:
						panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected child of OBJECTLIT: "+j2gort.StrValString(child.ToStringTree()))))
					}
					if !(hasAnyValidProps) {
						validProps.Clear()
						hasAnyValidProps = true
					}
					j2gort.SetAddStr(validProps, child.GetString())
					var childVal pkg_rhino.Node = child.GetFirstChild()
					_ = childVal
					{
						var j int32 = int32(0)
						_ = j
						for ; j < size; j = j + 1 {
							var t *Reference = refs.Get(j)
							_ = t
							var refNode pkg_rhino.Node = t.GetParent()
							_ = refNode
							for !(NodeUtil_IsStatementBlock(refNode)) {
								if refNode == childVal {
									return false
								}
								refNode = refNode.GetParent()
							}
						}
					}
				}
			}
			ret = true
		}
	}
	return ret
}

// Java method: InliningBehavior.isVarOrAssignExprLhs(com.google.javascript.rhino.Node n)
func (this_ *InlineObjectLiterals_InliningBehavior) priv_InliningBehavior_isVarOrAssignExprLhs_3336(n pkg_rhino.Node) bool {
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	return (NodeUtil_IsNameDeclaration(parent) || ((parent.IsAssign() && (parent.GetFirstChild() == n)) && parent.GetParent().IsExprResult()))
}

// Java method: InliningBehavior.computeVarList(com.google.javascript.jscomp.ReferenceCollection referenceInfo)
func (this_ *InlineObjectLiterals_InliningBehavior) priv_InliningBehavior_computeVarList_3337(referenceInfo *ReferenceCollection) j2gort.Ext_java_util_Map {
	_ = referenceInfo
	var varmap j2gort.Ext_java_util_Map = j2gort.NewExt_java_util_LinkedHashMap()
	_ = varmap
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = referenceInfo.Size()
		_ = size
		for ; i < size; i = i + 1 {
			var ref *Reference = referenceInfo.Get(i)
			_ = ref
			if ref.IsLvalue() || ref.IsInitializingDeclaration() {
				var val pkg_rhino.Node = ref.GetAssignedValue()
				_ = val
				if val != 0 {
					if !(val.IsObjectLit()) {
						j2gort.CheckStateMsg(false, j2gort.BoxNode(val))
					}
					{
						var child pkg_rhino.Node = val.GetFirstChild()
						_ = child
						for ; child != 0; child = child.GetNext() {
							var varname string = child.GetString()
							_ = varname
							if j2gort.MapContainsKeyStr(varmap, varname) {
								continue
							}
							var var_ string = varname
							_ = var_
							if !(pkg_rhino.TokenStream_IsJSIdentifier(varname)) {
								var_ = "string_key"
							}
							var_ = ("JSCompiler_object_inline_" + j2gort.StrValString(var_) + "_" + j2gort.StrValString(j2gort.CastString(this_.outer_.safeNameIdSupplier.Get())))
							j2gort.MapPutStr(varmap, varname, j2gort.BoxString(var_))
						}
					}
				}
			} else if NodeUtil_IsNameDeclaration(ref.GetParent()) {
			} else {
				var getprop pkg_rhino.Node = ref.GetParent()
				_ = getprop
				if !(NodeUtil_IsNormalOrOptChainGetProp(getprop)) {
					j2gort.CheckStateMsg(false, j2gort.BoxNode(getprop))
				}
				var varname_2 string = getprop.GetString()
				_ = varname_2
				if j2gort.MapContainsKeyStr(varmap, varname_2) {
					continue
				}
				var var__2 string = ("JSCompiler_object_inline_" + j2gort.StrValString(varname_2) + "_" + j2gort.StrValString(j2gort.CastString(this_.outer_.safeNameIdSupplier.Get())))
				_ = var__2
				j2gort.MapPutStr(varmap, varname_2, j2gort.BoxString(var__2))
			}
		}
	}
	return varmap
}

// Java method: InliningBehavior.fillInitialValues(com.google.javascript.jscomp.Reference init, java.util.Map<java.lang.String,com.google.javascript.rhino.Node> initvals)
func (this_ *InlineObjectLiterals_InliningBehavior) priv_InliningBehavior_fillInitialValues_3338(init *Reference, initvals *j2gort.OrderedMap[string, pkg_rhino.Node]) {
	_ = init
	_ = initvals
	var object pkg_rhino.Node = init.GetAssignedValue()
	_ = object
	if !(object.IsObjectLit()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(object))
	}
	{
		var key pkg_rhino.Node = object.GetFirstChild()
		_ = key
		for ; key != 0; key = key.GetNext() {
			initvals.Put(key.GetString(), key.RemoveFirstChild())
		}
	}
}

// Java method: InliningBehavior.replaceAssignmentExpression(com.google.javascript.jscomp.Var v, com.google.javascript.jscomp.Reference ref, java.util.Map<java.lang.String,java.lang.String> varmap)
func (this_ *InlineObjectLiterals_InliningBehavior) priv_InliningBehavior_replaceAssignmentExpression_3339(v *Var, ref *Reference, varmap j2gort.Ext_java_util_Map) {
	_ = v
	_ = ref
	_ = varmap
	var nodes []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = nodes
	var val pkg_rhino.Node = ref.GetAssignedValue()
	_ = val
	this_.priv_InliningBehavior_recordStaleVarReferencesInTree_3333(val, j2gort.CastPtr[Scope](v.GetScope()))
	if !(val.IsObjectLit()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(val))
	}
	var all j2gort.LocalOrderedSet[string] = j2gort.LocalOrderedSetFromStrCol(j2gort.AsIface[j2gort.Ext_java_util_Set](varmap.KeySet_0()))
	_ = all
	{
		var key pkg_rhino.Node = val.GetFirstChild()
		_ = key
		for ; key != 0; key = key.GetNext() {
			var var_ string = key.GetString()
			_ = var_
			var value pkg_rhino.Node = key.RemoveFirstChild()
			_ = value
			nodes = append(nodes, pkg_rhino.IR_Assign(pkg_rhino.IR_Name(j2gort.CastString(j2gort.MapGetStr(varmap, var_))), value))
			all.Remove(var_)
		}
	}
	{
		cur_4 := all.Iter()
		for cur_4.HasNext() {
			var__2 := cur_4.
				Next()
			_ = var__2
			nodes = append(nodes, pkg_rhino.IR_Assign(pkg_rhino.IR_Name(j2gort.CastString(j2gort.MapGetStr(varmap, var__2))), NodeUtil_NewUndefinedNode(pkg_rhino.Node(0))))
		}
	}
	var replacement pkg_rhino.Node = pkg_rhino.Node(0)
	_ = replacement
	if len(nodes) == 0 {
		replacement = pkg_rhino.IR_TrueNode()
	} else {
		nodes = append(nodes, pkg_rhino.IR_TrueNode())
		nodes = j2gort.SliceReverseInPlace(nodes)
		replacement = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_COMMA())
		var cur pkg_rhino.Node = replacement
		_ = cur
		var i int32 = int32(0)
		_ = i
		{
			i = int32(0)
			for ; i < (int32(len(nodes)) - int32(2)); i = i + 1 {
				cur.AddChildToFront(nodes[i])
				var t pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_COMMA())
				_ = t
				cur.AddChildToFront(t)
				cur = t
			}
		}
		cur.AddChildToFront(nodes[i])
		cur.AddChildToFront(nodes[(i + int32(1))])
	}
	var replace pkg_rhino.Node = ref.GetParent()
	_ = replace
	replacement.SrcrefTreeIfMissing(replace)
	if NodeUtil_IsNameDeclaration(replace) {
		replace.ReplaceWith(NodeUtil_NewExpr(replacement))
	} else {
		replace.ReplaceWith(replacement)
	}
}

// Java method: InliningBehavior.splitObject(com.google.javascript.jscomp.Var v, com.google.javascript.jscomp.Reference init, com.google.javascript.jscomp.ReferenceCollection referenceInfo)
func (this_ *InlineObjectLiterals_InliningBehavior) priv_InliningBehavior_splitObject_3340(v *Var, init *Reference, referenceInfo *ReferenceCollection) {
	_ = v
	_ = init
	_ = referenceInfo
	var varmap j2gort.Ext_java_util_Map = this_.priv_InliningBehavior_computeVarList_3337(referenceInfo)
	_ = varmap
	var initvals j2gort.OrderedMap[string, pkg_rhino.Node] = j2gort.NewOrderedMap[string, pkg_rhino.Node]()
	_ = initvals
	var vnode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = vnode
	var isBlockScoped bool = (v.IsLet() || v.IsConst())
	_ = isBlockScoped
	var defined bool = (referenceInfo.IsWellDefined() && NodeUtil_IsNameDeclaration(init.GetParent()))
	_ = defined
	if defined {
		vnode = init.GetParent()
		this_.priv_InliningBehavior_fillInitialValues_3338(init, &initvals)
	} else if isBlockScoped {
		var scopeRoot pkg_rhino.Node = j2gort.CastPtr[Scope](v.GetScope()).GetRootNode_0()
		_ = scopeRoot
		vnode = scopeRoot.GetFirstChild()
		if !(NodeUtil_IsStatement(vnode)) && NodeUtil_IsStatement(scopeRoot) {
			vnode = scopeRoot
		} else if scopeRoot.IsSwitchBody() {
			vnode = scopeRoot.GetParent()
		}
	} else {
		vnode = j2gort.CastPtr[Scope](j2gort.CastPtr[Scope](v.GetScope()).GetClosestHoistScope()).GetRootNode_0().GetFirstChild()
	}
	if !(NodeUtil_IsStatement(vnode)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(vnode))
	}
	if isBlockScoped && !(varmap.IsEmpty()) {
		var script pkg_rhino.Node = NodeUtil_GetEnclosingScript(vnode)
		_ = script
		if script != 0 {
			NodeUtil_AddFeatureToScript(script, pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), this_.outer_.compiler)
		}
	}
	{
		cur_5 := j2gort.NewMapEntryCursor(varmap)
		for cur_5.NextEntry() {
			var val pkg_rhino.Node = initvals.Get(cur_5.KeyStr())
			_ = val
			var lhs pkg_rhino.Node = pkg_rhino.IR_Name(j2gort.CastString(cur_5.GetValue()))
			_ = lhs
			var declarationType pkg_rhino.Token
			_ = declarationType
			if isBlockScoped {
				declarationType = pkg_rhino.StaticGet_Token_LET()
			} else {
				declarationType = pkg_rhino.StaticGet_Token_VAR()
			}
			var newVarNode pkg_rhino.Node = pkg_rhino.Node(0)
			_ = newVarNode
			if val == 0 {
				newVarNode = pkg_rhino.IR_Declaration_2(lhs, declarationType).SrcrefTreeIfMissing(vnode)
			} else {
				newVarNode = pkg_rhino.IR_Declaration_3(lhs, val, declarationType).SrcrefTreeIfMissing(val)
				this_.priv_InliningBehavior_recordStaleVarReferencesInTree_3333(val, j2gort.CastPtr[Scope](v.GetScope()))
			}
			newVarNode.InsertBefore(vnode)
			this_.outer_.compiler.ReportChangeToEnclosingScope(vnode)
		}
	}
	if defined {
		this_.outer_.compiler.ReportChangeToEnclosingScope(vnode.GetParent())
		vnode.Detach()
	}
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = referenceInfo.Size()
		_ = size
		for ; i < size; i = i + 1 {
			var ref *Reference = referenceInfo.Get(i)
			_ = ref
			if defined && (ref == init) {
				continue
			}
			this_.outer_.compiler.ReportChangeToEnclosingScope(ref.GetNode())
			if ref.IsLvalue() {
				this_.priv_InliningBehavior_replaceAssignmentExpression_3339(v, ref, varmap)
			} else if NodeUtil_IsNameDeclaration(ref.GetParent()) {
				ref.GetParent().Detach()
			} else {
				var getprop pkg_rhino.Node = ref.GetParent()
				_ = getprop
				if !(NodeUtil_IsNormalOrOptChainGetProp(getprop)) {
					j2gort.CheckStateMsg(false, j2gort.BoxNode(getprop))
				}
				var var_ string = getprop.GetString()
				_ = var_
				if !(j2gort.MapContainsKeyStr(varmap, var_)) {
					j2gort.CheckState(false)
				}
				var replacement pkg_rhino.Node = pkg_rhino.IR_Name(j2gort.CastString(j2gort.MapGetStr(varmap, var_)))
				_ = replacement
				replacement.SrcrefIfMissing(getprop)
				ref.GetParent().ReplaceWith(replacement)
			}
		}
	}
}

func (this_ *InlineObjectLiterals_InliningBehavior) BeforeTraverseNode(p0_ *NodeTraversal, p1_ pkg_rhino.Node, p2_ pkg_rhino.Node) {
	IfaceDefault_ReferenceCollector_Behavior_BeforeTraverseNode(this_, p0_, p1_, p2_)
}

func (this_ *InlineObjectLiterals_InliningBehavior) AfterVisitNode(p0_ *NodeTraversal, p1_ pkg_rhino.Node, p2_ pkg_rhino.Node, p3_ *Reference) {
	IfaceDefault_ReferenceCollector_Behavior_AfterVisitNode(this_, p0_, p1_, p2_, p3_)
}
func (this_ *InlineObjectLiterals) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_InlineObjectLiterals_InliningBehavior.Release()
}
