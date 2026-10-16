package jutil

// Go-only: language-adapter — Guava's com.google.common.escape.Escaper.

// Escaper is com.google.common.escape.Escaper.
type Escaper interface {
	// Escape is Escaper.escape(String).
	Escape(s string) string
}

// Escaper_AsFunction is Escaper.asFunction().
func Escaper_AsFunction(e Escaper) func(string) string {
	return e.Escape
}
