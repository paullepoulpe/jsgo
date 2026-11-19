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
// Port of com/google/javascript/jscomp/SyntacticScopeCreator.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.SyntacticScopeCreator
type SyntacticScopeCreator struct {
	compiler                      AbstractCompiler                           `gson:"compiler"`
	redeclarationHandler          SyntacticScopeCreator_RedeclarationHandler `gson:"redeclarationHandler"`
	reusableScanner               *SyntacticScopeCreator_ScopeScanner        `gson:"reusableScanner"`
	idHash_                       atomic.Int32
	treatProvidesAsRedeclarations bool `gson:"treatProvidesAsRedeclarations"`
	scannerInUse                  bool `gson:"scannerInUse"`
	slab_Var                      j2gort.PooledSlab[Var]
}

func (this_ *SyntacticScopeCreator) BindSelf_(s any) {
	_ = s
}
func (this_ *SyntacticScopeCreator) IsJavaNil_() bool { return this_ == nil }
func (*SyntacticScopeCreator) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SyntacticScopeCreator", "SyntacticScopeCreator"
}
func (this_ *SyntacticScopeCreator) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SyntacticScopeCreator) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SyntacticScopeCreator", "com.google.javascript.jscomp.SyntacticScopeCreator", func(v_ any) bool {
		return j2gort.InstanceOf[*SyntacticScopeCreator](v_)
	})
}
func (*SyntacticScopeCreator) Implements_jscomp_ScopeCreator() {}

func StaticGet_SyntacticScopeCreator_ARGUMENTS() string { return "arguments" }

var static_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER SyntacticScopeCreator_RedeclarationHandler

func StaticGet_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER() SyntacticScopeCreator_RedeclarationHandler {
	clinit_SyntacticScopeCreator()
	return static_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER
}
func StaticSet_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER(v SyntacticScopeCreator_RedeclarationHandler) {
	clinit_SyntacticScopeCreator()
	static_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER = v
}

var clinitDone_SyntacticScopeCreator uint32
var clinitState_SyntacticScopeCreator j2gort.ClinitState

func clinit_SyntacticScopeCreator() {
	if atomic.LoadUint32(&clinitDone_SyntacticScopeCreator) != 0 {
		return
	}
	clinitSlow_SyntacticScopeCreator()
}

func clinitSlow_SyntacticScopeCreator() {
	if !j2gort.EnterClinit(&clinitState_SyntacticScopeCreator, &clinitDone_SyntacticScopeCreator) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_SyntacticScopeCreator, &clinitDone_SyntacticScopeCreator)
	static_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER = NewSyntacticScopeCreator_DefaultRedeclarationHandler()
}

func (this_ *SyntacticScopeCreator) initFields_SyntacticScopeCreator() {
	_ = this_
}

// Java constructor: SyntacticScopeCreator(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewSyntacticScopeCreator_1(compiler AbstractCompiler) *SyntacticScopeCreator {
	clinit_SyntacticScopeCreator()
	this_ := &SyntacticScopeCreator{}
	InitSyntacticScopeCreator_1(this_, compiler)
	return this_
}

func InitSyntacticScopeCreator_1(this_ *SyntacticScopeCreator, compiler AbstractCompiler) {
	_ = compiler
	InitSyntacticScopeCreator_2(this_, compiler, StaticGet_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER())
}

// Java constructor: SyntacticScopeCreator(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.SyntacticScopeCreator.RedeclarationHandler redeclarationHandler)
func NewSyntacticScopeCreator_2(compiler AbstractCompiler, redeclarationHandler SyntacticScopeCreator_RedeclarationHandler) *SyntacticScopeCreator {
	clinit_SyntacticScopeCreator()
	this_ := &SyntacticScopeCreator{}
	InitSyntacticScopeCreator_2(this_, compiler, redeclarationHandler)
	return this_
}

func InitSyntacticScopeCreator_2(this_ *SyntacticScopeCreator, compiler AbstractCompiler, redeclarationHandler SyntacticScopeCreator_RedeclarationHandler) {
	_ = compiler
	_ = redeclarationHandler
	InitSyntacticScopeCreator_3(this_, compiler, redeclarationHandler, false)
}

