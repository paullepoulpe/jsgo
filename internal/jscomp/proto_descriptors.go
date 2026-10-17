// Runtime support for the descriptors and reflection of the reconstructed
// protobuf messages (proto_*_reflection.go).
//
// Go-only: language-adapter — the reconstructed protobuf messages (see
// proto_runtime.go) have no Java descriptors. The few reflective operations
// Closure uses on them — Message.getAllFields() and the text-format parser
// (TextFormat.merge, proto_text_format.go) — work on these minimal
// equivalents of com.google.protobuf.Descriptors' Descriptor,
// FieldDescriptor, EnumDescriptor and EnumValueDescriptor, which
// tools/oracle/java/ProtoGoGen.java generates from the jar's descriptors for
// the messages that need them, together with the generated builders'
// MessageReflection.MergeTarget methods.

package jscomp

import "github.com/paullepoulpe/jsgo/internal/jutil"

// Descriptors_Descriptor is Descriptors.Descriptor (a message type).
type Descriptors_Descriptor struct {
	name     string
	fullName string
	// fields in declaration order (Descriptor.getFields()).
	fields []*Descriptors_FieldDescriptor
}

// Java: Descriptors.Descriptor.getName()
func (d *Descriptors_Descriptor) GetName() string { return d.name }

// Java: Descriptors.Descriptor.getFullName()
func (d *Descriptors_Descriptor) GetFullName() string { return d.fullName }

// Java: Descriptors.Descriptor.getFields()
func (d *Descriptors_Descriptor) GetFields() []*Descriptors_FieldDescriptor { return d.fields }

// Java: Descriptors.Descriptor.findFieldByName(String)
func (d *Descriptors_Descriptor) FindFieldByName(name string) *Descriptors_FieldDescriptor {
	for _, f := range d.fields {
		if f.name == name {
			return f
		}
	}
	return nil
}

// Descriptors_FieldDescriptor_JavaType is Descriptors.FieldDescriptor.JavaType.
type Descriptors_FieldDescriptor_JavaType int32

const (
	Descriptors_FieldDescriptor_JavaType_INT Descriptors_FieldDescriptor_JavaType = iota
	Descriptors_FieldDescriptor_JavaType_LONG
	Descriptors_FieldDescriptor_JavaType_FLOAT
	Descriptors_FieldDescriptor_JavaType_DOUBLE
	Descriptors_FieldDescriptor_JavaType_BOOLEAN
	Descriptors_FieldDescriptor_JavaType_STRING
	Descriptors_FieldDescriptor_JavaType_BYTE_STRING
	Descriptors_FieldDescriptor_JavaType_ENUM
	Descriptors_FieldDescriptor_JavaType_MESSAGE
)

// Descriptors_FieldDescriptor_Type is Descriptors.FieldDescriptor.Type.
type Descriptors_FieldDescriptor_Type int32

const (
	Descriptors_FieldDescriptor_Type_DOUBLE Descriptors_FieldDescriptor_Type = iota
	Descriptors_FieldDescriptor_Type_FLOAT
	Descriptors_FieldDescriptor_Type_INT64
	Descriptors_FieldDescriptor_Type_UINT64
	Descriptors_FieldDescriptor_Type_INT32
	Descriptors_FieldDescriptor_Type_FIXED64
	Descriptors_FieldDescriptor_Type_FIXED32
	Descriptors_FieldDescriptor_Type_BOOL
	Descriptors_FieldDescriptor_Type_STRING
	Descriptors_FieldDescriptor_Type_GROUP
	Descriptors_FieldDescriptor_Type_MESSAGE
	Descriptors_FieldDescriptor_Type_BYTES
	Descriptors_FieldDescriptor_Type_UINT32
	Descriptors_FieldDescriptor_Type_ENUM
	Descriptors_FieldDescriptor_Type_SFIXED32
	Descriptors_FieldDescriptor_Type_SFIXED64
	Descriptors_FieldDescriptor_Type_SINT32
	Descriptors_FieldDescriptor_Type_SINT64
)

// Descriptors_FieldDescriptor is Descriptors.FieldDescriptor.
type Descriptors_FieldDescriptor struct {
	name        string
	fullName    string
	number      int32
	type_       Descriptors_FieldDescriptor_Type
	javaType    Descriptors_FieldDescriptor_JavaType
	repeated    bool
	groupLike   bool
	enumType    *Descriptors_EnumDescriptor
	messageType *Descriptors_Descriptor
}

// Java: Descriptors.FieldDescriptor.getName()
func (f *Descriptors_FieldDescriptor) GetName() string { return f.name }

// Java: Descriptors.FieldDescriptor.getFullName()
func (f *Descriptors_FieldDescriptor) GetFullName() string { return f.fullName }

