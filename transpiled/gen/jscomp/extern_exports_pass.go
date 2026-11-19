/*
 * Copyright 2009 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ExternExportsPass.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_extern_exports_pass_1 j2gort.Ext_java_util_function_Function

func init() {
	lambdaSingleton_extern_exports_pass_1 = j2gort.NewLambda_Ext_java_util_function_Function(func(mrp0_ any) any {
		return j2gort.BoxString(j2gort.AsIface[ExternExportsPass_Export](mrp0_).GetExportedPath())
	})
}

// Java class: com.google.javascript.jscomp.ExternExportsPass
type ExternExportsPass struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	exports                     j2gort.Ext_java_util_List                         `gson:"exports"`
	definitionMap               j2gort.OrderedMap[string, pkg_rhino.Node]         `gson:"definitionMap"`
	compiler                    AbstractCompiler                                  `gson:"compiler"`
	mappedPaths                 j2gort.OrderedMap[string, string]                 `gson:"mappedPaths"`
	alreadyExportedPaths        j2gort.OrderedSet[string]                         `gson:"alreadyExportedPaths"`
	exportSymbolFunctionNames   j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"exportSymbolFunctionNames"`
	exportPropertyFunctionNames j2gort.Ext_com_google_common_collect_ImmutableSet `gson:"exportPropertyFunctionNames"`
	externsRoot                 pkg_rhino.Node                                    `gson:"externsRoot"`
}

func (this_ *ExternExportsPass) BindSelf_(s any) {
	_ = s
}
func (this_ *ExternExportsPass) IsJavaNil_() bool { return this_ == nil }
func (*ExternExportsPass) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ExternExportsPass", "ExternExportsPass"
}
func (this_ *ExternExportsPass) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*ExternExportsPass) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ExternExportsPass", "com.google.javascript.jscomp.ExternExportsPass", func(v_ any) bool {
		return j2gort.InstanceOf[*ExternExportsPass](v_)
	})
}
func (*ExternExportsPass) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {}
func (*ExternExportsPass) Implements_jscomp_CompilerPass()                            {}
func (*ExternExportsPass) Implements_jscomp_NodeTraversal_Callback()                  {}

var static_ExternExportsPass_Q_NAME_JOINER *j2gort.Ext_com_google_common_base_Joiner

func StaticGet_ExternExportsPass_Q_NAME_JOINER() *j2gort.Ext_com_google_common_base_Joiner {
	clinit_ExternExportsPass()
	return static_ExternExportsPass_Q_NAME_JOINER
}
func StaticSet_ExternExportsPass_Q_NAME_JOINER(v *j2gort.Ext_com_google_common_base_Joiner) {
	clinit_ExternExportsPass()
	static_ExternExportsPass_Q_NAME_JOINER = v
}

var static_ExternExportsPass_Q_NAME_SPLITTER *j2gort.Ext_com_google_common_base_Splitter

func StaticGet_ExternExportsPass_Q_NAME_SPLITTER() *j2gort.Ext_com_google_common_base_Splitter {
	clinit_ExternExportsPass()
	return static_ExternExportsPass_Q_NAME_SPLITTER
}
func StaticSet_ExternExportsPass_Q_NAME_SPLITTER(v *j2gort.Ext_com_google_common_base_Splitter) {
	clinit_ExternExportsPass()
	static_ExternExportsPass_Q_NAME_SPLITTER = v
}

var clinitDone_ExternExportsPass uint32
var clinitState_ExternExportsPass j2gort.ClinitState

func clinit_ExternExportsPass() {
	if atomic.LoadUint32(&clinitDone_ExternExportsPass) != 0 {
		return
	}
	clinitSlow_ExternExportsPass()
}

func clinitSlow_ExternExportsPass() {
	if !j2gort.EnterClinit(&clinitState_ExternExportsPass, &clinitDone_ExternExportsPass) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ExternExportsPass, &clinitDone_ExternExportsPass)
	static_ExternExportsPass_Q_NAME_JOINER = j2gort.Ext_com_google_common_base_Joiner_On_1_uint16(uint16(46))
	static_ExternExportsPass_Q_NAME_SPLITTER = j2gort.Ext_com_google_common_base_Splitter_On_1_uint16(uint16(46))
}

func (this_ *ExternExportsPass) initFields_ExternExportsPass() {
	_ = this_
}

// Java constructor: ExternExportsPass(com.google.javascript.jscomp.AbstractCompiler compiler)
func NewExternExportsPass(compiler AbstractCompiler) *ExternExportsPass {
	clinit_ExternExportsPass()
	this_ := &ExternExportsPass{}
	InitExternExportsPass(this_, compiler)
	return this_
}

func InitExternExportsPass(this_ *ExternExportsPass, compiler AbstractCompiler) {
	_ = compiler
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_ExternExportsPass()
	this_.exports = j2gort.NewExt_java_util_ArrayList()
	this_.compiler = compiler
	this_.definitionMap = j2gort.NewOrderedMap[string, pkg_rhino.Node]()
	this_.externsRoot = pkg_rhino.IR_Script_0()
	this_.alreadyExportedPaths = j2gort.NewOrderedSet[string]()
	this_.mappedPaths = j2gort.NewOrderedMap[string, string]()
	this_.priv_ExternExportsPass_initExportMethods_2817()
}

// Java static method: ExternExportsPass.computePathPrefixes(java.lang.String path)
func ExternExportsPass_ComputePathPrefixes(path string) j2gort.Ext_com_google_common_collect_ImmutableList {
	clinit_ExternExportsPass()
	_ = path
	var pieces j2gort.Ext_java_util_List = StaticGet_ExternExportsPass_Q_NAME_SPLITTER().SplitToList(j2gort.BoxString(path))
	_ = pieces
	var pathPrefixes []string = []string{}
	_ = pathPrefixes
	var partial string = j2gort.ListGetStr(pieces, int32(0))
	_ = partial
	pathPrefixes = append(pathPrefixes, j2gort.CheckNotNullString(partial))
	{
		var i int32 = int32(1)
		_ = i
		for ; i < pieces.Size(); i = i + 1 {
			partial = StaticGet_ExternExportsPass_Q_NAME_JOINER().Join_3(j2gort.BoxString(partial), j2gort.BoxString(j2gort.ListGetStr(pieces, i)), []any{})
			pathPrefixes = append(pathPrefixes, j2gort.CheckNotNullString(partial))
		}
	}
	return j2gort.ImmListFromTypedSlice(pathPrefixes)
}

// Java method: ExternExportsPass.initExportMethods()
func (this_ *ExternExportsPass) priv_ExternExportsPass_initExportMethods_2817() {
	var convention CodingConvention = this_.compiler.GetCodingConvention()
	_ = convention
	this_.exportSymbolFunctionNames = j2gort.ImmSetOfStr2(convention.GetExportSymbolFunction(), "google_exportSymbol")
	this_.exportPropertyFunctionNames = j2gort.ImmSetOfStr2(convention.GetExportPropertyFunction(), "google_exportProperty")
}

// Java method: ExternExportsPass.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ExternExportsPass) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	NodeTraversal_Traverse(this_.compiler, root, this_)
	var sorted j2gort.Ext_java_util_Set = j2gort.NewExt_java_util_TreeSet_1_Ext_java_util_Comparator(j2gort.Ext_java_util_Comparator_Comparing_1(lambdaSingleton_extern_exports_pass_1))
	_ = sorted
	sorted.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.exports))
	{
		cur_1 := j2gort.NewIterCursor(sorted)
		for cur_1.HasNext() {
			rawElem_2 := cur_1.Next()
			export := j2gort.AsIface[ExternExportsPass_Export](rawElem_2)
			_ = export
			export.GenerateExterns()
		}
	}
	this_.priv_ExternExportsPass_setGeneratedExternsOnCompiler_2818()
}

// Java method: ExternExportsPass.setGeneratedExternsOnCompiler()
func (this_ *ExternExportsPass) priv_ExternExportsPass_setGeneratedExternsOnCompiler_2818() {
	var builder *CodePrinter_Builder = NewCodePrinter_Builder(this_.externsRoot).SetPrettyPrint(true).SetOutputTypes(true).SetTypeRegistry(this_.compiler.GetTypeRegistry())
	_ = builder
	this_.compiler.SetExternExports(("/**\n * @fileoverview Generated externs.\n * @externs\n */\n" + j2gort.StrValString(builder.Build())))
}