// Java constructor: SyntacticScopeCreator(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.SyntacticScopeCreator.RedeclarationHandler redeclarationHandler, boolean treatProvidesAsRedeclarations)
func NewSyntacticScopeCreator_3(compiler AbstractCompiler, redeclarationHandler SyntacticScopeCreator_RedeclarationHandler, treatProvidesAsRedeclarations bool) *SyntacticScopeCreator {
	clinit_SyntacticScopeCreator()
	this_ := &SyntacticScopeCreator{}
	InitSyntacticScopeCreator_3(this_, compiler, redeclarationHandler, treatProvidesAsRedeclarations)
	return this_
}

func InitSyntacticScopeCreator_3(this_ *SyntacticScopeCreator, compiler AbstractCompiler, redeclarationHandler SyntacticScopeCreator_RedeclarationHandler, treatProvidesAsRedeclarations bool) {
	_ = compiler
	_ = redeclarationHandler
	_ = treatProvidesAsRedeclarations
	this_.initFields_SyntacticScopeCreator()
	this_.compiler = compiler
	this_.redeclarationHandler = redeclarationHandler
	this_.treatProvidesAsRedeclarations = treatProvidesAsRedeclarations
	this_.reusableScanner = NewSyntacticScopeCreator_ScopeScanner(compiler, redeclarationHandler, treatProvidesAsRedeclarations)
}

// Java method: SyntacticScopeCreator.createScope(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.AbstractScope<?,?> parent)
func (this_ *SyntacticScopeCreator) CreateScope_2_Node_AbstractScope(n pkg_rhino.Node, parent AbstractScope) any {
	_ = n
	_ = parent
	return j2gort.PtrToAny(this_.CreateScope_2_Node_Scope(n, j2gort.CastPtr[Scope](parent)))
}

// Java method: SyntacticScopeCreator.createScope(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.Scope parent)
func (this_ *SyntacticScopeCreator) CreateScope_2_Node_Scope(n pkg_rhino.Node, parent *Scope) *Scope {
	_ = n
	_ = parent
	var scope *Scope
	_ = scope
	if parent == nil {
		scope = Scope_CreateGlobalScope(n)
	} else {
		scope = Scope_CreateChildScope(parent, n)
	}
	scope.varSlab_ = &this_.slab_Var
	if !(this_.scannerInUse) {
		this_.scannerInUse = true
		func() {
			defer func() {
				this_.reusableScanner.Reset()
				this_.scannerInUse = false
			}()
			this_.reusableScanner.Populate_1(scope)
		}()
	} else {
		{
			var val_1 SyntacticScopeCreator_ScopeScanner
			InitSyntacticScopeCreator_ScopeScanner(&val_1, this_.compiler, this_.redeclarationHandler, this_.treatProvidesAsRedeclarations)
			(&val_1).Populate_1(scope)
		}
	}
	return scope
}

// Java class: com.google.javascript.jscomp.SyntacticScopeCreator.ScopeScanner
type SyntacticScopeCreator_ScopeScanner struct {
	scope                         *Scope                                     `gson:"scope"`
	compiler                      AbstractCompiler                           `gson:"compiler"`
	redeclarationHandler          SyntacticScopeCreator_RedeclarationHandler `gson:"redeclarationHandler"`
	inputId                       *pkg_rhino.InputId                         `gson:"inputId"`
	changeRootSet                 j2gort.Ext_java_util_Set                   `gson:"changeRootSet"`
	idHash_                       atomic.Int32
	treatProvidesAsRedeclarations bool `gson:"treatProvidesAsRedeclarations"`
}

func (this_ *SyntacticScopeCreator_ScopeScanner) BindSelf_(s any) {
	_ = s
}
func (this_ *SyntacticScopeCreator_ScopeScanner) IsJavaNil_() bool { return this_ == nil }
func (*SyntacticScopeCreator_ScopeScanner) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SyntacticScopeCreator$ScopeScanner", "ScopeScanner"
}
func (this_ *SyntacticScopeCreator_ScopeScanner) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SyntacticScopeCreator_ScopeScanner) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SyntacticScopeCreator.ScopeScanner", "com.google.javascript.jscomp.SyntacticScopeCreator$ScopeScanner", func(v_ any) bool {
		return j2gort.InstanceOf[*SyntacticScopeCreator_ScopeScanner](v_)
	})
}