// Java: Descriptors.FieldDescriptor.getNumber()
func (f *Descriptors_FieldDescriptor) GetNumber() int32 { return f.number }

// Java: Descriptors.FieldDescriptor.getType()
func (f *Descriptors_FieldDescriptor) GetType() Descriptors_FieldDescriptor_Type { return f.type_ }

// Java: Descriptors.FieldDescriptor.getJavaType()
func (f *Descriptors_FieldDescriptor) GetJavaType() Descriptors_FieldDescriptor_JavaType {
	return f.javaType
}

// Java: Descriptors.FieldDescriptor.isRepeated()
func (f *Descriptors_FieldDescriptor) IsRepeated() bool { return f.repeated }

// Java: Descriptors.FieldDescriptor.isGroupLike()
func (f *Descriptors_FieldDescriptor) IsGroupLike() bool { return f.groupLike }

// Java: Descriptors.FieldDescriptor.getEnumType()
func (f *Descriptors_FieldDescriptor) GetEnumType() *Descriptors_EnumDescriptor { return f.enumType }

// Java: Descriptors.FieldDescriptor.getMessageType()
func (f *Descriptors_FieldDescriptor) GetMessageType() *Descriptors_Descriptor { return f.messageType }

// Java: Descriptors.FieldDescriptor.toString() (its full name)
func (f *Descriptors_FieldDescriptor) String() string { return f.fullName }

// Descriptors_EnumDescriptor is Descriptors.EnumDescriptor.
type Descriptors_EnumDescriptor struct {
	fullName string
	closed   bool
	values   []*Descriptors_EnumValueDescriptor
}

// Java: Descriptors.EnumDescriptor.getFullName()
func (e *Descriptors_EnumDescriptor) GetFullName() string { return e.fullName }

// Java: Descriptors.EnumDescriptor.isClosed()
func (e *Descriptors_EnumDescriptor) IsClosed() bool { return e.closed }

// Java: Descriptors.EnumDescriptor.findValueByName(String)
func (e *Descriptors_EnumDescriptor) FindValueByName(name string) *Descriptors_EnumValueDescriptor {
	for _, v := range e.values {
		if v.name == name {
			return v
		}
	}
	return nil
}

// Java: Descriptors.EnumDescriptor.findValueByNumber(int)
func (e *Descriptors_EnumDescriptor) FindValueByNumber(number int32) *Descriptors_EnumValueDescriptor {
	// The first value declared with the number (aliases resolve to it).
	for _, v := range e.values {
		if v.number == number {
			return v
		}
	}
	return nil
}

// Descriptors_EnumValueDescriptor is Descriptors.EnumValueDescriptor.
type Descriptors_EnumValueDescriptor struct {
	name   string
	number int32
}

// Java: Descriptors.EnumValueDescriptor.getName()
func (v *Descriptors_EnumValueDescriptor) GetName() string { return v.name }

// Java: Descriptors.EnumValueDescriptor.getNumber()
func (v *Descriptors_EnumValueDescriptor) GetNumber() int32 { return v.number }

// protoAllFields is the body of a generated Message.getAllFields(): the
// fields that are set, by field number, each with its value (a list for a
// repeated field). Enum values are the Go enum constants, not
// EnumValueDescriptors.
func protoAllFields(byNumber []*Descriptors_FieldDescriptor, value func(f *Descriptors_FieldDescriptor) (any, bool)) *jutil.LinkedHashMap[*Descriptors_FieldDescriptor, any] {
	result := jutil.NewLinkedHashMap[*Descriptors_FieldDescriptor, any]()
	for _, f := range byNumber {
		if v, ok := value(f); ok {
			result.Put(f, v)
		}
	}
	return result
}

// TextFormat_MergeTarget is MessageReflection.MergeTarget: what the
// text-format parser needs of a builder. The generated builders implement it
// (proto_*_reflection.go).
type TextFormat_MergeTarget interface {
	// GetDescriptorForType is MergeTarget.getDescriptorForType().
	GetDescriptorForType() *Descriptors_Descriptor
	// protoHasField is MergeTarget.hasField(FieldDescriptor) (singular fields).
	protoHasField(f *Descriptors_FieldDescriptor) bool
	// protoSetField is MergeTarget.setField(FieldDescriptor, Object).
	protoSetField(f *Descriptors_FieldDescriptor, value any)
	// protoAddRepeatedField is MergeTarget.addRepeatedField(FieldDescriptor, Object).
	protoAddRepeatedField(f *Descriptors_FieldDescriptor, value any)
	// protoNewMergeTargetForField is
	// MergeTarget.newMergeTargetForField(FieldDescriptor, Message) for a
	// message field (no extensions here).
	protoNewMergeTargetForField(f *Descriptors_FieldDescriptor) TextFormat_MergeTarget
	// protoFinish is MergeTarget.finish(): the built (partial) message.
	protoFinish() any
}