// Java method: ExternExportsPass.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ExternExportsPass) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	this_.priv_ExternExportsPass_lookForQnameDefinition_2819(n)
	this_.priv_ExternExportsPass_lookForAtExportOnThisDotProperty_2823(t, n)
	this_.priv_ExternExportsPass_lookForSymbolExportCall_2820(n)
	this_.priv_ExternExportsPass_lookForPropertyExportCall_2821(n)
}

// Java method: ExternExportsPass.lookForQnameDefinition(com.google.javascript.rhino.Node n)
func (this_ *ExternExportsPass) priv_ExternExportsPass_lookForQnameDefinition_2819(n pkg_rhino.Node) {
	_ = n
	if n.IsClass() {
		if NodeUtil_IsClassDeclaration(n) {
			this_.definitionMap.Put(n.GetFirstChild().GetString(), n)
		}
	} else if n.IsFunction() {
		if NodeUtil_IsFunctionDeclaration(n) {
			this_.definitionMap.Put(n.GetFirstChild().GetString(), n)
		}
	} else if n.IsAssign() {
		var lhs pkg_rhino.Node = n.GetFirstChild()
		_ = lhs
		if lhs.IsQualifiedName() && (!(lhs.HasChildren()) || !(lhs.GetFirstChild().IsThis())) {
			this_.definitionMap.Put(lhs.GetQualifiedName(), n.GetLastChild())
		}
	} else if n.IsName() {
		var parent pkg_rhino.Node = j2gort.CheckNotNullNode(n.GetParent())
		_ = parent
		if NodeUtil_IsNameDeclaration(parent) {
			var value pkg_rhino.Node = n.GetFirstChild()
			_ = value
			if value != 0 {
				this_.definitionMap.Put(n.GetString(), value)
			}
		}
	} else if n.IsMemberFunctionDef() {
		var lvalueName string = NodeUtil_GetBestLValueName(n)
		_ = lvalueName
		if lvalueName != j2gort.NullString {
			this_.definitionMap.Put(lvalueName, n.GetOnlyChild())
		}
	}
}

