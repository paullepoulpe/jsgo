// Clean-room Gson 2.9.1 Gson, GsonBuilder, TypeAdapter, and exceptions (Apache 2.0).

package gson

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"sync"

	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javautil"
)

// JsonParseException implements com.google.gson.JsonParseException.
type JsonParseException struct {
	javalang.RuntimeException
}

func (*JsonParseException) jsonParseException() {}

type IJsonParseException interface {
	javalang.IRuntimeException
	jsonParseException()
}

func NewJsonParseException(msg string) *JsonParseException {
	e := &JsonParseException{}
	javalang.InitThrowable(e, "com.google.gson.JsonParseException", msg, nil)
	return e
}

// JsonSyntaxException implements com.google.gson.JsonSyntaxException.
type JsonSyntaxException struct {
	JsonParseException
}

func NewJsonSyntaxException(msg string) *JsonSyntaxException {
	e := &JsonSyntaxException{}
	javalang.InitThrowable(e, "com.google.gson.JsonSyntaxException", msg, nil)
	return e
}

func NewJsonSyntaxExceptionCause(cause javalang.Throwable) *JsonSyntaxException {
	e := &JsonSyntaxException{}
	javalang.InitThrowable(e, "com.google.gson.JsonSyntaxException", "", cause)
	return e
}

// JsonIOException implements com.google.gson.JsonIOException.
type JsonIOException struct {
	JsonParseException
}

func NewJsonIOExceptionCause(cause javalang.Throwable) *JsonIOException {
	e := &JsonIOException{}
	javalang.InitThrowable(e, "com.google.gson.JsonIOException", "", cause)
	return e
}

// TypeAdapter implements com.google.gson.TypeAdapter.
type TypeAdapter interface {
	Write(out *JsonWriter, value any)
}

// GsonWritable allows types to stream their own JSON representation without reflection.
type GsonWritable interface {
	WriteGson(out *JsonWriter)
}

type typeHierarchyFactory struct {
	baseType    func(value any) bool
	typeAdapter TypeAdapter
}

// Gson implements com.google.gson.Gson.
type Gson struct {
	lenient            bool
	serializeNulls     bool
	htmlSafe           bool
	prettyPrinting     bool
	factoriesToBeAdded []*typeHierarchyFactory
}

// NewGson creates a default Gson instance (htmlSafe = true).
func NewGson() *Gson {
	return &Gson{htmlSafe: true}
}

// GsonBuilder implements com.google.gson.GsonBuilder.
type GsonBuilder struct {
	hierarchyFactories []*typeHierarchyFactory
	serializeNullsFlag bool
	prettyPrinting     bool
	escapeHtmlChars    bool
	lenientFlag        bool
}

func NewGsonBuilder() *GsonBuilder {
	return &GsonBuilder{escapeHtmlChars: true}
}

func (b *GsonBuilder) SerializeNulls() *GsonBuilder {
	b.serializeNullsFlag = true
	return b
}

func (b *GsonBuilder) DisableHtmlEscaping() *GsonBuilder {
	b.escapeHtmlChars = false
	return b
}

func (b *GsonBuilder) SetPrettyPrinting() *GsonBuilder {
	b.prettyPrinting = true
	return b
}

func (b *GsonBuilder) SetLenient() *GsonBuilder {
	b.lenientFlag = true
	return b
}

func (b *GsonBuilder) RegisterTypeHierarchyAdapter(baseType func(value any) bool, adapter TypeAdapter) *GsonBuilder {
	if adapter == nil {
		panic(javalang.NewIllegalArgumentException("typeAdapter == null"))
	}
	b.hierarchyFactories = append(b.hierarchyFactories, &typeHierarchyFactory{
		baseType:    baseType,
		typeAdapter: adapter,
	})
	return b
}

func (b *GsonBuilder) Create() *Gson {
	factories := make([]*typeHierarchyFactory, 0, len(b.hierarchyFactories))
	for i := len(b.hierarchyFactories) - 1; i >= 0; i-- {
		factories = append(factories, b.hierarchyFactories[i])
	}
	return &Gson{
		lenient:            b.lenientFlag,
		serializeNulls:     b.serializeNullsFlag,
		htmlSafe:           b.escapeHtmlChars,
		prettyPrinting:     b.prettyPrinting,
		factoriesToBeAdded: factories,
	}
}

