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
// Port of com/google/javascript/jscomp/Es6RewriteBlockScopedDeclaration.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_colors "github.com/paullepoulpe/jsgo/transpiled/gen/colors"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_parser "github.com/paullepoulpe/jsgo/transpiled/gen/parser"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

var lambdaSingleton_es6_rewrite_block_scoped_declaration_1 j2gort.Ext_com_google_common_base_Predicate

func init() {
	lambdaSingleton_es6_rewrite_block_scoped_declaration_1 = j2gort.NewNodeLambda_Ext_com_google_common_base_Predicate(func(rawNode_ uint32) bool {
		return Es6RewriteBlockScopedDeclaration_IsLoopOrFunction(pkg_rhino.Node(rawNode_))
	})
}

// Java class: com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration
type Es6RewriteBlockScopedDeclaration struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	compiler         AbstractCompiler                  `gson:"compiler"`
	astFactory       *AstFactory                       `gson:"astFactory"`
	letConsts        j2gort.OrderedSet[pkg_rhino.Node] `gson:"letConsts"`
	uniqueIdSupplier *UniqueIdSupplier                 `gson:"uniqueIdSupplier"`
}

func (this_ *Es6RewriteBlockScopedDeclaration) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteBlockScopedDeclaration) IsJavaNil_() bool { return this_ == nil }
func (*Es6RewriteBlockScopedDeclaration) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration", "Es6RewriteBlockScopedDeclaration"
}
func (this_ *Es6RewriteBlockScopedDeclaration) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*Es6RewriteBlockScopedDeclaration) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration", "com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteBlockScopedDeclaration](v_)
	})
}
func (*Es6RewriteBlockScopedDeclaration) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*Es6RewriteBlockScopedDeclaration) Implements_jscomp_CompilerPass()           {}
func (*Es6RewriteBlockScopedDeclaration) Implements_jscomp_NodeTraversal_Callback() {}

var static_Es6RewriteBlockScopedDeclaration_transpiledFeatures *pkg_parser.FeatureSet

func StaticGet_Es6RewriteBlockScopedDeclaration_TranspiledFeatures() *pkg_parser.FeatureSet {
	clinit_Es6RewriteBlockScopedDeclaration()
	return static_Es6RewriteBlockScopedDeclaration_transpiledFeatures
}
func StaticSet_Es6RewriteBlockScopedDeclaration_TranspiledFeatures(v *pkg_parser.FeatureSet) {
	clinit_Es6RewriteBlockScopedDeclaration()
	static_Es6RewriteBlockScopedDeclaration_transpiledFeatures = v
}
func StaticGet_Es6RewriteBlockScopedDeclaration_LOOP_PARAM_NAME_PREFIX() string {
	return "$jscomp$loop_param$"
}

var clinitDone_Es6RewriteBlockScopedDeclaration uint32
var clinitState_Es6RewriteBlockScopedDeclaration j2gort.ClinitState

func clinit_Es6RewriteBlockScopedDeclaration() {
	if atomic.LoadUint32(&clinitDone_Es6RewriteBlockScopedDeclaration) != 0 {
		return
	}
	clinitSlow_Es6RewriteBlockScopedDeclaration()
}

func clinitSlow_Es6RewriteBlockScopedDeclaration() {
	if !j2gort.EnterClinit(&clinitState_Es6RewriteBlockScopedDeclaration, &clinitDone_Es6RewriteBlockScopedDeclaration) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_Es6RewriteBlockScopedDeclaration, &clinitDone_Es6RewriteBlockScopedDeclaration)
	static_Es6RewriteBlockScopedDeclaration_transpiledFeatures = pkg_parser.StaticGet_FeatureSet_BARE_MINIMUM().With_1_FeatureSet_Feature_7192([]*pkg_parser.FeatureSet_Feature{pkg_parser.StaticGet_FeatureSet_Feature_LET_DECLARATIONS(), pkg_parser.StaticGet_FeatureSet_Feature_CONST_DECLARATIONS()})
}

func (this_ *Es6RewriteBlockScopedDeclaration) initFields_Es6RewriteBlockScopedDeclaration() {
	_ = this_
	this_.letConsts = j2gort.NewOrderedSet[pkg_rhino.Node]()
}

// Java constructor: Es6RewriteBlockScopedDeclaration(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewEs6RewriteBlockScopedDeclaration(compiler AbstractCompiler) *Es6RewriteBlockScopedDeclaration {
	clinit_Es6RewriteBlockScopedDeclaration()
	this_ := &Es6RewriteBlockScopedDeclaration{}
	InitEs6RewriteBlockScopedDeclaration(this_, compiler)
	return this_
}

func InitEs6RewriteBlockScopedDeclaration(this_ *Es6RewriteBlockScopedDeclaration, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_Es6RewriteBlockScopedDeclaration()
	this_.compiler = compiler
	this_.uniqueIdSupplier = compiler.GetUniqueIdSupplier()
	this_.astFactory = compiler.CreateAstFactory()
}

// Java method: Es6RewriteBlockScopedDeclaration.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteBlockScopedDeclaration) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !((parent == 0) || !(parent.IsForOf())) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(parent))
	}
	if n.IsLet() || n.IsConst() {
		this_.letConsts.Add(n)
		{
			var nameNode pkg_rhino.Node = n.GetFirstChild()
			_ = nameNode
			for ; nameNode != 0; nameNode = nameNode.GetNext() {
				this_.priv_Es6RewriteBlockScopedDeclaration_visitBlockScopedNameDeclaration_2564(n, nameNode)
			}
		}
	}
}

