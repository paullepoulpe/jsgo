// Port of com/google/javascript/jscomp/ExtractPrototypeMemberDeclarations.java @ closure-compiler v20261021

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

package jscomp

import (
	"slices"

	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// ExtractPrototypeMemberDeclarations: when there are multiple prototype
// member declarations to the same class, use a temp variable to alias the
// prototype object.
//
// Example:
//
//	function B() { ... }                 \
//	B.prototype.foo = function() { ... }  \___ ExtractionInstance
//	...                                   /
//	B.prototype.bar = function() { ... } /
//	         ^---------------------------------PrototypeMemberDeclaration
//
// becomes
//
//	function B() { ... }
//	x = B.prototype;
//	x.foo = function() { ... }
//	...
//	x.bar = function() { ... }
//
// Works almost like a redundant load elimination but limited to only
// recognizing the class prototype declaration idiom. First it only works
// within a basic block because we avoided DataFlowAnalysis for compilation
// performance. Secondly, we can avoid having to compute how long to
// sub-expressing has to be. Example:
//
//	a.b.c.d = ...
//	a.b.c = ...
//	a.b = ...
//	a.b.c = ...
//
// Further more, we only introduce one temp variable to hold a single
// prototype at a time. So all the PrototypeMemberDeclaration to be extracted
// must be in a single line. We call this a single ExtractionInstance.
//
// Alternatively, for users who do not want a global variable to be
// introduced, we will create an anonymous function instead.
//
//	function B() { ... }
//	(function (x) {
//	  x.foo = function() { ... }
//	  ...
//	  x.bar = function() { ... }
//	)(B.prototype)
//
// The RHS of the declarations can have side effects, however, one good way
// to break this is the following:
//
//	function B() { ... }
//	B.prototype.foo = (function() { B.prototype = somethingElse(); return 0 })();
//	...
//
// Such logic is highly unlikely and we will assume that it never occurs.
type ExtractPrototypeMemberDeclarations struct {
	compiler AbstractCompiler

	pattern ExtractPrototypeMemberDeclarations_Pattern
}

// The name of variable that will temporary hold the pointer to the prototype
// object. Of course, we assume that it'll be renamed by RenameVars.
const ExtractPrototypeMemberDeclarations_PROTOTYPE_ALIAS = "JSCompiler_prototypeAlias"

// ExtractPrototypeMemberDeclarations_Pattern is the nested enum Pattern.
type ExtractPrototypeMemberDeclarations_Pattern int32

const (
	ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP ExtractPrototypeMemberDeclarations_Pattern = iota
	ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP
	ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE
)

var extractPrototypeMemberDeclarations_Pattern_names = [...]string{"USE_GLOBAL_TEMP", "USE_CHUNK_TEMP", "USE_IIFE"}

// The enum fields globalOverhead, perExtractionOverhead and
// perMemberOverhead, indexed by ordinal.
var (
	extractPrototypeMemberDeclarations_Pattern_globalOverhead = [...]int{
		// USE_GLOBAL_TEMP: Global Overhead.
		// We need a temp variable to hold all the prototype.
		len("var t;"),
		// USE_CHUNK_TEMP: Per Chunk Overhead.
		// We need a temp variable to hold all the prototype.
		len("var t;"),
		// USE_IIFE: Global Overhead:
		0,
	}
	extractPrototypeMemberDeclarations_Pattern_perExtractionOverhead = [...]int{
		// USE_GLOBAL_TEMP: Per Extract overhead:
		// Every extraction instance must first use the temp variable to point
		// to the prototype object.
		len("t=y.prototype;"),
		// USE_CHUNK_TEMP: Per Extract overhead:
		// Every extraction instance must first use the temp variable to point
		// to the prototype object.
		len("t=y.prototype;"),
		// USE_IIFE: Per-extraction overhead:
		// This is the cost of a single anoynmous function.
		len("(function(t){})(y.prototype);"),
	}
	extractPrototypeMemberDeclarations_Pattern_perMemberOverhead = [...]int{
		// USE_GLOBAL_TEMP:
		// TODO(user): Check to to see if AliasExterns is on
		// The gain we get per prototype declaration. Assuming it can be
		// aliased.
		len("t.y=") - len("x[p].y="),
		// USE_CHUNK_TEMP:
		// TODO(user): Check to to see if AliasExterns is on
		// The gain we get per prototype declaration. Assuming it can be
		// aliased.
		len("t.y=") - len("x[p].y="),
		// USE_IIFE: Per-prototype member declaration overhead:
		// Here we assumes that they don't have AliasExterns on (in SIMPLE mode).
		len("t.y=") - len("x.prototype.y="),
	}
)

// ExtractPrototypeMemberDeclarations_Pattern_values is Pattern.values().
func ExtractPrototypeMemberDeclarations_Pattern_values() []ExtractPrototypeMemberDeclarations_Pattern {
	return []ExtractPrototypeMemberDeclarations_Pattern{
		ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP,
		ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP,
		ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE,
	}
}

// String is Enum.name().
func (p ExtractPrototypeMemberDeclarations_Pattern) String() string {
	return extractPrototypeMemberDeclarations_Pattern_names[p]
}

// Java: ExtractPrototypeMemberDeclarations(AbstractCompiler, Pattern)
func NewExtractPrototypeMemberDeclarations(compiler AbstractCompiler, pattern ExtractPrototypeMemberDeclarations_Pattern) *ExtractPrototypeMemberDeclarations {
	return &ExtractPrototypeMemberDeclarations{compiler: compiler, pattern: pattern}
}

// Java: ExtractPrototypeMemberDeclarations.process(Node, Node)
func (e *ExtractPrototypeMemberDeclarations) Process(externs rhino.Node, root rhino.Node) {
	extractionInfo := NewExtractPrototypeMemberDeclarations_GatherExtractionInfo(e)
	NodeTraversal_Traverse(e.compiler, root, extractionInfo)
	e.MaybeDoExtraction(extractionInfo)
}

// MaybeDoExtraction declares the temp variable to point to prototype objects
// and iterates through all ExtractInstance and performs extraction there.
//
// Java: ExtractPrototypeMemberDeclarations.maybeDoExtraction(GatherExtractionInfo)
func (e *ExtractPrototypeMemberDeclarations) MaybeDoExtraction(info *ExtractPrototypeMemberDeclarations_GatherExtractionInfo) {
	if (e.pattern == ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE || e.pattern == ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP) &&
		!info.ShouldExtractGlobal() {
		return
	}
	if e.pattern == ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP {
		injectionPoint := e.compiler.GetNodeForCodeInsertion(nil)

		var_ := NodeUtil_NewVarNode(ExtractPrototypeMemberDeclarations_PROTOTYPE_ALIAS, rhino.NilNode).SrcrefTreeIfMissing(injectionPoint)

		injectionPoint.AddChildToFront(var_)
		e.compiler.ReportChangeToEnclosingScope(var_)
	}
	// Go through all extraction instances and extract each of them.
	for entryKey, entryValue := range info.instancesByChunk.All() {
		alias := ExtractPrototypeMemberDeclarations_PROTOTYPE_ALIAS
		if e.pattern == ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP {
			// Rather than a truly global variable, use a unique variable per output chunk.
			// This prevents RescopeGlobalSymbolNames from converting these references to
			// namespace properties which reduces the benefit of the alias.
			if info.ShouldExtractChunk(entryKey) {
				injectionPoint := e.compiler.GetNodeForCodeInsertion(entryKey)
				alias = ExtractPrototypeMemberDeclarations_PROTOTYPE_ALIAS + jutil.ToJavaString(entryKey.GetIndex())
				var_ := NodeUtil_NewVarNode(alias, rhino.NilNode).SrcrefTreeIfMissing(injectionPoint)

				injectionPoint.AddChildToFront(var_)
				e.compiler.ReportChangeToEnclosingScope(var_)
			} else {
				continue
			}
		}

		for _, instance := range entryValue.instances {
			e.ExtractInstance(instance, alias)
		}
	}
}

// ExtractInstance, at a given ExtractionInstance, stores and prototype
// object in the temp variable and rewrite each member declaration to assign
// to the temp variable instead.
//
// Java: ExtractPrototypeMemberDeclarations.extractInstance(ExtractionInstance, String)
func (e *ExtractPrototypeMemberDeclarations) ExtractInstance(instance *ExtractPrototypeMemberDeclarations_ExtractionInstance, alias string) {
	first := &instance.declarations[0]
	className, _ := first.classNode.GetQualifiedNameString()
	if e.pattern == ExtractPrototypeMemberDeclarations_Pattern_USE_GLOBAL_TEMP || e.pattern == ExtractPrototypeMemberDeclarations_Pattern_USE_CHUNK_TEMP {
		// Use the temp variable to hold the prototype.
		classNameNode := NodeUtil_NewQName(e.compiler, className)
		classNameNode.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, first.constant)
		stmt := rhino.IR_ExprResult(rhino.IR_Assign(rhino.IR_Name(alias), rhino.IR_Getprop(classNameNode, "prototype"))).
			SrcrefTreeIfMissing(first.node)

		stmt.InsertBefore(first.node)
		e.compiler.ReportChangeToEnclosingScope(stmt)
	} else if e.pattern == ExtractPrototypeMemberDeclarations_Pattern_USE_IIFE {
		block := rhino.IR_Block()
		func_ := rhino.IR_Function(rhino.IR_Name(""), rhino.IR_ParamList(rhino.IR_Name(alias)), block)

		protoName := className + ".prototype"
		call := rhino.IR_Call(
			func_,
			NodeUtil_NewQNameWithOriginalName(
				e.compiler, protoName, instance.parent, protoName))
		call.PutIntProp(rhino.Node_FREE_CALL, 1)

		stmt := rhino.IR_ExprResult(call)
		stmt.SrcrefTreeIfMissing(first.node)
		stmt.InsertBefore(first.node)
		e.compiler.ReportChangeToEnclosingScope(stmt)
		for i := range instance.declarations {
			declar := &instance.declarations[i]
			e.compiler.ReportChangeToEnclosingScope(declar.node)
			block.AddChildToBack(declar.node.Detach())
		}
	}
	// Go through each member declaration and replace it with an assignment
	// to the prototype variable.
	for i := range instance.declarations {
		e.ReplacePrototypeMemberDeclaration(&instance.declarations[i], alias)
	}
}

