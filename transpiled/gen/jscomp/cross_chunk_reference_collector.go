/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CrossChunkReferenceCollector.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.CrossChunkReferenceCollector
type CrossChunkReferenceCollector struct {
	topLevelStatements     []*CrossChunkReferenceCollector_TopLevelStatement `gson:"topLevelStatements"`
	compiler               AbstractCompiler                                  `gson:"compiler"`
	referenceCollector     *ReferenceCollector                               `gson:"referenceCollector"`
	topLevelStatementDraft *CrossChunkReferenceCollector_TopLevelStatement   `gson:"topLevelStatementDraft"`
	statementCounter       int32                                             `gson:"statementCounter"`
	idHash_                atomic.Int32
	ownsScopeCreator_      bool
}

func (this_ *CrossChunkReferenceCollector) BindSelf_(s any) {
	_ = s
}
func (this_ *CrossChunkReferenceCollector) IsJavaNil_() bool { return this_ == nil }
func (*CrossChunkReferenceCollector) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CrossChunkReferenceCollector", "CrossChunkReferenceCollector"
}
func (this_ *CrossChunkReferenceCollector) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CrossChunkReferenceCollector) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CrossChunkReferenceCollector", "com.google.javascript.jscomp.CrossChunkReferenceCollector", func(v_ any) bool {
		return j2gort.InstanceOf[*CrossChunkReferenceCollector](v_)
	})
}
func (*CrossChunkReferenceCollector) Implements_jscomp_CompilerPass() {}

func (this_ *CrossChunkReferenceCollector) initFields_CrossChunkReferenceCollector() {
	_ = this_
	this_.topLevelStatements = []*CrossChunkReferenceCollector_TopLevelStatement{}

}

// Java constructor: CrossChunkReferenceCollector(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.ScopeCreator creator)
func NewCrossChunkReferenceCollector(compiler AbstractCompiler, creator ScopeCreator) *CrossChunkReferenceCollector {
	this_ := &CrossChunkReferenceCollector{}
	InitCrossChunkReferenceCollector(this_, compiler, creator)
	return this_
}

func InitCrossChunkReferenceCollector(this_ *CrossChunkReferenceCollector, compiler AbstractCompiler, creator ScopeCreator) {
	_ = compiler
	_ = creator
	this_.initFields_CrossChunkReferenceCollector()
	this_.compiler = compiler
	this_.referenceCollector = NewReferenceCollector_5(compiler, NewCrossChunkReferenceCollector_TopLevelStatementCollectingBehavior(this_), creator, func() j2gort.Ext_com_google_common_base_Predicate {
		bound_1 := this_
		return j2gort.NewLambda_Ext_com_google_common_base_Predicate(func(mrp0_ any) bool {
			return bound_1.priv_CrossChunkReferenceCollector_shouldCollectVar_2281(j2gort.CastPtr[Var](mrp0_))
		})
	}(), false)
}

// Java method: CrossChunkReferenceCollector.shouldCollectVar(com.google.javascript.jscomp.Var v)
func (this_ *CrossChunkReferenceCollector) priv_CrossChunkReferenceCollector_shouldCollectVar_2281(v *Var) bool {
	_ = v
	return (v.IsGlobal() && !(this_.compiler.GetCodingConvention().IsExported_2(v.GetName(), false)))
}

// Java method: CrossChunkReferenceCollector.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CrossChunkReferenceCollector) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if len(this_.topLevelStatements) != 0 {
		j2gort.CheckStateStr(false, "process() called more than once")
	}
	this_.referenceCollector.Process_2(externs, root)
}

// Java method: CrossChunkReferenceCollector.process(com.google.javascript.rhino.Node root)
func (this_ *CrossChunkReferenceCollector) Process_1(root pkg_rhino.Node) {
	_ = root
	if len(this_.topLevelStatements) != 0 {
		j2gort.CheckStateStr(false, "process() called more than once")
	}
	this_.referenceCollector.Process_1(root)
}