// Java method: Es6RewriteBlockScopedDeclaration.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Es6RewriteBlockScopedDeclaration) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
	var transformer *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer = NewEs6RewriteBlockScopedDeclaration_LoopClosureTransformer(this_)
	_ = transformer
	sc_ := NewSyntacticScopeCreator_1(this_.compiler)
	defer sc_.ReleaseSlabs_()
	{
		clinit_NodeTraversal()
		var ntb_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_)
		(&ntb_).SetCompiler(this_.compiler)
		(&ntb_).SetCallback_1_NodeTraversal_Callback(transformer)
		(&ntb_).SetScopeCreator(sc_)
		(&ntb_).Traverse(root)
	}

	transformer.priv_LoopClosureTransformer_transformLoopClosure_2574()
	this_.priv_Es6RewriteBlockScopedDeclaration_rewriteDeclsToVars_2571()
	TranspilationPasses_MaybeMarkFeaturesAsTranspiledAway(this_.compiler, root, StaticGet_Es6RewriteBlockScopedDeclaration_TranspiledFeatures())
}

// Java method: Es6RewriteBlockScopedDeclaration.visitBlockScopedNameDeclaration(com.google.javascript.rhino.Node decl, com.google.javascript.rhino.Node nameNode)
func (this_ *Es6RewriteBlockScopedDeclaration) priv_Es6RewriteBlockScopedDeclaration_visitBlockScopedNameDeclaration_2564(decl pkg_rhino.Node, nameNode pkg_rhino.Node) {
	_ = decl
	_ = nameNode
	var parent pkg_rhino.Node = decl.GetParent()
	_ = parent
	if (!(nameNode.HasChildren()) && ((parent == 0) || !(parent.IsForIn()))) && this_.priv_Es6RewriteBlockScopedDeclaration_inLoop_2565(decl) {
		var undefined pkg_rhino.Node = this_.astFactory.CreateUndefinedValue().SrcrefTree(nameNode)
		_ = undefined
		nameNode.AddChildToFront(undefined)
		this_.compiler.ReportChangeToEnclosingScope(undefined)
	}
}

// Java method: Es6RewriteBlockScopedDeclaration.inLoop(com.google.javascript.rhino.Node n)
func (this_ *Es6RewriteBlockScopedDeclaration) priv_Es6RewriteBlockScopedDeclaration_inLoop_2565(n pkg_rhino.Node) bool {
	_ = n
	var enclosingNode pkg_rhino.Node = NodeUtil_GetEnclosingNode(n, lambdaSingleton_es6_rewrite_block_scoped_declaration_1)
	_ = enclosingNode
	return ((enclosingNode != 0) && !(enclosingNode.IsFunction()))
}

// Java static method: Es6RewriteBlockScopedDeclaration.isLoopOrFunction(com.google.javascript.rhino.Node n)
func Es6RewriteBlockScopedDeclaration_IsLoopOrFunction(n pkg_rhino.Node) bool {
	clinit_Es6RewriteBlockScopedDeclaration()
	_ = n
	return (n.IsFunction() || NodeUtil_IsLoopStructure(n))
}

// Java static method: Es6RewriteBlockScopedDeclaration.extractInlineJSDoc(com.google.javascript.rhino.Node srcDeclaration, com.google.javascript.rhino.Node srcName, com.google.javascript.rhino.Node destDeclaration)
func Es6RewriteBlockScopedDeclaration_ExtractInlineJSDoc(srcDeclaration pkg_rhino.Node, srcName pkg_rhino.Node, destDeclaration pkg_rhino.Node) {
	clinit_Es6RewriteBlockScopedDeclaration()
	_ = srcDeclaration
	_ = srcName
	_ = destDeclaration
	var existingInfo *pkg_rhino.JSDocInfo = srcDeclaration.GetJSDocInfo()
	_ = existingInfo
	if existingInfo == nil {
		existingInfo = srcName.GetJSDocInfo()
		srcName.SetJSDocInfo(nil)
	}
	var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(existingInfo)
	_ = builder
	destDeclaration.SetJSDocInfo(builder.Build_0())
}

// Java static method: Es6RewriteBlockScopedDeclaration.maybeAddConstJSDoc(com.google.javascript.rhino.Node srcDeclaration, com.google.javascript.rhino.Node srcName, com.google.javascript.rhino.Node destDeclaration)
func Es6RewriteBlockScopedDeclaration_MaybeAddConstJSDoc(srcDeclaration pkg_rhino.Node, srcName pkg_rhino.Node, destDeclaration pkg_rhino.Node) {
	clinit_Es6RewriteBlockScopedDeclaration()
	_ = srcDeclaration
	_ = srcName
	_ = destDeclaration
	if srcDeclaration.IsConst() {
		Es6RewriteBlockScopedDeclaration_ExtractInlineJSDoc(srcDeclaration, srcName, destDeclaration)
		var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_Builder_MaybeCopyFrom(destDeclaration.GetJSDocInfo())
		_ = builder
		builder.RecordConstancy()
		destDeclaration.SetJSDocInfo(builder.Build_0())
	}
}

// Java method: Es6RewriteBlockScopedDeclaration.handleLetConstDeclarationList(com.google.javascript.rhino.Node declarationList, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteBlockScopedDeclaration) priv_Es6RewriteBlockScopedDeclaration_handleLetConstDeclarationList_2566(declarationList pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = declarationList
	_ = parent
	if declarationList.IsVar() {
		if !(!(declarationList.HasChildren()) || declarationList.HasOneChild()) {
			j2gort.CheckStateStr(false, "var declaration can have maximum one child post normalization.")
		}
		return
	}
	if parent.IsVanillaFor() {
		this_.priv_Es6RewriteBlockScopedDeclaration_handleLetConstDeclarationListInVanillaForInitializer_2568(declarationList, parent)
		return
	}
	if parent.IsForIn() {
		this_.priv_Es6RewriteBlockScopedDeclaration_handleLetConstDeclarationListInForInInitializer_2567(declarationList, parent)
		return
	}
	for declarationList.HasChildren() {
		var name pkg_rhino.Node = declarationList.GetLastChild()
		_ = name
		var newDeclaration pkg_rhino.Node = pkg_rhino.IR_Var_1(name.Detach()).Srcref(declarationList)
		_ = newDeclaration
		Es6RewriteBlockScopedDeclaration_ExtractInlineJSDoc(declarationList, name, newDeclaration)
		Es6RewriteBlockScopedDeclaration_MaybeAddConstJSDoc(declarationList, name, newDeclaration)
		newDeclaration.InsertAfter(declarationList)
		this_.compiler.ReportChangeToEnclosingScope(parent)
	}
	declarationList.Detach()
	this_.compiler.ReportChangeToEnclosingScope(parent)
}

