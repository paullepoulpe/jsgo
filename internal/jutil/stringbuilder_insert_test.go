package jutil

import "testing"

// StringBuilder.insert(int, String) and insert(int, char) shift the tail in
// place; the results are those of java.lang.StringBuilder.
func TestStringBuilderInsert(t *testing.T) {
	b := NewStringBuilderFrom("abcdef")
	b.InsertChar(3, '\n')
	b.Insert(0, "<")
	b.Insert(b.Length(), ">")
	b.Insert(2, "é\U0001F600")
	b.InsertChar(1, 0xD800)
	want := "<" + StringFromChar(0xD800) + "aé\U0001F600bc\ndef>"
	if got := b.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if b.Length() != 13 {
		t.Errorf("length %d", b.Length())
	}
	for _, off := range []int{-1, 14} {
		func() {
			defer func() {
				if _, ok := recover().(*StringIndexOutOfBoundsException); !ok {
					t.Errorf("insert at %d: want StringIndexOutOfBoundsException", off)
				}
			}()
			b.InsertChar(off, 'x')
		}()
	}
}
