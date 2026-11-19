// Clean-room IO, NIO, Zip, Regex, Guava, Gson, XML/SAX, and JDK runtime helpers for j2gort.
package j2gort

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"unicode/utf16"

	"github.com/paullepoulpe/jsgo/transpiled/shims/gson"
	"github.com/paullepoulpe/jsgo/transpiled/shims/guava"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javalang"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javanio"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javaregex"
	"github.com/paullepoulpe/jsgo/transpiled/shims/javazip"
	shimsres "github.com/paullepoulpe/jsgo/transpiled/shims/resources"
)

var (
	_ = bytes.NewReader
	_ = fmt.Sprintf
	_ = strconv.Itoa
	_ = javaregex.Pattern_Quote
	_ = shimsres.ResourceBundle_GetBundle
)

func int8ToBytes(s []int8) []byte {
	out := make([]byte, len(s))
	for i, v := range s {
		out[i] = byte(v)
	}
	return out
}

func bytesToInt8(b []byte) []int8 {
	out := make([]int8, len(b))
	for i, v := range b {
		out[i] = int8(v)
	}
	return out
}

func ioErr(err error) {
	if err == nil {
		return
	}
	if t, ok := err.(javalang.Throwable); ok {
		panic(t)
	}
	switch err.(type) {
	case *javanio.CharacterCodingException, *javanio.MalformedInputException, *javanio.UnmappableCharacterException:
		ex := &Ext_java_nio_charset_CharacterCodingExceptionBase{}
		ex.BindExtSelf_(ex)
		InitUserThrowable(ex, "java.nio.charset.CharacterCodingException", err.Error())
		panic(ex)
	}
	if os.IsNotExist(err) {
		panic(javaio.NewFileNotFoundException(err.Error()))
	}
	panic(javaio.NewIOException(err.Error()))
}

// --- File & PrintStream Wrappers ---

func unwrapFile(f *Ext_java_io_File) *javaio.File {
	if f == nil {
		return nil
	}
	if jf, ok := f.Core_.(*javaio.File); ok {
		return jf
	}
	return nil
}

func wrapFile(jf *javaio.File) *Ext_java_io_File {
	if jf == nil {
		return nil
	}
	f := &Ext_java_io_File{Core_: jf}
	f.BindExtSelf_(f)
	return f
}

func wrapPrintStream(ps *javaio.PrintStream) *Ext_java_io_PrintStream {
	p := &Ext_java_io_PrintStream{}
	p.Core_ = ps
	p.BindExtSelf_(p)
	return p
}

func unwrapPrintStream(p *Ext_java_io_PrintStream) *javaio.PrintStream {
	CheckNotNullPtr(p)
	if ps, ok := p.Core_.(*javaio.PrintStream); ok && ps != nil {
		return ps
	}
	return javaio.System_Out
}

func wrapStdinStream() Ext_java_io_InputStream {
	return wrapInputStream(javaio.NewReaderInputStream(os.Stdin))
}

func wrapInputStream(in javaio.InputStream) Ext_java_io_InputStream {
	if in == nil {
		return nil
	}
	s := &Ext_java_io_InputStreamBase{Core_: in}
	s.BindExtSelf_(s)
	return s
}

func wrapOutputStream(out javaio.OutputStream) Ext_java_io_OutputStream {
	if out == nil {
		return nil
	}
	s := &Ext_java_io_OutputStreamBase{Core_: out}
	s.BindExtSelf_(s)
	return s
}

// --- Polymorphic InputStream / OutputStream / Reader / Writer Bridges ---

type extInputStreamAdapter struct {
	ext   Ext_java_io_InputStream
	inner javaio.InputStream
}

func (a *extInputStreamAdapter) Read(p []byte) (int, error) {
	return a.inner.Read(p)
}
func (a *extInputStreamAdapter) ReadByte() (byte, error) {
	return a.inner.ReadByte()
}
func (a *extInputStreamAdapter) Skip(n int64) (int64, error) {
	return a.inner.Skip(n)
}
func (a *extInputStreamAdapter) Available() int {
	return a.inner.Available()
}
func (a *extInputStreamAdapter) Close() error {
	a.ext.Close_0()
	return nil
}

func asInputStream(ext Ext_java_io_InputStream) javaio.InputStream {
	if ext == nil || javalang.IsNil(ext) {
		panic(javalang.NewNullPointerException())
	}
	base := ext.GetExt_java_io_InputStreamBase()
	var underlying javaio.InputStream
	switch c := base.Core_.(type) {
	case javaio.InputStream:
		underlying = c
	case []byte:
		bais := javaio.NewByteArrayInputStream(c)
		base.Core_ = bais
		underlying = bais
	case Ext_java_io_InputStream:
		underlying = asInputStream(c)
	case io.Reader:
		rs := javaio.NewReaderInputStream(c)
		base.Core_ = rs
		underlying = rs
	default:
		panic(fmt.Sprintf("UNIMPLEMENTED_EXT: unsupported InputStream Core_ %T", base.Core_))
	}
	// Always route Close() through ext.Close_0() so NonClosingInputStream is respected!
	return &extInputStreamAdapter{ext: ext, inner: underlying}
}

type extOutputStreamAdapter struct {
	ext   Ext_java_io_OutputStream
	inner javaio.OutputStream
}

func (a *extOutputStreamAdapter) Write(p []byte) (int, error) {
	if a.inner != nil {
		return a.inner.Write(p)
	}
	a.ext.Write_1_slice_int8(bytesToInt8(p))
	return len(p), nil
}
func (a *extOutputStreamAdapter) WriteByte(c byte) error {
	if a.inner != nil {
		return a.inner.WriteByte(c)
	}
	a.ext.Write_1_int32(int32(c))
	return nil
}
func (a *extOutputStreamAdapter) Flush() error {
	if a.inner != nil {
		return a.inner.Flush()
	}
	a.ext.Flush()
	return nil
}
func (a *extOutputStreamAdapter) Close() error {
	a.ext.Close_0()
	return nil
}

type jByteSinkCore interface {
	Write([]byte) (int, error)
	WriteByte(byte) error
}

type jByteSinkOutputAdapter struct {
	sink jByteSinkCore
}

func (a *jByteSinkOutputAdapter) Write(p []byte) (int, error) { return a.sink.Write(p) }
func (a *jByteSinkOutputAdapter) WriteByte(c byte) error      { return a.sink.WriteByte(c) }
func (a *jByteSinkOutputAdapter) Flush() error                { return nil }
func (a *jByteSinkOutputAdapter) Close() error                { return nil }

func unwrapOutputStream(ext Ext_java_io_OutputStream) (javaio.OutputStream, bool) {
	base := ext.GetExt_java_io_OutputStreamBase()
	if base == nil {
		return nil, false
	}
	ci, ok := any(ext).(interface{ JavaClassInfo_() (string, string) })
	if !ok {
		return nil, false
	}
	binName, _ := ci.JavaClassInfo_()
	switch binName {
	case "java.io.OutputStream",
		"java.io.ByteArrayOutputStream",
		"java.io.BufferedOutputStream",
		"java.io.FileOutputStream",
		"java.io.PrintStream",
		"java.util.zip.DeflaterOutputStream",
		"java.util.zip.GZIPOutputStream",
		"com.google.javascript.jscomp.serialization.FastGzipOutputStream":
		if os, ok := base.Core_.(javaio.OutputStream); ok && os != nil {
			return os, true
		}
	case "java.io.FilterOutputStream":
		if innerExt, ok := base.Core_.(Ext_java_io_OutputStream); ok && innerExt != nil && !javalang.IsNil(innerExt) {
			return unwrapOutputStream(innerExt)
		}
		if os, ok := base.Core_.(javaio.OutputStream); ok && os != nil {
			return os, true
		}
	case "com.google.javascript.jscomp.Compiler$NonClosingOutputStream":
		if innerExt, ok := base.Core_.(Ext_java_io_OutputStream); ok && innerExt != nil && !javalang.IsNil(innerExt) {
			if os, _ := unwrapOutputStream(innerExt); os != nil {
				return os, false
			}
		}
		if os, ok := base.Core_.(javaio.OutputStream); ok && os != nil {
			return os, false
		}
	case "com.google.protobuf.ByteString$Output":
		if os, ok := base.Core_.(javaio.OutputStream); ok && os != nil {
			return os, true
		}
		if sink, ok := base.Core_.(jByteSinkCore); ok && sink != nil {
			return &jByteSinkOutputAdapter{sink: sink}, true
		}
	}
	return nil, false
}

func asOutputStream(ext Ext_java_io_OutputStream) javaio.OutputStream {
	if ext == nil || javalang.IsNil(ext) {
		panic(javalang.NewNullPointerException())
	}
	if inner, directClose := unwrapOutputStream(ext); inner != nil {
		if directClose {
			return inner
		}
		return &extOutputStreamAdapter{ext: ext, inner: inner}
	}
	return &extOutputStreamAdapter{ext: ext}
}

func newBufferedExtOutputStream(target Ext_java_io_OutputStream) javaio.OutputStream {
	return javaio.NewBufferedOutputStream(asOutputStream(target))
}

func extInputStreamClose(core any) {
	switch c := core.(type) {
	case nil:
		return
	case Ext_java_io_InputStream:
		c.Close_0()
	case javaio.InputStream:
		ioErr(c.Close())
	case io.Closer:
		ioErr(c.Close())
	}
}

func extOutputStreamWriteByte(core any, b byte) {
	switch c := core.(type) {
	case Ext_java_io_OutputStream:
		c.Write_1_int32(int32(b))
	case javaio.OutputStream:
		ioErr(c.WriteByte(b))
	case io.Writer:
		_, err := c.Write([]byte{b})
		ioErr(err)
	default:
		panic(fmt.Sprintf("UNIMPLEMENTED_EXT: OutputStream.write(byte) on %T", core))
	}
}

func extOutputStreamWriteSlice(core any, s []int8) {
	switch c := core.(type) {
	case Ext_java_io_OutputStream:
		c.Write_1_slice_int8(s)
	case javaio.OutputStream:
		_, err := c.Write(int8ToBytes(s))
		ioErr(err)
	case io.Writer:
		_, err := c.Write(int8ToBytes(s))
		ioErr(err)
	default:
		panic(fmt.Sprintf("UNIMPLEMENTED_EXT: OutputStream.write(byte[]) on %T", core))
	}
}

func extOutputStreamFlush(core any) {
	switch c := core.(type) {
	case nil:
		return
	case Ext_java_io_OutputStream:
		c.Flush()
	case javaio.OutputStream:
		ioErr(c.Flush())
	}
}

func extOutputStreamClose(core any) {
	switch c := core.(type) {
	case nil:
		return
	case Ext_java_io_OutputStream:
		c.Close_0()
	case javaio.OutputStream:
		ioErr(c.Close())
	case io.Closer:
		ioErr(c.Close())
	}
}

