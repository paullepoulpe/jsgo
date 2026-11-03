package jutil

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"sync/atomic"
	"testing"
)

type syntheticArchiveEntry struct {
	name    string
	deleted bool // locrec_offset == 0 tombstone
}

func buildSyntheticLocaleArchive(namehashSize uint32, entries []syntheticArchiveEntry) []byte {
	const hdrSize = glibcLocArchiveHeaderBytes
	tableBytes := int(namehashSize) * glibcLocArchiveEntryBytes
	strStart := hdrSize + tableBytes

	var strTable []byte
	type placedEntry struct {
		hashval   uint32
		nameOff   uint32
		locrecOff uint32
	}
	buckets := make([]placedEntry, namehashSize)

	for i, e := range entries {
		nameOff := uint32(strStart + len(strTable))
		strTable = append(strTable, e.name...)
		strTable = append(strTable, 0)

		hval := glibcArchiveHash(e.name)
		idx := hval % namehashSize
		incr := uint32(1)
		if namehashSize > 2 {
			incr = 1 + hval%(namehashSize-2)
		}
		for buckets[idx].nameOff != 0 {
			idx = (idx + incr) % namehashSize
		}
		locrecOff := uint32(0x10000 + i*64)
		if e.deleted {
			locrecOff = 0
		}
		buckets[idx] = placedEntry{
			hashval:   hval,
			nameOff:   nameOff,
			locrecOff: locrecOff,
		}
	}

	buf := make([]byte, strStart+len(strTable))
	binary.LittleEndian.PutUint32(buf[0:4], glibcLocArchiveMagic)
	binary.LittleEndian.PutUint32(buf[8:12], hdrSize)
	binary.LittleEndian.PutUint32(buf[12:16], uint32(len(entries)))
	binary.LittleEndian.PutUint32(buf[16:20], namehashSize)
	binary.LittleEndian.PutUint32(buf[20:24], uint32(strStart))
	binary.LittleEndian.PutUint32(buf[24:28], uint32(len(strTable)))
	binary.LittleEndian.PutUint32(buf[28:32], uint32(len(strTable)))

	for idx, b := range buckets {
		off := hdrSize + idx*glibcLocArchiveEntryBytes
		binary.LittleEndian.PutUint32(buf[off:off+4], b.hashval)
		binary.LittleEndian.PutUint32(buf[off+4:off+8], b.nameOff)
		binary.LittleEndian.PutUint32(buf[off+8:off+12], b.locrecOff)
	}
	copy(buf[strStart:], strTable)
	return buf
}

type countingVirtualReaderAt struct {
	prefix      []byte
	virtualSize int64
	bytesRead   atomic.Int64
	readCalls   atomic.Int64
}

func (r *countingVirtualReaderAt) ReadAt(p []byte, off int64) (int, error) {
	r.readCalls.Add(1)
	if off < 0 || off >= r.virtualSize {
		return 0, io.EOF
	}
	avail := r.virtualSize - off
	n := len(p)
	var err error
	if int64(n) > avail {
		n = int(avail)
		err = io.EOF
	}
	for i := 0; i < n; i++ {
		pos := off + int64(i)
		if pos < int64(len(r.prefix)) {
			p[i] = r.prefix[pos]
		} else {
			p[i] = 0
		}
	}
	r.bytesRead.Add(int64(n))
	return n, err
}

