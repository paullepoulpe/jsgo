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
// Port of com/google/javascript/jscomp/GenerateExports.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.GenerateExports
type GenerateExports struct {
	compiler               AbstractCompiler         `gson:"compiler"`
	exportSymbolFunction   string                   `gson:"exportSymbolFunction"`
	exportPropertyFunction string                   `gson:"exportPropertyFunction"`
	exportedVariables      j2gort.Ext_java_util_Set `gson:"exportedVariables"`
	idHash_                atomic.Int32
	allowNonGlobalExports  bool `gson:"allowNonGlobalExports"`
}

func (this_ *GenerateExports) BindSelf_(s any) {
	_ = s
}
func (this_ *GenerateExports) IsJavaNil_() bool { return this_ == nil }
func (*GenerateExports) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GenerateExports", "GenerateExports"
}
func (this_ *GenerateExports) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GenerateExports) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GenerateExports", "com.google.javascript.jscomp.GenerateExports", func(v_ any) bool {
		return j2gort.InstanceOf[*GenerateExports](v_)
	})
}
func (*GenerateExports) Implements_jscomp_CompilerPass() {}

func StaticGet_GenerateExports_PROTOTYPE_PROPERTY() string { return "prototype" }

var static_GenerateExports_MISSING_EXPORT_CONVENTION *DiagnosticType

func StaticGet_GenerateExports_MISSING_EXPORT_CONVENTION() *DiagnosticType {
	clinit_GenerateExports()
	return static_GenerateExports_MISSING_EXPORT_CONVENTION
}
func StaticSet_GenerateExports_MISSING_EXPORT_CONVENTION(v *DiagnosticType) {
	clinit_GenerateExports()
	static_GenerateExports_MISSING_EXPORT_CONVENTION = v
}

var static_GenerateExports_MISSING_GOOG_FOR_EXPORT *DiagnosticType

func StaticGet_GenerateExports_MISSING_GOOG_FOR_EXPORT() *DiagnosticType {
	clinit_GenerateExports()
	return static_GenerateExports_MISSING_GOOG_FOR_EXPORT
}
func StaticSet_GenerateExports_MISSING_GOOG_FOR_EXPORT(v *DiagnosticType) {
	clinit_GenerateExports()
	static_GenerateExports_MISSING_GOOG_FOR_EXPORT = v
}

var clinitDone_GenerateExports uint32
var clinitState_GenerateExports j2gort.ClinitState

func clinit_GenerateExports() {
	if atomic.LoadUint32(&clinitDone_GenerateExports) != 0 {
		return
	}
	clinitSlow_GenerateExports()
}

func clinitSlow_GenerateExports() {
	if !j2gort.EnterClinit(&clinitState_GenerateExports, &clinitDone_GenerateExports) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_GenerateExports, &clinitDone_GenerateExports)
	static_GenerateExports_MISSING_EXPORT_CONVENTION = DiagnosticType_Error("JSC_MISSING_EXPORT_CONVENTION", "@export cannot be used without defining a exportProperty and exportSymbol function in the coding convention")
	static_GenerateExports_MISSING_GOOG_FOR_EXPORT = DiagnosticType_Error("JSC_MISSING_EXPORT_SYMBOL_DEFINITION", "@export cannot be used without including closure/base.js or other definition of {0} and {1}")
}

func (this_ *GenerateExports) initFields_GenerateExports() {
	_ = this_
	this_.exportSymbolFunction = j2gort.NullString
	this_.exportPropertyFunction = j2gort.NullString
	this_.exportedVariables = j2gort.NewExt_java_util_LinkedHashSet()
}

// Java constructor: GenerateExports(com.google.javascript.jscomp.AbstractCompiler compiler, boolean allowNonGlobalExports, java.lang.@org.jspecify.annotations.Nullable String exportSymbolFunction, java.lang.@org.jspecify.annotations.Nullable String exportPropertyFunction)
func NewGenerateExports(compiler AbstractCompiler, allowNonGlobalExports bool, exportSymbolFunction string, exportPropertyFunction string) *GenerateExports {
	clinit_GenerateExports()
	this_ := &GenerateExports{}
	InitGenerateExports(this_, compiler, allowNonGlobalExports, exportSymbolFunction, exportPropertyFunction)
	return this_
}

