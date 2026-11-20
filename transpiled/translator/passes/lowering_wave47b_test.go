package passes

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const wave47bNodestoreStub = `package nodestore

import "sync/atomic"

type Node uint32

const BlockShift = 14
const blockLen = 1 << BlockShift
const blockMask = blockLen - 1

type NodeData struct {
	Next   Node
	First  Node
	Parent Node
	Token  uint8
}

type NodeBlock [blockLen]NodeData

var nodeBlocks [64]atomic.Pointer[NodeBlock]

var BlockOfCalls atomic.Int64
var FirstNextTokNodeCalls atomic.Int64

func ResetCounters() {
	BlockOfCalls.Store(0)
	FirstNextTokNodeCalls.Store(0)
}

func EnsureBlock(blkIdx uint32) *NodeBlock {
	if blkIdx == 0 {
		panic("block 0 is reserved for nil/unmapped NPE")
	}
	if b := nodeBlocks[blkIdx].Load(); b != nil {
		return b
	}
	b := new(NodeBlock)
	nodeBlocks[blkIdx].Store(b)
	return b
}

func MakeNode(blkIdx uint32, slot uint32, tok uint8) Node {
	b := EnsureBlock(blkIdx)
	n := Node((blkIdx << BlockShift) | (slot & blockMask))
	b[slot&blockMask] = NodeData{Token: tok}
	return n
}

func (n Node) Rec() *NodeData {
	return &nodeBlocks[n>>BlockShift].Load()[n&blockMask]
}

func BlockOf(n Node) *NodeBlock {
	BlockOfCalls.Add(1)
	b := nodeBlocks[n>>BlockShift].Load()
	_ = &b[0] // immediate runtime.Error nil deref on unmapped/0 handle
	return b
}

func (b *NodeBlock) FirstNextTok(n Node) (first, next Node, tok uint8) {
	rec := &b[n&blockMask]
	return rec.First, rec.Next, rec.Token
}

func (n Node) FirstNextTok() (first, next Node, tok uint8) {
	FirstNextTokNodeCalls.Add(1)
	rec := n.Rec()
	return rec.First, rec.Next, rec.Token
}

func (n Node) FirstChild() Node {
	return n.Rec().First
}

func (n Node) SetFirstChild(c Node) {
	n.Rec().First = c
}

func (n Node) Next() Node {
	return n.Rec().Next
}

func (n Node) SetNext(nx Node) {
	n.Rec().Next = nx
}

func (n Node) Parent() Node {
	return n.Rec().Parent
}

func (n Node) SetParent(p Node) {
	n.Rec().Parent = p
}

func AddChildToBack(parent, child Node) {
	child.SetParent(parent)
	child.SetNext(0)
	if parent.FirstChild() == 0 {
		parent.SetFirstChild(child)
		return
	}
	cur := parent.FirstChild()
	for cur.Next() != 0 {
		cur = cur.Next()
	}
	cur.SetNext(child)
}

func AddChildToFront(parent, child Node) {
	child.SetParent(parent)
	child.SetNext(parent.FirstChild())
	parent.SetFirstChild(child)
}

func RemoveChildren(parent Node) {
	parent.SetFirstChild(0)
}
`

const wave47bRhinoStub = `package rhino

import "example.com/w47btest/shims/nodestore"

type Node nodestore.Node

const (
	TokSCRIPT            int32 = 94
	TokFUNCTION          int32 = 73
	TokMODULE_BODY       int32 = 123
	TokCLASS             int32 = 105
	TokCLASS_MEMBERS     int32 = 106
	TokDEFAULT_VALUE     int32 = 134
	TokDESTRUCTURING_LHS int32 = 104

	TokBLOCK        int32 = 90
	TokFOR          int32 = 80
	TokFOR_IN       int32 = 81
	TokFOR_OF       int32 = 112
	TokFOR_AWAIT_OF int32 = 113
	TokSWITCH_BODY  int32 = 91

	TokEXPR_RESULT int32 = 30
	TokADD         int32 = 31
	TokNAME        int32 = 32
	TokNUMBER      int32 = 33
	TokCATCH       int32 = 40
	TokCASE        int32 = 41
)

type NullTokenPanic struct{}

var nodeNullPtrEx any = NullTokenPanic{}

func (this_ Node) TokenSwitchOrd_() int32 {
	t := nodestore.Node(this_).Rec().Token
	if t == 0xFF {
		panic(nodeNullPtrEx)
	}
	return int32(t)
}

func (this_ Node) FirstNextTok_() (Node, Node, uint8) {
	f, nx, tok := nodestore.Node(this_).FirstNextTok()
	return Node(f), Node(nx), tok
}

func (this_ Node) GetFirstChild() Node {
	return Node(nodestore.Node(this_).FirstChild())
}

func (this_ Node) GetNext() Node {
	return Node(nodestore.Node(this_).Next())
}

func (this_ Node) GetParent() Node {
	return Node(nodestore.Node(this_).Parent())
}

func (this_ Node) IsSyntheticBlock() bool { return false }
func (this_ Node) IsCatch() bool          { return nodestore.Node(this_).Rec().Token == uint8(TokCATCH) }
`

const wave47bJ2gortStub = `package j2gort

func CheckStateMsg(cond bool, msg any) {
	if !cond {
		panic(msg)
	}
}

func BoxNode(n any) any { return n }
`

