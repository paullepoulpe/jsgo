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
// Port of com/google/javascript/jscomp/Normalize.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.Normalize
type Normalize struct {
	compiler       AbstractCompiler `gson:"compiler"`
	astFactory     *AstFactory      `gson:"astFactory"`
	idHash_        atomic.Int32
	assertOnChange bool `gson:"assertOnChange"`
}

func (this_ *Normalize) BindSelf_(s any) {
	_ = s
}
func (this_ *Normalize) IsJavaNil_() bool { return this_ == nil }
func (*Normalize) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Normalize", "Normalize"
}
func (this_ *Normalize) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Normalize) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Normalize", "com.google.javascript.jscomp.Normalize", func(v_ any) bool {
		return j2gort.InstanceOf[*Normalize](v_)
	})
}
func (*Normalize) Implements_jscomp_CompilerPass() {}

func (this_ *Normalize) initFields_Normalize() {
	_ = this_
}

// Java constructor: Normalize(com.google.javascript.jscomp.Normalize.Builder builder)
func NewNormalize(builder *Normalize_Builder) *Normalize {
	this_ := &Normalize{}
	InitNormalize(this_, builder)
	return this_
}

func InitNormalize(this_ *Normalize, builder *Normalize_Builder) {
	_ = builder
	this_.initFields_Normalize()
	this_.compiler = builder.compiler
	this_.assertOnChange = builder.assertOnChange
	this_.astFactory = builder.compiler.CreateAstFactory()
}

// Java static method: Normalize.createNormalizeForOptimizations(com.google.javascript.jscomp.AbstractCompiler compiler)
func Normalize_CreateNormalizeForOptimizations(compiler AbstractCompiler) *Normalize {
	_ = compiler
	return Normalize_S_Builder(compiler).Build()
}

// Java static method: Normalize.builder(com.google.javascript.jscomp.AbstractCompiler compiler)
func Normalize_S_Builder(compiler AbstractCompiler) *Normalize_Builder {
	_ = compiler
	return NewNormalize_Builder(compiler)
}

// Java method: Normalize.reportCodeChange(java.lang.String changeDescription, com.google.javascript.rhino.Node n)
func (this_ *Normalize) priv_Normalize_reportCodeChange_4045(changeDescription string, n pkg_rhino.Node) {
	_ = changeDescription
	_ = n
	if this_.assertOnChange {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Normalize constraints violated:\n"+j2gort.StrValString(changeDescription))))
	}
	this_.compiler.ReportChangeToEnclosingScope(n)
}

// Java static method: Normalize.isFunctionExpressionNameNode(com.google.javascript.rhino.Node n)
func Normalize_IsFunctionExpressionNameNode(n pkg_rhino.Node) bool {
	_ = n
	if (n == 0) || !(n.IsName()) {
		return false
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if parent == 0 {
		return false
	}
	return (NodeUtil_IsFunctionExpression(parent) && n.IsFirstChildOf(parent))
}

// Java method: Normalize.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *Normalize) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var renamer *MakeDeclaredNamesUnique = MakeDeclaredNamesUnique_S_Builder().WithAssertOnChange(this_.assertOnChange).Build()
	_ = renamer
	defer renamer.ReleaseSlabs_()

	NodeTraversal_TraverseRoots(this_.compiler, NewNormalize_NormalizeStatements(this_.compiler, this_.assertOnChange, renamer), externs, root)
	{
		sc_ := NewSyntacticScopeCreator_2(this_.compiler, NewNormalize_DuplicateDeclarationHandler(this_))
		defer sc_.ReleaseSlabs_()
		{
			clinit_NodeTraversal()
			var ntb_ NodeTraversal_Builder
			InitNodeTraversal_Builder(&ntb_)
			(&ntb_).SetCompiler(this_.compiler)
			(&ntb_).SetCallback_1_NodeTraversal_Callback(NewNormalize_PropagateConstantPropertyOverVars(this_.compiler, this_.assertOnChange))
			(&ntb_).SetScopeCreator(sc_)
			(&ntb_).TraverseRoots(externs, root)
		}
	}

	if !(this_.compiler.GetLifeCycleStage().IsNormalized()) {
		this_.compiler.SetLifeCycleStage(StaticGet_AbstractCompiler_LifeCycleStage_NORMALIZED())
	}
}

// Java class: com.google.javascript.jscomp.Normalize.Builder
type Normalize_Builder struct {
	compiler       AbstractCompiler `gson:"compiler"`
	idHash_        atomic.Int32
	assertOnChange bool `gson:"assertOnChange"`
}

func (this_ *Normalize_Builder) BindSelf_(s any) {
	_ = s
}
func (this_ *Normalize_Builder) IsJavaNil_() bool { return this_ == nil }
func (*Normalize_Builder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Normalize$Builder", "Builder"
}
func (this_ *Normalize_Builder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Normalize_Builder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Normalize.Builder", "com.google.javascript.jscomp.Normalize$Builder", func(v_ any) bool {
		return j2gort.InstanceOf[*Normalize_Builder](v_)
	})
}

func (this_ *Normalize_Builder) initFields_Normalize_Builder() {
	_ = this_

}

