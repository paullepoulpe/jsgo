// Clean-room j2gort bridge to transpiled/shims/protobuf.
package j2gort

import (
	"strings"
	"sync"
	"sync/atomic"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/protobuf"
)

var _ = (*sync.Pool)(nil)
var _ = (*atomic.Bool)(nil)

type protoFieldDesc struct {
	name string
	idx  int32
}

func int8SliceToBytes(s []int8) []byte {
	if s == nil {
		return nil
	}
	out := make([]byte, len(s))
	for i, v := range s {
		out[i] = byte(v)
	}
	return out
}

func bytesToInt8Slice(b []byte) []int8 {
	if b == nil {
		return nil
	}
	out := make([]int8, len(b))
	for i, v := range b {
		out[i] = int8(v)
	}
	return out
}

func wrapByteString(bs protobuf.ByteString) Ext_com_google_protobuf_ByteString {
	return Ext_com_google_protobuf_ByteString{raw: bs, valid: true}
}

func (b Ext_com_google_protobuf_ByteString) bs() protobuf.ByteString {
	if !b.valid {
		ThrowNPE()
	}
	return b.raw
}

func unwrapByteString(bs Ext_com_google_protobuf_ByteString) protobuf.ByteString {
	if !bs.valid {
		ThrowNPE()
	}
	return bs.raw
}

func BoxByteString(bs Ext_com_google_protobuf_ByteString) any {
	if !bs.valid {
		return nil
	}
	return bs
}

func CheckNotNullByteString(bs Ext_com_google_protobuf_ByteString) Ext_com_google_protobuf_ByteString {
	if !bs.valid {
		ThrowNPE()
	}
	return bs
}

func ByteStringRefEquals(a, b Ext_com_google_protobuf_ByteString) bool {
	if !a.valid || !b.valid {
		return !a.valid && !b.valid
	}
	la := a.raw.Size()
	lb := b.raw.Size()
	if la != lb {
		return false
	}
	if la == 0 {
		return true
	}
	return &a.raw.BytesView()[0] == &b.raw.BytesView()[0]
}

func ByteStringObjEquals(a, b Ext_com_google_protobuf_ByteString) bool {
	if !a.valid || !b.valid {
		return !a.valid && !b.valid
	}
	return a.raw.Equals(b.raw)
}

func ByteStringObjHash(a Ext_com_google_protobuf_ByteString) int32 {
	if !a.valid {
		return 0
	}
	return a.raw.HashCode()
}

func byteStringEquals(a protobuf.ByteString, other any) bool {
	if other == nil {
		return false
	}
	if obs, ok := other.(Ext_com_google_protobuf_ByteString); ok {
		return obs.valid && a.Equals(obs.raw)
	}
	return false
}

func newByteStringIter(bs protobuf.ByteString) Ext_java_util_Iterator {
	raw := bs.BytesView()
	items := make([]any, len(raw))
	for i, v := range raw {
		items[i] = ValueOfByte(int8(v))
	}
	return newSliceObjIter(items)
}

func byteStringForEach(bs protobuf.ByteString, consumer Ext_java_util_function_Consumer) {
	CheckNotNull(consumer)
	for _, v := range bs.BytesView() {
		consumer.Accept(ValueOfByte(int8(v)))
	}
}

func wrapCodedInputStream(c *protobuf.CodedInputStream) Ext_com_google_protobuf_CodedInputStream {
	if c == nil {
		return nil
	}
	w := &Ext_com_google_protobuf_CodedInputStreamBase{Core_: c}
	w.Self_ = w
	return w
}

func (c *Ext_com_google_protobuf_CodedInputStreamBase) cis() *protobuf.CodedInputStream {
	CheckNotNullPtr(c)
	return c.Core_.(*protobuf.CodedInputStream)
}

func unwrapCodedInputStream(c Ext_com_google_protobuf_CodedInputStream) *protobuf.CodedInputStream {
	CheckNotNull(c)
	return c.GetExt_com_google_protobuf_CodedInputStreamBase().cis()
}

func wrapInvalidProtoEx(msg string) Ext_com_google_protobuf_InvalidProtocolBufferException {
	ex := &Ext_com_google_protobuf_InvalidProtocolBufferExceptionBase{}
	ex.Self_ = ex
	InitUserThrowable(ex, "com.google.protobuf.InvalidProtocolBufferException", msg)
	return ex
}

func checkProtoErr(err error) {
	if err == nil {
		return
	}
	if ipbe, ok := err.(*protobuf.InvalidProtocolBufferException); ok {
		panic(wrapInvalidProtoEx(ipbe.GetMessage()))
	}
	if th, ok := err.(javalang.Throwable); ok {
		panic(th)
	}
	panic(javaio.NewIOException(err.Error()))
}

func runProtoOp[R any](fn func() R) (res R) {
	defer func() {
		if rec := recover(); rec != nil {
			if ipbe, ok := rec.(*protobuf.InvalidProtocolBufferException); ok {
				panic(wrapInvalidProtoEx(ipbe.GetMessage()))
			}
			if pe, ok := rec.(*protobuf.TextFormat_ParseException); ok {
				panic(javaio.NewIOException(pe.Error()))
			}
			panic(rec)
		}
	}()
	return fn()
}

func lookupExtEnumByOrdinal[T any](className string, ord int32) *T {
	vals := LookupEnumConstantsByName(className)
	if int(ord) >= 0 && int(ord) < len(vals) {
		return vals[ord].(*T)
	}
	if len(vals) > 0 {
		return vals[len(vals)-1].(*T)
	}
	return nil
}

// Immutable random-access List wrapper for protobuf repeated fields.
type jProtoList struct {
	Ext_java_util_AbstractListBase
	items []any
}

func (l *jProtoList) Size() int32       { return int32(len(l.items)) }
func (l *jProtoList) IsEmpty() bool     { return len(l.items) == 0 }
func (l *jProtoList) Get(i int32) any   { CheckElementIndex(i, int32(len(l.items))); return l.items[i] }
func (l *jProtoList) Iterator() any     { return newSliceObjIter(l.items) }
func (l *jProtoList) IterSlice_() []any { return l.items }
func (l *jProtoList) ToArray_0() []any {
	if len(l.items) == 0 {
		return []any{}
	}
	return CloneSlice(l.items)
}

type jProtoStringList struct {
	jProtoList
}

func (*jProtoStringList) Implements_j2gort_Ext_com_google_protobuf_ProtocolStringList() {}

type jProtoByteStringList struct {
	Ext_java_util_AbstractListBase
	raw []protobuf.ByteString
}

func (l *jProtoByteStringList) Size() int32   { return int32(len(l.raw)) }
func (l *jProtoByteStringList) IsEmpty() bool { return len(l.raw) == 0 }
func (l *jProtoByteStringList) Get(i int32) any {
	CheckElementIndex(i, int32(len(l.raw)))
	return Ext_com_google_protobuf_ByteString{raw: l.raw[i], valid: true}
}
func (l *jProtoByteStringList) toAnySlice() []any {
	if len(l.raw) == 0 {
		return []any{}
	}
	out := make([]any, len(l.raw))
	for i, bs := range l.raw {
		out[i] = Ext_com_google_protobuf_ByteString{raw: bs, valid: true}
	}
	return out
}
func (l *jProtoByteStringList) Iterator() any     { return newSliceObjIter(l.toAnySlice()) }
func (l *jProtoByteStringList) IterSlice_() []any { return l.toAnySlice() }
func (l *jProtoByteStringList) ToArray_0() []any  { return l.toAnySlice() }

func ListGetByteString(l any, i int32) Ext_com_google_protobuf_ByteString {
	if l == nil {
		ThrowNPE()
	}
	if bl, ok := l.(*jProtoByteStringList); ok {
		if bl == nil {
			ThrowNPE()
		}
		CheckElementIndex(i, int32(len(bl.raw)))
		return Ext_com_google_protobuf_ByteString{raw: bl.raw[i], valid: true}
	}
	return ListGetVal[Ext_com_google_protobuf_ByteString](l, i)
}

var emptyProtoList_ *jProtoList
var emptyProtoStringList_ *jProtoStringList
var emptyProtoByteStringList_ *jProtoByteStringList

func init() {
	el := &jProtoList{}
	el.BindExtSelf_(el)
	emptyProtoList_ = el
	es := &jProtoStringList{}
	es.BindExtSelf_(es)
	emptyProtoStringList_ = es
	eb := &jProtoByteStringList{}
	eb.BindExtSelf_(eb)
	emptyProtoByteStringList_ = eb
}

func newProtoByteStringList(raw []protobuf.ByteString) *jProtoByteStringList {
	if len(raw) == 0 {
		return emptyProtoByteStringList_
	}
	l := &jProtoByteStringList{raw: raw}
	l.BindExtSelf_(l)
	return l
}

func newProtoList(items []any) *jProtoList {
	if len(items) == 0 {
		return emptyProtoList_
	}
	l := &jProtoList{items: items}
	l.BindExtSelf_(l)
	return l
}

func newProtoStringList(strs []string) *jProtoStringList {
	if len(strs) == 0 {
		return emptyProtoStringList_
	}
	items := make([]any, len(strs))
	for i, s := range strs {
		items[i] = s
	}
	l := &jProtoStringList{}
	l.items = items
	l.BindExtSelf_(l)
	return l
}

func wrapProtoStringList(strs []string) Ext_com_google_protobuf_ProtocolStringList {
	return newProtoStringList(strs)
}

func wrapProtoInt32List(nums []int32) *jProtoList {
	if len(nums) == 0 {
		return emptyProtoList_
	}
	items := make([]any, len(nums))
	for i, n := range nums {
		items[i] = ValueOfInt(n)
	}
	return newProtoList(items)
}

func unwrapProtoReadCore(v any) any {
	if pb, ok := v.(interface{ protoBuildCore_() any }); ok {
		return pb.protoBuildCore_()
	}
	if gc, ok := v.(interface{ getCore() any }); ok {
		return gc.getCore()
	}
	return v
}

func unwrapProtoMutCore(v any) any {
	if pm, ok := v.(interface{ protoMutCore_() any }); ok {
		return pm.protoMutCore_()
	}
	return v.(interface{ getCore() any }).getCore()
}

func unwrapProtoStringIterable(it Ext_java_lang_Iterable) []string {
	CheckNotNull(it)
	var out []string
	for iter := NewIterCursor(it); iter.HasNext(); {
		v := CheckNotNull(iter.Next())
		out = append(out, v.(string))
	}
	return out
}

func unwrapProtoInt32Iterable(it Ext_java_lang_Iterable) []int32 {
	CheckNotNull(it)
	var out []int32
	for iter := NewIterCursor(it); iter.HasNext(); {
		v := CheckNotNull(iter.Next())
		out = append(out, v.(*BoxInt).Val)
	}
	return out
}

func unwrapProtoMsgIterable[T any](it Ext_java_lang_Iterable) []T {
	CheckNotNull(it)
	var out []T
	for iter := NewIterCursor(it); iter.HasNext(); {
		v := CheckNotNull(iter.Next())
		if pb, ok := v.(interface{ protoBuildCore_() any }); ok {
			out = append(out, pb.protoBuildCore_().(T))
		} else {
			c := v.(interface{ getCore() any }).getCore().(T)
			out = append(out, c)
		}
	}
	return out
}

func protoMsgToByteString(core any) Ext_com_google_protobuf_ByteString {
	if ps, ok := core.(protobuf.ProtoSerializable); ok {
		return wrapByteString(protobuf.MessageToByteString(ps))
	}
	if tba, ok := core.(interface{ ToByteArray() []byte }); ok {
		return wrapByteString(protobuf.ByteString_FromBytesUnsafeWrap(tba.ToByteArray()))
	}
	return wrapByteString(protobuf.ByteString_EMPTY)
}

func protoMsgWriteTo(core any, out Ext_java_io_OutputStream) {
	CheckNotNull(out)
	w := asOutputStream(out)
	if ps, ok := core.(protobuf.ProtoSerializable); ok {
		protobuf.WriteMessageTo(ps, w)
		return
	}
	if tba, ok := core.(interface{ ToByteArray() []byte }); ok {
		if _, err := w.Write(tba.ToByteArray()); err != nil {
			checkProtoErr(err)
		}
	}
}

func protoMsgWriteDelimitedTo(core any, out Ext_java_io_OutputStream) {
	CheckNotNull(out)
	w := asOutputStream(out)
	if ps, ok := core.(protobuf.ProtoSerializable); ok {
		protobuf.WriteMessageDelimitedTo(ps, w)
		return
	}
	if tba, ok := core.(interface{ ToByteArray() []byte }); ok {
		raw := tba.ToByteArray()
		var hdr [10]byte
		n := 0
		v := uint64(len(raw))
		for v >= 0x80 {
			hdr[n] = byte(v) | 0x80
			v >>= 7
			n++
		}
		hdr[n] = byte(v)
		n++
		if _, err := w.Write(hdr[:n]); err != nil {
			checkProtoErr(err)
		}
		if _, err := w.Write(raw); err != nil {
			checkProtoErr(err)
		}
	}
}

func protoMsgEquals(self any, other any) bool {
	if self == other {
		return true
	}
	if other == nil || javalang.IsNil(other) {
		return false
	}
	r1 := unwrapProtoReadCore(self)
	r2 := unwrapProtoReadCore(other)
	if r1 == self || r2 == other {
		return false
	}
	if r1 == r2 {
		return true
	}
	if eq, ok := r1.(interface{ Equals(any) bool }); ok {
		return eq.Equals(r2)
	}
	ps1, s1 := r1.(protobuf.ProtoSerializable)
	ps2, s2 := r2.(protobuf.ProtoSerializable)
	if s1 && s2 {
		return string(protobuf.MessageToByteArray(ps1)) == string(protobuf.MessageToByteArray(ps2))
	}
	return false
}

func protoGetAllFields(core any) Ext_java_util_Map {
	mc := NewJMapCore(false)
	if req, ok := core.(*protobuf.Requirement); ok && req != nil {
		for i, name := range req.PopulatedFieldNames() {
			fd := &Ext_com_google_protobuf_Descriptors_FieldDescriptor{}
			fd.Core_ = protoFieldDesc{name: name, idx: int32(i)}
			mc.Put(fd, Boolean_TRUE)
		}
	}
	return wrapImmutableMap(mc)
}

