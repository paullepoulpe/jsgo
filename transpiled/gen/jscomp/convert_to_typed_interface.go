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
// Port of com/google/javascript/jscomp/ijs/ConvertToTypedInterface.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_convert_to_typed_interface_1 j2gort.Ext_java_util_function_Function
var lambdaSingleton_convert_to_typed_interface_2 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_convert_to_typed_interface_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.ValueOfInt(ConvertToTypedInterface_SimplifyDeclarations_CountDots(j2gort.CastString(mrp0_)))
	})
	lambdaSingleton_convert_to_typed_interface_2 = j2gort.NewLambda_Ext_java_util_function_Function(func(lp_3 any) any {
		decl := j2gort.AsIface[PotentialDeclaration](lp_3)
		_ = decl
		var jsdoc *pkg_rhino.JSDocInfo = decl.GetJsDoc()
		_ = jsdoc
		if jsdoc == nil {
			return j2gort.BoxNullableVal(StaticGet_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_NO_JSDOC())
		}
		if !(jsdoc.HasTypeNodes()) {
			return j2gort.BoxNullableVal(StaticGet_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_UNTYPED_JSDOC_DECLARATION())
		}
		return j2gort.BoxNullableVal(StaticGet_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_TYPED_JSDOC_DECLARATION())
	})
}

// Java class: com.google.javascript.jscomp.ijs.ConvertToTypedInterface
type ConvertToTypedInterface struct {
	compiler AbstractCompiler `gson:"compiler"`
	idHash_  atomic.Int32
}

func (this_ *ConvertToTypedInterface) BindSelf_(s any) {
	_ = s
}
func (this_ *ConvertToTypedInterface) IsJavaNil_() bool { return this_ == nil }
func (*ConvertToTypedInterface) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ijs.ConvertToTypedInterface", "ConvertToTypedInterface"
}
func (this_ *ConvertToTypedInterface) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ConvertToTypedInterface) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.ConvertToTypedInterface", "com.google.javascript.jscomp.ijs.ConvertToTypedInterface", func(v_ any) bool {
		return j2gort.InstanceOf[*ConvertToTypedInterface](v_)
	})
}
func (*ConvertToTypedInterface) Implements_jscomp_CompilerPass() {}

var static_ConvertToTypedInterface_CONSTANT_WITH_SUGGESTED_TYPE *DiagnosticType

func StaticGet_ConvertToTypedInterface_CONSTANT_WITH_SUGGESTED_TYPE() *DiagnosticType {
	clinit_ConvertToTypedInterface()
	return static_ConvertToTypedInterface_CONSTANT_WITH_SUGGESTED_TYPE
}
func StaticSet_ConvertToTypedInterface_CONSTANT_WITH_SUGGESTED_TYPE(v *DiagnosticType) {
	clinit_ConvertToTypedInterface()
	static_ConvertToTypedInterface_CONSTANT_WITH_SUGGESTED_TYPE = v
}

var static_ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE *DiagnosticType

func StaticGet_ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE() *DiagnosticType {
	clinit_ConvertToTypedInterface()
	return static_ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE
}
func StaticSet_ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE(v *DiagnosticType) {
	clinit_ConvertToTypedInterface()
	static_ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE = v
}

var static_ConvertToTypedInterface_GOOG_SCOPE_HIDDEN_TYPE *DiagnosticType

func StaticGet_ConvertToTypedInterface_GOOG_SCOPE_HIDDEN_TYPE() *DiagnosticType {
	clinit_ConvertToTypedInterface()
	return static_ConvertToTypedInterface_GOOG_SCOPE_HIDDEN_TYPE
}
func StaticSet_ConvertToTypedInterface_GOOG_SCOPE_HIDDEN_TYPE(v *DiagnosticType) {
	clinit_ConvertToTypedInterface()
	static_ConvertToTypedInterface_GOOG_SCOPE_HIDDEN_TYPE = v
}

var static_ConvertToTypedInterface_CALLS_TO_PRESERVE j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_ConvertToTypedInterface_CALLS_TO_PRESERVE() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_ConvertToTypedInterface()
	return static_ConvertToTypedInterface_CALLS_TO_PRESERVE
}
func StaticSet_ConvertToTypedInterface_CALLS_TO_PRESERVE(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_ConvertToTypedInterface()
	static_ConvertToTypedInterface_CALLS_TO_PRESERVE = v
}

var clinitDone_ConvertToTypedInterface uint32
var clinitState_ConvertToTypedInterface j2gort.ClinitState

func clinit_ConvertToTypedInterface() {
	if atomic.LoadUint32(&clinitDone_ConvertToTypedInterface) != 0 {
		return
	}
	clinitSlow_ConvertToTypedInterface()
}

func clinitSlow_ConvertToTypedInterface() {
	if !j2gort.EnterClinit(&clinitState_ConvertToTypedInterface, &clinitDone_ConvertToTypedInterface) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ConvertToTypedInterface, &clinitDone_ConvertToTypedInterface)
	static_ConvertToTypedInterface_CONSTANT_WITH_SUGGESTED_TYPE = DiagnosticType_Warning("JSC_CONSTANT_WITH_SUGGESTED_TYPE", "Constants in top-level should have types explicitly specified.\nYou may want specify this type as:\t@const '{'{0}'}'")
	static_ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE = DiagnosticType_Warning("JSC_CONSTANT_WITHOUT_EXPLICIT_TYPE", "Constants in top-level should have types explicitly specified.")
	static_ConvertToTypedInterface_GOOG_SCOPE_HIDDEN_TYPE = DiagnosticType_Warning("JSC_GOOG_SCOPE_HIDDEN_TYPE", "Found a goog.scope local type declaration.\nIf you can't yet migrate this file to goog.module, use a /** @private */ declaration within the goog.provide namespace.")
	static_ConvertToTypedInterface_CALLS_TO_PRESERVE = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_7("Polymer", "goog.addSingletonGetter", "goog.define", "goog.forwardDeclare", "goog.module", "goog.module.declareLegacyNamespace", []any{"goog.declareModuleId", "goog.provide", "goog.require", "goog.requireType"})
}

func (this_ *ConvertToTypedInterface) initFields_ConvertToTypedInterface() {
	_ = this_
}

// Java constructor: ConvertToTypedInterface(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewConvertToTypedInterface(compiler AbstractCompiler) *ConvertToTypedInterface {
	clinit_ConvertToTypedInterface()
	this_ := &ConvertToTypedInterface{}
	InitConvertToTypedInterface(this_, compiler)
	return this_
}

func InitConvertToTypedInterface(this_ *ConvertToTypedInterface, compiler AbstractCompiler) {
	_ = compiler
	this_.initFields_ConvertToTypedInterface()
	this_.compiler = compiler
}