func asReader(ext Ext_java_io_Reader) javaio.Reader {
	if ext == nil || javalang.IsNil(ext) {
		panic(javalang.NewNullPointerException())
	}
	base := ext.GetExt_java_io_ReaderBase()
	if r, ok := base.Core_.(javaio.Reader); ok && r != nil {
		return r
	}
	panic(fmt.Sprintf("UNIMPLEMENTED_EXT: unsupported Reader Core_ %T", base.Core_))
}

type extWriterAdapter struct {
	ext Ext_java_io_Writer
}

func (a *extWriterAdapter) Write(s string) {
	a.ext.Write_1_string(s)
}
func (a *extWriterAdapter) AppendCharSequence(s string) {
	a.ext.Write_1_string(s)
}
func (a *extWriterAdapter) AppendCharacter(c uint16) {
	a.ext.Write_1_string(javalang.StringFromChar(c))
}
func (a *extWriterAdapter) Flush() {
	a.ext.Flush()
}
func (a *extWriterAdapter) Close() {
	a.ext.Close_0()
}

func asWriter(ext Ext_java_io_Writer) javaio.Writer {
	if ext == nil || javalang.IsNil(ext) {
		panic(javalang.NewNullPointerException())
	}
	base := ext.GetExt_java_io_WriterBase()
	if w, ok := base.Core_.(javaio.Writer); ok && w != nil {
		if base.Self_ == nil || RefEquals(base.Self_, ext) {
			return w
		}
	}
	return &extWriterAdapter{ext: ext}
}

func extReaderReadOne(core any) int32 {
	r, ok := core.(javaio.Reader)
	if !ok || r == nil {
		return -1
	}
	var buf [1]uint16
	n := r.Read(buf[:], 0, 1)
	if n <= 0 {
		return -1
	}
	return int32(buf[0])
}

func extBufferedReaderReadLine(core any) string {
	br, ok := core.(*javaio.BufferedReader)
	if !ok || br == nil {
		return NullString
	}
	line, err := br.ReadLine()
	if err == io.EOF {
		return NullString
	}
	ioErr(err)
	return line
}

func extReaderClose(core any) {
	if r, ok := core.(javaio.Reader); ok && r != nil {
		ioErr(r.Close())
	}
}

func extWriterWriteString(core any, s string) {
	w, ok := core.(javaio.Writer)
	if !ok || w == nil {
		panic("UNIMPLEMENTED_EXT: Writer Core_ is nil")
	}
	w.Write(s)
}

func extWriterFlush(core any) {
	if w, ok := core.(javaio.Writer); ok && w != nil {
		w.Flush()
	}
}

func extWriterClose(core any) {
	if w, ok := core.(javaio.Writer); ok && w != nil {
		w.Close()
	}
}

func newInputStreamReaderCharset(in Ext_java_io_InputStream, cs Ext_java_nio_charset_Charset) javaio.Reader {
	return javanio.NewInputStreamReaderWithCharset(asInputStream(in), unwrapCharset(cs))
}

func newInputStreamReaderCharsetName(in Ext_java_io_InputStream, csName string) javaio.Reader {
	return javanio.NewInputStreamReaderWithCharset(asInputStream(in), unwrapCharset(charsetForName(csName)))
}

func newOutputStreamWriterCharset(out Ext_java_io_OutputStream, cs Ext_java_nio_charset_Charset) javaio.Writer {
	return javanio.NewOutputStreamWriterWithCharset(asOutputStream(out), unwrapCharset(cs))
}

func newOutputStreamWriterCharsetName(out Ext_java_io_OutputStream, csName string) javaio.Writer {
	return javanio.NewOutputStreamWriterWithCharset(asOutputStream(out), unwrapCharset(charsetForName(csName)))
}

func newOutputStreamWriterEncoder(out Ext_java_io_OutputStream, enc Ext_java_nio_charset_CharsetEncoder) javaio.Writer {
	return javanio.NewOutputStreamWriterWithEncoder(asOutputStream(out), unwrapCharsetEncoder(CheckNotNull(enc)))
}

// --- Zip & GZIP Helpers ---

func newGzipInputStream(in Ext_java_io_InputStream) *javazip.GZIPInputStream {
	gz, err := javazip.NewGZIPInputStream(asInputStream(in))
	ioErr(err)
	return gz
}

type gzipDeflaterBridge struct {
	out   io.Writer
	gz    *javazip.GZIPOutputStream
	owner *Ext_java_util_zip_DeflaterOutputStreamBase
}

func initGzipOutputStream(base *Ext_java_util_zip_DeflaterOutputStreamBase, out Ext_java_io_OutputStream) {
	w := asOutputStream(CheckNotNull(out))
	gz, err := javazip.NewGZIPOutputStream(w)
	ioErr(err)
	base.Core_ = gz
	bridge := &gzipDeflaterBridge{out: w, gz: gz, owner: base}
	base.F_Def = &Ext_java_util_zip_Deflater{Core_: bridge}
}

func zipDeflaterSetLevel(core any, level int) {
	if b, ok := core.(*gzipDeflaterBridge); ok && b != nil {
		gz, err := javazip.NewGZIPOutputStreamLevel(b.out, level)
		ioErr(err)
		b.gz = gz
		b.owner.Core_ = gz
	}
}

func gzipOutputStreamFinish(core any) {
	if gz, ok := core.(*javazip.GZIPOutputStream); ok && gz != nil {
		ioErr(gz.Finish())
	}
}

func newZipFileFromFile(f *Ext_java_io_File) *javazip.ZipFile {
	zf, err := javazip.NewZipFileFromFile(CheckNotNullPtr(unwrapFile(f)))
	ioErr(err)
	return zf
}

func newZipFileFromName(name string) *javazip.ZipFile {
	zf, err := javazip.NewZipFile(name)
	ioErr(err)
	return zf
}

func newZipInputStream(in Ext_java_io_InputStream, cs Ext_java_nio_charset_Charset) *javazip.ZipInputStream {
	if cs == nil || javalang.IsNil(cs) {
		return javazip.NewZipInputStream(asInputStream(in))
	}
	return javazip.NewZipInputStreamCharset(asInputStream(in), unwrapCharset(cs))
}

// --- NIO Charset, Path, Files, FileSystem, URI ---

func wrapCharset(cs *javanio.Charset) Ext_java_nio_charset_Charset {
	if cs == nil {
		return nil
	}
	c := &Ext_java_nio_charset_CharsetBase{Core_: cs}
	c.BindExtSelf_(c)
	return c
}

func wrapCharsetByName(name string) Ext_java_nio_charset_Charset {
	return wrapCharset(javanio.Charset_ForName(name))
}

func unwrapCharset(cs Ext_java_nio_charset_Charset) *javanio.Charset {
	if cs == nil || javalang.IsNil(cs) {
		return javanio.Charset_DefaultCharset()
	}
	if c, ok := cs.GetExt_java_nio_charset_CharsetBase().Core_.(*javanio.Charset); ok && c != nil {
		return c
	}
	return javanio.Charset_DefaultCharset()
}

func NewStringFromBytes(b []int8, off, length int32, csArg any) string {
	if b == nil {
		ThrowNPE()
	}
	if off < 0 || length < 0 || int(off) > len(b)-int(length) {
		panic(javalang.NewStringIndexOutOfBoundsException("String index out of range"))
	}
	sub := b[off : off+length]
	raw := make([]byte, len(sub))
	for i, v := range sub {
		raw[i] = byte(v)
	}
	if csArg == nil {
		return javanio.Charset_DefaultCharset().Decode(raw)
	}
	switch x := csArg.(type) {
	case Ext_java_nio_charset_Charset:
		if javalang.IsNil(x) {
			ThrowNPE()
		}
		return unwrapCharset(x).Decode(raw)
	case string:
		if x == NullString {
			ThrowNPE()
		}
		return unwrapCharset(charsetForName(x)).Decode(raw)
	default:
		s := CastString(csArg)
		if s == NullString {
			ThrowNPE()
		}
		return unwrapCharset(charsetForName(s)).Decode(raw)
	}
}

func charsetForName(name string) Ext_java_nio_charset_Charset {
	var cs *javanio.Charset
	func() {
		defer func() {
			if r := recover(); r != nil {
				switch e := r.(type) {
				case *javanio.IllegalCharsetNameException:
					panic(javalang.NewIllegalArgumentException(e.Error()))
				case *javanio.UnsupportedCharsetException:
					panic(javalang.NewIllegalArgumentException(e.Error()))
				default:
					panic(r)
				}
			}
		}()
		cs = javanio.Charset_ForName(name)
	}()
	return wrapCharset(cs)
}

func charsetIsSupported(name string) bool {
	var ok bool
	func() {
		defer func() {
			if r := recover(); r != nil {
				switch e := r.(type) {
				case *javanio.IllegalCharsetNameException:
					panic(javalang.NewIllegalArgumentException(e.Error()))
				default:
					panic(r)
				}
			}
		}()
		ok = javanio.Charset_IsSupported(name)
	}()
	return ok
}

func newCodingErrorAction(name string) *Ext_java_nio_charset_CodingErrorAction {
	act := javanio.CodingErrorAction_REPLACE
	switch name {
	case "IGNORE":
		act = javanio.CodingErrorAction_IGNORE
	case "REPORT":
		act = javanio.CodingErrorAction_REPORT
	}
	return &Ext_java_nio_charset_CodingErrorAction{Core_: act}
}

func unwrapCodingAction(a *Ext_java_nio_charset_CodingErrorAction) javanio.CodingErrorAction {
	CheckNotNullPtr(a)
	if act, ok := a.Core_.(javanio.CodingErrorAction); ok {
		return act
	}
	return javanio.CodingErrorAction_REPLACE
}

func wrapCharsetEncoder(enc *javanio.CharsetEncoder) Ext_java_nio_charset_CharsetEncoder {
	e := &Ext_java_nio_charset_CharsetEncoderBase{Core_: enc}
	e.BindExtSelf_(e)
	return e
}

func unwrapCharsetEncoder(enc Ext_java_nio_charset_CharsetEncoder) *javanio.CharsetEncoder {
	return enc.GetExt_java_nio_charset_CharsetEncoderBase().Core_.(*javanio.CharsetEncoder)
}

type jNioPath struct {
	p *javanio.Path
}

func (*jNioPath) Implements_j2gort_Ext_java_nio_file_Path()      {}
func (*jNioPath) Implements_j2gort_Ext_java_lang_Iterable()      {}
func (*jNioPath) Implements_j2gort_Ext_java_nio_file_Watchable() {}

func wrapPath(p *javanio.Path) Ext_java_nio_file_Path {
	if p == nil {
		return nil
	}
	return &jNioPath{p: p}
}

func unwrapPath(ep Ext_java_nio_file_Path) *javanio.Path {
	CheckNotNull(ep)
	return ep.(*jNioPath).p
}

