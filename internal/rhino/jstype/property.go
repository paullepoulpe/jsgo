// Port of com/google/javascript/rhino/jstype/Property.java @ closure-compiler v20260415
//
// A property slot of an object. Implements StaticTypedSlot and StaticTypedRef
// (Java interfaces). In Go we keep Property as a plain struct and expose the
// interface methods directly — callers use *Property as a value rather than
// through interface variables, matching Java's `final class Property`.
//
// Limitations vs. Java:
//   - StaticTypedSlot / StaticTypedRef are not modelled as Go interfaces here;
//     they are not needed by any current tier. The methods are present so call
//     sites that were forwarded from those interfaces still compile.
//   - The sealed Key hierarchy uses a Go interface with kind discrimination
//     rather than Java's sealed interface + record types, which Go lacks.
//   - KnownSymbolType is a minimal stub (name + display name) sufficient for
//     property-key tracking; the full JSType hierarchy for symbol values is a
//     TODO for a later tier.

package jstype

import (
	"fmt"

	"github.com/paullepoulpe/jsgo/internal/rhino"
)

// KeyKind mirrors Property.KeyKind — the two categories of property key.
type KeyKind int

const (
	KeyKindString KeyKind = iota
	KeyKindSymbol
)

// Key mirrors Property.Key — a sealed interface for property keys.
// Implementations: StringKey, SymbolKey.
type Key interface {
	Kind() KeyKind
	// StringVal returns the string value; panics if Kind == KeyKindSymbol.
	StringVal() string
	// SymbolVal returns the KnownSymbolType; panics if Kind == KeyKindString.
	SymbolVal() *KnownSymbolType
	// HumanReadableName mirrors Key.humanReadableName().
	HumanReadableName() string
	// Matches mirrors Key.matches(String).
	Matches(stringKey string) bool
}

// StringKey mirrors Property.StringKey — a string property key.
type StringKey struct {
	str string
}

// NewStringKey constructs a StringKey. Panics if str is empty (mirrors
// Java's checkNotNull on the record component).
func NewStringKey(str string) StringKey {
	if str == "" {
		panic("jstype.NewStringKey: empty string key")
	}
	return StringKey{str: str}
}

func (k StringKey) Kind() KeyKind            { return KeyKindString }
func (k StringKey) StringVal() string        { return k.str }
func (k StringKey) SymbolVal() *KnownSymbolType { panic("jstype.StringKey.SymbolVal: not a symbol key") }
func (k StringKey) HumanReadableName() string { return k.str }
func (k StringKey) Matches(s string) bool    { return k.str == s }

// SymbolKey mirrors Property.SymbolKey — a well-known symbol property key.
type SymbolKey struct {
	sym *KnownSymbolType
}

// NewSymbolKey constructs a SymbolKey. Panics if sym is nil (mirrors
// Java's checkNotNull on the record component).
func NewSymbolKey(sym *KnownSymbolType) SymbolKey {
	if sym == nil {
		panic("jstype.NewSymbolKey: nil KnownSymbolType")
	}
	return SymbolKey{sym: sym}
}

func (k SymbolKey) Kind() KeyKind            { return KeyKindSymbol }
func (k SymbolKey) StringVal() string        { panic("jstype.SymbolKey.StringVal: not a string key") }
func (k SymbolKey) SymbolVal() *KnownSymbolType { return k.sym }
func (k SymbolKey) HumanReadableName() string { return k.sym.GetDisplayName() }
func (k SymbolKey) Matches(string) bool      { return false }

// KnownSymbolType is defined in known_symbol_type.go (full 1:1 port of
// com/google/javascript/rhino/jstype/KnownSymbolType.java). Property uses
// the same type as the property-key identity for symbol-keyed properties.

// OwnedProperty mirrors Property.OwnedProperty — a property paired with
// the ObjectType that directly owns it (as opposed to inheriting it).
type OwnedProperty struct {
	owner *ObjectType
	value *Property
}

// NewOwnedProperty constructs an OwnedProperty.
func NewOwnedProperty(owner *ObjectType, value *Property) *OwnedProperty {
	return &OwnedProperty{owner: owner, value: value}
}

