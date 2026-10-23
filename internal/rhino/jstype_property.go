// Port of com/google/javascript/rhino/jstype/Property.java @ closure-compiler v20261015

/*
 *
 * ***** BEGIN LICENSE BLOCK *****
 * Version: MPL 1.1
 *
 * The contents of this file are subject to the Mozilla Public License Version
 * 1.1 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 * http://www.mozilla.org/MPL/
 *
 * Software distributed under the License is distributed on an "AS IS" basis,
 * WITHOUT WARRANTY OF ANY KIND, either express or implied. See the License
 * for the specific language governing rights and limitations under the
 * License.
 *
 * The Original Code is Rhino code, released
 * May 6, 1999.
 *
 * The Initial Developer of the Original Code is
 * Netscape Communications Corporation.
 * Portions created by the Initial Developer are Copyright (C) 1997-1999
 * the Initial Developer. All Rights Reserved.
 *
 * Contributor(s):
 *   Nick Santos
 *   Google Inc.
 *
 * This Go translation is distributed under the MPL 1.1 only, as the
 * original's dual licence allows; the other licence it offers is not
 * used.
 *
 * ***** END LICENSE BLOCK ***** */

package rhino

import (
	"github.com/paullepoulpe/jsgo/internal/jutil"
)

// propertySymbolID returns the 1-based index of sym in r's list of symbols
// used as property keys, adding it if new.
//
// Go-only: performance — a Property names a symbol by this index (it shares
// the 32-bit name field with string IDs), resolved through the registry of
// the property's type, so that the list dies with the registry.
func (r *JSTypeRegistry) propertySymbolID(sym *KnownSymbolType) rhinoStringID {
	for i, existing := range r.propertySymbols {
		if existing == sym {
			return rhinoStringID(i + 1)
		}
	}
	r.propertySymbols = append(r.propertySymbols, sym)
	return rhinoStringID(len(r.propertySymbols))
}

// symbol is the key of a symbol-keyed property.
func (p *Property) symbol() *KnownSymbolType {
	return p.type_.GetJSTypeBase().registry.propertySymbols[p.nameID-1]
}

// Property is a property slot of an object.
type Property struct {
	// Property's type.
	type_ JSType

	// The node corresponding to this property, e.g., a GETPROP node that
	// declares this property.
	propertyNode *Node

	// The JSDocInfo for this property.
	docInfo *JSDocInfo

	// Property's name: either a 32-bit RhinoStringPool ID (!isSymbol) or a
	// 1-based index into the registry's propertySymbols (isSymbol).
	nameID rhinoStringID

	// Whether the property's type is inferred.
	inferred bool
	isSymbol bool
}

// Property_KeyKind is what kind of property this is - currently the spec
// only supports strings & symbols.
type Property_KeyKind int32

const (
	Property_KeyKind_STRING Property_KeyKind = iota
	Property_KeyKind_SYMBOL
)

var property_KeyKind_names = [...]string{"STRING", "SYMBOL"}

// Property_KeyKind_values is Property.KeyKind.values().
func Property_KeyKind_values() []Property_KeyKind {
	return []Property_KeyKind{Property_KeyKind_STRING, Property_KeyKind_SYMBOL}
}

// String is Enum.name().
func (k Property_KeyKind) String() string { return property_KeyKind_names[k] }

// Property_Key is a valid key for a property. May be either a string or a
// well-known symbol such as Symbol.iterator.
//
// Go-only: language-adapter — Java's sealed interface Property.Key and its
// two records (Property.StringKey, Property.SymbolKey) are one value struct
// tagged by kind, so that keys (created on every string-keyed property
// lookup) are passed without allocation. NewProperty_StringKey and
// NewProperty_SymbolKey are the record constructors; Equals, HashCode and
// String are the records' generated methods.
type Property_Key struct {
	kind    Property_KeyKind
	string_ string
	symbol  *KnownSymbolType
}

// Kind returns the kind of the key.
//
// Java: Property.Key.kind()
func (k Property_Key) Kind() Property_KeyKind {
	return k.kind
}