// Java constructor: Builder(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewNormalize_Builder(compiler AbstractCompiler) *Normalize_Builder {
	this_ := &Normalize_Builder{}
	InitNormalize_Builder(this_, compiler)
	return this_
}

func InitNormalize_Builder(this_ *Normalize_Builder, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_Normalize_Builder()
	this_.compiler = compiler
}

// Java method: Builder.assertOnChange(boolean assertOnChange)
func (this_ *Normalize_Builder) AssertOnChange(assertOnChange bool) *Normalize_Builder {
	_ = assertOnChange
	this_.assertOnChange = assertOnChange
	return this_
}

// Java method: Builder.build()
func (this_ *Normalize_Builder) Build() *Normalize {
	return NewNormalize(this_)
}

// Java class: com.google.javascript.jscomp.Normalize.PropagateConstantPropertyOverVars
type Normalize_PropagateConstantPropertyOverVars struct {
	compiler       AbstractCompiler `gson:"compiler"`
	idHash_        atomic.Int32
	assertOnChange bool `gson:"assertOnChange"`
}

func (this_ *Normalize_PropagateConstantPropertyOverVars) BindSelf_(s any) {
	_ = s
}
func (this_ *Normalize_PropagateConstantPropertyOverVars) IsJavaNil_() bool { return this_ == nil }
func (*Normalize_PropagateConstantPropertyOverVars) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Normalize$PropagateConstantPropertyOverVars", "PropagateConstantPropertyOverVars"
}
func (this_ *Normalize_PropagateConstantPropertyOverVars) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Normalize_PropagateConstantPropertyOverVars) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Normalize.PropagateConstantPropertyOverVars", "com.google.javascript.jscomp.Normalize$PropagateConstantPropertyOverVars", func(v_ any) bool {
		return j2gort.InstanceOf[*Normalize_PropagateConstantPropertyOverVars](v_)
	})
}
func (*Normalize_PropagateConstantPropertyOverVars) Implements_jscomp_NodeTraversal_ScopedCallback() {
}
func (*Normalize_PropagateConstantPropertyOverVars) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *Normalize_PropagateConstantPropertyOverVars) initFields_Normalize_PropagateConstantPropertyOverVars() {
	_ = this_
}

// Java constructor: PropagateConstantPropertyOverVars(com.google.javascript.jscomp.AbstractCompiler compiler, boolean forbidChanges)
func NewNormalize_PropagateConstantPropertyOverVars(compiler AbstractCompiler, forbidChanges bool) *Normalize_PropagateConstantPropertyOverVars {
	this_ := &Normalize_PropagateConstantPropertyOverVars{}
	InitNormalize_PropagateConstantPropertyOverVars(this_, compiler, forbidChanges)
	return this_
}

func InitNormalize_PropagateConstantPropertyOverVars(this_ *Normalize_PropagateConstantPropertyOverVars, compiler AbstractCompiler, forbidChanges bool) {
	_ = compiler
	_ = forbidChanges
	this_.initFields_Normalize_PropagateConstantPropertyOverVars()
	this_.compiler = compiler
	this_.assertOnChange = forbidChanges
}

// Java method: PropagateConstantPropertyOverVars.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Normalize_PropagateConstantPropertyOverVars) EnterScope(t *NodeTraversal) {
	_ = t
	t.GetScope()
}

// Java method: PropagateConstantPropertyOverVars.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Normalize_PropagateConstantPropertyOverVars) ExitScope(t *NodeTraversal) {
	_ = t
}

// Java method: PropagateConstantPropertyOverVars.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Normalize_PropagateConstantPropertyOverVars) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return true
}

// Java method: PropagateConstantPropertyOverVars.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Normalize_PropagateConstantPropertyOverVars) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(n.IsName()) || (j2gort.CheckNotNullString(n.GetString()) == "") {
		return
	}
	var var_ *Var = j2gort.CastPtr[Var](t.GetScope().GetVar(n.GetString()))
	_ = var_
	var info *pkg_rhino.JSDocInfo
	_ = info
	if var_ != nil {
		info = var_.GetJSDocInfo()
	} else {
		info = nil
	}
	var shouldBeConstant bool = ((((var_ != nil) && (var_.IsConst() || Normalize_IsFunctionExpressionNameNode(var_.GetNode()))) || ((info != nil) && info.IsConstant())) || NodeUtil_IsConstantByConvention(this_.compiler.GetCodingConvention(), n))
	_ = shouldBeConstant
	var isMarkedConstant bool = n.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME())
	_ = isMarkedConstant
	if shouldBeConstant && !(isMarkedConstant) {
		if this_.assertOnChange {
			var name string = n.GetString()
			_ = name
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Should be const but not marked as const.\n  name: " + j2gort.StrValString(name) + "\n  parent:" + j2gort.StrValString(n.GetParent().ToStringTree()))))))
		}
		n.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	}
}

