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
// Port of com/google/javascript/jscomp/GlobalNamespace.java (generated mechanically by j2go).

package jscomp

import (
	"fmt"
	"sync/atomic"

	pkg_diagnostic "github.com/paullepoulpe/jsgo/transpiled/gen/diagnostic"
	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

var lambdaSingleton_global_namespace_1 j2gort.Ext_java_util_function_Predicate

func init() {
	lambdaSingleton_global_namespace_1 = j2gort.NewNodeLambda_Ext_java_util_function_Predicate(func(rawNode_1 uint32) bool {
		n := pkg_rhino.Node(rawNode_1)
		_ = n
		return true
	})
}

// Java class: com.google.javascript.jscomp.GlobalNamespace
type GlobalNamespace struct {
	compiler                      AbstractCompiler                                                                          `gson:"compiler"`
	spreadSiblingCache            j2gort.OrderedMap[pkg_rhino.Node, *j2gort.BoxBoolean]                                     `gson:"spreadSiblingCache"`
	decisionsLog                  pkg_diagnostic.LogFile                                                                    `gson:"decisionsLog"`
	globalNames                   []*GlobalNamespace_Name                                                                   `gson:"globalNames"`
	nameMap                       j2gort.OrderedMap[string, *GlobalNamespace_Name]                                          `gson:"nameMap"`
	nameMapByModule               j2gort.CustomOrderedMap[ModuleMetadataMap_ModuleMetadata, *GlobalNamespace_ModuleNameMap] `gson:"nameMapByModule"`
	shouldTraverseScript          j2gort.Ext_java_util_function_Predicate                                                   `gson:"shouldTraverseScript"`
	root                          pkg_rhino.Node                                                                            `gson:"root"`
	externsRoot                   pkg_rhino.Node                                                                            `gson:"externsRoot"`
	globalRoot                    pkg_rhino.Node                                                                            `gson:"globalRoot"`
	idHash_                       atomic.Int32
	enableImplicitlyAliasedValues bool                       `gson:"enableImplicitlyAliasedValues"`
	sourceKind                    GlobalNamespace_SourceKind `gson:"sourceKind"`
	generated                     bool                       `gson:"generated"`
	scopeCreator_                 SyntacticScopeCreator
	slab_GlobalNamespace_Ref      j2gort.PooledSlab[GlobalNamespace_Ref]
}

func (this_ *GlobalNamespace) BindSelf_(s any) {
	_ = s
}
func (this_ *GlobalNamespace) IsJavaNil_() bool { return this_ == nil }
func (*GlobalNamespace) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GlobalNamespace", "GlobalNamespace"
}
func (this_ *GlobalNamespace) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GlobalNamespace) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GlobalNamespace", "com.google.javascript.jscomp.GlobalNamespace", func(v_ any) bool {
		return j2gort.InstanceOf[*GlobalNamespace](v_)
	})
}
func (*GlobalNamespace) Implements_rhino_StaticScope()       {}
func (*GlobalNamespace) Implements_rhino_StaticSymbolTable() {}

var static_GlobalNamespace_GOOG_PROVIDE pkg_rhino.QualifiedName

func StaticGet_GlobalNamespace_GOOG_PROVIDE() pkg_rhino.QualifiedName {
	clinit_GlobalNamespace()
	return static_GlobalNamespace_GOOG_PROVIDE
}
func StaticSet_GlobalNamespace_GOOG_PROVIDE(v pkg_rhino.QualifiedName) {
	clinit_GlobalNamespace()
	static_GlobalNamespace_GOOG_PROVIDE = v
}

var clinitDone_GlobalNamespace uint32
var clinitState_GlobalNamespace j2gort.ClinitState

func clinit_GlobalNamespace() {
	if atomic.LoadUint32(&clinitDone_GlobalNamespace) != 0 {
		return
	}
	clinitSlow_GlobalNamespace()
}

func clinitSlow_GlobalNamespace() {
	if !j2gort.EnterClinit(&clinitState_GlobalNamespace, &clinitDone_GlobalNamespace) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_GlobalNamespace, &clinitDone_GlobalNamespace)
	static_GlobalNamespace_GOOG_PROVIDE = pkg_rhino.QualifiedName_Of("goog.provide")
}

func (this_ *GlobalNamespace) initFields_GlobalNamespace() {
	_ = this_
	this_.globalRoot = pkg_rhino.IR_Root([]pkg_rhino.Node{})
	this_.spreadSiblingCache = j2gort.NewOrderedMap[pkg_rhino.Node, *j2gort.BoxBoolean]()

	this_.globalNames = make([]*GlobalNamespace_Name, 0, int32(1024))
	this_.nameMap = j2gort.AcquireOrderedMapCap[string, *GlobalNamespace_Name](int32(1024))
	this_.nameMapByModule = j2gort.AcquireCustomOrderedMap[ModuleMetadataMap_ModuleMetadata, *GlobalNamespace_ModuleNameMap]()
	this_.shouldTraverseScript = lambdaSingleton_global_namespace_1
}

// Java constructor: GlobalNamespace(com.google.javascript.jscomp.diagnostic.LogFile decisionsLog, com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node root)
func NewGlobalNamespace_3_LogFile_AbstractCompiler_Node(decisionsLog pkg_diagnostic.LogFile, compiler AbstractCompiler, root pkg_rhino.Node) *GlobalNamespace {
	this_ := &GlobalNamespace{}
	InitGlobalNamespace_3_LogFile_AbstractCompiler_Node(this_, decisionsLog, compiler, root)
	return this_
}

func InitGlobalNamespace_3_LogFile_AbstractCompiler_Node(this_ *GlobalNamespace, decisionsLog pkg_diagnostic.LogFile, compiler AbstractCompiler, root pkg_rhino.Node) {
	_ = decisionsLog
	_ = compiler
	_ = root
	InitGlobalNamespace_4(this_, decisionsLog, compiler, pkg_rhino.Node(0), root)
}

// Java constructor: GlobalNamespace(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node root)
func NewGlobalNamespace_2(compiler AbstractCompiler, root pkg_rhino.Node) *GlobalNamespace {
	this_ := &GlobalNamespace{}
	InitGlobalNamespace_2(this_, compiler, root)
	return this_
}

func InitGlobalNamespace_2(this_ *GlobalNamespace, compiler AbstractCompiler, root pkg_rhino.Node) {
	_ = compiler
	_ = root
	InitGlobalNamespace_4(this_, nil, compiler, pkg_rhino.Node(0), root)
}

// Java constructor: GlobalNamespace(com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.Node externsRoot, com.google.javascript.rhino.Node root)
func NewGlobalNamespace_3_AbstractCompiler_Node_Node(compiler AbstractCompiler, externsRoot pkg_rhino.Node, root pkg_rhino.Node) *GlobalNamespace {
	this_ := &GlobalNamespace{}
	InitGlobalNamespace_3_AbstractCompiler_Node_Node(this_, compiler, externsRoot, root)
	return this_
}

func InitGlobalNamespace_3_AbstractCompiler_Node_Node(this_ *GlobalNamespace, compiler AbstractCompiler, externsRoot pkg_rhino.Node, root pkg_rhino.Node) {
	_ = compiler
	_ = externsRoot
	_ = root
	InitGlobalNamespace_4(this_, nil, compiler, externsRoot, root)
}

// Java constructor: GlobalNamespace(com.google.javascript.jscomp.diagnostic.@org.jspecify.annotations.Nullable LogFile decisionsLog, com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node externsRoot, com.google.javascript.rhino.Node root)
func NewGlobalNamespace_4(decisionsLog pkg_diagnostic.LogFile, compiler AbstractCompiler, externsRoot pkg_rhino.Node, root pkg_rhino.Node) *GlobalNamespace {
	this_ := &GlobalNamespace{}
	InitGlobalNamespace_4(this_, decisionsLog, compiler, externsRoot, root)
	return this_
}

func InitGlobalNamespace_4(this_ *GlobalNamespace, decisionsLog pkg_diagnostic.LogFile, compiler AbstractCompiler, externsRoot pkg_rhino.Node, root pkg_rhino.Node) {
	_ = decisionsLog
	_ = compiler
	_ = externsRoot
	_ = root
	this_.initFields_GlobalNamespace()
	this_.decisionsLog = decisionsLog
	this_.compiler = compiler
	this_.externsRoot = externsRoot
	this_.root = root
	this_.enableImplicitlyAliasedValues = !(compiler.GetOptions().GetAssumeStaticInheritanceIsNotUsed())
}

// Java method: GlobalNamespace.getModuleName(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata, java.lang.String name)
func (this_ *GlobalNamespace) priv_GlobalNamespace_getModuleName_3047(metadata ModuleMetadataMap_ModuleMetadata, name string) *GlobalNamespace_Name {
	_ = metadata
	_ = name
	var map_ *GlobalNamespace_ModuleNameMap = this_.nameMapByModule.Get(metadata)
	_ = map_
	if map_ == nil {
		return nil
	} else {
		return map_.names.Get(name)
	}
}

// Java method: GlobalNamespace.putModuleName(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata, java.lang.String name, com.google.javascript.jscomp.GlobalNamespace.Name node)
func (this_ *GlobalNamespace) priv_GlobalNamespace_putModuleName_3048(metadata ModuleMetadataMap_ModuleMetadata, name string, node *GlobalNamespace_Name) {
	_ = metadata
	_ = name
	_ = node
	var map_ *GlobalNamespace_ModuleNameMap = this_.nameMapByModule.Get(metadata)
	_ = map_
	if map_ == nil {
		map_ = NewGlobalNamespace_ModuleNameMap()
		this_.nameMapByModule.Put(metadata, map_)
	}
	map_.names.Put(name, node)
}

// Java method: GlobalNamespace.setShouldTraverseScriptPredicate(java.util.function.Predicate<com.google.javascript.rhino.Node> shouldTraverseScript)
func (this_ *GlobalNamespace) SetShouldTraverseScriptPredicate(shouldTraverseScript j2gort.Ext_java_util_function_Predicate) {
	_ = shouldTraverseScript
	this_.shouldTraverseScript = shouldTraverseScript
}

// Java method: GlobalNamespace.hasExternsRoot()
func (this_ *GlobalNamespace) HasExternsRoot() bool {
	return (this_.externsRoot != 0)
}

// Java method: GlobalNamespace.getRootNode()
func (this_ *GlobalNamespace) GetRootNode_0() pkg_rhino.Node {
	return this_.root.GetParent()
}

// Java method: GlobalNamespace.getRootNode(java.lang.String name, com.google.javascript.jscomp.Scope s)
func (this_ *GlobalNamespace) priv_GlobalNamespace_getRootNode_3051(name string, s *Scope) pkg_rhino.Node {
	_ = name
	_ = s
	return this_.priv_GlobalNamespace_getRootNodeForTopVar_3052(GlobalNamespace_GetTopVarName(name), s)
}

// Java method: GlobalNamespace.getRootNodeForTopVar(java.lang.String topVarName, com.google.javascript.jscomp.Scope s)
func (this_ *GlobalNamespace) priv_GlobalNamespace_getRootNodeForTopVar_3052(topVarName string, s *Scope) pkg_rhino.Node {
	_ = topVarName
	_ = s
	var v *Var = s.LookupVar(topVarName)
	_ = v
	if v == nil {
		var providedName *GlobalNamespace_Name = this_.nameMap.Get(topVarName)
		_ = providedName
		if (providedName != nil) && providedName.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_IS_PROVIDED()) {
			return this_.globalRoot
		} else {
			return pkg_rhino.Node(0)
		}
	}
	if v.IsLocal() {
		return v.GetScopeRoot()
	} else {
		return this_.globalRoot
	}
}

// Java method: GlobalNamespace.getParentScope()
func (this_ *GlobalNamespace) GetParentScope() any {
	return nil
}

// Java method: GlobalNamespace.getSlot(java.lang.String name)
func (this_ *GlobalNamespace) GetSlot(name string) any {
	_ = name
	return j2gort.PtrToAny(j2gort.CastPtr[GlobalNamespace_Name](this_.GetOwnSlot(name)))
}

// Java method: GlobalNamespace.getOwnSlot(java.lang.String name)
func (this_ *GlobalNamespace) GetOwnSlot(name string) any {
	_ = name
	this_.EnsureGenerated()
	return j2gort.PtrToAny(this_.nameMap.Get(name))
}

// Java method: GlobalNamespace.getReferences(com.google.javascript.jscomp.GlobalNamespace.Name slot)
func (this_ *GlobalNamespace) GetReferences(arg_2 any) any {
	slot := j2gort.CastPtr[GlobalNamespace_Name](arg_2)
	_ = slot
	this_.EnsureGenerated()
	return j2gort.Ext_java_util_Collections_UnmodifiableCollection(slot.GetRefs())
}

// Java method: GlobalNamespace.getScope(com.google.javascript.jscomp.GlobalNamespace.Name slot)
func (this_ *GlobalNamespace) GetScope(arg_3 any) any {
	slot := j2gort.CastPtr[GlobalNamespace_Name](arg_3)
	_ = slot
	return this_
}

// Java method: GlobalNamespace.getAllSymbols()
func (this_ *GlobalNamespace) GetAllSymbols() any {
	this_.EnsureGenerated()
	return j2gort.Ext_java_util_Collections_UnmodifiableCollection(j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.nameMap.LiveValues()))
}

// Java method: GlobalNamespace.ensureGenerated()
func (this_ *GlobalNamespace) EnsureGenerated() {
	if !(this_.generated) {
		this_.priv_GlobalNamespace_process_3060()
	}
}

// Java method: GlobalNamespace.containsName(java.lang.String name)
func (this_ *GlobalNamespace) ContainsName(name string) bool {
	_ = name
	this_.EnsureGenerated()
	return this_.nameMap.ContainsKey(name)
}

// Java method: GlobalNamespace.scanNewNodes(java.util.Set<com.google.javascript.jscomp.GlobalNamespace.AstChange> newNodes)
func (this_ *GlobalNamespace) ScanNewNodes(newNodes j2gort.Ext_java_util_Set) {
	_ = newNodes
	var builder *GlobalNamespace_BuildGlobalNamespace = NewGlobalNamespace_BuildGlobalNamespace(this_)
	_ = builder
	{
		cur_4 := j2gort.NewIterCursor(j2gort.CheckNotNull(newNodes))
		for cur_4.HasNext() {
			rawElem_5 := cur_4.Next()
			info := j2gort.AsIface[GlobalNamespace_AstChange](rawElem_5)
			_ = info
			if !(info.Node().IsQualifiedName()) && !(NodeUtil_MayBeObjectLitKey(info.Node())) {
				continue
			}
			this_.priv_GlobalNamespace_scanFromNode_3059(builder, info.Scope(), info.Node(), info.Chunk())
		}
	}
}

// Java method: GlobalNamespace.scanFromNode(com.google.javascript.jscomp.GlobalNamespace.BuildGlobalNamespace builder, com.google.javascript.jscomp.Scope scope, com.google.javascript.rhino.Node n, com.google.javascript.jscomp.JSChunk chunk)
func (this_ *GlobalNamespace) priv_GlobalNamespace_scanFromNode_3059(builder *GlobalNamespace_BuildGlobalNamespace, scope *Scope, n pkg_rhino.Node, chunk *JSChunk) {
	_ = builder
	_ = scope
	_ = n
	_ = chunk
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	if (n.IsName() || n.IsGetProp()) && parent.IsGetProp() {
		this_.priv_GlobalNamespace_scanFromNode_3059(builder, scope, n.GetParent(), chunk)
	} else if (n.GetPrevious() != 0) && n.GetPrevious().IsObjectPattern() {
		var pattern pkg_rhino.Node = n.GetPrevious()
		_ = pattern
		{
			var key pkg_rhino.Node = pattern.GetFirstChild()
			_ = key
			for ; key != 0; key = key.GetNext() {
				if key.IsStringKey() {
					this_.priv_GlobalNamespace_scanFromNode_3059(builder, scope, key, chunk)
				}
			}
		}
	}
	builder.priv_BuildGlobalNamespace_collect_3070(scope, chunk, n)
}

// Java method: GlobalNamespace.process()
func (this_ *GlobalNamespace) priv_GlobalNamespace_process_3060() {
	if this_.scopeCreator_.reusableScanner == nil {
		InitSyntacticScopeCreator_1(&this_.scopeCreator_, this_.compiler)
	}

	var traversal *NodeTraversal_Builder = NodeTraversal_S_Builder().SetCompiler(this_.compiler).SetCallback_1_NodeTraversal_Callback(NewGlobalNamespace_BuildGlobalNamespace(this_)).SetScopeCreator(&this_.scopeCreator_)
	_ = traversal
	if this_.HasExternsRoot() {
		traversal.TraverseRoots(this_.externsRoot, this_.root)
	} else {
		traversal.Traverse(this_.root)
	}
	this_.generated = true
}

// Java static method: GlobalNamespace.getTopVarName(java.lang.String name)
func GlobalNamespace_GetTopVarName(name string) string {
	_ = name
	var firstDotIndex int32 = j2gort.JStrIndexOfCp(name, int32(uint16(46)))
	_ = firstDotIndex
	if firstDotIndex == int32(-1) {
		return name
	} else {
		return j2gort.JStrSub(name, int32(0), firstDotIndex)
	}
}

// Java method: GlobalNamespace.getNameFromModule(com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata moduleMetadata, java.lang.String name)
func (this_ *GlobalNamespace) GetNameFromModule(moduleMetadata ModuleMetadataMap_ModuleMetadata, name string) *GlobalNamespace_Name {
	_ = moduleMetadata
	_ = name
	j2gort.CheckNotNullIface(moduleMetadata)
	j2gort.CheckNotNullString(name)
	this_.EnsureGenerated()
	return this_.priv_GlobalNamespace_getModuleName_3047(moduleMetadata, name)
}

// Java method: GlobalNamespace.declarationHasFollowingObjectSpreadSibling(com.google.javascript.rhino.Node declaration)
func (this_ *GlobalNamespace) priv_GlobalNamespace_declarationHasFollowingObjectSpreadSibling_3062(declaration pkg_rhino.Node) bool {
	_ = declaration
	if !(declaration.GetParent().IsObjectLit()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(declaration))
	}
	var cached *j2gort.BoxBoolean = this_.spreadSiblingCache.Get(declaration)
	_ = cached
	if cached != nil {
		return j2gort.UnboxBoolean(cached)
	}
	var toCache bool = false
	_ = toCache
	{
		var sibling pkg_rhino.Node = declaration.GetParent().GetLastChild()
		_ = sibling
		for ; sibling != 0; sibling = sibling.GetPrevious() {
			if sibling.IsSpread() {
				toCache = true
			}
			this_.spreadSiblingCache.Put(sibling, j2gort.ValueOfBoolean(toCache))
		}
	}
	return j2gort.UnboxBoolean(this_.spreadSiblingCache.Get(declaration))
}

// Java method: GlobalNamespace.createNameForTesting(java.lang.String name)
func (this_ *GlobalNamespace) CreateNameForTesting(name string) *GlobalNamespace_Name {
	_ = name
	return NewGlobalNamespace_Name(this_, name, nil, StaticGet_GlobalNamespace_SourceKind_CODE())
}

// Java method: GlobalNamespace.getOutermostObjectLit(com.google.javascript.rhino.Node objLitKey)
func (this_ *GlobalNamespace) priv_GlobalNamespace_getOutermostObjectLit_3064(objLitKey pkg_rhino.Node) pkg_rhino.Node {
	_ = objLitKey
	var outermostObjectLit pkg_rhino.Node = objLitKey.GetParent()
	_ = outermostObjectLit
	if !(outermostObjectLit.IsObjectLit()) {
		j2gort.CheckStateMsg(false, j2gort.BoxNode(outermostObjectLit))
	}
	for {
		var objLitGrandparent pkg_rhino.Node = outermostObjectLit.GetGrandparent()
		_ = objLitGrandparent
		if (objLitGrandparent != 0) && objLitGrandparent.IsObjectLit() {
			outermostObjectLit = objLitGrandparent
		} else {
			return outermostObjectLit
		}
	}
	panic("unreachable")
}

// Java static method: GlobalNamespace.isQnameDeclarationWithoutAssignment(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node node)
func GlobalNamespace_IsQnameDeclarationWithoutAssignment(node pkg_rhino.Node) bool {
	_ = node
	return (((node != 0) && node.IsGetProp()) && node.GetParent().IsExprResult())
}

func (this_ *GlobalNamespace) GetTopmostScopeOfEventualDeclaration(p0_ string) pkg_rhino.StaticScope {
	return pkg_rhino.IfaceDefault_StaticScope_GetTopmostScopeOfEventualDeclaration(this_, p0_)
}

// Java enum: com.google.javascript.jscomp.GlobalNamespace.SourceKind
type GlobalNamespace_SourceKind uint8

const (
	enum_GlobalNamespace_SourceKind_EXTERN       GlobalNamespace_SourceKind = 1
	enum_GlobalNamespace_SourceKind_TYPE_SUMMARY GlobalNamespace_SourceKind = 2
	enum_GlobalNamespace_SourceKind_CODE         GlobalNamespace_SourceKind = 3
)

var enumBaseTable_GlobalNamespace_SourceKind = [...]j2gort.EnumBase{
	{Name_: "EXTERN", Ordinal_: 0},
	{Name_: "TYPE_SUMMARY", Ordinal_: 1},
	{Name_: "CODE", Ordinal_: 2},
}
var enumValues_GlobalNamespace_SourceKind = []GlobalNamespace_SourceKind{1, 2, 3}

func StaticGet_GlobalNamespace_SourceKind_EXTERN() GlobalNamespace_SourceKind {
	return enum_GlobalNamespace_SourceKind_EXTERN
}
func StaticGet_GlobalNamespace_SourceKind_TYPE_SUMMARY() GlobalNamespace_SourceKind {
	return enum_GlobalNamespace_SourceKind_TYPE_SUMMARY
}
func StaticGet_GlobalNamespace_SourceKind_CODE() GlobalNamespace_SourceKind {
	return enum_GlobalNamespace_SourceKind_CODE
}

func GlobalNamespace_SourceKind_Values() []GlobalNamespace_SourceKind {
	return j2gort.CloneSlice(enumValues_GlobalNamespace_SourceKind)
}

