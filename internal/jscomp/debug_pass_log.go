package jscomp

// Go-only: measurement hook for the per-pass profile (cmd/jsgo,
// JSGO_PASS_LOG). It has no Java counterpart and never changes what the
// compiler does: when DebugPassLogHook is nil (the default) nothing below
// runs.

// DebugPassEvent describes the start or the end of one pass run (a
// PhaseOptimizer pass) or of one Compiler tracer section (parseInputs,
// toSource, ...).
type DebugPassEvent struct {
	Name   string
	Tracer bool // a Compiler.newTracer section rather than a PhaseOptimizer pass
	Start  bool
	// Loop is the 1-based ordinal of the fixed-point loop the pass runs in
	// (0 outside loops), Iter the 1-based iteration of that loop.
	Loop, Iter int
	// End events only: change reports made by the pass, and the number of
	// distinct change scopes (functions/scripts) it marked changed.
	Reports, ChangedScopes int
}

// DebugPassLogHook, when set, receives every DebugPassEvent.
var DebugPassLogHook func(ev DebugPassEvent)

const debugPassTimelineKey = "__debugPassLog__"

type debugPassChangeCounter struct{ n int }

func (h *debugPassChangeCounter) ReportChange() { h.n++ }

// debugPassStart reports the start of a PhaseOptimizer pass and returns the
// function that reports its end.
func (p *PhaseOptimizer) debugPassStart(name string) func() {
	ev := DebugPassEvent{Name: name, Start: true}
	if p.inLoop {
		ev.Loop, ev.Iter = p.debugLoopOrdinal, p.debugLoopIter
	}
	DebugPassLogHook(ev)
	counter := &debugPassChangeCounter{}
	p.changeTracker.AddChangeHandler(counter)
	p.changeTracker.changeTimeline.Mark(debugPassTimelineKey)
	return func() {
		p.changeTracker.RemoveChangeHandler(counter)
		scopes := p.changeTracker.changeTimeline.GetSince(debugPassTimelineKey)
		p.changeTracker.changeTimeline.Unmark(debugPassTimelineKey)
		ev.Start = false
		ev.Reports = counter.n
		ev.ChangedScopes = len(scopes)
		DebugPassLogHook(ev)
	}
}