func newExtPath(first string, more []string) Ext_java_nio_file_Path {
	for _, m := range more {
		CheckNotNullString(m)
	}
	var p *javanio.Path
	func() {
		defer func() {
			if r := recover(); r != nil {
				if ipe, ok := r.(*javanio.InvalidPathException); ok {
					panic(javalang.NewIllegalArgumentException(ipe.Error()))
				}
				panic(r)
			}
		}()
		p = javanio.Path_Of(first, more...)
	}()
	return wrapPath(p)
}

func (jp *jNioPath) GetParent() Ext_java_nio_file_Path {
	return wrapPath(jp.p.GetParent())
}
func (jp *jNioPath) Normalize() Ext_java_nio_file_Path {
	return wrapPath(jp.p.Normalize())
}
func (jp *jNioPath) Resolve_1_Ext_java_nio_file_Path(other Ext_java_nio_file_Path) Ext_java_nio_file_Path {
	return wrapPath(jp.p.Resolve(unwrapPath(other)))
}
func (jp *jNioPath) Resolve_1_string(other string) Ext_java_nio_file_Path {
	return wrapPath(jp.p.ResolveString(CheckNotNullString(other)))
}
func (jp *jNioPath) ResolveSibling_1_string(other string) Ext_java_nio_file_Path {
	return wrapPath(jp.p.ResolveSiblingString(CheckNotNullString(other)))
}
func (jp *jNioPath) Relativize(other Ext_java_nio_file_Path) Ext_java_nio_file_Path {
	return wrapPath(jp.p.Relativize(unwrapPath(other)))
}
func (jp *jNioPath) ToAbsolutePath() Ext_java_nio_file_Path {
	return wrapPath(jp.p.ToAbsolutePath())
}
func (jp *jNioPath) ToFile() *Ext_java_io_File {
	return wrapFile(jp.p.ToFile())
}
func (jp *jNioPath) CompareTo(other any) int32 {
	o := unwrapPath(AsIface[Ext_java_nio_file_Path](CheckNotNull(other)))
	return int32(jp.p.CompareTo(o))
}
func (jp *jNioPath) Equals(other any) bool {
	if o, ok := other.(*jNioPath); ok && o != nil {
		return jp.p.Equals(o.p)
	}
	return false
}
func (jp *jNioPath) HashCode() int32 { return jp.p.HashCode() }
func (jp *jNioPath) String() string  { return jp.p.String() }
func (jp *jNioPath) Iterator() any {
	n := jp.p.GetNameCount()
	items := make([]any, n)
	for i := 0; i < n; i++ {
		items[i] = wrapPath(jp.p.GetName(i))
	}
	return newSliceObjIter(items)
}
func (jp *jNioPath) ForEach(action Ext_java_util_function_Consumer) {
	ExtDefault_Iterable_ForEach(jp, action)
}

type jPathMatcher struct {
	pm javanio.PathMatcher
}

func (*jPathMatcher) Implements_j2gort_Ext_java_nio_file_PathMatcher() {}
func (m *jPathMatcher) Matches(p Ext_java_nio_file_Path) bool {
	return m.pm.Matches(unwrapPath(p))
}

func newExtPathMatcher(syntaxAndPattern string) Ext_java_nio_file_PathMatcher {
	return &jPathMatcher{pm: javanio.FileSystems_GetDefault().GetPathMatcher(syntaxAndPattern)}
}

func newDefaultFileSystem() Ext_java_nio_file_FileSystem {
	fs := &Ext_java_nio_file_FileSystemBase{Core_: javanio.FileSystems_GetDefault()}
	fs.BindExtSelf_(fs)
	return fs
}

func convertOpenOptions(opts []Ext_java_nio_file_OpenOption) []javanio.StandardOpenOption {
	var out []javanio.StandardOpenOption
	for _, o := range opts {
		if so, ok := o.(*Ext_java_nio_file_StandardOpenOption); ok && so != nil {
			switch so.Name_ {
			case "READ":
				out = append(out, javanio.StandardOpenOption_READ)
			case "WRITE":
				out = append(out, javanio.StandardOpenOption_WRITE)
			case "APPEND":
				out = append(out, javanio.StandardOpenOption_APPEND)
			case "TRUNCATE_EXISTING":
				out = append(out, javanio.StandardOpenOption_TRUNCATE_EXISTING)
			case "CREATE":
				out = append(out, javanio.StandardOpenOption_CREATE)
			case "CREATE_NEW":
				out = append(out, javanio.StandardOpenOption_CREATE_NEW)
			}
		}
	}
	return out
}

func nioFilesCreateDirectories(dir Ext_java_nio_file_Path) Ext_java_nio_file_Path {
	p, err := javanio.Files_CreateDirectories(unwrapPath(dir))
	ioErr(err)
	return wrapPath(p)
}

func nioFilesGetLastModifiedTime(p Ext_java_nio_file_Path) *Ext_java_nio_file_attribute_FileTime {
	ft, err := javanio.Files_GetLastModifiedTime(unwrapPath(p))
	ioErr(err)
	return &Ext_java_nio_file_attribute_FileTime{Core_: ft}
}

func nioFilesNewBufferedReader(p Ext_java_nio_file_Path, cs Ext_java_nio_charset_Charset) *Ext_java_io_BufferedReader {
	br, err := javanio.Files_NewBufferedReader(unwrapPath(p), unwrapCharset(cs))
	ioErr(err)
	r := &Ext_java_io_BufferedReader{}
	r.Core_ = br
	r.BindExtSelf_(r)
	return r
}

func nioFilesNewBufferedWriter(p Ext_java_nio_file_Path, cs Ext_java_nio_charset_Charset, opts []Ext_java_nio_file_OpenOption) *Ext_java_io_BufferedWriter {
	bw, err := javanio.Files_NewBufferedWriter(unwrapPath(p), unwrapCharset(cs), convertOpenOptions(opts)...)
	ioErr(err)
	w := &Ext_java_io_BufferedWriter{}
	w.Core_ = bw
	w.BindExtSelf_(w)
	return w
}

func nioFilesNewOutputStream(p Ext_java_nio_file_Path, opts []Ext_java_nio_file_OpenOption) Ext_java_io_OutputStream {
	fos, err := javanio.Files_NewOutputStream(unwrapPath(p), convertOpenOptions(opts)...)
	ioErr(err)
	return wrapOutputStream(fos)
}

func nioFilesReadString(p Ext_java_nio_file_Path, cs Ext_java_nio_charset_Charset) string {
	s, err := javanio.Files_ReadStringCharset(unwrapPath(p), unwrapCharset(cs))
	ioErr(err)
	return s
}

type jBasicFileAttrs struct {
	a *javanio.BasicFileAttributes
}

func (*jBasicFileAttrs) Implements_j2gort_Ext_java_nio_file_attribute_BasicFileAttributes() {}

type jFileVisitorAdapter struct {
	v Ext_java_nio_file_FileVisitor
}

func mapVisitResult(r *Ext_java_nio_file_FileVisitResult) javanio.FileVisitResult {
	if r == nil {
		return javanio.FileVisitResult_CONTINUE
	}
	switch r.Name_ {
	case "TERMINATE":
		return javanio.FileVisitResult_TERMINATE
	case "SKIP_SUBTREE":
		return javanio.FileVisitResult_SKIP_SUBTREE
	case "SKIP_SIBLINGS":
		return javanio.FileVisitResult_SKIP_SIBLINGS
	default:
		return javanio.FileVisitResult_CONTINUE
	}
}

func (a *jFileVisitorAdapter) PreVisitDirectory(_ *javanio.Path, _ *javanio.BasicFileAttributes) (javanio.FileVisitResult, error) {
	return javanio.FileVisitResult_CONTINUE, nil
}

func (a *jFileVisitorAdapter) VisitFile(file *javanio.Path, attrs *javanio.BasicFileAttributes) (javanio.FileVisitResult, error) {
	r := a.v.VisitFile(wrapPath(file), &jBasicFileAttrs{a: attrs})
	return mapVisitResult(r), nil
}

func (a *jFileVisitorAdapter) VisitFileFailed(file *javanio.Path, exc error) (javanio.FileVisitResult, error) {
	var ioEx javaio.IIOException
	if e, ok := exc.(javaio.IIOException); ok {
		ioEx = e
	} else if exc != nil {
		ioEx = javaio.NewIOException(exc.Error())
	}
	r := a.v.VisitFileFailed(wrapPath(file), ioEx)
	return mapVisitResult(r), nil
}

func (a *jFileVisitorAdapter) PostVisitDirectory(_ *javanio.Path, exc error) (javanio.FileVisitResult, error) {
	if exc != nil {
		return javanio.FileVisitResult_TERMINATE, exc
	}
	return javanio.FileVisitResult_CONTINUE, nil
}

func nioFilesWalkFileTree(start Ext_java_nio_file_Path, options Ext_java_util_Set, maxDepth int32, visitor Ext_java_nio_file_FileVisitor) Ext_java_nio_file_Path {
	follow := false
	if options != nil && !javalang.IsNil(options) {
		for it := Iterate(options); it.HasNext(); {
			if opt, ok := it.Next().(*Ext_java_nio_file_FileVisitOption); ok && opt != nil && opt.Name_ == "FOLLOW_LINKS" {
				follow = true
			}
		}
	}
	_, err := javanio.Files_WalkFileTree(unwrapPath(start), follow, int(maxDepth), &jFileVisitorAdapter{v: visitor})
	ioErr(err)
	return start
}

func newJavaURI(str string) *javanio.URI {
	u, err := javanio.NewURI(str)
	if err != nil {
		ex := &Ext_java_net_URISyntaxExceptionBase{}
		ex.BindExtSelf_(ex)
		InitUserThrowable(ex, "java.net.URISyntaxException", err.Error())
		panic(ex)
	}
	return u
}

// --- Guava Base, Escape, Hash, IO Helpers ---

func guavaAsciiEqualsIgnoreCase(a, b any) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil || javalang.IsNil(a) || javalang.IsNil(b) {
		return false
	}
	sa := StrVal(a)
	sb := StrVal(b)
	if sa == NullString || sb == NullString {
		return sa == sb
	}
	return guava.Ascii_EqualsIgnoreCase(sa, sb)
}

func wrapCharMatcher(cm *guava.CharMatcher) Ext_com_google_common_base_CharMatcher {
	obj := &Ext_com_google_common_base_CharMatcherBase{Core_: cm}
	obj.BindExtSelf_(obj)
	return obj
}

func unwrapCharMatcher(cm Ext_com_google_common_base_CharMatcher) *guava.CharMatcher {
	return cm.GetExt_com_google_common_base_CharMatcherBase().Core_.(*guava.CharMatcher)
}