const wave47bJscompCommon = `package jscomp

import pkg_rhino "example.com/w47btest/gen/rhino"

type AbstractCompiler interface{}

type ScopeCreator interface {
	IsJavaNil_() bool
	CreateScope(pkg_rhino.Node, AbstractScope) AbstractScope
}

type AbstractScope interface{}

type SyntacticScopeCreator_Scanner struct {
	compiler                      AbstractCompiler
	redeclarationHandler          any
	treatProvidesAsRedeclarations bool
}

func (s *SyntacticScopeCreator_Scanner) Reset() {}

type SyntacticScopeCreator struct {
	compiler                      AbstractCompiler
	redeclarationHandler          any
	treatProvidesAsRedeclarations bool
	scannerInUse                  bool
	reusableScanner               *SyntacticScopeCreator_Scanner
}

var Static_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER any

func StaticGet_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER() any {
	return Static_SyntacticScopeCreator_DEFAULT_REDECLARATION_HANDLER
}

func InitSyntacticScopeCreator_1(this_ *SyntacticScopeCreator, c AbstractCompiler) *SyntacticScopeCreator {
	this_.compiler = c
	return this_
}

func NewSyntacticScopeCreator_1(c AbstractCompiler) *SyntacticScopeCreator {
	return InitSyntacticScopeCreator_1(&SyntacticScopeCreator{}, c)
}

func (s *SyntacticScopeCreator) IsJavaNil_() bool { return s == nil }

func (s *SyntacticScopeCreator) CreateScope(n pkg_rhino.Node, parent AbstractScope) AbstractScope {
	return n
}

type NodeTraversal_Callback interface {
	ShouldTraverse_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node) bool
	Visit_3(*NodeTraversal, pkg_rhino.Node, pkg_rhino.Node)
}

type NodeTraversal_ScopedCallback interface {
	NodeTraversal_Callback
	EnterScope(*NodeTraversal)
	ExitScope(*NodeTraversal)
}

type NodeTraversal_Builder struct {
	compiler                                       AbstractCompiler
	callback                                       NodeTraversal_Callback
	scopeCreator                                   ScopeCreator
	obeyDestructuringAndDefaultValueExecutionOrder bool
}

func NodeTraversal_S_Builder() *NodeTraversal_Builder {
	return &NodeTraversal_Builder{}
}

func (b *NodeTraversal_Builder) SetCompiler(c AbstractCompiler) *NodeTraversal_Builder {
	b.compiler = c
	return b
}

func (b *NodeTraversal_Builder) SetCallback_1_NodeTraversal_Callback(cb NodeTraversal_Callback) *NodeTraversal_Builder {
	b.callback = cb
	return b
}

func (b *NodeTraversal_Builder) SetScopeCreator(sc ScopeCreator) *NodeTraversal_Builder {
	b.scopeCreator = sc
	return b
}

func (b *NodeTraversal_Builder) Build() *NodeTraversal {
	return NewNodeTraversal(b)
}

func NodeUtil_IsClassStaticBlock(n pkg_rhino.Node) bool  { return false }
func NodeUtil_IsSwitchCase(n pkg_rhino.Node) bool        { return false }
func NodeUtil_IsFunctionDeclaration(n pkg_rhino.Node) bool { return false }

// PostOrderCallback always returns true in ShouldTraverse_3 -> gets AlwaysShouldTraverse_() marker!
type PostOrderRecorder struct {
	Events []string
}

func (r *PostOrderRecorder) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	return true
}

func (r *PostOrderRecorder) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	if t.currentNode != n {
		panic("currentNode mismatch in Visit_3")
	}
	r.Events = append(r.Events, fmtEvent("V", n, parent))
}

// ScopedPostOrderCallback has AlwaysShouldTraverse_() + EnterScope mutation hook!
type ScopedPostOrderMutator struct {
	Events       []string
	OnEnterScope func(t *NodeTraversal, scopeRoot pkg_rhino.Node)
}

func (r *ScopedPostOrderMutator) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	return true
}

func (r *ScopedPostOrderMutator) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	if t.currentNode != n {
		panic("currentNode mismatch in Visit_3")
	}
	r.Events = append(r.Events, fmtEvent("V", n, parent))
}

func (r *ScopedPostOrderMutator) EnterScope(t *NodeTraversal) {
	root := t.GetScopeRoot()
	r.Events = append(r.Events, fmtEvent("ES", root, 0))
	if r.OnEnterScope != nil {
		r.OnEnterScope(t, root)
	}
}

func (r *ScopedPostOrderMutator) ExitScope(t *NodeTraversal) {
	root := t.GetScopeRoot()
	r.Events = append(r.Events, fmtEvent("XS", root, 0))
}

// GeneralMutatingCallback has non-trivial ShouldTraverse_3 (!alwaysTraverse_).
type GeneralMutatingCallback struct {
	Events    []string
	OnShould  func(t *NodeTraversal, n, parent pkg_rhino.Node) bool
}

func (r *GeneralMutatingCallback) ShouldTraverse_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) bool {
	if t.currentNode != n {
		panic("currentNode mismatch in ShouldTraverse_3")
	}
	r.Events = append(r.Events, fmtEvent("ST", n, parent))
	if r.OnShould != nil {
		return r.OnShould(t, n, parent)
	}
	return true
}

func (r *GeneralMutatingCallback) Visit_3(t *NodeTraversal, n pkg_rhino.Node, parent pkg_rhino.Node) {
	if t.currentNode != n {
		panic("currentNode mismatch in Visit_3")
	}
	r.Events = append(r.Events, fmtEvent("V", n, parent))
}

func fmtEvent(kind string, n, p pkg_rhino.Node) string {
	return kind + ":" + itoa32(uint32(n)) + ":" + itoa32(uint32(p))
}

func itoa32(v uint32) string {
	if v == 0 {
		return "0"
	}
	var buf [10]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	return string(buf[i:])
}
`