func InitGenerateExports(this_ *GenerateExports, compiler AbstractCompiler, allowNonGlobalExports bool, exportSymbolFunction string, exportPropertyFunction string) {
	_ = compiler
	_ = allowNonGlobalExports
	_ = exportSymbolFunction
	_ = exportPropertyFunction
	this_.initFields_GenerateExports()
	j2gort.CheckNotNullIface(compiler)
	this_.compiler = compiler
	this_.allowNonGlobalExports = allowNonGlobalExports
	this_.exportSymbolFunction = exportSymbolFunction
	this_.exportPropertyFunction = exportPropertyFunction
}

// Java method: GenerateExports.getExportedVariableNames()
func (this_ *GenerateExports) GetExportedVariableNames() j2gort.Ext_java_util_Set {
	return this_.exportedVariables
}

// Java method: GenerateExports.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *GenerateExports) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	var findExportableNodes *FindExportableNodes = NewFindExportableNodes(this_.compiler, this_.allowNonGlobalExports)
	_ = findExportableNodes
	NodeTraversal_Traverse(this_.compiler, root, findExportableNodes)
	var exports j2gort.Ext_java_util_Map = j2gort.AsIface[j2gort.Ext_java_util_Map](findExportableNodes.GetExports())
	_ = exports
	var es6Exports j2gort.Ext_java_util_Map = j2gort.AsIface[j2gort.Ext_java_util_Map](findExportableNodes.GetEs6ClassExports())
	_ = es6Exports
	var localExports j2gort.Ext_java_util_Set = j2gort.AsIface[j2gort.Ext_java_util_Set](findExportableNodes.GetLocalExports())
	_ = localExports
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(localExports))
		for cur_1.HasNext() {
			export := cur_1.NextStr()
			_ = export
			this_.priv_GenerateExports_addExtern_3050(export)
		}
	}
	if !(exports.IsEmpty()) || !(es6Exports.IsEmpty()) {
		if !(this_.priv_GenerateExports_validateExportMethodsIncluded_3046(exports, es6Exports, root)) {
			return
		}
	}
	{
		cur_2 := j2gort.NewMapEntryCursor(es6Exports)
		for cur_2.NextEntry() {
			this_.priv_GenerateExports_addExportForEs6Method_3052(pkg_rhino.Node(cur_2.KeyNode()), j2gort.CastString(cur_2.GetValue()))
		}
	}
	{
		cur_3 := j2gort.NewMapEntryCursor(exports)
		for cur_3.NextEntry() {
			var export_2 string = cur_3.KeyStr()
			_ = export_2
			var context pkg_rhino.Node = j2gort.CastNode[pkg_rhino.Node](cur_3.GetValue())
			_ = context
			this_.priv_GenerateExports_addExportMethod_3053(exports, export_2, context)
		}
	}
}