func GlobalNamespace_SourceKind_ValueOf(name string) GlobalNamespace_SourceKind {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_GlobalNamespace_SourceKind {
		if enumBaseTable_GlobalNamespace_SourceKind[i_].Name_ == name {
			return GlobalNamespace_SourceKind(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.GlobalNamespace.SourceKind." + name))
}

// Java static method: SourceKind.fromScriptNode(com.google.javascript.rhino.Node n)
func GlobalNamespace_SourceKind_FromScriptNode(n pkg_rhino.Node) GlobalNamespace_SourceKind {
	_ = n
	if !(n.IsFromExterns()) {
		return StaticGet_GlobalNamespace_SourceKind_CODE()
	} else if NodeUtil_IsFromTypeSummary(n) {
		return StaticGet_GlobalNamespace_SourceKind_TYPE_SUMMARY()
	} else {
		return StaticGet_GlobalNamespace_SourceKind_EXTERN()
	}
}

func (this_ GlobalNamespace_SourceKind) IsJavaNil_() bool                              { return this_ == 0 }
func (GlobalNamespace_SourceKind) HasIdentityEquals_()                                 {}
func (GlobalNamespace_SourceKind) Implements_j2gort_Ext_java_lang_Enum()               {}
func (GlobalNamespace_SourceKind) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ GlobalNamespace_SourceKind) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (GlobalNamespace_SourceKind) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GlobalNamespace$SourceKind", "SourceKind"
}
func (this_ GlobalNamespace_SourceKind) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_GlobalNamespace_SourceKind[this_-1]
}
func (this_ GlobalNamespace_SourceKind) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ GlobalNamespace_SourceKind) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ GlobalNamespace_SourceKind) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_GlobalNamespace_SourceKind[this_-1].Name_
}
func (this_ GlobalNamespace_SourceKind) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_GlobalNamespace_SourceKind[this_-1].Name_
}
func (this_ GlobalNamespace_SourceKind) Equals(other_ any) bool {
	o_, ok_ := other_.(GlobalNamespace_SourceKind)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ GlobalNamespace_SourceKind) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_GlobalNamespace_SourceKind[this_-1].HashCode()
}
func (this_ GlobalNamespace_SourceKind) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_GlobalNamespace_SourceKind[this_-1].CompareTo(other_)
}
func (this_ GlobalNamespace_SourceKind) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_GlobalNamespace_SourceKind[this_-1].Name_)
}
func init() {
	vals_ := []any{GlobalNamespace_SourceKind(1), GlobalNamespace_SourceKind(2), GlobalNamespace_SourceKind(3)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.GlobalNamespace.SourceKind", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.GlobalNamespace$SourceKind", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GlobalNamespace.SourceKind", "com.google.javascript.jscomp.GlobalNamespace$SourceKind", func(v_ any) bool {
		return j2gort.InstanceOf[GlobalNamespace_SourceKind](v_)
	})
}

// Java class: com.google.javascript.jscomp.GlobalNamespace.ModuleNameMap
type GlobalNamespace_ModuleNameMap struct {
	names   j2gort.OrderedMap[string, *GlobalNamespace_Name] `gson:"names"`
	idHash_ atomic.Int32
}

func (this_ *GlobalNamespace_ModuleNameMap) BindSelf_(s any) {
	_ = s
}
func (this_ *GlobalNamespace_ModuleNameMap) IsJavaNil_() bool { return this_ == nil }
func (*GlobalNamespace_ModuleNameMap) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GlobalNamespace$ModuleNameMap", "ModuleNameMap"
}
func (this_ *GlobalNamespace_ModuleNameMap) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GlobalNamespace_ModuleNameMap) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GlobalNamespace.ModuleNameMap", "com.google.javascript.jscomp.GlobalNamespace$ModuleNameMap", func(v_ any) bool {
		return j2gort.InstanceOf[*GlobalNamespace_ModuleNameMap](v_)
	})
}

func (this_ *GlobalNamespace_ModuleNameMap) initFields_GlobalNamespace_ModuleNameMap() {
	_ = this_
	this_.names = j2gort.AcquireOrderedMap[string, *GlobalNamespace_Name](

	// Java constructor: ModuleNameMap()
	)
}

func NewGlobalNamespace_ModuleNameMap() *GlobalNamespace_ModuleNameMap {
	this_ := &GlobalNamespace_ModuleNameMap{}
	InitGlobalNamespace_ModuleNameMap(this_)
	return this_
}

func InitGlobalNamespace_ModuleNameMap(this_ *GlobalNamespace_ModuleNameMap) {
	this_.initFields_GlobalNamespace_ModuleNameMap()
}

// Java interface: com.google.javascript.jscomp.GlobalNamespace.AstChange
type GlobalNamespace_AstChange interface {
	Implements_jscomp_GlobalNamespace_AstChange()
	IsJavaNil_() bool
	Node() pkg_rhino.Node
	Chunk() *JSChunk
	Scope() *Scope
}

func ToIface_GlobalNamespace_AstChange[T interface {
	GlobalNamespace_AstChange
	comparable
}](v T) GlobalNamespace_AstChange {
	var z T
	if v == z {
		return nil
	}
	return v
}

func Upcast_GlobalNamespace_AstChange[T GlobalNamespace_AstChange](v T) GlobalNamespace_AstChange {
	var vi GlobalNamespace_AstChange = v
	if vi == nil || vi.IsJavaNil_() {
		return nil
	}
	return vi
}

func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GlobalNamespace.AstChange", "com.google.javascript.jscomp.GlobalNamespace$AstChange", func(v_ any) bool {
		return j2gort.InstanceOf[GlobalNamespace_AstChange](v_)
	})
}

// Java record: com.google.javascript.jscomp.GlobalNamespace.SimpleAstChange
type GlobalNamespace_SimpleAstChange struct {
	chunk *JSChunk       `gson:"chunk"`
	scope *Scope         `gson:"scope"`
	node  pkg_rhino.Node `gson:"node"`
}

func (this_ *GlobalNamespace_SimpleAstChange) BindSelf_(s any) {
	_ = s
}
func (this_ *GlobalNamespace_SimpleAstChange) IsJavaNil_() bool { return this_ == nil }
func (*GlobalNamespace_SimpleAstChange) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GlobalNamespace$SimpleAstChange", "SimpleAstChange"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GlobalNamespace.SimpleAstChange", "com.google.javascript.jscomp.GlobalNamespace$SimpleAstChange", func(v_ any) bool {
		return j2gort.InstanceOf[*GlobalNamespace_SimpleAstChange](v_)
	})
}
func (*GlobalNamespace_SimpleAstChange) Implements_j2gort_Ext_java_lang_Record()      {}
func (*GlobalNamespace_SimpleAstChange) Implements_jscomp_GlobalNamespace_AstChange() {}

func (this_ *GlobalNamespace_SimpleAstChange) initFields_GlobalNamespace_SimpleAstChange() {
	_ = this_
}

// Java constructor: SimpleAstChange(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.jscomp.Scope scope)
func NewGlobalNamespace_SimpleAstChange(node pkg_rhino.Node, chunk *JSChunk, scope *Scope) *GlobalNamespace_SimpleAstChange {
	this_ := &GlobalNamespace_SimpleAstChange{}
	InitGlobalNamespace_SimpleAstChange(this_, node, chunk, scope)
	return this_
}

func InitGlobalNamespace_SimpleAstChange(this_ *GlobalNamespace_SimpleAstChange, node pkg_rhino.Node, chunk *JSChunk, scope *Scope) {
	_ = node
	_ = chunk
	_ = scope
	this_.initFields_GlobalNamespace_SimpleAstChange()
	this_.node = node
	this_.chunk = chunk
	this_.scope = scope
}

func (this_ *GlobalNamespace_SimpleAstChange) Node() pkg_rhino.Node { return this_.node }
func (this_ *GlobalNamespace_SimpleAstChange) Chunk() *JSChunk      { return this_.chunk }
func (this_ *GlobalNamespace_SimpleAstChange) Scope() *Scope        { return this_.scope }
func (this_ *GlobalNamespace_SimpleAstChange) Equals(other_ any) bool {
	o_, ok_ := other_.(*GlobalNamespace_SimpleAstChange)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.node != o_.node {
		return false
	}
	if this_.chunk != o_.chunk {
		return false
	}
	if this_.scope != o_.scope {
		return false
	}
	return true
}
func (this_ *GlobalNamespace_SimpleAstChange) HashCode() int32 {
	var h_ int32 = 0
	h_ = 31*h_ + int32(this_.node)
	{
		hp_6 := this_.chunk
		h_ = 31 * h_
		if hp_6 != nil {
			h_ += j2gort.IdentityHash(hp_6)
		}
	}
	{
		hp_7 := this_.scope
		h_ = 31 * h_
		if hp_7 != nil {
			h_ += j2gort.IdentityHash(hp_7)
		}
	}
	return h_
}
func (this_ *GlobalNamespace_SimpleAstChange) String() string {
	_ = fmt.Sprint
	return "SimpleAstChange[" + "node=" + j2gort.StrVal(j2gort.BoxNode(this_.node)) + ", " + "chunk=" + j2gort.StrVal(j2gort.PtrToAny(this_.chunk)) + ", " + "scope=" + j2gort.StrVal(j2gort.PtrToAny(this_.scope)) + "]"
}

// Java record: com.google.javascript.jscomp.GlobalNamespace.RefBasedAstChange
type GlobalNamespace_RefBasedAstChange struct {
	ref  *GlobalNamespace_Ref `gson:"ref"`
	node pkg_rhino.Node       `gson:"node"`
}

func (this_ *GlobalNamespace_RefBasedAstChange) BindSelf_(s any) {
	_ = s
}
func (this_ *GlobalNamespace_RefBasedAstChange) IsJavaNil_() bool { return this_ == nil }
func (*GlobalNamespace_RefBasedAstChange) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GlobalNamespace$RefBasedAstChange", "RefBasedAstChange"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GlobalNamespace.RefBasedAstChange", "com.google.javascript.jscomp.GlobalNamespace$RefBasedAstChange", func(v_ any) bool {
		return j2gort.InstanceOf[*GlobalNamespace_RefBasedAstChange](v_)
	})
}
func (*GlobalNamespace_RefBasedAstChange) Implements_j2gort_Ext_java_lang_Record()      {}
func (*GlobalNamespace_RefBasedAstChange) Implements_jscomp_GlobalNamespace_AstChange() {}

func (this_ *GlobalNamespace_RefBasedAstChange) initFields_GlobalNamespace_RefBasedAstChange() {
	_ = this_
}

// Java constructor: RefBasedAstChange(com.google.javascript.jscomp.GlobalNamespace.Ref ref, com.google.javascript.rhino.Node node)
func NewGlobalNamespace_RefBasedAstChange(ref *GlobalNamespace_Ref, node pkg_rhino.Node) *GlobalNamespace_RefBasedAstChange {
	this_ := &GlobalNamespace_RefBasedAstChange{}
	InitGlobalNamespace_RefBasedAstChange(this_, ref, node)
	return this_
}

func InitGlobalNamespace_RefBasedAstChange(this_ *GlobalNamespace_RefBasedAstChange, ref *GlobalNamespace_Ref, node pkg_rhino.Node) {
	_ = ref
	_ = node
	this_.initFields_GlobalNamespace_RefBasedAstChange()
	this_.ref = ref
	this_.node = node
}

// Java method: RefBasedAstChange.chunk()
func (this_ *GlobalNamespace_RefBasedAstChange) Chunk() *JSChunk {
	return this_.ref.GetChunk()
}

// Java method: RefBasedAstChange.scope()
func (this_ *GlobalNamespace_RefBasedAstChange) Scope() *Scope {
	return this_.ref.scope
}

func (this_ *GlobalNamespace_RefBasedAstChange) Ref() *GlobalNamespace_Ref { return this_.ref }
func (this_ *GlobalNamespace_RefBasedAstChange) Node() pkg_rhino.Node      { return this_.node }
func (this_ *GlobalNamespace_RefBasedAstChange) Equals(other_ any) bool {
	o_, ok_ := other_.(*GlobalNamespace_RefBasedAstChange)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.ref != o_.ref {
		return false
	}
	if this_.node != o_.node {
		return false
	}
	return true
}
func (this_ *GlobalNamespace_RefBasedAstChange) HashCode() int32 {
	var h_ int32 = 0
	{
		hp_6 := this_.ref
		h_ = 31 * h_
		if hp_6 != nil {
			h_ += j2gort.IdentityHash(hp_6)
		}
	}
	h_ = 31*h_ + int32(this_.node)
	return h_
}
func (this_ *GlobalNamespace_RefBasedAstChange) String() string {
	_ = fmt.Sprint
	return "RefBasedAstChange[" + "ref=" + j2gort.StrVal(j2gort.PtrToAny(this_.ref)) + ", " + "node=" + j2gort.StrVal(j2gort.BoxNode(this_.node)) + "]"
}

// Java class: com.google.javascript.jscomp.GlobalNamespace.BuildGlobalNamespace
type GlobalNamespace_BuildGlobalNamespace struct {
	NodeTraversal_AbstractPreOrderCallbackBase
	outer_                    *GlobalNamespace
	curMetadata               ModuleMetadataMap_ModuleMetadata `gson:"curMetadata"`
	curModuleRoot             pkg_rhino.Node                   `gson:"curModuleRoot"`
	scratchObjLitKeyAnalysis_ GlobalNamespace_ObjLitStringKeyAnalysis
}

func (this_ *GlobalNamespace_BuildGlobalNamespace) BindSelf_(s any) {
	_ = s
}
func (this_ *GlobalNamespace_BuildGlobalNamespace) IsJavaNil_() bool { return this_ == nil }
func (*GlobalNamespace_BuildGlobalNamespace) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GlobalNamespace$BuildGlobalNamespace", "BuildGlobalNamespace"
}
func (this_ *GlobalNamespace_BuildGlobalNamespace) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	return this_.NodeTraversal_AbstractPreOrderCallbackBase.IdentityHash_()
}
func (*GlobalNamespace_BuildGlobalNamespace) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GlobalNamespace.BuildGlobalNamespace", "com.google.javascript.jscomp.GlobalNamespace$BuildGlobalNamespace", func(v_ any) bool {
		return j2gort.InstanceOf[*GlobalNamespace_BuildGlobalNamespace](v_)
	})
}
func (this_ *GlobalNamespace_BuildGlobalNamespace) SetOuter_(o *GlobalNamespace) { this_.outer_ = o }
func (this_ *GlobalNamespace_BuildGlobalNamespace) GetOuter_() *GlobalNamespace  { return this_.outer_ }
func (*GlobalNamespace_BuildGlobalNamespace) Implements_jscomp_NodeTraversal_AbstractPreOrderCallback() {
}
func (*GlobalNamespace_BuildGlobalNamespace) Implements_jscomp_NodeTraversal_Callback() {}

func (this_ *GlobalNamespace_BuildGlobalNamespace) initFields_GlobalNamespace_BuildGlobalNamespace() {
	_ = this_

}

// Java constructor: BuildGlobalNamespace()
func NewGlobalNamespace_BuildGlobalNamespace(outer_ *GlobalNamespace) *GlobalNamespace_BuildGlobalNamespace {
	this_ := &GlobalNamespace_BuildGlobalNamespace{}
	this_.outer_ = outer_
	InitGlobalNamespace_BuildGlobalNamespace(this_)
	return this_
}

func InitGlobalNamespace_BuildGlobalNamespace(this_ *GlobalNamespace_BuildGlobalNamespace) {
	InitNodeTraversal_AbstractPreOrderCallback(&this_.NodeTraversal_AbstractPreOrderCallbackBase)
	this_.initFields_GlobalNamespace_BuildGlobalNamespace()
}

// Java method: BuildGlobalNamespace.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *GlobalNamespace_BuildGlobalNamespace) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if n.IsScript() && !(j2gort.PredTestNode(this_.outer_.shouldTraverseScript, uint32(n))) {
		return false
	}
	if this_.outer_.HasExternsRoot() && n.IsScript() {
		this_.outer_.sourceKind = GlobalNamespace_SourceKind_FromScriptNode(n)
	} else if n == this_.outer_.root {
		this_.outer_.sourceKind = StaticGet_GlobalNamespace_SourceKind_CODE()
	}
	if n.IsModuleBody() || NodeUtil_IsBundledGoogModuleScopeRoot(n) {
		this_.priv_BuildGlobalNamespace_setupModuleMetadata_3069(n)
	} else if n.IsScript() || NodeUtil_IsBundledGoogModuleCall(n) {
		this_.curModuleRoot = pkg_rhino.Node(0)
		this_.curMetadata = nil
	}
	this_.priv_BuildGlobalNamespace_collect_3070(t.GetScope(), t.GetChunk(), n)
	return true
}

// Java method: BuildGlobalNamespace.setupModuleMetadata(com.google.javascript.rhino.Node moduleRoot)
func (this_ *GlobalNamespace_BuildGlobalNamespace) priv_BuildGlobalNamespace_setupModuleMetadata_3069(moduleRoot pkg_rhino.Node) {
	_ = moduleRoot
	var moduleMap *ModuleMap = this_.outer_.compiler.GetModuleMap()
	_ = moduleMap
	if moduleMap == nil {
		return
	}
	this_.curModuleRoot = moduleRoot
	this_.curMetadata = j2gort.CheckNotNullIface(ModuleImportResolver_GetModuleFromScopeRoot(moduleMap, Upcast_CompilerInputProvider(this_.outer_.compiler), moduleRoot).Metadata())
	if this_.curMetadata.GetModuleMetadataMap_ModuleMetadataBase().IsGoogModule() {
		this_.GetOrCreateName("exports", this_.curMetadata)
	}
}

// Java static method: BuildGlobalNamespace.getTopVarNameFromQualifiedNode(com.google.javascript.rhino.Node n)
func GlobalNamespace_BuildGlobalNamespace_GetTopVarNameFromQualifiedNode(n pkg_rhino.Node) string {
	_ = n
	var cur pkg_rhino.Node = n
	_ = cur
	for cur.IsGetProp() {
		cur = cur.GetFirstChild()
	}
	switch cur.TokenSwitchOrd_() {
	case int32(32) /* Token.NAME */ :
		var s string = cur.GetString()
		_ = s
		{

			if j2gort.CheckNotNullString(s) == "" {
				return j2gort.NullString
			}
			return GlobalNamespace_GetTopVarName(s)
		}

	case int32(37) /* Token.THIS */ :
		return "this"
	case int32(110) /* Token.SUPER */ :
		return "super"
	default:
		return j2gort.NullString
	}
	panic("unreachable switch expression")
}

// Java method: BuildGlobalNamespace.collect(com.google.javascript.jscomp.Scope scope, com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.rhino.Node n)
func (this_ *GlobalNamespace_BuildGlobalNamespace) priv_BuildGlobalNamespace_collect_3070(scope *Scope, chunk *JSChunk, n pkg_rhino.Node) {
	_ = scope
	_ = chunk
	_ = n
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	var name string = j2gort.NullString
	_ = name
	var root pkg_rhino.Node = pkg_rhino.Node(0)
	_ = root
	var isSet bool = false
	_ = isSet
	var type_ *GlobalNamespace_NameProp = StaticGet_GlobalNamespace_NameProp_OTHER_OBJECT()
	_ = type_
	switch n.TokenSwitchOrd_() {
	case int32(95) /* Token.GETTER_DEF */, int32(96) /* Token.SETTER_DEF */, int32(107) /* Token.MEMBER_FUNCTION_DEF */, int32(108) /* Token.MEMBER_FIELD_DEF */ :
		if parent.IsClassMembers() && !(n.IsStaticMember()) {
			return
		}
		name = NodeUtil_GetBestLValueName(n)
		isSet = true
		if n.IsMemberFunctionDef() {
			type_ = StaticGet_GlobalNamespace_NameProp_FUNCTION()
		} else if n.IsMemberFieldDef() {
			type_ = StaticGet_GlobalNamespace_NameProp_OTHER_OBJECT()
		} else {
			type_ = StaticGet_GlobalNamespace_NameProp_GET_SET()
		}
	case int32(100) /* Token.STRING_KEY */ :
		name = j2gort.NullString
		if parent.IsObjectLit() {
			var analysis *GlobalNamespace_ObjLitStringKeyAnalysis = this_.priv_BuildGlobalNamespace_createObjLitStringKeyAnalysis_3071(n)
			_ = analysis
			name = analysis.NameString()
			type_ = analysis.NameType()
			isSet = true
		} else if parent.IsObjectPattern() {
			name = this_.GetNameForObjectPatternKey(n)
			type_ = this_.GetValueType(n.GetFirstChild())
		}
	case int32(32) /* Token.NAME */, int32(26) /* Token.GETPROP */, int32(29) /* Token.OPTCHAIN_GETPROP */ :
		{
			var rvalue pkg_rhino.Node = pkg_rhino.Node(0)
			_ = rvalue
			var grandparent pkg_rhino.Node = pkg_rhino.Node(0)
			_ = grandparent
			switch parent.TokenSwitchOrd_() {
			case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
				isSet = true
				rvalue = n.GetFirstChild()
				if rvalue == 0 {
					type_ = StaticGet_GlobalNamespace_NameProp_OTHER_OBJECT()
				} else {
					type_ = this_.GetValueType(rvalue)
				}
			case int32(51) /* Token.ASSIGN */ :
				if parent.GetFirstChild() == n {
					isSet = true
					type_ = this_.GetValueType(n.GetNext())
				}
			case int32(26) /* Token.GETPROP */ :
				return
			case int32(73) /* Token.FUNCTION */ :
				grandparent = parent.GetParent()
				if (grandparent == 0) || NodeUtil_IsFunctionExpression(parent) {
					return
				}
				isSet = true
				type_ = StaticGet_GlobalNamespace_NameProp_FUNCTION()
			case int32(86) /* Token.CATCH */, int32(71) /* Token.INC */, int32(72) /* Token.DEC */ :
				isSet = true
				type_ = StaticGet_GlobalNamespace_NameProp_OTHER_OBJECT()
			case int32(81) /* Token.FOR_IN */, int32(112) /* Token.FOR_OF */, int32(113) /* Token.FOR_AWAIT_OF */ :
				if parent.GetFirstChild() == n {
					isSet = true
					type_ = StaticGet_GlobalNamespace_NameProp_OTHER_OBJECT()
				}
			case int32(105) /* Token.CLASS */ :
				if parent.GetFirstChild() == n {
					isSet = true
					type_ = StaticGet_GlobalNamespace_NameProp_CLASS()
				}
			case int32(100) /* Token.STRING_KEY */, int32(102) /* Token.ARRAY_PATTERN */, int32(134) /* Token.DEFAULT_VALUE */, int32(129) /* Token.COMPUTED_PROP */, int32(125) /* Token.ITER_REST */, int32(126) /* Token.OBJECT_REST */ :
				if n.IsName() && NodeUtil_IsLhsByDestructuring(n) {
					isSet = true
					type_ = StaticGet_GlobalNamespace_NameProp_OTHER_OBJECT()
				}
			case int32(127) /* Token.ITER_SPREAD */, int32(128) /* Token.OBJECT_SPREAD */ :
			case int32(28) /* Token.CALL */, int32(31) /* Token.OPTCHAIN_CALL */ :
				if n.IsFirstChildOf(parent) && this_.priv_BuildGlobalNamespace_isObjectHasOwnPropertyCall_3083(parent) {
					var qname string = n.GetFirstChild().GetQualifiedName()
					_ = qname
					var globalName *GlobalNamespace_Name = this_.GetOrCreateName(qname, this_.curMetadata)
					_ = globalName
					globalName.priv_Name_setBooleanProperty_3114(StaticGet_GlobalNamespace_NameProp_IS_USED_HAS_OWN_PROPERTY())
				}
			default:
				if NodeUtil_IsAssignmentOp(parent) && (parent.GetFirstChild() == n) {
					isSet = true
					type_ = StaticGet_GlobalNamespace_NameProp_OTHER_OBJECT()
				}
			}
		}
		var topVarName string = GlobalNamespace_BuildGlobalNamespace_GetTopVarNameFromQualifiedNode(n)
		_ = topVarName
		if topVarName == j2gort.NullString {
			return
		}
		root = this_.outer_.priv_GlobalNamespace_getRootNodeForTopVar_3052(topVarName, scope)
		if !(this_.priv_BuildGlobalNamespace_isTopLevelScopeRoot_3073(root)) {
			return
		}
		name = n.GetQualifiedName()
	case int32(28) /* Token.CALL */ :
		if (parent.IsExprResult() && StaticGet_GlobalNamespace_GOOG_PROVIDE().Matches(n.GetFirstChild())) && n.GetSecondChild().IsStringLit() {
			this_.priv_BuildGlobalNamespace_createNamesFromProvide_3072(n.GetSecondChild().GetString())
			return
		}
		return
	default:
		return
	}
	if name == j2gort.NullString {
		return
	}
	if root == 0 {
		root = this_.outer_.priv_GlobalNamespace_getRootNode_3051(name, scope)
		if !(this_.priv_BuildGlobalNamespace_isTopLevelScopeRoot_3073(root)) {
			return
		}
	}
	var nameMetadata ModuleMetadataMap_ModuleMetadata
	_ = nameMetadata
	if root == this_.outer_.globalRoot {
		nameMetadata = nil
	} else {
		nameMetadata = this_.curMetadata
	}
	if isSet {
		var hoistScope *Scope = j2gort.CastPtr[Scope](scope.GetClosestHoistScope())
		_ = hoistScope
		if hoistScope.IsGlobal() || ((root != this_.outer_.globalRoot) && (hoistScope.GetRootNode_0() == this_.curModuleRoot)) {
			this_.HandleSetFromGlobal(chunk, scope, n, name, type_, nameMetadata)
		} else {
			this_.HandleSetFromLocal(chunk, scope, n, name, nameMetadata)
		}
	} else {
		this_.HandleGet(chunk, scope, n, name, nameMetadata)
	}
}

