// Clean-room runtime support bridge for j2go mechanical Java-to-Go output.
package j2gort

import (
	"fmt"
	"hash/maphash"
	"math"
	"math/big"
	"math/bits"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
	"unique"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaregex"
	"github.com/paullepoulpe/jsgo/transpiled/shims/nodestore"
	shimsres "github.com/paullepoulpe/jsgo/transpiled/shims/resources"
)

var (
	_ = javaio.System_Out
	_ = javaregex.StringSplit
	_ = shimsres.ResourceExists
	_ = nodestore.NilNode
	_ = big.NewInt
	_ = regexp.Compile
	_ = sort.SliceStable
	_ = atomic.LoadUint32
	_ = unicode.IsDigit
	_ = utf16.DecodeRune
	_ = utf8.DecodeRune
	_ = unique.Make[string]
)

func LenientFormat(template string, args ...any) string {
	if template == NullString {
		template = "null"
	}
	var sb strings.Builder
	i, argIdx := 0, 0
	for i < len(template) {
		p := strings.Index(template[i:], "%s")
		if p < 0 || argIdx >= len(args) {
			sb.WriteString(template[i:])
			i = len(template)
			break
		}
		sb.WriteString(template[i : i+p])
		sb.WriteString(StrVal(args[argIdx]))
		argIdx++
		i += p + 2
	}
	if argIdx < len(args) {
		sb.WriteString(" [")
		for j := argIdx; j < len(args); j++ {
			if j > argIdx {
				sb.WriteString(", ")
			}
			sb.WriteString(StrVal(args[j]))
		}
		sb.WriteByte(']')
	}
	return sb.String()
}

func normFmtArg(v any) any {
	if v == nil {
		return nil
	}
	if s, ok := v.(string); ok {
		if s == NullString {
			return nil
		}
		return s
	}
	if rn, ok := v.(jRawNodeProvider); ok {
		if rn.RawNodeId_() == 0 {
			return nil
		}
		return v
	}
	return ToAny(v)
}

//go:noinline
func checkStateFmtFail(tmpl string, args ...any) {
	for i, a := range args {
		args[i] = normFmtArg(a)
	}
	panic(javalang.NewIllegalStateException(LenientFormat(tmpl, args...)))
}

//go:noinline
func checkArgFmtFail(tmpl string, args ...any) {
	for i, a := range args {
		args[i] = normFmtArg(a)
	}
	panic(javalang.NewIllegalArgumentException(LenientFormat(tmpl, args...)))
}

func CheckState(b bool) {
	if !b {
		panic(javalang.NewIllegalStateException(""))
	}
}
func CheckStateStr(b bool, msg string) {
	if !b {
		panic(javalang.NewIllegalStateException(StrValString(msg)))
	}
}
func CheckStateMsg(b bool, msg any) {
	if !b {
		panic(javalang.NewIllegalStateException(StrVal(msg)))
	}
}
func CheckStateFmt(b bool, tmpl string, args ...any) {
	if !b {
		panic(javalang.NewIllegalStateException(LenientFormat(tmpl, args...)))
	}
}
func CheckStateFmt1[A any](b bool, tmpl string, a0 A) {
	if !b {
		checkStateFmtFail(tmpl, a0)
	}
}
func CheckStateFmt2[A, B any](b bool, tmpl string, a0 A, a1 B) {
	if !b {
		checkStateFmtFail(tmpl, a0, a1)
	}
}
func CheckStateFmt3[A, B, C any](b bool, tmpl string, a0 A, a1 B, a2 C) {
	if !b {
		checkStateFmtFail(tmpl, a0, a1, a2)
	}
}
func CheckStateFmt4[A, B, C, D any](b bool, tmpl string, a0 A, a1 B, a2 C, a3 D) {
	if !b {
		checkStateFmtFail(tmpl, a0, a1, a2, a3)
	}
}
func CheckArgument(b bool) {
	if !b {
		panic(javalang.NewIllegalArgumentException(""))
	}
}
func CheckArgumentStr(b bool, msg string) {
	if !b {
		panic(javalang.NewIllegalArgumentException(StrValString(msg)))
	}
}
func CheckArgumentMsg(b bool, msg any) {
	if !b {
		panic(javalang.NewIllegalArgumentException(StrVal(msg)))
	}
}
func CheckArgumentFmt(b bool, tmpl string, args ...any) {
	if !b {
		panic(javalang.NewIllegalArgumentException(LenientFormat(tmpl, args...)))
	}
}
func CheckArgumentFmt1[A any](b bool, tmpl string, a0 A) {
	if !b {
		checkArgFmtFail(tmpl, a0)
	}
}
func CheckArgumentFmt2[A, B any](b bool, tmpl string, a0 A, a1 B) {
	if !b {
		checkArgFmtFail(tmpl, a0, a1)
	}
}
func CheckArgumentFmt3[A, B, C any](b bool, tmpl string, a0 A, a1 B, a2 C) {
	if !b {
		checkArgFmtFail(tmpl, a0, a1, a2)
	}
}
func CheckArgumentFmt4[A, B, C, D any](b bool, tmpl string, a0 A, a1 B, a2 C, a3 D) {
	if !b {
		checkArgFmtFail(tmpl, a0, a1, a2, a3)
	}
}
func CheckElementIndex(index, size int32) int32 {
	if index < 0 || index >= size {
		panic(javalang.NewIndexOutOfBoundsExceptionForIndex(int(index), int(size)))
	}
	return index
}
func CheckPositionIndex(index, size int32) int32 {
	if index < 0 || index > size {
		panic(javalang.NewIndexOutOfBoundsExceptionForIndex(int(index), int(size)))
	}
	return index
}
func CheckPositionIndexes(start, end, size int32) {
	if start < 0 || end < start || end > size {
		panic(javalang.NewIndexOutOfBoundsException(fmt.Sprintf("start=%d end=%d size=%d", start, end, size)))
	}
}
func ArraysFill[T any](a []T, val T) {
	if a == nil {
		ThrowNPE()
	}
	for i := range a {
		a[i] = val
	}
}
func ArraysFillZero[T any](a []T) {
	if a == nil {
		ThrowNPE()
	}
	clear(a)
}
func ArraysFillRange[T any](a []T, from, to int32, val T) {
	if a == nil {
		ThrowNPE()
	}
	CheckPositionIndexes(from, to, int32(len(a)))
	for i := from; i < to; i++ {
		a[i] = val
	}
}
func ArraysFillRangeZero[T any](a []T, from, to int32) {
	if a == nil {
		ThrowNPE()
	}
	CheckPositionIndexes(from, to, int32(len(a)))
	clear(a[from:to])
}

// NullString is a WTF-8/UTF-8 impossible byte sentinel representing Java `(String) null`
// inside a value-typed Go `string` with zero heap allocation.
const NullString = "\xff\x00j2go_null"

func IsNullString(s string) bool {
	return s == NullString
}

var fastStrHashSeed = maphash.MakeSeed()

func FastStrHash32(s string) int32 {
	h := uint32(maphash.String(fastStrHashSeed, s))
	if h == 0 {
		return 1
	}
	return int32(h)
}

type jJavaNilChecker interface {
	IsJavaNil_() bool
}

func isJavaNilAny(v any) bool {
	if v == nil {
		return true
	}
	if nc, ok := v.(jJavaNilChecker); ok {
		return nc.IsJavaNil_()
	}
	switch x := v.(type) {
	case string:
		return x == NullString
	case nodestore.Node:
		return x == 0
	case []any:
		return x == nil
	case []string:
		return x == nil
	case []byte:
		return x == nil
	case []int8:
		return x == nil
	case []int32:
		return x == nil
	case []uint16:
		return x == nil
	case bool, int32, int64, uint16, int8, int16, float32, float64:
		return false
	}
	return javalang.IsNil(v)
}

var asciiBoxedStrings_ [128]any
var asciiInternStrings_ [128]string

func init() {
	for i := 0; i < 128; i++ {
		s := string(byte(i))
		asciiInternStrings_[i] = s
		asciiBoxedStrings_[i] = s
	}
}

func BoxString(s string) any {
	if len(s) == 1 && s[0] < 128 {
		return asciiBoxedStrings_[s[0]]
	}
	if s == NullString {
		return nil
	}
	return s
}

func CastString(v any) string {
	if v == nil {
		return NullString
	}
	if s, ok := v.(string); ok {
		return s
	}
	if isJavaNilAny(v) {
		return NullString
	}
	panic(javalang.NewClassCastException(fmt.Sprintf("cannot cast %T to java.lang.String", v)))
}

type jRawNodeProvider interface {
	RawNodeId_() uint32
}

var globalNodeBoxer func(uint32) any

func RegisterNodeBoxer(fn func(uint32) any) {
	globalNodeBoxer = fn
}

func RegisterNodeChunkFiller(_ func(uint32, *[4096]any)) {}

func BoxNode[N ~uint32](n N) any {
	if n == 0 {
		return nil
	}
	return n
}

func boxNodeById(u uint32) any {
	if u == 0 {
		return nil
	}
	if globalNodeBoxer != nil {
		return globalNodeBoxer(u)
	}
	return nodestore.Node(u)
}

var globalColorIdBoxer func(int64) any

func RegisterColorIdBoxer(fn func(int64) any) {
	globalColorIdBoxer = fn
}

func boxColorId(raw int64) any {
	if globalColorIdBoxer != nil {
		return globalColorIdBoxer(raw)
	}
	return raw
}

type RawColorIdKey interface {
	comparable
	RawColorId_() (int64, bool)
}

type jRawColorIdProvider interface {
	RawColorId_() (int64, bool)
}

type FastEqValKey interface {
	comparable
	FastValHash_() (int32, bool, bool)
}

type jSingleNonNullElem interface {
	singleNonNullElem_() any
}

func CastNode[N ~uint32](v any) N {
	if v == nil {
		return 0
	}
	if n, ok := v.(N); ok {
		return n
	}
	if rn, ok := v.(jRawNodeProvider); ok {
		return N(rn.RawNodeId_())
	}
	if n2, ok := v.(nodestore.Node); ok {
		return N(n2)
	}
	if isJavaNilAny(v) {
		return 0
	}
	panic(javalang.NewClassCastException(fmt.Sprintf("cannot cast %T to Node", v)))
}

// IdentityComparable marks generated Java classes that inherit Object.equals and
// Object.hashCode (guaranteed comparable concrete struct pointers).
type IdentityComparable interface {
	HasIdentityEquals_()
}

// SliceIterable allows zero-allocation enhanced-for iteration over backing slices.
type SliceIterable interface {
	IterSlice_() []any
}

// StrSliceIterable allows zero-box enhanced-for iteration over unboxed string slices.
type StrSliceIterable interface {
	IterStrSlice_() ([]string, bool)
}

type FastStrList interface {
	ListGetStr_(int32) string
}

func IterSlice(v any) ([]any, bool) {
	if v == nil {
		return nil, false
	}
	if si, ok := v.(SliceIterable); ok {
		if nc, ok2 := v.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
			return nil, false
		}
		return si.IterSlice_(), true
	}
	return nil, false
}

func IfaceEquals[I interface{ IsJavaNil_() bool }](a, b I) bool {
	if any(a) == any(b) {
		return true
	}
	if any(a) == nil {
		return b.IsJavaNil_()
	}
	if any(b) == nil {
		return a.IsJavaNil_()
	}
	return a.IsJavaNil_() && b.IsJavaNil_()
}

func ObjEqStr(a, b string) bool { return a == b }

func ObjEqNode[N ~uint32](a, b N) bool { return a == b }

func ObjEqVal[T comparable](a, b T) bool { return a == b }

func ObjEqPtrIdent[T any](a, b *T) bool { return a == b }

func ObjEqCustomPtr[T any, P interface {
	*T
	Equals(any) bool
}](a, b P) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Equals(b)
}

func ObjEqIface[T interface {
	comparable
	IsJavaNil_() bool
	Equals(any) bool
}](a, b T) bool {
	if a == b {
		return true
	}
	var z T
	aNil := a == z || a.IsJavaNil_()
	bNil := b == z || b.IsJavaNil_()
	if aNil || bNil {
		return aNil && bNil
	}
	return a.Equals(b)
}

func ObjEqNilIface[T interface {
	comparable
	IsJavaNil_() bool
}](a, b T) bool {
	if a == b {
		return true
	}
	var z T
	aNil := a == z || a.IsJavaNil_()
	bNil := b == z || b.IsJavaNil_()
	if aNil || bNil {
		return aNil && bNil
	}
	if eq, ok := any(a).(jEqualer); ok {
		return eq.Equals(b)
	}
	return false
}

func CheckNotNullIface[I interface{ IsJavaNil_() bool }](v I) I {
	if any(v) == nil || v.IsJavaNil_() {
		ThrowNPE()
	}
	return v
}

// AsIface converts a concrete pointer or interface to target interface I while turning
// typed nil pointers into untyped Go `nil` (preventing Go's typed-nil interface trap).
func AsIface[I any](v any) I {
	var zero I
	if v == nil {
		return zero
	}
	if res, ok := v.(I); ok {
		if nc, ok2 := v.(jJavaNilChecker); ok2 {
			if nc.IsJavaNil_() {
				return zero
			}
			return res
		}
		if isJavaNilAny(v) {
			return zero
		}
		return res
	}
	if isJavaNilAny(v) {
		return zero
	}
	panic(javalang.NewClassCastException(fmt.Sprintf("cannot convert %T to %T", v, zero)))
}

func ToAny(v any) any {
	if v == nil {
		return nil
	}
	if nc, ok := v.(jJavaNilChecker); ok {
		if nc.IsJavaNil_() {
			return nil
		}
		return v
	}
	switch x := v.(type) {
	case string:
		if x == NullString {
			return nil
		}
		return v
	case nodestore.Node:
		if x == 0 {
			return nil
		}
		return v
	case bool, int32, int64, uint16, int8, int16, float32, float64:
		return v
	case []any:
		if x == nil {
			return nil
		}
		return v
	case []string:
		if x == nil {
			return nil
		}
		return v
	case []byte:
		if x == nil {
			return nil
		}
		return v
	case []int8:
		if x == nil {
			return nil
		}
		return v
	case []int32:
		if x == nil {
			return nil
		}
		return v
	case []uint16:
		if x == nil {
			return nil
		}
		return v
	}
	if javalang.IsNil(v) {
		return nil
	}
	return v
}

func PtrToAny[T any](p *T) any {
	if p == nil {
		return nil
	}
	return p
}

func SliceToAny[E any](s []E) any {
	if s == nil {
		return nil
	}
	return s
}

func BoxNullableVal[T comparable](v T) any {
	var zero T
	if v == zero {
		return nil
	}
	return v
}

// Cast performs a Java checked reference cast `(T) v`, returning zero for nil and
// throwing ClassCastException on mismatch.
func Cast[T any](v any) T {
	var zero T
	if v == nil {
		return zero
	}
	if res, ok := v.(T); ok {
		if nc, ok2 := v.(jJavaNilChecker); ok2 {
			if nc.IsJavaNil_() {
				return zero
			}
			return res
		}
		if isJavaNilAny(v) {
			return zero
		}
		return res
	}
	if isJavaNilAny(v) {
		return zero
	}
	panic(javalang.NewClassCastException(fmt.Sprintf("cannot cast %T to %T", v, zero)))
}

func CastPtr[T any](v any) *T {
	if v == nil {
		return nil
	}
	if res, ok := v.(*T); ok {
		return res
	}
	if isJavaNilAny(v) {
		return nil
	}
	var zero *T
	panic(javalang.NewClassCastException(fmt.Sprintf("cannot cast %T to %T", v, zero)))
}

func UncheckedCast[T any](v any) T {
	var zero T
	if v == nil {
		return zero
	}
	if res, ok := v.(T); ok {
		if nc, ok2 := v.(jJavaNilChecker); ok2 {
			if nc.IsJavaNil_() {
				return zero
			}
			return res
		}
		if isJavaNilAny(v) {
			return zero
		}
		return res
	}
	return zero
}

func TryCast[T any](v any) (T, bool) {
	var zero T
	if v == nil {
		return zero, false
	}
	res, ok := v.(T)
	if !ok {
		return zero, false
	}
	if nc, ok2 := v.(jJavaNilChecker); ok2 {
		if nc.IsJavaNil_() {
			return zero, false
		}
		return res, true
	}
	if isJavaNilAny(v) {
		return zero, false
	}
	return res, true
}

func TryCastString(v any) (string, bool) {
	if v == nil {
		return NullString, false
	}
	if s, ok := v.(string); ok && s != NullString {
		return s, true
	}
	return NullString, false
}

func InstanceOf[T any](v any) bool {
	if v == nil {
		return false
	}
	if _, ok := v.(T); !ok {
		return false
	}
	if nc, ok2 := v.(jJavaNilChecker); ok2 {
		return !nc.IsJavaNil_()
	}
	return !isJavaNilAny(v)
}

func InstanceOfBind[T any](v any, out *T) bool {
	if v == nil {
		return false
	}
	t, ok := v.(T)
	if !ok {
		return false
	}
	if nc, ok2 := v.(jJavaNilChecker); ok2 {
		if nc.IsJavaNil_() {
			return false
		}
	} else if isJavaNilAny(v) {
		return false
	}
	*out = t
	return true
}

func InstanceOfString(v any) bool {
	if v == nil {
		return false
	}
	s, ok := v.(string)
	return ok && s != NullString
}

func InstanceOfStringBind(v any, out *string) bool {
	if v == nil {
		return false
	}
	if s, ok := v.(string); ok && s != NullString {
		*out = s
		return true
	}
	return false
}

// RefEquals implements Java reference `==` across interfaces, pointers, slices, and boxed values.
func RefEquals(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return isJavaNilAny(b)
	}
	if b == nil {
		return isJavaNilAny(a)
	}
	if ba, ok := a.(Ext_com_google_protobuf_ByteString); ok {
		if bb, ok2 := b.(Ext_com_google_protobuf_ByteString); ok2 {
			return ByteStringRefEquals(ba, bb)
		}
		return !ba.valid && isJavaNilAny(b)
	}
	if bb, ok := b.(Ext_com_google_protobuf_ByteString); ok {
		return !bb.valid && isJavaNilAny(a)
	}
	ncA, okA := a.(jJavaNilChecker)
	ncB, okB := b.(jJavaNilChecker)
	if okA && okB {
		na := ncA.IsJavaNil_()
		nb := ncB.IsJavaNil_()
		if na || nb {
			return na && nb
		}
		ta := reflect.TypeOf(a)
		tb := reflect.TypeOf(b)
		if ta == tb {
			if ta.Comparable() {
				return a == b
			}
			return false
		}
		if ta.Kind() == reflect.Ptr && tb.Kind() == reflect.Ptr {
			return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
		}
		return false
	}
	na := isJavaNilAny(a)
	nb := isJavaNilAny(b)
	if na || nb {
		return na && nb
	}
	switch xa := a.(type) {
	case string:
		xb, ok := b.(string)
		return ok && xa == xb
	case nodestore.Node:
		xb, ok := b.(nodestore.Node)
		return ok && xa == xb
	}
	ta := reflect.TypeOf(a)
	tb := reflect.TypeOf(b)
	if ta == tb {
		if ta.Kind() == reflect.Slice {
			ra := reflect.ValueOf(a)
			rb := reflect.ValueOf(b)
			return ra.Pointer() == rb.Pointer() && ra.Len() == rb.Len()
		}
		if ta.Comparable() {
			return a == b
		}
		return false
	}
	if ta.Kind() == reflect.Ptr && tb.Kind() == reflect.Ptr {
		return reflect.ValueOf(a).Pointer() == reflect.ValueOf(b).Pointer()
	}
	return false
}

func SliceRefEquals[T any](a, b []T) bool {
	if (a == nil) != (b == nil) || len(a) != len(b) || cap(a) != cap(b) {
		return false
	}
	if a == nil {
		return true
	}
	if cap(a) > 0 {
		return &a[:1][0] == &b[:1][0]
	}
	return true
}

func ThrowNPE() {
	panic(javalang.NewNullPointerException())
}

func CheckNotNullString(s string) string {
	if s == NullString {
		ThrowNPE()
	}
	return s
}

func CheckNotNullNode[N ~uint32](n N) N {
	if n == 0 {
		ThrowNPE()
	}
	return n
}

func CheckNotNullEnum[E ~uint8](v E) E {
	if v == 0 {
		ThrowNPE()
	}
	return v
}

func CheckNotNullVal[T comparable](v T) T {
	var zero T
	if v == zero {
		ThrowNPE()
	}
	return v
}

func CheckNotNullPtr[T any](p *T) *T {
	if p == nil {
		ThrowNPE()
	}
	return p
}

func CheckNotNullSlice[T any](s []T) []T {
	if s == nil {
		ThrowNPE()
	}
	return s
}

func CheckNotNull[T any](v T) T {
	if any(v) == nil {
		ThrowNPE()
	}
	if nc, ok := any(v).(jJavaNilChecker); ok {
		if nc.IsJavaNil_() {
			ThrowNPE()
		}
		return v
	}
	if isJavaNilAny(v) {
		ThrowNPE()
	}
	return v
}

func EnumSwitchOrd[E interface {
	comparable
	Ordinal() int32
}](e E) int32 {
	var zero E
	if e == zero {
		return -1
	}
	return e.Ordinal()
}

func BoolToStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func StrValString(s string) string {
	if s == NullString {
		return "null"
	}
	return s
}

func StrValInt32(v int32) string {
	return strconv.FormatInt(int64(v), 10)
}

func StrValInt64(v int64) string {
	return strconv.FormatInt(v, 10)
}

func StrValBool(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func StrValChar(v uint16) string {
	if v < 0x80 {
		return string(rune(v))
	}
	return javalang.StringFromChar(v)
}

func StrVal(v any) string {
	if v == nil {
		return "null"
	}
	switch x := v.(type) {
	case string:
		if x == NullString {
			return "null"
		}
		return x
	case float64:
		return javalang.DoubleToString(x)
	case float32:
		return javalang.FloatToString(x)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case bool:
		if x {
			return "true"
		}
		return "false"
	case uint16:
		return javalang.StringFromChar(x)
	case int8:
		return strconv.FormatInt(int64(x), 10)
	case int16:
		return strconv.FormatInt(int64(x), 10)
	case int:
		return strconv.Itoa(x)
	case *BoxDouble:
		if x == nil {
			return "null"
		}
		return javalang.DoubleToString(x.Val)
	case *BoxFloat:
		if x == nil {
			return "null"
		}
		return javalang.FloatToString(x.Val)
	case *BoxInt:
		if x == nil {
			return "null"
		}
		return strconv.FormatInt(int64(x.Val), 10)
	case *BoxLong:
		if x == nil {
			return "null"
		}
		return strconv.FormatInt(x.Val, 10)
	case *BoxBoolean:
		if x == nil {
			return "null"
		}
		if x.Val {
			return "true"
		}
		return "false"
	case *BoxChar:
		if x == nil {
			return "null"
		}
		return javalang.StringFromChar(x.Val)
	case *BoxByte:
		if x == nil {
			return "null"
		}
		return strconv.FormatInt(int64(x.Val), 10)
	case *BoxShort:
		if x == nil {
			return "null"
		}
		return strconv.FormatInt(int64(x.Val), 10)
	}
	if isJavaNilAny(v) {
		return "null"
	}
	return javalang.ToJavaString(v)
}

func StrConcat(parts ...any) string {
	n := len(parts)
	if n == 0 {
		return ""
	}
	if n == 1 {
		return StrVal(parts[0])
	}
	var stackBuf [16]string
	var strs []string
	if n <= len(stackBuf) {
		strs = stackBuf[:n]
	} else {
		strs = make([]string, n)
	}
	total := 0
	for i, p := range parts {
		s := StrVal(p)
		strs[i] = s
		total += len(s)
	}
	if total == 0 {
		return ""
	}
	var sb strings.Builder
	sb.Grow(total)
	for _, s := range strs {
		sb.WriteString(s)
	}
	return sb.String()
}

func NewStringArray(n int32) []string {
	if n < 0 {
		panic(javalang.NewRuntimeExceptionMsg("NegativeArraySizeException"))
	}
	s := make([]string, n)
	for i := range s {
		s[i] = NullString
	}
	return s
}

func Make2D[T any](d1, d2 int32, initElem func() T) [][]T {
	out := make([][]T, d1)
	for i := int32(0); i < d1; i++ {
		row := make([]T, d2)
		if initElem != nil {
			for j := int32(0); j < d2; j++ {
				row[j] = initElem()
			}
		}
		out[i] = row
	}
	return out
}

func CloneSlice[T any](s []T) []T {
	if s == nil {
		return nil
	}
	out := make([]T, len(s))
	copy(out, s)
	return out
}

// --- Goroutine ID & Reentrant Monitor for Java `synchronized` ---

func currentGoroutineID() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// Format: "goroutine <id> [..."
	const prefix = "goroutine "
	if n <= len(prefix) {
		return 1
	}
	var id uint64
	for i := len(prefix); i < n; i++ {
		b := buf[i]
		if b < '0' || b > '9' {
			break
		}
		id = id*10 + uint64(b-'0')
	}
	return id
}

type reentrantMonitor struct {
	stateMu sync.Mutex
	lockCh  chan struct{}
	owner   uint64
	depth   int
}

func newReentrantMonitor() *reentrantMonitor {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	return &reentrantMonitor{lockCh: ch}
}

func (rm *reentrantMonitor) lock() {
	gid := currentGoroutineID()
	rm.stateMu.Lock()
	if rm.depth > 0 && rm.owner == gid {
		rm.depth++
		rm.stateMu.Unlock()
		return
	}
	rm.stateMu.Unlock()
	<-rm.lockCh
	rm.stateMu.Lock()
	rm.owner = gid
	rm.depth = 1
	rm.stateMu.Unlock()
}

func (rm *reentrantMonitor) unlock() {
	rm.stateMu.Lock()
	rm.depth--
	if rm.depth == 0 {
		rm.owner = 0
		rm.stateMu.Unlock()
		rm.lockCh <- struct{}{}
		return
	}
	rm.stateMu.Unlock()
}

var monitorMap sync.Map // key -> *reentrantMonitor

func LockMonitor(obj any) func() {
	if obj == nil || javalang.IsNil(obj) {
		panic(javalang.NewNullPointerException())
	}
	rv := reflect.ValueOf(obj)
	var key any = obj
	if rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Slice {
		key = rv.Pointer()
	}
	m, ok := monitorMap.Load(key)
	if !ok {
		m, _ = monitorMap.LoadOrStore(key, newReentrantMonitor())
	}
	rm := m.(*reentrantMonitor)
	rm.lock()
	return rm.unlock
}

func CloseResource(res any) {
	if res == nil || javalang.IsNil(res) {
		return
	}
	if c, ok := res.(interface{ Close() }); ok {
		c.Close()
	} else if c2, ok := res.(interface{ Close_0() }); ok {
		c2.Close_0()
	} else if c3, ok := res.(interface{ Close() error }); ok {
		_ = c3.Close()
	}
}

type ClinitState struct {
	mu    sync.Mutex
	owner uint64
	wait  chan struct{}
}

func EnterClinit(st *ClinitState, done *uint32) bool {
	gid := currentGoroutineID()
	st.mu.Lock()
	for {
		if atomic.LoadUint32(done) != 0 {
			st.mu.Unlock()
			return false
		}
		if st.owner == gid {
			st.mu.Unlock()
			return false
		}
		if st.owner == 0 {
			st.owner = gid
			st.wait = make(chan struct{})
			st.mu.Unlock()
			return true
		}
		ch := st.wait
		st.mu.Unlock()
		<-ch
		st.mu.Lock()
	}
}

func ExitClinit(st *ClinitState, done *uint32) {
	st.mu.Lock()
	atomic.StoreUint32(done, 1)
	ch := st.wait
	st.owner = 0
	st.wait = nil
	st.mu.Unlock()
	if ch != nil {
		close(ch)
	}
}

// --- Boxed Primitives & Number ---

type Number interface {
	IntValue() int32
	LongValue() int64
	FloatValue() float32
	DoubleValue() float64
	ByteValue() int8
	ShortValue() int16
	Implements_j2gort_Number()
}

func ToIface_Number[T interface {
	Number
	comparable
}](v T) Number {
	var z T
	if v == z {
		return nil
	}
	return v
}

type NumberBase struct{}

func (n *NumberBase) IsJavaNil_() bool        { return n == nil }
func (*NumberBase) Implements_j2gort_Number() {}
func (*NumberBase) IntValue() int32           { return 0 }
func (*NumberBase) LongValue() int64          { return 0 }
func (*NumberBase) FloatValue() float32       { return 0 }
func (*NumberBase) DoubleValue() float64      { return 0 }
func (*NumberBase) ByteValue() int8           { return 0 }
func (*NumberBase) ShortValue() int16         { return 0 }
func (*NumberBase) BindExtSelf_(any)          {}
func InitExt_java_lang_Number(*NumberBase)    {}

type BoxBoolean struct{ Val bool }

func (b *BoxBoolean) IsJavaNil_() bool   { return b == nil }
func (b *BoxBoolean) BooleanValue() bool { return CheckNotNullPtr(b).Val }
func (b *BoxBoolean) Equals(o any) bool {
	x, ok := o.(*BoxBoolean)
	return ok && x != nil && x.Val == b.Val
}
func (b *BoxBoolean) HashCode() int32 {
	if b.Val {
		return 1231
	}
	return 1237
}
func (b *BoxBoolean) String() string {
	if b.Val {
		return "true"
	}
	return "false"
}
func (*BoxBoolean) JavaClassInfo_() (string, string) { return "java.lang.Boolean", "Boolean" }

var (
	Boolean_TRUE  = &BoxBoolean{Val: true}
	Boolean_FALSE = &BoxBoolean{Val: false}
)

func ValueOfBoolean(v bool) *BoxBoolean {
	if v {
		return Boolean_TRUE
	}
	return Boolean_FALSE
}
func UnboxBoolean(b *BoxBoolean) bool { return CheckNotNullPtr(b).Val }

type BoxByte struct{ Val int8 }

func (b *BoxByte) IsJavaNil_() bool               { return b == nil }
func (b *BoxByte) Implements_j2gort_Number()      {}
func (b *BoxByte) ByteValue() int8                { return CheckNotNullPtr(b).Val }
func (b *BoxByte) ShortValue() int16              { return int16(CheckNotNullPtr(b).Val) }
func (b *BoxByte) IntValue() int32                { return int32(CheckNotNullPtr(b).Val) }
func (b *BoxByte) LongValue() int64               { return int64(CheckNotNullPtr(b).Val) }
func (b *BoxByte) FloatValue() float32            { return float32(CheckNotNullPtr(b).Val) }
func (b *BoxByte) DoubleValue() float64           { return float64(CheckNotNullPtr(b).Val) }
func (b *BoxByte) Equals(o any) bool              { x, ok := o.(*BoxByte); return ok && x != nil && x.Val == b.Val }
func (b *BoxByte) HashCode() int32                { return int32(b.Val) }
func (b *BoxByte) String() string                 { return strconv.Itoa(int(b.Val)) }
func (*BoxByte) JavaClassInfo_() (string, string) { return "java.lang.Byte", "Byte" }

var (
	boxByteCache_  [256]BoxByte
	boxShortCache_ [256]BoxShort
	boxCharCache_  [128]BoxChar
	boxIntCache_   [65664]BoxInt // [-128..65535]
	boxLongCache_  [256]BoxLong
)

func init() {
	for i := 0; i < 256; i++ {
		boxByteCache_[i].Val = int8(i - 128)
		boxShortCache_[i].Val = int16(i - 128)
		boxLongCache_[i].Val = int64(i - 128)
	}
	for i := 0; i < 128; i++ {
		boxCharCache_[i].Val = uint16(i)
	}
	for i := 0; i < len(boxIntCache_); i++ {
		boxIntCache_[i].Val = int32(i - 128)
	}
}

func ValueOfByte(v int8) *BoxByte { return &boxByteCache_[int(v)+128] }
func UnboxByte(b *BoxByte) int8   { return CheckNotNullPtr(b).Val }

type BoxShort struct{ Val int16 }

func (b *BoxShort) IsJavaNil_() bool          { return b == nil }
func (b *BoxShort) Implements_j2gort_Number() {}
func (b *BoxShort) ByteValue() int8           { return int8(CheckNotNullPtr(b).Val) }
func (b *BoxShort) ShortValue() int16         { return CheckNotNullPtr(b).Val }
func (b *BoxShort) IntValue() int32           { return int32(CheckNotNullPtr(b).Val) }
func (b *BoxShort) LongValue() int64          { return int64(CheckNotNullPtr(b).Val) }
func (b *BoxShort) FloatValue() float32       { return float32(CheckNotNullPtr(b).Val) }
func (b *BoxShort) DoubleValue() float64      { return float64(CheckNotNullPtr(b).Val) }
func (b *BoxShort) Equals(o any) bool {
	x, ok := o.(*BoxShort)
	return ok && x != nil && x.Val == b.Val
}
func (b *BoxShort) HashCode() int32                { return int32(b.Val) }
func (b *BoxShort) String() string                 { return strconv.Itoa(int(b.Val)) }
func (*BoxShort) JavaClassInfo_() (string, string) { return "java.lang.Short", "Short" }
func ValueOfShort(v int16) *BoxShort {
	if v >= -128 && v <= 127 {
		return &boxShortCache_[int(v)+128]
	}
	return &BoxShort{Val: v}
}
func UnboxShort(b *BoxShort) int16 { return CheckNotNullPtr(b).Val }

type BoxChar struct{ Val uint16 }

func (b *BoxChar) IsJavaNil_() bool               { return b == nil }
func (b *BoxChar) CharValue() uint16              { return CheckNotNullPtr(b).Val }
func (b *BoxChar) Equals(o any) bool              { x, ok := o.(*BoxChar); return ok && x != nil && x.Val == b.Val }
func (b *BoxChar) HashCode() int32                { return int32(b.Val) }
func (b *BoxChar) String() string                 { return javalang.StringFromChar(b.Val) }
func (*BoxChar) JavaClassInfo_() (string, string) { return "java.lang.Character", "Character" }
func ValueOfChar(v uint16) *BoxChar {
	if v < 128 {
		return &boxCharCache_[v]
	}
	return &BoxChar{Val: v}
}
func UnboxChar(b *BoxChar) uint16 { return CheckNotNullPtr(b).Val }

type BoxInt struct{ Val int32 }

func (b *BoxInt) IsJavaNil_() bool               { return b == nil }
func (b *BoxInt) Implements_j2gort_Number()      {}
func (b *BoxInt) ByteValue() int8                { return int8(CheckNotNullPtr(b).Val) }
func (b *BoxInt) ShortValue() int16              { return int16(CheckNotNullPtr(b).Val) }
func (b *BoxInt) IntValue() int32                { return CheckNotNullPtr(b).Val }
func (b *BoxInt) LongValue() int64               { return int64(CheckNotNullPtr(b).Val) }
func (b *BoxInt) FloatValue() float32            { return float32(CheckNotNullPtr(b).Val) }
func (b *BoxInt) DoubleValue() float64           { return float64(CheckNotNullPtr(b).Val) }
func (b *BoxInt) Equals(o any) bool              { x, ok := o.(*BoxInt); return ok && x != nil && x.Val == b.Val }
func (b *BoxInt) HashCode() int32                { return b.Val }
func (b *BoxInt) String() string                 { return strconv.FormatInt(int64(b.Val), 10) }
func (*BoxInt) JavaClassInfo_() (string, string) { return "java.lang.Integer", "Integer" }
func ValueOfInt(v int32) *BoxInt {
	if u := uint32(v) + 128; u < uint32(len(boxIntCache_)) {
		return &boxIntCache_[u]
	}
	return &BoxInt{Val: v}
}
func UnboxInt(b *BoxInt) int32 { return CheckNotNullPtr(b).Val }

type BoxLong struct{ Val int64 }

func (b *BoxLong) IsJavaNil_() bool               { return b == nil }
func (b *BoxLong) Implements_j2gort_Number()      {}
func (b *BoxLong) ByteValue() int8                { return int8(CheckNotNullPtr(b).Val) }
func (b *BoxLong) ShortValue() int16              { return int16(CheckNotNullPtr(b).Val) }
func (b *BoxLong) IntValue() int32                { return int32(CheckNotNullPtr(b).Val) }
func (b *BoxLong) LongValue() int64               { return CheckNotNullPtr(b).Val }
func (b *BoxLong) FloatValue() float32            { return float32(CheckNotNullPtr(b).Val) }
func (b *BoxLong) DoubleValue() float64           { return float64(CheckNotNullPtr(b).Val) }
func (b *BoxLong) Equals(o any) bool              { x, ok := o.(*BoxLong); return ok && x != nil && x.Val == b.Val }
func (b *BoxLong) HashCode() int32                { return int32(uint64(b.Val) ^ (uint64(b.Val) >> 32)) }
func (b *BoxLong) String() string                 { return strconv.FormatInt(b.Val, 10) }
func (*BoxLong) JavaClassInfo_() (string, string) { return "java.lang.Long", "Long" }
func ValueOfLong(v int64) *BoxLong {
	if v >= -128 && v <= 127 {
		return &boxLongCache_[int(v)+128]
	}
	return &BoxLong{Val: v}
}
func UnboxLong(b *BoxLong) int64 { return CheckNotNullPtr(b).Val }

type BoxFloat struct{ Val float32 }

func (b *BoxFloat) IsJavaNil_() bool          { return b == nil }
func (b *BoxFloat) Implements_j2gort_Number() {}
func (b *BoxFloat) ByteValue() int8           { return int8(javalang.D2I(float64(CheckNotNullPtr(b).Val))) }
func (b *BoxFloat) ShortValue() int16         { return int16(javalang.D2I(float64(CheckNotNullPtr(b).Val))) }
func (b *BoxFloat) IntValue() int32           { return javalang.D2I(float64(CheckNotNullPtr(b).Val)) }
func (b *BoxFloat) LongValue() int64          { return javalang.D2L(float64(CheckNotNullPtr(b).Val)) }
func (b *BoxFloat) FloatValue() float32       { return CheckNotNullPtr(b).Val }
func (b *BoxFloat) DoubleValue() float64      { return float64(CheckNotNullPtr(b).Val) }
func (b *BoxFloat) Equals(o any) bool {
	x, ok := o.(*BoxFloat)
	if !ok || x == nil {
		return false
	}
	if math.IsNaN(float64(b.Val)) && math.IsNaN(float64(x.Val)) {
		return true
	}
	return math.Float32bits(b.Val) == math.Float32bits(x.Val)
}
func (b *BoxFloat) HashCode() int32 {
	if math.IsNaN(float64(b.Val)) {
		return 0x7fc00000
	}
	return int32(math.Float32bits(b.Val))
}
func (b *BoxFloat) String() string                 { return javalang.FloatToString(b.Val) }
func (b *BoxFloat) IsNaN() bool                    { return math.IsNaN(float64(CheckNotNullPtr(b).Val)) }
func (b *BoxFloat) IsInfinite() bool               { return math.IsInf(float64(CheckNotNullPtr(b).Val), 0) }
func (*BoxFloat) JavaClassInfo_() (string, string) { return "java.lang.Float", "Float" }
func ValueOfFloat(v float32) *BoxFloat             { return &BoxFloat{Val: v} }
func UnboxFloat(b *BoxFloat) float32               { return CheckNotNullPtr(b).Val }