// Java static method: ConvertToTypedInterface.maybeReport(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node node, com.google.javascript.jscomp.DiagnosticType diagnostic, java.lang.String[] fillers)
func ConvertToTypedInterface_MaybeReport(compiler AbstractCompiler, node pkg_rhino.Node, diagnostic *DiagnosticType, fillers []string) {
	clinit_ConvertToTypedInterface()
	_ = compiler
	_ = node
	_ = diagnostic
	_ = fillers
	var sourceName string = NodeUtil_GetSourceName(node)
	_ = sourceName
	if (javalang.StringEndsWith(j2gort.CheckNotNullString(sourceName), "_test.js") || javalang.StringEndsWith(j2gort.CheckNotNullString(sourceName), "_test.closure.js")) || javalang.StringEndsWith(j2gort.CheckNotNullString(sourceName), "_test.tsx.cl.js") {
		return
	}
	compiler.Report_1(JSError_Make_3(node, diagnostic, fillers))
}

// Java static method: ConvertToTypedInterface.maybeWarnForConstWithoutExplicitType(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.ijs.PotentialDeclaration decl)
func ConvertToTypedInterface_MaybeWarnForConstWithoutExplicitType(compiler AbstractCompiler, decl PotentialDeclaration) {
	clinit_ConvertToTypedInterface()
	_ = compiler
	_ = decl
	if (decl.IsConstToBeInferred() && !(decl.GetLhs().IsFromExterns())) && !(JsdocUtil_IsPrivate(decl.GetJsDoc())) {
		var nameNode pkg_rhino.Node = decl.GetLhs()
		_ = nameNode
		if j2gort.IfaceEquals[pkg_rhino.JSType](nameNode.GetJSType(), nil) {
			ConvertToTypedInterface_MaybeReport(compiler, nameNode, StaticGet_ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE(), []string{})
		} else {
			ConvertToTypedInterface_MaybeReport(compiler, nameNode, StaticGet_ConvertToTypedInterface_CONSTANT_WITH_SUGGESTED_TYPE(), []string{nameNode.GetJSType().ToAnnotationString(pkg_rhino.StaticGet_JSType_Nullability_EXPLICIT())})
		}
	}
}

// Java method: ConvertToTypedInterface.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ConvertToTypedInterface) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	{
		var script pkg_rhino.Node = root.GetFirstChild()
		_ = script
		var next pkg_rhino.Node = pkg_rhino.Node(0)
		_ = next
		for ; script != 0; script = next {
			next = script.GetNext()
			this_.priv_ConvertToTypedInterface_processFile_6380(script)
		}
	}
}

// Java method: ConvertToTypedInterface.processFile(com.google.javascript.rhino.Node scriptNode)
func (this_ *ConvertToTypedInterface) priv_ConvertToTypedInterface_processFile_6380(scriptNode pkg_rhino.Node) {
	_ = scriptNode
	if !(scriptNode.IsScript()) {
		j2gort.CheckArgument(false)
	}
	var sourceFileName string = scriptNode.GetSourceFileName()
	_ = sourceFileName
	if AbstractCompiler_IsFillFileName(sourceFileName) {
		scriptNode.Detach()
		this_.compiler.ReportChangeToChangeScope(scriptNode)
		scriptNode.SetDeleted(true)
		return
	}
	var scriptJsDoc *pkg_rhino.JSDocInfo = scriptNode.GetJSDocInfo()
	_ = scriptJsDoc
	if (scriptJsDoc != nil) && scriptJsDoc.IsClosureUnawareCode() {
		var scriptJsDocBuilder *pkg_rhino.JSDocInfo_Builder = scriptJsDoc.ToBuilder_0()
		_ = scriptJsDocBuilder
		scriptJsDocBuilder.RemoveClosureUnawareCode()
		scriptNode.SetJSDocInfo(scriptJsDocBuilder.Build_0())
	}
	var val_1 FileInfo
	InitFileInfo(&val_1, sourceFileName)
	var currentFile *FileInfo = &val_1
	_ = currentFile
	NodeTraversal_Traverse(this_.compiler, scriptNode, NewConvertToTypedInterface_RemoveNonDeclarations())
	NodeTraversal_Traverse(this_.compiler, scriptNode, NewConvertToTypedInterface_PropagateConstJsdoc(currentFile))
	{
		clinit_ConvertToTypedInterface_SimplifyDeclarations()
		var val_2 ConvertToTypedInterface_SimplifyDeclarations
		InitConvertToTypedInterface_SimplifyDeclarations(&val_2, this_.compiler, currentFile)
		(&val_2).SimplifyAll()
	}
}

// Java static method: ConvertToTypedInterface.findNameDeclaration(com.google.javascript.jscomp.Scope scope, com.google.javascript.rhino.Node rhs)
func ConvertToTypedInterface_FindNameDeclaration(scope *Scope, rhs pkg_rhino.Node) *Var {
	clinit_ConvertToTypedInterface()
	_ = scope
	_ = rhs
	if !(rhs.IsName()) {
		return nil
	}
	return j2gort.CastPtr[Var](scope.GetVar(rhs.GetString()))
}

// Java static method: ConvertToTypedInterface.isSymbolProp(com.google.javascript.rhino.Node lhs)
func ConvertToTypedInterface_IsSymbolProp(lhs pkg_rhino.Node) bool {
	clinit_ConvertToTypedInterface()
	_ = lhs
	return (lhs.IsGetProp() && lhs.GetFirstChild().MatchesName_1_string("Symbol"))
}

// Java static method: ConvertToTypedInterface.shouldPreserveAssignment(com.google.javascript.rhino.Node expr, com.google.javascript.jscomp.NodeTraversal t)
func ConvertToTypedInterface_ShouldPreserveAssignment(expr pkg_rhino.Node, t *NodeTraversal) bool {
	clinit_ConvertToTypedInterface()
	_ = expr
	_ = t
	var lhs pkg_rhino.Node = expr.GetFirstChild()
	_ = lhs
	if !(t.InGlobalHoistScope()) && !(t.InModuleHoistScope()) {
		return ClassUtil_IsThisPropInsideClassWithName(lhs)
	}
	if lhs.IsGetElem() && ConvertToTypedInterface_IsSymbolProp(lhs.GetSecondChild()) {
		return lhs.GetFirstChild().IsQualifiedName()
	}
	if lhs.IsName() {
		return (t.InGlobalScope() || t.InModuleScope())
	}
	return lhs.IsQualifiedName()
}

// Java class: com.google.javascript.jscomp.ijs.ConvertToTypedInterface.RemoveNonDeclarations
type ConvertToTypedInterface_RemoveNonDeclarations struct {
	idHash_ atomic.Int32
}