// ReplacePrototypeMemberDeclaration replaces a member declaration to an
// assignment to the temp prototype object.
//
// Java: ExtractPrototypeMemberDeclarations.replacePrototypeMemberDeclaration(PrototypeMemberDeclaration, String)
func (e *ExtractPrototypeMemberDeclarations) ReplacePrototypeMemberDeclaration(declar *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration, alias string) {
	// x.prototype.y = ...  ->  t.y = ...
	assignment := declar.node.GetFirstChild()
	lhs := assignment.GetFirstChild()
	name := NodeUtil_NewQNameWithOriginalName(
		e.compiler, alias+"."+declar.memberName, declar.node, declar.memberName)

	// Save the full prototype path on the left hand side of the assignment for debugging purposes.
	// declar.lhs = x.prototype.y so first child of the first child is 'x'.
	accessNode := declar.lhs.GetFirstFirstChild()
	originalName := accessNode.GetOriginalName()
	className := "?"
	if originalName != nil {
		className = *originalName
	}
	name.GetFirstChild().SrcrefTree(lhs)
	name.PutBooleanProp(rhino.Node_IS_CONSTANT_NAME, lhs.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME))
	newOriginalName := className + ".prototype"
	name.GetFirstChild().SetOriginalName(&newOriginalName)

	lhs.ReplaceWith(name)
	e.compiler.ReportChangeToEnclosingScope(name)
}

// ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo is per-chunk info
// needed for prototype extraction
type ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo struct {
	instances  []*ExtractPrototypeMemberDeclarations_ExtractionInstance
	totalDelta int
}

// Java: ExtractPrototypeMemberDeclarations.ExtractionInstanceInfo()
func NewExtractPrototypeMemberDeclarations_ExtractionInstanceInfo() *ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo {
	return &ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo{}
}

// ExtractPrototypeMemberDeclarations_GatherExtractionInfo collects all the
// possible extraction instances in a node traversal.
type ExtractPrototypeMemberDeclarations_GatherExtractionInfo struct {
	NodeTraversal_AbstractShallowCallbackBase
	outer            *ExtractPrototypeMemberDeclarations
	instancesByChunk *jutil.LinkedHashMap[*JSChunk, *ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo]
	scratchDecls     []ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration
}

// Java: ExtractPrototypeMemberDeclarations.GatherExtractionInfo()
func NewExtractPrototypeMemberDeclarations_GatherExtractionInfo(outer *ExtractPrototypeMemberDeclarations) *ExtractPrototypeMemberDeclarations_GatherExtractionInfo {
	return &ExtractPrototypeMemberDeclarations_GatherExtractionInfo{
		outer:            outer,
		instancesByChunk: jutil.NewLinkedHashMap[*JSChunk, *ExtractPrototypeMemberDeclarations_ExtractionInstanceInfo](),
	}
}

