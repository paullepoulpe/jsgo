// Reconstruction of com.google.debugging.sourcemap.proto.Mapping (generated
// by protoc from src/com/google/debugging/sourcemap/proto/mapping.proto) @
// closure-compiler v20261021.
//
// Go-only: language-adapter — the Java class is protobuf-generated and not
// part of the Java sources; it is rebuilt here from mapping.proto and the
// class signatures in the pinned closure-compiler.jar (javap), as a plain
// struct with the proto2 semantics the port uses: has-bits, defaults,
// builders, and value equality. Protobuf runtime features (descriptors,
// serialization, parsers, unknown fields, toString's text format and the
// descriptor-based hashCode) are not reconstructed.
//
//	syntax = "proto2";
//	package sourcemap;
//	message OriginalMapping {
//	  optional string original_file = 1;
//	  optional int32 line_number = 2;
//	  optional int32 column_position = 3;
//	  optional string identifier = 4;
//	  enum Precision { UNKNOWN_PRECISION = 0; EXACT = 1; APPROXIMATE_LINE = 2; }
//	  optional Precision precision = 5 [default = UNKNOWN_PRECISION];
//	}

package sourcemap

// Has-bits of bitField0_.
const (
	mapping_originalFileBit   = 0x00000001
	mapping_lineNumberBit     = 0x00000002
	mapping_columnPositionBit = 0x00000004
	mapping_identifierBit     = 0x00000008
	mapping_precisionBit      = 0x00000010
)

// Mapping_OriginalMapping is the immutable message
// sourcemap.OriginalMapping.
type Mapping_OriginalMapping struct {
	bitField0_      int32
	precision_      Mapping_OriginalMapping_Precision
	originalFile_   string
	lineNumber_     int
	columnPosition_ int
	identifier_     string
}

// Mapping_OriginalMapping_Precision is OriginalMapping.Precision.
type Mapping_OriginalMapping_Precision int32

const (
	Mapping_OriginalMapping_Precision_UNKNOWN_PRECISION Mapping_OriginalMapping_Precision = 0
	Mapping_OriginalMapping_Precision_EXACT             Mapping_OriginalMapping_Precision = 1
	Mapping_OriginalMapping_Precision_APPROXIMATE_LINE  Mapping_OriginalMapping_Precision = 2
)

var mapping_OriginalMapping_Precision_names = [...]string{"UNKNOWN_PRECISION", "EXACT", "APPROXIMATE_LINE"}

// Mapping_OriginalMapping_Precision_values is Precision.values().
func Mapping_OriginalMapping_Precision_values() []Mapping_OriginalMapping_Precision {
	return []Mapping_OriginalMapping_Precision{
		Mapping_OriginalMapping_Precision_UNKNOWN_PRECISION,
		Mapping_OriginalMapping_Precision_EXACT,
		Mapping_OriginalMapping_Precision_APPROXIMATE_LINE,
	}
}

// String is Enum.name().
func (p Mapping_OriginalMapping_Precision) String() string {
	return mapping_OriginalMapping_Precision_names[p]
}

// GetNumber is the enum value's number.
//
// Java: OriginalMapping.Precision.getNumber()
func (p Mapping_OriginalMapping_Precision) GetNumber() int {
	return int(p)
}

// Mapping_OriginalMapping_Precision_ForNumber returns the enum value with
// the given number; ok is false where Java returns null.
//
// Java: OriginalMapping.Precision.forNumber(int)
func Mapping_OriginalMapping_Precision_ForNumber(value int) (Mapping_OriginalMapping_Precision, bool) {
	if uint(value) <= 2 {
		return Mapping_OriginalMapping_Precision(value), true
	}
	return 0, false
}

var mapping_OriginalMapping_DEFAULT_INSTANCE = &Mapping_OriginalMapping{}

// Mapping_OriginalMapping_GetDefaultInstance returns the message with no
// field set.
//
// Java: OriginalMapping.getDefaultInstance()
func Mapping_OriginalMapping_GetDefaultInstance() *Mapping_OriginalMapping {
	return mapping_OriginalMapping_DEFAULT_INSTANCE
}

