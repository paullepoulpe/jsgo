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
// Port of com/google/javascript/jscomp/VariableReferenceCheck.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.VariableReferenceCheck
type VariableReferenceCheck struct {
	compiler               AbstractCompiler               `gson:"compiler"`
	blocksWithDeclarations j2gort.OrderedSet[*BasicBlock] `gson:"blocksWithDeclarations"`
	idHash_                atomic.Int32
	checkUnusedLocals      bool `gson:"checkUnusedLocals"`
}

func (this_ *VariableReferenceCheck) BindSelf_(s any) {
	_ = s
}
func (this_ *VariableReferenceCheck) IsJavaNil_() bool { return this_ == nil }
func (*VariableReferenceCheck) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.VariableReferenceCheck", "VariableReferenceCheck"
}
func (this_ *VariableReferenceCheck) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*VariableReferenceCheck) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.VariableReferenceCheck", "com.google.javascript.jscomp.VariableReferenceCheck", func(v_ any) bool {
		return j2gort.InstanceOf[*VariableReferenceCheck](v_)
	})
}
func (*VariableReferenceCheck) Implements_jscomp_CompilerPass() {}

var static_VariableReferenceCheck_EARLY_REFERENCE *DiagnosticType

func StaticGet_VariableReferenceCheck_EARLY_REFERENCE() *DiagnosticType {
	clinit_VariableReferenceCheck()
	return static_VariableReferenceCheck_EARLY_REFERENCE
}
func StaticSet_VariableReferenceCheck_EARLY_REFERENCE(v *DiagnosticType) {
	clinit_VariableReferenceCheck()
	static_VariableReferenceCheck_EARLY_REFERENCE = v
}

var static_VariableReferenceCheck_EARLY_EXPORTS_REFERENCE *DiagnosticType

func StaticGet_VariableReferenceCheck_EARLY_EXPORTS_REFERENCE() *DiagnosticType {
	clinit_VariableReferenceCheck()
	return static_VariableReferenceCheck_EARLY_EXPORTS_REFERENCE
}
func StaticSet_VariableReferenceCheck_EARLY_EXPORTS_REFERENCE(v *DiagnosticType) {
	clinit_VariableReferenceCheck()
	static_VariableReferenceCheck_EARLY_EXPORTS_REFERENCE = v
}

var static_VariableReferenceCheck_REDECLARED_VARIABLE *DiagnosticType

func StaticGet_VariableReferenceCheck_REDECLARED_VARIABLE() *DiagnosticType {
	clinit_VariableReferenceCheck()
	return static_VariableReferenceCheck_REDECLARED_VARIABLE
}
func StaticSet_VariableReferenceCheck_REDECLARED_VARIABLE(v *DiagnosticType) {
	clinit_VariableReferenceCheck()
	static_VariableReferenceCheck_REDECLARED_VARIABLE = v
}

var static_VariableReferenceCheck_EARLY_REFERENCE_ERROR *DiagnosticType

func StaticGet_VariableReferenceCheck_EARLY_REFERENCE_ERROR() *DiagnosticType {
	clinit_VariableReferenceCheck()
	return static_VariableReferenceCheck_EARLY_REFERENCE_ERROR
}
func StaticSet_VariableReferenceCheck_EARLY_REFERENCE_ERROR(v *DiagnosticType) {
	clinit_VariableReferenceCheck()
	static_VariableReferenceCheck_EARLY_REFERENCE_ERROR = v
}

var static_VariableReferenceCheck_REASSIGNED_CONSTANT *DiagnosticType

func StaticGet_VariableReferenceCheck_REASSIGNED_CONSTANT() *DiagnosticType {
	clinit_VariableReferenceCheck()
	return static_VariableReferenceCheck_REASSIGNED_CONSTANT
}
func StaticSet_VariableReferenceCheck_REASSIGNED_CONSTANT(v *DiagnosticType) {
	clinit_VariableReferenceCheck()
	static_VariableReferenceCheck_REASSIGNED_CONSTANT = v
}

var static_VariableReferenceCheck_REDECLARED_VARIABLE_ERROR *DiagnosticType

func StaticGet_VariableReferenceCheck_REDECLARED_VARIABLE_ERROR() *DiagnosticType {
	clinit_VariableReferenceCheck()
	return static_VariableReferenceCheck_REDECLARED_VARIABLE_ERROR
}
func StaticSet_VariableReferenceCheck_REDECLARED_VARIABLE_ERROR(v *DiagnosticType) {
	clinit_VariableReferenceCheck()
	static_VariableReferenceCheck_REDECLARED_VARIABLE_ERROR = v
}

var static_VariableReferenceCheck_DECLARATION_NOT_DIRECTLY_IN_BLOCK *DiagnosticType

func StaticGet_VariableReferenceCheck_DECLARATION_NOT_DIRECTLY_IN_BLOCK() *DiagnosticType {
	clinit_VariableReferenceCheck()
	return static_VariableReferenceCheck_DECLARATION_NOT_DIRECTLY_IN_BLOCK
}
func StaticSet_VariableReferenceCheck_DECLARATION_NOT_DIRECTLY_IN_BLOCK(v *DiagnosticType) {
	clinit_VariableReferenceCheck()
	static_VariableReferenceCheck_DECLARATION_NOT_DIRECTLY_IN_BLOCK = v
}

