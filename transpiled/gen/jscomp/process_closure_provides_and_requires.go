/*
 * Copyright 2019 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/ProcessClosureProvidesAndRequires.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
)

// Java class: com.google.javascript.jscomp.ProcessClosureProvidesAndRequires
type ProcessClosureProvidesAndRequires struct {
	compiler                        AbstractCompiler                  `gson:"compiler"`
	chunkGraph                      *JSChunkGraph                     `gson:"chunkGraph"`
	providedNames                   j2gort.Ext_java_util_Map          `gson:"providedNames"`
	providedNameHashes              []int32                           `gson:"providedNameHashes"`
	exportedVariables               j2gort.Ext_java_util_Set          `gson:"exportedVariables"`
	requiresToBeRemoved             []pkg_rhino.Node                  `gson:"requiresToBeRemoved"`
	forwardDeclaresToRemove         j2gort.OrderedSet[pkg_rhino.Node] `gson:"forwardDeclaresToRemove"`
	astFactory                      *AstFactory                       `gson:"astFactory"`
	providedNameHashCount           int32                             `gson:"providedNameHashCount"`
	idHash_                         atomic.Int32
	preserveGoogProvidesAndRequires bool `gson:"preserveGoogProvidesAndRequires"`
	hasRewritingOccurred            bool `gson:"hasRewritingOccurred"`
}

func (this_ *ProcessClosureProvidesAndRequires) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessClosureProvidesAndRequires) IsJavaNil_() bool { return this_ == nil }
func (*ProcessClosureProvidesAndRequires) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessClosureProvidesAndRequires", "ProcessClosureProvidesAndRequires"
}
func (this_ *ProcessClosureProvidesAndRequires) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessClosureProvidesAndRequires) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessClosureProvidesAndRequires", "com.google.javascript.jscomp.ProcessClosureProvidesAndRequires", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessClosureProvidesAndRequires](v_)
	})
}
func (*ProcessClosureProvidesAndRequires) Implements_jscomp_CompilerPass() {}

var static_ProcessClosureProvidesAndRequires_TYPEDEF_CHILD_OF_PROVIDE *DiagnosticType

func StaticGet_ProcessClosureProvidesAndRequires_TYPEDEF_CHILD_OF_PROVIDE() *DiagnosticType {
	clinit_ProcessClosureProvidesAndRequires()
	return static_ProcessClosureProvidesAndRequires_TYPEDEF_CHILD_OF_PROVIDE
}
func StaticSet_ProcessClosureProvidesAndRequires_TYPEDEF_CHILD_OF_PROVIDE(v *DiagnosticType) {
	clinit_ProcessClosureProvidesAndRequires()
	static_ProcessClosureProvidesAndRequires_TYPEDEF_CHILD_OF_PROVIDE = v
}
func StaticGet_ProcessClosureProvidesAndRequires_GOOG() string { return "goog" }

var clinitDone_ProcessClosureProvidesAndRequires uint32
var clinitState_ProcessClosureProvidesAndRequires j2gort.ClinitState

func clinit_ProcessClosureProvidesAndRequires() {
	if atomic.LoadUint32(&clinitDone_ProcessClosureProvidesAndRequires) != 0 {
		return
	}
	clinitSlow_ProcessClosureProvidesAndRequires()
}

func clinitSlow_ProcessClosureProvidesAndRequires() {
	if !j2gort.EnterClinit(&clinitState_ProcessClosureProvidesAndRequires, &clinitDone_ProcessClosureProvidesAndRequires) {
		return
	}
	defer j2gort.ExitClinit(&clinitState_ProcessClosureProvidesAndRequires, &clinitDone_ProcessClosureProvidesAndRequires)
	static_ProcessClosureProvidesAndRequires_TYPEDEF_CHILD_OF_PROVIDE = DiagnosticType_Error("JSC_TYPEDEF_CHILD_OF_PROVIDE", "invalid @typedef goog.provide {0}\nParent namespace {1} is goog.provided and initialized in the same file")
}

func (this_ *ProcessClosureProvidesAndRequires) initFields_ProcessClosureProvidesAndRequires() {
	_ = this_
	this_.providedNames = j2gort.NewExt_java_util_LinkedHashMap()
	this_.exportedVariables = j2gort.NewExt_java_util_LinkedHashSet()
	this_.requiresToBeRemoved = []pkg_rhino.Node{}

	this_.forwardDeclaresToRemove = j2gort.NewOrderedSet[pkg_rhino.Node]()
}

// Java constructor: ProcessClosureProvidesAndRequires(com.google.javascript.jscomp.AbstractCompiler compiler, boolean preserveGoogProvidesAndRequires)
func NewProcessClosureProvidesAndRequires(compiler AbstractCompiler, preserveGoogProvidesAndRequires bool) *ProcessClosureProvidesAndRequires {
	clinit_ProcessClosureProvidesAndRequires()
	this_ := &ProcessClosureProvidesAndRequires{}
	InitProcessClosureProvidesAndRequires(this_, compiler, preserveGoogProvidesAndRequires)
	return this_
}

func InitProcessClosureProvidesAndRequires(this_ *ProcessClosureProvidesAndRequires, compiler AbstractCompiler, preserveGoogProvidesAndRequires bool) {
	_ = compiler
	_ = preserveGoogProvidesAndRequires
	this_.initFields_ProcessClosureProvidesAndRequires()
	this_.compiler = compiler
	this_.chunkGraph = compiler.GetChunkGraph()
	this_.preserveGoogProvidesAndRequires = preserveGoogProvidesAndRequires
	this_.astFactory = compiler.CreateAstFactory()
}

// Java method: ProcessClosureProvidesAndRequires.putProvidedName(java.lang.String ns, com.google.javascript.jscomp.ProcessClosureProvidesAndRequires.ProvidedName pn)
func (this_ *ProcessClosureProvidesAndRequires) priv_ProcessClosureProvidesAndRequires_putProvidedName_4467(ns string, pn *ProcessClosureProvidesAndRequires_ProvidedName) {
	_ = ns
	_ = pn
	j2gort.MapPutStr(this_.providedNames, ns, j2gort.PtrToAny(pn))
	var h int32 = pkg_rhino.Node_HashQualifiedNameString(ns)
	_ = h
	if h == int32(0) {
		return
	}
	var table []int32 = this_.providedNameHashes
	_ = table
	if table == nil {
		table = make([]int32, int32(16))
		this_.providedNameHashes = table
	} else if (this_.providedNameHashCount * int32(4)) >= (int32(len(table)) * int32(3)) {
		var next []int32 = make([]int32, (int32(len(table)) * int32(2)))
		_ = next
		var nextMask int32 = (int32(len(next)) - int32(1))
		_ = nextMask
		{
			var i int32 = int32(0)
			_ = i
			var len_ int32 = int32(len(table))
			_ = len_
			for ; i < len_; i = i + 1 {
				var existing int32 = table[i]
				_ = existing
				if existing != int32(0) {
					var s int32 = (existing & nextMask)
					_ = s
					for next[s] != int32(0) {
						s = ((s + int32(1)) & nextMask)
					}
					{
						a_1 := next
						i_2 := s
						a_1[i_2] = existing
					}
				}
			}
		}
		table = next
		this_.providedNameHashes = table
	}
	var mask int32 = (int32(len(table)) - int32(1))
	_ = mask
	var slot int32 = (h & mask)
	_ = slot
	for {
		var existing_2 int32 = table[slot]
		_ = existing_2
		if existing_2 == int32(0) {
			{
				a_3 := table
				i_4 := slot
				a_3[i_4] = h
			}
			this_.providedNameHashCount = this_.providedNameHashCount + 1
			return
		}
		if existing_2 == h {
			return
		}
		slot = ((slot + int32(1)) & mask)
	}
}

// Java method: ProcessClosureProvidesAndRequires.containsProvidedNameHash(int h)
func (this_ *ProcessClosureProvidesAndRequires) priv_ProcessClosureProvidesAndRequires_containsProvidedNameHash_4468(h int32) bool {
	_ = h
	var table []int32 = this_.providedNameHashes
	_ = table
	if (table == nil) || (h == int32(0)) {
		return false
	}
	var mask int32 = (int32(len(table)) - int32(1))
	_ = mask
	var slot int32 = (h & mask)
	_ = slot
	for {
		var existing int32 = table[slot]
		_ = existing
		if existing == int32(0) {
			return false
		}
		if existing == h {
			return true
		}
		slot = ((slot + int32(1)) & mask)
	}
	panic("unreachable")
}

// Java method: ProcessClosureProvidesAndRequires.process(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ProcessClosureProvidesAndRequires) Process_2(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	this_.RewriteProvidesAndRequires(externs, root)
}

// Java method: ProcessClosureProvidesAndRequires.getExportedVariableNames()
func (this_ *ProcessClosureProvidesAndRequires) GetExportedVariableNames() j2gort.Ext_java_util_Set {
	return this_.exportedVariables
}

// Java method: ProcessClosureProvidesAndRequires.collectProvidedNames(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ProcessClosureProvidesAndRequires) CollectProvidedNames(externs pkg_rhino.Node, root pkg_rhino.Node) j2gort.Ext_java_util_Map {
	_ = externs
	_ = root
	if this_.providedNames.IsEmpty() {
		this_.priv_ProcessClosureProvidesAndRequires_putProvidedName_4467("goog", NewProcessClosureProvidesAndRequires_ProvidedNameBuilder(this_).SetNamespace("goog").SetNode(pkg_rhino.Node(0)).SetChunk(nil).Build())
		NodeTraversal_TraverseRoots(this_.compiler, NewProcessClosureProvidesAndRequires_CollectDefinitions(this_), externs, root)
	}
	return this_.providedNames
}

// Java method: ProcessClosureProvidesAndRequires.rewriteProvidesAndRequires(com.google.javascript.rhino.Node externs, com.google.javascript.rhino.Node root)
func (this_ *ProcessClosureProvidesAndRequires) RewriteProvidesAndRequires(externs pkg_rhino.Node, root pkg_rhino.Node) {
	_ = externs
	_ = root
	if !(!(this_.hasRewritingOccurred)) {
		j2gort.CheckStateStr(false, "Cannot call rewriteProvidesAndRequires twice per instance")
	}
	this_.hasRewritingOccurred = true
	this_.CollectProvidedNames(externs, root)
	{
		cur_5 := j2gort.NewMapValCursor(this_.providedNames)
		for cur_5.HasNext() {
			rawElem_6 := cur_5.Next()
			pn := j2gort.CastPtr[ProcessClosureProvidesAndRequires_ProvidedName](rawElem_6)
			_ = pn
			pn.priv_ProvidedName_replace_4501(this_.preserveGoogProvidesAndRequires, this_.providedNames)
		}
	}
	for _, closureRequire := range this_.requiresToBeRemoved {
		_ = closureRequire
		this_.compiler.ReportChangeToEnclosingScope(closureRequire)
		closureRequire.Detach()
	}

	{
		cur_8 := this_.forwardDeclaresToRemove.Iter()
		for cur_8.HasNext() {
			forwardDeclare := cur_8.
				Next()
			_ = forwardDeclare
			NodeUtil_DeleteNode(forwardDeclare, this_.compiler)
		}
	}
}

// Java method: ProcessClosureProvidesAndRequires.visitGoogMethodCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node parent, com.google.javascript.rhino.Node n, java.lang.String methodName)
func (this_ *ProcessClosureProvidesAndRequires) priv_ProcessClosureProvidesAndRequires_visitGoogMethodCall_4472(t *NodeTraversal, parent pkg_rhino.Node, n pkg_rhino.Node, methodName string) {
	_ = t
	_ = parent
	_ = n
	_ = methodName
	switch j2gort.CheckNotNullString(methodName) {
	case "exportSymbol":
		var arg pkg_rhino.Node = n.GetSecondChild()
		_ = arg
		if arg.IsStringLit() {
			var argString string = arg.GetString()
			_ = argString
			var dot int32 = j2gort.JStrIndexOfCp(argString, int32(uint16(46)))
			_ = dot
			if dot == int32(-1) {
				j2gort.SetAddStr(this_.exportedVariables, argString)
			} else {
				j2gort.SetAddStr(this_.exportedVariables, j2gort.JStrSub(argString, int32(0), dot))
			}
		}
	case "require", "requireType":
		if this_.priv_ProcessClosureProvidesAndRequires_isValidPrimitiveCall_4473(t, n) {
			this_.priv_ProcessClosureProvidesAndRequires_processRequireCall_4474(n, parent)
		}
	case "provide":
		if this_.priv_ProcessClosureProvidesAndRequires_isValidPrimitiveCall_4473(t, n) {
			this_.priv_ProcessClosureProvidesAndRequires_processProvideCall_4476(t, n, parent)
		}
	case "forwardDeclare":
		if this_.priv_ProcessClosureProvidesAndRequires_isValidPrimitiveCall_4473(t, n) {
			this_.priv_ProcessClosureProvidesAndRequires_processForwardDeclare_4481(n, parent)
		}
	default:
	}
}

// Java method: ProcessClosureProvidesAndRequires.isValidPrimitiveCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n)
func (this_ *ProcessClosureProvidesAndRequires) priv_ProcessClosureProvidesAndRequires_isValidPrimitiveCall_4473(t *NodeTraversal, n pkg_rhino.Node) bool {
	_ = t
	_ = n
	if this_.compiler.GetOptions().ShouldPreserveGoogModule() {
		return true
	}
	return (t.InGlobalHoistScope() && n.GetParent().IsExprResult())
}

// Java method: ProcessClosureProvidesAndRequires.processRequireCall(com.google.javascript.rhino.Node call, com.google.javascript.rhino.Node parent)
func (this_ *ProcessClosureProvidesAndRequires) priv_ProcessClosureProvidesAndRequires_processRequireCall_4474(call pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = call
	_ = parent
	if !(this_.priv_ProcessClosureProvidesAndRequires_verifyOnlyArgumentIsString_4482(call)) {
		return
	}
	if !(this_.preserveGoogProvidesAndRequires) {
		this_.requiresToBeRemoved = append(this_.requiresToBeRemoved, parent)
	}
}

// Java method: ProcessClosureProvidesAndRequires.processLegacyModuleCall(java.lang.String namespace, com.google.javascript.rhino.Node googModuleCall, com.google.javascript.jscomp.JSChunk chunk)
func (this_ *ProcessClosureProvidesAndRequires) priv_ProcessClosureProvidesAndRequires_processLegacyModuleCall_4475(namespace string, googModuleCall pkg_rhino.Node, chunk *JSChunk) {
	_ = namespace
	_ = googModuleCall
	_ = chunk
	this_.priv_ProcessClosureProvidesAndRequires_registerAnyProvidedPrefixes_4483(namespace, googModuleCall, chunk)
	this_.priv_ProcessClosureProvidesAndRequires_putProvidedName_4467(namespace, NewProcessClosureProvidesAndRequires_ProvidedNameBuilder(this_).SetNamespace(namespace).SetNode(googModuleCall).SetChunk(chunk).SetExplicit(true).SetFromLegacyModule(true).Build())
}

// Java method: ProcessClosureProvidesAndRequires.processProvideCall(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node call, com.google.javascript.rhino.Node parent)
func (this_ *ProcessClosureProvidesAndRequires) priv_ProcessClosureProvidesAndRequires_processProvideCall_4476(t *NodeTraversal, call pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = call
	_ = parent
	if !(call.IsCall()) {
		j2gort.CheckState(false)
	}
	if !(this_.priv_ProcessClosureProvidesAndRequires_verifyOnlyArgumentIsString_4482(call)) {
		return
	}
	var left pkg_rhino.Node = call.GetFirstChild()
	_ = left
	var arg pkg_rhino.Node = left.GetNext()
	_ = arg
	var ns string = arg.GetString()
	_ = ns
	var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(call)
	_ = info
	var isImplicitlyInitialized bool = ((info != nil) && info.IsProvideAlreadyProvided())
	_ = isImplicitlyInitialized
	if j2gort.MapContainsKeyStr(this_.providedNames, ns) {
		var previouslyProvided *ProcessClosureProvidesAndRequires_ProvidedName = j2gort.CastPtr[ProcessClosureProvidesAndRequires_ProvidedName](j2gort.MapGetStr(this_.providedNames, ns))
		_ = previouslyProvided
		if !(previouslyProvided.IsExplicitlyProvided()) {
			previouslyProvided.AddProvide(parent, t.GetChunk(), true, this_.chunkGraph)
		}
	} else {
		this_.priv_ProcessClosureProvidesAndRequires_registerAnyProvidedPrefixes_4483(ns, parent, t.GetChunk())
		this_.priv_ProcessClosureProvidesAndRequires_putProvidedName_4467(ns, NewProcessClosureProvidesAndRequires_ProvidedNameBuilder(this_).SetNamespace(ns).SetNode(parent).SetChunk(t.GetChunk()).SetExplicit(true).SetHasImplicitInitialization(isImplicitlyInitialized).Build())
	}
}

// Java method: ProcessClosureProvidesAndRequires.handleStubDefinition(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node exprResult)
func (this_ *ProcessClosureProvidesAndRequires) priv_ProcessClosureProvidesAndRequires_handleStubDefinition_4477(t *NodeTraversal, exprResult pkg_rhino.Node) {
	_ = t
	_ = exprResult
	if !(t.InGlobalHoistScope()) {
		return
	}
	var isExternStub bool = exprResult.IsFromExterns()
	_ = isExternStub
	var isTypedefStub bool = this_.priv_ProcessClosureProvidesAndRequires_isTypedefStubDeclaration_4479(exprResult)
	_ = isTypedefStub
	var isValidTypedefStubDefinition bool = (isTypedefStub && !(this_.hasRewritingOccurred))
	_ = isValidTypedefStubDefinition
	if isValidTypedefStubDefinition || isExternStub {
		var target pkg_rhino.Node = exprResult.GetFirstChild()
		_ = target
		var h int32 = target.HashQualifiedName()
		_ = h
		if (h != int32(0)) && this_.priv_ProcessClosureProvidesAndRequires_containsProvidedNameHash_4468(h) {
			var name string = target.GetQualifiedName()
			_ = name
			if name != j2gort.NullString {
				var pn *ProcessClosureProvidesAndRequires_ProvidedName = j2gort.CastPtr[ProcessClosureProvidesAndRequires_ProvidedName](j2gort.MapGetStr(this_.providedNames, name))
				_ = pn
				if pn != nil {
					pn.priv_ProvidedName_addDefinition_4499(exprResult, t.GetChunk(), this_.chunkGraph)
				}
			}
		}
	}
	if isTypedefStub {
		this_.priv_ProcessClosureProvidesAndRequires_checkNestedTypedefProvide_4478(exprResult)
	}
}

// Java method: ProcessClosureProvidesAndRequires.checkNestedTypedefProvide(com.google.javascript.rhino.Node exprResult)
func (this_ *ProcessClosureProvidesAndRequires) priv_ProcessClosureProvidesAndRequires_checkNestedTypedefProvide_4478(exprResult pkg_rhino.Node) {
	_ = exprResult
	var target pkg_rhino.Node = exprResult.GetFirstChild()
	_ = target
	var h int32 = target.HashQualifiedName()
	_ = h
	if (h == int32(0)) || !(this_.priv_ProcessClosureProvidesAndRequires_containsProvidedNameHash_4468(h)) {
		return
	}
	var name string = target.GetQualifiedName()
	_ = name
	if (name == j2gort.NullString) || !(j2gort.JStrIndexOf(name, ".") >= 0) {
		return
	}
	if !(j2gort.MapContainsKeyStr(this_.providedNames, name)) {
		return
	}
	var parentName string = j2gort.JStrSub(name, int32(0), j2gort.JStrLastIndexOf(name, "."))
	_ = parentName
	var parent *ProcessClosureProvidesAndRequires_ProvidedName = j2gort.CastPtr[ProcessClosureProvidesAndRequires_ProvidedName](j2gort.MapGetStr(this_.providedNames, parentName))
	_ = parent
	var parentDefinition pkg_rhino.Node = parent.GetCandidateDefinition()
	_ = parentDefinition
	if ((parentDefinition == 0) || !(javalang.ObjectsEquals(j2gort.CheckNotNull(parentDefinition.GetStaticSourceFile()), exprResult.GetStaticSourceFile()))) || this_.priv_ProcessClosureProvidesAndRequires_isTypedefStubDeclaration_4479(parentDefinition) {
		return
	}
	this_.compiler.Report_1(JSError_Make_3(exprResult, StaticGet_ProcessClosureProvidesAndRequires_TYPEDEF_CHILD_OF_PROVIDE(), []string{name, parentName}))
}

// Java method: ProcessClosureProvidesAndRequires.isTypedefStubDeclaration(com.google.javascript.rhino.Node statement)
func (this_ *ProcessClosureProvidesAndRequires) priv_ProcessClosureProvidesAndRequires_isTypedefStubDeclaration_4479(statement pkg_rhino.Node) bool {
	_ = statement
	if !(statement.IsExprResult()) {
		return false
	}
	var info *pkg_rhino.JSDocInfo = NodeUtil_GetBestJSDocInfo(statement)
	_ = info
	return ((info != nil) && info.HasTypedefType())
}

// Java method: ProcessClosureProvidesAndRequires.handleCandidateProvideDefinition(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ProcessClosureProvidesAndRequires) priv_ProcessClosureProvidesAndRequires_handleCandidateProvideDefinition_4480(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	if !(t.InGlobalHoistScope()) {
		return
	}
	var name string
iifeSwitch_1:
	switch {
	default:
		switch n.GetParent().TokenSwitchOrd_() {
		case int32(111) /* Token.LET */, int32(97) /* Token.CONST */ :
			{
				if t.InGlobalScope() {
					name = n.GetString()
					break iifeSwitch_1
				}
				name = j2gort.NullString
				break iifeSwitch_1
			}
		case int32(84) /* Token.VAR */ :
			name = n.GetString()
			break iifeSwitch_1
		case int32(93) /* Token.EXPR_RESULT */ :
			if !(n.IsAssign()) {
				name = j2gort.NullString
				break iifeSwitch_1
			}
			var target pkg_rhino.Node = n.GetFirstChild()
			_ = target
			var h int32 = target.HashQualifiedName()
			_ = h
			if (h == int32(0)) || !(this_.priv_ProcessClosureProvidesAndRequires_containsProvidedNameHash_4468(h)) {
				name = j2gort.NullString
				break iifeSwitch_1
			}
			name = target.GetQualifiedName()
			break iifeSwitch_1
		case int32(105) /* Token.CLASS */, int32(73) /* Token.FUNCTION */ :
			name = j2gort.NullString
			break iifeSwitch_1
		default:
			name = j2gort.NullString
			break iifeSwitch_1
		}
		panic("unreachable switch expression")
	}

	_ = name
	if name == j2gort.NullString {
		return
	}
	var pn *ProcessClosureProvidesAndRequires_ProvidedName = j2gort.CastPtr[ProcessClosureProvidesAndRequires_ProvidedName](j2gort.MapGetStr(this_.providedNames, name))
	_ = pn
	if pn != nil {
		pn.priv_ProvidedName_addDefinition_4499(parent, t.GetChunk(), this_.chunkGraph)
	}
}

