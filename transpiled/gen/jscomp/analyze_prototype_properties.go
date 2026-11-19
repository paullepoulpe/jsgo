/*
 * Copyright 2006 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/AnalyzePrototypeProperties.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	pkg_graph "github.com/paullepoulpe/jsgo/transpiled/gen/graph"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.AnalyzePrototypeProperties
type AnalyzePrototypeProperties struct {
	compiler               AbstractCompiler                                                              `gson:"compiler"`
	chunkGraph             *JSChunkGraph                                                                 `gson:"chunkGraph"`
	firstChunk             *JSChunk                                                                      `gson:"firstChunk"`
	symbolGraph            pkg_graph.LinkedDirectedGraph[*AnalyzePrototypeProperties_NameInfo, *JSChunk] `gson:"symbolGraph"`
	globalNode             *AnalyzePrototypeProperties_NameInfo                                          `gson:"globalNode"`
	externNode             *AnalyzePrototypeProperties_NameInfo                                          `gson:"externNode"`
	anonymousNode          *AnalyzePrototypeProperties_NameInfo                                          `gson:"anonymousNode"`
	propertyNameInfo       j2gort.OrderedMap[string, *AnalyzePrototypeProperties_NameInfo]               `gson:"propertyNameInfo"`
	varNameInfo            j2gort.OrderedMap[string, *AnalyzePrototypeProperties_NameInfo]               `gson:"varNameInfo"`
	idHash_                atomic.Int32
	canModifyExterns       bool `gson:"canModifyExterns"`
	anchorUnusedVars       bool `gson:"anchorUnusedVars"`
	rootScopeUsesAreGlobal bool `gson:"rootScopeUsesAreGlobal"`
}

func (this_ *AnalyzePrototypeProperties) BindSelf_(s any) {
	_ = s
}
func (this_ *AnalyzePrototypeProperties) IsJavaNil_() bool { return this_ == nil }
func (*AnalyzePrototypeProperties) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AnalyzePrototypeProperties", "AnalyzePrototypeProperties"
}
func (this_ *AnalyzePrototypeProperties) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AnalyzePrototypeProperties) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AnalyzePrototypeProperties", "com.google.javascript.jscomp.AnalyzePrototypeProperties", func(v_ any) bool {
		return j2gort.InstanceOf[*AnalyzePrototypeProperties](v_)
	})
}
func (*AnalyzePrototypeProperties) Implements_jscomp_CompilerPass() {}

var static_AnalyzePrototypeProperties_PROPERTY AnalyzePrototypeProperties_SymbolType

func StaticGet_AnalyzePrototypeProperties_PROPERTY() AnalyzePrototypeProperties_SymbolType {
	clinit_AnalyzePrototypeProperties()
	return static_AnalyzePrototypeProperties_PROPERTY
}
func StaticSet_AnalyzePrototypeProperties_PROPERTY(v AnalyzePrototypeProperties_SymbolType) {
	clinit_AnalyzePrototypeProperties()
	static_AnalyzePrototypeProperties_PROPERTY = v
}

var static_AnalyzePrototypeProperties_VAR AnalyzePrototypeProperties_SymbolType

func StaticGet_AnalyzePrototypeProperties_VAR() AnalyzePrototypeProperties_SymbolType {
	clinit_AnalyzePrototypeProperties()
	return static_AnalyzePrototypeProperties_VAR
}
func StaticSet_AnalyzePrototypeProperties_VAR(v AnalyzePrototypeProperties_SymbolType) {
	clinit_AnalyzePrototypeProperties()
	static_AnalyzePrototypeProperties_VAR = v
}

var static_AnalyzePrototypeProperties_IMPLICITLY_USED_PROPERTIES j2gort.Ext_com_google_common_collect_ImmutableSet

func StaticGet_AnalyzePrototypeProperties_IMPLICITLY_USED_PROPERTIES() j2gort.Ext_com_google_common_collect_ImmutableSet {
	clinit_AnalyzePrototypeProperties()
	return static_AnalyzePrototypeProperties_IMPLICITLY_USED_PROPERTIES
}
func StaticSet_AnalyzePrototypeProperties_IMPLICITLY_USED_PROPERTIES(v j2gort.Ext_com_google_common_collect_ImmutableSet) {
	clinit_AnalyzePrototypeProperties()
	static_AnalyzePrototypeProperties_IMPLICITLY_USED_PROPERTIES = v
}

var clinitDone_AnalyzePrototypeProperties uint32
var clinitState_AnalyzePrototypeProperties j2gort.ClinitState

func clinit_AnalyzePrototypeProperties() {
	if atomic.LoadUint32(&clinitDone_AnalyzePrototypeProperties) != 0 {
		return
	}
	clinitSlow_AnalyzePrototypeProperties()
}

func clinitSlow_AnalyzePrototypeProperties() {
	if !j2gort.EnterClinit(&clinitState_AnalyzePrototypeProperties, &clinitDone_AnalyzePrototypeProperties) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_AnalyzePrototypeProperties, &clinitDone_AnalyzePrototypeProperties)
	static_AnalyzePrototypeProperties_PROPERTY = StaticGet_AnalyzePrototypeProperties_SymbolType_PROPERTY()
	static_AnalyzePrototypeProperties_VAR = StaticGet_AnalyzePrototypeProperties_SymbolType_VAR()
	static_AnalyzePrototypeProperties_IMPLICITLY_USED_PROPERTIES = j2gort.Ext_com_google_common_collect_ImmutableSet_Of_3("length", "toString", "valueOf")
}

func (this_ *AnalyzePrototypeProperties) initFields_AnalyzePrototypeProperties() {
	_ = this_
	this_.symbolGraph = pkg_graph.LinkedDirectedGraph_CreateWithoutAnnotations[*AnalyzePrototypeProperties_NameInfo, *JSChunk]()
	this_.globalNode = NewAnalyzePrototypeProperties_NameInfo(this_, "[global]")
	this_.externNode = NewAnalyzePrototypeProperties_NameInfo(this_, "[extern]")
	this_.anonymousNode = NewAnalyzePrototypeProperties_NameInfo(this_, "[anonymous]")
	this_.propertyNameInfo = j2gort.NewOrderedMap[string, *AnalyzePrototypeProperties_NameInfo]()
	this_.varNameInfo = j2gort.NewOrderedMap[string, *AnalyzePrototypeProperties_NameInfo](

	// Java constructor: AnalyzePrototypeProperties(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.JSChunkGraph chunkGraph, boolean canModifyExterns, boolean anchorUnusedVars, boolean rootScopeUsesAreGlobal)
	)
}

func NewAnalyzePrototypeProperties(compiler AbstractCompiler, chunkGraph *JSChunkGraph, canModifyExterns bool, anchorUnusedVars bool, rootScopeUsesAreGlobal bool) *AnalyzePrototypeProperties {
	this_ := &AnalyzePrototypeProperties{}
	InitAnalyzePrototypeProperties(this_, compiler, chunkGraph, canModifyExterns, anchorUnusedVars, rootScopeUsesAreGlobal)
	return this_
}

func InitAnalyzePrototypeProperties(this_ *AnalyzePrototypeProperties, compiler AbstractCompiler, chunkGraph *JSChunkGraph, canModifyExterns bool, anchorUnusedVars bool, rootScopeUsesAreGlobal bool) {
	_ = compiler
	_ = chunkGraph
	_ = canModifyExterns
	_ = anchorUnusedVars
	_ = rootScopeUsesAreGlobal
	this_.initFields_AnalyzePrototypeProperties()
	this_.compiler = compiler
	this_.chunkGraph = chunkGraph
	this_.canModifyExterns = canModifyExterns
	this_.anchorUnusedVars = anchorUnusedVars
	this_.rootScopeUsesAreGlobal = rootScopeUsesAreGlobal
	if chunkGraph.GetChunkCount() > int32(1) {
		this_.firstChunk = chunkGraph.GetRootChunk()
	} else {
		this_.firstChunk = nil
	}
	this_.globalNode.MarkReference(nil)
	this_.externNode.MarkReference(nil)
	this_.symbolGraph.CreateNode(this_.globalNode)
	this_.symbolGraph.CreateNode(this_.externNode)
	{
		cur_1 := j2gort.NewIterCursor(j2gort.CheckNotNull(StaticGet_AnalyzePrototypeProperties_IMPLICITLY_USED_PROPERTIES()))
		for cur_1.HasNext() {
			property := cur_1.NextStr()
			_ = property
			var nameInfo *AnalyzePrototypeProperties_NameInfo = this_.priv_AnalyzePrototypeProperties_getNameInfoForName_565(property, StaticGet_AnalyzePrototypeProperties_PROPERTY())
			_ = nameInfo
			if chunkGraph == nil {
				this_.symbolGraph.Connect_3_any_any_any(this_.externNode, nil, nameInfo)
			} else {
				{
					cur_2 := j2gort.NewIterCursor(j2gort.CheckNotNull(chunkGraph.GetAllChunks()))
					for cur_2.HasNext() {
						rawElem_3 := cur_2.Next()
						chunk := j2gort.CastPtr[JSChunk](rawElem_3)
						_ = chunk
						this_.symbolGraph.Connect_3_any_any_any(this_.externNode, chunk, nameInfo)
					}
				}
			}
		}
	}
}

// Java method: AnalyzePrototypeProperties.process(com.google.javascript.rhino.Node externRoot, com.google.javascript.rhino.Node root)
func (this_ *AnalyzePrototypeProperties) Process_2(externRoot pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externRoot
	_ = root
	if !(this_.compiler.GetLifeCycleStage().IsNormalized()) {
		j2gort.CheckState(false)
	}
	if !(this_.canModifyExterns) {
		NodeTraversal_Traverse(this_.compiler, externRoot, NewAnalyzePrototypeProperties_ProcessExternProperties(this_))
	}
	NodeTraversal_Traverse(this_.compiler, root, NewAnalyzePrototypeProperties_ProcessProperties(this_))
	var t *pkg_graph.FixedPointGraphTraversal[*AnalyzePrototypeProperties_NameInfo, *JSChunk] = pkg_graph.FixedPointGraphTraversal_NewTraversal[*AnalyzePrototypeProperties_NameInfo, *JSChunk](pkg_graph.ToIface_FixedPointGraphTraversal_EdgeCallback[*AnalyzePrototypeProperties_NameInfo, *JSChunk](NewAnalyzePrototypeProperties_PropagateReferences(this_)))
	_ = t
	t.ComputeFixedPoint_2_DiGraph_any_any__Ext_java_util_Set(j2gort.AsIface[pkg_graph.DiGraph[*AnalyzePrototypeProperties_NameInfo, *JSChunk]](this_.symbolGraph), j2gort.AsIface[j2gort.Ext_java_util_Set](j2gort.Ext_com_google_common_collect_ImmutableSet_Of_2(j2gort.PtrToAny(this_.externNode), j2gort.PtrToAny(this_.globalNode))))
}

// Java method: AnalyzePrototypeProperties.getAllNameInfo()
func (this_ *AnalyzePrototypeProperties) GetAllNameInfo() j2gort.Ext_java_util_Collection {
	var result j2gort.Ext_java_util_ArrayList = j2gort.NewExt_java_util_ArrayList_1_Ext_java_util_Collection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.propertyNameInfo.LiveValues()))
	_ = result
	result.AddAll_1(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.varNameInfo.LiveValues()))
	return result
}

// Java method: AnalyzePrototypeProperties.getNameInfoForName(java.lang.String name, com.google.javascript.jscomp.AnalyzePrototypeProperties.SymbolType type)
func (this_ *AnalyzePrototypeProperties) priv_AnalyzePrototypeProperties_getNameInfoForName_565(name string, type_ AnalyzePrototypeProperties_SymbolType) *AnalyzePrototypeProperties_NameInfo {
	_ = name
	_ = type_
	if type_ == StaticGet_AnalyzePrototypeProperties_PROPERTY() {
		var nameInfo *AnalyzePrototypeProperties_NameInfo = this_.propertyNameInfo.Get(name)
		_ = nameInfo
		if nameInfo == nil {
			nameInfo = NewAnalyzePrototypeProperties_NameInfo(this_, name)
			this_.propertyNameInfo.Put(name, nameInfo)
			this_.symbolGraph.CreateNode(nameInfo)
		}
		return nameInfo
	} else {
		var nameInfo_2 *AnalyzePrototypeProperties_NameInfo = this_.varNameInfo.Get(name)
		_ = nameInfo_2
		if nameInfo_2 == nil {
			nameInfo_2 = NewAnalyzePrototypeProperties_NameInfo(this_, name)
			this_.varNameInfo.Put(name, nameInfo_2)
			this_.symbolGraph.CreateNode(nameInfo_2)
		}
		return nameInfo_2
	}
}

// Java class: com.google.javascript.jscomp.AnalyzePrototypeProperties.ProcessProperties
type AnalyzePrototypeProperties_ProcessProperties struct {
	outer_      *AnalyzePrototypeProperties
	symbolStack j2gort.Deque[*AnalyzePrototypeProperties_NameContext] `gson:"symbolStack"`
	idHash_     atomic.Int32
}

func (this_ *AnalyzePrototypeProperties_ProcessProperties) BindSelf_(s any) {
	_ = s
}
func (this_ *AnalyzePrototypeProperties_ProcessProperties) IsJavaNil_() bool { return this_ == nil }
func (*AnalyzePrototypeProperties_ProcessProperties) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AnalyzePrototypeProperties$ProcessProperties", "ProcessProperties"
}
func (this_ *AnalyzePrototypeProperties_ProcessProperties) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AnalyzePrototypeProperties_ProcessProperties) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AnalyzePrototypeProperties.ProcessProperties", "com.google.javascript.jscomp.AnalyzePrototypeProperties$ProcessProperties", func(v_ any) bool {
		return j2gort.InstanceOf[*AnalyzePrototypeProperties_ProcessProperties](v_)
	})
}
func (this_ *AnalyzePrototypeProperties_ProcessProperties) SetOuter_(o *AnalyzePrototypeProperties) {
	this_.outer_ = o
}
func (this_ *AnalyzePrototypeProperties_ProcessProperties) GetOuter_() *AnalyzePrototypeProperties {
	return this_.outer_
}
func (*AnalyzePrototypeProperties_ProcessProperties) Implements_jscomp_NodeTraversal_ScopedCallback() {
}
func (*AnalyzePrototypeProperties_ProcessProperties) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *AnalyzePrototypeProperties_ProcessProperties) initFields_AnalyzePrototypeProperties_ProcessProperties() {
	_ = this_
	this_.symbolStack = j2gort.Deque[*AnalyzePrototypeProperties_NameContext]{}
}

// Java constructor: ProcessProperties()
func NewAnalyzePrototypeProperties_ProcessProperties(outer_ *AnalyzePrototypeProperties) *AnalyzePrototypeProperties_ProcessProperties {
	this_ := &AnalyzePrototypeProperties_ProcessProperties{}
	this_.outer_ = outer_
	InitAnalyzePrototypeProperties_ProcessProperties(this_)
	return this_
}

func InitAnalyzePrototypeProperties_ProcessProperties(this_ *AnalyzePrototypeProperties_ProcessProperties) {
	this_.initFields_AnalyzePrototypeProperties_ProcessProperties()
}

// Java method: ProcessProperties.enterScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *AnalyzePrototypeProperties_ProcessProperties) EnterScope(t *NodeTraversal) {
	_ = t
	var n pkg_rhino.Node = t.GetCurrentNode()
	_ = n
	var scope *Scope = t.GetScope()
	_ = scope
	var root pkg_rhino.Node = scope.GetRootNode_0()
	_ = root
	if root.IsFunction() {
		var propName string = this_.priv_ProcessProperties_getPrototypePropertyNameFromRValue_572(n)
		_ = propName
		if propName != j2gort.NullString {
			this_.symbolStack.Push(NewAnalyzePrototypeProperties_NameContext(this_.outer_.priv_AnalyzePrototypeProperties_getNameInfoForName_565(propName, StaticGet_AnalyzePrototypeProperties_PROPERTY()), scope))
		} else if this_.priv_ProcessProperties_isGlobalFunctionDeclaration_570(t, n) {
			var parent pkg_rhino.Node = n.GetParent()
			_ = parent
			var name string
			_ = name
			if parent.IsName() {
				name = parent.GetString()
			} else {
				name = n.GetFirstChild().GetString()
			}
			this_.symbolStack.Push(NewAnalyzePrototypeProperties_NameContext(this_.outer_.priv_AnalyzePrototypeProperties_getNameInfoForName_565(name, StaticGet_AnalyzePrototypeProperties_VAR()), j2gort.CastPtr[Scope](scope.GetClosestHoistScope())))
		} else {
			this_.symbolStack.Push(NewAnalyzePrototypeProperties_NameContext(this_.outer_.anonymousNode, scope))
		}
	} else if t.InGlobalScope() {
		this_.symbolStack.Push(NewAnalyzePrototypeProperties_NameContext(this_.outer_.globalNode, scope))
	} else {
		if !(((NodeUtil_CreatesBlockScope(root) || root.IsModuleBody()) || root.IsComputedFieldDef()) || root.IsMemberFieldDef()) {
			j2gort.CheckStateMsg(false, j2gort.PtrToAny(scope))
		}
		this_.symbolStack.Push(NewAnalyzePrototypeProperties_NameContext(this_.outer_.anonymousNode, scope))
	}
}

// Java method: ProcessProperties.exitScope(com.google.javascript.jscomp.NodeTraversal t)
func (this_ *AnalyzePrototypeProperties_ProcessProperties) ExitScope(t *NodeTraversal) {
	_ = t
	this_.symbolStack.Pop()
}

// Java method: ProcessProperties.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *AnalyzePrototypeProperties_ProcessProperties) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if !(this_.outer_.rootScopeUsesAreGlobal) {
		var propName string = this_.priv_ProcessProperties_processNonFunctionPrototypeAssign_569(n, parent)
		_ = propName
		if propName != j2gort.NullString {
			this_.symbolStack.Push(NewAnalyzePrototypeProperties_NameContext(this_.outer_.priv_AnalyzePrototypeProperties_getNameInfoForName_565(propName, StaticGet_AnalyzePrototypeProperties_PROPERTY()), nil))
		}
	}
	return true
}

// Java method: ProcessProperties.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *AnalyzePrototypeProperties_ProcessProperties) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	{
		var propName string = j2gort.NullString
		_ = propName
		var lValueName string = j2gort.NullString
		_ = lValueName
		var classMembers pkg_rhino.Node = pkg_rhino.Node(0)
		_ = classMembers
		var name string = j2gort.NullString
		_ = name
		var var_ *Var = nil
		_ = var_
		switch n.TokenSwitchOrd_() {
		case int32(110) /* Token.SUPER */ :
			{
				cur_4 := this_.symbolStack.Iter()
				for cur_4.HasNext() {

					context := cur_4.Next()
					_ = context
					{
						r_6 := context.name
						r_6.referencesSuper = true
					}
					if NodeUtil_IsMethodDeclaration(context.scope.GetRootNode_0()) {
						break
					}
				}
			}
		case int32(29) /* Token.OPTCHAIN_GETPROP */ :
			this_.priv_ProcessProperties_addSymbolUse_568(n.GetString(), t.GetChunk(), StaticGet_AnalyzePrototypeProperties_PROPERTY())
		case int32(26) /* Token.GETPROP */ :
			propName = n.GetString()
			if n.IsQualifiedName() {
				if j2gort.CheckNotNullString(propName) == "prototype" {
					if this_.priv_ProcessProperties_handlePossibleAssignmentToPrototype_574(t, n) {
						return
					}
				} else if this_.outer_.compiler.GetCodingConvention().IsExported_2(propName, false) {
					this_.priv_ProcessProperties_addGlobalUseOfSymbol_577(propName, t.GetChunk(), StaticGet_AnalyzePrototypeProperties_PROPERTY())
					return
				} else {
					if parent.IsAssign() && (n == parent.GetFirstChild()) {
						var rValueName string = this_.priv_ProcessProperties_getPrototypePropertyNameFromRValue_572(n)
						_ = rValueName
						if rValueName != j2gort.NullString {
							return
						}
					}
				}
			}
			this_.priv_ProcessProperties_addSymbolUse_568(propName, t.GetChunk(), StaticGet_AnalyzePrototypeProperties_PROPERTY())
		case int32(47) /* Token.OBJECTLIT */ :
			lValueName = NodeUtil_GetBestLValueName(NodeUtil_GetBestLValue(n))
			if (lValueName != j2gort.NullString) && javalang.StringEndsWith(j2gort.CheckNotNullString(lValueName), ".prototype") {
				return
			}
			fallthrough
		case int32(103) /* Token.OBJECT_PATTERN */ :
			{
				var propNode pkg_rhino.Node = n.GetFirstChild()
				_ = propNode
				for ; propNode != 0; propNode = propNode.GetNext() {
					switch propNode.TokenSwitchOrd_() {
					case int32(129) /* Token.COMPUTED_PROP */, int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */, int32(127) /* Token.ITER_SPREAD */, int32(128) /* Token.OBJECT_SPREAD */ :
					case int32(100) /* Token.STRING_KEY */, int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(107) /* Token.MEMBER_FUNCTION_DEF */ :
						if !(propNode.IsQuotedStringKey()) {
							this_.priv_ProcessProperties_addSymbolUse_568(propNode.GetString(), t.GetChunk(), StaticGet_AnalyzePrototypeProperties_PROPERTY())
						}
					default:
						panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Unexpected child of " + j2gort.StrVal(n.GetToken()) + ": " + j2gort.StrValString(propNode.ToStringTree()))))))
					}
				}
			}
		case int32(105) /* Token.CLASS */ :
			classMembers = n.GetLastChild()
			{
				var child pkg_rhino.Node = classMembers.GetFirstChild()
				_ = child
				for ; child != 0; child = child.GetNext() {
					if (child.IsMemberFunctionDef() || child.IsSetterDef()) || child.IsGetterDef() {
						this_.priv_ProcessProperties_processMemberDef_575(t, child)
					}
				}
			}
		case int32(32) /* Token.NAME */ :
			name = n.GetString()
			var_ = j2gort.CastPtr[Var](t.GetScope().GetVar(name))
			if var_ != nil {
				if var_.IsGlobal() {
					if (var_.GetInitialValue() != 0) && var_.GetInitialValue().IsFunction() {
						if t.InGlobalHoistScope() {
							if !(this_.priv_ProcessProperties_processGlobalFunctionDeclaration_573(t, n, var_)) {
								this_.priv_ProcessProperties_addGlobalUseOfSymbol_577(name, t.GetChunk(), StaticGet_AnalyzePrototypeProperties_VAR())
							}
						} else {
							this_.priv_ProcessProperties_addSymbolUse_568(name, t.GetChunk(), StaticGet_AnalyzePrototypeProperties_VAR())
						}
					}
				} else if j2gort.CastPtr[Scope](var_.GetScope()) != t.GetScope() {
					{
						cur_7 := this_.symbolStack.Iter()
						for cur_7.HasNext() {

							context_2 := cur_7.Next()
							_ = context_2
							if context_2.scope == j2gort.CastPtr[Scope](var_.GetScope()) {
								break
							}
							{
								r_9 := context_2.name
								r_9.readClosureVariables = true
							}
						}
					}
				}
			}
		default:
		}
	}
	if !(this_.outer_.rootScopeUsesAreGlobal) && (this_.priv_ProcessProperties_processNonFunctionPrototypeAssign_569(n, parent) != j2gort.NullString) {
		this_.symbolStack.Pop()
	}
}