type BoxDouble struct{ Val float64 }

func (b *BoxDouble) IsJavaNil_() bool          { return b == nil }
func (b *BoxDouble) Implements_j2gort_Number() {}
func (b *BoxDouble) ByteValue() int8           { return int8(javalang.D2I(CheckNotNullPtr(b).Val)) }
func (b *BoxDouble) ShortValue() int16         { return int16(javalang.D2I(CheckNotNullPtr(b).Val)) }
func (b *BoxDouble) IntValue() int32           { return javalang.D2I(CheckNotNullPtr(b).Val) }
func (b *BoxDouble) LongValue() int64          { return javalang.D2L(CheckNotNullPtr(b).Val) }
func (b *BoxDouble) FloatValue() float32       { return float32(CheckNotNullPtr(b).Val) }
func (b *BoxDouble) DoubleValue() float64      { return CheckNotNullPtr(b).Val }
func (b *BoxDouble) Equals(o any) bool {
	x, ok := o.(*BoxDouble)
	if !ok || x == nil {
		return false
	}
	if math.IsNaN(b.Val) && math.IsNaN(x.Val) {
		return true
	}
	return math.Float64bits(b.Val) == math.Float64bits(x.Val)
}
func (b *BoxDouble) HashCode() int32 {
	if math.IsNaN(b.Val) {
		return 0x7ff80000
	}
	u := math.Float64bits(b.Val)
	return int32(u ^ (u >> 32))
}
func (b *BoxDouble) String() string                 { return javalang.DoubleToString(b.Val) }
func (b *BoxDouble) IsNaN() bool                    { return math.IsNaN(CheckNotNullPtr(b).Val) }
func (b *BoxDouble) IsInfinite() bool               { return math.IsInf(CheckNotNullPtr(b).Val, 0) }
func (*BoxDouble) JavaClassInfo_() (string, string) { return "java.lang.Double", "Double" }
func ValueOfDouble(v float64) *BoxDouble            { return &BoxDouble{Val: v} }
func UnboxDouble(b *BoxDouble) float64              { return CheckNotNullPtr(b).Val }

func UnboxToDouble(v any) float64 {
	if n, ok := v.(Number); ok && !isJavaNilAny(n) {
		return n.DoubleValue()
	}
	if c, ok := v.(*BoxChar); ok && c != nil {
		return float64(c.Val)
	}
	panic(javalang.NewNullPointerException())
}

// --- Enum & Type Predicate Registries ---

type EnumBase struct {
	Name_     string                         `gson:"-"`
	Ordinal_  int32                          `gson:"-"`
	equalPred atomic.Pointer[predEqualEnum_] `gson:"-"`
}

func (e *EnumBase) IsJavaNil_() bool { return e == nil }
func (e *EnumBase) Name() string     { return e.Name_ }
func (e *EnumBase) Ordinal() int32   { return e.Ordinal_ }
func (e *EnumBase) String() string   { return e.Name_ }
func (e *EnumBase) HashCode() int32 {
	return javalang.StringHashCode(e.Name_)*31 + e.Ordinal_
}
func (e *EnumBase) Equals(other any) bool {
	if isJavaNilAny(other) {
		return false
	}
	if oa, ok := other.(EnumAccessor); ok {
		return e == oa.GetEnumBase()
	}
	return false
}
func (e *EnumBase) CompareTo(other any) int32 {
	if isJavaNilAny(other) {
		panic(javalang.NewNullPointerException())
	}
	oa := other.(EnumAccessor).GetEnumBase()
	if e.Ordinal_ < oa.Ordinal_ {
		return -1
	}
	if e.Ordinal_ > oa.Ordinal_ {
		return 1
	}
	return 0
}

type EnumAccessor interface {
	GetEnumBase() *EnumBase
}

func (e *EnumBase) GetEnumBase() *EnumBase { return e }
func (e *EnumBase) EnumOrdinal_() int32 {
	if e == nil {
		return -1
	}
	return e.Ordinal_
}
func (*EnumBase) HasIdentityEquals_()                                 {}
func (*EnumBase) Implements_j2gort_Ext_java_lang_Enum()               {}
func (*EnumBase) Implements_j2gort_Ext_java_lang_constant_Constable() {}
func (e *EnumBase) GetExt_java_lang_EnumBase() *Ext_java_lang_EnumBase {
	return &Ext_java_lang_EnumBase{Core_: e}
}

func EnumName(v any) string {
	if isJavaNilAny(v) {
		panic(javalang.NewNullPointerException())
	}
	return v.(EnumAccessor).GetEnumBase().Name_
}

func EnumOrdinal(v any) int32 {
	if isJavaNilAny(v) {
		panic(javalang.NewNullPointerException())
	}
	return v.(EnumAccessor).GetEnumBase().Ordinal_
}

func EnumCompareTo(a, b any) int32 {
	return EnumOrdinal(a) - EnumOrdinal(b)
}

var (
	enumRegistry  sync.Map // string -> []any
	enumLoaderMap sync.Map // string -> func()
	typePredMap   sync.Map // string -> func(any) bool
)

func RegisterEnumConstants(className string, vals []any) {
	enumRegistry.Store(className, vals)
}

func RegisterEnumConstantsLoader(qn, binName string, loader func()) {
	enumLoaderMap.Store(qn, loader)
	if binName != qn {
		enumLoaderMap.Store(binName, loader)
	}
}

func LookupEnumConstantsByName(name string) []any {
	if name == "" {
		return nil
	}
	if v, ok := enumRegistry.Load(name); ok {
		return v.([]any)
	}
	if loader, ok := enumLoaderMap.Load(name); ok {
		loader.(func())()
	}
	if v, ok := enumRegistry.Load(name); ok {
		return v.([]any)
	}
	return nil
}

func LookupEnumConstants(cls *javalang.Class) []any {
	if cls == nil {
		return nil
	}
	return LookupEnumConstantsByName(cls.GetName())
}

func RegisterTypePredicate(qn, binName string, pred func(any) bool) {
	typePredMap.Store(qn, pred)
	if binName != qn {
		typePredMap.Store(binName, pred)
	}
}

func CheckTypePredicate(cls *javalang.Class, v any) bool {
	if cls == nil || v == nil || javalang.IsNil(v) {
		return false
	}
	if pred, ok := typePredMap.Load(cls.GetName()); ok {
		return pred.(func(any) bool)(v)
	}
	return GetObjectClass(v).GetName() == cls.GetName()
}

// --- Zero-Closure Guava & JDK Concrete Predicates ---

type jEqualTargetProvider interface {
	equalTarget_() any
}

type predAlwaysTrueType_ struct{ _ byte }

func (*predAlwaysTrueType_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*predAlwaysTrueType_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (*predAlwaysTrueType_) HasIdentityEquals_()                                     {}
func (p *predAlwaysTrueType_) IsJavaNil_() bool                                      { return p == nil }
func (p *predAlwaysTrueType_) Apply(_ any) bool                                      { CheckNotNullPtr(p); return true }
func (p *predAlwaysTrueType_) Test(_ any) bool                                       { CheckNotNullPtr(p); return true }
func (p *predAlwaysTrueType_) TestNode_(_ uint32) bool                               { CheckNotNullPtr(p); return true }
func (p *predAlwaysTrueType_) Equals(o any) bool                                     { CheckNotNullPtr(p); return p == o }
func (p *predAlwaysTrueType_) HashCode() int32                                       { CheckNotNullPtr(p); return 0x54525545 }
func (p *predAlwaysTrueType_) String() string                                        { CheckNotNullPtr(p); return "Predicates.alwaysTrue()" }
func (p *predAlwaysTrueType_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return CheckNotNull(o)
}
func (p *predAlwaysTrueType_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	CheckNotNull(o)
	return p
}
func (p *predAlwaysTrueType_) Negate() Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return predAlwaysFalseIface_
}

type predAlwaysFalseType_ struct{ _ byte }

func (*predAlwaysFalseType_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*predAlwaysFalseType_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (*predAlwaysFalseType_) HasIdentityEquals_()                                     {}
func (p *predAlwaysFalseType_) IsJavaNil_() bool                                      { return p == nil }
func (p *predAlwaysFalseType_) Apply(_ any) bool                                      { CheckNotNullPtr(p); return false }
func (p *predAlwaysFalseType_) Test(_ any) bool                                       { CheckNotNullPtr(p); return false }
func (p *predAlwaysFalseType_) TestNode_(_ uint32) bool                               { CheckNotNullPtr(p); return false }
func (p *predAlwaysFalseType_) Equals(o any) bool                                     { CheckNotNullPtr(p); return p == o }
func (p *predAlwaysFalseType_) HashCode() int32                                       { CheckNotNullPtr(p); return 0x46414c53 }
func (p *predAlwaysFalseType_) String() string                                        { CheckNotNullPtr(p); return "Predicates.alwaysFalse()" }
func (p *predAlwaysFalseType_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	CheckNotNull(o)
	return p
}
func (p *predAlwaysFalseType_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return CheckNotNull(o)
}
func (p *predAlwaysFalseType_) Negate() Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return predAlwaysTrueIface_
}

type predIsNullType_ struct{ _ byte }

func (*predIsNullType_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*predIsNullType_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (*predIsNullType_) HasIdentityEquals_()                                     {}
func (p *predIsNullType_) IsJavaNil_() bool                                      { return p == nil }
func (p *predIsNullType_) Apply(x any) bool                                      { CheckNotNullPtr(p); return x == nil || isJavaNilAny(x) }
func (p *predIsNullType_) Test(x any) bool                                       { CheckNotNullPtr(p); return x == nil || isJavaNilAny(x) }
func (p *predIsNullType_) TestNode_(n uint32) bool                               { CheckNotNullPtr(p); return n == 0 }
func (p *predIsNullType_) Equals(o any) bool                                     { CheckNotNullPtr(p); return p == o }
func (p *predIsNullType_) HashCode() int32                                       { CheckNotNullPtr(p); return 0 }
func (p *predIsNullType_) String() string                                        { CheckNotNullPtr(p); return "Predicates.isNull()" }
func (p *predIsNullType_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predAnd2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predIsNullType_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predOr2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predIsNullType_) Negate() Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return predNotNullIface_
}

type predNotNullType_ struct{ _ byte }

func (*predNotNullType_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*predNotNullType_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (*predNotNullType_) HasIdentityEquals_()                                     {}
func (p *predNotNullType_) IsJavaNil_() bool                                      { return p == nil }
func (p *predNotNullType_) Apply(x any) bool                                      { CheckNotNullPtr(p); return x != nil && !isJavaNilAny(x) }
func (p *predNotNullType_) Test(x any) bool                                       { CheckNotNullPtr(p); return x != nil && !isJavaNilAny(x) }
func (p *predNotNullType_) TestNode_(n uint32) bool                               { CheckNotNullPtr(p); return n != 0 }
func (p *predNotNullType_) Equals(o any) bool                                     { CheckNotNullPtr(p); return p == o }
func (p *predNotNullType_) HashCode() int32                                       { CheckNotNullPtr(p); return 0x4e4e554c }
func (p *predNotNullType_) String() string                                        { CheckNotNullPtr(p); return "Predicates.notNull()" }
func (p *predNotNullType_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predAnd2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predNotNullType_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predOr2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predNotNullType_) Negate() Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return predIsNullIface_
}

var (
	predAlwaysTrueSingleton_  predAlwaysTrueType_
	predAlwaysFalseSingleton_ predAlwaysFalseType_
	predIsNullSingleton_      predIsNullType_
	predNotNullSingleton_     predNotNullType_

	predAlwaysTrueIface_  Ext_com_google_common_base_Predicate = &predAlwaysTrueSingleton_
	predAlwaysFalseIface_ Ext_com_google_common_base_Predicate = &predAlwaysFalseSingleton_
	predIsNullIface_      Ext_com_google_common_base_Predicate = &predIsNullSingleton_
	predNotNullIface_     Ext_com_google_common_base_Predicate = &predNotNullSingleton_

	predAlwaysTrueAny_ any = &predAlwaysTrueSingleton_
)

type predEqualNode_ uint32

func (predEqualNode_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (predEqualNode_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (predEqualNode_) HasIdentityEquals_()                                     {}
func (predEqualNode_) IsJavaNil_() bool                                        { return false }
func (pn predEqualNode_) equalTarget_() any                                    { return boxNodeById(uint32(pn)) }
func (pn predEqualNode_) Apply(x any) bool {
	if x == nil {
		return false
	}
	if rn, ok := x.(jRawNodeProvider); ok {
		return rn.RawNodeId_() == uint32(pn)
	}
	if n, ok := x.(nodestore.Node); ok {
		return uint32(n) == uint32(pn)
	}
	return false
}
func (pn predEqualNode_) Test(x any) bool         { return pn.Apply(x) }
func (pn predEqualNode_) TestNode_(n uint32) bool { return uint32(pn) == n && n != 0 }
func (pn predEqualNode_) Equals(o any) bool {
	if op, ok := o.(predEqualNode_); ok {
		return op == pn
	}
	if op, ok := o.(jEqualTargetProvider); ok && !isJavaNilAny(o) {
		return pn.Apply(op.equalTarget_())
	}
	return false
}
func (pn predEqualNode_) HashCode() int32 { return int32(pn) }
func (pn predEqualNode_) String() string {
	return "Predicates.equalTo(" + javalang.ToJavaString(boxNodeById(uint32(pn))) + ")"
}
func (pn predEqualNode_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	return &predAnd2_{p1: pn, p2: CheckNotNull(o)}
}
func (pn predEqualNode_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	return &predOr2_{p1: pn, p2: CheckNotNull(o)}
}
func (pn predEqualNode_) Negate() Ext_java_util_function_Predicate {
	return &predNot_{inner: pn}
}

type predEqualStr_ string

var predEqualEmptyStrIface_ Ext_com_google_common_base_Predicate = predEqualStr_("")

func (predEqualStr_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (predEqualStr_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (predEqualStr_) HasIdentityEquals_()                                     {}
func (predEqualStr_) IsJavaNil_() bool                                        { return false }
func (ps predEqualStr_) equalTarget_() any                                    { return string(ps) }
func (ps predEqualStr_) Apply(x any) bool {
	s, ok := x.(string)
	return ok && s == string(ps)
}
func (ps predEqualStr_) Test(x any) bool {
	s, ok := x.(string)
	return ok && s == string(ps)
}
func (ps predEqualStr_) TestNode_(_ uint32) bool { return false }
func (ps predEqualStr_) Equals(o any) bool {
	if op, ok := o.(predEqualStr_); ok {
		return op == ps
	}
	if op, ok := o.(jEqualTargetProvider); ok && !isJavaNilAny(o) {
		return ps.Apply(op.equalTarget_())
	}
	return false
}
func (ps predEqualStr_) HashCode() int32 { return javalang.StringHashCode(string(ps)) }
func (ps predEqualStr_) String() string  { return "Predicates.equalTo(" + string(ps) + ")" }
func (ps predEqualStr_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	return &predAnd2_{p1: ps, p2: CheckNotNull(o)}
}
func (ps predEqualStr_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	return &predOr2_{p1: ps, p2: CheckNotNull(o)}
}
func (ps predEqualStr_) Negate() Ext_java_util_function_Predicate {
	return &predNot_{inner: ps}
}

type predEqualEnum_ struct {
	eb *EnumBase
}

func (*predEqualEnum_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*predEqualEnum_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (*predEqualEnum_) HasIdentityEquals_()                                     {}
func (p *predEqualEnum_) IsJavaNil_() bool                                      { return p == nil }
func (p *predEqualEnum_) equalTarget_() any                                     { return p.eb }
func (p *predEqualEnum_) Apply(x any) bool {
	CheckNotNullPtr(p)
	if x == nil {
		return false
	}
	if oa, ok := x.(EnumAccessor); ok {
		if nc, ok2 := x.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
			return false
		}
		return oa.GetEnumBase() == p.eb
	}
	return false
}
func (p *predEqualEnum_) Test(x any) bool         { return p.Apply(x) }
func (p *predEqualEnum_) TestNode_(_ uint32) bool { CheckNotNullPtr(p); return false }
func (p *predEqualEnum_) Equals(o any) bool {
	CheckNotNullPtr(p)
	if op, ok := o.(*predEqualEnum_); ok {
		return op != nil && op.eb == p.eb
	}
	if op, ok := o.(jEqualTargetProvider); ok && !isJavaNilAny(o) {
		return p.Apply(op.equalTarget_())
	}
	return false
}
func (p *predEqualEnum_) HashCode() int32 { CheckNotNullPtr(p); return p.eb.HashCode() }
func (p *predEqualEnum_) String() string {
	CheckNotNullPtr(p)
	return "Predicates.equalTo(" + p.eb.Name_ + ")"
}
func (p *predEqualEnum_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predAnd2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predEqualEnum_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predOr2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predEqualEnum_) Negate() Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predNot_{inner: p}
}

type predEqualIdentity_ struct {
	target any
}

func (*predEqualIdentity_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*predEqualIdentity_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (p *predEqualIdentity_) IsJavaNil_() bool                                      { return p == nil }
func (p *predEqualIdentity_) equalTarget_() any                                     { return p.target }
func (p *predEqualIdentity_) Apply(x any) bool                                      { CheckNotNullPtr(p); return x == p.target }
func (p *predEqualIdentity_) Test(x any) bool                                       { CheckNotNullPtr(p); return x == p.target }
func (p *predEqualIdentity_) TestNode_(_ uint32) bool                               { CheckNotNullPtr(p); return false }
func (p *predEqualIdentity_) Equals(o any) bool {
	CheckNotNullPtr(p)
	if op, ok := o.(jEqualTargetProvider); ok && !isJavaNilAny(o) {
		return op.equalTarget_() == p.target
	}
	return false
}
func (p *predEqualIdentity_) HashCode() int32 {
	CheckNotNullPtr(p)
	return javalang.ObjectHashCode(p.target)
}
func (p *predEqualIdentity_) String() string {
	CheckNotNullPtr(p)
	return "Predicates.equalTo(" + javalang.ToJavaString(p.target) + ")"
}
func (p *predEqualIdentity_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predAnd2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predEqualIdentity_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predOr2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predEqualIdentity_) Negate() Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predNot_{inner: p}
}

type predEqualAny_ struct {
	target any
}

func (*predEqualAny_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*predEqualAny_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (p *predEqualAny_) IsJavaNil_() bool                                      { return p == nil }
func (p *predEqualAny_) equalTarget_() any                                     { return p.target }
func (p *predEqualAny_) Apply(x any) bool {
	CheckNotNullPtr(p)
	if x == nil {
		return false
	}
	return javalang.ObjectsEquals(p.target, ToAny(x))
}
func (p *predEqualAny_) Test(x any) bool { return p.Apply(x) }
func (p *predEqualAny_) TestNode_(n uint32) bool {
	CheckNotNullPtr(p)
	if n == 0 {
		return false
	}
	return javalang.ObjectsEquals(p.target, boxNodeById(n))
}
func (p *predEqualAny_) Equals(o any) bool {
	CheckNotNullPtr(p)
	if op, ok := o.(jEqualTargetProvider); ok && !isJavaNilAny(o) {
		return javalang.ObjectsEquals(p.target, op.equalTarget_())
	}
	return false
}
func (p *predEqualAny_) HashCode() int32 {
	CheckNotNullPtr(p)
	return javalang.ObjectHashCode(p.target)
}
func (p *predEqualAny_) String() string {
	CheckNotNullPtr(p)
	return "Predicates.equalTo(" + javalang.ToJavaString(p.target) + ")"
}
func (p *predEqualAny_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predAnd2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predEqualAny_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predOr2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predEqualAny_) Negate() Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predNot_{inner: p}
}

var (
	predEqualBoolTrue_  Ext_com_google_common_base_Predicate = &predEqualAny_{target: Boolean_TRUE}
	predEqualBoolFalse_ Ext_com_google_common_base_Predicate = &predEqualAny_{target: Boolean_FALSE}
)

func guavaPredEqualTo(raw any) Ext_com_google_common_base_Predicate {
	if raw == nil || isJavaNilAny(raw) {
		return predIsNullIface_
	}
	if ea, ok := raw.(EnumAccessor); ok {
		if eb := ea.GetEnumBase(); eb != nil {
			if p := eb.equalPred.Load(); p != nil {
				return p
			}
			np := &predEqualEnum_{eb: eb}
			if eb.equalPred.CompareAndSwap(nil, np) {
				return np
			}
			return eb.equalPred.Load()
		}
	}
	if rn, ok := raw.(jRawNodeProvider); ok {
		id := rn.RawNodeId_()
		if id == 0 {
			return predIsNullIface_
		}
		return predEqualNode_(id)
	}
	switch v := raw.(type) {
	case nodestore.Node:
		if v == 0 {
			return predIsNullIface_
		}
		return predEqualNode_(v)
	case string:
		if v == "" {
			return predEqualEmptyStrIface_
		}
		return predEqualStr_(v)
	case *BoxBoolean:
		if v.Val {
			return predEqualBoolTrue_
		}
		return predEqualBoolFalse_
	}
	if _, ok := raw.(IdentityComparable); ok {
		return &predEqualIdentity_{target: raw}
	}
	return &predEqualAny_{target: ToAny(raw)}
}

type predNot_ struct {
	inner Ext_java_util_function_Predicate
}

func (*predNot_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*predNot_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (p *predNot_) IsJavaNil_() bool                                      { return p == nil }
func (p *predNot_) Apply(x any) bool                                      { CheckNotNullPtr(p); return !p.inner.Test(x) }
func (p *predNot_) Test(x any) bool                                       { CheckNotNullPtr(p); return !p.inner.Test(x) }
func (p *predNot_) TestNode_(n uint32) bool                               { CheckNotNullPtr(p); return !PredTestNode(p.inner, n) }
func (p *predNot_) Equals(o any) bool {
	CheckNotNullPtr(p)
	op, ok := o.(*predNot_)
	return ok && op != nil && javalang.ObjectsEquals(p.inner, op.inner)
}
func (p *predNot_) HashCode() int32 { CheckNotNullPtr(p); return ^javalang.ObjectHashCode(p.inner) }
func (p *predNot_) String() string {
	CheckNotNullPtr(p)
	return "Predicates.not(" + javalang.ToJavaString(p.inner) + ")"
}
func (p *predNot_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predAnd2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predNot_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predOr2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predNot_) Negate() Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return p.inner
}

func guavaPredNot(target Ext_com_google_common_base_Predicate) Ext_com_google_common_base_Predicate {
	CheckNotNull(target)
	return &predNot_{inner: target}
}

func jdkPredNot(target Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNull(target)
	return &predNot_{inner: target}
}

type predAnd2_ struct {
	p1 Ext_java_util_function_Predicate
	p2 Ext_java_util_function_Predicate
}

func (*predAnd2_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*predAnd2_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (p *predAnd2_) IsJavaNil_() bool                                      { return p == nil }
func (p *predAnd2_) Apply(x any) bool                                      { CheckNotNullPtr(p); return p.p1.Test(x) && p.p2.Test(x) }
func (p *predAnd2_) Test(x any) bool                                       { CheckNotNullPtr(p); return p.p1.Test(x) && p.p2.Test(x) }
func (p *predAnd2_) TestNode_(n uint32) bool {
	CheckNotNullPtr(p)
	return PredTestNode(p.p1, n) && PredTestNode(p.p2, n)
}
func (p *predAnd2_) Equals(o any) bool {
	CheckNotNullPtr(p)
	op, ok := o.(*predAnd2_)
	return ok && op != nil && javalang.ObjectsEquals(p.p1, op.p1) && javalang.ObjectsEquals(p.p2, op.p2)
}
func (p *predAnd2_) HashCode() int32 {
	CheckNotNullPtr(p)
	return javalang.ObjectHashCode(p.p1)*31 + javalang.ObjectHashCode(p.p2) + 0x12472c2c
}
func (p *predAnd2_) String() string {
	CheckNotNullPtr(p)
	return "Predicates.and(" + javalang.ToJavaString(p.p1) + "," + javalang.ToJavaString(p.p2) + ")"
}
func (p *predAnd2_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predAnd2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predAnd2_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predOr2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predAnd2_) Negate() Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predNot_{inner: p}
}

func guavaPredAnd2(a, b Ext_com_google_common_base_Predicate) Ext_com_google_common_base_Predicate {
	CheckNotNull(a)
	CheckNotNull(b)
	return &predAnd2_{p1: a, p2: b}
}

type predOr2_ struct {
	p1 Ext_java_util_function_Predicate
	p2 Ext_java_util_function_Predicate
}

func (*predOr2_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*predOr2_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (p *predOr2_) IsJavaNil_() bool                                      { return p == nil }
func (p *predOr2_) Apply(x any) bool                                      { CheckNotNullPtr(p); return p.p1.Test(x) || p.p2.Test(x) }
func (p *predOr2_) Test(x any) bool                                       { CheckNotNullPtr(p); return p.p1.Test(x) || p.p2.Test(x) }
func (p *predOr2_) TestNode_(n uint32) bool {
	CheckNotNullPtr(p)
	return PredTestNode(p.p1, n) || PredTestNode(p.p2, n)
}
func (p *predOr2_) Equals(o any) bool {
	CheckNotNullPtr(p)
	op, ok := o.(*predOr2_)
	return ok && op != nil && javalang.ObjectsEquals(p.p1, op.p1) && javalang.ObjectsEquals(p.p2, op.p2)
}
func (p *predOr2_) HashCode() int32 {
	CheckNotNullPtr(p)
	return javalang.ObjectHashCode(p.p1)*31 + javalang.ObjectHashCode(p.p2) + 0x053c91cf
}
func (p *predOr2_) String() string {
	CheckNotNullPtr(p)
	return "Predicates.or(" + javalang.ToJavaString(p.p1) + "," + javalang.ToJavaString(p.p2) + ")"
}
func (p *predOr2_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predAnd2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predOr2_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predOr2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predOr2_) Negate() Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predNot_{inner: p}
}

func guavaPredOr2(a, b Ext_com_google_common_base_Predicate) Ext_com_google_common_base_Predicate {
	CheckNotNull(a)
	CheckNotNull(b)
	return &predOr2_{p1: a, p2: b}
}

type predAndSlice_ struct {
	preds []Ext_java_util_function_Predicate
}

func (*predAndSlice_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*predAndSlice_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (p *predAndSlice_) IsJavaNil_() bool                                      { return p == nil }
func (p *predAndSlice_) Apply(x any) bool {
	CheckNotNullPtr(p)
	for _, q := range p.preds {
		if !q.Test(x) {
			return false
		}
	}
	return true
}
func (p *predAndSlice_) Test(x any) bool { return p.Apply(x) }
func (p *predAndSlice_) TestNode_(n uint32) bool {
	CheckNotNullPtr(p)
	for _, q := range p.preds {
		if !PredTestNode(q, n) {
			return false
		}
	}
	return true
}
func (p *predAndSlice_) Equals(o any) bool {
	CheckNotNullPtr(p)
	op, ok := o.(*predAndSlice_)
	if !ok || op == nil || len(p.preds) != len(op.preds) {
		return false
	}
	for i := range p.preds {
		if !javalang.ObjectsEquals(p.preds[i], op.preds[i]) {
			return false
		}
	}
	return true
}
func (p *predAndSlice_) HashCode() int32 {
	CheckNotNullPtr(p)
	h := int32(1)
	for _, q := range p.preds {
		h = 31*h + javalang.ObjectHashCode(q)
	}
	return h + 0x12472c2c
}
func (p *predAndSlice_) String() string { CheckNotNullPtr(p); return "Predicates.and(...)" }
func (p *predAndSlice_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predAnd2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predAndSlice_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predOr2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predAndSlice_) Negate() Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predNot_{inner: p}
}

type predOrSlice_ struct {
	preds []Ext_java_util_function_Predicate
}

func (*predOrSlice_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*predOrSlice_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (p *predOrSlice_) IsJavaNil_() bool                                      { return p == nil }
func (p *predOrSlice_) Apply(x any) bool {
	CheckNotNullPtr(p)
	for _, q := range p.preds {
		if q.Test(x) {
			return true
		}
	}
	return false
}
func (p *predOrSlice_) Test(x any) bool { return p.Apply(x) }
func (p *predOrSlice_) TestNode_(n uint32) bool {
	CheckNotNullPtr(p)
	for _, q := range p.preds {
		if PredTestNode(q, n) {
			return true
		}
	}
	return false
}
func (p *predOrSlice_) Equals(o any) bool {
	CheckNotNullPtr(p)
	op, ok := o.(*predOrSlice_)
	if !ok || op == nil || len(p.preds) != len(op.preds) {
		return false
	}
	for i := range p.preds {
		if !javalang.ObjectsEquals(p.preds[i], op.preds[i]) {
			return false
		}
	}
	return true
}
func (p *predOrSlice_) HashCode() int32 {
	CheckNotNullPtr(p)
	h := int32(1)
	for _, q := range p.preds {
		h = 31*h + javalang.ObjectHashCode(q)
	}
	return h + 0x053c91cf
}
func (p *predOrSlice_) String() string { CheckNotNullPtr(p); return "Predicates.or(...)" }
func (p *predOrSlice_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predAnd2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predOrSlice_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predOr2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predOrSlice_) Negate() Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predNot_{inner: p}
}

func collectPredIter_(src Ext_java_lang_Iterable) []Ext_java_util_function_Predicate {
	CheckNotNull(src)
	if si, ok := src.(SliceIterable); ok {
		raw := si.IterSlice_()
		out := make([]Ext_java_util_function_Predicate, len(raw))
		for i, x := range raw {
			out[i] = AsIface[Ext_java_util_function_Predicate](CheckNotNull(x))
		}
		return out
	}
	var out []Ext_java_util_function_Predicate
	for it := NewIterCursor(src); it.HasNext(); {
		out = append(out, AsIface[Ext_java_util_function_Predicate](CheckNotNull(it.Next())))
	}
	return out
}

func guavaPredAndFromList_(elems []Ext_java_util_function_Predicate) Ext_com_google_common_base_Predicate {
	switch len(elems) {
	case 0:
		return predAlwaysTrueIface_
	case 1:
		if gp, ok := elems[0].(Ext_com_google_common_base_Predicate); ok {
			return gp
		}
	case 2:
		return &predAnd2_{p1: elems[0], p2: elems[1]}
	}
	return &predAndSlice_{preds: elems}
}

func guavaPredOrFromList_(elems []Ext_java_util_function_Predicate) Ext_com_google_common_base_Predicate {
	switch len(elems) {
	case 0:
		return predAlwaysFalseIface_
	case 1:
		if gp, ok := elems[0].(Ext_com_google_common_base_Predicate); ok {
			return gp
		}
	case 2:
		return &predOr2_{p1: elems[0], p2: elems[1]}
	}
	return &predOrSlice_{preds: elems}
}

func guavaPredAndIter(src Ext_java_lang_Iterable) Ext_com_google_common_base_Predicate {
	return guavaPredAndFromList_(collectPredIter_(src))
}

func guavaPredOrIter(src Ext_java_lang_Iterable) Ext_com_google_common_base_Predicate {
	return guavaPredOrFromList_(collectPredIter_(src))
}

func guavaPredAndSlice[E any](src []E) Ext_com_google_common_base_Predicate {
	CheckNotNullSlice(src)
	elems := make([]Ext_java_util_function_Predicate, len(src))
	for i, x := range src {
		elems[i] = AsIface[Ext_java_util_function_Predicate](CheckNotNull(x))
	}
	return guavaPredAndFromList_(elems)
}

func guavaPredOrSlice[E any](src []E) Ext_com_google_common_base_Predicate {
	CheckNotNullSlice(src)
	elems := make([]Ext_java_util_function_Predicate, len(src))
	for i, x := range src {
		elems[i] = AsIface[Ext_java_util_function_Predicate](CheckNotNull(x))
	}
	return guavaPredOrFromList_(elems)
}

type predInCollection_ struct {
	col Ext_java_util_Collection
}

func (*predInCollection_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*predInCollection_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (p *predInCollection_) IsJavaNil_() bool                                      { return p == nil }
func (p *predInCollection_) Apply(x any) (matched bool) {
	CheckNotNullPtr(p)
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case *javalang.NullPointerException, *javalang.ClassCastException:
				matched = false
			default:
				panic(r)
			}
		}
	}()
	return predColContainsAny_(p.col, ToAny(x))
}
func (p *predInCollection_) Test(x any) bool { return p.Apply(x) }
func (p *predInCollection_) TestNode_(n uint32) bool {
	CheckNotNullPtr(p)
	if fs, ok := p.col.(FastNodeSet); ok {
		return fs.SetContainsNode_(n)
	}
	return p.Apply(boxNodeById(n))
}
func (p *predInCollection_) Equals(o any) bool {
	CheckNotNullPtr(p)
	op, ok := o.(*predInCollection_)
	return ok && op != nil && javalang.ObjectsEquals(p.col, op.col)
}
func (p *predInCollection_) HashCode() int32 {
	CheckNotNullPtr(p)
	return javalang.ObjectHashCode(p.col)
}
func (p *predInCollection_) String() string {
	CheckNotNullPtr(p)
	return "Predicates.in(" + javalang.ToJavaString(p.col) + ")"
}
func (p *predInCollection_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predAnd2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predInCollection_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predOr2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predInCollection_) Negate() Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predNot_{inner: p}
}

func guavaPredInCollection(col Ext_java_util_Collection) Ext_com_google_common_base_Predicate {
	CheckNotNull(col)
	return &predInCollection_{col: col}
}

type predInstanceOf_ struct {
	cls *javalang.Class
}

func (*predInstanceOf_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*predInstanceOf_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (p *predInstanceOf_) IsJavaNil_() bool                                      { return p == nil }
func (p *predInstanceOf_) Apply(x any) bool                                      { CheckNotNullPtr(p); return CheckTypePredicate(p.cls, x) }
func (p *predInstanceOf_) Test(x any) bool                                       { CheckNotNullPtr(p); return CheckTypePredicate(p.cls, x) }
func (p *predInstanceOf_) TestNode_(n uint32) bool {
	CheckNotNullPtr(p)
	if n == 0 {
		return false
	}
	return CheckTypePredicate(p.cls, boxNodeById(n))
}
func (p *predInstanceOf_) Equals(o any) bool {
	CheckNotNullPtr(p)
	op, ok := o.(*predInstanceOf_)
	return ok && op != nil && p.cls == op.cls
}
func (p *predInstanceOf_) HashCode() int32 { CheckNotNullPtr(p); return javalang.ObjectHashCode(p.cls) }
func (p *predInstanceOf_) String() string {
	CheckNotNullPtr(p)
	return "Predicates.instanceOf(" + p.cls.GetName() + ")"
}
func (p *predInstanceOf_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predAnd2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predInstanceOf_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predOr2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predInstanceOf_) Negate() Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predNot_{inner: p}
}

func guavaPredInstanceOf(cls *javalang.Class) Ext_com_google_common_base_Predicate {
	CheckNotNullPtr(cls)
	return &predInstanceOf_{cls: cls}
}

type predCompose_ struct {
	pred Ext_com_google_common_base_Predicate
	fn   Ext_com_google_common_base_Function
}

func (*predCompose_) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*predCompose_) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (p *predCompose_) IsJavaNil_() bool                                      { return p == nil }
func (p *predCompose_) Apply(x any) bool                                      { CheckNotNullPtr(p); return p.pred.Apply(p.fn.Apply(x)) }
func (p *predCompose_) Test(x any) bool                                       { CheckNotNullPtr(p); return p.pred.Apply(p.fn.Apply(x)) }
func (p *predCompose_) TestNode_(n uint32) bool {
	CheckNotNullPtr(p)
	return p.pred.Apply(p.fn.Apply(boxNodeById(n)))
}
func (p *predCompose_) Equals(o any) bool {
	CheckNotNullPtr(p)
	op, ok := o.(*predCompose_)
	return ok && op != nil && javalang.ObjectsEquals(p.fn, op.fn) && javalang.ObjectsEquals(p.pred, op.pred)
}
func (p *predCompose_) HashCode() int32 {
	CheckNotNullPtr(p)
	return javalang.ObjectHashCode(p.fn) ^ javalang.ObjectHashCode(p.pred)
}
func (p *predCompose_) String() string {
	CheckNotNullPtr(p)
	return javalang.ToJavaString(p.pred) + "(" + javalang.ToJavaString(p.fn) + ")"
}
func (p *predCompose_) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predAnd2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predCompose_) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predOr2_{p1: p, p2: CheckNotNull(o)}
}
func (p *predCompose_) Negate() Ext_java_util_function_Predicate {
	CheckNotNullPtr(p)
	return &predNot_{inner: p}
}

func guavaPredCompose(pred Ext_com_google_common_base_Predicate, fn Ext_com_google_common_base_Function) Ext_com_google_common_base_Predicate {
	CheckNotNull(pred)
	CheckNotNull(fn)
	return &predCompose_{pred: pred, fn: fn}
}

// --- Class Literal Cache ---

var classCache sync.Map // string -> *javalang.Class

func ClassLiteral(qualName, simpleName string) *javalang.Class {
	if c, ok := classCache.Load(qualName); ok {
		return c.(*javalang.Class)
	}
	c := javalang.NewClass(qualName, simpleName)
	actual, _ := classCache.LoadOrStore(qualName, c)
	return actual.(*javalang.Class)
}

func GetObjectClass(v any) *javalang.Class {
	if v == nil || javalang.IsNil(v) {
		panic(javalang.NewNullPointerException())
	}
	if s, ok := v.(string); ok {
		if s == NullString {
			panic(javalang.NewNullPointerException())
		}
		return ClassLiteral("java.lang.String", "String")
	}
	if n, ok := v.(nodestore.Node); ok {
		if n == 0 {
			panic(javalang.NewNullPointerException())
		}
		switch n.KindFlags() {
		case 1:
			return ClassLiteral("com.google.javascript.rhino.Node$NumberNode", "NumberNode")
		case 2:
			return ClassLiteral("com.google.javascript.rhino.Node$BigIntNode", "BigIntNode")
		case 3:
			return ClassLiteral("com.google.javascript.rhino.Node$StringNode", "StringNode")
		case 4:
			return ClassLiteral("com.google.javascript.rhino.Node$TemplateLiteralSubstringNode", "TemplateLiteralSubstringNode")
		default:
			return ClassLiteral("com.google.javascript.rhino.Node", "Node")
		}
	}
	if ci, ok := v.(interface{ JavaClassInfo_() (string, string) }); ok {
		qn, sn := ci.JavaClassInfo_()
		return ClassLiteral(qn, sn)
	}
	t := reflect.TypeOf(v)
	return ClassLiteral(t.String(), t.Name())
}

