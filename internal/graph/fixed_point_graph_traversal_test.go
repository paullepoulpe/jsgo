package graph

import (
	"testing"

	"github.com/paullepoulpe/jsgo/internal/jutil"
)

type embeddedLinkedGraph[N comparable, E any] struct {
	LinkedDirectedGraph[N, E]
}

func newEmbeddedLinkedGraph[N comparable, E any]() *embeddedLinkedGraph[N, E] {
	g := &embeddedLinkedGraph[N, E]{}
	g.InitLinkedDirectedGraph(g, true, true)
	return g
}

func TestFixedPointGraphTraversal_EmbeddedLinkedGraph(t *testing.T) {
	var g DiGraph[string, string] = newEmbeddedLinkedGraph[string, string]()
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

func TestFixedPointGraphTraversal_ReverseAndNilEntryHandling(t *testing.T) {
	var g DiGraph[string, string] = newEmbeddedLinkedGraph[string, string]()
	g.CreateNode("a")
	g.CreateNode("b")
	g.Connect("a", "->", "b")

	// Nil entrySet throws NullPointerException.
	tr := FixedPointGraphTraversal_NewReverseTraversal[string, string](
		FixedPointGraphTraversal_EdgeCallbackFunc[string, string](func(src, _, dst string) bool {
			return false
		}))
	assertThrows[jutil.INullPointerException](t, func() {
		tr.ComputeFixedPointFromSet(g, nil)
	})
	var typedNil *jutil.LinkedHashSet[string]
	assertThrows[jutil.INullPointerException](t, func() {
		tr.ComputeFixedPointFromSet(g, typedNil)
	})

	// Missing entry after valid entry visits valid entry first, then throws NullPointerException,
	// leaving inWorkSet flags clean for subsequent traversals.
	var visited []string
	trVisit := FixedPointGraphTraversal_NewTraversal[string, string](
		FixedPointGraphTraversal_EdgeCallbackFunc[string, string](func(src, _, dst string) bool {
			visited = append(visited, src+"->"+dst)
			return true
		}))
	entries := jutil.NewLinkedHashSet[string]()
	entries.Add("a")
	entries.Add("missing1")
	entries.Add("missing2")
	assertThrows[jutil.INullPointerException](t, func() {
		trVisit.ComputeFixedPointFromSet(g, entries)
	})
	assertSliceEqual(t, "visited before missing entry NPE", visited, []string{"a->b"})

	// Subsequent traversal on the same graph still works cleanly (defer cleaned inWorkSet on "b").
	visited = nil
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