// Java method: ProcessClosureProvidesAndRequires.processForwardDeclare(com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ProcessClosureProvidesAndRequires) priv_ProcessClosureProvidesAndRequires_processForwardDeclare_4481(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	var convention CodingConvention = this_.compiler.GetCodingConvention()
	_ = convention
	var typeDeclarations j2gort.Ext_java_util_List = convention.IdentifyTypeDeclarationCall(n)
	_ = typeDeclarations
	if (!(this_.preserveGoogProvidesAndRequires) && (typeDeclarations != nil)) && (typeDeclarations.Size() == int32(1)) {
		var toRemove pkg_rhino.Node
		_ = toRemove
		if parent.IsExprResult() {
			toRemove = parent
		} else {
			toRemove = parent.GetParent()
		}
		this_.forwardDeclaresToRemove.Add(toRemove)
	}
}

// Java method: ProcessClosureProvidesAndRequires.verifyOnlyArgumentIsString(com.google.javascript.rhino.Node call)
func (this_ *ProcessClosureProvidesAndRequires) priv_ProcessClosureProvidesAndRequires_verifyOnlyArgumentIsString_4482(call pkg_rhino.Node) bool {
	_ = call
	var arg pkg_rhino.Node = call.GetSecondChild()
	_ = arg
	return (((arg != 0) && arg.IsStringLit()) && (arg.GetNext() == 0))
}

