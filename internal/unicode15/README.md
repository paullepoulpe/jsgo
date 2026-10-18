This is the Go 1.24.7 standard library unicode package, with its original
BSD notices and tests, renamed to unicode15. Its data version is 15.0.0,
matching the pinned JDK 21. Java-semantics helpers must not silently adopt
new Unicode classifications or case mappings when the Go toolchain changes.

The source is Go code and public Unicode data, not JDK implementation code.
Package/import names and provenance headers are adjusted. One upstream test
uses strings.Map with the local SpecialCase method because the standard
strings.ToLowerSpecial API accepts the standard library's distinct Go type.
Update intentionally when moving the reference JDK and regenerate/recheck
the Java public-API vectors.
