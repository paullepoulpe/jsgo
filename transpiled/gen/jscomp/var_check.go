/*
 * Copyright 2004 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/VarCheck.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.VarCheck
type VarCheck struct {
	undefinedNamesFromExterns j2gort.OrderedSet[string]           `gson:"undefinedNamesFromExterns"`
	compiler                  AbstractCompiler                    `gson:"compiler"`
	dupHandler                *VarCheck_RedeclarationCheckHandler `gson:"dupHandler"`
	idHash_                   atomic.Int32
	validityCheck             bool `gson:"validityCheck"`
	strictExternCheck         bool `gson:"strictExternCheck"`
}

func (this_ *VarCheck) BindSelf_(s any) {
	_ = s
}
func (this_ *VarCheck) IsJavaNil_() bool { return this_ == nil }
func (*VarCheck) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.VarCheck", "VarCheck"
}
func (this_ *VarCheck) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*VarCheck) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.VarCheck", "com.google.javascript.jscomp.VarCheck", func(v_ any) bool {
		return j2gort.InstanceOf[*VarCheck](v_)
	})
}
func (*VarCheck) Implements_jscomp_NodeTraversal_ScopedCallback() {}
func (*VarCheck) Implements_jscomp_CompilerPass()                 {}
func (*VarCheck) Implements_jscomp_NodeTraversal_Callback()       {}

var static_VarCheck_UNDEFINED_VAR_ERROR *DiagnosticType

func StaticGet_VarCheck_UNDEFINED_VAR_ERROR() *DiagnosticType {
	clinit_VarCheck()
	return static_VarCheck_UNDEFINED_VAR_ERROR
}
func StaticSet_VarCheck_UNDEFINED_VAR_ERROR(v *DiagnosticType) {
	clinit_VarCheck()
	static_VarCheck_UNDEFINED_VAR_ERROR = v
}

var static_VarCheck_VIOLATED_CHUNK_DEP_ERROR *DiagnosticType

func StaticGet_VarCheck_VIOLATED_CHUNK_DEP_ERROR() *DiagnosticType {
	clinit_VarCheck()
	return static_VarCheck_VIOLATED_CHUNK_DEP_ERROR
}
func StaticSet_VarCheck_VIOLATED_CHUNK_DEP_ERROR(v *DiagnosticType) {
	clinit_VarCheck()
	static_VarCheck_VIOLATED_CHUNK_DEP_ERROR = v
}

var static_VarCheck_MISSING_CHUNK_DEP_ERROR *DiagnosticType

func StaticGet_VarCheck_MISSING_CHUNK_DEP_ERROR() *DiagnosticType {
	clinit_VarCheck()
	return static_VarCheck_MISSING_CHUNK_DEP_ERROR
}
func StaticSet_VarCheck_MISSING_CHUNK_DEP_ERROR(v *DiagnosticType) {
	clinit_VarCheck()
	static_VarCheck_MISSING_CHUNK_DEP_ERROR = v
}

var static_VarCheck_STRICT_CHUNK_DEP_ERROR *DiagnosticType

func StaticGet_VarCheck_STRICT_CHUNK_DEP_ERROR() *DiagnosticType {
	clinit_VarCheck()
	return static_VarCheck_STRICT_CHUNK_DEP_ERROR
}
func StaticSet_VarCheck_STRICT_CHUNK_DEP_ERROR(v *DiagnosticType) {
	clinit_VarCheck()
	static_VarCheck_STRICT_CHUNK_DEP_ERROR = v
}

var static_VarCheck_NAME_REFERENCE_IN_EXTERNS_ERROR *DiagnosticType

func StaticGet_VarCheck_NAME_REFERENCE_IN_EXTERNS_ERROR() *DiagnosticType {
	clinit_VarCheck()
	return static_VarCheck_NAME_REFERENCE_IN_EXTERNS_ERROR
}
func StaticSet_VarCheck_NAME_REFERENCE_IN_EXTERNS_ERROR(v *DiagnosticType) {
	clinit_VarCheck()
	static_VarCheck_NAME_REFERENCE_IN_EXTERNS_ERROR = v
}

var static_VarCheck_UNDEFINED_EXTERN_VAR_ERROR *DiagnosticType

func StaticGet_VarCheck_UNDEFINED_EXTERN_VAR_ERROR() *DiagnosticType {
	clinit_VarCheck()
	return static_VarCheck_UNDEFINED_EXTERN_VAR_ERROR
}
func StaticSet_VarCheck_UNDEFINED_EXTERN_VAR_ERROR(v *DiagnosticType) {
	clinit_VarCheck()
	static_VarCheck_UNDEFINED_EXTERN_VAR_ERROR = v
}

var static_VarCheck_VAR_MULTIPLY_DECLARED_ERROR *DiagnosticType

func StaticGet_VarCheck_VAR_MULTIPLY_DECLARED_ERROR() *DiagnosticType {
	clinit_VarCheck()
	return static_VarCheck_VAR_MULTIPLY_DECLARED_ERROR
}
func StaticSet_VarCheck_VAR_MULTIPLY_DECLARED_ERROR(v *DiagnosticType) {
	clinit_VarCheck()
	static_VarCheck_VAR_MULTIPLY_DECLARED_ERROR = v
}

var static_VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR *DiagnosticType

func StaticGet_VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR() *DiagnosticType {
	clinit_VarCheck()
	return static_VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR
}
func StaticSet_VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR(v *DiagnosticType) {
	clinit_VarCheck()
	static_VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR = v
}

var static_VarCheck_VAR_ARGUMENTS_SHADOWED_ERROR *DiagnosticType

func StaticGet_VarCheck_VAR_ARGUMENTS_SHADOWED_ERROR() *DiagnosticType {
	clinit_VarCheck()
	return static_VarCheck_VAR_ARGUMENTS_SHADOWED_ERROR
}
func StaticSet_VarCheck_VAR_ARGUMENTS_SHADOWED_ERROR(v *DiagnosticType) {
	clinit_VarCheck()
	static_VarCheck_VAR_ARGUMENTS_SHADOWED_ERROR = v
}
func StaticGet_VarCheck_ARGUMENTS() string { return "arguments" }
func StaticGet_VarCheck_JSCOMP_CLOSURE_UNAWARE_CODE_SHADOW_HOST_NAME() string {
	return "$jscomp_wrap_closure_unaware_code"
}

var static_VarCheck_googForwardDeclare pkg_rhino.Node

func StaticGet_VarCheck_GoogForwardDeclare() pkg_rhino.Node {
	clinit_VarCheck()
	return static_VarCheck_googForwardDeclare
}
func StaticSet_VarCheck_GoogForwardDeclare(v pkg_rhino.Node) {
	clinit_VarCheck()
	static_VarCheck_googForwardDeclare = v
}

var static_VarCheck_REQUIRED_SYMBOLS j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_VarCheck_REQUIRED_SYMBOLS() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_VarCheck()
	return static_VarCheck_REQUIRED_SYMBOLS
}
func StaticSet_VarCheck_REQUIRED_SYMBOLS(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_VarCheck()
	static_VarCheck_REQUIRED_SYMBOLS = v
}

var clinitDone_VarCheck uint32
var clinitState_VarCheck j2gort.ClinitState

func clinit_VarCheck() {
	if atomic.LoadUint32(&clinitDone_VarCheck) != 0 {
		return
	}
	clinitSlow_VarCheck()
}

func clinitSlow_VarCheck() {
	if !j2gort.EnterClinit(&clinitState_VarCheck, &clinitDone_VarCheck) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_VarCheck, &clinitDone_VarCheck)
	static_VarCheck_UNDEFINED_VAR_ERROR = DiagnosticType_Error("JSC_UNDEFINED_VARIABLE", "variable {0} is undeclared")
	static_VarCheck_VIOLATED_CHUNK_DEP_ERROR = DiagnosticType_Error("JSC_VIOLATED_CHUNK_DEPENDENCY", "chunk {0} cannot reference {2}, defined in chunk {1}, since {1} loads after {0}")
	static_VarCheck_MISSING_CHUNK_DEP_ERROR = DiagnosticType_Warning("JSC_MISSING_CHUNK_DEPENDENCY", "missing chunk dependency; chunk {0} should depend on chunk {1} because it references {2}")
	static_VarCheck_STRICT_CHUNK_DEP_ERROR = DiagnosticType_Disabled("JSC_STRICT_CHUNK_DEPENDENCY", "cannot reference {2} because of a missing chunk dependency\ndefined in chunk {1}, referenced from chunk {0}")
	static_VarCheck_NAME_REFERENCE_IN_EXTERNS_ERROR = DiagnosticType_Warning("JSC_NAME_REFERENCE_IN_EXTERNS", "accessing name {0} in externs has no effect. Perhaps you forgot to add a var keyword?")
	static_VarCheck_UNDEFINED_EXTERN_VAR_ERROR = DiagnosticType_Warning("JSC_UNDEFINED_EXTERN_VAR_ERROR", "name {0} is not defined in the externs.")
	static_VarCheck_VAR_MULTIPLY_DECLARED_ERROR = DiagnosticType_Error("JSC_VAR_MULTIPLY_DECLARED_ERROR", "Variable {0} declared more than once. First occurrence: {1}")
	static_VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR = DiagnosticType_Error("JSC_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR", "Block-scoped variable {0} declared more than once. First occurrence: {1}")
	static_VarCheck_VAR_ARGUMENTS_SHADOWED_ERROR = DiagnosticType_Error("JSC_VAR_ARGUMENTS_SHADOWED_ERROR", "Shadowing \"arguments\" is not allowed")
	static_VarCheck_googForwardDeclare = pkg_rhino.IR_Getprop_2(pkg_rhino.IR_Name("goog"), "forwardDeclare")
	static_VarCheck_REQUIRED_SYMBOLS = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_7("AggregateError", "Array", "Error", "Float32Array", "Function", "Infinity", []any{"Iterator", "JSCOMPILER_PRESERVE", "JSCompiler_renameProperty", "Map", "Math", "NaN", "Number", "Object", "Promise", "RangeError", "ReferenceError", "Reflect", "RegExp", "Set", "String", "SuppressedError", "Symbol", "TypeError", "WeakMap", "global", "globalThis", "isNaN", "parseFloat", "parseInt", "self", "undefined", "window"})
}

func (this_ *VarCheck) initFields_VarCheck() {
	_ = this_
	this_.undefinedNamesFromExterns = j2gort.NewOrderedSet[string]()
}

// Java constructor: VarCheck(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewVarCheck_1(compiler AbstractCompiler) *VarCheck {
	clinit_VarCheck()
	this_ := &VarCheck{}
	InitVarCheck_1(this_, compiler)
	return this_
}

func InitVarCheck_1(this_ *VarCheck, compiler AbstractCompiler) {
	_ = compiler
	InitVarCheck_2(this_, compiler, false)
}

// Java constructor: VarCheck(com.google.javascript.jscomp.AbstractCompiler compiler, boolean validityCheck)
func NewVarCheck_2(compiler AbstractCompiler, validityCheck bool) *VarCheck {
	clinit_VarCheck()
	this_ := &VarCheck{}
	InitVarCheck_2(this_, compiler, validityCheck)
	return this_
}

func InitVarCheck_2(this_ *VarCheck, compiler AbstractCompiler, validityCheck bool) {
	_ = compiler
	_ = validityCheck
	this_.initFields_VarCheck()
	this_.compiler = compiler
	this_.strictExternCheck = (compiler.GetErrorLevel(JSError_Make_5("", int32(0), int32(0), StaticGet_VarCheck_UNDEFINED_EXTERN_VAR_ERROR(), []string{})) == StaticGet_CheckLevel_ERROR())
	this_.validityCheck = validityCheck
}

// Java method: VarCheck.createScopeCreator()
func (this_ *VarCheck) priv_VarCheck_createScopeCreator_5878() *SyntacticScopeCreator {
	if this_.validityCheck {
		return NewSyntacticScopeCreator_1(this_.compiler)
	} else {
		this_.dupHandler = NewVarCheck_RedeclarationCheckHandler(this_)
		return NewSyntacticScopeCreator_2(this_.compiler, ToIface_SyntacticScopeCreator_RedeclarationHandler(this_.dupHandler))
	}
}

// Java method: VarCheck.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *VarCheck) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var scopeCreator ScopeCreator = ToIface_ScopeCreator(this_.priv_VarCheck_createScopeCreator_5878())
	_ = scopeCreator
	if r_, ok_ := any(scopeCreator).(interface{ ReleaseSlabs_() }); ok_ {
		defer r_.ReleaseSlabs_()
	}

	if !(this_.validityCheck) {
		{
			clinit_NodeTraversal()
			var ntb_ NodeTraversal_Builder
			InitNodeTraversal_Builder(&ntb_)
			(&ntb_).SetCompiler(this_.compiler)
			(&ntb_).SetCallback_1_NodeTraversal_Callback(NewVarCheck_NameRefInExternsCheck(this_))
			(&ntb_).SetScopeCreator(scopeCreator)
			(&ntb_).Traverse(externs)
		}
	}
	{
		clinit_NodeTraversal()
		var ntb_1_ NodeTraversal_Builder
		InitNodeTraversal_Builder(&ntb_1_)
		(&ntb_1_).SetCompiler(this_.compiler)
		(&ntb_1_).SetCallback_1_NodeTraversal_Callback(this_)
		(&ntb_1_).SetScopeCreator(scopeCreator)
		(&ntb_1_).TraverseRoots(externs, root)
	}
	{
		cur_1 := this_.undefinedNamesFromExterns.Iter()
		for cur_1.HasNext() {
			varName := cur_1.
				Next()
			_ = varName
			this_.priv_VarCheck_createSynthesizedExternVar_5881(varName, false)
		}
	}
	if this_.dupHandler != nil {
		this_.dupHandler.RemoveDuplicates()
	}
}

// Java method: VarCheck.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *VarCheck) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return true
}

// Java method: VarCheck.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *VarCheck) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsName() {
		this_.priv_VarCheck_checkName_5879(t, n, parent)
	}
}

// Java method: VarCheck.checkName(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *VarCheck) priv_VarCheck_checkName_5879(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	var varName string = n.GetString()
	_ = varName
	if j2gort.CheckNotNullString(varName) == "" {
		if !(NodeUtil_IsFunctionExpression(parent) || NodeUtil_IsMethodDeclaration(parent)) {
			j2gort.CheckState(false)
		}
		return
	}
	var scope *Scope = t.GetScope()
	_ = scope
	var var_ *Var = j2gort.CastPtr[Var](scope.GetVar(varName))
	_ = var_
	if var_ == nil {
		if (NodeUtil_IsFunctionExpression(parent) || NodeUtil_IsClassExpression(parent)) && n.IsFirstChildOf(parent) {
			return
		}
		if NodeUtil_IsNonlocalModuleExportName(n) {
			return
		}
		this_.priv_VarCheck_handleUndeclaredVariableRef_5880(t, n)
		j2gort.CastPtr[Scope](scope.GetGlobalScope()).Declare(varName, n, this_.compiler.GetSynthesizedExternsInput())
		return
	}
	if (var_.IsImplicitGoogNamespace() && func() bool {
		re_2 := var_.GetImplicitGoogNamespaceStrength()
		ae_3 := pkg_rhino.StaticGet_StaticSourceFile_SourceKind_WEAK()
		if re_2 == 0 {
			j2gort.ThrowNPE()
		}
		return re_2 == ae_3
	}()) && func() bool {
		re_4 := VarCheck_StrengthOf(n)
		ae_5 := pkg_rhino.StaticGet_StaticSourceFile_SourceKind_STRONG()
		if re_4 == 0 {
			j2gort.ThrowNPE()
		}
		return re_4 == ae_5
	}() {
		this_.priv_VarCheck_handleUndeclaredVariableRef_5880(t, n)
		return
	}
	var currInput *CompilerInput = t.GetInput()
	_ = currInput
	var varInput *CompilerInput = var_.GetInput()
	_ = varInput
	if ((currInput == varInput) || (currInput == nil)) || (varInput == nil) {
		return
	}
	var currChunk *JSChunk = currInput.GetChunk()
	_ = currChunk
	var varChunk *JSChunk = varInput.GetChunk()
	_ = varChunk
	var chunkGraph *JSChunkGraph = this_.compiler.GetChunkGraph()
	_ = chunkGraph
	if ((!(this_.validityCheck) && (varChunk != currChunk)) && (varChunk != nil)) && (currChunk != nil) {
		if varChunk.IsWeak() {
			this_.priv_VarCheck_handleUndeclaredVariableRef_5880(t, n)
		}
		if chunkGraph.DependsOn(currChunk, varChunk) {
		} else {
			if scope.IsGlobal() {
				if chunkGraph.DependsOn(varChunk, currChunk) {
					t.Report_3(n, StaticGet_VarCheck_VIOLATED_CHUNK_DEP_ERROR(), []string{currChunk.GetName(), varChunk.GetName(), varName})
				} else {
					t.Report_3(n, StaticGet_VarCheck_MISSING_CHUNK_DEP_ERROR(), []string{currChunk.GetName(), varChunk.GetName(), varName})
				}
			} else {
				t.Report_3(n, StaticGet_VarCheck_STRICT_CHUNK_DEP_ERROR(), []string{currChunk.GetName(), varChunk.GetName(), varName})
			}
		}
	}
}

// Java static method: VarCheck.strengthOf(com.google.javascript.rhino.Node n)
func VarCheck_StrengthOf(n pkg_rhino.Node) pkg_rhino.StaticSourceFile_SourceKind {
	clinit_VarCheck()
	_ = n
	var source pkg_rhino.StaticSourceFile = n.GetStaticSourceFile()
	_ = source
	if source == nil || source.IsJavaNil_() {
		return pkg_rhino.StaticGet_StaticSourceFile_SourceKind_EXTERN()
	}
	return source.GetKind()
}

// Java method: VarCheck.handleUndeclaredVariableRef(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *VarCheck) priv_VarCheck_handleUndeclaredVariableRef_5880(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if !(n.IsName()) {
		j2gort.CheckState(false)
	}
	var varName string = n.GetString()
	_ = varName
	if n.GetParent().IsTypeOf() {
	} else {
		var iifeCond_1 bool
		{

			sr_6 := varName
			se_7 := "$jscomp_wrap_closure_unaware_code"
			iifeCond_1 = j2gort.CheckNotNullString(sr_6) == se_7
		}
		if iifeCond_1 {
		} else if this_.strictExternCheck && t.GetInput().IsExtern() {
		} else {
			t.Report_3(n, StaticGet_VarCheck_UNDEFINED_VAR_ERROR(), []string{varName})
		}
	}
	if this_.validityCheck {
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Unexpected variable "+j2gort.StrValString(varName))))
	} else if !(this_.undefinedNamesFromExterns.Contains(varName)) {
		this_.priv_VarCheck_createSynthesizedExternVar_5881(varName, true)
	}
}

// Java method: VarCheck.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *VarCheck) EnterScope(t *NodeTraversal) {
	_ = t
}

// Java method: VarCheck.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *VarCheck) ExitScope(t *NodeTraversal) {
	_ = t
	if !(this_.validityCheck) && t.InGlobalScope() {
		var scope *Scope = t.GetScope()
		_ = scope
		{
			cur_8 := j2gort.NewIterCursor(j2gort.CheckNotNull(StaticGet_VarCheck_REQUIRED_SYMBOLS()))
			for cur_8.HasNext() {
				requiredSymbol := cur_8.NextStr()
				_ = requiredSymbol
				var var_ *Var = j2gort.CastPtr[Var](scope.GetVar(requiredSymbol))
				_ = var_
				if var_ == nil {
					this_.undefinedNamesFromExterns.Add(requiredSymbol)
				}
			}
		}
	}
}

// Java method: VarCheck.createSynthesizedExternVar(java.lang.String varName, boolean isFromUndefinedCodeRef)
func (this_ *VarCheck) priv_VarCheck_createSynthesizedExternVar_5881(varName string, isFromUndefinedCodeRef bool) {
	_ = varName
	_ = isFromUndefinedCodeRef
	var nameNode pkg_rhino.Node = pkg_rhino.IR_Name(varName)
	_ = nameNode
	if this_.compiler.GetCodingConvention().IsConstant(varName) {
		nameNode.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	}
	var syntheticExternVar pkg_rhino.Node = pkg_rhino.IR_Var_1(nameNode)
	_ = syntheticExternVar
	syntheticExternVar.SetIsSynthesizedUnfulfilledNameDeclaration(isFromUndefinedCodeRef)
	VarCheck_GetSynthesizedExternsRoot(this_.compiler).AddChildToBack(syntheticExternVar)
	this_.compiler.ReportChangeToEnclosingScope(syntheticExternVar)
}

// Java static method: VarCheck.hasDuplicateDeclarationSuppression(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node origVar)
func VarCheck_HasDuplicateDeclarationSuppression(compiler AbstractCompiler, n pkg_rhino.Node, origVar pkg_rhino.Node) bool {
	clinit_VarCheck()
	_ = compiler
	_ = n
	_ = origVar
	if VarCheck_IsExternNamespace(n) {
		return true
	}
	return TypeValidator_HasDuplicateDeclarationSuppression(compiler, origVar)
}

// Java static method: VarCheck.isExternNamespace(com.google.javascript.rhino.Node n)
func VarCheck_IsExternNamespace(n pkg_rhino.Node) bool {
	clinit_VarCheck()
	_ = n
	return ((n.GetParent().IsVar() && n.IsFromExterns()) && NodeUtil_IsNamespaceDecl(n))
}

// Java static method: VarCheck.reportVarMultiplyDeclared(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node current, java.lang.String name, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node original)
func VarCheck_ReportVarMultiplyDeclared(compiler AbstractCompiler, current pkg_rhino.Node, name string, original pkg_rhino.Node) {
	clinit_VarCheck()
	_ = compiler
	_ = current
	_ = name
	_ = original
	compiler.Report_1(JSError_Make_3(current, StaticGet_VarCheck_VAR_MULTIPLY_DECLARED_ERROR(), []string{name, VarCheck_LocationOf(original)}))
}

// Java method: VarCheck.reportBlockScopedMultipleDeclaration(com.google.javascript.rhino.Node current, java.lang.String name, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node original)
func (this_ *VarCheck) priv_VarCheck_reportBlockScopedMultipleDeclaration_5882(current pkg_rhino.Node, name string, original pkg_rhino.Node) {
	_ = current
	_ = name
	_ = original
	this_.compiler.Report_1(JSError_Make_3(current, StaticGet_VarCheck_BLOCK_SCOPED_DECL_MULTIPLY_DECLARED_ERROR(), []string{name, VarCheck_LocationOf(original)}))
}

// Java static method: VarCheck.locationOf(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node n)
func VarCheck_LocationOf(n pkg_rhino.Node) string {
	clinit_VarCheck()
	_ = n
	if n == 0 {
		return "<unknown>"
	} else {
		return n.GetLocation()
	}
}

// Java static method: VarCheck.getSynthesizedExternsRoot(com.google.javascript.jscomp.AbstractCompiler compiler)
func VarCheck_GetSynthesizedExternsRoot(compiler AbstractCompiler) pkg_rhino.Node {
	clinit_VarCheck()
	_ = compiler
	return compiler.GetSynthesizedExternsInput().GetAstRoot(compiler)
}

// Java class: com.google.javascript.jscomp.VarCheck.NameRefInExternsCheck
type VarCheck_NameRefInExternsCheck struct {
	outer_  *VarCheck
	idHash_ atomic.Int32
}

func (this_ *VarCheck_NameRefInExternsCheck) BindSelf_(s any) {
	_ = s
}
func (this_ *VarCheck_NameRefInExternsCheck) IsJavaNil_() bool { return this_ == nil }
func (*VarCheck_NameRefInExternsCheck) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.VarCheck$NameRefInExternsCheck", "NameRefInExternsCheck"
}
func (this_ *VarCheck_NameRefInExternsCheck) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*VarCheck_NameRefInExternsCheck) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.VarCheck.NameRefInExternsCheck", "com.google.javascript.jscomp.VarCheck$NameRefInExternsCheck", func(v_ any) bool {
		return j2gort.InstanceOf[*VarCheck_NameRefInExternsCheck](v_)
	})
}
func (this_ *VarCheck_NameRefInExternsCheck) SetOuter_(o *VarCheck)               { this_.outer_ = o }
func (this_ *VarCheck_NameRefInExternsCheck) GetOuter_() *VarCheck                { return this_.outer_ }
func (*VarCheck_NameRefInExternsCheck) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *VarCheck_NameRefInExternsCheck) initFields_VarCheck_NameRefInExternsCheck() {
	_ = this_
}

// Java constructor: NameRefInExternsCheck()
func NewVarCheck_NameRefInExternsCheck(outer_ *VarCheck) *VarCheck_NameRefInExternsCheck {
	this_ := &VarCheck_NameRefInExternsCheck{}
	this_.outer_ = outer_
	InitVarCheck_NameRefInExternsCheck(this_)
	return this_
}

func InitVarCheck_NameRefInExternsCheck(this_ *VarCheck_NameRefInExternsCheck) {
	this_.initFields_VarCheck_NameRefInExternsCheck()
}

// Java method: NameRefInExternsCheck.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *VarCheck_NameRefInExternsCheck) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	return (!(n.IsScript()) || !(NodeUtil_IsFromTypeSummary(n)))
}

// Java method: NameRefInExternsCheck.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *VarCheck_NameRefInExternsCheck) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsName() {
		switch parent.TokenSwitchOrd_() {
		case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */, int32(73) /* Token.FUNCTION */, int32(105) /* Token.CLASS */, int32(49) /* Token.PARAM_LIST */, int32(134) /* Token.DEFAULT_VALUE */, int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */, int32(102) /* Token.ARRAY_PATTERN */ :
			return
		case int32(100) /* Token.STRING_KEY */ :
			if parent.GetParent().IsObjectPattern() {
				return
			}
		case int32(26) /* Token.GETPROP */ :
			if n == parent.GetFirstChild() {
				var scope *Scope = t.GetScope()
				_ = scope
				var var_ *Var = j2gort.CastPtr[Var](scope.GetVar(n.GetString()))
				_ = var_
				if var_ != nil {
					return
				}
				if parent.MatchesQualifiedName_1_Node(StaticGet_VarCheck_GoogForwardDeclare()) {
					return
				}
				t.Report_3(n, StaticGet_VarCheck_UNDEFINED_EXTERN_VAR_ERROR(), []string{n.GetString()})
				this_.outer_.undefinedNamesFromExterns.Add(n.GetString())
			}
			return
		case int32(51) /* Token.ASSIGN */ :
			if ((n == parent.GetLastChild()) && n.IsQualifiedName()) && parent.GetFirstChild().IsQualifiedName() {
				return
			}
		case int32(32) /* Token.NAME */ :
			if NodeUtil_IsNameDeclaration(parent.GetParent()) {
				return
			}
		case int32(68) /* Token.OR */ :
			if NodeUtil_IsNamespaceDecl(parent.GetParent()) {
				return
			}
		default:
		}
		t.Report_3(n, StaticGet_VarCheck_NAME_REFERENCE_IN_EXTERNS_ERROR(), []string{n.GetString()})
		var scope_2 *Scope = t.GetScope()
		_ = scope_2
		var var__2 *Var = j2gort.CastPtr[Var](scope_2.GetVar(n.GetString()))
		_ = var__2
		if var__2 == nil {
			this_.outer_.undefinedNamesFromExterns.Add(n.GetString())
		}
	}
}