// Java method: BuildGlobalNamespace.createObjLitStringKeyAnalysis(com.google.javascript.rhino.Node stringKeyNode)
func (this_ *GlobalNamespace_BuildGlobalNamespace) priv_BuildGlobalNamespace_createObjLitStringKeyAnalysis_3071(stringKeyNode pkg_rhino.Node) *GlobalNamespace_ObjLitStringKeyAnalysis {
	_ = stringKeyNode
	var nameString string = NodeUtil_GetBestLValueName(stringKeyNode)
	_ = nameString
	if nameString != j2gort.NullString {
		return allocScratchVia_GlobalNamespace_ObjLitStringKeyAnalysis_ForObjLitAssignment(&this_.scratchObjLitKeyAnalysis_, nameString, this_.GetValueType(stringKeyNode.GetOnlyChild()))
	} else {
		var objLitNode pkg_rhino.Node = stringKeyNode.GetParent()
		_ = objLitNode
		if !(objLitNode.IsObjectLit()) {
			j2gort.CheckArgumentMsg(false, j2gort.BoxNode(objLitNode))
		}
		var objLitParentNode pkg_rhino.Node = objLitNode.GetParent()
		_ = objLitParentNode
		if NodeUtil_IsObjectDefinePropertiesDefinition(objLitParentNode) {
			var receiverNode pkg_rhino.Node = objLitParentNode.GetSecondChild()
			_ = receiverNode
			if receiverNode.IsQualifiedName() {
				if !(objLitNode == receiverNode.GetNext()) {
					j2gort.CheckStateMsg(false, j2gort.BoxNode(objLitParentNode))
				}
				nameString = (j2gort.StrValString(receiverNode.GetQualifiedName()) + "." + j2gort.StrValString(stringKeyNode.GetString()))
				return allocScratchVia_GlobalNamespace_ObjLitStringKeyAnalysis_ForObjectDefineProperty(&this_.scratchObjLitKeyAnalysis_, nameString)
			}
		}
		return allocScratchVia_GlobalNamespace_ObjLitStringKeyAnalysis_ForNonReference(&this_.scratchObjLitKeyAnalysis_)
	}
}

// Java method: BuildGlobalNamespace.createNamesFromProvide(java.lang.String namespace)
func (this_ *GlobalNamespace_BuildGlobalNamespace) priv_BuildGlobalNamespace_createNamesFromProvide_3072(namespace string) {
	_ = namespace
	var name *GlobalNamespace_Name = nil
	_ = name
	var dot int32 = int32(0)
	_ = dot
	for dot >= int32(0) {
		dot = j2gort.JStrIndexOfCpFrom(namespace, int32(uint16(46)), (dot + int32(1)))
		var subNamespace string
		_ = subNamespace
		if dot < int32(0) {
			subNamespace = namespace
		} else {
			subNamespace = j2gort.JStrSub(namespace, int32(0), dot)
		}
		if !(!(j2gort.CheckNotNullString(subNamespace) == "")) {
			j2gort.CheckState(false)
		}
		name = this_.GetOrCreateName(subNamespace, nil)
		name.priv_Name_setBooleanProperty_3114(StaticGet_GlobalNamespace_NameProp_IS_PROVIDED())
	}
	var newName *GlobalNamespace_Name = this_.GetOrCreateName(namespace, nil)
	_ = newName
	newName.priv_Name_setBooleanProperty_3114(StaticGet_GlobalNamespace_NameProp_IS_PROVIDED())
}

// Java method: BuildGlobalNamespace.isTopLevelScopeRoot(com.google.javascript.rhino.Node root)
func (this_ *GlobalNamespace_BuildGlobalNamespace) priv_BuildGlobalNamespace_isTopLevelScopeRoot_3073(root pkg_rhino.Node) bool {
	_ = root
	if root == 0 {
		return false
	} else if root == this_.outer_.globalRoot {
		return true
	} else if root == this_.curModuleRoot {
		return true
	}
	return (((this_.curModuleRoot != 0) && this_.curModuleRoot.IsBlock()) && (root == this_.curModuleRoot.GetParent()))
}

// Java method: BuildGlobalNamespace.getNameForObjectPatternKey(com.google.javascript.rhino.Node stringKey)
func (this_ *GlobalNamespace_BuildGlobalNamespace) GetNameForObjectPatternKey(stringKey pkg_rhino.Node) string {
	_ = stringKey
	var parent pkg_rhino.Node = stringKey.GetParent()
	_ = parent
	if !(parent.IsObjectPattern()) {
		j2gort.CheckState(false)
	}
	var patternParent pkg_rhino.Node = parent.GetParent()
	_ = patternParent
	if patternParent.IsAssign() || patternParent.IsDestructuringLhs() {
		var rhs pkg_rhino.Node = patternParent.GetSecondChild()
		_ = rhs
		if (rhs == 0) || !(rhs.IsQualifiedName()) {
			return j2gort.NullString
		}
		return (j2gort.StrValString(rhs.GetQualifiedName()) + "." + j2gort.StrValString(stringKey.GetString()))
	} else {
		return j2gort.NullString
	}
}

// Java method: BuildGlobalNamespace.getValueType(com.google.javascript.rhino.Node n)
func (this_ *GlobalNamespace_BuildGlobalNamespace) GetValueType(n pkg_rhino.Node) *GlobalNamespace_NameProp {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(105) /* Token.CLASS */ :
		return StaticGet_GlobalNamespace_NameProp_CLASS()
	case int32(47) /* Token.OBJECTLIT */ :
		return StaticGet_GlobalNamespace_NameProp_OBJECTLIT()
	case int32(73) /* Token.FUNCTION */ :
		return StaticGet_GlobalNamespace_NameProp_FUNCTION()
	case int32(68) /* Token.OR */ :
		return this_.GetValueType(n.GetLastChild())
	case int32(67) /* Token.HOOK */ :
		var second pkg_rhino.Node = n.GetSecondChild()
		_ = second
		var t *GlobalNamespace_NameProp = this_.GetValueType(second)
		_ = t
		if t != StaticGet_GlobalNamespace_NameProp_OTHER_OBJECT() {
			return t
		}
		var third pkg_rhino.Node = second.GetNext()
		_ = third
		return this_.GetValueType(third)
	default:
	}
	return StaticGet_GlobalNamespace_NameProp_OTHER_OBJECT()
}

// Java method: BuildGlobalNamespace.handleSetFromGlobal(com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.jscomp.Scope scope, com.google.javascript.rhino.Node n, java.lang.String name, com.google.javascript.jscomp.GlobalNamespace.NameProp type, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata)
func (this_ *GlobalNamespace_BuildGlobalNamespace) HandleSetFromGlobal(chunk *JSChunk, scope *Scope, n pkg_rhino.Node, name string, type_ *GlobalNamespace_NameProp, metadata ModuleMetadataMap_ModuleMetadata) {
	_ = chunk
	_ = scope
	_ = n
	_ = name
	_ = type_
	_ = metadata
	if this_.MaybeHandlePrototypePrefix(chunk, scope, n, name, metadata) {
		return
	}
	var nameObj *GlobalNamespace_Name = this_.GetOrCreateName(name, metadata)
	_ = nameObj
	if !(nameObj.IsGetOrSetDefinition()) {
		nameObj.priv_Name_setNameType_3158(type_)
	}
	if n.GetBooleanProp(pkg_rhino.StaticGet_Node_MODULE_EXPORT()) {
		nameObj.priv_Name_setBooleanProperty_3114(StaticGet_GlobalNamespace_NameProp_IS_MODULE_PROP())
	}
	if this_.IsNestedAssign(n.GetParent()) {
		var refType GlobalNamespace_Ref_Type = StaticGet_GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL()
		_ = refType
		this_.priv_BuildGlobalNamespace_addOrConfirmRef_3081(nameObj, n, refType, scope, chunk)
	} else {
		this_.priv_BuildGlobalNamespace_addOrConfirmRef_3081(nameObj, n, StaticGet_GlobalNamespace_Ref_Type_SET_FROM_GLOBAL(), scope, chunk)
		nameObj.priv_Name_setDeclaredTypeKind_3157(this_.priv_BuildGlobalNamespace_getDeclaredTypeKind_3077(n))
	}
}

// Java method: BuildGlobalNamespace.getDeclaredTypeKind(com.google.javascript.rhino.Node n)
func (this_ *GlobalNamespace_BuildGlobalNamespace) priv_BuildGlobalNamespace_getDeclaredTypeKind_3077(n pkg_rhino.Node) *GlobalNamespace_NameProp {
	_ = n
	var valueNode pkg_rhino.Node = NodeUtil_GetRValueOfLValue(n)
	_ = valueNode
	var kind *GlobalNamespace_NameProp = nil
	_ = kind
	if valueNode == 0 {
		kind = StaticGet_GlobalNamespace_NameProp_NOT_A_TYPE()
	} else if valueNode.IsClass() {
		kind = StaticGet_GlobalNamespace_NameProp_CONSTRUCTOR_TYPE()
	} else {
		var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(n)
		_ = info
		if info == nil {
			kind = StaticGet_GlobalNamespace_NameProp_NOT_A_TYPE()
		} else if info.IsConstructor() && valueNode.IsFunction() {
			kind = StaticGet_GlobalNamespace_NameProp_CONSTRUCTOR_TYPE()
		} else if info.IsInterface() && valueNode.IsFunction() {
			kind = StaticGet_GlobalNamespace_NameProp_INTERFACE_TYPE()
		} else if info.HasEnumParameterType() && valueNode.IsObjectLit() {
			kind = StaticGet_GlobalNamespace_NameProp_ENUM_TYPE()
		} else {
			kind = StaticGet_GlobalNamespace_NameProp_NOT_A_TYPE()
		}
	}
	return kind
}

// Java method: BuildGlobalNamespace.handleSetFromLocal(com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.jscomp.Scope scope, com.google.javascript.rhino.Node n, java.lang.String name, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata)
func (this_ *GlobalNamespace_BuildGlobalNamespace) HandleSetFromLocal(chunk *JSChunk, scope *Scope, n pkg_rhino.Node, name string, metadata ModuleMetadataMap_ModuleMetadata) {
	_ = chunk
	_ = scope
	_ = n
	_ = name
	_ = metadata
	if this_.MaybeHandlePrototypePrefix(chunk, scope, n, name, metadata) {
		return
	}
	var nameObj *GlobalNamespace_Name = this_.GetOrCreateName(name, metadata)
	_ = nameObj
	if n.GetBooleanProp(pkg_rhino.StaticGet_Node_MODULE_EXPORT()) {
		nameObj.priv_Name_setBooleanProperty_3114(StaticGet_GlobalNamespace_NameProp_IS_MODULE_PROP())
	}
	if this_.IsNestedAssign(n.GetParent()) {
		this_.priv_BuildGlobalNamespace_addOrConfirmRef_3081(nameObj, n, StaticGet_GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL(), scope, chunk)
	} else {
		this_.priv_BuildGlobalNamespace_addOrConfirmRef_3081(nameObj, n, StaticGet_GlobalNamespace_Ref_Type_SET_FROM_LOCAL(), scope, chunk)
	}
}

// Java method: BuildGlobalNamespace.handleGet(com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.jscomp.Scope scope, com.google.javascript.rhino.Node n, java.lang.String name, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata)
func (this_ *GlobalNamespace_BuildGlobalNamespace) HandleGet(chunk *JSChunk, scope *Scope, n pkg_rhino.Node, name string, metadata ModuleMetadataMap_ModuleMetadata) {
	_ = chunk
	_ = scope
	_ = n
	_ = name
	_ = metadata
	if this_.MaybeHandlePrototypePrefix(chunk, scope, n, name, metadata) {
		return
	}
	var type_ GlobalNamespace_Ref_Type = this_.priv_BuildGlobalNamespace_determineRefTypeForGet_3080(n, n, name)
	_ = type_
	this_.priv_BuildGlobalNamespace_addOrConfirmRef_3081(this_.GetOrCreateName(name, metadata), n, type_, scope, chunk)
}

// Java method: BuildGlobalNamespace.determineRefTypeForGet(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node referenceNode, java.lang.String name)
func (this_ *GlobalNamespace_BuildGlobalNamespace) priv_BuildGlobalNamespace_determineRefTypeForGet_3080(n pkg_rhino.Node, referenceNode pkg_rhino.Node, name string) GlobalNamespace_Ref_Type {
	_ = n
	_ = referenceNode
	_ = name
	var type_ GlobalNamespace_Ref_Type = GlobalNamespace_Ref_Type(0)
	_ = type_
	var parent pkg_rhino.Node = n.GetParent()
	_ = parent
	{
		var lhs pkg_rhino.Node = pkg_rhino.Node(0)
		_ = lhs
		switch parent.TokenSwitchOrd_() {
		case int32(93) /* Token.EXPR_RESULT */, int32(74) /* Token.IF */, int32(78) /* Token.WHILE */, int32(80) /* Token.FOR */, int32(45) /* Token.INSTANCEOF */, int32(25) /* Token.TYPEOF */, int32(87) /* Token.VOID */, int32(19) /* Token.NOT */, int32(20) /* Token.BITNOT */, int32(21) /* Token.POS */, int32(22) /* Token.NEG */, int32(40) /* Token.SHEQ */, int32(4) /* Token.EQ */, int32(41) /* Token.SHNE */, int32(5) /* Token.NE */, int32(6) /* Token.LT */, int32(7) /* Token.LE */, int32(8) /* Token.GT */, int32(9) /* Token.GE */, int32(13) /* Token.ADD */, int32(14) /* Token.SUB */, int32(15) /* Token.MUL */, int32(16) /* Token.DIV */, int32(17) /* Token.MOD */, int32(18) /* Token.EXPONENT */, int32(3) /* Token.BITAND */, int32(1) /* Token.BITOR */, int32(2) /* Token.BITXOR */, int32(10) /* Token.LSH */, int32(11) /* Token.RSH */, int32(12) /* Token.URSH */ :
			type_ = StaticGet_GlobalNamespace_Ref_Type_DIRECT_GET()
		case int32(31) /* Token.OPTCHAIN_CALL */, int32(28) /* Token.CALL */ :
			if n == parent.GetFirstChild() {
				type_ = StaticGet_GlobalNamespace_Ref_Type_CALL_GET()
			} else if this_.priv_BuildGlobalNamespace_isClassDefiningCall_3082(parent) {
				type_ = StaticGet_GlobalNamespace_Ref_Type_DIRECT_GET()
			} else {
				type_ = StaticGet_GlobalNamespace_Ref_Type_ALIASING_GET()
			}
		case int32(23) /* Token.NEW */ :
			if n == parent.GetFirstChild() {
				type_ = StaticGet_GlobalNamespace_Ref_Type_DIRECT_GET()
			} else {
				type_ = StaticGet_GlobalNamespace_Ref_Type_ALIASING_GET()
			}
		case int32(101) /* Token.CAST */, int32(68) /* Token.OR */, int32(69) /* Token.AND */, int32(70) /* Token.COALESCE */ :
			type_ = this_.priv_BuildGlobalNamespace_determineRefTypeForGet_3080(parent, referenceNode, name)
		case int32(32) /* Token.NAME */ :
			if (n != referenceNode) && func() bool { sr_6 := name; se_7 := parent.GetString(); return j2gort.CheckNotNullString(sr_6) == se_7 }() {
				type_ = StaticGet_GlobalNamespace_Ref_Type_DIRECT_GET()
			} else {
				type_ = StaticGet_GlobalNamespace_Ref_Type_ALIASING_GET()
			}
		case int32(50) /* Token.COMMA */, int32(67) /* Token.HOOK */ :
			if n != parent.GetFirstChild() {
				type_ = this_.priv_BuildGlobalNamespace_determineRefTypeForGet_3080(parent, referenceNode, name)
			} else {
				type_ = StaticGet_GlobalNamespace_Ref_Type_DIRECT_GET()
			}
		case int32(24) /* Token.DELPROP */ :
			type_ = StaticGet_GlobalNamespace_Ref_Type_DELETE_PROP()
		case int32(105) /* Token.CLASS */ :
			type_ = StaticGet_GlobalNamespace_Ref_Type_SUBCLASSING_GET()
		case int32(104) /* Token.DESTRUCTURING_LHS */, int32(51) /* Token.ASSIGN */ :
			lhs = n.GetPrevious()
			if lhs == 0 {
				type_ = StaticGet_GlobalNamespace_Ref_Type_ALIASING_GET()
				break
			}
			for lhs.IsCast() {
				lhs = lhs.GetOnlyChild()
			}
			if n != referenceNode {
				if lhs.MatchesQualifiedName_1_string(name) {
					return StaticGet_GlobalNamespace_Ref_Type_DIRECT_GET()
				}
			}
			switch lhs.TokenSwitchOrd_() {
			case int32(32) /* Token.NAME */, int32(26) /* Token.GETPROP */, int32(27) /* Token.GETELEM */, int32(102) /* Token.ARRAY_PATTERN */, int32(103) /* Token.OBJECT_PATTERN */ :
				type_ = StaticGet_GlobalNamespace_Ref_Type_ALIASING_GET()
			default:
				panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(("Unexpected previous sibling of " + j2gort.StrVal(n.GetToken()) + ": " + j2gort.StrVal(n.GetPrevious()))))))
			}
		case int32(103) /* Token.OBJECT_PATTERN */, int32(127) /* Token.ITER_SPREAD */, int32(128) /* Token.OBJECT_SPREAD */, int32(0) /* Token.RETURN */, int32(43) /* Token.THROW */ :
			fallthrough
		default:
			type_ = StaticGet_GlobalNamespace_Ref_Type_ALIASING_GET()
		}
	}
	return type_
}

// Java method: BuildGlobalNamespace.addOrConfirmRef(com.google.javascript.jscomp.GlobalNamespace.Name nameObj, com.google.javascript.rhino.Node node, com.google.javascript.jscomp.GlobalNamespace.Ref.Type refType, com.google.javascript.jscomp.Scope scope, com.google.javascript.jscomp.JSChunk chunk)
func (this_ *GlobalNamespace_BuildGlobalNamespace) priv_BuildGlobalNamespace_addOrConfirmRef_3081(nameObj *GlobalNamespace_Name, node pkg_rhino.Node, refType GlobalNamespace_Ref_Type, scope *Scope, chunk *JSChunk) {
	_ = nameObj
	_ = node
	_ = refType
	_ = scope
	_ = chunk
	var existingRef *GlobalNamespace_Ref = nameObj.GetRefForNode(node)
	_ = existingRef
	if existingRef == nil {
		nameObj.priv_Name_addRef_3116(chunk, scope, node, refType)
	} else {
		var existingRefType GlobalNamespace_Ref_Type = existingRef.type_
		_ = existingRefType
		if !(existingRefType == refType) {
			j2gort.CheckStateFmt(false, "existing ref type: %s expected: %s", j2gort.BoxNullableVal(existingRefType), j2gort.BoxNullableVal(refType))
		}
	}
}

// Java method: BuildGlobalNamespace.isClassDefiningCall(com.google.javascript.rhino.Node callNode)
func (this_ *GlobalNamespace_BuildGlobalNamespace) priv_BuildGlobalNamespace_isClassDefiningCall_3082(callNode pkg_rhino.Node) bool {
	_ = callNode
	var convention CodingConvention = this_.outer_.compiler.GetCodingConvention()
	_ = convention
	var classes *CodingConvention_SubclassRelationship = convention.GetClassesDefinedByCall(callNode)
	_ = classes
	if classes != nil {
		return true
	}
	var className string = convention.GetSingletonGetterClassName(callNode)
	_ = className
	return (className != j2gort.NullString)
}

// Java method: BuildGlobalNamespace.isObjectHasOwnPropertyCall(com.google.javascript.rhino.Node callNode)
func (this_ *GlobalNamespace_BuildGlobalNamespace) priv_BuildGlobalNamespace_isObjectHasOwnPropertyCall_3083(callNode pkg_rhino.Node) bool {
	_ = callNode
	if !(callNode.IsCall() || callNode.IsOptChainCall()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(callNode))
	}
	if !(callNode.HasTwoChildren()) {
		return false
	}
	var callee pkg_rhino.Node = callNode.GetFirstChild()
	_ = callee
	if !(callee.IsGetProp()) && !(callee.IsOptChainGetProp()) {
		return false
	}
	var receiver pkg_rhino.Node = callee.GetFirstChild()
	_ = receiver
	var propName string = callee.GetString()
	_ = propName
	return ((("hasOwnProperty" == propName) || ("propertyIsEnumerable" == propName)) && receiver.IsQualifiedName())
}

// Java method: BuildGlobalNamespace.maybeHandlePrototypePrefix(com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.jscomp.Scope scope, com.google.javascript.rhino.Node n, java.lang.String name, com.google.javascript.jscomp.modules.ModuleMetadataMap.ModuleMetadata metadata)
func (this_ *GlobalNamespace_BuildGlobalNamespace) MaybeHandlePrototypePrefix(chunk *JSChunk, scope *Scope, n pkg_rhino.Node, name string, metadata ModuleMetadataMap_ModuleMetadata) bool {
	_ = chunk
	_ = scope
	_ = n
	_ = name
	_ = metadata
	var numLevelsToRemove int32 = int32(0)
	_ = numLevelsToRemove
	var prefix string = j2gort.NullString
	_ = prefix
	if javalang.StringEndsWith(j2gort.CheckNotNullString(name), ".prototype") {
		numLevelsToRemove = int32(1)
		prefix = j2gort.JStrSub(name, int32(0), (j2gort.JStrLen(name) - int32(10)))
	} else {
		var i int32 = j2gort.JStrIndexOf(name, ".prototype.")
		_ = i
		if i == int32(-1) {
			return false
		}
		prefix = j2gort.JStrSub(name, int32(0), i)
		numLevelsToRemove = int32(2)
		i = j2gort.JStrIndexOfCpFrom(name, int32(uint16(46)), (i + int32(11)))
		for i >= int32(0) {
			numLevelsToRemove = numLevelsToRemove + 1
			i = j2gort.JStrIndexOfCpFrom(name, int32(uint16(46)), (i + int32(1)))
		}
	}
	if NodeUtil_MayBeObjectLitKey(n) {
		return true
	}
	{
		var i_2 int32 = int32(0)
		_ = i_2
		for ; i_2 < numLevelsToRemove; i_2 = i_2 + 1 {
			n = n.GetFirstChild()
		}
	}
	this_.priv_BuildGlobalNamespace_addOrConfirmRef_3081(this_.GetOrCreateName(prefix, metadata), n, StaticGet_GlobalNamespace_Ref_Type_PROTOTYPE_GET(), scope, chunk)
	return true
}