func (this_ *ConvertToTypedInterface_RemoveNonDeclarations) BindSelf_(s any) {
	_ = s
}
func (this_ *ConvertToTypedInterface_RemoveNonDeclarations) IsJavaNil_() bool { return this_ == nil }
func (*ConvertToTypedInterface_RemoveNonDeclarations) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ijs.ConvertToTypedInterface$RemoveNonDeclarations", "RemoveNonDeclarations"
}
func (this_ *ConvertToTypedInterface_RemoveNonDeclarations) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ConvertToTypedInterface_RemoveNonDeclarations) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.ConvertToTypedInterface.RemoveNonDeclarations", "com.google.javascript.jscomp.ijs.ConvertToTypedInterface$RemoveNonDeclarations", func(v_ any) bool {
		return j2gort.InstanceOf[*ConvertToTypedInterface_RemoveNonDeclarations](v_)
	})
}
func (*ConvertToTypedInterface_RemoveNonDeclarations) Implements_jscomp_NodeTraversal_Callback() {}

var static_ConvertToTypedInterface_RemoveNonDeclarations_GOOG_SCOPE pkg_rhino.QualifiedName

func StaticGet_ConvertToTypedInterface_RemoveNonDeclarations_GOOG_SCOPE() pkg_rhino.QualifiedName {
	clinit_ConvertToTypedInterface_RemoveNonDeclarations()
	return static_ConvertToTypedInterface_RemoveNonDeclarations_GOOG_SCOPE
}
func StaticSet_ConvertToTypedInterface_RemoveNonDeclarations_GOOG_SCOPE(v pkg_rhino.QualifiedName) {
	clinit_ConvertToTypedInterface_RemoveNonDeclarations()
	static_ConvertToTypedInterface_RemoveNonDeclarations_GOOG_SCOPE = v
}

var clinitDone_ConvertToTypedInterface_RemoveNonDeclarations uint32
var clinitState_ConvertToTypedInterface_RemoveNonDeclarations j2gort.ClinitState

func clinit_ConvertToTypedInterface_RemoveNonDeclarations() {
	if atomic.LoadUint32(&clinitDone_ConvertToTypedInterface_RemoveNonDeclarations) != 0 {
		return
	}
	clinitSlow_ConvertToTypedInterface_RemoveNonDeclarations()
}

func clinitSlow_ConvertToTypedInterface_RemoveNonDeclarations() {
	if !j2gort.EnterClinit(&clinitState_ConvertToTypedInterface_RemoveNonDeclarations, &clinitDone_ConvertToTypedInterface_RemoveNonDeclarations) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ConvertToTypedInterface_RemoveNonDeclarations, &clinitDone_ConvertToTypedInterface_RemoveNonDeclarations)
	static_ConvertToTypedInterface_RemoveNonDeclarations_GOOG_SCOPE = pkg_rhino.QualifiedName_Of("goog.scope")
}

func (this_ *ConvertToTypedInterface_RemoveNonDeclarations) initFields_ConvertToTypedInterface_RemoveNonDeclarations() {
	_ = this_
}

// Java constructor: RemoveNonDeclarations()
func NewConvertToTypedInterface_RemoveNonDeclarations() *ConvertToTypedInterface_RemoveNonDeclarations {
	this_ := &ConvertToTypedInterface_RemoveNonDeclarations{}
	InitConvertToTypedInterface_RemoveNonDeclarations(this_)
	return this_
}

func InitConvertToTypedInterface_RemoveNonDeclarations(this_ *ConvertToTypedInterface_RemoveNonDeclarations) {
	this_.initFields_ConvertToTypedInterface_RemoveNonDeclarations()
}

// Java method: RemoveNonDeclarations.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ConvertToTypedInterface_RemoveNonDeclarations) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	{
		var expr pkg_rhino.Node = pkg_rhino.Node(0)
		_ = expr
		var initializer pkg_rhino.Node = pkg_rhino.Node(0)
		_ = initializer
		switch n.TokenSwitchOrd_() {
		case int32(73) /* Token.FUNCTION */ :
			if !(ClassUtil_IsConstructor(n)) || !(ClassUtil_HasNamedClass(n)) {
				var body pkg_rhino.Node = n.GetLastChild()
				_ = body
				if !(body.IsBlock()) || body.HasChildren() {
					t.ReportCodeChange_1(body)
					body.ReplaceWith(pkg_rhino.IR_Block_0().Srcref(body))
					NodeUtil_MarkFunctionsDeleted_2_Node_AbstractCompiler(body, t.GetCompiler())
				}
			}
			return true
		case int32(108) /* Token.MEMBER_FIELD_DEF */ :
			if ClassUtil_IsMemberFieldDefInsideClassWithName(n) {
				return true
			}
			NodeUtil_DeleteNode(n, t.GetCompiler())
			return false
		case int32(93) /* Token.EXPR_RESULT */ :
			expr = n.GetFirstChild()
			{
				var callee pkg_rhino.Node = pkg_rhino.Node(0)
				_ = callee
				switch expr.TokenSwitchOrd_() {
				case int32(28) /* Token.CALL */ :
					callee = expr.GetFirstChild()
					if !(!(StaticGet_ConvertToTypedInterface_RemoveNonDeclarations_GOOG_SCOPE().Matches(callee))) {
						j2gort.CheckState(false)
					}
					if j2gort.SetContainsStr(StaticGet_ConvertToTypedInterface_CALLS_TO_PRESERVE(), callee.GetQualifiedName()) {
						return true
					}
					NodeUtil_DeleteNode(n, t.GetCompiler())
					return false
				case int32(51) /* Token.ASSIGN */ :
					if ConvertToTypedInterface_ShouldPreserveAssignment(expr, t) {
						return true
					}
					NodeUtil_DeleteNode(n, t.GetCompiler())
					return false
				case int32(26) /* Token.GETPROP */ :
					if !(expr.IsQualifiedName()) || (expr.GetJSDocInfo() == nil) {
						NodeUtil_DeleteNode(n, t.GetCompiler())
						return false
					}
					return true
				case int32(27) /* Token.GETELEM */ :
					if ConvertToTypedInterface_IsSymbolProp(expr.GetSecondChild()) && (expr.GetJSDocInfo() != nil) {
						return true
					}
					NodeUtil_DeleteNode(n, t.GetCompiler())
					return false
				default:
					NodeUtil_DeleteNode(n, t.GetCompiler())
					return false
				}
			}
			fallthrough
		case int32(129) /* Token.COMPUTED_PROP */ :
			if ClassUtil_IsComputedMemberInsideClassWithName(n) {
				return true
			}
			if !(NodeUtil_IsLhsByDestructuring(n.GetSecondChild())) {
				NodeUtil_DeleteNode(n, t.GetCompiler())
			}
			return false
		case int32(109) /* Token.COMPUTED_FIELD_DEF */ :
			if ClassUtil_IsComputedMemberInsideClassWithName(n) {
				return true
			}
			NodeUtil_DeleteNode(n, t.GetCompiler())
			return false
		case int32(43) /* Token.THROW */, int32(0) /* Token.RETURN */, int32(82) /* Token.BREAK */, int32(83) /* Token.CONTINUE */, int32(98) /* Token.DEBUGGER */, int32(88) /* Token.EMPTY */ :
			if NodeUtil_IsStatementParent(parent) {
				NodeUtil_DeleteNode(n, t.GetCompiler())
			}
			return false
		case int32(92) /* Token.LABEL */, int32(74) /* Token.IF */, int32(75) /* Token.SWITCH */, int32(76) /* Token.CASE */, int32(78) /* Token.WHILE */ :
			NodeUtil_DeleteNode(n.GetFirstChild(), t.GetCompiler())
			return true
		case int32(48) /* Token.TRY */, int32(79) /* Token.DO */ :
			NodeUtil_DeleteNode(n.GetSecondChild(), t.GetCompiler())
			return true
		case int32(80) /* Token.FOR */ :
			NodeUtil_DeleteNode(n.GetSecondChild(), t.GetCompiler())
			fallthrough
		case int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(81) /* Token.FOR_IN */ :
			NodeUtil_DeleteNode(n.GetSecondChild(), t.GetCompiler())
			initializer = n.RemoveFirstChild()
			if initializer.IsVar() {
				n.GetLastChild().AddChildToFront(initializer)
			}
			return true
		case int32(97) /* Token.CONST */, int32(111) /* Token.LET */ :
			if !(t.InGlobalScope()) && !(t.InModuleScope()) {
				NodeUtil_RemoveChild(parent, n)
				t.ReportCodeChange_1(parent)
				return false
			}
			return true
		case int32(84) /* Token.VAR */ :
			if !(t.InGlobalHoistScope()) && !(t.InModuleHoistScope()) {
				NodeUtil_RemoveChild(parent, n)
				t.ReportCodeChange_1(parent)
				return false
			}
			return true
		case int32(123) /* Token.MODULE_BODY */, int32(105) /* Token.CLASS */, int32(77) /* Token.DEFAULT_CASE */, int32(90) /* Token.BLOCK */, int32(120) /* Token.EXPORT */, int32(116) /* Token.IMPORT */ :
			return true
		default:
			if !(!(NodeUtil_IsStatement(n))) {
				j2gort.CheckStateMsg(false, j2gort.BoxNullableVal(n.GetToken()))
			}
			return true
		}
	}
	panic("unreachable")
}