// Java method: Es6RewriteBlockScopedDeclaration.handleLetConstDeclarationListInForInInitializer(com.google.javascript.rhino.Node declarationList, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteBlockScopedDeclaration) priv_Es6RewriteBlockScopedDeclaration_handleLetConstDeclarationListInForInInitializer_2567(declarationList pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = declarationList
	_ = parent
	if !(parent.IsForIn()) {
		j2gort.CheckState(false)
	}
	var first pkg_rhino.Node = declarationList
	_ = first
	var lhs pkg_rhino.Node = first.GetFirstChild()
	_ = lhs
	var loopNode pkg_rhino.Node = parent
	_ = loopNode
	if lhs.IsDestructuringLhs() {
		if !(false) {
			j2gort.CheckStateStr(false, "Destructuring syntax is unsupported in ES6RewriteBlockScopedDeclarations pass")
		}
		NodeUtil_VisitLhsNodesInNode(lhs, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_1 any) {
			name := j2gort.CastNode[pkg_rhino.Node](lp_1)
			_ = name
			if !(name.IsName()) {
				j2gort.CheckStateFmt(false, "lhs in destructuring declaration should be a simple name. (%s)", j2gort.BoxNode(name))
			}
			var newName pkg_rhino.Node = pkg_rhino.IR_Name(name.GetString()).Srcref(name)
			_ = newName
			var newVar pkg_rhino.Node = pkg_rhino.IR_Var_1(newName).Srcref(name)
			_ = newVar
			Es6RewriteBlockScopedDeclaration_ExtractInlineJSDoc(declarationList, name, newVar)
			this_.priv_Es6RewriteBlockScopedDeclaration_addNodeBeforeLoop_2569(newVar, loopNode)
		}))
		var destructuringPattern pkg_rhino.Node = lhs.RemoveFirstChild()
		_ = destructuringPattern
		first.ReplaceWith(destructuringPattern)
	} else {
		var newStatement pkg_rhino.Node = first.CloneTree_0()
		_ = newStatement
		var name pkg_rhino.Node = newStatement.GetFirstChild().CloneNode_0()
		_ = name
		if name.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME()) {
		}
		newStatement.SetToken(pkg_rhino.StaticGet_Token_VAR())
		Es6RewriteBlockScopedDeclaration_ExtractInlineJSDoc(first, first.GetFirstChild(), newStatement)
		first.ReplaceWith(name)
		this_.priv_Es6RewriteBlockScopedDeclaration_addNodeBeforeLoop_2569(newStatement, loopNode)
	}
}

// Java method: Es6RewriteBlockScopedDeclaration.handleLetConstDeclarationListInVanillaForInitializer(com.google.javascript.rhino.Node declarationList, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteBlockScopedDeclaration) priv_Es6RewriteBlockScopedDeclaration_handleLetConstDeclarationListInVanillaForInitializer_2568(declarationList pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = declarationList
	_ = parent
	if !(parent.IsVanillaFor()) {
		j2gort.CheckState(false)
	}
	var insertSpot pkg_rhino.Node = this_.priv_Es6RewriteBlockScopedDeclaration_getInsertSpotBeforeLoop_2570(parent)
	_ = insertSpot
	for declarationList.HasChildren() {
		var name pkg_rhino.Node = declarationList.GetLastChild()
		_ = name
		var newDeclaration pkg_rhino.Node = pkg_rhino.IR_Var_1(name.Detach()).Srcref(declarationList)
		_ = newDeclaration
		Es6RewriteBlockScopedDeclaration_ExtractInlineJSDoc(declarationList, name, newDeclaration)
		if name.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME()) {
		}
		newDeclaration.InsertBefore(insertSpot)
		insertSpot = newDeclaration
		this_.compiler.ReportChangeToEnclosingScope(parent)
	}
	var empty pkg_rhino.Node = this_.astFactory.CreateEmpty().Srcref(declarationList)
	_ = empty
	declarationList.ReplaceWith(empty)
	this_.compiler.ReportChangeToEnclosingScope(empty)
}

// Java method: Es6RewriteBlockScopedDeclaration.addNodeBeforeLoop(com.google.javascript.rhino.Node newNode, com.google.javascript.rhino.Node loopNode)
func (this_ *Es6RewriteBlockScopedDeclaration) priv_Es6RewriteBlockScopedDeclaration_addNodeBeforeLoop_2569(newNode pkg_rhino.Node, loopNode pkg_rhino.Node) {
	_ = newNode
	_ = loopNode
	var insertSpot pkg_rhino.Node = this_.priv_Es6RewriteBlockScopedDeclaration_getInsertSpotBeforeLoop_2570(loopNode)
	_ = insertSpot
	newNode.InsertBefore(insertSpot)
	this_.compiler.ReportChangeToEnclosingScope(newNode)
}

// Java method: Es6RewriteBlockScopedDeclaration.getInsertSpotBeforeLoop(com.google.javascript.rhino.Node loopNode)
func (this_ *Es6RewriteBlockScopedDeclaration) priv_Es6RewriteBlockScopedDeclaration_getInsertSpotBeforeLoop_2570(loopNode pkg_rhino.Node) pkg_rhino.Node {
	_ = loopNode
	var insertSpot pkg_rhino.Node = loopNode
	_ = insertSpot
	for insertSpot.GetParent().IsLabel() {
		insertSpot = insertSpot.GetParent()
	}
	return insertSpot
}