// Java class: com.google.javascript.jscomp.Normalize.NormalizeStatements
type Normalize_NormalizeStatements struct {
	compiler                                AbstractCompiler                         `gson:"compiler"`
	astFactory                              *AstFactory                              `gson:"astFactory"`
	makeDeclaredNamesUnique                 *MakeDeclaredNamesUnique                 `gson:"makeDeclaredNamesUnique"`
	rewriteLogicalAssignmentOperatorsHelper *RewriteLogicalAssignmentOperatorsHelper `gson:"rewriteLogicalAssignmentOperatorsHelper"`
	idHash_                                 atomic.Int32
	assertOnChange                          bool `gson:"assertOnChange"`
}

func (this_ *Normalize_NormalizeStatements) BindSelf_(s any) {
	_ = s
}
func (this_ *Normalize_NormalizeStatements) IsJavaNil_() bool { return this_ == nil }
func (*Normalize_NormalizeStatements) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Normalize$NormalizeStatements", "NormalizeStatements"
}
func (this_ *Normalize_NormalizeStatements) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Normalize_NormalizeStatements) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Normalize.NormalizeStatements", "com.google.javascript.jscomp.Normalize$NormalizeStatements", func(v_ any) bool {
		return j2gort.InstanceOf[*Normalize_NormalizeStatements](v_)
	})
}
func (*Normalize_NormalizeStatements) Implements_jscomp_NodeTraversal_ScopedCallback() {}
func (*Normalize_NormalizeStatements) Implements_jscomp_NodeTraversal_Callback()       {}

func (this_ *Normalize_NormalizeStatements) initFields_Normalize_NormalizeStatements() {
	_ = this_
}

// Java constructor: NormalizeStatements(com.google.javascript.jscomp.AbstractCompiler compiler, boolean assertOnChange, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable MakeDeclaredNamesUnique makeDeclaredNamesUnique)
func NewNormalize_NormalizeStatements(compiler AbstractCompiler, assertOnChange bool, makeDeclaredNamesUnique *MakeDeclaredNamesUnique) *Normalize_NormalizeStatements {
	this_ := &Normalize_NormalizeStatements{}
	InitNormalize_NormalizeStatements(this_, compiler, assertOnChange, makeDeclaredNamesUnique)
	return this_
}

func InitNormalize_NormalizeStatements(this_ *Normalize_NormalizeStatements, compiler AbstractCompiler, assertOnChange bool, makeDeclaredNamesUnique *MakeDeclaredNamesUnique) {
	_ = compiler
	_ = assertOnChange
	_ = makeDeclaredNamesUnique
	this_.initFields_Normalize_NormalizeStatements()
	this_.compiler = compiler
	this_.assertOnChange = assertOnChange
	this_.makeDeclaredNamesUnique = makeDeclaredNamesUnique
	this_.astFactory = compiler.CreateAstFactory()
	this_.rewriteLogicalAssignmentOperatorsHelper = NewRewriteLogicalAssignmentOperatorsHelper(compiler, this_.astFactory, compiler.GetUniqueIdSupplier())
}

// Java method: NormalizeStatements.reportCodeChange(java.lang.String changeDescription, com.google.javascript.rhino.Node n)
func (this_ *Normalize_NormalizeStatements) priv_NormalizeStatements_reportCodeChange_4048(changeDescription string, n pkg_rhino.Node) {
	_ = changeDescription
	_ = n
	if this_.assertOnChange {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Normalize constraints violated:\n"+j2gort.StrValString(changeDescription))))
	}
	this_.compiler.ReportChangeToEnclosingScope(n)
}

// Java method: NormalizeStatements.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Normalize_NormalizeStatements) EnterScope(t *NodeTraversal) {
	_ = t
	if this_.makeDeclaredNamesUnique != nil {
		this_.makeDeclaredNamesUnique.EnterScope(t)
	}
}

// Java method: NormalizeStatements.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *Normalize_NormalizeStatements) ExitScope(t *NodeTraversal) {
	_ = t
	if this_.makeDeclaredNamesUnique != nil {
		this_.makeDeclaredNamesUnique.ExitScope(t)
	}
}

// Java method: NormalizeStatements.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Normalize_NormalizeStatements) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if this_.makeDeclaredNamesUnique != nil {
		var unused bool = this_.makeDeclaredNamesUnique.ShouldTraverse_3(t, n, parent)
		_ = unused
	}
	var oldParent pkg_rhino.Node = n.GetParent()
	_ = oldParent
	this_.priv_NormalizeStatements_doStatementNormalizations_4054(n)
	if !(n.GetParent() == oldParent) {
		j2gort.CheckStateStr(false, "n's parent changed during normalization. It needs to stay the same so that we properly visit any added nodes.")
	}
	return true
}