// Java method: ProcessProperties.addSymbolUse(java.lang.String name, com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.jscomp.AnalyzePrototypeProperties.SymbolType type)
func (this_ *AnalyzePrototypeProperties_ProcessProperties) priv_ProcessProperties_addSymbolUse_568(name string, chunk *JSChunk, type_ AnalyzePrototypeProperties_SymbolType) {
	_ = name
	_ = chunk
	_ = type_
	var info *AnalyzePrototypeProperties_NameInfo = this_.outer_.priv_AnalyzePrototypeProperties_getNameInfoForName_565(name, type_)
	_ = info
	var def *AnalyzePrototypeProperties_NameInfo = nil
	_ = def
	{
		cur_10 := this_.symbolStack.Iter()
		for cur_10.HasNext() {

			context := cur_10.Next()
			_ = context
			def = context.name
			if def != this_.outer_.anonymousNode {
				break
			}
		}
	}
	if !(j2gort.CheckNotNullPtr(def) == info) {
		this_.outer_.symbolGraph.Connect_3_any_any_any(def, chunk, info)
	}
}

// Java method: ProcessProperties.processNonFunctionPrototypeAssign(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *AnalyzePrototypeProperties_ProcessProperties) priv_ProcessProperties_processNonFunctionPrototypeAssign_569(n pkg_rhino.Node, parent pkg_rhino.Node) string {
	_ = n
	_ = parent
	if this_.priv_ProcessProperties_isAssignRValue_571(n, parent) && !(n.IsFunction()) {
		return this_.priv_ProcessProperties_getPrototypePropertyNameFromRValue_572(n)
	}
	return j2gort.NullString
}