// --- Throwable Helpers ---

func InitUserThrowable(t javalang.Throwable, className string, args ...any) {
	msg := ""
	hasMsg := false
	var cause javalang.Throwable
	for _, a := range args {
		if s, ok := a.(string); ok {
			if s != NullString {
				msg = s
				hasMsg = true
			}
		} else if th, ok := a.(javalang.Throwable); ok {
			cause = th
			if !hasMsg && th != nil && !javalang.IsNil(th) {
				msg = th.Error()
				hasMsg = true
			}
		} else if a != nil && !javalang.IsNil(a) {
			msg = StrVal(a)
			hasMsg = true
		}
	}
	if hasMsg {
		javalang.InitThrowable(t, className, msg, cause)
	} else {
		javalang.InitThrowable(t, className, "", cause)
	}
}

func ConstructShimThrowable(qualName string, args ...any) javalang.Throwable {
	var t javalang.Throwable
	switch qualName {
	case "java.lang.Throwable", "java.lang.Exception":
		t = &javalang.JavaException{}
	case "java.lang.RuntimeException":
		t = &javalang.RuntimeException{}
	case "java.lang.IllegalArgumentException":
		t = &javalang.IllegalArgumentException{}
	case "java.lang.IllegalStateException":
		t = &javalang.IllegalStateException{}
	case "java.lang.NullPointerException":
		t = &javalang.NullPointerException{}
	case "java.lang.IndexOutOfBoundsException":
		t = &javalang.IndexOutOfBoundsException{}
	case "java.lang.ArrayIndexOutOfBoundsException":
		t = &javalang.ArrayIndexOutOfBoundsException{}
	case "java.lang.StringIndexOutOfBoundsException":
		t = &javalang.StringIndexOutOfBoundsException{}
	case "java.lang.UnsupportedOperationException":
		t = &javalang.UnsupportedOperationException{}
	case "java.lang.ClassCastException":
		t = &javalang.ClassCastException{}
	case "java.util.NoSuchElementException":
		t = &javalang.NoSuchElementException{}
	case "java.lang.ArithmeticException":
		t = &javalang.ArithmeticException{}
	case "java.lang.NumberFormatException":
		t = &javalang.NumberFormatException{}
	case "java.lang.Error":
		t = &javalang.JavaLangError{}
	case "java.lang.AssertionError":
		t = &javalang.AssertionError{}
	case "java.io.IOException":
		t = &javaio.IOException{}
	case "java.io.FileNotFoundException":
		t = &javaio.FileNotFoundException{}
	case "java.io.EOFException":
		t = &javaio.EOFException{}
	default:
		t = &javalang.RuntimeException{}
	}
	InitUserThrowable(t, qualName, args...)
	return t
}

func ThrowableGetMessage(t javalang.Throwable) string {
	if t == nil || javalang.IsNil(t) {
		panic(javalang.NewNullPointerException())
	}
	return t.GetMessage()
}

// --- String & CharSequence Intrinsics ---

func StringEquals(s string, other any) bool {
	if s == NullString {
		panic(javalang.NewNullPointerException())
	}
	o, ok := other.(string)
	return ok && o != NullString && s == o
}

func StringFormat(format string, args ...any) string {
	CheckNotNullString(format)
	var out strings.Builder
	out.Grow(len(format) + 16)
	seqIdx := 0
	lastIdx := 0
	for i := 0; i < len(format); {
		if format[i] != '%' {
			j := strings.IndexByte(format[i:], '%')
			if j < 0 {
				out.WriteString(format[i:])
				break
			}
			out.WriteString(format[i : i+j])
			i += j
		}
		i++ // consume '%'
		if i >= len(format) {
			out.WriteByte('%')
			break
		}
		if format[i] == '%' {
			out.WriteByte('%')
			i++
			continue
		}
		if format[i] == 'n' {
			out.WriteByte('\n')
			i++
			continue
		}
		// Optional [argument_index$]
		explicitIdx := -1
		usePrev := false
		save := i
		if format[i] >= '1' && format[i] <= '9' {
			n := 0
			for i < len(format) && format[i] >= '0' && format[i] <= '9' {
				n = n*10 + int(format[i]-'0')
				i++
			}
			if i < len(format) && format[i] == '$' {
				explicitIdx = n - 1
				i++
			} else {
				i = save
			}
		}
		// Flags
		var flags strings.Builder
		leftAlign := false
		commaFlag := false
		parenFlag := false
		for i < len(format) {
			c := format[i]
			if c == '<' {
				usePrev = true
				i++
			} else if c == '-' {
				leftAlign = true
				flags.WriteByte(c)
				i++
			} else if c == ',' {
				commaFlag = true
				i++
			} else if c == '(' {
				parenFlag = true
				i++
			} else if c == '#' || c == '+' || c == ' ' || c == '0' {
				flags.WriteByte(c)
				i++
			} else {
				break
			}
		}
		_ = commaFlag
		_ = parenFlag
		// Width
		width := -1
		widthStart := i
		if i < len(format) && format[i] >= '1' && format[i] <= '9' {
			w := 0
			for i < len(format) && format[i] >= '0' && format[i] <= '9' {
				w = w*10 + int(format[i]-'0')
				i++
			}
			width = w
		}
		widthStr := format[widthStart:i]
		// Precision
		prec := -1
		precStr := ""
		if i < len(format) && format[i] == '.' {
			pStart := i
			i++
			p := 0
			for i < len(format) && format[i] >= '0' && format[i] <= '9' {
				p = p*10 + int(format[i]-'0')
				i++
			}
			prec = p
			precStr = format[pStart:i]
		}
		if i >= len(format) {
			break
		}
		convCh := format[i]
		i++
		var argIdx int
		if explicitIdx >= 0 {
			argIdx = explicitIdx
		} else if usePrev {
			argIdx = lastIdx
		} else {
			argIdx = seqIdx
			seqIdx++
		}
		lastIdx = argIdx
		var arg any
		if argIdx >= 0 && argIdx < len(args) {
			arg = args[argIdx]
		}
		padAndWrite := func(s string) {
			if prec >= 0 && javalang.StringLength(s) > prec {
				s = javalang.StringSubstring(s, 0, prec)
			}
			sLen := javalang.StringLength(s)
			if width > sLen {
				pad := strings.Repeat(" ", width-sLen)
				if leftAlign {
					out.WriteString(s)
					out.WriteString(pad)
				} else {
					out.WriteString(pad)
					out.WriteString(s)
				}
			} else {
				out.WriteString(s)
			}
		}
		switch convCh {
		case 's', 'S':
			s := StrVal(arg)
			if convCh == 'S' {
				s = strings.ToUpper(s)
			}
			padAndWrite(s)
		case 'b', 'B':
			bv := false
			if arg != nil && !javalang.IsNil(arg) {
				switch x := arg.(type) {
				case bool:
					bv = x
				case *BoxBoolean:
					bv = x != nil && x.Val
				default:
					bv = true
				}
			}
			s := "false"
			if bv {
				s = "true"
			}
			if convCh == 'B' {
				s = strings.ToUpper(s)
			}
			padAndWrite(s)
		case 'h', 'H':
			s := "null"
			if arg != nil && !javalang.IsNil(arg) {
				s = strconv.FormatUint(uint64(uint32(javalang.ObjectHashCode(arg))), 16)
				if convCh == 'H' {
					s = strings.ToUpper(s)
				}
			}
			padAndWrite(s)
		case 'c', 'C':
			if arg == nil || javalang.IsNil(arg) {
				padAndWrite("null")
			} else {
				var r rune
				switch x := arg.(type) {
				case uint16:
					r = rune(x)
				case *BoxChar:
					r = rune(x.Val)
				case int32:
					r = x
				case *BoxInt:
					r = x.Val
				case int8:
					r = rune(uint8(x))
				case *BoxByte:
					r = rune(uint8(x.Val))
				case int16:
					r = rune(uint16(x))
				case *BoxShort:
					r = rune(uint16(x.Val))
				}
				s := string(r)
				if convCh == 'C' {
					s = strings.ToUpper(s)
				}
				prec = -1
				padAndWrite(s)
			}
		case 'd', 'o', 'x', 'X':
			if arg == nil || javalang.IsNil(arg) {
				padAndWrite("null")
			} else {
				goVerb := "%" + flags.String() + widthStr + string(convCh)
				var num any = arg
				if convCh == 'd' {
					switch x := arg.(type) {
					case *BoxInt:
						num = x.Val
					case *BoxLong:
						num = x.Val
					case *BoxShort:
						num = x.Val
					case *BoxByte:
						num = x.Val
					case Number:
						num = x.LongValue()
					}
				} else {
					switch x := arg.(type) {
					case int8:
						num = uint8(x)
					case *BoxByte:
						num = uint8(x.Val)
					case int16:
						num = uint16(x)
					case *BoxShort:
						num = uint16(x.Val)
					case int32:
						num = uint32(x)
					case *BoxInt:
						num = uint32(x.Val)
					case int64:
						num = uint64(x)
					case *BoxLong:
						num = uint64(x.Val)
					case Number:
						num = uint64(x.LongValue())
					}
				}
				out.WriteString(fmt.Sprintf(goVerb, num))
			}
		case 'e', 'E', 'f', 'g', 'G', 'a', 'A':
			if arg == nil || javalang.IsNil(arg) {
				padAndWrite("null")
			} else {
				var fv float64
				switch x := arg.(type) {
				case float64:
					fv = x
				case float32:
					fv = float64(x)
				case *BoxDouble:
					fv = x.Val
				case *BoxFloat:
					fv = float64(x.Val)
				case Number:
					fv = x.DoubleValue()
				}
				if math.IsNaN(fv) {
					prec = -1
					padAndWrite("NaN")
				} else if math.IsInf(fv, 1) {
					prec = -1
					if convCh == 'E' || convCh == 'G' || convCh == 'A' {
						padAndWrite("INFINITY")
					} else {
						padAndWrite("Infinity")
					}
				} else if math.IsInf(fv, -1) {
					prec = -1
					if convCh == 'E' || convCh == 'G' || convCh == 'A' {
						padAndWrite("-INFINITY")
					} else {
						padAndWrite("-Infinity")
					}
				} else {
					vCh := convCh
					if vCh == 'a' {
						vCh = 'x'
					} else if vCh == 'A' {
						vCh = 'X'
					}
					goVerb := "%" + flags.String() + widthStr + precStr + string(vCh)
					out.WriteString(fmt.Sprintf(goVerb, fv))
				}
			}
		default:
			padAndWrite(StrVal(arg))
		}
	}
	return out.String()
}

func JStrLen(s string) int32 {
	if len(s) < 48 {
		for i := 0; i < len(s); i++ {
			if s[i] >= 0x80 {
				if s == NullString {
					ThrowNPE()
				}
				return int32(javalang.StringLength(s))
			}
		}
		return int32(len(s))
	}
	if s == NullString {
		ThrowNPE()
	}
	return int32(javalang.StringLength(s))
}
func JStrCharAt(s string, idx int32) uint16 {
	CheckNotNullString(s)
	return javalang.StringCharAt(s, int(idx))
}
func JStrCodePointAt(s string, idx int32) int32 {
	CheckNotNullString(s)
	return javalang.StringCodePointAt(s, int(idx))
}
func JStrSub(s string, begin, end int32) string {
	CheckNotNullString(s)
	return javalang.StringSubstring(s, int(begin), int(end))
}
func JStrSubFrom(s string, begin int32) string {
	CheckNotNullString(s)
	return javalang.StringSubstringFrom(s, int(begin))
}
func JStrIndexOf(s, sub string) int32 {
	CheckNotNullString(s)
	CheckNotNullString(sub)
	return int32(javalang.StringIndexOf(s, sub))
}
func JStrIndexOfFrom(s, sub string, from int32) int32 {
	CheckNotNullString(s)
	CheckNotNullString(sub)
	return int32(javalang.StringIndexOfFrom(s, sub, int(from)))
}
func JStrIndexOfCp(s string, cp int32) int32 {
	CheckNotNullString(s)
	if cp >= 0 && cp <= 0xFFFF {
		return int32(javalang.StringIndexOfChar(s, uint16(cp)))
	}
	if cp > 0xFFFF && cp <= 0x10FFFF {
		return JStrIndexOf(s, string(cp))
	}
	return -1
}
func JStrIndexOfCpFrom(s string, cp, from int32) int32 {
	CheckNotNullString(s)
	if cp >= 0 && cp <= 0xFFFF {
		return JStrIndexOfFrom(s, javalang.StringFromChar(uint16(cp)), from)
	}
	if cp > 0xFFFF && cp <= 0x10FFFF {
		return JStrIndexOfFrom(s, string(cp), from)
	}
	return -1
}
func JStrLastIndexOf(s, sub string) int32 {
	CheckNotNullString(s)
	CheckNotNullString(sub)
	return int32(javalang.StringLastIndexOf(s, sub))
}
func JStrLastIndexOfFrom(s, sub string, from int32) int32 {
	CheckNotNullString(s)
	CheckNotNullString(sub)
	return int32(javalang.StringLastIndexOfFrom(s, sub, int(from)))
}
func JStrLastIndexOfCp(s string, cp int32) int32 {
	CheckNotNullString(s)
	if cp >= 0 && cp <= 0xFFFF {
		return int32(javalang.StringLastIndexOfCharFrom(s, uint16(cp), javalang.StringLength(s)))
	}
	if cp > 0xFFFF && cp <= 0x10FFFF {
		return JStrLastIndexOf(s, string(cp))
	}
	return -1
}
func JStrLastIndexOfCpFrom(s string, cp, from int32) int32 {
	CheckNotNullString(s)
	if cp >= 0 && cp <= 0xFFFF {
		return int32(javalang.StringLastIndexOfCharFrom(s, uint16(cp), int(from)))
	}
	if cp > 0xFFFF && cp <= 0x10FFFF {
		return JStrLastIndexOfFrom(s, string(cp), from)
	}
	return -1
}
func JStrStartsWithOffset(s, prefix string, offset int32) bool {
	CheckNotNullString(s)
	CheckNotNullString(prefix)
	return javalang.StringRegionMatches(s, int(offset), prefix, 0, javalang.StringLength(prefix))
}
func JStrReplaceChar(s string, oldCh, newCh uint16) string {
	CheckNotNullString(s)
	if oldCh == newCh {
		return s
	}
	return javalang.StringReplace(s, javalang.StringFromChar(oldCh), javalang.StringFromChar(newCh))
}
func JStrCompareTo(a, b string) int32 {
	CheckNotNullString(a)
	CheckNotNullString(b)
	return int32(javalang.StringCompareTo(a, b))
}
func JStrSplit(s, regex string) []string {
	CheckNotNullString(s)
	CheckNotNullString(regex)
	return javaregex.StringSplit(s, regex)
}
func JStrSplitLimit(s, regex string, limit int32) []string {
	CheckNotNullString(s)
	CheckNotNullString(regex)
	return javaregex.StringSplitLimit(s, regex, int(limit))
}
func JStrMatches(s, regex string) bool {
	CheckNotNullString(s)
	CheckNotNullString(regex)
	return javaregex.StringMatches(s, regex)
}
func JStrReplaceAll(s, regex, repl string) string {
	CheckNotNullString(s)
	CheckNotNullString(regex)
	CheckNotNullString(repl)
	return javaregex.StringReplaceAll(s, regex, repl)
}
func JStrReplaceFirst(s, regex, repl string) string {
	CheckNotNullString(s)
	CheckNotNullString(regex)
	CheckNotNullString(repl)
	return javaregex.StringReplaceFirst(s, regex, repl)
}
func JStrRegionMatches(s string, ignoreCase bool, toffset int32, other string, ooffset, length int32) bool {
	CheckNotNullString(s)
	CheckNotNullString(other)
	if !ignoreCase {
		return javalang.StringRegionMatches(s, int(toffset), other, int(ooffset), int(length))
	}
	if toffset < 0 || ooffset < 0 || length < 0 {
		return false
	}
	sLen := javalang.StringLength(s)
	oLen := javalang.StringLength(other)
	if int(toffset)+int(length) > sLen || int(ooffset)+int(length) > oLen {
		return false
	}
	a := javalang.StringSubstring(s, int(toffset), int(toffset+length))
	b := javalang.StringSubstring(other, int(ooffset), int(ooffset+length))
	return strings.EqualFold(a, b)
}
func JStrGetBytes(s string) []int8 {
	CheckNotNullString(s)
	out := make([]int8, len(s))
	for i := 0; i < len(s); i++ {
		out[i] = int8(s[i])
	}
	return out
}

// --- O(1) UTF-16 Length & WTF-8 StringBuilder Core ---

type jStringBuilderCore struct {
	buf         []byte
	cachedStr   string
	utf16Len    int32  // >= 0 exact UTF-16 code unit count; -1 dirty
	pendingHigh uint16 // 0 = none; 0xD800..0xDBFF = trailing lone high surrogate held out of buf
	hasCache    bool
}

func (c *jStringBuilderCore) invalidateCache() {
	if c.hasCache {
		c.cachedStr = ""
		c.hasCache = false
	}
}

func (c *jStringBuilderCore) grow(capHint int) {
	if capHint > 0 && capHint <= (1<<24) && cap(c.buf) < capHint {
		nb := make([]byte, len(c.buf), capHint)
		copy(nb, c.buf)
		c.buf = nb
	}
}

func (c *jStringBuilderCore) reset() {
	c.buf = c.buf[:0]
	c.cachedStr = ""
	c.hasCache = false
	c.utf16Len = 0
	c.pendingHigh = 0
}

func (c *jStringBuilderCore) resetTo(s string) {
	c.reset()
	c.appendString(s)
}

func (c *jStringBuilderCore) flushPendingHigh() {
	if c.pendingHigh != 0 {
		ch := c.pendingHigh
		c.pendingHigh = 0
		c.invalidateCache()
		c.buf = append(c.buf,
			byte(0xE0|(ch>>12)),
			byte(0x80|((ch>>6)&0x3F)),
			byte(0x80|(ch&0x3F)))
	}
}

func (c *jStringBuilderCore) str() string {
	if len(c.buf) == 0 && c.pendingHigh == 0 {
		return ""
	}
	if c.pendingHigh == 0 {
		if !c.hasCache {
			c.cachedStr = string(c.buf)
			c.hasCache = true
		}
		return c.cachedStr
	}
	ch := c.pendingHigh
	out := make([]byte, len(c.buf)+3)
	copy(out, c.buf)
	out[len(c.buf)] = byte(0xE0 | (ch >> 12))
	out[len(c.buf)+1] = byte(0x80 | ((ch >> 6) & 0x3F))
	out[len(c.buf)+2] = byte(0x80 | (ch & 0x3F))
	return string(out)
}

func (c *jStringBuilderCore) length() int32 {
	if c.utf16Len < 0 {
		n := int32(javalang.StringLength(string(c.buf)))
		if c.pendingHigh != 0 {
			n++
		}
		c.utf16Len = n
	}
	return c.utf16Len
}

func (c *jStringBuilderCore) appendChar(ch uint16) {
	if c.hasCache {
		c.cachedStr = ""
		c.hasCache = false
	}
	if c.pendingHigh == 0 && ch < 0x80 {
		c.buf = append(c.buf, byte(ch))
		if c.utf16Len >= 0 {
			c.utf16Len++
		}
		return
	}
	if c.utf16Len >= 0 {
		c.utf16Len++
	}
	if c.pendingHigh != 0 {
		if ch >= 0xDC00 && ch <= 0xDFFF {
			r := utf16.DecodeRune(rune(c.pendingHigh), rune(ch))
			c.pendingHigh = 0
			var tmp [4]byte
			n := utf8.EncodeRune(tmp[:], r)
			c.buf = append(c.buf, tmp[:n]...)
			return
		}
		c.flushPendingHigh()
	}
	if ch >= 0xD800 && ch <= 0xDBFF {
		c.pendingHigh = ch
		return
	}
	if ch < 0x80 {
		c.buf = append(c.buf, byte(ch))
	} else if ch < 0x800 {
		c.buf = append(c.buf, byte(0xC0|(ch>>6)), byte(0x80|(ch&0x3F)))
	} else {
		c.buf = append(c.buf,
			byte(0xE0|(ch>>12)),
			byte(0x80|((ch>>6)&0x3F)),
			byte(0x80|(ch&0x3F)))
	}
}

func (c *jStringBuilderCore) appendString(raw string) {
	s := StrValString(raw)
	if len(s) == 0 {
		return
	}
	if c.hasCache {
		c.cachedStr = ""
		c.hasCache = false
	}
	if c.pendingHigh == 0 && len(s) < 48 {
		ascii := true
		for i := 0; i < len(s); i++ {
			if s[i] >= 0x80 {
				ascii = false
				break
			}
		}
		if ascii {
			c.buf = append(c.buf, s...)
			if c.utf16Len >= 0 {
				c.utf16Len += int32(len(s))
			}
			return
		}
	}
	if c.utf16Len >= 0 {
		c.utf16Len += JStrLen(s)
	}
	if c.pendingHigh != 0 {
		if len(s) >= 3 && s[0] == 0xED && s[1] >= 0xB0 && s[1] <= 0xBF && (s[2]&0xC0) == 0x80 {
			low := (uint16(s[0]&0x0F) << 12) | (uint16(s[1]&0x3F) << 6) | uint16(s[2]&0x3F)
			r := utf16.DecodeRune(rune(c.pendingHigh), rune(low))
			c.pendingHigh = 0
			var tmp [4]byte
			n := utf8.EncodeRune(tmp[:], r)
			c.buf = append(c.buf, tmp[:n]...)
			s = s[3:]
			if len(s) == 0 {
				return
			}
		} else {
			c.flushPendingHigh()
		}
	}
	n := len(s)
	if n >= 3 && s[n-3] == 0xED && s[n-2] >= 0xA0 && s[n-2] <= 0xAF && (s[n-1]&0xC0) == 0x80 {
		high := (uint16(s[n-3]&0x0F) << 12) | (uint16(s[n-2]&0x3F) << 6) | uint16(s[n-1]&0x3F)
		c.buf = append(c.buf, s[:n-3]...)
		c.pendingHigh = high
		return
	}
	c.buf = append(c.buf, s...)
}

func (c *jStringBuilderCore) appendASCII(s string) {
	if len(s) == 0 {
		return
	}
	if c.hasCache {
		c.cachedStr = ""
		c.hasCache = false
	}
	if c.pendingHigh != 0 {
		c.flushPendingHigh()
	}
	c.buf = append(c.buf, s...)
	if c.utf16Len >= 0 {
		c.utf16Len += int32(len(s))
	}
}

func (c *jStringBuilderCore) appendInt64(v int64) {
	if c.hasCache {
		c.cachedStr = ""
		c.hasCache = false
	}
	if c.pendingHigh != 0 {
		c.flushPendingHigh()
	}
	var tmp [20]byte
	b := strconv.AppendInt(tmp[:0], v, 10)
	c.buf = append(c.buf, b...)
	if c.utf16Len >= 0 {
		c.utf16Len += int32(len(b))
	}
}

func (c *jStringBuilderCore) appendInt32(v int32) {
	c.appendInt64(int64(v))
}

func (c *jStringBuilderCore) appendBool(v bool) {
	if v {
		c.appendASCII("true")
	} else {
		c.appendASCII("false")
	}
}

func (c *jStringBuilderCore) appendCodePoint(cp int32) {
	if cp >= 0 && cp <= 0xFFFF {
		c.appendChar(uint16(cp))
		return
	}
	if cp > 0xFFFF && cp <= 0x10FFFF {
		if c.hasCache {
			c.cachedStr = ""
			c.hasCache = false
		}
		if c.pendingHigh != 0 {
			c.flushPendingHigh()
		}
		var tmp [4]byte
		n := utf8.EncodeRune(tmp[:], cp)
		c.buf = append(c.buf, tmp[:n]...)
		if c.utf16Len >= 0 {
			c.utf16Len += 2
		}
		return
	}
	panic(javalang.NewIllegalArgumentException("Not a valid Unicode code point"))
}

// byteOffsetOfUTF16 returns the byte offset in c.buf for UTF-16 index uIdx (0 <= uIdx <= effN),
// or -1 if uIdx splits a 4-byte astral surrogate pair or encounters non-WTF-8 raw bytes.
func (c *jStringBuilderCore) byteOffsetOfUTF16(uIdx, effN int32) int {
	if uIdx == 0 {
		return 0
	}
	if uIdx == effN {
		return len(c.buf)
	}
	if int(effN) == len(c.buf) {
		return int(uIdx)
	}
	if uIdx > effN/2 {
		b := len(c.buf)
		rem := effN - uIdx
		for rem > 0 {
			if b <= 0 {
				return -1
			}
			b--
			steps := 1
			for b > 0 && steps < 4 && (c.buf[b]&0xC0) == 0x80 {
				b--
				steps++
			}
			lead := c.buf[b]
			if lead < 0x80 {
				if steps != 1 {
					return -1
				}
				rem--
			} else if lead >= 0xC2 && lead < 0xE0 && steps == 2 {
				rem--
			} else if lead >= 0xE0 && lead < 0xF0 && steps == 3 {
				rem--
			} else if lead >= 0xF0 && lead < 0xF5 && steps == 4 {
				if rem < 2 {
					return -1
				}
				rem -= 2
			} else {
				return -1
			}
		}
		return b
	}
	b := 0
	rem := uIdx
	for rem > 0 {
		if b >= len(c.buf) {
			return -1
		}
		lead := c.buf[b]
		if lead < 0x80 {
			b++
			rem--
		} else if lead >= 0xC2 && lead < 0xE0 && b+2 <= len(c.buf) {
			b += 2
			rem--
		} else if lead >= 0xE0 && lead < 0xF0 && b+3 <= len(c.buf) {
			b += 3
			rem--
		} else if lead >= 0xF0 && lead < 0xF5 && b+4 <= len(c.buf) {
			if rem < 2 {
				return -1
			}
			b += 4
			rem -= 2
		} else {
			return -1
		}
	}
	return b
}

func (c *jStringBuilderCore) charAt(idx int32) uint16 {
	n := c.length()
	if idx < 0 || idx >= n {
		panic(javalang.NewStringIndexOutOfBoundsExceptionForIndex(int(idx), int(n)))
	}
	if c.pendingHigh != 0 && idx == n-1 {
		return c.pendingHigh
	}
	effN := n
	if c.pendingHigh != 0 {
		effN--
	}
	if int(effN) == len(c.buf) {
		return uint16(c.buf[idx])
	}
	if idx == effN-1 && len(c.buf) > 0 && c.buf[len(c.buf)-1] < 0x80 {
		return uint16(c.buf[len(c.buf)-1])
	}
	off := c.byteOffsetOfUTF16(idx, effN)
	if off >= 0 && off < len(c.buf) && c.buf[off] < 0x80 {
		return uint16(c.buf[off])
	}
	return javalang.StringCharAt(c.str(), int(idx))
}

func (c *jStringBuilderCore) insertChar(idx int32, ch uint16) {
	n := c.length()
	CheckPositionIndex(idx, n)
	if idx == n {
		c.appendChar(ch)
		return
	}
	if ch < 0x80 {
		effN := n
		if c.pendingHigh != 0 {
			effN--
		}
		if idx <= effN {
			off := c.byteOffsetOfUTF16(idx, effN)
			if off >= 0 {
				c.invalidateCache()
				c.buf = append(c.buf, 0)
				copy(c.buf[off+1:], c.buf[off:len(c.buf)-1])
				c.buf[off] = byte(ch)
				c.utf16Len = n + 1
				return
			}
		}
	}
	s := c.str()
	c.resetTo(JStrSub(s, 0, idx))
	c.appendChar(ch)
	c.appendString(JStrSubFrom(s, idx))
}

func (c *jStringBuilderCore) insertString(idx int32, raw string) {
	sIns := StrValString(raw)
	n := c.length()
	CheckPositionIndex(idx, n)
	if len(sIns) == 0 {
		return
	}
	if idx == n {
		c.appendString(sIns)
		return
	}
	if len(sIns) < 64 {
		ascii := true
		for i := 0; i < len(sIns); i++ {
			if sIns[i] >= 0x80 {
				ascii = false
				break
			}
		}
		if ascii {
			effN := n
			if c.pendingHigh != 0 {
				effN--
			}
			if idx <= effN {
				off := c.byteOffsetOfUTF16(idx, effN)
				if off >= 0 {
					c.invalidateCache()
					k := len(sIns)
					oldLen := len(c.buf)
					c.buf = append(c.buf, sIns...)
					copy(c.buf[off+k:], c.buf[off:oldLen])
					copy(c.buf[off:], sIns)
					c.utf16Len = n + int32(k)
					return
				}
			}
		}
	}
	s := c.str()
	c.resetTo(JStrSub(s, 0, idx))
	c.appendString(sIns)
	c.appendString(JStrSubFrom(s, idx))
}

func (c *jStringBuilderCore) setCharAt(idx int32, ch uint16) {
	n := c.length()
	CheckElementIndex(idx, n)
	if ch < 0x80 && c.pendingHigh == 0 && int(n) == len(c.buf) {
		c.invalidateCache()
		c.buf[idx] = byte(ch)
		return
	}
	s := c.str()
	c.resetTo(JStrSub(s, 0, idx))
	c.appendChar(ch)
	c.appendString(JStrSubFrom(s, idx+1))
}

func (c *jStringBuilderCore) deleteCharAt(idx int32) {
	n := c.length()
	CheckElementIndex(idx, n)
	if c.pendingHigh != 0 && idx == n-1 {
		c.pendingHigh = 0
		c.utf16Len = n - 1
		return
	}
	if c.pendingHigh == 0 && int(n) == len(c.buf) {
		c.invalidateCache()
		copy(c.buf[idx:], c.buf[idx+1:])
		c.buf = c.buf[:len(c.buf)-1]
		c.utf16Len = n - 1
		return
	}
	s := c.str()
	c.resetTo(JStrSub(s, 0, idx))
	c.appendString(JStrSubFrom(s, idx+1))
}

func (c *jStringBuilderCore) deleteRange(start, end int32) {
	n := c.length()
	if start < 0 || start > n || start > end {
		panic(javalang.NewStringIndexOutOfBoundsExceptionForRange(int(start), int(end), int(n)))
	}
	if end > n {
		end = n
	}
	if start == end {
		return
	}
	if start == 0 && end == n {
		c.reset()
		return
	}
	if c.pendingHigh == 0 && int(n) == len(c.buf) {
		c.invalidateCache()
		copy(c.buf[start:], c.buf[end:])
		c.buf = c.buf[:len(c.buf)-int(end-start)]
		c.utf16Len = n - (end - start)
		return
	}
	s := c.str()
	c.resetTo(JStrSub(s, 0, start))
	c.appendString(JStrSubFrom(s, end))
}

func (c *jStringBuilderCore) replaceRange(start, end int32, raw string) {
	rpl := CheckNotNullString(raw)
	n := c.length()
	if start < 0 || start > n || start > end {
		panic(javalang.NewStringIndexOutOfBoundsExceptionForRange(int(start), int(end), int(n)))
	}
	if end > n {
		end = n
	}
	s := c.str()
	c.resetTo(JStrSub(s, 0, start))
	c.appendString(rpl)
	c.appendString(JStrSubFrom(s, end))
}

func (c *jStringBuilderCore) reverse() {
	n := c.length()
	if n <= 1 {
		return
	}
	if c.pendingHigh == 0 && int(n) == len(c.buf) {
		c.invalidateCache()
		for i, j := 0, len(c.buf)-1; i < j; i, j = i+1, j-1 {
			c.buf[i], c.buf[j] = c.buf[j], c.buf[i]
		}
		return
	}
	s := c.str()
	u := make([]uint16, n)
	for i := int32(0); i < n; i++ {
		u[i] = javalang.StringCharAt(s, int(i))
	}
	// Reverse UTF-16 code units while preserving valid surrogate pairs per Java StringBuilder.reverse()
	hasSurrogates := false
	for i, j := int32(0), n-1; i < j; i, j = i+1, j-1 {
		ci, cj := u[i], u[j]
		if (ci >= 0xD800 && ci <= 0xDFFF) || (cj >= 0xD800 && cj <= 0xDFFF) {
			hasSurrogates = true
		}
		u[i], u[j] = cj, ci
	}
	if len(u) == 1 && u[0] >= 0xD800 && u[0] <= 0xDFFF {
		hasSurrogates = true
	}
	if hasSurrogates {
		for i := int32(0); i < n-1; i++ {
			c2 := u[i]
			if c2 >= 0xDC00 && c2 <= 0xDFFF {
				c1 := u[i+1]
				if c1 >= 0xD800 && c1 <= 0xDBFF {
					u[i] = c1
					u[i+1] = c2
					i++
				}
			}
		}
	}
	c.resetTo(javalang.StringFromUTF16(u))
}

func (c *jStringBuilderCore) setLength(newLen int32) {
	if newLen < 0 {
		panic(javalang.NewStringIndexOutOfBoundsExceptionForIndex(int(newLen), 0))
	}
	if newLen == 0 {
		c.reset()
		return
	}
	cur := c.length()
	if newLen == cur {
		return
	}
	if newLen > cur {
		c.invalidateCache()
		if c.pendingHigh != 0 {
			c.flushPendingHigh()
		}
		add := int(newLen - cur)
		c.buf = append(c.buf, make([]byte, add)...)
		c.utf16Len = newLen
		return
	}
	if c.pendingHigh == 0 && int(cur) == len(c.buf) {
		c.invalidateCache()
		c.buf = c.buf[:newLen]
		c.utf16Len = newLen
		return
	}
	s := c.str()
	c.resetTo(JStrSub(s, 0, newLen))
}
func StringFromInt8Slice(b []int8) string {
	CheckNotNullSlice(b)
	raw := make([]byte, len(b))
	for i, v := range b {
		raw[i] = byte(v)
	}
	return string(raw)
}
func StringFromInt8Range(b []int8, off, length int32) string {
	CheckNotNullSlice(b)
	if off < 0 || length < 0 || int(off)+int(length) > len(b) {
		panic(javalang.NewStringIndexOutOfBoundsExceptionForRange(int(off), int(off+length), len(b)))
	}
	raw := make([]byte, int(length))
	for i := 0; i < int(length); i++ {
		raw[i] = byte(b[int(off)+i])
	}
	return string(raw)
}
func Arraycopy(src any, srcPos int32, dest any, destPos int32, length int32) {
	CheckNotNull(src)
	CheckNotNull(dest)
	if srcPos < 0 || destPos < 0 || length < 0 {
		panic(javalang.NewArrayIndexOutOfBoundsException("arraycopy negative index"))
	}
	sp, dp, n := int(srcPos), int(destPos), int(length)
	vs := reflect.ValueOf(src)
	vd := reflect.ValueOf(dest)
	if sp+n > vs.Len() || dp+n > vd.Len() {
		panic(javalang.NewArrayIndexOutOfBoundsException("arraycopy out of bounds"))
	}
	if n == 0 {
		return
	}
	if vs.Type() == vd.Type() {
		reflect.Copy(vd.Slice(dp, dp+n), vs.Slice(sp, sp+n))
		return
	}
	for i := 0; i < n; i++ {
		vd.Index(dp + i).Set(vs.Index(sp + i))
	}
}
func ClassPackageName(c *javalang.Class) string {
	CheckNotNullPtr(c)
	n := c.GetName()
	if idx := strings.LastIndexByte(n, '.'); idx >= 0 {
		return n[:idx]
	}
	return ""
}
func ClassResourceExists(c *javalang.Class, path string) bool {
	CheckNotNullPtr(c)
	CheckNotNullString(path)
	return shimsres.ResourceExists(ClassPackageName(c), path)
}
func ClassResourceBytes(c *javalang.Class, path string) ([]byte, bool) {
	CheckNotNullPtr(c)
	CheckNotNullString(path)
	return shimsres.GetResourceBytes(ClassPackageName(c), path)
}
func ClassForName(name string) *javalang.Class {
	CheckNotNullString(name)
	simple := name
	if idx := strings.LastIndexAny(name, ".$"); idx >= 0 {
		simple = name[idx+1:]
	}
	return ClassLiteral(name, simple)
}

func CharSeqLength(cs any) int32 {
	if cs == nil || javalang.IsNil(cs) {
		panic(javalang.NewNullPointerException())
	}
	if s, ok := cs.(string); ok {
		if s == NullString {
			panic(javalang.NewNullPointerException())
		}
		return int32(javalang.StringLength(s))
	}
	if l, ok := cs.(interface{ Length() int32 }); ok {
		return l.Length()
	}
	return int32(javalang.StringLength(StrVal(cs)))
}

func CharSeqCharAt(cs any, idx int32) uint16 {
	if cs == nil || javalang.IsNil(cs) {
		panic(javalang.NewNullPointerException())
	}
	if s, ok := cs.(string); ok {
		if s == NullString {
			panic(javalang.NewNullPointerException())
		}
		return javalang.StringCharAt(s, int(idx))
	}
	if c, ok := cs.(interface{ CharAt(int32) uint16 }); ok {
		return c.CharAt(idx)
	}
	return javalang.StringCharAt(StrVal(cs), int(idx))
}

func CompareValues(a, b any) int32 {
	a = CheckNotNull(ToAny(a))
	b = CheckNotNull(ToAny(b))
	if sa, ok := a.(string); ok {
		return int32(javalang.StringCompareTo(sa, b.(string)))
	}
	if ea, ok := a.(EnumAccessor); ok {
		return ea.GetEnumBase().Ordinal_ - b.(EnumAccessor).GetEnumBase().Ordinal_
	}
	if ia, ok := a.(*BoxInt); ok {
		ib := b.(*BoxInt)
		if ia.Val < ib.Val {
			return -1
		}
		if ia.Val > ib.Val {
			return 1
		}
		return 0
	}
	if la, ok := a.(*BoxLong); ok {
		lb := b.(*BoxLong)
		if la.Val < lb.Val {
			return -1
		}
		if la.Val > lb.Val {
			return 1
		}
		return 0
	}
	if da, ok := a.(*BoxDouble); ok {
		db := b.(*BoxDouble)
		return int32(javalang.DoubleCompare(da.Val, db.Val))
	}
	if fa, ok := a.(*BoxFloat); ok {
		fb := b.(*BoxFloat)
		return int32(javalang.DoubleCompare(float64(fa.Val), float64(fb.Val)))
	}
	if ba, ok := a.(*BoxBoolean); ok {
		bb := b.(*BoxBoolean)
		if ba.Val == bb.Val {
			return 0
		}
		if ba.Val {
			return 1
		}
		return -1
	}
	if cha, ok := a.(*BoxChar); ok {
		chb := b.(*BoxChar)
		return int32(cha.Val) - int32(chb.Val)
	}
	if ca, ok := a.(interface{ CompareTo(any) int32 }); ok {
		return ca.CompareTo(b)
	}
	panic(javalang.NewClassCastException(fmt.Sprintf("%T is not Comparable", a)))
}

