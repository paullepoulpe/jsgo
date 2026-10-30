// Port of com/google/javascript/jscomp/ReferenceCollector.java @ closure-compiler v20261021

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
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

type basicBlockEntry struct {
	root  rhino.Node
	block *BasicBlock
}

type savedBlockFrame struct {
	base int
	len  int
}

type basicBlockTracker struct {
	entries          []basicBlockEntry
	entriesBuf       [16]basicBlockEntry
	frameBase        int
	blocks           jutil.Slab[BasicBlock]
	collectedHoisted map[rhino.Node]struct{}
}

func (t *basicBlockTracker) push(root rhino.Node) {
	if cap(t.entries) == 0 {
		t.entries = t.entriesBuf[:0]
	}
	t.entries = append(t.entries, basicBlockEntry{root: root})
}

func (t *basicBlockTracker) pop(root rhino.Node) {
	if len(t.entries) <= t.frameBase {
		panic(jutil.NewNoSuchElementException())
	}
	idx := len(t.entries) - 1
	lastRoot := t.entries[idx].root
	t.entries[idx] = basicBlockEntry{}
	t.entries = t.entries[:idx]
	jutil.CheckStateMsg1(root == lastRoot, "%s", root)
}

func (t *basicBlockTracker) popUnchecked() {
	if len(t.entries) <= t.frameBase {
		panic(jutil.NewIndexOutOfBoundsException("Index -1 out of bounds for length 0"))
	}
	idx := len(t.entries) - 1
	t.entries[idx] = basicBlockEntry{}
	t.entries = t.entries[:idx]
}

func (t *basicBlockTracker) setTopBlock(b *BasicBlock) {
	t.entries[len(t.entries)-1].block = b
}

func (t *basicBlockTracker) currentBlock() *BasicBlock {
	if len(t.entries) <= t.frameBase {
		panic(jutil.NewNoSuchElementException())
	}
	idx := len(t.entries) - 1
	if b := t.entries[idx].block; b != nil {
		return b
	}
	return t.materializeAt(idx)
}

func (t *basicBlockTracker) materializeAt(idx int) *BasicBlock {
	base := t.frameBase
	i := idx
	for i > base && t.entries[i-1].block == nil {
		i--
	}
	var parent *BasicBlock
	if i > base {
		parent = t.entries[i-1].block
	}
	for ; i <= idx; i++ {
		parent = initBasicBlock(t.blocks.New(), parent, t.entries[i].root)
		t.entries[i].block = parent
	}
	return parent
}

func (t *basicBlockTracker) saveHoistedFrame(varScope *Scope) savedBlockFrame {
	base := t.frameBase
	oldLen := len(t.entries)
	saved := savedBlockFrame{base: base, len: oldLen}
	if varScope.IsGlobal() {
		if oldLen <= base {
			panic(jutil.NewNoSuchElementException())
		}
		if t.entries[base].block == nil {
			t.materializeAt(base)
		}
		entry0 := t.entries[base]
		t.frameBase = oldLen
		t.entries = append(t.entries, entry0)
	} else {
		rootNode := varScope.GetRootNode()
		end := oldLen - base
		for k := 0; k < oldLen-base; k++ {
			if t.entries[base+k].root == rootNode {
				end = k + 1
				break
			}
		}
		if end > 0 && t.entries[base+end-1].block == nil {
			t.materializeAt(base + end - 1)
		}
		t.frameBase = oldLen
		t.entries = append(t.entries, t.entries[base:base+end]...)
	}
	return saved
}

func (t *basicBlockTracker) pushRootBlockFrame(rootBlock *BasicBlock) savedBlockFrame {
	if cap(t.entries) == 0 {
		t.entries = t.entriesBuf[:0]
	}
	oldLen := len(t.entries)
	saved := savedBlockFrame{base: t.frameBase, len: oldLen}
	t.frameBase = oldLen
	t.entries = append(t.entries, basicBlockEntry{root: rootBlock.root, block: rootBlock})
	return saved
}

func (t *basicBlockTracker) restoreFrame(saved savedBlockFrame) {
	clear(t.entries[saved.len:])
	t.entries = t.entries[:saved.len]
	t.frameBase = saved.base
}

func (t *basicBlockTracker) markHoistedCollected(fnNode rhino.Node) bool {
	if _, ok := t.collectedHoisted[fnNode]; ok {
		return false
	}
	if t.collectedHoisted == nil {
		t.collectedHoisted = make(map[rhino.Node]struct{})
	}
	t.collectedHoisted[fnNode] = struct{}{}
	return true
}

func (t *basicBlockTracker) isHoistedCollected(fnNode rhino.Node) bool {
	_, ok := t.collectedHoisted[fnNode]
	return ok
}

func (t *basicBlockTracker) resetHoistedCollected() {
	if t.collectedHoisted != nil {
		clear(t.collectedHoisted)
	}
}

func (t *basicBlockTracker) uncollectedHoistedTarget(v *Var) (rhino.Node, bool) {
	nameNode := v.nameNode
	if nameNode.IsNil() {
		return rhino.NilNode, false
	}
	fnNode := nameNode.GetParent()
	if fnNode.IsNil() || !fnNode.IsFunction() {
		return rhino.NilNode, false
	}
	if t.isHoistedCollected(fnNode) || !rcIsHoistedFunctionDeclaration(fnNode) {
		return rhino.NilNode, false
	}
	return fnNode, true
}

func rcIsHoistedFunctionDeclaration(fnNode rhino.Node) bool {
	parent := fnNode.GetParent()
	if parent.IsNil() {
		return false
	}
	switch parent.GetToken() {
	case rhino.Token_SCRIPT, rhino.Token_MODULE_BODY, rhino.Token_EXPORT:
		first := fnNode.GetFirstChild()
		return !first.IsNil() && first.IsName() && first.GetStringID() != rhino.RhinoStringPool_EmptyID
	case rhino.Token_BLOCK, rhino.Token_LABEL, rhino.Token_DECLARE, rhino.Token_NAMESPACE_ELEMENTS:
		if gp := parent.GetParent(); !gp.IsNil() && gp.IsFunction() {
			first := fnNode.GetFirstChild()
			return !first.IsNil() && first.IsName() && first.GetStringID() != rhino.RhinoStringPool_EmptyID
		}
	}
	return false
}