// String_ is the string of a StringKey; it throws for a SymbolKey.
//
// Java: Property.Key.string()
func (k Property_Key) String_() string {
	if k.kind == Property_KeyKind_SYMBOL {
		panic(jutil.NewUnsupportedOperationException())
	}
	return k.string_
}

// Symbol is the symbol of a SymbolKey; it throws for a StringKey.
//
// Java: Property.Key.symbol()
func (k Property_Key) Symbol() *KnownSymbolType {
	if k.kind == Property_KeyKind_STRING {
		panic(jutil.NewUnsupportedOperationException())
	}
	return k.symbol
}

// HumanReadableName is the string key or the symbol's display name.
//
// Java: Property.Key.humanReadableName()
func (k Property_Key) HumanReadableName() string {
	switch k.kind {
	case Property_KeyKind_STRING:
		return k.String_()
	}
	return *k.Symbol().GetDisplayName()
}

// Java: Property.Key.matches(String) (default method)
func (k Property_Key) Matches(stringKey string) bool {
	return k.Kind() == Property_KeyKind_STRING && k.String_() == stringKey
}

// Equals is the records' equals(Object).
func (k Property_Key) Equals(o any) bool {
	other, ok := o.(Property_Key)
	if !ok || other.kind != k.kind {
		return false
	}
	if k.kind == Property_KeyKind_STRING {
		return k.string_ == other.string_
	}
	return jutil.ObjectsEquals(k.symbol, other.symbol)
}

// HashCode is the records' hashCode().
func (k Property_Key) HashCode() int32 {
	if k.kind == Property_KeyKind_STRING {
		return jutil.StringHashCode(k.string_)
	}
	return jutil.ObjectHashCode(k.symbol)
}

// String is the records' toString().
func (k Property_Key) String() string {
	if k.kind == Property_KeyKind_STRING {
		return "StringKey[string=" + k.string_ + "]"
	}
	return "SymbolKey[symbol=" + jutil.ToJavaString(k.symbol) + "]"
}

// NewProperty_StringKey is a property string key, such as a in {a: 0}.
//
// Java: Property.StringKey(String)
func NewProperty_StringKey(string_ string) Property_Key {
	return Property_Key{kind: Property_KeyKind_STRING, string_: string_}
}

// NewProperty_SymbolKey is a property well-known symbol key, such as
// Symbol.iterator.
//
// Java: Property.SymbolKey(KnownSymbolType)
func NewProperty_SymbolKey(symbol *KnownSymbolType) Property_Key {
	jutil.CheckNotNull(symbol)
	return Property_Key{kind: Property_KeyKind_SYMBOL, symbol: symbol}
}

// Property_OwnedProperty is a property instance associated with particular
// owner type.
type Property_OwnedProperty struct {
	owner ObjectType
	value *Property
}

// Java: Property.OwnedProperty(ObjectType, Property)
func NewProperty_OwnedProperty(owner ObjectType, value *Property) *Property_OwnedProperty {
	return &Property_OwnedProperty{owner: owner, value: value}
}

// Java: Property.OwnedProperty.getOwner()
func (p *Property_OwnedProperty) GetOwner() ObjectType {
	return p.owner
}

// Java: Property.OwnedProperty.getValue()
func (p *Property_OwnedProperty) GetValue() *Property {
	return p.value
}

// Java: Property.OwnedProperty.getOwnerInstanceType()
func (p *Property_OwnedProperty) GetOwnerInstanceType() ObjectType {
	if p.owner.IsFunctionPrototypeType() {
		return p.owner.GetOwnerFunction().GetInstanceType()
	}
	return p.owner
}

// Java: Property.OwnedProperty.isOwnedByInterface()
func (p *Property_OwnedProperty) IsOwnedByInterface() bool {
	if p.owner.IsFunctionPrototypeType() {
		return p.owner.GetOwnerFunction().IsInterface()
	}
	return p.owner.IsInterface()
}

// Java: Property.serialVersionUID
const Property_serialVersionUID int64 = 1