// --- Math & Numeric Intrinsics ---

func ParseInt(s string, radix int32) int32 {
	CheckNotNullString(s)
	v, err := strconv.ParseInt(strings.TrimPrefix(s, "+"), int(radix), 32)
	if err != nil {
		panic(javalang.NewNumberFormatException(s))
	}
	return int32(v)
}

func ParseLong(s string, radix int32) int64 {
	CheckNotNullString(s)
	v, err := strconv.ParseInt(strings.TrimPrefix(s, "+"), int(radix), 64)
	if err != nil {
		panic(javalang.NewNumberFormatException(s))
	}
	return v
}

func ParseDouble(s string) float64 {
	CheckNotNullString(s)
	return javalang.DoubleParseDouble(s)
}

func IntBitCount(x int32) int32      { return int32(bits.OnesCount32(uint32(x))) }
func IntLeadingZeros(x int32) int32  { return int32(bits.LeadingZeros32(uint32(x))) }
func IntTrailingZeros(x int32) int32 { return int32(bits.TrailingZeros32(uint32(x))) }
func IntLowestOneBit(x int32) int32  { return x & (-x) }
func IntHighestOneBit(x int32) int32 {
	if x == 0 {
		return 0
	}
	return int32(uint32(1) << (31 - bits.LeadingZeros32(uint32(x))))
}
func LongBitCount(x int64) int32      { return int32(bits.OnesCount64(uint64(x))) }
func LongLeadingZeros(x int64) int32  { return int32(bits.LeadingZeros64(uint64(x))) }
func LongTrailingZeros(x int64) int32 { return int32(bits.TrailingZeros64(uint64(x))) }
func LongLowestOneBit(x int64) int64  { return x & (-x) }
func LongHighestOneBit(x int64) int64 {
	if x == 0 {
		return 0
	}
	return int64(uint64(1) << (63 - bits.LeadingZeros64(uint64(x))))
}

func StringHash(s string) int32 {
	if s == NullString {
		return 0
	}
	return javalang.StringHashCode(s)
}

func BoolHash(b bool) int32 {
	if b {
		return 1231
	}
	return 1237
}

func LongHash(x int64) int32 {
	u := uint64(x)
	return int32(u ^ (u >> 32))
}

type IdentityHasher interface {
	IdentityHash_() int32
}

var nextIdentHash_ atomic.Uint32

//go:noinline
func AssignIdentHash(p *atomic.Int32) int32 {
	for {
		h := int32(nextIdentHash_.Add(0x9e3779b9))
		if h == 0 {
			continue
		}
		if p.CompareAndSwap(0, h) {
			return h
		}
		return p.Load()
	}
}

func IdentityHash(v any) int32 {
	if v == nil {
		return 0
	}
	if ea, ok := v.(EnumAccessor); ok {
		eb := ea.GetEnumBase()
		if eb == nil {
			return 0
		}
		return eb.HashCode()
	}
	if nc, ok := v.(jJavaNilChecker); ok && nc.IsJavaNil_() {
		return 0
	}
	if ih, ok := v.(IdentityHasher); ok {
		return ih.IdentityHash_()
	}
	return javalang.IdentityHashCode(v)
}

func ObjectHash(v any) int32 {
	if v == nil {
		return 0
	}
	if fk, ok := v.(FastCustomKey); ok {
		if fk.IsJavaNil_() {
			return 0
		}
		return fk.HashCode()
	}
	if _, hasHC := v.(jHashCoder); !hasHC {
		if ih, ok := v.(IdentityHasher); ok {
			if nc, ok2 := v.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
				return 0
			}
			return ih.IdentityHash_()
		}
	}
	if isJavaNilAny(v) {
		return 0
	}
	return javalang.ObjectHashCode(v)
}

func ObjectsHash(vals ...any) int32 {
	if vals == nil {
		return 0
	}
	h := int32(1)
	for _, v := range vals {
		h = 31*h + ObjectHash(ToAny(v))
	}
	return h
}

// --- Deterministic Slim Zero-Bloat Collections Core ---

type JListCore struct {
	Items   []any
	capHint int32
}

func NewJListCore(items ...any) *JListCore {
	c := &JListCore{}
	n := len(items)
	if n == 0 {
		return c
	}
	s := make([]any, n)
	for i, v := range items {
		s[i] = ToAny(v)
	}
	c.Items = s
	return c
}

func (l *JListCore) InitFromSlice(items []any) {
	n := len(items)
	if n == 0 {
		l.Items = nil
		return
	}
	out := make([]any, n)
	copy(out, items)
	l.Items = out
}

func (l *JListCore) InitOwnSlice(items []any) {
	if len(items) == 0 {
		l.Items = nil
		return
	}
	l.Items = items
}

func (l *JListCore) ToArray() []any {
	if len(l.Items) == 0 {
		return []any{}
	}
	return CloneSlice(l.Items)
}

func (l *JListCore) InitCap(initialCap int32) {
	if initialCap < 0 {
		panic(javalang.NewIllegalArgumentException("Illegal Capacity: " + fmt.Sprint(initialCap)))
	}
	if initialCap <= 0 {
		return
	}
	if initialCap > 1024 {
		initialCap = 1024
	}
	if l.Items == nil {
		if initialCap > l.capHint {
			l.capHint = initialCap
		}
		return
	}
	if int(initialCap) > cap(l.Items) {
		next := make([]any, len(l.Items), initialCap)
		copy(next, l.Items)
		l.Items = next
	}
}

func (l *JListCore) Size() int32   { return int32(len(l.Items)) }
func (l *JListCore) IsEmpty() bool { return len(l.Items) == 0 }
func (l *JListCore) Get(i int32) any {
	if i < 0 || int(i) >= len(l.Items) {
		panic(javalang.NewIndexOutOfBoundsExceptionForIndex(int(i), len(l.Items)))
	}
	return l.Items[i]
}
func (l *JListCore) Set(i int32, v any) any {
	old := l.Get(i)
	l.Items[i] = ToAny(v)
	return old
}
func (l *JListCore) Add(v any) bool {
	v = ToAny(v)
	if cap(l.Items) == 0 {
		c := int(l.capHint)
		if c < 4 {
			c = 4
		}
		l.Items = make([]any, 1, c)
		l.Items[0] = v
		return true
	}
	l.Items = append(l.Items, v)
	return true
}
func (l *JListCore) AddAt(i int32, v any) {
	if i < 0 || int(i) > len(l.Items) {
		panic(javalang.NewIndexOutOfBoundsExceptionForIndex(int(i), len(l.Items)))
	}
	v = ToAny(v)
	if cap(l.Items) == 0 {
		c := int(l.capHint)
		if c < 4 {
			c = 4
		}
		l.Items = make([]any, 1, c)
		l.Items[0] = v
		return
	}
	l.Items = append(l.Items, nil)
	copy(l.Items[int(i)+1:], l.Items[int(i):])
	l.Items[i] = v
}
func (l *JListCore) RemoveAt(i int32) any {
	old := l.Get(i)
	n := len(l.Items)
	if n == 1 {
		l.Items[0] = nil
		l.Items = l.Items[:0]
		return old
	}
	if i == 0 && n > 8 {
		l.Items[0] = nil
		l.Items = l.Items[1:]
		if cap(l.Items) > 64 && len(l.Items)*4 <= cap(l.Items) {
			next := make([]any, len(l.Items), len(l.Items)*2)
			copy(next, l.Items)
			l.Items = next
		}
		return old
	}
	last := n - 1
	copy(l.Items[int(i):], l.Items[int(i)+1:])
	l.Items[last] = nil
	l.Items = l.Items[:last]
	return old
}
func (l *JListCore) RemoveVal(v any) bool {
	v = ToAny(v)
	for i, x := range l.Items {
		if javalang.ObjectsEquals(x, v) {
			l.RemoveAt(int32(i))
			return true
		}
	}
	return false
}
func (l *JListCore) Contains(v any) bool {
	return l.IndexOf(v) >= 0
}
func (l *JListCore) IndexOf(v any) int32 {
	v = ToAny(v)
	for i, x := range l.Items {
		if javalang.ObjectsEquals(x, v) {
			return int32(i)
		}
	}
	return -1
}
func (l *JListCore) LastIndexOf(v any) int32 {
	v = ToAny(v)
	for i := len(l.Items) - 1; i >= 0; i-- {
		if javalang.ObjectsEquals(l.Items[i], v) {
			return int32(i)
		}
	}
	return -1
}
func (l *JListCore) Clear() {
	if l.Items != nil {
		clear(l.Items)
		l.Items = l.Items[:0]
	}
}

// JDequeCore is an O(1) power-of-2 circular ring buffer backing java.util.ArrayDeque,
// with a homogeneous noscan `nodes []uint32` fast mode for AST Node worklists.
type JDequeCore struct {
	buf      []any
	nodes    []uint32
	colorIds []int64
	head     int32
	mask     int32 // int32(len(buf)-1), int32(len(nodes)-1), or int32(len(colorIds)-1)
	size     int32
	capHint  int32
}

func checkDequeElem(v any) any {
	if v == nil {
		ThrowNPE()
	}
	if nc, ok := v.(jJavaNilChecker); ok {
		if nc.IsJavaNil_() {
			ThrowNPE()
		}
		return v
	}
	if s, ok := v.(string); ok {
		if s == NullString {
			ThrowNPE()
		}
		return v
	}
	v = ToAny(v)
	if v == nil {
		ThrowNPE()
	}
	return v
}

func extractDequeNodeId(v any) (uint32, bool) {
	if rn, ok := v.(jRawNodeProvider); ok {
		if isJavaNilAny(v) {
			return 0, false
		}
		u := rn.RawNodeId_()
		return u, u != 0
	}
	if n, ok := v.(nodestore.Node); ok {
		return uint32(n), n != 0
	}
	if u, ok := v.(uint32); ok {
		return u, u != 0
	}
	return 0, false
}

func (d *JDequeCore) InitCap(initialCap int32) {
	if initialCap <= 0 {
		return
	}
	c := int32(8)
	limit := initialCap
	if limit > (1 << 20) {
		limit = 1 << 20
	}
	for c < limit {
		c <<= 1
	}
	if d.buf == nil && d.nodes == nil && d.colorIds == nil {
		if c > d.capHint {
			d.capHint = c
		}
	}
}

func (d *JDequeCore) growColorIds() {
	if d.colorIds == nil {
		c := d.capHint
		if c < 8 {
			c = 8
		}
		d.colorIds = make([]int64, c)
		d.mask = c - 1
		return
	}
	oldCap := int32(len(d.colorIds))
	newCap := oldCap << 1
	nb := make([]int64, newCap)
	firstPart := oldCap - d.head
	copy(nb[:firstPart], d.colorIds[d.head:])
	if d.head > 0 {
		copy(nb[firstPart:], d.colorIds[:d.head])
	}
	d.colorIds = nb
	d.head = 0
	d.mask = newCap - 1
}

func (d *JDequeCore) upgradeColorIdsToBuf() {
	if d.size == 0 {
		if int32(len(d.colorIds)) > d.capHint {
			d.capHint = int32(len(d.colorIds))
		}
		d.colorIds = nil
		d.head = 0
		d.mask = 0
		return
	}
	capLen := int32(len(d.colorIds))
	nb := make([]any, capLen)
	for i := int32(0); i < d.size; i++ {
		nb[i] = boxColorId(d.colorIds[(d.head+i)&d.mask])
	}
	d.colorIds = nil
	d.buf = nb
	d.head = 0
	d.mask = capLen - 1
}

func (d *JDequeCore) grow() {
	if d.buf == nil {
		c := d.capHint
		if c < 8 {
			c = 8
		}
		d.buf = make([]any, c)
		d.mask = c - 1
		return
	}
	oldCap := int32(len(d.buf))
	newCap := oldCap << 1
	nb := make([]any, newCap)
	firstPart := oldCap - d.head
	copy(nb[:firstPart], d.buf[d.head:])
	if d.head > 0 {
		copy(nb[firstPart:], d.buf[:d.head])
	}
	d.buf = nb
	d.head = 0
	d.mask = newCap - 1
}

func (d *JDequeCore) growNodes() {
	if d.nodes == nil {
		c := d.capHint
		if c < 8 {
			c = 8
		}
		d.nodes = make([]uint32, c)
		d.mask = c - 1
		return
	}
	oldCap := int32(len(d.nodes))
	newCap := oldCap << 1
	nb := make([]uint32, newCap)
	firstPart := oldCap - d.head
	copy(nb[:firstPart], d.nodes[d.head:])
	if d.head > 0 {
		copy(nb[firstPart:], d.nodes[:d.head])
	}
	d.nodes = nb
	d.head = 0
	d.mask = newCap - 1
}

func (d *JDequeCore) upgradeNodesToBuf() {
	if d.size == 0 {
		if int32(len(d.nodes)) > d.capHint {
			d.capHint = int32(len(d.nodes))
		}
		d.nodes = nil
		d.head = 0
		d.mask = 0
		return
	}
	capLen := int32(len(d.nodes))
	nb := make([]any, capLen)
	for i := int32(0); i < d.size; i++ {
		nb[i] = boxNodeById(d.nodes[(d.head+i)&d.mask])
	}
	d.nodes = nil
	d.buf = nb
	d.head = 0
	d.mask = capLen - 1
}

func (d *JDequeCore) elemAt(idx int32) any {
	if d.nodes != nil {
		return boxNodeById(d.nodes[idx])
	}
	if d.colorIds != nil {
		return boxColorId(d.colorIds[idx])
	}
	return d.buf[idx]
}

func (d *JDequeCore) Size() int32   { return d.size }
func (d *JDequeCore) IsEmpty() bool { return d.size == 0 }

func (d *JDequeCore) AddFirstColorId(raw int64) {
	if d.nodes != nil {
		d.upgradeNodesToBuf()
	}
	if d.buf != nil {
		if int(d.size) == len(d.buf) {
			d.grow()
		}
		d.head = (d.head - 1) & d.mask
		d.buf[d.head] = boxColorId(raw)
		d.size++
		return
	}
	if int(d.size) == len(d.colorIds) {
		d.growColorIds()
	}
	d.head = (d.head - 1) & d.mask
	d.colorIds[d.head] = raw
	d.size++
}

func (d *JDequeCore) AddLastColorId(raw int64) {
	if d.nodes != nil {
		d.upgradeNodesToBuf()
	}
	if d.buf != nil {
		if int(d.size) == len(d.buf) {
			d.grow()
		}
		d.buf[(d.head+d.size)&d.mask] = boxColorId(raw)
		d.size++
		return
	}
	if int(d.size) == len(d.colorIds) {
		d.growColorIds()
	}
	d.colorIds[(d.head+d.size)&d.mask] = raw
	d.size++
}

func (d *JDequeCore) RemoveFirstVoid() {
	if d.size == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	if d.nodes != nil {
		d.nodes[d.head] = 0
		d.head = (d.head + 1) & d.mask
		d.size--
		return
	}
	if d.colorIds != nil {
		d.head = (d.head + 1) & d.mask
		d.size--
		return
	}
	d.buf[d.head] = nil
	d.head = (d.head + 1) & d.mask
	d.size--
}

func (d *JDequeCore) RemoveLastVoid() {
	if d.size == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	idx := (d.head + d.size - 1) & d.mask
	if d.nodes != nil {
		d.nodes[idx] = 0
		d.size--
		return
	}
	if d.colorIds != nil {
		d.size--
		return
	}
	d.buf[idx] = nil
	d.size--
}

func (d *JDequeCore) AddFirstNode(n uint32) {
	if n == 0 {
		ThrowNPE()
	}
	if d.colorIds != nil {
		d.upgradeColorIdsToBuf()
	}
	if d.buf != nil {
		if int(d.size) == len(d.buf) {
			d.grow()
		}
		d.head = (d.head - 1) & d.mask
		d.buf[d.head] = boxNodeById(n)
		d.size++
		return
	}
	if int(d.size) == len(d.nodes) {
		d.growNodes()
	}
	d.head = (d.head - 1) & d.mask
	d.nodes[d.head] = n
	d.size++
}

func (d *JDequeCore) AddLastNode(n uint32) {
	if n == 0 {
		ThrowNPE()
	}
	if d.colorIds != nil {
		d.upgradeColorIdsToBuf()
	}
	if d.buf != nil {
		if int(d.size) == len(d.buf) {
			d.grow()
		}
		d.buf[(d.head+d.size)&d.mask] = boxNodeById(n)
		d.size++
		return
	}
	if int(d.size) == len(d.nodes) {
		d.growNodes()
	}
	d.nodes[(d.head+d.size)&d.mask] = n
	d.size++
}

func (d *JDequeCore) PopNode() uint32 {
	if d.size == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	if d.nodes != nil {
		u := d.nodes[d.head]
		d.nodes[d.head] = 0
		d.head = (d.head + 1) & d.mask
		d.size--
		return u
	}
	return CastNode[uint32](d.RemoveFirst())
}

func (d *JDequeCore) RemoveLastNode() uint32 {
	if d.size == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	if d.nodes != nil {
		idx := (d.head + d.size - 1) & d.mask
		u := d.nodes[idx]
		d.nodes[idx] = 0
		d.size--
		return u
	}
	return CastNode[uint32](d.RemoveLast())
}

func (d *JDequeCore) PollFirstNode() uint32 {
	if d.size == 0 {
		return 0
	}
	if d.nodes != nil {
		u := d.nodes[d.head]
		d.nodes[d.head] = 0
		d.head = (d.head + 1) & d.mask
		d.size--
		return u
	}
	return CastNode[uint32](d.PollFirst())
}

func (d *JDequeCore) PeekFirstNode() uint32 {
	if d.size == 0 {
		return 0
	}
	if d.nodes != nil {
		return d.nodes[d.head]
	}
	return CastNode[uint32](d.PeekFirst())
}

func (d *JDequeCore) PeekLastNode() uint32 {
	if d.size == 0 {
		return 0
	}
	if d.nodes != nil {
		return d.nodes[(d.head+d.size-1)&d.mask]
	}
	return CastNode[uint32](d.PeekLast())
}

func (d *JDequeCore) AddLast(v any) {
	v = checkDequeElem(v)
	if d.buf == nil {
		if d.colorIds == nil {
			if u, ok := extractDequeNodeId(v); ok {
				d.AddLastNode(u)
				return
			}
		}
		if d.nodes == nil {
			if rc, ok := v.(jRawColorIdProvider); ok {
				if raw, valid := rc.RawColorId_(); valid {
					d.AddLastColorId(raw)
					return
				}
			}
		}
		if d.nodes != nil {
			d.upgradeNodesToBuf()
		}
		if d.colorIds != nil {
			d.upgradeColorIdsToBuf()
		}
	}
	if int(d.size) == len(d.buf) {
		d.grow()
	}
	d.buf[(d.head+d.size)&d.mask] = v
	d.size++
}

func (d *JDequeCore) OfferLast(v any) bool {
	d.AddLast(v)
	return true
}

func (d *JDequeCore) Add(v any) bool {
	d.AddLast(v)
	return true
}

func (d *JDequeCore) AddFirst(v any) {
	v = checkDequeElem(v)
	if d.buf == nil {
		if d.colorIds == nil {
			if u, ok := extractDequeNodeId(v); ok {
				d.AddFirstNode(u)
				return
			}
		}
		if d.nodes == nil {
			if rc, ok := v.(jRawColorIdProvider); ok {
				if raw, valid := rc.RawColorId_(); valid {
					d.AddFirstColorId(raw)
					return
				}
			}
		}
		if d.nodes != nil {
			d.upgradeNodesToBuf()
		}
		if d.colorIds != nil {
			d.upgradeColorIdsToBuf()
		}
	}
	if int(d.size) == len(d.buf) {
		d.grow()
	}
	d.head = (d.head - 1) & d.mask
	d.buf[d.head] = v
	d.size++
}

func (d *JDequeCore) OfferFirst(v any) bool {
	d.AddFirst(v)
	return true
}

func (d *JDequeCore) Push(v any) {
	d.AddFirst(v)
}

func (d *JDequeCore) PollFirst() any {
	if d.size == 0 {
		return nil
	}
	if d.nodes != nil {
		u := d.nodes[d.head]
		d.nodes[d.head] = 0
		d.head = (d.head + 1) & d.mask
		d.size--
		return boxNodeById(u)
	}
	if d.colorIds != nil {
		raw := d.colorIds[d.head]
		d.head = (d.head + 1) & d.mask
		d.size--
		return boxColorId(raw)
	}
	v := d.buf[d.head]
	d.buf[d.head] = nil
	d.head = (d.head + 1) & d.mask
	d.size--
	return v
}

func (d *JDequeCore) RemoveFirst() any {
	if d.size == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	if d.nodes != nil {
		u := d.nodes[d.head]
		d.nodes[d.head] = 0
		d.head = (d.head + 1) & d.mask
		d.size--
		return boxNodeById(u)
	}
	if d.colorIds != nil {
		raw := d.colorIds[d.head]
		d.head = (d.head + 1) & d.mask
		d.size--
		return boxColorId(raw)
	}
	v := d.buf[d.head]
	d.buf[d.head] = nil
	d.head = (d.head + 1) & d.mask
	d.size--
	return v
}

func (d *JDequeCore) Pop() any {
	return d.RemoveFirst()
}

func (d *JDequeCore) PollLast() any {
	if d.size == 0 {
		return nil
	}
	idx := (d.head + d.size - 1) & d.mask
	if d.nodes != nil {
		u := d.nodes[idx]
		d.nodes[idx] = 0
		d.size--
		return boxNodeById(u)
	}
	if d.colorIds != nil {
		raw := d.colorIds[idx]
		d.size--
		return boxColorId(raw)
	}
	v := d.buf[idx]
	d.buf[idx] = nil
	d.size--
	return v
}

func (d *JDequeCore) RemoveLast() any {
	if d.size == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	idx := (d.head + d.size - 1) & d.mask
	if d.nodes != nil {
		u := d.nodes[idx]
		d.nodes[idx] = 0
		d.size--
		return boxNodeById(u)
	}
	if d.colorIds != nil {
		raw := d.colorIds[idx]
		d.size--
		return boxColorId(raw)
	}
	v := d.buf[idx]
	d.buf[idx] = nil
	d.size--
	return v
}

func (d *JDequeCore) PeekFirst() any {
	if d.size == 0 {
		return nil
	}
	return d.elemAt(d.head)
}

func (d *JDequeCore) GetFirst() any {
	if d.size == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return d.elemAt(d.head)
}

func (d *JDequeCore) PeekLast() any {
	if d.size == 0 {
		return nil
	}
	return d.elemAt((d.head + d.size - 1) & d.mask)
}

func (d *JDequeCore) GetLast() any {
	if d.size == 0 {
		panic(javalang.NewNoSuchElementException("empty"))
	}
	return d.elemAt((d.head + d.size - 1) & d.mask)
}

func (d *JDequeCore) Clear() {
	if d.nodes != nil || d.colorIds != nil {
		d.head = 0
		d.size = 0
		return
	}
	for i := int32(0); i < d.size; i++ {
		d.buf[(d.head+i)&d.mask] = nil
	}
	d.head = 0
	d.size = 0
}

func (d *JDequeCore) ForEach(fn func(any)) {
	if d.nodes != nil {
		for i := int32(0); i < d.size; i++ {
			fn(boxNodeById(d.nodes[(d.head+i)&d.mask]))
		}
		return
	}
	if d.colorIds != nil {
		for i := int32(0); i < d.size; i++ {
			fn(boxColorId(d.colorIds[(d.head+i)&d.mask]))
		}
		return
	}
	for i := int32(0); i < d.size; i++ {
		fn(d.buf[(d.head+i)&d.mask])
	}
}

func (d *JDequeCore) RemoveIf(pred func(any) bool) bool {
	if d.size == 0 {
		return false
	}
	if d.colorIds != nil {
		d.upgradeColorIdsToBuf()
	}
	oldSize := d.size
	writeOff := int32(0)
	if d.nodes != nil {
		for readOff := int32(0); readOff < oldSize; readOff++ {
			u := d.nodes[(d.head+readOff)&d.mask]
			if !pred(boxNodeById(u)) {
				if writeOff != readOff {
					d.nodes[(d.head+writeOff)&d.mask] = u
				}
				writeOff++
			}
		}
		if writeOff == oldSize {
			return false
		}
		d.size = writeOff
		return true
	}
	for readOff := int32(0); readOff < oldSize; readOff++ {
		v := d.buf[(d.head+readOff)&d.mask]
		if !pred(v) {
			if writeOff != readOff {
				d.buf[(d.head+writeOff)&d.mask] = v
			}
			writeOff++
		}
	}
	if writeOff == oldSize {
		return false
	}
	for i := writeOff; i < oldSize; i++ {
		d.buf[(d.head+i)&d.mask] = nil
	}
	d.size = writeOff
	return true
}

func (d *JDequeCore) deleteAtOffset(off int32) {
	if d.colorIds != nil {
		d.upgradeColorIdsToBuf()
	}
	if d.nodes != nil {
		if off < (d.size >> 1) {
			for j := off; j > 0; j-- {
				d.nodes[(d.head+j)&d.mask] = d.nodes[(d.head+j-1)&d.mask]
			}
			d.nodes[d.head] = 0
			d.head = (d.head + 1) & d.mask
			d.size--
			return
		}
		last := d.size - 1
		for j := off; j < last; j++ {
			d.nodes[(d.head+j)&d.mask] = d.nodes[(d.head+j+1)&d.mask]
		}
		d.nodes[(d.head+last)&d.mask] = 0
		d.size--
		return
	}
	if off < (d.size >> 1) {
		for j := off; j > 0; j-- {
			d.buf[(d.head+j)&d.mask] = d.buf[(d.head+j-1)&d.mask]
		}
		d.buf[d.head] = nil
		d.head = (d.head + 1) & d.mask
		d.size--
		return
	}
	last := d.size - 1
	for j := off; j < last; j++ {
		d.buf[(d.head+j)&d.mask] = d.buf[(d.head+j+1)&d.mask]
	}
	d.buf[(d.head+last)&d.mask] = nil
	d.size--
}

func (d *JDequeCore) Contains(v any) bool {
	if v == nil || d.size == 0 {
		return false
	}
	if d.colorIds != nil {
		rc, ok := v.(jRawColorIdProvider)
		if !ok {
			return false
		}
		raw, valid := rc.RawColorId_()
		if !valid {
			return false
		}
		for i := int32(0); i < d.size; i++ {
			if d.colorIds[(d.head+i)&d.mask] == raw {
				return true
			}
		}
		return false
	}
	if d.nodes != nil {
		u, ok := extractDequeNodeId(v)
		if !ok {
			return false
		}
		for i := int32(0); i < d.size; i++ {
			if d.nodes[(d.head+i)&d.mask] == u {
				return true
			}
		}
		return false
	}
	v = ToAny(v)
	if v == nil {
		return false
	}
	for i := int32(0); i < d.size; i++ {
		if javalang.ObjectsEquals(d.buf[(d.head+i)&d.mask], v) {
			return true
		}
	}
	return false
}

func (d *JDequeCore) RemoveFirstOccurrence(v any) bool {
	if v == nil || d.size == 0 {
		return false
	}
	if d.colorIds != nil {
		d.upgradeColorIdsToBuf()
	}
	if d.nodes != nil {
		u, ok := extractDequeNodeId(v)
		if !ok {
			return false
		}
		for i := int32(0); i < d.size; i++ {
			if d.nodes[(d.head+i)&d.mask] == u {
				d.deleteAtOffset(i)
				return true
			}
		}
		return false
	}
	v = ToAny(v)
	if v == nil {
		return false
	}
	for i := int32(0); i < d.size; i++ {
		if javalang.ObjectsEquals(d.buf[(d.head+i)&d.mask], v) {
			d.deleteAtOffset(i)
			return true
		}
	}
	return false
}

func (d *JDequeCore) RemoveLastOccurrence(v any) bool {
	if v == nil || d.size == 0 {
		return false
	}
	if d.colorIds != nil {
		d.upgradeColorIdsToBuf()
	}
	if d.nodes != nil {
		u, ok := extractDequeNodeId(v)
		if !ok {
			return false
		}
		for i := d.size - 1; i >= 0; i-- {
			if d.nodes[(d.head+i)&d.mask] == u {
				d.deleteAtOffset(i)
				return true
			}
		}
		return false
	}
	v = ToAny(v)
	if v == nil {
		return false
	}
	for i := d.size - 1; i >= 0; i-- {
		if javalang.ObjectsEquals(d.buf[(d.head+i)&d.mask], v) {
			d.deleteAtOffset(i)
			return true
		}
	}
	return false
}

func (d *JDequeCore) CloneInto(dst *JDequeCore) {
	if d.size == 0 {
		return
	}
	if d.colorIds != nil {
		c := int32(len(d.colorIds))
		dst.colorIds = make([]int64, c)
		dst.mask = c - 1
		for i := int32(0); i < d.size; i++ {
			dst.colorIds[i] = d.colorIds[(d.head+i)&d.mask]
		}
		dst.size = d.size
		return
	}
	if d.nodes != nil {
		c := int32(len(d.nodes))
		dst.nodes = make([]uint32, c)
		dst.mask = c - 1
		for i := int32(0); i < d.size; i++ {
			dst.nodes[i] = d.nodes[(d.head+i)&d.mask]
		}
		dst.size = d.size
		return
	}
	c := int32(len(d.buf))
	dst.buf = make([]any, c)
	dst.mask = c - 1
	for i := int32(0); i < d.size; i++ {
		dst.buf[i] = d.buf[(d.head+i)&d.mask]
	}
	dst.size = d.size
}

func (d *JDequeCore) ToSlice() []any {
	if d.size == 0 {
		return []any{}
	}
	out := make([]any, d.size)
	if d.colorIds != nil {
		for i := int32(0); i < d.size; i++ {
			out[i] = boxColorId(d.colorIds[(d.head+i)&d.mask])
		}
		return out
	}
	if d.nodes != nil {
		for i := int32(0); i < d.size; i++ {
			out[i] = boxNodeById(d.nodes[(d.head+i)&d.mask])
		}
		return out
	}
	capLen := int32(len(d.buf))
	if d.head+d.size <= capLen {
		copy(out, d.buf[d.head:d.head+d.size])
	} else {
		firstPart := capLen - d.head
		copy(out[:firstPart], d.buf[d.head:])
		copy(out[firstPart:], d.buf[:d.size-firstPart])
	}
	return out
}

const (
	jTombstoneStr         = "\xff\x00j2go_tomb"
	jTombstoneNode uint32 = 0xFFFFFFFF

	jSetModeEmpty uint8 = 0
	jSetModeNode  uint8 = 1
	jSetModeStr   uint8 = 2
	jSetModeAny   uint8 = 3

	jSmallLimit     int = 8
	jStrSmallLimit  int = 13
	jNodeSmallLimit int = 16
)

type jTombstoneType struct{}

var (
	jTombstoneObj  any      = jTombstoneType{}
	emptyAnySlice_ []any    = []any{}
	emptyStrSlice_ []string = []string{}
)

type jIdentityKey interface {
	HasIdentityEquals_()
	IsJavaNil_() bool
}

type FastCustomKey interface {
	HashCode() int32
	Equals(any) bool
	IsJavaNil_() bool
}

type jSourceObjPtr interface {
	BindSelf_(any)
	IsJavaNil_() bool
}

type jHashCoder interface {
	HashCode() int32
}

type jEqualer interface {
	Equals(any) bool
}

func smearHash(h int32) uint32 {
	u := uint32(h)
	return u ^ (u >> 16)
}

func elemJavaHashCode(v any, identity bool) int32 {
	if v == nil {
		return 0
	}
	if ea, ok := v.(EnumAccessor); ok {
		eb := ea.GetEnumBase()
		if eb == nil {
			return 0
		}
		return eb.HashCode()
	}
	if identity {
		if nc, ok := v.(jJavaNilChecker); ok && nc.IsJavaNil_() {
			return 0
		}
		if ih, ok := v.(IdentityHasher); ok {
			return ih.IdentityHash_()
		}
		return javalang.IdentityHashCode(v)
	}
	if ik, ok := v.(jIdentityKey); ok {
		if ik.IsJavaNil_() {
			return 0
		}
		if ih, ok := v.(IdentityHasher); ok {
			return ih.IdentityHash_()
		}
		return javalang.IdentityHashCode(v)
	}
	if rn, ok := v.(jRawNodeProvider); ok {
		return int32(rn.RawNodeId_())
	}
	if fk, ok := v.(FastCustomKey); ok {
		if fk.IsJavaNil_() {
			return 0
		}
		return fk.HashCode()
	}
	if s, ok := v.(string); ok {
		if s == NullString {
			return 0
		}
		return javalang.StringHashCode(s)
	}
	if nd, ok := v.(nodestore.Node); ok {
		return int32(nd)
	}
	if _, hasHC := v.(jHashCoder); !hasHC {
		if ih, ok := v.(IdentityHasher); ok {
			if nc, ok2 := v.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
				return 0
			}
			return ih.IdentityHash_()
		}
	}
	return javalang.ObjectHashCode(v)
}

func jSetElemEquals(a, b any) bool {
	if b == nil {
		if a == nil {
			return true
		}
		if sa, ok := a.(string); ok {
			return sa == NullString
		}
		if na, ok := a.(jJavaNilChecker); ok {
			return na.IsJavaNil_()
		}
		return isJavaNilAny(a)
	}
	if sb, ok := b.(string); ok {
		if sb == NullString {
			return a == nil || isJavaNilAny(a)
		}
		sa, ok2 := a.(string)
		return ok2 && sa == sb
	}
	if ib, ok := b.(jIdentityKey); ok {
		if ib.IsJavaNil_() {
			return a == nil || isJavaNilAny(a)
		}
		return a == b
	}
	if _, ok := b.(IdentityComparable); ok {
		if nb, ok2 := b.(jJavaNilChecker); ok2 && nb.IsJavaNil_() {
			return a == nil || isJavaNilAny(a)
		}
		return a == b
	}
	if _, ok := b.(EnumAccessor); ok {
		if nb, ok2 := b.(jJavaNilChecker); ok2 && nb.IsJavaNil_() {
			return a == nil || isJavaNilAny(a)
		}
		return a == b
	}
	if rnb, ok := b.(jRawNodeProvider); ok {
		ub := rnb.RawNodeId_()
		if ub == 0 {
			return a == nil || isJavaNilAny(a)
		}
		if rna, ok2 := a.(jRawNodeProvider); ok2 {
			return rna.RawNodeId_() == ub
		}
		if nda, ok2 := a.(nodestore.Node); ok2 {
			return uint32(nda) == ub
		}
		return false
	}
	if ib, ok := b.(*BoxInt); ok {
		if ib == nil {
			return a == nil || isJavaNilAny(a)
		}
		if ia, ok2 := a.(*BoxInt); ok2 {
			return ia != nil && ia.Val == ib.Val
		}
		if iv, ok2 := a.(int32); ok2 {
			return iv == ib.Val
		}
		return false
	}
	if ivb, ok := b.(int32); ok {
		if iva, ok2 := a.(int32); ok2 {
			return iva == ivb
		}
		if ia, ok2 := a.(*BoxInt); ok2 {
			return ia != nil && ia.Val == ivb
		}
		return false
	}
	if fkb, ok := b.(FastCustomKey); ok {
		if fkb.IsJavaNil_() {
			return a == nil || isJavaNilAny(a)
		}
		if a == nil || isJavaNilAny(a) {
			return false
		}
		return a == b || fkb.Equals(a)
	}
	if a == nil {
		return isJavaNilAny(b)
	}
	if ia, ok := a.(jIdentityKey); ok {
		if ia.IsJavaNil_() {
			return isJavaNilAny(b)
		}
		return a == b
	}
	if _, ok := a.(IdentityComparable); ok {
		if na, ok2 := a.(jJavaNilChecker); ok2 && na.IsJavaNil_() {
			return isJavaNilAny(b)
		}
		return a == b
	}
	if _, ok := a.(EnumAccessor); ok {
		if na, ok2 := a.(jJavaNilChecker); ok2 && na.IsJavaNil_() {
			return isJavaNilAny(b)
		}
		return a == b
	}
	if rna, ok := a.(jRawNodeProvider); ok {
		ua := rna.RawNodeId_()
		if ua == 0 {
			return isJavaNilAny(b)
		}
		if ndb, ok2 := b.(nodestore.Node); ok2 {
			return ua == uint32(ndb)
		}
		return false
	}
	if ia, ok := a.(*BoxInt); ok {
		if ia == nil {
			return isJavaNilAny(b)
		}
		if iv, ok2 := b.(int32); ok2 {
			return ia.Val == iv
		}
		return false
	}
	if fka, ok := a.(FastCustomKey); ok {
		if fka.IsJavaNil_() {
			return isJavaNilAny(b)
		}
		if isJavaNilAny(b) {
			return false
		}
		return a == b || fka.Equals(b)
	}
	if nda, ok := a.(nodestore.Node); ok {
		if nda == 0 {
			return isJavaNilAny(b)
		}
		if ndb, ok2 := b.(nodestore.Node); ok2 {
			return nda == ndb
		}
		return false
	}
	return javalang.ObjectsEquals(a, b)
}

func elemJavaEquals(a, b any, identity bool) bool {
	if identity {
		return RefEquals(a, b)
	}
	return jSetElemEquals(a, b)
}