// Java method: Es6RewriteBlockScopedDeclaration.rewriteDeclsToVars()
func (this_ *Es6RewriteBlockScopedDeclaration) priv_Es6RewriteBlockScopedDeclaration_rewriteDeclsToVars_2571() {
	if !(this_.letConsts.IsEmpty()) {
		{
			cur_2 := this_.letConsts.Iter()
			for cur_2.HasNext() {
				n := cur_2.
					Next()
				_ = n
				this_.priv_Es6RewriteBlockScopedDeclaration_handleLetConstDeclarationList_2566(n, n.GetParent())
			}
		}
	}
}

// Java class: com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration.LoopClosureTransformer
type Es6RewriteBlockScopedDeclaration_LoopClosureTransformer struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	outer_                              *Es6RewriteBlockScopedDeclaration
	loopObjectMap                       j2gort.OrderedMap[pkg_rhino.Node, *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject] `gson:"loopObjectMap"`
	nodesRequiringLoopObjectsClosureMap j2gort.Ext_com_google_common_collect_SetMultimap                                                       `gson:"nodesRequiringLoopObjectsClosureMap"`
	nodesHandledForLoopObjectClosure    j2gort.Ext_com_google_common_collect_SetMultimap                                                       `gson:"nodesHandledForLoopObjectClosure"`
	referenceMap                        j2gort.Ext_com_google_common_collect_SetMultimap                                                       `gson:"referenceMap"`
}

func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) IsJavaNil_() bool {
	return this_ == nil
}
func (*Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration$LoopClosureTransformer", "LoopClosureTransformer"
}
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPreOrderCallbackBase.IdentityHash_()
}
func (*Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration.LoopClosureTransformer", "com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration$LoopClosureTransformer", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteBlockScopedDeclaration_LoopClosureTransformer](v_)
	})
}
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) SetOuter_(o *Es6RewriteBlockScopedDeclaration) {
	this_.outer_ = o
}
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) GetOuter_() *Es6RewriteBlockScopedDeclaration {
	return this_.outer_
}
func (*Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {
}
func (*Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) Implements_jscomp_NodeTraversal_Callback() {
}

func StaticGet_Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LOOP_OBJECT_NAME() string {
	return "$jscomp$loop"
}
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) initFields_Es6RewriteBlockScopedDeclaration_LoopClosureTransformer() {
	_ = this_
	this_.loopObjectMap = j2gort.NewOrderedMap[pkg_rhino.Node, *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject]()
	this_.nodesRequiringLoopObjectsClosureMap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_SetMultimap](j2gort.Ext_com_google_common_collect_LinkedHashMultimap_Create_0())
	this_.nodesHandledForLoopObjectClosure = j2gort.AsIface[j2gort.Ext_com_google_common_collect_SetMultimap](j2gort.Ext_com_google_common_collect_HashMultimap_Create_0())
	this_.referenceMap = j2gort.AsIface[j2gort.Ext_com_google_common_collect_SetMultimap](j2gort.Ext_com_google_common_collect_LinkedHashMultimap_Create_0())
}

// Java constructor: LoopClosureTransformer()
func NewEs6RewriteBlockScopedDeclaration_LoopClosureTransformer(outer_ *Es6RewriteBlockScopedDeclaration) *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer {
	this_ := &Es6RewriteBlockScopedDeclaration_LoopClosureTransformer{}
	this_.outer_ = outer_
	InitEs6RewriteBlockScopedDeclaration_LoopClosureTransformer(this_)
	return this_
}

func InitEs6RewriteBlockScopedDeclaration_LoopClosureTransformer(this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) {
	InitNodeTraversal_AbstractPreOrderCallback(&this_.NodeTraversal_AbstractPreOrderCallbackBase)
	this_.initFields_Es6RewriteBlockScopedDeclaration_LoopClosureTransformer()
}

// Java method: LoopClosureTransformer.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if !(NodeUtil_IsReferenceName(n)) {
		return true
	}
	var name string = n.GetString()
	_ = name
	var referencedIn *Scope = t.GetScope()
	_ = referencedIn
	var var_ *Var = j2gort.CastPtr[Var](referencedIn.GetVar(name))
	_ = var_
	if var_ == nil {
		return true
	}
	if !(var_.IsLet()) && !(var_.IsConst()) {
		return true
	}
	var declaredIn *Scope = j2gort.CastPtr[Scope](var_.GetScope())
	_ = declaredIn
	var loopNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = loopNode
	{
		var s *Scope = declaredIn
		_ = s
		for ; ; s = j2gort.CastPtr[Scope](s.GetParent()) {
			var scopeRoot pkg_rhino.Node = s.GetRootNode_0()
			_ = scopeRoot
			if NodeUtil_IsLoopStructure(scopeRoot) {
				loopNode = scopeRoot
				break
			} else if scopeRoot.HasParent() && NodeUtil_IsLoopStructure(scopeRoot.GetParent()) {
				loopNode = scopeRoot.GetParent()
				break
			} else if s.IsFunctionBlockScope() || s.IsGlobal() {
				return true
			}
		}
	}
	this_.referenceMap.Put(j2gort.PtrToAny(var_), j2gort.BoxNode(n))
	var outerMostFunctionScope *Scope = nil
	_ = outerMostFunctionScope
	{
		var s_2 *Scope = referencedIn
		_ = s_2
		for ; (s_2 != declaredIn) && (s_2.GetRootNode_0() != loopNode); s_2 = j2gort.CastPtr[Scope](s_2.GetParent()) {
			if s_2.IsFunctionScope() {
				outerMostFunctionScope = s_2
			}
		}
	}
	if outerMostFunctionScope != nil {
		var enclosingFunction pkg_rhino.Node = outerMostFunctionScope.GetRootNode_0()
		_ = enclosingFunction
		var nodeToWrapInClosure pkg_rhino.Node = pkg_rhino.Node(0)
		_ = nodeToWrapInClosure
		if enclosingFunction.GetParent().IsGetterDef() || enclosingFunction.GetParent().IsSetterDef() {
			nodeToWrapInClosure = enclosingFunction.GetGrandparent()
			if !(nodeToWrapInClosure.IsObjectLit()) {
				j2gort.CheckState(false)
			}
		} else {
			nodeToWrapInClosure = enclosingFunction
		}
		if j2gort.MultimapContainsEntryNode(this_.nodesHandledForLoopObjectClosure, nodeToWrapInClosure, j2gort.BoxString(name)) {
			return true
		}
		j2gort.MultimapPutNode(this_.nodesHandledForLoopObjectClosure, nodeToWrapInClosure, j2gort.BoxString(name))
		var object *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject = j2gort.OrderedMapComputeIfAbsent(&this_.loopObjectMap, loopNode, nil, func(lp_7 pkg_rhino.Node) *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject {
			k := lp_7
			_ = k
			return NewEs6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject(this_.priv_LoopClosureTransformer_createUniqueObjectName_2572(t.GetInput()))
		})

		_ = object
		object.vars.Add(var_)
		j2gort.MultimapPutNode(this_.nodesRequiringLoopObjectsClosureMap, nodeToWrapInClosure, j2gort.PtrToAny(object))
	}
	return true
}

