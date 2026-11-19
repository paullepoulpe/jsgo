/*
 * Copyright 2017 The Closure Compiler Authors.
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
// Port of com/google/javascript/jscomp/BasicBlock.java (generated mechanically by j2go).

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/gen/j2gort"
	pkg_rhino "github.com/paullepoulpe/jsgo/transpiled/gen/rhino"
)

// Java class: com.google.javascript.jscomp.BasicBlock
type BasicBlock struct {
	parent     *BasicBlock    `gson:"parent"`
	root       pkg_rhino.Node `gson:"root"`
	idHash_    atomic.Int32
	isFunction bool `gson:"isFunction"`
	isLoop     bool `gson:"isLoop"`
}

func (this_ *BasicBlock) BindSelf_(s any) {
	_ = s
}
func (this_ *BasicBlock) IsJavaNil_() bool { return this_ == nil }
func (*BasicBlock) JavaClassInfo_() (string, string) {
	return "com.google.javascript.jscomp.BasicBlock", "BasicBlock"
}
func (this_ *BasicBlock) IdentityHash_() int32 {
	if this_ == nil {
		return 0
	}
	if h := this_.idHash_.Load(); h != 0 {
		return h
	}
	return j2gort.AssignIdentHash(&this_.idHash_)
}
func (*BasicBlock) HasIdentityEquals_() {}
func init() {
	j2gort.RegisterTypePredicate("com.google.javascript.jscomp.BasicBlock", "com.google.javascript.jscomp.BasicBlock", func(v_ any) bool {
		return j2gort.InstanceOf[*BasicBlock](v_)
	})
}

func (this_ *BasicBlock) initFields_BasicBlock() {
	_ = this_
}

// Java constructor: BasicBlock(com.google.javascript.jscomp.BasicBlock parent, com.google.javascript.rhino.Node root)
func NewBasicBlock(parent *BasicBlock, root pkg_rhino.Node) *BasicBlock {
	this_ := &BasicBlock{}
	InitBasicBlock(this_, parent, root)
	return this_
}
func allocPooled_BasicBlock(slab_ *j2gort.PooledSlab[BasicBlock], parent *BasicBlock, root pkg_rhino.Node) *BasicBlock {
	this_ := slab_.Alloc()
	InitBasicBlock(this_, parent, root)
	return this_
}

func InitBasicBlock(this_ *BasicBlock, parent *BasicBlock, root pkg_rhino.Node) {
	_ = parent
	_ = root
	this_.initFields_BasicBlock()
	this_.parent = parent
	this_.root = root
	this_.isFunction = root.IsFunction()
	if root.HasParent() {
		var pType pkg_rhino.Token = root.GetParent().GetToken()
		_ = pType
		this_.isLoop = ((((((pType == pkg_rhino.StaticGet_Token_DO()) || (pType == pkg_rhino.StaticGet_Token_WHILE())) || (pType == pkg_rhino.StaticGet_Token_FOR())) || (pType == pkg_rhino.StaticGet_Token_FOR_OF())) || (pType == pkg_rhino.StaticGet_Token_FOR_AWAIT_OF())) || (pType == pkg_rhino.StaticGet_Token_FOR_IN()))
	} else {
		this_.isLoop = false
	}
}

// Java method: BasicBlock.getParent()
func (this_ *BasicBlock) GetParent() *BasicBlock {
	return this_.parent
}

// Java method: BasicBlock.isGlobalScopeBlock()
func (this_ *BasicBlock) IsGlobalScopeBlock() bool {
	return (this_.GetParent() == nil)
}

// Java method: BasicBlock.provablyExecutesBefore(com.google.javascript.jscomp.BasicBlock thatBlock)
func (this_ *BasicBlock) ProvablyExecutesBefore(thatBlock *BasicBlock) bool {
	_ = thatBlock
	var currentBlock *BasicBlock = nil
	_ = currentBlock
	{
		currentBlock = thatBlock
		for ; (currentBlock != nil) && (currentBlock != this_); currentBlock = currentBlock.GetParent() {
		}
	}
	if currentBlock == this_ {
		return true
	}
	return (this_.IsGlobalScopeBlock() && thatBlock.IsGlobalScopeBlock())
}

// Java method: BasicBlock.isFunction()
func (this_ *BasicBlock) IsFunction() bool {
	return this_.isFunction
}

// Java method: BasicBlock.isLoop()
func (this_ *BasicBlock) IsLoop() bool {
	return this_.isLoop
}

// Java method: BasicBlock.getRoot()
func (this_ *BasicBlock) GetRoot() pkg_rhino.Node {
	return this_.root
}

// Java method: BasicBlock.toString()
func (this_ *BasicBlock) String() string {
	return ("BasicBlock @ " + j2gort.StrVal(this_.root))
}