func TestLocaleArchiveContainsGlibcFormatAndSuffixSafety(t *testing.T) {
	// Use a compact 7-slot hash table with 6 entries (including a deleted tombstone
	// and "nde.utf8" suffix trap) to force multi-step double-hashing probe chains.
	raw := buildSyntheticLocaleArchive(7, []syntheticArchiveEntry{
		{name: "en_US.utf8"},
		{name: "nde.utf8"},
		{name: "es_ES.utf8", deleted: true},
		{name: "fr_FR.iso88591"},
		{name: "ja_JP.utf8"},
		{name: "zh_CN.utf8"},
	})

	const virtual200MiB = 200 << 20
	reader := &countingVirtualReaderAt{
		prefix:      raw,
		virtualSize: virtual200MiB,
	}

	wantInstalled := []string{"en_US.utf8", "nde.utf8", "fr_FR.iso88591", "ja_JP.utf8", "zh_CN.utf8"}
	for _, loc := range wantInstalled {
		reader.bytesRead.Store(0)
		if !localeArchiveContains(reader, virtual200MiB, loc) {
			t.Fatalf("expected %q to be found in synthetic locale-archive", loc)
		}
		if gotBytes := reader.bytesRead.Load(); gotBytes > 512 {
			t.Fatalf("lookup of %q read %d bytes (expected tiny header+probe ReadAt, not full archive)", loc, gotBytes)
		}
	}

	// Suffix, prefix, missing, and deleted tombstone entries MUST return false.
	wantMissing := []string{
		"de.utf8",    // suffix of installed "nde.utf8" — must NOT match!
		"e.utf8",     // shorter suffix of "nde.utf8"
		"nde.utf",    // prefix without NUL match
		"es_ES.utf8", // deleted tombstone (locrec_offset == 0)
		"pt_BR.utf8", // absent locale probing across crowded table
		"",
	}
	for _, loc := range wantMissing {
		if localeArchiveContains(reader, virtual200MiB, loc) {
			t.Fatalf("expected %q NOT to match in synthetic locale-archive", loc)
		}
	}

	// Verify exact-EOF boundary when string NUL is the very last byte of bytes.Reader.
	exactReader := bytes.NewReader(raw)
	if !localeArchiveContains(exactReader, int64(len(raw)), "zh_CN.utf8") {
		t.Fatalf("expected final entry zh_CN.utf8 to match at exact archive EOF")
	}
	if localeArchiveContains(exactReader, int64(len(raw)), "de.utf8") {
		t.Fatalf("expected suffix de.utf8 not to match in exactReader")
	}

	// Corrupt / adversarial headers must return false cleanly without panicking or spinning.
	if localeArchiveContains(bytes.NewReader([]byte("short")), 5, "en_US.utf8") {
		t.Fatal("truncated header should return false")
	}
	badMagic := append([]byte(nil), raw...)
	binary.LittleEndian.PutUint32(badMagic[0:4], 0xdeadbeef)
	if localeArchiveContains(bytes.NewReader(badMagic), int64(len(badMagic)), "en_US.utf8") {
		t.Fatal("bad magic should return false")
	}
	badOffset := append([]byte(nil), raw...)
	binary.LittleEndian.PutUint32(badOffset[8:12], uint32(len(badOffset)+100))
	if localeArchiveContains(bytes.NewReader(badOffset), int64(len(badOffset)), "en_US.utf8") {
		t.Fatal("out-of-bounds namehash_offset should return false")
	}
	fullTable := buildSyntheticLocaleArchive(3, []syntheticArchiveEntry{
		{name: "a.utf8"},
		{name: "b.utf8"},
		{name: "c.utf8"},
	})
	if localeArchiveContains(bytes.NewReader(fullTable), int64(len(fullTable)), "missing.utf8") {
		t.Fatal("full table probe should terminate and return false")
	}
}

func TestSystemNativeCharsetCachingAndEnvInvalidation(t *testing.T) {
	t.Setenv("LC_ALL", "C.UTF-8")
	t.Setenv("LC_CTYPE", "")
	t.Setenv("LANG", "")
	if got := System_NativeCharset(); got != StandardCharsets_UTF_8 {
		t.Fatalf("expected UTF-8 for C.UTF-8, got %v", got)
	}

	allocs := testing.AllocsPerRun(100, func() {
		if got := System_NativeCharset(); got != StandardCharsets_UTF_8 {
			t.Fatalf("unexpected charset: %v", got)
		}
	})
	if allocs != 0 {
		t.Fatalf("expected 0 allocs on cached System_NativeCharset(), got %v", allocs)
	}

	// Changing env vars must immediately invalidate the cached entry.
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_CTYPE", "C.ISO-8859-1")
	if got := System_NativeCharset(); got != standardCharsets_ISO_8859_1 {
		t.Fatalf("expected ISO-8859-1 for C.ISO-8859-1, got %v", got)
	}

	t.Setenv("LC_CTYPE", "")
	t.Setenv("LANG", "POSIX")
	if got := System_NativeCharset(); got != StandardCharsets_US_ASCII {
		t.Fatalf("expected US-ASCII for POSIX, got %v", got)
	}
}

func TestStringReaderReadSignedOverflowBounds(t *testing.T) {
	cbuf := make([]uint16, 8)
	cases := [][2]int{
		{1, math.MaxInt},
		{4, math.MaxInt - 2},
		{8, 1},
		{9, 0},
		{-1, 1},
		{0, -1},
	}
	for _, tc := range cases {
		off, length := tc[0], tc[1]
		func() {
			defer func() {
				rec := recover()
				if rec == nil {
					t.Fatalf("StringReader.Read(off=%d, len=%d) did not panic", off, length)
				}
				if _, ok := rec.(IIndexOutOfBoundsException); !ok {
					t.Fatalf("StringReader.Read(off=%d, len=%d) panicked with %T (%v), want IIndexOutOfBoundsException", off, length, rec, rec)
				}
			}()
			r := NewStringReader("hello")
			_ = r.Read(cbuf, off, length)
		}()
	}

	// Exact boundary reads still succeed.
	r := NewStringReader("hello")
	if n := r.Read(cbuf, 8, 0); n != 0 {
		t.Fatalf("expected 0 at exact end with length 0, got %d", n)
	}
	if n := r.Read(cbuf, 3, 5); n != 5 {
		t.Fatalf("expected 5 chars read into cbuf[3:8], got %d", n)
	}
}