func readTextFormatInput(in any) string {
	CheckNotNull(in)
	if s, ok := in.(string); ok {
		return s
	}
	if rd, ok := in.(Ext_java_io_Reader); ok {
		return readAllFromReadable(rd)
	}
	var sb strings.Builder
	_ = sb
	return StrVal(in)
}

func (x *protoFieldDesc) IsJavaNil_() bool       { return x == nil }
func (x *jProtoList) IsJavaNil_() bool           { return x == nil }
func (x *jProtoStringList) IsJavaNil_() bool     { return x == nil }
func (x *jProtoByteStringList) IsJavaNil_() bool { return x == nil }
var protoEnumCache_Mapping_OriginalMapping_Precision_ atomic.Pointer[[]*Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_Precision]
func lookupProtoEnum_Mapping_OriginalMapping_Precision(ord int32) *Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_Precision {
	s := protoEnumCache_Mapping_OriginalMapping_Precision_.Load()
	if s == nil {
		raw := LookupEnumConstantsByName("com.google.debugging.sourcemap.proto.Mapping.OriginalMapping.Precision")
		typed := make([]*Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_Precision, len(raw))
		for i, v := range raw { typed[i] = v.(*Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_Precision) }
		protoEnumCache_Mapping_OriginalMapping_Precision_.Store(&typed)
		s = &typed
	}
	vals := *s
	if int(ord) >= 0 && int(ord) < len(vals) { return vals[ord] }
	if len(vals) > 0 { return vals[len(vals)-1] }
	return nil
}
var protoEnumCache_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_ atomic.Pointer[[]*Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior]
func lookupProtoEnum_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior(ord int32) *Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior {
	s := protoEnumCache_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_.Load()
	if s == nil {
		raw := LookupEnumConstantsByName("com.google.javascript.jscomp.ConformanceConfig.LibraryLevelNonAllowlistedConformanceViolationsBehavior")
		typed := make([]*Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior, len(raw))
		for i, v := range raw { typed[i] = v.(*Ext_com_google_javascript_jscomp_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior) }
		protoEnumCache_ConformanceConfig_LibraryLevelNonAllowlistedConformanceViolationsBehavior_.Store(&typed)
		s = &typed
	}
	vals := *s
	if int(ord) >= 0 && int(ord) < len(vals) { return vals[ord] }
	if len(vals) > 0 { return vals[len(vals)-1] }
	return nil
}
var protoEnumCache_FeatureProto_ atomic.Pointer[[]*Ext_com_google_javascript_jscomp_FeatureProto]
func lookupProtoEnum_FeatureProto(ord int32) *Ext_com_google_javascript_jscomp_FeatureProto {
	s := protoEnumCache_FeatureProto_.Load()
	if s == nil {
		raw := LookupEnumConstantsByName("com.google.javascript.jscomp.FeatureProto")
		typed := make([]*Ext_com_google_javascript_jscomp_FeatureProto, len(raw))
		for i, v := range raw { typed[i] = v.(*Ext_com_google_javascript_jscomp_FeatureProto) }
		protoEnumCache_FeatureProto_.Store(&typed)
		s = &typed
	}
	vals := *s
	if int(ord) >= 0 && int(ord) < len(vals) { return vals[ord] }
	if len(vals) > 0 { return vals[len(vals)-1] }
	return nil
}
var protoEnumCache_LifeCycleStageProto_ atomic.Pointer[[]*Ext_com_google_javascript_jscomp_LifeCycleStageProto]
func lookupProtoEnum_LifeCycleStageProto(ord int32) *Ext_com_google_javascript_jscomp_LifeCycleStageProto {
	s := protoEnumCache_LifeCycleStageProto_.Load()
	if s == nil {
		raw := LookupEnumConstantsByName("com.google.javascript.jscomp.LifeCycleStageProto")
		typed := make([]*Ext_com_google_javascript_jscomp_LifeCycleStageProto, len(raw))
		for i, v := range raw { typed[i] = v.(*Ext_com_google_javascript_jscomp_LifeCycleStageProto) }
		protoEnumCache_LifeCycleStageProto_.Store(&typed)
		s = &typed
	}
	vals := *s
	if int(ord) >= 0 && int(ord) < len(vals) { return vals[ord] }
	if len(vals) > 0 { return vals[len(vals)-1] }
	return nil
}
var protoEnumCache_PropertyAccessKindProto_ atomic.Pointer[[]*Ext_com_google_javascript_jscomp_PropertyAccessKindProto]
func lookupProtoEnum_PropertyAccessKindProto(ord int32) *Ext_com_google_javascript_jscomp_PropertyAccessKindProto {
	s := protoEnumCache_PropertyAccessKindProto_.Load()
	if s == nil {
		raw := LookupEnumConstantsByName("com.google.javascript.jscomp.PropertyAccessKindProto")
		typed := make([]*Ext_com_google_javascript_jscomp_PropertyAccessKindProto, len(raw))
		for i, v := range raw { typed[i] = v.(*Ext_com_google_javascript_jscomp_PropertyAccessKindProto) }
		protoEnumCache_PropertyAccessKindProto_.Store(&typed)
		s = &typed
	}
	vals := *s
	if int(ord) >= 0 && int(ord) < len(vals) { return vals[ord] }
	if len(vals) > 0 { return vals[len(vals)-1] }
	return nil
}
var protoEnumCache_Requirement_Severity_ atomic.Pointer[[]*Ext_com_google_javascript_jscomp_Requirement_Severity]
func lookupProtoEnum_Requirement_Severity(ord int32) *Ext_com_google_javascript_jscomp_Requirement_Severity {
	s := protoEnumCache_Requirement_Severity_.Load()
	if s == nil {
		raw := LookupEnumConstantsByName("com.google.javascript.jscomp.Requirement.Severity")
		typed := make([]*Ext_com_google_javascript_jscomp_Requirement_Severity, len(raw))
		for i, v := range raw { typed[i] = v.(*Ext_com_google_javascript_jscomp_Requirement_Severity) }
		protoEnumCache_Requirement_Severity_.Store(&typed)
		s = &typed
	}
	vals := *s
	if int(ord) >= 0 && int(ord) < len(vals) { return vals[ord] }
	if len(vals) > 0 { return vals[len(vals)-1] }
	return nil
}
var protoEnumCache_Requirement_Type_ atomic.Pointer[[]*Ext_com_google_javascript_jscomp_Requirement_Type]
func lookupProtoEnum_Requirement_Type(ord int32) *Ext_com_google_javascript_jscomp_Requirement_Type {
	s := protoEnumCache_Requirement_Type_.Load()
	if s == nil {
		raw := LookupEnumConstantsByName("com.google.javascript.jscomp.Requirement.Type")
		typed := make([]*Ext_com_google_javascript_jscomp_Requirement_Type, len(raw))
		for i, v := range raw { typed[i] = v.(*Ext_com_google_javascript_jscomp_Requirement_Type) }
		protoEnumCache_Requirement_Type_.Store(&typed)
		s = &typed
	}
	vals := *s
	if int(ord) >= 0 && int(ord) < len(vals) { return vals[ord] }
	if len(vals) > 0 { return vals[len(vals)-1] }
	return nil
}
var protoEnumCache_Requirement_TypeMatchingStrategy_ atomic.Pointer[[]*Ext_com_google_javascript_jscomp_Requirement_TypeMatchingStrategy]
func lookupProtoEnum_Requirement_TypeMatchingStrategy(ord int32) *Ext_com_google_javascript_jscomp_Requirement_TypeMatchingStrategy {
	s := protoEnumCache_Requirement_TypeMatchingStrategy_.Load()
	if s == nil {
		raw := LookupEnumConstantsByName("com.google.javascript.jscomp.Requirement.TypeMatchingStrategy")
		typed := make([]*Ext_com_google_javascript_jscomp_Requirement_TypeMatchingStrategy, len(raw))
		for i, v := range raw { typed[i] = v.(*Ext_com_google_javascript_jscomp_Requirement_TypeMatchingStrategy) }
		protoEnumCache_Requirement_TypeMatchingStrategy_.Store(&typed)
		s = &typed
	}
	vals := *s
	if int(ord) >= 0 && int(ord) < len(vals) { return vals[ord] }
	if len(vals) > 0 { return vals[len(vals)-1] }
	return nil
}
var protoEnumCache_InstrumentationPoint_Type_ atomic.Pointer[[]*Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint_Type]
func lookupProtoEnum_InstrumentationPoint_Type(ord int32) *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint_Type {
	s := protoEnumCache_InstrumentationPoint_Type_.Load()
	if s == nil {
		raw := LookupEnumConstantsByName("com.google.javascript.jscomp.instrumentation.reporter.proto.InstrumentationPoint.Type")
		typed := make([]*Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint_Type, len(raw))
		for i, v := range raw { typed[i] = v.(*Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint_Type) }
		protoEnumCache_InstrumentationPoint_Type_.Store(&typed)
		s = &typed
	}
	vals := *s
	if int(ord) >= 0 && int(ord) < len(vals) { return vals[ord] }
	if len(vals) > 0 { return vals[len(vals)-1] }
	return nil
}
var protoEnumCache_InstrumentationPointStats_Presence_ atomic.Pointer[[]*Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_Presence]
func lookupProtoEnum_InstrumentationPointStats_Presence(ord int32) *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_Presence {
	s := protoEnumCache_InstrumentationPointStats_Presence_.Load()
	if s == nil {
		raw := LookupEnumConstantsByName("com.google.javascript.jscomp.instrumentation.reporter.proto.InstrumentationPointStats.Presence")
		typed := make([]*Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_Presence, len(raw))
		for i, v := range raw { typed[i] = v.(*Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_Presence) }
		protoEnumCache_InstrumentationPointStats_Presence_.Store(&typed)
		s = &typed
	}
	vals := *s
	if int(ord) >= 0 && int(ord) < len(vals) { return vals[ord] }
	if len(vals) > 0 { return vals[len(vals)-1] }
	return nil
}
var protoEnumCache_JsdocTag_ atomic.Pointer[[]*Ext_com_google_javascript_jscomp_serialization_JsdocTag]
func lookupProtoEnum_JsdocTag(ord int32) *Ext_com_google_javascript_jscomp_serialization_JsdocTag {
	s := protoEnumCache_JsdocTag_.Load()
	if s == nil {
		raw := LookupEnumConstantsByName("com.google.javascript.jscomp.serialization.JsdocTag")
		typed := make([]*Ext_com_google_javascript_jscomp_serialization_JsdocTag, len(raw))
		for i, v := range raw { typed[i] = v.(*Ext_com_google_javascript_jscomp_serialization_JsdocTag) }
		protoEnumCache_JsdocTag_.Store(&typed)
		s = &typed
	}
	vals := *s
	if int(ord) >= 0 && int(ord) < len(vals) { return vals[ord] }
	if len(vals) > 0 { return vals[len(vals)-1] }
	return nil
}
var protoEnumCache_NodeKind_ atomic.Pointer[[]*Ext_com_google_javascript_jscomp_serialization_NodeKind]
func lookupProtoEnum_NodeKind(ord int32) *Ext_com_google_javascript_jscomp_serialization_NodeKind {
	s := protoEnumCache_NodeKind_.Load()
	if s == nil {
		raw := LookupEnumConstantsByName("com.google.javascript.jscomp.serialization.NodeKind")
		typed := make([]*Ext_com_google_javascript_jscomp_serialization_NodeKind, len(raw))
		for i, v := range raw { typed[i] = v.(*Ext_com_google_javascript_jscomp_serialization_NodeKind) }
		protoEnumCache_NodeKind_.Store(&typed)
		s = &typed
	}
	vals := *s
	if int(ord) >= 0 && int(ord) < len(vals) { return vals[ord] }
	if len(vals) > 0 { return vals[len(vals)-1] }
	return nil
}
var protoEnumCache_NodeProperty_ atomic.Pointer[[]*Ext_com_google_javascript_jscomp_serialization_NodeProperty]
func lookupProtoEnum_NodeProperty(ord int32) *Ext_com_google_javascript_jscomp_serialization_NodeProperty {
	s := protoEnumCache_NodeProperty_.Load()
	if s == nil {
		raw := LookupEnumConstantsByName("com.google.javascript.jscomp.serialization.NodeProperty")
		typed := make([]*Ext_com_google_javascript_jscomp_serialization_NodeProperty, len(raw))
		for i, v := range raw { typed[i] = v.(*Ext_com_google_javascript_jscomp_serialization_NodeProperty) }
		protoEnumCache_NodeProperty_.Store(&typed)
		s = &typed
	}
	vals := *s
	if int(ord) >= 0 && int(ord) < len(vals) { return vals[ord] }
	if len(vals) > 0 { return vals[len(vals)-1] }
	return nil
}
var protoEnumCache_SourceFileProto_LoaderCase_ atomic.Pointer[[]*Ext_com_google_javascript_jscomp_serialization_SourceFileProto_LoaderCase]
func lookupProtoEnum_SourceFileProto_LoaderCase(ord int32) *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_LoaderCase {
	s := protoEnumCache_SourceFileProto_LoaderCase_.Load()
	if s == nil {
		raw := LookupEnumConstantsByName("com.google.javascript.jscomp.serialization.SourceFileProto.LoaderCase")
		typed := make([]*Ext_com_google_javascript_jscomp_serialization_SourceFileProto_LoaderCase, len(raw))
		for i, v := range raw { typed[i] = v.(*Ext_com_google_javascript_jscomp_serialization_SourceFileProto_LoaderCase) }
		protoEnumCache_SourceFileProto_LoaderCase_.Store(&typed)
		s = &typed
	}
	vals := *s
	if int(ord) >= 0 && int(ord) < len(vals) { return vals[ord] }
	if len(vals) > 0 { return vals[len(vals)-1] }
	return nil
}
var protoEnumCache_SourceFileProto_SourceKind_ atomic.Pointer[[]*Ext_com_google_javascript_jscomp_serialization_SourceFileProto_SourceKind]
func lookupProtoEnum_SourceFileProto_SourceKind(ord int32) *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_SourceKind {
	s := protoEnumCache_SourceFileProto_SourceKind_.Load()
	if s == nil {
		raw := LookupEnumConstantsByName("com.google.javascript.jscomp.serialization.SourceFileProto.SourceKind")
		typed := make([]*Ext_com_google_javascript_jscomp_serialization_SourceFileProto_SourceKind, len(raw))
		for i, v := range raw { typed[i] = v.(*Ext_com_google_javascript_jscomp_serialization_SourceFileProto_SourceKind) }
		protoEnumCache_SourceFileProto_SourceKind_.Store(&typed)
		s = &typed
	}
	vals := *s
	if int(ord) >= 0 && int(ord) < len(vals) { return vals[ord] }
	if len(vals) > 0 { return vals[len(vals)-1] }
	return nil
}
var protoEnumCache_TypeProto_KindCase_ atomic.Pointer[[]*Ext_com_google_javascript_jscomp_serialization_TypeProto_KindCase]
func lookupProtoEnum_TypeProto_KindCase(ord int32) *Ext_com_google_javascript_jscomp_serialization_TypeProto_KindCase {
	s := protoEnumCache_TypeProto_KindCase_.Load()
	if s == nil {
		raw := LookupEnumConstantsByName("com.google.javascript.jscomp.serialization.TypeProto.KindCase")
		typed := make([]*Ext_com_google_javascript_jscomp_serialization_TypeProto_KindCase, len(raw))
		for i, v := range raw { typed[i] = v.(*Ext_com_google_javascript_jscomp_serialization_TypeProto_KindCase) }
		protoEnumCache_TypeProto_KindCase_.Store(&typed)
		s = &typed
	}
	vals := *s
	if int(ord) >= 0 && int(ord) < len(vals) { return vals[ord] }
	if len(vals) > 0 { return vals[len(vals)-1] }
	return nil
}
var protoDefault_Mapping_OriginalMapping_ *Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping
var protoDefaultRaw_Mapping_OriginalMapping_ *protobuf.Mapping_OriginalMapping
func (m *Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping) p_Mapping_OriginalMapping() *protobuf.Mapping_OriginalMapping {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_Mapping_OriginalMapping_ }
	return m.Core_.(*protobuf.Mapping_OriginalMapping)
}
func (m *Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping) buildRaw_Mapping_OriginalMapping() *protobuf.Mapping_OriginalMapping { return m.p_Mapping_OriginalMapping() }
func (m *Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping) protoBuildCore_() any { return m.buildRaw_Mapping_OriginalMapping() }
func (b *Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_Builder) p_Mapping_OriginalMapping() *protobuf.Mapping_OriginalMapping {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.Mapping_OriginalMapping)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_Builder) pMut_Mapping_OriginalMapping() *protobuf.Mapping_OriginalMapping {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		b.Raw_ = b.Raw_.ToBuilder()
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_Mapping_OriginalMapping()
}
func (b *Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_Builder) buildRaw_Mapping_OriginalMapping() *protobuf.Mapping_OriginalMapping {
	p := b.p_Mapping_OriginalMapping()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_Builder) buildMsg_Mapping_OriginalMapping() *Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping {
	raw := b.buildRaw_Mapping_OriginalMapping()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_Builder) protoMutCore_() any { return b.pMut_Mapping_OriginalMapping() }
