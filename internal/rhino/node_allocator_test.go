package rhino

import (
	"math/big"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"weak"

	"github.com/paullepoulpe/jsgo/internal/colors"
)

func TestNodeAllocatorCloneMatchesOriginal(t *testing.T) {
	cooked := "cooked"
	pool := NewRhinoStringPool_LazyInternedStringList([]string{"lazy"})
	for _, source := range []*Node{
		NewNode(Token_SCRIPT), Node_NewStringWithToken(Token_NAME, "name"),
		Node_NewStringFromStringPool(Token_NAME, pool, 0),
		Node_NewNumber(0.25), Node_NewBigInt(big.NewInt(12345)),
		Node_NewTemplateLitString(&cooked, "raw"), Node_NewTemplateLitString(nil, "bad\\x"),
	} {
		source.SetLinenoCharno(3, 7)
		source.SetLength(12)
		source.SetColor(colors.StandardColors_NUMBER)
		source.SetSourceFileForTesting("allocator.js")
		name := "original"
		source.SetOriginalName(&name)
		source.AddChildToBack(NewNode(Token_EMPTY))
		var allocator NodeAllocator
		got, want := allocator.CloneNode(source), source.CloneNode()
		if !reflect.DeepEqual(got, want) || got == source || got.HasChildren() || got.GetParent() != nil {
			t.Fatalf("owned clone differs for %v: %s, want %s", source.GetToken(), got, want)
		}
		if got.propListHead != want.propListHead {
			t.Fatal("clone stopped sharing immutable properties")
		}
	}
}

func discardedOwnedGraph() weak.Pointer[Node] {
	var allocator NodeAllocator
	root := allocator.NewNode(Token_SCRIPT)
	for range 1000 {
		child := allocator.NewString(Token_NAME, "discarded")
		root.AddChildToBack(child)
	}
	return weak.Make(root)
}

func TestNodeAllocatorDiscardDoesNotShareGlobalOrOtherOwner(t *testing.T) {
	before := NewNode(Token_EMPTY)
	discarded := discardedOwnedGraph()
	var survivorAllocator NodeAllocator
	survivor := survivorAllocator.NewNode(Token_SCRIPT)
	survivor.AddChildToBack(survivorAllocator.NewNumber(1))
	after := NewNode(Token_EMPTY)
	for attempt := 0; attempt < 10; attempt++ {
		runtime.GC()
		if discarded.Value() == nil {
			runtime.KeepAlive(before)
			runtime.KeepAlive(after)
			runtime.KeepAlive(survivor)
			return
		}
	}
	t.Fatal("a surviving independent allocation retained a discarded graph")
}

func TestNodeAllocatorIndependentConcurrentOwners(t *testing.T) {
	const workers, count = 4, 1000
	groups := make([][]*Node, workers)
	var wait sync.WaitGroup
	for worker := range groups {
		wait.Add(1)
		go func(worker int) {
			defer wait.Done()
			var allocator NodeAllocator
			for range count {
				groups[worker] = append(groups[worker], allocator.NewNumber(float64(worker)))
			}
		}(worker)
	}
	wait.Wait()
	seen := make(map[*Node]bool, workers*count)
	for worker, group := range groups {
		for _, node := range group {
			if seen[node] || !node.IsNumber() || node.GetDouble() != float64(worker) {
				t.Fatal("independent allocation streams shared or corrupted nodes")
			}
			seen[node] = true
		}
	}
}
