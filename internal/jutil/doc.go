// Package jutil is the Java-semantics runtime for the Closure Compiler port.
//
// The port translates Java source 1:1 (see PORTING.md). Wherever Java
// semantics differ from Go's in a way that can change observable behaviour,
// the ported code calls into this package instead of using the Go builtin:
//
//   - Strings: Java strings are UTF-16. Go strings in this port hold WTF-8
//     (UTF-8 that may also encode lone surrogates, which JavaScript string
//     literals can contain). Index-based Java String APIs (length, charAt,
//     substring, indexOf, compareTo, hashCode) go through the String* helpers
//     here so that indices and orderings are in UTF-16 code units.
//   - Collections: Java's LinkedHashMap/LinkedHashSet/ImmutableMap/ImmutableSet
//     iterate in insertion order; Go maps iterate randomly. Ported code never
//     ranges over a Go map where the order is observable.
//   - Exceptions: Java exceptions become Go panics carrying a *Throwable-like
//     value; Try/catch sites use the helpers in exceptions.go.
//   - Numbers: Java's narrowing casts, shift masking and Double.toString
//     formatting are reproduced exactly.
//
// This package must not depend on any other package in the port.
package jutil