// Java method: NormalizeStatements.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *Normalize_NormalizeStatements) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if this_.makeDeclaredNamesUnique != nil {
		this_.makeDeclaredNamesUnique.Visit_3(t, n, parent)
	}
	switch n.TokenSwitchOrd_() {
	case int32(78) /* Token.WHILE */ :
		var expr pkg_rhino.Node = n.GetFirstChild()
		_ = expr
		n.SetToken(pkg_rhino.StaticGet_Token_FOR())
		var empty pkg_rhino.Node = pkg_rhino.IR_Empty()
		_ = empty
		empty.SrcrefIfMissing(n)
		empty.InsertBefore(expr)
		empty.CloneNode_0().InsertAfter(expr)
		this_.priv_NormalizeStatements_reportCodeChange_4048("WHILE node", n)
	case int32(46) /* Token.ARRAYLIT */, int32(28) /* Token.CALL */, int32(49) /* Token.PARAM_LIST */, int32(23) /* Token.NEW */, int32(47) /* Token.OBJECTLIT */, int32(31) /* Token.OPTCHAIN_CALL */ :
		n.SetTrailingComma(false)
	case int32(32) /* Token.NAME */ :
		this_.priv_NormalizeStatements_annotateConstantsByConvention_4049(n)
		this_.priv_NormalizeStatements_annotateFunctionExpressionNameAsConstant_4050(n)
	case int32(104) /* Token.DESTRUCTURING_LHS */ :
		this_.priv_NormalizeStatements_normalizeDestructuringLhs_4055(n)
	case int32(64) /* Token.ASSIGN_OR */, int32(65) /* Token.ASSIGN_AND */, int32(66) /* Token.ASSIGN_COALESCE */ :
		this_.rewriteLogicalAssignmentOperatorsHelper.VisitLogicalAssignmentOperator(t, n)
	default:
	}
}

// Java method: NormalizeStatements.annotateConstantsByConvention(com.google.javascript.rhino.Node n)
func (this_ *Normalize_NormalizeStatements) priv_NormalizeStatements_annotateConstantsByConvention_4049(n pkg_rhino.Node) {
	_ = n
	if !(n.IsName()) {
		j2gort.CheckState(false)
	}
	if this_.compiler.GetLifeCycleStage().IsNormalizedObfuscated() {
		return
	}
	if !(n.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME())) && NodeUtil_IsConstantByConvention(this_.compiler.GetCodingConvention(), n) {
		if !(!(this_.assertOnChange)) {
			j2gort.CheckStateFmt(false, "Not marked as constant when it should be: %s", j2gort.BoxNode(n))
		}
		n.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	}
}

// Java method: NormalizeStatements.annotateFunctionExpressionNameAsConstant(com.google.javascript.rhino.Node n)
func (this_ *Normalize_NormalizeStatements) priv_NormalizeStatements_annotateFunctionExpressionNameAsConstant_4050(n pkg_rhino.Node) {
	_ = n
	if !(n.IsName()) {
		j2gort.CheckStateFmt(false, "Expected NAME node but got %s", j2gort.BoxString(j2gort.CheckNotNullEnum(n.GetToken()).String()))
	}
	if Normalize_IsFunctionExpressionNameNode(n) {
		n.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	}
}

// Java method: NormalizeStatements.splitExportDeclaration(com.google.javascript.rhino.Node n)
func (this_ *Normalize_NormalizeStatements) priv_NormalizeStatements_splitExportDeclaration_4051(n pkg_rhino.Node) {
	_ = n
	if n.GetBooleanProp(pkg_rhino.StaticGet_Node_EXPORT_DEFAULT()) {
		return
	}
	var destructuringLhsNodes []pkg_rhino.Node = []pkg_rhino.Node{}
	_ = destructuringLhsNodes
	var c pkg_rhino.Node = n.GetFirstChild()
	_ = c
	if NodeUtil_IsDeclaration(c) {
		c.Detach()
		var exportSpecs pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_EXPORT_SPECS()).Srcref(n)
		_ = exportSpecs
		n.AddChildToFront(exportSpecs)
		if c.IsClass() || c.IsFunction() {
			var name pkg_rhino.Node = c.GetFirstChild()
			_ = name
			c.InsertBefore(n)
			this_.priv_NormalizeStatements_addNameNodeToExportSpecs_4052(exportSpecs, name)
		} else {
			NodeUtil_VisitLhsNodesInNode(c, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_1 any) {
				name_2 := j2gort.CastNode[pkg_rhino.Node](lp_1)
				_ = name_2
				this_.priv_NormalizeStatements_addNameNodeToExportSpecs_4052(exportSpecs, name_2)
			}))
			{
				var child pkg_rhino.Node = c.GetFirstChild()
				_ = child
				for child != 0 {
					var next pkg_rhino.Node = child.GetNext()
					_ = next
					child.Detach()
					var newDeclaration pkg_rhino.Node = pkg_rhino.NewNode_2(c.GetToken(), child).Srcref(n)
					_ = newDeclaration
					newDeclaration.InsertBefore(n)
					if child.IsDestructuringLhs() {
						destructuringLhsNodes = append(destructuringLhsNodes, child)
					}
					child = next
				}
			}
		}
		for _, destructuringLhsNode := range destructuringLhsNodes {
			_ = destructuringLhsNode
			this_.priv_NormalizeStatements_normalizeDestructuringLhs_4055(destructuringLhsNode)
		}

		this_.priv_NormalizeStatements_reportCodeChange_4048("combined export and declaration", n.GetParent())
	}
}