// Java: ExtractPrototypeMemberDeclarations.GatherExtractionInfo.visit(NodeTraversal, Node, Node)
func (g *ExtractPrototypeMemberDeclarations_GatherExtractionInfo) Visit(t *NodeTraversal, n rhino.Node, parent rhino.Node) {

	if !n.IsScript() && !n.IsBlock() {
		return
	}

	for cur := n.GetFirstChild(); !cur.IsNil(); cur = cur.GetNext() {
		if !ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration_IsPrototypePropertyDeclaration(cur) {
			continue
		}
		g.scratchDecls = g.scratchDecls[:0]
		head := makePrototypeMemberDeclaration(cur.GetFirstFirstChild(), cur)
		g.scratchDecls = append(g.scratchDecls, head)
		delta := extractPrototypeMemberDeclarations_Pattern_perExtractionOverhead[g.outer.pattern] +
			extractPrototypeMemberDeclarations_Pattern_perMemberOverhead[g.outer.pattern]

		for next := head.node.GetNext(); !next.IsNil(); next = next.GetNext() {
			if next.IsFunction() {
				continue
			}
			if !ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration_IsPrototypePropertyDeclaration(next) {
				break
			}
			cand := makePrototypeMemberDeclaration(next.GetFirstFirstChild(), next)
			if !head.IsSameClass(&cand) {
				break
			}
			g.scratchDecls = append(g.scratchDecls, cand)
			delta += extractPrototypeMemberDeclarations_Pattern_perMemberOverhead[g.outer.pattern]
		}
		cur = g.scratchDecls[len(g.scratchDecls)-1].node

		// Only add it to our work list if the extraction at this instance makes the code smaller.
		if delta <= 0 {
			instance := &ExtractPrototypeMemberDeclarations_ExtractionInstance{
				declarations: slices.Clone(g.scratchDecls),
				delta:        delta,
				parent:       n,
			}

			chunk := t.GetChunk()
			instanceInfo := g.instancesByChunk.Get(chunk)
			if instanceInfo == nil {
				instanceInfo = NewExtractPrototypeMemberDeclarations_ExtractionInstanceInfo()
				g.instancesByChunk.PutNew(chunk, instanceInfo)
			}
			instanceInfo.instances = append(instanceInfo.instances, instance)
			instanceInfo.totalDelta += instance.delta
		}
	}
}

// ShouldExtractGlobal returns true if the sum of all the extraction instance
// gain outweighs the overhead of the temp variable declaration.
//
// Java: ExtractPrototypeMemberDeclarations.GatherExtractionInfo.shouldExtractGlobal()
func (g *ExtractPrototypeMemberDeclarations_GatherExtractionInfo) ShouldExtractGlobal() bool {
	allModulesDelta := 0
	for instanceInfo := range g.instancesByChunk.Values() {
		allModulesDelta += instanceInfo.totalDelta
	}
	return allModulesDelta+extractPrototypeMemberDeclarations_Pattern_globalOverhead[g.outer.pattern] < 0
}

// ShouldExtractChunk returns true if the sum of all the extraction instance
// gain outweighs the overhead of the temp variable declaration.
//
// Java: ExtractPrototypeMemberDeclarations.GatherExtractionInfo.shouldExtractChunk(JSChunk)
func (g *ExtractPrototypeMemberDeclarations_GatherExtractionInfo) ShouldExtractChunk(chunk *JSChunk) bool {
	instanceInfo := g.instancesByChunk.Get(chunk)
	if instanceInfo == nil {
		return false
	}
	return instanceInfo.totalDelta+extractPrototypeMemberDeclarations_Pattern_globalOverhead[g.outer.pattern] < 0
}

// ExtractPrototypeMemberDeclarations_ExtractionInstance is the inner class
// ExtractionInstance.
type ExtractPrototypeMemberDeclarations_ExtractionInstance struct {
	declarations []ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration
	delta        int
	parent       rhino.Node
}