// Java method: ProcessProperties.isGlobalFunctionDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *AnalyzePrototypeProperties_ProcessProperties) priv_ProcessProperties_isGlobalFunctionDeclaration_570(t *NodeTraversal, n pkg_rhino.Node) bool {
	_ = t
	_ = n
	if !(t.InGlobalHoistScope() || ((n.IsFunction() && (t.GetScopeRoot() == n)) && j2gort.CastPtr[Scope](j2gort.CastPtr[Scope](t.GetScope().GetParent()).GetClosestHoistScope()).IsGlobal())) {
		return false
	}
	return (NodeUtil_IsFunctionDeclaration(n) || (n.IsFunction() && n.GetParent().IsName()))
}

// Java method: ProcessProperties.isAssignRValue(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *AnalyzePrototypeProperties_ProcessProperties) priv_ProcessProperties_isAssignRValue_571(n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = n
	_ = parent
	return (((parent != 0) && parent.IsAssign()) && (parent.GetFirstChild() != n))
}

// Java method: ProcessProperties.getPrototypePropertyNameFromRValue(com.google.javascript.rhino.Node rValue)
func (this_ *AnalyzePrototypeProperties_ProcessProperties) priv_ProcessProperties_getPrototypePropertyNameFromRValue_572(rValue pkg_rhino.Node) string {
	_ = rValue
	var lValue pkg_rhino.Node = NodeUtil_GetBestLValue(rValue)
	_ = lValue
	if (lValue == 0) || !((NodeUtil_MayBeObjectLitKey(lValue) && !(lValue.IsQuotedStringKey())) || NodeUtil_IsExprAssign(lValue.GetGrandparent())) {
		return j2gort.NullString
	}
	var lValueName string = NodeUtil_GetBestLValueName(lValue)
	_ = lValueName
	if lValueName == j2gort.NullString {
		return j2gort.NullString
	}
	var lastDot int32 = j2gort.JStrLastIndexOfCp(lValueName, int32(uint16(46)))
	_ = lastDot
	if lastDot == int32(-1) {
		return j2gort.NullString
	}
	var firstPart string = j2gort.JStrSub(lValueName, int32(0), lastDot)
	_ = firstPart
	if !(javalang.StringEndsWith(j2gort.CheckNotNullString(firstPart), ".prototype")) {
		return j2gort.NullString
	}
	return j2gort.JStrSubFrom(lValueName, (lastDot + int32(1)))
}