// Java method: ProcessClosureProvidesAndRequires.registerAnyProvidedPrefixes(java.lang.String ns, com.google.javascript.rhino.Node node, com.google.javascript.jscomp.JSChunk chunk)
func (this_ *ProcessClosureProvidesAndRequires) priv_ProcessClosureProvidesAndRequires_registerAnyProvidedPrefixes_4483(ns string, node pkg_rhino.Node, chunk *JSChunk) {
	_ = ns
	_ = node
	_ = chunk
	var pos int32 = j2gort.JStrIndexOfCp(ns, int32(uint16(46)))
	_ = pos
	for pos != int32(-1) {
		var prefixNs string = j2gort.JStrSub(ns, int32(0), pos)
		_ = prefixNs
		pos = j2gort.JStrIndexOfCpFrom(ns, int32(uint16(46)), (pos + int32(1)))
		if j2gort.MapContainsKeyStr(this_.providedNames, prefixNs) {
			j2gort.CastPtr[ProcessClosureProvidesAndRequires_ProvidedName](j2gort.MapGetStr(this_.providedNames, prefixNs)).AddProvide(node, chunk, false, this_.chunkGraph)
		} else {
			this_.priv_ProcessClosureProvidesAndRequires_putProvidedName_4467(prefixNs, NewProcessClosureProvidesAndRequires_ProvidedNameBuilder(this_).SetNamespace(prefixNs).SetNode(node).SetChunk(chunk).SetExplicit(false).Build())
		}
	}
}