// Java method: BuildGlobalNamespace.isNestedAssign(com.google.javascript.rhino.Node parent)
func (this_ *GlobalNamespace_BuildGlobalNamespace) IsNestedAssign(parent pkg_rhino.Node) bool {
	_ = parent
	return (parent.IsAssign() && !(parent.GetParent().IsExprResult()))
}

// Java method: BuildGlobalNamespace.getOrCreateName(java.lang.String name, com.google.javascript.jscomp.modules.ModuleMetadataMap.@org.jspecify.annotations.Nullable ModuleMetadata metadata)
func (this_ *GlobalNamespace_BuildGlobalNamespace) GetOrCreateName(name string, metadata ModuleMetadataMap_ModuleMetadata) *GlobalNamespace_Name {
	_ = name
	_ = metadata
	var node *GlobalNamespace_Name
	_ = node
	if metadata == nil || metadata.IsJavaNil_() {
		node = this_.outer_.nameMap.Get(name)
	} else {
		node = this_.outer_.priv_GlobalNamespace_getModuleName_3047(metadata, name)
	}
	if node == nil {
		var i int32 = j2gort.JStrLastIndexOfCp(name, int32(uint16(46)))
		_ = i
		if i >= int32(0) {
			var parentName string = j2gort.JStrSub(name, int32(0), i)
			_ = parentName
			var parent *GlobalNamespace_Name = this_.GetOrCreateName(parentName, metadata)
			_ = parent
			node = parent.AddProperty(j2gort.JStrSubFrom(name, (i+int32(1))), this_.outer_.sourceKind)
			if metadata == nil || metadata.IsJavaNil_() {
				this_.outer_.nameMap.Put(name, node)
			} else {
				this_.outer_.priv_GlobalNamespace_putModuleName_3048(metadata, name, node)
			}
		} else {
			node = NewGlobalNamespace_Name(this_.outer_, name, nil, this_.outer_.sourceKind)
			if metadata == nil || metadata.IsJavaNil_() {
				this_.outer_.globalNames = append(this_.outer_.globalNames, node)
				this_.outer_.nameMap.Put(name, node)
			} else {
				this_.outer_.priv_GlobalNamespace_putModuleName_3048(metadata, name, node)
			}
		}
	}
	return node
}

// Java enum: com.google.javascript.jscomp.GlobalNamespace.Inlinability
type GlobalNamespace_Inlinability uint8

const (
	enum_GlobalNamespace_Inlinability_INLINE_COMPLETELY                            GlobalNamespace_Inlinability = 1
	enum_GlobalNamespace_Inlinability_INLINE_BUT_KEEP_DECLARATION                  GlobalNamespace_Inlinability = 2
	enum_GlobalNamespace_Inlinability_INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY GlobalNamespace_Inlinability = 3
	enum_GlobalNamespace_Inlinability_DO_NOT_INLINE                                GlobalNamespace_Inlinability = 4
)

var enumBaseTable_GlobalNamespace_Inlinability = [...]j2gort.EnumBase{
	{Name_: "INLINE_COMPLETELY", Ordinal_: 0},
	{Name_: "INLINE_BUT_KEEP_DECLARATION", Ordinal_: 1},
	{Name_: "INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY", Ordinal_: 2},
	{Name_: "DO_NOT_INLINE", Ordinal_: 3},
}
var enumValues_GlobalNamespace_Inlinability = []GlobalNamespace_Inlinability{1, 2, 3, 4}

var enumField_GlobalNamespace_Inlinability_shouldInlineUsages = [...]bool{
	true,
	true,
	true,
	false,
}

func (this_ GlobalNamespace_Inlinability) FieldGet_shouldInlineUsages() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_GlobalNamespace_Inlinability_shouldInlineUsages[this_-1]
}

var enumField_GlobalNamespace_Inlinability_shouldRemoveDeclaration = [...]bool{
	true,
	false,
	false,
	false,
}

func (this_ GlobalNamespace_Inlinability) FieldGet_shouldRemoveDeclaration() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_GlobalNamespace_Inlinability_shouldRemoveDeclaration[this_-1]
}

var enumField_GlobalNamespace_Inlinability_canCollapse = [...]bool{
	true,
	true,
	true,
	false,
}

func (this_ GlobalNamespace_Inlinability) FieldGet_canCollapse() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumField_GlobalNamespace_Inlinability_canCollapse[this_-1]
}

func StaticGet_GlobalNamespace_Inlinability_INLINE_COMPLETELY() GlobalNamespace_Inlinability {
	return enum_GlobalNamespace_Inlinability_INLINE_COMPLETELY
}
func StaticGet_GlobalNamespace_Inlinability_INLINE_BUT_KEEP_DECLARATION() GlobalNamespace_Inlinability {
	return enum_GlobalNamespace_Inlinability_INLINE_BUT_KEEP_DECLARATION
}
func StaticGet_GlobalNamespace_Inlinability_INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY() GlobalNamespace_Inlinability {
	return enum_GlobalNamespace_Inlinability_INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY
}
func StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE() GlobalNamespace_Inlinability {
	return enum_GlobalNamespace_Inlinability_DO_NOT_INLINE
}

func GlobalNamespace_Inlinability_Values() []GlobalNamespace_Inlinability {
	return j2gort.CloneSlice(enumValues_GlobalNamespace_Inlinability)
}

func GlobalNamespace_Inlinability_ValueOf(name string) GlobalNamespace_Inlinability {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_GlobalNamespace_Inlinability {
		if enumBaseTable_GlobalNamespace_Inlinability[i_].Name_ == name {
			return GlobalNamespace_Inlinability(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.GlobalNamespace.Inlinability." + name))
}

// Java method: Inlinability.shouldInlineUsages()
func (this_ GlobalNamespace_Inlinability) ShouldInlineUsages() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return this_.FieldGet_shouldInlineUsages()
}

// Java method: Inlinability.shouldRemoveDeclaration()
func (this_ GlobalNamespace_Inlinability) ShouldRemoveDeclaration() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return this_.FieldGet_shouldRemoveDeclaration()
}

// Java method: Inlinability.canCollapse()
func (this_ GlobalNamespace_Inlinability) CanCollapse() bool {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return this_.FieldGet_canCollapse()
}

func (this_ GlobalNamespace_Inlinability) IsJavaNil_() bool                              { return this_ == 0 }
func (GlobalNamespace_Inlinability) HasIdentityEquals_()                                 {}
func (GlobalNamespace_Inlinability) Implements_j2gort_Ext_java_lang_Enum()               {}
func (GlobalNamespace_Inlinability) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ GlobalNamespace_Inlinability) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (GlobalNamespace_Inlinability) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GlobalNamespace$Inlinability", "Inlinability"
}
func (this_ GlobalNamespace_Inlinability) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_GlobalNamespace_Inlinability[this_-1]
}
func (this_ GlobalNamespace_Inlinability) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ GlobalNamespace_Inlinability) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ GlobalNamespace_Inlinability) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_GlobalNamespace_Inlinability[this_-1].Name_
}
func (this_ GlobalNamespace_Inlinability) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_GlobalNamespace_Inlinability[this_-1].Name_
}
func (this_ GlobalNamespace_Inlinability) Equals(other_ any) bool {
	o_, ok_ := other_.(GlobalNamespace_Inlinability)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ GlobalNamespace_Inlinability) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_GlobalNamespace_Inlinability[this_-1].HashCode()
}
func (this_ GlobalNamespace_Inlinability) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_GlobalNamespace_Inlinability[this_-1].CompareTo(other_)
}
func (this_ GlobalNamespace_Inlinability) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_GlobalNamespace_Inlinability[this_-1].Name_)
}
func init() {
	vals_ := []any{GlobalNamespace_Inlinability(1), GlobalNamespace_Inlinability(2), GlobalNamespace_Inlinability(3), GlobalNamespace_Inlinability(4)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.GlobalNamespace.Inlinability", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.GlobalNamespace$Inlinability", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GlobalNamespace.Inlinability", "com.google.javascript.jscomp.GlobalNamespace$Inlinability", func(v_ any) bool {
		return j2gort.InstanceOf[GlobalNamespace_Inlinability](v_)
	})
}

// Java class: com.google.javascript.jscomp.GlobalNamespace.Name
type GlobalNamespace_Name struct {
	outer_                                          *GlobalNamespace
	baseName                                        string                                                   `gson:"baseName"`
	parent                                          *GlobalNamespace_Name                                    `gson:"parent"`
	fullName                                        string                                                   `gson:"fullName"`
	prop0                                           *GlobalNamespace_Name                                    `gson:"prop0"`
	prop1                                           *GlobalNamespace_Name                                    `gson:"prop1"`
	prop2                                           *GlobalNamespace_Name                                    `gson:"prop2"`
	prop3                                           *GlobalNamespace_Name                                    `gson:"prop3"`
	props                                           []*GlobalNamespace_Name                                  `gson:"props"`
	declaration                                     *GlobalNamespace_Ref                                     `gson:"declaration"`
	initialization                                  *GlobalNamespace_Ref                                     `gson:"initialization"`
	ref0                                            *GlobalNamespace_Ref                                     `gson:"ref0"`
	ref1                                            *GlobalNamespace_Ref                                     `gson:"ref1"`
	ref2                                            *GlobalNamespace_Ref                                     `gson:"ref2"`
	ref3                                            *GlobalNamespace_Ref                                     `gson:"ref3"`
	refsForNodeMap                                  *j2gort.OrderedMap[pkg_rhino.Node, *GlobalNamespace_Ref] `gson:"refsForNodeMap"`
	firstDeclarationJSDocInfo                       *pkg_rhino.JSDocInfo                                     `gson:"firstDeclarationJSDocInfo"`
	firstQnameDeclarationWithoutAssignmentJsDocInfo *pkg_rhino.JSDocInfo                                     `gson:"firstQnameDeclarationWithoutAssignmentJsDocInfo"`
	propCount                                       int32                                                    `gson:"propCount"`
	globalSets                                      int32                                                    `gson:"globalSets"`
	localSets                                       int32                                                    `gson:"localSets"`
	localSetsWithNoCollapse                         int32                                                    `gson:"localSetsWithNoCollapse"`
	aliasingGets                                    int32                                                    `gson:"aliasingGets"`
	totalGets                                       int32                                                    `gson:"totalGets"`
	callGets                                        int32                                                    `gson:"callGets"`
	deleteProps                                     int32                                                    `gson:"deleteProps"`
	subclassingGets                                 int32                                                    `gson:"subclassingGets"`
	propertyBitSet                                  int32                                                    `gson:"propertyBitSet"`
	idHash_                                         atomic.Int32
}

func (this_ *GlobalNamespace_Name) BindSelf_(s any) {
	_ = s
}
func (this_ *GlobalNamespace_Name) IsJavaNil_() bool { return this_ == nil }
func (*GlobalNamespace_Name) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GlobalNamespace$Name", "Name"
}
func (this_ *GlobalNamespace_Name) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GlobalNamespace_Name) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GlobalNamespace.Name", "com.google.javascript.jscomp.GlobalNamespace$Name", func(v_ any) bool {
		return j2gort.InstanceOf[*GlobalNamespace_Name](v_)
	})
}
func (this_ *GlobalNamespace_Name) SetOuter_(o *GlobalNamespace) { this_.outer_ = o }
func (this_ *GlobalNamespace_Name) GetOuter_() *GlobalNamespace  { return this_.outer_ }
func (*GlobalNamespace_Name) Implements_rhino_StaticSlot()       {}

func (this_ *GlobalNamespace_Name) initFields_GlobalNamespace_Name() {
	_ = this_
	this_.baseName = j2gort.NullString
	this_.fullName = j2gort.NullString

}

// Java constructor: Name(java.lang.String name, com.google.javascript.jscomp.GlobalNamespace.Name parent, com.google.javascript.jscomp.GlobalNamespace.SourceKind sourceKind)
func NewGlobalNamespace_Name(outer_ *GlobalNamespace, name string, parent *GlobalNamespace_Name, sourceKind GlobalNamespace_SourceKind) *GlobalNamespace_Name {
	this_ := &GlobalNamespace_Name{}
	this_.outer_ = outer_
	InitGlobalNamespace_Name(this_, name, parent, sourceKind)
	return this_
}

func InitGlobalNamespace_Name(this_ *GlobalNamespace_Name, name string, parent *GlobalNamespace_Name, sourceKind GlobalNamespace_SourceKind) {
	_ = name
	_ = parent
	_ = sourceKind
	this_.initFields_GlobalNamespace_Name()
	this_.baseName = name
	this_.parent = parent
	if parent == nil {
		this_.fullName = name
	} else {
		this_.fullName = j2gort.NullString
	}
	this_.priv_Name_setBooleanProperty_3114(StaticGet_GlobalNamespace_NameProp_NOT_A_TYPE())
	this_.priv_Name_setBooleanProperty_3114(StaticGet_GlobalNamespace_NameProp_OTHER_OBJECT())
	{
		var iifeCond_1 bool
		{

			re_8 := sourceKind
			ae_9 := StaticGet_GlobalNamespace_SourceKind_EXTERN()
			if re_8 == 0 {
				j2gort.ThrowNPE()
			}
			iifeCond_1 = re_8 == ae_9
		}
		if iifeCond_1 {
			this_.priv_Name_setBooleanProperty_3114(StaticGet_GlobalNamespace_NameProp_IS_EXTERN())
		}
	}
}

// Java method: Name.addProperty(java.lang.String name, com.google.javascript.jscomp.GlobalNamespace.SourceKind sourceKind)
func (this_ *GlobalNamespace_Name) AddProperty(name string, sourceKind GlobalNamespace_SourceKind) *GlobalNamespace_Name {
	_ = name
	_ = sourceKind
	var val_10 GlobalNamespace_Name
	val_10.outer_ = this_.outer_
	InitGlobalNamespace_Name(&val_10, name, this_, sourceKind)
	var node *GlobalNamespace_Name = &val_10
	_ = node
	var c int32 = this_.propCount
	_ = c
	var arr []*GlobalNamespace_Name = this_.props
	_ = arr
	if arr == nil {
		switch c {
		case int32(0):
			this_.prop0 = node
			this_.propCount = int32(1)
			return node
		case int32(1):
			this_.prop1 = node
			this_.propCount = int32(2)
			return node
		case int32(2):
			this_.prop2 = node
			this_.propCount = int32(3)
			return node
		case int32(3):
			this_.prop3 = node
			this_.propCount = int32(4)
			return node
		default:
			arr = make([]*GlobalNamespace_Name, int32(16))
			{
				a_11 := arr
				i_12 := int32(0)
				a_11[i_12] = this_.prop0
			}
			{
				a_13 := arr
				i_14 := int32(1)
				a_13[i_14] = this_.prop1
			}
			{
				a_15 := arr
				i_16 := int32(2)
				a_15[i_16] = this_.prop2
			}
			{
				a_17 := arr
				i_18 := int32(3)
				a_17[i_18] = this_.prop3
			}
			{
				a_19 := arr
				i_20 := int32(4)
				a_19[i_20] = node
			}
			this_.prop0 = nil
			this_.prop1 = nil
			this_.prop2 = nil
			this_.prop3 = nil
			this_.props = arr
			this_.propCount = int32(5)
			return node
		}
	}
	if c == int32(len(arr)) {
		{

			arrSrc_21 := arr
			arrLen_22 := (c * int32(2))
			j2gort.CheckNotNullSlice(arrSrc_21)
			if arrLen_22 < 0 {
				panic(j2gort.ConstructShimThrowable("java.lang.NegativeArraySizeException"))
			}
			arrOut_23 := make([]*GlobalNamespace_Name, arrLen_22)
			copy(arrOut_23, arrSrc_21)
			arr = arrOut_23
		}

		this_.props = arr
	}
	{
		a_24 := arr
		i_25 := c
		a_24[i_25] = node
	}
	this_.propCount = (c + int32(1))
	return node
}

// Java method: Name.hasProps()
func (this_ *GlobalNamespace_Name) HasProps() bool {
	return (this_.propCount != int32(0))
}

// Java method: Name.getPropCount()
func (this_ *GlobalNamespace_Name) GetPropCount() int32 {
	return this_.propCount
}

// Java method: Name.getPropAt(int index)
func (this_ *GlobalNamespace_Name) GetPropAt(index int32) *GlobalNamespace_Name {
	_ = index
	if (index < int32(0)) || (index >= this_.propCount) {
		panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", j2gort.ValueOfInt(index))))
	}
	var arr []*GlobalNamespace_Name = this_.props
	_ = arr
	if arr != nil {
		return arr[index]
	}
	switch index {
	case int32(0):
		return this_.prop0
	case int32(1):
		return this_.prop1
	case int32(2):
		return this_.prop2
	case int32(3):
		return this_.prop3
	default:
		panic(j2gort.AsIface[javalang.IIndexOutOfBoundsException](j2gort.ConstructShimThrowable("java.lang.IndexOutOfBoundsException", j2gort.ValueOfInt(index))))
	}
	panic("unreachable switch expression")
}

// Java method: Name.getBaseName()
func (this_ *GlobalNamespace_Name) GetBaseName() string {
	return this_.baseName
}

// Java method: Name.inExterns()
func (this_ *GlobalNamespace_Name) InExterns() bool {
	return this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_IS_EXTERN())
}

// Java method: Name.subclassingGetCount()
func (this_ *GlobalNamespace_Name) SubclassingGetCount() int32 {
	return this_.subclassingGets
}

// Java method: Name.getName()
func (this_ *GlobalNamespace_Name) GetName() string {
	return this_.GetFullName()
}

// Java method: Name.getFullName()
func (this_ *GlobalNamespace_Name) GetFullName() string {
	var fn string = this_.fullName
	_ = fn
	if fn == j2gort.NullString {
		if this_.parent == nil {
			fn = this_.baseName
		} else {
			fn = (j2gort.StrValString(this_.parent.GetFullName()) + "." + j2gort.StrValString(this_.baseName))
		}
		this_.fullName = fn
	}
	return fn
}

// Java method: Name.usesHasOwnProperty()
func (this_ *GlobalNamespace_Name) UsesHasOwnProperty() bool {
	return this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_IS_USED_HAS_OWN_PROPERTY())
}

// Java method: Name.getDeclaration()
func (this_ *GlobalNamespace_Name) GetDeclaration() any {
	return j2gort.PtrToAny(this_.declaration)
}

// Java method: Name.getInitialization()
func (this_ *GlobalNamespace_Name) GetInitialization() *GlobalNamespace_Ref {
	return this_.initialization
}

// Java method: Name.isFunction()
func (this_ *GlobalNamespace_Name) IsFunction() bool {
	return this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_FUNCTION())
}

// Java method: Name.isClass()
func (this_ *GlobalNamespace_Name) IsClass() bool {
	return this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_CLASS())
}

// Java method: Name.isObjectLiteral()
func (this_ *GlobalNamespace_Name) IsObjectLiteral() bool {
	return this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_OBJECTLIT())
}

// Java method: Name.getAliasingGets()
func (this_ *GlobalNamespace_Name) GetAliasingGets() int32 {
	return this_.aliasingGets
}

// Java method: Name.getSubclassingGets()
func (this_ *GlobalNamespace_Name) GetSubclassingGets() int32 {
	return this_.subclassingGets
}

// Java method: Name.getLocalSets()
func (this_ *GlobalNamespace_Name) GetLocalSets() int32 {
	return this_.localSets
}

// Java method: Name.getGlobalSets()
func (this_ *GlobalNamespace_Name) GetGlobalSets() int32 {
	return this_.globalSets
}

// Java method: Name.getTotalSets()
func (this_ *GlobalNamespace_Name) GetTotalSets() int32 {
	return (this_.globalSets + this_.localSets)
}

// Java method: Name.getCallGets()
func (this_ *GlobalNamespace_Name) GetCallGets() int32 {
	return this_.callGets
}

// Java method: Name.getTotalGets()
func (this_ *GlobalNamespace_Name) GetTotalGets() int32 {
	return this_.totalGets
}

// Java method: Name.getDeleteProps()
func (this_ *GlobalNamespace_Name) GetDeleteProps() int32 {
	return this_.deleteProps
}

// Java method: Name.getParent()
func (this_ *GlobalNamespace_Name) GetParent() *GlobalNamespace_Name {
	return this_.parent
}

// Java method: Name.getDeepestCommonAncestorChunk(com.google.javascript.jscomp.JSChunkGraph chunkGraph)
func (this_ *GlobalNamespace_Name) GetDeepestCommonAncestorChunk(chunkGraph *JSChunkGraph) *JSChunk {
	_ = chunkGraph
	var commonAncestor *JSChunk = j2gort.CastPtr[GlobalNamespace_Ref](this_.GetDeclaration()).GetChunk()
	_ = commonAncestor
	if this_.ref0 != nil {
		commonAncestor = chunkGraph.GetDeepestCommonDependencyInclusive_2(this_.ref0.GetChunk(), commonAncestor)
		if this_.ref1 != nil {
			commonAncestor = chunkGraph.GetDeepestCommonDependencyInclusive_2(this_.ref1.GetChunk(), commonAncestor)
			if this_.ref2 != nil {
				commonAncestor = chunkGraph.GetDeepestCommonDependencyInclusive_2(this_.ref2.GetChunk(), commonAncestor)
				if this_.ref3 != nil {
					commonAncestor = chunkGraph.GetDeepestCommonDependencyInclusive_2(this_.ref3.GetChunk(), commonAncestor)
				}
			}
		}
		return commonAncestor
	}
	if this_.refsForNodeMap != nil {
		{
			cur_26 := this_.refsForNodeMap.IterValues()
			for cur_26.HasNext() {

				r := cur_26.
					Next()
				_ = r
				commonAncestor = chunkGraph.GetDeepestCommonDependencyInclusive_2(r.GetChunk(), commonAncestor)
			}
		}
	}
	return commonAncestor
}