// Java method: ProcessProperties.processGlobalFunctionDeclaration(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node nameNode, com.google.javascript.jscomp.Var v)
func (this_ *AnalyzePrototypeProperties_ProcessProperties) priv_ProcessProperties_processGlobalFunctionDeclaration_573(t *NodeTraversal, nameNode pkg_rhino.Node, v *Var) bool {
	_ = t
	_ = nameNode
	_ = v
	var firstChild pkg_rhino.Node = nameNode.GetFirstChild()
	_ = firstChild
	var parent pkg_rhino.Node = nameNode.GetParent()
	_ = parent
	if this_.priv_ProcessProperties_isGlobalFunctionDeclaration_570(t, parent) || ((firstChild != 0) && this_.priv_ProcessProperties_isGlobalFunctionDeclaration_570(t, firstChild)) {
		var name string = nameNode.GetString()
		_ = name
		this_.outer_.priv_AnalyzePrototypeProperties_getNameInfoForName_565(name, StaticGet_AnalyzePrototypeProperties_VAR()).GetDeclarations().Add_1(NewAnalyzePrototypeProperties_GlobalFunction(nameNode, v, t.GetChunk()))
		if this_.outer_.compiler.GetCodingConvention().IsExported_2(name, false) || this_.outer_.anchorUnusedVars {
			this_.priv_ProcessProperties_addGlobalUseOfSymbol_577(name, t.GetChunk(), StaticGet_AnalyzePrototypeProperties_VAR())
		}
		return true
	}
	return false
}

// Java method: ProcessProperties.handlePossibleAssignmentToPrototype(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node ref)
func (this_ *AnalyzePrototypeProperties_ProcessProperties) priv_ProcessProperties_handlePossibleAssignmentToPrototype_574(t *NodeTraversal, ref pkg_rhino.Node) bool {
	_ = t
	_ = ref
	var root pkg_rhino.Node = NodeUtil_GetRootOfQualifiedName_1_Node(ref)
	_ = root
	var n pkg_rhino.Node = ref.GetParent()
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(26) /* Token.GETPROP */ :
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		var grandParent pkg_rhino.Node = parent.GetParent()
		_ = grandParent
		if NodeUtil_IsExprAssign(grandParent) && NodeUtil_IsNameDeclOrSimpleAssignLhs(n, parent) {
			var prop AnalyzePrototypeProperties_PrototypeProperty = NewAnalyzePrototypeProperties_AssignmentPrototypeProperty(grandParent, this_.priv_ProcessProperties_maybeGetVar_576(t, root), t.GetChunk())
			_ = prop
			this_.outer_.priv_AnalyzePrototypeProperties_getNameInfoForName_565(n.GetString(), StaticGet_AnalyzePrototypeProperties_PROPERTY()).GetDeclarations().Add_1(prop)
			return true
		}
	case int32(51) /* Token.ASSIGN */ :
		var map_ pkg_rhino.Node = n.GetSecondChild()
		_ = map_
		if map_.IsObjectLit() {
			{
				var key pkg_rhino.Node = map_.GetFirstChild()
				_ = key
				for ; key != 0; key = key.GetNext() {
					if !(key.IsQuotedStringKey()) && !(key.IsComputedProp()) {
						var name string = key.GetString()
						_ = name
						var prop_2 AnalyzePrototypeProperties_PrototypeProperty = NewAnalyzePrototypeProperties_LiteralPrototypeProperty(key.GetFirstChild(), n, this_.priv_ProcessProperties_maybeGetVar_576(t, root), t.GetChunk())
						_ = prop_2
						this_.outer_.priv_AnalyzePrototypeProperties_getNameInfoForName_565(name, StaticGet_AnalyzePrototypeProperties_PROPERTY()).GetDeclarations().Add_1(prop_2)
					}
				}
			}
			return true
		}
	default:
	}
	return false
}

// Java method: ProcessProperties.processMemberDef(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *AnalyzePrototypeProperties_ProcessProperties) priv_ProcessProperties_processMemberDef_575(t *NodeTraversal, n pkg_rhino.Node) {
	_ = t
	_ = n
	if !((n.IsMemberFunctionDef() || n.IsGetterDef()) || n.IsSetterDef()) {
		j2gort.CheckState(false)
	}
	var name string = n.GetString()
	_ = name
	if NodeUtil_IsEs6ConstructorMemberFunctionDef(n) || n.IsStaticMember() {
		return
	}
	var classNameNode pkg_rhino.Node = NodeUtil_GetNameNode(n.GetGrandparent())
	_ = classNameNode
	var var_ *Var
	_ = var_
	if (classNameNode != 0) && classNameNode.IsName() {
		var_ = j2gort.CastPtr[Var](t.GetScope().GetVar(classNameNode.GetString()))
	} else {
		var_ = nil
	}
	this_.outer_.priv_AnalyzePrototypeProperties_getNameInfoForName_565(name, StaticGet_AnalyzePrototypeProperties_PROPERTY()).GetDeclarations().Add_1(NewAnalyzePrototypeProperties_ClassMemberFunction(n, var_, t.GetChunk()))
}

// Java method: ProcessProperties.maybeGetVar(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node maybeName)
func (this_ *AnalyzePrototypeProperties_ProcessProperties) priv_ProcessProperties_maybeGetVar_576(t *NodeTraversal, maybeName pkg_rhino.Node) *Var {
	_ = t
	_ = maybeName
	if maybeName.IsName() {
		return j2gort.CastPtr[Var](t.GetScope().GetVar(maybeName.GetString()))
	} else {
		return nil
	}
}

// Java method: ProcessProperties.addGlobalUseOfSymbol(java.lang.String name, com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.jscomp.AnalyzePrototypeProperties.SymbolType type)
func (this_ *AnalyzePrototypeProperties_ProcessProperties) priv_ProcessProperties_addGlobalUseOfSymbol_577(name string, chunk *JSChunk, type_ AnalyzePrototypeProperties_SymbolType) {
	_ = name
	_ = chunk
	_ = type_
	this_.outer_.symbolGraph.Connect_3_any_any_any(this_.outer_.globalNode, chunk, this_.outer_.priv_AnalyzePrototypeProperties_getNameInfoForName_565(name, type_))
}

// Java class: com.google.javascript.jscomp.AnalyzePrototypeProperties.ProcessExternProperties
type AnalyzePrototypeProperties_ProcessExternProperties struct {
	NodeTraversal_AbstractPostOrderCallbackBase
	outer_ *AnalyzePrototypeProperties
}

