package graph

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type embeddedLinkedGraph[N, E comparable] struct {
	LinkedDirectedGraph[N, E]
}

func newEmbeddedLinkedGraph[N, E comparable]() *embeddedLinkedGraph[N, E] {
	g := &embeddedLinkedGraph[N, E]{}
	g.Init(true, true)
	return g
}

func TestFixedPointGraphTraversal_EmbeddedLinkedGraph(t *testing.T) {
	g := newEmbeddedLinkedGraph[string, string]()
	g.CreateNode("a")
	g.CreateNode("b")
	g.CreateNode("c")
	g.Connect("a", "->", "b")
	g.Connect("b", "->", "c")

	reached := map[string]bool{"a": true}
	tr := FixedPointGraphTraversal_NewTraversal[string, string](
		FixedPointGraphTraversal_EdgeCallbackFunc[string, string](func(src, _, dst string) bool {
			if reached[src] && !reached[dst] {
				reached[dst] = true
				return true
			}
			return false
		}))
	tr.ComputeFixedPointFrom(g, "a")
	assertTrue(t, "b reached", reached["b"])
	assertTrue(t, "c reached", reached["c"])
}

func TestFixedPointGraphTraversal_ReverseAndMissingEntryHandling(t *testing.T) {
	g := newEmbeddedLinkedGraph[string, string]()
	g.CreateNode("a")
	g.CreateNode("b")
	g.Connect("a", "->", "b")

	// Missing entry throws NullPointerException.
	tr := FixedPointGraphTraversal_NewReverseTraversal[string, string](
		FixedPointGraphTraversal_EdgeCallbackFunc[string, string](func(src, _, dst string) bool {
			return false
		}))
	assertThrows[jutil.INullPointerException](t, func() {
		tr.ComputeFixedPointFrom(g, "missing")
	})

	// Subsequent traversal on the same graph still works cleanly.
	var visited []string
	trOnce := FixedPointGraphTraversal_NewReverseTraversal[string, string](
		FixedPointGraphTraversal_EdgeCallbackFunc[string, string](func(src, _, dst string) bool {
			if len(visited) == 0 {
				visited = append(visited, src+"<-"+dst)
				return true
			}
			return false
		}))
	trOnce.ComputeFixedPoint(g)
	assertSliceEqual(t, "reverse traversal after cleanup", visited, []string{"b<-a"})
}
