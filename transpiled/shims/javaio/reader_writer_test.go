package javaio

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"testing/iotest"
)

func TestStringReaderAndWriter(t *testing.T) {
	sr := NewStringReader("line1\nline2\r\nline3")
	br := NewBufferedReader(sr)

	l1, err := br.ReadLine()
	if err != nil || l1 != "line1" {
		t.Errorf("ReadLine 1 failed: got %q, err=%v", l1, err)
	}

	l2, err := br.ReadLine()
	if err != nil || l2 != "line2" {
		t.Errorf("ReadLine 2 failed: got %q, err=%v", l2, err)
	}

	l3, err := br.ReadLine()
	if err != nil || l3 != "line3" {
		t.Errorf("ReadLine 3 failed: got %q, err=%v", l3, err)
	}

	_, err = br.ReadLine()
	if err != io.EOF {
		t.Errorf("expected EOF, got %v", err)
	}

	sw := NewStringWriter()
	sw.Write("hello ")
	sw.AppendCharSequence("world")
	sw.AppendCharacter('!')
	if sw.String() != "hello world!" {
		t.Errorf("StringWriter got %q, want \"hello world!\"", sw.String())
	}
}

func TestStreams(t *testing.T) {
	data := []byte{10, 20, 30, 40, 50}
	in := NewByteArrayInputStream(data)
	if in.Available() != 5 {
		t.Errorf("Available got %d, want 5", in.Available())
	}

	b0, err := in.ReadByte()
	if err != nil || b0 != 10 {
		t.Errorf("ReadByte got %d, want 10", b0)
	}

	buf := make([]byte, 2)
	n, err := in.Read(buf)
	if err != nil || n != 2 || buf[0] != 20 || buf[1] != 30 {
		t.Errorf("Read failed: n=%d, buf=%v", n, buf)
	}

	out := NewByteArrayOutputStream()
	out.Write([]byte("foo"))
	out.WriteByte('-')
	out.Write([]byte("bar"))
	if out.ToString() != "foo-bar" || out.Size() != 7 {
		t.Errorf("ByteArrayOutputStream got %q (size %d)", out.ToString(), out.Size())
	}
}