var static_VariableReferenceCheck_UNUSED_LOCAL_ASSIGNMENT *DiagnosticType

func StaticGet_VariableReferenceCheck_UNUSED_LOCAL_ASSIGNMENT() *DiagnosticType {
	clinit_VariableReferenceCheck()
	return static_VariableReferenceCheck_UNUSED_LOCAL_ASSIGNMENT
}
func StaticSet_VariableReferenceCheck_UNUSED_LOCAL_ASSIGNMENT(v *DiagnosticType) {
	clinit_VariableReferenceCheck()
	static_VariableReferenceCheck_UNUSED_LOCAL_ASSIGNMENT = v
}

var static_VariableReferenceCheck_GOOG_REQUIRE pkg_rhino.QualifiedName

func StaticGet_VariableReferenceCheck_GOOG_REQUIRE() pkg_rhino.QualifiedName {
	clinit_VariableReferenceCheck()
	return static_VariableReferenceCheck_GOOG_REQUIRE
}
func StaticSet_VariableReferenceCheck_GOOG_REQUIRE(v pkg_rhino.QualifiedName) {
	clinit_VariableReferenceCheck()
	static_VariableReferenceCheck_GOOG_REQUIRE = v
}

var static_VariableReferenceCheck_GOOG_REQUIRE_TYPE pkg_rhino.QualifiedName

func StaticGet_VariableReferenceCheck_GOOG_REQUIRE_TYPE() pkg_rhino.QualifiedName {
	clinit_VariableReferenceCheck()
	return static_VariableReferenceCheck_GOOG_REQUIRE_TYPE
}
func StaticSet_VariableReferenceCheck_GOOG_REQUIRE_TYPE(v pkg_rhino.QualifiedName) {
	clinit_VariableReferenceCheck()
	static_VariableReferenceCheck_GOOG_REQUIRE_TYPE = v
}

var static_VariableReferenceCheck_GOOG_FORWARD_DECLARE pkg_rhino.QualifiedName

func StaticGet_VariableReferenceCheck_GOOG_FORWARD_DECLARE() pkg_rhino.QualifiedName {
	clinit_VariableReferenceCheck()
	return static_VariableReferenceCheck_GOOG_FORWARD_DECLARE
}
func StaticSet_VariableReferenceCheck_GOOG_FORWARD_DECLARE(v pkg_rhino.QualifiedName) {
	clinit_VariableReferenceCheck()
	static_VariableReferenceCheck_GOOG_FORWARD_DECLARE = v
}

var static_VariableReferenceCheck_BLOCKLESS_DECLARATION_FORBIDDEN_STATEMENTS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_VariableReferenceCheck_BLOCKLESS_DECLARATION_FORBIDDEN_STATEMENTS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_VariableReferenceCheck()
	return static_VariableReferenceCheck_BLOCKLESS_DECLARATION_FORBIDDEN_STATEMENTS
}
func StaticSet_VariableReferenceCheck_BLOCKLESS_DECLARATION_FORBIDDEN_STATEMENTS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_VariableReferenceCheck()
	static_VariableReferenceCheck_BLOCKLESS_DECLARATION_FORBIDDEN_STATEMENTS = v
}

var static_VariableReferenceCheck_GOOG_SCOPE pkg_rhino.QualifiedName

func StaticGet_VariableReferenceCheck_GOOG_SCOPE() pkg_rhino.QualifiedName {
	clinit_VariableReferenceCheck()
	return static_VariableReferenceCheck_GOOG_SCOPE
}
func StaticSet_VariableReferenceCheck_GOOG_SCOPE(v pkg_rhino.QualifiedName) {
	clinit_VariableReferenceCheck()
	static_VariableReferenceCheck_GOOG_SCOPE = v
}

var clinitDone_VariableReferenceCheck uint32
var clinitState_VariableReferenceCheck j2gort.ClinitState

func clinit_VariableReferenceCheck() {
	if atomic.LoadUint32(&clinitDone_VariableReferenceCheck) != 0 {
		return
	}
	clinitSlow_VariableReferenceCheck()
}