// Java method: Name.getScope()
func (this_ *GlobalNamespace_Name) GetScope() any {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: Name.setBooleanProperty(com.google.javascript.jscomp.GlobalNamespace.NameProp property)
func (this_ *GlobalNamespace_Name) priv_Name_setBooleanProperty_3114(property *GlobalNamespace_NameProp) {
	_ = property
	this_.propertyBitSet = (this_.propertyBitSet | property.bit)
}

// Java method: Name.getBooleanProperty(com.google.javascript.jscomp.GlobalNamespace.NameProp property)
func (this_ *GlobalNamespace_Name) priv_Name_getBooleanProperty_3115(property *GlobalNamespace_NameProp) bool {
	_ = property
	return ((this_.propertyBitSet & property.bit) != int32(0))
}

// Java method: Name.addRef(com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.jscomp.Scope scope, com.google.javascript.rhino.Node node, com.google.javascript.jscomp.GlobalNamespace.Ref.Type type)
func (this_ *GlobalNamespace_Name) priv_Name_addRef_3116(chunk *JSChunk, scope *Scope, node pkg_rhino.Node, type_ GlobalNamespace_Ref_Type) {
	_ = chunk
	_ = scope
	_ = node
	_ = type_
	this_.priv_Name_checkNoExistingRefsForNode_3117(node)
	var ref *GlobalNamespace_Ref = this_.priv_Name_createNewRef_3118(chunk, scope, node, type_)
	_ = ref
	this_.priv_Name_putRef_3119(node, ref)
	this_.priv_Name_updateStateForAddedRef_3123(ref)
}

// Java method: Name.checkNoExistingRefsForNode(com.google.javascript.rhino.Node node)
func (this_ *GlobalNamespace_Name) priv_Name_checkNoExistingRefsForNode_3117(node pkg_rhino.Node) {
	_ = node
	if this_.ref0 != nil {
		if !(this_.ref0.node != node) {
			j2gort.CheckStateFmt(false, "Ref already exists for node: %s", j2gort.PtrToAny(this_.ref0))
		}
		if this_.ref1 != nil {
			if !(this_.ref1.node != node) {
				j2gort.CheckStateFmt(false, "Ref already exists for node: %s", j2gort.PtrToAny(this_.ref1))
			}
			if this_.ref2 != nil {
				if !(this_.ref2.node != node) {
					j2gort.CheckStateFmt(false, "Ref already exists for node: %s", j2gort.PtrToAny(this_.ref2))
				}
				if this_.ref3 != nil {
					if !(this_.ref3.node != node) {
						j2gort.CheckStateFmt(false, "Ref already exists for node: %s", j2gort.PtrToAny(this_.ref3))
					}
				}
			}
		}
		return
	}
	if this_.refsForNodeMap != nil {
		var refForNode *GlobalNamespace_Ref = this_.refsForNodeMap.Get(node)
		_ = refForNode
		if !(refForNode == nil) {
			j2gort.CheckStateFmt(false, "Ref already exists for node: %s", j2gort.PtrToAny(refForNode))
		}
	}
}

// Java method: Name.createNewRef(com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.jscomp.Scope scope, com.google.javascript.rhino.Node node, com.google.javascript.jscomp.GlobalNamespace.Ref.Type type)
func (this_ *GlobalNamespace_Name) priv_Name_createNewRef_3118(chunk *JSChunk, scope *Scope, node pkg_rhino.Node, type_ GlobalNamespace_Ref_Type) *GlobalNamespace_Ref {
	_ = chunk
	_ = scope
	_ = node
	_ = type_
	return allocPooled_GlobalNamespace_Ref(this_.outer_, chunk, j2gort.CheckNotNullPtr(scope), j2gort.CheckNotNullNode(node), type_)
}

// Java method: Name.putRef(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.GlobalNamespace.Ref ref)
func (this_ *GlobalNamespace_Name) priv_Name_putRef_3119(node pkg_rhino.Node, ref *GlobalNamespace_Ref) {
	_ = node
	_ = ref
	if this_.refsForNodeMap != nil {
		cn_ := node
		if this_.refsForNodeMap.Size() == int32(256) && !this_.refsForNodeMap.ContainsKey(cn_) {
			big_ := j2gort.AcquireOrderedMapPtrCap[pkg_rhino.Node, *GlobalNamespace_Ref](int32(512))
			big_.PutAllMapPtr(this_.refsForNodeMap)
			j2gort.ReleaseOrderedMap(this_.refsForNodeMap)
			this_.refsForNodeMap = big_
		}
		this_.refsForNodeMap.Put(cn_, ref)
		return

	}
	if this_.ref0 == nil {
		this_.ref0 = ref
		return
	}
	if this_.ref0.node == node {
		this_.ref0 = ref
		return
	}
	this_.priv_Name_setBooleanProperty_3114(StaticGet_GlobalNamespace_NameProp_MULTI_REF_MODE())
	if this_.ref1 == nil {
		this_.ref1 = ref
		return
	}
	if this_.ref1.node == node {
		this_.ref1 = ref
		return
	}
	if this_.ref2 == nil {
		this_.ref2 = ref
		return
	}
	if this_.ref2.node == node {
		this_.ref2 = ref
		return
	}
	if this_.ref3 == nil {
		this_.ref3 = ref
		return
	}
	if this_.ref3.node == node {
		this_.ref3 = ref
		return
	}
	var r0 *GlobalNamespace_Ref = this_.ref0
	_ = r0
	var r1 *GlobalNamespace_Ref = this_.ref1
	_ = r1
	var r2 *GlobalNamespace_Ref = this_.ref2
	_ = r2
	var r3 *GlobalNamespace_Ref = this_.ref3
	_ = r3
	this_.ref0 = nil
	this_.ref1 = nil
	this_.ref2 = nil
	this_.ref3 = nil
	this_.refsForNodeMap = j2gort.AcquireOrderedMapPtrCap[pkg_rhino.Node, *GlobalNamespace_Ref](int32(8))
	this_.refsForNodeMap.Put(r0.node, r0)
	this_.refsForNodeMap.Put(r1.node, r1)
	this_.refsForNodeMap.Put(r2.node, r2)
	this_.refsForNodeMap.Put(r3.node, r3)
	this_.refsForNodeMap.Put(node, ref)
}

// Java method: Name.findInlineRefForNode(com.google.javascript.rhino.Node node)
func (this_ *GlobalNamespace_Name) priv_Name_findInlineRefForNode_3120(node pkg_rhino.Node) *GlobalNamespace_Ref {
	_ = node
	if this_.ref0 != nil {
		if this_.ref0.node == node {
			return this_.ref0
		}
		if this_.ref1 != nil {
			if this_.ref1.node == node {
				return this_.ref1
			}
			if this_.ref2 != nil {
				if this_.ref2.node == node {
					return this_.ref2
				}
				if (this_.ref3 != nil) && (this_.ref3.node == node) {
					return this_.ref3
				}
			}
		}
	}
	return nil
}

// Java method: Name.addSingleRefForTesting(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.GlobalNamespace.Ref.Type type)
func (this_ *GlobalNamespace_Name) AddSingleRefForTesting(node pkg_rhino.Node, type_ GlobalNamespace_Ref_Type) *GlobalNamespace_Ref {
	_ = node
	_ = type_
	val_28 := allocPooled_GlobalNamespace_Ref(this_.outer_, nil, nil, node, type_)
	var ref *GlobalNamespace_Ref = val_28
	_ = ref
	this_.priv_Name_putRef_3119(node, ref)
	this_.priv_Name_updateStateForAddedRef_3123(ref)
	return ref
}

// Java method: Name.addAliasingGetClonedFromDeclaration(com.google.javascript.rhino.Node newRefNode)
func (this_ *GlobalNamespace_Name) AddAliasingGetClonedFromDeclaration(newRefNode pkg_rhino.Node) {
	_ = newRefNode
	var declRef *GlobalNamespace_Ref = j2gort.CheckNotNullPtr(this_.declaration)
	_ = declRef
	this_.priv_Name_addRef_3116(declRef.chunk, declRef.scope, newRefNode, StaticGet_GlobalNamespace_Ref_Type_ALIASING_GET())
}

// Java method: Name.updateStateForAddedRef(com.google.javascript.jscomp.GlobalNamespace.Ref ref)
func (this_ *GlobalNamespace_Name) priv_Name_updateStateForAddedRef_3123(ref *GlobalNamespace_Ref) {
	_ = ref
	switch (ref.type_).Ordinal() {
	case int32(2) /* Type.GET_AND_SET_FROM_GLOBAL */, int32(0) /* Type.SET_FROM_GLOBAL */ :
		if this_.declaration == nil {
			this_.declaration = ref
		}
		if (this_.initialization == nil) && !(ref.IsUninitializedDeclaration()) {
			this_.initialization = ref
		}
		if this_.firstDeclarationJSDocInfo == nil {
			this_.firstDeclarationJSDocInfo = this_.priv_Name_getDocInfoForDeclaration_3165(ref)
		}
		this_.globalSets = this_.globalSets + 1
		{
			var iifeCond_1 bool
			{

				re_29 := ref.type_
				ae_30 := StaticGet_GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL()
				if re_29 == 0 {
					j2gort.ThrowNPE()
				}
				iifeCond_1 = re_29 == ae_30
			}
			if iifeCond_1 {
				this_.aliasingGets = this_.aliasingGets + 1
				this_.totalGets = this_.totalGets + 1
			}
		}
	case int32(3) /* Type.GET_AND_SET_FROM_LOCAL */, int32(1) /* Type.SET_FROM_LOCAL */ :
		this_.localSets = this_.localSets + 1
		var info *pkg_rhino.JSDocInfo
		_ = info
		if ref.GetNode() == 0 {
			info = nil
		} else {
			info = NodeUtil_GetBestJSDocInfo(ref.GetNode())
		}
		if (info != nil) && info.IsNoCollapse() {
			this_.localSetsWithNoCollapse = this_.localSetsWithNoCollapse + 1
		}
		{
			var iifeCond_1 bool
			{

				re_31 := ref.type_
				ae_32 := StaticGet_GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL()
				if re_31 == 0 {
					j2gort.ThrowNPE()
				}
				iifeCond_1 = re_31 == ae_32
			}
			if iifeCond_1 {
				this_.aliasingGets = this_.aliasingGets + 1
				this_.totalGets = this_.totalGets + 1
			}
		}
	case int32(4) /* Type.PROTOTYPE_GET */, int32(6) /* Type.DIRECT_GET */ :
		var node pkg_rhino.Node = ref.GetNode()
		_ = node
		if (this_.firstQnameDeclarationWithoutAssignmentJsDocInfo == nil) && GlobalNamespace_IsQnameDeclarationWithoutAssignment(node) {
			this_.firstQnameDeclarationWithoutAssignmentJsDocInfo = node.GetJSDocInfo()
		}
		this_.totalGets = this_.totalGets + 1
	case int32(5) /* Type.ALIASING_GET */ :
		this_.aliasingGets = this_.aliasingGets + 1
		this_.totalGets = this_.totalGets + 1
	case int32(7) /* Type.CALL_GET */ :
		this_.callGets = this_.callGets + 1
		this_.totalGets = this_.totalGets + 1
	case int32(8) /* Type.DELETE_PROP */ :
		this_.deleteProps = this_.deleteProps + 1
	case int32(9) /* Type.SUBCLASSING_GET */ :
		this_.subclassingGets = this_.subclassingGets + 1
		this_.totalGets = this_.totalGets + 1
	}
}

// Java method: Name.updateRefNode(com.google.javascript.jscomp.GlobalNamespace.Ref ref, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node newNode)
func (this_ *GlobalNamespace_Name) UpdateRefNode(ref *GlobalNamespace_Ref, newNode pkg_rhino.Node) {
	_ = ref
	_ = newNode
	if !(ref.node != newNode) {
		j2gort.CheckArgumentFmt(false, "redundant update to Ref node: %s", j2gort.PtrToAny(ref))
	}
	var oldNode pkg_rhino.Node = ref.GetNode()
	_ = oldNode
	if !(oldNode != 0) {
		j2gort.CheckStateFmt(false, "Ref's node is already null: %s", j2gort.PtrToAny(ref))
	}
	{
		r_33 := ref
		r_33.node = newNode
	}
	if this_.refsForNodeMap != nil {
		this_.refsForNodeMap.Remove(oldNode)
		if newNode != 0 {
			var existingRefForNewNode *GlobalNamespace_Ref = this_.refsForNodeMap.Get(newNode)
			_ = existingRefForNewNode
			if !(existingRefForNewNode == nil) {
				j2gort.CheckArgumentFmt(false, "refs already exist: %s", j2gort.PtrToAny(existingRefForNewNode))
			}
			this_.refsForNodeMap.Put(newNode, ref)
		}
	} else if !(this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_MULTI_REF_MODE())) {
		if this_.ref0 == nil {
			this_.ref0 = ref
		} else {
			if !(this_.ref0 == ref) {
				j2gort.CheckState(false)
			}
		}
	} else {
		if (this_.ref0 == ref) || ((this_.ref0 != nil) && (this_.ref0.node == oldNode)) {
			this_.ref0 = this_.ref1
			this_.ref1 = this_.ref2
			this_.ref2 = this_.ref3
			this_.ref3 = nil
		} else if (this_.ref1 == ref) || ((this_.ref1 != nil) && (this_.ref1.node == oldNode)) {
			this_.ref1 = this_.ref2
			this_.ref2 = this_.ref3
			this_.ref3 = nil
		} else if (this_.ref2 == ref) || ((this_.ref2 != nil) && (this_.ref2.node == oldNode)) {
			this_.ref2 = this_.ref3
			this_.ref3 = nil
		} else if (this_.ref3 == ref) || ((this_.ref3 != nil) && (this_.ref3.node == oldNode)) {
			this_.ref3 = nil
		}
		if newNode != 0 {
			var existingRefForNewNode_2 *GlobalNamespace_Ref = this_.priv_Name_findInlineRefForNode_3120(newNode)
			_ = existingRefForNewNode_2
			if !(existingRefForNewNode_2 == nil) {
				j2gort.CheckArgumentFmt(false, "refs already exist: %s", j2gort.PtrToAny(existingRefForNewNode_2))
			}
			this_.priv_Name_putRef_3119(newNode, ref)
		}
	}
}

// Java method: Name.removeRef(com.google.javascript.jscomp.GlobalNamespace.Ref ref)
func (this_ *GlobalNamespace_Name) RemoveRef(ref *GlobalNamespace_Ref) {
	_ = ref
	if this_.refsForNodeMap != nil {
		if !(this_.refsForNodeMap.ContainsKey(ref.GetNode())) {
			j2gort.CheckStateFmt(false, "removeRef(%s): unknown ref", j2gort.PtrToAny(ref))
		}
		var refNode pkg_rhino.Node = ref.GetNode()
		_ = refNode
		if refNode != 0 {
			var existing *GlobalNamespace_Ref = this_.refsForNodeMap.Get(refNode)
			_ = existing
			if !(existing == ref) {
				j2gort.CheckStateFmt(false, "Unexpected Refs for Node: %s: when removing Ref: %s", j2gort.PtrToAny(existing), j2gort.PtrToAny(ref))
			}
			this_.refsForNodeMap.Remove(refNode)
		}
	} else if !(this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_MULTI_REF_MODE())) {
		if this_.ref0 == nil {
			panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "removeRef("+j2gort.StrVal(ref)+"): unknown ref")))
		}
		if !(this_.ref0 == ref) {
			j2gort.CheckStateFmt(false, "removeRef(%s): unknown ref", j2gort.PtrToAny(ref))
		}
		this_.ref0 = nil
	} else {
		var refNode_2 pkg_rhino.Node = ref.GetNode()
		_ = refNode_2
		var existing_2 *GlobalNamespace_Ref
		_ = existing_2
		if refNode_2 == 0 {
			existing_2 = nil
		} else {
			existing_2 = this_.priv_Name_findInlineRefForNode_3120(refNode_2)
		}
		if !(existing_2 != nil) {
			j2gort.CheckStateFmt(false, "removeRef(%s): unknown ref", j2gort.PtrToAny(ref))
		}
		if !(existing_2 == ref) {
			j2gort.CheckStateFmt(false, "Unexpected Refs for Node: %s: when removing Ref: %s", j2gort.PtrToAny(existing_2), j2gort.PtrToAny(ref))
		}
		if existing_2 == this_.ref0 {
			this_.ref0 = this_.ref1
			this_.ref1 = this_.ref2
			this_.ref2 = this_.ref3
			this_.ref3 = nil
		} else if existing_2 == this_.ref1 {
			this_.ref1 = this_.ref2
			this_.ref2 = this_.ref3
			this_.ref3 = nil
		} else if existing_2 == this_.ref2 {
			this_.ref2 = this_.ref3
			this_.ref3 = nil
		} else {
			this_.ref3 = nil
		}
	}
	this_.priv_Name_removeRefAndUpdateState_3126(ref)
}

// Java method: Name.removeRefAndUpdateState(com.google.javascript.jscomp.GlobalNamespace.Ref ref)
func (this_ *GlobalNamespace_Name) priv_Name_removeRefAndUpdateState_3126(ref *GlobalNamespace_Ref) {
	_ = ref
	if ref == this_.declaration {
		this_.declaration = nil
		if this_.refsForNodeMap != nil {
			{
				cur_34 := this_.refsForNodeMap.IterValues()
				for cur_34.HasNext() {

					maybeNewDecl := cur_34.
						Next()
					_ = maybeNewDecl
					if maybeNewDecl.type_ == StaticGet_GlobalNamespace_Ref_Type_SET_FROM_GLOBAL() {
						this_.declaration = maybeNewDecl
						break
					}
				}
			}
		} else if this_.ref0 != nil {
			if this_.ref0.type_ == StaticGet_GlobalNamespace_Ref_Type_SET_FROM_GLOBAL() {
				this_.declaration = this_.ref0
			} else if this_.ref1 != nil {
				if this_.ref1.type_ == StaticGet_GlobalNamespace_Ref_Type_SET_FROM_GLOBAL() {
					this_.declaration = this_.ref1
				} else if this_.ref2 != nil {
					if this_.ref2.type_ == StaticGet_GlobalNamespace_Ref_Type_SET_FROM_GLOBAL() {
						this_.declaration = this_.ref2
					} else if (this_.ref3 != nil) && (this_.ref3.type_ == StaticGet_GlobalNamespace_Ref_Type_SET_FROM_GLOBAL()) {
						this_.declaration = this_.ref3
					}
				}
			}
		}
	}
	var info *pkg_rhino.JSDocInfo = nil
	_ = info
	switch (ref.type_).Ordinal() {
	case int32(0) /* Type.SET_FROM_GLOBAL */ :
		this_.globalSets = this_.globalSets - 1
	case int32(2) /* Type.GET_AND_SET_FROM_GLOBAL */ :
		this_.aliasingGets = this_.aliasingGets - 1
		this_.totalGets = this_.totalGets - 1
		this_.globalSets = this_.globalSets - 1
	case int32(1) /* Type.SET_FROM_LOCAL */, int32(3) /* Type.GET_AND_SET_FROM_LOCAL */ :
		this_.localSets = this_.localSets - 1
		if ref.GetNode() == 0 {
			info = nil
		} else {
			info = NodeUtil_GetBestJSDocInfo(ref.GetNode())
		}
		if (info != nil) && info.IsNoCollapse() {
			this_.localSetsWithNoCollapse = this_.localSetsWithNoCollapse - 1
		}
		{
			var iifeCond_1 bool
			{

				re_36 := ref.type_
				ae_37 := StaticGet_GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL()
				if re_36 == 0 {
					j2gort.ThrowNPE()
				}
				iifeCond_1 = re_36 == ae_37
			}
			if iifeCond_1 {
				this_.aliasingGets = this_.aliasingGets - 1
				this_.totalGets = this_.totalGets - 1
			}
		}
	case int32(4) /* Type.PROTOTYPE_GET */, int32(6) /* Type.DIRECT_GET */ :
		this_.totalGets = this_.totalGets - 1
	case int32(5) /* Type.ALIASING_GET */ :
		this_.aliasingGets = this_.aliasingGets - 1
		this_.totalGets = this_.totalGets - 1
	case int32(7) /* Type.CALL_GET */ :
		this_.callGets = this_.callGets - 1
		this_.totalGets = this_.totalGets - 1
	case int32(8) /* Type.DELETE_PROP */ :
		this_.deleteProps = this_.deleteProps - 1
	case int32(9) /* Type.SUBCLASSING_GET */ :
		this_.subclassingGets = this_.subclassingGets - 1
		this_.totalGets = this_.totalGets - 1
	}
}

// Java method: Name.hasInlineRefs()
func (this_ *GlobalNamespace_Name) HasInlineRefs() bool {
	return (this_.refsForNodeMap == nil)
}

// Java method: Name.inlineRef0()
func (this_ *GlobalNamespace_Name) InlineRef0() *GlobalNamespace_Ref {
	return this_.ref0
}

// Java method: Name.inlineRef1()
func (this_ *GlobalNamespace_Name) InlineRef1() *GlobalNamespace_Ref {
	return this_.ref1
}

// Java method: Name.inlineRef2()
func (this_ *GlobalNamespace_Name) InlineRef2() *GlobalNamespace_Ref {
	return this_.ref2
}

// Java method: Name.inlineRef3()
func (this_ *GlobalNamespace_Name) InlineRef3() *GlobalNamespace_Ref {
	return this_.ref3
}

// Java method: Name.getOnlyRef()
func (this_ *GlobalNamespace_Name) GetOnlyRef() *GlobalNamespace_Ref {
	if this_.refsForNodeMap == nil {
		if this_.ref0 == nil {
			panic(j2gort.AsIface[javalang.INoSuchElementException](j2gort.ConstructShimThrowable("java.util.NoSuchElementException")))
		}
		if this_.ref1 == nil {
			return this_.ref0
		}
		panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "expected one element")))
	}
	var size int32 = this_.refsForNodeMap.Size()
	_ = size
	if size == int32(0) {
		panic(j2gort.AsIface[javalang.INoSuchElementException](j2gort.ConstructShimThrowable("java.util.NoSuchElementException")))
	}
	if size == int32(1) {
		{
			cur_38 := this_.refsForNodeMap.IterValues()
			for cur_38.HasNext() {

				ref := cur_38.
					Next()
				_ = ref
				return ref
			}
		}
	}
	panic(j2gort.AsIface[javalang.IIllegalArgumentException](j2gort.ConstructShimThrowable("java.lang.IllegalArgumentException", "expected one element")))
}

// Java method: Name.getRefs()
func (this_ *GlobalNamespace_Name) GetRefs() j2gort.Ext_java_util_Collection {
	if this_.refsForNodeMap != nil {
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](this_.refsForNodeMap.LiveValues())
	}
	if this_.ref0 != nil {
		if this_.ref1 == nil {
			return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.ImmListOfPtr(this_.ref0))
		}
		if this_.ref2 == nil {
			return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_2(j2gort.PtrToAny(this_.ref0), j2gort.PtrToAny(this_.ref1)))
		}
		if this_.ref3 == nil {
			return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_3(j2gort.PtrToAny(this_.ref0), j2gort.PtrToAny(this_.ref1), j2gort.PtrToAny(this_.ref2)))
		}
		return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_4(j2gort.PtrToAny(this_.ref0), j2gort.PtrToAny(this_.ref1), j2gort.PtrToAny(this_.ref2), j2gort.PtrToAny(this_.ref3)))
	}
	return j2gort.AsIface[j2gort.Ext_java_util_Collection](j2gort.Ext_com_google_common_collect_ImmutableList_Of_0())
}

// Java method: Name.getRefForNode(com.google.javascript.rhino.Node node)
func (this_ *GlobalNamespace_Name) GetRefForNode(node pkg_rhino.Node) *GlobalNamespace_Ref {
	_ = node
	j2gort.CheckNotNullNode(node)
	if this_.refsForNodeMap != nil {
		return this_.refsForNodeMap.Get(node)
	}
	return this_.priv_Name_findInlineRefForNode_3120(node)
}