func (this_ *AnalyzePrototypeProperties_ProcessExternProperties) BindSelf_(s any) {
	_ = s
}
func (this_ *AnalyzePrototypeProperties_ProcessExternProperties) IsJavaNil_() bool {
	return this_ == nil
}
func (*AnalyzePrototypeProperties_ProcessExternProperties) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AnalyzePrototypeProperties$ProcessExternProperties", "ProcessExternProperties"
}
func (this_ *AnalyzePrototypeProperties_ProcessExternProperties) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPostOrderCallbackBase.IdentityHash_()
}
func (*AnalyzePrototypeProperties_ProcessExternProperties) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AnalyzePrototypeProperties.ProcessExternProperties", "com.google.javascript.jscomp.AnalyzePrototypeProperties$ProcessExternProperties", func(v_ any) bool {
		return j2gort.InstanceOf[*AnalyzePrototypeProperties_ProcessExternProperties](v_)
	})
}
func (this_ *AnalyzePrototypeProperties_ProcessExternProperties) SetOuter_(o *AnalyzePrototypeProperties) {
	this_.outer_ = o
}
func (this_ *AnalyzePrototypeProperties_ProcessExternProperties) GetOuter_() *AnalyzePrototypeProperties {
	return this_.outer_
}
func (*AnalyzePrototypeProperties_ProcessExternProperties) Implements_jscomp_NodeTraversal_AbstractPostOrderCallback() {
}
func (*AnalyzePrototypeProperties_ProcessExternProperties) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *AnalyzePrototypeProperties_ProcessExternProperties) initFields_AnalyzePrototypeProperties_ProcessExternProperties() {
	_ = this_
}

// Java constructor: ProcessExternProperties()
func NewAnalyzePrototypeProperties_ProcessExternProperties(outer_ *AnalyzePrototypeProperties) *AnalyzePrototypeProperties_ProcessExternProperties {
	this_ := &AnalyzePrototypeProperties_ProcessExternProperties{}
	this_.outer_ = outer_
	InitAnalyzePrototypeProperties_ProcessExternProperties(this_)
	return this_
}

func InitAnalyzePrototypeProperties_ProcessExternProperties(this_ *AnalyzePrototypeProperties_ProcessExternProperties) {
	InitNodeTraversal_AbstractPostOrderCallback(&this_.NodeTraversal_AbstractPostOrderCallbackBase)
	this_.initFields_AnalyzePrototypeProperties_ProcessExternProperties()
}

// Java method: ProcessExternProperties.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *AnalyzePrototypeProperties_ProcessExternProperties) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if n.IsGetProp() {
		this_.outer_.symbolGraph.Connect_3_any_any_any(this_.outer_.externNode, this_.outer_.firstChunk, this_.outer_.priv_AnalyzePrototypeProperties_getNameInfoForName_565(n.GetString(), StaticGet_AnalyzePrototypeProperties_PROPERTY()))
	} else if (n.IsMemberFunctionDef() || n.IsGetterDef()) || n.IsSetterDef() {
		this_.outer_.symbolGraph.Connect_3_any_any_any(this_.outer_.externNode, this_.outer_.firstChunk, this_.outer_.priv_AnalyzePrototypeProperties_getNameInfoForName_565(n.GetString(), StaticGet_AnalyzePrototypeProperties_PROPERTY()))
	}
}

// Java class: com.google.javascript.jscomp.AnalyzePrototypeProperties.PropagateReferences
type AnalyzePrototypeProperties_PropagateReferences struct {
	outer_  *AnalyzePrototypeProperties
	idHash_ atomic.Int32
}

func (this_ *AnalyzePrototypeProperties_PropagateReferences) BindSelf_(s any) {
	_ = s
}
func (this_ *AnalyzePrototypeProperties_PropagateReferences) IsJavaNil_() bool { return this_ == nil }
func (*AnalyzePrototypeProperties_PropagateReferences) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AnalyzePrototypeProperties$PropagateReferences", "PropagateReferences"
}
func (this_ *AnalyzePrototypeProperties_PropagateReferences) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AnalyzePrototypeProperties_PropagateReferences) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AnalyzePrototypeProperties.PropagateReferences", "com.google.javascript.jscomp.AnalyzePrototypeProperties$PropagateReferences", func(v_ any) bool {
		return j2gort.InstanceOf[*AnalyzePrototypeProperties_PropagateReferences](v_)
	})
}
func (this_ *AnalyzePrototypeProperties_PropagateReferences) SetOuter_(o *AnalyzePrototypeProperties) {
	this_.outer_ = o
}
func (this_ *AnalyzePrototypeProperties_PropagateReferences) GetOuter_() *AnalyzePrototypeProperties {
	return this_.outer_
}
func (*AnalyzePrototypeProperties_PropagateReferences) Implements_graph_FixedPointGraphTraversal_EdgeCallback() {
}

func (this_ *AnalyzePrototypeProperties_PropagateReferences) initFields_AnalyzePrototypeProperties_PropagateReferences() {
	_ = this_
}

// Java constructor: PropagateReferences()
func NewAnalyzePrototypeProperties_PropagateReferences(outer_ *AnalyzePrototypeProperties) *AnalyzePrototypeProperties_PropagateReferences {
	this_ := &AnalyzePrototypeProperties_PropagateReferences{}
	this_.outer_ = outer_
	InitAnalyzePrototypeProperties_PropagateReferences(this_)
	return this_
}

func InitAnalyzePrototypeProperties_PropagateReferences(this_ *AnalyzePrototypeProperties_PropagateReferences) {
	this_.initFields_AnalyzePrototypeProperties_PropagateReferences()
}

// Java method: PropagateReferences.traverseEdge(com.google.javascript.jscomp.AnalyzePrototypeProperties.NameInfo start, com.google.javascript.jscomp.JSChunk edge, com.google.javascript.jscomp.AnalyzePrototypeProperties.NameInfo dest)
func (this_ *AnalyzePrototypeProperties_PropagateReferences) TraverseEdge(start *AnalyzePrototypeProperties_NameInfo, edge *JSChunk, dest *AnalyzePrototypeProperties_NameInfo) bool {
	_ = start
	_ = edge
	_ = dest
	if start.IsReferenced() {
		var startChunk *JSChunk = start.GetDeepestCommonChunkRef()
		_ = startChunk
		if (startChunk != nil) && this_.outer_.chunkGraph.DependsOn(startChunk, edge) {
			return dest.MarkReference(startChunk)
		} else {
			return dest.MarkReference(edge)
		}
	}
	return false
}

// Java interface: com.google.javascript.jscomp.AnalyzePrototypeProperties.Symbol
type AnalyzePrototypeProperties_Symbol interface {
	Implements_jscomp_AnalyzePrototypeProperties_Symbol()
	IsJavaNil_() bool
	GetRootVar() *Var
	GetChunk() *JSChunk
}

func ToIface_AnalyzePrototypeProperties_Symbol[T interface {
	AnalyzePrototypeProperties_Symbol
	comparable
}](v T) AnalyzePrototypeProperties_Symbol {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_AnalyzePrototypeProperties_Symbol[T AnalyzePrototypeProperties_Symbol](v T) AnalyzePrototypeProperties_Symbol {
	var vi AnalyzePrototypeProperties_Symbol = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AnalyzePrototypeProperties.Symbol", "com.google.javascript.jscomp.AnalyzePrototypeProperties$Symbol", func(v_ any) bool {
		return j2gort.InstanceOf[AnalyzePrototypeProperties_Symbol](v_)
	})
}

// Java enum: com.google.javascript.jscomp.AnalyzePrototypeProperties.SymbolType
type AnalyzePrototypeProperties_SymbolType uint8

const (
	enum_AnalyzePrototypeProperties_SymbolType_PROPERTY AnalyzePrototypeProperties_SymbolType = 1
	enum_AnalyzePrototypeProperties_SymbolType_VAR      AnalyzePrototypeProperties_SymbolType = 2
)

var enumBaseTable_AnalyzePrototypeProperties_SymbolType = [...]j2gort.EnumBase{
	{Name_: "PROPERTY", Ordinal_: 0},
	{Name_: "VAR", Ordinal_: 1},
}
var enumValues_AnalyzePrototypeProperties_SymbolType = []AnalyzePrototypeProperties_SymbolType{1, 2}

func StaticGet_AnalyzePrototypeProperties_SymbolType_PROPERTY() AnalyzePrototypeProperties_SymbolType {
	return enum_AnalyzePrototypeProperties_SymbolType_PROPERTY
}
func StaticGet_AnalyzePrototypeProperties_SymbolType_VAR() AnalyzePrototypeProperties_SymbolType {
	return enum_AnalyzePrototypeProperties_SymbolType_VAR
}

func AnalyzePrototypeProperties_SymbolType_Values() []AnalyzePrototypeProperties_SymbolType {
	return j2gort.CloneSlice(enumValues_AnalyzePrototypeProperties_SymbolType)
}