// Streams_Write writes a JsonElement tree to writer.
func Streams_Write(element JsonElement, writer *JsonWriter) {
	if element == nil {
		writer.NullValue()
		return
	}
	switch v := element.(type) {
	case *JsonNull:
		writer.NullValue()
	case *JsonPrimitive:
		if v.IsNumber() {
			writer.ValueNumber(v.GetAsNumber())
		} else if v.IsBoolean() {
			writer.ValueBoolean(v.GetAsBoolean())
		} else {
			writer.ValueString(v.GetAsString())
		}
	case *JsonArray:
		writer.BeginArray()
		for e := range v.All() {
			Streams_Write(e, writer)
		}
		writer.EndArray()
	case *JsonObject:
		writer.BeginObject()
		for k, elem := range v.EntrySet().All() {
			writer.Name(k)
			Streams_Write(elem, writer)
		}
		writer.EndObject()
	default:
		panic(javalang.NewIllegalArgumentException("Unknown JsonElement type"))
	}
}

// ReadJsonElement reads a single JsonElement tree from r.
func ReadJsonElement(r *JsonReader) (JsonElement, javalang.Throwable) {
	p, err := r.peekCode()
	if err != nil {
		return nil, err
	}
	readScalar := func(tok uint8) (JsonElement, javalang.Throwable) {
		switch tok {
		case peekDoubleQuoted, peekSingleQuoted, peekUnquoted:
			var s string
			var sErr javalang.Throwable
			switch tok {
			case peekDoubleQuoted:
				s, sErr = r.readQuotedValue('"')
			case peekSingleQuoted:
				s, sErr = r.readQuotedValue('\'')
			default:
				s, sErr = r.readUnquotedValue()
			}
			if sErr != nil {
				return nil, sErr
			}
			r.peeked = peekNone
			r.incTopIndex()
			return NewJsonPrimitiveString(s), nil
		case peekLong, peekNumber:
			s := r.NextString()
			return NewJsonPrimitiveNumber(NewLazilyParsedNumber(s)), nil
		case peekTrue:
			r.peeked = peekNone
			r.incTopIndex()
			return NewJsonPrimitiveBoolean(true), nil
		case peekFalse:
			r.peeked = peekNone
			r.incTopIndex()
			return NewJsonPrimitiveBoolean(false), nil
		case peekNull:
			r.peeked = peekNone
			r.incTopIndex()
			return JsonNull_INSTANCE, nil
		default:
			jt, tErr := r.peekTokenErr()
			if tErr != nil {
				return nil, tErr
			}
			return nil, javalang.NewIllegalStateException("Unexpected token: " + jt.String())
		}
	}

	var current JsonElement
	switch p {
	case peekBeginArray:
		r.BeginArray()
		current = NewJsonArray()
	case peekBeginObject:
		r.BeginObject()
		current = NewJsonObject()
	default:
		return readScalar(p)
	}

	var stack []JsonElement
	for {
		for {
			vpCheck, pErr := r.peekCode()
			if pErr != nil {
				return nil, pErr
			}
			if vpCheck == peekEndObject || vpCheck == peekEndArray || vpCheck == peekEOF {
				break
			}
			var name string
			_, isObj := current.(*JsonObject)
			if isObj {
				name = r.NextName()
			}
			vp, vErr := r.peekCode()
			if vErr != nil {
				return nil, vErr
			}
			var val JsonElement
			nesting := false
			switch vp {
			case peekBeginArray:
				r.BeginArray()
				val = NewJsonArray()
				nesting = true
			case peekBeginObject:
				r.BeginObject()
				val = NewJsonObject()
				nesting = true
			default:
				val, vErr = readScalar(vp)
				if vErr != nil {
					return nil, vErr
				}
			}
			if isObj {
				current.(*JsonObject).Add(name, val)
			} else {
				current.(*JsonArray).Add(val)
			}
			if nesting {
				stack = append(stack, current)
				current = val
			}
		}
		if _, isArr := current.(*JsonArray); isArr {
			r.EndArray()
		} else {
			r.EndObject()
		}
		if len(stack) == 0 {
			return current, nil
		}
		current = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
	}
}

