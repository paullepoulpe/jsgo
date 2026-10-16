package jutil

// Go-only: language-adapter — java.lang.Appendable. An IOException is
// thrown as a panic (jutil.IOException), like every other Java exception.

// Appendable is java.lang.Appendable.
type Appendable interface {
	// AppendCharSequence is Appendable.append(CharSequence).
	AppendCharSequence(csq string)
	// AppendCharacter is Appendable.append(char).
	AppendCharacter(c uint16)
}

// AppendCharSequence is StringBuilder.append(CharSequence) as an
// Appendable.
func (b *StringBuilder) AppendCharSequence(csq string) { b.Append(csq) }

// AppendCharacter is StringBuilder.append(char) as an Appendable.
func (b *StringBuilder) AppendCharacter(c uint16) { b.AppendChar(c) }
