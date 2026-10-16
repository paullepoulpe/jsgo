package jutil

// Go-only: language-adapter — com.google.gson.TypeAdapter (Gson 2.9.1, as
// shaded into the pinned closure-compiler.jar).

// TypeAdapter is com.google.gson.TypeAdapter: converts Java objects to and
// from JSON.
type TypeAdapter interface {
	// Write writes one JSON value (an array, object, string, number,
	// boolean or null) for value.
	//
	// Java: TypeAdapter.write(JsonWriter, T)
	Write(out *JsonWriter, value any)

	// Read reads one JSON value and converts it to a Java object.
	//
	// Java: TypeAdapter.read(JsonReader)
	Read(in *JsonReader) any
}