// Java method: NormalizeStatements.addNameNodeToExportSpecs(com.google.javascript.rhino.Node exportSpecs, com.google.javascript.rhino.Node name)
func (this_ *Normalize_NormalizeStatements) priv_NormalizeStatements_addNameNodeToExportSpecs_4052(exportSpecs pkg_rhino.Node, name pkg_rhino.Node) {
	_ = exportSpecs
	_ = name
	var exportSpec pkg_rhino.Node = pkg_rhino.NewNode_1(pkg_rhino.StaticGet_Token_EXPORT_SPEC()).Srcref(name)
	_ = exportSpec
	exportSpec.AddChildToFront(name.CloneNode_0())
	exportSpec.AddChildToFront(name.CloneNode_0())
	exportSpecs.AddChildToBack(exportSpec)
}

// Java method: NormalizeStatements.normalizeBlocklessFunction(com.google.javascript.rhino.Node n)
func (this_ *Normalize_NormalizeStatements) priv_NormalizeStatements_normalizeBlocklessFunction_4053(n pkg_rhino.Node) {
	_ = n
	if !(n.IsFunction()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
	}
	if !(NodeUtil_GetFunctionBody(n).IsBlock()) {
		var returnValue pkg_rhino.Node = NodeUtil_GetFunctionBody(n)
		_ = returnValue
		var body pkg_rhino.Node = pkg_rhino.IR_Block_1_Node(pkg_rhino.IR_ReturnNode_1(returnValue.Detach()))
		_ = body
		body.SrcrefTreeIfMissing(returnValue)
		n.AddChildToBack(body)
		this_.priv_NormalizeStatements_reportCodeChange_4048("blockless arrow function", body)
	}
}

// Java method: NormalizeStatements.doStatementNormalizations(com.google.javascript.rhino.Node n)
func (this_ *Normalize_NormalizeStatements) priv_NormalizeStatements_doStatementNormalizations_4054(n pkg_rhino.Node) {
	_ = n
	if n.IsLabel() {
		this_.priv_NormalizeStatements_normalizeLabels_4058(n)
	}
	if NodeUtil_IsStatementBlock(n) || n.IsLabel() {
		this_.priv_NormalizeStatements_extractForInitializer_4059(n, pkg_rhino.Node(0), pkg_rhino.Node(0))
	}
	if NodeUtil_IsStatementBlock(n) {
		this_.priv_NormalizeStatements_splitVarDeclarations_4060(n)
	}
	if n.IsFunction() {
		this_.priv_NormalizeStatements_normalizeBlocklessFunction_4053(n)
		this_.priv_NormalizeStatements_moveNamedFunctions_4061(n.GetLastChild())
	}
	if NodeUtil_IsClassStaticBlock(n) {
		this_.priv_NormalizeStatements_moveNamedFunctions_4061(n)
	}
	if n.IsExport() {
		this_.priv_NormalizeStatements_splitExportDeclaration_4051(n)
	}
	if NodeUtil_IsCompoundAssignmentOp(n) && !(NodeUtil_IsLogicalAssignmentOp(n)) {
		this_.priv_NormalizeStatements_normalizeAssignShorthand_4062(n)
	}
}

// Java method: NormalizeStatements.normalizeDestructuringLhs(com.google.javascript.rhino.Node n)
func (this_ *Normalize_NormalizeStatements) priv_NormalizeStatements_normalizeDestructuringLhs_4055(n pkg_rhino.Node) {
	_ = n
	if this_.priv_NormalizeStatements_isVarDestructuringDeclaration_4056(n) {
		this_.priv_NormalizeStatements_rewriteVarDestructuringDeclaration_4057(n)
	}
}

// Java method: NormalizeStatements.isVarDestructuringDeclaration(com.google.javascript.rhino.Node n)
func (this_ *Normalize_NormalizeStatements) priv_NormalizeStatements_isVarDestructuringDeclaration_4056(n pkg_rhino.Node) bool {
	_ = n
	return ((n.IsDestructuringLhs() && (n.GetFirstChild().IsArrayPattern() || n.GetFirstChild().IsObjectPattern())) && n.GetParent().IsVar())
}

// Java method: NormalizeStatements.rewriteVarDestructuringDeclaration(com.google.javascript.rhino.Node destructuringLhs)
func (this_ *Normalize_NormalizeStatements) priv_NormalizeStatements_rewriteVarDestructuringDeclaration_4057(destructuringLhs pkg_rhino.Node) {
	_ = destructuringLhs
	if !(!(destructuringLhs.GetGrandparent().IsExport())) {
		j2gort.CheckStateStr(false, "Export destructuring declarations should be already handled in `splitExportDeclaration` during pre-order traversal.")
	}
	var var_ pkg_rhino.Node = destructuringLhs.GetParent()
	_ = var_
	NodeUtil_VisitLhsNodesInNode(destructuringLhs, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_3 any) {
		name := j2gort.CastNode[pkg_rhino.Node](lp_3)
		_ = name
		if !(name.IsName()) {
			j2gort.CheckStateFmt(false, "lhs in destructuring declaration should be a simple name. (%s)", j2gort.BoxNode(name))
		}
		var newName pkg_rhino.Node = pkg_rhino.IR_Name(name.GetString()).Srcref(name)
		_ = newName
		if name.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME()) {
			newName.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
		}
		var newVar pkg_rhino.Node = pkg_rhino.IR_Var_1(newName).Srcref(name)
		_ = newVar
		newVar.InsertBefore(var_)
	}))
	var destructuringPattern pkg_rhino.Node = destructuringLhs.RemoveFirstChild()
	_ = destructuringPattern
	if !(destructuringPattern.IsDestructuringPattern()) {
		j2gort.CheckStateStr(false, "Expected destructuring pattern.")
	}
	var rhs pkg_rhino.Node = destructuringLhs.RemoveFirstChild()
	_ = rhs
	var assign pkg_rhino.Node = this_.astFactory.CreateAssign_2_Node_Node(destructuringPattern, rhs)
	_ = assign
	assign.Srcref(var_)
	var expr pkg_rhino.Node = this_.astFactory.ExprResult(assign)
	_ = expr
	expr.Srcref(var_)
	var_.ReplaceWith(expr)
	destructuringLhs.Detach()
	this_.priv_NormalizeStatements_reportCodeChange_4048("Var destructuring declaration rewritten", expr)
}