func (g *Gson) parseElement(json string) (elem JsonElement, prevPath string, retErr IJsonParseException) {
	defer func() {
		if rec := recover(); rec != nil {
			th := javalang.RecoverException(rec)
			switch e := th.(type) {
			case IJsonParseException:
				retErr = e
			case javaio.IIOException, javalang.IIllegalStateException:
				retErr = NewJsonSyntaxExceptionCause(e)
			default:
				panic(rec)
			}
		}
	}()
	r := NewJsonReaderString(json)
	r.SetLenient(true)
	_, emptyEOF, err := r.doPeek(true)
	if err != nil {
		return nil, "", NewJsonSyntaxExceptionCause(err)
	}
	if emptyEOF {
		return nil, "", nil
	}
	obj, err := ReadJsonElement(r)
	if err != nil {
		return nil, "", NewJsonSyntaxExceptionCause(err)
	}
	prevPath = r.GetPreviousPath()
	r.SetLenient(g.lenient)
	endTok, err := r.peekCode()
	if err != nil {
		if _, isMJE := err.(*MalformedJsonException); isMJE {
			return nil, "", NewJsonSyntaxExceptionCause(err)
		}
		if _, isIO := err.(javaio.IIOException); isIO {
			return nil, "", NewJsonIOExceptionCause(err)
		}
		return nil, "", NewJsonSyntaxExceptionCause(err)
	}
	if endTok != peekEOF {
		return nil, "", NewJsonSyntaxException("JSON document was not fully consumed.")
	}
	return obj, prevPath, nil
}

// FromJsonElement parses json as a JsonElement (nil on empty input).
func (g *Gson) FromJsonElement(json string) JsonElement {
	elem, _, err := g.parseElement(json)
	if err != nil {
		panic(err)
	}
	return elem
}

// FromJsonObject parses json as a *JsonObject.
func (g *Gson) FromJsonObject(json string) *JsonObject {
	elem, path, err := g.parseElement(json)
	if err != nil {
		panic(err)
	}
	if elem == nil {
		return nil
	}
	obj, ok := elem.(*JsonObject)
	if !ok {
		panic(NewJsonSyntaxException("Expected a com.google.gson.JsonObject at path " + path))
	}
	return obj
}

// FromJsonArray parses json as a *JsonArray.
func (g *Gson) FromJsonArray(json string) *JsonArray {
	elem, path, err := g.parseElement(json)
	if err != nil {
		panic(err)
	}
	if elem == nil {
		return nil
	}
	arr, ok := elem.(*JsonArray)
	if !ok {
		panic(NewJsonSyntaxException("Expected a com.google.gson.JsonArray at path " + path))
	}
	return arr
}

// FromJsonStringList parses a JSON array of strings into an *ArrayList[string] (ReportDecoder helper).
func (g *Gson) FromJsonStringList(json string) *javautil.ArrayList[string] {
	arr := g.FromJsonArray(json)
	if arr == nil {
		return nil
	}
	out := javautil.NewArrayListWithCapacity[string](arr.Size())
	for e := range arr.All() {
		if e.IsJsonNull() {
			out.Add("")
		} else {
			out.Add(e.GetAsString())
		}
	}
	return out
}

// FromJsonStringLongMap parses a JSON object into a *LinkedHashMap[string, int64] (ProductionInstrumentationReporter helper).
func (g *Gson) FromJsonStringLongMap(json string) *javautil.LinkedHashMap[string, int64] {
	obj := g.FromJsonObject(json)
	if obj == nil {
		return nil
	}
	out := javautil.NewLinkedHashMapWithCapacity[string, int64](obj.Size())
	for k, v := range obj.EntrySet().All() {
		out.Put(k, v.GetAsLong())
	}
	return out
}