func fastSetKey(k any, identity bool) (any, bool) {
	if k == nil {
		return nil, true
	}
	if ik, ok := k.(jIdentityKey); ok {
		if ik.IsJavaNil_() {
			return nil, true
		}
		return k, true
	}
	if identity {
		if sp, ok := k.(jSourceObjPtr); ok {
			if sp.IsJavaNil_() {
				return nil, true
			}
			return k, true
		}
		if _, ok := k.(EnumAccessor); ok {
			if nc, ok2 := k.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
				return nil, true
			}
			return k, true
		}
		if rn, ok := k.(jRawNodeProvider); ok {
			u := rn.RawNodeId_()
			if u == 0 {
				return nil, true
			}
			return nodestore.Node(u), true
		}
		if nd, ok := k.(nodestore.Node); ok {
			if nd == 0 {
				return nil, true
			}
			return nd, true
		}
		return nil, false
	}
	if _, ok := k.(IdentityComparable); ok {
		if nc, ok2 := k.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
			return nil, true
		}
		return k, true
	}
	if _, ok := k.(EnumAccessor); ok {
		if nc, ok2 := k.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
			return nil, true
		}
		return k, true
	}
	switch x := k.(type) {
	case string:
		if x == NullString {
			return nil, true
		}
		return x, true
	case nodestore.Node:
		if x == 0 {
			return nil, true
		}
		return x, true
	case *BoxInt:
		if x == nil {
			return nil, true
		}
		return x.Val, true
	case int32:
		return x, true
	}
	if rn, ok := k.(jRawNodeProvider); ok {
		u := rn.RawNodeId_()
		if u == 0 {
			return nil, true
		}
		return nodestore.Node(u), true
	}
	if sp, ok := k.(jSourceObjPtr); ok {
		if sp.IsJavaNil_() {
			return nil, true
		}
		if _, hasEq := k.(javalang.Equaler); !hasEq {
			if _, hasHc := k.(jHashCoder); !hasHc {
				return k, true
			}
		}
	}
	return nil, false
}

func fastSetMapKey(k any) (any, bool) {
	return fastSetKey(k, false)
}

// --- Dedicated Compact JSetCore (HashSet, LinkedHashSet, ImmutableSet.Builder, >8 ImmutableSet) ---

type jSetStrState struct {
	strs   []string
	strIdx map[string]int32
}

type jSetStrState4 struct {
	state jSetStrState
	buf   [4]string
}

type jSetGenBuckets struct {
	hashes   []int32
	hashNext []int32
	buckets  []int32
	genLive  int32
}

type jSetAnyState struct {
	items    []any
	ptrIdx   map[any]int32
	gen      *jSetGenBuckets
	reads    uint16
	allIdent bool
}

func isNonNilIdentityElem(k any) bool {
	if k == nil {
		return false
	}
	ik, ok := k.(jIdentityKey)
	return ok && !ik.IsJavaNil_()
}

func isFastIdentElem(v any) bool {
	if v == nil {
		return false
	}
	if ik, ok := v.(jIdentityKey); ok {
		return !ik.IsJavaNil_()
	}
	if _, ok := v.(IdentityComparable); ok {
		if nc, ok2 := v.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
			return false
		}
		return true
	}
	if _, ok := v.(EnumAccessor); ok {
		if nc, ok2 := v.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
			return false
		}
		return true
	}
	if cls, ok := v.(*javalang.Class); ok {
		return cls != nil
	}
	if globalNodeBoxer == nil {
		if nd, ok := v.(nodestore.Node); ok {
			return nd != 0
		}
	}
	if sp, ok := v.(jSourceObjPtr); ok {
		if sp.IsJavaNil_() {
			return false
		}
		if _, hasEq := v.(javalang.Equaler); !hasEq {
			if _, hasHc := v.(jHashCoder); !hasHc {
				return reflect.TypeOf(v).Comparable()
			}
		}
	}
	return false
}

type jSetAnyState4 struct {
	state jSetAnyState
	buf   [4]any
}

type JSetCore struct {
	nodes    []uint32
	ext      any
	live     int32
	holes    int32
	firstIdx int32
	nodeBuf  [1]uint32
	capHint  int16
	mode     uint8
	identity bool
	imm      bool
}

func NewJSetCore(identity bool) *JSetCore {
	return &JSetCore{identity: identity}
}

func (s *JSetCore) Size() int32   { return s.live }
func (s *JSetCore) IsEmpty() bool { return s.live == 0 }

func (s *JSetCore) InitExpectedSize(expected int32) {
	if expected < 0 {
		panic(javalang.NewIllegalArgumentException("Illegal initial capacity: " + fmt.Sprint(expected)))
	}
	if expected > 4 && s.mode == jSetModeEmpty {
		if expected > 1024 {
			s.capHint = 1024
		} else {
			s.capHint = int16(expected)
		}
	}
}

func (s *JSetCore) initCap() int {
	if s.capHint > 4 {
		return min(int(s.capHint), 256)
	}
	return 4
}

func (s *JSetCore) compactNodes() {
	if s.holes == 0 {
		s.firstIdx = 0
		return
	}
	var idxMap map[uint32]int32
	if s.ext != nil {
		idxMap = s.ext.(map[uint32]int32)
	}
	w := 0
	for i := int(s.firstIdx); i < len(s.nodes); i++ {
		u := s.nodes[i]
		if u != jTombstoneNode {
			s.nodes[w] = u
			if idxMap != nil {
				idxMap[u] = int32(w)
			}
			w++
		}
	}
	s.nodes = s.nodes[:w]
	s.holes = 0
	s.firstIdx = 0
}

func (s *JSetCore) compactStrs(st *jSetStrState) {
	if s.holes == 0 {
		s.firstIdx = 0
		return
	}
	w := 0
	for i := int(s.firstIdx); i < len(st.strs); i++ {
		str := st.strs[i]
		if str != jTombstoneStr {
			st.strs[w] = str
			if st.strIdx != nil {
				st.strIdx[str] = int32(w)
			}
			w++
		}
	}
	clear(st.strs[w:])
	st.strs = st.strs[:w]
	s.holes = 0
	s.firstIdx = 0
}

func (s *JSetCore) rebuildSetGenBuckets(st *jSetAnyState) {
	g := st.gen
	if g == nil {
		return
	}
	capB := 16
	for capB < int(g.genLive)*2 {
		capB <<= 1
	}
	if len(g.buckets) < capB {
		g.buckets = make([]int32, capB)
	} else {
		clear(g.buckets)
	}
	if cap(g.hashes) < len(st.items) {
		g.hashes = make([]int32, len(st.items))
		g.hashNext = make([]int32, len(st.items))
	} else {
		g.hashes = g.hashes[:len(st.items)]
		g.hashNext = g.hashNext[:len(st.items)]
		clear(g.hashNext)
	}
	mask := uint32(len(g.buckets) - 1)
	var cnt int32
	for i, it := range st.items {
		if it == jTombstoneObj {
			continue
		}
		if _, fast := fastSetKey(it, s.identity); fast {
			continue
		}
		h := elemJavaHashCode(it, s.identity)
		g.hashes[i] = h
		slot := smearHash(h) & mask
		g.hashNext[i] = g.buckets[slot]
		g.buckets[slot] = int32(i + 1)
		cnt++
	}
	g.genLive = cnt
}

func (s *JSetCore) compactAny(st *jSetAnyState) {
	if s.holes == 0 {
		s.firstIdx = 0
		return
	}
	w := 0
	for i := int(s.firstIdx); i < len(st.items); i++ {
		it := st.items[i]
		if it != jTombstoneObj {
			st.items[w] = it
			if st.ptrIdx != nil {
				if fk, ok := fastSetKey(it, s.identity); ok {
					st.ptrIdx[fk] = int32(w)
				}
			}
			w++
		}
	}
	clear(st.items[w:])
	st.items = st.items[:w]
	s.holes = 0
	s.firstIdx = 0
	if st.gen != nil {
		s.rebuildSetGenBuckets(st)
	}
}

func (s *JSetCore) materializeAnyIdx(st *jSetAnyState) {
	var needPtr, needGen int32
	for _, it := range st.items {
		if it == jTombstoneObj {
			continue
		}
		if _, ok := fastSetKey(it, s.identity); ok {
			needPtr++
		} else {
			needGen++
		}
	}
	if needPtr > 0 {
		if st.ptrIdx == nil {
			st.ptrIdx = make(map[any]int32, int(needPtr)+4)
		}
		for i, it := range st.items {
			if it == jTombstoneObj {
				continue
			}
			if fk, ok := fastSetKey(it, s.identity); ok {
				st.ptrIdx[fk] = int32(i)
			}
		}
	}
	if needGen > 0 {
		if st.gen == nil {
			st.gen = &jSetGenBuckets{genLive: needGen}
		}
		s.rebuildSetGenBuckets(st)
	}
}

func newSetAnyState(capHint int) *jSetAnyState {
	if capHint <= 4 {
		s4 := &jSetAnyState4{}
		s4.state.items = s4.buf[:0]
		s4.state.allIdent = true
		return &s4.state
	}
	return &jSetAnyState{items: make([]any, 0, capHint), allIdent: true}
}

func (s *JSetCore) upgradeToAny() {
	switch s.mode {
	case jSetModeEmpty:
		s.ext = newSetAnyState(s.initCap())
		s.mode = jSetModeAny
		s.firstIdx = 0
	case jSetModeNode:
		st := newSetAnyState(max(int(s.live)+2, s.initCap()))
		st.allIdent = false
		for i := int(s.firstIdx); i < len(s.nodes); i++ {
			if u := s.nodes[i]; u != jTombstoneNode {
				st.items = append(st.items, boxNodeById(u))
			}
		}
		s.nodes = nil
		s.holes = 0
		s.firstIdx = 0
		s.ext = st
		s.mode = jSetModeAny
		if len(st.items) > jSmallLimit {
			s.materializeAnyIdx(st)
		}
	case jSetModeStr:
		old := s.ext.(*jSetStrState)
		st := newSetAnyState(max(int(s.live)+2, s.initCap()))
		st.allIdent = false
		for i := int(s.firstIdx); i < len(old.strs); i++ {
			str := old.strs[i]
			if str == jTombstoneStr {
				continue
			}
			if str == NullString {
				st.items = append(st.items, nil)
			} else {
				st.items = append(st.items, str)
			}
		}
		s.holes = 0
		s.firstIdx = 0
		s.ext = st
		s.mode = jSetModeAny
		if len(st.items) > jSmallLimit {
			s.materializeAnyIdx(st)
		}
	}
}

func (s *JSetCore) findGenIdx(st *jSetAnyState, k any, h int32) int32 {
	g := st.gen
	if g == nil || len(g.buckets) == 0 {
		return -1
	}
	slot := smearHash(h) & uint32(len(g.buckets)-1)
	if !s.identity {
		if _, isId := k.(jIdentityKey); !isId {
			if fk, ok := k.(FastCustomKey); ok {
				for cur := g.buckets[slot]; cur != 0; cur = g.hashNext[cur-1] {
					i := cur - 1
					if g.hashes[i] == h {
						it := st.items[i]
						if it == k || fk.Equals(it) {
							return i
						}
					}
				}
				return -1
			}
		}
	}
	for cur := g.buckets[slot]; cur != 0; cur = g.hashNext[cur-1] {
		i := cur - 1
		if g.hashes[i] == h {
			it := st.items[i]
			if s.identity {
				if RefEquals(it, k) {
					return i
				}
			} else if javalang.ObjectsEquals(it, k) {
				return i
			}
		}
	}
	return -1
}

func (s *JSetCore) unlinkGenIdx(st *jSetAnyState, idx int32, h int32) {
	g := st.gen
	if g == nil || len(g.buckets) == 0 {
		return
	}
	slot := smearHash(h) & uint32(len(g.buckets)-1)
	target := idx + 1
	var prev int32
	for cur := g.buckets[slot]; cur != 0; cur = g.hashNext[cur-1] {
		if cur == target {
			next := g.hashNext[cur-1]
			if prev == 0 {
				g.buckets[slot] = next
			} else {
				g.hashNext[prev-1] = next
			}
			g.genLive--
			return
		}
		prev = cur
	}
}

func (s *JSetCore) containsAllIdent(st *jSetAnyState, k any) bool {
	if k == nil {
		return false
	}
	if st.ptrIdx != nil {
		if isFastIdentElem(k) {
			_, ok := st.ptrIdx[k]
			return ok
		}
		switch k.(type) {
		case string, *BoxInt, int32, int64, FastCustomKey, jRawNodeProvider, nodestore.Node:
			return false
		}
		if !reflect.TypeOf(k).Comparable() {
			return false
		}
		_, ok := st.ptrIdx[k]
		return ok
	}
	if s.live >= 2 && !s.imm {
		st.reads++
		if st.reads >= 32 {
			s.materializeAnyIdx(st)
			if isFastIdentElem(k) {
				_, ok := st.ptrIdx[k]
				return ok
			}
			switch k.(type) {
			case string, *BoxInt, int32, int64, FastCustomKey, jRawNodeProvider, nodestore.Node:
				return false
			}
			if !reflect.TypeOf(k).Comparable() {
				return false
			}
			_, ok := st.ptrIdx[k]
			return ok
		}
	}
	for _, it := range st.items {
		if it == k {
			return true
		}
	}
	return false
}

func (s *JSetCore) containsAnyNormalized(k any) bool {
	st := s.ext.(*jSetAnyState)
	if st.allIdent {
		return s.containsAllIdent(st, k)
	}
	if st.ptrIdx != nil || st.gen != nil {
		if fk, ok := fastSetKey(k, s.identity); ok {
			if st.ptrIdx == nil {
				return false
			}
			_, found := st.ptrIdx[fk]
			return found
		}
		h := elemJavaHashCode(k, s.identity)
		return s.findGenIdx(st, k, h) >= 0
	}
	if k != nil {
		if ik, ok := k.(jIdentityKey); ok {
			if !ik.IsJavaNil_() {
				for _, it := range st.items {
					if it == k {
						return true
					}
				}
				return false
			}
			k = nil
		} else if !s.identity {
			if _, isRN := k.(jRawNodeProvider); !isRN {
				if _, isBI := k.(*BoxInt); !isBI {
					if fk, ok := k.(FastCustomKey); ok && !fk.IsJavaNil_() {
						for _, it := range st.items {
							if it != jTombstoneObj && (it == k || fk.Equals(it)) {
								return true
							}
						}
						return false
					}
				}
			}
		}
	}
	for _, it := range st.items {
		if it == jTombstoneObj {
			continue
		}
		if s.identity {
			if RefEquals(it, k) {
				return true
			}
		} else if jSetElemEquals(it, k) {
			return true
		}
	}
	return false
}

func (s *JSetCore) addAnyNormalized(k any) bool {
	st := s.ext.(*jSetAnyState)
	if st.ptrIdx != nil || st.gen != nil {
		if fk, ok := fastSetKey(k, s.identity); ok {
			if st.ptrIdx == nil {
				st.ptrIdx = make(map[any]int32, 8)
			} else if _, exists := st.ptrIdx[fk]; exists {
				return false
			}
			if int(s.holes)*2 >= len(st.items) && len(st.items) >= jSmallLimit {
				s.compactAny(st)
			}
			if st.allIdent && !isNonNilIdentityElem(k) {
				st.allIdent = false
			}
			idx := int32(len(st.items))
			st.items = append(st.items, k)
			st.ptrIdx[fk] = idx
			if st.gen != nil {
				st.gen.hashes = append(st.gen.hashes, 0)
				st.gen.hashNext = append(st.gen.hashNext, 0)
			}
			s.live++
			return true
		}
		h := elemJavaHashCode(k, s.identity)
		if st.gen != nil && s.findGenIdx(st, k, h) >= 0 {
			return false
		}
		if int(s.holes)*2 >= len(st.items) && len(st.items) >= jSmallLimit {
			s.compactAny(st)
		}
		st.allIdent = false
		idx := int32(len(st.items))
		st.items = append(st.items, k)
		s.live++
		if st.gen == nil {
			st.gen = &jSetGenBuckets{genLive: 1}
			s.rebuildSetGenBuckets(st)
			return true
		}
		st.gen.genLive++
		if int(st.gen.genLive) >= len(st.gen.buckets) {
			s.rebuildSetGenBuckets(st)
		} else {
			slot := smearHash(h) & uint32(len(st.gen.buckets)-1)
			st.gen.hashes = append(st.gen.hashes, h)
			st.gen.hashNext = append(st.gen.hashNext, st.gen.buckets[slot])
			st.gen.buckets[slot] = idx + 1
		}
		return true
	}
	if k != nil {
		if ik, ok := k.(jIdentityKey); ok {
			if !ik.IsJavaNil_() {
				for _, it := range st.items {
					if it == k {
						return false
					}
				}
				if s.holes > 0 && int(s.holes)*2 >= len(st.items) {
					s.compactAny(st)
				}
				st.items = append(st.items, k)
				s.live++
				if len(st.items) > jSmallLimit {
					s.materializeAnyIdx(st)
				}
				return true
			}
			k = nil
		} else if !s.identity {
			if _, isRN := k.(jRawNodeProvider); !isRN {
				if _, isBI := k.(*BoxInt); !isBI {
					if fk, ok := k.(FastCustomKey); ok && !fk.IsJavaNil_() {
						for _, it := range st.items {
							if it != jTombstoneObj && (it == k || fk.Equals(it)) {
								return false
							}
						}
						if s.holes > 0 && int(s.holes)*2 >= len(st.items) {
							s.compactAny(st)
						}
						st.allIdent = false
						st.items = append(st.items, k)
						s.live++
						if len(st.items) > jSmallLimit {
							s.materializeAnyIdx(st)
						}
						return true
					}
				}
			}
		}
	}
	for _, it := range st.items {
		if it == jTombstoneObj {
			continue
		}
		if s.identity {
			if RefEquals(it, k) {
				return false
			}
		} else if jSetElemEquals(it, k) {
			return false
		}
	}
	if s.holes > 0 && int(s.holes)*2 >= len(st.items) {
		s.compactAny(st)
	}
	st.allIdent = false
	st.items = append(st.items, k)
	s.live++
	if len(st.items) > jSmallLimit {
		s.materializeAnyIdx(st)
	}
	return true
}

func (s *JSetCore) removeAnyAt(idx int) bool {
	st := s.ext.(*jSetAnyState)
	it := st.items[idx]
	if it == jTombstoneObj {
		return false
	}
	if st.ptrIdx != nil {
		if fk, ok := fastSetKey(it, s.identity); ok {
			delete(st.ptrIdx, fk)
		} else if st.gen != nil {
			s.unlinkGenIdx(st, int32(idx), st.gen.hashes[idx])
		}
	} else if st.gen != nil {
		s.unlinkGenIdx(st, int32(idx), st.gen.hashes[idx])
	}
	st.items[idx] = jTombstoneObj
	s.live--
	if s.live == 0 {
		clear(st.items)
		st.items = st.items[:0]
		if st.gen != nil {
			st.gen.hashes = st.gen.hashes[:0]
			st.gen.hashNext = st.gen.hashNext[:0]
		}
		st.reads = 0
		st.allIdent = true
		s.holes = 0
		s.firstIdx = 0
	} else if idx == len(st.items)-1 {
		st.items = st.items[:idx]
		if st.gen != nil && idx < len(st.gen.hashes) {
			st.gen.hashes = st.gen.hashes[:idx]
			st.gen.hashNext = st.gen.hashNext[:idx]
		}
	} else {
		s.holes++
		if idx == int(s.firstIdx) {
			s.firstIdx++
		}
	}
	return true
}

func (s *JSetCore) removeAnyNormalized(k any) bool {
	st := s.ext.(*jSetAnyState)
	if st.ptrIdx != nil || st.gen != nil {
		if fk, ok := fastSetKey(k, s.identity); ok {
			if st.ptrIdx == nil {
				return false
			}
			idx, found := st.ptrIdx[fk]
			if !found {
				return false
			}
			return s.removeAnyAt(int(idx))
		}
		h := elemJavaHashCode(k, s.identity)
		idx := s.findGenIdx(st, k, h)
		if idx < 0 {
			return false
		}
		return s.removeAnyAt(int(idx))
	}
	if k != nil {
		if ik, ok := k.(jIdentityKey); ok {
			if !ik.IsJavaNil_() {
				for i, it := range st.items {
					if it == k {
						return s.removeAnyAt(i)
					}
				}
				return false
			}
			k = nil
		} else if !s.identity {
			if _, isRN := k.(jRawNodeProvider); !isRN {
				if _, isBI := k.(*BoxInt); !isBI {
					if fk, ok := k.(FastCustomKey); ok && !fk.IsJavaNil_() {
						for i, it := range st.items {
							if it != jTombstoneObj && (it == k || fk.Equals(it)) {
								return s.removeAnyAt(i)
							}
						}
						return false
					}
				}
			}
		}
	}
	for i, it := range st.items {
		if it == jTombstoneObj {
			continue
		}
		if s.identity {
			if RefEquals(it, k) {
				return s.removeAnyAt(i)
			}
		} else if jSetElemEquals(it, k) {
			return s.removeAnyAt(i)
		}
	}
	return false
}

func (s *JSetCore) ContainsNode(u uint32) bool {
	if s.live == 0 {
		return false
	}
	switch s.mode {
	case jSetModeNode:
		if u == jTombstoneNode {
			return false
		}
		if s.ext != nil {
			_, ok := s.ext.(map[uint32]int32)[u]
			return ok
		}
		for _, x := range s.nodes {
			if x == u {
				return true
			}
		}
		return false
	case jSetModeStr:
		if u == 0 {
			return s.ContainsStr(NullString)
		}
		return false
	case jSetModeAny:
		if u == 0 {
			return s.containsAnyNormalized(nil)
		}
		return s.containsAnyNormalized(boxNodeById(u))
	}
	return false
}