func AnalyzePrototypeProperties_SymbolType_ValueOf(name string) AnalyzePrototypeProperties_SymbolType {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_AnalyzePrototypeProperties_SymbolType {
		if enumBaseTable_AnalyzePrototypeProperties_SymbolType[i_].Name_ == name {
			return AnalyzePrototypeProperties_SymbolType(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.AnalyzePrototypeProperties.SymbolType." + name))
}

func (this_ AnalyzePrototypeProperties_SymbolType) IsJavaNil_() bool                              { return this_ == 0 }
func (AnalyzePrototypeProperties_SymbolType) HasIdentityEquals_()                                 {}
func (AnalyzePrototypeProperties_SymbolType) Implements_j2gort_Ext_java_lang_Enum()               {}
func (AnalyzePrototypeProperties_SymbolType) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ AnalyzePrototypeProperties_SymbolType) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (AnalyzePrototypeProperties_SymbolType) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AnalyzePrototypeProperties$SymbolType", "SymbolType"
}
func (this_ AnalyzePrototypeProperties_SymbolType) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_AnalyzePrototypeProperties_SymbolType[this_-1]
}
func (this_ AnalyzePrototypeProperties_SymbolType) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ AnalyzePrototypeProperties_SymbolType) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ AnalyzePrototypeProperties_SymbolType) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_AnalyzePrototypeProperties_SymbolType[this_-1].Name_
}
func (this_ AnalyzePrototypeProperties_SymbolType) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_AnalyzePrototypeProperties_SymbolType[this_-1].Name_
}
func (this_ AnalyzePrototypeProperties_SymbolType) Equals(other_ any) bool {
	o_, ok_ := other_.(AnalyzePrototypeProperties_SymbolType)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ AnalyzePrototypeProperties_SymbolType) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_AnalyzePrototypeProperties_SymbolType[this_-1].HashCode()
}
func (this_ AnalyzePrototypeProperties_SymbolType) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_AnalyzePrototypeProperties_SymbolType[this_-1].CompareTo(other_)
}
func (this_ AnalyzePrototypeProperties_SymbolType) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_AnalyzePrototypeProperties_SymbolType[this_-1].Name_)
}
func init() {
	vals_ := []any{AnalyzePrototypeProperties_SymbolType(1), AnalyzePrototypeProperties_SymbolType(2)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AnalyzePrototypeProperties.SymbolType", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.AnalyzePrototypeProperties$SymbolType", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AnalyzePrototypeProperties.SymbolType", "com.google.javascript.jscomp.AnalyzePrototypeProperties$SymbolType", func(v_ any) bool {
		return j2gort.InstanceOf[AnalyzePrototypeProperties_SymbolType](v_)
	})
}

// Java class: com.google.javascript.jscomp.AnalyzePrototypeProperties.GlobalFunction
type AnalyzePrototypeProperties_GlobalFunction struct {
	var_    *Var     `gson:"var"`
	chunk   *JSChunk `gson:"chunk"`
	idHash_ atomic.Int32
}

func (this_ *AnalyzePrototypeProperties_GlobalFunction) BindSelf_(s any) {
	_ = s
}
func (this_ *AnalyzePrototypeProperties_GlobalFunction) IsJavaNil_() bool { return this_ == nil }
func (*AnalyzePrototypeProperties_GlobalFunction) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AnalyzePrototypeProperties$GlobalFunction", "GlobalFunction"
}
func (this_ *AnalyzePrototypeProperties_GlobalFunction) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AnalyzePrototypeProperties_GlobalFunction) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AnalyzePrototypeProperties.GlobalFunction", "com.google.javascript.jscomp.AnalyzePrototypeProperties$GlobalFunction", func(v_ any) bool {
		return j2gort.InstanceOf[*AnalyzePrototypeProperties_GlobalFunction](v_)
	})
}
func (*AnalyzePrototypeProperties_GlobalFunction) Implements_jscomp_AnalyzePrototypeProperties_Symbol() {
}

func (this_ *AnalyzePrototypeProperties_GlobalFunction) initFields_AnalyzePrototypeProperties_GlobalFunction() {
	_ = this_
}

// Java constructor: GlobalFunction(com.google.javascript.rhino.Node nameNode, com.google.javascript.jscomp.Var var, com.google.javascript.jscomp.JSChunk chunk)
func NewAnalyzePrototypeProperties_GlobalFunction(nameNode pkg_rhino.Node, var_ *Var, chunk *JSChunk) *AnalyzePrototypeProperties_GlobalFunction {
	this_ := &AnalyzePrototypeProperties_GlobalFunction{}
	InitAnalyzePrototypeProperties_GlobalFunction(this_, nameNode, var_, chunk)
	return this_
}

func InitAnalyzePrototypeProperties_GlobalFunction(this_ *AnalyzePrototypeProperties_GlobalFunction, nameNode pkg_rhino.Node, var_ *Var, chunk *JSChunk) {
	_ = nameNode
	_ = var_
	_ = chunk
	this_.initFields_AnalyzePrototypeProperties_GlobalFunction()
	var parent pkg_rhino.Node = nameNode.GetParent()
	_ = parent
	if !((NodeUtil_IsNameDeclaration(parent) && var_.IsGlobal()) || NodeUtil_IsFunctionDeclaration(parent)) {
		j2gort.CheckState(false)
	}
	this_.var_ = var_
	this_.chunk = chunk
}

// Java method: GlobalFunction.getRootVar()
func (this_ *AnalyzePrototypeProperties_GlobalFunction) GetRootVar() *Var {
	return this_.var_
}

// Java method: GlobalFunction.getChunk()
func (this_ *AnalyzePrototypeProperties_GlobalFunction) GetChunk() *JSChunk {
	return this_.chunk
}

// Java interface: com.google.javascript.jscomp.AnalyzePrototypeProperties.Property
type AnalyzePrototypeProperties_Property interface {
	Implements_jscomp_AnalyzePrototypeProperties_Property()
	IsJavaNil_() bool
	AnalyzePrototypeProperties_Symbol
}

func ToIface_AnalyzePrototypeProperties_Property[T interface {
	AnalyzePrototypeProperties_Property
	comparable
}](v T) AnalyzePrototypeProperties_Property {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_AnalyzePrototypeProperties_Property[T AnalyzePrototypeProperties_Property](v T) AnalyzePrototypeProperties_Property {
	var vi AnalyzePrototypeProperties_Property = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AnalyzePrototypeProperties.Property", "com.google.javascript.jscomp.AnalyzePrototypeProperties$Property", func(v_ any) bool {
		return j2gort.InstanceOf[AnalyzePrototypeProperties_Property](v_)
	})
}

// Java class: com.google.javascript.jscomp.AnalyzePrototypeProperties.ClassMemberFunction
type AnalyzePrototypeProperties_ClassMemberFunction struct {
	var_    *Var           `gson:"var"`
	chunk   *JSChunk       `gson:"chunk"`
	node    pkg_rhino.Node `gson:"node"`
	idHash_ atomic.Int32
}

func (this_ *AnalyzePrototypeProperties_ClassMemberFunction) BindSelf_(s any) {
	_ = s
}
func (this_ *AnalyzePrototypeProperties_ClassMemberFunction) IsJavaNil_() bool { return this_ == nil }
func (*AnalyzePrototypeProperties_ClassMemberFunction) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AnalyzePrototypeProperties$ClassMemberFunction", "ClassMemberFunction"
}
func (this_ *AnalyzePrototypeProperties_ClassMemberFunction) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AnalyzePrototypeProperties_ClassMemberFunction) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AnalyzePrototypeProperties.ClassMemberFunction", "com.google.javascript.jscomp.AnalyzePrototypeProperties$ClassMemberFunction", func(v_ any) bool {
		return j2gort.InstanceOf[*AnalyzePrototypeProperties_ClassMemberFunction](v_)
	})
}
func (*AnalyzePrototypeProperties_ClassMemberFunction) Implements_jscomp_AnalyzePrototypeProperties_Property() {
}
func (*AnalyzePrototypeProperties_ClassMemberFunction) Implements_jscomp_AnalyzePrototypeProperties_Symbol() {
}

func (this_ *AnalyzePrototypeProperties_ClassMemberFunction) initFields_AnalyzePrototypeProperties_ClassMemberFunction() {
	_ = this_
}

// Java constructor: ClassMemberFunction(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.Var var, com.google.javascript.jscomp.JSChunk chunk)
func NewAnalyzePrototypeProperties_ClassMemberFunction(node pkg_rhino.Node, var_ *Var, chunk *JSChunk) *AnalyzePrototypeProperties_ClassMemberFunction {
	this_ := &AnalyzePrototypeProperties_ClassMemberFunction{}
	InitAnalyzePrototypeProperties_ClassMemberFunction(this_, node, var_, chunk)
	return this_
}

func InitAnalyzePrototypeProperties_ClassMemberFunction(this_ *AnalyzePrototypeProperties_ClassMemberFunction, node pkg_rhino.Node, var_ *Var, chunk *JSChunk) {
	_ = node
	_ = var_
	_ = chunk
	this_.initFields_AnalyzePrototypeProperties_ClassMemberFunction()
	if !(node.GetParent().IsClassMembers()) {
		j2gort.CheckState(false)
	}
	if !((node.IsMemberFunctionDef() || node.IsSetterDef()) || node.IsGetterDef()) {
		j2gort.CheckState(false)
	}
	this_.node = node
	this_.var_ = var_
	this_.chunk = chunk
}