// NewJsonWriter returns a configured JsonWriter writing to sb.
func (g *Gson) NewJsonWriter(sb *strings.Builder) *JsonWriter {
	jw := NewJsonWriterSink(sb)
	if g.prettyPrinting {
		jw.SetIndent("  ")
	}
	jw.SetHtmlSafe(g.htmlSafe)
	jw.SetLenient(g.lenient)
	jw.SetSerializeNulls(g.serializeNulls)
	return jw
}

// ToJson serializes src to a JSON string.
func (g *Gson) ToJson(src any) string {
	if src == nil || javalang.IsNil(src) {
		return "null"
	}
	var sb strings.Builder
	jw := g.NewJsonWriter(&sb)
	g.ToJsonWriter(src, jw)
	return sb.String()
}

// ToJsonWriter serializes src to writer using g's settings.
func (g *Gson) ToJsonWriter(src any, writer *JsonWriter) {
	oldLenient := writer.IsLenient()
	oldHtmlSafe := writer.IsHtmlSafe()
	oldNulls := writer.GetSerializeNulls()
	writer.SetLenient(true)
	writer.SetHtmlSafe(g.htmlSafe)
	writer.SetSerializeNulls(g.serializeNulls)
	defer func() {
		writer.SetLenient(oldLenient)
		writer.SetHtmlSafe(oldHtmlSafe)
		writer.SetSerializeNulls(oldNulls)
	}()
	g.writeValue(writer, src)
}

func (g *Gson) writeValue(out *JsonWriter, value any) {
	if value == nil || javalang.IsNil(value) {
		out.NullValue()
		return
	}
	for _, f := range g.factoriesToBeAdded {
		if f.baseType(value) {
			f.typeAdapter.Write(out, value)
			return
		}
	}
	if w, ok := value.(GsonWritable); ok {
		w.WriteGson(out)
		return
	}
	if e, ok := value.(JsonElement); ok {
		Streams_Write(e, out)
		return
	}
	switch v := value.(type) {
	case string:
		if v == "\xff\x00j2go_null" {
			out.NullValue()
			return
		}
		out.ValueString(v)
		return
	case *string:
		out.ValueStringPtr(v)
		return
	case bool:
		out.ValueBoolean(v)
		return
	case int:
		out.ValueLong(int64(v))
		return
	case int32:
		out.ValueLong(int64(v))
		return
	case int64:
		out.ValueLong(v)
		return
	case int16:
		out.ValueLong(int64(v))
		return
	case int8:
		out.ValueLong(int64(v))
		return
	case uint16:
		out.ValueString(javalang.StringFromChar(v))
		return
	case float64:
		if math.IsNaN(v) || math.IsInf(v, 0) {
			panic(javalang.NewIllegalArgumentException(javalang.DoubleToString(v) + " is not a valid double value as per JSON specification."))
		}
		out.ValueDouble(v)
		return
	case []string:
		out.BeginArray()
		for _, s := range v {
			out.ValueString(s)
		}
		out.EndArray()
		return
	}
	g.writeReflectValue(out, reflect.ValueOf(value))
}

type gsonMapWriter interface {
	gsonWriteEntries(g *Gson, out *JsonWriter)
}

type gsonSeqWriter interface {
	gsonWriteItems(g *Gson, out *JsonWriter)
}