func clinitSlow_VariableReferenceCheck() {
	if !j2gort.EnterClinit(&clinitState_VariableReferenceCheck, &clinitDone_VariableReferenceCheck) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_VariableReferenceCheck, &clinitDone_VariableReferenceCheck)
	static_VariableReferenceCheck_EARLY_REFERENCE = DiagnosticType_Warning("JSC_REFERENCE_BEFORE_DECLARE", "Variable referenced before declaration: {0}")
	static_VariableReferenceCheck_EARLY_EXPORTS_REFERENCE = DiagnosticType_Error("JSC_EXPORTS_REFERENCE_BEFORE_ASSIGN", "Illegal reference to `exports` before assignment `exports = ...`")
	static_VariableReferenceCheck_REDECLARED_VARIABLE = DiagnosticType_Warning("JSC_REDECLARED_VARIABLE", "Redeclared variable: {0}")
	static_VariableReferenceCheck_EARLY_REFERENCE_ERROR = DiagnosticType_Error("JSC_REFERENCE_BEFORE_DECLARE_ERROR", "Illegal variable reference before declaration: {0}")
	static_VariableReferenceCheck_REASSIGNED_CONSTANT = DiagnosticType_Error("JSC_REASSIGNED_CONSTANT", "Constant reassigned: {0}")
	static_VariableReferenceCheck_REDECLARED_VARIABLE_ERROR = DiagnosticType_Error("JSC_REDECLARED_VARIABLE_ERROR", "Illegal redeclared variable: {0}")
	static_VariableReferenceCheck_DECLARATION_NOT_DIRECTLY_IN_BLOCK = DiagnosticType_Error("JSC_DECLARATION_NOT_DIRECTLY_IN_BLOCK", "Block-scoped declaration not directly within block: {0}")
	static_VariableReferenceCheck_UNUSED_LOCAL_ASSIGNMENT = DiagnosticType_Disabled("JSC_UNUSED_LOCAL_ASSIGNMENT", "Value assigned to local variable {0} is never read")
	static_VariableReferenceCheck_GOOG_REQUIRE = pkg_rhino.QualifiedName_Of("goog.require")
	static_VariableReferenceCheck_GOOG_REQUIRE_TYPE = pkg_rhino.QualifiedName_Of("goog.requireType")
	static_VariableReferenceCheck_GOOG_FORWARD_DECLARE = pkg_rhino.QualifiedName_Of("goog.forwardDeclare")
	static_VariableReferenceCheck_BLOCKLESS_DECLARATION_FORBIDDEN_STATEMENTS = j2gort.Ext_com_google_common_collect_Sets_ImmutableEnumSet_2(j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_IF()), []j2gort.Ext_java_lang_Enum{j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_FOR()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_FOR_IN()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_FOR_OF()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_FOR_AWAIT_OF()), j2gort.ToIface_Ext_java_lang_Enum(pkg_rhino.StaticGet_Token_WHILE())})
	static_VariableReferenceCheck_GOOG_SCOPE = pkg_rhino.QualifiedName_Of("goog.scope")
}

func (this_ *VariableReferenceCheck) initFields_VariableReferenceCheck() {
	_ = this_
	this_.blocksWithDeclarations = j2gort.NewOrderedSet[*BasicBlock]()
}

// Java constructor: VariableReferenceCheck(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewVariableReferenceCheck(compiler AbstractCompiler) *VariableReferenceCheck {
	clinit_VariableReferenceCheck()
	this_ := &VariableReferenceCheck{}
	InitVariableReferenceCheck(this_, compiler)
	return this_
}

func InitVariableReferenceCheck(this_ *VariableReferenceCheck, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_VariableReferenceCheck()
	this_.compiler = compiler
	this_.checkUnusedLocals = compiler.GetOptions().Enables(DiagnosticGroup_ForType(StaticGet_VariableReferenceCheck_UNUSED_LOCAL_ASSIGNMENT()))
}

// Java method: VariableReferenceCheck.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *VariableReferenceCheck) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	{
		clinit_ReferenceCollector()
		var val_1 ReferenceCollector
		InitReferenceCollector_3(&val_1, this_.compiler, NewVariableReferenceCheck_ReferenceCheckingBehavior(this_), NewSyntacticScopeCreator_1(this_.compiler))
		(&val_1).ownsScopeCreator_ = true
		defer func() {
			this_.blocksWithDeclarations.Clear()
			(&val_1).ReleaseSlabs_()
		}()

		(&val_1).Process_2(externs, root)
	}
}

// Java method: VariableReferenceCheck.lookForHoistedFunction(com.google.javascript.jscomp.ReferenceCollection references)
func (this_ *VariableReferenceCheck) priv_VariableReferenceCheck_lookForHoistedFunction_5852(references *ReferenceCollection) *Reference {
	_ = references
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = references.Size()
		_ = size
		for ; i < size; i = i + 1 {
			var reference *Reference = references.Get(i)
			_ = reference
			if reference.IsHoistedFunction() {
				this_.blocksWithDeclarations.Add(reference.GetBasicBlock())
				return reference
			}
		}
	}
	return nil
}

// Java method: VariableReferenceCheck.checkBlocklessDeclaration(com.google.javascript.jscomp.Var v, com.google.javascript.jscomp.Reference reference, com.google.javascript.rhino.Node referenceNode)
func (this_ *VariableReferenceCheck) priv_VariableReferenceCheck_checkBlocklessDeclaration_5853(v *Var, reference *Reference, referenceNode pkg_rhino.Node) {
	_ = v
	_ = reference
	_ = referenceNode
	if (!(reference.IsVarDeclaration()) && reference.GetGrandparent().IsAddedBlock()) && j2gort.ImmSetContainsTypedU8_(StaticGet_VariableReferenceCheck_BLOCKLESS_DECLARATION_FORBIDDEN_STATEMENTS(), uint8(reference.GetGrandparent().GetParent().GetToken())) {
		this_.compiler.Report_1(JSError_Make_3(referenceNode, StaticGet_VariableReferenceCheck_DECLARATION_NOT_DIRECTLY_IN_BLOCK(), []string{v.GetName()}))
	}
}