// Java method: RemoveNonDeclarations.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ConvertToTypedInterface_RemoveNonDeclarations) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(48) /* Token.TRY */, int32(92) /* Token.LABEL */, int32(77) /* Token.DEFAULT_CASE */, int32(76) /* Token.CASE */, int32(79) /* Token.DO */, int32(78) /* Token.WHILE */, int32(80) /* Token.FOR */, int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */, int32(74) /* Token.IF */ :
		if n.HasParent() {
			var children pkg_rhino.Node = n.RemoveChildren()
			_ = children
			parent.AddChildrenAfter(children, n)
			n.Detach()
			t.ReportCodeChange_0()
		}
	case int32(75) /* Token.SWITCH */ :
		if n.HasParent() {
			if !(n.HasOneChild()) {
				j2gort.CheckStateFmt(false, "malfrmed sWITCH %s", j2gort.BoxString(n.ToStringTree()))
			}
			var children_2 pkg_rhino.Node = n.GetFirstChild().RemoveChildren()
			_ = children_2
			parent.AddChildrenAfter(children_2, n)
			n.Detach()
			t.ReportCodeChange_0()
		}
	case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
		ConvertToTypedInterface_RemoveNonDeclarations_SplitNameDeclarationsAndRemoveDestructuring(n, t)
	case int32(90) /* Token.BLOCK */ :
		if !(parent.IsFunction()) {
			parent.AddChildrenAfter(n.RemoveChildren(), n)
			n.Detach()
			t.ReportCodeChange_1(parent)
		}
	default:
	}
}

// Java static method: RemoveNonDeclarations.splitNameDeclarationsAndRemoveDestructuring(com.google.javascript.rhino.Node n, com.google.javascript.jscomp.NodeTraversal t)
func ConvertToTypedInterface_RemoveNonDeclarations_SplitNameDeclarationsAndRemoveDestructuring(n pkg_rhino.Node, t *NodeTraversal) {
	_ = n
	_ = t
	if !(NodeUtil_IsNameDeclaration(n)) {
		j2gort.CheckArgument(false)
	}
	var sharedJsdoc *pkg_rhino.JSDocInfo = n.GetJSDocInfo()
	_ = sharedJsdoc
	var isExport bool = n.GetParent().IsExport()
	_ = isExport
	var statement pkg_rhino.Node
	_ = statement
	if isExport {
		statement = n.GetParent()
	} else {
		statement = n
	}
	for n.HasChildren() {
		var lhsToSplit pkg_rhino.Node = n.GetLastChild()
		_ = lhsToSplit
		var nameJsdoc *pkg_rhino.JSDocInfo = lhsToSplit.GetJSDocInfo()
		_ = nameJsdoc
		lhsToSplit.SetJSDocInfo(nil)
		var mergedJsdoc *pkg_rhino.JSDocInfo = JsdocUtil_MergeJsdocs(sharedJsdoc, nameJsdoc)
		_ = mergedJsdoc
		if n.HasOneChild() {
			n.SetJSDocInfo(mergedJsdoc)
			return
		}
		var rhs pkg_rhino.Node
		_ = rhs
		if lhsToSplit.HasChildren() {
			rhs = lhsToSplit.RemoveFirstChild()
		} else {
			rhs = pkg_rhino.Node(0)
		}
		var newDeclaration pkg_rhino.Node = NodeUtil_NewDeclaration(lhsToSplit.Detach(), rhs, n.GetToken()).Srcref(n)
		_ = newDeclaration
		newDeclaration.SetJSDocInfo(mergedJsdoc)
		if isExport {
			newDeclaration = pkg_rhino.IR_Export(newDeclaration).Srcref(statement)
		}
		newDeclaration.InsertAfter(statement)
		t.ReportCodeChange_0()
	}
}

// Java class: com.google.javascript.jscomp.ijs.ConvertToTypedInterface.PropagateConstJsdoc
type ConvertToTypedInterface_PropagateConstJsdoc struct {
	ProcessConstJsdocCallbackBase
}

