package jutil

// Throwable_PrintStackTrace is Throwable.printStackTrace(PrintStream): the
// throwable's toString(), then "Caused by: " and the toString() of each
// cause in turn.
//
// Go-only: language-adapter — a Go panic carries no Java stack frames, so
// the "\tat ..." and "\t... n more" lines that Java prints after each of
// these lines are absent.
func Throwable_PrintStackTrace(t Throwable, s *PrintStream) {
	s.Println(t.Error())
	seen := map[Throwable]bool{t: true}
	for cause := t.GetCause(); cause != nil; cause = cause.GetCause() {
		if seen[cause] {
			s.Println("\t[CIRCULAR REFERENCE: " + cause.Error() + "]")
			break
		}
		seen[cause] = true
		s.Println("Caused by: " + cause.Error())
	}
}