// Java method: NormalizeStatements.normalizeLabels(com.google.javascript.rhino.Node n)
func (this_ *Normalize_NormalizeStatements) priv_NormalizeStatements_normalizeLabels_4058(n pkg_rhino.Node) {
	_ = n
	if !(n.IsLabel()) {
		j2gort.CheckArgument(false)
	}
	var last pkg_rhino.Node = n.GetLastChild()
	_ = last
	switch last.TokenSwitchOrd_() {
	case int32(92) /* Token.LABEL */, int32(90) /* Token.BLOCK */, int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(78) /* Token.WHILE */, int32(79) /* Token.DO */ :
		return
	default:
		var block pkg_rhino.Node = pkg_rhino.IR_Block_0()
		_ = block
		block.SrcrefIfMissing(last)
		last.ReplaceWith(block)
		block.AddChildToFront(last)
		this_.priv_NormalizeStatements_reportCodeChange_4048("LABEL normalization", n)
		return
	}
}

// Java method: NormalizeStatements.extractForInitializer(com.google.javascript.rhino.Node n, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node before, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node beforeParent)
func (this_ *Normalize_NormalizeStatements) priv_NormalizeStatements_extractForInitializer_4059(n pkg_rhino.Node, before pkg_rhino.Node, beforeParent pkg_rhino.Node) {
	_ = n
	_ = before
	_ = beforeParent
	{
		var next pkg_rhino.Node = pkg_rhino.Node(0)
		_ = next
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = next {
			next = c.GetNext()
			var insertBefore pkg_rhino.Node
			_ = insertBefore
			if before == 0 {
				insertBefore = c
			} else {
				insertBefore = before
			}
			var insertBeforeParent pkg_rhino.Node
			_ = insertBeforeParent
			if before == 0 {
				insertBeforeParent = n
			} else {
				insertBeforeParent = beforeParent
			}
			switch c.TokenSwitchOrd_() {
			case int32(92) /* Token.LABEL */ :
				this_.priv_NormalizeStatements_extractForInitializer_4059(c, insertBefore, insertBeforeParent)
			case int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
				var first pkg_rhino.Node = c.GetFirstChild()
				_ = first
				if first.IsVar() {
					var lhs pkg_rhino.Node = first.GetFirstChild()
					_ = lhs
					if lhs.IsDestructuringLhs() {
						NodeUtil_VisitLhsNodesInNode(lhs, j2gort.NewLambda_Ext_java_util_function_Consumer(func(lp_4 any) {
							name := j2gort.CastNode[pkg_rhino.Node](lp_4)
							_ = name
							if !(name.IsName()) {
								j2gort.CheckStateFmt(false, "lhs in destructuring declaration should be a simple name. (%s)", j2gort.BoxNode(name))
							}
							var newName pkg_rhino.Node = pkg_rhino.IR_Name(name.GetString()).Srcref(name)
							_ = newName
							var newVar pkg_rhino.Node = pkg_rhino.IR_Var_1(newName).Srcref(name)
							_ = newVar
							newVar.InsertBefore(insertBefore)
						}))
						var destructuringPattern pkg_rhino.Node = lhs.RemoveFirstChild()
						_ = destructuringPattern
						first.ReplaceWith(destructuringPattern)
					} else {
						var newStatement pkg_rhino.Node = first
						_ = newStatement
						var name pkg_rhino.Node = newStatement.GetFirstChild().CloneNode_0()
						_ = name
						first.ReplaceWith(name)
						newStatement.InsertBefore(insertBefore)
					}
					this_.priv_NormalizeStatements_reportCodeChange_4048("FOR-IN var declaration", n)
				}
			case int32(80) /* Token.FOR */ :
				if !(c.GetFirstChild().IsEmpty()) {
					var init pkg_rhino.Node = c.GetFirstChild()
					_ = init
					if ((init.IsLet() || init.IsConst()) || init.IsClass()) || init.IsFunction() {
						continue
					}
					var empty pkg_rhino.Node = pkg_rhino.IR_Empty()
					_ = empty
					empty.SrcrefIfMissing(c)
					init.ReplaceWith(empty)
					var newStatement_2 pkg_rhino.Node = pkg_rhino.Node(0)
					_ = newStatement_2
					if init.IsVar() {
						newStatement_2 = init
					} else {
						newStatement_2 = NodeUtil_NewExpr(init)
					}
					newStatement_2.InsertBefore(insertBefore)
					this_.priv_NormalizeStatements_reportCodeChange_4048("FOR initializer", n)
				}
			default:
			}
		}
	}
}