func (this_ *SyntacticScopeCreator_ScopeScanner) initFields_SyntacticScopeCreator_ScopeScanner() {
	_ = this_
}

// Java constructor: ScopeScanner(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.SyntacticScopeCreator.RedeclarationHandler redeclarationHandler, boolean treatProvidesAsRedeclarations)
func NewSyntacticScopeCreator_ScopeScanner(compiler AbstractCompiler, redeclarationHandler SyntacticScopeCreator_RedeclarationHandler, treatProvidesAsRedeclarations bool) *SyntacticScopeCreator_ScopeScanner {
	this_ := &SyntacticScopeCreator_ScopeScanner{}
	InitSyntacticScopeCreator_ScopeScanner(this_, compiler, redeclarationHandler, treatProvidesAsRedeclarations)
	return this_
}

func InitSyntacticScopeCreator_ScopeScanner(this_ *SyntacticScopeCreator_ScopeScanner, compiler AbstractCompiler, redeclarationHandler SyntacticScopeCreator_RedeclarationHandler, treatProvidesAsRedeclarations bool) {
	_ = compiler
	_ = redeclarationHandler
	_ = treatProvidesAsRedeclarations
	this_.initFields_SyntacticScopeCreator_ScopeScanner()
	this_.compiler = compiler
	this_.redeclarationHandler = redeclarationHandler
	this_.treatProvidesAsRedeclarations = treatProvidesAsRedeclarations
}

// Java method: ScopeScanner.populate(com.google.javascript.jscomp.Scope scope)
func (this_ *SyntacticScopeCreator_ScopeScanner) Populate_1(scope *Scope) {
	_ = scope
	this_.scope = scope
	this_.changeRootSet = nil
	this_.Populate_0()
}

// Java method: ScopeScanner.reset()
func (this_ *SyntacticScopeCreator_ScopeScanner) Reset() {
	this_.scope = nil
	this_.inputId = nil
	this_.changeRootSet = nil
}

// Java method: ScopeScanner.populate()
func (this_ *SyntacticScopeCreator_ScopeScanner) Populate_0() {
	var n pkg_rhino.Node = this_.scope.GetRootNode_0()
	_ = n
	this_.inputId = NodeUtil_GetInputId(n)
	switch n.TokenSwitchOrd_() {
	case int32(73) /* Token.FUNCTION */ :
		var fnNameNode pkg_rhino.Node = n.GetFirstChild()
		_ = fnNameNode
		var args pkg_rhino.Node = fnNameNode.GetNext()
		_ = args
		if !(args.IsParamList()) {
			j2gort.CheckState(false)
		}
		this_.priv_ScopeScanner_declareLHS_5308(this_.scope, args)
		var fnName string = fnNameNode.GetString()
		_ = fnName
		if !(j2gort.CheckNotNullString(fnName) == "") && NodeUtil_IsFunctionExpression(n) {
			this_.priv_ScopeScanner_declareVar_5311(this_.scope, fnNameNode)
		}
		return
	case int32(105) /* Token.CLASS */ :
		var classNameNode pkg_rhino.Node = n.GetFirstChild()
		_ = classNameNode
		if !(classNameNode.IsEmpty()) && NodeUtil_IsClassExpression(n) {
			this_.priv_ScopeScanner_declareVar_5311(this_.scope, classNameNode)
		}
		return
	case int32(89) /* Token.ROOT */, int32(94) /* Token.SCRIPT */ :
		if !(this_.scope.IsGlobal()) {
			j2gort.CheckStateMsg(false, j2gort.PtrToAny(this_.scope))
		}
		this_.priv_ScopeScanner_scanVars_5310(n, this_.scope, this_.scope)
		return
	case int32(123) /* Token.MODULE_BODY */ :
		this_.priv_ScopeScanner_scanVars_5310(n, this_.scope, this_.scope)
		return
	case int32(80) /* Token.FOR */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(81) /* Token.FOR_IN */, int32(91) /* Token.SWITCH_BODY */ :
		this_.priv_ScopeScanner_scanVars_5310(n, nil, this_.scope)
		return
	case int32(90) /* Token.BLOCK */ :
		if NodeUtil_IsFunctionBlock(n) || NodeUtil_IsClassStaticBlock(n) {
			this_.priv_ScopeScanner_scanVars_5310(n, this_.scope, this_.scope)
		} else {
			this_.priv_ScopeScanner_scanVars_5310(n, nil, this_.scope)
		}
		return
	case int32(109) /* Token.COMPUTED_FIELD_DEF */, int32(108) /* Token.MEMBER_FIELD_DEF */ :
		return
	default:
		panic(j2gort.AsIface[javalang.IRuntimeException](j2gort.ConstructShimThrowable("java.lang.RuntimeException", "Illegal scope root: "+j2gort.StrVal(n))))
	}
}