// Java method: Name.getFirstRef()
func (this_ *GlobalNamespace_Name) GetFirstRef() *GlobalNamespace_Ref {
	if this_.ref0 != nil {
		return this_.ref0
	}
	if this_.refsForNodeMap != nil {
		{
			cur_40 := this_.refsForNodeMap.IterValues()
			for cur_40.HasNext() {

				ref := cur_40.
					Next()
				_ = ref
				return ref
			}
		}
		panic(j2gort.AsIface[javalang.INoSuchElementException](j2gort.ConstructShimThrowable("java.util.NoSuchElementException")))
	}
	if this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_MULTI_REF_MODE()) {
		panic(j2gort.AsIface[javalang.INoSuchElementException](j2gort.ConstructShimThrowable("java.util.NoSuchElementException")))
	}
	panic(j2gort.AsIface[javalang.INullPointerException](j2gort.ConstructShimThrowable("java.lang.NullPointerException", "no first Ref to get")))
}

// Java method: Name.canEliminate()
func (this_ *GlobalNamespace_Name) CanEliminate() bool {
	if !(this_.CanCollapseUnannotatedChildNames()) || (this_.totalGets > int32(0)) {
		return false
	}
	{
		var i int32 = int32(0)
		_ = i
		var size int32 = this_.GetPropCount()
		_ = size
		for ; i < size; i = i + 1 {
			if !(this_.GetPropAt(i).CanCollapse()) {
				return false
			}
		}
	}
	return true
}

// Java method: Name.isSimpleStubDeclaration()
func (this_ *GlobalNamespace_Name) IsSimpleStubDeclaration() bool {
	if ((this_.ref0 != nil) && (this_.ref1 == nil)) || ((this_.refsForNodeMap != nil) && (this_.refsForNodeMap.Size() == int32(1))) {
		var ref *GlobalNamespace_Ref = this_.GetFirstRef()
		_ = ref
		if ref.node.GetParent().IsExprResult() {
			return true
		}
	}
	return false
}

// Java method: Name.isCollapsingExplicitlyDenied()
func (this_ *GlobalNamespace_Name) IsCollapsingExplicitlyDenied() bool {
	var docInfo *pkg_rhino.JSDocInfo = this_.GetJSDocInfo()
	_ = docInfo
	return ((docInfo != nil) && docInfo.IsNoCollapse())
}

// Java method: Name.calculateInlinability()
func (this_ *GlobalNamespace_Name) CalculateInlinability() GlobalNamespace_Inlinability {
	if (this_.InExterns() || !(this_.priv_Name_hasOneRealGlobalSet_3141())) || (this_.localSets != int32(0)) {
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	}
	var inlinability GlobalNamespace_Inlinability = this_.CanCollapseOrInline()
	_ = inlinability
	if !(inlinability.ShouldInlineUsages()) {
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	}
	var initialChunk *JSChunk
	_ = initialChunk
	if this_.GetInitialization() != nil {
		initialChunk = this_.GetInitialization().GetChunk()
	} else {
		initialChunk = j2gort.CastPtr[GlobalNamespace_Ref](this_.GetDeclaration()).GetChunk()
	}
	if this_.ref0 != nil {
		inlinability = this_.priv_Name_checkRefInlinability_3140(this_.ref0, initialChunk, inlinability)
		if inlinability == StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE() {
			return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
		}
		if this_.ref1 != nil {
			inlinability = this_.priv_Name_checkRefInlinability_3140(this_.ref1, initialChunk, inlinability)
			if inlinability == StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE() {
				return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
			}
			if this_.ref2 != nil {
				inlinability = this_.priv_Name_checkRefInlinability_3140(this_.ref2, initialChunk, inlinability)
				if inlinability == StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE() {
					return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
				}
				if this_.ref3 != nil {
					inlinability = this_.priv_Name_checkRefInlinability_3140(this_.ref3, initialChunk, inlinability)
				}
			}
		}
	} else if this_.refsForNodeMap != nil {
		{
			cur_42 := this_.refsForNodeMap.IterValues()
			for cur_42.HasNext() {

				ref := cur_42.
					Next()
				_ = ref
				inlinability = this_.priv_Name_checkRefInlinability_3140(ref, initialChunk, inlinability)
				if inlinability == StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE() {
					return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
				}
			}
		}
	}
	return inlinability
}

// Java method: Name.checkRefInlinability(com.google.javascript.jscomp.GlobalNamespace.Ref ref, com.google.javascript.jscomp.JSChunk initialChunk, com.google.javascript.jscomp.GlobalNamespace.Inlinability inlinability)
func (this_ *GlobalNamespace_Name) priv_Name_checkRefInlinability_3140(ref *GlobalNamespace_Ref, initialChunk *JSChunk, inlinability GlobalNamespace_Inlinability) GlobalNamespace_Inlinability {
	_ = ref
	_ = initialChunk
	_ = inlinability
	switch (ref.type_).Ordinal() {
	case int32(0) /* Type.SET_FROM_GLOBAL */, int32(2) /* Type.GET_AND_SET_FROM_GLOBAL */ :
		if !(ref.IsUninitializedDeclaration() || (j2gort.CheckNotNullPtr(ref.GetChunk()) == initialChunk)) {
			j2gort.CheckStateMsg(false, j2gort.PtrToAny(ref))
		}
		return inlinability
	case int32(1) /* Type.SET_FROM_LOCAL */, int32(3) /* Type.GET_AND_SET_FROM_LOCAL */ :
		panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException")))
	case int32(5) /* Type.ALIASING_GET */, int32(6) /* Type.DIRECT_GET */, int32(4) /* Type.PROTOTYPE_GET */, int32(7) /* Type.CALL_GET */, int32(9) /* Type.SUBCLASSING_GET */ :
		var refChunk *JSChunk = ref.GetChunk()
		_ = refChunk
		if ((inlinability != StaticGet_GlobalNamespace_Inlinability_INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY()) && (refChunk != initialChunk)) && !(this_.outer_.compiler.GetChunkGraph().DependsOn(refChunk, initialChunk)) {
			return StaticGet_GlobalNamespace_Inlinability_INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY()
		}
		return inlinability
	case int32(8) /* Type.DELETE_PROP */ :
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	}
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError")))
}

// Java method: Name.hasOneRealGlobalSet()
func (this_ *GlobalNamespace_Name) priv_Name_hasOneRealGlobalSet_3141() bool {
	return ((this_.globalSets == int32(1)) || ((this_.globalSets == int32(2)) && this_.declaration.IsUninitializedDeclaration()))
}

// Java method: Name.canCollapse()
func (this_ *GlobalNamespace_Name) CanCollapse() bool {
	return this_.CanCollapseOrInline().CanCollapse()
}

// Java method: Name.canCollapseOrInline()
func (this_ *GlobalNamespace_Name) CanCollapseOrInline() GlobalNamespace_Inlinability {
	if this_.InExterns() {
		this_.priv_Name_logDecision_3144(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), "declared in externs")
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	}
	if this_.IsGetOrSetDefinition() {
		this_.priv_Name_logDecision_3144(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), "getter / setter")
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	}
	if this_.IsCollapsingExplicitlyDenied() {
		this_.priv_Name_logDecision_3144(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), "@nocollapse")
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	}
	if this_.ReferencesSuperOrInnerClassName() {
		this_.priv_Name_logDecision_3144(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), "references super or inner class name")
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	}
	if this_.IsToStringValueOfInObjectLiteral() {
		this_.priv_Name_logDecision_3144(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), "references explicit definition of toString/valueOf functions used implicitly in the JS language")
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	}
	if this_.deleteProps > int32(0) {
		this_.priv_Name_logDecision_3144(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), "delete operator is used on this property")
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	}
	if j2gort.CastPtr[GlobalNamespace_Ref](this_.GetDeclaration()) != nil {
		var declaration pkg_rhino.Node = j2gort.CastPtr[GlobalNamespace_Ref](this_.GetDeclaration()).GetNode()
		_ = declaration
		if declaration.GetParent().IsObjectLit() {
			if this_.outer_.priv_GlobalNamespace_declarationHasFollowingObjectSpreadSibling_3062(declaration) {
				this_.priv_Name_logDecision_3144(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), "obj lit property followed by spread")
				return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
			}
			{
				var child pkg_rhino.Node = this_.outer_.priv_GlobalNamespace_getOutermostObjectLit_3064(declaration)
				_ = child
				var p pkg_rhino.Node = child.GetParent()
				_ = p
				for ; ((p != 0) && (p.GetParent() != 0)) && !(NodeUtil_IsStatement(p)); func() { child = p; p = p.GetParent() }() {
					if ((p.IsOr() && (p.GetFirstChild() != child)) || (p.IsHook() && (p.GetFirstChild() != child))) || (p.IsAnd() && (p.GetFirstChild() != child)) {
						this_.priv_Name_logDecision_3144(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), "conditional definition")
						return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
					}
				}
			}
		}
	}
	var isUnchangedThroughFullName bool = ((((this_.globalSets > int32(0)) || (this_.localSets > int32(0))) && (this_.localSetsWithNoCollapse == int32(0))) && (this_.deleteProps == int32(0)))
	_ = isUnchangedThroughFullName
	var parentInlinability GlobalNamespace_Inlinability
	_ = parentInlinability
	if this_.parent == nil {
		parentInlinability = StaticGet_GlobalNamespace_Inlinability_INLINE_COMPLETELY()
	} else {
		parentInlinability = this_.parent.CanCollapseOrInlineChildNames()
	}
	switch (parentInlinability).Ordinal() {
	case int32(0) /* Inlinability.INLINE_COMPLETELY */ :
		if isUnchangedThroughFullName {
			this_.priv_Name_logDecision_3144(StaticGet_GlobalNamespace_Inlinability_INLINE_COMPLETELY(), "parent inlineable: unchanged through full name")
			return StaticGet_GlobalNamespace_Inlinability_INLINE_COMPLETELY()
		} else {
			var unsafeInlinablility GlobalNamespace_Inlinability = this_.priv_Name_getUnsafeInlinabilityBasedOnDeclaredType_3145()
			_ = unsafeInlinablility
			this_.priv_Name_logDecision_3144(unsafeInlinablility, "parent inlineable: changed through full name")
			return unsafeInlinablility
		}
	case int32(1) /* Inlinability.INLINE_BUT_KEEP_DECLARATION */, int32(2) /* Inlinability.INLINE_UNLESS_INVALID_CROSS_CHUNK_DEPENDENCY */ :
		if this_.IsDeclaredType() {
			var unsafeInlinability GlobalNamespace_Inlinability = this_.priv_Name_getUnsafeInlinabilityBasedOnDeclaredType_3145()
			_ = unsafeInlinability
			this_.priv_Name_logDecision_3144(unsafeInlinability, "parent unsafely inlineable & is declared type")
			return unsafeInlinability
		} else if isUnchangedThroughFullName {
			this_.priv_Name_logDecision_3144(parentInlinability, "parent unsafely inlineable & unchanged through full name")
			return parentInlinability
		} else {
			this_.priv_Name_logDecision_3144(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), "parent unsafely inlineable & changed through full name")
			return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
		}
	case int32(3) /* Inlinability.DO_NOT_INLINE */ :
		var unsafeInlinability_2 GlobalNamespace_Inlinability = this_.priv_Name_getUnsafeInlinabilityBasedOnDeclaredType_3145()
		_ = unsafeInlinability_2
		this_.priv_Name_logDecision_3144(unsafeInlinability_2, "parent cannot be inlined")
		return unsafeInlinability_2
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", "unknown enum value "+j2gort.StrVal(parentInlinability))))
}

// Java method: Name.logDecision(com.google.javascript.jscomp.GlobalNamespace.Inlinability inlinability, java.lang.String reason)
func (this_ *GlobalNamespace_Name) priv_Name_logDecision_3144(inlinability GlobalNamespace_Inlinability, reason string) {
	_ = inlinability
	_ = reason
	if (this_.outer_.decisionsLog != nil && !this_.outer_.decisionsLog.IsJavaNil_()) && this_.outer_.decisionsLog.IsLogging() {
		this_.outer_.decisionsLog.Log_2("%s: %s: %s", []any{j2gort.BoxString(this_.GetFullName()), j2gort.BoxNullableVal(inlinability), j2gort.BoxString(reason)})
	}
}

// Java method: Name.getUnsafeInlinabilityBasedOnDeclaredType()
func (this_ *GlobalNamespace_Name) priv_Name_getUnsafeInlinabilityBasedOnDeclaredType_3145() GlobalNamespace_Inlinability {
	if this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_CONSTRUCTOR_TYPE()) {
		return StaticGet_GlobalNamespace_Inlinability_INLINE_BUT_KEEP_DECLARATION()
	} else if this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_INTERFACE_TYPE()) {
		return StaticGet_GlobalNamespace_Inlinability_INLINE_BUT_KEEP_DECLARATION()
	} else if this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_ENUM_TYPE()) {
		return StaticGet_GlobalNamespace_Inlinability_INLINE_BUT_KEEP_DECLARATION()
	} else if this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_NOT_A_TYPE()) {
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	}
	panic(j2gort.AsIface[javalang.IIllegalStateException](j2gort.ConstructShimThrowable("java.lang.IllegalStateException", j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("name missing declaredType value: %s", []any{this_})))))
}

// Java method: Name.referencesSuperOrInnerClassName()
func (this_ *GlobalNamespace_Name) ReferencesSuperOrInnerClassName() bool {
	var ref *GlobalNamespace_Ref = j2gort.CastPtr[GlobalNamespace_Ref](this_.GetDeclaration())
	_ = ref
	if ref == nil {
		return false
	}
	var member pkg_rhino.Node = ref.GetNode()
	_ = member
	if (member == 0) || !(member.IsStaticMember() && member.GetParent().IsClassMembers()) {
		return false
	}
	var body pkg_rhino.Node
	_ = body
	if member.IsMemberFieldDef() {
		body = member.GetFirstChild()
	} else {
		body = NodeUtil_GetFunctionBody(member.GetFirstChild())
	}
	if (body != 0) && NodeUtil_ReferencesSuper(body) {
		return true
	}
	var classNode pkg_rhino.Node = member.GetGrandparent()
	_ = classNode
	if NodeUtil_IsClassDeclaration(classNode) {
		return false
	}
	var innerNameNode pkg_rhino.Node = classNode.GetFirstChild()
	_ = innerNameNode
	return (!(innerNameNode.IsEmpty()) && NodeUtil_IsNameReferenced_2(member, innerNameNode.GetString()))
}

// Java method: Name.isSetInLoop()
func (this_ *GlobalNamespace_Name) priv_Name_isSetInLoop_3147() bool {
	var ref *GlobalNamespace_Ref = j2gort.CastPtr[GlobalNamespace_Ref](this_.GetDeclaration())
	_ = ref
	if ref != nil {
		var n pkg_rhino.Node = ref.GetNode()
		_ = n
		if n != 0 {
			return NodeUtil_IsWithinLoop(n)
		}
	}
	return false
}

// Java method: Name.isGetOrSetDefinition()
func (this_ *GlobalNamespace_Name) IsGetOrSetDefinition() bool {
	return this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_GET_SET())
}

// Java method: Name.canCollapseUnannotatedChildNames()
func (this_ *GlobalNamespace_Name) CanCollapseUnannotatedChildNames() bool {
	return this_.CanCollapseOrInlineChildNames().CanCollapse()
}

// Java method: Name.isToStringValueOfInObjectLiteral()
func (this_ *GlobalNamespace_Name) IsToStringValueOfInObjectLiteral() bool {
	var parent *GlobalNamespace_Name = this_.GetParent()
	_ = parent
	var baseName string = this_.GetBaseName()
	_ = baseName
	return (((this_.IsFunction() && (parent != nil)) && parent.IsObjectLiteral()) && ((j2gort.CheckNotNullString(baseName) == "toString") || (j2gort.CheckNotNullString(baseName) == "valueOf")))
}

// Java method: Name.canCollapseOrInlineChildNames()
func (this_ *GlobalNamespace_Name) CanCollapseOrInlineChildNames() GlobalNamespace_Inlinability {
	if this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_OTHER_OBJECT()) {
		this_.priv_Name_logChildNamesDecision_3152(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), "NameProp.OTHER_OBJECT")
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	} else if this_.IsGetOrSetDefinition() {
		this_.priv_Name_logChildNamesDecision_3152(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), "getter/setter")
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	} else if this_.globalSets != int32(1) {
		if (this_.outer_.decisionsLog != nil && !this_.outer_.decisionsLog.IsJavaNil_()) && this_.outer_.decisionsLog.IsLogging() {
			this_.priv_Name_logChildNamesDecision_3153(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
				return j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("set %d times globally", []any{j2gort.ValueOfInt(this_.globalSets)}))
			}))
		}
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	} else if this_.localSets != int32(0) {
		if (this_.outer_.decisionsLog != nil && !this_.outer_.decisionsLog.IsJavaNil_()) && this_.outer_.decisionsLog.IsLogging() {
			this_.priv_Name_logChildNamesDecision_3153(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
				return j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("set %d times locally", []any{j2gort.ValueOfInt(this_.localSets)}))
			}))
		}
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	} else if this_.deleteProps != int32(0) {
		if (this_.outer_.decisionsLog != nil && !this_.outer_.decisionsLog.IsJavaNil_()) && this_.outer_.decisionsLog.IsLogging() {
			this_.priv_Name_logChildNamesDecision_3153(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
				return j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("properties are deleted %d times", []any{j2gort.ValueOfInt(this_.deleteProps)}))
			}))
		}
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	}
	j2gort.CheckNotNullPtr(this_.declaration)
	if this_.declaration.IsTwin() {
		this_.priv_Name_logChildNamesDecision_3152(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), "twinned declaration")
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	}
	if this_.IsCollapsingExplicitlyDenied() {
		this_.priv_Name_logChildNamesDecision_3152(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), "@nocollapse")
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	}
	if this_.priv_Name_isSetInLoop_3147() {
		this_.priv_Name_logChildNamesDecision_3152(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), "set in a loop")
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	}
	if this_.UsesHasOwnProperty() {
		this_.priv_Name_logChildNamesDecision_3152(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), "hasOwnProperty() call exists")
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	}
	if this_.priv_Name_valueImplicitlySupportsAliasing_3154() {
		this_.priv_Name_logChildNamesDecision_3152(StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE(), "value implicitly supports aliasing")
		return StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE()
	}
	if (this_.parent != nil) && this_.parent.ShouldKeepKeys() {
		var unsafeInlinability GlobalNamespace_Inlinability = this_.priv_Name_getUnsafeInlinabilityBasedOnDeclaredType_3145()
		_ = unsafeInlinability
		this_.priv_Name_logChildNamesDecision_3152(unsafeInlinability, "parent.shouldKeepKeys()")
		return unsafeInlinability
	}
	if this_.aliasingGets > int32(0) {
		var unsafeInlinability_2 GlobalNamespace_Inlinability = this_.priv_Name_getUnsafeInlinabilityBasedOnDeclaredType_3145()
		_ = unsafeInlinability_2
		if (this_.outer_.decisionsLog != nil && !this_.outer_.decisionsLog.IsJavaNil_()) && this_.outer_.decisionsLog.IsLogging() {
			this_.priv_Name_logChildNamesDecision_3153(unsafeInlinability_2, j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
				return j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("%d aliasing gets exist", []any{j2gort.ValueOfInt(this_.aliasingGets)}))
			}))
		}
		return unsafeInlinability_2
	}
	if this_.parent == nil {
		this_.priv_Name_logChildNamesDecision_3152(StaticGet_GlobalNamespace_Inlinability_INLINE_COMPLETELY(), "no reason not to inline")
		return StaticGet_GlobalNamespace_Inlinability_INLINE_COMPLETELY()
	}
	var parentInlinability GlobalNamespace_Inlinability = this_.parent.CanCollapseOrInlineChildNames()
	_ = parentInlinability
	if parentInlinability == StaticGet_GlobalNamespace_Inlinability_DO_NOT_INLINE() {
		var unsafeInlinability_3 GlobalNamespace_Inlinability = this_.priv_Name_getUnsafeInlinabilityBasedOnDeclaredType_3145()
		_ = unsafeInlinability_3
		this_.priv_Name_logChildNamesDecision_3152(unsafeInlinability_3, "parent is not inlineable")
		return unsafeInlinability_3
	}
	this_.priv_Name_logChildNamesDecision_3152(parentInlinability, "inherited from parent")
	return parentInlinability
}

// Java method: Name.logChildNamesDecision(com.google.javascript.jscomp.GlobalNamespace.Inlinability inlinability, java.lang.String reason)
func (this_ *GlobalNamespace_Name) priv_Name_logChildNamesDecision_3152(inlinability GlobalNamespace_Inlinability, reason string) {
	_ = inlinability
	_ = reason
	if (this_.outer_.decisionsLog != nil && !this_.outer_.decisionsLog.IsJavaNil_()) && this_.outer_.decisionsLog.IsLogging() {
		this_.outer_.decisionsLog.Log_2("%s: children: %s: %s", []any{j2gort.BoxString(this_.GetFullName()), j2gort.BoxNullableVal(inlinability), j2gort.BoxString(reason)})
	}
}

// Java method: Name.logChildNamesDecision(com.google.javascript.jscomp.GlobalNamespace.Inlinability inlinability, java.util.function.Supplier<java.lang.String> reasonSupplier)
func (this_ *GlobalNamespace_Name) priv_Name_logChildNamesDecision_3153(inlinability GlobalNamespace_Inlinability, reasonSupplier j2gort.Ext_java_util_function_Supplier) {
	_ = inlinability
	_ = reasonSupplier
	if (this_.outer_.decisionsLog != nil && !this_.outer_.decisionsLog.IsJavaNil_()) && this_.outer_.decisionsLog.IsLogging() {
		this_.outer_.decisionsLog.Log_1_Ext_java_util_function_Supplier(j2gort.NewLambda_Ext_java_util_function_Supplier(func() any {
			return j2gort.BoxString(j2gort.Ext_java_lang_String_Format_2("%s: children: %s: %s", []any{j2gort.BoxString(this_.GetFullName()), j2gort.BoxNullableVal(inlinability), j2gort.BoxString(j2gort.CastString(reasonSupplier.Get()))}))
		}))
	}
}

// Java method: Name.valueImplicitlySupportsAliasing()
func (this_ *GlobalNamespace_Name) priv_Name_valueImplicitlySupportsAliasing_3154() bool {
	if !(this_.outer_.enableImplicitlyAliasedValues) {
		return false
	}
	if this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_CLASS()) {
		return true
	}
	if this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_FUNCTION()) {
		var jsdoc *pkg_rhino.JSDocInfo = this_.GetJSDocInfo()
		_ = jsdoc
		return ((jsdoc != nil) && jsdoc.IsConstructorOrInterface())
	}
	return false
}

// Java method: Name.shouldKeepKeys()
func (this_ *GlobalNamespace_Name) ShouldKeepKeys() bool {
	return (this_.IsObjectLiteral() && ((this_.aliasingGets > int32(0)) || this_.IsCollapsingExplicitlyDenied()))
}

// Java method: Name.needsToBeStubbed()
func (this_ *GlobalNamespace_Name) NeedsToBeStubbed() bool {
	return ((((this_.globalSets == int32(0)) && (this_.localSets > int32(0))) && (this_.localSetsWithNoCollapse == int32(0))) && !(this_.IsCollapsingExplicitlyDenied()))
}

