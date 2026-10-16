// Port of com/google/javascript/rhino/jstype/SimpleReference.java @ closure-compiler v20260415

package rhino

import "github.com/paullepoulpe/jsgo/internal/jutil"

// SimpleReference is a simple immutable reference.
//
// Go-only: language-adapter — SimpleReference<T extends StaticTypedSlot> is
// a generic class whose methods make no virtual calls; subclasses (jscomp's
// SymbolTable.Reference, PreprocessorSymbolTable.Reference) embed the
// struct. getSymbol() keeps StaticRef's return type (see StaticTypedScope);
// GetTypedSymbol returns the T.
type SimpleReference[T StaticTypedSlot] struct {
	symbol T
	node   *Node
}

// Java: SimpleReference(T, Node)
func NewSimpleReference[T StaticTypedSlot](symbol T, node *Node) *SimpleReference[T] {
	r := &SimpleReference[T]{}
	r.InitSimpleReference(symbol, node)
	return r
}

// InitSimpleReference is the Go form of the constructor for subclasses.
//
// Java: SimpleReference(T, Node)
func (r *SimpleReference[T]) InitSimpleReference(symbol T, node *Node) {
	r.symbol = symbol
	r.node = node
}

// Java: SimpleReference.getSymbol()
func (r *SimpleReference[T]) GetSymbol() StaticSlot {
	if jutil.IsNil(r.symbol) {
		return nil
	}
	return r.symbol
}

// GetTypedSymbol is getSymbol() with its Java return type T.
//
// Go-only: language-adapter — covariant return of getSymbol().
func (r *SimpleReference[T]) GetTypedSymbol() T {
	return r.symbol
}

// Java: SimpleReference.getNode()
func (r *SimpleReference[T]) GetNode() *Node {
	return r.node
}

// Java: SimpleReference.getSourceFile()
func (r *SimpleReference[T]) GetSourceFile() StaticSourceFile {
	return r.node.GetStaticSourceFile()
}

// Java: SimpleReference.toString()
func (r *SimpleReference[T]) String() string {
	var sourceName *string
	if r.node != nil {
		sourceName = r.node.GetSourceFileName()
	}
	lineNo := -1
	if r.node != nil {
		lineNo = r.node.GetLineno()
	}
	return jutil.ToJavaString(r.node.GetQualifiedName()) + "@" + jutil.ToJavaString(sourceName) + ":" + jutil.ToJavaString(lineNo)
}