// Java method: VariableReferenceCheck.checkRedeclaration(com.google.javascript.jscomp.Var v, com.google.javascript.jscomp.Reference reference, com.google.javascript.rhino.Node referenceNode, com.google.javascript.jscomp.Reference hoistedFn, com.google.javascript.jscomp.BasicBlock basicBlock)
func (this_ *VariableReferenceCheck) priv_VariableReferenceCheck_checkRedeclaration_5854(v *Var, reference *Reference, referenceNode pkg_rhino.Node, hoistedFn *Reference, basicBlock *BasicBlock) bool {
	_ = v
	_ = reference
	_ = referenceNode
	_ = hoistedFn
	_ = basicBlock
	var letConstShadowsVar bool = (v.GetParentNode().IsVar() && (reference.IsLetDeclaration() || reference.IsConstDeclaration()))
	_ = letConstShadowsVar
	var isVarNodeSameAsReferenceNode bool = (v.GetNode() == reference.GetNode())
	_ = isVarNodeSameAsReferenceNode
	var shadowCatchVar bool = (v.GetParentNode().IsCatch() && !(isVarNodeSameAsReferenceNode))
	_ = shadowCatchVar
	if this_.priv_VariableReferenceCheck_isRedeclaration_5855(basicBlock) && !(VarCheck_HasDuplicateDeclarationSuppression(this_.compiler, referenceNode, v.GetNameNode())) {
		var diagnosticType *DiagnosticType = nil
		_ = diagnosticType
		var warningNode pkg_rhino.Node = referenceNode
		_ = warningNode
		var shadowParam bool = ((v.IsParam() && NodeUtil_IsBlockScopedDeclaration(referenceNode.GetParent())) && (j2gort.CastPtr[Scope](v.GetScope()) == j2gort.CastPtr[Scope](reference.GetScope().GetParent())))
		_ = shadowParam
		var isFunctionDecl bool = (((v.GetParentNode() != 0) && v.GetParentNode().IsFunction()) && (v.GetParentNode().GetFirstChild() == referenceNode))
		_ = isFunctionDecl
		if ((((((v.IsLet() || v.IsConst()) || v.IsClass()) || letConstShadowsVar) || shadowCatchVar) || shadowParam) || v.IsImport()) || isFunctionDecl {
			diagnosticType = StaticGet_VariableReferenceCheck_REDECLARED_VARIABLE_ERROR()
		} else if reference.GetNode().GetParent().IsCatch() {
			return false
		} else {
			if j2gort.CastPtr[Scope](v.GetScope()).IsGlobal() {
				diagnosticType = StaticGet_VarCheck_VAR_MULTIPLY_DECLARED_ERROR()
			} else {
				diagnosticType = StaticGet_VariableReferenceCheck_REDECLARED_VARIABLE()
			}
			if (isVarNodeSameAsReferenceNode && (hoistedFn != nil)) && func() bool {
				sr_2 := v.GetName()
				se_3 := hoistedFn.GetNode().GetString()
				return j2gort.CheckNotNullString(sr_2) == se_3
			}() {
				warningNode = hoistedFn.GetNode()
			}
		}
		this_.compiler.Report_1(JSError_Make_3(warningNode, diagnosticType, []string{v.GetName(), VariableReferenceCheck_LocationOf(v.GetNode())}))
		return true
	}
	if (letConstShadowsVar || shadowCatchVar) && (j2gort.CastPtr[Scope](v.GetScope()) == reference.GetScope()) {
		this_.compiler.Report_1(JSError_Make_3(referenceNode, StaticGet_VariableReferenceCheck_REDECLARED_VARIABLE_ERROR(), []string{v.GetName()}))
		return true
	}
	return false
}

// Java method: VariableReferenceCheck.isRedeclaration(com.google.javascript.jscomp.BasicBlock newDeclaration)
func (this_ *VariableReferenceCheck) priv_VariableReferenceCheck_isRedeclaration_5855(newDeclaration *BasicBlock) bool {
	_ = newDeclaration
	{
		cur_4 := this_.blocksWithDeclarations.Iter()
		for cur_4.HasNext() {

			previousDeclaration := cur_4.
				Next()
			_ = previousDeclaration
			if previousDeclaration.ProvablyExecutesBefore(newDeclaration) {
				return true
			}
		}
	}
	return false
}

// Java static method: VariableReferenceCheck.locationOf(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node n)
func VariableReferenceCheck_LocationOf(n pkg_rhino.Node) string {
	clinit_VariableReferenceCheck()
	_ = n
	if n == 0 {
		return "<unknown>"
	} else {
		return n.GetLocation()
	}
}