// Java method: ScopeScanner.declareLHS(com.google.javascript.jscomp.Scope s, com.google.javascript.rhino.Node n)
func (this_ *SyntacticScopeCreator_ScopeScanner) priv_ScopeScanner_declareLHS_5308(s *Scope, n pkg_rhino.Node) {
	_ = s
	_ = n
	var first pkg_rhino.Node = n.GetFirstChild()
	_ = first
	if ((first != 0) && (first.GetNext() == 0)) && first.IsName() {
		this_.priv_ScopeScanner_declareVar_5311(s, first)
		return
	}
	if !(((((((NodeUtil_IsNameDeclaration(n) || n.IsParamList()) || NodeUtil_IsAssignmentOp(n)) || n.IsCatch()) || n.IsDestructuringLhs()) || n.IsDefaultValue()) || n.IsImport()) || NodeUtil_IsEnhancedFor(n)) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(n))
	}
	this_.priv_ScopeScanner_declareLhsHelper_5309(s, n)
}

// Java method: ScopeScanner.declareLhsHelper(com.google.javascript.jscomp.Scope s, com.google.javascript.rhino.Node n)
func (this_ *SyntacticScopeCreator_ScopeScanner) priv_ScopeScanner_declareLhsHelper_5309(s *Scope, n pkg_rhino.Node) {
	_ = s
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(116) /* Token.IMPORT */ :
		this_.priv_ScopeScanner_declareLhsHelper_5309(s, n.GetFirstChild())
		this_.priv_ScopeScanner_declareLhsHelper_5309(s, n.GetSecondChild())
	case int32(84) /* Token.VAR */, int32(97) /* Token.CONST */, int32(111) /* Token.LET */, int32(103) /* Token.OBJECT_PATTERN */, int32(102) /* Token.ARRAY_PATTERN */, int32(49) /* Token.PARAM_LIST */, int32(117) /* Token.IMPORT_SPECS */ :
		{
			var child pkg_rhino.Node = n.GetFirstChild()
			_ = child
			for ; child != 0; child = child.GetNext() {
				if child.IsName() {
					this_.priv_ScopeScanner_declareVar_5311(s, child)
				} else {
					this_.priv_ScopeScanner_declareLhsHelper_5309(s, child)
				}
			}
		}
	case int32(104) /* Token.DESTRUCTURING_LHS */, int32(134) /* Token.DEFAULT_VALUE */, int32(86) /* Token.CATCH */, int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */, int32(101) /* Token.CAST */ :
		this_.priv_ScopeScanner_declareLhsHelper_5309(s, n.GetFirstChild())
	case int32(118) /* Token.IMPORT_SPEC */, int32(129) /* Token.COMPUTED_PROP */, int32(100) /* Token.STRING_KEY */ :
		this_.priv_ScopeScanner_declareLhsHelper_5309(s, n.GetLastChild())
	case int32(32) /* Token.NAME */, int32(119) /* Token.IMPORT_STAR */, int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */ :
		this_.priv_ScopeScanner_declareVar_5311(s, n)
	case int32(88) /* Token.EMPTY */ :
	case int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
		this_.priv_ScopeScanner_declareLhsHelper_5309(s, n.GetFirstChild())
	default:
		if NodeUtil_IsAssignmentOp(n) {
			this_.priv_ScopeScanner_declareLhsHelper_5309(s, n.GetFirstChild())
		} else {
			panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "Invalid node in lhs: "+j2gort.StrVal(n))))
		}
	}
}