const wave47bNodeTraversalFresh = `package jscomp

import (
	"example.com/w47btest/gen/j2gort"
	pkg_rhino "example.com/w47btest/gen/rhino"
)

func clinit_NodeTraversal() {}

type NodeTraversal struct {
	compiler                                        AbstractCompiler
	callback                                        NodeTraversal_Callback
	scopeCallback                                   NodeTraversal_ScopedCallback
	scopeCreator                                    ScopeCreator
	currentNode                                     pkg_rhino.Node
	currentHoistScopeRoot                           pkg_rhino.Node
	currentFunction                                 pkg_rhino.Node
	scopeNodes                                      []pkg_rhino.Node
	scopeObjs                                       []AbstractScope
	scopeDepth                                      int32
	mayContainSyntheticBlocks                       bool
	obeyDestructuringAndDefaultValueExecutionOrder  bool
}

func InitNodeTraversal(this_ *NodeTraversal, builder *NodeTraversal_Builder) {
	this_.compiler = builder.compiler
	this_.callback = builder.callback
	if scb, ok := builder.callback.(NodeTraversal_ScopedCallback); ok {
		this_.scopeCallback = scb
	}
	if builder.scopeCreator == nil || builder.scopeCreator.IsJavaNil_() {
		this_.scopeCreator = NewSyntacticScopeCreator_1(this_.compiler)
	} else {
		this_.scopeCreator = builder.scopeCreator
	}
}

func NewNodeTraversal(builder *NodeTraversal_Builder) *NodeTraversal {
	t := &NodeTraversal{}
	InitNodeTraversal(t, builder)
	return t
}

func NewNodeTraversal_3(compiler AbstractCompiler, cb NodeTraversal_Callback, sc ScopeCreator) *NodeTraversal {
	return NewNodeTraversal(&NodeTraversal_Builder{compiler: compiler, callback: cb, scopeCreator: sc})
}

func (this_ *NodeTraversal) GetScopeRoot() pkg_rhino.Node {
	if this_.scopeDepth <= 0 {
		return 0
	}
	return this_.scopeNodes[this_.scopeDepth-1]
}

func (this_ *NodeTraversal) priv_NodeTraversal_ensureScopeCapacity_4000() {
	if this_.scopeDepth == int32(len(this_.scopeNodes)) {
		newCap := int32(len(this_.scopeNodes) * 2)
		if newCap < 12 {
			newCap = 12
		}
		ns := make([]pkg_rhino.Node, newCap)
		copy(ns, this_.scopeNodes)
		this_.scopeNodes = ns
		os := make([]AbstractScope, newCap)
		copy(os, this_.scopeObjs)
		this_.scopeObjs = os
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_pushScope_4001(n pkg_rhino.Node) {
	this_.priv_NodeTraversal_ensureScopeCapacity_4000()
	this_.scopeNodes[this_.scopeDepth] = n
	this_.scopeObjs[this_.scopeDepth] = nil
	this_.scopeDepth++
	if this_.scopeCallback != nil {
		this_.scopeCallback.EnterScope(this_)
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_popScope_4004() {
	if this_.scopeCallback != nil {
		this_.scopeCallback.ExitScope(this_)
	}
	this_.scopeDepth--
	this_.scopeNodes[this_.scopeDepth] = 0
	this_.scopeObjs[this_.scopeDepth] = nil
}

func (this_ *NodeTraversal) priv_NodeTraversal_handleScript_3988(n pkg_rhino.Node, parent pkg_rhino.Node) {
	this_.currentNode = n
	if this_.callback.ShouldTraverse_3(this_, n, parent) {
		this_.priv_NodeTraversal_traverseChildren_3998(n)
		this_.currentNode = n
		this_.callback.Visit_3(this_, n, parent)
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_handleFunction_3989(n pkg_rhino.Node, parent pkg_rhino.Node) {
	this_.currentNode = n
	if this_.callback.ShouldTraverse_3(this_, n, parent) {
		this_.priv_NodeTraversal_traverseChildren_3998(n)
		this_.currentNode = n
		this_.callback.Visit_3(this_, n, parent)
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_handleModule_3990(n pkg_rhino.Node, parent pkg_rhino.Node) {
	this_.currentNode = n
	if this_.callback.ShouldTraverse_3(this_, n, parent) {
		this_.priv_NodeTraversal_traverseChildren_3998(n)
		this_.currentNode = n
		this_.callback.Visit_3(this_, n, parent)
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_handleClass_3995(n pkg_rhino.Node, parent pkg_rhino.Node) {
	this_.currentNode = n
	if this_.callback.ShouldTraverse_3(this_, n, parent) {
		this_.priv_NodeTraversal_traverseChildren_3998(n)
		this_.currentNode = n
		this_.callback.Visit_3(this_, n, parent)
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_handleClassMembers_3996(n pkg_rhino.Node, parent pkg_rhino.Node) {
	this_.currentNode = n
	if this_.callback.ShouldTraverse_3(this_, n, parent) {
		this_.priv_NodeTraversal_traverseChildren_3998(n)
		this_.currentNode = n
		this_.callback.Visit_3(this_, n, parent)
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_handleDestructuringOrDefaultValue_3991(n pkg_rhino.Node, parent pkg_rhino.Node) {
	this_.currentNode = n
	if this_.callback.ShouldTraverse_3(this_, n, parent) {
		var first pkg_rhino.Node = n.GetFirstChild()
		var second pkg_rhino.Node = first.GetNext()
		if second != 0 {
			if !(second.GetNext() == 0) {
				j2gort.CheckStateMsg(false, j2gort.BoxNode(second))
			}
			this_.priv_NodeTraversal_traverseBranch_3993(second, n)
		}
		this_.priv_NodeTraversal_traverseBranch_3993(first, n)
		this_.currentNode = n
		this_.callback.Visit_3(this_, n, parent)
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_createsBlockScope_3992(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(90):
		if this_.mayContainSyntheticBlocks && n.IsSyntheticBlock() {
			return false
		}
		var parent pkg_rhino.Node = n.GetParent()
		_ = parent
		return (((parent != 0) && !(NodeUtil_IsSwitchCase(parent))) && !(parent.IsCatch()))
	case int32(80), int32(81), int32(112), int32(113), int32(91), int32(105):
		return true
	default:
		return false
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseBranch_3993(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	switch n.TokenSwitchOrd_() {
	case int32(94):
		this_.priv_NodeTraversal_handleScript_3988(n, parent)
		return
	case int32(73):
		this_.priv_NodeTraversal_handleFunction_3989(n, parent)
		return
	case int32(123):
		this_.priv_NodeTraversal_handleModule_3990(n, parent)
		return
	case int32(105):
		this_.priv_NodeTraversal_handleClass_3995(n, parent)
		return
	case int32(106):
		this_.priv_NodeTraversal_handleClassMembers_3996(n, parent)
		return
	case int32(134), int32(104):
		if this_.obeyDestructuringAndDefaultValueExecutionOrder {
			this_.priv_NodeTraversal_handleDestructuringOrDefaultValue_3991(n, parent)
			return
		}
	default:
	}
	this_.currentNode = n
	if !(this_.callback.ShouldTraverse_3(this_, n, parent)) {
		return
	}
	var createsBlockScope bool = this_.priv_NodeTraversal_createsBlockScope_3992(n)
	_ = createsBlockScope
	var previousHoistScopeRoot pkg_rhino.Node = this_.currentHoistScopeRoot
	_ = previousHoistScopeRoot
	if createsBlockScope {
		this_.priv_NodeTraversal_pushScope_4001(n)
		if NodeUtil_IsClassStaticBlock(n) {
			this_.currentHoistScopeRoot = n
		}
	}
	for child := n.GetFirstChild(); child != 0; {
		next := child.GetNext()
		this_.priv_NodeTraversal_traverseBranch_3993(child, n)
		child = next
	}
	if createsBlockScope {
		this_.priv_NodeTraversal_popScope_4004()
		this_.currentHoistScopeRoot = previousHoistScopeRoot
	}
	this_.currentNode = n
	this_.callback.Visit_3(this_, n, parent)
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseChildren_3998(n pkg_rhino.Node) {
	_ = n
	for child := n.GetFirstChild(); child != 0; {
		next := child.GetNext()
		this_.priv_NodeTraversal_traverseBranch_3993(child, n)
		child = next
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverse_3980(root pkg_rhino.Node) {
	this_.priv_NodeTraversal_traverseBranch_3993(root, 0)
}

func (this_ *NodeTraversal) Traverse(root pkg_rhino.Node) {
	this_.priv_NodeTraversal_traverse_3980(root)
}

func (this_ *NodeTraversal) TraverseChildrenExported(root pkg_rhino.Node) {
	this_.priv_NodeTraversal_traverseChildren_3998(root)
}

// Reference unlowered traversal for differential parity verification.
func (this_ *NodeTraversal) ReferenceTraverseBranch(n, parent pkg_rhino.Node) {
	switch n.TokenSwitchOrd_() {
	case int32(94):
		this_.currentNode = n
		if this_.callback.ShouldTraverse_3(this_, n, parent) {
			this_.ReferenceTraverseChildren(n)
			this_.currentNode = n
			this_.callback.Visit_3(this_, n, parent)
		}
		return
	case int32(73):
		this_.currentNode = n
		if this_.callback.ShouldTraverse_3(this_, n, parent) {
			this_.ReferenceTraverseChildren(n)
			this_.currentNode = n
			this_.callback.Visit_3(this_, n, parent)
		}
		return
	case int32(123), int32(105), int32(106):
		this_.currentNode = n
		if this_.callback.ShouldTraverse_3(this_, n, parent) {
			this_.ReferenceTraverseChildren(n)
			this_.currentNode = n
			this_.callback.Visit_3(this_, n, parent)
		}
		return
	}
	this_.currentNode = n
	if !this_.callback.ShouldTraverse_3(this_, n, parent) {
		return
	}
	createsBlockScope := this_.priv_NodeTraversal_createsBlockScope_3992(n)
	prevHoist := this_.currentHoistScopeRoot
	if createsBlockScope {
		this_.priv_NodeTraversal_pushScope_4001(n)
		if NodeUtil_IsClassStaticBlock(n) {
			this_.currentHoistScopeRoot = n
		}
	}
	for child := n.GetFirstChild(); child != 0; {
		next := child.GetNext()
		this_.ReferenceTraverseBranch(child, n)
		child = next
	}
	if createsBlockScope {
		this_.priv_NodeTraversal_popScope_4004()
		this_.currentHoistScopeRoot = prevHoist
	}
	this_.currentNode = n
	this_.callback.Visit_3(this_, n, parent)
}

func (this_ *NodeTraversal) ReferenceTraverseChildren(n pkg_rhino.Node) {
	for child := n.GetFirstChild(); child != 0; {
		next := child.GetNext()
		this_.ReferenceTraverseBranch(child, n)
		child = next
	}
}
`