// Java method: ExternExportsPass.lookForSymbolExportCall(com.google.javascript.rhino.Node n)
func (this_ *ExternExportsPass) priv_ExternExportsPass_lookForSymbolExportCall_2820(n pkg_rhino.Node) {
	_ = n
	if !(this_.priv_ExternExportsPass_isCallToOneOf_2822(n, this_.exportSymbolFunctionNames)) {
		return
	}
	if !(n.HasXChildren(int32(3))) {
		return
	}
	var thisNode pkg_rhino.Node = n.GetFirstChild()
	_ = thisNode
	var nameArg pkg_rhino.Node = thisNode.GetNext()
	_ = nameArg
	var valueArg pkg_rhino.Node = nameArg.GetNext()
	_ = valueArg
	if !(nameArg.IsStringLit()) {
		return
	}
	this_.exports.Add_1(NewExternExportsPass_SymbolExport(this_, nameArg.GetString(), valueArg))
}

// Java method: ExternExportsPass.lookForPropertyExportCall(com.google.javascript.rhino.Node n)
func (this_ *ExternExportsPass) priv_ExternExportsPass_lookForPropertyExportCall_2821(n pkg_rhino.Node) {
	_ = n
	if !(this_.priv_ExternExportsPass_isCallToOneOf_2822(n, this_.exportPropertyFunctionNames)) {
		return
	}
	if !(n.HasXChildren(int32(4))) {
		return
	}
	var thisNode pkg_rhino.Node = n.GetFirstChild()
	_ = thisNode
	var objectArg pkg_rhino.Node = thisNode.GetNext()
	_ = objectArg
	var nameArg pkg_rhino.Node = objectArg.GetNext()
	_ = nameArg
	var valueArg pkg_rhino.Node = nameArg.GetNext()
	_ = valueArg
	if !(objectArg.IsQualifiedName()) {
		return
	}
	if !(nameArg.IsStringLit()) {
		return
	}
	this_.exports.Add_1(NewExternExportsPass_PropertyExport(this_, objectArg.GetQualifiedName(), nameArg.GetString(), valueArg))
}

// Java method: ExternExportsPass.isCallToOneOf(com.google.javascript.rhino.Node n, com.google.common.collect.ImmutableSet<java.lang.String> functionQnames)
func (this_ *ExternExportsPass) priv_ExternExportsPass_isCallToOneOf_2822(n pkg_rhino.Node, functionQnames j2gort.Ext_com_google_common_collect_ImmutableSet) bool {
	_ = n
	_ = functionQnames
	if !(n.IsCall()) {
		return false
	} else {
		var callee pkg_rhino.Node = n.GetFirstChild()
		_ = callee
		return (callee.IsQualifiedName() && j2gort.SetContainsStr(functionQnames, callee.GetQualifiedName()))
	}
}

// Java method: ExternExportsPass.lookForAtExportOnThisDotProperty(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node thisDotPropName)
func (this_ *ExternExportsPass) priv_ExternExportsPass_lookForAtExportOnThisDotProperty_2823(t *NodeTraversal, thisDotPropName pkg_rhino.Node) {
	_ = t
	_ = thisDotPropName
	if !(thisDotPropName.IsGetProp()) || !(thisDotPropName.GetFirstChild().IsThis()) {
		return
	}
	var jsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(thisDotPropName)
	_ = jsdoc
	if (jsdoc == nil) || !(jsdoc.IsExport()) {
		return
	}
	var constructorNode pkg_rhino.Node = t.GetEnclosingFunction()
	_ = constructorNode
	if !(NodeUtil_IsConstructor(constructorNode)) {
		return
	}
	var classNode pkg_rhino.Node
	_ = classNode
	if NodeUtil_IsEs6Constructor(constructorNode) {
		classNode = NodeUtil_GetEnclosingClass(constructorNode)
	} else {
		classNode = constructorNode
	}
	var className string = NodeUtil_GetName(classNode)
	_ = className
	var propertyName string = thisDotPropName.GetString()
	_ = propertyName
	var prototypeName string = (j2gort.StrValString(className) + ".prototype")
	_ = prototypeName
	var propertyNameNode pkg_rhino.Node = NodeUtil_NewQName_2_AbstractCompiler_string(this_.compiler, ("this." + j2gort.StrValString(propertyName)))
	_ = propertyNameNode
	this_.exports.Add_1(NewExternExportsPass_PropertyExport(this_, prototypeName, propertyName, propertyNameNode))
}

