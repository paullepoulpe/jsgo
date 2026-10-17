// Package parser is the port of Closure Compiler's
// com.google.javascript.jscomp.parsing.parser package together with its
// .trees and .util subpackages (one Go package because the three Java
// packages import each other; see PORTING.md §2).
//
// File names map Java classes mechanically: parser/Scanner.java →
// scanner.go, parser/trees/BlockTree.java → trees_block_tree.go,
// parser/util/SourceRange.java → util_source_range.go.
//
// Source text is held as UTF-16 code units ([]uint16, like Java chars), so
// every offset, line and column is in UTF-16 units exactly as in Java.
package parser
