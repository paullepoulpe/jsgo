/*
 * Copyright 2008 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/CrossChunkCodeMotion.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_diagnostic "github.com/paullepoulpe/jsgo/transpiled/gen/diagnostic"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.CrossChunkCodeMotion
type CrossChunkCodeMotion struct {
	cccmLog                                    pkg_diagnostic.LogFile                      `gson:"cccmLog"`
	compiler                                   AbstractCompiler                            `gson:"compiler"`
	graph                                      *JSChunkGraph                               `gson:"graph"`
	chunkInsertionPointMap                     j2gort.OrderedMap[*JSChunk, pkg_rhino.Node] `gson:"chunkInsertionPointMap"`
	idHash_                                    atomic.Int32
	parentChunkCanSeeSymbolsDeclaredInChildren bool `gson:"parentChunkCanSeeSymbolsDeclaredInChildren"`
	slab_InstanceofReference                   j2gort.PooledSlab[CrossChunkCodeMotion_InstanceofReference]
	slab_GlobalSymbol                          j2gort.PooledSlab[CrossChunkCodeMotion_GlobalSymbol]
	slab_DeclarationStatementGroup             j2gort.PooledSlab[CrossChunkCodeMotion_DeclarationStatementGroup]
	slab_DeclarationStatementGroupCycle        j2gort.PooledSlab[CrossChunkCodeMotion_DeclarationStatementGroupCycle]
	slab_GlobalSymbolCycle                     j2gort.PooledSlab[CrossChunkCodeMotion_GlobalSymbolCycle]
}

func (this_ *CrossChunkCodeMotion) BindSelf_(s any) {
	_ = s
}
func (this_ *CrossChunkCodeMotion) IsJavaNil_() bool { return this_ == nil }
func (*CrossChunkCodeMotion) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CrossChunkCodeMotion", "CrossChunkCodeMotion"
}
func (this_ *CrossChunkCodeMotion) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CrossChunkCodeMotion) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CrossChunkCodeMotion", "com.google.javascript.jscomp.CrossChunkCodeMotion", func(v_ any) bool {
		return j2gort.InstanceOf[*CrossChunkCodeMotion](v_)
	})
}
func (*CrossChunkCodeMotion) Implements_jscomp_CompilerPass() {}

func (this_ *CrossChunkCodeMotion) initFields_CrossChunkCodeMotion() {
	_ = this_
	this_.chunkInsertionPointMap = j2gort.NewOrderedMap[*JSChunk, pkg_rhino.Node]()
}

// Java constructor: CrossChunkCodeMotion(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.JSChunkGraph graph, boolean parentChunkCanSeeSymbolsDeclaredInChildren)
func NewCrossChunkCodeMotion(compiler AbstractCompiler, graph *JSChunkGraph, parentChunkCanSeeSymbolsDeclaredInChildren bool) *CrossChunkCodeMotion {
	this_ := &CrossChunkCodeMotion{}
	InitCrossChunkCodeMotion(this_, compiler, graph, parentChunkCanSeeSymbolsDeclaredInChildren)
	return this_
}

func InitCrossChunkCodeMotion(this_ *CrossChunkCodeMotion, compiler AbstractCompiler, graph *JSChunkGraph, parentChunkCanSeeSymbolsDeclaredInChildren bool) {
	_ = compiler
	_ = graph
	_ = parentChunkCanSeeSymbolsDeclaredInChildren
	this_.initFields_CrossChunkCodeMotion()
	this_.compiler = compiler
	this_.graph = graph
	this_.parentChunkCanSeeSymbolsDeclaredInChildren = parentChunkCanSeeSymbolsDeclaredInChildren
}

// Java method: CrossChunkCodeMotion.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *CrossChunkCodeMotion) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	defer this_.ReleaseSlabs_()
	func() {
		defer func() {
			this_.cccmLog = nil
		}()
		var logFile pkg_diagnostic.LogFile = this_.compiler.CreateOrReopenIndexedLog(j2gort.GetObjectClass(this_), "cccm.log", []string{})
		_ = logFile
		defer j2gort.CloseResource(logFile)
		this_.cccmLog = logFile
		if this_.cccmLog.IsLogging() {
			this_.cccmLog.Log_1_string("doing cccm")
		}
		if this_.graph.GetChunkCount() > int32(1) {
			var val_1 CrossChunkReferenceCollector
			InitCrossChunkReferenceCollector(&val_1, this_.compiler, NewSyntacticScopeCreator_1(this_.compiler))
			var referenceCollector *CrossChunkReferenceCollector = &val_1
			_ = referenceCollector
			referenceCollector.ownsScopeCreator_ = true
			defer referenceCollector.ReleaseSlabs_()

			referenceCollector.Process_1(root)
			var val_2 CrossChunkCodeMotion_GlobalSymbolCollector
			val_2.outer_ = this_
			InitCrossChunkCodeMotion_GlobalSymbolCollector(&val_2)
			var collector *CrossChunkCodeMotion_GlobalSymbolCollector = &val_2
			_ = collector
			collector.CollectGlobalSymbols(referenceCollector)
			this_.priv_CrossChunkCodeMotion_moveGlobalSymbols_2237(collector)
			this_.priv_CrossChunkCodeMotion_addInstanceofGuards_2236(collector)
		} else {
			if this_.cccmLog.IsLogging() {
				this_.cccmLog.Log_1_string("only one chunk exists")
			}
		}
	}()
}

// Java method: CrossChunkCodeMotion.addInstanceofGuards(com.google.javascript.jscomp.CrossChunkCodeMotion.GlobalSymbolCollector collector)
func (this_ *CrossChunkCodeMotion) priv_CrossChunkCodeMotion_addInstanceofGuards_2236(collector *CrossChunkCodeMotion_GlobalSymbolCollector) {
	_ = collector
	{
		cur_3 := collector.symbolStack.Iter()
		for cur_3.HasNext() {

			globalSymbol := cur_3.Next()
			_ = globalSymbol
			{
				cur_5 := globalSymbol.instanceofReferencesToGuard.Iter()
				for cur_5.HasNext() {

					instanceofReference := cur_5.Next()
					_ = instanceofReference
					if !(globalSymbol.DeclarationsCoverChunk(instanceofReference.GetChunk())) {
						this_.priv_CrossChunkCodeMotion_addGuardToInstanceofReference_2238(instanceofReference.GetReference().GetNode())
					}
				}
			}
		}
	}
}

// Java method: CrossChunkCodeMotion.moveGlobalSymbols(com.google.javascript.jscomp.CrossChunkCodeMotion.GlobalSymbolCollector collector)
func (this_ *CrossChunkCodeMotion) priv_CrossChunkCodeMotion_moveGlobalSymbols_2237(collector *CrossChunkCodeMotion_GlobalSymbolCollector) {
	_ = collector
	var val_7 CrossChunkCodeMotion_OrderAndCombineGlobalSymbols
	val_7.outer_ = this_
	InitCrossChunkCodeMotion_OrderAndCombineGlobalSymbols(&val_7, collector)
	var orderer *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols = &val_7
	_ = orderer
	orderer.OrderAndCombine()
	{
		cur_8 := orderer.stronglyConnectedSymbols.Iter()
		for cur_8.HasNext() {

			globalSymbolCycle := cur_8.Next()
			_ = globalSymbolCycle
			globalSymbolCycle.MoveDeclarationStatements()
		}
	}
}

// Java method: CrossChunkCodeMotion.addGuardToInstanceofReference(com.google.javascript.rhino.Node referenceNode)
func (this_ *CrossChunkCodeMotion) priv_CrossChunkCodeMotion_addGuardToInstanceofReference_2238(referenceNode pkg_rhino.Node) {
	_ = referenceNode
	if !(this_.priv_CrossChunkCodeMotion_isUnguardedInstanceofReference_2242(referenceNode)) {
		j2gort.CheckStateFmt(false, "instanceof Reference is already guarded: %s", j2gort.BoxNode(referenceNode))
	}
	var instanceofNode pkg_rhino.Node = j2gort.CheckNotNullNode(referenceNode.GetParent())
	_ = instanceofNode
	var referenceForTypeOf pkg_rhino.Node = referenceNode.CloneNode_0()
	_ = referenceForTypeOf
	var tmp pkg_rhino.Node = pkg_rhino.IR_Block_0()
	_ = tmp
	instanceofNode.ReplaceWith(tmp)
	var and pkg_rhino.Node = pkg_rhino.IR_And(pkg_rhino.NewNode_3(pkg_rhino.StaticGet_Token_EQ(), pkg_rhino.IR_String("function"), pkg_rhino.NewNode_2(pkg_rhino.StaticGet_Token_TYPEOF(), referenceForTypeOf)), instanceofNode)
	_ = and
	and.SrcrefTreeIfMissing(instanceofNode)
	tmp.ReplaceWith(and)
	this_.compiler.ReportChangeToEnclosingScope(and)
}

// Java method: CrossChunkCodeMotion.isUndefinedTypeofGuardReference(com.google.javascript.rhino.Node reference)
func (this_ *CrossChunkCodeMotion) priv_CrossChunkCodeMotion_isUndefinedTypeofGuardReference_2239(reference pkg_rhino.Node) bool {
	_ = reference
	var maybeTypeofGuard pkg_rhino.Node = reference.GetGrandparent()
	_ = maybeTypeofGuard
	if (maybeTypeofGuard != 0) && this_.priv_CrossChunkCodeMotion_isExistenceTypeofGuardFor_2240(maybeTypeofGuard, reference) {
		var andNode pkg_rhino.Node = maybeTypeofGuard.GetParent()
		_ = andNode
		return (((andNode != 0) && andNode.IsAnd()) && this_.priv_CrossChunkCodeMotion_isInstanceofFor_2243(andNode.GetLastChild(), reference))
	} else {
		return false
	}
}

// Java method: CrossChunkCodeMotion.isExistenceTypeofGuardFor(com.google.javascript.rhino.Node expression, com.google.javascript.rhino.Node reference)
func (this_ *CrossChunkCodeMotion) priv_CrossChunkCodeMotion_isExistenceTypeofGuardFor_2240(expression pkg_rhino.Node, reference pkg_rhino.Node) bool {
	_ = expression
	_ = reference
	if expression.IsNE() || expression.IsSHNE() {
		var undefinedString pkg_rhino.Node = expression.GetFirstChild()
		_ = undefinedString
		var typeofNode pkg_rhino.Node = expression.GetLastChild()
		_ = typeofNode
		return (((undefinedString.IsStringLit() && (j2gort.CheckNotNullString(undefinedString.GetString()) == "undefined")) && typeofNode.IsTypeOf()) && typeofNode.GetFirstChild().IsEquivalentTo_1(reference))
	} else if expression.IsEQ() || expression.IsSHEQ() {
		var functionString pkg_rhino.Node = expression.GetFirstChild()
		_ = functionString
		var typeofNode_2 pkg_rhino.Node = expression.GetLastChild()
		_ = typeofNode_2
		return (((functionString.IsStringLit() && (j2gort.CheckNotNullString(functionString.GetString()) == "function")) && typeofNode_2.IsTypeOf()) && typeofNode_2.GetFirstChild().IsEquivalentTo_1(reference))
	} else {
		return false
	}
}

// Java method: CrossChunkCodeMotion.isGuardedInstanceofReference(com.google.javascript.rhino.Node reference)
func (this_ *CrossChunkCodeMotion) priv_CrossChunkCodeMotion_isGuardedInstanceofReference_2241(reference pkg_rhino.Node) bool {
	_ = reference
	var instanceofNode pkg_rhino.Node = reference.GetParent()
	_ = instanceofNode
	if this_.priv_CrossChunkCodeMotion_isInstanceofFor_2243(instanceofNode, reference) {
		var andNode pkg_rhino.Node = instanceofNode.GetParent()
		_ = andNode
		return (((andNode != 0) && andNode.IsAnd()) && this_.priv_CrossChunkCodeMotion_isExistenceTypeofGuardFor_2240(andNode.GetFirstChild(), reference))
	} else {
		return false
	}
}

// Java method: CrossChunkCodeMotion.isUnguardedInstanceofReference(com.google.javascript.rhino.Node reference)
func (this_ *CrossChunkCodeMotion) priv_CrossChunkCodeMotion_isUnguardedInstanceofReference_2242(reference pkg_rhino.Node) bool {
	_ = reference
	var instanceofNode pkg_rhino.Node = reference.GetParent()
	_ = instanceofNode
	if this_.priv_CrossChunkCodeMotion_isInstanceofFor_2243(instanceofNode, reference) {
		if this_.compiler.GetAstAnalyzer().MayHaveSideEffects(instanceofNode.GetFirstChild()) {
			return false
		}
		return !(this_.priv_CrossChunkCodeMotion_isGuardedInstanceofReference_2241(reference))
	} else {
		return false
	}
}

// Java method: CrossChunkCodeMotion.isInstanceofFor(com.google.javascript.rhino.Node expression, com.google.javascript.rhino.Node reference)
func (this_ *CrossChunkCodeMotion) priv_CrossChunkCodeMotion_isInstanceofFor_2243(expression pkg_rhino.Node, reference pkg_rhino.Node) bool {
	_ = expression
	_ = reference
	return (expression.IsInstanceOf() && expression.GetLastChild().IsEquivalentTo_1(reference))
}

// Java class: com.google.javascript.jscomp.CrossChunkCodeMotion.GlobalSymbolCollector
type CrossChunkCodeMotion_GlobalSymbolCollector struct {
	outer_             *CrossChunkCodeMotion
	globalSymbolforVar j2gort.CustomOrderedMap[*Var, *CrossChunkCodeMotion_GlobalSymbol] `gson:"globalSymbolforVar"`
	symbolStack        j2gort.Deque[*CrossChunkCodeMotion_GlobalSymbol]                  `gson:"symbolStack"`
	idHash_            atomic.Int32
}

func (this_ *CrossChunkCodeMotion_GlobalSymbolCollector) BindSelf_(s any) {
	_ = s
}
func (this_ *CrossChunkCodeMotion_GlobalSymbolCollector) IsJavaNil_() bool { return this_ == nil }
func (*CrossChunkCodeMotion_GlobalSymbolCollector) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CrossChunkCodeMotion$GlobalSymbolCollector", "GlobalSymbolCollector"
}
func (this_ *CrossChunkCodeMotion_GlobalSymbolCollector) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CrossChunkCodeMotion_GlobalSymbolCollector) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CrossChunkCodeMotion.GlobalSymbolCollector", "com.google.javascript.jscomp.CrossChunkCodeMotion$GlobalSymbolCollector", func(v_ any) bool {
		return j2gort.InstanceOf[*CrossChunkCodeMotion_GlobalSymbolCollector](v_)
	})
}
func (this_ *CrossChunkCodeMotion_GlobalSymbolCollector) SetOuter_(o *CrossChunkCodeMotion) {
	this_.outer_ = o
}
func (this_ *CrossChunkCodeMotion_GlobalSymbolCollector) GetOuter_() *CrossChunkCodeMotion {
	return this_.outer_
}

func (this_ *CrossChunkCodeMotion_GlobalSymbolCollector) initFields_CrossChunkCodeMotion_GlobalSymbolCollector() {
	_ = this_
	this_.globalSymbolforVar = j2gort.NewCustomOrderedMap[*Var, *CrossChunkCodeMotion_GlobalSymbol]()
	this_.symbolStack = j2gort.Deque[*CrossChunkCodeMotion_GlobalSymbol]{}
}

// Java constructor: GlobalSymbolCollector()
func NewCrossChunkCodeMotion_GlobalSymbolCollector(outer_ *CrossChunkCodeMotion) *CrossChunkCodeMotion_GlobalSymbolCollector {
	this_ := &CrossChunkCodeMotion_GlobalSymbolCollector{}
	this_.outer_ = outer_
	InitCrossChunkCodeMotion_GlobalSymbolCollector(this_)
	return this_
}

func InitCrossChunkCodeMotion_GlobalSymbolCollector(this_ *CrossChunkCodeMotion_GlobalSymbolCollector) {
	this_.initFields_CrossChunkCodeMotion_GlobalSymbolCollector()
}

// Java method: GlobalSymbolCollector.collectGlobalSymbols(com.google.javascript.jscomp.CrossChunkReferenceCollector referenceCollector)
func (this_ *CrossChunkCodeMotion_GlobalSymbolCollector) CollectGlobalSymbols(referenceCollector *CrossChunkReferenceCollector) {
	_ = referenceCollector
	{
		var i int32 = int32(0)
		_ = i
		var count int32 = referenceCollector.GetTopLevelStatementCount()
		_ = count
		for ; i < count; i = i + 1 {
			var statement *CrossChunkReferenceCollector_TopLevelStatement = referenceCollector.GetTopLevelStatement(i)
			_ = statement
			if statement.IsDeclarationStatement() {
				this_.priv_GlobalSymbolCollector_processDeclarationStatement_2247(statement)
			} else {
				this_.priv_GlobalSymbolCollector_processImmovableStatement_2245(statement)
			}
		}
	}
}

// Java method: GlobalSymbolCollector.processImmovableStatement(com.google.javascript.jscomp.CrossChunkReferenceCollector.TopLevelStatement statement)
func (this_ *CrossChunkCodeMotion_GlobalSymbolCollector) priv_GlobalSymbolCollector_processImmovableStatement_2245(statement *CrossChunkReferenceCollector_TopLevelStatement) {
	_ = statement
	var chunk *JSChunk = statement.GetChunk()
	_ = chunk
	{
		var i int32 = int32(0)
		_ = i
		var count int32 = statement.GetNonDeclarationReferenceCount()
		_ = count
		for ; i < count; i = i + 1 {
			this_.priv_GlobalSymbolCollector_processImmovableReference_2246(statement.GetNonDeclarationReference(i), chunk)
		}
	}
}

// Java method: GlobalSymbolCollector.processImmovableReference(com.google.javascript.jscomp.Reference ref, com.google.javascript.jscomp.JSChunk chunk)
func (this_ *CrossChunkCodeMotion_GlobalSymbolCollector) priv_GlobalSymbolCollector_processImmovableReference_2246(ref *Reference, chunk *JSChunk) {
	_ = ref
	_ = chunk
	var globalSymbol *CrossChunkCodeMotion_GlobalSymbol = this_.priv_GlobalSymbolCollector_getGlobalSymbol_2249(j2gort.CastPtr[Var](ref.GetSymbol()))
	_ = globalSymbol
	if this_.outer_.parentChunkCanSeeSymbolsDeclaredInChildren {
		var n pkg_rhino.Node = ref.GetNode()
		_ = n
		if this_.outer_.priv_CrossChunkCodeMotion_isGuardedInstanceofReference_2241(n) || this_.outer_.priv_CrossChunkCodeMotion_isUndefinedTypeofGuardReference_2239(n) {
			return
		} else if this_.outer_.priv_CrossChunkCodeMotion_isUnguardedInstanceofReference_2242(n) {
			globalSymbol.instanceofReferencesToGuard.Push(allocPooledVia_CrossChunkCodeMotion_InstanceofReference_Immovable(&this_.outer_.slab_InstanceofReference, chunk, ref))
			return
		}
	}
	globalSymbol.AddImmovableReference(chunk)
}

// Java method: GlobalSymbolCollector.processDeclarationStatement(com.google.javascript.jscomp.CrossChunkReferenceCollector.TopLevelStatement statement)
func (this_ *CrossChunkCodeMotion_GlobalSymbolCollector) priv_GlobalSymbolCollector_processDeclarationStatement_2247(statement *CrossChunkReferenceCollector_TopLevelStatement) {
	_ = statement
	var declaredSymbol *CrossChunkCodeMotion_GlobalSymbol = this_.priv_GlobalSymbolCollector_getGlobalSymbol_2249(j2gort.CastPtr[Var](statement.GetDeclaredNameReference().GetSymbol()))
	_ = declaredSymbol
	var dsg *CrossChunkCodeMotion_DeclarationStatementGroup = declaredSymbol.AddDeclarationStatement(statement)
	_ = dsg
	this_.priv_GlobalSymbolCollector_processDeclarationStatementContainedReferences_2248(statement, declaredSymbol, dsg)
}

// Java method: GlobalSymbolCollector.processDeclarationStatementContainedReferences(com.google.javascript.jscomp.CrossChunkReferenceCollector.TopLevelStatement statement, com.google.javascript.jscomp.CrossChunkCodeMotion.GlobalSymbol declaredSymbol, com.google.javascript.jscomp.CrossChunkCodeMotion.DeclarationStatementGroup dsg)
func (this_ *CrossChunkCodeMotion_GlobalSymbolCollector) priv_GlobalSymbolCollector_processDeclarationStatementContainedReferences_2248(statement *CrossChunkReferenceCollector_TopLevelStatement, declaredSymbol *CrossChunkCodeMotion_GlobalSymbol, dsg *CrossChunkCodeMotion_DeclarationStatementGroup) {
	_ = statement
	_ = declaredSymbol
	_ = dsg
	{
		var i int32 = int32(0)
		_ = i
		var count int32 = statement.GetNonDeclarationReferenceCount()
		_ = count
		for ; i < count; i = i + 1 {
			var ref *Reference = statement.GetNonDeclarationReference(i)
			_ = ref
			var refSymbol *CrossChunkCodeMotion_GlobalSymbol = this_.priv_GlobalSymbolCollector_getGlobalSymbol_2249(j2gort.CastPtr[Var](ref.GetSymbol()))
			_ = refSymbol
			if j2gort.CheckNotNullPtr(refSymbol) == declaredSymbol {
				continue
			}
			if this_.outer_.parentChunkCanSeeSymbolsDeclaredInChildren {
				var n pkg_rhino.Node = ref.GetNode()
				_ = n
				if this_.outer_.priv_CrossChunkCodeMotion_isGuardedInstanceofReference_2241(n) || this_.outer_.priv_CrossChunkCodeMotion_isUndefinedTypeofGuardReference_2239(n) {
					continue
				} else if this_.outer_.priv_CrossChunkCodeMotion_isUnguardedInstanceofReference_2242(n) {
					refSymbol.instanceofReferencesToGuard.Push(allocPooledVia_CrossChunkCodeMotion_InstanceofReference_Movable(&this_.outer_.slab_InstanceofReference, dsg, ref))
					continue
				}
			}
			dsg.AddReferenceToGlobalSymbol(refSymbol)
			refSymbol.AddReferringGlobalSymbol(declaredSymbol)
		}
	}
}

// Java method: GlobalSymbolCollector.getGlobalSymbol(com.google.javascript.jscomp.Var var)
func (this_ *CrossChunkCodeMotion_GlobalSymbolCollector) priv_GlobalSymbolCollector_getGlobalSymbol_2249(var_ *Var) *CrossChunkCodeMotion_GlobalSymbol {
	_ = var_
	var globalSymbol *CrossChunkCodeMotion_GlobalSymbol = this_.globalSymbolforVar.Get(var_)
	_ = globalSymbol
	if globalSymbol == nil {
		globalSymbol = allocPooled_CrossChunkCodeMotion_GlobalSymbol(&this_.outer_.slab_GlobalSymbol, this_.outer_, var_)
		this_.globalSymbolforVar.Put(var_, globalSymbol)
		this_.symbolStack.Push(globalSymbol)
	}
	return globalSymbol
}

// Java class: com.google.javascript.jscomp.CrossChunkCodeMotion.GlobalSymbol
type CrossChunkCodeMotion_GlobalSymbol struct {
	outer_                                       *CrossChunkCodeMotion
	var_                                         *Var                                                          `gson:"var"`
	dsgStack                                     j2gort.Deque[*CrossChunkCodeMotion_DeclarationStatementGroup] `gson:"dsgStack"`
	chunksWithImmovableReferences                *j2gort.Ext_java_util_BitSet                                  `gson:"chunksWithImmovableReferences"`
	referencingGlobalSymbols                     j2gort.OrderedSet[*CrossChunkCodeMotion_GlobalSymbol]         `gson:"referencingGlobalSymbols"`
	instanceofReferencesToGuard                  j2gort.Deque[*CrossChunkCodeMotion_InstanceofReference]       `gson:"instanceofReferencesToGuard"`
	preorderNumber                               int32                                                         `gson:"preorderNumber"`
	idHash_                                      atomic.Int32
	hasBeenAssignedToAStronglyConnectedComponent bool `gson:"hasBeenAssignedToAStronglyConnectedComponent"`
	isMoveDeclarationStatementsDone              bool `gson:"isMoveDeclarationStatementsDone"`
}

func (this_ *CrossChunkCodeMotion_GlobalSymbol) BindSelf_(s any) {
	_ = s
}
func (this_ *CrossChunkCodeMotion_GlobalSymbol) IsJavaNil_() bool { return this_ == nil }
func (*CrossChunkCodeMotion_GlobalSymbol) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CrossChunkCodeMotion$GlobalSymbol", "GlobalSymbol"
}
func (this_ *CrossChunkCodeMotion_GlobalSymbol) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CrossChunkCodeMotion_GlobalSymbol) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CrossChunkCodeMotion.GlobalSymbol", "com.google.javascript.jscomp.CrossChunkCodeMotion$GlobalSymbol", func(v_ any) bool {
		return j2gort.InstanceOf[*CrossChunkCodeMotion_GlobalSymbol](v_)
	})
}
func (this_ *CrossChunkCodeMotion_GlobalSymbol) SetOuter_(o *CrossChunkCodeMotion) { this_.outer_ = o }
func (this_ *CrossChunkCodeMotion_GlobalSymbol) GetOuter_() *CrossChunkCodeMotion {
	return this_.outer_
}

func (this_ *CrossChunkCodeMotion_GlobalSymbol) initFields_CrossChunkCodeMotion_GlobalSymbol() {
	_ = this_
	this_.dsgStack = j2gort.Deque[*CrossChunkCodeMotion_DeclarationStatementGroup]{}
	this_.chunksWithImmovableReferences = j2gort.NewExt_java_util_BitSet_1(this_.outer_.graph.GetChunkCount())
	this_.referencingGlobalSymbols = j2gort.NewOrderedSet[*CrossChunkCodeMotion_GlobalSymbol]()
	this_.instanceofReferencesToGuard = j2gort.Deque[*CrossChunkCodeMotion_InstanceofReference]{}
	this_.preorderNumber = int32(-1)

}

// Java constructor: GlobalSymbol(com.google.javascript.jscomp.Var var)
func NewCrossChunkCodeMotion_GlobalSymbol(outer_ *CrossChunkCodeMotion, var_ *Var) *CrossChunkCodeMotion_GlobalSymbol {
	this_ := &CrossChunkCodeMotion_GlobalSymbol{}
	this_.outer_ = outer_
	InitCrossChunkCodeMotion_GlobalSymbol(this_, var_)
	return this_
}
func allocPooled_CrossChunkCodeMotion_GlobalSymbol(slab_ *j2gort.PooledSlab[CrossChunkCodeMotion_GlobalSymbol], outer_ *CrossChunkCodeMotion, var_ *Var) *CrossChunkCodeMotion_GlobalSymbol {
	this_ := slab_.Alloc()
	this_.outer_ = outer_
	InitCrossChunkCodeMotion_GlobalSymbol(this_, var_)
	return this_
}

func InitCrossChunkCodeMotion_GlobalSymbol(this_ *CrossChunkCodeMotion_GlobalSymbol, var_ *Var) {
	_ = var_
	this_.initFields_CrossChunkCodeMotion_GlobalSymbol()
	this_.var_ = var_
}

// Java method: GlobalSymbol.toString()
func (this_ *CrossChunkCodeMotion_GlobalSymbol) String() string {
	return this_.var_.GetName()
}

// Java method: GlobalSymbol.addImmovableReference(com.google.javascript.jscomp.JSChunk chunk)
func (this_ *CrossChunkCodeMotion_GlobalSymbol) AddImmovableReference(chunk *JSChunk) {
	_ = chunk
	this_.chunksWithImmovableReferences.Set_1(chunk.GetIndex())
}

// Java method: GlobalSymbol.addDeclarationStatement(com.google.javascript.jscomp.CrossChunkReferenceCollector.TopLevelStatement statement)
func (this_ *CrossChunkCodeMotion_GlobalSymbol) AddDeclarationStatement(statement *CrossChunkReferenceCollector_TopLevelStatement) *CrossChunkCodeMotion_DeclarationStatementGroup {
	_ = statement
	var chunk *JSChunk = statement.GetChunk()
	_ = chunk
	var lastDsg *CrossChunkCodeMotion_DeclarationStatementGroup = this_.dsgStack.Peek()
	_ = lastDsg
	if lastDsg == nil {
		lastDsg = allocPooled_CrossChunkCodeMotion_DeclarationStatementGroup(&this_.outer_.slab_DeclarationStatementGroup, this_, chunk)
		this_.dsgStack.Push(lastDsg)
	}
	var statementDsg *CrossChunkCodeMotion_DeclarationStatementGroup = nil
	_ = statementDsg
	{
		var iifeCond_1 bool
		{

			rp_10 := chunk
			ap_11 := lastDsg.currentChunk
			iifeCond_1 = j2gort.CheckNotNullPtr(rp_10) == ap_11
		}
		if iifeCond_1 {
			statementDsg = lastDsg
		} else {
			statementDsg = allocPooled_CrossChunkCodeMotion_DeclarationStatementGroup(&this_.outer_.slab_DeclarationStatementGroup, this_, chunk)
			this_.dsgStack.Push(statementDsg)
		}
	}
	statementDsg.statementStack.Push(statement)
	return statementDsg
}

// Java method: GlobalSymbol.addReferringGlobalSymbol(com.google.javascript.jscomp.CrossChunkCodeMotion.GlobalSymbol declaredSymbol)
func (this_ *CrossChunkCodeMotion_GlobalSymbol) AddReferringGlobalSymbol(declaredSymbol *CrossChunkCodeMotion_GlobalSymbol) {
	_ = declaredSymbol
	this_.referencingGlobalSymbols.Add(declaredSymbol)
}

// Java method: GlobalSymbol.declarationsCoverChunk(com.google.javascript.jscomp.JSChunk chunk)
func (this_ *CrossChunkCodeMotion_GlobalSymbol) DeclarationsCoverChunk(chunk *JSChunk) bool {
	_ = chunk
	{
		cur_12 := this_.dsgStack.Iter()
		for cur_12.HasNext() {

			dsg := cur_12.Next()
			_ = dsg
			{
				var iifeCond_1 bool
				{

					rp_14 := chunk
					ap_15 := dsg.currentChunk
					iifeCond_1 = j2gort.CheckNotNullPtr(rp_14) == ap_15
				}
				if iifeCond_1 || this_.outer_.graph.DependsOn(chunk, dsg.currentChunk) {
					return true
				}
			}
		}
	}
	return false
}

// Java class: com.google.javascript.jscomp.CrossChunkCodeMotion.GlobalSymbolCycle
type CrossChunkCodeMotion_GlobalSymbolCycle struct {
	outer_  *CrossChunkCodeMotion
	symbols j2gort.Deque[*CrossChunkCodeMotion_GlobalSymbol] `gson:"symbols"`
	idHash_ atomic.Int32
}

func (this_ *CrossChunkCodeMotion_GlobalSymbolCycle) BindSelf_(s any) {
	_ = s
}
func (this_ *CrossChunkCodeMotion_GlobalSymbolCycle) IsJavaNil_() bool { return this_ == nil }
func (*CrossChunkCodeMotion_GlobalSymbolCycle) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CrossChunkCodeMotion$GlobalSymbolCycle", "GlobalSymbolCycle"
}
func (this_ *CrossChunkCodeMotion_GlobalSymbolCycle) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CrossChunkCodeMotion_GlobalSymbolCycle) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CrossChunkCodeMotion.GlobalSymbolCycle", "com.google.javascript.jscomp.CrossChunkCodeMotion$GlobalSymbolCycle", func(v_ any) bool {
		return j2gort.InstanceOf[*CrossChunkCodeMotion_GlobalSymbolCycle](v_)
	})
}
func (this_ *CrossChunkCodeMotion_GlobalSymbolCycle) SetOuter_(o *CrossChunkCodeMotion) {
	this_.outer_ = o
}
func (this_ *CrossChunkCodeMotion_GlobalSymbolCycle) GetOuter_() *CrossChunkCodeMotion {
	return this_.outer_
}

func (this_ *CrossChunkCodeMotion_GlobalSymbolCycle) initFields_CrossChunkCodeMotion_GlobalSymbolCycle() {
	_ = this_
	this_.symbols = j2gort.Deque[*CrossChunkCodeMotion_GlobalSymbol]{}
}

// Java constructor: GlobalSymbolCycle()
func NewCrossChunkCodeMotion_GlobalSymbolCycle(outer_ *CrossChunkCodeMotion) *CrossChunkCodeMotion_GlobalSymbolCycle {
	this_ := &CrossChunkCodeMotion_GlobalSymbolCycle{}
	this_.outer_ = outer_
	InitCrossChunkCodeMotion_GlobalSymbolCycle(this_)
	return this_
}
func allocPooled_CrossChunkCodeMotion_GlobalSymbolCycle(slab_ *j2gort.PooledSlab[CrossChunkCodeMotion_GlobalSymbolCycle], outer_ *CrossChunkCodeMotion) *CrossChunkCodeMotion_GlobalSymbolCycle {
	this_ := slab_.Alloc()
	this_.outer_ = outer_
	InitCrossChunkCodeMotion_GlobalSymbolCycle(this_)
	return this_
}

func InitCrossChunkCodeMotion_GlobalSymbolCycle(this_ *CrossChunkCodeMotion_GlobalSymbolCycle) {
	this_.initFields_CrossChunkCodeMotion_GlobalSymbolCycle()
}

// Java method: GlobalSymbolCycle.addSymbol(com.google.javascript.jscomp.CrossChunkCodeMotion.GlobalSymbol symbol)
func (this_ *CrossChunkCodeMotion_GlobalSymbolCycle) AddSymbol(symbol *CrossChunkCodeMotion_GlobalSymbol) {
	_ = symbol
	this_.symbols.AddLast(symbol)
}

// Java method: GlobalSymbolCycle.moveDeclarationStatements()
func (this_ *CrossChunkCodeMotion_GlobalSymbolCycle) MoveDeclarationStatements() {
	{
		cur_16 := this_.symbols.Iter()
		for cur_16.HasNext() {

			symbol := cur_16.Next()
			_ = symbol
			if !(!(symbol.isMoveDeclarationStatementsDone)) {
				j2gort.CheckStateFmt(false, "duplicate attempt to move %s", j2gort.PtrToAny(symbol))
			}
		}
	}
	var chunksWithImmovableReferences_bsVal j2gort.Ext_java_util_BitSet
	j2gort.InitExt_java_util_BitSet_1(&chunksWithImmovableReferences_bsVal, this_.outer_.graph.GetChunkCount())
	var chunksWithImmovableReferences *j2gort.Ext_java_util_BitSet = &chunksWithImmovableReferences_bsVal
	_ = chunksWithImmovableReferences
	var cyclesLatestFirst []*CrossChunkCodeMotion_DeclarationStatementGroupCycle = []*CrossChunkCodeMotion_DeclarationStatementGroupCycle{}
	_ = cyclesLatestFirst
	this_.priv_GlobalSymbolCycle_populateDsgCyclesLatestFirst_2257(&cyclesLatestFirst)
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = int32(len(cyclesLatestFirst))
		_ = size
		for ; i < size; i = i + 1 {
			var dsgCycle *CrossChunkCodeMotion_DeclarationStatementGroupCycle = cyclesLatestFirst[i]
			_ = dsgCycle
			{
				cur_18 := this_.symbols.Iter()
				for cur_18.HasNext() {

					symbol_2 := cur_18.Next()
					_ = symbol_2
					chunksWithImmovableReferences.Or(symbol_2.chunksWithImmovableReferences)
				}
			}
			dsgCycle.MoveToPreferredChunk(chunksWithImmovableReferences)
		}
	}
	{
		cur_20 := this_.symbols.Iter()
		for cur_20.HasNext() {

			symbol_3 := cur_20.Next()
			_ = symbol_3
			{
				r_22 := symbol_3
				r_22.isMoveDeclarationStatementsDone = true
			}
		}
	}
}

// Java method: GlobalSymbolCycle.populateDsgCyclesLatestFirst(java.util.ArrayList<com.google.javascript.jscomp.CrossChunkCodeMotion.DeclarationStatementGroupCycle> cyclesLatestFirst)
func (this_ *CrossChunkCodeMotion_GlobalSymbolCycle) priv_GlobalSymbolCycle_populateDsgCyclesLatestFirst_2257(cyclesLatestFirst *[]*CrossChunkCodeMotion_DeclarationStatementGroupCycle) {
	_ = cyclesLatestFirst
	var dsgsLatestFirst []*CrossChunkCodeMotion_DeclarationStatementGroup = []*CrossChunkCodeMotion_DeclarationStatementGroup{}
	_ = dsgsLatestFirst
	{
		cur_23 := this_.symbols.Iter()
		for cur_23.HasNext() {

			symbol := cur_23.Next()
			_ = symbol
			if len(dsgsLatestFirst) == 0 {
				{
					cur_25 := symbol.dsgStack.Iter()
					for cur_25.HasNext() {

						dsg := cur_25.Next()
						_ = dsg
						dsgsLatestFirst = append(dsgsLatestFirst, dsg)
					}
				}
			} else {
				var prev []*CrossChunkCodeMotion_DeclarationStatementGroup = j2gort.CloneNonNilSlice(dsgsLatestFirst)
				_ = prev
				clear(dsgsLatestFirst)
				dsgsLatestFirst = dsgsLatestFirst[:0]
				var list2 []*CrossChunkCodeMotion_DeclarationStatementGroup = []*CrossChunkCodeMotion_DeclarationStatementGroup{}
				_ = list2
				{
					cur_27 := symbol.dsgStack.Iter()
					for cur_27.HasNext() {

						dsg_2 := cur_27.Next()
						_ = dsg_2
						list2 = append(list2, dsg_2)
					}
				}
				var i1 int32 = int32(0)
				_ = i1
				var n1 int32 = int32(len(prev))
				_ = n1
				var i2 int32 = int32(0)
				_ = i2
				var n2 int32 = int32(len(list2))
				_ = n2
				for (i1 < n1) && (i2 < n2) {
					var dsg1 *CrossChunkCodeMotion_DeclarationStatementGroup = prev[i1]
					_ = dsg1
					var dsg2 *CrossChunkCodeMotion_DeclarationStatementGroup = list2[i2]
					_ = dsg2
					if dsg1.currentChunk.GetIndex() > dsg2.currentChunk.GetIndex() {
						i1 = i1 + 1
						dsgsLatestFirst = append(dsgsLatestFirst, dsg1)
						if !((i1 == n1) || (prev[i1].currentChunk.GetIndex() <= dsg1.currentChunk.GetIndex())) {
							j2gort.CheckStateStr(false, "DSG stacks are out of order.")
						}
					} else {
						i2 = i2 + 1
						dsgsLatestFirst = append(dsgsLatestFirst, dsg2)
						if !((i2 == n2) || (list2[i2].currentChunk.GetIndex() <= dsg2.currentChunk.GetIndex())) {
							j2gort.CheckStateStr(false, "DSG stacks are out of order.")
						}
					}
				}
				for i1 < n1 {
					dsgsLatestFirst = append(dsgsLatestFirst, prev[func() int32 { old_29 := i1; nw_30 := old_29 + 1; i1 = nw_30; return old_29 }()])
				}
				for i2 < n2 {
					dsgsLatestFirst = append(dsgsLatestFirst, list2[func() int32 { old_31 := i2; nw_32 := old_31 + 1; i2 = nw_32; return old_31 }()])
				}
			}
		}
	}
	var cycle *CrossChunkCodeMotion_DeclarationStatementGroupCycle = nil
	_ = cycle
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = int32(len(dsgsLatestFirst))
		_ = size
		for ; i < size; i = i + 1 {
			var dsg_3 *CrossChunkCodeMotion_DeclarationStatementGroup = dsgsLatestFirst[i]
			_ = dsg_3
			if (cycle == nil) || !(func() bool {
				rp_33 := cycle.currentChunk
				ap_34 := dsg_3.currentChunk
				return j2gort.CheckNotNullPtr(rp_33) == ap_34
			}()) {
				cycle = allocPooled_CrossChunkCodeMotion_DeclarationStatementGroupCycle(&this_.outer_.slab_DeclarationStatementGroupCycle, this_.outer_, dsg_3.currentChunk)
				*cyclesLatestFirst = append(*cyclesLatestFirst, cycle)
			}
			cycle.dsgs.AddLast(dsg_3)
		}
	}
}

// Java class: com.google.javascript.jscomp.CrossChunkCodeMotion.DeclarationStatementGroup
type CrossChunkCodeMotion_DeclarationStatementGroup struct {
	declaredGlobalSymbol    *CrossChunkCodeMotion_GlobalSymbol                            `gson:"declaredGlobalSymbol"`
	referencedGlobalSymbols j2gort.OrderedSet[*CrossChunkCodeMotion_GlobalSymbol]         `gson:"referencedGlobalSymbols"`
	currentChunk            *JSChunk                                                      `gson:"currentChunk"`
	statementStack          j2gort.Deque[*CrossChunkReferenceCollector_TopLevelStatement] `gson:"statementStack"`
	idHash_                 atomic.Int32
}

func (this_ *CrossChunkCodeMotion_DeclarationStatementGroup) BindSelf_(s any) {
	_ = s
}
func (this_ *CrossChunkCodeMotion_DeclarationStatementGroup) IsJavaNil_() bool { return this_ == nil }
func (*CrossChunkCodeMotion_DeclarationStatementGroup) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CrossChunkCodeMotion$DeclarationStatementGroup", "DeclarationStatementGroup"
}
func (this_ *CrossChunkCodeMotion_DeclarationStatementGroup) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CrossChunkCodeMotion_DeclarationStatementGroup) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CrossChunkCodeMotion.DeclarationStatementGroup", "com.google.javascript.jscomp.CrossChunkCodeMotion$DeclarationStatementGroup", func(v_ any) bool {
		return j2gort.InstanceOf[*CrossChunkCodeMotion_DeclarationStatementGroup](v_)
	})
}

func (this_ *CrossChunkCodeMotion_DeclarationStatementGroup) initFields_CrossChunkCodeMotion_DeclarationStatementGroup() {
	_ = this_
	this_.referencedGlobalSymbols = j2gort.NewOrderedSet[*CrossChunkCodeMotion_GlobalSymbol]()
	this_.statementStack = j2gort.Deque[*CrossChunkReferenceCollector_TopLevelStatement]{}
}

// Java constructor: DeclarationStatementGroup(com.google.javascript.jscomp.CrossChunkCodeMotion.GlobalSymbol declaredGlobalSymbol, com.google.javascript.jscomp.JSChunk currentChunk)
func NewCrossChunkCodeMotion_DeclarationStatementGroup(declaredGlobalSymbol *CrossChunkCodeMotion_GlobalSymbol, currentChunk *JSChunk) *CrossChunkCodeMotion_DeclarationStatementGroup {
	this_ := &CrossChunkCodeMotion_DeclarationStatementGroup{}
	InitCrossChunkCodeMotion_DeclarationStatementGroup(this_, declaredGlobalSymbol, currentChunk)
	return this_
}
func allocPooled_CrossChunkCodeMotion_DeclarationStatementGroup(slab_ *j2gort.PooledSlab[CrossChunkCodeMotion_DeclarationStatementGroup], declaredGlobalSymbol *CrossChunkCodeMotion_GlobalSymbol, currentChunk *JSChunk) *CrossChunkCodeMotion_DeclarationStatementGroup {
	this_ := slab_.Alloc()
	InitCrossChunkCodeMotion_DeclarationStatementGroup(this_, declaredGlobalSymbol, currentChunk)
	return this_
}

func InitCrossChunkCodeMotion_DeclarationStatementGroup(this_ *CrossChunkCodeMotion_DeclarationStatementGroup, declaredGlobalSymbol *CrossChunkCodeMotion_GlobalSymbol, currentChunk *JSChunk) {
	_ = declaredGlobalSymbol
	_ = currentChunk
	this_.initFields_CrossChunkCodeMotion_DeclarationStatementGroup()
	this_.declaredGlobalSymbol = declaredGlobalSymbol
	this_.currentChunk = currentChunk
}

// Java method: DeclarationStatementGroup.allStatementsCanMove()
func (this_ *CrossChunkCodeMotion_DeclarationStatementGroup) AllStatementsCanMove() bool {
	{
		cur_35 := this_.statementStack.Iter()
		for cur_35.HasNext() {

			s := cur_35.Next()
			_ = s
			if !(s.IsMovableDeclaration()) {
				return false
			}
		}
	}
	return true
}

// Java method: DeclarationStatementGroup.addReferenceToGlobalSymbol(com.google.javascript.jscomp.CrossChunkCodeMotion.GlobalSymbol refSymbol)
func (this_ *CrossChunkCodeMotion_DeclarationStatementGroup) AddReferenceToGlobalSymbol(refSymbol *CrossChunkCodeMotion_GlobalSymbol) {
	_ = refSymbol
	this_.referencedGlobalSymbols.Add(refSymbol)
}

// Java method: DeclarationStatementGroup.makeReferencesImmovable()
func (this_ *CrossChunkCodeMotion_DeclarationStatementGroup) MakeReferencesImmovable() {
	this_.declaredGlobalSymbol.AddImmovableReference(this_.currentChunk)
	{
		cur_37 := this_.referencedGlobalSymbols.Iter()
		for cur_37.HasNext() {

			symbol := cur_37.
				Next()
			_ = symbol
			if !(!(symbol.isMoveDeclarationStatementsDone)) {
				j2gort.CheckStateFmt(false, "symbol %s moved before referring symbol %s", j2gort.PtrToAny(symbol), j2gort.PtrToAny(this_.declaredGlobalSymbol))
			}
			symbol.AddImmovableReference(this_.currentChunk)
		}
	}
}

// Java class: com.google.javascript.jscomp.CrossChunkCodeMotion.OrderAndCombineGlobalSymbols
type CrossChunkCodeMotion_OrderAndCombineGlobalSymbols struct {
	outer_                   *CrossChunkCodeMotion
	collector                *CrossChunkCodeMotion_GlobalSymbolCollector           `gson:"collector"`
	componentContents        j2gort.Deque[*CrossChunkCodeMotion_GlobalSymbol]      `gson:"componentContents"`
	componentRoots           j2gort.Deque[*CrossChunkCodeMotion_GlobalSymbol]      `gson:"componentRoots"`
	stronglyConnectedSymbols j2gort.Deque[*CrossChunkCodeMotion_GlobalSymbolCycle] `gson:"stronglyConnectedSymbols"`
	preorderCounter          int32                                                 `gson:"preorderCounter"`
	idHash_                  atomic.Int32
}

func (this_ *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols) BindSelf_(s any) {
	_ = s
}
func (this_ *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols) IsJavaNil_() bool {
	return this_ == nil
}
func (*CrossChunkCodeMotion_OrderAndCombineGlobalSymbols) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CrossChunkCodeMotion$OrderAndCombineGlobalSymbols", "OrderAndCombineGlobalSymbols"
}
func (this_ *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CrossChunkCodeMotion_OrderAndCombineGlobalSymbols) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CrossChunkCodeMotion.OrderAndCombineGlobalSymbols", "com.google.javascript.jscomp.CrossChunkCodeMotion$OrderAndCombineGlobalSymbols", func(v_ any) bool {
		return j2gort.InstanceOf[*CrossChunkCodeMotion_OrderAndCombineGlobalSymbols](v_)
	})
}
func (this_ *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols) SetOuter_(o *CrossChunkCodeMotion) {
	this_.outer_ = o
}
func (this_ *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols) GetOuter_() *CrossChunkCodeMotion {
	return this_.outer_
}

func (this_ *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols) initFields_CrossChunkCodeMotion_OrderAndCombineGlobalSymbols() {
	_ = this_
	this_.componentContents = j2gort.Deque[*CrossChunkCodeMotion_GlobalSymbol]{}
	this_.componentRoots = j2gort.Deque[*CrossChunkCodeMotion_GlobalSymbol]{}

}

// Java constructor: OrderAndCombineGlobalSymbols(com.google.javascript.jscomp.CrossChunkCodeMotion.GlobalSymbolCollector collector)
func NewCrossChunkCodeMotion_OrderAndCombineGlobalSymbols(outer_ *CrossChunkCodeMotion, collector *CrossChunkCodeMotion_GlobalSymbolCollector) *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols {
	this_ := &CrossChunkCodeMotion_OrderAndCombineGlobalSymbols{}
	this_.outer_ = outer_
	InitCrossChunkCodeMotion_OrderAndCombineGlobalSymbols(this_, collector)
	return this_
}

func InitCrossChunkCodeMotion_OrderAndCombineGlobalSymbols(this_ *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols, collector *CrossChunkCodeMotion_GlobalSymbolCollector) {
	_ = collector
	this_.initFields_CrossChunkCodeMotion_OrderAndCombineGlobalSymbols()
	this_.collector = collector
	this_.stronglyConnectedSymbols = j2gort.NewDequeCap[*CrossChunkCodeMotion_GlobalSymbolCycle](collector.symbolStack.Size())
}

// Java method: OrderAndCombineGlobalSymbols.orderAndCombine()
func (this_ *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols) OrderAndCombine() {
	{
		cur_39 := this_.collector.symbolStack.Iter()
		for cur_39.HasNext() {

			globalSymbol := cur_39.Next()
			_ = globalSymbol
			if globalSymbol.preorderNumber < int32(0) {
				this_.ProcessGlobalSymbol(globalSymbol)
			}
		}
	}
}

// Java method: OrderAndCombineGlobalSymbols.processGlobalSymbol(com.google.javascript.jscomp.CrossChunkCodeMotion.GlobalSymbol symbol)
func (this_ *CrossChunkCodeMotion_OrderAndCombineGlobalSymbols) ProcessGlobalSymbol(symbol *CrossChunkCodeMotion_GlobalSymbol) {
	_ = symbol
	if !(symbol.preorderNumber < int32(0)) {
		j2gort.CheckStateFmt(false, "already processed: %s", j2gort.PtrToAny(symbol))
	}
	{
		r_43 := symbol
		r_43.preorderNumber = func() int32 {
			old_41 := this_.preorderCounter
			nw_42 := old_41 + 1
			this_.preorderCounter = nw_42
			return old_41
		}()
	}
	this_.componentRoots.Push(symbol)
	this_.componentContents.Push(symbol)
	{
		cur_44 := symbol.referencingGlobalSymbols.Iter()
		for cur_44.HasNext() {

			referringSymbol := cur_44.
				Next()
			_ = referringSymbol
			if referringSymbol.preorderNumber < int32(0) {
				this_.ProcessGlobalSymbol(referringSymbol)
			} else {
				if !(referringSymbol.hasBeenAssignedToAStronglyConnectedComponent) {
					for this_.componentRoots.Peek().preorderNumber > referringSymbol.preorderNumber {
						this_.componentRoots.Pop()
					}
				}
			}
		}
	}
	if j2gort.CheckNotNullPtr(this_.componentRoots.Peek()) == symbol {
		this_.componentRoots.Pop()
		val_46 := allocPooled_CrossChunkCodeMotion_GlobalSymbolCycle(&this_.outer_.slab_GlobalSymbolCycle, this_.outer_)

		var cycle *CrossChunkCodeMotion_GlobalSymbolCycle = val_46
		_ = cycle
		var connectedSymbol *CrossChunkCodeMotion_GlobalSymbol = nil
		_ = connectedSymbol
		for doFirst_47 := true; doFirst_47 || (!(j2gort.CheckNotNullPtr(connectedSymbol) == symbol)); doFirst_47 = false {
			connectedSymbol = this_.componentContents.Pop()
			cycle.AddSymbol(connectedSymbol)
			{
				r_48 := connectedSymbol
				r_48.hasBeenAssignedToAStronglyConnectedComponent = true
			}
		}
		this_.stronglyConnectedSymbols.AddLast(cycle)
	}
}

// Java class: com.google.javascript.jscomp.CrossChunkCodeMotion.DeclarationStatementGroupCycle
type CrossChunkCodeMotion_DeclarationStatementGroupCycle struct {
	outer_       *CrossChunkCodeMotion
	currentChunk *JSChunk                                                      `gson:"currentChunk"`
	dsgs         j2gort.Deque[*CrossChunkCodeMotion_DeclarationStatementGroup] `gson:"dsgs"`
	idHash_      atomic.Int32
}

func (this_ *CrossChunkCodeMotion_DeclarationStatementGroupCycle) BindSelf_(s any) {
	_ = s
}
func (this_ *CrossChunkCodeMotion_DeclarationStatementGroupCycle) IsJavaNil_() bool {
	return this_ == nil
}
func (*CrossChunkCodeMotion_DeclarationStatementGroupCycle) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CrossChunkCodeMotion$DeclarationStatementGroupCycle", "DeclarationStatementGroupCycle"
}
func (this_ *CrossChunkCodeMotion_DeclarationStatementGroupCycle) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CrossChunkCodeMotion_DeclarationStatementGroupCycle) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CrossChunkCodeMotion.DeclarationStatementGroupCycle", "com.google.javascript.jscomp.CrossChunkCodeMotion$DeclarationStatementGroupCycle", func(v_ any) bool {
		return j2gort.InstanceOf[*CrossChunkCodeMotion_DeclarationStatementGroupCycle](v_)
	})
}
func (this_ *CrossChunkCodeMotion_DeclarationStatementGroupCycle) SetOuter_(o *CrossChunkCodeMotion) {
	this_.outer_ = o
}
func (this_ *CrossChunkCodeMotion_DeclarationStatementGroupCycle) GetOuter_() *CrossChunkCodeMotion {
	return this_.outer_
}

func (this_ *CrossChunkCodeMotion_DeclarationStatementGroupCycle) initFields_CrossChunkCodeMotion_DeclarationStatementGroupCycle() {
	_ = this_
}

// Java constructor: DeclarationStatementGroupCycle(com.google.javascript.jscomp.JSChunk currentChunk)
func NewCrossChunkCodeMotion_DeclarationStatementGroupCycle(outer_ *CrossChunkCodeMotion, currentChunk *JSChunk) *CrossChunkCodeMotion_DeclarationStatementGroupCycle {
	this_ := &CrossChunkCodeMotion_DeclarationStatementGroupCycle{}
	this_.outer_ = outer_
	InitCrossChunkCodeMotion_DeclarationStatementGroupCycle(this_, currentChunk)
	return this_
}
func allocPooled_CrossChunkCodeMotion_DeclarationStatementGroupCycle(slab_ *j2gort.PooledSlab[CrossChunkCodeMotion_DeclarationStatementGroupCycle], outer_ *CrossChunkCodeMotion, currentChunk *JSChunk) *CrossChunkCodeMotion_DeclarationStatementGroupCycle {
	this_ := slab_.Alloc()
	this_.outer_ = outer_
	InitCrossChunkCodeMotion_DeclarationStatementGroupCycle(this_, currentChunk)
	return this_
}

func InitCrossChunkCodeMotion_DeclarationStatementGroupCycle(this_ *CrossChunkCodeMotion_DeclarationStatementGroupCycle, currentChunk *JSChunk) {
	_ = currentChunk
	this_.initFields_CrossChunkCodeMotion_DeclarationStatementGroupCycle()
	this_.currentChunk = currentChunk
	this_.dsgs = j2gort.Deque[*CrossChunkCodeMotion_DeclarationStatementGroup]{}
}

// Java method: DeclarationStatementGroupCycle.moveToPreferredChunk(java.util.BitSet chunksWithImmovableReferences)
func (this_ *CrossChunkCodeMotion_DeclarationStatementGroupCycle) MoveToPreferredChunk(chunksWithImmovableReferences *j2gort.Ext_java_util_BitSet) {
	_ = chunksWithImmovableReferences
	var preferredChunk *JSChunk = this_.priv_DeclarationStatementGroupCycle_getPreferredChunk_2266(chunksWithImmovableReferences)
	_ = preferredChunk
	{
		var iifeCond_1 bool
		{

			rp_49 := preferredChunk
			ap_50 := this_.currentChunk
			iifeCond_1 = j2gort.CheckNotNullPtr(rp_49) == ap_50
		}
		if !(iifeCond_1) {
			if this_.outer_.cccmLog.IsLogging() {
				this_.outer_.cccmLog.Log_1_string("start DSG Cycle move")
				{
					cur_51 := j2gort.NewIterCursor(j2gort.CheckNotNull(this_.priv_DeclarationStatementGroupCycle_getGlobalSymbolNames_2264()))
					for cur_51.HasNext() {
						globalSymbolName := cur_51.NextStr()
						_ = globalSymbolName
						this_.outer_.cccmLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
							return j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("Moving DSG for %s from chunk %s to chunk %s", []any{j2gort.BoxString(globalSymbolName), j2gort.PtrToAny(this_.currentChunk), j2gort.PtrToAny(preferredChunk)}))
						}))
					}
				}
			}
			this_.priv_DeclarationStatementGroupCycle_moveStatementsToChunk_2265(preferredChunk)
		}
	}
	{
		cur_52 := this_.dsgs.Iter()
		for cur_52.HasNext() {

			dsg := cur_52.Next()
			_ = dsg
			{
				r_54 := dsg
				r_54.currentChunk = preferredChunk
			}
			dsg.MakeReferencesImmovable()
		}
	}
}

// Java method: DeclarationStatementGroupCycle.getGlobalSymbolNames()
func (this_ *CrossChunkCodeMotion_DeclarationStatementGroupCycle) priv_DeclarationStatementGroupCycle_getGlobalSymbolNames_2264() j2gort.Ext_com_google_common_collect_ImmutableList {
	var builder []string = []string{}
	_ = builder
	{
		cur_55 := this_.dsgs.Iter()
		for cur_55.HasNext() {

			dsg := cur_55.Next()
			_ = dsg
			builder = append(builder, j2gort.CheckNotNullString(dsg.declaredGlobalSymbol.var_.GetName()))
		}
	}
	return j2gort.ImmListFromTypedSlice(builder)
}

// Java method: DeclarationStatementGroupCycle.moveStatementsToChunk(com.google.javascript.jscomp.JSChunk preferredChunk)
func (this_ *CrossChunkCodeMotion_DeclarationStatementGroupCycle) priv_DeclarationStatementGroupCycle_moveStatementsToChunk_2265(preferredChunk *JSChunk) {
	_ = preferredChunk
	var destParent pkg_rhino.Node = func() pkg_rhino.Node {
		bound_57 := this_.outer_.compiler
		return j2gort.OrderedMapComputeIfAbsent(&this_.outer_.chunkInsertionPointMap, preferredChunk, pkg_rhino.Node(0), func(mrp0_ *JSChunk) pkg_rhino.Node {
			return bound_57.GetNodeForCodeInsertion(mrp0_)
		})
	}()

	_ = destParent
	var statementsLastFirst []*CrossChunkReferenceCollector_TopLevelStatement = []*CrossChunkReferenceCollector_TopLevelStatement{}
	_ = statementsLastFirst
	{
		cur_58 := this_.dsgs.Iter()
		for cur_58.HasNext() {

			dsg := cur_58.Next()
			_ = dsg
			if len(statementsLastFirst) == 0 {
				{
					cur_60 := dsg.statementStack.Iter()
					for cur_60.HasNext() {

						stmt := cur_60.Next()
						_ = stmt
						statementsLastFirst = append(statementsLastFirst, stmt)
					}
				}
			} else {
				var list1 []*CrossChunkReferenceCollector_TopLevelStatement = []*CrossChunkReferenceCollector_TopLevelStatement{}
				_ = list1
				{
					cur_62 := dsg.statementStack.Iter()
					for cur_62.HasNext() {

						stmt_2 := cur_62.Next()
						_ = stmt_2
						list1 = append(list1, stmt_2)
					}
				}
				var prev []*CrossChunkReferenceCollector_TopLevelStatement = j2gort.CloneNonNilSlice(statementsLastFirst)
				_ = prev
				clear(statementsLastFirst)
				statementsLastFirst = statementsLastFirst[:0]
				var i1 int32 = int32(0)
				_ = i1
				var n1 int32 = int32(len(list1))
				_ = n1
				var i2 int32 = int32(0)
				_ = i2
				var n2 int32 = int32(len(prev))
				_ = n2
				for (i1 < n1) && (i2 < n2) {
					var s1 *CrossChunkReferenceCollector_TopLevelStatement = list1[i1]
					_ = s1
					var s2 *CrossChunkReferenceCollector_TopLevelStatement = prev[i2]
					_ = s2
					if s1.GetOriginalOrder() > s2.GetOriginalOrder() {
						i1 = i1 + 1
						statementsLastFirst = append(statementsLastFirst, s1)
						if !((i1 == n1) || (list1[i1].GetOriginalOrder() < s1.GetOriginalOrder())) {
							j2gort.CheckStateStr(false, "Statements are recorded in the wrong order.")
						}
					} else {
						i2 = i2 + 1
						statementsLastFirst = append(statementsLastFirst, s2)
						if !((i2 == n2) || (prev[i2].GetOriginalOrder() < s2.GetOriginalOrder())) {
							j2gort.CheckStateStr(false, "Statements are recorded in the wrong order.")
						}
					}
				}
				for i1 < n1 {
					statementsLastFirst = append(statementsLastFirst, list1[func() int32 { old_64 := i1; nw_65 := old_64 + 1; i1 = nw_65; return old_64 }()])
				}
				for i2 < n2 {
					statementsLastFirst = append(statementsLastFirst, prev[func() int32 { old_66 := i2; nw_67 := old_66 + 1; i2 = nw_67; return old_66 }()])
				}
			}
		}
	}
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = int32(len(statementsLastFirst))
		_ = size
		for ; i < size; i = i + 1 {
			var statementNode pkg_rhino.Node = statementsLastFirst[i].GetStatementNode()
			_ = statementNode
			this_.outer_.compiler.ReportChangeToEnclosingScope(statementNode)
			var originalScript pkg_rhino.Node = NodeUtil_GetEnclosingScript(statementNode)
			_ = originalScript
			statementNode.Detach()
			destParent.AddChildToFront(statementNode)
			NodeUtil_AddFeaturesToScript(destParent, NodeUtil_GetFeatureSetOfScript(originalScript), this_.outer_.compiler)
			this_.outer_.compiler.ReportChangeToEnclosingScope(statementNode)
		}
	}
}

// Java method: DeclarationStatementGroupCycle.getPreferredChunk(java.util.BitSet chunksWithImmovableReferences)
func (this_ *CrossChunkCodeMotion_DeclarationStatementGroupCycle) priv_DeclarationStatementGroupCycle_getPreferredChunk_2266(chunksWithImmovableReferences *j2gort.Ext_java_util_BitSet) *JSChunk {
	_ = chunksWithImmovableReferences
	if chunksWithImmovableReferences.IsEmpty() {
		return this_.currentChunk
	} else if !(this_.AllStatementsCanMove()) {
		return this_.currentChunk
	} else {
		return this_.outer_.graph.GetSmallestCoveringSubtree(this_.currentChunk, chunksWithImmovableReferences)
	}
}

// Java method: DeclarationStatementGroupCycle.allStatementsCanMove()
func (this_ *CrossChunkCodeMotion_DeclarationStatementGroupCycle) AllStatementsCanMove() bool {
	{
		cur_68 := this_.dsgs.Iter()
		for cur_68.HasNext() {

			dsg := cur_68.Next()
			_ = dsg
			if !(dsg.AllStatementsCanMove()) {
				return false
			}
		}
	}
	return true
}

// Java class: com.google.javascript.jscomp.CrossChunkCodeMotion.InstanceofReference
type CrossChunkCodeMotion_InstanceofReference struct {
	chunk         *JSChunk                                        `gson:"chunk"`
	containingDsg *CrossChunkCodeMotion_DeclarationStatementGroup `gson:"containingDsg"`
	reference     *Reference                                      `gson:"reference"`
	idHash_       atomic.Int32
}

func (this_ *CrossChunkCodeMotion_InstanceofReference) BindSelf_(s any) {
	_ = s
}
func (this_ *CrossChunkCodeMotion_InstanceofReference) IsJavaNil_() bool { return this_ == nil }
func (*CrossChunkCodeMotion_InstanceofReference) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.CrossChunkCodeMotion$InstanceofReference", "InstanceofReference"
}
func (this_ *CrossChunkCodeMotion_InstanceofReference) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*CrossChunkCodeMotion_InstanceofReference) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.CrossChunkCodeMotion.InstanceofReference", "com.google.javascript.jscomp.CrossChunkCodeMotion$InstanceofReference", func(v_ any) bool {
		return j2gort.InstanceOf[*CrossChunkCodeMotion_InstanceofReference](v_)
	})
}

func (this_ *CrossChunkCodeMotion_InstanceofReference) initFields_CrossChunkCodeMotion_InstanceofReference() {
	_ = this_
}

// Java constructor: InstanceofReference(com.google.javascript.jscomp.@org.jspecify.annotations.Nullable JSChunk chunk, com.google.javascript.jscomp.CrossChunkCodeMotion.@org.jspecify.annotations.Nullable DeclarationStatementGroup containingDsg, com.google.javascript.jscomp.Reference reference)
func NewCrossChunkCodeMotion_InstanceofReference(chunk *JSChunk, containingDsg *CrossChunkCodeMotion_DeclarationStatementGroup, reference *Reference) *CrossChunkCodeMotion_InstanceofReference {
	this_ := &CrossChunkCodeMotion_InstanceofReference{}
	InitCrossChunkCodeMotion_InstanceofReference(this_, chunk, containingDsg, reference)
	return this_
}
func allocPooled_CrossChunkCodeMotion_InstanceofReference(slab_ *j2gort.PooledSlab[CrossChunkCodeMotion_InstanceofReference], chunk *JSChunk, containingDsg *CrossChunkCodeMotion_DeclarationStatementGroup, reference *Reference) *CrossChunkCodeMotion_InstanceofReference {
	this_ := slab_.Alloc()
	InitCrossChunkCodeMotion_InstanceofReference(this_, chunk, containingDsg, reference)
	return this_
}

func InitCrossChunkCodeMotion_InstanceofReference(this_ *CrossChunkCodeMotion_InstanceofReference, chunk *JSChunk, containingDsg *CrossChunkCodeMotion_DeclarationStatementGroup, reference *Reference) {
	_ = chunk
	_ = containingDsg
	_ = reference
	this_.initFields_CrossChunkCodeMotion_InstanceofReference()
	this_.chunk = chunk
	this_.containingDsg = containingDsg
	this_.reference = reference
}

// Java static method: InstanceofReference.immovable(com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.jscomp.Reference reference)
func CrossChunkCodeMotion_InstanceofReference_Immovable(chunk *JSChunk, reference *Reference) *CrossChunkCodeMotion_InstanceofReference {
	_ = chunk
	_ = reference
	return NewCrossChunkCodeMotion_InstanceofReference(chunk, nil, reference)
}
func allocPooledVia_CrossChunkCodeMotion_InstanceofReference_Immovable(slab_ *j2gort.PooledSlab[CrossChunkCodeMotion_InstanceofReference], chunk *JSChunk, reference *Reference) *CrossChunkCodeMotion_InstanceofReference {
	_ = chunk
	_ = reference
	return allocPooled_CrossChunkCodeMotion_InstanceofReference(slab_, chunk, nil, reference)
}

// Java static method: InstanceofReference.movable(com.google.javascript.jscomp.CrossChunkCodeMotion.DeclarationStatementGroup containingDsg, com.google.javascript.jscomp.Reference reference)
func CrossChunkCodeMotion_InstanceofReference_Movable(containingDsg *CrossChunkCodeMotion_DeclarationStatementGroup, reference *Reference) *CrossChunkCodeMotion_InstanceofReference {
	_ = containingDsg
	_ = reference
	return NewCrossChunkCodeMotion_InstanceofReference(nil, containingDsg, reference)
}
func allocPooledVia_CrossChunkCodeMotion_InstanceofReference_Movable(slab_ *j2gort.PooledSlab[CrossChunkCodeMotion_InstanceofReference], containingDsg *CrossChunkCodeMotion_DeclarationStatementGroup, reference *Reference) *CrossChunkCodeMotion_InstanceofReference {
	_ = containingDsg
	_ = reference
	return allocPooled_CrossChunkCodeMotion_InstanceofReference(slab_, nil, containingDsg, reference)
}

// Java method: InstanceofReference.getChunk()
func (this_ *CrossChunkCodeMotion_InstanceofReference) GetChunk() *JSChunk {
	if this_.containingDsg != nil {
		return this_.containingDsg.currentChunk
	} else {
		return this_.chunk
	}
}

// Java method: InstanceofReference.getReference()
func (this_ *CrossChunkCodeMotion_InstanceofReference) GetReference() *Reference {
	return this_.reference
}
func (this_ *CrossChunkCodeMotion) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_GlobalSymbol.ReleaseWith(func(x_ *CrossChunkCodeMotion_GlobalSymbol) {
		var z_ CrossChunkCodeMotion_GlobalSymbol
		*x_ = z_
	})
	this_.slab_GlobalSymbolCycle.ReleaseWith(func(x_ *CrossChunkCodeMotion_GlobalSymbolCycle) {
		var z_ CrossChunkCodeMotion_GlobalSymbolCycle
		*x_ = z_
	})
	this_.slab_DeclarationStatementGroup.ReleaseWith(func(x_ *CrossChunkCodeMotion_DeclarationStatementGroup) {
		var z_ CrossChunkCodeMotion_DeclarationStatementGroup
		*x_ = z_
	})
	this_.slab_DeclarationStatementGroupCycle.ReleaseWith(func(x_ *CrossChunkCodeMotion_DeclarationStatementGroupCycle) {
		var z_ CrossChunkCodeMotion_DeclarationStatementGroupCycle
		*x_ = z_
	})
	this_.slab_InstanceofReference.ReleaseWith(func(x_ *CrossChunkCodeMotion_InstanceofReference) {
		var z_ CrossChunkCodeMotion_InstanceofReference
		*x_ = z_
	})
}