// Java polymorphic class interface: com.google.javascript.jscomp.ExternExportsPass.Export
type ExternExportsPass_Export interface {
	Implements_jscomp_ExternExportsPass_Export()
	IsJavaNil_() bool
	GetExternExportsPass_ExportBase() *ExternExportsPass_ExportBase
	GenerateExterns()
	GetExportedPath() string
	AppendExtern(string, pkg_rhino.Node)
	GetValue() pkg_rhino.Node
}

func ToIface_ExternExportsPass_Export[T interface {
	ExternExportsPass_Export
	comparable
}](v T) ExternExportsPass_Export {
	var z T
	if v == z {
		return nil
	}
	return v
}

// Java class: com.google.javascript.jscomp.ExternExportsPass.Export
type ExternExportsPass_ExportBase struct {
	self_      ExternExportsPass_Export
	outer_     *ExternExportsPass
	symbolName string         `gson:"symbolName"`
	value      pkg_rhino.Node `gson:"value"`
	idHash_    atomic.Int32
}

func (this_ *ExternExportsPass_ExportBase) BindSelf_(s any) {
	_ = s
	this_.self_ = s.(ExternExportsPass_Export)
}
func (this_ *ExternExportsPass_ExportBase) IsJavaNil_() bool { return this_ == nil }
func (*ExternExportsPass_ExportBase) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ExternExportsPass$Export", "Export"
}
func (this_ *ExternExportsPass_ExportBase) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ExternExportsPass.Export", "com.google.javascript.jscomp.ExternExportsPass$Export", func(v_ any) bool {
		return j2gort.InstanceOf[ExternExportsPass_Export](v_)
	})
}
func (b *ExternExportsPass_ExportBase) GetExternExportsPass_ExportBase() *ExternExportsPass_ExportBase {
	return b
}
func (this_ *ExternExportsPass_ExportBase) SetOuter_(o *ExternExportsPass)        { this_.outer_ = o }
func (this_ *ExternExportsPass_ExportBase) GetOuter_() *ExternExportsPass         { return this_.outer_ }
func (*ExternExportsPass_ExportBase) Implements_jscomp_ExternExportsPass_Export() {}

func (this_ *ExternExportsPass_ExportBase) initFields_ExternExportsPass_Export() {
	_ = this_
	this_.symbolName = j2gort.NullString
}

func InitExternExportsPass_Export(this_ *ExternExportsPass_ExportBase, symbolName string, value pkg_rhino.Node) {
	_ = symbolName
	_ = value
	this_.initFields_ExternExportsPass_Export()
	this_.symbolName = j2gort.CheckNotNullString(symbolName)
	this_.value = j2gort.CheckNotNullNode(value)
}

// Java method: Export.generateExterns()
func (this_ *ExternExportsPass_ExportBase) GenerateExterns() {
	this_.AppendExtern(this_.self_.GetExportedPath(), this_.GetValue())
}

// Java method: Export.getExportedPath()
func (this_ *ExternExportsPass_ExportBase) GetExportedPath() string {
	panic("abstract or native method invoked")
}