// GetOwner mirrors OwnedProperty.getOwner().
func (op *OwnedProperty) GetOwner() *ObjectType { return op.owner }

// GetValue mirrors OwnedProperty.getValue().
func (op *OwnedProperty) GetValue() *Property { return op.value }

// GetOwnerInstanceType mirrors OwnedProperty.getOwnerInstanceType().
// Since this port does not yet have FunctionPrototypeType, we always
// return the owner directly.
//
// TODO(later-tier): handle owner.isFunctionPrototypeType() once
// FunctionPrototypeType is ported.
func (op *OwnedProperty) GetOwnerInstanceType() *ObjectType {
	return op.owner
}

// GetOwnerInstantiation mirrors the intent of OwnedProperty.getOwnerInstanceType()
// when the caller needs the declaring *instantiation* of the owner for
// template-type resolution (e.g. TypeValidator.checkPropertyType uses the
// owner's TemplateTypeMap to substitute template params in the required type).
// In this port, where TemplateTypeMap resolution for owner instantiations is
// not yet wired, it returns the same value as GetOwnerInstanceType.
//
// TODO(later-tier): when TemplatizedType / FunctionPrototypeType are fully
// ported, resolve the owner through getOwnerFunction().getInstanceType() and
// apply its TemplateTypeMap.
func (op *OwnedProperty) GetOwnerInstantiation() *ObjectType {
	return op.GetOwnerInstanceType()
}

// IsOwnedByInterface mirrors OwnedProperty.isOwnedByInterface().
//
// TODO(later-tier): implement isFunctionPrototypeType() dispatch.
func (op *OwnedProperty) IsOwnedByInterface() bool {
	return op.owner.IsStructural
}

// Property mirrors Property.java — a property slot of an object.
//
// The Java field `name` is stored as an `any` (either `string` or
// `*KnownSymbolType`) to mirror the memory-saving note in the Java source.
// All public methods use Key for type safety.
type Property struct {
	// name is either a string or *KnownSymbolType, matching Java's Object field.
	name any

	// typ is the property's declared or inferred JSType.
	typ JSType

	// inferred is true when the type came from RHS inference rather than
	// an explicit @type annotation.
	inferred bool

	// propertyNode is the AST node that declares this property.
	propertyNode *rhino.Node

	// docInfo is the JSDoc comment attached to this property, if any.
	docInfo *rhino.JSDocInfo
}

// NewPropertyFromString mirrors Property(String, JSType, boolean, Node).
// Panics if name is empty or typ is nil (mirrors Java's checkNotNull).
func NewPropertyFromString(name string, typ JSType, inferred bool, propertyNode *rhino.Node) *Property {
	if name == "" {
		panic("jstype.NewPropertyFromString: empty name")
	}
	if typ == nil {
		panic(fmt.Sprintf("jstype.NewPropertyFromString: nil type for property %q", name))
	}
	return &Property{
		name:         name,
		typ:          typ,
		inferred:     inferred,
		propertyNode: propertyNode,
	}
}

// NewPropertyFromKey mirrors Property(Key, JSType, boolean, Node).
// Panics if key is nil or typ is nil.
func NewPropertyFromKey(key Key, typ JSType, inferred bool, propertyNode *rhino.Node) *Property {
	if key == nil {
		panic("jstype.NewPropertyFromKey: nil key")
	}
	if typ == nil {
		panic(fmt.Sprintf("jstype.NewPropertyFromKey: nil type for property %q", key.HumanReadableName()))
	}
	var rawName any
	switch key.Kind() {
	case KeyKindString:
		rawName = key.StringVal()
	case KeyKindSymbol:
		rawName = key.SymbolVal()
	}
	return &Property{
		name:         rawName,
		typ:          typ,
		inferred:     inferred,
		propertyNode: propertyNode,
	}
}

// GetName mirrors Property.getName(). Returns the human-readable name:
// the string value for string keys, or the symbol's display name for
// symbol keys.
func (p *Property) GetName() string {
	switch v := p.name.(type) {
	case string:
		return v
	case *KnownSymbolType:
		return v.GetDisplayName()
	default:
		panic(fmt.Sprintf("jstype.Property.GetName: unexpected name type %T", p.name))
	}
}