// Java method: VariableReferenceCheck.checkEarlyReference(com.google.javascript.jscomp.Var v, com.google.javascript.jscomp.Reference reference, com.google.javascript.rhino.Node referenceNode)
func (this_ *VariableReferenceCheck) priv_VariableReferenceCheck_checkEarlyReference_5856(v *Var, reference *Reference, referenceNode pkg_rhino.Node) bool {
	_ = v
	_ = reference
	_ = referenceNode
	if referenceNode.IsFromExterns() || v.IsImplicitGoogNamespace() {
		return false
	}
	if v.IsVar() {
		var curr pkg_rhino.Node = reference.GetParent()
		_ = curr
		for curr.IsOr() && (curr.GetParent().GetFirstChild() == curr) {
			curr = curr.GetParent()
		}
		if curr.IsName() && func() bool {
			sr_6 := curr.GetString()
			se_7 := v.GetName()
			return j2gort.CheckNotNullString(sr_6) == se_7
		}() {
			return false
		}
	}
	var referenceScopeRoot pkg_rhino.Node = reference.GetScope().GetRootNode_0()
	_ = referenceScopeRoot
	if referenceScopeRoot.IsMemberFieldDef() && !(referenceScopeRoot.IsStaticMember()) {
		return false
	}
	if reference.GetScope().HasSameContainerScope(ToIface_AbstractScope(j2gort.CastPtr[Scope](v.GetScope()))) && !(j2gort.CheckNotNullString(v.GetName()) == "goog") {
		this_.compiler.Report_1(JSError_Make_3(reference.GetNode(), func() *DiagnosticType {
			if v.IsGoogModuleExports() {
				return StaticGet_VariableReferenceCheck_EARLY_EXPORTS_REFERENCE()
			}
			return func() *DiagnosticType {
				if ((v.IsLet() || v.IsConst()) || v.IsClass()) || v.IsParam() {
					return StaticGet_VariableReferenceCheck_EARLY_REFERENCE_ERROR()
				}
				return StaticGet_VariableReferenceCheck_EARLY_REFERENCE()
			}()
		}(), []string{v.GetName()}))
		return true
	}
	return false
}

// Java method: VariableReferenceCheck.checkForUnusedLocalVar(com.google.javascript.jscomp.Var v, com.google.javascript.jscomp.Reference unusedAssignment)
func (this_ *VariableReferenceCheck) priv_VariableReferenceCheck_checkForUnusedLocalVar_5857(v *Var, unusedAssignment *Reference) {
	_ = v
	_ = unusedAssignment
	if !(v.IsLocal()) {
		return
	}
	var jsDoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(unusedAssignment.GetNode())
	_ = jsDoc
	if (jsDoc != nil) && jsDoc.HasTypedefType() {
		return
	}
	var inGoogScope bool = false
	_ = inGoogScope
	var s *Scope = j2gort.CastPtr[Scope](v.GetScope())
	_ = s
	if s.IsFunctionBlockScope() {
		var function pkg_rhino.Node = s.GetRootNode_0().GetParent()
		_ = function
		var callee pkg_rhino.Node = function.GetPrevious()
		_ = callee
		inGoogScope = ((callee != 0) && StaticGet_VariableReferenceCheck_GOOG_SCOPE().Matches(callee))
	}
	if inGoogScope {
		return
	}
	if s.IsModuleScope() {
		var statement pkg_rhino.Node = NodeUtil_GetEnclosingStatement(v.GetNode())
		_ = statement
		if NodeUtil_IsNameDeclaration(statement) {
			var lhs pkg_rhino.Node = statement.GetFirstChild()
			_ = lhs
			var rhs pkg_rhino.Node = lhs.GetFirstChild()
			_ = rhs
			if (rhs != 0) && (((NodeUtil_IsCallTo_2_Node_QualifiedName(rhs, StaticGet_VariableReferenceCheck_GOOG_FORWARD_DECLARE()) || NodeUtil_IsCallTo_2_Node_QualifiedName(rhs, StaticGet_VariableReferenceCheck_GOOG_REQUIRE_TYPE())) || NodeUtil_IsCallTo_2_Node_QualifiedName(rhs, StaticGet_VariableReferenceCheck_GOOG_REQUIRE())) || rhs.IsQualifiedName()) {
				return
			}
		}
	}
	this_.compiler.Report_1(JSError_Make_3(unusedAssignment.GetNode(), StaticGet_VariableReferenceCheck_UNUSED_LOCAL_ASSIGNMENT(), []string{v.GetName()}))
}

// Java class: com.google.javascript.jscomp.VariableReferenceCheck.ReferenceCheckingBehavior
type VariableReferenceCheck_ReferenceCheckingBehavior struct {
	outer_             *VariableReferenceCheck
	varsInFunctionBody j2gort.Ext_java_util_Set `gson:"varsInFunctionBody"`
	idHash_            atomic.Int32
}