// Java method: ClassMemberFunction.getRootVar()
func (this_ *AnalyzePrototypeProperties_ClassMemberFunction) GetRootVar() *Var {
	return this_.var_
}

// Java method: ClassMemberFunction.getChunk()
func (this_ *AnalyzePrototypeProperties_ClassMemberFunction) GetChunk() *JSChunk {
	return this_.chunk
}

// Java method: ClassMemberFunction.getFunctionNode()
func (this_ *AnalyzePrototypeProperties_ClassMemberFunction) GetFunctionNode() pkg_rhino.Node {
	return this_.node.GetOnlyChild()
}

// Java method: ClassMemberFunction.getDefinitionNode()
func (this_ *AnalyzePrototypeProperties_ClassMemberFunction) GetDefinitionNode() pkg_rhino.Node {
	return this_.node
}

// Java interface: com.google.javascript.jscomp.AnalyzePrototypeProperties.PrototypeProperty
type AnalyzePrototypeProperties_PrototypeProperty interface {
	Implements_jscomp_AnalyzePrototypeProperties_PrototypeProperty()
	IsJavaNil_() bool
	AnalyzePrototypeProperties_Property
	GetPrototype() pkg_rhino.Node
	GetValue() pkg_rhino.Node
}

func ToIface_AnalyzePrototypeProperties_PrototypeProperty[T interface {
	AnalyzePrototypeProperties_PrototypeProperty
	comparable
}](v T) AnalyzePrototypeProperties_PrototypeProperty {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_AnalyzePrototypeProperties_PrototypeProperty[T AnalyzePrototypeProperties_PrototypeProperty](v T) AnalyzePrototypeProperties_PrototypeProperty {
	var vi AnalyzePrototypeProperties_PrototypeProperty = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AnalyzePrototypeProperties.PrototypeProperty", "com.google.javascript.jscomp.AnalyzePrototypeProperties$PrototypeProperty", func(v_ any) bool {
		return j2gort.InstanceOf[AnalyzePrototypeProperties_PrototypeProperty](v_)
	})
}

// Java class: com.google.javascript.jscomp.AnalyzePrototypeProperties.AssignmentPrototypeProperty
type AnalyzePrototypeProperties_AssignmentPrototypeProperty struct {
	rootVar  *Var           `gson:"rootVar"`
	chunk    *JSChunk       `gson:"chunk"`
	exprNode pkg_rhino.Node `gson:"exprNode"`
	idHash_  atomic.Int32
}

func (this_ *AnalyzePrototypeProperties_AssignmentPrototypeProperty) BindSelf_(s any) {
	_ = s
}
func (this_ *AnalyzePrototypeProperties_AssignmentPrototypeProperty) IsJavaNil_() bool {
	return this_ == nil
}
func (*AnalyzePrototypeProperties_AssignmentPrototypeProperty) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AnalyzePrototypeProperties$AssignmentPrototypeProperty", "AssignmentPrototypeProperty"
}
func (this_ *AnalyzePrototypeProperties_AssignmentPrototypeProperty) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AnalyzePrototypeProperties_AssignmentPrototypeProperty) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AnalyzePrototypeProperties.AssignmentPrototypeProperty", "com.google.javascript.jscomp.AnalyzePrototypeProperties$AssignmentPrototypeProperty", func(v_ any) bool {
		return j2gort.InstanceOf[*AnalyzePrototypeProperties_AssignmentPrototypeProperty](v_)
	})
}
func (*AnalyzePrototypeProperties_AssignmentPrototypeProperty) Implements_jscomp_AnalyzePrototypeProperties_PrototypeProperty() {
}
func (*AnalyzePrototypeProperties_AssignmentPrototypeProperty) Implements_jscomp_AnalyzePrototypeProperties_Property() {
}
func (*AnalyzePrototypeProperties_AssignmentPrototypeProperty) Implements_jscomp_AnalyzePrototypeProperties_Symbol() {
}

func (this_ *AnalyzePrototypeProperties_AssignmentPrototypeProperty) initFields_AnalyzePrototypeProperties_AssignmentPrototypeProperty() {
	_ = this_
}

// Java constructor: AssignmentPrototypeProperty(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.Var rootVar, com.google.javascript.jscomp.JSChunk chunk)
func NewAnalyzePrototypeProperties_AssignmentPrototypeProperty(node pkg_rhino.Node, rootVar *Var, chunk *JSChunk) *AnalyzePrototypeProperties_AssignmentPrototypeProperty {
	this_ := &AnalyzePrototypeProperties_AssignmentPrototypeProperty{}
	InitAnalyzePrototypeProperties_AssignmentPrototypeProperty(this_, node, rootVar, chunk)
	return this_
}

func InitAnalyzePrototypeProperties_AssignmentPrototypeProperty(this_ *AnalyzePrototypeProperties_AssignmentPrototypeProperty, node pkg_rhino.Node, rootVar *Var, chunk *JSChunk) {
	_ = node
	_ = rootVar
	_ = chunk
	this_.initFields_AnalyzePrototypeProperties_AssignmentPrototypeProperty()
	this_.exprNode = node
	this_.rootVar = rootVar
	this_.chunk = chunk
}

// Java method: AssignmentPrototypeProperty.getRootVar()
func (this_ *AnalyzePrototypeProperties_AssignmentPrototypeProperty) GetRootVar() *Var {
	return this_.rootVar
}

// Java method: AssignmentPrototypeProperty.getPrototype()
func (this_ *AnalyzePrototypeProperties_AssignmentPrototypeProperty) GetPrototype() pkg_rhino.Node {
	return this_.priv_AssignmentPrototypeProperty_getAssignNode_585().GetFirstFirstChild()
}

// Java method: AssignmentPrototypeProperty.getValue()
func (this_ *AnalyzePrototypeProperties_AssignmentPrototypeProperty) GetValue() pkg_rhino.Node {
	return this_.priv_AssignmentPrototypeProperty_getAssignNode_585().GetLastChild()
}

// Java method: AssignmentPrototypeProperty.getAssignNode()
func (this_ *AnalyzePrototypeProperties_AssignmentPrototypeProperty) priv_AssignmentPrototypeProperty_getAssignNode_585() pkg_rhino.Node {
	return this_.exprNode.GetFirstChild()
}

// Java method: AssignmentPrototypeProperty.getChunk()
func (this_ *AnalyzePrototypeProperties_AssignmentPrototypeProperty) GetChunk() *JSChunk {
	return this_.chunk
}

// Java class: com.google.javascript.jscomp.AnalyzePrototypeProperties.LiteralPrototypeProperty
type AnalyzePrototypeProperties_LiteralPrototypeProperty struct {
	rootVar *Var           `gson:"rootVar"`
	chunk   *JSChunk       `gson:"chunk"`
	value   pkg_rhino.Node `gson:"value"`
	assign  pkg_rhino.Node `gson:"assign"`
	idHash_ atomic.Int32
}

func (this_ *AnalyzePrototypeProperties_LiteralPrototypeProperty) BindSelf_(s any) {
	_ = s
}
func (this_ *AnalyzePrototypeProperties_LiteralPrototypeProperty) IsJavaNil_() bool {
	return this_ == nil
}
func (*AnalyzePrototypeProperties_LiteralPrototypeProperty) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AnalyzePrototypeProperties$LiteralPrototypeProperty", "LiteralPrototypeProperty"
}
func (this_ *AnalyzePrototypeProperties_LiteralPrototypeProperty) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AnalyzePrototypeProperties_LiteralPrototypeProperty) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AnalyzePrototypeProperties.LiteralPrototypeProperty", "com.google.javascript.jscomp.AnalyzePrototypeProperties$LiteralPrototypeProperty", func(v_ any) bool {
		return j2gort.InstanceOf[*AnalyzePrototypeProperties_LiteralPrototypeProperty](v_)
	})
}
func (*AnalyzePrototypeProperties_LiteralPrototypeProperty) Implements_jscomp_AnalyzePrototypeProperties_PrototypeProperty() {
}
func (*AnalyzePrototypeProperties_LiteralPrototypeProperty) Implements_jscomp_AnalyzePrototypeProperties_Property() {
}
func (*AnalyzePrototypeProperties_LiteralPrototypeProperty) Implements_jscomp_AnalyzePrototypeProperties_Symbol() {
}

func (this_ *AnalyzePrototypeProperties_LiteralPrototypeProperty) initFields_AnalyzePrototypeProperties_LiteralPrototypeProperty() {
	_ = this_
}