// Java method: Export.appendExtern(java.lang.String path, com.google.javascript.rhino.Node valueToExport)
func (this_ *ExternExportsPass_ExportBase) AppendExtern(path string, valueToExport pkg_rhino.Node) {
	_ = path
	_ = valueToExport
	var pathPrefixes []string = j2gort.ImmListToSliceNullable[string](ExternExportsPass_ComputePathPrefixes(path))
	_ = pathPrefixes
	{
		var i int32 = int32(0)
		_ = i
		for ; i < int32(len(pathPrefixes)); i = i + 1 {
			var pathPrefix string = pathPrefixes[i]
			_ = pathPrefix
			var isCompletePathPrefix bool = (i == (int32(len(pathPrefixes)) - int32(1)))
			_ = isCompletePathPrefix
			var skipPathPrefix bool = (javalang.StringEndsWith(j2gort.CheckNotNullString(pathPrefix), ".prototype") || (this_.outer_.alreadyExportedPaths.Contains(pathPrefix) && !(isCompletePathPrefix)))
			_ = skipPathPrefix
			if skipPathPrefix {
				continue
			}
			var exportedValueDefinesNewType bool = false
			_ = exportedValueDefinesNewType
			if valueToExport != 0 {
				var jsdoc *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(valueToExport)
				_ = jsdoc
				if valueToExport.IsClass() || ((jsdoc != nil) && jsdoc.ContainsTypeDefinition()) {
					exportedValueDefinesNewType = true
				}
			}
			var initializer pkg_rhino.Node = pkg_rhino.Node(0)
			_ = initializer
			var jsdoc_2 *pkg_rhino.JSDocInfo = nil
			_ = jsdoc_2
			if isCompletePathPrefix && (valueToExport != 0) {
				if valueToExport.IsFunction() {
					initializer = this_.priv_Export_createExternFunction_2828(valueToExport)
				} else if valueToExport.IsClass() {
					initializer = this_.priv_Export_createExternFunctionForEs6Class_2833(valueToExport)
				} else {
					if !(valueToExport.IsObjectLit()) {
						j2gort.CheckState(false)
					}
					initializer = this_.priv_Export_createExternObjectLit_2836(valueToExport)
				}
			} else if !(isCompletePathPrefix) && exportedValueDefinesNewType {
				jsdoc_2 = this_.priv_Export_buildNamespaceJSDoc_2835()
				initializer = this_.priv_Export_createExternObjectLit_2836(pkg_rhino.IR_Objectlit([]pkg_rhino.Node{}))
				initializer.SetJSDocInfo(nil)
			} else {
				initializer = pkg_rhino.IR_Empty()
			}
			this_.priv_Export_appendPathDefinition_2827(pathPrefix, initializer, jsdoc_2)
		}
	}
}

// Java method: Export.appendPathDefinition(java.lang.String path, com.google.javascript.rhino.Node initializer, com.google.javascript.rhino.JSDocInfo jsdoc)
func (this_ *ExternExportsPass_ExportBase) priv_Export_appendPathDefinition_2827(path string, initializer pkg_rhino.Node, jsdoc *pkg_rhino.JSDocInfo) {
	_ = path
	_ = initializer
	_ = jsdoc
	var pathDefinition pkg_rhino.Node = pkg_rhino.Node(0)
	_ = pathDefinition
	if j2gort.JStrIndexOf(path, ".") >= 0 {
		var qualifiedPath pkg_rhino.Node = NodeUtil_NewQName_2_AbstractCompiler_string(this_.outer_.compiler, path)
		_ = qualifiedPath
		if initializer.IsEmpty() {
			pathDefinition = NodeUtil_NewExpr(qualifiedPath)
		} else {
			pathDefinition = NodeUtil_NewExpr(pkg_rhino.IR_Assign(qualifiedPath, initializer))
		}
	} else {
		if initializer.IsEmpty() {
			pathDefinition = pkg_rhino.IR_Var_1(pkg_rhino.IR_Name(path))
		} else {
			pathDefinition = NodeUtil_NewVarNode_2_string_Node(path, initializer)
		}
	}
	if jsdoc != nil {
		if pathDefinition.IsExprResult() {
			pathDefinition.GetFirstChild().SetJSDocInfo(jsdoc)
		} else {
			if !(pathDefinition.IsVar()) {
				j2gort.CheckState(false)
			}
			pathDefinition.SetJSDocInfo(jsdoc)
		}
	}
	this_.outer_.externsRoot.AddChildToBack(pathDefinition)
	this_.outer_.alreadyExportedPaths.Add(path)
}

// Java method: Export.createExternFunction(com.google.javascript.rhino.Node exportedFunction)
func (this_ *ExternExportsPass_ExportBase) priv_Export_createExternFunction_2828(exportedFunction pkg_rhino.Node) pkg_rhino.Node {
	_ = exportedFunction
	var paramList pkg_rhino.Node = this_.priv_Export_createExternsParamListFromOriginalFunction_2829(exportedFunction)
	_ = paramList
	var externFunction pkg_rhino.Node = pkg_rhino.IR_Function(pkg_rhino.IR_Name(""), paramList, pkg_rhino.IR_Block_0())
	_ = externFunction
	externFunction.SetJSType(exportedFunction.GetJSType())
	return externFunction
}

// Java method: Export.createExternsParamListFromOriginalFunction(com.google.javascript.rhino.Node exportedFunction)
func (this_ *ExternExportsPass_ExportBase) priv_Export_createExternsParamListFromOriginalFunction_2829(exportedFunction pkg_rhino.Node) pkg_rhino.Node {
	_ = exportedFunction
	var originalParamList pkg_rhino.Node = NodeUtil_GetFunctionParameters(exportedFunction)
	_ = originalParamList
	var originalParamNames j2gort.Ext_java_util_List = j2gort.NewExt_java_util_ArrayList()
	_ = originalParamNames
	{
		var originalParam pkg_rhino.Node = originalParamList.GetFirstChild()
		_ = originalParam
		for ; originalParam != 0; originalParam = originalParam.GetNext() {
			originalParamNames.Add_1(j2gort.BoxString(this_.priv_Export_getOriginalNameForParam_2832(originalParam)))
		}
	}
	return this_.priv_Export_createExternsParamListFromOriginalParamList_2831(originalParamNames)
}