const wave47bNodeTraversalWave22CPreLowered = `package jscomp

import (
	"example.com/w47btest/gen/j2gort"
	pkg_rhino "example.com/w47btest/gen/rhino"
	"example.com/w47btest/shims/nodestore"
)

func clinit_NodeTraversal() {}

var nodeTraversalBlockScopeTok_ = [256]bool{80: true, 81: true, 90: true, 91: true, 105: true, 112: true, 113: true}
var nodeTraversalSpecialOrBlockTok_ = [256]bool{73: true, 80: true, 81: true, 90: true, 91: true, 94: true, 104: true, 105: true, 106: true, 112: true, 113: true, 123: true, 134: true, 255: true}

type NodeTraversal struct {
	compiler                                        AbstractCompiler
	callback                                        NodeTraversal_Callback
	scopeCallback                                   NodeTraversal_ScopedCallback
	scopeCreator                                    ScopeCreator
	currentNode                                     pkg_rhino.Node
	currentHoistScopeRoot                           pkg_rhino.Node
	currentFunction                                 pkg_rhino.Node
	scopeNodes                                      []pkg_rhino.Node
	scopeObjs                                       []AbstractScope
	scopeDepth                                      int32
	mayContainSyntheticBlocks                       bool
	obeyDestructuringAndDefaultValueExecutionOrder  bool
	scopeNodesBuf_                                  [12]pkg_rhino.Node
	scopeObjsBuf_                                   [12]AbstractScope
	defaultScopeCreator_                            SyntacticScopeCreator
	alwaysTraverse_                                 bool
	noOpVisit_                                      bool
}

func InitNodeTraversal(this_ *NodeTraversal, builder *NodeTraversal_Builder) {
	this_.compiler = builder.compiler
	this_.callback = builder.callback
	if scb, ok := builder.callback.(NodeTraversal_ScopedCallback); ok {
		this_.scopeCallback = scb
	}
	if builder.scopeCreator == nil || builder.scopeCreator.IsJavaNil_() {
		this_.scopeCreator = NewSyntacticScopeCreator_1(this_.compiler)
	} else {
		this_.scopeCreator = builder.scopeCreator
	}
	_, this_.alwaysTraverse_ = any(this_.callback).(interface{ AlwaysShouldTraverse_() })
	_, this_.noOpVisit_ = any(this_.callback).(interface{ NoOpVisit_() })
}

func NewNodeTraversal(builder *NodeTraversal_Builder) *NodeTraversal {
	t := &NodeTraversal{}
	InitNodeTraversal(t, builder)
	return t
}

func NewNodeTraversal_3(compiler AbstractCompiler, cb NodeTraversal_Callback, sc ScopeCreator) *NodeTraversal {
	return NewNodeTraversal(&NodeTraversal_Builder{compiler: compiler, callback: cb, scopeCreator: sc})
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverse_3980(root pkg_rhino.Node) {
	this_.priv_NodeTraversal_traverseBranch_3993(root, 0)
}

func (this_ *NodeTraversal) GetScopeRoot() pkg_rhino.Node {
	if this_.scopeDepth <= 0 {
		return 0
	}
	return this_.scopeNodes[this_.scopeDepth-1]
}

func (this_ *NodeTraversal) priv_NodeTraversal_pushScope_4001(n pkg_rhino.Node) {
	this_.scopeNodes = append(this_.scopeNodes, n)
	this_.scopeObjs = append(this_.scopeObjs, nil)
	this_.scopeDepth++
	if this_.scopeCallback != nil {
		this_.scopeCallback.EnterScope(this_)
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_popScope_4004() {
	if this_.scopeCallback != nil {
		this_.scopeCallback.ExitScope(this_)
	}
	this_.scopeDepth--
	this_.scopeNodes = this_.scopeNodes[:this_.scopeDepth]
	this_.scopeObjs = this_.scopeObjs[:this_.scopeDepth]
}

func (this_ *NodeTraversal) priv_NodeTraversal_handleScript_3988(n, parent pkg_rhino.Node)              { this_.priv_NodeTraversal_traverseChildren_3998(n) }
func (this_ *NodeTraversal) priv_NodeTraversal_handleFunction_3989(n, parent pkg_rhino.Node)            { this_.priv_NodeTraversal_traverseChildren_3998(n) }
func (this_ *NodeTraversal) priv_NodeTraversal_handleModule_3990(n, parent pkg_rhino.Node)              { this_.priv_NodeTraversal_traverseChildren_3998(n) }
func (this_ *NodeTraversal) priv_NodeTraversal_handleClass_3995(n, parent pkg_rhino.Node)               { this_.priv_NodeTraversal_traverseChildren_3998(n) }
func (this_ *NodeTraversal) priv_NodeTraversal_handleClassMembers_3996(n, parent pkg_rhino.Node)        { this_.priv_NodeTraversal_traverseChildren_3998(n) }
func (this_ *NodeTraversal) priv_NodeTraversal_handleDestructuringOrDefaultValue_3991(n, p pkg_rhino.Node) {
	_ = j2gort.BoxNode(n)
	this_.priv_NodeTraversal_traverseChildren_3998(n)
}

func (this_ *NodeTraversal) priv_NodeTraversal_createsBlockScope_3992(n pkg_rhino.Node) bool {
	_ = n
	switch n.TokenSwitchOrd_() {
	case int32(90):
		var parent pkg_rhino.Node = n.GetParent()
		return parent != 0
	case int32(80), int32(81), int32(112), int32(113), int32(91), int32(105):
		return true
	default:
		return false
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseBranch_3993(n pkg_rhino.Node, parent pkg_rhino.Node) {
	_ = n
	_ = parent
	_, _, tok_ := n.FirstNextTok_()
	if tok_ == 0xFF {
		_ = n.TokenSwitchOrd_()
	}
	switch int32(tok_) {
	case int32(94):
		this_.priv_NodeTraversal_handleScript_3988(n, parent)
		return
	case int32(73):
		this_.priv_NodeTraversal_handleFunction_3989(n, parent)
		return
	case int32(123):
		this_.priv_NodeTraversal_handleModule_3990(n, parent)
		return
	case int32(105):
		this_.priv_NodeTraversal_handleClass_3995(n, parent)
		return
	case int32(106):
		this_.priv_NodeTraversal_handleClassMembers_3996(n, parent)
		return
	case int32(134), int32(104):
		if this_.obeyDestructuringAndDefaultValueExecutionOrder {
			this_.priv_NodeTraversal_handleDestructuringOrDefaultValue_3991(n, parent)
			return
		}
	default:
	}
	this_.currentNode = n
	if !this_.alwaysTraverse_ {
		if !(this_.callback.ShouldTraverse_3(this_, n, parent)) {
			return
		}
	}
	var createsBlockScope bool
	if nodeTraversalBlockScopeTok_[tok_] {
		createsBlockScope = this_.priv_NodeTraversal_createsBlockScope_3992(n)
	}
	_ = createsBlockScope
	var previousHoistScopeRoot pkg_rhino.Node = this_.currentHoistScopeRoot
	_ = previousHoistScopeRoot
	if createsBlockScope {
		this_.priv_NodeTraversal_pushScope_4001(n)
		if NodeUtil_IsClassStaticBlock(n) {
			this_.currentHoistScopeRoot = n
		}
	}
	var curBlkIdx_ uint32 = 0xFFFFFFFF
	var curBlk_ *nodestore.NodeBlock
	for child_ := n.GetFirstChild(); child_ != 0; {
		var cFirst_, next_ pkg_rhino.Node
		var cTok_ uint8
		if idx_ := uint32(child_) >> nodestore.BlockShift; idx_ == curBlkIdx_ {
			f_, nx_, t_ := curBlk_.FirstNextTok(nodestore.Node(child_))
			cFirst_, next_, cTok_ = pkg_rhino.Node(f_), pkg_rhino.Node(nx_), t_
		} else {
			curBlk_ = nodestore.BlockOf(nodestore.Node(child_))
			curBlkIdx_ = idx_
			f_, nx_, t_ := curBlk_.FirstNextTok(nodestore.Node(child_))
			cFirst_, next_, cTok_ = pkg_rhino.Node(f_), pkg_rhino.Node(nx_), t_
		}
		if cFirst_ == 0 && cTok_ != 0xFF && !nodeTraversalSpecialOrBlockTok_[cTok_] {
			if this_.alwaysTraverse_ {
				this_.currentNode = child_
				if !this_.noOpVisit_ {
					this_.callback.Visit_3(this_, child_, n)
				}
				child_ = next_
				continue
			}
			this_.currentNode = child_
			if !this_.callback.ShouldTraverse_3(this_, child_, n) {
				child_ = next_
				continue
			}
			if child_.GetFirstChild() == 0 {
				this_.currentNode = child_
				if !this_.noOpVisit_ {
					this_.callback.Visit_3(this_, child_, n)
				}
				child_ = next_
				continue
			}
			for gc_ := child_.GetFirstChild(); gc_ != 0; {
				var gNext_ pkg_rhino.Node
				if gIdx_ := uint32(gc_) >> nodestore.BlockShift; gIdx_ == curBlkIdx_ {
					_, gnx_, _ := curBlk_.FirstNextTok(nodestore.Node(gc_))
					gNext_ = pkg_rhino.Node(gnx_)
				} else {
					curBlk_ = nodestore.BlockOf(nodestore.Node(gc_))
					curBlkIdx_ = gIdx_
					_, gnx_, _ := curBlk_.FirstNextTok(nodestore.Node(gc_))
					gNext_ = pkg_rhino.Node(gnx_)
				}
				this_.priv_NodeTraversal_traverseBranch_3993(gc_, child_)
				gc_ = gNext_
			}
			this_.currentNode = child_
			if !this_.noOpVisit_ {
				this_.callback.Visit_3(this_, child_, n)
			}
			child_ = next_
			continue
		}
		this_.priv_NodeTraversal_traverseBranch_3993(child_, n)
		child_ = next_
	}
	if createsBlockScope {
		this_.priv_NodeTraversal_popScope_4004()
		this_.currentHoistScopeRoot = previousHoistScopeRoot
	}
	this_.currentNode = n
	if !this_.noOpVisit_ {
		this_.callback.Visit_3(this_, n, parent)
	}
}

func (this_ *NodeTraversal) priv_NodeTraversal_traverseChildren_3998(n pkg_rhino.Node) {
	_ = n
	var curBlkIdx_ uint32 = 0xFFFFFFFF
	var curBlk_ *nodestore.NodeBlock
	for child_ := n.GetFirstChild(); child_ != 0; {
		var cFirst_, next_ pkg_rhino.Node
		var cTok_ uint8
		if idx_ := uint32(child_) >> nodestore.BlockShift; idx_ == curBlkIdx_ {
			f_, nx_, t_ := curBlk_.FirstNextTok(nodestore.Node(child_))
			cFirst_, next_, cTok_ = pkg_rhino.Node(f_), pkg_rhino.Node(nx_), t_
		} else {
			curBlk_ = nodestore.BlockOf(nodestore.Node(child_))
			curBlkIdx_ = idx_
			f_, nx_, t_ := curBlk_.FirstNextTok(nodestore.Node(child_))
			cFirst_, next_, cTok_ = pkg_rhino.Node(f_), pkg_rhino.Node(nx_), t_
		}
		if cFirst_ == 0 && cTok_ != 0xFF && !nodeTraversalSpecialOrBlockTok_[cTok_] {
			if this_.alwaysTraverse_ {
				this_.currentNode = child_
				if !this_.noOpVisit_ {
					this_.callback.Visit_3(this_, child_, n)
				}
				child_ = next_
				continue
			}
			this_.currentNode = child_
			if !this_.callback.ShouldTraverse_3(this_, child_, n) {
				child_ = next_
				continue
			}
			if child_.GetFirstChild() == 0 {
				this_.currentNode = child_
				if !this_.noOpVisit_ {
					this_.callback.Visit_3(this_, child_, n)
				}
				child_ = next_
				continue
			}
			for gc_ := child_.GetFirstChild(); gc_ != 0; {
				var gNext_ pkg_rhino.Node
				if gIdx_ := uint32(gc_) >> nodestore.BlockShift; gIdx_ == curBlkIdx_ {
					_, gnx_, _ := curBlk_.FirstNextTok(nodestore.Node(gc_))
					gNext_ = pkg_rhino.Node(gnx_)
				} else {
					curBlk_ = nodestore.BlockOf(nodestore.Node(gc_))
					curBlkIdx_ = gIdx_
					_, gnx_, _ := curBlk_.FirstNextTok(nodestore.Node(gc_))
					gNext_ = pkg_rhino.Node(gnx_)
				}
				this_.priv_NodeTraversal_traverseBranch_3993(gc_, child_)
				gc_ = gNext_
			}
			this_.currentNode = child_
			if !this_.noOpVisit_ {
				this_.callback.Visit_3(this_, child_, n)
			}
			child_ = next_
			continue
		}
		this_.priv_NodeTraversal_traverseBranch_3993(child_, n)
		child_ = next_
	}
}
`

