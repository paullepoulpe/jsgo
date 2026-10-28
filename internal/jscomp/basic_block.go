// Port of com/google/javascript/jscomp/BasicBlock.java @ closure-compiler v20261021

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

package jscomp

import (
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

var nextBasicBlockRootID atomic.Uint32

// BasicBlock represents a section of code that is uninterrupted by control
// structures (conditional or iterative logic).
type BasicBlock struct {
	parent *BasicBlock

	root   rhino.Node
	depth  int32
	rootID uint32

	// Whether this block denotes a function scope.
	isFunction bool

	// Whether this block denotes a loop.
	isLoop bool

	hasFnAncestor   bool
	hasLoopBeforeFn bool
}

// NewBasicBlock creates a new block.
//
// parent is the containing block; root is the root node of the block.
//
// Java: BasicBlock(BasicBlock, Node)
func NewBasicBlock(parent *BasicBlock, root rhino.Node) *BasicBlock {
	return initBasicBlock(&BasicBlock{}, parent, root)
}

func initBasicBlock(b *BasicBlock, parent *BasicBlock, root rhino.Node) *BasicBlock {
	b.parent = parent
	b.root = root
	if parent != nil {
		b.depth = parent.depth + 1
		b.rootID = parent.rootID
	} else {
		b.depth = 0
		b.rootID = nextBasicBlockRootID.Add(1)
	}

	b.isFunction = root.IsFunction()

	if p := root.GetParent(); !p.IsNil() {
		pType := p.GetToken()
		b.isLoop = pType == rhino.Token_DO ||
			pType == rhino.Token_WHILE ||
			pType == rhino.Token_FOR ||
			pType == rhino.Token_FOR_OF ||
			pType == rhino.Token_FOR_AWAIT_OF ||
			pType == rhino.Token_FOR_IN
	} else {
		b.isLoop = false
	}

	if b.isFunction {
		b.hasFnAncestor = true
		b.hasLoopBeforeFn = false
	} else if parent != nil {
		b.hasFnAncestor = parent.hasFnAncestor
		b.hasLoopBeforeFn = b.isLoop || parent.hasLoopBeforeFn
	} else {
		b.hasFnAncestor = false
		b.hasLoopBeforeFn = b.isLoop
	}
	return b
}

// Java: BasicBlock.getParent()
func (b *BasicBlock) GetParent() *BasicBlock {
	return b.parent
}

// IsGlobalScopeBlock determines whether this block is equivalent to the
// very first block that is created when reference collection traversal
// enters global scope. Note that when traversing a single script in a
// hot-swap fashion a new instance of BasicBlock is created.
//
// Java: BasicBlock.isGlobalScopeBlock()
func (b *BasicBlock) IsGlobalScopeBlock() bool {
	return b.parent == nil
}

// ProvablyExecutesBefore determines whether this block is guaranteed to
// begin executing before the given block does.
//
// Java: BasicBlock.provablyExecutesBefore(BasicBlock)
func (b *BasicBlock) ProvablyExecutesBefore(thatBlock *BasicBlock) bool {
	if thatBlock == b {
		return true
	}
	if b.parent == nil {
		return thatBlock.parent == nil || thatBlock.rootID == b.rootID
	}
	// If thatBlock is a descendant of this block, and there are no hoisted
	// blocks between them, then this block must start before thatBlock.
	bDepth := b.depth
	if thatBlock.depth <= bDepth || thatBlock.rootID != b.rootID {
		return false
	}
	currentBlock := thatBlock.parent
	for currentBlock.depth > bDepth {
		currentBlock = currentBlock.parent
	}
	return currentBlock == b
}

// Java: BasicBlock.isFunction()
func (b *BasicBlock) IsFunction() bool {
	return b.isFunction
}

// Java: BasicBlock.isLoop()
func (b *BasicBlock) IsLoop() bool {
	return b.isLoop
}

// Java: BasicBlock.getRoot()
func (b *BasicBlock) GetRoot() rhino.Node {
	return b.root
}

// Java: BasicBlock.toString()
func (b *BasicBlock) String() string {
	return "BasicBlock @ " + b.root.String()
}