// Java: OriginalMapping.hasOriginalFile()
func (m *Mapping_OriginalMapping) HasOriginalFile() bool {
	return (m.bitField0_ & mapping_originalFileBit) != 0
}

// Java: OriginalMapping.getOriginalFile()
func (m *Mapping_OriginalMapping) GetOriginalFile() string {
	return m.originalFile_
}

// Java: OriginalMapping.hasLineNumber()
func (m *Mapping_OriginalMapping) HasLineNumber() bool {
	return (m.bitField0_ & mapping_lineNumberBit) != 0
}

// Java: OriginalMapping.getLineNumber()
func (m *Mapping_OriginalMapping) GetLineNumber() int {
	return m.lineNumber_
}

// Java: OriginalMapping.hasColumnPosition()
func (m *Mapping_OriginalMapping) HasColumnPosition() bool {
	return (m.bitField0_ & mapping_columnPositionBit) != 0
}

// Java: OriginalMapping.getColumnPosition()
func (m *Mapping_OriginalMapping) GetColumnPosition() int {
	return m.columnPosition_
}

// SetOriginalFileAndColumnInPlace mutates OriginalFile and ColumnPosition in place.
func (m *Mapping_OriginalMapping) SetOriginalFileAndColumnInPlace(originalFile string, columnPosition int) {
	m.originalFile_ = originalFile
	m.columnPosition_ = columnPosition
	m.bitField0_ |= mapping_originalFileBit | mapping_columnPositionBit
}

// Java: OriginalMapping.hasIdentifier()
func (m *Mapping_OriginalMapping) HasIdentifier() bool {
	return (m.bitField0_ & mapping_identifierBit) != 0
}

// Java: OriginalMapping.getIdentifier()
func (m *Mapping_OriginalMapping) GetIdentifier() string {
	return m.identifier_
}

// Java: OriginalMapping.hasPrecision()
func (m *Mapping_OriginalMapping) HasPrecision() bool {
	return (m.bitField0_ & mapping_precisionBit) != 0
}

// Java: OriginalMapping.getPrecision()
func (m *Mapping_OriginalMapping) GetPrecision() Mapping_OriginalMapping_Precision {
	if uint(m.precision_) <= 2 {
		return m.precision_
	}
	return Mapping_OriginalMapping_Precision_UNKNOWN_PRECISION
}

// Equals is the generated message equality: the same fields set, with equal
// values.
//
// Java: OriginalMapping.equals(Object)
func (m *Mapping_OriginalMapping) Equals(obj any) bool {
	other, ok := obj.(*Mapping_OriginalMapping)
	return ok && other != nil && *m == *other
}

// Mapping_OriginalMapping_NewBuilder returns an empty builder.
//
// Java: OriginalMapping.newBuilder()
func Mapping_OriginalMapping_NewBuilder() *Mapping_OriginalMapping_Builder {
	return &Mapping_OriginalMapping_Builder{}
}

// Mapping_OriginalMapping_NewBuilderFrom returns a builder initialized from
// prototype.
//
// Java: OriginalMapping.newBuilder(OriginalMapping)
func Mapping_OriginalMapping_NewBuilderFrom(prototype *Mapping_OriginalMapping) *Mapping_OriginalMapping_Builder {
	return (&Mapping_OriginalMapping_Builder{}).MergeFrom(prototype)
}

// ToBuilder returns a builder initialized from this message.
//
// Java: OriginalMapping.toBuilder()
func (m *Mapping_OriginalMapping) ToBuilder() *Mapping_OriginalMapping_Builder {
	if m == mapping_OriginalMapping_DEFAULT_INSTANCE {
		return &Mapping_OriginalMapping_Builder{}
	}
	return &Mapping_OriginalMapping_Builder{Mapping_OriginalMapping: *m}
}

// IsInitialized is always true: the message has no required fields.
//
// Java: OriginalMapping.isInitialized()
func (m *Mapping_OriginalMapping) IsInitialized() bool {
	return true
}