// Java method: ScopeScanner.scanVars(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Scope hoistScope, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Scope blockScope)
func (this_ *SyntacticScopeCreator_ScopeScanner) priv_ScopeScanner_scanVars_5310(n pkg_rhino.Node, hoistScope *Scope, blockScope *Scope) {
	_ = n
	_ = hoistScope
	_ = blockScope
	switch n.TokenSwitchOrd_() {
	case int32(84) /* Token.VAR */ :
		if hoistScope != nil {
			this_.priv_ScopeScanner_declareLHS_5308(hoistScope, n)
		}
		return
	case int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
		if blockScope != nil {
			this_.priv_ScopeScanner_declareLHS_5308(blockScope, n)
		}
		return
	case int32(116) /* Token.IMPORT */ :
		this_.priv_ScopeScanner_declareLHS_5308(hoistScope, n)
		return
	case int32(120) /* Token.EXPORT */ :
		this_.priv_ScopeScanner_scanVars_5310(n.GetFirstChild(), hoistScope, blockScope)
		return
	case int32(73) /* Token.FUNCTION */ :
		if NodeUtil_IsFunctionExpression(n) || (blockScope == nil) {
			return
		}
		var fnName string = n.GetFirstChild().GetString()
		_ = fnName
		if j2gort.CheckNotNullString(fnName) == "" {
			return
		}
		this_.priv_ScopeScanner_declareVar_5311(blockScope, n.GetFirstChild())
		return
	case int32(105) /* Token.CLASS */ :
		if NodeUtil_IsClassExpression(n) || (blockScope == nil) {
			return
		}
		var className string = n.GetFirstChild().GetString()
		_ = className
		if j2gort.CheckNotNullString(className) == "" {
			return
		}
		this_.priv_ScopeScanner_declareVar_5311(blockScope, n.GetFirstChild())
		return
	case int32(86) /* Token.CATCH */ :
		if !(n.HasTwoChildren()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(n))
		}
		if blockScope != nil {
			this_.priv_ScopeScanner_declareLHS_5308(blockScope, n)
		}
		var block pkg_rhino.Node = n.GetSecondChild()
		_ = block
		this_.priv_ScopeScanner_scanVars_5310(block, hoistScope, blockScope)
		return
	case int32(94) /* Token.SCRIPT */ :
		if (this_.changeRootSet != nil) && !(j2gort.SetContainsNode(this_.changeRootSet, n)) {
			return
		}
		this_.inputId = n.GetInputId()
	case int32(123) /* Token.MODULE_BODY */ :
		if hoistScope.IsGlobal() {
			var expr pkg_rhino.Node = n.GetFirstChild()
			_ = expr
			if (expr != 0) && SyntacticScopeCreator_ScopeScanner_IsLegacyGoogModule(expr) {
				this_.priv_ScopeScanner_declareImplicitGoogNamespaceFromCall_5312(hoistScope, expr)
			}
			return
		}
	case int32(93) /* Token.EXPR_RESULT */ :
		if !(n.GetParent().IsScript()) {
			break
		}
		if NodeUtil_IsGoogProvideCall(n) {
			this_.priv_ScopeScanner_declareImplicitGoogNamespaceFromCall_5312(j2gort.CastPtr[Scope](hoistScope.GetGlobalScope()), n)
		} else if NodeUtil_IsBundledGoogModuleCall(n.GetFirstChild()) && n.GetFirstChild().GetSecondChild().IsFunction() {
			var fn pkg_rhino.Node = n.GetFirstChild().GetSecondChild()
			_ = fn
			var moduleCall pkg_rhino.Node = NodeUtil_GetFunctionBody(fn).GetFirstChild()
			_ = moduleCall
			if SyntacticScopeCreator_ScopeScanner_IsLegacyGoogModule(moduleCall) {
				this_.priv_ScopeScanner_declareImplicitGoogNamespaceFromCall_5312(j2gort.CastPtr[Scope](hoistScope.GetGlobalScope()), moduleCall)
			}
		}
	default:
	}
	var isBlockStart bool = ((blockScope != nil) && (n == blockScope.GetRootNode_0()))
	_ = isBlockStart
	var enteringNewBlock bool = (!(isBlockStart) && NodeUtil_CreatesBlockScope(n))
	_ = enteringNewBlock
	if enteringNewBlock && (hoistScope == nil) {
		return
	}
	if NodeUtil_IsShallowStatementTree(n) {
		{
			var child pkg_rhino.Node = n.GetFirstChild()
			_ = child
			for child != 0 {
				var next pkg_rhino.Node = child.GetNext()
				_ = next
				this_.priv_ScopeScanner_scanVars_5310(child, hoistScope, func() *Scope {
					if enteringNewBlock {
						return nil
					}
					return blockScope
				}())
				child = next
			}
		}
	}
}