func (g *Gson) writeReflectValue(out *JsonWriter, rv reflect.Value) {
	switch rv.Kind() {
	case reflect.String:
		if rv.String() == "\xff\x00j2go_null" {
			out.NullValue()
			return
		}
		out.ValueString(rv.String())
	case reflect.Bool:
		out.ValueBoolean(rv.Bool())
	case reflect.Int32:
		if s, ok := rv.Interface().(fmt.Stringer); ok && rv.Type().PkgPath() != "" {
			out.ValueString(s.String())
			return
		}
		out.ValueLong(rv.Int())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64:
		out.ValueLong(rv.Int())
	case reflect.Uint16:
		out.ValueString(javalang.StringFromChar(uint16(rv.Uint())))
	case reflect.Float64, reflect.Float32:
		d := rv.Float()
		if math.IsNaN(d) || math.IsInf(d, 0) {
			panic(javalang.NewIllegalArgumentException(javalang.DoubleToString(d) + " is not a valid double value as per JSON specification."))
		}
		out.ValueDouble(d)
	case reflect.Slice, reflect.Array:
		if rv.Kind() == reflect.Slice && rv.IsNil() {
			out.NullValue()
			return
		}
		out.BeginArray()
		for i := 0; i < rv.Len(); i++ {
			g.writeValue(out, rv.Index(i).Interface())
		}
		out.EndArray()
	case reflect.Pointer:
		if rv.IsNil() {
			out.NullValue()
			return
		}
		// Support *javautil.ArrayList[T], *javautil.LinkedHashSet[T], *javautil.LinkedHashMap[K, V] via ToSlice / Entries
		if m := rv.MethodByName("Entries"); m.IsValid() && m.Type().NumIn() == 0 && m.Type().NumOut() == 1 {
			entries := m.Call(nil)[0]
			if entries.Kind() == reflect.Slice {
				out.BeginObject()
				for i := 0; i < entries.Len(); i++ {
					ent := entries.Index(i)
					kField := ent.FieldByName("Key")
					vField := ent.FieldByName("Value")
					if kField.IsValid() && vField.IsValid() {
						out.Name(javalang.ToJavaString(kField.Interface()))
						g.writeValue(out, vField.Interface())
					}
				}
				out.EndObject()
				return
			}
		}
		if ts := rv.MethodByName("ToSlice"); ts.IsValid() && ts.Type().NumIn() == 0 && ts.Type().NumOut() == 1 {
			sl := ts.Call(nil)[0]
			out.BeginArray()
			for i := 0; i < sl.Len(); i++ {
				g.writeValue(out, sl.Index(i).Interface())
			}
			out.EndArray()
			return
		}
		if rv.Elem().Kind() == reflect.Struct {
			g.writeStruct(out, rv.Elem(), rv)
			return
		}
		g.writeValue(out, rv.Elem().Interface())
	case reflect.Interface:
		if rv.IsNil() {
			out.NullValue()
			return
		}
		g.writeValue(out, rv.Elem().Interface())
	case reflect.Struct:
		g.writeStruct(out, rv, reflect.Value{})
	default:
		panic(javalang.NewUnsupportedOperationException("Gson cannot serialize Go type " + rv.Type().String()))
	}
}

type boundField struct {
	name  string
	index []int
}

var boundFieldsCache sync.Map // reflect.Type -> []boundField

func getBoundFields(t reflect.Type) []boundField {
	if cached, ok := boundFieldsCache.Load(t); ok {
		return cached.([]boundField)
	}
	var fields []boundField
	seen := javautil.NewLinkedHashSet[string]()
	var collect func(t reflect.Type, prefix []int)
	collect = func(t reflect.Type, prefix []int) {
		var embedded []int
		for i := 0; i < t.NumField(); i++ {
			sf := t.Field(i)
			if sf.Anonymous && sf.Type.Kind() == reflect.Struct {
				if sf.IsExported() {
					embedded = append(embedded, i)
				}
				continue
			}
			tag := sf.Tag.Get("gson")
			if tag == "-" || !sf.IsExported() {
				continue
			}
			name := sf.Name
			if tag != "" {
				name = tag
			}
			if !seen.Add(name) {
				panic(javalang.NewIllegalArgumentException(t.String() + " declares multiple JSON fields named " + name))
			}
			fields = append(fields, boundField{name: name, index: append(append([]int{}, prefix...), i)})
		}
		for _, i := range embedded {
			collect(t.Field(i).Type, append(append([]int{}, prefix...), i))
		}
	}
	collect(t, nil)
	boundFieldsCache.Store(t, fields)
	return fields
}

func (g *Gson) writeStruct(out *JsonWriter, sv reflect.Value, self reflect.Value) {
	fields := getBoundFields(sv.Type())
	out.BeginObject()
	for _, f := range fields {
		fv := sv.FieldByIndex(f.index).Interface()
		if self.IsValid() && !javalang.IsNil(fv) {
			if fp := reflect.ValueOf(fv); fp.Kind() == reflect.Pointer && fp.Pointer() == self.Pointer() {
				continue
			}
		}
		out.Name(f.name)
		g.writeValue(out, fv)
	}
	out.EndObject()
}