// Java static method: ProcessClosureProvidesAndRequires.isNamespacePlaceholder(com.google.javascript.rhino.Node n)
func ProcessClosureProvidesAndRequires_IsNamespacePlaceholder(n pkg_rhino.Node) bool {
	clinit_ProcessClosureProvidesAndRequires()
	_ = n
	if !(n.GetBooleanProp(pkg_rhino.StaticGet_Node_IS_NAMESPACE())) {
		return false
	}
	var value pkg_rhino.Node = pkg_rhino.Node(0)
	_ = value
	if n.IsExprResult() {
		var assign pkg_rhino.Node = n.GetFirstChild()
		_ = assign
		value = assign.GetLastChild()
	} else if n.IsVar() {
		var name pkg_rhino.Node = n.GetFirstChild()
		_ = name
		value = name.GetFirstChild()
	}
	if value == 0 {
		return false
	}
	if value.IsCast() {
		value = value.GetOnlyChild()
	}
	return (value.IsObjectLit() && !(value.HasChildren()))
}

// Java class: com.google.javascript.jscomp.ProcessClosureProvidesAndRequires.CollectDefinitions
type ProcessClosureProvidesAndRequires_CollectDefinitions struct {
	outer_  *ProcessClosureProvidesAndRequires
	idHash_ atomic.Int32
}