func (this_ *ConvertToTypedInterface_PropagateConstJsdoc) BindSelf_(s any) {
	_ = s
	this_.ProcessConstJsdocCallbackBase.BindSelf_(s)
}
func (this_ *ConvertToTypedInterface_PropagateConstJsdoc) IsJavaNil_() bool { return this_ == nil }
func (*ConvertToTypedInterface_PropagateConstJsdoc) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ijs.ConvertToTypedInterface$PropagateConstJsdoc", "PropagateConstJsdoc"
}
func (this_ *ConvertToTypedInterface_PropagateConstJsdoc) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ProcessConstJsdocCallbackBase.IdentityHash_()
}
func (*ConvertToTypedInterface_PropagateConstJsdoc) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.ConvertToTypedInterface.PropagateConstJsdoc", "com.google.javascript.jscomp.ijs.ConvertToTypedInterface$PropagateConstJsdoc", func(v_ any) bool {
		return j2gort.InstanceOf[*ConvertToTypedInterface_PropagateConstJsdoc](v_)
	})
}
func (*ConvertToTypedInterface_PropagateConstJsdoc) Implements_jscomp_ProcessConstJsdocCallback() {}
func (*ConvertToTypedInterface_PropagateConstJsdoc) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*ConvertToTypedInterface_PropagateConstJsdoc) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *ConvertToTypedInterface_PropagateConstJsdoc) initFields_ConvertToTypedInterface_PropagateConstJsdoc() {
	_ = this_
}

// Java constructor: PropagateConstJsdoc(com.google.javascript.jscomp.ijs.FileInfo currentFile)
func NewConvertToTypedInterface_PropagateConstJsdoc(currentFile *FileInfo) *ConvertToTypedInterface_PropagateConstJsdoc {
	this_ := &ConvertToTypedInterface_PropagateConstJsdoc{}
	this_.BindSelf_(this_)
	InitConvertToTypedInterface_PropagateConstJsdoc(this_, currentFile)
	return this_
}

func InitConvertToTypedInterface_PropagateConstJsdoc(this_ *ConvertToTypedInterface_PropagateConstJsdoc, currentFile *FileInfo) {
	_ = currentFile
	InitProcessConstJsdocCallback(&this_.ProcessConstJsdocCallbackBase, currentFile)
	this_.initFields_ConvertToTypedInterface_PropagateConstJsdoc()
}

// Java method: PropagateConstJsdoc.processConstWithRhs(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node nameNode)
func (this_ *ConvertToTypedInterface_PropagateConstJsdoc) ProcessConstWithRhs(t *NodeTraversal, nameNode pkg_rhino.Node) {
	_ = t
	_ = nameNode
	if !(((nameNode.IsQualifiedName() || nameNode.IsStringKey()) || nameNode.IsDestructuringLhs()) || nameNode.IsMemberFieldDef()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(nameNode))
	}
	var jsdocNode pkg_rhino.Node = NodeUtil_GetBestJSDocInfoNode(nameNode)
	_ = jsdocNode
	var originalJsdoc *pkg_rhino.JSDocInfo = jsdocNode.GetJSDocInfo()
	_ = originalJsdoc
	var rhs pkg_rhino.Node = NodeUtil_GetRValueOfLValue(nameNode)
	_ = rhs
	var newJsdoc *pkg_rhino.JSDocInfo = JsdocUtil_GetJSDocForRhs(rhs, originalJsdoc)
	_ = newJsdoc
	if (newJsdoc == nil) && ClassUtil_IsThisPropInsideClassWithName(nameNode) {
		var decl *Var = ConvertToTypedInterface_FindNameDeclaration(t.GetScope(), rhs)
		_ = decl
		newJsdoc = JsdocUtil_GetJSDocForName(decl, originalJsdoc)
	}
	if newJsdoc != nil {
		jsdocNode.SetJSDocInfo(newJsdoc)
		t.ReportCodeChange_0()
	}
}

// Java class: com.google.javascript.jscomp.ijs.ConvertToTypedInterface.SimplifyDeclarations
type ConvertToTypedInterface_SimplifyDeclarations struct {
	compiler    AbstractCompiler `gson:"compiler"`
	currentFile *FileInfo        `gson:"currentFile"`
	idHash_     atomic.Int32
}

func (this_ *ConvertToTypedInterface_SimplifyDeclarations) BindSelf_(s any) {
	_ = s
}
func (this_ *ConvertToTypedInterface_SimplifyDeclarations) IsJavaNil_() bool { return this_ == nil }
func (*ConvertToTypedInterface_SimplifyDeclarations) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ijs.ConvertToTypedInterface$SimplifyDeclarations", "SimplifyDeclarations"
}
func (this_ *ConvertToTypedInterface_SimplifyDeclarations) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ConvertToTypedInterface_SimplifyDeclarations) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.ConvertToTypedInterface.SimplifyDeclarations", "com.google.javascript.jscomp.ijs.ConvertToTypedInterface$SimplifyDeclarations", func(v_ any) bool {
		return j2gort.InstanceOf[*ConvertToTypedInterface_SimplifyDeclarations](v_)
	})
}

var static_ConvertToTypedInterface_SimplifyDeclarations_SHORT_TO_LONG j2gort.Ext_java_util_Comparator

func StaticGet_ConvertToTypedInterface_SimplifyDeclarations_SHORT_TO_LONG() j2gort.Ext_java_util_Comparator {
	clinit_ConvertToTypedInterface_SimplifyDeclarations()
	return static_ConvertToTypedInterface_SimplifyDeclarations_SHORT_TO_LONG
}
func StaticSet_ConvertToTypedInterface_SimplifyDeclarations_SHORT_TO_LONG(v j2gort.Ext_java_util_Comparator) {
	clinit_ConvertToTypedInterface_SimplifyDeclarations()
	static_ConvertToTypedInterface_SimplifyDeclarations_SHORT_TO_LONG = v
}

var static_ConvertToTypedInterface_SimplifyDeclarations_DECLARATIONS_FIRST j2gort.Ext_java_util_Comparator

func StaticGet_ConvertToTypedInterface_SimplifyDeclarations_DECLARATIONS_FIRST() j2gort.Ext_java_util_Comparator {
	clinit_ConvertToTypedInterface_SimplifyDeclarations()
	return static_ConvertToTypedInterface_SimplifyDeclarations_DECLARATIONS_FIRST
}
func StaticSet_ConvertToTypedInterface_SimplifyDeclarations_DECLARATIONS_FIRST(v j2gort.Ext_java_util_Comparator) {
	clinit_ConvertToTypedInterface_SimplifyDeclarations()
	static_ConvertToTypedInterface_SimplifyDeclarations_DECLARATIONS_FIRST = v
}

var clinitDone_ConvertToTypedInterface_SimplifyDeclarations uint32
var clinitState_ConvertToTypedInterface_SimplifyDeclarations j2gort.ClinitState

