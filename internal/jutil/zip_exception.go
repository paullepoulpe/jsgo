package jutil

// ZipException is java.util.zip.ZipException: an I/O error reading a ZIP
// archive.
type ZipException struct{ IOException }

func (*ZipException) zipException() {}

// IZipException is the catch marker for ZipException.
type IZipException interface {
	IIOException
	zipException()
}

// NewZipException is new ZipException(String).
func NewZipException(message string) *ZipException {
	e := &ZipException{}
	InitThrowable(e, "java.util.zip.ZipException", &message, nil)
	return e
}
