// Package sourcemap is a 1:1 port of com.google.debugging.sourcemap
// (closure-compiler v20261021): the source map V3 generator
// (SourceMapGeneratorV3), consumer (SourceMapConsumerV3), their JSON
// parsers (SourceMapObjectParser on Gson, and the SourceMapJsonLexer fast
// path), and the Base64 VLQ encoding. proto_mapping.go reconstructs the
// protobuf-generated com.google.debugging.sourcemap.proto.Mapping.
//
// Java Strings that may be null are *string; Java arrays that may be null
// are nil slices (a non-null empty array is a non-nil empty slice); a Java
// Appendable is an io.StringWriter whose write errors become IOExceptions.
// See PORTING.md for the translation rules.
package sourcemap