func (this_ *ProcessClosureProvidesAndRequires_CollectDefinitions) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessClosureProvidesAndRequires_CollectDefinitions) IsJavaNil_() bool {
	return this_ == nil
}
func (*ProcessClosureProvidesAndRequires_CollectDefinitions) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessClosureProvidesAndRequires$CollectDefinitions", "CollectDefinitions"
}
func (this_ *ProcessClosureProvidesAndRequires_CollectDefinitions) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessClosureProvidesAndRequires_CollectDefinitions) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessClosureProvidesAndRequires.CollectDefinitions", "com.google.javascript.jscomp.ProcessClosureProvidesAndRequires$CollectDefinitions", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessClosureProvidesAndRequires_CollectDefinitions](v_)
	})
}
func (this_ *ProcessClosureProvidesAndRequires_CollectDefinitions) SetOuter_(o *ProcessClosureProvidesAndRequires) {
	this_.outer_ = o
}
func (this_ *ProcessClosureProvidesAndRequires_CollectDefinitions) GetOuter_() *ProcessClosureProvidesAndRequires {
	return this_.outer_
}
func (*ProcessClosureProvidesAndRequires_CollectDefinitions) Implements_jscomp_NodeTraversal_Callback() {
}

func (this_ *ProcessClosureProvidesAndRequires_CollectDefinitions) initFields_ProcessClosureProvidesAndRequires_CollectDefinitions() {
	_ = this_
}

// Java constructor: CollectDefinitions()
func NewProcessClosureProvidesAndRequires_CollectDefinitions(outer_ *ProcessClosureProvidesAndRequires) *ProcessClosureProvidesAndRequires_CollectDefinitions {
	this_ := &ProcessClosureProvidesAndRequires_CollectDefinitions{}
	this_.outer_ = outer_
	InitProcessClosureProvidesAndRequires_CollectDefinitions(this_)
	return this_
}

func InitProcessClosureProvidesAndRequires_CollectDefinitions(this_ *ProcessClosureProvidesAndRequires_CollectDefinitions) {
	this_.initFields_ProcessClosureProvidesAndRequires_CollectDefinitions()
}

// Java method: CollectDefinitions.shouldTraverse(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ProcessClosureProvidesAndRequires_CollectDefinitions) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	_ = t
	_ = n
	_ = parent
	if (n.IsModuleBody() && n.GetParent().GetBooleanProp(pkg_rhino.StaticGet_Node_GOOG_MODULE())) || NodeUtil_IsBundledGoogModuleScopeRoot(n) {
		var googModuleCall pkg_rhino.Node = n.GetFirstChild()
		_ = googModuleCall
		var closureNamespace string = googModuleCall.GetFirstChild().GetSecondChild().GetString()
		_ = closureNamespace
		var maybeLegacyNamespaceCall pkg_rhino.Node = googModuleCall.GetNext()
		_ = maybeLegacyNamespaceCall
		if (maybeLegacyNamespaceCall != 0) && NodeUtil_IsGoogModuleDeclareLegacyNamespaceCall(maybeLegacyNamespaceCall) {
			this_.outer_.priv_ProcessClosureProvidesAndRequires_processLegacyModuleCall_4475(closureNamespace, googModuleCall, t.GetChunk())
		}
		return false
	}
	return !(n.IsModuleBody())
}

// Java method: CollectDefinitions.visit(com.google.javascript.jscomp.NodeTraversal t, com.google.javascript.rhino.Node n, com.google.javascript.rhino.Node parent)
func (this_ *ProcessClosureProvidesAndRequires_CollectDefinitions) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = t
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(28) /* Token.CALL */ :
		var left pkg_rhino.Node = n.GetFirstChild()
		_ = left
		if left.IsGetProp() {
			var name pkg_rhino.Node = left.GetFirstChild()
			_ = name
			if name.MatchesName_1_string("goog") {
				this_.outer_.priv_ProcessClosureProvidesAndRequires_visitGoogMethodCall_4472(t, parent, n, left.GetString())
			}
		}
	case int32(51) /* Token.ASSIGN */, int32(32) /* Token.NAME */ :
		this_.outer_.priv_ProcessClosureProvidesAndRequires_handleCandidateProvideDefinition_4480(t, n, parent)
	case int32(93) /* Token.EXPR_RESULT */ :
		this_.outer_.priv_ProcessClosureProvidesAndRequires_handleStubDefinition_4477(t, n)
	default:
	}
}

// Java class: com.google.javascript.jscomp.ProcessClosureProvidesAndRequires.ProvidedNameBuilder
type ProcessClosureProvidesAndRequires_ProvidedNameBuilder struct {
	outer_                    *ProcessClosureProvidesAndRequires
	namespace                 string         `gson:"namespace"`
	chunk                     *JSChunk       `gson:"chunk"`
	node                      pkg_rhino.Node `gson:"node"`
	idHash_                   atomic.Int32
	explicit                  bool `gson:"explicit"`
	fromLegacyModule          bool `gson:"fromLegacyModule"`
	hasImplicitInitialization bool `gson:"hasImplicitInitialization"`
}

func (this_ *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) IsJavaNil_() bool {
	return this_ == nil
}
func (*ProcessClosureProvidesAndRequires_ProvidedNameBuilder) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessClosureProvidesAndRequires$ProvidedNameBuilder", "ProvidedNameBuilder"
}
func (this_ *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessClosureProvidesAndRequires_ProvidedNameBuilder) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessClosureProvidesAndRequires.ProvidedNameBuilder", "com.google.javascript.jscomp.ProcessClosureProvidesAndRequires$ProvidedNameBuilder", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessClosureProvidesAndRequires_ProvidedNameBuilder](v_)
	})
}
func (this_ *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) SetOuter_(o *ProcessClosureProvidesAndRequires) {
	this_.outer_ = o
}
func (this_ *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) GetOuter_() *ProcessClosureProvidesAndRequires {
	return this_.outer_
}