func (s *JSetCore) AddNode(u uint32) bool {
	if s.imm {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if s.identity || u == jTombstoneNode {
		if s.mode != jSetModeAny {
			s.upgradeToAny()
		}
		return s.addAnyNormalized(boxNodeById(u))
	}
	switch s.mode {
	case jSetModeEmpty:
		if s.capHint > 4 {
			s.nodes = make([]uint32, 1, min(int(s.capHint), 256))
			s.nodes[0] = u
		} else {
			s.nodeBuf[0] = u
			s.nodes = s.nodeBuf[:1]
		}
		s.mode = jSetModeNode
		s.live = 1
		return true
	case jSetModeNode:
		if s.ext != nil {
			idxMap := s.ext.(map[uint32]int32)
			if _, ok := idxMap[u]; ok {
				return false
			}
			if int(s.holes)*2 >= len(s.nodes) && len(s.nodes) >= jNodeSmallLimit {
				s.compactNodes()
			}
			idxMap[u] = int32(len(s.nodes))
			s.nodes = append(s.nodes, u)
			s.live++
			return true
		}
		for _, x := range s.nodes {
			if x == u {
				return false
			}
		}
		if s.holes > 0 && int(s.holes)*2 >= len(s.nodes) {
			s.compactNodes()
		}
		if len(s.nodes) >= jNodeSmallLimit {
			idxMap := make(map[uint32]int32, max(len(s.nodes)+4, s.initCap()))
			for i, x := range s.nodes {
				if x != jTombstoneNode {
					idxMap[x] = int32(i)
				}
			}
			idxMap[u] = int32(len(s.nodes))
			s.ext = idxMap
		}
		if len(s.nodes) == 1 && cap(s.nodes) == 1 {
			nb := make([]uint32, 2, max(8, s.initCap()))
			nb[0] = s.nodes[0]
			nb[1] = u
			s.nodes = nb
		} else {
			s.nodes = append(s.nodes, u)
		}
		s.live++
		return true
	case jSetModeStr:
		if u == 0 {
			return s.AddStr(NullString)
		}
		s.upgradeToAny()
		return s.addAnyNormalized(boxNodeById(u))
	default:
		if u == 0 {
			return s.addAnyNormalized(nil)
		}
		return s.addAnyNormalized(boxNodeById(u))
	}
}

func (s *JSetCore) removeNodeAt(idx int) bool {
	u := s.nodes[idx]
	if u == jTombstoneNode {
		return false
	}
	if s.ext != nil {
		delete(s.ext.(map[uint32]int32), u)
	}
	s.nodes[idx] = jTombstoneNode
	s.live--
	if s.live == 0 {
		s.nodes = s.nodes[:0]
		s.holes = 0
		s.firstIdx = 0
	} else if idx == len(s.nodes)-1 {
		s.nodes = s.nodes[:idx]
	} else {
		s.holes++
		if idx == int(s.firstIdx) {
			s.firstIdx++
		}
	}
	return true
}

func (s *JSetCore) RemoveNode(u uint32) bool {
	if s.imm {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if s.live == 0 {
		return false
	}
	switch s.mode {
	case jSetModeNode:
		if u == jTombstoneNode {
			return false
		}
		if s.ext != nil {
			idx, ok := s.ext.(map[uint32]int32)[u]
			if !ok {
				return false
			}
			return s.removeNodeAt(int(idx))
		}
		for i, x := range s.nodes {
			if x == u {
				return s.removeNodeAt(i)
			}
		}
		return false
	case jSetModeStr:
		if u == 0 {
			return s.RemoveStr(NullString)
		}
		return false
	case jSetModeAny:
		if u == 0 {
			return s.removeAnyNormalized(nil)
		}
		return s.removeAnyNormalized(boxNodeById(u))
	}
	return false
}

func (s *JSetCore) ContainsStr(k string) bool {
	if s.live == 0 {
		return false
	}
	switch s.mode {
	case jSetModeStr:
		if k == jTombstoneStr {
			return false
		}
		st := s.ext.(*jSetStrState)
		if st.strIdx != nil {
			_, ok := st.strIdx[k]
			return ok
		}
		for _, x := range st.strs {
			if x == k {
				return true
			}
		}
		return false
	case jSetModeNode:
		if k == NullString {
			return s.ContainsNode(0)
		}
		return false
	case jSetModeAny:
		if k == NullString {
			return s.containsAnyNormalized(nil)
		}
		return s.containsAnyNormalized(k)
	}
	return false
}

func (s *JSetCore) AddStr(k string) bool {
	if s.imm {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if s.identity || k == jTombstoneStr {
		if s.mode != jSetModeAny {
			s.upgradeToAny()
		}
		if k == NullString {
			return s.addAnyNormalized(nil)
		}
		return s.addAnyNormalized(k)
	}
	switch s.mode {
	case jSetModeEmpty:
		var st *jSetStrState
		if s.capHint > 4 {
			st = &jSetStrState{strs: make([]string, 1, min(int(s.capHint), 256))}
			st.strs[0] = k
		} else {
			s4 := &jSetStrState4{}
			s4.buf[0] = k
			st = &s4.state
			st.strs = s4.buf[:1]
		}
		s.ext = st
		s.mode = jSetModeStr
		s.live = 1
		return true
	case jSetModeStr:
		st := s.ext.(*jSetStrState)
		if st.strIdx != nil {
			if _, ok := st.strIdx[k]; ok {
				return false
			}
			if int(s.holes)*2 >= len(st.strs) && len(st.strs) >= jStrSmallLimit {
				s.compactStrs(st)
			}
			st.strIdx[k] = int32(len(st.strs))
			st.strs = append(st.strs, k)
			s.live++
			return true
		}
		for _, x := range st.strs {
			if x == k {
				return false
			}
		}
		if s.holes > 0 && int(s.holes)*2 >= len(st.strs) {
			s.compactStrs(st)
		}
		if len(st.strs) >= jStrSmallLimit {
			st.strIdx = make(map[string]int32, max(len(st.strs)+1, s.initCap()))
			for i, x := range st.strs {
				if x != jTombstoneStr {
					st.strIdx[x] = int32(i)
				}
			}
			st.strIdx[k] = int32(len(st.strs))
		}
		st.strs = append(st.strs, k)
		s.live++
		return true
	case jSetModeNode:
		if k == NullString {
			return s.AddNode(0)
		}
		s.upgradeToAny()
		return s.addAnyNormalized(k)
	default:
		if k == NullString {
			return s.addAnyNormalized(nil)
		}
		return s.addAnyNormalized(k)
	}
}

func (s *JSetCore) removeStrAt(idx int) bool {
	st := s.ext.(*jSetStrState)
	k := st.strs[idx]
	if k == jTombstoneStr {
		return false
	}
	if st.strIdx != nil {
		delete(st.strIdx, k)
	}
	st.strs[idx] = jTombstoneStr
	s.live--
	if s.live == 0 {
		clear(st.strs)
		st.strs = st.strs[:0]
		s.holes = 0
		s.firstIdx = 0
	} else if idx == len(st.strs)-1 {
		st.strs = st.strs[:idx]
	} else {
		s.holes++
		if idx == int(s.firstIdx) {
			s.firstIdx++
		}
	}
	return true
}

func (s *JSetCore) RemoveStr(k string) bool {
	if s.imm {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if s.live == 0 {
		return false
	}
	switch s.mode {
	case jSetModeStr:
		if k == jTombstoneStr {
			return false
		}
		st := s.ext.(*jSetStrState)
		if st.strIdx != nil {
			idx, ok := st.strIdx[k]
			if !ok {
				return false
			}
			return s.removeStrAt(int(idx))
		}
		for i, x := range st.strs {
			if x == k {
				return s.removeStrAt(i)
			}
		}
		return false
	case jSetModeNode:
		if k == NullString {
			return s.RemoveNode(0)
		}
		return false
	case jSetModeAny:
		if k == NullString {
			return s.removeAnyNormalized(nil)
		}
		return s.removeAnyNormalized(k)
	}
	return false
}

func (s *JSetCore) Contains(k any) bool {
	if s.live == 0 {
		return false
	}
	if s.mode == jSetModeAny {
		st := s.ext.(*jSetAnyState)
		if st.allIdent {
			return s.containsAllIdent(st, k)
		}
	}
	if !s.identity {
		if rn, ok := k.(jRawNodeProvider); ok {
			return s.ContainsNode(rn.RawNodeId_())
		}
		if str, ok := k.(string); ok {
			return s.ContainsStr(str)
		}
		if nd, ok := k.(nodestore.Node); ok {
			return s.ContainsNode(uint32(nd))
		}
	}
	k = ToAny(k)
	if k == nil {
		switch s.mode {
		case jSetModeNode:
			return s.ContainsNode(0)
		case jSetModeStr:
			return s.ContainsStr(NullString)
		case jSetModeAny:
			return s.containsAnyNormalized(nil)
		default:
			return false
		}
	}
	if s.mode != jSetModeAny {
		return false
	}
	return s.containsAnyNormalized(k)
}

func (s *JSetCore) Add(k any) bool {
	if s.imm {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if !s.identity {
		if rn, ok := k.(jRawNodeProvider); ok {
			return s.AddNode(rn.RawNodeId_())
		}
		if str, ok := k.(string); ok {
			return s.AddStr(str)
		}
		if nd, ok := k.(nodestore.Node); ok {
			return s.AddNode(uint32(nd))
		}
	}
	k = ToAny(k)
	if k == nil {
		switch s.mode {
		case jSetModeNode:
			return s.AddNode(0)
		case jSetModeStr:
			return s.AddStr(NullString)
		case jSetModeEmpty:
			s.upgradeToAny()
			return s.addAnyNormalized(nil)
		default:
			return s.addAnyNormalized(nil)
		}
	}
	if s.mode != jSetModeAny {
		s.upgradeToAny()
	}
	return s.addAnyNormalized(k)
}

func (s *JSetCore) Remove(k any) bool {
	if s.imm {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if s.live == 0 {
		return false
	}
	if !s.identity {
		if rn, ok := k.(jRawNodeProvider); ok {
			return s.RemoveNode(rn.RawNodeId_())
		}
		if str, ok := k.(string); ok {
			return s.RemoveStr(str)
		}
		if nd, ok := k.(nodestore.Node); ok {
			return s.RemoveNode(uint32(nd))
		}
	}
	k = ToAny(k)
	if k == nil {
		switch s.mode {
		case jSetModeNode:
			return s.RemoveNode(0)
		case jSetModeStr:
			return s.RemoveStr(NullString)
		case jSetModeAny:
			return s.removeAnyNormalized(nil)
		default:
			return false
		}
	}
	if s.mode != jSetModeAny {
		return false
	}
	return s.removeAnyNormalized(k)
}

func (s *JSetCore) Clear() {
	if s.imm {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if s.live == 0 && s.holes == 0 && s.firstIdx == 0 {
		return
	}
	switch s.mode {
	case jSetModeNode:
		s.nodes = s.nodes[:0]
		if s.ext != nil {
			clear(s.ext.(map[uint32]int32))
		}
	case jSetModeStr:
		st := s.ext.(*jSetStrState)
		clear(st.strs)
		st.strs = st.strs[:0]
		if st.strIdx != nil {
			clear(st.strIdx)
		}
	case jSetModeAny:
		st := s.ext.(*jSetAnyState)
		clear(st.items)
		st.items = st.items[:0]
		if st.ptrIdx != nil {
			clear(st.ptrIdx)
		}
		if st.gen != nil {
			st.gen.hashes = st.gen.hashes[:0]
			st.gen.hashNext = st.gen.hashNext[:0]
			clear(st.gen.buckets)
			st.gen.genLive = 0
		}
		st.reads = 0
		st.allIdent = true
	}
	s.live = 0
	s.holes = 0
	s.firstIdx = 0
}

func (s *JSetCore) ForEach(fn func(any) bool) {
	if s == nil || s.live == 0 {
		return
	}
	start := int(s.firstIdx)
	switch s.mode {
	case jSetModeNode:
		for i := start; i < len(s.nodes); i++ {
			if u := s.nodes[i]; u != jTombstoneNode {
				if !fn(boxNodeById(u)) {
					return
				}
			}
		}
	case jSetModeStr:
		st := s.ext.(*jSetStrState)
		for i := start; i < len(st.strs); i++ {
			str := st.strs[i]
			if str == jTombstoneStr {
				continue
			}
			if !fn(BoxString(str)) {
				return
			}
		}
	case jSetModeAny:
		st := s.ext.(*jSetAnyState)
		for i := start; i < len(st.items); i++ {
			if it := st.items[i]; it != jTombstoneObj {
				if !fn(it) {
					return
				}
			}
		}
	}
}

func (s *JSetCore) LiveKeys() []any {
	if s == nil || s.live == 0 {
		return []any{}
	}
	out := make([]any, 0, s.live)
	s.ForEach(func(k any) bool {
		out = append(out, k)
		return true
	})
	return out
}

func (s *JSetCore) DenseStrSlice() ([]string, bool) {
	if s == nil {
		return nil, false
	}
	if s.live == 0 {
		return emptyStrSlice_, true
	}
	if s.mode == jSetModeStr && s.holes == 0 {
		return s.ext.(*jSetStrState).strs, true
	}
	return nil, false
}

func (s *JSetCore) SetHash() int32 {
	if s == nil || s.live == 0 {
		return 0
	}
	var h int32
	start := int(s.firstIdx)
	switch s.mode {
	case jSetModeNode:
		for i := start; i < len(s.nodes); i++ {
			if u := s.nodes[i]; u != jTombstoneNode {
				h += int32(u)
			}
		}
	case jSetModeStr:
		st := s.ext.(*jSetStrState)
		for i := start; i < len(st.strs); i++ {
			if str := st.strs[i]; str != jTombstoneStr && str != NullString {
				h += javalang.StringHashCode(str)
			}
		}
	case jSetModeAny:
		st := s.ext.(*jSetAnyState)
		for i := start; i < len(st.items); i++ {
			if it := st.items[i]; it != jTombstoneObj {
				h += elemJavaHashCode(it, s.identity)
			}
		}
	}
	return h
}

func (s *JSetCore) FirstKey() any {
	if s == nil || s.live == 0 {
		panic(javalang.NewNoSuchElementException("empty set"))
	}
	start := int(s.firstIdx)
	switch s.mode {
	case jSetModeNode:
		for i := start; i < len(s.nodes); i++ {
			if u := s.nodes[i]; u != jTombstoneNode {
				if i != start && !s.imm {
					s.firstIdx = int32(i)
				}
				return boxNodeById(u)
			}
		}
	case jSetModeStr:
		st := s.ext.(*jSetStrState)
		for i := start; i < len(st.strs); i++ {
			if str := st.strs[i]; str != jTombstoneStr {
				if i != start && !s.imm {
					s.firstIdx = int32(i)
				}
				if str == NullString {
					return nil
				}
				return str
			}
		}
	case jSetModeAny:
		st := s.ext.(*jSetAnyState)
		for i := start; i < len(st.items); i++ {
			if it := st.items[i]; it != jTombstoneObj {
				if i != start && !s.imm {
					s.firstIdx = int32(i)
				}
				return it
			}
		}
	}
	panic(javalang.NewNoSuchElementException("empty set"))
}

func (s *JSetCore) LastKey() any {
	if s == nil || s.live == 0 {
		panic(javalang.NewNoSuchElementException("empty set"))
	}
	switch s.mode {
	case jSetModeNode:
		for i := len(s.nodes) - 1; i >= 0; i-- {
			if s.nodes[i] != jTombstoneNode {
				return boxNodeById(s.nodes[i])
			}
		}
	case jSetModeStr:
		st := s.ext.(*jSetStrState)
		for i := len(st.strs) - 1; i >= 0; i-- {
			if str := st.strs[i]; str != jTombstoneStr {
				if str == NullString {
					return nil
				}
				return str
			}
		}
	case jSetModeAny:
		st := s.ext.(*jSetAnyState)
		for i := len(st.items) - 1; i >= 0; i-- {
			if it := st.items[i]; it != jTombstoneObj {
				return it
			}
		}
	}
	panic(javalang.NewNoSuchElementException("empty set"))
}

func (s *JSetCore) RemoveFirstKey() any {
	if s == nil || s.live == 0 {
		panic(javalang.NewNoSuchElementException("empty set"))
	}
	if s.imm {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	start := int(s.firstIdx)
	switch s.mode {
	case jSetModeNode:
		for i := start; i < len(s.nodes); i++ {
			if u := s.nodes[i]; u != jTombstoneNode {
				s.firstIdx = int32(i)
				s.removeNodeAt(i)
				return boxNodeById(u)
			}
		}
	case jSetModeStr:
		st := s.ext.(*jSetStrState)
		for i := start; i < len(st.strs); i++ {
			if str := st.strs[i]; str != jTombstoneStr {
				s.firstIdx = int32(i)
				s.removeStrAt(i)
				if str == NullString {
					return nil
				}
				return str
			}
		}
	case jSetModeAny:
		st := s.ext.(*jSetAnyState)
		for i := start; i < len(st.items); i++ {
			if it := st.items[i]; it != jTombstoneObj {
				s.firstIdx = int32(i)
				s.removeAnyAt(i)
				return it
			}
		}
	}
	panic(javalang.NewNoSuchElementException("empty set"))
}

func (s *JSetCore) RemoveLastKey() any {
	if s == nil || s.live == 0 {
		panic(javalang.NewNoSuchElementException("empty set"))
	}
	switch s.mode {
	case jSetModeNode:
		for i := len(s.nodes) - 1; i >= 0; i-- {
			if u := s.nodes[i]; u != jTombstoneNode {
				s.removeNodeAt(i)
				return boxNodeById(u)
			}
		}
	case jSetModeStr:
		st := s.ext.(*jSetStrState)
		for i := len(st.strs) - 1; i >= 0; i-- {
			if str := st.strs[i]; str != jTombstoneStr {
				s.removeStrAt(i)
				if str == NullString {
					return nil
				}
				return str
			}
		}
	case jSetModeAny:
		st := s.ext.(*jSetAnyState)
		for i := len(st.items) - 1; i >= 0; i-- {
			if it := st.items[i]; it != jTombstoneObj {
				s.removeAnyAt(i)
				return it
			}
		}
	}
	panic(javalang.NewNoSuchElementException("empty set"))
}

func (s *JSetCore) AddLast(k any) {
	if s.imm {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	s.Remove(k)
	s.Add(k)
}

func (s *JSetCore) AddFirst(k any) {
	if s.imm {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if s.live == 0 {
		s.Add(k)
		return
	}
	s.Remove(k)
	rest := s.LiveKeys()
	s.Clear()
	s.Add(k)
	for _, x := range rest {
		s.Add(x)
	}
}

func (s *JSetCore) sealForImmutable() {
	if s == nil || s.imm {
		return
	}
	switch s.mode {
	case jSetModeNode:
		s.compactNodes()
		if s.ext == nil && len(s.nodes) > jNodeSmallLimit {
			idxMap := make(map[uint32]int32, len(s.nodes))
			for i, u := range s.nodes {
				idxMap[u] = int32(i)
			}
			s.ext = idxMap
		}
	case jSetModeStr:
		st := s.ext.(*jSetStrState)
		s.compactStrs(st)
		if st.strIdx == nil && len(st.strs) > jStrSmallLimit {
			st.strIdx = make(map[string]int32, len(st.strs))
			for i, str := range st.strs {
				st.strIdx[str] = int32(i)
			}
		}
	case jSetModeAny:
		st := s.ext.(*jSetAnyState)
		s.compactAny(st)
		if st.ptrIdx == nil && st.gen == nil && len(st.items) > jSmallLimit {
			s.materializeAnyIdx(st)
		}
	}
	s.firstIdx = 0
	s.imm = true
}

// --- OrderedMap-Style Slim JMapCore (56 bytes, 24B/32B homogeneous entries, 0 inline bloat) ---

const (
	jMapModeEmpty     uint8 = 0
	jMapModeStr       uint8 = 1
	jMapModeNode      uint8 = 2
	jMapModePtr       uint8 = 3
	jMapModeGeneral   uint8 = 4
	jMapModeSortedNat uint8 = 5
	jMapModeSortedCmp uint8 = 6
	jMapModeNullOnly  uint8 = 7
	jMapModeNode1     uint8 = 8

	jKeyGeneral   uint8  = 0
	jKeyNull      uint8  = 1
	jKeyStr       uint8  = 2
	jKeyNode      uint8  = 3
	jKeyInt       uint8  = 4
	jKeyPtr       uint8  = 5
	jKeyFree      uint8  = 7
	jEntryPosMask uint32 = 0x1fffffff
)

type jStrMapEntry struct {
	key string
	val any
}

type jNodeMapEntry struct {
	key    uint32
	active bool
	val    any
}

type jNodeMapState struct {
	entries []jNodeMapEntry
	nodeIdx map[uint32]int32
}

type jNodeMapState4 struct {
	state jNodeMapState
	buf   [4]jNodeMapEntry
}

type jNodeMapState8 struct {
	state jNodeMapState
	buf   [8]jNodeMapEntry
}

func newNodeMapState2(capHint int16) *jNodeMapState {
	if capHint > 8 {
		return &jNodeMapState{entries: make([]jNodeMapEntry, 2, min(int(capHint), 256))}
	}
	if capHint >= 1 && capHint <= 4 {
		s4 := &jNodeMapState4{}
		s4.state.entries = s4.buf[:2]
		return &s4.state
	}
	s8 := &jNodeMapState8{}
	s8.state.entries = s8.buf[:2]
	return &s8.state
}

type jPtrMapEntry struct {
	key any
	val any
}

type jPtrMapState struct {
	entries []jPtrMapEntry
	ptrIdx  map[any]int32
}

type jPtrMapState4 struct {
	state jPtrMapState
	buf   [4]jPtrMapEntry
}

type jTreeNode struct {
	left   int32
	right  int32
	height int8
	isStr  bool
	strKey string
	key    any
	val    any
}

type jSortedMapState struct {
	cmp      func(a, b any) int32
	tree     []jTreeNode
	root     int32
	freeTree int32
	natOrder bool
}

type JMapEntry struct {
	prev        int32
	nextAndKind uint32
	hashNext    int32
	u32OrHash   uint32
	Key         any
	Val         any
}

func makeNextKind(n int32, k uint8) uint32 {
	return (uint32(k) << 29) | (uint32(n) & jEntryPosMask)
}

func (e *JMapEntry) next() int32 {
	return int32(e.nextAndKind & jEntryPosMask)
}

func (e *JMapEntry) setNext(n int32) {
	e.nextAndKind = (e.nextAndKind &^ jEntryPosMask) | (uint32(n) & jEntryPosMask)
}

func (e *JMapEntry) kind() uint8 {
	return uint8(e.nextAndKind >> 29)
}

type jGenMapIdx struct {
	strIdx    map[string]int32
	u64Idx    map[uint64]int32
	ptrIdx    map[any]int32
	buckets   []int32
	otherLive int32
}

type jGenMapState struct {
	entries  []JMapEntry
	idx      *jGenMapIdx
	head     int32
	tail     int32
	freeHead int32
	nullPos  int32
}

type jGenMapState1 struct {
	state  jGenMapState
	inline [1]JMapEntry
}

type jGenMapState2 struct {
	state  jGenMapState
	inline [2]JMapEntry
}

type jGenMapState4 struct {
	state  jGenMapState
	inline [4]JMapEntry
}

func newGenMapState(capHint int16) *jGenMapState {
	if capHint > 4 {
		return &jGenMapState{entries: make([]JMapEntry, 0, min(int(capHint), 256))}
	}
	if capHint == 1 {
		s1 := &jGenMapState1{}
		s1.state.entries = s1.inline[:0]
		return &s1.state
	}
	if capHint == 2 {
		s2 := &jGenMapState2{}
		s2.state.entries = s2.inline[:0]
		return &s2.state
	}
	s4 := &jGenMapState4{}
	s4.state.entries = s4.inline[:0]
	return &s4.state
}

func newGenMapStateSized(n int) *jGenMapState {
	switch {
	case n <= 1:
		s1 := &jGenMapState1{}
		s1.state.entries = s1.inline[:n]
		return &s1.state
	case n == 2:
		s2 := &jGenMapState2{}
		s2.state.entries = s2.inline[:n]
		return &s2.state
	case n <= 4:
		s4 := &jGenMapState4{}
		s4.state.entries = s4.inline[:n]
		return &s4.state
	default:
		return &jGenMapState{entries: make([]JMapEntry, n)}
	}
}

type JMapCore struct {
	strEntries   []jStrMapEntry
	ext          any
	Live         int32
	holes        int32
	firstIdx     int32
	capHint      int16
	mode         uint8
	Identity     bool
	immutable    bool
	accessOrd    bool
	genNonCustom bool
}

type FastComputeMap interface {
	MapComputeIfAbsentStr_(string, Ext_java_util_function_Function) any
	MapComputeIfAbsentNode_(uint32, Ext_java_util_function_Function) any
	MapTryGetComputeStr_(string) (any, bool)
	MapTryGetComputeNode_(uint32) (any, bool)
	MapTryGetComputeAny_(any) (any, bool)
}

type FastStrMap interface {
	MapGetStr_(string) any
	MapGetOrDefaultStr_(string, any) any
	MapContainsStr_(string) bool
	MapPutStr_(string, any) any
	MapRemoveStr_(string) any
}

type FastStrSet interface {
	SetContainsStr_(string) bool
	SetAddStr_(string) bool
	SetRemoveStr_(string) bool
}

type FastNodeMap interface {
	MapGetNode_(uint32) any
	MapGetOrDefaultNode_(uint32, any) any
	MapContainsNode_(uint32) bool
	MapPutNode_(uint32, any) any
	MapRemoveNode_(uint32) any
}

type FastNodeSet interface {
	SetContainsNode_(uint32) bool
	SetAddNode_(uint32) bool
	SetRemoveNode_(uint32) bool
}

type FastMultimap interface {
	MultimapContainsKeyStr_(string) bool
	MultimapContainsKeyNode_(uint32) bool
	MultimapContainsEntryStr_(string, any) bool
	MultimapContainsEntryNode_(uint32, any) bool
	MultimapGetStr_(string) any
	MultimapGetNode_(uint32) any
	MultimapPutStr_(string, any) bool
	MultimapPutNode_(uint32, any) bool
	MultimapRemoveStr_(string, any) bool
	MultimapRemoveNode_(uint32, any) bool
	MultimapRemoveAllStr_(string) any
	MultimapRemoveAllNode_(uint32) any
	MultimapPutAllStr_(string, Ext_java_lang_Iterable) bool
	MultimapPutAllNode_(uint32, Ext_java_lang_Iterable) bool
}

type FastMultiset interface {
	MultisetCountStr_(string) int32
	MultisetContainsStr_(string) bool
	MultisetAddStr1_(string) bool
	MultisetAddStr2_(string, int32) int32
	MultisetRemoveStr1_(string) bool
	MultisetSetCountStr3_(string, int32, int32) bool
}

type NodePredicate interface {
	TestNode_(uint32) bool
}

type jFastDequeCoreProvider interface {
	fastDequeCore_() *JDequeCore
}

type jLinkedDequeCoreProvider interface {
	fastLinkedDequeCore_() *JListCore
}

func NewJMapCore(identity bool) *JMapCore {
	return &JMapCore{Identity: identity, genNonCustom: identity}
}

func NewSortedJMapCore(cmp func(a, b any) int32) *JMapCore {
	m := &JMapCore{}
	if cmp == nil {
		m.mode = jMapModeSortedNat
	} else {
		m.mode = jMapModeSortedCmp
		m.ext = &jSortedMapState{cmp: cmp}
	}
	return m
}

func (m *JMapCore) IsSorted() bool {
	return m != nil && (m.mode == jMapModeSortedNat || m.mode == jMapModeSortedCmp)
}

func (m *JMapCore) SortedCmp() func(a, b any) int32 {
	if m == nil {
		return nil
	}
	if m.mode == jMapModeSortedNat {
		return CompareValues
	}
	if m.mode == jMapModeSortedCmp {
		return m.ext.(*jSortedMapState).cmp
	}
	return nil
}

func (m *JMapCore) IsNaturalOrder() bool {
	return m != nil && m.mode == jMapModeSortedNat
}

func (m *JMapCore) ensureSortedState() *jSortedMapState {
	if m.ext != nil {
		return m.ext.(*jSortedMapState)
	}
	st := &jSortedMapState{cmp: CompareValues, natOrder: true}
	m.ext = st
	return st
}

func (m *JMapCore) InitExpectedSize(expected int32) {
	if expected < 0 {
		panic(javalang.NewIllegalArgumentException("Illegal initial capacity: " + fmt.Sprint(expected)))
	}
	if expected > 0 && (m.mode == jMapModeEmpty || m.mode == jMapModeNullOnly || m.mode == jMapModeNode1) {
		if expected > 4096 {
			expected = 4096
		}
		if int16(expected) > m.capHint {
			m.capHint = int16(expected)
		}
	}
}

func (m *JMapCore) InitAccessOrder(initialCap int32, loadFactor float32, accessOrder bool) {
	if initialCap < 0 {
		panic(javalang.NewIllegalArgumentException("Illegal initial capacity: " + fmt.Sprint(initialCap)))
	}
	if loadFactor <= 0 || math.IsNaN(float64(loadFactor)) {
		panic(javalang.NewIllegalArgumentException("Illegal load factor: " + fmt.Sprint(loadFactor)))
	}
	m.InitExpectedSize(initialCap)
	m.accessOrd = accessOrder
}

func (m *JMapCore) initCap() int {
	if m.capHint > 0 {
		return min(int(m.capHint), 256)
	}
	return 4
}

func (m *JMapCore) strInitCap() int {
	if m.capHint > 0 {
		return min(int(m.capHint), 256)
	}
	return 2
}

func (m *JMapCore) Size() int32   { return m.Live }
func (m *JMapCore) IsEmpty() bool { return m.Live == 0 }

// --- AVL Tree Helpers (m.IsSorted(); 1-based node ids, 0 = nil) ---

func (s *jSortedMapState) treeNodeKey(id int32) any {
	n := &s.tree[id-1]
	if n.key != nil {
		return n.key
	}
	if n.isStr {
		return n.strKey
	}
	return nil
}

func (s *jSortedMapState) avlHeight(id int32) int8 {
	if id == 0 {
		return 0
	}
	return s.tree[id-1].height
}

func (s *jSortedMapState) avlUpdateHeight(id int32) {
	lh := s.avlHeight(s.tree[id-1].left)
	rh := s.avlHeight(s.tree[id-1].right)
	if lh > rh {
		s.tree[id-1].height = lh + 1
	} else {
		s.tree[id-1].height = rh + 1
	}
}

func (s *jSortedMapState) avlBalance(id int32) int8 {
	return s.avlHeight(s.tree[id-1].left) - s.avlHeight(s.tree[id-1].right)
}

func (s *jSortedMapState) avlRotateRight(y int32) int32 {
	x := s.tree[y-1].left
	t2 := s.tree[x-1].right
	s.tree[x-1].right = y
	s.tree[y-1].left = t2
	s.avlUpdateHeight(y)
	s.avlUpdateHeight(x)
	return x
}

func (s *jSortedMapState) avlRotateLeft(x int32) int32 {
	y := s.tree[x-1].right
	t2 := s.tree[y-1].left
	s.tree[y-1].left = x
	s.tree[x-1].right = t2
	s.avlUpdateHeight(x)
	s.avlUpdateHeight(y)
	return y
}

func (s *jSortedMapState) avlRebalance(id int32) int32 {
	s.avlUpdateHeight(id)
	bf := s.avlBalance(id)
	if bf > 1 {
		if s.avlBalance(s.tree[id-1].left) < 0 {
			s.tree[id-1].left = s.avlRotateLeft(s.tree[id-1].left)
		}
		return s.avlRotateRight(id)
	}
	if bf < -1 {
		if s.avlBalance(s.tree[id-1].right) > 0 {
			s.tree[id-1].right = s.avlRotateRight(s.tree[id-1].right)
		}
		return s.avlRotateLeft(id)
	}
	return id
}

func (s *jSortedMapState) avlAllocNode(isStr bool, strKey string, key, val any) int32 {
	if s.freeTree != 0 {
		id := s.freeTree
		s.freeTree = s.tree[id-1].left
		s.tree[id-1] = jTreeNode{height: 1, isStr: isStr, strKey: strKey, key: key, val: val}
		return id
	}
	node := jTreeNode{height: 1, isStr: isStr, strKey: strKey, key: key, val: val}
	if s.tree == nil {
		s.tree = make([]jTreeNode, 1, 4)
		s.tree[0] = node
		return 1
	}
	s.tree = append(s.tree, node)
	return int32(len(s.tree))
}

func (s *jSortedMapState) avlFreeNode(id int32) {
	s.tree[id-1] = jTreeNode{left: s.freeTree}
	s.freeTree = id
}

func (s *jSortedMapState) avlCmpStr(k string, id int32) int32 {
	n := &s.tree[id-1]
	if s.natOrder && n.isStr {
		return int32(javalang.StringCompareTo(k, n.strKey))
	}
	return s.cmp(k, s.treeNodeKey(id))
}

func (s *jSortedMapState) avlCmpAny(k any, id int32) int32 {
	n := &s.tree[id-1]
	if s.natOrder && n.isStr {
		if str, ok := k.(string); ok {
			return int32(javalang.StringCompareTo(str, n.strKey))
		}
	}
	return s.cmp(k, s.treeNodeKey(id))
}

func (s *jSortedMapState) avlFindStr(k string) int32 {
	cur := s.root
	for cur != 0 {
		n := &s.tree[cur-1]
		var c int32
		if s.natOrder && n.isStr {
			c = int32(javalang.StringCompareTo(k, n.strKey))
		} else {
			c = s.cmp(k, s.treeNodeKey(cur))
		}
		if c == 0 {
			return cur
		}
		if c < 0 {
			cur = n.left
		} else {
			cur = n.right
		}
	}
	return 0
}

func (s *jSortedMapState) avlFindAny(k any) int32 {
	cur := s.root
	for cur != 0 {
		c := s.avlCmpAny(k, cur)
		if c == 0 {
			return cur
		}
		if c < 0 {
			cur = s.tree[cur-1].left
		} else {
			cur = s.tree[cur-1].right
		}
	}
	return 0
}

func (m *JMapCore) avlInsertStr(s *jSortedMapState, cur int32, k string, v any) (int32, any, bool) {
	if cur == 0 {
		if s.natOrder {
			_ = javalang.StringCompareTo(k, k)
		} else {
			_ = s.cmp(k, k)
		}
		id := s.avlAllocNode(true, k, nil, v)
		m.Live++
		return id, nil, true
	}
	c := s.avlCmpStr(k, cur)
	if c == 0 {
		old := s.tree[cur-1].val
		s.tree[cur-1].val = v
		return cur, old, false
	}
	if c < 0 {
		nl, old, ins := m.avlInsertStr(s, s.tree[cur-1].left, k, v)
		s.tree[cur-1].left = nl
		if ins {
			return s.avlRebalance(cur), nil, true
		}
		return cur, old, false
	}
	nr, old, ins := m.avlInsertStr(s, s.tree[cur-1].right, k, v)
	s.tree[cur-1].right = nr
	if ins {
		return s.avlRebalance(cur), nil, true
	}
	return cur, old, false
}

func (m *JMapCore) avlInsertAny(s *jSortedMapState, cur int32, k, v any) (int32, any, bool) {
	if cur == 0 {
		_ = s.cmp(k, k)
		if str, ok := k.(string); ok {
			id := s.avlAllocNode(true, str, nil, v)
			m.Live++
			return id, nil, true
		}
		id := s.avlAllocNode(false, "", k, v)
		m.Live++
		return id, nil, true
	}
	c := s.avlCmpAny(k, cur)
	if c == 0 {
		old := s.tree[cur-1].val
		s.tree[cur-1].val = v
		return cur, old, false
	}
	if c < 0 {
		nl, old, ins := m.avlInsertAny(s, s.tree[cur-1].left, k, v)
		s.tree[cur-1].left = nl
		if ins {
			return s.avlRebalance(cur), nil, true
		}
		return cur, old, false
	}
	nr, old, ins := m.avlInsertAny(s, s.tree[cur-1].right, k, v)
	s.tree[cur-1].right = nr
	if ins {
		return s.avlRebalance(cur), nil, true
	}
	return cur, old, false
}

func (s *jSortedMapState) avlExtractMin(cur int32) (int32, int32) {
	if s.tree[cur-1].left == 0 {
		return s.tree[cur-1].right, cur
	}
	nl, minId := s.avlExtractMin(s.tree[cur-1].left)
	s.tree[cur-1].left = nl
	return s.avlRebalance(cur), minId
}

func (s *jSortedMapState) avlExtractMax(cur int32) (int32, int32) {
	if s.tree[cur-1].right == 0 {
		return s.tree[cur-1].left, cur
	}
	nr, maxId := s.avlExtractMax(s.tree[cur-1].right)
	s.tree[cur-1].right = nr
	return s.avlRebalance(cur), maxId
}

func (m *JMapCore) avlUnlinkRoot(s *jSortedMapState, cur int32) int32 {
	l := s.tree[cur-1].left
	r := s.tree[cur-1].right
	s.avlFreeNode(cur)
	m.Live--
	if l == 0 {
		return r
	}
	if r == 0 {
		return l
	}
	nr, succ := s.avlExtractMin(r)
	s.tree[succ-1].left = l
	s.tree[succ-1].right = nr
	return s.avlRebalance(succ)
}

func (m *JMapCore) avlDeleteStr(s *jSortedMapState, cur int32, k string) (int32, any, bool) {
	if cur == 0 {
		return 0, nil, false
	}
	c := s.avlCmpStr(k, cur)
	if c == 0 {
		old := s.tree[cur-1].val
		return m.avlUnlinkRoot(s, cur), old, true
	}
	if c < 0 {
		nl, old, del := m.avlDeleteStr(s, s.tree[cur-1].left, k)
		if !del {
			return cur, nil, false
		}
		s.tree[cur-1].left = nl
		return s.avlRebalance(cur), old, true
	}
	nr, old, del := m.avlDeleteStr(s, s.tree[cur-1].right, k)
	if !del {
		return cur, nil, false
	}
	s.tree[cur-1].right = nr
	return s.avlRebalance(cur), old, true
}

func (m *JMapCore) avlDeleteAny(s *jSortedMapState, cur int32, k any) (int32, any, bool) {
	if cur == 0 {
		return 0, nil, false
	}
	c := s.avlCmpAny(k, cur)
	if c == 0 {
		old := s.tree[cur-1].val
		return m.avlUnlinkRoot(s, cur), old, true
	}
	if c < 0 {
		nl, old, del := m.avlDeleteAny(s, s.tree[cur-1].left, k)
		if !del {
			return cur, nil, false
		}
		s.tree[cur-1].left = nl
		return s.avlRebalance(cur), old, true
	}
	nr, old, del := m.avlDeleteAny(s, s.tree[cur-1].right, k)
	if !del {
		return cur, nil, false
	}
	s.tree[cur-1].right = nr
	return s.avlRebalance(cur), old, true
}

func (s *jSortedMapState) avlForEach(cur int32, fn func(k, v any) bool) bool {
	if cur == 0 {
		return true
	}
	l := s.tree[cur-1].left
	if l != 0 && !s.avlForEach(l, fn) {
		return false
	}
	k := s.treeNodeKey(cur)
	v := s.tree[cur-1].val
	r := s.tree[cur-1].right
	if !fn(k, v) {
		return false
	}
	if r != 0 && !s.avlForEach(r, fn) {
		return false
	}
	return true
}

func (s *jSortedMapState) avlTransformValues(cur int32, fn func(k, v any) any) {
	if cur == 0 {
		return
	}
	s.avlTransformValues(s.tree[cur-1].left, fn)
	s.tree[cur-1].val = ToAny(fn(s.treeNodeKey(cur), s.tree[cur-1].val))
	s.avlTransformValues(s.tree[cur-1].right, fn)
}

func (s *jSortedMapState) avlSetHash(cur int32) int32 {
	if cur == 0 {
		return 0
	}
	n := &s.tree[cur-1]
	var kh int32
	if n.isStr {
		kh = javalang.StringHashCode(n.strKey)
	} else {
		kh = elemJavaHashCode(n.key, false)
	}
	return s.avlSetHash(n.left) + kh + s.avlSetHash(n.right)
}

// --- Homogeneous Mode Compaction & Promotion Helpers ---

func (m *JMapCore) compactStr() {
	if m.holes == 0 {
		m.firstIdx = 0
		return
	}
	var idxMap map[string]int32
	if m.ext != nil {
		idxMap = m.ext.(map[string]int32)
	}
	w := 0
	for i := int(m.firstIdx); i < len(m.strEntries); i++ {
		e := m.strEntries[i]
		if e.key != jTombstoneStr {
			m.strEntries[w] = e
			if idxMap != nil {
				idxMap[e.key] = int32(w)
			}
			w++
		}
	}
	clear(m.strEntries[w:])
	m.strEntries = m.strEntries[:w]
	m.holes = 0
	m.firstIdx = 0
}

func (m *JMapCore) compactNode(ns *jNodeMapState) {
	if m.holes == 0 {
		m.firstIdx = 0
		return
	}
	w := 0
	for i := int(m.firstIdx); i < len(ns.entries); i++ {
		e := ns.entries[i]
		if e.active {
			ns.entries[w] = e
			if ns.nodeIdx != nil {
				ns.nodeIdx[e.key] = int32(w)
			}
			w++
		}
	}
	clear(ns.entries[w:])
	ns.entries = ns.entries[:w]
	m.holes = 0
	m.firstIdx = 0
}

func (m *JMapCore) compactPtr(ps *jPtrMapState) {
	if m.holes == 0 {
		m.firstIdx = 0
		return
	}
	w := 0
	for i := int(m.firstIdx); i < len(ps.entries); i++ {
		e := ps.entries[i]
		if e.key != jTombstoneObj {
			ps.entries[w] = e
			if ps.ptrIdx != nil {
				ps.ptrIdx[e.key] = int32(w)
			}
			w++
		}
	}
	clear(ps.entries[w:])
	ps.entries = ps.entries[:w]
	m.holes = 0
	m.firstIdx = 0
}

func (m *JMapCore) removeStrAt(idx int) any {
	e := &m.strEntries[idx]
	if e.key == jTombstoneStr {
		return nil
	}
	old := e.val
	if m.ext != nil {
		delete(m.ext.(map[string]int32), e.key)
	}
	*e = jStrMapEntry{key: jTombstoneStr}
	m.Live--
	if m.Live == 0 {
		clear(m.strEntries)
		m.strEntries = m.strEntries[:0]
		m.holes = 0
		m.firstIdx = 0
	} else if idx == len(m.strEntries)-1 {
		m.strEntries = m.strEntries[:idx]
	} else {
		m.holes++
		if idx == int(m.firstIdx) {
			m.firstIdx++
		}
	}
	return old
}

func (m *JMapCore) removeNodeAt(idx int) any {
	ns := m.ext.(*jNodeMapState)
	e := &ns.entries[idx]
	if !e.active {
		return nil
	}
	old := e.val
	if ns.nodeIdx != nil {
		delete(ns.nodeIdx, e.key)
	}
	*e = jNodeMapEntry{}
	m.Live--
	if m.Live == 0 {
		clear(ns.entries)
		ns.entries = ns.entries[:0]
		m.holes = 0
		m.firstIdx = 0
	} else if idx == len(ns.entries)-1 {
		ns.entries = ns.entries[:idx]
	} else {
		m.holes++
		if idx == int(m.firstIdx) {
			m.firstIdx++
		}
	}
	return old
}

func (m *JMapCore) removePtrAt(idx int) any {
	ps := m.ext.(*jPtrMapState)
	e := &ps.entries[idx]
	if e.key == jTombstoneObj {
		return nil
	}
	old := e.val
	if ps.ptrIdx != nil {
		delete(ps.ptrIdx, e.key)
	}
	*e = jPtrMapEntry{key: jTombstoneObj}
	m.Live--
	if m.Live == 0 {
		clear(ps.entries)
		ps.entries = ps.entries[:0]
		m.holes = 0
		m.firstIdx = 0
	} else if idx == len(ps.entries)-1 {
		ps.entries = ps.entries[:idx]
	} else {
		m.holes++
		if idx == int(m.firstIdx) {
			m.firstIdx++
		}
	}
	return old
}

// --- General / Mixed / LRU Access-Order State on *JMapCore ---

func (g *jGenMapState) ensureIdx() *jGenMapIdx {
	if g.idx == nil {
		g.idx = &jGenMapIdx{}
	}
	return g.idx
}

func (g *jGenMapState) entryKey(pos int32) any {
	e := &g.entries[pos-1]
	if e.Key != nil {
		return e.Key
	}
	switch e.kind() {
	case jKeyNull:
		return nil
	case jKeyNode:
		return boxNodeById(e.u32OrHash)
	case jKeyInt:
		return ValueOfInt(int32(e.u32OrHash))
	default:
		return e.Key
	}
}

func (m *JMapCore) allocGenEntry(g *jGenMapState, entry JMapEntry) int32 {
	var pos int32
	if g.freeHead != 0 {
		pos = g.freeHead
		g.freeHead = g.entries[pos-1].next()
		g.entries[pos-1] = entry
	} else {
		if len(g.entries) == cap(g.entries) {
			newCap := cap(g.entries) * 2
			if newCap < 4 {
				newCap = 4
			}
			if hint := min(int(m.capHint), 256); hint > newCap {
				newCap = hint
			}
			ne := make([]JMapEntry, len(g.entries)+1, newCap)
			copy(ne, g.entries)
			ne[len(g.entries)] = entry
			g.entries = ne
		} else {
			g.entries = append(g.entries, entry)
		}
		pos = int32(len(g.entries))
	}
	g.entries[pos-1].prev = g.tail
	g.entries[pos-1].setNext(0)
	if g.tail != 0 {
		g.entries[g.tail-1].setNext(pos)
	} else {
		g.head = pos
	}
	g.tail = pos
	m.Live++
	return pos
}

func (g *jGenMapState) moveToTail(pos int32) {
	if pos == 0 || pos == g.tail {
		return
	}
	prev := g.entries[pos-1].prev
	next := g.entries[pos-1].next()
	if prev != 0 {
		g.entries[prev-1].setNext(next)
	} else {
		g.head = next
	}
	g.entries[next-1].prev = prev
	g.entries[pos-1].prev = g.tail
	g.entries[pos-1].setNext(0)
	g.entries[g.tail-1].setNext(pos)
	g.tail = pos
}

func (g *jGenMapState) moveTailToFront() {
	pos := g.tail
	if pos == 0 || pos == g.head {
		return
	}
	prev := g.entries[pos-1].prev
	g.entries[prev-1].setNext(0)
	g.tail = prev
	g.entries[pos-1].prev = 0
	g.entries[pos-1].setNext(g.head)
	g.entries[g.head-1].prev = pos
	g.head = pos
}

func (m *JMapCore) hashGeneral(k any) int32 {
	return elemJavaHashCode(k, m.Identity)
}

func (m *JMapCore) eqGeneral(a, b any) bool {
	if m.Identity {
		return RefEquals(a, b)
	}
	return javalang.ObjectsEquals(a, b)
}

func (m *JMapCore) materializeGenBuckets(g *jGenMapState) *jGenMapIdx {
	idx := g.ensureIdx()
	if len(idx.buckets) == 0 {
		capB := 16
		for capB*3 <= int(m.capHint)*4 && capB < 2048 {
			capB <<= 1
		}
		idx.buckets = make([]int32, capB)
		mask := uint32(capB - 1)
		var cnt int32
		for pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
			e := &g.entries[pos-1]
			if e.kind() == jKeyGeneral {
				slot := smearHash(int32(e.u32OrHash)) & mask
				e.hashNext = idx.buckets[slot]
				idx.buckets[slot] = pos
				cnt++
			}
		}
		idx.otherLive = cnt
		return idx
	}
	if int(idx.otherLive)*4 < len(idx.buckets)*3 {
		return idx
	}
	newCap := len(idx.buckets) * 2
	nb := make([]int32, newCap)
	mask := uint32(newCap - 1)
	for _, headPos := range idx.buckets {
		for pos := headPos; pos != 0; {
			e := &g.entries[pos-1]
			nextHash := e.hashNext
			slot := smearHash(int32(e.u32OrHash)) & mask
			e.hashNext = nb[slot]
			nb[slot] = pos
			pos = nextHash
		}
	}
	idx.buckets = nb
	return idx
}

func findGenFastCustom(g *jGenMapState, k any, fk FastCustomKey, h int32) int32 {
	if idx := g.idx; idx != nil && len(idx.buckets) > 0 {
		entries := g.entries
		slot := smearHash(h) & uint32(len(idx.buckets)-1)
		for pos := idx.buckets[slot]; pos != 0; {
			e := &entries[pos-1]
			if int32(e.u32OrHash) == h && (e.Key == k || fk.Equals(e.Key)) {
				return pos
			}
			pos = e.hashNext
		}
		return 0
	}
	entries := g.entries
	uh := uint32(h)
	for i := range entries {
		e := &entries[i]
		if e.u32OrHash == uh && e.kind() == jKeyGeneral && (e.Key == k || fk.Equals(e.Key)) {
			return int32(i + 1)
		}
	}
	return 0
}

func (m *JMapCore) findGenPos(g *jGenMapState, k any, h int32) int32 {
	if idx := g.idx; idx != nil && len(idx.buckets) > 0 {
		entries := g.entries
		slot := smearHash(h) & uint32(len(idx.buckets)-1)
		for pos := idx.buckets[slot]; pos != 0; {
			e := &entries[pos-1]
			if int32(e.u32OrHash) == h && m.eqGeneral(e.Key, k) {
				return pos
			}
			pos = e.hashNext
		}
		return 0
	}
	entries := g.entries
	uh := uint32(h)
	for i := range entries {
		e := &entries[i]
		if e.u32OrHash == uh && e.kind() == jKeyGeneral && m.eqGeneral(e.Key, k) {
			return int32(i + 1)
		}
	}
	return 0
}

func (g *jGenMapState) unlinkGenBucket(pos int32, h int32) {
	if g.idx == nil || len(g.idx.buckets) == 0 {
		return
	}
	slot := smearHash(h) & uint32(len(g.idx.buckets)-1)
	var prev int32
	for cur := g.idx.buckets[slot]; cur != 0; cur = g.entries[cur-1].hashNext {
		if cur == pos {
			next := g.entries[cur-1].hashNext
			if prev == 0 {
				g.idx.buckets[slot] = next
			} else {
				g.entries[prev-1].hashNext = next
			}
			g.idx.otherLive--
			return
		}
		prev = cur
	}
}

func (m *JMapCore) removeGenPos(pos int32) any {
	g := m.ext.(*jGenMapState)
	e := &g.entries[pos-1]
	old := e.Val
	switch e.kind() {
	case jKeyNull:
		g.nullPos = 0
	case jKeyStr:
		if g.idx != nil && g.idx.strIdx != nil {
			delete(g.idx.strIdx, e.Key.(string))
		}
	case jKeyNode:
		if g.idx != nil && g.idx.u64Idx != nil {
			delete(g.idx.u64Idx, (uint64(1)<<32)|uint64(e.u32OrHash))
		}
	case jKeyInt:
		if g.idx != nil && g.idx.u64Idx != nil {
			delete(g.idx.u64Idx, (uint64(2)<<32)|uint64(e.u32OrHash))
		}
	case jKeyPtr:
		if g.idx != nil && g.idx.ptrIdx != nil {
			delete(g.idx.ptrIdx, e.Key)
		}
	default:
		g.unlinkGenBucket(pos, int32(e.u32OrHash))
	}
	prev, next := e.prev, e.next()
	if prev != 0 {
		g.entries[prev-1].setNext(next)
	} else {
		g.head = next
	}
	if next != 0 {
		g.entries[next-1].prev = prev
	} else {
		g.tail = prev
	}
	g.entries[pos-1] = JMapEntry{nextAndKind: makeNextKind(g.freeHead, jKeyFree)}
	g.freeHead = pos
	m.Live--
	if m.Live == 0 {
		g.head = 0
		g.tail = 0
		g.freeHead = 0
		g.nullPos = 0
		clear(g.entries)
		g.entries = g.entries[:0]
	}
	return old
}

func (g *jGenMapState) materializeStrIdx() *jGenMapIdx {
	idx := g.ensureIdx()
	idx.strIdx = make(map[string]int32, 16)
	for pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
		e := &g.entries[pos-1]
		if e.kind() == jKeyStr {
			idx.strIdx[e.Key.(string)] = pos
		}
	}
	return idx
}

func (g *jGenMapState) materializeU64Idx() *jGenMapIdx {
	idx := g.ensureIdx()
	idx.u64Idx = make(map[uint64]int32, 16)
	for pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
		e := &g.entries[pos-1]
		switch e.kind() {
		case jKeyNode:
			idx.u64Idx[(uint64(1)<<32)|uint64(e.u32OrHash)] = pos
		case jKeyInt:
			idx.u64Idx[(uint64(2)<<32)|uint64(e.u32OrHash)] = pos
		}
	}
	return idx
}

func (g *jGenMapState) materializePtrIdx() *jGenMapIdx {
	idx := g.ensureIdx()
	idx.ptrIdx = make(map[any]int32, 16)
	for pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
		e := &g.entries[pos-1]
		if e.kind() == jKeyPtr {
			idx.ptrIdx[e.Key] = pos
		}
	}
	return idx
}

func (m *JMapCore) putGenNull(g *jGenMapState, v any) (any, bool) {
	if g.nullPos != 0 {
		old := g.entries[g.nullPos-1].Val
		g.entries[g.nullPos-1].Val = v
		if m.accessOrd {
			g.moveToTail(g.nullPos)
		}
		return old, false
	}
	g.nullPos = m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyNull), Val: v})
	return nil, true
}

func (m *JMapCore) putGenStr(g *jGenMapState, k string, boxedKey any, v any) (any, bool) {
	m.genNonCustom = true
	if !m.Identity {
		if boxedKey == nil {
			boxedKey = k
		}
		if g.idx != nil && g.idx.strIdx != nil {
			if pos := g.idx.strIdx[k]; pos != 0 {
				old := g.entries[pos-1].Val
				g.entries[pos-1].Val = v
				if m.accessOrd {
					g.moveToTail(pos)
				}
				return old, false
			}
			pos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyStr), Key: boxedKey, Val: v})
			g.idx.strIdx[k] = pos
			return nil, true
		}
		for pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
			e := &g.entries[pos-1]
			if e.kind() == jKeyStr && e.Key.(string) == k {
				old := e.Val
				e.Val = v
				if m.accessOrd {
					g.moveToTail(pos)
				}
				return old, false
			}
		}
		if int(m.Live) < jSmallLimit {
			m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyStr), Key: boxedKey, Val: v})
			return nil, true
		}
		idx := g.materializeStrIdx()
		pos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyStr), Key: boxedKey, Val: v})
		idx.strIdx[k] = pos
		return nil, true
	}
	if boxedKey != nil {
		return m.putGenOther(g, boxedKey, v)
	}
	return m.putGenOther(g, k, v)
}

func (m *JMapCore) putGenNode(g *jGenMapState, u uint32, v any) (any, bool) {
	m.genNonCustom = true
	if g.idx != nil && g.idx.u64Idx != nil {
		uk := (uint64(1) << 32) | uint64(u)
		if pos := g.idx.u64Idx[uk]; pos != 0 {
			old := g.entries[pos-1].Val
			g.entries[pos-1].Val = v
			if m.accessOrd {
				g.moveToTail(pos)
			}
			return old, false
		}
		pos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyNode), u32OrHash: u, Val: v})
		g.idx.u64Idx[uk] = pos
		return nil, true
	}
	for pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
		e := &g.entries[pos-1]
		if e.kind() == jKeyNode && e.u32OrHash == u {
			old := e.Val
			e.Val = v
			if m.accessOrd {
				g.moveToTail(pos)
			}
			return old, false
		}
	}
	if int(m.Live) < jSmallLimit {
		m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyNode), u32OrHash: u, Val: v})
		return nil, true
	}
	idx := g.materializeU64Idx()
	pos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyNode), u32OrHash: u, Val: v})
	idx.u64Idx[(uint64(1)<<32)|uint64(u)] = pos
	return nil, true
}