func wrapCharMatcherAnyOf(seq any) Ext_com_google_common_base_CharMatcher {
	return wrapCharMatcher(guava.CharMatcher_AnyOf(StrVal(CheckNotNull(seq))))
}
func wrapCharMatcherInRange(lo, hi uint16) Ext_com_google_common_base_CharMatcher {
	return wrapCharMatcher(guava.CharMatcher_InRange(lo, hi))
}
func wrapCharMatcherIs(ch uint16) Ext_com_google_common_base_CharMatcher {
	return wrapCharMatcher(guava.CharMatcher_Is(ch))
}
func wrapCharMatcherIsNot(ch uint16) Ext_com_google_common_base_CharMatcher {
	return wrapCharMatcher(guava.CharMatcher_IsNot(ch))
}
func wrapCharMatcherJavaLetterOrDigit() Ext_com_google_common_base_CharMatcher {
	return wrapCharMatcher(guava.CharMatcher_JavaLetterOrDigit())
}
func wrapCharMatcherWhitespace() Ext_com_google_common_base_CharMatcher {
	return wrapCharMatcher(guava.CharMatcher_Whitespace())
}
func wrapCharMatcherAscii() Ext_com_google_common_base_CharMatcher {
	return wrapCharMatcher(guava.CharMatcher_Ascii())
}

func charMatcherOr(a *Ext_com_google_common_base_CharMatcherBase, b Ext_com_google_common_base_CharMatcher) Ext_com_google_common_base_CharMatcher {
	return wrapCharMatcher(a.Core_.(*guava.CharMatcher).Or(unwrapCharMatcher(b)))
}
func charMatcherMatchesAllOf(a *Ext_com_google_common_base_CharMatcherBase, seq any) bool {
	return a.Core_.(*guava.CharMatcher).MatchesAllOf(StrVal(seq))
}
func charMatcherLastIndexIn(a *Ext_com_google_common_base_CharMatcherBase, seq any) int32 {
	return int32(a.Core_.(*guava.CharMatcher).LastIndexIn(StrVal(seq)))
}
func charMatcherTrimFrom(a *Ext_com_google_common_base_CharMatcherBase, seq any) string {
	return a.Core_.(*guava.CharMatcher).TrimFrom(StrVal(seq))
}
func charMatcherCollapseFrom(a *Ext_com_google_common_base_CharMatcherBase, seq any, repl uint16) string {
	return a.Core_.(*guava.CharMatcher).CollapseFrom(StrVal(seq), repl)
}
func charMatcherMatchesBoxed(a *Ext_com_google_common_base_CharMatcherBase, ch any) bool {
	c := UnboxChar(Cast[*BoxChar](ch))
	return a.Core_.(*guava.CharMatcher).Matches(c)
}

type jGuavaFunc struct {
	fn func(any) any
}

func (*jGuavaFunc) Implements_j2gort_Ext_com_google_common_base_Function() {}
func (*jGuavaFunc) Implements_j2gort_Ext_java_util_function_Function()     {}
func (f *jGuavaFunc) Apply(x any) any                                      { return f.fn(x) }
func (f *jGuavaFunc) Equals(o any) bool                                    { return RefEquals(f, o) }
func (f *jGuavaFunc) AndThen(after Ext_java_util_function_Function) Ext_java_util_function_Function {
	return ExtDefault_Function_AndThen(f, after)
}
func (f *jGuavaFunc) Compose(before Ext_java_util_function_Function) Ext_java_util_function_Function {
	return ExtDefault_Function_Compose(f, before)
}

func newGuavaIdentityFunction() Ext_com_google_common_base_Function {
	return &jGuavaFunc{fn: func(x any) any { return x }}
}

type jGuavaPred struct {
	fn     func(any) bool
	nodeFn func(uint32) bool
}

var (
	guavaAlwaysTruePred_  Ext_com_google_common_base_Predicate = &jGuavaPred{fn: func(_ any) bool { return true }, nodeFn: func(_ uint32) bool { return true }}
	guavaAlwaysFalsePred_ Ext_com_google_common_base_Predicate = &jGuavaPred{fn: func(_ any) bool { return false }, nodeFn: func(_ uint32) bool { return false }}
)

func (*jGuavaPred) Implements_j2gort_Ext_com_google_common_base_Predicate() {}
func (*jGuavaPred) Implements_j2gort_Ext_java_util_function_Predicate()     {}
func (*jGuavaPred) HasIdentityEquals_()                                     {}
func (p *jGuavaPred) Apply(x any) bool                                      { return p.fn(x) }
func (p *jGuavaPred) Test(x any) bool                                       { return p.fn(x) }
func (p *jGuavaPred) TestNode_(n uint32) bool {
	if p == nil {
		ThrowNPE()
	}
	if p.nodeFn != nil {
		return p.nodeFn(n)
	}
	return p.fn(boxNodeById(n))
}
func (p *jGuavaPred) Equals(o any) bool                                     { return RefEquals(p, o) }
func (p *jGuavaPred) And(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	return ExtDefault_Predicate_And(p, o)
}
func (p *jGuavaPred) Or(o Ext_java_util_function_Predicate) Ext_java_util_function_Predicate {
	return ExtDefault_Predicate_Or(p, o)
}
func (p *jGuavaPred) Negate() Ext_java_util_function_Predicate {
	return ExtDefault_Predicate_Negate(p)
}

func newGuavaPredicate(fn func(any) bool) Ext_com_google_common_base_Predicate {
	return &jGuavaPred{fn: fn}
}

func newGuavaNotPredicate(target Ext_com_google_common_base_Predicate) Ext_com_google_common_base_Predicate {
	return &jGuavaPred{
		fn:     func(x any) bool { return !target.Apply(x) },
		nodeFn: func(n uint32) bool { return !PredTestNode(target, n) },
	}
}

type jMemoSupplier struct {
	once     sync.Once
	delegate Ext_com_google_common_base_Supplier
	val      any
}

func (*jMemoSupplier) Implements_j2gort_Ext_com_google_common_base_Supplier() {}
func (*jMemoSupplier) Implements_j2gort_Ext_java_util_function_Supplier()     {}
func (s *jMemoSupplier) Get() any {
	s.once.Do(func() { s.val = s.delegate.Get() })
	return s.val
}

func newMemoizingSupplier(d Ext_com_google_common_base_Supplier) Ext_com_google_common_base_Supplier {
	return &jMemoSupplier{delegate: d}
}

func newGuavaJoiner(sep string) *Ext_com_google_common_base_Joiner {
	j := &Ext_com_google_common_base_Joiner{Core_: guava.Joiner_On(sep)}
	j.BindExtSelf_(j)
	return j
}

func guavaJoinerJoinIter(j *Ext_com_google_common_base_Joiner, iter Ext_java_lang_Iterable) string {
	if ssi, ok := any(iter).(StrSliceIterable); ok {
		if ss, ok2 := ssi.IterStrSlice_(); ok2 {
			for _, s := range ss {
				CheckNotNullString(s)
			}
			return j.Core_.(*guava.Joiner).JoinStrings(ss)
		}
	}
	var parts []string
	for it := Iterate(iter); it.HasNext(); {
		v := it.Next()
		if v == nil || javalang.IsNil(v) {
			panic(javalang.NewNullPointerException())
		}
		parts = append(parts, StrVal(v))
	}
	return j.Core_.(*guava.Joiner).JoinStrings(parts)
}

func guavaJoinerJoinSlice(j *Ext_com_google_common_base_Joiner, parts []any) string {
	ss := make([]string, len(parts))
	for i, v := range parts {
		if v == nil || javalang.IsNil(v) {
			panic(javalang.NewNullPointerException())
		}
		ss[i] = StrVal(v)
	}
	return j.Core_.(*guava.Joiner).JoinStrings(ss)
}

func guavaJoinerJoinVarargs(j *Ext_com_google_common_base_Joiner, first, second any, rest []any) string {
	all := make([]any, 0, 2+len(rest))
	all = append(all, first, second)
	all = append(all, rest...)
	return guavaJoinerJoinSlice(j, all)
}

func newGuavaSplitterChar(sep uint16) *Ext_com_google_common_base_Splitter {
	s := &Ext_com_google_common_base_Splitter{Core_: guava.Splitter_OnChar(sep)}
	s.BindExtSelf_(s)
	return s
}

func newGuavaSplitterStr(sep string) *Ext_com_google_common_base_Splitter {
	s := &Ext_com_google_common_base_Splitter{Core_: guava.Splitter_On(sep)}
	s.BindExtSelf_(s)
	return s
}

func splitterOmitEmpty(s *Ext_com_google_common_base_Splitter) *Ext_com_google_common_base_Splitter {
	out := &Ext_com_google_common_base_Splitter{Core_: s.Core_.(*guava.Splitter).OmitEmptyStrings()}
	out.BindExtSelf_(out)
	return out
}

func splitterTrimResults(s *Ext_com_google_common_base_Splitter) *Ext_com_google_common_base_Splitter {
	out := &Ext_com_google_common_base_Splitter{Core_: s.Core_.(*guava.Splitter).TrimResults()}
	out.BindExtSelf_(out)
	return out
}

func splitterLimit(s *Ext_com_google_common_base_Splitter, lim int32) *Ext_com_google_common_base_Splitter {
	out := &Ext_com_google_common_base_Splitter{Core_: s.Core_.(*guava.Splitter).Limit(int(lim))}
	out.BindExtSelf_(out)
	return out
}

func SplitterSplitToListStr(s *Ext_com_google_common_base_Splitter, seq string) Ext_java_util_List {
	CheckNotNullPtr(s)
	CheckNotNullString(seq)
	parts := s.Core_.(*guava.Splitter).SplitToList(seq)
	return wrapImmStrSliceOwn(parts)
}

func SplitterSplitStr(s *Ext_com_google_common_base_Splitter, seq string) Ext_java_lang_Iterable {
	return SplitterSplitToListStr(s, seq)
}

func splitterSplitToList(s *Ext_com_google_common_base_Splitter, seq any) Ext_java_util_List {
	return SplitterSplitToListStr(s, StrVal(seq))
}

func splitterSplit(s *Ext_com_google_common_base_Splitter, seq any) Ext_java_lang_Iterable {
	return SplitterSplitToListStr(s, StrVal(seq))
}

type jToStringHelperData struct {
	name      string
	omitNulls bool
	fields    []struct {
		k string
		v any
	}
}

func newToStringHelper(name string) *Ext_com_google_common_base_MoreObjects_ToStringHelper {
	h := &Ext_com_google_common_base_MoreObjects_ToStringHelper{Core_: &jToStringHelperData{name: name}}
	h.BindExtSelf_(h)
	return h
}

func (h *Ext_com_google_common_base_MoreObjects_ToStringHelper) toStringHelperOmitNulls() {
	h.Core_.(*jToStringHelperData).omitNulls = true
}

func (h *Ext_com_google_common_base_MoreObjects_ToStringHelper) toStringHelperAdd(name string, val any) {
	d := h.Core_.(*jToStringHelperData)
	d.fields = append(d.fields, struct {
		k string
		v any
	}{k: name, v: ToAny(val)})
}