// Java method: Export.createExternsParamListFromFunctionType(com.google.javascript.rhino.jstype.JSType functionType)
func (this_ *ExternExportsPass_ExportBase) priv_Export_createExternsParamListFromFunctionType_2830(functionType pkg_rhino.JSType) pkg_rhino.Node {
	_ = functionType
	var emptyParamNames j2gort.Ext_java_util_List = j2gort.Ext_java_util_Collections_NCopies(functionType.AssertFunctionType().GetParameterCount(), "")
	_ = emptyParamNames
	return this_.priv_Export_createExternsParamListFromOriginalParamList_2831(emptyParamNames)
}

// Java method: Export.createExternsParamListFromOriginalParamList(java.util.List<java.lang.String> originalParamNames)
func (this_ *ExternExportsPass_ExportBase) priv_Export_createExternsParamListFromOriginalParamList_2831(originalParamNames j2gort.Ext_java_util_List) pkg_rhino.Node {
	_ = originalParamNames
	var paramList pkg_rhino.Node = pkg_rhino.IR_ParamList([]pkg_rhino.Node{})
	_ = paramList
	var nameGenerator NameGenerator = NewDefaultNameGenerator_3(j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_CopyOf_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](originalParamNames))), "", j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_0()))
	_ = nameGenerator
	{
		cur_3 := j2gort.NewIterCursor(j2gort.CheckNotNull(originalParamNames))
		for cur_3.HasNext() {
			originalParamName := cur_3.NextStr()
			_ = originalParamName
			var externParamName string
			_ = externParamName
			if j2gort.CheckNotNullString(originalParamName) == "" {
				externParamName = nameGenerator.GenerateNextName()
			} else {
				externParamName = originalParamName
			}
			paramList.AddChildToBack(pkg_rhino.IR_Name(externParamName))
		}
	}
	return paramList
}

// Java method: Export.getOriginalNameForParam(com.google.javascript.rhino.Node paramNode)
func (this_ *ExternExportsPass_ExportBase) priv_Export_getOriginalNameForParam_2832(paramNode pkg_rhino.Node) string {
	_ = paramNode
	var nameOrPatternNode pkg_rhino.Node = pkg_rhino.Node(0)
	_ = nameOrPatternNode
	if paramNode.IsRest() {
		nameOrPatternNode = paramNode.GetOnlyChild()
	} else if paramNode.IsDefaultValue() {
		nameOrPatternNode = paramNode.GetFirstChild()
	} else {
		nameOrPatternNode = paramNode
	}
	if nameOrPatternNode.IsName() {
		var originalName string = nameOrPatternNode.GetOriginalName()
		_ = originalName
		if originalName != j2gort.NullString {
			return originalName
		} else {
			return nameOrPatternNode.GetString()
		}
	} else {
		if !(nameOrPatternNode.IsDestructuringPattern()) {
			j2gort.CheckStateMsg(false, j2gort.BoxNode(nameOrPatternNode))
		}
		return ""
	}
}

// Java method: Export.createExternFunctionForEs6Class(com.google.javascript.rhino.Node exportedClass)
func (this_ *ExternExportsPass_ExportBase) priv_Export_createExternFunctionForEs6Class_2833(exportedClass pkg_rhino.Node) pkg_rhino.Node {
	_ = exportedClass
	var constructorMethodDefinition pkg_rhino.Node = NodeUtil_GetEs6ClassConstructorMemberFunctionDef(exportedClass)
	_ = constructorMethodDefinition
	if constructorMethodDefinition == 0 {
		var classJSType pkg_rhino.JSType = exportedClass.GetJSType()
		_ = classJSType
		var paramList pkg_rhino.Node = this_.priv_Export_createExternsParamListFromFunctionType_2830(classJSType)
		_ = paramList
		var externFunction pkg_rhino.Node = pkg_rhino.IR_Function(pkg_rhino.IR_Name(""), paramList, pkg_rhino.IR_Block_0())
		_ = externFunction
		externFunction.SetJSType(classJSType)
		return externFunction
	} else {
		return this_.priv_Export_createExternFunction_2828(constructorMethodDefinition.GetOnlyChild())
	}
}

// Java method: Export.buildEmptyJSDoc()
func (this_ *ExternExportsPass_ExportBase) priv_Export_buildEmptyJSDoc_2834() *pkg_rhino.JSDocInfo {
	return pkg_rhino.JSDocInfo_S_Builder().Build_1(true)
}