// Java method: CrossChunkReferenceCollector.getReferences(com.google.javascript.jscomp.Var v)
func (this_ *CrossChunkReferenceCollector) GetReferences(v *Var) *ReferenceCollection {
	_ = v
	return j2gort.CastPtr[ReferenceCollection](this_.referenceCollector.GetReferences(j2gort.PtrToAny(v)))
}

// Java method: CrossChunkReferenceCollector.getGlobalVariableNamesMap()
func (this_ *CrossChunkReferenceCollector) GetGlobalVariableNamesMap() j2gort.Ext_com_google_common_collect_ImmutableMap {
	var builder j2gort.Ext_com_google_common_collect_ImmutableMap_Builder = j2gort.Ext_com_google_common_collect_ImmutableMap_S_Builder()
	_ = builder
	{
		cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_lang_Iterable](this_.referenceCollector.GetAllSymbols())))
		for cur_2.HasNext() {
			rawElem_3 := cur_2.Next()
			v := j2gort.CastPtr[Var](rawElem_3)
			_ = v
			builder.Put_2(j2gort.BoxString(v.GetName()), j2gort.PtrToAny(v))
		}
	}
	return j2gort.AsIface[j2gort.Ext_com_google_common_collect_ImmutableMap](builder.BuildOrThrow())
}

// Java method: CrossChunkReferenceCollector.initializeDraftStatement(com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.rhino.Node statementNode)
func (this_ *CrossChunkReferenceCollector) priv_CrossChunkReferenceCollector_initializeDraftStatement_2285(chunk *JSChunk, statementNode pkg_rhino.Node) *CrossChunkReferenceCollector_TopLevelStatement {
	_ = chunk
	_ = statementNode
	var val_4 CrossChunkReferenceCollector_TopLevelStatement
	val_4.outer_ = this_
	var iifeArg_1 int32
	{

		old_5 := this_.statementCounter
		nw_6 := old_5 + 1
		this_.statementCounter = nw_6
		iifeArg_1 = old_5
	}
	InitCrossChunkReferenceCollector_TopLevelStatement(&val_4, iifeArg_1,

		chunk, statementNode)
	var draft *CrossChunkReferenceCollector_TopLevelStatement = &val_4
	_ = draft
	if NodeUtil_IsNameDeclaration(statementNode) {
		{
			r_7 := draft
			r_7.declaredNameNode = statementNode.GetFirstChild()
		}
		{
			r_8 := draft
			r_8.declaredValueNode = statementNode.GetFirstFirstChild()
		}
	} else if statementNode.IsClass() {
		{
			r_9 := draft
			r_9.declaredNameNode = statementNode.GetFirstChild()
		}
		{
			r_10 := draft
			r_10.declaredValueNode = statementNode
		}
	} else if statementNode.IsFunction() {
		{
			r_11 := draft
			r_11.declaredNameNode = statementNode.GetFirstChild()
		}
		{
			r_12 := draft
			r_12.declaredValueNode = statementNode
		}
	} else if statementNode.IsExprResult() {
		var expr pkg_rhino.Node = j2gort.CheckNotNullNode(statementNode.GetFirstChild())
		_ = expr
		if expr.IsAssign() {
			var lhs pkg_rhino.Node = j2gort.CheckNotNullNode(expr.GetFirstChild())
			_ = lhs
			var rhs pkg_rhino.Node = j2gort.CheckNotNullNode(expr.GetSecondChild())
			_ = rhs
			if lhs.IsName() {
				{
					r_13 := draft
					r_13.declaredNameNode = lhs
				}
				{
					r_14 := draft
					r_14.declaredValueNode = rhs
				}
			} else if lhs.IsGetProp() || CrossChunkReferenceCollector_IsPotentiallyMoveableGetElem(lhs) {
				var nameNode pkg_rhino.Node = j2gort.CheckNotNullNode(lhs.GetFirstChild())
				_ = nameNode
				for nameNode.IsGetProp() {
					nameNode = j2gort.CheckNotNullNode(nameNode.GetFirstChild())
				}
				if nameNode.IsName() {
					{
						r_15 := draft
						r_15.declaredNameNode = nameNode
					}
					{
						r_16 := draft
						r_16.declaredValueNode = rhs
					}
				}
			}
		} else if expr.IsCall() {
			var nameNode_2 pkg_rhino.Node = pkg_rhino.Node(0)
			_ = nameNode_2
			var valueNode pkg_rhino.Node = pkg_rhino.Node(0)
			_ = valueNode
			var relationship *CodingConvention_SubclassRelationship = this_.compiler.GetCodingConvention().GetClassesDefinedByCall(expr)
			_ = relationship
			if relationship != nil {
				var declaredName string = j2gort.CheckNotNullString(relationship.subclassName)
				_ = declaredName
				{
					var callArg pkg_rhino.Node = expr.GetSecondChild()
					_ = callArg
					for ; callArg != 0; callArg = callArg.GetNext() {
						if callArg.IsName() && func() bool {
							sr_17 := declaredName
							se_18 := callArg.GetString()
							return j2gort.CheckNotNullString(sr_17) == se_18
						}() {
							nameNode_2 = callArg
							break
						}
					}
				}
			} else if NodeUtil_IsObjectDefinePropertiesDefinition(expr) {
				var targetObject pkg_rhino.Node = expr.GetSecondChild()
				_ = targetObject
				if targetObject.IsName() {
					nameNode_2 = targetObject
				} else {
					nameNode_2 = targetObject.GetFirstChild()
				}
				valueNode = targetObject.GetNext()
			}
			if nameNode_2 != 0 {
				{
					r_19 := draft
					r_19.declaredNameNode = nameNode_2
				}
				{
					r_20 := draft
					r_20.declaredValueNode = valueNode
				}
			}
		}
	}
	return draft
}

