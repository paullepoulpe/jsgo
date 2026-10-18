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
	scriptOnly bool
	lastFile   string
	nodes      []*rhino.Node
}

func (l *parallelSourceLicenses) TrackLicensesForNode(node *rhino.Node) {
	if (node.IsRoot() || node.IsScript()) != l.scriptOnly {
		return
	}
	file := node.GetStaticSourceFile()
	if file == nil || file.GetName() == l.lastFile {
		return
	}
	l.lastFile = file.GetName()
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
	if workers <= 1 || len(roots) <= 1 || !eligible || c.options.ShouldGenerateTypedExterns() {
		for i, root := range roots {
			emit(i, root, nil)
		}
		return
	}
	results := make([]parallelSource, len(roots))
	failures := parallelPrebuild(len(roots), workers, func(i int) {
		result := &results[i]
		result.licenses.scriptOnly = scriptOnly
		result.code = c.toSourceAndMappings(roots[i], i == 0, &result.licenses)
	})
	for i, root := range roots {
		if failures[i] != nil {
			panic(failures[i])
		}
		emit(i, root, &results[i])
		results[i] = parallelSource{}
	}
}
