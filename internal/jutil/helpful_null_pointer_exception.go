package jutil

// The JVM's "helpful" NullPointerException messages (JEP 358), written from
// the JEP and checked black-box against the pinned compiler's output on
// JDK 21 (tools/clie2e).
//
// When the JVM itself throws a NullPointerException (a method call, a field
// or array access on null), getMessage() describes the bytecode that failed,
// for example `Cannot invoke "com.google.javascript.rhino.Node.getSecondChild()"
// because "clazz" is null`. With -XX:-ShowCodeDetailsInExceptionMessages it
// is null instead. The option is on by default since JDK 15, so
// `java -jar closure-compiler.jar` prints the message wherever one reaches
// its output, as in an INTERNAL COMPILER ERROR report.
//
// A Go nil dereference carries no message (AsThrowable): the JVM derives
// it from bytecode. So a port spells the message out, with
// NewHelpfulNullPointerException, at each site where Java's output shows
// it. The oracle dumps made with the option off (or that replace these
// messages by "null") turn it off with SetShowCodeDetailsInExceptionMessages.

// showCodeDetailsInExceptionMessages is -XX:[+-]ShowCodeDetailsInExceptionMessages.
var showCodeDetailsInExceptionMessages = true

// SetShowCodeDetailsInExceptionMessages sets what
// -XX:+ShowCodeDetailsInExceptionMessages (on) or
// -XX:-ShowCodeDetailsInExceptionMessages (off) sets on the JVM, and returns
// the previous setting. It is on by default, as on JDK 21.
func SetShowCodeDetailsInExceptionMessages(on bool) (was bool) {
	was = showCodeDetailsInExceptionMessages
	showCodeDetailsInExceptionMessages = on
	return was
}

// NewHelpfulNullPointerException is the NullPointerException that the JVM
// throws on a null dereference whose helpful message is detail: its message
// is detail, or null when the option is off.
func NewHelpfulNullPointerException(detail string) *NullPointerException {
	if !showCodeDetailsInExceptionMessages {
		return NewNullPointerException()
	}
	return NewNullPointerExceptionMsg(detail)
}