// Java static method: CrossChunkReferenceCollector.isPotentiallyMoveableGetElem(com.google.javascript.rhino.Node maybeGetElem)
func CrossChunkReferenceCollector_IsPotentiallyMoveableGetElem(maybeGetElem pkg_rhino.Node) bool {
	_ = maybeGetElem
	if !(maybeGetElem.IsGetElem()) {
		return false
	}
	var property pkg_rhino.Node = maybeGetElem.GetLastChild()
	_ = property
	if !(property.IsGetProp()) {
		return false
	}
	var name pkg_rhino.Node = property.GetFirstChild()
	_ = name
	return (name.IsName() && (j2gort.CheckNotNullString(name.GetString()) == "Symbol"))
}

// Java method: CrossChunkReferenceCollector.getTopLevelStatementCount()
func (this_ *CrossChunkReferenceCollector) GetTopLevelStatementCount() int32 {
	return int32(len(this_.topLevelStatements))
}

// Java method: CrossChunkReferenceCollector.getTopLevelStatement(int index)
func (this_ *CrossChunkReferenceCollector) GetTopLevelStatement(index int32) *CrossChunkReferenceCollector_TopLevelStatement {
	_ = index
	return this_.topLevelStatements[index]
}

// Java method: CrossChunkReferenceCollector.canMoveValue(com.google.javascript.jscomp.Scope scope, com.google.javascript.rhino.Node valueNode)
func (this_ *CrossChunkReferenceCollector) priv_CrossChunkReferenceCollector_canMoveValue_2288(scope *Scope, valueNode pkg_rhino.Node) bool {
	_ = scope
	_ = valueNode
	if ((valueNode == 0) || NodeUtil_IsLiteralValue(valueNode, true)) || valueNode.IsFunction() {
		return true
	}
	var jsDocInfo *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(valueNode)
	_ = jsDocInfo
	if (jsDocInfo != nil) && jsDocInfo.IsPureOrBreakMyCode() {
		return true
	}
	switch valueNode.TokenSwitchOrd_() {
	case int32(105) /* Token.CLASS */ :
		var classMembers pkg_rhino.Node = valueNode.GetLastChild()
		_ = classMembers
		{
			var member pkg_rhino.Node = classMembers.GetFirstChild()
			_ = member
			for ; member != 0; member = member.GetNext() {
				if member.IsComputedProp() {
					var keyExpr pkg_rhino.Node = member.GetFirstChild()
					_ = keyExpr
					var method pkg_rhino.Node = member.GetLastChild()
					_ = method
					if !(method.IsFunction()) {
						j2gort.CheckStateMsg(false, j2gort.BoxNode(method))
					}
					if !(this_.priv_CrossChunkReferenceCollector_canMoveValue_2288(scope, keyExpr)) {
						return false
					}
				} else if member.IsComputedFieldDef() {
					var keyExpr_2 pkg_rhino.Node = member.GetFirstChild()
					_ = keyExpr_2
					if !(this_.priv_CrossChunkReferenceCollector_canMoveValue_2288(scope, keyExpr_2)) {
						return false
					}
				} else if member.IsMemberFieldDef() {
					if member.IsStaticMember() && !(this_.priv_CrossChunkReferenceCollector_canMoveValue_2288(scope, member.GetFirstChild())) {
						return false
					}
				} else if member.IsBlock() {
					return false
				} else {
					if !(member.IsMemberFunctionDef() || NodeUtil_IsGetOrSetKey(member)) {
						j2gort.CheckStateMsg(false, j2gort.BoxNode(member))
					}
				}
			}
		}
		return true
	case int32(28) /* Token.CALL */ :
		var functionName pkg_rhino.Node = j2gort.CheckNotNullNode(valueNode.GetFirstChild())
		_ = functionName
		return (functionName.IsName() && func() bool {
			sr_21 := functionName.GetString()
			se_22 := "JSCompiler_stubMethod"
			return j2gort.CheckNotNullString(sr_21) == se_22
		}())
	case int32(46) /* Token.ARRAYLIT */ :
		{
			var child pkg_rhino.Node = valueNode.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				if !(this_.priv_CrossChunkReferenceCollector_canMoveValue_2288(scope, child)) {
					return false
				}
			}
		}
		return true
	case int32(47) /* Token.OBJECTLIT */ :
		{
			var child_2 pkg_rhino.Node = valueNode.GetFirstChild()
			_ = child_2
			for ; child_2 != 0; child_2 = child_2.GetNext() {
				switch child_2.TokenSwitchOrd_() {
				case int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
				case int32(129) /* Token.COMPUTED_PROP */ :
					if !(this_.priv_CrossChunkReferenceCollector_canMoveValue_2288(scope, child_2.GetFirstChild())) || !(this_.priv_CrossChunkReferenceCollector_canMoveValue_2288(scope, child_2.GetLastChild())) {
						return false
					}
				case int32(100) /* Token.STRING_KEY */, int32(128) /* Token.OBJECT_SPREAD */ :
					if !(this_.priv_CrossChunkReferenceCollector_canMoveValue_2288(scope, child_2.GetOnlyChild())) {
						return false
					}
				default:
					panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected child of OBJECTLIT: "+j2gort.StrValString(child_2.ToStringTree()))))
				}
			}
		}
		return true
	case int32(32) /* Token.NAME */ :
		var v *Var = j2gort.CastPtr[Var](scope.GetVar(valueNode.GetString()))
		_ = v
		if (v != nil) && v.IsGlobal() {
			var refCollection *ReferenceCollection = this_.GetReferences(v)
			_ = refCollection
			if ((refCollection != nil) && refCollection.IsWellDefined()) && refCollection.IsAssignedOnceInLifetime() {
				return true
			}
		}
	case int32(131) /* Token.TEMPLATELIT */ :
		{
			var child_3 pkg_rhino.Node = valueNode.GetFirstChild()
			_ = child_3
			for ; child_3 != 0; child_3 = child_3.GetNext() {
				if child_3.IsTemplateLitSub() {
					if !(this_.priv_CrossChunkReferenceCollector_canMoveValue_2288(scope, child_3.GetFirstChild())) {
						return false
					}
				} else {
					if !(child_3.IsTemplateLitString()) {
						j2gort.CheckStateMsg(false, j2gort.BoxNode(child_3))
					}
				}
			}
		}
		return true
	case int32(26) /* Token.GETPROP */ :
		var child_4 pkg_rhino.Node = valueNode.GetFirstChild()
		_ = child_4
		if (child_4 != 0) && child_4.IsName() {
			var symbol string = child_4.GetString()
			_ = symbol
			if (symbol != j2gort.NullString) && (j2gort.CheckNotNullString(symbol) == "Symbol") {
				return true
			}
		}
		return false
	default:
	}
	return false
}