func (b *Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping_Builder) protoBuildCore_() any { return b.buildRaw_Mapping_OriginalMapping() }
var protoDefault_AccessorSummaryEntryProto_ *Ext_com_google_javascript_jscomp_AccessorSummaryEntryProto
var protoDefaultRaw_AccessorSummaryEntryProto_ *protobuf.AccessorSummaryEntryProto
func (m *Ext_com_google_javascript_jscomp_AccessorSummaryEntryProto) p_AccessorSummaryEntryProto() *protobuf.AccessorSummaryEntryProto {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_AccessorSummaryEntryProto_ }
	return m.Core_.(*protobuf.AccessorSummaryEntryProto)
}
func (m *Ext_com_google_javascript_jscomp_AccessorSummaryEntryProto) buildRaw_AccessorSummaryEntryProto() *protobuf.AccessorSummaryEntryProto { return m.p_AccessorSummaryEntryProto() }
func (m *Ext_com_google_javascript_jscomp_AccessorSummaryEntryProto) protoBuildCore_() any { return m.buildRaw_AccessorSummaryEntryProto() }
func (b *Ext_com_google_javascript_jscomp_AccessorSummaryEntryProto_Builder) p_AccessorSummaryEntryProto() *protobuf.AccessorSummaryEntryProto {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.AccessorSummaryEntryProto)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_AccessorSummaryEntryProto_Builder) pMut_AccessorSummaryEntryProto() *protobuf.AccessorSummaryEntryProto {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_AccessorSummaryEntryProto()
}
func (b *Ext_com_google_javascript_jscomp_AccessorSummaryEntryProto_Builder) buildRaw_AccessorSummaryEntryProto() *protobuf.AccessorSummaryEntryProto {
	p := b.p_AccessorSummaryEntryProto()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_AccessorSummaryEntryProto_Builder) buildMsg_AccessorSummaryEntryProto() *Ext_com_google_javascript_jscomp_AccessorSummaryEntryProto {
	raw := b.buildRaw_AccessorSummaryEntryProto()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_AccessorSummaryEntryProto{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_AccessorSummaryEntryProto_Builder) protoMutCore_() any { return b.pMut_AccessorSummaryEntryProto() }
func (b *Ext_com_google_javascript_jscomp_AccessorSummaryEntryProto_Builder) protoBuildCore_() any { return b.buildRaw_AccessorSummaryEntryProto() }
var protoDefault_AccessorSummaryProto_ *Ext_com_google_javascript_jscomp_AccessorSummaryProto
var protoDefaultRaw_AccessorSummaryProto_ *protobuf.AccessorSummaryProto
func (m *Ext_com_google_javascript_jscomp_AccessorSummaryProto) p_AccessorSummaryProto() *protobuf.AccessorSummaryProto {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_AccessorSummaryProto_ }
	return m.Core_.(*protobuf.AccessorSummaryProto)
}
func (m *Ext_com_google_javascript_jscomp_AccessorSummaryProto) buildRaw_AccessorSummaryProto() *protobuf.AccessorSummaryProto { return m.p_AccessorSummaryProto() }
func (m *Ext_com_google_javascript_jscomp_AccessorSummaryProto) protoBuildCore_() any { return m.buildRaw_AccessorSummaryProto() }
func (b *Ext_com_google_javascript_jscomp_AccessorSummaryProto_Builder) p_AccessorSummaryProto() *protobuf.AccessorSummaryProto {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.AccessorSummaryProto)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_AccessorSummaryProto_Builder) pMut_AccessorSummaryProto() *protobuf.AccessorSummaryProto {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_AccessorSummaryProto()
}
func (b *Ext_com_google_javascript_jscomp_AccessorSummaryProto_Builder) buildRaw_AccessorSummaryProto() *protobuf.AccessorSummaryProto {
	p := b.p_AccessorSummaryProto()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_AccessorSummaryProto_Builder) buildMsg_AccessorSummaryProto() *Ext_com_google_javascript_jscomp_AccessorSummaryProto {
	raw := b.buildRaw_AccessorSummaryProto()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_AccessorSummaryProto{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_AccessorSummaryProto_Builder) protoMutCore_() any { return b.pMut_AccessorSummaryProto() }
func (b *Ext_com_google_javascript_jscomp_AccessorSummaryProto_Builder) protoBuildCore_() any { return b.buildRaw_AccessorSummaryProto() }
var protoDefault_ChunkProto_ *Ext_com_google_javascript_jscomp_ChunkProto
var protoDefaultRaw_ChunkProto_ *protobuf.ChunkProto
func (m *Ext_com_google_javascript_jscomp_ChunkProto) p_ChunkProto() *protobuf.ChunkProto {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_ChunkProto_ }
	return m.Core_.(*protobuf.ChunkProto)
}
func (m *Ext_com_google_javascript_jscomp_ChunkProto) buildRaw_ChunkProto() *protobuf.ChunkProto { return m.p_ChunkProto() }
func (m *Ext_com_google_javascript_jscomp_ChunkProto) protoBuildCore_() any { return m.buildRaw_ChunkProto() }
func (b *Ext_com_google_javascript_jscomp_ChunkProto_Builder) p_ChunkProto() *protobuf.ChunkProto {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.ChunkProto)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_ChunkProto_Builder) pMut_ChunkProto() *protobuf.ChunkProto {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_ChunkProto()
}
func (b *Ext_com_google_javascript_jscomp_ChunkProto_Builder) buildRaw_ChunkProto() *protobuf.ChunkProto {
	p := b.p_ChunkProto()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_ChunkProto_Builder) buildMsg_ChunkProto() *Ext_com_google_javascript_jscomp_ChunkProto {
	raw := b.buildRaw_ChunkProto()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_ChunkProto{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_ChunkProto_Builder) protoMutCore_() any { return b.pMut_ChunkProto() }
func (b *Ext_com_google_javascript_jscomp_ChunkProto_Builder) protoBuildCore_() any { return b.buildRaw_ChunkProto() }
var protoDefault_ConformanceConfig_ *Ext_com_google_javascript_jscomp_ConformanceConfig
var protoDefaultRaw_ConformanceConfig_ *protobuf.ConformanceConfig
func (m *Ext_com_google_javascript_jscomp_ConformanceConfig) p_ConformanceConfig() *protobuf.ConformanceConfig {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_ConformanceConfig_ }
	return m.Core_.(*protobuf.ConformanceConfig)
}
func (m *Ext_com_google_javascript_jscomp_ConformanceConfig) buildRaw_ConformanceConfig() *protobuf.ConformanceConfig { return m.p_ConformanceConfig() }
func (m *Ext_com_google_javascript_jscomp_ConformanceConfig) protoBuildCore_() any { return m.buildRaw_ConformanceConfig() }
func (b *Ext_com_google_javascript_jscomp_ConformanceConfig_Builder) p_ConformanceConfig() *protobuf.ConformanceConfig {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.ConformanceConfig)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_ConformanceConfig_Builder) pMut_ConformanceConfig() *protobuf.ConformanceConfig {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_ConformanceConfig()
}
func (b *Ext_com_google_javascript_jscomp_ConformanceConfig_Builder) buildRaw_ConformanceConfig() *protobuf.ConformanceConfig {
	p := b.p_ConformanceConfig()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_ConformanceConfig_Builder) buildMsg_ConformanceConfig() *Ext_com_google_javascript_jscomp_ConformanceConfig {
	raw := b.buildRaw_ConformanceConfig()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_ConformanceConfig{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_ConformanceConfig_Builder) protoMutCore_() any { return b.pMut_ConformanceConfig() }
func (b *Ext_com_google_javascript_jscomp_ConformanceConfig_Builder) protoBuildCore_() any { return b.buildRaw_ConformanceConfig() }
var protoDefault_JSCompilerStateProto_ *Ext_com_google_javascript_jscomp_JSCompilerStateProto
var protoDefaultRaw_JSCompilerStateProto_ *protobuf.JSCompilerStateProto
func (m *Ext_com_google_javascript_jscomp_JSCompilerStateProto) p_JSCompilerStateProto() *protobuf.JSCompilerStateProto {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_JSCompilerStateProto_ }
	return m.Core_.(*protobuf.JSCompilerStateProto)
}
func (m *Ext_com_google_javascript_jscomp_JSCompilerStateProto) buildRaw_JSCompilerStateProto() *protobuf.JSCompilerStateProto { return m.p_JSCompilerStateProto() }
func (m *Ext_com_google_javascript_jscomp_JSCompilerStateProto) protoBuildCore_() any { return m.buildRaw_JSCompilerStateProto() }
func (this_ *Ext_com_google_javascript_jscomp_JSCompilerStateProto) getOrBuildChunksList_() *jProtoList {
	if this_ == nil { ThrowNPE() }
	if p := this_.cache_GetChunksList.Load(); p != nil { return p }
	raw := this_.p_JSCompilerStateProto().GetChunksList()
	if len(raw) == 0 {
		this_.cache_GetChunksList.CompareAndSwap(nil, emptyProtoList_)
		return emptyProtoList_
	}
	wrappers := make([]Ext_com_google_javascript_jscomp_ChunkProto, len(raw))
	items := make([]any, len(raw))
	for i, x := range raw {
		wrappers[i].Core_ = x
		items[i] = &wrappers[i]
	}
	l := newProtoList(items)
	if this_.cache_GetChunksList.CompareAndSwap(nil, l) { return l }
	return this_.cache_GetChunksList.Load()
}
func (b *Ext_com_google_javascript_jscomp_JSCompilerStateProto_Builder) p_JSCompilerStateProto() *protobuf.JSCompilerStateProto {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.JSCompilerStateProto)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_JSCompilerStateProto_Builder) pMut_JSCompilerStateProto() *protobuf.JSCompilerStateProto {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		b.Raw_ = b.Raw_.ToBuilder()
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_JSCompilerStateProto()
}
func (b *Ext_com_google_javascript_jscomp_JSCompilerStateProto_Builder) buildRaw_JSCompilerStateProto() *protobuf.JSCompilerStateProto {
	p := b.p_JSCompilerStateProto()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_JSCompilerStateProto_Builder) buildMsg_JSCompilerStateProto() *Ext_com_google_javascript_jscomp_JSCompilerStateProto {
	raw := b.buildRaw_JSCompilerStateProto()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_JSCompilerStateProto{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_JSCompilerStateProto_Builder) protoMutCore_() any { return b.pMut_JSCompilerStateProto() }
func (b *Ext_com_google_javascript_jscomp_JSCompilerStateProto_Builder) protoBuildCore_() any { return b.buildRaw_JSCompilerStateProto() }
var protoDefault_Requirement_ *Ext_com_google_javascript_jscomp_Requirement
var protoDefaultRaw_Requirement_ *protobuf.Requirement
func (m *Ext_com_google_javascript_jscomp_Requirement) p_Requirement() *protobuf.Requirement {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_Requirement_ }
	return m.Core_.(*protobuf.Requirement)
}
func (m *Ext_com_google_javascript_jscomp_Requirement) buildRaw_Requirement() *protobuf.Requirement { return m.p_Requirement() }
func (m *Ext_com_google_javascript_jscomp_Requirement) protoBuildCore_() any { return m.buildRaw_Requirement() }
func (b *Ext_com_google_javascript_jscomp_Requirement_Builder) p_Requirement() *protobuf.Requirement {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.Requirement)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_Requirement_Builder) pMut_Requirement() *protobuf.Requirement {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		b.Raw_ = b.Raw_.ToBuilder()
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_Requirement()
}
func (b *Ext_com_google_javascript_jscomp_Requirement_Builder) buildRaw_Requirement() *protobuf.Requirement {
	p := b.p_Requirement()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_Requirement_Builder) buildMsg_Requirement() *Ext_com_google_javascript_jscomp_Requirement {
	raw := b.buildRaw_Requirement()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_Requirement{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_Requirement_Builder) protoMutCore_() any { return b.pMut_Requirement() }
func (b *Ext_com_google_javascript_jscomp_Requirement_Builder) protoBuildCore_() any { return b.buildRaw_Requirement() }
var protoDefault_RequirementScopeEntry_ *Ext_com_google_javascript_jscomp_RequirementScopeEntry
var protoDefaultRaw_RequirementScopeEntry_ *protobuf.RequirementScopeEntry
func (m *Ext_com_google_javascript_jscomp_RequirementScopeEntry) p_RequirementScopeEntry() *protobuf.RequirementScopeEntry {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_RequirementScopeEntry_ }
	return m.Core_.(*protobuf.RequirementScopeEntry)
}
func (m *Ext_com_google_javascript_jscomp_RequirementScopeEntry) buildRaw_RequirementScopeEntry() *protobuf.RequirementScopeEntry { return m.p_RequirementScopeEntry() }
func (m *Ext_com_google_javascript_jscomp_RequirementScopeEntry) protoBuildCore_() any { return m.buildRaw_RequirementScopeEntry() }
func (b *Ext_com_google_javascript_jscomp_RequirementScopeEntry_Builder) p_RequirementScopeEntry() *protobuf.RequirementScopeEntry {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.RequirementScopeEntry)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_RequirementScopeEntry_Builder) pMut_RequirementScopeEntry() *protobuf.RequirementScopeEntry {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_RequirementScopeEntry()
}
func (b *Ext_com_google_javascript_jscomp_RequirementScopeEntry_Builder) buildRaw_RequirementScopeEntry() *protobuf.RequirementScopeEntry {
	p := b.p_RequirementScopeEntry()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_RequirementScopeEntry_Builder) buildMsg_RequirementScopeEntry() *Ext_com_google_javascript_jscomp_RequirementScopeEntry {
	raw := b.buildRaw_RequirementScopeEntry()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_RequirementScopeEntry{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_RequirementScopeEntry_Builder) protoMutCore_() any { return b.pMut_RequirementScopeEntry() }
func (b *Ext_com_google_javascript_jscomp_RequirementScopeEntry_Builder) protoBuildCore_() any { return b.buildRaw_RequirementScopeEntry() }
var protoDefault_UniqueIdProto_ *Ext_com_google_javascript_jscomp_UniqueIdProto
var protoDefaultRaw_UniqueIdProto_ *protobuf.UniqueIdProto
func (m *Ext_com_google_javascript_jscomp_UniqueIdProto) p_UniqueIdProto() *protobuf.UniqueIdProto {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_UniqueIdProto_ }
	return m.Core_.(*protobuf.UniqueIdProto)
}
func (m *Ext_com_google_javascript_jscomp_UniqueIdProto) buildRaw_UniqueIdProto() *protobuf.UniqueIdProto { return m.p_UniqueIdProto() }
func (m *Ext_com_google_javascript_jscomp_UniqueIdProto) protoBuildCore_() any { return m.buildRaw_UniqueIdProto() }
func (b *Ext_com_google_javascript_jscomp_UniqueIdProto_Builder) p_UniqueIdProto() *protobuf.UniqueIdProto {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.UniqueIdProto)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_UniqueIdProto_Builder) pMut_UniqueIdProto() *protobuf.UniqueIdProto {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_UniqueIdProto()
}
func (b *Ext_com_google_javascript_jscomp_UniqueIdProto_Builder) buildRaw_UniqueIdProto() *protobuf.UniqueIdProto {
	p := b.p_UniqueIdProto()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_UniqueIdProto_Builder) buildMsg_UniqueIdProto() *Ext_com_google_javascript_jscomp_UniqueIdProto {
	raw := b.buildRaw_UniqueIdProto()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_UniqueIdProto{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_UniqueIdProto_Builder) protoMutCore_() any { return b.pMut_UniqueIdProto() }
func (b *Ext_com_google_javascript_jscomp_UniqueIdProto_Builder) protoBuildCore_() any { return b.buildRaw_UniqueIdProto() }
var protoDefault_VariableMapEntryProto_ *Ext_com_google_javascript_jscomp_VariableMapEntryProto
var protoDefaultRaw_VariableMapEntryProto_ *protobuf.VariableMapEntryProto
func (m *Ext_com_google_javascript_jscomp_VariableMapEntryProto) p_VariableMapEntryProto() *protobuf.VariableMapEntryProto {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_VariableMapEntryProto_ }
	return m.Core_.(*protobuf.VariableMapEntryProto)
}
func (m *Ext_com_google_javascript_jscomp_VariableMapEntryProto) buildRaw_VariableMapEntryProto() *protobuf.VariableMapEntryProto { return m.p_VariableMapEntryProto() }
func (m *Ext_com_google_javascript_jscomp_VariableMapEntryProto) protoBuildCore_() any { return m.buildRaw_VariableMapEntryProto() }
func (b *Ext_com_google_javascript_jscomp_VariableMapEntryProto_Builder) p_VariableMapEntryProto() *protobuf.VariableMapEntryProto {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.VariableMapEntryProto)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_VariableMapEntryProto_Builder) pMut_VariableMapEntryProto() *protobuf.VariableMapEntryProto {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_VariableMapEntryProto()
}
func (b *Ext_com_google_javascript_jscomp_VariableMapEntryProto_Builder) buildRaw_VariableMapEntryProto() *protobuf.VariableMapEntryProto {
	p := b.p_VariableMapEntryProto()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_VariableMapEntryProto_Builder) buildMsg_VariableMapEntryProto() *Ext_com_google_javascript_jscomp_VariableMapEntryProto {
	raw := b.buildRaw_VariableMapEntryProto()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_VariableMapEntryProto{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_VariableMapEntryProto_Builder) protoMutCore_() any { return b.pMut_VariableMapEntryProto() }
func (b *Ext_com_google_javascript_jscomp_VariableMapEntryProto_Builder) protoBuildCore_() any { return b.buildRaw_VariableMapEntryProto() }
var protoDefault_FileProfile_ *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_FileProfile
var protoDefaultRaw_FileProfile_ *protobuf.FileProfile
func (m *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_FileProfile) p_FileProfile() *protobuf.FileProfile {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_FileProfile_ }
	return m.Core_.(*protobuf.FileProfile)
}
func (m *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_FileProfile) buildRaw_FileProfile() *protobuf.FileProfile { return m.p_FileProfile() }
func (m *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_FileProfile) protoBuildCore_() any { return m.buildRaw_FileProfile() }
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_FileProfile_Builder) p_FileProfile() *protobuf.FileProfile {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.FileProfile)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_FileProfile_Builder) pMut_FileProfile() *protobuf.FileProfile {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_FileProfile()
}
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_FileProfile_Builder) buildRaw_FileProfile() *protobuf.FileProfile {
	p := b.p_FileProfile()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_FileProfile_Builder) buildMsg_FileProfile() *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_FileProfile {
	raw := b.buildRaw_FileProfile()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_FileProfile{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_FileProfile_Builder) protoMutCore_() any { return b.pMut_FileProfile() }
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_FileProfile_Builder) protoBuildCore_() any { return b.buildRaw_FileProfile() }
var protoDefault_InstrumentationPoint_ *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint
var protoDefaultRaw_InstrumentationPoint_ *protobuf.InstrumentationPoint
func (m *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint) p_InstrumentationPoint() *protobuf.InstrumentationPoint {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_InstrumentationPoint_ }
	return m.Core_.(*protobuf.InstrumentationPoint)
}
func (m *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint) buildRaw_InstrumentationPoint() *protobuf.InstrumentationPoint { return m.p_InstrumentationPoint() }
func (m *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint) protoBuildCore_() any { return m.buildRaw_InstrumentationPoint() }
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint_Builder) p_InstrumentationPoint() *protobuf.InstrumentationPoint {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.InstrumentationPoint)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint_Builder) pMut_InstrumentationPoint() *protobuf.InstrumentationPoint {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_InstrumentationPoint()
}
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint_Builder) buildRaw_InstrumentationPoint() *protobuf.InstrumentationPoint {
	p := b.p_InstrumentationPoint()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint_Builder) buildMsg_InstrumentationPoint() *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint {
	raw := b.buildRaw_InstrumentationPoint()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint_Builder) protoMutCore_() any { return b.pMut_InstrumentationPoint() }
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint_Builder) protoBuildCore_() any { return b.buildRaw_InstrumentationPoint() }
var protoDefault_InstrumentationPointStats_ *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats
var protoDefaultRaw_InstrumentationPointStats_ *protobuf.InstrumentationPointStats
func (m *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats) p_InstrumentationPointStats() *protobuf.InstrumentationPointStats {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_InstrumentationPointStats_ }
	return m.Core_.(*protobuf.InstrumentationPointStats)
}
func (m *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats) buildRaw_InstrumentationPointStats() *protobuf.InstrumentationPointStats { return m.p_InstrumentationPointStats() }
func (m *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats) protoBuildCore_() any { return m.buildRaw_InstrumentationPointStats() }
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_Builder) p_InstrumentationPointStats() *protobuf.InstrumentationPointStats {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.InstrumentationPointStats)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_Builder) pMut_InstrumentationPointStats() *protobuf.InstrumentationPointStats {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_InstrumentationPointStats()
}
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_Builder) buildRaw_InstrumentationPointStats() *protobuf.InstrumentationPointStats {
	p := b.p_InstrumentationPointStats()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_Builder) buildMsg_InstrumentationPointStats() *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats {
	raw := b.buildRaw_InstrumentationPointStats()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_Builder) protoMutCore_() any { return b.pMut_InstrumentationPointStats() }
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats_Builder) protoBuildCore_() any { return b.buildRaw_InstrumentationPointStats() }
var protoDefault_ReportProfile_ *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile
var protoDefaultRaw_ReportProfile_ *protobuf.ReportProfile
func (m *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile) p_ReportProfile() *protobuf.ReportProfile {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_ReportProfile_ }
	return m.Core_.(*protobuf.ReportProfile)
}
func (m *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile) buildRaw_ReportProfile() *protobuf.ReportProfile { return m.p_ReportProfile() }
func (m *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile) protoBuildCore_() any { return m.buildRaw_ReportProfile() }
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile_Builder) p_ReportProfile() *protobuf.ReportProfile {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.ReportProfile)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile_Builder) pMut_ReportProfile() *protobuf.ReportProfile {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_ReportProfile()
}
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile_Builder) buildRaw_ReportProfile() *protobuf.ReportProfile {
	p := b.p_ReportProfile()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile_Builder) buildMsg_ReportProfile() *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile {
	raw := b.buildRaw_ReportProfile()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile_Builder) protoMutCore_() any { return b.pMut_ReportProfile() }