// Java method: ScopeScanner.declareVar(com.google.javascript.jscomp.Scope s, com.google.javascript.rhino.Node n)
func (this_ *SyntacticScopeCreator_ScopeScanner) priv_ScopeScanner_declareVar_5311(s *Scope, n pkg_rhino.Node) {
	_ = s
	_ = n
	if !(n.IsName() || n.IsImportStar()) {
		j2gort.CheckStateFmt(false, "Invalid node for declareVar: %s", j2gort.BoxNode(n))
	}
	var name string = n.GetString()
	_ = name
	var v *Var = j2gort.CastPtr[Var](s.GetOwnSlot(name))
	_ = v
	if v != nil {
		if v.GetNode() == n {
			return
		} else if v.IsImplicitGoogNamespace() {
			s.Undeclare(ToIface_AbstractVar(v))
			v = nil
		}
	}
	var input *CompilerInput = this_.compiler.GetInput(this_.inputId)
	_ = input
	if ((v != nil) || !(SyntacticScopeCreator_ScopeScanner_IsShadowingAllowed(name, s))) || ((s.IsFunctionScope() || s.IsFunctionBlockScope()) && func() bool { sr_2 := name; se_3 := "arguments"; return j2gort.CheckNotNullString(sr_2) == se_3 }()) {
		this_.redeclarationHandler.OnRedeclaration(s, name, n, input)
	} else {
		s.Declare(name, n, input)
	}
}

// Java method: ScopeScanner.declareImplicitGoogNamespaceFromCall(com.google.javascript.jscomp.Scope s, com.google.javascript.rhino.Node exprCall)
func (this_ *SyntacticScopeCreator_ScopeScanner) priv_ScopeScanner_declareImplicitGoogNamespaceFromCall_5312(s *Scope, exprCall pkg_rhino.Node) {
	_ = s
	_ = exprCall
	var namespaceNode pkg_rhino.Node = exprCall.GetFirstChild().GetSecondChild()
	_ = namespaceNode
	if (namespaceNode == 0) || !(namespaceNode.IsString()) {
		return
	}
	var namespace string = namespaceNode.GetString()
	_ = namespace
	var root string = NodeUtil_GetRootOfQualifiedName_1_string(namespace)
	_ = root
	if j2gort.CheckNotNullString(root) == "" {
		return
	}
	if this_.treatProvidesAsRedeclarations {
		var existing *Var = j2gort.CastPtr[Var](s.GetOwnSlot(root))
		_ = existing
		if (existing != nil) && !(existing.IsImplicitGoogNamespace()) {
			this_.redeclarationHandler.OnRedeclaration(s, root, existing.GetNode(), this_.compiler.GetInput(this_.inputId))
		}
	}
	s.DeclareImplicitGoogNamespaceIfAbsent(root, namespaceNode)
}

// Java static method: ScopeScanner.isShadowingAllowed(java.lang.String name, com.google.javascript.jscomp.Scope s)
func SyntacticScopeCreator_ScopeScanner_IsShadowingAllowed(name string, s *Scope) bool {
	_ = name
	_ = s
	if s.IsFunctionBlockScope() {
		var maybeParam *Var = j2gort.CastPtr[Var](j2gort.CastPtr[Scope](s.GetParent()).GetOwnSlot(name))
		_ = maybeParam
		return ((maybeParam == nil) || !(maybeParam.IsParam()))
	}
	return true
}

