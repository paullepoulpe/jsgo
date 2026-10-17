package jutil

// AppendCharSequence is Writer.append(CharSequence), making a
// *BufferedWriter an Appendable (java.io.Writer implements Appendable).
func (w *BufferedWriter) AppendCharSequence(csq string) { w.Write(csq) }

// AppendCharacter is Writer.append(char) (Appendable).
func (w *BufferedWriter) AppendCharacter(c uint16) { w.Write(StringFromChar(c)) }