// Java: ExtractPrototypeMemberDeclarations.ExtractionInstance(PrototypeMemberDeclaration, Node)
func NewExtractPrototypeMemberDeclarations_ExtractionInstance(
	outer *ExtractPrototypeMemberDeclarations, head *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration, parent rhino.Node) *ExtractPrototypeMemberDeclarations_ExtractionInstance {
	i := &ExtractPrototypeMemberDeclarations_ExtractionInstance{
		parent: parent,
	}
	i.declarations = append(i.declarations, *head)
	i.delta = extractPrototypeMemberDeclarations_Pattern_perExtractionOverhead[outer.pattern] +
		extractPrototypeMemberDeclarations_Pattern_perMemberOverhead[outer.pattern]

	for cur := head.node.GetNext(); !cur.IsNil(); cur = cur.GetNext() {

		// We can skip over any named functions because they have no effect on
		// the control flow. In fact, they are lifted to the beginning of the
		// block. This happens a lot when devirtualization breaks the whole chain.
		if cur.IsFunction() {
			continue
		}

		prototypeMember := ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration_ExtractDeclaration(cur)
		if prototypeMember == nil || !head.IsSameClass(prototypeMember) {
			break
		}
		i.declarations = append(i.declarations, *prototypeMember)
		i.delta += extractPrototypeMemberDeclarations_Pattern_perMemberOverhead[outer.pattern]
	}
	return i
}

// IsFavorable returns true if extracting all the declarations at this
// instance will overweight the overhead of aliasing the prototype object.
//
// Java: ExtractPrototypeMemberDeclarations.ExtractionInstance.isFavorable()
func (i *ExtractPrototypeMemberDeclarations_ExtractionInstance) IsFavorable() bool {
	return i.delta <= 0
}

// ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration is the
// abstraction for a prototype member declaration.
//
// `a.b.c.prototype.d = ....`
type ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration struct {
	memberName string
	node       rhino.Node
	classNode  rhino.Node
	lhs        rhino.Node
	constant   bool
}

func makePrototypeMemberDeclaration(lhs rhino.Node, node rhino.Node) ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration {
	jutil.CheckStateMsg1(NodeUtil_IsExprAssign(node), "%s", node)
	classNode := ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration_GetPrototypeClassName(lhs)
	return ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration{
		memberName: NodeUtil_GetPrototypePropertyName(lhs),
		node:       node,
		classNode:  classNode,
		lhs:        lhs,
		constant:   classNode.GetBooleanProp(rhino.Node_IS_CONSTANT_NAME),
	}
}

// Java: ExtractPrototypeMemberDeclarations.PrototypeMemberDeclaration(Node, Node)
func NewExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration(lhs rhino.Node, node rhino.Node) *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration {
	d := makePrototypeMemberDeclaration(lhs, node)
	return &d
}

// Java: ExtractPrototypeMemberDeclarations.PrototypeMemberDeclaration.isSameClass(PrototypeMemberDeclaration)
func (d *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration) IsSameClass(other *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration) bool {
	// String.equals(Object): other's name may be null.
	return !d.classNode.IsNil() && !other.classNode.IsNil() && d.classNode.MatchesQualifiedNameNode(other.classNode)
}

// Java: ExtractPrototypeMemberDeclarations.PrototypeMemberDeclaration.getPrototypeClassName(Node)
func ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration_GetPrototypeClassName(qName rhino.Node) rhino.Node {
	cur := qName
	for cur.IsGetProp() {
		if cur.GetString() == "prototype" {
			return cur.GetFirstChild()
		} else {
			cur = cur.GetFirstChild()
		}
	}
	return rhino.NilNode
}

// Java: ExtractPrototypeMemberDeclarations.PrototypeMemberDeclaration.isPrototypePropertyDeclaration(Node)
func ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration_IsPrototypePropertyDeclaration(n rhino.Node) bool {
	if !NodeUtil_IsExprAssign(n) {
		return false
	}
	lvalue := n.GetFirstFirstChild()
	if lvalue.IsGetProp() {
		cur := lvalue.GetFirstChild()
		for cur.IsGetProp() {
			if cur.GetString() == "prototype" {
				return cur.IsQualifiedName()
			}
			cur = cur.GetFirstChild()
		}
	}
	return false
}

// ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration_ExtractDeclaration
// returns a prototype member declaration representation if there is one else
// it returns null.
//
// Java: ExtractPrototypeMemberDeclarations.PrototypeMemberDeclaration.extractDeclaration(Node)
func ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration_ExtractDeclaration(n rhino.Node) *ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration {
	if !ExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration_IsPrototypePropertyDeclaration(n) {
		return nil
	}
	lhs := n.GetFirstFirstChild()
	return NewExtractPrototypeMemberDeclarations_PrototypeMemberDeclaration(lhs, n)
}