// Java: Property(String, JSType, boolean, Node)
func NewProperty(name string, type_ JSType, inferred bool, propertyNode *Node) *Property {
	jutil.CheckNotNullMsg1(type_, "Null type specified for %s", name)
	p := &Property{}
	p.nameID = rhinoStringPoolIntern(name)
	p.type_ = type_
	p.inferred = inferred
	p.propertyNode = propertyNode
	return p
}

// Java: Property(Key, JSType, boolean, Node)
func NewPropertyByKey(name Property_Key, type_ JSType, inferred bool, propertyNode *Node) *Property {
	if jutil.IsNil(type_) {
		jutil.CheckNotNullMsg1(type_, "Null type specified for %s", name)
	}
	p := &Property{}
	switch name.Kind() {
	case Property_KeyKind_STRING:
		p.nameID = rhinoStringPoolIntern(name.String_())
	case Property_KeyKind_SYMBOL:
		p.nameID = type_.GetJSTypeBase().registry.propertySymbolID(name.Symbol())
		p.isSymbol = true
	}
	p.type_ = type_
	p.inferred = inferred
	p.propertyNode = propertyNode
	return p
}

func (p *Property) nameObject() any {
	if !p.isSymbol {
		return rhinoStringFromID(p.nameID)
	}
	return p.symbol()
}

// Java: Property.getName()
func (p *Property) GetName() string {
	if !p.isSymbol {
		return rhinoStringFromID(p.nameID)
	}
	return *p.symbol().GetDisplayName()
}

// Java: Property.getNode()
func (p *Property) GetNode() *Node {
	return p.propertyNode
}

// Java: Property.getSourceFile()
func (p *Property) GetSourceFile() StaticSourceFile {
	if p.propertyNode == nil {
		return nil
	}
	return p.propertyNode.GetStaticSourceFile()
}

// GetSymbol returns this property (Java's covariant return type is
// Property; see StaticTypedScope).
//
// Java: Property.getSymbol()
func (p *Property) GetSymbol() StaticSlot {
	return p
}

// GetDeclaration returns this property, or null if it has no node (Java's
// covariant return type is Property; see StaticTypedScope).
//
// Java: Property.getDeclaration()
func (p *Property) GetDeclaration() StaticRef {
	if p.propertyNode == nil {
		return nil
	}
	return p
}

// Java: Property.getType()
func (p *Property) GetType() JSType {
	return p.type_
}

// Java: Property.isTypeInferred()
func (p *Property) IsTypeInferred() bool {
	return p.inferred
}

// Java: Property.isFromExterns()
func (p *Property) IsFromExterns() bool {
	if p.propertyNode == nil {
		return false
	}
	return p.propertyNode.IsFromExterns()
}

// Java: Property.setType(JSType)
func (p *Property) SetType(type_ JSType) {
	if jutil.IsNil(type_) {
		jutil.CheckNotNullMsg1(type_, "Null type specified for property %s", jutil.ToJavaString(p.nameObject()))
	}
	p.type_ = type_
}

// Java: Property.getJSDocInfo()
func (p *Property) GetJSDocInfo() *JSDocInfo {
	return p.docInfo
}

// Java: Property.setJSDocInfo(JSDocInfo)
func (p *Property) SetJSDocInfo(info *JSDocInfo) {
	p.docInfo = info
}

// Java: Property.setNode(Node)
func (p *Property) SetNode(n *Node) {
	p.propertyNode = n
}

// Java: Property.toString()
func (p *Property) String() string {
	return "Property { " +
		" name: " + jutil.ToJavaString(p.nameObject()) +
		", type:" + jutil.ToJavaString(p.type_) +
		", inferred: " + jutil.ToJavaString(p.inferred) +
		"}"
}

// Java: Property.getScope()
func (p *Property) GetScope() StaticScope {
	panic(jutil.NewUnsupportedOperationException())
}

// Java: Property.hashCode()
func (p *Property) HashCode() int32 {
	if !p.isSymbol {
		return 31*(31+jutil.StringHashCode(rhinoStringFromID(p.nameID))) + jutil.ObjectHashCode(p.type_)
	}
	return jutil.ObjectsHash(p.symbol(), p.type_)
}

// Equals is Object.equals (Property does not override it): identity.
func (p *Property) Equals(o any) bool {
	return any(p) == o
}