// Java method: NormalizeStatements.splitVarDeclarations(com.google.javascript.rhino.Node n)
func (this_ *Normalize_NormalizeStatements) priv_NormalizeStatements_splitVarDeclarations_4060(n pkg_rhino.Node) {
	_ = n
	{
		var next pkg_rhino.Node = pkg_rhino.Node(0)
		_ = next
		var c pkg_rhino.Node = n.GetFirstChild()
		_ = c
		for ; c != 0; c = next {
			next = c.GetNext()
			if NodeUtil_IsNameDeclaration(c) {
				if this_.assertOnChange && !(c.HasChildren()) {
					panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Empty VAR node.")))
				}
				for c.GetFirstChild() != c.GetLastChild() {
					var name pkg_rhino.Node = c.GetFirstChild()
					_ = name
					name.Detach()
					var newVar pkg_rhino.Node = pkg_rhino.NewNode_2(c.GetToken(), name).Srcref(n)
					_ = newVar
					newVar.InsertBefore(c)
					this_.priv_NormalizeStatements_reportCodeChange_4048("VAR with multiple children", n)
				}
			}
		}
	}
}

// Java method: NormalizeStatements.moveNamedFunctions(com.google.javascript.rhino.Node body)
func (this_ *Normalize_NormalizeStatements) priv_NormalizeStatements_moveNamedFunctions_4061(body pkg_rhino.Node) {
	_ = body
	if !(body.GetParent().IsFunction() || NodeUtil_IsClassStaticBlock(body)) {
		j2gort.CheckState(false)
	}
	var insertAfter pkg_rhino.Node = pkg_rhino.Node(0)
	_ = insertAfter
	var current pkg_rhino.Node = body.GetFirstChild()
	_ = current
	for (current != 0) && NodeUtil_IsFunctionDeclaration(current) {
		insertAfter = current
		current = current.GetNext()
	}
	for current != 0 {
		var next pkg_rhino.Node = current.GetNext()
		_ = next
		if NodeUtil_IsFunctionDeclaration(current) {
			current.Detach()
			insertAfter = Normalize_NormalizeStatements_AddToFront(body, current, insertAfter)
			this_.priv_NormalizeStatements_reportCodeChange_4048("Move function declaration not at top of function or class static block", body)
		}
		current = next
	}
}

// Java method: NormalizeStatements.normalizeAssignShorthand(com.google.javascript.rhino.Node shorthand)
func (this_ *Normalize_NormalizeStatements) priv_NormalizeStatements_normalizeAssignShorthand_4062(shorthand pkg_rhino.Node) {
	_ = shorthand
	if !(shorthand.GetFirstChild().IsName()) {
		return
	}
	if !(shorthand.HasTwoChildren()) {
		j2gort.CheckArgumentStr(false, "Expected two children for shorthand assignment.")
	}
	var lhs pkg_rhino.Node = shorthand.GetFirstChild()
	_ = lhs
	var rhs pkg_rhino.Node = shorthand.GetSecondChild()
	_ = rhs
	lhs.Detach()
	rhs.Detach()
	var op pkg_rhino.Token = NodeUtil_GetOpFromAssignmentOp(shorthand)
	_ = op
	var opNode pkg_rhino.Node = pkg_rhino.NewNode_3(op, lhs.CloneNode_0().Srcref(lhs), rhs).Srcref(shorthand).CopyTypeFrom(shorthand)
	_ = opNode
	shorthand.SetToken(pkg_rhino.StaticGet_Token_ASSIGN())
	lhs.CopyTypeFrom(shorthand)
	shorthand.AddChildToBack(lhs)
	shorthand.AddChildToBack(opNode)
	this_.priv_NormalizeStatements_reportCodeChange_4048("assign shorthand", shorthand)
}

// Java static method: NormalizeStatements.addToFront(com.google.javascript.rhino.Node parent, com.google.javascript.rhino.Node newChild, com.google.javascript.rhino.Node after)
func Normalize_NormalizeStatements_AddToFront(parent pkg_rhino.Node, newChild pkg_rhino.Node, after pkg_rhino.Node) pkg_rhino.Node {
	_ = parent
	_ = newChild
	_ = after
	if after == 0 {
		parent.AddChildToFront(newChild)
	} else {
		newChild.InsertAfter(after)
	}
	return newChild
}

// Java class: com.google.javascript.jscomp.Normalize.DuplicateDeclarationHandler
type Normalize_DuplicateDeclarationHandler struct {
	outer_                    *Normalize
	hasOkDuplicateDeclaration j2gort.CustomOrderedSet[*Var] `gson:"hasOkDuplicateDeclaration"`
	idHash_                   atomic.Int32
}