// Java method: Name.setDeclaredTypeKind(com.google.javascript.jscomp.GlobalNamespace.NameProp declaredType)
func (this_ *GlobalNamespace_Name) priv_Name_setDeclaredTypeKind_3157(declaredType *GlobalNamespace_NameProp) {
	_ = declaredType
	if !((declaredType.bit & StaticGet_GlobalNamespace_NameProp_DECLARED_TYPE_KIND_MASK()) != int32(0)) {
		j2gort.CheckArgumentFmt(false, "Unexpected NameProp for declaredType %s", j2gort.PtrToAny(declaredType))
	}
	this_.propertyBitSet = ((this_.propertyBitSet & ^(StaticGet_GlobalNamespace_NameProp_DECLARED_TYPE_KIND_MASK())) | declaredType.bit)
	if declaredType != StaticGet_GlobalNamespace_NameProp_NOT_A_TYPE() {
		{
			var ancestor *GlobalNamespace_Name = this_.parent
			_ = ancestor
			for ; ancestor != nil; ancestor = ancestor.parent {
				ancestor.priv_Name_setBooleanProperty_3114(StaticGet_GlobalNamespace_NameProp_IS_DECLARED())
			}
		}
	}
}

// Java method: Name.setNameType(com.google.javascript.jscomp.GlobalNamespace.NameProp type)
func (this_ *GlobalNamespace_Name) priv_Name_setNameType_3158(type_ *GlobalNamespace_NameProp) {
	_ = type_
	if !((type_.bit & StaticGet_GlobalNamespace_NameProp_NAME_KIND_MASK()) != int32(0)) {
		j2gort.CheckArgumentFmt(false, "Unexpected NameProp for nameType %s", j2gort.PtrToAny(type_))
	}
	this_.propertyBitSet = ((this_.propertyBitSet & ^(StaticGet_GlobalNamespace_NameProp_NAME_KIND_MASK())) | type_.bit)
}

// Java method: Name.isDeclaredType()
func (this_ *GlobalNamespace_Name) IsDeclaredType() bool {
	if !((this_.propertyBitSet & StaticGet_GlobalNamespace_NameProp_DECLARED_TYPE_KIND_MASK()) != int32(0)) {
		j2gort.CheckStateMsg(false, j2gort.ValueOfInt(this_.propertyBitSet))
	}
	return !(this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_NOT_A_TYPE()))
}

// Java method: Name.isConstructor()
func (this_ *GlobalNamespace_Name) IsConstructor() bool {
	return this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_CONSTRUCTOR_TYPE())
}

// Java method: Name.isNamespaceObjectLit()
func (this_ *GlobalNamespace_Name) IsNamespaceObjectLit() bool {
	return (this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_IS_DECLARED()) && this_.IsObjectLiteral())
}

// Java method: Name.isSimpleName()
func (this_ *GlobalNamespace_Name) IsSimpleName() bool {
	return (this_.parent == nil)
}

// Java method: Name.getTypeDebugString()
func (this_ *GlobalNamespace_Name) GetTypeDebugString() string {
	if this_.IsClass() {
		return "CLASS"
	} else if this_.IsFunction() {
		return "FUNCTION"
	} else if this_.IsObjectLiteral() {
		return "OBJECTLIT"
	} else if this_.IsGetOrSetDefinition() {
		return "GET_SET"
	} else if this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_OTHER_OBJECT()) {
		return "OTHER"
	}
	panic(j2gort.AsIface[javalang.IAssertionError](j2gort.ConstructShimThrowable("java.lang.AssertionError", "Missing NameProp for name kind in "+j2gort.StrValInt32(this_.propertyBitSet))))
}

// Java method: Name.toString()
func (this_ *GlobalNamespace_Name) String() string {
	return (j2gort.StrValString(this_.GetFullName()) + " (" + j2gort.StrValString(this_.GetTypeDebugString()) + "): " + j2gort.StrValString(j2gort.Ext_com_google_common_base_Joiner_On_1_string(", ").Join_3("globalSets="+j2gort.StrValInt32(this_.globalSets), "localSets="+j2gort.StrValInt32(this_.localSets), []any{"totalGets=" + j2gort.StrValInt32(this_.totalGets), "aliasingGets=" + j2gort.StrValInt32(this_.aliasingGets), "callGets=" + j2gort.StrValInt32(this_.callGets), "subclassingGets=" + j2gort.StrValInt32(this_.subclassingGets)})))
}

// Java method: Name.getJSDocInfo()
func (this_ *GlobalNamespace_Name) GetJSDocInfo() *pkg_rhino.JSDocInfo {
	if this_.firstDeclarationJSDocInfo != nil {
		return this_.firstDeclarationJSDocInfo
	} else {
		return this_.firstQnameDeclarationWithoutAssignmentJsDocInfo
	}
}

// Java method: Name.getDocInfoForDeclaration(com.google.javascript.jscomp.GlobalNamespace.Ref ref)
func (this_ *GlobalNamespace_Name) priv_Name_getDocInfoForDeclaration_3165(ref *GlobalNamespace_Ref) *pkg_rhino.JSDocInfo {
	_ = ref
	if ref.node != 0 {
		var refParent pkg_rhino.Node = ref.node.GetParent()
		_ = refParent
		if refParent == 0 {
			return nil
		}
		switch refParent.TokenSwitchOrd_() {
		case int32(73) /* Token.FUNCTION */, int32(51) /* Token.ASSIGN */, int32(105) /* Token.CLASS */ :
			return refParent.GetJSDocInfo()
		case int32(84) /* Token.VAR */, int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
			if ref.node == refParent.GetFirstChild() {
				return refParent.GetJSDocInfo()
			} else {
				return ref.node.GetJSDocInfo()
			}
		case int32(47) /* Token.OBJECTLIT */, int32(106) /* Token.CLASS_MEMBERS */ :
			return ref.node.GetJSDocInfo()
		default:
		}
	}
	return nil
}

// Java method: Name.isModuleExport()
func (this_ *GlobalNamespace_Name) IsModuleExport() bool {
	return this_.priv_Name_getBooleanProperty_3115(StaticGet_GlobalNamespace_NameProp_IS_MODULE_PROP())
}

// Java class: com.google.javascript.jscomp.GlobalNamespace.Ref
type GlobalNamespace_Ref struct {
	chunk   *JSChunk       `gson:"chunk"`
	scope   *Scope         `gson:"scope"`
	node    pkg_rhino.Node `gson:"node"`
	idHash_ atomic.Int32
	type_   GlobalNamespace_Ref_Type `gson:"type"`
}

func (this_ *GlobalNamespace_Ref) BindSelf_(s any) {
	_ = s
}
func (this_ *GlobalNamespace_Ref) IsJavaNil_() bool { return this_ == nil }
func (*GlobalNamespace_Ref) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GlobalNamespace$Ref", "Ref"
}
func (this_ *GlobalNamespace_Ref) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*GlobalNamespace_Ref) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GlobalNamespace.Ref", "com.google.javascript.jscomp.GlobalNamespace$Ref", func(v_ any) bool {
		return j2gort.InstanceOf[*GlobalNamespace_Ref](v_)
	})
}
func (*GlobalNamespace_Ref) Implements_rhino_StaticRef() {}

func (this_ *GlobalNamespace_Ref) initFields_GlobalNamespace_Ref() {
	_ = this_
}

// Java constructor: Ref(com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.jscomp.@org.jspecify.annotations.Nullable Scope scope, com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node node, com.google.javascript.jscomp.GlobalNamespace.Ref.Type type)
func NewGlobalNamespace_Ref(chunk *JSChunk, scope *Scope, node pkg_rhino.Node, type_ GlobalNamespace_Ref_Type) *GlobalNamespace_Ref {
	this_ := &GlobalNamespace_Ref{}
	InitGlobalNamespace_Ref(this_, chunk, scope, node, type_)
	return this_
}
func allocPooled_GlobalNamespace_Ref(outer_ *GlobalNamespace, chunk *JSChunk, scope *Scope, node pkg_rhino.Node, type_ GlobalNamespace_Ref_Type) *GlobalNamespace_Ref {
	var this_ *GlobalNamespace_Ref
	if outer_ != nil {
		this_ = outer_.slab_GlobalNamespace_Ref.Alloc()
	} else {
		this_ = &GlobalNamespace_Ref{}
	}
	InitGlobalNamespace_Ref(this_, chunk, scope, node, type_)
	return this_
}

func InitGlobalNamespace_Ref(this_ *GlobalNamespace_Ref, chunk *JSChunk, scope *Scope, node pkg_rhino.Node, type_ GlobalNamespace_Ref_Type) {
	_ = chunk
	_ = scope
	_ = node
	_ = type_
	this_.initFields_GlobalNamespace_Ref()
	this_.chunk = chunk
	this_.node = node
	this_.type_ = type_
	this_.scope = scope
}

// Java method: Ref.getChunk()
func (this_ *GlobalNamespace_Ref) GetChunk() *JSChunk {
	return this_.chunk
}

// Java method: Ref.getNode()
func (this_ *GlobalNamespace_Ref) GetNode() pkg_rhino.Node {
	return this_.node
}

// Java method: Ref.getSourceFile()
func (this_ *GlobalNamespace_Ref) GetSourceFile() pkg_rhino.StaticSourceFile {
	if this_.node != 0 {
		return this_.node.GetStaticSourceFile()
	} else {
		return nil
	}
}

// Java method: Ref.getSymbol()
func (this_ *GlobalNamespace_Ref) GetSymbol() any {
	panic(j2gort.AsIface[javalang.IUnsupportedOperationException](j2gort.ConstructShimThrowable("java.lang.UnsupportedOperationException")))
}

// Java method: Ref.isDeleteProp()
func (this_ *GlobalNamespace_Ref) IsDeleteProp() bool {
	return (this_.type_ == StaticGet_GlobalNamespace_Ref_Type_DELETE_PROP())
}

// Java method: Ref.isSubclassingGet()
func (this_ *GlobalNamespace_Ref) IsSubclassingGet() bool {
	return (this_.type_ == StaticGet_GlobalNamespace_Ref_Type_SUBCLASSING_GET())
}