func (h *Ext_com_google_common_base_MoreObjects_ToStringHelper) toStringHelperBuild() string {
	d := h.Core_.(*jToStringHelperData)
	var sb strings.Builder
	sb.WriteString(d.name)
	sb.WriteByte('{')
	first := true
	for _, f := range d.fields {
		if d.omitNulls && f.v == nil {
			continue
		}
		if !first {
			sb.WriteString(", ")
		}
		first = false
		sb.WriteString(f.k)
		sb.WriteByte('=')
		sb.WriteString(StrVal(f.v))
	}
	sb.WriteByte('}')
	return sb.String()
}

func guavaPadStart(s string, minLen int32, padChar uint16) string {
	return guava.Strings_PadStart(s, int(minLen), padChar)
}

func splitCaseWords(format, s string) []string {
	switch format {
	case "LOWER_HYPHEN":
		return strings.Split(s, "-")
	case "LOWER_UNDERSCORE", "UPPER_UNDERSCORE":
		return strings.Split(s, "_")
	case "LOWER_CAMEL", "UPPER_CAMEL":
		if s == "" {
			return []string{""}
		}
		var words []string
		start := 0
		for i := 1; i < len(s); i++ {
			if s[i] >= 'A' && s[i] <= 'Z' {
				words = append(words, s[start:i])
				start = i
			}
		}
		words = append(words, s[start:])
		return words
	default:
		return []string{s}
	}
}

func capAsciiWord(w string) string {
	if w == "" {
		return ""
	}
	low := javalang.StringToLowerCaseASCII(w)
	first := low[0]
	if first >= 'a' && first <= 'z' {
		first = first - 'a' + 'A'
	}
	return string(first) + low[1:]
}

func guavaCaseFormatConvert(srcFmt, dstFmt, s string) string {
	if srcFmt == dstFmt || s == "" {
		return s
	}
	words := splitCaseWords(srcFmt, s)
	switch dstFmt {
	case "LOWER_HYPHEN":
		out := make([]string, len(words))
		for i, w := range words {
			out[i] = javalang.StringToLowerCaseASCII(w)
		}
		return strings.Join(out, "-")
	case "LOWER_UNDERSCORE":
		out := make([]string, len(words))
		for i, w := range words {
			out[i] = javalang.StringToLowerCaseASCII(w)
		}
		return strings.Join(out, "_")
	case "UPPER_UNDERSCORE":
		out := make([]string, len(words))
		for i, w := range words {
			out[i] = javalang.StringToUpperCaseASCII(w)
		}
		return strings.Join(out, "_")
	case "UPPER_CAMEL":
		var sb strings.Builder
		for _, w := range words {
			sb.WriteString(capAsciiWord(w))
		}
		return sb.String()
	case "LOWER_CAMEL":
		var sb strings.Builder
		for i, w := range words {
			if i == 0 {
				sb.WriteString(javalang.StringToLowerCaseASCII(w))
			} else {
				sb.WriteString(capAsciiWord(w))
			}
		}
		return sb.String()
	default:
		return s
	}
}

func guavaStandardSystemPropertyValue(propName string) string {
	switch propName {
	case "USER_DIR":
		wd, err := os.Getwd()
		if err != nil {
			return "."
		}
		return wd
	case "FILE_SEPARATOR":
		return "/"
	case "PATH_SEPARATOR":
		return ":"
	case "LINE_SEPARATOR":
		return "\n"
	case "JAVA_IO_TMPDIR":
		return os.TempDir()
	case "USER_HOME":
		h, _ := os.UserHomeDir()
		return h
	default:
		return NullString
	}
}

func newPercentEscaper(safeChars string, plusForSpace bool) *guava.PercentEscaper {
	return guava.NewPercentEscaper(safeChars, plusForSpace)
}

func escapePercentCore(core any, s string) string {
	if pe, ok := core.(*guava.PercentEscaper); ok && pe != nil {
		return pe.Escape(s)
	}
	panic("UNIMPLEMENTED_EXT: Escaper.escape without override")
}

func newEscaperFunction(esc Ext_com_google_common_escape_Escaper) Ext_com_google_common_base_Function {
	return &jGuavaFunc{fn: func(x any) any {
		if x == nil || javalang.IsNil(x) {
			return nil
		}
		return BoxString(esc.Escape_1_string(CastString(x)))
	}}
}

type jHashFunction struct {
	hf guava.HashFunction
}

func (*jHashFunction) Implements_j2gort_Ext_com_google_common_hash_HashFunction() {}
func (f *jHashFunction) NewHasher_0() Ext_com_google_common_hash_Hasher {
	return &jHasher{h: f.hf.NewHasher()}
}

type jHasher struct {
	h guava.Hasher
}

func (*jHasher) Implements_j2gort_Ext_com_google_common_hash_Hasher()        {}
func (*jHasher) Implements_j2gort_Ext_com_google_common_hash_PrimitiveSink() {}
func (hr *jHasher) PutInt(v int32) any {
	hr.h.PutInt(v)
	return hr
}
func (hr *jHasher) PutLong(v int64) any {
	hr.h.PutLong(v)
	return hr
}
func (hr *jHasher) PutUnencodedChars(cs any) any {
	hr.h.PutUnencodedChars(StrVal(CheckNotNull(cs)))
	return hr
}
func (hr *jHasher) PutString(cs any, charset Ext_java_nio_charset_Charset) any {
	hr.h.PutString(StrVal(CheckNotNull(cs)), unwrapCharset(charset))
	return hr
}
func (hr *jHasher) Hash() Ext_com_google_common_hash_HashCode {
	hc := hr.h.Hash()
	obj := &Ext_com_google_common_hash_HashCodeBase{Core_: hc}
	obj.BindExtSelf_(obj)
	return obj
}

func newFarmHashFunction() Ext_com_google_common_hash_HashFunction {
	return &jHashFunction{hf: guava.Hashing_FarmHashFingerprint64()}
}

func newMurmur3_128HashFunction() Ext_com_google_common_hash_HashFunction {
	return &jHashFunction{hf: guava.Hashing_Murmur3_128()}
}

func hashCodeAsInt(core any) int32 {
	return core.(guava.HashCode).AsInt()
}

func hashCodeAsLong(core any) int64 {
	return core.(guava.HashCode).AsLong()
}

func newGuavaBase64Encoding() Ext_com_google_common_io_BaseEncoding {
	b := &Ext_com_google_common_io_BaseEncodingBase{Core_: guava.BaseEncoding_Base64()}
	b.BindExtSelf_(b)
	return b
}

func base64Encode(b []int8) string {
	return guava.BaseEncoding_Base64().Encode(int8ToBytes(b))
}

func base64Decode(s string) []int8 {
	raw, err := guava.BaseEncoding_Base64().DecodeChecked(s)
	if err != nil {
		panic(javalang.NewIllegalArgumentException(err.Error()))
	}
	return bytesToInt8(raw)
}

func limitExtInputStream(in Ext_java_io_InputStream, limit int64) Ext_java_io_InputStream {
	if limit < 0 {
		panic(javalang.NewIllegalArgumentException("limit must be non-negative"))
	}
	lr := io.LimitReader(asInputStream(in), limit)
	return wrapInputStream(javaio.NewReaderInputStream(lr))
}

func readAllExtInputStreamBytes(in Ext_java_io_InputStream) []int8 {
	b, err := guava.ByteStreams_ToByteArray(asInputStream(in))
	ioErr(err)
	return bytesToInt8(b)
}

type jCharSourceSpec struct {
	content string
	file    *javaio.File
	cs      *javanio.Charset
}

func wrapStringCharSource(s string) Ext_com_google_common_io_CharSource {
	obj := &Ext_com_google_common_io_CharSourceBase{Core_: &jCharSourceSpec{content: s}}
	obj.BindExtSelf_(obj)
	return obj
}

func newGuavaFileCharSource(f *Ext_java_io_File, cs Ext_java_nio_charset_Charset) Ext_com_google_common_io_CharSource {
	obj := &Ext_com_google_common_io_CharSourceBase{
		Core_: &jCharSourceSpec{file: unwrapFile(f), cs: unwrapCharset(cs)},
	}
	obj.BindExtSelf_(obj)
	return obj
}

func charSourceRead(core any) string {
	spec := core.(*jCharSourceSpec)
	if spec.file == nil {
		return spec.content
	}
	b, err := os.ReadFile(spec.file.GetPath())
	ioErr(err)
	return spec.cs.Decode(b)
}

func charSourceOpenStream(core any) Ext_java_io_Reader {
	spec := core.(*jCharSourceSpec)
	if spec.file == nil {
		r := &Ext_java_io_StringReader{}
		r.Core_ = javaio.NewStringReader(spec.content)
		r.BindExtSelf_(r)
		return r
	}
	fis, err := javaio.NewFileInputStreamFromFile(spec.file)
	ioErr(err)
	r := &Ext_java_io_InputStreamReader{}
	r.Core_ = javanio.NewInputStreamReaderWithCharset(fis, spec.cs)
	r.BindExtSelf_(r)
	return r
}

type jCharSinkSpec struct {
	file   *javaio.File
	cs     *javanio.Charset
	append bool
}

func newGuavaFileCharSink(f *Ext_java_io_File, cs Ext_java_nio_charset_Charset, modes []*Ext_com_google_common_io_FileWriteMode) Ext_com_google_common_io_CharSink {
	app := false
	for _, m := range modes {
		if m != nil && m.Name_ == "APPEND" {
			app = true
		}
	}
	obj := &Ext_com_google_common_io_CharSinkBase{
		Core_: &jCharSinkSpec{file: unwrapFile(f), cs: unwrapCharset(cs), append: app},
	}
	obj.BindExtSelf_(obj)
	return obj
}

func charSinkWrite(core any, content string) {
	spec := core.(*jCharSinkSpec)
	data := spec.cs.Encode(content)
	flags := os.O_WRONLY | os.O_CREATE
	if spec.append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}
	f, err := os.OpenFile(spec.file.GetPath(), flags, 0o666)
	ioErr(err)
	_, wErr := f.Write(data)
	cErr := f.Close()
	ioErr(wErr)
	ioErr(cErr)
}

func newGuavaFileBufferedReader(f *Ext_java_io_File, cs Ext_java_nio_charset_Charset) *Ext_java_io_BufferedReader {
	fis, err := javaio.NewFileInputStreamFromFile(unwrapFile(f))
	ioErr(err)
	br := &Ext_java_io_BufferedReader{}
	br.Core_ = javaio.NewBufferedReader(javanio.NewInputStreamReaderWithCharset(fis, unwrapCharset(cs)))
	br.BindExtSelf_(br)
	return br
}

func guavaFileToByteArray(f *Ext_java_io_File) []int8 {
	b, err := guava.Files_ToByteArray(unwrapFile(f))
	ioErr(err)
	return bytesToInt8(b)
}

func guavaFileWriteBytes(b []int8, f *Ext_java_io_File) {
	ioErr(guava.Files_WriteBytes(int8ToBytes(b), unwrapFile(f)))
}