// Java method: Export.buildNamespaceJSDoc()
func (this_ *ExternExportsPass_ExportBase) priv_Export_buildNamespaceJSDoc_2835() *pkg_rhino.JSDocInfo {
	var builder *pkg_rhino.JSDocInfo_Builder = pkg_rhino.JSDocInfo_S_Builder()
	_ = builder
	builder.RecordConstancy()
	builder.RecordSuppressions_1(j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.ImmSetOfStr2("const", "duplicate")))
	return builder.Build_0()
}

// Java method: Export.createExternObjectLit(com.google.javascript.rhino.Node exportedObjectLit)
func (this_ *ExternExportsPass_ExportBase) priv_Export_createExternObjectLit_2836(exportedObjectLit pkg_rhino.Node) pkg_rhino.Node {
	_ = exportedObjectLit
	var lit pkg_rhino.Node = pkg_rhino.IR_Objectlit([]pkg_rhino.Node{})
	_ = lit
	lit.SetJSType(exportedObjectLit.GetJSType())
	lit.SetJSDocInfo(this_.priv_Export_buildEmptyJSDoc_2834())
	var index int32 = int32(1)
	_ = index
	{
		var child pkg_rhino.Node = exportedObjectLit.GetFirstChild()
		_ = child
		for ; child != 0; child = child.GetNext() {
			if child.IsStringKey() {
				lit.AddChildToBack(pkg_rhino.IR_Propdef(pkg_rhino.IR_StringKey_1(child.GetString()), pkg_rhino.IR_Number(float64(func() int32 { old_4 := index; nw_5 := old_4 + 1; index = nw_5; return old_4 }()))))
			}
		}
	}
	return lit
}

// Java method: Export.getValue()
func (this_ *ExternExportsPass_ExportBase) GetValue() pkg_rhino.Node {
	var qualifiedName string = this_.value.GetQualifiedName()
	_ = qualifiedName
	if qualifiedName == j2gort.NullString {
		return pkg_rhino.Node(0)
	}
	var definition pkg_rhino.Node = this_.outer_.definitionMap.Get(qualifiedName)
	_ = definition
	if definition == 0 {
		return pkg_rhino.Node(0)
	}
	if (definition.IsFunction() || definition.IsClass()) || definition.IsObjectLit() {
		return definition
	}
	return pkg_rhino.Node(0)
}

// Java class: com.google.javascript.jscomp.ExternExportsPass.SymbolExport
type ExternExportsPass_SymbolExport struct {
	ExternExportsPass_ExportBase
	outer_ *ExternExportsPass
}

func (this_ *ExternExportsPass_SymbolExport) BindSelf_(s any) {
	_ = s
	this_.ExternExportsPass_ExportBase.BindSelf_(s)
}
func (this_ *ExternExportsPass_SymbolExport) IsJavaNil_() bool { return this_ == nil }
func (*ExternExportsPass_SymbolExport) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ExternExportsPass$SymbolExport", "SymbolExport"
}
func (this_ *ExternExportsPass_SymbolExport) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ExternExportsPass_ExportBase.IdentityHash_()
}
func (*ExternExportsPass_SymbolExport) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ExternExportsPass.SymbolExport", "com.google.javascript.jscomp.ExternExportsPass$SymbolExport", func(v_ any) bool {
		return j2gort.InstanceOf[*ExternExportsPass_SymbolExport](v_)
	})
}
func (this_ *ExternExportsPass_SymbolExport) SetOuter_(o *ExternExportsPass)        { this_.outer_ = o }
func (this_ *ExternExportsPass_SymbolExport) GetOuter_() *ExternExportsPass         { return this_.outer_ }
func (*ExternExportsPass_SymbolExport) Implements_jscomp_ExternExportsPass_Export() {}

func (this_ *ExternExportsPass_SymbolExport) initFields_ExternExportsPass_SymbolExport() {
	_ = this_
}

