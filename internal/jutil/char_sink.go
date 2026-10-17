package jutil

import "os"

// Go-only: language-adapter — Guava com.google.common.io.CharSink, as
// returned by Files.asCharSink(File, Charset) (no FileWriteMode), written
// from Guava's documentation: write(CharSequence) opens a new
// OutputStreamWriter over a new FileOutputStream for the file (which
// creates the file, or truncates it when it exists), writes the characters
// and closes the stream. Only UTF-8 is supported (as by
// OutputStreamWriter).

// CharSink is Guava's CharSink over a file.
type CharSink struct {
	file    *File
	charset *Charset
}

// Files_AsCharSink is Guava Files.asCharSink(File, Charset) with no
// FileWriteMode.
func Files_AsCharSink(file *File, charset *Charset) *CharSink {
	return &CharSink{file: CheckNotNull(file), charset: CheckNotNull(charset)}
}

// OpenStream is CharSink.openStream(): a new OutputStreamWriter over a new
// FileOutputStream(file, false). It panics with FileNotFoundException when
// the file cannot be opened for writing.
func (c *CharSink) OpenStream() *OutputStreamWriter {
	if c.charset != StandardCharsets_UTF_8 {
		panic(NewUnsupportedOperationExceptionMsg("CharSink charset " + c.charset.Name()))
	}
	f, err := os.OpenFile(c.file.GetPath(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		panic(FileNotFoundExceptionFromError(err, c.file.GetPath()))
	}
	return NewOutputStreamWriter(f)
}

// Write is CharSink.write(CharSequence): the chars of charSequence are
// written to the file, replacing its contents. It panics with an
// IOException when the file cannot be written.
func (c *CharSink) Write(charSequence string) {
	out := c.OpenStream()
	defer out.Close()
	out.Write(charSequence)
}