// Java class: com.google.javascript.jscomp.CrossChunkReferenceCollector.TopLevelStatementCollectingBehavior
type CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior struct {
	outer_  *CrossChunkReferenceCollector
	idHash_ atomic.Int32
}

func (this_ *CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior) BindSelf_(s any) {
	_ = s
}
func (this_ *CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior) IsJavaNil_() bool {
	return this_ == nil
}
func (*CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CrossChunkReferenceCollector$TopLevelStatementCollectingBehavior", "TopLevelStatementCollectingBehavior"
}
func (this_ *CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CrossChunkReferenceCollector.TopLevelStatementCollectingBehavior", "com.google.javascript.jscomp.CrossChunkReferenceCollector$TopLevelStatementCollectingBehavior", func(v_ any) bool {
		return j2gort.InstanceOf[*CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior](v_)
	})
}
func (this_ *CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior) SetOuter_(o *CrossChunkReferenceCollector) {
	this_.outer_ = o
}
func (this_ *CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior) GetOuter_() *CrossChunkReferenceCollector {
	return this_.outer_
}
func (*CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior) Implements_jscomp_ReferenceCollector_Behavior() {
}

func (this_ *CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior) initFields_CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior() {
	_ = this_
}

// Java constructor: TopLevelStatementCollectingBehavior()
func NewCrossChunkReferenceCollector_TopLevelStatementCollectingBehavior(outer_ *CrossChunkReferenceCollector) *CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior {
	this_ := &CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior{}
	this_.outer_ = outer_
	InitCrossChunkReferenceCollector_TopLevelStatementCollectingBehavior(this_)
	return this_
}

