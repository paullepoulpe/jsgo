// Package regex is the port of Closure Compiler's
// com.google.javascript.jscomp.regex package: the JavaScript regular
// expression AST (RegExpTree and its node classes), the sparse code-unit set
// CharRanges and the ES5 Canonicalize tables (CaseCanonicalize).
//
// Patterns are parsed as UTF-16 code units, like Java chars, so offsets in
// exception messages and the code units that appendSourceCode escapes are
// exactly Java's. reg_exp_tree_oracle_test.go diffs this package against the
// Java classes (tools/oracle/java/RegExpDump.java).
package regex