func readAllFromReadable(r Ext_java_lang_Readable) string {
	if rd, ok := r.(Ext_java_io_Reader); ok {
		return guava.CharStreams_ToString(asReader(rd))
	}
	panic(fmt.Sprintf("UNIMPLEMENTED_EXT: CharStreams.toString on %T", r))
}

func readLinesFromReadable(r Ext_java_lang_Readable) Ext_java_util_List {
	rd, ok := r.(Ext_java_io_Reader)
	if !ok {
		panic(fmt.Sprintf("UNIMPLEMENTED_EXT: CharStreams.readLines on %T", r))
	}
	lines := guava.CharStreams_ReadLines(asReader(rd))
	items := make([]any, len(lines))
	for i, ln := range lines {
		items[i] = BoxString(ln)
	}
	return wrapArrayList(items)
}

// --- Gson & JsonElement Helpers + GsonWritable Implementations ---

type jGsonAdapterBridge struct {
	ext Ext_com_google_gson_TypeAdapter
}

func (b *jGsonAdapterBridge) Write(out *gson.JsonWriter, val any) {
	jw := &Ext_com_google_gson_stream_JsonWriter{Core_: out}
	jw.BindExtSelf_(jw)
	b.ext.Write(jw, val)
}

type jGsonBuilderCore struct {
	b *gson.GsonBuilder
}

func newGsonBuilderCore() *jGsonBuilderCore {
	return &jGsonBuilderCore{b: gson.NewGsonBuilder()}
}

func newDefaultGson() *gson.Gson {
	return gson.NewGson()
}

func gsonBuilderSetPretty(gb *Ext_com_google_gson_GsonBuilder) {
	gb.Core_.(*jGsonBuilderCore).b.SetPrettyPrinting()
}

func gsonBuilderDisableHtml(gb *Ext_com_google_gson_GsonBuilder) {
	gb.Core_.(*jGsonBuilderCore).b.DisableHtmlEscaping()
}

func gsonBuilderRegisterHierarchyAdapter(gb *Ext_com_google_gson_GsonBuilder, baseClass *javalang.Class, adapter any) {
	CheckNotNullPtr(baseClass)
	ta := AsIface[Ext_com_google_gson_TypeAdapter](CheckNotNull(adapter))
	gb.Core_.(*jGsonBuilderCore).b.RegisterTypeHierarchyAdapter(
		func(v any) bool { return CheckTypePredicate(baseClass, v) },
		&jGsonAdapterBridge{ext: ta},
	)
}

func gsonBuilderCreate(gb *Ext_com_google_gson_GsonBuilder) *Ext_com_google_gson_Gson {
	g := &Ext_com_google_gson_Gson{Core_: gb.Core_.(*jGsonBuilderCore).b.Create()}
	g.BindExtSelf_(g)
	return g
}

func unwrapGson(g *Ext_com_google_gson_Gson) *gson.Gson {
	CheckNotNullPtr(g)
	if gs, ok := g.Core_.(*gson.Gson); ok && gs != nil {
		return gs
	}
	return gson.NewGson()
}

func gsonToJson(g *Ext_com_google_gson_Gson, v any) string {
	if v == nil || javalang.IsNil(v) {
		return "null"
	}
	if je, ok := v.(Ext_com_google_gson_JsonElement); ok {
		return unwrapGson(g).ToJson(unwrapJsonElement(je))
	}
	return unwrapGson(g).ToJson(ToAny(v))
}

func gsonToJsonWriter(g *Ext_com_google_gson_Gson, v any, w *Ext_com_google_gson_stream_JsonWriter) {
	jw := gsonWriter(CheckNotNullPtr(w))
	if v == nil || javalang.IsNil(v) {
		jw.NullValue()
		return
	}
	if je, ok := v.(Ext_com_google_gson_JsonElement); ok {
		gson.Streams_Write(unwrapJsonElement(je), jw)
		return
	}
	unwrapGson(g).ToJsonWriter(ToAny(v), jw)
}

func wrapGsonPanic(fn func() any) (res any) {
	defer func() {
		if r := recover(); r != nil {
			if gp, ok := r.(gson.IJsonParseException); ok {
				ex := &Ext_com_google_gson_JsonParseExceptionBase{}
				ex.BindExtSelf_(ex)
				InitUserThrowable(ex, "com.google.gson.JsonParseException", gp.GetMessage(), gp)
				panic(ex)
			}
			panic(r)
		}
	}()
	return fn()
}

func gsonFromJsonClass(g *Ext_com_google_gson_Gson, jsonStr string, cls *javalang.Class) any {
	if jsonStr == NullString {
		return nil
	}
	CheckNotNullPtr(cls)
	gs := unwrapGson(g)
	return wrapGsonPanic(func() any {
		switch cls.GetName() {
		case "com.google.gson.JsonObject":
			jo := gs.FromJsonObject(jsonStr)
			if jo == nil {
				return nil
			}
			return wrapJsonObject(jo)
		case "com.google.gson.JsonArray":
			ja := gs.FromJsonArray(jsonStr)
			if ja == nil {
				return nil
			}
			return wrapJsonArray(ja)
		case "com.google.gson.JsonElement":
			je := gs.FromJsonElement(jsonStr)
			if je == nil {
				return nil
			}
			return wrapJsonElement(je)
		default:
			panic("UNIMPLEMENTED_EXT: Gson.fromJson(String, " + cls.GetName() + ")")
		}
	})
}

func gsonFromJsonType(g *Ext_com_google_gson_Gson, jsonStr string, typ any) any {
	if jsonStr == NullString {
		return nil
	}
	if cls, ok := typ.(*javalang.Class); ok {
		return gsonFromJsonClass(g, jsonStr, cls)
	}
	ts := fmt.Sprint(typ)
	gs := unwrapGson(g)
	return wrapGsonPanic(func() any {
		if strings.Contains(ts, "Map<java.lang.String, java.lang.Long>") || strings.Contains(ts, "Map<String, Long>") {
			m := gs.FromJsonStringLongMap(jsonStr)
			if m == nil {
				return nil
			}
			mc := NewJMapCore(false)
			for k, v := range m.All() {
				mc.Put(BoxString(k), ValueOfLong(v))
			}
			return wrapLinkedHashMap(mc)
		}
		if strings.Contains(ts, "List<java.lang.String>") || strings.Contains(ts, "List<String>") {
			l := gs.FromJsonStringList(jsonStr)
			if l == nil {
				return nil
			}
			s := l.ToSlice()
			items := make([]any, len(s))
			for i, str := range s {
				items[i] = BoxString(str)
			}
			return wrapArrayList(items)
		}
		panic("UNIMPLEMENTED_EXT: Gson.fromJson(String, Type=" + ts + ")")
	})
}

func newGsonJsonArray() *gson.JsonArray   { return gson.NewJsonArray() }
func newGsonJsonObject() *gson.JsonObject { return gson.NewJsonObject() }
func newGsonJsonPrimitive(v any) *gson.JsonPrimitive {
	if s, ok := v.(string); ok {
		return gson.NewJsonPrimitiveString(CheckNotNullString(s))
	}
	if b, ok := v.(bool); ok {
		return gson.NewJsonPrimitiveBoolean(b)
	}
	if bb, ok := v.(*BoxBoolean); ok {
		return gson.NewJsonPrimitiveBoolean(UnboxBoolean(bb))
	}
	if n, ok := v.(Number); ok {
		return gson.NewJsonPrimitiveLong(n.LongValue())
	}
	return gson.NewJsonPrimitiveString(StrVal(v))
}

func unwrapJsonElement(je Ext_com_google_gson_JsonElement) gson.JsonElement {
	if je == nil || javalang.IsNil(je) {
		return gson.JsonNull_INSTANCE
	}
	if e, ok := je.GetExt_com_google_gson_JsonElementBase().Core_.(gson.JsonElement); ok && e != nil {
		return e
	}
	return gson.JsonNull_INSTANCE
}

func unwrapJsonArray(ja *Ext_com_google_gson_JsonArray) *gson.JsonArray {
	CheckNotNullPtr(ja)
	return ja.Core_.(*gson.JsonArray)
}

func unwrapJsonObject(jo *Ext_com_google_gson_JsonObject) *gson.JsonObject {
	CheckNotNullPtr(jo)
	return jo.Core_.(*gson.JsonObject)
}

func unwrapJsonPrimitive(jp *Ext_com_google_gson_JsonPrimitive) *gson.JsonPrimitive {
	CheckNotNullPtr(jp)
	return jp.Core_.(*gson.JsonPrimitive)
}

func wrapJsonObject(jo *gson.JsonObject) *Ext_com_google_gson_JsonObject {
	if jo == nil {
		return nil
	}
	o := &Ext_com_google_gson_JsonObject{}
	o.Core_ = jo
	o.BindExtSelf_(o)
	return o
}

func wrapJsonArray(ja *gson.JsonArray) *Ext_com_google_gson_JsonArray {
	if ja == nil {
		return nil
	}
	a := &Ext_com_google_gson_JsonArray{}
	a.Core_ = ja
	a.BindExtSelf_(a)
	return a
}

func wrapJsonPrimitive(jp *gson.JsonPrimitive) *Ext_com_google_gson_JsonPrimitive {
	if jp == nil {
		return nil
	}
	p := &Ext_com_google_gson_JsonPrimitive{}
	p.Core_ = jp
	p.BindExtSelf_(p)
	return p
}

func wrapJsonElement(je gson.JsonElement) Ext_com_google_gson_JsonElement {
	if je == nil {
		return nil
	}
	switch x := je.(type) {
	case *gson.JsonObject:
		return wrapJsonObject(x)
	case *gson.JsonArray:
		return wrapJsonArray(x)
	case *gson.JsonPrimitive:
		return wrapJsonPrimitive(x)
	case *gson.JsonNull:
		n := &Ext_com_google_gson_JsonNull{}
		n.Core_ = x
		n.BindExtSelf_(n)
		return n
	default:
		b := &Ext_com_google_gson_JsonElementBase{Core_: je}
		b.BindExtSelf_(b)
		return b
	}
}

func gsonJsonArrayAdd(ja *Ext_com_google_gson_JsonArray, elem Ext_com_google_gson_JsonElement) {
	unwrapJsonArray(ja).Add(unwrapJsonElement(elem))
}

func gsonJsonArrayIterator(ja *Ext_com_google_gson_JsonArray) Ext_java_util_Iterator {
	a := unwrapJsonArray(ja)
	n := a.Size()
	items := make([]any, n)
	for i := 0; i < n; i++ {
		items[i] = wrapJsonElement(a.Get(i))
	}
	return newSliceObjIter(items)
}

func gsonJsonArrayForEach(ja *Ext_com_google_gson_JsonArray, action Ext_java_util_function_Consumer) {
	a := unwrapJsonArray(ja)
	n := a.Size()
	for i := 0; i < n; i++ {
		action.Accept(wrapJsonElement(a.Get(i)))
	}
}