// Mapping_OriginalMapping_Builder is OriginalMapping.Builder.
type Mapping_OriginalMapping_Builder struct {
	Mapping_OriginalMapping
}

// Clear resets every field.
//
// Java: OriginalMapping.Builder.clear()
func (b *Mapping_OriginalMapping_Builder) Clear() *Mapping_OriginalMapping_Builder {
	*b = Mapping_OriginalMapping_Builder{}
	return b
}

// Build returns the message (no field is required).
//
// Java: OriginalMapping.Builder.build()
func (b *Mapping_OriginalMapping_Builder) Build() *Mapping_OriginalMapping {
	m := b.Mapping_OriginalMapping
	return &m
}

// MergeFrom sets every field that other has set.
//
// Java: OriginalMapping.Builder.mergeFrom(OriginalMapping)
func (b *Mapping_OriginalMapping_Builder) MergeFrom(other *Mapping_OriginalMapping) *Mapping_OriginalMapping_Builder {
	if other == mapping_OriginalMapping_DEFAULT_INSTANCE {
		return b
	}
	if other.HasOriginalFile() {
		b.SetOriginalFile(other.originalFile_)
	}
	if other.HasLineNumber() {
		b.SetLineNumber(other.lineNumber_)
	}
	if other.HasColumnPosition() {
		b.SetColumnPosition(other.columnPosition_)
	}
	if other.HasIdentifier() {
		b.SetIdentifier(other.identifier_)
	}
	if other.HasPrecision() {
		b.SetPrecision(other.GetPrecision())
	}
	return b
}

// Java: OriginalMapping.Builder.setOriginalFile(String)
func (b *Mapping_OriginalMapping_Builder) SetOriginalFile(value string) *Mapping_OriginalMapping_Builder {
	b.originalFile_ = value
	b.bitField0_ |= mapping_originalFileBit
	return b
}

// Java: OriginalMapping.Builder.clearOriginalFile()
func (b *Mapping_OriginalMapping_Builder) ClearOriginalFile() *Mapping_OriginalMapping_Builder {
	b.originalFile_ = ""
	b.bitField0_ &^= mapping_originalFileBit
	return b
}

// Java: OriginalMapping.Builder.setLineNumber(int)
func (b *Mapping_OriginalMapping_Builder) SetLineNumber(value int) *Mapping_OriginalMapping_Builder {
	b.lineNumber_ = value
	b.bitField0_ |= mapping_lineNumberBit
	return b
}

// Java: OriginalMapping.Builder.setColumnPosition(int)
func (b *Mapping_OriginalMapping_Builder) SetColumnPosition(value int) *Mapping_OriginalMapping_Builder {
	b.columnPosition_ = value
	b.bitField0_ |= mapping_columnPositionBit
	return b
}

// Java: OriginalMapping.Builder.setIdentifier(String)
func (b *Mapping_OriginalMapping_Builder) SetIdentifier(value string) *Mapping_OriginalMapping_Builder {
	b.identifier_ = value
	b.bitField0_ |= mapping_identifierBit
	return b
}

// Java: OriginalMapping.Builder.clearIdentifier()
func (b *Mapping_OriginalMapping_Builder) ClearIdentifier() *Mapping_OriginalMapping_Builder {
	b.identifier_ = ""
	b.bitField0_ &^= mapping_identifierBit
	return b
}

// Java: OriginalMapping.Builder.setPrecision(OriginalMapping.Precision)
func (b *Mapping_OriginalMapping_Builder) SetPrecision(value Mapping_OriginalMapping_Precision) *Mapping_OriginalMapping_Builder {
	if uint(value) > 2 {
		value = Mapping_OriginalMapping_Precision_UNKNOWN_PRECISION
	}
	b.precision_ = value
	b.bitField0_ |= mapping_precisionBit
	return b
}

// Java: OriginalMapping.Builder.clearPrecision()
func (b *Mapping_OriginalMapping_Builder) ClearPrecision() *Mapping_OriginalMapping_Builder {
	b.precision_ = 0
	b.bitField0_ &^= mapping_precisionBit
	return b
}
