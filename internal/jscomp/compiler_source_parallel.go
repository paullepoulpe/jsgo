package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// parallelSource owns one printer's text, relative source mappings and license
// events. No output builder, source map or actual license tracker is shared.
type parallelSource struct {
	code     *CodePrinter_SourceAndMappings
	licenses parallelSourceLicenses
}

type parallelSourceLicenses struct {
	scriptOnly  bool
	lastFile    string
	hasLastFile bool
	nodes       []*rhino.Node
}

func (l *parallelSourceLicenses) TrackLicensesForNode(node *rhino.Node) {
	if (node.IsRoot() || node.IsScript()) != l.scriptOnly {
		return
	}
	file := node.GetStaticSourceFile()
	if file == nil {
		return
	}
	name := file.GetName()
	if l.hasLastFile && name == l.lastFile {
		return
	}
	l.lastFile = name
	l.hasLastFile = true
	l.nodes = append(l.nodes, node)
}

func (l *parallelSourceLicenses) EmitLicenses() *jutil.LinkedHashSet[string] {
	panic("printer must not emit licenses")
}

// toSourceRoots computes independent printers concurrently, then publishes in
// input order. Known trackers only observe the eligible source-file transitions
// recorded here; custom trackers retain their exact sequential callback path.
// All compiler state, source-map offsets, license de-duplication and chunk graph
// processing are still accessed exclusively by the calling compiler goroutine.
func (c *Compiler) toSourceRoots(roots []*rhino.Node, tracker CodePrinter_LicenseTracker,
	emit func(int, *rhino.Node, *parallelSource)) {
	workers := c.options.GetNumParallelThreads()
	scriptOnly := false
	eligible := true
	switch tracker.(type) {
	case *Compiler_ScriptNodeLicensesOnlyTracker:
		scriptOnly = true
	case *Compiler_SingleBinaryLicenseTracker, *Compiler_ChunkGraphAwareLicenseTracker:
	default:
		eligible = false
	}
	if workers <= 1 || len(roots) <= 1 || !eligible || c.options.ShouldGenerateTypedExterns() || c.typedAstFilesystem != nil {
		for i, root := range roots {
			emit(i, root, nil)
		}
		return
	}
	// Bound private text/mapping buffers to one window, rather than retaining
	// a second whole output for a many-file bundle. Publishing a window before
	// starting the next also releases its AST/mapping references promptly.
	workers = min(workers, len(roots))
	results := make([]parallelSource, workers)
	for start := 0; start < len(roots); start += workers {
		count := min(workers, len(roots)-start)
		failures := parallelPrebuild(count, workers, func(slot int) {
			i := start + slot
			result := &results[slot]
			result.licenses.scriptOnly = scriptOnly
			result.code = c.toSourceAndMappings(roots[i], i == 0, &result.licenses)
		})
		for slot := range count {
			if failures[slot] != nil {
				panic(failures[slot])
			}
			i := start + slot
			emit(i, roots[i], &results[slot])
			results[slot] = parallelSource{}
		}
	}
}