func (m *JMapCore) putGenInt32(g *jGenMapState, iv int32, boxedKey any, v any) (any, bool) {
	m.genNonCustom = true
	u := uint32(iv)
	if g.idx != nil && g.idx.u64Idx != nil {
		uk := (uint64(2) << 32) | uint64(u)
		if pos := g.idx.u64Idx[uk]; pos != 0 {
			old := g.entries[pos-1].Val
			g.entries[pos-1].Val = v
			if m.accessOrd {
				g.moveToTail(pos)
			}
			return old, false
		}
		pos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyInt), u32OrHash: u, Key: boxedKey, Val: v})
		g.idx.u64Idx[uk] = pos
		return nil, true
	}
	for pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
		e := &g.entries[pos-1]
		if e.kind() == jKeyInt && e.u32OrHash == u {
			old := e.Val
			e.Val = v
			if m.accessOrd {
				g.moveToTail(pos)
			}
			return old, false
		}
	}
	if int(m.Live) < jSmallLimit {
		m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyInt), u32OrHash: u, Key: boxedKey, Val: v})
		return nil, true
	}
	idx := g.materializeU64Idx()
	pos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyInt), u32OrHash: u, Key: boxedKey, Val: v})
	idx.u64Idx[(uint64(2)<<32)|uint64(u)] = pos
	return nil, true
}

func (m *JMapCore) putGenPtr(g *jGenMapState, k, v any) (any, bool) {
	m.genNonCustom = true
	if g.idx != nil && g.idx.ptrIdx != nil {
		if pos := g.idx.ptrIdx[k]; pos != 0 {
			old := g.entries[pos-1].Val
			g.entries[pos-1].Val = v
			if m.accessOrd {
				g.moveToTail(pos)
			}
			return old, false
		}
		pos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyPtr), Key: k, Val: v})
		g.idx.ptrIdx[k] = pos
		return nil, true
	}
	for pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
		e := &g.entries[pos-1]
		if e.kind() == jKeyPtr && e.Key == k {
			old := e.Val
			e.Val = v
			if m.accessOrd {
				g.moveToTail(pos)
			}
			return old, false
		}
	}
	if int(m.Live) < jSmallLimit {
		m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyPtr), Key: k, Val: v})
		return nil, true
	}
	idx := g.materializePtrIdx()
	pos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyPtr), Key: k, Val: v})
	idx.ptrIdx[k] = pos
	return nil, true
}

func (m *JMapCore) putGenFastCustom(g *jGenMapState, k any, fk FastCustomKey, v any) (any, bool) {
	h := fk.HashCode()
	if pos := findGenFastCustom(g, k, fk, h); pos != 0 {
		old := g.entries[pos-1].Val
		g.entries[pos-1].Val = v
		if m.accessOrd {
			g.moveToTail(pos)
		}
		return old, false
	}
	if int(m.Live) < jSmallLimit && (g.idx == nil || len(g.idx.buckets) == 0) {
		m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyGeneral), u32OrHash: uint32(h), Key: k, Val: v})
		return nil, true
	}
	idx := m.materializeGenBuckets(g)
	slot := smearHash(h) & uint32(len(idx.buckets)-1)
	pos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyGeneral), u32OrHash: uint32(h), hashNext: idx.buckets[slot], Key: k, Val: v})
	idx.buckets[slot] = pos
	idx.otherLive++
	return nil, true
}

func (m *JMapCore) putGenOther(g *jGenMapState, k, v any) (any, bool) {
	h := m.hashGeneral(k)
	if pos := m.findGenPos(g, k, h); pos != 0 {
		old := g.entries[pos-1].Val
		g.entries[pos-1].Val = v
		if m.accessOrd {
			g.moveToTail(pos)
		}
		return old, false
	}
	if int(m.Live) < jSmallLimit && (g.idx == nil || len(g.idx.buckets) == 0) {
		m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyGeneral), u32OrHash: uint32(h), Key: k, Val: v})
		return nil, true
	}
	idx := m.materializeGenBuckets(g)
	slot := smearHash(h) & uint32(len(idx.buckets)-1)
	pos := m.allocGenEntry(g, JMapEntry{nextAndKind: makeNextKind(0, jKeyGeneral), u32OrHash: uint32(h), hashNext: idx.buckets[slot], Key: k, Val: v})
	idx.buckets[slot] = pos
	idx.otherLive++
	return nil, true
}

func newGenMapStateCap(need int) *jGenMapState {
	if need <= 4 {
		s4 := &jGenMapState4{}
		s4.state.entries = s4.inline[:0]
		return &s4.state
	}
	return &jGenMapState{entries: make([]JMapEntry, 0, need)}
}

func (m *JMapCore) upgradeToGeneral() *jGenMapState {
	switch m.mode {
	case jMapModeGeneral:
		return m.ext.(*jGenMapState)
	case jMapModeEmpty:
		g := newGenMapState(m.capHint)
		m.ext = g
		m.mode = jMapModeGeneral
		m.firstIdx = 0
		m.genNonCustom = m.Identity
		return g
	case jMapModeNullOnly:
		nullVal := m.ext
		c := m.capHint
		if c > 0 && c < 2 {
			c = 2
		}
		g := newGenMapState(c)
		m.holes = 0
		m.firstIdx = 0
		m.Live = 0
		m.ext = g
		m.mode = jMapModeGeneral
		m.genNonCustom = m.Identity
		m.putGenNull(g, nullVal)
		return g
	case jMapModeNode1:
		u1 := uint32(m.firstIdx)
		v1 := m.ext
		c := m.capHint
		if c > 0 && c < 2 {
			c = 2
		}
		g := newGenMapState(c)
		m.holes = 0
		m.firstIdx = 0
		m.Live = 0
		m.ext = g
		m.mode = jMapModeGeneral
		m.genNonCustom = m.Identity
		m.putGenNode(g, u1, v1)
		return g
	case jMapModeStr:
		old := m.strEntries
		start := int(m.firstIdx)
		g := newGenMapStateCap(max(int(m.Live)+2, m.initCap()))
		m.strEntries = nil
		m.holes = 0
		m.firstIdx = 0
		m.Live = 0
		m.ext = g
		m.mode = jMapModeGeneral
		m.genNonCustom = m.Identity
		for i := start; i < len(old); i++ {
			e := old[i]
			if e.key == jTombstoneStr {
				continue
			}
			if e.key == NullString {
				m.putGenNull(g, e.val)
			} else {
				m.putGenStr(g, e.key, nil, e.val)
			}
		}
		return g
	case jMapModeNode:
		ns := m.ext.(*jNodeMapState)
		start := int(m.firstIdx)
		g := newGenMapStateCap(max(int(m.Live)+2, m.initCap()))
		m.holes = 0
		m.firstIdx = 0
		m.Live = 0
		m.ext = g
		m.mode = jMapModeGeneral
		m.genNonCustom = m.Identity
		for i := start; i < len(ns.entries); i++ {
			e := ns.entries[i]
			if !e.active {
				continue
			}
			if e.key == 0 {
				m.putGenNull(g, e.val)
			} else {
				m.putGenNode(g, e.key, e.val)
			}
		}
		return g
	case jMapModePtr:
		ps := m.ext.(*jPtrMapState)
		start := int(m.firstIdx)
		g := newGenMapStateCap(max(int(m.Live)+2, m.initCap()))
		m.holes = 0
		m.firstIdx = 0
		m.Live = 0
		m.ext = g
		m.mode = jMapModeGeneral
		m.genNonCustom = m.Identity
		for i := start; i < len(ps.entries); i++ {
			e := ps.entries[i]
			if e.key == jTombstoneObj {
				continue
			}
			if e.key == nil {
				m.putGenNull(g, e.val)
			} else {
				m.putGenPtr(g, e.key, e.val)
			}
		}
		return g
	}
	panic("upgradeToGeneral called on sorted JMapCore")
}

// --- Fast Typed String, Node, Ptr & General Dispatch on *JMapCore ---

func (m *JMapCore) getGenNull(g *jGenMapState, touch bool) (any, bool) {
	if g.nullPos != 0 {
		if touch && m.accessOrd {
			g.moveToTail(g.nullPos)
		}
		return g.entries[g.nullPos-1].Val, true
	}
	return nil, false
}

func (m *JMapCore) getStr(k string, touch bool) (any, bool) {
	if m.Live == 0 {
		return nil, false
	}
	switch m.mode {
	case jMapModeStr:
		if k == jTombstoneStr {
			return nil, false
		}
		if m.ext != nil {
			if idx, ok := m.ext.(map[string]int32)[k]; ok {
				return m.strEntries[idx].val, true
			}
			return nil, false
		}
		for i := range m.strEntries {
			if m.strEntries[i].key == k {
				return m.strEntries[i].val, true
			}
		}
		return nil, false
	case jMapModeNode:
		if k == NullString {
			return m.getNode(0, touch)
		}
		return nil, false
	case jMapModePtr:
		if k == NullString {
			return m.getPtr(nil)
		}
		return nil, false
	case jMapModeNullOnly:
		if k == NullString {
			return m.ext, true
		}
		return nil, false
	case jMapModeNode1:
		return nil, false
	case jMapModeSortedNat, jMapModeSortedCmp:
		if k == NullString {
			return m.getSortedNull()
		}
		st := m.ext.(*jSortedMapState)
		id := st.avlFindStr(k)
		if id == 0 {
			return nil, false
		}
		return st.tree[id-1].val, true
	case jMapModeGeneral:
		g := m.ext.(*jGenMapState)
		if k == NullString {
			return m.getGenNull(g, touch)
		}
		if !m.Identity {
			if g.idx != nil && g.idx.strIdx != nil {
				if pos := g.idx.strIdx[k]; pos != 0 {
					if touch && m.accessOrd {
						g.moveToTail(pos)
					}
					return g.entries[pos-1].Val, true
				}
				return nil, false
			}
			for pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
				e := &g.entries[pos-1]
				if e.kind() == jKeyStr && e.Key.(string) == k {
					if touch && m.accessOrd {
						g.moveToTail(pos)
					}
					return e.Val, true
				}
			}
			return nil, false
		}
		h := m.hashGeneral(k)
		if pos := m.findGenPos(g, k, h); pos != 0 {
			if touch && m.accessOrd {
				g.moveToTail(pos)
			}
			return g.entries[pos-1].Val, true
		}
		return nil, false
	}
	return nil, false
}

func (m *JMapCore) getSortedNull() (any, bool) {
	if m.mode == jMapModeSortedCmp && m.ext != nil && m.Live > 0 {
		st := m.ext.(*jSortedMapState)
		id := st.avlFindAny(nil)
		if id != 0 {
			return st.tree[id-1].val, true
		}
	}
	return nil, false
}

func (m *JMapCore) ContainsStr(k string) bool {
	_, ok := m.getStr(k, false)
	return ok
}

func (m *JMapCore) GetStr(k string) any {
	v, _ := m.getStr(k, true)
	return v
}

func (m *JMapCore) GetOrDefaultStr(k string, def any) any {
	if v, ok := m.getStr(k, true); ok {
		return v
	}
	return ToAny(def)
}

func (m *JMapCore) putStrBoxed(k string, boxedKey any, v any) (any, bool) {
	if m.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	v = ToAny(v)
	switch m.mode {
	case jMapModeEmpty:
		if k == NullString && !m.accessOrd {
			m.mode = jMapModeNullOnly
			m.ext = v
			m.Live = 1
			m.firstIdx = 0
			return nil, true
		}
		if m.Identity || m.accessOrd || k == jTombstoneStr {
			g := m.upgradeToGeneral()
			if k == NullString {
				return m.putGenNull(g, v)
			}
			return m.putGenStr(g, k, boxedKey, v)
		}
		m.strEntries = make([]jStrMapEntry, 1, m.strInitCap())
		m.strEntries[0] = jStrMapEntry{key: k, val: v}
		m.mode = jMapModeStr
		m.Live = 1
		return nil, true
	case jMapModeNullOnly:
		if k == NullString {
			old := m.ext
			m.ext = v
			return old, false
		}
		if m.Identity || m.accessOrd || k == jTombstoneStr {
			g := m.upgradeToGeneral()
			return m.putGenStr(g, k, boxedKey, v)
		}
		nullVal := m.ext
		m.strEntries = make([]jStrMapEntry, 2, max(m.initCap(), 2))
		m.strEntries[0] = jStrMapEntry{key: NullString, val: nullVal}
		m.strEntries[1] = jStrMapEntry{key: k, val: v}
		m.ext = nil
		m.mode = jMapModeStr
		m.Live = 2
		m.holes = 0
		m.firstIdx = 0
		return nil, true
	case jMapModeStr:
		if k == jTombstoneStr {
			g := m.upgradeToGeneral()
			return m.putGenStr(g, k, boxedKey, v)
		}
		if m.ext != nil {
			idxMap := m.ext.(map[string]int32)
			if idx, ok := idxMap[k]; ok {
				old := m.strEntries[idx].val
				m.strEntries[idx].val = v
				return old, false
			}
			if int(m.holes)*2 >= len(m.strEntries) && len(m.strEntries) >= jStrSmallLimit {
				m.compactStr()
			}
			idxMap[k] = int32(len(m.strEntries))
			m.strEntries = append(m.strEntries, jStrMapEntry{key: k, val: v})
			m.Live++
			return nil, true
		}
		for i := range m.strEntries {
			if m.strEntries[i].key == k {
				old := m.strEntries[i].val
				m.strEntries[i].val = v
				return old, false
			}
		}
		if m.holes > 0 && int(m.holes)*2 >= len(m.strEntries) {
			m.compactStr()
		}
		if len(m.strEntries) >= jStrSmallLimit {
			idxMap := make(map[string]int32, max(len(m.strEntries)+1, min(int(m.capHint), 256)))
			for i := range m.strEntries {
				if ek := m.strEntries[i].key; ek != jTombstoneStr {
					idxMap[ek] = int32(i)
				}
			}
			idxMap[k] = int32(len(m.strEntries))
			m.ext = idxMap
		}
		m.strEntries = append(m.strEntries, jStrMapEntry{key: k, val: v})
		m.Live++
		return nil, true
	case jMapModeSortedNat, jMapModeSortedCmp:
		if k == NullString {
			if m.mode == jMapModeSortedNat {
				ThrowNPE()
			}
			st := m.ensureSortedState()
			nr, old, ins := m.avlInsertAny(st, st.root, nil, v)
			st.root = nr
			return old, ins
		}
		st := m.ensureSortedState()
		nr, old, ins := m.avlInsertStr(st, st.root, k, v)
		st.root = nr
		return old, ins
	case jMapModeNode1:
		if k == NullString && !m.Identity && !m.accessOrd {
			return m.putNodeInternal(0, v)
		}
		g := m.upgradeToGeneral()
		if k == NullString {
			return m.putGenNull(g, v)
		}
		return m.putGenStr(g, k, boxedKey, v)
	case jMapModeNode:
		if k == NullString {
			return m.putNodeInternal(0, v)
		}
		g := m.upgradeToGeneral()
		return m.putGenStr(g, k, boxedKey, v)
	case jMapModePtr:
		if k == NullString {
			return m.putPtr(nil, v)
		}
		g := m.upgradeToGeneral()
		return m.putGenStr(g, k, boxedKey, v)
	default:
		g := m.ext.(*jGenMapState)
		if k == NullString {
			return m.putGenNull(g, v)
		}
		return m.putGenStr(g, k, boxedKey, v)
	}
}

func (m *JMapCore) PutStr(k string, v any) any {
	old, _ := m.putStrBoxed(k, nil, v)
	return old
}

func (m *JMapCore) AddKeyStr(k string) bool {
	if _, ok := m.getStr(k, false); ok {
		return false
	}
	_, ins := m.putStrBoxed(k, nil, true)
	return ins
}

func (m *JMapCore) removeStrInternal(k string) (any, bool) {
	if m.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if m.Live == 0 {
		return nil, false
	}
	switch m.mode {
	case jMapModeStr:
		if k == jTombstoneStr {
			return nil, false
		}
		if m.ext != nil {
			idx, ok := m.ext.(map[string]int32)[k]
			if !ok {
				return nil, false
			}
			return m.removeStrAt(int(idx)), true
		}
		for i := range m.strEntries {
			if m.strEntries[i].key == k {
				return m.removeStrAt(i), true
			}
		}
		return nil, false
	case jMapModeNode:
		if k == NullString {
			return m.removeNodeInternal(0)
		}
		return nil, false
	case jMapModePtr:
		if k == NullString {
			return m.removePtr(nil)
		}
		return nil, false
	case jMapModeNullOnly:
		if k == NullString {
			old := m.ext
			m.ext = nil
			m.mode = jMapModeEmpty
			m.Live = 0
			m.holes = 0
			m.firstIdx = 0
			return old, true
		}
		return nil, false
	case jMapModeNode1:
		return nil, false
	case jMapModeSortedNat, jMapModeSortedCmp:
		st := m.ext.(*jSortedMapState)
		if k == NullString {
			if m.mode == jMapModeSortedNat {
				return nil, false
			}
			nr, old, del := m.avlDeleteAny(st, st.root, nil)
			st.root = nr
			return old, del
		}
		nr, old, del := m.avlDeleteStr(st, st.root, k)
		st.root = nr
		return old, del
	case jMapModeGeneral:
		g := m.ext.(*jGenMapState)
		if k == NullString {
			if g.nullPos == 0 {
				return nil, false
			}
			return m.removeGenPos(g.nullPos), true
		}
		if !m.Identity {
			if g.idx != nil && g.idx.strIdx != nil {
				pos := g.idx.strIdx[k]
				if pos == 0 {
					return nil, false
				}
				return m.removeGenPos(pos), true
			}
			for pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
				e := &g.entries[pos-1]
				if e.kind() == jKeyStr && e.Key.(string) == k {
					return m.removeGenPos(pos), true
				}
			}
			return nil, false
		}
		h := m.hashGeneral(k)
		pos := m.findGenPos(g, k, h)
		if pos == 0 {
			return nil, false
		}
		return m.removeGenPos(pos), true
	}
	return nil, false
}

func (m *JMapCore) RemoveStr(k string) any {
	old, _ := m.removeStrInternal(k)
	return old
}

func (m *JMapCore) RemoveKeyStr(k string) bool {
	_, del := m.removeStrInternal(k)
	return del
}

func (m *JMapCore) getNode(u uint32, touch bool) (any, bool) {
	if m.Live == 0 {
		return nil, false
	}
	switch m.mode {
	case jMapModeNode1:
		if u != 0 && uint32(m.firstIdx) == u {
			return m.ext, true
		}
		return nil, false
	case jMapModeNode:
		ns := m.ext.(*jNodeMapState)
		if ns.nodeIdx != nil {
			if idx, ok := ns.nodeIdx[u]; ok {
				return ns.entries[idx].val, true
			}
			return nil, false
		}
		for i := range ns.entries {
			if ns.entries[i].active && ns.entries[i].key == u {
				return ns.entries[i].val, true
			}
		}
		return nil, false
	case jMapModeStr:
		if u == 0 {
			return m.getStr(NullString, false)
		}
		return nil, false
	case jMapModePtr:
		if u == 0 {
			return m.getPtr(nil)
		}
		return nil, false
	case jMapModeNullOnly:
		if u == 0 {
			return m.ext, true
		}
		return nil, false
	case jMapModeSortedNat, jMapModeSortedCmp:
		if u == 0 {
			return m.getSortedNull()
		}
		return m.getAny(boxNodeById(u), false)
	case jMapModeGeneral:
		g := m.ext.(*jGenMapState)
		if u == 0 {
			return m.getGenNull(g, touch)
		}
		if g.idx != nil && g.idx.u64Idx != nil {
			uk := (uint64(1) << 32) | uint64(u)
			if pos := g.idx.u64Idx[uk]; pos != 0 {
				if touch && m.accessOrd {
					g.moveToTail(pos)
				}
				return g.entries[pos-1].Val, true
			}
			return nil, false
		}
		for pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
			e := &g.entries[pos-1]
			if e.kind() == jKeyNode && e.u32OrHash == u {
				if touch && m.accessOrd {
					g.moveToTail(pos)
				}
				return e.Val, true
			}
		}
		return nil, false
	}
	return nil, false
}

func (m *JMapCore) ContainsNode(u uint32) bool {
	_, ok := m.getNode(u, false)
	return ok
}

func (m *JMapCore) GetNode(u uint32) any {
	v, _ := m.getNode(u, true)
	return v
}

func (m *JMapCore) GetOrDefaultNode(u uint32, def any) any {
	if v, ok := m.getNode(u, true); ok {
		return v
	}
	return ToAny(def)
}

func (m *JMapCore) putNodeInternal(u uint32, v any) (any, bool) {
	if m.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	v = ToAny(v)
	switch m.mode {
	case jMapModeEmpty:
		if m.accessOrd {
			g := m.upgradeToGeneral()
			if u == 0 {
				return m.putGenNull(g, v)
			}
			return m.putGenNode(g, u, v)
		}
		if u == 0 {
			m.mode = jMapModeNullOnly
			m.ext = v
			m.Live = 1
			m.firstIdx = 0
			return nil, true
		}
		m.mode = jMapModeNode1
		m.firstIdx = int32(u)
		m.ext = v
		m.Live = 1
		return nil, true
	case jMapModeNode1:
		u1 := uint32(m.firstIdx)
		if u == u1 {
			old := m.ext
			m.ext = v
			return old, false
		}
		if u == 0 && m.Identity {
			g := m.upgradeToGeneral()
			return m.putGenNull(g, v)
		}
		v1 := m.ext
		ns := newNodeMapState2(m.capHint)
		ns.entries[0] = jNodeMapEntry{key: u1, active: true, val: v1}
		ns.entries[1] = jNodeMapEntry{key: u, active: true, val: v}
		m.ext = ns
		m.mode = jMapModeNode
		m.Live = 2
		m.holes = 0
		m.firstIdx = 0
		return nil, true
	case jMapModeNullOnly:
		if u == 0 {
			old := m.ext
			m.ext = v
			return old, false
		}
		if m.Identity || m.accessOrd {
			g := m.upgradeToGeneral()
			return m.putGenNode(g, u, v)
		}
		nullVal := m.ext
		ns := newNodeMapState2(m.capHint)
		ns.entries[0] = jNodeMapEntry{key: 0, active: true, val: nullVal}
		ns.entries[1] = jNodeMapEntry{key: u, active: true, val: v}
		m.ext = ns
		m.mode = jMapModeNode
		m.Live = 2
		m.holes = 0
		m.firstIdx = 0
		return nil, true
	case jMapModeNode:
		ns := m.ext.(*jNodeMapState)
		if ns.nodeIdx != nil {
			if idx, ok := ns.nodeIdx[u]; ok {
				old := ns.entries[idx].val
				ns.entries[idx].val = v
				return old, false
			}
			if int(m.holes)*2 >= len(ns.entries) && len(ns.entries) >= jNodeSmallLimit {
				m.compactNode(ns)
			}
			ns.nodeIdx[u] = int32(len(ns.entries))
			ns.entries = append(ns.entries, jNodeMapEntry{key: u, active: true, val: v})
			m.Live++
			return nil, true
		}
		for i := range ns.entries {
			if ns.entries[i].active && ns.entries[i].key == u {
				old := ns.entries[i].val
				ns.entries[i].val = v
				return old, false
			}
		}
		if m.holes > 0 && int(m.holes)*2 >= len(ns.entries) {
			m.compactNode(ns)
		}
		if len(ns.entries) >= jNodeSmallLimit {
			ns.nodeIdx = make(map[uint32]int32, max(len(ns.entries)+4, min(int(m.capHint), 256)))
			for i := range ns.entries {
				if ns.entries[i].active {
					ns.nodeIdx[ns.entries[i].key] = int32(i)
				}
			}
			ns.nodeIdx[u] = int32(len(ns.entries))
		}
		ns.entries = append(ns.entries, jNodeMapEntry{key: u, active: true, val: v})
		m.Live++
		return nil, true
	case jMapModeStr:
		if u == 0 {
			return m.putStrBoxed(NullString, nil, v)
		}
		g := m.upgradeToGeneral()
		return m.putGenNode(g, u, v)
	case jMapModePtr:
		if u == 0 {
			return m.putPtr(nil, v)
		}
		g := m.upgradeToGeneral()
		return m.putGenNode(g, u, v)
	case jMapModeSortedNat, jMapModeSortedCmp:
		return m.putAnyInternal(boxNodeById(u), v)
	default:
		g := m.ext.(*jGenMapState)
		if u == 0 {
			return m.putGenNull(g, v)
		}
		return m.putGenNode(g, u, v)
	}
}

func (m *JMapCore) PutNode(u uint32, v any) any {
	old, _ := m.putNodeInternal(u, v)
	return old
}

func (m *JMapCore) AddKeyNode(u uint32) bool {
	if _, ok := m.getNode(u, false); ok {
		return false
	}
	_, ins := m.putNodeInternal(u, true)
	return ins
}

func (m *JMapCore) removeNodeInternal(u uint32) (any, bool) {
	if m.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if m.Live == 0 {
		return nil, false
	}
	switch m.mode {
	case jMapModeNode1:
		if u != 0 && uint32(m.firstIdx) == u {
			old := m.ext
			m.ext = nil
			m.mode = jMapModeEmpty
			m.Live = 0
			m.holes = 0
			m.firstIdx = 0
			return old, true
		}
		return nil, false
	case jMapModeNode:
		ns := m.ext.(*jNodeMapState)
		if ns.nodeIdx != nil {
			idx, ok := ns.nodeIdx[u]
			if !ok {
				return nil, false
			}
			return m.removeNodeAt(int(idx)), true
		}
		for i := range ns.entries {
			if ns.entries[i].active && ns.entries[i].key == u {
				return m.removeNodeAt(i), true
			}
		}
		return nil, false
	case jMapModeStr:
		if u == 0 {
			return m.removeStrInternal(NullString)
		}
		return nil, false
	case jMapModePtr:
		if u == 0 {
			return m.removePtr(nil)
		}
		return nil, false
	case jMapModeNullOnly:
		if u == 0 {
			old := m.ext
			m.ext = nil
			m.mode = jMapModeEmpty
			m.Live = 0
			m.holes = 0
			m.firstIdx = 0
			return old, true
		}
		return nil, false
	case jMapModeSortedNat, jMapModeSortedCmp:
		return m.removeAnyInternal(boxNodeById(u))
	case jMapModeGeneral:
		g := m.ext.(*jGenMapState)
		if u == 0 {
			if g.nullPos == 0 {
				return nil, false
			}
			return m.removeGenPos(g.nullPos), true
		}
		if g.idx != nil && g.idx.u64Idx != nil {
			uk := (uint64(1) << 32) | uint64(u)
			pos := g.idx.u64Idx[uk]
			if pos == 0 {
				return nil, false
			}
			return m.removeGenPos(pos), true
		}
		for pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
			e := &g.entries[pos-1]
			if e.kind() == jKeyNode && e.u32OrHash == u {
				return m.removeGenPos(pos), true
			}
		}
		return nil, false
	}
	return nil, false
}

func (m *JMapCore) RemoveNode(u uint32) any {
	old, _ := m.removeNodeInternal(u)
	return old
}

func (m *JMapCore) RemoveKeyNode(u uint32) bool {
	_, del := m.removeNodeInternal(u)
	return del
}

func (m *JMapCore) getPtr(k any) (any, bool) {
	ps := m.ext.(*jPtrMapState)
	if ps.ptrIdx != nil {
		if idx, ok := ps.ptrIdx[k]; ok {
			return ps.entries[idx].val, true
		}
		return nil, false
	}
	for i := range ps.entries {
		if ps.entries[i].key == k {
			return ps.entries[i].val, true
		}
	}
	return nil, false
}

func (m *JMapCore) putPtr(k, v any) (any, bool) {
	v = ToAny(v)
	if m.mode == jMapModeEmpty {
		if m.accessOrd {
			g := m.upgradeToGeneral()
			if k == nil {
				return m.putGenNull(g, v)
			}
			return m.putGenPtr(g, k, v)
		}
		if k == nil {
			m.mode = jMapModeNullOnly
			m.ext = v
			m.Live = 1
			m.firstIdx = 0
			return nil, true
		}
		var ps *jPtrMapState
		if m.capHint > 4 {
			ps = &jPtrMapState{entries: make([]jPtrMapEntry, 1, min(int(m.capHint), 256))}
		} else {
			s4 := &jPtrMapState4{}
			ps = &s4.state
			ps.entries = s4.buf[:1]
		}
		ps.entries[0] = jPtrMapEntry{key: k, val: v}
		m.ext = ps
		m.mode = jMapModePtr
		m.Live = 1
		return nil, true
	}
	if m.mode == jMapModeNullOnly {
		if k == nil {
			old := m.ext
			m.ext = v
			return old, false
		}
		if m.accessOrd {
			g := m.upgradeToGeneral()
			return m.putGenPtr(g, k, v)
		}
		nullVal := m.ext
		var ps *jPtrMapState
		if m.capHint > 4 {
			ps = &jPtrMapState{entries: make([]jPtrMapEntry, 2, min(int(m.capHint), 256))}
		} else {
			s4 := &jPtrMapState4{}
			ps = &s4.state
			ps.entries = s4.buf[:2]
		}
		ps.entries[0] = jPtrMapEntry{key: nil, val: nullVal}
		ps.entries[1] = jPtrMapEntry{key: k, val: v}
		m.ext = ps
		m.mode = jMapModePtr
		m.Live = 2
		m.holes = 0
		m.firstIdx = 0
		return nil, true
	}
	if m.mode != jMapModePtr {
		g := m.upgradeToGeneral()
		if k == nil {
			return m.putGenNull(g, v)
		}
		return m.putGenPtr(g, k, v)
	}
	ps := m.ext.(*jPtrMapState)
	if ps.ptrIdx != nil {
		if idx, ok := ps.ptrIdx[k]; ok {
			old := ps.entries[idx].val
			ps.entries[idx].val = v
			return old, false
		}
		if int(m.holes)*2 >= len(ps.entries) && len(ps.entries) >= jSmallLimit {
			m.compactPtr(ps)
		}
		ps.ptrIdx[k] = int32(len(ps.entries))
		ps.entries = append(ps.entries, jPtrMapEntry{key: k, val: v})
		m.Live++
		return nil, true
	}
	for i := range ps.entries {
		if ps.entries[i].key == k {
			old := ps.entries[i].val
			ps.entries[i].val = v
			return old, false
		}
	}
	if m.holes > 0 && int(m.holes)*2 >= len(ps.entries) {
		m.compactPtr(ps)
	}
	if len(ps.entries) >= jSmallLimit {
		ps.ptrIdx = make(map[any]int32, max(len(ps.entries)+4, min(int(m.capHint), 256)))
		for i := range ps.entries {
			if ek := ps.entries[i].key; ek != jTombstoneObj {
				ps.ptrIdx[ek] = int32(i)
			}
		}
		ps.ptrIdx[k] = int32(len(ps.entries))
	}
	ps.entries = append(ps.entries, jPtrMapEntry{key: k, val: v})
	m.Live++
	return nil, true
}

func (m *JMapCore) removePtr(k any) (any, bool) {
	ps := m.ext.(*jPtrMapState)
	if ps.ptrIdx != nil {
		idx, ok := ps.ptrIdx[k]
		if !ok {
			return nil, false
		}
		return m.removePtrAt(int(idx)), true
	}
	for i := range ps.entries {
		if ps.entries[i].key == k {
			return m.removePtrAt(i), true
		}
	}
	return nil, false
}

func isIdentityOrEnumKey(k any) bool {
	return isFastPtrMapKey(k, false)
}

func isFastPtrMapKey(k any, identity bool) bool {
	if _, ok := k.(IdentityComparable); ok {
		return true
	}
	if _, ok := k.(EnumAccessor); ok {
		return true
	}
	if _, ok := k.(jSourceObjPtr); ok {
		if identity {
			return true
		}
		if _, hasEq := k.(javalang.Equaler); hasEq {
			return false
		}
		_, hasHc := k.(jHashCoder)
		return !hasHc
	}
	return false
}

func (m *JMapCore) getAny(k any, touch bool) (any, bool) {
	if m.Live == 0 {
		return nil, false
	}
	switch m.mode {
	case jMapModeEmpty:
		return nil, false
	case jMapModeStr:
		if s, ok := k.(string); ok {
			return m.getStr(s, false)
		}
		if k == nil || isJavaNilAny(k) {
			return m.getStr(NullString, false)
		}
		return nil, false
	case jMapModeNode1:
		if rn, ok := k.(jRawNodeProvider); ok {
			u := rn.RawNodeId_()
			if u != 0 && uint32(m.firstIdx) == u {
				return m.ext, true
			}
			return nil, false
		}
		if nd, ok := k.(nodestore.Node); ok {
			if nd != 0 && uint32(m.firstIdx) == uint32(nd) {
				return m.ext, true
			}
			return nil, false
		}
		return nil, false
	case jMapModeNode:
		if rn, ok := k.(jRawNodeProvider); ok {
			return m.getNode(rn.RawNodeId_(), false)
		}
		if nd, ok := k.(nodestore.Node); ok {
			return m.getNode(uint32(nd), false)
		}
		if k == nil || isJavaNilAny(k) {
			return m.getNode(0, false)
		}
		return nil, false
	case jMapModePtr:
		if k == nil {
			return m.getPtr(nil)
		}
		if ik, ok := k.(jIdentityKey); ok {
			if ik.IsJavaNil_() {
				return m.getPtr(nil)
			}
			return m.getPtr(k)
		}
		k = ToAny(k)
		if k == nil {
			return m.getPtr(nil)
		}
		if isFastPtrMapKey(k, m.Identity) {
			return m.getPtr(k)
		}
		return nil, false
	case jMapModeGeneral:
		g := m.ext.(*jGenMapState)
		if k == nil {
			return m.getGenNull(g, touch)
		}
		if ik, ok := k.(jIdentityKey); ok {
			if ik.IsJavaNil_() {
				return m.getGenNull(g, touch)
			}
			return m.getGenPtrKey(g, k, touch)
		}
		if rn, ok := k.(jRawNodeProvider); ok {
			return m.getNode(rn.RawNodeId_(), touch)
		}
		if !m.Identity {
			switch x := k.(type) {
			case string:
				return m.getStr(x, touch)
			case nodestore.Node:
				return m.getNode(uint32(x), touch)
			case *BoxInt:
				if x == nil {
					return m.getGenNull(g, touch)
				}
				return m.getGenInt32(g, x.Val, touch)
			case int32:
				return m.getGenInt32(g, x, touch)
			case FastCustomKey:
				if x.IsJavaNil_() {
					return m.getGenNull(g, touch)
				}
				if pos := findGenFastCustom(g, k, x, x.HashCode()); pos != 0 {
					if touch && m.accessOrd {
						g.moveToTail(pos)
					}
					return g.entries[pos-1].Val, true
				}
				return nil, false
			}
		}
		k = ToAny(k)
		if k == nil {
			return m.getGenNull(g, touch)
		}
		if isFastPtrMapKey(k, m.Identity) {
			return m.getGenPtrKey(g, k, touch)
		}
		h := m.hashGeneral(k)
		if pos := m.findGenPos(g, k, h); pos != 0 {
			if touch && m.accessOrd {
				g.moveToTail(pos)
			}
			return g.entries[pos-1].Val, true
		}
		return nil, false
	case jMapModeNullOnly:
		if k == nil || isJavaNilAny(k) {
			return m.ext, true
		}
		return nil, false
	case jMapModeSortedNat, jMapModeSortedCmp:
		if s, ok := k.(string); ok {
			if s == NullString {
				return m.getSortedNull()
			}
			st := m.ext.(*jSortedMapState)
			id := st.avlFindStr(s)
			if id == 0 {
				return nil, false
			}
			return st.tree[id-1].val, true
		}
		k = ToAny(k)
		if k == nil {
			return m.getSortedNull()
		}
		st := m.ext.(*jSortedMapState)
		id := st.avlFindAny(k)
		if id == 0 {
			return nil, false
		}
		return st.tree[id-1].val, true
	}
	return nil, false
}

func (m *JMapCore) getNullKey(touch bool) (any, bool) {
	switch m.mode {
	case jMapModeNullOnly:
		return m.ext, true
	case jMapModeStr:
		return m.getStr(NullString, false)
	case jMapModeNode:
		return m.getNode(0, false)
	case jMapModePtr:
		return m.getPtr(nil)
	case jMapModeGeneral:
		return m.getGenNull(m.ext.(*jGenMapState), touch)
	case jMapModeSortedNat, jMapModeSortedCmp:
		return m.getSortedNull()
	}
	return nil, false
}

func (m *JMapCore) getGenPtrKey(g *jGenMapState, k any, touch bool) (any, bool) {
	if g.idx != nil && g.idx.ptrIdx != nil {
		if pos := g.idx.ptrIdx[k]; pos != 0 {
			if touch && m.accessOrd {
				g.moveToTail(pos)
			}
			return g.entries[pos-1].Val, true
		}
		return nil, false
	}
	for pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
		e := &g.entries[pos-1]
		if e.kind() == jKeyPtr && e.Key == k {
			if touch && m.accessOrd {
				g.moveToTail(pos)
			}
			return e.Val, true
		}
	}
	return nil, false
}

func (m *JMapCore) getGenInt32(g *jGenMapState, iv int32, touch bool) (any, bool) {
	u := uint32(iv)
	if g.idx != nil && g.idx.u64Idx != nil {
		uk := (uint64(2) << 32) | uint64(u)
		if pos := g.idx.u64Idx[uk]; pos != 0 {
			if touch && m.accessOrd {
				g.moveToTail(pos)
			}
			return g.entries[pos-1].Val, true
		}
		return nil, false
	}
	for pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
		e := &g.entries[pos-1]
		if e.kind() == jKeyInt && e.u32OrHash == u {
			if touch && m.accessOrd {
				g.moveToTail(pos)
			}
			return e.Val, true
		}
	}
	return nil, false
}