const wave47bRuntimeTestSrc = `package jscomp_test

import (
	"reflect"
	"runtime"
	"testing"

	"example.com/w47btest/gen/jscomp"
	pkg_rhino "example.com/w47btest/gen/rhino"
	"example.com/w47btest/shims/nodestore"
)

func mk(blk, slot uint32, tok int32) pkg_rhino.Node {
	return pkg_rhino.Node(nodestore.MakeNode(blk, slot, uint8(tok)))
}

func addBack(p, c pkg_rhino.Node) {
	nodestore.AddChildToBack(nodestore.Node(p), nodestore.Node(c))
}

func addFront(p, c pkg_rhino.Node) {
	nodestore.AddChildToFront(nodestore.Node(p), nodestore.Node(c))
}

func TestWave47BRuntime_AlwaysTraverseFastPathAndZeroRedundantBlockOf(t *testing.T) {
	// Build a multi-level tree all in block 1 (zero special/block tokens inside):
	// EXPR_RESULT -> ADD -> (ADD -> (NAME, NUMBER), NAME, NUMBER)
	root := mk(1, 1, pkg_rhino.TokEXPR_RESULT)
	add1 := mk(1, 2, pkg_rhino.TokADD)
	add2 := mk(1, 3, pkg_rhino.TokADD)
	leaf1 := mk(1, 4, pkg_rhino.TokNAME)
	leaf2 := mk(1, 5, pkg_rhino.TokNUMBER)
	leaf3 := mk(1, 6, pkg_rhino.TokNAME)
	leaf4 := mk(1, 7, pkg_rhino.TokNUMBER)
	addBack(root, add1)
	addBack(add1, add2)
	addBack(add2, leaf1)
	addBack(add2, leaf2)
	addBack(add1, leaf3)
	addBack(add1, leaf4)

	refRec := &jscomp.PostOrderRecorder{}
	jscomp.NewNodeTraversal_3(nil, refRec, nil).ReferenceTraverseBranch(root, 0)

	nodestore.ResetCounters()
	optRec := &jscomp.PostOrderRecorder{}
	jscomp.NewNodeTraversal_3(nil, optRec, nil).Traverse(root)

	if !reflect.DeepEqual(optRec.Events, refRec.Events) {
		t.Fatalf("event mismatch:\n got: %v\nwant: %v", optRec.Events, refRec.Events)
	}
	if got := nodestore.BlockOfCalls.Load(); got != 1 {
		t.Fatalf("expected exactly 1 BlockOf call for entire same-block tree (trampoline entry only), got %d", got)
	}
	if got := nodestore.FirstNextTokNodeCalls.Load(); got != 0 {
		t.Fatalf("expected 0 Node.FirstNextTok_() calls, got %d", got)
	}
}

func TestWave47BRuntime_ShouldTraverseMutatesChildrenOfNAndLeafChild(t *testing.T) {
	buildTree := func(blk uint32) (root, mid, leafA, leafB, extra1, extra2, extraFn, fnLeaf pkg_rhino.Node) {
		root = mk(blk, 10, pkg_rhino.TokEXPR_RESULT)
		mid = mk(blk, 11, pkg_rhino.TokADD)
		leafA = mk(blk, 12, pkg_rhino.TokNAME)
		leafB = mk(blk, 13, pkg_rhino.TokNUMBER)
		extra1 = mk(blk, 14, pkg_rhino.TokNAME)
		extra2 = mk(blk+1, 15, pkg_rhino.TokNUMBER) // cross-block grandchild!
		extraFn = mk(blk+2, 16, pkg_rhino.TokFUNCTION)
		fnLeaf = mk(blk+2, 17, pkg_rhino.TokNAME)
		addBack(extraFn, fnLeaf)
		addBack(root, mid)
		addBack(mid, leafA)
		addBack(mid, leafB)
		return
	}

	// Reference run
	rRoot, rMid, rLeafA, rLeafB, rExtra1, rExtra2, rExtraFn, _ := buildTree(2)
	refCb := &jscomp.GeneralMutatingCallback{
		OnShould: func(tr *jscomp.NodeTraversal, n, parent pkg_rhino.Node) bool {
			if n == rMid {
				// Prepend extra1 to mid's children during ShouldTraverse_3!
				addFront(rMid, rExtra1)
			}
			if n == rLeafA {
				// Attach cross-block grandchildren (including a FUNCTION) to leafA during ShouldTraverse_3!
				addBack(rLeafA, rExtra2)
				addBack(rLeafA, rExtraFn)
			}
			if n == rLeafB {
				// Skip leafB
				return false
			}
			return true
		},
	}
	jscomp.NewNodeTraversal_3(nil, refCb, nil).ReferenceTraverseBranch(rRoot, 0)

	// Lowered run (using same node handles in fresh blocks 2..4 after rebuilding identically)
	oRoot, oMid, oLeafA, oLeafB, oExtra1, oExtra2, oExtraFn, _ := buildTree(2)
	optCb := &jscomp.GeneralMutatingCallback{
		OnShould: func(tr *jscomp.NodeTraversal, n, parent pkg_rhino.Node) bool {
			if n == oMid {
				addFront(oMid, oExtra1)
			}
			if n == oLeafA {
				addBack(oLeafA, oExtra2)
				addBack(oLeafA, oExtraFn)
			}
			if n == oLeafB {
				return false
			}
			return true
		},
	}
	nodestore.ResetCounters()
	jscomp.NewNodeTraversal_3(nil, optCb, nil).Traverse(oRoot)

	if !reflect.DeepEqual(optCb.Events, refCb.Events) {
		t.Fatalf("mutating ShouldTraverse_3 event mismatch:\n got: %v\nwant: %v", optCb.Events, refCb.Events)
	}
	if got := nodestore.FirstNextTokNodeCalls.Load(); got != 0 {
		t.Fatalf("expected 0 Node.FirstNextTok_() calls, got %d", got)
	}
}

func TestWave47BRuntime_CreatesBlockScopeEnterScopeMutatesChildren(t *testing.T) {
	buildTree := func() (root, blk, origChild, replacedChild pkg_rhino.Node) {
		root = mk(5, 1, pkg_rhino.TokEXPR_RESULT)
		blk = mk(5, 2, pkg_rhino.TokBLOCK) // parent != 0 -> createsBlockScope == true!
		origChild = mk(5, 3, pkg_rhino.TokNAME)
		replacedChild = mk(5, 4, pkg_rhino.TokNUMBER)
		addBack(root, blk)
		addBack(blk, origChild)
		return
	}

	rRoot, rBlk, _, rRepl := buildTree()
	refCb := &jscomp.ScopedPostOrderMutator{
		OnEnterScope: func(tr *jscomp.NodeTraversal, scopeRoot pkg_rhino.Node) {
			if scopeRoot == rBlk {
				nodestore.RemoveChildren(nodestore.Node(rBlk))
				addBack(rBlk, rRepl)
			}
		},
	}
	jscomp.NewNodeTraversal_3(nil, refCb, nil).ReferenceTraverseBranch(rRoot, 0)

	oRoot, oBlk, _, oRepl := buildTree()
	optCb := &jscomp.ScopedPostOrderMutator{
		OnEnterScope: func(tr *jscomp.NodeTraversal, scopeRoot pkg_rhino.Node) {
			if scopeRoot == oBlk {
				nodestore.RemoveChildren(nodestore.Node(oBlk))
				addBack(oBlk, oRepl)
			}
		},
	}
	nodestore.ResetCounters()
	jscomp.NewNodeTraversal_3(nil, optCb, nil).Traverse(oRoot)

	if !reflect.DeepEqual(optCb.Events, refCb.Events) {
		t.Fatalf("EnterScope-mutated block scope event mismatch:\n got: %v\nwant: %v", optCb.Events, refCb.Events)
	}
	// All nodes are in block 5 -> BlockOf must be called ONLY ONCE (on trampoline entry)!
	if got := nodestore.BlockOfCalls.Load(); got != 1 {
		t.Fatalf("expected exactly 1 BlockOf call for same-block EnterScope mutation, got %d", got)
	}
}

func TestWave47BRuntime_CrossBlockTreeAndTraverseChildren(t *testing.T) {
	// Nodes spread across blocks 6, 7, 8, 9
	root := mk(6, 1, pkg_rhino.TokSCRIPT)
	c1 := mk(7, 2, pkg_rhino.TokADD)
	gc1 := mk(8, 3, pkg_rhino.TokNAME)
	gc2 := mk(6, 4, pkg_rhino.TokNUMBER)
	c2 := mk(9, 5, pkg_rhino.TokNAME)
	c3 := mk(7, 6, pkg_rhino.TokEXPR_RESULT)
	gc3 := mk(9, 7, pkg_rhino.TokNUMBER)
	addBack(root, c1)
	addBack(c1, gc1)
	addBack(c1, gc2)
	addBack(root, c2)
	addBack(root, c3)
	addBack(c3, gc3)

	// 1. Traverse from SCRIPT (goes through handleScript -> traverseChildren -> traverseBranchBlk_)
	refRec := &jscomp.PostOrderRecorder{}
	jscomp.NewNodeTraversal_3(nil, refRec, nil).ReferenceTraverseBranch(root, 0)

	optRec := &jscomp.PostOrderRecorder{}
	jscomp.NewNodeTraversal_3(nil, optRec, nil).Traverse(root)
	if !reflect.DeepEqual(optRec.Events, refRec.Events) {
		t.Fatalf("cross-block SCRIPT mismatch:\n got: %v\nwant: %v", optRec.Events, refRec.Events)
	}

	// 2. Direct TraverseChildrenExported
	refRec2 := &jscomp.PostOrderRecorder{}
	jscomp.NewNodeTraversal_3(nil, refRec2, nil).ReferenceTraverseChildren(root)

	optRec2 := &jscomp.PostOrderRecorder{}
	jscomp.NewNodeTraversal_3(nil, optRec2, nil).TraverseChildrenExported(root)
	if !reflect.DeepEqual(optRec2.Events, refRec2.Events) {
		t.Fatalf("cross-block traverseChildren mismatch:\n got: %v\nwant: %v", optRec2.Events, refRec2.Events)
	}
}

func TestWave47BRuntime_ZeroNodeNPEAndNullTokenPanic(t *testing.T) {
	assertRuntimeErr := func(name string, fn func()) {
		t.Helper()
		defer func() {
			r := recover()
			if r == nil {
				t.Fatalf("%s: expected panic on nil/0 node, got nil", name)
			}
			if _, ok := r.(runtime.Error); !ok {
				t.Fatalf("%s: expected runtime.Error (NPE), got %T (%v)", name, r, r)
			}
		}()
		fn()
	}

	nt := jscomp.NewNodeTraversal_3(nil, &jscomp.PostOrderRecorder{}, nil)
	assertRuntimeErr("Traverse(0)", func() { nt.Traverse(0) })
	assertRuntimeErr("TraverseChildren(0)", func() { nt.TraverseChildrenExported(0) })

	// 0xFF token at root and as child -> must panic with NullTokenPanic (via TokenSwitchOrd_())
	assertNullTokPanic := func(name string, fn func()) {
		t.Helper()
		defer func() {
			r := recover()
			if _, ok := r.(pkg_rhino.NullTokenPanic); !ok {
				t.Fatalf("%s: expected NullTokenPanic, got %T (%v)", name, r, r)
			}
		}()
		fn()
	}

	badRoot := pkg_rhino.Node(nodestore.MakeNode(10, 1, 0xFF))
	assertNullTokPanic("Traverse(0xFF root)", func() { nt.Traverse(badRoot) })

	goodParent := mk(10, 2, pkg_rhino.TokEXPR_RESULT)
	badChild := pkg_rhino.Node(nodestore.MakeNode(10, 3, 0xFF))
	addBack(goodParent, badChild)
	assertNullTokPanic("Traverse(0xFF child)", func() { nt.Traverse(goodParent) })
}
`

