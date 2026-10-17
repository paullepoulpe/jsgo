package jutil

import (
	"bufio"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

// TestDeflaterVectors replays testdata/deflater_vectors.txt, the vectors of
// tools/oracle/java/JutilDeflaterVectors.java (whose header documents the
// case specs), against Deflater, DeflaterOutputStream and
// GZIPOutputStream. jutil compresses with compress/flate, so the compressed
// bytes are not compared with the JDK's: each case's output must decompress
// to exactly its input (which is checked against the JDK side's hash), the
// GZIP header must be the JDK's, a sync flush must make everything written
// so far decodable, and the exceptions must be the JDK's.

// deflaterVectorRng is the generator's SplitMix64.
type deflaterVectorRng struct{ state uint64 }

func (r *deflaterVectorRng) next() uint64 {
	r.state += 0x9e3779b97f4a7c15
	z := r.state
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

func (r *deflaterVectorRng) below(n int) int { return int(r.next() % uint64(n)) }

func deflaterVectorInput(t *testing.T, spec string) []byte {
	p := strings.Split(spec, ":")
	shape := p[0]
	size, _ := strconv.Atoi(p[1])
	seed, _ := strconv.ParseInt(p[2], 10, 64)
	r := &deflaterVectorRng{uint64(seed)}
	b := make([]byte, size)
	switch shape {
	case "rand":
		for i := range b {
			b[i] = byte(r.next())
		}
	case "text":
		alphabet := []byte("abcdefgh \n")
		for i := range b {
			b[i] = alphabet[r.below(len(alphabet))]
		}
	case "words", "mix":
		vr := &deflaterVectorRng{12345}
		vocab := make([][]byte, 300)
		for w := range vocab {
			vocab[w] = make([]byte, 1+vr.below(10))
			for k := range vocab[w] {
				vocab[w][k] = byte('a' + vr.below(26))
			}
		}
		seps := []byte(" , . ;\n(){}=")
		i := 0
		for i < size {
			if shape == "mix" && r.below(8) == 0 {
				n := 1 + r.below(40000)
				for k := 0; k < n && i < size; k++ {
					b[i] = byte(r.next())
					i++
				}
				continue
			}
			w := r.below(1 + r.below(len(vocab)))
			for k := 0; k < len(vocab[w]) && i < size; k++ {
				b[i] = vocab[w][k]
				i++
			}
			if i < size {
				b[i] = seps[r.below(len(seps))]
				i++
			}
		}
	case "runs":
		i := 0
		for i < size {
			v := byte(r.below(4))
			n := 1 + r.below(300)
			for k := 0; k < n && i < size; k++ {
				b[i] = v
				i++
			}
		}
	case "zeros":
	default:
		t.Fatalf("unknown shape %q", shape)
	}
	return b
}

func deflaterVectorChunks(t *testing.T, spec string, size int) [][2]int {
	var out [][2]int
	p := strings.Split(spec, ":")
	switch p[0] {
	case "all":
		out = append(out, [2]int{0, size})
	case "fixed":
		n, _ := strconv.Atoi(p[1])
		for i := 0; i < size; i += n {
			out = append(out, [2]int{i, min(n, size-i)})
		}
	case "rand":
		seed, _ := strconv.ParseInt(p[1], 10, 64)
		r := &deflaterVectorRng{uint64(seed)}
		max_, _ := strconv.Atoi(p[2])
		for i := 0; i < size; {
			n := min(1+r.below(max_), size-i)
			out = append(out, [2]int{i, n})
			i += n
		}
	default:
		t.Fatalf("unknown chunks %q", spec)
	}
	if len(out) == 0 {
		out = append(out, [2]int{0, 0})
	}
	return out
}

// deflaterVectorRun runs a case and returns its input and output.
func deflaterVectorRun(t *testing.T, spec string) (data, out []byte) {
	m := map[string]string{}
	for _, kv := range strings.Split(spec, " ") {
		k, v, _ := strings.Cut(kv, "=")
		m[k] = v
	}
	get := func(k, def string) string {
		if v, ok := m[k]; ok {
			return v
		}
		return def
	}
	atoi := func(s string) int {
		n, err := strconv.Atoi(s)
		if err != nil {
			t.Fatalf("%s: %v", spec, err)
		}
		return n
	}
	data = deflaterVectorInput(t, m["input"])
	chunks := deflaterVectorChunks(t, get("chunks", "all"), len(data))
	outSize := atoi(get("out", "512"))
	var bos bytes.Buffer
	switch m["kind"] {
	case "deflater":
		d := NewDeflaterWithLevelAndNowrap(atoi(get("level", "-1")), get("nowrap", "0") == "1")
		if strategy := atoi(get("strategy", "0")); strategy != 0 {
			d.SetStrategy(strategy)
		}
		midAt, midLevel, midStrategy := -1, 0, 0
		if mid, ok := m["mid"]; ok {
			mp := strings.Split(mid, ":")
			midAt, midLevel, midStrategy = atoi(mp[0]), atoi(mp[1]), atoi(mp[2])
		}
		buf := make([]byte, outSize)
		fed := 0
		for _, c := range chunks {
			d.SetInput(data[c[0] : c[0]+c[1]])
			for !d.NeedsInput() {
				n := d.Deflate(buf)
				bos.Write(buf[:n])
			}
			fed += c[1]
			if midAt >= 0 && fed >= midAt {
				d.SetLevel(midLevel)
				d.SetStrategy(midStrategy)
				midAt = -1
			}
		}
		d.Finish()
		for !d.Finished() {
			n := d.Deflate(buf)
			bos.Write(buf[:n])
		}
		d.End()
	case "gzip", "gzipsync":
		sync := m["kind"] == "gzipsync"
		g := NewGZIPOutputStreamWithSizeAndSyncFlush(&bos, outSize, sync)
		if prelevel, ok := m["prelevel"]; ok {
			g.Def.SetLevel(atoi(prelevel))
		}
		for _, c := range chunks {
			g.Write(data[c[0] : c[0]+c[1]])
			if sync {
				g.Flush()
				// Everything written so far must be decodable now.
				got, err := io.ReadAll(flate.NewReader(bytes.NewReader(bos.Bytes()[10:])))
				if !errors.Is(err, io.ErrUnexpectedEOF) || !bytes.Equal(got, data[:c[0]+c[1]]) {
					t.Errorf("%s: after flush at %d, decoded %d bytes (%v)", spec, c[0]+c[1], len(got), err)
				}
			}
		}
		g.Close()
	default:
		t.Fatalf("unknown kind in %q", spec)
	}
	return data, bos.Bytes()
}

// deflaterVectorDecompress decompresses a case's output.
func deflaterVectorDecompress(spec string, out []byte) ([]byte, error) {
	var r io.Reader
	var err error
	switch {
	case strings.Contains(spec, "kind=gzip"):
		var z *gzip.Reader
		z, err = gzip.NewReader(bytes.NewReader(out))
		if err == nil {
			z.Multistream(false)
			r = z
		}
	case strings.Contains(spec, "nowrap=1"):
		r = flate.NewReader(bytes.NewReader(out))
	default:
		r, err = zlib.NewReader(bytes.NewReader(out))
	}
	if err != nil {
		return nil, err
	}
	return io.ReadAll(r)
}

func deflaterVectorExc(body func()) (s string) {
	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(Throwable)
			if !ok {
				panic(r)
			}
			s = e.throwable().className + ": " + ToJavaString(GetMessageOrNull(e))
		}
	}()
	body()
	return "none"
}