// Java method: GenerateExports.validateExportMethodsIncluded(java.util.Map<java.lang.String,com.google.javascript.rhino.Node> exports, java.util.Map<com.google.javascript.rhino.Node,java.lang.String> es6Exports, com.google.javascript.rhino.Node root)
func (this_ *GenerateExports) priv_GenerateExports_validateExportMethodsIncluded_3046(exports j2gort.Ext_java_util_Map, es6Exports j2gort.Ext_java_util_Map, root pkg_rhino.Node) bool {
	_ = exports
	_ = es6Exports
	_ = root
	var errorLocation pkg_rhino.Node = pkg_rhino.Node(0)
	_ = errorLocation
	if !(exports.IsEmpty()) {
		errorLocation = j2gort.CastNode[pkg_rhino.Node](j2gort.AsIface[j2gort.Ext_java_util_Collection](exports.Values()).Stream().FindFirst().Get())
	} else {
		errorLocation = j2gort.CastNode[pkg_rhino.Node](j2gort.AsIface[j2gort.Ext_java_util_Set](es6Exports.KeySet_0()).Stream().FindFirst().Get())
	}
	if (this_.exportSymbolFunction == j2gort.NullString) || (this_.exportPropertyFunction == j2gort.NullString) {
		this_.compiler.Report_1(JSError_Make_3(errorLocation, StaticGet_GenerateExports_MISSING_EXPORT_CONVENTION(), []string{}))
		return false
	}
	if !(this_.priv_GenerateExports_includesExportMethods_3047(root.GetParent())) {
		this_.compiler.Report_1(JSError_Make_3(errorLocation, StaticGet_GenerateExports_MISSING_GOOG_FOR_EXPORT(), []string{this_.exportSymbolFunction, this_.exportPropertyFunction}))
		this_.priv_GenerateExports_declareExportMethodsInExterns_3048(errorLocation)
		return true
	}
	return true
}

// Java method: GenerateExports.includesExportMethods(com.google.javascript.rhino.Node root)
func (this_ *GenerateExports) priv_GenerateExports_includesExportMethods_3047(root pkg_rhino.Node) bool {
	_ = root
	var globalScope *Scope = NewSyntacticScopeCreator_1(this_.compiler).CreateScope_2_Node_Scope(root, nil)
	_ = globalScope
	var exportPropertyRoot string = NodeUtil_GetRootOfQualifiedName_1_string(this_.exportPropertyFunction)
	_ = exportPropertyRoot
	var exportSymbolRoot string = NodeUtil_GetRootOfQualifiedName_1_string(this_.exportSymbolFunction)
	_ = exportSymbolRoot
	return (globalScope.HasSlot(exportPropertyRoot) && globalScope.HasSlot(exportSymbolRoot))
}

// Java method: GenerateExports.declareExportMethodsInExterns(com.google.javascript.rhino.Node srcref)
func (this_ *GenerateExports) priv_GenerateExports_declareExportMethodsInExterns_3048(srcref pkg_rhino.Node) {
	_ = srcref
	var exportPropertyRoot string = NodeUtil_GetRootOfQualifiedName_1_string(this_.exportPropertyFunction)
	_ = exportPropertyRoot
	var exportSymbolRoot string = NodeUtil_GetRootOfQualifiedName_1_string(this_.exportSymbolFunction)
	_ = exportSymbolRoot
	if j2gort.CheckNotNullString(exportPropertyRoot) == exportSymbolRoot {
		this_.priv_GenerateExports_declareSyntheticExternsVar_3049(exportPropertyRoot, srcref)
	} else {
		this_.priv_GenerateExports_declareSyntheticExternsVar_3049(exportPropertyRoot, srcref)
		this_.priv_GenerateExports_declareSyntheticExternsVar_3049(exportSymbolRoot, srcref)
	}
}

// Java method: GenerateExports.declareSyntheticExternsVar(java.lang.String name, com.google.javascript.rhino.Node srcref)
func (this_ *GenerateExports) priv_GenerateExports_declareSyntheticExternsVar_3049(name string, srcref pkg_rhino.Node) {
	_ = name
	_ = srcref
	var syntheticInput *CompilerInput = this_.compiler.GetSynthesizedExternsInput()
	_ = syntheticInput
	var syntheticVarRoot pkg_rhino.Node = syntheticInput.GetAstRoot(this_.compiler)
	_ = syntheticVarRoot
	var varDeclaration pkg_rhino.Node = pkg_rhino.IR_Var_1(pkg_rhino.IR_Name(name)).SrcrefTree(srcref)
	_ = varDeclaration
	varDeclaration.SetStaticSourceFile(pkg_rhino.ToIface_StaticSourceFile(this_.compiler.GetSynthesizedExternsInput().GetSourceFile()))
	syntheticVarRoot.AddChildToBack(varDeclaration)
	this_.compiler.ReportChangeToEnclosingScope(varDeclaration)
}