func clinit_ConvertToTypedInterface_SimplifyDeclarations() {
	if atomic.LoadUint32(&clinitDone_ConvertToTypedInterface_SimplifyDeclarations) != 0 {
		return
	}
	clinitSlow_ConvertToTypedInterface_SimplifyDeclarations()
}

func clinitSlow_ConvertToTypedInterface_SimplifyDeclarations() {
	if !j2gort.EnterClinit(&clinitState_ConvertToTypedInterface_SimplifyDeclarations, &clinitDone_ConvertToTypedInterface_SimplifyDeclarations) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ConvertToTypedInterface_SimplifyDeclarations, &clinitDone_ConvertToTypedInterface_SimplifyDeclarations)
	static_ConvertToTypedInterface_SimplifyDeclarations_SHORT_TO_LONG = j2gort.Ext_java_util_Comparator_Comparing_1(lambdaSingleton_convert_to_typed_interface_1)
	static_ConvertToTypedInterface_SimplifyDeclarations_DECLARATIONS_FIRST = j2gort.Ext_java_util_Comparator_Comparing_1(lambdaSingleton_convert_to_typed_interface_2)
}

func (this_ *ConvertToTypedInterface_SimplifyDeclarations) initFields_ConvertToTypedInterface_SimplifyDeclarations() {
	_ = this_
}

// Java constructor: SimplifyDeclarations(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.ijs.FileInfo currentFile)
func NewConvertToTypedInterface_SimplifyDeclarations(compiler AbstractCompiler, currentFile *FileInfo) *ConvertToTypedInterface_SimplifyDeclarations {
	clinit_ConvertToTypedInterface_SimplifyDeclarations()
	this_ := &ConvertToTypedInterface_SimplifyDeclarations{}
	InitConvertToTypedInterface_SimplifyDeclarations(this_, compiler, currentFile)
	return this_
}

func InitConvertToTypedInterface_SimplifyDeclarations(this_ *ConvertToTypedInterface_SimplifyDeclarations, compiler AbstractCompiler, currentFile *FileInfo) {
	_ = compiler
	_ = currentFile
	this_.initFields_ConvertToTypedInterface_SimplifyDeclarations()
	this_.compiler = compiler
	this_.currentFile = currentFile
}

// Java static method: SimplifyDeclarations.countDots(java.lang.String name)
func ConvertToTypedInterface_SimplifyDeclarations_CountDots(name string) int32 {
	clinit_ConvertToTypedInterface_SimplifyDeclarations()
	_ = name
	var count int32 = int32(0)
	_ = count
	{
		var i int32 = int32(0)
		_ = i
		for ; i < j2gort.JStrLen(name); i = i + 1 {
			if int32(j2gort.JStrCharAt(name, i)) == int32(uint16(46)) {
				count = count + 1
			}
		}
	}
	return count
}

// Java method: SimplifyDeclarations.removeDuplicateDeclarations()
func (this_ *ConvertToTypedInterface_SimplifyDeclarations) priv_SimplifyDeclarations_removeDuplicateDeclarations_6382() {
	{
		cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_Set](this_.currentFile.GetDeclarations().KeySet())))
		for cur_4.HasNext() {
			name := cur_4.NextStr()
			_ = name
			if javalang.StringStartsWith(j2gort.CheckNotNullString(name), "this.") {
				continue
			}
			var declList j2gort.Ext_java_util_List = j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.MultimapGetStr(this_.currentFile.GetDeclarations(), name))
			_ = declList
			declList.Sort(StaticGet_ConvertToTypedInterface_SimplifyDeclarations_DECLARATIONS_FIRST())
			for declList.Size() > int32(1) {
				var decl PotentialDeclaration = j2gort.AsIface[PotentialDeclaration](declList.Remove_1_int32(int32(1)))
				_ = decl
				decl.Remove(this_.compiler)
			}
		}
	}
}

// Java method: SimplifyDeclarations.simplifyAll()
func (this_ *ConvertToTypedInterface_SimplifyDeclarations) SimplifyAll() {
	this_.priv_SimplifyDeclarations_removeDuplicateDeclarations_6382()
	var seenNames []string = j2gort.ImmListToSliceNullable[string](j2gort.Ext_com_google_common_collect_ImmutableList_SortedCopyOf_2(StaticGet_ConvertToTypedInterface_SimplifyDeclarations_SHORT_TO_LONG(), j2gort.AsIface[j2gort.Ext_java_util_Set](this_.currentFile.GetDeclarations().KeySet())))
	_ = seenNames
	{

		for _, name := range seenNames {

			_ = name
			{
				cur_6 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.AsIface[j2gort.Ext_java_util_List](j2gort.MultimapGetStr(this_.currentFile.GetDeclarations(), name))))
				for cur_6.HasNext() {
					rawElem_7 := cur_6.Next()
					decl := j2gort.AsIface[PotentialDeclaration](rawElem_7)
					_ = decl
					this_.priv_SimplifyDeclarations_processDeclaration_6384(name, decl)
				}
			}
		}
	}
}

// Java method: SimplifyDeclarations.processDeclaration(java.lang.String name, com.google.javascript.jscomp.ijs.PotentialDeclaration decl)
func (this_ *ConvertToTypedInterface_SimplifyDeclarations) priv_SimplifyDeclarations_processDeclaration_6384(name string, decl PotentialDeclaration) {
	_ = name
	_ = decl
	if this_.priv_SimplifyDeclarations_shouldRemove_6388(name, decl) {
		decl.Remove(this_.compiler)
		return
	}
	if decl.BreakDownDestructure(this_.compiler) {
		ConvertToTypedInterface_MaybeReport(this_.compiler, decl.GetLhs(), StaticGet_ConvertToTypedInterface_CONSTANT_WITHOUT_EXPLICIT_TYPE(), []string{})
	}
	if (decl.GetRhs() != 0) && decl.GetRhs().IsFunction() {
		this_.priv_SimplifyDeclarations_processFunction_6386(decl.GetRhs())
	} else if (decl.GetRhs() != 0) && ConvertToTypedInterface_SimplifyDeclarations_IsClass(decl.GetRhs()) {
		this_.priv_SimplifyDeclarations_processClass_6385(decl.GetRhs())
	}
	this_.priv_SimplifyDeclarations_setUndeclaredToUnusableType_6389(decl)
	this_.priv_SimplifyDeclarations_elidePrivateMemberType_6390(decl)
	decl.Simplify(this_.compiler)
}