// Java method: LoopClosureTransformer.createUniqueObjectName(com.google.javascript.jscomp.CompilerInput input)
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) priv_LoopClosureTransformer_createUniqueObjectName_2572(input *CompilerInput) string {
	_ = input
	return ("$jscomp$loop" + "$" + j2gort.StrValString(this_.outer_.uniqueIdSupplier.GetUniqueId(input)))
}

// Java method: LoopClosureTransformer.getLoopObjPropName(com.google.javascript.jscomp.Var var)
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) priv_LoopClosureTransformer_getLoopObjPropName_2573(var_ *Var) string {
	_ = var_
	return var_.GetNameNode().GetString()
}

// Java method: LoopClosureTransformer.transformLoopClosure()
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) priv_LoopClosureTransformer_transformLoopClosure_2574() {
	if this_.loopObjectMap.IsEmpty() {
		return
	}
	var wrapperFunctions j2gort.OrderedSet[pkg_rhino.Node] = this_.priv_LoopClosureTransformer_createWrapperFunctions_2576()
	_ = wrapperFunctions
	{
		cur_8 := this_.loopObjectMap.IterKeys()
		for cur_8.HasNext() {
			loopNode := cur_8.
				Next()
			_ = loopNode
			var loopObject *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject = this_.loopObjectMap.Get(loopNode)
			_ = loopObject
			var objectLitNextIteration pkg_rhino.Node = this_.outer_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{})
			_ = objectLitNextIteration
			this_.priv_LoopClosureTransformer_renameVarsToProperties_2583(loopObject, objectLitNextIteration, loopNode)
			var updateLoopObject pkg_rhino.Node = this_.outer_.astFactory.CreateAssign_2_Node_Node(this_.priv_LoopClosureTransformer_createLoopObjectNameNode_2586(loopObject), objectLitNextIteration)
			_ = updateLoopObject
			var objectLit pkg_rhino.Node = pkg_rhino.IR_Var_2(this_.priv_LoopClosureTransformer_createLoopObjectNameNode_2586(loopObject), this_.outer_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{})).SrcrefTree(loopNode)
			_ = objectLit
			this_.outer_.priv_Es6RewriteBlockScopedDeclaration_addNodeBeforeLoop_2569(objectLit, loopNode)
			if loopNode.IsVanillaFor() {
				this_.priv_LoopClosureTransformer_changeVanillaForLoopHeader_2578(loopNode, updateLoopObject)
			} else {
				var loopBody pkg_rhino.Node = NodeUtil_GetLoopCodeBlock(loopNode)
				_ = loopBody
				loopBody.AddChildToFront(pkg_rhino.IR_ExprResult(updateLoopObject).SrcrefTreeIfMissing(loopNode))
			}
			this_.outer_.compiler.ReportChangeToEnclosingScope(loopNode)
			this_.priv_LoopClosureTransformer_changeLoopLocalVariablesToProperties_2579(loopNode, loopObject)
		}
	}
	this_.priv_LoopClosureTransformer_updateNamesInWrapperFunctions_2575(&wrapperFunctions)
}

// Java method: LoopClosureTransformer.updateNamesInWrapperFunctions(java.util.Set<com.google.javascript.rhino.Node> wrapperFunctions)
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) priv_LoopClosureTransformer_updateNamesInWrapperFunctions_2575(wrapperFunctions *j2gort.OrderedSet[pkg_rhino.Node]) {
	_ = wrapperFunctions
	{
		cur_9 := wrapperFunctions.Iter()
		for cur_9.HasNext() {
			func_ := cur_9.
				Next()
			_ = func_
			var paramList pkg_rhino.Node = func_.GetSecondChild()
			_ = paramList
			{
				var param pkg_rhino.Node = paramList.GetFirstChild()
				_ = param
				for ; param != 0; param = param.GetNext() {
					var loopObjectName string = j2gort.JStrSub(param.GetString(), int32(0), j2gort.JStrIndexOf(param.GetString(), "$jscomp$loop_param$"))
					_ = loopObjectName
					this_.priv_LoopClosureTransformer_updateNames_2577(func_.GetLastChild(), param, loopObjectName)
				}
			}
		}
	}
}