func TestFileAndPrintStream(t *testing.T) {
	dir := t.TempDir()
	f := NewFile(filepath.Join(dir, "sub", "test.txt"))

	parent := NewFile(f.GetParent())
	if !parent.Mkdirs() {
		t.Errorf("parent.Mkdirs() failed")
	}
	if !parent.Exists() || !parent.IsDirectory() {
		t.Errorf("parent should exist as directory")
	}

	err := os.WriteFile(f.GetPath(), []byte("sample content"), 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	if !f.Exists() || !f.IsFile() {
		t.Errorf("f should exist as a normal file")
	}
	if f.Length() != 14 {
		t.Errorf("f.Length() got %d, want 14", f.Length())
	}
	if f.GetName() != "test.txt" {
		t.Errorf("f.GetName() got %q, want \"test.txt\"", f.GetName())
	}

	var buf bytes.Buffer
	ps := NewPrintStream(&buf)
	ps.Print("count=")
	ps.Println(42)
	ps.Printf("hex=%x\n", 255)
	expected := "count=42\nhex=ff\n"
	if buf.String() != expected {
		t.Errorf("PrintStream output got %q, want %q", buf.String(), expected)
	}
}

func TestBufferedReaderBoundariesAndAllocations(t *testing.T) {
	// 1. \r\n split right across bufferedReaderBufSize (4096): \r at 4095, \n at 4096.
	prefix := strings.Repeat("x", bufferedReaderBufSize-1)
	src := prefix + "\r\nsecond\rthird\n\n" + strings.Repeat("y", bufferedReaderBufSize*2+17)
	br := NewBufferedReader(NewInputStreamReader(strings.NewReader(src)))

	wantLines := []string{
		prefix,
		"second",
		"third",
		"",
		strings.Repeat("y", bufferedReaderBufSize*2+17),
	}
	for i, want := range wantLines {
		got, err := br.ReadLine()
		if err != nil {
			t.Fatalf("line %d unexpected error: %v", i, err)
		}
		if got != want {
			t.Fatalf("line %d len=%d, want len=%d", i, len(got), len(want))
		}
	}
	if _, err := br.ReadLine(); err != io.EOF {
		t.Fatalf("expected EOF, got %v", err)
	}

	// 2. \r at index 4095 followed immediately by \n and EOF -> only 1 line!
	br2 := NewBufferedReader(NewStringReader(prefix + "\r\n"))
	l1, err := br2.ReadLine()
	if err != nil || l1 != prefix {
		t.Fatalf("expected prefix line, got len=%d err=%v", len(l1), err)
	}
	if _, err := br2.ReadLine(); err != io.EOF {
		t.Fatalf("expected EOF after boundary CRLF, got %v", err)
	}

	// 3. Interleaved ReadLine + Read across \r\n boundary.
	br3 := NewBufferedReader(NewStringReader(prefix + "\r\nhello"))
	if l, err := br3.ReadLine(); err != nil || l != prefix {
		t.Fatalf("ReadLine before Read failed: err=%v", err)
	}
	var out [5]uint16
	if n := br3.Read(out[:], 0, 5); n != 5 || out != [5]uint16{'h', 'e', 'l', 'l', 'o'} {
		t.Fatalf("Read after boundary CRLF got n=%d, out=%v", n, out)
	}

	// 4. UTF-8 4-byte emoji split across 8192-byte InputStreamReader chunk boundary + OneByteReader.
	utf8Boundary := strings.Repeat("a", isrByteBufSize-2) + "😀" + "Ω\nend"
	br4 := NewBufferedReader(NewInputStreamReader(iotest.OneByteReader(strings.NewReader(utf8Boundary))))
	uLine1, err := br4.ReadLine()
	if err != nil || uLine1 != strings.Repeat("a", isrByteBufSize-2)+"😀Ω" {
		t.Fatalf("UTF-8 boundary ReadLine mismatch: err=%v, len=%d", err, len(uLine1))
	}
	uLine2, err := br4.ReadLine()
	if err != nil || uLine2 != "end" {
		t.Fatalf("UTF-8 line 2 got %q, err=%v", uLine2, err)
	}

	// 5. Verify zero per-character allocations: reading 200 lines of 40 chars should allocate <= 205 objects total
	// (1 BufferedReader + 1 StringReader + 200 returned strings, vs ~10,000+ previously!).
	sampleJS := strings.Repeat("const x = 1234567890; // comment line!!\n", 200)
	allocs := testing.AllocsPerRun(20, func() {
		r := NewBufferedReader(NewStringReader(sampleJS))
		for {
			_, err := r.ReadLine()
			if err != nil {
				break
			}
		}
	})
	if allocs > 205 {
		t.Errorf("BufferedReader.ReadLine allocated %v objects for 200 lines, expected <= 205", allocs)
	}
}

func TestDrainStringFastPaths(t *testing.T) {
	// 1. Unread non-ASCII StringReader skips UTF-16 conversion completely (0 allocations on DrainString).
	nonASCII := strings.Repeat("const greeting = 'こんにちは😀';\n", 50)
	allocs := testing.AllocsPerRun(20, func() {
		sr := NewStringReader(nonASCII)
		s, ok := sr.DrainString()
		if !ok || s != nonASCII {
			t.Fatalf("unread StringReader.DrainString mismatch")
		}
		var buf [4]uint16
		if n := sr.Read(buf[:], 0, 4); n != -1 {
			t.Fatalf("expected EOF (-1) after DrainString, got %d", n)
		}
	})
	if allocs > 1 {
		t.Errorf("unread non-ASCII StringReader.DrainString allocated %v objects, expected <= 1", allocs)
	}

	// 2. Partially read non-ASCII StringReader still returns exact remaining substring.
	srPart := NewStringReader("αβ😀γ")
	if ch, ok := srPart.ReadChar(); !ok || ch != 'α' {
		t.Fatalf("ReadChar got %#x, ok=%v", ch, ok)
	}
	if rem, ok := srPart.DrainString(); !ok || rem != "β😀γ" {
		t.Fatalf("partially read StringReader.DrainString got %q (ok=%v), want %q", rem, ok, "β😀γ")
	}

	// 3. Unread BufferedReader wrapping InputStreamReader (UTF-8 and ISO-8859-1).
	brUTF8 := NewBufferedReader(NewInputStreamReader(strings.NewReader(nonASCII)))
	if got, ok := brUTF8.DrainString(); !ok || got != nonASCII {
		t.Fatalf("unread BufferedReader(InputStreamReader) DrainString mismatch: ok=%v", ok)
	}
	if _, err := brUTF8.ReadLine(); err != io.EOF {
		t.Fatalf("expected EOF after BufferedReader.DrainString, got %v", err)
	}

	isrLatin1 := NewInputStreamReaderCharset(bytes.NewReader([]byte{'c', 'a', 'f', 0xe9}), "ISO-8859-1")
	if got, ok := isrLatin1.DrainString(); !ok || got != "café" {
		t.Fatalf("Latin-1 DrainString got %q (ok=%v), want café", got, ok)
	}
}

func TestRemainingBytesHintAndReadAllBytesFastPath(t *testing.T) {
	payload := bytes.Repeat([]byte("0123456789abcdef"), 8192) // 128 KiB

	bais := NewByteArrayInputStream(payload)
	if _, err := bais.Skip(16); err != nil {
		t.Fatal(err)
	}
	if rem, ok := bais.RemainingBytesHint(); !ok || rem != int64(len(payload)-16) {
		t.Fatalf("ByteArrayInputStream.RemainingBytesHint got (%d, %v), want %d", rem, ok, len(payload)-16)
	}
	allocs := testing.AllocsPerRun(20, func() {
		s := NewByteArrayInputStream(payload)
		got, err := ReadAllBytes(s)
		if err != nil || len(got) != len(payload) {
			t.Fatalf("ReadAllBytes failed: len=%d err=%v", len(got), err)
		}
	})
	if allocs > 2 {
		t.Fatalf("expected single-allocation ReadAllBytes on ByteArrayInputStream, got %.1f allocs", allocs)
	}

	tmpPath := filepath.Join(t.TempDir(), "data.bin")
	if err := os.WriteFile(tmpPath, payload, 0644); err != nil {
		t.Fatal(err)
	}
	fis, err := NewFileInputStream(tmpPath)
	if err != nil {
		t.Fatal(err)
	}
	defer fis.Close()

	bis := NewBufferedInputStream(fis)
	var head [10]byte
	if n, err := bis.Read(head[:]); err != nil || n != 10 {
		t.Fatalf("bis.Read head failed: n=%d err=%v", n, err)
	}
	if rem, ok := bis.RemainingBytesHint(); !ok || rem != int64(len(payload)-10) {
		t.Fatalf("BufferedInputStream(FileInputStream).RemainingBytesHint got (%d, %v), want %d", rem, ok, len(payload)-10)
	}
	rest, err := ReadAllBytes(bis)
	if err != nil || !bytes.Equal(rest, payload[10:]) {
		t.Fatalf("ReadAllBytes(BufferedInputStream) mismatch: len=%d err=%v", len(rest), err)
	}
	_ = bis.Close()
}

type writeCountingSink struct {
	buf    bytes.Buffer
	writes int
	closes int
}

func (s *writeCountingSink) Write(p []byte) (int, error) {
	s.writes++
	return s.buf.Write(p)
}

func (s *writeCountingSink) Close() error {
	s.closes++
	return nil
}

func TestWave42CRawUnixFDFileOutputStreamAnd64KiBBufferedStreams(t *testing.T) {
	dir := t.TempDir()

	// 1. Create, truncate, append, WriteByte, small + >64 KiB WriteString, and NewFileOutputStreamFromOSFile.
	p1 := filepath.Join(dir, "chunk.js")
	fos1, err := NewFileOutputStream(p1)
	if err != nil {
		t.Fatalf("NewFileOutputStream failed: %v", err)
	}
	if _, err := fos1.Write([]byte("stale-initial-content-that-must-be-truncated")); err != nil {
		t.Fatal(err)
	}
	if err := fos1.Close(); err != nil {
		t.Fatal(err)
	}
	// Double Close() must be idempotent.
	if err := fos1.Close(); err != nil {
		t.Fatalf("double Close() failed: %v", err)
	}
	if _, err := fos1.Write([]byte("x")); err == nil {
		t.Fatal("expected error writing to closed FileOutputStream")
	}

	// Re-open via NewFileOutputStream -> must truncate existing content.
	fos2, err := NewFileOutputStream(p1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fos2.WriteString("const header = 1;\n"); err != nil {
		t.Fatal(err)
	}
	if err := fos2.WriteByte('/'); err != nil {
		t.Fatal(err)
	}
	if err := fos2.WriteByte('/'); err != nil {
		t.Fatal(err)
	}
	largeBody := strings.Repeat("export const line = '0123456789abcdef';\n", 2500) // 100 KiB (> 64 KiB scratch chunk)
	if n, err := fos2.WriteString(largeBody); err != nil || n != len(largeBody) {
		t.Fatalf("large WriteString failed: n=%d err=%v", n, err)
	}
	if err := fos2.Close(); err != nil {
		t.Fatal(err)
	}

	// Append via NewFileOutputStreamAppend(..., true).
	fosApp, err := NewFileOutputStreamAppend(p1, true)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fosApp.Write([]byte("// trailer\n")); err != nil {
		t.Fatal(err)
	}
	if err := fosApp.Close(); err != nil {
		t.Fatal(err)
	}

	wantFull := "const header = 1;\n//" + largeBody + "// trailer\n"
	fis, err := NewFileInputStream(p1)
	if err != nil {
		t.Fatal(err)
	}
	bis := NewBufferedInputStream(fis)
	gotBytes, err := ReadAllBytes(bis)
	if err != nil {
		t.Fatal(err)
	}
	if err := bis.Close(); err != nil {
		t.Fatal(err)
	}
	if string(gotBytes) != wantFull {
		t.Fatalf("round-trip mismatch: got len=%d, want len=%d", len(gotBytes), len(wantFull))
	}

	// NewFileOutputStreamFromOSFile compatibility.
	osFilePath := filepath.Join(dir, "from_os.txt")
	osf, err := os.Create(osFilePath)
	if err != nil {
		t.Fatal(err)
	}
	fosOS := NewFileOutputStreamFromOSFile(osf)
	if _, err := fosOS.WriteString("os-file-ok"); err != nil {
		t.Fatal(err)
	}
	if err := fosOS.Close(); err != nil {
		t.Fatal(err)
	}
	if b, err := os.ReadFile(osFilePath); err != nil || string(b) != "os-file-ok" {
		t.Fatalf("NewFileOutputStreamFromOSFile mismatch: %q err=%v", string(b), err)
	}

	// Error handling on bad path (missing parent directory -> FileNotFoundException).
	if _, err := NewFileOutputStream(filepath.Join(dir, "no_such_dir", "out.js")); err == nil {
		t.Fatal("expected FileNotFoundException on missing parent directory")
	} else if _, ok := err.(*FileNotFoundException); !ok {
		t.Fatalf("expected *FileNotFoundException, got %T (%v)", err, err)
	}

	// 2. 64 KiB pooled BufferedOutputStream: 50 KiB of small writes must coalesce into ONE flush on Close()!
	sink := &writeCountingSink{}
	bos := NewBufferedOutputStream(sink)
	chunk := bytes.Repeat([]byte("abcdefghij"), 100) // 1,000 B
	for i := 0; i < 50; i++ {
		if _, err := bos.Write(chunk); err != nil {
			t.Fatal(err)
		}
	}
	if sink.writes != 0 {
		t.Fatalf("expected 0 flushes for 50 KB inside 64 KiB BufferedOutputStream, got %d", sink.writes)
	}
	if err := bos.Close(); err != nil {
		t.Fatal(err)
	}
	if sink.writes != 1 || sink.closes != 1 || sink.buf.Len() != 50000 {
		t.Fatalf("after Close(): writes=%d closes=%d len=%d", sink.writes, sink.closes, sink.buf.Len())
	}

	// 3. Concurrent -race stress: parallel file creates/writes/reads + concurrent Close() on shared stream.
	sharedPath := filepath.Join(dir, "shared_close.js")
	sharedOut, err := NewFileOutputStream(sharedPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sharedOut.WriteString("race-safe-close"); err != nil {
		t.Fatal(err)
	}
	var closeWG sync.WaitGroup
	for g := 0; g < 16; g++ {
		closeWG.Add(1)
		go func() {
			defer closeWG.Done()
			_ = sharedOut.Close()
		}()
	}
	closeWG.Wait()

	var wg sync.WaitGroup
	for g := 0; g < 16; g++ {
		g := g
		wg.Add(1)
		go func() {
			defer wg.Done()
			p := filepath.Join(dir, "worker_"+strings.Repeat("x", g%5)+"_"+string(rune('a'+g))+".bin")
			out, err := NewFileOutputStream(p)
			if err != nil {
				t.Errorf("worker %d open failed: %v", g, err)
				return
			}
			bw := NewBufferedOutputStream(out)
			payload := bytes.Repeat([]byte{byte('A' + g)}, 80*1024) // > 64 KiB
			if _, err := bw.Write(payload); err != nil {
				t.Errorf("worker %d write failed: %v", g, err)
				return
			}
			if err := bw.Close(); err != nil {
				t.Errorf("worker %d close failed: %v", g, err)
				return
			}
			in, err := NewFileInputStream(p)
			if err != nil {
				t.Errorf("worker %d reopen failed: %v", g, err)
				return
			}
			bi := NewBufferedInputStream(in)
			got, err := ReadAllBytes(bi)
			_ = bi.Close()
			if err != nil || !bytes.Equal(got, payload) {
				t.Errorf("worker %d verify mismatch: len=%d err=%v", g, len(got), err)
			}
		}()
	}
	wg.Wait()
}