func (this_ *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) initFields_ProcessClosureProvidesAndRequires_ProvidedNameBuilder() {
	_ = this_
	this_.namespace = j2gort.NullString
}

// Java constructor: ProvidedNameBuilder()
func NewProcessClosureProvidesAndRequires_ProvidedNameBuilder(outer_ *ProcessClosureProvidesAndRequires) *ProcessClosureProvidesAndRequires_ProvidedNameBuilder {
	this_ := &ProcessClosureProvidesAndRequires_ProvidedNameBuilder{}
	this_.outer_ = outer_
	InitProcessClosureProvidesAndRequires_ProvidedNameBuilder(this_)
	return this_
}

func InitProcessClosureProvidesAndRequires_ProvidedNameBuilder(this_ *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) {
	this_.initFields_ProcessClosureProvidesAndRequires_ProvidedNameBuilder()
}

// Java method: ProvidedNameBuilder.setNamespace(java.lang.String namespace)
func (this_ *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) SetNamespace(namespace string) *ProcessClosureProvidesAndRequires_ProvidedNameBuilder {
	_ = namespace
	this_.namespace = namespace
	return this_
}

// Java method: ProvidedNameBuilder.setNode(com.google.javascript.rhino.@org.jspecify.annotations.Nullable Node node)
func (this_ *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) SetNode(node pkg_rhino.Node) *ProcessClosureProvidesAndRequires_ProvidedNameBuilder {
	_ = node
	this_.node = node
	return this_
}

// Java method: ProvidedNameBuilder.setChunk(com.google.javascript.jscomp.@org.jspecify.annotations.Nullable JSChunk chunk)
func (this_ *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) SetChunk(chunk *JSChunk) *ProcessClosureProvidesAndRequires_ProvidedNameBuilder {
	_ = chunk
	this_.chunk = chunk
	return this_
}

// Java method: ProvidedNameBuilder.setExplicit(boolean explicit)
func (this_ *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) SetExplicit(explicit bool) *ProcessClosureProvidesAndRequires_ProvidedNameBuilder {
	_ = explicit
	this_.explicit = explicit
	return this_
}

// Java method: ProvidedNameBuilder.setHasImplicitInitialization(boolean alreadyInitialized)
func (this_ *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) SetHasImplicitInitialization(alreadyInitialized bool) *ProcessClosureProvidesAndRequires_ProvidedNameBuilder {
	_ = alreadyInitialized
	this_.hasImplicitInitialization = alreadyInitialized
	return this_
}

// Java method: ProvidedNameBuilder.setFromLegacyModule(boolean fromLegacyModule)
func (this_ *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) SetFromLegacyModule(fromLegacyModule bool) *ProcessClosureProvidesAndRequires_ProvidedNameBuilder {
	_ = fromLegacyModule
	this_.fromLegacyModule = fromLegacyModule
	return this_
}

// Java method: ProvidedNameBuilder.build()
func (this_ *ProcessClosureProvidesAndRequires_ProvidedNameBuilder) Build() *ProcessClosureProvidesAndRequires_ProvidedName {
	return NewProcessClosureProvidesAndRequires_ProvidedName(this_, this_.outer_.chunkGraph, this_.outer_.compiler, this_.outer_.astFactory)
}

// Java class: com.google.javascript.jscomp.ProcessClosureProvidesAndRequires.ProvidedName
type ProcessClosureProvidesAndRequires_ProvidedName struct {
	namespace                 string           `gson:"namespace"`
	firstChunk                *JSChunk         `gson:"firstChunk"`
	minimumChunk              *JSChunk         `gson:"minimumChunk"`
	compiler                  AbstractCompiler `gson:"compiler"`
	astFactory                *AstFactory      `gson:"astFactory"`
	firstNode                 pkg_rhino.Node   `gson:"firstNode"`
	explicitNode              pkg_rhino.Node   `gson:"explicitNode"`
	candidateDefinition       pkg_rhino.Node   `gson:"candidateDefinition"`
	replacementNode           pkg_rhino.Node   `gson:"replacementNode"`
	idHash_                   atomic.Int32
	hasImplicitInitialization bool `gson:"hasImplicitInitialization"`
	fromLegacyModule          bool `gson:"fromLegacyModule"`
}

func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) BindSelf_(s any) {
	_ = s
}
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) IsJavaNil_() bool { return this_ == nil }
func (*ProcessClosureProvidesAndRequires_ProvidedName) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.ProcessClosureProvidesAndRequires$ProvidedName", "ProvidedName"
}
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*ProcessClosureProvidesAndRequires_ProvidedName) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.ProcessClosureProvidesAndRequires.ProvidedName", "com.google.javascript.jscomp.ProcessClosureProvidesAndRequires$ProvidedName", func(v_ any) bool {
		return j2gort.InstanceOf[*ProcessClosureProvidesAndRequires_ProvidedName](v_)
	})
}

func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) initFields_ProcessClosureProvidesAndRequires_ProvidedName() {
	_ = this_
	this_.namespace = j2gort.NullString

}

// Java constructor: ProvidedName(com.google.javascript.jscomp.ProcessClosureProvidesAndRequires.ProvidedNameBuilder builder, com.google.javascript.jscomp.JSChunkGraph chunkGraph, com.google.javascript.jscomp.AbstractCompiler compiler, com.google.javascript.jscomp.AstFactory astFactory)
func NewProcessClosureProvidesAndRequires_ProvidedName(builder *ProcessClosureProvidesAndRequires_ProvidedNameBuilder, chunkGraph *JSChunkGraph, compiler AbstractCompiler, astFactory *AstFactory) *ProcessClosureProvidesAndRequires_ProvidedName {
	this_ := &ProcessClosureProvidesAndRequires_ProvidedName{}
	InitProcessClosureProvidesAndRequires_ProvidedName(this_, builder, chunkGraph, compiler, astFactory)
	return this_
}

func InitProcessClosureProvidesAndRequires_ProvidedName(this_ *ProcessClosureProvidesAndRequires_ProvidedName, builder *ProcessClosureProvidesAndRequires_ProvidedNameBuilder, chunkGraph *JSChunkGraph, compiler AbstractCompiler, astFactory *AstFactory) {
	_ = builder
	_ = chunkGraph
	_ = compiler
	_ = astFactory
	this_.initFields_ProcessClosureProvidesAndRequires_ProvidedName()
	this_.compiler = compiler
	this_.astFactory = astFactory
	var node pkg_rhino.Node = builder.node
	_ = node
	if !(((node == 0) || NodeUtil_IsExprCall(node)) || (!(builder.explicit) && ((NodeUtil_IsExprAssign(node) || NodeUtil_IsNameDeclaration(node)) || (node.IsExprResult() && node.GetFirstChild().IsQualifiedName())))) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(node))
	}
	this_.namespace = builder.namespace
	this_.firstNode = builder.node
	this_.firstChunk = builder.chunk
	this_.fromLegacyModule = builder.fromLegacyModule
	this_.hasImplicitInitialization = builder.hasImplicitInitialization
	this_.AddProvide(node, builder.chunk, builder.explicit, chunkGraph)
}