func (b *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile_Builder) protoBuildCore_() any { return b.buildRaw_ReportProfile() }
var protoDefault_AstNode_ *Ext_com_google_javascript_jscomp_serialization_AstNode
var protoDefaultRaw_AstNode_ *protobuf.AstNode
func (m *Ext_com_google_javascript_jscomp_serialization_AstNode) p_AstNode() *protobuf.AstNode {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_AstNode_ }
	if r, ok := m.Core_.(*protobuf.AstNode); ok { return r }
	return &m.Core_.(*Ext_com_google_javascript_jscomp_serialization_AstNode_Builder).inlineRaw
}
func (m *Ext_com_google_javascript_jscomp_serialization_AstNode) buildRaw_AstNode() *protobuf.AstNode { return m.p_AstNode() }
func (m *Ext_com_google_javascript_jscomp_serialization_AstNode) protoBuildCore_() any { return m.buildRaw_AstNode() }
type astNodeExtSlabCache_ struct {
	wraps       []Ext_com_google_javascript_jscomp_serialization_AstNode
	firstHandle *protobuf.AstNode
	initLen     int
}
func initAstNodeExtSlabCache_(prev any, root *protobuf.AstNode, handles []protobuf.AstNode) any {
	c, _ := prev.(*astNodeExtSlabCache_)
	if c == nil {
		c = &astNodeExtSlabCache_{}
	}
	total := 1 + len(handles)
	if cap(c.wraps) < total {
		c.wraps = make([]Ext_com_google_javascript_jscomp_serialization_AstNode, total, max(total, cap(c.wraps)*2))
		c.firstHandle = nil
		c.initLen = 0
	} else {
		c.wraps = c.wraps[:total]
	}
	if len(handles) > 0 && c.firstHandle != &handles[0] {
		c.firstHandle = &handles[0]
		c.initLen = 0
	}
	if c.initLen == 0 {
		c.wraps[0].Core_ = root
		c.initLen = 1
	}
	for i := c.initLen; i < total; i++ {
		c.wraps[i].Core_ = &handles[i-1]
	}
	if c.initLen < total {
		c.initLen = total
	}
	return c
}
//go:noinline
func getAstNodeChildSlow_(r *protobuf.AstNode, idx int) *Ext_com_google_javascript_jscomp_serialization_AstNode {
	if cIdx, cache, ok := r.SlabChildInit(idx, initAstNodeExtSlabCache_); ok {
		return &cache.(*astNodeExtSlabCache_).wraps[cIdx]
	}
	v := r.GetChild(idx)
	if v == nil {
		return nil
	}
	out := &Ext_com_google_javascript_jscomp_serialization_AstNode{}
	out.Core_ = v
	return out
}
func (this_ *Ext_com_google_javascript_jscomp_serialization_AstNode) Release_() {
	if this_ == nil {
		return
	}
	if r, ok := this_.Core_.(*protobuf.AstNode); ok && r != nil {
		this_.Core_ = nil
		r.ReleaseSlab()
	}
}
func (b *Ext_com_google_javascript_jscomp_serialization_AstNode_Builder) p_AstNode() *protobuf.AstNode {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.AstNode)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_AstNode_Builder) pMut_AstNode() *protobuf.AstNode {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		b.pinned_ = true
		if b.usedBuiltMsg_ && b.builtMsg.Core_ == b { b.builtMsg.Core_ = b.Raw_ }
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_AstNode()
}
func (b *Ext_com_google_javascript_jscomp_serialization_AstNode_Builder) buildRaw_AstNode() *protobuf.AstNode {
	p := b.p_AstNode()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_AstNode_Builder) buildMsg_AstNode() *Ext_com_google_javascript_jscomp_serialization_AstNode {
	raw := b.buildRaw_AstNode()
	if b.fromPool_ && !b.pinned_ && raw == &b.inlineRaw {
		if !b.usedBuiltMsg_ {
			b.usedBuiltMsg_ = true
			b.builtMsg.Core_ = b
			return &b.builtMsg
		}
		if b.builtMsg.Core_ == b {
			return &b.builtMsg
		}
		b.pinned_ = true
	}
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_AstNode{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_AstNode_Builder) protoMutCore_() any { return b.pMut_AstNode() }
func (b *Ext_com_google_javascript_jscomp_serialization_AstNode_Builder) protoBuildCore_() any { return b.buildRaw_AstNode() }
type astNodeBuilderBatch struct {
	idx  atomic.Uint32
	done atomic.Uint32
	ptrs [128]*Ext_com_google_javascript_jscomp_serialization_AstNode_Builder
}
type astNodeBuilderChunk struct {
	items [128]Ext_com_google_javascript_jscomp_serialization_AstNode_Builder
}
var (
	astNodeBuilderBatch_     atomic.Pointer[astNodeBuilderBatch]
	astNodeBuilderPoolMu_    sync.Mutex
	astNodeBuilderCurBatch_  *astNodeBuilderBatch
	astNodeBuilderBatchFree_ []*astNodeBuilderBatch
	astNodeBuilderRetired_   []*astNodeBuilderBatch
	astNodeBuilderFree_      []*Ext_com_google_javascript_jscomp_serialization_AstNode_Builder
)
func allocAstNodeBuilder_() *Ext_com_google_javascript_jscomp_serialization_AstNode_Builder {
	if batch := astNodeBuilderBatch_.Load(); batch != nil {
		if i := batch.idx.Add(1) - 1; i < 128 {
			b := batch.ptrs[i]
			batch.done.Add(1)
			return b
		}
	}
	return allocAstNodeBuilderSlow_()
}
func allocAstNodeBuilderSlow_() *Ext_com_google_javascript_jscomp_serialization_AstNode_Builder {
	astNodeBuilderPoolMu_.Lock()
	if batch := astNodeBuilderBatch_.Load(); batch != nil {
		if i := batch.idx.Add(1) - 1; i < 128 {
			b := batch.ptrs[i]
			batch.done.Add(1)
			astNodeBuilderPoolMu_.Unlock()
			return b
		}
	}
	if old := astNodeBuilderCurBatch_; old != nil {
		astNodeBuilderCurBatch_ = nil
		astNodeBuilderBatch_.Store(nil)
		if old.done.Load() == 128 {
			clear(old.ptrs[:])
			if len(astNodeBuilderBatchFree_) < 16 {
				astNodeBuilderBatchFree_ = append(astNodeBuilderBatchFree_, old)
			}
		} else {
			astNodeBuilderRetired_ = append(astNodeBuilderRetired_, old)
		}
	}
	if len(astNodeBuilderRetired_) > 0 {
		w := 0
		for _, rb := range astNodeBuilderRetired_ {
			if rb.done.Load() == 128 {
				clear(rb.ptrs[:])
				if len(astNodeBuilderBatchFree_) < 16 {
					astNodeBuilderBatchFree_ = append(astNodeBuilderBatchFree_, rb)
				}
			} else {
				astNodeBuilderRetired_[w] = rb
				w++
			}
		}
		clear(astNodeBuilderRetired_[w:])
		astNodeBuilderRetired_ = astNodeBuilderRetired_[:w]
	}
	var nb *astNodeBuilderBatch
	if nbf := len(astNodeBuilderBatchFree_); nbf > 0 {
		nb = astNodeBuilderBatchFree_[nbf-1]
		astNodeBuilderBatchFree_[nbf-1] = nil
		astNodeBuilderBatchFree_ = astNodeBuilderBatchFree_[:nbf-1]
	} else {
		nb = &astNodeBuilderBatch{}
		nb.idx.Store(128)
	}
	nFree := len(astNodeBuilderFree_)
	if nFree >= 128 {
		start := nFree - 128
		copy(nb.ptrs[:], astNodeBuilderFree_[start:])
		clear(astNodeBuilderFree_[start:])
		astNodeBuilderFree_ = astNodeBuilderFree_[:start]
	} else {
		copied := 0
		if nFree > 0 {
			copied = copy(nb.ptrs[:], astNodeBuilderFree_)
			clear(astNodeBuilderFree_)
			astNodeBuilderFree_ = astNodeBuilderFree_[:0]
		}
		nc := &astNodeBuilderChunk{}
		need := 128 - copied
		for i := 0; i < need; i++ {
			b := &nc.items[i]
			b.Raw_ = &b.inlineRaw
			b.Core_ = b.Raw_
			b.fromPool_ = true
			nb.ptrs[copied+i] = b
		}
		for i := need; i < 128; i++ {
			b := &nc.items[i]
			b.Raw_ = &b.inlineRaw
			b.Core_ = b.Raw_
			b.fromPool_ = true
			astNodeBuilderFree_ = append(astNodeBuilderFree_, b)
		}
	}
	ret := nb.ptrs[0]
	nb.done.Store(1)
	nb.idx.Store(1)
	astNodeBuilderCurBatch_ = nb
	astNodeBuilderBatch_.Store(nb)
	astNodeBuilderPoolMu_.Unlock()
	return ret
}
func recycleAstNodeBuilderTree_(root *Ext_com_google_javascript_jscomp_serialization_AstNode_Builder) {
	if root == nil || root.pinned_ {
		return
	}
	root.poolNextSibling_ = nil
	workHead := root
	var freeHead *Ext_com_google_javascript_jscomp_serialization_AstNode_Builder
	for workHead != nil {
		cur := workHead
		workHead = cur.poolNextSibling_
		cur.poolNextSibling_ = nil
		if cur.pinned_ {
			continue
		}
		children := cur.poolChildHead_
		cur.poolChildHead_ = nil
		for c := children; c != nil; {
			next := c.poolNextSibling_
			c.poolNextSibling_ = nil
			c.linked_ = false
			if !c.pinned_ {
				c.poolNextSibling_ = workHead
				workHead = c
			}
			c = next
		}
		cur.inlineRaw.ResetForPool()
		cur.Raw_ = &cur.inlineRaw
		cur.Core_ = cur.Raw_
		cur.builtMsg.Core_ = nil
		cur.built_ = false
		cur.usedBuiltMsg_ = false
		cur.linked_ = false
		cur.pinned_ = false
		cur.poolChildHead_ = freeHead
		freeHead = cur
	}
	if freeHead == nil {
		return
	}
	astNodeBuilderPoolMu_.Lock()
	for freeHead != nil {
		n := freeHead
		freeHead = n.poolChildHead_
		n.poolChildHead_ = nil
		if len(astNodeBuilderFree_) < 262144 {
			astNodeBuilderFree_ = append(astNodeBuilderFree_, n)
		}
	}
	astNodeBuilderPoolMu_.Unlock()
}
var protoDefault_ExternsSummary_ *Ext_com_google_javascript_jscomp_serialization_ExternsSummary
var protoDefaultRaw_ExternsSummary_ *protobuf.ExternsSummary
func (m *Ext_com_google_javascript_jscomp_serialization_ExternsSummary) p_ExternsSummary() *protobuf.ExternsSummary {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_ExternsSummary_ }
	return m.Core_.(*protobuf.ExternsSummary)
}
func (m *Ext_com_google_javascript_jscomp_serialization_ExternsSummary) buildRaw_ExternsSummary() *protobuf.ExternsSummary { return m.p_ExternsSummary() }
func (m *Ext_com_google_javascript_jscomp_serialization_ExternsSummary) protoBuildCore_() any { return m.buildRaw_ExternsSummary() }
func (b *Ext_com_google_javascript_jscomp_serialization_ExternsSummary_Builder) p_ExternsSummary() *protobuf.ExternsSummary {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.ExternsSummary)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_ExternsSummary_Builder) pMut_ExternsSummary() *protobuf.ExternsSummary {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_ExternsSummary()
}
func (b *Ext_com_google_javascript_jscomp_serialization_ExternsSummary_Builder) buildRaw_ExternsSummary() *protobuf.ExternsSummary {
	p := b.p_ExternsSummary()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_ExternsSummary_Builder) buildMsg_ExternsSummary() *Ext_com_google_javascript_jscomp_serialization_ExternsSummary {
	raw := b.buildRaw_ExternsSummary()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_ExternsSummary{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_ExternsSummary_Builder) protoMutCore_() any { return b.pMut_ExternsSummary() }
func (b *Ext_com_google_javascript_jscomp_serialization_ExternsSummary_Builder) protoBuildCore_() any { return b.buildRaw_ExternsSummary() }
var protoDefault_LazyAst_ *Ext_com_google_javascript_jscomp_serialization_LazyAst
var protoDefaultRaw_LazyAst_ *protobuf.LazyAst
func (m *Ext_com_google_javascript_jscomp_serialization_LazyAst) p_LazyAst() *protobuf.LazyAst {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_LazyAst_ }
	return m.Core_.(*protobuf.LazyAst)
}
func (m *Ext_com_google_javascript_jscomp_serialization_LazyAst) buildRaw_LazyAst() *protobuf.LazyAst { return m.p_LazyAst() }
func (m *Ext_com_google_javascript_jscomp_serialization_LazyAst) protoBuildCore_() any { return m.buildRaw_LazyAst() }
func (b *Ext_com_google_javascript_jscomp_serialization_LazyAst_Builder) p_LazyAst() *protobuf.LazyAst {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.LazyAst)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_LazyAst_Builder) pMut_LazyAst() *protobuf.LazyAst {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_LazyAst()
}
func (b *Ext_com_google_javascript_jscomp_serialization_LazyAst_Builder) buildRaw_LazyAst() *protobuf.LazyAst {
	p := b.p_LazyAst()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_LazyAst_Builder) buildMsg_LazyAst() *Ext_com_google_javascript_jscomp_serialization_LazyAst {
	raw := b.buildRaw_LazyAst()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_LazyAst{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_LazyAst_Builder) protoMutCore_() any { return b.pMut_LazyAst() }
func (b *Ext_com_google_javascript_jscomp_serialization_LazyAst_Builder) protoBuildCore_() any { return b.buildRaw_LazyAst() }
var protoDefault_ObjectTypeProto_ *Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto
var protoDefaultRaw_ObjectTypeProto_ *protobuf.ObjectTypeProto
func (m *Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto) p_ObjectTypeProto() *protobuf.ObjectTypeProto {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_ObjectTypeProto_ }
	return m.Core_.(*protobuf.ObjectTypeProto)
}
func (m *Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto) buildRaw_ObjectTypeProto() *protobuf.ObjectTypeProto { return m.p_ObjectTypeProto() }
func (m *Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto) protoBuildCore_() any { return m.buildRaw_ObjectTypeProto() }
func (b *Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto_Builder) p_ObjectTypeProto() *protobuf.ObjectTypeProto {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.ObjectTypeProto)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto_Builder) pMut_ObjectTypeProto() *protobuf.ObjectTypeProto {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_ObjectTypeProto()
}
func (b *Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto_Builder) buildRaw_ObjectTypeProto() *protobuf.ObjectTypeProto {
	p := b.p_ObjectTypeProto()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto_Builder) buildMsg_ObjectTypeProto() *Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto {
	raw := b.buildRaw_ObjectTypeProto()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto_Builder) protoMutCore_() any { return b.pMut_ObjectTypeProto() }