// Java method: SimplifyDeclarations.processClass(com.google.javascript.rhino.Node n)
func (this_ *ConvertToTypedInterface_SimplifyDeclarations) priv_SimplifyDeclarations_processClass_6385(n pkg_rhino.Node) {
	_ = n
	if !(ConvertToTypedInterface_SimplifyDeclarations_IsClass(n)) {
		j2gort.CheckArgument(false)
	}
	{
		var member pkg_rhino.Node = n.GetLastChild().GetFirstChild()
		_ = member
		for member != 0 {
			var next pkg_rhino.Node = member.GetNext()
			_ = next
			switch member.TokenSwitchOrd_() {
			case int32(108) /* Token.MEMBER_FIELD_DEF */, int32(109) /* Token.COMPUTED_FIELD_DEF */ :
			case int32(88) /* Token.EMPTY */ :
				NodeUtil_DeleteNode(member, this_.compiler)
			case int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */ :
				this_.priv_SimplifyDeclarations_processFunction_6386(member.GetLastChild())
			case int32(129) /* Token.COMPUTED_PROP */ :
				if !(member.GetSecondChild().IsFunction()) {
					j2gort.CheckStateFmt(false, "Non-function computed class member: %s", j2gort.BoxNode(member))
				}
				this_.priv_SimplifyDeclarations_processFunction_6386(member.GetSecondChild())
			default:
				panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", j2gort.StrVal(member.GetToken())+" should not be handled by processClass")))
			}
			member = next
		}
	}
}

// Java method: SimplifyDeclarations.processFunction(com.google.javascript.rhino.Node n)
func (this_ *ConvertToTypedInterface_SimplifyDeclarations) priv_SimplifyDeclarations_processFunction_6386(n pkg_rhino.Node) {
	_ = n
	if !(n.IsFunction()) {
		j2gort.CheckArgument(false)
	}
	this_.priv_SimplifyDeclarations_processFunctionParameters_6387(n.GetSecondChild())
}

// Java method: SimplifyDeclarations.processFunctionParameters(com.google.javascript.rhino.Node paramList)
func (this_ *ConvertToTypedInterface_SimplifyDeclarations) priv_SimplifyDeclarations_processFunctionParameters_6387(paramList pkg_rhino.Node) {
	_ = paramList
	if !(paramList.IsParamList()) {
		j2gort.CheckArgument(false)
	}
	{
		var arg pkg_rhino.Node = paramList.GetFirstChild()
		_ = arg
		for ; arg != 0; arg = arg.GetNext() {
			if arg.IsDefaultValue() {
				var rhs pkg_rhino.Node = arg.GetLastChild()
				_ = rhs
				if rhs.IsCall() && StaticGet_RewriteCallerCodeLocation_GOOG_CALLER_LOCATION_QUALIFIED_NAME().Matches(rhs.GetFirstChild()) {
					continue
				}
				rhs.ReplaceWith(NodeUtil_NewUndefinedNode(rhs))
				this_.compiler.ReportChangeToEnclosingScope(arg)
			}
		}
	}
}

// Java static method: SimplifyDeclarations.isClass(com.google.javascript.rhino.Node n)
func ConvertToTypedInterface_SimplifyDeclarations_IsClass(n pkg_rhino.Node) bool {
	clinit_ConvertToTypedInterface_SimplifyDeclarations()
	_ = n
	return n.IsClass()
}

// Java static method: SimplifyDeclarations.rootName(java.lang.String qualifiedName)
func ConvertToTypedInterface_SimplifyDeclarations_RootName(qualifiedName string) string {
	clinit_ConvertToTypedInterface_SimplifyDeclarations()
	_ = qualifiedName
	var dotIndex int32 = j2gort.JStrIndexOfCp(qualifiedName, int32(uint16(46)))
	_ = dotIndex
	if dotIndex == int32(-1) {
		return qualifiedName
	}
	return j2gort.JStrSub(qualifiedName, int32(0), dotIndex)
}

// Java method: SimplifyDeclarations.shouldRemove(java.lang.String name, com.google.javascript.jscomp.ijs.PotentialDeclaration decl)
func (this_ *ConvertToTypedInterface_SimplifyDeclarations) priv_SimplifyDeclarations_shouldRemove_6388(name string, decl PotentialDeclaration) bool {
	_ = name
	_ = decl
	if decl.IsDetached() {
		return true
	}
	if javalang.StringStartsWith(j2gort.CheckNotNullString(ConvertToTypedInterface_SimplifyDeclarations_RootName(name)), "$jscomp") {
		if ((decl.GetRhs() != 0) && decl.GetRhs().IsClass()) || ((decl.GetJsDoc() != nil) && decl.GetJsDoc().ContainsTypeDefinition()) {
			ConvertToTypedInterface_MaybeReport(this_.compiler, decl.GetLhs(), StaticGet_ConvertToTypedInterface_GOOG_SCOPE_HIDDEN_TYPE(), []string{})
		}
		return true
	}
	return ((((!(javalang.StringStartsWith(j2gort.CheckNotNullString(name), "this.")) && !(decl.IsDefiniteDeclaration(this_.compiler))) && !(decl.GetLhs().IsMemberFieldDef())) && !(this_.currentFile.IsPrefixProvided(name))) && !(this_.currentFile.IsStrictPrefixDeclared(name)))
}

// Java method: SimplifyDeclarations.setUndeclaredToUnusableType(com.google.javascript.jscomp.ijs.PotentialDeclaration decl)
func (this_ *ConvertToTypedInterface_SimplifyDeclarations) priv_SimplifyDeclarations_setUndeclaredToUnusableType_6389(decl PotentialDeclaration) {
	_ = decl
	var nameNode pkg_rhino.Node = decl.GetLhs()
	_ = nameNode
	var jsdoc *pkg_rhino.JSDocInfo = decl.GetJsDoc()
	_ = jsdoc
	if ((decl.ShouldPreserve() || NodeUtil_IsNamespaceDecl(nameNode)) || ((decl.GetRhs() != 0) && NodeUtil_IsCallTo_2_Node_string(decl.GetRhs(), "Symbol"))) || (((jsdoc != nil) && jsdoc.ContainsDeclaration()) && !(decl.IsConstToBeInferred())) {
		return
	}
	ConvertToTypedInterface_MaybeWarnForConstWithoutExplicitType(this_.compiler, decl)
	var jsdocNode pkg_rhino.Node = NodeUtil_GetBestJSDocInfoNode(nameNode)
	_ = jsdocNode
	jsdocNode.SetJSDocInfo(JsdocUtil_GetUnusableTypeJSDoc(jsdoc))
}

// Java method: SimplifyDeclarations.elidePrivateMemberType(com.google.javascript.jscomp.ijs.PotentialDeclaration decl)
func (this_ *ConvertToTypedInterface_SimplifyDeclarations) priv_SimplifyDeclarations_elidePrivateMemberType_6390(decl PotentialDeclaration) {
	_ = decl
	if this_.priv_SimplifyDeclarations_shouldElidePrivateMemberType_6391(decl) {
		var jsdoc *pkg_rhino.JSDocInfo = decl.GetJsDoc()
		_ = jsdoc
		var jsdocNode pkg_rhino.Node = NodeUtil_GetBestJSDocInfoNode(decl.GetLhs())
		_ = jsdocNode
		jsdocNode.SetJSDocInfo(JsdocUtil_ReplaceAnnotatedTypesWithUnusableType(jsdoc))
	}
}