// Java method: Ref.isTwin()
func (this_ *GlobalNamespace_Ref) IsTwin() bool {
	switch (this_.type_).Ordinal() {
	case int32(2) /* Type.GET_AND_SET_FROM_GLOBAL */, int32(3) /* Type.GET_AND_SET_FROM_LOCAL */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Ref.isGet()
func (this_ *GlobalNamespace_Ref) IsGet() bool {
	switch (this_.type_).Ordinal() {
	case int32(6) /* Type.DIRECT_GET */, int32(5) /* Type.ALIASING_GET */, int32(9) /* Type.SUBCLASSING_GET */, int32(7) /* Type.CALL_GET */, int32(2) /* Type.GET_AND_SET_FROM_GLOBAL */, int32(3) /* Type.GET_AND_SET_FROM_LOCAL */, int32(4) /* Type.PROTOTYPE_GET */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Ref.isAliasingGet()
func (this_ *GlobalNamespace_Ref) IsAliasingGet() bool {
	switch (this_.type_).Ordinal() {
	case int32(5) /* Type.ALIASING_GET */, int32(2) /* Type.GET_AND_SET_FROM_GLOBAL */, int32(3) /* Type.GET_AND_SET_FROM_LOCAL */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Ref.isSet()
func (this_ *GlobalNamespace_Ref) IsSet() bool {
	switch (this_.type_).Ordinal() {
	case int32(0) /* Type.SET_FROM_GLOBAL */, int32(1) /* Type.SET_FROM_LOCAL */, int32(2) /* Type.GET_AND_SET_FROM_GLOBAL */, int32(3) /* Type.GET_AND_SET_FROM_LOCAL */ :
		return true
	default:
		return false
	}
	panic("unreachable switch expression")
}

// Java method: Ref.isSetFromGlobal()
func (this_ *GlobalNamespace_Ref) IsSetFromGlobal() bool {
	return ((this_.type_ == StaticGet_GlobalNamespace_Ref_Type_SET_FROM_GLOBAL()) || (this_.type_ == StaticGet_GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL()))
}

// Java method: Ref.isUninitializedDeclaration()
func (this_ *GlobalNamespace_Ref) IsUninitializedDeclaration() bool {
	return ((this_.node.IsName() && NodeUtil_IsNameDeclaration(this_.node.GetParent())) && !(this_.node.HasChildren()))
}

// Java method: Ref.toString()
func (this_ *GlobalNamespace_Ref) String() string {
	return j2gort.Ext_com_google_common_base_MoreObjects_S_ToStringHelper_1_any(this_).OmitNullValues().Add_2_string_any("type", j2gort.BoxNullableVal(this_.type_)).Add_2_string_any("node", j2gort.BoxNode(this_.node)).Add_2_string_any("scope", j2gort.PtrToAny(this_.scope)).String()
}

// Java enum: com.google.javascript.jscomp.GlobalNamespace.Ref.Type
type GlobalNamespace_Ref_Type uint8

const (
	enum_GlobalNamespace_Ref_Type_SET_FROM_GLOBAL         GlobalNamespace_Ref_Type = 1
	enum_GlobalNamespace_Ref_Type_SET_FROM_LOCAL          GlobalNamespace_Ref_Type = 2
	enum_GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL GlobalNamespace_Ref_Type = 3
	enum_GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL  GlobalNamespace_Ref_Type = 4
	enum_GlobalNamespace_Ref_Type_PROTOTYPE_GET           GlobalNamespace_Ref_Type = 5
	enum_GlobalNamespace_Ref_Type_ALIASING_GET            GlobalNamespace_Ref_Type = 6
	enum_GlobalNamespace_Ref_Type_DIRECT_GET              GlobalNamespace_Ref_Type = 7
	enum_GlobalNamespace_Ref_Type_CALL_GET                GlobalNamespace_Ref_Type = 8
	enum_GlobalNamespace_Ref_Type_DELETE_PROP             GlobalNamespace_Ref_Type = 9
	enum_GlobalNamespace_Ref_Type_SUBCLASSING_GET         GlobalNamespace_Ref_Type = 10
)

var enumBaseTable_GlobalNamespace_Ref_Type = [...]j2gort.EnumBase{
	{Name_: "SET_FROM_GLOBAL", Ordinal_: 0},
	{Name_: "SET_FROM_LOCAL", Ordinal_: 1},
	{Name_: "GET_AND_SET_FROM_GLOBAL", Ordinal_: 2},
	{Name_: "GET_AND_SET_FROM_LOCAL", Ordinal_: 3},
	{Name_: "PROTOTYPE_GET", Ordinal_: 4},
	{Name_: "ALIASING_GET", Ordinal_: 5},
	{Name_: "DIRECT_GET", Ordinal_: 6},
	{Name_: "CALL_GET", Ordinal_: 7},
	{Name_: "DELETE_PROP", Ordinal_: 8},
	{Name_: "SUBCLASSING_GET", Ordinal_: 9},
}
var enumValues_GlobalNamespace_Ref_Type = []GlobalNamespace_Ref_Type{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func StaticGet_GlobalNamespace_Ref_Type_SET_FROM_GLOBAL() GlobalNamespace_Ref_Type {
	return enum_GlobalNamespace_Ref_Type_SET_FROM_GLOBAL
}
func StaticGet_GlobalNamespace_Ref_Type_SET_FROM_LOCAL() GlobalNamespace_Ref_Type {
	return enum_GlobalNamespace_Ref_Type_SET_FROM_LOCAL
}
func StaticGet_GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL() GlobalNamespace_Ref_Type {
	return enum_GlobalNamespace_Ref_Type_GET_AND_SET_FROM_GLOBAL
}
func StaticGet_GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL() GlobalNamespace_Ref_Type {
	return enum_GlobalNamespace_Ref_Type_GET_AND_SET_FROM_LOCAL
}
func StaticGet_GlobalNamespace_Ref_Type_PROTOTYPE_GET() GlobalNamespace_Ref_Type {
	return enum_GlobalNamespace_Ref_Type_PROTOTYPE_GET
}
func StaticGet_GlobalNamespace_Ref_Type_ALIASING_GET() GlobalNamespace_Ref_Type {
	return enum_GlobalNamespace_Ref_Type_ALIASING_GET
}
func StaticGet_GlobalNamespace_Ref_Type_DIRECT_GET() GlobalNamespace_Ref_Type {
	return enum_GlobalNamespace_Ref_Type_DIRECT_GET
}
func StaticGet_GlobalNamespace_Ref_Type_CALL_GET() GlobalNamespace_Ref_Type {
	return enum_GlobalNamespace_Ref_Type_CALL_GET
}
func StaticGet_GlobalNamespace_Ref_Type_DELETE_PROP() GlobalNamespace_Ref_Type {
	return enum_GlobalNamespace_Ref_Type_DELETE_PROP
}
func StaticGet_GlobalNamespace_Ref_Type_SUBCLASSING_GET() GlobalNamespace_Ref_Type {
	return enum_GlobalNamespace_Ref_Type_SUBCLASSING_GET
}

func GlobalNamespace_Ref_Type_Values() []GlobalNamespace_Ref_Type {
	return j2gort.CloneSlice(enumValues_GlobalNamespace_Ref_Type)
}

func GlobalNamespace_Ref_Type_ValueOf(name string) GlobalNamespace_Ref_Type {
	if name == j2gort.NullString {
		j2gort.ThrowNPE()
	}
	for i_ := range enumBaseTable_GlobalNamespace_Ref_Type {
		if enumBaseTable_GlobalNamespace_Ref_Type[i_].Name_ == name {
			return GlobalNamespace_Ref_Type(i_ + 1)
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.GlobalNamespace.Ref.Type." + name))
}

func (this_ GlobalNamespace_Ref_Type) IsJavaNil_() bool                              { return this_ == 0 }
func (GlobalNamespace_Ref_Type) HasIdentityEquals_()                                 {}
func (GlobalNamespace_Ref_Type) Implements_j2gort_Ext_java_lang_Enum()               {}
func (GlobalNamespace_Ref_Type) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ GlobalNamespace_Ref_Type) GetExt_java_lang_EnumBase() *j2gort.Ext_java_lang_EnumBase {
	return &j2gort.Ext_java_lang_EnumBase{Core_: this_.GetEnumBase()}
}
func (GlobalNamespace_Ref_Type) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GlobalNamespace$Ref$Type", "Type"
}
func (this_ GlobalNamespace_Ref_Type) GetEnumBase() *j2gort.EnumBase {
	if this_ == 0 {
		return nil
	}
	return &enumBaseTable_GlobalNamespace_Ref_Type[this_-1]
}
func (this_ GlobalNamespace_Ref_Type) Ordinal() int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return int32(this_ - 1)
}
func (this_ GlobalNamespace_Ref_Type) EnumOrdinal_() int32 { return int32(this_) - 1 }
func (this_ GlobalNamespace_Ref_Type) Name() string {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_GlobalNamespace_Ref_Type[this_-1].Name_
}
func (this_ GlobalNamespace_Ref_Type) String() string {
	if this_ == 0 {
		return "null"
	}
	return enumBaseTable_GlobalNamespace_Ref_Type[this_-1].Name_
}
func (this_ GlobalNamespace_Ref_Type) Equals(other_ any) bool {
	o_, ok_ := other_.(GlobalNamespace_Ref_Type)
	return this_ != 0 && ok_ && this_ == o_
}
func (this_ GlobalNamespace_Ref_Type) HashCode() int32 {
	if this_ == 0 {
		return 0
	}
	return enumBaseTable_GlobalNamespace_Ref_Type[this_-1].HashCode()
}
func (this_ GlobalNamespace_Ref_Type) CompareTo(other_ any) int32 {
	if this_ == 0 {
		j2gort.ThrowNPE()
	}
	return enumBaseTable_GlobalNamespace_Ref_Type[this_-1].CompareTo(other_)
}
func (this_ GlobalNamespace_Ref_Type) WriteGson(out_ *gson.JsonWriter) {
	if this_ == 0 {
		out_.NullValue()
		return
	}
	out_.ValueString(enumBaseTable_GlobalNamespace_Ref_Type[this_-1].Name_)
}
func init() {
	vals_ := []any{GlobalNamespace_Ref_Type(1), GlobalNamespace_Ref_Type(2), GlobalNamespace_Ref_Type(3), GlobalNamespace_Ref_Type(4), GlobalNamespace_Ref_Type(5), GlobalNamespace_Ref_Type(6), GlobalNamespace_Ref_Type(7), GlobalNamespace_Ref_Type(8), GlobalNamespace_Ref_Type(9), GlobalNamespace_Ref_Type(10)}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.GlobalNamespace.Ref.Type", vals_)
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.GlobalNamespace$Ref$Type", vals_)
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GlobalNamespace.Ref.Type", "com.google.javascript.jscomp.GlobalNamespace$Ref$Type", func(v_ any) bool {
		return j2gort.InstanceOf[GlobalNamespace_Ref_Type](v_)
	})
}

// Java enum: com.google.javascript.jscomp.GlobalNamespace.NameProp
type GlobalNamespace_NameProp struct {
	j2gort.EnumBase
	bit int32
}

func (this_ *GlobalNamespace_NameProp) GetEnumBase() *j2gort.EnumBase {
	if this_ == nil {
		return nil
	}
	return &this_.EnumBase
}
func (this_ *GlobalNamespace_NameProp) EnumOrdinal_() int32 {
	if this_ == nil {
		return -1
	}
	return this_.EnumBase.Ordinal_
}
func (*GlobalNamespace_NameProp) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*GlobalNamespace_NameProp) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (this_ *GlobalNamespace_NameProp) initFields_GlobalNamespace_NameProp() {
	_ = this_
}

func initEnumGlobalNamespace_NameProp(this_ *GlobalNamespace_NameProp, enumName_ string, enumOrd_ int32, index int32) *GlobalNamespace_NameProp {
	this_.EnumBase.Name_ = enumName_
	this_.EnumBase.Ordinal_ = enumOrd_
	_ = index
	this_.initFields_GlobalNamespace_NameProp()
	this_.bit = javalang.Shl32(int32(1), int(index))
	return this_
}

var enum_GlobalNamespace_NameProp_IS_DECLARED *GlobalNamespace_NameProp

func StaticGet_GlobalNamespace_NameProp_IS_DECLARED() *GlobalNamespace_NameProp {
	clinit_GlobalNamespace_NameProp()
	return enum_GlobalNamespace_NameProp_IS_DECLARED
}

var enum_GlobalNamespace_NameProp_IS_MODULE_PROP *GlobalNamespace_NameProp

func StaticGet_GlobalNamespace_NameProp_IS_MODULE_PROP() *GlobalNamespace_NameProp {
	clinit_GlobalNamespace_NameProp()
	return enum_GlobalNamespace_NameProp_IS_MODULE_PROP
}

var enum_GlobalNamespace_NameProp_IS_PROVIDED *GlobalNamespace_NameProp

func StaticGet_GlobalNamespace_NameProp_IS_PROVIDED() *GlobalNamespace_NameProp {
	clinit_GlobalNamespace_NameProp()
	return enum_GlobalNamespace_NameProp_IS_PROVIDED
}

var enum_GlobalNamespace_NameProp_IS_USED_HAS_OWN_PROPERTY *GlobalNamespace_NameProp

func StaticGet_GlobalNamespace_NameProp_IS_USED_HAS_OWN_PROPERTY() *GlobalNamespace_NameProp {
	clinit_GlobalNamespace_NameProp()
	return enum_GlobalNamespace_NameProp_IS_USED_HAS_OWN_PROPERTY
}

var enum_GlobalNamespace_NameProp_IS_EXTERN *GlobalNamespace_NameProp

func StaticGet_GlobalNamespace_NameProp_IS_EXTERN() *GlobalNamespace_NameProp {
	clinit_GlobalNamespace_NameProp()
	return enum_GlobalNamespace_NameProp_IS_EXTERN
}

var enum_GlobalNamespace_NameProp_CONSTRUCTOR_TYPE *GlobalNamespace_NameProp

func StaticGet_GlobalNamespace_NameProp_CONSTRUCTOR_TYPE() *GlobalNamespace_NameProp {
	clinit_GlobalNamespace_NameProp()
	return enum_GlobalNamespace_NameProp_CONSTRUCTOR_TYPE
}

var enum_GlobalNamespace_NameProp_INTERFACE_TYPE *GlobalNamespace_NameProp

func StaticGet_GlobalNamespace_NameProp_INTERFACE_TYPE() *GlobalNamespace_NameProp {
	clinit_GlobalNamespace_NameProp()
	return enum_GlobalNamespace_NameProp_INTERFACE_TYPE
}

var enum_GlobalNamespace_NameProp_ENUM_TYPE *GlobalNamespace_NameProp

func StaticGet_GlobalNamespace_NameProp_ENUM_TYPE() *GlobalNamespace_NameProp {
	clinit_GlobalNamespace_NameProp()
	return enum_GlobalNamespace_NameProp_ENUM_TYPE
}

var enum_GlobalNamespace_NameProp_NOT_A_TYPE *GlobalNamespace_NameProp

func StaticGet_GlobalNamespace_NameProp_NOT_A_TYPE() *GlobalNamespace_NameProp {
	clinit_GlobalNamespace_NameProp()
	return enum_GlobalNamespace_NameProp_NOT_A_TYPE
}

var enum_GlobalNamespace_NameProp_CLASS *GlobalNamespace_NameProp

func StaticGet_GlobalNamespace_NameProp_CLASS() *GlobalNamespace_NameProp {
	clinit_GlobalNamespace_NameProp()
	return enum_GlobalNamespace_NameProp_CLASS
}

var enum_GlobalNamespace_NameProp_OBJECTLIT *GlobalNamespace_NameProp

func StaticGet_GlobalNamespace_NameProp_OBJECTLIT() *GlobalNamespace_NameProp {
	clinit_GlobalNamespace_NameProp()
	return enum_GlobalNamespace_NameProp_OBJECTLIT
}

var enum_GlobalNamespace_NameProp_FUNCTION *GlobalNamespace_NameProp

func StaticGet_GlobalNamespace_NameProp_FUNCTION() *GlobalNamespace_NameProp {
	clinit_GlobalNamespace_NameProp()
	return enum_GlobalNamespace_NameProp_FUNCTION
}

var enum_GlobalNamespace_NameProp_GET_SET *GlobalNamespace_NameProp

func StaticGet_GlobalNamespace_NameProp_GET_SET() *GlobalNamespace_NameProp {
	clinit_GlobalNamespace_NameProp()
	return enum_GlobalNamespace_NameProp_GET_SET
}

var enum_GlobalNamespace_NameProp_OTHER_OBJECT *GlobalNamespace_NameProp

func StaticGet_GlobalNamespace_NameProp_OTHER_OBJECT() *GlobalNamespace_NameProp {
	clinit_GlobalNamespace_NameProp()
	return enum_GlobalNamespace_NameProp_OTHER_OBJECT
}

var enum_GlobalNamespace_NameProp_MULTI_REF_MODE *GlobalNamespace_NameProp

func StaticGet_GlobalNamespace_NameProp_MULTI_REF_MODE() *GlobalNamespace_NameProp {
	clinit_GlobalNamespace_NameProp()
	return enum_GlobalNamespace_NameProp_MULTI_REF_MODE
}

var static_GlobalNamespace_NameProp_DECLARED_TYPE_KIND_MASK int32

func StaticGet_GlobalNamespace_NameProp_DECLARED_TYPE_KIND_MASK() int32 {
	clinit_GlobalNamespace_NameProp()
	return static_GlobalNamespace_NameProp_DECLARED_TYPE_KIND_MASK
}
func StaticSet_GlobalNamespace_NameProp_DECLARED_TYPE_KIND_MASK(v int32) {
	clinit_GlobalNamespace_NameProp()
	static_GlobalNamespace_NameProp_DECLARED_TYPE_KIND_MASK = v
}

var static_GlobalNamespace_NameProp_NAME_KIND_MASK int32

func StaticGet_GlobalNamespace_NameProp_NAME_KIND_MASK() int32 {
	clinit_GlobalNamespace_NameProp()
	return static_GlobalNamespace_NameProp_NAME_KIND_MASK
}
func StaticSet_GlobalNamespace_NameProp_NAME_KIND_MASK(v int32) {
	clinit_GlobalNamespace_NameProp()
	static_GlobalNamespace_NameProp_NAME_KIND_MASK = v
}

var clinitDone_GlobalNamespace_NameProp uint32
var clinitState_GlobalNamespace_NameProp j2gort.ClinitState

func clinit_GlobalNamespace_NameProp() {
	if atomic.LoadUint32(&clinitDone_GlobalNamespace_NameProp) != 0 {
		return
	}
	clinitSlow_GlobalNamespace_NameProp()
}

func clinitSlow_GlobalNamespace_NameProp() {
	if !j2gort.EnterClinit(&clinitState_GlobalNamespace_NameProp, &clinitDone_GlobalNamespace_NameProp) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_GlobalNamespace_NameProp, &clinitDone_GlobalNamespace_NameProp)
	enum_GlobalNamespace_NameProp_IS_DECLARED = &GlobalNamespace_NameProp{}
	initEnumGlobalNamespace_NameProp(enum_GlobalNamespace_NameProp_IS_DECLARED, "IS_DECLARED", int32(0), int32(0))
	enum_GlobalNamespace_NameProp_IS_MODULE_PROP = &GlobalNamespace_NameProp{}
	initEnumGlobalNamespace_NameProp(enum_GlobalNamespace_NameProp_IS_MODULE_PROP, "IS_MODULE_PROP", int32(1), int32(1))
	enum_GlobalNamespace_NameProp_IS_PROVIDED = &GlobalNamespace_NameProp{}
	initEnumGlobalNamespace_NameProp(enum_GlobalNamespace_NameProp_IS_PROVIDED, "IS_PROVIDED", int32(2), int32(2))
	enum_GlobalNamespace_NameProp_IS_USED_HAS_OWN_PROPERTY = &GlobalNamespace_NameProp{}
	initEnumGlobalNamespace_NameProp(enum_GlobalNamespace_NameProp_IS_USED_HAS_OWN_PROPERTY, "IS_USED_HAS_OWN_PROPERTY", int32(3), int32(3))
	enum_GlobalNamespace_NameProp_IS_EXTERN = &GlobalNamespace_NameProp{}
	initEnumGlobalNamespace_NameProp(enum_GlobalNamespace_NameProp_IS_EXTERN, "IS_EXTERN", int32(4), int32(4))
	enum_GlobalNamespace_NameProp_CONSTRUCTOR_TYPE = &GlobalNamespace_NameProp{}
	initEnumGlobalNamespace_NameProp(enum_GlobalNamespace_NameProp_CONSTRUCTOR_TYPE, "CONSTRUCTOR_TYPE", int32(5), int32(5))
	enum_GlobalNamespace_NameProp_INTERFACE_TYPE = &GlobalNamespace_NameProp{}
	initEnumGlobalNamespace_NameProp(enum_GlobalNamespace_NameProp_INTERFACE_TYPE, "INTERFACE_TYPE", int32(6), int32(6))
	enum_GlobalNamespace_NameProp_ENUM_TYPE = &GlobalNamespace_NameProp{}
	initEnumGlobalNamespace_NameProp(enum_GlobalNamespace_NameProp_ENUM_TYPE, "ENUM_TYPE", int32(7), int32(7))
	enum_GlobalNamespace_NameProp_NOT_A_TYPE = &GlobalNamespace_NameProp{}
	initEnumGlobalNamespace_NameProp(enum_GlobalNamespace_NameProp_NOT_A_TYPE, "NOT_A_TYPE", int32(8), int32(8))
	enum_GlobalNamespace_NameProp_CLASS = &GlobalNamespace_NameProp{}
	initEnumGlobalNamespace_NameProp(enum_GlobalNamespace_NameProp_CLASS, "CLASS", int32(9), int32(9))
	enum_GlobalNamespace_NameProp_OBJECTLIT = &GlobalNamespace_NameProp{}
	initEnumGlobalNamespace_NameProp(enum_GlobalNamespace_NameProp_OBJECTLIT, "OBJECTLIT", int32(10), int32(10))
	enum_GlobalNamespace_NameProp_FUNCTION = &GlobalNamespace_NameProp{}
	initEnumGlobalNamespace_NameProp(enum_GlobalNamespace_NameProp_FUNCTION, "FUNCTION", int32(11), int32(11))
	enum_GlobalNamespace_NameProp_GET_SET = &GlobalNamespace_NameProp{}
	initEnumGlobalNamespace_NameProp(enum_GlobalNamespace_NameProp_GET_SET, "GET_SET", int32(12), int32(12))
	enum_GlobalNamespace_NameProp_OTHER_OBJECT = &GlobalNamespace_NameProp{}
	initEnumGlobalNamespace_NameProp(enum_GlobalNamespace_NameProp_OTHER_OBJECT, "OTHER_OBJECT", int32(13), int32(13))
	enum_GlobalNamespace_NameProp_MULTI_REF_MODE = &GlobalNamespace_NameProp{}
	initEnumGlobalNamespace_NameProp(enum_GlobalNamespace_NameProp_MULTI_REF_MODE, "MULTI_REF_MODE", int32(14), int32(14))
	enumValues_GlobalNamespace_NameProp = []*GlobalNamespace_NameProp{enum_GlobalNamespace_NameProp_IS_DECLARED, enum_GlobalNamespace_NameProp_IS_MODULE_PROP, enum_GlobalNamespace_NameProp_IS_PROVIDED, enum_GlobalNamespace_NameProp_IS_USED_HAS_OWN_PROPERTY, enum_GlobalNamespace_NameProp_IS_EXTERN, enum_GlobalNamespace_NameProp_CONSTRUCTOR_TYPE, enum_GlobalNamespace_NameProp_INTERFACE_TYPE, enum_GlobalNamespace_NameProp_ENUM_TYPE, enum_GlobalNamespace_NameProp_NOT_A_TYPE, enum_GlobalNamespace_NameProp_CLASS, enum_GlobalNamespace_NameProp_OBJECTLIT, enum_GlobalNamespace_NameProp_FUNCTION, enum_GlobalNamespace_NameProp_GET_SET, enum_GlobalNamespace_NameProp_OTHER_OBJECT, enum_GlobalNamespace_NameProp_MULTI_REF_MODE}
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.GlobalNamespace.NameProp", []any{enum_GlobalNamespace_NameProp_IS_DECLARED, enum_GlobalNamespace_NameProp_IS_MODULE_PROP, enum_GlobalNamespace_NameProp_IS_PROVIDED, enum_GlobalNamespace_NameProp_IS_USED_HAS_OWN_PROPERTY, enum_GlobalNamespace_NameProp_IS_EXTERN, enum_GlobalNamespace_NameProp_CONSTRUCTOR_TYPE, enum_GlobalNamespace_NameProp_INTERFACE_TYPE, enum_GlobalNamespace_NameProp_ENUM_TYPE, enum_GlobalNamespace_NameProp_NOT_A_TYPE, enum_GlobalNamespace_NameProp_CLASS, enum_GlobalNamespace_NameProp_OBJECTLIT, enum_GlobalNamespace_NameProp_FUNCTION, enum_GlobalNamespace_NameProp_GET_SET, enum_GlobalNamespace_NameProp_OTHER_OBJECT, enum_GlobalNamespace_NameProp_MULTI_REF_MODE})
	j2gort.RegisterEnumConstants("com.google.javascript.jscomp.GlobalNamespace$NameProp", []any{enum_GlobalNamespace_NameProp_IS_DECLARED, enum_GlobalNamespace_NameProp_IS_MODULE_PROP, enum_GlobalNamespace_NameProp_IS_PROVIDED, enum_GlobalNamespace_NameProp_IS_USED_HAS_OWN_PROPERTY, enum_GlobalNamespace_NameProp_IS_EXTERN, enum_GlobalNamespace_NameProp_CONSTRUCTOR_TYPE, enum_GlobalNamespace_NameProp_INTERFACE_TYPE, enum_GlobalNamespace_NameProp_ENUM_TYPE, enum_GlobalNamespace_NameProp_NOT_A_TYPE, enum_GlobalNamespace_NameProp_CLASS, enum_GlobalNamespace_NameProp_OBJECTLIT, enum_GlobalNamespace_NameProp_FUNCTION, enum_GlobalNamespace_NameProp_GET_SET, enum_GlobalNamespace_NameProp_OTHER_OBJECT, enum_GlobalNamespace_NameProp_MULTI_REF_MODE})
	static_GlobalNamespace_NameProp_DECLARED_TYPE_KIND_MASK = (((StaticGet_GlobalNamespace_NameProp_CONSTRUCTOR_TYPE().bit | StaticGet_GlobalNamespace_NameProp_INTERFACE_TYPE().bit) | StaticGet_GlobalNamespace_NameProp_ENUM_TYPE().bit) | StaticGet_GlobalNamespace_NameProp_NOT_A_TYPE().bit)
	static_GlobalNamespace_NameProp_NAME_KIND_MASK = ((((StaticGet_GlobalNamespace_NameProp_CLASS().bit | StaticGet_GlobalNamespace_NameProp_OBJECTLIT().bit) | StaticGet_GlobalNamespace_NameProp_FUNCTION().bit) | StaticGet_GlobalNamespace_NameProp_GET_SET().bit) | StaticGet_GlobalNamespace_NameProp_OTHER_OBJECT().bit)
}

func init() {
	j2gort.RegisterEnumConstantsLoader("com.google.javascript.jscomp.GlobalNamespace.NameProp", "com.google.javascript.jscomp.GlobalNamespace$NameProp", func() {
		clinit_GlobalNamespace_NameProp()
	})
}

var enumValues_GlobalNamespace_NameProp []*GlobalNamespace_NameProp

func GlobalNamespace_NameProp_Values() []*GlobalNamespace_NameProp {
	clinit_GlobalNamespace_NameProp()
	return j2gort.CloneSlice(enumValues_GlobalNamespace_NameProp)
}

func GlobalNamespace_NameProp_ValueOf(name string) *GlobalNamespace_NameProp {
	clinit_GlobalNamespace_NameProp()
	for _, c := range enumValues_GlobalNamespace_NameProp {
		if c.Name() == name {
			return c
		}
	}
	panic(javalang.NewIllegalArgumentException("No enum constant com.google.javascript.jscomp.GlobalNamespace.NameProp." + name))
}

func (this_ *GlobalNamespace_NameProp) IsJavaNil_() bool { return this_ == nil }
func (*GlobalNamespace_NameProp) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GlobalNamespace$NameProp", "NameProp"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GlobalNamespace.NameProp", "com.google.javascript.jscomp.GlobalNamespace$NameProp", func(v_ any) bool {
		return j2gort.InstanceOf[*GlobalNamespace_NameProp](v_)
	})
}

// Java record: com.google.javascript.jscomp.GlobalNamespace.ObjLitStringKeyAnalysis
type GlobalNamespace_ObjLitStringKeyAnalysis struct {
	nameString string                    `gson:"nameString"`
	nameType   *GlobalNamespace_NameProp `gson:"nameType"`
}

func (this_ *GlobalNamespace_ObjLitStringKeyAnalysis) BindSelf_(s any) {
	_ = s
}
func (this_ *GlobalNamespace_ObjLitStringKeyAnalysis) IsJavaNil_() bool { return this_ == nil }
func (*GlobalNamespace_ObjLitStringKeyAnalysis) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.GlobalNamespace$ObjLitStringKeyAnalysis", "ObjLitStringKeyAnalysis"
}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.GlobalNamespace.ObjLitStringKeyAnalysis", "com.google.javascript.jscomp.GlobalNamespace$ObjLitStringKeyAnalysis", func(v_ any) bool {
		return j2gort.InstanceOf[*GlobalNamespace_ObjLitStringKeyAnalysis](v_)
	})
}
func (*GlobalNamespace_ObjLitStringKeyAnalysis) Implements_j2gort_Ext_java_lang_Record() {}

func (this_ *GlobalNamespace_ObjLitStringKeyAnalysis) initFields_GlobalNamespace_ObjLitStringKeyAnalysis() {
	_ = this_
}

// Java constructor: ObjLitStringKeyAnalysis(java.lang.@org.jspecify.annotations.Nullable String nameString, com.google.javascript.jscomp.GlobalNamespace.@org.jspecify.annotations.Nullable NameProp nameType)
func NewGlobalNamespace_ObjLitStringKeyAnalysis(nameString string, nameType *GlobalNamespace_NameProp) *GlobalNamespace_ObjLitStringKeyAnalysis {
	this_ := &GlobalNamespace_ObjLitStringKeyAnalysis{}
	InitGlobalNamespace_ObjLitStringKeyAnalysis(this_, nameString, nameType)
	return this_
}
func allocScratch_GlobalNamespace_ObjLitStringKeyAnalysis(dst_ *GlobalNamespace_ObjLitStringKeyAnalysis, nameString string, nameType *GlobalNamespace_NameProp) *GlobalNamespace_ObjLitStringKeyAnalysis {
	*dst_ = GlobalNamespace_ObjLitStringKeyAnalysis{}
	this_ := dst_
	InitGlobalNamespace_ObjLitStringKeyAnalysis(this_, nameString, nameType)
	return this_
}

func InitGlobalNamespace_ObjLitStringKeyAnalysis(this_ *GlobalNamespace_ObjLitStringKeyAnalysis, nameString string, nameType *GlobalNamespace_NameProp) {
	_ = nameString
	_ = nameType
	this_.initFields_GlobalNamespace_ObjLitStringKeyAnalysis()
	this_.nameString = nameString
	this_.nameType = nameType
}

// Java static method: ObjLitStringKeyAnalysis.forObjectDefineProperty(java.lang.String nameString)
func GlobalNamespace_ObjLitStringKeyAnalysis_ForObjectDefineProperty(nameString string) *GlobalNamespace_ObjLitStringKeyAnalysis {
	_ = nameString
	return NewGlobalNamespace_ObjLitStringKeyAnalysis(j2gort.CheckNotNullString(nameString), StaticGet_GlobalNamespace_NameProp_GET_SET())
}
func allocScratchVia_GlobalNamespace_ObjLitStringKeyAnalysis_ForObjectDefineProperty(dst_ *GlobalNamespace_ObjLitStringKeyAnalysis, nameString string) *GlobalNamespace_ObjLitStringKeyAnalysis {
	_ = nameString
	return allocScratch_GlobalNamespace_ObjLitStringKeyAnalysis(dst_, j2gort.CheckNotNullString(nameString), StaticGet_GlobalNamespace_NameProp_GET_SET())
}

// Java static method: ObjLitStringKeyAnalysis.forObjLitAssignment(java.lang.String nameString, com.google.javascript.jscomp.GlobalNamespace.NameProp nameType)
func GlobalNamespace_ObjLitStringKeyAnalysis_ForObjLitAssignment(nameString string, nameType *GlobalNamespace_NameProp) *GlobalNamespace_ObjLitStringKeyAnalysis {
	_ = nameString
	_ = nameType
	return NewGlobalNamespace_ObjLitStringKeyAnalysis(j2gort.CheckNotNullString(nameString), nameType)
}
func allocScratchVia_GlobalNamespace_ObjLitStringKeyAnalysis_ForObjLitAssignment(dst_ *GlobalNamespace_ObjLitStringKeyAnalysis, nameString string, nameType *GlobalNamespace_NameProp) *GlobalNamespace_ObjLitStringKeyAnalysis {
	_ = nameString
	_ = nameType
	return allocScratch_GlobalNamespace_ObjLitStringKeyAnalysis(dst_, j2gort.CheckNotNullString(nameString), nameType)
}

// Java static method: ObjLitStringKeyAnalysis.forNonReference()
func GlobalNamespace_ObjLitStringKeyAnalysis_ForNonReference() *GlobalNamespace_ObjLitStringKeyAnalysis {
	return NewGlobalNamespace_ObjLitStringKeyAnalysis(j2gort.NullString, StaticGet_GlobalNamespace_NameProp_OTHER_OBJECT())
}
func allocScratchVia_GlobalNamespace_ObjLitStringKeyAnalysis_ForNonReference(dst_ *GlobalNamespace_ObjLitStringKeyAnalysis) *GlobalNamespace_ObjLitStringKeyAnalysis {
	return allocScratch_GlobalNamespace_ObjLitStringKeyAnalysis(dst_, j2gort.NullString, StaticGet_GlobalNamespace_NameProp_OTHER_OBJECT())
}

func (this_ *GlobalNamespace_ObjLitStringKeyAnalysis) NameString() string { return this_.nameString }
func (this_ *GlobalNamespace_ObjLitStringKeyAnalysis) NameType() *GlobalNamespace_NameProp {
	return this_.nameType
}
func (this_ *GlobalNamespace_ObjLitStringKeyAnalysis) Equals(other_ any) bool {
	o_, ok_ := other_.(*GlobalNamespace_ObjLitStringKeyAnalysis)
	if !ok_ || o_ == nil {
		return false
	}
	if this_ == o_ {
		return true
	}
	if this_.nameString != o_.nameString {
		return false
	}
	if this_.nameType != o_.nameType {
		return false
	}
	return true
}
func (this_ *GlobalNamespace_ObjLitStringKeyAnalysis) HashCode() int32 {
	var h_ int32 = 0
	{
		hs_44 := this_.nameString
		h_ = 31 * h_
		if hs_44 != j2gort.NullString {
			h_ += javalang.StringHashCode(hs_44)
		}
	}
	{
		hp_45 := this_.nameType
		h_ = 31 * h_
		if hp_45 != nil {
			h_ += hp_45.HashCode()
		}
	}
	return h_
}
func (this_ *GlobalNamespace_ObjLitStringKeyAnalysis) String() string {
	_ = fmt.Sprint
	return "ObjLitStringKeyAnalysis[" + "nameString=" + j2gort.StrVal(j2gort.BoxString(this_.nameString)) + ", " + "nameType=" + j2gort.StrVal(j2gort.PtrToAny(this_.nameType)) + "]"
}
func (this_ *GlobalNamespace) ReleaseSlabs_() {
	if this_ == nil {
		return
	}
	this_.globalNames = nil
	for cur_ := this_.nameMap.IterValues(); cur_.HasNext(); {
		nm_ := cur_.Next()
		if nm_ != nil && nm_.refsForNodeMap != nil {
			j2gort.ReleaseOrderedMap(nm_.refsForNodeMap)
			nm_.refsForNodeMap = nil
		}
	}
	this_.nameMap.ReleaseToPool()
	for curMod_ := this_.nameMapByModule.IterValues(); curMod_.HasNext(); {
		mnm_ := curMod_.Next()
		if mnm_ != nil {
			for curNm_ := mnm_.names.IterValues(); curNm_.HasNext(); {
				nm_ := curNm_.Next()
				if nm_ != nil && nm_.refsForNodeMap != nil {
					j2gort.ReleaseOrderedMap(nm_.refsForNodeMap)
					nm_.refsForNodeMap = nil
				}
			}
			mnm_.names.ReleaseToPool()
		}
	}
	this_.nameMapByModule.ReleaseToPool()
	this_.slab_GlobalNamespace_Ref.Release()
	this_.scopeCreator_.ReleaseSlabs_()
}
func (*GlobalNamespace_BuildGlobalNamespace) NoOpVisit_() {}