func InitCrossChunkReferenceCollector_TopLevelStatementCollectingBehavior(this_ *CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior) {
	this_.initFields_CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior()
}

// Java method: TopLevelStatementCollectingBehavior.beforeTraverseNode(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node parent)
func (this_ *CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior) BeforeTraverseNode(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if (parent != 0) && NodeUtil_IsTopLevel(parent) {
		if !(this_.outer_.topLevelStatementDraft == nil) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
		}
		this_.outer_.topLevelStatementDraft = this_.outer_.priv_CrossChunkReferenceCollector_initializeDraftStatement_2285(t.GetChunk(), n)
	}
}

// Java method: TopLevelStatementCollectingBehavior.afterVisitNode(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node parent, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Reference addedReference)
func (this_ *CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior) AfterVisitNode(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node, addedReference *Reference) {
	_ = t
	_ = n
	_ = parent
	_ = addedReference
	if this_.outer_.topLevelStatementDraft == nil {
		return
	}
	{
		var iifeCond_1 bool
		{

			rn_23 := n
			an_24 := this_.outer_.topLevelStatementDraft.statementNode
			iifeCond_1 = j2gort.CheckNotNullNode(rn_23) == an_24
		}
		if iifeCond_1 {
			this_.outer_.topLevelStatements = append(this_.outer_.topLevelStatements, this_.outer_.topLevelStatementDraft)
			this_.outer_.topLevelStatementDraft = nil
		} else if (addedReference != nil) && n.IsName() {
			if addedReference.GetNode() == this_.outer_.topLevelStatementDraft.declaredNameNode {
				{
					r_25 := this_.outer_.topLevelStatementDraft
					r_25.declaredNameReference = addedReference
				}
			} else {
				this_.outer_.topLevelStatementDraft.AddNonDeclarationReference(addedReference)
			}
		}
	}
}

// Java method: TopLevelStatementCollectingBehavior.afterExitScope(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.ReferenceMap referenceMap)
func (this_ *CrossChunkReferenceCollector_TopLevelStatementCollectingBehavior) AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap) {
	_ = t
	_ = referenceMap
}