func writeWave47BFiles(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		full := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func extractFuncBodySnippet(t *testing.T, src, funcSigPrefix string) string {
	t.Helper()
	idx := strings.Index(src, funcSigPrefix)
	if idx < 0 {
		t.Fatalf("function %q not found in:\n%s", funcSigPrefix, src)
	}
	rest := src[idx:]
	// Find matching closing brace at column 0 ("\n}")
	end := strings.Index(rest, "\n}\n")
	if end < 0 {
		if strings.HasSuffix(rest, "\n}") {
			return rest
		}
		t.Fatalf("could not locate closing brace for %q", funcSigPrefix)
	}
	return rest[:end+2]
}

func assertWave47BLoweredNodeTraversalShape(t *testing.T, ntSrc string) {
	t.Helper()
	for _, want := range []string{
		"func (this_ *NodeTraversal) traverseBranchBlk_(n, parent, nFirst pkg_rhino.Node, nTok uint8, curBlkIdx_ uint32, curBlk_ *nodestore.NodeBlock)",
		"if this_.alwaysTraverse_ && !createsBlockScope {",
		"child_ = nFirst",
		"this_.traverseBranchBlk_(child_, n, cFirst_, cTok_, curBlkIdx_, curBlk_)",
		"this_.traverseBranchBlk_(gc_, child_, gFirst_, gTok_, curBlkIdx_, curBlk_)",
	} {
		if !strings.Contains(ntSrc, want) {
			t.Fatalf("lowered node_traversal.go missing %q\n--- got ---\n%s", want, ntSrc)
		}
	}

	// Verify priv_NodeTraversal_traverseBranch_3993 is a lean 4-line trampoline!
	tbFunc := extractFuncBodySnippet(t, ntSrc, "func (this_ *NodeTraversal) priv_NodeTraversal_traverseBranch_3993(")
	for _, want := range []string{
		"curBlkIdx_ := uint32(n) >> nodestore.BlockShift",
		"curBlk_ := nodestore.BlockOf(nodestore.Node(n))",
		"f_, _, tok_ := curBlk_.FirstNextTok(nodestore.Node(n))",
		"this_.traverseBranchBlk_(n, parent, pkg_rhino.Node(f_), tok_, curBlkIdx_, curBlk_)",
	} {
		if !strings.Contains(tbFunc, want) {
			t.Fatalf("traverseBranch trampoline missing %q:\n%s", want, tbFunc)
		}
	}
	if strings.Contains(tbFunc, "FirstNextTok_()") || strings.Contains(tbFunc, "GetFirstChild()") {
		t.Fatalf("traverseBranch trampoline must NOT call FirstNextTok_() or GetFirstChild():\n%s", tbFunc)
	}
	if lines := strings.Count(tbFunc, "\n"); lines > 8 {
		t.Fatalf("traverseBranch trampoline should be lean (<=8 lines), got %d:\n%s", lines, tbFunc)
	}

	// Verify traverseBranchBlk_ has ZERO FirstNextTok_() and ZERO GetFirstChild()!
	blkFunc := extractFuncBodySnippet(t, ntSrc, "func (this_ *NodeTraversal) traverseBranchBlk_(")
	if strings.Contains(blkFunc, "FirstNextTok_()") {
		t.Fatalf("traverseBranchBlk_ must NOT call FirstNextTok_():\n%s", blkFunc)
	}
	if strings.Contains(blkFunc, "GetFirstChild()") {
		t.Fatalf("traverseBranchBlk_ must NOT call GetFirstChild():\n%s", blkFunc)
	}
	if !strings.Contains(blkFunc, "if nTok == 0xFF") || !strings.Contains(blkFunc, "switch int32(nTok)") || !strings.Contains(blkFunc, "nodeTraversalBlockScopeTok_[nTok]") {
		t.Fatalf("traverseBranchBlk_ must use nTok directly:\n%s", blkFunc)
	}

	// Verify priv_NodeTraversal_traverseChildren_3998 initializes curBlk_ from n and calls traverseBranchBlk_!
	tcFunc := extractFuncBodySnippet(t, ntSrc, "func (this_ *NodeTraversal) priv_NodeTraversal_traverseChildren_3998(")
	for _, want := range []string{
		"curBlkIdx_ := uint32(n) >> nodestore.BlockShift",
		"curBlk_ := nodestore.BlockOf(nodestore.Node(n))",
		"nFirst_, _, _ := curBlk_.FirstNextTok(nodestore.Node(n))",
		"this_.traverseBranchBlk_(child_, n, cFirst_, cTok_, curBlkIdx_, curBlk_)",
		"this_.traverseBranchBlk_(gc_, child_, gFirst_, gTok_, curBlkIdx_, curBlk_)",
	} {
		if !strings.Contains(tcFunc, want) {
			t.Fatalf("traverseChildren missing %q:\n%s", want, tcFunc)
		}
	}
	if strings.Contains(tcFunc, "GetFirstChild()") || strings.Contains(tcFunc, "FirstNextTok_()") {
		t.Fatalf("traverseChildren must NOT call GetFirstChild() or FirstNextTok_():\n%s", tcFunc)
	}
}

func TestWave47B_NodeTraversalBlkLoweringAndRuntimeParity(t *testing.T) {
	// 1. Fresh unlowered NodeTraversal -> Run 1 lowering + Run 2 idempotency + compiled runtime parity!
	tmpFresh := t.TempDir()
	writeWave47BFiles(t, tmpFresh, map[string]string{
		"go.mod":                    "module example.com/w47btest\n\ngo 1.24\n",
		"shims/nodestore/store.go":  wave47bNodestoreStub,
		"gen/j2gort/j2gort.go":      wave47bJ2gortStub,
		"gen/rhino/node.go":         wave47bRhinoStub,
		"gen/jscomp/common.go":      wave47bJscompCommon,
		"gen/jscomp/node_traversal.go": wave47bNodeTraversalFresh,
		"gen/jscomp/runtime_test.go": wave47bRuntimeTestSrc,
	})

	jscompFreshDir := filepath.Join(tmpFresh, "gen", "jscomp")
	stats1, err := RunOnPackageDir(jscompFreshDir, []Pass{NewNodeTraversalLoweringPass()}, Options{})
	if err != nil {
		t.Fatalf("Run 1 on fresh jscomp failed: %v", err)
	}
	if stats1.TotalSites == 0 || stats1.FilesModified == 0 {
		t.Fatalf("expected non-zero Run 1 sites on fresh jscomp, got %+v", stats1)
	}

	ntFreshBytes, err := os.ReadFile(filepath.Join(jscompFreshDir, "node_traversal.go"))
	if err != nil {
		t.Fatal(err)
	}
	assertWave47BLoweredNodeTraversalShape(t, string(ntFreshBytes))

	// Run 2 idempotency on fresh lowered package!
	stats2, err := RunOnPackageDir(jscompFreshDir, []Pass{NewNodeTraversalLoweringPass()}, Options{})
	if err != nil {
		t.Fatalf("Run 2 on fresh jscomp failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT idempotent on fresh jscomp: %+v", stats2)
	}

	// Execute compiled runtime parity tests on the lowered package!
	cmd := exec.Command("go", "test", "-p", "1", "-parallel", "1", "-v", "./gen/jscomp")
	cmd.Dir = tmpFresh
	cmd.Env = append(os.Environ(), "GOFLAGS=-p=2")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("compiled runtime parity tests failed: %v\n%s", err, string(out))
	}

	// 2. Pre-lowered Wave 22C NodeTraversal -> Run 1 upgrade + Run 2 idempotency + compile check!
	tmpUpg := t.TempDir()
	writeWave47BFiles(t, tmpUpg, map[string]string{
		"go.mod":                       "module example.com/w47btest\n\ngo 1.24\n",
		"shims/nodestore/store.go":     wave47bNodestoreStub,
		"gen/j2gort/j2gort.go":         wave47bJ2gortStub,
		"gen/rhino/node.go":            wave47bRhinoStub,
		"gen/jscomp/common.go":         wave47bJscompCommon,
		"gen/jscomp/node_traversal.go": wave47bNodeTraversalWave22CPreLowered,
	})
	jscompUpgDir := filepath.Join(tmpUpg, "gen", "jscomp")
	upg1, err := RunOnPackageDir(jscompUpgDir, []Pass{NewNodeTraversalLoweringPass()}, Options{})
	if err != nil {
		t.Fatalf("Run 1 upgrade from Wave 22C failed: %v", err)
	}
	if upg1.TotalSites == 0 || upg1.FilesModified == 0 {
		t.Fatalf("expected non-zero upgrade sites from Wave 22C, got %+v", upg1)
	}
	ntUpgBytes, err := os.ReadFile(filepath.Join(jscompUpgDir, "node_traversal.go"))
	if err != nil {
		t.Fatal(err)
	}
	assertWave47BLoweredNodeTraversalShape(t, string(ntUpgBytes))

	upg2, err := RunOnPackageDir(jscompUpgDir, []Pass{NewNodeTraversalLoweringPass()}, Options{})
	if err != nil {
		t.Fatalf("Run 2 idempotency after Wave 22C upgrade failed: %v", err)
	}
	if upg2.TotalSites != 0 || upg2.FilesModified != 0 {
		t.Fatalf("Run 2 NOT idempotent after Wave 22C upgrade: %+v", upg2)
	}

	// Verify upgraded Wave 22C package compiles cleanly!
	buildCmd := exec.Command("go", "build", "./gen/jscomp")
	buildCmd.Dir = tmpUpg
	buildCmd.Env = append(os.Environ(), "GOFLAGS=-p=2")
	if bOut, bErr := buildCmd.CombinedOutput(); bErr != nil {
		t.Fatalf("go build on upgraded Wave 22C package failed: %v\n%s", bErr, string(bOut))
	}
}