func (this_ *Normalize_DuplicateDeclarationHandler) BindSelf_(s any) {
	_ = s
}
func (this_ *Normalize_DuplicateDeclarationHandler) IsJavaNil_() bool { return this_ == nil }
func (*Normalize_DuplicateDeclarationHandler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.Normalize$DuplicateDeclarationHandler", "DuplicateDeclarationHandler"
}
func (this_ *Normalize_DuplicateDeclarationHandler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*Normalize_DuplicateDeclarationHandler) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.Normalize.DuplicateDeclarationHandler", "com.google.javascript.jscomp.Normalize$DuplicateDeclarationHandler", func(v_ any) bool {
		return j2gort.InstanceOf[*Normalize_DuplicateDeclarationHandler](v_)
	})
}
func (this_ *Normalize_DuplicateDeclarationHandler) SetOuter_(o *Normalize) { this_.outer_ = o }
func (this_ *Normalize_DuplicateDeclarationHandler) GetOuter_() *Normalize  { return this_.outer_ }
func (*Normalize_DuplicateDeclarationHandler) Implements_jscomp_SyntacticScopeCreator_RedeclarationHandler() {
}

func (this_ *Normalize_DuplicateDeclarationHandler) initFields_Normalize_DuplicateDeclarationHandler() {
	_ = this_
	this_.hasOkDuplicateDeclaration = j2gort.NewCustomOrderedSet[*Var]()
}

// Java constructor: DuplicateDeclarationHandler()
func NewNormalize_DuplicateDeclarationHandler(outer_ *Normalize) *Normalize_DuplicateDeclarationHandler {
	this_ := &Normalize_DuplicateDeclarationHandler{}
	this_.outer_ = outer_
	InitNormalize_DuplicateDeclarationHandler(this_)
	return this_
}

func InitNormalize_DuplicateDeclarationHandler(this_ *Normalize_DuplicateDeclarationHandler) {
	this_.initFields_Normalize_DuplicateDeclarationHandler()
}

// Java method: DuplicateDeclarationHandler.onRedeclaration(com.google.javascript.jscomp.Scope s, java.lang.String name, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.CompilerInput input)
func (this_ *Normalize_DuplicateDeclarationHandler) OnRedeclaration(s *Scope, name string, n pkg_rhino.Node, input *CompilerInput) {
	_ = s
	_ = name
	_ = n
	_ = input
	if !(n.IsName()) {
		j2gort.CheckState(false)
	}
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var v *Var = j2gort.CastPtr[Var](s.GetVar(name))
	_ = v
	if s.IsGlobal() {
		if v.IsExtern() && !(input.IsExtern()) {
			if this_.hasOkDuplicateDeclaration.Add(v) {
				return
			}
		}
	}
	if parent.IsFunction() {
		if v.GetParentNode().IsVar() {
			s.Undeclare(ToIface_AbstractVar(v))
			s.Declare(name, n, v.GetInput())
			this_.priv_DuplicateDeclarationHandler_replaceVarWithAssignment_4064(v.GetNameNode(), v.GetParentNode(), v.GetParentNode().GetParent())
		}
	} else if parent.IsVar() {
		if !(parent.HasOneChild()) {
			j2gort.CheckState(false)
		}
		this_.priv_DuplicateDeclarationHandler_replaceVarWithAssignment_4064(n, parent, parent.GetParent())
	}
}

// Java method: DuplicateDeclarationHandler.replaceVarWithAssignment(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent, com.google.javascript.rhino.Node grandparent)
func (this_ *Normalize_DuplicateDeclarationHandler) priv_DuplicateDeclarationHandler_replaceVarWithAssignment_4064(n pkg_rhino.Node, parent pkg_rhino.Node, grandparent pkg_rhino.Node) {
	_ = n
	_ = parent
	_ = grandparent
	if n.HasChildren() {
		n.Detach()
		var value pkg_rhino.Node = n.GetFirstChild()
		_ = value
		value.Detach()
		var replacement pkg_rhino.Node = this_.outer_.astFactory.CreateAssign_2_Node_Node(n, value)
		_ = replacement
		replacement.SetJSDocInfo(parent.GetJSDocInfo())
		replacement.SrcrefIfMissing(parent)
		var statement pkg_rhino.Node = NodeUtil_NewExpr(replacement)
		_ = statement
		parent.ReplaceWith(statement)
		this_.outer_.priv_Normalize_reportCodeChange_4045("Duplicate VAR declaration", statement)
	} else {
		if !(NodeUtil_IsStatementBlock(grandparent)) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(grandparent))
		}
		parent.Detach()
		this_.outer_.priv_Normalize_reportCodeChange_4045("Duplicate VAR declaration", grandparent)
	}
}
func (*Normalize_PropagateConstantPropertyOverVars) AlwaysShouldTraverse_()      {}
func (*Normalize_PropagateConstantPropertyOverVars) CanReleaseSyntacticScopes_() {}
func (*Normalize_NormalizeStatements) CanReleaseSyntacticScopes_()               {}