// Java constructor: LiteralPrototypeProperty(com.google.javascript.rhino.Node value, com.google.javascript.rhino.Node assign, com.google.javascript.jscomp.Var rootVar, com.google.javascript.jscomp.JSChunk chunk)
func NewAnalyzePrototypeProperties_LiteralPrototypeProperty(value pkg_rhino.Node, assign pkg_rhino.Node, rootVar *Var, chunk *JSChunk) *AnalyzePrototypeProperties_LiteralPrototypeProperty {
	this_ := &AnalyzePrototypeProperties_LiteralPrototypeProperty{}
	InitAnalyzePrototypeProperties_LiteralPrototypeProperty(this_, value, assign, rootVar, chunk)
	return this_
}

func InitAnalyzePrototypeProperties_LiteralPrototypeProperty(this_ *AnalyzePrototypeProperties_LiteralPrototypeProperty, value pkg_rhino.Node, assign pkg_rhino.Node, rootVar *Var, chunk *JSChunk) {
	_ = value
	_ = assign
	_ = rootVar
	_ = chunk
	this_.initFields_AnalyzePrototypeProperties_LiteralPrototypeProperty()
	this_.value = value
	this_.assign = assign
	this_.rootVar = rootVar
	this_.chunk = chunk
}

// Java method: LiteralPrototypeProperty.getRootVar()
func (this_ *AnalyzePrototypeProperties_LiteralPrototypeProperty) GetRootVar() *Var {
	return this_.rootVar
}

// Java method: LiteralPrototypeProperty.getPrototype()
func (this_ *AnalyzePrototypeProperties_LiteralPrototypeProperty) GetPrototype() pkg_rhino.Node {
	return this_.assign.GetFirstChild()
}

// Java method: LiteralPrototypeProperty.getValue()
func (this_ *AnalyzePrototypeProperties_LiteralPrototypeProperty) GetValue() pkg_rhino.Node {
	return this_.value
}

// Java method: LiteralPrototypeProperty.getChunk()
func (this_ *AnalyzePrototypeProperties_LiteralPrototypeProperty) GetChunk() *JSChunk {
	return this_.chunk
}

// Java class: com.google.javascript.jscomp.AnalyzePrototypeProperties.NameContext
type AnalyzePrototypeProperties_NameContext struct {
	name    *AnalyzePrototypeProperties_NameInfo `gson:"name"`
	scope   *Scope                               `gson:"scope"`
	idHash_ atomic.Int32
}

func (this_ *AnalyzePrototypeProperties_NameContext) BindSelf_(s any) {
	_ = s
}
func (this_ *AnalyzePrototypeProperties_NameContext) IsJavaNil_() bool { return this_ == nil }
func (*AnalyzePrototypeProperties_NameContext) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AnalyzePrototypeProperties$NameContext", "NameContext"
}
func (this_ *AnalyzePrototypeProperties_NameContext) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AnalyzePrototypeProperties_NameContext) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AnalyzePrototypeProperties.NameContext", "com.google.javascript.jscomp.AnalyzePrototypeProperties$NameContext", func(v_ any) bool {
		return j2gort.InstanceOf[*AnalyzePrototypeProperties_NameContext](v_)
	})
}

func (this_ *AnalyzePrototypeProperties_NameContext) initFields_AnalyzePrototypeProperties_NameContext() {
	_ = this_
}

// Java constructor: NameContext(com.google.javascript.jscomp.AnalyzePrototypeProperties.NameInfo name, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Scope scope)
func NewAnalyzePrototypeProperties_NameContext(name *AnalyzePrototypeProperties_NameInfo, scope *Scope) *AnalyzePrototypeProperties_NameContext {
	this_ := &AnalyzePrototypeProperties_NameContext{}
	InitAnalyzePrototypeProperties_NameContext(this_, name, scope)
	return this_
}

func InitAnalyzePrototypeProperties_NameContext(this_ *AnalyzePrototypeProperties_NameContext, name *AnalyzePrototypeProperties_NameInfo, scope *Scope) {
	_ = name
	_ = scope
	this_.initFields_AnalyzePrototypeProperties_NameContext()
	this_.name = name
	this_.scope = scope
}

// Java class: com.google.javascript.jscomp.AnalyzePrototypeProperties.NameInfo
type AnalyzePrototypeProperties_NameInfo struct {
	outer_                *AnalyzePrototypeProperties
	name                  string                     `gson:"name"`
	declarations          j2gort.Ext_java_util_Deque `gson:"declarations"`
	deepestCommonChunkRef *JSChunk                   `gson:"deepestCommonChunkRef"`
	idHash_               atomic.Int32
	referenced            bool `gson:"referenced"`
	readClosureVariables  bool `gson:"readClosureVariables"`
	referencesSuper       bool `gson:"referencesSuper"`
}

func (this_ *AnalyzePrototypeProperties_NameInfo) BindSelf_(s any) {
	_ = s
}
func (this_ *AnalyzePrototypeProperties_NameInfo) IsJavaNil_() bool { return this_ == nil }
func (*AnalyzePrototypeProperties_NameInfo) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.AnalyzePrototypeProperties$NameInfo", "NameInfo"
}
func (this_ *AnalyzePrototypeProperties_NameInfo) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*AnalyzePrototypeProperties_NameInfo) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.AnalyzePrototypeProperties.NameInfo", "com.google.javascript.jscomp.AnalyzePrototypeProperties$NameInfo", func(v_ any) bool {
		return j2gort.InstanceOf[*AnalyzePrototypeProperties_NameInfo](v_)
	})
}
func (this_ *AnalyzePrototypeProperties_NameInfo) SetOuter_(o *AnalyzePrototypeProperties) {
	this_.outer_ = o
}
func (this_ *AnalyzePrototypeProperties_NameInfo) GetOuter_() *AnalyzePrototypeProperties {
	return this_.outer_
}

func (this_ *AnalyzePrototypeProperties_NameInfo) initFields_AnalyzePrototypeProperties_NameInfo() {
	_ = this_
	this_.name = j2gort.NullString

	this_.declarations = j2gort.NewExt_java_util_ArrayDeque()

}

// Java constructor: NameInfo(java.lang.String name)
func NewAnalyzePrototypeProperties_NameInfo(outer_ *AnalyzePrototypeProperties, name string) *AnalyzePrototypeProperties_NameInfo {
	this_ := &AnalyzePrototypeProperties_NameInfo{}
	this_.outer_ = outer_
	InitAnalyzePrototypeProperties_NameInfo(this_, name)
	return this_
}

func InitAnalyzePrototypeProperties_NameInfo(this_ *AnalyzePrototypeProperties_NameInfo, name string) {
	_ = name
	this_.initFields_AnalyzePrototypeProperties_NameInfo()
	this_.name = name
}

// Java method: NameInfo.toString()
func (this_ *AnalyzePrototypeProperties_NameInfo) String() string {
	return this_.name
}

// Java method: NameInfo.isReferenced()
func (this_ *AnalyzePrototypeProperties_NameInfo) IsReferenced() bool {
	return this_.referenced
}

// Java method: NameInfo.readsClosureVariables()
func (this_ *AnalyzePrototypeProperties_NameInfo) ReadsClosureVariables() bool {
	return this_.readClosureVariables
}

// Java method: NameInfo.referencesSuper()
func (this_ *AnalyzePrototypeProperties_NameInfo) ReferencesSuper() bool {
	return this_.referencesSuper
}

// Java method: NameInfo.markReference(com.google.javascript.jscomp.@org.jspecify.annotations.Nullable JSChunk chunk)
func (this_ *AnalyzePrototypeProperties_NameInfo) MarkReference(chunk *JSChunk) bool {
	_ = chunk
	var hasChanged bool = false
	_ = hasChanged
	if !(this_.referenced) {
		this_.referenced = true
		hasChanged = true
	}
	var originalDeepestCommon *JSChunk = this_.deepestCommonChunkRef
	_ = originalDeepestCommon
	if this_.deepestCommonChunkRef == nil {
		this_.deepestCommonChunkRef = chunk
	} else {
		this_.deepestCommonChunkRef = this_.outer_.chunkGraph.GetDeepestCommonDependencyInclusive_2(this_.deepestCommonChunkRef, chunk)
	}
	if originalDeepestCommon != this_.deepestCommonChunkRef {
		hasChanged = true
	}
	return hasChanged
}

// Java method: NameInfo.getDeepestCommonChunkRef()
func (this_ *AnalyzePrototypeProperties_NameInfo) GetDeepestCommonChunkRef() *JSChunk {
	return this_.deepestCommonChunkRef
}

// Java method: NameInfo.getDeclarations()
func (this_ *AnalyzePrototypeProperties_NameInfo) GetDeclarations() j2gort.Ext_java_util_Deque {
	return this_.declarations
}
func (*AnalyzePrototypeProperties_ProcessExternProperties) AlwaysShouldTraverse_()      {}
func (*AnalyzePrototypeProperties_ProcessExternProperties) CanReleaseSyntacticScopes_() {}