func gsonJsonObjectAdd(jo *Ext_com_google_gson_JsonObject, key string, val Ext_com_google_gson_JsonElement) {
	unwrapJsonObject(jo).Add(key, unwrapJsonElement(val))
}

func gsonJsonObjectEntrySet(jo *Ext_com_google_gson_JsonObject) Ext_java_util_Set {
	o := unwrapJsonObject(jo)
	mc := NewJMapCore(false)
	for k, v := range o.EntrySet().All() {
		mc.Put(BoxString(k), wrapJsonElement(v))
	}
	return newLiveMapEntrySet(mc)
}

func newGsonJsonReader(r Ext_java_io_Reader) *gson.JsonReader {
	return gson.NewJsonReader(asReader(r))
}

func newGsonJsonWriter(w Ext_java_io_Writer) *gson.JsonWriter {
	return gson.NewJsonWriter(asWriter(w))
}

func gsonReader(r *Ext_com_google_gson_stream_JsonReader) *gson.JsonReader {
	CheckNotNullPtr(r)
	return r.Core_.(*gson.JsonReader)
}

func gsonWriter(w *Ext_com_google_gson_stream_JsonWriter) *gson.JsonWriter {
	CheckNotNullPtr(w)
	return w.Core_.(*gson.JsonWriter)
}

func writeGsonAny(out *gson.JsonWriter, v any) {
	v = ToAny(v)
	if v == nil {
		out.NullValue()
		return
	}
	switch x := v.(type) {
	case string:
		if x == NullString {
			out.NullValue()
		} else {
			out.ValueString(x)
		}
	case bool:
		out.ValueBoolean(x)
	case *BoxBoolean:
		out.ValueBoolean(x.Val)
	case *BoxInt:
		out.ValueInt(x.Val)
	case *BoxLong:
		out.ValueLong(x.Val)
	case *BoxByte:
		out.ValueInt(int32(x.Val))
	case *BoxShort:
		out.ValueInt(int32(x.Val))
	case *BoxFloat:
		out.ValueDouble(float64(x.Val))
	case *BoxDouble:
		out.ValueDouble(x.Val)
	case *BoxChar:
		out.ValueString(javalang.StringFromChar(x.Val))
	case EnumAccessor:
		out.ValueString(x.GetEnumBase().Name_)
	case Ext_com_google_gson_JsonElement:
		gson.Streams_Write(unwrapJsonElement(x), out)
	case gson.GsonWritable:
		x.WriteGson(out)
	case Ext_java_util_Map:
		out.BeginObject()
		for it := Iterate(x.EntrySet()); it.HasNext(); {
			e := AsIface[Ext_java_util_Map_Entry](it.Next())
			out.Name(StrVal(e.GetKey()))
			writeGsonAny(out, e.GetValue())
		}
		out.EndObject()
	case Ext_java_lang_Iterable:
		out.BeginArray()
		for it := Iterate(x); it.HasNext(); {
			writeGsonAny(out, it.Next())
		}
		out.EndArray()
	default:
		// Struct reflection via default Gson
		out.JsonValue(gson.NewGson().ToJson(v))
	}
}

func (b *BoxBoolean) WriteGson(out *gson.JsonWriter) { out.ValueBoolean(b.Val) }
func (b *BoxByte) WriteGson(out *gson.JsonWriter)    { out.ValueInt(int32(b.Val)) }
func (b *BoxShort) WriteGson(out *gson.JsonWriter)   { out.ValueInt(int32(b.Val)) }
func (b *BoxChar) WriteGson(out *gson.JsonWriter)    { out.ValueString(javalang.StringFromChar(b.Val)) }
func (b *BoxInt) WriteGson(out *gson.JsonWriter)     { out.ValueInt(b.Val) }
func (b *BoxLong) WriteGson(out *gson.JsonWriter)    { out.ValueLong(b.Val) }
func (b *BoxFloat) WriteGson(out *gson.JsonWriter)   { out.ValueDouble(float64(b.Val)) }
func (b *BoxDouble) WriteGson(out *gson.JsonWriter)  { out.ValueDouble(b.Val) }
func (e *EnumBase) WriteGson(out *gson.JsonWriter)   { out.ValueString(e.Name_) }

func (l *JListCore) WriteGson(out *gson.JsonWriter) {
	out.BeginArray()
	for _, item := range l.Items {
		writeGsonAny(out, item)
	}
	out.EndArray()
}

func (m *JMapCore) WriteGson(out *gson.JsonWriter) {
	out.BeginObject()
	m.ForEachEntry(func(k, v any) bool {
		out.Name(StrVal(k))
		writeGsonAny(out, v)
		return true
	})
	out.EndObject()
}

func (b *Ext_java_util_AbstractCollectionBase) WriteGson(out *gson.JsonWriter) {
	var target Ext_java_lang_Iterable = b
	if b.Self_ != nil {
		target = b.Self_
	}
	out.BeginArray()
	for it := Iterate(target); it.HasNext(); {
		writeGsonAny(out, it.Next())
	}
	out.EndArray()
}

func (b *Ext_java_util_AbstractMapBase) WriteGson(out *gson.JsonWriter) {
	var target Ext_java_util_Map = b
	if b.Self_ != nil {
		target = b.Self_
	}
	out.BeginObject()
	for it := Iterate(target.EntrySet()); it.HasNext(); {
		e := AsIface[Ext_java_util_Map_Entry](it.Next())
		out.Name(StrVal(e.GetKey()))
		writeGsonAny(out, e.GetValue())
	}
	out.EndObject()
}

func (b *Ext_com_google_common_collect_ImmutableCollectionBase) WriteGson(out *gson.JsonWriter) {
	var target Ext_java_lang_Iterable = b
	if b.Self_ != nil {
		target = b.Self_
	}
	out.BeginArray()
	for it := Iterate(target); it.HasNext(); {
		writeGsonAny(out, it.Next())
	}
	out.EndArray()
}

func (b *Ext_com_google_common_collect_ImmutableMapBase) WriteGson(out *gson.JsonWriter) {
	var target Ext_java_util_Map = b
	if b.Self_ != nil {
		target = b.Self_
	}
	out.BeginObject()
	for it := Iterate(target.EntrySet()); it.HasNext(); {
		e := AsIface[Ext_java_util_Map_Entry](it.Next())
		out.Name(StrVal(e.GetKey()))
		writeGsonAny(out, e.GetValue())
	}
	out.EndObject()
}

// --- Clean-Room SAX XML Reader Adapter (over Go encoding/xml) ---

type jSaxAttrs struct {
	attrs []xml.Attr
}

func (*jSaxAttrs) Implements_j2gort_Ext_org_xml_sax_Attributes() {}

func (a *jSaxAttrs) GetValue_1_string(qName string) string {
	for _, attr := range a.attrs {
		full := attr.Name.Local
		if attr.Name.Space != "" && !strings.HasPrefix(attr.Name.Space, "http") {
			full = attr.Name.Space + ":" + attr.Name.Local
		}
		if full == qName || attr.Name.Local == qName {
			return attr.Value
		}
	}
	return NullString
}

type jSaxXMLReader struct {
	entityResolver Ext_org_xml_sax_EntityResolver
	contentHandler Ext_org_xml_sax_ContentHandler
}

func (*jSaxXMLReader) Implements_j2gort_Ext_org_xml_sax_XMLReader() {}

func (xr *jSaxXMLReader) SetEntityResolver(er Ext_org_xml_sax_EntityResolver) {
	xr.entityResolver = er
}

func (xr *jSaxXMLReader) SetContentHandler(ch Ext_org_xml_sax_ContentHandler) {
	xr.contentHandler = ch
}

func throwSAXException(msg string) {
	ex := &Ext_org_xml_sax_SAXExceptionBase{}
	ex.BindExtSelf_(ex)
	InitUserThrowable(ex, "org.xml.sax.SAXException", msg)
	panic(ex)
}

func (xr *jSaxXMLReader) Parse_1_Ext_org_xml_sax_InputSource(src *Ext_org_xml_sax_InputSource) {
	CheckNotNullPtr(src)
	var raw []byte
	switch c := src.Core_.(type) {
	case Ext_java_io_InputStream:
		b, err := io.ReadAll(asInputStream(c))
		ioErr(err)
		raw = b
	case Ext_java_io_Reader:
		s := guava.CharStreams_ToString(asReader(c))
		raw = []byte(s)
	default:
		throwSAXException("Unsupported InputSource")
	}

	dec := xml.NewDecoder(bytes.NewReader(raw))
	dec.Strict = false
	// Note: custom entities in XTB files if any
	dec.Entity = xml.HTMLEntity

	if xr.contentHandler != nil {
		xr.contentHandler.StartDocument()
	}
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			throwSAXException(err.Error())
		}
		if xr.contentHandler == nil {
			continue
		}
		switch t := tok.(type) {
		case xml.StartElement:
			xr.contentHandler.StartElement(t.Name.Space, t.Name.Local, t.Name.Local, &jSaxAttrs{attrs: t.Attr})
		case xml.EndElement:
			xr.contentHandler.EndElement(t.Name.Space, t.Name.Local, t.Name.Local)
		case xml.CharData:
			u16 := utf16.Encode([]rune(string(t)))
			if len(u16) > 0 {
				xr.contentHandler.Characters(u16, 0, int32(len(u16)))
			}
		case xml.ProcInst:
			xr.contentHandler.ProcessingInstruction(t.Target, string(t.Inst))
		}
	}
	if xr.contentHandler != nil {
		xr.contentHandler.EndDocument()
	}
}

// --- Regex, MessageFormat, ResourceBundle, StringTokenizer, Logging, JMX, System ---

func wrapPatternCompile(regex string, flags int) *Ext_java_util_regex_Pattern {
	p := &Ext_java_util_regex_Pattern{Core_: javaregex.Pattern_CompileFlags(regex, flags)}
	p.BindExtSelf_(p)
	return p
}

func javaMessageFormat(pattern string, args []any) string {
	var sb strings.Builder
	inQuote := false
	n := len(pattern)
	for i := 0; i < n; i++ {
		c := pattern[i]
		if c == '\'' {
			if i+1 < n && pattern[i+1] == '\'' {
				sb.WriteByte('\'')
				i++
				continue
			}
			inQuote = !inQuote
			continue
		}
		if !inQuote && c == '{' {
			end := strings.IndexByte(pattern[i+1:], '}')
			if end > 0 {
				spec := pattern[i+1 : i+1+end]
				if comma := strings.IndexByte(spec, ','); comma >= 0 {
					spec = spec[:comma]
				}
				if idx, err := strconv.Atoi(strings.TrimSpace(spec)); err == nil && idx >= 0 {
					if idx < len(args) {
						sb.WriteString(StrVal(args[idx]))
					} else {
						sb.WriteByte('{')
						sb.WriteString(pattern[i+1 : i+1+end])
						sb.WriteByte('}')
					}
					i += end + 1
					continue
				}
			}
		}
		sb.WriteByte(c)
	}
	return sb.String()
}

