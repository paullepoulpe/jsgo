// Port of com/google/javascript/jscomp/CheckDebuggerStatement.java @ closure-compiler v20260415
//
// CheckDebuggerStatement checks for the presence of the "debugger" statement
// in JavaScript code. It is appropriate to use this statement while developing
// JavaScript; however, it is generally undesirable to include it in production
// code.
//
// DEBUGGER_STATEMENT_PRESENT is disabled by default (matches Java's
// DiagnosticType.disabled), so it only fires when the diagnostic group
// "checkDebuggerStatement" is explicitly enabled via CompilerOptions.

package jscomp

import (
	"github.com/paullepoulpe/jsgo/internal/diag"
	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// CheckDebuggerStatement checks for the presence of the "debugger" statement
// in JavaScript code. Implements AbstractPostOrderCallback + CompilerPass.
type CheckDebuggerStatement struct {
	path  string
	opts  *CompilerOptions
	diags []diag.Diagnostic
}

// NewCheckDebuggerStatement mirrors the Java constructor
// CheckDebuggerStatement(AbstractCompiler compiler).
func NewCheckDebuggerStatement(path string, opts *CompilerOptions) *CheckDebuggerStatement {
	return &CheckDebuggerStatement{path: path, opts: opts}
}

// Check walks the AST and returns any diagnostics. Mirrors process() +
// NodeTraversal.traverse() in Java, using a simple recursive walk instead of
// the full NodeTraversal infrastructure.
func (c *CheckDebuggerStatement) Check(script *rhino.Node) []diag.Diagnostic {
	// JSC_DEBUGGER_STATEMENT_PRESENT is disabled by default. Skip the walk
	// entirely when the group is off — matches Java's behaviour where a
	// DiagnosticType.disabled check never surfaces without opt-in.
	if !c.opts.IsOn(GroupCheckDebuggerStatement) {
		return nil
	}
	c.diags = nil
	c.walk(script)
	return c.diags
}

// walk performs a post-order traversal, mirroring AbstractPostOrderCallback.
func (c *CheckDebuggerStatement) walk(n *rhino.Node) {
	if n == nil {
		return
	}
	for _, child := range n.Children() {
		c.walk(child)
	}
	c.visit(n)
}

// visit implements the visitor logic. Mirrors
// visit(NodeTraversal t, Node n, Node parent) in Java.
func (c *CheckDebuggerStatement) visit(n *rhino.Node) {
	if n.Token() == rhino.DEBUGGER {
		c.diags = append(c.diags, diag.Diagnostic{
			Path:  c.path,
			Line:  n.Lineno(),
			Col:   n.Charno(),
			Code:  diag.JSC_DEBUGGER_STATEMENT_PRESENT,
			Level: diag.Warning,
			Message: "Using the debugger statement can halt your application if the user has a" +
				" JavaScript debugger running.\n" +
				"To disable this check when you want to do debugging, you can suppress this" +
				" message like this:\n" +
				"  /** @suppress {checkDebuggerStatement} */\n" +
				"  debugger;",
		})
	}
}
