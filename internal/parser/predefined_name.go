// Port of com/google/javascript/jscomp/parsing/parser/PredefinedName.java @ closure-compiler v20260415

package parser

// PredefinedName: the set of all non-keyword, non-reserved words used in
// javascript.
//
// Java: final class PredefinedName (static constants only; the private
// constructor has no Go counterpart).
const (
	PredefinedName_AS   = "as"
	PredefinedName_FROM = "from"
	PredefinedName_GET  = "get"
	PredefinedName_OF   = "of"
	PredefinedName_SET  = "set"
)