func (this_ *VariableReferenceCheck_ReferenceCheckingBehavior) BindSelf_(s any) {
	_ = s
}
func (this_ *VariableReferenceCheck_ReferenceCheckingBehavior) IsJavaNil_() bool { return this_ == nil }
func (*VariableReferenceCheck_ReferenceCheckingBehavior) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.VariableReferenceCheck$ReferenceCheckingBehavior", "ReferenceCheckingBehavior"
}
func (this_ *VariableReferenceCheck_ReferenceCheckingBehavior) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*VariableReferenceCheck_ReferenceCheckingBehavior) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.VariableReferenceCheck.ReferenceCheckingBehavior", "com.google.javascript.jscomp.VariableReferenceCheck$ReferenceCheckingBehavior", func(v_ any) bool {
		return j2gort.InstanceOf[*VariableReferenceCheck_ReferenceCheckingBehavior](v_)
	})
}
func (this_ *VariableReferenceCheck_ReferenceCheckingBehavior) SetOuter_(o *VariableReferenceCheck) {
	this_.outer_ = o
}
func (this_ *VariableReferenceCheck_ReferenceCheckingBehavior) GetOuter_() *VariableReferenceCheck {
	return this_.outer_
}
func (*VariableReferenceCheck_ReferenceCheckingBehavior) Implements_jscomp_ReferenceCollector_Behavior() {
}

func (this_ *VariableReferenceCheck_ReferenceCheckingBehavior) initFields_VariableReferenceCheck_ReferenceCheckingBehavior() {
	_ = this_
}

// Java constructor: ReferenceCheckingBehavior()
func NewVariableReferenceCheck_ReferenceCheckingBehavior(outer_ *VariableReferenceCheck) *VariableReferenceCheck_ReferenceCheckingBehavior {
	this_ := &VariableReferenceCheck_ReferenceCheckingBehavior{}
	this_.outer_ = outer_
	InitVariableReferenceCheck_ReferenceCheckingBehavior(this_)
	return this_
}

func InitVariableReferenceCheck_ReferenceCheckingBehavior(this_ *VariableReferenceCheck_ReferenceCheckingBehavior) {
	this_.initFields_VariableReferenceCheck_ReferenceCheckingBehavior()
	this_.varsInFunctionBody = j2gort.NewExt_java_util_LinkedHashSet()
}

// Java method: ReferenceCheckingBehavior.afterExitScope(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.jscomp.ReferenceMap referenceMap)
func (this_ *VariableReferenceCheck_ReferenceCheckingBehavior) AfterExitScope(t *NodeTraversal, referenceMap ReferenceMap) {
	_ = t
	_ = referenceMap
	var scope *Scope = t.GetScope()
	_ = scope
	if scope.IsFunctionBlockScope() {
		this_.varsInFunctionBody.Clear()
		{
			cur_8 := (scope).iterVars_()
			for cur_8.HasNext() {

				v := cur_8.Next()

				_ = v
				j2gort.SetAddStr(this_.varsInFunctionBody, v.GetName())
			}
		}
	}
	{
		cur_10 := (scope).iterVars_()
		for cur_10.HasNext() {

			v_2 := cur_10.Next()

			_ = v_2
			var referenceCollection *ReferenceCollection = j2gort.CastPtr[ReferenceCollection](referenceMap.GetReferences(j2gort.PtrToAny(v_2)))
			_ = referenceCollection
			if referenceCollection != nil {
				if scope.GetRootNode_0().IsFunction() && v_2.IsDefaultParam() {
					this_.priv_ReferenceCheckingBehavior_checkDefaultParam_5858(v_2, scope, this_.varsInFunctionBody)
				}
				if scope.GetRootNode_0().IsFunction() {
					this_.priv_ReferenceCheckingBehavior_checkShadowParam_5859(v_2, scope, referenceCollection)
				}
				this_.priv_ReferenceCheckingBehavior_checkVar_5861(v_2, referenceCollection)
			}
		}
	}
	if scope.HasOwnImplicitSlot(StaticGet_AbstractScope_ImplicitVar_EXPORTS()) {
		this_.priv_ReferenceCheckingBehavior_checkGoogModuleExports_5860(j2gort.CastPtr[Var](scope.MakeImplicitVar(StaticGet_AbstractScope_ImplicitVar_EXPORTS())), referenceMap)
	}
}

// Java method: ReferenceCheckingBehavior.checkDefaultParam(com.google.javascript.jscomp.Var param, com.google.javascript.jscomp.Scope scope, java.util.Set<java.lang.String> varsInFunctionBody)
func (this_ *VariableReferenceCheck_ReferenceCheckingBehavior) priv_ReferenceCheckingBehavior_checkDefaultParam_5858(param *Var, scope *Scope, varsInFunctionBody j2gort.Ext_java_util_Set) {
	_ = param
	_ = scope
	_ = varsInFunctionBody
	NodeTraversal_Traverse(this_.outer_.compiler, param.GetParentNode().GetSecondChild(), NewVariableReferenceCheck__Anon_99(this_, varsInFunctionBody, scope))
}