// GetNode mirrors Property.getNode() / StaticRef.getNode().
func (p *Property) GetNode() *rhino.Node { return p.propertyNode }

// GetSourceFile mirrors Property.getSourceFile() / StaticRef.getSourceFile().
// Returns nil when there is no property node.
func (p *Property) GetSourceFile() string {
	if p.propertyNode == nil {
		return ""
	}
	// TODO(later-tier): expose StaticSourceFile; return filename for now.
	return ""
}

// GetSymbol mirrors StaticTypedRef.getSymbol() — returns the Property itself
// as the canonical slot (Property implements both StaticTypedSlot and
// StaticTypedRef in Java).
func (p *Property) GetSymbol() *Property { return p }

// GetDeclaration mirrors StaticTypedSlot.getDeclaration().
// Returns nil when there is no declaring node.
func (p *Property) GetDeclaration() *Property {
	if p.propertyNode == nil {
		return nil
	}
	return p
}

// GetType mirrors Property.getType() / StaticTypedSlot.getType().
func (p *Property) GetType() JSType { return p.typ }

// IsTypeInferred mirrors Property.isTypeInferred() / StaticTypedSlot.isTypeInferred().
func (p *Property) IsTypeInferred() bool { return p.inferred }

// IsFromExterns mirrors Property.isFromExterns().
func (p *Property) IsFromExterns() bool {
	// TODO(later-tier): expose Node.isFromExterns().
	return false
}

// IsSystemProp reports whether this property was synthesized by the type
// system rather than declared by user code. A property is a "system"
// property when it has no declaring AST node AND its type was inferred
// (i.e. it was added via defineSynthesizedProperty, not from a user
// annotation). This predicate is used by interface-conflict diagnostics
// to skip properties that carry no source location.
//
// There is no direct Java equivalent; in Java the same condition is
// expressed inline as `prop.getDeclaration() == null && prop.isTypeInferred()`.
func (p *Property) IsSystemProp() bool {
	return p.propertyNode == nil && p.inferred
}

// SetType mirrors Property.setType(JSType). Panics if typ is nil.
func (p *Property) SetType(typ JSType) {
	if typ == nil {
		panic(fmt.Sprintf("jstype.Property.SetType: nil type for property %q", p.GetName()))
	}
	p.typ = typ
}

// GetJSDocInfo mirrors Property.getJSDocInfo() / StaticTypedSlot.getJSDocInfo().
func (p *Property) GetJSDocInfo() *rhino.JSDocInfo { return p.docInfo }

// SetJSDocInfo mirrors Property.setJSDocInfo(JSDocInfo).
func (p *Property) SetJSDocInfo(info *rhino.JSDocInfo) { p.docInfo = info }

// SetNode mirrors Property.setNode(Node).
func (p *Property) SetNode(n *rhino.Node) { p.propertyNode = n }

// String mirrors Property.toString().
func (p *Property) String() string {
	return fmt.Sprintf("Property { name: %v, type:%v, inferred: %v}", p.name, p.typ, p.inferred)
}

// GetScope mirrors StaticTypedSlot.getScope(). Panics — not supported for
// Property (matches Java's UnsupportedOperationException).
func (p *Property) GetScope() {
	panic("jstype.Property.GetScope: unsupported operation")
}

// HashCode mirrors Property.hashCode(). Uses name + type identity.
func (p *Property) HashCode() int {
	// Simple hash combining name string and type pointer.
	h := 17
	if s, ok := p.name.(string); ok {
		for _, c := range s {
			h = h*31 + int(c)
		}
	} else if sym, ok := p.name.(*KnownSymbolType); ok {
		for _, c := range sym.GetDisplayName() {
			h = h*31 + int(c)
		}
	}
	// Mix in type pointer identity, matching Java's Objects.hash(name, type).
	if p.typ != nil {
		// Use fmt.Sprintf to get a stable identity contribution.
		for _, c := range fmt.Sprintf("%p", &p.typ) {
			h = h*31 + int(c)
		}
	}
	return h
}
