// Generated from the protobuf enum com.google.javascript.jscomp.serialization.NodeProperty
// in the pinned closure-compiler.jar (v20260415) with `javap -c`: constant names,
// ordinals and proto numbers are copied from its static initializer. Updated to
// v20261015 (PRIVATE_IDENTIFIER = 51) from NodeProperty.values() and
// getNumber() of that jar.
//
// Go-only: language-adapter — generated protobuf classes are a leaf package in
// PORTING.md §2; NodeProperty is only used by rhino (Node, PropTranslator) and
// the TypedAST serializer, so it lives here. The protobuf reflection API
// (getValueDescriptor, getDescriptorForType, getDescriptor, internalGetValueMap,
// valueOf(EnumValueDescriptor)) is not ported: nothing in the port uses
// protobuf descriptors.

package rhino

import "github.com/paullepoulpe/jsgo/internal/jutil"

// NodeProperty is the TypedAST NodeProperty enum.
type NodeProperty int32

const (
	NodeProperty_NODE_PROPERTY_UNSPECIFIED NodeProperty = iota
	NodeProperty_IS_PARENTHESIZED
	NodeProperty_SYNTHETIC
	NodeProperty_ADDED_BLOCK
	NodeProperty_IS_CONSTANT_NAME
	NodeProperty_IS_NAMESPACE
	NodeProperty_DIRECT_EVAL
	NodeProperty_FREE_CALL
	NodeProperty_UNUSED_11
	NodeProperty_REFLECTED_OBJECT
	NodeProperty_STATIC_MEMBER
	NodeProperty_GENERATOR_FN
	NodeProperty_ARROW_FN
	NodeProperty_ASYNC_FN
	NodeProperty_YIELD_ALL
	NodeProperty_EXPORT_DEFAULT
	NodeProperty_EXPORT_ALL_FROM
	NodeProperty_CONSTANT_VAR_FLAGS
	NodeProperty_IS_GENERATOR_MARKER
	NodeProperty_IS_GENERATOR_SAFE
	NodeProperty_COMPUTED_PROP_METHOD
	NodeProperty_COMPUTED_PROP_GETTER
	NodeProperty_COMPUTED_PROP_SETTER
	NodeProperty_COMPUTED_PROP_VARIABLE
	NodeProperty_COLOR_FROM_CAST
	NodeProperty_NON_INDEXABLE
	NodeProperty_GOOG_MODULE
	NodeProperty_DELETED
	NodeProperty_MODULE_ALIAS
	NodeProperty_IS_UNUSED_PARAMETER
	NodeProperty_MODULE_EXPORT
	NodeProperty_IS_SHORTHAND_PROPERTY
	NodeProperty_ES6_MODULE
	NodeProperty_START_OF_OPT_CHAIN
	NodeProperty_TRAILING_COMMA
	NodeProperty_IS_INFERRED_CONSTANT
	NodeProperty_IS_DECLARED_CONSTANT
	NodeProperty_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION
	NodeProperty_MUTATES_GLOBAL_STATE
	NodeProperty_MUTATES_THIS
	NodeProperty_MUTATES_ARGUMENTS
	NodeProperty_THROWS
	NodeProperty_CLOSURE_UNAWARE_SHADOW
	NodeProperty_PRIVATE_IDENTIFIER
	NodeProperty_UNRECOGNIZED
)

