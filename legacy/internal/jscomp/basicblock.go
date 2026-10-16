// Port of com/google/javascript/jscomp/BasicBlock.java @ closure-compiler v20260415

package jscomp

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/legacy/internal/rhino"
)

// BasicBlock represents a section of code that is uninterrupted by control
// structures (conditional or iterative logic).
type BasicBlock struct {
	parent *BasicBlock

	root *rhino.Node

	// isFunction indicates whether this block denotes a function scope.
	isFunction bool

	// isLoop indicates whether this block denotes a loop.
	isLoop bool
}

// NewBasicBlock creates a new block.
//
// parent is the containing block.
// root is the root node of the block.
func NewBasicBlock(parent *BasicBlock, root *rhino.Node) *BasicBlock {
	b := &BasicBlock{
		parent:     parent,
		root:       root,
		isFunction: root.Token() == rhino.FUNCTION,
	}

	if root.Parent() != nil {
		pType := root.Parent().Token()
		b.isLoop = pType == rhino.DO_WHILE ||
			pType == rhino.WHILE ||
			pType == rhino.FOR ||
			pType == rhino.FOR_OF ||
			pType == rhino.FOR_AWAIT_OF ||
			pType == rhino.FOR_IN
	}

	return b
}

func (b *BasicBlock) GetParent() *BasicBlock {
	return b.parent
}

// IsGlobalScopeBlock determines whether this block is equivalent to the very
// first block that is created when reference collection traversal enters global
// scope. Note that when traversing a single script in a hot-swap fashion a new
// instance of BasicBlock is created.
//
// Returns true if this is global scope block.
func (b *BasicBlock) IsGlobalScopeBlock() bool {
	return b.GetParent() == nil
}

// ProvablyExecutesBefore determines whether this block is guaranteed to begin
// executing before the given block does.
func (b *BasicBlock) ProvablyExecutesBefore(thatBlock *BasicBlock) bool {
	// If thatBlock is a descendant of this block, and there are no hoisted
	// blocks between them, then this block must start before thatBlock.
	var currentBlock *BasicBlock
	for currentBlock = thatBlock; currentBlock != nil && currentBlock != b; currentBlock = currentBlock.GetParent() {
	}

	if currentBlock == b {
		return true
	}
	return b.IsGlobalScopeBlock() && thatBlock.IsGlobalScopeBlock()
}

func (b *BasicBlock) IsFunction() bool {
	return b.isFunction
}

func (b *BasicBlock) IsLoop() bool {
	return b.isLoop
}

func (b *BasicBlock) GetRoot() *rhino.Node {
	return b.root
}

func (b *BasicBlock) String() string {
	return fmt.Sprintf("BasicBlock @ %v", b.root)
}