// Java method: ProvidedName.addProvide(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.JSChunk chunk, boolean explicit, com.google.javascript.jscomp.JSChunkGraph chunkGraph)
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) AddProvide(node pkg_rhino.Node, chunk *JSChunk, explicit bool, chunkGraph *JSChunkGraph) {
	_ = node
	_ = chunk
	_ = explicit
	_ = chunkGraph
	if explicit {
		if !(this_.explicitNode == 0) {
			j2gort.CheckState(false)
		}
		if !(node.IsExprResult()) {
			j2gort.CheckArgumentMsg(false, j2gort.BoxNode(node))
		}
		this_.explicitNode = node
	}
	this_.priv_ProvidedName_updateMinimumChunk_4500(chunk, chunkGraph)
}

// Java method: ProvidedName.isExplicitlyProvided()
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) IsExplicitlyProvided() bool {
	return (this_.explicitNode != 0)
}

// Java method: ProvidedName.hasImplicitInitialization()
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) HasImplicitInitialization() bool {
	return this_.hasImplicitInitialization
}

// Java method: ProvidedName.isFromLegacyModule()
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) IsFromLegacyModule() bool {
	return this_.fromLegacyModule
}

// Java method: ProvidedName.hasCandidateDefinition()
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) priv_ProvidedName_hasCandidateDefinition_4495() bool {
	return (this_.candidateDefinition != 0)
}

// Java method: ProvidedName.getFirstProvideCall()
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) GetFirstProvideCall() pkg_rhino.Node {
	return this_.firstNode
}

// Java method: ProvidedName.getCandidateDefinition()
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) GetCandidateDefinition() pkg_rhino.Node {
	return this_.candidateDefinition
}

// Java method: ProvidedName.getNamespace()
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) GetNamespace() string {
	return this_.namespace
}

// Java method: ProvidedName.addDefinition(com.google.javascript.rhino.Node node, com.google.javascript.jscomp.JSChunk chunk, com.google.javascript.jscomp.JSChunkGraph chunkGraph)
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) priv_ProvidedName_addDefinition_4499(node pkg_rhino.Node, chunk *JSChunk, chunkGraph *JSChunkGraph) {
	_ = node
	_ = chunk
	_ = chunkGraph
	if !((node.IsExprResult() || node.IsFunction()) || NodeUtil_IsNameDeclaration(node)) {
		j2gort.CheckArgument(false)
	}
	if !(this_.explicitNode != node) {
		j2gort.CheckArgument(false)
	}
	if (this_.candidateDefinition == 0) || !(node.IsExprResult()) {
		this_.candidateDefinition = node
		this_.priv_ProvidedName_updateMinimumChunk_4500(chunk, chunkGraph)
	}
}

// Java method: ProvidedName.updateMinimumChunk(com.google.javascript.jscomp.JSChunk newChunk, com.google.javascript.jscomp.JSChunkGraph chunkGraph)
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) priv_ProvidedName_updateMinimumChunk_4500(newChunk *JSChunk, chunkGraph *JSChunkGraph) {
	_ = newChunk
	_ = chunkGraph
	if this_.minimumChunk == nil {
		this_.minimumChunk = newChunk
	} else if chunkGraph.GetChunkCount() > int32(1) {
		this_.minimumChunk = chunkGraph.GetDeepestCommonDependencyInclusive_2(this_.minimumChunk, newChunk)
	} else {
		if !(newChunk == this_.minimumChunk) {
			j2gort.CheckStateStr(false, "Missing chunk graph")
		}
	}
}

// Java method: ProvidedName.replace(boolean preserveGoogProvidesAndRequires, java.util.Map<java.lang.String,com.google.javascript.jscomp.ProcessClosureProvidesAndRequires.ProvidedName> providedNames)
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) priv_ProvidedName_replace_4501(preserveGoogProvidesAndRequires bool, providedNames j2gort.Ext_java_util_Map) {
	_ = preserveGoogProvidesAndRequires
	_ = providedNames
	if !(!(this_.IsFromLegacyModule())) {
		j2gort.CheckStateFmt(false, "Cannot rewrite provides without having rewritten goog.modules, found %s", j2gort.BoxNode(this_.firstNode))
	}
	if this_.firstNode == 0 {
		this_.replacementNode = this_.candidateDefinition
		return
	}
	if this_.priv_ProvidedName_hasCandidateDefinition_4495() && (this_.explicitNode != 0) {
		this_.replacementNode = this_.candidateDefinition
		if this_.candidateDefinition.IsExprResult() {
			var exprNode pkg_rhino.Node = this_.candidateDefinition.GetOnlyChild()
			_ = exprNode
			if exprNode.IsAssign() {
				var nameNode pkg_rhino.Node = exprNode.GetFirstChild()
				_ = nameNode
				if nameNode.IsName() {
					this_.priv_ProvidedName_convertProvideAssignmentToVarDeclaration_4502(exprNode, nameNode)
				} else {
					this_.candidateDefinition.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_NAMESPACE(), true)
				}
			}
		}
	} else {
		if !(this_.hasImplicitInitialization) {
			this_.priv_ProvidedName_createNamespaceInitialization_4503(this_.priv_ProvidedName_createDeclarationNode_4504(this_.astFactory.CreateObjectLit_1([]pkg_rhino.Node{})), providedNames)
		}
	}
	if this_.explicitNode != 0 {
		if preserveGoogProvidesAndRequires {
			return
		}
		this_.compiler.ReportChangeToEnclosingScope(this_.explicitNode)
		this_.explicitNode.Detach()
	}
}

// Java method: ProvidedName.convertProvideAssignmentToVarDeclaration(com.google.javascript.rhino.Node assignNode, com.google.javascript.rhino.Node nameNode)
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) priv_ProvidedName_convertProvideAssignmentToVarDeclaration_4502(assignNode pkg_rhino.Node, nameNode pkg_rhino.Node) {
	_ = assignNode
	_ = nameNode
	if !(assignNode.IsAssign()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(assignNode))
	}
	if !(nameNode.IsName()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(nameNode))
	}
	var valueNode pkg_rhino.Node = nameNode.GetNext()
	_ = valueNode
	nameNode.Detach()
	valueNode.Detach()
	var varNode pkg_rhino.Node = pkg_rhino.IR_Var_2(nameNode, valueNode).Srcref(this_.candidateDefinition)
	_ = varNode
	varNode.SetJSDocInfo(assignNode.GetJSDocInfo())
	varNode.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_NAMESPACE(), true)
	this_.candidateDefinition.ReplaceWith(varNode)
	this_.replacementNode = varNode
	this_.compiler.ReportChangeToEnclosingScope(varNode)
}