func (b *Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto_Builder) protoBuildCore_() any { return b.buildRaw_ObjectTypeProto() }
var protoDefault_OptimizationJsdoc_ *Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc
var protoDefaultRaw_OptimizationJsdoc_ *protobuf.OptimizationJsdoc
func (m *Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc) p_OptimizationJsdoc() *protobuf.OptimizationJsdoc {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_OptimizationJsdoc_ }
	return m.Core_.(*protobuf.OptimizationJsdoc)
}
func (m *Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc) buildRaw_OptimizationJsdoc() *protobuf.OptimizationJsdoc { return m.p_OptimizationJsdoc() }
func (m *Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc) protoBuildCore_() any { return m.buildRaw_OptimizationJsdoc() }
func (b *Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc_Builder) p_OptimizationJsdoc() *protobuf.OptimizationJsdoc {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.OptimizationJsdoc)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc_Builder) pMut_OptimizationJsdoc() *protobuf.OptimizationJsdoc {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_OptimizationJsdoc()
}
func (b *Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc_Builder) buildRaw_OptimizationJsdoc() *protobuf.OptimizationJsdoc {
	p := b.p_OptimizationJsdoc()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc_Builder) buildMsg_OptimizationJsdoc() *Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc {
	raw := b.buildRaw_OptimizationJsdoc()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc_Builder) protoMutCore_() any { return b.pMut_OptimizationJsdoc() }