// Java method: LoopClosureTransformer.createWrapperFunctions()
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) priv_LoopClosureTransformer_createWrapperFunctions_2576() j2gort.OrderedSet[pkg_rhino.Node] {
	var wrapperFunctions j2gort.OrderedSet[pkg_rhino.Node] = j2gort.NewOrderedSet[pkg_rhino.Node]()
	_ = wrapperFunctions
	{
		cur_10 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.nodesRequiringLoopObjectsClosureMap.KeySet())))
		for cur_10.HasNext() {
			functionOrObjectLit := pkg_rhino.Node(cur_10.NextNode())
			_ = functionOrObjectLit
			if functionOrObjectLit.IsObjectLit() {
				var enclosingFunction pkg_rhino.Node = NodeUtil_GetEnclosingFunction(functionOrObjectLit)
				_ = enclosingFunction
				var scopeBody pkg_rhino.Node
				_ = scopeBody
				if enclosingFunction != 0 {
					scopeBody = enclosingFunction.GetLastChild()
				} else {
					scopeBody = NodeUtil_GetEnclosingScript(functionOrObjectLit)
				}
				var isConstructor bool = ((enclosingFunction != 0) && NodeUtil_IsEs6Constructor(enclosingFunction))
				_ = isConstructor
				var script pkg_rhino.Node = NodeUtil_GetEnclosingScript(functionOrObjectLit)
				_ = script
				var input *CompilerInput = this_.outer_.compiler.GetInput(script.GetInputId())
				_ = input
				var thisContext *ThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext = NewThisAndArgumentsReferenceUpdater_ThisAndArgumentsContext(scopeBody, isConstructor, this_.outer_.uniqueIdSupplier.GetUniqueId(input))
				_ = thisContext
				var updater *ThisAndArgumentsReferenceUpdater = NewThisAndArgumentsReferenceUpdater(this_.outer_.compiler, thisContext, this_.outer_.astFactory)
				_ = updater
				NodeTraversal_Traverse(this_.outer_.compiler, functionOrObjectLit, updater)
				thisContext.AddVarDeclarations(this_.outer_.compiler, this_.outer_.astFactory, script)
				{
					var child pkg_rhino.Node = scopeBody.GetFirstChild()
					_ = child
					for ; child != 0; child = child.GetNext() {
						if child.IsLet() || child.IsConst() {
							this_.outer_.letConsts.Add(child)
						}
					}
				}
			}
			var returnNode pkg_rhino.Node = pkg_rhino.IR_ReturnNode_0()
			_ = returnNode
			var objects j2gort.Ext_java_util_Set = j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.MultimapGetNode(this_.nodesRequiringLoopObjectsClosureMap, functionOrObjectLit))
			_ = objects
			var parameterNames []pkg_rhino.Node = make([]pkg_rhino.Node, objects.Size())
			_ = parameterNames
			var objectNamesForCall []pkg_rhino.Node = make([]pkg_rhino.Node, objects.Size())
			_ = objectNamesForCall
			var loopObjectIndex int32 = int32(0)
			_ = loopObjectIndex
			{
				cur_11 := j2gort.NewIterCursor(j2gort.CheckNotNull(objects))
				for cur_11.HasNext() {
					rawElem_12 := cur_11.Next()
					object := j2gort.CastPtr[Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject](rawElem_12)
					_ = object
					{
						a_13 := parameterNames
						i_14 := loopObjectIndex
						a_13[i_14] = this_.priv_LoopClosureTransformer_createUniqueParameterNameToUseWithinLoop_2585(functionOrObjectLit, object)
					}
					{
						a_15 := objectNamesForCall
						i_16 := loopObjectIndex
						a_15[i_16] = this_.priv_LoopClosureTransformer_createLoopObjectNameNode_2586(object)
					}
					loopObjectIndex = loopObjectIndex + 1
				}
			}
			var iife pkg_rhino.Node = this_.outer_.astFactory.CreateFunction("", pkg_rhino.IR_ParamList(parameterNames), pkg_rhino.IR_Block_1_Node(returnNode), AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_TOP_OBJECT()))
			_ = iife
			this_.outer_.compiler.ReportChangeToChangeScope(iife)
			var call pkg_rhino.Node = this_.outer_.astFactory.CreateCall(iife, AstFactory_S_Type_1_Node(functionOrObjectLit), objectNamesForCall)
			_ = call
			call.PutBooleanProp(pkg_rhino.StaticGet_Node_FREE_CALL(), true)
			var replacement pkg_rhino.Node = pkg_rhino.Node(0)
			_ = replacement
			if !(!(NodeUtil_IsFunctionDeclaration(functionOrObjectLit))) {
				j2gort.CheckStateFmt(false, "block-scoped function declarations should not exist now: %s", j2gort.BoxNode(functionOrObjectLit))
			}
			replacement = call.SrcrefTreeIfMissing(functionOrObjectLit)
			functionOrObjectLit.ReplaceWith(replacement)
			wrapperFunctions.Add(iife)
			returnNode.AddChildToFront(functionOrObjectLit)
			this_.outer_.compiler.ReportChangeToEnclosingScope(replacement)
		}
	}
	return wrapperFunctions
}

// Java method: LoopClosureTransformer.updateNames(com.google.javascript.rhino.Node node, com.google.javascript.rhino.Node paramName, java.lang.String loopObjectName)
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) priv_LoopClosureTransformer_updateNames_2577(node pkg_rhino.Node, paramName pkg_rhino.Node, loopObjectName string) {
	_ = node
	_ = paramName
	_ = loopObjectName
	if node.IsName() && !(j2gort.CheckNotNullString(node.GetString()) == "") {
		if j2gort.CheckNotNullString(node.GetString()) == loopObjectName {
			node.SetString(paramName.GetString())
			return
		}
	}
	{
		var child pkg_rhino.Node = node.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			this_.priv_LoopClosureTransformer_updateNames_2577(child, paramName, loopObjectName)
		}
	}
}