// Java method: ProvidedName.createNamespaceInitialization(com.google.javascript.rhino.Node replacement, java.util.Map<java.lang.String,com.google.javascript.jscomp.ProcessClosureProvidesAndRequires.ProvidedName> providedNames)
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) priv_ProvidedName_createNamespaceInitialization_4503(replacement pkg_rhino.Node, providedNames j2gort.Ext_java_util_Map) {
	_ = replacement
	_ = providedNames
	this_.replacementNode = replacement
	if this_.firstChunk == this_.minimumChunk {
		this_.replacementNode.InsertBefore(this_.firstNode)
	} else {
		var indexOfDot int32 = j2gort.JStrLastIndexOfCp(this_.namespace, int32(uint16(46)))
		_ = indexOfDot
		if indexOfDot == int32(-1) {
			this_.compiler.GetNodeForCodeInsertion(this_.minimumChunk).AddChildToBack(this_.replacementNode)
		} else {
			var parentName *ProcessClosureProvidesAndRequires_ProvidedName = j2gort.CastPtr[ProcessClosureProvidesAndRequires_ProvidedName](j2gort.MapGetStr(providedNames, j2gort.JStrSub(this_.namespace, int32(0), indexOfDot)))
			_ = parentName
			j2gort.CheckNotNullPtr(parentName)
			j2gort.CheckNotNullNode(parentName.replacementNode)
			this_.replacementNode.InsertAfter(parentName.replacementNode)
		}
	}
	this_.compiler.ReportChangeToEnclosingScope(this_.replacementNode)
}

// Java method: ProvidedName.createDeclarationNode(com.google.javascript.rhino.Node value)
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) priv_ProvidedName_createDeclarationNode_4504(value pkg_rhino.Node) pkg_rhino.Node {
	_ = value
	if !(value.IsObjectLit() || value.IsCast()) {
		j2gort.CheckArgumentMsg(false, j2gort.BoxNode(value))
	}
	if j2gort.JStrIndexOfCp(this_.namespace, int32(uint16(46))) == int32(-1) {
		return this_.priv_ProvidedName_makeVarDeclNode_4505(value)
	} else {
		return this_.priv_ProvidedName_makeAssignmentExprNode_4506(value)
	}
}

// Java method: ProvidedName.makeVarDeclNode(com.google.javascript.rhino.Node value)
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) priv_ProvidedName_makeVarDeclNode_4505(value pkg_rhino.Node) pkg_rhino.Node {
	_ = value
	var name pkg_rhino.Node = pkg_rhino.IR_Name(this_.namespace)
	_ = name
	name.AddChildToFront(value)
	var decl pkg_rhino.Node = pkg_rhino.IR_Var_1(name)
	_ = decl
	decl.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_NAMESPACE(), true)
	if this_.compiler.GetCodingConvention().IsConstant(this_.namespace) {
		name.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_CONSTANT_NAME(), true)
	}
	if !(this_.priv_ProvidedName_hasCandidateDefinition_4495()) {
		decl.SetJSDocInfo(NodeUtil_CreateConstantJsDoc())
	}
	if !(ProcessClosureProvidesAndRequires_IsNamespacePlaceholder(decl)) {
		j2gort.CheckState(false)
	}
	this_.priv_ProvidedName_setSourceInfo_4507(decl)
	return decl
}

// Java method: ProvidedName.makeAssignmentExprNode(com.google.javascript.rhino.Node value)
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) priv_ProvidedName_makeAssignmentExprNode_4506(value pkg_rhino.Node) pkg_rhino.Node {
	_ = value
	var lhs pkg_rhino.Node = this_.astFactory.CreateQNameWithUnknownType_1_string(this_.namespace).SrcrefTree(this_.firstNode)
	_ = lhs
	var decl pkg_rhino.Node = pkg_rhino.IR_ExprResult(this_.astFactory.CreateAssign_2_Node_Node(lhs, value))
	_ = decl
	decl.PutBooleanProp(pkg_rhino.StaticGet_Node_IS_NAMESPACE(), true)
	if !(this_.priv_ProvidedName_hasCandidateDefinition_4495()) {
		decl.GetFirstChild().SetJSDocInfo(NodeUtil_CreateConstantJsDoc())
	}
	if !(ProcessClosureProvidesAndRequires_IsNamespacePlaceholder(decl)) {
		j2gort.CheckState(false)
	}
	this_.priv_ProvidedName_setSourceInfo_4507(decl)
	lhs.GetFirstChild().MakeNonIndexableRecursive()
	return decl
}

// Java method: ProvidedName.setSourceInfo(com.google.javascript.rhino.Node newNode)
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) priv_ProvidedName_setSourceInfo_4507(newNode pkg_rhino.Node) {
	_ = newNode
	var sourceInfoNode pkg_rhino.Node = this_.firstNode
	_ = sourceInfoNode
	var provideStringNode pkg_rhino.Node = this_.priv_ProvidedName_getProvideStringNode_4508()
	_ = provideStringNode
	if provideStringNode != 0 {
		var firstCharIndex int32 = (j2gort.JStrLastIndexOfCp(this_.namespace, int32(uint16(46))) + int32(1))
		_ = firstCharIndex
		sourceInfoNode = provideStringNode.CloneNode_0()
		sourceInfoNode.SetLinenoCharno(sourceInfoNode.GetLineno(), ((sourceInfoNode.GetCharno() + firstCharIndex) + int32(1)))
		sourceInfoNode.SetLength((j2gort.JStrLen(this_.namespace) - firstCharIndex))
	}
	newNode.SrcrefTree(sourceInfoNode)
}

// Java method: ProvidedName.getProvideStringNode()
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) priv_ProvidedName_getProvideStringNode_4508() pkg_rhino.Node {
	if this_.firstNode.HasChildren() && NodeUtil_IsExprCall(this_.firstNode) {
		return this_.firstNode.GetFirstChild().GetLastChild()
	} else {
		return pkg_rhino.Node(0)
	}
}

// Java method: ProvidedName.toString()
func (this_ *ProcessClosureProvidesAndRequires_ProvidedName) String() string {
	var explicitOrImplicit string
	_ = explicitOrImplicit
	if this_.IsExplicitlyProvided() {
		explicitOrImplicit = "explicit"
	} else {
		explicitOrImplicit = "implicit"
	}
	return j2gort.Ext_java_lang_String_Format_2("ProvidedName: %s, %s", []any{j2gort.BoxString(this_.namespace), j2gort.BoxString(explicitOrImplicit)})
}
func (*ProcessClosureProvidesAndRequires_CollectDefinitions) CanReleaseSyntacticScopes_() {}
