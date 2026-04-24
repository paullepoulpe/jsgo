// Port of com/google/javascript/jscomp/DotFormatter.java
// @ closure-compiler v20260415
//
// DotFormatter converts an AST (and optionally its ControlFlowGraph) to a
// DOT-format graph for visualization with Graphviz.
//
// Typical usage:
//
//	dot, err := DotFormatterToDot(scriptNode, nil)
//	// or with CFG:
//	dot, err := DotFormatterToDot(scriptNode, cfg)
//
// Mirrors the Java final class DotFormatter in the closure-compiler.

package jscomp

import (
	"fmt"
	"sort"
	"strings"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

const (
	dotIndent          = "  "
	dotArrow           = " -> "
	dotMaxLabelLength  = 10
)

// DotFormatter formats an AST and optional CFG as a DOT graph string.
// Mirrors com.google.javascript.jscomp.DotFormatter in Java.
type DotFormatter struct {
	// assignments maps each rhino.Node pointer to a unique integer key.
	// Mirrors DotFormatter.assignments (LinkedHashMap<Node, Integer>) in Java.
	assignments map[*rhino.Node]int

	// keyCount is the next key to assign.
	// Mirrors DotFormatter.keyCount in Java.
	keyCount int

	// builder accumulates the DOT output.
	// Mirrors DotFormatter.builder (Appendable) in Java.
	builder strings.Builder

	// cfg is the optional ControlFlowGraph whose edges are rendered in red.
	// Mirrors DotFormatter.cfg in Java.
	cfg *ControlFlowGraph

	// printAnnotations controls whether per-node annotations are appended.
	// Mirrors DotFormatter.printAnnotations in Java.
	printAnnotations bool
}

// newDotFormatter creates a DotFormatter for testing internal methods.
// Mirrors DotFormatter.newInstanceForTesting() in Java.
func newDotFormatter() *DotFormatter {
	return &DotFormatter{
		assignments: make(map[*rhino.Node]int),
	}
}

// DotFormatterToDot converts an AST rooted at n to a DOT graph string.
// If cfg is non-nil, CFG edges are rendered as red arrows.
// Mirrors DotFormatter.toDot(Node n, ControlFlowGraph inCFG) in Java.
func DotFormatterToDot(n *rhino.Node, cfg *ControlFlowGraph) (string, error) {
	df := &DotFormatter{
		assignments:      make(map[*rhino.Node]int),
		cfg:              cfg,
		printAnnotations: false,
	}
	if err := df.formatPreamble(); err != nil {
		return "", err
	}
	if err := df.traverseNodes(n); err != nil {
		return "", err
	}
	if err := df.formatConclusion(); err != nil {
		return "", err
	}
	return df.builder.String(), nil
}

// DotFormatterASTOnly converts an AST rooted at n to DOT (no CFG edges).
// Mirrors DotFormatter.toDot(Node n) in Java.
func DotFormatterASTOnly(n *rhino.Node) (string, error) {
	return DotFormatterToDot(n, nil)
}

// DotFormatterAppend builds a DOT representation and appends it to the given
// strings.Builder. Mirrors DotFormatter.appendDot in Java.
func DotFormatterAppend(n *rhino.Node, cfg *ControlFlowGraph, buf *strings.Builder) error {
	dot, err := DotFormatterToDot(n, cfg)
	if err != nil {
		return err
	}
	buf.WriteString(dot)
	return nil
}

// DotFormatterCFG renders only the CFG (without AST structure) as a DOT graph.
// Each CfgNode is a box labelled with its AST token; edges carry branch labels.
// This is a Go-only convenience not present in the Java DotFormatter.
// Go-only: Go port adds a CFG-only renderer for debugging.
func DotFormatterCFG(cfg *ControlFlowGraph) string {
	if cfg == nil {
		return "digraph CFG {}\n"
	}
	var b strings.Builder
	b.WriteString("digraph CFG {\n")
	b.WriteString(dotIndent + "node [color=lightblue2, style=filled];\n")

	// Assign stable IDs to nodes.
	ids := make(map[*CfgNode]int)
	for i, n := range cfg.Nodes {
		ids[n] = i
	}
	ids[cfg.Entry] = len(cfg.Nodes)
	ids[cfg.Exit] = len(cfg.Nodes) + 1

	label := func(n *CfgNode) string {
		if n == cfg.Entry {
			return "ENTRY"
		}
		if n == cfg.Exit {
			return "EXIT"
		}
		if n.AST != nil {
			s := n.AST.Token().String()
			str := n.AST.String_()
			if len(str) > dotMaxLabelLength {
				str = str[:dotMaxLabelLength]
			}
			if str != "" {
				s += "(" + str + ")"
			}
			return s
		}
		return n.Kind
	}

	// Emit node declarations.
	for _, n := range append(cfg.Nodes, cfg.Entry, cfg.Exit) {
		fmt.Fprintf(&b, "%snode%d [label=%q];\n", dotIndent, ids[n], label(n))
	}

	// Emit edges.
	type edge struct {
		from, to int
		label    string
	}
	var edges []edge
	for _, n := range append(cfg.Nodes, cfg.Entry, cfg.Exit) {
		fromID := ids[n]
		for _, e := range n.Successors {
			toID := ids[e.To]
			edges = append(edges, edge{fromID, toID, e.Branch.String()})
		}
	}
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].from != edges[j].from {
			return edges[i].from < edges[j].from
		}
		return edges[i].to < edges[j].to
	})
	for _, e := range edges {
		fmt.Fprintf(&b, "%snode%d%snode%d [label=%q, color=\"red\"];\n",
			dotIndent, e.from, dotArrow, e.to, e.label)
	}
	b.WriteString("}\n")
	return b.String()
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// traverseNodes walks the AST rooted at parent, emitting edges and node
// declarations into df.builder.
// Mirrors DotFormatter.traverseNodes(Node parent) in Java.
func (df *DotFormatter) traverseNodes(parent *rhino.Node) error {
	if parent == nil {
		return nil
	}
	keyParent, err := df.key(parent)
	if err != nil {
		return err
	}

	// Emit AST child edges.
	for _, child := range parent.Children() {
		keyChild, err := df.key(child)
		if err != nil {
			return err
		}
		fmt.Fprintf(&df.builder, "%s%s%s%s [weight=1];\n",
			dotIndent,
			dotFormatNodeName(keyParent),
			dotArrow,
			dotFormatNodeName(keyChild))
		if err := df.traverseNodes(child); err != nil {
			return err
		}
	}

	// Emit CFG flow edges if a CFG was provided and this node has a CfgNode.
	if df.cfg != nil {
		cfgNode := df.cfg.Get(parent)
		if cfgNode != nil {
			// Collect and sort the edge strings for deterministic output.
			var edgeStrs []string
			for _, e := range cfgNode.Successors {
				var toName string
				if e.To == df.cfg.Exit {
					toName = "RETURN"
				} else if e.To.AST != nil {
					keySucc, err := df.key(e.To.AST)
					if err != nil {
						return err
					}
					toName = dotFormatNodeName(keySucc)
				} else {
					// Synthetic node (LOOP_AFTER, etc.) — label it by kind.
					toName = fmt.Sprintf("%q", e.To.Kind)
				}
				edgeStrs = append(edgeStrs,
					fmt.Sprintf("%s%s%s%s [label=%q, fontcolor=\"red\", weight=0.01, color=\"red\"];\n",
						dotIndent,
						dotFormatNodeName(keyParent),
						dotArrow,
						toName,
						e.Branch.String()))
			}
			sort.Strings(edgeStrs)
			for _, s := range edgeStrs {
				df.builder.WriteString(s)
			}
		}
	}
	return nil
}