// The proto field numbers (NodeProperty.X_VALUE).
const (
	NodeProperty_NODE_PROPERTY_UNSPECIFIED_VALUE                int32 = 0
	NodeProperty_IS_PARENTHESIZED_VALUE                         int32 = 1
	NodeProperty_SYNTHETIC_VALUE                                int32 = 3
	NodeProperty_ADDED_BLOCK_VALUE                              int32 = 4
	NodeProperty_IS_CONSTANT_NAME_VALUE                         int32 = 6
	NodeProperty_IS_NAMESPACE_VALUE                             int32 = 7
	NodeProperty_DIRECT_EVAL_VALUE                              int32 = 9
	NodeProperty_FREE_CALL_VALUE                                int32 = 10
	NodeProperty_UNUSED_11_VALUE                                int32 = 11
	NodeProperty_REFLECTED_OBJECT_VALUE                         int32 = 12
	NodeProperty_STATIC_MEMBER_VALUE                            int32 = 13
	NodeProperty_GENERATOR_FN_VALUE                             int32 = 14
	NodeProperty_ARROW_FN_VALUE                                 int32 = 15
	NodeProperty_ASYNC_FN_VALUE                                 int32 = 16
	NodeProperty_YIELD_ALL_VALUE                                int32 = 17
	NodeProperty_EXPORT_DEFAULT_VALUE                           int32 = 18
	NodeProperty_EXPORT_ALL_FROM_VALUE                          int32 = 19
	NodeProperty_CONSTANT_VAR_FLAGS_VALUE                       int32 = 20
	NodeProperty_IS_GENERATOR_MARKER_VALUE                      int32 = 21
	NodeProperty_IS_GENERATOR_SAFE_VALUE                        int32 = 22
	NodeProperty_COMPUTED_PROP_METHOD_VALUE                     int32 = 23
	NodeProperty_COMPUTED_PROP_GETTER_VALUE                     int32 = 24
	NodeProperty_COMPUTED_PROP_SETTER_VALUE                     int32 = 25
	NodeProperty_COMPUTED_PROP_VARIABLE_VALUE                   int32 = 26
	NodeProperty_COLOR_FROM_CAST_VALUE                          int32 = 28
	NodeProperty_NON_INDEXABLE_VALUE                            int32 = 29
	NodeProperty_GOOG_MODULE_VALUE                              int32 = 30
	NodeProperty_DELETED_VALUE                                  int32 = 35
	NodeProperty_MODULE_ALIAS_VALUE                             int32 = 36
	NodeProperty_IS_UNUSED_PARAMETER_VALUE                      int32 = 37
	NodeProperty_MODULE_EXPORT_VALUE                            int32 = 38
	NodeProperty_IS_SHORTHAND_PROPERTY_VALUE                    int32 = 39
	NodeProperty_ES6_MODULE_VALUE                               int32 = 40
	NodeProperty_START_OF_OPT_CHAIN_VALUE                       int32 = 41
	NodeProperty_TRAILING_COMMA_VALUE                           int32 = 42
	NodeProperty_IS_INFERRED_CONSTANT_VALUE                     int32 = 43
	NodeProperty_IS_DECLARED_CONSTANT_VALUE                     int32 = 44
	NodeProperty_SYNTHESIZED_UNFULFILLED_NAME_DECLARATION_VALUE int32 = 45
	NodeProperty_MUTATES_GLOBAL_STATE_VALUE                     int32 = 46
	NodeProperty_MUTATES_THIS_VALUE                             int32 = 47
	NodeProperty_MUTATES_ARGUMENTS_VALUE                        int32 = 48
	NodeProperty_THROWS_VALUE                                   int32 = 49
	NodeProperty_CLOSURE_UNAWARE_SHADOW_VALUE                   int32 = 50
	NodeProperty_PRIVATE_IDENTIFIER_VALUE                       int32 = 51
)

var nodeProperty_names = [...]string{
	"NODE_PROPERTY_UNSPECIFIED",
	"IS_PARENTHESIZED",
	"SYNTHETIC",
	"ADDED_BLOCK",
	"IS_CONSTANT_NAME",
	"IS_NAMESPACE",
	"DIRECT_EVAL",
	"FREE_CALL",
	"UNUSED_11",
	"REFLECTED_OBJECT",
	"STATIC_MEMBER",
	"GENERATOR_FN",
	"ARROW_FN",
	"ASYNC_FN",
	"YIELD_ALL",
	"EXPORT_DEFAULT",
	"EXPORT_ALL_FROM",
	"CONSTANT_VAR_FLAGS",
	"IS_GENERATOR_MARKER",
	"IS_GENERATOR_SAFE",
	"COMPUTED_PROP_METHOD",
	"COMPUTED_PROP_GETTER",
	"COMPUTED_PROP_SETTER",
	"COMPUTED_PROP_VARIABLE",
	"COLOR_FROM_CAST",
	"NON_INDEXABLE",
	"GOOG_MODULE",
	"DELETED",
	"MODULE_ALIAS",
	"IS_UNUSED_PARAMETER",
	"MODULE_EXPORT",
	"IS_SHORTHAND_PROPERTY",
	"ES6_MODULE",
	"START_OF_OPT_CHAIN",
	"TRAILING_COMMA",
	"IS_INFERRED_CONSTANT",
	"IS_DECLARED_CONSTANT",
	"SYNTHESIZED_UNFULFILLED_NAME_DECLARATION",
	"MUTATES_GLOBAL_STATE",
	"MUTATES_THIS",
	"MUTATES_ARGUMENTS",
	"THROWS",
	"CLOSURE_UNAWARE_SHADOW",
	"PRIVATE_IDENTIFIER",
	"UNRECOGNIZED",
}