// Java method: LoopClosureTransformer.changeVanillaForLoopHeader(com.google.javascript.rhino.Node loopNode, com.google.javascript.rhino.Node updateLoopObject)
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) priv_LoopClosureTransformer_changeVanillaForLoopHeader_2578(loopNode pkg_rhino.Node, updateLoopObject pkg_rhino.Node) {
	_ = loopNode
	_ = updateLoopObject
	var initializer pkg_rhino.Node = loopNode.GetFirstChild()
	_ = initializer
	initializer.ReplaceWith(pkg_rhino.IR_Empty())
	if !(initializer.IsEmpty()) {
		if !(NodeUtil_IsNameDeclaration(initializer)) {
			initializer = pkg_rhino.IR_ExprResult(initializer).Srcref(initializer)
		}
		this_.outer_.priv_Es6RewriteBlockScopedDeclaration_addNodeBeforeLoop_2569(initializer, loopNode)
	}
	var increment pkg_rhino.Node = loopNode.GetChildAtIndex(int32(2))
	_ = increment
	if increment.IsEmpty() {
		increment.ReplaceWith(updateLoopObject.SrcrefTreeIfMissing(loopNode))
	} else {
		var placeHolder pkg_rhino.Node = pkg_rhino.IR_Empty()
		_ = placeHolder
		increment.ReplaceWith(placeHolder)
		placeHolder.ReplaceWith(this_.outer_.astFactory.CreateComma(updateLoopObject, increment).SrcrefTreeIfMissing(loopNode))
	}
}

// Java method: LoopClosureTransformer.changeLoopLocalVariablesToProperties(com.google.javascript.rhino.Node loopNode, com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.LoopObject loopObject)
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) priv_LoopClosureTransformer_changeLoopLocalVariablesToProperties_2579(loopNode pkg_rhino.Node, loopObject *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject) {
	_ = loopNode
	_ = loopObject
	{
		cur_17 := loopObject.vars.Iter()
		for cur_17.HasNext() {

			var_ := cur_17.
				Next()
			_ = var_
			var newPropertyName string = this_.priv_LoopClosureTransformer_getLoopObjPropName_2573(var_)
			_ = newPropertyName
			{
				cur_19 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.referenceMap.Get(j2gort.PtrToAny(var_)))))
				for cur_19.HasNext() {
					reference := pkg_rhino.Node(cur_19.NextNode())
					_ = reference
					if !(!(loopNode.IsForOf())) {
						j2gort.CheckStateMsg(false, j2gort.BoxNode(loopNode))
					}
					if loopNode.IsForIn() && (loopNode.GetFirstChild() == reference.GetParent()) {
						this_.priv_LoopClosureTransformer_assignLoopVarToLoopObjectProperty_2582(loopNode, loopObject, var_, newPropertyName, reference)
					} else {
						if NodeUtil_IsNameDeclaration(reference.GetParent()) {
							this_.priv_LoopClosureTransformer_replaceDeclarationWithProperty_2580(loopObject, newPropertyName, reference)
						} else {
							this_.priv_LoopClosureTransformer_replaceReferenceWithProperty_2581(loopObject, newPropertyName, reference)
						}
					}
				}
			}
		}
	}
}

// Java method: LoopClosureTransformer.replaceDeclarationWithProperty(com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.LoopObject loopObject, java.lang.String newPropertyName, com.google.javascript.rhino.Node reference)
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) priv_LoopClosureTransformer_replaceDeclarationWithProperty_2580(loopObject *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject, newPropertyName string, reference pkg_rhino.Node) {
	_ = loopObject
	_ = newPropertyName
	_ = reference
	var declaration pkg_rhino.Node = reference.GetParent()
	_ = declaration
	var grandParent pkg_rhino.Node = declaration.GetParent()
	_ = grandParent
	this_.outer_.letConsts.Remove(declaration)
	this_.outer_.priv_Es6RewriteBlockScopedDeclaration_handleLetConstDeclarationList_2566(declaration, grandParent)
	declaration = reference.GetParent()
	if reference.HasChildren() {
		var newReference pkg_rhino.Node = this_.priv_LoopClosureTransformer_createLoopVarReferenceReplacement_2584(loopObject, reference, newPropertyName)
		_ = newReference
		var assign pkg_rhino.Node = this_.outer_.astFactory.CreateAssign_2_Node_Node(newReference, reference.RemoveFirstChild())
		_ = assign
		Es6RewriteBlockScopedDeclaration_ExtractInlineJSDoc(declaration, reference, declaration)
		Es6RewriteBlockScopedDeclaration_MaybeAddConstJSDoc(declaration, reference, declaration)
		assign.SetJSDocInfo(declaration.GetJSDocInfo())
		var replacement pkg_rhino.Node = pkg_rhino.IR_ExprResult(assign).SrcrefTreeIfMissing(declaration)
		_ = replacement
		declaration.ReplaceWith(replacement)
	} else {
		declaration.Detach()
	}
	this_.outer_.compiler.ReportChangeToEnclosingScope(grandParent)
}

// Java method: LoopClosureTransformer.replaceReferenceWithProperty(com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.LoopObject loopObject, java.lang.String newPropertyName, com.google.javascript.rhino.Node reference)
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) priv_LoopClosureTransformer_replaceReferenceWithProperty_2581(loopObject *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject, newPropertyName string, reference pkg_rhino.Node) {
	_ = loopObject
	_ = newPropertyName
	_ = reference
	var referenceParent pkg_rhino.Node = reference.GetParent()
	_ = referenceParent
	reference.ReplaceWith(this_.priv_LoopClosureTransformer_createLoopVarReferenceReplacement_2584(loopObject, reference, newPropertyName))
	this_.outer_.compiler.ReportChangeToEnclosingScope(referenceParent)
}