func (b *Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc_Builder) protoBuildCore_() any { return b.buildRaw_OptimizationJsdoc() }
var protoDefault_SourceFilePool_ *Ext_com_google_javascript_jscomp_serialization_SourceFilePool
var protoDefaultRaw_SourceFilePool_ *protobuf.SourceFilePool
func (m *Ext_com_google_javascript_jscomp_serialization_SourceFilePool) p_SourceFilePool() *protobuf.SourceFilePool {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_SourceFilePool_ }
	return m.Core_.(*protobuf.SourceFilePool)
}
func (m *Ext_com_google_javascript_jscomp_serialization_SourceFilePool) buildRaw_SourceFilePool() *protobuf.SourceFilePool { return m.p_SourceFilePool() }
func (m *Ext_com_google_javascript_jscomp_serialization_SourceFilePool) protoBuildCore_() any { return m.buildRaw_SourceFilePool() }
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFilePool_Builder) p_SourceFilePool() *protobuf.SourceFilePool {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.SourceFilePool)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFilePool_Builder) pMut_SourceFilePool() *protobuf.SourceFilePool {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_SourceFilePool()
}
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFilePool_Builder) buildRaw_SourceFilePool() *protobuf.SourceFilePool {
	p := b.p_SourceFilePool()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFilePool_Builder) buildMsg_SourceFilePool() *Ext_com_google_javascript_jscomp_serialization_SourceFilePool {
	raw := b.buildRaw_SourceFilePool()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_SourceFilePool{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFilePool_Builder) protoMutCore_() any { return b.pMut_SourceFilePool() }
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFilePool_Builder) protoBuildCore_() any { return b.buildRaw_SourceFilePool() }
var protoDefault_SourceFileProto_ *Ext_com_google_javascript_jscomp_serialization_SourceFileProto
var protoDefaultRaw_SourceFileProto_ *protobuf.SourceFileProto
func (m *Ext_com_google_javascript_jscomp_serialization_SourceFileProto) p_SourceFileProto() *protobuf.SourceFileProto {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_SourceFileProto_ }
	return m.Core_.(*protobuf.SourceFileProto)
}
func (m *Ext_com_google_javascript_jscomp_serialization_SourceFileProto) buildRaw_SourceFileProto() *protobuf.SourceFileProto { return m.p_SourceFileProto() }
func (m *Ext_com_google_javascript_jscomp_serialization_SourceFileProto) protoBuildCore_() any { return m.buildRaw_SourceFileProto() }
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_Builder) p_SourceFileProto() *protobuf.SourceFileProto {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.SourceFileProto)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_Builder) pMut_SourceFileProto() *protobuf.SourceFileProto {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_SourceFileProto()
}
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_Builder) buildRaw_SourceFileProto() *protobuf.SourceFileProto {
	p := b.p_SourceFileProto()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_Builder) buildMsg_SourceFileProto() *Ext_com_google_javascript_jscomp_serialization_SourceFileProto {
	raw := b.buildRaw_SourceFileProto()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_SourceFileProto{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_Builder) protoMutCore_() any { return b.pMut_SourceFileProto() }
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_Builder) protoBuildCore_() any { return b.buildRaw_SourceFileProto() }
var protoDefault_SourceFileProto_FileOnDisk_ *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_FileOnDisk
var protoDefaultRaw_SourceFileProto_FileOnDisk_ *protobuf.SourceFileProto_FileOnDisk
func (m *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_FileOnDisk) p_SourceFileProto_FileOnDisk() *protobuf.SourceFileProto_FileOnDisk {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_SourceFileProto_FileOnDisk_ }
	return m.Core_.(*protobuf.SourceFileProto_FileOnDisk)
}
func (m *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_FileOnDisk) buildRaw_SourceFileProto_FileOnDisk() *protobuf.SourceFileProto_FileOnDisk { return m.p_SourceFileProto_FileOnDisk() }
func (m *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_FileOnDisk) protoBuildCore_() any { return m.buildRaw_SourceFileProto_FileOnDisk() }
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_FileOnDisk_Builder) p_SourceFileProto_FileOnDisk() *protobuf.SourceFileProto_FileOnDisk {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.SourceFileProto_FileOnDisk)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_FileOnDisk_Builder) pMut_SourceFileProto_FileOnDisk() *protobuf.SourceFileProto_FileOnDisk {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_SourceFileProto_FileOnDisk()
}
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_FileOnDisk_Builder) buildRaw_SourceFileProto_FileOnDisk() *protobuf.SourceFileProto_FileOnDisk {
	p := b.p_SourceFileProto_FileOnDisk()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_FileOnDisk_Builder) buildMsg_SourceFileProto_FileOnDisk() *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_FileOnDisk {
	raw := b.buildRaw_SourceFileProto_FileOnDisk()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_SourceFileProto_FileOnDisk{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_FileOnDisk_Builder) protoMutCore_() any { return b.pMut_SourceFileProto_FileOnDisk() }
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_FileOnDisk_Builder) protoBuildCore_() any { return b.buildRaw_SourceFileProto_FileOnDisk() }
var protoDefault_SourceFileProto_ZipEntryOnDisk_ *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk
var protoDefaultRaw_SourceFileProto_ZipEntryOnDisk_ *protobuf.SourceFileProto_ZipEntryOnDisk
func (m *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk) p_SourceFileProto_ZipEntryOnDisk() *protobuf.SourceFileProto_ZipEntryOnDisk {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_SourceFileProto_ZipEntryOnDisk_ }
	return m.Core_.(*protobuf.SourceFileProto_ZipEntryOnDisk)
}
func (m *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk) buildRaw_SourceFileProto_ZipEntryOnDisk() *protobuf.SourceFileProto_ZipEntryOnDisk { return m.p_SourceFileProto_ZipEntryOnDisk() }
func (m *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk) protoBuildCore_() any { return m.buildRaw_SourceFileProto_ZipEntryOnDisk() }
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk_Builder) p_SourceFileProto_ZipEntryOnDisk() *protobuf.SourceFileProto_ZipEntryOnDisk {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.SourceFileProto_ZipEntryOnDisk)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk_Builder) pMut_SourceFileProto_ZipEntryOnDisk() *protobuf.SourceFileProto_ZipEntryOnDisk {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_SourceFileProto_ZipEntryOnDisk()
}
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk_Builder) buildRaw_SourceFileProto_ZipEntryOnDisk() *protobuf.SourceFileProto_ZipEntryOnDisk {
	p := b.p_SourceFileProto_ZipEntryOnDisk()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk_Builder) buildMsg_SourceFileProto_ZipEntryOnDisk() *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk {
	raw := b.buildRaw_SourceFileProto_ZipEntryOnDisk()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk_Builder) protoMutCore_() any { return b.pMut_SourceFileProto_ZipEntryOnDisk() }
func (b *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk_Builder) protoBuildCore_() any { return b.buildRaw_SourceFileProto_ZipEntryOnDisk() }
var protoDefault_StringPoolProto_ *Ext_com_google_javascript_jscomp_serialization_StringPoolProto
var protoDefaultRaw_StringPoolProto_ *protobuf.StringPoolProto
func (m *Ext_com_google_javascript_jscomp_serialization_StringPoolProto) p_StringPoolProto() *protobuf.StringPoolProto {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_StringPoolProto_ }
	return m.Core_.(*protobuf.StringPoolProto)
}
func (m *Ext_com_google_javascript_jscomp_serialization_StringPoolProto) buildRaw_StringPoolProto() *protobuf.StringPoolProto { return m.p_StringPoolProto() }
func (m *Ext_com_google_javascript_jscomp_serialization_StringPoolProto) protoBuildCore_() any { return m.buildRaw_StringPoolProto() }
func (b *Ext_com_google_javascript_jscomp_serialization_StringPoolProto_Builder) p_StringPoolProto() *protobuf.StringPoolProto {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.StringPoolProto)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_StringPoolProto_Builder) pMut_StringPoolProto() *protobuf.StringPoolProto {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_StringPoolProto()
}
func (b *Ext_com_google_javascript_jscomp_serialization_StringPoolProto_Builder) buildRaw_StringPoolProto() *protobuf.StringPoolProto {
	p := b.p_StringPoolProto()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_StringPoolProto_Builder) buildMsg_StringPoolProto() *Ext_com_google_javascript_jscomp_serialization_StringPoolProto {
	raw := b.buildRaw_StringPoolProto()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_StringPoolProto{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_StringPoolProto_Builder) protoMutCore_() any { return b.pMut_StringPoolProto() }
func (b *Ext_com_google_javascript_jscomp_serialization_StringPoolProto_Builder) protoBuildCore_() any { return b.buildRaw_StringPoolProto() }
var protoDefault_SubtypingEdge_ *Ext_com_google_javascript_jscomp_serialization_SubtypingEdge
var protoDefaultRaw_SubtypingEdge_ *protobuf.SubtypingEdge
func (m *Ext_com_google_javascript_jscomp_serialization_SubtypingEdge) p_SubtypingEdge() *protobuf.SubtypingEdge {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_SubtypingEdge_ }
	return m.Core_.(*protobuf.SubtypingEdge)
}
func (m *Ext_com_google_javascript_jscomp_serialization_SubtypingEdge) buildRaw_SubtypingEdge() *protobuf.SubtypingEdge { return m.p_SubtypingEdge() }
func (m *Ext_com_google_javascript_jscomp_serialization_SubtypingEdge) protoBuildCore_() any { return m.buildRaw_SubtypingEdge() }
func (b *Ext_com_google_javascript_jscomp_serialization_SubtypingEdge_Builder) p_SubtypingEdge() *protobuf.SubtypingEdge {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.SubtypingEdge)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_SubtypingEdge_Builder) pMut_SubtypingEdge() *protobuf.SubtypingEdge {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_SubtypingEdge()
}
func (b *Ext_com_google_javascript_jscomp_serialization_SubtypingEdge_Builder) buildRaw_SubtypingEdge() *protobuf.SubtypingEdge {
	p := b.p_SubtypingEdge()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_SubtypingEdge_Builder) buildMsg_SubtypingEdge() *Ext_com_google_javascript_jscomp_serialization_SubtypingEdge {
	raw := b.buildRaw_SubtypingEdge()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_SubtypingEdge{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_SubtypingEdge_Builder) protoMutCore_() any { return b.pMut_SubtypingEdge() }
func (b *Ext_com_google_javascript_jscomp_serialization_SubtypingEdge_Builder) protoBuildCore_() any { return b.buildRaw_SubtypingEdge() }
var protoDefault_TemplateStringValue_ *Ext_com_google_javascript_jscomp_serialization_TemplateStringValue
var protoDefaultRaw_TemplateStringValue_ *protobuf.TemplateStringValue
func (m *Ext_com_google_javascript_jscomp_serialization_TemplateStringValue) p_TemplateStringValue() *protobuf.TemplateStringValue {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_TemplateStringValue_ }
	return m.Core_.(*protobuf.TemplateStringValue)
}
func (m *Ext_com_google_javascript_jscomp_serialization_TemplateStringValue) buildRaw_TemplateStringValue() *protobuf.TemplateStringValue { return m.p_TemplateStringValue() }
func (m *Ext_com_google_javascript_jscomp_serialization_TemplateStringValue) protoBuildCore_() any { return m.buildRaw_TemplateStringValue() }
func (b *Ext_com_google_javascript_jscomp_serialization_TemplateStringValue_Builder) p_TemplateStringValue() *protobuf.TemplateStringValue {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.TemplateStringValue)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_TemplateStringValue_Builder) pMut_TemplateStringValue() *protobuf.TemplateStringValue {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_TemplateStringValue()
}
func (b *Ext_com_google_javascript_jscomp_serialization_TemplateStringValue_Builder) buildRaw_TemplateStringValue() *protobuf.TemplateStringValue {
	p := b.p_TemplateStringValue()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_TemplateStringValue_Builder) buildMsg_TemplateStringValue() *Ext_com_google_javascript_jscomp_serialization_TemplateStringValue {
	raw := b.buildRaw_TemplateStringValue()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_TemplateStringValue{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_TemplateStringValue_Builder) protoMutCore_() any { return b.pMut_TemplateStringValue() }
func (b *Ext_com_google_javascript_jscomp_serialization_TemplateStringValue_Builder) protoBuildCore_() any { return b.buildRaw_TemplateStringValue() }
var protoDefault_TypePool_ *Ext_com_google_javascript_jscomp_serialization_TypePool
var protoDefaultRaw_TypePool_ *protobuf.TypePool
func (m *Ext_com_google_javascript_jscomp_serialization_TypePool) p_TypePool() *protobuf.TypePool {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_TypePool_ }
	return m.Core_.(*protobuf.TypePool)
}
func (m *Ext_com_google_javascript_jscomp_serialization_TypePool) buildRaw_TypePool() *protobuf.TypePool { return m.p_TypePool() }
func (m *Ext_com_google_javascript_jscomp_serialization_TypePool) protoBuildCore_() any { return m.buildRaw_TypePool() }
func (this_ *Ext_com_google_javascript_jscomp_serialization_TypePool) getOrBuildTypeList_() *jProtoList {
	if this_ == nil { ThrowNPE() }
	if p := this_.cache_GetTypeList.Load(); p != nil { return p }
	raw := this_.p_TypePool().GetTypeList()
	if len(raw) == 0 {
		this_.cache_GetTypeList.CompareAndSwap(nil, emptyProtoList_)
		return emptyProtoList_
	}
	wrappers := make([]Ext_com_google_javascript_jscomp_serialization_TypeProto, len(raw))
	items := make([]any, len(raw))
	for i, x := range raw {
		wrappers[i].Core_ = x
		items[i] = &wrappers[i]
	}
	l := newProtoList(items)
	if this_.cache_GetTypeList.CompareAndSwap(nil, l) { return l }
	return this_.cache_GetTypeList.Load()
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypePool_Builder) p_TypePool() *protobuf.TypePool {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.TypePool)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypePool_Builder) pMut_TypePool() *protobuf.TypePool {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_TypePool()
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypePool_Builder) buildRaw_TypePool() *protobuf.TypePool {
	p := b.p_TypePool()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypePool_Builder) buildMsg_TypePool() *Ext_com_google_javascript_jscomp_serialization_TypePool {
	raw := b.buildRaw_TypePool()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_TypePool{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypePool_Builder) protoMutCore_() any { return b.pMut_TypePool() }
func (b *Ext_com_google_javascript_jscomp_serialization_TypePool_Builder) protoBuildCore_() any { return b.buildRaw_TypePool() }
var protoDefault_TypePool_DebugInfo_ *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo
var protoDefaultRaw_TypePool_DebugInfo_ *protobuf.TypePool_DebugInfo
func (m *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo) p_TypePool_DebugInfo() *protobuf.TypePool_DebugInfo {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_TypePool_DebugInfo_ }
	return m.Core_.(*protobuf.TypePool_DebugInfo)
}
func (m *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo) buildRaw_TypePool_DebugInfo() *protobuf.TypePool_DebugInfo { return m.p_TypePool_DebugInfo() }
func (m *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo) protoBuildCore_() any { return m.buildRaw_TypePool_DebugInfo() }
func (b *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Builder) p_TypePool_DebugInfo() *protobuf.TypePool_DebugInfo {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.TypePool_DebugInfo)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Builder) pMut_TypePool_DebugInfo() *protobuf.TypePool_DebugInfo {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_TypePool_DebugInfo()
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Builder) buildRaw_TypePool_DebugInfo() *protobuf.TypePool_DebugInfo {
	p := b.p_TypePool_DebugInfo()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Builder) buildMsg_TypePool_DebugInfo() *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo {
	raw := b.buildRaw_TypePool_DebugInfo()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Builder) protoMutCore_() any { return b.pMut_TypePool_DebugInfo() }
