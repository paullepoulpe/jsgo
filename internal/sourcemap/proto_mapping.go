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

// Mapping is the generated outer class of mapping.proto.
type Mapping struct{}

// Mapping_OriginalMappingOrBuilder is the accessor interface implemented by
// OriginalMapping and its Builder.
type Mapping_OriginalMappingOrBuilder interface {
	// The original source file.
	HasOriginalFile() bool
	GetOriginalFile() string
	// The line in the original file.
	HasLineNumber() bool
	GetLineNumber() int
	// The column number on the line.
	HasColumnPosition() bool
	GetColumnPosition() int
	// The original name of the identifier.
	HasIdentifier() bool
	GetIdentifier() string
	// The type of retrieval performed to get this mapping
	HasPrecision() bool
	GetPrecision() Mapping_OriginalMapping_Precision
}

// Field numbers of OriginalMapping.
const (
	Mapping_OriginalMapping_ORIGINAL_FILE_FIELD_NUMBER   = 1
	Mapping_OriginalMapping_LINE_NUMBER_FIELD_NUMBER     = 2
	Mapping_OriginalMapping_COLUMN_POSITION_FIELD_NUMBER = 3
	Mapping_OriginalMapping_IDENTIFIER_FIELD_NUMBER      = 4
	Mapping_OriginalMapping_PRECISION_FIELD_NUMBER       = 5
)

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
	originalFile_   string
	lineNumber_     int
	columnPosition_ int
	identifier_     string
	precision_      int32
}

// Mapping_OriginalMapping_Precision is OriginalMapping.Precision.
type Mapping_OriginalMapping_Precision int32

const (
	Mapping_OriginalMapping_Precision_UNKNOWN_PRECISION Mapping_OriginalMapping_Precision = 0
	Mapping_OriginalMapping_Precision_EXACT             Mapping_OriginalMapping_Precision = 1
	Mapping_OriginalMapping_Precision_APPROXIMATE_LINE  Mapping_OriginalMapping_Precision = 2
)