// Java method: ReferenceCheckingBehavior.checkShadowParam(com.google.javascript.jscomp.Var v, com.google.javascript.jscomp.Scope functionScope, com.google.javascript.jscomp.ReferenceCollection references)
func (this_ *VariableReferenceCheck_ReferenceCheckingBehavior) priv_ReferenceCheckingBehavior_checkShadowParam_5859(v *Var, functionScope *Scope, references *ReferenceCollection) {
	_ = v
	_ = functionScope
	_ = references
	var maybeParam *Var = j2gort.CastPtr[Var](functionScope.GetVar(v.GetName()))
	_ = maybeParam
	if ((maybeParam != nil) && maybeParam.IsParam()) && (j2gort.CastPtr[Scope](maybeParam.GetScope()) == functionScope) {
		{
			var i int32 = int32(0)
			_ = i
			var size int32 = references.Size()
			_ = size
			for ; i < size; i = i + 1 {
				var r *Reference = references.Get(i)
				_ = r
				if (r.IsVarDeclaration() || r.IsHoistedFunction()) && (r.GetNode() != v.GetNameNode()) {
					this_.outer_.compiler.Report_1(JSError_Make_3(r.GetNode(), StaticGet_VariableReferenceCheck_REDECLARED_VARIABLE(), []string{v.GetName()}))
				}
			}
		}
	}
}

// Java method: ReferenceCheckingBehavior.checkGoogModuleExports(com.google.javascript.jscomp.Var exportsVar, com.google.javascript.jscomp.ReferenceMap referenceMap)
func (this_ *VariableReferenceCheck_ReferenceCheckingBehavior) priv_ReferenceCheckingBehavior_checkGoogModuleExports_5860(exportsVar *Var, referenceMap ReferenceMap) {
	_ = exportsVar
	_ = referenceMap
	var references *ReferenceCollection = j2gort.CastPtr[ReferenceCollection](referenceMap.GetReferences(j2gort.PtrToAny(exportsVar)))
	_ = references
	if (references == nil) || references.IsNeverAssigned() {
		return
	}
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = references.Size()
		_ = size
		for ; i < size; i = i + 1 {
			var reference *Reference = references.Get(i)
			_ = reference
			if reference.IsLvalue() {
				break
			}
			this_.outer_.priv_VariableReferenceCheck_checkEarlyReference_5856(exportsVar, reference, reference.GetNode())
		}
	}
}

// Java method: ReferenceCheckingBehavior.checkVar(com.google.javascript.jscomp.Var v, com.google.javascript.jscomp.ReferenceCollection references)
func (this_ *VariableReferenceCheck_ReferenceCheckingBehavior) priv_ReferenceCheckingBehavior_checkVar_5861(v *Var, references *ReferenceCollection) {
	_ = v
	_ = references
	this_.outer_.blocksWithDeclarations.Clear()
	var hasSeenDeclaration bool = false
	_ = hasSeenDeclaration
	var hasErrors bool = false
	_ = hasErrors
	var isRead bool = false
	_ = isRead
	var unusedAssignment *Reference = nil
	_ = unusedAssignment
	var hoistedFn *Reference = this_.outer_.priv_VariableReferenceCheck_lookForHoistedFunction_5852(references)
	_ = hoistedFn
	if hoistedFn != nil {
		hasSeenDeclaration = true
	}
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = references.Size()
		_ = size
		for ; i < size; i = i + 1 {
			var reference *Reference = references.Get(i)
			_ = reference
			if reference == hoistedFn {
				continue
			}
			var referenceNode pkg_rhino.Node = reference.GetNode()
			_ = referenceNode
			var basicBlock *BasicBlock = reference.GetBasicBlock()
			_ = basicBlock
			var isDeclaration bool = reference.IsDeclaration()
			_ = isDeclaration
			var isAssignment bool = (isDeclaration || reference.IsLvalue())
			_ = isAssignment
			if isDeclaration {
				hasSeenDeclaration = true
				hasErrors = this_.outer_.priv_VariableReferenceCheck_checkRedeclaration_5854(v, reference, referenceNode, hoistedFn, basicBlock)
				this_.outer_.blocksWithDeclarations.Add(basicBlock)
				this_.outer_.priv_VariableReferenceCheck_checkBlocklessDeclaration_5853(v, reference, referenceNode)
				if reference.GetGrandparent().IsExport() {
					isRead = true
				}
			} else {
				if !(hasSeenDeclaration) {
					hasErrors = this_.outer_.priv_VariableReferenceCheck_checkEarlyReference_5856(v, reference, referenceNode)
				}
				if (!(hasErrors) && v.IsConst()) && reference.IsLvalue() {
					this_.outer_.compiler.Report_1(JSError_Make_3(referenceNode, StaticGet_VariableReferenceCheck_REASSIGNED_CONSTANT(), []string{v.GetName()}))
				}
				if ((v.IsLet() || v.IsConst()) && (j2gort.CastPtr[Scope](v.GetScope()) == reference.GetScope())) && NodeUtil_IsEnhancedFor(reference.GetScope().GetRootNode_0()) {
					this_.outer_.compiler.Report_1(JSError_Make_3(referenceNode, StaticGet_VariableReferenceCheck_EARLY_REFERENCE_ERROR(), []string{v.GetName()}))
				}
			}
			if isAssignment {
				var decl *Reference = references.Get(int32(0))
				_ = decl
				var declNode pkg_rhino.Node = decl.GetNode()
				_ = declNode
				var gp pkg_rhino.Node = declNode.GetGrandparent()
				_ = gp
				var lhsOfForInLoop bool = (gp.IsForIn() && (gp.GetFirstFirstChild() == declNode))
				_ = lhsOfForInLoop
				if ((decl.GetScope().IsLocal() && ((decl.IsVarDeclaration() || decl.IsLetDeclaration()) || decl.IsConstDeclaration())) && !(decl.GetNode().IsFromExterns())) && !(lhsOfForInLoop) {
					unusedAssignment = reference
				}
				if ((reference.GetParent().IsDec() || reference.GetParent().IsInc()) || NodeUtil_IsCompoundAssignmentOp(reference.GetParent())) && NodeUtil_IsExpressionResultUsed(reference.GetNode()) {
					isRead = true
				}
			} else {
				isRead = true
			}
		}
	}
	if ((this_.outer_.checkUnusedLocals && (unusedAssignment != nil)) && !(isRead)) && !(hasErrors) {
		this_.outer_.priv_VariableReferenceCheck_checkForUnusedLocalVar_5857(v, unusedAssignment)
	}
}

