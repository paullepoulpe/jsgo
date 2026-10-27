// Package graph provides generic directed and undirected graph data structures
// and algorithms used by control-flow, data-flow, call-graph, dependency, and
// register/property coloring passes in jscomp.
//
// Core interfaces and implementations:
//   - [Graph], [DiGraph], [UndiGraph], [AdjacencyGraph], [SubGraph], [GraphvizGraph]
//   - [LinkedDirectedGraph] and [LinkedUndirectedGraph] — insertion-ordered
//     graphs backed by slab-allocated nodes and inline 2-edge lists
//   - [FixedPointGraphTraversal], [GraphReachability], [CheckPathsBetweenNodes]
//   - [DominatorTree], [LowestCommonAncestorFinder]
//   - [GraphColoring] / [GraphColoring_GreedyGraphColoring]
//   - [UnionFind] / [StandardUnionFind]
package graph