// Java class: com.google.javascript.jscomp.VarCheck.RedeclarationCheckHandler
type VarCheck_RedeclarationCheckHandler struct {
	outer_       *VarCheck
	dupDeclNodes []pkg_rhino.Node `gson:"dupDeclNodes"`
	idHash_      atomic.Int32
}

func (this_ *VarCheck_RedeclarationCheckHandler) BindSelf_(s any) {
	_ = s
}
func (this_ *VarCheck_RedeclarationCheckHandler) IsJavaNil_() bool { return this_ == nil }
func (*VarCheck_RedeclarationCheckHandler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.VarCheck$RedeclarationCheckHandler", "RedeclarationCheckHandler"
}
func (this_ *VarCheck_RedeclarationCheckHandler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*VarCheck_RedeclarationCheckHandler) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.VarCheck.RedeclarationCheckHandler", "com.google.javascript.jscomp.VarCheck$RedeclarationCheckHandler", func(v_ any) bool {
		return j2gort.InstanceOf[*VarCheck_RedeclarationCheckHandler](v_)
	})
}
func (this_ *VarCheck_RedeclarationCheckHandler) SetOuter_(o *VarCheck) { this_.outer_ = o }
func (this_ *VarCheck_RedeclarationCheckHandler) GetOuter_() *VarCheck  { return this_.outer_ }
func (*VarCheck_RedeclarationCheckHandler) Implements_jscomp_SyntacticScopeCreator_RedeclarationHandler() {
}

