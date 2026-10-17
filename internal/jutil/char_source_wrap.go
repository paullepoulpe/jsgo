package jutil

// Go-only: language-adapter — Guava CharSource.wrap(CharSequence), written
// from Guava's documentation: wrap views the characters of a string as a
// source (CharSource.read() is in char_source.go).

// CharSource_Wrap is CharSource.wrap(CharSequence).
func CharSource_Wrap(charSequence string) *CharSource {
	return &CharSource{wrapped: &charSequence}
}