// Java class: com.google.javascript.jscomp.CrossChunkReferenceCollector.TopLevelStatement
type CrossChunkReferenceCollector_TopLevelStatement struct {
	outer_                *CrossChunkReferenceCollector
	chunk                 *JSChunk       `gson:"chunk"`
	nonDeclRef0           *Reference     `gson:"nonDeclRef0"`
	nonDeclRef1           *Reference     `gson:"nonDeclRef1"`
	nonDeclRef2           *Reference     `gson:"nonDeclRef2"`
	nonDeclRef3           *Reference     `gson:"nonDeclRef3"`
	nonDeclRefsRest       []*Reference   `gson:"nonDeclRefsRest"`
	declaredNameReference *Reference     `gson:"declaredNameReference"`
	originalOrder         int32          `gson:"originalOrder"`
	statementNode         pkg_rhino.Node `gson:"statementNode"`
	nonDeclRefCount       int32          `gson:"nonDeclRefCount"`
	declaredNameNode      pkg_rhino.Node `gson:"declaredNameNode"`
	declaredValueNode     pkg_rhino.Node `gson:"declaredValueNode"`
	idHash_               atomic.Int32
}

func (this_ *CrossChunkReferenceCollector_TopLevelStatement) BindSelf_(s any) {
	_ = s
}
func (this_ *CrossChunkReferenceCollector_TopLevelStatement) IsJavaNil_() bool { return this_ == nil }
func (*CrossChunkReferenceCollector_TopLevelStatement) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CrossChunkReferenceCollector$TopLevelStatement", "TopLevelStatement"
}
func (this_ *CrossChunkReferenceCollector_TopLevelStatement) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CrossChunkReferenceCollector_TopLevelStatement) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CrossChunkReferenceCollector.TopLevelStatement", "com.google.javascript.jscomp.CrossChunkReferenceCollector$TopLevelStatement", func(v_ any) bool {
		return j2gort.InstanceOf[*CrossChunkReferenceCollector_TopLevelStatement](v_)
	})
}
func (this_ *CrossChunkReferenceCollector_TopLevelStatement) SetOuter_(o *CrossChunkReferenceCollector) {
	this_.outer_ = o
}
func (this_ *CrossChunkReferenceCollector_TopLevelStatement) GetOuter_() *CrossChunkReferenceCollector {
	return this_.outer_
}

func (this_ *CrossChunkReferenceCollector_TopLevelStatement) initFields_CrossChunkReferenceCollector_TopLevelStatement() {
	_ = this_

}

// Java constructor: TopLevelStatement(int originalOrder, com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.rhino.Node statementNode)
func NewCrossChunkReferenceCollector_TopLevelStatement(outer_ *CrossChunkReferenceCollector, originalOrder int32, chunk *JSChunk, statementNode pkg_rhino.Node) *CrossChunkReferenceCollector_TopLevelStatement {
	this_ := &CrossChunkReferenceCollector_TopLevelStatement{}
	this_.outer_ = outer_
	InitCrossChunkReferenceCollector_TopLevelStatement(this_, originalOrder, chunk, statementNode)
	return this_
}

func InitCrossChunkReferenceCollector_TopLevelStatement(this_ *CrossChunkReferenceCollector_TopLevelStatement, originalOrder int32, chunk *JSChunk, statementNode pkg_rhino.Node) {
	_ = originalOrder
	_ = chunk
	_ = statementNode
	this_.initFields_CrossChunkReferenceCollector_TopLevelStatement()
	this_.originalOrder = originalOrder
	this_.chunk = chunk
	this_.statementNode = statementNode
}