func TestWave47B_RealNodeTraversalDryRunAndIdempotency(t *testing.T) {
	realJscompDir, err := filepath.Abs(filepath.Join("..", "..", "gen", "jscomp"))
	if err != nil {
		t.Fatal(err)
	}
	ntPath := filepath.Join(realJscompDir, "node_traversal.go")
	origBytes, err := os.ReadFile(ntPath)
	if err != nil {
		t.Skipf("real node_traversal.go not found: %v", err)
	}
	alreadyLowered := strings.Contains(string(origBytes), "func (this_ *NodeTraversal) traverseBranchBlk_(")

	// 1. In-memory pipeline targeting ONLY node_traversal.go (ZERO disk writes to gen/jscomp!)
	filePaths, err := listPackageGoFiles(realJscompDir, false)
	if err != nil {
		t.Fatal(err)
	}
	memFiles := make([]memFile, len(filePaths))
	ntIdx := -1
	for i, fp := range filePaths {
		raw, rErr := os.ReadFile(fp)
		if rErr != nil {
			t.Fatal(rErr)
		}
		isTarget := fp == ntPath
		if isTarget {
			ntIdx = i
		}
		memFiles[i] = memFile{name: fp, src: raw, target: isTarget}
	}
	if ntIdx < 0 {
		t.Fatalf("node_traversal.go not found in gen/jscomp package file list")
	}

	passList := []Pass{NewNodeTraversalLoweringPass()}
	imp1 := NewModuleImporter(realJscompDir)
	outFiles, stats1, err := executePipeline("github.com/paullepoulpe/jsgo/transpiled/gen/jscomp", memFiles, passList, imp1)
	if err != nil {
		t.Fatalf("executePipeline Run 1 on real node_traversal.go failed: %v", err)
	}
	if !alreadyLowered && stats1.TotalSites < 2 {
		t.Fatalf("expected >= 2 sites (traverseBranch + traverseChildren) on real node_traversal.go, got %+v", stats1)
	}
	gotNT := string(outFiles[ntIdx].src)
	assertWave47BLoweredNodeTraversalShape(t, gotNT)

	// 2. Round-2 idempotency on the in-memory lowered real node_traversal.go (also type-checks lowered output!)
	imp2 := NewModuleImporter(realJscompDir)
	_, stats2, err := executePipeline("github.com/paullepoulpe/jsgo/transpiled/gen/jscomp", outFiles, passList, imp2)
	if err != nil {
		t.Fatalf("executePipeline Run 2 on real node_traversal.go failed: %v", err)
	}
	if stats2.TotalSites != 0 || stats2.FilesModified != 0 {
		t.Fatalf("Run 2 on real node_traversal.go NOT idempotent: %+v", stats2)
	}

	// 3. DryRun on real gen/jscomp via RunOnPackageDir
	dryStats, err := RunOnPackageDir(realJscompDir, passList, Options{DryRun: true})
	if err != nil {
		t.Fatalf("DryRun on real gen/jscomp failed: %v", err)
	}
	if !alreadyLowered && dryStats.TotalSites < 2 {
		t.Fatalf("expected >= 2 DryRun sites on real gen/jscomp, got %+v", dryStats)
	}
	t.Logf("Real node_traversal.go: Run 1 sites=%d (modified=%d), Run 2 sites=%d, DryRun totalSites=%d filesModified=%d byPass=%v",
		stats1.TotalSites, stats1.FilesModified, stats2.TotalSites, dryStats.TotalSites, dryStats.FilesModified, dryStats.ByPass)
}