// Java method: GenerateExports.addExtern(java.lang.String export)
func (this_ *GenerateExports) priv_GenerateExports_addExtern_3050(export string) {
	_ = export
	var objectPrototype pkg_rhino.Node = NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, "Object.prototype")
	_ = objectPrototype
	var objCtor pkg_rhino.JSType = this_.compiler.GetTypeRegistry().GetNativeType(pkg_rhino.StaticGet_JSTypeNative_OBJECT_FUNCTION_TYPE())
	_ = objCtor
	objectPrototype.GetFirstChild().SetJSType(objCtor)
	var propstmt pkg_rhino.Node = pkg_rhino.IR_ExprResult(pkg_rhino.IR_Getprop_2(objectPrototype, export))
	_ = propstmt
	propstmt.SrcrefTree(this_.priv_GenerateExports_getSynthesizedExternsRoot_3057())
	propstmt.SetOriginalName(export)
	this_.priv_GenerateExports_getSynthesizedExternsRoot_3057().AddChildToBack(propstmt)
	this_.compiler.ReportChangeToEnclosingScope(propstmt)
}

// Java method: GenerateExports.recordExportSymbol(java.lang.String qname)
func (this_ *GenerateExports) priv_GenerateExports_recordExportSymbol_3051(qname string) {
	_ = qname
	var dot int32 = j2gort.JStrIndexOfCp(qname, int32(uint16(46)))
	_ = dot
	if dot == int32(-1) {
		j2gort.SetAddStr(this_.exportedVariables, qname)
	} else {
		j2gort.SetAddStr(this_.exportedVariables, j2gort.JStrSub(qname, int32(0), dot))
	}
}

// Java method: GenerateExports.addExportForEs6Method(com.google.javascript.rhino.Node memberFunction, java.lang.String ownerName)
func (this_ *GenerateExports) priv_GenerateExports_addExportForEs6Method_3052(memberFunction pkg_rhino.Node, ownerName string) {
	_ = memberFunction
	_ = ownerName
	if !(memberFunction.IsMemberFunctionDef()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(memberFunction))
	}
	if !(!(j2gort.CheckNotNullString(ownerName) == "")) {
		j2gort.CheckArgumentStr(false, ownerName)
	}
	var fullExport string = (j2gort.StrValString(ownerName) + "." + j2gort.StrValString(memberFunction.GetString()))
	_ = fullExport
	this_.priv_GenerateExports_addExportPropertyCall_3054(ownerName, memberFunction, fullExport, memberFunction.GetString())
}

// Java method: GenerateExports.addExportMethod(java.util.Map<java.lang.String,com.google.javascript.rhino.Node> exports, java.lang.String export, com.google.javascript.rhino.Node context)
func (this_ *GenerateExports) priv_GenerateExports_addExportMethod_3053(exports j2gort.Ext_java_util_Map, export string, context pkg_rhino.Node) {
	_ = exports
	_ = export
	_ = context
	var methodOwnerName string = j2gort.NullString
	_ = methodOwnerName
	var isEs5StylePrototypeAssignment bool = false
	_ = isEs5StylePrototypeAssignment
	var propertyName string = j2gort.NullString
	_ = propertyName
	var child pkg_rhino.Node = context.GetFirstChild()
	_ = child
	if (child != 0) && child.IsGetProp() {
		var node pkg_rhino.Node = context.GetFirstChild()
		_ = node
		var ownerNode pkg_rhino.Node = node.GetFirstChild()
		_ = ownerNode
		methodOwnerName = ownerNode.GetQualifiedName()
		if ownerNode.IsGetProp() && func() bool {
			sr_4 := ownerNode.GetString()
			se_5 := "prototype"
			return j2gort.CheckNotNullString(sr_4) == se_5
		}() {
			isEs5StylePrototypeAssignment = true
		}
		propertyName = node.GetString()
	}
	var useExportSymbol bool = true
	_ = useExportSymbol
	if isEs5StylePrototypeAssignment {
		useExportSymbol = false
	} else if (methodOwnerName != j2gort.NullString) && j2gort.MapContainsKeyStr(exports, methodOwnerName) {
		useExportSymbol = false
	}
	if useExportSymbol {
		this_.priv_GenerateExports_addExportSymbolCall_3055(export, context)
	} else {
		this_.priv_GenerateExports_addExportPropertyCall_3054(methodOwnerName, context, export, propertyName)
	}
}