// Java method: TopLevelStatement.addNonDeclarationReference(com.google.javascript.jscomp.Reference ref)
func (this_ *CrossChunkReferenceCollector_TopLevelStatement) AddNonDeclarationReference(ref *Reference) {
	_ = ref
	var c int32 = this_.nonDeclRefCount
	_ = c
	switch c {
	case int32(0):
		this_.nonDeclRef0 = ref
	case int32(1):
		this_.nonDeclRef1 = ref
	case int32(2):
		this_.nonDeclRef2 = ref
	case int32(3):
		this_.nonDeclRef3 = ref
	default:
		var rest []*Reference = this_.nonDeclRefsRest
		_ = rest
		var idx int32 = (c - int32(4))
		_ = idx
		if rest == nil {
			rest = make([]*Reference, int32(8))
			this_.nonDeclRefsRest = rest
		} else if idx == int32(len(rest)) {
			{

				arrSrc_26 := rest
				arrLen_27 := (idx * int32(2))
				j2gort.CheckNotNullSlice(arrSrc_26)
				if arrLen_27 < 0 {
					panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
				}
				arrOut_28 := make([]*Reference, arrLen_27)
				copy(arrOut_28, arrSrc_26)
				rest = arrOut_28
			}

			this_.nonDeclRefsRest = rest
		}
		{
			a_29 := rest
			i_30 := idx
			a_29[i_30] = ref
		}
	}
	this_.nonDeclRefCount = (c + int32(1))
}

// Java method: TopLevelStatement.getOriginalOrder()
func (this_ *CrossChunkReferenceCollector_TopLevelStatement) GetOriginalOrder() int32 {
	return this_.originalOrder
}

// Java method: TopLevelStatement.getChunk()
func (this_ *CrossChunkReferenceCollector_TopLevelStatement) GetChunk() *JSChunk {
	return this_.chunk
}

// Java method: TopLevelStatement.getStatementNode()
func (this_ *CrossChunkReferenceCollector_TopLevelStatement) GetStatementNode() pkg_rhino.Node {
	return this_.statementNode
}

// Java method: TopLevelStatement.getNonDeclarationReferenceCount()
func (this_ *CrossChunkReferenceCollector_TopLevelStatement) GetNonDeclarationReferenceCount() int32 {
	return this_.nonDeclRefCount
}

// Java method: TopLevelStatement.getNonDeclarationReference(int index)
func (this_ *CrossChunkReferenceCollector_TopLevelStatement) GetNonDeclarationReference(index int32) *Reference {
	_ = index
	if (index < int32(0)) || (index >= this_.nonDeclRefCount) {
		panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", j2gort.ValueOfInt(index))))
	}
	switch index {
	case int32(0):
		return this_.nonDeclRef0
	case int32(1):
		return this_.nonDeclRef1
	case int32(2):
		return this_.nonDeclRef2
	case int32(3):
		return this_.nonDeclRef3
	default:
		return this_.nonDeclRefsRest[(index - int32(4))]
	}
	panic("unreachable switch expression")
}

// Java method: TopLevelStatement.isDeclarationStatement()
func (this_ *CrossChunkReferenceCollector_TopLevelStatement) IsDeclarationStatement() bool {
	return (this_.declaredNameReference != nil)
}

// Java method: TopLevelStatement.getDeclaredNameReference()
func (this_ *CrossChunkReferenceCollector_TopLevelStatement) GetDeclaredNameReference() *Reference {
	return j2gort.CheckNotNullPtr(this_.declaredNameReference)
}

// Java method: TopLevelStatement.getDeclaredValueNode()
func (this_ *CrossChunkReferenceCollector_TopLevelStatement) GetDeclaredValueNode() pkg_rhino.Node {
	return this_.declaredValueNode
}

// Java method: TopLevelStatement.isMovableDeclaration()
func (this_ *CrossChunkReferenceCollector_TopLevelStatement) IsMovableDeclaration() bool {
	return (this_.IsDeclarationStatement() && this_.outer_.priv_CrossChunkReferenceCollector_canMoveValue_2288(this_.declaredNameReference.GetScope(), this_.declaredValueNode))
}
func (this_ *CrossChunkReferenceCollector) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.topLevelStatements = nil
	this_.topLevelStatementDraft = nil
	if this_.referenceCollector != nil {
		if this_.ownsScopeCreator_ {
			this_.referenceCollector.ownsScopeCreator_ = true
		}
		this_.referenceCollector.ReleaseSlabs_()
		this_.referenceCollector = nil
	}
}