func TestDeflaterVectors(t *testing.T) {
	f, err := os.Open("testdata/deflater_vectors.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	excs := map[string]func(){
		"newlevel10":   func() { NewDeflaterWithLevelAndNowrap(10, true) },
		"setlevel-2":   func() { NewDeflaterWithLevelAndNowrap(1, true).SetLevel(-2) },
		"setstrategy3": func() { NewDeflaterWithLevelAndNowrap(1, true).SetStrategy(3) },
		"flush1": func() {
			NewDeflaterWithLevelAndNowrap(1, true).DeflateWithFlush(make([]byte, 10), 1)
		},
		"ended": func() {
			d := NewDeflaterWithLevelAndNowrap(1, true)
			d.End()
			d.SetInput(make([]byte, 3))
			d.SetLevel(5)
			d.NeedsInput()
			d.Finished()
			d.Deflate(make([]byte, 10))
		},
		"gzipsize0": func() { NewGZIPOutputStreamWithSize(&bytes.Buffer{}, 0) },
		"writeafterfinish": func() {
			g := NewGZIPOutputStream(&bytes.Buffer{})
			g.Finish()
			g.WriteByte(1)
		},
	}
	sc := bufio.NewScanner(f)
	sc.Buffer(nil, 1<<20)
	cases := 0
	for sc.Scan() {
		spec, want, ok := strings.Cut(sc.Text(), " => ")
		if !ok {
			t.Fatalf("bad line %q", sc.Text())
		}
		cases++
		switch {
		case strings.HasPrefix(spec, "exc="):
			name := strings.TrimPrefix(spec, "exc=")
			body, ok := excs[name]
			if !ok {
				t.Fatalf("unknown exception case %q", name)
			}
			if got := "exc=" + deflaterVectorExc(body); got != want {
				t.Errorf("%s: got %s, want %s", name, got, want)
			}
		case spec == "header":
			var b bytes.Buffer
			NewGZIPOutputStream(&b)
			if got := "hex=" + hex.EncodeToString(b.Bytes()); got != want {
				t.Errorf("header: got %s, want %s", got, want)
			}
		default:
			data, out := deflaterVectorRun(t, spec)
			sum := sha256.Sum256(data)
			if got := fmt.Sprintf("in=%d sha=%s", len(data), hex.EncodeToString(sum[:])); got != want {
				t.Fatalf("%s: the input generators disagree:\n got %s\nwant %s", spec, got, want)
			}
			if strings.Contains(spec, "kind=gzip") && !bytes.HasPrefix(out, gzipOutputStreamHeader) {
				t.Errorf("%s: header % x", spec, out[:min(10, len(out))])
			}
			got, err := deflaterVectorDecompress(spec, out)
			if err != nil || !bytes.Equal(got, data) {
				t.Errorf("%s: decompresses to %d bytes (%v), want the %d input bytes", spec, len(got), err, len(data))
			}
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	if cases < 300 {
		t.Fatalf("only %d cases", cases)
	}
}
