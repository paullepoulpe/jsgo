package jutil

import "strconv"

// Go-only: language-adapter — JsonWriter.value(long) and JsonWriter.flush()
// (Gson), extending gson_json_writer.go.

// ValueLong encodes value.
//
// Java: JsonWriter.value(long)
func (w *JsonWriter) ValueLong(value int64) *JsonWriter {
	w.writeDeferredName()
	w.beforeValue()
	w.write(strconv.FormatInt(value, 10))
	return w
}

// Flush ensures all buffered data is written to the underlying Writer (when
// it has a Flush method) and flushes that writer.
//
// Java: JsonWriter.flush()
func (w *JsonWriter) Flush() {
	if w.stackSize == 0 {
		panic(NewIllegalStateExceptionMsg("JsonWriter is closed."))
	}
	if f, ok := w.out.(interface{ Flush() }); ok {
		f.Flush()
	}
}