var nodeProperty_numbers = [...]int32{
	0,  // NODE_PROPERTY_UNSPECIFIED
	1,  // IS_PARENTHESIZED
	3,  // SYNTHETIC
	4,  // ADDED_BLOCK
	6,  // IS_CONSTANT_NAME
	7,  // IS_NAMESPACE
	9,  // DIRECT_EVAL
	10, // FREE_CALL
	11, // UNUSED_11
	12, // REFLECTED_OBJECT
	13, // STATIC_MEMBER
	14, // GENERATOR_FN
	15, // ARROW_FN
	16, // ASYNC_FN
	17, // YIELD_ALL
	18, // EXPORT_DEFAULT
	19, // EXPORT_ALL_FROM
	20, // CONSTANT_VAR_FLAGS
	21, // IS_GENERATOR_MARKER
	22, // IS_GENERATOR_SAFE
	23, // COMPUTED_PROP_METHOD
	24, // COMPUTED_PROP_GETTER
	25, // COMPUTED_PROP_SETTER
	26, // COMPUTED_PROP_VARIABLE
	28, // COLOR_FROM_CAST
	29, // NON_INDEXABLE
	30, // GOOG_MODULE
	35, // DELETED
	36, // MODULE_ALIAS
	37, // IS_UNUSED_PARAMETER
	38, // MODULE_EXPORT
	39, // IS_SHORTHAND_PROPERTY
	40, // ES6_MODULE
	41, // START_OF_OPT_CHAIN
	42, // TRAILING_COMMA
	43, // IS_INFERRED_CONSTANT
	44, // IS_DECLARED_CONSTANT
	45, // SYNTHESIZED_UNFULFILLED_NAME_DECLARATION
	46, // MUTATES_GLOBAL_STATE
	47, // MUTATES_THIS
	48, // MUTATES_ARGUMENTS
	49, // THROWS
	50, // CLOSURE_UNAWARE_SHADOW
	51, // PRIVATE_IDENTIFIER
	-1, // UNRECOGNIZED
}

// NodeProperty_values is NodeProperty.values().
func NodeProperty_values() []NodeProperty {
	out := make([]NodeProperty, len(nodeProperty_names))
	for i := range out {
		out[i] = NodeProperty(i)
	}
	return out
}

// String is Enum.name().
func (p NodeProperty) String() string { return nodeProperty_names[p] }

// Ordinal is Enum.ordinal().
func (p NodeProperty) Ordinal() int { return int(p) }

// GetNumber returns the proto number; UNRECOGNIZED has none.
//
// Java: NodeProperty.getNumber()
func (p NodeProperty) GetNumber() int32 {
	if p == NodeProperty_UNRECOGNIZED {
		panic(jutil.NewIllegalArgumentExceptionMsg("Can't get the number of an unknown enum value."))
	}
	return nodeProperty_numbers[p]
}

var nodePropertyByNumber = func() [64]*NodeProperty {
	var out [64]*NodeProperty
	var vals [len(nodeProperty_numbers)]NodeProperty
	for i, n := range nodeProperty_numbers {
		if n >= 0 && int(n) < len(out) && NodeProperty(i) != NodeProperty_UNRECOGNIZED {
			vals[i] = NodeProperty(i)
			out[n] = &vals[i]
		}
	}
	return out
}()

// NodeProperty_ForNumber returns the constant with the given proto number, or
// null.
//
// Java: NodeProperty.forNumber(int)
func NodeProperty_ForNumber(value int32) *NodeProperty {
	if uint32(value) < uint32(len(nodePropertyByNumber)) {
		return nodePropertyByNumber[value]
	}
	return nil
}

// NodeProperty_ValueOf is the deprecated NodeProperty.valueOf(int).
//
// Java: NodeProperty.valueOf(int)
func NodeProperty_ValueOf(value int32) *NodeProperty {
	return NodeProperty_ForNumber(value)
}

// NodeProperty_ValueOfName is Enum.valueOf(String).
//
// Java: NodeProperty.valueOf(String)
func NodeProperty_ValueOfName(name string) NodeProperty {
	for i, n := range nodeProperty_names {
		if n == name {
			return NodeProperty(i)
		}
	}
	panic(jutil.NewIllegalArgumentExceptionMsg("No enum constant com.google.javascript.jscomp.serialization.NodeProperty." + name))
}
