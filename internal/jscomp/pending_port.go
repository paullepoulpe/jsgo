package jscomp

// PendingPort is the panic value of the declarations that tools/pending
// generates in pending_gen.go for the Java members that are not ported yet
// (PORTING.md §1).
//
// Go-only: language-adapter — not a Java exception: it is
// jutil.Uncatchable, so ported catch blocks (jutil.Try) do not catch it.
// RunJUnit4 and SkipIfPendingPort report a test that reaches one as
// skipped.
type PendingPort struct {
	// Signature is the Java member, as on its `// Java:` line.
	Signature string
}

// Error renders "pending port: <Java signature>".
func (p *PendingPort) Error() string { return "pending port: " + p.Signature }

// Uncatchable marks PendingPort as jutil.Uncatchable.
func (p *PendingPort) Uncatchable() {}