func (this_ *VarCheck_RedeclarationCheckHandler) initFields_VarCheck_RedeclarationCheckHandler() {
	_ = this_
	this_.dupDeclNodes = []pkg_rhino.Node{}
}

// Java constructor: RedeclarationCheckHandler()
func NewVarCheck_RedeclarationCheckHandler(outer_ *VarCheck) *VarCheck_RedeclarationCheckHandler {
	this_ := &VarCheck_RedeclarationCheckHandler{}
	this_.outer_ = outer_
	InitVarCheck_RedeclarationCheckHandler(this_)
	return this_
}

func InitVarCheck_RedeclarationCheckHandler(this_ *VarCheck_RedeclarationCheckHandler) {
	this_.initFields_VarCheck_RedeclarationCheckHandler()
}

// Java method: RedeclarationCheckHandler.onRedeclaration(com.google.javascript.jscomp.Scope s, java.lang.String name, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.CompilerInput input)
func (this_ *VarCheck_RedeclarationCheckHandler) OnRedeclaration(s *Scope, name string, n pkg_rhino.Node, input *CompilerInput) {
	_ = s
	_ = name
	_ = n
	_ = input
	var parent pkg_rhino.Node = NodeUtil_GetDeclaringParent(n)
	_ = parent
	var origVar *Var = j2gort.CastPtr[Var](s.GetVar(name))
	_ = origVar
	var origNode pkg_rhino.Node = origVar.GetNode()
	_ = origNode
	var origParent pkg_rhino.Node
	_ = origParent
	if origNode == 0 {
		origParent = pkg_rhino.Node(0)
	} else {
		origParent = NodeUtil_GetDeclaringParent(origNode)
	}
	switch parent.TokenSwitchOrd_() {
	case int32(105) /* Token.CLASS */, int32(97) /* Token.CONST */, int32(111) /* Token.LET */ :
		this_.outer_.priv_VarCheck_reportBlockScopedMultipleDeclaration_5882(n, name, origNode)
		return
	default:
	}
	if origParent != 0 {
		switch origParent.TokenSwitchOrd_() {
		case int32(105) /* Token.CLASS */, int32(97) /* Token.CONST */, int32(111) /* Token.LET */ :
			this_.outer_.priv_VarCheck_reportBlockScopedMultipleDeclaration_5882(n, name, origNode)
			return
		case int32(73) /* Token.FUNCTION */ :
			if !(s.IsGlobal()) && parent.IsFunction() {
				this_.outer_.priv_VarCheck_reportBlockScopedMultipleDeclaration_5882(n, name, origNode)
				return
			}
		default:
		}
	}
	if s.IsGlobal() {
		if origParent.IsCatch() && parent.IsCatch() {
			return
		}
		var allowDupe bool = VarCheck_HasDuplicateDeclarationSuppression(this_.outer_.compiler, n, origVar.GetNameNode())
		_ = allowDupe
		if VarCheck_IsExternNamespace(n) {
			this_.dupDeclNodes = append(this_.dupDeclNodes, parent)
			return
		}
		if !(allowDupe) {
			VarCheck_ReportVarMultiplyDeclared(this_.outer_.compiler, n, name, origNode)
		}
	} else {
		var iifeCond_1 bool
		{

			sr_9 := name
			se_10 := "arguments"
			iifeCond_1 = j2gort.CheckNotNullString(sr_9) == se_10
		}
		if iifeCond_1 && !(NodeUtil_IsNameDeclaration(n.GetParent()) && n.IsName()) {
			this_.outer_.compiler.Report_1(JSError_Make_3(n, StaticGet_VarCheck_VAR_ARGUMENTS_SHADOWED_ERROR(), []string{}))
		}
	}
}

// Java method: RedeclarationCheckHandler.removeDuplicates()
func (this_ *VarCheck_RedeclarationCheckHandler) RemoveDuplicates() {
	for _, n := range this_.dupDeclNodes {
		_ = n
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		if parent != 0 {
			n.Detach()
			this_.outer_.compiler.ReportChangeToEnclosingScope(parent)
		}
	}

}
func (*VarCheck) AlwaysShouldTraverse_()                            {}
func (*VarCheck_NameRefInExternsCheck) CanReleaseSyntacticScopes_() {}
func (*VarCheck) CanReleaseSyntacticScopes_()                       {}