func (b *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Builder) protoBuildCore_() any { return b.buildRaw_TypePool_DebugInfo() }
var protoDefault_TypePool_DebugInfo_Mismatch_ *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Mismatch
var protoDefaultRaw_TypePool_DebugInfo_Mismatch_ *protobuf.TypePool_DebugInfo_Mismatch
func (m *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Mismatch) p_TypePool_DebugInfo_Mismatch() *protobuf.TypePool_DebugInfo_Mismatch {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_TypePool_DebugInfo_Mismatch_ }
	return m.Core_.(*protobuf.TypePool_DebugInfo_Mismatch)
}
func (m *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Mismatch) buildRaw_TypePool_DebugInfo_Mismatch() *protobuf.TypePool_DebugInfo_Mismatch { return m.p_TypePool_DebugInfo_Mismatch() }
func (m *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Mismatch) protoBuildCore_() any { return m.buildRaw_TypePool_DebugInfo_Mismatch() }
func (b *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Mismatch_Builder) p_TypePool_DebugInfo_Mismatch() *protobuf.TypePool_DebugInfo_Mismatch {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.TypePool_DebugInfo_Mismatch)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Mismatch_Builder) pMut_TypePool_DebugInfo_Mismatch() *protobuf.TypePool_DebugInfo_Mismatch {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_TypePool_DebugInfo_Mismatch()
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Mismatch_Builder) buildRaw_TypePool_DebugInfo_Mismatch() *protobuf.TypePool_DebugInfo_Mismatch {
	p := b.p_TypePool_DebugInfo_Mismatch()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Mismatch_Builder) buildMsg_TypePool_DebugInfo_Mismatch() *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Mismatch {
	raw := b.buildRaw_TypePool_DebugInfo_Mismatch()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Mismatch{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Mismatch_Builder) protoMutCore_() any { return b.pMut_TypePool_DebugInfo_Mismatch() }
func (b *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Mismatch_Builder) protoBuildCore_() any { return b.buildRaw_TypePool_DebugInfo_Mismatch() }
var protoDefault_TypeProto_ *Ext_com_google_javascript_jscomp_serialization_TypeProto
var protoDefaultRaw_TypeProto_ *protobuf.TypeProto
func (m *Ext_com_google_javascript_jscomp_serialization_TypeProto) p_TypeProto() *protobuf.TypeProto {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_TypeProto_ }
	return m.Core_.(*protobuf.TypeProto)
}
func (m *Ext_com_google_javascript_jscomp_serialization_TypeProto) buildRaw_TypeProto() *protobuf.TypeProto { return m.p_TypeProto() }
func (m *Ext_com_google_javascript_jscomp_serialization_TypeProto) protoBuildCore_() any { return m.buildRaw_TypeProto() }
func (b *Ext_com_google_javascript_jscomp_serialization_TypeProto_Builder) p_TypeProto() *protobuf.TypeProto {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.TypeProto)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypeProto_Builder) pMut_TypeProto() *protobuf.TypeProto {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_TypeProto()
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypeProto_Builder) buildRaw_TypeProto() *protobuf.TypeProto {
	p := b.p_TypeProto()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypeProto_Builder) buildMsg_TypeProto() *Ext_com_google_javascript_jscomp_serialization_TypeProto {
	raw := b.buildRaw_TypeProto()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_TypeProto{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypeProto_Builder) protoMutCore_() any { return b.pMut_TypeProto() }
func (b *Ext_com_google_javascript_jscomp_serialization_TypeProto_Builder) protoBuildCore_() any { return b.buildRaw_TypeProto() }
var protoDefault_TypedAst_ *Ext_com_google_javascript_jscomp_serialization_TypedAst
var protoDefaultRaw_TypedAst_ *protobuf.TypedAst
func (m *Ext_com_google_javascript_jscomp_serialization_TypedAst) p_TypedAst() *protobuf.TypedAst {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_TypedAst_ }
	return m.Core_.(*protobuf.TypedAst)
}
func (m *Ext_com_google_javascript_jscomp_serialization_TypedAst) buildRaw_TypedAst() *protobuf.TypedAst { return m.p_TypedAst() }
func (m *Ext_com_google_javascript_jscomp_serialization_TypedAst) protoBuildCore_() any { return m.buildRaw_TypedAst() }
func (b *Ext_com_google_javascript_jscomp_serialization_TypedAst_Builder) p_TypedAst() *protobuf.TypedAst {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.TypedAst)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypedAst_Builder) pMut_TypedAst() *protobuf.TypedAst {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_TypedAst()
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypedAst_Builder) buildRaw_TypedAst() *protobuf.TypedAst {
	p := b.p_TypedAst()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypedAst_Builder) buildMsg_TypedAst() *Ext_com_google_javascript_jscomp_serialization_TypedAst {
	raw := b.buildRaw_TypedAst()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_TypedAst{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypedAst_Builder) protoMutCore_() any { return b.pMut_TypedAst() }
func (b *Ext_com_google_javascript_jscomp_serialization_TypedAst_Builder) protoBuildCore_() any { return b.buildRaw_TypedAst() }
var protoDefault_TypedAst_List_ *Ext_com_google_javascript_jscomp_serialization_TypedAst_List
var protoDefaultRaw_TypedAst_List_ *protobuf.TypedAst_List
func (m *Ext_com_google_javascript_jscomp_serialization_TypedAst_List) p_TypedAst_List() *protobuf.TypedAst_List {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_TypedAst_List_ }
	return m.Core_.(*protobuf.TypedAst_List)
}
func (m *Ext_com_google_javascript_jscomp_serialization_TypedAst_List) buildRaw_TypedAst_List() *protobuf.TypedAst_List { return m.p_TypedAst_List() }
func (m *Ext_com_google_javascript_jscomp_serialization_TypedAst_List) protoBuildCore_() any { return m.buildRaw_TypedAst_List() }
func (b *Ext_com_google_javascript_jscomp_serialization_TypedAst_List_Builder) p_TypedAst_List() *protobuf.TypedAst_List {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.TypedAst_List)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypedAst_List_Builder) pMut_TypedAst_List() *protobuf.TypedAst_List {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_TypedAst_List()
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypedAst_List_Builder) buildRaw_TypedAst_List() *protobuf.TypedAst_List {
	p := b.p_TypedAst_List()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypedAst_List_Builder) buildMsg_TypedAst_List() *Ext_com_google_javascript_jscomp_serialization_TypedAst_List {
	raw := b.buildRaw_TypedAst_List()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_TypedAst_List{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_TypedAst_List_Builder) protoMutCore_() any { return b.pMut_TypedAst_List() }
func (b *Ext_com_google_javascript_jscomp_serialization_TypedAst_List_Builder) protoBuildCore_() any { return b.buildRaw_TypedAst_List() }
var protoDefault_UnionTypeProto_ *Ext_com_google_javascript_jscomp_serialization_UnionTypeProto
var protoDefaultRaw_UnionTypeProto_ *protobuf.UnionTypeProto
func (m *Ext_com_google_javascript_jscomp_serialization_UnionTypeProto) p_UnionTypeProto() *protobuf.UnionTypeProto {
	if m == nil { ThrowNPE() }
	if m.Core_ == nil { return protoDefaultRaw_UnionTypeProto_ }
	return m.Core_.(*protobuf.UnionTypeProto)
}
func (m *Ext_com_google_javascript_jscomp_serialization_UnionTypeProto) buildRaw_UnionTypeProto() *protobuf.UnionTypeProto { return m.p_UnionTypeProto() }
func (m *Ext_com_google_javascript_jscomp_serialization_UnionTypeProto) protoBuildCore_() any { return m.buildRaw_UnionTypeProto() }
func (b *Ext_com_google_javascript_jscomp_serialization_UnionTypeProto_Builder) p_UnionTypeProto() *protobuf.UnionTypeProto {
	if b == nil { ThrowNPE() }
	if b.Raw_ != nil { return b.Raw_ }
	if b.Core_ != nil {
		r := b.Core_.(*protobuf.UnionTypeProto)
		b.Raw_ = r
		return r
	}
	b.Raw_ = &b.inlineRaw
	b.Core_ = b.Raw_
	return b.Raw_
}
func (b *Ext_com_google_javascript_jscomp_serialization_UnionTypeProto_Builder) pMut_UnionTypeProto() *protobuf.UnionTypeProto {
	if b == nil { ThrowNPE() }
	if !b.built_ && b.Raw_ != nil { return b.Raw_ }
	if b.built_ {
		cp := *b.Raw_
		b.Raw_ = &cp
		b.Core_ = b.Raw_
		b.built_ = false
		return b.Raw_
	}
	return b.p_UnionTypeProto()
}
func (b *Ext_com_google_javascript_jscomp_serialization_UnionTypeProto_Builder) buildRaw_UnionTypeProto() *protobuf.UnionTypeProto {
	p := b.p_UnionTypeProto()
	out := p.Build()
	if out == p { b.built_ = true }
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_UnionTypeProto_Builder) buildMsg_UnionTypeProto() *Ext_com_google_javascript_jscomp_serialization_UnionTypeProto {
	raw := b.buildRaw_UnionTypeProto()
	if !b.usedBuiltMsg_ {
		b.usedBuiltMsg_ = true
		b.builtMsg.Core_ = raw
		return &b.builtMsg
	}
	if b.builtMsg.Core_ == raw {
		return &b.builtMsg
	}
	out := &Ext_com_google_javascript_jscomp_serialization_UnionTypeProto{}
	out.Core_ = raw
	return out
}
func (b *Ext_com_google_javascript_jscomp_serialization_UnionTypeProto_Builder) protoMutCore_() any { return b.pMut_UnionTypeProto() }
func (b *Ext_com_google_javascript_jscomp_serialization_UnionTypeProto_Builder) protoBuildCore_() any { return b.buildRaw_UnionTypeProto() }
var protoDefSub_JSCompilerStateProto_GetAccessorSummary *Ext_com_google_javascript_jscomp_AccessorSummaryProto
var protoDefSub_InstrumentationPointStats_GetPoint *Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint
var protoDefSub_AstNode_GetTemplateStringValue *Ext_com_google_javascript_jscomp_serialization_TemplateStringValue
var protoDefSub_AstNode_GetJsdoc *Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc
var protoDefSub_SourceFileProto_GetFileOnDisk *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_FileOnDisk
var protoDefSub_SourceFileProto_GetZipEntry *Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk
var protoDefSub_TypePool_GetDebugInfo *Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo
var protoDefSub_TypeProto_GetObject *Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto
var protoDefSub_TypeProto_GetUnion *Ext_com_google_javascript_jscomp_serialization_UnionTypeProto
var protoDefSub_TypedAst_GetTypePool *Ext_com_google_javascript_jscomp_serialization_TypePool
var protoDefSub_TypedAst_GetStringPool *Ext_com_google_javascript_jscomp_serialization_StringPoolProto
var protoDefSub_TypedAst_GetSourceFilePool *Ext_com_google_javascript_jscomp_serialization_SourceFilePool
var protoDefSub_TypedAst_GetExternsSummary *Ext_com_google_javascript_jscomp_serialization_ExternsSummary
func init() {
	protoDefaultRaw_Mapping_OriginalMapping_ = protobuf.Mapping_OriginalMapping_GetDefaultInstance()
	protoDefault_Mapping_OriginalMapping_ = &Ext_com_google_debugging_sourcemap_proto_Mapping_OriginalMapping{}
	protoDefault_Mapping_OriginalMapping_.Core_ = protoDefaultRaw_Mapping_OriginalMapping_
	protoDefaultRaw_AccessorSummaryEntryProto_ = &protobuf.AccessorSummaryEntryProto{}
	protoDefault_AccessorSummaryEntryProto_ = &Ext_com_google_javascript_jscomp_AccessorSummaryEntryProto{}
	protoDefault_AccessorSummaryEntryProto_.Core_ = protoDefaultRaw_AccessorSummaryEntryProto_
	protoDefaultRaw_AccessorSummaryProto_ = (&protobuf.JSCompilerStateProto{}).GetAccessorSummary()
	protoDefault_AccessorSummaryProto_ = &Ext_com_google_javascript_jscomp_AccessorSummaryProto{}
	protoDefault_AccessorSummaryProto_.Core_ = protoDefaultRaw_AccessorSummaryProto_
	protoDefaultRaw_ChunkProto_ = &protobuf.ChunkProto{}
	protoDefault_ChunkProto_ = &Ext_com_google_javascript_jscomp_ChunkProto{}
	protoDefault_ChunkProto_.Core_ = protoDefaultRaw_ChunkProto_
	protoDefaultRaw_ConformanceConfig_ = &protobuf.ConformanceConfig{}
	protoDefault_ConformanceConfig_ = &Ext_com_google_javascript_jscomp_ConformanceConfig{}
	protoDefault_ConformanceConfig_.Core_ = protoDefaultRaw_ConformanceConfig_
	protoDefaultRaw_JSCompilerStateProto_ = &protobuf.JSCompilerStateProto{}
	protoDefault_JSCompilerStateProto_ = &Ext_com_google_javascript_jscomp_JSCompilerStateProto{}
	protoDefault_JSCompilerStateProto_.Core_ = protoDefaultRaw_JSCompilerStateProto_
	protoDefaultRaw_Requirement_ = &protobuf.Requirement{}
	protoDefault_Requirement_ = &Ext_com_google_javascript_jscomp_Requirement{}
	protoDefault_Requirement_.Core_ = protoDefaultRaw_Requirement_
	protoDefaultRaw_RequirementScopeEntry_ = &protobuf.RequirementScopeEntry{}
	protoDefault_RequirementScopeEntry_ = &Ext_com_google_javascript_jscomp_RequirementScopeEntry{}
	protoDefault_RequirementScopeEntry_.Core_ = protoDefaultRaw_RequirementScopeEntry_
	protoDefaultRaw_UniqueIdProto_ = &protobuf.UniqueIdProto{}
	protoDefault_UniqueIdProto_ = &Ext_com_google_javascript_jscomp_UniqueIdProto{}
	protoDefault_UniqueIdProto_.Core_ = protoDefaultRaw_UniqueIdProto_
	protoDefaultRaw_VariableMapEntryProto_ = &protobuf.VariableMapEntryProto{}
	protoDefault_VariableMapEntryProto_ = &Ext_com_google_javascript_jscomp_VariableMapEntryProto{}
	protoDefault_VariableMapEntryProto_.Core_ = protoDefaultRaw_VariableMapEntryProto_
	protoDefaultRaw_FileProfile_ = &protobuf.FileProfile{}
	protoDefault_FileProfile_ = &Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_FileProfile{}
	protoDefault_FileProfile_.Core_ = protoDefaultRaw_FileProfile_
	protoDefaultRaw_InstrumentationPoint_ = &protobuf.InstrumentationPoint{}
	protoDefault_InstrumentationPoint_ = &Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint{}
	protoDefault_InstrumentationPoint_.Core_ = protoDefaultRaw_InstrumentationPoint_
	protoDefaultRaw_InstrumentationPointStats_ = &protobuf.InstrumentationPointStats{}
	protoDefault_InstrumentationPointStats_ = &Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPointStats{}
	protoDefault_InstrumentationPointStats_.Core_ = protoDefaultRaw_InstrumentationPointStats_
	protoDefaultRaw_ReportProfile_ = &protobuf.ReportProfile{}
	protoDefault_ReportProfile_ = &Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_ReportProfile{}
	protoDefault_ReportProfile_.Core_ = protoDefaultRaw_ReportProfile_
	protoDefaultRaw_AstNode_ = &protobuf.AstNode{}
	protoDefault_AstNode_ = &Ext_com_google_javascript_jscomp_serialization_AstNode{}
	protoDefault_AstNode_.Core_ = protoDefaultRaw_AstNode_
	registerCompPoolReset(func() {
		astNodeBuilderPoolMu_.Lock()
		astNodeBuilderBatch_.Store(nil)
		astNodeBuilderCurBatch_ = nil
		astNodeBuilderBatchFree_ = nil
		astNodeBuilderRetired_ = nil
		astNodeBuilderFree_ = nil
		astNodeBuilderPoolMu_.Unlock()
	})
	protoDefaultRaw_ExternsSummary_ = (&protobuf.TypedAst{}).GetExternsSummary()
	protoDefault_ExternsSummary_ = &Ext_com_google_javascript_jscomp_serialization_ExternsSummary{}
	protoDefault_ExternsSummary_.Core_ = protoDefaultRaw_ExternsSummary_
	protoDefaultRaw_LazyAst_ = &protobuf.LazyAst{}
	protoDefault_LazyAst_ = &Ext_com_google_javascript_jscomp_serialization_LazyAst{}
	protoDefault_LazyAst_.Core_ = protoDefaultRaw_LazyAst_
	protoDefaultRaw_ObjectTypeProto_ = (&protobuf.TypeProto{}).GetObject()
	protoDefault_ObjectTypeProto_ = &Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto{}
	protoDefault_ObjectTypeProto_.Core_ = protoDefaultRaw_ObjectTypeProto_
	protoDefaultRaw_OptimizationJsdoc_ = protobuf.OptimizationJsdoc_GetDefaultInstance()
	protoDefault_OptimizationJsdoc_ = &Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc{}
	protoDefault_OptimizationJsdoc_.Core_ = protoDefaultRaw_OptimizationJsdoc_
	protoDefaultRaw_SourceFilePool_ = (&protobuf.TypedAst{}).GetSourceFilePool()
	protoDefault_SourceFilePool_ = &Ext_com_google_javascript_jscomp_serialization_SourceFilePool{}
	protoDefault_SourceFilePool_.Core_ = protoDefaultRaw_SourceFilePool_
	protoDefaultRaw_SourceFileProto_ = &protobuf.SourceFileProto{}
	protoDefault_SourceFileProto_ = &Ext_com_google_javascript_jscomp_serialization_SourceFileProto{}
	protoDefault_SourceFileProto_.Core_ = protoDefaultRaw_SourceFileProto_
	protoDefaultRaw_SourceFileProto_FileOnDisk_ = (&protobuf.SourceFileProto{}).GetFileOnDisk()
	protoDefault_SourceFileProto_FileOnDisk_ = &Ext_com_google_javascript_jscomp_serialization_SourceFileProto_FileOnDisk{}
	protoDefault_SourceFileProto_FileOnDisk_.Core_ = protoDefaultRaw_SourceFileProto_FileOnDisk_
	protoDefaultRaw_SourceFileProto_ZipEntryOnDisk_ = (&protobuf.SourceFileProto{}).GetZipEntry()
	protoDefault_SourceFileProto_ZipEntryOnDisk_ = &Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk{}
	protoDefault_SourceFileProto_ZipEntryOnDisk_.Core_ = protoDefaultRaw_SourceFileProto_ZipEntryOnDisk_
	protoDefaultRaw_StringPoolProto_ = (&protobuf.TypedAst{}).GetStringPool()
	protoDefault_StringPoolProto_ = &Ext_com_google_javascript_jscomp_serialization_StringPoolProto{}
	protoDefault_StringPoolProto_.Core_ = protoDefaultRaw_StringPoolProto_
	protoDefaultRaw_SubtypingEdge_ = &protobuf.SubtypingEdge{}
	protoDefault_SubtypingEdge_ = &Ext_com_google_javascript_jscomp_serialization_SubtypingEdge{}
	protoDefault_SubtypingEdge_.Core_ = protoDefaultRaw_SubtypingEdge_
	protoDefaultRaw_TemplateStringValue_ = (&protobuf.AstNode{}).GetTemplateStringValue()
	protoDefault_TemplateStringValue_ = &Ext_com_google_javascript_jscomp_serialization_TemplateStringValue{}
	protoDefault_TemplateStringValue_.Core_ = protoDefaultRaw_TemplateStringValue_
	protoDefaultRaw_TypePool_ = protobuf.TypePool_GetDefaultInstance()
	protoDefault_TypePool_ = &Ext_com_google_javascript_jscomp_serialization_TypePool{}
	protoDefault_TypePool_.Core_ = protoDefaultRaw_TypePool_
	protoDefaultRaw_TypePool_DebugInfo_ = (&protobuf.TypePool{}).GetDebugInfo()
	protoDefault_TypePool_DebugInfo_ = &Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo{}
	protoDefault_TypePool_DebugInfo_.Core_ = protoDefaultRaw_TypePool_DebugInfo_
	protoDefaultRaw_TypePool_DebugInfo_Mismatch_ = &protobuf.TypePool_DebugInfo_Mismatch{}
	protoDefault_TypePool_DebugInfo_Mismatch_ = &Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo_Mismatch{}
	protoDefault_TypePool_DebugInfo_Mismatch_.Core_ = protoDefaultRaw_TypePool_DebugInfo_Mismatch_
	protoDefaultRaw_TypeProto_ = &protobuf.TypeProto{}
	protoDefault_TypeProto_ = &Ext_com_google_javascript_jscomp_serialization_TypeProto{}
	protoDefault_TypeProto_.Core_ = protoDefaultRaw_TypeProto_
	protoDefaultRaw_TypedAst_ = &protobuf.TypedAst{}
	protoDefault_TypedAst_ = &Ext_com_google_javascript_jscomp_serialization_TypedAst{}
	protoDefault_TypedAst_.Core_ = protoDefaultRaw_TypedAst_
	protoDefaultRaw_TypedAst_List_ = &protobuf.TypedAst_List{}
	protoDefault_TypedAst_List_ = &Ext_com_google_javascript_jscomp_serialization_TypedAst_List{}
	protoDefault_TypedAst_List_.Core_ = protoDefaultRaw_TypedAst_List_
	protoDefaultRaw_UnionTypeProto_ = (&protobuf.TypeProto{}).GetUnion()
	protoDefault_UnionTypeProto_ = &Ext_com_google_javascript_jscomp_serialization_UnionTypeProto{}
	protoDefault_UnionTypeProto_.Core_ = protoDefaultRaw_UnionTypeProto_
	if def := (&protobuf.JSCompilerStateProto{}).GetAccessorSummary(); def != nil {
		if protoDefault_AccessorSummaryProto_ != nil && protoDefault_AccessorSummaryProto_.Core_ == def {
			protoDefSub_JSCompilerStateProto_GetAccessorSummary = protoDefault_AccessorSummaryProto_
		} else {
			w := &Ext_com_google_javascript_jscomp_AccessorSummaryProto{}
			w.Core_ = def
			protoDefSub_JSCompilerStateProto_GetAccessorSummary = w
		}
	}
	if def := (&protobuf.InstrumentationPointStats{}).GetPoint(); def != nil {
		if protoDefault_InstrumentationPoint_ != nil && protoDefault_InstrumentationPoint_.Core_ == def {
			protoDefSub_InstrumentationPointStats_GetPoint = protoDefault_InstrumentationPoint_
		} else {
			w := &Ext_com_google_javascript_jscomp_instrumentation_reporter_proto_InstrumentationPoint{}
			w.Core_ = def
			protoDefSub_InstrumentationPointStats_GetPoint = w
		}
	}
	if def := (&protobuf.AstNode{}).GetTemplateStringValue(); def != nil {
		if protoDefault_TemplateStringValue_ != nil && protoDefault_TemplateStringValue_.Core_ == def {
			protoDefSub_AstNode_GetTemplateStringValue = protoDefault_TemplateStringValue_
		} else {
			w := &Ext_com_google_javascript_jscomp_serialization_TemplateStringValue{}
			w.Core_ = def
			protoDefSub_AstNode_GetTemplateStringValue = w
		}
	}
	if def := (&protobuf.AstNode{}).GetJsdoc(); def != nil {
		if protoDefault_OptimizationJsdoc_ != nil && protoDefault_OptimizationJsdoc_.Core_ == def {
			protoDefSub_AstNode_GetJsdoc = protoDefault_OptimizationJsdoc_
		} else {
			w := &Ext_com_google_javascript_jscomp_serialization_OptimizationJsdoc{}
			w.Core_ = def
			protoDefSub_AstNode_GetJsdoc = w
		}
	}
	if def := (&protobuf.SourceFileProto{}).GetFileOnDisk(); def != nil {
		if protoDefault_SourceFileProto_FileOnDisk_ != nil && protoDefault_SourceFileProto_FileOnDisk_.Core_ == def {
			protoDefSub_SourceFileProto_GetFileOnDisk = protoDefault_SourceFileProto_FileOnDisk_
		} else {
			w := &Ext_com_google_javascript_jscomp_serialization_SourceFileProto_FileOnDisk{}
			w.Core_ = def
			protoDefSub_SourceFileProto_GetFileOnDisk = w
		}
	}
	if def := (&protobuf.SourceFileProto{}).GetZipEntry(); def != nil {
		if protoDefault_SourceFileProto_ZipEntryOnDisk_ != nil && protoDefault_SourceFileProto_ZipEntryOnDisk_.Core_ == def {
			protoDefSub_SourceFileProto_GetZipEntry = protoDefault_SourceFileProto_ZipEntryOnDisk_
		} else {
			w := &Ext_com_google_javascript_jscomp_serialization_SourceFileProto_ZipEntryOnDisk{}
			w.Core_ = def
			protoDefSub_SourceFileProto_GetZipEntry = w
		}
	}
	if def := (&protobuf.TypePool{}).GetDebugInfo(); def != nil {
		if protoDefault_TypePool_DebugInfo_ != nil && protoDefault_TypePool_DebugInfo_.Core_ == def {
			protoDefSub_TypePool_GetDebugInfo = protoDefault_TypePool_DebugInfo_
		} else {
			w := &Ext_com_google_javascript_jscomp_serialization_TypePool_DebugInfo{}
			w.Core_ = def
			protoDefSub_TypePool_GetDebugInfo = w
		}
	}
	if def := (&protobuf.TypeProto{}).GetObject(); def != nil {
		if protoDefault_ObjectTypeProto_ != nil && protoDefault_ObjectTypeProto_.Core_ == def {
			protoDefSub_TypeProto_GetObject = protoDefault_ObjectTypeProto_
		} else {
			w := &Ext_com_google_javascript_jscomp_serialization_ObjectTypeProto{}
			w.Core_ = def
			protoDefSub_TypeProto_GetObject = w
		}
	}
	if def := (&protobuf.TypeProto{}).GetUnion(); def != nil {
		if protoDefault_UnionTypeProto_ != nil && protoDefault_UnionTypeProto_.Core_ == def {
			protoDefSub_TypeProto_GetUnion = protoDefault_UnionTypeProto_
		} else {
			w := &Ext_com_google_javascript_jscomp_serialization_UnionTypeProto{}
			w.Core_ = def
			protoDefSub_TypeProto_GetUnion = w
		}
	}
	if def := (&protobuf.TypedAst{}).GetTypePool(); def != nil {
		if protoDefault_TypePool_ != nil && protoDefault_TypePool_.Core_ == def {
			protoDefSub_TypedAst_GetTypePool = protoDefault_TypePool_
		} else {
			w := &Ext_com_google_javascript_jscomp_serialization_TypePool{}
			w.Core_ = def
			protoDefSub_TypedAst_GetTypePool = w
		}
	}
	if def := (&protobuf.TypedAst{}).GetStringPool(); def != nil {
		if protoDefault_StringPoolProto_ != nil && protoDefault_StringPoolProto_.Core_ == def {
			protoDefSub_TypedAst_GetStringPool = protoDefault_StringPoolProto_
		} else {
			w := &Ext_com_google_javascript_jscomp_serialization_StringPoolProto{}
			w.Core_ = def
			protoDefSub_TypedAst_GetStringPool = w
		}
	}
	if def := (&protobuf.TypedAst{}).GetSourceFilePool(); def != nil {
		if protoDefault_SourceFilePool_ != nil && protoDefault_SourceFilePool_.Core_ == def {
			protoDefSub_TypedAst_GetSourceFilePool = protoDefault_SourceFilePool_
		} else {
			w := &Ext_com_google_javascript_jscomp_serialization_SourceFilePool{}
			w.Core_ = def
			protoDefSub_TypedAst_GetSourceFilePool = w
		}
	}
	if def := (&protobuf.TypedAst{}).GetExternsSummary(); def != nil {
		if protoDefault_ExternsSummary_ != nil && protoDefault_ExternsSummary_.Core_ == def {
			protoDefSub_TypedAst_GetExternsSummary = protoDefault_ExternsSummary_
		} else {
			w := &Ext_com_google_javascript_jscomp_serialization_ExternsSummary{}
			w.Core_ = def
			protoDefSub_TypedAst_GetExternsSummary = w
		}
	}
}
