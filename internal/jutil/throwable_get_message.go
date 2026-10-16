package jutil

// GetMessageOrNull is Throwable.getMessage() with Java null as nil:
// GetMessage returns "" for both a null and an empty message, which is not
// enough where Java concatenates the message ("null" + ...).
func GetMessageOrNull(t Throwable) *string {
	msg := t.GetMessage()
	if msg == "" && !t.throwable().hasMsg {
		return nil
	}
	return &msg
}