func (this_ *VariableReferenceCheck_ReferenceCheckingBehavior) BeforeTraverseNode(p0_ *NodeTraversal, p1_ pkg_rhino.Node, p2_ pkg_rhino.Node) {
	IfaceDefault_ReferenceCollector_Behavior_BeforeTraverseNode(this_, p0_, p1_, p2_)
}

func (this_ *VariableReferenceCheck_ReferenceCheckingBehavior) AfterVisitNode(p0_ *NodeTraversal, p1_ pkg_rhino.Node, p2_ pkg_rhino.Node, p3_ *Reference) {
	IfaceDefault_ReferenceCollector_Behavior_AfterVisitNode(this_, p0_, p1_, p2_, p3_)
}

// Java class:
type VariableReferenceCheck__Anon_99 struct {
	NodeTraversal_AbstractShallowCallbackBase
	outer_                   *VariableReferenceCheck_ReferenceCheckingBehavior
	cap_0_varsInFunctionBody j2gort.Ext_java_util_Set
	cap_1_scope              *Scope
}

func (this_ *VariableReferenceCheck__Anon_99) BindSelf_(s any) {
	_ = s
}
func (this_ *VariableReferenceCheck__Anon_99) IsJavaNil_() bool { return this_ == nil }
func (*VariableReferenceCheck__Anon_99) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.VariableReferenceCheck$ReferenceCheckingBehavior$1", ""
}
func (this_ *VariableReferenceCheck__Anon_99) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractShallowCallbackBase.IdentityHash_()
}
func (*VariableReferenceCheck__Anon_99) HasIdentityEquals_() {}
func (this_ *VariableReferenceCheck__Anon_99) SetOuter_(o *VariableReferenceCheck_ReferenceCheckingBehavior) {
	this_.outer_ = o
}
func (this_ *VariableReferenceCheck__Anon_99) GetOuter_() *VariableReferenceCheck_ReferenceCheckingBehavior {
	return this_.outer_
}
func (*VariableReferenceCheck__Anon_99) Implements_jscomp_NodeTraversal_AbstractShallowCallback() {}
func (*VariableReferenceCheck__Anon_99) Implements_jscomp_NodeTraversal_Callback()                {}

func (this_ *VariableReferenceCheck__Anon_99) initFields_VariableReferenceCheck__Anon_99() {
	_ = this_
}

func NewVariableReferenceCheck__Anon_99(outer_ *VariableReferenceCheck_ReferenceCheckingBehavior, cap_p_varsInFunctionBody j2gort.Ext_java_util_Set, cap_p_scope *Scope) *VariableReferenceCheck__Anon_99 {
	this_ := &VariableReferenceCheck__Anon_99{}
	this_.outer_ = outer_
	this_.cap_0_varsInFunctionBody = cap_p_varsInFunctionBody
	this_.cap_1_scope = cap_p_scope
	InitNodeTraversal_AbstractShallowCallback(&this_.NodeTraversal_AbstractShallowCallbackBase)
	this_.initFields_VariableReferenceCheck__Anon_99()
	return this_
}

// Java method: .visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *VariableReferenceCheck__Anon_99) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(NodeUtil_IsReferenceName(n)) {
		return
	}
	var refName string = n.GetString()
	_ = refName
	if j2gort.SetContainsStr(this_.cap_0_varsInFunctionBody, refName) && !(this_.cap_1_scope.HasSlot(refName)) {
		this_.outer_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_VariableReferenceCheck_EARLY_REFERENCE_ERROR(), []string{refName}))
	}
}
func (*VariableReferenceCheck__Anon_99) CanReleaseSyntacticScopes_() {}