// Java method: GenerateExports.addExportPropertyCall(java.lang.String methodOwnerName, com.google.javascript.rhino.Node context, java.lang.String export, java.lang.String propertyName)
func (this_ *GenerateExports) priv_GenerateExports_addExportPropertyCall_3054(methodOwnerName string, context pkg_rhino.Node, export string, propertyName string) {
	_ = methodOwnerName
	_ = context
	_ = export
	_ = propertyName
	j2gort.CheckNotNullString(methodOwnerName)
	var call pkg_rhino.Node = pkg_rhino.IR_Call(NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.compiler, this_.exportPropertyFunction, context, this_.exportPropertyFunction), []pkg_rhino.Node{NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.compiler, methodOwnerName, context, this_.exportPropertyFunction), pkg_rhino.IR_String(propertyName), NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.compiler, export, context, this_.exportPropertyFunction)})
	_ = call
	var expression pkg_rhino.Node = pkg_rhino.IR_ExprResult(call).SrcrefTreeIfMissing(context)
	_ = expression
	this_.priv_GenerateExports_addStatement_3056(context, expression)
}

// Java method: GenerateExports.addExportSymbolCall(java.lang.String export, com.google.javascript.rhino.Node context)
func (this_ *GenerateExports) priv_GenerateExports_addExportSymbolCall_3055(export string, context pkg_rhino.Node) {
	_ = export
	_ = context
	this_.priv_GenerateExports_recordExportSymbol_3051(export)
	var call pkg_rhino.Node = pkg_rhino.IR_Call(NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.compiler, this_.exportSymbolFunction, context, export), []pkg_rhino.Node{pkg_rhino.IR_String(export), NodeUtil_NewQName_4_AbstractCompiler_string_Node_string(this_.compiler, export, context, export)})
	_ = call
	var expression pkg_rhino.Node = pkg_rhino.IR_ExprResult(call).SrcrefTreeIfMissing(context)
	_ = expression
	this_.priv_GenerateExports_addStatement_3056(context, expression)
}

// Java method: GenerateExports.addStatement(com.google.javascript.rhino.Node context, com.google.javascript.rhino.Node stmt)
func (this_ *GenerateExports) priv_GenerateExports_addStatement_3056(context pkg_rhino.Node, stmt pkg_rhino.Node) {
	_ = context
	_ = stmt
	var convention CodingConvention = this_.compiler.GetCodingConvention()
	_ = convention
	var n pkg_rhino.Node = context
	_ = n
	var exprRoot pkg_rhino.Node = n
	_ = exprRoot
	for !(NodeUtil_IsStatementBlock(exprRoot.GetParent())) {
		exprRoot = exprRoot.GetParent()
	}
	for {
		var next pkg_rhino.Node = exprRoot.GetNext()
		_ = next
		if ((next != 0) && NodeUtil_IsExprCall(next)) && (convention.GetClassesDefinedByCall(next.GetFirstChild()) != nil) {
			exprRoot = next
		} else {
			break
		}
	}
	stmt.InsertAfter(exprRoot)
	this_.compiler.ReportChangeToEnclosingScope(stmt)
}

// Java method: GenerateExports.getSynthesizedExternsRoot()
func (this_ *GenerateExports) priv_GenerateExports_getSynthesizedExternsRoot_3057() pkg_rhino.Node {
	return this_.compiler.GetSynthesizedExternsInput().GetAstRoot(this_.compiler)
}