// key returns the integer key for node n, creating one if needed.
// It also emits the node declaration into df.builder the first time.
// Mirrors DotFormatter.key(Node n) in Java.
func (df *DotFormatter) key(n *rhino.Node) (int, error) {
	if k, ok := df.assignments[n]; ok {
		return k, nil
	}
	k := df.keyCount
	df.keyCount++
	df.assignments[n] = k

	// Node declaration: node<k> [label="TOKEN(label)"]
	label := n.Token().String()
	str := dotGetNodeLabel(n)
	if str != "" {
		label += str
	}
	// JSType annotation — our AST doesn't yet carry JSTypes on nodes.
	// When it does, append " : typeStr" here.

	nodeName := dotFormatNodeName(k)
	attrs := fmt.Sprintf("label=%q", label)
	if n.JSDocInfo() != nil {
		attrs += ` color="green"`
	}
	fmt.Fprintf(&df.builder, "%s%s [%s];\n", dotIndent, nodeName, attrs)
	return k, nil
}

// dotGetNodeLabel returns the parenthesised label string for name-bearing nodes.
// Mirrors DotFormatter.getNodeLabel(Node n) in Java.
func dotGetNodeLabel(n *rhino.Node) string {
	tok := n.Token()
	switch tok {
	case rhino.NAME, rhino.STRINGLIT, rhino.GETPROP:
		// carry the string value
	default:
		return ""
	}
	content := n.String_()
	if len(content) > dotMaxLabelLength {
		content = content[:dotMaxLabelLength]
	}
	if content == "" {
		return ""
	}
	return "(" + content + ")"
}

// dotFormatNodeName returns the DOT node identifier for an integer key.
// Mirrors DotFormatter.formatNodeName(Integer key) in Java.
func dotFormatNodeName(key int) string {
	return fmt.Sprintf("node%d", key)
}

// formatPreamble writes the DOT digraph header.
// Mirrors DotFormatter.formatPreamble() in Java.
func (df *DotFormatter) formatPreamble() error {
	_, err := fmt.Fprintf(&df.builder,
		"digraph AST {\n%snode [color=lightblue2, style=filled];\n",
		dotIndent)
	return err
}

// formatConclusion writes the closing brace.
// Mirrors DotFormatter.formatConclusion() in Java.
func (df *DotFormatter) formatConclusion() error {
	_, err := df.builder.WriteString("}\n")
	return err
}