// Java constructor: SymbolExport(java.lang.String symbolName, com.google.javascript.rhino.Node value)
func NewExternExportsPass_SymbolExport(outer_ *ExternExportsPass, symbolName string, value pkg_rhino.Node) *ExternExportsPass_SymbolExport {
	this_ := &ExternExportsPass_SymbolExport{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitExternExportsPass_SymbolExport(this_, symbolName, value)
	return this_
}

func InitExternExportsPass_SymbolExport(this_ *ExternExportsPass_SymbolExport, symbolName string, value pkg_rhino.Node) {
	_ = symbolName
	_ = value
	this_.ExternExportsPass_ExportBase.outer_ = this_.outer_
	InitExternExportsPass_Export(&this_.ExternExportsPass_ExportBase, symbolName, value)
	this_.initFields_ExternExportsPass_SymbolExport()
	var qualifiedName string = value.GetQualifiedName()
	_ = qualifiedName
	if qualifiedName != j2gort.NullString {
		this_.outer_.mappedPaths.Put(qualifiedName, symbolName)
	}
}

// Java method: SymbolExport.getExportedPath()
func (this_ *ExternExportsPass_SymbolExport) GetExportedPath() string {
	return this_.GetExternExportsPass_ExportBase().symbolName
}

// Java class: com.google.javascript.jscomp.ExternExportsPass.PropertyExport
type ExternExportsPass_PropertyExport struct {
	ExternExportsPass_ExportBase
	outer_     *ExternExportsPass
	exportPath string `gson:"exportPath"`
}

func (this_ *ExternExportsPass_PropertyExport) BindSelf_(s any) {
	_ = s
	this_.ExternExportsPass_ExportBase.BindSelf_(s)
}
func (this_ *ExternExportsPass_PropertyExport) IsJavaNil_() bool { return this_ == nil }
func (*ExternExportsPass_PropertyExport) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ExternExportsPass$PropertyExport", "PropertyExport"
}
func (this_ *ExternExportsPass_PropertyExport) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.ExternExportsPass_ExportBase.IdentityHash_()
}
func (*ExternExportsPass_PropertyExport) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ExternExportsPass.PropertyExport", "com.google.javascript.jscomp.ExternExportsPass$PropertyExport", func(v_ any) bool {
		return j2gort.InstanceOf[*ExternExportsPass_PropertyExport](v_)
	})
}
func (this_ *ExternExportsPass_PropertyExport) SetOuter_(o *ExternExportsPass)        { this_.outer_ = o }
func (this_ *ExternExportsPass_PropertyExport) GetOuter_() *ExternExportsPass         { return this_.outer_ }
func (*ExternExportsPass_PropertyExport) Implements_jscomp_ExternExportsPass_Export() {}

func (this_ *ExternExportsPass_PropertyExport) initFields_ExternExportsPass_PropertyExport() {
	_ = this_
	this_.exportPath = j2gort.NullString
}

// Java constructor: PropertyExport(java.lang.String exportPath, java.lang.String symbolName, com.google.javascript.rhino.Node value)
func NewExternExportsPass_PropertyExport(outer_ *ExternExportsPass, exportPath string, symbolName string, value pkg_rhino.Node) *ExternExportsPass_PropertyExport {
	this_ := &ExternExportsPass_PropertyExport{}
	this_.BindSelf_(this_)
	this_.outer_ = outer_
	InitExternExportsPass_PropertyExport(this_, exportPath, symbolName, value)
	return this_
}

func InitExternExportsPass_PropertyExport(this_ *ExternExportsPass_PropertyExport, exportPath string, symbolName string, value pkg_rhino.Node) {
	_ = exportPath
	_ = symbolName
	_ = value
	this_.ExternExportsPass_ExportBase.outer_ = this_.outer_
	InitExternExportsPass_Export(&this_.ExternExportsPass_ExportBase, symbolName, value)
	this_.initFields_ExternExportsPass_PropertyExport()
	this_.exportPath = j2gort.CheckNotNullString(exportPath)
}

// Java method: PropertyExport.getExportedPath()
func (this_ *ExternExportsPass_PropertyExport) GetExportedPath() string {
	{
		cur_6 := j2gort.NewIterCursor(j2gort.CheckNotNull(j2gort.Ext_com_google_common_collect_Lists_Reverse(j2gort.AsIface[j2gort.Ext_java_util_List](ExternExportsPass_ComputePathPrefixes(this_.exportPath)))))
		for cur_6.HasNext() {
			currentPath := cur_6.NextStr()
			_ = currentPath
			if !(j2gort.JStrLen(currentPath) > int32(0)) {
				j2gort.CheckState(false)
			}
			var mappedPath string = j2gort.OrderedMapGetStr(this_.outer_.mappedPaths, currentPath)
			_ = mappedPath
			if mappedPath == j2gort.NullString {
				continue
			}
			mappedPath = (j2gort.StrValString(mappedPath) + j2gort.StrValString(j2gort.JStrSubFrom(this_.exportPath, j2gort.JStrLen(currentPath))))
			return StaticGet_ExternExportsPass_Q_NAME_JOINER().Join_3(j2gort.BoxString(mappedPath), j2gort.BoxString(this_.GetExternExportsPass_ExportBase().symbolName), []any{})
		}
	}
	return StaticGet_ExternExportsPass_Q_NAME_JOINER().Join_3(j2gort.BoxString(this_.exportPath), j2gort.BoxString(this_.GetExternExportsPass_ExportBase().symbolName), []any{})
}
func (*ExternExportsPass) AlwaysShouldTraverse_()      {}
func (*ExternExportsPass) CanReleaseSyntacticScopes_() {}