func wrapResourceBundle(baseName string) Ext_java_util_ResourceBundle {
	rb := shimsres.ResourceBundle_GetBundle(baseName)
	obj := &Ext_java_util_ResourceBundleBase{Core_: rb}
	obj.BindExtSelf_(obj)
	return obj
}

type jStringTokenizer struct {
	str    string
	delims string
	pos    int
}

func newStringTokenizer(str, delims string) *jStringTokenizer {
	return &jStringTokenizer{str: str, delims: delims, pos: 0}
}

func (st *jStringTokenizer) isDelim(b byte) bool {
	return strings.IndexByte(st.delims, b) >= 0
}

func (st *jStringTokenizer) skipDelims(p int) int {
	for p < len(st.str) && st.isDelim(st.str[p]) {
		p++
	}
	return p
}

func (st *jStringTokenizer) hasMoreTokens() bool {
	st.pos = st.skipDelims(st.pos)
	return st.pos < len(st.str)
}

func (st *jStringTokenizer) nextToken() string {
	st.pos = st.skipDelims(st.pos)
	if st.pos >= len(st.str) {
		panic(javalang.NewNoSuchElementException("StringTokenizer"))
	}
	start := st.pos
	for st.pos < len(st.str) && !st.isDelim(st.str[st.pos]) {
		st.pos++
	}
	return st.str[start:st.pos]
}

var loggerCache sync.Map // string -> *Ext_java_util_logging_Logger

func getJavaLogger(name string) *Ext_java_util_logging_Logger {
	if v, ok := loggerCache.Load(name); ok {
		return v.(*Ext_java_util_logging_Logger)
	}
	l := &Ext_java_util_logging_Logger{Core_: Ext_java_util_logging_Level_F_OFF}
	l.BindExtSelf_(l)
	actual, _ := loggerCache.LoadOrStore(name, l)
	return actual.(*Ext_java_util_logging_Logger)
}

func loggerIsLoggable(l *Ext_java_util_logging_Logger, lvl *Ext_java_util_logging_Level) bool {
	if l == nil || lvl == nil {
		return false
	}
	cur, _ := l.Core_.(*Ext_java_util_logging_Level)
	if cur == nil || cur.GetName() == "OFF" || lvl.GetName() == "OFF" {
		return false
	}
	return false
}

var (
	mainJavaThreadOnce sync.Once
	mainJavaThread     *Ext_java_lang_Thread
	goroutineThreads   sync.Map
)

func currentJavaThread() *Ext_java_lang_Thread {
	gid := currentGoroutineID()
	if gid == 1 {
		mainJavaThreadOnce.Do(func() {
			t := &Ext_java_lang_Thread{}
			t.BindExtSelf_(t)
			mainJavaThread = t
		})
		return mainJavaThread
	}
	if v, ok := goroutineThreads.Load(gid); ok {
		return v.(*Ext_java_lang_Thread)
	}
	t := &Ext_java_lang_Thread{}
	t.BindExtSelf_(t)
	actual, _ := goroutineThreads.LoadOrStore(gid, t)
	return actual.(*Ext_java_lang_Thread)
}

type jThreadLocal struct {
	mu   sync.Mutex
	vals map[uint64]any
}

func (tl *jThreadLocal) get() any {
	gid := currentGoroutineID()
	tl.mu.Lock()
	defer tl.mu.Unlock()
	if tl.vals == nil {
		return nil
	}
	return tl.vals[gid]
}

func (tl *jThreadLocal) set(v any) {
	gid := currentGoroutineID()
	tl.mu.Lock()
	defer tl.mu.Unlock()
	if tl.vals == nil {
		tl.vals = make(map[uint64]any)
	}
	tl.vals[gid] = ToAny(v)
}

func (tl *jThreadLocal) remove() {
	gid := currentGoroutineID()
	tl.mu.Lock()
	defer tl.mu.Unlock()
	delete(tl.vals, gid)
}

type jMemoryUsage struct {
	init, used, committed, max int64
}

func runtimeFreeMemory() int64 {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	if ms.HeapIdle > ms.HeapReleased {
		return int64(ms.HeapIdle - ms.HeapReleased)
	}
	return 0
}

func runtimeTotalMemory() int64 {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return int64(ms.HeapSys)
}

func runtimeMaxMemory() int64 {
	return 4 * 1024 * 1024 * 1024
}

type jCompilationMXBean struct{}

func (*jCompilationMXBean) Implements_j2gort_Ext_java_lang_management_CompilationMXBean()     {}
func (*jCompilationMXBean) Implements_j2gort_Ext_java_lang_management_PlatformManagedObject() {}
func (*jCompilationMXBean) GetName() string                                                   { return "GoCompiler" }
func (*jCompilationMXBean) GetTotalCompilationTime() int64                                    { return 0 }

type jGarbageCollectorMXBean struct{}

func (*jGarbageCollectorMXBean) Implements_j2gort_Ext_java_lang_management_GarbageCollectorMXBean() {}
func (*jGarbageCollectorMXBean) Implements_j2gort_Ext_java_lang_management_MemoryManagerMXBean()    {}
func (*jGarbageCollectorMXBean) Implements_j2gort_Ext_java_lang_management_PlatformManagedObject()  {}
func (*jGarbageCollectorMXBean) GetName() string                                                    { return "GoGC" }
func (*jGarbageCollectorMXBean) GetCollectionCount() int64 {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return int64(ms.NumGC)
}
func (*jGarbageCollectorMXBean) GetCollectionTime() int64 {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return int64(ms.PauseTotalNs / 1_000_000)
}

type jMemoryMXBean struct{}

func (*jMemoryMXBean) Implements_j2gort_Ext_java_lang_management_MemoryMXBean()          {}
func (*jMemoryMXBean) Implements_j2gort_Ext_java_lang_management_PlatformManagedObject() {}
func (*jMemoryMXBean) GetHeapMemoryUsage() *Ext_java_lang_management_MemoryUsage {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return &Ext_java_lang_management_MemoryUsage{
		Core_: &jMemoryUsage{init: 0, used: int64(ms.HeapAlloc), committed: int64(ms.HeapSys), max: runtimeMaxMemory()},
	}
}
func (*jMemoryMXBean) GetNonHeapMemoryUsage() *Ext_java_lang_management_MemoryUsage {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return &Ext_java_lang_management_MemoryUsage{
		Core_: &jMemoryUsage{init: 0, used: int64(ms.StackInuse), committed: int64(ms.StackSys), max: -1},
	}
}

type jMemoryPoolMXBean struct{}

func (*jMemoryPoolMXBean) Implements_j2gort_Ext_java_lang_management_MemoryPoolMXBean()      {}
func (*jMemoryPoolMXBean) Implements_j2gort_Ext_java_lang_management_PlatformManagedObject() {}
func (*jMemoryPoolMXBean) GetName() string                                                   { return "GoHeap" }
func (*jMemoryPoolMXBean) GetUsage() *Ext_java_lang_management_MemoryUsage {
	return (&jMemoryMXBean{}).GetHeapMemoryUsage()
}
func (*jMemoryPoolMXBean) GetPeakUsage() *Ext_java_lang_management_MemoryUsage {
	return (&jMemoryMXBean{}).GetHeapMemoryUsage()
}

func javaSystemGetProperty(key, def string) string {
	switch key {
	case "line.separator":
		return "\n"
	case "file.separator":
		return "/"
	case "path.separator":
		return ":"
	case "user.dir":
		wd, err := os.Getwd()
		if err != nil {
			return "."
		}
		return wd
	case "user.home":
		h, err := os.UserHomeDir()
		if err == nil {
			return h
		}
	case "java.io.tmpdir":
		return os.TempDir()
	}
	if val, ok := os.LookupEnv("JSGO_PROP_" + strings.ReplaceAll(key, ".", "_")); ok {
		return val
	}
	return def
}

func javaSystemGetenv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return NullString
	}
	return val
}

func ThrowablePrintStackTrace(t javalang.Throwable, dest any) {
	CheckNotNull(t)
	msg := t.Error()
	if dest == nil || javalang.IsNil(dest) {
		javaio.System_Err.Println(msg)
		return
	}
	switch d := dest.(type) {
	case *Ext_java_io_PrintStream:
		unwrapPrintStream(d).Println(msg)
	case *Ext_java_io_PrintWriter:
		extWriterWriteString(d.Core_, msg+"\n")
		extWriterFlush(d.Core_)
	default:
		javaio.System_Err.Println(msg)
	}
}

func (x *extInputStreamAdapter) IsJavaNil_() bool   { return x == nil }
func (x *extOutputStreamAdapter) IsJavaNil_() bool  { return x == nil }
func (x *extWriterAdapter) IsJavaNil_() bool        { return x == nil }
func (x *gzipDeflaterBridge) IsJavaNil_() bool      { return x == nil }
func (x *jNioPath) IsJavaNil_() bool                { return x == nil }
func (x *jPathMatcher) IsJavaNil_() bool            { return x == nil }
func (x *jBasicFileAttrs) IsJavaNil_() bool         { return x == nil }
func (x *jFileVisitorAdapter) IsJavaNil_() bool     { return x == nil }
func (x *jGuavaFunc) IsJavaNil_() bool              { return x == nil }
func (x *jGuavaPred) IsJavaNil_() bool              { return x == nil }
func (x *jMemoSupplier) IsJavaNil_() bool           { return x == nil }
func (x *jToStringHelperData) IsJavaNil_() bool     { return x == nil }
func (x *jHashFunction) IsJavaNil_() bool           { return x == nil }
func (x *jHasher) IsJavaNil_() bool                 { return x == nil }
func (x *jCharSourceSpec) IsJavaNil_() bool         { return x == nil }
func (x *jCharSinkSpec) IsJavaNil_() bool           { return x == nil }
func (x *jGsonAdapterBridge) IsJavaNil_() bool      { return x == nil }
func (x *jGsonBuilderCore) IsJavaNil_() bool        { return x == nil }
func (x *jSaxAttrs) IsJavaNil_() bool               { return x == nil }
func (x *jSaxXMLReader) IsJavaNil_() bool           { return x == nil }
func (x *jStringTokenizer) IsJavaNil_() bool        { return x == nil }
func (x *jThreadLocal) IsJavaNil_() bool            { return x == nil }
func (x *jMemoryUsage) IsJavaNil_() bool            { return x == nil }
func (x *jCompilationMXBean) IsJavaNil_() bool      { return x == nil }
func (x *jGarbageCollectorMXBean) IsJavaNil_() bool { return x == nil }
func (x *jMemoryMXBean) IsJavaNil_() bool           { return x == nil }
func (x *jMemoryPoolMXBean) IsJavaNil_() bool       { return x == nil }