// The enum values' numbers.
const (
	Mapping_OriginalMapping_Precision_UNKNOWN_PRECISION_VALUE = 0
	Mapping_OriginalMapping_Precision_EXACT_VALUE             = 1
	Mapping_OriginalMapping_Precision_APPROXIMATE_LINE_VALUE  = 2
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
func Mapping_OriginalMapping_Precision_ForNumber(value int) (p Mapping_OriginalMapping_Precision, ok bool) {
	switch value {
	case 0:
		return Mapping_OriginalMapping_Precision_UNKNOWN_PRECISION, true
	case 1:
		return Mapping_OriginalMapping_Precision_EXACT, true
	case 2:
		return Mapping_OriginalMapping_Precision_APPROXIMATE_LINE, true
	default:
		return 0, false
	}
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
	result, ok := Mapping_OriginalMapping_Precision_ForNumber(int(m.precision_))
	if !ok {
		return Mapping_OriginalMapping_Precision_UNKNOWN_PRECISION
	}
	return result
}

// Equals is the generated message equality: the same fields set, with equal
// values.
//
// Java: OriginalMapping.equals(Object)
func (m *Mapping_OriginalMapping) Equals(obj any) bool {
	other, ok := obj.(*Mapping_OriginalMapping)
	if !ok || other == nil {
		return false
	}
	if m == other {
		return true
	}
	if m.HasOriginalFile() != other.HasOriginalFile() {
		return false
	}
	if m.HasOriginalFile() && m.GetOriginalFile() != other.GetOriginalFile() {
		return false
	}
	if m.HasLineNumber() != other.HasLineNumber() {
		return false
	}
	if m.HasLineNumber() && m.GetLineNumber() != other.GetLineNumber() {
		return false
	}
	if m.HasColumnPosition() != other.HasColumnPosition() {
		return false
	}
	if m.HasColumnPosition() && m.GetColumnPosition() != other.GetColumnPosition() {
		return false
	}
	if m.HasIdentifier() != other.HasIdentifier() {
		return false
	}
	if m.HasIdentifier() && m.GetIdentifier() != other.GetIdentifier() {
		return false
	}
	if m.HasPrecision() != other.HasPrecision() {
		return false
	}
	if m.HasPrecision() && m.precision_ != other.precision_ {
		return false
	}
	return true
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
	return Mapping_OriginalMapping_NewBuilder().MergeFrom(prototype)
}

// NewBuilderForType returns an empty builder.
//
// Java: OriginalMapping.newBuilderForType()
func (m *Mapping_OriginalMapping) NewBuilderForType() *Mapping_OriginalMapping_Builder {
	return Mapping_OriginalMapping_NewBuilder()
}

// ToBuilder returns a builder initialized from this message.
//
// Java: OriginalMapping.toBuilder()
func (m *Mapping_OriginalMapping) ToBuilder() *Mapping_OriginalMapping_Builder {
	if m == mapping_OriginalMapping_DEFAULT_INSTANCE {
		return &Mapping_OriginalMapping_Builder{}
	}
	return (&Mapping_OriginalMapping_Builder{}).MergeFrom(m)
}

// GetDefaultInstanceForType returns the default instance.
//
// Java: OriginalMapping.getDefaultInstanceForType()
func (m *Mapping_OriginalMapping) GetDefaultInstanceForType() *Mapping_OriginalMapping {
	return mapping_OriginalMapping_DEFAULT_INSTANCE
}

// IsInitialized is always true: the message has no required fields.
//
// Java: OriginalMapping.isInitialized()
func (m *Mapping_OriginalMapping) IsInitialized() bool {
	return true
}

// Mapping_OriginalMapping_Builder is OriginalMapping.Builder.
type Mapping_OriginalMapping_Builder struct {
	bitField0_      int32
	originalFile_   string
	lineNumber_     int
	columnPosition_ int
	identifier_     string
	precision_      int32
}

// Clear resets every field.
//
// Java: OriginalMapping.Builder.clear()
func (b *Mapping_OriginalMapping_Builder) Clear() *Mapping_OriginalMapping_Builder {
	*b = Mapping_OriginalMapping_Builder{}
	return b
}

// GetDefaultInstanceForType returns the default instance.
//
// Java: OriginalMapping.Builder.getDefaultInstanceForType()
func (b *Mapping_OriginalMapping_Builder) GetDefaultInstanceForType() *Mapping_OriginalMapping {
	return mapping_OriginalMapping_DEFAULT_INSTANCE
}

// Build returns the message (no field is required).
//
// Java: OriginalMapping.Builder.build()
func (b *Mapping_OriginalMapping_Builder) Build() *Mapping_OriginalMapping {
	return b.BuildPartial()
}

// BuildPartial returns the message.
//
// Java: OriginalMapping.Builder.buildPartial()
func (b *Mapping_OriginalMapping_Builder) BuildPartial() *Mapping_OriginalMapping {
	result := &Mapping_OriginalMapping{}
	if b.bitField0_ != 0 {
		b.buildPartial0(result)
	}
	return result
}

// Java: OriginalMapping.Builder.buildPartial0(OriginalMapping)
func (b *Mapping_OriginalMapping_Builder) buildPartial0(result *Mapping_OriginalMapping) {
	fromBitField0_ := b.bitField0_
	var toBitField0_ int32
	if (fromBitField0_ & mapping_originalFileBit) != 0 {
		result.originalFile_ = b.originalFile_
		toBitField0_ |= mapping_originalFileBit
	}
	if (fromBitField0_ & mapping_lineNumberBit) != 0 {
		result.lineNumber_ = b.lineNumber_
		toBitField0_ |= mapping_lineNumberBit
	}
	if (fromBitField0_ & mapping_columnPositionBit) != 0 {
		result.columnPosition_ = b.columnPosition_
		toBitField0_ |= mapping_columnPositionBit
	}
	if (fromBitField0_ & mapping_identifierBit) != 0 {
		result.identifier_ = b.identifier_
		toBitField0_ |= mapping_identifierBit
	}
	if (fromBitField0_ & mapping_precisionBit) != 0 {
		result.precision_ = b.precision_
		toBitField0_ |= mapping_precisionBit
	}
	result.bitField0_ |= toBitField0_
}

// MergeFrom sets every field that other has set.
//
// Java: OriginalMapping.Builder.mergeFrom(OriginalMapping)
func (b *Mapping_OriginalMapping_Builder) MergeFrom(other *Mapping_OriginalMapping) *Mapping_OriginalMapping_Builder {
	if other == mapping_OriginalMapping_DEFAULT_INSTANCE {
		return b
	}
	if other.HasOriginalFile() {
		b.originalFile_ = other.originalFile_
		b.bitField0_ |= mapping_originalFileBit
	}
	if other.HasLineNumber() {
		b.SetLineNumber(other.GetLineNumber())
	}
	if other.HasColumnPosition() {
		b.SetColumnPosition(other.GetColumnPosition())
	}
	if other.HasIdentifier() {
		b.identifier_ = other.identifier_
		b.bitField0_ |= mapping_identifierBit
	}
	if other.HasPrecision() {
		b.SetPrecision(other.GetPrecision())
	}
	return b
}

// IsInitialized is always true: the message has no required fields.
//
// Java: OriginalMapping.Builder.isInitialized()
func (b *Mapping_OriginalMapping_Builder) IsInitialized() bool {
	return true
}

// Java: OriginalMapping.Builder.hasOriginalFile()
func (b *Mapping_OriginalMapping_Builder) HasOriginalFile() bool {
	return (b.bitField0_ & mapping_originalFileBit) != 0
}

// Java: OriginalMapping.Builder.getOriginalFile()
func (b *Mapping_OriginalMapping_Builder) GetOriginalFile() string {
	return b.originalFile_
}

// Java: OriginalMapping.Builder.setOriginalFile(String)
func (b *Mapping_OriginalMapping_Builder) SetOriginalFile(value string) *Mapping_OriginalMapping_Builder {
	b.originalFile_ = value
	b.bitField0_ |= mapping_originalFileBit
	return b
}

// Java: OriginalMapping.Builder.clearOriginalFile()
func (b *Mapping_OriginalMapping_Builder) ClearOriginalFile() *Mapping_OriginalMapping_Builder {
	b.originalFile_ = Mapping_OriginalMapping_GetDefaultInstance().GetOriginalFile()
	b.bitField0_ = b.bitField0_ & ^mapping_originalFileBit
	return b
}

// Java: OriginalMapping.Builder.hasLineNumber()
func (b *Mapping_OriginalMapping_Builder) HasLineNumber() bool {
	return (b.bitField0_ & mapping_lineNumberBit) != 0
}

// Java: OriginalMapping.Builder.getLineNumber()
func (b *Mapping_OriginalMapping_Builder) GetLineNumber() int {
	return b.lineNumber_
}

// Java: OriginalMapping.Builder.setLineNumber(int)
func (b *Mapping_OriginalMapping_Builder) SetLineNumber(value int) *Mapping_OriginalMapping_Builder {
	b.lineNumber_ = value
	b.bitField0_ |= mapping_lineNumberBit
	return b
}

// Java: OriginalMapping.Builder.clearLineNumber()
func (b *Mapping_OriginalMapping_Builder) ClearLineNumber() *Mapping_OriginalMapping_Builder {
	b.bitField0_ = b.bitField0_ & ^mapping_lineNumberBit
	b.lineNumber_ = 0
	return b
}

// Java: OriginalMapping.Builder.hasColumnPosition()
func (b *Mapping_OriginalMapping_Builder) HasColumnPosition() bool {
	return (b.bitField0_ & mapping_columnPositionBit) != 0
}

// Java: OriginalMapping.Builder.getColumnPosition()
func (b *Mapping_OriginalMapping_Builder) GetColumnPosition() int {
	return b.columnPosition_
}

// Java: OriginalMapping.Builder.setColumnPosition(int)
func (b *Mapping_OriginalMapping_Builder) SetColumnPosition(value int) *Mapping_OriginalMapping_Builder {
	b.columnPosition_ = value
	b.bitField0_ |= mapping_columnPositionBit
	return b
}

// Java: OriginalMapping.Builder.clearColumnPosition()
func (b *Mapping_OriginalMapping_Builder) ClearColumnPosition() *Mapping_OriginalMapping_Builder {
	b.bitField0_ = b.bitField0_ & ^mapping_columnPositionBit
	b.columnPosition_ = 0
	return b
}

// Java: OriginalMapping.Builder.hasIdentifier()
func (b *Mapping_OriginalMapping_Builder) HasIdentifier() bool {
	return (b.bitField0_ & mapping_identifierBit) != 0
}

// Java: OriginalMapping.Builder.getIdentifier()
func (b *Mapping_OriginalMapping_Builder) GetIdentifier() string {
	return b.identifier_
}

// Java: OriginalMapping.Builder.setIdentifier(String)
func (b *Mapping_OriginalMapping_Builder) SetIdentifier(value string) *Mapping_OriginalMapping_Builder {
	b.identifier_ = value
	b.bitField0_ |= mapping_identifierBit
	return b
}

// Java: OriginalMapping.Builder.clearIdentifier()
func (b *Mapping_OriginalMapping_Builder) ClearIdentifier() *Mapping_OriginalMapping_Builder {
	b.identifier_ = Mapping_OriginalMapping_GetDefaultInstance().GetIdentifier()
	b.bitField0_ = b.bitField0_ & ^mapping_identifierBit
	return b
}

// Java: OriginalMapping.Builder.hasPrecision()
func (b *Mapping_OriginalMapping_Builder) HasPrecision() bool {
	return (b.bitField0_ & mapping_precisionBit) != 0
}

// Java: OriginalMapping.Builder.getPrecision()
func (b *Mapping_OriginalMapping_Builder) GetPrecision() Mapping_OriginalMapping_Precision {
	result, ok := Mapping_OriginalMapping_Precision_ForNumber(int(b.precision_))
	if !ok {
		return Mapping_OriginalMapping_Precision_UNKNOWN_PRECISION
	}
	return result
}

// Java: OriginalMapping.Builder.setPrecision(OriginalMapping.Precision)
func (b *Mapping_OriginalMapping_Builder) SetPrecision(value Mapping_OriginalMapping_Precision) *Mapping_OriginalMapping_Builder {
	b.bitField0_ |= mapping_precisionBit
	b.precision_ = int32(value.GetNumber())
	return b
}

// Java: OriginalMapping.Builder.clearPrecision()
func (b *Mapping_OriginalMapping_Builder) ClearPrecision() *Mapping_OriginalMapping_Builder {
	b.bitField0_ = b.bitField0_ & ^mapping_precisionBit
	b.precision_ = 0
	return b
}