func (m *JMapCore) putNullKey(v any) (any, bool) {
	v = ToAny(v)
	switch m.mode {
	case jMapModeEmpty:
		if m.accessOrd {
			g := m.upgradeToGeneral()
			return m.putGenNull(g, v)
		}
		m.mode = jMapModeNullOnly
		m.ext = v
		m.Live = 1
		m.firstIdx = 0
		return nil, true
	case jMapModeNullOnly:
		old := m.ext
		m.ext = v
		return old, false
	case jMapModeStr:
		return m.putStrBoxed(NullString, nil, v)
	case jMapModeNode1, jMapModeNode:
		return m.putNodeInternal(0, v)
	case jMapModePtr:
		return m.putPtr(nil, v)
	case jMapModeGeneral:
		return m.putGenNull(m.ext.(*jGenMapState), v)
	case jMapModeSortedNat:
		ThrowNPE()
		return nil, false
	default:
		st := m.ensureSortedState()
		nr, old, ins := m.avlInsertAny(st, st.root, nil, v)
		st.root = nr
		return old, ins
	}
}

func (m *JMapCore) putAnyInternal(k, v any) (any, bool) {
	if m.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if k == nil {
		return m.putNullKey(v)
	}
	if !m.IsSorted() {
		if ik, ok := k.(jIdentityKey); ok {
			if ik.IsJavaNil_() {
				return m.putNullKey(v)
			}
			return m.putPtr(k, v)
		}
		if rn, ok := k.(jRawNodeProvider); ok {
			return m.putNodeInternal(rn.RawNodeId_(), v)
		}
		if !m.Identity {
			switch x := k.(type) {
			case string:
				return m.putStrBoxed(x, k, v)
			case nodestore.Node:
				return m.putNodeInternal(uint32(x), v)
			case *BoxInt:
				if x == nil {
					return m.putNullKey(v)
				}
				g := m.upgradeToGeneral()
				return m.putGenInt32(g, x.Val, x, ToAny(v))
			case int32:
				g := m.upgradeToGeneral()
				return m.putGenInt32(g, x, nil, ToAny(v))
			case FastCustomKey:
				if x.IsJavaNil_() {
					return m.putNullKey(v)
				}
				var g *jGenMapState
				if m.mode == jMapModeGeneral {
					g = m.ext.(*jGenMapState)
				} else {
					g = m.upgradeToGeneral()
				}
				return m.putGenFastCustom(g, k, x, ToAny(v))
			}
		}
	} else if iv, ok := k.(int32); ok {
		k = ValueOfInt(iv)
	}
	k = ToAny(k)
	if k == nil {
		return m.putNullKey(v)
	}
	if m.IsSorted() {
		v = ToAny(v)
		st := m.ensureSortedState()
		nr, old, ins := m.avlInsertAny(st, st.root, k, v)
		st.root = nr
		return old, ins
	}
	if isFastPtrMapKey(k, m.Identity) {
		return m.putPtr(k, v)
	}
	g := m.upgradeToGeneral()
	return m.putGenOther(g, k, ToAny(v))
}

func (m *JMapCore) removeNullKey() (any, bool) {
	switch m.mode {
	case jMapModeNullOnly:
		old := m.ext
		m.ext = nil
		m.mode = jMapModeEmpty
		m.Live = 0
		m.holes = 0
		m.firstIdx = 0
		return old, true
	case jMapModeStr:
		return m.removeStrInternal(NullString)
	case jMapModeNode:
		return m.removeNodeInternal(0)
	case jMapModePtr:
		return m.removePtr(nil)
	case jMapModeGeneral:
		g := m.ext.(*jGenMapState)
		if g.nullPos == 0 {
			return nil, false
		}
		return m.removeGenPos(g.nullPos), true
	case jMapModeSortedCmp:
		if m.Live == 0 {
			return nil, false
		}
		st := m.ext.(*jSortedMapState)
		nr, old, del := m.avlDeleteAny(st, st.root, nil)
		st.root = nr
		return old, del
	}
	return nil, false
}

func (m *JMapCore) removeAnyInternal(k any) (any, bool) {
	if m.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if m.Live == 0 {
		return nil, false
	}
	switch m.mode {
	case jMapModeEmpty:
		return nil, false
	case jMapModeStr:
		if s, ok := k.(string); ok {
			return m.removeStrInternal(s)
		}
		if k == nil || isJavaNilAny(k) {
			return m.removeStrInternal(NullString)
		}
		return nil, false
	case jMapModeNode1:
		if rn, ok := k.(jRawNodeProvider); ok {
			return m.removeNodeInternal(rn.RawNodeId_())
		}
		if nd, ok := k.(nodestore.Node); ok {
			return m.removeNodeInternal(uint32(nd))
		}
		return nil, false
	case jMapModeNode:
		if rn, ok := k.(jRawNodeProvider); ok {
			return m.removeNodeInternal(rn.RawNodeId_())
		}
		if nd, ok := k.(nodestore.Node); ok {
			return m.removeNodeInternal(uint32(nd))
		}
		if k == nil || isJavaNilAny(k) {
			return m.removeNodeInternal(0)
		}
		return nil, false
	case jMapModePtr:
		if k == nil {
			return m.removePtr(nil)
		}
		if ik, ok := k.(jIdentityKey); ok {
			if ik.IsJavaNil_() {
				return m.removePtr(nil)
			}
			return m.removePtr(k)
		}
		k = ToAny(k)
		if k == nil {
			return m.removePtr(nil)
		}
		if isFastPtrMapKey(k, m.Identity) {
			return m.removePtr(k)
		}
		return nil, false
	case jMapModeGeneral:
		g := m.ext.(*jGenMapState)
		if k == nil {
			if g.nullPos == 0 {
				return nil, false
			}
			return m.removeGenPos(g.nullPos), true
		}
		if ik, ok := k.(jIdentityKey); ok {
			if ik.IsJavaNil_() {
				if g.nullPos == 0 {
					return nil, false
				}
				return m.removeGenPos(g.nullPos), true
			}
			return m.removeGenPtrKey(g, k)
		}
		if rn, ok := k.(jRawNodeProvider); ok {
			return m.removeNodeInternal(rn.RawNodeId_())
		}
		if !m.Identity {
			switch x := k.(type) {
			case string:
				return m.removeStrInternal(x)
			case nodestore.Node:
				return m.removeNodeInternal(uint32(x))
			case *BoxInt:
				if x == nil {
					if g.nullPos == 0 {
						return nil, false
					}
					return m.removeGenPos(g.nullPos), true
				}
				return m.removeGenInt32(g, x.Val)
			case int32:
				return m.removeGenInt32(g, x)
			case FastCustomKey:
				if x.IsJavaNil_() {
					if g.nullPos == 0 {
						return nil, false
					}
					return m.removeGenPos(g.nullPos), true
				}
				pos := findGenFastCustom(g, k, x, x.HashCode())
				if pos == 0 {
					return nil, false
				}
				return m.removeGenPos(pos), true
			}
		}
		k = ToAny(k)
		if k == nil {
			if g.nullPos == 0 {
				return nil, false
			}
			return m.removeGenPos(g.nullPos), true
		}
		if isFastPtrMapKey(k, m.Identity) {
			return m.removeGenPtrKey(g, k)
		}
		h := m.hashGeneral(k)
		pos := m.findGenPos(g, k, h)
		if pos == 0 {
			return nil, false
		}
		return m.removeGenPos(pos), true
	case jMapModeNullOnly:
		if k == nil || isJavaNilAny(k) {
			return m.removeNullKey()
		}
		return nil, false
	case jMapModeSortedNat, jMapModeSortedCmp:
		if s, ok := k.(string); ok {
			return m.removeStrInternal(s)
		}
		if iv, ok := k.(int32); ok {
			k = ValueOfInt(iv)
		} else {
			k = ToAny(k)
		}
		if k == nil {
			return m.removeNullKey()
		}
		st := m.ext.(*jSortedMapState)
		nr, old, del := m.avlDeleteAny(st, st.root, k)
		st.root = nr
		return old, del
	}
	return nil, false
}

func (m *JMapCore) removeGenPtrKey(g *jGenMapState, k any) (any, bool) {
	if g.idx != nil && g.idx.ptrIdx != nil {
		pos := g.idx.ptrIdx[k]
		if pos == 0 {
			return nil, false
		}
		return m.removeGenPos(pos), true
	}
	for pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
		e := &g.entries[pos-1]
		if e.kind() == jKeyPtr && e.Key == k {
			return m.removeGenPos(pos), true
		}
	}
	return nil, false
}

func (m *JMapCore) removeGenInt32(g *jGenMapState, iv int32) (any, bool) {
	u := uint32(iv)
	if g.idx != nil && g.idx.u64Idx != nil {
		uk := (uint64(2) << 32) | uint64(u)
		pos := g.idx.u64Idx[uk]
		if pos == 0 {
			return nil, false
		}
		return m.removeGenPos(pos), true
	}
	for pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
		e := &g.entries[pos-1]
		if e.kind() == jKeyInt && e.u32OrHash == u {
			return m.removeGenPos(pos), true
		}
	}
	return nil, false
}

func (m *JMapCore) GetEntry(k any) (any, bool) {
	return m.getAny(k, false)
}

func (m *JMapCore) CopyFromSameMode(src *JMapCore) bool {
	if m == nil || src == nil || m == src || m.immutable || m.mode != jMapModeEmpty || m.Live != 0 {
		return false
	}
	if m.IsSorted() || src.IsSorted() || m.Identity != src.Identity || (src.accessOrd && !m.accessOrd) {
		return false
	}
	if src.Live == 0 {
		return true
	}
	switch src.mode {
	case jMapModeNullOnly:
		if m.accessOrd {
			return false
		}
		m.mode = jMapModeNullOnly
		m.ext = src.ext
		m.Live = 1
		m.holes = 0
		m.firstIdx = 0
		return true
	case jMapModeNode1:
		if m.accessOrd {
			return false
		}
		m.mode = jMapModeNode1
		m.firstIdx = src.firstIdx
		m.ext = src.ext
		m.Live = 1
		m.holes = 0
		return true
	case jMapModeStr:
		if m.accessOrd {
			return false
		}
		if src.holes == 0 && src.firstIdx == 0 {
			n := len(src.strEntries)
			m.strEntries = make([]jStrMapEntry, n)
			copy(m.strEntries, src.strEntries)
			if n > jStrSmallLimit {
				dstIdx := make(map[string]int32, n)
				for i := range m.strEntries {
					dstIdx[m.strEntries[i].key] = int32(i)
				}
				m.ext = dstIdx
			}
		} else {
			live := int(src.Live)
			entries := make([]jStrMapEntry, live)
			var dstIdx map[string]int32
			if live > jStrSmallLimit {
				dstIdx = make(map[string]int32, live)
			}
			w := 0
			for i := int(src.firstIdx); i < len(src.strEntries); i++ {
				e := src.strEntries[i]
				if e.key != jTombstoneStr {
					if dstIdx != nil {
						dstIdx[e.key] = int32(w)
					}
					entries[w] = e
					w++
				}
			}
			m.strEntries = entries
			if dstIdx != nil {
				m.ext = dstIdx
			}
		}
		m.mode = jMapModeStr
		m.Live = src.Live
		m.holes = 0
		m.firstIdx = 0
		return true
	case jMapModeNode:
		if m.accessOrd {
			return false
		}
		sns := src.ext.(*jNodeMapState)
		n := int(src.Live)
		if n == 1 {
			for i := int(src.firstIdx); i < len(sns.entries); i++ {
				if e := sns.entries[i]; e.active {
					if e.key == 0 {
						m.mode = jMapModeNullOnly
						m.firstIdx = 0
					} else {
						m.mode = jMapModeNode1
						m.firstIdx = int32(e.key)
					}
					m.ext = e.val
					m.Live = 1
					m.holes = 0
					return true
				}
			}
		}
		var dns *jNodeMapState
		if n <= 4 {
			s4 := &jNodeMapState4{}
			dns = &s4.state
			dns.entries = s4.buf[:n]
		} else if n <= 8 {
			s8 := &jNodeMapState8{}
			dns = &s8.state
			dns.entries = s8.buf[:n]
		} else {
			dns = &jNodeMapState{entries: make([]jNodeMapEntry, n)}
		}
		if src.holes == 0 && src.firstIdx == 0 && len(sns.entries) == n {
			copy(dns.entries, sns.entries)
		} else {
			w := 0
			for i := int(src.firstIdx); i < len(sns.entries); i++ {
				if e := sns.entries[i]; e.active {
					dns.entries[w] = e
					w++
				}
			}
		}
		if n > jNodeSmallLimit {
			dstIdx := make(map[uint32]int32, n+4)
			for i := range dns.entries {
				dstIdx[dns.entries[i].key] = int32(i)
			}
			dns.nodeIdx = dstIdx
		}
		m.ext = dns
		m.mode = jMapModeNode
		m.Live = src.Live
		m.holes = 0
		m.firstIdx = 0
		return true
	case jMapModePtr:
		if m.accessOrd {
			return false
		}
		sps := src.ext.(*jPtrMapState)
		n := int(src.Live)
		var dps *jPtrMapState
		if n <= 4 {
			s4 := &jPtrMapState4{}
			dps = &s4.state
			dps.entries = s4.buf[:n]
		} else {
			dps = &jPtrMapState{entries: make([]jPtrMapEntry, n)}
		}
		if src.holes == 0 && src.firstIdx == 0 && len(sps.entries) == n {
			copy(dps.entries, sps.entries)
		} else {
			w := 0
			for i := int(src.firstIdx); i < len(sps.entries); i++ {
				if e := sps.entries[i]; e.key != jTombstoneObj {
					dps.entries[w] = e
					w++
				}
			}
		}
		if n > jSmallLimit {
			dstIdx := make(map[any]int32, n+4)
			for i := range dps.entries {
				dstIdx[dps.entries[i].key] = int32(i)
			}
			dps.ptrIdx = dstIdx
		}
		m.ext = dps
		m.mode = jMapModePtr
		m.Live = src.Live
		m.holes = 0
		m.firstIdx = 0
		return true
	case jMapModeGeneral:
		sg := src.ext.(*jGenMapState)
		if sg.freeHead == 0 {
			n := len(sg.entries)
			dg := newGenMapStateSized(n)
			entries := dg.entries
			copy(entries, sg.entries)
			dg.head = sg.head
			dg.tail = sg.tail
			dg.nullPos = sg.nullPos
			if sg.idx != nil && n > jSmallLimit {
				sidx := sg.idx
				didx := &jGenMapIdx{otherLive: sidx.otherLive}
				if len(sidx.buckets) > 0 {
					didx.buckets = make([]int32, len(sidx.buckets))
					copy(didx.buckets, sidx.buckets)
				}
				hasMapIdx := false
				if sidx.strIdx != nil {
					didx.strIdx = make(map[string]int32, len(sidx.strIdx))
					hasMapIdx = true
				}
				if sidx.u64Idx != nil {
					didx.u64Idx = make(map[uint64]int32, len(sidx.u64Idx))
					hasMapIdx = true
				}
				if sidx.ptrIdx != nil {
					didx.ptrIdx = make(map[any]int32, len(sidx.ptrIdx))
					hasMapIdx = true
				}
				if hasMapIdx {
					for i := range entries {
						e := &entries[i]
						pos := int32(i + 1)
						switch e.kind() {
						case jKeyStr:
							if didx.strIdx != nil {
								didx.strIdx[e.Key.(string)] = pos
							}
						case jKeyNode:
							if didx.u64Idx != nil {
								didx.u64Idx[(uint64(1)<<32)|uint64(e.u32OrHash)] = pos
							}
						case jKeyInt:
							if didx.u64Idx != nil {
								didx.u64Idx[(uint64(2)<<32)|uint64(e.u32OrHash)] = pos
							}
						case jKeyPtr:
							if didx.ptrIdx != nil {
								didx.ptrIdx[e.Key] = pos
							}
						}
					}
				}
				dg.idx = didx
			} else if sg.idx != nil {
				for i := range entries {
					entries[i].hashNext = 0
				}
			}
			m.ext = dg
			m.mode = jMapModeGeneral
			m.Live = src.Live
			m.holes = 0
			m.firstIdx = 0
			m.genNonCustom = src.genNonCustom
			return true
		}
		live := int(src.Live)
		dg := newGenMapStateSized(live)
		entries := dg.entries
		dg.head = 1
		dg.tail = src.Live
		var didx *jGenMapIdx
		var bucketMask uint32
		if sg.idx != nil && live > jSmallLimit {
			sidx := sg.idx
			didx = &jGenMapIdx{}
			if len(sidx.buckets) > 0 && sidx.otherLive > 0 {
				capB := 16
				for capB*3 <= int(sidx.otherLive)*4 && capB < 2048 {
					capB <<= 1
				}
				didx.buckets = make([]int32, capB)
				bucketMask = uint32(capB - 1)
			}
			if sidx.strIdx != nil {
				didx.strIdx = make(map[string]int32, len(sidx.strIdx))
			}
			if sidx.u64Idx != nil {
				didx.u64Idx = make(map[uint64]int32, len(sidx.u64Idx))
			}
			if sidx.ptrIdx != nil {
				didx.ptrIdx = make(map[any]int32, len(sidx.ptrIdx))
			}
			dg.idx = didx
		}
		i := 0
		for pos := sg.head; pos != 0; pos = sg.entries[pos-1].next() {
			se := sg.entries[pos-1]
			kind := se.kind()
			newPos := int32(i + 1)
			var nextPos int32
			if i+1 < live {
				nextPos = int32(i + 2)
			}
			de := JMapEntry{
				prev:        int32(i),
				nextAndKind: makeNextKind(nextPos, kind),
				u32OrHash:   se.u32OrHash,
				Key:         se.Key,
				Val:         se.Val,
			}
			switch kind {
			case jKeyNull:
				dg.nullPos = newPos
			case jKeyStr:
				if didx != nil && didx.strIdx != nil {
					didx.strIdx[se.Key.(string)] = newPos
				}
			case jKeyNode:
				if didx != nil && didx.u64Idx != nil {
					didx.u64Idx[(uint64(1)<<32)|uint64(se.u32OrHash)] = newPos
				}
			case jKeyInt:
				if didx != nil && didx.u64Idx != nil {
					didx.u64Idx[(uint64(2)<<32)|uint64(se.u32OrHash)] = newPos
				}
			case jKeyPtr:
				if didx != nil && didx.ptrIdx != nil {
					didx.ptrIdx[se.Key] = newPos
				}
			case jKeyGeneral:
				if didx != nil && len(didx.buckets) > 0 {
					slot := smearHash(int32(se.u32OrHash)) & bucketMask
					de.hashNext = didx.buckets[slot]
					didx.buckets[slot] = newPos
					didx.otherLive++
				}
			}
			entries[i] = de
			i++
		}
		m.ext = dg
		m.mode = jMapModeGeneral
		m.Live = src.Live
		m.holes = 0
		m.firstIdx = 0
		m.genNonCustom = src.genNonCustom
		return true
	}
	return false
}

func (m *JMapCore) ContainsKey(k any) bool {
	_, ok := m.getAny(k, false)
	return ok
}

func (m *JMapCore) ContainsValue(v any) bool {
	v = ToAny(v)
	found := false
	m.ForEachEntry(func(_, ev any) bool {
		if javalang.ObjectsEquals(ev, v) {
			found = true
			return false
		}
		return true
	})
	return found
}

func (m *JMapCore) Get(k any) any {
	v, _ := m.getAny(k, true)
	return v
}

func findMapCoreEqVal[K FastEqValKey](m *JMapCore, k K) (any, bool) {
	if m.Live == 0 {
		return nil, false
	}
	h, nonNull, isIdent := k.FastValHash_()
	if !nonNull {
		return m.getNullKey(false)
	}
	switch m.mode {
	case jMapModeGeneral:
		g := m.ext.(*jGenMapState)
		if idx := g.idx; idx != nil && len(idx.buckets) > 0 {
			entries := g.entries
			slot := smearHash(h) & uint32(len(idx.buckets)-1)
			for pos := idx.buckets[slot]; pos != 0; {
				e := &entries[pos-1]
				if int32(e.u32OrHash) == h {
					if ek, ok := e.Key.(K); ok && ek == k {
						if m.accessOrd {
							g.moveToTail(pos)
						}
						return e.Val, true
					}
				}
				pos = e.hashNext
			}
			return nil, false
		}
		entries := g.entries
		uh := uint32(h)
		for i := range entries {
			e := &entries[i]
			if e.u32OrHash == uh && e.kind() == jKeyGeneral {
				if ek, ok := e.Key.(K); ok && ek == k {
					if m.accessOrd {
						g.moveToTail(int32(i + 1))
					}
					return e.Val, true
				}
			}
		}
		return nil, false
	case jMapModePtr:
		if !isIdent {
			return nil, false
		}
		ps := m.ext.(*jPtrMapState)
		for i := range ps.entries {
			if ek, ok := ps.entries[i].key.(K); ok && ek == k {
				return ps.entries[i].val, true
			}
		}
		return nil, false
	case jMapModeEmpty, jMapModeStr, jMapModeNode, jMapModeNode1, jMapModeNullOnly:
		return nil, false
	default:
		return m.getAny(k, false)
	}
}

func (m *JMapCore) GetOrDefault(k, def any) any {
	if v, ok := m.getAny(k, true); ok {
		return v
	}
	return ToAny(def)
}

func (m *JMapCore) Put(k, v any) any {
	old, _ := m.putAnyInternal(k, v)
	return old
}

func (m *JMapCore) AddKey(k any) bool {
	if _, ok := m.getAny(k, false); ok {
		return false
	}
	_, ins := m.putAnyInternal(k, true)
	return ins
}

func (m *JMapCore) PutIfAbsent(k, v any) any {
	if cur, ok := m.getAny(k, true); ok && cur != nil {
		return cur
	}
	return m.Put(k, v)
}

func (m *JMapCore) Remove(k any) any {
	old, _ := m.removeAnyInternal(k)
	return old
}

func (m *JMapCore) RemoveKey(k any) bool {
	_, del := m.removeAnyInternal(k)
	return del
}

func (m *JMapCore) FirstEntry() (any, any) {
	if m == nil || m.Live == 0 {
		panic(javalang.NewNoSuchElementException("empty set/map"))
	}
	start := int(m.firstIdx)
	switch m.mode {
	case jMapModeSortedNat, jMapModeSortedCmp:
		st := m.ext.(*jSortedMapState)
		cur := st.root
		for st.tree[cur-1].left != 0 {
			cur = st.tree[cur-1].left
		}
		return st.treeNodeKey(cur), st.tree[cur-1].val
	case jMapModeStr:
		for i := start; i < len(m.strEntries); i++ {
			if ek := m.strEntries[i].key; ek != jTombstoneStr {
				if i != start && !m.immutable {
					m.firstIdx = int32(i)
				}
				if ek == NullString {
					return nil, m.strEntries[i].val
				}
				return ek, m.strEntries[i].val
			}
		}
	case jMapModeNode:
		ns := m.ext.(*jNodeMapState)
		for i := start; i < len(ns.entries); i++ {
			if ns.entries[i].active {
				if i != start && !m.immutable {
					m.firstIdx = int32(i)
				}
				return boxNodeById(ns.entries[i].key), ns.entries[i].val
			}
		}
	case jMapModePtr:
		ps := m.ext.(*jPtrMapState)
		for i := start; i < len(ps.entries); i++ {
			if ek := ps.entries[i].key; ek != jTombstoneObj {
				if i != start && !m.immutable {
					m.firstIdx = int32(i)
				}
				return ek, ps.entries[i].val
			}
		}
	case jMapModeGeneral:
		g := m.ext.(*jGenMapState)
		return g.entryKey(g.head), g.entries[g.head-1].Val
	case jMapModeNullOnly:
		return nil, m.ext
	case jMapModeNode1:
		return boxNodeById(uint32(m.firstIdx)), m.ext
	}
	panic(javalang.NewNoSuchElementException("empty set/map"))
}

func (m *JMapCore) FirstKey() any {
	k, _ := m.FirstEntry()
	return k
}

func (m *JMapCore) LastKey() any {
	if m == nil || m.Live == 0 {
		panic(javalang.NewNoSuchElementException("empty set/map"))
	}
	switch m.mode {
	case jMapModeSortedNat, jMapModeSortedCmp:
		st := m.ext.(*jSortedMapState)
		cur := st.root
		for st.tree[cur-1].right != 0 {
			cur = st.tree[cur-1].right
		}
		return st.treeNodeKey(cur)
	case jMapModeStr:
		for i := len(m.strEntries) - 1; i >= 0; i-- {
			if ek := m.strEntries[i].key; ek != jTombstoneStr {
				if ek == NullString {
					return nil
				}
				return ek
			}
		}
	case jMapModeNode:
		ns := m.ext.(*jNodeMapState)
		for i := len(ns.entries) - 1; i >= 0; i-- {
			if ns.entries[i].active {
				return boxNodeById(ns.entries[i].key)
			}
		}
	case jMapModePtr:
		ps := m.ext.(*jPtrMapState)
		for i := len(ps.entries) - 1; i >= 0; i-- {
			if ek := ps.entries[i].key; ek != jTombstoneObj {
				return ek
			}
		}
	case jMapModeGeneral:
		g := m.ext.(*jGenMapState)
		return g.entryKey(g.tail)
	case jMapModeNullOnly:
		return nil
	case jMapModeNode1:
		return boxNodeById(uint32(m.firstIdx))
	}
	panic(javalang.NewNoSuchElementException("empty set/map"))
}

func (m *JMapCore) RemoveFirstKey() any {
	if m == nil || m.Live == 0 {
		panic(javalang.NewNoSuchElementException("empty set/map"))
	}
	if m.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if m.IsSorted() {
		st := m.ext.(*jSortedMapState)
		nr, minId := st.avlExtractMin(st.root)
		st.root = nr
		k := st.treeNodeKey(minId)
		st.avlFreeNode(minId)
		m.Live--
		return k
	}
	start := int(m.firstIdx)
	switch m.mode {
	case jMapModeStr:
		for i := start; i < len(m.strEntries); i++ {
			if ek := m.strEntries[i].key; ek != jTombstoneStr {
				m.firstIdx = int32(i)
				m.removeStrAt(i)
				if ek == NullString {
					return nil
				}
				return ek
			}
		}
	case jMapModeNode:
		ns := m.ext.(*jNodeMapState)
		for i := start; i < len(ns.entries); i++ {
			if ns.entries[i].active {
				k := boxNodeById(ns.entries[i].key)
				m.firstIdx = int32(i)
				m.removeNodeAt(i)
				return k
			}
		}
	case jMapModePtr:
		ps := m.ext.(*jPtrMapState)
		for i := start; i < len(ps.entries); i++ {
			if ek := ps.entries[i].key; ek != jTombstoneObj {
				m.firstIdx = int32(i)
				m.removePtrAt(i)
				return ek
			}
		}
	case jMapModeGeneral:
		g := m.ext.(*jGenMapState)
		k := g.entryKey(g.head)
		m.removeGenPos(g.head)
		return k
	case jMapModeNullOnly:
		m.ext = nil
		m.mode = jMapModeEmpty
		m.Live = 0
		m.holes = 0
		m.firstIdx = 0
		return nil
	case jMapModeNode1:
		k := boxNodeById(uint32(m.firstIdx))
		m.ext = nil
		m.mode = jMapModeEmpty
		m.Live = 0
		m.holes = 0
		m.firstIdx = 0
		return k
	}
	panic(javalang.NewNoSuchElementException("empty set/map"))
}

func (m *JMapCore) RemoveLastKey() any {
	if m == nil || m.Live == 0 {
		panic(javalang.NewNoSuchElementException("empty set/map"))
	}
	if m.IsSorted() {
		st := m.ext.(*jSortedMapState)
		nr, maxId := st.avlExtractMax(st.root)
		st.root = nr
		k := st.treeNodeKey(maxId)
		st.avlFreeNode(maxId)
		m.Live--
		return k
	}
	k := m.LastKey()
	m.Remove(k)
	return k
}

func (m *JMapCore) AddFirstKey(k, v any) {
	if m.IsSorted() {
		_, _ = m.removeAnyInternal(k)
		_, _ = m.putAnyInternal(k, v)
		return
	}
	g := m.upgradeToGeneral()
	_, _ = m.removeAnyInternal(k)
	_, _ = m.putAnyInternal(k, v)
	g.moveTailToFront()
}

func (m *JMapCore) AddLastKey(k, v any) {
	_, _ = m.removeAnyInternal(k)
	_, _ = m.putAnyInternal(k, v)
}

func (m *JMapCore) Clear() {
	if m.immutable {
		panic(javalang.NewUnsupportedOperationException("immutable"))
	}
	if m.Live == 0 && m.holes == 0 && m.firstIdx == 0 {
		return
	}
	switch m.mode {
	case jMapModeNullOnly, jMapModeNode1:
		m.ext = nil
		m.mode = jMapModeEmpty
	case jMapModeStr:
		clear(m.strEntries)
		m.strEntries = m.strEntries[:0]
		if m.ext != nil {
			clear(m.ext.(map[string]int32))
		}
	case jMapModeNode:
		ns := m.ext.(*jNodeMapState)
		clear(ns.entries)
		ns.entries = ns.entries[:0]
		if ns.nodeIdx != nil {
			clear(ns.nodeIdx)
		}
	case jMapModePtr:
		ps := m.ext.(*jPtrMapState)
		clear(ps.entries)
		ps.entries = ps.entries[:0]
		if ps.ptrIdx != nil {
			clear(ps.ptrIdx)
		}
	case jMapModeGeneral:
		g := m.ext.(*jGenMapState)
		g.head = 0
		g.tail = 0
		g.freeHead = 0
		g.nullPos = 0
		if g.entries != nil {
			clear(g.entries)
			g.entries = g.entries[:0]
		}
		if g.idx != nil {
			if g.idx.strIdx != nil {
				clear(g.idx.strIdx)
			}
			if g.idx.u64Idx != nil {
				clear(g.idx.u64Idx)
			}
			if g.idx.ptrIdx != nil {
				clear(g.idx.ptrIdx)
			}
			if len(g.idx.buckets) > 0 {
				clear(g.idx.buckets)
			}
			g.idx.otherLive = 0
		}
		m.genNonCustom = m.Identity
	case jMapModeSortedNat, jMapModeSortedCmp:
		if m.ext != nil {
			st := m.ext.(*jSortedMapState)
			if st.tree != nil {
				clear(st.tree)
				st.tree = st.tree[:0]
			}
			st.root = 0
			st.freeTree = 0
		}
	}
	m.Live = 0
	m.holes = 0
	m.firstIdx = 0
}

func (m *JMapCore) sealImmutable() {
	if m == nil || m.immutable {
		return
	}
	switch m.mode {
	case jMapModeStr:
		m.compactStr()
		m.firstIdx = 0
	case jMapModeNode:
		m.compactNode(m.ext.(*jNodeMapState))
		m.firstIdx = 0
	case jMapModePtr:
		m.compactPtr(m.ext.(*jPtrMapState))
		m.firstIdx = 0
	case jMapModeNode1:
	default:
		m.firstIdx = 0
	}
	m.immutable = true
}

func (m *JMapCore) SetHash() int32 {
	if m == nil || m.Live == 0 {
		return 0
	}
	if m.IsSorted() {
		return m.ext.(*jSortedMapState).avlSetHash(m.ext.(*jSortedMapState).root)
	}
	var h int32
	m.ForEachEntry(func(k, _ any) bool {
		h += elemJavaHashCode(k, m.Identity)
		return true
	})
	return h
}

func (m *JMapCore) ImmutableSetHash() int32 {
	return m.SetHash()
}

func (m *JMapCore) ForEachEntry(fn func(k, v any) bool) {
	if m == nil || m.Live == 0 {
		return
	}
	start := int(m.firstIdx)
	switch m.mode {
	case jMapModeNullOnly:
		fn(nil, m.ext)
	case jMapModeNode1:
		fn(boxNodeById(uint32(m.firstIdx)), m.ext)
	case jMapModeStr:
		for i := start; i < len(m.strEntries); i++ {
			e := m.strEntries[i]
			if e.key == jTombstoneStr {
				continue
			}
			if !fn(BoxString(e.key), e.val) {
				return
			}
		}
	case jMapModeNode:
		ns := m.ext.(*jNodeMapState)
		for i := start; i < len(ns.entries); i++ {
			e := ns.entries[i]
			if !e.active {
				continue
			}
			if !fn(boxNodeById(e.key), e.val) {
				return
			}
		}
	case jMapModePtr:
		ps := m.ext.(*jPtrMapState)
		for i := start; i < len(ps.entries); i++ {
			e := ps.entries[i]
			if e.key == jTombstoneObj {
				continue
			}
			if !fn(e.key, e.val) {
				return
			}
		}
	case jMapModeGeneral:
		g := m.ext.(*jGenMapState)
		for pos := g.head; pos != 0; {
			nextPos := g.entries[pos-1].next()
			k := g.entryKey(pos)
			v := g.entries[pos-1].Val
			if !fn(k, v) {
				return
			}
			pos = nextPos
		}
	case jMapModeSortedNat, jMapModeSortedCmp:
		st := m.ext.(*jSortedMapState)
		st.avlForEach(st.root, fn)
	}
}

func (m *JMapCore) TransformValues(fn func(k, v any) any) {
	if m == nil || m.Live == 0 {
		return
	}
	start := int(m.firstIdx)
	switch m.mode {
	case jMapModeNullOnly:
		m.ext = ToAny(fn(nil, m.ext))
	case jMapModeNode1:
		m.ext = ToAny(fn(boxNodeById(uint32(m.firstIdx)), m.ext))
	case jMapModeStr:
		for i := start; i < len(m.strEntries); i++ {
			e := &m.strEntries[i]
			if e.key == jTombstoneStr {
				continue
			}
			var k any
			if e.key != NullString {
				k = e.key
			}
			e.val = ToAny(fn(k, e.val))
		}
	case jMapModeNode:
		ns := m.ext.(*jNodeMapState)
		for i := start; i < len(ns.entries); i++ {
			e := &ns.entries[i]
			if !e.active {
				continue
			}
			e.val = ToAny(fn(boxNodeById(e.key), e.val))
		}
	case jMapModePtr:
		ps := m.ext.(*jPtrMapState)
		for i := start; i < len(ps.entries); i++ {
			e := &ps.entries[i]
			if e.key == jTombstoneObj {
				continue
			}
			e.val = ToAny(fn(e.key, e.val))
		}
	case jMapModeGeneral:
		g := m.ext.(*jGenMapState)
		for pos := g.head; pos != 0; pos = g.entries[pos-1].next() {
			k := g.entryKey(pos)
			g.entries[pos-1].Val = ToAny(fn(k, g.entries[pos-1].Val))
		}
	case jMapModeSortedNat, jMapModeSortedCmp:
		st := m.ext.(*jSortedMapState)
		st.avlTransformValues(st.root, fn)
	}
}

func (m *JMapCore) LiveKeys() []any {
	if m == nil || m.Live == 0 {
		return []any{}
	}
	out := make([]any, 0, m.Live)
	m.ForEachEntry(func(k, _ any) bool {
		out = append(out, k)
		return true
	})
	return out
}

func (m *JMapCore) LiveStrKeys() ([]string, bool) {
	if m == nil || m.Live == 0 {
		return nil, true
	}
	if m.mode != jMapModeStr {
		return nil, false
	}
	out := make([]string, 0, m.Live)
	for i := int(m.firstIdx); i < len(m.strEntries); i++ {
		k := m.strEntries[i].key
		if k != jTombstoneStr {
			out = append(out, k)
		}
	}
	return out, true
}

func (m *JMapCore) LiveValues() []any {
	if m == nil || m.Live == 0 {
		return []any{}
	}
	out := make([]any, 0, m.Live)
	m.ForEachEntry(func(_, v any) bool {
		out = append(out, v)
		return true
	})
	return out
}

// --- Simple Slice Iterator for `for-each` Loops ---

type SimpleIter interface {
	HasNext() bool
	Next() any
}

type sliceIter struct {
	items []any
	idx   int
}

func (s *sliceIter) IsJavaNil_() bool { return s == nil }
func (s *sliceIter) HasNext() bool    { return s.idx < len(s.items) }
func (s *sliceIter) Next() any {
	if s.idx >= len(s.items) {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	v := s.items[s.idx]
	s.idx++
	return v
}

type strSliceIter struct {
	strs []string
	idx  int
}

func (s *strSliceIter) HasNext() bool { return s.idx < len(s.strs) }
func (s *strSliceIter) Next() any {
	if s.idx >= len(s.strs) {
		panic(javalang.NewNoSuchElementException("iterator exhausted"))
	}
	v := s.strs[s.idx]
	s.idx++
	if v == NullString {
		return nil
	}
	return v
}

func Iterate(iterable any) SimpleIter {
	if iterable == nil {
		panic(javalang.NewNullPointerException())
	}
	if ssi, ok := iterable.(StrSliceIterable); ok {
		if nc, ok2 := iterable.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
			panic(javalang.NewNullPointerException())
		}
		if strs, ok2 := ssi.IterStrSlice_(); ok2 {
			return &strSliceIter{strs: strs}
		}
	}
	if tv, ok := iterable.(jTypedViewCursorSource); ok {
		if nc, ok2 := iterable.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
			panic(javalang.NewNullPointerException())
		}
		if _, ok2 := tv.typedViewLen_(); ok2 {
			if itProvider, ok3 := iterable.(interface{ Iterator() any }); ok3 {
				if si, ok4 := itProvider.Iterator().(SimpleIter); ok4 {
					return si
				}
			}
		}
	}
	if si, ok := iterable.(SliceIterable); ok {
		if nc, ok2 := iterable.(jJavaNilChecker); ok2 && nc.IsJavaNil_() {
			panic(javalang.NewNullPointerException())
		}
		return &sliceIter{items: si.IterSlice_()}
	}
	if isJavaNilAny(iterable) {
		panic(javalang.NewNullPointerException())
	}
	if itProvider, ok := iterable.(interface{ Iterator() Ext_java_util_Iterator }); ok {
		return itProvider.Iterator()
	}
	if itProvider, ok := iterable.(interface{ Iterator_0() Ext_java_util_Iterator }); ok {
		return itProvider.Iterator_0()
	}
	if itProvider, ok := iterable.(interface{ Iterator() any }); ok {
		if si, ok := itProvider.Iterator().(SimpleIter); ok {
			return si
		}
	}
	if itProvider, ok := iterable.(interface{ Iterator_0() any }); ok {
		if si, ok := itProvider.Iterator_0().(SimpleIter); ok {
			return si
		}
	}
	rv := reflect.ValueOf(iterable)
	if rv.Kind() == reflect.Slice {
		items := make([]any, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			items[i] = rv.Index(i).Interface()
		}
		return &sliceIter{items: items}
	}
	panic(javalang.NewUnsupportedOperationException(fmt.Sprintf("Iterate unsupported on %T", iterable)))
}