// Java static method: ScopeScanner.isLegacyGoogModule(com.google.javascript.rhino.Node n)
func SyntacticScopeCreator_ScopeScanner_IsLegacyGoogModule(n pkg_rhino.Node) bool {
	_ = n
	if !(NodeUtil_IsGoogModuleCall(n)) {
		return false
	}
	return ((n.GetNext() != 0) && NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(n.GetNext()))
}

// Java interface: com.google.javascript.jscomp.SyntacticScopeCreator.RedeclarationHandler
type SyntacticScopeCreator_RedeclarationHandler interface {
	Implements_jscomp_SyntacticScopeCreator_RedeclarationHandler()
	IsJavaNil_() bool
	OnRedeclaration(*Scope, string, pkg_rhino.Node, *CompilerInput)
}

func ToIface_SyntacticScopeCreator_RedeclarationHandler[T interface {
	SyntacticScopeCreator_RedeclarationHandler
	comparable
}](v T) SyntacticScopeCreator_RedeclarationHandler {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_SyntacticScopeCreator_RedeclarationHandler[T SyntacticScopeCreator_RedeclarationHandler](v T) SyntacticScopeCreator_RedeclarationHandler {
	var vi SyntacticScopeCreator_RedeclarationHandler = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SyntacticScopeCreator.RedeclarationHandler", "com.google.javascript.jscomp.SyntacticScopeCreator$RedeclarationHandler", func(v_ any) bool {
		return j2gort.InstanceOf[SyntacticScopeCreator_RedeclarationHandler](v_)
	})
}

// Java class: com.google.javascript.jscomp.SyntacticScopeCreator.DefaultRedeclarationHandler
type SyntacticScopeCreator_DefaultRedeclarationHandler struct {
	idHash_ atomic.Int32
}

func (this_ *SyntacticScopeCreator_DefaultRedeclarationHandler) BindSelf_(s any) {
	_ = s
}
func (this_ *SyntacticScopeCreator_DefaultRedeclarationHandler) IsJavaNil_() bool {
	return this_ == nil
}
func (*SyntacticScopeCreator_DefaultRedeclarationHandler) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.SyntacticScopeCreator$DefaultRedeclarationHandler", "DefaultRedeclarationHandler"
}
func (this_ *SyntacticScopeCreator_DefaultRedeclarationHandler) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*SyntacticScopeCreator_DefaultRedeclarationHandler) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.SyntacticScopeCreator.DefaultRedeclarationHandler", "com.google.javascript.jscomp.SyntacticScopeCreator$DefaultRedeclarationHandler", func(v_ any) bool {
		return j2gort.InstanceOf[*SyntacticScopeCreator_DefaultRedeclarationHandler](v_)
	})
}
func (*SyntacticScopeCreator_DefaultRedeclarationHandler) Implements_jscomp_SyntacticScopeCreator_RedeclarationHandler() {
}

func (this_ *SyntacticScopeCreator_DefaultRedeclarationHandler) initFields_SyntacticScopeCreator_DefaultRedeclarationHandler() {
	_ = this_
}

// Java constructor: DefaultRedeclarationHandler()
func NewSyntacticScopeCreator_DefaultRedeclarationHandler() *SyntacticScopeCreator_DefaultRedeclarationHandler {
	this_ := &SyntacticScopeCreator_DefaultRedeclarationHandler{}
	InitSyntacticScopeCreator_DefaultRedeclarationHandler(this_)
	return this_
}

func InitSyntacticScopeCreator_DefaultRedeclarationHandler(this_ *SyntacticScopeCreator_DefaultRedeclarationHandler) {
	this_.initFields_SyntacticScopeCreator_DefaultRedeclarationHandler()
}

// Java method: DefaultRedeclarationHandler.onRedeclaration(com.google.javascript.jscomp.Scope s, java.lang.String name, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.CompilerInput input)
func (this_ *SyntacticScopeCreator_DefaultRedeclarationHandler) OnRedeclaration(s *Scope, name string, n pkg_rhino.Node, input *CompilerInput) {
	_ = s
	_ = name
	_ = n
	_ = input
}
func (this_ *SyntacticScopeCreator) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.slab_Var.Release()
}