// Java method: LoopClosureTransformer.assignLoopVarToLoopObjectProperty(com.google.javascript.rhino.Node loopNode, com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.LoopObject loopObject, com.google.javascript.jscomp.Var var, java.lang.String newPropertyName, com.google.javascript.rhino.Node reference)
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) priv_LoopClosureTransformer_assignLoopVarToLoopObjectProperty_2582(loopNode pkg_rhino.Node, loopObject *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject, var_ *Var, newPropertyName string, reference pkg_rhino.Node) {
	_ = loopNode
	_ = loopObject
	_ = var_
	_ = newPropertyName
	_ = reference
	if !(reference == var_.GetNameNode()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(reference))
	}
	var referenceParent pkg_rhino.Node = reference.GetParent()
	_ = referenceParent
	if !(NodeUtil_IsNameDeclaration(referenceParent)) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(referenceParent))
	}
	if !(reference.IsName()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(reference))
	}
	var loopVarReference pkg_rhino.Node = reference.CloneNode_0()
	_ = loopVarReference
	var forInPropAssignmentStatemnt pkg_rhino.Node = pkg_rhino.IR_ExprResult(this_.outer_.astFactory.CreateAssign_2_Node_Node(this_.priv_LoopClosureTransformer_createLoopVarReferenceReplacement_2584(loopObject, reference, newPropertyName), loopVarReference)).SrcrefTreeIfMissing(reference)
	_ = forInPropAssignmentStatemnt
	var loopUpdateStatement pkg_rhino.Node = loopNode.GetLastChild().GetFirstChild()
	_ = loopUpdateStatement
	forInPropAssignmentStatemnt.InsertAfter(loopUpdateStatement)
}

// Java method: LoopClosureTransformer.renameVarsToProperties(com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.LoopObject loopObject, com.google.javascript.rhino.Node objectLitNextIteration, com.google.javascript.rhino.Node scopeRoot)
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) priv_LoopClosureTransformer_renameVarsToProperties_2583(loopObject *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject, objectLitNextIteration pkg_rhino.Node, scopeRoot pkg_rhino.Node) {
	_ = loopObject
	_ = objectLitNextIteration
	_ = scopeRoot
	{
		cur_20 := loopObject.vars.Iter()
		for cur_20.HasNext() {

			var_ := cur_20.
				Next()
			_ = var_
			var newPropertyName string = this_.priv_LoopClosureTransformer_getLoopObjPropName_2573(var_)
			_ = newPropertyName
			var newPropertyValue pkg_rhino.Node
			_ = newPropertyValue
			if var_.GetScopeRoot() == scopeRoot {
				newPropertyValue = this_.priv_LoopClosureTransformer_createLoopVarReferenceReplacement_2584(loopObject, var_.GetNameNode(), newPropertyName)
			} else {
				newPropertyValue = this_.outer_.astFactory.CreateUndefinedValue()
			}
			objectLitNextIteration.AddChildToBack(this_.outer_.astFactory.CreateStringKey(newPropertyName, newPropertyValue))
		}
	}
}

// Java method: LoopClosureTransformer.createLoopVarReferenceReplacement(com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.LoopObject loopObject, com.google.javascript.rhino.Node reference, java.lang.String propertyName)
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) priv_LoopClosureTransformer_createLoopVarReferenceReplacement_2584(loopObject *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject, reference pkg_rhino.Node, propertyName string) pkg_rhino.Node {
	_ = loopObject
	_ = reference
	_ = propertyName
	var replacement pkg_rhino.Node = this_.outer_.astFactory.CreateGetProp(this_.priv_LoopClosureTransformer_createLoopObjectNameNode_2586(loopObject), propertyName, AstFactory_S_Type_1_Node(reference))
	_ = replacement
	replacement.SrcrefTree(reference)
	return replacement
}

// Java method: LoopClosureTransformer.createUniqueParameterNameToUseWithinLoop(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.LoopObject loopObject)
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) priv_LoopClosureTransformer_createUniqueParameterNameToUseWithinLoop_2585(node pkg_rhino.Node, loopObject *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject) pkg_rhino.Node {
	_ = node
	_ = loopObject
	return this_.outer_.astFactory.CreateName_2_string_AstFactory_Type((j2gort.StrValString(loopObject.name) + "$jscomp$loop_param$" + j2gort.StrValString(this_.outer_.uniqueIdSupplier.GetUniqueId(this_.outer_.compiler.GetInput(NodeUtil_GetInputId(node))))), AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_TOP_OBJECT()))
}

// Java method: LoopClosureTransformer.createLoopObjectNameNode(com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.LoopObject loopObject)
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) priv_LoopClosureTransformer_createLoopObjectNameNode_2586(loopObject *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject) pkg_rhino.Node {
	_ = loopObject
	return this_.outer_.astFactory.CreateName_2_string_AstFactory_Type(loopObject.name, AstFactory_S_Type_1_Color(pkg_colors.StaticGet_StandardColors_TOP_OBJECT()))
}

// Java class: com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.LoopObject
type Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject struct {
	name    string                        `gson:"name"`
	vars    j2gort.CustomOrderedSet[*Var] `gson:"vars"`
	idHash_ atomic.Int32
}

func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject) BindSelf_(s any) {
	_ = s
}
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject) IsJavaNil_() bool {
	return this_ == nil
}
func (*Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration$LoopClosureTransformer$LoopObject", "LoopObject"
}
func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration.LoopClosureTransformer.LoopObject", "com.google.javascript.jscomp.Es6RewriteBlockScopedDeclaration$LoopClosureTransformer$LoopObject", func(v_ any) bool {
		return j2gort.InstanceOf[*Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject](v_)
	})
}

func (this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject) initFields_Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject() {
	_ = this_
	this_.name = j2gort.NullString
	this_.vars = j2gort.NewCustomOrderedSet[*Var]()
}

// Java constructor: LoopObject(java.lang.String name)
func NewEs6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject(name string) *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject {
	this_ := &Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject{}
	InitEs6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject(this_, name)
	return this_
}

func InitEs6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject(this_ *Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject, name string) {
	_ = name
	this_.initFields_Es6RewriteBlockScopedDeclaration_LoopClosureTransformer_LoopObject()
	this_.name = name
}
func (*Es6RewriteBlockScopedDeclaration) AlwaysShouldTraverse_()             {}
func (*Es6RewriteBlockScopedDeclaration) CanReleaseSyntacticScopes_()        {}
func (*Es6RewriteBlockScopedDeclaration_LoopClosureTransformer) NoOpVisit_() {}