// Java method: SimplifyDeclarations.shouldElidePrivateMemberType(com.google.javascript.jscomp.ijs.PotentialDeclaration decl)
func (this_ *ConvertToTypedInterface_SimplifyDeclarations) priv_SimplifyDeclarations_shouldElidePrivateMemberType_6391(decl PotentialDeclaration) bool {
	_ = decl
	if !(this_.currentFile.IsFromTypeScript()) {
		return false
	}
	var jsdoc *pkg_rhino.JSDocInfo = decl.GetJsDoc()
	_ = jsdoc
	if (jsdoc == nil) || (jsdoc.GetVisibility() != pkg_rhino.StaticGet_JSDocInfo_Visibility_PRIVATE()) {
		return false
	}
	if (jsdoc.IsConstructorOrInterface() || jsdoc.HasEnumParameterType()) || jsdoc.HasTypedefType() {
		return false
	}
	if decl.GetLhs().IsName() {
		return false
	}
	return JsdocUtil_HasAnnotatedType(jsdoc)
}

// Java enum: com.google.javascript.jscomp.ijs.ConvertToTypedInterface.SimplifyDeclarations.TypingLevel
type ConvertToTypedInterface_SimplifyDeclarations_TypingLevel uint8

const (
	enum_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_TYPED_JSDOC_DECLARATION   ConvertToTypedInterface_SimplifyDeclarations_TypingLevel = 1
	enum_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_UNTYPED_JSDOC_DECLARATION ConvertToTypedInterface_SimplifyDeclarations_TypingLevel = 2
	enum_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_NO_JSDOC                  ConvertToTypedInterface_SimplifyDeclarations_TypingLevel = 3
)

var enumBaseTable_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel = [...]j2gort.EnumBase{
	{Name_: "TYPED_JSDOC_DECLARATION", Ordinal_: 0},
	{Name_: "UNTYPED_JSDOC_DECLARATION", Ordinal_: 1},
	{Name_: "NO_JSDOC", Ordinal_: 2},
}
var enumValues_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel = []ConvertToTypedInterface_SimplifyDeclarations_TypingLevel{1, 2, 3}

func StaticGet_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_TYPED_JSDOC_DECLARATION() ConvertToTypedInterface_SimplifyDeclarations_TypingLevel {
	return enum_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_TYPED_JSDOC_DECLARATION
}
func StaticGet_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_UNTYPED_JSDOC_DECLARATION() ConvertToTypedInterface_SimplifyDeclarations_TypingLevel {
	return enum_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_UNTYPED_JSDOC_DECLARATION
}
func StaticGet_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_NO_JSDOC() ConvertToTypedInterface_SimplifyDeclarations_TypingLevel {
	return enum_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_NO_JSDOC
}

func ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_Values() []ConvertToTypedInterface_SimplifyDeclarations_TypingLevel {
	return j2gort.CloneSlice(enumValues_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel)
}

func ConvertToTypedInterface_SimplifyDeclarations_TypingLevel_ValueOf(name string) ConvertToTypedInterface_SimplifyDeclarations_TypingLevel {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel {
		if enumBaseTable_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel[i_].Name_ == name {
			return ConvertToTypedInterface_SimplifyDeclarations_TypingLevel(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.ijs.ConvertToTypedInterface.SimplifyDeclarations.TypingLevel." + name))
}

func (this_ ConvertToTypedInterface_SimplifyDeclarations_TypingLevel) IsJavaNil_() bool {
	return this_ == 0
}
func (ConvertToTypedInterface_SimplifyDeclarations_TypingLevel) HasIdentityEquals_() {}
func (ConvertToTypedInterface_SimplifyDeclarations_TypingLevel) Implements_j2gort_Ext_java_lang_Enum() {
}
func (ConvertToTypedInterface_SimplifyDeclarations_TypingLevel) Implements_j2gort_Ext_java_lang_constant_Constable() {
}
func (this_ ConvertToTypedInterface_SimplifyDeclarations_TypingLevel) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (ConvertToTypedInterface_SimplifyDeclarations_TypingLevel) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ijs.ConvertToTypedInterface$SimplifyDeclarations$TypingLevel", "TypingLevel"
}
func (this_ ConvertToTypedInterface_SimplifyDeclarations_TypingLevel) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel[this_-1]
}
func (this_ ConvertToTypedInterface_SimplifyDeclarations_TypingLevel) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ ConvertToTypedInterface_SimplifyDeclarations_TypingLevel) EnumOrdinal_() int32 {
	return int32(this_) - 1
}
func (this_ ConvertToTypedInterface_SimplifyDeclarations_TypingLevel) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel[this_-1].Name_
}
func (this_ ConvertToTypedInterface_SimplifyDeclarations_TypingLevel) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel[this_-1].Name_
}
func (this_ ConvertToTypedInterface_SimplifyDeclarations_TypingLevel) Equals(other_ any) bool {
	o_, ok_ := other_.(ConvertToTypedInterface_SimplifyDeclarations_TypingLevel)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ ConvertToTypedInterface_SimplifyDeclarations_TypingLevel) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel[this_-1].HashCode()
}
func (this_ ConvertToTypedInterface_SimplifyDeclarations_TypingLevel) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel[this_-1].CompareTo(other_)
}
func (this_ ConvertToTypedInterface_SimplifyDeclarations_TypingLevel) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_ConvertToTypedInterface_SimplifyDeclarations_TypingLevel[this_-1].Name_)
}
func init() {
	vals_ := []any{ConvertToTypedInterface_SimplifyDeclarations_TypingLevel(1), ConvertToTypedInterface_SimplifyDeclarations_TypingLevel(2), ConvertToTypedInterface_SimplifyDeclarations_TypingLevel(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ijs.ConvertToTypedInterface.SimplifyDeclarations.TypingLevel", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.ijs.ConvertToTypedInterface$SimplifyDeclarations$TypingLevel", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ijs.ConvertToTypedInterface.SimplifyDeclarations.TypingLevel", "com.google.javascript.jscomp.ijs.ConvertToTypedInterface$SimplifyDeclarations$TypingLevel", func(v_ any) bool {
		return j2gort.InstanceOf[ConvertToTypedInterface_SimplifyDeclarations_TypingLevel](v_)
	})
}
func (*ConvertToTypedInterface_RemoveNonDeclarations) CanReleaseSyntacticScopes_() {}
func (*ConvertToTypedInterface_PropagateConstJsdoc) AlwaysShouldTraverse_()        {}
